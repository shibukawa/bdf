import type { Manifest, Rect, TextRun } from "@bdf/core";
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
  text(view: string, page: number): Promise<TextRun[]> {
    return this.call<TextRun[]>({ type: "text", view, page });
  }
  close(): Promise<null> {
    return this.call<null>({ type: "close" });
  }
  terminate(): void {
    this.worker.terminate();
  }
}
