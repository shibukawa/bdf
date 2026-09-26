// Package all registers every input format of this module with the
// converter registry. Import it for its side effect:
//
//	import _ "github.com/shibukawa/bdf/converter/all"
//
// Programs that need only some formats import those converter packages
// instead, and link only them.
package all

import (
	_ "github.com/shibukawa/bdf/converter/csv"   // CSV and TSV
	_ "github.com/shibukawa/bdf/converter/emf"   // Windows metafiles
	_ "github.com/shibukawa/bdf/converter/pdf"   // PDF
	_ "github.com/shibukawa/bdf/converter/pptx"  // PowerPoint
	_ "github.com/shibukawa/bdf/converter/visio" // Visio
	_ "github.com/shibukawa/bdf/converter/xlsx"  // Excel
)
