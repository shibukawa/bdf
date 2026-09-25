/// <reference lib="webworker" />
import { BdfDocument, BufferSource, RangeSource, SplitSource, fetchSingle, extractText, type TextRun, type Manifest, type Page, type Rect } from "@bdf/core";
import { fontString } from "./resources.js";
import { PageRenderer } from "./page.js";
import { DocumentSearch } from "./search.js";
import type { WorkerRequest, WorkerResponse, WorkerResult, OpenSource } from "./protocol.js";

let doc: BdfDocument | undefined;
let pages: PageRenderer | undefined;
let search: DocumentSearch | undefined;

const measureCtx = new OffscreenCanvas(1, 1).getContext("2d")!;
const measure = (font: string, text: string) => {
  measureCtx.font = font;
  return measureCtx.measureText(text).width;
};

async function open(source: OpenSource): Promise<Manifest> {
  switch (source.kind) {
    case "buffer": doc = await BdfDocument.open(new BufferSource(new Uint8Array(source.buffer))); break;
    case "single": doc = await BdfDocument.open(source.range ? new RangeSource(source.url) : await fetchSingle(source.url)); break;
    case "split": doc = await BdfDocument.open(new SplitSource(source.base)); break;
  }
  pages = new PageRenderer(doc!, {}, (self as unknown as { fonts?: FontFaceSet }).fonts);
  search = new DocumentSearch(doc!, measure);
  return doc!.manifest;
}

/**
 * Text runs of a page in the given space. Fonts are loaded first so runs
 * without an advance get one measured with the embedded font; the text layer
 * on the main thread stretches its fallback rendering to that width.
 */
async function pageText(page: Page, matrix?: [number, number, number, number, number, number]): Promise<TextRun[]> {
  await pages!.preparePage(page);
  const runs: TextRun[] = [];
  for (const layer of page.layers) {
    const obj = doc!.objectSync(layer.obj)!;
    runs.push(...extractText(obj, (h) => doc!.objectSync(h), matrix));
  }
  for (const r of runs) {
    if (r.advance === 0 && r.font && r.text) r.advance = measure(fontString(r.font, r.size), r.text);
  }
  return runs;
}

/** Runs of the pages whose body intersects viewport, moved into viewport coordinates. */
async function continuousText(view: ReturnType<BdfDocument["view"]>, viewport: Rect): Promise<TextRun[]> {
  const { offsets } = pages!.continuousLayout(view);
  const out: TextRun[] = [];
  const list = view.pages ?? [];
  for (let i = 0; i < list.length; i++) {
    const p = list[i];
    const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
    if (offsets[i] >= viewport.y + viewport.h || offsets[i] + b.h <= viewport.y) continue;
    const dx = -b.x - viewport.x, dy = offsets[i] - b.y - viewport.y;
    for (const r of await pageText(p, [1, 0, 0, 1, dx, dy])) {
      // Keep the runs whose anchor lies inside the body rectangle (the band clips to it).
      const px = r.x - dx, py = r.y - dy;
      if (px < b.x || px >= b.x + b.w || py < b.y || py > b.y + b.h) continue;
      out.push(r);
    }
  }
  return out;
}

function canvasFor(w: number, h: number): OffscreenCanvas {
  return new OffscreenCanvas(Math.max(1, Math.ceil(w)), Math.max(1, Math.ceil(h)));
}

async function handle(req: WorkerRequest): Promise<{ result: WorkerResult; transfer: Transferable[] }> {
  switch (req.type) {
    case "open":
      return { result: await open(req.source), transfer: [] };
    case "close":
      doc = pages = search = undefined;
      return { result: null, transfer: [] };
  }
  if (!doc || !pages || !search) throw new Error("bdf: no document open");
  const view = doc.view(req.view);
  switch (req.type) {
    case "page": {
      const page = view.pages?.[req.page];
      if (!page) throw new Error(`bdf: no page ${req.page}`);
      const canvas = canvasFor(page.w * req.scale, page.h * req.scale);
      await pages.renderPage(canvas.getContext("2d")!, page, { scale: req.scale, roles: req.roles });
      const bmp = canvas.transferToImageBitmap();
      return { result: bmp, transfer: [bmp] };
    }
    case "continuous": {
      const canvas = canvasFor(req.viewport.w * req.scale, req.viewport.h * req.scale);
      await pages.renderContinuous(canvas.getContext("2d")!, view, req.viewport, { scale: req.scale });
      const bmp = canvas.transferToImageBitmap();
      return { result: bmp, transfer: [bmp] };
    }
    case "sheet": {
      const canvas = canvasFor(req.viewport.w * req.scale, req.viewport.h * req.scale);
      await pages.renderSheet(canvas.getContext("2d")!, view, req.viewport, { scale: req.scale });
      const bmp = canvas.transferToImageBitmap();
      return { result: bmp, transfer: [bmp] };
    }
    case "text": {
      const page = view.pages?.[req.page];
      if (!page) throw new Error(`bdf: no page ${req.page}`);
      return { result: await pageText(page), transfer: [] };
    }
    case "continuousText":
      return { result: await continuousText(view, req.viewport), transfer: [] };
    case "search":
      return { result: await search.search(view, req.query, req.options), transfer: [] };
    case "locate": {
      // Fonts must be loaded for measureText; locate() loads the objects, which loads their fonts.
      const rects = [];
      for (const hit of req.hits) rects.push(await search.locate(view, hit));
      return { result: rects, transfer: [] };
    }
  }
}

self.onmessage = async (ev: MessageEvent<WorkerRequest>) => {
  const req = ev.data;
  try {
    const { result, transfer } = await handle(req);
    const res: WorkerResponse = { id: req.id, ok: true, result };
    (self as unknown as Worker).postMessage(res, transfer);
  } catch (e) {
    const res: WorkerResponse = { id: req.id, ok: false, error: e instanceof Error ? e.message : String(e) };
    (self as unknown as Worker).postMessage(res);
  }
};
