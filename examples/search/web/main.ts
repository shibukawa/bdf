// The document page of search: like examples/preview-server's, it opens the
// bdf ?src= names with the mini viewer (examples/miniviewer) — nothing is
// converted here either. The one addition is reading a search hit's place
// in the document off the URL fragment (#view=ID&page=N, set by
// search.ts) and jumping straight to it once the document is open.
import { MiniViewer } from "../../miniviewer/miniviewer.js";

const params = new URLSearchParams(location.search);
const src = params.get("src");
const target = new URLSearchParams(location.hash.slice(1));
const status = document.getElementById("status")!;
const stage = document.getElementById("stage")!;

function setStatus(text: string) {
  status.textContent = text;
}

async function main() {
  if (!src) {
    setStatus("no ?src= given — search from the home page");
    return;
  }
  document.title = `${src.split("/").pop()} – search`;
  const viewer = new MiniViewer(stage, { worker: new URL("lib/worker.js", location.href).href });
  setStatus("opening…");
  const manifest = await viewer.open(src);
  if (!manifest) {
    setStatus("encrypted: not opened");
    return;
  }
  const view = target.get("view");
  const page = Number(target.get("page") ?? "1") - 1;
  if (view && manifest.views.some((v) => v.id === view)) viewer.show(view, Math.max(0, page));
  setStatus(`opened (${manifest.views.length} view(s))`);
}

main().catch((e) => setStatus(`error: ${(e as Error).message ?? e}`));
