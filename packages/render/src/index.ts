export { ResourceCache, ImageHold, DEFAULT_IMAGE_BUDGET, DEFAULT_MAX_IMAGE_PIXELS, DEFAULT_HOLD_LIMIT, imageSize, buildPath2D, fontString, embeddedFamily, type ResourceOptions } from "./resources.js";
export { CanvasRenderer, cssColor, resetState, filterAllowed, MAX_GROUP_DEPTH, type Ctx2D, type RenderOptions } from "./canvas.js";
export { PageRenderer, tilesIn, type PageRenderOptions, type TileRange } from "./page.js";
export { concat, within } from "./content.js";
export { BdfWorkerClient, BdfWorkerError } from "./client.js";
import { BdfWorkerClient } from "./client.js";
/** Create the rendering Worker from this package's bundled asset. */
export function createRenderWorker(): BdfWorkerClient {
  return new BdfWorkerClient(new Worker(new URL("./worker.js", import.meta.url), { type: "module" }));
}
export type { WorkerRequest, WorkerResponse, WorkerResult, WorkerCall, WorkerErrorCode, WorkerOpenOptions, OpenSource, RasterizeRequest, RasterizeResponse, PlayData, AudioData } from "./protocol.js";
export { MusicPlayer, cursorAtTick, tickAt, type Cursor, type PlayState, type MusicPlayerOptions } from "./player.js";
export { AudioPlayer, lineAtTime, timeOfLine, type AudioPlayerOptions } from "./audio.js";
export { VectorImage, domSvgRasterizer, isSvg, svgSize, type SvgRasterizer } from "./svg.js";
export { DocumentSearch, runRect, hasExtent, type HitRect, type Measure } from "./search.js";
export { buildTextLayer, linkHref, internalLink, selectedRuns, selectionText, selectionCells, joinRuns, installCopyHandler, TEXT_LAYER_CSS, RUN_ATTR, type TextLayerOptions, type SelectedRun, type InternalLink } from "./textlayer.js";
export { tableCells, cellClipboard, MAX_CLIPBOARD_CELLS, type CellText, type CellRange, type CellClipboard, type CellClipboardOptions } from "./cells.js";
