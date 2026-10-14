import { test } from "node:test";
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import ts from "typescript";

const source = await readFile(new URL("../site/form/data.ts", import.meta.url), "utf8");
const js = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.ESNext, target: ts.ScriptTarget.ES2022 } }).outputText;
const { valueOf, setValue, pictureFiles } = await import(`data:text/javascript;base64,${Buffer.from(js).toString("base64")}`);

test("missing special field names do not expose or mutate inherited objects", () => {
  const data = {};
  for (const name of ["__proto__", "constructor", "toString"]) {
    assert.equal(valueOf(data, name), undefined);
    const nested = {};
    setValue(nested, "formValue", "typed");
    setValue(data, name, nested);
    assert.equal(valueOf(data, name), nested);
  }
  assert.equal(Object.getPrototypeOf(data), Object.prototype);
  assert.equal({}.formValue, undefined);
  assert.deepEqual(JSON.parse(JSON.stringify(data)), JSON.parse('{"__proto__":{"formValue":"typed"},"constructor":{"formValue":"typed"},"toString":{"formValue":"typed"}}'));
  const supplied = JSON.parse('{"__proto__":{"name":"existing"}}');
  assert.deepEqual(valueOf(supplied, "__proto__"), { name: "existing" });
});

test("picture attachments match string values in nested data, including special names", () => {
  const bytes = new Uint8Array([1, 2, 3]);
  const pictures = new Map(["logo.png", "__proto__", "only-a-key.png", "unused.png"].map((name) => [name, bytes]));
  const data = { items: [{ image: "logo.png" }, { image: "__proto__" }], "only-a-key.png": false, text: "prefix logo.png suffix" };
  const files = pictureFiles(data, pictures);
  assert.deepEqual(Object.keys(files), ["__proto__", "logo.png"]);
  assert.equal(Object.getPrototypeOf(files), Object.prototype);
  const copied = structuredClone(files);
  assert.deepEqual(copied.__proto__, bytes);
  assert.deepEqual(copied["logo.png"], bytes);
});

test("picture selection handles deeply nested JSON without recursive calls", () => {
  const bytes = new Uint8Array([1]);
  let data = "logo.png";
  for (let i = 0; i < 20_000; i++) data = { child: data };
  assert.deepEqual(pictureFiles(data, new Map([["logo.png", bytes]])), { "logo.png": bytes });
});
