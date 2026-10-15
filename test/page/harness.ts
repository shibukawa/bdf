// Browser integration harness: rendering, text selection, accessibility and search.
// Driven by test/browser.mjs.
import { BdfDocument, fetchSingle, RangeSource, SplitSource, type Rect } from "@bdfkit/core";
import { PageRenderer, BdfWorkerClient, buildTextLayer, selectionText, joinRuns, TEXT_LAYER_CSS, RUN_ATTR } from "@bdfkit/render";

interface Case {
  name: string;
  /** Document path (default: the generated fixture). */
  src?: string;
  kind: "page" | "continuous" | "sheet";
  view: string;
  page?: number;
  viewport?: Rect;
  scale: number;
  expectedSize: [number, number];
  roles?: string[];
}

// Exercise each rendering mode without storing screenshots or matching past pixels.
const CASES: Case[] = [
  { name: "page", kind: "page", view: "slides", page: 0, scale: 1, expectedSize: [960, 540] },
  { name: "page-body-only", kind: "page", view: "slides", page: 0, scale: 0.5, expectedSize: [480, 270], roles: ["body"] },
  { name: "page-no-layers", kind: "page", view: "slides", page: 0, scale: 0.5, expectedSize: [480, 270], roles: [] },
  { name: "continuous", kind: "continuous", view: "doc", viewport: { x: 0, y: 600, w: 451.3, h: 400 }, scale: 1, expectedSize: [452, 400] },
  { name: "sheet-tile-boundary", kind: "sheet", view: "sheet1", viewport: { x: 0, y: 1800, w: 800, h: 500 }, scale: 1, expectedSize: [800, 500] },
  { name: "sheet-zoomed", kind: "sheet", view: "sheet1", viewport: { x: 64, y: 20, w: 300, h: 150 }, scale: 2, expectedSize: [600, 300] },
];

const DEFAULT_SRC = "/testdata/demo.bdf";

function canvasSize(c: Case, doc: BdfDocument): [number, number] {
  if (c.kind === "page") {
    const p = doc.view(c.view).pages![c.page!];
    return [Math.ceil(p.w * c.scale), Math.ceil(p.h * c.scale)];
  }
  return [Math.ceil(c.viewport!.w * c.scale), Math.ceil(c.viewport!.h * c.scale)];
}

async function renderMain(doc: BdfDocument, pr: PageRenderer, c: Case): Promise<HTMLCanvasElement> {
  const [w, h] = canvasSize(c, doc);
  const canvas = document.createElement("canvas");
  canvas.width = w;
  canvas.height = h;
  const ctx = canvas.getContext("2d")!;
  const view = doc.view(c.view);
  switch (c.kind) {
    case "page": await pr.renderPage(ctx, view.pages![c.page!], { scale: c.scale, roles: c.roles }); break;
    case "continuous": await pr.renderContinuous(ctx, view, c.viewport!, { scale: c.scale }); break;
    case "sheet": await pr.renderSheet(ctx, view, c.viewport!, { scale: c.scale }); break;
  }
  return canvas;
}

async function renderWorker(client: BdfWorkerClient, c: Case): Promise<ImageBitmap> {
  switch (c.kind) {
    case "page": return client.page(c.view, c.page!, c.scale, c.roles);
    case "continuous": return client.continuous(c.view, c.viewport!, c.scale);
    case "sheet": return client.sheet(c.view, c.viewport!, c.scale);
  }
}

function toCanvas(bmp: ImageBitmap): HTMLCanvasElement {
  const canvas = document.createElement("canvas");
  canvas.width = bmp.width;
  canvas.height = bmp.height;
  canvas.getContext("2d")!.drawImage(bmp, 0, 0);
  return canvas;
}

interface Rendered {
  width: number;
  height: number;
  hasContent: boolean;
}

function rendered(canvas: HTMLCanvasElement): Rendered {
  const pixels = canvas.getContext("2d")!.getImageData(0, 0, canvas.width, canvas.height).data;
  let hasContent = false;
  for (let i = 0; i < pixels.length; i += 4) {
    if (pixels[i + 3] > 0 && Math.min(pixels[i], pixels[i + 1], pixels[i + 2]) < 240) {
      hasContent = true;
      break;
    }
  }
  return { width: canvas.width, height: canvas.height, hasContent };
}

interface Result {
  name: string;
  size: [number, number];
  expectsContent: boolean;
  main: Rendered;
  worker: Rendered;
}

interface Opened { doc: BdfDocument; pr: PageRenderer; client: BdfWorkerClient }

async function main() {
  const params = new URLSearchParams(location.search);
  const sourceKind = params.get("source") ?? "single";
  const workerUrl = params.get("worker") ?? "./worker.js";
  const opened = new Map<string, Opened>();
  const open = async (src: string): Promise<Opened> => {
    let o = opened.get(src);
    if (o) return o;
    const isDefault = src === DEFAULT_SRC;
    const doc = await BdfDocument.open(
      isDefault && sourceKind === "split" ? new SplitSource("/testdata/demo-split/")
        : isDefault && sourceKind === "range" ? new RangeSource(src)
          : await fetchSingle(src),
    );
    const pr = new PageRenderer(doc, {}, document.fonts);
    const client = new BdfWorkerClient(new Worker(workerUrl, { type: "module" }));
    await client.open(isDefault && sourceKind === "split" ? { kind: "split", base: location.origin + "/testdata/demo-split/" } : { kind: "single", url: location.origin + src, range: isDefault && sourceKind === "range" });
    o = { doc, pr, client };
    opened.set(src, o);
    return o;
  };
  const { client } = await open(DEFAULT_SRC);

  const results: Result[] = [];
  for (const c of CASES) {
    const { doc, pr, client } = await open(c.src ?? DEFAULT_SRC);
    const main = await renderMain(doc, pr, c);
    const bitmap = await renderWorker(client, c);
    const worker = toCanvas(bitmap);
    bitmap.close();
    results.push({
      name: c.name,
      size: c.expectedSize,
      expectsContent: c.roles?.length !== 0,
      main: rendered(main),
      worker: rendered(worker),
    });
  }
  // text layer: select runs of the first doc page through the DOM and
  // rebuild the text; compare with the runs joined directly.
  (window as unknown as { bdfSelection: unknown }).bdfSelection = await selectionCheck(client);
  // ALT_TEXT runs (ligatures drawn as one private-use glyph) are in the text layer
  const { client: chromeDoc } = await open("/testdata/pdf/chrome-doc.bdf");
  (window as unknown as { bdfAltText: unknown }).bdfAltText = await altTextCheck(chromeDoc);
  (window as unknown as { bdfForcedColors: unknown }).bdfForcedColors = forcedColorsCheck;
  const { client: xlsx } = await open("/testdata/xlsx/features.bdf");
  (window as unknown as { bdfStructure: unknown }).bdfStructure = await structureCheck(client, xlsx);
  // search through the worker: hits, then rectangles
  const hits = await client.search("doc", "list of objects");
  const rects = await client.locate("doc", hits);
  const sheetHits = await client.search("sheet1", "row 150");
  const sheetRects = await client.locate("sheet1", sheetHits);
  // Japanese text wrapped between characters: the hit spans two lines.
  const { client: pptx } = await open("/testdata/pptx/basic.bdf");
  const pptxHits = await pptx.search("slides", "改行します");
  const pptxRects = await pptx.locate("slides", pptxHits);
  // "fixture" is drawn with an fi ligature: the hit is on an ALT_TEXT run.
  const ligHits = await chromeDoc.search("pages", "fixture");
  const ligRects = await chromeDoc.locate("pages", ligHits);
  // text over a tile edge of a sheet is drawn in both tiles but found once
  const { client: xlsxBasic } = await open("/testdata/xlsx/basic.bdf");
  const xlsxHits = await xlsxBasic.search("sheet3", "straddles the tile");
  const xlsxRects = await xlsxBasic.locate("sheet3", xlsxHits);
  // draw.io: a page of the diagram is a view; the search runs in it, and
  // links to other pages are #view= links
  const { client: drawio } = await open("/testdata/drawio/multipage.bdf");
  const drawioHits = await drawio.search("details", "日本語の説明");
  const drawioRects = await drawio.locate("details", drawioHits);
  const drawioLinks = (await drawio.content("overview", 0)).links.map((l) => l.url);
  // formulas are found by their linear notation, typed with a hyphen-minus
  const { client: math } = await open("/testdata/docx/math.bdf");
  const mathHits = await math.search("pages", "b^2-4ac");
  const mathRects = await math.locate("pages", mathHits);
  (window as unknown as { bdfSearch: unknown }).bdfSearch = { hits, rects, sheetHits, sheetRects, pptxHits, pptxRects, ligHits, ligRects, xlsxHits, xlsxRects, drawioHits, drawioRects, drawioLinks, mathHits, mathRects };
  (window as unknown as { bdfResults: Result[] }).bdfResults = results;
  for (const { pr, client } of opened.values()) {
    pr.dispose();
    await client.close();
    client.terminate();
  }
  document.title = "done";
}

async function selectionCheck(client: BdfWorkerClient) {
  const style = document.createElement("style");
  style.textContent = TEXT_LAYER_CSS;
  document.head.appendChild(style);
  // the structured layer: runs nested in paragraphs, headings and table cells
  const content = await client.content("doc", 0);
  const runs = content.runs;
  const layer = buildTextLayer(content, 1);
  const host = document.createElement("div");
  host.style.cssText = "position: relative; width: 600px; height: 850px;";
  host.appendChild(layer);
  document.body.appendChild(host);
  const spans = layer.querySelectorAll("span");
  const sel = getSelection()!;
  // all runs
  sel.removeAllRanges();
  const all = document.createRange();
  all.selectNodeContents(layer);
  sel.addRange(all);
  const allText = selectionText(sel, host);
  const drawn = runs.filter((r) => r.text).map((r) => ({ ordinal: r.ordinal, sep: r.sep, text: r.text, layer }));
  const wantAll = joinRuns(drawn);
  // a partial range: from the third character of the second run to the
  // second character of the fourth run
  sel.removeAllRanges();
  const part = document.createRange();
  part.setStart(spans[1].firstChild!, 2);
  part.setEnd(spans[3].firstChild!, 2);
  sel.addRange(part);
  const partText = selectionText(sel, host);
  const wantPart = joinRuns([{ ...drawn[1], text: drawn[1].text.slice(2) }, drawn[2], { ...drawn[3], text: drawn[3].text.slice(0, 2) }]);
  sel.removeAllRanges();
  host.remove();
  return { spans: spans.length, allText, wantAll, partText, wantPart, breaks: (allText.match(/\n/g) ?? []).length };
}

async function altTextCheck(client: BdfWorkerClient) {
  const runs = await client.text("pages", 0);
  const layer = buildTextLayer(runs, 1);
  const host = document.createElement("div");
  host.style.cssText = "position: relative; width: 800px; height: 1100px;";
  host.appendChild(layer);
  document.body.appendChild(host);
  const spans = [...layer.querySelectorAll("span")];
  // each ALT_TEXT run has its span, as wide as the ligature run it stands for
  const alt = runs.filter((r) => r.text && r.altText).map((r) => {
    const span = spans.find((s) => s.getAttribute(RUN_ATTR.ordinal) === String(r.ordinal));
    return { text: r.text, span: span?.textContent ?? null, width: span?.getBoundingClientRect().width ?? 0, want: r.advance * r.matrix[0] };
  });
  host.remove();
  return { runs: runs.filter((r) => r.text).length, spans: spans.length, alt };
}

/** Roles, links and languages of structured text layers (spec §7.8). */
async function structureCheck(client: BdfWorkerClient, xlsx: BdfWorkerClient) {
  const roles = (el: Element) => {
    const out: Record<string, number> = {};
    for (const e of el.querySelectorAll("[role]")) out[e.getAttribute("role")!] = (out[e.getAttribute("role")!] ?? 0) + 1;
    return out;
  };
  const slide = buildTextLayer(await client.content("slides", 2), 1, { lang: "en" });
  const flow = buildTextLayer(await client.content("doc", 0), 1);
  const figures = buildTextLayer(await client.content("slides", 1), 2);
  const fig = figures.querySelectorAll<HTMLElement>("[role=img]")[1];
  const sheet = buildTextLayer(await client.sheetContent("sheet1", { x: 0, y: 0, w: 800, h: 500 }), 1, { sheet: { rows: 200, cols: 26, headerRows: 1, headerCols: 1 } });
  // a converted workbook: the header row of a table, a picture and a shape with descriptions
  let placed = 0;
  const workbook = buildTextLayer(await xlsx.sheetContent("sheet1", [{ x: 0, y: 0, w: 400, h: 700 }, { x: 400, y: 0, w: 400, h: 700 }]), 1, { sheet: { rows: 41, cols: 14 }, onSpan: () => placed++ });
  // links a document must not turn into anchors, and a run in another language
  const run = (text: string, x: number, lang?: string) => ({ text, x, y: 20, advance: 40, size: 10, font: undefined, align: 0, matrix: [1, 0, 0, 1, 0, 0] as [number, number, number, number, number, number], sep: 1, ordinal: 0, altText: false, node: 0, lang });
  const unsafe = buildTextLayer({
    runs: [run("safe", 0), run("日本語", 100, "ja"), run("evil", 200)],
    nodes: [{ kind: "paragraph", parent: -1 }],
    links: [
      { x: 0, y: 0, w: 60, h: 30, url: "https://example.com/", after: 2, node: 0 },
      { x: 200, y: 0, w: 60, h: 30, url: "javascript:alert(1)", after: 2, node: 0 },
      { x: 300, y: 0, w: 60, h: 30, url: "data:text/html,x", after: 2, node: 0 },
    ],
  }, 1, { lang: "en" });
  return {
    slide: roles(slide),
    headingLevel: slide.querySelector("[role=heading]")?.getAttribute("aria-level"),
    links: [...slide.querySelectorAll("a")].map((a) => [a.getAttribute("href"), a.getAttribute("data-bdf-page"), a.textContent]),
    flow: roles(flow),
    headers: [...flow.querySelectorAll("[role=columnheader]")].map((e) => e.textContent),
    figures: [...figures.querySelectorAll("[role=img]")].map((e) => e.getAttribute("aria-label")),
    figureBox: [fig.style.left, fig.style.top, fig.style.width, fig.style.height],
    sheet: roles(sheet),
    sheetCell: sheet.querySelector("[role=row]:nth-child(3) [role=cell]")?.getAttribute("aria-colindex"),
    xlsxHeaders: [...workbook.querySelectorAll("[role=columnheader]")].map((e) => e.textContent),
    xlsxFigures: [...workbook.querySelectorAll("[role=img]")].map((e) => e.getAttribute("aria-label")),
    xlsxSpans: [placed, workbook.querySelectorAll("span").length, workbook.querySelectorAll("[role=table]").length],
    unsafeAnchors: [...unsafe.querySelectorAll("a")].map((a) => a.getAttribute("href")),
    lang: [unsafe.lang, ...[...unsafe.querySelectorAll("span")].map((s) => s.lang)],
  };
}

/** Called by test/browser.mjs with forced colors emulated: the text layer's color. */
function forcedColorsCheck() {
  const layer = buildTextLayer([{ text: "x", x: 0, y: 10, advance: 0, size: 10, font: undefined, align: 0, matrix: [1, 0, 0, 1, 0, 0], sep: 0, ordinal: 0, altText: false }], 1);
  document.body.appendChild(layer);
  const out = { forced: matchMedia("(forced-colors: active)").matches, color: getComputedStyle(layer.firstElementChild!).color };
  layer.remove();
  return out;
}

main().catch((e) => {
  (window as unknown as { bdfError: string }).bdfError = String(e?.stack ?? e);
  document.title = "error";
});
