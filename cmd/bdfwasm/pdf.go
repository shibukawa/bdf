//go:build js && wasm && !officeonly && !webonly && !imageonly && !previewonly && !npm_pdfepub && !npm_office && !npm_custom

package main

import (
	_ "github.com/shibukawa/bdf/converter/ai" // Illustrator: a PDF, so it goes with the PDF converter
	_ "github.com/shibukawa/bdf/converter/pdf"
)
