import { ByteReader, BdfFormatError } from "./bytes.js";
import { Sep } from "./text.js";

/** One entry of a text index part (docs/spec.md §7.9). */
export interface IndexRun {
  /** Page and layer (fixed/flow) or tile x and y (sheet). */
  a: number;
  b: number;
  ordinal: number;
  sep: number;
  text: string;
}

export function decodeTextIndex(bytes: Uint8Array): IndexRun[] {
  const r = new ByteReader(bytes);
  const m = r.bytesN(4);
  if (m[0] !== 0x42 || m[1] !== 0x54 || m[2] !== 0x58 || m[3] !== 0x54) throw new BdfFormatError("not a text index part");
  const version = r.u16();
  if (version > 1) throw new BdfFormatError(`unsupported text index version ${version}`);
  const n = r.varuint();
  const out: IndexRun[] = new Array(n);
  for (let i = 0; i < n; i++) {
    out[i] = { a: r.varuint(), b: r.varuint(), ordinal: r.varuint(), sep: r.u8(), text: r.str() };
  }
  return out;
}

/** Part of a hit inside one run: [start, end) are code-unit offsets in the run's text. */
export interface HitSegment { run: number; a: number; b: number; ordinal: number; start: number; end: number }

/** A search hit; spans one or more consecutive runs. */
export interface SearchHit {
  segments: HitSegment[];
  /** Matched text as it appears in the document (unnormalized). */
  text: string;
  /** Some context around the hit, for result lists. */
  context: string;
}

export interface SearchOptions {
  /** Maximum number of hits (default 1000). */
  limit?: number;
  /** Case-sensitive matching (default false). */
  caseSensitive?: boolean;
  /** Context characters around each hit (default 40). */
  context?: number;
}

/** Normalize one character for matching: NFKC, case folding, and small kana to full-size kana are treated alike. */
export function normalizeChar(ch: string, caseSensitive: boolean): string {
  let s = ch.normalize("NFKC");
  if (!caseSensitive) s = s.toLowerCase();
  // Katakana to hiragana so that both scripts match each other.
  let out = "";
  for (const c of s) {
    const cp = c.codePointAt(0)!;
    out += cp >= 0x30a1 && cp <= 0x30f6 ? String.fromCodePoint(cp - 0x60) : c;
  }
  return out;
}

export function normalizeQuery(q: string, caseSensitive = false): string {
  let out = "";
  for (const c of q) out += normalizeChar(c, caseSensitive);
  return out;
}

/**
 * Searchable text built from index runs. The normalized string keeps a map
 * back to (run, offset) so hits can be located in the original runs.
 */
export class TextSearch {
  private norm = "";
  /** For each code unit of norm: run index and code-unit offset in that run's text (or -1 for separators). */
  private runOf: Int32Array;
  private offOf: Int32Array;

  constructor(readonly runs: IndexRun[], readonly caseSensitive = false) {
    const runIdx: number[] = [];
    const offIdx: number[] = [];
    let norm = "";
    runs.forEach((run, ri) => {
      if (ri > 0 && run.sep !== Sep.NONE) {
        // Paragraph breaks become a character no query contains; spaces stay spaces.
        norm += run.sep === Sep.BREAK ? "\u0000" : " ";
        runIdx.push(-1);
        offIdx.push(-1);
      }
      let off = 0;
      for (const ch of run.text) {
        const n = normalizeChar(ch, caseSensitive);
        for (let k = 0; k < n.length; k++) {
          runIdx.push(ri);
          offIdx.push(off);
        }
        norm += n;
        off += ch.length;
      }
    });
    this.norm = norm;
    this.runOf = Int32Array.from(runIdx);
    this.offOf = Int32Array.from(offIdx);
  }

  /** The searchable plain text (paragraph breaks as newlines). */
  get text(): string { return this.norm.replaceAll("\u0000", "\n"); }

  search(query: string, opts: SearchOptions = {}): SearchHit[] {
    const q = normalizeQuery(query, this.caseSensitive).replace(/\s+/g, " ");
    if (!q) return [];
    const limit = opts.limit ?? 1000;
    const ctxLen = opts.context ?? 40;
    const hits: SearchHit[] = [];
    let from = 0;
    while (hits.length < limit) {
      const at = this.norm.indexOf(q, from);
      if (at < 0) break;
      from = at + 1;
      const hit = this.locate(at, at + q.length, ctxLen);
      if (hit) hits.push(hit);
    }
    return hits;
  }

  private locate(start: number, end: number, ctxLen: number): SearchHit | undefined {
    const segments: HitSegment[] = [];
    for (let i = start; i < end; i++) {
      const ri = this.runOf[i];
      if (ri < 0) continue;
      const off = this.offOf[i];
      const run = this.runs[ri];
      const charEnd = off + ((run.text.codePointAt(off) ?? 0) > 0xffff ? 2 : 1);
      const last = segments[segments.length - 1];
      if (last && last.run === ri) last.end = Math.max(last.end, charEnd);
      else segments.push({ run: ri, a: run.a, b: run.b, ordinal: run.ordinal, start: off, end: charEnd });
    }
    if (!segments.length) return undefined;
    let text = "";
    for (const s of segments) {
      const run = this.runs[s.run];
      text += (text && run.sep === Sep.SPACE ? " " : "") + run.text.slice(s.start, s.end);
    }
    const plain = this.text;
    const context = plain.slice(Math.max(0, start - ctxLen), Math.min(plain.length, end + ctxLen)).replace(/\n/g, " ⏎ ");
    return { segments, text, context };
  }
}
