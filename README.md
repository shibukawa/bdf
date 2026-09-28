# bdf

English | [日本語](README.ja.md)

**bdf** (Browser-specific Document Format) is a draft document format for previews that browsers can draw straight onto Canvas 2D.

Previewing a PDF, a Word file, or a CAD drawing in a browser usually means running a full office suite, or a headless browser, on the server. bdf's converters need neither: they're plain Go, and the same code also builds as WebAssembly, so it can convert inside the browser too. PDF, Office files (Word, PowerPoint, Excel, CSV, Parquet, Visio), draw.io diagrams, CAD drawings (DXF, Jw_cad, SXF, CGM, HP-GL/2), circuit boards (Gerber, Excellon, KiCad), Illustrator, Photoshop, HTML, Markdown, EPUB, music (MML, MIDI, MusicXML, playable in the viewer) and font files all convert this way, and one renderer, running in a Web Worker, draws every one of them.

**[Try it](https://shibukawa.github.io/bdf/)** — drop a file on the demo site; it converts and draws inside your browser, without being uploaded. **[Read the docs](https://shibukawa.github.io/bdf/docs/)** for why bdf exists, every format it reads, three sample server architectures, and how it's built.

## Quick start

```sh
go run ./cmd/bdf generate report.pptx report.bdf   # convert
go run ./cmd/bdf ls report.bdf                     # inspect

npm ci && npm run demo                             # open it in a browser
# http://127.0.0.1:8765/examples/viewer/.out/?src=/report.bdf
```

Go 1.27+ and Node.js are needed to build from source; see [Getting started](https://shibukawa.github.io/bdf/docs/getting-started.html) for the rest, including converting from Go code directly.

## Layout

| Path | Contents |
|---|---|
| `*.go`, `converter/`, `raster/`, `internal/` | The Go encoder/decoder and one converter package per input format |
| `packages/core`, `packages/render` | `@bdf/core` and `@bdf/render`: the TypeScript decoder and Canvas renderer |
| `examples/`, `site/` | The demo viewer, four sample server architectures, and the site published on GitHub Pages |
| `docs/` | This documentation |

See [Repository layout](https://shibukawa.github.io/bdf/docs/architecture/layout.html) for the full breakdown, and [Building and testing](https://shibukawa.github.io/bdf/docs/architecture/development.html) for `go test`, the golden tests, and regenerating fixtures.

## License

[MIT](LICENSE).
