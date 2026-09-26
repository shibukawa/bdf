package drawio

import (
	"math"
	"strconv"
	"strings"
)

// Swimlanes and tables: mxSwimlane with draw.io's override of its
// paintVertexShape, TableShape, TableRowShape and TableLineShape from
// Shapes.js, and the helpers of Graph.js and mxGraph that find the lines
// of a table (getTableLines, paintTableCellLines, getActualStartSize).
// The collapse icon of swimlanes is not painted (it is part of the editor).

func init() {
	registerShape("swimlane", &shapeDef{
		paintVertex: func(s *shape, c *c2d, x, y, w, h float64) { paintSwimlaneShape(s, c, x, y, w, h, true) },
		labelBounds: func(s *shape, r rect) rect { return swimlaneLabelBounds(s, r, true) },
		roundable:   true,
	})
	registerShape("table", tableShape(false))
	registerShape("tableRow", tableShape(true))
	// TableLineShape paints the line it is constructed with; made from a
	// style it has none
	registerShape("tableLine", &shapeDef{paintVertex: func(*shape, *c2d, float64, float64, float64, float64) {}})
}

// defaultStartSize is mxConstants.DEFAULT_STARTSIZE.
const defaultStartSize = 40

// swimlaneTitleSize is mxSwimlane.getTitleSize.
func swimlaneTitleSize(s *shape) float64 {
	return math.Max(0, s.style.num("startSize", defaultStartSize))
}

// swimlaneHorizontal is mxSwimlane.isHorizontal.
func swimlaneHorizontal(s *shape) bool { return s.style.num("horizontal", 1) == 1 }

// swimlaneLaneFill is mxSwimlane.laneFill ("" for none).
func swimlaneLaneFill(s *shape) string { return colorOrNone(s.style.get("swimlaneFillColor", "none")) }

// swimlaneLabelBounds is mxSwimlane.getLabelBounds: the title area.
func swimlaneLabelBounds(s *shape, r rect, fixedHeaderDefault bool) rect {
	start := swimlaneTitleSize(s)
	if start == 0 && !styleBool(s.style, "fixedHeader", fixedHeaderDefault) {
		return shapeLabelBounds(s, r)
	}
	flipH := s.style.num("flipH", 0) == 1
	flipV := s.style.num("flipV", 0) == 1
	horizontal := swimlaneHorizontal(s)
	// east is the default
	shapeVertical := s.direction == "north" || s.direction == "south"
	realHorizontal := horizontal == !shapeVertical
	southWest := s.direction == "south" || s.direction == "west"
	realFlipH := !realHorizontal && flipH != southWest
	realFlipV := realHorizontal && flipV != southWest
	if !shapeVertical {
		tmp := math.Min(r.h, start)
		if realFlipH || realFlipV {
			r.y += r.h - tmp
		}
		r.h = tmp
	} else {
		tmp := math.Min(r.w, start)
		if realFlipH || realFlipV {
			r.x += r.w - tmp
		}
		r.w = tmp
	}
	return r
}

// swimlaneArcSize is mxSwimlane.getSwimlaneArcSize.
func swimlaneArcSize(s *shape, w, h, start float64) float64 {
	if s.style.num("absoluteArcSize", 0) == 1 {
		return math.Min(w/2, math.Min(h/2, s.style.num("arcSize", lineArcSize)/2))
	}
	f := s.style.num("arcSize", rectangleRoundingFactor*100) / 100
	return start * f * 3
}

// paintSwimlaneShape is mxSwimlane.paintVertexShape with draw.io's
// override, which repaints the table lines a filled row or cell hides.
func paintSwimlaneShape(s *shape, c *c2d, x, y, w, h float64, fixedHeaderDefault bool) {
	start := swimlaneTitleSize(s)
	if start != 0 || styleBool(s.style, "fixedHeader", fixedHeaderDefault) {
		horizontal := swimlaneHorizontal(s)
		if horizontal {
			start = math.Min(start, h)
		} else {
			start = math.Min(start, w)
		}
		c.translate(x, y)
		r := 0.0
		if !s.isRounded {
			paintSwimlane(s, c, w, h, start)
		} else {
			r = swimlaneArcSize(s, w, h, start)
			side := w
			if horizontal {
				side = h
			}
			r = math.Min(side-start, math.Min(start, r))
			paintRoundedSwimlane(s, c, w, h, start, r)
		}
		paintSwimlaneFooter(s, c, w, h, start, r)
		paintSwimlaneSeparator(s, c, w, h, start, colorOrNone(s.style.get("separatorColor", "none")))
		if src := s.style.get("image", ""); src != "" {
			if img := s.conv.imageFor(src); img != nil {
				// getImageBounds: the top right (top left when vertical) corner
				const imageSize = 16
				ix := 0.0
				if horizontal {
					ix = w - imageSize
				}
				c.image(ix, 0, imageSize, imageSize, img, false, false, false)
			}
		}
		if s.glass {
			s.paintGlassEffect(c, 0, 0, w, start, r)
		}
	}
	if s.fill != "" || swimlaneLaneFill(s) != "" {
		// a table row only fills its title strip, so only the lines there are hidden
		cl := s.st.cell
		if isTableRow(s.st.view, cl) {
			start := swimlaneTitleSize(s)
			if swimlaneHorizontal(s) {
				paintTableCellLines(s, c, cl, 0, 0, w, math.Min(start, h), s.stroke, s.strokewidth)
			} else {
				paintTableCellLines(s, c, cl, 0, 0, math.Min(start, w), h, s.stroke, s.strokewidth)
			}
		} else {
			paintTableCellLines(s, c, cl, 0, 0, w, h, s.stroke, s.strokewidth)
		}
	}
}

// swimlaneParts reads swimlaneLine, swimlaneHead and swimlaneBody.
func swimlaneParts(s *shape) (line, head, body bool) {
	return s.style.num("swimlaneLine", 1) == 1, s.style.num("swimlaneHead", 1) == 1, s.style.num("swimlaneBody", 1) == 1
}

// paintSwimlaneBody fills and strokes the content area path (the stroke
// only without a lane fill).
func paintSwimlaneBody(c *c2d, fill string, body bool) {
	if body {
		if fill == "" {
			c.stroke()
		} else {
			c.fillAndStroke()
		}
	} else if fill != "" {
		c.fill()
	}
}

func paintSwimlaneHead(c *c2d, head bool) {
	if head {
		c.fillAndStroke()
	} else {
		c.fill()
	}
}

// paintSwimlane is mxSwimlane.paintSwimlane: the title and the content
// area in coordinates relative to the bounds.
func paintSwimlane(s *shape, c *c2d, w, h, start float64) {
	fill := swimlaneLaneFill(s)
	line, head, body := swimlaneParts(s)
	if swimlaneHorizontal(s) {
		c.begin()
		c.moveTo(0, start)
		c.lineTo(0, 0)
		c.lineTo(w, 0)
		c.lineTo(w, start)
		paintSwimlaneHead(c, head)
		if start < h {
			if fill != "" {
				c.setFillColor(fill)
			}
			c.begin()
			c.moveTo(0, start)
			c.lineTo(0, h)
			c.lineTo(w, h)
			c.lineTo(w, start)
			paintSwimlaneBody(c, fill, body)
		}
	} else {
		c.begin()
		c.moveTo(start, 0)
		c.lineTo(0, 0)
		c.lineTo(0, h)
		c.lineTo(start, h)
		paintSwimlaneHead(c, head)
		if start < w {
			if fill != "" {
				c.setFillColor(fill)
			}
			c.begin()
			c.moveTo(start, 0)
			c.lineTo(w, 0)
			c.lineTo(w, h)
			c.lineTo(start, h)
			paintSwimlaneBody(c, fill, body)
		}
	}
	if line {
		paintSwimlaneDivider(s, c, w, h, start)
	}
}

// paintRoundedSwimlane is mxSwimlane.paintRoundedSwimlane.
func paintRoundedSwimlane(s *shape, c *c2d, w, h, start, r float64) {
	fill := swimlaneLaneFill(s)
	line, head, body := swimlaneParts(s)
	if swimlaneHorizontal(s) {
		c.begin()
		c.moveTo(w, start)
		c.lineTo(w, r)
		c.quadTo(w, 0, w-math.Min(w/2, r), 0)
		c.lineTo(math.Min(w/2, r), 0)
		c.quadTo(0, 0, 0, r)
		c.lineTo(0, start)
		paintSwimlaneHead(c, head)
		if start < h {
			if fill != "" {
				c.setFillColor(fill)
			}
			c.begin()
			c.moveTo(0, start)
			c.lineTo(0, h-r)
			c.quadTo(0, h, math.Min(w/2, r), h)
			c.lineTo(w-math.Min(w/2, r), h)
			c.quadTo(w, h, w, h-r)
			c.lineTo(w, start)
			paintSwimlaneBody(c, fill, body)
		}
	} else {
		c.begin()
		c.moveTo(start, 0)
		c.lineTo(r, 0)
		c.quadTo(0, 0, 0, math.Min(h/2, r))
		c.lineTo(0, h-math.Min(h/2, r))
		c.quadTo(0, h, r, h)
		c.lineTo(start, h)
		paintSwimlaneHead(c, head)
		if start < w {
			if fill != "" {
				c.setFillColor(fill)
			}
			c.begin()
			c.moveTo(start, h)
			c.lineTo(w-r, h)
			c.quadTo(w, h, w, h-math.Min(h/2, r))
			c.lineTo(w, math.Min(h/2, r))
			c.quadTo(w, 0, w-r, 0)
			c.lineTo(start, 0)
			paintSwimlaneBody(c, fill, body)
		}
	}
	if line {
		paintSwimlaneDivider(s, c, w, h, start)
	}
}

// paintSwimlaneFooter is mxSwimlane.paintFooter: footerSize paints a
// filled band at the end opposite the title.
func paintSwimlaneFooter(s *shape, c *c2d, w, h, start, r float64) {
	footer := math.Max(0, s.style.num("footerSize", 0))
	horizontal := swimlaneHorizontal(s)
	size := w
	if horizontal {
		size = h
	}
	// the footer never overlaps the title
	footer = math.Min(footer, size-start)
	if footer <= 0 {
		return
	}
	if s.fill != "" {
		c.setFillColor(s.fill)
	}
	c.begin()
	if horizontal {
		rr := math.Min(w/2, r)
		if footer >= r {
			c.moveTo(0, h-footer)
			c.lineTo(0, h-r)
			c.quadTo(0, h, rr, h)
			c.lineTo(w-rr, h)
			c.quadTo(w, h, w, h-r)
			c.lineTo(w, h-footer)
		} else {
			// the divider is inside the corners: trace the wall from where
			// its quadratic crosses the divider
			t := 1 - math.Sqrt(footer/r)
			px, bx := rr*t*t, rr*t
			c.moveTo(px, h-footer)
			c.quadTo(bx, h, rr, h)
			c.lineTo(w-rr, h)
			c.quadTo(w-bx, h, w-px, h-footer)
		}
	} else {
		rr := math.Min(h/2, r)
		if footer >= r {
			c.moveTo(w-footer, 0)
			c.lineTo(w-r, 0)
			c.quadTo(w, 0, w, rr)
			c.lineTo(w, h-rr)
			c.quadTo(w, h, w-r, h)
			c.lineTo(w-footer, h)
		} else {
			t := 1 - math.Sqrt(footer/r)
			py, by := rr*t*t, rr*t
			c.moveTo(w-footer, py)
			c.quadTo(w, by, w, rr)
			c.lineTo(w, h-rr)
			c.quadTo(w, h-by, w-footer, h-py)
		}
	}
	c.close()
	if s.fill != "" {
		c.fillAndStroke()
	} else {
		c.stroke()
	}
}

// paintSwimlaneDivider is mxSwimlane.paintDivider: the line under the title.
func paintSwimlaneDivider(s *shape, c *c2d, w, h, start float64) {
	if start == 0 {
		return
	}
	c.begin()
	if swimlaneHorizontal(s) {
		c.moveTo(0, start)
		c.lineTo(w, start)
	} else {
		c.moveTo(start, 0)
		c.lineTo(start, h)
	}
	c.stroke()
}

// paintSwimlaneSeparator is mxSwimlane.paintSeparator: a dashed line in
// separatorColor at the far side.
func paintSwimlaneSeparator(s *shape, c *c2d, w, h, start float64, color string) {
	if color == "" {
		return
	}
	c.setStrokeColor(color)
	c.setDashed(true, false)
	c.begin()
	if swimlaneHorizontal(s) {
		c.moveTo(w, start)
		c.lineTo(w, h)
	} else {
		c.moveTo(start, 0)
		c.lineTo(w, 0)
	}
	c.stroke()
	c.setDashed(false, false)
}

// tableShape ports TableShape ("table") and TableRowShape ("tableRow",
// which does not paint the table lines). A table is a swimlane that is
// not painted as one without a title unless fixedHeader is set.
func tableShape(row bool) *shapeDef {
	return &shapeDef{
		paintVertex: func(s *shape, c *c2d, x, y, w, h float64) {
			horizontal := swimlaneHorizontal(s)
			start := swimlaneTitleSize(s)
			fixedHeader := styleBool(s.style, "fixedHeader", false)
			switch {
			case start == 0 && s.isRounded:
				// headerless rounded table
				r := s.arcSize(w, h)
				c.begin()
				c.roundrect(x, y, w, h, r, r)
				c.fillAndStroke()
			case start == 0 && !fixedHeader:
				paintPartialRectangle(s, c, x, y, w, h)
			default:
				paintSwimlaneShape(s, c, x, y, w, h, false)
				c.translate(-x, -y)
			}
			if !row && !s.st.cell.collapsed && (horizontal && start < h || !horizontal && start < w) {
				paintTableForeground(s, c, x, y, w, h)
			}
		},
		labelBounds: func(s *shape, r rect) rect { return swimlaneLabelBounds(s, r, false) },
		roundable:   true,
	}
}

// paintTableForeground is TableShape.paintForeground: the row and column
// lines, in the table's unrotated bounds.
func paintTableForeground(s *shape, c *c2d, x, y, w, h float64) {
	flipH, flipV := s.flipH, s.flipV
	if s.direction == "north" || s.direction == "south" {
		flipH, flipV = flipV, flipH
	}
	// a negative transform instead of save and restore
	c.rotate(-s.shapeRotation(), flipH, flipV, x+w/2, y+h/2)
	b := s.bounds
	lines := tableLines(s.st.view, s.st.cell, s.style.num("rowLines", 1) != 0, s.style.num("columnLines", 1) != 0)
	for _, l := range lines {
		paintTableLine(c, l, b.x, b.y)
	}
}

// paintTableLine is TableLineShape.paintTableLine: a polyline with gaps
// (nil points).
func paintTableLine(c *c2d, line []*point, dx, dy float64) {
	if line == nil {
		return
	}
	var last *point
	c.begin()
	for _, curr := range line {
		if curr != nil {
			if last == nil {
				c.moveTo(curr.x+dx, curr.y+dy)
			} else {
				c.lineTo(curr.x+dx, curr.y+dy)
			}
		}
		last = curr
	}
	c.stroke()
}

// tableIter is a cell visited by visitTableCells.
type tableIter struct {
	cell             *cell
	rowspan, colspan int
	row, col         int
	point            point // the bottom right corner, relative to the table
	actual           *tableIter
}

// visitTableCells calls visitor for each cell of a table with the corner
// of the cell and the origin of the lines (Graph.visitTableCells).
func visitTableCells(v *graphView, table *cell, visitor func(it *tableIter, colCount, rowCount int, x0, y0 float64)) {
	var lastRow []*tableIter
	rows := childVertices(table)
	start := actualStartSize(v, table)
	for i, rc := range rows {
		rowStart := actualStartSize(v, rc)
		cols := childVertices(rc)
		rowStyle := modelStyle(rc)
		var lastCol *tableIter
		var row []*tableIter
		for j, cc := range cols {
			geo := cc.geo
			if geo == nil {
				// keeps lastRow[j] aligned with the columns
				row = append(row, nil)
				continue
			}
			gw, gh := geo.w, geo.h
			if geo.alternateBounds != nil {
				gw, gh = geo.alternateBounds.w, geo.alternateBounds.h
			}
			it := &tableIter{cell: cc, row: i, col: j}
			px := start.x + rowStart.x
			if lastCol != nil {
				px = lastCol.point.x
			}
			py := start.y + rowStart.y
			if lastRow != nil && len(lastRow) > 0 && lastRow[0] != nil {
				py = lastRow[0].point.y
			}
			it.point = point{gw + px, gh + py}
			it.actual = it
			switch {
			case j < len(lastRow) && lastRow[j] != nil && lastRow[j].rowspan > 1:
				it.rowspan = lastRow[j].rowspan - 1
				it.colspan = lastRow[j].colspan
				it.actual = lastRow[j].actual
			case lastCol != nil && lastCol.colspan > 1:
				it.rowspan = lastCol.rowspan
				it.colspan = lastCol.colspan - 1
				it.actual = lastCol.actual
			default:
				st := modelStyle(cc)
				it.rowspan = styleInt(st, "rowspan", 1)
				it.colspan = styleInt(st, "colspan", 1)
			}
			head := rowStyle.num("swimlaneHead", 1) == 1 && colorOrNone(rowStyle.get("strokeColor", "none")) != ""
			x0, y0 := start.x, start.y
			if head {
				x0 += rowStart.x
				y0 += rowStart.y
			}
			visitor(it, len(cols), len(rows), x0, y0)
			row = append(row, it)
			lastCol = it
		}
		lastRow = row
	}
}

// tableLines returns the row lines and then the column lines of a table
// (Graph.getTableLines).
func tableLines(v *graphView, table *cell, horizontal, vertical bool) [][]*point {
	var hl, vl [][]*point
	grow := func(l [][]*point, i int) [][]*point {
		for len(l) <= i {
			l = append(l, nil)
		}
		return l
	}
	if horizontal || vertical {
		visitTableCells(v, table, func(it *tableIter, colCount, rowCount int, x0, y0 float64) {
			p := it.point
			if horizontal && it.row < rowCount-1 {
				hl = grow(hl, it.row)
				if hl[it.row] == nil {
					hl[it.row] = []*point{{x0, p.y}}
				}
				if it.rowspan > 1 {
					hl[it.row] = append(hl[it.row], nil)
				}
				hl[it.row] = append(hl[it.row], &p)
			}
			if vertical && it.col < colCount-1 {
				vl = grow(vl, it.col)
				if vl[it.col] == nil {
					vl[it.col] = []*point{{p.x, y0}}
				}
				if it.colspan > 1 {
					vl[it.col] = append(vl[it.col], nil)
				}
				vl[it.col] = append(vl[it.col], &p)
			}
		})
	}
	return append(hl, vl...)
}

// paintTableCellLines strokes the table lines along the sides of a filled
// table row or cell over its fill (Graph.paintTableCellLines): the table
// paints them below its children, where a fill hides them.
func paintTableCellLines(s *shape, c *c2d, cl *cell, x, y, w, h float64, stroke string, strokeWidth float64) {
	if stroke == "" || cl == nil || s.st == nil {
		return
	}
	v := s.st.view
	var top, bottom, left, right bool
	switch {
	case isTableRow(v, cl):
		rows := childVertices(cl.parent)
		i := cellIndex(rows, cl)
		rowLines := currentCellStyle(v, cl.parent).num("rowLines", 1) != 0
		top = rowLines && i > 0
		bottom = rowLines && i < len(rows)-1
	case isTableCell(v, cl):
		row := cl.parent
		table := row.parent
		rows := childVertices(table)
		cols := childVertices(row)
		st := currentCellStyle(v, table)
		rowLines := st.num("rowLines", 1) != 0
		columnLines := st.num("columnLines", 1) != 0
		ri, ci := cellIndex(rows, row), cellIndex(cols, cl)
		top = rowLines && ri > 0
		bottom = rowLines && ri < len(rows)-1
		left = columnLines && ci > 0
		right = columnLines && ci < len(cols)-1
	default:
		return
	}
	if !top && !bottom && !left && !right {
		return
	}
	c.setStrokeColor(stroke)
	c.setStrokeWidth(strokeWidth)
	c.setDashed(false, false)
	line := func(x0, y0, x1, y1 float64) {
		c.begin()
		c.moveTo(x0, y0)
		c.lineTo(x1, y1)
		c.stroke()
	}
	if top {
		line(x, y, x+w, y)
	}
	if bottom {
		line(x, y+h, x+w, y+h)
	}
	if left {
		line(x, y, x, y+h)
	}
	if right {
		line(x+w, y, x+w, y+h)
	}
}

// --- color keywords ---

// resolveColors replaces the keywords inherit, fillColor, strokeColor and
// swimlane in the gradient, stroke and fill colors of a shape, as table
// rows and cells use them (mxCellRenderer.postConfigureShape).
func (s *shape) resolveColors() {
	if s.st == nil {
		return
	}
	s.gradient = resolveColor(s.st, s.gradient, "gradientColor")
	s.stroke = resolveColor(s.st, s.stroke, "strokeColor")
	s.fill = resolveColor(s.st, s.fill, "fillColor")
}

// resolveColor resolves a color keyword in the state's field color
// (mxCellRenderer.resolveColor): inherit takes the color of the parent's
// shape, fillColor and strokeColor the parent's style value of that key,
// swimlane the color of the nearest swimlane's shape.
func resolveColor(st *cellState, value, field string) string {
	key := field
	var referenced *cell
	switch value {
	case "inherit", "fillColor", "strokeColor":
		referenced = st.cell.parent
		if value != "inherit" {
			key = value
		}
	case "swimlane":
		value = "#ffffff"
		if field == "strokeColor" {
			value = "#000000"
		}
		referenced = st.cell
		if st.cell.target != nil {
			referenced = st.cell.target
		}
		referenced = swimlaneOf(st.view, referenced)
		key = "fillColor" // mxGraph.swimlaneIndicatorColorAttribute
	default:
		return value
	}
	if referenced == nil {
		return value
	}
	rst := st.view.state(referenced)
	if rst == nil || referenced == st.cell {
		return "" // a shape that refers to itself reads its color after clearing it
	}
	if (referenced.vertex || referenced.edge) && value != "strokeColor" && value != "fillColor" {
		// the color of the referenced shape, resolved in turn
		return resolveColor(rst, colorOrNone(rst.style.get(field, "")), field)
	}
	return colorOrNone(rst.style.get(key, ""))
}

// swimlaneOf returns the nearest swimlane at or above a cell (mxGraph.getSwimlane).
func swimlaneOf(v *graphView, c *cell) *cell {
	for c != nil {
		if st := v.state(c); st != nil && !c.edge && c.parent != v.m.root && c != v.m.root {
			switch st.style.get("shape", "") {
			case "swimlane", "table", "tableRow":
				return c
			}
		}
		c = c.parent
	}
	return nil
}

// --- graph helpers ---

// modelStyle is the style of a cell from the model, without its state
// (mxGraph.getCellStyle).
func modelStyle(c *cell) style {
	if c == nil {
		return parseStyle("", false)
	}
	return parseStyle(c.styleStr, c.edge)
}

// currentCellStyle is the style of the cell's state, or else its model
// style (mxGraph.getCurrentCellStyle).
func currentCellStyle(v *graphView, c *cell) style {
	if st := v.state(c); st != nil {
		return st.style
	}
	return modelStyle(c)
}

// childVertices returns the vertex children of a cell (mxGraphModel.getChildCells).
func childVertices(c *cell) []*cell {
	if c == nil {
		return nil
	}
	var out []*cell
	for _, k := range c.children {
		if k.vertex {
			out = append(out, k)
		}
	}
	return out
}

func cellIndex(cells []*cell, c *cell) int {
	for i, k := range cells {
		if k == c {
			return i
		}
	}
	return -1
}

// isTable reports whether a cell is a table (Graph.isTable).
func isTable(c *cell) bool {
	return c != nil && modelStyle(c).get("childLayout", "") == "tableLayout"
}

// isTableRow is Graph.isTableRow.
func isTableRow(v *graphView, c *cell) bool { return c != nil && c.vertex && isTable(c.parent) }

// isTableCell is Graph.isTableCell.
func isTableCell(v *graphView, c *cell) bool { return c != nil && c.vertex && isTableRow(v, c.parent) }

// isSwimlaneCell is Graph.isSwimlane from the model style: swimlanes,
// tables and table rows below the layers.
func isSwimlaneCell(v *graphView, c *cell) bool {
	if c == nil || c.edge || c.parent == v.m.root {
		return false
	}
	switch modelStyle(c).get("shape", "") {
	case "swimlane", "table", "tableRow":
		return true
	}
	return false
}

// swimlaneDirection is mxGraph.getSwimlaneDirection: the side of the
// title.
func swimlaneDirection(st style) string {
	dir := st.get("direction", "east")
	flipH := st.num("flipH", 0) == 1
	flipV := st.num("flipV", 0) == 1
	n := 3
	if styleBool(st, "horizontal", true) {
		n = 0
	}
	switch dir {
	case "north":
		n--
	case "west":
		n += 2
	case "south":
		n++
	}
	m := mod(n, 2)
	if flipH && m == 1 {
		n += 2
	}
	if flipV && m == 0 {
		n += 2
	}
	return [4]string{"north", "east", "south", "west"}[mod(n, 4)]
}

// actualStartSize is mxGraph.getActualStartSize from the model style: the
// title size of a swimlane on its side (as left, top, right, bottom in x,
// y, w, h).
func actualStartSize(v *graphView, c *cell) rect {
	var r rect
	if !isSwimlaneCell(v, c) {
		return r
	}
	st := modelStyle(c)
	size := math.Trunc(st.num("startSize", defaultStartSize))
	switch swimlaneDirection(st) {
	case "north":
		r.y = size
	case "west":
		r.x = size
	case "south":
		r.h = size
	default:
		r.w = size
	}
	return r
}

// styleInt is parseInt(style[key] || def).
func styleInt(st style, key string, def int) int {
	v := strings.TrimSpace(st.get(key, ""))
	if v == "" {
		return def
	}
	end := 0
	for end < len(v) && (v[end] >= '0' && v[end] <= '9' || end == 0 && (v[0] == '-' || v[0] == '+')) {
		end++
	}
	n, err := strconv.Atoi(v[:end])
	if err != nil {
		return def
	}
	return n
}

// --- label bounds ---

// shapeLabelBounds is mxShape.getLabelBounds: the label rectangle inset
// by the shape's label margins, which are turned for the direction and
// flips of the shape. Shapes that override getLabelBounds and call the
// base implementation use it; so can the label layout for shapes without
// labelBounds.
func shapeLabelBounds(s *shape, r rect) rect {
	if s.def.labelMargins == nil {
		return r
	}
	// mxText.isPaintBoundsInverted: vertical labels of vertices
	textInverted := s.st != nil && s.st.cell.vertex && !s.style.flag("horizontal", true)
	bounds := r
	if s.direction != "south" && s.direction != "north" && textInverted {
		bounds.w, bounds.h = bounds.h, bounds.w
	}
	mp := s.def.labelMargins(s, bounds)
	if mp == nil {
		return r
	}
	m := *mp
	flipH := s.style.get("flipH", "0") == "1"
	flipV := s.style.get("flipV", "0") == "1"
	if textInverted {
		m = rect{m.h, m.x, m.y, m.w}
		flipH, flipV = flipV, flipH
	}
	return directedBounds(r, m, s.direction, flipH, flipV)
}

// directedBounds insets rect by the margins m (left, top, right, bottom in
// x, y, w, h) given for the east direction, turned for dir and the flips
// (mxUtils.getDirectedBounds).
func directedBounds(r rect, m rect, dir string, flipH, flipV bool) rect {
	m.x = jsRound(math.Max(0, math.Min(r.w, m.x)))
	m.y = jsRound(math.Max(0, math.Min(r.h, m.y)))
	m.w = jsRound(math.Max(0, math.Min(r.w, m.w)))
	m.h = jsRound(math.Max(0, math.Min(r.h, m.h)))
	vertical := dir == "south" || dir == "north"
	if flipV && vertical || flipH && !vertical {
		m.x, m.w = m.w, m.x
	}
	if flipH && vertical || flipV && !vertical {
		m.y, m.h = m.h, m.y
	}
	m2 := m
	switch dir {
	case "south":
		m2 = rect{m.h, m.x, m.y, m.w}
	case "west":
		m2 = rect{m.w, m.h, m.x, m.y}
	case "north":
		m2 = rect{m.y, m.w, m.h, m.x}
	}
	return rect{r.x + m2.x, r.y + m2.y, r.w - m2.w - m2.x, r.h - m2.h - m2.y}
}
