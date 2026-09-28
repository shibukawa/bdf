# Visio, draw.io

Visio and draw.io (diagrams.net) diagrams are shapes connected by routed edges, organized into the pages a diagramming tool itself keeps and switches between. bdf keeps that structure instead of flattening it. Each page in the file becomes a page in bdf, and the viewer switches between them with tabs — not laid out one after another the way a PDF export would.

## Try it

Drop a file on the [viewer](../../viewer/), or try a sample: a Visio drawing ([shapes.vsdx](../../viewer/?file=samples/shapes.vsdx)), a draw.io diagram of three pages ([multipage.drawio](../../viewer/?file=samples/multipage.drawio)), or an [AWS architecture diagram](../../viewer/?file=samples/aws.drawio).

## Visio (.vsdx, .vdx)

`converter/visio` reads Visio 2013 and later packages (.vsdx, .vsdm, .vstx, .vstm) and the XML drawings of Visio 2003–2010 (.vdx, .vtx) into one ShapeSheet model. The binary .vsd format is not read.

Each foreground page becomes a page, sized from the page's own width and scale; background pages (which a Visio drawing can chain) become a shared background layer drawn once and reused across every page that shares it. Shapes inherit unset cells from their master and its style, and cells a document theme sets are resolved from the theme and each shape's quick style — colors, line and fill schemes, and the connector-specific scheme used by shapes styled as connectors. Every geometry row, fill pattern, gradient, line pattern and all 45 arrowheads are drawn, including the line jumps Visio computes when routing connectors rather than storing them. Text is laid out by the same DrawingML text engine the PowerPoint converter uses, with fonts embedded as WOFF2 subsets.

`converter/visio` has no `-param` options.

See [design.md §3.8](../design.html#38-visio--bdf-変換器convertervisioの構造) for the internals.

## draw.io

`converter/drawio` reads a diagram's mxGraphModel XML directly rather than going through draw.io's own SVG or PDF export: `.drawio` files (including pages draw.io compresses) and the `.drawio.svg` / `.drawio.png` exports that embed the diagram.

Every page becomes its own page in bdf, so the viewer switches between them with tabs the way a spreadsheet switches sheets. draw.io's layers become the page's layer objects, and links to other pages become `#view=` links. Cell geometry, edge routing (orthogonal, elbow and the other mxGraph edge styles, and shape perimeters), shapes, arrows, stencils, and the wrapping and formatting of HTML labels are ported from draw.io's own (mxGraph's) rendering code, since draw.io never saves the computed edge routes to the file. Stencil libraries (flowchart, basic, arrows, AWS, BPMN, networking, and others) are embedded from draw.io's own stencil XML, and fonts are subset-embedded the same way as for PowerPoint. AWS diagrams are drawn with the current AWS icon set, including diagrams originally drawn with older AWS icon sets, whose shapes are mapped onto their current counterparts. With math typesetting on (`math=1`), LaTeX in labels is laid out as formulas by the same engine PowerPoint and Word formulas use. Hand-drawn styles (`sketch=1`) are drawn as a normal, non-sketchy rendering.

| Option | Values | Default |
|---|---|---|
| `-param border=` | margin around each page's drawing, in pixels | `10` (`0` for none) |

See [design.md §3.11](../design.html#311-drawio--bdf-変換器converterdrawioの構造) for the internals.
