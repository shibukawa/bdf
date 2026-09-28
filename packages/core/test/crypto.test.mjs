import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { BufferSource, RangeSource, BdfDocument, BdfPasswordError, BdfFormatError, MAX_ITERATIONS, MAX_MANIFEST_SIZE, parseHeader } from "../dist/index.js";

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

test("a sealed manifest is bounded before its bytes are requested", async () => {
  const outer = await new BufferSource(sealedBytes).manifest();
  outer.parts.find((p) => p.h === outer.encryption.manifest.part).len = MAX_MANIFEST_SIZE + 12 + 16 + 1;
  const source = {
    manifest: async () => outer,
    stored: () => assert.fail("oversized sealed manifest must not be loaded"),
  };
  await assert.rejects(BdfDocument.open(source, { password: PASSWORD }), (e) => e instanceof BdfFormatError && /sealed manifest size/.test(e.message));
});

test("the iterations of the key slots are limited one by one and in sum", async () => {
  const outer = await new BufferSource(sealedBytes).manifest();
  const [slot] = outer.encryption.keys;
  const open = (keys) => BdfDocument.open({ manifest: async () => ({ ...outer, encryption: { ...outer.encryption, keys } }), stored: () => assert.fail("no part is read") }, { password: "wrong" });
  const tooMany = (e) => e instanceof BdfFormatError && /key slots take more than 10000000 iterations/.test(e.message);
  for (const iter of [0, MAX_ITERATIONS + 1, 1.5]) {
    await assert.rejects(open([{ ...slot, iter }]), (e) => e instanceof BdfFormatError && /out of range/.test(e.message), `${iter}`);
  }
  // slots that are tried: the second one would take the sum past the limit, and is not
  const t0 = performance.now();
  await assert.rejects(open([{ ...slot, iter: 1 }, { ...slot, iter: MAX_ITERATIONS }]), tooMany);
  await assert.rejects(open([{ ...slot, iter: 1000 }, { ...slot, iter: 1000 }, { ...slot, iter: MAX_ITERATIONS - 1999 }]), tooMany);
  assert.ok(performance.now() - t0 < 1000, "the last slot was derived");
  // up to the limit they are all tried
  await assert.rejects(open([{ ...slot, iter: 1000 }, { ...slot, iter: 1000 }]), (e) => e instanceof BdfPasswordError && e.reason === "wrong");
});
