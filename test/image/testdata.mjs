// Converts the images of converter/image/testdata into testdata/image/*.bdf
// with the Go converter (the images are stored as they are; see
// test/image/gen.sh for how they were made).
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, extname, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const dir = join(root, "converter/image/testdata");
mkdirSync(join(root, "testdata/image"), { recursive: true });
for (const f of readdirSync(dir)) {
  const out = join(root, "testdata/image", basename(f, extname(f)) + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", join(dir, f), out], { cwd: root, stdio: "inherit" });
}
