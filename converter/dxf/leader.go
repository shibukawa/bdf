package dxf

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

// dimVar returns a dimension variable of an entity: its override in the
// extended data (the DSTYLE list of code and value pairs), or the value
// of its dimension style.
func (c *converter) dimVar(e *entity, code int, def float64) float64 {
	xd := e.xdata("ACAD")
	for i := 0; i+1 < len(xd); i++ {
		if xd[i].code == 1070 && int(xd[i].f) == code {
			return xd[i+1].f
		}
	}
	if ds := c.d.dimstyles[key(e.str(3))]; ds != nil {
		return ds.num(code, def)
	}
	if ds := c.d.dimstyles["STANDARD"]; ds != nil {
		return ds.num(code, def)
	}
	return def
}

func (c *converter) dimStr(e *entity, code int) string {
	xd := e.xdata("ACAD")
	for i := 0; i+1 < len(xd); i++ {
		if xd[i].code == 1070 && int(xd[i].f) == code {
			return xd[i+1].s
		}
	}
	if ds := c.d.dimstyles[key(e.str(3))]; ds != nil {
		return ds.str(code)
	}
	return ""
}

// arrowName returns the name of the block of an arrowhead handle ("" for
// the default closed filled arrow).
func (c *converter) arrowName(handle string) string {
	if handle == "" {
		return ""
	}
	if r := c.d.records[strings.ToUpper(handle)]; r != nil {
		return strings.ToUpper(r.name)
	}
	return ""
}

// arrow draws an arrowhead with its tip at tip, pointing away from from.
func (c *converter) arrow(out *cad.Drawing, name string, tip, from cad.Point, size float64, m canvas.Matrix, pen cad.Pen) {
	if size <= 0 {
		return
	}
	d := tip.Sub(from)
	l := d.Len()
	if l == 0 {
		return
	}
	u := d.Mul(1 / l)
	n := cad.Point{X: -u.Y, Y: u.X}
	back := tip.Sub(u.Mul(size))
	pen.Dash = nil
	switch name {
	case "_NONE":
	case "_DOT", "_DOTBLANK":
		circle := (&cad.Path{}).Circle(tip, size/4).Transform(m)
		if name == "_DOT" {
			out.Fill(circle, cad.Fill{Color: pen.Color}, false)
		} else {
			out.Stroke(circle, pen)
		}
	case "_DOTSMALL":
		out.Fill((&cad.Path{}).Circle(tip, size/8).Transform(m), cad.Fill{Color: pen.Color}, false)
	case "_ARCHTICK", "_OBLIQUE":
		v := u.Add(n).Mul(size / 2 / math.Sqrt2)
		a, b := tip.Sub(v), tip.Add(v)
		out.Stroke((&cad.Path{}).MoveTo(a.X, a.Y).LineTo(b.X, b.Y).Transform(m), pen)
	case "_OPEN", "_OPEN30", "_OPEN90":
		w := size / 6
		switch name {
		case "_OPEN30":
			w = size * math.Tan(15*math.Pi/180)
		case "_OPEN90":
			w = size
		}
		a, b := back.Add(n.Mul(w)), back.Sub(n.Mul(w))
		out.Stroke((&cad.Path{}).MoveTo(a.X, a.Y).LineTo(tip.X, tip.Y).LineTo(b.X, b.Y).Transform(m), pen)
	case "_CLOSEDBLANK", "_CLOSED":
		a, b := back.Add(n.Mul(size/6)), back.Sub(n.Mul(size/6))
		out.Stroke((&cad.Path{}).Polyline([]cad.Point{tip, a, b}, true).Transform(m), pen)
	default: // closed filled
		a, b := back.Add(n.Mul(size/6)), back.Sub(n.Mul(size/6))
		out.Fill((&cad.Path{}).Polyline([]cad.Point{tip, a, b}, true).Transform(m), cad.Fill{Color: pen.Color}, false)
	}
}

func (c *converter) leader(e *entity, x *ctx, p props) {
	var pts []cad.Point
	for _, t := range e.tags {
		switch t.code {
		case 10:
			pts = append(pts, cad.Point{X: t.f})
		case 20:
			if n := len(pts); n > 0 {
				pts[n-1].Y = t.f
			}
		}
	}
	if len(pts) < 2 {
		return
	}
	pen := c.pen(e, x, p, x.m)
	path := &cad.Path{}
	if e.int(72, 0) == 1 {
		path.Interpolate(pts, nil, nil)
	} else {
		path.Polyline(pts, false)
	}
	x.out.Stroke(path.Transform(x.m), pen)
	if e.int(71, 1) != 0 {
		size := c.dimVar(e, 41, 2.5) * math.Max(c.dimVar(e, 40, 1), 0)
		if size == 0 {
			size = c.dimVar(e, 41, 2.5)
		}
		c.arrow(x.out, c.arrowName(c.dimStr(e, 341)), pts[0], pts[1], size, x.m, pen)
	}
}

// rawColor converts the raw colors of MULTILEADER (0xC2RRGGBB for a true
// color, 0xC3 with an index, 0xC0 by layer, 0xC1 by block).
func (c *converter) rawColor(v int, def bdf.Color, dark bool) bdf.Color {
	switch uint32(v) >> 24 {
	case 0xc2:
		return rgb(v & 0xffffff)
	case 0xc3:
		return c.aci(v&0xff, dark)
	}
	return def
}

type mleaderLeader struct {
	last    cad.Point
	dogleg  cad.Point
	doglen  float64
	lines   [][]cad.Point
	hasLast bool
}

func (c *converter) multileader(e *entity, x *ctx, p props) {
	var (
		leaders   []*mleaderLeader
		cur       *mleaderLeader
		line      *[]cad.Point
		text      string
		textPos   cad.Point
		textDir   = cad.Point{X: 1}
		textH     = 2.5
		textW     float64
		attach    = 1
		textColor = p.color
		arrowSize = -1.0
		hasText   bool
		blockID   string
		blockPos  cad.Point
		blockSc   = cad.Point{X: 1, Y: 1}
		blockRot  float64
		styleID   string
		depth     int
	)
	lineType := 1
	arrowID := ""
	for _, t := range e.tags {
		if t.code >= 1000 {
			break
		}
		switch t.code {
		case 300, 302, 304:
			switch t.s {
			case "CONTEXT_DATA{":
				depth = 1
			case "LEADER{":
				depth = 2
				cur = &mleaderLeader{}
				leaders = append(leaders, cur)
			case "LEADER_LINE{":
				depth = 3
				cur.lines = append(cur.lines, nil)
				line = &cur.lines[len(cur.lines)-1]
			default:
				if t.code == 304 && depth == 1 {
					text = t.s
				}
			}
			continue
		case 301:
			depth = 0
			continue
		case 303:
			depth = 1
			continue
		case 305:
			depth = 2
			continue
		}
		switch depth {
		case 3:
			switch t.code {
			case 10:
				*line = append(*line, cad.Point{X: t.f})
			case 20:
				if n := len(*line); n > 0 {
					(*line)[n-1].Y = t.f
				}
			}
		case 2:
			switch t.code {
			case 10:
				cur.last.X, cur.hasLast = t.f, true
			case 20:
				cur.last.Y = t.f
			case 11:
				cur.dogleg.X = t.f
			case 21:
				cur.dogleg.Y = t.f
			case 40:
				cur.doglen = t.f
			}
		case 1:
			switch t.code {
			case 41:
				textH = t.f
			case 140:
				arrowSize = t.f
			case 290:
				hasText = t.f != 0
			case 12:
				textPos.X = t.f
			case 22:
				textPos.Y = t.f
			case 13:
				textDir.X = t.f
			case 23:
				textDir.Y = t.f
			case 43:
				textW = t.f
			case 90:
				textColor = c.rawColor(int(t.f), p.color, x.dark)
			case 171:
				attach = int(t.f)
			case 340:
				styleID = strings.ToUpper(t.s)
			case 341:
				blockID = strings.ToUpper(t.s)
			case 15:
				blockPos.X = t.f
			case 25:
				blockPos.Y = t.f
			case 16:
				blockSc.X = t.f
			case 26:
				blockSc.Y = t.f
			case 46:
				blockRot = t.f
			}
		case 0:
			switch t.code {
			case 170:
				lineType = int(t.f)
			case 342:
				arrowID = t.s
			case 42:
				if arrowSize < 0 {
					arrowSize = t.f
				}
			}
		}
	}
	pen := c.pen(e, x, p, x.m)
	if arrowSize < 0 {
		arrowSize = 2.5
	}
	for _, l := range leaders {
		end := l.last
		for _, ln := range l.lines {
			if len(ln) == 0 || lineType == 0 {
				continue
			}
			pts := append(append([]cad.Point(nil), ln...), end)
			path := &cad.Path{}
			if lineType == 2 {
				path.Interpolate(pts, nil, nil)
			} else {
				path.Polyline(pts, false)
			}
			x.out.Stroke(path.Transform(x.m), pen)
			next := end
			if len(ln) > 1 {
				next = ln[1]
			}
			c.arrow(x.out, c.arrowName(arrowID), ln[0], next, arrowSize, x.m, pen)
		}
		if l.hasLast && l.doglen > 0 && lineType != 0 {
			d := end.Add(l.dogleg.Mul(l.doglen))
			x.out.Stroke((&cad.Path{}).MoveTo(end.X, end.Y).LineTo(d.X, d.Y).Transform(x.m), pen)
		}
	}
	if hasText && text != "" {
		style := ""
		if st := c.styleByHandle(styleID); st != nil {
			style = st.name
		}
		mt := &entity{typ: "MTEXT", tags: []tag{
			{code: 7, s: style}, {code: 10, f: textPos.X}, {code: 20, f: textPos.Y},
			{code: 11, f: textDir.X}, {code: 21, f: textDir.Y}, {code: 40, f: textH},
			{code: 41, f: textW}, {code: 71, f: float64(attach)}, {code: 1, s: text},
		}}
		tp := p
		tp.color = textColor
		c.mtext(mt, x, tp)
	}
	if blockID != "" {
		if r := c.d.records[blockID]; r != nil {
			if b := c.d.blocks[key(r.name)]; b != nil {
				m := x.m.Mul(canvas.Translate(blockPos.X, blockPos.Y)).Mul(canvas.Rotate(blockRot * 180 / math.Pi)).
					Mul(canvas.Scale(blockSc.X, blockSc.Y)).Mul(canvas.Translate(-b.base.X, -b.base.Y))
				c.block(b, x, m, p, c.layer(e.str(8)))
			}
		}
	}
}

func (c *converter) styleByHandle(h string) *style {
	if h == "" {
		return nil
	}
	return c.d.styleHandles[h]
}

// mline draws a multiline: each element of its style is a polyline
// through the vertices, offset along their miter directions.
func (c *converter) mline(e *entity, x *ctx, p props) {
	type mvertex struct {
		p, miter cad.Point
		params   [][]float64
	}
	var vs []*mvertex
	var cur *mvertex
	var elem int
	for _, t := range e.tags {
		switch t.code {
		case 11:
			cur = &mvertex{p: cad.Point{X: t.f}}
			vs = append(vs, cur)
			elem = -1
		case 21:
			if cur != nil {
				cur.p.Y = t.f
			}
		case 13:
			if cur != nil {
				cur.miter.X = t.f
			}
		case 23:
			if cur != nil {
				cur.miter.Y = t.f
			}
		case 74:
			if cur != nil {
				cur.params = append(cur.params, nil)
				elem = len(cur.params) - 1
			}
		case 41:
			if cur != nil && elem >= 0 {
				cur.params[elem] = append(cur.params[elem], t.f)
			}
		}
	}
	if len(vs) < 2 {
		return
	}
	n := e.int(73, len(vs[0].params))
	closed := e.int(71, 0)&2 != 0
	colors := c.mlineColors(e, x.dark)
	for j := range n {
		var pts []cad.Point
		for _, v := range vs {
			if j >= len(v.params) || len(v.params[j]) == 0 {
				continue
			}
			pts = append(pts, v.p.Add(v.miter.Mul(v.params[j][0])))
		}
		if len(pts) < 2 {
			continue
		}
		pen := c.pen(e, x, p, x.m)
		if j < len(colors) && colors[j] != 0 {
			pen.Color = colors[j]
		}
		x.out.Stroke((&cad.Path{}).Polyline(pts, closed).Transform(x.m), pen)
	}
}

// mlineColors returns the colors of the elements of a multiline's style
// (0 for BYLAYER and BYBLOCK).
func (c *converter) mlineColors(e *entity, dark bool) []bdf.Color {
	st := c.d.objects[strings.ToUpper(e.str(340))]
	if st == nil {
		return nil
	}
	var out []bdf.Color
	inElems := false
	for _, t := range st.tags {
		switch t.code {
		case 71:
			inElems = true
		case 62:
			if inElems {
				v := int(t.f)
				if v == 0 || v == 256 {
					out = append(out, 0)
				} else {
					out = append(out, c.aci(v, dark))
				}
			}
		case 420:
			if inElems && len(out) > 0 {
				out[len(out)-1] = rgb(int(t.f))
			}
		}
	}
	return out
}
