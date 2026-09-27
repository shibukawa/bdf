//go:build js && wasm && !pdfonly && !officeonly && !imageonly

package main

import (
	_ "github.com/shibukawa/bdf/converter/epub"
	_ "github.com/shibukawa/bdf/converter/html"
	_ "github.com/shibukawa/bdf/converter/markdown"
)
