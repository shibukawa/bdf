import { BdfFormatError } from "./bytes.js";
import { checkHash } from "./container.js";
import type { Manifest, View, RectDef } from "./types.js";

// The numbers of a manifest decide how long the loops of a reader run and
// how large what it allocates is (the tiles of a region, the gridlines of a
// sheet, the bands of a continuous layout), so they are checked once, when
// the document opens. What is refused here no writer has a use for.

/**
 * The sides of pages and of their bodies are this many units at most: object
 * coordinates are f32, which cannot tell whole units apart beyond it.
 */
export const MAX_PAGE_SIZE = 1 << 24;
/** A sheet has this many rows, and this many columns, at most (Excel has 1,048,576 rows). */
export const MAX_SHEET_ENTRIES = 1 << 24;
/** Rows and columns with a size (0: hidden) are this many units at least. */
export const MIN_ENTRY_SIZE = 1 / 64;
/** The tile size of a sheet view that states none. */
export const DEFAULT_TILE = 2048;

/** The tile size of a sheet view; a size of 0 or less is the default one, as the Go reader has it. */
export function tileSize(view: View): number {
  return view.tile !== undefined && view.tile > 0 ? view.tile : DEFAULT_TILE;
}

const isNumber = (v: unknown): v is number => typeof v === "number" && Number.isFinite(v);
const isCount = (v: unknown): v is number => Number.isInteger(v) && (v as number) >= 0;

function checkSize(v: unknown, what: string): void {
  if (!isNumber(v) || v < 0 || v > MAX_PAGE_SIZE) throw new BdfFormatError(`${what} ${JSON.stringify(v)} out of range`);
}

function checkRect(r: RectDef, what: string): void {
  for (const v of [r.x, r.y]) {
    if (!isNumber(v) || Math.abs(v) > MAX_PAGE_SIZE) throw new BdfFormatError(`${what} at ${JSON.stringify(v)} out of range`);
  }
  checkSize(r.w, `${what} width`);
  checkSize(r.h, `${what} height`);
}

function checkRuns(runs: unknown, what: string): void {
  if (runs === undefined) return;
  if (!Array.isArray(runs)) throw new BdfFormatError(`${what} are not a list`);
  let total = 0;
  for (const run of runs) {
    if (!Array.isArray(run)) throw new BdfFormatError(`${what}: a run is not a pair`);
    const [n, size] = run;
    if (!isCount(n) || (total += n) > MAX_SHEET_ENTRIES) throw new BdfFormatError(`more than ${MAX_SHEET_ENTRIES} ${what}`);
    if (!isNumber(size) || size < 0 || size > MAX_PAGE_SIZE || (size > 0 && size < MIN_ENTRY_SIZE)) throw new BdfFormatError(`${what} of size ${JSON.stringify(size)}`);
  }
}

function checkView(v: View): void {
  if (v.pages !== undefined) {
    if (!Array.isArray(v.pages)) throw new BdfFormatError(`view ${v.id}: pages are not a list`);
    for (const p of v.pages) {
      checkSize(p?.w, "page width");
      checkSize(p.h, "page height");
      if (p.body !== undefined) checkRect(p.body, "page body");
    }
  }
  if (v.continuous !== undefined && (!isNumber(v.continuous?.gap) || Math.abs(v.continuous.gap) > MAX_PAGE_SIZE)) {
    throw new BdfFormatError(`view ${v.id}: bad gap`);
  }
  if (v.tile !== undefined && (!isNumber(v.tile) || (v.tile > 0 && v.tile < 1))) throw new BdfFormatError(`tile size ${JSON.stringify(v.tile)} out of range`);
  checkRuns(v.cols, "columns");
  checkRuns(v.rows, "rows");
  if (v.freeze !== undefined) {
    for (const n of [v.freeze?.cols ?? 0, v.freeze?.rows ?? 0]) if (!isCount(n)) throw new BdfFormatError(`view ${v.id}: bad frozen panes`);
  }
}

/**
 * Check the part names and the numbers of the views of a manifest (not the
 * outer manifest of an encrypted document, which has no views); throws
 * BdfFormatError.
 */
export function checkManifest(m: Manifest): void {
  if (!Array.isArray(m.parts)) throw new BdfFormatError("manifest without parts");
  for (const e of m.parts) {
    checkHash(e?.h);
    if (e.sealed !== undefined) checkHash(e.sealed, "sealed part name");
  }
  if (!Array.isArray(m.views)) throw new BdfFormatError("manifest without views");
  for (const v of m.views) {
    if (typeof v?.id !== "string") throw new BdfFormatError("view without an id");
    checkView(v);
  }
}
