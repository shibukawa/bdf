// Converts converter/docx/testdata/*.docx into testdata/docx/*.bdf with the
// Go converter, laying text out with the test fonts only.
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
