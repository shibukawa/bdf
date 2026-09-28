# Why bdf

## No office suite to run

Previewing an Office file in a browser usually means running LibreOffice or OpenOffice headless on a server and converting it to PDF. That install is over a gigabyte before it does anything, and a headless office suite needs a process to start, keep alive, and isolate — one malformed file can wedge it.

bdf's converters are Go packages that use no cgo and no external program: PDF, Word, PowerPoint, Excel, CSV, Parquet, Visio, draw.io, DXF, Jw_cad, SXF, CGM, HP-GL/2, Gerber, Excellon, KiCad, TIFF, Illustrator, Photoshop, Windows metafiles, HTML, Markdown, EPUB, MML, MIDI, MusicXML and font files all convert in one binary. Built as WebAssembly, the same code converts inside the browser too — gzipped, about 7.3 MB for the PDF module and about 7.8 MB for the Office/CAD/KiCad/font module — so a file never has to leave the browser at all. Thumbnails and search text come from the same process. Pages are drawn by a pure-Go rasterizer; no browser is needed on the server either.

## In the shape of the content

A PDF cuts everything into sheets of paper. Print a spreadsheet and a wide table is sliced across page boundaries, so you lose the row you were following and can't find the cell you wanted; frozen headers and gridlines disappear, and the sheets you switched between in the workbook become one run of pages. A Word document is stuck reading page by page too.

bdf has a layout model for each kind of content. Slides, drawings and PDFs get fixed-size pages. Worksheets get an endless plane per sheet, drawn in tiles, with frozen panes, row and column headers, and gridlines. Word processor documents get a flow you can read as pages or as one continuous scroll, and a page-free view re-laid-out as one long column. Illustrator and Photoshop artboards become pages. A workbook's sheets, a draw.io diagram's pages, a DXF drawing's model space and layouts, and a circuit board's front, back and each layer each become one view, switched with tabs in the viewer.

## Made for browsers

The instruction set maps one to one onto Canvas 2D. Fonts are WOFF2 for `FontFace`, images are in whatever format the browser can decode itself, and a Part's compression is what `DecompressionStream` can inflate. Font rasterization, image decoding and decompression — the tens of thousands of lines a PDF viewer like pdf.js reimplements — are left to the browser, so bdf's decoder and renderer are about 3,400 lines of TypeScript (the renderer's worker is 21 KB gzipped). Drawing happens on an `OffscreenCanvas` in a Worker; the main thread only places the bitmap. Parts are content-addressed, so a master or a repeated element is stored once, and the viewer fetches only the Parts a visible page needs — by range request, or from a split layout on a CDN. A PDF converted in the browser is drawn page by page as it converts, the visible page first.

## Many formats, one renderer

PDF, Word (.docx), PowerPoint (.pptx), Excel (.xlsx), CSV/TSV, Apache Parquet, Visio (.vsdx, .vdx), draw.io (.drawio, and SVG/PNG exports with the diagram embedded), AutoCAD DXF, Jw_cad (.jww), SXF (.p21, .p2z, .sfc), CGM (.cgm), HP-GL/2 plot files (.plt), Gerber (RS-274X) and Excellon drill files (one by one, or a board's files zipped together), KiCad schematics and boards (.kicad_sch, .kicad_pcb, project .kicad_pro, or a project zipped), TIFF, Illustrator (.ai), Photoshop (.psd, .psb), Windows metafiles (.emf, .wmf), HTML (reader mode), Markdown, EPUB (reflowable books with Japanese vertical text, fixed-layout manga), scores engraved from MML (.mml), MIDI (.mid, .kar) and MusicXML (.musicxml, .mxl), and font files (.ttf, .otf, .ttc, .woff, .woff2) showing their characters, glyphs and OpenType features — all of it, including password-protected Office documents and PDFs, becomes the same format. Every format is drawn by the same renderer, and search, text selection and the accessible text layer (headings, lists, tables, alt text) work the same way across all of them.

See [Formats](formats/index.html) for what each converter keeps and how it lays its input out, and [Architecture](architecture/index.html) for how conversion and rendering fit together.
