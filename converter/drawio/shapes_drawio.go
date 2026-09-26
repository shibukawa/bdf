package drawio

import (
	"math"
	"strconv"
	"strings"
)

// draw.io's own shapes from Shapes.js (its mxCellRenderer.registerShape
// calls), ported class by class. The redrawPath of an mxActor subclass is
// the path function of actorShape, the redrawPath of an mxCylinder
// subclass the path function of cylinderShape; subclasses of
// mxRectangleShape, mxEllipse and mxRhombus reuse the painters of those
// core shapes. Sizes are read as the JavaScript reads them: relative sizes
// are fractions of the width or height, fixedSize=1 makes them absolute.

func init() {
	registerShape("process", processShape)
	registerShape("process2", processShape) // backwards compatibility
	registerShape("transparent", transparentShape)
	registerShape("plus", plusShape)
	registerShape("ext", extendedShape)
	registerShape("internalStorage", internalStorageShape)
	registerShape("associativeEntity", associativeEntityShape)

	registerShape("document", documentShape)
	registerShape("parallelogram", parallelogramShape)
	registerShape("trapezoid", trapezoidShape)
	registerShape("step", stepShape)
	registerShape("hexagon", hexagonShape)
	registerShape("card", cardShape)
	registerShape("tape", tapeShape)
	registerShape("curlyBracket", curlyBracketShape)
	registerShape("parallelMarker", parallelMarkerShape)
	registerShape("switch", switchShape)
	registerShape("isoRectangle", isoRectangleShape)
	registerShape("callout", calloutShape)
	registerShape("wedgeCallout", wedgeCalloutShape)
	registerShape("manualInput", manualInputShape)
	registerShape("corner", cornerShape)
	registerShape("crossbar", crossbarShape)
	registerShape("tee", teeShape)
	registerShape("singleArrow", singleArrowShape)
	registerShape("doubleArrow", doubleArrowShape)
	registerShape("dataStorage", dataStorageShape)
	registerShape("or", orShape)
	registerShape("xor", xorShape)
	registerShape("loopLimit", loopLimitShape)
	registerShape("offPageConnector", offPageConnectorShape)
	registerShape("delay", delayShape)
	registerShape("cross", crossShape)
	registerShape("display", displayShape)

	registerShape("cube", cubeShape)
	registerShape("isoCube", isoCubeShape)
	registerShape("isoCube2", isoCube2Shape)
	registerShape("datastore", dataStoreShape)
	registerShape("note", noteShape)
	registerShape("note2", note2Shape)
	registerShape("cylinder2", cylinder2Shape)
	registerShape("cylinder3", cylinder3Shape)
	registerShape("folder", folderShape)
	registerShape("message", messageShape)
	registerShape("waypoint", waypointShape)

	registerShape("tapeData", tapeDataShape)
	registerShape("orEllipse", orEllipseShape)
	registerShape("sumEllipse", sumEllipseShape)
	registerShape("sortShape", sortShape)
	registerShape("collate", collateShape)
	registerShape("dimension", dimensionShape)
	registerShape("partialRectangle", partialRectangleShape)
	registerShape("lineEllipse", lineEllipseShape)
	registerShape("endState", stateShape(true))
	registerShape("startState", stateShape(false))
	registerShape("smileyFace", smileyFaceShape)
	registerShape("zigzag", zigzagShape)
}

// styleBool is mxUtils.getValue(style, key, def) used as a condition, as
// JavaScript's truthiness sees the style value: numbers are true unless 0,
// other strings unless empty.
func styleBool(st style, key string, def bool) bool {
	v, ok := st[key]
	if !ok {
		return def
	}
	if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
		return f != 0
	}
	return v != ""
}



// lineArc is the corner radius of rounded polygons
// (mxUtils.getValue(style, 'arcSize', LINE_ARCSIZE) / 2).
func (s *shape) lineArc() float64 { return s.style.num("arcSize", lineArcSize) / 2 }

// rectArc is the rounding of rounded rectangles in the fraction form
// (arcSize / 100, RECTANGLE_ROUNDING_FACTOR by default) as inset shapes
// use it.
func (s *shape) rectArc(w, h float64) float64 {
	f := s.style.num("arcSize", rectangleRoundingFactor*100) / 100
	return math.Min(w*f, h*f)
}

// fixedOrRelative returns the size of shapes with a fixedSize style: the
// size style in units clamped to limit when fixedSize is set, else the
// size style as a fraction (clamped to maxRel) of ref.
func (s *shape) fixedOrRelative(def, fixedDef, limit, ref, maxRel float64) float64 {
	if styleBool(s.style, "fixedSize", false) {
		return clamp(s.style.num("size", fixedDef), 0, limit)
	}
	return ref * clamp(s.style.num("size", def), 0, maxRel)
}

func boundedLbl(s *shape) bool { return styleBool(s.style, "boundedLbl", false) }

// --- mxRectangleShape subclasses ---

// processShape ports ProcessShape: a rectangle with inner vertical lines.
var processShape = &shapeDef{
	paintBackground: paintRectangleBackground,
	paintForeground: func(s *shape, c *c2d, x, y, w, h float64) {
		inset := s.style.num("size", 0.1)
		if styleBool(s.style, "fixedSize", false) {
			inset = clamp(inset, 0, w)
		} else {
			inset = w * clamp(inset, 0, 1)
		}
		if s.isRounded {
			inset = math.Max(inset, s.rectArc(w, h))
		}
		// crisp rendering of inner lines
		inset = jsRound(inset)
		c.begin()
		c.moveTo(x+inset, y)
		c.lineTo(x+inset, y+h)
		c.moveTo(x+w-inset, y)
		c.lineTo(x+w-inset, y+h)
		c.stroke()
		paintRectangleForeground(s, c, x, y, w, h)
	},
	labelBounds: func(s *shape, r rect) rect {
		dirH := s.direction == "" || s.direction == "east" || s.direction == "west"
		// horizontal == (direction is east or west), with JavaScript's ==
		match := dirH
		if v, ok := s.style["horizontal"]; ok {
			f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			match = err == nil && (dirH && f == 1 || !dirH && f == 0)
		}
		if !match {
			return r
		}
		w, h := r.w, r.h
		inset := s.style.num("size", 0.1)
		if styleBool(s.style, "fixedSize", false) {
			inset = clamp(inset, 0, w)
		} else {
			inset = w * clamp(inset, 0, 1)
			if s.isRounded {
				inset = math.Max(inset, s.rectArc(w, h))
			}
		}
		r.x += jsRound(inset)
		r.w -= jsRound(2 * inset)
		return r
	},
	roundable: true,
}

// transparentShape ports TransparentShape: an invisible rectangle.
var transparentShape = &shapeDef{
	paintBackground: func(s *shape, c *c2d, x, y, w, h float64) {
		c.setFillColor("")
		c.rect(x, y, w, h)
		c.fill()
	},
	labelBounds: rectangleLabelBounds,
	roundable:   true,
}

// plusShape ports PlusShape: a rectangle with a cross inside.
var plusShape = &shapeDef{
	paintBackground: paintRectangleBackground,
	paintForeground: func(s *shape, c *c2d, x, y, w, h float64) {
		border := math.Min(w/5, h/5) + 1
		c.begin()
		c.moveTo(x+w/2, y+border)
		c.lineTo(x+w/2, y+h-border)
		c.moveTo(x+border, y+h/2)
		c.lineTo(x+w-border, y+h/2)
		c.stroke()
		paintRectangleForeground(s, c, x, y, w, h)
	},
	labelBounds: rectangleLabelBounds,
	roundable:   true,
}

// extMargin is the inset of the inner rectangle of double=1 (ExtendedShape).
func extMargin(s *shape) float64 {
	return math.Max(2, s.strokewidth+1) + s.style.num("margin", 0)
}

// extendedShape ports ExtendedShape ("ext"): a rectangle with an optional
// inner rectangle (double=1) and symbols (symbol0, symbol1, … naming
// registered shapes, placed by the symbolN* styles).
var extendedShape = &shapeDef{
	paintBackground: paintRectangleBackground,
	paintForeground: func(s *shape, c *c2d, x, y, w, h float64) {
		if s.style.num("double", 0) == 1 {
			m := extMargin(s)
			x, y, w, h = x+m, y+m, w-2*m, h-2*m
			if w > 0 && h > 0 {
				paintRectangleBackground(s, c, x, y, w, h)
			}
		}
		c.setDashed(false, false)
		for i := 0; ; i++ {
			key := "symbol" + itoa(i)
			d := shapeRegistry[s.style.get(key, "")]
			if d == nil {
				break
			}
			if s.style.get(key, "") == "ext" {
				continue // would recurse forever
			}
			width := s.style.num(key+"Width", 0)
			height := s.style.num(key+"Height", 0)
			spacing := s.style.num(key+"Spacing", 0)
			vspacing := s.style.num(key+"VSpacing", spacing)
			if s.style.has(key + "ArcSpacing") {
				arc := s.arcSize(w+s.strokewidth, h+s.strokewidth) * s.style.num(key+"ArcSpacing", 0)
				spacing += arc
				vspacing += arc
			}
			x2, y2 := x, y
			switch s.style.get(key+"Align", "") {
			case "center":
				x2 += (w - width) / 2
			case "right":
				x2 += w - width - spacing
			default:
				x2 += spacing
			}
			switch s.style.get(key+"VerticalAlign", "") {
			case "middle":
				y2 += (h - height) / 2
			case "bottom":
				y2 += h - height - vspacing
			default:
				y2 += vspacing
			}
			c.save()
			// a new shape of the symbol's class that only has the style
			// (the canvas keeps the colors of this shape)
			sym := &shape{conv: s.conv, st: s.st, style: s.style, def: d, name: s.style.get(key, ""), strokewidth: 1, opacity: 100, fillOpacity: 100, strokeOpacity: 100}
			sym.paintVertexShape(c, x2, y2, width, height)
			c.restore()
		}
		paintRectangleForeground(s, c, x, y, w, h)
	},
	labelBounds: func(s *shape, r rect) rect {
		if s.style.num("double", 0) == 1 {
			m := extMargin(s)
			return rect{r.x + m, r.y + m, r.w - 2*m, r.h - 2*m}
		}
		return r
	},
	roundable: true,
}

// internalStorageShape ports InternalStorageShape: a rectangle with a
// horizontal line at dy and a vertical line at dx.
var internalStorageShape = &shapeDef{
	paintBackground: paintRectangleBackground,
	paintForeground: func(s *shape, c *c2d, x, y, w, h float64) {
		paintRectangleForeground(s, c, x, y, w, h)
		inset := 0.0
		if s.isRounded {
			inset = math.Max(inset, s.rectArc(w, h))
		}
		dx := math.Max(inset, math.Min(w, s.style.num("dx", 20)))
		dy := math.Max(inset, math.Min(h, s.style.num("dy", 20)))
		c.begin()
		c.moveTo(x, y+dy)
		c.lineTo(x+w, y+dy)
		c.stroke()
		c.begin()
		c.moveTo(x+dx, y)
		c.lineTo(x+dx, y+h)
		c.stroke()
	},
	labelBounds: rectangleLabelBounds,
	roundable:   true,
}

// associativeEntityShape ports AssociativeEntity: a rectangle with a
// rhombus inside.
var associativeEntityShape = &shapeDef{
	paintBackground: paintRectangleBackground,
	paintForeground: func(s *shape, c *c2d, x, y, w, h float64) {
		hw, hh := w/2, h/2
		c.begin()
		s.addPoints(c, []point{{x + hw, y}, {x + w, y + hh}, {x + hw, y + h}, {x, y + hh}}, s.isRounded, s.lineArc(), true, nil, true)
		c.stroke()
		paintRectangleForeground(s, c, x, y, w, h)
	},
	labelBounds: rectangleLabelBounds,
	roundable:   true,
}

// --- mxActor subclasses ---

// roundableActor is an actorShape that supports rounded.
func roundableActor(path func(s *shape, c *c2d, x, y, w, h float64)) *shapeDef {
	d := actorShape(path)
	d.roundable = true
	return d
}

// documentShape ports DocumentShape: a rectangle with a wavy bottom.
var documentShape = func() *shapeDef {
	d := actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
		dy := h * clamp(s.style.num("size", 0.3), 0, 1)
		const fy = 1.4
		c.moveTo(0, 0)
		c.lineTo(w, 0)
		c.lineTo(w, h-dy/2)
		c.quadTo(w*3/4, h-dy*fy, w/2, h-dy/2)
		c.quadTo(w/4, h-dy*(1-fy), 0, h-dy/2)
		c.lineTo(0, dy/2)
		c.close()
	})
	d.labelMargins = func(s *shape, r rect) *rect {
		if boundedLbl(s) {
			return &rect{0, 0, 0, s.style.num("size", 0.3) * r.h}
		}
		return nil
	}
	return d
}()

// parallelogramShape ports ParallelogramShape.
var parallelogramShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	dx := s.fixedOrRelative(0.2, 20, w, w, 1)
	s.addPoints(c, []point{{0, h}, {dx, 0}, {w, 0}, {w - dx, h}}, s.isRounded, s.lineArc(), true, nil, true)
})

// trapezoidShape ports TrapezoidShape.
var trapezoidShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	dx := s.fixedOrRelative(0.2, 20, w*0.5, w, 0.5)
	s.addPoints(c, []point{{0, h}, {dx, 0}, {w - dx, 0}, {w, h}}, s.isRounded, s.lineArc(), true, nil, true)
})

// stepShape ports StepShape: a chevron.
var stepShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	sz := s.fixedOrRelative(0.2, 20, w, w, 1)
	s.addPoints(c, []point{{0, 0}, {w - sz, 0}, {w, h / 2}, {w - sz, h}, {0, h}, {sz, h / 2}}, s.isRounded, s.lineArc(), true, nil, true)
})

// hexagonShape ports HexagonShape, which replaces mxHexagon.
var hexagonShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	sz := s.fixedOrRelative(0.25, 20, w*0.5, w, 1)
	s.addPoints(c, []point{{sz, 0}, {w - sz, 0}, {w, 0.5 * h}, {w - sz, h}, {sz, h}, {0, 0.5 * h}}, s.isRounded, s.lineArc(), true, nil, true)
})

// cardShape ports CardShape: a rectangle with a cut top left corner.
var cardShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	sz := clamp(s.style.num("size", 30), 0, math.Min(w, h))
	s.addPoints(c, []point{{sz, 0}, {w, 0}, {w, h}, {0, h}, {0, sz}}, s.isRounded, s.lineArc(), true, nil, true)
})

// tapeShape ports TapeShape: a rectangle with wavy top and bottom.
var tapeShape = func() *shapeDef {
	d := actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
		dy := h * clamp(s.style.num("size", 0.4), 0, 1)
		const fy = 1.4
		c.moveTo(0, dy/2)
		c.quadTo(w/4, dy*fy, w/2, dy/2)
		c.quadTo(w*3/4, dy*(1-fy), w, dy/2)
		c.lineTo(w, h-dy/2)
		c.quadTo(w*3/4, h-dy*fy, w/2, h-dy/2)
		c.quadTo(w/4, h-dy*(1-fy), 0, h-dy/2)
		c.lineTo(0, dy/2)
		c.close()
	})
	d.labelBounds = func(s *shape, r rect) rect {
		if !boundedLbl(s) {
			return r
		}
		size := s.style.num("size", 0.4)
		if s.direction == "" || s.direction == "east" || s.direction == "west" {
			dy := r.h * size
			return rect{r.x, r.y + dy, r.w, r.h - 2*dy}
		}
		dx := r.w * size
		return rect{r.x + dx, r.y, r.w - 2*dx, r.h}
	}
	return d
}()

// curlyBracketShape ports CurlyBracketShape: an unfilled bracket.
var curlyBracketShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	c.setFillColor("")
	sz := w * clamp(s.style.num("size", 0.5), 0, 1)
	s.addPoints(c, []point{{w, 0}, {sz, 0}, {sz, h / 2}, {0, h / 2}, {sz, h / 2}, {sz, h}, {w, h}}, s.isRounded, s.lineArc(), false, nil, true)
})

// parallelMarkerShape ports ParallelMarkerShape: three bars.
var parallelMarkerShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	c.setStrokeWidth(1)
	c.setFillColor(s.stroke)
	w2 := w / 5
	c.rect(0, 0, w2, h)
	c.fillAndStroke()
	c.rect(2*w2, 0, w2, h)
	c.fillAndStroke()
	c.rect(4*w2, 0, w2, h)
	c.fillAndStroke()
})

// switchShape ports SwitchShape: a rectangle with concave sides.
var switchShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	const curve = 0.5
	c.moveTo(0, 0)
	c.quadTo(w/2, h*curve, w, 0)
	c.quadTo(w*(1-curve), h/2, w, h)
	c.quadTo(w/2, h*(1-curve), 0, h)
	c.quadTo(w*curve, h/2, 0, 0)
})

var (
	tan30   = math.Tan(toRadians(30))
	tan30Dx = (0.5 - tan30) / 2
)

// isoRectangleShape ports IsoRectangleShape: an isometric rhombus.
var isoRectangleShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	m := math.Min(w, h/tan30)
	c.translate((w-m)/2, (h-m)/2+m/4)
	c.moveTo(0, 0.25*m)
	c.lineTo(0.5*m, m*tan30Dx)
	c.lineTo(m, 0.25*m)
	c.lineTo(0.5*m, (0.5-tan30Dx)*m)
	c.lineTo(0, 0.25*m)
	c.close()
})

// calloutShape ports CalloutShape: a speech bubble with its tail at the
// bottom (size is the tail height, position and position2 place the
// tail's base and tip, base is the base width).
var calloutShape = func() *shapeDef {
	d := roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
		sz := clamp(s.style.num("size", 30), 0, h)
		dx := w * clamp(s.style.num("position", 0.5), 0, 1)
		dx2 := w * clamp(s.style.num("position2", 0.5), 0, 1)
		base := clamp(s.style.num("base", 20), 0, w)
		s.addPoints(c, []point{{0, 0}, {w, 0}, {w, h - sz}, {math.Min(w, dx+base), h - sz}, {dx2, h},
			{math.Max(0, dx), h - sz}, {0, h - sz}}, s.isRounded, s.lineArc(), true, []int{4}, true)
	})
	d.labelMargins = func(s *shape, r rect) *rect {
		return &rect{0, 0, 0, s.style.num("size", 30)}
	}
	return d
}()

// wedgeTail returns the tail of a wedge callout (base, tip, base) and
// where it goes in the clockwise bubble path, or nil when the tip is
// inside the bubble (WedgeCalloutShape.getTailPoints).
func wedgeTail(s *shape, w, h float64) ([]point, int) {
	const maxTipOffset = 100
	tx := clamp(s.style.num("tipX", -0.25), -maxTipOffset, maxTipOffset)
	ty := clamp(s.style.num("tipY", 1), -maxTipOffset, maxTipOffset)
	dx, dy := tx*w, ty*h
	if math.Abs(dx) <= w/2 && math.Abs(dy) <= h/2 {
		return nil, 0
	}
	base := s.style.num("base", 20)
	if !finite(base) {
		base = 20
	}
	base = math.Max(0, base)
	inset := 0.0
	if s.isRounded {
		inset = s.lineArc()
	}
	tp := point{w/2 + dx, h/2 + dy}
	if math.Abs(dx)*h >= math.Abs(dy)*w && dx != 0 {
		// the tail leaves through the left or right side
		inset = math.Min(inset, h/2)
		hb := math.Max(0, math.Min(base, h-2*inset)) / 2
		ey := h/2 + dy*(w/2)/math.Abs(dx)
		ey = math.Max(inset+hb, math.Min(h-inset-hb, ey))
		if dx > 0 {
			return []point{{w, ey - hb}, tp, {w, ey + hb}}, 2
		}
		return []point{{0, ey + hb}, tp, {0, ey - hb}}, 4
	}
	// the tail leaves through the top or bottom side
	inset = math.Min(inset, w/2)
	hb := math.Max(0, math.Min(base, w-2*inset)) / 2
	ex := w/2 + dx*(h/2)/math.Abs(dy)
	ex = math.Max(inset+hb, math.Min(w-inset-hb, ex))
	if dy > 0 {
		return []point{{ex + hb, h}, tp, {ex - hb, h}}, 3
	}
	return []point{{ex - hb, 0}, tp, {ex + hb, 0}}, 1
}

// wedgeCalloutShape ports WedgeCalloutShape: a rectangle with a tail to
// a free tip (tipX, tipY relative to the size, from the center).
var wedgeCalloutShape = func() *shapeDef {
	d := roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
		pts := []point{{0, 0}, {w, 0}, {w, h}, {0, h}}
		var exclude []int
		if tail, i := wedgeTail(s, w, h); tail != nil {
			pts = append(pts[:i], append(tail, pts[i:]...)...)
			exclude = []int{i, i + 1, i + 2}
		}
		s.addPoints(c, pts, s.isRounded, s.lineArc(), true, exclude, true)
	})
	// the tip is outside of the bounds (WedgeCalloutShape.getShapeBoundingBox)
	d.augmentBounds = func(s *shape, r rect) rect {
		b := s.bounds
		if s.isPaintBoundsInverted() {
			b = b.rotate90()
		}
		tail, _ := wedgeTail(s, b.w, b.h)
		if tail == nil {
			return r
		}
		tp := tail[1]
		cx, cy := b.cx(), b.cy()
		fx, fy := 1.0, 1.0
		if s.flipH {
			fx = -1
		}
		if s.flipV {
			fy = -1
		}
		t := rect{cx + (tp.x-b.w/2)*fx, cy + (tp.y-b.h/2)*fy, 0, 0}.grow(s.strokewidth / 2)
		// the miter join at the acute tip
		v1x, v1y := tail[0].x-tp.x, tail[0].y-tp.y
		v2x, v2y := tail[2].x-tp.x, tail[2].y-tp.y
		n1, n2 := math.Hypot(v1x, v1y), math.Hypot(v2x, v2y)
		if n1 > 0 && n2 > 0 {
			cos := clamp((v1x*v2x+v1y*v2y)/(n1*n2), -1, 1)
			sinHalf := math.Sqrt((1 - cos) / 2)
			miter := 10.0
			if sinHalf > 0 {
				miter = math.Min(10, 1/sinHalf)
			}
			t = t.grow((miter - 1) * s.strokewidth / 2)
		}
		if rot := s.shapeRotation(); rot != 0 {
			rad := toRadians(rot)
			cos, sin := math.Cos(rad), math.Sin(rad)
			var bb rect
			for i, p := range []point{{t.x, t.y}, {t.x + t.w, t.y}, {t.x, t.y + t.h}, {t.x + t.w, t.y + t.h}} {
				q := rotatePoint(p, cos, sin, point{cx, cy})
				if i == 0 {
					bb = rect{q.x, q.y, 0, 0}
				} else {
					bb = bb.union(rect{q.x, q.y, 0, 0})
				}
			}
			t = bb
		}
		return r.union(t)
	}
	return d
}()

// manualInputShape ports ManualInputShape: a slanted top.
var manualInputShape = func() *shapeDef {
	d := roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
		sz := math.Min(h, s.style.num("size", 30))
		s.addPoints(c, []point{{0, h}, {0, sz}, {w, 0}, {w, h}}, s.isRounded, s.lineArc(), true, nil, true)
	})
	d.labelMargins = func(s *shape, r rect) *rect {
		if boundedLbl(s) {
			return &rect{0, s.style.num("size", 30), 0, 0}
		}
		return nil
	}
	return d
}()

// cornerShape ports CornerShape: an L.
var cornerShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	dx := clamp(s.style.num("dx", 20), 0, w)
	dy := clamp(s.style.num("dy", 20), 0, h)
	s.addPoints(c, []point{{0, 0}, {w, 0}, {w, dy}, {dx, dy}, {dx, h}, {0, h}}, s.isRounded, s.lineArc(), true, nil, true)
})

// crossbarShape ports CrossbarShape: an H on its side.
var crossbarShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	c.moveTo(0, 0)
	c.lineTo(0, h)
	c.moveTo(w, 0)
	c.lineTo(w, h)
	c.moveTo(0, h/2)
	c.lineTo(w, h/2)
})

// teeShape ports TeeShape: a T.
var teeShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	dx := clamp(s.style.num("dx", 20), 0, w)
	dy := clamp(s.style.num("dy", 20), 0, h)
	s.addPoints(c, []point{{0, 0}, {w, 0}, {w, dy}, {(w + dx) / 2, dy}, {(w + dx) / 2, h}, {(w - dx) / 2, h},
		{(w - dx) / 2, dy}, {0, dy}}, s.isRounded, s.lineArc(), true, nil, true)
})

// singleArrowShape ports SingleArrowShape (arrowWidth is the height of the
// shaft, arrowSize the length of the head, both relative).
var singleArrowShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	aw := h * clamp(s.style.num("arrowWidth", 0.3), 0, 1)
	as := w * clamp(s.style.num("arrowSize", 0.2), 0, 1)
	at := (h - aw) / 2
	ab := at + aw
	s.addPoints(c, []point{{0, at}, {w - as, at}, {w - as, 0}, {w, h / 2}, {w - as, h}, {w - as, ab}, {0, ab}},
		s.isRounded, s.lineArc(), true, nil, true)
})

// doubleArrowShape ports DoubleArrowShape.
var doubleArrowShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	aw := h * clamp(s.style.num("arrowWidth", 0.3), 0, 1)
	as := w * clamp(s.style.num("arrowSize", 0.2), 0, 1)
	at := (h - aw) / 2
	ab := at + aw
	s.addPoints(c, []point{{0, h / 2}, {as, 0}, {as, at}, {w - as, at}, {w - as, 0}, {w, h / 2}, {w - as, h},
		{w - as, ab}, {as, ab}, {as, h}}, s.isRounded, s.lineArc(), true, nil, true)
})

// dataStorageShape ports DataStorageShape: a rectangle with a concave
// right side and a convex left side.
var dataStorageShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	sz := s.fixedOrRelative(0.1, 20, w, w, 1)
	c.moveTo(sz, 0)
	c.lineTo(w, 0)
	c.quadTo(w-sz*2, h/2, w, h)
	c.lineTo(sz, h)
	c.quadTo(sz-sz*2, h/2, sz, 0)
	c.close()
})

// orShape ports OrShape: a D.
var orShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	c.moveTo(0, 0)
	c.quadTo(w, 0, w, h/2)
	c.quadTo(w, h, 0, h)
	c.close()
})

// xorShape ports XorShape: a D with a concave back.
var xorShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	c.moveTo(0, 0)
	c.quadTo(w, 0, w, h/2)
	c.quadTo(w, h, 0, h)
	c.quadTo(w/2, h/2, 0, 0)
	c.close()
})

// loopLimitShape ports LoopLimitShape: a rectangle with cut top corners.
var loopLimitShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	sz := math.Min(w/2, math.Min(h, s.style.num("size", 20)))
	s.addPoints(c, []point{{sz, 0}, {w - sz, 0}, {w, sz * 0.8}, {w, h}, {0, h}, {0, sz * 0.8}}, s.isRounded, s.lineArc(), true, nil, true)
})

// offPageConnectorShape ports OffPageConnectorShape: a pentagon pointing down.
var offPageConnectorShape = roundableActor(func(s *shape, c *c2d, x, y, w, h float64) {
	sz := h * clamp(s.style.num("size", 3.0/8), 0, 1)
	s.addPoints(c, []point{{0, 0}, {w, 0}, {w, h - sz}, {w / 2, h}, {0, h - sz}}, s.isRounded, s.lineArc(), true, nil, true)
})

// delayShape ports DelayShape: a rectangle with a round right side.
var delayShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	dx := math.Min(w, h/2)
	c.moveTo(0, 0)
	c.lineTo(w-dx, 0)
	c.quadTo(w, 0, w, h/2)
	c.quadTo(w, h, w-dx, h)
	c.lineTo(0, h)
	c.close()
})

// crossShape ports CrossShape: a plus sign (size is the relative bar width).
var crossShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	m := math.Min(h, w)
	sz := clamp(m*s.style.num("size", 0.2), 0, m)
	t := (h - sz) / 2
	b := t + sz
	l := (w - sz) / 2
	r := l + sz
	c.moveTo(0, t)
	c.lineTo(l, t)
	c.lineTo(l, 0)
	c.lineTo(r, 0)
	c.lineTo(r, t)
	c.lineTo(w, t)
	c.lineTo(w, b)
	c.lineTo(r, b)
	c.lineTo(r, h)
	c.lineTo(l, h)
	c.lineTo(l, b)
	c.lineTo(0, b)
	c.close()
})

// displayShape ports DisplayShape.
var displayShape = actorShape(func(s *shape, c *c2d, x, y, w, h float64) {
	dx := math.Min(w, h/2)
	sz := math.Min(w-dx, math.Max(0, s.style.num("size", 0.25))*w)
	c.moveTo(0, h/2)
	c.lineTo(sz, 0)
	c.lineTo(w-dx, 0)
	c.quadTo(w, 0, w, h/2)
	c.quadTo(w, h, w-dx, h)
	c.lineTo(sz, h)
	c.close()
})

// --- mxCylinder subclasses ---

// cylinderShape is an mxCylinder: redraw adds the background path
// (foreground false) or the lines painted over it (foreground true), in
// coordinates relative to the bounds.
func cylinderShape(redraw func(s *shape, c *c2d, w, h float64, foreground bool)) *shapeDef {
	return &shapeDef{
		paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
			c.translate(x, y)
			c.begin()
			redraw(s, c, w, h, false)
			c.fillAndStroke()
			c.begin()
			redraw(s, c, w, h, true)
			c.stroke()
		},
		labelMargins: cylinderLabelMargins,
	}
}

// cylinderLabelMargins is draw.io's mxCylinder.getLabelMargins, which the
// subclasses of mxCylinder inherit: boundedLbl puts the label below the
// top.
func cylinderLabelMargins(s *shape, r rect) *rect {
	if boundedLbl(s) {
		size := s.style.num("size", 0.15) * 2
		return &rect{0, math.Min(cylinderMaxHeight, r.h*size), 0, 0}
	}
	return nil
}

// cylinderMaxHeight is mxCylinder.maxHeight.
const cylinderMaxHeight = 40

// cubeShape ports CubeShape: a box in cabinet projection (darkOpacity and
// darkOpacity2 shade the top and the side).
var cubeShape = &shapeDef{
	paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		sz := clamp(s.style.num("size", 20), 0, math.Min(w, h))
		op := clamp(s.style.num("darkOpacity", 0), -1, 1)
		op2 := clamp(s.style.num("darkOpacity2", 0), -1, 1)
		c.translate(x, y)
		c.begin()
		c.moveTo(0, 0)
		c.lineTo(w-sz, 0)
		c.lineTo(w, sz)
		c.lineTo(w, h)
		c.lineTo(sz, h)
		c.lineTo(0, h-sz)
		c.lineTo(0, 0)
		c.close()
		c.fillAndStroke()
		if op != 0 {
			c.setFillAlpha(math.Abs(op))
			c.setFillColor(darkColor(op))
			c.begin()
			c.moveTo(0, 0)
			c.lineTo(w-sz, 0)
			c.lineTo(w, sz)
			c.lineTo(sz, sz)
			c.close()
			c.fill()
		}
		if op2 != 0 {
			c.setFillAlpha(math.Abs(op2))
			c.setFillColor(darkColor(op2))
			c.begin()
			c.moveTo(0, 0)
			c.lineTo(sz, sz)
			c.lineTo(sz, h)
			c.lineTo(0, h-sz)
			c.close()
			c.fill()
		}
		c.begin()
		c.moveTo(sz, h)
		c.lineTo(sz, sz)
		c.lineTo(0, 0)
		c.moveTo(sz, sz)
		c.lineTo(w, sz)
		c.stroke()
	},
	labelMargins: func(s *shape, r rect) *rect {
		if boundedLbl(s) {
			sz := s.style.num("size", 20)
			return &rect{sz, sz, 0, 0}
		}
		return nil
	},
}

// darkColor is the color of the darkOpacity shading: white when negative.
func darkColor(op float64) string {
	if op < 0 {
		return "#FFFFFF"
	}
	return "#000000"
}

// isoCubeShape ports IsoCubeShape: an isometric cube.
var isoCubeShape = cylinderShape(func(s *shape, c *c2d, w, h float64, fg bool) {
	m := math.Min(w, h/(0.5+tan30))
	if fg {
		c.moveTo(0, 0.25*m)
		c.lineTo(0.5*m, (0.5-tan30Dx)*m)
		c.lineTo(m, 0.25*m)
		c.moveTo(0.5*m, (0.5-tan30Dx)*m)
		c.lineTo(0.5*m, (1-tan30Dx)*m)
		return
	}
	// the translation stays for the foreground
	c.translate((w-m)/2, (h-m)/2)
	c.moveTo(0, 0.25*m)
	c.lineTo(0.5*m, m*tan30Dx)
	c.lineTo(m, 0.25*m)
	c.lineTo(m, 0.75*m)
	c.lineTo(0.5*m, (1-tan30Dx)*m)
	c.lineTo(0, 0.75*m)
	c.close()
})

// isoCube2Shape ports IsoCubeShape2: an isometric cube filling the bounds
// (isoAngle sets the angle of the top).
var isoCube2Shape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	isoAngle := clamp(s.style.num("isoAngle", 15), 0.01, 94) * math.Pi / 200
	isoH := math.Min(w*math.Tan(isoAngle), h*0.5)
	c.translate(x, y)
	c.begin()
	c.moveTo(w*0.5, 0)
	c.lineTo(w, isoH)
	c.lineTo(w, h-isoH)
	c.lineTo(w*0.5, h)
	c.lineTo(0, h-isoH)
	c.lineTo(0, isoH)
	c.close()
	c.fillAndStroke()
	c.begin()
	c.moveTo(0, isoH)
	c.lineTo(w*0.5, 2*isoH)
	c.lineTo(w, isoH)
	c.moveTo(w*0.5, 2*isoH)
	c.lineTo(w*0.5, h)
	c.stroke()
}}

// dataStoreShape ports DataStoreShape ("datastore"): a cylinder with three
// rings at the top.
var dataStoreShape = func() *shapeDef {
	d := cylinderShape(func(s *shape, c *c2d, w, h float64, fg bool) {
		dy := math.Min(h/2, jsRound(h/8)+s.strokewidth-1)
		if fg && s.fill != "" || !fg && s.fill == "" {
			c.moveTo(0, dy)
			c.curveTo(0, 2*dy, w, 2*dy, w, dy)
			if !fg {
				c.stroke()
				c.begin()
			}
			c.translate(0, dy/2)
			c.moveTo(0, dy)
			c.curveTo(0, 2*dy, w, 2*dy, w, dy)
			if !fg {
				c.stroke()
				c.begin()
			}
			c.translate(0, dy/2)
			c.moveTo(0, dy)
			c.curveTo(0, 2*dy, w, 2*dy, w, dy)
			if !fg {
				c.stroke()
				c.begin()
			}
			c.translate(0, -dy)
		}
		if !fg {
			c.moveTo(0, dy)
			c.curveTo(0, -dy/3, w, -dy/3, w, dy)
			c.lineTo(w, h-dy)
			c.curveTo(w, h+dy/3, 0, h+dy/3, 0, h-dy)
			c.close()
		}
	})
	d.labelMargins = func(s *shape, r rect) *rect {
		return &rect{0, 2.5 * math.Min(r.h/2, jsRound(r.h/8)+s.strokewidth-1), 0, 0}
	}
	return d
}()

// paintNote paints NoteShape: a sheet with a folded top right corner.
func paintNote(s *shape, c *c2d, x, y, w, h float64) {
	sz := clamp(s.style.num("size", 30), 0, math.Min(w, h))
	op := clamp(s.style.num("darkOpacity", 0), -1, 1)
	c.translate(x, y)
	c.begin()
	c.moveTo(0, 0)
	c.lineTo(w-sz, 0)
	c.lineTo(w, sz)
	c.lineTo(w, h)
	c.lineTo(0, h)
	c.lineTo(0, 0)
	c.close()
	c.fillAndStroke()
	if op != 0 {
		c.setFillAlpha(math.Abs(op))
		c.setFillColor(darkColor(op))
		c.begin()
		c.moveTo(w-sz, 0)
		c.lineTo(w-sz, sz)
		c.lineTo(w, sz)
		c.close()
		c.fill()
	}
	c.begin()
	c.moveTo(w-sz, 0)
	c.lineTo(w-sz, sz)
	c.lineTo(w, sz)
	c.stroke()
}

// noteShape ports NoteShape (an mxCylinder, whose label margins it keeps).
var noteShape = &shapeDef{paintVertex: paintNote, labelMargins: cylinderLabelMargins}

// note2Shape ports NoteShape2, which only changes the label margins.
var note2Shape = &shapeDef{
	paintVertex: paintNote,
	labelMargins: func(s *shape, r rect) *rect {
		if boundedLbl(s) {
			sz := s.style.num("size", 15)
			return &rect{0, math.Min(r.h, sz), 0, math.Max(0, sz)}
		}
		return nil
	},
}

// cylinder2Shape ports CylinderShape ("cylinder2", legacy): a cylinder
// drawn with elliptic arcs.
var cylinder2Shape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	sz := math.Max(0, math.Min(h*0.5, s.style.num("size", 15)))
	c.translate(x, y)
	if sz == 0 {
		c.rect(0, 0, w, h)
		c.fillAndStroke()
		return
	}
	c.begin()
	c.moveTo(0, sz)
	c.arcTo(w*0.5, sz, 0, false, true, w*0.5, 0)
	c.arcTo(w*0.5, sz, 0, false, true, w, sz)
	c.lineTo(w, h-sz)
	c.arcTo(w*0.5, sz, 0, false, true, w*0.5, h)
	c.arcTo(w*0.5, sz, 0, false, true, 0, h-sz)
	c.close()
	c.fillAndStroke()
	c.begin()
	c.moveTo(w, sz)
	c.arcTo(w*0.5, sz, 0, false, true, w*0.5, 2*sz)
	c.arcTo(w*0.5, sz, 0, false, true, 0, sz)
	c.stroke()
}}

// cylinder3Shape ports CylinderShape3: cylinder2 with an optional lid
// (lid=0 draws an open top) and its own label margins.
var cylinder3Shape = &shapeDef{
	paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		sz := math.Max(0, math.Min(h*0.5, s.style.num("size", 15)))
		lid := styleBool(s.style, "lid", true)
		c.translate(x, y)
		if sz == 0 {
			c.rect(0, 0, w, h)
			c.fillAndStroke()
			return
		}
		c.begin()
		if lid {
			c.moveTo(0, sz)
			c.arcTo(w*0.5, sz, 0, false, true, w*0.5, 0)
			c.arcTo(w*0.5, sz, 0, false, true, w, sz)
		} else {
			c.moveTo(0, 0)
			c.arcTo(w*0.5, sz, 0, false, false, w*0.5, sz)
			c.arcTo(w*0.5, sz, 0, false, false, w, 0)
		}
		c.lineTo(w, h-sz)
		c.arcTo(w*0.5, sz, 0, false, true, w*0.5, h)
		c.arcTo(w*0.5, sz, 0, false, true, 0, h-sz)
		c.close()
		c.fillAndStroke()
		if lid {
			c.begin()
			c.moveTo(w, sz)
			c.arcTo(w*0.5, sz, 0, false, true, w*0.5, 2*sz)
			c.arcTo(w*0.5, sz, 0, false, true, 0, sz)
			c.stroke()
		}
	},
	labelMargins: func(s *shape, r rect) *rect {
		if !boundedLbl(s) {
			return nil
		}
		sz := s.style.num("size", 15)
		if !styleBool(s.style, "lid", true) {
			sz /= 2
		}
		return &rect{0, math.Min(r.h, sz*2), 0, math.Max(0, sz*0.3)}
	},
}

// folderArc is the corner radius of FolderShape (arcSize is a fraction of
// the smaller side unless absoluteArcSize is set, 0 unless rounded).
func folderArc(s *shape, w, h, dy float64) float64 {
	arc := s.style.num("arcSize", 0.1)
	if !styleBool(s.style, "absoluteArcSize", false) {
		arc = math.Min(w, h) * arc
	}
	return math.Min(arc, math.Min(w*0.5, (h-dy)*0.5))
}

// folderShape ports FolderShape: a rectangle with a tab (tabWidth,
// tabHeight, tabPosition=left or right).
var folderShape = &shapeDef{
	paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.translate(x, y)
		dx := clamp(s.style.num("tabWidth", 60), 0, w)
		dy := clamp(s.style.num("tabHeight", 20), 0, h)
		tp := s.style.get("tabPosition", "right")
		rounded := styleBool(s.style, "rounded", false)
		arc := folderArc(s, w, h, dy)
		dx = math.Max(dx, arc)
		dx = math.Min(w-arc, dx)
		if !rounded {
			arc = 0
		}
		c.begin()
		if tp == "left" {
			c.moveTo(math.Max(arc, 0), dy)
			c.lineTo(math.Max(arc, 0), 0)
			c.lineTo(dx, 0)
			c.lineTo(dx, dy)
		} else {
			c.moveTo(w-dx, dy)
			c.lineTo(w-dx, 0)
			c.lineTo(w-math.Max(arc, 0), 0)
			c.lineTo(w-math.Max(arc, 0), dy)
		}
		if rounded {
			c.moveTo(0, arc+dy)
			c.arcTo(arc, arc, 0, false, true, arc, dy)
			c.lineTo(w-arc, dy)
			c.arcTo(arc, arc, 0, false, true, w, arc+dy)
			c.lineTo(w, h-arc)
			c.arcTo(arc, arc, 0, false, true, w-arc, h)
			c.lineTo(arc, h)
			c.arcTo(arc, arc, 0, false, true, 0, h-arc)
		} else {
			c.moveTo(0, dy)
			c.lineTo(w, dy)
			c.lineTo(w, h)
			c.lineTo(0, h)
		}
		c.close()
		c.fillAndStroke()
		if s.style.get("folderSymbol", "") == "triangle" {
			c.begin()
			c.moveTo(w-30, dy+20)
			c.lineTo(w-20, dy+10)
			c.lineTo(w-10, dy+20)
			c.close()
			c.stroke()
		}
	},
	labelMargins: func(s *shape, r rect) *rect {
		if !boundedLbl(s) {
			return nil
		}
		sizeY := s.style.num("tabHeight", 15)
		if !styleBool(s.style, "labelInHeader", false) {
			return &rect{0, math.Min(r.h, sizeY), 0, 0}
		}
		sizeX := s.style.num("tabWidth", 15)
		arc := folderArc(s, r.w, r.h, sizeY)
		if !styleBool(s.style, "rounded", false) {
			arc = 0
		}
		if s.style.get("tabPosition", "right") == "left" {
			return &rect{arc, 0, math.Min(r.w, r.w-sizeX), math.Min(r.h, r.h-sizeY)}
		}
		return &rect{math.Min(r.w, r.w-sizeX), 0, arc, math.Min(r.h, r.h-sizeY)}
	},
	roundable: true,
}

// messageShape ports MessageShape: an envelope.
var messageShape = cylinderShape(func(s *shape, c *c2d, w, h float64, fg bool) {
	if fg {
		c.moveTo(0, 0)
		c.lineTo(w/2, h/2)
		c.lineTo(w, 0)
		return
	}
	c.moveTo(0, 0)
	c.lineTo(w, 0)
	c.lineTo(w, h)
	c.lineTo(0, h)
	c.close()
})

// waypointShape ports WaypointShape: a dot in the stroke color.
var waypointShape = &shapeDef{
	paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.setFillColor(s.stroke)
		sz := math.Max(0, s.style.num("size", 6)-2) + 2*s.strokewidth
		c.ellipse(x+(w-sz)*0.5, y+(h-sz)*0.5, sz, sz)
		c.fill()
		c.setFillColor("")
		c.rect(x, y, w, h)
		c.fill()
	},
	labelMargins: cylinderLabelMargins,
}

// --- mxEllipse and mxRhombus subclasses ---

// tapeDataShape ports TapeDataShape: a circle with a tangent at the bottom.
var tapeDataShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	paintEllipseShape(s, c, x, y, w, h)
	c.begin()
	c.moveTo(x+w/2, y+h)
	c.lineTo(x+w, y+h)
	c.stroke()
}}

// orEllipseShape ports OrEllipseShape: a circle with a plus.
var orEllipseShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	paintEllipseShape(s, c, x, y, w, h)
	c.begin()
	c.moveTo(x, y+h/2)
	c.lineTo(x+w, y+h/2)
	c.stroke()
	c.begin()
	c.moveTo(x+w/2, y)
	c.lineTo(x+w/2, y+h)
	c.stroke()
}}

// sumEllipseShape ports SumEllipseShape: a circle with an x.
var sumEllipseShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	paintEllipseShape(s, c, x, y, w, h)
	const s2 = 0.145
	c.begin()
	c.moveTo(x+w*s2, y+h*s2)
	c.lineTo(x+w*(1-s2), y+h*(1-s2))
	c.stroke()
	c.begin()
	c.moveTo(x+w*(1-s2), y+h*s2)
	c.lineTo(x+w*s2, y+h*(1-s2))
	c.stroke()
}}

// sortShape ports SortShape: a rhombus with a horizontal line.
var sortShape = &shapeDef{
	paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		paintRhombusShape(s, c, x, y, w, h)
		c.begin()
		c.moveTo(x, y+h/2)
		c.lineTo(x+w, y+h/2)
		c.stroke()
	},
	labelBounds: rhombusLabelBounds,
	roundable:   true,
}

// collateShape ports CollateShape: an hourglass.
var collateShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	c.begin()
	c.moveTo(x, y)
	c.lineTo(x+w, y)
	c.lineTo(x+w/2, y+h/2)
	c.close()
	c.fillAndStroke()
	c.begin()
	c.moveTo(x, y+h)
	c.lineTo(x+w, y+h)
	c.lineTo(x+w/2, y+h/2)
	c.close()
	c.fillAndStroke()
}}

// dimensionShape ports DimensionShape: a dimension line with arrows at
// the bottom between two vertical lines.
var dimensionShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	sw := c.st.strokeWidth / 2
	al := 10 + 2*sw // arrow size
	cy := y + h - al/2
	c.begin()
	c.moveTo(x, y)
	c.lineTo(x, y+h)
	c.moveTo(x+sw, cy)
	c.lineTo(x+sw+al, cy-al/2)
	c.moveTo(x+sw, cy)
	c.lineTo(x+sw+al, cy+al/2)
	c.moveTo(x+sw, cy)
	c.lineTo(x+w-sw, cy)
	// opposite side
	c.moveTo(x+w, y)
	c.lineTo(x+w, y+h)
	c.moveTo(x+w-sw, cy)
	c.lineTo(x+w-al-sw, cy-al/2)
	c.moveTo(x+w-sw, cy)
	c.lineTo(x+w-al-sw, cy+al/2)
	c.stroke()
}}

// partialRectangleShape ports PartialRectangleShape: a filled rectangle
// whose sides are stroked when top, left, right and bottom are 1.
var partialRectangleShape = &shapeDef{paintVertex: paintPartialRectangle}

func paintPartialRectangle(s *shape, c *c2d, x, y, w, h float64) {
	c.setStrokeColor("")
	top := s.style.get("top", "1") == "1"
	left := s.style.get("left", "1") == "1"
	right := s.style.get("right", "1") == "1"
	bottom := s.style.get("bottom", "1") == "1"
	// drawHidden is true: the rectangle is always painted
	c.rect(x, y, w, h)
	c.fill()
	c.setStrokeColor(s.stroke)
	c.setLineCap("square")
	c.begin()
	c.moveTo(x, y)
	if top {
		c.lineTo(x+w, y)
	} else {
		c.moveTo(x+w, y)
	}
	if right {
		c.lineTo(x+w, y+h)
	} else {
		c.moveTo(x+w, y+h)
	}
	if bottom {
		c.lineTo(x, y+h)
	} else {
		c.moveTo(x, y+h)
	}
	if left {
		c.lineTo(x, y)
	}
	c.stroke()
	c.setLineCap("flat")
	if s.fill != "" {
		paintTableCellLines(s, c, s.st.cell, x, y, w, h, s.stroke, s.strokewidth)
	}
}

// lineEllipseShape ports LineEllipseShape: a circle with a horizontal (or
// line=vertical) line.
var lineEllipseShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	paintEllipseShape(s, c, x, y, w, h)
	c.begin()
	if s.style.get("line", "") == "vertical" {
		c.moveTo(x+w/2, y)
		c.lineTo(x+w/2, y+h)
	} else {
		c.moveTo(x, y+h/2)
		c.lineTo(x+w, y+h/2)
	}
	c.stroke()
}}

// stateShape ports StateShape ("endState", outer stroke) and
// StartStateShape ("startState"): a filled circle, inset in a ring. They
// keep mxDoubleEllipse's label bounds.
func stateShape(outerStroke bool) *shapeDef {
	return &shapeDef{
		paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
			inset := math.Min(4, math.Min(w/5, h/5))
			if w > 0 && h > 0 {
				c.ellipse(x+inset, y+inset, w-2*inset, h-2*inset)
				c.fillAndStroke()
			}
			if outerStroke {
				c.ellipse(x, y, w, h)
				c.stroke()
			}
		},
		labelBounds: doubleEllipseLabelBounds,
	}
}

// smileyFaceShape ports SmileyFaceShape (smileyType=happy, neutral or sad;
// smileyFeatureColor for the eyes and mouth).
var smileyFaceShape = &shapeDef{paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
	smileyType := s.style.get("smileyType", "happy")
	featureColor := s.style.get("smileyFeatureColor", "#666666")
	c.translate(x, y)
	r := math.Min(w, h) / 2
	c.ellipse(w/2-r, h/2-r, r*2, r*2)
	c.fillAndStroke()
	sc := math.Min(w, h) / 30
	cx, cy := w/2, h/2
	c.setFillColor(featureColor)
	c.setStrokeColor(featureColor)
	c.setStrokeWidth(2 * sc)
	eyeR := 1.5 * sc
	eyeXOffset := 5 * sc
	eyeYOffset := 5 * sc
	c.ellipse(cx-eyeXOffset-eyeR, cy-eyeYOffset-eyeR, eyeR*2, eyeR*2)
	c.fillAndStroke()
	c.ellipse(cx+eyeXOffset-eyeR, cy-eyeYOffset-eyeR, eyeR*2, eyeR*2)
	c.fillAndStroke()
	c.setStrokeWidth(1 * sc)
	c.setFillColor("#000000")
	switch smileyType {
	case "happy":
		mcx, mcy := cx, cy+2*sc
		c.begin()
		c.moveTo(mcx+7.5*sc, mcy)
		c.arcTo(7.5*sc, 7.5*sc, 0, true, true, mcx-7.5*sc, mcy)
		c.lineTo(mcx-6.818*sc, mcy)
		c.arcTo(6.818*sc, 6.818*sc, 0, true, false, mcx+6.818*sc, mcy)
		c.close()
		c.fillAndStroke()
	case "sad":
		mcx, mcy := cx, cy+7*sc
		c.begin()
		c.moveTo(mcx-7.5*sc, mcy)
		c.arcTo(7.5*sc, 7.5*sc, 0, true, true, mcx+7.5*sc, mcy)
		c.lineTo(mcx+6.818*sc, mcy)
		c.arcTo(6.818*sc, 6.818*sc, 0, true, false, mcx-6.818*sc, mcy)
		c.close()
		c.fillAndStroke()
	default:
		c.begin()
		c.moveTo(cx-5*sc, cy+7*sc)
		c.lineTo(cx+5*sc, cy+7*sc)
		c.stroke()
	}
}}

// maxZigzagTeeth bounds the work for absurd widths.
const maxZigzagTeeth = 100000

// zigzagShape ports ZigzagShape: a zigzag (or with rounded a wave) line
// across the middle; size is the half wavelength, the height the amplitude.
var zigzagShape = &shapeDef{
	paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
		c.translate(x, y)
		if fill := colorOrNone(s.style.get("fillColor", "")); fill != "" {
			c.setStrokeColor("")
			c.begin()
			c.rect(0, 0, w, h)
			c.fillAndStroke()
			c.setStrokeColor(s.stroke)
		}
		size := math.Max(5, s.style.num("size", 10))
		centerY := h / 2
		// peaks are inset so that the stroke stays in the bounds
		sw := s.strokewidth
		inset := sw
		if s.isRounded {
			inset = sw / 2
		}
		topY, bottomY := inset, h-inset
		numFull := math.Max(1, jsRound(w/size)-1)
		if !finite(numFull) || numFull > maxZigzagTeeth {
			return // degenerate widths
		}
		halfWave := w / (numFull + 1)
		halfEnd := halfWave / 2
		n := int(numFull)
		c.begin()
		c.moveTo(0, centerY)
		if s.isRounded {
			const k = 0.4
			c.curveTo(k*halfEnd, centerY-(centerY-topY)*k, (1-k)*halfEnd, topY, halfEnd, topY)
			for j := 0; j < n; j++ {
				sx := halfEnd + float64(j)*halfWave
				ex := sx + halfWave
				sy, ey := topY, bottomY
				if j%2 != 0 {
					sy, ey = bottomY, topY
				}
				c.curveTo(sx+k*halfWave, sy, ex-k*halfWave, ey, ex, ey)
			}
			lastX := halfEnd + numFull*halfWave
			lastY := topY
			if n%2 != 0 {
				lastY = bottomY
			}
			dirSign := -1.0
			if lastY == topY {
				dirSign = 1
			}
			c.curveTo(lastX+k*halfEnd, lastY, w-k*halfEnd, centerY-dirSign*(centerY-topY)*k, w, centerY)
		} else {
			c.lineTo(halfEnd, topY)
			for j := 0; j < n; j++ {
				ey := bottomY
				if j%2 != 0 {
					ey = topY
				}
				c.lineTo(halfEnd+float64(j+1)*halfWave, ey)
			}
			c.lineTo(w, centerY)
		}
		c.stroke()
	},
	labelBounds: rectangleLabelBounds,
	roundable:   true,
}
