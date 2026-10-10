# Conversion and output

What the converters share, and what is made from a converted document. It works the same for every format, so it is gathered here rather than on the pages of the formats.

## Formulas

One formula engine lays out all of it — Word's Office Math, PowerPoint's and Excel's formulas (built from Office Math, not the fallback image some files also save), HTML and EPUB's MathML (including the MathML that KaTeX, MathJax and Wikipedia render alongside their own display), and the LaTeX in Markdown and draw.io labels (`math=1`). It uses the constants and glyph variants of an OpenType MATH font (STIX Two Math, Cambria Math, Latin Modern Math and others) for fractions, radicals, subscripts and limits, large operators, stretchy brackets and radicals built from glyph variants and assembly parts, matrices, aligned equations and accents. Search and copy produce a linear form — `x=(−b±√(b^2−4ac))/(2a)` — findable even typed with a plain hyphen-minus. See [design.md §3.23](../design.md#323-数式internalmathlayoutconverterinternalequationformula) for the internals.

The same engine is available by itself as the public package `formula`: it turns a formula in LaTeX or MathML into an object of paths only, to draw into an image, on an Ebitengine screen or on a browser canvas ([Rendering](../rendering.md#formulas)).

## Scores and playback

MML, MIDI and MusicXML are engraved as staff notation by their converters. A view carries that music as a Standard MIDI File, with cues tying playback time to a place on the page ([spec §4.4](../spec.md#44-演奏play)). The demo viewer plays it with the Web Audio API — an oscillator per General MIDI instrument family, synthesized percussion, square-wave MML chiptunes. It shows a cursor on the system currently playing, turns pages or scrolls to follow it, and starts playback from a system you click. See [design.md §3.27](../design.md#327-楽譜と演奏convertermmlconvertermidiconvertermusicxmlconverterinternalmusic).

An audio file's view plays through the same `play` member, but carries the file itself instead of a Standard MIDI File: the viewer hands it to the browser's `<audio>` element, so no converter or viewer decodes sound. Its cues are in milliseconds and point at lines of the card — the lyrics when they carry their times (ID3 `SYLT`, or LRC text in a lyrics field), otherwise the chapters — so the viewer marks the line being sung and plays from a line you click. See [design.md §3.31](../design.md#331-音声ファイル--bdf-変換器converteraudio).

## Metadata

The manifest can carry Dublin Core metadata (title, creator, subject, language, dates) inherited from a PDF's document info, PowerPoint's/Excel's/Word's core properties, Visio's document properties, Photoshop's document XMP, HTML `meta` elements, Markdown front matter, an EPUB package document, an image's XMP/EXIF/IPTC, or an audio file's tags (ID3, iTunes-style MP4 metadata, Vorbis comments, RIFF INFO).

## Server-side thumbnails and search text

`imagebdf` runs the same instructions the viewer draws, in pure Go, over any page (or, for a sheet or a scroll view, any range). It draws antialiased paths, lines, clips, gradients and patterns, images, embedded fonts (unpacking WOFF2) and name-referenced fonts (found on the system), groups, soft masks, shadows and SVG images. Rendered against the same pages the browser golden tests use and compared at reduced size, its output averages under 4/255 difference on most pages — the gap comes from not hinting, kerning or applying ligatures, and not drawing AVIF images. `thumbnail` picks a layout by the kind of document: Word/HTML/Markdown/scores/tall PDF get the top-left square of page one (an audio file's card too: the square is its cover); Excel/CSV get the range starting at A1; slides/drawings/images/EPUB covers get the whole first page — written as PNG, JPEG or WebP. `Document.SearchText` returns metadata and per-page text for a search index. Neither is encrypted, so neither is written for an encrypted document unless asked (`-allow-plaintext`). See [design.md §3.25](../design.md#325-サーバー側のサムネイルと検索用テキストimagebdfthumbnailsearchtext).

Try the [thumbnail](https://shibukawa.github.io/bdf/thumbnail/) and [search text](https://shibukawa.github.io/bdf/text/) pages — they run this same code, built as WebAssembly, on a file you drop on them.
