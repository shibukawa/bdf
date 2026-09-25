package sfnt

import (
	"os"
	"testing"
)

func TestSFNTRebuild(t *testing.T) {
	data, err := os.ReadFile("../../../fixture/fonts/DejaVuSans-sub.ttf")
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
	out := sf.Rebuild(map[uint32]uint16{'Z': gidA, 0xE000: gidA, 0x1F600: gidA}, "Test", 700, true)
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
		data, err = os.ReadFile("../../../fixture/fonts/DejaVuSans-sub.ttf")
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
	out := sf.Rebuild(map[uint32]uint16{'A': sf.Cmap['A'], 0xE9: sf.Cmap[0xE9]}, "Pruned", 400, false)
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
