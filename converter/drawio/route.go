package drawio

import "math"

// Edge routing: terminal points, edge styles and perimeters
// (mxGraphView.updateFixedTerminalPoints, updatePoints,
// updateFloatingTerminalPoints and getPerimeterPoint).
//
// This is a minimal version: straight segments through the control
// points, ending on the terminals' bounding rectangles.

func (v *graphView) updateFixedTerminalPoints(st *cellState, src, trg *cellState) {
	v.updateFixedTerminalPoint(st, src, true)
	v.updateFixedTerminalPoint(st, trg, false)
}

func (v *graphView) updateFixedTerminalPoint(st *cellState, term *cellState, source bool) {
	var pt *point
	if term == nil {
		geo := st.cell.geo
		tp := geo.targetPoint
		if source {
			tp = geo.sourcePoint
		}
		if tp != nil {
			p := point{tp.x + st.origin.x, tp.y + st.origin.y}
			pt = &p
		}
	}
	st.setAbsoluteTerminalPoint(pt, source)
}

func (v *graphView) updatePoints(st *cellState, points []point, src, trg *cellState) {
	if len(st.absPoints) == 0 {
		return
	}
	pts := []*point{st.absPoints[0]}
	for _, p := range points {
		q := v.transformControlPoint(st, p)
		pts = append(pts, &q)
	}
	pts = append(pts, st.absPoints[len(st.absPoints)-1])
	st.absPoints = pts
}

func (v *graphView) updateFloatingTerminalPoints(st *cellState, src, trg *cellState) {
	pts := st.absPoints
	if len(pts) == 0 {
		return
	}
	if pts[len(pts)-1] == nil && trg != nil {
		p := v.floatingTerminalPoint(st, trg, src, false)
		st.setAbsoluteTerminalPoint(&p, false)
	}
	if st.absPoints[0] == nil && src != nil {
		p := v.floatingTerminalPoint(st, src, trg, true)
		st.setAbsoluteTerminalPoint(&p, true)
	}
}

func (v *graphView) floatingTerminalPoint(st *cellState, start, end *cellState, source bool) point {
	next := v.nextPoint(st, end, source)
	b := start.bounds()
	c := b.center()
	dx, dy := next.x-c.x, next.y-c.y
	if dx == 0 && dy == 0 || b.w == 0 || b.h == 0 {
		return c
	}
	t := math.Min(b.w/2/math.Abs(dx), b.h/2/math.Abs(dy))
	return point{c.x + dx*t, c.y + dy*t}
}

// nextPoint returns the point an edge heads to from a terminal: the
// nearest control point, or the center of the opposite terminal
// (mxGraphView.getNextPoint).
func (v *graphView) nextPoint(st *cellState, opposite *cellState, source bool) point {
	pts := st.absPoints
	if len(pts) >= 2 {
		i := len(pts) - 2
		if source {
			i = 1
		}
		if p := pts[i]; p != nil {
			return *p
		}
	}
	if opposite != nil {
		return point{opposite.cx(), opposite.cy()}
	}
	return point{}
}
