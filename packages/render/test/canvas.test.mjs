import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { BdfDocument, BufferSource } from "@bdf/core";
import { ResourceCache } from "../dist/resources.js";
import { CanvasRenderer } from "../dist/canvas.js";

// The first tile of the fixture's sheet: 2,678 runs of 10 pt text, and a
// picture (a child object used at 136, 1900).
const root = new URL("../../../", import.meta.url);
const doc = await BdfDocument.open(new BufferSource(new Uint8Array(await readFile(new URL("testdata/demo.bdf", root)))));
const tileHash = doc.view("sheet1").tiles["0,0"];

// Node has no fonts to load, nor pictures to decode.
globalThis.FontFace ??= class { load() { return Promise.resolve(this); } };
const decodeImage = async () => ({ width: 10, height: 10, close() {} });

/**
 * A 2D context that records where text is drawn (following translations)
 * and counts measureText calls; measured widths are off, so that runs are
 * corrected to their advance as with fallback fonts.
 */
function recorder() {
  const rec = { runs: [], images: 0, measured: 0 };
  const state = { font: "10px sans-serif", letterSpacing: "0px", direction: "inherit" };
  let at = [0, 0];
  const stack = [];
  rec.ctx = new Proxy(state, {
    get(t, k) {
      if (k in t) return t[k];
      switch (k) {
        case "save": return () => stack.push([...at]);
        case "restore": return () => { at = stack.pop() ?? at; };
        case "translate": return (x, y) => { at = [at[0] + x, at[1] + y]; };
        case "fillText": return (text, x, y) => rec.runs.push({ text, x: at[0] + x, y: at[1] + y });
        case "measureText": return (text) => { rec.measured++; return { width: text.length * 7 }; };
        case "drawImage": return () => rec.images++;
        case "getTransform": return () => ({ a: 1, b: 0, c: 0, d: 1, e: 0, f: 0 });
        default: return () => {};
      }
    },
    set(t, k, v) { t[k] = v; return true; },
  });
  return rec;
}

async function renderer(fontSet) {
  const res = new ResourceCache(doc, fontSet, { decodeImage });
  await res.prepare(tileHash);
  return { r: new CanvasRenderer(res), tile: res.object(tileHash) };
}

test("text outside the visible box is skipped, and everything in it drawn", async () => {
  const { r, tile } = await renderer();
  const all = recorder();
  r.draw(all.ctx, tile);
  assert.equal(all.runs.length, 2678);
  assert.equal(all.images, 1);

  const box = { x: 100, y: 400, w: 300, h: 60 };
  const seen = recorder();
  r.draw(seen.ctx, tile, true, box);
  const reach = 4 * 10; // INK_REACH font sizes of 10 pt text
  assert.ok(seen.runs.length > 0 && seen.runs.length < all.runs.length / 10, `${seen.runs.length} runs drawn`);
  for (const run of all.runs) {
    const inside = run.x >= box.x && run.x <= box.x + box.w && run.y >= box.y && run.y <= box.y + box.h;
    if (inside) assert.ok(seen.runs.some((s) => s.text === run.text && s.x === run.x && s.y === run.y), `${run.text} at ${run.x}, ${run.y} is not drawn`);
  }
  for (const run of seen.runs) assert.ok(run.y >= box.y - reach && run.y <= box.y + box.h + reach, `${run.text} at ${run.y} is drawn`);
  // what the child object draws is not text: it is drawn whatever the box
  assert.equal(seen.images, 1);
});

/** An object of text drawn by ops: [name, ...args], strings for FILL_TEXT and FILTER given as they are. */
function textObject(ops) {
  const bytes = [], strings = [];
  const f32 = (v) => bytes.push(...new Uint8Array(new Float32Array([v]).buffer));
  const u32 = (v) => bytes.push(v & 0xff, (v >>> 8) & 0xff, (v >>> 16) & 0xff, (v >>> 24) & 0xff);
  const str = (s) => { strings.push(s); bytes.push(strings.length - 1); };
  const code = { SAVE: 0x01, RESTORE: 0x02, TRANSFORM: 0x03, TRANSLATE: 0x04, SCALE: 0x05, SHADOW: 0x18, FILTER: 0x19, FONT: 0x1a, FILL_TEXT: 0x30 };
  for (const [name, ...a] of ops) {
    bytes.push(code[name]);
    switch (name) {
      case "FONT": bytes.push(a[0]); f32(a[1]); break;
      case "FILL_TEXT": str(a[0]); a.slice(1).forEach(f32); break;
      case "SHADOW": u32(a[0]); a.slice(1).forEach(f32); break;
      case "FILTER": str(a[0]); break;
      default: a.forEach(f32);
    }
  }
  return {
    opset: 1, bbox: { x: 0, y: 0, w: 10000, h: 10000 }, strings, paths: [], paints: [], images: [], objects: [],
    fonts: [{ kind: 1, family: "sans-serif", weight: 400, style: 0 }], ops: new Uint8Array(bytes),
  };
}

test("the box follows transforms, and shadows, filters and unknown states draw everything", async () => {
  const { r } = await renderer();
  const obj = textObject([
    ["FONT", 0, 10],
    // all at y ≈ 1000 once transformed
    ["SAVE"], ["TRANSLATE", 0, 1000], ["FILL_TEXT", "moved", 0, 5, 30], ["RESTORE"],
    ["SAVE"], ["SCALE", 2, 2], ["FILL_TEXT", "scaled", 0, 502, 30], ["RESTORE"],
    ["SAVE"], ["TRANSFORM", 0, 1, -1, 0, 0, 0], ["FILL_TEXT", "rotated", 1005, 0, 30], ["RESTORE"],
    ["FILL_TEXT", "far", 0, 5000, 30],
    ["SAVE"], ["SHADOW", 0x000000ff, 4, 2, 2], ["FILL_TEXT", "shadowed", 0, 5000, 30], ["RESTORE"],
    ["SAVE"], ["FILTER", "blur(2px)"], ["FILL_TEXT", "filtered", 0, 5000, 30], ["RESTORE"],
    ["FILL_TEXT", "far again", 0, 5000, 30],
    // a RESTORE with no SAVE: what the state is now is not known
    ["RESTORE"], ["FILL_TEXT", "unknown", 0, 5000, 30],
  ]);
  const rec = recorder();
  r.draw(rec.ctx, obj, true, { x: 0, y: 990, w: 100, h: 30 });
  assert.deepEqual(rec.runs.map((run) => run.text), ["moved", "scaled", "rotated", "shadowed", "filtered", "unknown"]);
});

test("widths are measured once per font and text", async () => {
  const { r, tile } = await renderer();
  const first = recorder();
  r.draw(first.ctx, tile);
  assert.ok(first.measured > 0);
  const again = recorder();
  r.draw(again.ctx, tile);
  assert.equal(again.measured, 0);
  assert.deepEqual(again.runs, first.runs);
});

test("widths in a font still loading are measured each time", async () => {
  const loading = { check: () => false, add() {}, delete() {} };
  const { r, tile } = await renderer(loading);
  const first = recorder(), again = recorder();
  r.draw(first.ctx, tile);
  r.draw(again.ctx, tile);
  assert.ok(first.measured > 0);
  assert.equal(again.measured, first.measured);
});
