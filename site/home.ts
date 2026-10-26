// The top page: the document it shows, as a book of the demo viewer
// (examples/viewer/book.ts) by default: two facing pages, the first alone
// like a cover, turned by the buttons, the keys, a tap or a pull on a corner.
// The menu shows it one page at a time instead, or as pages to scroll
// through. Where two pages would be too small (a phone), it opens one page at
// a time. With a search box, zoom buttons and the page buttons.
//
// Under the pages, the way back out of BDF: the same document written as a
// PDF file and as Markdown, its text for a search index, and that text
// searched with MiniSearch, all inside the browser by the preview module
// (cmd/bdfwasm built with -tags previewonly), as the site's tool pages do.
import MiniSearch from "minisearch";
import { dcValues, normalizeQuery, type View } from "@bdfkit/core";
import { BdfWorkerClient, buildTextLayer, findKey, installCopyHandler, RUN_ATTR, TEXT_LAYER_CSS, type HitRect } from "@bdfkit/render";
import { ConverterClient } from "../examples/common/convert.js";
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

q.onkeydown = (e) => {
  if (e.key === "Enter") search(e.shiftKey ? -1 : 1);
  else if (e.key === "Escape" && !e.isComposing && e.keyCode !== 229) {
    // the search ends: back to the pages (not the Escape that drops the text of an input method)
    e.preventDefault();
    q.value = "";
    search(1);
    stage.focus({ preventScroll: true });
  }
};
q.oninput = () => { if (!q.value) search(1); };
// While the focus is in the viewer (its pages or its tools), the find shortcut (Cmd+F, Ctrl+F) goes to its
// search box, which searches the whole document; the browser's own search finds the pages shown only. Elsewhere
// on the page, and pressed again in the box, the shortcut is the browser's.
stage.parentElement!.addEventListener("keydown", (e) => {
  const key = findKey(e);
  if (!key || e.defaultPrevented || !view) return;
  if (key === "open") {
    if (document.activeElement === q) return;
    e.preventDefault();
    q.focus({ preventScroll: true });
    q.select();
  } else if (q.value) {
    e.preventDefault();
    search(key === "next" ? 1 : -1);
  }
});
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

// --- The way back out of BDF -------------------------------------------------

/** The version of MiniSearch bundled (set by the build). */
declare const MINISEARCH_VERSION: string;

/** The text of a document for a search index, as bdf text writes it. */
interface SearchText {
  meta: { source?: string; dc?: Record<string, unknown> };
  views: { id: string; kind: string; title?: string; pages: { page?: number; text: string }[] }[];
}

const ja = document.documentElement.lang === "ja";
/** What the tools say while they work and when they are done. */
const L = ja ? {
  fetching: "demo.bdf を取得中…", writing: "書き出し中…", extracting: "テキストを取り出し中…",
  pdf: (pages: number, fonts: number, size: string, ms: number) => `${pages} ページ、フォント ${fonts} 個、${size}。${ms} ms で書き出しました`,
  md: (images: number, tables: number, size: string, ms: number) => `${size}、画像 ${images} 個${tables ? `、表 ${tables} 個` : ""}。${ms} ms で書き出しました`,
  text: (pages: number, chars: number, ms: number) => `${pages} ページ、${chars.toLocaleString("ja")} 文字。${ms} ms で取り出しました`,
  warnings: (n: number) => `。警告 ${n} 件（コンソールに出力）`,
  indexed: (n: number, ms: number) => `${n} 件（1 ページ 1 件）を MiniSearch ${MINISEARCH_VERSION} に ${ms} ms で索引づけ`,
  hits: (n: number, ms: number) => `${n} 件、${ms} ms`, none: "ヒットなし",
  page: (n: number) => `${n} ページ`, copied: "コピーしました",
  error: (m: string) => `エラー: ${m}`,
} : {
  fetching: "fetching demo.bdf…", writing: "writing…", extracting: "extracting the text…",
  pdf: (pages: number, fonts: number, size: string, ms: number) => `${pages} pages, ${fonts} fonts, ${size}, written in ${ms} ms`,
  md: (images: number, tables: number, size: string, ms: number) => `${size}, ${images} ${images === 1 ? "image" : "images"}${tables ? `, ${tables} ${tables === 1 ? "table" : "tables"}` : ""}, written in ${ms} ms`,
  text: (pages: number, chars: number, ms: number) => `${pages} pages, ${chars.toLocaleString("en")} characters, extracted in ${ms} ms`,
  warnings: (n: number) => `, ${n} ${n === 1 ? "warning" : "warnings"} (see the console)`,
  indexed: (n: number, ms: number) => `${n} records (one a page) indexed by MiniSearch ${MINISEARCH_VERSION} in ${ms} ms`,
  hits: (n: number, ms: number) => `${n} ${n === 1 ? "hit" : "hits"} in ${ms} ms`, none: "no hits",
  page: (n: number) => `Page ${n}`, copied: "copied",
  error: (m: string) => `error: ${m}`,
};

const kb = (n: number) => (n < 1024 ? `${n} B` : n < 1024 * 1024 ? `${(n / 1024).toFixed(1)} KB` : `${(n / 1024 / 1024).toFixed(1)} MB`);
const ms = (t0: number) => Math.round(performance.now() - t0);

function setStatus(el: HTMLElement, text: string, error = false) {
  el.textContent = text;
  el.toggleAttribute("data-error", error);
}
function fail(el: HTMLElement, e: unknown) {
  console.error("bdf:", e);
  setStatus(el, L.error((e as Error).message ?? String(e)), true);
}

// the tabs of the tools: one panel shown at a time, the arrow keys move between them
const tabs = [...document.querySelectorAll<HTMLButtonElement>(".back [role=tab]")];
function selectTab(tab: HTMLButtonElement) {
  for (const t of tabs) {
    const on = t === tab;
    t.setAttribute("aria-selected", String(on));
    t.tabIndex = on ? 0 : -1;
    document.getElementById(t.getAttribute("aria-controls")!)!.hidden = !on;
  }
}
tabs.forEach((t, i) => {
  t.onclick = () => selectTab(t);
  t.onkeydown = (e) => {
    const d = e.key === "ArrowRight" ? 1 : e.key === "ArrowLeft" ? -1 : 0;
    if (!d) return;
    e.preventDefault();
    const next = tabs[(i + d + tabs.length) % tabs.length];
    selectTab(next);
    next.focus();
  };
});

const PREVIEW = new URL("lib/bdf-preview.wasm", location.href).href;
const FONTS = new URL("fonts/", location.href).href;
let converter: ConverterClient | undefined;
/** The converter worker, started the first time a tool asks for it. */
const converterClient = () => (converter ??= new ConverterClient(new Worker(new URL("lib/convert-worker.js", location.href).href)));
let bytes: Promise<ArrayBuffer> | undefined;
/** The document above as one file, fetched whole the first time a tool asks (the viewer reads it by ranges). */
const demoBytes = () => (bytes ??= fetch(new URL(stage.dataset.src!, location.href)).then((r) => (r.ok ? r.arrayBuffer() : Promise.reject(new Error(`demo.bdf: HTTP ${r.status}`)))));

/** An object URL for a download link or a frame; the one before is revoked. */
const urls = new Map<string, string>();
function objectURL(key: string, file: Blob): string {
  const old = urls.get(key);
  if (old) URL.revokeObjectURL(old);
  const url = URL.createObjectURL(file);
  urls.set(key, url);
  return url;
}
function offer(a: HTMLAnchorElement, url: string, name: string) {
  a.href = url;
  a.download = name;
}
function copyOf(button: HTMLButtonElement, area: HTMLTextAreaElement, status: HTMLElement) {
  button.onclick = () => {
    const was = status.textContent ?? "";
    navigator.clipboard.writeText(area.value).then(() => {
      setStatus(status, L.copied);
      setTimeout(() => { if (status.textContent === L.copied) setStatus(status, was); }, 1500);
    }, (e) => fail(status, e));
  };
}

// a PDF file of the pages: the browser's own PDF viewer shows it
{
  const button = $<HTMLButtonElement>("pdfWrite"), status = $("pdfStatus"), frame = $<HTMLIFrameElement>("pdfFrame");
  button.onclick = async () => {
    button.disabled = true;
    setStatus(status, bytes ? L.writing : L.fetching);
    try {
      const data = await demoBytes();
      setStatus(status, L.writing);
      const t0 = performance.now();
      const res = await converterClient().pdf(PREVIEW, data.slice(0), { fonts: FONTS, tagged: $<HTMLInputElement>("pdfTagged").checked });
      const file = new Blob([res.pdf as BlobPart], { type: "application/pdf" });
      const url = objectURL("pdf", file);
      $<HTMLAnchorElement>("pdfOpen").href = url;
      offer($<HTMLAnchorElement>("pdfSave"), url, "demo-from-bdf.pdf");
      frame.src = url;
      frame.hidden = false;
      $("pdfLinks").hidden = false;
      for (const w of res.warnings) console.warn(`bdf: pdf: ${w}`);
      setStatus(status, L.pdf(res.pages, res.fonts, kb(file.size), ms(t0)) + (res.warnings.length ? L.warnings(res.warnings.length) : ""));
    } catch (e) {
      fail(status, e);
    } finally {
      button.disabled = false;
    }
  };
}

// the document as Markdown: the text, with its pictures in a zip beside it, embedded, or left out
{
  const button = $<HTMLButtonElement>("mdWrite"), status = $("mdStatus"), out = $<HTMLTextAreaElement>("mdOut");
  const images = $<HTMLSelectElement>("mdImages");
  copyOf($<HTMLButtonElement>("mdCopy"), out, status);
  button.onclick = async () => {
    button.disabled = true;
    setStatus(status, bytes ? L.writing : L.fetching);
    try {
      const data = await demoBytes();
      setStatus(status, L.writing);
      const t0 = performance.now();
      const how = images.value as "zip" | "embed" | "none";
      const res = await converterClient().markdown(PREVIEW, data.slice(0), { images: how === "zip" ? "files" : how, zip: how === "zip", name: "demo" });
      out.value = new TextDecoder().decode(res.markdown);
      // the text alone when it refers to no file
      const zipped = how === "zip" && res.images.length > 0 && res.zip;
      const file = zipped ? new Blob([res.zip as BlobPart], { type: "application/zip" }) : new Blob([res.markdown as BlobPart], { type: "text/markdown" });
      const save = $<HTMLAnchorElement>("mdSave");
      offer(save, objectURL("md", file), zipped ? "demo.zip" : "demo.md");
      save.textContent = `${save.dataset.label} (${zipped ? "demo.zip" : "demo.md"}, ${kb(file.size)})`;
      $("mdLinks").hidden = false;
      for (const w of res.warnings) console.warn(`bdf: markdown: ${w}`);
      setStatus(status, L.md(res.images.length, res.tables, kb(res.markdown.length), ms(t0)) + (res.warnings.length ? L.warnings(res.warnings.length) : ""));
    } catch (e) {
      fail(status, e);
    } finally {
      button.disabled = false;
    }
  };
}

/** The text of the document for a search index, extracted once for the text tool and the search. */
let extracted: Promise<{ json: string; text: SearchText; ms: number }> | undefined;
function extract(status: HTMLElement) {
  return (extracted ??= (async () => {
    setStatus(status, bytes ? L.extracting : L.fetching);
    const data = await demoBytes();
    setStatus(status, L.extracting);
    const t0 = performance.now();
    const res = await converterClient().text(PREVIEW, data.slice(0), {});
    return { json: res.json, text: JSON.parse(res.json) as SearchText, ms: ms(t0) };
  })().catch((e) => {
    extracted = undefined; // asked again next time
    throw e;
  }));
}

/** The text alone: the pages one after another, under the title of their view when there are several. */
function plain(text: SearchText): string {
  const parts: string[] = [];
  for (const v of text.views) {
    for (const p of v.pages) {
      const head = [text.views.length > 1 ? v.title || v.id : "", p.page ? `page ${p.page}` : ""].filter(Boolean).join(", ");
      parts.push(`${head ? `— ${head} —\n` : ""}${p.text}`);
    }
  }
  return parts.join("\n\n");
}

// the text for a search index, as JSON or as text alone
{
  const button = $<HTMLButtonElement>("textExtract"), status = $("textStatus"), out = $<HTMLTextAreaElement>("textOut");
  const shape = $<HTMLSelectElement>("textShape");
  copyOf($<HTMLButtonElement>("textCopy"), out, status);
  const show = async () => {
    const { json, text, ms: took } = await extract(status);
    const asJSON = shape.value === "json";
    out.value = asJSON ? json : plain(text);
    const file = new Blob([out.value], { type: asJSON ? "application/json" : "text/plain" });
    const save = $<HTMLAnchorElement>("textSave");
    offer(save, objectURL("text", file), asJSON ? "demo-text.json" : "demo-text.txt");
    save.textContent = `${save.dataset.label} (${kb(file.size)})`;
    $("textLinks").hidden = false;
    const pages = text.views.reduce((n, v) => n + v.pages.length, 0);
    const chars = text.views.reduce((n, v) => v.pages.reduce((n, p) => n + p.text.length, n), 0);
    setStatus(status, L.text(pages, chars, took));
  };
  button.onclick = () => {
    button.disabled = true;
    show().catch((e) => fail(status, e)).finally(() => { button.disabled = false; });
  };
  shape.onchange = () => { if (extracted) show().catch((e) => fail(status, e)); };
}

/** Runs of Chinese, Japanese and Korean characters, which no space divides into words. */
const CJK = /[\p{scx=Han}\p{scx=Hiragana}\p{scx=Katakana}\p{scx=Hangul}]+/gu;

/**
 * The terms of a text for the index and for a query alike: its words, and
 * the pairs of neighboring characters (bigrams) of its runs of Chinese,
 * Japanese and Korean, the usual way to index those scripts without a
 * dictionary (the browser's own segmenter knows the words of its dictionary
 * only, and cuts ぶらうざ into four characters).
 */
function tokenize(text: string): string[] {
  const out: string[] = [];
  for (const word of text.normalize("NFKC").split(/[\s\p{P}\p{S}]+/u)) {
    let at = 0;
    for (const m of word.matchAll(CJK)) {
      if (m.index > at) out.push(word.slice(at, m.index));
      const run = [...m[0]];
      if (run.length === 1) out.push(run[0]);
      for (let i = 1; i < run.length; i++) out.push(run[i - 1] + run[i]);
      at = m.index + m[0].length;
    }
    if (at < word.length) out.push(word.slice(at));
  }
  return out;
}

// the same records searched by MiniSearch, in the browser: a record a page, each term normalized as the viewer's
// search normalizes (NFKC, case and kana folded)
{
  interface Rec { id: string; page: number; title: string; text: string }
  const box = $<HTMLInputElement>("findQ"), status = $("findStatus"), list = $<HTMLOListElement>("findHits");
  const processTerm = (term: string) => normalizeQuery(term) || null;
  let index: Promise<{ engine: MiniSearch<Rec>; records: Map<string, Rec> }> | undefined;
  /** What the status says once the index is built, which an empty search shows again. */
  let built = "";
  const build = () => (index ??= extract(status).then(({ text }) => {
    const t0 = performance.now();
    const records = new Map<string, Rec>();
    for (const v of text.views) {
      for (const p of v.pages) {
        const id = `${v.id}#${p.page ?? 0}`;
        records.set(id, { id, page: p.page ?? 0, title: v.title ?? "", text: p.text });
      }
    }
    const engine = new MiniSearch<Rec>({ fields: ["text", "title"], storeFields: ["page"], tokenize, processTerm });
    engine.addAll([...records.values()]);
    setStatus(status, built = L.indexed(records.size, ms(t0)));
    return { engine, records };
  }).catch((e) => {
    index = undefined;
    throw e;
  }));

  /** The text around the first of the terms matched, the terms marked. */
  function snippet(rec: Rec, terms: string[]): DocumentFragment {
    const text = rec.text.normalize("NFKC").replace(/\s+/g, " ");
    const folded = normalizeQuery(text);
    const same = folded.length === text.length; // folding kept each character's place
    let at = -1;
    for (const t of terms) {
      const i = folded.indexOf(t);
      if (i >= 0 && (at < 0 || i < at)) at = i;
    }
    const start = Math.max(0, (at < 0 || !same ? 0 : at) - 60);
    const end = Math.min(text.length, start + 180);
    const frag = document.createDocumentFragment();
    if (start > 0) frag.append("…");
    let pos = start;
    if (same) {
      // every term, where it occurs in the window; the pairs of characters of a Japanese word overlap, and are one mark
      const found: [number, number][] = [];
      for (const t of terms) {
        for (let i = folded.indexOf(t, start); i >= 0 && i < end; i = folded.indexOf(t, i + 1)) found.push([i, i + t.length]);
      }
      found.sort((a, b) => a[0] - b[0]);
      const marks: [number, number][] = [];
      for (const [a, b] of found) {
        const last = marks[marks.length - 1];
        if (last && a <= last[1]) last[1] = Math.max(last[1], b);
        else marks.push([a, b]);
      }
      for (const [a, b] of marks) {
        frag.append(text.slice(pos, a));
        const m = document.createElement("mark");
        m.textContent = text.slice(a, b);
        frag.append(m);
        pos = b;
      }
    }
    frag.append(text.slice(pos, end));
    if (end < text.length) frag.append("…");
    return frag;
  }

  async function find() {
    const words = box.value.trim();
    const { engine, records } = await build();
    if (words !== box.value.trim()) return; // typed on
    if (!words) {
      list.replaceChildren();
      setStatus(status, built);
      return;
    }
    const t0 = performance.now();
    const hits = engine.search(words, { prefix: true, fuzzy: 0.2, combineWith: "AND" });
    setStatus(status, hits.length ? L.hits(hits.length, ms(t0)) : L.none);
    list.replaceChildren(...hits.slice(0, 20).map((h) => {
      const rec = records.get(String(h.id))!;
      const li = document.createElement("li");
      const a = document.createElement("a");
      a.href = `#page=${rec.page}`;
      a.textContent = L.page(rec.page);
      a.onclick = (e) => {
        e.preventDefault();
        // the pages above turn to the hit, and the viewer's own search marks the words on them
        if (q.value !== words) {
          q.value = words;
          search(1).catch((err) => console.error(err));
        }
        go(rec.page - 1, true);
        stage.scrollIntoView({ block: "nearest", behavior: "smooth" });
      };
      const p = document.createElement("p");
      p.append(snippet(rec, h.terms));
      li.append(a, p);
      return li;
    }));
  }
  let timer: ReturnType<typeof setTimeout> | undefined;
  box.oninput = () => {
    clearTimeout(timer);
    timer = setTimeout(() => { find().catch((e) => fail(status, e)); }, 150);
  };
  box.onkeydown = (e) => { if (e.key === "Enter") { clearTimeout(timer); find().catch((err) => fail(status, err)); } };
  $("findGo").onclick = () => { clearTimeout(timer); find().catch((e) => fail(status, e)); };
  // the index is built as the tab opens, so that the first search is quick
  $("tab-find").addEventListener("click", () => { build().catch((e) => fail(status, e)); }, { once: true });
}
