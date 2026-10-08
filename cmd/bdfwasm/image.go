//go:build js && wasm && !pdfonly && !webonly && !previewonly && !npm_pdfepub && !npm_office && !npm_custom

package main

// Images and audio files are in the Office module too, and alone in the
// image module (-tags imageonly), which is small: both store what the
// browser decodes and read only the metadata.
import (
	_ "github.com/shibukawa/bdf/converter/audio"
	_ "github.com/shibukawa/bdf/converter/image"
)
