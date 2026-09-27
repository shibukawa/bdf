package equation

import (
	"math"

	"github.com/shibukawa/bdf"
)

// Without a formula font, delimiters, radical signs and braces that grow
// are drawn as strokes: a text font's glyph scaled up would thicken with
// it, and text fonts often lack these characters altogether.

// pathDelimiters are the delimiters drawn as strokes.
const pathDelimiters = "()[]{}|‖⟨⟩⌊⌋⌈⌉/"

// strokeWidth is the width of the strokes, which grows a little with
// their length.
func (e *Engine) strokeWidth(v env, length float64) float64 {
	size := e.size(v)
	return math.Min(0.05*size+0.01*length, 0.09*size)
}

// pathDelimiter draws a delimiter total points tall with its bottom on the
// baseline.
func (e *Engine) pathDelimiter(r rune, total float64, v env) *Box {
	size := e.size(v)
	h := total
	w := math.Min(0.25*size+0.08*h, 0.6*size)
	lw := e.strokeWidth(v, h)
	p := &bdf.Path{}
	pt := func(x, y float64) (float32, float32) { return f32(x), f32(y) }
	move := func(x, y float64) { p.MoveTo(pt(x, y)) }
	line := func(x, y float64) { p.LineTo(pt(x, y)) }
	curve := func(x1, y1, x2, y2, x, y float64) {
		a, b := pt(x1, y1)
		c, d := pt(x2, y2)
		ex, ey := pt(x, y)
		p.CubicTo(a, b, c, d, ex, ey)
	}
	// mirror flips a shape drawn for the opening form
	mirror := false
	switch r {
	case ')', ']', '}', '⟩', '⌋', '⌉':
		mirror = true
	}
	x := func(f float64) float64 {
		if mirror {
			return w * (1 - f)
		}
		return w * f
	}
	switch r {
	case '(', ')':
		move(x(0.85), -h)
		curve(x(0.1), -h*0.78, x(0.1), -h*0.22, x(0.85), 0)
	case '[', ']':
		move(x(0.85), -h)
		line(x(0.3), -h)
		line(x(0.3), 0)
		line(x(0.85), 0)
	case '{', '}':
		move(x(0.9), -h)
		curve(x(0.45), -h, x(0.55), -h*0.5, x(0.1), -h*0.5)
		curve(x(0.55), -h*0.5, x(0.45), 0, x(0.9), 0)
	case '|':
		move(w/2, -h)
		line(w/2, 0)
	case '‖':
		move(w*0.3, -h)
		line(w*0.3, 0)
		move(w*0.7, -h)
		line(w*0.7, 0)
	case '⟨', '⟩':
		move(x(0.85), -h)
		line(x(0.15), -h/2)
		line(x(0.85), 0)
	case '⌊', '⌋':
		move(x(0.3), -h)
		line(x(0.3), 0)
		line(x(0.85), 0)
	case '⌈', '⌉':
		move(x(0.85), -h)
		line(x(0.3), -h)
		line(x(0.3), 0)
	case '/':
		move(w*0.9, -h)
		line(w*0.1, 0)
	}
	b := newBox()
	b.stroke(p, 0, 0, lw, v.color, [4]float64{0, -h, w, 0})
	b.W, b.H = w, h
	return b
}

// pathRadical draws a radical sign whose top (where the rule over the
// radicand starts) is top points above the baseline and whose foot is
// bottom points below it.
func (e *Engine) pathRadical(top, bottom float64, v env) *Box {
	size := e.size(v)
	h := top + bottom
	w := math.Min(0.45*size+0.05*h, 0.8*size)
	lw := e.strokeWidth(v, h)
	p := (&bdf.Path{}).MoveTo(0, f32(bottom-0.45*size)).LineTo(f32(0.15*size), f32(bottom-0.52*size)).
		LineTo(f32(0.38*size), f32(bottom)).LineTo(f32(w), f32(-top+lw/2))
	b := newBox()
	b.stroke(p, 0, 0, lw, v.color, [4]float64{0, -top, w, bottom})
	b.W, b.H, b.D = w, top, bottom
	return b
}

// pathBrace draws a horizontal brace (over or under) width points wide
// sitting on the baseline.
func (e *Engine) pathBrace(r rune, width float64, v env) *Box {
	size := e.size(v)
	hb := math.Min(0.3*size, 0.1*size+0.05*width)
	lw := e.strokeWidth(v, width/4)
	y := func(f float64) float32 {
		if r == '⏟' || r == '⏝' || r == '⎵' {
			return f32(-hb + f*hb) // opening downwards
		}
		return f32(-f * hb)
	}
	w := width
	p := &bdf.Path{}
	switch r {
	case '⏞', '⏟':
		p.MoveTo(0, y(0))
		p.CubicTo(0, y(0.5), f32(hb*0.2), y(0.5), f32(hb*0.6), y(0.5))
		p.LineTo(f32(w/2-hb*0.6), y(0.5))
		p.CubicTo(f32(w/2-hb*0.2), y(0.5), f32(w/2), y(0.5), f32(w/2), y(1))
		p.CubicTo(f32(w/2), y(0.5), f32(w/2+hb*0.2), y(0.5), f32(w/2+hb*0.6), y(0.5))
		p.LineTo(f32(w-hb*0.6), y(0.5))
		p.CubicTo(f32(w-hb*0.2), y(0.5), f32(w), y(0.5), f32(w), y(0))
	case '⏜', '⏝':
		p.MoveTo(0, y(0))
		p.CubicTo(f32(w*0.2), y(1), f32(w*0.8), y(1), f32(w), y(0))
	default: // ⎴ ⎵
		p.MoveTo(0, y(0))
		p.LineTo(0, y(0.7))
		p.LineTo(f32(w), y(0.7))
		p.LineTo(f32(w), y(0))
	}
	b := newBox()
	b.stroke(p, 0, 0, lw, v.color, [4]float64{0, -hb, w, 0})
	b.W, b.H = w, hb
	return b
}
