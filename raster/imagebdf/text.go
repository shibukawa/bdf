package imagebdf

import (
	"github.com/shibukawa/bdf/internal/fontchain"
	"github.com/shibukawa/bdf/internal/sfnt"
	"github.com/shibukawa/bdf/internal/textrun"
)

// fontFace is a parsed font program.
type fontFace struct {
	fontchain.Program
	out    *sfnt.Outlines
	glyphs map[uint16]*path
}

func newFace(f *sfnt.Font, _ bool) *fontFace {
	return &fontFace{Program: fontchain.NewProgram(f), out: sfnt.NewOutlines(f), glyphs: map[uint16]*path{}}
}

// The faces of a document and the chains of a font reference are found by
// internal/fontchain, which lays text out on them too.
type (
	fonts  = fontchain.Fonts[*fontFace]
	placed = fontchain.Placed[*fontFace]
)

// outline returns the outline of a glyph in font units (y up).
func (ff *fontFace) outline(g uint16) *path {
	if p, ok := ff.glyphs[g]; ok {
		return p
	}
	pp := &pathPen{p: &path{}}
	if !ff.out.Outline(g, pp) {
		pp.p = nil
	}
	ff.glyphs[g] = pp.p
	return pp.p
}

// pathPen builds a path from a glyph outline.
type pathPen struct{ p *path }

func (pp *pathPen) MoveTo(x, y float64)                 { pp.p.moveTo(x, y) }
func (pp *pathPen) LineTo(x, y float64)                 { pp.p.lineTo(x, y) }
func (pp *pathPen) QuadTo(cx, cy, x, y float64)         { pp.p.quadTo(cx, cy, x, y) }
func (pp *pathPen) CubeTo(ax, ay, bx, by, x, y float64) { pp.p.cubicTo(ax, ay, bx, by, x, y) }
func (pp *pathPen) Close()                              { pp.p.close() }

// fakeBoldWidth is the outline growth of synthesized bold for a size in
// device pixels.
func fakeBoldWidth(size float64) float64 { return textrun.FakeBoldWidth(size) }

// fakeItalicSkew is the slant of synthesized italic (x per unit of height).
const fakeItalicSkew = textrun.FakeItalicSkew
