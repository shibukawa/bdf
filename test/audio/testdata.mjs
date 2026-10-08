// Converts the audio files of converter/audio/testdata into
// testdata/audio/*.bdf with the Go converter (see test/audio/gen.sh for how
// the files were made).
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, extname, basename } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const dir = join(root, "converter/audio/testdata");
mkdirSync(join(root, "testdata/audio"), { recursive: true });
for (const f of readdirSync(dir)) {
  const out = join(root, "testdata/audio", basename(f, extname(f)) + "-" + extname(f).slice(1) + ".bdf");
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", join(dir, f), out], { cwd: root, stdio: "inherit" });
}
