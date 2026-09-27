package equation

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf/internal/sfnt"
)

// stretchV returns a box of the character r grown vertically to at least
// total points: the first variant of the formula font that is tall enough,
// else an assembly of its parts, else (without a formula font or
// variants) the glyph scaled. A variant is drawn at its own baseline, an
// assembly with its bottom on the baseline. grown is false when the
// glyph is the character's own (it was tall enough already).
func (e *Engine) stretchV(r rune, total float64, v env) (b *Box, grown bool) {
	size := e.size(v)
	gl := e.resolve(r, "", false, false)
	if gl.fc == e.face && e.face != nil {
		base := e.metrics(e.face, gl.gid)
		if (base.ink.YMax-base.ink.YMin)*size >= total {
			return e.glyphBox(r, gl.gid, v), false
		}
		if con := e.mt.Vertical(gl.gid); con != nil {
			for _, vr := range con.Variants {
				if float64(vr.Advance)/e.upem*size >= total {
					return e.glyphBox(r, vr.Glyph, v), vr.Glyph != gl.gid
				}
			}
			if len(con.Parts) > 0 {
				return e.assembly(r, con, total, v, true), true
			}
			if n := len(con.Variants); n > 0 {
				last := con.Variants[n-1].Glyph
				return e.glyphBox(r, last, v), last != gl.gid
			}
		}
		return e.glyphBox(r, gl.gid, v), false
	}
	b = e.token(&Atom{Kind: Op, Text: string(r)}, string(r), v)
	b.single = false
	h := b.H + b.D
	if strings.ContainsRune(pathDelimiters, r) && (total > h*1.1 || !e.hasGlyph(gl)) {
		return e.pathDelimiter(r, total, v), true
	}
	if h > 0 && total > h*1.05 {
		sy := total / h
		return scaleBox(b, 1, sy), true
	}
	return b, false
}

// hasGlyph reports whether the face picked for a character has it.
func (e *Engine) hasGlyph(gl glyph) bool {
	return gl.fc != nil && gl.fc.Loaded != nil && gl.fc.Loaded.Has(gl.r)
}

// stretchH returns a box of the character r grown horizontally to at least
// width points (braces, arrows, wide accents), at its own baseline.
func (e *Engine) stretchH(r rune, width float64, v env) *Box {
	size := e.size(v)
	gl := e.resolve(r, "", false, false)
	var b *Box
	if gl.fc == e.face && e.face != nil {
		base := e.metrics(e.face, gl.gid)
		con := e.mt.Horizontal(gl.gid)
		switch {
		case base.adv*size >= width || con == nil:
			b = e.glyphBox(r, gl.gid, v)
		default:
			for _, vr := range con.Variants {
				if float64(vr.Advance)/e.upem*size >= width {
					b = e.glyphBox(r, vr.Glyph, v)
					break
				}
			}
			if b == nil && len(con.Parts) > 0 {
				b = e.assembly(r, con, width, v, false)
			}
			if b == nil {
				b = e.glyphBox(r, con.Variants[len(con.Variants)-1].Glyph, v)
			}
		}
	} else if strings.ContainsRune("⏞⏟⏜⏝⎴⎵", r) {
		b = e.pathBrace(r, width, v)
	} else {
		b = e.token(&Atom{Kind: Op, Text: string(r)}, string(r), v)
		if w := inkWidth(b); w > 0 && width > w*1.05 && !isCombining(r) {
			b = scaleBox(b, width/w, 1)
		}
	}
	b.single = false
	b.stretchH = true
	return b
}

// inkWidth is the width of what b draws.
func inkWidth(b *Box) float64 {
	x0, x1 := math.Inf(1), math.Inf(-1)
	for _, it := range b.items {
		x0, x1 = math.Min(x0, it.x+it.ink[0]), math.Max(x1, it.x+it.ink[2])
	}
	if x1 < x0 {
		return 0
	}
	return x1 - x0
}

// inkCenter is the middle of what b draws, horizontally.
func inkCenter(b *Box) float64 {
	x0, x1 := math.Inf(1), math.Inf(-1)
	for _, it := range b.items {
		x0, x1 = math.Min(x0, it.x+it.ink[0]), math.Max(x1, it.x+it.ink[2])
	}
	if x1 < x0 {
		return b.W / 2
	}
	return (x0 + x1) / 2
}

// assembly builds a glyph of at least size points from the parts of a
// construction (OpenType MATH's glyph assembly): each extender is repeated
// the same number of times, the fewest that reach the size, and the parts
// overlap evenly, at least by the font's MinConnectorOverlap and at most
// by their connectors.
func (e *Engine) assembly(r rune, con *sfnt.GlyphConstruction, size float64, v env, vertical bool) *Box {
	fs := e.size(v)
	scale := fs / e.upem
	minO := e.c.minOverlap * fs
	nExt := 0
	for _, p := range con.Parts {
		if p.Extender {
			nExt++
		}
	}
	var list []sfnt.GlyphPart
	for reps := 0; ; reps++ {
		list = list[:0]
		for _, p := range con.Parts {
			n := 1
			if p.Extender {
				n = reps
			}
			for i := 0; i < n; i++ {
				list = append(list, p)
			}
		}
		sum := 0.0
		for _, p := range list {
			sum += float64(p.FullAdvance) * scale
		}
		if len(list) > 0 && sum-float64(len(list)-1)*minO >= size || nExt == 0 || reps >= 200 {
			break
		}
	}
	n := len(list)
	sum := 0.0
	maxO := math.Inf(1)
	for i, p := range list {
		sum += float64(p.FullAdvance) * scale
		if i+1 < n {
			maxO = math.Min(maxO, float64(min(p.EndConnector, list[i+1].StartConnector))*scale)
		}
	}
	o := minO
	if n > 1 {
		o = (sum - size) / float64(n-1)
		o = math.Min(o, maxO)
		o = math.Max(o, minO)
	}
	b := newBox()
	pos := 0.0
	for _, p := range list {
		part := e.glyphBox(r, p.Glyph, v)
		if vertical {
			b.add(part, 0, -pos)
			b.W = math.Max(b.W, part.W)
		} else {
			b.add(part, pos, 0)
			b.grow(part, 0)
		}
		pos += float64(p.FullAdvance)*scale - o
	}
	total := pos + o
	if vertical {
		b.H, b.D = total, 0
	} else {
		b.W = total
	}
	b.Italic = float64(con.ItalicsCorrection) * scale
	return b
}
