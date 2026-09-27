// Converts the KiCad test projects in converter/kicad/testdata (the demo
// project's archive, and the project with its own drawing sheet) into
// testdata/kicad/*.bdf with the Go converter, laying out the Japanese text
// with the test fonts only.
import { execFileSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/kicad"), { recursive: true });
for (const [src, out] of [["demo.zip", "demo"], ["frame/frame.kicad_pro", "frame"]]) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts",
    join(root, "converter/kicad/testdata", src), join(root, "testdata/kicad", out + ".bdf")], { cwd: root, stdio: "inherit" });
}
