//go:build !(js && wasm)

package parquet

import (
	"bytes"

	"github.com/andybalholm/brotli"
)

// brotliDecompress decompresses Brotli. The browser build leaves it out
// (brotli_js.go): the package's tables add 1.7 MB to the module, for a
// codec few files use.
func brotliDecompress(src []byte, size int) ([]byte, error) {
	return readSized(brotli.NewReader(bytes.NewReader(src)), size)
}
