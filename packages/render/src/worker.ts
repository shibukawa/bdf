/// <reference lib="webworker" />
import { BdfDocument, BdfPasswordError, BdfSegmentError, SegmentLoader, BufferSource, RangeSource, SplitSource, UseLimits, fetchSingle, extractContent, tileSize, type PartSource, type TextContent, type Manifest, type Page, type Rect, type View } from "@bdfkit/core";
import { fontString } from "./resources.js";
import { PageRenderer, tilesIn, type TileRange } from "./page.js";
import { concat, within } from "./content.js";
import { DocumentSearch } from "./search.js";
import type { WorkerRequest, WorkerResponse, WorkerResult, WorkerErrorCode, OpenSource, WorkerOpenOptions, RasterizeRequest, RasterizeResponse } from "./protocol.js";
import type { SvgRasterizer } from "./svg.js";

/** An open document, with what draws and searches it. */
interface Opened {
  doc: BdfDocument;
  pages: PageRenderer;
  search: DocumentSearch;
  /** The server that hands out the pages of a document opened by segments. */
  segments?: SegmentLoader;
}

let open: Opened | undefined;
/** A source opened but still locked: kept for "unlock", so nothing is fetched or sent again. */
let locked: PartSource | undefined;
/** Documents closed or replaced, whose fonts and images go once no request is running. */
const retired: Opened[] = [];
let running = 0;
/** The options of the last open, which a replacing document keeps. */
let settings: WorkerOpenOptions = {};

const measureCtx = new OffscreenCanvas(1, 1).getContext("2d")!;
const measure = (font: string, text: string) => {
  measureCtx.font = font;
  return measureCtx.measureText(text).width;
};

/** SVG images being drawn by the page, by request id. */
const rasterizing = new Map<number, { resolve: (b: ImageBitmap) => void; reject: (e: Error) => void }>();
let nextRid = 1;

/**
 * Workers cannot decode SVG: the page draws SVG images for them
 * (BdfWorkerClient). A page that does not answer leaves them undrawn.
 */
const rasterizeOnPage: SvgRasterizer = (hash, data, width, height) =>
  new Promise((resolve, reject) => {
    const rid = nextRid++;
    const timer = setTimeout(() => {
      rasterizing.delete(rid);
      reject(new Error("the page did not draw the SVG image"));
    }, 30000);
    const done = () => { clearTimeout(timer); rasterizing.delete(rid); };
    rasterizing.set(rid, { resolve: (b) => { done(); resolve(b); }, reject: (e) => { done(); reject(e); } });
    const req: RasterizeRequest = { type: "rasterize", rid, hash, data, width, height };
    (self as unknown as Worker).postMessage(req);
  });

function sourceOf(source: OpenSource): PartSource | Promise<PartSource> {
  switch (source.kind) {
    case "buffer": return new BufferSource(new Uint8Array(source.buffer));
    case "single": return source.range ? new RangeSource(source.url) : fetchSingle(source.url);
    case "split": return new SplitSource(source.base);
    case "segments": throw new Error("bdf: a document of segments is opened, not replaced");
  }
}

function opened(doc: BdfDocument): Opened {
  const pages = new PageRenderer(doc, {}, (self as unknown as { fonts?: FontFaceSet }).fonts, { imageBudget: settings.imageBudget, maxImagePixels: settings.maxImagePixels, holdLimit: settings.holdLimit, rasterizeSvg: rasterizeOnPage });
  // hits are measured in the fonts of their objects, which are loaded with them
  return { doc, pages, search: new DocumentSearch(doc, measure, (h) => pages.res.prepareText(h)) };
}

function retire() {
  if (open) {
    open.segments?.close();
    retired.push(open);
  }
  open = locked = undefined;
}

async function openSource(source: OpenSource, password?: string, options: WorkerOpenOptions = {}): Promise<Manifest> {
  retire();
  settings = options;
  if (source.kind === "segments") {
    const segments = await SegmentLoader.open(source.url, { view: source.view, page: source.page });
    const o = { ...opened(segments.doc), segments };
    // the search of a view covers the pages that came
    segments.onSegment = (s) => o.search.forget(o.doc.view(s.view));
    open = o;
    return o.doc.manifest;
  }
  locked = await sourceOf(source);
  return unlock(password);
}

/** Fetch the pages of a document of segments that a request draws or reads, unless they came. */
async function ensurePages(o: Opened, view: View, pages: number[]): Promise<void> {
  if (o.segments) await Promise.all(pages.map((i) => o.segments!.ensure(view.id, i)));
}

/** The pages of a continuous layout whose body intersects viewport. */
function pagesIn(o: Opened, view: View, viewport: Rect): number[] {
  const { offsets } = o.pages.continuousLayout(view);
  const out: number[] = [];
  (view.pages ?? []).forEach((p, i) => {
    const h = (p.body ?? p).h;
    if (offsets[i] < viewport.y + viewport.h && offsets[i] + h > viewport.y) out.push(i);
  });
  return out;
}

/** Open the pending source; an encrypted one stays pending until a password opens it. */
async function unlock(password?: string): Promise<Manifest> {
  if (!locked) throw new Error("bdf: no document to unlock");
  const doc = await BdfDocument.open(locked, { password });
  locked = undefined;
  open = opened(doc);
  return doc.manifest;
}

/**
 * Swap in another document with the same views and pages (the finished
 * conversion of a streamed one). Requests running on the old one finish
 * with it.
 */
async function replace(source: OpenSource): Promise<Manifest> {
  const doc = await BdfDocument.open(await sourceOf(source));
  if (open) retired.push(open);
  open = opened(doc);
  return doc.manifest;
}

type Matrix = [number, number, number, number, number, number];

/**
 * Text content of a page in the given space. Fonts are loaded first so runs
 * without an advance get one measured with the embedded font; the text layer
 * on the main thread stretches its fallback rendering to that width. ALT_TEXT
 * runs keep theirs: only text drawn with a font has a known extent. limits
 * counts the objects drawn with USE across the pages of a request.
 */
async function pageContent({ doc, pages }: Opened, page: Page, matrix?: Matrix, roles?: string[], limits = new UseLimits()): Promise<TextContent> {
  await pages.preparePageText(page);
  const layers = roles ? page.layers.filter((l) => roles.includes(l.role)) : page.layers;
  const c = concat(layers.map((layer) => extractContent(doc.objectSync(layer.obj)!, (h) => doc.objectSync(h), matrix, limits)));
  for (const r of c.runs) {
    if (r.advance === 0 && r.font && r.text && !r.altText) r.advance = measure(fontString(r.font, r.size), r.text);
  }
  return c;
}

/** Content of the pages whose body intersects viewport, moved into viewport coordinates. */
async function continuousContent(o: Opened, view: View, viewport: Rect): Promise<TextContent> {
  const { offsets } = o.pages.continuousLayout(view);
  const parts: TextContent[] = [];
  const list = view.pages ?? [];
  const limits = new UseLimits();
  for (let i = 0; i < list.length; i++) {
    const p = list[i];
    const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
    if (offsets[i] >= viewport.y + viewport.h || offsets[i] + b.h <= viewport.y) continue;
    const dx = -b.x - viewport.x, dy = offsets[i] - b.y - viewport.y;
    // The layers continuous mode draws, and what lies inside the body rectangle (the band clips to it).
    parts.push(within(await pageContent(o, p, [1, 0, 0, 1, dx, dy], ["body", "annotation"], limits), b, dx, dy));
  }
  return concat(parts);
}

/**
 * Content of the sheet tiles that intersect viewport (or any of several
 * rectangles), in sheet coordinates, tiles in reading order.
 * A run belongs to the tile its anchor lies in; tiles repeat what straddles them.
 */
async function sheetContent({ doc, pages }: Opened, view: View, viewport: Rect | Rect[]): Promise<TextContent> {
  const tile = tileSize(view);
  const ranges: TileRange[] = [];
  for (const r of Array.isArray(viewport) ? viewport : [viewport]) {
    if (r.w <= 0 || r.h <= 0) continue;
    const tx0 = Math.max(0, Math.floor(r.x / tile)), ty0 = Math.max(0, Math.floor(r.y / tile));
    const tx1 = Math.floor((r.x + r.w - 1e-6) / tile), ty1 = Math.floor((r.y + r.h - 1e-6) / tile);
    ranges.push({ tx0, ty0, tx1, ty1 });
  }
  const parts: TextContent[] = [];
  const limits = new UseLimits();
  for (const [tx, ty, h] of tilesIn(view, ranges)) {
    await pages.res.prepareText(h);
    const c = extractContent(doc.objectSync(h)!, (hh) => doc.objectSync(hh), [1, 0, 0, 1, tx * tile, ty * tile], limits);
    parts.push(within(c, { x: 0, y: 0, w: tile, h: tile - 1e-6 }, tx * tile, ty * tile)); // the rule of the text index (spec §4.1)
  }
  return concat(parts);
}

/**
 * The pixels of the canvas of a render at most, and those of a side: what
 * browsers give a canvas. The size of a page is the document's to say.
 */
const MAX_CANVAS_PIXELS = 1 << 28, MAX_CANVAS_SIDE = 32767;

function canvasFor(w: number, h: number): OffscreenCanvas {
  const cw = Math.max(1, Math.ceil(w)), ch = Math.max(1, Math.ceil(h));
  if (!(cw <= MAX_CANVAS_SIDE && ch <= MAX_CANVAS_SIDE && cw * ch <= MAX_CANVAS_PIXELS)) throw new Error(`bdf: a canvas of ${cw} by ${ch} pixels is too large`);
  return new OffscreenCanvas(cw, ch);
}

async function handle(req: WorkerRequest): Promise<{ result: WorkerResult; transfer: Transferable[] }> {
  switch (req.type) {
    case "open":
      return { result: await openSource(req.source, req.password, req.options), transfer: [] };
    case "unlock":
      return { result: await unlock(req.password), transfer: [] };
    case "replace":
      return { result: await replace(req.source), transfer: [] };
    case "close":
      retire();
      return { result: null, transfer: [] };
  }
  // the document the request started on, even if another replaces it meanwhile
  const o = open;
  if (!o) throw new Error("bdf: no document open");
  const { doc, pages, search } = o;
  const view = doc.view(req.view);
  if (req.type === "addPage") {
    doc.addPage(view.id, req.page, await BdfDocument.open(new BufferSource(new Uint8Array(req.buffer))));
    search.forget(view);
    return { result: null, transfer: [] };
  }
  switch (req.type) {
    case "page": case "text": case "content":
      await ensurePages(o, view, [req.page]);
      break;
    case "continuous": case "continuousText": case "continuousContent":
      await ensurePages(o, view, pagesIn(o, view, req.viewport));
      break;
    case "locate":
      await ensurePages(o, view, req.hits.flatMap((h) => h.segments.map((s) => s.a)));
      break;
  }
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
    case "text": case "content": {
      const page = view.pages?.[req.page];
      if (!page) throw new Error(`bdf: no page ${req.page}`);
      const c = await pageContent(o, page);
      return { result: req.type === "text" ? c.runs : c, transfer: [] };
    }
    case "continuousText": case "continuousContent": {
      const c = await continuousContent(o, view, req.viewport);
      return { result: req.type === "continuousText" ? c.runs : c, transfer: [] };
    }
    case "sheetContent":
      return { result: await sheetContent(o, view, req.viewport), transfer: [] };
    case "search":
      return { result: await search.search(view, req.query, req.options), transfer: [] };
    case "locate": {
      // Fonts must be loaded for measureText; locate() loads the objects, which loads their fonts.
      const rects = [];
      for (const hit of req.hits) rects.push(await search.locate(view, hit));
      return { result: rects, transfer: [] };
    }
    case "play": {
      const play = await doc.play(view);
      if (!play) return { result: null, transfer: [] };
      // a copy: the document keeps its part (which may be a view of the whole file)
      const seq = play.seq.slice().buffer;
      return { result: { seq, cues: play.cues }, transfer: [seq] };
    }
    case "audio": {
      const audio = await doc.audio(view);
      if (!audio) return { result: null, transfer: [] };
      // a blob: the browser keeps the file (on disk, when it is large), and the page gets it without another copy
      return { result: { blob: new Blob([audio.data as BlobPart], { type: audio.type }), cues: audio.cues }, transfer: [] };
    }
  }
}

function errorCode(e: unknown): WorkerErrorCode | undefined {
  if (e instanceof BdfPasswordError) return e.reason === "required" ? "password-required" : "wrong-password";
  if (e instanceof BdfSegmentError) return e.status === 429 ? "rate-limited" : e.status === 401 || e.status === 403 ? "not-allowed" : undefined;
  return undefined;
}

self.onmessage = async (ev: MessageEvent<WorkerRequest | RasterizeResponse>) => {
  if ("rid" in ev.data) {
    const res = ev.data;
    const p = rasterizing.get(res.rid);
    if (res.ok) {
      if (p) p.resolve(res.bitmap); else res.bitmap.close(); // too late
    } else {
      p?.reject(new Error(res.error));
    }
    return;
  }
  const req = ev.data;
  running++;
  try {
    const { result, transfer } = await handle(req);
    const res: WorkerResponse = { id: req.id, ok: true, result };
    (self as unknown as Worker).postMessage(res, transfer);
  } catch (e) {
    const code = errorCode(e);
    const res: WorkerResponse = { id: req.id, ok: false, error: e instanceof Error ? e.message : String(e), code };
    (self as unknown as Worker).postMessage(res);
  } finally {
    // nothing uses the retired documents any more
    if (--running === 0) for (const r of retired.splice(0)) r.pages.dispose();
  }
};
