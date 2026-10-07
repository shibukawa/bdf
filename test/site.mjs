// Smoke test of the demo site (site/build.mjs): serves it, runs
// its converter modules in Node with Go's wasm_exec.js, and converts each
// sample with the site's fonts, fetched from the server as in a browser.
// PDFs are also converted a page at a time, as the viewer does, and the
// pages put into the outline with @bdfkit/core (npm run build first). The
// preview module draws the thumbnail of each converted sample and gives its
// text.
//
//   node test/site.mjs [site dir]
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { pathToFileURL } from "node:url";
import { BdfDocument, BufferSource, extractText } from "../packages/core/dist/index.js";
import { serve } from "./serve.mjs";

const site = resolve(process.argv[2] ?? "site/dist");
await import(pathToFileURL(join(site, "lib/wasm_exec.js")).href);
const { server, port } = await serve(site);
const base = `http://127.0.0.1:${port}/`;

/** Load a module; the Go program puts its API on globalThis and waits. */
async function load(file) {
  const go = new globalThis.Go();
  const { instance } = await WebAssembly.instantiate(await readFile(join(site, "lib", file)), go.importObject);
  go.run(instance);
  const conv = globalThis.bdfConverter;
  delete globalThis.bdfConverter;
  return conv;
}
const modules = {
  pdf: await load("bdf-pdf.wasm"), office: await load("bdf-office.wasm"), image: await load("bdf-image.wasm"), web: await load("bdf-web.wasm"),
  preview: await load("bdf-preview.wasm"),
};
assert.deepEqual(modules.pdf.formats.map((f) => f.name), ["ai", "pdf"]);
assert.deepEqual(modules.office.formats.map((f) => f.name), ["cgm", "csv", "docx", "drawio", "dxf", "emf", "font", "gerber", "hpgl", "image", "jww", "kicad", "midi", "mml", "musicxml", "parquet", "pptx", "psd", "sxf", "visio", "xlsx"]);
assert.deepEqual(modules.image.formats.map((f) => f.name), ["image"]);
assert.deepEqual(modules.web.formats.map((f) => f.name), ["epub", "html", "markdown"]);
assert.deepEqual(modules.preview.formats, []);
const imageExtensions = modules.image.formats[0].extensions;
const webExtensions = modules.web.formats.flatMap((f) => f.extensions);

let failed = 0;
/** The converted samples, for the preview module. */
const converted = [];
const samples = JSON.parse(await readFile(join(site, "samples/index.json"), "utf8"));
for (const { name } of samples) {
  if (name.endsWith(".bdf")) continue;
  const data = new Uint8Array(await readFile(join(site, "samples", name)));
  const image = imageExtensions.some((e) => name.endsWith(e));
  const web = webExtensions.some((e) => name.endsWith(e));
  const conv = /\.(pdf|ai)$/.test(name) ? modules.pdf : image ? modules.image : web ? modules.web : modules.office;
  const t0 = performance.now();
  try {
    const res = await conv.convert(data, { fonts: image ? undefined : `${base}fonts/`, name });
    assert.deepEqual([...res.bdf.subarray(0, 4)], [0x62, 0x64, 0x66, 0], "bdf magic");
    // the converters that lay text out found the site's fonts and embedded them (HTML, Markdown and EPUB refer to
    // them by name, as a web page does)
    if (conv !== modules.image && conv !== modules.web && !["pdf", "ai", "psd", "gerber"].includes(res.format)) {
      assert.match(res.summary, /[1-9]\d* embedded font/);
    }
    converted.push({ name, format: res.format, bdf: res.bdf });
    console.log(`ok   ${name}: ${res.format}, ${res.summary}, ${res.bdf.length} bytes, ${(performance.now() - t0).toFixed(0)} ms`);
    for (const w of res.warnings) console.log(`     warning: ${w}`);
  } catch (e) {
    failed++;
    console.log(`FAIL ${name}: ${e.message}`);
  }
}

/** The text of every page of the first view, drawn from the parts the document holds. */
async function pageTexts(doc) {
  const texts = [];
  for (const page of doc.manifest.views[0].pages) {
    let text = "";
    for (const l of page.layers) {
      const obj = await doc.ensure(l.obj);
      for (const r of extractText(obj, (h) => doc.objectSync(h))) text += r.text;
    }
    texts.push(text);
  }
  return texts;
}

/**
 * Convert a PDF a page at a time in the order given, putting each page into
 * the outline as the viewer does, and compare with the whole conversion:
 * the same text on every page, and the same bytes when the pages came in
 * order.
 */
async function streamed(name, data, order) {
  const whole = await modules.pdf.convert(data, {});
  const opened = await modules.pdf.open(data, {});
  const n = opened.pages;
  assert.ok(n > 0 && opened.stream, "a PDF streams");
  const doc = await BdfDocument.open(new BufferSource(opened.bdf));
  assert.equal(doc.manifest.views[0].pages.length, n);
  assert.ok(doc.manifest.views[0].pages.every((p) => p.layers.length === 0), "the outline has no layers");
  const pages = order(n);
  for (const i of pages) {
    const page = await opened.stream.page(i);
    doc.addPage(doc.manifest.views[0].id, i, await BdfDocument.open(new BufferSource(page.bdf)));
  }
  const finished = await opened.stream.finish();
  opened.stream.close();
  const want = await pageTexts(await BdfDocument.open(new BufferSource(whole.bdf)));
  assert.deepEqual(await pageTexts(doc), want, `${name}: the streamed pages' text`);
  assert.deepEqual(await pageTexts(await BdfDocument.open(new BufferSource(finished.bdf))), want, `${name}: the finished document's text`);
  if (pages.every((p, i) => p === i)) assert.deepEqual(finished.bdf, whole.bdf, `${name}: finished in order, the same bytes as convert`);
  assert.equal(finished.summary, whole.summary);
}
for (const { name } of samples) {
  if (!name.endsWith(".pdf")) continue;
  const data = new Uint8Array(await readFile(join(site, "samples", name)));
  try {
    await streamed(name, data, (n) => [...Array(n).keys()]);
    await streamed(name, data, (n) => [...Array(n).keys()].reverse());
    console.log(`ok   ${name}: converted a page at a time, in order and reversed`);
  } catch (e) {
    failed++;
    console.log(`FAIL ${name} a page at a time: ${e.message}`);
  }
}
// formats converted whole come back whole from open
{
  const name = samples.find((s) => s.name.endsWith(".pptx")).name;
  const opened = await modules.office.open(new Uint8Array(await readFile(join(site, "samples", name))), { fonts: `${base}fonts/` });
  assert.equal(opened.pages, 0);
  assert.equal(opened.stream, undefined);
  assert.match(opened.summary, /slide/);
}

/** The width and height of a PNG image. */
function pngSize(b) {
  assert.deepEqual([...b.subarray(0, 8)], [0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a], "PNG signature");
  const dv = new DataView(b.buffer, b.byteOffset, b.byteLength);
  return [dv.getUint32(16), dv.getUint32(20)];
}
/** The layout the thumbnail package picks for the formats that have one. */
const LAYOUT = { docx: "crop", html: "crop", markdown: "crop", xlsx: "crop", csv: "crop", parquet: "crop", font: "crop", pptx: "fit", vsdx: "fit", drawio: "fit", epub: "fit", svg: "fit", jpeg: "fit" };
for (const { name, format, bdf } of converted) {
  const t0 = performance.now();
  try {
    const th = await modules.preview.thumbnail(bdf, { size: 64, fonts: `${base}fonts/` });
    assert.equal(th.format, "png");
    assert.deepEqual(pngSize(th.image), [th.width, th.height]);
    assert.equal(Math.max(th.width, th.height), 64);
    if (th.mode === "crop") assert.deepEqual([th.width, th.height], [64, 64]);
    if (LAYOUT[format]) assert.equal(th.mode, LAYOUT[format], "layout");
    const text = JSON.parse((await modules.preview.text(bdf)).json);
    assert.equal(text.meta.source, format);
    const pages = text.views.reduce((n, v) => n + v.pages.length, 0);
    console.log(`ok   ${name}: ${th.width}×${th.height} thumbnail (${th.mode}), text of ${pages} page(s), ${(performance.now() - t0).toFixed(0)} ms`);
    for (const w of th.warnings) console.log(`     warning: ${w}`);
  } catch (e) {
    failed++;
    console.log(`FAIL ${name} thumbnail and text: ${e.message}`);
  }
}
// JPEG, the other sizes, and an encrypted document with its password
{
  const plain = new Uint8Array(await readFile(new URL("../testdata/demo.bdf", import.meta.url)));
  const enc = new Uint8Array(await readFile(new URL("../testdata/demo-encrypted.bdf", import.meta.url)));
  const jpeg = await modules.preview.thumbnail(plain, { size: 256, format: "jpeg", mode: "fit" });
  assert.deepEqual([...jpeg.image.subarray(0, 3)], [0xff, 0xd8, 0xff], "JPEG signature");
  assert.equal(jpeg.mode, "fit");
  assert.equal(Math.max(jpeg.width, jpeg.height), 256);
  await assert.rejects(modules.preview.thumbnail(enc, {}), { code: "password-required" });
  await assert.rejects(modules.preview.text(enc, { password: "wrong" }), { code: "wrong-password" });
  const password = "demo-パスワード"; // package.json's testdata script
  const th = await modules.preview.thumbnail(enc, { size: 128, password });
  assert.deepEqual(pngSize(th.image), [th.width, th.height]);
  assert.deepEqual((await modules.preview.thumbnail(plain, { size: 128 })).image, th.image, "the encrypted document's thumbnail");
  assert.equal((await modules.preview.text(enc, { password })).json, (await modules.preview.text(plain)).json);
  await assert.rejects(modules.preview.thumbnail(plain, { size: 0 }), /size/);
  await assert.rejects(modules.preview.thumbnail(plain, { format: "gif" }), /format/);
  await assert.rejects(modules.preview.text(new Uint8Array([1, 2, 3])));
  console.log("ok   testdata/demo.bdf and demo-encrypted.bdf: JPEG, sizes, passwords");
}

// the home pages show the measured comparison below the viewer in both languages
for (const page of ["index.html", "index.ja.html"]) {
  const html = await readFile(join(site, page), "utf8");
  const introduction = page.endsWith(".ja.html")
    ? ["BDF は、ブラウザ向けの文書プレビュースイートです", "gzip 圧縮後約 23 KB", "PDF が PostScript をベースにした", "Canvas 2D の描画命令"]
    : ["BDF is a browser document-preview suite", "about 23 KB gzipped", "PDF is a portable, PostScript-based format for printing", "Canvas 2D drawing commands"];
  for (const phrase of introduction) assert.ok(html.includes(phrase), `${page}: introduction includes ${phrase}`);
  const viewer = html.indexOf('id="demo-title"');
  const benchmark = html.indexOf('id="benchmark-title"');
  const gallery = html.indexOf('id="gallery-title"');
  assert.ok(viewer >= 0 && viewer < benchmark && benchmark < gallery, `${page}: benchmark follows the viewer`);
  const timeChart = html.indexOf('id="benchmark-time-title"');
  const memoryChart = html.indexOf('id="benchmark-memory-title"');
  const energyChart = html.indexOf('id="benchmark-energy-title"');
  assert.ok(benchmark < timeChart && timeChart < memoryChart && memoryChart < energyChart && energyChart < gallery,
    `${page}: charts are ordered by time, memory and energy`);
  const energyRanges = page.endsWith(".ja.html")
    ? ["3.30–7.16、", "0.27–0.28。"]
    : ["3.30–7.16 J/document", "0.27–0.28 J/document"];
  for (const result of ["4.27 J", "0.27 J", ...energyRanges, "56.38 s", "0.18 s", "1,918.9 MiB", "42.6 MiB", "94%", "313×", "45×", "LibreOffice + Poppler"]) {
    assert.ok(html.includes(result), `${page}: benchmark includes ${result}`);
  }
}

// the documentation pages
for (const page of ["index.html", "index.ja.html", "formats/spreadsheet.html", "examples/search.ja.html", "api.html", "spec.html", "design.html"]) {
  const html = await readFile(join(site, "docs", page), "utf8");
  assert.doesNotMatch(html, /href="(?!https?:)[^"]*\.md(#[^"]*)?"/, `docs/${page} links to Markdown`);
}

// the graphs in why are relative to their pages and copied into the site
for (const page of ["why.html", "why.ja.html"]) {
  const html = await readFile(join(site, "docs", page), "utf8");
  const image = html.match(/<img src="([^"]*why-economy\.(?:en|ja)\.svg)"/);
  assert.ok(image, `docs/${page} includes its comparison graph`);
  const response = await fetch(new URL(image[1], new URL(page, `${base}docs/`)));
  assert.equal(response.status, 200, `docs/${page}'s graph is reachable`);
  assert.match(response.headers.get("content-type") ?? "", /^image\/svg\+xml(?:;|$)/, `docs/${page}'s graph is served as SVG`);
}

// errors carry the codes the page acts on
await assert.rejects(modules.office.convert(new Uint8Array([1, 2, 3]), {}), { code: "unknown-format" });
await assert.rejects(modules.pdf.open(new Uint8Array([1, 2, 3]), {}), { code: "unknown-format" });

server.close();
if (failed) {
  console.log(`${failed} sample(s) failed`);
  process.exit(1);
}
process.exit(0); // the Go programs wait for calls forever
