// Package all registers every input format of this module with the
// converter registry. Import it for its side effect:
//
//	import _ "github.com/shibukawa/bdf/converter/all"
//
// Programs that need only some formats import those converter packages
// instead, and link only them.
package all

import (
	_ "github.com/shibukawa/bdf/converter/ai"       // Illustrator
	_ "github.com/shibukawa/bdf/converter/cgm"      // CGM
	_ "github.com/shibukawa/bdf/converter/csv"      // CSV and TSV
	_ "github.com/shibukawa/bdf/converter/docx"     // Word
	_ "github.com/shibukawa/bdf/converter/drawio"   // draw.io
	_ "github.com/shibukawa/bdf/converter/dxf"      // AutoCAD DXF
	_ "github.com/shibukawa/bdf/converter/emf"      // Windows metafiles
	_ "github.com/shibukawa/bdf/converter/epub"     // EPUB, in reader mode
	_ "github.com/shibukawa/bdf/converter/gerber"   // Gerber and Excellon (PCB fabrication data)
	_ "github.com/shibukawa/bdf/converter/hpgl"     // HP-GL/2 plot files
	_ "github.com/shibukawa/bdf/converter/html"     // HTML, in reader mode
	_ "github.com/shibukawa/bdf/converter/image"    // images browsers display (PNG, JPEG, SVG …)
	_ "github.com/shibukawa/bdf/converter/jww"      // Jw_cad
	_ "github.com/shibukawa/bdf/converter/markdown" // Markdown, in reader mode
	_ "github.com/shibukawa/bdf/converter/parquet"  // Apache Parquet
	_ "github.com/shibukawa/bdf/converter/pdf"      // PDF
	_ "github.com/shibukawa/bdf/converter/pptx"     // PowerPoint
	_ "github.com/shibukawa/bdf/converter/psd"      // Photoshop
	_ "github.com/shibukawa/bdf/converter/sxf"      // SXF (P21, SFC)
	_ "github.com/shibukawa/bdf/converter/tiff"     // TIFF images
	_ "github.com/shibukawa/bdf/converter/visio"    // Visio
	_ "github.com/shibukawa/bdf/converter/xlsx"     // Excel
)
