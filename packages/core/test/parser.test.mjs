import { test } from "node:test";
import assert from "node:assert/strict";
import { ByteReader, BdfFormatError, decodeObject, decodePath, decodeTextIndex, walk, NoopSink, Op } from "../dist/index.js";
import { ascii, bytes, cat, f32, object, u16, varuint } from "./build.mjs";

test("varuint preserves safe integers and refuses rounded values", () => {
  for (const n of [0, 127, 128, 0xffffffff, 2 ** 32, Number.MAX_SAFE_INTEGER]) {
    assert.equal(new ByteReader(varuint(n)).varuint(), n);
  }
  for (const n of [2n ** 53n, 2n ** 53n + 1n, 2n ** 64n - 1n]) {
    assert.throws(() => new ByteReader(varuint(n)).varuint(), BdfFormatError);
  }
});

test("invalid byte counts cannot move the reader backwards or lose its position", () => {
  const reader = new ByteReader(bytes(1, 2, 3));
  assert.equal(reader.u8(), 1);
  for (const n of [-1, 0.5, NaN, Infinity, Number.MAX_SAFE_INTEGER + 1]) {
    assert.throws(() => reader.bytesN(n), BdfFormatError);
    assert.throws(() => reader.f32array(n), BdfFormatError);
    assert.equal(reader.pos, 1);
  }
  assert.deepEqual(reader.bytesN(2), bytes(2, 3));
});

test("truncated float arrays are refused before allocating decoded storage", (t) => {
  // Observe the allocation without risking the gigabytes a forged DASH
  // count asks for, including when this regression test runs against a bug.
  const dash = decodeObject(object({ ops: [bytes(Op.DASH), varuint(1 << 28)] }));
  const FloatArray = globalThis.Float32Array;
  let allocations = 0;
  t.mock.method(globalThis, "Float32Array", new Proxy(FloatArray, {
    construct() {
      allocations++;
      throw new Error("unexpected float array allocation");
    },
  }));
  assert.throws(() => new ByteReader(bytes()).f32array(1 << 28), BdfFormatError);
  assert.throws(() => walk(dash, new NoopSink()), BdfFormatError);
  // An inline path also computes its float count from untrusted verbs.
  assert.throws(() => decodePath(new ByteReader(cat(varuint(1), bytes(1)))), BdfFormatError);
  assert.equal(allocations, 0);
});

test("truncated glyph runs and text indexes are refused before allocating arrays", (t) => {
  const ArrayClass = globalThis.Array;
  let allocations = 0;
  t.mock.method(globalThis, "Array", new Proxy(ArrayClass, {
    construct() {
      allocations++;
      throw new Error("unexpected array allocation");
    },
  }));
  const glyphs = decodeObject(object({ ops: [bytes(Op.FILL_PATH_RUN, 0), varuint(1 << 28)] }));
  assert.throws(() => walk(glyphs, new NoopSink()), BdfFormatError);
  assert.throws(() => decodeTextIndex(cat(ascii("BTXT"), u16(1), varuint(1 << 28))), BdfFormatError);
  assert.equal(allocations, 0);
});

test("minimum-size glyph and text records still decode", () => {
  let received;
  class Sink extends NoopSink {
    fillPathRun(rule, glyphs) { received = { rule, glyphs }; }
  }
  const glyphs = decodeObject(object({ ops: [bytes(Op.FILL_PATH_RUN, 0), varuint(1), varuint(0), f32(1, 2)] }));
  walk(glyphs, new Sink());
  assert.deepEqual(received, { rule: 0, glyphs: [{ path: 0, x: 1, y: 2 }] });
  assert.deepEqual(decodeTextIndex(cat(ascii("BTXT"), u16(1), varuint(1), bytes(0, 0, 0, 0, 0))), [
    { a: 0, b: 0, ordinal: 0, sep: 0, text: "" },
  ]);
  assert.deepEqual([...new ByteReader(f32(1, -2)).f32array(2)], [1, -2]);
});
