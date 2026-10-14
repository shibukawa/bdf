// Smoke test of the demo site (site/build.mjs): serves it, runs
// its converter modules in Node with Go's wasm_exec.js, and converts each
// sample with the site's fonts, fetched from the server as in a browser.
// PDFs are also converted a page at a time, as the viewer does, and the
// pages put into the outline with @bdfkit/core (npm run build first). The
// preview module draws the thumbnail of each converted sample, gives its
// text and writes it as a PDF file. The form module reads the tags of the
// templates of the form page (and the fields of a PDF form and of the
// templates of pdfme), fills them with their data, and writes the filled
// forms as PDF files: plain, tagged, and PDF/A.
//
//   node test/site.mjs [site dir]
import assert from "node:assert/strict";
import { readFile, readdir } from "node:fs/promises";
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
  preview: await load("bdf-preview.wasm"), form: await load("bdf-form.wasm"),
};
assert.deepEqual(modules.pdf.formats.map((f) => f.name), ["ai", "pdf"]);
assert.deepEqual(modules.office.formats.map((f) => f.name), ["audio", "cgm", "csv", "docx", "drawio", "dxf", "emf", "font", "gerber", "hpgl", "idml", "image", "jww", "kicad", "midi", "mml", "musicxml", "parquet", "pptx", "psd", "sxf", "visio", "xlsx"]);
assert.deepEqual(modules.image.formats.map((f) => f.name), ["audio", "image"]);
assert.deepEqual(modules.web.formats.map((f) => f.name), ["epub", "html", "markdown"]);
assert.deepEqual(modules.preview.formats, []);
// (the formats of the templates, and of the tables of data that package form reads)
// (and the PDF converter: a PDF with form fields, and the base PDF of a template of pdfme)
assert.deepEqual(modules.form.formats.map((f) => f.name), ["csv", "docx", "pdf", "pptx", "xlsx"]);
// what each module does besides converting
const calls = (m) => ["thumbnail", "text", "pdf", "fields", "fill"].filter((name) => typeof m[name] === "function");
assert.deepEqual(calls(modules.preview), ["thumbnail", "text", "pdf"]);
assert.deepEqual(calls(modules.form), ["pdf", "fields", "fill"]);
assert.deepEqual(calls(modules.office), []);
// the image module converts images and audio files (their covers and tags), neither needing fonts
const imageExtensions = modules.image.formats.flatMap((f) => f.extensions);
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
// an audio file's document carries the file, which the viewer plays, and cues for its timed lyrics; the options of a
// call leave it out (the format's own params, and the pages a thumbnail asks for)
try {
  const name = "synced.mp3";
  const data = new Uint8Array(await readFile(join(site, "samples", name)));
  const play = async (options) => {
    const res = await modules.image.convert(data, { name, ...options });
    const doc = await BdfDocument.open(new BufferSource(res.bdf));
    return { res, doc, view: doc.manifest.views[0] };
  };
  const whole = await play({});
  const audio = await whole.doc.audio(whole.view);
  assert.deepEqual(audio.data, data, "the audio is the file");
  assert.equal(audio.type, "audio/mpeg");
  assert.deepEqual(audio.cues.cues.map((c) => c.tick), [1500, 4500, 7500, 10500], "a cue a line of the lyrics, in milliseconds");
  assert.equal(whole.doc.manifest.parts.find((p) => p.t === "audio").enc, "identity");
  for (const options of [{ params: { play: "false" } }, { params: { play: false } }, { pages: "1" }]) {
    const silent = await play(options);
    assert.equal(silent.view.play, undefined, `${JSON.stringify(options)}: no audio`);
    assert.deepEqual(silent.res.warnings, []);
    assert.ok(silent.res.bdf.length < data.length, "a document without the audio is smaller than the file");
  }
  // a file over the limit converts to its card, with a warning
  const over = await play({ params: { maxaudio: 0.01 } });
  assert.equal(over.view.play, undefined);
  assert.match(over.res.warnings.join("\n"), /larger than .* its audio is not stored/);
  assert.ok((await play({ pages: "1", params: { play: "true" } })).view.play.audio, "play=true stores it for a selection of pages too");
  await assert.rejects(play({ params: { play: "perhaps" } }), /parameter play/);
  console.log(`ok   ${name}: the audio and the cues of its lyrics; left out by params and pages`);
} catch (e) {
  failed++;
  console.log(`FAIL synced.mp3 audio: ${e.message}`);
}

/** What a PDF file says of itself: its pages, and the fonts embedded in it. */
function pdfInfo(b) {
  const text = new TextDecoder("latin1").decode(b);
  assert.ok(text.startsWith("%PDF-1.") && text.trimEnd().endsWith("%%EOF"), "a PDF file");
  return { pages: Number(/\/Type \/Pages \/Count (\d+)/.exec(text)?.[1]), fonts: (text.match(/\/FontFile[23] /g) ?? []).length };
}

// the form page: each shape of form, as a template of each format, filled with the data it comes with
const formFonts = `${base}form/fonts/`;
const formFontDir = join(site, "form/fonts");
assert.deepEqual(await readdir(formFontDir), ["index.json"], "form fonts are fetched remotely, with no bundled font files");
const remoteFonts = JSON.parse(await readFile(join(formFontDir, "index.json"), "utf8"));
assert.equal(remoteFonts.length, 2);
for (const font of remoteFonts) {
  const url = new URL(font.url);
  assert.equal(url.origin, "https://fonts.gstatic.com");
  assert.equal(url.search, "", "the font URL contains no document text");
  assert.ok(font.size > 0 && font.scan.length > 0);
}
const forms = JSON.parse(await readFile(join(site, "forms/index.json"), "utf8"));
assert.deepEqual(forms.map((f) => f.id), ["fixed", "rows", "pages", "application", "pdfme-invoice", "pdfme-certificate"]);
/** The format of a template of the page, by its name. */
const formatOf = (name) => (name.startsWith("pdfme-") ? "pdfme" : name.split(".").pop());
const styledSchemas = new Map();
function withoutLocations(value) {
  if (Array.isArray(value)) return value.map(withoutLocations);
  if (value && typeof value === "object") return Object.fromEntries(Object.entries(value)
    .filter(([key]) => key !== "at").map(([key, v]) => [key, withoutLocations(v)]));
  return value;
}
for (const shape of forms) {
  const data = await readFile(join(site, "forms", shape.data), "utf8");
  const files = Object.fromEntries(await Promise.all(shape.files.map(async (f) => [f, new Uint8Array(await readFile(join(site, "forms", f)))])));
  for (const { name, style } of shape.templates) {
    const t0 = performance.now();
    try {
      const template = new Uint8Array(await readFile(join(site, "forms", name)));
      // the tags: the names, and the JSON Schema that the form is made of, with the sample and without
      const asked = await modules.form.fields(template, { sample: data });
      const fields = JSON.parse(asked.fields), schema = JSON.parse(asked.schema);
      assert.ok(fields.length > 0 && fields.every((f) => f.name && f.kind && f.at.length > 0), "fields");
      assert.equal(schema.$schema, "https://json-schema.org/draft/2020-12/schema");
      assert.equal(schema.type, "object");
      if (style) {
        const key = `${shape.id}/${formatOf(name)}`;
        const normalized = withoutLocations(schema);
        if (styledSchemas.has(key)) assert.deepEqual(normalized, styledSchemas.get(key), "both styles ask for the same data");
        else styledSchemas.set(key, normalized);
      }
      const sample = JSON.parse(data);
      const props = schema.properties;
      // every name of the sample is a property, of the type of its value
      for (const [key, value] of Object.entries(sample)) {
        assert.ok(props[key], `the schema has ${key}`);
        const type = Array.isArray(value) ? "array" : typeof value;
        // (a field of a PDF form or of pdfme says what it takes; a tag takes what the sample has)
        if (props[key].type !== "string" || type !== "number") assert.equal(props[key].type, type, `the type of ${key}`);
      }
      assert.equal(asked.format, formatOf(name));
      const own = await modules.form.fields(template, {});
      const bare = JSON.parse(own.schema);
      // a template of pdfme has a sample of its own: what its designer left in its fields
      assert.equal(typeof own.sample, name.startsWith("pdfme-") ? "string" : "undefined");
      if (shape.id === "application") {
        assert.deepEqual(props.plan.enum, ["basic", "pro", "team"]);
        assert.deepEqual(props.country["x-bdf"].labels, ["Japan", "United States", "Germany"]);
        assert.equal(props.zip.maxLength, 7);
        assert.equal(props.agree.type, "boolean");
        assert.equal(props.applicant.properties.name.type, "string");
      }
      if (shape.id === "pdfme-invoice") {
        assert.deepEqual(Object.keys(props.items.items.properties), ["項目", "数量", "単価", "金額"]);
        assert.equal(props.items["x-bdf"].mode, "rows");
        assert.equal(props.due.format, "date");
        assert.deepEqual(Object.keys(props.totals.properties).sort(), ["subtotal", "tax", "total"]);
        assert.equal(JSON.parse(own.sample).number, "INV-2026-1014");
        // filled with its own sample, it is the page its designer saw
        const seen = await modules.form.fill(template, { data: own.sample, fonts: formFonts, strict: true });
        assert.equal(seen.pages, 1);
        assert.deepEqual(seen.warnings, []);
      }
      if (shape.id === "pdfme-certificate") {
        assert.equal(props.signature["x-bdf"].kind, "image");
        assert.equal(props.honours.type, "boolean");
        assert.equal(props.points.type, "number");
      }
      if (shape.id === "fixed") {
        assert.equal(props.jan["x-bdf"].code, "ean13");
        assert.equal(props.url["x-bdf"].code, "qr");
        assert.equal(props.photo["x-bdf"].kind, "image");
        assert.equal(props.date.format, "date");
        assert.equal(props.customer.properties.name.type, "string");
        assert.equal(props.tags.items.type, "string");
        assert.deepEqual(bare.properties.no.type, ["string", "number"], "without a sample a value is a text or a number");
      }
      if (shape.id === "rows") assert.deepEqual(Object.keys(props.items.items.properties), ["name", "qty", "price", "amount"]);
      if (shape.id === "pages") {
        assert.equal(props.pages.items.properties.rows.maxItems, 10);
        assert.equal(props.pages.items.properties.rows["x-bdf"].mode, "slots");
        // (the sample has them around the pages; the tags alone leave them in a page)
        assert.equal(props.slip.type, "string");
        assert.ok(bare.properties.pages.items.properties.slip && !bare.properties.slip);
        assert.ok(props.pages.items.properties.page_total);
      }
      // filled: a document whose text has the values, and no warnings
      const res = await modules.form.fill(template, { data, files, fonts: formFonts, name });
      assert.deepEqual([...res.bdf.subarray(0, 4)], [0x62, 0x64, 0x66, 0], "bdf magic");
      assert.equal(res.format, formatOf(name));
      assert.deepEqual(res.warnings, []);
      const doc = await BdfDocument.open(new BufferSource(res.bdf));
      assert.equal(doc.manifest.views.reduce((n, v) => n + (v.pages?.length ?? 0), 0), res.pages);
      const text = (await Promise.all(doc.manifest.views.map(async (v) => {
        let out = "";
        for (const page of v.pages ?? []) for (const l of page.layers) for (const r of extractText(await doc.ensure(l.obj), (h) => doc.objectSync(h))) out += r.text;
        return out;
      }))).join("");
      const want = {
        fixed: ["明細書", "Date", "株式会社みなと商事", "デスクオーガナイザー", "1024", "13,200"],
        rows: ["請求書", "INV-2026-1014", "付箋セット", "23,958"],
        pages: ["納品書", "DN-2026-1014", "当頁計", "36,180"], application: ["Taro Yamada", "12800"],
        "pdfme-invoice": ["請求書", "INV-2026-1014", "ファイルボックス", "総額", "23,958"],
        "pdfme-certificate": ["Hanako Tanaka", "Score 92 / 100 (with distinction)", "October 14, 2026"],
      }[shape.id];
      for (const w of want) assert.ok(text.includes(w), `the filled form has ${w}`);
      assert.ok(!text.includes("{{"), "no tag is left in the form");
      // as a PDF file, with its text in embedded fonts
      const pdf = await modules.form.pdf(res.bdf, { fonts: formFonts });
      const info = pdfInfo(pdf.pdf);
      assert.equal(info.pages, res.pages);
      assert.equal(pdf.pages, res.pages);
      assert.ok(pdf.fonts > 0 && info.fonts === pdf.fonts, "embedded fonts");
      assert.deepEqual(pdf.warnings, []);
      // tagged, and as a PDF/A file that is tagged
      const plain = new TextDecoder("latin1").decode(pdf.pdf);
      assert.ok(!plain.includes("/StructTreeRoot") && !plain.includes("pdfaid:part"), "a plain PDF file");
      const tagged = new TextDecoder("latin1").decode((await modules.form.pdf(res.bdf, { fonts: formFonts, tagged: true })).pdf);
      assert.ok(tagged.includes("/StructTreeRoot") && tagged.includes("/MarkInfo << /Marked true >>") && !tagged.includes("pdfaid:part"), "a tagged PDF file");
      const kept = await modules.form.pdf(res.bdf, { fonts: formFonts, pdfa: "2a" });
      const pdfa = new TextDecoder("latin1").decode(kept.pdf);
      assert.ok(pdfa.includes("/StructTreeRoot") && pdfa.includes("<pdfaid:part>2</pdfaid:part>") && pdfa.includes("<pdfaid:conformance>A</pdfaid:conformance>") && pdfa.includes("/OutputIntents"), "a PDF/A-2a file");
      assert.equal(pdfInfo(kept.pdf).pages, res.pages);
      await assert.rejects(modules.form.pdf(res.bdf, { pdfa: "9z" }), /pdfa|PDF\/A/i);
      console.log(`ok   forms/${name}: ${fields.length} fields, filled into ${res.pages} page(s), a PDF of ${pdf.pdf.length} bytes with ${pdf.fonts} font(s), ${(performance.now() - t0).toFixed(0)} ms`);
    } catch (e) {
      failed++;
      console.log(`FAIL forms/${name}: ${e.message}`);
    }
  }
}
// An expanded Excel invoice repeats only its table heading on continuation
// pages. Totals, payment details and the note remain whole after the table.
for (const style of ["", "-modern"]) {
  const name = `form-rows${style}.xlsx`;
  const template = new Uint8Array(await readFile(join(site, "forms", name)));
  const sample = JSON.parse(await readFile(join(site, "forms/form-rows.json"), "utf8"));
  try {
    for (const count of [0, 5, 12, 13, 14, 16, 25, 35, 60, 100]) {
      const data = { ...sample, customer: "株式会社みなと商事営業部",
        note: "改ページ確認用の備考。\n内容が分断されないことを確認します。",
        items: Array.from({ length: count }, (_, i) => ({
          name: `明細項目[${String(i + 1).padStart(3, "0")}]`, qty: 1, price: 100, amount: 100,
        })), subtotal: count * 100, tax: count * 10, total: count * 110 };
      const res = await modules.form.fill(template, { data: JSON.stringify(data), fonts: formFonts, strict: true });
      assert.deepEqual(res.warnings, []);
      const pages = await pageTexts(await BdfDocument.open(new BufferSource(res.bdf)));
      const all = pages.join("\n");
      assert.equal(all.split(data.customer).length - 1, 1, `${name}/${count}: customer appears once`);
      for (const item of data.items) assert.equal(all.split(item.name).length - 1, 1, `${name}/${count}: ${item.name} appears once`);
      for (const [i, page] of pages.entries()) {
        if (i > 0) assert.equal(page.includes("単価"), page.includes("明細項目["), `${name}/${count}/page ${i + 1}: headings only with continued details`);
        for (const group of [
          ["小計", "消費税", "総額"],
          ["お振込先", "みらい銀行", "口座名義", "お支払期限までに上記口座へお振り込みください。"],
          ["備考", ...data.note.split("\n"), "いつもお引き立てを賜り、誠にありがとうございます。"],
        ]) {
          const seen = group.filter((text) => page.includes(text));
          assert.ok(seen.length === 0 || seen.length === group.length, `${name}/${count}/page ${i + 1}: split section ${group.join(" / ")}`);
        }
      }
    }
    console.log(`ok   forms/${name}: 0–100 detail rows, scoped table headings, single customer, intact totals/payment/notes`);
  } catch (e) {
    failed++;
    console.log(`FAIL forms/${name} pagination: ${e.message}`);
  }
}

// data that lacks a name leaves its tag empty, with a warning; strictly it is an error; what is no template is refused
{
  const template = new Uint8Array(await readFile(join(site, "forms/form-rows.xlsx")));
  const res = await modules.form.fill(template, { data: JSON.stringify({ no: 1, items: [] }), fonts: formFonts });
  assert.match(res.warnings.join("\n"), /no "customer"/);
  await assert.rejects(modules.form.fill(template, { data: "{}", strict: true, fonts: formFonts }), /the data has no/);
  await assert.rejects(modules.form.fill(template, { data: "[1]" }), /object of JSON/);
  await assert.rejects(modules.form.fill(template, { data: "{", fonts: formFonts }), /data:/);
  await assert.rejects(modules.form.fill(template, { data: "{}", files: { "a.png": "not bytes" } }), /Uint8Array/);
  // (a PDF without form fields, a picture, JSON that is no template of pdfme)
  await assert.rejects(modules.form.fields(new Uint8Array(await readFile(join(site, "samples/demo.pdf"))), {}), /no form fields/);
  await assert.rejects(modules.form.fields(new Uint8Array(await readFile(join(site, "forms/brand-mark.png"))), {}), /not an Excel workbook/);
  await assert.rejects(modules.form.fields(new TextEncoder().encode('{"items": []}'), {}), /not an Excel workbook/);
  await assert.rejects(modules.form.fields(new TextEncoder().encode('{"schemas": [[]]}'), {}), /no base PDF/);
  // a compiled template holds its schema and its sample
  const bdft = new Uint8Array(await readFile(new URL("../testdata/form/pages.bdft", import.meta.url)));
  const asked = await modules.form.fields(bdft, {});
  assert.equal(asked.format, "xlsx");
  assert.equal(JSON.parse(asked.sample).slip, "D-1");
  assert.equal(JSON.parse(asked.schema).properties.slip.type, "string");
  const filled = await modules.form.fill(bdft, { data: asked.sample, fonts: `${base}fonts/` });
  assert.equal(filled.pages, 3);
  console.log("ok   forms: missing names, strict filling, what is refused, a compiled template");
}

// the viewer's PDF files: of a document converted on the page, of one whose fonts are WOFF2 (made by the
// bdf command: the preview module reads them), and of an encrypted one with its password
try {
  for (const name of ["features.pptx", "basic.docx", "reportlab-master.pdf", "aws.drawio"]) {
    const { bdf } = converted.find((c) => c.name === name);
    const doc = await BdfDocument.open(new BufferSource(bdf));
    const res = await modules.preview.pdf(bdf, { fonts: `${base}fonts/` });
    const info = pdfInfo(res.pdf);
    const paged = doc.manifest.views.filter((v) => v.kind === "fixed" || v.kind === "flow");
    assert.equal(info.pages, res.pages);
    assert.equal(res.pages, paged.reduce((n, v) => n + v.pages.length, 0), `${name}: the pages of its views`);
    assert.ok(res.fonts > 0 && info.fonts === res.fonts, `${name}: embedded fonts`);
    assert.ok(!res.warnings.some((w) => /font/.test(w)), `${name}: ${res.warnings.join("; ")}`);
  }
  const plain = new Uint8Array(await readFile(new URL("../testdata/demo.bdf", import.meta.url)));
  const enc = new Uint8Array(await readFile(new URL("../testdata/demo-encrypted.bdf", import.meta.url)));
  const demo = await modules.preview.pdf(plain, { fonts: `${base}fonts/` });
  assert.ok(demo.fonts > 0 && pdfInfo(demo.pdf).fonts === demo.fonts, "the WOFF2 fonts of a document are embedded");
  assert.deepEqual(demo.warnings, []);
  // (without the font directory, the text of the fonts it refers to by name is left out, and said so)
  assert.match((await modules.preview.pdf(plain, {})).warnings.join("\n"), /no font for "serif"/);
  // one view of it, and without a layer
  const first = (await BdfDocument.open(new BufferSource(plain))).manifest.views[0];
  const one = await modules.preview.pdf(plain, { views: [first.id], skip: ["background"], fonts: `${base}fonts/` });
  assert.equal(one.pages, first.pages.length);
  await assert.rejects(modules.preview.pdf(enc, {}), { code: "password-required" });
  await assert.rejects(modules.preview.pdf(enc, { password: "wrong" }), { code: "wrong-password" });
  assert.deepEqual((await modules.preview.pdf(enc, { password: "demo-パスワード", fonts: `${base}fonts/` })).pdf, demo.pdf, "the encrypted document's PDF file");
  await assert.rejects(modules.preview.pdf(new Uint8Array([1, 2, 3]), {}));
  console.log(`ok   PDF files of converted samples, of testdata/demo.bdf (${demo.pages} pages, ${demo.fonts} fonts) and of the encrypted one`);
} catch (e) {
  failed++;
  console.log(`FAIL PDF files: ${e.message}`);
}

/** The width and height of a PNG image. */
function pngSize(b) {
  assert.deepEqual([...b.subarray(0, 8)], [0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a], "PNG signature");
  const dv = new DataView(b.buffer, b.byteOffset, b.byteLength);
  return [dv.getUint32(16), dv.getUint32(20)];
}
/** The layout the thumbnail package picks for the formats that have one. */
const LAYOUT = { docx: "crop", html: "crop", markdown: "crop", xlsx: "crop", csv: "crop", parquet: "crop", font: "crop", mp3: "crop", pptx: "fit", vsdx: "fit", drawio: "fit", epub: "fit", svg: "fit", jpeg: "fit" };
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

// the home pages show the measured comparison below the viewer in both
// languages, with the numbers of docs/benchmarks/latest.json as the page
// prints them (two decimals of a second, a tenth of a MiB, hundredths of a
// joule, rounded ratios), so a remeasurement changes the test's input, not
// the test
const benchmark = JSON.parse(await readFile(new URL("../docs/benchmarks/latest.json", import.meta.url), "utf8"));
function benchmarkNumbers(lang) {
  const number = (value, digits) => new Intl.NumberFormat(lang, { minimumFractionDigits: digits, maximumFractionDigits: digits }).format(value);
  const results = Object.values(benchmark.process.results);
  const energy = benchmark.power.median_estimated_soc_j_per_document;
  const ranges = benchmark.power.range_estimated_soc_j_per_document;
  const range = (w) => `${number(ranges[w][0], 2)}–${number(ranges[w][1], 2)}`;
  return [
    ...results.flatMap((r) => ["libreoffice", "bdf"].flatMap((w) => [`${number(r[w].elapsed_seconds, 2)} s`, `${number(r[w].max_rss_bytes / 2 ** 20, 1)} MiB`])),
    `${number(energy.libreoffice, 2)} J`, `${number(energy.bdf, 2)} J`,
    // the highlights: energy saved, and the largest speed and memory ratios
    `${Math.round(Math.abs(1 - energy.bdf / energy.libreoffice) * 100)}%`,
    `${number(Math.max(...results.map((r) => r.libreoffice.elapsed_seconds / r.bdf.elapsed_seconds)), 0)}×`,
    `${number(Math.max(...results.map((r) => r.libreoffice.max_rss_bytes / r.bdf.max_rss_bytes)), 0)}×`,
    lang === "ja" ? `${range("libreoffice")}、BDF ${range("bdf")}。` : `${range("libreoffice")} J/document; BDF ${range("bdf")} J/document.`,
    "LibreOffice + Poppler",
  ];
}
// every page of the site links to the form page, which loads its script
for (const page of ["index.html", "index.ja.html", "viewer/index.html", "thumbnail/index.html", "text/index.html", "form/index.html", "docs/index.html"]) {
  const html = await readFile(join(site, page), "utf8");
  assert.match(html, /<a href="(\.\.\/|\.\/)?form\/"/, `${page} links to the form page`);
}
{
  const html = await readFile(join(site, "form/index.html"), "utf8");
  assert.match(html, /<a href="\.\.\/form\/" aria-current="page">Forms<\/a>/);
  assert.match(html, /<script type="module" src="\.\/main\.js">/);
  assert.equal((await fetch(`${base}form/main.js`)).status, 200);
  assert.match(await readFile(join(site, "viewer/index.html"), "utf8"), /<button id="savePdf"/);
}

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
  for (const result of benchmarkNumbers(page.endsWith(".ja.html") ? "ja" : "en")) {
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
