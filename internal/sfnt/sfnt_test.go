package sfnt

import (
	"encoding/binary"
	"os"
	"strings"
	"testing"
)

func TestSFNTRebuild(t *testing.T) {
	data, err := os.ReadFile("../../fixture/testdata/fonts/DejaVuSans-sub.ttf")
	if err != nil {
		t.Skip("fixture font missing")
	}
	sf, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	gidA, ok := sf.Cmap['A']
	if !ok || gidA == 0 {
		t.Fatal("no cmap entry for A")
	}
	out := sf.Rebuild(map[uint32]uint16{'Z': gidA, 0xE000: gidA, 0x1F600: gidA}, FontInfo{Family: "Test", Weight: 700, Italic: true})
	sf2, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if sf2.Cmap['Z'] != gidA || sf2.Cmap[0xE000] != gidA || sf2.Cmap[0x1F600] != gidA {
		t.Fatalf("rebuilt cmap = %v", sf2.Cmap)
	}
	if _, has := sf2.Cmap['A']; has {
		t.Fatal("old mapping survived")
	}
	for _, tag := range []string{"glyf", "loca", "head", "hhea", "hmtx", "maxp", "cmap", "name", "OS/2", "post"} {
		if sf2.Tables[tag] == nil {
			t.Fatalf("missing table %s", tag)
		}
	}
	if be16(sf2.Tables["OS/2"], 4) != 700 {
		t.Fatal("weight not written")
	}
}

func TestPruneGlyphs(t *testing.T) {
	data, err := os.ReadFile("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf")
	if err != nil {
		data, err = os.ReadFile("../../fixture/testdata/fonts/DejaVuSans-sub.ttf")
		if err != nil {
			t.Skip("no font available")
		}
	}
	sf, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	before := len(sf.Tables["glyf"])
	// é is a composite of e and acute in DejaVu; both components must survive.
	keep := map[uint16]bool{sf.Cmap['A']: true, sf.Cmap[0xE9]: true}
	sf.PruneGlyphs(keep)
	after := len(sf.Tables["glyf"])
	if after >= before/4 {
		t.Fatalf("glyf not pruned: %d -> %d", before, after)
	}
	out := sf.Rebuild(map[uint32]uint16{'A': sf.Cmap['A'], 0xE9: sf.Cmap[0xE9]}, FontInfo{Family: "Pruned", Weight: 400})
	sf2, err := Parse(out)
	if err != nil {
		t.Fatal(err)
	}
	if sf2.NumGlyphs != sf.NumGlyphs {
		t.Fatal("glyph count changed")
	}
	loca, head := sf2.Tables["loca"], sf2.Tables["head"]
	long := be16(head, 50) != 0
	off := func(g uint16) (uint32, uint32) {
		if long {
			return be32(loca, int(g)*4), be32(loca, int(g)*4+4)
		}
		return uint32(be16(loca, int(g)*2)) * 2, uint32(be16(loca, int(g)*2+2)) * 2
	}
	if s, e := off(sf.Cmap['A']); e <= s {
		t.Fatal("kept glyph is empty")
	}
	if s, e := off(sf.Cmap['Z']); e != s {
		t.Fatal("unused glyph still has data")
	}
	// Components of the composite: the accent glyph must still have outlines.
	if acute, ok := sf.Cmap[0xB4]; ok {
		if s, e := off(acute); e <= s {
			t.Fatal("composite component was pruned")
		}
	}
}

// TestPruneCFF empties the glyphs not kept of a CFF font, keeping their
// numbers: the kept glyphs keep their outlines.
func TestPruneCFF(t *testing.T) {
	data, err := os.ReadFile("testdata/STIXTwoMath-cff-subset.otf")
	if err != nil {
		t.Fatal(err)
	}
	f, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	x, a := f.Cmap['x'], f.Cmap['a']
	before := NewOutlines(f)
	want, ok := before.Bounds(x)
	if !ok || x == 0 || a == 0 {
		t.Fatal("no x or a")
	}
	pruned := *f
	pruned.Tables = map[string][]byte{}
	for tag, b := range f.Tables {
		pruned.Tables[tag] = b
	}
	if !pruned.PruneCFF(map[uint16]bool{x: true}) {
		t.Fatal("not pruned")
	}
	g, err := Parse(Build(pruned.Tables, true))
	if err != nil {
		t.Fatal(err)
	}
	if g.NumGlyphs != f.NumGlyphs || len(g.Tables["CFF "]) >= len(f.Tables["CFF "]) {
		t.Fatalf("%d glyphs (of %d), CFF %d bytes (of %d)", g.NumGlyphs, f.NumGlyphs, len(g.Tables["CFF "]), len(f.Tables["CFF "]))
	}
	after := NewOutlines(g)
	if got, ok := after.Bounds(x); !ok || got != want {
		t.Errorf("x: %v %v, want %v", got, ok, want)
	}
	if _, ok := after.Bounds(a); ok {
		t.Error("a is not emptied")
	}
}

// TestCmapGroups reads format 12 groups that overlap and reach past
// Unicode: at most as many mappings as there are code points.
func TestCmapGroups(t *testing.T) {
	format12 := func(groups [][3]uint32) []byte {
		b := make([]byte, 16+len(groups)*12)
		be := func(i int, v uint32) { b[i], b[i+1], b[i+2], b[i+3] = byte(v>>24), byte(v>>16), byte(v>>8), byte(v) }
		b[1] = 12
		be(12, uint32(len(groups)))
		for i, g := range groups {
			be(16+i*12, g[0])
			be(16+i*12+4, g[1])
			be(16+i*12+8, g[2])
		}
		return b
	}
	var overlapping [][3]uint32
	for i := range uint32(1000) {
		overlapping = append(overlapping, [3]uint32{i * 16, 0xFFFFFFFF, 1})
	}
	left := maxCmap
	if m := parseCmapSubtable(format12(overlapping), &left); len(m) > 0x110000 {
		t.Errorf("overlapping groups: %d mappings", len(m))
	}
	left = maxCmap
	m := parseCmapSubtable(format12([][3]uint32{{0x10FFF0, 0xFFFFFFFF, 5}, {0x110000, 0x120000, 5}}), &left)
	if len(m) != 16 || m[0x10FFFF] != 20 {
		t.Errorf("groups past Unicode: %d mappings, U+10FFFF → %d", len(m), m[0x10FFFF])
	}
}

// words writes 16-bit numbers.
func words(v ...int) []byte {
	var b []byte
	for _, x := range v {
		b = binary.BigEndian.AppendUint16(b, uint16(x))
	}
	return b
}

// cmapOf makes a cmap table of records of platform 3 that all point at
// the subtables, in turn.
func cmapOf(records int, subtables ...[]byte) []byte {
	cm := words(0, records)
	off := 4 + records*8
	offs := make([]int, len(subtables))
	for i, s := range subtables {
		offs[i] = off
		off += len(s)
	}
	for i := range records {
		cm = append(cm, words(3, 1)...)
		cm = binary.BigEndian.AppendUint32(cm, uint32(offs[i%len(offs)]))
	}
	for _, s := range subtables {
		cm = append(cm, s...)
	}
	return cm
}

// TestCmapBounded reads cmap tables that map more than they hold: records
// that share a subtable read it once, and segments that map the same
// characters again and again end in an error.
func TestCmapBounded(t *testing.T) {
	format4 := func(segments int) []byte {
		b := words(4, 0, 0, segments*2, 0, 0, 0)
		for _, v := range []int{0xfffe, -1, 0, 1, 0} { // ends, a pad, starts, deltas, range offsets
			if v < 0 {
				b = append(b, 0, 0)
				continue
			}
			for range segments {
				b = append(b, words(v)...)
			}
		}
		return b
	}
	head := make([]byte, 54)
	font := func(cm []byte) []byte {
		return Build(map[string][]byte{"cmap": cm, "head": head}, false)
	}
	// 4000 records of one subtable of 65535 characters
	var f *Font
	var err error
	shared := font(cmapOf(4000, format4(1)))
	allocs := testing.AllocsPerRun(1, func() { f, err = Parse(shared) })
	if err != nil || len(f.Cmap) != 0xffff || f.Cmap['A'] != 'A'+1 {
		t.Fatalf("records of one subtable: %d characters, error %v", len(f.Cmap), err)
	}
	if allocs > 1000 {
		t.Errorf("records of one subtable: %v allocations, the subtable is read more than once", allocs)
	}
	// segments of the same 65535 characters: more than maxCmap
	if _, err := Parse(font(cmapOf(1, format4(maxCmap/0xffff+1)))); err == nil || !strings.Contains(err.Error(), "cmap") {
		t.Errorf("segments that overlap: error %v", err)
	}
	// subtables that each stay under the limit pass it together
	if _, err := Parse(font(cmapOf(40, format4(1), format4(1), format4(2), format4(3), format4(4), format4(5), format4(6), format4(7), format4(8), format4(9), format4(10), format4(11), format4(12)))); err == nil {
		t.Errorf("subtables that overlap: no error")
	}
}

// TestNoticesBounded reads a name table whose records share one long
// string: the records that cannot replace the notice kept are not decoded,
// and those that are decoded are bounded together.
func TestNoticesBounded(t *testing.T) {
	const length = 60000
	name := func(records int, platform, lang int, text []byte) []byte {
		b := words(0, records, 6+records*12)
		for range records {
			b = append(b, words(platform, 0, lang, 0, len(text), 0)...)
		}
		return append(b, text...)
	}
	head := make([]byte, 54)
	font := func(nt []byte) []byte {
		return Build(map[string][]byte{"name": nt, "head": head}, false)
	}
	var f *Font
	var err error
	shared := font(name(2000, 3, 0x0409, []byte(strings.Repeat("\x00c", length/2))))
	allocs := testing.AllocsPerRun(1, func() { f, err = Parse(shared) })
	if err != nil || len(f.Notices) != 1 || len(f.Notices[0].Text) != length/2 {
		t.Fatalf("records of one string: notices %d, error %v", len(f.Notices), err)
	}
	if allocs > 100 {
		t.Errorf("records of one string: %v allocations, the string is decoded more than once", allocs)
	}
	// Macintosh records that are not ASCII replace nothing: each is read
	if _, err := Parse(font(name(100, 1, 0, []byte(strings.Repeat("\xa9", length))))); err == nil || !strings.Contains(err.Error(), "name") {
		t.Errorf("records that are read again and again: error %v", err)
	}
}
