# Features

## Three layout models

- **Fixed pages** — slides, drawings, PDFs: a fixed size per page.
- **An endless plane** — worksheets: drawn in tiles, with frozen panes, row/column headers and gridlines. A sheet is never cut into printed pages.
- **Flow** — word processor documents: pages you can read as pages, or as one continuous scroll, plus an optional page-free view laid out as a single long column at the body's width.

## Multiple views, and links between them

A document can hold several views — a workbook's sheets, a draw.io diagram's pages, a KiCad project's schematic sheets and board layers. The viewer switches between them with tabs, like a spreadsheet's sheet tabs. A document can also link from one view into another (`#view=ID`) — a draw.io page linking to another page, a KiCad hierarchical sheet linking to its parent.

## Reading pages, or turning them

The demo viewer lays pages out in a scroll, or shows them one at a time or as a spread (two pages side by side, left- or right-bound) and turns them like a book's. A page turns with the keyboard, page buttons, or a tap (past the spine goes forward, before it goes back), or by dragging a page's corner or outer edge — drawn in WebGL, with the facing pages pre-rendered as textures. Turning off "animate" (already off if the system asks to reduce motion) changes pages instantly. Text inside a page stays selectable while it's shown. Right-bound books (`direction: "rtl"` — vertical Japanese text, manga read right to left) turn right to left, and EPUB opens as a spread from the start (`?layout=pages`, `single`, `spread` or `spread-rtl` picks the layout up front).

## Search and selection

- **Full-text search** runs in the worker, across line breaks, with NFKC and kana normalization (searching in hiragana finds katakana and vice versa), and returns the rectangles of every hit.
- **Text selection** is a transparent DOM layer over the bitmap; copying rebuilds spaces and line breaks from the layer's own markers, works across pages, and works in the continuous layout too.
- **Table cells select as cells**: drag from one cell to another and the rectangle between them is selected. The demo viewer's sheets (Excel, CSV, Parquet) select cells the way a spreadsheet does — drag, Shift, row/column headers, arrow keys — and copying puts them on the clipboard as tab-separated values and an HTML table, so they paste back into a spreadsheet as cells, not as a picture.

## Accessible text

A structural MARK layer carries headings, lists, tables, figures with alt text, links and language to screen readers. It's converted from tagged PDF, PowerPoint's and Word's structure, Excel/CSV/Parquet's cells and table headers, and HTML/Markdown/EPUB's elements.

## Formulas

One formula engine lays out all of it — Word's Office Math, PowerPoint's and Excel's formulas (built from Office Math, not the fallback image some files also save), HTML and EPUB's MathML (including the MathML that KaTeX, MathJax and Wikipedia render alongside their own display), and the LaTeX in Markdown and draw.io labels (`math=1`). It uses the constants and glyph variants of an OpenType MATH font (STIX Two Math, Cambria Math, Latin Modern Math and others) for fractions, radicals, subscripts and limits, large operators, stretchy brackets and radicals built from glyph variants and assembly parts, matrices, aligned equations and accents. Search and copy produce a linear form — `x=(−b±√(b^2−4ac))/(2a)` — findable even typed with a plain hyphen-minus. See [design.md §3.23](design.html#323-数式converterinternalequation) for the internals.

## Music

MML, MIDI and MusicXML are engraved as staff notation by their converters. A view carries that music as a Standard MIDI File, with cues tying playback time to a place on the page ([spec §4.4](spec.html#44-演奏play)). The demo viewer plays it with the Web Audio API — an oscillator per General MIDI instrument family, synthesized percussion, square-wave MML chiptunes. It shows a cursor on the system currently playing, turns pages or scrolls to follow it, and starts playback from a system you click. See [design.md §3.27](design.html#327-楽譜と演奏convertermmlconvertermidiconvertermusicxmlconverterinternalmusic).

## Metadata

The manifest can carry Dublin Core metadata (title, creator, subject, language, dates) inherited from a PDF's document info, PowerPoint's/Excel's/Word's core properties, Visio's document properties, Photoshop's document XMP, HTML `meta` elements, Markdown front matter, an EPUB package document, or an image's XMP/EXIF/IPTC.

## Encrypted input

A password-protected input (a read-password Office document, a user-password PDF) is opened with its password and converted. The bdf output is then encrypted with the same password, sealed part by part (AES-256-GCM), so range requests and split layouts keep working. The viewer decrypts with WebCrypto, and a server never stores the password ([spec §3.5](spec.html#35-暗号化)).

## Segments sealed for each request

A server can keep a document and hand it to a logged-in reader a few pages at a time: each segment is sealed for a key pair the reader's browser made for that one request, and neither side keeps the key afterwards. A recording of the traffic stays sealed even if the server's keys or the reader's session are taken later. The Go package `segment` answers the requests, `SegmentLoader` of `@bdf/core` makes them, and the rendering worker fetches the pages it is asked to draw ([spec §3.6](spec.html#36-区間の文書segment), [secure-reader](examples/secure-reader.html)).

## Server-side thumbnails and search text

`raster` runs the same instructions the viewer draws, in pure Go, over any page (or, for a sheet or a scroll view, any range). It draws antialiased paths, lines, clips, gradients and patterns, images, embedded fonts (unpacking WOFF2) and name-referenced fonts (found on the system), groups, soft masks, shadows and SVG images. Rendered against the same pages the browser golden tests use and compared at reduced size, its output averages under 4/255 difference on most pages — the gap comes from not hinting, kerning or applying ligatures, and not drawing AVIF images. `thumbnail` picks a layout by the kind of document: Word/HTML/Markdown/scores/tall PDF get the top-left square of page one; Excel/CSV get the range starting at A1; slides/drawings/images/EPUB covers get the whole first page — written as PNG, JPEG or WebP. `Document.SearchText` returns metadata and per-page text for a search index. Neither is encrypted, so neither is written for an encrypted document unless asked (`-allow-plaintext`). See [design.md §3.25](design.html#325-サーバー側のサムネイルと検索用テキストrasterthumbnailsearchtext).

Try the [thumbnail](../thumbnail/) and [search text](../text/) pages — they run this same code, built as WebAssembly, on a file you drop on them.
