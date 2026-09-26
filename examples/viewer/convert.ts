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

export interface ConvertOptions {
  format?: string;
  password?: string;
  /** URL of the font directory the Office converters lay text out with. */
  fonts?: string;
}

export type ConvertRequest = { id: number; module: string; data: ArrayBuffer; options: ConvertOptions };
export type ConvertResponse = { id: number; ok: true; result: Converted } | { id: number; ok: false; error: string; code?: string };

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

export class ConverterClient {
  private next = 1;
  private pending = new Map<number, { resolve: (v: Converted) => void; reject: (e: Error) => void }>();
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

  /** Convert data (copied, so that it can be converted again) with the module at the URL. */
  convert(module: string, data: ArrayBuffer, options: ConvertOptions = {}): Promise<Converted> {
    if (this.failed) return Promise.reject(this.failed);
    const id = this.next++;
    return new Promise<Converted>((resolve, reject) => {
      this.pending.set(id, { resolve, reject });
      this.worker.postMessage({ id, module, data, options } satisfies ConvertRequest);
    });
  }
}
