# InDesign

An InDesign document is a page layout: frames placed on pages, text flowing from frame to frame across them, and master pages whose items every page inherits. BDF reads the IDML package InDesign exports (File → Export → InDesign Markup) and lays the text out itself, with the fonts it then embeds, rather than rasterizing pages or going through a PDF export. Each page of the document becomes a page in BDF; the viewer turns them as a book, in the reading direction of the binding.

## Try it

Drop a file on the [viewer](https://shibukawa.github.io/bdf/viewer/), or try a sample: [basic.idml](https://shibukawa.github.io/bdf/viewer/?file=samples/basic.idml) (three pages of shapes, images and a story threaded through them), [vertical.idml](https://shibukawa.github.io/bdf/viewer/?file=samples/vertical.idml) (a right-to-left bound document in Japanese vertical text), or a magazine page InDesign exported, [interview.idml](https://shibukawa.github.io/bdf/viewer/?file=samples/interview.idml).

## InDesign Markup (.idml)

`converter/idml` reads IDML, the XML package InDesign exports and opens again: the design map, the spreads with their pages and items, the master spreads, the stories, and the resources (swatches, stroke styles, paragraph and character styles). The native `.indd` format is binary and undocumented and is not read: export the document as IDML first.

What's read:
- Pages, in document order, sized from each page's own bounds; a document bound on the right (Japanese books) opens right to left. Sections number the pages, and the page number marker in a text frame shows that number.
- Master pages: the items of the master a page applies (and of the master's own master) become a `master` layer, chosen by the page's side of the spine; items the page overrides are left out. Hidden layers and invisible items are skipped.
- Rectangles, ovals, polygons and lines from their Bézier paths, rounded corners, groups (with their transforms), fills from swatches of every kind — CMYK, RGB and Lab colors, tints, linear and radial gradients — and opacity, strokes with weight, dashes (the canned dashed and dotted styles, custom dashed and dotted styles), caps, joins and alignment inside or outside the path.
- Placed images: raster images embedded in the document, or linked ones given with the document (by file name: `bdf generate -with photo.jpg`, or a server's file list), scaled and clipped by their frames. TIFF images are converted.
- Text: stories flow through their chains of threaded text frames, across pages, with columns, insets and vertical alignment of each frame; a story set vertically flows right to left, with the Japanese rules for upright characters. Paragraph and character styles (and the styles they are based on), justification, indents, spacing, leading (fixed or automatic), tabs, bullets and numbering, font families and styles, size, tracking, superscripts and subscripts, baseline shift, underline and strikethrough, capitals, and languages. The text is laid out by the DrawingML text engine the PowerPoint, Visio and Word converters share, so line breaking (Japanese kinsoku included), the fonts and the text layer's structure are the same.
- Metadata from the package's XMP, as Dublin Core.

Not drawn, with a warning: tables and footnotes in text, objects anchored in text, placed EPS, PDF and PICT files, arrowheads, blend modes other than normal, striped strokes (drawn solid), text overset beyond the last frame of a chain. Fonts the document names are looked for on the machine (and under `-font-dir`); a font that is not found is replaced by a metric-compatible or generic one, and the layout changes accordingly.

The converter was checked against documents InDesign itself exported (CS5.5 to InDesign 2020), from the regression tests of [SimpleIDML](https://github.com/Starou/SimpleIDML) (BSD 3-clause), kept in `converter/idml/testdata/simpleidml/`; the viewer samples above include one of them.

`converter/idml` has no `-param` options. `-pages` selects pages, and the general font flags `-font-dir`, `-fonts`, `-no-subset`, `-no-woff2`, `-no-system-fonts` and `-ignore-fstype` apply as they do to PowerPoint.

See [design.md §3.32](../design.md#332-indesign--bdf-変換器converteridmlの構造) for the internals.
