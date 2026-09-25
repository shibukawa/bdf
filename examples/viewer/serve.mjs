// Builds the demo viewer with esbuild and serves the repository root.
import { build } from "esbuild";
import { copyFile, mkdir } from "node:fs/promises";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { serve } from "../../test/serve.mjs";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const out = join(root, "examples/viewer/.out");
await mkdir(out, { recursive: true });
await build({
  entryPoints: [
    { in: join(root, "examples/viewer/main.ts"), out: "main" },
    { in: join(root, "packages/render/src/worker.ts"), out: "worker" },
  ],
  bundle: true, format: "esm", outdir: out, sourcemap: true, target: "es2022", logLevel: "info",
});
await copyFile(join(root, "examples/viewer/index.html"), join(out, "index.html"));
const { port } = await serve(root, Number(process.env.PORT ?? 8765));
console.log(`viewer: http://127.0.0.1:${port}/examples/viewer/.out/`);
