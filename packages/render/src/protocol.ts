import type { Manifest, Rect, TextRun, TextContent, SearchHit, SearchOptions, Cues } from "@bdfkit/core";
import type { HitRect } from "./search.js";

/** How the worker should open a document. */
export type OpenSource =
  | { kind: "buffer"; buffer: ArrayBuffer }
  | { kind: "single"; url: string; range?: boolean }
  | { kind: "split"; base: string }
  /**
   * A document a server hands out a few pages at a time, each segment sealed
   * for a key pair made for its request (spec §3.6, SegmentLoader): pages
   * are fetched as they are drawn. view and page name the page whose
   * segment comes first.
   */
  | { kind: "segments"; url: string; view?: string; page?: number };

/** The music of a view (spec §4.4): its Standard MIDI File, and its cues when it has them. */
export interface PlayData {
  seq: ArrayBuffer;
  cues: Cues | null;
}

/** Settings of the worker for a document. */
export interface WorkerOpenOptions {
  /** Bytes of decoded images the worker keeps (see ResourceOptions.imageBudget). */
  imageBudget?: number;
  /** Images of more pixels are not decoded (see ResourceOptions.maxImagePixels). */
  maxImagePixels?: number;
  /** Bytes of decoded images one render may hold (see ResourceOptions.holdLimit). */
  holdLimit?: number;
}

export type WorkerRequest =
  /** Open a document; an encrypted one needs its password, here or with "unlock". */
  | { id: number; type: "open"; source: OpenSource; password?: string; options?: WorkerOpenOptions }
  /** Retry the document "open" left locked with another password. */
  | { id: number; type: "unlock"; password: string }
  /** Put a page a streamed conversion returned (a page document) in place of a page of the open document. */
  | { id: number; type: "addPage"; view: string; page: number; buffer: ArrayBuffer }
  /** Swap in another document with the same views and pages: the finished conversion of a streamed one. */
  | { id: number; type: "replace"; source: OpenSource }
  | { id: number; type: "page"; view: string; page: number; scale: number; roles?: string[] }
  | { id: number; type: "continuous"; view: string; viewport: Rect; scale: number }
  | { id: number; type: "sheet"; view: string; viewport: Rect; scale: number }
  | { id: number; type: "text"; view: string; page: number }
  | { id: number; type: "continuousText"; view: string; viewport: Rect }
  | { id: number; type: "content"; view: string; page: number }
  | { id: number; type: "continuousContent"; view: string; viewport: Rect }
  | { id: number; type: "sheetContent"; view: string; viewport: Rect | Rect[] }
  | { id: number; type: "search"; view: string; query: string; options?: SearchOptions }
  | { id: number; type: "locate"; view: string; hits: SearchHit[] }
  /** The music of a view; null when it has none. */
  | { id: number; type: "play"; view: string }
  | { id: number; type: "close" };

export type WorkerResult = Manifest | ImageBitmap | TextRun[] | TextContent | SearchHit[] | HitRect[][] | PlayData | null;

/**
 * Why a request failed, when the viewer has something to do about it: ask
 * for a password, or tell the reader that the server does not give these
 * pages (401, 403) or not so fast (429).
 */
export type WorkerErrorCode = "password-required" | "wrong-password" | "not-allowed" | "rate-limited";

export type WorkerResponse =
  | { id: number; ok: true; result: WorkerResult }
  | { id: number; ok: false; error: string; code?: WorkerErrorCode };

/**
 * The worker asks the page to draw an SVG image into a bitmap of width ×
 * height pixels: workers cannot decode SVG (spec §6.2). BdfWorkerClient
 * answers with a RasterizeResponse.
 */
export type RasterizeRequest = { type: "rasterize"; rid: number; hash: string; data: Uint8Array; width: number; height: number };
export type RasterizeResponse = { rid: number; ok: true; bitmap: ImageBitmap } | { rid: number; ok: false; error: string };

/** A request without its id; the id is assigned by the client. */
export type WorkerCall = WorkerRequest extends infer R ? (R extends { id: number } ? Omit<R, "id"> : never) : never;
