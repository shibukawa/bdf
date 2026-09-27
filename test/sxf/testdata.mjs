// Converts the SXF test drawing in converter/sxf/testdata (its P21 file)
// into testdata/sxf/shapes.bdf with the Go converter, laying text out with
// the test fonts only.
import { execFileSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/sxf"), { recursive: true });
execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts",
  join(root, "converter/sxf/testdata/shapes.p21"), join(root, "testdata/sxf/shapes.bdf")], { cwd: root, stdio: "inherit" });
