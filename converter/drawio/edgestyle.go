package drawio

import (
	"math"
	"strings"
)

// Edge styles (mxEdgeStyle and draw.io's additions in Shapes.js): each
// computes the points of an edge between its ends. result holds the
// source point (nil while it floats) and the style appends the points
// that follow; the target point is added by updatePoints. source and
// target are the terminal states (or ports), either may be nil for a
// dangling end, and points are the control points of the geometry.

// edgeStyleFunc is an mxEdgeStyle function; it returns result with the
// route appended.
type edgeStyleFunc func(state, source, target *cellState, points []point, result []*point) []*point

// edgeStyle is a registered edge style.
type edgeStyle struct {
	route edgeStyleFunc
	// orthogonal marks the styles whose ends are projected orthogonally
	// onto the perimeters (mxGraph.isOrthogonalEdgeStyle).
	orthogonal bool
}

// loopStyle is mxGraph.defaultLoopStyle.
var loopStyle = &edgeStyle{route: loopRoute}

// edgeStyles maps edgeStyle and loopStyle values to edge styles
// (mxStyleRegistry).
var edgeStyles = map[string]*edgeStyle{
	"elbowEdgeStyle":          {route: elbowConnector, orthogonal: true},
	"entityRelationEdgeStyle": {route: entityRelation, orthogonal: true},
	"loopEdgeStyle":           loopStyle,
	"sideToSideEdgeStyle":     {route: sideToSide, orthogonal: true},
	"topToBottomEdgeStyle":    {route: topToBottom, orthogonal: true},
	"orthogonalEdgeStyle":     {route: orthConnector, orthogonal: true},
	"segmentEdgeStyle":        {route: segmentConnector, orthogonal: true},
	"isometricEdgeStyle":      {route: isometricConnector},
	"sequenceEdgeStyle":       {route: sequenceMessage},
}

// Direction masks (mxConstants.DIRECTION_MASK_*).
const (
	dirMaskNone  = 0
	dirMaskWest  = 1
	dirMaskNorth = 2
	dirMaskSouth = 4
	dirMaskEast  = 8
	dirMaskAll   = 15
)

// entitySegment is the default segment of entity relations
// (mxConstants.ENTITY_SEGMENT).
const entitySegment = 30

// pointState returns a zero-size state at p, for routing to a fixed end
// (new mxCellState() with x and y set).
func pointState(p *point) *cellState { return &cellState{x: p.x, y: p.y} }

// entityRelation routes an edge out of the left or right side of both
// terminals, with a horizontal segment at each end (mxEdgeStyle.EntityRelation).
func entityRelation(state, source, target *cellState, points []point, result []*point) []*point {
	view := state.view
	segment := state.style.num("segment", entitySegment)
	pts := state.absPoints
	p0, pe := pts[0], pts[len(pts)-1]

	isSourceLeft := false
	if source != nil {
		if sg := cellGeometry(source); sg != nil && sg.relative {
			isSourceLeft = sg.x <= 0.5
		} else if target != nil {
			a := target.x + target.w
			if pe != nil {
				a = pe.x
			}
			b := source.x
			if p0 != nil {
				b = p0.x
			}
			isSourceLeft = a < b
		}
	}
	if p0 != nil {
		source = pointState(p0)
	} else if source != nil {
		c := getPortConstraints(source, state, true, dirMaskNone)
		if c != dirMaskNone && c != dirMaskWest+dirMaskEast {
			isSourceLeft = c == dirMaskWest
		}
	} else {
		return result
	}

	isTargetLeft := true
	if target != nil {
		if tg := cellGeometry(target); tg != nil && tg.relative {
			isTargetLeft = tg.x <= 0.5
		} else if source != nil {
			a := source.x + source.w
			if p0 != nil {
				a = p0.x
			}
			b := target.x
			if pe != nil {
				b = pe.x
			}
			isTargetLeft = a < b
		}
	}
	if pe != nil {
		target = pointState(pe)
	} else if target != nil {
		c := getPortConstraints(target, state, false, dirMaskNone)
		if c != dirMaskNone && c != dirMaskWest+dirMaskEast {
			isTargetLeft = c == dirMaskWest
		}
	}

	if source != nil && target != nil {
		x0 := source.x + source.w
		if isSourceLeft {
			x0 = source.x
		}
		y0 := view.getRoutingCenterY(source)
		xe := target.x + target.w
		if isTargetLeft {
			xe = target.x
		}
		ye := view.getRoutingCenterY(target)
		seg := segment
		dx := seg
		if isSourceLeft {
			dx = -seg
		}
		dep := point{x0 + dx, y0}
		dx = seg
		if isTargetLeft {
			dx = -seg
		}
		arr := point{xe + dx, ye}
		// adds intermediate points if both go out on the same side
		if isSourceLeft == isTargetLeft {
			x := math.Max(x0, xe) + segment
			if isSourceLeft {
				x = math.Min(x0, xe) - segment
			}
			result = append(result, &point{x, y0}, &point{x, ye})
		} else if (dep.x < arr.x) == isSourceLeft {
			midY := y0 + (ye-y0)/2
			result = append(result, &dep, &point{dep.x, midY}, &point{arr.x, midY}, &arr)
		} else {
			result = append(result, &dep, &arr)
		}
	}
	return result
}

// cellGeometry returns the geometry of a state's cell (mxGraph.getCellGeometry).
func cellGeometry(st *cellState) *geometry {
	if st.cell == nil {
		return nil
	}
	return st.cell.geo
}

// loopRoute routes a self-reference beside its terminal: to the side the
// direction style names (west by default), segment (the grid size) away,
// or through the first control point (mxEdgeStyle.Loop).
func loopRoute(state, source, target *cellState, points []point, result []*point) []*point {
	view := state.view
	pts := state.absPoints
	p0, pe := pts[0], pts[len(pts)-1]
	if p0 != nil && pe != nil {
		for _, p := range points {
			q := view.transformControlPoint(state, p)
			result = append(result, &q)
		}
		return result
	}
	if source == nil {
		return result
	}
	var pt *point
	if len(points) > 0 {
		q := view.transformControlPoint(state, points[0])
		if !source.bounds().contains(q.x, q.y) {
			pt = &q
		}
	}
	x, dx, y, dy := 0.0, 0.0, 0.0, 0.0
	seg := state.style.num("segment", view.gridSize())
	dir := state.style.get("direction", "west")
	if dir == "north" || dir == "south" {
		x = view.getRoutingCenterX(source)
		dx = seg
	} else {
		y = view.getRoutingCenterY(source)
		dy = seg
	}
	if pt == nil || pt.x < source.x || pt.x > source.x+source.w {
		if pt != nil {
			x = pt.x
			dy = math.Max(math.Abs(y-pt.y), dy)
		} else {
			switch dir {
			case "north":
				y = source.y - 2*dx
			case "south":
				y = source.y + source.h + 2*dx
			case "east":
				x = source.x - 2*dy
			default:
				x = source.x + source.w + 2*dy
			}
		}
	} else {
		x = view.getRoutingCenterX(source)
		dx = math.Max(math.Abs(x-pt.x), dy)
		y = pt.y
		dy = 0
	}
	return append(result, &point{x - dx, y - dy}, &point{x + dx, y + dy})
}

// elbowConnector uses topToBottom or sideToSide depending on how the
// terminals and the first control point lie and on the elbow style
// (mxEdgeStyle.ElbowConnector).
func elbowConnector(state, source, target *cellState, points []point, result []*point) []*point {
	vertical, horizontal := false, false
	if source != nil && target != nil {
		if len(points) > 0 {
			left := math.Min(source.x, target.x)
			right := math.Max(source.x+source.w, target.x+target.w)
			top := math.Min(source.y, target.y)
			bottom := math.Max(source.y+source.h, target.y+target.h)
			pt := state.view.transformControlPoint(state, points[0])
			vertical = pt.y < top || pt.y > bottom
			horizontal = pt.x < left || pt.x > right
		} else {
			left := math.Max(source.x, target.x)
			right := math.Min(source.x+source.w, target.x+target.w)
			vertical = left == right
			if !vertical {
				top := math.Max(source.y, target.y)
				bottom := math.Min(source.y+source.h, target.y+target.h)
				horizontal = top == bottom
			}
		}
	}
	if !horizontal && (vertical || state.style["elbow"] == "vertical") {
		return topToBottom(state, source, target, points, result)
	}
	return sideToSide(state, source, target, points, result)
}

// sideToSide routes a horizontal elbow: out of the sides, with a vertical
// segment in the middle or at the control point (mxEdgeStyle.SideToSide).
func sideToSide(state, source, target *cellState, points []point, result []*point) []*point {
	view := state.view
	var pt *point
	if len(points) > 0 {
		q := view.transformControlPoint(state, points[0])
		pt = &q
	}
	pts := state.absPoints
	p0, pe := pts[0], pts[len(pts)-1]
	if p0 != nil {
		source = pointState(p0)
	}
	if pe != nil {
		target = pointState(pe)
	}
	if source == nil || target == nil {
		return result
	}
	l := math.Max(source.x, target.x)
	r := math.Min(source.x+source.w, target.x+target.w)
	x := jsRound(r + (l-r)/2)
	if pt != nil {
		x = pt.x
	}
	y1 := view.getRoutingCenterY(source)
	y2 := view.getRoutingCenterY(target)
	if pt != nil {
		if pt.y >= source.y && pt.y <= source.y+source.h {
			y1 = pt.y
		}
		if pt.y >= target.y && pt.y <= target.y+target.h {
			y2 = pt.y
		}
	}
	if !target.bounds().contains(x, y1) && !source.bounds().contains(x, y1) {
		result = append(result, &point{x, y1})
	}
	if !target.bounds().contains(x, y2) && !source.bounds().contains(x, y2) {
		result = append(result, &point{x, y2})
	}
	if len(result) == 1 {
		if pt != nil {
			if !target.bounds().contains(x, pt.y) && !source.bounds().contains(x, pt.y) {
				result = append(result, &point{x, pt.y})
			}
		} else {
			t := math.Max(source.y, target.y)
			b := math.Min(source.y+source.h, target.y+target.h)
			result = append(result, &point{x, t + (b-t)/2})
		}
	}
	return result
}

// topToBottom routes a vertical elbow: out of the top or bottom, with a
// horizontal segment in the middle or at the control point
// (mxEdgeStyle.TopToBottom).
func topToBottom(state, source, target *cellState, points []point, result []*point) []*point {
	view := state.view
	var pt *point
	if len(points) > 0 {
		q := view.transformControlPoint(state, points[0])
		pt = &q
	}
	pts := state.absPoints
	p0, pe := pts[0], pts[len(pts)-1]
	if p0 != nil {
		source = pointState(p0)
	}
	if pe != nil {
		target = pointState(pe)
	}
	if source == nil || target == nil {
		return result
	}
	t := math.Max(source.y, target.y)
	b := math.Min(source.y+source.h, target.y+target.h)
	x := view.getRoutingCenterX(source)
	if pt != nil && pt.x >= source.x && pt.x <= source.x+source.w {
		x = pt.x
	}
	y := jsRound(b + (t-b)/2)
	if pt != nil {
		y = pt.y
	}
	if !target.bounds().contains(x, y) && !source.bounds().contains(x, y) {
		result = append(result, &point{x, y})
	}
	if pt != nil && pt.x >= target.x && pt.x <= target.x+target.w {
		x = pt.x
	} else {
		x = view.getRoutingCenterX(target)
	}
	if !target.bounds().contains(x, y) && !source.bounds().contains(x, y) {
		result = append(result, &point{x, y})
	}
	if len(result) == 1 {
		if pt != nil {
			if !target.bounds().contains(pt.x, y) && !source.bounds().contains(pt.x, y) {
				result = append(result, &point{pt.x, y})
			}
		} else {
			l := math.Max(source.x, target.x)
			r := math.Min(source.x+source.w, target.x+target.w)
			result = append(result, &point{l + (r-l)/2, y})
		}
	}
	return result
}

// round1 rounds to a tenth, as the orthogonal styles do at scale 1.
func round1(v float64) float64 { return jsRound(v*10) / 10 }

// scalePointArray copies the points of an edge rounded to a tenth
// (mxEdgeStyle.scalePointArray at scale 1).
func scalePointArray(pts []*point) []*point {
	out := make([]*point, len(pts))
	for i, p := range pts {
		if p != nil {
			out[i] = &point{round1(p.x), round1(p.y)}
		}
	}
	return out
}

// scaleCellState copies a state with its bounds rounded to a tenth
// (mxEdgeStyle.scaleCellState at scale 1).
func scaleCellState(st *cellState) *cellState {
	if st == nil {
		return nil
	}
	c := *st
	c.setRect(rect{round1(st.x), round1(st.y), round1(st.w), round1(st.h)})
	return &c
}

// segmentConnector routes an orthogonal edge through its control points,
// alternating horizontal and vertical segments (mxEdgeStyle.SegmentConnector).
func segmentConnector(state, sourceScaled, targetScaled *cellState, controlHints []point, result []*point) []*point {
	view := state.view
	// all way- and terminal points, rounded
	pts := scalePointArray(state.absPoints)
	source := scaleCellState(sourceScaled)
	target := scaleCellState(targetScaled)
	const tol = 1.0

	// translated unscaled points for precise collision checks
	var tempPoints []point
	addPoint := func(p point) { tempPoints = append(tempPoints, p) }

	// whether the first segment outgoing from the source end is horizontal
	var lastPushed *point
	if len(result) > 0 {
		lastPushed = result[0]
	}
	horizontal := true
	var hint *point

	// adds waypoints only if outside of tolerance
	pushPoint := func(p point) {
		p.x = round1(p.x)
		p.y = round1(p.y)
		if lastPushed == nil || math.Abs(lastPushed.x-p.x) >= tol || math.Abs(lastPushed.y-p.y) >= 1 {
			q := p
			result = append(result, &q)
			lastPushed = &q
		}
	}

	// the first point
	var pt *point
	if pts[0] == nil && source != nil {
		pt = &point{view.getRoutingCenterX(source), view.getRoutingCenterY(source)}
	} else if pts[0] != nil {
		c := *pts[0]
		pt = &c
	}
	lastInx := len(pts) - 1
	// pe is only set with control points (a hoisted var in the JavaScript)
	var pe *point

	if len(controlHints) > 0 {
		hints := make([]*point, 0, len(controlHints))
		for _, h := range controlHints {
			q := view.transformControlPoint(state, h)
			hints = append(hints, &q)
		}
		// aligns source and target hint to fixed points
		if pt != nil {
			if math.Abs(hints[0].x-pt.x) < tol {
				hints[0].x = pt.x
			}
			if math.Abs(hints[0].y-pt.y) < tol {
				hints[0].y = pt.y
			}
		}
		pe = pts[lastInx]
		if pe != nil {
			h := hints[len(hints)-1]
			if math.Abs(h.x-pe.x) < tol {
				h.x = pe.x
			}
			if math.Abs(h.y-pe.y) < tol {
				h.y = pe.y
			}
		}
		hint = hints[0]

		currentTerm := source
		currentPt := pts[0]
		currentHint := hint
		if currentPt != nil {
			currentTerm = nil
		}
		// checks for alignment with fixed points and with channels at
		// source and target segments only
		for i := 0; i < 2; i++ {
			fixedVertAlign := currentPt != nil && currentPt.x == currentHint.x
			fixedHozAlign := currentPt != nil && currentPt.y == currentHint.y
			inHozChan := currentTerm != nil && currentHint.y >= currentTerm.y && currentHint.y <= currentTerm.y+currentTerm.h
			inVertChan := currentTerm != nil && currentHint.x >= currentTerm.x && currentHint.x <= currentTerm.x+currentTerm.w
			hozChan := fixedHozAlign || (currentPt == nil && inHozChan)
			vertChan := fixedVertAlign || (currentPt == nil && inVertChan)
			// if the hint falls in both channels of a floating port, or
			// coincides with a fixed point, ignore the source and work out
			// the orientation from the target end
			if !(i == 0 && ((hozChan && vertChan) || (fixedVertAlign && fixedHozAlign))) {
				if currentPt != nil && !fixedHozAlign && !fixedVertAlign && (inHozChan || inVertChan) {
					horizontal = !inHozChan
					break
				}
				if vertChan || hozChan {
					horizontal = hozChan
					if i == 1 {
						// work back from target end
						if len(hints)%2 == 0 {
							horizontal = hozChan
						} else {
							horizontal = vertChan
						}
					}
					break
				}
			}
			currentTerm = target
			currentPt = pts[lastInx]
			if currentPt != nil {
				currentTerm = nil
			}
			if len(hints) > 0 {
				currentHint = hints[len(hints)-1]
			}
			if fixedVertAlign && fixedHozAlign && len(hints) > 0 {
				hints = hints[1:]
			}
		}

		if horizontal && ((pts[0] != nil && pts[0].y != hint.y) ||
			(pts[0] == nil && source != nil && (hint.y < source.y || hint.y > source.y+source.h))) {
			addPoint(point{pt.x, hint.y})
		} else if !horizontal && ((pts[0] != nil && pts[0].x != hint.x) ||
			(pts[0] == nil && source != nil && (hint.x < source.x || hint.x > source.x+source.w))) {
			addPoint(point{hint.x, pt.y})
		}
		if horizontal {
			pt.y = hint.y
		} else {
			pt.x = hint.x
		}
		for _, h := range hints {
			horizontal = !horizontal
			hint = h
			if horizontal {
				pt.y = hint.y
			} else {
				pt.x = hint.x
			}
			addPoint(*pt)
		}
	} else {
		hint = pt
		horizontal = true
	}

	// the last point
	pt = pts[lastInx]
	if pt == nil && target != nil {
		pt = &point{view.getRoutingCenterX(target), view.getRoutingCenterY(target)}
	}
	if pt != nil && hint != nil {
		if horizontal && ((pts[lastInx] != nil && pts[lastInx].y != hint.y) ||
			(pts[lastInx] == nil && target != nil && (hint.y < target.y || hint.y > target.y+target.h))) {
			addPoint(point{pt.x, hint.y})
		} else if !horizontal && ((pts[lastInx] != nil && pts[lastInx].x != hint.x) ||
			(pts[lastInx] == nil && target != nil && (hint.x < target.x || hint.x > target.x+target.w))) {
			addPoint(point{hint.x, pt.y})
		}
	}

	// keeps bends inside the shape for self-loops with innerLoopWaypoints
	if sourceScaled == nil || sourceScaled != targetScaled || !styleIsOne(state.style, "innerLoopWaypoints", false) {
		// removes bends inside the source terminal
		if pts[0] == nil && source != nil {
			for len(tempPoints) > 0 && source.bounds().contains(tempPoints[0].x, tempPoints[0].y) {
				tempPoints = tempPoints[1:]
			}
		}
		// removes bends inside the target terminal
		if pts[lastInx] == nil && target != nil {
			for len(tempPoints) > 0 && target.bounds().contains(tempPoints[len(tempPoints)-1].x, tempPoints[len(tempPoints)-1].y) {
				tempPoints = tempPoints[:len(tempPoints)-1]
			}
		}
	}

	for _, p := range tempPoints {
		pushPoint(p)
	}

	// removes the last point if inside tolerance with the end point
	if n := len(result); pe != nil && n > 0 && result[n-1] != nil &&
		math.Abs(pe.x-result[n-1].x) <= tol && math.Abs(pe.y-result[n-1].y) <= tol {
		result = result[:n-1]
		// lines up the second last point with the end point
		if n := len(result); n > 0 && result[n-1] != nil {
			last := result[n-1]
			if math.Abs(last.x-pe.x) < tol {
				last.x = pe.x
			}
			if math.Abs(last.y-pe.y) < tol {
				last.y = pe.y
			}
		}
	}
	return result
}

// Orthogonal connector tables (mxEdgeStyle).
const (
	orthBuffer         = 10 // mxEdgeStyle.orthBuffer
	orthPointsFallback = true

	orthSideMask   = 480 // LEFT_MASK | TOP_MASK | RIGHT_MASK | BOTTOM_MASK
	orthCenterMask = 512
	orthSourceMask = 1024
	orthTargetMask = 2048
)

var orthDirVectors = [][2]float64{{-1, 0}, {0, -1}, {1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 0}}

var orthRoutePatterns = [4][4][]int{
	{{513, 2308, 2081, 2562}, {513, 1090, 514, 2184, 2114, 2561},
		{513, 1090, 514, 2564, 2184, 2562},
		{513, 2308, 2561, 1090, 514, 2568, 2308}},
	{{514, 1057, 513, 2308, 2081, 2562}, {514, 2184, 2114, 2561},
		{514, 2184, 2562, 1057, 513, 2564, 2184},
		{514, 1057, 513, 2568, 2308, 2561}},
	{{1090, 514, 1057, 513, 2308, 2081, 2562}, {2114, 2561},
		{1090, 2562, 1057, 513, 2564, 2184},
		{1090, 514, 1057, 513, 2308, 2561, 2568}},
	{{2081, 2562}, {1057, 513, 1090, 514, 2184, 2114, 2561},
		{1057, 513, 1090, 514, 2184, 2562, 2564},
		{1057, 2561, 1090, 514, 2568, 2308}},
}

// getJettySize returns the length of the first or last segment of an
// orthogonal edge: sourceJettySize/targetJettySize, jettySize or 10; "auto"
// makes room for the end marker (mxEdgeStyle.getJettySize).
func getJettySize(state *cellState, isSource bool) float64 {
	key, arrow, size := "targetJettySize", "endArrow", "endSize"
	if isSource {
		key, arrow, size = "sourceJettySize", "startArrow", "startSize"
	}
	value, ok := state.style[key]
	if !ok {
		value, ok = state.style["jettySize"]
	}
	if !ok {
		return orthBuffer
	}
	if value == "auto" {
		// computes the automatic jetty size
		if a, ok := state.style[arrow]; ok && a != "none" {
			sz := float64(defaultMarkerSize)
			if s, ok := state.style[size]; ok {
				sz, _ = jsNumber(s)
			}
			return math.Max(2, math.Ceil((sz+orthBuffer)/orthBuffer)) * orthBuffer
		}
		return 2 * orthBuffer
	}
	f, _ := parseFloat(value)
	return f
}

// getPortConstraints returns the sides of a terminal an edge may connect
// to, as a direction mask: the terminal's portConstraint or the edge's
// sourcePortConstraint/targetPortConstraint, turned with the terminal when
// portConstraintRotation=1 (mxUtils.getPortConstraints).
func getPortConstraints(terminal, edge *cellState, source bool, def int) int {
	value, ok := terminal.style["portConstraint"]
	if !ok {
		key := "targetPortConstraint"
		if source {
			key = "sourcePortConstraint"
		}
		value, ok = edge.style[key]
	}
	if !ok {
		return def
	}
	rotation := 0.0
	if styleIsOne(terminal.style, "portConstraintRotation", false) {
		rotation = terminal.style.num("rotation", 0)
	}
	quad := 0
	if rotation > 45 {
		quad = 1
		if rotation >= 135 {
			quad = 2
		}
	} else if rotation < -45 {
		quad = 3
		if rotation <= -135 {
			quad = 2
		}
	}
	// each side turned by quad quarters clockwise
	turn := func(masks [4]int) int { return masks[quad] }
	result := dirMaskNone
	if strings.Contains(value, "north") {
		result |= turn([4]int{dirMaskNorth, dirMaskEast, dirMaskSouth, dirMaskWest})
	}
	if strings.Contains(value, "west") {
		result |= turn([4]int{dirMaskWest, dirMaskNorth, dirMaskEast, dirMaskSouth})
	}
	if strings.Contains(value, "south") {
		result |= turn([4]int{dirMaskSouth, dirMaskWest, dirMaskNorth, dirMaskEast})
	}
	if strings.Contains(value, "east") {
		result |= turn([4]int{dirMaskEast, dirMaskSouth, dirMaskWest, dirMaskNorth})
	}
	return result
}

// reversePortConstraints swaps opposite sides of a direction mask
// (mxUtils.reversePortConstraints).
func reversePortConstraints(c int) int {
	r := (c & dirMaskWest) << 3
	r |= (c & dirMaskNorth) << 1
	r |= (c & dirMaskSouth) >> 1
	r |= (c & dirMaskEast) >> 3
	return r
}

// rotatedBounds returns the bounding box of r turned by deg degrees about
// its center (mxUtils.getBoundingBox).
func rotatedBounds(r rect, deg float64) rect {
	rad := toRadians(deg)
	cos, sin := math.Cos(rad), math.Sin(rad)
	c := r.center()
	var out rect
	for i, p := range []point{{r.x, r.y}, {r.x + r.w, r.y}, {r.x + r.w, r.y + r.h}, {r.x, r.y + r.h}} {
		q := rotatePoint(p, cos, sin, c)
		if i == 0 {
			out = rect{q.x, q.y, 0, 0}
		} else {
			out = rectAdd(out, q)
		}
	}
	return out
}

// rectAdd grows r to hold p (mxRectangle.add of an empty rectangle).
func rectAdd(r rect, p point) rect {
	x0, y0 := math.Min(r.x, p.x), math.Min(r.y, p.y)
	x1, y1 := math.Max(r.x+r.w, p.x), math.Max(r.y+r.h, p.y)
	return rect{x0, y0, x1 - x0, y1 - y0}
}

// orthConnector routes an edge with horizontal and vertical segments
// around its terminals, picking the sides from their relative positions,
// port constraints and fixed ends (mxEdgeStyle.OrthConnector). Edges with
// control points, and fixed ends closer than the jetties, fall back to
// segmentConnector.
func orthConnector(state, sourceScaled, targetScaled *cellState, controlHints []point, result []*point) []*point {
	pts := scalePointArray(state.absPoints)
	source := scaleCellState(sourceScaled)
	target := scaleCellState(targetScaled)
	p0, pe := pts[0], pts[len(pts)-1]

	var sourceX, sourceY, sourceWidth, sourceHeight float64
	if source != nil {
		sourceX, sourceY, sourceWidth, sourceHeight = source.x, source.y, source.w, source.h
	} else {
		sourceX, sourceY, sourceWidth, sourceHeight = p0.x, p0.y, 1, 1
	}
	var targetX, targetY, targetWidth, targetHeight float64
	if target != nil {
		targetX, targetY, targetWidth, targetHeight = target.x, target.y, target.w, target.h
	} else {
		targetX, targetY, targetWidth, targetHeight = pe.x, pe.y, 1, 1
	}

	sourceBuffer := getJettySize(state, true)
	targetBuffer := getJettySize(state, false)
	// workaround for loop routing within buffer zone (the scaled states
	// are copies, so this never applies, as in mxGraph)
	if source != nil && target == source {
		targetBuffer = math.Max(sourceBuffer, targetBuffer)
		sourceBuffer = targetBuffer
	}
	totalBuffer := targetBuffer + sourceBuffer
	tooShort := false
	// checks minimum distance for fixed points and falls back to segment connector
	if p0 != nil && pe != nil {
		dx, dy := pe.x-p0.x, pe.y-p0.y
		tooShort = dx*dx+dy*dy < totalBuffer*totalBuffer
	}
	if tooShort || (orthPointsFallback && len(controlHints) > 0) {
		return segmentConnector(state, sourceScaled, targetScaled, controlHints, result)
	}

	// the sides of the source and target the edge may connect to
	portConstraint := [2]int{dirMaskAll, dirMaskAll}
	if source != nil {
		portConstraint[0] = getPortConstraints(source, state, true, dirMaskAll)
		if rotation := source.style.num("rotation", 0); rotation != 0 {
			r := rotatedBounds(rect{sourceX, sourceY, sourceWidth, sourceHeight}, rotation)
			sourceX, sourceY, sourceWidth, sourceHeight = r.x, r.y, r.w, r.h
		}
	}
	if target != nil {
		portConstraint[1] = getPortConstraints(target, state, false, dirMaskAll)
		if rotation := target.style.num("rotation", 0); rotation != 0 {
			r := rotatedBounds(rect{targetX, targetY, targetWidth, targetHeight}, rotation)
			targetX, targetY, targetWidth, targetHeight = r.x, r.y, r.w, r.h
		}
	}
	if sourceWidth == 0 || sourceHeight == 0 || targetWidth == 0 || targetHeight == 0 {
		return result
	}

	var dir [2]int
	// which faces of the vertices present against each other in a way that
	// would allow a 3-segment connection if port constraints permitted
	geo := [2][4]float64{{sourceX, sourceY, sourceWidth, sourceHeight}, {targetX, targetY, targetWidth, targetHeight}}
	buffer := [2]float64{sourceBuffer, targetBuffer}
	var limits [2][9]float64
	for i := 0; i < 2; i++ {
		limits[i][1] = geo[i][0] - buffer[i]
		limits[i][2] = geo[i][1] - buffer[i]
		limits[i][4] = geo[i][0] + geo[i][2] + buffer[i]
		limits[i][8] = geo[i][1] + geo[i][3] + buffer[i]
	}

	// which quad the target is in:
	//   0 | 1
	//   -----
	//   3 | 2
	sourceCenX := geo[0][0] + geo[0][2]/2
	sourceCenY := geo[0][1] + geo[0][3]/2
	targetCenX := geo[1][0] + geo[1][2]/2
	targetCenY := geo[1][1] + geo[1][3]/2
	dx := sourceCenX - targetCenX
	dy := sourceCenY - targetCenY
	quad := 0
	if dx < 0 {
		if dy < 0 {
			quad = 2
		} else {
			quad = 1
		}
	} else if dy <= 0 {
		quad = 3
		// special case on x = 0 and negative y
		if dx == 0 {
			quad = 2
		}
	}

	// connection constraints: fixed ends pick their side
	var currentTerm *point
	if source != nil {
		currentTerm = p0
	}
	constraint := [2][2]float64{{0.5, 0.5}, {0.5, 0.5}}
	// an unattached end has no size: its top-left corner is the exact point
	if source == nil {
		constraint[0] = [2]float64{0, 0}
	}
	if target == nil {
		constraint[1] = [2]float64{0, 0}
	}
	for i := 0; i < 2; i++ {
		if currentTerm != nil {
			constraint[i][0] = (currentTerm.x - geo[i][0]) / geo[i][2]
			if math.Abs(currentTerm.x-geo[i][0]) <= 1 {
				dir[i] = dirMaskWest
			} else if math.Abs(currentTerm.x-geo[i][0]-geo[i][2]) <= 1 {
				dir[i] = dirMaskEast
			}
			constraint[i][1] = (currentTerm.y - geo[i][1]) / geo[i][3]
			if math.Abs(currentTerm.y-geo[i][1]) <= 1 {
				dir[i] = dirMaskNorth
			} else if math.Abs(currentTerm.y-geo[i][1]-geo[i][3]) <= 1 {
				dir[i] = dirMaskSouth
			}
		}
		currentTerm = nil
		if target != nil {
			currentTerm = pe
		}
	}

	sourceTopDist := geo[0][1] - (geo[1][1] + geo[1][3])
	sourceLeftDist := geo[0][0] - (geo[1][0] + geo[1][2])
	sourceBottomDist := geo[1][1] - (geo[0][1] + geo[0][3])
	sourceRightDist := geo[1][0] - (geo[0][0] + geo[0][2])
	var vertexSeperations [5]float64
	vertexSeperations[1] = math.Max(sourceLeftDist-totalBuffer, 0)
	vertexSeperations[2] = math.Max(sourceTopDist-totalBuffer, 0)
	vertexSeperations[4] = math.Max(sourceBottomDist-totalBuffer, 0)
	vertexSeperations[3] = math.Max(sourceRightDist-totalBuffer, 0)

	// the preferred orientations by relative positioning of the vertices,
	// in preferred and available order
	var dirPref, horPref, vertPref [2]int
	horPref[0], vertPref[0] = dirMaskEast, dirMaskSouth
	if sourceLeftDist >= sourceRightDist {
		horPref[0] = dirMaskWest
	}
	if sourceTopDist >= sourceBottomDist {
		vertPref[0] = dirMaskNorth
	}
	horPref[1] = reversePortConstraints(horPref[0])
	vertPref[1] = reversePortConstraints(vertPref[0])
	preferredHorizDist := math.Max(sourceLeftDist, sourceRightDist)
	preferredVertDist := math.Max(sourceTopDist, sourceBottomDist)

	var prefOrdering [2][2]int
	preferredOrderSet := false
	// if the preferred port isn't available, switch it
	for i := 0; i < 2; i++ {
		if dir[i] != 0 {
			continue
		}
		if horPref[i]&portConstraint[i] == 0 {
			horPref[i] = reversePortConstraints(horPref[i])
		}
		if vertPref[i]&portConstraint[i] == 0 {
			vertPref[i] = reversePortConstraints(vertPref[i])
		}
		prefOrdering[i] = [2]int{vertPref[i], horPref[i]}
	}
	if preferredVertDist > 0 && preferredHorizDist > 0 {
		// possibility of two segment edge connection
		if horPref[0]&portConstraint[0] > 0 && vertPref[1]&portConstraint[1] > 0 {
			prefOrdering[0] = [2]int{horPref[0], vertPref[0]}
			prefOrdering[1] = [2]int{vertPref[1], horPref[1]}
			preferredOrderSet = true
		} else if vertPref[0]&portConstraint[0] > 0 && horPref[1]&portConstraint[1] > 0 {
			prefOrdering[0] = [2]int{vertPref[0], horPref[0]}
			prefOrdering[1] = [2]int{horPref[1], vertPref[1]}
			preferredOrderSet = true
		}
	}
	if preferredVertDist > 0 && !preferredOrderSet {
		prefOrdering[0] = [2]int{vertPref[0], horPref[0]}
		prefOrdering[1] = [2]int{vertPref[1], horPref[1]}
		preferredOrderSet = true
	}
	if preferredHorizDist > 0 && !preferredOrderSet {
		prefOrdering[0] = [2]int{horPref[0], vertPref[0]}
		prefOrdering[1] = [2]int{horPref[1], vertPref[1]}
	}

	// the lists of preferred port selections, compacted (32-bit shifts as
	// in JavaScript)
	for i := 0; i < 2; i++ {
		if dir[i] != 0 {
			continue
		}
		if prefOrdering[i][0]&portConstraint[i] == 0 {
			prefOrdering[i][0] = prefOrdering[i][1]
		}
		dirPref[i] = prefOrdering[i][0] & portConstraint[i]
		dirPref[i] |= (prefOrdering[i][1] & portConstraint[i]) << 8
		dirPref[i] |= (prefOrdering[1-i][i] & portConstraint[i]) << 16
		dirPref[i] |= (prefOrdering[1-i][1-i] & portConstraint[i]) << 24
		if dirPref[i]&0xF == 0 {
			dirPref[i] = (dirPref[i] << 8) & 0xFFFFFFFF
		}
		if dirPref[i]&0xF00 == 0 {
			dirPref[i] = (dirPref[i] & 0xF) | dirPref[i]>>8
		}
		if dirPref[i]&0xF0000 == 0 {
			dirPref[i] = (dirPref[i] & 0xFFFF) | ((dirPref[i] & 0xF000000) >> 8)
		}
		dir[i] = dirPref[i] & 0xF
		switch portConstraint[i] {
		case dirMaskWest, dirMaskNorth, dirMaskEast, dirMaskSouth:
			dir[i] = portConstraint[i]
		}
	}

	sideIndex := func(d int) int {
		if d == dirMaskEast {
			return 3
		}
		return d
	}
	sourceIndex := sideIndex(dir[0]) - quad
	targetIndex := sideIndex(dir[1]) - quad
	if sourceIndex < 1 {
		sourceIndex += 4
	}
	if targetIndex < 1 {
		targetIndex += 4
	}
	routePattern := orthRoutePatterns[sourceIndex-1][targetIndex-1]

	var wayPoints [12][2]float64
	wayPoints[0] = [2]float64{geo[0][0], geo[0][1]}
	switch dir[0] {
	case dirMaskWest:
		wayPoints[0][0] -= sourceBuffer
		wayPoints[0][1] += constraint[0][1] * geo[0][3]
	case dirMaskSouth:
		wayPoints[0][0] += constraint[0][0] * geo[0][2]
		wayPoints[0][1] += geo[0][3] + sourceBuffer
	case dirMaskEast:
		wayPoints[0][0] += geo[0][2] + sourceBuffer
		wayPoints[0][1] += constraint[0][1] * geo[0][3]
	case dirMaskNorth:
		wayPoints[0][0] += constraint[0][0] * geo[0][2]
		wayPoints[0][1] -= sourceBuffer
	}

	currentIndex := 0
	// orientation: 0 horizontal, 1 vertical
	lastOrientation := 1
	if dir[0]&(dirMaskEast|dirMaskWest) > 0 {
		lastOrientation = 0
	}
	initialOrientation := lastOrientation
	for _, step := range routePattern {
		nextDirection := step & 0xF
		// rotate the index of this direction by the quad to get the real direction
		directionIndex := sideIndex(nextDirection) + quad
		if directionIndex > 4 {
			directionIndex -= 4
		}
		direction := orthDirVectors[directionIndex-1]
		currentOrientation := 1
		if directionIndex%2 > 0 {
			currentOrientation = 0
		}
		// only update the current index if the point moved in the
		// direction of the current segment move, otherwise the same point
		// is moved until there is a segment direction change
		if currentOrientation != lastOrientation {
			currentIndex++
			// copy the previous way point into the new one
			wayPoints[currentIndex] = wayPoints[currentIndex-1]
		}
		tar := step&orthTargetMask > 0
		sou := step&orthSourceMask > 0
		side := (step & orthSideMask) >> 5
		side <<= quad
		if side > 0xF {
			side >>= 4
		}
		center := step&orthCenterMask > 0
		if (sou || tar) && side < 9 {
			souTar := 1
			if sou {
				souTar = 0
			}
			var limit float64
			switch {
			case center && currentOrientation == 0:
				limit = geo[souTar][0] + constraint[souTar][0]*geo[souTar][2]
			case center:
				limit = geo[souTar][1] + constraint[souTar][1]*geo[souTar][3]
			default:
				limit = limits[souTar][side]
			}
			if currentOrientation == 0 {
				lastX := wayPoints[currentIndex][0]
				if deltaX := (limit - lastX) * direction[0]; deltaX > 0 {
					wayPoints[currentIndex][0] += direction[0] * deltaX
				}
			} else {
				lastY := wayPoints[currentIndex][1]
				if deltaY := (limit - lastY) * direction[1]; deltaY > 0 {
					wayPoints[currentIndex][1] += direction[1] * deltaY
				}
			}
		} else if center {
			// which center we're travelling to depends on the current direction
			wayPoints[currentIndex][0] += direction[0] * math.Abs(vertexSeperations[directionIndex]/2)
			wayPoints[currentIndex][1] += direction[1] * math.Abs(vertexSeperations[directionIndex]/2)
		}
		if currentIndex > 0 && wayPoints[currentIndex][currentOrientation] == wayPoints[currentIndex-1][currentOrientation] {
			currentIndex--
		} else {
			lastOrientation = currentOrientation
		}
	}

	for i := 0; i <= currentIndex; i++ {
		if i == currentIndex {
			// the last point can make the last segment run in the direction
			// of the jetty: same orientations of source and target need an
			// even number of turns (points), different ones an odd number
			targetOrientation := 1
			if dir[1]&(dirMaskEast|dirMaskWest) > 0 {
				targetOrientation = 0
			}
			sameOrient := 1
			if targetOrientation == initialOrientation {
				sameOrient = 0
			}
			if sameOrient != (currentIndex+1)%2 {
				// the last point isn't required
				break
			}
		}
		result = append(result, &point{round1(wayPoints[i][0]), round1(wayPoints[i][1])})
	}

	// removes duplicates
	for index := 1; index < len(result); {
		a, b := result[index-1], result[index]
		if a == nil || b == nil || a.x != b.x || a.y != b.y {
			index++
		} else {
			result = append(result[:index], result[index+1:]...)
		}
	}
	return result
}

// isometricConnector routes an elbow along the isometric axes through the
// first control point or the middle of the ends (mxEdgeStyle.IsometricConnector
// in Shapes.js).
func isometricConnector(state, source, target *cellState, points []point, result []*point) []*point {
	view := state.view
	var pt *point
	if len(points) > 0 {
		q := view.transformControlPoint(state, points[0])
		pt = &q
	}
	pts := state.absPoints
	p0, pe := pts[0], pts[len(pts)-1]
	if p0 == nil && source != nil {
		p0 = &point{source.cx(), source.cy()}
	}
	if pe == nil && target != nil {
		pe = &point{target.cx(), target.cy()}
	}
	if p0 == nil || pe == nil {
		return result
	}
	a1, a2 := isoHVector.x, isoHVector.y
	b1, b2 := isoVVector.x, isoVVector.y
	elbow := state.style.get("elbow", "horizontal") == "horizontal"
	last := *p0
	isoLineTo := func(x, y float64, ignoreFirst bool) {
		c1, c2 := x-last.x, y-last.y
		// solves for isometric base vectors
		h := (b2*c1 - b1*c2) / (a1*b2 - a2*b1)
		v := (a2*c1 - a1*c2) / (a2*b1 - a1*b2)
		if elbow {
			if ignoreFirst {
				last = point{last.x + a1*h, last.y + a2*h}
				q := last
				result = append(result, &q)
			}
			last = point{last.x + b1*v, last.y + b2*v}
		} else {
			if ignoreFirst {
				last = point{last.x + b1*v, last.y + b2*v}
				q := last
				result = append(result, &q)
			}
			last = point{last.x + a1*h, last.y + a2*h}
		}
		q := last
		result = append(result, &q)
	}
	if pt == nil {
		pt = &point{p0.x + (pe.x-p0.x)/2, p0.y + (pe.y-p0.y)/2}
	}
	isoLineTo(pt.x, pt.y, true)
	isoLineTo(pe.x, pe.y, false)
	return result
}

// The isometric axes: the x axis turned by -30° and by -150°.
var (
	isoHVector = rotatePoint(point{1, 0}, math.Cos(toRadians(-30)), math.Sin(toRadians(-30)), point{})
	isoVVector = rotatePoint(point{1, 0}, math.Cos(toRadians(-150)), math.Sin(toRadians(-150)), point{})
)

// sequenceSelfCallSize is the width and default height of a self-call
// (Graph.sequenceSelfCallSize).
const sequenceSelfCallSize = 30

// sequenceMessage routes a UML sequence message as a horizontal line at
// the message's own y: that of a fixed end, of the single control point,
// or the middle of the terminals, clamped to where both terminals can
// anchor. A self-call on one lifeline loops out to the right
// (mxEdgeStyle.SequenceMessage in Shapes.js).
func sequenceMessage(state, source, target *cellState, points []point, result []*point) []*point {
	view := state.view
	pts := state.absPoints
	p0, pe := pts[0], pts[len(pts)-1]

	// a self-call leaves at the y of the first control point (or the fixed
	// source) and returns at the y of the last one (or the fixed target),
	// at the x of the first control point
	if source != nil && target != nil && len(points) <= 2 &&
		view.sequenceLifeline(source.cell) == view.sequenceLifeline(target.cell) {
		var first, last *point
		if len(points) > 0 {
			q := view.transformControlPoint(state, points[0])
			first = &q
		}
		if len(points) > 1 {
			q := view.transformControlPoint(state, points[len(points)-1])
			last = &q
		}
		y1 := source.cy()
		if p0 != nil {
			y1 = p0.y
		} else if first != nil {
			y1 = first.y
		}
		y2 := y1 + sequenceSelfCallSize
		if pe != nil {
			y2 = pe.y
		} else if last != nil {
			y2 = last.y
		}
		var x float64
		if first != nil {
			x = first.x
		} else {
			x = view.sequenceSelfCallX(source, y1)
		}
		if finite(x) && finite(y1) && finite(y2) {
			result = append(result, &point{x, y1}, &point{x, y2})
		}
		return result
	}

	// a hand-routed message (more than one control point) or a
	// self-message keeps its control points
	if len(points) > 1 || source == nil || target == nil || source == target {
		for _, p := range points {
			q := view.transformControlPoint(state, p)
			result = append(result, &q)
		}
		return result
	}
	// both ends fixed: the ends define the message
	if p0 != nil && pe != nil {
		return result
	}

	var y float64
	if p0 != nil || pe != nil {
		floating := target
		if p0 != nil {
			y = p0.y
		} else {
			floating = source
			y = pe.y
		}
		// a fixed y out of reach of the floating end gives a straight line
		if y < view.sequenceAnchorY(floating, true) || y > view.sequenceAnchorY(floating, false) {
			return result
		}
	} else {
		y = (source.cy() + target.cy()) / 2
		if len(points) > 0 {
			if q := view.transformControlPoint(state, points[0]); finite(q.y) {
				y = q.y
			}
		}
		top := math.Max(view.sequenceAnchorY(source, true), view.sequenceAnchorY(target, true))
		bottom := math.Min(view.sequenceAnchorY(source, false), view.sequenceAnchorY(target, false))
		// terminals with no common range give a straight line between them
		if !(top <= bottom) {
			return result
		}
		y = math.Min(bottom, math.Max(top, y))
	}
	if finite(y) {
		result = append(result, &point{sequenceMessageX(source, target), y})
	}
	return result
}

// sequenceAnchorY returns the topmost or bottommost y at which a terminal
// anchors a connection, by probing its perimeter far above or below
// (getSequenceAnchorY in Shapes.js).
func (v *graphView) sequenceAnchorY(terminal *cellState, top bool) float64 {
	far := terminal.cy() + 1e5
	if top {
		far = terminal.cy() - 1e5
	}
	if pt := v.getPerimeterPoint(terminal, &point{terminal.cx(), far}, false, 0); pt != nil {
		return pt.y
	}
	if top {
		return terminal.y
	}
	return terminal.y + terminal.h
}

// sequenceMessageX returns the x of the routed point of a sequence message:
// the middle between the terminals, pushed out of both
// (getSequenceMessageX in Shapes.js).
func sequenceMessageX(source, target *cellState) float64 {
	x := (source.cx() + target.cx()) / 2
	dir := -1.0
	if source.cx() <= target.cx() {
		dir = 1
	}
	clear := func(terminal *cellState, side float64) {
		if x > terminal.x && x < terminal.x+terminal.w {
			if side > 0 {
				x = terminal.x + terminal.w + 1
			} else {
				x = terminal.x - 1
			}
		}
	}
	clear(source, dir)
	clear(target, -dir)
	return x
}

// sequenceSelfCallX returns the x of a self-call leaving a terminal at y:
// sequenceSelfCallSize right of its anchor (Graph.getSequenceSelfCallPoints).
func (v *graphView) sequenceSelfCallX(st *cellState, y float64) float64 {
	pt := v.getPerimeterPoint(st, &point{st.x + st.w + sequenceSelfCallSize, y}, false, 0)
	x := st.x + st.w
	if pt != nil {
		x = pt.x
	}
	return x + sequenceSelfCallSize
}

// sequenceLifeline returns the lifeline of a cell: the cell itself or the
// lifeline an activation bar is on (Graph.getSequenceLifeline).
func (v *graphView) sequenceLifeline(c *cell) *cell {
	if c == nil {
		return nil
	}
	if p := c.parent; !v.isLifeline(c) && v.isLifeline(p) {
		return p
	}
	return c
}

// isLifeline reports whether a cell is a UML lifeline (Graph.isLifeline).
func (v *graphView) isLifeline(c *cell) bool {
	if c == nil || !c.vertex {
		return false
	}
	if st := v.states[c]; st != nil {
		return st.style["shape"] == "umlLifeline"
	}
	return parseStyle(c.styleStr, false)["shape"] == "umlLifeline"
}
