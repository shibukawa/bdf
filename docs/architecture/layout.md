# Repository layout

## Core

| Path | Contents |
|---|---|
| `*.go`, `cmd/bdf` | The Go encoder, decoder, container I/O and CLI |
| `converter/` | Registry of the input formats, shared options, format detection, page ranges |
| `converter/internal/` | Shared by the Office and CAD converters — see below |
| `raster/imagebdf/` | Draws pages and single objects into images in pure Go, as the viewer does: the software rasterizer, and an SVG renderer for SVG images (read by `image/svg`) |
| `raster/ebitenginebdf/` | Draws objects of paths (formulas, say) with Ebitengine. A module of its own, with its own `go.mod`, since it depends on Ebitengine; `raster/internal/shapes` is the shared part that turns an object into a list of filled and stroked shapes |
| `image/formula/` | The public package that lays out LaTeX and MathML formulas as objects of paths; embeds the formula font (STIX Two Math) |
| `thumbnail/` | Thumbnails of documents: the layout by kind of document, and PNG/JPEG/WebP encoding |
| `internal/` | Font lookup, measurement and subsetting (`fontdb`); TrueType/OpenType reading, writing and glyph outlines (`sfnt`); CFF reading and subsetting (`cff`); the OpenType layout tables GSUB, GPOS and GDEF (`otlayout`); the formula tree, its readers (LaTeX, MathML, Office Math) and the layout (`mathlayout`), and the XML trees of Office documents (`xmltree`) — shared by the converters and `imagebdf` |
| `contrib/` | Pieces that lean on the platform: `otf` opens font files mapped into memory on Unix and Windows (read whole elsewhere), so that only the tables and glyphs in use become resident; `fontdb` opens the fonts it loads with it |
| `font/woff2/` | TrueType/OpenType ↔ WOFF2 (the glyf transform, and Brotli) |
| `image/svg/` | The public reader of SVG: the tree of elements with the style sheets applied, and the values of attributes (colours, lengths, transforms). `raster/imagebdf` draws what it reads |
| `image/imgconv/` | How images are stored (as they are, or converted to WebP); bundles a pure-Go libwebp |
| `fixture/` | Generates the sample document used by tests and `bdf demo`, with embedded fonts |

## Converters, one package per input format

| Path | Converts |
|---|---|
| `converter/pdf` | PDF |
| `converter/ai` | Illustrator (.ai), over the PDF converter |
| `converter/psd` | Photoshop (.psd, .psb) |
| `converter/pptx` | PowerPoint (.pptx) |
| `converter/xlsx` | Excel (.xlsx) |
| `converter/csv` | CSV and TSV (drawn by `converter/xlsx`) |
| `converter/parquet` | Apache Parquet (drawn by `converter/xlsx`) |
| `converter/docx` | Word (.docx) |
| `converter/visio` | Visio (.vsdx) |
| `converter/idml` | InDesign (.idml) |
| `converter/drawio` | draw.io (.drawio / .drawio.svg / .drawio.png) |
| `converter/dxf` | AutoCAD DXF |
| `converter/jww` | Jw_cad (.jww) |
| `converter/sxf` | SXF (.p21, .p2z, .sfc) |
| `converter/cgm` | CGM (.cgm, .cgz) |
| `converter/hpgl` | HP-GL/2 plot files (.plt; HP-GL, PJL, PCL, HP RTL) |
| `converter/gerber` | Gerber and Excellon (and zip archives of a board's files) |
| `converter/kicad` | KiCad schematics, boards and projects (and zip archives of a project) |
| `converter/emf` | Windows metafiles (.emf, .wmf) |
| `converter/tiff` | TIFF (.tif, .tiff) |
| `converter/html` | HTML (.html, .xhtml, .mhtml), in reader mode |
| `converter/markdown` | Markdown, through HTML |
| `converter/epub` | EPUB (reflowable books in reader mode, fixed-layout books of pictures) |
| `converter/mml` | MML (.mml; generic, Mabinogi, PPMCK dialects), as a score |
| `converter/midi` | Standard MIDI Files (.mid, .kar, .rmi), as a score |
| `converter/musicxml` | MusicXML (.musicxml, .mxl) |
| `converter/font` | Font files (.ttf, .otf, .ttc, .woff, .woff2): a preview of characters, glyphs and OpenType features |
| `converter/image` | PNG, JPEG, GIF, WebP, AVIF, BMP, ICO, SVG (stored as they are, metadata read) |
| `converter/audio` | MP3, M4A, FLAC, Ogg, WAV, AIFF (a card of the cover art and the tags; the audio is not stored) |
| `converter/all` | Registers every input format (imported for its side effect) |

`converter/internal/` holds what several of these share: OOXML packages and XML (`ooxml`), DrawingML shapes/text/tables/charts (`ooxml/drawingml`), font choice/measuring/embedding for text layout (`fontset`, used by draw.io too), objects under construction (`canvas`, used by draw.io too), EMF/WMF replay (`metafile`), the boxes of the ISO base media file format (`isobmff`: AVIF images, M4A audio), CAD drawings plotted onto pages (`cad`), line-breaking rules (`linebreak`), compound files (`cfb`) and the decryption of password-protected Office documents (`offcrypto`); for PDF, Adobe's predefined CJK CMaps (`cjkcmap`) and the JPEG 2000/JBIG2 decoders (`jpx`, `jbig2`); the TIFF reader with its CCITT fax decoder (`tiff`); the layout engine shared by Word, HTML, Markdown and EPUB (`wordproc`: paragraphs, tables, pages and scroll views) and the HTML/XHTML reader they share (`webdoc`); Dublin Core from XMP metadata (`xmp`); the part that runs the formula engine (`internal/mathlayout`) with the fonts of a document and draws formulas as text in the embedded font (`equation`); and the engraving and playback of MML/MIDI/MusicXML scores (`music`, with Bravura's SMuFL glyphs in `music/smufl`).

## TypeScript packages

| Path | Contents |
|---|---|
| `packages/core` | `@bdfkit/core` — decoder, container loading, text extraction (no dependencies) |
| `packages/render` | `@bdfkit/render` — Canvas renderer, page/continuous/sheet rendering (a scroll view renders as continuous), the Worker and its client (SVG images are drawn on the main thread, which workers can't decode) |
| `cmd/bdfwasm` | The converters built as WebAssembly for in-browser conversion (modules for PDF, the Office formats, HTML/Markdown/EPUB, images, and a converter-free module for thumbnails and search text) |

## The demo site and its pages

| Path | Contents |
|---|---|
| `examples/viewer` | The full demo viewer (Worker-based rendering, the accessible text layer, page turning) |
| `examples/common` | What the viewer and the rest of the site share: the converter Worker's client, and the site's esbuild helpers |
| `examples/miniviewer` | A small, embeddable viewer (used by the sample architectures below) |
| `site/` | The site published on GitHub Pages: the top page, the thumbnail and search-text pages, and the documentation renderer (`site/docs.mjs`, which builds this page from `docs/*.md`) |
| `docs/` | This documentation, as Markdown |

## Sample projects, the CLI and tests

| Path | Contents |
|---|---|
| `examples/light-server`, `examples/preview-server`, `examples/search`, `examples/secure-reader` | Four small Go servers, each built around a different way of fitting BDF into a system — see [Sample architectures](../examples/index.md) |
| `cmd/bdf` | The `bdf` CLI: `generate`, `thumbnail`, `text`, `render`, `ls`, `manifest`, `disasm`, `extract`, `split`, `join`, `encrypt`, `decrypt`, `demo` — see [The BDF command](cli.md) |
| `tools/` | One-off generators run by hand when their source changes: Unicode data, SMuFL glyph paths, draw.io stencils and the AWS icon map, table styles, cmaps, and the pure-Go WebP codec's generated tables |
| `testdata/` | Generated samples and golden images, checked byte-for-byte (or pixel-for-pixel) in CI |
| `test/` | The Playwright-driven golden tests and the fixture generators for each format's `testdata/` |

See the [API reference](../api.md) for what each package exports, and [Building and testing](development.md) for how these pieces are built, tested and regenerated.
