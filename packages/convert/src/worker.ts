/// <reference lib="webworker" />
/** One module per Worker avoids overlapping TinyGo globals. */
interface Stream {
  page(index: number): Promise<{ bdf: Uint8Array; warnings: string[] }>;
  finish(): Promise<{ bdf: Uint8Array }>;
  close(): void;
}
interface Module {
  formats: { name: string; description: string; extensions: string[] }[];
  convert(data: Uint8Array, options: object): Promise<{ bdf: Uint8Array }>;
  open(data: Uint8Array, options: object): Promise<{ bdf: Uint8Array; stream?: Stream }>;
}
declare class Go {
  importObject: WebAssembly.Imports;
  run(instance: WebAssembly.Instance): Promise<void>;
}

let module: Module;
const streams = new Map<number, Stream>();
let nextStream = 1;

async function start(wasm: string): Promise<Module> {
  const go = new Go();
  const response = await fetch(wasm);
  if (!response.ok) throw new Error(`wasm fetch failed: ${response.status}`);
  // ArrayBuffer works on hosts that do not serve application/wasm.
  const { instance } = await WebAssembly.instantiate(await response.arrayBuffer(), go.importObject);
  let failure: unknown;
  go.run(instance).catch((error) => { failure = error; });
  const global = self as unknown as { bdfConverter?: Module };
  if (!global.bdfConverter) {
    await new Promise((resolve) => setTimeout(resolve, 0));
    throw new Error(`wasm did not start${failure ? `: ${failure}` : ""}`);
  }
  const result = global.bdfConverter;
  delete global.bdfConverter;
  return result;
}

self.onmessage = async (ev: MessageEvent<Record<string, unknown>>) => {
  const request = ev.data;
  const id = request.type === "init" ? 0 : Number(request.id);
  try {
    if (request.type === "init") {
      module = await start(String(request.wasm));
      self.postMessage({ id, ok: true, formats: module.formats });
      return;
    }
    if (!module) throw new Error("converter is not ready");
    let result: unknown;
    switch (request.type) {
      case "convert":
        result = await module.convert(new Uint8Array(request.data as ArrayBuffer), request.options as object);
        break;
      case "open": {
        const { stream, ...opened } = await module.open(new Uint8Array(request.data as ArrayBuffer), request.options as object);
        if (stream) {
          const handle = nextStream++;
          streams.set(handle, stream);
          result = { ...opened, stream: handle };
        } else result = opened;
        break;
      }
      case "page": {
        const stream = streams.get(Number(request.stream));
        if (!stream) throw new Error(`no conversion ${request.stream}`);
        result = await stream.page(Number(request.index));
        break;
      }
      case "finish": {
        const stream = streams.get(Number(request.stream));
        if (!stream) throw new Error(`no conversion ${request.stream}`);
        result = await stream.finish();
        break;
      }
      case "close":
        streams.get(Number(request.stream))?.close();
        streams.delete(Number(request.stream));
        result = null;
        break;
      default: throw new Error(`unknown request ${request.type}`);
    }
    const transfer: Transferable[] = [];
    if (result && typeof result === "object" && "bdf" in result) transfer.push((result as { bdf: Uint8Array }).bdf.buffer as ArrayBuffer);
    self.postMessage({ id, ok: true, result }, transfer);
  } catch (error) {
    const e = error as Error & { code?: string };
    self.postMessage({ id, ok: false, error: e.message ?? String(error), code: e.code });
  }
};
