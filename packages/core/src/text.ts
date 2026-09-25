import { NoopSink, walk } from "./object.js";
import type { ObjectPart, Font, Hash } from "./types.js";

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

class TextSink extends NoopSink {
  runs: TextRun[] = [];
  private stack: Matrix[] = [];
  private fontStack: { font: Font | undefined; size: number; align: number }[] = [];
  private cur: { font: Font | undefined; size: number; align: number } = { font: undefined, size: 10, align: 0 };

  constructor(private m: Matrix, private obj: ObjectPart, private resolve: (h: Hash) => ObjectPart | undefined) {
    super();
  }
  override save() { this.stack.push(this.m); this.fontStack.push({ ...this.cur }); }
  override restore() { this.m = this.stack.pop() ?? this.m; this.cur = this.fontStack.pop() ?? this.cur; }
  override transform(a: number, b: number, c: number, d: number, e: number, f: number) { this.m = mul(this.m, [a, b, c, d, e, f]); }
  override translate(x: number, y: number) { this.transform(1, 0, 0, 1, x, y); }
  override scale(x: number, y: number) { this.transform(x, 0, 0, y, 0, 0); }
  override font(font: number, size: number) { this.cur = { ...this.cur, font: this.obj.fonts[font], size }; }
  override textStyle(align: number) { this.cur = { ...this.cur, align }; }
  override fillText(text: string, x: number, y: number, advance: number) { this.emit(text, x, y, advance); }
  override strokeText(text: string, x: number, y: number, advance: number) { this.emit(text, x, y, advance); }
  private emit(text: string, x: number, y: number, advance: number) {
    const m = this.m;
    this.runs.push({
      text, advance, size: this.cur.size, font: this.cur.font, align: this.cur.align,
      x: m[0] * x + m[2] * y + m[4], y: m[1] * x + m[3] * y + m[5], matrix: m,
    });
  }
  override use(obj: number) { this.useAt(obj, 0, 0); }
  override useAt(obj: number, x: number, y: number) {
    const child = this.resolve(this.obj.objects[obj]);
    if (!child) return;
    const sub = new TextSink(mul(this.m, [1, 0, 0, 1, x, y]), child, this.resolve);
    walk(child, sub);
    this.runs.push(...sub.runs);
  }
}

/**
 * Extract text runs from an object (and the objects it USEs) in the space
 * defined by matrix. resolve must return already-loaded child objects.
 */
export function extractText(obj: ObjectPart, resolve: (h: Hash) => ObjectPart | undefined, matrix: Matrix = [1, 0, 0, 1, 0, 0]): TextRun[] {
  const sink = new TextSink(matrix, obj, resolve);
  walk(obj, sink);
  return sink.runs;
}
