// Converts the CGM test metafiles in converter/cgm/testdata (the binary
// ones) into testdata/cgm with the Go converter, laying text out with the
// test fonts only.
import { execFileSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/cgm"), { recursive: true });
for (const name of ["shapes", "illustration", "sjis"]) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts",
    join(root, `converter/cgm/testdata/${name}.cgm`), join(root, `testdata/cgm/${name}.bdf`)], { cwd: root, stdio: "inherit" });
}
