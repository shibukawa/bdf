package gerber

import (
	"errors"
	"math"
)

// image is the image of a Gerber or drill file (or of a block within
// one): runs of objects of one polarity, in the order they are added. The
// objects of a run can be drawn in any order, since a run only adds (dark)
// or only takes away (clear).
type image struct {
	runs   []*run
	bounds box
	// count is the number of segments, contours and flashes, budgeted
	// (see maxObjects)
	count *int
}

// run is objects of one polarity.
type run struct {
	clear bool
	// strokes are lines drawn with round apertures, by width
	strokes []*stroke
	byWidth map[float64]int
	// fills are closed contours filled with the nonzero rule: regions,
	// the lines drawn with other apertures, dots
	fills []contour
	// flashes are apertures flashed
	flashes []flash
}

// stroke is chains of segments drawn with a round aperture of a width
// (with round ends and joins).
type stroke struct {
	width  float64
	chains []contour
}

// flash is an aperture flashed at a point.
type flash struct {
	s  *shape
	at vec
}

// maxObjects bounds the segments, contours and flashes of the images of
// a conversion (step and repeat and block apertures copy theirs).
const maxObjects = 5_000_000

var errTooLarge = errors.New("the image holds more than 5 million objects; the rest is left out")

func newImage(count *int) *image {
	if count == nil {
		count = new(int)
	}
	return &image{count: count}
}

// spend counts n objects against the budget.
func (im *image) spend(n int) bool {
	if *im.count+n > maxObjects {
		*im.count = maxObjects + 1
		return false
	}
	*im.count += n
	return true
}

// over reports whether the budget ran out.
func (im *image) over() bool { return *im.count > maxObjects }

// empty reports whether the image draws nothing.
func (im *image) empty() bool {
	for _, r := range im.runs {
		if len(r.strokes) > 0 || len(r.fills) > 0 || len(r.flashes) > 0 {
			return false
		}
	}
	return true
}

// hasClear reports whether the image takes anything away.
func (im *image) hasClear() bool {
	for _, r := range im.runs {
		if r.clear {
			return true
		}
	}
	return false
}

// run returns the run objects of a polarity go to.
func (im *image) run(clear bool) *run {
	if n := len(im.runs); n > 0 && im.runs[n-1].clear == clear {
		return im.runs[n-1]
	}
	r := &run{clear: clear}
	im.runs = append(im.runs, r)
	return r
}

// line adds a segment drawn from p with a round aperture of width w.
func (im *image) line(clear bool, w float64, p vec, s seg) {
	if !finite(s.to.X) || !finite(s.to.Y) || !im.spend(1) {
		return
	}
	r := im.run(clear)
	if r.byWidth == nil {
		r.byWidth = map[float64]int{}
	}
	i, ok := r.byWidth[w]
	if !ok {
		i = len(r.strokes)
		r.byWidth[w] = i
		r.strokes = append(r.strokes, &stroke{width: w})
	}
	st := r.strokes[i]
	if n := len(st.chains); n > 0 && st.chains[n-1].end() == p {
		st.chains[n-1].segs = append(st.chains[n-1].segs, s)
	} else {
		st.chains = append(st.chains, contour{start: p, segs: []seg{s}})
	}
	b := box{}.add(p).add(s.to)
	if s.arc {
		b = b.union(arcBounds(p, s))
	}
	if !r.clear {
		im.bounds = im.bounds.union(b.grow(w / 2))
	}
}

// fill adds a closed contour, counterclockwise.
func (im *image) fill(clear bool, c contour) {
	if !im.spend(1 + len(c.segs)) {
		return
	}
	r := im.run(clear)
	r.fills = append(r.fills, c)
	if !clear {
		im.bounds = im.bounds.union(c.bounds())
	}
}

// flash adds a flash of s at p.
func (im *image) flash(clear bool, s *shape, p vec) {
	if s.empty() || !finite(p.X) || !finite(p.Y) {
		return
	}
	if s.block != nil {
		// a block aperture: its objects, in their polarities (inverted by
		// a clear flash)
		im.insert(s.block, translate(p), clear)
		return
	}
	if !im.spend(1) {
		return
	}
	r := im.run(clear)
	r.flashes = append(r.flashes, flash{s: s, at: p})
	if !clear {
		im.bounds = im.bounds.union(s.bounds.offset(p))
	}
}

// insert adds the objects of another image transformed by m (a
// similarity), inverting their polarities when invert is set.
func (im *image) insert(src *image, m affine, invert bool) {
	shapes := map[*shape]*shape{}
	lin := m
	lin[4], lin[5] = 0, 0
	for _, sr := range src.runs {
		clear := sr.clear != invert
		for _, st := range sr.strokes {
			w := st.width * m.scale()
			for _, ch := range st.chains {
				c := ch.transform(m)
				p := c.start
				for _, s := range c.segs {
					im.line(clear, w, p, s)
					p = s.to
				}
			}
		}
		for _, f := range sr.fills {
			c := f.transform(m)
			if m.det() < 0 {
				c = c.reverse()
			}
			im.fill(clear, c)
		}
		for _, f := range sr.flashes {
			s, ok := shapes[f.s]
			if !ok {
				s = f.s.transform(lin)
				shapes[f.s] = s
			}
			im.flash(clear, s, m.apply(f.at))
		}
		if im.over() {
			return
		}
	}
}

// transformed returns a copy of the image transformed by m, counting its
// objects against count (nil: the image's budget).
func (im *image) transformed(m affine, count *int) *image {
	if count == nil {
		count = im.count
	}
	out := newImage(count)
	out.insert(im, m, false)
	return out
}

// dot is the contour of a dot drawn by a round aperture of width w (a draw
// of zero length).
func dot(p vec, w float64) contour { return circle(p, w/2) }

// sweep returns the outline that a flat aperture outline swept from a to b
// covers: the convex hull of the outline at both ends (exact for convex
// apertures, the only ones that should draw).
func sweep(outline []vec, a, b vec) (contour, bool) {
	pts := make([]vec, 0, 2*len(outline))
	for _, p := range outline {
		pts = append(pts, p.add(a), p.add(b))
	}
	h := hull(pts)
	if len(h) < 3 {
		return contour{}, false
	}
	return polygon(h...), true
}

// arcAngles returns the start and end angles (radians) and the radius of
// an arc from p.
func arcAngles(p vec, s seg) (a0, a1, r float64) {
	r = p.sub(s.c).len()
	a0 = math.Atan2(p.Y-s.c.Y, p.X-s.c.X)
	return a0, a0 + s.sweep, r
}
