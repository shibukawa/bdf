package bdf

import (
	"bytes"
	stdflate "compress/flate"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/klauspost/compress/flate"
)

// compressedParts returns the parts of the test documents that the writer
// compresses.
func compressedParts(t testing.TB, patterns ...string) [][]byte {
	t.Helper()
	var parts [][]byte
	for _, pattern := range patterns {
		files, _ := filepath.Glob(pattern)
		for _, name := range files {
			r, err := OpenSingleFile(name)
			if err != nil || r.Locked() {
				continue
			}
			d, err := r.ToDocument()
			r.Close()
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			for _, p := range d.Parts() {
				if d.shouldCompress(p) {
					parts = append(parts, p.Data)
				}
			}
		}
	}
	if len(parts) == 0 {
		t.Fatal("no parts")
	}
	return parts
}

// TestFlateMatchesStandardLibrary checks that compress, which writes the
// parts with klauspost/compress, writes the stream the standard library
// writes at every level: the documents in testdata were written with the
// standard library, and a document is to be the same bytes whichever Go
// wrote it.
func TestFlateMatchesStandardLibrary(t *testing.T) {
	parts := compressedParts(t, "testdata/*.bdf", "testdata/*/*.bdf")
	for _, level := range []int{flate.BestCompression, flate.DefaultCompression, 7, 6, flate.BestSpeed, flate.HuffmanOnly} {
		var a bytes.Buffer
		w, err := stdflate.NewWriter(&a, level)
		if err != nil {
			t.Fatal(err)
		}
		for i, p := range parts {
			a.Reset()
			w.Reset(&a)
			w.Write(p)
			w.Close()
			b, err := compress(p, level)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(a.Bytes(), b) {
				t.Fatalf("level %d: part %d of %d bytes compresses to other bytes (%d vs %d)", level, i, len(p), len(a.Bytes()), len(b))
			}
		}
	}
}

// TestFlateCandidates compares the deflate of the standard library with
// klauspost/compress on the parts the writer compresses: time and size
// (BDF_FLATE_BENCH=1 go test -run TestFlateCandidates -v).
func TestFlateCandidates(t *testing.T) {
	if os.Getenv("BDF_FLATE_BENCH") == "" {
		t.Skip("BDF_FLATE_BENCH is not set")
	}
	parts := compressedParts(t, "testdata/*.bdf", "testdata/*/*.bdf")
	total := 0
	for _, p := range parts {
		total += len(p)
	}
	t.Logf("%d parts, %.1f MB raw", len(parts), float64(total)/1e6)
	type cfg struct {
		name string
		enc  func(dst *bytes.Buffer, src []byte)
	}
	std := func(level int) func(*bytes.Buffer, []byte) {
		w, _ := stdflate.NewWriter(io.Discard, level)
		return func(dst *bytes.Buffer, src []byte) { w.Reset(dst); w.Write(src); w.Close() }
	}
	kp := func(level int) func(*bytes.Buffer, []byte) {
		w, _ := flate.NewWriter(io.Discard, level)
		return func(dst *bytes.Buffer, src []byte) { w.Reset(dst); w.Write(src); w.Close() }
	}
	cfgs := []cfg{{"std 9", std(9)}, {"std 6", std(6)}, {"klauspost 9", kp(9)}, {"klauspost 7", kp(7)}, {"klauspost 6", kp(6)}, {"klauspost 5", kp(5)}}
	const rounds = 5
	var out bytes.Buffer
	var compressed [][]byte
	for i, c := range cfgs {
		size := 0
		start := time.Now()
		for round := range rounds {
			for _, p := range parts {
				out.Reset()
				c.enc(&out, p)
				if round == 0 {
					size += out.Len()
					if i == 0 {
						compressed = append(compressed, append([]byte(nil), out.Bytes()...))
					}
				}
			}
		}
		el := time.Since(start) / rounds
		t.Logf("%-14s %7.1f ms  %6.1f MB/s  %8d bytes (%.1f%% of raw)", c.name, float64(el.Microseconds())/1000, float64(total)/el.Seconds()/1e6, size, 100*float64(size)/float64(total))
	}
	for _, d := range []struct {
		name string
		run  func(src []byte) int
	}{
		{"std inflate", func(src []byte) int {
			n, _ := io.Copy(io.Discard, stdflate.NewReader(bytes.NewReader(src)))
			return int(n)
		}},
		{"klauspost inflate", func(src []byte) int {
			n, _ := io.Copy(io.Discard, flate.NewReader(bytes.NewReader(src)))
			return int(n)
		}},
	} {
		start := time.Now()
		n := 0
		for range rounds {
			n = 0
			for _, c := range compressed {
				n += d.run(c)
			}
		}
		el := time.Since(start) / rounds
		if n != total {
			t.Fatalf("%s: %d bytes, want %d", d.name, n, total)
		}
		t.Logf("%-18s %7.1f ms  %6.1f MB/s (decoded)", d.name, float64(el.Microseconds())/1000, float64(total)/el.Seconds()/1e6)
	}
}
