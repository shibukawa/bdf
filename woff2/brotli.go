//go:build !bdf_noconv

package woff2

import (
	"bytes"
	"io"
	"math/bits"

	"github.com/andybalholm/brotli"
)

// Available reports whether the Brotli encoder is compiled in.
func Available() bool { return true }

func compress(b []byte) ([]byte, error) {
	// Quality 9 keeps font data lossless while avoiding quality 11's
	// expensive match search for a few percent smaller embedded subsets.
	const quality = 9
	var out bytes.Buffer
	// A Brotli window holds (1<<LGWin)-16 bytes. A font is passed in one
	// write, so a larger window than its table stream cannot find more
	// matches. Quality 11 allocates eight bytes per window position for
	// its match tree; always using LGWin=22 costs 32 MiB per small subset.
	lgwin := 22
	if len(b) < (1<<22)-16 {
		lgwin = max(10, bits.Len(uint(len(b)+15)))
	}
	w := brotli.NewWriterOptions(&out, brotli.WriterOptions{Quality: quality, LGWin: lgwin})
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
