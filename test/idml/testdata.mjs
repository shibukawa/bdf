// Converts the InDesign test documents in converter/idml/testdata (and two
// of the InDesign exports in its simpleidml directory) into
// testdata/idml/*.bdf with the Go converter, laying text out with the test
// fonts only.
import { execFileSync } from "node:child_process";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/pptx/testdata/fonts");
for (const [src, name] of [["basic.idml", "basic"], ["vertical.idml", "vertical"],
  ["simpleidml/interview.idml", "interview"], ["simpleidml/magazineA-courrier-des-lecteurs-3pages.idml", "courrier"]]) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts",
    join(root, "converter/idml/testdata", src), join(root, "testdata/idml", name + ".bdf")], { cwd: root, stdio: "inherit" });
}
