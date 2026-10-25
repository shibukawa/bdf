package pdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"math/rand"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/filter"
)

// zlibOf compresses b with the zlib wrapper that FlateDecode data has.
func zlibOf(b []byte) []byte {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(b)
	w.Close()
	return buf.Bytes()
}

// pdfcpuInflate decodes as pdfcpu's filter did, the oracle for the rows a
// predictor reconstructs.
func pdfcpuInflate(t *testing.T, data []byte, parms map[string]int, max int64) []byte {
	t.Helper()
	fi, err := filter.NewFilter(filter.Flate, parms)
	if err != nil {
		t.Fatal(err)
	}
	r, err := fi.DecodeLength(bytes.NewReader(data), max+1)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	return r.(*bytes.Buffer).Bytes()
}

// TestPredictorsAgree reconstructs rows of every predictor, depth and
// width from random sample data, as pdfcpu did.
func TestPredictorsAgree(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	for _, c := range []struct{ predictor, colors, bpc, columns int }{
		{2, 1, 8, 1}, {2, 3, 8, 17}, {2, 4, 16, 9}, {2, 1, 1, 13}, {2, 2, 2, 7}, {2, 3, 4, 5},
		{10, 1, 8, 1}, {11, 3, 8, 16}, {12, 4, 8, 7}, {13, 1, 8, 100}, {14, 3, 16, 11}, {15, 1, 1, 37}, {15, 2, 4, 6},
	} {
		name := fmt.Sprintf("predictor %d, %d colours, %d bits, %d columns", c.predictor, c.colors, c.bpc, c.columns)
		parms := map[string]int{"Predictor": c.predictor, "Colors": c.colors, "BitsPerComponent": c.bpc, "Columns": c.columns}
		rowSize := (c.bpc*c.colors*c.columns + 7) / 8
		rows := 1 + rng.Intn(40)
		var raw []byte
		for r := 0; r < rows; r++ {
			if c.predictor >= 10 {
				raw = append(raw, byte(rng.Intn(5))) // a filter type a row
			}
			for i := 0; i < rowSize; i++ {
				raw = append(raw, byte(rng.Intn(256)))
			}
		}
		data := zlibOf(raw)
		for _, max := range []int64{1 << 20, int64(rows*rowSize) - 1, int64(rowSize) / 2} {
			if err := checkPredictor(parms, len(data)); err != nil {
				t.Fatal(name, err)
			}
			got, err := inflateUpTo(data, parms, max, false)
			if err != nil {
				t.Fatalf("%s, max %d: %v", name, max, err)
			}
			want := pdfcpuInflate(t, data, parms, max)
			if !bytes.Equal(got, want) {
				t.Errorf("%s, max %d: %d bytes, want %d (%x…, want %x…)", name, max, len(got), len(want), got[:min(8, len(got))], want[:min(8, len(want))])
			}
		}
	}
}

func TestInflateUpTo(t *testing.T) {
	text := bytes.Repeat([]byte("0 0 10 10 re f\n"), 1000)
	z := zlibOf(text)
	none := map[string]int{}
	check := func(name string, got []byte, err error, want []byte, wantErr bool) {
		t.Helper()
		if (err != nil) != wantErr {
			t.Errorf("%s: error %v, want error %v", name, err, wantErr)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s: %d bytes, want %d", name, len(got), len(want))
		}
	}
	got, err := inflateUpTo(z, none, 1<<20, false)
	check("whole", got, err, text, false)
	got, err = inflateUpTo(z, none, 100, false)
	check("over the limit: a byte more", got, err, text[:101], false)
	got, err = inflateUpTo(z, none, int64(len(text)), false)
	check("at the limit", got, err, text, false)
	got, err = inflateUpTo(z[2:], none, 1<<20, false)
	check("bare deflate", got, err, text, false)
	bad := bytes.Clone(z)
	bad[len(bad)-1]++
	got, err = inflateUpTo(bad, none, 1<<20, false)
	check("wrong checksum", got, err, text, false)
	got, err = inflateUpTo(z[:len(z)/2], none, 1<<20, false)
	if err != nil || len(got) == 0 || len(got) >= len(text) || !bytes.HasPrefix(text, got) {
		t.Errorf("cut short: %d bytes, %v", len(got), err)
	}
	got, err = inflateUpTo(nil, none, 1<<20, false)
	check("no data", got, err, nil, true)
	corrupt := bytes.Clone(z)
	for i := 2; i < len(corrupt)-4; i++ {
		corrupt[i] = 0xff
	}
	if _, err := inflateUpTo(corrupt, none, 1<<20, false); err == nil {
		t.Error("corrupt data: no error")
	}
	if _, err := inflateUpTo([]byte{0x78, 0xbb, 0, 0, 0, 0, 1}, none, 1<<20, false); err == nil {
		t.Error("preset dictionary: no error")
	}
	if _, err := inflateUpTo(z, map[string]int{"Predictor": 5}, 1<<20, false); err == nil {
		t.Error("predictor 5: no error")
	}
	// A predictor's last row cut short is dropped; a wrong filter type fails.
	rowsRaw := []byte{2, 1, 2, 3, 0, 4, 5, 6, 1, 7}
	got, err = inflateUpTo(zlibOf(rowsRaw), map[string]int{"Predictor": 12, "Columns": 3}, 1<<20, false)
	check("short last row", got, err, []byte{1, 2, 3, 4, 5, 6}, false)
	got, err = inflateUpTo(zlibOf([]byte{9, 1, 2, 3}), map[string]int{"Predictor": 12, "Columns": 3}, 1<<20, false)
	check("filter type 9", got, err, nil, true)
}

// TestInflateLarge grows the output past its first guess, and stops at the
// limit without reading the rest.
func TestInflateLarge(t *testing.T) {
	big := make([]byte, 3<<20)
	for i := range big {
		big[i] = byte(i >> 10)
	}
	z := zlibOf(big) // a few KiB
	got, err := inflateUpTo(z, nil, 1<<30, false)
	if err != nil || !bytes.Equal(got, big) {
		t.Fatalf("%d bytes, %v", len(got), err)
	}
	got, err = inflateUpTo(z, nil, 1000, false)
	if err != nil || len(got) != 1001 || cap(got) > 2048 {
		t.Fatalf("limited: %d bytes (cap %d), %v", len(got), cap(got), err)
	}
	// The second slice is sized by the ratio so far: a little over the
	// stream's size. A sized stream gets its slice at once.
	got, err = inflateUpTo(z, nil, 1<<30, false)
	if err != nil || cap(got) > len(got)+len(got)/6 {
		t.Fatalf("projected: %d bytes in %d, %v", len(got), cap(got), err)
	}
	got, err = inflateUpTo(z, nil, int64(len(big)), true)
	if err != nil || len(got) != len(big) || cap(got) != len(big)+1 {
		t.Fatalf("sized: %d bytes in %d, %v", len(got), cap(got), err)
	}
}

func BenchmarkInflateUpTo(b *testing.B) {
	text := bytes.Repeat([]byte("0 0 10 10 re f BT /F1 12 Tf 72 712 Td (Hello) Tj ET\n"), 2000)
	z := zlibOf(text)
	b.Run("own", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(text)))
		for i := 0; i < b.N; i++ {
			if _, err := inflateUpTo(z, nil, 1<<24, false); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("pdfcpu", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(text)))
		fi, _ := filter.NewFilter(filter.Flate, nil)
		for i := 0; i < b.N; i++ {
			if _, err := fi.DecodeLength(bytes.NewReader(z), 1<<24); err != nil && err != io.EOF {
				b.Fatal(err)
			}
		}
	})
}
