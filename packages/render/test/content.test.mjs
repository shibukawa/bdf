import { test } from "node:test";
import assert from "node:assert/strict";
import { concat, within } from "../dist/content.js";

const run = (text, node, x = 0, y = 0) => ({ text, x, y, advance: 10, size: 10, font: undefined, align: 0, matrix: [1, 0, 0, 1, x, y], sep: 0, ordinal: 0, altText: false, node });

test("concat renumbers the nodes, runs and links of the contents it joins", () => {
  const a = {
    runs: [run("a1", 1), run("a2", undefined)],
    nodes: [{ kind: "list", parent: -1 }, { kind: "item", parent: 0 }],
    links: [{ x: 0, y: 0, w: 1, h: 1, url: "#page=1", after: -1, node: -1 }, { x: 0, y: 0, w: 1, h: 1, url: "#page=2", after: 1, node: 1 }],
  };
  const b = {
    runs: [run("b1", 0)],
    nodes: [{ kind: "paragraph", parent: -1 }],
    links: [{ x: 0, y: 0, w: 1, h: 1, url: "#page=3", after: 0, node: 0 }],
  };
  const c = concat([a, b, a]);
  assert.deepEqual(c.runs.map((r) => [r.text, r.node]), [["a1", 1], ["a2", undefined], ["b1", 2], ["a1", 4], ["a2", undefined]]);
  assert.deepEqual(c.nodes.map((n) => [n.kind, n.parent]), [["list", -1], ["item", 0], ["paragraph", -1], ["list", -1], ["item", 3]]);
  assert.deepEqual(c.links.map((l) => [l.url, l.after, l.node]), [["#page=1", -1, -1], ["#page=2", 1, 1], ["#page=3", 2, 2], ["#page=1", 2, -1], ["#page=2", 4, 4]]);
  // the contents are left as they were
  assert.deepEqual(a.runs.map((r) => r.node), [1, undefined]);
  assert.deepEqual(concat([]), { runs: [], nodes: [], links: [] });
});

test("concat joins layers of more runs than a call takes arguments", () => {
  const n = 300_000;
  const layer = {
    runs: Array.from({ length: n }, (_, i) => run("x", i)),
    nodes: Array.from({ length: n }, () => ({ kind: "paragraph", parent: -1 })),
    links: Array.from({ length: n }, (_, i) => ({ x: 0, y: 0, w: 1, h: 1, url: "#page=1", after: i, node: i })),
  };
  const c = concat([layer, layer]);
  assert.equal(c.runs.length, 2 * n);
  assert.equal(c.runs.at(-1).node, 2 * n - 1);
  assert.equal(c.nodes.length, 2 * n);
  assert.deepEqual([c.links.at(-1).after, c.links.at(-1).node], [2 * n - 1, 2 * n - 1]);
});

test("within keeps the runs and links of a rectangle", () => {
  const c = {
    runs: [run("in", 0, 10, 10), run("out", 0, 10, 200), run("moved in", 0, 110, 50)],
    nodes: [{ kind: "paragraph", parent: -1 }],
    links: [{ x: 5, y: 5, w: 10, h: 10, url: "#page=1", after: 1, node: 0 }, { x: 5, y: 195, w: 10, h: 10, url: "#page=2", after: 1, node: 0 }],
  };
  const w = within(c, { x: 0, y: 0, w: 100, h: 100 });
  assert.deepEqual(w.runs.map((r) => r.text), ["in"]);
  assert.deepEqual(w.links.map((l) => [l.url, l.after]), [["#page=1", 0]]);
  assert.deepEqual(within(c, { x: 0, y: 0, w: 100, h: 100 }, 100, 0).runs.map((r) => r.text), ["moved in"]);
});
