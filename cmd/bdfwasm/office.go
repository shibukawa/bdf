//go:build js && wasm && !pdfonly

package main

import (
	_ "github.com/shibukawa/bdf/converter/csv"
	_ "github.com/shibukawa/bdf/converter/docx"
	_ "github.com/shibukawa/bdf/converter/drawio"
	_ "github.com/shibukawa/bdf/converter/dxf"
	_ "github.com/shibukawa/bdf/converter/emf"
	_ "github.com/shibukawa/bdf/converter/pptx"
	_ "github.com/shibukawa/bdf/converter/psd"
	_ "github.com/shibukawa/bdf/converter/visio"
	_ "github.com/shibukawa/bdf/converter/xlsx"
)
