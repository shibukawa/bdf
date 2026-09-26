/// <reference lib="webworker" />
// Converter worker: runs the converter modules (cmd/bdfwasm, built as wasm)
// and converts the files the page sends into bdf documents, whole or a page
// at a time. A classic worker, so that it can load Go's wasm_exec.js.
import type { ConvertRequest, ConvertResponse, ConvertOptions, Converted, ConvertedPage, Opened } from "./convert.js";

/** A conversion done a page at a time, as cmd/bdfwasm returns it. */
interface GoStream {
  page(index: number): Promise<ConvertedPage>;
  finish(): Promise<Converted>;
  close(): void;
}
/** What cmd/bdfwasm sets on globalThis. */
interface Converter {
  convert(data: Uint8Array, options: ConvertOptions): Promise<Converted>;
  open(data: Uint8Array, options: ConvertOptions): Promise<Omit<Opened, "stream"> & { stream?: GoStream }>;
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
      go.run(instance).catch(() => {}).finally(() => modules.delete(url)); // the program ended: load it again next time
      const g = self as unknown as { bdfConverter?: Converter };
      const converter = g.bdfConverter;
      delete g.bdfConverter;
      if (!converter) throw new Error(`${url}: not a bdf converter module`);
      return converter;
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

async function handle(req: ConvertRequest): Promise<Converted | Opened | ConvertedPage | null> {
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
  }
}

self.onmessage = async (ev: MessageEvent<ConvertRequest>) => {
  const { id } = ev.data;
  let res: ConvertResponse;
  const transfer: Transferable[] = [];
  try {
    const result = await handle(ev.data);
    res = { id, ok: true, result };
    if (result) transfer.push(result.bdf.buffer);
  } catch (e) {
    const err = e as Error & { code?: string };
    res = { id, ok: false, error: err.message ?? String(e), code: err.code };
  }
  self.postMessage(res, transfer);
};
