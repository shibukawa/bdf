//go:build !bdf_noconv

package woff2

import (
	"bytes"
	"testing"
)

func TestFontCompressionWindowBoundaries(t *testing.T) {
	for _, n := range []int{0, 1, 1008, 1009, (1 << 16) - 16, (1 << 16) - 15, 1 << 20, (1 << 20) + 1} {
		input := make([]byte, n)
		for i := range input {
			input[i] = byte(i*37 + i/251)
		}
		packed, err := compress(input)
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := decompress(packed, len(input))
		if err != nil || !bytes.Equal(decoded, input) {
			t.Fatalf("%d-byte input: round trip failed: %v", n, err)
		}
	}
}
