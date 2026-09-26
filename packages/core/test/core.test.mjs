import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";
import { BufferSource, BdfDocument, parseHeader, walk, NoopSink, opHistogram, extractText, objectDeps, SplitSource, dcValues } from "../dist/index.js";

const root = new URL("../../../", import.meta.url);
const fixture = new Uint8Array(await readFile(new URL("fixtures/demo.bdf", root)));

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

test("split source with a fake fetch", async () => {
  const base = new URL("fixtures/demo-split/", root);
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
