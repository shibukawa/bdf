import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { BdfDocument, BufferSource } from "@bdfkit/core";
import { SvgExporter, sheetBounds, saveTarget, imageScale, svgPathData } from "../dist/svgexport.js";

const root = new URL("../../../", import.meta.url);
const open = async (path) => BdfDocument.open(new BufferSource(new Uint8Array(await readFile(new URL(path, root)))));

/**
 * Check that an SVG stands alone and holds together: its groups are closed,
 * what it refers to by id is defined in it once, and it refers to no file.
 */
function checkSvg(svg) {
  assert.match(svg, /^<svg xmlns="http:\/\/www\.w3\.org\/2000\/svg" viewBox="[-\d. ]+" width="[\d.]+pt" height="[\d.]+pt"/);
  assert.ok(svg.endsWith("</svg>\n"));
  for (const tag of ["g", "text", "defs", "mask", "clipPath", "symbol", "pattern", "filter", "a"]) {
    const opened = (svg.match(new RegExp(`<${tag}[ >]`, "g")) ?? []).length;
    const closed = (svg.match(new RegExp(`</${tag}>`, "g")) ?? []).length;
    assert.equal(opened, closed, `<${tag}> opened ${opened} times, closed ${closed}`);
  }
  const ids = [...svg.matchAll(/ id="([^"]+)"/g)].map((m) => m[1]);
  assert.equal(new Set(ids).size, ids.length, "an id is defined twice");
  for (const m of svg.matchAll(/url\(#([^)]+)\)|href="#([^"]+)"/g)) assert.ok(ids.includes(m[1] ?? m[2]), `id ${m[1] ?? m[2]} is not defined`);
  for (const m of svg.matchAll(/(?:href="|url\(")([^"#]+)"/g)) assert.match(m[1], /^(data:|https?:\/\/|mailto:)/, `${m[1].slice(0, 40)} is a file beside the SVG`);
}

test("the path data of Canvas paths", () => {
  const path = (verbs, args) => ({ verbs: new Uint8Array(verbs), args: new Float32Array(args) });
  // a line without a sub-path starts one
  assert.equal(svgPathData(path([1, 1], [1, 2, 3, 4]), 2), "M1 2L3 4");
  assert.equal(svgPathData(path([5, 1], [1, 2, 3, 4, 9, 9]), 2), "M1 2h3v4h-3ZL9 9");
  assert.equal(svgPathData(path([3], [1, 2, 3, 4, 5, 6]), 2), "M1 2C1 2 3 4 5 6");
  assert.equal(svgPathData(path([0, 1], [0, 0, -1.5, -0.25]), 2), "M0 0L-1.5-.25");
  // a full circle: three arcs, each one reading
  assert.equal(svgPathData(path([6], [0, 0, 10, 10, 0, 0, 6.2831855, 0]), 2), "M10 0A10 10 0 0 1 -5 8.66A10 10 0 0 1 -5 -8.66A10 10 0 0 1 10 0");
  assert.equal(svgPathData(path([8], [0, 0, 10, 10, 2]), 2), "M2 0L8 0L8 0A2 2 0 0 1 10 2L10 8L10 8A2 2 0 0 1 8 10L2 10L2 10A2 2 0 0 1 0 8L0 2L0 2A2 2 0 0 1 2 0Z");
  assert.equal(svgPathData(path([], []), 2), undefined);
  assert.equal(svgPathData(path([0, 1], [0, 0, NaN, 1]), 2), undefined);
});

test("the pages of the fixture stand alone", async () => {
  const doc = await open("testdata/demo.bdf");
  const out = new SvgExporter(doc);
  const pages = doc.manifest.views.find((v) => v.kind === "fixed").pages;
  let texts = 0, fonts = 0;
  for (const page of pages) {
    const { svg, warnings } = await out.page(page);
    checkSvg(svg);
    assert.deepEqual(warnings, []);
    assert.ok(svg.includes(`viewBox="0 0 ${page.w} ${page.h}"`));
    texts += (svg.match(/<text/g) ?? []).length;
    fonts += (svg.match(/@font-face\{font-family:"bdf-[0-9a-f]{10}";src:url\("data:font\/\w+;base64,/g) ?? []).length;
  }
  // text stays text, in the fonts of the document
  assert.ok(texts > 10, `${texts} text elements`);
  assert.ok(fonts > 0, "no embedded font");
  // a part of a page: the view box is the part
  const { svg } = await out.page(pages[0], { x: 10, y: 20, w: 100, h: 50 });
  checkSvg(svg);
  assert.ok(svg.includes(`viewBox="10 20 100 50" width="100pt" height="50pt"`));
  // no background
  const bare = await new SvgExporter(doc, { background: null }).page(pages[0]);
  assert.ok(!bare.svg.includes(`<path d="M0 0h${pages[0].w}v${pages[0].h}h-${pages[0].w}Z" fill="#fff"/>`));
  assert.ok(svg.includes(`<path d="M0 0h${pages[0].w}v${pages[0].h}h-${pages[0].w}Z" fill="#fff"/>`));
});

test("a sheet is saved as far as its cells reach", async () => {
  const doc = await open("testdata/demo.bdf");
  const sheet = doc.view("sheet1");
  const bounds = sheetBounds(sheet);
  assert.ok(bounds.w > 0 && bounds.h > 0);
  assert.deepEqual(saveTarget(sheet, 3), { target: bounds, width: bounds.w, height: bounds.h });
  const { svg } = await new SvgExporter(doc).sheet(sheet, bounds);
  checkSvg(svg);
  // the gridlines are one path, and each tile is cut to its square
  assert.equal((svg.match(/vector-effect="non-scaling-stroke"/g) ?? []).length, sheet.gridlines ? 1 : 0);
  assert.equal((svg.match(/<clipPath/g) ?? []).length, Object.keys(sheet.tiles).length);
  assert.ok((svg.match(/<text/g) ?? []).length > 1000);
});

test("a scroll view is saved whole, its strips stacked", async () => {
  const doc = await open("testdata/html/article.bdf");
  const view = doc.manifest.views.find((v) => v.kind === "scroll");
  const what = saveTarget(view, 0);
  const height = view.pages.reduce((h, p) => h + (p.body ?? p).h, 0);
  assert.deepEqual(what.target, { x: 0, y: 0, w: what.width, h: height });
  const { svg } = await new SvgExporter(doc).continuous(view, what.target);
  checkSvg(svg);
  assert.equal((svg.match(/<clipPath/g) ?? []).length, view.pages.length);
  // the links that leave the file stay, those to its pages do not
  assert.ok(/<a href="https?:\/\//.test(svg));
  assert.ok(!/<a href="(?!https?:|mailto:)/.test(svg));
  // a part of the layout draws the strips it shows only
  const part = await new SvgExporter(doc).continuous(view, { x: 0, y: 0, w: what.width, h: 10 });
  checkSvg(part.svg);
  assert.ok(part.svg.length < svg.length);
});

test("what a view saves, and how large an image of it may be", async () => {
  const doc = await open("testdata/demo.bdf");
  const fixed = doc.manifest.views.find((v) => v.kind === "fixed");
  assert.deepEqual(saveTarget(fixed, 1), { target: 1, width: fixed.pages[1].w, height: fixed.pages[1].h });
  assert.equal(saveTarget(fixed, 99), undefined);
  assert.equal(saveTarget({ id: "empty", kind: "sheet" }, 0), undefined);
  assert.equal(imageScale(600, 800, 2), 2);
  // a long strip is drawn smaller than asked for: no canvas holds it
  assert.equal(imageScale(600, 100000, 2), 16384 / 100000);
  assert.ok(imageScale(20000, 20000, 2) * 20000 <= 16384);
});
