import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { BdfDocument, BufferSource } from "@bdf/core";
import { ResourceCache } from "../dist/resources.js";
import { DocumentSearch } from "../dist/search.js";

const root = new URL("../../../", import.meta.url);
const fixture = new Uint8Array(await readFile(new URL("testdata/demo.bdf", root)));

// Node has no fonts: faces that say their family, and a font set that notes what it is given.
globalThis.FontFace = class {
  constructor(family) { this.family = family; }
  load() { return Promise.resolve(this); }
};
function fontSet() {
  const added = new Set();
  return { added, add: (face) => added.add(face.family), delete: (face) => added.delete(face.family), check: () => true };
}

test("hits are measured once the fonts of their objects are loaded", async () => {
  const doc = await BdfDocument.open(new BufferSource(fixture));
  const fonts = fontSet();
  let decoded = 0;
  const res = new ResourceCache(doc, fonts, { decodeImage: async () => ({ width: 1, height: 1, close() {}, n: decoded++ }) });
  const measured = [];
  const measure = (font, text) => {
    const embedded = /"(bdf-[0-9a-f]{32})"/.exec(font)?.[1];
    if (embedded) measured.push(fonts.added.has(embedded));
    return text.length * 5;
  };
  const view = doc.view("slides");
  const search = new DocumentSearch(doc, measure, (h) => res.prepareText(h));
  const hits = await search.search(view, "shapes and paths");
  assert.equal(hits.length, 1);
  const rects = await search.locate(view, hits[0]);
  assert.equal(rects.length, 1);
  assert.ok(measured.length > 0 && measured.every(Boolean), `measured with the font loaded: ${measured}`);
  // the text needs no pictures
  assert.equal(decoded, 0);

  // without a renderer's fonts the objects alone are loaded, as before
  const plain = new DocumentSearch(await BdfDocument.open(new BufferSource(fixture)), (font, text) => text.length * 5);
  assert.deepEqual(await plain.locate(view, (await plain.search(view, "shapes and paths"))[0]), rects);
});
