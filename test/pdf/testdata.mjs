// Converts converter/pdf/testdata/*.pdf into testdata/pdf/*.bdf with the Go converter.
import { execFileSync } from "node:child_process";
import { readdirSync } from "node:fs";
import { dirname, join, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
for (const f of readdirSync(join(root, "converter/pdf/testdata"))) {
  if (!f.endsWith(".pdf")) continue;
  const out = join(root, "testdata/pdf", basename(f, ".pdf") + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", join(root, "converter/pdf/testdata", f), out], { cwd: root, stdio: "inherit" });
}
