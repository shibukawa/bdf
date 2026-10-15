# The bdf command

```
bdf generate [flags] <input> <out.bdf | dir/>
                                   convert a document or a drawing (bdf generate -h for flags and formats)
bdf fill [flags] <template.xlsx | .docx | .pptx | .pdf | .json | .bdft> <out.bdf | dir/>
                                   fill a template of a form with data: an Office document with tags, a
                                   PDF with form fields, a template of pdfme (bdf fill -h for flags)
bdf template compile [flags] <template> <out.bdft>
bdf template schema [-json-schema] <template> [out.json]
                                   check a template and write it compiled; write what it asks of data
bdf pdf [flags] <file.bdf | dir> <out.pdf>
                                   write the pages as a PDF file (bdf pdf -h for flags)
bdf thumbnail [flags] <file.bdf | dir | input> <out.png | .jpg | .webp>
                                   draw a thumbnail: the top of the first page or the whole slide
                                   (of an input to convert, only the first page is converted)
bdf text [flags] <file.bdf | dir> [out.json]
                                   write the metadata and the text of each page for a search index
bdf render [flags] <file.bdf | dir> <out.png | .jpg | .webp>
                                   draw a page (bdf thumbnail -h, text -h, render -h for flags)
bdf ls <file.bdf | dir>            list views and parts
bdf manifest <file.bdf | dir>      print the manifest as JSON
bdf disasm <file.bdf | dir> <hash> disassemble an object part
bdf extract <file.bdf | dir> <hash> <out>
bdf split <file.bdf> <dir>         write the split form
bdf join <dir> <file.bdf>          write the single-file form
bdf encrypt [-password-file f] <file.bdf | dir> <out.bdf | dir/>
                                   encrypt a document with a password
bdf decrypt [-password-file f] <file.bdf | dir> <out.bdf | dir/>
                                   write an encrypted document without its encryption
bdf demo <file.bdf | dir/>         write the fixture document (dir/ ends with a slash)
bdf conformance [-testdata dir] <dir>
                                   write the conformance suite for other readers (package conformance)
```

An encrypted document is read with the password in `$BDF_PASSWORD`; `split` and `join` copy it without needing the password.

## generate

Converts one input into BDF. The format is detected from the input's content, or set explicitly with `-format`; the output is a single file, or a directory ending in `/` for the split form. `-pages 1-3,5,8-` selects pages, slides, sheets or artboards; `-thumbnail` and `-text` write a preview next to the output (not for an encrypted output, unless `-allow-plaintext`); `-dc name=value` overrides a Dublin Core element (`-dc name=` removes it). See [Formats](../formats/index.md) for the `-param` options specific to each input format, or `bdf generate -h` for the current, complete list.

```sh
bdf generate report.pptx report.bdf              # single file
bdf generate report.pptx out/                    # split form
bdf generate -pages 1-3 book.pdf book.bdf         # first three pages only
bdf generate -password-file pw.txt secret.xlsx out.bdf  # a protected input
```

## fill, template

Fill a template of a form (an Excel workbook, a Word document or a PowerPoint presentation whose text holds tags) with data. The format is told from the content. See [Making forms](../forms.md) for how to write one.

```sh
bdf fill -data order.yaml -data items=lines.csv invoice.xlsx invoice.bdf   # fill, as a document of the printed pages
bdf fill -data order.yaml notice.docx notice.bdf                           # a Word document as its pages, a presentation as its slides
bdf fill -each orders.csv invoice.xlsx 'out/{{no}}.bdf'                    # a form for each record
bdf template compile -sample sample.yaml invoice.xlsx invoice.bdft         # check it and write a compiled template
bdf template schema invoice.bdft                                           # the names its tags look up, as JSON
bdf template schema -json-schema invoice.xlsx                              # the JSON Schema of the data, to make a form of
bdf fill -data applicant.json application.pdf application.bdf              # a PDF with form fields: the values in the places of its fields
bdf fill -data invoice.json pdfme-template.json invoice.bdf                # a template of pdfme, as its generator fills it
```

`-data` is a JSON or YAML object (`file`) or a value under a name (`name=file`; a CSV or TSV file or the sheets of a workbook are lists of objects), and repeats. `-strict` makes a tag whose name the data lacks an error, and `-files` names the directory of the pictures the data names. The font and image flags (`-font-dir`, `-no-system-fonts`, `-fonts`, `-images` …), `-param` and `-pages` mean what they mean for `generate` (`-pages` selects the sheets of a workbook, the slides of a presentation's template, and the pages of the form of a Word template). A `.bdft` is a BDF document: `ls`, `thumbnail` and `render` open it as it is.

## pdf

Write the pages of a document as a PDF file: its text as text in the embedded fonts, its drawings as vectors ([Writing a PDF](../forms.md#writing-a-pdf)).

```sh
bdf pdf invoice.bdf invoice.pdf
bdf pdf -view print1 book.bdf sheet1.pdf            # choose views (repeatable)
bdf pdf -skip background -skip master slip.bdf slip.pdf   # leave layers out by their roles
bdf pdf -tagged report.bdf report.pdf               # a tagged PDF: with the structure of the document, for screen readers
bdf pdf -pdfa 2a invoice.bdf invoice.pdf            # a PDF/A file (2b, 2u, 2a, 3b, 3u, 3a), for keeping
```

By default all the views of pages (`fixed`, `flow`) are written, and a document without pages becomes one page. `-font-dir` and `-no-system-fonts` say where the fonts the document refers to by name are looked for. `-ignore-fstype` embeds fonts of the system that do not allow it (only with the rights to; by default their text is drawn as outlines). `-no-compression` writes the streams as they are, to read them. An encrypted document needs `-allow-plaintext` (the PDF is not encrypted).

## thumbnail, text, render

Three ways to look at a `.bdf` without a browser, using the same Go code (`imagebdf`, `thumbnail`, `Document.SearchText`) a server would.

```sh
bdf thumbnail -size 256 out.bdf thumb.png    # -mode auto|crop|fit, -view, -sheet-dpi, -font-dir, -no-system-fonts
bdf thumbnail -size 256 book.pdf thumb.png   # an input bdf generate converts: only its first page is converted
bdf text out.bdf text.json                   # metadata and per-page text as JSON (default: stdout)
bdf render -page 2 -scale 2 out.bdf page2.png  # a sheet takes -width/-height too (from A1)
bdf render -antialias=off -scale 8.333333 out.bdf print.png  # 600 dpi with faster, unsmoothed edges
bdf render -pages=all -workers=4 -antialias=off out.bdf 'page-{page}.png'  # parallel independent pages, ordered output
bdf render -pages=all -antialias=off -mono=bilevel -scale 8.333333 out.bdf 'page-{page}.png'  # 600 dpi black and white, 1 bit a pixel
```

`render -antialias=default|on|off` controls the edges of shapes, text and geometric clips. `default` and `on` preserve the existing smoothing at every scale; `off` skips edge coverage calculations. There is no dpi-based automatic switch. Transparency, blending, image interpolation, soft masks, blur and the reduced alpha of subpixel-width lines are preserved. `thumbnail` and generated previews keep their existing smoothing. See [Rendering](../rendering.md#raster-quality-and-speed) for the Go option and output tradeoffs. `render -mono=gray|bilevel` writes a monochrome image instead of a colour one: 8 bits of luminance a pixel, or black and white (1 bit a pixel in a PNG; black where the gray value is below 128), for a page or a batch ([Monochrome output](../rendering.md#monochrome-output)).

Given an input that is not a `.bdf` (a PDF, a `.docx` …), `thumbnail` converts only its first page, all a thumbnail shows (the whole input with `-view`): a 500-page PDF takes milliseconds instead of seconds. A password-protected input needs `$BDF_PASSWORD` and `-allow-plaintext`, as an encrypted document does.

`-allow-plaintext` is required for `thumbnail` and `text` on an encrypted document (`$BDF_PASSWORD` decrypts it first) — neither output is itself encrypted, so this is opt-in. See [Conversion and output → Server-side thumbnails and search text](conversion.md#server-side-thumbnails-and-search-text) for what each one draws or extracts.

## Inspecting a document

```sh
bdf ls out.bdf              # every part: hash, kind, encoding, compressed/uncompressed size
bdf manifest out.bdf        # the manifest as JSON: views, pages, metadata
bdf disasm out.bdf <hash>   # disassemble one object part's drawing instructions
bdf extract out.bdf <hash> out.bin  # a part's raw bytes (a font, an image, …)
```

## Converting between the two container forms

```sh
bdf split out.bdf out/    # single file → split form
bdf join out/ out.bdf     # split form → single file
```

Neither needs a password: they copy the container's parts as they are, encrypted or not.

## encrypt, decrypt

```sh
BDF_PASSWORD=… bdf encrypt in.bdf out.bdf   # seal an existing bdf with a password
BDF_PASSWORD=… bdf decrypt in.bdf out.bdf   # remove the encryption (needs the password to open it)
```

`generate` can do this in one step too, with `-encrypt auto` (the default: only when the input needed a password), `always`, or `never`.

## demo

```sh
bdf demo out.bdf      # the fixture document used by the golden tests and the demo viewer's default
bdf demo out/         # its split form
```

## conformance

```bash
bdf conformance out/                 # the suite: suite.json and inputs/, with testdata/ as accept cases
bdf conformance -testdata "" out/    # the synthetic cases alone
```

The suite is for readers written in other languages (libbdf): inputs and what this implementation does with them, as `suite.json` records it (see package `conformance`). CI uploads it as the `bdf-conformance` artifact, and the release workflow attaches `bdf-conformance-<tag>.tar.gz` to the release of a tag.

## Font embedding flags, across formats

A few flags apply to every format that lays out and embeds text (PowerPoint, Excel, Word, CSV, Parquet, metafiles, HTML, Markdown, EPUB, font files):

| Flag | Effect |
|---|---|
| `-font-dir DIR` | Search this directory for fonts before the system ones (repeatable) |
| `-fonts embed \| system` | Embed the fonts laid out with (subset), or refer to them by name (the viewer's own fonts draw the text — the default for HTML, Markdown, EPUB) |
| `-no-system-fonts` | Only the fonts under `-font-dir`, none from the system |
| `-no-woff2` | Store embedded fonts as plain TrueType/OpenType instead of WOFF2 |
| `-no-subset` | Embed whole font files instead of just the glyphs in use |
| `-ignore-fstype` | Embed a font whose OS/2 `fsType` forbids embedding or subsetting — only with the rights to do so |
| `-compression normal \| best \| fast` | How hard the parts are deflated: `normal` (level 6, the default), `best` (level 9: a few percent smaller, four times the time) or `fast` (level 1) |

See the [API reference](../api.md) for the equivalent Go calls, and [Building and testing](development.md) for how `bdf generate` and `bdf demo` are used to regenerate the fixtures under `testdata/`.
