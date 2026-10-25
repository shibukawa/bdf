//go:build !bdf_noconv

package woff2

import (
	"bytes"
	"fmt"
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

// TestDecompressPooled reads streams one after another and at once through
// the pooled readers: each read gets its own bytes.
func TestDecompressPooled(t *testing.T) {
	inputs := make([][]byte, 6)
	packed := make([][]byte, len(inputs))
	for i := range inputs {
		inputs[i] = bytes.Repeat([]byte{byte('a' + i)}, 1000*(i+1))
		var err error
		if packed[i], err = compress(inputs[i]); err != nil {
			t.Fatal(err)
		}
	}
	for round := 0; round < 3; round++ {
		for i := range inputs {
			got, err := decompress(packed[i], len(inputs[i]))
			if err != nil || !bytes.Equal(got, inputs[i]) {
				t.Fatalf("round %d, stream %d: %v", round, i, err)
			}
		}
	}
	done := make(chan error, 4*len(inputs))
	for g := 0; g < 4; g++ {
		go func() {
			for i := range inputs {
				got, err := decompress(packed[i], len(inputs[i]))
				if err == nil && !bytes.Equal(got, inputs[i]) {
					err = fmt.Errorf("stream %d: wrong bytes", i)
				}
				done <- err
			}
		}()
	}
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			t.Fatal(err)
		}
	}
}
