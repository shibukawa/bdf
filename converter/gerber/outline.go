package gerber

import (
	"math"
	"slices"
)

// edge is a segment of an outline with the point it starts at.
type edge struct {
	from vec
	s    seg
}

// boardShape returns the shape of the board that outline layers draw: the
// closed loops their lines make (joined end to end, whatever order and
// direction they are drawn in), and the regions they fill; filled with the
// even-odd rule, loops inside the board are cut-outs. ok is false when no
// loop closes.
func boardShape(layers []*layer) (loops []contour, ok bool) {
	var edges []edge
	width := 0.0
	for _, l := range layers {
		for _, r := range l.img.runs {
			if r.clear {
				continue
			}
			for _, st := range r.strokes {
				width = math.Max(width, st.width)
				for _, ch := range st.chains {
					p := ch.start
					for _, s := range ch.segs {
						if s.arc && math.Abs(s.sweep) >= 2*math.Pi-1e-9 {
							loops = append(loops, contour{start: p, segs: []seg{s}})
						} else if p != s.to {
							edges = append(edges, edge{p, s})
						}
						p = s.to
					}
				}
			}
			loops = append(loops, r.fills...)
		}
	}
	// gaps up to 0.02 mm (or half the line width, up to 0.1 mm) close
	tol := math.Min(math.Max(0.02, width/2), 0.1)
	loops = append(loops, joinEdges(edges, tol)...)
	if len(loops) == 0 {
		return nil, false
	}
	// loops outside the largest (drawings beside the board) are left out
	big := 0
	for i := range loops {
		if math.Abs(loops[i].area()) > math.Abs(loops[big].area()) {
			big = i
		}
	}
	bb := loops[big].bounds()
	out := []contour{loops[big]}
	for i, c := range loops {
		if i == big {
			continue
		}
		b := c.bounds()
		if b.Min.X >= bb.Min.X-tol && b.Max.X <= bb.Max.X+tol && b.Min.Y >= bb.Min.Y-tol && b.Max.Y <= bb.Max.Y+tol {
			out = append(out, c)
		}
	}
	return out, true
}

// joinEdges joins edges end to end into closed loops; chains that do not
// close are dropped.
func joinEdges(edges []edge, tol float64) []contour {
	if len(edges) == 0 {
		return nil
	}
	// a grid of the edges' ends, in cells of tol
	type key struct{ x, y int64 }
	cell := func(p vec) key { return key{int64(math.Floor(p.X / tol)), int64(math.Floor(p.Y / tol))} }
	grid := map[key][]int{} // edge index × 2 + end (0 from, 1 to)
	for i, e := range edges {
		for end, p := range []vec{e.from, e.s.to} {
			k := cell(p)
			grid[k] = append(grid[k], i*2+end)
		}
	}
	used := make([]bool, len(edges))
	// find returns an unused edge with an end within tol of p, and which end
	find := func(p vec) (int, int) {
		best, bestEnd, bestD := -1, 0, tol
		k := cell(p)
		for dx := int64(-1); dx <= 1; dx++ {
			for dy := int64(-1); dy <= 1; dy++ {
				for _, v := range grid[key{k.x + dx, k.y + dy}] {
					i, end := v/2, v%2
					if used[i] {
						continue
					}
					q := edges[i].from
					if end == 1 {
						q = edges[i].s.to
					}
					if d := q.sub(p).len(); d <= bestD && (best < 0 || d < bestD || v < best*2+bestEnd) {
						best, bestEnd, bestD = i, end, d
					}
				}
			}
		}
		return best, bestEnd
	}
	var loops []contour
	for i := range edges {
		if used[i] {
			continue
		}
		used[i] = true
		c := contour{start: edges[i].from, segs: []seg{edges[i].s}}
		for c.end().sub(c.start).len() > tol {
			j, end := find(c.end())
			if j < 0 {
				break
			}
			used[j] = true
			e := edges[j]
			if end == 0 {
				c.segs = append(c.segs, e.s)
			} else {
				// the edge drawn the other way
				r := contour{start: e.from, segs: []seg{e.s}}
				rev := r.reverse()
				c.segs = append(c.segs, rev.segs[0])
			}
		}
		if c.end().sub(c.start).len() <= tol && (len(c.segs) > 2 || slices.ContainsFunc(c.segs, func(s seg) bool { return s.arc })) {
			if c.end() != c.start {
				c.lineTo(c.start)
			}
			loops = append(loops, c)
		}
	}
	return loops
}
