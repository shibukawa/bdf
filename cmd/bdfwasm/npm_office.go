//go:build js && wasm && npm_office

package main

import (
	_ "github.com/shibukawa/bdf/converter/pdf"
	_ "github.com/shibukawa/bdf/converter/epub"
	_ "github.com/shibukawa/bdf/converter/docx"
	_ "github.com/shibukawa/bdf/converter/xlsx"
	_ "github.com/shibukawa/bdf/converter/pptx"
	_ "github.com/shibukawa/bdf/converter/visio"
	_ "github.com/shibukawa/bdf/converter/idml"
	_ "github.com/shibukawa/bdf/converter/csv"
	_ "github.com/shibukawa/bdf/converter/parquet"
)
