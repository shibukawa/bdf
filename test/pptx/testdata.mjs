// Converts converter/pptx/testdata/*.pptx into testdata/pptx/*.bdf with the
// Go converter, laying text out with the test fonts only. The templates
// form-*.pptx are also filled with the data of the Excel templates
// (test/xlsx/data; testdata/form/pptx-*.bdf), and the form of one slide is
// written as a PDF file (testdata/form/pptx-fixed.pdf).
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const dir = join(root, "converter/pptx/testdata");
for (const f of readdirSync(dir)) {
  if (!f.endsWith(".pptx")) continue;
  const out = join(root, "testdata/pptx", basename(f, ".pptx") + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", join(dir, "fonts"), "-no-system-fonts", join(dir, f), out], { cwd: root, stdio: "inherit" });
}
const data = join(root, "test/xlsx/data");
mkdirSync(join(root, "testdata/form"), { recursive: true });
for (const n of ["fixed", "rows", "pages"]) {
  execFileSync("go", ["run", "./cmd/bdf", "fill", "-q", "-font-dir", join(dir, "fonts"), "-no-system-fonts", "-files", data, "-data", join(data, `form-${n}.json`),
    join(dir, `form-${n}.pptx`), join(root, `testdata/form/pptx-${n}.bdf`)], { cwd: root, stdio: "inherit" });
}
execFileSync("go", ["run", "./cmd/bdf", "pdf", "-q", "-font-dir", join(dir, "fonts"), "-no-system-fonts", join(root, "testdata/form/pptx-fixed.bdf"),
  join(root, "testdata/form/pptx-fixed.pdf")], { cwd: root, stdio: "inherit" });
