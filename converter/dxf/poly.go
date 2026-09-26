package dxf

import (
	"math"

	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

type vertex struct {
	p            cad.Point
	sw, ew       float64 // start and end widths
	bulge        float64
	hasSW, hasEW bool
}

// lwVertices reads the vertices of an LWPOLYLINE, whether it is closed and
// its constant width.
func lwVertices(e *entity) (vs []vertex, closed bool, width float64) {
	for _, t := range e.tags {
		n := len(vs)
		switch t.code {
		case 10:
			vs = append(vs, vertex{p: cad.Point{X: t.f}})
		case 20:
			if n > 0 {
				vs[n-1].p.Y = t.f
			}
		case 40:
			if n > 0 {
				vs[n-1].sw, vs[n-1].hasSW = t.f, true
			}
		case 41:
			if n > 0 {
				vs[n-1].ew, vs[n-1].hasEW = t.f, true
			}
		case 42:
			if n > 0 {
				vs[n-1].bulge = t.f
			}
		}
	}
	width = e.num(43, 0)
	for i := range vs {
		if !vs[i].hasSW {
			vs[i].sw = width
		}
		if !vs[i].hasEW {
			vs[i].ew = width
		}
	}
	return vs, e.int(70, 0)&1 != 0, width
}

// polyline draws a 2D polyline in the coordinates that m maps into the
// drawing: stroked with the entity's pen, or with its width when all
// segments have the same one, or filled segment by segment when the
// widths vary.
func (c *converter) polyline(e *entity, x *ctx, p props, m canvas.Matrix, vs []vertex, closed bool) {
	if len(vs) == 0 {
		return
	}
	nseg := len(vs) - 1
	if closed {
		nseg = len(vs)
	}
	uniform, w0 := true, vs[0].sw
	for i := range nseg {
		if vs[i].sw != w0 || vs[i].ew != w0 {
			uniform = false
			break
		}
	}
	pen := c.pen(e, x, p, m)
	if uniform {
		path := &cad.Path{}
		path.MoveTo(vs[0].p.X, vs[0].p.Y)
		for i := range nseg {
			q := vs[(i+1)%len(vs)].p
			path.BulgeTo(q, vs[i].bulge)
		}
		if closed {
			path.Close()
		}
		if len(vs) == 1 {
			path.LineTo(vs[0].p.X, vs[0].p.Y)
		}
		if w0 > 0 {
			pen.WorldWidth = w0 * cad.Scale(m)
			pen.Cap, pen.Join = cad.CapButt, cad.JoinMiter
		}
		x.out.Stroke(path.Transform(m), pen)
		return
	}
	// widths that vary: each segment is a filled band
	for i := range nseg {
		a, b := vs[i], vs[(i+1)%len(vs)]
		band := segmentBand(a.p, b.p, a.bulge, a.sw, a.ew)
		if band == nil {
			continue
		}
		if a.sw == 0 && a.ew == 0 {
			path := (&cad.Path{}).MoveTo(a.p.X, a.p.Y).BulgeTo(b.p, a.bulge)
			x.out.Stroke(path.Transform(m), pen)
			continue
		}
		x.out.Fill(band.Transform(m), cad.Fill{Color: p.color}, false)
	}
}

// segmentBand returns the outline of a polyline segment from a to b whose
// width goes from w0 to w1.
func segmentBand(a, b cad.Point, bulge, w0, w1 float64) *cad.Path {
	if a == b {
		return nil
	}
	var left, right []cad.Point
	if bulge == 0 {
		d := b.Sub(a)
		n := cad.Point{X: -d.Y, Y: d.X}.Mul(1 / d.Len())
		left = []cad.Point{a.Add(n.Mul(w0 / 2)), b.Add(n.Mul(w1 / 2))}
		right = []cad.Point{a.Sub(n.Mul(w0 / 2)), b.Sub(n.Mul(w1 / 2))}
	} else {
		ctr, r, a0, a1 := cad.BulgeArc(a, b, bulge)
		steps := max(int(math.Ceil(math.Abs(a1-a0)/(math.Pi/32))), 2)
		for k := 0; k <= steps; k++ {
			t := float64(k) / float64(steps)
			ang := a0 + (a1-a0)*t
			w := w0 + (w1-w0)*t
			dir := cad.Point{X: math.Cos(ang), Y: math.Sin(ang)}
			left = append(left, ctr.Add(dir.Mul(r+w/2)))
			right = append(right, ctr.Add(dir.Mul(r-w/2)))
		}
	}
	pts := append([]cad.Point(nil), left...)
	for i := len(right) - 1; i >= 0; i-- {
		pts = append(pts, right[i])
	}
	return (&cad.Path{}).Polyline(pts, true)
}

// heavyPolyline draws a POLYLINE: a 2D or 3D polyline, a polyface mesh or
// a polygon mesh.
func (c *converter) heavyPolyline(e *entity, x *ctx, p props) {
	flags := e.int(70, 0)
	switch {
	case flags&64 != 0:
		c.polyface(e, x, p)
		return
	case flags&16 != 0:
		c.polygonMesh(e, x, p, flags)
		return
	}
	sw, ew := e.num(40, 0), e.num(41, 0)
	var vs []vertex
	for _, v := range e.kids {
		vf := v.int(70, 0)
		if flags&4 != 0 && vf&16 != 0 {
			continue // the frame of a spline-fit polyline
		}
		vx := vertex{p: v.pt(10), sw: v.num(40, sw), ew: v.num(41, ew)}
		if flags&8 == 0 {
			vx.bulge = v.num(42, 0)
		}
		vs = append(vs, vx)
	}
	m := x.m
	if flags&8 == 0 {
		m = m.Mul(ocs(e, e.vec3(10, [3]float64{})[2]))
	} else {
		for i := range vs {
			vs[i].sw, vs[i].ew = 0, 0
		}
	}
	c.polyline(e, x, p, m, vs, flags&1 != 0)
}

func (c *converter) polyface(e *entity, x *ctx, p props) {
	var verts []cad.Point
	path := &cad.Path{}
	for _, v := range e.kids {
		vf := v.int(70, 0)
		if vf&64 != 0 {
			verts = append(verts, v.pt(10))
			continue
		}
		if vf&128 == 0 {
			continue
		}
		idx := []int{v.int(71, 0), v.int(72, 0), v.int(73, 0), v.int(74, 0)}
		for len(idx) > 0 && idx[len(idx)-1] == 0 {
			idx = idx[:len(idx)-1]
		}
		for k := range idx {
			a, b := idx[k], idx[(k+1)%len(idx)]
			if a <= 0 {
				continue // an invisible edge
			}
			if b < 0 {
				b = -b
			}
			if a > len(verts) || b == 0 || b > len(verts) {
				continue
			}
			pa, pb := verts[a-1], verts[b-1]
			path.MoveTo(pa.X, pa.Y).LineTo(pb.X, pb.Y)
		}
	}
	x.out.Stroke(path.Transform(x.m), c.pen(e, x, p, x.m))
}

func (c *converter) polygonMesh(e *entity, x *ctx, p props, flags int) {
	mc, nc := e.int(71, 0), e.int(72, 0)
	var verts []cad.Point
	for _, v := range e.kids {
		verts = append(verts, v.pt(10))
	}
	if mc <= 0 || nc <= 0 || mc*nc > len(verts) {
		return
	}
	at := func(i, j int) cad.Point { return verts[i*nc+j] }
	path := &cad.Path{}
	for i := range mc {
		row := make([]cad.Point, 0, nc)
		for j := range nc {
			row = append(row, at(i, j))
		}
		path.Polyline(row, flags&32 != 0)
	}
	for j := range nc {
		col := make([]cad.Point, 0, mc)
		for i := range mc {
			col = append(col, at(i, j))
		}
		path.Polyline(col, flags&1 != 0)
	}
	x.out.Stroke(path.Transform(x.m), c.pen(e, x, p, x.m))
}
