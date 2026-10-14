//go:build !bdf_noconv || bdf_woff2dec

package woff2

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"
)

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
