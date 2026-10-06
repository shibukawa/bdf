# Getting started

The fastest way to see BDF work is the [demo site](https://shibukawa.github.io/bdf/viewer/): drop a file on the page and it converts and draws inside your browser, without uploading anything. From here, you build and run BDF on your own machine instead.

## The CLI

BDF needs Go 1.27 or later.

```sh
git clone https://github.com/shibukawa/bdf
cd bdf
go run ./cmd/bdf generate report.pptx report.bdf
```

`generate` picks the input format from the file's content (Markdown, which can't be told from its bytes, is picked from the `.md` extension). The output can also be a directory ending in `/`, which writes the [split form](spec.md#34-分割形式split) instead of one file — the shape a file server or CDN serves without repacking.

```sh
go run ./cmd/bdf generate report.pptx out/          # split form: manifest.json and one file per part
go run ./cmd/bdf generate -pages 1-3 book.pdf out.bdf   # only the first three pages
go run ./cmd/bdf generate -thumbnail thumb.png -text text.json report.docx out.bdf  # a thumbnail and search text too
```

`go run ./cmd/bdf generate -h` lists every flag and the `-param` options specific to each input format.

## Look at what you made

```sh
go run ./cmd/bdf ls out.bdf          # the parts: fonts, images, object streams, their sizes
go run ./cmd/bdf manifest out.bdf    # the manifest as JSON: views, pages, metadata
```

## Open it in a browser

Clone the repository, then:

```sh
npm ci
npm run demo
# http://127.0.0.1:8765/examples/viewer/.out/?src=/out.bdf
```

`examples/viewer` is the same viewer the demo site uses (search, text selection, the pages that turn like a book's, a score that plays); `?src=` points it at any `.bdf` under the repository. `npm run site:serve` runs the fuller demo site locally, with the in-browser converters — the same experience as [the one published on GitHub Pages](https://shibukawa.github.io/bdf/), rebuilt from your checkout.

## From Go code

```go
import (
	"github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/pptx" // or converter/all for every format
)

res, err := converter.ConvertFile("report.pptx", "", &converter.Options{}) // "" detects the format
if err != nil {
	// converter.ErrPasswordRequired / ErrWrongPassword if the input needs a password
}
// res.Doc is a *bdf.Document; write it with res.Doc.WriteSingle(w) or WriteSplit(dir)
```

Each converter package registers its format as a side effect of being imported, so a program only handles the formats it links. See the [API reference](api.md) for the full surface (Go, the CLI, the browser wasm modules, and the two TypeScript packages), and [Sample architectures](examples/index.md) for three small servers built around this.

## Next

- **[Why BDF](why.md)** — the formats at a glance, and what stays usable after conversion: page turning, search, selection, accessibility
- **[Formats](formats/index.md)** — every input format and its layout
- **[Architecture](architecture/index.md)** — how conversion and rendering fit together, and how to build and test the repository
