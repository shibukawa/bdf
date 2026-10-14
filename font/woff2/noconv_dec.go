//go:build bdf_noconv && !bdf_woff2dec

package woff2

// (the decoder of Brotli is left out with the encoder, unless the build
// asks for it by itself with the tag bdf_woff2dec: a program that reads
// the fonts of documents made elsewhere, and makes none)
func decompress([]byte, int) ([]byte, error) { return nil, ErrDecodeNotAvailable }
