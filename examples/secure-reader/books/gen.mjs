// Makes the two books secure-reader serves: PDFs long enough that a reader
// meets several segments (of 10 pages) in each. Every page says which
// segment it travels in. They are printed by headless Chromium where only
// free fonts are installed, as site/demo/gen.mjs does: in the Playwright
// image, which has the Liberation fonts (SIL OFL).
//
//   docker run --rm -v "$PWD":/work -w /work mcr.microsoft.com/playwright:v1.56.1-noble \
//     sh -c 'npm ci --ignore-scripts >/dev/null && node examples/secure-reader/books/gen.mjs'
import { chromium } from "playwright-core";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const SEGMENT = 10;

const BOOKS = [
  {
    file: "reading-room.pdf",
    title: "The Reading Room",
    subtitle: "A book that is never sent whole",
    color: "#1d6b40",
    chapters: [
      ["Opening the book", [
        "You logged in, and the server knows which books are yours. It does not send this one as a file: it sends ten pages at a time, as you come to them.",
        "Each request carries a public key that your browser made for that request alone. The private key never leaves the page's worker, and WebCrypto will not hand it out.",
      ]],
      ["What a segment is", [
        "A segment is a small bdf document: every page of the book is in it, with its size, but only ten of them have anything to draw. Its parts are those ten pages need.",
        "Fonts and pictures that several pages share come with the first segment that needs them. The next request says which segments you hold, and the server leaves them out.",
      ]],
      ["Sealed for one request", [
        "The server makes a key pair of its own for each segment, agrees a key with yours, wraps the segment's content key with it, and forgets its private key before it answers.",
        "Your worker agrees the same key, opens the segment, and drops its private key too. After that no one can open what was sent, not even the server.",
      ]],
      ["What a recording holds", [
        "Someone who recorded the traffic holds sealed segments. Stealing the server's TLS key later, or your password, or your session cookie, opens none of them.",
        "The session only lets its holder ask for new segments, as you could: the server checks your rights, and how fast pages are read, before it seals anything.",
      ]],
      ["What stays on the device", [
        "Nothing is written to the disk: the answers are sent with Cache-Control no-store, and the keys live in the worker's memory while the book is open.",
        "Close the tab, or log out, and there is nothing left to open. This site does not read books offline; that would need a key kept on the device.",
      ]],
      ["What this does not do", [
        "A reader who may read the book can still copy what the screen shows. What the server can do is decide who reads which pages, and how fast.",
        "The server keeps the book in the clear, converted once. Whoever takes the server's disk takes the books; that is a matter for the server, not for the wire.",
      ]],
    ],
    pagesPerChapter: 5,
  },
  {
    file: "field-notes.pdf",
    title: "Field Notes",
    subtitle: "The pages a sample shows",
    color: "#6b3d1d",
    chapters: [
      ["The first segment", [
        "A reader who does not own this book may read its first segment: the first ten pages. The server asks Handler.Allow before it seals a segment.",
        "The pages after them are in the outline, with their sizes, so the viewer lays the whole book out. Asking for them is answered 403.",
      ]],
      ["The pages after", [
        "Only a reader who owns the book gets these pages. The segment that holds them is sealed for that reader's request like any other.",
        "Nothing about this page was sent to a reader of the sample, not even its fonts: parts go with the pages that need them.",
      ]],
      ["Reading fast", [
        "The server counts the segments each reader asks for. Asking for many in a short time is answered 429, and the viewer says to slow down.",
        "A person turning pages does not come near the limit; a script that walks through the whole book does.",
      ]],
    ],
    pagesPerChapter: 7,
  },
];

const esc = (s) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;");

function bookHTML(b) {
  const pages = [];
  pages.push(`<section class="page cover"><div><h1>${esc(b.title)}</h1><p class="lead">${esc(b.subtitle)}</p><p class="small">A sample book of bdf's secure-reader example</p></div></section>`);
  pages.push(`<section class="page"><h2>Contents</h2><ol>${b.chapters.map(([t], i) => `<li>${esc(t)} <span>${3 + i * b.pagesPerChapter}</span></li>`).join("")}</ol></section>`);
  b.chapters.forEach(([title, paras], ci) => {
    for (let k = 0; k < b.pagesPerChapter; k++) {
      const n = pages.length + 1;
      const seg = Math.floor((n - 1) / SEGMENT);
      const body = k === 0 ? `<h2>${ci + 1}. ${esc(title)}</h2>${paras.map((p) => `<p>${esc(p)}</p>`).join("")}` : `<h3>${esc(title)}, continued</h3><p>${esc(paras[k % paras.length])}</p>`;
      pages.push(`<section class="page">${body}<div class="seg">This is page ${n}. It travels in segment ${seg + 1}, with pages ${seg * SEGMENT + 1} to ${seg * SEGMENT + SEGMENT}, sealed for the request that asked for it.</div></section>`);
    }
  });
  return `<!doctype html><html lang="en"><head><meta charset="utf-8"><title>${esc(b.title)}</title><meta name="author" content="bdf">
<style>
@page { size: 148mm 210mm; margin: 0; }
* { box-sizing: border-box; }
html { font-family: "Liberation Serif", serif; font-size: 11pt; line-height: 1.55; color: #222; }
body { margin: 0; counter-reset: page; }
.page { position: relative; width: 148mm; height: 210mm; padding: 22mm 18mm 26mm; overflow: hidden; page-break-after: always; counter-increment: page; }
.page:last-child { page-break-after: auto; }
.page::after { content: counter(page); position: absolute; bottom: 12mm; left: 0; right: 0; text-align: center; font: 9pt "Liberation Sans", sans-serif; color: #777; }
.cover { display: flex; align-items: center; background: ${b.color}; color: #fff; }
.cover::after { content: none; }
h1 { font: 700 26pt/1.2 "Liberation Sans", sans-serif; margin: 0 0 4mm; }
h2 { font: 700 16pt/1.3 "Liberation Sans", sans-serif; color: ${b.color}; margin: 0 0 6mm; }
h3 { font: 700 11pt/1.3 "Liberation Sans", sans-serif; color: #555; margin: 0 0 5mm; }
.lead { font-size: 13pt; }
.small { font: 9pt "Liberation Sans", sans-serif; opacity: .8; }
ol { padding-left: 5mm; } li { margin: 2mm 0; } li span { float: right; }
.seg { position: absolute; left: 18mm; right: 18mm; bottom: 20mm; padding: 3mm 4mm; border-left: 2.5pt solid ${b.color}; background: #f4f4f1; font: 9pt/1.45 "Liberation Sans", sans-serif; color: #444; }
</style></head><body>${pages.join("\n")}</body></html>`;
}

const browser = await chromium.launch(process.env.CHROMIUM_PATH ? { executablePath: process.env.CHROMIUM_PATH } : {});
const page = await browser.newPage();
for (const b of BOOKS) {
  await page.setContent(bookHTML(b));
  await page.waitForTimeout(200);
  const path = join(here, b.file);
  await page.pdf({ path, preferCSSPageSize: true, printBackground: true, tagged: true });
  console.log("wrote", path);
}
await browser.close();
