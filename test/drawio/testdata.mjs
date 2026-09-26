// Converts converter/drawio/testdata/*.drawio into testdata/drawio/*.bdf with
// the Go converter, laying text out with the test fonts only (those of the
// PowerPoint converter), so that the output does not depend on the machine.
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const dir = join(root, "converter/drawio/testdata");
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/drawio"), { recursive: true });
for (const f of readdirSync(dir)) {
  if (!f.endsWith(".drawio")) continue;
  const out = join(root, "testdata/drawio", basename(f, ".drawio") + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts", join(dir, f), out], { cwd: root, stdio: "inherit" });
}
