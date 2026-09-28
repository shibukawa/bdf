package font

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"maps"
	"runtime"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/internal/otlayout"
	"github.com/shibukawa/bdf/internal/sfnt"
)

// Fonts are untrusted input. These tests read fonts that say they hold
// more than they do: tables, names and substitutions that share their
// bytes, so that a small file would be read into gigabytes, or for hours.

// allocated returns the bytes fn allocates.
func allocated(fn func()) uint64 {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	fn()
	runtime.ReadMemStats(&after)
	return after.TotalAlloc - before.TotalAlloc
}

// words writes 16-bit numbers.
func words(v ...int) []byte {
	var b []byte
	for _, x := range v {
		b = binary.BigEndian.AppendUint16(b, uint16(x))
	}
	return b
}

// tablesOf returns the tables of a test font, to be changed.
func tablesOf(t *testing.T, name string) map[string][]byte {
	t.Helper()
	f, err := sfnt.Parse(testData(t, name))
	if err != nil {
		t.Fatal(err)
	}
	return maps.Clone(f.Tables)
}

// TestWOFFTablesBounded reads a WOFF file whose tables are one compressed
// table of 4 MiB, 65 times: more than maxSize together.
func TestWOFFTablesBounded(t *testing.T) {
	const size, tables = 4 << 20, 65
	var z bytes.Buffer
	w := zlib.NewWriter(&z)
	chunk := make([]byte, 1<<20)
	for range size >> 20 {
		w.Write(chunk)
	}
	w.Close()
	file := make([]byte, 44)
	copy(file, "wOFF")
	binary.BigEndian.PutUint32(file[4:], 0x00010000)
	binary.BigEndian.PutUint16(file[12:], tables)
	for i := range tables {
		rec := make([]byte, 20)
		copy(rec, fmt.Sprintf("T%03d", i))
		binary.BigEndian.PutUint32(rec[4:], 44+20*tables)
		binary.BigEndian.PutUint32(rec[8:], uint32(z.Len()))
		binary.BigEndian.PutUint32(rec[12:], size)
		file = append(file, rec...)
	}
	file = append(file, z.Bytes()...)
	var err error
	if n := allocated(func() { _, err = load(file) }); err == nil || !strings.Contains(err.Error(), "larger than") || n > 1<<20 {
		t.Errorf("tables of %d MiB together: %d bytes allocated, error %v", tables*size>>20, n, err)
	}
	// the web fonts of the tests are read as before
	if _, err := load(testData(t, "stix.woff")); err != nil {
		t.Error(err)
	}
}

// TestCollectionFaces reads a collection of 64 fonts that are one font of
// a megabyte: its fonts are not copied, those that are not shown are not
// read further, and the fonts that share a layout table read it once.
func TestCollectionFaces(t *testing.T) {
	tb := tablesOf(t, "features.ttf")
	tb["DATA"] = make([]byte, 1<<20)
	data := collection(sfnt.Build(tb, false), 64)
	var fl *file
	var err error
	if n := allocated(func() { fl, err = load(data) }); err != nil || len(fl.faces) != 64 || n > 1<<20 {
		t.Fatalf("%d bytes allocated, error %v", n, err)
	}
	for _, fc := range fl.faces {
		if fc.ready || fc.gsub != nil {
			t.Fatalf("font %d is read before it is shown", fc.index+1)
		}
	}
	a, b := fl.faces[0], fl.faces[63]
	a.prepare()
	b.prepare()
	if a.gsub == nil || a.gsub != b.gsub || a.gpos != b.gpos || a.gdef != b.gdef {
		t.Error("the layout tables the fonts share are read for each")
	}
	if !a.has("abc") || len(a.chars) != a.f.NumGlyphs {
		t.Error("the font is not read")
	}
	// the size of a font is that of the file made of its tables
	for _, name := range []string{"bold.ttc", "stix.woff", "stix.ttf"} {
		fl, err := load(testData(t, name))
		if err != nil {
			t.Fatal(err)
		}
		for _, fc := range fl.faces {
			want := len(sfnt.Build(fc.f.Tables, fc.f.IsCFF))
			if name == "stix.ttf" {
				want = fl.size
			}
			if fc.size != want {
				t.Errorf("%s: font %d of %d bytes, want %d", name, fc.index+1, fc.size, want)
			}
		}
	}
}

// TestNamesBounded reads a name table whose records are one string of
// 60000 bytes: those past maxNames are left out, with a warning.
func TestNamesBounded(t *testing.T) {
	const records, length = 200, 60000
	name := words(0, records, 6+records*12)
	for i := range records {
		name = append(name, words(3, 1, 0x0409, 256+i, length, 0)...)
	}
	name = append(name, strings.Repeat("\x00a", length/2)...)
	var names []nameRec
	var cut bool
	if n := allocated(func() { names, cut = parseNames(name) }); !cut || len(names) != maxNames/length || n > 16*maxNames {
		t.Errorf("%d names, cut %v, %d bytes allocated", len(names), cut, n)
	}
	// language tags are counted with them
	tags := append(words(1, 0, 8+records*4, records), bytes.Repeat(words(length, 0), records)...)
	if names, cut = parseNames(append(tags, name[len(name)-length:]...)); !cut || len(names) != 0 {
		t.Errorf("language tags: %d names, cut %v", len(names), cut)
	}
	// in a font: its own records, then those of the long string
	tb := tablesOf(t, "features.ttf")
	own := tb["name"]
	count, storage := int(binary.BigEndian.Uint16(own[2:])), int(binary.BigEndian.Uint16(own[4:]))
	long := words(0, count+records, 6+(count+records)*12)
	long = append(long, own[6:6+count*12]...)
	for i := range records {
		long = append(long, words(3, 1, 0x0409, 256+i, length, len(own)-storage)...)
	}
	long = append(append(long, own[storage:]...), strings.Repeat("\x00a", length/2)...)
	tb["name"] = long
	var warnings []string
	if _, err := convert(sfnt.Build(tb, false), testOptions(nil), func(m string) { warnings = append(warnings, m) }); err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "name table") {
		t.Errorf("warnings %q", warnings)
	}
	if names, cut := parseNames(own); cut || len(names) == 0 {
		t.Errorf("the names of the test font: %d, cut %v", len(names), cut)
	}
}

// reachBefore is reach as it was: every substitution of every lookup is
// listed again until a pass adds no glyph.
func reachBefore(fc *face, text string) map[uint16]bool {
	keep := map[uint16]bool{0: true}
	for _, r := range text + "◌ -" {
		if g := fc.glyph(r); g != 0 {
			keep[g] = true
		}
	}
	if fc.gsub == nil {
		return keep
	}
	for changed := true; changed; {
		changed = false
		for i := range fc.gsub.Lookups {
			fc.gsub.Substs(i, func(s otlayout.Subst) bool {
				for _, g := range s.In {
					if !keep[g] {
						return true
					}
				}
				for _, g := range s.Out {
					if !keep[g] {
						keep[g] = true
						changed = true
					}
				}
				return true
			})
		}
	}
	return keep
}

// TestReachUnchanged compares the glyphs text reaches with those the
// passes over the lookups found, for the fonts of the tests.
func TestReachUnchanged(t *testing.T) {
	for _, name := range []string{"features.ttf", "stix.ttf", "stix.otf", "bold.ttc", "color.ttf", "variable.ttf", "restricted.ttf"} {
		fl, err := load(testData(t, name))
		if err != nil {
			t.Fatal(err)
		}
		for _, fc := range fl.faces {
			fc.prepare()
			plan := fc.plan("")
			texts := []string{"", "a", "fi", "abcdef", "ffi", "The quick brown fox jumps over the lazy dog.", "0123456789/½", "Âẫ", plan.all(), string(fc.runes)}
			for _, text := range texts {
				got, want := fc.reach(text), reachBefore(fc, text)
				if !maps.Equal(got, want) {
					t.Errorf("%s, font %d, %.20q: %d glyphs, %d before", name, fc.index+1, text, len(got), len(want))
				}
			}
		}
	}
}

// TestClosure applies substitutions that each add the glyph the one
// before them needs: they are listed once, not once for each glyph added.
func TestClosure(t *testing.T) {
	const glyphs = 20000
	listed := 0
	keep := map[uint16]bool{0: true, glyphs: true}
	closure(keep, func(fn func(otlayout.Subst) bool) {
		listed++
		for g := uint16(2); g <= glyphs; g++ {
			fn(otlayout.Subst{In: []uint16{g}, Out: []uint16{g - 1}})
		}
		// glyphs that come twice, one that is never added, none at all
		fn(otlayout.Subst{In: []uint16{5, 5, 7}, Out: []uint16{30001, 30002}})
		fn(otlayout.Subst{In: []uint16{5, 40000}, Out: []uint16{30003}})
		fn(otlayout.Subst{Out: []uint16{30004}})
	})
	if listed != 1 || len(keep) != glyphs+1+3 || !keep[1] || !keep[30001] || !keep[30002] || keep[30003] || !keep[30004] {
		t.Errorf("listed %d times, %d glyphs", listed, len(keep))
	}
}

// TestVariationSequencesBounded counts the variation sequences of
// selectors that share their ranges, and of ranges that overlap.
func TestVariationSequencesBounded(t *testing.T) {
	cmap := func(selectors, ranges int, offset func(i int) int) []byte {
		st := binary.BigEndian.AppendUint32(binary.BigEndian.AppendUint32(words(14), 0), uint32(selectors))
		for i := range selectors {
			st = append(st, 0, 0xfe, byte(i))
			st = binary.BigEndian.AppendUint32(st, uint32(10+selectors*11+offset(i)))
			st = binary.BigEndian.AppendUint32(st, 0)
		}
		for range ranges {
			st = binary.BigEndian.AppendUint32(st, uint32(ranges)) // a count, and a range of one character
		}
		return append(binary.BigEndian.AppendUint32(words(0, 1, 0, 5), 12), append(st, make([]byte, 4)...)...)
	}
	// 1000 selectors of the same ranges: 999 of 233 characters (the low
	// byte of the count, and one), and one of the zeros that follow
	if n, more := variationSequences(cmap(1000, 1000, func(int) int { return 0 })); more || n != 1000*(999*(1000&0xff+1)+1) {
		t.Errorf("selectors that share ranges: %d sequences, more %v", n, more)
	}
	// 3000 selectors of ranges that start one after the other: 4.5 million
	var n int
	var more bool
	allocs := testing.AllocsPerRun(1, func() { n, more = variationSequences(cmap(3000, 3000, func(i int) int { return i * 4 })) })
	if !more || n == 0 || allocs > 5000 {
		t.Errorf("ranges that overlap: %d sequences, more %v", n, more)
	}
	if n, more := variationSequences(tablesOf(t, "stix.ttf")["cmap"]); n != 0 || more {
		t.Errorf("stix.ttf: %d sequences", n)
	}
}

// TestNoGlyphs converts a font that says it has no glyph, may not be
// embedded, and is asked for sample text.
func TestNoGlyphs(t *testing.T) {
	tb := tablesOf(t, "restricted.ttf")
	tb["maxp"] = bytes.Clone(tb["maxp"])
	binary.BigEndian.PutUint16(tb["maxp"][4:], 0)
	_, err := convert(sfnt.Build(tb, false), testOptions(map[string]string{"text": "abc"}), func(string) {})
	if err == nil || !strings.Contains(err.Error(), "no glyphs") {
		t.Errorf("error %v", err)
	}
	// so is a font that may be embedded
	tb = tablesOf(t, "features.ttf")
	tb["maxp"] = bytes.Clone(tb["maxp"])
	binary.BigEndian.PutUint16(tb["maxp"][4:], 0)
	if _, err := load(sfnt.Build(tb, false)); err == nil {
		t.Error("a font without glyphs is read")
	}
}
