package pdf

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

func strconvFormat(v float64) string { return strconv.FormatFloat(v, 'f', 4, 64) }

func (in *interp) beginText() {
	in.endText()
	in.inText = true
	in.text = textState{tm: identity, tlm: identity, base: identity}
}

func (in *interp) endText() {
	in.flushRun()
	in.closeTextBlock()
	in.inText = false
}

func (in *interp) setLineMatrix(m matrix) {
	t := &in.text
	// Producers that position every glyph with Td (Skia, some TeX drivers) keep the
	// same line: when the new matrix is the old one shifted along its baseline to
	// roughly where the pen already is, continue the current run.
	if !in.curRun.empty && !in.curRun.font.type3 && m[0] == t.base[0] && m[1] == t.base[1] && m[2] == t.base[2] && m[3] == t.base[3] {
		ax, ay := t.base[0], t.base[1]
		dxv, dyv := m[4]-t.base[4], m[5]-t.base[5]
		if l2 := ax*ax + ay*ay; l2 > 0 {
			dx := (dxv*ax + dyv*ay) / l2
			// residual perpendicular to the baseline
			rx, ry := dxv-dx*ax, dyv-dx*ay
			size := in.gs.size
			if size <= 0 {
				size = 1
			}
			if rx*rx+ry*ry < 1e-6*size*size && math.Abs(dx-t.tx) <= 0.25*size {
				in.curRun.adv += (dx - t.tx) / nonZero(in.gs.hscale)
				t.tx = dx
				t.tm = m
				t.tlm = m
				return
			}
		}
	}
	in.flushRun()
	t.tlm = m
	t.tm = m
	t.base = m
	t.tx = 0
}

// closeTextBlock ends the SAVE/TRANSFORM block opened for the current line matrix.
func (in *interp) closeTextBlock() {
	if in.text.open {
		in.textRestore()
		in.text.open = false
	}
}

// ensureTextBlock opens a transform block for the current line matrix so text
// can be placed with plain x offsets. flip selects a y-down glyph space
// (for fillText) or the PDF y-up space (for Type3 glyph objects).
func (in *interp) ensureTextBlock(flip bool) {
	t := &in.text
	if t.open && t.openTlm == t.base && t.openH == in.gs.hscale && t.openRise == in.gs.rise && t.openFlip == flip {
		return
	}
	in.closeTextBlock()
	in.textSave()
	if !t.base.isIdentity() {
		in.transform(t.base)
	}
	if in.gs.hscale != 1 || in.gs.rise != 0 {
		in.transform(matrix{in.gs.hscale, 0, 0, 1, 0, in.gs.rise})
	}
	if flip {
		in.transform(matrix{1, 0, 0, -1, 0, 0})
	}
	t.open, t.openTlm, t.openH, t.openRise, t.openFlip = true, t.base, in.gs.hscale, in.gs.rise, flip
}

type textRun struct {
	draw  strings.Builder
	text  strings.Builder
	x0    float64 // start position in line space (text space / hscale)
	adv   float64 // advance in line space
	alt   bool
	empty bool
	font  *pdfFont
	tgt   *mcTarget // structure of the marked content the run started in
}

// flushRun emits the run that is being accumulated across TJ elements.
func (in *interp) flushRun() {
	if in.curRun.empty {
		return
	}
	in.emitRun(in.curRun.font, &in.curRun)
	in.curRun = textRun{empty: true}
}

// adjustText applies a TJ number: a horizontal shift in thousandths of text space.
// Small shifts (kerning) are absorbed into the current run; large ones end it.
func (in *interp) adjustText(v float64) {
	shift := -v / 1000 * in.gs.size * in.gs.hscale
	if math.Abs(v) <= 150 && !in.curRun.empty && !in.curRun.font.type3 {
		in.curRun.adv += shift / nonZero(in.gs.hscale)
	} else {
		in.flushRun()
	}
	in.text.tx += shift
}

func nonZero(h float64) float64 {
	if h == 0 {
		return 1
	}
	return h
}

// showText draws a string operand with the current text state.
func (in *interp) showText(s []byte) {
	f := in.gs.font
	if f == nil {
		f = in.c.defaultFont()
		in.gs.font = f
	}
	if !in.inText {
		in.beginText()
	}
	h := nonZero(in.gs.hscale)
	codes := f.decode(s)
	if in.hidden > 0 {
		// Hidden optional content: the glyphs are not drawn, but they still
		// move the pen.
		for _, g := range codes {
			adv := f.width(g)*in.gs.size + in.gs.charSp
			if g.nbytes == 1 && g.code == 32 {
				adv += in.gs.wordSp
			}
			in.text.tx += adv * h
		}
		return
	}
	if f.type3 {
		in.flushRun()
		in.showType3(f, codes)
		return
	}
	if !in.curRun.empty && in.curRun.font != f {
		in.flushRun()
	}
	for _, g := range codes {
		w0 := f.width(g)
		adv := w0*in.gs.size + in.gs.charSp
		isSpace := g.nbytes == 1 && g.code == 32
		if isSpace {
			adv += in.gs.wordSp
		}
		draw, text := f.use(g)
		if in.actual != nil {
			// /ActualText marked content overrides the Unicode of the glyphs it wraps.
			text = ""
			if !in.actual.used {
				text = in.actual.text
				in.actual.used = true
			}
		}
		if isSpace && in.gs.wordSp != 0 {
			// Word spacing only affects this glyph: give it a run of its own so the
			// advance correction stretches the space and not the words around it.
			in.flushRun()
		}
		if in.curRun.empty {
			in.curRun.x0 = in.text.tx / h
			in.curRun.empty = false
			in.curRun.font = f
			in.curRun.tgt = in.curTarget()
		}
		in.curRun.draw.WriteString(draw)
		in.curRun.text.WriteString(text)
		if draw != text {
			in.curRun.alt = true
		}
		in.curRun.adv += adv
		in.text.tx += adv * h
		if isSpace && in.gs.wordSp != 0 {
			in.flushRun()
		}
	}
}

func (in *interp) emitRun(f *pdfFont, run *textRun) {
	in.syncText(run.tgt)
	mode := in.gs.render
	if mode == 7 || in.gs.size == 0 {
		mode = 3
	}
	in.ensureTextBlock(true)
	ref := in.p.fontRef(f)
	if in.em.font != f || in.em.fontSize != in.gs.size {
		in.obj.Font(ref, float32(in.gs.size))
		in.em.font, in.em.fontSize = f, in.gs.size
	}
	if in.em.letterSp != in.gs.charSp {
		in.obj.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, float32(in.gs.charSp))
		in.em.letterSp = in.gs.charSp
	}
	if run.alt {
		in.obj.Mark(bdf.MarkAltText, run.text.String())
	}
	draw := run.draw.String()
	x, adv := float32(run.x0), float32(run.adv)
	switch mode {
	case 0, 4:
		in.syncFill()
		in.syncAlpha(false)
		in.obj.FillText(draw, x, 0, adv)
	case 1, 5:
		in.syncStroke()
		in.syncLine()
		in.syncAlpha(true)
		in.obj.StrokeText(draw, x, 0, adv)
	case 2, 6:
		in.syncFill()
		in.syncAlpha(false)
		in.obj.FillText(draw, x, 0, adv)
		in.syncStroke()
		in.syncLine()
		in.syncAlpha(true)
		in.obj.StrokeText(draw, x, 0, adv)
	default: // invisible: keep the text for search and selection
		transparent := bdf.RGBA(0, 0, 0, 0)
		if in.em.fillColor != transparent || in.em.fillIsPaint {
			in.obj.FillColor(transparent)
			in.em.fillColor, in.em.fillIsPaint = transparent, false
		}
		in.obj.FillText(draw, x, 0, adv)
	}
}

// showType3 draws Type3 glyphs as child objects.
func (in *interp) showType3(f *pdfFont, codes []glyphCode) {
	h := in.gs.hscale
	if h == 0 {
		h = 1
	}
	var text strings.Builder
	for _, g := range codes {
		_, uni := f.use(g)
		text.WriteString(uni)
	}
	altPending := text.Len() > 0
	synced := false
	sync := func() {
		if !synced {
			in.syncText(in.curTarget())
			synced = true
		}
	}
	for _, g := range codes {
		w0 := f.width(g)
		adv := w0*in.gs.size + in.gs.charSp
		if g.nbytes == 1 && g.code == 32 {
			adv += in.gs.wordSp
		}
		glyph := in.c.type3Glyph(f, f.glyphName(g.code), in.depth+1)
		if glyph != nil && in.gs.render != 3 && in.gs.render != 7 {
			sync()
			in.ensureTextBlock(false)
			if altPending {
				in.obj.Mark(bdf.MarkAltText, text.String())
				altPending = false
			}
			in.syncFill()
			in.syncStroke()
			in.syncAlpha(false)
			in.save()
			in.transform(matrix{in.gs.size, 0, 0, in.gs.size, in.text.tx / h, 0})
			in.transform(f.fontMatrix)
			in.obj.Use(in.p.childRef(glyph))
			in.restore()
		}
		in.text.tx += adv * h
	}
	if altPending {
		// Nothing was drawn (invisible or missing glyphs): keep the text for search.
		sync()
		in.obj.Mark(bdf.MarkAltText, text.String())
	}
}
