import { NoopSink, walk } from "./object.js";
import type { ObjectPart, Font, Hash, PathData, Rect, Glyph } from "./types.js";

/** MARK kinds (docs/spec.md §7.8). */
export const Mark = {
  PARAGRAPH: 0, LINE: 1, CELL: 2, BOX: 3, ALT_TEXT: 4, WRAP: 5,
  HEADING: 6, LIST: 7, LIST_ITEM: 8, TABLE: 9, FIGURE: 10, END: 11, LANG: 12,
} as const;

/** Separators between consecutive runs (docs/spec.md §7.9). */
export const Sep = { NONE: 0, SPACE: 1, BREAK: 2 } as const;

/** A text run in the coordinate space passed to extractText. */
export interface TextRun {
  text: string;
  /** Anchor position after applying the transform. */
  x: number;
  y: number;
  /** Advance in the object's units (0 when unknown). */
  advance: number;
  size: number;
  font: Font | undefined;
  align: number;
  /** Full 2D transform at the time of drawing (a b c d e f). */
  matrix: [number, number, number, number, number, number];
  /** Separator from the previous run (Sep.*). */
  sep: number;
  /** Position in the extraction order of the top-level object; matches the text index. */
  ordinal: number;
  /** True for ALT_TEXT marks (nothing was drawn with fillText). */
  altText: boolean;
  /** Index of the structure node (TextContent.nodes) the run belongs to. */
  node?: number;
  /** Language (BCP 47) set by MARK LANG; absent for the document default. */
  lang?: string;
}

/**
 * Kinds of structure nodes (docs/spec.md §7.8). Leaves (paragraph, heading,
 * caption, and cells outside a table) hold runs; containers (list, item,
 * table, cells inside a table, figure) hold nodes.
 */
export type TextNodeKind = "paragraph" | "heading" | "caption" | "list" | "item" | "table" | "cell" | "figure";

/** A structure node. Nodes and runs point at their parent by index. */
export interface TextNode {
  kind: TextNodeKind;
  /** Index of the enclosing node, -1 at the top level. */
  parent: number;
  /** Heading level 1–6; absent when the payload is not a level. */
  level?: number;
  /** Cell position (0-based) and span: relative to the table inside a TABLE, absolute in a sheet. */
  row?: number;
  col?: number;
  rows?: number;
  cols?: number;
  /** Header cell. */
  scope?: "col" | "row";
  /** Figure: alternative text ("" when decorative). */
  alt?: string;
  /** Figure: bounding box of what it draws, clipped, in the extraction space. */
  bounds?: Rect;
}

/** A LINK area: the bounding box of its rectangle in the extraction space. */
export interface TextLink extends Rect {
  url: string;
  /** Index of the last run before the LINK op (-1 before the first run). */
  after: number;
  /** Node open at the LINK op (-1 at the top level). */
  node: number;
}

/** Runs with the structure and links of an object. */
export interface TextContent {
  runs: TextRun[];
  nodes: TextNode[];
  links: TextLink[];
}

/** A parsed cell reference ("B2", "A4:B4 col"). */
export interface CellRef { row: number; col: number; rows: number; cols: number; scope?: "col" | "row" }

/** Parse a CELL payload: an A1 reference or range, optionally followed by " col" or " row". */
export function parseCellRef(payload: string): CellRef | undefined {
  const m = /^\s*([A-Za-z]+)(\d+)(?::([A-Za-z]+)(\d+))?(?:\s+(col|row))?\s*$/.exec(payload);
  if (!m) return undefined;
  const colOf = (s: string) => [...s.toUpperCase()].reduce((a, ch) => a * 26 + ch.charCodeAt(0) - 64, 0) - 1;
  const row = Number(m[2]) - 1, col = colOf(m[1]);
  if (row < 0) return undefined;
  const row2 = m[4] ? Number(m[4]) - 1 : row, col2 = m[3] ? colOf(m[3]) : col;
  const ref: CellRef = { row, col, rows: Math.max(1, row2 - row + 1), cols: Math.max(1, col2 - col + 1) };
  if (m[5]) ref.scope = m[5] as "col" | "row";
  return ref;
}

type Matrix = [number, number, number, number, number, number];

function mul(m: Matrix, n: Matrix): Matrix {
  return [
    m[0] * n[0] + m[2] * n[1],
    m[1] * n[0] + m[3] * n[1],
    m[0] * n[2] + m[2] * n[3],
    m[1] * n[2] + m[3] * n[3],
    m[0] * n[4] + m[2] * n[5] + m[4],
    m[1] * n[4] + m[3] * n[5] + m[5],
  ];
}

/** Bounding box of a rectangle moved by m. */
function boxOf(m: Matrix, x: number, y: number, w: number, h: number): Rect {
  let x0 = Infinity, y0 = Infinity, x1 = -Infinity, y1 = -Infinity;
  for (const [px, py] of [[x, y], [x + w, y], [x, y + h], [x + w, y + h]]) {
    const tx = m[0] * px + m[2] * py + m[4], ty = m[1] * px + m[3] * py + m[5];
    x0 = Math.min(x0, tx); y0 = Math.min(y0, ty); x1 = Math.max(x1, tx); y1 = Math.max(y1, ty);
  }
  return { x: x0, y: y0, w: x1 - x0, h: y1 - y0 };
}

function intersect(a: Rect, b: Rect | undefined): Rect | undefined {
  if (!b) return a;
  const x0 = Math.max(a.x, b.x), y0 = Math.max(a.y, b.y);
  const x1 = Math.min(a.x + a.w, b.x + b.w), y1 = Math.min(a.y + a.h, b.y + b.h);
  return x1 >= x0 && y1 >= y0 ? { x: x0, y: y0, w: x1 - x0, h: y1 - y0 } : undefined;
}

function union(a: Rect | undefined, b: Rect): Rect {
  if (!a) return b;
  const x0 = Math.min(a.x, b.x), y0 = Math.min(a.y, b.y);
  return { x: x0, y: y0, w: Math.max(a.x + a.w, b.x + b.w) - x0, h: Math.max(a.y + a.h, b.y + b.h) - y0 };
}

/** Bounding box of the control points of a path (arcs and ellipses by their extent). */
function pathBox(p: PathData): Rect | undefined {
  let x0 = Infinity, y0 = Infinity, x1 = -Infinity, y1 = -Infinity;
  const add = (x: number, y: number) => { x0 = Math.min(x0, x); y0 = Math.min(y0, y); x1 = Math.max(x1, x); y1 = Math.max(y1, y); };
  const a = p.args;
  let i = 0;
  for (const v of p.verbs) {
    switch (v) {
      case 0: case 1: add(a[i], a[i + 1]); i += 2; break; // MOVE, LINE
      case 2: add(a[i], a[i + 1]); add(a[i + 2], a[i + 3]); i += 4; break; // QUAD
      case 3: add(a[i], a[i + 1]); add(a[i + 2], a[i + 3]); add(a[i + 4], a[i + 5]); i += 6; break; // CUBIC
      case 4: break; // CLOSE
      case 5: add(a[i], a[i + 1]); add(a[i] + a[i + 2], a[i + 1] + a[i + 3]); i += 4; break; // RECT
      case 6: { const r = Math.max(Math.abs(a[i + 2]), Math.abs(a[i + 3])); add(a[i] - r, a[i + 1] - r); add(a[i] + r, a[i + 1] + r); i += 8; break; } // ELLIPSE
      case 7: add(a[i], a[i + 1]); add(a[i + 2], a[i + 3]); i += 5; break; // ARC_TO
      case 8: add(a[i], a[i + 1]); add(a[i] + a[i + 2], a[i + 1] + a[i + 3]); i += 5; break; // ROUND_RECT
    }
  }
  return x1 >= x0 ? { x: x0, y: y0, w: x1 - x0, h: y1 - y0 } : undefined;
}

interface State { m: Matrix; font: Font | undefined; size: number; align: number; clip?: Rect }

type NodeTemplate = Omit<TextNode, "parent">;

/**
 * Shared run list and MARK state across nested objects: separators, the
 * pending ALT_TEXT, and the structure (open containers, the current leaf,
 * the language), which run on through USE in walk order (spec §7.8).
 */
class Extraction {
  runs: TextRun[] = [];
  pending = Sep.NONE as number;
  hasMark = false;
  /** ALT_TEXT waiting for the drawing op it describes. */
  alt: string | undefined;
  nodes: TextNode[] = [];
  links: TextLink[] = [];
  /** Open containers, outermost first. */
  open: number[] = [];
  /** Leaf the next run joins (-1: start a new one). */
  leaf = -1;
  /** Leaf kind reserved by a leaf MARK for the next run. */
  reserved: NodeTemplate | undefined;
  lang = "";
  /** Open figures (their bounds grow with what is drawn). */
  figures: number[] = [];
  /** Where each open figure started, for figures that draw nothing. */
  private origins = new Map<number, Rect>();
  constructor(readonly resolve: (h: Hash) => ObjectPart | undefined) {}

  mark(sep: number) {
    if (!this.hasMark || sep > this.pending) this.pending = sep;
    this.hasMark = true;
  }

  /** Emit an ALT_TEXT that no drawing op consumed, without a position. */
  flushAlt(st: State) {
    if (this.alt !== undefined) {
      const t = this.alt;
      this.alt = undefined;
      this.emit(t, 0, 0, 0, st, true);
    }
  }

  emit(text: string, x: number, y: number, advance: number, st: State, altText: boolean) {
    const m = st.m;
    const run: TextRun = {
      text, advance, size: st.size, font: st.font, align: st.align, matrix: m, altText,
      x: m[0] * x + m[2] * y + m[4], y: m[1] * x + m[3] * y + m[5],
      sep: Sep.NONE, ordinal: this.runs.length,
    };
    if (this.hasMark) run.sep = this.pending;
    else if (this.runs.length) run.sep = guessSep(this.runs[this.runs.length - 1], run);
    this.pending = Sep.NONE;
    this.hasMark = false;
    run.node = this.leafNode();
    if (this.lang) run.lang = this.lang;
    if (this.figures.length) {
      const w = advance > 0 ? advance : st.size * text.length * 0.5;
      const anchor = st.align === 1 ? -w : st.align === 2 ? -w / 2 : 0;
      this.grow(boxOf(m, x + anchor, y - st.size * 0.8, w, st.size), st.clip);
    }
    this.runs.push(run);
  }

  private top(): number { return this.open.length ? this.open[this.open.length - 1] : -1; }
  private topKind(): TextNodeKind | undefined { const t = this.top(); return t >= 0 ? this.nodes[t].kind : undefined; }

  private add(n: TextNode): number {
    this.nodes.push(n);
    return this.nodes.length - 1;
  }

  /** The leaf of the next run, created from the reserved kind or implicitly. */
  private leafNode(): number {
    if (this.leaf >= 0 && !this.reserved) return this.leaf;
    const r = this.reserved;
    this.reserved = undefined;
    // runs directly in a list form an implicit item, in a table its caption
    if (this.topKind() === "list") this.openNode({ kind: "item" });
    const parent = this.top();
    if (this.topKind() === "table") this.leaf = this.add({ kind: "caption", parent });
    else this.leaf = this.add({ ...(r ?? { kind: "paragraph" }), parent });
    return this.leaf;
  }

  private openNode(t: NodeTemplate, st?: State): number {
    const i = this.add({ ...t, parent: this.top() });
    this.open.push(i);
    this.leaf = -1;
    this.reserved = undefined;
    if (t.kind === "figure") {
      this.figures.push(i);
      this.origins.set(i, { x: st?.m[4] ?? 0, y: st?.m[5] ?? 0, w: 0, h: 0 });
    }
    return i;
  }

  private closeTop() {
    const i = this.open.pop();
    this.leaf = -1;
    this.reserved = undefined;
    if (i === undefined || this.nodes[i].kind !== "figure") return;
    this.figures = this.figures.filter((f) => f !== i);
    this.nodes[i].bounds ??= this.origins.get(i);
    this.origins.delete(i);
  }

  closeAll() { while (this.open.length) this.closeTop(); }

  /** Apply a structure MARK (spec §7.8). */
  structure(kind: number, payload: string, st: State) {
    switch (kind) {
      case Mark.PARAGRAPH: case Mark.BOX: this.reserved = { kind: "paragraph" }; break;
      case Mark.HEADING: {
        const level = Number(payload);
        this.reserved = Number.isInteger(level) && level >= 1 && level <= 6 ? { kind: "heading", level } : { kind: "heading" };
        break;
      }
      case Mark.LIST: this.openNode({ kind: "list" }); break;
      case Mark.LIST_ITEM:
        if (this.topKind() === "item") this.closeTop();
        if (this.topKind() === "list") this.openNode({ kind: "item" });
        else this.reserved = { kind: "paragraph" };
        break;
      case Mark.TABLE: this.openNode({ kind: "table" }); break;
      case Mark.CELL: {
        const ref = parseCellRef(payload);
        if (this.topKind() === "cell") this.closeTop();
        if (this.topKind() === "table") this.openNode({ kind: "cell", ...ref });
        else this.reserved = { kind: "cell", ...ref };
        break;
      }
      case Mark.FIGURE: this.openNode({ kind: "figure", alt: payload }, st); break;
      case Mark.END: {
        const k = this.topKind();
        if (k === "item" || k === "cell") this.closeTop();
        if (this.open.length) this.closeTop();
        break;
      }
    }
  }

  /** Add a drawn area to the open figures. */
  grow(box: Rect | undefined, clip: Rect | undefined) {
    const b = box && intersect(box, clip);
    if (!b) return;
    for (const f of this.figures) this.nodes[f].bounds = union(this.nodes[f].bounds, b);
  }

  link(x: number, y: number, w: number, h: number, url: string, m: Matrix) {
    this.links.push({ ...boxOf(m, x, y, w, h), url, after: this.runs.length - 1, node: this.leaf >= 0 ? this.leaf : this.top() });
  }
}

/** Estimate the separator between runs that carry no MARK information. */
export function guessSep(prev: TextRun, cur: TextRun): number {
  const size = prev.size > 0 ? prev.size : 10;
  if (Math.abs(cur.y - prev.y) > size * 0.5) return Sep.SPACE;
  const end = prev.x + prev.advance;
  if (prev.advance === 0 || cur.x - end > Math.fround(0.2) * size) return Sep.SPACE;
  return Sep.NONE;
}

class TextSink extends NoopSink {
  private stack: State[] = [];
  private st: State;

  constructor(private ex: Extraction, private obj: ObjectPart, m: Matrix, clip?: Rect) {
    super();
    this.st = { m, font: undefined, size: 10, align: 0, clip };
  }
  override save() { this.stack.push({ ...this.st }); }
  override restore() { this.st = this.stack.pop() ?? this.st; }
  override transform(a: number, b: number, c: number, d: number, e: number, f: number) { this.st.m = mul(this.st.m, [a, b, c, d, e, f]); }
  override translate(x: number, y: number) { this.transform(1, 0, 0, 1, x, y); }
  override scale(x: number, y: number) { this.transform(x, 0, 0, y, 0, 0); }
  override font(font: number, size: number) { this.st.font = this.obj.fonts[font]; this.st.size = size; }
  override textStyle(align: number) { this.st.align = align; }
  override mark(kind: number, payload: string) {
    switch (kind) {
      case Mark.LINE: this.ex.flushAlt(this.st); this.ex.mark(Sep.SPACE); break;
      case Mark.PARAGRAPH: case Mark.CELL: case Mark.BOX:
      case Mark.HEADING: case Mark.LIST: case Mark.LIST_ITEM: case Mark.TABLE: case Mark.FIGURE: case Mark.END:
        this.ex.flushAlt(this.st); this.ex.mark(Sep.BREAK); this.ex.structure(kind, payload, this.st); break;
      case Mark.ALT_TEXT: this.ex.flushAlt(this.st); this.ex.alt = payload; break;
      case Mark.WRAP: this.ex.flushAlt(this.st); this.ex.mark(Sep.NONE); break;
      // no separator, and a pending ALT_TEXT still belongs to the next drawing op
      case Mark.LANG: this.ex.lang = payload; break;
    }
  }
  override link(x: number, y: number, w: number, h: number, url: string) { this.ex.link(x, y, w, h, url, this.st.m); }

  // --- figure bounds: only computed while a figure is open ---
  private path(i: number): PathData | undefined {
    const p = this.obj.paths[i];
    return p && "inline" in p ? p.inline : undefined;
  }
  private grow(x: number, y: number, w: number, h: number) {
    if (this.ex.figures.length) this.ex.grow(boxOf(this.st.m, x, y, w, h), this.st.clip);
  }
  private growPath(i: number, dx = 0, dy = 0) {
    if (!this.ex.figures.length) return;
    const b = this.path(i) && pathBox(this.path(i)!);
    if (b) this.grow(b.x + dx, b.y + dy, b.w, b.h);
  }
  override clipRect(x: number, y: number, w: number, h: number) { this.st.clip = intersect(boxOf(this.st.m, x, y, w, h), this.st.clip) ?? { x: 0, y: 0, w: 0, h: 0 }; }
  override clipPath(path: number) {
    const b = this.path(path) && pathBox(this.path(path)!);
    if (b) this.st.clip = intersect(boxOf(this.st.m, b.x, b.y, b.w, b.h), this.st.clip) ?? { x: 0, y: 0, w: 0, h: 0 };
  }
  override fillRect(x: number, y: number, w: number, h: number) { this.grow(x, y, w, h); }
  override strokeRect(x: number, y: number, w: number, h: number) { this.grow(x, y, w, h); }
  override fillPath(path: number) { this.growPath(path); }
  override strokePath(path: number) { this.growPath(path); }
  override image(_img: number, x: number, y: number, w: number, h: number) { this.grow(x, y, w, h); }
  override imageSub(_img: number, _sx: number, _sy: number, _sw: number, _sh: number, dx: number, dy: number, dw: number, dh: number) { this.grow(dx, dy, dw, dh); }

  /** The drawing op right after ALT_TEXT renders that text; returns true when consumed. */
  private takeAlt(x: number, y: number, advance: number): boolean {
    if (this.ex.alt === undefined) return false;
    const t = this.ex.alt;
    this.ex.alt = undefined;
    this.ex.emit(t, x, y, advance, this.st, true);
    return true;
  }
  override fillText(text: string, x: number, y: number, advance: number) { if (!this.takeAlt(x, y, advance)) this.ex.emit(text, x, y, advance, this.st, false); }
  override strokeText(text: string, x: number, y: number, advance: number) { if (!this.takeAlt(x, y, advance)) this.ex.emit(text, x, y, advance, this.st, false); }
  override fillPathAt(path: number, _rule: number, x: number, y: number) {
    this.growPath(path, x, y);
    this.takeAlt(x, y, 0);
  }
  override fillPathRun(_rule: number, glyphs: Glyph[]) {
    for (const g of glyphs) this.growPath(g.path, g.x, g.y);
    this.takeAlt(glyphs[0]?.x ?? 0, glyphs[0]?.y ?? 0, 0);
  }
  override use(obj: number) { this.useAt(obj, 0, 0); }
  override useAt(obj: number, x: number, y: number) {
    if (this.ex.alt !== undefined) {
      // text drawn by the child: its extent still counts for a figure
      const child = this.ex.figures.length ? this.ex.resolve(this.obj.objects[obj]) : undefined;
      if (child) this.grow(child.bbox.x + x, child.bbox.y + y, child.bbox.w, child.bbox.h);
      this.takeAlt(x, y, 0);
      return;
    }
    const child = this.ex.resolve(this.obj.objects[obj]);
    if (!child) return;
    walk(child, new TextSink(this.ex, child, mul(this.st.m, [1, 0, 0, 1, x, y]), this.st.clip));
  }
}

function extract(obj: ObjectPart, resolve: (h: Hash) => ObjectPart | undefined, matrix: Matrix): Extraction {
  const ex = new Extraction(resolve);
  walk(obj, new TextSink(ex, obj, matrix));
  ex.flushAlt({ m: matrix, font: undefined, size: 10, align: 0 });
  ex.closeAll();
  return ex;
}

/**
 * Extract text runs from an object (and the objects it USEs) in the space
 * defined by matrix. resolve must return already-loaded child objects.
 * The order and ordinals match the Go extractor and the text index part.
 */
export function extractText(obj: ObjectPart, resolve: (h: Hash) => ObjectPart | undefined, matrix: Matrix = [1, 0, 0, 1, 0, 0]): TextRun[] {
  return extract(obj, resolve, matrix).runs;
}

/**
 * Extract the runs of an object together with its structure (headings,
 * lists, tables, figures; MARK kinds of spec §7.8) and its links, for
 * accessible text layers. The runs are the same as extractText's.
 */
export function extractContent(obj: ObjectPart, resolve: (h: Hash) => ObjectPart | undefined, matrix: Matrix = [1, 0, 0, 1, 0, 0]): TextContent {
  const ex = extract(obj, resolve, matrix);
  return { runs: ex.runs, nodes: ex.nodes, links: ex.links };
}
