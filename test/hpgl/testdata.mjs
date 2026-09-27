// Converts the HP-GL/2 test plots in converter/hpgl/testdata into
// testdata/hpgl/*.bdf with the Go converter, laying text out with the test
// fonts only.
import { execFileSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/hpgl"), { recursive: true });
for (const n of ["shapes", "job", "hpgl1"]) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts",
    join(root, "converter/hpgl/testdata", n + ".plt"), join(root, "testdata/hpgl", n + ".bdf")], { cwd: root, stdio: "inherit" });
}
