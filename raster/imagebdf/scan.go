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
	// full is the rectangle within r whose pixels a mask with coverage
	// is known to cover whole: that of a rectangle with fractional
	// edges but for its edges (see rectMaskIn), or the product of such
	// (see intersect). Empty when not known. The draws through such a
	// mask take the fast paths of fill within it (see mask.whole).
	full image.Rectangle
}

func (m *mask) empty() bool { return m == nil || m.r.Empty() }

// whole returns the rectangle within which the mask covers its pixels
// whole: r for a mask without coverage, else full.
func (m *mask) whole() image.Rectangle {
	if m.a == nil {
		return m.r
	}
	return m.full
}

// within returns the columns of [x0, x1) of row y that the mask covers
// whole, [lo, hi) within [x0, x1): none (lo == hi) when the row is
// outside the rows it covers whole.
func (m *mask) within(y, x0, x1 int) (lo, hi int) {
	f := m.whole()
	if y < f.Min.Y || y >= f.Max.Y {
		return x0, x0
	}
	lo = min(max(x0, f.Min.X), x1)
	hi = max(min(x1, f.Max.X), lo)
	return lo, hi
}

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
func intersect(a, b *mask) *mask { return (*Renderer)(nil).intersect(a, b) }

// intersect is the function of that name with the memory and the
// goroutines of a renderer (see newMaskBuf and splitRows).
func (r *Renderer) intersect(a, b *mask) *mask {
	if a.empty() || b.empty() {
		return &mask{}
	}
	rect := a.r.Intersect(b.r)
	if rect.Empty() {
		return &mask{}
	}
	if a.a == nil && b.a == nil {
		return &mask{r: rect}
	}
	w := rect.Dx()
	threads := 1
	if r != nil {
		threads = r.threads
	}
	out := &mask{r: rect, a: r.newMaskBuf(w * rect.Dy()), full: a.whole().Intersect(b.whole())}
	splitRows(threads, rect.Min.Y, rect.Max.Y, w, func(y0, y1, _ int) {
		for y := y0; y < y1; y++ {
			ra, rb := a.row(y, rect.Min.X, rect.Max.X), b.row(y, rect.Min.X, rect.Max.X)
			row := out.a[(y-rect.Min.Y)*w : (y-rect.Min.Y+1)*w]
			switch {
			case ra == nil:
				copy(row, rb)
			case rb == nil:
				copy(row, ra)
			default:
				for x := range row {
					row[x] = ra[x] * rb[x]
				}
			}
		}
	})
	return out
}

// rectMask is the exact coverage of an axis-aligned rectangle in device
// space, limited to bounds.
func rectMask(x0, y0, x1, y1 float64, bounds image.Rectangle) *mask {
	return rectMaskIn(nil, x0, y0, x1, y1, bounds)
}

// rectMaskIn is rectMask with the memory of a renderer (see newMaskBuf).
func rectMaskIn(r *Renderer, x0, y0, x1, y1 float64, bounds image.Rectangle) *mask {
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	if !(x1 > x0 && y1 > y0) || math.IsInf(x0, 0) && math.IsInf(x1, 0) {
		return &mask{}
	}
	rect := image.Rect(int(math.Floor(clampF(x0))), int(math.Floor(clampF(y0))), int(math.Ceil(clampF(x1))), int(math.Ceil(clampF(y1)))).Intersect(bounds)
	if rect.Empty() {
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
		return &mask{r: rect}
	}
	w := rect.Dx()
	m := &mask{r: rect, a: r.newMaskBuf(w * rect.Dy())}
	// the pixels within the edges are covered whole: none when the
	// edges lie within a pixel (image.Rect would swap such corners)
	full := image.Rectangle{Min: image.Pt(int(math.Ceil(clampF(x0))), int(math.Ceil(clampF(y0)))), Max: image.Pt(int(math.Floor(clampF(x1))), int(math.Floor(clampF(y1))))}
	if !full.Empty() {
		m.full = full.Intersect(rect)
	}
	xs := make([]float32, w)
	for x := range xs {
		xs[x] = cov(x0, x1, rect.Min.X+x)
	}
	threads := 1
	if r != nil {
		threads = r.threads
	}
	splitRows(threads, rect.Min.Y, rect.Max.Y, w, func(ya, yb, _ int) {
		for y := ya; y < yb; y++ {
			cy := cov(y0, y1, y)
			row := m.a[(y-rect.Min.Y)*w : (y-rect.Min.Y+1)*w]
			for x, cx := range xs {
				row[x] = cx * cy
			}
		}
	})
	return m
}

// rectMask follows the renderer's edge policy, including rectangular clips.
func (r *Renderer) rectMask(x0, y0, x1, y1 float64, bounds image.Rectangle) *mask {
	if r.opts.Antialias != AntialiasOff {
		return rectMaskIn(r, x0, y0, x1, y1, bounds)
	}
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	if !(x1 > x0 && y1 > y0) {
		return &mask{}
	}
	// Pixel centers lie in the half-open rectangle. Full coverage is
	// implicit, so even fractional rectangles need no coverage array.
	return &mask{r: image.Rect(
		int(math.Ceil(clampF(x0)-0.5)), int(math.Ceil(clampF(y0)-0.5)),
		int(math.Ceil(clampF(x1)-0.5)), int(math.Ceil(clampF(y1)-0.5)),
	).Intersect(bounds)}
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
	edges []edge
	// band is the memory of scanning the rows of a mask, and bands that
	// of the goroutines that scan the rows of a large one in bands (see
	// rasterizeBands)
	band  scanBand
	bands []scanBand
	// pts holds the points of the paths flattened for a drawing
	pts   []point
	polys []polyline
	// These limits apply to an entire drawing instruction, including all
	// glyphs and stroke outlines. An incomplete fill is never painted.
	limited bool
	work    int
	r       *Renderer
	noAA    bool
}

// scanBand is the memory of scanning a band of rows of a mask (see scan).
type scanBand struct {
	active []int
	xs     []crossing
	spans  []int
	area   []float32
	cover  []float32
	row    []float32
	// work counts the active edges visited (see maxScanWork), limited
	// says that the band passed the limit, and whole that every row of
	// the band was covered whole
	work    int
	limited bool
	whole   bool
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

// insertionMoves is how many moves a crossing of a longer scanline may
// take to sort by insertion before the scanline is sorted by pdqsort
// instead: nearly sorted crossings take few, those in no order many.
const insertionMoves = 4

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
	w, h := r.Dx(), r.Dy()
	threads := 1
	if sc.r != nil {
		threads = sc.r.threads
	}
	if threads > 1 && !plain && h >= 2 && w*h >= parallelMinPixels {
		return sc.rasterizeBands(edges, rule, r, threads)
	}
	// the coverage is made when a row is not covered whole: a mask whose
	// rows all are covers its rectangle whole and needs none (see mask.a)
	m := &mask{r: r}
	if plain {
		m.a = make([]float32, w*h)
	}
	b := &sc.band
	b.work = sc.work
	b.scan(sc.r, edges, rule, r, r.Min.Y, r.Max.Y, m, sc.noAA)
	sc.work = b.work
	if b.limited {
		sc.stop()
		return &mask{}
	}
	return m
}

// rasterizeBands scans the rows of a large mask in bands, a goroutine each
// (see splitRows), into coverage made before: a mask whose bands were all
// covered whole gives it back and needs none.
func (sc *scratch) rasterizeBands(edges []edge, rule byte, r image.Rectangle, threads int) *mask {
	w, h := r.Dx(), r.Dy()
	m := &mask{r: r, a: sc.r.newMaskBuf(w * h)}
	for len(sc.bands) < threads {
		sc.bands = append(sc.bands, scanBand{})
	}
	bands := sc.bands[:threads]
	before := sc.work
	for i := range bands {
		bands[i].work, bands[i].limited, bands[i].whole = before, false, true
	}
	splitRows(threads, r.Min.Y, r.Max.Y, w, func(y0, y1, k int) {
		bands[k].scan(sc.r, edges, rule, r, y0, y1, m, sc.noAA)
	})
	whole, limited := true, false
	for i := range bands {
		b := &bands[i]
		sc.work += b.work - before
		limited = limited || b.limited
		whole = whole && b.whole
	}
	if limited || sc.work > maxScanWork {
		sc.r.releaseMask(m)
		sc.stop()
		return &mask{}
	}
	if whole {
		sc.r.releaseMask(m)
		return &mask{r: r}
	}
	return m
}

// scan scans the rows [y0, y1) of the mask m of r for edges sorted by
// their top: with antialiasing, subSamples scanlines a row with spans
// covered exactly along x; without, pixel centers once a row, with whole
// pixel spans. The rows go into the coverage of m when it has one (a
// band of a large mask); when it has none yet, it is made at the first
// row that is not covered whole, with the rows before covered (a mask
// scanned alone: a mask whose rows are all covered whole needs none).
func (b *scanBand) scan(rr *Renderer, edges []edge, rule byte, r image.Rectangle, y0, y1 int, m *mask, noAA bool) {
	w, h := r.Dx(), r.Dy()
	b.limited, b.whole = false, true
	if cap(b.area) < w+1 || cap(b.cover) < w+2 || cap(b.row) < w {
		b.area, b.cover, b.row = make([]float32, w+1), make([]float32, w+2), make([]float32, w)
	}
	area, cover, row := b.area[:w+1], b.cover[:w+2], b.row[:w]
	clear(area)
	clear(cover)
	active, xs, spans := b.active[:0], b.xs[:0], b.spans[:0]
	defer func() { b.active, b.xs, b.spans = keep(active), keep(xs), keep(spans) }()
	samples := subSamples
	if noAA {
		samples = 1
	}
	wgt := 1 / float32(samples)
	left := float64(r.Min.X)
	// lo and hi bound the pixels the spans of a row reach: [lo, hi) is
	// where its coverage can be other than 0 (hi may be w + 1, for a
	// span that ends at the right edge)
	lo, hi := w, 0
	// span covers [xa, xb) of the scanline, in device x
	span := func(xa, xb float64) {
		xa, xb = math.Max(xa-left, 0), math.Min(xb-left, float64(w))
		if !(xb > xa) {
			return
		}
		ia, ib := int(xa), int(xb)
		lo, hi = min(lo, ia), max(hi, ib+1)
		if ia == ib {
			area[ia] += float32((xb - xa) * float64(wgt))
			return
		}
		area[ia] += float32((float64(ia+1) - xa) * float64(wgt))
		cover[ia+1] += wgt
		cover[ib] -= wgt
		if ib < w {
			area[ib] += float32((xb - float64(ib)) * float64(wgt))
		}
	}
	next := 0
	for y := y0; y < y1; y++ {
		spans = spans[:0]
		lo, hi = w, 0
		whole := false
		for s := 0; s < samples; s++ {
			sy := float64(y) + (float64(s)+0.5)/float64(samples)
			for next < len(edges) && edges[next].y0 <= sy {
				active = append(active, next)
				next++
			}
			if len(active) > maxScanWork-b.work {
				b.limited = true
				return
			}
			b.work += len(active)
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
			sortCrossings(xs, active)
			if !noAA {
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
				continue
			}
			// the spans of the row, in order, and whether they cover it whole
			wind, covered := 0, 0
			whole = true
			for i := 0; i < len(xs)-1; i++ {
				wind += int(xs[i].dir)
				inside := wind != 0
				if rule == bdf.EvenOdd {
					inside = wind&1 != 0
				}
				if !inside {
					continue
				}
				xa := math.Max(xs[i].x, float64(r.Min.X))
				xb := math.Min(xs[i+1].x, float64(r.Max.X))
				if !(xb > xa) {
					continue
				}
				// Left-inclusive, right-exclusive at half-integer centers.
				x0 := int(math.Ceil(xa-0.5)) - r.Min.X
				x1 := int(math.Ceil(xb-0.5)) - r.Min.X
				if x1 <= x0 {
					continue
				}
				if x0 > covered {
					whole = false
				}
				covered = max(covered, x1)
				spans = append(spans, x0, x1)
			}
			whole = whole && covered >= w
		}
		x1 := min(hi, w)
		if !noAA && lo < hi {
			// the coverage is settled over [lo, hi) alone: the pixels
			// outside are 0, and so are their area and cover (the
			// running sum of cover is 0 past the spans, which add as
			// much as they take)
			var run float32
			whole = lo == 0 && hi >= w
			for x := lo; x < x1; x++ {
				run += cover[x]
				v := area[x] + run
				if v > 1 {
					v = 1
				} else if v < 0 {
					v = 0
				}
				row[x] = v
				whole = whole && v == 1
				area[x], cover[x] = 0, 0
			}
			clear(area[x1:hi])
			clear(cover[x1:hi])
		}
		if !whole {
			b.whole = false
			m.coverFull(rr, y-r.Min.Y, w, h)
		}
		if m.a != nil {
			out := m.a[(y-r.Min.Y)*w : (y-r.Min.Y+1)*w]
			if !noAA {
				// the rest of the row is 0 already: the coverage is
				// made of zeros, and each row written once
				if lo < x1 {
					copy(out[lo:x1], row[lo:x1])
				}
			} else {
				for i := 0; i < len(spans); i += 2 {
					ones := out[spans[i]:spans[i+1]]
					ones[0] = 1
					for n := 1; n < len(ones); n *= 2 {
						copy(ones[n:], ones[:n])
					}
				}
			}
		}
		if next >= len(edges) && len(active) == 0 {
			if y+1 < y1 {
				// the rows after the last edge are not covered
				b.whole = false
				m.coverFull(rr, y+1-r.Min.Y, w, h)
			}
			break
		}
	}
}

// coverFull makes the coverage of a mask of w × h pixels whose first rows
// are covered whole, if it has none yet, with the memory of a renderer.
func (m *mask) coverFull(r *Renderer, rows, w, h int) {
	if m.a != nil {
		return
	}
	m.a = r.newMaskBuf(w * h)
	ones := m.a[:rows*w]
	if len(ones) > 0 {
		ones[0] = 1
		for n := 1; n < len(ones); n *= 2 {
			copy(ones[n:], ones[:n])
		}
	}
}

// sortCrossings preserves the edge order for the next scanline when useful.
func sortCrossings(xs []crossing, active []int) {
	sorted := len(xs) <= insertionLimit || plain
	if sorted {
		// insertion sort: the crossings are few and nearly sorted
		for i := 1; i < len(xs); i++ {
			for j := i; j > 0 && xs[j].x < xs[j-1].x; j-- {
				xs[j], xs[j-1] = xs[j-1], xs[j]
			}
		}
	} else {
		// the edges are in the order of their crossings on the
		// scanline before, so insertion sorts them in a few
		// moves, unless this scanline orders them anew (the
		// first one, or edges that cross)
		moves := insertionMoves * len(xs)
		sorted = true
	insertion:
		for i := 1; i < len(xs); i++ {
			for j := i; j > 0 && xs[j].x < xs[j-1].x; j-- {
				xs[j], xs[j-1] = xs[j-1], xs[j]
				if moves--; moves < 0 {
					sorted = false
					break insertion
				}
			}
		}
	}
	if !sorted {
		slices.SortFunc(xs, func(a, b crossing) int {
			switch {
			case a.x < b.x:
				return -1
			case a.x > b.x:
				return 1
			}
			return 0
		})
	}
	// every edge that is active crosses: in this order next
	if len(xs) > insertionLimit && len(xs) == len(active) && !plain {
		for i := range xs {
			active[i] = xs[i].edge
		}
	}
}

func finite(p point) bool {
	return !math.IsNaN(p.x) && !math.IsNaN(p.y) && !math.IsInf(p.x, 0) && !math.IsInf(p.y, 0)
}
