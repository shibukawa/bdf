// Make @bdfkit/render/worker a self-contained browser asset.
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { build } from "esbuild";

const root = dirname(dirname(fileURLToPath(import.meta.url)));
await build({
  entryPoints: [join(root, "packages/render/src/worker.ts")],
  outfile: join(root, "packages/render/dist/worker.js"),
  bundle: true,
  format: "esm",
  target: "es2022",
});
