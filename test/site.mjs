// Smoke test of the demo site (examples/viewer/site.mjs): serves it, runs
// its converter modules in Node with Go's wasm_exec.js, and converts each
// sample with the site's fonts, fetched from the server as in a browser.
//
//   node test/site.mjs [site dir]
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { join, resolve } from "node:path";
import { pathToFileURL } from "node:url";
import { serve } from "./serve.mjs";

const site = resolve(process.argv[2] ?? "examples/viewer/.site");
await import(pathToFileURL(join(site, "wasm_exec.js")).href);
const { server, port } = await serve(site);
const base = `http://127.0.0.1:${port}/`;

/** Load a module; the Go program puts its API on globalThis and waits. */
async function load(file) {
  const go = new globalThis.Go();
  const { instance } = await WebAssembly.instantiate(await readFile(join(site, file)), go.importObject);
  go.run(instance);
  const conv = globalThis.bdfConverter;
  delete globalThis.bdfConverter;
  return conv;
}
const modules = { pdf: await load("bdf-pdf.wasm"), office: await load("bdf-office.wasm") };
assert.deepEqual(modules.pdf.formats.map((f) => f.name), ["pdf"]);
assert.deepEqual(modules.office.formats.map((f) => f.name), ["csv", "docx", "drawio", "dxf", "emf", "pptx", "visio", "xlsx"]);

let failed = 0;
const samples = JSON.parse(await readFile(join(site, "samples/index.json"), "utf8"));
for (const { name } of samples) {
  if (name.endsWith(".bdf")) continue;
  const data = new Uint8Array(await readFile(join(site, "samples", name)));
  const conv = name.endsWith(".pdf") ? modules.pdf : modules.office;
  const t0 = performance.now();
  try {
    const res = await conv.convert(data, { fonts: `${base}fonts/` });
    assert.deepEqual([...res.bdf.subarray(0, 4)], [0x62, 0x64, 0x66, 0], "bdf magic");
    // the Office converters found the site's fonts and embedded them
    if (res.format !== "pdf") assert.match(res.summary, /[1-9]\d* embedded font/);
    console.log(`ok   ${name}: ${res.format}, ${res.summary}, ${res.bdf.length} bytes, ${(performance.now() - t0).toFixed(0)} ms`);
    for (const w of res.warnings) console.log(`     warning: ${w}`);
  } catch (e) {
    failed++;
    console.log(`FAIL ${name}: ${e.message}`);
  }
}

// errors carry the codes the page acts on
await assert.rejects(modules.office.convert(new Uint8Array([1, 2, 3]), {}), { code: "unknown-format" });

server.close();
if (failed) {
  console.log(`${failed} sample(s) failed`);
  process.exit(1);
}
process.exit(0); // the Go programs wait for calls forever
