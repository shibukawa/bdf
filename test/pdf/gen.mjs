// Generates test PDFs from HTML with headless Chromium (Skia PDF backend).
import { chromium } from "playwright-core";
import { existsSync } from "node:fs";
import { readdir } from "node:fs/promises";
import { dirname, join, basename } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const root = dirname(dirname(here));
const out = join(root, "pdf2bdf/testdata");
const executablePath = process.env.CHROMIUM_PATH ?? ["/opt/pw-browsers/chromium-1194/chrome-linux/chrome"].find((p) => existsSync(p));
const browser = await chromium.launch(executablePath ? { executablePath } : {});
const page = await browser.newPage();
for (const f of await readdir(join(here, "html"))) {
  if (!f.endsWith(".html")) continue;
  await page.goto("file://" + join(here, "html", f));
  await page.waitForTimeout(300);
  const path = join(out, `chrome-${basename(f, ".html")}.pdf`);
  await page.pdf({ path, preferCSSPageSize: true, printBackground: true, tagged: true });
  console.log("wrote", path);
}
await browser.close();
