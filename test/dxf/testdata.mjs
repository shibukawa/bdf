// Converts the DXF test drawings in converter/dxf/testdata into
// testdata/dxf/*.bdf with the Go converter, laying text out with the test
// fonts only. shapes-bin.dxf is not converted: it is shapes.dxf in binary,
// which the Go tests check converts to the same objects.
import { execFileSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/dxf"), { recursive: true });
for (const n of ["shapes", "layout", "r12-sjis"]) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts",
    join(root, "converter/dxf/testdata", n + ".dxf"), join(root, "testdata/dxf", n + ".bdf")], { cwd: root, stdio: "inherit" });
}
