// Converts converter/markdown/testdata/*.md into testdata/markdown/*.bdf and
// converter/html/testdata/*.html into testdata/html/*.bdf with the Go
// converter. Text is laid out with the metrics of the Word converter's test
// fonts and refers to the fonts by name (the default for HTML and
// Markdown), so viewers draw it with their own fonts; nothing is fetched.
import { execFileSync } from "node:child_process";
import { mkdirSync, readdirSync } from "node:fs";
import { dirname, join, basename, extname } from "node:path";
import { fileURLToPath } from "node:url";

const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));
const fonts = join(root, "converter/docx/testdata/fonts");
for (const [kind, ext] of [["markdown", ".md"], ["html", ".html"]]) {
  const dir = join(root, "converter", kind, "testdata");
  mkdirSync(join(root, "testdata", kind), { recursive: true });
  for (const f of readdirSync(dir)) {
    if (extname(f) !== ext) continue;
    const out = join(root, "testdata", kind, basename(f, ext) + ".bdf");
    execFileSync("go", ["run", "./cmd/bdf", "generate", "-q", "-font-dir", fonts, "-no-system-fonts", "-param", "remote=false", join(dir, f), out],
      { cwd: root, stdio: "inherit" });
  }
}
