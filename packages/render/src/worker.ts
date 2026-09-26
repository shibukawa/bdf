/// <reference lib="webworker" />
import { BdfDocument, BdfPasswordError, BufferSource, RangeSource, SplitSource, fetchSingle, extractContent, type PartSource, type TextRun, type TextContent, type Manifest, type Page, type Rect, type View } from "@bdf/core";
import { fontString } from "./resources.js";
import { PageRenderer } from "./page.js";
import { DocumentSearch } from "./search.js";
import type { WorkerRequest, WorkerResponse, WorkerResult, OpenSource, WorkerOpenOptions, RasterizeRequest, RasterizeResponse } from "./protocol.js";
import type { SvgRasterizer } from "./svg.js";

let doc: BdfDocument | undefined;
/** A source opened but still locked: kept for "unlock", so nothing is fetched or sent again. */
let locked: PartSource | undefined;
let pages: PageRenderer | undefined;
let search: DocumentSearch | undefined;
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

/** Forget the document, closing its decoded images. */
function close(): void {
  pages?.res.dispose();
  doc = pages = search = locked = undefined;
}

async function open(source: OpenSource, password?: string, options: WorkerOpenOptions = {}): Promise<Manifest> {
  close();
  settings = options;
  switch (source.kind) {
    case "buffer": locked = new BufferSource(new Uint8Array(source.buffer)); break;
    case "single": locked = source.range ? new RangeSource(source.url) : await fetchSingle(source.url); break;
    case "split": locked = new SplitSource(source.base); break;
  }
  return unlock(password);
}

/** Open the pending source; an encrypted one stays pending until a password opens it. */
async function unlock(password?: string): Promise<Manifest> {
  if (!locked) throw new Error("bdf: no document to unlock");
  doc = await BdfDocument.open(locked, { password });
  locked = undefined;
  pages = new PageRenderer(doc!, {}, (self as unknown as { fonts?: FontFaceSet }).fonts, { imageBudget: settings.imageBudget, rasterizeSvg: rasterizeOnPage });
  search = new DocumentSearch(doc!, measure);
  return doc!.manifest;
}

type Matrix = [number, number, number, number, number, number];

/** Join the contents of several objects (layers, pages, tiles), renumbering nodes and runs. */
function concat(parts: TextContent[]): TextContent {
  const out: TextContent = { runs: [], nodes: [], links: [] };
  for (const c of parts) {
    const nodeBase = out.nodes.length, runBase = out.runs.length;
    const shift = (i: number | undefined) => (i === undefined || i < 0 ? i : i + nodeBase);
    out.nodes.push(...c.nodes.map((n) => ({ ...n, parent: shift(n.parent)! })));
    out.runs.push(...c.runs.map((r) => ({ ...r, node: shift(r.node) })));
    out.links.push(...c.links.map((l) => ({ ...l, after: l.after + runBase, node: shift(l.node)! })));
  }
  return out;
}

/**
 * Keep the runs and links inside a rectangle (by run anchor, by link
 * center), for a band or a tile; the text layer leaves out nodes left empty.
 */
function within(c: TextContent, r: Rect, dx = 0, dy = 0): TextContent {
  const inside = (x: number, y: number) => x - dx >= r.x && x - dx < r.x + r.w && y - dy >= r.y && y - dy <= r.y + r.h;
  const runs: TextRun[] = [];
  const kept: number[] = []; // new index of the last kept run up to each old index
  for (const run of c.runs) {
    if (inside(run.x, run.y)) runs.push(run);
    kept.push(runs.length - 1);
  }
  const links = c.links.filter((l) => inside(l.x + l.w / 2, l.y + l.h / 2)).map((l) => ({ ...l, after: l.after >= 0 ? kept[l.after] : -1 }));
  return { runs, nodes: c.nodes, links };
}

/**
 * Text content of a page in the given space. Fonts are loaded first so runs
 * without an advance get one measured with the embedded font; the text layer
 * on the main thread stretches its fallback rendering to that width. ALT_TEXT
 * runs keep theirs: only text drawn with a font has a known extent.
 */
async function pageContent(page: Page, matrix?: Matrix, roles?: string[]): Promise<TextContent> {
  await pages!.preparePageText(page);
  const layers = roles ? page.layers.filter((l) => roles.includes(l.role)) : page.layers;
  const c = concat(layers.map((layer) => extractContent(doc!.objectSync(layer.obj)!, (h) => doc!.objectSync(h), matrix)));
  for (const r of c.runs) {
    if (r.advance === 0 && r.font && r.text && !r.altText) r.advance = measure(fontString(r.font, r.size), r.text);
  }
  return c;
}

/** Content of the pages whose body intersects viewport, moved into viewport coordinates. */
async function continuousContent(view: View, viewport: Rect): Promise<TextContent> {
  const { offsets } = pages!.continuousLayout(view);
  const parts: TextContent[] = [];
  const list = view.pages ?? [];
  for (let i = 0; i < list.length; i++) {
    const p = list[i];
    const b = p.body ?? { x: 0, y: 0, w: p.w, h: p.h };
    if (offsets[i] >= viewport.y + viewport.h || offsets[i] + b.h <= viewport.y) continue;
    const dx = -b.x - viewport.x, dy = offsets[i] - b.y - viewport.y;
    // The layers continuous mode draws, and what lies inside the body rectangle (the band clips to it).
    parts.push(within(await pageContent(p, [1, 0, 0, 1, dx, dy], ["body", "annotation"]), b, dx, dy));
  }
  return concat(parts);
}

/**
 * Content of the sheet tiles that intersect viewport (or any of several
 * rectangles), in sheet coordinates, tiles in reading order.
 * A run belongs to the tile its anchor lies in; tiles repeat what straddles them.
 */
async function sheetContent(view: View, viewport: Rect | Rect[]): Promise<TextContent> {
  const tile = view.tile ?? 2048;
  const keys = new Map<string, [number, number]>();
  for (const r of Array.isArray(viewport) ? viewport : [viewport]) {
    if (r.w <= 0 || r.h <= 0) continue;
    const tx0 = Math.max(0, Math.floor(r.x / tile)), ty0 = Math.max(0, Math.floor(r.y / tile));
    const tx1 = Math.floor((r.x + r.w - 1e-6) / tile), ty1 = Math.floor((r.y + r.h - 1e-6) / tile);
    for (let ty = ty0; ty <= ty1; ty++) for (let tx = tx0; tx <= tx1; tx++) keys.set(`${tx},${ty}`, [tx, ty]);
  }
  const parts: TextContent[] = [];
  for (const [tx, ty] of [...keys.values()].sort((a, b) => a[1] - b[1] || a[0] - b[0])) {
    const h = view.tiles?.[`${tx},${ty}`];
    if (!h) continue;
    await pages!.res.prepareText(h);
    const c = extractContent(doc!.objectSync(h)!, (hh) => doc!.objectSync(hh), [1, 0, 0, 1, tx * tile, ty * tile]);
    parts.push(within(c, { x: 0, y: 0, w: tile, h: tile - 1e-6 }, tx * tile, ty * tile)); // the rule of the text index (spec §4.1)
  }
  return concat(parts);
}

function canvasFor(w: number, h: number): OffscreenCanvas {
  return new OffscreenCanvas(Math.max(1, Math.ceil(w)), Math.max(1, Math.ceil(h)));
}

async function handle(req: WorkerRequest): Promise<{ result: WorkerResult; transfer: Transferable[] }> {
  switch (req.type) {
    case "open":
      return { result: await open(req.source, req.password, req.options), transfer: [] };
    case "unlock":
      return { result: await unlock(req.password), transfer: [] };
    case "close":
      close();
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
    case "text": case "content": {
      const page = view.pages?.[req.page];
      if (!page) throw new Error(`bdf: no page ${req.page}`);
      const c = await pageContent(page);
      return { result: req.type === "text" ? c.runs : c, transfer: [] };
    }
    case "continuousText": case "continuousContent": {
      const c = await continuousContent(view, req.viewport);
      return { result: req.type === "continuousText" ? c.runs : c, transfer: [] };
    }
    case "sheetContent":
      return { result: await sheetContent(view, req.viewport), transfer: [] };
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
  try {
    const { result, transfer } = await handle(req);
    const res: WorkerResponse = { id: req.id, ok: true, result };
    (self as unknown as Worker).postMessage(res, transfer);
  } catch (e) {
    const code = e instanceof BdfPasswordError ? (e.reason === "required" ? "password-required" : "wrong-password") : undefined;
    const res: WorkerResponse = { id: req.id, ok: false, error: e instanceof Error ? e.message : String(e), code };
    (self as unknown as Worker).postMessage(res);
  }
};
