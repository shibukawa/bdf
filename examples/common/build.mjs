// What the examples and the demo site are built from: the pages' scripts
// bundled with esbuild, the workers (the rendering worker of @bdfkit/render and
// the converter worker, a classic worker that loads Go's wasm_exec.js), the
// converters as wasm (cmd/bdfwasm: one module for PDF, one for the Office
// formats, one for HTML, Markdown and EPUB, one for images, and one without
// converters for thumbnails and search text) and the fonts the converters
// lay text out with.
import { execFile as execFileCb } from "node:child_process";
import { copyFile, mkdir, mkdtemp, readdir, readFile, rm, stat, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { delimiter, dirname, join, relative } from "node:path";
import { fileURLToPath } from "node:url";
import { promisify } from "node:util";
import { build } from "esbuild";

const execFile = promisify(execFileCb);

export const root = dirname(dirname(dirname(fileURLToPath(import.meta.url))));

/** Converter modules: file name and the build tags that select the formats. */
export const MODULES = [
  { file: "bdf-pdf.wasm", tags: "pdfonly" },
  { file: "bdf-office.wasm", tags: "officeonly" },
  { file: "bdf-web.wasm", tags: "webonly" },
  { file: "bdf-image.wasm", tags: "imageonly" },
  { file: "bdf-preview.wasm", tags: "previewonly" },
];

const FONT_EXT = /\.(ttf|otf|ttc|otc)$/i;

/**
 * Bundle scripts into out: entries are [source (from the repository root),
 * name of the output without .js]. Modules by default; format "iife" for a
 * classic worker.
 */
export async function bundle(out, entries, { format = "esm", define = {} } = {}) {
  await mkdir(out, { recursive: true });
  await build({
    bundle: true, outdir: out, sourcemap: true, target: "es2022", logLevel: "info", format, define,
    entryPoints: entries.map(([src, name]) => ({ in: join(root, src), out: name })),
  });
}

/**
 * Copy a page of the examples (from the repository root). Blocks between
 * <!-- site --> and <!-- /site --> are kept when the page is one of the demo
 * site's (they link to its other pages), those between <!-- standalone -->
 * and <!-- /standalone --> when it is built by itself.
 */
export async function copyPage(src, dst, { site = false } = {}) {
  const drop = site ? "standalone" : "site";
  const html = (await readFile(join(root, src), "utf8"))
    .replace(new RegExp(`[ \\t]*<!-- ${drop} -->[\\s\\S]*?<!-- /${drop} -->\\n?`, "g"), "")
    .replace(/[ \t]*<!-- \/?(site|standalone) -->\n?/g, "");
  await mkdir(dirname(dst), { recursive: true });
  await writeFile(dst, html);
}

/** Build the workers into lib: worker.js draws, convert-worker.js (with convert) runs the converter modules. */
export async function buildWorkers(lib, { convert = true } = {}) {
  await bundle(lib, [["packages/render/src/worker.ts", "worker"]]);
  if (convert) await bundle(lib, [["examples/common/convert-worker.ts", "convert-worker"]], { format: "iife" });
}

/**
 * The byte ranges of a font file that the converters' font scan reads (the
 * collection header, the table directories, the name, OS/2 and post tables,
 * and the cmap tables used to check character coverage), merged, as
 * [offset, length] pairs.
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
      const tableName = tag(rec);
      const tableSize = dv.getUint32(rec + 12);
      if (["name", "OS/2", "post"].includes(tableName) || (tableName === "cmap" && tableSize <= 16 * 1024 * 1024)) {
        ranges.push([dv.getUint32(rec + 8), tableSize]);
      }
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
export async function patchGoldmark(dir) {
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

/** Build the converter modules (all of them, or those named) and Go's wasm_exec.js into lib. Requires Go. */
export async function buildModules(lib, files = MODULES.map((m) => m.file)) {
  await mkdir(lib, { recursive: true });
  const env = { ...process.env, GOOS: "js", GOARCH: "wasm" };
  const tmp = await mkdtemp(join(tmpdir(), "bdf-site-"));
  try {
    const mod = await patchGoldmark(tmp);
    await Promise.all(MODULES.filter((m) => files.includes(m.file)).map(async (m) => {
      const t0 = performance.now();
      await execFile("go", ["build", "-tags", `bdf_noconv,${m.tags}`, "-trimpath", `-modfile=${mod}`, "-ldflags=-s -w", "-o", join(lib, m.file), "./cmd/bdfwasm"], { cwd: root, env });
      const { size } = await stat(join(lib, m.file));
      console.log(`${m.file}: ${(size / 1e6).toFixed(1)} MB (${((performance.now() - t0) / 1000).toFixed(1)} s)`);
    }));
  } finally {
    await rm(tmp, { recursive: true, force: true });
  }
  // wasm_exec.js moved from misc/wasm to lib/wasm in Go 1.24
  const goroot = await goEnv("GOROOT");
  for (const dir of ["lib/wasm", "misc/wasm"]) {
    const src = join(goroot, dir, "wasm_exec.js");
    if (await stat(src).catch(() => null)) return copyFile(src, join(lib, "wasm_exec.js"));
  }
  throw new Error(`wasm_exec.js not found under ${goroot}`);
}

/**
 * Publish fonts under dst: the files of the directories BDF_SITE_FONTS
 * lists (separated like PATH; the test fonts of converter/pptx/testdata/fonts
 * by default). Font files are listed in index.json with the ranges the font
 * scan reads, other files (licenses) are copied as they are.
 */
export async function copyFonts(dst) {
  const dirs = (process.env.BDF_SITE_FONTS ?? join(root, "converter/pptx/testdata/fonts")).split(delimiter).filter(Boolean);
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
  const scan = index.reduce((a, e) => a + e.scan.reduce((n, [_, size]) => n + size, 0), 0);
  const size = (n) => n >= 1e6 ? `${(n / 1e6).toFixed(1)} MB` : `${(n / 1024).toFixed(1)} KiB`;
  console.log(`fonts: ${index.length} files, ${size(total)} total; ${size(scan)} prefetched, full files fetched as documents use them`);
}
