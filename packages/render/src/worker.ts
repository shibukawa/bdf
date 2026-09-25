/// <reference lib="webworker" />
import { BdfDocument, BufferSource, RangeSource, SplitSource, fetchSingle, extractText, type TextRun, type Manifest } from "@bdf/core";
import { PageRenderer } from "./page.js";
import type { WorkerRequest, WorkerResponse, WorkerResult, OpenSource } from "./protocol.js";

let doc: BdfDocument | undefined;
let pages: PageRenderer | undefined;

async function open(source: OpenSource): Promise<Manifest> {
  switch (source.kind) {
    case "buffer": doc = await BdfDocument.open(new BufferSource(new Uint8Array(source.buffer))); break;
    case "single": doc = await BdfDocument.open(source.range ? new RangeSource(source.url) : await fetchSingle(source.url)); break;
    case "split": doc = await BdfDocument.open(new SplitSource(source.base)); break;
  }
  pages = new PageRenderer(doc!, {}, (self as unknown as { fonts?: FontFaceSet }).fonts);
  return doc!.manifest;
}

function canvasFor(w: number, h: number): OffscreenCanvas {
  return new OffscreenCanvas(Math.max(1, Math.ceil(w)), Math.max(1, Math.ceil(h)));
}

async function handle(req: WorkerRequest): Promise<{ result: WorkerResult; transfer: Transferable[] }> {
  switch (req.type) {
    case "open":
      return { result: await open(req.source), transfer: [] };
    case "close":
      doc = pages = undefined;
      return { result: null, transfer: [] };
  }
  if (!doc || !pages) throw new Error("bdf: no document open");
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
      const runs: TextRun[] = [];
      for (const layer of page.layers) {
        const obj = await doc.ensure(layer.obj);
        runs.push(...extractText(obj, (h) => doc!.objectSync(h)));
      }
      return { result: runs, transfer: [] };
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
