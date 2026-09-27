package kicad

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// effects is the (effects ...) of a text.
type effects struct {
	sizeX, sizeY float64 // mm
	thickness    float64 // mm; 0: the default
	bold, italic bool
	color        bdf.Color
	hasColor     bool
	face         string // a TrueType font ("" for KiCad's stroke font)
	lineSpacing  float64
	h            hAlign
	v            vAlign
	mirror       bool
	hide         bool
}

// defaultTextSize is the size of text without one (50 mils).
const defaultTextSize = 1.27

func parseEffects(n *node) effects {
	e := effects{sizeX: defaultTextSize, sizeY: defaultTextSize, lineSpacing: 1}
	if n == nil {
		return e
	}
	if f := n.child("font"); f != nil {
		if s := f.child("size"); s != nil {
			e.sizeY, e.sizeX = s.num(0), s.num(1)
		}
		e.thickness = f.numOf("thickness", 0)
		e.bold = f.flag("bold")
		e.italic = f.flag("italic")
		if c, ok := parseColor(f.child("color")); ok {
			e.color, e.hasColor = c, true
		}
		e.face = f.str("face")
		if e.face == "KiCad Font" {
			e.face = ""
		}
		if ls := f.numOf("line_spacing", 0); ls > 0 {
			e.lineSpacing = ls
		}
	}
	if j := n.child("justify"); j != nil {
		for i := 0; ; i++ {
			a := j.arg(i)
			if a == "" {
				break
			}
			switch a {
			case "left":
				e.h = alignLeft
			case "right":
				e.h = alignRight
			case "top":
				e.v = alignTop
			case "bottom":
				e.v = alignBottom
			case "mirror":
				e.mirror = true
			}
		}
	}
	e.hide = n.flag("hide")
	return e
}

// parseColor reads (color r g b a); transparent black (0 0 0 0) means the
// theme's color.
func parseColor(n *node) (bdf.Color, bool) {
	if n == nil {
		return 0, false
	}
	r, g, b, a := n.num(0), n.num(1), n.num(2), n.num(3)
	if r == 0 && g == 0 && b == 0 && a == 0 {
		return 0, false
	}
	u8 := func(v float64) uint8 { return uint8(math.Max(0, math.Min(255, math.Round(v)))) }
	return bdf.RGBA(u8(r), u8(g), u8(b), u8(a*255)), true
}

// placedText is a text to draw: its string (with markup), anchor (mm, y
// down), angle (degrees, counterclockwise) and effects.
type placedText struct {
	s         string
	x, y      float64
	angle     float64
	e         effects
	color     bdf.Color
	thickness float64 // the pen, mm (0: from the effects)
	brk       cad.Break
}

// penWidth returns the width of a text's strokes (EDA_TEXT::
// GetEffectiveTextPenWidth): its thickness, or one made from its size.
func penWidth(e effects, def float64) float64 {
	size := math.Min(math.Abs(e.sizeX), math.Abs(e.sizeY))
	t := e.thickness
	if t <= 0 {
		t = def
		if e.bold {
			t = boldThickness(size)
		} else if t <= 0 {
			t = normalThickness(size)
		}
	}
	return clampThickness(t, size, false)
}

// place returns the transform that puts a line's own space (x along the
// baseline, y down, mm) at its place in the drawing.
func place(x, y, angle float64, mirror bool) canvas.Matrix {
	m := canvas.Translate(x, y).Mul(canvas.Rotate(-angle))
	if mirror {
		m = m.Mul(canvas.Scale(-1, 1))
	}
	return m
}

// text draws a text into d.
func (c *conv) text(d *cad.Drawing, t placedText) {
	s := strings.ReplaceAll(t.s, "\r\n", "\n")
	if strings.TrimSpace(s) == "" || t.e.sizeX <= 0 || t.e.sizeY <= 0 {
		return
	}
	if t.e.face != "" {
		c.outlineText(d, t, s)
		return
	}
	pen := t.thickness
	if pen <= 0 {
		pen = penWidth(t.e, c.defaultPen)
	}
	l := &strokeLayout{sizeX: t.e.sizeX, sizeY: t.e.sizeY, italic: t.e.italic, pen: pen, iu: c.iu}
	b := l.layoutBlock(s, t.e.h, t.e.v, t.e.lineSpacing)
	pm := place(t.x, t.y, t.angle, t.e.mirror)
	for i, lg := range b.lines {
		if strings.TrimSpace(lg.text) == "" {
			continue
		}
		st := &cad.StrokeText{
			S:     lg.text,
			M:     pm.Mul(canvas.Translate(b.pos[i].X, b.pos[i].Y)).Mul(canvas.Scale(1, -1)),
			Pen:   cad.Pen{Color: t.color, WorldWidth: pen, Cap: cad.CapRound, Join: cad.JoinRound},
			Break: t.brk,
		}
		if i > 0 {
			st.Break = cad.BreakLine
		}
		if len(lg.strokes) > 0 {
			p := &cad.Path{}
			for _, s := range lg.strokes {
				for j, q := range s {
					if j == 0 {
						p.MoveTo(q.X, -q.Y)
					} else {
						p.LineTo(q.X, -q.Y)
					}
				}
				if len(s) == 1 {
					// a dot
					p.LineTo(s[0].X, -s[0].Y)
				}
			}
			st.Strokes = p
		}
		for _, part := range lg.parts {
			spec := cad.FontSpec{Family: "sans-serif", Bold: t.e.bold}
			adv := c.fonts.RuneAdvance(spec, part.r)
			x := part.x + (part.cell-adv*part.em)/2
			txt := c.fonts.NewText(spec, string(part.r), t.color,
				canvas.Translate(x, part.dy).Mul(canvas.Scale(part.em, part.em)))
			txt.Break = cad.BreakNone
			st.Parts = append(st.Parts, txt)
		}
		d.StrokeText(st)
	}
}

// outlineSize is how KiCad sizes TrueType text: the em square is the
// text's height times this (OUTLINE_FONT, which scales on the full height
// of ascenders and descenders).
const outlineSize = 1.4

// outlineText draws a text in a TrueType font, its lines placed as KiCad
// places them (FONT::getLinePositions, without the stroke font's pen
// offsets).
func (c *conv) outlineText(d *cad.Drawing, t placedText, s string) {
	spec := cad.FontSpec{Family: t.e.face, EastAsian: "sans-serif", Bold: t.e.bold, Italic: t.e.italic}
	em := t.e.sizeY * outlineSize
	lines := splitLines(s)
	interline := t.e.sizeY * t.e.lineSpacing * 1.68
	height := t.e.sizeY*firstLineHeight + float64(len(lines)-1)*interline
	offY := t.e.sizeY
	switch t.e.v {
	case alignMiddle:
		offY -= height / 2
	case alignBottom:
		offY -= height
	}
	pm := place(t.x, t.y, t.angle, t.e.mirror)
	for i, line := range lines {
		plain := parseMarkup(line).plain()
		if strings.TrimSpace(plain) == "" {
			continue
		}
		w := c.fonts.Measure(spec, plain) * em
		x := 0.0
		switch t.e.h {
		case alignCenter:
			x = -w / 2
		case alignRight:
			x = -w
		}
		y := offY + float64(i)*interline
		txt := c.fonts.NewText(spec, plain, t.color, pm.Mul(canvas.Translate(x, y)).Mul(canvas.Scale(em, -em)))
		txt.Break = t.brk
		if i > 0 {
			txt.Break = cad.BreakLine
		}
		d.Text(txt)
	}
}

// outlineHeight is the height of a line of TrueType text as KiCad measures
// it: the font's ascent and descent.
func (c *conv) outlineHeight(e effects) float64 {
	asc, desc, _ := c.fonts.Metrics(cad.FontSpec{Family: e.face, Bold: e.bold, Italic: e.italic})
	if asc+desc <= 0 {
		asc, desc = 0.9, 0.25
	}
	return (asc + desc) * e.sizeY * outlineSize
}

// textBox returns the box of a text anchored at (0, 0) as KiCad computes it
// (EDA_TEXT::GetTextBox), before it is turned: its top left corner and
// size (mm, y down). The box is of the text's own pen, never the default
// one it may be drawn with.
func (c *conv) textBox(s string, e effects) (x0, y0, w, h float64) {
	pen := penWidth(e, 0)
	if e.face != "" {
		spec := cad.FontSpec{Family: e.face, Bold: e.bold, Italic: e.italic}
		lines := splitLines(s)
		for _, line := range lines {
			w = math.Max(w, c.fonts.Measure(spec, parseMarkup(line).plain())*e.sizeY*outlineSize)
		}
		h = c.outlineHeight(e) + float64(len(lines)-1)*e.sizeY*e.lineSpacing*1.68
		switch e.h {
		case alignCenter:
			x0 = -w / 2
		case alignRight:
			x0 = -w
		}
		switch e.v {
		case alignMiddle:
			y0 = -h / 2
		case alignBottom:
			y0 = -h
		}
		return x0, y0, w, h
	}
	l := &strokeLayout{sizeX: e.sizeX, sizeY: e.sizeY, italic: e.italic, pen: pen, iu: c.iu}
	return l.textBox(s, e.h, e.v, e.lineSpacing, e.mirror)
}
