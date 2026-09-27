// Converts test fonts of converter/font/testdata into testdata/font/*.bdf
// with the Go converter, laying the views' own text out with the test fonts
// only (those of the PowerPoint tests): a CFF OpenType font, embedded, and
// a font whose license forbids embedding it, drawn as outlines.
import { execFileSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/font"), { recursive: true });
for (const [src, name] of [["stix.otf", "stix"], ["restricted.ttf", "restricted"]]) {
  const out = join(root, "testdata/font", name + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts", join(root, "converter/font/testdata", src), out], { cwd: root, stdio: "inherit" });
}
