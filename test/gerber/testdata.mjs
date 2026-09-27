// Converts the Gerber test data in converter/gerber/testdata (the board's
// archive and the features file) into testdata/gerber/*.bdf with the Go
// converter. Gerber files have no text, so no fonts are needed.
import { execFileSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
mkdirSync(join(root, "testdata/gerber"), { recursive: true });
for (const [src, out] of [["board.zip", "board"], ["features.gbr", "features"]]) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", join(root, "converter/gerber/testdata", src),
    join(root, "testdata/gerber", out + ".bdf")], { cwd: root, stdio: "inherit" });
}
