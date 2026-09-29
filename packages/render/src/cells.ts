// Tables as cells: the text of each cell of a table on a page (MARK TABLE)
// or of a sheet (cells outside tables), and a rectangle of them in the two
// forms spreadsheets paste as cells: tab-separated values and an HTML table.
import { Sep, type TextContent, type TextRun } from "@bdfkit/core";

/** A cell and its text. Positions are 0-based: relative to the table, or absolute in a sheet. */
export interface CellText {
  row: number;
  col: number;
  /** Rows and columns a merged cell spans (1 for others). */
  rows: number;
  cols: number;
  /** The cell's runs joined as their separators say: nothing, a space, or a line break between paragraphs. */
  text: string;
  /** Header cell. */
  scope?: "col" | "row";
}

/** A rectangle of cells: its first and last rows and columns, included. */
export interface CellRange { row0: number; col0: number; row1: number; col1: number }

/** Cells as a spreadsheet takes them from the clipboard: text/plain and text/html. */
export interface CellClipboard { text: string; html: string }

export interface CellClipboardOptions {
  /** Leave out the rows and columns after the last cell with text (whole rows or columns selected). */
  trim?: boolean;
  /** Rows to leave out (hidden ones). */
  skipRow?: (row: number) => boolean;
  /** Columns to leave out (hidden ones). */
  skipCol?: (col: number) => boolean;
}

/**
 * The cells of a table of a text content (table is the index of its table
 * node), or with table -1 the cells outside tables (the cells of a sheet),
 * sorted by position. A cell that appears more than once (a sheet cell
 * drawn by two tiles) is one cell. keep chooses the runs (default: those
 * with text).
 */
export function tableCells(content: TextContent, table = -1, keep: (run: TextRun) => boolean = (r) => !!r.text): CellText[] {
  const { nodes, runs } = content;
  const isCell = (i: number) => {
    const n = nodes[i];
    if (n.kind !== "cell" || n.row === undefined || n.col === undefined) return false;
    if (table >= 0) return n.parent === table;
    for (let p = n.parent; p >= 0; p = nodes[p].parent) if (nodes[p].kind === "table") return false;
    return true;
  };
  // the cell each node lies in (-1: none)
  const owner = new Map<number, number>();
  const cellOf = (i: number | undefined): number => {
    if (i === undefined || i < 0 || i >= nodes.length) return -1;
    let c = owner.get(i);
    if (c === undefined) {
      c = isCell(i) ? i : cellOf(nodes[i].parent);
      owner.set(i, c);
    }
    return c;
  };
  const texts = new Map<number, string>();
  for (const r of runs) {
    if (!keep(r)) continue;
    const c = cellOf(r.node);
    if (c < 0) continue;
    const t = texts.get(c);
    texts.set(c, t === undefined ? r.text : t + (r.sep === Sep.BREAK ? "\n" : r.sep === Sep.SPACE ? " " : "") + r.text);
  }
  const out = new Map<string, CellText>();
  nodes.forEach((n, i) => {
    if (!isCell(i)) return;
    const key = `${n.row},${n.col}`;
    const text = texts.get(i) ?? "";
    const seen = out.get(key);
    if (seen) {
      // parts of one cell from two tiles: their separator is not known
      if (text) seen.text = seen.text ? `${seen.text} ${text}` : text;
      return;
    }
    const cell: CellText = { row: n.row!, col: n.col!, rows: n.rows ?? 1, cols: n.cols ?? 1, text };
    if (n.scope) cell.scope = n.scope;
    out.set(key, cell);
  });
  return [...out.values()].sort((a, b) => a.row - b.row || a.col - b.col);
}

/** A field of tab-separated values: quoted (quotes doubled) when it holds a tab, a line break or a quote. */
function field(s: string): string {
  return /[\t\n\r"]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s;
}

function escapeHTML(s: string): string {
  return s.replace(/[&<>"]/g, (c) => (c === "&" ? "&amp;" : c === "<" ? "&lt;" : c === ">" ? "&gt;" : "&quot;"));
}

/** The first index of a sorted list at or after v. */
function lowerBound(list: number[], v: number): number {
  let lo = 0, hi = list.length;
  while (lo < hi) {
    const mid = (lo + hi) >> 1;
    if (list[mid] < v) lo = mid + 1;
    else hi = mid;
  }
  return lo;
}

/**
 * A rectangle of more cells than this is not put on the clipboard, unless it
 * has a 64th as many cells with something in them: the references of two
 * cells can be any number of rows apart.
 */
export const MAX_CLIPBOARD_CELLS = 1 << 22;

/**
 * The cells in range as tab-separated values (rows end with a line break
 * but the last; a field holding a tab, a line break or a quote is quoted, as
 * Excel copies it) and as an HTML table (merged cells span their rows and
 * columns, header cells are th, line breaks stay in the cell). A merged cell
 * that reaches into the range is put at its first row and column inside it,
 * the others it covers are empty. Throws for a rectangle of too many cells
 * (MAX_CLIPBOARD_CELLS).
 */
export function cellClipboard(cells: CellText[], range: CellRange, opts: CellClipboardOptions = {}): CellClipboard {
  const { row0, col0 } = range;
  let { row1, col1 } = range;
  type Clipped = { cell: CellText; r0: number; c0: number; r1: number; c1: number };
  const inside: Clipped[] = [];
  let lastRow = row0 - 1, lastCol = col0 - 1;
  for (const cell of cells) {
    const r0 = Math.max(cell.row, row0), c0 = Math.max(cell.col, col0);
    const r1 = Math.min(cell.row + cell.rows - 1, row1), c1 = Math.min(cell.col + cell.cols - 1, col1);
    if (r0 > r1 || c0 > c1) continue;
    inside.push({ cell, r0, c0, r1, c1 });
    if (cell.text) {
      lastRow = Math.max(lastRow, r0);
      lastCol = Math.max(lastCol, c0);
    }
  }
  if (opts.trim) {
    row1 = Math.min(row1, lastRow);
    col1 = Math.min(col1, lastCol);
  }
  const limit = Math.max(MAX_CLIPBOARD_CELLS, 64 * inside.length);
  const tooMany = () => new Error("bdf: too many cells to copy");
  if (!(Math.max(0, row1 - row0 + 1) * Math.max(0, col1 - col0 + 1) <= limit)) throw tooMany();
  const rows: number[] = [], cols: number[] = [];
  for (let r = row0; r <= row1; r++) if (!opts.skipRow?.(r)) rows.push(r);
  for (let c = col0; c <= col1; c++) if (!opts.skipCol?.(c)) cols.push(c);
  // each cell at its first row and column shown, spanning those shown
  const n = cols.length;
  const at = new Map<number, { cell: CellText; rows: number; cols: number }>();
  const covered = new Set<number>();
  /** Rows with a cell, or a part of one. */
  const used = new Uint8Array(rows.length);
  let spanned = 0;
  for (const { cell, r0, c0, r1, c1 } of inside) {
    const ri = lowerBound(rows, r0), ci = lowerBound(cols, c0);
    const rn = lowerBound(rows, Math.min(r1, row1) + 1) - ri, cn = lowerBound(cols, Math.min(c1, col1) + 1) - ci;
    if (rn <= 0 || cn <= 0) continue;
    // cells do not overlap: those that do may not take longer than the others
    if ((spanned += rn * cn) > limit) throw tooMany();
    at.set(ri * n + ci, { cell, rows: rn, cols: cn });
    for (let i = 0; i < rn; i++) {
      used[ri + i] = 1;
      for (let j = 0; j < cn; j++) if (i || j) covered.add((ri + i) * n + ci + j);
    }
  }
  const emptyLine = "\t".repeat(Math.max(0, n - 1)), emptyRow = `<tr>${"<td></td>".repeat(n)}</tr>`;
  const lines: string[] = [];
  let html = '<meta charset="utf-8"><table>';
  for (let i = 0; i < rows.length; i++) {
    if (!used[i]) {
      html += emptyRow;
      lines.push(emptyLine);
      continue;
    }
    const fields: string[] = [];
    html += "<tr>";
    for (let j = 0; j < n; j++) {
      const got = at.get(i * n + j);
      fields.push(got ? field(got.cell.text) : "");
      if (covered.has(i * n + j)) continue;
      if (!got) {
        html += "<td></td>";
        continue;
      }
      const tag = got.cell.scope ? "th" : "td";
      let attrs = got.cell.scope ? ` scope="${got.cell.scope}"` : "";
      if (got.rows > 1) attrs += ` rowspan="${got.rows}"`;
      if (got.cols > 1) attrs += ` colspan="${got.cols}"`;
      // Excel starts a new row at <br> unless it is told to stay in the cell
      html += `<${tag}${attrs}>${escapeHTML(got.cell.text).replace(/\n/g, '<br style="mso-data-placement:same-cell">')}</${tag}>`;
    }
    html += "</tr>";
    lines.push(fields.join("\t"));
  }
  html += "</table>";
  return { text: lines.join("\n"), html };
}
