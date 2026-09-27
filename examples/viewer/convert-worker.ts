/// <reference lib="webworker" />
// Converter worker: runs the converter modules (cmd/bdfwasm, built as wasm)
// and converts the files the page sends into bdf documents, whole or a page
// at a time; the preview module draws the thumbnails of bdf documents and
// gives their text for a search index. A classic worker, so that it can
// load Go's wasm_exec.js.
import type { ConvertRequest, ConvertResponse, ConvertOptions, ConvertResult, Converted, ConvertedPage, Opened, SearchText, Thumbnail, ThumbnailOptions } from "./convert.js";

/** A conversion done a page at a time, as cmd/bdfwasm returns it. */
interface GoStream {
  page(index: number): Promise<ConvertedPage>;
  finish(): Promise<Converted>;
  close(): void;
}
/** What cmd/bdfwasm sets on globalThis; thumbnail and text are the preview module's. */
interface Converter {
  convert(data: Uint8Array, options: ConvertOptions): Promise<Converted>;
  open(data: Uint8Array, options: ConvertOptions): Promise<Omit<Opened, "stream"> & { stream?: GoStream }>;
  thumbnail?(data: Uint8Array, options: ThumbnailOptions): Promise<Thumbnail>;
  text?(data: Uint8Array, options: { password?: string }): Promise<SearchText>;
}
declare class Go {
  importObject: WebAssembly.Imports;
  run(instance: WebAssembly.Instance): Promise<void>;
}

importScripts("./wasm_exec.js");

/** Loaded modules by URL. Several can run in the worker, each a Go program of its own. */
const modules = new Map<string, Promise<Converter>>();
/** Streams opened, by the number the page knows them by. */
const streams = new Map<number, GoStream>();
let nextStream = 1;

function load(url: string): Promise<Converter> {
  let p = modules.get(url);
  if (!p) {
    p = (async () => {
      const go = new Go();
      const { instance } = await WebAssembly.instantiateStreaming(fetch(url), go.importObject);
      // run() executes main until it waits for calls, so the global is set when it returns
      let failure: unknown;
      go.run(instance).catch((e) => { failure = e; }).finally(() => modules.delete(url)); // the program ended: load it again next time
      const g = self as unknown as { bdfConverter?: Converter };
      const converter = g.bdfConverter;
      delete g.bdfConverter;
      if (converter) return converter;
      // why main stopped (the stack overflowed, say): run() has rejected, the catch comes a task later
      await new Promise((r) => setTimeout(r, 0));
      throw new Error(failure ? `${url} did not start: ${failure}` : `${url}: not a bdf converter module`);
    })();
    p.catch(() => modules.delete(url));
    modules.set(url, p);
  }
  return p;
}

function stream(id: number): GoStream {
  const s = streams.get(id);
  if (!s) throw new Error(`no conversion ${id}`);
  return s;
}

/** The preview module at a URL. */
async function preview(url: string): Promise<Required<Converter>> {
  const m = await load(url);
  if (!m.thumbnail || !m.text) throw new Error(`${url}: not the preview module`);
  return m as Required<Converter>;
}

async function handle(req: ConvertRequest): Promise<ConvertResult> {
  switch (req.type) {
    case "convert":
      return (await load(req.module)).convert(new Uint8Array(req.data), req.options);
    case "open": {
      const { stream: s, ...opened } = await (await load(req.module)).open(new Uint8Array(req.data), req.options);
      if (!s) return opened;
      const id = nextStream++;
      streams.set(id, s);
      return { ...opened, stream: id };
    }
    case "page":
      return stream(req.stream).page(req.index);
    case "finish":
      return stream(req.stream).finish();
    case "close":
      streams.get(req.stream)?.close();
      streams.delete(req.stream);
      return null;
    case "thumbnail":
      return (await preview(req.module)).thumbnail(new Uint8Array(req.data), req.options);
    case "text":
      return (await preview(req.module)).text(new Uint8Array(req.data), req.options);
  }
}

self.onmessage = async (ev: MessageEvent<ConvertRequest>) => {
  const { id } = ev.data;
  let res: ConvertResponse;
  const transfer: Transferable[] = [];
  try {
    const result = await handle(ev.data);
    res = { id, ok: true, result };
    if (result && "bdf" in result) transfer.push(result.bdf.buffer);
    if (result && "image" in result) transfer.push(result.image.buffer);
  } catch (e) {
    const err = e as Error & { code?: string };
    res = { id, ok: false, error: err.message ?? String(e), code: err.code };
  }
  self.postMessage(res, transfer);
};
