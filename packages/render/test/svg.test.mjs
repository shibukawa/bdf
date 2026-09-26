import { test } from "node:test";
import assert from "node:assert/strict";
import { isSvg, svgSize, VectorImage } from "../dist/svg.js";

const enc = (s) => new TextEncoder().encode(s);

test("isSvg tells markup from encoded bitmaps", () => {
  assert.ok(isSvg(enc("<svg/>")));
  assert.ok(isSvg(enc("﻿ \n<?xml version='1.0'?><svg/>")));
  assert.ok(isSvg(new Uint8Array([0xff, 0xfe, 0x3c, 0])));
  assert.ok(!isSvg(new Uint8Array([0x89, 0x50, 0x4e, 0x47])));
  assert.ok(!isSvg(new Uint8Array([0xff, 0xd8, 0xff])));
  assert.ok(!isSvg(new Uint8Array([])));
});

test("svgSize follows the root's width, height and view box (spec §6.2)", () => {
  const size = (attrs) => svgSize(enc(`<?xml version="1.0"?>\n<!-- c -->\n<svg xmlns="http://www.w3.org/2000/svg" ${attrs}><rect/></svg>`));
  assert.deepEqual(size(`width="200" height="100"`), { width: 200, height: 100 });
  assert.deepEqual(size(`width='72pt' height="1in"`), { width: 96, height: 96 });
  assert.deepEqual(size(`viewBox="0 0 40 30"`), { width: 40, height: 30 });
  assert.deepEqual(size(`width="80" viewBox="0,0,40,30"`), { width: 80, height: 60 });
  assert.deepEqual(size(`width="100%" height="50" viewBox="0 0 40 20"`), { width: 100, height: 50 });
  assert.deepEqual(size(``), { width: 300, height: 150 });
  assert.deepEqual(size(`width="50"`), { width: 50, height: 150 });
  assert.deepEqual(svgSize(enc(`<svg:svg xmlns:svg="http://www.w3.org/2000/svg" width="10mm" height="5em"/>`)), { width: 10 * (96 / 25.4), height: 80 });
});

/** A rasterizer that records its calls and makes fake bitmaps. */
function fakeRasterizer(fail = false) {
  const calls = [];
  const rasterize = async (hash, data, width, height) => {
    calls.push([width, height]);
    if (fail) throw new Error("cannot draw");
    return { width, height, closed: false, close() { this.closed = true; } };
  };
  return { calls, rasterize };
}

test("VectorImage draws a raster once per scale, and reuses a slightly larger one", async () => {
  const v = new VectorImage("h", enc(`<svg xmlns="http://www.w3.org/2000/svg" width="100" height="50"/>`));
  const { calls, rasterize } = fakeRasterizer();
  let r = v.raster(2);
  assert.equal(r.raster, undefined);
  assert.equal(r.missing, 2);
  // two draws missing the same raster draw it once
  assert.deepEqual(await Promise.all([v.draw(2, rasterize), v.draw(2, rasterize)]), [true, true]);
  assert.deepEqual(calls, [[200, 100]]);
  r = v.raster(2);
  assert.equal(r.raster.bitmap.width, 200);
  assert.equal(r.missing, undefined);
  // within 25 % larger: used as it is
  assert.equal(v.raster(1.7).raster.k, 2);
  // smaller than needed: the closest meanwhile, and the scale missed
  r = v.raster(3);
  assert.equal(r.raster.k, 2);
  assert.equal(r.missing, 3);
});

test("VectorImage keeps the rasters used last and bounds their size", async () => {
  const v = new VectorImage("h", enc(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1000 1000"/>`));
  const { calls, rasterize } = fakeRasterizer();
  const first = [];
  for (const k of [0.5, 1, 2, 3, 4]) {
    await v.draw(k, rasterize);
    first.push(v.raster(k).raster?.bitmap);
  }
  assert.ok(first[0].closed, "the least recently used raster is closed");
  assert.equal(v.raster(0.5).missing, 0.5);
  // 4096 × 4096 pixels at most
  assert.equal(v.raster(100).missing, 4.096);
  assert.deepEqual(calls.at(-1), [4000, 4000]);
});

test("VectorImage gives up on a scale it cannot draw", async () => {
  const v = new VectorImage("h", enc(`<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10"/>`));
  const { rasterize } = fakeRasterizer(true);
  const warn = console.warn;
  console.warn = () => {};
  try {
    assert.equal(await v.draw(1, rasterize), false);
    assert.equal(await v.draw(2, undefined), false);
  } finally {
    console.warn = warn;
  }
  assert.deepEqual(v.raster(1), { raster: undefined, missing: undefined });
});
