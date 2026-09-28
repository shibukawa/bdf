//go:build !bdf_noconv

package woff2

import (
	"bytes"
	"io"

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

// decompress reads the first size bytes of a Brotli stream: what follows
// the tables of the directory is not read, so that a short stream cannot
// expand to more than the font says it holds. It reads one byte more to
// find the end of a stream of that size, and the damage there may be.
func decompress(b []byte, size int) ([]byte, error) {
	out, err := io.ReadAll(io.LimitReader(brotli.NewReader(bytes.NewReader(b)), int64(size)+1))
	if len(out) > size {
		out = out[:size]
	}
	return out, err
}
