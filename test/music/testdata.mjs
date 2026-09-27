// Converts the test scores of the music converters (converter/mml,
// converter/midi, converter/musicxml) into testdata/music/*.bdf with the Go
// converters, setting the words with the test fonts only.
import { execFileSync } from "node:child_process";
import { mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/pptx/testdata/fonts");
mkdirSync(join(root, "testdata/music"), { recursive: true });
// [input, output name]
const inputs = [
  ["converter/mml/testdata/frere.mml", "frere-mml"],
  ["converter/mml/testdata/ode.mml", "ode-mml"],
  ["converter/midi/testdata/twinkle.kar", "twinkle-kar"],
  ["converter/musicxml/testdata/minuet.musicxml", "minuet-musicxml"],
  ["converter/musicxml/testdata/ode_to_joy.mxl", "ode-mxl"],
];
for (const [src, name] of inputs) {
  execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts",
    join(root, src), join(root, "testdata/music", name + ".bdf")], { cwd: root, stdio: "inherit" });
}
