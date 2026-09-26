package metafile

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// textItem is one character of a text record, placed along the baseline.
type textItem struct {
	r    rune
	fc   *fontset.Choice // nil for tabs and line breaks, which are not drawn
	x, w float64
}

// text draws a string at the logical reference point; dx are the logical
// advances of the characters (nil: the font's).
func (g *gdi) text(x, y float64, s []rune, dx []float64, opaque *[4]float64) {
	if len(s) == 0 {
		return
	}
	m := g.matrix()
	if opaque != nil {
		save := g.st.brush
		g.st.brush = &gdiBrush{color: g.st.bkColor, hatch: -1}
		pen := g.st.pen
		g.st.pen = &gdiPen{null: true}
		g.rect(opaque[0], opaque[1], opaque[2], opaque[3])
		g.st.brush, g.st.pen = save, pen
	}
	if g.st.textAlign&1 != 0 { // TA_UPDATECP
		x, y = g.st.cur[0], g.st.cur[1]
	}
	f := g.st.font
	// scale of the logical space along the text and across it
	sx := math.Hypot(m[0], m[1])
	sy := math.Hypot(m[2], m[3])
	em := math.Abs(f.height)
	face, bold := f.face, f.weight >= 600
	text := string(s)
	if fontset.IsSymbol(face) {
		text = fontset.MapSymbol(face, text)
		face = ""
	}
	fonts := g.opts.Fonts
	fc0 := fonts.FaceFor(face, face, bold, f.italic, []rune(text)[0])
	if f.height > 0 {
		// cell height: the em is the cell minus the internal leading
		em = f.height / math.Max(fc0.Asc+fc0.Desc, 0.5)
	}
	if em == 0 {
		em = 12 / math.Max(sy, 1e-9)
	}
	size := em * sy
	var items []textItem
	for _, r := range text {
		switch r {
		case '\r':
			continue
		case '\v', '\t':
			items = append(items, textItem{r: r})
			continue
		}
		fc := fonts.FaceFor(face, face, bold, f.italic, r)
		items = append(items, textItem{r: r, fc: fc, w: fonts.Advance(fc, r) * size})
	}
	adv := 0.0
	for i := range items {
		items[i].x = adv
		if i < len(dx) {
			items[i].w = dx[i] * sx
		}
		adv += items[i].w
	}
	// reference point to baseline
	asc, desc := fc0.Asc*size, fc0.Desc*size
	var oy float64
	switch g.st.textAlign & 24 {
	case 0: // TA_TOP
		oy = asc
	case 8: // TA_BOTTOM
		oy = -desc
	}
	var ox float64
	switch g.st.textAlign & 6 {
	case 6: // TA_CENTER
		ox = -adv / 2
	case 2: // TA_RIGHT
		ox = -adv
	}
	if g.st.textAlign&1 != 0 {
		// advance the current position
		g.st.cur[0] += adv / math.Max(sx, 1e-9)
	}
	px, py := m.Apply(x, y)
	// the text stands upright in output space even when the logical y
	// axis points up; escapement turns it counter-clockwise
	rot := math.Atan2(m[1], m[0])*180/math.Pi - f.escapement/10
	g.end()
	g.obj.Save()
	for _, poly := range g.st.clips {
		g.obj.ClipPath(g.obj.AddPath(polyPath(poly, true)), 0)
	}
	g.cv.Transform(canvas.Translate(px, py).Mul(canvas.Rotate(rot)))
	g.obj.Mark(bdf.MarkBox, "")
	g.emitText(items, ox, oy, size, toBDF(g.color(g.st.textColor)), f.underline, f.strike)
	g.obj.Restore()
}

// emitText draws the characters as runs of one face each, with the line
// decorations of the font; trailing spaces are not drawn.
func (g *gdi) emitText(items []textItem, dx, base, size float64, c bdf.Color, underline, strike bool) {
	var font bdf.FontRef
	hasFont, hasColor := false, false
	var cur bdf.Color
	setColor := func(c bdf.Color) {
		if !hasColor || c != cur {
			g.obj.FillColor(c)
			cur, hasColor = c, true
		}
	}
	for i := 0; i < len(items); {
		it := items[i]
		if it.fc == nil {
			i++
			continue
		}
		j := i + 1
		for j < len(items) && items[j].fc == it.fc {
			j++
		}
		k := j
		if j == len(items) {
			for k > i && isBreakSpace(items[k-1].r) {
				k--
			}
		}
		if k == i {
			i = j
			continue
		}
		run := items[i:k]
		i = j
		g.setLang(run)
		var b strings.Builder
		adv := 0.0
		for _, r := range run {
			b.WriteRune(r.r)
			adv += r.w
		}
		x := it.x + dx
		if ref := g.cv.Font(it.fc.Use); !hasFont || ref != font {
			g.obj.Font(ref, f32(size))
			font, hasFont = ref, true
		}
		setColor(c)
		g.obj.FillText(b.String(), f32(x), f32(base), f32(adv))
		g.cv.Drawn = true
		th := math.Max(it.fc.ULThick*size, 0.5)
		if underline {
			setColor(c)
			g.obj.FillRect(f32(x), f32(base-it.fc.ULPos*size), f32(adv), f32(th))
		}
		if strike {
			setColor(c)
			g.obj.FillRect(f32(x), f32(base-0.3*size), f32(adv), f32(th))
		}
	}
}

// setLang sets the language of a run: East Asian text takes the language
// its kana or hangul imply, and otherwise keeps an East Asian language in
// effect or takes the document's (Han characters alone could be Chinese as
// well as Japanese or Korean); other text takes the document's.
func (g *gdi) setLang(run []textItem) {
	cjk := false
	for _, it := range run {
		if fontdb.IsCJK(it.r) {
			cjk = true
			break
		}
	}
	lang := ""
	if cjk {
		for _, it := range run {
			if lang = fontset.RuneLang(it.r); lang != "" {
				break
			}
		}
		if lang == "" && fontset.Script(g.cv.Lang()) != "" {
			return
		}
	}
	g.cv.SetLang(lang)
}

func isBreakSpace(r rune) bool {
	return r == ' ' || r == '　' || r == ' ' || r == ' ' || r == ' '
}
