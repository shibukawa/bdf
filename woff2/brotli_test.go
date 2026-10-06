//go:build !bdf_noconv

package woff2

import (
	"bytes"
	"math/rand"
	"runtime"
	"testing"
)

// Small embedded subsets must not allocate a multi-megabyte match tree.
func TestSmallFontCompressionMemory(t *testing.T) {
	input := make([]byte, 4096)
	rand.New(rand.NewSource(1)).Read(input)
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	packed, err := compress(input)
	runtime.ReadMemStats(&after)
	if err != nil {
		t.Fatal(err)
	}
	if n := after.TotalAlloc - before.TotalAlloc; n > 8<<20 {
		t.Fatalf("compressing 4 KiB allocated %d bytes", n)
	}
	decoded, err := decompress(packed, len(input))
	if err != nil || !bytes.Equal(decoded, input) {
		t.Fatalf("round trip: %v", err)
	}
}

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
