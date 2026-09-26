import type { Manifest, Rect, TextRun, TextContent, SearchHit, SearchOptions } from "@bdf/core";
import type { HitRect } from "./search.js";

/** How the worker should open a document. */
export type OpenSource =
  | { kind: "buffer"; buffer: ArrayBuffer }
  | { kind: "single"; url: string; range?: boolean }
  | { kind: "split"; base: string };

export type WorkerRequest =
  | { id: number; type: "open"; source: OpenSource }
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
  | { id: number; type: "close" };

export type WorkerResult = Manifest | ImageBitmap | TextRun[] | TextContent | SearchHit[] | HitRect[][] | null;

export type WorkerResponse =
  | { id: number; ok: true; result: WorkerResult }
  | { id: number; ok: false; error: string };

/** A request without its id; the id is assigned by the client. */
export type WorkerCall = WorkerRequest extends infer R ? (R extends { id: number } ? Omit<R, "id"> : never) : never;
