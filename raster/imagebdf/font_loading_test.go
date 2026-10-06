package imagebdf

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/sfnt"
)

func TestEmbeddedFontLoadsFallbacksOnlyForMissingGlyphs(t *testing.T) {
	data, err := os.ReadFile("../../converter/pptx/testdata/fonts/MPLUS1p-Regular-subset.ttf")
	if err != nil {
		t.Fatal(err)
	}
	f, err := sfnt.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	// Both programs have valid outlines, but expose different characters.
	program := func(family string, chars string) []byte {
		cmap := map[uint32]uint16{}
		for _, r := range chars {
			g := f.Cmap[uint32(r)]
			if g == 0 {
				t.Fatalf("test font lacks %c", r)
			}
			cmap[uint32(r)] = g
		}
		return f.Rebuild(cmap, sfnt.FontInfo{Family: family, Weight: 400})
	}
	d := bdf.NewDocument()
	hash := d.AddFont(program("Embedded", "A"))
	r := New(d, &Options{NoSystemFonts: true, FontFS: fstest.MapFS{
		"fallback.ttf": {Data: program("Fallback", "B")},
		"unused.ttf":   {Data: program("Unused", "C")},
	}})
	ref := bdf.Font{Kind: bdf.FontEmbedded, Hash: hash, Family: "Fallback", Weight: 400}
	c := r.fonts.chain(&ref)
	ref.Family = "Unused" // must not alter the cached chain's deferred lookup
	g, width := c.layout("AA", 12, 0, bdf.DirLTR)
	if len(g) != 2 || width <= 0 || r.fonts.db != nil || len(r.fonts.files) != 0 {
		t.Fatal("drawing the embedded characters loaded external fonts")
	}
	primary := g[0].use.face
	baseline := c.baselineShift(bdf.BaselineTop, 12)
	g, _ = c.layout("AB", 12, 0, bdf.DirLTR)
	if len(g) != 2 || g[0].use.face != primary || g[1].use.face == primary || g[1].gid == 0 {
		t.Fatal("missing character did not use its fallback")
	}
	if len(r.fonts.files) != 1 || c.baselineShift(bdf.BaselineTop, 12) != baseline {
		t.Fatal("unused fallback loaded or primary baseline changed")
	}
	g, _ = c.layout("CBA", 12, 0, bdf.DirLTR)
	if len(g) != 3 || g[0].gid == 0 || g[1].gid == 0 || g[2].use.face != primary || len(r.fonts.files) != 2 {
		t.Fatal("later fallback changed the priority of earlier faces")
	}
	g, _ = c.layout("\U0010ffff", 12, 0, bdf.DirLTR)
	if len(g) != 1 || g[0].gid != 0 || g[0].use.face != primary {
		t.Fatal("missing character must use the primary face's .notdef")
	}
}
