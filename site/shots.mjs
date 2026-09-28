// Takes the pictures of documents in the viewer that the top page and the
// documentation show (docs/images/*.webp): it opens the samples of the built
// site (npm run site) in the demo viewer with Chromium and captures the
// document and the tabs of its views, 800 × 500 CSS pixels at twice the
// resolution.
//
//   node site/shots.mjs [name ...]      (all of them by default)
//
// Needs a Chromium (CHROMIUM_PATH, or the one playwright-core installed)
// and cwebp. The pictures are committed: fonts and rasterizers differ
// between machines, so they are taken again only when what they show
// changed.
import { execFile as execFileCb } from "node:child_process";
import { mkdir, mkdtemp, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { promisify } from "node:util";
import { chromium } from "playwright-core";
import { serve } from "../test/serve.mjs";
import { root } from "../examples/common/build.mjs";

const execFile = promisify(execFileCb);
const WIDTH = 800, HEIGHT = 500, SCALE = 2;

/**
 * The pictures: the sample shown, and how. zoom is the viewer's zoom ("fit"
 * makes the widest page as wide as the picture); view the tab to show (its
 * number, or part of its label); layout the viewer's layout and page how
 * many times its pages are turned; top how far down the document the
 * picture starts, in CSS pixels.
 */
const SHOTS = [
  { name: "pdf", file: "demo.pdf", zoom: "fit" },
  { name: "ai", file: "artboards.ai", zoom: 1.25 },
  { name: "docx", file: "basic.docx", zoom: "fit" },
  { name: "docx-vertical", file: "vertical.docx", zoom: "fit" },
  { name: "docx-math", file: "math.docx", zoom: "fit" },
  { name: "markdown", file: "basic.md", zoom: 1 },
  { name: "pptx", file: "features.pptx", zoom: "fit" },
  { name: "xlsx", file: "features.xlsx", zoom: 1 },
  { name: "csv", file: "japanese.tsv", zoom: 1 },
  { name: "parquet", file: "basic.parquet", zoom: 1 },
  { name: "visio", file: "shapes.vsdx", zoom: "fit" },
  { name: "drawio", file: "aws.drawio", zoom: "fit" },
  { name: "drawio-pages", file: "multipage.drawio", zoom: "fit" },
  { name: "dxf", file: "layout.dxf", zoom: "fit" },
  { name: "dxf-layout", file: "layout.dxf", zoom: "fit", view: 1 },
  { name: "jww", file: "shapes.jww", zoom: "fit" },
  { name: "sxf", file: "shapes.p21", zoom: "fit" },
  { name: "cgm", file: "shapes.cgm", zoom: "fit" },
  { name: "hpgl", file: "shapes.plt", zoom: "fit" },
  { name: "gerber", file: "board.zip", zoom: "fit" },
  { name: "gerber-layers", file: "board.zip", zoom: "fit", view: "Top copper" },
  { name: "kicad", file: "demo.zip", zoom: "fit" },
  { name: "kicad-board", file: "demo.zip", zoom: "fit", view: "Front" },
  { name: "epub", file: "vertical.epub", layout: "spread", page: 2 },
  { name: "epub-fixed", file: "fixed.epub", layout: "spread", page: 1 },
  { name: "music", file: "minuet.musicxml", zoom: "fit" },
  { name: "music-mml", file: "frere.mml", zoom: "fit" },
  { name: "font", file: "stix.otf", zoom: 1 },
  { name: "font-features", file: "stix.otf", zoom: 1, view: "Features" },
  { name: "psd", file: "artboards.psd", zoom: 1.25 },
  { name: "image", file: "drawing.svg", zoom: "fit" },
];

const site = join(root, "site/dist");
const out = join(root, "docs/images");
const only = process.argv.slice(2);
const shots = only.length ? SHOTS.filter((s) => only.includes(s.name)) : SHOTS;
await mkdir(out, { recursive: true });
const tmp = await mkdtemp(join(tmpdir(), "bdf-shots-"));
const { server, port } = await serve(site);
const browser = await chromium.launch(process.env.CHROMIUM_PATH ? { executablePath: process.env.CHROMIUM_PATH } : {});
try {
  for (const s of shots) {
    const context = await browser.newContext({ viewport: { width: WIDTH, height: HEIGHT + 200 }, deviceScaleFactor: SCALE, colorScheme: "light" });
    const page = await context.newPage();
    page.on("pageerror", (e) => console.error(`${s.name}: ${e.message}`));
    const layout = s.layout ? `&layout=${s.layout}` : "";
    await page.goto(`http://127.0.0.1:${port}/viewer/?file=samples/${s.file}${layout}`);
    // the document is open when its first view is shown
    await page.waitForFunction(() => /converted in/.test(document.getElementById("status").textContent), null, { timeout: 120000 });
    // the header wraps at this width: the document gets the height of the picture under it
    const fitViewport = async () => {
      const header = await page.evaluate(() => document.querySelector("header").offsetHeight);
      await page.setViewportSize({ width: WIDTH, height: header + HEIGHT });
    };
    await fitViewport();
    if (s.view !== undefined) {
      const tabs = page.locator("#tabs [role=tab]");
      await (typeof s.view === "number" ? tabs.nth(s.view) : tabs.filter({ hasText: s.view }).first()).click();
    }
    if (s.zoom !== undefined) {
      const zoom = s.zoom === "fit" ? await page.evaluate(() => {
        const stage = document.getElementById("stage");
        // pages are laid out at zoom 1 when the document opens
        const widest = Math.max(0, ...[...stage.querySelectorAll(".page")].map((p) => p.offsetWidth));
        return widest ? Math.min(1.5, (stage.clientWidth - 32) / widest) : 1;
      }) : s.zoom;
      await page.evaluate((z) => {
        const input = document.getElementById("zoom");
        input.step = "any";
        input.min = "0.05";
        input.value = String(z);
        input.dispatchEvent(new Event("input"));
      }, zoom);
    }
    if (s.page) for (let i = 0; i < s.page; i++) await page.click("#nextPage");
    await fitViewport();
    if (s.top) await page.evaluate((y) => document.getElementById("stage").scrollTo(0, y), s.top);
    // pages and sheet chunks are drawn by the worker as they come into view
    await page.waitForFunction(() => {
      const stage = document.getElementById("stage");
      const box = stage.getBoundingClientRect();
      const inView = [...stage.querySelectorAll(".page")].filter((p) => {
        const r = p.getBoundingClientRect();
        return r.bottom > box.top && r.top < box.bottom;
      });
      return inView.length ? inView.every((p) => p.querySelector("canvas")) : !!stage.querySelector("canvas");
    }, null, { timeout: 60000 });
    await page.mouse.move(0, 0);
    await page.waitForTimeout(1200);
    const clip = await page.evaluate(() => {
      const top = document.getElementById("stage").getBoundingClientRect().top;
      return { x: 0, y: top, width: innerWidth, height: innerHeight - top };
    });
    const png = join(tmp, `${s.name}.png`);
    await page.screenshot({ path: png, clip });
    await execFile("cwebp", ["-quiet", "-q", "82", "-m", "6", png, "-o", join(out, `${s.name}.webp`)]);
    console.log(`${s.name}.webp: ${s.file}, ${clip.width * SCALE} × ${clip.height * SCALE}`);
    await context.close();
  }
} finally {
  await browser.close();
  server.close();
  await rm(tmp, { recursive: true, force: true });
}
