# preview-server: convert on the server

[`examples/preview-server`](https://github.com/shibukawa/bdf/tree/main/examples/preview-server) — the server converts every document to bdf once and caches the result; the browser only ever fetches and draws an already-converted bdf, never a wasm converter.

```sh
cd examples/preview-server
node web/build.mjs     # builds web/main.js and web/lib/worker.js — no converters
go run .
open http://127.0.0.1:8082/
```

## What the server does

`main.go` walks a directory of documents (`-dir`, `../sample-files` by default) and, for each one with no cached bdf yet, converts it with `converter.ConvertFile`, writes the result with `(*bdf.Document).WriteSingle` under `-cache`, and draws its thumbnail the same way [light-server](light-server.md) does. Both are kept this time. It serves:

- `GET /` — an HTML page listing the documents, with their thumbnails and what the converter said of each (`res.Summary`), linking to `/view/?src=/cache/NAME.bdf`
- `GET /cache/NAME.bdf` — the converted document, through plain `http.FileServer`
- `GET /view/` — the built front end

`http.FileServer` answers Range requests on its own (`net/http`'s `ServeContent`), so the viewer fetches only the parts of the bdf it needs for the pages in view — the same as fetching a bundle from a CDN in [the server-side conversion path](../architecture/index.md#①-server-side-conversion).

## What the browser does

`web/main.ts` opens `?src=` directly with [`MiniViewer`](https://github.com/shibukawa/bdf/tree/main/examples/miniviewer) — no converter, no `sniff()`, because there's nothing left to convert. `web/build.mjs` only builds the rendering worker (`buildWorkers(lib, { convert: false })`); this front end is a few hundred KB, not tens of megabytes.

## Why this shape

This is the opposite tradeoff from [light-server](light-server.md). Conversion happens once, on the server, using the same converter packages any other bdf tool does (still no cgo, no external process — see [Why bdf](../why.md)) — not once per browser, per view. A document opens exactly as fast as its bdf can be range-fetched, on any device, without downloading a converter first. The cost is the server's own conversion capacity, and somewhere to keep the cache; for a system where documents are edited often and viewed by many different readers, that trade usually wins. See [search](search.md) for what indexing the same converted text for full-text search adds on top.
