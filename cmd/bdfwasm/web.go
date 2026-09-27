//go:build js && wasm && !pdfonly && !officeonly

package main

import (
	_ "github.com/shibukawa/bdf/converter/html"
	_ "github.com/shibukawa/bdf/converter/markdown"
)
