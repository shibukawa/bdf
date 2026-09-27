package sfnt

import (
	"os"
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
	if m := parseCmapSubtable(format12(overlapping)); len(m) > 0x110000 {
		t.Errorf("overlapping groups: %d mappings", len(m))
	}
	m := parseCmapSubtable(format12([][3]uint32{{0x10FFF0, 0xFFFFFFFF, 5}, {0x110000, 0x120000, 5}}))
	if len(m) != 16 || m[0x10FFFF] != 20 {
		t.Errorf("groups past Unicode: %d mappings, U+10FFFF → %d", len(m), m[0x10FFFF])
	}
}
