import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import ts from "typescript";

test("search snippets render document HTML as text, with only mark highlights", async () => {
  class Element {
    constructor(tagName) { this.tagName = tagName; this.children = []; this.listeners = {}; }
    append(...nodes) { this.children.push(...nodes); }
    replaceChildren(...nodes) { this.children = nodes; }
    addEventListener(name, listener) { this.listeners[name] = listener; }
    set innerHTML(_) { throw new Error("snippet was inserted as HTML"); }
  }
  const form = new Element("form"), input = new Element("input"), list = new Element("ul");
  const nodes = { searchForm: form, q: input, hits: list };
  const oldDocument = globalThis.document, oldFetch = globalThis.fetch;
  globalThis.document = {
    getElementById: (id) => nodes[id],
    createElement: (name) => new Element(name),
    createTextNode: (text) => ({ textContent: text }),
  };
  globalThis.fetch = async () => ({ json: async () => [{
    bdf: "book.bdf", view: "main", title: "Book", page: 1,
    snippet: 'before <img src=x onerror=alert(1)> <mark>needle</mark> after',
  }] });
  try {
    const source = await readFile(new URL("./search.ts", import.meta.url), "utf8");
    const js = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.ESNext } }).outputText;
    await import(`data:text/javascript;base64,${Buffer.from(js).toString("base64")}`);
    input.value = "needle";
    form.listeners.submit({ preventDefault() {} });
    await new Promise(setImmediate);
    const snippet = list.children[0].children[1];
    assert.deepEqual(snippet.children.map((n) => [n.tagName ?? "text", n.textContent ?? ""]), [
      ["text", "before <img src=x onerror=alert(1)> "], ["mark", ""], ["text", " after"],
    ]);
    assert.equal(snippet.children[1].children[0].textContent, "needle");
  } finally {
    globalThis.document = oldDocument;
    globalThis.fetch = oldFetch;
  }
});
