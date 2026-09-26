// Converts converter/csv/testdata/*.csv and *.tsv into testdata/csv/*.bdf
// with the Go converter, laying text out with the test fonts only (those of
// the PowerPoint tests).
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, extname, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const dir = join(root, "converter/csv/testdata");
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/csv"), { recursive: true });
for (const f of readdirSync(dir)) {
  if (![".csv", ".tsv"].includes(extname(f))) continue;
  const out = join(root, "testdata/csv", basename(f, extname(f)) + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts", join(dir, f), out], { cwd: root, stdio: "inherit" });
}
