// Test harness: renders fixture pages on the main thread and via the worker,
// and compares canvases against golden PNGs. Driven by test/golden.mjs.
import { BdfDocument, fetchSingle, RangeSource, SplitSource, type Rect } from "@bdf/core";
import { PageRenderer, BdfWorkerClient, buildTextLayer, selectionText, joinRuns, TEXT_LAYER_CSS, RUN_ATTR } from "@bdf/render";

export interface Case {
  name: string;
  /** Document path (default: the generated fixture). */
  src?: string;
  kind: "page" | "continuous" | "sheet";
  view: string;
  page?: number;
  viewport?: Rect;
  scale: number;
  roles?: string[];
}

export const CASES: Case[] = [
  { name: "slides-1", kind: "page", view: "slides", page: 0, scale: 1 },
  { name: "slides-2", kind: "page", view: "slides", page: 1, scale: 1 },
  { name: "slides-3", kind: "page", view: "slides", page: 2, scale: 1 },
  { name: "slides-1-body-only", kind: "page", view: "slides", page: 0, scale: 0.5, roles: ["body"] },
  { name: "doc-1", kind: "page", view: "doc", page: 0, scale: 1 },
  { name: "doc-continuous", kind: "continuous", view: "doc", viewport: { x: 0, y: 600, w: 451.3, h: 400 }, scale: 1 },
  { name: "sheet-tile-boundary", kind: "sheet", view: "sheet1", viewport: { x: 0, y: 1800, w: 800, h: 500 }, scale: 1 },
  { name: "sheet-zoomed", kind: "sheet", view: "sheet1", viewport: { x: 64, y: 20, w: 300, h: 150 }, scale: 2 },
  // Documents converted from PDF (converter/pdf); see test/pdf.
  { name: "pdf-chrome-slides-1", src: "/testdata/pdf/chrome-slides.bdf", kind: "page", view: "pages", page: 0, scale: 1 },
  { name: "pdf-chrome-slides-2", src: "/testdata/pdf/chrome-slides.bdf", kind: "page", view: "pages", page: 1, scale: 1 },
  { name: "pdf-chrome-doc-1", src: "/testdata/pdf/chrome-doc.bdf", kind: "page", view: "pages", page: 0, scale: 1 },
  { name: "pdf-reportlab-1", src: "/testdata/pdf/reportlab-mixed.bdf", kind: "page", view: "pages", page: 0, scale: 1 },
  { name: "pdf-reportlab-2", src: "/testdata/pdf/reportlab-mixed.bdf", kind: "page", view: "pages", page: 1, scale: 1.5 },
  // Three pages that share a master prefix (PDF converter prefix sharing); the pages must render as if unshared.
  { name: "pdf-master-1", src: "/testdata/pdf/reportlab-master.bdf", kind: "page", view: "pages", page: 0, scale: 1 },
  { name: "pdf-master-2", src: "/testdata/pdf/reportlab-master.bdf", kind: "page", view: "pages", page: 1, scale: 1 },
  { name: "pdf-master-3", src: "/testdata/pdf/reportlab-master.bdf", kind: "page", view: "pages", page: 2, scale: 1 },
  // OpenType CFF web fonts (.otf, .woff2) and license-restricted TrueType; bare Type1C and CIDFontType0C from cairo.
  { name: "pdf-webfonts-1", src: "/testdata/pdf/weasyprint-webfonts.bdf", kind: "page", view: "pages", page: 0, scale: 1 },
  { name: "pdf-cairo-cff-1", src: "/testdata/pdf/cairo-cff.bdf", kind: "page", view: "pages", page: 0, scale: 1.5 },
  // Type 1 programs (FontFile) converted to CFF: cairo subsets (seac accents, built-in encoding) and a whole font.
  { name: "pdf-cairo-type1-1", src: "/testdata/pdf/cairo-type1.bdf", kind: "page", view: "pages", page: 0, scale: 1.5 },
  { name: "pdf-reportlab-type1-1", src: "/testdata/pdf/reportlab-type1.bdf", kind: "page", view: "pages", page: 0, scale: 1.5 },
  // PowerPoint decks rendered by converter/pptx with the test fonts; see test/pptx.
  { name: "pptx-basic-2", src: "/testdata/pptx/basic.bdf", kind: "page", view: "slides", page: 1, scale: 1 },
  { name: "pptx-basic-3", src: "/testdata/pptx/basic.bdf", kind: "page", view: "slides", page: 2, scale: 1 },
  { name: "pptx-basic-4", src: "/testdata/pptx/basic.bdf", kind: "page", view: "slides", page: 3, scale: 1 },
  { name: "pptx-basic-5", src: "/testdata/pptx/basic.bdf", kind: "page", view: "slides", page: 4, scale: 1 },
  { name: "pptx-features-1", src: "/testdata/pptx/features.bdf", kind: "page", view: "slides", page: 0, scale: 1 },
  { name: "pptx-features-3", src: "/testdata/pptx/features.bdf", kind: "page", view: "slides", page: 2, scale: 1 },
  { name: "pptx-features-4", src: "/testdata/pptx/features.bdf", kind: "page", view: "slides", page: 3, scale: 1 },
  { name: "pptx-features-5", src: "/testdata/pptx/features.bdf", kind: "page", view: "slides", page: 4, scale: 0.75 },
  // draw.io diagrams rendered by converter/drawio with the test fonts: every page is a view of its own.
  { name: "drawio-labels", src: "/testdata/drawio/labels.bdf", kind: "page", view: "text", page: 0, scale: 1.5 },
  { name: "drawio-multipage-1", src: "/testdata/drawio/multipage.bdf", kind: "page", view: "overview", page: 0, scale: 1.5 },
  { name: "drawio-multipage-2", src: "/testdata/drawio/multipage.bdf", kind: "page", view: "details", page: 0, scale: 1.5 },
  { name: "drawio-multipage-3", src: "/testdata/drawio/multipage.bdf", kind: "page", view: "layers", page: 0, scale: 1.5 },
  // shapes, stencils, markers and edge shapes; swimlanes, orthogonal and curved routing, line jumps, a table
  { name: "drawio-showcase-1", src: "/testdata/drawio/showcase.bdf", kind: "page", view: "shapes", page: 0, scale: 1.5 },
  { name: "drawio-showcase-2", src: "/testdata/drawio/showcase.bdf", kind: "page", view: "flow", page: 0, scale: 1.5 },
  // AWS: current icons and groups (aws4), and a diagram in an older icon set drawn with them
  { name: "drawio-aws-1", src: "/testdata/drawio/aws.bdf", kind: "page", view: "current", page: 0, scale: 1.5 },
  { name: "drawio-aws-2", src: "/testdata/drawio/aws.bdf", kind: "page", view: "legacy", page: 0, scale: 1 },
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

export interface Diff { total: number; different: number; maxDelta: number; width: number; height: number; png: string }

function compare(a: HTMLCanvasElement, b: HTMLCanvasElement, channelTol: number): Diff {
  const w = Math.max(a.width, b.width), h = Math.max(a.height, b.height);
  const da = a.getContext("2d")!.getImageData(0, 0, a.width, a.height).data;
  const db = b.getContext("2d")!.getImageData(0, 0, b.width, b.height).data;
  const diff = document.createElement("canvas");
  diff.width = w;
  diff.height = h;
  const dctx = diff.getContext("2d")!;
  const dimg = dctx.createImageData(w, h);
  let different = 0, maxDelta = 0;
  for (let y = 0; y < h; y++) for (let x = 0; x < w; x++) {
    let delta = 0;
    for (let ch = 0; ch < 4; ch++) {
      const va = x < a.width && y < a.height ? da[(y * a.width + x) * 4 + ch] : -1;
      const vb = x < b.width && y < b.height ? db[(y * b.width + x) * 4 + ch] : -1;
      delta = Math.max(delta, Math.abs(va - vb));
    }
    const o = (y * w + x) * 4;
    dimg.data[o] = 255; dimg.data[o + 1] = delta > channelTol ? 0 : 255; dimg.data[o + 2] = delta > channelTol ? 0 : 255; dimg.data[o + 3] = 255;
    if (delta > channelTol) different++;
    maxDelta = Math.max(maxDelta, delta);
  }
  dctx.putImageData(dimg, 0, 0);
  return { total: w * h, different, maxDelta, width: w, height: h, png: different ? diff.toDataURL("image/png") : "" };
}

async function loadImage(url: string): Promise<HTMLCanvasElement | null> {
  const res = await fetch(url);
  if (!res.ok) return null;
  const bmp = await createImageBitmap(await res.blob());
  return toCanvas(bmp);
}

export interface Result {
  name: string;
  png: string; // data URL of the main-thread render
  workerPng: string; // data URL of the worker render
  golden: Diff | null; // null when no golden exists
  worker: Diff; // main thread vs worker
  textRuns: number;
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
    main.title = c.name;
    document.body.appendChild(main);
    const worker = toCanvas(await renderWorker(client, c));
    const golden = await loadImage(`/testdata/golden/${c.name}.png`);
    const textRuns = c.kind === "page" ? (await client.text(c.view, c.page!)).length : 0;
    results.push({
      name: c.name,
      png: main.toDataURL("image/png"),
      workerPng: worker.toDataURL("image/png"),
      golden: golden ? compare(main, golden, 24) : null,
      worker: compare(main, worker, 0),
      textRuns,
    });
  }
  // text layer: select runs of the first doc page through the DOM and
  // rebuild the text; compare with the runs joined directly.
  (window as unknown as { bdfSelection: unknown }).bdfSelection = await selectionCheck(client);
  // ALT_TEXT runs (ligatures drawn as one private-use glyph) are in the text layer
  const { client: chromeDoc } = await open("/testdata/pdf/chrome-doc.bdf");
  (window as unknown as { bdfAltText: unknown }).bdfAltText = await altTextCheck(chromeDoc);
  (window as unknown as { bdfForcedColors: unknown }).bdfForcedColors = forcedColorsCheck;
  (window as unknown as { bdfStructure: unknown }).bdfStructure = await structureCheck(client);
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
  // draw.io: a page of the diagram is a view; the search runs in it, and
  // links to other pages are #view= links
  const { client: drawio } = await open("/testdata/drawio/multipage.bdf");
  const drawioHits = await drawio.search("details", "日本語の説明");
  const drawioRects = await drawio.locate("details", drawioHits);
  const drawioLinks = (await drawio.content("overview", 0)).links.map((l) => l.url);
  (window as unknown as { bdfSearch: unknown }).bdfSearch = { hits, rects, sheetHits, sheetRects, pptxHits, pptxRects, ligHits, ligRects, drawioHits, drawioRects, drawioLinks };
  (window as unknown as { bdfResults: Result[] }).bdfResults = results;
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
async function structureCheck(client: BdfWorkerClient) {
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
    unsafeAnchors: [...unsafe.querySelectorAll("a")].map((a) => a.getAttribute("href")),
    lang: [unsafe.lang, ...[...unsafe.querySelectorAll("span")].map((s) => s.lang)],
  };
}

/** Called by test/golden.mjs with forced colors emulated: the text layer's color. */
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
