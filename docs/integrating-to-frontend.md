# Integrating into a frontend

Use `@bdfkit/viewer` to display BDF on a page. If a server already produces BDF, the viewer is enough. To convert input files in the browser, add a Wasm converter runtime. The viewer handles page placement, drawing, and scrolling; your application supplies the file picker, toolbar, and other controls.

These npm packages are currently workspaces in this repository. The `npm install` commands below show their use after publication. See [npm packages](npm.md) for the package profiles and measured sizes.

## Display BDF produced on a server

```sh
npm install @bdfkit/viewer
```

Give the host element a width and height:

```html
<div id="preview" style="width:100%;height:70vh"></div>
```

```ts
import { Viewer } from "@bdfkit/viewer";

const host = document.querySelector<HTMLElement>("#preview")!;
const viewer = new Viewer(host, { mode: "embedded", layout: "single" });
await viewer.open({ kind: "single", url: "/documents/report.bdf" });
// For split BDF: await viewer.open({ kind: "split", base: "/documents/report/" });

viewer.setLayout("spread");      // "single" / "spread" / "continuous"
viewer.setZoom(1.25);
viewer.scrollToPage(2);           // zero-based
// Call viewer.destroy() when removing the view.
```

A single-file BDF can be fetched in parts through HTTP Range requests. See [Integrating into a server](integrate-to-server.md) for producing it. If you want to place pages yourself, use the lower-level `@bdfkit/render` and `@bdfkit/core` packages.

## Convert files in the browser

Choose `@bdfkit/viewer` and one converter package matching the input formats you need.

| Input formats | Package | Factory |
|---|---|---|
| PDF, EPUB | `@bdfkit/convert-pdf-epub` | `createPdfEpubConverter()` |
| Word, Excel, PowerPoint, Visio, CSV/TSV, Parquet, PDF, EPUB | `@bdfkit/convert-office` | `createOfficeConverter()` |
| All formats | `@bdfkit/convert-all` | `createAllConverter()` |

```sh
npm install @bdfkit/viewer @bdfkit/convert-office
```

```ts
import { Viewer } from "@bdfkit/viewer";
import { createOfficeConverter } from "@bdfkit/convert-office";

const viewer = new Viewer(document.querySelector<HTMLElement>("#preview")!, {
  mode: "embedded", layout: "single",
});
const converter = createOfficeConverter();
await converter.ready; // also returns the formats included in the module

const input = document.querySelector<HTMLInputElement>("#file")!;
input.addEventListener("change", async () => {
  const file = input.files?.[0];
  if (!file) return;
  try {
    await viewer.openFile(converter, await file.arrayBuffer(), {
      name: file.name,
      fonts: "/fonts/", // when the format needs fonts for text layout
    });
  } catch (error) {
    console.error(error); // replace with your application's error UI
  }
});

// Call viewer.destroy() and converter.terminate() when removing the view.
```

If you pass `fonts`, serve font files and an `index.json` at that URL. For formats that support page-by-page conversion, `openFile` shows page outlines first and fills pages as conversion progresses. Listen for `conversionprogress`, `conversiondone`, and `error` events to update your UI.

## Find in the document

The viewer draws pages as bitmaps, so the browser's own find has nothing to search in them. While the focus is in the viewer, Cmd+F (macOS) or Ctrl+F opens the viewer's find bar instead. It searches the text of the whole view, including pages that are not drawn yet, marks the matches, and scrolls to the one shown. Enter or Cmd/Ctrl+G goes to the next match, Shift+Enter or Shift+Cmd/Ctrl+G to the previous one, and Escape closes the bar. Pressing the shortcut again inside the bar's field leaves it to the browser.

```ts
new Viewer(host, { find: "page" }); // the shortcut is the viewer's wherever the focus is: a page that is only the viewer
new Viewer(host, { find: false });  // no find bar; the shortcut stays the browser's
```

The default is `"focus"` for an embedded viewer and `"page"` for a lightbox while it is shown. `findLabels` replaces the bar's English words. To search from your own controls, call `viewer.find(query, step)`: a new query shows the first match from the page being read, and the same query again steps to the next (`1`) or previous (`-1`) match. It resolves to `{ query, index, total, more }`, which the `findchange` event also carries. `viewer.clearFind()` ends the search and `viewer.openFind()` shows the bar.

## Presentation and controls

For a lightbox, create `new Viewer(host, { mode: "lightbox" })` and call `viewer.show()` / `viewer.hide()` from your controls. Use `setView(id)` to choose a BDF view and `setLayout(...)` to switch single, spread, or continuous pages. Read the available views from `viewer.document?.views` and the current one from `viewer.view`.

Music presentation is outside the default viewer import. Import `@bdfkit/viewer/music` only on screens that need piano roll or `MusicPlayer`. Without the `music` option, the viewer can still draw score pages, but `setMusicMode("piano-roll")` is unavailable.

```ts
import { Viewer } from "@bdfkit/viewer";
import { pianoRoll, MusicPlayer } from "@bdfkit/viewer/music";

const viewer = new Viewer(host, { music: pianoRoll });
viewer.setMusicMode("piano-roll"); // switch back with "score"
// Use MusicPlayer when building playback controls.
```

If your application only converts files, omit `@bdfkit/viewer` and use `@bdfkit/convert` with a converter preset. `converter.convert(...)` returns BDF bytes that you can upload or download. [npm packages](npm.md#payload-sizes) also measures the profiles without a renderer.

For your own combination of formats, follow the TinyGo steps in [Build a Wasm runtime](build-wasm-runtime.md), then create `createConverter(new URL("/my-converter.wasm", location.href))` from `@bdfkit/convert`. Pass that converter to the same `openFile` method. Configure your bundler to serve the package's Worker and Wasm assets, and allow Workers and WebAssembly in your content security policy.
