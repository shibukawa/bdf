# The bdf command

```
bdf generate [flags] <input> <out.bdf | dir/>
                                   convert a document or a drawing (bdf generate -h for flags and formats)
bdf fill [flags] <template.xlsx | .docx | .pptx | .pdf | .json | .txt | .bdft> <out.bdf | out.txt | dir/>
                                   fill a template of a form with data: an Office document with tags, a
                                   PDF with form fields, a template of pdfme, a template of text (bdf fill -h for flags)
bdf template compile [flags] <template> <out.bdft>
bdf template schema [-json-schema] <template> [out.json]
                                   check a template and write it compiled; write what it asks of data
bdf pdf [flags] <file.bdf | dir> <out.pdf>
                                   write the pages as a PDF file (bdf pdf -h for flags)
bdf print [flags] -printer <url | name> <file.bdf | dir | input>
                                   print on an IPP printer or a CUPS queue, without a driver
                                   (-server: on a queue of a print server)
bdf printers [flags]               list the IPP printers on the network, the queues of CUPS
                                   (-cups), what a printer supports (-url), or the queues and
                                   printers of a print server (-server)
bdf printd -config <printd.yaml>   run the print server in the foreground
bdf jobs [flags] [show | cancel | hold | release | retry | move | priority | password | accounting]
                                   list and change the jobs of a print server
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
bdf sign -key key.pem -cert chain.pem [flags] <file.bdf | dir> [out.bdf]
                                   add a signature (CMS, X.509) after those a document has
bdf verify [-roots roots.pem] <file.bdf | dir>
                                   check the signatures of a document
bdf demo <file.bdf | dir/>         write the fixture document (dir/ ends with a slash)
bdf conformance [-testdata dir] <dir>
                                   write the conformance suite for other readers (package conformance)
```

An encrypted document is read with the password in `$BDF_PASSWORD`; `split` and `join` copy it without needing the password, and keep the signatures of a document.

## generate

Converts one input into BDF. The format is detected from the input's content, or set explicitly with `-format`; the output is a single file, or a directory ending in `/` for the split form. `-pages 1-3,5,8-` selects pages, slides, sheets or artboards; `-thumbnail` and `-text` write a preview next to the output (not for an encrypted output, unless `-allow-plaintext`); `-dc name=value` overrides a Dublin Core element (`-dc name=` removes it). See [Formats](../formats/index.md) for the `-param` options specific to each input format, or `bdf generate -h` for the current, complete list.

```sh
bdf generate report.pptx report.bdf              # single file
bdf generate report.pptx out/                    # split form
bdf generate -pages 1-3 book.pdf book.bdf         # first three pages only
bdf generate -password-file pw.txt secret.xlsx out.bdf  # a protected input
```

## fill, template

Fill a template of a form (an Excel workbook, a Word document or a PowerPoint presentation whose text holds tags, a PDF with form fields, a template of pdfme, or a text file of the lines of a form with tags) with data. The format is told from the content. See [Making forms](../forms.md) for how to write one.

```sh
bdf fill -data order.yaml -data items=lines.csv invoice.xlsx invoice.bdf   # fill, as a document of the printed pages
bdf fill -data order.yaml notice.docx notice.bdf                           # a Word document as its pages, a presentation as its slides
bdf fill -each orders.csv invoice.xlsx 'out/{{no}}.bdf'                    # a form for each record
bdf template compile -sample sample.yaml invoice.xlsx invoice.bdft         # check it and write a compiled template
bdf template schema invoice.bdft                                           # the names its tags look up, as JSON
bdf template schema -json-schema invoice.xlsx                              # the JSON Schema of the data, to make a form of
bdf fill -data applicant.json application.pdf application.bdf              # a PDF with form fields: the values in the places of its fields
bdf fill -data invoice.json pdfme-template.json invoice.bdf                # a template of pdfme, as its generator fills it
bdf fill -data order.json slip.txt slip.bdf                                # a template of text: characters on a grid of columns and lines
bdf fill -data order.json -param encoding=shift_jis slip.txt slip.txt      # the form of a template of text as plain text
```

`-data` is a JSON or YAML object (`file`) or a value under a name (`name=file`; a CSV or TSV file or the sheets of a workbook are lists of objects), and repeats. `-strict` makes a tag whose name the data lacks an error, and `-files` names the directory of the pictures the data names. The font and image flags (`-font-dir`, `-no-system-fonts`, `-fonts`, `-images` …), `-param` and `-pages` mean what they mean for `generate` (`-pages` selects the sheets of a workbook, the slides of a presentation's template, and the pages of the form of a Word template or a template of text). An output path ending in `.txt` writes the form of a template of text as plain text, with `-param encoding=`, `newline=crlf` and `paper=preprinted` ([Templates of text](../forms.md#templates-of-text)). A `.bdft` is a BDF document: `ls`, `thumbnail` and `render` open it as it is.

## pdf

Write the pages of a document as a PDF file: its text as text in the embedded fonts, its drawings as vectors ([Writing a PDF](../forms.md#writing-a-pdf)).

```sh
bdf pdf invoice.bdf invoice.pdf
bdf pdf -view print1 book.bdf sheet1.pdf            # choose views (repeatable)
bdf pdf -skip background -skip master slip.bdf slip.pdf   # leave layers out by their roles
bdf pdf -tagged report.bdf report.pdf               # a tagged PDF: with the structure of the document, for screen readers
bdf pdf -pdfa 2a invoice.bdf invoice.pdf            # a PDF/A file (2b, 2u, 2a, 3b, 3u, 3a), for keeping
bdf pdf -sign-key alice.key -sign-cert alice-chain.pem -sign-field party_a.sign contract.bdf contract.pdf
                                                    # signed (PAdES), shown on the picture of a signature field
```

`-sign-key` and `-sign-cert` sign the file as PAdES does (a CMS signature over the bytes of the file around it, `/SubFilter /ETSI.CAdES.detached`), with a key and certificates as `bdf sign` takes them. Acrobat and other readers check it. `-sign-field` names the picture of a signature field of the form (`{{%sign | signature}}`, a signature field of a PDF form, a `signature` of pdfme): the signature shows on its box, over the drawn signature; without it the signature shows on no page. `-tsa` adds a time-stamp (B-T), `-certify 1|2|3` makes it the certification of the file with the changes it permits after it, and `-sign-name`, `-sign-reason` and `-sign-location` are what it says of the signer. It works with `-tagged` and `-pdfa`. The signature of a `.bdf` (`bdf sign`) does not carry over to the PDF: whoever writes the PDF signs it.

By default all the views of pages (`fixed`, `flow`) are written, and a document without pages becomes one page. `-font-dir` and `-no-system-fonts` say where the fonts the document refers to by name are looked for. `-ignore-fstype` embeds fonts of the system that do not allow it (only with the rights to; by default their text is drawn as outlines). `-no-compression` writes the streams as they are, to read them. An encrypted document needs `-allow-plaintext` (the PDF is not encrypted).

## print, printers

Print a document on an IPP printer, without a driver of the system. When the printer takes PDF the document is sent as a PDF file (its drawings stay vectors, its text is drawn at the printer's resolution); when it takes PostScript and not PDF, as PostScript; when it takes neither, the pages are drawn here as PWG Raster (what IPP Everywhere requires), Apple Raster (AirPrint) or PCLm (Mopria). A CUPS queue of macOS or Linux (`ipp://localhost:631/printers/<name>`) prints the same way. On Windows, since Windows Ready Print, printers speak IPP themselves: give the printer's URL.

```sh
bdf printers                                         # the IPP printers on the network (mDNS) and of USB on this machine
bdf printers -url ipp://192.0.2.10/ipp/print         # what one supports: media, sides, colour, resolutions, formats …
bdf printers -cups local                             # the CUPS queues of this machine

bdf print -printer ipp://192.0.2.10/ipp/print invoice.bdf
bdf print -printer "Canon MF Series" -sides long -media b5 -copies 3 -mono invoice.bdf   # found by name
bdf print -printer ipps://mfp.example.local/ipp/print -pages 1-3 report.xlsx            # converted first (a workbook as the pages Excel prints)
bdf print -printer ipp://localhost:631/printers/office -format pwg-raster -mono -bilevel -resolution 300 slip.bdf
bdf print -printer ipp://192.0.2.10/ipp/print -skip background slip.bdf                 # onto preprinted paper
bdf print -printer socket://192.0.2.30:9100 -sides long -media a4 invoice.bdf          # a printer without IPP: PostScript
bdf print -printer "lpd://192.0.2.31/lp1?format=pdf" invoice.bdf                       # PDF over LPD, with a PJL header

bdf print -profile list                                                                 # the models of labels, receipts and forms
bdf print -printer "socket://192.0.2.40:9100?profile=zd421-203" -media 104x150mm label.bdf      # ZPL labels
bdf print -printer "serial:///dev/ttyUSB0?lang=escpos&baud=19200" -drawer 1 receipt.bdf          # ESC/POS receipts, the drawer opened
bdf print -printer "winspool:///EPSON%20TM-T88VII%20Receipt?profile=tm-t88" receipt.bdf          # through the queue of a driver of Windows
bdf print -printer "usb:///dev/usb/lp0?profile=vp-f2000" -perforation-skip 12.7 slip.bdf         # ESC/P on a continuous form
bdf print -printer ipp://localhost:631/printers/zebra -profile zd421-300 label.bdf               # a raw queue of CUPS
```

- **Paper**: `-media` takes a PWG name (`iso_a4_210x297mm`) or a short one (`a4`, `a3`, `b4` and `b5` for the JIS sizes, `letter`, `legal`, `hagaki`, `chou3`, `kaku2` …). Without it, the paper of the size of the first page is picked among the printer's, the loaded ones first.
- **Format**: `-format auto` (the default) is PDF, else PostScript, else PWG Raster, Apple Raster or PCLm (`-format pdf|ps|pwg-raster|urf|pclm` asks for one). PostScript is Level 3: its text is drawn with Type 3 fonts of the outlines of the glyphs, and what PostScript cannot draw (transparency, blend modes, shadows, filters) is drawn as pictures of just the rectangles where it is, the rest staying vectors. Rasters are 600 dpi gray for black, 300 dpi sRGB for colour, black and white with `-bilevel`, at the resolutions the printer has. The back of a two-sided sheet is turned as the printer asks (`pwg-raster-document-sheet-back`, `urf-supported`'s `DM`, `pclm-raster-back-side`). `-raster-if-lossy` sends PWG Raster when the PDF file would leave out what the page draws (`FILTER`, Canvas's own composite operations).
- **Without IPP**: `socket://host:9100` (the raw port) and `lpd://host/queue` are sent PostScript, or PDF with `?format=pdf`, which gets a PJL header (`?pjl=false` leaves it out, `?pjl=true` puts one before PostScript too). The paper, copies, sides and media type of PostScript are asked in the file (`setpagedevice`), those of PDF in PJL; the rest (trays, output bins, finishings) is reported unsupported. Such a printer tells nothing of a job: bdf print returns once the printer has it.
- **Labels, receipts and continuous forms** take their printer's own language: ZPL (Zebra), SBPL (SATO), TPCL (Toshiba TEC), ESC/POS (Epson TM and the receipt printers like it) and ESC/P (Epson VP and the dot-matrix printers like it). The URL's `lang=` (or `format=`) names it and `profile=` the model (`-profile list`: its dots and width), `dpi=` and `width=` (dots, or `104mm`, `4in`) override them. They are reached over `socket://` and `lpd://`, `serial:///dev/ttyUSB0` (`serial://COM3` on Windows; `baud`, `bits`, `parity`, `stop`, `flow=rtscts|xonxoff`; a Bluetooth printer of the serial port profile is a serial port of the system), `usb:///dev/usb/lp0` (the device file of Linux and the BSDs; macOS has none, so print to the printer's CUPS queue), `winspool:///<queue>` (Windows: the queue of the maker's driver passes the data on as it is, `RAW`), or a raw queue of CUPS with `-format zpl|sbpl|tpcl|escpos|escp` or `-profile`. A page is drawn in black and white at the printer's dots: text and lines cut at the half, images dithered, the modules of QR codes and barcodes on whole dots so that they read at 203 dpi. A label is as long as `-media` (`104x150mm`, `4x6in`) or as the page; a receipt and a form as the page (shrunk to the width). A label printer prints the QR codes and barcodes of filled forms with its own commands (`-no-printer-barcodes` draws them); `-darkness` (-100 to 100), `-speed` (mm a second), `-media-tracking continuous|gap|mark`, `-label-gap` (mm, TPCL), `-cut none|page|job` (receipts are cut after each page), `-drawer 1|2` (the cash drawer, at the end) and `-perforation-skip` (mm, ESC/P) ask the rest. A form filled from a [template of text](../forms.md#templates-of-text) is printed by ESC/P in the text mode: its characters on the grid of columns and lines, in the printer's own fonts, which a dot-matrix printer prints much faster and darker through carbon copies than dots (the pitch, line spacing, page length and perforation skip of the template; Kanji in the Kanji mode). `-print-mode text` prints it so on ESC/POS too (Font A spaced to the grid), and `-print-mode raster` prints the dots; `-skip background` leaves out what the paper has printed already.
- **What the printer does not support** is left out with a warning, or with `-strict` is an error. Copies a printer cannot make are made in the document (with a blank sheet between the copies of an odd two-sided document).
- **Waiting**: the job is waited for until it is done (`-wait 10m`); `-wait 0` returns once it is sent. Ctrl-C cancels the job at the printer.
- **Exit status**: 0 printed, 1 canceled, aborted or another error, 2 usage, 3 the printer could not be reached (or refused the certificate or password), 4 the printer refused the job (or, with `-strict`, what it asks), 5 the printer is stopped, 6 the job is not done within `-wait`.
- **TLS**: the certificate of an `ipps://` printer is verified with the system's roots, else its fingerprint is recorded on the first connection in `bdf/printers.json` of the user's config directory and required after (as CUPS does; `-trust-store` moves it, `-insecure` verifies nothing). Basic authentication takes `-username` and `$BDF_PRINTER_PASSWORD` (or `-password-file`).
- **An encrypted document** is opened with `$BDF_PASSWORD`, and is not sent over a connection without TLS unless `-allow-plaintext`.
- **On a print server**: `-server` (or `$BDF_SERVER`, when there is no `-printer`) sends the document as it is to a queue of [`bdf printd`](../print-server.md) (`-queue`; the server's only queue by default), which converts and prints it. It writes the job's id and returns; `-wait` follows the job until it is done, with the same exit statuses (5 is a queue that takes no jobs, 3 a server not reached or a token refused). The ticket's flags are the same; `-priority` and `-hold` order and hold the job in the spool. The token is `$BDF_TOKEN` or `-token-file`; `-user` names the user of the job to a server without tokens or to a token that may act for others. `bdf printers -server` lists the server's queues and printers.

## printd, jobs

`bdf printd` runs the [print server](../print-server.md): a spool in one file, its queues and printers from a configuration file, and its HTTP API; with `ipp:` also its queues as IPP printers that the systems add without a driver, advertised with DNS-SD, and with `upstreams:` the jobs of another server's queue it fetches as the proxy of a branch. It runs in the foreground until SIGINT or SIGTERM (Ctrl-C or Ctrl-Break on Windows), and lets a service manager keep it running.

```sh
bdf printd -config printd.yaml                 # serve (the log is JSON on the standard error; -log text for people)
bdf printd -config printd.yaml -check          # ask each printer what it supports, without opening the spool
bdf printd -config printd.yaml -compact        # rewrite the stopped spool's file without the space of removed jobs
bdf printd -config printd.yaml -listen 127.0.0.1:9631
bdf printd -config printd.yaml -dns-sd example.com             # the DNS-SD records of the IPP printers, for a DNS zone
bdf printd -config printd.yaml -airprint-profile print.example.com > printers.mobileconfig   # for iPhones and iPads of other subnets

bdf jobs                                       # the jobs not done ($BDF_SERVER, or -server)
bdf jobs -all -queue office                    # with those done
bdf jobs show 42
bdf jobs cancel 42                             # also hold, release, retry
bdf jobs move 42 warehouse
bdf jobs priority 42 90
bdf jobs password 42                           # an encrypted document's password, from $BDF_PASSWORD or -password-file
bdf jobs accounting -from 2026-10-01 -to 2026-10-31 -csv -o october.csv
```

Without the manage role a token sees and changes its own jobs; `-user` lists a user's jobs and, with a token that may act for others, acts as that user. `-json` writes JSON.

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
bdf render -workers=1 -mono=bilevel -scale 8.333333 out.bdf print.tif  # the bands of a page one at a time: the least memory
```

`render -antialias=default|on|off` controls the edges of shapes, text and geometric clips. `default` and `on` preserve the existing smoothing at every scale; `off` skips edge coverage calculations. There is no dpi-based automatic switch. Transparency, blending, image interpolation, soft masks, blur and the reduced alpha of subpixel-width lines are preserved. `thumbnail` and generated previews keep their existing smoothing. See [Rendering](../rendering.md#raster-quality-and-speed) for the Go option and output tradeoffs. `render -mono=gray|bilevel` writes a monochrome image instead of a colour one: 8 bits of luminance a pixel, or black and white (1 bit a pixel in a PNG; black where the gray value is below 128), for a page or a batch ([Monochrome output](../rendering.md#monochrome-output)). A PNG or TIFF page is drawn in bands (`-band=N` rows, default 128) and written as the bands are drawn, up to eight bands at a time on as many cores; `-workers=N` sets how many at a time (`-workers=1` for the least memory, about 80 MiB for an A4 page at 600 dpi), and the PNG is compressed by as many goroutines, at zlib level 1 by default: `-compression=normal` (6) or `best` (9, a half to a third of the size in one and a half to two times the time) for a file to keep.

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

## sign, verify

```sh
bdf sign -key alice.key -cert alice-chain.pem -role author contract.bdf signed.bdf
bdf sign -key bob.key -cert bob-chain.pem -role counterparty -tsa https://tsa.example/ signed.bdf
bdf verify -roots ca.pem signed.bdf
```

`sign` adds a signature after those a document has (spec §3.7): a CMS signature with the X.509 certificate of the signer, over the SHA-256 of the manifest and of every part, and of the signatures before it. The key is PEM (PKCS #8, EC or PKCS #1; a PKCS #12 file is written as PEM by `openssl pkcs12 -in key.p12 -nodes`), and `-cert` is the certificate of the key followed by its chain. `-role` and `-claims file.json` are signed with the document; `-tsa` has a time-stamp authority (RFC 3161) stamp the signature. Without an output, the file is signed in its place; a document of the split form is signed in its directory (`signatures.json`).

`verify` reports, for each signature, whether it holds (nothing changed after it: no part, no page, no signature before it) and whether its signer is trusted (`-roots`, or the roots of the system), at the time of its time-stamp when it has one. It exits with 1 when a signature does not hold or is not trusted, or when there is none; `-json` writes the report as JSON.

A signature finds changes; it does not prevent them. `split` and `join` keep the signatures, which hold in either form. A drawn signature in a form is a picture, not a signature of `sign`.

## demo

```sh
bdf demo out.bdf      # the fixture document used by browser integration tests and the demo viewer's default
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
