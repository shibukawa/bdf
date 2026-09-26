import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { BufferSource, BdfDocument, TextSearch, normalizeQuery, extractText, Sep } from "../dist/index.js";

const root = new URL("../../../", import.meta.url);
const fixture = new Uint8Array(await readFile(new URL("testdata/demo.bdf", root)));

test("normalization folds width, case and kana", () => {
  assert.equal(normalizeQuery("ＡＢＣ ｶﾀｶﾅ Ⅸ"), "abc かたかな ix");
  assert.equal(normalizeQuery("ABC", true), "ABC");
});

test("text index matches the extractor's ordinals", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const view = doc.view("doc");
  assert.ok(view.textIndex);
  const runs = await doc.textIndex(view);
  const body = await doc.ensure(view.pages[0].layers[1].obj);
  const extracted = extractText(body, (h) => doc.objectSync(h));
  const bodyRuns = runs.filter((r) => r.a === 0 && r.b === 1);
  assert.equal(bodyRuns.length, extracted.length);
  bodyRuns.forEach((r, i) => {
    assert.equal(r.ordinal, extracted[i].ordinal);
    assert.equal(r.text, extracted[i].text);
    if (i > 0) assert.equal(r.sep, extracted[i].sep);
  });
  // MARK LINE => space, MARK PARAGRAPH => break, MARK CELL => break
  assert.equal(extracted[1].sep, Sep.BREAK, "first paragraph after heading");
  assert.equal(extracted[2].sep, Sep.SPACE, "line inside a paragraph");
  const cell = extracted.find((r) => r.text === "R2C2");
  assert.equal(cell.sep, Sep.BREAK);
});

test("search across line breaks, case-insensitively, with locations", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const view = doc.view("doc");
  const ts = new TextSearch(await doc.textIndex(view));
  // "list of" ends a line and "objects" starts the next one; the LINE mark makes it one phrase.
  const hits = ts.search("LIST OF OBJECTS");
  assert.equal(hits.length, 6); // 3 repeats x 2 pages
  assert.equal(hits[0].segments.length, 2);
  assert.equal(hits[0].text, "list of objects");
  assert.deepEqual(hits[0].segments.map((s) => [s.a, s.b]), [[0, 1], [0, 1]]);
  assert.equal(hits[0].segments[0].start, "BDF is a display-list format for browsers. Each page is a ".length);
  assert.equal(hits[0].segments[1].end, "objects".length);
  assert.match(hits[0].context, /each page is a list of objects/); // context is normalized text
  // paragraph breaks are never crossed
  assert.equal(ts.search("small. Repeated").length, 0);
  // header text is in layer 0 of both pages
  const header = ts.search("flow view");
  assert.deepEqual(header.map((h) => h.segments[0].a), [0, 1]);
  assert.equal(ts.search("R3C2").length, 2);
  assert.equal(ts.search("").length, 0);
  assert.equal(ts.search("zzz", { limit: 5 }).length, 0);
});

test("sheet index uses tile coordinates and cell marks", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const view = doc.view("sheet1");
  const ts = new TextSearch(await doc.textIndex(view));
  const hits = ts.search("Row 150");
  assert.equal(hits.length, 1);
  assert.deepEqual([hits[0].segments[0].a, hits[0].segments[0].b], [0, 1]);
  // numbers in adjacent cells never merge: "1001092" must not match although 1001 and 092... are neighbours
  assert.equal(ts.search("10921116").length, 0);
});

test("fallback index without a text index part", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const view = { ...doc.view("slides"), textIndex: undefined };
  const runs = await doc.textIndex(view);
  const ts = new TextSearch(runs);
  assert.equal(ts.search("BDF fixture deck").length, 3);
  assert.equal(ts.search("advance correction").length, 1);
});

test("MARK WRAP joins East Asian lines without a separator", async () => {
  const pptx = new Uint8Array(await readFile(new URL("../../../testdata/pptx/basic.bdf", import.meta.url)));
  const doc = await BdfDocument.open(new BufferSource(pptx));
  const view = doc.view("slides");
  const runs = await doc.textIndex(view);
  // the index and a walk of the objects agree
  const page = view.pages[3];
  const bodyLayer = page.layers.findIndex((l) => l.role === "body");
  const body = await doc.ensure(page.layers[bodyLayer].obj);
  const extracted = extractText(body, (h) => doc.objectSync(h));
  const indexed = runs.filter((r) => r.a === 3 && r.b === bodyLayer);
  assert.deepEqual(indexed.map((r) => [r.text, r.sep]).slice(1), extracted.map((r) => [r.text, r.sep]).slice(1));
  const wrapped = extracted.filter((r) => r.sep === Sep.NONE && /[\u3040-\u9fff]/.test(r.text));
  assert.ok(wrapped.length >= 2, "wrapped Japanese lines");
  const hits = new TextSearch(runs).search("改行します");
  assert.equal(hits.length, 1);
  assert.equal(hits[0].segments.length, 2);
});
