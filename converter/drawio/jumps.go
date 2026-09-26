package drawio

import "math"

// Line jumps: where an edge with a jumpStyle crosses an edge that comes
// before it in the model, draw.io draws a jump (an arc, a gap, a sharp
// bump or two ticks). The crossings are found when the view is validated
// (mxGraphView.updateLineJumps in Graph.js) and drawn in place of the
// line (draw.io's mxConnector.paintLine override).

// routedPoint is a point of an edge's route; jump marks a crossing.
type routedPoint struct {
	x, y float64
	jump bool
}

// defaultJumpSize is Graph.defaultJumpSize.
const defaultJumpSize = 6

// updateLineJumps computes the routed points of the edges with a jump
// style: each crosses the straight or orthogonal edges before it in the
// model's order (curved edges neither jump nor are jumped over).
func (v *graphView) updateLineJumps() {
	var valid []*cellState
	for _, st := range v.order {
		if !st.cell.edge || st.style.is("curved") || len(st.absPoints) < 2 {
			continue
		}
		if st.style.get("jumpStyle", "none") != "none" {
			st.routedPoints = lineJumps(st, valid)
		}
		valid = append(valid, st)
	}
}

// lineJumps returns the route of an edge with its crossings with the
// edges in valid inserted, in order along each segment.
func lineJumps(st *cellState, valid []*cellState) []routedPoint {
	pts := st.edgePoints()
	if len(pts) < 2 {
		return nil
	}
	const thresh = 0.5
	var out []routedPoint
	at := func(pts []point, i int) *point {
		if i < len(pts) {
			return &pts[i]
		}
		return nil
	}
	for i := 0; i < len(pts)-1; i++ {
		p0, p1 := pts[i], pts[i+1]
		// ignores waypoints on straight segments
		pn := at(pts, i+2)
		for i < len(pts)-2 && ptSegDistSq(p0.x, p0.y, pn.x, pn.y, p1.x, p1.y) < 1 {
			p1 = *pn
			i++
			pn = at(pts, i+2)
		}
		out = append(out, routedPoint{x: p0.x, y: p0.y})
		type crossing struct{ distSq, x, y float64 }
		var list []crossing
		for _, st2 := range valid {
			pts2 := st2.edgePoints()
			if len(pts2) < 2 || !intersects(st.bounds(), st2.bounds()) || st2.style.get("noJump", "") == "1" {
				continue
			}
			var pl *point
			for j := 0; j < len(pts2)-1; j++ {
				p2, p3 := pts2[j], pts2[j+1]
				pn := at(pts2, j+2)
				for j < len(pts2)-2 && ptSegDistSq(p2.x, p2.y, pn.x, pn.y, p3.x, p3.y) < 1 {
					p3 = *pn
					j++
					pn = at(pts2, j+2)
				}
				pt, ok := intersection(p0.x, p0.y, p1.x, p1.y, p2.x, p2.y, p3.x, p3.y)
				if ok && (math.Abs(pt.x-p0.x) > thresh || math.Abs(pt.y-p0.y) > thresh) &&
					(math.Abs(pt.x-p1.x) > thresh || math.Abs(pt.y-p1.y) > thresh) &&
					// removes jumps on overlapping incoming segments
					(pl == nil || ptLineDist(p0.x, p0.y, p1.x, p1.y, pl.x, pl.y) > thresh ||
						ptLineDist(p0.x, p0.y, p1.x, p1.y, p2.x, p2.y) > thresh) &&
					// removes jumps on overlapping outgoing segments
					(pn == nil || ptLineDist(p0.x, p0.y, p1.x, p1.y, pn.x, pn.y) > thresh ||
						ptLineDist(p0.x, p0.y, p1.x, p1.y, p3.x, p3.y) > thresh) {
					dx, dy := pt.x-p0.x, pt.y-p0.y
					c := crossing{dx*dx + dy*dy, pt.x, pt.y}
					// ordered by distance from the start of the segment
					idx := len(list)
					for t := range list {
						if list[t].distSq > c.distSq {
							idx = t
							break
						}
					}
					// ignores multiple crossings at joints and overlapping edges
					if (idx == 0 || math.Abs(list[idx-1].x-c.x) > thresh || math.Abs(list[idx-1].y-c.y) > thresh) &&
						(idx == len(list) || math.Abs(list[idx].x-c.x) > thresh || math.Abs(list[idx].y-c.y) > thresh) {
						list = append(list[:idx], append([]crossing{c}, list[idx:]...)...)
					}
				}
				p := p2
				pl = &p
			}
		}
		for _, c := range list {
			out = append(out, routedPoint{x: c.x, y: c.y, jump: true})
		}
	}
	last := pts[len(pts)-1]
	return append(out, routedPoint{x: last.x, y: last.y})
}

// intersects reports whether two rectangles overlap (mxUtils.intersects).
func intersects(a, b rect) bool {
	tw, th, rw, rh := a.w, a.h, b.w, b.h
	if rw <= 0 || rh <= 0 || tw <= 0 || th <= 0 {
		return false
	}
	tx, ty, rx, ry := a.x, a.y, b.x, b.y
	rw += rx
	rh += ry
	tw += tx
	th += ty
	return (rw < rx || rw > tx) && (rh < ry || rh > ty) && (tw < tx || tw > rx) && (th < ty || th > ry)
}

// ptLineDist is the distance of (px, py) from the line through (x1, y1)
// and (x2, y2) (mxUtils.ptLineDist).
func ptLineDist(x1, y1, x2, y2, px, py float64) float64 {
	return math.Abs((y2-y1)*px-(x2-x1)*py+x2*y1-y2*x1) / math.Sqrt((y2-y1)*(y2-y1)+(x2-x1)*(x2-x1))
}

// paintLineJumps strokes an edge's line through its routed points with a
// jump at every crossing (draw.io's mxConnector.paintLine); absPts are the
// painted points, whose ends the markers have moved. It reports false when
// the edge has no routed points.
func paintLineJumps(s *shape, c *c2d, absPts []point, rounded bool) bool {
	routed := s.st.routedPoints
	if len(routed) == 0 {
		return false
	}
	arcSize := s.style.num("arcSize", lineArcSize) / 2
	size := float64(parseInt(s.style.get("jumpSize", ""), defaultJumpSize)-2)/2 + s.strokewidth
	style := s.style.get("jumpStyle", "none")
	moveTo := true
	var last *point
	var n *point
	var length float64
	var pts []point
	c.begin()
	for i, rpt := range routed {
		pt := point{rpt.x, rpt.y}
		switch i {
		case 0:
			pt = absPts[0]
		case len(routed) - 1:
			pt = absPts[len(absPts)-1]
		}
		done := false
		if last != nil && rpt.jump {
			// checks if the next and previous points are too close
			next := routed[i+1]
			dx, dy := next.x-pt.x, next.y-pt.y
			dist := dx*dx + dy*dy
			if n == nil {
				nn := point{pt.x - last.x, pt.y - last.y}
				length = math.Sqrt(nn.x*nn.x + nn.y*nn.y)
				if length > 0 {
					nn.x = nn.x * size / length
					nn.y = nn.y * size / length
					n = &nn
				}
			}
			if n != nil && dist > size*size && length > 0 {
				dx, dy = last.x-pt.x, last.y-pt.y
				if dx*dx+dy*dy > size*size {
					p0 := point{pt.x - n.x, pt.y - n.y}
					p1 := point{pt.x + n.x, pt.y + n.y}
					pts = append(pts, p0)
					s.addPoints(c, pts, rounded, arcSize, false, nil, moveTo)
					f := -1.0
					if math.Round(n.x) < 0 || math.Round(n.x) == 0 && math.Round(n.y) <= 0 {
						f = 1
					}
					moveTo = false
					switch style {
					case "sharp":
						c.lineTo(p0.x-n.y*f, p0.y+n.x*f)
						c.lineTo(p1.x-n.y*f, p1.y+n.x*f)
						c.lineTo(p1.x, p1.y)
					case "line":
						c.moveTo(p0.x+n.y*f, p0.y-n.x*f)
						c.lineTo(p0.x-n.y*f, p0.y+n.x*f)
						c.moveTo(p1.x-n.y*f, p1.y+n.x*f)
						c.lineTo(p1.x+n.y*f, p1.y-n.x*f)
						c.moveTo(p1.x, p1.y)
					case "arc":
						f *= 1.3
						c.curveTo(p0.x-n.y*f, p0.y+n.x*f, p1.x-n.y*f, p1.y+n.x*f, p1.x, p1.y)
					default:
						c.moveTo(p1.x, p1.y)
						moveTo = true
					}
					pts = []point{p1}
					done = true
				}
			}
		} else {
			n = nil
		}
		if !done {
			pts = append(pts, pt)
			p := pt
			last = &p
		}
	}
	s.addPoints(c, pts, rounded, arcSize, false, nil, moveTo)
	c.stroke()
	return true
}
