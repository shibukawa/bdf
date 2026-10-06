// Builds the site published on GitHub Pages:
//
//   index.html, index.ja.html   the top page (home.mjs), which shows a PDF from bdf
//   viewer/                     the demo viewer (examples/viewer)
//   thumbnail/, text/           thumbnails and search text of a file dropped on the page
//   docs/                       the documentation (docs.mjs: docs/*.md as HTML)
//   lib/                        the workers and the converters as wasm (examples/common/build.mjs)
//   fonts/, samples/            the fonts the converters lay text out with, and sample files
//
// Files opened on the site are converted inside the browser.
//
//   node site/build.mjs [--serve] [--out dir]
//
// BDF_SITE_FONTS lists directories (separated like PATH) whose files are
// published under fonts/ (the test fonts by default). Requires Go.
import { execFile as execFileCb } from "node:child_process";
import { copyFile, mkdir, readFile, rm, stat, writeFile } from "node:fs/promises";
import { basename, join, resolve } from "node:path";
import { promisify } from "node:util";
import { gzipSync } from "node:zlib";
import { build } from "esbuild";
import { serve } from "../test/serve.mjs";
import { buildModules, bundle, copyFonts, root } from "../examples/common/build.mjs";
import { buildViewer } from "../examples/viewer/build.mjs";
import { buildDocs, head, header } from "./docs.mjs";
import { home } from "./home.mjs";

const execFile = promisify(execFileCb);
const args = process.argv.slice(2);
const outArg = args.indexOf("--out");
const out = resolve(outArg >= 0 ? args[outArg + 1] : join(root, "site/dist"));

/** Samples offered by the viewer and the other pages: repository path and label. */
const SAMPLES = [
  { path: "site/demo/demo.pdf", label: "PDF (the one of the top page)" },
  { path: "converter/pdf/testdata/reportlab-master.pdf", label: "PDF (3 pages)" },
  { path: "converter/pdf/testdata/chrome-doc.pdf", label: "PDF from Chrome" },
  { path: "converter/docx/testdata/basic.docx", label: "Word" },
  { path: "converter/docx/testdata/vertical.docx", label: "Word (vertical text)" },
  { path: "converter/docx/testdata/math.docx", label: "Word (formulas)" },
  { path: "converter/markdown/testdata/basic.md", label: "Markdown" },
  { path: "converter/markdown/testdata/math.md", label: "Markdown (formulas)" },
  { path: "converter/html/testdata/article.html", label: "HTML (article, reader mode)" },
  { path: "converter/html/testdata/math.html", label: "HTML (formulas in MathML)" },
  { path: "converter/pptx/testdata/features.pptx", label: "PowerPoint" },
  { path: "converter/pptx/testdata/math.pptx", label: "PowerPoint (formulas)" },
  { path: "converter/xlsx/testdata/features.xlsx", label: "Excel" },
  { path: "converter/xlsx/testdata/math.xlsx", label: "Excel (equation objects)" },
  { path: "converter/csv/testdata/japanese.tsv", label: "TSV (Shift_JIS)" },
  { path: "converter/parquet/testdata/basic.parquet", label: "Parquet" },
  { path: "converter/visio/testdata/shapes.vsdx", label: "Visio" },
  { path: "converter/visio/testdata/flow.vdx", label: "Visio XML (.vdx)" },
  { path: "converter/drawio/testdata/multipage.drawio", label: "draw.io (3 pages)" },
  { path: "converter/drawio/testdata/aws.drawio", label: "draw.io (AWS)" },
  { path: "converter/drawio/testdata/math.drawio", label: "draw.io (formulas in labels)" },
  { path: "converter/dxf/testdata/layout.dxf", label: "DXF (model space and a layout)" },
  { path: "converter/jww/testdata/shapes.jww", label: "Jw_cad" },
  { path: "converter/sxf/testdata/shapes.p21", label: "SXF (P21)" },
  { path: "converter/cgm/testdata/shapes.cgm", label: "CGM" },
  { path: "converter/hpgl/testdata/shapes.plt", label: "HP-GL/2 (.plt)" },
  { path: "converter/gerber/testdata/board.zip", label: "Gerber and Excellon (a board, zipped)" },
  { path: "converter/kicad/testdata/demo.zip", label: "KiCad (a project: schematic and board, zipped)" },
  { path: "converter/ai/testdata/artboards.ai", label: "Illustrator (3 artboards)" },
  { path: "converter/psd/testdata/artboards.psd", label: "Photoshop (3 artboards)" },
  { path: "converter/image/testdata/drawing.svg", label: "SVG image" },
  { path: "converter/image/testdata/photo.jpg", label: "JPEG (EXIF)" },
  { path: "converter/epub/testdata/basic.epub", label: "EPUB" },
  { path: "converter/epub/testdata/vertical.epub", label: "EPUB (Japanese, vertical)" },
  { path: "converter/epub/testdata/fixed.epub", label: "EPUB (fixed layout)" },
  { path: "converter/mml/testdata/frere.mml", label: "MML (NES, a round in four parts)" },
  { path: "converter/midi/testdata/twinkle.kar", label: "MIDI (karaoke, with words)" },
  { path: "converter/musicxml/testdata/minuet.musicxml", label: "MusicXML (piano)" },
  { path: "converter/font/testdata/stix.otf", label: "Font (OpenType, STIX Two Text)" },
  { path: "testdata/demo.bdf", label: "BDF" },
];

async function copySamples() {
  const dst = join(out, "samples");
  await mkdir(dst, { recursive: true });
  const index = [];
  for (const s of SAMPLES) {
    const name = basename(s.path);
    if (index.some((e) => e.name === name)) throw new Error(`samples: two files named ${name}`);
    await copyFile(join(root, s.path), join(dst, name));
    index.push({ name, label: s.label });
  }
  await writeFile(join(out, "samples/index.json"), JSON.stringify(index));
}

/** A page of the site whose source is an HTML file: the site's head and header take the place of its marks. */
async function copyTool(name) {
  const file = `${name}/index.html`;
  const html = (await readFile(join(root, "site", file), "utf8"))
    .replace("<!-- head -->", head({ out: file, lang: "en" }))
    .replace("<!-- header -->", header({ lang: "en", out: file, current: name }));
  await writeFile(join(out, file), html);
}

/** The thumbnail and search text pages. */
async function buildTools() {
  for (const name of ["thumbnail", "text"]) {
    await bundle(join(out, name), [[`site/${name}/main.ts`, "main"]], { define: { SITE_ROOT: JSON.stringify("../") } });
    await copyTool(name);
  }
}

/** The input formats of the converters, as bdf generate -h lists them. */
async function formats() {
  const { stderr } = await execFile("go", ["run", "./cmd/bdf", "generate", "-h"], { cwd: root }).catch((e) => e); // -h exits with 2
  const names = [...stderr.slice(stderr.indexOf("input formats:")).matchAll(/^ {2}(\w+) /gm)].map((m) => m[1]);
  if (names.length < 20) throw new Error(`bdf generate -h lists ${names.length} formats`);
  return names;
}

/** The size of the rendering worker as a server sends it: minified and gzipped. */
async function workerSize() {
  const res = await build({ entryPoints: [join(root, "packages/render/src/worker.ts")], bundle: true, minify: true, format: "esm", target: "es2022", write: false, logLevel: "warning" });
  return gzipSync(res.outputFiles[0].contents).length;
}

/** The top page: the document it shows, converted as a server converts it, and the page in each language. */
async function buildHome() {
  const pdf = join(root, "site/demo/demo.pdf"), bdf = join(out, "demo.bdf");
  await execFile("go", ["run", "./cmd/bdf", "generate", "-q", pdf, bdf], { cwd: root });
  await bundle(out, [["site/home.ts", "home"]]);
  await copyFile(join(root, "site/site.css"), join(out, "site.css"));
  const benchmark = JSON.parse(await readFile(join(root, "docs/benchmarks/latest.json"), "utf8"));
  const info = { demo: { pdf: (await stat(pdf)).size, bdf: (await stat(bdf)).size }, formats: (await formats()).length, worker: await workerSize(), benchmark };
  await writeFile(join(out, "index.html"), home("en", info));
  await writeFile(join(out, "index.ja.html"), home("ja", info));
  console.log(`home: demo.bdf ${info.demo.bdf} bytes, ${info.formats} formats, worker ${(info.worker / 1024).toFixed(1)} KB gzipped`);
}

await rm(out, { recursive: true, force: true });
await buildViewer(join(out, "viewer"), { defaultSrc: "", siteRoot: "../", site: true });
await Promise.all([buildModules(join(out, "lib")), copyFonts(join(out, "fonts")), copySamples(), buildDocs(out), buildTools(), buildHome()]);
console.log(`site: ${out}`);

if (args.includes("--serve")) {
  const { port } = await serve(out, Number(process.env.PORT ?? 8766));
  console.log(`site: http://127.0.0.1:${port}/`);
}
