// Converts pdf2bdf/testdata/*.pdf into fixtures/pdf/*.bdf with the Go converter.
import { execFileSync } from "node:child_process";
import { readdirSync } from "node:fs";
import { dirname, join, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
for (const f of readdirSync(join(root, "pdf2bdf/testdata"))) {
  if (!f.endsWith(".pdf")) continue;
  const out = join(root, "fixtures/pdf", basename(f, ".pdf") + ".bdf");
  execFileSync("go", ["run", "./cmd/pdf2bdf", "-q", join(root, "pdf2bdf/testdata", f), out], { cwd: root, stdio: "inherit" });
}
