import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";
import { BufferSource, BdfDocument, parseHeader, walk, NoopSink, opHistogram, extractText, extractContent, parseCellRef, objectDeps, SplitSource, dcValues } from "../dist/index.js";

const root = new URL("../../../", import.meta.url);
const fixture = new Uint8Array(await readFile(new URL("testdata/demo.bdf", root)));

test("header parses", () => {
  assert.equal(new TextDecoder().decode(fixture.subarray(0, 4)), "bdf\0");
  const h = parseHeader(fixture);
  assert.equal(h.version, 1);
  assert.equal(h.manifestOff, 32);
  assert.ok(h.manifestLen > 0);
  const old = fixture.slice(0, 32);
  old.set([0x42, 0x44, 0x46, 0x31]); // "BDF1"
  assert.throws(() => parseHeader(old), /bad magic/);
});

test("Dublin Core metadata", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const dc = doc.manifest.meta.dc;
  assert.equal(dc.title, "BDF fixture");
  assert.deepEqual(dcValues(dc.title), ["BDF fixture"]);
  assert.deepEqual(dcValues(dc.creator), ["bdf-go", "BDF fixture generator"]);
  assert.deepEqual(dcValues(dc.language), ["en", "ja"]);
  assert.deepEqual(dcValues(dc.rights), []);
  assert.equal(doc.manifest.meta.source, "fixture");
});

test("manifest and views", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const m = doc.manifest;
  assert.equal(m.bdf, 1);
  assert.equal(m.opset, 1);
  assert.deepEqual(m.views.map((v) => v.kind), ["fixed", "flow", "sheet"]);
  assert.equal(doc.view("slides").pages.length, 3);
  assert.equal(doc.view("doc").pages[0].body.x, 72);
  assert.equal(Object.keys(doc.view("sheet1").tiles).length, 2);
  // every layer and tile references a known object part
  for (const v of m.views) {
    for (const p of v.pages ?? []) for (const l of p.layers) assert.equal(doc.entry(l.obj).t, "obj");
    for (const h of Object.values(v.tiles ?? {})) assert.equal(doc.entry(h).t, "obj");
  }
});

test("all objects decode and walk (deflate-raw via DecompressionStream)", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  let compressed = 0;
  for (const e of doc.manifest.parts) {
    if (e.t !== "obj") continue;
    if (e.enc === "deflate-raw") compressed++;
    const bytes = await doc.part(e.h);
    assert.equal(bytes.length, e.size);
    const o = await doc.object(e.h);
    assert.equal(o.opset, 1);
    walk(o, new NoopSink());
    for (const dep of objectDeps(o)) doc.entry(dep);
  }
  assert.ok(compressed >= 4, `expected compressed objects, got ${compressed}`);
});

test("master is shared and USEs the icon", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const slides = doc.view("slides");
  const masters = new Set(slides.pages.map((p) => p.layers[0].obj));
  assert.equal(masters.size, 1);
  const master = await doc.ensure([...masters][0]);
  const hist = opHistogram(master);
  assert.equal(hist.useAt, 3);
  assert.equal(hist.fillPaint, 1);
  assert.equal(master.objects.length, 1);
  const icon = doc.objectSync(master.objects[0]);
  assert.equal(opHistogram(icon).fillPath, 1);
});

test("text extraction follows transforms and USE", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const page = doc.view("slides").pages[0];
  const resources = [];
  const body = await doc.ensure(page.layers[1].obj, (e) => { resources.push(e.t); });
  assert.ok(resources.includes("font"));
  const runs = extractText(body, (h) => doc.objectSync(h), [2, 0, 0, 2, 10, 10]);
  const title = runs.find((r) => r.text === "Shapes and paths");
  assert.ok(title);
  assert.equal(title.x, 2 * 48 + 10);
  assert.equal(title.y, 2 * 80 + 10);
  assert.equal(title.size, 36);
  assert.equal(title.font.weight, 700);
  assert.ok(title.advance > 200);
  // master text is reachable through the master layer
  const master = await doc.ensure(page.layers[0].obj);
  const mruns = extractText(master, (h) => doc.objectSync(h));
  assert.equal(mruns[0].text, "BDF fixture deck");
});

test("cell references", () => {
  assert.deepEqual(parseCellRef("B12"), { row: 11, col: 1, rows: 1, cols: 1 });
  assert.deepEqual(parseCellRef("A4:B5 col"), { row: 3, col: 0, rows: 2, cols: 2, scope: "col" });
  assert.deepEqual(parseCellRef("AA1 row"), { row: 0, col: 26, rows: 1, cols: 1, scope: "row" });
  assert.equal(parseCellRef("R1C1"), undefined);
  assert.equal(parseCellRef("A0"), undefined);
});

test("content: headings, lists, tables, figures and links", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const content = async (view, page, layer) => {
    const o = await doc.ensure(doc.view(view).pages[page].layers[layer].obj);
    return extractContent(o, (h) => doc.objectSync(h));
  };
  const nodeOf = (c, text) => c.nodes[c.runs.find((r) => r.text.startsWith(text)).node];
  // flow page: heading, then a table whose first row holds column headers
  const flow = await content("doc", 0, 1);
  assert.deepEqual(nodeOf(flow, "Section 1"), { kind: "heading", level: 1, parent: -1 });
  const cell = flow.nodes[nodeOf(flow, "R1C2").parent];
  assert.deepEqual({ kind: cell.kind, row: cell.row, col: cell.col, scope: cell.scope }, { kind: "cell", row: 0, col: 1, scope: "col" });
  const r3 = flow.nodes[nodeOf(flow, "R3C1").parent];
  assert.equal(r3.scope, undefined);
  assert.equal(flow.nodes[r3.parent].kind, "table");
  // slide 3: a list of four items, and links moved into page space
  const slide = await content("slides", 2, 1);
  const items = slide.runs.filter((r) => r.text.startsWith("Regular")).map((r) => slide.nodes[slide.nodes[r.node].parent]);
  assert.deepEqual(items.map((n) => n.kind), ["item", "item", "item", "item"]);
  assert.equal(new Set(items.map((n) => n.parent)).size, 1);
  assert.equal(slide.nodes[items[0].parent].kind, "list");
  assert.deepEqual(slide.links.map((l) => l.url), ["https://example.com/", "#page=1"]);
  // slide 2: figures bounded by what they draw, clipped
  const figs = (await content("slides", 1, 1)).nodes.filter((n) => n.kind === "figure");
  assert.deepEqual(figs.map((f) => [f.alt, f.bounds]), [
    ["Checkerboard", { x: 60, y: 130, w: 200, h: 200 }],
    ["Enlarged corner of the checkerboard", { x: 780, y: 130, w: 120, h: 200 }],
  ]);
  // footer: a figure drawn by a USE_AT child
  const logo = (await content("doc", 0, 2)).nodes.find((n) => n.kind === "figure");
  assert.equal(logo.alt, "BDF logo");
  assert.ok(logo.bounds.w > 20 && logo.bounds.w < 25 && logo.bounds.y > 780);
});

test("converted PDF: structure from the tagged tree, and links over their text", async () => {
  // Chrome's content leaves a cm in effect; the LINK must still land on its text
  const doc = await BdfDocument.open(new BufferSource(new Uint8Array(await readFile(new URL("testdata/pdf/chrome-doc.bdf", root)))));
  assert.equal(dcValues(doc.manifest.meta.dc.language)[0], "en-US");
  const c = extractContent(await doc.ensure(doc.view("pages").pages[0].layers[0].obj), (h) => doc.objectSync(h));
  const kinds = new Set(c.nodes.map((n) => n.kind));
  for (const k of ["heading", "paragraph", "list", "item", "table", "cell"]) assert.ok(kinds.has(k), k);
  assert.equal(c.links.length, 1);
  const [l] = c.links, run = c.runs.find((r) => r.text === "hyperlink");
  const cx = run.x + run.matrix[0] * run.advance / 2, cy = run.y - run.matrix[3] * run.size * 0.3;
  assert.ok(cx > l.x && cx < l.x + l.w && cy > l.y && cy < l.y + l.h, JSON.stringify({ l, cx, cy }));
});

test("converted PDF: soft masks are drawn with MASK_BEGIN/MASK_END and are not text", async () => {
  const doc = await BdfDocument.open(new BufferSource(new Uint8Array(await readFile(new URL("testdata/pdf/chrome-masks.bdf", root)))));
  const page = await doc.ensure(doc.view("pages").pages[0].layers[0].obj);
  // Chrome puts each masked element in a form: count the MASK ops of every object.
  const counts = {}, seen = new Set();
  const visit = (o) => {
    for (const [k, v] of Object.entries(opHistogram(o))) counts[k] = (counts[k] ?? 0) + v;
    for (const h of o.objects) if (!seen.has(h)) { seen.add(h); visit(doc.objectSync(h)); }
  };
  visit(page);
  assert.ok(counts.maskBegin >= 5 && counts.maskBegin === counts.maskEnd, JSON.stringify(counts));
  const text = extractText(page, (h) => doc.objectSync(h)).map((r) => r.text).join(" ");
  assert.match(text, /Soft masks fade this heading out/);
  assert.match(text, /luminance: white to black/);

  // A MASK_BEGIN … MASK_END section keeps its text out and its transform in.
  let masked = 0;
  class Sink extends NoopSink {
    maskBegin(kind, backdrop, transfer) { masked++; assert.ok(kind === 0 || kind === 1); assert.ok(transfer.length === 0 || transfer.length === 256); }
  }
  for (const h of seen) walk(doc.objectSync(h), new Sink());
  assert.equal(masked, counts.maskBegin);
});

test("split source with a fake fetch", async () => {
  const base = new URL("testdata/demo-split/", root);
  const realFetch = globalThis.fetch;
  globalThis.fetch = async (url) => {
    const path = fileURLToPath(new URL(url));
    const body = await readFile(path);
    return new Response(body, { status: 200 });
  };
  try {
    const doc = await BdfDocument.open(new SplitSource(base.href));
    assert.equal(doc.manifest.views.length, 3);
    const tile = await doc.ensure(doc.view("sheet1").tiles["0,1"]);
    assert.ok(opHistogram(tile).fillText > 1000);
    assert.equal(tile.objects.length, 1, "second tile references the chart object");
  } finally {
    globalThis.fetch = realFetch;
  }
});
