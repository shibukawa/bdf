// Main-thread side of the converter worker: files opened in the viewer are
// converted into bdf inside the browser by the Go converters built as wasm.

/** A converted document (see cmd/bdfwasm). */
export interface Converted {
  bdf: Uint8Array;
  /** The input format ("pdf", "pptx", …). */
  format: string;
  summary: string;
  warnings: string[];
  /** The input needed its password (the document is not encrypted). */
  protected: boolean;
}

/**
 * A conversion done a page at a time (see cmd/bdfwasm's open): with pages
 * > 0, bdf is the outline (every page sized, the pages of the first view
 * without layers) and stream names the conversion for page() and finish();
 * otherwise bdf is the whole document.
 */
export interface Opened {
  bdf: Uint8Array;
  format: string;
  /** Pages left to convert with page(); 0 when bdf is the whole document. */
  pages: number;
  warnings: string[];
  /** The whole document's summary (pages 0). */
  summary?: string;
  stream?: number;
}

/** A page a stream converted: a bdf whose only page is that page, with the parts it brings. */
export interface ConvertedPage {
  bdf: Uint8Array;
  /** Warnings since the last call. */
  warnings: string[];
}

export interface ConvertOptions {
  format?: string;
  password?: string;
  /** URL of the font directory the Office converters lay text out with. */
  fonts?: string;
}

export type ConvertRequest =
  | { id: number; type: "convert" | "open"; module: string; data: ArrayBuffer; options: ConvertOptions }
  | { id: number; type: "page"; stream: number; index: number }
  | { id: number; type: "finish" | "close"; stream: number };
export type ConvertResponse = { id: number; ok: true; result: Converted | Opened | ConvertedPage | null } | { id: number; ok: false; error: string; code?: string };

/** A failed conversion; code is "password-required", "wrong-password" or "unknown-format" when it is one of those. */
export class ConvertError extends Error {
  constructor(message: string, readonly code?: string) {
    super(message);
    this.name = "ConvertError";
  }
}

/** What a file is, from its first bytes: a bdf document, a PDF, or (possibly) an Office document. */
export function sniff(head: Uint8Array): "bdf" | "pdf" | "office" {
  if (head[0] === 0x62 && head[1] === 0x64 && head[2] === 0x66 && head[3] === 0) return "bdf";
  // the header may follow some garbage in the first 1024 bytes
  const text = String.fromCharCode(...head.subarray(0, 1024));
  return text.includes("%PDF-") ? "pdf" : "office";
}

type Call = ConvertRequest extends infer R ? (R extends { id: number } ? Omit<R, "id"> : never) : never;

export class ConverterClient {
  private next = 1;
  private pending = new Map<number, { resolve: (v: unknown) => void; reject: (e: Error) => void }>();
  /** Why the worker stopped (it could not load wasm_exec.js, say); later calls fail with it. */
  private failed?: ConvertError;

  constructor(readonly worker: Worker) {
    worker.onmessage = (ev: MessageEvent<ConvertResponse>) => {
      const res = ev.data;
      const p = this.pending.get(res.id);
      if (!p) return;
      this.pending.delete(res.id);
      if (res.ok) p.resolve(res.result); else p.reject(new ConvertError(res.error, res.code));
    };
    worker.onerror = (ev) => {
      this.failed = new ConvertError(`the converter worker failed: ${ev.message || "is the site built with its converters?"}`);
      for (const p of this.pending.values()) p.reject(this.failed);
      this.pending.clear();
    };
  }

  private call<T>(req: Call): Promise<T> {
    if (this.failed) return Promise.reject(this.failed);
    const id = this.next++;
    return new Promise<T>((resolve, reject) => {
      this.pending.set(id, { resolve: resolve as (v: unknown) => void, reject });
      this.worker.postMessage({ id, ...req });
    });
  }

  /** Convert data (copied, so that it can be converted again) with the module at the URL. */
  convert(module: string, data: ArrayBuffer, options: ConvertOptions = {}): Promise<Converted> {
    return this.call<Converted>({ type: "convert", module, data, options });
  }
  /** Start a conversion done a page at a time (data is copied); formats that cannot come back whole. */
  open(module: string, data: ArrayBuffer, options: ConvertOptions = {}): Promise<Opened> {
    return this.call<Opened>({ type: "open", module, data, options });
  }
  /** Convert page index of the first view of an opened stream. */
  page(stream: number, index: number): Promise<ConvertedPage> {
    return this.call<ConvertedPage>({ type: "page", stream, index });
  }
  /** Convert the pages left and return the whole document. */
  finish(stream: number): Promise<Converted> {
    return this.call<Converted>({ type: "finish", stream });
  }
  /** Let a stream go (finished or not). */
  close(stream: number): Promise<null> {
    return this.call<null>({ type: "close", stream });
  }
}
