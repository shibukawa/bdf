// Makes demo.pdf, the document the top page of the site shows, of demo.html
// with headless Chromium. The PDF embeds the fonts it is laid out with, so
// it is made where only free fonts are installed: in the Playwright image,
// which has Liberation (SIL OFL) and IPA (IPA Font License) fonts.
//
//   docker run --rm -v "$PWD":/work -w /work mcr.microsoft.com/playwright:v1.56.1-noble \
//     sh -c 'npm ci --ignore-scripts >/dev/null && node site/demo/gen.mjs'
//
// A page whose content runs past its foot is cut off in the PDF: the script
// refuses to write one.
import { chromium } from "playwright-core";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const browser = await chromium.launch(process.env.CHROMIUM_PATH ? { executablePath: process.env.CHROMIUM_PATH } : {});
const page = await browser.newPage();
await page.goto("file://" + join(here, "demo.html"));
await page.waitForTimeout(300);
// the pages, and the ones whose content does not fit (overflow: hidden would cut it)
const pages = await page.evaluate(() => [...document.querySelectorAll(".page")].map((p, i) => {
  const foot = p.querySelector(".footer");
  const bottom = foot ? foot.getBoundingClientRect().top - p.getBoundingClientRect().top : p.clientHeight;
  // the lowest edge of the content above the footer
  const low = Math.max(0, ...[...p.children].filter((c) => c !== foot && !c.classList.contains("band")).map((c) => c.getBoundingClientRect().bottom - p.getBoundingClientRect().top));
  return { page: i + 1, room: Math.round(bottom - low), over: p.scrollHeight > p.clientHeight || low > bottom };
}));
for (const p of pages) console.log(`page ${p.page}: ${p.room} px to the footer${p.over ? " — OVERFLOWS" : ""}`);
const over = pages.filter((p) => p.over);
if (over.length) {
  await browser.close();
  throw new Error(`the content of page(s) ${over.map((p) => p.page).join(", ")} does not fit`);
}
const path = join(here, "demo.pdf");
await page.pdf({ path, preferCSSPageSize: true, printBackground: true, tagged: true });
console.log("wrote", path, `(${pages.length} pages)`);
await browser.close();
