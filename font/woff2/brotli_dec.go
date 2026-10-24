//go:build !bdf_noconv || bdf_woff2dec

package woff2

import (
	"bytes"
	"io"
	"sync"

	brrr "github.com/molecule-man/go-brrr"
)

// The readers are pooled: a new one is 170 KiB of state and buffers, a
// Reset of a pooled one a few KiB. Its ring buffer, the window of the last
// stream it read (4 MiB at most for a font this package wrote), stays with
// it in the pool until a collection.
var brotliReaders sync.Pool

// brotliReader is a pooled Brotli reader with the reader it reads from.
type brotliReader struct {
	in bytes.Reader
	r  *brrr.Reader
}

// decompress reads the first size bytes of a Brotli stream: what follows
// the tables of the directory is not read, so that a short stream cannot
// expand to more than the font says it holds. It reads one byte more to
// find the end of a stream of that size, and the damage there may be.
// (go-brrr 1.3.0 gives out what a command past the end of its metablock
// wrote before it fails, molecule-man/go-brrr#71: such a stream is read
// as tables when that error comes after size bytes.)
func decompress(b []byte, size int) ([]byte, error) {
	rd, _ := brotliReaders.Get().(*brotliReader)
	if rd == nil {
		rd = &brotliReader{}
		rd.r = brrr.NewReader(&rd.in)
	}
	rd.in.Reset(b)
	rd.r.Reset(&rd.in)
	defer func() {
		rd.in.Reset(nil) // not the font's bytes, while pooled
		rd.r.Reset(&rd.in)
		brotliReaders.Put(rd)
	}()
	out, err := io.ReadAll(io.LimitReader(rd.r, int64(size)+1))
	if len(out) > size {
		out = out[:size]
	}
	return out, err
}
