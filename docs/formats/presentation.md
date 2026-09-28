# PowerPoint

A PowerPoint deck is built from DrawingML: shapes, text runs and tables placed on slides that inherit from a layout and a master, rather than one rasterized picture per slide. bdf reads that DrawingML directly, instead of converting the file through an office suite's PDF export. As a result, the pptx converter's output keeps the master/layout sharing and the paragraph and bullet-list structure of the original text.

## Try it

Drop a file on the [viewer](https://shibukawa.github.io/bdf/viewer/), or try a sample: [PowerPoint deck](https://shibukawa.github.io/bdf/viewer/?file=samples/features.pptx).

## PowerPoint (.pptx)

`converter/pptx` reads the OPC zip package (.pptx, .pptm, .ppsx, .ppsm, .potx, .potm) and turns its DrawingML into bdf's own drawing instructions. The XML is parsed into a generic element tree rather than typed Go structs, since DrawingML has many optional elements and several layers of style inheritance. The same OPC and element-tree reader (`converter/internal/ooxml`) is shared with the Excel, Word and Visio converters, and the DrawingML drawing code (`converter/internal/ooxml/drawingml`) is shared with Excel's shapes and charts.

What's read:
- Slide masters, layouts and slides. Each slide becomes a page with up to four layers: `background` (the first one found looking slide → layout → master), two `master` layers (the master's own shapes, then the layout's), and `body` (the slide's shapes). Slides that share a master and layout produce the same master/layout layer objects, so they're stored once.
- Placeholder inheritance: a slide's placeholder inherits position, shape, fill, line, effects and text styles from the layout's matching placeholder (matched by `idx`, then by type), which in turn inherits from the master's.
- The 187 preset shapes defined by ECMA-376, custom geometry, gradient/image/pattern fills, line ends, joins, dashes and arrowheads, and shadows.
- Text layout done by bdf itself, not by the browser: line breaking with Japanese rules, indentation and bullet lists (including Wingdings/Symbol bullet glyphs), tabs, line and paragraph spacing, alignment, autofit, columns, and vertical text (`vert`, `vert270`, `eaVert`).
- Tables (with PowerPoint's built-in table styles) and charts (bar, line, area, pie/doughnut and scatter, drawn from their cached values), SmartArt's own drawing part, and embedded EMF/WMF pictures, which are replayed as their original GDI drawing commands rather than stored as bitmap images.
- Formulas written as Office Math (the fallback image some files also save alongside them is ignored), laid out by the same formula engine used for Word, Excel and HTML.
- Fonts, resolved the same way as for Word and Excel: the real family, then a metric-compatible substitute, then a generic family; only the glyphs actually used are embedded, as WOFF2.

How it's laid out: each slide becomes one page, in slide order; hidden slides are excluded unless asked for. Titles become headings, bullet levels become nested lists, and tables and pictures/shapes/charts become tables and figures with alt text in the accessible text layer. The language of each run of text is inferred from its script when that would otherwise be misread from a template's default language tag.

Options (`-param`, from `bdf generate -h`):

| Option | Values | Default |
|---|---|---|
| `-param hidden=` | `true` | hidden slides excluded |

`-hidden` is shorthand for `-param hidden=true` (Excel's `-hidden` does the same for hidden sheets). `-pages` selects which slides to convert, and the general font flags `-font-dir`, `-fonts system`, `-no-subset`, `-no-woff2`, `-no-system-fonts` and `-ignore-fstype` apply to PowerPoint the same way they do to Word and Excel.

See [design.md §3.4](../design.md#34-powerpoint--bdf-変換器converterpptxの構造) for the internals.
