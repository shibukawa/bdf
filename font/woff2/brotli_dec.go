//go:build !bdf_noconv || bdf_woff2dec

package woff2

import (
	"bytes"
	"io"

	brrr "github.com/molecule-man/go-brrr"
)

// decompress reads the first size bytes of a Brotli stream: what follows
// the tables of the directory is not read, so that a short stream cannot
// expand to more than the font says it holds. It reads one byte more to
// find the end of a stream of that size, and the damage there may be.
// (go-brrr 1.3.0 gives out what a command past the end of its metablock
// wrote before it fails, molecule-man/go-brrr#71: such a stream is read
// as tables when that error comes after size bytes.)
func decompress(b []byte, size int) ([]byte, error) {
	r := brrr.NewReader(bytes.NewReader(b))
	defer r.Close()
	out, err := io.ReadAll(io.LimitReader(r, int64(size)+1))
	if len(out) > size {
		out = out[:size]
	}
	return out, err
}
