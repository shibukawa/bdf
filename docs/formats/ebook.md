# EPUB

EPUB is a zip archive of XHTML, CSS and images, so a reflowable book converts through much the same reader-mode layout as HTML. What's different: a book is several documents read in order, and its reader wants book pages, not a scrolling web page. Many Japanese books are set vertically, and comics and photo books use a fixed layout with no text to flow at all. `converter/epub` handles both EPUB 3 and EPUB 2.

## Try it

Drop a file on the [viewer](https://shibukawa.github.io/bdf/viewer/), or try a sample: a [reflowable book](https://shibukawa.github.io/bdf/viewer/?file=samples/basic.epub), a [vertical Japanese book](https://shibukawa.github.io/bdf/viewer/?file=samples/vertical.epub), or a [fixed-layout book](https://shibukawa.github.io/bdf/viewer/?file=samples/fixed.epub).

## EPUB (converter/epub)

### Reflowable books

The package document (found through `META-INF/container.xml`) gives the metadata, manifest and reading order (spine). Content documents are read in spine order — documents marked `linear="no"`, such as end notes only reached by a link, are moved to the end — and each becomes a section that starts a new page, the way a Word section does. Because EPUB's XHTML breaks a standard HTML parser (an empty `<a id="p5"/>` swallows what follows, for instance), it's read as XML into the same tree shape the HTML converter uses. Links between chapters resolve to the right page, and images load straight from the zip archive. EPUB doesn't allow documents to pull in outside images, so none are fetched from the network.

Almost none of a book's own CSS is kept — layout, colors and fonts stay whatever the reader-mode stylesheet gives them — except for a handful of properties that carry meaning for how a book is meant to be read: `writing-mode` (see below), tate-chu-yoko, emphasis marks, text alignment, `display: none`, and the pixel size an image is meant to render at (used for full-page illustrations and the "gaiji" character images some Japanese publishers embed). MathML formulas are laid out by BDF's formula engine, matching Word and HTML.

### Vertical Japanese text

A chapter's writing direction comes from `writing-mode` on `body` or `html` (including the `-epub-`/`-webkit-` prefixed forms publishers commonly use); failing that, from Kindle's `primary-writing-mode`; failing that, from the language and `page-progression-direction` together. A vertically-set chapter lays out the same way Word's vertical sections do: upright kanji and kana, vertical punctuation forms, rotated Latin text, tate-chu-yoko for short horizontal runs, and emphasis dots. Right-bound books (the common case for vertical Japanese, and for right-to-left comics) get a view with `direction: "rtl"` ([spec §4.1](../spec.md#41-view-の種類)).

### Fixed-layout books

A book marked `pre-paginated` (or Kindle's `fixed-layout`) whose every document is just one picture — an `img`, an inline `<svg>`, an SVG document, or an image placed directly on the spine — skips the text layout engine entirely and becomes a page per picture, sized from the document's viewport metadata (or the image's own size, or the SVG viewBox), each picture centered and scaled to fit. If a book mixes flowed text pages with picture-only pages, BDF falls back to laying the whole thing out as a reflowable book (it doesn't interpret the absolute CSS positioning some such mixed books use).

### Formulas

`math` elements (MathML) in the XHTML of a book are laid out by the same formula engine as in HTML. In a chapter of vertical text a formula stays horizontal inside its line.

```html
<p>Pythagoras: <math><msup><mi>a</mi><mn>2</mn></msup><mo>+</mo><msup><mi>b</mi><mn>2</mn></msup><mo>=</mo><msup><mi>c</mi><mn>2</mn></msup></math>.</p>
```

See [formulas in Word, HTML and Markdown](document.md#formulas) for the layout.

### Options

| Option | Values | Default |
|---|---|---|
| `views` | `pages` (book pages), `scroll` (one long horizontal-writing column), `both` | `pages` |
| `paper` | `a5`, `a4`, `a6`, `b5`, `b6`, `letter`, or `WIDTHxHEIGHT` in millimeters | `a5` |
| `size` | text size, in points | `12` |
| `font` | font family of the text | an installed sans-serif family |
| `mono` | font family of code | an installed monospaced family |
| `width` | text width of the scroll view, in points | 36 ems of the text |

DRM-encrypted books are refused. See [design.md §3.24](../design.md#324-epub--bdf-変換器converterepubの構造) for the internals.
