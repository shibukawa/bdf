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
  /** The input needed its password (pages 0). */
  protected?: boolean;
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
  /** The pages (slides, sheets) to convert, as bdf generate -pages takes them: "1" for a thumbnail. All when absent. */
  pages?: string;
  /** The format's own options, as bdf generate -param takes them: { play: "false" } leaves the audio of an audio file out. */
  params?: Record<string, string>;
}

/** A thumbnail the preview module drew (see cmd/bdfwasm), as the thumbnail package does on a server. */
export interface Thumbnail {
  /** The encoded image. */
  image: Uint8Array;
  format: "png" | "jpeg";
  width: number;
  height: number;
  /** The layout used. */
  mode: "crop" | "fit";
  /** What the drawing left out. */
  warnings: string[];
}

export interface ThumbnailOptions {
  /** The side of a cropped thumbnail, the longer side of a fitted one, in pixels (default 256). */
  size?: number;
  /** The layout: auto (default: from the kind of document), crop or fit. */
  mode?: "auto" | "crop" | "fit";
  format?: "png" | "jpeg";
  /** The resolution a sheet is drawn at (default 72): a lower value shows more cells, smaller. */
  sheetDpi?: number;
  /** The id of the view to draw (default: the first one). */
  view?: string;
  /** The password of an encrypted document. */
  password?: string;
  /** URL of the font directory, for text in fonts referred to by name. */
  fonts?: string;
}

/** The text of a document for a search index, as bdf text writes it (JSON). */
export interface SearchText {
  json: string;
}

/** A PDF file the preview or the form module wrote of a bdf document (see cmd/bdfwasm), as bdf pdf does. */
export interface PdfFile {
  pdf: Uint8Array;
  pages: number;
  /** The fonts embedded. */
  fonts: number;
  /** What the file leaves out. */
  warnings: string[];
}

export interface PdfOptions {
  /** The password of an encrypted document (the file is not encrypted). */
  password?: string;
  /** URL of the font directory, for text in fonts referred to by name. */
  fonts?: string;
  /** The ids of the views to write (default: those with pages; the others when there are none). */
  views?: string[];
  /** Roles of layers to leave out ("background", "master"): for paper that has them printed. */
  skip?: string[];
  /** A tagged PDF: with the structure of the document (headings, paragraphs, tables), for screen readers. */
  tagged?: boolean;
  /** A PDF/A file, for keeping: "2b", "2u", "2a" (tagged), or "3b", "3u", "3a". */
  pdfa?: string;
}

/** A document written as Markdown (see cmd/bdfwasm's markdown). */
export interface MarkdownFile {
  /** The text, UTF-8. */
  markdown: Uint8Array;
  /** The images the text refers to, as files beside it (options.images "files"); none otherwise. */
  images: { name: string; data: Uint8Array }[];
  /** A zip file of the text (name.md) and the images, when asked for. */
  zip?: Uint8Array;
  /** The views written. */
  views: number;
  /** What the text leaves out. */
  warnings: string[];
}

export interface MarkdownOptions {
  /** The password of an encrypted document (the text is plain). */
  password?: string;
  /** The ids of the views to write (default: the scroll view when there is one, else all). */
  views?: string[];
  /** Roles of layers to leave out (default: header, footer, master and background). */
  skip?: string[];
  /** Keep the running heads, feet and page numbers of pages. */
  keepRunning?: boolean;
  /** Keep links into the document itself (#page=N), which have no target in Markdown. */
  internalLinks?: boolean;
  /** Images: as files returned beside the text (default), embedded in it as data URLs, or left out. */
  images?: "files" | "embed" | "none";
  /** Also a zip file of the text and the image files. */
  zip?: boolean;
  /** The name of the text in the zip file (default "document"). */
  name?: string;
}

/** What a template of a form asks of data (see cmd/bdfwasm's fields), as JSON texts. */
export interface TemplateFields {
  /** The names the tags look up, as bdf template schema lists them. */
  fields: string;
  /** The JSON Schema of the data (form.Schema). */
  schema: string;
  /** The format of the template ("xlsx", "docx", "pptx", "pdf", "pdfme"; of a compiled template, that of the template it holds). */
  format?: string;
  /** The sample of the data that the template has for itself: a compiled template (.bdft), a template of pdfme. */
  sample?: string;
}

export interface FillOptions extends ConvertOptions {
  /** The data, as the JSON of an object. */
  data: string;
  /** The pictures that the data names, by their paths. */
  files?: Record<string, Uint8Array>;
  /** Fail when the data lacks a name of a tag (default: the tag is left empty, with a warning). */
  strict?: boolean;
}

/** A filled template: a bdf document. */
export interface Filled {
  bdf: Uint8Array;
  /** The format of the template ("xlsx", "docx", "pptx", "pdf" or "pdfme"). */
  format: string;
  pages: number;
  warnings: string[];
}

export type ConvertRequest =
  | { id: number; type: "convert" | "open"; module: string; data: ArrayBuffer; options: ConvertOptions }
  | { id: number; type: "page"; stream: number; index: number }
  | { id: number; type: "finish" | "close"; stream: number }
  | { id: number; type: "thumbnail"; module: string; data: ArrayBuffer; options: ThumbnailOptions }
  | { id: number; type: "text"; module: string; data: ArrayBuffer; options: { password?: string } }
  | { id: number; type: "pdf"; module: string; data: ArrayBuffer; options: PdfOptions }
  | { id: number; type: "markdown"; module: string; data: ArrayBuffer; options: MarkdownOptions }
  | { id: number; type: "fields"; module: string; data: ArrayBuffer; options: { sample?: string } }
  | { id: number; type: "fill"; module: string; data: ArrayBuffer; options: FillOptions };
export type ConvertResult = Converted | Opened | ConvertedPage | Thumbnail | SearchText | PdfFile | MarkdownFile | TemplateFields | Filled | null;
export type ConvertResponse = { id: number; ok: true; result: ConvertResult } | { id: number; ok: false; error: string; code?: string };

/** A failed conversion; code is "password-required", "wrong-password" or "unknown-format" when it is one of those. */
export class ConvertError extends Error {
  constructor(message: string, readonly code?: string) {
    super(message);
    this.name = "ConvertError";
  }
}

/** Extensions of the files the web module converts (HTML, Markdown and EPUB). */
const WEB = [".html", ".htm", ".xhtml", ".mhtml", ".mht", ".md", ".markdown", ".mdown", ".mkd", ".mdx", ".epub"];

/** Extensions of audio files, for those whose content does not say (an MP3 without a tag, an MP4 file). */
const AUDIO = [".mp3", ".m4a", ".m4b", ".m4p", ".aac", ".flac", ".ogg", ".oga", ".opus", ".spx", ".wav", ".wave", ".aif", ".aiff", ".aifc"];

/**
 * What a file is, from its content and its name: a bdf document, an EPUB,
 * an HTML or Markdown document (by its extension: a README may start with
 * an SVG picture), an image the browser displays by itself or an audio
 * file (its cover and tags; both by the small image module), a PDF, or
 * (possibly) an Office document or a diagram. draw.io's PNG and SVG
 * exports are diagrams.
 */
export function sniff(data: Uint8Array, name = ""): "bdf" | "image" | "pdf" | "web" | "office" {
  if (data[0] === 0x62 && data[1] === 0x64 && data[2] === 0x66 && data[3] === 0) return "bdf";
  // an EPUB starts with its mimetype file, stored
  const zip = String.fromCharCode(...data.subarray(0, 58));
  if (zip.startsWith("PK\x03\x04") && zip.slice(30) === "mimetypeapplication/epub+zip") return "web";
  const dot = name.lastIndexOf(".");
  const ext = dot >= 0 ? name.slice(dot).toLowerCase() : "";
  if (WEB.includes(ext)) return "web";
  if (isImage(data)) return "image";
  if (isAudio(data) || AUDIO.includes(ext)) return "image";
  // the header may follow some garbage in the first 1024 bytes
  const text = String.fromCharCode(...data.subarray(0, 1024));
  return text.includes("%PDF-") ? "pdf" : "office";
}

/** Whether data starts as an audio file does: an ID3 tag, FLAC, Ogg, WAVE, AIFF or an M4A brand. */
function isAudio(b: Uint8Array): boolean {
  const s = (at: number, n: number) => String.fromCharCode(...b.subarray(at, at + n));
  if (s(0, 3) === "ID3" && b[3] >= 2 && b[3] <= 4) return true;
  if (s(0, 4) === "fLaC" || s(0, 4) === "OggS") return true;
  if ((s(0, 4) === "RIFF" || s(0, 4) === "RF64") && s(8, 4) === "WAVE") return true;
  if (s(0, 4) === "FORM" && (s(8, 4) === "AIFF" || s(8, 4) === "AIFC")) return true;
  return s(4, 4) === "ftyp" && /^M4[ABP] /.test(s(8, 4));
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

  private call<T>(req: Call, transfer: Transferable[] = []): Promise<T> {
    if (this.failed) return Promise.reject(this.failed);
    const id = this.next++;
    return new Promise<T>((resolve, reject) => {
      this.pending.set(id, { resolve: resolve as (v: unknown) => void, reject });
      this.worker.postMessage({ id, ...req }, transfer);
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
  /** Draw the thumbnail of a single-file bdf (transferred) with the preview module at the URL. */
  thumbnail(module: string, data: ArrayBuffer, options: ThumbnailOptions = {}): Promise<Thumbnail> {
    return this.call<Thumbnail>({ type: "thumbnail", module, data, options }, [data]);
  }
  /** The text of a single-file bdf (transferred) for a search index, with the preview module at the URL. */
  text(module: string, data: ArrayBuffer, options: { password?: string } = {}): Promise<SearchText> {
    return this.call<SearchText>({ type: "text", module, data, options }, [data]);
  }
  /** Write a single-file bdf (transferred) as a PDF file, with the preview or the form module at the URL. */
  pdf(module: string, data: ArrayBuffer, options: PdfOptions = {}): Promise<PdfFile> {
    return this.call<PdfFile>({ type: "pdf", module, data, options }, [data]);
  }
  /** Write a single-file bdf (transferred) as Markdown, with the preview module at the URL. */
  markdown(module: string, data: ArrayBuffer, options: MarkdownOptions = {}): Promise<MarkdownFile> {
    return this.call<MarkdownFile>({ type: "markdown", module, data, options }, [data]);
  }
  /** What a template of a form (copied) asks of data, with the form module at the URL. */
  fields(module: string, data: ArrayBuffer, options: { sample?: string } = {}): Promise<TemplateFields> {
    return this.call<TemplateFields>({ type: "fields", module, data, options });
  }
  /** Fill a template of a form (copied) with data, with the form module at the URL. */
  fill(module: string, data: ArrayBuffer, options: FillOptions): Promise<Filled> {
    return this.call<Filled>({ type: "fill", module, data, options });
  }
}
