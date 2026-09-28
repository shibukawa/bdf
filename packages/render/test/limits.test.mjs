// What a document may not make the renderer do (see also the limits of @bdf/core).
import { test } from "node:test";
import assert from "node:assert/strict";
import { BdfDocument, BdfFormatError, BufferSource, UseLimits } from "@bdf/core";
import { ResourceCache, imageSize } from "../dist/resources.js";
import { CanvasRenderer, filterAllowed, MAX_GROUP_DEPTH } from "../dist/canvas.js";
import { PageRenderer, tilesIn } from "../dist/page.js";
import { svgSize } from "../dist/svg.js";
import { single, object, fanOut, page, name, cat, bytes, ascii, varuint, f32, systemFont, Op } from "../../core/test/build.mjs";

// Node has no fonts to load.
globalThis.FontFace ??= class { load() { return Promise.resolve(this); } };
const fonts = { check: () => true, add() {}, delete() {} };

const open = (file) => BdfDocument.open(new BufferSource(file));
const formatError = (message) => (e) => e instanceof BdfFormatError && message.test(e.message);
const text = { strings: ["x"], fonts: [systemFont()], ops: [bytes(Op.FONT), varuint(0), f32(10), bytes(Op.FILL_TEXT), varuint(0), f32(0, 10, 5)] };
const fixed = (pages) => ({ views: [{ id: "v", kind: "fixed", pages }] });

/** A 2D context on a canvas of a size, scaled by 2, that counts the calls of each method. */
function context(width = 200, height = 200) {
  const calls = {};
  const state = { font: "10px sans-serif", canvas: { width, height }, calls, shadowOffsetX: 0, shadowOffsetY: 0, shadowBlur: 0, filter: "none" };
  return new Proxy(state, {
    get(t, k) {
      if (k in t) return t[k];
      if (k === "getTransform") return () => ({ a: 2, b: 0, c: 0, d: 2, e: 0, f: 0, transformPoint: ({ x, y }) => ({ x: 2 * x, y: 2 * y }) });
      if (k === "measureText") return (s) => ({ width: s.length * 5 });
      return () => { calls[k] = (calls[k] ?? 0) + 1; };
    },
    set(t, k, v) { t[k] = v; return true; },
  });
}

/** Temporary canvases that note their sizes. */
function canvases() {
  const made = [];
  const createCanvas = (w, h) => {
    const canvas = { width: w, height: h, getContext: () => ctx };
    const ctx = context(w, h);
    ctx.canvas = canvas;
    made.push({ w, h, canvas });
    return canvas;
  };
  return { made, createCanvas };
}

async function prepared(file, opts, resources) {
  const doc = await open(file);
  const pr = new PageRenderer(doc, opts, fonts, resources);
  const p = doc.view("v").pages[0];
  await pr.preparePage(p);
  return { doc, pr, page: p, top: pr.res.object(p.layers[0].obj) };
}

test("an object that draws itself is not drawn", async () => {
  const self = object({ objects: [name(1)], ops: [bytes(Op.USE), varuint(0)] });
  const { pr, page: p, top } = await prepared(single(fixed([page(name(1))]), [{ h: name(1), t: "obj", bytes: self }]));
  assert.throws(() => pr.renderer.draw(context(), top), formatError(/too deep/));
  await assert.rejects(pr.renderPage(context(), p, { scale: 1 }), formatError(/too deep/));
  // the renderer draws the next page as ever
  const ok = await prepared(fanOut(3, 10, text));
  const ctx = context();
  await ok.pr.renderPage(ctx, ok.page, { scale: 1 });
  assert.equal(ctx.calls.fillText, 1000);
});

test("the instructions of objects drawn again are counted for a render", async () => {
  const { pr, top } = await prepared(fanOut(3, 10, text));
  // as for the text (the limits test of @bdf/core): every object is walked once for nothing
  const again = 10 + 100 + 1000 + 2 * 1000 - (10 + 10 + 10 + 2);
  const draw = (limits) => {
    const ctx = context();
    pr.renderer.draw(ctx, top, true, undefined, limits);
    return ctx.calls.fillText;
  };
  assert.equal(draw(new UseLimits(again)), 1000);
  assert.throws(() => draw(new UseLimits(again - 1)), formatError(/drawn too many times/));
  // the objects of a render are counted together, those of the next render anew
  const limits = new UseLimits(again);
  draw(limits);
  assert.throws(() => draw(limits), formatError(/drawn too many times/));
  assert.equal(draw(), 1000);
  assert.equal(draw(), 1000);
});

const group = (x, y, w, h) => cat(bytes(Op.GROUP_BEGIN), f32(1), bytes(0), f32(x, y, w, h));
const rect = cat(bytes(Op.FILL_RECT), f32(0, 0, 10, 10));

test("the canvas of a group is the part of it on the canvas it is drawn on", async () => {
  const ops = [group(0, 0, 8000, 8000), group(10, 20, 8000, 8000), rect, bytes(Op.GROUP_END), group(-50, -50, 60, 70), rect, bytes(Op.GROUP_END), group(500, 500, 10, 10), bytes(Op.GROUP_END), bytes(Op.GROUP_END)];
  const { made, createCanvas } = canvases();
  const { pr, top } = await prepared(single(fixed([page(name(1))]), [{ h: name(1), t: "obj", bytes: object({ ops }) }]), { createCanvas });
  pr.renderer.draw(context(200, 150), top);
  // in device pixels (the scale is 2): cut to the 200 by 150 canvas, then to the canvas of the group around
  assert.deepEqual(made.map(({ w, h }) => [w, h]), [[200, 150], [180, 110], [20, 40], [1, 1]]);
  // they are let go once drawn
  assert.deepEqual(made.map(({ canvas }) => [canvas.width, canvas.height]), [[0, 0], [0, 0], [0, 0], [0, 0]]);
});

test("the canvas of a group reaches as far as its shadow or filter", async () => {
  const draw = async (set) => {
    const { made, createCanvas } = canvases();
    const { pr, top } = await prepared(single(fixed([page(name(1))]), [{ h: name(1), t: "obj", bytes: object({ ops: [group(-5000, -5000, 20000, 20000), rect, bytes(Op.GROUP_END)] }) }]), { createCanvas });
    const ctx = context(200, 150);
    Object.assign(ctx, set);
    // the state of the canvas when the group begins is that it is composited with
    pr.renderer.draw(ctx, top, false);
    return [made[0].w, made[0].h];
  };
  assert.deepEqual(await draw({}), [200, 150]);
  // 6 px to the side and a blur of 4 px reach 12 px
  assert.deepEqual(await draw({ shadowOffsetX: -6, shadowOffsetY: 2, shadowBlur: 4 }), [224, 174]);
  assert.deepEqual(await draw({ filter: "blur(2px) drop-shadow(1px 1px 3px #000)" }), [200 + 2 * 21, 150 + 2 * 21]);
  // what a document says of its shadow is limited too
  assert.deepEqual(await draw({ shadowBlur: 1e9 }), [200 + 2 * 4096, 150 + 2 * 4096]);
  assert.deepEqual(await draw({ shadowOffsetX: Infinity }), [200 + 2 * 4096, 150 + 2 * 4096]);
});

test("a group within the canvas has the canvas it had", async () => {
  const ops = [group(10, 10, 50.2, 30), rect, bytes(Op.GROUP_END)];
  const { made, createCanvas } = canvases();
  const { pr, top } = await prepared(single(fixed([page(name(1))]), [{ h: name(1), t: "obj", bytes: object({ ops }) }]), { createCanvas });
  pr.renderer.draw(context(200, 150), top);
  assert.deepEqual(made.map(({ w, h }) => [w, h]), [[101, 60]]);
});

test("groups and masks nest 64 deep", async () => {
  const nested = (n) => object({ ops: [...Array.from({ length: n }, () => group(0, 0, 10, 10)), rect, ...Array.from({ length: n }, () => bytes(Op.GROUP_END))] });
  const draw = async (bytes) => {
    const { made, createCanvas } = canvases();
    const { pr, top } = await prepared(single(fixed([page(name(1))]), [{ h: name(1), t: "obj", bytes }]), { createCanvas });
    pr.renderer.draw(context(), top);
    return made.length;
  };
  assert.equal(await draw(nested(MAX_GROUP_DEPTH)), 64);
  await assert.rejects(draw(nested(MAX_GROUP_DEPTH + 1)), formatError(/groups nested too deep/));
  const mask = cat(bytes(Op.MASK_BEGIN, 0), bytes(0, 0, 0, 0), varuint(0));
  const masks = object({ ops: [group(0, 0, 10, 10), ...Array.from({ length: MAX_GROUP_DEPTH }, () => mask)] });
  await assert.rejects(draw(masks), formatError(/groups nested too deep/));
});

test("FILTER takes filter functions only", async () => {
  for (const ok of ["", "none", "blur(2px)", "blur(1.5px) brightness(120%)", "drop-shadow(2px 2px 4px rgba(0, 0, 0, 0.5))", "drop-shadow(0 0 3px #336699)", "hue-rotate(-90deg) SEPIA(1)", "opacity(50%) drop-shadow(1px 1px 0 rgb(0 0 0 / 50%))"]) {
    assert.ok(filterAllowed(ok), ok);
  }
  for (const bad of ["url(#f)", "URL(https://example.com/f.svg#f)", "blur(2px) url(a.svg#f)", "u\\72l(#f)", 'url("#f")', "blur(2px); color: red", "image-set(a)", "var(--f)", "(1)", "blur(2px) /* */", "drop-shadow(0 0 url(x))"]) {
    assert.ok(!filterAllowed(bad), bad);
  }
  // the filter of the context as the instruction leaves it
  const filtered = async (css) => {
    const o = object({ strings: [css], ops: [bytes(Op.FILTER), varuint(0)] });
    const { pr, top } = await prepared(single(fixed([page(name(1))]), [{ h: name(1), t: "obj", bytes: o }]));
    const ctx = context();
    ctx.filter = "";
    pr.renderer.draw(ctx, top);
    return ctx.filter;
  };
  assert.equal(await filtered("blur(2px)"), "blur(2px)");
  assert.equal(await filtered("url(https://example.com/f.svg#f)"), "none");
});

/** The start of a PNG image that says it is width by height pixels. */
const png = (width, height) => {
  const b = new Uint8Array(33);
  b.set([0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0, 0, 0, 13]);
  b.set(ascii("IHDR"), 12);
  new DataView(b.buffer).setUint32(16, width);
  new DataView(b.buffer).setUint32(20, height);
  return b;
};

test("the size of an image is read from its header", () => {
  assert.deepEqual(imageSize(png(640, 480)), { width: 640, height: 480 });
  assert.deepEqual(imageSize(png(30000, 30000)), { width: 30000, height: 30000 });
  assert.equal(imageSize(png(640, 480).subarray(0, 20)), undefined);
  assert.deepEqual(imageSize(cat(ascii("GIF89a"), bytes(0x80, 0x02, 0xe0, 0x01, 0))), { width: 640, height: 480 });
  const bmp = new Uint8Array(30);
  bmp.set(ascii("BM"));
  const dv = new DataView(bmp.buffer);
  dv.setUint32(14, 40, true);
  dv.setInt32(18, 640, true);
  dv.setInt32(22, -480, true); // top-down
  assert.deepEqual(imageSize(bmp), { width: 640, height: 480 });
  // JPEG: the frame header after an APP0 segment and fill bytes
  const jpeg = cat(bytes(0xff, 0xd8, 0xff, 0xe0, 0, 4, 1, 2, 0xff, 0xff, 0xc2, 0, 17, 8), bytes(0x01, 0xe0, 0x02, 0x80), bytes(3, 0, 0, 0));
  assert.deepEqual(imageSize(jpeg), { width: 640, height: 480 });
  assert.equal(imageSize(bytes(0xff, 0xd8, 0xff, 0xda, 0, 2, 0, 0)), undefined);
  const webp = (kind, rest) => cat(ascii("RIFF"), bytes(0, 0, 0, 0), ascii("WEBP"), ascii(kind), bytes(0, 0, 0, 0), rest, bytes(0, 0, 0, 0, 0, 0));
  assert.deepEqual(imageSize(webp("VP8X", bytes(0, 0, 0, 0, 0x7f, 0x02, 0, 0xdf, 0x01, 0))), { width: 640, height: 480 });
  assert.deepEqual(imageSize(webp("VP8 ", bytes(0, 0, 0, 0x9d, 0x01, 0x2a, 0x80, 0x02, 0xe0, 0x01))), { width: 640, height: 480 });
  // 14 bits of the width less one, then of the height less one
  const v = 639 | (479 << 14);
  assert.deepEqual(imageSize(webp("VP8L", bytes(0x2f, v & 0xff, (v >> 8) & 0xff, (v >> 16) & 0xff, (v >> 24) & 0xff))), { width: 640, height: 480 });
  for (const other of [ascii("<svg/>"), bytes(), bytes(0, 0, 1, 0), ascii("....ftypavif")]) assert.equal(imageSize(other), undefined);
});

/** A page that draws images, each a part of its own. */
function imagePage(images) {
  const ops = images.map((_, i) => cat(bytes(Op.IMAGE), varuint(i), f32(0, 0, 10, 10)));
  const parts = images.map((b, i) => ({ h: name(i + 2), t: "img", bytes: b }));
  return single(fixed([page(name(1))]), [{ h: name(1), t: "obj", bytes: object({ images: parts.map((p) => p.h), ops }) }, ...parts]);
}

test("an image of too many pixels is not decoded, and drawn as nothing", async (t) => {
  t.mock.method(console, "warn", () => {});
  const decoded = [];
  const decodeImage = async (b) => {
    const { width, height } = imageSize(b);
    decoded.push(width);
    return { width, height, close() {} };
  };
  // a PNG of 33 bytes that says it is a gigabyte and a half
  const doc = await open(imagePage([png(20000, 20000), png(640, 480)]));
  const pr = new PageRenderer(doc, {}, fonts, { decodeImage });
  const ctx = context();
  await pr.renderPage(ctx, doc.view("v").pages[0], { scale: 1 });
  assert.deepEqual(decoded, [640]);
  assert.equal(ctx.calls.drawImage, 1);
  assert.equal(pr.res.image(name(2)), undefined);
  assert.equal(pr.res.image(name(3)).width, 640);
  assert.throws(() => pr.res.image(name(4)), /not loaded/);
  assert.equal(console.warn.mock.callCount(), 1);
  // the limit can be set; an image whose header is not read is measured once decoded
  let closed = 0;
  const small = new PageRenderer(doc, {}, fonts, { maxImagePixels: 640 * 480 - 1, decodeImage: async () => ({ width: 640, height: 480, close() { closed++; } }) });
  const avif = await open(imagePage([ascii("....ftypavif")]));
  const other = new PageRenderer(avif, {}, fonts, { maxImagePixels: 640 * 480 - 1, decodeImage: async () => ({ width: 640, height: 480, close() { closed++; } }) });
  await small.renderPage(context(), doc.view("v").pages[0], { scale: 1 });
  assert.equal(closed, 0);
  assert.equal(small.res.imageStats.count, 0);
  await other.renderPage(context(), avif.view("v").pages[0], { scale: 1 });
  assert.equal(closed, 1);
  assert.equal(other.res.image(name(2)), undefined);
});

test("the images a render holds are limited", async () => {
  const decodeImage = async (b) => ({ ...(imageSize(b) ?? { width: 1000, height: 1000 }), close() {} });
  const images = [png(1000, 1000), png(1000, 1000), ascii("....ftypavif")];
  const doc = await open(imagePage(images));
  const p = doc.view("v").pages[0];
  const fits = new PageRenderer(doc, {}, fonts, { holdLimit: 3 * 4e6, decodeImage });
  const ctx = context();
  await fits.renderPage(ctx, p, { scale: 1 });
  assert.equal(ctx.calls.drawImage, 3);
  const more = new PageRenderer(doc, {}, fonts, { holdLimit: 3 * 4e6 - 1, decodeImage });
  await assert.rejects(more.renderPage(context(), p, { scale: 1 }), formatError(/images of a render take more/));
  // what it held is let go
  more.res.release(more.res.hold());
  assert.ok(more.res.imageStats.bytes <= more.res.imageBudget);
  // an image prepared twice for a render counts once
  const hold = fits.res.hold();
  for (let i = 0; i < 5; i++) await fits.res.prepare(p.layers[0].obj, hold);
  assert.equal(hold.bytes, 3 * 4e6);
});

test("the tiles of a region are found by their coordinates or among the tiles", () => {
  const tiles = { "0,0": name(1), "1,0": name(2), "0,1": name(3), "5,7": name(4), "01,1": name(5), "x": name(6), "2,2": "" };
  const view = { id: "s", kind: "sheet", tiles };
  const all = [[0, 0, name(1)], [1, 0, name(2)], [0, 1, name(3)], [5, 7, name(4)]];
  // few coordinates: each is looked up
  assert.deepEqual(tilesIn(view, [{ tx0: 0, ty0: 0, tx1: 1, ty1: 1 }]), all.slice(0, 3));
  assert.deepEqual(tilesIn(view, [{ tx0: -1, ty0: 0, tx1: 0, ty1: 0 }]), all.slice(0, 1));
  // more coordinates than tiles: the tiles are looked through, with the same result
  assert.deepEqual(tilesIn(view, [{ tx0: 0, ty0: 0, tx1: 1, ty1: 2 }]), all.slice(0, 3));
  assert.deepEqual(tilesIn(view, [{ tx0: 0, ty0: 0, tx1: 1e15, ty1: 1e15 }]), all);
  assert.deepEqual(tilesIn(view, [{ tx0: 0, ty0: 0, tx1: Infinity, ty1: Infinity }]), all);
  assert.deepEqual(tilesIn(view, [{ tx0: 5, ty0: 0, tx1: 1e9, ty1: 1e9 }]), all.slice(3));
  // several ranges: each tile once, by row and column
  assert.deepEqual(tilesIn(view, [{ tx0: 0, ty0: 1, tx1: 0, ty1: 1 }, { tx0: 0, ty0: 0, tx1: 1, ty1: 1 }, { tx0: 0, ty0: 0, tx1: 0, ty1: 0 }]), all.slice(0, 3));
  assert.deepEqual(tilesIn(view, [{ tx0: 3, ty0: 3, tx1: 1e9, ty1: 1e9 }, { tx0: 0, ty0: 0, tx1: 0, ty1: 1e9 }]), [all[0], all[2], all[3]]);
  assert.deepEqual(tilesIn(view, []), []);
  assert.deepEqual(tilesIn(view, [{ tx0: 2, ty0: 0, tx1: 1, ty1: 5 }, { tx0: NaN, ty0: 0, tx1: 1, ty1: 1 }]), []);
  assert.deepEqual(tilesIn({ id: "s", kind: "sheet" }, [{ tx0: 0, ty0: 0, tx1: 1e9, ty1: 1e9 }]), []);
});

test("a sheet is drawn whatever its manifest says of its size", async () => {
  const sheet = (v) => single({ views: [{ id: "v", kind: "sheet", tiles: { "0,0": name(1) }, gridlines: true, ...v }] }, [{ h: name(1), t: "obj", bytes: object(text) }]);
  const draw = async (v, viewport) => {
    const doc = await open(sheet(v));
    const ctx = context();
    await new PageRenderer(doc, {}, fonts).renderSheet(ctx, doc.view("v"), viewport, { scale: 1 });
    return ctx.calls;
  };
  const chunk = { x: 0, y: 0, w: 512, h: 512 };
  const valid = await draw({ tile: 2048, cols: [[16384, 64]], rows: [[1048576, 20]] }, chunk);
  assert.equal(valid.fillText, 1);
  // a tile size of 0 is the default one: the same is drawn
  assert.deepEqual(await draw({ tile: 0, cols: [[16384, 64]], rows: [[1048576, 20]] }, chunk), valid);
  assert.equal((await draw({ tile: 0 }, { x: 10, y: 10, w: 512, h: 512 })).fillText, 1);
  // the whole of a sheet as large as a sheet may be
  const whole = await draw({ tile: 1, gridlines: false, cols: [[1 << 24, 1 << 24]], rows: [[1, 1 << 24]] }, { x: 0, y: 0, w: 2 ** 48, h: 2 ** 24 });
  assert.equal(whole.fillText, 1);
  await assert.rejects(draw({ tile: 1e-3 }, chunk), BdfFormatError);
  await assert.rejects(draw({ rows: [[1e11, 0]] }, chunk), BdfFormatError);
});

test("the size of an SVG image is found in a head of many starts of tags", () => {
  const head = ascii("<svg".repeat(16384));
  const t0 = performance.now();
  assert.deepEqual(svgSize(head), { width: 300, height: 150 });
  assert.ok(performance.now() - t0 < 300, `${performance.now() - t0} ms`);
  // the start tag of the root, as before
  const nested = ascii('<?xml version="1.0"?><!-- a picture --><s:svg xmlns:s="http://www.w3.org/2000/svg" width="10" height=\'20\'><s:svg width="1" height="2"/></s:svg>');
  assert.deepEqual(svgSize(nested), { width: 10, height: 20 });
});
