// Text content of several objects put together, as the worker answers
// requests for the text of a page, a band or a region of a sheet.
import type { Rect, TextContent, TextRun } from "@bdfkit/core";

/** Join the contents of several objects (layers, pages, tiles), renumbering nodes and runs. */
export function concat(parts: TextContent[]): TextContent {
  const out: TextContent = { runs: [], nodes: [], links: [] };
  for (const c of parts) {
    const nodeBase = out.nodes.length, runBase = out.runs.length;
    const shift = (i: number | undefined) => (i === undefined || i < 0 ? i : i + nodeBase);
    // one by one: a call takes so many arguments, and a layer may have more runs
    for (const n of c.nodes) out.nodes.push({ ...n, parent: shift(n.parent)! });
    for (const r of c.runs) out.runs.push({ ...r, node: shift(r.node) });
    for (const l of c.links) out.links.push({ ...l, after: l.after + runBase, node: shift(l.node)! });
  }
  return out;
}

/**
 * Keep the runs and links inside a rectangle (by run anchor, by link
 * center), for a band or a tile; the text layer leaves out nodes left empty.
 */
export function within(c: TextContent, r: Rect, dx = 0, dy = 0): TextContent {
  const inside = (x: number, y: number) => x - dx >= r.x && x - dx < r.x + r.w && y - dy >= r.y && y - dy <= r.y + r.h;
  const runs: TextRun[] = [];
  const kept: number[] = []; // new index of the last kept run up to each old index
  for (const run of c.runs) {
    if (inside(run.x, run.y)) runs.push(run);
    kept.push(runs.length - 1);
  }
  const links = c.links.filter((l) => inside(l.x + l.w / 2, l.y + l.h / 2)).map((l) => ({ ...l, after: l.after >= 0 ? kept[l.after] : -1 }));
  return { runs, nodes: c.nodes, links };
}
