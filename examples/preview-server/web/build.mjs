// Builds preview-server's front end into web/ (next to this script): main.js
// and lib/worker.js (@bdfkit/render's rendering worker — no converter modules:
// this server's browser never converts anything, see main.go). index.html
// is checked in as it is.
//
//   node web/build.mjs
import { join } from "node:path";
import { buildWorkers, bundle, root } from "../../common/build.mjs";

const out = join(root, "examples/preview-server/web");
await bundle(out, [["examples/preview-server/web/main.ts", "main"]]);
await buildWorkers(join(out, "lib"), { convert: false });
console.log(`preview-server: web front end built into ${out}`);
