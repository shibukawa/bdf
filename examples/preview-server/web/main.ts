// The front end of preview-server: opens the bdf ?src= names — already
// converted by the server (main.go's makePreviews) — with the mini viewer
// (examples/miniviewer). There is no converter here at all: unlike
// light-server, this page never sees the original file, only a bdf one,
// fetched by range as its pages come into view.
import { MiniViewer } from "../../miniviewer/miniviewer.js";

const params = new URLSearchParams(location.search);
const src = params.get("src");
const status = document.getElementById("status")!;
const stage = document.getElementById("stage")!;

function setStatus(text: string) {
  status.textContent = text;
}

async function main() {
  if (!src) {
    setStatus("no ?src= given — pick a document from the file list");
    return;
  }
  document.title = `${src.split("/").pop()} – preview-server`;
  const viewer = new MiniViewer(stage, { worker: new URL("lib/worker.js", location.href).href });
  setStatus("opening…");
  const manifest = await viewer.open(src);
  setStatus(manifest ? `opened (${manifest.views.length} view(s))` : "encrypted: not opened");
}

main().catch((e) => setStatus(`error: ${(e as Error).message ?? e}`));
