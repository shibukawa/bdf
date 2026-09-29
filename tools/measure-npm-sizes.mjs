// Browser payload sizes for the four npm installation profiles.
// Counts the runtime .js and .wasm files in each package's dist directory,
// once per profile; source maps and TypeScript declarations are not sent.
import { gzipSync } from "node:zlib";
import { readFile, readdir } from "node:fs/promises";
import { join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { build } from "esbuild";

const root = fileURLToPath(new URL("..", import.meta.url));
const common = ["core", "render", "viewer"];
const profiles = [
  ["core + render（低レベル）", ["core", "render"]],
  ["レンダラーのみ", common],
  ["PDF + EPUB（レンダラーなし）", ["convert", "convert-pdf-epub"]],
  ["PDF + EPUB", [...common, "convert", "convert-pdf-epub"]],
  ["Office + PDF + EPUB（レンダラーなし）", ["convert", "convert-office"]],
  ["Office + PDF + EPUB", [...common, "convert", "convert-office"]],
  ["All（レンダラーなし）", ["convert", "convert-all"]],
  ["All", [...common, "convert", "convert-all"]],
];

async function files(directory) {
  const found = [];
  for (const item of await readdir(directory, { withFileTypes: true })) {
    const path = join(directory, item.name);
    if (item.isDirectory()) found.push(...await files(path));
    else if (item.isFile() && /\.(js|wasm)$/.test(item.name)) found.push(path);
  }
  return found;
}

const size = (bytes) => bytes >= 1024 * 1024
  ? `${(bytes / (1024 * 1024)).toFixed(2)} MiB`
  : `${(bytes / 1024).toFixed(1)} KiB`;

console.log("| 構成 | 無圧縮 | gzip | 対象ファイル |");
console.log("|---|---:|---:|---:|");
for (const [name, packages] of profiles) {
  const paths = (await Promise.all(packages.map((pkg) => files(join(root, "packages", pkg, "dist"))))).flat();
  if (!paths.length || (packages.includes("convert") && !paths.some((path) => path.endsWith("converter.wasm")))) {
    throw new Error(`${name}: build first with npm run build && npm run build:wasm`);
  }
  let raw = 0, gzip = 0;
  for (const path of paths) {
    const content = await readFile(path);
    raw += content.length;
    gzip += gzipSync(content, { level: 9 }).length;
  }
  console.log(`| ${name} | ${size(raw)} | ${size(gzip)} | ${paths.length} |`);
}

// The package table counts every dist file. Compare the application-side JS
// that a bundler actually reaches with and without the optional music entry.
console.log("\nViewer application JS (minified; render Worker and wasm excluded):");
console.log("| 音楽機能 | 無圧縮 | gzip |");
console.log("|---|---:|---:|");
for (const music of [false, true]) {
  const entry = `import { Viewer } from './packages/viewer/dist/index.js';\n`
    + (music ? `import { pianoRoll, MusicPlayer } from './packages/viewer/dist/music.js';\n` : "")
    + `globalThis.bdfViewer = Viewer;\n`
    + (music ? `globalThis.bdfMusic = { pianoRoll, MusicPlayer };\n` : "");
  const result = await build({ stdin: { contents: entry, resolveDir: root, sourcefile: "bdf-size-entry.js" }, bundle: true, write: false, platform: "browser", format: "esm", minify: true });
  const js = result.outputFiles[0].contents;
  console.log(`| ${music ? "on" : "off"} | ${size(js.length)} | ${size(gzipSync(js, { level: 9 }).length)} |`);
}

const customAt = process.argv.indexOf("--custom-wasm");
if (customAt >= 0) {
  if (!process.argv[customAt + 1]) throw new Error("--custom-wasm requires a file path");
  const wasm = await readFile(resolve(process.argv[customAt + 1]));
  console.log(`\nCustom wasm: ${size(wasm.length)} raw; ${size(gzipSync(wasm, { level: 9 }).length)} gzip`);
}
