//go:build !(js && wasm)

package parquet

import (
	"bytes"
	"sync"

	brrr "github.com/molecule-man/go-brrr"
)

// brotliDecompress decompresses Brotli. The browser build leaves it out
// (brotli_js.go): the decoder adds 0.45 MB to the module, for a codec few
// files use. The readers are pooled like the Zstandard ones: a new one is
// 170 KiB of state and buffers, a Reset a few KiB.
func brotliDecompress(src []byte, size int) ([]byte, error) {
	rd, _ := brotliReaders.Get().(*brotliReader)
	if rd == nil {
		rd = &brotliReader{}
		rd.r = brrr.NewReader(&rd.in)
	}
	rd.in.Reset(src)
	rd.r.Reset(&rd.in)
	defer func() {
		// Do not keep the compressed page alive while the reader is pooled.
		rd.in.Reset(nil)
		rd.r.Reset(&rd.in)
		brotliReaders.Put(rd)
	}()
	return readSized(rd.r, size)
}

var brotliReaders sync.Pool

// brotliReader is a pooled Brotli reader with the reader it reads from.
type brotliReader struct {
	in bytes.Reader
	r  *brrr.Reader
}
