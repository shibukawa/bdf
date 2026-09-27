// Converts the Illustrator test file converter/ai/testdata/artboards.ai into
// testdata/ai/artboards.bdf with the Go converter. nopdf.ai and legacy.ai are
// the files the converter refuses; the Go tests check them.
import { execFileSync } from "node:child_process";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", join(root, "converter/ai/testdata/artboards.ai"),
  join(root, "testdata/ai/artboards.bdf")], { cwd: root, stdio: "inherit" });
