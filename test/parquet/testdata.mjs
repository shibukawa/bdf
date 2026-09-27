// Converts converter/parquet/testdata/*.parquet into testdata/parquet/*.bdf
// with the Go converter, laying text out with the test fonts only (those of
// the PowerPoint tests, as for CSV).
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, basename, extname } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const dir = join(root, "converter/parquet/testdata");
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/parquet"), { recursive: true });
for (const f of readdirSync(dir)) {
  if (extname(f) !== ".parquet") continue;
  const out = join(root, "testdata/parquet", basename(f, ".parquet") + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts", join(dir, f), out], { cwd: root, stdio: "inherit" });
}
