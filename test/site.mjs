// Smoke test of the demo site (examples/viewer/site.mjs): serves it, runs
// its converter modules in Node with Go's wasm_exec.js, and converts each
// sample with the site's fonts, fetched from the server as in a browser.
// PDFs are also converted a page at a time, as the viewer does, and the
// pages put into the outline with @bdf/core (npm run build first).
//
//   node test/site.mjs [site dir]
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { pathToFileURL } from "node:url";
import { BdfDocument, BufferSource, extractText } from "../packages/core/dist/index.js";
import { serve } from "./serve.mjs";

const site = resolve(process.argv[2] ?? "examples/viewer/.site");
await import(pathToFileURL(join(site, "wasm_exec.js")).href);
const { server, port } = await serve(site);
const base = `http://127.0.0.1:${port}/`;

/** Load a module; the Go program puts its API on globalThis and waits. */
async function load(file) {
  const go = new globalThis.Go();
  const { instance } = await WebAssembly.instantiate(await readFile(join(site, file)), go.importObject);
  go.run(instance);
  const conv = globalThis.bdfConverter;
  delete globalThis.bdfConverter;
  return conv;
}
const modules = { pdf: await load("bdf-pdf.wasm"), office: await load("bdf-office.wasm") };
assert.deepEqual(modules.pdf.formats.map((f) => f.name), ["pdf"]);
assert.deepEqual(modules.office.formats.map((f) => f.name), ["csv", "docx", "emf", "pptx", "visio", "xlsx"]);

let failed = 0;
const samples = JSON.parse(await readFile(join(site, "samples/index.json"), "utf8"));
for (const { name } of samples) {
  if (name.endsWith(".bdf")) continue;
  const data = new Uint8Array(await readFile(join(site, "samples", name)));
  const conv = name.endsWith(".pdf") ? modules.pdf : modules.office;
  const t0 = performance.now();
  try {
    const res = await conv.convert(data, { fonts: `${base}fonts/` });
    assert.deepEqual([...res.bdf.subarray(0, 4)], [0x62, 0x64, 0x66, 0], "bdf magic");
    // the Office converters found the site's fonts and embedded them
    if (res.format !== "pdf") assert.match(res.summary, /[1-9]\d* embedded font/);
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

// the documentation pages
for (const page of ["index.html", "ja.html", "api.html", "spec.html", "design.html"]) {
  const html = await readFile(join(site, "docs", page), "utf8");
  assert.doesNotMatch(html, /href="(?!https?:)[^"]*\.md(#[^"]*)?"/, `docs/${page} links to Markdown`);
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
