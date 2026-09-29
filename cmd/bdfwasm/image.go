//go:build js && wasm && !pdfonly && !webonly && !previewonly && !npm_pdfepub && !npm_office && !npm_custom

package main

// Images are in the Office module too, and alone in the image module
// (-tags imageonly), which is small.
import _ "github.com/shibukawa/bdf/converter/image"
