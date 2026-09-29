// Builds search's front end into web/ (next to this script): main.js (the
// document page), search.js (the search box on the home page, served by
// main.go's own template), and lib/worker.js (@bdfkit/render's rendering
// worker — no converter modules: this server's browser never converts
// anything, see main.go). index.html is checked in as it is.
//
//   node web/build.mjs
import { join } from "node:path";
import { buildWorkers, bundle, root } from "../../common/build.mjs";

const out = join(root, "examples/search/web");
await bundle(out, [
  ["examples/search/web/main.ts", "main"],
  ["examples/search/web/search.ts", "search"],
]);
await buildWorkers(join(out, "lib"), { convert: false });
console.log(`search: web front end built into ${out}`);
