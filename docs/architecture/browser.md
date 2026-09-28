# In the browser

Two workers do the work; the main thread only places bitmaps and a transparent text layer.

## The rendering worker (`@bdf/render`)

`packages/render/src/worker.ts` loads a document, decodes its parts and draws pages, sheet tiles and continuous-layout bands into `ImageBitmap`s on an `OffscreenCanvas`. `BdfWorkerClient` (`packages/render/src/client.ts`) is the main-thread handle to it — one `postMessage`-based call per method:

```ts
import { BdfWorkerClient } from "@bdf/render";

const client = new BdfWorkerClient(new Worker("./worker.js", { type: "module" }));
const manifest = await client.open({ kind: "single", url, range: true }); // or {kind: "split", base} / {kind: "buffer", buffer}
const bitmap = await client.page(manifest.views[0].id, 0, scale);         // scale: device pixels per unit
const content = await client.content(manifest.views[0].id, 0);           // structured text, for buildTextLayer
```

`open` takes a single file (optionally read by HTTP Range requests), a split-form directory, or an in-memory buffer (a document just converted in the browser); an encrypted document rejects with a `BdfWorkerError` coded `"password-required"`, and the worker keeps it locked until `unlock(password)` succeeds. `page`/`continuous`/`sheet` draw bitmaps; `text`/`continuousText`/`content`/`sheetContent` return the text runs `buildTextLayer` (from `@bdf/render`) turns into the selectable, accessible DOM layer placed over each bitmap. `search`/`locate` run full-text search inside the worker and return hit rectangles. `play` returns a view's music as a Standard MIDI File plus cues, for `MusicPlayer`. SVG images can't be decoded inside a worker, so the worker asks the main thread to rasterize them (`RasterizeRequest`/`RasterizeResponse`), which `BdfWorkerClient` handles automatically via `domSvgRasterizer`.

A viewer's own job, on top of this, is placing bitmaps and text layers as pages scroll into view, driving zoom and search UI, and — if it wants page-turning, cell selection, or music playback — building those interactions. [`examples/miniviewer`](https://github.com/shibukawa/bdf/tree/main/examples/miniviewer) is a minimal one (~600 lines) that does the first part only; [`examples/viewer`](https://github.com/shibukawa/bdf/tree/main/examples/viewer) (the demo site's full viewer) adds the rest.

## The converter worker (`cmd/bdfwasm`)

For in-browser conversion, `cmd/bdfwasm` builds the same Go converter packages as WebAssembly. It's a classic (non-module) worker, since it needs Go's `wasm_exec.js`:

```
bdfConverter.convert(data: Uint8Array, options?: {format?, password?, fonts?, name?})
  → Promise<{bdf, format, summary, warnings, protected}>
bdfConverter.open(data, options?) → Promise<{bdf, format, pages, warnings, stream?}>
  // pages > 0: bdf is the outline (every page sized, no layers yet); stream converts the rest
```

`-tags pdfonly|officeonly|webonly|imageonly` builds a smaller module — PDF and Illustrator; the Office/CAD/KiCad/music/font formats; HTML, Markdown and EPUB; or the images browsers decode themselves — so a page only fetches what it needs. `-tags previewonly` builds a converter-free module exposing `thumbnail()` and `text()`, the same Go code (`thumbnail`, `Document.SearchText`) a server runs, for the [thumbnail](../thumbnail/) and [search text](../text/) pages. `examples/common/convert-worker.ts` is the client-side glue that loads these modules on demand and exposes `convert`/`open`/`page`/`finish`/`thumbnail`/`text` over `postMessage`; `sniff()` in `examples/common/convert.ts` picks which module a dropped file needs from its content and extension.

A PDF converts a page at a time: `open()` returns pages sized but empty, then `page(i)` fills each one in as it converts — the pages in view first — and `finish()` returns the complete document, which replaces the streaming one. See [design.md §2](../design.html#2-go-と-wasm-について) for why this shape, and the [API reference](../api.html#ブラウザ内変換cmdbdfwasm) for the full call surface.

## Fonts in the browser

The Office/CAD/HTML/Markdown/EPUB converters need to measure and lay out text before they can convert a page; run in the browser, they read a font directory's `index.json` (name, size, and the byte ranges the font scan reads — table directories, `name`/`OS/2`/`post`) built by [`examples/common/build.mjs`](https://github.com/shibukawa/bdf/blob/main/examples/common/build.mjs), fetching those ranges ahead and a font's whole file only once a document actually uses it.

## Building the site

[`site/build.mjs`](https://github.com/shibukawa/bdf/blob/main/site/build.mjs) builds everything published on GitHub Pages: the top page, the full viewer, the thumbnail and search-text pages, the rendered documentation, and the shared `lib/` (workers + wasm modules), `fonts/` and `samples/` every page fetches from. `npm run site` builds it locally into `site/dist`; `npm run site:serve` also serves it. See [Building and testing](development.html) for the full command list, including the standalone viewer (`npm run demo`) that needs no Go build.
