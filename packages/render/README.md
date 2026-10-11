# @bdfkit/render

Canvas and Web Worker rendering APIs for BDF documents. This package provides the drawing layer, text selection helpers, document search, table selection, SVG support, and optional music playback primitives.

```sh
npm install @bdfkit/render
```

The default worker is bundled with the package:

```ts
import { createRenderWorker } from "@bdfkit/render";

const renderer = createRenderWorker();
const manifest = await renderer.open({ kind: "buffer", buffer: bdfBytes });
```

`renderer.svg(view, page)` writes a page as an SVG file that a browser shows without a viewer (text as text, with its fonts and images in the file), and `renderer.image(view, page, scale)` as a WebP, PNG or JPEG image; both also take a rectangle of a sheet or of a view's continuous layout.

`@bdfkit/render` depends on `@bdfkit/core`. It supplies rendering and page operations; it does not provide a document surface or toolbar. Use [`@bdfkit/viewer`](https://github.com/shibukawa/bdf/tree/main/packages/viewer) for scrolling, embedded/lightbox presentation, and page layout.

The bundler must preserve the package's `worker.js` asset and allow module Workers. See the [npm integration guide](https://github.com/shibukawa/bdf/blob/main/docs/npm.md).

## License

MIT. See [LICENSE](LICENSE).
