// Converts the JWW test drawings in converter/jww/testdata into
// testdata/jww/*.bdf with the Go converter, laying text out with the test
// fonts only.
import { execFileSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/jww"), { recursive: true });
for (const n of ["shapes", "old"]) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts",
    join(root, "converter/jww/testdata", n + ".jww"), join(root, "testdata/jww", n + ".bdf")], { cwd: root, stdio: "inherit" });
}
