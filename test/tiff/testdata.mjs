// Converts the TIFF test files in converter/tiff/testdata into
// testdata/tiff/*.bdf with the Go converter.
import { execFileSync } from "node:child_process";
import { readdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const src = join(root, "converter/tiff/testdata");
for (const f of readdirSync(src).filter((f) => f.endsWith(".tif")).sort()) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", join(src, f), join(root, "testdata/tiff", f.replace(/\.tif$/, ".bdf"))],
    { cwd: root, stdio: "inherit" });
}
