// Builds light-server's front end into web/ (next to this script): main.js
// (from main.ts), lib/ (the rendering worker, the converter worker and the
// wasm converter modules — reusing examples/common/build.mjs, the same
// helper the demo site is built with) and fonts/ (the free fonts the Office
// converters lay text out with). index.html is checked in as it is.
//
//   node web/build.mjs
import { join } from "node:path";
import { buildModules, buildWorkers, bundle, copyFonts, root } from "../../common/build.mjs";

const out = join(root, "examples/light-server/web");
await bundle(out, [["examples/light-server/web/main.ts", "main"]]);
await buildWorkers(join(out, "lib"));
await buildModules(join(out, "lib"), ["bdf-pdf.wasm", "bdf-office.wasm", "bdf-web.wasm", "bdf-image.wasm"]); // no bdf-preview.wasm: this server never draws a thumbnail in the browser
await copyFonts(join(out, "fonts"));
console.log(`light-server: web front end built into ${out}`);
