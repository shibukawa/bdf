// Builds the static demo site (published on GitHub Pages): the viewer, the
// converters as wasm (cmd/bdfwasm: one module for PDF, one for the Office
// formats, one for HTML and Markdown, one for images), the fonts the Office,
// HTML and Markdown converters lay text out with, sample files, and the
// documentation (docs.mjs: the READMEs and docs/ as HTML under docs/).
// Files opened on the site are converted inside the browser.
//
//   node examples/viewer/site.mjs [--serve] [--out dir]
//
// BDF_SITE_FONTS lists directories (separated like PATH) whose files are
// published under fonts/: their font files are listed in fonts/index.json,
// other files (licenses) are copied as they are. It defaults to the test
// fonts of converter/pptx/testdata/fonts. Requires Go.
import { execFile as execFileCb } from "node:child_process";
import { copyFile, mkdir, mkdtemp, readdir, readFile, rm, stat, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { basename, delimiter, dirname, join, relative, resolve } from "node:path";
import { promisify } from "node:util";
import { serve } from "../../test/serve.mjs";
import { buildViewer, root } from "./build.mjs";
import { buildDocs } from "./docs.mjs";

const execFile = promisify(execFileCb);
const args = process.argv.slice(2);
const outArg = args.indexOf("--out");
const out = resolve(outArg >= 0 ? args[outArg + 1] : join(root, "examples/viewer/.site"));

/** Converter modules: file name and the build tags that select the formats. */
const MODULES = [
  { file: "bdf-pdf.wasm", tags: "pdfonly" },
  { file: "bdf-office.wasm", tags: "officeonly" },
  { file: "bdf-web.wasm", tags: "webonly" },
  { file: "bdf-image.wasm", tags: "imageonly" },
];

/** Samples offered on the start page: repository path and label. */
const SAMPLES = [
  { path: "converter/pdf/testdata/reportlab-master.pdf", label: "PDF (3 pages)" },
  { path: "converter/pdf/testdata/chrome-doc.pdf", label: "PDF from Chrome" },
  { path: "converter/docx/testdata/basic.docx", label: "Word" },
  { path: "converter/docx/testdata/vertical.docx", label: "Word (vertical text)" },
  { path: "converter/docx/testdata/math.docx", label: "Word (formulas)" },
  { path: "converter/pptx/testdata/features.pptx", label: "PowerPoint" },
  { path: "converter/xlsx/testdata/features.xlsx", label: "Excel" },
  { path: "converter/csv/testdata/japanese.tsv", label: "TSV (Shift_JIS)" },
  { path: "converter/parquet/testdata/basic.parquet", label: "Parquet" },
  { path: "converter/visio/testdata/shapes.vsdx", label: "Visio" },
  { path: "converter/visio/testdata/flow.vdx", label: "Visio XML (.vdx)" },
  { path: "converter/drawio/testdata/multipage.drawio", label: "draw.io (3 pages)" },
  { path: "converter/drawio/testdata/aws.drawio", label: "draw.io (AWS)" },
  { path: "converter/dxf/testdata/layout.dxf", label: "DXF (model space and a layout)" },
  { path: "converter/jww/testdata/shapes.jww", label: "Jw_cad" },
  { path: "converter/sxf/testdata/shapes.p21", label: "SXF (P21)" },
  { path: "converter/cgm/testdata/shapes.cgm", label: "CGM" },
  { path: "converter/hpgl/testdata/shapes.plt", label: "HP-GL/2 (.plt)" },
  { path: "converter/gerber/testdata/board.zip", label: "Gerber and Excellon (a board, zipped)" },
  { path: "converter/ai/testdata/artboards.ai", label: "Illustrator (3 artboards)" },
  { path: "converter/psd/testdata/artboards.psd", label: "Photoshop (3 artboards)" },
  { path: "converter/image/testdata/drawing.svg", label: "SVG image" },
  { path: "converter/image/testdata/photo.jpg", label: "JPEG (EXIF)" },
  { path: "converter/epub/testdata/basic.epub", label: "EPUB" },
  { path: "converter/epub/testdata/vertical.epub", label: "EPUB (Japanese, vertical)" },
  { path: "converter/epub/testdata/fixed.epub", label: "EPUB (fixed layout)" },
  { path: "testdata/demo.bdf", label: "bdf" },
];

const FONT_EXT = /\.(ttf|otf|ttc|otc)$/i;

/**
 * The byte ranges of a font file that the converters' font scan reads (the
 * collection header, the table directories and the name, OS/2 and post
 * tables), merged, as [offset, length] pairs.
 */
export function scanRanges(b) {
  const dv = new DataView(b.buffer, b.byteOffset, b.byteLength);
  const tag = (o) => String.fromCharCode(b[o], b[o + 1], b[o + 2], b[o + 3]);
  const ranges = [];
  let faces = [0];
  if (tag(0) === "ttcf") {
    const n = dv.getUint32(8);
    ranges.push([0, 12 + 4 * n]);
    faces = Array.from({ length: n }, (_, i) => dv.getUint32(12 + 4 * i));
  }
  for (const off of faces) {
    const n = dv.getUint16(off + 4);
    ranges.push([off, 12 + 16 * n]);
    for (let i = 0; i < n; i++) {
      const rec = off + 12 + 16 * i;
      if (["name", "OS/2", "post"].includes(tag(rec))) ranges.push([dv.getUint32(rec + 8), dv.getUint32(rec + 12)]);
    }
  }
  ranges.sort((a, b) => a[0] - b[0]);
  const merged = [];
  for (const [o, n] of ranges) {
    const last = merged.at(-1);
    // close ranges become one request
    if (last && o <= last[0] + last[1] + 4096) last[1] = Math.max(last[1], o + n - last[0]);
    else merged.push([o, n]);
  }
  return merged.filter(([o, n]) => n > 0 && o + n <= b.length);
}

async function goEnv(name) {
  return (await execFile("go", ["env", name], { cwd: root })).stdout.trim();
}

/**
 * Write into dir a copy of goldmark (the Markdown converter's parser) whose
 * linkify extension compiles shallower patterns, and a go.mod that replaces
 * goldmark with it; return that go.mod's path, for go build -modfile.
 * regexp/syntax turns [...]{1,256} into 255 nested groups and compiles them
 * recursing once a group, deeper than Safari lets a worker's stack go: the
 * web module, which compiles the patterns when it starts, stopped with
 * "Maximum call stack size exceeded" before it set bdfConverter. [...]+
 * finds the same links, but for host names longer than 256 characters. (A
 * -overlay cannot replace files in the module cache.)
 */
async function patchGoldmark(dir) {
  const { Dir: src } = JSON.parse((await execFile("go", ["mod", "download", "-json", "github.com/yuin/goldmark"], { cwd: root })).stdout);
  const dst = join(dir, "goldmark");
  // copied file by file: the module cache's directories are read-only, and so would be the copy's
  for (const e of await readdir(src, { recursive: true, withFileTypes: true })) {
    if (!e.isFile()) continue;
    const to = join(dst, relative(src, join(e.parentPath, e.name)));
    await mkdir(dirname(to), { recursive: true });
    await writeFile(to, await readFile(join(e.parentPath, e.name)));
  }
  const linkify = join(dst, "extension/linkify.go");
  const code = await readFile(linkify, "utf8");
  if (code.split("]{1,256}").length !== 3) throw new Error(`${join(src, "extension/linkify.go")}: not the two patterns patchGoldmark rewrites`);
  await writeFile(linkify, code.replaceAll("]{1,256}", "]+"));
  const mod = join(dir, "go.mod");
  await writeFile(mod, `${await readFile(join(root, "go.mod"), "utf8")}\nreplace github.com/yuin/goldmark => ${dst}\n`);
  await copyFile(join(root, "go.sum"), join(dir, "go.sum"));
  return mod;
}

async function buildModules() {
  const env = { ...process.env, GOOS: "js", GOARCH: "wasm" };
  const tmp = await mkdtemp(join(tmpdir(), "bdf-site-"));
  try {
    const mod = await patchGoldmark(tmp);
    await Promise.all(MODULES.map(async (m) => {
      const t0 = performance.now();
      await execFile("go", ["build", "-tags", `bdf_noconv,${m.tags}`, "-trimpath", `-modfile=${mod}`, "-ldflags=-s -w", "-o", join(out, m.file), "./cmd/bdfwasm"], { cwd: root, env });
      const { size } = await stat(join(out, m.file));
      console.log(`${m.file}: ${(size / 1e6).toFixed(1)} MB (${((performance.now() - t0) / 1000).toFixed(1)} s)`);
    }));
  } finally {
    await rm(tmp, { recursive: true, force: true });
  }
  // wasm_exec.js moved from misc/wasm to lib/wasm in Go 1.24
  const goroot = await goEnv("GOROOT");
  for (const dir of ["lib/wasm", "misc/wasm"]) {
    const src = join(goroot, dir, "wasm_exec.js");
    if (await stat(src).catch(() => null)) return copyFile(src, join(out, "wasm_exec.js"));
  }
  throw new Error(`wasm_exec.js not found under ${goroot}`);
}

async function copyFonts() {
  const dirs = (process.env.BDF_SITE_FONTS ?? join(root, "converter/pptx/testdata/fonts")).split(delimiter).filter(Boolean);
  const dst = join(out, "fonts");
  await mkdir(dst, { recursive: true });
  const index = [];
  for (const dir of dirs) {
    for (const name of (await readdir(dir)).sort()) {
      const src = join(dir, name);
      if (!(await stat(src)).isFile()) continue;
      if (index.some((e) => e.name === name)) throw new Error(`fonts: ${name} is in more than one directory`);
      await copyFile(src, join(dst, name));
      if (!FONT_EXT.test(name)) continue;
      const data = await readFile(src);
      index.push({ name, size: data.length, scan: scanRanges(data) });
    }
  }
  await writeFile(join(dst, "index.json"), JSON.stringify(index));
  const total = index.reduce((a, e) => a + e.size, 0);
  console.log(`fonts: ${index.length} files, ${(total / 1e6).toFixed(1)} MB (fetched as documents use them)`);
}

async function copySamples() {
  const dst = join(out, "samples");
  await mkdir(dst, { recursive: true });
  const index = [];
  for (const s of SAMPLES) {
    const name = basename(s.path);
    await copyFile(join(root, s.path), join(dst, name));
    index.push({ name, label: s.label });
  }
  await writeFile(join(dst, "index.json"), JSON.stringify(index));
}

await rm(out, { recursive: true, force: true });
await buildViewer(out, { defaultSrc: "" });
await Promise.all([buildModules(), copyFonts(), copySamples(), buildDocs(out)]);
console.log(`site: ${out}`);

if (args.includes("--serve")) {
  const { port } = await serve(out, Number(process.env.PORT ?? 8766));
  console.log(`site: http://127.0.0.1:${port}/`);
}
