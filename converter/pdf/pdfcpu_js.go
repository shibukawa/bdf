//go:build js

package pdf

import "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

// A browser has no configuration directory for pdfcpu to keep its
// config.yml in, and pdfcpu's js parser of that file rejects the
// hexadecimal permissions it writes there itself: use its built-in
// defaults instead.
func init() { model.ConfigPath = "disable" }
