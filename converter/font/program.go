package font

import (
	"github.com/shibukawa/bdf/internal/otlayout"
	"github.com/shibukawa/bdf/internal/sfnt"
)

// layoutTables are the tables that make browsers substitute or move glyphs
// (OpenType, AAT and Graphite layout, kerning), left out of the glyph font
// so that each glyph is drawn as it is.
var layoutTables = map[string]bool{
	"GSUB": true, "GPOS": true, "GDEF": true, "BASE": true, "JSTF": true, "MATH": true, "kern": true,
	"morx": true, "mort": true, "kerx": true, "feat": true, "trak": true, "ankr": true, "prop": true,
	"lcar": true, "opbd": true, "just": true, "bsln": true, "Silf": true, "Glat": true, "Gloc": true,
	"Sill": true, "Feat": true,
}

// otherTables are left out of both embedded fonts: vertical metrics
// (Canvas lays text out horizontally), a signature the changes break, and
// what only applications read.
var otherTables = map[string]bool{
	"vhea": true, "vmtx": true, "VORG": true, "DSIG": true, "PCLT": true, "meta": true, "STAT": true,
}

// variationTables are those of a variable font. The embedded fonts draw the
// default instance: TrueType outlines are that instance without them,
// while CFF2 charstrings need fvar to be read.
var variationTables = map[string]bool{
	"fvar": true, "gvar": true, "avar": true, "cvar": true, "HVAR": true, "VVAR": true, "MVAR": true,
}

// colorBitmapTables hold color glyphs as images; the sample font of a
// font where they are large is left out rather than holding them again.
var colorBitmapTables = []string{"sbix", "CBDT", "SVG "}

// maxWholeSample is the largest font embedded whole as the sample font
// when it cannot be pruned.
const maxWholeSample = 4 << 20

// tablesFor copies the tables of a font the embedded fonts keep.
func (fc *face) tablesFor(drop map[string]bool) map[string][]byte {
	cff2 := fc.f.Tables["CFF2"] != nil
	out := map[string][]byte{}
	for tag, b := range fc.f.Tables {
		if drop[tag] || otherTables[tag] || variationTables[tag] && !(cff2 && (tag == "fvar" || tag == "avar")) {
			continue
		}
		out[tag] = b
	}
	return out
}

// glyphProgram builds the font the glyphs are drawn with: the font's
// outlines and color glyphs with a cmap that maps every glyph from a
// private use character (puaRune), without the layout tables, so that the
// browser draws each glyph as it is.
func glyphProgram(fc *face) []byte {
	t := fc.tablesFor(layoutTables)
	cmap := make(map[uint32]uint16, fc.f.NumGlyphs)
	for g := 1; g < fc.f.NumGlyphs; g++ {
		cmap[uint32(puaRune(uint16(g)))] = uint16(g)
	}
	t["cmap"] = sfnt.BuildCmap(cmap)
	if sb := oneStrike(t["sbix"], fc.f.NumGlyphs); sb != nil {
		t["sbix"] = sb
	}
	if post := t["post"]; len(post) >= 32 {
		// version 3: no glyph names
		p := append([]byte(nil), post[:32]...)
		p[0], p[1], p[2], p[3] = 0, 3, 0, 0
		t["post"] = p
	}
	return sfnt.Build(t, fc.f.IsCFF)
}

// sampleProgram builds the font sample text is drawn with: the font with
// its layout tables, so that the browser shapes the text as applications
// do, holding only the glyphs the text can reach (those of its characters
// and what substitutions make of them) under their own numbers. ok is false
// when no such font is worth embedding: the glyphs are then drawn one by
// one from the glyph font.
func sampleProgram(fc *face, text string, subset bool) (data []byte, ok bool) {
	size := 0
	for _, tag := range colorBitmapTables {
		size += len(fc.f.Tables[tag])
	}
	if size > 1<<20 {
		return nil, false
	}
	t := fc.tablesFor(nil)
	// what AAT and Graphite substitute is not followed (reach): such a
	// font is embedded whole
	whole := t["morx"] != nil || t["mort"] != nil || t["Silf"] != nil
	if !subset || whole {
		if fc.size > maxWholeSample {
			return nil, false
		}
		return sfnt.Build(t, fc.f.IsCFF), true
	}
	keep := fc.reach(text)
	f := *fc.f
	f.Tables = t
	switch {
	case t["glyf"] != nil:
		f.PruneGlyphs(keep)
	case t["CFF "] != nil:
		if !f.PruneCFF(keep) && fc.size > maxWholeSample {
			return nil, false
		}
	default:
		// CFF2 and bitmap fonts are not pruned
		if fc.size > maxWholeSample {
			return nil, false
		}
	}
	return sfnt.Build(f.Tables, f.IsCFF), true
}

// reach returns the glyphs text can end up drawn with: the glyphs of its
// characters, then what any substitution makes of glyphs in the set, until
// nothing is added. Contexts are not checked, so the set may hold more
// than shaping would use.
func (fc *face) reach(text string) map[uint16]bool {
	keep := map[uint16]bool{0: true}
	for _, r := range text + "◌ -" { // the dotted circle shapers insert
		if g := fc.glyph(r); g != 0 {
			keep[g] = true
		}
	}
	if fc.gsub == nil {
		return keep
	}
	closure(keep, func(fn func(otlayout.Subst) bool) {
		for i := range fc.gsub.Lookups {
			fc.gsub.Substs(i, fn)
		}
	})
	return keep
}

// closure adds to keep what substitutions make of glyphs in it, until
// nothing is added. list lists the substitutions, and is called once: a
// substitution whose glyphs are not all in the set waits for those that
// are missing, and is applied when the last of them is added.
func closure(keep map[uint16]bool, list func(fn func(otlayout.Subst) bool)) {
	type waiting struct {
		missing int
		out     []uint16
	}
	var subs []waiting
	waits := map[uint16][]int{} // glyph → the substitutions that wait for it
	var added []uint16
	add := func(out []uint16) {
		for _, g := range out {
			if !keep[g] {
				keep[g] = true
				added = append(added, g)
			}
		}
	}
	for g := range keep {
		added = append(added, g)
	}
	list(func(s otlayout.Subst) bool {
		if len(s.In) == 0 {
			add(s.Out)
			return true
		}
		// every glyph waits to be taken from added, those in the set too
		for _, g := range s.In {
			waits[g] = append(waits[g], len(subs))
		}
		subs = append(subs, waiting{len(s.In), s.Out})
		return true
	})
	for len(added) > 0 {
		g := added[len(added)-1]
		added = added[:len(added)-1]
		for _, i := range waits[g] {
			if subs[i].missing--; subs[i].missing == 0 {
				add(subs[i].out)
			}
		}
		delete(waits, g)
	}
}

// strikePPEM is the resolution of the one bitmap strike the glyph font
// keeps of an sbix table: enough for the glyphs of the views on a display
// of twice the density (the specimen, much larger, is then soft).
const strikePPEM = 64

// oneStrike returns an sbix table holding one of the strikes of b: the
// smallest of at least strikePPEM, or the largest; nil when b has at most
// one strike or cannot be read. Emoji fonts hold a strike for each of many
// sizes, together most of the file.
func oneStrike(b []byte, numGlyphs int) []byte {
	if len(b) < 8 {
		return nil
	}
	n := be32(b, 4)
	if n < 2 || 8+n*4 > len(b) {
		return nil
	}
	best, bestPPEM := -1, 0
	for i := range n {
		o := be32(b, 8+i*4)
		if o+4+(numGlyphs+1)*4 > len(b) {
			return nil
		}
		ppem := be16(b, o)
		switch {
		case best < 0,
			bestPPEM < strikePPEM && ppem > bestPPEM,
			ppem >= strikePPEM && ppem < bestPPEM:
			best, bestPPEM = o, ppem
		}
	}
	// the strike: its header, glyph data offsets (from the strike) and data
	end := best + be32(b, best+4+numGlyphs*4)
	if end <= best || end > len(b) {
		return nil
	}
	out := make([]byte, 0, 12+end-best)
	out = append(out, b[:4]...) // version and flags
	out = append(out, 0, 0, 0, 1, 0, 0, 0, 12)
	return append(out, b[best:end]...)
}
