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
  /** The file name: its extension tells the format of a file whose content does not, and images are named after it. */
  name?: string;
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

/** Extensions of the files the web module converts (HTML and Markdown). */
const WEB = [".html", ".htm", ".xhtml", ".mhtml", ".mht", ".md", ".markdown", ".mdown", ".mkd", ".mdx"];

/**
 * What a file is, from its content and its name: a bdf document, an image
 * the browser displays by itself (stored as it is by a small module), a
 * PDF, an HTML or Markdown document (by its extension), or (possibly) an
 * Office document or a diagram. draw.io's PNG and SVG exports are diagrams.
 */
export function sniff(data: Uint8Array, name = ""): "bdf" | "image" | "pdf" | "web" | "office" {
  if (data[0] === 0x62 && data[1] === 0x64 && data[2] === 0x66 && data[3] === 0) return "bdf";
  if (isImage(data)) return "image";
  // the header may follow some garbage in the first 1024 bytes
  const text = String.fromCharCode(...data.subarray(0, 1024));
  if (text.includes("%PDF-")) return "pdf";
  const dot = name.lastIndexOf(".");
  return dot >= 0 && WEB.includes(name.slice(dot).toLowerCase()) ? "web" : "office";
}

/** Whether data is a PNG, JPEG, GIF, WebP, AVIF, BMP, ICO or SVG image (and not a draw.io export). */
function isImage(b: Uint8Array): boolean {
  const s = (at: number, n: number) => String.fromCharCode(...b.subarray(at, at + n));
  const dv = new DataView(b.buffer, b.byteOffset, b.byteLength);
  if (s(0, 8) === "\x89PNG\r\n\x1a\n") {
    // draw.io keeps the diagram in a text chunk
    for (let at = 8; at + 12 <= b.length; at += 12 + dv.getUint32(at)) {
      const type = s(at + 4, 4);
      if ((type === "tEXt" || type === "zTXt") && /^(mxfile|mxGraphModel)\0/.test(s(at + 8, 13))) return false;
      if (type === "IEND") break;
    }
    return true;
  }
  if (b[0] === 0xff && b[1] === 0xd8 && b[2] === 0xff) return true;
  if (s(0, 6) === "GIF87a" || s(0, 6) === "GIF89a") return true;
  if (s(0, 4) === "RIFF" && s(8, 4) === "WEBP") return true;
  if (b.length >= 16 && s(4, 4) === "ftyp") {
    // the major brand, then the compatible ones
    for (let at = 8; at + 4 <= Math.min(b.length, dv.getUint32(0)); at += 4) if (at !== 12 && /^avi[fs]$/.test(s(at, 4))) return true;
    return false;
  }
  if (b.length >= 18 && s(0, 2) === "BM" && [12, 16, 40, 52, 56, 64, 108, 124].includes(dv.getUint32(14, true))) return true;
  if (b.length >= 22 && dv.getUint16(0, true) === 0 && dv.getUint16(2, true) === 1 && dv.getUint16(4, true) > 0 && b[9] === 0) {
    const size = dv.getUint32(14, true), at = dv.getUint32(18, true);
    return size > 0 && at + size <= b.length; // an icon
  }
  // SVG: the root element, after the XML declaration, comments and the document type
  let text = new TextDecoder().decode(b.subarray(0, 65536));
  for (let prev = ""; prev !== text; ) {
    prev = text;
    text = text.replace(/^\s+|^<\?[\s\S]*?\?>|^<!--[\s\S]*?-->|^<!DOCTYPE[^[>]*(\[[\s\S]*?\])?\s*>/i, "");
  }
  const root = /^<(?:[\w.-]+:)?svg[\s/>][^>]*/.exec(text);
  return !!root && !/\scontent="[^"]*&lt;(mxfile|mxGraphModel)/.test(root[0]);
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
