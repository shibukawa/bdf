/** Options passed to the browser converter. */
export interface ConvertOptions {
  format?: string;
  password?: string;
  fonts?: string;
  name?: string;
  pages?: string;
}

export interface Converted {
  bdf: Uint8Array;
  format: string;
  summary: string;
  warnings: string[];
  protected: boolean;
}

export interface ConvertedPage { bdf: Uint8Array; warnings: string[] }
export interface Opened {
  bdf: Uint8Array;
  format: string;
  pages: number;
  warnings: string[];
  summary?: string;
  protected?: boolean;
  stream?: number;
}
export interface Format { name: string; description: string; extensions: string[] }

export class ConvertError extends Error {
  constructor(message: string, readonly code?: string) {
    super(message);
    this.name = "ConvertError";
  }
}

type Reply =
  | { id: number; ok: true; result: unknown }
  | { id: number; ok: false; error: string; code?: string };
type Request =
  | { type: "convert" | "open"; data: ArrayBuffer; options: ConvertOptions }
  | { type: "page"; stream: number; index: number }
  | { type: "finish" | "close"; stream: number };

/** A dedicated Worker owns one TinyGo module. Close it when the viewer is removed. */
export class Converter {
  private next = 1;
  private pending = new Map<number, { resolve(value: unknown): void; reject(error: Error): void }>();
  private failed?: ConvertError;
  readonly ready: Promise<Format[]>;

  constructor(readonly worker: Worker, wasmUrl: string | URL) {
    this.ready = new Promise<Format[]>((resolve, reject) => {
      const startup = (ev: MessageEvent<Reply & { formats?: Format[] }>) => {
        if (ev.data.id !== 0) return;
        worker.removeEventListener("message", startup);
        if (ev.data.ok) resolve(ev.data.formats ?? []);
        else reject(new ConvertError(ev.data.error, ev.data.code));
      };
      worker.addEventListener("message", startup);
      worker.addEventListener("error", (ev) => reject(new ConvertError(`converter worker failed: ${ev.message}`)), { once: true });
      worker.postMessage({ type: "init", wasm: String(wasmUrl) });
    });
    worker.addEventListener("message", (ev: MessageEvent<Reply>) => {
      const reply = ev.data;
      if (typeof reply?.id !== "number" || reply.id === 0) return;
      const call = this.pending.get(reply.id);
      if (!call) return;
      this.pending.delete(reply.id);
      if (reply.ok) call.resolve(reply.result);
      else call.reject(new ConvertError(reply.error, reply.code));
    });
    worker.addEventListener("error", (ev) => {
      this.failed = new ConvertError(`converter worker failed: ${ev.message}`);
      for (const call of this.pending.values()) call.reject(this.failed);
      this.pending.clear();
    });
  }

  private async call<T>(request: Request): Promise<T> {
    await this.ready;
    if (this.failed) throw this.failed;
    const id = this.next++;
    return new Promise<T>((resolve, reject) => {
      this.pending.set(id, { resolve: resolve as (value: unknown) => void, reject });
      this.worker.postMessage({ id, ...request });
    });
  }

  convert(data: ArrayBuffer | Uint8Array, options: ConvertOptions = {}): Promise<Converted> {
    return this.call({ type: "convert", data: bytes(data), options });
  }
  open(data: ArrayBuffer | Uint8Array, options: ConvertOptions = {}): Promise<Opened> {
    return this.call({ type: "open", data: bytes(data), options });
  }
  page(stream: number, index: number): Promise<ConvertedPage> {
    return this.call({ type: "page", stream, index });
  }
  finish(stream: number): Promise<Converted> {
    return this.call({ type: "finish", stream });
  }
  close(stream: number): Promise<null> {
    return this.call({ type: "close", stream });
  }
  terminate(): void {
    this.failed = new ConvertError("converter terminated");
    for (const call of this.pending.values()) call.reject(this.failed);
    this.pending.clear();
    this.worker.terminate();
  }
}

function bytes(data: ArrayBuffer | Uint8Array): ArrayBuffer {
  if (data instanceof ArrayBuffer) return data;
  return data.buffer.slice(data.byteOffset, data.byteOffset + data.byteLength) as ArrayBuffer;
}

/** The worker and TinyGo runtime are assets of @bdfkit/convert. */
export function createConverter(wasmUrl: string | URL): Converter {
  return new Converter(new Worker(new URL("./worker.js", import.meta.url)), wasmUrl);
}
