# Integrating the npm packages

Use `@bdfkit/viewer` alone when the server already converts files to BDF. Add one browser converter preset when the browser should convert input files. These packages are npm workspaces in this repository; the `npm install` commands below show how to use them after publication. Before publishing them, run `npm run build && npm run build:wasm`; the resulting TinyGo wasm file is included in each converter package, so consumers do not need Go or TinyGo.

For focused guides, see [Integrating into a frontend](integrating-to-frontend.md), [Integrating into a server](integrate-to-server.md), and [Build a Wasm runtime](build-wasm-runtime.md).

| Profile | Install | Input |
|---|---|---|
| Renderer only | `@bdfkit/viewer` | BDF produced on the server |
| PDF + EPUB | `@bdfkit/viewer`, `@bdfkit/convert-pdf-epub` | PDF, EPUB |
| Office + PDF + EPUB | `@bdfkit/viewer`, `@bdfkit/convert-office` | Word, Excel, PowerPoint, Visio, CSV, TSV, Parquet, PDF, EPUB |
| All | `@bdfkit/viewer`, `@bdfkit/convert-all` | Every browser converter in this repository |

To omit the renderer, install just the converter package and leave out `@bdfkit/viewer`. The default viewer import also leaves out piano roll and playback; import `@bdfkit/viewer/music` only when needed. To omit MML, MIDI, and MusicXML converters from the All profile, [build custom Wasm](build-wasm-runtime.md#tinygo-select-your-formats) with `--formats all --without-music`.

`@bdfkit/viewer` depends on `@bdfkit/render` and `@bdfkit/core`; the converter presets depend on `@bdfkit/convert`. Your bundler must serve the `worker.js` and `converter.wasm` files referenced by `new URL(..., import.meta.url)` in these packages. The TinyGo runtime is embedded in the converter Worker. Allow Workers and WebAssembly in your content security policy.

## Renderer only: server conversion

```sh
npm install @bdfkit/viewer
```

Give the host a height and open a BDF document served by your application:

```html
<div id="preview" style="width:100%;height:70vh"></div>
```

```ts
import { Viewer } from "@bdfkit/viewer";

const viewer = new Viewer(document.querySelector("#preview")!, { layout: "single" });
await viewer.open({ kind: "single", url: "/documents/report.bdf" });
// Or: await viewer.open({ kind: "split", base: "/documents/report/" });
viewer.setLayout("spread");
// Call viewer.destroy() when removing the view.
```

For applications that provide their own page placement, use `@bdfkit/render` and `@bdfkit/core` directly. `@bdfkit/render` provides the drawing Worker; `@bdfkit/viewer` adds a scrolling document surface and presentation modes.

## Convert in the browser

```sh
npm install @bdfkit/viewer @bdfkit/convert-office
```

```ts
import { Viewer } from "@bdfkit/viewer";
import { createOfficeConverter } from "@bdfkit/convert-office";

const viewer = new Viewer(document.querySelector("#preview")!, {
  mode: "embedded", // or "lightbox"
  layout: "single", // or "spread" / "continuous"
});
const converter = createOfficeConverter();
const file = (document.querySelector("input[type=file]") as HTMLInputElement).files![0];
await viewer.openFile(converter, await file.arrayBuffer(), { name: file.name, fonts: "/fonts/" });

viewer.setLayout("spread");
viewer.setView(viewer.document?.views[0]?.id);
viewer.scrollToPage(3); // zero-based
// In lightbox mode, call viewer.show() and viewer.hide().

viewer.destroy();
converter.terminate();
```

The host application owns its file picker, toolbar, view tabs, zoom controls and lightbox close button. Use `setView`, `setLayout`, `setZoom`, `show` and `hide` to connect them. For formats that support page streaming, the viewer shows the outline first and fills pages as they convert. It emits `documentchange`, `viewchange`, `conversionprogress`, `conversiondone` and `error` events. It places and scrolls pages and renders sheet viewports.

## Payload sizes

Measured on September 30, 2026 with TinyGo 0.42.0. Run `npm run build && npm run build:wasm && npm run size:npm` to repeat the measurements. The table counts runtime `.js` and `.wasm` files in each profile's `dist` directories, with dependencies counted once. gzip is the sum of individual files compressed with gzip level 9. Type declarations, source maps, fonts, documents, HTML and CSS are excluded. The profiles without a renderer omit `@bdfkit/viewer`. A production bundler may remove unused JavaScript.

<!-- npm-size-table -->
| Profile | Uncompressed | gzip | Files |
|---|---:|---:|---:|
| core + render (low level) | 356.7 KiB | 108.5 KiB | 29 |
| Renderer only | 371.4 KiB | 112.7 KiB | 31 |
| PDF + EPUB (no renderer) | 13.79 MiB | 6.09 MiB | 4 |
| PDF + EPUB | 14.15 MiB | 6.20 MiB | 35 |
| Office + PDF + EPUB (no renderer) | 14.74 MiB | 6.44 MiB | 4 |
| Office + PDF + EPUB | 15.10 MiB | 6.55 MiB | 35 |
| All (no renderer) | 20.11 MiB | 8.81 MiB | 4 |
| All | 20.48 MiB | 8.92 MiB | 35 |
<!-- /npm-size-table -->

The converter Wasm files are built with TinyGo 0.42.0. `tools/build-npm-converters.mjs` invokes `tinygo build -target wasm -no-debug`; TinyGo's default `-opt=z` size optimization is used. The Wasm is a statically linked converter module, not only a small language runtime: each selected format contributes its parser, layout code, rasterizer, and lookup data. The `All` preset therefore contains every converter and cannot be reduced by the frontend bundler's JavaScript tree shaking.

The renderer-only profile avoids converter Wasm entirely. When browser conversion is needed, use a narrower preset or build a custom module with `--formats`; `--without-music` removes the music converters from a custom build. For network transfer, serve the `.wasm` with HTTP compression such as Brotli. The table above uses gzip level 9 so that the measurements remain reproducible across environments.

The package file totals include optional entries even when an app does not import them. For minified application JS, the viewer without music was 11.7 KiB (4.4 KiB gzip); importing piano roll and `MusicPlayer` from `@bdfkit/viewer/music` raised this to 29.5 KiB (10.9 KiB gzip). These JS bundle measurements exclude the render Worker and Wasm.

## Choose formats yourself

Build a custom wasm module from this repository. The names come from `converter/`; `tsv` aliases `csv`.

```sh
npm ci
node tools/build-npm-converters.mjs --formats pdf,epub,docx,csv --out ./public/my-converter.wasm
```

```ts
import { createConverter } from "@bdfkit/convert";
const converter = createConverter(new URL("/my-converter.wasm", location.href));
console.log(await converter.ready); // formats included in this wasm
```

For inputs that need fonts, pass a font directory URL as the `fonts` conversion option. It needs an `index.json`; see `copyFonts` in `examples/common/build.mjs` for generating one.

Before publishing, inspect the assets in a package:

```sh
npm run build
npm run build:wasm
npm pack --dry-run --workspace @bdfkit/convert-office
```

Publish `@bdfkit/convert`, the chosen preset, `@bdfkit/viewer`, `@bdfkit/render` and `@bdfkit/core` at the same version.

## Publish from GitHub Actions

`.github/workflows/publish.yml` publishes all seven packages when a `v*` tag is pushed. It checks that the tag version matches every package, runs the TypeScript and Wasm builds and tests, inspects the package tarballs, and then publishes in dependency order.

Because the packages are unpublished initially, the first release needs an npm granular access token. Create a token with read and write access to the `@bdfkit` packages (and permission to bypass 2FA for publishing), then save it as the repository Actions secret `NPM_TOKEN`. From a clean checkout, bump all package versions together, commit, tag and push:

```sh
git tag v0.1.0
git push origin main v0.1.0
```

The workflow uses npm provenance (`--provenance`) and requests the GitHub OIDC identity token. After the first packages exist on npm, configure GitHub Actions as a Trusted Publisher for each package and remove `NPM_TOKEN`; subsequent releases can then use OIDC without a long-lived npm token.
