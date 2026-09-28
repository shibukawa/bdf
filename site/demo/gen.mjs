// Makes demo.pdf, the document the top page of the site shows, of demo.html
// with headless Chromium. The PDF embeds the fonts it is laid out with, so
// it is made where only free fonts are installed: in the Playwright image,
// which has Liberation (SIL OFL) and IPA (IPA Font License) fonts.
//
//   docker run --rm -v "$PWD":/work -w /work mcr.microsoft.com/playwright:v1.56.1-noble \
//     sh -c 'npm ci --ignore-scripts >/dev/null && node site/demo/gen.mjs'
import { chromium } from "playwright-core";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const browser = await chromium.launch(process.env.CHROMIUM_PATH ? { executablePath: process.env.CHROMIUM_PATH } : {});
const page = await browser.newPage();
await page.goto("file://" + join(here, "demo.html"));
await page.waitForTimeout(300);
const path = join(here, "demo.pdf");
await page.pdf({ path, preferCSSPageSize: true, printBackground: true, tagged: true });
console.log("wrote", path);
await browser.close();
