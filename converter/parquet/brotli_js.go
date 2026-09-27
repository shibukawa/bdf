//go:build js && wasm

package parquet

import "errors"

func brotliDecompress([]byte, int) ([]byte, error) {
	return nil, errors.New("not supported in the browser build")
}
