// What a document may not make a reader do: documents that ask for more
// are refused with BdfFormatError.
import { test } from "node:test";
import assert from "node:assert/strict";
import {
  BdfDocument, BdfFormatError, BufferSource, RangeSource, SplitSource, UseLimits, decode, extractText, extractContent, tileSize,
  MAX_MANIFEST_SIZE, MAX_REUSED_INSTRUCTIONS, MAX_TEXT_RUNS, MAX_USE_DEPTH,
} from "../dist/index.js";
import { single, object, fanOut, page, name, cat, bytes, varuint, f32, systemFont, Op } from "./build.mjs";

const open = (file) => BdfDocument.open(new BufferSource(file));
const formatError = (message) => (e) => e instanceof BdfFormatError && message.test(e.message);
const text = { strings: ["x"], fonts: [systemFont()], ops: [bytes(Op.FONT), varuint(0), f32(10), bytes(Op.FILL_TEXT), varuint(0), f32(0, 10, 5)] };

test("the limits are those of the Go reader", () => {
  assert.equal(MAX_USE_DEPTH, 64);
  assert.equal(MAX_REUSED_INSTRUCTIONS, 1 << 27);
  assert.equal(MAX_TEXT_RUNS, 1 << 22);
  assert.equal(MAX_MANIFEST_SIZE, 256 << 20);
});

test("a part inflates to the size its manifest states at most", async () => {
  const zeros = new Uint8Array(1 << 20);
  const doc = await open(single({}, [
    { h: name(1), t: "img", bytes: zeros, deflate: true },
    { h: name(2), t: "img", bytes: zeros, deflate: true, size: zeros.length - 1 },
    { h: name(3), t: "img", bytes: zeros, deflate: true, size: -1 },
    { h: name(4), t: "img", bytes: zeros, deflate: true, entry: { size: undefined } },
    { h: name(5), t: "img", bytes: zeros, size: 1 },
  ]));
  assert.equal((await doc.part(name(1))).length, zeros.length);
  await assert.rejects(doc.part(name(2)), formatError(/inflates to more than its stated size/));
  await assert.rejects(doc.part(name(3)), formatError(/bad part size/));
  // no size stated: none to inflate to
  await assert.rejects(doc.part(name(4)), formatError(/inflates to more/));
  // stored as it is, a part is as large as the file has it
  assert.equal((await doc.part(name(5))).length, zeros.length);
});

test("decode stops at the limit", async () => {
  const { deflateRawSync } = await import("node:zlib");
  const data = new Uint8Array(deflateRawSync(new Uint8Array(3 << 20)));
  assert.equal((await decode(data, "deflate-raw", 3 << 20)).length, 3 << 20);
  await assert.rejects(decode(data, "deflate-raw", (3 << 20) - 1), formatError(/inflates to more/));
  await assert.rejects(decode(data, "deflate-raw", 0), formatError(/inflates to more/));
});

test("a manifest inflates to 256 MiB at most", async () => {
  const pad = "x".repeat(MAX_MANIFEST_SIZE);
  const file = single({ meta: { generator: pad } }, [], { deflateManifest: true });
  assert.ok(file.length < 2 << 20, `${file.length} bytes`);
  await assert.rejects(open(file), formatError(/inflates to more/));
});

test("a stored manifest is bounded before its bytes are requested", async () => {
  const file = single({});
  new DataView(file.buffer, file.byteOffset).setBigUint64(16, BigInt(MAX_MANIFEST_SIZE + 1), true);
  await assert.rejects(open(file), formatError(/manifest too large/));
  const realFetch = globalThis.fetch;
  let requests = 0;
  globalThis.fetch = async () => {
    requests++;
    return new Response(file.subarray(0, 32), { status: 206 });
  };
  try {
    await assert.rejects(new RangeSource("https://example.com/doc.bdf").manifest(), formatError(/manifest too large/));
    assert.equal(requests, 1);
  } finally {
    globalThis.fetch = realFetch;
  }
});

test("parts and the manifest lie within the file", async () => {
  const file = single({}, [{ h: name(1), t: "img", bytes: bytes(1, 2, 3), entry: { off: -2 } }, { h: name(2), t: "img", bytes: bytes(1, 2, 3), entry: { len: 1e9 } }]);
  const doc = await open(file);
  await assert.rejects(doc.part(name(1)), formatError(/part out of range/));
  await assert.rejects(doc.part(name(2)), formatError(/part out of range/));
  const cut = file.slice(0, 40);
  await assert.rejects(open(cut), formatError(/manifest out of range/));
});

test("an object that draws itself is refused", async () => {
  const self = object({ objects: [name(1)], ops: [bytes(Op.USE), varuint(0)] });
  const doc = await open(single({ views: [{ id: "v", kind: "fixed", pages: [page(name(1))] }] }, [{ h: name(1), t: "obj", bytes: self }]));
  const top = await doc.ensure(name(1));
  assert.throws(() => extractText(top, (h) => doc.objectSync(h)), formatError(/too deep/));
  assert.throws(() => extractContent(top, (h) => doc.objectSync(h)), formatError(/too deep/));
  await assert.rejects(doc.textIndex(doc.view("v")), formatError(/too deep/));
});

test("objects draw objects 64 deep", async () => {
  // a chain: the object at depth 64 draws text, the one at depth 65 does too
  const chain = async (depth) => {
    const doc = await open(fanOut(depth, 1, text));
    return extractText(await doc.ensure(name(1)), (h) => doc.objectSync(h));
  };
  assert.deepEqual((await chain(MAX_USE_DEPTH)).map((r) => r.text), ["x"]);
  await assert.rejects(chain(MAX_USE_DEPTH + 1), formatError(/too deep/));
});

test("the instructions of objects drawn again are counted", async () => {
  // 10 objects of 10 USE each, 3 deep: the last one is drawn 1000 times
  const doc = await open(fanOut(3, 10, text));
  const top = await doc.ensure(name(1));
  const resolve = (h) => doc.objectSync(h);
  assert.equal(extractText(top, resolve).length, 1000);
  // 10, 100 and 1000 instructions of the three, and 2 of the last object each time
  const all = 10 + 100 + 1000 + 2 * 1000;
  // what is walked for the first time does not count: each object once
  const again = all - (10 + 10 + 10 + 2);
  assert.equal(extractText(top, resolve, undefined, new UseLimits(again)).length, 1000);
  assert.throws(() => extractText(top, resolve, undefined, new UseLimits(again - 1)), formatError(/drawn too many times/));
  // one count for the objects walked together: the second time, all but the top-level object is walked again
  const both = again + all - 10;
  const enough = new UseLimits(both), short = new UseLimits(both - 1);
  for (const limits of [enough, enough, short]) extractText(top, resolve, undefined, limits);
  assert.throws(() => extractText(top, resolve, undefined, short), formatError(/drawn too many times/));
});

test("the runs of a walk are counted", async () => {
  const doc = await open(fanOut(2, 10, text));
  const top = await doc.ensure(name(1));
  const resolve = (h) => doc.objectSync(h);
  assert.equal(extractText(top, resolve, undefined, new UseLimits(undefined, 100)).length, 100);
  assert.throws(() => extractText(top, resolve, undefined, new UseLimits(undefined, 99)), formatError(/too many text runs/));
});

test("part names are 32 hex characters", async () => {
  const bad = ["../../../api/delete?x=", "0".repeat(31), "0".repeat(33), "A".repeat(32), "", 12, null];
  for (const h of bad) {
    await assert.rejects(open(single({}, [{ h, t: "img", bytes: bytes(1) }])), formatError(/bad part name/), JSON.stringify(h));
    await assert.rejects(open(single({}, [{ h: name(1), t: "img", bytes: bytes(1), entry: { sealed: h ?? 0 } }])), formatError(/bad sealed part name/), JSON.stringify(h));
  }
});

test("a split document does not fetch what is not a part name", async () => {
  const evil = "../../../api/account/delete?confirm=1&x=";
  const manifest = { bdf: 1, opset: 1, unit: "pt", views: [], parts: [{ h: evil, t: "img", enc: "identity", len: 1, size: 1 }] };
  const fetched = [];
  const realFetch = globalThis.fetch;
  globalThis.fetch = async (url) => {
    fetched.push(String(url));
    return new Response(url.endsWith("manifest.json") ? JSON.stringify(manifest) : "x", { status: 200 });
  };
  try {
    const source = new SplitSource("https://example.com/files/1234/");
    await assert.rejects(BdfDocument.open(source), formatError(/bad part name/));
    await assert.rejects(source.stored(manifest.parts[0]), formatError(/bad part name/));
    assert.deepEqual(fetched, ["https://example.com/files/1234/manifest.json"]);
  } finally {
    globalThis.fetch = realFetch;
  }
});

test("a server that ignores Range is asked for the file once", async () => {
  const leaf = object(text);
  const file = single({ views: [{ id: "v", kind: "fixed", pages: [page(name(1))] }] }, [
    { h: name(1), t: "obj", bytes: object({ objects: [name(2)], ops: [bytes(Op.USE), varuint(0)] }) },
    { h: name(2), t: "obj", bytes: leaf, deflate: true },
  ]);
  let requests = 0;
  const realFetch = globalThis.fetch;
  globalThis.fetch = async () => {
    requests++;
    return new Response(file, { status: 200 });
  };
  try {
    const doc = await BdfDocument.open(new RangeSource("https://example.com/doc.bdf"));
    const top = await doc.ensure(name(1));
    assert.deepEqual(extractText(top, (h) => doc.objectSync(h)).map((r) => r.text), ["x"]);
    assert.equal(requests, 1);
    // a part that is not in the file
    await assert.rejects(doc.source.stored({ h: name(3), t: "img", enc: "identity", len: 10, size: 10, off: file.length }), formatError(/part out of range/));
  } finally {
    globalThis.fetch = realFetch;
  }
});

test("a server that answers ranges is asked for each", async () => {
  const file = single({}, [{ h: name(1), t: "img", bytes: bytes(1, 2, 3) }, { h: name(2), t: "img", bytes: bytes(4, 5) }]);
  const ranges = [];
  const realFetch = globalThis.fetch;
  globalThis.fetch = async (url, init) => {
    const [, from, to] = /^bytes=(\d+)-(\d+)$/.exec(init.headers.Range);
    ranges.push([Number(from), Number(to)]);
    return new Response(file.slice(Number(from), Number(to) + 1), { status: 206 });
  };
  try {
    const doc = await BdfDocument.open(new RangeSource("https://example.com/doc.bdf"));
    assert.deepEqual([...await doc.part(name(2))], [4, 5]);
    assert.deepEqual([...await doc.part(name(1))], [1, 2, 3]);
    assert.equal(ranges.length, 4);
    assert.deepEqual(ranges[0], [0, 31]);
  } finally {
    globalThis.fetch = realFetch;
  }
});

test("range responses must contain exactly the requested bytes", async () => {
  const realFetch = globalThis.fetch;
  try {
    for (const [body, headers, message] of [
      [new Uint8Array(31), {}, /shorter than its stated size/],
      [new Uint8Array(33), {}, /exceeds its stated size/],
      [new Uint8Array(32), { "Content-Range": "bytes 1-32/100" }, /does not match the request/],
    ]) {
      globalThis.fetch = async () => new Response(body, { status: 206, headers });
      await assert.rejects(new RangeSource("https://example.com/doc.bdf").manifest(), formatError(message));
    }
  } finally {
    globalThis.fetch = realFetch;
  }
});

test("split responses cannot exceed the sizes in the manifest", async () => {
  const realFetch = globalThis.fetch;
  try {
    globalThis.fetch = async () => new Response("abc", { status: 200 });
    const source = new SplitSource("https://example.com/doc/");
    await assert.rejects(source.stored({ h: name(1), t: "img", enc: "identity", len: 2, size: 2 }), formatError(/exceeds its stated size/));
    await assert.rejects(source.stored({ h: name(1), t: "img", enc: "identity", len: 4, size: 4 }), formatError(/shorter than its stated size/));
    globalThis.fetch = async () => new Response("abc", { status: 200, headers: { "Content-Encoding": "gzip", "Content-Length": "23" } });
    assert.deepEqual([...await source.stored({ h: name(1), t: "img", enc: "identity", len: 3, size: 3 })], [97, 98, 99]);
    globalThis.fetch = async () => new Response("x", { status: 200, headers: { "Content-Length": String(MAX_MANIFEST_SIZE + 1) } });
    await assert.rejects(source.manifest(), formatError(/manifest exceeds its stated size/));
  } finally {
    globalThis.fetch = realFetch;
  }
});

test("the numbers of a view are checked when the document opens", async () => {
  const view = (v) => open(single({ views: [{ id: "v", kind: "sheet", ...v }] }));
  const pages = (p, v = {}) => open(single({ views: [{ id: "v", kind: "flow", pages: [{ w: 100, h: 100, layers: [], ...p }], ...v }] }));
  // Excel's own size, hidden rows, and the smallest row that is one
  const doc = await view({ tile: 2048, cols: [[16384, 64]], rows: [[1048576, 20], [10, 0], [1, 1 / 64]], freeze: { rows: 1 } });
  assert.equal(tileSize(doc.view("v")), 2048);
  // a tile size of 0 or less is the default one, as it is for the Go reader
  assert.equal(tileSize((await view({ tile: 0 })).view("v")), 2048);
  assert.equal(tileSize((await view({ tile: -5 })).view("v")), 2048);
  assert.equal(tileSize((await view({})).view("v")), 2048);
  assert.equal(tileSize((await view({ tile: 512 })).view("v")), 512);
  for (const bad of [
    { tile: 0.001 }, { tile: "2048" },
    { rows: [[1e11, 0]] }, { rows: [[1 << 23, 20], [(1 << 23) + 1, 20]] }, { cols: [[1.5, 20]] }, { cols: [[-1, 20]] },
    { rows: [[10, 1e-7]] }, { rows: [[10, -1]] }, { cols: [[10, "20"]] }, { cols: [[10, 1e9]] }, { rows: [10, 20] }, { rows: "many" },
    { freeze: { rows: -1 } }, { freeze: { cols: 1.5 } },
  ]) {
    await assert.rejects(view(bad), BdfFormatError, JSON.stringify(bad));
  }
  await pages({ w: 0, h: 0 });
  await pages({ body: { x: -10, y: 0, w: 50, h: 50 } }, { continuous: { gap: 24 } });
  for (const bad of [{ w: 1e9 }, { h: -1 }, { w: "100" }, { h: null }, { body: { x: 0, y: 0, w: 1e30, h: 10 } }, { body: { x: 0, y: "0", w: 10, h: 10 } }]) {
    await assert.rejects(pages(bad), BdfFormatError, JSON.stringify(bad));
  }
  // JSON has no infinity, but what is read of a split document may be anything
  for (const gap of [1e300, "24", null]) await assert.rejects(pages({}, { continuous: { gap } }), BdfFormatError, JSON.stringify(gap));
  await assert.rejects(open(single({ views: [{ kind: "fixed" }] })), formatError(/view without an id/));
});
