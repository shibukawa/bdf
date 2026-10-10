# BDF

English | [日本語](README.ja.md)

**BDF** (Browser-specific Document Format) is a document format for previews that browsers can draw straight onto Canvas 2D. Its stable wire format is version 1.0.

Previewing a PDF, a Word file, or a CAD drawing in a browser usually means running a full office suite, or a headless browser, on the server. BDF's converters need neither: they're plain Go, and the same code also builds as WebAssembly, so it can convert inside the browser too. PDF, Office files (Word, PowerPoint, Excel, CSV, Parquet, Visio), InDesign (IDML), draw.io diagrams, CAD drawings (DXF, Jw_cad, SXF, CGM, HP-GL/2), circuit boards (Gerber, Excellon, KiCad), Illustrator, Photoshop, HTML, Markdown, EPUB, music (MML, MIDI, MusicXML, playable in the viewer), audio files (MP3, M4A, FLAC, Ogg, WAV, AIFF: their cover art, tags and lyrics, playable in the viewer) and font files all convert this way, and one renderer, running in a Web Worker, draws every one of them.

**[Try it](https://shibukawa.github.io/bdf/)** — drop a file on the demo site; it converts and draws inside your browser, without being uploaded. **[Read the docs](https://shibukawa.github.io/bdf/docs/)** for why BDF exists, every format it reads, three sample server architectures, and how it's built.

## Quick start

```sh
go run ./cmd/bdf generate report.pptx report.bdf   # convert
go run ./cmd/bdf ls report.bdf                     # inspect

npm ci && npm run demo                             # open it in a browser
# http://127.0.0.1:8765/examples/viewer/.out/?src=/report.bdf
```

Go 1.27+ and Node.js are needed to build from source; see [Getting started](https://shibukawa.github.io/bdf/docs/getting-started.html) for the rest, including converting from Go code directly.

For an npm integration, including server-converted BDF and browser conversion presets, see [npm packages](docs/npm.md).

## Layout

| Path | Contents |
|---|---|
| `*.go`, `converter/`, `raster/`, `contrib/`, `internal/` | The Go encoder/decoder and one converter package per input format |
| `packages/` | The decoder, renderer, viewer surface and TinyGo-backed browser converter npm packages |
| `examples/`, `site/` | The demo viewer, four sample server architectures, and the site published on GitHub Pages |
| `docs/` | This documentation |

See [Repository layout](https://shibukawa.github.io/bdf/docs/architecture/layout.html) for the full breakdown, and [Building and testing](https://shibukawa.github.io/bdf/docs/architecture/development.html) for `go test`, the golden tests, and regenerating fixtures.

## License

[MIT](LICENSE).

### Fonts

The fonts in this repository have licenses of their own:

- **STIX Two Math** 2.13 b171 is embedded in the package `formula`, as the font formulas are laid out with. It is the file of the upstream release, unmodified ([stipub/stixfonts](https://github.com/stipub/stixfonts), tag `v2.13b171`; [`formula/fonts/README.md`](formula/fonts/README.md) has its SHA-256), under the [SIL Open Font License 1.1](formula/fonts/OFL.txt). Copyright 2001-2021 The STIX Fonts Project Authors, with Reserved Font Name "TM Math". A program that imports `formula` holds this font.
- **Bravura** (music symbols, as a table in `converter/internal/music/smufl`) is under the [SIL Open Font License 1.1](converter/internal/music/smufl/OFL.txt). Copyright © 2015 Steinberg Media Technologies GmbH, with Reserved Font Name "Bravura".
- **NewStroke** (KiCad's stroke font, in `converter/kicad`) is [CC0](converter/kicad/NEWSTROKE.txt).
- **DejaVu Sans** subsets (in `fixture`, for the sample document) are under the [DejaVu Fonts License](fixture/testdata/fonts/LICENSE.txt).
- The test fonts (subsets of M PLUS 1p, STIX Two Math and STIX Two Text) sit beside their licenses in `testdata` directories.

[Fonts and licenses](https://shibukawa.github.io/bdf/docs/licenses.html) lists them with the fonts the demo site serves.
