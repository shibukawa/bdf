//go:build !bdf_noconv

package woff2

import (
	"bytes"
	"testing"
)

// TestEncodeIsDeterministic checks that a font encodes to the same bytes
// however many fonts were encoded before it: the Brotli writers are reused,
// and a reused one must not remember the font before.
func TestEncodeIsDeterministic(t *testing.T) {
	fonts := testFonts(t)
	first := map[string][]byte{}
	for round := range 3 {
		for name, font := range fonts {
			out, err := Encode(font)
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			if round == 0 {
				first[name] = out
			} else if !bytes.Equal(out, first[name]) {
				t.Fatalf("%s: round %d encodes to other bytes than the first", name, round)
			}
		}
	}
}
