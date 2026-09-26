import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { BdfDocument, BufferSource } from "@bdf/core";
import { ResourceCache } from "../dist/resources.js";
import { PageRenderer } from "../dist/page.js";

// testdata/tiff/orientation.bdf has eight pages, each one image of its own.
const root = new URL("../../../", import.meta.url);
const bytes = new Uint8Array(await readFile(new URL("testdata/tiff/orientation.bdf", root)));
const open = () => BdfDocument.open(new BufferSource(bytes));

const SIZE = 100 * 100 * 4; // what a fake bitmap counts for

/**
 * A decoder of fake 100×100 bitmaps that counts decodings and closings.
 * With manual, each decoding waits until release(n) lets the n-th finish.
 */
function fakeDecoder({ manual = false } = {}) {
  const d = { decoded: 0, closed: 0, waiting: [] };
  d.decode = () => {
    d.decoded++;
    const bmp = { width: 100, height: 100, closed: false, close() { this.closed = true; d.closed++; } };
    if (!manual) return Promise.resolve(bmp);
    return new Promise((resolve) => d.waiting.push(() => resolve(bmp)));
  };
  d.release = async (n) => {
    d.waiting[n]();
    for (let i = 0; i < 10; i++) await new Promise((r) => setTimeout(r)); // let the awaiting renders run
  };
  return d;
}

/** The image hash each page draws. */
async function pageImages(doc) {
  return Promise.all(doc.view("pages").pages.map(async (p) => (await doc.object(p.layers[0].obj)).images[0]));
}

/** A 2D context that does nothing but record the images drawn, and whether they were closed then. */
function fakeContext(drawn) {
  const state = {};
  return new Proxy(state, {
    get(t, k) {
      if (k in t) return t[k];
      if (k === "drawImage") return (img) => drawn.push({ img, closed: img.closed });
      if (k === "getTransform") return () => ({ a: 1, b: 0, c: 0, d: 1, e: 0, f: 0 });
      if (k === "measureText") return () => ({ width: 0 });
      return () => {};
    },
    set(t, k, v) { t[k] = v; return true; },
  });
}

test("images beyond the budget are closed, least recently prepared first", async () => {
  const doc = await open();
  const dec = fakeDecoder();
  const res = new ResourceCache(doc, undefined, { imageBudget: 2 * SIZE, decodeImage: dec.decode });
  const pages = doc.view("pages").pages;
  for (const p of pages) {
    const hold = res.hold();
    await res.prepare(p.layers[0].obj, hold);
    res.release(hold);
  }
  const images = await pageImages(doc);
  assert.equal(dec.decoded, 8);
  assert.deepEqual(res.imageStats, { count: 2, bytes: 2 * SIZE });
  assert.equal(dec.closed, 6);
  // The last two pages' images are kept, the others are gone.
  assert.equal(res.image(images[7]).closed, false);
  assert.equal(res.image(images[6]).closed, false);
  assert.throws(() => res.image(images[0]), /not loaded/);

  // Preparing page 7 again makes it the most recent: page 8 goes next.
  const hold = res.hold();
  await res.prepare(pages[6].layers[0].obj, hold);
  await res.prepare(pages[0].layers[0].obj, hold); // decoded again
  res.release(hold);
  assert.equal(dec.decoded, 9);
  assert.equal(res.image(images[0]).closed, false);
  assert.equal(res.image(images[6]).closed, false);
  assert.throws(() => res.image(images[7]), /not loaded/);
});

test("images a render holds are not closed, even over the budget", async () => {
  const doc = await open();
  const dec = fakeDecoder();
  // A budget smaller than one image: nothing stays that is not held.
  const res = new ResourceCache(doc, undefined, { imageBudget: SIZE / 2, decodeImage: dec.decode });
  const pages = doc.view("pages").pages;
  const first = res.hold();
  await res.prepare(pages[0].layers[0].obj, first);
  for (const p of pages.slice(1)) {
    const hold = res.hold();
    await res.prepare(p.layers[0].obj, hold);
    res.release(hold);
  }
  const images = await pageImages(doc);
  assert.equal(res.image(images[0]).closed, false, "held image closed");
  assert.deepEqual(res.imageStats, { count: 1, bytes: SIZE });
  assert.equal(dec.closed, 7);
  res.release(first);
  assert.deepEqual(res.imageStats, { count: 0, bytes: 0 });
  assert.throws(() => res.image(images[0]), /not loaded/);
});

test("a render that has prepared keeps its images while another render finishes", async () => {
  // The race holds prevent: render A has decoded page 1's image and waits
  // for page 2's; render B draws page 3 and finishes, which trims the cache
  // to its budget of one image. Without A's hold, page 1's image (the least
  // recently prepared) would be closed before A draws it.
  const doc = await open();
  const dec = fakeDecoder({ manual: true });
  const pr = new PageRenderer(doc, {}, undefined, { imageBudget: SIZE, decodeImage: dec.decode });
  const view = doc.view("pages");
  const images = await pageImages(doc);
  const drawnA = [], drawnB = [];
  const opts = { scale: 1 };
  const a = pr.renderContinuous(fakeContext(drawnA), view, { x: 0, y: 0, w: 96, h: 100 }, opts); // pages 1 and 2 (64 pt each)
  await new Promise((r) => setTimeout(r, 10));
  assert.equal(dec.waiting.length, 2);
  await dec.release(0); // page 1 decoded
  const b = pr.renderPage(fakeContext(drawnB), view.pages[2], opts);
  await new Promise((r) => setTimeout(r, 10));
  await dec.release(2); // page 3 decoded
  await b;
  assert.equal(drawnB.length, 1);
  await dec.release(1); // page 2 decoded
  await a;
  assert.deepEqual(drawnA.map((d) => d.closed), [false, false]);
  // Both renders are over: only the most recent image is left.
  assert.deepEqual(pr.res.imageStats, { count: 1, bytes: SIZE });
  assert.ok(pr.res.image(images[1]));
});

test("renders that need the same image share one decoding", async () => {
  const doc = await open();
  const dec = fakeDecoder({ manual: true });
  const pr = new PageRenderer(doc, {}, undefined, { decodeImage: dec.decode });
  const page = doc.view("pages").pages[0];
  const drawn = [];
  const r1 = pr.renderPage(fakeContext(drawn), page, { scale: 1 });
  const r2 = pr.renderPage(fakeContext(drawn), page, { scale: 2 });
  await new Promise((r) => setTimeout(r, 10));
  assert.equal(dec.waiting.length, 1);
  await dec.release(0);
  await Promise.all([r1, r2]);
  assert.equal(drawn.length, 2);
  assert.equal(drawn[0].img, drawn[1].img);
  assert.equal(dec.decoded, 1);
});

test("preparing for text does not decode images", async () => {
  const doc = await open();
  const dec = fakeDecoder();
  const pr = new PageRenderer(doc, {}, undefined, { decodeImage: dec.decode });
  for (const p of doc.view("pages").pages) await pr.preparePageText(p);
  assert.equal(dec.decoded, 0);
  assert.equal(pr.res.imageStats.count, 0);
});

test("dispose closes the images, and those decoded afterwards", async () => {
  const doc = await open();
  const dec = fakeDecoder({ manual: true });
  const res = new ResourceCache(doc, undefined, { decodeImage: dec.decode });
  const pages = doc.view("pages").pages;
  const p1 = res.prepare(pages[0].layers[0].obj);
  await new Promise((r) => setTimeout(r, 10));
  await dec.release(0);
  await p1;
  const p2 = res.prepare(pages[1].layers[0].obj);
  await new Promise((r) => setTimeout(r, 10));
  res.dispose();
  assert.equal(dec.closed, 1);
  await dec.release(1);
  await p2;
  assert.equal(dec.closed, 2);
  assert.deepEqual(res.imageStats, { count: 0, bytes: 0 });
});

test("dispose leaves the images of a render in progress until it releases them", async () => {
  // The worker disposes the cache of a document it closes; a render of
  // that document may still be between preparing and drawing.
  const doc = await open();
  const dec = fakeDecoder();
  const res = new ResourceCache(doc, undefined, { decodeImage: dec.decode });
  const pages = doc.view("pages").pages;
  const images = await pageImages(doc);
  const hold = res.hold();
  await res.prepare(pages[0].layers[0].obj, hold);
  await res.prepare(pages[1].layers[0].obj);
  res.dispose();
  assert.equal(dec.closed, 1);
  assert.equal(res.image(images[0]).closed, false);
  res.release(hold);
  assert.equal(dec.closed, 2);
  assert.deepEqual(res.imageStats, { count: 0, bytes: 0 });
});
