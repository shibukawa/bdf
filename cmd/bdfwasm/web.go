//go:build js && wasm && !pdfonly && !officeonly && !imageonly && !previewonly && !npm_pdfepub && !npm_office && !npm_custom

package main

import (
	_ "github.com/shibukawa/bdf/converter/epub"
	_ "github.com/shibukawa/bdf/converter/html"
	_ "github.com/shibukawa/bdf/converter/markdown"
)
