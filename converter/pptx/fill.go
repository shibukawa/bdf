package pptx

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"

	"github.com/shibukawa/bdf"
)

// Fill kinds.
const (
	fillNone = iota
	fillSolid
	fillGrad
	fillBlip
	fillPatt
)

// fill is a resolved DrawingML fill.
type fill struct {
	kind  int
	color rgba
	stops []gradStop // gradient
	grad  *node      // a:gradFill (angle, path)
	blip  *node      // a:blipFill
	patt  string     // preset pattern name
	fg    rgba       // pattern foreground
	bg    rgba       // pattern background
	part  string     // part the blip's relationship belongs to
	cc    *colorCtx  // colors of the blip's effects (phClr of style references)
}

type gradStop struct {
	pos float64
	c   rgba
}

// fillElem returns the fill choice among n's children.
func fillElem(n *node) *node {
	if n == nil {
		return nil
	}
	for _, k := range n.Kids {
		switch k.Name {
		case "noFill", "solidFill", "gradFill", "blipFill", "pattFill", "grpFill":
			return k
		}
	}
	return nil
}

// resolveFill evaluates a fill choice element.
func (s *slideCtx) resolveFill(n *node, part string, cc *colorCtx, grp *groupCtx) fill {
	if n == nil {
		return fill{}
	}
	switch n.Name {
	case "solidFill":
		if c, ok := cc.color(n); ok {
			return fill{kind: fillSolid, color: c}
		}
		return fill{}
	case "gradFill":
		f := fill{kind: fillGrad, grad: n}
		for _, gs := range n.path("gsLst").children("gs") {
			if c, ok := cc.color(gs); ok {
				f.stops = append(f.stops, gradStop{pct(gs, "pos", 0), c})
			}
		}
		if len(f.stops) == 0 {
			return fill{}
		}
		sortStops(f.stops)
		if len(f.stops) == 1 {
			return fill{kind: fillSolid, color: f.stops[0].c}
		}
		return f
	case "blipFill":
		return fill{kind: fillBlip, blip: n, part: part, cc: cc}
	case "pattFill":
		f := fill{kind: fillPatt, patt: n.attrStr("prst", "pct5"), fg: black, bg: white}
		if c, ok := cc.color(n.child("fgClr")); ok {
			f.fg = c
		}
		if c, ok := cc.color(n.child("bgClr")); ok {
			f.bg = c
		}
		return f
	case "grpFill":
		for g := grp; g != nil; g = g.parent {
			if e := fillElem(g.spPr); e != nil && e.Name != "grpFill" {
				return s.resolveFill(e, g.part, cc, g.parent)
			}
		}
	}
	return fill{}
}

func sortStops(st []gradStop) {
	for i := 1; i < len(st); i++ {
		for j := i; j > 0 && st[j].pos < st[j-1].pos; j-- {
			st[j], st[j-1] = st[j-1], st[j]
		}
	}
}

// styleFill resolves a p:style fillRef (or a background bgRef).
func (s *slideCtx) styleFill(ref *node, cc *colorCtx) fill {
	if ref == nil {
		return fill{}
	}
	idx := int(ref.attrInt("idx", 0))
	phc, ok := cc.color(ref)
	if ok {
		cc = cc.withPh(phc)
	}
	var e *node
	switch {
	case idx >= 1001:
		e = styleEntry(s.th.bgFills, idx-1000)
	case idx >= 1:
		e = styleEntry(s.th.fills, idx)
	default:
		return fill{}
	}
	if e == nil {
		if ok {
			return fill{kind: fillSolid, color: phc}
		}
		return fill{}
	}
	return s.resolveFill(e, s.th.part, cc, nil)
}

// adjust applies a path's lighten/darken fill mode.
func (f fill) adjust(mode string) fill {
	var mod func(rgba) rgba
	switch mode {
	case "darken":
		mod = func(c rgba) rgba { return shadeRGB(c, 0.6) }
	case "darkenLess":
		mod = func(c rgba) rgba { return shadeRGB(c, 0.8) }
	case "lighten":
		mod = func(c rgba) rgba { return tintRGB(c, 0.6) }
	case "lightenLess":
		mod = func(c rgba) rgba { return tintRGB(c, 0.8) }
	default:
		return f
	}
	switch f.kind {
	case fillSolid:
		f.color = mod(f.color)
	case fillGrad:
		st := make([]gradStop, len(f.stops))
		for i, s := range f.stops {
			st[i] = gradStop{s.pos, mod(s.c)}
		}
		f.stops = st
	case fillPatt:
		f.fg, f.bg = mod(f.fg), mod(f.bg)
	}
	return f
}

func shadeRGB(c rgba, v float64) rgba {
	return rgba{fromLinear(toLinear(c.R) * v), fromLinear(toLinear(c.G) * v), fromLinear(toLinear(c.B) * v), c.A}
}

func tintRGB(c rgba, v float64) rgba {
	return rgba{fromLinear(1 - (1-toLinear(c.R))*v), fromLinear(1 - (1-toLinear(c.G))*v), fromLinear(1 - (1-toLinear(c.B))*v), c.A}
}

// average returns a representative solid color (for text drawn with a
// gradient fill and similar approximations).
func (f fill) average() (rgba, bool) {
	switch f.kind {
	case fillSolid:
		return f.color, true
	case fillGrad:
		var r rgba
		for _, s := range f.stops {
			r.R += s.c.R
			r.G += s.c.G
			r.B += s.c.B
			r.A += s.c.A
		}
		n := float64(len(f.stops))
		return rgba{r.R / n, r.G / n, r.B / n, r.A / n}, true
	case fillPatt:
		return f.fg, true
	}
	return rgba{}, false
}

// setFill emits the fill style for a box of w×h (local coordinates); it
// returns false when the fill needs clipping and drawing instead (images).
func (cv *canvas) setFill(f fill, w, h float64) bool {
	switch f.kind {
	case fillSolid:
		cv.obj.FillColor(f.color.bdf())
		return true
	case fillGrad:
		cv.obj.FillPaint(cv.obj.AddPaint(gradientPaint(f, w, h)))
		return true
	case fillPatt:
		img := cv.c.patternImage(f)
		ref := cv.image(img)
		p := bdf.Pattern(ref, bdf.RepeatBoth)
		p.Matrix = [6]float32{0.75, 0, 0, 0.75, 0, 0} // 96 dpi pixels
		cv.obj.FillPaint(cv.obj.AddPaint(p))
		return true
	}
	return false
}

// gradientPaint builds the Canvas gradient for a DrawingML gradient over a
// w×h box: linear gradients run through the center at the given angle
// (scaled with the box when "scaled" is set) and cover the corners; path
// gradients become radial ones from the fillToRect center.
func gradientPaint(f fill, w, h float64) bdf.Paint {
	stops := make([]bdf.Stop, len(f.stops))
	for i, s := range f.stops {
		stops[i] = bdf.Stop{Offset: f32(s.pos), Color: s.c.bdf()}
	}
	if path := f.grad.child("path"); path != nil {
		r := path.child("fillToRect")
		l, t, rr, b := pct(r, "l", 0), pct(r, "t", 0), pct(r, "r", 0), pct(r, "b", 0)
		cx, cy := w*(l+(1-rr))/2, h*(t+(1-b))/2
		// distance to the farthest corner
		rad := 0.0
		for _, p := range [][2]float64{{0, 0}, {w, 0}, {0, h}, {w, h}} {
			rad = math.Max(rad, math.Hypot(p[0]-cx, p[1]-cy))
		}
		if path.attrStr("path", "circle") != "circle" {
			// rect/shape: approximate the rectangular spread with the half diagonal
			rad = math.Hypot(math.Max(cx, w-cx), math.Max(cy, h-cy))
		}
		return bdf.RadialGradient(f32(cx), f32(cy), 0, f32(cx), f32(cy), f32(rad), stops...)
	}
	lin := f.grad.child("lin")
	ang := float64(lin.attrInt("ang", 0)) * angUnit
	gx, gy := math.Cos(ang), math.Sin(ang)
	if lin.attrBool("scaled", false) && w > 0 && h > 0 {
		gx, gy = gx/w, gy/h
	}
	// Normalize so that the gradient parameter spans the box corners.
	t := (math.Abs(gx)*w + math.Abs(gy)*h) / 2
	g2 := gx*gx + gy*gy
	if g2 == 0 || t == 0 {
		return bdf.LinearGradient(0, 0, f32(w), 0, stops...)
	}
	cx, cy := w/2, h/2
	dx, dy := gx*t/g2, gy*t/g2
	return bdf.LinearGradient(f32(cx-dx), f32(cy-dy), f32(cx+dx), f32(cy+dy), stops...)
}

// patternImage renders a preset pattern as an 8×8 tile.
func (c *converter) patternImage(f fill) bdf.Hash {
	key := f.patt + f.fg.bdf().CSS() + f.bg.bdf().CSS()
	if h, ok := c.patterns[key]; ok {
		return h
	}
	bits := patternBits(f.patt)
	img := image.NewNRGBA(image.Rect(0, 0, 8, 8))
	fg := color.NRGBA{uint8(f.fg.R * 255), uint8(f.fg.G * 255), uint8(f.fg.B * 255), uint8(f.fg.A * 255)}
	bg := color.NRGBA{uint8(f.bg.R * 255), uint8(f.bg.G * 255), uint8(f.bg.B * 255), uint8(f.bg.A * 255)}
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if bits[y]&(0x80>>x) != 0 {
				img.SetNRGBA(x, y, fg)
			} else {
				img.SetNRGBA(x, y, bg)
			}
		}
	}
	var b bytes.Buffer
	png.Encode(&b, img)
	h := c.doc.AddImage(b.Bytes())
	c.patterns[key] = h
	return h
}

// patternBits are 8×8 bitmaps of the preset patterns (set bit = foreground).
func patternBits(name string) [8]byte {
	switch name {
	case "pct5":
		return [8]byte{0x80, 0, 0, 0, 0x08, 0, 0, 0}
	case "pct10":
		return [8]byte{0x80, 0, 0x08, 0, 0x80, 0, 0x08, 0}
	case "pct20":
		return [8]byte{0x88, 0, 0x22, 0, 0x88, 0, 0x22, 0}
	case "pct25":
		return [8]byte{0x88, 0x22, 0x88, 0x22, 0x88, 0x22, 0x88, 0x22}
	case "pct30":
		return [8]byte{0xaa, 0x44, 0xaa, 0x11, 0xaa, 0x44, 0xaa, 0x11}
	case "pct40":
		return [8]byte{0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x11}
	case "pct50":
		return [8]byte{0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55}
	case "pct60":
		return [8]byte{0xee, 0x55, 0xbb, 0x55, 0xee, 0x55, 0xbb, 0x55}
	case "pct70":
		return [8]byte{0xee, 0x77, 0xbb, 0xdd, 0xee, 0x77, 0xbb, 0xdd}
	case "pct75":
		return [8]byte{0xee, 0xff, 0xbb, 0xff, 0xee, 0xff, 0xbb, 0xff}
	case "pct80":
		return [8]byte{0xf7, 0xff, 0x7f, 0xff, 0xf7, 0xff, 0x7f, 0xff}
	case "pct90":
		return [8]byte{0xff, 0xff, 0xf7, 0xff, 0xff, 0xff, 0x7f, 0xff}
	case "horz", "ltHorz":
		return [8]byte{0xff, 0, 0, 0, 0, 0, 0, 0}
	case "dkHorz":
		return [8]byte{0xff, 0xff, 0, 0, 0xff, 0xff, 0, 0}
	case "narHorz":
		return [8]byte{0xff, 0, 0xff, 0, 0xff, 0, 0xff, 0}
	case "wdHorz":
		return [8]byte{0xff, 0xff, 0xff, 0, 0, 0, 0, 0}
	case "vert", "ltVert":
		return [8]byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}
	case "dkVert":
		return [8]byte{0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc, 0xcc}
	case "narVert":
		return [8]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa}
	case "wdVert":
		return [8]byte{0xe0, 0xe0, 0xe0, 0xe0, 0xe0, 0xe0, 0xe0, 0xe0}
	case "dnDiag", "ltDnDiag":
		return [8]byte{0x80, 0x40, 0x20, 0x10, 0x08, 0x04, 0x02, 0x01}
	case "upDiag", "ltUpDiag":
		return [8]byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}
	case "dkDnDiag", "wdDnDiag":
		return [8]byte{0xc1, 0xe0, 0x70, 0x38, 0x1c, 0x0e, 0x07, 0x83}
	case "dkUpDiag", "wdUpDiag":
		return [8]byte{0x83, 0x07, 0x0e, 0x1c, 0x38, 0x70, 0xe0, 0xc1}
	case "dashDnDiag":
		return [8]byte{0x80, 0x40, 0x20, 0, 0x08, 0x04, 0x02, 0}
	case "dashUpDiag":
		return [8]byte{0x01, 0x02, 0x04, 0, 0x10, 0x20, 0x40, 0}
	case "dashHorz":
		return [8]byte{0xf0, 0, 0, 0, 0x0f, 0, 0, 0}
	case "dashVert":
		return [8]byte{0x80, 0x80, 0x80, 0x80, 0x08, 0x08, 0x08, 0x08}
	case "cross", "lgGrid":
		return [8]byte{0xff, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}
	case "smGrid":
		return [8]byte{0xff, 0x88, 0x88, 0x88, 0xff, 0x88, 0x88, 0x88}
	case "diagCross", "openDmnd":
		return [8]byte{0x81, 0x42, 0x24, 0x18, 0x18, 0x24, 0x42, 0x81}
	case "smCheck":
		return [8]byte{0x99, 0x66, 0x66, 0x99, 0x99, 0x66, 0x66, 0x99}
	case "lgCheck":
		return [8]byte{0xf0, 0xf0, 0xf0, 0xf0, 0x0f, 0x0f, 0x0f, 0x0f}
	case "solidDmnd":
		return [8]byte{0x10, 0x38, 0x7c, 0xfe, 0x7c, 0x38, 0x10, 0}
	case "dotDmnd":
		return [8]byte{0x80, 0, 0x22, 0, 0x08, 0, 0x22, 0}
	case "smConfetti", "lgConfetti":
		return [8]byte{0x80, 0x04, 0x40, 0x02, 0x20, 0x01, 0x10, 0x08}
	case "horzBrick":
		return [8]byte{0xff, 0x80, 0x80, 0x80, 0xff, 0x08, 0x08, 0x08}
	case "diagBrick":
		return [8]byte{0x80, 0x40, 0x20, 0x10, 0x18, 0x24, 0x42, 0x81}
	case "zigZag", "wave":
		return [8]byte{0x81, 0x42, 0x24, 0x18, 0, 0, 0, 0}
	case "trellis":
		return [8]byte{0xff, 0x66, 0xff, 0x99, 0xff, 0x66, 0xff, 0x99}
	case "sphere", "weave", "plaid", "divot", "shingle", "dotGrid":
		return [8]byte{0x88, 0x55, 0x22, 0x55, 0x88, 0x55, 0x22, 0x55}
	}
	return [8]byte{0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55, 0xaa, 0x55}
}

// line is a resolved outline.
type line struct {
	fill  fill
	width float64
	cap   byte
	join  byte
	miter float64
	dash  []float64 // in multiples of the width
	head  arrow
	tail  arrow
}

type arrow struct {
	kind string // triangle, stealth, diamond, oval, arrow
	w, l float64
}

// resolveLine merges the a:ln elements (most specific first) and the
// theme line style of the p:style lnRef.
func (s *slideCtx) resolveLine(lns []*node, parts []string, ref *node, cc *colorCtx) *line {
	themeIdx := -1
	if ref != nil {
		if idx := int(ref.attrInt("idx", 0)); idx > 0 {
			if e := styleEntry(s.th.lines, idx); e != nil {
				themeIdx = len(lns)
				lns = append(lns, e)
				parts = append(parts, "")
			}
		}
	}
	if len(lns) == 0 {
		return nil
	}
	refCC := cc
	if phc, ok := cc.color(ref); ok {
		refCC = cc.withPh(phc)
	}
	l := &line{width: -1, miter: 8}
	var haveFill, haveCap, haveJoin, haveDash, haveHead, haveTail bool
	for i, ln := range lns {
		if ln == nil {
			continue
		}
		c := cc
		if i == themeIdx {
			c = refCC
		}
		if l.width < 0 {
			if _, ok := ln.attr("w"); ok {
				l.width = ln.emuAttr("w", 0.75)
			}
		}
		if !haveFill {
			if e := fillElem(ln); e != nil {
				l.fill = s.resolveFill(e, parts[i], c, nil)
				haveFill = true
			}
		}
		if !haveCap {
			if v, ok := ln.attr("cap"); ok {
				haveCap = true
				switch v {
				case "rnd":
					l.cap = 1
				case "sq":
					l.cap = 2
				}
			}
		}
		if !haveJoin {
			switch {
			case ln.child("round") != nil:
				l.join, haveJoin = 1, true
			case ln.child("bevel") != nil:
				l.join, haveJoin = 2, true
			case ln.child("miter") != nil:
				l.join, haveJoin = 0, true
				if lim := pct(ln.child("miter"), "lim", 0); lim > 0 {
					l.miter = lim
				}
			}
		}
		if !haveDash {
			if d := ln.child("prstDash"); d != nil {
				l.dash, haveDash = presetDash(d.attrStr("val", "solid")), true
			} else if d := ln.child("custDash"); d != nil {
				haveDash = true
				for _, ds := range d.children("ds") {
					l.dash = append(l.dash, pct(ds, "d", 1), pct(ds, "sp", 1))
				}
			}
		}
		if !haveHead {
			if e := ln.child("headEnd"); e != nil {
				l.head, haveHead = parseArrow(e), true
			}
		}
		if !haveTail {
			if e := ln.child("tailEnd"); e != nil {
				l.tail, haveTail = parseArrow(e), true
			}
		}
	}
	if l.fill.kind == fillNone || l.fill.kind == fillBlip {
		return nil
	}
	if l.width < 0 {
		l.width = 0.75
	}
	if l.width == 0 {
		l.width = 0.5 // hairline
	}
	return l
}

func presetDash(v string) []float64 {
	switch v {
	case "dot":
		return []float64{1, 3}
	case "dash":
		return []float64{4, 3}
	case "lgDash":
		return []float64{8, 3}
	case "dashDot":
		return []float64{4, 3, 1, 3}
	case "lgDashDot":
		return []float64{8, 3, 1, 3}
	case "lgDashDotDot":
		return []float64{8, 3, 1, 3, 1, 3}
	case "sysDash":
		return []float64{3, 1}
	case "sysDot":
		return []float64{1, 1}
	case "sysDashDot":
		return []float64{3, 1, 1, 1}
	case "sysDashDotDot":
		return []float64{3, 1, 1, 1, 1, 1}
	}
	return nil
}

func parseArrow(e *node) arrow {
	a := arrow{kind: e.attrStr("type", "none"), w: 3, l: 3}
	size := func(v string) float64 {
		switch v {
		case "sm":
			return 2
		case "lg":
			return 5
		}
		return 3
	}
	a.w, a.l = size(e.attrStr("w", "med")), size(e.attrStr("len", "med"))
	if a.kind == "none" {
		a.kind = ""
	}
	return a
}

// setLine emits the stroke style.
func (cv *canvas) setLine(l *line, w, h float64) {
	cv.obj.Line(f32(l.width), l.cap, l.join, f32(l.miter))
	if len(l.dash) > 0 {
		seg := make([]float32, len(l.dash))
		for i, d := range l.dash {
			seg[i] = f32(d * l.width)
		}
		cv.obj.Dash(seg, 0)
	}
	switch l.fill.kind {
	case fillSolid, fillPatt:
		c, _ := l.fill.average()
		cv.obj.StrokeColor(c.bdf())
	case fillGrad:
		cv.obj.StrokePaint(cv.obj.AddPaint(gradientPaint(l.fill, w, h)))
	}
}

// arrowHead draws a line end at (x, y) pointing along (dx, dy).
func (cv *canvas) arrowHead(l *line, a arrow, x, y, dx, dy float64) {
	if a.kind == "" {
		return
	}
	n := math.Hypot(dx, dy)
	if n == 0 {
		return
	}
	ux, uy := dx/n, dy/n // direction of travel at the tip
	px, py := -uy, ux    // perpendicular
	lw := math.Max(l.width, 1)
	hw, hl := a.w*lw/2, a.l*lw
	c, _ := l.fill.average()
	p := &bdf.Path{}
	pt := func(along, across float64) (float32, float32) {
		return f32(x + ux*along + px*across), f32(y + uy*along + py*across)
	}
	switch a.kind {
	case "oval":
		p.Ellipse(f32(x), f32(y), f32(hl/2), f32(hw), f32(math.Atan2(uy, ux)), 0, 2*math.Pi, false)
	case "diamond":
		p.MoveTo(pt(hl/2, 0))
		p.LineTo(pt(0, hw))
		p.LineTo(pt(-hl/2, 0))
		p.LineTo(pt(0, -hw))
		p.Close()
	case "arrow":
		p.MoveTo(pt(-hl, hw))
		p.LineTo(pt(0, 0))
		p.LineTo(pt(-hl, -hw))
		cv.obj.Save()
		cv.obj.Line(f32(l.width), 0, 0, 10)
		cv.obj.StrokeColor(c.bdf())
		cv.obj.StrokePath(cv.obj.AddPath(p))
		cv.obj.Restore()
		return
	case "stealth":
		p.MoveTo(pt(0, 0))
		p.LineTo(pt(-hl, hw))
		p.LineTo(pt(-hl*0.6, 0))
		p.LineTo(pt(-hl, -hw))
		p.Close()
	default: // triangle
		p.MoveTo(pt(0, 0))
		p.LineTo(pt(-hl, hw))
		p.LineTo(pt(-hl, -hw))
		p.Close()
	}
	cv.obj.FillColor(c.bdf())
	cv.obj.FillPath(cv.obj.AddPath(p), 0)
}

// shadow is an outer shadow effect.
type shadow struct {
	c            rgba
	blur, dx, dy float64
}

// resolveShadow finds an outer shadow among effect lists (most specific
// first) and the theme effect style of an effectRef.
func (s *slideCtx) resolveShadow(effects []*node, ref *node, cc *colorCtx) *shadow {
	lst := (*node)(nil)
	for _, e := range effects {
		if e != nil {
			lst = e
			break
		}
	}
	if lst == nil && ref != nil {
		if e := styleEntry(s.th.effects, int(ref.attrInt("idx", 0))); e != nil {
			lst = e.child("effectLst")
			if phc, ok := cc.color(ref); ok {
				cc = cc.withPh(phc)
			}
		}
	}
	sh := lst.child("outerShdw")
	if sh == nil {
		return nil
	}
	c, ok := cc.color(sh)
	if !ok {
		c = rgba{0, 0, 0, 0.5}
	}
	dist := sh.emuAttr("dist", 0)
	dir := float64(sh.attrInt("dir", 0)) * angUnit
	return &shadow{c: c, blur: sh.emuAttr("blurRad", 0), dx: dist * math.Cos(dir), dy: dist * math.Sin(dir)}
}
