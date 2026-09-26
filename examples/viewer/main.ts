// Demo viewer: everything is decoded and rendered in a worker; the main thread
// only places bitmaps and a selectable text layer.
import { dcValues, type Manifest, type View, type SearchHit } from "@bdf/core";
import { BdfWorkerClient, buildTextLayer, installCopyHandler, TEXT_LAYER_CSS, type HitRect } from "@bdf/render";

const params = new URLSearchParams(location.search);
const src = params.get("src") ?? "/testdata/demo.bdf";
const client = new BdfWorkerClient(new Worker("./worker.js", { type: "module" }));

const $ = <T extends HTMLElement>(id: string) => document.getElementById(id) as T;
const stage = $<HTMLDivElement>("stage");
const tabs = $<HTMLSpanElement>("tabs");
const zoomInput = $<HTMLInputElement>("zoom");
const continuousBox = $<HTMLInputElement>("continuous");
const status = $<HTMLSpanElement>("status");
const timing = $<HTMLSpanElement>("timing");
const hitsBox = $<HTMLSpanElement>("hits");

let manifest: Manifest;
let current: View;
let zoom = 1;

/** Search state for the current view. */
const found = { query: "", hits: [] as SearchHit[], rects: [] as HitRect[][], index: -1 };
const sheetCanvases = new WeakMap<View, { canvas: HTMLCanvasElement; viewport: () => { x: number; y: number } }>();
const dpr = () => window.devicePixelRatio || 1;

/** Announced to screen readers (a polite live region). */
function setStatus(s: string) { status.textContent = s; }
/** Render timings: shown only, as they change on every scroll. */
function setTiming(s: string) { timing.textContent = s; }

const style = document.createElement("style");
style.textContent = TEXT_LAYER_CSS;
document.head.appendChild(style);

// Copy puts the selected runs' text (with the document's spaces and line
// breaks) on the clipboard, across pages.
installCopyHandler(stage);

async function main() {
  const source = src.endsWith("/") ? { kind: "split" as const, base: new URL(src, location.href).href } : { kind: "single" as const, url: new URL(src, location.href).href, range: params.has("range") };
  setStatus("loading…");
  manifest = await client.open(source);
  document.title = `${dcValues(manifest.meta?.dc?.title)[0] ?? "BDF"} – viewer`;
  manifest.views.forEach((v, i) => {
    const b = document.createElement("button");
    b.textContent = `${v.title ?? v.id} (${v.kind})`;
    b.onclick = () => show(v);
    b.id = `tab-${i}`;
    b.dataset.id = v.id;
    b.setAttribute("role", "tab");
    b.setAttribute("aria-controls", stage.id);
    tabs.appendChild(b);
  });
  // tablist keys: arrows, Home and End move to a view and show it
  tabs.onkeydown = (e) => {
    const all = [...tabs.querySelectorAll<HTMLButtonElement>("[role=tab]")];
    const i = all.findIndex((b) => b.dataset.id === current.id);
    const keys: Record<string, number> = { ArrowRight: i + 1, ArrowLeft: i - 1, Home: 0, End: all.length - 1 };
    const to = keys[e.key];
    if (to === undefined) return;
    e.preventDefault();
    const b = all[(to + all.length) % all.length];
    b.focus();
    show(manifest.views.find((v) => v.id === b.dataset.id)!);
  };
  const q = $<HTMLInputElement>("q");
  q.onkeydown = (e) => { if (e.key === "Enter") void runSearch(q.value, e.shiftKey ? -1 : 1); };
  q.oninput = () => { if (!q.value) void runSearch("", 0); };
  $("next").onclick = () => void runSearch(q.value, 1);
  $("prev").onclick = () => void runSearch(q.value, -1);
  zoomInput.oninput = () => {
    zoom = Number(zoomInput.value);
    const pct = `${Math.round(zoom * 100)}%`;
    $("zoomv").textContent = pct;
    zoomInput.setAttribute("aria-valuetext", pct);
    show(current);
  };
  continuousBox.onchange = () => show(current);
  show(manifest.views[0]);
}

function show(v: View) {
  if (current !== v) {
    found.query = ""; found.hits = []; found.rects = []; found.index = -1; hitsBox.textContent = "";
    setStatus(describe(v));
  }
  current = v;
  for (const b of tabs.querySelectorAll<HTMLButtonElement>("[role=tab]")) {
    const selected = b.dataset.id === v.id;
    b.setAttribute("aria-selected", String(selected));
    b.tabIndex = selected ? 0 : -1;
    if (selected) stage.setAttribute("aria-labelledby", b.id);
  }
  $("modeBox").hidden = v.kind !== "flow";
  stage.replaceChildren();
  stage.scrollTop = 0;
  if (v.kind === "sheet") showSheet(v);
  else if (v.kind === "flow" && continuousBox.checked) showContinuous(v);
  else showPages(v);
}

/** What a view holds, for the status line. */
function describe(v: View): string {
  if (v.kind !== "sheet") {
    const n = v.pages?.length ?? 0;
    return `${n} ${n === 1 ? "page" : "pages"}`;
  }
  const count = (runs?: [number, number][]) => (runs ?? []).reduce((a, [n]) => a + n, 0);
  return `${count(v.rows)} rows, ${count(v.cols)} columns`;
}

/** Lazily render each page when it scrolls into view. */
function showPages(v: View) {
  const list = document.createElement("div");
  list.className = "pages";
  const observer = new IntersectionObserver((entries) => {
    for (const e of entries) {
      if (!e.isIntersecting) continue;
      observer.unobserve(e.target);
      void renderPage(v, Number((e.target as HTMLElement).dataset.index), e.target as HTMLDivElement);
    }
  }, { root: stage, rootMargin: "400px" });
  (v.pages ?? []).forEach((p, i) => {
    const el = document.createElement("div");
    el.className = "page";
    el.dataset.index = String(i);
    el.style.width = `${p.w * zoom}px`;
    el.style.height = `${p.h * zoom}px`;
    list.appendChild(el);
    observer.observe(el);
  });
  stage.appendChild(list);
}

async function renderPage(v: View, index: number, el: HTMLDivElement) {
  const page = v.pages![index];
  const t0 = performance.now();
  const [bmp, runs] = await Promise.all([client.page(v.id, index, zoom * dpr()), client.text(v.id, index)]);
  const canvas = document.createElement("canvas");
  canvas.width = bmp.width;
  canvas.height = bmp.height;
  canvas.style.width = `${page.w * zoom}px`;
  canvas.style.height = `${page.h * zoom}px`;
  canvas.getContext("2d")!.drawImage(bmp, 0, 0);
  bmp.close();
  el.replaceChildren(canvas, buildTextLayer(runs, zoom), highlightLayer(index));
  setTiming(`page ${index + 1} rendered in ${(performance.now() - t0).toFixed(0)} ms (${runs.length} text runs)`);
}

/** Search: ask the worker for hits, locate them once, then highlight and step through them. */
async function runSearch(query: string, step: number) {
  if (query !== found.query) {
    found.query = query;
    found.hits = query ? await client.search(current.id, query, { limit: 500, context: 30 }) : [];
    found.rects = found.hits.length ? await client.locate(current.id, found.hits) : [];
    found.index = -1;
  }
  if (found.hits.length) found.index = (found.index + step + found.hits.length) % found.hits.length;
  showHitCount(query);
  // refresh highlights on everything that is rendered
  for (const el of stage.querySelectorAll<HTMLDivElement>(".page[data-index]")) {
    const old = el.querySelector(".hlLayer");
    if (old) old.replaceWith(highlightLayer(Number(el.dataset.index)));
  }
  const sheet = sheetCanvases.get(current);
  if (sheet) drawSheetHighlights(sheet.canvas, sheet.viewport());
  // scroll to the current hit
  const rect = found.rects[found.index]?.[0];
  if (!rect) return;
  if (current.kind === "sheet") {
    stage.scrollTo({ left: Math.max(0, rect.x * zoom - stage.clientWidth / 2), top: Math.max(0, rect.y * zoom - stage.clientHeight / 2), behavior: "smooth" });
  } else {
    const el = stage.querySelector<HTMLDivElement>(`.page[data-index="${rect.a}"]`);
    el?.scrollIntoView({ block: "center", behavior: "smooth" });
  }
}

/**
 * "3 of 12" next to the search box. It is a live region, and the highlights
 * are only drawn, so screen readers also get the current hit's page and
 * its surrounding text.
 */
function showHitCount(query: string) {
  if (!query) { hitsBox.textContent = ""; return; }
  if (!found.hits.length) { hitsBox.textContent = "no matches"; return; }
  const hit = found.hits[found.index];
  const page = current.kind === "sheet" ? "" : `, page ${hit.segments[0].a + 1}`;
  const detail = document.createElement("span");
  detail.className = "sr-only";
  detail.textContent = `${page}: ${hit.context.replace(/ ⏎ /g, " ")}`;
  hitsBox.replaceChildren(`${found.index + 1} of ${found.hits.length}`, detail);
}

/** Highlight rectangles of the hits on one page. */
function highlightLayer(pageIndex: number): HTMLDivElement {
  const layer = document.createElement("div");
  layer.className = "textLayer hlLayer";
  found.rects.forEach((rects, hi) => {
    for (const r of rects) {
      if (r.a !== pageIndex) continue;
      const d = document.createElement("div");
      d.className = hi === found.index ? "hl current" : "hl";
      d.style.left = `${r.x * zoom}px`;
      d.style.top = `${r.y * zoom}px`;
      d.style.width = `${r.w * zoom}px`;
      d.style.height = `${r.h * zoom}px`;
      layer.appendChild(d);
    }
  });
  return layer;
}

/** Sheets: paint the hits over the bitmap of the visible region. */
function drawSheetHighlights(canvas: HTMLCanvasElement, viewport: { x: number; y: number }) {
  const ctx = canvas.getContext("2d")!;
  const s = zoom * dpr();
  ctx.save();
  ctx.setTransform(s, 0, 0, s, -viewport.x * s, -viewport.y * s);
  found.rects.forEach((rects, hi) => {
    ctx.fillStyle = hi === found.index ? "rgba(255,120,0,.5)" : "rgba(255,210,0,.45)";
    for (const r of rects) ctx.fillRect(r.x, r.y, r.w, r.h);
  });
  ctx.restore();
}

/** Continuous flow: body rectangles stacked, rendered in 800-unit bands on demand. */
function showContinuous(v: View) {
  const gap = v.continuous?.gap ?? 0;
  let height = 0, width = 0;
  for (const p of v.pages ?? []) {
    const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
    height += b.h + gap;
    width = Math.max(width, b.w);
  }
  height = Math.max(0, height - gap);
  const list = document.createElement("div");
  list.className = "pages";
  const band = 800;
  const observer = new IntersectionObserver((entries) => {
    for (const e of entries) {
      if (!e.isIntersecting) continue;
      observer.unobserve(e.target);
      const el = e.target as HTMLDivElement;
      const y = Number(el.dataset.y);
      const h = Math.min(band, height - y);
      const viewport = { x: 0, y, w: width, h };
      void Promise.all([client.continuous(v.id, viewport, zoom * dpr()), client.continuousText(v.id, viewport)]).then(([bmp, runs]) => {
        const canvas = document.createElement("canvas");
        canvas.width = bmp.width;
        canvas.height = bmp.height;
        canvas.style.width = `${width * zoom}px`;
        canvas.style.height = `${h * zoom}px`;
        canvas.getContext("2d")!.drawImage(bmp, 0, 0);
        bmp.close();
        el.replaceChildren(canvas, buildTextLayer(runs, zoom));
      });
    }
  }, { root: stage, rootMargin: "400px" });
  for (let y = 0; y < height; y += band) {
    const el = document.createElement("div");
    el.className = "page";
    el.dataset.y = String(y);
    el.style.width = `${width * zoom}px`;
    el.style.height = `${Math.min(band, height - y) * zoom}px`;
    el.style.boxShadow = "none";
    list.appendChild(el);
    observer.observe(el);
  }
  list.style.gap = "0";
  stage.appendChild(list);
}

/** Sheet: a scroll area the size of the sheet with one sticky canvas showing the visible region. */
function showSheet(v: View) {
  const sum = (runs?: [number, number][]) => (runs ?? []).reduce((a, [n, s]) => a + n * s, 0);
  const width = sum(v.cols), height = sum(v.rows);
  const wrap = document.createElement("div");
  wrap.className = "sheet";
  wrap.style.width = `${width * zoom}px`;
  wrap.style.height = `${height * zoom}px`;
  const canvas = document.createElement("canvas");
  wrap.appendChild(canvas);
  stage.appendChild(wrap);

  let last = { x: 0, y: 0 };
  sheetCanvases.set(v, { canvas, viewport: () => last });
  let pending = false;
  const draw = async () => {
    if (pending) return;
    pending = true;
    await new Promise((r) => requestAnimationFrame(r));
    pending = false;
    const vw = stage.clientWidth, vh = stage.clientHeight;
    const viewport = { x: stage.scrollLeft / zoom, y: stage.scrollTop / zoom, w: Math.min(vw, width * zoom - stage.scrollLeft) / zoom, h: Math.min(vh, height * zoom - stage.scrollTop) / zoom };
    const t0 = performance.now();
    const bmp = await client.sheet(v.id, viewport, zoom * dpr());
    canvas.width = bmp.width;
    canvas.height = bmp.height;
    canvas.style.width = `${viewport.w * zoom}px`;
    canvas.style.height = `${viewport.h * zoom}px`;
    canvas.getContext("2d")!.drawImage(bmp, 0, 0);
    bmp.close();
    last = viewport;
    drawSheetHighlights(canvas, viewport);
    setTiming(`sheet region ${Math.round(viewport.x)},${Math.round(viewport.y)} rendered in ${(performance.now() - t0).toFixed(0)} ms`);
  };
  stage.onscroll = () => { if (current === v) void draw(); };
  void draw();
}

main().catch((e) => setStatus(`error: ${e.message ?? e}`));
