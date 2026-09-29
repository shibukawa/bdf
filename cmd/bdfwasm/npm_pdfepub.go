//go:build js && wasm && npm_pdfepub

package main

import (
	_ "github.com/shibukawa/bdf/converter/pdf"
	_ "github.com/shibukawa/bdf/converter/epub"
)
