// Converts converter/epub/testdata/*.epub into testdata/epub/*.bdf with the
// Go converter. Text is laid out with the metrics of the Word converter's
// test fonts and refers to the fonts by name (the default for EPUB, as for
// HTML), so viewers draw it with their own fonts.
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, basename, extname } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/docx/testdata/fonts");
const dir = join(root, "converter/epub/testdata");
mkdirSync(join(root, "testdata/epub"), { recursive: true });
for (const f of readdirSync(dir)) {
  if (extname(f) !== ".epub") continue;
  const out = join(root, "testdata/epub", basename(f, ".epub") + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts", join(dir, f), out],
    { cwd: root, stdio: "inherit" });
}
