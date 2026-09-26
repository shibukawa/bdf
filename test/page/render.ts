// Renders every page of a BDF document (given by ?src=) into canvases; used by test/render.mjs.
import { BdfDocument, fetchSingle } from "@bdf/core";
import { PageRenderer } from "@bdf/render";

async function main() {
  const params = new URLSearchParams(location.search);
  const scale = Number(params.get("scale") ?? "1");
  const doc = await BdfDocument.open(await fetchSingle(params.get("src")!));
  const pr = new PageRenderer(doc, {}, document.fonts);
  const out: { name: string; png: string }[] = [];
  for (const view of doc.manifest.views) {
    if (view.kind === "sheet") {
      // the whole sheet, up to 4096 px each way
      const { width, height } = pr.sheetSize(view);
      const viewport = { x: 0, y: 0, w: Math.min(width, 4096 / scale), h: Math.min(height, 4096 / scale) };
      const canvas = document.createElement("canvas");
      canvas.width = Math.ceil(viewport.w * scale);
      canvas.height = Math.ceil(viewport.h * scale);
      document.body.appendChild(canvas);
      await pr.renderSheet(canvas.getContext("2d")!, view, viewport, { scale });
      out.push({ name: view.id, png: canvas.toDataURL("image/png") });
      continue;
    }
    const pages = view.pages ?? [];
    for (let i = 0; i < pages.length; i++) {
      const p = pages[i];
      const canvas = document.createElement("canvas");
      canvas.width = Math.ceil(p.w * scale);
      canvas.height = Math.ceil(p.h * scale);
      document.body.appendChild(canvas);
      await pr.renderPage(canvas.getContext("2d")!, p, { scale });
      out.push({ name: `${view.id}-${i + 1}`, png: canvas.toDataURL("image/png") });
    }
  }
  (window as unknown as { bdfPages: unknown }).bdfPages = out;
  document.title = "done";
}
main().catch((e) => {
  (window as unknown as { bdfError: string }).bdfError = String(e?.stack ?? e);
  document.title = "error";
});
