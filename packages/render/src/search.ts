import { type BdfDocument, type View, type Rect, type Hash, type TextRun, TextSearch, extractText, type SearchHit, type SearchOptions } from "@bdf/core";
import { fontString } from "./resources.js";

/** A highlight rectangle for part of a hit, in view coordinates. */
export interface HitRect extends Rect {
  /** Page index (fixed/flow) or tile x (sheet). */
  a: number;
  /** Layer index (fixed/flow) or tile y (sheet). */
  b: number;
}

/** Measures text width for a CSS font string (a 2D context's measureText). */
export type Measure = (font: string, text: string) => number;

/**
 * Search inside a document's views and locate hits as rectangles.
 * Hit location needs the object that drew the text, so it is lazy.
 */
export class DocumentSearch {
  private indexes = new Map<string, Promise<TextSearch>>();
  private runs = new Map<Hash, TextRun[]>();

  constructor(readonly doc: BdfDocument, readonly measure: Measure) {}

  private index(view: View): Promise<TextSearch> {
    let p = this.indexes.get(view.id);
    if (!p) {
      p = this.doc.textIndex(view).then((runs) => new TextSearch(runs));
      this.indexes.set(view.id, p);
    }
    return p;
  }

  async search(view: View, query: string, opts?: SearchOptions): Promise<SearchHit[]> {
    return (await this.index(view)).search(query, opts);
  }

  /** The searchable plain text of a view. */
  async text(view: View): Promise<string> {
    return (await this.index(view)).text;
  }

  private objectFor(view: View, a: number, b: number): Hash | undefined {
    if (view.kind === "sheet") return view.tiles?.[`${a},${b}`];
    return view.pages?.[a]?.layers[b]?.obj;
  }

  private async runsOf(hash: Hash): Promise<TextRun[]> {
    let runs = this.runs.get(hash);
    if (!runs) {
      const obj = await this.doc.ensure(hash);
      runs = extractText(obj, (h) => this.doc.objectSync(h));
      this.runs.set(hash, runs);
    }
    return runs;
  }

  /** Rectangles covering a hit, in page coordinates (or sheet coordinates for sheets). */
  async locate(view: View, hit: SearchHit): Promise<HitRect[]> {
    const out: HitRect[] = [];
    const tile = view.tile ?? 2048;
    for (const seg of hit.segments) {
      const hash = this.objectFor(view, seg.a, seg.b);
      if (!hash) continue;
      const run = (await this.runsOf(hash))[seg.ordinal];
      if (!run || run.altText || !run.font) continue;
      const r = runRect(run, seg.start, seg.end, this.measure);
      if (!r) continue;
      if (view.kind === "sheet") {
        r.x += seg.a * tile;
        r.y += seg.b * tile;
      }
      out.push({ ...r, a: seg.a, b: seg.b });
    }
    return out;
  }
}

/** Bounding box of text[start:end] of a run, in the run's coordinate space. */
export function runRect(run: TextRun, start: number, end: number, measure: Measure): Rect | undefined {
  if (!run.font) return undefined;
  const font = fontString(run.font, run.size);
  const full = measure(font, run.text);
  const width = run.advance > 0 ? run.advance : full;
  const scale = run.advance > 0 && full > 0 ? run.advance / full : 1;
  const w1 = measure(font, run.text.slice(0, start)) * scale;
  const w2 = measure(font, run.text.slice(0, end)) * scale;
  const anchor = run.align === 1 ? -width : run.align === 2 ? -width / 2 : 0;
  const m = run.matrix;
  // Local offsets relative to the anchor, before the run's transform.
  const xs = [anchor + w1, anchor + w2];
  const ys = [-run.size * 0.8, run.size * 0.2];
  let x0 = Infinity, y0 = Infinity, x1 = -Infinity, y1 = -Infinity;
  for (const dx of xs) for (const dy of ys) {
    const px = m[0] * dx + m[2] * dy + run.x;
    const py = m[1] * dx + m[3] * dy + run.y;
    x0 = Math.min(x0, px); y0 = Math.min(y0, py); x1 = Math.max(x1, px); y1 = Math.max(y1, py);
  }
  return { x: x0, y: y0, w: x1 - x0, h: y1 - y0 };
}
