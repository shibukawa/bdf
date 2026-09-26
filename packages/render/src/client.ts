import type { Manifest, Rect, TextRun, TextContent, SearchHit, SearchOptions } from "@bdf/core";
import type { HitRect } from "./search.js";
import type { WorkerCall, WorkerResponse, OpenSource } from "./protocol.js";

/** Main-thread handle to a rendering worker. */
export class BdfWorkerClient {
  private next = 1;
  private pending = new Map<number, { resolve: (v: unknown) => void; reject: (e: Error) => void }>();

  constructor(readonly worker: Worker) {
    worker.onmessage = (ev: MessageEvent<WorkerResponse>) => {
      const res = ev.data;
      const p = this.pending.get(res.id);
      if (!p) return;
      this.pending.delete(res.id);
      if (res.ok) p.resolve(res.result); else p.reject(new Error(res.error));
    };
  }

  private call<T>(req: WorkerCall, transfer: Transferable[] = []): Promise<T> {
    const id = this.next++;
    return new Promise<T>((resolve, reject) => {
      this.pending.set(id, { resolve: resolve as (v: unknown) => void, reject });
      this.worker.postMessage({ id, ...req }, transfer);
    });
  }

  open(source: OpenSource): Promise<Manifest> {
    return this.call<Manifest>({ type: "open", source }, source.kind === "buffer" ? [source.buffer] : []);
  }
  page(view: string, page: number, scale: number, roles?: string[]): Promise<ImageBitmap> {
    return this.call<ImageBitmap>({ type: "page", view, page, scale, roles });
  }
  continuous(view: string, viewport: Rect, scale: number): Promise<ImageBitmap> {
    return this.call<ImageBitmap>({ type: "continuous", view, viewport, scale });
  }
  sheet(view: string, viewport: Rect, scale: number): Promise<ImageBitmap> {
    return this.call<ImageBitmap>({ type: "sheet", view, viewport, scale });
  }
  /** Text runs of a page in page units; advances are measured with the embedded fonts. */
  text(view: string, page: number): Promise<TextRun[]> {
    return this.call<TextRun[]>({ type: "text", view, page });
  }
  /** Text runs of the continuous layout inside viewport, in viewport coordinates. */
  continuousText(view: string, viewport: Rect): Promise<TextRun[]> {
    return this.call<TextRun[]>({ type: "continuousText", view, viewport });
  }
  /** Runs of a page with its structure and links (for accessible text layers), in page units. */
  content(view: string, page: number): Promise<TextContent> {
    return this.call<TextContent>({ type: "content", view, page });
  }
  /** Content of the continuous layout inside viewport, in viewport coordinates. */
  continuousContent(view: string, viewport: Rect): Promise<TextContent> {
    return this.call<TextContent>({ type: "continuousContent", view, viewport });
  }
  /**
   * Content of the sheet tiles that intersect viewport (or any of several
   * rectangles: the frozen panes and the scrolled region), in sheet
   * coordinates.
   */
  sheetContent(view: string, viewport: Rect | Rect[]): Promise<TextContent> {
    return this.call<TextContent>({ type: "sheetContent", view, viewport });
  }
  search(view: string, query: string, options?: SearchOptions): Promise<SearchHit[]> {
    return this.call<SearchHit[]>({ type: "search", view, query, options });
  }
  /** Rectangles for each hit, in page or sheet coordinates. */
  locate(view: string, hits: SearchHit[]): Promise<HitRect[][]> {
    return this.call<HitRect[][]>({ type: "locate", view, hits });
  }
  close(): Promise<null> {
    return this.call<null>({ type: "close" });
  }
  terminate(): void {
    this.worker.terminate();
  }
}
