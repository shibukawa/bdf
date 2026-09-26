package cad

import (
	"math"
	"slices"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
)

// Plotter writes drawings onto pages.
type Plotter struct {
	Fonts *Fonts
	// Thin is the width in pt of the lines whose width is 0.
	Thin float64
}

// state is what has been written of the drawing state, per SAVE level.
type state struct {
	fill, stroke       bdf.Color
	fillSet, strokeSet bool
	width              float64
	cap, join          byte
	lineSet            bool
	dash               []float64
	dashOff            float64
	dashSet            bool
	font               bdf.FontRef
	size               float64
	fontSet            bool
}

type plot struct {
	pl    *Plotter
	cv    *canvas.Canvas
	obj   *bdf.Object
	st    state
	stack []state
	// view is the part of the page that is drawn on (page units); items
	// wholly outside are left out
	view Rect
}

// Plot draws d into cv; m maps the drawing into the page (pt, y down).
// view is the visible part of the page: items outside it are left out.
func (pl *Plotter) Plot(cv *canvas.Canvas, d *Drawing, m canvas.Matrix, view Rect) {
	p := &plot{pl: pl, cv: cv, obj: cv.Obj, view: view}
	for _, it := range d.Items {
		p.item(it, m)
	}
}

func (p *plot) save() {
	p.obj.Save()
	st := p.st
	st.dash = slices.Clone(p.st.dash)
	p.stack = append(p.stack, st)
}

func (p *plot) restore() {
	p.obj.Restore()
	n := len(p.stack)
	p.st = p.stack[n-1]
	p.stack = p.stack[:n-1]
}

func (p *plot) visible(r Rect, m canvas.Matrix) bool {
	if !p.view.ok {
		return true
	}
	if !r.ok {
		return false
	}
	pr := r.Transform(m)
	// the widest lines reach a little outside their path
	const margin = 10
	return pr.Max.X >= p.view.Min.X-margin && pr.Min.X <= p.view.Max.X+margin &&
		pr.Max.Y >= p.view.Min.Y-margin && pr.Min.Y <= p.view.Max.Y+margin
}

func (p *plot) item(it Item, m canvas.Matrix) {
	if !p.visible(it.bounds, m) {
		return
	}
	switch it.kind {
	case kStroke:
		p.strokeItem(it, m)
	case kFill:
		p.fillItem(it, m)
	case kText:
		p.text(it.text, m)
	case kGroup:
		p.save()
		if it.clip != nil {
			p.obj.ClipPath(p.obj.AddPath(toBDF(it.clip, m)), 1)
		}
		for _, c := range it.items {
			p.item(c, m)
		}
		p.restore()
	}
}

// toBDF converts a path to page coordinates.
func toBDF(path *Path, m canvas.Matrix) *bdf.Path {
	out := &bdf.Path{}
	path.walk(func(v byte, pts []Point) {
		switch v {
		case moveTo:
			q := Apply(m, pts[0])
			out.MoveTo(f32(q.X), f32(q.Y))
		case lineTo:
			q := Apply(m, pts[0])
			out.LineTo(f32(q.X), f32(q.Y))
		case cubicTo:
			a, b, c := Apply(m, pts[0]), Apply(m, pts[1]), Apply(m, pts[2])
			out.CubicTo(f32(a.X), f32(a.Y), f32(b.X), f32(b.Y), f32(c.X), f32(c.Y))
		case closePath:
			out.Close()
		}
	})
	return out
}

func (p *plot) strokeItem(it Item, m canvas.Matrix) {
	pen := it.pen
	s := Scale(m)
	w := pen.Width
	if pen.WorldWidth > 0 {
		w = pen.WorldWidth * s
	}
	if w <= 0 {
		w = p.pl.Thin
	}
	var dash []float64
	off := 0.0
	if len(pen.Dash) > 0 {
		sum := 0.0
		dots := false
		for _, v := range pen.Dash {
			v = math.Abs(v) * s
			dash = append(dash, v)
			sum += v
			if v == 0 {
				dots = true
			}
		}
		// a pattern finer than the line is drawn as a continuous line
		if sum < math.Max(w*2, 0.2) || sum > 1e6 {
			dash = nil
		} else {
			off = pen.DashOffset * s
			if dots {
				pen.Cap = CapRound
			}
		}
	}
	p.setStroke(pen.Color)
	p.setLine(w, pen.Cap, pen.Join)
	p.setDash(dash, off)
	p.obj.StrokePath(p.obj.AddPath(toBDF(it.path, m)))
	p.cv.Drawn = true
}

func (p *plot) fillItem(it Item, m canvas.Matrix) {
	rule := byte(0)
	if it.evenOdd {
		rule = 1
	}
	if g := it.fill.Gradient; g != nil && len(g.Stops) > 0 {
		s := Scale(m)
		a, b := Apply(m, g.P0), Apply(m, g.P1)
		var paint bdf.Paint
		if g.Radial {
			paint = bdf.RadialGradient(f32(a.X), f32(a.Y), f32(g.R0*s), f32(b.X), f32(b.Y), f32(g.R1*s), g.Stops...)
		} else {
			paint = bdf.LinearGradient(f32(a.X), f32(a.Y), f32(b.X), f32(b.Y), g.Stops...)
		}
		p.obj.FillPaint(p.obj.AddPaint(paint))
		p.st.fillSet = false
	} else {
		p.setFill(it.fill.Color)
	}
	p.obj.FillPath(p.obj.AddPath(toBDF(it.path, m)), rule)
	p.cv.Drawn = true
}

func (p *plot) setFill(c bdf.Color) {
	if !p.st.fillSet || p.st.fill != c {
		p.obj.FillColor(c)
		p.st.fill, p.st.fillSet = c, true
	}
}

func (p *plot) setStroke(c bdf.Color) {
	if !p.st.strokeSet || p.st.stroke != c {
		p.obj.StrokeColor(c)
		p.st.stroke, p.st.strokeSet = c, true
	}
}

func (p *plot) setLine(w float64, cap, join byte) {
	if !p.st.lineSet || p.st.width != w || p.st.cap != cap || p.st.join != join {
		p.obj.Line(f32(w), cap, join, 10)
		p.st.width, p.st.cap, p.st.join, p.st.lineSet = w, cap, join, true
	}
}

func (p *plot) setDash(dash []float64, off float64) {
	if p.st.dashSet && slices.Equal(p.st.dash, dash) && p.st.dashOff == off {
		return
	}
	if !p.st.dashSet && dash == nil {
		// the initial state is a continuous line
		return
	}
	seg := make([]float32, len(dash))
	for i, v := range dash {
		seg[i] = f32(v)
	}
	p.obj.Dash(seg, f32(off))
	p.st.dash, p.st.dashOff, p.st.dashSet = dash, off, true
}

func (p *plot) setFont(ref bdf.FontRef, size float64) {
	if !p.st.fontSet || p.st.font != ref || p.st.size != size {
		p.obj.Font(ref, f32(size))
		p.st.font, p.st.size, p.st.fontSet = ref, size, true
	}
}

// run is characters of a text drawn with one face.
type run struct {
	fc   *fontset.Choice
	s    []rune
	x, w float64 // ems
}

func (p *plot) text(t *Text, m canvas.Matrix) {
	// canvas text space (y down) to the page
	tm := m.Mul(t.M).Mul(canvas.Scale(1, -1))
	// the font size is the height of the em square across the baseline;
	// the width along the baseline (width factors, fitted text) goes into
	// the advances, which the renderer stretches the glyphs to, so that
	// the transform only turns, slants and mirrors the text and the
	// positions of the runs stay comparable with their advances
	ex := math.Hypot(tm[0], tm[1])
	det := tm[0]*tm[3] - tm[1]*tm[2]
	if ex < 1e-9 {
		return
	}
	size := math.Abs(det) / ex
	if size < 1e-3 || size > 1e5 || math.IsNaN(size) {
		return
	}
	hx := ex / size // horizontal scale of the glyphs
	norm := canvas.Matrix{tm[0] / ex, tm[1] / ex, tm[2] / size, tm[3] / size, tm[4], tm[5]}
	fonts := p.pl.Fonts
	var runs []run
	x := 0.0
	i := 0
	for _, r := range t.S {
		fc := fonts.face(t.Font, r)
		w := fonts.Set.Advance(fc, r)
		rx := x
		if t.Pos != nil && i < len(t.Pos) {
			rx = t.Pos[i]
		}
		n := len(runs)
		if n > 0 && runs[n-1].fc.Use == fc.Use && (t.Pos == nil || math.Abs(runs[n-1].x+runs[n-1].w-rx) < 1e-6) {
			runs[n-1].s = append(runs[n-1].s, r)
			runs[n-1].w += w
		} else {
			runs = append(runs, run{fc: fc, s: []rune{r}, x: rx, w: w})
		}
		x = rx + w
		i++
	}
	p.save()
	p.cv.Transform(norm)
	switch t.Break {
	case BreakBox:
		p.obj.Mark(bdf.MarkBox, "")
	case BreakParagraph:
		p.obj.Mark(bdf.MarkParagraph, "")
	case BreakLine:
		p.obj.Mark(bdf.MarkLine, "")
	case BreakWrap:
		p.obj.Mark(bdf.MarkWrap, "")
	}
	k := size * hx // ems to units along the baseline
	for _, r := range runs {
		s := string(r.s)
		if strings.TrimSpace(s) == "" {
			continue
		}
		p.setLang(r.s)
		p.setFont(p.cv.Font(r.fc.Use), size)
		p.setFill(t.Color)
		p.obj.FillText(s, f32(r.x*k), 0, f32(r.w*k))
		p.cv.Drawn = true
	}
	if t.Underline || t.Overline || t.Strike {
		fc := fonts.face(t.Font, 'H')
		th := math.Max(fc.ULThick, 0.04) * size
		w := x * k
		p.setFill(t.Color)
		if t.Underline {
			p.obj.FillRect(0, f32(-fc.ULPos*size), f32(w), f32(th))
		}
		if t.Overline {
			p.obj.FillRect(0, f32(-(t.Cap+0.15)*size-th), f32(w), f32(th))
		}
		if t.Strike {
			p.obj.FillRect(0, f32(-t.Cap*0.5*size), f32(w), f32(th))
		}
	}
	p.restore()
}

// setLang sets the language of a run: East Asian text takes the language
// its kana or hangul imply, and otherwise keeps an East Asian language in
// effect or takes the document's (Han characters alone could be Chinese as
// well as Japanese or Korean); other text takes the document's.
func (p *plot) setLang(s []rune) {
	cjk := false
	for _, r := range s {
		if fontdb.IsCJK(r) {
			cjk = true
			break
		}
	}
	lang := ""
	if cjk {
		for _, r := range s {
			if lang = fontset.RuneLang(r); lang != "" {
				break
			}
		}
		if lang == "" && fontset.Script(p.cv.Lang()) != "" {
			return
		}
	}
	p.cv.SetLang(lang)
}

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}
