package xlsx

import (
	"math"

	"github.com/shibukawa/bdf"
)

// The icons of icon sets, drawn as paths in a unit square. Each set lists
// its icons from the lowest values to the highest.

type iconKind byte

const (
	icArrowUp iconKind = iota
	icArrowUpRight
	icArrowSide
	icArrowDownRight
	icArrowDown
	icCircle
	icRimCircle
	icDiamond
	icTriangle
	icTriangleUp
	icTriangleDown
	icDash
	icFlag
	icCheck
	icExclaim
	icCross
	icCircleCheck
	icCircleExclaim
	icCircleCross
	icStar
	icRating
	icQuarter
	icBoxes
)

type icon struct {
	kind  iconKind
	color uint32
	level float64 // stars, ratings, quarters, boxes: how full (0..1)
}

const (
	icRed    = 0xE0301E
	icYellow = 0xFFB628
	icGreen  = 0x1D9E48
	icGray   = 0x808080
	icBlack  = 0x404040
	icPink   = 0xF4A6A6
	icBlue   = 0x4472C4
	icGold   = 0xF6B31B
)

var iconSets = map[string][]icon{
	"3Arrows":         {{icArrowDown, icRed, 0}, {icArrowSide, icYellow, 0}, {icArrowUp, icGreen, 0}},
	"3ArrowsGray":     {{icArrowDown, icGray, 0}, {icArrowSide, icGray, 0}, {icArrowUp, icGray, 0}},
	"3Flags":          {{icFlag, icRed, 0}, {icFlag, icYellow, 0}, {icFlag, icGreen, 0}},
	"3TrafficLights1": {{icCircle, icRed, 0}, {icCircle, icYellow, 0}, {icCircle, icGreen, 0}},
	"3TrafficLights2": {{icRimCircle, icRed, 0}, {icRimCircle, icYellow, 0}, {icRimCircle, icGreen, 0}},
	"3Signs":          {{icDiamond, icRed, 0}, {icTriangle, icYellow, 0}, {icCircle, icGreen, 0}},
	"3Symbols":        {{icCircleCross, icRed, 0}, {icCircleExclaim, icYellow, 0}, {icCircleCheck, icGreen, 0}},
	"3Symbols2":       {{icCross, icRed, 0}, {icExclaim, icYellow, 0}, {icCheck, icGreen, 0}},
	"3Stars":          {{icStar, icGold, 0}, {icStar, icGold, 0.5}, {icStar, icGold, 1}},
	"3Triangles":      {{icTriangleDown, icRed, 0}, {icDash, icYellow, 0}, {icTriangleUp, icGreen, 0}},
	"4Arrows":         {{icArrowDown, icRed, 0}, {icArrowDownRight, icYellow, 0}, {icArrowUpRight, icYellow, 0}, {icArrowUp, icGreen, 0}},
	"4ArrowsGray":     {{icArrowDown, icGray, 0}, {icArrowDownRight, icGray, 0}, {icArrowUpRight, icGray, 0}, {icArrowUp, icGray, 0}},
	"4RedToBlack":     {{icCircle, icBlack, 0}, {icCircle, icGray, 0}, {icCircle, icPink, 0}, {icCircle, icRed, 0}},
	"4Rating":         {{icRating, icBlue, 0.25}, {icRating, icBlue, 0.5}, {icRating, icBlue, 0.75}, {icRating, icBlue, 1}},
	"4TrafficLights":  {{icCircle, icBlack, 0}, {icCircle, icRed, 0}, {icCircle, icYellow, 0}, {icCircle, icGreen, 0}},
	"5Arrows":         {{icArrowDown, icRed, 0}, {icArrowDownRight, icYellow, 0}, {icArrowSide, icYellow, 0}, {icArrowUpRight, icYellow, 0}, {icArrowUp, icGreen, 0}},
	"5ArrowsGray":     {{icArrowDown, icGray, 0}, {icArrowDownRight, icGray, 0}, {icArrowSide, icGray, 0}, {icArrowUpRight, icGray, 0}, {icArrowUp, icGray, 0}},
	"5Rating":         {{icRating, icBlue, 0}, {icRating, icBlue, 0.25}, {icRating, icBlue, 0.5}, {icRating, icBlue, 0.75}, {icRating, icBlue, 1}},
	"5Quarters":       {{icQuarter, icBlack, 0}, {icQuarter, icBlack, 0.25}, {icQuarter, icBlack, 0.5}, {icQuarter, icBlack, 0.75}, {icQuarter, icBlack, 1}},
	"5Boxes":          {{icBoxes, icBlue, 0}, {icBoxes, icBlue, 0.25}, {icBoxes, icBlue, 0.5}, {icBoxes, icBlue, 0.75}, {icBoxes, icBlue, 1}},
}

// drawIcon draws the icon of a cell at its left, 12 points high at most.
func (s *sheetCtx) drawIcon(b box, ic *iconRef) {
	set, ok := iconSets[ic.set]
	if !ok {
		s.c.warnOnce("iconset:"+ic.set, "icon set %s is drawn as traffic lights", ic.set)
		set = iconSets["3TrafficLights1"]
	}
	if ic.idx < 0 || ic.idx >= len(set) {
		return
	}
	size := math.Min(b.h-2*pxPt, 16*pxPt)
	if size <= 0 {
		return
	}
	x, y := b.x+padL, b.y+b.h-padV-size
	for _, t := range s.tilesIn(x, y, x+size, y+size) {
		t.save()
		t.cv.Obj.Transform(f32(size), 0, 0, f32(size), f32(x-t.ox), f32(y-t.oy))
		drawIconShape(t, set[ic.idx])
		t.restore()
		t.cv.Drawn = true
	}
}

// path builder in the unit square
func poly(pts ...float64) *bdf.Path {
	p := &bdf.Path{}
	p.MoveTo(float32(pts[0]), float32(pts[1]))
	for i := 2; i+1 < len(pts); i += 2 {
		p.LineTo(float32(pts[i]), float32(pts[i+1]))
	}
	p.Close()
	return p
}

func drawIconShape(t *tileCv, ic icon) {
	o := t.cv.Obj
	col := hexRGB(ic.color).bdf()
	fill := func(p *bdf.Path, c bdf.Color) {
		t.fillColor(c)
		o.FillPath(o.AddPath(p), 0)
	}
	circle := func(cx, cy, r float64) *bdf.Path {
		p := &bdf.Path{}
		p.Circle(float32(cx), float32(cy), float32(r))
		return p
	}
	white := bdf.RGB(255, 255, 255)
	switch ic.kind {
	case icArrowUp:
		fill(poly(0.5, 0.05, 0.92, 0.5, 0.65, 0.5, 0.65, 0.95, 0.35, 0.95, 0.35, 0.5, 0.08, 0.5), col)
	case icArrowDown:
		fill(poly(0.5, 0.95, 0.92, 0.5, 0.65, 0.5, 0.65, 0.05, 0.35, 0.05, 0.35, 0.5, 0.08, 0.5), col)
	case icArrowSide:
		fill(poly(0.95, 0.5, 0.5, 0.92, 0.5, 0.65, 0.05, 0.65, 0.05, 0.35, 0.5, 0.35, 0.5, 0.08), col)
	case icArrowUpRight, icArrowDownRight:
		o.Save()
		a := -45.0
		if ic.kind == icArrowDownRight {
			a = 45
		}
		s, c := math.Sincos(a * math.Pi / 180)
		o.Transform(float32(c), float32(s), float32(-s), float32(c), float32(0.5-0.5*c+0.5*s), float32(0.5-0.5*s-0.5*c))
		fill(poly(0.95, 0.5, 0.5, 0.92, 0.5, 0.65, 0.05, 0.65, 0.05, 0.35, 0.5, 0.35, 0.5, 0.08), col)
		o.Restore()
		t.st.hasFill = false
	case icCircle:
		fill(circle(0.5, 0.5, 0.42), col)
	case icRimCircle:
		fill(circle(0.5, 0.5, 0.48), bdf.RGB(0x30, 0x30, 0x30))
		fill(circle(0.5, 0.5, 0.36), col)
	case icDiamond:
		fill(poly(0.5, 0.05, 0.95, 0.5, 0.5, 0.95, 0.05, 0.5), col)
	case icTriangle:
		fill(poly(0.5, 0.08, 0.95, 0.9, 0.05, 0.9), col)
	case icTriangleUp:
		fill(poly(0.5, 0.2, 0.9, 0.8, 0.1, 0.8), col)
	case icTriangleDown:
		fill(poly(0.1, 0.2, 0.9, 0.2, 0.5, 0.8), col)
	case icDash:
		fill(poly(0.15, 0.4, 0.85, 0.4, 0.85, 0.6, 0.15, 0.6), col)
	case icFlag:
		fill(poly(0.15, 0.05, 0.25, 0.05, 0.25, 0.95, 0.15, 0.95), bdf.RGB(0x40, 0x40, 0x40))
		fill(poly(0.25, 0.08, 0.9, 0.3, 0.25, 0.55), col)
	case icCheck:
		fill(poly(0.08, 0.55, 0.22, 0.4, 0.4, 0.6, 0.8, 0.15, 0.94, 0.3, 0.4, 0.88), col)
	case icExclaim:
		fill(poly(0.4, 0.08, 0.6, 0.08, 0.56, 0.65, 0.44, 0.65), col)
		fill(circle(0.5, 0.82, 0.1), col)
	case icCross:
		fill(crossPath(0.15, 0.85, 0.1), col)
	case icCircleCheck, icCircleExclaim, icCircleCross:
		fill(circle(0.5, 0.5, 0.46), col)
		switch ic.kind {
		case icCircleCheck:
			fill(poly(0.22, 0.52, 0.32, 0.42, 0.44, 0.56, 0.7, 0.28, 0.8, 0.38, 0.44, 0.74), white)
		case icCircleExclaim:
			fill(poly(0.44, 0.2, 0.56, 0.2, 0.54, 0.6, 0.46, 0.6), white)
			fill(circle(0.5, 0.74, 0.07), white)
		default:
			fill(crossPath(0.3, 0.7, 0.07), white)
		}
	case icStar:
		star := starPath()
		if ic.level < 1 {
			fill(star, bdf.RGB(0xd9, 0xd9, 0xd9))
		}
		if ic.level > 0 {
			o.Save()
			o.ClipRect(0, 0, float32(ic.level), 1)
			fill(star, col)
			o.Restore()
			t.st.hasFill = false
		}
	case icRating:
		for i := 0; i < 4; i++ {
			h := 0.25 + 0.2*float64(i)
			c := bdf.RGB(0xd9, 0xd9, 0xd9)
			if float64(i+1) <= ic.level*4+1e-9 {
				c = col
			}
			x := 0.1 + 0.21*float64(i)
			fill(poly(x, 0.9-h, x+0.16, 0.9-h, x+0.16, 0.9, x, 0.9), c)
		}
	case icQuarter:
		fill(circle(0.5, 0.5, 0.44), col)
		fill(circle(0.5, 0.5, 0.36), white)
		if ic.level > 0 {
			p := &bdf.Path{}
			p.MoveTo(0.5, 0.5)
			p.Ellipse(0.5, 0.5, 0.36, 0.36, 0, float32(-math.Pi/2), float32(-math.Pi/2+2*math.Pi*ic.level), false)
			p.Close()
			fill(p, col)
		}
	case icBoxes:
		for i := 0; i < 4; i++ {
			x, y := 0.1+0.42*float64(i%2), 0.1+0.42*float64(i/2)
			c := bdf.RGB(0xd9, 0xd9, 0xd9)
			if float64(i+1) <= ic.level*4+1e-9 {
				c = col
			}
			fill(poly(x, y, x+0.38, y, x+0.38, y+0.38, x, y+0.38), c)
		}
	}
}

func crossPath(a, b, w float64) *bdf.Path {
	p := &bdf.Path{}
	for _, seg := range [][4]float64{{a, a, b, b}, {a, b, b, a}} {
		dx, dy := seg[2]-seg[0], seg[3]-seg[1]
		l := math.Hypot(dx, dy)
		nx, ny := -dy/l*w, dx/l*w
		p.MoveTo(float32(seg[0]+nx), float32(seg[1]+ny))
		p.LineTo(float32(seg[2]+nx), float32(seg[3]+ny))
		p.LineTo(float32(seg[2]-nx), float32(seg[3]-ny))
		p.LineTo(float32(seg[0]-nx), float32(seg[1]-ny))
		p.Close()
	}
	return p
}

func starPath() *bdf.Path {
	var pts []float64
	for i := 0; i < 10; i++ {
		r := 0.47
		if i%2 == 1 {
			r = 0.2
		}
		a := -math.Pi/2 + float64(i)*math.Pi/5
		pts = append(pts, 0.5+r*math.Cos(a), 0.53+r*math.Sin(a))
	}
	return poly(pts...)
}
