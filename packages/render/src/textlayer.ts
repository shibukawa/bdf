// Selectable, accessible text over a rendered page (pdf.js style): one
// transparent, absolutely positioned span per text run, so the browser's own
// selection works on top of the bitmap. The spans sit inside elements that
// carry the document's structure for assistive technology (paragraphs,
// headings, lists, tables, figures with their alternative text, links,
// languages; spec §7.8); those elements take no space, except figures and
// links, which cover their area. Copying goes through selectionText(), which
// rebuilds the text from the runs' separators (spaces at line breaks, line
// breaks at paragraph/cell/box boundaries) instead of the browser's
// serialization of unrelated spans.
import { Sep, type TextRun, type TextContent, type TextNode, type TextLink } from "@bdf/core";
import { fontString } from "./resources.js";
import { hasExtent } from "./search.js";

export interface TextLayerOptions {
  /** Class of the layer element (default "bdfTextLayer"). */
  className?: string;
  /**
   * Include ALT_TEXT runs; default true. They carry the real text of what was
   * drawn: ligatures and other glyphs mapped to private-use characters (placed
   * like any run), and text drawn as paths or child objects (placed at their
   * anchor along the baseline direction). Leaving them out drops that text
   * from selection and from screen readers.
   */
  altText?: boolean;
  /** Width measurement in CSS px for a CSS font string (default: a shared 2D context). */
  measure?: (font: string, text: string) => number;
  /**
   * Language of the layer (the first of the document's meta.dc.language; ""
   * for unknown). Set it:
   * the page around the layer is usually in the viewer's own language. Runs
   * in another language get their own lang attribute.
   */
  lang?: string;
  /**
   * Lay out cells outside tables (sheet tiles) as one table of rows sorted by
   * position, with the sheet's size and its header rows and columns.
   */
  sheet?: { rows: number; cols: number; headerRows?: number; headerCols?: number };
  /** Accessible name of a link that covers no text (default: the URL, "page N" or "view ID"). */
  linkLabel?: (url: string) => string;
}

let sharedMeasure: ((font: string, text: string) => number) | undefined;
function defaultMeasure(): (font: string, text: string) => number {
  if (!sharedMeasure) {
    const ctx = document.createElement("canvas").getContext("2d")!;
    sharedMeasure = (font, text) => {
      ctx.font = font;
      return ctx.measureText(text).width;
    };
  }
  return sharedMeasure;
}

/** Attribute names: on the run spans, on the layer element, and on links to a page or a view. */
export const RUN_ATTR = { ordinal: "data-bdf-ordinal", sep: "data-bdf-sep", layer: "data-bdf-layer", page: "data-bdf-page", view: "data-bdf-view" } as const;

/**
 * A link within the document: "#page=N" (a page of the same view) or
 * "#view=ID" with an optional "&page=N" (another view, e.g. another page
 * of a diagram), spec §7.7. Pages are 1-based.
 */
export interface InternalLink { view?: string; page?: number }

/** Parse an internal link url; undefined for anything else. */
export function internalLink(url: string): InternalLink | undefined {
  const page = /^#page=(\d{1,7})$/.exec(url);
  if (page) return Number(page[1]) >= 1 ? { page: Number(page[1]) } : undefined;
  const view = /^#view=([^&#]+)(?:&page=(\d{1,7}))?$/.exec(url);
  if (!view) return undefined;
  let id: string;
  try {
    id = decodeURIComponent(view[1].replace(/\+/g, " "));
  } catch {
    return undefined;
  }
  if (view[2] !== undefined && Number(view[2]) < 1) return undefined;
  return view[2] !== undefined ? { view: id, page: Number(view[2]) } : { view: id };
}

/**
 * The href of a LINK url, or undefined when a link must not be made: only
 * http, https and mailto URLs, "#page=N" and "#view=ID" (spec §7.7).
 */
export function linkHref(url: string): string | undefined {
  if (url.startsWith("#")) return internalLink(url) ? url : undefined;
  try {
    const u = new URL(url);
    return u.protocol === "http:" || u.protocol === "https:" || u.protocol === "mailto:" ? u.href : undefined;
  } catch {
    return undefined;
  }
}

/** An element and the layer position (CSS px) that its absolutely positioned children are placed from. */
interface Placed { el: HTMLElement; ox: number; oy: number }

/**
 * Build the text layer of one page, band or sheet region: the extractor's
 * runs in reading order, or its content (runs with their structure and
 * links), in page units; scale is CSS px per unit. The element must be
 * placed over the page (position: absolute; inset: 0) and needs
 * TEXT_LAYER_CSS.
 */
export function buildTextLayer(input: TextRun[] | TextContent, scale: number, opts: TextLayerOptions = {}): HTMLDivElement {
  const layer = document.createElement("div");
  layer.className = opts.className ?? "bdfTextLayer";
  layer.setAttribute(RUN_ATTR.layer, "");
  if (opts.lang !== undefined) layer.lang = opts.lang;
  const content: TextContent = Array.isArray(input) ? { runs: input, nodes: [], links: [] } : input;
  new LayerBuilder(layer, content, scale, opts).build();
  return layer;
}

class LayerBuilder {
  private readonly root: Placed;
  private readonly measure: (font: string, text: string) => number;
  private readonly placed = new Map<number, Placed>();
  /** Current row of each table element. */
  private readonly rows = new Map<HTMLElement, { row: number; el: HTMLElement }>();
  /** Link covering each run (-1: none). */
  private linkOf: number[] = [];
  private textLinks = new Set<number>();
  /** The anchor the previous run went into, to continue it. */
  private anchor: { parent: HTMLElement; link: number; placed: Placed } | undefined;

  constructor(layer: HTMLElement, private c: TextContent, private scale: number, private opts: TextLayerOptions) {
    this.root = { el: layer, ox: 0, oy: 0 };
    this.measure = opts.measure ?? defaultMeasure();
  }

  build() {
    const { runs, nodes, links } = this.c;
    const keep = (r: TextRun) => !!r.text && !(r.altText && this.opts.altText === false);
    this.linkOf = runs.map((r) => (keep(r) ? this.coveringLink(r) : -1));
    for (const l of this.linkOf) if (l >= 0) this.textLinks.add(l);
    if (this.opts.sheet) this.sheetTable();
    // Nodes are created in stream order and a run always belongs to the newest
    // node, so walking both in step keeps the DOM in reading order: figures
    // without text are placed where they were drawn, other nodes on demand.
    let nextNode = 0, nextLink = 0;
    const standalone = (upTo: number) => {
      for (; nextLink < links.length && links[nextLink].after < upTo; nextLink++) {
        if (!this.textLinks.has(nextLink)) this.standaloneLink(nextLink);
      }
    };
    runs.forEach((r, i) => {
      for (; r.node !== undefined && nextNode <= r.node && nextNode < nodes.length; nextNode++) this.figure(nextNode);
      standalone(i);
      if (keep(r)) this.run(r, i);
    });
    for (; nextNode < nodes.length; nextNode++) this.figure(nextNode);
    standalone(Infinity);
  }

  /** Place a figure that has alternative text even when it covers no text. */
  private figure(i: number) {
    const n = this.c.nodes[i];
    if (n.kind === "figure" && n.alt) this.node(i);
  }

  /** The element of a node, created with its ancestors on first use. */
  private node(i: number | undefined): Placed {
    if (i === undefined || i < 0 || i >= this.c.nodes.length) return this.root;
    const got = this.placed.get(i);
    if (got) return got;
    const n = this.c.nodes[i];
    const parent = this.node(n.parent);
    const parentNode = n.parent >= 0 ? this.c.nodes[n.parent] : undefined;
    let placed: Placed;
    if (n.kind === "cell" && parentNode?.kind === "table") {
      placed = { ...parent, el: this.cell(n, this.row(parent.el, n.row ?? 0), n.scope) };
    } else if (n.kind === "cell") {
      // a sheet cell outside sheet mode (sheetTable places the others)
      placed = { ...parent, el: this.div(parent.el, "paragraph") };
    } else if (n.kind === "figure") {
      const b = n.bounds;
      if (n.alt && b) {
        const el = this.div(parent.el, "img");
        el.setAttribute("aria-label", n.alt);
        placed = this.position(el, parent, b.x, b.y, b.w, b.h);
      } else {
        placed = { ...parent, el: this.div(parent.el) }; // decorative: no role
      }
    } else {
      const el = this.div(parent.el, ROLE[n.kind]);
      if (n.kind === "heading" && n.level) el.setAttribute("aria-level", String(n.level));
      if (n.kind === "table") this.tableSize(el, i);
      placed = { ...parent, el };
    }
    this.placed.set(i, placed);
    return placed;
  }

  private div(parent: HTMLElement, role?: string): HTMLElement {
    const el = document.createElement("div");
    if (role) el.setAttribute("role", role);
    parent.appendChild(el);
    return el;
  }

  /** Absolutely position el over a rectangle (page units); returns it as the origin of its children. */
  private position(el: HTMLElement, parent: Placed, x: number, y: number, w: number, h: number): Placed {
    const s = this.scale;
    el.style.left = `${x * s - parent.ox}px`;
    el.style.top = `${y * s - parent.oy}px`;
    el.style.width = `${w * s}px`;
    el.style.height = `${h * s}px`;
    return { el, ox: x * s, oy: y * s };
  }

  private tableSize(el: HTMLElement, table: number) {
    let rows = 0, cols = 0;
    for (const n of this.c.nodes) {
      if (n.parent !== table || n.kind !== "cell") continue;
      rows = Math.max(rows, (n.row ?? 0) + (n.rows ?? 1));
      cols = Math.max(cols, (n.col ?? 0) + (n.cols ?? 1));
    }
    if (rows) el.setAttribute("aria-rowcount", String(rows));
    if (cols) el.setAttribute("aria-colcount", String(cols));
  }

  /** The row element for a cell starting at row; a new one when the row changes (cells are never reordered). */
  private row(table: HTMLElement, row: number): HTMLElement {
    const cur = this.rows.get(table);
    if (cur && cur.row === row) return cur.el;
    const el = this.div(table, "row");
    el.setAttribute("aria-rowindex", String(row + 1));
    this.rows.set(table, { row, el });
    return el;
  }

  private cell(n: TextNode, row: HTMLElement, scope: "col" | "row" | undefined): HTMLElement {
    const el = this.div(row, scope === "col" ? "columnheader" : scope === "row" ? "rowheader" : "cell");
    el.setAttribute("aria-colindex", String((n.col ?? 0) + 1));
    if ((n.rows ?? 1) > 1) el.setAttribute("aria-rowspan", String(n.rows));
    if ((n.cols ?? 1) > 1) el.setAttribute("aria-colspan", String(n.cols));
    return el;
  }

  /** Sheet mode: one table of the cells outside tables, rows sorted by position. */
  private sheetTable() {
    const sheet = this.opts.sheet!;
    const cells: number[] = [];
    const used = new Set(this.c.runs.map((r) => r.node));
    this.c.nodes.forEach((n, i) => {
      if (n.kind === "cell" && n.parent < 0 && n.row !== undefined && used.has(i)) cells.push(i);
    });
    if (!cells.length) return;
    const table = this.div(this.root.el, "table");
    table.setAttribute("aria-rowcount", String(sheet.rows));
    table.setAttribute("aria-colcount", String(sheet.cols));
    const at = (i: number) => this.c.nodes[i];
    cells.sort((a, b) => at(a).row! - at(b).row! || at(a).col! - at(b).col!);
    for (const i of cells) {
      const n = at(i);
      const scope = n.row! < (sheet.headerRows ?? 0) ? "col" : n.col! < (sheet.headerCols ?? 0) ? "row" : undefined;
      this.placed.set(i, { ...this.root, el: this.cell(n, this.row(table, n.row!), scope) });
    }
  }

  /** The smallest link whose rectangle contains the middle of a run (-1: none). */
  private coveringLink(r: TextRun): number {
    const links = this.c.links;
    if (!links.length) return -1;
    const w = r.advance > 0 ? r.advance : r.size * r.text.length * 0.5;
    const cx = (r.align === 1 ? -w / 2 : r.align === 2 ? 0 : w / 2), cy = -r.size * 0.3;
    const m = r.matrix;
    const x = m[0] * cx + m[2] * cy + r.x, y = m[1] * cx + m[3] * cy + r.y;
    let best = -1;
    links.forEach((l, i) => {
      if (!linkHref(l.url) || x < l.x || x > l.x + l.w || y < l.y || y > l.y + l.h) return;
      if (best < 0 || l.w * l.h < links[best].w * links[best].h) best = i;
    });
    return best;
  }

  private anchorFor(link: TextLink, parent: Placed): Placed {
    const a = document.createElement("a");
    const href = linkHref(link.url)!;
    a.href = href;
    a.draggable = false;
    const internal = internalLink(href);
    if (internal) {
      if (internal.view !== undefined) a.setAttribute(RUN_ATTR.view, internal.view);
      if (internal.page !== undefined) a.setAttribute(RUN_ATTR.page, String(internal.page));
    } else {
      a.target = "_blank";
      a.rel = "noopener noreferrer";
    }
    parent.el.appendChild(a);
    return this.position(a, parent, link.x, link.y, link.w, link.h);
  }

  private standaloneLink(i: number) {
    const link = this.c.links[i];
    if (!linkHref(link.url)) return;
    const placed = this.anchorFor(link, this.node(link.node));
    const internal = internalLink(link.url);
    const fallback = internal?.view !== undefined ? `view ${internal.view}` : internal?.page !== undefined ? `page ${internal.page}` : link.url;
    placed.el.setAttribute("aria-label", this.opts.linkLabel?.(link.url) ?? fallback);
  }

  private run(r: TextRun, i: number) {
    let target = this.node(r.node);
    const link = this.linkOf[i];
    if (link >= 0) {
      // consecutive runs of one leaf under the same link share an anchor
      const a = this.anchor;
      if (a && a.parent === target.el && a.link === link && a.placed.el === target.el.lastElementChild) target = a.placed;
      else {
        const placed = this.anchorFor(this.c.links[link], target);
        this.anchor = { parent: target.el, link, placed };
        target = placed;
      }
    }
    target.el.appendChild(this.span(r, target));
  }

  private span(r: TextRun, at: Placed): HTMLSpanElement {
    const span = document.createElement("span");
    span.textContent = r.text;
    span.setAttribute(RUN_ATTR.ordinal, String(r.ordinal));
    span.setAttribute(RUN_ATTR.sep, String(r.sep));
    if (r.lang && r.lang !== (this.opts.lang ?? "")) span.lang = r.lang;
    const scale = this.scale;
    if (hasExtent(r)) {
      const font = fontString(r.font!, r.size);
      span.style.font = font;
      // The main thread may not have the embedded font: stretch the fallback
      // rendering to the advance the worker measured with the real one.
      const w = this.measure(font, r.text);
      const sx = r.advance > 0 && w > 0 ? r.advance / w : 1;
      const width = r.advance > 0 ? r.advance : w;
      const anchor = r.align === 1 ? -width : r.align === 2 ? -width / 2 : 0;
      const m = r.matrix;
      span.style.transform = `matrix(${m[0] * scale}, ${m[1] * scale}, ${m[2] * scale}, ${m[3] * scale}, ${r.x * scale - at.ox}, ${r.y * scale - at.oy}) translate(${anchor}px, ${-r.size * 0.8}px) scaleX(${sx})`;
    } else {
      // ALT_TEXT for paths, a child object or no drawing op: only the anchor
      // and the baseline direction are known, so lay the text there.
      const size = (r.size > 0 ? r.size : 10) * scale;
      const angle = Math.atan2(r.matrix[1], r.matrix[0]);
      span.style.font = `${size}px sans-serif`;
      span.style.transform = `translate(${r.x * scale - at.ox}px, ${r.y * scale - at.oy}px) rotate(${angle}rad) translate(0, ${-size * 0.8}px)`;
    }
    return span;
  }
}

const ROLE: Record<TextNode["kind"], string> = {
  paragraph: "paragraph", heading: "heading", caption: "caption", list: "list", item: "listitem",
  table: "table", cell: "cell", figure: "img",
};

/**
 * CSS that a text layer needs; the class names follow TextLayerOptions.className.
 * forced-color-adjust keeps the text transparent in forced colors (Windows
 * high contrast), which would otherwise paint it opaque over the canvas.
 * Figures and links are positioned boxes; a figure lets pointer events
 * through to what lies under it, except on its text.
 */
export const TEXT_LAYER_CSS = `
.bdfTextLayer { position: absolute; inset: 0; overflow: hidden; line-height: 1; forced-color-adjust: none; }
.bdfTextLayer span { position: absolute; color: transparent; white-space: pre; transform-origin: 0 0; cursor: text; }
.bdfTextLayer span::selection { background: rgba(0, 120, 255, .35); }
.bdfTextLayer [role=img], .bdfTextLayer a { position: absolute; }
.bdfTextLayer [role=img] { pointer-events: none; }
.bdfTextLayer [role=img] span { pointer-events: auto; }
.bdfTextLayer a, .bdfTextLayer a span { cursor: pointer; -webkit-user-drag: none; }
.bdfTextLayer a:focus-visible { outline: 2px solid #1a73e8; outline-offset: 1px; }
`;

/** A selected piece of a run. */
export interface SelectedRun {
  ordinal: number;
  sep: number;
  text: string;
  /** Element holding the run; layers differ when the selection spans pages. */
  layer: Element | null;
}

/**
 * The runs (or parts of runs) a DOM range covers, in document order: every
 * run span under root that intersects the range, clipped to it.
 */
export function selectedRuns(range: Range, root: ParentNode = document): SelectedRun[] {
  const out: SelectedRun[] = [];
  for (const span of root.querySelectorAll<HTMLElement>(`[${RUN_ATTR.ordinal}]`)) {
    const node = span.firstChild;
    if (!node || !range.intersectsNode(span)) continue;
    const full = node.textContent ?? "";
    let start = 0, end = full.length;
    // The range starts inside this run when its first character lies before the
    // range start; the start offset is only meaningful within the text node.
    if (range.comparePoint(node, 0) < 0) start = range.startContainer === node ? range.startOffset : 0;
    if (range.comparePoint(node, full.length) > 0) end = range.endContainer === node ? range.endOffset : full.length;
    if (end <= start) continue;
    out.push({ ordinal: Number(span.getAttribute(RUN_ATTR.ordinal)), sep: Number(span.getAttribute(RUN_ATTR.sep)), text: full.slice(start, end), layer: span.closest(`[${RUN_ATTR.layer}]`) });
  }
  return out;
}

/**
 * Join selected runs into plain text: Sep.SPACE becomes a space, Sep.BREAK a
 * line break, and a change of layer (page) a blank line. The first run's
 * own separator is dropped.
 */
export function joinRuns(parts: SelectedRun[]): string {
  let s = "";
  for (let i = 0; i < parts.length; i++) {
    const p = parts[i];
    if (i > 0) {
      if (p.layer !== parts[i - 1].layer) s += "\n\n";
      else if (p.sep === Sep.BREAK) s += "\n";
      else if (p.sep === Sep.SPACE) s += " ";
    }
    s += p.text;
  }
  return s;
}

/** Plain text of the current selection within root (empty when nothing is selected there). */
export function selectionText(sel: Selection | null = typeof getSelection === "function" ? getSelection() : null, root: ParentNode = document): string {
  if (!sel || sel.rangeCount === 0 || sel.isCollapsed) return "";
  let s = "";
  for (let i = 0; i < sel.rangeCount; i++) {
    const t = joinRuns(selectedRuns(sel.getRangeAt(i), root));
    if (t) s = s ? `${s}\n${t}` : t;
  }
  return s;
}

/**
 * Make copy (and cut) inside container put the selection's run text on the
 * clipboard instead of the browser's serialization. Returns a function that
 * removes the handler.
 */
export function installCopyHandler(container: HTMLElement): () => void {
  const onCopy = (e: ClipboardEvent) => {
    const text = selectionText(getSelection(), container);
    if (!text || !e.clipboardData) return;
    e.clipboardData.setData("text/plain", text);
    e.preventDefault();
  };
  container.addEventListener("copy", onCopy);
  container.addEventListener("cut", onCopy);
  return () => {
    container.removeEventListener("copy", onCopy);
    container.removeEventListener("cut", onCopy);
  };
}
