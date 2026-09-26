//go:build js && wasm && !pdfonly

package main

// Images are in the Office module too, and alone in the image module
// (-tags imageonly), which is small.
import _ "github.com/shibukawa/bdf/converter/image"
