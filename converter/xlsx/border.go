package xlsx

import (
	"fmt"
	"sort"

	"github.com/shibukawa/bdf"
)

// Borders are drawn on the grid lines, centered on them. A line between two
// cells can be set on either; the heavier style wins. Lines inside merged
// ranges are not drawn. Collinear edges of one style join into one line, and
// the lines of a tile are stroked as one path per style and color.

type lineStyle struct {
	w      float64   // points
	dash   []float64 // points; nil for solid
	weight int
	double bool
}

// borderStyles are the line styles at 100% zoom (1 px = 0.75 pt).
var borderStyles = map[string]lineStyle{
	"hair":             {0.5, []float64{0.75, 0.75}, 1, false},
	"dotted":           {0.75, []float64{0.75, 0.75}, 2, false},
	"dashDotDot":       {0.75, []float64{6.75, 2.25, 2.25, 2.25, 2.25, 2.25}, 3, false},
	"dashDot":          {0.75, []float64{6.75, 2.25, 2.25, 2.25}, 4, false},
	"dashed":           {0.75, []float64{2.25, 0.75}, 5, false},
	"thin":             {0.75, nil, 6, false},
	"mediumDashDotDot": {1.5, []float64{6.75, 2.25, 2.25, 2.25, 2.25, 2.25}, 7, false},
	"slantDashDot":     {1.5, []float64{7.5, 0.75, 3, 0.75}, 8, false},
	"mediumDashDot":    {1.5, []float64{6.75, 2.25, 2.25, 2.25}, 9, false},
	"mediumDashed":     {1.5, []float64{6.75, 2.25}, 10, false},
	"medium":           {1.5, nil, 11, false},
	"double":           {0.75, nil, 12, true},
	"thick":            {2.25, nil, 13, false},
}

type edge struct {
	style string
	color rgb
}

func (e edge) weight() int { return borderStyles[e.style].weight }

type edgeSet struct {
	s     *sheetCtx
	h, v  map[[2]int]edge // top edge of (row, col); left edge of (row, col)
	diags []diagLine
}

type diagLine struct {
	b        box
	e        edge
	up, down bool
}

func (es *edgeSet) put(m map[[2]int]edge, k [2]int, side borderSide) {
	if _, ok := borderStyles[side.style]; !ok {
		return
	}
	e := edge{style: side.style, color: es.s.c.st.color(side.color, black)}
	if old, ok := m[k]; ok && old.weight() >= e.weight() {
		return
	}
	m[k] = e
}

// sides adds the borders of a cell's format; top, left, bottom and right
// tell which sides are outer edges (all but merged ranges' inner ones).
func (es *edgeSet) sides(r, c int, f *cellFmt, top, left, bottom, right bool) {
	b := &f.border
	if top {
		es.put(es.h, [2]int{r, c}, b.top)
	}
	if bottom {
		es.put(es.h, [2]int{r + 1, c}, b.bottom)
	}
	if left {
		es.put(es.v, [2]int{r, c}, b.left)
	}
	if right {
		es.put(es.v, [2]int{r, c + 1}, b.right)
	}
}

func hasBorder(f *cellFmt) bool {
	b := &f.border
	return b.top.style != "" || b.bottom.style != "" || b.left.style != "" || b.right.style != "" ||
		(b.diagUp || b.diagDown) && b.diag.style != ""
}

// paintBorders collects the edges of every formatted cell and draws them.
func (s *sheetCtx) paintBorders() {
	es := &edgeSet{s: s, h: map[[2]int]edge{}, v: map[[2]int]edge{}}
	for r := 0; r < s.nRows; r++ {
		if s.rows.at(r) == 0 {
			continue
		}
		s.eachFormatted(r, func(c0, c1 int, f *cellFmt) {
			if f == nil || !hasBorder(f) {
				return
			}
			for c := c0; c <= c1; c++ {
				es.sides(r, c, f, true, true, true, true)
				if f.border.diag.style != "" && (f.border.diagUp || f.border.diagDown) {
					es.diags = append(es.diags, diagLine{b: s.cellBox(r, c), e: edge{f.border.diag.style, s.c.st.color(f.border.diag.color, black)},
						up: f.border.diagUp, down: f.border.diagDown})
				}
			}
		})
	}
	for _, m := range s.merges {
		for r := m.r0; r <= m.r1; r++ {
			rw := s.ws.rowAt(r)
			for c := m.c0; c <= m.c1; c++ {
				if r != m.r0 && r != m.r1 && c != m.c0 && c != m.c1 {
					continue
				}
				f := s.formatAt(r, c, rw, s.ws.cellAt(r, c))
				if f == nil || !hasBorder(f) {
					continue
				}
				es.sides(r, c, f, r == m.r0, c == m.c0, r == m.r1, c == m.c1)
				if r == m.r0 && c == m.c0 && f.border.diag.style != "" && (f.border.diagUp || f.border.diagDown) {
					es.diags = append(es.diags, diagLine{b: s.rangeBox(m), e: edge{f.border.diag.style, s.c.st.color(f.border.diag.color, black)},
						up: f.border.diagUp, down: f.border.diagDown})
				}
			}
		}
	}
	es.draw()
}

type segment struct {
	x0, y0, x1, y1 float64
}

// draw joins the edges into lines and strokes them per tile.
func (es *edgeSet) draw() {
	s := es.s
	groups := map[edge][]segment{}
	// horizontal: by row, then column
	hk := make([][2]int, 0, len(es.h))
	for k := range es.h {
		hk = append(hk, k)
	}
	sort.Slice(hk, func(i, j int) bool {
		if hk[i][0] != hk[j][0] {
			return hk[i][0] < hk[j][0]
		}
		return hk[i][1] < hk[j][1]
	})
	for i := 0; i < len(hk); {
		k := hk[i]
		e := es.h[k]
		j := i + 1
		for j < len(hk) && hk[j][0] == k[0] && hk[j][1] == hk[j-1][1]+1 && es.h[hk[j]] == e {
			j++
		}
		if k[0] <= s.nRows && k[1] < s.nCols {
			y := s.rows.pos(k[0])
			x0, x1 := s.cols.pos(k[1]), s.cols.pos(min(hk[j-1][1]+1, s.nCols))
			if x1 > x0 {
				groups[e] = append(groups[e], segment{x0, y, x1, y})
			}
		}
		i = j
	}
	vk := make([][2]int, 0, len(es.v))
	for k := range es.v {
		vk = append(vk, k)
	}
	sort.Slice(vk, func(i, j int) bool {
		if vk[i][1] != vk[j][1] {
			return vk[i][1] < vk[j][1]
		}
		return vk[i][0] < vk[j][0]
	})
	for i := 0; i < len(vk); {
		k := vk[i]
		e := es.v[k]
		j := i + 1
		for j < len(vk) && vk[j][1] == k[1] && vk[j][0] == vk[j-1][0]+1 && es.v[vk[j]] == e {
			j++
		}
		if k[1] <= s.nCols && k[0] < s.nRows {
			x := s.cols.pos(k[1])
			y0, y1 := s.rows.pos(k[0]), s.rows.pos(min(vk[j-1][0]+1, s.nRows))
			if y1 > y0 {
				groups[e] = append(groups[e], segment{x, y0, x, y1})
			}
		}
		i = j
	}
	for _, d := range es.diags {
		if d.down {
			groups[d.e] = append(groups[d.e], segment{d.b.x, d.b.y, d.b.x + d.b.w, d.b.y + d.b.h})
		}
		if d.up {
			groups[d.e] = append(groups[d.e], segment{d.b.x, d.b.y + d.b.h, d.b.x + d.b.w, d.b.y})
		}
	}
	// deterministic order of the groups
	keys := make([]edge, 0, len(groups))
	for e := range groups {
		keys = append(keys, e)
	}
	sort.Slice(keys, func(i, j int) bool {
		a, b := keys[i], keys[j]
		if a.weight() != b.weight() {
			return a.weight() < b.weight()
		}
		return fmt.Sprint(a.color) < fmt.Sprint(b.color)
	})
	for _, e := range keys {
		ls := borderStyles[e.style]
		perTile := map[*tileCv]*bdf.Path{}
		var order []*tileCv
		for _, sg := range groups[e] {
			pad := ls.w * 2
			for _, t := range s.tilesIn(minf(sg.x0, sg.x1)-pad, minf(sg.y0, sg.y1)-pad, maxf(sg.x0, sg.x1)+pad, maxf(sg.y0, sg.y1)+pad) {
				p := perTile[t]
				if p == nil {
					p = &bdf.Path{}
					perTile[t] = p
					order = append(order, t)
				}
				if ls.double {
					// two hairlines a pixel apart on each side of the grid line
					dx, dy := 0.0, pxPt
					if sg.x0 == sg.x1 {
						dx, dy = pxPt, 0
					}
					for _, o := range []float64{-1, 1} {
						p.MoveTo(f32(sg.x0+o*dx-t.ox), f32(sg.y0+o*dy-t.oy)).LineTo(f32(sg.x1+o*dx-t.ox), f32(sg.y1+o*dy-t.oy))
					}
					continue
				}
				p.MoveTo(f32(sg.x0-t.ox), f32(sg.y0-t.oy)).LineTo(f32(sg.x1-t.ox), f32(sg.y1-t.oy))
			}
		}
		for _, t := range order {
			o := t.cv.Obj
			// solid lines are capped square so that they meet at the corners
			var capStyle byte = bdf.CapSquare
			if ls.dash != nil {
				capStyle = bdf.CapButt
			}
			o.Line(f32(ls.w), capStyle, bdf.JoinMiter, 10)
			if ls.dash != nil {
				d := make([]float32, len(ls.dash))
				for i, v := range ls.dash {
					d[i] = f32(v)
				}
				o.Dash(d, 0)
			} else if t.st.dashed {
				o.Dash(nil, 0)
			}
			t.st.dashed = ls.dash != nil
			o.StrokeColor(e.color.bdf())
			o.StrokePath(o.AddPath(perTile[t]))
			t.cv.Drawn = true
		}
	}
	// drawings (children) inherit the line state and expect no dashes
	for _, t := range s.tiles {
		if t.st.dashed {
			t.cv.Obj.Dash(nil, 0)
			t.st.dashed = false
		}
	}
}

func minf(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxf(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
