package drawio

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

// Edge routing: the points of an edge come from its terminals, control
// points and edge style, as mxGraphView computes them in updateEdgeState:
//
//  1. updateFixedTerminalPoints sets the ends that are fixed by a
//     connection constraint (exitX/exitY, entryX/entryY) or, for a
//     dangling end, by the geometry's source or target point;
//  2. updatePoints runs the edge style (edgestyle.go) or takes the
//     control points as they are;
//  3. updateFloatingTerminalPoints puts the remaining ends on the
//     terminals' perimeters (perimeter.go), heading for the next point.
//
// The functions follow mxGraphView and mxGraph with draw.io's overrides
// from Graph.js; scale is 1 and translate 0 throughout.

// connectionConstraint is an mxConnectionConstraint: a point relative to
// the terminal's bounds (nil for a floating end), whether it is projected
// onto the perimeter, and an absolute offset.
type connectionConstraint struct {
	point     *point
	perimeter bool
	dx, dy    float64
}

// updateFixedTerminalPoints sets the ends of an edge that do not depend on
// the route (mxGraphView.updateFixedTerminalPoints).
func (v *graphView) updateFixedTerminalPoints(edge, source, target *cellState) {
	v.updateFixedTerminalPoint(edge, source, true, v.getConnectionConstraint(edge, source, true))
	v.updateFixedTerminalPoint(edge, target, false, v.getConnectionConstraint(edge, target, false))
}

// updateFixedTerminalPoint (mxGraphView.updateFixedTerminalPoint).
func (v *graphView) updateFixedTerminalPoint(edge, terminal *cellState, source bool, c *connectionConstraint) {
	edge.setAbsoluteTerminalPoint(v.getFixedTerminalPoint(edge, terminal, source, c), source)
}

// getFixedTerminalPoint returns the connection point of a constrained end,
// the geometry's terminal point of a dangling end, or nil for a floating
// end (mxGraphView.getFixedTerminalPoint). Terminals with centerPerimeter
// always connect at their center (Graph.js override).
func (v *graphView) getFixedTerminalPoint(edge, terminal *cellState, source bool, c *connectionConstraint) *point {
	if terminal != nil && terminal.style["perimeter"] == "centerPerimeter" {
		return &point{terminal.cx(), terminal.cy()}
	}
	var pt *point
	if c != nil {
		pt = v.getConnectionPoint(terminal, c, false)
	}
	if pt == nil && terminal == nil {
		geo := edge.cell.geo
		tp := geo.targetPoint
		if source {
			tp = geo.sourcePoint
		}
		if tp != nil {
			pt = &point{tp.x + edge.origin.x, tp.y + edge.origin.y}
		}
	}
	return pt
}

// updateBoundsFromStencil fits the bounds of a terminal to the aspect of
// its fixed-aspect stencil for routing and returns the previous bounds,
// or nil when nothing changed (mxGraphView.updateBoundsFromStencil).
func (v *graphView) updateBoundsFromStencil(st *cellState) *rect {
	if st == nil {
		return nil
	}
	sten := v.stencilFor(st)
	if sten == nil || sten.aspect != "fixed" {
		return nil
	}
	prev := st.bounds()
	st.setRect(fitStencilBounds(prev, sten, st.style["direction"]))
	return &prev
}

// fitStencilBounds centers the largest rectangle of the stencil's aspect in
// r (the stencil turned for north and south directions). Stencils without
// a size leave r as it is.
func fitStencilBounds(r rect, sten *stencil, direction string) rect {
	sw, sh := sten.w0, sten.h0
	if direction == "north" || direction == "south" {
		sw, sh = sh, sw
	}
	if sw <= 0 || sh <= 0 {
		return r
	}
	s := math.Min(r.w/sw, r.h/sh)
	return rect{r.x + (r.w-sw*s)/2, r.y + (r.h-sh*s)/2, sw * s, sh * s}
}

// stencilFor returns the stencil a state is drawn with (state.shape.stencil):
// only shape names that are not built-in shapes resolve to stencils
// (mxCellRenderer.createShape).
func (v *graphView) stencilFor(st *cellState) *stencil {
	if st == nil || st.cell == nil {
		return nil
	}
	name := st.style["shape"]
	if name == "" {
		return nil
	}
	if _, ok := shapeRegistry[name]; ok {
		return nil
	}
	return v.stencilOf(name)
}

// updatePoints replaces the points of an edge between its ends with the
// route of its edge style, or with its control points when it has none
// (mxGraphView.updatePoints).
func (v *graphView) updatePoints(edge *cellState, points []point, source, target *cellState) {
	if len(edge.absPoints) == 0 {
		return
	}
	pts := []*point{edge.absPoints[0]}
	if es := v.getEdgeStyle(edge, points, source, target); es != nil {
		src := v.getTerminalPort(edge, source, true)
		trg := v.getTerminalPort(edge, target, false)
		// uses the stencil bounds for routing and restores them after
		srcBounds := v.updateBoundsFromStencil(src)
		var trgBounds *rect
		if src != trg {
			trgBounds = v.updateBoundsFromStencil(trg)
		}
		pts = es.route(edge, src, trg, points, pts)
		if srcBounds != nil {
			src.setRect(*srcBounds)
		}
		if trgBounds != nil {
			trg.setRect(*trgBounds)
		}
	} else {
		for _, p := range points {
			q := v.transformControlPoint(edge, p)
			pts = append(pts, &q)
		}
	}
	edge.absPoints = append(pts, edge.absPoints[len(edge.absPoints)-1])
}

// isLoopStyleEnabled reports whether an edge is routed with its loop style:
// a loop with fewer than two control points, unless orthogonalLoop is set
// and one of its ends is constrained (mxGraphView.isLoopStyleEnabled).
func (v *graphView) isLoopStyleEnabled(edge *cellState, points []point, source, target *cellState) bool {
	sc := v.getConnectionConstraint(edge, source, true)
	tc := v.getConnectionConstraint(edge, target, false)
	if len(points) < 2 && (!styleTruthy(edge.style, "orthogonalLoop", false) ||
		((sc == nil || sc.point == nil) && (tc == nil || tc.point == nil))) {
		return source != nil && source == target
	}
	return false
}

// getEdgeStyle returns the edge style an edge is routed with, nil for
// straight segments through its control points (mxGraphView.getEdgeStyle):
// the loop style (default Loop) for loops, else the edgeStyle, unless
// noEdgeStyle is set. Unknown names resolve to nil.
func (v *graphView) getEdgeStyle(edge *cellState, points []point, source, target *cellState) *edgeStyle {
	if v.isLoopStyleEnabled(edge, points, source, target) {
		name, ok := edge.style["loopStyle"]
		if !ok {
			return loopStyle // mxGraph.defaultLoopStyle
		}
		return edgeStyles[name]
	}
	if styleTruthy(edge.style, "noEdgeStyle", false) {
		return nil
	}
	return edgeStyles[edge.style["edgeStyle"]]
}

// isOrthogonal reports whether the floating ends of an edge are projected
// orthogonally onto the perimeters (mxGraph.isOrthogonal): the orthogonal
// style when set, else whether the edge style routes orthogonally.
func (v *graphView) isOrthogonal(edge *cellState) bool {
	if _, ok := edge.style["orthogonal"]; ok {
		return styleTruthy(edge.style, "orthogonal", false)
	}
	es := v.getEdgeStyle(edge, nil, nil, nil)
	return es != nil && es.orthogonal
}

// updateFloatingTerminalPoints puts the ends that were not fixed on the
// terminals' perimeters (mxGraphView.updateFloatingTerminalPoints).
func (v *graphView) updateFloatingTerminalPoints(st *cellState, source, target *cellState) {
	pts := st.absPoints
	if len(pts) == 0 {
		return
	}
	p0, pe := pts[0], pts[len(pts)-1]
	if pe == nil && target != nil {
		v.updateFloatingTerminalPoint(st, target, source, false)
	}
	if p0 == nil && source != nil {
		v.updateFloatingTerminalPoint(st, source, target, true)
	}
}

// updateFloatingTerminalPoint sets a floating end of an edge
// (mxGraphView.updateFloatingTerminalPoint), snapped to the nearest
// connection point with snapToPoint (Graph.js override).
func (v *graphView) updateFloatingTerminalPoint(edge, start, end *cellState, source bool) {
	pt := v.getFloatingTerminalPoint(edge, start, end, source)
	if start != nil && (styleIsOne(start.style, "snapToPoint", false) || styleIsOne(edge.style, "snapToPoint", false)) {
		pt = v.snapToAnchorPoint(edge, v.getTerminalPort(edge, start, source), end, source, pt)
	}
	edge.setAbsoluteTerminalPoint(pt, source)
}

// getFloatingTerminalPoint returns where an edge meets the perimeter of
// its terminal (or port) on the way to the next point, in the terminal's
// unrotated frame (mxGraphView.getFloatingTerminalPoint).
func (v *graphView) getFloatingTerminalPoint(edge, start, end *cellState, source bool) *point {
	start = v.getTerminalPort(edge, start, source)
	next := v.getNextPoint(edge, end, source)
	orth := v.isOrthogonal(edge)
	alpha := toRadians(start.style.num("rotation", 0))
	center := point{start.cx(), start.cy()}
	if alpha != 0 && next != nil {
		q := rotatePoint(*next, math.Cos(-alpha), math.Sin(-alpha), center)
		next = &q
	}
	border := edge.style.num("perimeterSpacing", 0)
	if source {
		border += edge.style.num("sourcePerimeterSpacing", 0)
	} else {
		border += edge.style.num("targetPerimeterSpacing", 0)
	}
	pt := v.getPerimeterPoint(start, next, alpha == 0 && orth, border)
	if alpha != 0 && pt != nil {
		q := rotatePoint(*pt, math.Cos(alpha), math.Sin(alpha), center)
		pt = &q
	}
	return pt
}

// snapToAnchorPoint moves a floating end to the nearest connection point
// of its terminal (mxGraphView.snapToAnchorPoint in Graph.js). Only the
// connection points of the points style are known here; the built-in
// points of shapes and stencils are not, and leave pt as it is.
func (v *graphView) snapToAnchorPoint(edge, start, end *cellState, source bool, pt *point) *point {
	if start == nil || edge == nil || pt == nil {
		return pt
	}
	var nearest *point
	dist := 0.0
	for _, c := range v.getAllConnectionConstraints(start) {
		cp := v.getConnectionPoint(start, c, true)
		if cp == nil {
			continue
		}
		tmp := (cp.x-pt.x)*(cp.x-pt.x) + (cp.y-pt.y)*(cp.y-pt.y)
		if nearest == nil || tmp < dist {
			nearest, dist = cp, tmp
		}
	}
	if nearest != nil {
		return nearest
	}
	return pt
}

// getAllConnectionConstraints returns the connection points of a terminal
// given by its points style, [[x, y, perimeter, dx, dy], …]
// (Graph.getAllConnectionConstraints).
func (v *graphView) getAllConnectionConstraints(terminal *cellState) []*connectionConstraint {
	s, ok := terminal.style["points"]
	if !ok {
		return nil
	}
	var raw [][]any
	if json.Unmarshal([]byte(s), &raw) != nil {
		return nil
	}
	var out []*connectionConstraint
	for _, tmp := range raw {
		n := func(i int) float64 {
			if i < len(tmp) {
				if f, ok := jsNumberValue(tmp[i]); ok {
					return f
				}
			}
			return 0
		}
		if len(tmp) < 2 {
			continue
		}
		c := &connectionConstraint{point: &point{n(0), n(1)}, perimeter: true, dx: n(3), dy: n(4)}
		if len(tmp) > 2 {
			// tmp[2] != '0'
			f, ok := jsNumberValue(tmp[2])
			c.perimeter = !ok || f != 0
		}
		out = append(out, c)
	}
	return out
}

// jsNumberValue converts a decoded JSON value to a number as JavaScript's
// loose comparisons do.
func jsNumberValue(x any) (float64, bool) {
	switch x := x.(type) {
	case float64:
		return x, true
	case string:
		return jsNumber(x)
	case bool:
		if x {
			return 1, true
		}
		return 0, true
	}
	return 0, false
}

// getTerminalPort returns the port an edge connects to (the cell named by
// sourcePort or targetPort) or the terminal itself
// (mxGraphView.getTerminalPort).
func (v *graphView) getTerminalPort(st, terminal *cellState, source bool) *cellState {
	key := "targetPort"
	if source {
		key = "sourcePort"
	}
	if id, ok := st.style[key]; ok {
		// only uses ports where a cell state exists
		if tmp := v.validate(v.state(v.m.cells[id])); tmp != nil {
			terminal = tmp
		}
	}
	return terminal
}

// getPerimeterPoint returns where the line from next to the center of a
// terminal crosses its perimeter, or with orthogonal the orthogonal
// projection of next onto it; border is added around the bounds. Flipped
// terminals mirror next and the result (mxGraphView.getPerimeterPoint).
func (v *graphView) getPerimeterPoint(terminal *cellState, next *point, orthogonal bool, border float64) *point {
	if terminal == nil {
		return nil
	}
	var pt *point
	if perimeter := v.getPerimeterFunction(terminal); perimeter != nil && next != nil {
		bounds := v.getPerimeterBounds(terminal, border)
		if bounds.w > 0 || bounds.h > 0 {
			p := *next
			flipH, flipV := false, false
			if terminal.cell != nil && terminal.cell.vertex {
				flipH = styleIsOne(terminal.style, "flipH", false)
				flipV = styleIsOne(terminal.style, "flipV", false)
				// legacy support for stencilFlipH/V
				if v.stencilFor(terminal) != nil {
					flipH = styleIsOne(terminal.style, "stencilFlipH", false) || flipH
					flipV = styleIsOne(terminal.style, "stencilFlipV", false) || flipV
				}
				if flipH {
					p.x = 2*bounds.cx() - p.x
				}
				if flipV {
					p.y = 2*bounds.cy() - p.y
				}
			}
			pt = perimeter(bounds, terminal, p, orthogonal)
			if pt != nil {
				if flipH {
					pt.x = 2*bounds.cx() - pt.x
				}
				if flipV {
					pt.y = 2*bounds.cy() - pt.y
				}
			}
		}
	}
	if pt == nil {
		q := v.getPoint(terminal, nil)
		pt = &q
	}
	return pt
}

// getRoutingCenterX is the x the edge styles route to (the center moved by
// routingCenterX) (mxGraphView.getRoutingCenterX).
func (v *graphView) getRoutingCenterX(st *cellState) float64 {
	return st.cx() + st.style.num("routingCenterX", 0)*st.w
}

// getRoutingCenterY (mxGraphView.getRoutingCenterY).
func (v *graphView) getRoutingCenterY(st *cellState) float64 {
	return st.cy() + st.style.num("routingCenterY", 0)*st.h
}

// getPerimeterBounds returns the bounds of a terminal grown by border and
// its perimeterSpacing (mxGraphView.getPerimeterBounds).
func (v *graphView) getPerimeterBounds(terminal *cellState, border float64) rect {
	border += terminal.style.num("perimeterSpacing", 0)
	return v.cellPerimeterBounds(terminal, border)
}

// cellPerimeterBounds returns the bounds of a state, fitted to the aspect
// of a fixed-aspect stencil and grown by border
// (mxCellState.getPerimeterBounds).
func (v *graphView) cellPerimeterBounds(st *cellState, border float64) rect {
	if sten := v.stencilFor(st); sten != nil && sten.aspect == "fixed" {
		r := fitStencilBounds(st.bounds(), sten, st.style["direction"])
		if border != 0 {
			r = r.grow(border)
		}
		return r
	}
	return st.perimeterBounds(border)
}

// getPerimeterFunction returns the perimeter of a state's perimeter style,
// nil when it has none (mxGraphView.getPerimeterFunction).
func (v *graphView) getPerimeterFunction(st *cellState) perimeterFunc {
	return perimeters[st.style["perimeter"]]
}

// getNextPoint returns the point an edge heads to from an end: the nearest
// other point of the route, or the center of the opposite terminal
// (mxGraphView.getNextPoint).
func (v *graphView) getNextPoint(edge, opposite *cellState, source bool) *point {
	pts := edge.absPoints
	var pt *point
	if n := len(pts); n >= 2 {
		if source {
			pt = pts[min(1, n-1)]
		} else {
			pt = pts[max(0, n-2)]
		}
	}
	if pt == nil && opposite != nil {
		pt = &point{opposite.cx(), opposite.cy()}
	}
	return pt
}

// getConnectionConstraint returns the constraint of an end of an edge:
// exitX/exitY (entryX/entryY) with exitPerimeter and exitDx/exitDy
// (mxGraph.getConnectionConstraint). The point is nil for floating ends.
func (v *graphView) getConnectionConstraint(edge, terminal *cellState, source bool) *connectionConstraint {
	kx, ky, kp, kdx, kdy := "entryX", "entryY", "entryPerimeter", "entryDx", "entryDy"
	if source {
		kx, ky, kp, kdx, kdy = "exitX", "exitY", "exitPerimeter", "exitDx", "exitDy"
	}
	c := &connectionConstraint{}
	if xs, ok := edge.style[kx]; ok {
		if ys, ok := edge.style[ky]; ok {
			// parseFloat: values that are not numbers would make NaN
			// points; they leave the end floating instead
			x, okx := parseFloat(xs)
			y, oky := parseFloat(ys)
			if okx && oky {
				c.point = &point{x, y}
			}
		}
	}
	if c.point != nil {
		c.perimeter = styleTruthy(edge.style, kp, true)
		c.dx = edge.style.num(kdx, 0)
		c.dy = edge.style.num(kdy, 0)
	}
	return c
}

// getConnectionPoint returns the absolute point of a connection constraint
// on a terminal (Graph.getConnectionPoint): draw.io uses the legacy order
// of direction, flipping and rotation unless the terminal's style has
// legacyAnchorPoints=0.
func (v *graphView) getConnectionPoint(vertex *cellState, c *connectionConstraint, round bool) *point {
	if vertex != nil && styleIsOne(vertex.style, "legacyAnchorPoints", true) {
		return v.getLegacyConnectionPoint(vertex, c, round)
	}
	return v.mxGetConnectionPoint(vertex, c, round)
}

// getLegacyConnectionPoint (Graph.getLegacyConnectionPoint): the point is
// turned for the direction and projected onto the perimeter, or else
// flipped, and then rotated with the shape.
func (v *graphView) getLegacyConnectionPoint(vertex *cellState, c *connectionConstraint, round bool) *point {
	if vertex == nil || c.point == nil {
		return nil
	}
	bounds := v.getPerimeterBounds(vertex, 0)
	cx := bounds.center()
	direction, hasDirection := vertex.style["direction"]
	r1 := 0.0
	if hasDirection && styleIsOne(vertex.style, "anchorPointDirection", true) {
		r1 = directionRotation(direction)
		// bounds need to be rotated by 90 degrees for further computation
		if direction == "north" || direction == "south" {
			bounds = bounds.rotate90()
		}
	}
	p := point{bounds.x + c.point.x*bounds.w + c.dx, bounds.y + c.point.y*bounds.h + c.dy}
	pt := &p
	// rotation for direction before projection on perimeter
	r2 := vertex.style.num("rotation", 0)
	if c.perimeter {
		if r1 != 0 {
			cos, sin := quarterTurn(r1)
			p = rotatePoint(p, cos, sin, cx)
		}
		pt = v.getPerimeterPoint(vertex, &p, false, 0)
	} else {
		r2 += r1
		if vertex.cell != nil && vertex.cell.vertex {
			flipH, flipV := v.connectionFlips(vertex, direction)
			if flipH {
				p.x = 2*bounds.cx() - p.x
			}
			if flipV {
				p.y = 2*bounds.cy() - p.y
			}
		}
	}
	// generic rotation after projection on perimeter
	if r2 != 0 && pt != nil {
		rad := toRadians(r2)
		q := rotatePoint(*pt, math.Cos(rad), math.Sin(rad), cx)
		pt = &q
	}
	if round && pt != nil {
		pt.x, pt.y = jsRound(pt.x), jsRound(pt.y)
	}
	return pt
}

// mxGetConnectionPoint (mxGraph.getConnectionPoint): the point is flipped,
// turned for the direction, projected onto the perimeter and then rotated
// with the shape.
func (v *graphView) mxGetConnectionPoint(vertex *cellState, c *connectionConstraint, round bool) *point {
	if vertex == nil || c.point == nil {
		return nil
	}
	bounds := v.getPerimeterBounds(vertex, 0)
	direction, hasDirection := vertex.style["direction"]
	cx := bounds.center()
	// bounds need to be rotated by 90 degrees for further computation
	if direction == "north" || direction == "south" {
		bounds = bounds.rotate90()
	}
	p := point{bounds.x + c.point.x*bounds.w + c.dx, bounds.y + c.point.y*bounds.h + c.dy}
	pt := &p
	// horizontal and vertical shape flipping
	if vertex.cell != nil && vertex.cell.vertex {
		flipH, flipV := v.connectionFlips(vertex, direction)
		if flipH {
			p.x = 2*bounds.cx() - p.x
		}
		if flipV {
			p.y = 2*bounds.cy() - p.y
		}
	}
	// shape direction
	r1 := 0.0
	if hasDirection && styleIsOne(vertex.style, "anchorPointDirection", true) {
		r1 = directionRotation(direction)
	}
	if r1 != 0 {
		cos, sin := quarterTurn(r1)
		p = rotatePoint(p, cos, sin, cx)
	}
	// projection on perimeter
	if c.perimeter {
		pt = v.getPerimeterPoint(vertex, &p, false, 0)
	}
	// general shape rotation
	if r2 := vertex.style.num("rotation", 0); r2 != 0 && pt != nil {
		rad := toRadians(r2)
		q := rotatePoint(*pt, math.Cos(rad), math.Sin(rad), cx)
		pt = &q
	}
	if round && pt != nil {
		pt.x, pt.y = jsRound(pt.x), jsRound(pt.y)
	}
	return pt
}

// connectionFlips returns the flips applied to connection points: flipH
// and flipV (or the legacy stencilFlipH/V of stencils), swapped for north
// and south directions.
func (v *graphView) connectionFlips(vertex *cellState, direction string) (flipH, flipV bool) {
	flipH = styleIsOne(vertex.style, "flipH", false)
	flipV = styleIsOne(vertex.style, "flipV", false)
	if v.stencilFor(vertex) != nil {
		flipH = styleIsOne(vertex.style, "stencilFlipH", false) || flipH
		flipV = styleIsOne(vertex.style, "stencilFlipV", false) || flipV
	}
	if direction == "north" || direction == "south" {
		flipH, flipV = flipV, flipH
	}
	return flipH, flipV
}

// directionRotation is the turn of a direction in degrees.
func directionRotation(direction string) float64 {
	switch direction {
	case "north":
		return 270
	case "west":
		return 180
	case "south":
		return 90
	}
	return 0
}

// quarterTurn returns the cosine and sine of 90, 180 or 270 degrees
// without trigonometry, as mxGraph does.
func quarterTurn(deg float64) (cos, sin float64) {
	switch deg {
	case 90:
		return 0, 1
	case 180:
		return -1, 0
	case 270:
		return 0, -1
	}
	return 0, 0
}

// gridSize is the grid size of the model (mxGraph.gridSize), the default
// segment of loops.
func (v *graphView) gridSize() float64 {
	if f, ok := parseFloat(v.m.attrs["gridSize"]); ok && f != 0 {
		return f
	}
	return 10
}

// Style values as JavaScript sees them: mxStylesheet stores numeric values
// as numbers and everything else as strings, and mxGraph tests them for
// truthiness or compares them with == 1.

// jsNumber converts s like JavaScript's Number when s is numeric in the
// sense of mxUtils.isNumeric.
func jsNumber(s string) (float64, bool) {
	t := strings.TrimSpace(s)
	if t == "" || strings.Contains(strings.ToLower(t), "0x") {
		return 0, false
	}
	f, err := strconv.ParseFloat(t, 64)
	if err != nil || !finite(f) {
		return 0, false
	}
	return f, true
}

// styleTruthy reports whether a style value is truthy in JavaScript
// (mxUtils.getValue(style, key, def) used as a condition).
func styleTruthy(s style, key string, def bool) bool {
	v, ok := s[key]
	if !ok {
		return def
	}
	if f, ok := jsNumber(v); ok {
		return f != 0
	}
	return v != ""
}

// styleIsOne reports whether a style value == 1 in JavaScript
// (mxUtils.getValue(style, key, def) == 1).
func styleIsOne(s style, key string, def bool) bool {
	v, ok := s[key]
	if !ok {
		return def
	}
	f, ok := jsNumber(v)
	return ok && f == 1
}

// jsRound rounds halves up like JavaScript's Math.round.
func jsRound(x float64) float64 {
	r := math.Floor(x)
	if x-r >= 0.5 {
		r++
	}
	return r
}
