// Builds the demo viewer with esbuild and serves the repository root.
import { join } from "node:path";
import { serve } from "../../test/serve.mjs";
import { buildViewer, root } from "./build.mjs";

await buildViewer(join(root, "examples/viewer/.out"), { defaultSrc: "/testdata/demo.bdf" });
const { port } = await serve(root, Number(process.env.PORT ?? 8765));
console.log(`viewer: http://127.0.0.1:${port}/examples/viewer/.out/`);
