package font

import (
	"math"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// Colors of the views.
var (
	colText    = bdf.RGB(0x1f, 0x1f, 0x1f)
	colGray    = bdf.RGB(0x6b, 0x6b, 0x6b)
	colRule    = bdf.RGB(0xd9, 0xd9, 0xd9)
	colEmpty   = bdf.RGB(0xf3, 0xf3, 0xf3)
	colVoid    = bdf.RGB(0xe2, 0xe2, 0xe2)
	colAccent  = bdf.RGB(0x1f, 0x5f, 0xbf)
	colBadge   = bdf.RGB(0xe6, 0xef, 0xfb)
	colGuide   = bdf.RGB(0xc8, 0xd8, 0xf0)
	colWhite   = bdf.RGB(0xff, 0xff, 0xff)
	colMissing = bdf.RGB(0xc0, 0x30, 0x30)
)

// puaRune is the private use character the glyph font maps glyph g from:
// Supplementary Private Use Area-A, and B for the last two glyph IDs, which
// would be noncharacters.
func puaRune(g uint16) rune {
	if g < 0xFFFE {
		return 0xF0000 + rune(g)
	}
	return 0x100000 + rune(g-0xFFFE)
}

// painter draws the glyphs of the font being shown.
type painter struct {
	fc *face
	// glyphFont draws any glyph by its private use character (see
	// glyphProgram); nil when the glyphs are drawn as outlines.
	glyphFont *bdf.Font
	// sampleFont draws text as the browser shapes it (see sampleProgram);
	// nil when text is drawn glyph by glyph.
	sampleFont *bdf.Font
}

// pen draws into one object (a strip of a view): it keeps the state the
// instructions so far have set, so as not to repeat them.
type pen struct {
	c     *converter
	cv    *canvas.Canvas
	font  bdf.FontRef
	size  float32
	fill  bdf.Color
	set   bool // font set
	fset  bool // fill set
	paths map[uint16]bdf.PathRef
	// link resolves an anchor to the strip it is in.
	link func(anchor string) string
}

func (p *pen) obj() *bdf.Object { return p.cv.Obj }

func (p *pen) setFont(ref bdf.FontRef, size float64) {
	s := f32(size)
	if !p.set || p.font != ref || p.size != s {
		p.obj().Font(ref, s)
		p.font, p.size, p.set = ref, s, true
	}
}

func (p *pen) setFill(c bdf.Color) {
	if !p.fset || p.fill != c {
		p.obj().FillColor(c)
		p.fill, p.fset = c, true
	}
}

func (p *pen) rect(x, y, w, h float64, c bdf.Color) {
	p.setFill(c)
	p.obj().FillRect(f32(x), f32(y), f32(w), f32(h))
	p.cv.Drawn = true
}

// frame draws the outline of a rectangle with hairlines of color c.
func (p *pen) frame(x, y, w, h, t float64, c bdf.Color) {
	p.rect(x, y, w, t, c)
	p.rect(x, y+h-t, w, t, c)
	p.rect(x, y, t, h, c)
	p.rect(x+w-t, y, t, h, c)
}

// linkTo makes a rectangle a link to an anchor of the view.
func (p *pen) linkTo(x, y, w, h float64, anchor string) {
	if p.link == nil {
		return
	}
	if u := p.link(anchor); u != "" {
		p.obj().Link(f32(x), f32(y), f32(w), f32(h), u)
	}
}

// glyphPath returns the reference of glyph g's outline in the object: a
// path in font units, y down.
func (p *pen) glyphPath(g uint16) (bdf.PathRef, bool) {
	if r, ok := p.paths[g]; ok {
		return r, true
	}
	pp := &pathPen{p: &bdf.Path{}}
	if !p.c.fc.out.Outline(g, pp) || len(pp.p.Verbs) == 0 {
		return 0, false
	}
	if p.paths == nil {
		p.paths = map[uint16]bdf.PathRef{}
	}
	r := p.obj().AddPath(pp.p)
	p.paths[g] = r
	return r, true
}

// pathPen turns an outline into a path, flipping y.
type pathPen struct{ p *bdf.Path }

func (q *pathPen) MoveTo(x, y float64) { q.p.MoveTo(f32(x), f32(-y)) }
func (q *pathPen) LineTo(x, y float64) { q.p.LineTo(f32(x), f32(-y)) }
func (q *pathPen) QuadTo(cx, cy, x, y float64) {
	q.p.QuadTo(f32(cx), f32(-cy), f32(x), f32(-y))
}
func (q *pathPen) CubeTo(ax, ay, bx, by, x, y float64) {
	q.p.CubicTo(f32(ax), f32(-ay), f32(bx), f32(-by), f32(x), f32(-y))
}
func (q *pathPen) Close() { q.p.Close() }

// glyph draws glyph g with its origin at (x, y) at size pt. alt is the text
// the glyph stands for in text extraction ("" for none).
func (p *pen) glyph(g uint16, x, y, size float64, alt string) {
	p.run([]uint16{g}, []float64{0}, x, y, size, alt)
}

// run draws glyphs with their origins at x+dx[i], y (dx in em) at size pt,
// as one piece of text: alt, when not "", is what they stand for.
func (p *pen) run(gs []uint16, dx []float64, x, y, size float64, alt string) {
	if len(gs) == 0 {
		return
	}
	pt := p.c.paint
	p.setFill(colText)
	outline := pt.glyphFont == nil
	for _, g := range gs {
		// glyph 0 cannot be mapped from a character: browsers take a
		// character mapped to it as missing
		outline = outline || g == 0
	}
	if !outline {
		// consecutive glyphs at their advances are one string
		ok := true
		adv := 0.0
		for i, g := range gs {
			if math.Abs(dx[i]-adv) > 1e-6 {
				ok = false
			}
			adv += p.c.fc.advances[g]
		}
		if ok {
			s := make([]byte, 0, len(gs)*4)
			for _, g := range gs {
				s = utf8.AppendRune(s, puaRune(g))
			}
			// the private use characters are not text
			p.obj().Mark(bdf.MarkAltText, alt)
			p.setFont(p.cv.FixedFont(*pt.glyphFont), size)
			p.obj().FillText(string(s), f32(x), f32(y), f32(adv*size))
			p.cv.Drawn = true
			return
		}
	}
	// as outlines (or glyph by glyph where the positions are not the
	// advances)
	k := size / p.c.fc.upem
	var run []bdf.Glyph
	for i, g := range gs {
		if !outline && g != 0 {
			if len(run) > 0 {
				p.pathRun(run, x, y, k, alt)
				run, alt = nil, ""
			}
			p.run([]uint16{g}, []float64{0}, x+dx[i]*size, y, size, alt)
			alt = ""
			continue
		}
		if ref, ok := p.glyphPath(g); ok {
			run = append(run, bdf.Glyph{Path: ref, X: f32(dx[i] * p.c.fc.upem)})
		}
	}
	if len(run) > 0 {
		p.pathRun(run, x, y, k, alt)
	}
}

func (p *pen) pathRun(run []bdf.Glyph, x, y, k float64, alt string) {
	o := p.obj()
	o.Save()
	o.Transform(f32(k), 0, 0, f32(k), f32(x), f32(y))
	if alt != "" {
		o.Mark(bdf.MarkAltText, alt)
	}
	o.FillPathRun(0, run)
	o.Restore()
	p.cv.Drawn = true
}

// sample draws text in the font being shown at size pt with its start at
// (x, y), as the browser shapes it when the sample font is embedded, else
// glyph by glyph; rtl sets the direction of right-to-left text, whose x is
// then the right end.
func (p *pen) sample(s string, x, y, size float64, rtl bool, lang string) {
	p.cv.SetLang(lang)
	defer p.cv.SetLang("")
	if sf := p.c.paint.sampleFont; sf != nil {
		p.setFill(colText)
		p.setFont(p.cv.FixedFont(*sf), size)
		if rtl {
			p.obj().TextStyle(bdf.AlignRight, bdf.BaselineAlphabetic, bdf.DirRTL, 0)
		}
		// no expected width: the browser's shaping decides it
		p.obj().FillText(s, f32(x), f32(y), 0)
		if rtl {
			p.obj().TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirInherit, 0)
		}
		p.cv.Drawn = true
		return
	}
	fc := p.c.fc
	var gs []uint16
	var dx []float64
	w := 0.0
	for _, r := range s {
		g := fc.glyph(r)
		gs = append(gs, g)
		dx = append(dx, w)
		w += fc.advances[g]
	}
	if rtl {
		// right to left, glyph by glyph from the right end
		for i := range dx {
			dx[i] = w - dx[i] - fc.advances[gs[i]]
		}
		x -= w * size
	}
	p.run(gs, dx, x, y, size, s)
}

// text draws a line of UI text measured with the font set.
func (p *pen) text(t *textLine, x, y float64, c bdf.Color) {
	p.setFill(c)
	for _, r := range t.runs {
		if r.s == "" {
			continue
		}
		p.cv.SetLang(r.lang)
		p.setFont(p.cv.Font(r.fc.Use), t.size)
		p.obj().FillText(r.s, f32(x+r.x), f32(y), f32(r.w))
	}
	p.cv.SetLang("")
	p.cv.Drawn = true
}

// label draws a line of UI text; align 0 left, 1 center, 2 right of x.
func (p *pen) label(s string, x, y, size float64, bold bool, c bdf.Color, align int) float64 {
	t := p.c.ui.line(s, size, bold)
	p.text(t, x-t.width*float64(align)/2, y, c)
	return t.width
}

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}

// uiFamily is the family the views' own text is set in; characters it
// lacks fall back to other available fonts.
const uiFamily = "Arial"

// texter measures the views' own text with the fonts of the document.
type texter struct {
	set *fontset.Set
}

// textLine is a line of UI text split into runs of one face.
type textLine struct {
	runs  []textRun
	width float64
	size  float64
}

type textRun struct {
	fc   *fontset.Choice
	s    string
	lang string
	x, w float64
}

func (t *texter) line(s string, size float64, bold bool) *textLine {
	l := &textLine{size: size}
	var cur *textRun
	x := 0.0
	for _, r := range s {
		fc := t.set.FaceFor(uiFamily, "", bold, false, r)
		adv := t.set.Advance(fc, r) * size
		lang := fontset.RuneLang(r)
		if cur == nil || cur.fc != fc || lang != cur.lang && lang != "" {
			l.runs = append(l.runs, textRun{fc: fc, x: x, lang: lang})
			cur = &l.runs[len(l.runs)-1]
		}
		cur.s += string(r)
		cur.w += adv
		x += adv
	}
	l.width = x
	return l
}

func (t *texter) width(s string, size float64, bold bool) float64 {
	return t.line(s, size, bold).width
}

// fit shortens s with an ellipsis until it is at most w wide.
func (t *texter) fit(s string, size float64, bold bool, w float64) string {
	adv := func(r rune) float64 { return t.set.Advance(t.set.FaceFor(uiFamily, "", bold, false, r), r) * size }
	rs := []rune(s)
	x, n := 0.0, 0
	for n < len(rs) && x+adv(rs[n]) <= w {
		x += adv(rs[n])
		n++
	}
	if n == len(rs) {
		return s
	}
	// room for the ellipsis
	e := adv('…')
	for n > 0 && x+e > w {
		n--
		x -= adv(rs[n])
	}
	return string(rs[:n]) + "…"
}

// wrap breaks s into lines at most w wide, at spaces, or anywhere in text
// without spaces (East Asian names, long URLs).
func (t *texter) wrap(s string, size float64, bold bool, w float64) []string {
	var out []string
	for _, para := range splitLines(s) {
		rs := []rune(para)
		for len(rs) > 0 {
			x, brk, end := 0.0, -1, len(rs)
			for i, r := range rs {
				x += t.set.Advance(t.set.FaceFor(uiFamily, "", bold, false, r), r) * size
				if x > w && i > 0 {
					end = i
					break
				}
				if r == ' ' || r == '/' || r == '-' || r == ',' {
					brk = i + 1
				}
			}
			if end < len(rs) && brk > 0 {
				end = brk
			}
			line := string(rs[:end])
			out = append(out, trimSpace(line))
			rs = rs[end:]
			for len(rs) > 0 && rs[0] == ' ' {
				rs = rs[1:]
			}
		}
	}
	return out
}

func splitLines(s string) []string {
	var out []string
	cur := []rune{}
	for _, r := range s {
		switch r {
		case '\r':
		case '\n', 0x2028, 0x2029:
			out = append(out, string(cur))
			cur = cur[:0]
		default:
			cur = append(cur, r)
		}
	}
	return append(out, string(cur))
}

func trimSpace(s string) string {
	for len(s) > 0 && s[len(s)-1] == ' ' {
		s = s[:len(s)-1]
	}
	return s
}
