import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { BufferSource, BdfDocument, extractContent } from "@bdfkit/core";
import { tableCells, cellClipboard, MAX_CLIPBOARD_CELLS } from "../dist/cells.js";

const root = new URL("../../../", import.meta.url);
const doc = await BdfDocument.open(new BufferSource(new Uint8Array(await readFile(new URL("testdata/demo.bdf", root)))));

/** Content of an object and the objects it uses. */
async function content(hash, matrix) {
  const obj = await doc.ensure(hash);
  return extractContent(obj, (h) => doc.objectSync(h), matrix);
}

// Sep values (docs/spec.md §7.9): NONE 0, SPACE 1, BREAK 2.
const run = (text, node, sep = 2) => ({ text, node, sep, x: 0, y: 0, advance: 0, size: 10, font: undefined, align: 0, matrix: [1, 0, 0, 1, 0, 0], ordinal: 0, altText: false });

test("tableCells reads the cells of a table on a page", async () => {
  const page = doc.view("doc").pages[0];
  const c = await content(page.layers.find((l) => l.role === "body").obj);
  const table = c.nodes.findIndex((n) => n.kind === "table");
  assert.ok(table >= 0);
  const cells = tableCells(c, table);
  assert.equal(cells.length, 12);
  assert.deepEqual(cells[0], { row: 0, col: 0, rows: 1, cols: 1, text: "R1C1", scope: "col" });
  assert.deepEqual(cells[11], { row: 3, col: 2, rows: 1, cols: 1, text: "R4C3" });
  // outside tables: none on this page
  assert.deepEqual(tableCells(c), []);
  const { text, html } = cellClipboard(cells, { row0: 0, col0: 1, row1: 1, col1: 2 });
  assert.equal(text, "R1C2\tR1C3\nR2C2\tR2C3");
  assert.equal(html, '<meta charset="utf-8"><table><tr><th scope="col">R1C2</th><th scope="col">R1C3</th></tr><tr><td>R2C2</td><td>R2C3</td></tr></table>');
});

test("tableCells reads the cells of a sheet", async () => {
  const v = doc.view("sheet1");
  const c = await content(v.tiles["0,0"], [1, 0, 0, 1, 0, 0]);
  const cells = tableCells(c);
  const at = (r, col) => cells.find((x) => x.row === r && x.col === col)?.text;
  assert.equal(at(0, 0), "Col A");
  assert.equal(at(0, 2), "Col C");
  assert.equal(at(4, 0), "Row 5");
  assert.equal(at(4, 2), "15");
  assert.equal(cellClipboard(cells, { row0: 0, col0: 0, row1: 2, col1: 2 }).text, "Col A\tCol B\tCol C\nRow 2\t4\t6\nRow 3\t6\t9");
});

test("tableCells joins a cell's runs and puts a cell drawn twice together", () => {
  const nodes = [
    { kind: "table", parent: -1 },
    { kind: "cell", parent: 0, row: 0, col: 0, rows: 1, cols: 1 },
    { kind: "paragraph", parent: 1 },
    { kind: "paragraph", parent: 1 },
    { kind: "cell", parent: 0, row: 0, col: 1, rows: 1, cols: 1 },
    { kind: "paragraph", parent: 4 },
  ];
  const runs = [run("first", 2), run("line", 2, 1), run("wrap", 2, 0), run("", 2, 1), run("second", 3), run("B", 5)];
  assert.deepEqual(tableCells({ runs, nodes, links: [] }, 0).map((c) => c.text), ["first linewrap\nsecond", "B"]);

  // a sheet: cells outside tables; the same cell from two tiles is one
  const sheet = [
    { kind: "cell", parent: -1, row: 3, col: 1, rows: 2, cols: 1 },
    { kind: "cell", parent: -1, row: 0, col: 0, rows: 1, cols: 1, scope: "col" },
    { kind: "cell", parent: -1, row: 3, col: 1, rows: 2, cols: 1 },
    { kind: "table", parent: -1 },
    { kind: "cell", parent: 3, row: 0, col: 0, rows: 1, cols: 1 },
  ];
  const cells = tableCells({ runs: [run("top", 0), run("head", 1), run("bottom", 2), run("in a table", 4)], nodes: sheet, links: [] });
  assert.deepEqual(cells, [
    { row: 0, col: 0, rows: 1, cols: 1, text: "head", scope: "col" },
    { row: 3, col: 1, rows: 2, cols: 1, text: "top bottom" },
  ]);
});

test("cellClipboard quotes fields as Excel does and keeps empty cells", () => {
  const cells = [
    { row: 0, col: 0, rows: 1, cols: 1, text: "plain" },
    { row: 0, col: 2, rows: 1, cols: 1, text: 'say "hi"' },
    { row: 1, col: 1, rows: 1, cols: 1, text: "two\nlines" },
    { row: 1, col: 2, rows: 1, cols: 1, text: "a\tb" },
    { row: 2, col: 0, rows: 1, cols: 1, text: "<b>&amp;</b>" },
  ];
  const { text, html } = cellClipboard(cells, { row0: 0, col0: 0, row1: 2, col1: 2 });
  assert.equal(text, 'plain\t\t"say ""hi"""\n\t"two\nlines"\t"a\tb"\n<b>&amp;</b>\t\t');
  assert.match(html, /<td>two<br style="mso-data-placement:same-cell">lines<\/td>/);
  assert.match(html, /<td>&lt;b&gt;&amp;amp;&lt;\/b&gt;<\/td><td><\/td><td><\/td>/);
  assert.match(html, /<td>say &quot;hi&quot;<\/td>/);
});

test("cellClipboard puts a merged cell at its first position in the range and spans it", () => {
  const cells = [
    { row: 0, col: 0, rows: 2, cols: 2, text: "merged" },
    { row: 0, col: 2, rows: 1, cols: 1, text: "C1" },
    { row: 1, col: 2, rows: 1, cols: 1, text: "C2" },
    { row: 2, col: 0, rows: 1, cols: 3, text: "wide" },
  ];
  const all = cellClipboard(cells, { row0: 0, col0: 0, row1: 2, col1: 2 });
  assert.equal(all.text, "merged\t\tC1\n\t\tC2\nwide\t\t");
  assert.equal(all.html, '<meta charset="utf-8"><table><tr><td rowspan="2" colspan="2">merged</td><td>C1</td></tr><tr><td>C2</td></tr><tr><td colspan="3">wide</td></tr></table>');
  // a range that starts inside the merged cells
  const part = cellClipboard(cells, { row0: 1, col0: 1, row1: 2, col1: 2 });
  assert.equal(part.text, "merged\tC2\nwide\t");
  assert.equal(part.html, '<meta charset="utf-8"><table><tr><td>merged</td><td>C2</td></tr><tr><td colspan="2">wide</td></tr></table>');
});

test("cellClipboard leaves out hidden rows and columns, and trims whole rows and columns", () => {
  const cells = [
    { row: 0, col: 0, rows: 1, cols: 1, text: "A1" },
    { row: 1, col: 0, rows: 1, cols: 1, text: "hidden row" },
    { row: 2, col: 0, rows: 1, cols: 3, text: "spans a hidden column" },
    { row: 0, col: 2, rows: 1, cols: 1, text: "C1" },
    { row: 0, col: 1, rows: 1, cols: 1, text: "hidden column" },
  ];
  const got = cellClipboard(cells, { row0: 0, col0: 0, row1: 1_048_575, col1: 16_383 }, { trim: true, skipRow: (r) => r === 1, skipCol: (c) => c === 1 });
  assert.equal(got.text, "A1\tC1\nspans a hidden column\t");
  assert.match(got.html, /<td colspan="2">spans a hidden column<\/td>/);
  // nothing with text: nothing
  assert.deepEqual(cellClipboard([], { row0: 5, col0: 5, row1: 9, col1: 9 }, { trim: true }), { text: "", html: '<meta charset="utf-8"><table></table>' });
  // without trim, empty rows and columns stay
  assert.equal(cellClipboard([], { row0: 0, col0: 0, row1: 1, col1: 1 }).text, "\t\n\t");
});

/** cellClipboard as it was written first: every position of the rectangle looked up by its name. */
function reference(cells, range, opts = {}) {
  const field = (s) => (/[\t\n\r"]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s);
  const escapeHTML = (s) => s.replace(/[&<>"]/g, (c) => (c === "&" ? "&amp;" : c === "<" ? "&lt;" : c === ">" ? "&gt;" : "&quot;"));
  const lowerBound = (list, v) => {
    const i = list.findIndex((x) => x >= v);
    return i < 0 ? list.length : i;
  };
  const { row0, col0 } = range;
  let { row1, col1 } = range;
  const inside = [];
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
  const rows = [], cols = [];
  for (let r = row0; r <= row1; r++) if (!opts.skipRow?.(r)) rows.push(r);
  for (let c = col0; c <= col1; c++) if (!opts.skipCol?.(c)) cols.push(c);
  const at = new Map(), covered = new Set();
  for (const { cell, r0, c0, r1, c1 } of inside) {
    const ri = lowerBound(rows, r0), ci = lowerBound(cols, c0);
    const rn = lowerBound(rows, Math.min(r1, row1) + 1) - ri, cn = lowerBound(cols, Math.min(c1, col1) + 1) - ci;
    if (rn <= 0 || cn <= 0) continue;
    at.set(`${ri},${ci}`, { cell, rows: rn, cols: cn });
    for (let i = 0; i < rn; i++) for (let j = 0; j < cn; j++) if (i || j) covered.add(`${ri + i},${ci + j}`);
  }
  const lines = [];
  let html = '<meta charset="utf-8"><table>';
  for (let i = 0; i < rows.length; i++) {
    const fields = [];
    html += "<tr>";
    for (let j = 0; j < cols.length; j++) {
      const got = at.get(`${i},${j}`);
      fields.push(got ? field(got.cell.text) : "");
      if (covered.has(`${i},${j}`)) continue;
      if (!got) {
        html += "<td></td>";
        continue;
      }
      const tag = got.cell.scope ? "th" : "td";
      let attrs = got.cell.scope ? ` scope="${got.cell.scope}"` : "";
      if (got.rows > 1) attrs += ` rowspan="${got.rows}"`;
      if (got.cols > 1) attrs += ` colspan="${got.cols}"`;
      html += `<${tag}${attrs}>${escapeHTML(got.cell.text).replace(/\n/g, '<br style="mso-data-placement:same-cell">')}</${tag}>`;
    }
    html += "</tr>";
    lines.push(fields.join("\t"));
  }
  html += "</table>";
  return { text: lines.join("\n"), html };
}

test("cellClipboard gives what it gave, for any cells and rectangle", () => {
  // the same numbers every time
  let seed = 12345;
  const random = (n) => {
    seed = (seed * 1103515245 + 12345) % 2147483648;
    return seed % n;
  };
  const texts = ["", "a", "two\nlines", 'q"', "t\tab", "<b>&</b>", "0"];
  for (let round = 0; round < 400; round++) {
    const size = 1 + random(12);
    const cells = [];
    for (let i = random(30); i > 0; i--) {
      // merged cells and cells that lie on others, as no sheet has them
      const cell = { row: random(size), col: random(size), rows: 1 + (random(4) ? 0 : random(4)), cols: 1 + (random(4) ? 0 : random(4)), text: texts[random(texts.length)] };
      if (!random(5)) cell.scope = random(2) ? "col" : "row";
      cells.push(cell);
    }
    const r = [random(size), random(size)].sort((a, b) => a - b), c = [random(size), random(size)].sort((a, b) => a - b);
    const range = { row0: r[0], col0: c[0], row1: r[1] + random(3), col1: c[1] + random(3) };
    const hiddenRow = random(3) ? -1 : random(size), hiddenCol = random(3) ? -1 : random(size);
    const opts = { trim: !random(3), skipRow: (i) => i === hiddenRow, skipCol: (i) => i === hiddenCol };
    assert.deepEqual(cellClipboard(cells, range, opts), reference(cells, range, opts), JSON.stringify({ cells, range, trim: opts.trim, hiddenRow, hiddenCol }));
  }
  // a rectangle upside down, or with no columns left
  for (const range of [{ row0: 3, col0: 0, row1: 1, col1: 2 }, { row0: 0, col0: 3, row1: 2, col1: 1 }]) {
    assert.deepEqual(cellClipboard([], range), reference([], range));
  }
  const none = { skipCol: () => true };
  assert.deepEqual(cellClipboard([], { row0: 0, col0: 0, row1: 2, col1: 2 }, none), reference([], { row0: 0, col0: 0, row1: 2, col1: 2 }, none));
});

test("cellClipboard does not walk a rectangle of too many cells", () => {
  const cell = (row, col, rows = 1, cols = 1) => ({ row, col, rows, cols, text: "x" });
  const tooMany = /too many cells to copy/;
  // two cells whose references are far apart
  const far = [cell(0, 0), cell(3_999_999_999, 0)];
  assert.throws(() => cellClipboard(far, { row0: 0, col0: 0, row1: 3_999_999_999, col1: 0 }), tooMany);
  assert.throws(() => cellClipboard([cell(0, 0), cell(0, 1e12)], { row0: 0, col0: 0, row1: 0, col1: 1e12 }), tooMany);
  assert.throws(() => cellClipboard([], { row0: 0, col0: 0, row1: Infinity, col1: 0 }), tooMany);
  // as many as are copied: MAX_CLIPBOARD_CELLS, or 64 for each cell with something in it
  const column = cellClipboard(far, { row0: 0, col0: 0, row1: MAX_CLIPBOARD_CELLS - 1, col1: 0 });
  assert.equal(column.text.length, 1 + MAX_CLIPBOARD_CELLS - 1);
  assert.throws(() => cellClipboard(far, { row0: 0, col0: 0, row1: MAX_CLIPBOARD_CELLS, col1: 0 }), tooMany);
  // whole columns of a sheet end with the last cell with text
  assert.equal(cellClipboard([cell(0, 0), cell(2, 1)], { row0: 0, col0: 0, row1: 1e9, col1: 1e9 }, { trim: true }).text, "x\t\n\t\n\tx");
  // cells that lie on each other are counted as often
  const merged = Array.from({ length: 100 }, (_, i) => cell(0, i % 10, 2000, 2000));
  assert.throws(() => cellClipboard(merged, { row0: 0, col0: 0, row1: 1999, col1: 1999 }), tooMany);
});
