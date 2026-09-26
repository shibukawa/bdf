// Demo viewer: everything is decoded and rendered in a worker; the main thread
// only places bitmaps and a selectable, accessible text layer.
import { dcValues, type Manifest, type View, type SearchHit, type TextContent } from "@bdf/core";
import { BdfWorkerClient, buildTextLayer, installCopyHandler, TEXT_LAYER_CSS, RUN_ATTR, type HitRect, type TextLayerOptions } from "@bdf/render";

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
/** Bumped by show(): work started for an earlier view or zoom is dropped. */
let generation = 0;

/** Search state for the current view. */
const found = { query: "", hits: [] as SearchHit[], rects: [] as HitRect[][], index: -1 };
/** Shown sheets: redraw (with the search highlights) and scroll a hit into view. */
const sheetViews = new WeakMap<View, { redraw: () => void; reveal: (r: { x: number; y: number }) => void }>();
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

// Links to a page ("#page=N") scroll there and move the focus to it.
stage.addEventListener("click", (e) => {
  const a = (e.target as Element).closest(`a[${RUN_ATTR.page}]`);
  if (!a) return;
  e.preventDefault();
  goToPage(Number(a.getAttribute(RUN_ATTR.page)) - 1);
});

/** Text layer options: the document's language (the UI around it is English). */
const layerOptions = (): TextLayerOptions => ({ lang: dcValues(manifest.meta?.dc?.language)[0] ?? "" });

/** A text layer, which screen readers read ahead of the bitmaps: built over a wider margin. */
const TEXT_MARGIN = "300% 0px";
const BITMAP_MARGIN = "400px";

/** Call render once for each element as it comes near the visible area. */
function onNear(margin: string, render: (el: HTMLDivElement) => void): IntersectionObserver {
  const observer = new IntersectionObserver((entries) => {
    for (const e of entries) {
      if (!e.isIntersecting) continue;
      observer.unobserve(e.target);
      render(e.target as HTMLDivElement);
    }
  }, { root: stage, rootMargin: margin });
  return observer;
}

/** Put a bitmap on a page or band element, under its text layer. */
function placeBitmap(el: HTMLElement, bmp: ImageBitmap, w: number, h: number) {
  const canvas = document.createElement("canvas");
  canvas.width = bmp.width;
  canvas.height = bmp.height;
  canvas.style.width = `${w * zoom}px`;
  canvas.style.height = `${h * zoom}px`;
  canvas.setAttribute("aria-hidden", "true");
  canvas.getContext("2d")!.drawImage(bmp, 0, 0);
  bmp.close();
  el.querySelector("canvas")?.remove();
  el.prepend(canvas);
}

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
  generation++;
  stage.onscroll = null;
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

/**
 * Pages: each a labelled group. Bitmaps are rendered when a page comes near
 * the visible area, text layers from further away, so that a screen reader
 * can read on: moving its cursor scrolls the pages, which builds the next
 * layers. The whole document is not loaded up front (split or range sources
 * would fetch every page).
 */
function showPages(v: View) {
  const gen = generation;
  const list = document.createElement("div");
  list.className = "pages";
  const pagesOf = v.pages ?? [];
  const noun = manifest.meta?.source === "pptx" ? "Slide" : "Page";
  const bitmaps = onNear(BITMAP_MARGIN, (el) => void renderPage(v, Number(el.dataset.index), el, gen));
  const texts = onNear(TEXT_MARGIN, (el) => void renderText(v, Number(el.dataset.index), el, gen));
  pagesOf.forEach((p, i) => {
    const el = document.createElement("div");
    el.className = "page";
    el.dataset.index = String(i);
    el.setAttribute("role", "group");
    el.setAttribute("aria-label", `${noun} ${i + 1} of ${pagesOf.length}`);
    el.tabIndex = -1; // target of page links
    el.style.width = `${p.w * zoom}px`;
    el.style.height = `${p.h * zoom}px`;
    list.appendChild(el);
    bitmaps.observe(el);
    texts.observe(el);
  });
  stage.appendChild(list);
}

async function renderPage(v: View, index: number, el: HTMLDivElement, gen: number) {
  const page = v.pages![index];
  const t0 = performance.now();
  const bmp = await client.page(v.id, index, zoom * dpr());
  if (gen !== generation) return bmp.close();
  placeBitmap(el, bmp, page.w, page.h);
  setTiming(`page ${index + 1} rendered in ${(performance.now() - t0).toFixed(0)} ms`);
}

async function renderText(v: View, index: number, el: HTMLDivElement, gen: number) {
  const content = await client.content(v.id, index);
  if (gen !== generation) return;
  el.append(buildTextLayer(content, zoom, layerOptions()), highlightLayer(index));
}

/** Follow a link to a page (0-based): scroll to it and focus it. */
function goToPage(index: number) {
  if (current.kind === "flow" && continuousBox.checked) {
    const { top, band } = continuousPosition(current, index);
    stage.scrollTo({ top: top * zoom });
    stage.querySelector<HTMLElement>(`.page[data-y="${band}"]`)?.focus({ preventScroll: true });
    return;
  }
  const el = stage.querySelector<HTMLElement>(`.page[data-index="${index}"]`);
  el?.scrollIntoView({ block: "start" });
  el?.focus({ preventScroll: true });
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
  sheetViews.get(current)?.redraw();
  // scroll to the current hit
  const rect = found.rects[found.index]?.[0];
  if (!rect) return;
  if (current.kind === "sheet") {
    sheetViews.get(current)?.reveal(rect);
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

const BAND = 800;

/** Continuous layout: stacked body rectangles. */
function continuousSize(v: View) {
  const gap = v.continuous?.gap ?? 0;
  let height = 0, width = 0;
  const tops: number[] = [];
  for (const p of v.pages ?? []) {
    const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
    tops.push(height);
    height += b.h + gap;
    width = Math.max(width, b.w);
  }
  return { width, height: Math.max(0, height - gap), tops };
}

/** Where a page starts in the continuous layout, and the band that holds it. */
function continuousPosition(v: View, index: number) {
  const top = continuousSize(v).tops[index] ?? 0;
  return { top, band: Math.floor(top / BAND) * BAND };
}

/** Continuous flow: body rectangles stacked, rendered in 800-unit bands on demand (text layers further ahead). */
function showContinuous(v: View) {
  const gen = generation;
  const { width, height } = continuousSize(v);
  const list = document.createElement("div");
  list.className = "pages";
  const viewportOf = (el: HTMLElement) => {
    const y = Number(el.dataset.y);
    return { x: 0, y, w: width, h: Math.min(BAND, height - y) };
  };
  const bitmaps = onNear(BITMAP_MARGIN, (el) => {
    const vp = viewportOf(el);
    void client.continuous(v.id, vp, zoom * dpr()).then((bmp) => (gen === generation ? placeBitmap(el, bmp, vp.w, vp.h) : bmp.close()));
  });
  const texts = onNear(TEXT_MARGIN, (el) => {
    void client.continuousContent(v.id, viewportOf(el)).then((c) => { if (gen === generation) el.append(buildTextLayer(c, zoom, layerOptions())); });
  });
  for (let y = 0; y < height; y += BAND) {
    const el = document.createElement("div");
    el.className = "page";
    el.dataset.y = String(y);
    el.tabIndex = -1; // target of page links
    el.style.width = `${width * zoom}px`;
    el.style.height = `${Math.min(BAND, height - y) * zoom}px`;
    el.style.boxShadow = "none";
    list.appendChild(el);
    bitmaps.observe(el);
    texts.observe(el);
  }
  list.style.gap = "0";
  stage.appendChild(list);
}

/** Row or column sizes of a sheet, from the manifest's run-length list. */
class SheetAxis {
  readonly count: number;
  readonly total: number;
  constructor(private runs: [number, number][] = []) {
    this.count = runs.reduce((a, [n]) => a + n, 0);
    this.total = runs.reduce((a, [n, s]) => a + n * s, 0);
  }
  /** Where entry i starts. */
  pos(i: number): number {
    let p = 0;
    for (const [n, s] of this.runs) {
      if (i <= n) return p + i * s;
      p += n * s;
      i -= n;
    }
    return p;
  }
  /** Call fn for the entries with a size that overlap [from, to). */
  each(from: number, to: number, fn: (i: number, start: number, size: number) => void) {
    let i = 0, p = 0;
    for (const [n, s] of this.runs) {
      if (s > 0 && p + n * s > from) {
        for (let k = Math.max(0, Math.floor((from - p) / s)); k < n; k++) {
          const start = p + k * s;
          if (start >= to) return;
          fn(i + k, start, s);
        }
      }
      p += n * s;
      i += n;
    }
  }
}

/** Column letters: A … Z, AA … */
function columnLabel(c: number): string {
  let s = "";
  for (c++; c > 0; c = Math.floor((c - 1) / 26)) s = String.fromCharCode(65 + ((c - 1) % 26)) + s;
  return s;
}

const SHEET_HEADER = 20;

/**
 * Sheet: a scroll area the size of the sheet with one sticky canvas that
 * shows the column and row headers and the visible region, split into the
 * frozen panes of the view (each rendered by the worker as a region of its
 * own). A text layer (a table of the cells, for screen readers, selection
 * and copy) covers the visible region and a screen around it; the cells of
 * the frozen panes are pinned with a translation that follows the scroll.
 */
function showSheet(v: View) {
  const cols = new SheetAxis(v.cols), rows = new SheetAxis(v.rows);
  const fc = Math.min(v.freeze?.cols ?? 0, cols.count), fr = Math.min(v.freeze?.rows ?? 0, rows.count);
  const fw = cols.pos(fc), fh = rows.pos(fr); // frozen extent in units
  const hw = Math.max(32, String(rows.count).length * 7 + 14), hh = SHEET_HEADER; // headers in CSS px
  const wrap = document.createElement("div");
  wrap.className = "sheet";
  wrap.style.width = `${hw + cols.total * zoom}px`;
  wrap.style.height = `${hh + rows.total * zoom}px`;
  const canvas = document.createElement("canvas");
  canvas.setAttribute("aria-hidden", "true");
  const layerHost = document.createElement("div");
  layerHost.className = "sheetText";
  layerHost.style.cssText = `left: ${hw}px; top: ${hh}px; width: ${cols.total * zoom}px; height: ${rows.total * zoom}px`;
  wrap.append(canvas, layerHost);
  stage.appendChild(wrap);

  const gen = generation;
  const sheet = { rows: rows.count, cols: cols.count, headerRows: fr, headerCols: fc };
  // cells of the frozen panes stay where they are while the sheet scrolls
  const pin = (span: HTMLSpanElement, r: { x: number; y: number }) => {
    const x = r.x < fw, y = r.y < fh;
    if (!x && !y) return;
    span.style.translate = `${x ? "var(--sx)" : "0px"} ${y ? "var(--sy)" : "0px"}`;
    span.style.zIndex = "1";
  };
  let covered: { x: number; y: number; w: number; h: number } | undefined;
  let textTimer: ReturnType<typeof setTimeout> | undefined;
  const updateText = (vp: { x: number; y: number; w: number; h: number }) => {
    const inside = covered && vp.x >= covered.x && vp.y >= covered.y && vp.x + vp.w <= covered.x + covered.w && vp.y + vp.h <= covered.y + covered.h;
    if (inside) return;
    clearTimeout(textTimer);
    textTimer = setTimeout(async () => {
      const region = { x: Math.max(0, vp.x - vp.w), y: Math.max(0, vp.y - vp.h), w: vp.w * 3, h: vp.h * 3 };
      // the frozen panes along the region
      const frozen = [{ x: region.x, y: 0, w: region.w, h: fh }, { x: 0, y: region.y, w: fw, h: region.h }, { x: 0, y: 0, w: fw, h: fh }];
      const content: TextContent = await client.sheetContent(v.id, [region, ...frozen]);
      if (gen !== generation) return;
      covered = region;
      const layer = buildTextLayer(content, zoom, { ...layerOptions(), sheet, onSpan: pin });
      layerHost.replaceChildren(layer);
    }, 150);
  };

  type Pane = { src: { x: number; y: number; w: number; h: number }; dx: number; dy: number };
  let panes: Pane[] = [];
  let bitmaps: ImageBitmap[] = [];
  const paint = () => {
    const d = dpr();
    const vw = Math.min(stage.clientWidth, hw + cols.total * zoom), vh = Math.min(stage.clientHeight, hh + rows.total * zoom);
    canvas.width = Math.ceil(vw * d);
    canvas.height = Math.ceil(vh * d);
    canvas.style.width = `${vw}px`;
    canvas.style.height = `${vh}px`;
    const ctx = canvas.getContext("2d")!;
    ctx.fillStyle = "#fff";
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    panes.forEach((p, i) => ctx.drawImage(bitmaps[i], p.dx * d, p.dy * d, p.src.w * zoom * d, p.src.h * zoom * d));
    // search hits, clipped to each pane
    ctx.save();
    ctx.scale(d, d);
    for (const p of panes) {
      ctx.save();
      ctx.beginPath();
      ctx.rect(p.dx, p.dy, p.src.w * zoom, p.src.h * zoom);
      ctx.clip();
      found.rects.forEach((rects, hi) => {
        ctx.fillStyle = hi === found.index ? "rgba(255,120,0,.5)" : "rgba(255,210,0,.45)";
        for (const r of rects) ctx.fillRect(p.dx + (r.x - p.src.x) * zoom, p.dy + (r.y - p.src.y) * zoom, r.w * zoom, r.h * zoom);
      });
      ctx.restore();
    }
    drawHeaders(ctx, vw, vh);
    ctx.restore();
  };
  const drawHeaders = (ctx: CanvasRenderingContext2D, vw: number, vh: number) => {
    const sx = stage.scrollLeft, sy = stage.scrollTop;
    ctx.fillStyle = "#f3f4f6";
    ctx.fillRect(0, 0, vw, hh);
    ctx.fillRect(0, 0, hw, vh);
    ctx.font = "11px system-ui, sans-serif";
    ctx.textAlign = "center";
    ctx.textBaseline = "middle";
    ctx.strokeStyle = "#c8ccd2";
    ctx.lineWidth = 1;
    const header = (x: number, y: number, w: number, h: number, label: string) => {
      ctx.fillStyle = "#444";
      ctx.fillText(label, x + w / 2, y + h / 2);
      ctx.strokeRect(Math.round(x) + 0.5, Math.round(y) + 0.5, Math.round(w), Math.round(h));
    };
    // frozen columns, then the scrolled ones clipped to their pane
    const colsIn = (from: number, to: number, shift: number, clipX: number) => {
      ctx.save();
      ctx.beginPath();
      ctx.rect(clipX, 0, vw - clipX, hh);
      ctx.clip();
      cols.each(from, to, (c, start, size) => header(hw + start * zoom - shift, 0, size * zoom, hh, columnLabel(c)));
      ctx.restore();
    };
    colsIn(0, fw, 0, hw);
    colsIn(fw + sx / zoom, fw + sx / zoom + (vw - hw) / zoom, sx, hw + fw * zoom);
    const rowsIn = (from: number, to: number, shift: number, clipY: number) => {
      ctx.save();
      ctx.beginPath();
      ctx.rect(0, clipY, hw, vh - clipY);
      ctx.clip();
      rows.each(from, to, (r, start, size) => header(0, hh + start * zoom - shift, hw, size * zoom, String(r + 1)));
      ctx.restore();
    };
    rowsIn(0, fh, 0, hh);
    rowsIn(fh + sy / zoom, fh + sy / zoom + (vh - hh) / zoom, sy, hh + fh * zoom);
    ctx.fillStyle = "#e5e7eb";
    ctx.fillRect(0, 0, hw, hh);
    // the edges of the frozen panes
    ctx.strokeStyle = "#8a9099";
    ctx.beginPath();
    if (fc > 0) { ctx.moveTo(hw + fw * zoom + 0.5, 0); ctx.lineTo(hw + fw * zoom + 0.5, vh); }
    if (fr > 0) { ctx.moveTo(0, hh + fh * zoom + 0.5); ctx.lineTo(vw, hh + fh * zoom + 0.5); }
    ctx.stroke();
  };

  let pending = false;
  const draw = async () => {
    if (pending) return;
    pending = true;
    await new Promise((r) => requestAnimationFrame(r));
    pending = false;
    const vw = stage.clientWidth, vh = stage.clientHeight;
    const sx = stage.scrollLeft, sy = stage.scrollTop;
    layerHost.style.setProperty("--sx", `${sx}px`);
    layerHost.style.setProperty("--sy", `${sy}px`);
    // the scrolled region starts where the frozen panes end
    const mx = fw + sx / zoom, my = fh + sy / zoom;
    const pw = Math.max(0, Math.min((vw - hw) / zoom - fw, cols.total - mx)), ph = Math.max(0, Math.min((vh - hh) / zoom - fh, rows.total - my));
    const next: Pane[] = [
      { src: { x: mx, y: my, w: pw, h: ph }, dx: hw + fw * zoom, dy: hh + fh * zoom },
      { src: { x: mx, y: 0, w: pw, h: fh }, dx: hw + fw * zoom, dy: hh },
      { src: { x: 0, y: my, w: fw, h: ph }, dx: hw, dy: hh + fh * zoom },
      { src: { x: 0, y: 0, w: fw, h: fh }, dx: hw, dy: hh },
    ].filter((p) => p.src.w > 0 && p.src.h > 0);
    const t0 = performance.now();
    const bmps = await Promise.all(next.map((p) => client.sheet(v.id, p.src, zoom * dpr())));
    if (gen !== generation) return bmps.forEach((b) => b.close());
    bitmaps.forEach((b) => b.close());
    panes = next;
    bitmaps = bmps;
    paint();
    updateText(next[0]?.src ?? { x: 0, y: 0, w: 1, h: 1 });
    setTiming(`sheet region ${Math.round(mx)},${Math.round(my)} rendered in ${(performance.now() - t0).toFixed(0)} ms`);
  };
  sheetViews.set(v, {
    redraw: () => { if (bitmaps.length) paint(); },
    reveal: (r) => {
      const pw = stage.clientWidth - hw - fw * zoom, ph = stage.clientHeight - hh - fh * zoom;
      stage.scrollTo({
        left: r.x < fw ? stage.scrollLeft : Math.max(0, (r.x - fw) * zoom - pw / 2),
        top: r.y < fh ? stage.scrollTop : Math.max(0, (r.y - fh) * zoom - ph / 2),
        behavior: "smooth",
      });
    },
  });
  stage.onscroll = () => { if (current === v) void draw(); };
  // the visible region changes with the window too
  const resize = new ResizeObserver(() => { if (current === v && gen === generation) void draw(); else resize.disconnect(); });
  resize.observe(stage);
  void draw();
}

main().catch((e) => setStatus(`error: ${e.message ?? e}`));
