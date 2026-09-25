//go:build !bdf_noconv

package woff2

import (
	"bytes"

	"github.com/andybalholm/brotli"
)

// Available reports whether the Brotli encoder is compiled in.
func Available() bool { return true }

func compress(b []byte) ([]byte, error) {
	// Quality 11 is 20–30 times slower than 9 for a few percent less; keep it
	// for the usual subset font and use 9 for whole multi-megabyte fonts.
	quality := 11
	if len(b) > 1<<20 {
		quality = 9
	}
	var out bytes.Buffer
	w := brotli.NewWriterOptions(&out, brotli.WriterOptions{Quality: quality, LGWin: 22})
	if _, err := w.Write(b); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
