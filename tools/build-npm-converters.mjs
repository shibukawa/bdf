// Build the npm presets with TinyGo. A custom set can be built as:
// node tools/build-npm-converters.mjs --formats pdf,epub,docx --out public/converter.wasm
import { execFile as execFileCb } from "node:child_process";
import { copyFile, mkdir, mkdtemp, readFile, readdir, rm, stat, writeFile } from "node:fs/promises";
import { basename, dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { promisify } from "node:util";
import { build } from "esbuild";
import { patchGoldmark, patchPdfcpuForTinyGo, patchRunewidthForTinyGo } from "../examples/common/build.mjs";

const execFile = promisify(execFileCb);
const root = dirname(dirname(fileURLToPath(import.meta.url)));
const presets = [
  ["pdf-epub", "npm_pdfepub", "packages/convert-pdf-epub/dist/converter.wasm"],
  ["office", "npm_office", "packages/convert-office/dist/converter.wasm"],
  ["all", "npm_all", "packages/convert-all/dist/converter.wasm"],
];
const allowed = new Set((await readdir(join(root, "converter"), { withFileTypes: true }))
  .filter((entry) => entry.isDirectory() && !entry.name.startsWith("internal"))
  .map((entry) => entry.name));
const musicFormats = new Set(["midi", "mml", "musicxml"]);

async function buildWasm(dir, tags, output, mod) {
  await mkdir(dirname(output), { recursive: true });
  const args = ["build", "-target", "wasm", "-no-debug", "-tags", `bdf_noconv${tags ? `,${tags}` : ""}`, "-o", output];
  // The broad preset and arbitrary combinations need more room during init.
  const stack = process.env.BDF_TINYGO_STACK_SIZE ?? (["npm_all", "npm_custom"].includes(tags) ? "1MB" : "");
  if (stack) args.push(`-stack-size=${stack}`);
  args.push(dir);
  const t0 = performance.now();
  await execFile("tinygo", args, { cwd: root, env: mod ? { ...process.env, GOFLAGS: `-modfile=${mod}` } : process.env, maxBuffer: 8 << 20 });
  const { size } = await stat(output);
  console.log(`${basename(output)}: ${(size / 1e6).toFixed(1)} MB (${((performance.now() - t0) / 1000).toFixed(1)} s)`);
}

const args = process.argv.slice(2);
const selected = args.indexOf("--formats");
const outputAt = args.indexOf("--out");
const presetAt = args.indexOf("--preset");
const withoutMusic = args.includes("--without-music");
if ((selected >= 0) !== (outputAt >= 0) || (selected >= 0 && (!args[selected + 1] || !args[outputAt + 1]))) {
  throw new Error("use --formats pdf,epub,... --out path/to/converter.wasm together");
}
if (withoutMusic && selected < 0) throw new Error("--without-music requires --formats (use --formats all for every non-music format)");
if (presetAt >= 0 && (selected >= 0 || !presets.some(([name]) => name === args[presetAt + 1]))) {
  throw new Error("use --preset pdf-epub, office or all without --formats");
}

const tmp = await mkdtemp(join(root, "cmd/bdfnpm-"));
let tmpPatch;
try {
  let mod;
  tmpPatch = await mkdtemp(join(root, "cmd/bdfmod-"));
  if (selected < 0 || args[selected + 1].split(",").some((name) => ["all", "html", "markdown"].includes(name.trim().toLowerCase()))) {
    mod = await patchGoldmark(tmpPatch);
  }
  mod = await patchPdfcpuForTinyGo(tmpPatch, mod);
  mod = await patchRunewidthForTinyGo(tmpPatch, mod);
  if (selected >= 0) {
    const requested = args[selected + 1].split(",").map((name) => name.trim().toLowerCase()).map((name) => name === "tsv" ? "csv" : name);
    if (requested.some((name) => !allowed.has(name))) throw new Error(`unknown format in ${args[selected + 1]}`);
    const names = [...new Set(requested.flatMap((name) => name === "all" ? [...allowed].filter((format) => format !== "all") : [name]))]
      .filter((name) => !withoutMusic || !musicFormats.has(name));
    if (!names.length) throw new Error("no formats left to build");
    console.log(`formats: ${names.join(", ")}`);
    for (const file of (await readdir(join(root, "cmd/bdfwasm"))).filter((name) => name.endsWith(".go") && !name.endsWith("_test.go"))) {
      await copyFile(join(root, "cmd/bdfwasm", file), join(tmp, file));
    }
    const imports = names.map((name) => `\t_ "github.com/shibukawa/bdf/converter/${name}"`).join("\n");
    await writeFile(join(tmp, "custom.go"), `//go:build js && wasm && npm_custom\n\npackage main\n\nimport (\n${imports}\n)\n`);
    await buildWasm(`./cmd/${basename(tmp)}`, "npm_custom", resolve(root, args[outputAt + 1]), mod);
  } else {
    // Put TinyGo's runtime in the Worker so it is one asset for bundlers.
    const dist = join(root, "packages/convert/dist");
    await mkdir(dist, { recursive: true });
    const { stdout } = await execFile("tinygo", ["env", "TINYGOROOT"], { cwd: root });
    const runtime = await readFile(join(stdout.trim(), "targets/wasm_exec.js"), "utf8");
    await build({ entryPoints: [join(root, "packages/convert/src/worker.ts")], outfile: join(dist, "worker.js"), bundle: true, format: "iife", target: "es2022", banner: { js: runtime } });
    await rm(join(dist, "wasm_exec.js"), { force: true });
    for (const [name, tag, file] of presets.filter(([name]) => presetAt < 0 || name === args[presetAt + 1])) {
      await buildWasm("./cmd/bdfwasm", tag, join(root, file), mod);
    }
  }
} finally {
  await rm(tmp, { recursive: true, force: true });
  if (tmpPatch) await rm(tmpPatch, { recursive: true, force: true });
}
