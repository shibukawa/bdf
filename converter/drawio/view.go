package drawio

import (
	"math"
)

// The view computes where cells are, as mxGraphView does at scale 1 with
// no translation: absolute vertex bounds from nested geometries, and the
// points of edges from their terminals, control points and edge style.

// cellState is the computed state of a visible cell (mxCellState).
type cellState struct {
	cell  *cell
	style style
	view  *graphView

	// x, y, w, h are the absolute bounds: the shape's rectangle for
	// vertices, the bounding box of the points for edges.
	x, y, w, h float64
	// origin is the absolute position of the cell's geometry origin.
	origin point
	// absoluteOffset is the label position: an offset from the vertex
	// label box, or the absolute label anchor of an edge.
	absoluteOffset       point
	unscaledW, unscaledH float64

	// Edges: the points of the route; nil entries are terminal points not
	// yet computed.
	absPoints        []*point
	terminalDistance float64
	length           float64
	segments         []float64
	visibleSource    *cellState
	visibleTarget    *cellState
	// routedPoints is the route with line jumps (jumpStyle), see jumps.go.
	routedPoints []routedPoint

	parent *cellState // state of the parent cell (nil for layers)
	layer  *cell      // the layer the cell belongs to

	validated bool
	visiting  bool
}

func (s *cellState) cx() float64    { return s.x + s.w/2 }
func (s *cellState) cy() float64    { return s.y + s.h/2 }
func (s *cellState) bounds() rect   { return rect{s.x, s.y, s.w, s.h} }
func (s *cellState) setRect(r rect) { s.x, s.y, s.w, s.h = r.x, r.y, r.w, r.h }

// perimeterBounds returns the bounds grown by border (mxCellState.getPerimeterBounds).
func (s *cellState) perimeterBounds(border float64) rect {
	r := s.bounds()
	if border != 0 {
		r = r.grow(border)
	}
	return r
}

// setAbsoluteTerminalPoint sets the first (source) or last point of an edge.
func (s *cellState) setAbsoluteTerminalPoint(p *point, source bool) {
	if source {
		if len(s.absPoints) == 0 {
			s.absPoints = []*point{p}
		} else {
			s.absPoints[0] = p
		}
		return
	}
	if len(s.absPoints) == 0 {
		s.absPoints = []*point{nil, p}
	} else if len(s.absPoints) == 1 {
		s.absPoints = append(s.absPoints, p)
	} else {
		s.absPoints[len(s.absPoints)-1] = p
	}
}

// graphView holds the states of the visible cells of one page.
type graphView struct {
	m      *model
	states map[*cell]*cellState
	order  []*cellState // painting order: depth-first, document order of children
	warn   func(key, format string, args ...any)
	// stencilOf returns the stencil a shape name stands for (nil when it
	// is not a stencil), for routing edges to the bounds of fixed-aspect
	// stencils (mxGraphView.updateBoundsFromStencil).
	stencilOf func(name string) *stencil
}

// newView computes the states of every visible cell of m.
func newView(m *model, warn func(key, format string, args ...any), stencilOf func(string) *stencil) *graphView {
	v := &graphView{m: m, states: map[*cell]*cellState{}, warn: warn, stencilOf: stencilOf}
	if warn == nil {
		v.warn = func(string, string, ...any) {}
	}
	if stencilOf == nil {
		v.stencilOf = func(string) *stencil { return nil }
	}
	if m.root == nil {
		return v
	}
	// create states for visible cells (mxGraphView.validateCell)
	var create func(c *cell, layer *cell, parent *cellState)
	create = func(c *cell, layer *cell, parent *cellState) {
		if !c.visible {
			return
		}
		st := &cellState{cell: c, style: parseStyle(c.styleStr, c.edge), view: v, parent: parent, layer: layer}
		v.states[c] = st
		v.order = append(v.order, st)
		if c.collapsed && c.vertex {
			return
		}
		for _, k := range c.children {
			if layer == nil {
				create(k, c, st) // children of the root are layers
			} else {
				create(k, layer, st)
			}
		}
	}
	for _, l := range m.root.children {
		create(l, nil, nil)
	}
	for _, st := range v.order {
		v.validate(st)
	}
	v.updateLineJumps()
	return v
}

func (v *graphView) state(c *cell) *cellState {
	if c == nil {
		return nil
	}
	return v.states[c]
}

// validate computes a state after its parent and terminals
// (mxGraphView.validateCellState).
func (v *graphView) validate(st *cellState) *cellState {
	if st == nil || st.validated {
		return st
	}
	if st.visiting {
		return st // a cycle through terminals: use what is known
	}
	st.visiting = true
	defer func() { st.visiting = false }()
	if st.parent != nil {
		v.validate(st.parent)
	}
	if st.cell.edge {
		st.visibleSource = v.validate(v.state(v.visibleTerminal(st.cell, true)))
		st.visibleTarget = v.validate(v.state(v.visibleTerminal(st.cell, false)))
	}
	v.updateCellState(st)
	st.validated = true
	return st
}

// visibleTerminal returns the terminal of an edge, or its nearest
// collapsed ancestor that stands for it (mxGraphView.getVisibleTerminal).
func (v *graphView) visibleTerminal(edge *cell, source bool) *cell {
	result := edge.target
	if source {
		result = edge.source
	}
	best := result
	for result != nil && result != v.m.root {
		if v.states[best] == nil || result.collapsed {
			best = result
		}
		result = result.parent
	}
	if best != nil && (best.parent == v.m.root || best == v.m.root) {
		best = nil
	}
	return best
}

// updateCellState computes the bounds of a cell (mxGraphView.updateCellState).
func (v *graphView) updateCellState(st *cellState) {
	st.absoluteOffset = point{}
	st.origin = point{}
	st.length = 0
	p := st.parent
	if p != nil && p.cell != v.m.root {
		st.origin = p.origin
	}
	geo := st.cell.geo
	if geo == nil {
		return
	}
	if !st.cell.edge {
		off := point{}
		if geo.offset != nil {
			off = *geo.offset
		}
		if geo.relative && p != nil {
			if p.cell.edge {
				o := v.getPoint(p, geo)
				st.origin.x += o.x - p.origin.x
				st.origin.y += o.y - p.origin.y
			} else {
				st.origin.x += geo.x*p.unscaledW + off.x
				st.origin.y += geo.y*p.unscaledH + off.y
			}
		} else {
			st.absoluteOffset = off
			st.origin.x += geo.x
			st.origin.y += geo.y
		}
	}
	st.x, st.y = st.origin.x, st.origin.y
	st.w, st.h = geo.w, geo.h
	st.unscaledW, st.unscaledH = geo.w, geo.h
	if st.cell.vertex {
		v.updateVertexState(st, geo)
	}
	if st.cell.edge {
		v.updateEdgeState(st, geo)
	}
}

// updateVertexState places relative children of rotated parents and the
// vertex label (mxGraphView.updateVertexState).
func (v *graphView) updateVertexState(st *cellState, geo *geometry) {
	p := st.parent
	if geo.relative && p != nil && !p.cell.edge {
		if alpha := toRadians(p.style.num("rotation", 0)); alpha != 0 {
			cos, sin := math.Cos(alpha), math.Sin(alpha)
			pt := rotatePoint(point{st.cx(), st.cy()}, cos, sin, point{p.cx(), p.cy()})
			st.x, st.y = pt.x-st.w/2, pt.y-st.h/2
		}
	}
	v.updateVertexLabelOffset(st)
}

// updateVertexLabelOffset moves the label box for labelPosition and
// verticalLabelPosition (mxGraphView.updateVertexLabelOffset).
func (v *graphView) updateVertexLabelOffset(st *cellState) {
	switch st.style.get("labelPosition", "center") {
	case "left":
		lw := st.style.num("labelWidth", st.w)
		st.absoluteOffset.x -= lw
	case "right":
		st.absoluteOffset.x += st.w
	case "center":
		if st.style.has("labelWidth") {
			lw := st.style.num("labelWidth", 0)
			dx := 0.0
			switch st.style.get("align", "center") {
			case "center":
				dx = 0.5
			case "right":
				dx = 1
			}
			if dx != 0 {
				st.absoluteOffset.x -= (lw - st.w) * dx
			}
		}
	}
	switch st.style.get("verticalLabelPosition", "middle") {
	case "top":
		st.absoluteOffset.y -= st.h
	case "bottom":
		st.absoluteOffset.y += st.h
	}
}

// updateEdgeState computes the route of an edge (mxGraphView.updateEdgeState).
func (v *graphView) updateEdgeState(st *cellState, geo *geometry) {
	src, trg := st.visibleSource, st.visibleTarget
	if (st.cell.source != nil && src == nil) || (src == nil && geo.sourcePoint == nil) ||
		(st.cell.target != nil && trg == nil) || (trg == nil && geo.targetPoint == nil) {
		st.absPoints = nil
		return
	}
	st.absPoints = nil
	v.updateFixedTerminalPoints(st, src, trg)
	v.updatePoints(st, geo.points, src, trg)
	v.updateFloatingTerminalPoints(st, src, trg)
	pts := st.absPoints
	if len(pts) < 2 || pts[0] == nil || pts[len(pts)-1] == nil {
		st.absPoints = nil
		return
	}
	v.updateEdgeBounds(st)
	v.updateEdgeLabelOffset(st)
}

// edgePoints returns the resolved route of an edge (nil when it is not drawn).
func (st *cellState) edgePoints() []point {
	if len(st.absPoints) < 2 {
		return nil
	}
	out := make([]point, 0, len(st.absPoints))
	for _, p := range st.absPoints {
		if p != nil {
			out = append(out, *p)
		}
	}
	return out
}

// transformControlPoint makes a control point of an edge absolute.
func (v *graphView) transformControlPoint(st *cellState, p point) point {
	return point{p.x + st.origin.x, p.y + st.origin.y}
}

// updateEdgeBounds sets the length, segments and bounding box of an edge
// (mxGraphView.updateEdgeBounds).
func (v *graphView) updateEdgeBounds(st *cellState) {
	pts := st.absPoints
	p0, pe := pts[0], pts[len(pts)-1]
	if p0.x != pe.x || p0.y != pe.y {
		st.terminalDistance = math.Hypot(pe.x-p0.x, pe.y-p0.y)
	} else {
		st.terminalDistance = 0
	}
	length := 0.0
	st.segments = st.segments[:0]
	pt := p0
	minX, minY, maxX, maxY := pt.x, pt.y, pt.x, pt.y
	for _, tmp := range pts[1:] {
		if tmp == nil {
			continue
		}
		seg := math.Hypot(pt.x-tmp.x, pt.y-tmp.y)
		st.segments = append(st.segments, seg)
		length += seg
		pt = tmp
		minX, minY = math.Min(pt.x, minX), math.Min(pt.y, minY)
		maxX, maxY = math.Max(pt.x, maxX), math.Max(pt.y, maxY)
	}
	st.length = length
	st.x, st.y = minX, minY
	st.w, st.h = math.Max(1, maxX-minX), math.Max(1, maxY-minY)
}

// getPoint returns the point at a relative geometry along an edge: x from
// -1 (source) to 1 (target), y a perpendicular distance
// (mxGraphView.getPoint).
func (v *graphView) getPoint(st *cellState, geo *geometry) point {
	x, y := st.cx(), st.cy()
	if st.segments != nil && len(st.absPoints) > 1 && (geo == nil || geo.relative) {
		gx := 0.0
		if geo != nil {
			gx = geo.x / 2
		}
		pointCount := len(st.absPoints)
		dist := math.Round((gx + 0.5) * st.length)
		segment := 0.0
		if len(st.segments) > 0 {
			segment = st.segments[0]
		}
		length := 0.0
		index := 1
		for dist >= math.Round(length+segment) && index < pointCount-1 && index < len(st.segments) {
			length += segment
			segment = st.segments[index]
			index++
		}
		factor := 0.0
		if segment != 0 {
			factor = (dist - length) / segment
		}
		p0, pe := st.absPoints[index-1], st.absPoints[index]
		if p0 != nil && pe != nil {
			gy, ox, oy := 0.0, 0.0, 0.0
			if geo != nil {
				gy = geo.y
				if geo.offset != nil {
					ox, oy = geo.offset.x, geo.offset.y
				}
			}
			dx, dy := pe.x-p0.x, pe.y-p0.y
			nx, ny := 0.0, 0.0
			if segment != 0 {
				nx, ny = dy/segment, dx/segment
			}
			x = p0.x + dx*factor + (nx*gy + ox)
			y = p0.y + dy*factor - (ny*gy - oy)
		}
	} else if geo != nil && geo.offset != nil {
		x += geo.offset.x
		y += geo.offset.y
	}
	return point{x, y}
}

// updateEdgeLabelOffset sets the anchor of an edge's own label
// (mxGraphView.updateEdgeLabelOffset).
func (v *graphView) updateEdgeLabelOffset(st *cellState) {
	pts := st.absPoints
	st.absoluteOffset = point{st.cx(), st.cy()}
	if len(pts) == 0 || st.segments == nil {
		return
	}
	geo := st.cell.geo
	if geo.relative {
		st.absoluteOffset = v.getPoint(st, geo)
		return
	}
	p0, pe := pts[0], pts[len(pts)-1]
	if p0 != nil && pe != nil {
		x0, y0 := 0.0, 0.0
		if geo.offset != nil {
			x0, y0 = geo.offset.x, geo.offset.y
		}
		st.absoluteOffset = point{p0.x + (pe.x-p0.x)/2 + x0, p0.y + (pe.y-p0.y)/2 + y0}
	}
}
