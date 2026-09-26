// Builds the demo viewer with esbuild: the page, the rendering worker and the
// converter worker (a classic worker, which loads Go's wasm_exec.js).
import { build } from "esbuild";
import { copyFile, mkdir } from "node:fs/promises";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

export const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));

/**
 * Build the viewer into out. defaultSrc is the document shown when the URL
 * has no ?src= ("" shows the page for opening a file).
 */
export async function buildViewer(out, { defaultSrc }) {
  await mkdir(out, { recursive: true });
  const common = { bundle: true, outdir: out, sourcemap: true, target: "es2022", logLevel: "info" };
  await build({
    ...common,
    entryPoints: [
      { in: join(root, "examples/viewer/main.ts"), out: "main" },
      { in: join(root, "packages/render/src/worker.ts"), out: "worker" },
    ],
    format: "esm",
    define: { DEFAULT_SRC: JSON.stringify(defaultSrc) },
  });
  await build({ ...common, entryPoints: [{ in: join(root, "examples/viewer/convert-worker.ts"), out: "convert-worker" }], format: "iife" });
  await copyFile(join(root, "examples/viewer/index.html"), join(out, "index.html"));
}
