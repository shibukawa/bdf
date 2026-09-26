package jbig2

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestFixtures decodes the streams in testdata and compares them with the pages
// Ghostscript rendered from them, or for the features jbig2dec lacks, with the
// reference rendering (see fixture.reference).
func TestFixtures(t *testing.T) {
	names, err := filepath.Glob(filepath.Join("testdata", "*.jb2"))
	if err != nil || len(names) == 0 {
		t.Fatalf("no fixtures: %v", err)
	}
	for _, name := range names {
		base := strings.TrimSuffix(name, ".jb2")
		t.Run(filepath.Base(base), func(t *testing.T) {
			page, globals, want := readFixture(t, base)
			got, err := Decode(page, globals)
			if err != nil {
				t.Fatal(err)
			}
			if d := diff(got, want); d != "" {
				t.Fatal(d)
			}
			if got.Stride != (got.Width+7)/8 || len(got.Data) != got.Stride*got.Height {
				t.Fatalf("stride %d and %d bytes for %dx%d", got.Stride, len(got.Data), got.Width, got.Height)
			}
			if pad := got.Width & 7; pad != 0 {
				for y := 0; y < got.Height; y++ {
					if got.Data[y*got.Stride+got.Stride-1]&(0xff>>uint(pad)) != 0 {
						t.Fatalf("padding bits set in row %d", y)
					}
				}
			}
		})
	}
}

func readFixture(t *testing.T, base string) (page, globals []byte, want *Bitmap) {
	t.Helper()
	page, err := os.ReadFile(base + ".jb2")
	if err != nil {
		t.Fatal(err)
	}
	if g, err := os.ReadFile(base + ".glob"); err == nil {
		globals = g
	}
	p, err := os.ReadFile(base + ".pbm")
	if err != nil {
		t.Fatal(err)
	}
	if want, err = decodePBM(p); err != nil {
		t.Fatal(err)
	}
	return page, globals, want
}

// TestEncoderFixtures checks that the test encoder still produces the fixtures'
// pages, so that the other tests exercise what Ghostscript verified.
func TestEncoderFixtures(t *testing.T) {
	for _, f := range fixtures() {
		t.Run(f.name, func(t *testing.T) {
			got, err := Decode(f.page, f.globals)
			if err != nil {
				t.Fatal(err)
			}
			if d := diff(got, f.want); d != "" {
				t.Fatal(d)
			}
		})
	}
}

// TestStandardTables checks the prefix codes assigned to the standard Huffman
// tables against the codes printed in T.88 Annex B.
func TestStandardTables(t *testing.T) {
	// Codes of each table's lines in the order of standardLines.
	want := [16][]int{
		1:  {0b0, 0b10, 0b110, 0b111},
		2:  {0b0, 0b10, 0b110, 0b1110, 0b11110, 0b111110, 0b111111},
		3:  {0b11111110, 0b0, 0b10, 0b110, 0b1110, 0b11110, 0b11111111, 0b1111110, 0b111110},
		4:  {0b0, 0b10, 0b110, 0b1110, 0b11110, 0b11111},
		5:  {0b1111110, 0b0, 0b10, 0b110, 0b1110, 0b11110, 0b1111111, 0b111110},
		6:  {0b11100, 0b1000, 0b1001, 0b1010, 0b11101, 0b11110, 0b1011, 0b00, 0b010, 0b011, 0b1100, 0b1101, 0b111110, 0b111111},
		7:  {0b1000, 0b000, 0b1001, 0b11010, 0b11011, 0b1010, 0b1011, 0b11100, 0b11101, 0b1100, 0b001, 0b010, 0b011, 0b11110, 0b11111},
		8:  {0b11111100, 0b111111100, 0b11111101, 0b111111101, 0b1111100, 0b1010, 0b00, 0b11010, 0b111010, 0b100, 0b111011, 0b1011, 0b1100, 0b11011, 0b11100, 0b111100, 0b1111101, 0b111101, 0b111111110, 0b111111111, 0b01},
		9:  {0b11111100, 0b111111100, 0b11111101, 0b111111101, 0b1111100, 0b1010, 0b010, 0b011, 0b11010, 0b111010, 0b100, 0b111011, 0b1011, 0b1100, 0b11011, 0b11100, 0b111100, 0b1111101, 0b111101, 0b111111110, 0b111111111, 0b00},
		10: {0b1111010, 0b11111100, 0b1111011, 0b11000, 0b00, 0b11001, 0b110110, 0b1111100, 0b11111101, 0b01, 0b11010, 0b110111, 0b111000, 0b111001, 0b111010, 0b111011, 0b111100, 0b1111101, 0b11111110, 0b11111111, 0b10},
		11: {0b0, 0b10, 0b1100, 0b1101, 0b11100, 0b11101, 0b111100, 0b1111010, 0b1111011, 0b1111100, 0b1111101, 0b1111110, 0b1111111},
		12: {0b0, 0b10, 0b110, 0b11100, 0b11101, 0b111100, 0b1111010, 0b1111011, 0b1111100, 0b1111101, 0b1111110, 0b11111110, 0b11111111},
		13: {0b0, 0b100, 0b1100, 0b11100, 0b1101, 0b101, 0b111010, 0b111011, 0b111100, 0b111101, 0b111110, 0b1111110, 0b1111111},
		14: {0b100, 0b101, 0b0, 0b110, 0b111},
		15: {0b1111100, 0b111100, 0b11100, 0b1100, 0b100, 0b0, 0b101, 0b1101, 0b11101, 0b111101, 0b1111101, 0b1111110, 0b1111111},
	}
	for i := 1; i <= 15; i++ {
		got := huffCodes(standardTables[i].lines)
		for j := range want[i] {
			if got[j] != want[i][j] {
				t.Errorf("table B.%d line %d: code %b, want %b", i, j, got[j], want[i][j])
			}
		}
		// Every value round-trips.
		for _, v := range []int{-3000, -300, -26, -17, -5, -1, 0, 1, 2, 3, 5, 17, 76, 142, 2100, 70000} {
			var w bitWriter
			if !inTable(standardTables[i], v) {
				continue
			}
			encodeHuff(&w, standardTables[i], v, false)
			got, ok, err := standardTables[i].decode(&bitReader{data: w.out})
			if err != nil || !ok || got != v {
				t.Errorf("table B.%d: %d decodes as %d, %v, %v", i, v, got, ok, err)
			}
		}
	}
}

func inTable(t *huffTable, v int) bool {
	for _, l := range t.lines {
		switch {
		case l.kind == lineNormal && int64(v) >= l.low && int64(v) < l.low+1<<uint(l.rangeLen),
			l.kind == lineLower && int64(v) <= l.low,
			l.kind == lineUpper && int64(v) >= l.low:
			return true
		}
	}
	return false
}

// TestArithmetic round-trips decisions and integers through the MQ coder.
func TestArithmetic(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	e := newMQEncoder()
	var cx [8]uint8
	var ic intContext
	iaid := make([]uint8, 1<<9)
	bits := make([]int, 20000)
	ints := make([]int, 500)
	for i := range bits {
		if r.Intn(10) < 1+i%8 {
			bits[i] = 1
		}
		e.encode(&cx[i%8], bits[i])
	}
	for i := range ints {
		ints[i] = int(r.Int63n(1<<(r.Intn(31)))) - 1<<20
		if i%50 == 0 {
			ints[i] = 1<<31 - 1
		}
		e.encodeInt(&ic, ints[i], i%17 == 3)
		e.encodeIAID(iaid, 9, i%512)
	}
	data := e.flush()
	d := newArithDecoder(data)
	cx = [8]uint8{}
	ic = intContext{}
	iaid = make([]uint8, 1<<9)
	for i := range bits {
		if b := d.decode(&cx[i%8]); b != bits[i] {
			t.Fatalf("decision %d: %d, want %d", i, b, bits[i])
		}
	}
	for i := range ints {
		v, ok := d.decodeInt(&ic)
		if oob := i%17 == 3; ok == oob || !oob && v != ints[i] {
			t.Fatalf("integer %d: %d, %v, want %d", i, v, ok, ints[i])
		}
		if id := d.decodeIAID(iaid, 9); id != i%512 {
			t.Fatalf("symbol ID %d: %d", i, id)
		}
	}
}

// TestMMR round-trips bitmaps with long runs and dense changes through MMR.
func TestMMR(t *testing.T) {
	for _, w := range []int{1, 7, 64, 333, 3000, 6000} {
		b := newBitmap(w, 40)
		r := rand.New(rand.NewSource(int64(w)))
		for y := 0; y < 40; y++ {
			x := 0
			for x < w {
				n := 1 + r.Intn(1+[]int{1, 5, 60, 3000}[y%4])
				if y%3 == 0 {
					fillRect(b, x, y, n, 1)
				}
				x += n + r.Intn(1+y%5*20)
			}
		}
		for _, eofb := range []bool{false, true} {
			data := append(encodeMMR(b, eofb), 0xaa)
			d := &decoder{}
			got, n, err := d.decodeMMR(data, w, 40)
			if err != nil {
				t.Fatalf("width %d: %v", w, err)
			}
			if df := diff(got, b); df != "" {
				t.Fatalf("width %d: %s", w, df)
			}
			if n != len(data)-1 {
				t.Errorf("width %d, EOFB %v: used %d bytes of %d", w, eofb, n, len(data)-1)
			}
		}
	}
}

// TestRobustness truncates and corrupts the fixtures and checks that decoding
// neither panics nor takes long.
func TestRobustness(t *testing.T) {
	names, _ := filepath.Glob(filepath.Join("testdata", "*.jb2"))
	r := rand.New(rand.NewSource(1))
	n := 150
	if testing.Short() {
		n = 20
	}
	var slowest time.Duration
	check := func(name string, page, globals []byte) {
		start := time.Now()
		bm, err := Decode(page, globals)
		if err == nil && (bm == nil || len(bm.Data) != bm.Stride*bm.Height) {
			t.Fatalf("%s: inconsistent result", name)
		}
		if err != nil && strings.Contains(err.Error(), "malformed data (") {
			t.Errorf("%s: recovered from a panic: %v", name, err)
		}
		el := time.Since(start)
		slowest = max(slowest, el)
		if el > 200*time.Millisecond {
			t.Logf("%s: %v: %v", name, el, err)
		}
		if el > 3*time.Second {
			t.Errorf("%s: took %v", name, el)
		}
	}
	for _, name := range names {
		page, globals, _ := readFixture(t, strings.TrimSuffix(name, ".jb2"))
		for i := 0; i < n; i++ {
			check(name, page[:r.Intn(len(page)+1)], globals)
			c := append([]byte{}, page...)
			for k := 0; k <= i%4; k++ {
				c[r.Intn(len(c))] ^= 1 << uint(r.Intn(8))
			}
			check(name, c, globals)
			if i%4 == 0 {
				c[r.Intn(len(c))] = byte(r.Intn(256))
				check(name, c, globals)
			}
			if globals != nil {
				g := append([]byte{}, globals...)
				g[r.Intn(len(g))] ^= 1 << uint(r.Intn(8))
				check(name, page, g[:r.Intn(len(g)+1)])
			}
		}
	}
	t.Logf("slowest decode: %v", slowest)
}

// FuzzDecode checks that no input makes Decode panic. The fixtures seed it; run
// go test -fuzz FuzzDecode to explore further.
func FuzzDecode(f *testing.F) {
	names, _ := filepath.Glob(filepath.Join("testdata", "*.jb2"))
	for _, name := range names {
		page, _ := os.ReadFile(name)
		globals, _ := os.ReadFile(strings.TrimSuffix(name, ".jb2") + ".glob")
		f.Add(page, globals)
	}
	type result struct {
		bm  *Bitmap
		err error
	}
	f.Fuzz(func(t *testing.T, page, globals []byte) {
		done := make(chan result, 1)
		go func() {
			bm, err := Decode(page, globals)
			done <- result{bm, err}
		}()
		select {
		case r := <-done:
			if r.err != nil && strings.Contains(r.err.Error(), "malformed data (") {
				t.Fatal(r.err)
			}
			if r.err == nil && len(r.bm.Data) != r.bm.Stride*r.bm.Height {
				t.Fatal("inconsistent bitmap")
			}
		case <-time.After(3 * time.Second):
			t.Fatal("decoding takes too long")
		}
	})
}

// largePage is an A4 page at 300 dpi.
type largePage struct {
	name string
	page []byte
	want *Bitmap
	hash string // of the page Ghostscript renders
}

func a4Text() (*Bitmap, []*Bitmap, *textConfig, []textInst) {
	const w, h = 2480, 3508
	g := glyphs(2024, 80, 30, 34, 38, 42)
	sortByHeight(g)
	r := rand.New(rand.NewSource(2025))
	tc := &textConfig{syms: g, corner: cornerBottomLeft, logStrips: 2}
	var insts []textInst
	for y := 200; y+60 < h-200; y += 50 {
		for x := 240; x < w-300; {
			id := r.Intn(len(g))
			insts = append(insts, textInst{id: id, s: x, t: y + 40})
			x += g[id].Width + 2 + r.Intn(3)
			if r.Intn(8) == 0 {
				x += 12
			}
		}
	}
	return renderText(w, h, tc, insts), g, tc, insts
}

func largePages() []largePage {
	img, g, tc, insts := a4Text()
	w, h := img.Width, img.Height
	s, _ := newPage(w, h, 0)
	sd := plainSD(g, 0, nominalAT(0))
	n := s.add(typeSymbolDict, symbolDictData(sd), segOpts{})
	s.add(typeImmediateText, textRegion(w, h, 0, 0, opOr, tc, insts), segOpts{refs: []uint32{n}})
	s.add(typeEndOfPage, nil, segOpts{})
	text := largePage{"a4_text", s.buf, img, "bac3f49ab60a39584e58bec01fcf2171bd0806417207e61f66099d1c7f4c4133"}

	s, _ = newPage(w, h, 0)
	s.add(typeImmediateGeneric, genericRegion(img, 0, 0, opOr, false, 0, nominalAT(0), true), segOpts{})
	generic := largePage{"a4_generic", s.buf, img, text.hash}

	s, _ = newPage(w, h, 0)
	s.add(typeImmediateGeneric, genericRegion(img, 0, 0, opOr, true, 0, nil, false), segOpts{})
	mmr := largePage{"a4_mmr", s.buf, img, text.hash}
	return []largePage{text, generic, mmr}
}

// TestLargePages decodes A4 pages at 300 dpi: a text region, a generic region
// and an MMR region.
func TestLargePages(t *testing.T) {
	if testing.Short() {
		t.Skip("slow encoder")
	}
	for _, p := range largePages() {
		if h := bitmapHash(p.want); h != p.hash {
			t.Errorf("%s: the encoder's page has hash %s, want %s", p.name, h, p.hash)
		}
		start := time.Now()
		got, err := Decode(p.page, nil)
		el := time.Since(start)
		if err != nil {
			t.Fatalf("%s: %v", p.name, err)
		}
		if d := diff(got, p.want); d != "" {
			t.Fatalf("%s: %s", p.name, d)
		}
		t.Logf("%s: %d bytes decoded in %v", p.name, len(p.page), el)
		if el > 5*time.Second {
			t.Errorf("%s: decoding took %v", p.name, el)
		}
	}
}

func BenchmarkLargePages(b *testing.B) {
	for _, p := range largePages() {
		b.Run(p.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := Decode(p.page, nil); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
