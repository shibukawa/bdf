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
  // Three references, a separator and a string length take at least five
  // bytes per run, even when every reference and string is empty.
  r.need(n * 5);
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

/**
 * Fold normalized text for matching, a character at a time: case, and small
 * and large kana alike (katakana to hiragana). The minus sign of formulas
 * matches the hyphen-minus that queries are typed with.
 */
function fold(s: string, caseSensitive: boolean): string {
  let out = "";
  for (const c of s) {
    if (c === "\u2212") {
      out += "-";
      continue;
    }
    for (const d of caseSensitive ? c : c.toLowerCase()) {
      const cp = d.codePointAt(0)!;
      // Katakana to hiragana so that both scripts match each other.
      out += cp >= 0x30a1 && cp <= 0x30f6 ? String.fromCodePoint(cp - 0x60) : d;
    }
  }
  return out;
}

/**
 * Normalize a string for matching: NFKC of the string as a whole, so that a
 * letter and the marks after it are the letter they compose (e and a
 * combining acute are é, a half-width ka and its voiced mark are が), then
 * case and kana folded.
 */
export function normalizeQuery(q: string, caseSensitive = false): string {
  return fold(q.normalize("NFKC"), caseSensitive);
}

/** A string that starts with a combining mark: it belongs to what is before it. */
const MARK = /^\p{M}/u;

/**
 * Splits text into the pieces that NFKC normalizes on their own, and calls
 * back with each: where it starts and ends in the text, in code units, and
 * what it normalizes to. The pieces normalized one after another are the
 * text normalized, and they are as small as that allows: a character alone,
 * unless it changes with the ones before it (marks, the letters of a Hangul
 * syllable).
 */
function pieces(text: string, each: (start: number, end: number, norm: string) => void): void {
  let start = 0;
  let off = 0;
  let piece = "";
  let norm = "";
  for (const ch of text) {
    // (ASCII is as it is, and changes nothing before it)
    const ascii = ch.length === 1 && ch.charCodeAt(0) < 0x80;
    const n = ascii ? ch : ch.normalize("NFKC");
    if (piece !== "") {
      if (!ascii && (MARK.test(n) || (piece + ch).normalize("NFKC") !== norm + n)) {
        piece += ch;
        norm = piece.normalize("NFKC");
        off += ch.length;
        continue;
      }
      each(start, off, norm);
    }
    start = off;
    piece = ch;
    norm = n;
    off += ch.length;
  }
  if (piece !== "") each(start, off, norm);
}

/**
 * Searchable text built from index runs. The normalized string keeps a map
 * back to (run, offset) so hits can be located in the original runs.
 */
export class TextSearch {
  private norm = "";
  /**
   * For each code unit of norm: the run it is of, and where the piece of the
   * run's text it was normalized from starts and ends, in code units (-1 for
   * separators).
   */
  private runOf: Int32Array;
  private offOf: Int32Array;
  private endOf: Int32Array;
  /** The plain text, once it was asked for (every hit takes its context from it). */
  private plain: string | undefined;

  constructor(readonly runs: IndexRun[], readonly caseSensitive = false) {
    const runIdx: number[] = [];
    const offIdx: number[] = [];
    const endIdx: number[] = [];
    let norm = "";
    runs.forEach((run, ri) => {
      if (ri > 0 && run.sep !== Sep.NONE) {
        // Paragraph breaks become a character no query contains; spaces stay spaces.
        norm += run.sep === Sep.BREAK ? "\u0000" : " ";
        runIdx.push(-1);
        offIdx.push(-1);
        endIdx.push(-1);
      }
      pieces(run.text, (start, end, piece) => {
        const n = fold(piece, caseSensitive);
        for (let k = 0; k < n.length; k++) {
          runIdx.push(ri);
          offIdx.push(start);
          endIdx.push(end);
        }
        norm += n;
      });
    });
    this.norm = norm;
    this.runOf = Int32Array.from(runIdx);
    this.offOf = Int32Array.from(offIdx);
    this.endOf = Int32Array.from(endIdx);
  }

  /** The searchable plain text (paragraph breaks as newlines). */
  get text(): string { return (this.plain ??= this.norm.replaceAll("\u0000", "\n")); }

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
      const end = this.endOf[i];
      const run = this.runs[ri];
      const last = segments[segments.length - 1];
      if (last && last.run === ri) last.end = Math.max(last.end, end);
      else segments.push({ run: ri, a: run.a, b: run.b, ordinal: run.ordinal, start: off, end });
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
