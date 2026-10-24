//go:build !(js && wasm)

package parquet

import (
	"bytes"

	brrr "github.com/molecule-man/go-brrr"
)

// brotliDecompress decompresses Brotli. The browser build leaves it out
// (brotli_js.go): the decoder adds 0.45 MB to the module, for a codec few
// files use.
func brotliDecompress(src []byte, size int) ([]byte, error) {
	r := brrr.NewReader(bytes.NewReader(src))
	defer r.Close()
	return readSized(r, size)
}
