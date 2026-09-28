# Word, HTML, Markdown

Word documents, HTML pages and Markdown files are all prose laid out into lines, paragraphs, tables and pages, so bdf lays all three out with the same engine (`converter/internal/wordproc`, shared with EPUB) rather than three separate ones. BDF itself has no relayout. The line breaks and page breaks a viewer draws are exactly the ones the converter chose — so each converter finishes, once, the layout work Word or a browser would otherwise redo every time the file opens.

## Try it

Drop a file on the [viewer](../../viewer/), or try a sample: a [Word document](../../viewer/?file=samples/basic.docx), the same kind of document [set in Japanese vertical text](../../viewer/?file=samples/vertical.docx), an [HTML article in reader mode](../../viewer/?file=samples/article.html), or a [Markdown file](../../viewer/?file=samples/basic.md).

## Word .docx (`converter/docx`)

`converter/docx` reads a `.docx` and lays WordprocessingML out itself — the same line-breaking and page-breaking work Word does every time it opens the file. It handles line breaking (Japanese kinsoku rules and spacing, tab stops and leaders, justification, the document grid), lists and numbered paragraphs, tables (table styles, merged cells, rows split across a page break, repeated header rows), floating pictures and text boxes with text wrapping around them, columns, sections, headers and footers with page numbers, footnotes, East Asian vertical text (upright kanji and kana, vertical punctuation, turned Latin text and tables), and formulas (Office Math).

Two views come out of a single document: the pages themselves (a flow view with header, body and footer layers), and a second, page-free view that lays the same paragraphs and tables out again as one long column at the body's width — the same idea as Word's draft or web layout view, sharing the pages' fonts and images rather than duplicating them. Drawings go through the same DrawingML renderer PowerPoint uses, and fonts are embedded the same way.

| Option | Values | Default |
|---|---|---|
| `-param views=` | `both`, `pages`, `scroll` | `both` — pages and the scroll view |

See [design.md §3.9](../design.html#39-word--bdf-変換器converterdocxの構造) for the internals.

## HTML .html / .xhtml / .mhtml (`converter/html`)

`converter/html` lays a page out like a browser's reader view. The author's CSS is discarded, and elements are laid out with a fixed style sheet instead, by the same layout engine the Word converter uses. The article is picked out of a web page with go-readability (a Go port of Mozilla's Readability), and headings, paragraphs, lists, quotations, code, tables (with merged cells and content-driven column widths), figures, links and MathML formulas are all laid out the same way a Word document's would be. Since bdf never relays text out, the result is a scroll view of a fixed width (36 ems of text by default; `-param views=pages` adds A4 pages). Images are read from files beside the page, from an MHTML archive, from `data:` URLs, and from the network (fetched by default; `-param remote=false` leaves them out). SVG — both `.svg` files and inline `<svg>` elements — is stored as-is and drawn by the viewer, sized the way a browser would size it. Fonts aren't embedded: they're referred to by name, so a viewer draws the text in its own fonts the way it would a web page (`-fonts embed` embeds a subset instead).

| Option | Values | Default |
|---|---|---|
| `-param extract=` | `auto`, `article`, `none` | `auto` — pick the article out of a page that looks like one |
| `-param views=` | `scroll`, `pages`, `both` | `scroll` |
| `-param width=` | text width of the scroll view, in points | 36 ems of the text |
| `-param size=` | text size, in points | `12` |
| `-param font=` | font family of the text | an installed sans-serif family |
| `-param mono=` | font family of code | an installed monospaced family |
| `-param remote=` | `true`, `false` | `true` — fetch images from the network |
| `-param base=` | URL the document came from, for resolving relative links and images | — |

See [design.md §3.16](../design.html#316-htmlmarkdown--bdf-変換器converterhtmlconvertermarkdownの構造) for the internals.

## Markdown .md (`converter/markdown`)

`converter/markdown` renders a document to HTML with goldmark (CommonMark plus GitHub's tables, task lists, strikethrough and autolinks, footnotes and definition lists) and lays it out the same way `converter/html` lays out a page — without picking an article out, since the whole file is the content. Raw HTML that README files often carry (`<p align="center">`, `<details>`) is laid out along with the rest, and headings get GitHub's own ids, so links to them keep working. Math is written the way GitHub renders it — `$…$`, `$$…$$`, and ```` ```math ```` code blocks, all in LaTeX. YAML or TOML front matter becomes the document's Dublin Core metadata.

| Option | Values | Default |
|---|---|---|
| `-param views=` | `scroll`, `pages`, `both` | `scroll` |
| `-param width=` | text width of the scroll view, in points | 36 ems of the text |
| `-param size=` | text size, in points | `12` |
| `-param font=` | font family of the text | an installed sans-serif family |
| `-param mono=` | font family of code | an installed monospaced family |
| `-param remote=` | `true`, `false` | `true` — fetch images from the network |
| `-param base=` | URL the document came from, for resolving relative links and images | — |

See [design.md §3.16](../design.html#316-htmlmarkdown--bdf-変換器converterhtmlconvertermarkdownの構造) for the internals — the HTML and Markdown converters share one design-notes section, since Markdown is rendered to HTML before this layout stage.
