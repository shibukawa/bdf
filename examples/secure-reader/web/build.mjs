// Builds secure-reader's front end into web/ (next to this script): main.js
// and lib/worker.js (@bdf/render's rendering worker, which fetches and
// opens the segments; no converter modules: the server converts, see
// main.go). index.html and reader.css are checked in as they are.
//
//   node web/build.mjs
import { join } from "node:path";
import { buildWorkers, bundle, root } from "../../common/build.mjs";

const out = join(root, "examples/secure-reader/web");
await bundle(out, [["examples/secure-reader/web/main.ts", "main"]]);
await buildWorkers(join(out, "lib"), { convert: false });
console.log(`secure-reader: web front end built into ${out}`);
