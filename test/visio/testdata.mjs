// Converts the Visio test drawings in converter/visio/testdata into
// testdata/visio/*.bdf with the Go converter, laying text out with the test
// fonts only. flow.vsdx is not converted: it is flow.vdx as a package, which
// the Go tests check converts to the same objects.
import { execFileSync } from "node:child_process";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/pptx/testdata/fonts");
for (const [src, out] of [["shapes.vsdx", "shapes.bdf"], ["flow.vdx", "flow.bdf"]]) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts",
    join(root, "converter/visio/testdata", src), join(root, "testdata/visio", out)], { cwd: root, stdio: "inherit" });
}
