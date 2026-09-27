package gerber

import (
	"math"

	"github.com/shibukawa/bdf"
)

// mapping places the plane of the files on pages: k pt per mm, the point
// origin (mm) at the top left of the page, y down.
type mapping struct {
	k      float64
	origin vec
}

func (m mapping) pt(v vec) (float32, float32) {
	return f32((v.X - m.origin.X) * m.k), f32((m.origin.Y - v.Y) * m.k)
}

// rect returns a box in page coordinates.
func (m mapping) rect(b box) bdf.Rect {
	x0, y1 := m.pt(b.Min)
	x1, y0 := m.pt(b.Max)
	return bdf.Rect{X: x0, Y: y0, W: x1 - x0, H: y1 - y0}
}

func f32(v float64) float32 {
	if !finite(v) {
		return 0
	}
	return float32(v)
}

// hairline is the width in pt of lines drawn with apertures of no size.
const hairline = 0.35

// pathWriter appends contours to a path in page coordinates, with the origin
// of the file's plane at the point (0, 0) when rel is set (the paths of
// flashes, drawn at their flash points).
type pathWriter struct {
	m   mapping
	rel bool
	p   *bdf.Path
}

func (w *pathWriter) pt(v vec) (float32, float32) {
	if w.rel {
		return f32(v.X * w.m.k), f32(-v.Y * w.m.k)
	}
	return w.m.pt(v)
}

// contour adds a contour; closed closes it.
func (w *pathWriter) contour(c *contour, closed bool) {
	x, y := w.pt(c.start)
	w.p.MoveTo(x, y)
	cur := c.start
	for _, s := range c.segs {
		if s.arc {
			a0, _, r := arcAngles(cur, s)
			cx, cy := w.pt(s.c)
			// y down: angles turn the other way
			b0 := -a0
			b1 := b0 - s.sweep
			if math.Abs(s.sweep) >= 2*math.Pi-1e-12 {
				b1 = b0 - math.Copysign(2*math.Pi, s.sweep)
			}
			rr := f32(r * w.m.k)
			w.p.Ellipse(cx, cy, rr, rr, 0, f32(b0), f32(b1), s.sweep > 0)
			if e := s.c.add(vec{r * math.Cos(a0+s.sweep), r * math.Sin(a0+s.sweep)}); e.sub(s.to).len() > 1e-6 {
				// an end off the circle: on to it
				x, y := w.pt(s.to)
				w.p.LineTo(x, y)
			}
		} else {
			x, y := w.pt(s.to)
			w.p.LineTo(x, y)
		}
		cur = s.to
	}
	if closed {
		w.p.Close()
	}
}

// encoder writes the images of layers as objects of a document.
type encoder struct {
	doc *bdf.Document
	m   mapping
}

// layerObject is an image written as an object: drawn with the fill and
// stroke colors of whoever uses it (spec §8), so that one object serves
// every color a layer is shown in. Clear objects are drawn with
// destination-out, so an image with any (clear) must be used inside a
// group of its own.
type layerObject struct {
	hash  bdf.Hash
	bbox  bdf.Rect
	clear bool
}

// object writes an image.
func (e *encoder) object(im *image) layerObject {
	o := bdf.NewObject()
	bb := e.m.rect(im.bounds.grow(0.01))
	o.SetBBox(bb.X, bb.Y, bb.W, bb.H)
	glyphs := map[*shape]bdf.PathRef{}
	children := map[*shape]bdf.ObjRef{}
	for _, r := range im.runs {
		if r.clear {
			o.Save()
			o.Blend(bdf.BlendDestinationOut)
		}
		if len(r.fills) > 0 {
			w := &pathWriter{m: e.m, p: &bdf.Path{}}
			for i := range r.fills {
				w.contour(&r.fills[i], true)
			}
			o.FillPath(o.AddPath(w.p), bdf.NonZero)
		}
		for _, st := range r.strokes {
			w := &pathWriter{m: e.m, p: &bdf.Path{}}
			for i := range st.chains {
				w.contour(&st.chains[i], false)
			}
			width := f32(st.width * e.m.k)
			if st.width <= 0 {
				width = hairline
			}
			o.Line(width, bdf.CapRound, bdf.JoinRound, 10)
			o.StrokePath(o.AddPath(w.p))
		}
		var run []bdf.Glyph
		flush := func() {
			if len(run) > 0 {
				o.FillPathRun(bdf.NonZero, run)
				run = run[:0]
			}
		}
		for _, f := range r.flashes {
			x, y := e.m.pt(f.at)
			if f.s.simple {
				ref, ok := glyphs[f.s]
				if !ok {
					w := &pathWriter{m: e.m, rel: true, p: &bdf.Path{}}
					for _, p := range f.s.prims {
						for i := range p.contours {
							w.contour(&p.contours[i], true)
						}
					}
					ref = o.AddPath(w.p)
					glyphs[f.s] = ref
				}
				run = append(run, bdf.Glyph{Path: ref, X: x, Y: y})
				continue
			}
			// a macro that takes parts away: drawn in a group of its own
			flush()
			ref, ok := children[f.s]
			if !ok {
				ref = e.child(o, f.s)
				children[f.s] = ref
			}
			b := e.m.rect(f.s.bounds.offset(f.at).grow(0.01))
			blend := bdf.BlendSourceOver
			if r.clear {
				blend = bdf.BlendDestinationOut
			}
			o.GroupBegin(1, blend, b.X, b.Y, b.W, b.H)
			o.UseAt(ref, x, y)
			o.GroupEnd()
		}
		flush()
		if r.clear {
			o.Restore()
		}
	}
	h, _ := e.doc.AddObject(o)
	return layerObject{hash: h, bbox: bb, clear: im.hasClear()}
}

// child writes the object of a shape whose parts take away from the parts
// before them, about the flash point, and adds it to o.
func (e *encoder) child(o *bdf.Object, s *shape) bdf.ObjRef {
	c := bdf.NewObject()
	b := s.bounds.grow(0.01)
	x0, y0 := f32(b.Min.X*e.m.k), f32(-b.Max.Y*e.m.k)
	w, h := f32(b.w()*e.m.k), f32(b.h()*e.m.k)
	c.SetBBox(x0, y0, w, h)
	for _, p := range s.prims {
		pw := &pathWriter{m: e.m, rel: true, p: &bdf.Path{}}
		for i := range p.contours {
			pw.contour(&p.contours[i], true)
		}
		if p.on {
			c.FillPath(c.AddPath(pw.p), bdf.NonZero)
			continue
		}
		c.Save()
		c.Blend(bdf.BlendDestinationOut)
		c.FillPath(c.AddPath(pw.p), bdf.NonZero)
		c.Restore()
	}
	hash, bbox := e.doc.AddObject(c)
	return o.AddObject(hash, bbox)
}
