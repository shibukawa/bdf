// Builds the static demo site (published on GitHub Pages): the viewer, the
// converters as wasm (cmd/bdfwasm: one module for PDF, one for the Office
// formats), the fonts the Office converters lay text out with, and sample
// files. Files opened on the site are converted inside the browser.
//
//   node examples/viewer/site.mjs [--serve] [--out dir]
//
// BDF_SITE_FONTS lists directories (separated like PATH) whose files are
// published under fonts/: their font files are listed in fonts/index.json,
// other files (licenses) are copied as they are. It defaults to the test
// fonts of converter/pptx/testdata/fonts. Requires Go.
import { execFile as execFileCb } from "node:child_process";
import { copyFile, mkdir, readdir, readFile, rm, stat, writeFile } from "node:fs/promises";
import { basename, delimiter, join, resolve } from "node:path";
import { promisify } from "node:util";
import { serve } from "../../test/serve.mjs";
import { buildViewer, root } from "./build.mjs";

const execFile = promisify(execFileCb);
const args = process.argv.slice(2);
const outArg = args.indexOf("--out");
const out = resolve(outArg >= 0 ? args[outArg + 1] : join(root, "examples/viewer/.site"));

/** Converter modules: file name and the build tags that select the formats. */
const MODULES = [
  { file: "bdf-pdf.wasm", tags: "pdfonly" },
  { file: "bdf-office.wasm", tags: "officeonly" },
];

/** Samples offered on the start page: repository path and label. */
const SAMPLES = [
  { path: "converter/pdf/testdata/reportlab-master.pdf", label: "PDF (3 pages)" },
  { path: "converter/pdf/testdata/chrome-doc.pdf", label: "PDF from Chrome" },
  { path: "converter/docx/testdata/basic.docx", label: "Word" },
  { path: "converter/docx/testdata/vertical.docx", label: "Word (vertical text)" },
  { path: "converter/pptx/testdata/features.pptx", label: "PowerPoint" },
  { path: "converter/xlsx/testdata/features.xlsx", label: "Excel" },
  { path: "converter/csv/testdata/japanese.tsv", label: "TSV (Shift_JIS)" },
  { path: "converter/visio/testdata/shapes.vsdx", label: "Visio" },
  { path: "converter/visio/testdata/flow.vdx", label: "Visio XML (.vdx)" },
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

async function buildModules() {
  const env = { ...process.env, GOOS: "js", GOARCH: "wasm" };
  await Promise.all(MODULES.map(async (m) => {
    const t0 = performance.now();
    await execFile("go", ["build", "-tags", `bdf_noconv,${m.tags}`, "-trimpath", "-ldflags=-s -w", "-o", join(out, m.file), "./cmd/bdfwasm"], { cwd: root, env });
    const { size } = await stat(join(out, m.file));
    console.log(`${m.file}: ${(size / 1e6).toFixed(1)} MB (${((performance.now() - t0) / 1000).toFixed(1)} s)`);
  }));
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
await Promise.all([buildModules(), copyFonts(), copySamples()]);
console.log(`site: ${out}`);

if (args.includes("--serve")) {
  const { port } = await serve(out, Number(process.env.PORT ?? 8766));
  console.log(`site: http://127.0.0.1:${port}/`);
}
