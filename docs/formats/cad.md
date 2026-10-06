# CAD drawings and plots

CAD drawings and plotter output carry vector geometry drawn to scale on paper, without the page or slide structure an office document keeps for itself. BDF gives each of these formats' own idea of a sheet — DXF's model space and paper-space layouts, Jw_cad's and SXF's drawing sheet, a CGM picture, a plotted page — one BDF page each, and resolves colors, line types and line weights the way a plotter or the format's own viewer draws them, down to the dark background CAD programs use for model space.

## Try it

Drop a file on the [viewer](https://shibukawa.github.io/bdf/viewer/), or try a sample: an [AutoCAD DXF drawing](https://shibukawa.github.io/bdf/viewer/?file=samples/layout.dxf) with a model-space page and a paper-space layout, a [Jw_cad drawing](https://shibukawa.github.io/bdf/viewer/?file=samples/shapes.jww), an [SXF drawing](https://shibukawa.github.io/bdf/viewer/?file=samples/shapes.p21), a [CGM drawing](https://shibukawa.github.io/bdf/viewer/?file=samples/shapes.cgm), or an [HP-GL/2 plot](https://shibukawa.github.io/bdf/viewer/?file=samples/shapes.plt).

## AutoCAD DXF

`converter/dxf` reads AutoCAD DXF drawings, text and binary, from R12 to 2018. DWG is not read, since its native format isn't published and the libraries that read it aren't compatible with this project's license.

Model space becomes one page fitted to the drawing, drawn on the dark background AutoCAD itself uses for model space by default. Each paper-space layout becomes a page sized to its own paper, with its viewports showing model space at their own scale. Layers, colors, line types and lineweights are resolved the way a plotter would draw them. Polylines with bulges and variable widths, splines drawn as exact Bézier curves, hatches (patterns, islands, gradients), blocks and block arrays with attributes, dimensions, leaders and multileaders, and single- and multiline text (formatting codes, wrapping under Japanese line-breaking rules, stacked fractions) are all drawn. SHX fonts are replaced by a sans-serif font and big fonts by an East Asian font, with the fonts actually used embedded as WOFF2 subsets; strings in older files are read in their declared code page (Shift_JIS and others).

| Option | Values | Default |
|---|---|---|
| `-param views=` | `all`, `model`, `layouts` | `all` |
| `-param background=` | `dark`, `light` | `dark` |

See [design.md §3.12](../design.md#312-dxf--bdf-変換器converterdxfの構造) for the internals.

## Jw_cad

`converter/jww` reads Jw_cad's .jww drawings by its published data format, from version 2 through current version 7 and later files. The DOS-era .jwc format and the shape files .jws/.jwk are not read.

A drawing becomes one page of its sheet, grown to include anything drawn outside it. By default it's drawn in the screen colors saved in the file, on the file's own background color — the way Jw_cad shows it on screen (`-param colors=print` switches to the printer's output colors on white, `colors=mono` to black on white), with the printed line widths and line types. Lines, arcs and ellipses, points, text in Jw_cad's fixed pitch (full-width characters as wide as the text size, half-width characters half that, vertical text too), dimensions, solids including circle solids, and nested blocks are drawn. Hidden layers and auxiliary lines are left out, since Jw_cad doesn't print them either.

| Option | Values | Default |
|---|---|---|
| `-param colors=` | `screen`, `print`, `mono` | `screen` |

See [design.md §3.13](../design.md#313-jw_cad--bdf-変換器converterjwwの構造) for the internals.

## SXF

`converter/sxf` reads SXF, the CAD exchange format used for Japanese public-works deliveries (電子納品), versions 2 through 3.1, in both of its encodings: STEP AP202 files (.p21, and .p2z, a zipped .p21) and the feature-comment files (.sfc) CAD programs use to exchange drawings with each other.

A drawing becomes one page of its sheet, grown to include anything drawn outside it, on the background color the drawing itself names — black when it names none, matching the way SXF viewers show drawings (`-param background=light` switches to white paper, with white lines and text redrawn in black). Predefined and user-defined colors, line types and widths, lines, polylines, circles, arcs, ellipses, splines, clothoids, point markers, text placed at any of its nine anchor points (rotated, slanted, spaced, vertical), compound figures (placed, scaled, nested, or in geodetic coordinates), dimensions, leaders and balloons, color fills, hatching, and blank areas are all drawn; hidden layers are left out. Files written by the SCADEC library, which most SXF-capable CAD programs use to read and write the format, are read the way SCADEC itself reads them back — SCADEC diverges from the written specification in a few details.

| Option | Values | Default |
|---|---|---|
| `-param background=` | `file`, `light` | `file` |

See [design.md §3.14](../design.md#314-sxf--bdf-変換器convertersxfの構造) for the internals.

## CGM

`converter/cgm` reads Computer Graphics Metafiles (ISO/IEC 8632, versions 1–4), the format used for CAD plotter output, technical illustration exchange (S1000D, ATA) and WebCGM. Both the binary encoding and the clear-text encoding are read, gzip-compressed files (.cgz) too.

Each picture in the metafile becomes one page, sized to its own metric scale where the file gives one; a picture with no scale (an "abstract" picture) is fitted to A4. Lines with their types, widths and caps, markers, filled areas (hollow, solid, hatched, patterned, gradient) and their edges, closed figures with holes, circles, ellipses, circular, elliptical, hyperbolic and parabolic arcs, Bézier curves and B-splines, text (rotated and slanted by its orientation vectors, in any of four text paths including vertical, aligned, boxed to a fixed size, or appended to a previous string), cell arrays and tiles (JPEG, PNG, CCITT fax), and segments and their copies are drawn; WebCGM application structures marked invisible are left out. Strings are read in their declared character set (JIS X 0208 and other ISO 2022 sets, UTF-8, UTF-16), and undeclared Shift_JIS text, common in metafiles from Japanese CAD programs, is recognized on its own.

`converter/cgm` has no `-param` options.

See [design.md §3.20](../design.md#320-cgm--bdf-変換器convertercgmの構造) for the internals.

## HP-GL/2

`converter/hpgl` reads the plot files HP plotters and printers produce: HP-GL/2 and the older HP-GL of pen plotters and cutters. These arrive either bare or wrapped in the PJL job and PCL/HP RTL escape sequences that files sent to a large-format plotter (a DesignJet driver printing to a file) or a PCL printer carry, including the raster images HP RTL embeds.

Each plotted page becomes a page, sized from the plot size (PS) grown to include anything drawn outside it, shown in the coordinate system the file's own rotation (RO) turns it to — the way plotter drivers rotate a plot to fit the paper — or, for a PCL job, the printer's own page and picture frame. Pens draw in the file's palette colors and widths (`-param colors=mono` forces black), with fixed, adaptive and user-defined line types, caps and joins; arcs and circles (drawn as the plotter's own coarse chords when a chord angle wider than the default is given); Bézier curves; encoded polylines; and polygons, rectangles and wedges filled solid, hatched, cross-hatched, shaded or with a raster pattern. Labels become text: the plotter's built-in stick font is drawn with a fixed-pitch font at its cap height and pitch, other designated fonts with a font of a matching kind, all rotated, slanted, mirrored, and placed in any of the four label text paths and any of HP-GL/2's label origins, reading Roman-8 and the other built-in symbol sets as well as Japanese (Shift_JIS, 16-bit JIS). Raster images (compression methods 0–3, 5 and 9, indexed or direct colors) are scaled down to BDF's resolution cap as they're read.

| Option | Values | Default |
|---|---|---|
| `-param colors=` | `pens`, `mono` | `pens` (the file's own pen colors) |

See [design.md §3.22](../design.md#322-hp-gl2--bdf-変換器converterhpgl) for the internals.
