// The top page: the document it shows, as a book of the demo viewer
// (examples/viewer/book.ts) by default: two facing pages, the first alone
// like a cover, turned by the buttons, the keys, a tap or a pull on a corner.
// The menu shows it one page at a time instead, or as pages to scroll
// through. Where two pages would be too small (a phone), it opens one page at
// a time. With a search box, zoom buttons and the page buttons.
import { dcValues, type View } from "@bdf/core";
import { BdfWorkerClient, buildTextLayer, installCopyHandler, RUN_ATTR, TEXT_LAYER_CSS, type HitRect } from "@bdf/render";
import { Book, type BookLayout } from "../examples/viewer/book.js";

type Layout = BookLayout | "scroll";

const $ = <T extends HTMLElement>(id: string) => document.getElementById(id) as T;
const stage = $("demo");
const q = $<HTMLInputElement>("q"), hits = $("hits"), pageNum = $("pageNum");
const prevPage = $<HTMLButtonElement>("prevPage"), nextPage = $<HTMLButtonElement>("nextPage");
const layoutSelect = $<HTMLSelectElement>("layout");

const style = document.createElement("style");
style.textContent = TEXT_LAYER_CSS;
document.head.append(style);

const client = new BdfWorkerClient(new Worker(new URL("lib/worker.js", location.href).href, { type: "module" }));
installCopyHandler(stage);

let view: View | undefined;
let lang = "";
let layout: Layout = "spread";
/** The reader picked the layout: it no longer follows the size of the stage. */
let picked = false;
let book: Book | undefined;
/** The page shown first in the scrolled layout, and CSS px per unit there. */
let scrolled = { page: 0, scale: 1 };
/** Bumped as the pages are laid out anew: work started for the ones before is dropped. */
let generation = 0;
/** Times the size that fits the stage. */
let zoom = 1;
/** The hits of the last search, their rectangles, and the one shown. */
let query = "";
let found: HitRect[][] = [];
let hit = -1;

const dpr = () => window.devicePixelRatio || 1;

/** The hits on a page, over its text. */
function hitLayer(index: number, scale: number): HTMLDivElement {
  const layer = document.createElement("div");
  layer.className = "hits";
  found.forEach((rects, i) => {
    for (const r of rects) {
      if (r.a !== index) continue;
      const d = document.createElement("div");
      d.className = i === hit ? "hit current" : "hit";
      d.style.cssText = `left: ${r.x * scale}px; top: ${r.y * scale}px; width: ${r.w * scale}px; height: ${r.h * scale}px`;
      layer.append(d);
    }
  });
  return layer;
}

/** Two facing pages, unless each would be much smaller than one page alone. */
function fitting(v: View): BookLayout {
  const pages = v.pages ?? [];
  const w = pages.reduce((m, p) => Math.max(m, p.w), 1), h = pages.reduce((m, p) => Math.max(m, p.h), 1);
  const cw = stage.clientWidth - 48, ch = stage.clientHeight - 48;
  const one = Math.min(cw / w, ch / h), two = Math.min(cw / (2 * w), ch / h);
  return pages.length > 1 && two >= 0.7 * one ? "spread" : "single";
}

/** The first page shown. */
function current(): number {
  return book ? book.pages[0] ?? 0 : scrolled.page;
}

/** Lay the pages out again (another layout, a zoom, a stage of another size), at the page given. */
function show(start = current()) {
  const v = view;
  if (!v) return;
  generation++;
  book?.destroy();
  book = undefined;
  stage.replaceChildren();
  stage.scrollTo(0, 0);
  if (!picked) layoutSelect.value = layout = fitting(v);
  if (layout === "scroll") showScroll(v, start);
  else showBook(v, start, layout);
}

function showBook(v: View, start: number, layout: BookLayout) {
  const pages = v.pages ?? [];
  const b: Book = book = new Book(stage, {
    pages,
    layout,
    rtl: v.direction === "rtl",
    zoom,
    animate: !matchMedia("(prefers-reduced-motion: reduce)").matches,
    start,
    label: (i) => `Page ${i + 1} of ${pages.length}`,
    bitmap: (i, scale) => client.page(v.id, i, scale),
    content: (i) => client.content(v.id, i),
    layers: (i, content, scale) => [buildTextLayer(content, scale, { lang }), hitLayer(i, scale)],
    pending: () => false,
    failed: () => false,
    onTurn: (shown) => turned(shown),
    onError: (e) => console.error(e),
  });
  turned(b.pages);
}

/**
 * The pages one above the other, as wide as the stage (times the zoom).
 * Each is drawn as it comes near the visible area.
 */
function showScroll(v: View, start: number) {
  const gen = generation, pages = v.pages ?? [];
  const live = () => gen === generation;
  const w = pages.reduce((m, p) => Math.max(m, p.w), 1);
  const scale = scrolled.scale = Math.min(1.5, Math.max(120, stage.clientWidth - 48) / w) * zoom;
  const near = (margin: string, render: (el: HTMLElement, i: number) => void) => {
    const o = new IntersectionObserver((entries) => {
      for (const e of entries) {
        if (!e.isIntersecting) continue;
        o.unobserve(e.target);
        render(e.target as HTMLElement, Number((e.target as HTMLElement).dataset.index));
      }
    }, { root: stage, rootMargin: margin });
    return o;
  };
  const bitmaps = near("400px", (el, i) => {
    client.page(v.id, i, scale * dpr()).then((bmp) => {
      if (!live()) return bmp.close();
      const canvas = document.createElement("canvas");
      canvas.width = bmp.width;
      canvas.height = bmp.height;
      canvas.style.width = `${pages[i].w * scale}px`;
      canvas.style.height = `${pages[i].h * scale}px`;
      canvas.setAttribute("aria-hidden", "true");
      canvas.getContext("2d")!.drawImage(bmp, 0, 0);
      bmp.close();
      el.prepend(canvas);
    }).catch((e) => console.error(e));
  });
  // text from further away, which screen readers read ahead of the pictures
  const texts = near("200% 0px", (el, i) => {
    client.content(v.id, i).then((content) => {
      if (live()) el.append(buildTextLayer(content, scale, { lang }), hitLayer(i, scale));
    }).catch((e) => console.error(e));
  });
  const inView = new Set<number>();
  const visible = new IntersectionObserver((entries) => {
    for (const e of entries) {
      const i = Number((e.target as HTMLElement).dataset.index);
      if (e.isIntersecting) inView.add(i); else inView.delete(i);
    }
    if (!live() || !inView.size) return;
    scrolled.page = Math.min(...inView);
    turned([scrolled.page]);
  }, { root: stage });
  const list = document.createElement("div");
  list.className = "scroll";
  pages.forEach((p, i) => {
    const el = document.createElement("div");
    el.className = "page";
    el.dataset.index = String(i);
    el.setAttribute("role", "group");
    el.setAttribute("aria-label", `Page ${i + 1} of ${pages.length}`);
    el.tabIndex = -1; // target of page links
    el.style.width = `${p.w * scale}px`;
    el.style.height = `${p.h * scale}px`;
    list.append(el);
    bitmaps.observe(el);
    texts.observe(el);
    visible.observe(el);
  });
  stage.append(list);
  scrolled.page = start;
  if (start > 0) scrollToPage(start);
  turned([start]);
}

function pageEl(index: number): HTMLElement | null {
  return stage.querySelector<HTMLElement>(`.page[data-index="${index}"]`);
}

function scrollToPage(index: number, focus = false) {
  const el = pageEl(index);
  if (!el) return;
  stage.scrollTo({ top: el.offsetTop - 12 });
  if (focus) el.focus({ preventScroll: true });
}

/** Go to a page: turn to it, or scroll to it. */
function go(index: number, focus = false) {
  if (book) book.go(index, true, focus);
  else scrollToPage(index, focus);
}

/** The page number between the page buttons. */
function turned(shown: number[]) {
  const n = view?.pages?.length ?? 0;
  const range = shown.length > 1 ? `${shown[0] + 1}–${shown[shown.length - 1] + 1}` : `${(shown[0] ?? 0) + 1}`;
  pageNum.textContent = `${range} / ${n}`;
  prevPage.disabled = book ? !book.canTurn(-1) : shown[0] <= 0;
  nextPage.disabled = book ? !book.canTurn(1) : shown[shown.length - 1] >= n - 1;
}

function turnBy(delta: 1 | -1) {
  if (book) book.turnBy(delta);
  else scrollToPage(Math.min((view?.pages?.length ?? 1) - 1, Math.max(0, scrolled.page + delta)));
}

function count() {
  hits.textContent = !q.value ? "" : found.length ? `${hit + 1} / ${found.length}` : "0";
}

/** Search, or with the same words go on to the next match (step -1: the one before). */
async function search(step: number) {
  const v = view;
  if (!v) return;
  if (q.value !== query) {
    const words = q.value;
    const res = words ? await client.search(v.id, words, { limit: 500 }) : [];
    const rects = res.length ? await client.locate(v.id, res) : [];
    if (v !== view) return;
    query = words;
    found = rects;
    hit = step < 0 ? 0 : -1;
  }
  if (found.length) hit = (hit + step + found.length) % found.length;
  count();
  const scale = book?.scale ?? scrolled.scale;
  for (const el of stage.querySelectorAll<HTMLElement>(".page[data-index]")) {
    el.querySelector(".hits")?.replaceWith(hitLayer(Number(el.dataset.index), scale));
  }
  const r = found[hit]?.[0];
  if (!r) return;
  if (book) book.go(r.a, true);
  else {
    const el = pageEl(r.a);
    if (el) stage.scrollTo({ top: Math.max(0, el.offsetTop + r.y * scale - stage.clientHeight / 3) });
  }
}

/** Zoom in or out: 1 fits the pages to the stage; scrolled pages go smaller too. */
function setZoom(z: number) {
  zoom = Math.min(4, Math.max(layout === "scroll" ? 0.5 : 1, z));
  show();
}

q.onkeydown = (e) => { if (e.key === "Enter") search(e.shiftKey ? -1 : 1); };
q.oninput = () => { if (!q.value) search(1); };
$("next").onclick = () => search(1);
$("prev").onclick = () => search(-1);
$("zoomIn").onclick = () => setZoom(zoom * 1.25);
$("zoomOut").onclick = () => setZoom(zoom / 1.25);
prevPage.onclick = () => turnBy(-1);
nextPage.onclick = () => turnBy(1);
layoutSelect.onchange = () => {
  picked = true;
  layout = layoutSelect.value as Layout;
  if (layout !== "scroll") zoom = Math.max(1, zoom);
  show();
};

// links to a page ("#page=N") go to it
stage.addEventListener("click", (e) => {
  const a = (e.target as Element).closest(`a[${RUN_ATTR.page}]`);
  if (!a || !view) return;
  e.preventDefault();
  go(Number(a.getAttribute(RUN_ATTR.page)) - 1, true);
});

// the book fits itself to the stage, but opens one page or two as it grows or
// shrinks; scrolled pages are as wide as the stage
let width = 0, timer: ReturnType<typeof setTimeout> | undefined;
new ResizeObserver(() => {
  if (!view) return;
  const changed = stage.clientWidth !== width;
  width = stage.clientWidth;
  if (layout === "scroll" ? changed : !picked && fitting(view) !== layout) {
    clearTimeout(timer);
    timer = setTimeout(() => show(), 100);
  }
}).observe(stage);

function message(text: string) {
  const p = document.createElement("p");
  p.setAttribute("role", "status");
  p.textContent = text;
  stage.replaceChildren(p);
}

async function open(src: string) {
  message("loading…");
  const manifest = await client.open({ kind: "single", url: new URL(src, location.href).href, range: true });
  lang = dcValues(manifest.meta?.dc?.language)[0] ?? "";
  view = manifest.views[0];
  width = stage.clientWidth;
  show(0);
}

open(stage.dataset.src!).catch((e) => {
  message(`The document could not be opened: ${(e as Error).message ?? e}`);
  console.error(e);
});
