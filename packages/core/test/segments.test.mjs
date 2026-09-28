import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { BufferSource, BdfDocument, BdfKeyError, BdfPasswordError, BdfSegmentError, SegmentLoader } from "../dist/index.js";

// testdata/segments/basic-{1,4,7}.bdf are the segments of 3 pages of
// testdata/epub/basic.bdf that hold pages 1, 4 and 7, each written by the Go
// writer (npm run testdata: bdf segment -size 3 -have …) as a server sends
// them to a reader that holds the segments before it, sealed for the key
// pair below. It was made for these tests only: its public key is
// testdata/segments/reader.pem.
const READER = {"kty":"EC","crv":"P-256","x":"TL6kYiUfNKSRKbTiahzyxwcAyOUy0pTL2c6JfeoqL1E","y":"OK8MOP4eHuQHLjAnrz7hubCKED1E8wPUbZYJQY-3z5I","d":"_loj2CexGXt2g1YkkUeTui4c9IeJ-BtcOYRzVL81Zi4"};
const root = new URL("../../../", import.meta.url);
const read = async (path) => new Uint8Array(await readFile(new URL(path, root)));
const plainBytes = await read("testdata/epub/basic.bdf");
const segmentBytes = { 0: await read("testdata/segments/basic-1.bdf"), 3: await read("testdata/segments/basic-4.bdf"), 6: await read("testdata/segments/basic-7.bdf") };

const { x, y } = READER;
const keyPair = {
  privateKey: await crypto.subtle.importKey("jwk", READER, { name: "ECDH", namedCurve: "P-256" }, false, ["deriveKey"]),
  publicKey: await crypto.subtle.importKey("jwk", { kty: "EC", crv: "P-256", x, y }, { name: "ECDH", namedCurve: "P-256" }, true, []),
};
const openSegment = (from) => BdfDocument.open(new BufferSource(segmentBytes[from]), { keyPair });

test("a segment opens with the key pair it was sealed for only", async () => {
  const source = new BufferSource(segmentBytes[0]);
  await assert.rejects(BdfDocument.open(source), BdfKeyError, "no key: a password would not open it either");
  await assert.rejects(BdfDocument.open(source, { password: "x" }), (e) => e instanceof BdfPasswordError && e.reason === "wrong");
  const other = await crypto.subtle.generateKey({ name: "ECDH", namedCurve: "P-256" }, false, ["deriveKey"]);
  await assert.rejects(BdfDocument.open(source, { keyPair: other }), BdfKeyError);
  const doc = await openSegment(0);
  assert.deepEqual(doc.manifest.segment, { view: "pages", from: 0, to: 3 });
  assert.ok(!new TextDecoder("latin1").decode(segmentBytes[0]).includes("A Small Book"), "the title is in the clear");
});

test("segments put together are the document's pages", async () => {
  const plain = await BdfDocument.open(new BufferSource(plainBytes));
  const doc = await openSegment(0);
  const pages = doc.view("pages").pages;
  assert.equal(pages.length, 7);
  assert.deepEqual(pages.map((p) => p.layers.length > 0), [true, true, true, false, false, false, false]);
  assert.equal(doc.view("pages").textIndex, undefined, "a segment carries no text index");
  // a later segment leaves out what an earlier one brought
  const later = await openSegment(3);
  const earlier = new Set(doc.manifest.parts.map((e) => e.h));
  assert.ok(later.manifest.parts.every((e) => !earlier.has(e.h)));
  doc.addSegment(later);
  doc.addSegment(await openSegment(6));
  const want = plain.view("pages").pages;
  assert.deepEqual(pages, want);
  for (const p of want) {
    for (const l of p.layers) {
      await doc.ensure(l.obj);
      assert.deepEqual(await doc.part(l.obj), await plain.part(l.obj));
    }
  }
  for (const e of doc.manifest.parts) assert.deepEqual(await doc.part(e.h), await plain.part(e.h), e.h);
});

test("a segment of another document is refused", async () => {
  const doc = await openSegment(0);
  const other = await BdfDocument.open(new BufferSource(await read("testdata/demo.bdf")));
  assert.throws(() => doc.addSegment(other), /not a segment document/);
  const seg = await openSegment(3);
  seg.manifest.views[0].pages.pop();
  assert.throws(() => doc.addSegment(seg), /another document/);
});

test("the loader asks for each segment once, with what it holds", async () => {
  const realFetch = globalThis.fetch;
  const realGenerate = crypto.subtle.generateKey;
  const asked = [];
  // the fixtures are sealed for the test key: the loader gets it as the key of each request
  crypto.subtle.generateKey = async () => keyPair;
  globalThis.fetch = async (url, init) => {
    const body = JSON.parse(init.body);
    asked.push(body);
    assert.equal(url, "https://books.example/segments/basic");
    assert.equal(init.method, "POST");
    assert.equal(init.headers["Content-Type"], "application/json");
    assert.equal(init.credentials, "same-origin");
    assert.equal(init.cache, "no-store");
    assert.equal(atob(body.key).length, 65);
    if (body.page === 6) return new Response("", { status: 403 });
    const from = body.page - (body.page % 3);
    return new Response(segmentBytes[from]);
  };
  try {
    const loader = await SegmentLoader.open("https://books.example/segments/basic", { init: { credentials: "same-origin" }, ahead: 0 });
    const came = [];
    loader.onSegment = (s) => came.push(s.from);
    assert.equal(loader.has("pages", 2), true);
    assert.equal(loader.has("pages", 3), false);
    // pages of one segment asked for at once: one request
    await Promise.all([loader.ensure("pages", 3), loader.ensure("pages", 4), loader.ensure("pages", 5)]);
    assert.deepEqual(came, [3]);
    assert.deepEqual(asked.map((b) => [b.page, b.have]), [
      [0, []],
      [3, [{ view: "pages", from: 0, to: 3 }]],
    ]);
    await assert.rejects(loader.ensure("pages", 6), (e) => e instanceof BdfSegmentError && e.status === 403);
    await assert.rejects(loader.ensure("pages", 7), /no page 7/);
    const doc = loader.doc;
    const obj = await doc.ensure(doc.view("pages").pages[4].layers[0].obj);
    assert.ok(obj.ops.length > 0);
  } finally {
    globalThis.fetch = realFetch;
    crypto.subtle.generateKey = realGenerate;
  }
});
