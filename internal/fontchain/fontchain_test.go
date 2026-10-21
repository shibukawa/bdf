package fontchain

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/sfnt"
)

// face is the least a user of the package keeps of a font program.
type face struct{ Program }

func newFace(f *sfnt.Font, _ bool) *face { return &face{NewProgram(f)} }

// An embedded face is drawn with alone while it has the glyphs: the fonts
// of the system are looked up when a character is missing from it, and
// only those that have the character are loaded.
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
	fs := &Fonts[*face]{Doc: d, NoSystemFonts: true, FontFS: fstest.MapFS{
		"fallback.ttf": {Data: program("Fallback", "B")},
		"unused.ttf":   {Data: program("Unused", "C")},
	}, New: newFace, Warnf: func(format string, args ...any) { t.Errorf(format, args...) }, NotDone: "not drawn"}
	ref := bdf.Font{Kind: bdf.FontEmbedded, Hash: hash, Family: "Fallback", Weight: 400}
	c := fs.Chain(&ref)
	ref.Family = "Unused" // must not alter the cached chain's deferred lookup
	g, width := c.Layout("AA", 12, 0, bdf.DirLTR)
	if len(g) != 2 || width <= 0 || fs.db != nil || len(fs.files) != 0 {
		t.Fatal("drawing the embedded characters loaded external fonts")
	}
	primary := g[0].Use.Face
	baseline := c.BaselineShift(bdf.BaselineTop, 12)
	g, _ = c.Layout("AB", 12, 0, bdf.DirLTR)
	if len(g) != 2 || g[0].Use.Face != primary || g[1].Use.Face == primary || g[1].GID == 0 {
		t.Fatal("missing character did not use its fallback")
	}
	if len(fs.files) != 1 || c.BaselineShift(bdf.BaselineTop, 12) != baseline {
		t.Fatal("unused fallback loaded or primary baseline changed")
	}
	g, _ = c.Layout("CBA", 12, 0, bdf.DirLTR)
	if len(g) != 3 || g[0].GID == 0 || g[1].GID == 0 || g[2].Use.Face != primary || len(fs.files) != 2 {
		t.Fatal("later fallback changed the priority of earlier faces")
	}
	g, _ = c.Layout("\U0010ffff", 12, 0, bdf.DirLTR)
	if len(g) != 1 || g[0].GID != 0 || g[0].Use.Face != primary {
		t.Fatal("missing character must use the primary face's .notdef")
	}
	if len(fs.All) != 3 {
		t.Fatalf("%d faces loaded", len(fs.All))
	}
}
