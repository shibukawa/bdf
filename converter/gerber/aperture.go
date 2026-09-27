package gerber

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// prim is a part of an aperture's shape: its contours (outer ones
// counterclockwise, holes clockwise), added to the shape (on) or taken out
// of what the parts before it made (off: macro primitives of exposure 0).
type prim struct {
	on       bool
	contours []contour
}

// shape is the image an aperture flashes, about the flash point, in mm.
type shape struct {
	prims []prim
	// simple: every part is on, so the shape is the union of the contours
	// (filled with the nonzero rule)
	simple bool
	bounds box
	// round: a solid circle of diameter d, which draws round-capped lines
	round bool
	d     float64
	// block is the image of a block aperture (AB), flashed by copying it
	block *image
}

func newShape(prims []prim) *shape {
	s := &shape{prims: prims, simple: true}
	for _, p := range prims {
		if !p.on {
			s.simple = false
			continue
		}
		for i := range p.contours {
			s.bounds = s.bounds.union(p.contours[i].bounds())
		}
	}
	return s
}

// empty reports whether the shape draws nothing.
func (s *shape) empty() bool { return s.block == nil && !s.bounds.ok }

// transform returns the shape transformed by the linear transform m.
func (s *shape) transform(m affine) *shape {
	if m == identity {
		return s
	}
	flip := m.det() < 0
	prims := make([]prim, len(s.prims))
	for i, p := range s.prims {
		cs := make([]contour, len(p.contours))
		for j := range p.contours {
			cs[j] = p.contours[j].transform(m)
			if flip {
				// keep outer contours counterclockwise
				cs[j] = cs[j].reverse()
			}
		}
		prims[i] = prim{on: p.on, contours: cs}
	}
	t := newShape(prims)
	if s.round {
		t.round, t.d = true, s.d*m.scale()
	}
	return t
}

// outline returns the points of the outlines of the parts that are on:
// what a draw with the aperture sweeps (by their convex hull).
func (s *shape) outline() []vec {
	var pts []vec
	for _, p := range s.prims {
		if p.on {
			for i := range p.contours {
				pts = p.contours[i].flatten(pts, math.Pi/16)
			}
		}
	}
	return pts
}

// aperture is an aperture defined by AD (or AB): its shape in mm.
type aperture struct {
	shape *shape
	// transformed caches the shape under the aperture transformations in
	// effect (LM, LR, LS) and the image transformation
	transformed map[affine]*shape
}

func (a *aperture) at(m affine) *shape {
	if m == identity {
		return a.shape
	}
	if s, ok := a.transformed[m]; ok {
		return s
	}
	if a.transformed == nil {
		a.transformed = map[affine]*shape{}
	}
	var s *shape
	if a.shape.block != nil {
		b := a.shape.block.transformed(m, nil)
		s = &shape{block: b, bounds: b.bounds}
	} else {
		s = a.shape.transform(m)
	}
	a.transformed[m] = s
	return s
}

// holeContour returns the hole of a standard aperture (the last one or
// two parameters: a round hole, or the rectangular hole of older files),
// clockwise.
func holeContour(params []float64) (contour, bool) {
	switch len(params) {
	case 1:
		if params[0] > 0 {
			c := circle(vec{}, params[0]/2)
			return c.reverse(), true
		}
	case 2:
		if params[0] > 0 && params[1] > 0 {
			c := rect(vec{}, params[0], params[1])
			return c.reverse(), true
		}
	}
	return contour{}, false
}

// standardShape returns the shape of a standard aperture (C, R, O, P) with
// its parameters in mm.
func standardShape(kind string, p []float64) (*shape, error) {
	need := map[string]int{"C": 1, "R": 2, "O": 2, "P": 2}[kind]
	if len(p) < need {
		return nil, fmt.Errorf("aperture %s needs %d parameters", kind, need)
	}
	for _, v := range p {
		if !finite(v) || math.Abs(v) > 1e6 {
			return nil, fmt.Errorf("aperture %s has a parameter out of range", kind)
		}
	}
	var outer contour
	var hole []float64
	switch kind {
	case "C":
		d := math.Abs(p[0])
		hole = p[1:]
		if d == 0 {
			s := newShape(nil)
			s.round = true
			return s, nil
		}
		outer = circle(vec{}, d/2)
	case "R":
		w, h := math.Abs(p[0]), math.Abs(p[1])
		hole = p[2:]
		if w == 0 || h == 0 {
			return newShape(nil), nil
		}
		outer = rect(vec{}, w, h)
	case "O":
		w, h := math.Abs(p[0]), math.Abs(p[1])
		hole = p[2:]
		if w == 0 || h == 0 {
			return newShape(nil), nil
		}
		outer = obround(w, h)
	case "P":
		d := math.Abs(p[0])
		n := int(p[1])
		rot := 0.0
		if len(p) > 2 {
			rot = p[2]
			hole = p[3:]
		}
		if n < 3 || n > 12 {
			return nil, fmt.Errorf("polygon aperture with %d vertices", n)
		}
		if d == 0 {
			return newShape(nil), nil
		}
		outer = regular(vec{}, d, n, rot)
	}
	cs := []contour{outer}
	if h, ok := holeContour(hole); ok {
		cs = append(cs, h)
	}
	s := newShape([]prim{{on: true, contours: cs}})
	if kind == "C" {
		s.round, s.d = true, math.Abs(p[0])
	}
	return s, nil
}

// obround returns a rectangle w × h whose short sides are half circles,
// centred on the origin, counterclockwise.
func obround(w, h float64) contour {
	if w == h {
		return circle(vec{}, w/2)
	}
	if w > h {
		r, x := h/2, w/2-h/2
		c := contour{start: vec{-x, -r}}
		c.lineTo(vec{x, -r})
		c.arcTo(vec{x, r}, vec{x, 0}, math.Pi)
		c.lineTo(vec{-x, r})
		c.arcTo(vec{-x, -r}, vec{-x, 0}, math.Pi)
		return c
	}
	r, y := w/2, h/2-w/2
	c := contour{start: vec{r, -y}}
	c.lineTo(vec{r, y})
	c.arcTo(vec{-r, y}, vec{0, y}, math.Pi)
	c.lineTo(vec{-r, -y})
	c.arcTo(vec{r, -y}, vec{0, -y}, math.Pi)
	return c
}

// regular returns a regular polygon of n vertices on a circle of diameter
// d about c, the first at rot degrees, counterclockwise.
func regular(c vec, d float64, n int, rot float64) contour {
	pts := make([]vec, n)
	for i := range n {
		a := (rot + 360*float64(i)/float64(n)) * math.Pi / 180
		pts[i] = vec{c.X + d/2*math.Cos(a), c.Y + d/2*math.Sin(a)}
	}
	return polygon(pts...)
}

// parseParams parses the parameters of an AD command: numbers separated
// by X.
func parseParams(s string) ([]float64, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, "X")
	out := make([]float64, 0, len(parts))
	for _, f := range parts {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, fmt.Errorf("aperture parameter %q", f)
		}
		out = append(out, v)
	}
	return out, nil
}
