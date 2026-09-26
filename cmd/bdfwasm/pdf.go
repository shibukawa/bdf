//go:build js && wasm && !officeonly

package main

import (
	_ "github.com/shibukawa/bdf/converter/ai" // Illustrator: a PDF, so it goes with the PDF converter
	_ "github.com/shibukawa/bdf/converter/pdf"
)
