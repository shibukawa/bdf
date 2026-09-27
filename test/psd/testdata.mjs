// Converts the Photoshop test documents layers.psd and artboards.psd in
// converter/psd/testdata into testdata/psd/*.bdf with the Go converter. The
// small documents of the other colour modes and depths are for the Go tests.
import { execFileSync } from "node:child_process";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
for (const name of ["layers", "artboards"]) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", join(root, "converter/psd/testdata", name + ".psd"),
    join(root, "testdata/psd", name + ".bdf")], { cwd: root, stdio: "inherit" });
}
