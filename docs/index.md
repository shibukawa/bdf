# Documentation

BDF is a document format made for previews in the browser. The [demo site](https://shibukawa.github.io/bdf/) is the fastest way to see that idea work; the pages below explain why it exists and how the pieces underneath it fit together.

## Read next

- **[Why BDF](why.md)** — what problem this solves, the formats at a glance, page turning, search and accessibility
- **[Getting started](getting-started.md)** — convert a file and open it, in five minutes
- **[Integrating into a frontend](integrating-to-frontend.md)** — show server-produced BDF or convert files in the browser
- **[Integrating into a server](integrate-to-server.md)** — register only the Go converters your service needs
- **[Rendering](rendering.md)** — draw documents and formulas into a Go `image.Image`, with Ebitengine, or on a browser canvas
- **[Build a Wasm runtime](build-wasm-runtime.md)** — choose converter formats and build for TinyGo or Go
- **[Formats](formats/index.md)** — every input format, and how BDF lays each one out
- **[Sample architectures](examples/index.md)** — three ways to fit BDF into a system, each a small server you can run
- **[Architecture](architecture/index.md)** — how conversion and rendering work, formulas and thumbnails, passwords and protected mode, the repository layout, the `bdf` command, building and testing

## Reference (Japanese)

The detailed reference documents are written in Japanese, for now:

- **[Format specification v1.0](spec.md)** — the stable container, manifest, instruction set, and compatibility policy
- **[API reference](api.md)** — every package and command, by name
- **[Design notes](design.md)** — why each converter is built the way it is, section by section
- **[Conversion quality report](quality.md)** — Japanese results for fixtures, rendering goldens, browser conversion, and known limits

- **[Fonts and licenses](licenses.md)** — the fonts BDF holds and the demo site serves, and where each comes from (in English and Japanese)

## Source

The documentation lives in [`docs/`](https://github.com/shibukawa/bdf/tree/main/docs) as Markdown, next to the code it describes; this site renders it. Corrections and questions are welcome as [issues or pull requests](https://github.com/shibukawa/bdf).
