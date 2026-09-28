// The front end of light-server: fetches the document ?file= names from
// /files/, converts it to bdf inside this page — the server never sees the
// converted document — and shows it with the mini viewer
// (examples/miniviewer). This is exactly what the demo viewer's own
// openFile does (examples/viewer/main.ts), pared down to one file.
//
// Built by web/build.mjs into web/ next to this source: main.js, and lib/
// (the rendering worker, the converter worker and the wasm modules built by
// examples/common/build.mjs — the same helper the demo site uses) and
// fonts/ (the free fonts the Office converters lay text out with).
import { ConverterClient, sniff } from "../../common/convert.js";
import { MiniViewer } from "../../miniviewer/miniviewer.js";

const params = new URLSearchParams(location.search);
const name = params.get("file");
const status = document.getElementById("status")!;
const stage = document.getElementById("stage")!;
/** A URL next to this page, absolute: workers resolve relative URLs against their own script, not the page's. */
const here = (path: string) => new URL(path, location.href).href;

const MODULES = { pdf: "lib/bdf-pdf.wasm", office: "lib/bdf-office.wasm", web: "lib/bdf-web.wasm", image: "lib/bdf-image.wasm" };

function setStatus(text: string) {
  status.textContent = text;
}

async function main() {
  if (!name) {
    setStatus("no ?file= given — pick one from the file list");
    return;
  }
  document.title = `${name} – light-server`;
  setStatus(`fetching ${name}…`);
  const res = await fetch(`/files/${encodeURIComponent(name)}`);
  if (!res.ok) throw new Error(`${name}: HTTP ${res.status}`);
  const data = await res.arrayBuffer();
  const kind = sniff(new Uint8Array(data), name);

  const viewer = new MiniViewer(stage, { worker: here("lib/worker.js") });
  if (kind === "bdf") {
    setStatus("opening…");
    await viewer.open({ kind: "buffer", buffer: data });
    setStatus(`${name}: opened as bdf directly`);
    return;
  }

  setStatus(`converting ${name} in your browser…`);
  const conv = new ConverterClient(new Worker(here("lib/convert-worker.js")));
  const t0 = performance.now();
  const out = await conv.convert(here(MODULES[kind]), data, { fonts: kind === "office" ? here("fonts/") : undefined, name });
  setStatus(`opening ${name} (converted in ${(performance.now() - t0).toFixed(0)} ms: ${out.summary})…`);
  await viewer.open({ kind: "buffer", buffer: out.bdf.buffer as ArrayBuffer });
  setStatus(`${name}: ${out.summary}`);
}

main().catch((e) => setStatus(`error: ${(e as Error).message ?? e}`));
