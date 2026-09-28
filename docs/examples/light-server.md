# light-server: convert in the browser

[`examples/light-server`](https://github.com/shibukawa/bdf/tree/main/examples/light-server) — the server never converts a document to bdf. It serves the documents as they are and a thumbnail image of each, and the browser converts and renders.

```sh
cd examples/light-server
node web/build.mjs     # builds web/main.js, web/lib/ (workers + wasm converters), web/fonts/
go run .
open http://127.0.0.1:8081/
```

## What the server does

`main.go` walks a directory of documents (`-dir`, `../sample-files` by default) and, for each one that has no cached thumbnail yet, converts it in memory with `converter.ConvertFile` and draws a 320px thumbnail with `thumbnail.Make` — then **throws the converted bdf away**. Only the PNG is kept, under `-thumbs`. The server holds no bdf, ever; it doesn't even know if a file converts a second time the same way. It serves:

- `GET /` — an HTML page listing the documents that have a thumbnail, each linking to `/view/?file=NAME`
- `GET /files/NAME` — the original file, plain `http.FileServer`
- `GET /thumbs/NAME.png` — its cached thumbnail
- `GET /view/` — the built front end

## What the browser does

`web/main.ts` reads `?file=`, fetches the bytes from `/files/`, and converts them with the same wasm converter modules the demo site uses (`examples/common/convert.ts`'s `ConverterClient`, loaded from `web/lib/bdf-*.wasm`), then opens the result in [`MiniViewer`](https://github.com/shibukawa/bdf/tree/main/examples/miniviewer). The sequence itself is nothing new — the same `sniff → convert → open` that [`examples/viewer/main.ts`](https://github.com/shibukawa/bdf/blob/main/examples/viewer/main.ts) uses, pared down here to one file with no page-turning, cell selection, or streaming PDF pages.

`web/build.mjs` builds this front end by calling into [`examples/common/build.mjs`](https://github.com/shibukawa/bdf/blob/main/examples/common/build.mjs) — the same helper [the demo site](../architecture/browser.html) is built with — for the rendering worker, the converter worker, the wasm modules (every one but the preview-only module, which this server never needs), and the fonts the Office converters lay text out with.

## Why this shape

See [Why bdf](../why.html): the converters build as WebAssembly precisely so a server can be this thin. The cost lands on the browser instead — a wasm module to fetch (7–27 MB depending on the format, gzipped and cached after the first document of that kind) and a moment of local conversion on every open, not just the first. That's the right trade when there are many documents, the server has no spare capacity to convert them, or documents are opened rarely enough that pre-converting them all would be wasted work. For the opposite trade, see [preview-server](preview-server.html).
