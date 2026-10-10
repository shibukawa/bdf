import { test } from "node:test";
import assert from "node:assert/strict";
import { findKey, FindHits } from "../dist/find.js";

const key = (key, mods = {}) => ({ key, code: mods.code ?? "", metaKey: false, ctrlKey: false, shiftKey: false, altKey: false, isComposing: false, ...mods });

test("the find shortcut is Command's on macOS and Control's elsewhere", () => {
  assert.equal(findKey(key("f", { metaKey: true }), true), "open");
  assert.equal(findKey(key("f", { ctrlKey: true }), false), "open");
  // Control+F moves the caret in a text field on macOS
  assert.equal(findKey(key("f", { ctrlKey: true }), true), undefined);
  assert.equal(findKey(key("f", { metaKey: true }), false), undefined);
  assert.equal(findKey(key("f"), true), undefined);
  assert.equal(findKey(key("f", { metaKey: true, altKey: true }), true), undefined);
  assert.equal(findKey(key("F", { metaKey: true, shiftKey: true }), true), undefined);
  // not while an input method composes text
  assert.equal(findKey(key("f", { metaKey: true, isComposing: true }), true), undefined);
});

test("the next and previous match keys", () => {
  assert.equal(findKey(key("g", { metaKey: true }), true), "next");
  assert.equal(findKey(key("G", { metaKey: true, shiftKey: true }), true), "previous");
  assert.equal(findKey(key("g", { ctrlKey: true }), false), "next");
  assert.equal(findKey(key("F3"), false), "next");
  assert.equal(findKey(key("F3", { shiftKey: true }), true), "previous");
  assert.equal(findKey(key("F3", { ctrlKey: true }), false), undefined);
});

test("a keyboard of another script is read by the place of the key, one of Latin letters by its letter", () => {
  // the F key of a Russian keyboard types а
  assert.equal(findKey(key("а", { ctrlKey: true, code: "KeyF" }), false), "open");
  // the key at the place of F types u on a Dvorak keyboard: not the shortcut
  assert.equal(findKey(key("u", { ctrlKey: true, code: "KeyF" }), false), undefined);
  assert.equal(findKey(key("f", { ctrlKey: true, code: "KeyY" }), false), "open");
});

const hit = (a, b, ordinal, start = 0) => ({ segments: [{ run: 0, a, b, ordinal, start, end: start + 1 }], text: "x", context: "x" });

/** A worker that finds the hits given and notes what it is asked to locate. */
function fakeClient(hits) {
  const located = [];
  return {
    located,
    searches: [],
    async search(view, query, options) {
      this.searches.push({ view, query, options });
      return hits.slice(0, options.limit);
    },
    async locate(view, asked) {
      located.push(asked.map((h) => `${h.segments[0].a}:${h.segments[0].ordinal}`));
      return asked.map((h) => [{ a: h.segments[0].a, b: h.segments[0].b, x: h.segments[0].ordinal, y: 0, w: 1, h: 1 }]);
    },
  };
}

test("hits are located a page at a time, each once", async () => {
  const client = fakeClient([hit(0, 0, 1), hit(0, 0, 2), hit(3, 0, 1), hit(7, 1, 4)]);
  const found = await FindHits.search(client, { id: "pages", kind: "fixed" }, "x");
  assert.equal(found.length, 4);
  assert.equal(found.more, false);
  assert.equal(found.rects(0), undefined);
  assert.equal(await found.locate([0]), true);
  assert.deepEqual(client.located, [["0:1", "0:2"]]);
  assert.equal(found.rects(1)[0].x, 2);
  assert.equal(found.rects(2), undefined);
  // located already: nothing is asked for, and pages without hits ask nothing
  assert.equal(await found.locate([0, 1, 2]), false);
  assert.equal(client.located.length, 1);
  // a hit is located with the others of its page
  assert.deepEqual(await found.locateHit(3), [{ a: 7, b: 1, x: 4, y: 0, w: 1, h: 1 }]);
  assert.deepEqual(client.located[1], ["7:4"]);
  const seen = [];
  found.each((r, i) => seen.push(`${i}@${r.a}`));
  assert.deepEqual(seen, ["0@0", "1@0", "3@7"]);
});

test("pages asked for while they are being located are not asked for twice", async () => {
  const client = fakeClient([hit(2, 0, 1), hit(2, 0, 2)]);
  const found = await FindHits.search(client, { id: "pages", kind: "fixed" }, "x");
  const [first, second, rects] = await Promise.all([found.locate([2]), found.locate([2]), found.locateHit(1)]);
  assert.equal(first, true);
  assert.equal(second, true);
  assert.equal(rects[0].x, 2);
  assert.equal(client.located.length, 1);
});

test("a failed location is asked for again", async () => {
  const client = fakeClient([hit(1, 0, 1)]);
  const locate = client.locate;
  client.locate = async () => { throw new Error("offline"); };
  const found = await FindHits.search(client, { id: "pages", kind: "fixed" }, "x");
  await assert.rejects(found.locate([1]), /offline/);
  client.locate = locate;
  assert.equal(await found.locate([1]), true);
  assert.equal(found.rects(0).length, 1);
});

test("a search starts at the page being read and wraps", async () => {
  const client = fakeClient([hit(0, 0, 1), hit(3, 0, 1), hit(3, 0, 2), hit(7, 0, 1)]);
  const found = await FindHits.search(client, { id: "pages", kind: "fixed" }, "x");
  assert.equal(found.from(0), 0);
  assert.equal(found.from(1), 1);
  assert.equal(found.from(3), 1);
  assert.equal(found.from(4), 3);
  assert.equal(found.from(8), 0);
});

test("the hits of a sheet are located by tile", async () => {
  const client = fakeClient([hit(0, 0, 1), hit(1, 0, 1), hit(0, 5, 1)]);
  const found = await FindHits.search(client, { id: "sheet", kind: "sheet" }, "x");
  assert.equal(found.from(3), 0);
  assert.equal(await found.locate(["0,5", "4,4"]), true);
  assert.deepEqual(client.located, [["0:1"]]);
  assert.equal(found.rects(2).length, 1);
  assert.equal(found.rects(0), undefined);
});

test("more hits than the limit are told, and spaces alone find nothing", async () => {
  const client = fakeClient(Array.from({ length: 30 }, (_, i) => hit(i, 0, 1)));
  const found = await FindHits.search(client, { id: "pages", kind: "fixed" }, "x", 10);
  assert.equal(client.searches[0].options.limit, 11);
  assert.equal(found.length, 10);
  assert.equal(found.more, true);
  const none = await FindHits.search(client, { id: "pages", kind: "fixed" }, "  ");
  assert.equal(none.length, 0);
  assert.equal(client.searches.length, 1);
});
