import { NoopSink, walk } from "./object.js";
import type { ObjectPart, Font, Hash } from "./types.js";

/** MARK kinds (docs/spec.md §7.8). */
export const Mark = { PARAGRAPH: 0, LINE: 1, CELL: 2, BOX: 3, ALT_TEXT: 4, WRAP: 5 } as const;

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

interface State { m: Matrix; font: Font | undefined; size: number; align: number }

/** Shared run list and MARK state across nested objects. */
class Extraction {
  runs: TextRun[] = [];
  pending = Sep.NONE as number;
  hasMark = false;
  /** ALT_TEXT waiting for the drawing op it describes. */
  alt: string | undefined;
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
    this.runs.push(run);
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

  constructor(private ex: Extraction, private obj: ObjectPart, m: Matrix) {
    super();
    this.st = { m, font: undefined, size: 10, align: 0 };
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
      case Mark.PARAGRAPH: case Mark.CELL: case Mark.BOX: this.ex.flushAlt(this.st); this.ex.mark(Sep.BREAK); break;
      case Mark.ALT_TEXT: this.ex.flushAlt(this.st); this.ex.alt = payload; break;
      case Mark.WRAP: this.ex.flushAlt(this.st); this.ex.mark(Sep.NONE); break;
    }
  }
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
  override fillPathAt(_path: number, _rule: number, x: number, y: number) { this.takeAlt(x, y, 0); }
  override fillPathRun(_rule: number, glyphs: { x: number; y: number }[]) { this.takeAlt(glyphs[0]?.x ?? 0, glyphs[0]?.y ?? 0, 0); }
  override use(obj: number) { this.useAt(obj, 0, 0); }
  override useAt(obj: number, x: number, y: number) {
    if (this.takeAlt(x, y, 0)) return;
    const child = this.ex.resolve(this.obj.objects[obj]);
    if (!child) return;
    walk(child, new TextSink(this.ex, child, mul(this.st.m, [1, 0, 0, 1, x, y])));
  }
}

/**
 * Extract text runs from an object (and the objects it USEs) in the space
 * defined by matrix. resolve must return already-loaded child objects.
 * The order and ordinals match the Go extractor and the text index part.
 */
export function extractText(obj: ObjectPart, resolve: (h: Hash) => ObjectPart | undefined, matrix: Matrix = [1, 0, 0, 1, 0, 0]): TextRun[] {
  const ex = new Extraction(resolve);
  walk(obj, new TextSink(ex, obj, matrix));
  ex.flushAlt({ m: matrix, font: undefined, size: 10, align: 0 });
  return ex.runs;
}
