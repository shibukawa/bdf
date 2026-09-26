import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { BufferSource, RangeSource, BdfDocument, BdfPasswordError, BdfFormatError, parseHeader } from "../dist/index.js";

// testdata/demo-encrypted.bdf is testdata/demo.bdf encrypted by the Go writer
// (npm run testdata) with this password.
const PASSWORD = "demo-パスワード";
const root = new URL("../../../", import.meta.url);
const plainBytes = new Uint8Array(await readFile(new URL("testdata/demo.bdf", root)));
const sealedBytes = new Uint8Array(await readFile(new URL("testdata/demo-encrypted.bdf", root)));

test("header flag and outer manifest", async () => {
  assert.equal(parseHeader(sealedBytes).flags, 1);
  const m = await new BufferSource(sealedBytes).manifest();
  assert.equal(m.encryption.cipher, "A256GCM");
  assert.equal(m.views, undefined);
  assert.ok(m.parts.every((p) => p.t === "sealed"));
  const text = new TextDecoder("latin1").decode(sealedBytes);
  assert.ok(!text.includes("BDF fixture"), "the title is stored in the clear");
});

test("a password is required, and must be right", async () => {
  const source = new BufferSource(sealedBytes);
  await assert.rejects(BdfDocument.open(source), (e) => e instanceof BdfPasswordError && e.reason === "required");
  for (const pw of ["", "wrong", PASSWORD + " "]) {
    await assert.rejects(BdfDocument.open(source, { password: pw }), (e) => e instanceof BdfPasswordError && e.reason === "wrong", pw);
  }
});

test("decrypted parts equal the plain document's", async () => {
  const plain = await BdfDocument.open(new BufferSource(plainBytes));
  const doc = await BdfDocument.open(new BufferSource(sealedBytes), { password: PASSWORD });
  assert.deepEqual(doc.manifest.views, plain.manifest.views);
  assert.deepEqual(doc.manifest.meta, plain.manifest.meta);
  assert.equal(doc.manifest.parts.length, plain.manifest.parts.length);
  for (const e of plain.manifest.parts) {
    assert.deepEqual(await doc.part(e.h), await plain.part(e.h), e.h);
  }
  const slide = doc.view("slides").pages[0].layers.at(-1).obj;
  const obj = await doc.ensure(slide);
  assert.ok(obj.ops.length > 0);
});

test("passwords are compared in NFC", async () => {
  const nfd = PASSWORD.normalize("NFD");
  assert.notEqual(nfd, PASSWORD);
  const doc = await BdfDocument.open(new BufferSource(sealedBytes), { password: nfd });
  assert.equal(doc.manifest.views.length, 3);
});

test("range requests read sealed parts", async () => {
  const realFetch = globalThis.fetch;
  let ranges = 0;
  globalThis.fetch = async (_url, init) => {
    const m = /bytes=(\d+)-(\d+)/.exec(init.headers.Range);
    ranges++;
    return new Response(sealedBytes.slice(Number(m[1]), Number(m[2]) + 1), { status: 206 });
  };
  try {
    const doc = await BdfDocument.open(new RangeSource("https://example.invalid/demo.bdf"), { password: PASSWORD });
    const tile = await doc.ensure(doc.view("sheet1").tiles["0,0"]);
    assert.ok(tile.ops.length > 0);
    assert.ok(ranges >= 4, `${ranges} range requests`);
  } finally {
    globalThis.fetch = realFetch;
  }
});

test("a tampered part fails authentication", async () => {
  const bytes = sealedBytes.slice();
  const source = new BufferSource(bytes);
  const outer = await source.manifest();
  const doc = await BdfDocument.open(source, { password: PASSWORD });
  const e = doc.manifest.parts.find((p) => p.t === "obj");
  const o = outer.parts.find((p) => p.h === e.sealed);
  const h = parseHeader(bytes);
  bytes[h.manifestOff + h.manifestLen + o.off + 20] ^= 1;
  await assert.rejects(doc.part(e.h), (err) => err instanceof BdfFormatError && /authentication/.test(err.message));
});
