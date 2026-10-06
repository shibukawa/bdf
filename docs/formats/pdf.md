# PDF, Illustrator

PDF pages are a sequence of drawing instructions, not a document model, so BDF reads them the same way a PDF viewer does: replaying the content stream into BDF's own instruction set instead of trying to recover paragraphs or tables. Illustrator's `.ai` format builds directly on that. From Illustrator 9 onward, an `.ai` file is a PDF with Illustrator's own data attached beside it — so BDF reads it with the same PDF converter and treats each artboard as a page.

## Try it

Drop a file on the [viewer](https://shibukawa.github.io/bdf/viewer/), or try a sample: a [PDF](https://shibukawa.github.io/bdf/viewer/?file=samples/demo.pdf), a [three-page PDF made with reportlab](https://shibukawa.github.io/bdf/viewer/?file=samples/reportlab-master.pdf), or an [Illustrator file with three artboards](https://shibukawa.github.io/bdf/viewer/?file=samples/artboards.ai).

## PDF (`converter/pdf`)

`converter/pdf` rebuilds embedded fonts (TrueType, CFF, OpenType, Type1) into WOFF2 files holding only the glyphs actually used, checking the OS/2 embedding permission (`fsType`) first and carrying the original copyright notice over. Repeated form XObjects become shared objects, text becomes searchable runs, and the content every page starts with — a repeated header, logo or footer — is split out into one shared object instead of being stored once per page. Soft masks (fades, drop shadows, the masks behind Chrome's CSS `mask-image`) are drawn by the viewer rather than baked into a bitmap. JPEG 2000 and JBIG2 images, which browsers can't decode themselves, go through BDF's own pure-Go decoders instead. Text set in a CJK font that isn't embedded is read through Adobe's predefined CMaps (Shift_JIS, EUC, UCS-2, and the others). Vertical writing (`WMode 1`) is laid out as vertical lines. Optional content (layers) is fixed at whatever a viewer would show on first opening the document — hidden layers aren't drawn — and a tagged PDF's structure tree becomes BDF's accessible-text MARKs (headings, paragraphs, lists, tables with cell spans, figures with alt text, language).

Each PDF page becomes one BDF page by default (a `fixed` view); `-kind flow` makes a flow view instead, the same kind of view Word produces, so the viewer can also offer it as a continuous scroll.

| Option | Values | Default |
|---|---|---|
| `-kind` / `-param kind=` | `fixed`, `flow` | `fixed` |
| `-param box=` | `crop`, `media`, `bleed`, `trim`, `art` | `crop` |
| `-no-share` / `-param no-share=` | `true`: don't move the prefix pages have in common into a shared object | shared |

See [design.md §3.1](../design.md#31-pdf--bdf-変換器converterpdfの構造) for the internals.

## Illustrator .ai (`converter/ai`)

An Illustrator file saved with "Create PDF Compatible File" — on by default since Illustrator 9 — is a PDF with Illustrator's own data attached alongside it, in a private stream the PDF part never refers to. `converter/ai` draws that PDF part with the PDF converter above and never reads Illustrator's private data at all. Each PDF page is one artboard, and BDF crops it to the trim box — the artboard itself, without the bleed — rather than the larger media or bleed box a page also carries. Hidden layers, which Illustrator stores the same way any PDF stores optional content, stay hidden for the same reason they do in a plain PDF.

Files saved with "Create PDF Compatible File" turned off keep only a one-page PDF that says so, and the PostScript-based `.ai` format Illustrator 8 and earlier wrote has no PDF at all. Both are reported as errors rather than silently converted, since drawing either would mean reading Illustrator's own native format, which BDF doesn't do.

| Option | Values | Default |
|---|---|---|
| `-param box=` | `trim`, `bleed`, `media`, `crop`, `art` | `trim` (the artboard, without the bleed) |

See [design.md §3.17](../design.md#317-illustrator--bdf-変換器converterai) for the internals.
