# Circuit boards, KiCad

Two converters read printed-circuit-board data: `converter/gerber` for the fabrication files a board house receives (Gerber and Excellon), and `converter/kicad` for the KiCad project the board was designed in. Both draw the board as it actually looks — substrate, copper, mask, silkscreen, holes — rather than as an abstract schematic of layers.

## Try it

Drop a file on the [viewer](../../viewer/), or try a sample: a fabricated [circuit board (Gerber and Excellon, zipped)](../../viewer/?file=samples/board.zip), or a [KiCad project (schematic and board, zipped)](../../viewer/?file=samples/demo.zip).

## Gerber, Excellon (converter/gerber)

`converter/gerber` reads the fabrication data of a printed circuit board: Gerber files (RS-274X with X2 attributes, plus the deprecated commands older files still use) and Excellon drill files. RS-274-D, which keeps aperture definitions in a separate file, isn't read. A board's files are normally handed to a manufacturer as a zip archive (with its Gerber job file); bdf converts that zip as one board. A single Gerber or drill file converts too, becoming a view of just that layer.

Every standard aperture and macro primitive is drawn, arcs in both quadrant modes, regions with clear polarity, step-and-repeat, block apertures and aperture transformations. Pads are drawn as runs of shared paths, the same way glyphs are. Which layer each file represents (top/bottom copper, an inner layer, solder mask, silkscreen, paste, outline, drill) comes from the job file, from X2 `TF.FileFunction` attributes, or — failing that — from the filename conventions of KiCad, Altium (Protel), Eagle, EasyEDA and others.

When there are two or more files and at least one is copper, mask or silkscreen, the board gets a `top` view and a `bottom` view (seen from below, mirrored) drawn the way the physical board looks: substrate, copper under the mask, the finish on pads the mask leaves bare, silkscreen, and holes — all clipped to the board's outline. A view per file follows, drawn in the dark-background colors PCB design tools use.

| Option | Values | Default |
|---|---|---|
| `views` | `all`, `board`, `layers` | `all` — the board's top and bottom, plus a view per file |
| `mask` | `green`, `red`, `blue`, `black`, `white`, `yellow`, `purple`, or `#rrggbb` | the job file's mask color, else `green` |
| `silkscreen` | `white`, `black`, `yellow`, or `#rrggbb` | the job file's, else `white` |
| `finish` | `gold`, `silver`, `copper` | the job file's finish, else `gold` |

See [design.md §3.21](../design.html#321-gerberexcellon--bdf-変換器convertergerberの構造) for the internals.

## KiCad (converter/kicad)

`converter/kicad` reads the schematics (`.kicad_sch`) and boards (`.kicad_pcb`) of KiCad 6 and later — KiCad 5's file format is refused, with a note to save it again in a newer KiCad. The input can be one schematic or board file, a project file (`.kicad_pro`), or a whole project zipped up.

A schematic becomes a view whose pages are its sheet instances, in the order of their page numbers: a hierarchical sheet used twice in the design becomes two pages, each keeping its own reference designators (R101, R201, and so on), and a sheet's box on the page links to that page. Every page is drawn on its paper in its drawing sheet — KiCad's default title block, or the project's `.kicad_wks`, with the title block and the project's text variables filled in — using KiCad's default colors and the project's drawing settings. Symbols draw with their unit, alternate body style and pin shape; wires, buses, junctions, every kind of label, notes, text boxes, tables, hatched shapes, images and do-not-place (DNP) marks all draw as KiCad draws them. Text uses KiCad's own stroke font, NewStroke (a CC0 build is embedded; CJK characters the font lacks fall back to a TrueType font).

A board becomes a `front` view and a `back` view (mirrored, seen from below) with every layer drawn the way KiCad's board editor draws it, plus a view for each individual layer. The sheets a hierarchical schematic references, and the drawing sheet a project names, are separate files: bdf reads them from the input's own directory, from a zip archive, or from files a server lists alongside the input (`converter.Options.Files`; `bdf generate -with`).

| Option | Values | Default |
|---|---|---|
| `views` | `all`, `schematic`, `board` | `all` — a project's schematic and board |
| `layers` | `all`, `board`, `layers` | `all` — a board's front, back, and a view per layer |

`-with` adds files a schematic refers to (the sheets of a hierarchical design, or the project file) by their path relative to the input, when the server or CLI invocation doesn't have them alongside it already.

See [design.md §3.29](../design.html#329-kicad--bdf-変換器converterkicadの構造) for the internals.
