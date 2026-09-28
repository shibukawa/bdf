// Builds the demo viewer with esbuild and serves the repository root: the
// viewer by itself, which shows the bdf documents of the repository
// (?src=/testdata/docx/basic.bdf). It converts no files: the demo site
// (npm run site:serve) has the converters.
import { join } from "node:path";
import { serve } from "../../test/serve.mjs";
import { buildViewer, root } from "./build.mjs";

const out = join(root, "examples/viewer/.out");
await buildViewer(out, { defaultSrc: "/testdata/demo.bdf" });
const { port } = await serve(root, Number(process.env.PORT ?? 8765));
console.log(`viewer: http://127.0.0.1:${port}/examples/viewer/.out/`);
