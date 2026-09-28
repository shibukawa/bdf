import { test } from "node:test";
import assert from "node:assert/strict";
import { joinRuns } from "../dist/textlayer.js";

// Sep values (docs/spec.md §7.9): NONE 0, SPACE 1, BREAK 2.
const run = (text, sep, layer = "p1") => ({ ordinal: 0, sep, text, layer });

test("joinRuns puts spaces at line breaks and newlines at paragraph breaks", () => {
  assert.equal(joinRuns([run("Hello", 2), run("world", 1), run("!", 0), run("Next", 2)]), "Hello world!\nNext");
});

test("joinRuns drops the first run's own separator", () => {
  assert.equal(joinRuns([run("a", 2), run("b", 2)]), "a\nb");
  assert.equal(joinRuns([run("a", 1)]), "a");
  assert.equal(joinRuns([]), "");
});

test("joinRuns separates pages with a blank line", () => {
  assert.equal(joinRuns([run("end of page", 0, "p1"), run("top of next", 2, "p2")]), "end of page\n\ntop of next");
});

test("linkHref allows web and mail URLs and page and view links only", async () => {
  const { linkHref } = await import("../dist/textlayer.js");
  assert.equal(linkHref("https://example.com/a?b"), "https://example.com/a?b");
  assert.equal(linkHref("mailto:a@example.com"), "mailto:a@example.com");
  assert.equal(linkHref("#page=12"), "#page=12");
  assert.equal(linkHref("#view=page2"), "#view=page2");
  assert.equal(linkHref("#view=a%20b&page=3"), "#view=a%20b&page=3");
  for (const bad of ["javascript:alert(1)", " JavaScript:alert(1)", "data:text/html,x", "file:///etc/passwd", "/relative", "#page=0", "#page=x", "vbscript:x",
    "#view=", "#view=a&page=0", "#view=%E0%A4%A", "#other"]) {
    assert.equal(linkHref(bad), undefined, bad);
  }
});

test("internalLink reads page and view links", async () => {
  const { internalLink } = await import("../dist/textlayer.js");
  assert.deepEqual(internalLink("#page=4"), { page: 4 });
  assert.deepEqual(internalLink("#view=Page-2"), { view: "Page-2" });
  assert.deepEqual(internalLink("#view=%E5%9B%B3+1&page=2"), { view: "図 1", page: 2 });
  assert.equal(internalLink("https://example.com/#view=x"), undefined);
});

/** Elements just enough for buildTextLayer, which Node has none of. */
function fakeDocument() {
  class Element {
    children = [];
    attributes = {};
    style = {};
    constructor(tag) { this.tag = tag; }
    setAttribute(k, v) { this.attributes[k] = String(v); }
    getAttribute(k) { return this.attributes[k] ?? null; }
    hasAttribute(k) { return k in this.attributes; }
    appendChild(c) { this.children.push(c); c.parent = this; return c; }
    get lastElementChild() { return this.children.at(-1) ?? null; }
    closest() { return null; }
  }
  return { createElement: (tag) => new Element(tag), addEventListener() {} };
}

test("runs lie in the links that cover them, and links are read once each", async () => {
  const { buildTextLayer } = await import("../dist/textlayer.js");
  const saved = { document: globalThis.document, window: globalThis.window, URL: globalThis.URL };
  let parsed = 0;
  globalThis.document = fakeDocument();
  globalThis.window = { addEventListener() {} };
  globalThis.URL = class extends saved.URL { constructor(...a) { super(...a); parsed++; } };
  try {
    const font = { kind: 1, family: "sans-serif", weight: 400, style: 0 };
    const run = (text, x, y) => ({ text, x, y, advance: 40, size: 10, font, align: 0, matrix: [1, 0, 0, 1, x, y], sep: 2, ordinal: 0, altText: false });
    const link = (url, x, y, w, h, after) => ({ url, x, y, w, h, after, node: -1 });
    const content = {
      runs: [run("plain", 0, 10), run("linked", 0, 30), run("too", 50, 30), run("unsafe", 0, 50), run("inner", 0, 70), run("page", 0, 90)],
      nodes: [],
      links: [
        link("https://example.com/a?b", 0, 20, 100, 12, 0),
        link("javascript:alert(1)", 0, 40, 100, 12, 2),
        // one within another: the smaller one has the run
        link("https://example.com/outer", 0, 55, 200, 40, 3),
        link("mailto:a@example.com", 0, 60, 60, 12, 3),
        link("#page=3", 0, 80, 100, 12, 4),
        // no run in them: a link of their own, but for the one that may not be one
        link("https://example.com/picture", 300, 300, 50, 50, 5),
        link("data:text/html,x", 400, 300, 50, 50, 5),
      ],
    };
    const layer = buildTextLayer(content, 2, { measure: () => 40 });
    const shape = (el) => (el.tag === "span" ? el.textContent : { [el.tag === "a" ? `a ${el.href}` : el.tag]: el.children.map(shape) });
    assert.deepEqual(shape(layer), { div: [
      "plain",
      { "a https://example.com/a?b": ["linked", "too"] },
      "unsafe",
      { "a https://example.com/outer": [] },
      { "a mailto:a@example.com": ["inner"] },
      { "a #page=3": ["page"] },
      { "a https://example.com/picture": [] },
      { div: [] }, // the end of the layer
    ] });
    const [, web, , , , page, picture] = layer.children;
    assert.deepEqual([web.target, web.rel, web.attributes["data-bdf-page"]], ["_blank", "noopener noreferrer", undefined]);
    assert.deepEqual([page.target, page.attributes["data-bdf-page"]], [undefined, "3"]);
    assert.equal(picture.attributes["aria-label"], "https://example.com/picture");
    assert.deepEqual([web.style.left, web.style.top, web.style.width], ["0px", "40px", "200px"]);
    // not once for every run and link
    assert.ok(parsed <= 2 * content.links.length, `${parsed} URLs read`);
  } finally {
    Object.assign(globalThis, saved);
    if (saved.document === undefined) delete globalThis.document;
    if (saved.window === undefined) delete globalThis.window;
  }
});
