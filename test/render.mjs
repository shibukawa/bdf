// Renders a .bdf file to PNGs with headless Chromium: node test/render.mjs file.bdf outdir [scale]
import { build } from "esbuild";
import { chromium } from "playwright-core";
import { writeFile, mkdir, copyFile } from "node:fs/promises";
import { existsSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join, resolve, relative, basename } from "node:path";
import { serve } from "./serve.mjs";

const root = dirname(dirname(fileURLToPath(import.meta.url)));
const [file, outDir, scaleArg] = process.argv.slice(2);
if (!file || !outDir) { console.error("usage: node test/render.mjs file.bdf outdir [scale]"); process.exit(2); }
const scale = Number(scaleArg ?? "1");
const bundleDir = join(root, "test", ".out");
await mkdir(bundleDir, { recursive: true });
await mkdir(outDir, { recursive: true });
await build({ entryPoints: [{ in: join(root, "test/page/render.ts"), out: "render" }], bundle: true, format: "esm", outdir: bundleDir, target: "es2022", logLevel: "warning" });
await writeFile(join(bundleDir, "render.html"), `<!doctype html><meta charset="utf-8"><title>render</title><script type="module" src="./render.js"></script>`);
// Serve the directory that contains the file (may be outside the repo).
const abs = resolve(file);
const serveRoot = abs.startsWith(root) ? root : dirname(abs);
const { server, port } = await serve(root);
let fileServer;
let src;
if (serveRoot === root) src = "/" + relative(root, abs).split("\\").join("/");
else {
  fileServer = await serve(dirname(abs));
  src = `http://127.0.0.1:${fileServer.port}/${basename(abs)}`;
}
const executablePath = process.env.CHROMIUM_PATH ?? ["/opt/pw-browsers/chromium-1194/chrome-linux/chrome"].find((p) => existsSync(p));
// The pixels must not depend on the machine: Canvas 2D is rasterized in software (no GPU or
// SwiftShader path, which also makes the worker's OffscreenCanvas match the main thread), text
// is drawn without hinting or LCD filtering so the fontconfig defaults (hint style, subpixel
// order) do not leak in, and playwright-core is pinned to the release whose Chromium build
// matches the goldens.
const launch = { args: ["--disable-gpu", "--disable-accelerated-2d-canvas", "--font-render-hinting=none", "--disable-lcd-text"] };
if (executablePath) launch.executablePath = executablePath;
const browser = await chromium.launch(launch);
try {
  const page = await browser.newPage();
  page.on("console", (m) => { if (m.type() === "error") console.error("[browser]", m.text()); });
  await page.goto(`http://127.0.0.1:${port}/test/.out/render.html?src=${encodeURIComponent(src)}&scale=${scale}`);
  await page.waitForFunction(() => document.title === "done" || document.title === "error", null, { timeout: 120000 });
  const error = await page.evaluate(() => window.bdfError);
  if (error) throw new Error(error);
  const pages = await page.evaluate(() => window.bdfPages);
  for (const p of pages) {
    const path = join(outDir, `${p.name}.png`);
    await writeFile(path, Buffer.from(p.png.split(",")[1], "base64"));
    console.log("wrote", path);
  }
} finally {
  await browser.close();
  server.close();
  fileServer?.server.close();
}
