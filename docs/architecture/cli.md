# The bdf command

```
bdf generate [flags] <input> <out.bdf | dir/>
                                   convert a document or a drawing (bdf generate -h for flags and formats)
bdf thumbnail [flags] <file.bdf | dir> <out.png | .jpg | .webp>
                                   draw a thumbnail: the top of the first page or the whole slide
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
```

An encrypted document is read with the password in `$BDF_PASSWORD`; `split` and `join` copy it without needing the password.

## generate

Converts one input into bdf. The format is detected from the input's content, or set explicitly with `-format`; the output is a single file, or a directory ending in `/` for the split form. `-pages 1-3,5,8-` selects pages, slides, sheets or artboards; `-thumbnail` and `-text` write a preview next to the output (not for an encrypted output, unless `-allow-plaintext`); `-dc name=value` overrides a Dublin Core element (`-dc name=` removes it). See [Formats](../formats/index.md) for the `-param` options specific to each input format, or `bdf generate -h` for the current, complete list.

```sh
bdf generate report.pptx report.bdf              # single file
bdf generate report.pptx out/                    # split form
bdf generate -pages 1-3 book.pdf book.bdf         # first three pages only
bdf generate -password-file pw.txt secret.xlsx out.bdf  # a protected input
```

## thumbnail, text, render

Three ways to look at a `.bdf` without a browser, using the same Go code (`raster`, `thumbnail`, `Document.SearchText`) a server would.

```sh
bdf thumbnail -size 256 out.bdf thumb.png    # -mode auto|crop|fit, -view, -sheet-dpi, -font-dir, -no-system-fonts
bdf text out.bdf text.json                   # metadata and per-page text as JSON (default: stdout)
bdf render -page 2 -scale 2 out.bdf page2.png  # a sheet takes -width/-height too (from A1)
```

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

See the [API reference](../api.md) for the equivalent Go calls, and [Building and testing](development.md) for how `bdf generate` and `bdf demo` are used to regenerate the fixtures under `testdata/`.
