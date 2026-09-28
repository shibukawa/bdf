# Formats

bdf converts each input from its own package under `converter/`; a program only handles the formats it links.

```go
import (
	"github.com/shibukawa/bdf/converter"
	_ "github.com/shibukawa/bdf/converter/pdf" // every format: converter/all
)

res, err := converter.ConvertFile("in.pdf", "", &converter.Options{}) // "" detects the format
```

The format is detected from the input's content where that's possible; Markdown, which can look like any text, is picked from the `.md`/`.markdown` extension. A password-protected input opens with `Options.Password`; when `res.Protected` says the input needed one, the output is encrypted with the same password (see [Features → Encrypted input](../features.html#encrypted-input)).

## By category

| Category | Formats | Page |
|---|---|---|
| Documents | PDF, Illustrator (.ai) | [PDF, Illustrator](pdf.html) |
| Word processing | Word (.docx), HTML, Markdown | [Word, HTML, Markdown](document.html) |
| Presentations | PowerPoint (.pptx) | [PowerPoint](presentation.html) |
| Spreadsheets | Excel (.xlsx), CSV/TSV, Apache Parquet | [Excel, CSV, Parquet](spreadsheet.html) |
| Diagrams | Visio (.vsdx, .vdx), draw.io | [Visio, draw.io](diagram.html) |
| CAD drawings and plots | AutoCAD DXF, Jw_cad (.jww), SXF, CGM, HP-GL/2 | [CAD drawings and plots](cad.html) |
| Electronics | Gerber/Excellon (circuit boards), KiCad | [Circuit boards, KiCad](electronics.html) |
| Books | EPUB | [EPUB](ebook.html) |
| Music | MML, MIDI, MusicXML | [Scores](music.html) |
| Fonts | TrueType/OpenType/WOFF/WOFF2 | [Font files](font.html) |
| Images and design | PNG/JPEG/GIF/WebP/AVIF/BMP/ICO/SVG, Photoshop, TIFF, Windows metafiles | [Images, Photoshop](image.html) |

## Every input format

| Format | Package | Extensions |
|---|---|---|
| PDF | `converter/pdf` | .pdf |
| Illustrator | `converter/ai` | .ai |
| Photoshop | `converter/psd` | .psd, .psb |
| PowerPoint | `converter/pptx` | .pptx, .pptm, .ppsx, .ppsm, .potx, .potm |
| Excel | `converter/xlsx` | .xlsx, .xlsm, .xltx, .xltm |
| CSV / TSV | `converter/csv` | .csv, .tsv, .tab |
| Apache Parquet | `converter/parquet` | .parquet, .parq, .pqt |
| Word | `converter/docx` | .docx, .docm, .dotx, .dotm |
| Visio | `converter/visio` | .vsdx, .vsdm, .vstx, .vstm, .vdx, .vtx |
| draw.io | `converter/drawio` | .drawio, .dio, .drawio.svg, .drawio.png |
| AutoCAD DXF | `converter/dxf` | .dxf |
| Jw_cad | `converter/jww` | .jww |
| SXF | `converter/sxf` | .p21, .p2z, .sfc |
| CGM | `converter/cgm` | .cgm, .cgz |
| HP-GL/2 | `converter/hpgl` | .plt, .hpgl, .hgl, .hpg, .hp2, .rtl |
| Gerber / Excellon | `converter/gerber` | .gbr and related, .drl, .xln, or a .zip of a board |
| KiCad | `converter/kicad` | .kicad_sch, .kicad_pcb, .kicad_pro, or a project .zip |
| Windows metafile | `converter/emf` | .emf, .wmf |
| TIFF | `converter/tiff` | .tif, .tiff |
| HTML | `converter/html` | .html, .xhtml, .mhtml |
| Markdown | `converter/markdown` | .md, .markdown, .mdown, .mkd, .mdx |
| EPUB | `converter/epub` | .epub |
| MML | `converter/mml` | .mml |
| MIDI | `converter/midi` | .mid, .midi, .smf, .kar, .rmi |
| MusicXML | `converter/musicxml` | .musicxml, .mxl |
| Font file | `converter/font` | .ttf, .otf, .ttc, .otc, .woff, .woff2 |
| Image | `converter/image` | .png, .jpg, .gif, .webp, .avif, .bmp, .ico, .svg |

Each package's `-param` options are listed by `go run ./cmd/bdf generate -h`, and again on the page for that format. For the container format each converter writes into, see the [format specification](../spec.html); for why each converter is built the way it is, [design notes](../design.html).
