// Test harness: renders fixture pages on the main thread and via the worker,
// and compares canvases against golden PNGs. Driven by test/golden.mjs.
import { BdfDocument, fetchSingle, RangeSource, SplitSource, type Rect } from "@bdf/core";
import { PageRenderer, BdfWorkerClient } from "@bdf/render";

export interface Case {
  name: string;
  kind: "page" | "continuous" | "sheet";
  view: string;
  page?: number;
  viewport?: Rect;
  scale: number;
  roles?: string[];
}

export const CASES: Case[] = [
  { name: "slides-1", kind: "page", view: "slides", page: 0, scale: 1 },
  { name: "slides-2", kind: "page", view: "slides", page: 1, scale: 1 },
  { name: "slides-3", kind: "page", view: "slides", page: 2, scale: 1 },
  { name: "slides-1-body-only", kind: "page", view: "slides", page: 0, scale: 0.5, roles: ["body"] },
  { name: "doc-1", kind: "page", view: "doc", page: 0, scale: 1 },
  { name: "doc-continuous", kind: "continuous", view: "doc", viewport: { x: 0, y: 600, w: 451.3, h: 400 }, scale: 1 },
  { name: "sheet-tile-boundary", kind: "sheet", view: "sheet1", viewport: { x: 0, y: 1800, w: 800, h: 500 }, scale: 1 },
  { name: "sheet-zoomed", kind: "sheet", view: "sheet1", viewport: { x: 64, y: 20, w: 300, h: 150 }, scale: 2 },
];

function canvasSize(c: Case, doc: BdfDocument): [number, number] {
  if (c.kind === "page") {
    const p = doc.view(c.view).pages![c.page!];
    return [Math.ceil(p.w * c.scale), Math.ceil(p.h * c.scale)];
  }
  return [Math.ceil(c.viewport!.w * c.scale), Math.ceil(c.viewport!.h * c.scale)];
}

async function renderMain(doc: BdfDocument, pr: PageRenderer, c: Case): Promise<HTMLCanvasElement> {
  const [w, h] = canvasSize(c, doc);
  const canvas = document.createElement("canvas");
  canvas.width = w;
  canvas.height = h;
  const ctx = canvas.getContext("2d")!;
  const view = doc.view(c.view);
  switch (c.kind) {
    case "page": await pr.renderPage(ctx, view.pages![c.page!], { scale: c.scale, roles: c.roles }); break;
    case "continuous": await pr.renderContinuous(ctx, view, c.viewport!, { scale: c.scale }); break;
    case "sheet": await pr.renderSheet(ctx, view, c.viewport!, { scale: c.scale }); break;
  }
  return canvas;
}

async function renderWorker(client: BdfWorkerClient, c: Case): Promise<ImageBitmap> {
  switch (c.kind) {
    case "page": return client.page(c.view, c.page!, c.scale, c.roles);
    case "continuous": return client.continuous(c.view, c.viewport!, c.scale);
    case "sheet": return client.sheet(c.view, c.viewport!, c.scale);
  }
}

function toCanvas(bmp: ImageBitmap): HTMLCanvasElement {
  const canvas = document.createElement("canvas");
  canvas.width = bmp.width;
  canvas.height = bmp.height;
  canvas.getContext("2d")!.drawImage(bmp, 0, 0);
  return canvas;
}

export interface Diff { total: number; different: number; maxDelta: number; width: number; height: number; png: string }

function compare(a: HTMLCanvasElement, b: HTMLCanvasElement, channelTol: number): Diff {
  const w = Math.max(a.width, b.width), h = Math.max(a.height, b.height);
  const da = a.getContext("2d")!.getImageData(0, 0, a.width, a.height).data;
  const db = b.getContext("2d")!.getImageData(0, 0, b.width, b.height).data;
  const diff = document.createElement("canvas");
  diff.width = w;
  diff.height = h;
  const dctx = diff.getContext("2d")!;
  const dimg = dctx.createImageData(w, h);
  let different = 0, maxDelta = 0;
  for (let y = 0; y < h; y++) for (let x = 0; x < w; x++) {
    let delta = 0;
    for (let ch = 0; ch < 4; ch++) {
      const va = x < a.width && y < a.height ? da[(y * a.width + x) * 4 + ch] : -1;
      const vb = x < b.width && y < b.height ? db[(y * b.width + x) * 4 + ch] : -1;
      delta = Math.max(delta, Math.abs(va - vb));
    }
    const o = (y * w + x) * 4;
    dimg.data[o] = 255; dimg.data[o + 1] = delta > channelTol ? 0 : 255; dimg.data[o + 2] = delta > channelTol ? 0 : 255; dimg.data[o + 3] = 255;
    if (delta > channelTol) different++;
    maxDelta = Math.max(maxDelta, delta);
  }
  dctx.putImageData(dimg, 0, 0);
  return { total: w * h, different, maxDelta, width: w, height: h, png: different ? diff.toDataURL("image/png") : "" };
}

async function loadImage(url: string): Promise<HTMLCanvasElement | null> {
  const res = await fetch(url);
  if (!res.ok) return null;
  const bmp = await createImageBitmap(await res.blob());
  return toCanvas(bmp);
}

export interface Result {
  name: string;
  png: string; // data URL of the main-thread render
  workerPng: string; // data URL of the worker render
  golden: Diff | null; // null when no golden exists
  worker: Diff; // main thread vs worker
  textRuns: number;
}

async function main() {
  const params = new URLSearchParams(location.search);
  const sourceKind = params.get("source") ?? "single";
  const workerUrl = params.get("worker") ?? "./worker.js";
  const doc = await BdfDocument.open(
    sourceKind === "split" ? new SplitSource("/fixtures/demo-split/")
      : sourceKind === "range" ? new RangeSource("/fixtures/demo.bdf")
        : await fetchSingle("/fixtures/demo.bdf"),
  );
  const pr = new PageRenderer(doc, {}, document.fonts);
  const client = new BdfWorkerClient(new Worker(workerUrl, { type: "module" }));
  await client.open(sourceKind === "split" ? { kind: "split", base: location.origin + "/fixtures/demo-split/" } : { kind: "single", url: location.origin + "/fixtures/demo.bdf", range: sourceKind === "range" });

  const results: Result[] = [];
  for (const c of CASES) {
    const main = await renderMain(doc, pr, c);
    main.title = c.name;
    document.body.appendChild(main);
    const worker = toCanvas(await renderWorker(client, c));
    const golden = await loadImage(`/fixtures/golden/${c.name}.png`);
    const textRuns = c.kind === "page" ? (await client.text(c.view, c.page!)).length : 0;
    results.push({
      name: c.name,
      png: main.toDataURL("image/png"),
      workerPng: worker.toDataURL("image/png"),
      golden: golden ? compare(main, golden, 24) : null,
      worker: compare(main, worker, 0),
      textRuns,
    });
  }
  (window as unknown as { bdfResults: Result[] }).bdfResults = results;
  document.title = "done";
}

main().catch((e) => {
  (window as unknown as { bdfError: string }).bdfError = String(e?.stack ?? e);
  document.title = "error";
});
