// Converts converter/xlsx/testdata/*.xlsx into testdata/xlsx/*.bdf with the
// Go converter, laying text out with the test fonts only (those of the
// PowerPoint tests). print.xlsx is also converted as the pages it prints
// (testdata/xlsx/print-pages.bdf), and the templates form-*.xlsx are filled
// with the data of test/xlsx/data (testdata/form/*.bdf; pages.bdft is the
// slip as a compiled template).
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const dir = join(root, "converter/xlsx/testdata");
const fonts = join(root, "converter/pptx/testdata/fonts");
const bdf = (...args) => execFileSync("go", ["run", "./cmd/bdf", ...args], { cwd: root, stdio: "inherit" });
mkdirSync(join(root, "testdata/xlsx"), { recursive: true });
mkdirSync(join(root, "testdata/form"), { recursive: true });
for (const f of readdirSync(dir)) {
  if (!f.endsWith(".xlsx")) continue;
  const name = basename(f, ".xlsx");
  bdf("generate", "-q", "-font-dir", fonts, "-no-system-fonts", join(dir, f), join(root, "testdata/xlsx", name + ".bdf"));
  if (name.startsWith("form-")) {
    const data = join(root, "test/xlsx/data");
    bdf("fill", "-q", "-font-dir", fonts, "-no-system-fonts", "-files", data, "-data", join(data, name + ".json"), join(dir, f), join(root, "testdata/form", name.slice(5) + ".bdf"));
  }
}
{
  // a compiled template: the slip, filled with its data as the sample
  const data = join(root, "test/xlsx/data");
  bdf("template", "compile", "-q", "-font-dir", fonts, "-no-system-fonts", "-files", data, "-sample", join(data, "form-pages.json"), join(dir, "form-pages.xlsx"), join(root, "testdata/form/pages.bdft"));
}
bdf("generate", "-q", "-font-dir", fonts, "-no-system-fonts", "-param", "views=print", join(dir, "print.xlsx"), join(root, "testdata/xlsx/print-pages.bdf"));
