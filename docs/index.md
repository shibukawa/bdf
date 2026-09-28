# Documentation

bdf is a document format made for previews in the browser. The [demo site](https://shibukawa.github.io/bdf/) is the fastest way to see that idea work; the pages below explain why it exists and how the pieces underneath it fit together.

## Read next

- **[Why bdf](why.md)** — what problem this solves, the formats at a glance, page turning, search and accessibility
- **[Getting started](getting-started.md)** — convert a file and open it, in five minutes
- **[Formats](formats/index.md)** — every input format, and how bdf lays each one out
- **[Sample architectures](examples/index.md)** — three ways to fit bdf into a system, each a small server you can run
- **[Architecture](architecture/index.md)** — how conversion and rendering work, formulas and thumbnails, passwords and protected mode, the repository layout, the `bdf` command, building and testing

## Reference (Japanese)

The detailed reference documents are written in Japanese, for now:

- **[Format specification](spec.md)** — the container, the manifest, the instruction set
- **[API reference](api.md)** — every package and command, by name
- **[Design notes](design.md)** — why each converter is built the way it is, section by section

## Source

The documentation lives in [`docs/`](https://github.com/shibukawa/bdf/tree/main/docs) as Markdown, next to the code it describes; this site renders it. Corrections and questions are welcome as [issues or pull requests](https://github.com/shibukawa/bdf).
