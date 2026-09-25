// Demo viewer: everything is decoded and rendered in a worker; the main thread
// only places bitmaps and a selectable text layer.
import type { Manifest, View, TextRun } from "@bdf/core";
import { BdfWorkerClient } from "@bdf/render";
import { fontString } from "@bdf/render";

const params = new URLSearchParams(location.search);
const src = params.get("src") ?? "/fixtures/demo.bdf";
const client = new BdfWorkerClient(new Worker("./worker.js", { type: "module" }));

const $ = <T extends HTMLElement>(id: string) => document.getElementById(id) as T;
const stage = $<HTMLDivElement>("stage");
const tabs = $<HTMLSpanElement>("tabs");
const zoomInput = $<HTMLInputElement>("zoom");
const continuousBox = $<HTMLInputElement>("continuous");
const status = $<HTMLSpanElement>("status");

let manifest: Manifest;
let current: View;
let zoom = 1;
const dpr = () => window.devicePixelRatio || 1;

function setStatus(s: string) { status.textContent = s; }

async function main() {
  const source = src.endsWith("/") ? { kind: "split" as const, base: new URL(src, location.href).href } : { kind: "single" as const, url: new URL(src, location.href).href, range: params.has("range") };
  manifest = await client.open(source);
  document.title = `${manifest.meta?.title ?? "BDF"} – viewer`;
  for (const v of manifest.views) {
    const b = document.createElement("button");
    b.textContent = `${v.title ?? v.id} (${v.kind})`;
    b.onclick = () => show(v);
    b.dataset.id = v.id;
    tabs.appendChild(b);
  }
  zoomInput.oninput = () => { zoom = Number(zoomInput.value); $("zoomv").textContent = `${Math.round(zoom * 100)}%`; show(current); };
  continuousBox.onchange = () => show(current);
  show(manifest.views[0]);
}

function show(v: View) {
  current = v;
  for (const b of tabs.querySelectorAll("button")) b.setAttribute("aria-selected", String(b.dataset.id === v.id));
  $("modeBox").hidden = v.kind !== "flow";
  stage.replaceChildren();
  stage.scrollTop = 0;
  if (v.kind === "sheet") showSheet(v);
  else if (v.kind === "flow" && continuousBox.checked) showContinuous(v);
  else showPages(v);
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
  el.replaceChildren(canvas, textLayer(runs, zoom));
  setStatus(`page ${index + 1} rendered in ${(performance.now() - t0).toFixed(0)} ms (${runs.length} text runs)`);
}

/** Transparent, selectable text positioned over the bitmap (pdf.js style). */
function textLayer(runs: TextRun[], scale: number): HTMLDivElement {
  const layer = document.createElement("div");
  layer.className = "textLayer";
  const meas = document.createElement("canvas").getContext("2d")!;
  for (const r of runs) {
    if (!r.font) continue;
    const span = document.createElement("span");
    span.textContent = r.text;
    const font = fontString(r.font, r.size);
    span.style.font = font;
    meas.font = font;
    const w = meas.measureText(r.text).width;
    const sx = r.advance > 0 && w > 0 ? r.advance / w : 1;
    const anchor = r.align === 1 ? -r.advance : r.align === 2 ? -r.advance / 2 : 0;
    const m = r.matrix;
    // matrix maps object space to page space; scale to CSS px.
    span.style.transform = `matrix(${m[0] * scale}, ${m[1] * scale}, ${m[2] * scale}, ${m[3] * scale}, ${r.x * scale}, ${r.y * scale}) translate(${anchor}px, -${r.size * 0.8}px) scaleX(${sx})`;
    layer.appendChild(span);
  }
  return layer;
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
      void client.continuous(v.id, { x: 0, y, w: width, h }, zoom * dpr()).then((bmp) => {
        const canvas = document.createElement("canvas");
        canvas.width = bmp.width;
        canvas.height = bmp.height;
        canvas.style.width = `${width * zoom}px`;
        canvas.style.height = `${h * zoom}px`;
        canvas.getContext("2d")!.drawImage(bmp, 0, 0);
        bmp.close();
        el.replaceChildren(canvas);
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
    setStatus(`sheet region ${Math.round(viewport.x)},${Math.round(viewport.y)} rendered in ${(performance.now() - t0).toFixed(0)} ms`);
  };
  stage.onscroll = () => { if (current === v) void draw(); };
  void draw();
}

main().catch((e) => setStatus(`error: ${e.message ?? e}`));
