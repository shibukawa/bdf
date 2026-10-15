import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { BufferSource, BdfDocument, TextSearch, normalizeQuery, extractText, Sep } from "../dist/index.js";

const root = new URL("../../../", import.meta.url);
const fixture = new Uint8Array(await readFile(new URL("testdata/demo.bdf", root)));

test("normalization folds width, case and kana", () => {
  assert.equal(normalizeQuery("ＡＢＣ ｶﾀｶﾅ Ⅸ"), "abc かたかな ix");
  assert.equal(normalizeQuery("ABC", true), "ABC");
  // the minus sign of a formula matches a typed hyphen-minus
  assert.equal(normalizeQuery("b^2\u22124ac"), "b^2-4ac");
});

test("normalization is of the string, not of its characters one by one", () => {
  // a letter and the mark after it are the letter they compose
  assert.equal(normalizeQuery("cafe\u0301"), "caf\u00e9");
  assert.equal(normalizeQuery("CAFE\u0301", true), "CAF\u00c9");
  // half-width katakana with its voiced marks, and a kana with a combining one
  assert.equal(normalizeQuery("\uff76\uff9e\uff72\uff84\uff9e"), "\u304c\u3044\u3069");
  assert.equal(normalizeQuery("\uff8a\uff9f"), "\u3071");
  assert.equal(normalizeQuery("\u304b\u3099"), "\u304c");
  // the letters of a Hangul syllable
  assert.equal(normalizeQuery("\u1100\u1161\u11a8"), "\uac01");
  // marks in another order than the canonical one
  assert.equal(normalizeQuery("a\u0301\u0323"), normalizeQuery("a\u0323\u0301"));
  // a minus sign that normalization makes of another character
  assert.equal(normalizeQuery("10\u207b\u00b3"), "10-3");
});

test("a hit is located in the text it was normalized from", () => {
  const run = (text, sep = Sep.SPACE) => ({ a: 0, b: 0, ordinal: 0, sep, text });
  const at = (text, query) => new TextSearch([run(text)]).search(query).map((h) => [h.segments[0].start, h.segments[0].end, h.text]);
  // decomposed in the document, composed in the query, and the other way round
  assert.deepEqual(at("un cafe\u0301 noir", "caf\u00e9"), [[3, 8, "cafe\u0301"]]);
  assert.deepEqual(at("un caf\u00e9 noir", "CAFE\u0301"), [[3, 7, "caf\u00e9"]]);
  // the letter without its mark is another letter, as it is when they are one character
  assert.deepEqual(at("cafe\u0301", "cafe"), []);
  assert.deepEqual(at("caf\u00e9", "cafe"), []);
  // half-width katakana: a letter and its voiced mark are located together
  const guide = "\uff76\uff9e\uff72\uff84\uff9e"; // ｶﾞｲﾄﾞ
  assert.deepEqual(at(guide, "\u30ac\u30a4\u30c9"), [[0, 5, guide]]);
  assert.deepEqual(at(guide, "\u3044"), [[2, 3, "\uff72"]]);
  assert.deepEqual(at(guide, "\u30c9"), [[3, 5, "\uff84\uff9e"]]);
  assert.deepEqual(at(guide, "\u30c8"), []); // ト is not in ド
  // marks out of order, and a character outside the BMP before them
  assert.deepEqual(at("\ud835\udc9c a\u0301\u0323b", "a\u0323\u0301"), [[3, 6, "a\u0301\u0323"]]);
  assert.deepEqual(at("\ud835\udc9ce\u0301", "a\u00e9"), [[0, 4, "\ud835\udc9ce\u0301"]]);
  // a character that normalizes to several is located whole
  assert.deepEqual(at("\u3314", "\u308d"), [[0, 1, "\u3314"]]); // ㌔ is キロ
  // the letters of a syllable
  assert.deepEqual(at("x\u1100\u1161\u11a8y", "\uac01"), [[1, 4, "\u1100\u1161\u11a8"]]);
  // a letter with more marks than a piece holds: the first 31 are with it, the
  // rest by themselves, and the index is built in time in proportion to the text
  const marks = "\u0301".repeat(200000);
  const t0 = performance.now();
  const many = new TextSearch([run("e" + marks + "x")]);
  assert.ok(performance.now() - t0 < 5000, "200,000 marks after a letter");
  assert.deepEqual(many.search("\u00e9").map((h) => [h.segments[0].start, h.segments[0].end]), [[0, 32]]);
  assert.equal(many.search("x").length, 1);
  // a run is normalized by itself: a mark at the start of the next is not joined
  const split = new TextSearch([run("e"), run("\u0301x", Sep.NONE)]);
  assert.equal(split.search("\u00e9").length, 0);
  assert.equal(split.search("e\u0301x").length, 0);
  assert.equal(split.text, "e\u0301x");
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

test("Go text index and TS extractor agree on every page layer", async () => {
  // the demo uses every structure MARK kind; both extractors must derive the same separators
  const doc = await BdfDocument.open(new BufferSource(fixture));
  for (const view of doc.manifest.views.filter((v) => v.pages)) {
    const index = await doc.textIndex(view);
    for (const [a, page] of view.pages.entries()) {
      for (const [b, layer] of page.layers.entries()) {
        const extracted = extractText(await doc.ensure(layer.obj), (h) => doc.objectSync(h));
        const indexed = index.filter((r) => r.a === a && r.b === b);
        assert.equal(indexed.length, extracted.length, `${view.id} ${a}/${b}`);
        indexed.forEach((r, i) => {
          assert.equal(r.text, extracted[i].text);
          if (i > 0) assert.equal(r.sep, extracted[i].sep, `${view.id} ${a}/${b} run ${i} ${r.text}`);
        });
      }
    }
  }
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

test("hits take their context from the plain text, including line breaks", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const index = await doc.textIndex(doc.view("doc"));
  const ts = new TextSearch(index);
  const hits = ts.search("page", { context: 12, limit: 4 });
  assert.deepEqual(hits.map((h) => [h.text, h.context]), [
    ["page", "wsers. each page is a list o"], ["page", " from every page. the body o"],
    ["page", "ody of each page is its own "], ["page", "wsers. each page is a list o"],
  ]);
  const plain = ts.text;
  // line breaks of the text are shown in the context
  const shown = plain.replace(/\n/g, " ⏎ ");
  const all = ts.search("e", { context: 20 });
  assert.ok(all.length > 100);
  for (const h of all) assert.ok(shown.includes(h.context), h.context);
  assert.ok(all.some((h) => h.context.includes(" ⏎ ")));
  // whether the text was asked for before the search or not
  const other = new TextSearch(index);
  assert.equal(other.text, plain);
  assert.deepEqual(other.search("e", { context: 20 }), all);
});
