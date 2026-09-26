// Builds the demo viewer with esbuild (and the documentation it links to)
// and serves the repository root.
import { join } from "node:path";
import { serve } from "../../test/serve.mjs";
import { buildViewer, root } from "./build.mjs";
import { buildDocs } from "./docs.mjs";

const out = join(root, "examples/viewer/.out");
await buildViewer(out, { defaultSrc: "/testdata/demo.bdf" });
await buildDocs(out);
const { port } = await serve(root, Number(process.env.PORT ?? 8765));
console.log(`viewer: http://127.0.0.1:${port}/examples/viewer/.out/`);
