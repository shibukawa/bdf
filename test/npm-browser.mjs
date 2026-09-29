// Browser smoke test for the npm entry points and their Worker assets.
// npm run build && npm run build:wasm && node test/npm-browser.mjs
import assert from "node:assert/strict";
import { mkdtemp, rm, writeFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { chromium } from "playwright-core";
import { serve } from "./serve.mjs";

const root = resolve(import.meta.dirname, "..");
const tmp = await mkdtemp(join(root, "test/npm-browser-"));
const html = `<!doctype html><meta charset="utf-8">
<script type="importmap">{"imports":{
  "@bdfkit/core":"/packages/core/dist/index.js",
  "@bdfkit/render":"/packages/render/dist/index.js",
  "@bdfkit/convert":"/packages/convert/dist/index.js"
}}</script>
<div id="preview" style="width:800px;height:600px"></div>
<script type="module">
import { Viewer } from "/packages/viewer/dist/index.js";
import { pianoRoll } from "/packages/viewer/dist/music.js";
import { createPdfEpubConverter } from "/packages/convert-pdf-epub/dist/index.js";
import { createAllConverter } from "/packages/convert-all/dist/index.js";
const viewer = new Viewer(document.querySelector("#preview"), { music: pianoRoll });
const converter = createPdfEpubConverter();
window.test = { viewer, converter, Viewer, createAllConverter };
await viewer.open({ kind: "single", url: "/testdata/epub/basic.bdf" });
window.loaded = true;
</script>`;
await writeFile(join(tmp, "index.html"), html);
let browser, server, port;
try {
  ({ server, port } = await serve(root));
  browser = await chromium.launch({ executablePath: process.env.BDF_CHROME ?? "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", headless: true });
  const page = await browser.newPage();
  const errors = [];
  page.on("pageerror", (error) => errors.push(error.message));
  await page.goto(`http://127.0.0.1:${port}/${tmp.slice(root.length + 1)}/index.html`);
  await page.waitForFunction(() => window.loaded && document.querySelector("#preview canvas")?.width > 0, null, { timeout: 30000 });
  const result = await page.evaluate(async () => {
    const { viewer, converter, Viewer, createAllConverter } = window.test;
    const viewCount = viewer.document.views.length;
    viewer.setLayout("spread");
    const rows = viewer.stage.children.length;
    const formats = (await converter.ready).map((format) => format.name);
    const pdf = new Uint8Array(await (await fetch("/examples/sample-files/report.pdf")).arrayBuffer());
    const converted = await converter.convert(pdf, { name: "report.pdf" });
    const finished = new Promise((resolve, reject) => {
      viewer.addEventListener("conversiondone", () => resolve(true), { once: true });
      viewer.addEventListener("error", (event) => reject(event.detail), { once: true });
    });
    await viewer.openFile(converter, pdf, { name: "report.pdf" });
    await finished;
    const pdfViews = viewer.document.views.length;
    await viewer.open({ kind: "single", url: "/testdata/csv/basic.bdf" });
    await new Promise((resolve, reject) => {
      const until = Date.now() + 10000;
      const check = () => {
        if (viewer.stage.querySelector("canvas")?.width > 0) resolve(true);
        else if (Date.now() > until) reject(new Error("sheet did not render"));
        else setTimeout(check, 50);
      };
      check();
    });
    const sheetKind = viewer.view.kind;
    await viewer.open({ kind: "single", url: "/testdata/music/minuet-musicxml.bdf" });
    viewer.setMusicMode("piano-roll");
    await new Promise((resolve, reject) => {
      const until = Date.now() + 10000;
      const check = () => {
        if (viewer.stage.querySelector("canvas[aria-label^='Piano roll']")) resolve(true);
        else if (Date.now() > until) reject(new Error("piano roll did not render"));
        else setTimeout(check, 50);
      };
      check();
    });
    const lightbox = new Viewer(document.body, { mode: "lightbox", renderer: viewer.renderer });
    lightbox.show();
    const lightboxShown = lightbox.surface.style.display === "block";
    lightbox.hide();
    lightbox.destroy();
    const basic = new Viewer(document.body, { renderer: viewer.renderer });
    let musicOff = false;
    try { basic.setMusicMode("piano-roll"); } catch { musicOff = true; }
    basic.destroy();
    const all = createAllConverter();
    const allFormats = (await all.ready).map((format) => format.name);
    const csv = new Uint8Array(await (await fetch("/examples/sample-files/basic.csv")).arrayBuffer());
    const allConverted = await all.convert(csv, { name: "basic.csv" });
    all.terminate();
    viewer.destroy();
    converter.terminate();
    return { viewCount, rows, formats, pdfViews, sheetKind, lightboxShown, musicOff, allFormats, allFormat: allConverted.format, magic: [...converted.bdf.slice(0, 4)] };
  });
  assert.ok(result.viewCount > 0);
  assert.ok(result.rows > 0);
  assert.deepEqual(result.formats, ["epub", "pdf"]);
  assert.ok(result.pdfViews > 0);
  assert.equal(result.sheetKind, "sheet");
  assert.equal(result.lightboxShown, true);
  assert.equal(result.musicOff, true);
  assert.equal(result.allFormats.length, 27);
  assert.ok(result.allFormats.includes("tiff"));
  assert.equal(result.allFormat, "csv");
  assert.deepEqual(result.magic, [98, 100, 102, 0]);
  assert.deepEqual(errors, []);
  console.log("npm browser smoke test passed");
} finally {
  await browser?.close();
  await new Promise((resolve) => server?.close(resolve) ?? resolve());
  await rm(tmp, { recursive: true, force: true });
}
