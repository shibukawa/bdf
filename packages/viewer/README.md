# @bdfkit/viewer

An unopinionated document surface for BDF. `Viewer` owns page placement, scrolling, lazy rendering, and presentation state while the host application owns its toolbar, file picker, tabs, and lightbox controls.

```sh
npm install @bdfkit/viewer
```

```ts
import { Viewer } from "@bdfkit/viewer";

const viewer = new Viewer(document.querySelector("#document")!, {
  mode: "embedded",       // or "lightbox"
  layout: "spread",        // single, spread, or continuous
});
await viewer.open({ kind: "url", url: "/reports/report.bdf" });
```

The package is renderer-only and expects the server to provide BDF. To convert files in the browser, pass a converter from `@bdfkit/convert-pdf-epub`, `@bdfkit/convert-office`, or `@bdfkit/convert-all` to `viewer.openFile`.

While the focus is in the viewer, the browser's find shortcut (Cmd+F on macOS, Ctrl+F elsewhere) opens the viewer's own find bar, which searches the text of the whole view, including pages that are not drawn yet. Pass `find: "page"` when the page is only the viewer, or `find: false` to leave the shortcut to the browser and call `viewer.find(query)` from your own controls.

Music score rendering is included in the normal viewer API. Import `@bdfkit/viewer/music` only when piano-roll rendering or playback is needed, so applications that do not use music can omit that code from their bundle.

See the [frontend integration guide](https://github.com/shibukawa/bdf/blob/main/docs/integrating-to-frontend.md) for CSP, bundler, layout, and conversion examples.

## License

MIT. See [LICENSE](LICENSE).
