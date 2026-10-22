// Browser integration tests for rendering, text selection, accessibility and search.
// Usage: node test/browser.mjs [--source single|split|range]
import { build } from "esbuild";
import { chromium } from "playwright-core";
import { mkdir, copyFile } from "node:fs/promises";
import { existsSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";
import { serve } from "./serve.mjs";

const root = dirname(dirname(fileURLToPath(import.meta.url)));
const args = process.argv.slice(2);
const source = args.includes("--source") ? args[args.indexOf("--source") + 1] : "single";
if (!(["single", "split", "range"].includes(source))) throw new Error(`unknown source: ${source}`);

const outDir = join(root, "test", ".out");
await mkdir(outDir, { recursive: true });
await build({
  entryPoints: [
    { in: join(root, "test/page/harness.ts"), out: "harness" },
    { in: join(root, "packages/render/src/worker.ts"), out: "worker" },
  ],
  bundle: true, format: "esm", outdir: outDir, sourcemap: true, target: "es2022", logLevel: "warning",
});

await copyFile(join(root, "test/page/index.html"), join(outDir, "page.html"));

const { server, port } = await serve(root);
// CHROMIUM_PATH supports an existing browser; otherwise use Playwright's install.
const executablePath = process.env.CHROMIUM_PATH ?? ["/opt/pw-browsers/chromium_headless_shell-1194/chrome-linux/headless_shell"].find((p) => existsSync(p));
let browser;
let failed = 0;
try {
  browser = await chromium.launch(executablePath ? { executablePath } : {});
  const page = await browser.newPage({ deviceScaleFactor: 1 });
  page.on("console", (m) => { if (m.type() === "error") console.error("[browser]", m.text()); });
  page.on("pageerror", (e) => console.error("[pageerror]", e.message));
  await page.goto(`http://127.0.0.1:${port}/test/.out/page.html?source=${source}&worker=/test/.out/worker.js`);
  await page.waitForFunction(() => document.title === "done" || document.title === "error", null, { timeout: 60000 });
  const error = await page.evaluate(() => window.bdfError);
  if (error) throw new Error(error);
  const results = await page.evaluate(() => window.bdfResults);
  const search = await page.evaluate(() => window.bdfSearch);
  const selection = await page.evaluate(() => window.bdfSelection);
  const selectionOk = selection.spans > 4 && selection.allText === selection.wantAll && selection.partText === selection.wantPart
    && selection.breaks > 0 && selection.allText.includes(" ");
  if (!selectionOk) failed++;
  console.log(`selection: ${selectionOk ? "ok" : "FAIL"} (${selection.spans} spans, ${selection.breaks} line breaks, partial ${JSON.stringify(selection.partText)})`);
  const alt = await page.evaluate(() => window.bdfAltText);
  const altOk = alt.alt.length > 0 && alt.spans === alt.runs && alt.alt.every((a) => a.span === a.text && Math.abs(a.width - a.want) < 1);
  if (!altOk) failed++;
  console.log(`alt text: ${altOk ? "ok" : "FAIL"} (${alt.spans} spans for ${alt.runs} runs, ${JSON.stringify(alt.alt)})`);
  const st = await page.evaluate(() => window.bdfStructure);
  const eq = (a, b) => JSON.stringify(a) === JSON.stringify(b);
  const structureOk = st.slide.heading === 1 && st.headingLevel === "1" && st.slide.list === 1 && st.slide.listitem === 4
    && eq(st.links, [["#page=1", "1", "Text styles"], ["https://example.com/", null, "system serif italic (no correction)"]])
    && st.flow.table === 1 && st.flow.row === 4 && st.flow.cell === 9 && eq(st.headers, ["R1C1", "R1C2", "R1C3"])
    && eq(st.figures, ["Checkerboard", "Enlarged corner of the checkerboard"]) && eq(st.figureBox, ["1560px", "260px", "240px", "400px"])
    && st.sheet.table === 1 && st.sheet.columnheader === 26 && st.sheet.rowheader > 50 && st.sheetCell === "2"
    && eq(st.lang, ["en", "", "ja", ""])
    && eq(st.xlsxHeaders, ["Item", "Units", "Price"]) && eq(st.xlsxFigures, ["Column chart: Units", "Picture", "A rounded note box"])
    && st.xlsxSpans[0] > 50 && st.xlsxSpans[0] === st.xlsxSpans[1] && st.xlsxSpans[2] === 1;
  if (!structureOk) failed++;
  console.log(`structure: ${structureOk ? "ok" : "FAIL"} (${JSON.stringify(st)})`);
  // Windows high contrast: forced colors must not paint the text layer over the canvas
  await page.emulateMedia({ forcedColors: "active" });
  const forced = await page.evaluate(() => window.bdfForcedColors());
  await page.emulateMedia({ forcedColors: "none" });
  const forcedOk = forced.forced && forced.color === "rgba(0, 0, 0, 0)";
  if (!forcedOk) failed++;
  console.log(`forced colors: ${forcedOk ? "ok" : "FAIL"} (text layer color ${forced.color})`);
  const docPage = 0;
  const searchOk = search.hits.length === 6 && search.rects.length === 6 && search.rects[0].length === 2
    && search.rects[0].every((r) => r.a === docPage && r.w > 5 && r.h > 5 && r.x >= 72 && r.x + r.w <= 595.3 - 72)
    && search.sheetHits.length === 1 && search.sheetRects[0].length === 1 && Math.abs(search.sheetRects[0][0].y - 149 * 20) < 10
    && search.pptxHits.length === 1 && search.pptxRects[0].length === 2 && search.pptxRects[0].every((r) => r.a === 3)
    && search.ligHits.length > 0 && search.ligRects.every((rs) => rs.length > 0 && rs.every((r) => r.w > 5 && r.h > 5))
    && search.xlsxHits.length === 1 && search.xlsxRects[0].length === 1 && Math.abs(search.xlsxRects[0][0].y - 2055) < 30
    && search.drawioHits.length === 1 && search.drawioRects[0].length === 1 && search.drawioRects[0][0].w > 20
    && eq(search.drawioLinks, ["#view=details"])
    && search.mathHits.length === 2 && search.mathRects.every((rs) => rs.length === 1 && rs[0].w > 20 && rs[0].h > 5);
  if (!searchOk) failed++;
  console.log(`search: ${searchOk ? "ok" : "FAIL"} (${search.hits.length} hits, first hit rects ${JSON.stringify(search.rects[0])}, sheet ${JSON.stringify(search.sheetRects[0])}, pptx ${JSON.stringify(search.pptxRects[0])}, ligature ${JSON.stringify(search.ligRects[0])}, xlsx ${JSON.stringify(search.xlsxRects)}, drawio ${JSON.stringify(search.drawioRects[0])} links ${JSON.stringify(search.drawioLinks)}, math ${JSON.stringify(search.mathRects)})`);
  for (const r of results) {
    const renderOk = [r.main, r.worker].every((output) =>
      output.width === r.size[0] && output.height === r.size[1] && output.hasContent === r.expectsContent);
    if (!renderOk) failed++;
    console.log(`render ${r.name}: ${renderOk ? "ok" : "FAIL"} (${JSON.stringify(r)})`);
  }
} finally {
  await browser?.close();
  server.close();
}
if (failed) {
  console.error(`${failed} browser check(s) failed`);
  process.exit(1);
}
