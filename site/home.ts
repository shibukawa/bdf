// The top page: the document it shows, in the small viewer of
// examples/miniviewer, with a search box and zoom buttons.
import { MiniViewer } from "../examples/miniviewer/miniviewer.js";

const $ = <T extends HTMLElement>(id: string) => document.getElementById(id) as T;
const host = $("demo");
const viewer = new MiniViewer(host, { worker: new URL("lib/worker.js", location.href).href });
const q = $<HTMLInputElement>("q"), hits = $("hits");

function count() {
  hits.textContent = !q.value ? "" : viewer.hits.length ? `${viewer.hit + 1} / ${viewer.hits.length}` : "0";
}
async function search(step: number) {
  // the same words again go on to the next match
  await viewer.search(q.value, step);
  count();
}
q.onkeydown = (e) => { if (e.key === "Enter") search(e.shiftKey ? -1 : 1); };
q.oninput = () => { if (!q.value) search(1); };
$("next").onclick = () => search(1);
$("prev").onclick = () => search(-1);
$("zoomIn").onclick = () => viewer.setZoom(viewer.zoom * 1.25);
$("zoomOut").onclick = () => viewer.setZoom(viewer.zoom / 1.25);

viewer.open(host.dataset.src!).catch((e) => console.error(e));
