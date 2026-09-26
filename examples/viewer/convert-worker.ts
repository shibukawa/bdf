/// <reference lib="webworker" />
// Converter worker: runs the converter modules (cmd/bdfwasm, built as wasm)
// and converts the files the page sends into bdf documents. A classic worker,
// so that it can load Go's wasm_exec.js.
import type { ConvertRequest, ConvertResponse, Converted } from "./convert.js";

/** What cmd/bdfwasm sets on globalThis. */
interface Converter {
  convert(data: Uint8Array, options: { format?: string; password?: string; fonts?: string }): Promise<Converted>;
}
declare class Go {
  importObject: WebAssembly.Imports;
  run(instance: WebAssembly.Instance): Promise<void>;
}

importScripts("./wasm_exec.js");

/** Loaded modules by URL. Several can run in the worker, each a Go program of its own. */
const modules = new Map<string, Promise<Converter>>();

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

self.onmessage = async (ev: MessageEvent<ConvertRequest>) => {
  const { id, module, data, options } = ev.data;
  let res: ConvertResponse;
  const transfer: Transferable[] = [];
  try {
    const converter = await load(module);
    const result = await converter.convert(new Uint8Array(data), options);
    res = { id, ok: true, result };
    transfer.push(result.bdf.buffer);
  } catch (e) {
    const err = e as Error & { code?: string };
    res = { id, ok: false, error: err.message ?? String(e), code: err.code };
  }
  self.postMessage(res, transfer);
};
