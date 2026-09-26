package sxf

import (
	"math"
	"slices"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/width"
)

const mm = 72 / 25.4

// renderer draws into a drawing in sheet coordinates (mm, y up).
type renderer struct {
	fonts  *cad.Fonts
	out    *cad.Drawing
	light  bool // on white paper: white is drawn black
	bg     bdf.Color
	warn   func(key, format string, args ...any)
	count  int // elements drawn
	points int // their points
	visits int // elements looked at
}

// The budget: maxItems bounds the elements drawn and maxPoints their
// points, compound figures included, and maxVisits the elements looked at
// (compound figures placed in compound figures multiply them).
const (
	maxItems  = 5_000_000
	maxPoints = 10_000_000
	maxVisits = 4 * maxItems
)

func (r *renderer) spent() bool {
	if r.count > maxItems || r.points > maxPoints || r.visits > maxVisits {
		r.warn("budget", "the drawing is too large; only its first %d elements or %d points are drawn", maxItems, maxPoints)
		return true
	}
	return false
}

// visit counts an element looked at and reports whether to go on.
func (r *renderer) visit() bool {
	r.visits++
	return !r.spent()
}

// budget counts an element of n points to draw and reports whether to
// draw it.
func (r *renderer) budget(n int) bool {
	r.count++
	r.points += n
	return !r.spent()
}

// lineStyle is how a curve is drawn: color, width (mm) and dash pattern
// (mm on the sheet, positive dashes and negative gaps).
type lineStyle struct {
	color bdf.Color
	width float64
	dash  []float64
}

func (r *renderer) paint(c bdf.Color) bdf.Color {
	if r.light && c == bdf.RGB(255, 255, 255) {
		return bdf.RGB(0, 0, 0)
	}
	return c
}

func (r *renderer) pen(s lineStyle) cad.Pen {
	p := cad.Pen{Color: r.paint(s.color), Width: s.width * mm, Cap: cad.CapRound, Join: cad.JoinRound}
	p.Dash, p.DashOffset = cad.DashPattern(s.dash)
	return p
}

// stroke draws a path given in the coordinates that m maps to the sheet.
func (r *renderer) stroke(p *cad.Path, s lineStyle, m canvas.Matrix) {
	if p.Empty() || !r.budget(p.Size()) {
		return
	}
	r.out.Stroke(p.Transform(m), r.pen(s))
}

// textBox is a text in its box: the box is w × h with the anchor at one of
// nine points (1 bottom left ... 9 top right, as SXF numbers them), turned
// by angle and slanted by slant (degrees, positive to the right); vertical
// text runs down the box.
type textBox struct {
	s             string
	font          string
	color         bdf.Color
	x, y          float64
	w, h, spacing float64
	angle, slant  float64
	anchor        int
	vertical      bool
	onBaseline    bool // x, y is on the baseline, not at an anchor (P21 files not written for SXF)
}

var sjisEnc = japanese.ShiftJIS.NewEncoder()

// units returns the share of a character in the width of a text: 2 for
// the full-width characters (two bytes in Shift_JIS), 1 for the others.
func units(c rune) float64 {
	if c < 0x80 {
		return 1
	}
	if b, err := sjisEnc.String(string(c)); err == nil {
		if len(b) == 2 {
			return 2
		}
		return 1
	}
	if k := width.LookupRune(c).Kind(); k == width.EastAsianWide || k == width.EastAsianFullwidth {
		return 2
	}
	return 1
}

// text draws a text box given in the coordinates that m maps to the sheet.
func (r *renderer) text(t textBox, m canvas.Matrix) {
	s := strings.TrimRight(t.s, " ")
	if strings.TrimSpace(s) == "" || t.h <= 0 || !r.budget(len(s)) {
		return
	}
	runes := []rune(s)
	n := float64(len(runes))
	total := 0.0
	for _, c := range runes {
		total += units(c)
	}
	// the frame of the text: at its anchor, turned, not mirrored
	a := t.angle * math.Pi / 180
	f := m.Mul(canvas.Translate(t.x, t.y)).Mul(canvas.Matrix{math.Cos(a), math.Sin(a), -math.Sin(a), math.Cos(a), 0, 0})
	if f[0]*f[3]-f[1]*f[2] < 0 {
		// a mirrored frame (geodetic coordinates): the text still reads
		// along its direction
		f = f.Mul(canvas.Scale(1, -1))
	}
	skew := canvas.Matrix{1, 0, math.Tan(t.slant * math.Pi / 180), 1, 0, 0}
	col := r.paint(t.color)
	spec := cad.FontSpec{Family: t.font}
	if !t.vertical {
		w := t.w
		cellsW := w - t.spacing*(n-1)
		if cellsW <= 0 {
			cellsW = w
		}
		unit := cellsW / total
		// the box relative to the anchor
		var bx, by float64
		switch (t.anchor - 1) % 3 {
		case 1:
			bx = -w / 2
		case 2:
			bx = -w
		}
		em := t.h
		switch {
		case t.onBaseline:
			by = -0.12 * em
		case (t.anchor-1)/3 == 1:
			by = -t.h / 2
		case (t.anchor-1)/3 == 2:
			by = -t.h
		}
		cells := make([]float64, len(runes))
		for i, c := range runes {
			cells[i] = units(c) * unit / em
		}
		tm := f.Mul(canvas.Translate(bx, by+0.12*em)).Mul(skew).Mul(canvas.Scale(em, em))
		r.out.Text(r.fonts.NewCellText(spec, s, cells, t.spacing/em, col, tm))
		return
	}
	// vertical: the column is w wide and h long; characters stand upright
	// down its middle
	h := t.h
	cellsH := h - t.spacing*(n-1)
	if cellsH <= 0 {
		cellsH = h
	}
	unit := cellsH / total
	var bx, by float64
	switch (t.anchor - 1) % 3 {
	case 1:
		bx = -t.w / 2
	case 2:
		bx = -t.w
	}
	switch {
	case t.onBaseline:
		// the column starts at the point
	case (t.anchor-1)/3 == 0:
		by = h
	case (t.anchor-1)/3 == 1:
		by = h / 2
	}
	em := t.w
	cells := make([]float64, len(runes))
	for i, c := range runes {
		cells[i] = units(c) * unit / em
	}
	tm := f.Mul(canvas.Translate(bx+t.w/2, by)).Mul(skew).Mul(canvas.Rotate(-90)).Mul(canvas.Scale(em, em))
	tx := r.fonts.NewCellText(spec, s, cells, t.spacing/em, col, tm)
	tx.Vertical = true
	r.out.Text(tx)
}

// arrowNames are the names of the arrows (terminator symbols), by code.
var arrowNames = []string{"blanked arrow", "blanked box", "blanked dot", "dimension origin", "filled box", "filled arrow",
	"filled dot", "integral symbol", "open arrow", "slash", "unfilled arrow"}

func arrowCode(name string) int {
	return slices.Index(arrowNames, strings.ToLower(strings.TrimSpace(name))) + 1
}

// arrow draws an arrow with its tip at tip pointing along dir, both in the
// coordinates that m maps to the sheet; its size (2.5 mm times scale) is
// on the sheet.
func (r *renderer) arrow(code int, tip, dir cad.Point, scale float64, s lineStyle, m canvas.Matrix) {
	if code <= 0 || scale <= 0 || !r.budget(16) {
		return
	}
	t := cad.Apply(m, tip)
	d := cad.Apply(m, tip.Add(dir)).Sub(t)
	l := d.Len()
	if l == 0 {
		return
	}
	u := d.Mul(1 / l)
	n := cad.Point{X: -u.Y, Y: u.X}
	size := 2.5 * scale
	pen := r.pen(lineStyle{color: s.color, width: s.width})
	fill := cad.Fill{Color: r.paint(s.color)}
	back := t.Sub(u.Mul(size))
	half := size * math.Tan(15*math.Pi/180)
	tri := (&cad.Path{}).Polyline([]cad.Point{t, back.Add(n.Mul(half)), back.Sub(n.Mul(half))}, true)
	box := (&cad.Path{}).Polyline([]cad.Point{t.Add(n.Mul(size / 2)), t.Sub(n.Mul(size / 2)), t.Sub(n.Mul(size / 2)).Sub(u.Mul(size)),
		t.Add(n.Mul(size / 2)).Sub(u.Mul(size))}, true)
	dot := (&cad.Path{}).Circle(t, size/2)
	switch code {
	case 1: // blanked arrow: an outline over the background
		r.out.Fill(tri, cad.Fill{Color: r.bg}, false)
		r.out.Stroke(tri, pen)
	case 2:
		r.out.Fill(box, cad.Fill{Color: r.bg}, false)
		r.out.Stroke(box, pen)
	case 3, 4:
		r.out.Fill(dot, cad.Fill{Color: r.bg}, false)
		r.out.Stroke(dot, pen)
		if code == 4 {
			r.out.Fill((&cad.Path{}).Circle(t, size/6), fill, false)
		}
	case 5:
		r.out.Fill(box, fill, false)
	case 6:
		r.out.Fill(tri, fill, false)
	case 7:
		r.out.Fill(dot, fill, false)
	case 8: // integral symbol: an S across the line
		p := &cad.Path{}
		a := t.Add(n.Mul(size * 1.5)).Sub(u.Mul(size / 2))
		b := t.Sub(n.Mul(size * 1.5)).Add(u.Mul(size / 2))
		p.MoveTo(a.X, a.Y).CubicTo(t.Add(n.Mul(size*1.5)).Add(u.Mul(size)), t.Sub(n.Mul(size*1.5)).Sub(u.Mul(size)), b)
		r.out.Stroke(p, pen)
	case 9:
		p := (&cad.Path{}).MoveTo(back.X+n.X*half, back.Y+n.Y*half).LineTo(t.X, t.Y).LineTo(back.X-n.X*half, back.Y-n.Y*half)
		r.out.Stroke(p, pen)
	case 10:
		v := u.Add(n).Mul(size / 2)
		p := (&cad.Path{}).MoveTo(t.X-v.X, t.Y-v.Y).LineTo(t.X+v.X, t.Y+v.Y)
		r.out.Stroke(p, pen)
	case 11:
		r.out.Stroke(tri, pen)
	}
}

// markerNames are the point markers, by code.
var markerNames = []string{"asterisk", "circle", "dot", "plus", "square", "triangle", "x"}

// marker draws a point marker at p turned by angle (degrees), both in the
// coordinates that m maps to the sheet, 2.5 mm times scale on the sheet.
func (r *renderer) marker(code int, p cad.Point, angle, scale float64, c bdf.Color, m canvas.Matrix) {
	if !r.budget(16) {
		return
	}
	q := cad.Apply(m, p)
	h := 1.25 * math.Max(scale, 0)
	if h == 0 {
		return
	}
	// the direction and handedness through m
	a := angle * math.Pi / 180
	d := cad.Apply(m, p.Add(cad.Point{X: math.Cos(a), Y: math.Sin(a)})).Sub(q)
	rot := canvas.Translate(q.X, q.Y).Mul(canvas.Rotate(math.Atan2(d.Y, d.X) * 180 / math.Pi)).Mul(canvas.Scale(h, h))
	if m[0]*m[3]-m[1]*m[2] < 0 {
		rot = rot.Mul(canvas.Scale(1, -1))
	}
	pen := r.pen(lineStyle{color: c, width: 0.13})
	path := &cad.Path{}
	switch code {
	case 1: // asterisk
		for i := range 3 {
			a := float64(i) * math.Pi / 3
			path.MoveTo(-math.Cos(a), -math.Sin(a)).LineTo(math.Cos(a), math.Sin(a))
		}
	case 2:
		path.Circle(cad.Point{}, 1)
	case 3:
		r.out.Fill((&cad.Path{}).Circle(cad.Point{}, 0.3).Transform(rot), cad.Fill{Color: r.paint(c)}, false)
		return
	case 4:
		path.MoveTo(-1, 0).LineTo(1, 0).MoveTo(0, -1).LineTo(0, 1)
	case 5:
		path.Polyline([]cad.Point{{X: -1, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: 1}, {X: -1, Y: 1}}, true)
	case 6:
		path.Polyline([]cad.Point{{X: 0, Y: 1}, {X: -math.Sqrt(3) / 2, Y: -0.5}, {X: math.Sqrt(3) / 2, Y: -0.5}}, true)
	default: // 7: x
		path.MoveTo(-1, -1).LineTo(1, 1).MoveTo(-1, 1).LineTo(1, -1)
	}
	r.out.Stroke(path.Transform(rot), pen)
}

// hatchLine is a family of hatch lines: through a point at an angle,
// spaced apart, in a line style.
type hatchLine struct {
	style        lineStyle
	x, y         float64
	spacing      float64
	angle        float64 // degrees
	patternShift float64
}

// hatch draws hatch lines inside a boundary given in the coordinates that
// m maps to the sheet (even-odd, the holes included in the path).
func (r *renderer) hatch(boundary *cad.Path, lines []hatchLine, m canvas.Matrix) {
	if boundary.Empty() || !r.budget(boundary.Size()) {
		return
	}
	b := boundary.Transform(m)
	bounds := b.Bounds()
	for _, l := range lines {
		if l.spacing == 0 {
			continue
		}
		a := l.angle * math.Pi / 180
		dir := cad.Point{X: math.Cos(a), Y: math.Sin(a)}
		off := cad.Point{X: -dir.Y * l.spacing, Y: dir.X * l.spacing}
		// through m: the angle, spacing and base point on the sheet
		base := cad.Apply(m, cad.Point{X: l.x, Y: l.y})
		md := cad.Apply(m, cad.Point{X: l.x, Y: l.y}.Add(dir)).Sub(base)
		mo := cad.Apply(m, cad.Point{X: l.x, Y: l.y}.Add(off)).Sub(base)
		// the lines across the outline, each clipped into a few pieces
		if n := math.Hypot(bounds.W(), bounds.H()) / math.Max(mo.Len(), 1e-9); !r.budget(4 * int(math.Min(n, cad.MaxHatchLines))) {
			return
		}
		pen := r.pen(l.style)
		pl := cad.PatternLine{Angle: math.Atan2(md.Y, md.X), Base: base, Offset: mo, Dash: l.style.dash}
		pen.Dash, pen.DashOffset = nil, 0
		if !r.out.Hatch(b, []cad.PatternLine{pl}, pen) {
			r.warn("densehatch", "hatches too dense to draw are filled with a tint of their color")
			r.out.Fill(b, cad.Fill{Color: withAlpha(pen.Color, 0.3)}, true)
		}
	}
}

func withAlpha(c bdf.Color, a float64) bdf.Color {
	return c&^0xff | bdf.Color(uint8(math.Round(a*255)))
}

// clothoid returns the points of a clothoid of parameter a from arc length
// s0 to s1, starting at the origin along the x axis (turning to the left;
// cw mirrors it).
func clothoid(a, s0, s1 float64, cw bool) []cad.Point {
	if a <= 0 || s1 <= s0 {
		return nil
	}
	steps := int(math.Min(math.Max((s1-s0)/a*64, 16), 2048))
	var pts []cad.Point
	// integrate cos and sin of s²/(2a²) with Simpson's rule per step
	x, y := 0.0, 0.0
	f := func(s float64) (float64, float64) {
		t := s * s / (2 * a * a)
		return math.Cos(t), math.Sin(t)
	}
	// the part before s0 moves the start
	if s0 > 0 {
		n := int(math.Min(math.Max(s0/a*64, 16), 2048))
		h := s0 / float64(n)
		for i := range n {
			sa, sb := float64(i)*h, float64(i+1)*h
			ca, ya := f(sa)
			cm, ym := f((sa + sb) / 2)
			cb, yb := f(sb)
			x += h / 6 * (ca + 4*cm + cb)
			y += h / 6 * (ya + 4*ym + yb)
		}
	}
	pts = append(pts, cad.Point{X: x, Y: y})
	h := (s1 - s0) / float64(steps)
	for i := range steps {
		sa, sb := s0+float64(i)*h, s0+float64(i+1)*h
		ca, ya := f(sa)
		cm, ym := f((sa + sb) / 2)
		cb, yb := f(sb)
		x += h / 6 * (ca + 4*cm + cb)
		y += h / 6 * (ya + 4*ym + yb)
		pts = append(pts, cad.Point{X: x, Y: y})
	}
	if cw {
		for i := range pts {
			pts[i].Y = -pts[i].Y
		}
	}
	return pts
}
