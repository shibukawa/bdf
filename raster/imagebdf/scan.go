package imagebdf

import (
	"image"
	"math"
	"slices"

	"github.com/shibukawa/bdf"
)

// mask is the coverage of device pixels in a rectangle: 0 (outside) to 1.
// A nil slice covers the whole rectangle.
type mask struct {
	r image.Rectangle
	a []float32
}

func (m *mask) empty() bool { return m == nil || m.r.Empty() }

// bytes is the memory the coverage takes.
func (m *mask) bytes() int { return 4 * len(m.a) }

// at returns the coverage of pixel (x, y), which must lie in m.r.
func (m *mask) at(x, y int) float32 {
	if m.a == nil {
		return 1
	}
	return m.a[(y-m.r.Min.Y)*m.r.Dx()+x-m.r.Min.X]
}

// row returns the coverage of the pixels [x0, x1) of row y (all in m.r);
// nil means full coverage.
func (m *mask) row(y, x0, x1 int) []float32 {
	if m.a == nil {
		return nil
	}
	i := (y-m.r.Min.Y)*m.r.Dx() - m.r.Min.X
	return m.a[i+x0 : i+x1]
}

// intersect returns the product of two masks over their common rectangle.
func intersect(a, b *mask) *mask {
	if a.empty() || b.empty() {
		return &mask{}
	}
	r := a.r.Intersect(b.r)
	if r.Empty() {
		return &mask{}
	}
	if a.a == nil && b.a == nil {
		return &mask{r: r}
	}
	out := &mask{r: r, a: make([]float32, r.Dx()*r.Dy())}
	i := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		ra, rb := a.row(y, r.Min.X, r.Max.X), b.row(y, r.Min.X, r.Max.X)
		for x := 0; x < r.Dx(); x++ {
			v := float32(1)
			if ra != nil {
				v = ra[x]
			}
			if rb != nil {
				v *= rb[x]
			}
			out.a[i] = v
			i++
		}
	}
	return out
}

// rectMask is the exact coverage of an axis-aligned rectangle in device
// space, limited to bounds.
func rectMask(x0, y0, x1, y1 float64, bounds image.Rectangle) *mask {
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	if !(x1 > x0 && y1 > y0) || math.IsInf(x0, 0) && math.IsInf(x1, 0) {
		return &mask{}
	}
	r := image.Rect(int(math.Floor(clampF(x0))), int(math.Floor(clampF(y0))), int(math.Ceil(clampF(x1))), int(math.Ceil(clampF(y1)))).Intersect(bounds)
	if r.Empty() {
		return &mask{}
	}
	cov := func(p0, p1 float64, i int) float32 {
		lo, hi := math.Max(p0, float64(i)), math.Min(p1, float64(i+1))
		if hi <= lo {
			return 0
		}
		return float32(hi - lo)
	}
	if x0 == math.Floor(x0) && x1 == math.Floor(x1) && y0 == math.Floor(y0) && y1 == math.Floor(y1) {
		return &mask{r: r}
	}
	m := &mask{r: r, a: make([]float32, r.Dx()*r.Dy())}
	xs := make([]float32, r.Dx())
	for x := range xs {
		xs[x] = cov(x0, x1, r.Min.X+x)
	}
	i := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		cy := cov(y0, y1, y)
		for _, cx := range xs {
			m.a[i] = cx * cy
			i++
		}
	}
	return m
}

// clampF keeps coordinates in a range that converts to int safely.
func clampF(v float64) float64 {
	const lim = 1 << 24
	switch {
	case math.IsNaN(v):
		return 0
	case v < -lim:
		return -lim
	case v > lim:
		return lim
	}
	return v
}

type edge struct {
	x0, y0, y1 float64
	slope      float64 // dx/dy
	dir        int8
}

// subSamples is the number of scanlines sampled per pixel row; spans are
// covered exactly along x.
const subSamples = 16

// crossing is where an edge crosses a scanline.
type crossing struct {
	x    float64
	dir  int8
	edge int
}

// scratch is the working memory of the rasterizer, which a renderer keeps
// from one drawing to the next instead of allocating it for each.
type scratch struct {
	edges  []edge
	active []int
	xs     []crossing
	area   []float32
	cover  []float32
	// pts holds the points of the paths flattened for a drawing
	pts   []point
	polys []polyline
	// These limits apply to an entire drawing instruction, including all
	// glyphs and stroke outlines. An incomplete fill is never painted.
	limited bool
	work    int
	r       *Renderer
}

const (
	maxPathPoints = 1 << 20
	maxScanWork   = 1 << 24 // active edges visited across sampled scanlines
	pathLimit     = "a path is too complex to rasterize: it is not drawn"
)

func (sc *scratch) stop() {
	sc.limited = true
	if sc.r != nil {
		sc.r.warnf(pathLimit)
	}
}

// scratchKept is the most entries of a buffer that are kept for the next
// drawing: a path of millions of edges does not hold its memory for ever.
const scratchKept = 1 << 20

func keep[T any](s []T) []T {
	if cap(s) > scratchKept {
		return nil
	}
	return s[:0]
}

// insertionLimit is the most crossings of a scanline that are sorted by
// insertion, which takes the square of their number when the edges come in
// no order, as those of a path drawn from right to left do. More are sorted
// by pdqsort, and the edges are kept in the order of their crossings for
// the next scanline, where they cross in nearly the same. Crossings in the
// same place may come in another order than by insertion: the spans between
// them are empty, and the winding after them is the same.
const insertionLimit = 16

// rasterize returns the coverage of the polylines (each closed implicitly)
// under a fill rule, limited to bounds.
func rasterize(polys []polyline, rule byte, bounds image.Rectangle) *mask {
	return (&scratch{}).rasterize(polys, rule, bounds)
}

func (sc *scratch) rasterize(polys []polyline, rule byte, bounds image.Rectangle) *mask {
	if sc.limited {
		return &mask{}
	}
	edges := sc.edges[:0]
	defer func() { sc.edges = keep(edges) }()
	minX, minY, maxX, maxY := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	add := func(a, b point) {
		if a.y == b.y || !finite(a) || !finite(b) {
			return
		}
		e := edge{dir: 1}
		if a.y > b.y {
			a, b = b, a
			e.dir = -1
		}
		e.x0, e.y0, e.y1 = a.x, a.y, b.y
		e.slope = (b.x - a.x) / (b.y - a.y)
		if math.IsNaN(e.slope) || math.IsInf(e.slope, 0) {
			return // the ends are too far apart for a float64 to tell where it crosses
		}
		edges = append(edges, e)
		minX, maxX = math.Min(minX, math.Min(a.x, b.x)), math.Max(maxX, math.Max(a.x, b.x))
		minY, maxY = math.Min(minY, a.y), math.Max(maxY, b.y)
	}
	for _, pl := range polys {
		n := len(pl.pts)
		if n < 2 {
			continue
		}
		for i := 1; i < n; i++ {
			add(pl.pts[i-1], pl.pts[i])
			if len(edges) > maxPathPoints {
				sc.stop()
				return &mask{}
			}
		}
		add(pl.pts[n-1], pl.pts[0])
		if len(edges) > maxPathPoints {
			sc.stop()
			return &mask{}
		}
	}
	if len(edges) == 0 {
		return &mask{}
	}
	r := image.Rect(int(math.Floor(clampF(minX))), int(math.Floor(clampF(minY))), int(math.Ceil(clampF(maxX))), int(math.Ceil(clampF(maxY)))).Intersect(bounds)
	if r.Empty() {
		return &mask{}
	}
	slices.SortFunc(edges, func(a, b edge) int {
		switch {
		case a.y0 < b.y0:
			return -1
		case a.y0 > b.y0:
			return 1
		}
		return 0
	})
	w := r.Dx()
	m := &mask{r: r, a: make([]float32, w*r.Dy())}
	if cap(sc.area) < w+1 || cap(sc.cover) < w+2 {
		sc.area, sc.cover = make([]float32, w+1), make([]float32, w+2)
	}
	area, cover := sc.area[:w+1], sc.cover[:w+2]
	clear(area)
	clear(cover)
	active, xs := sc.active[:0], sc.xs[:0]
	defer func() { sc.active, sc.xs = keep(active), keep(xs) }()
	next := 0
	const wgt = 1.0 / subSamples
	left := float64(r.Min.X)
	span := func(xa, xb float64) {
		xa, xb = math.Max(xa-left, 0), math.Min(xb-left, float64(w))
		if !(xb > xa) {
			return
		}
		ia, ib := int(xa), int(xb)
		if ia == ib {
			area[ia] += float32((xb - xa) * wgt)
			return
		}
		area[ia] += float32((float64(ia+1) - xa) * wgt)
		cover[ia+1] += wgt
		cover[ib] -= wgt
		if ib < w {
			area[ib] += float32((xb - float64(ib)) * wgt)
		}
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for s := 0; s < subSamples; s++ {
			sy := float64(y) + (float64(s)+0.5)/subSamples
			for next < len(edges) && edges[next].y0 <= sy {
				active = append(active, next)
				next++
			}
			if len(active) > maxScanWork-sc.work {
				sc.stop()
				return &mask{}
			}
			sc.work += len(active)
			xs = xs[:0]
			k := 0
			for _, i := range active {
				e := &edges[i]
				if e.y1 <= sy {
					continue // done
				}
				active[k] = i
				k++
				if e.y0 <= sy {
					xs = append(xs, crossing{e.x0 + (sy-e.y0)*e.slope, e.dir, i})
				}
			}
			active = active[:k]
			if len(xs) < 2 {
				continue
			}
			if len(xs) <= insertionLimit || plain {
				// insertion sort: the crossings are few and nearly sorted
				for i := 1; i < len(xs); i++ {
					for j := i; j > 0 && xs[j].x < xs[j-1].x; j-- {
						xs[j], xs[j-1] = xs[j-1], xs[j]
					}
				}
			} else {
				slices.SortFunc(xs, func(a, b crossing) int {
					switch {
					case a.x < b.x:
						return -1
					case a.x > b.x:
						return 1
					}
					return 0
				})
				// every edge that is active crosses: in this order next
				if len(xs) == len(active) {
					for i := range xs {
						active[i] = xs[i].edge
					}
				}
			}
			wind := 0
			for i := 0; i < len(xs)-1; i++ {
				wind += int(xs[i].dir)
				inside := wind != 0
				if rule == bdf.EvenOdd {
					inside = wind&1 != 0
				}
				if inside {
					span(xs[i].x, xs[i+1].x)
				}
			}
		}
		row := m.a[(y-r.Min.Y)*w : (y-r.Min.Y+1)*w]
		var run float32
		for x := 0; x < w; x++ {
			run += cover[x]
			v := area[x] + run
			if v > 1 {
				v = 1
			} else if v < 0 {
				v = 0
			}
			row[x] = v
			area[x], cover[x] = 0, 0
		}
		cover[w], cover[w+1], area[w] = 0, 0, 0
		if next >= len(edges) && len(active) == 0 {
			break
		}
	}
	return m
}

func finite(p point) bool {
	return !math.IsNaN(p.x) && !math.IsNaN(p.y) && !math.IsInf(p.x, 0) && !math.IsInf(p.y, 0)
}
