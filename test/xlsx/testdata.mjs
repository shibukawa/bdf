// Converts converter/xlsx/testdata/*.xlsx into testdata/xlsx/*.bdf with the
// Go converter, laying text out with the test fonts only (those of the
// PowerPoint tests).
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const dir = join(root, "converter/xlsx/testdata");
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/xlsx"), { recursive: true });
for (const f of readdirSync(dir)) {
  if (!f.endsWith(".xlsx")) continue;
  const out = join(root, "testdata/xlsx", basename(f, ".xlsx") + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts", join(dir, f), out], { cwd: root, stdio: "inherit" });
}
