// Renders the fixture in headless Chromium and compares against fixtures/golden/*.png.
// Usage: node test/golden.mjs [--update] [--source single|split|range]
import { build } from "esbuild";
import { chromium } from "playwright-core";
import { writeFile, mkdir, copyFile } from "node:fs/promises";
import { existsSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";
import { serve } from "./serve.mjs";

const root = dirname(dirname(fileURLToPath(import.meta.url)));
const args = process.argv.slice(2);
const update = args.includes("--update");
const source = args.includes("--source") ? args[args.indexOf("--source") + 1] : "single";
const PIXEL_TOLERANCE = 0.02; // fraction of pixels allowed to differ from the golden

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
// Use an explicit Chromium (CHROMIUM_PATH, or the one preinstalled in some sandboxes);
// otherwise fall back to the browser installed by `npx playwright-core install chromium`.
const executablePath = process.env.CHROMIUM_PATH ?? ["/opt/pw-browsers/chromium-1194/chrome-linux/chrome"].find((p) => existsSync(p));
const browser = await chromium.launch(executablePath ? { executablePath } : {});
let failed = 0;
try {
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
  const docPage = 0;
  const searchOk = search.hits.length === 6 && search.rects.length === 6 && search.rects[0].length === 2
    && search.rects[0].every((r) => r.a === docPage && r.w > 5 && r.h > 5 && r.x >= 72 && r.x + r.w <= 595.3 - 72)
    && search.sheetHits.length === 1 && search.sheetRects[0].length === 1 && Math.abs(search.sheetRects[0][0].y - 149 * 20) < 10;
  if (!searchOk) failed++;
  console.log(`search: ${searchOk ? "ok" : "FAIL"} (${search.hits.length} hits, first hit rects ${JSON.stringify(search.rects[0])}, sheet ${JSON.stringify(search.sheetRects[0])})`);
  for (const r of results) {
    const goldenPath = join(root, "fixtures/golden", `${r.name}.png`);
    const png = Buffer.from(r.png.split(",")[1], "base64");
    let status;
    if (update || !r.golden) {
      await writeFile(goldenPath, png);
      status = update ? "updated" : "created";
    } else {
      const frac = r.golden.different / r.golden.total;
      status = frac <= PIXEL_TOLERANCE ? "ok" : "FAIL";
      if (status === "FAIL") {
        failed++;
        await writeFile(join(outDir, `${r.name}.actual.png`), png);
        await writeFile(join(outDir, `${r.name}.golden-diff.png`), Buffer.from(r.golden.png.split(",")[1], "base64"));
      }
      status += ` golden diff ${(frac * 100).toFixed(3)}% (max delta ${r.golden.maxDelta})`;
    }
    // OffscreenCanvas in a worker antialiases clip edges slightly differently from a
    // main-thread canvas, so only the fraction of differing pixels is checked here.
    const workerOk = r.worker.different / r.worker.total <= PIXEL_TOLERANCE;
    if (!workerOk) {
      failed++;
      await writeFile(join(outDir, `${r.name}.main.png`), png);
      await writeFile(join(outDir, `${r.name}.worker.png`), Buffer.from(r.workerPng.split(",")[1], "base64"));
      await writeFile(join(outDir, `${r.name}.worker-diff.png`), Buffer.from(r.worker.png.split(",")[1], "base64"));
    }
    console.log(`${r.name.padEnd(22)} ${status}; worker ${r.worker.different === 0 ? "identical" : `${workerOk ? "close" : "DIFFERS"} (${r.worker.different} px, max delta ${r.worker.maxDelta})`}; text runs ${r.textRuns}`);
  }
} finally {
  await browser.close();
  server.close();
}
if (failed) {
  console.error(`${failed} check(s) failed; actual renders are in test/.out/`);
  process.exit(1);
}
