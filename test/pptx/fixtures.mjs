// Converts converter/pptx/testdata/*.pptx into fixtures/pptx/*.bdf with the
// Go converter, laying text out with the test fonts only.
import { execFileSync } from "node:child_process";
import { readdirSync } from "node:fs";
import { dirname, join, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const dir = join(root, "converter/pptx/testdata");
for (const f of readdirSync(dir)) {
  if (!f.endsWith(".pptx")) continue;
  const out = join(root, "fixtures/pptx", basename(f, ".pptx") + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", join(dir, "fonts"), "-no-system-fonts", join(dir, f), out], { cwd: root, stdio: "inherit" });
}
