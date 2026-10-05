package bdf

import (
	"bytes"
	"compress/flate"
	"fmt"
	"math/rand"
	"testing"
)

func TestCompressionReuse(t *testing.T) {
	noise := make([]byte, 70000)
	rand.New(rand.NewSource(1)).Read(noise)
	inputs := [][]byte{nil, []byte("a"), bytes.Repeat([]byte("document text and paths "), 4000), noise, []byte("after a large part")}
	for level := flate.HuffmanOnly; level <= flate.BestCompression; level++ {
		t.Run(fmt.Sprint(level), func(t *testing.T) {
			var previous, snapshot []byte
			for _, input := range inputs {
				var want bytes.Buffer
				w, err := flate.NewWriter(&want, level)
				if err != nil {
					t.Fatal(err)
				}
				w.Write(input)
				w.Close()
				got, err := compress(input, level)
				if err != nil || !bytes.Equal(got, want.Bytes()) {
					t.Fatalf("%d bytes: output differs from a fresh compressor: %v", len(input), err)
				}
				if !bytes.Equal(previous, snapshot) {
					t.Fatal("reusing the compressor changed its previous output")
				}
				previous, snapshot = got, bytes.Clone(got)
				// A failed or interrupted read must not affect the next part.
				if _, err := inflate([]byte{0xff}, 100, 0); err == nil {
					t.Fatal("accepted malformed deflate")
				}
				if len(input) > 0 {
					if _, err := inflate(got, len(input)-1, 0); err == nil {
						t.Fatal("accepted data over the limit")
					}
				}
				decoded, err := inflate(got, len(input), len(input))
				if err != nil || !bytes.Equal(decoded, input) {
					t.Fatalf("%d bytes: round trip failed: %v", len(input), err)
				}
			}
		})
	}
	for _, level := range []int{-3, 10} {
		if _, err := compress(noise, level); err == nil {
			t.Errorf("accepted compression level %d", level)
		}
	}
}

func TestConcurrentCompression(t *testing.T) {
	for i := 0; i < 8; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			for n := 0; n < 32; n++ {
				input := bytes.Repeat([]byte{byte(i), byte(n), 0xff}, 1000+n)
				encoded, err := compress(input, 1+n%9)
				if err != nil {
					t.Fatal(err)
				}
				decoded, err := inflate(encoded, len(input), len(input))
				if err != nil || !bytes.Equal(decoded, input) {
					t.Fatalf("round trip failed: %v", err)
				}
			}
		})
	}
}
