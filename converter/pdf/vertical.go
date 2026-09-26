package pdf

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cjkcmap"
)

// Vertical writing (WMode 1): glyphs go down the text space y axis, each
// placed by its position vector. A run of them is drawn in a child object
// laid along the run, with the glyphs turned back upright, and placed with
// an ALT_TEXT carrying the run's text, so that extraction, search and
// selection see one run per vertical line (as for vertical text from
// PowerPoint and Word, spec §7.8).

// vertRun is vertical text accumulated until the next flushRun.
type vertRun struct {
	glyphs   []vertGlyph
	font     *pdfFont
	size     float64
	tx0, ty0 float64 // start in text space
	text     strings.Builder
	tgt      *mcTarget
}

type vertGlyph struct {
	draw  string
	u     float64 // distance down the run from its start
	x, y  float64 // horizontal origin from the vertical one, y down
	adv   float64 // horizontal advance
	along float64 // vertical advance (the length it takes in the run)
	turn  bool    // a horizontal form drawn along the run (system fonts)
}

// turnable are characters a system font only has in their horizontal
// form, which vertical text turns with the line.
const turnable = "ー－―‐〜～…‥‖｜＝＿￣（）「」『』【】〔〕［］｛｝〈〉《》"

// showVertical adds the glyphs of a string to the vertical run.
func (in *interp) showVertical(f *pdfFont, codes []glyphCode) {
	if !in.curRun.empty {
		in.flushRun()
	}
	vr := &in.vrun
	if len(vr.glyphs) > 0 && (vr.font != f || vr.size != in.gs.size) {
		in.flushRun()
	}
	size := in.gs.size
	system := f.prog == nil || f.noEmbed
	for _, g := range codes {
		w0 := f.width(g)
		w1, vx, vy := f.vmetrics(g, w0)
		draw, text := f.use(g)
		if in.actual != nil {
			text = ""
			if !in.actual.used {
				text = in.actual.text
				in.actual.used = true
			}
		}
		if len(vr.glyphs) == 0 && vr.text.Len() == 0 {
			vr.font, vr.size = f, size
			vr.tx0, vr.ty0 = in.text.tx, in.text.ty
			vr.tgt = in.curTarget()
		}
		vg := vertGlyph{draw: draw, u: vr.ty0 - in.text.ty, x: -vx * size, y: vy * size, adv: w0 * size, along: math.Abs(w1 * size)}
		if system && draw != "" {
			// A system font draws the character of the CID, which for a
			// vertical CMap is often a vertical presentation form (、 →
			// ︑) whatever the code (a Unicode CMap's 、).
			if r, ok := cjkcmap.VerticalForm(f.collection, g.cid); ok {
				vg.draw = string(r)
			} else if r, vertical, ok := cjkcmap.Unicode(f.collection, g.cid); ok && vertical && strings.ContainsRune(turnable, r) {
				vg.draw, vg.turn = string(r), true
			}
		}
		if draw != "" {
			vr.glyphs = append(vr.glyphs, vg)
		}
		vr.text.WriteString(text)
		ty := w1*size + in.gs.charSp
		if g.nbytes == 1 && g.code == 32 {
			ty += in.gs.wordSp
		}
		in.text.ty += ty
	}
}

// continueVertical keeps a vertical run going when the line matrix only
// moves the pen along it (producers that place every glyph with Td).
func (in *interp) continueVertical(m matrix) bool {
	t := &in.text
	vr := &in.vrun
	if len(vr.glyphs) == 0 || m[0] != t.base[0] || m[1] != t.base[1] || m[2] != t.base[2] || m[3] != t.base[3] {
		return false
	}
	bx, by := t.base[2], t.base[3] // the text space y axis
	l2 := bx*bx + by*by
	if l2 == 0 {
		return false
	}
	dxv, dyv := m[4]-t.base[4], m[5]-t.base[5]
	dy := (dxv*bx + dyv*by) / l2
	rx, ry := dxv-dy*bx, dyv-dy*by
	if rx*rx+ry*ry >= 1e-6*vr.size*vr.size*l2 || math.Abs(dy-t.ty) > 0.25*vr.size || t.tx != vr.tx0 {
		return false
	}
	t.ty = dy
	t.tm, t.tlm = m, m
	return true
}

// flushVertical draws the vertical run.
func (in *interp) flushVertical() {
	vr := &in.vrun
	if len(vr.glyphs) == 0 && vr.text.Len() == 0 {
		return
	}
	defer func() { *vr = vertRun{} }()
	f, size := vr.font, vr.size
	in.syncText(vr.tgt)
	mode := in.gs.render
	if mode == 7 || size == 0 {
		mode = 3
	}
	if mode != 3 {
		in.enterMask()
	}
	in.ensureTextBlock(true)
	ref := in.p.fontRef(f)
	if in.em.font != f || in.em.fontSize != size {
		in.obj.Font(ref, float32(size))
		in.em.font, in.em.fontSize = f, size
	}
	if in.em.letterSp != 0 {
		in.obj.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, 0)
		in.em.letterSp = 0
	}
	switch mode {
	case 0, 4:
		in.syncFill()
		in.syncAlpha(false)
	case 1, 5:
		in.syncStroke()
		in.syncLine()
		in.syncAlpha(true)
	case 2, 6:
		in.syncFill()
		in.syncStroke()
		in.syncLine()
		in.syncAlpha(false)
	default:
		transparent := bdf.RGBA(0, 0, 0, 0)
		if in.em.fillColor != transparent || in.em.fillIsPaint {
			in.obj.FillColor(transparent)
			in.em.fillColor, in.em.fillIsPaint = transparent, false
		}
	}
	// The run's baseline sits below the line of vertical origins so that
	// the em boxes of its characters (0.88 above, 0.12 below) cover the
	// glyphs across it.
	b := 0.38 * size
	length := 0.0
	for _, g := range vr.glyphs {
		length = math.Max(length, g.u+g.along)
	}
	child := newPending(bdf.Rect{X: 0, Y: float32(-b - size), W: float32(length), H: float32(2 * size)})
	in.c.pendings = append(in.c.pendings, child)
	co := child.obj
	co.Font(child.fontRef(f), float32(size))
	for _, g := range vr.glyphs {
		co.Save()
		if g.turn {
			co.Translate(float32(g.u), 0)
			vertText(co, mode, g.draw, 0, 0, float32(g.along))
		} else {
			co.Transform(0, -1, 1, 0, float32(g.u), float32(-b))
			vertText(co, mode, g.draw, float32(g.x), float32(g.y), float32(g.adv))
		}
		co.Restore()
	}
	// Child space runs down the line: turn it by 90° onto the text block.
	in.obj.Save()
	in.obj.Transform(0, 1, -1, 0, float32(vr.tx0), float32(-vr.ty0))
	in.obj.Mark(bdf.MarkAltText, vr.text.String())
	in.obj.UseAt(in.p.childRef(child), 0, float32(b))
	in.obj.Restore()
}

func vertText(o *bdf.Object, mode int, s string, x, y, adv float32) {
	switch mode {
	case 1, 5:
		o.StrokeText(s, x, y, adv)
	case 2, 6:
		o.FillText(s, x, y, adv)
		o.StrokeText(s, x, y, adv)
	default:
		o.FillText(s, x, y, adv)
	}
}
