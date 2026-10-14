// Converts converter/docx/testdata/*.docx into testdata/docx/*.bdf with the
// Go converter, laying text out with the test fonts only. The templates
// form-*.docx are also filled with the data of the Excel templates
// (test/xlsx/data; testdata/form/docx-*.bdf).
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const dir = join(root, "converter/docx/testdata");
mkdirSync(join(root, "testdata/docx"), { recursive: true });
for (const f of readdirSync(dir)) {
  if (!f.endsWith(".docx")) continue;
  const out = join(root, "testdata/docx", basename(f, ".docx") + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", join(dir, "fonts"), "-no-system-fonts", join(dir, f), out], { cwd: root, stdio: "inherit" });
}
const data = join(root, "test/xlsx/data");
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/form"), { recursive: true });
for (const n of ["fixed", "rows", "pages"]) {
  execFileSync("go", ["run", "./cmd/bdf", "fill", "-q", "-font-dir", fonts, "-no-system-fonts", "-files", data, "-data", join(data, `form-${n}.json`),
    join(dir, `form-${n}.docx`), join(root, `testdata/form/docx-${n}.bdf`)], { cwd: root, stdio: "inherit" });
}
