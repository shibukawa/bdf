package drawio

import (
	"math"
	"strconv"
)

// Perimeters (mxPerimeter and draw.io's additions in Shapes.js): where an
// edge heading from next towards the center of a vertex meets its
// outline. bounds are the vertex's perimeter bounds; with orthogonal the
// point is the orthogonal projection of next where one exists. A nil
// result stands for none (the edge then ends at the center).

// perimeterFunc is an mxPerimeter function.
type perimeterFunc func(bounds rect, vertex *cellState, next point, orthogonal bool) *point

// perimeters maps perimeter style values to perimeters (mxStyleRegistry).
var perimeters = map[string]perimeterFunc{
	"rectanglePerimeter":     rectanglePerimeter,
	"ellipsePerimeter":       ellipsePerimeter,
	"rhombusPerimeter":       rhombusPerimeter,
	"trianglePerimeter":      trianglePerimeter,
	"hexagonPerimeter":       hexagonPerimeter,
	"hexagonPerimeter2":      hexagonPerimeter2,
	"calloutPerimeter":       calloutPerimeter,
	"parallelogramPerimeter": parallelogramPerimeter,
	"trapezoidPerimeter":     trapezoidPerimeter,
	"stepPerimeter":          stepPerimeter,
	"lifelinePerimeter":      lifelinePerimeter,
	"orthogonalPerimeter":    orthogonalPerimeter,
	"backbonePerimeter":      backbonePerimeter,
	"centerPerimeter":        centerPerimeter,
}

// Default sizes of the shapes whose perimeters depend on them (the
// prototype sizes in Shapes.js).
const (
	calloutSize            = 30   // CalloutShape.prototype.size
	parallelogramSize      = 0.2  // ParallelogramShape.prototype.size
	parallelogramFixedSize = 20   // ParallelogramShape.prototype.fixedSize
	trapezoidSize          = 0.2  // TrapezoidShape.prototype.size
	trapezoidFixedSize     = 20   // TrapezoidShape.prototype.fixedSize
	stepSize               = 0.2  // StepShape.prototype.size
	stepFixedSize          = 20   // StepShape.prototype.fixedSize
	hexagonSize            = 0.25 // HexagonShape.prototype.size
	hexagonFixedSize       = 20   // HexagonShape.prototype.fixedSize
	umlLifelineSize        = 40   // UmlLifeline.prototype.size
)

// intersectionPt returns the intersection of two segments or nil
// (mxUtils.intersection).
func intersectionPt(x0, y0, x1, y1, x2, y2, x3, y3 float64) *point {
	if p, ok := intersection(x0, y0, x1, y1, x2, y2, x3, y3); ok {
		return &p
	}
	return nil
}

// polygonPerimeterPoint returns the intersection of a polygon's outline
// with the line from center to p that is nearest to p
// (mxUtils.getPerimeterPoint).
func polygonPerimeterPoint(pts []point, center, p point) *point {
	var best *point
	bestDist := 0.0
	for i := 0; i < len(pts)-1; i++ {
		ip := intersectionPt(pts[i].x, pts[i].y, pts[i+1].x, pts[i+1].y, center.x, center.y, p.x, p.y)
		if ip == nil {
			continue
		}
		dx, dy := p.x-ip.x, p.y-ip.y
		if d := dy*dy + dx*dx; best == nil || bestDist > d {
			best, bestDist = ip, d
		}
	}
	return best
}

// jsParseInt is JavaScript's parseInt of a number: it truncates, except
// that numbers written in exponential notation (below 1e-6 or from 1e21)
// give their first digit.
func jsParseInt(v float64) float64 {
	a := math.Abs(v)
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return math.NaN()
	}
	if a != 0 && (a < 1e-6 || a >= 1e21) {
		d := float64(strconv.FormatFloat(a, 'e', -1, 64)[0] - '0')
		return math.Copysign(d, v)
	}
	return math.Trunc(v)
}

// rectanglePerimeter (mxPerimeter.RectanglePerimeter).
func rectanglePerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	cx, cy := bounds.cx(), bounds.cy()
	dx, dy := next.x-cx, next.y-cy
	alpha := math.Atan2(dy, dx)
	var p point
	pi, pi2 := math.Pi, math.Pi/2
	beta := pi2 - alpha
	t := math.Atan2(bounds.h, bounds.w)
	switch {
	case alpha < -pi+t || alpha > pi-t: // left edge
		p.x = bounds.x
		p.y = cy - bounds.w*math.Tan(alpha)/2
	case alpha < -t: // top edge
		p.y = bounds.y
		p.x = cx - bounds.h*math.Tan(beta)/2
	case alpha < t: // right edge
		p.x = bounds.x + bounds.w
		p.y = cy + bounds.w*math.Tan(alpha)/2
	default: // bottom edge
		p.y = bounds.y + bounds.h
		p.x = cx + bounds.h*math.Tan(beta)/2
	}
	if orthogonal {
		if next.x >= bounds.x && next.x <= bounds.x+bounds.w {
			p.x = next.x
		} else if next.y >= bounds.y && next.y <= bounds.y+bounds.h {
			p.y = next.y
		}
		if next.x < bounds.x {
			p.x = bounds.x
		} else if next.x > bounds.x+bounds.w {
			p.x = bounds.x + bounds.w
		}
		if next.y < bounds.y {
			p.y = bounds.y
		} else if next.y > bounds.y+bounds.h {
			p.y = bounds.y + bounds.h
		}
	}
	return &p
}

// ellipsePerimeter (mxPerimeter.EllipsePerimeter). The slope of the line
// comes from the offsets truncated to integers, as in mxGraph.
func ellipsePerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	x, y := bounds.x, bounds.y
	a, b := bounds.w/2, bounds.h/2
	cx, cy := x+a, y+b
	px, py := next.x, next.y
	// the straight line through the point and the center, y = d * x + h
	dx := jsParseInt(px - cx)
	dy := jsParseInt(py - cy)
	if dx == 0 && dy != 0 {
		return &point{cx, cy + b*dy/math.Abs(dy)}
	} else if dx == 0 && dy == 0 {
		return &point{px, py}
	}
	if orthogonal {
		if py >= y && py <= y+bounds.h {
			ty := py - cy
			tx := orZero(math.Sqrt(a * a * (1 - (ty*ty)/(b*b))))
			if px <= x {
				tx = -tx
			}
			return &point{cx + tx, py}
		}
		if px >= x && px <= x+bounds.w {
			tx := px - cx
			ty := orZero(math.Sqrt(b * b * (1 - (tx*tx)/(a*a))))
			if py <= y {
				ty = -ty
			}
			return &point{px, cy + ty}
		}
	}
	// the intersection
	d := dy / dx
	h := cy - d*cx
	e := a*a*d*d + b*b
	f := -2 * cx * e
	g := a*a*d*d*cx*cx + b*b*cx*cx - a*a*b*b
	det := math.Sqrt(f*f - 4*e*g)
	// two solutions (perimeter points)
	xout1 := (-f + det) / (2 * e)
	xout2 := (-f - det) / (2 * e)
	yout1 := d*xout1 + h
	yout2 := d*xout2 + h
	dist1 := math.Sqrt(math.Pow(xout1-px, 2) + math.Pow(yout1-py, 2))
	dist2 := math.Sqrt(math.Pow(xout2-px, 2) + math.Pow(yout2-py, 2))
	// the correct solution
	if dist1 < dist2 {
		return &point{xout1, yout1}
	}
	return &point{xout2, yout2}
}

// orZero is JavaScript's v || 0 for numbers.
func orZero(v float64) float64 {
	if math.IsNaN(v) {
		return 0
	}
	return v
}

// rhombusPerimeter (mxPerimeter.RhombusPerimeter).
func rhombusPerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	x, y, w, h := bounds.x, bounds.y, bounds.w, bounds.h
	cx, cy := x+w/2, y+h/2
	px, py := next.x, next.y
	// special case for intersecting the diamond's corners
	if cx == px {
		if cy > py {
			return &point{cx, y} // top
		}
		return &point{cx, y + h} // bottom
	} else if cy == py {
		if cx > px {
			return &point{x, cy} // left
		}
		return &point{x + w, cy} // right
	}
	tx, ty := cx, cy
	if orthogonal {
		if px >= x && px <= x+w {
			tx = px
		} else if py >= y && py <= y+h {
			ty = py
		}
	}
	// in which quadrant will the intersection be?
	if px < cx {
		if py < cy {
			return intersectionPt(px, py, tx, ty, cx, y, x, cy)
		}
		return intersectionPt(px, py, tx, ty, cx, y+h, x, cy)
	} else if py < cy {
		return intersectionPt(px, py, tx, ty, cx, y, x+w, cy)
	}
	return intersectionPt(px, py, tx, ty, cx, y+h, x+w, cy)
}

// trianglePerimeter (mxPerimeter.TrianglePerimeter): the triangle points
// east unless the vertex's direction says otherwise.
func trianglePerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	direction := ""
	if vertex != nil {
		direction = vertex.style["direction"]
	}
	vertical := direction == "north" || direction == "south"
	x, y, w, h := bounds.x, bounds.y, bounds.w, bounds.h
	cx, cy := x+w/2, y+h/2
	start := point{x, y}
	corner := point{x + w, cy}
	end := point{x, y + h}
	switch direction {
	case "north":
		start = end
		corner = point{cx, y}
		end = point{x + w, y + h}
	case "south":
		corner = point{cx, y + h}
		end = point{x + w, y}
	case "west":
		start = point{x + w, y}
		corner = point{x, cy}
		end = point{x + w, y + h}
	}
	dx, dy := next.x-cx, next.y-cy
	alpha := math.Atan2(dy, dx)
	t := math.Atan2(h, w)
	if vertical {
		alpha = math.Atan2(dx, dy)
		t = math.Atan2(w, h)
	}
	var base bool
	if direction == "north" || direction == "west" {
		base = alpha > -t && alpha < t
	} else {
		base = alpha < -math.Pi+t || alpha > math.Pi-t
	}
	var result *point
	if base {
		if orthogonal && ((vertical && next.x >= start.x && next.x <= end.x) ||
			(!vertical && next.y >= start.y && next.y <= end.y)) {
			if vertical {
				result = &point{next.x, start.y}
			} else {
				result = &point{start.x, next.y}
			}
		} else {
			switch direction {
			case "north":
				result = &point{x + w/2 + h*math.Tan(alpha)/2, y + h}
			case "south":
				result = &point{x + w/2 - h*math.Tan(alpha)/2, y}
			case "west":
				result = &point{x + w, y + h/2 + w*math.Tan(alpha)/2}
			default:
				result = &point{x, y + h/2 - w*math.Tan(alpha)/2}
			}
		}
	} else {
		if orthogonal {
			pt := point{cx, cy}
			if next.y >= y && next.y <= y+h {
				switch {
				case vertical:
					pt.x = cx
				case direction == "west":
					pt.x = x + w
				default:
					pt.x = x
				}
				pt.y = next.y
			} else if next.x >= x && next.x <= x+w {
				pt.x = next.x
				switch {
				case !vertical:
					pt.y = cy
				case direction == "north":
					pt.y = y + h
				default:
					pt.y = y
				}
			}
			cx, cy = pt.x, pt.y
		}
		if (vertical && next.x <= x+w/2) || (!vertical && next.y <= y+h/2) {
			result = intersectionPt(next.x, next.y, cx, cy, start.x, start.y, corner.x, corner.y)
		} else {
			result = intersectionPt(next.x, next.y, cx, cy, corner.x, corner.y, end.x, end.y)
		}
	}
	if result == nil {
		result = &point{cx, cy}
	}
	return result
}

// hexagonPerimeter (mxPerimeter.HexagonPerimeter): a hexagon with corners
// at a quarter of the width (or height, for north and south directions).
func hexagonPerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	x, y, w, h := bounds.x, bounds.y, bounds.w, bounds.h
	cx, cy := bounds.cx(), bounds.cy()
	px, py := next.x, next.y
	dx, dy := px-cx, py-cy
	alpha := -math.Atan2(dy, dx)
	pi, pi2 := math.Pi, math.Pi/2
	result := &point{cx, cy}
	direction := "east"
	if vertex != nil {
		direction = vertex.style.get("direction", "east")
	}
	vertical := direction == "north" || direction == "south"
	var a, b point
	fl := math.Floor

	// only consider correct quadrants for the orthogonal case
	if (px < x) && (py < y) || (px < x) && (py > y+h) || (px > x+w) && (py < y) || (px > x+w) && (py > y+h) {
		orthogonal = false
	}

	if orthogonal {
		if vertical {
			// special cases where intersects with hexagon corners
			if px == cx {
				if py <= y {
					return &point{cx, y}
				} else if py >= y+h {
					return &point{cx, y + h}
				}
			} else if px < x {
				if py == y+h/4 {
					return &point{x, y + h/4}
				} else if py == y+3*h/4 {
					return &point{x, y + 3*h/4}
				}
			} else if px > x+w {
				if py == y+h/4 {
					return &point{x + w, y + h/4}
				} else if py == y+3*h/4 {
					return &point{x + w, y + 3*h/4}
				}
			} else if px == x {
				if py < cy {
					return &point{x, y + h/4}
				} else if py > cy {
					return &point{x, y + 3*h/4}
				}
			} else if px == x+w {
				if py < cy {
					return &point{x + w, y + h/4}
				} else if py > cy {
					return &point{x + w, y + 3*h/4}
				}
			}
			if py == y {
				return &point{cx, y}
			} else if py == y+h {
				return &point{cx, y + h}
			}
			if px < cx {
				if (py > y+h/4) && (py < y+3*h/4) {
					a = point{x, y}
					b = point{x, y + h}
				} else if py < y+h/4 {
					a = point{x - fl(0.5*w), y + fl(0.5*h)}
					b = point{x + w, y - fl(0.25*h)}
				} else if py > y+3*h/4 {
					a = point{x - fl(0.5*w), y + fl(0.5*h)}
					b = point{x + w, y + fl(1.25*h)}
				}
			} else if px > cx {
				if (py > y+h/4) && (py < y+3*h/4) {
					a = point{x + w, y}
					b = point{x + w, y + h}
				} else if py < y+h/4 {
					a = point{x, y - fl(0.25*h)}
					b = point{x + fl(1.5*w), y + fl(0.5*h)}
				} else if py > y+3*h/4 {
					a = point{x + fl(1.5*w), y + fl(0.5*h)}
					b = point{x, y + fl(1.25*h)}
				}
			}
		} else {
			// special cases where intersects with hexagon corners
			if py == cy {
				if px <= x {
					return &point{x, y + h/2}
				} else if px >= x+w {
					return &point{x + w, y + h/2}
				}
			} else if py < y {
				if px == x+w/4 {
					return &point{x + w/4, y}
				} else if px == x+3*w/4 {
					return &point{x + 3*w/4, y}
				}
			} else if py > y+h {
				if px == x+w/4 {
					return &point{x + w/4, y + h}
				} else if px == x+3*w/4 {
					return &point{x + 3*w/4, y + h}
				}
			} else if py == y {
				if px < cx {
					return &point{x + w/4, y}
				} else if px > cx {
					return &point{x + 3*w/4, y}
				}
			} else if py == y+h {
				if px < cx {
					return &point{x + w/4, y + h}
				} else if py > cy {
					return &point{x + 3*w/4, y + h}
				}
			}
			if px == x {
				return &point{x, cy}
			} else if px == x+w {
				return &point{x + w, cy}
			}
			if py < cy {
				if (px > x+w/4) && (px < x+3*w/4) {
					a = point{x, y}
					b = point{x + w, y}
				} else if px < x+w/4 {
					a = point{x - fl(0.25*w), y + h}
					b = point{x + fl(0.5*w), y - fl(0.5*h)}
				} else if px > x+3*w/4 {
					a = point{x + fl(0.5*w), y - fl(0.5*h)}
					b = point{x + fl(1.25*w), y + h}
				}
			} else if py > cy {
				if (px > x+w/4) && (px < x+3*w/4) {
					a = point{x, y + h}
					b = point{x + w, y + h}
				} else if px < x+w/4 {
					a = point{x - fl(0.25*w), y}
					b = point{x + fl(0.5*w), y + fl(1.5*h)}
				} else if px > x+3*w/4 {
					a = point{x + fl(0.5*w), y + fl(1.5*h)}
					b = point{x + fl(1.25*w), y}
				}
			}
		}
		tx, ty := cx, cy
		if px >= x && px <= x+w {
			tx = px
			if py < cy {
				ty = y + h
			} else {
				ty = y
			}
		} else if py >= y && py <= y+h {
			ty = py
			if px < cx {
				tx = x + w
			} else {
				tx = x
			}
		}
		result = intersectionPt(tx, ty, next.x, next.y, a.x, a.y, b.x, b.y)
	} else {
		if vertical {
			beta := math.Atan2(h/4, w/2)
			// special cases where intersects with hexagon corners
			switch alpha {
			case beta:
				return &point{x + w, y + fl(0.25*h)}
			case pi2:
				return &point{x + fl(0.5*w), y}
			case pi - beta:
				return &point{x, y + fl(0.25*h)}
			case -beta:
				return &point{x + w, y + fl(0.75*h)}
			case -pi2:
				return &point{x + fl(0.5*w), y + h}
			case -pi + beta:
				return &point{x, y + fl(0.75*h)}
			}
			if (alpha < beta) && (alpha > -beta) {
				a = point{x + w, y}
				b = point{x + w, y + h}
			} else if (alpha > beta) && (alpha < pi2) {
				a = point{x, y - fl(0.25*h)}
				b = point{x + fl(1.5*w), y + fl(0.5*h)}
			} else if (alpha > pi2) && (alpha < (pi - beta)) {
				a = point{x - fl(0.5*w), y + fl(0.5*h)}
				b = point{x + w, y - fl(0.25*h)}
			} else if ((alpha > (pi - beta)) && (alpha <= pi)) || ((alpha < (-pi + beta)) && (alpha >= -pi)) {
				a = point{x, y}
				b = point{x, y + h}
			} else if (alpha < -beta) && (alpha > -pi2) {
				a = point{x + fl(1.5*w), y + fl(0.5*h)}
				b = point{x, y + fl(1.25*h)}
			} else if (alpha < -pi2) && (alpha > (-pi + beta)) {
				a = point{x - fl(0.5*w), y + fl(0.5*h)}
				b = point{x + w, y + fl(1.25*h)}
			}
		} else {
			beta := math.Atan2(h/2, w/4)
			// special cases where intersects with hexagon corners
			switch alpha {
			case beta:
				return &point{x + fl(0.75*w), y}
			case pi - beta:
				return &point{x + fl(0.25*w), y}
			case pi, -pi:
				return &point{x, y + fl(0.5*h)}
			case 0:
				return &point{x + w, y + fl(0.5*h)}
			case -beta:
				return &point{x + fl(0.75*w), y + h}
			case -pi + beta:
				return &point{x + fl(0.25*w), y + h}
			}
			if (alpha > 0) && (alpha < beta) {
				a = point{x + fl(0.5*w), y - fl(0.5*h)}
				b = point{x + fl(1.25*w), y + h}
			} else if (alpha > beta) && (alpha < (pi - beta)) {
				a = point{x, y}
				b = point{x + w, y}
			} else if (alpha > (pi - beta)) && (alpha < pi) {
				a = point{x - fl(0.25*w), y + h}
				b = point{x + fl(0.5*w), y - fl(0.5*h)}
			} else if (alpha < 0) && (alpha > -beta) {
				a = point{x + fl(0.5*w), y + fl(1.5*h)}
				b = point{x + fl(1.25*w), y}
			} else if (alpha < -beta) && (alpha > (-pi + beta)) {
				a = point{x, y + h}
				b = point{x + w, y + h}
			} else if (alpha < (-pi + beta)) && (alpha > -pi) {
				a = point{x - fl(0.25*w), y}
				b = point{x + fl(0.5*w), y + fl(1.5*h)}
			}
		}
		result = intersectionPt(cx, cy, next.x, next.y, a.x, a.y, b.x, b.y)
	}
	if result == nil {
		return &point{cx, cy}
	}
	return result
}

// shapeSize returns the size of a shape whose size is relative to its
// bounds, or absolute with fixedSize (the size style or the shape's
// default): the size helpers of the parallelogram, trapezoid, step and
// hexagon perimeters in Shapes.js.
func shapeSize(vertex *cellState, def, fixedDef float64) (size float64, fixed bool) {
	// mxUtils.getValue(style, 'fixedSize', '0') != '0'
	if v, ok := vertex.style["fixedSize"]; ok {
		f, num := jsNumber(v)
		fixed = !num || f != 0
	}
	size = def
	if fixed {
		size = fixedDef
	}
	return vertex.style.num("size", size), fixed
}

// orthogonalStart returns the point the perimeter line starts from: the
// center, moved to next's x or y for orthogonal projections (the Shapes.js
// polygon perimeters).
func orthogonalStart(bounds rect, next point, orthogonal bool) point {
	p1 := bounds.center()
	if orthogonal {
		if next.x < bounds.x || next.x > bounds.x+bounds.w {
			p1.y = next.y
		} else {
			p1.x = next.x
		}
	}
	return p1
}

// parallelogramPerimeter (mxPerimeter.ParallelogramPerimeter in Shapes.js).
func parallelogramPerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	size, fixed := shapeSize(vertex, parallelogramSize, parallelogramFixedSize)
	x, y, w, h := bounds.x, bounds.y, bounds.w, bounds.h
	direction := vertex.style.get("direction", "east")
	var pts []point
	if direction == "north" || direction == "south" {
		dy := h * math.Max(0, math.Min(1, size))
		if fixed {
			dy = math.Max(0, math.Min(h, size))
		}
		pts = []point{{x, y}, {x + w, y + dy}, {x + w, y + h}, {x, y + h - dy}, {x, y}}
	} else {
		dx := w * math.Max(0, math.Min(1, size))
		if fixed {
			dx = math.Max(0, math.Min(w*0.5, size))
		}
		pts = []point{{x + dx, y}, {x + w, y}, {x + w - dx, y + h}, {x, y + h}, {x + dx, y}}
	}
	return polygonPerimeterPoint(pts, orthogonalStart(bounds, next, orthogonal), next)
}

// trapezoidPerimeter (mxPerimeter.TrapezoidPerimeter in Shapes.js).
func trapezoidPerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	size, fixed := shapeSize(vertex, trapezoidSize, trapezoidFixedSize)
	x, y, w, h := bounds.x, bounds.y, bounds.w, bounds.h
	var pts []point
	switch vertex.style.get("direction", "east") {
	case "east":
		dx := w * math.Max(0, math.Min(1, size))
		if fixed {
			dx = math.Max(0, math.Min(w*0.5, size))
		}
		pts = []point{{x + dx, y}, {x + w - dx, y}, {x + w, y + h}, {x, y + h}, {x + dx, y}}
	case "west":
		dx := w * math.Max(0, math.Min(1, size))
		if fixed {
			dx = math.Max(0, math.Min(w, size))
		}
		pts = []point{{x, y}, {x + w, y}, {x + w - dx, y + h}, {x + dx, y + h}, {x, y}}
	case "north":
		dy := h * math.Max(0, math.Min(1, size))
		if fixed {
			dy = math.Max(0, math.Min(h, size))
		}
		pts = []point{{x, y + dy}, {x + w, y}, {x + w, y + h}, {x, y + h - dy}, {x, y + dy}}
	default:
		dy := h * math.Max(0, math.Min(1, size))
		if fixed {
			dy = math.Max(0, math.Min(h, size))
		}
		pts = []point{{x, y}, {x + w, y + dy}, {x + w, y + h - dy}, {x, y + h}, {x, y}}
	}
	return polygonPerimeterPoint(pts, orthogonalStart(bounds, next, orthogonal), next)
}

// stepPerimeter (mxPerimeter.StepPerimeter in Shapes.js).
func stepPerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	size, fixed := shapeSize(vertex, stepSize, stepFixedSize)
	x, y, w, h := bounds.x, bounds.y, bounds.w, bounds.h
	cx, cy := bounds.cx(), bounds.cy()
	var pts []point
	switch vertex.style.get("direction", "east") {
	case "east":
		dx := w * math.Max(0, math.Min(1, size))
		if fixed {
			dx = math.Max(0, math.Min(w, size))
		}
		pts = []point{{x, y}, {x + w - dx, y}, {x + w, cy}, {x + w - dx, y + h}, {x, y + h}, {x + dx, cy}, {x, y}}
	case "west":
		dx := w * math.Max(0, math.Min(1, size))
		if fixed {
			dx = math.Max(0, math.Min(w, size))
		}
		pts = []point{{x + dx, y}, {x + w, y}, {x + w - dx, cy}, {x + w, y + h}, {x + dx, y + h}, {x, cy}, {x + dx, y}}
	case "north":
		dy := h * math.Max(0, math.Min(1, size))
		if fixed {
			dy = math.Max(0, math.Min(h, size))
		}
		pts = []point{{x, y + dy}, {cx, y}, {x + w, y + dy}, {x + w, y + h}, {cx, y + h - dy}, {x, y + h}, {x, y + dy}}
	default:
		dy := h * math.Max(0, math.Min(1, size))
		if fixed {
			dy = math.Max(0, math.Min(h, size))
		}
		pts = []point{{x, y}, {cx, y + dy}, {x + w, y}, {x + w, y + h - dy}, {cx, y + h}, {x, y + h - dy}, {x, y}}
	}
	return polygonPerimeterPoint(pts, orthogonalStart(bounds, next, orthogonal), next)
}

// hexagonPerimeter2 is the perimeter of draw.io's hexagon shape
// (mxPerimeter.HexagonPerimeter2 in Shapes.js).
func hexagonPerimeter2(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	size, fixed := shapeSize(vertex, hexagonSize, hexagonFixedSize)
	x, y, w, h := bounds.x, bounds.y, bounds.w, bounds.h
	cx, cy := bounds.cx(), bounds.cy()
	direction := vertex.style.get("direction", "east")
	var pts []point
	if direction == "north" || direction == "south" {
		dy := h * math.Max(0, math.Min(1, size))
		if fixed {
			dy = math.Max(0, math.Min(h, size))
		}
		pts = []point{{cx, y}, {x + w, y + dy}, {x + w, y + h - dy}, {cx, y + h}, {x, y + h - dy}, {x, y + dy}, {cx, y}}
	} else {
		dx := w * math.Max(0, math.Min(1, size))
		if fixed {
			dx = math.Max(0, math.Min(w, size))
		}
		pts = []point{{x + dx, y}, {x + w - dx, y}, {x + w, cy}, {x + w - dx, y + h}, {x + dx, y + h}, {x, cy}, {x + dx, y}}
	}
	return polygonPerimeterPoint(pts, orthogonalStart(bounds, next, orthogonal), next)
}

// calloutPerimeter is the rectangle above the callout's pointer
// (mxPerimeter.CalloutPerimeter in Shapes.js).
func calloutPerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	size := math.Max(0, math.Min(bounds.h, vertex.style.num("size", calloutSize)))
	r := directedBounds(bounds, rect{0, 0, 0, size}, vertex.style,
		styleTruthy(vertex.style, "flipH", false), styleTruthy(vertex.style, "flipV", false))
	return rectanglePerimeter(r, vertex, next, orthogonal)
}

// directedBounds returns r without the margins m (left, top, right and
// bottom as x, y, w and h), turned and flipped with the style's direction
// and flips (mxUtils.getDirectedBounds).
func directedBounds(r rect, m rect, s style, flipH, flipV bool) rect {
	return directedBoundsDir(r, m, s.get("direction", "east"), flipH, flipV)
}

// directedBoundsDir is directedBounds for a direction ("" is east).
func directedBoundsDir(r rect, m rect, d string, flipH, flipV bool) rect {
	if d == "" {
		d = "east"
	}
	m.x = jsRound(math.Max(0, math.Min(r.w, m.x)))
	m.y = jsRound(math.Max(0, math.Min(r.h, m.y)))
	m.w = jsRound(math.Max(0, math.Min(r.w, m.w)))
	m.h = jsRound(math.Max(0, math.Min(r.h, m.h)))
	vertical := d == "south" || d == "north"
	horizontal := d == "east" || d == "west"
	if (flipV && vertical) || (flipH && horizontal) {
		m.x, m.w = m.w, m.x
	}
	if (flipH && vertical) || (flipV && horizontal) {
		m.y, m.h = m.h, m.y
	}
	m2 := m
	switch d {
	case "south":
		m2 = rect{m.h, m.x, m.y, m.w}
	case "west":
		m2 = rect{m.w, m.h, m.x, m.y}
	case "north":
		m2 = rect{m.y, m.w, m.h, m.x}
	}
	return rect{r.x + m2.x, r.y + m2.y, r.w - m2.w - m2.x, r.h - m2.h - m2.y}
}

// lifelinePerimeter keeps connections on the dashed line of a UML lifeline,
// below its head (mxPerimeter.LifelinePerimeter in Shapes.js).
func lifelinePerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	size := float64(umlLifelineSize)
	bottom := bounds.y + bounds.h
	if vertex != nil {
		size = vertex.style.num("size", size)
		// connections stay on the body line, off the mirrored foot box
		if styleIsOne(vertex.style, "lifelineMirror", false) {
			bottom -= size
		}
	}
	sw := vertex.style.num("strokeWidth", 1)/2 - 1
	if next.x < bounds.cx() {
		sw += 1
		sw *= -1
	}
	return &point{bounds.cx() + sw, math.Min(bottom, math.Max(bounds.y+size, next.y))}
}

// orthogonalPerimeter is the rectangle perimeter, always projected
// orthogonally (mxPerimeter.OrthogonalPerimeter in Shapes.js).
func orthogonalPerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	return rectanglePerimeter(bounds, vertex, next, true)
}

// backbonePerimeter keeps connections on the line through the middle of a
// backbone shape (mxPerimeter.BackbonePerimeter in Shapes.js).
func backbonePerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	sw := vertex.style.num("strokeWidth", 1)/2 - 1
	if vertex.style.has("backboneSize") {
		sw += vertex.style.num("backboneSize", 0)/2 - 1
	}
	if d := vertex.style["direction"]; d == "south" || d == "north" {
		if next.x < bounds.cx() {
			sw += 1
			sw *= -1
		}
		return &point{bounds.cx() + sw, math.Min(bounds.y+bounds.h, math.Max(bounds.y, next.y))}
	}
	if next.y < bounds.cy() {
		sw += 1
		sw *= -1
	}
	return &point{math.Min(bounds.x+bounds.w, math.Max(bounds.x, next.x)), bounds.cy() + sw}
}

// centerPerimeter connects at the center (mxPerimeter.CenterPerimeter in Shapes.js).
func centerPerimeter(bounds rect, vertex *cellState, next point, orthogonal bool) *point {
	return &point{bounds.cx(), bounds.cy()}
}
