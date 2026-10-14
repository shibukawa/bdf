package pdfbdf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
)

// painted is what a gradient or a pattern is drawn with.
type painted struct {
	ok bool
	// solid is set for a gradient of one color
	solid bool
	color bdf.Color
	// a pattern object, with the alpha of the whole and, for a gradient
	// whose stops have different alphas, an ExtGState with a soft mask
	name  string
	ref   ref
	alpha float64
	mask  func(toPage matrix) ref
	// an image that is drawn once, under matrix (a pattern that does not
	// repeat, for a fill)
	image  *pdfImage
	matrix matrix
}

// shading is a gradient of a document, as the shadings of its colors and
// of its alphas.
type shading struct {
	ok     bool
	solid  bool
	color  bdf.Color
	colors ref
	alpha  float64 // the alpha of all stops, when they have the same
	// alphas is the shading of the alphas of stops that have different
	// ones (0 for none), and masks the soft masks made of it
	alphas ref
	masks  map[matrix]ref
}

// patKey is a pattern object: a paint under a matrix.
type patKey struct {
	paint  *bdf.Paint
	matrix matrix
}

// paint returns what a gradient or a pattern is drawn with under m, which
// maps the space of the drawing to the default space of the content stream.
func (ex *exporter) paint(sty style, m matrix, stroke bool) painted {
	p := sty.paint
	if p.Kind == bdf.PaintPattern {
		return ex.tiling(sty, m, stroke)
	}
	sh := ex.shading(p)
	switch {
	case !sh.ok:
		return painted{}
	case sh.solid:
		return painted{ok: true, solid: true, color: sh.color}
	}
	key := patKey{p, m}
	pat, ok := ex.patterns[key]
	if !ok {
		pat = named{name: ex.newName("P"), ref: ex.pw.add(fmt.Sprintf("<< /Type /Pattern /PatternType 2 /Shading %s /Matrix [%s] >>", sh.colors, m))}
		ex.patterns[key] = pat
	}
	out := painted{ok: true, name: pat.name, ref: pat.ref, alpha: sh.alpha}
	if sh.alphas != 0 {
		// the alphas, as the luminosity of a soft mask
		out.mask = func(toPage matrix) ref {
			if r, ok := sh.masks[toPage]; ok {
				return r
			}
			in := newStream(identity, identity)
			in.op("%s cm %s sh", toPage.exact(), in.use("Shading", ex.newName("Sh"), sh.alphas))
			form := ex.pw.addStream(fmt.Sprintf(" /Type /XObject /Subtype /Form /BBox [%d %d %d %d] /Resources %s /Group << /Type /Group /S /Transparency /CS /DeviceGray >>",
				-huge, -huge, huge, huge, in.resources()), in.b.Bytes(), true)
			r := ex.pw.add("<< /Type /ExtGState /SMask << /Type /Mask /S /Luminosity /G " + form.String() + " >> >>")
			sh.masks[toPage] = r
			return r
		}
	}
	return out
}

// stop is a color stop of a gradient at an offset from 0 to 1.
type stop struct {
	off     float64
	r, g, b float64
	a       float64
}

// stops returns the stops of a gradient as Canvas uses them: in the order
// of their offsets, from 0 to 1.
func stopsOf(p *bdf.Paint) []stop {
	var out []stop
	last := 0.0
	for _, s := range p.Stops {
		off := float64(s.Offset)
		if math.IsNaN(off) {
			continue
		}
		// (a stop before the one before it is at its place)
		off = math.Min(math.Max(off, last), 1)
		last = off
		c := uint32(s.Color)
		out = append(out, stop{off, float64(c>>24) / 255, float64(c>>16&0xff) / 255, float64(c>>8&0xff) / 255, float64(c&0xff) / 255})
	}
	return out
}

// function writes the function of the colors (or of the alphas) of stops:
// from 0 to 1, with the first and the last color outside the stops.
func (ex *exporter) function(stops []stop, alpha bool) ref {
	val := func(s stop) string {
		if alpha {
			return fine(s.a)
		}
		return fine(s.r) + " " + fine(s.g) + " " + fine(s.b)
	}
	seg := func(a, b stop) string {
		return fmt.Sprintf("<< /FunctionType 2 /Domain [0 1] /C0 [%s] /C1 [%s] /N 1 >>", val(a), val(b))
	}
	// the stops from 0 to 1
	if stops[0].off > 0 {
		first := stops[0]
		first.off = 0
		stops = append([]stop{first}, stops...)
	}
	if n := len(stops); stops[n-1].off < 1 {
		last := stops[n-1]
		last.off = 1
		stops = append(stops, last)
	}
	if len(stops) == 2 {
		return ex.pw.add(seg(stops[0], stops[1]))
	}
	var fns, bounds, encode strings.Builder
	for i := 0; i+1 < len(stops); i++ {
		fns.WriteString(seg(stops[i], stops[i+1]) + " ")
		encode.WriteString("0 1 ")
		if i > 0 {
			bounds.WriteString(fine(stops[i].off) + " ")
		}
	}
	return ex.pw.add(fmt.Sprintf("<< /FunctionType 3 /Domain [0 1] /Functions [%s] /Bounds [%s] /Encode [%s] >>", fns.String(), bounds.String(), encode.String()))
}

// huge is a rectangle that holds every drawing: the place of what has no
// bounds of its own.
const huge = 1 << 15

// shading writes the shadings of a gradient.
func (ex *exporter) shading(p *bdf.Paint) *shading {
	if sh, ok := ex.shadings[p]; ok {
		return sh
	}
	sh := &shading{}
	ex.shadings[p] = sh
	stops := stopsOf(p)
	want := map[byte]int{bdf.PaintLinear: 4, bdf.PaintRadial: 6, bdf.PaintConic: 3}[p.Kind]
	if len(stops) == 0 || want == 0 || len(p.Coords) < want {
		return sh // a gradient without stops draws nothing
	}
	c := make([]float64, want)
	for i := range c {
		c[i] = float64(p.Coords[i])
		if math.IsNaN(c[i]) || math.IsInf(c[i], 0) {
			return sh
		}
	}
	sh.ok = true
	same, opaque := true, true
	for _, s := range stops {
		same = same && s == stop{s.off, stops[0].r, stops[0].g, stops[0].b, stops[0].a}
		opaque = opaque && s.a == stops[0].a
	}
	degenerate := p.Kind == bdf.PaintLinear && c[0] == c[2] && c[1] == c[3] ||
		p.Kind == bdf.PaintRadial && c[0] == c[3] && c[1] == c[4] && c[2] == c[5]
	if degenerate {
		// Canvas draws nothing with a gradient between a point and itself
		sh.ok = false
		return sh
	}
	if same {
		s := stops[0]
		sh.solid, sh.color = true, bdf.RGBA(uint8(math.Round(s.r*255)), uint8(math.Round(s.g*255)), uint8(math.Round(s.b*255)), uint8(math.Round(s.a*255)))
		return sh
	}
	sh.alpha = 1
	if opaque {
		sh.alpha = stops[0].a
	}
	write := func(alpha bool) ref {
		cs := "/DeviceRGB"
		if alpha {
			cs = "/DeviceGray"
		}
		switch p.Kind {
		case bdf.PaintLinear:
			return ex.pw.add(fmt.Sprintf("<< /ShadingType 2 /ColorSpace %s /Coords [%s] /Function %s /Extend [true true] >>", cs, nums(c...), ex.function(stops, alpha)))
		case bdf.PaintRadial:
			c[2], c[5] = math.Max(c[2], 0), math.Max(c[5], 0)
			return ex.pw.add(fmt.Sprintf("<< /ShadingType 3 /ColorSpace %s /Coords [%s] /Function %s /Extend [true true] >>", cs, nums(c...), ex.function(stops, alpha)))
		}
		return ex.conic(c, stops, alpha)
	}
	sh.colors = write(false)
	if !opaque {
		sh.alphas, sh.masks = write(true), map[matrix]ref{}
	}
	return sh
}

// conic writes a conic gradient (angle, x, y) as a mesh of triangles
// around its centre, a PDF having no conic shading: each two degrees wide,
// with the colors of the gradient at its corners.
func (ex *exporter) conic(c []float64, stops []stop, alpha bool) ref {
	const steps = 180
	at := func(t float64) stop {
		if t <= stops[0].off {
			return stops[0]
		}
		for i := 1; i < len(stops); i++ {
			if t <= stops[i].off {
				a, b := stops[i-1], stops[i]
				k := 0.0
				if b.off > a.off {
					k = (t - a.off) / (b.off - a.off)
				}
				return stop{t, a.r + (b.r-a.r)*k, a.g + (b.g-a.g)*k, a.b + (b.b-a.b)*k, a.a + (b.a-a.a)*k}
			}
		}
		return stops[len(stops)-1]
	}
	var data bytes.Buffer
	vertex := func(x, y float64, s stop) {
		data.WriteByte(0)
		coord := func(v float64) uint32 {
			return uint32(math.Round((v + huge) / (2 * huge) * math.MaxUint32))
		}
		binary.Write(&data, binary.BigEndian, coord(x-c[1]))
		binary.Write(&data, binary.BigEndian, coord(y-c[2]))
		if alpha {
			data.WriteByte(byte(math.Round(s.a * 255)))
		} else {
			data.Write([]byte{byte(math.Round(s.r * 255)), byte(math.Round(s.g * 255)), byte(math.Round(s.b * 255))})
		}
	}
	const radius = huge / 2
	for i := 0; i < steps; i++ {
		t0, t1 := float64(i)/steps, float64(i+1)/steps
		a0, a1 := c[0]+t0*2*math.Pi, c[0]+t1*2*math.Pi
		vertex(c[1], c[2], at((t0+t1)/2))
		vertex(c[1]+radius*math.Cos(a0), c[2]+radius*math.Sin(a0), at(t0))
		vertex(c[1]+radius*math.Cos(a1), c[2]+radius*math.Sin(a1), at(t1))
	}
	cs, decode := "/DeviceRGB", "0 1 0 1 0 1"
	if alpha {
		cs, decode = "/DeviceGray", "0 1"
	}
	return ex.pw.addStream(fmt.Sprintf(" /ShadingType 4 /ColorSpace %s /BitsPerCoordinate 32 /BitsPerComponent 8 /BitsPerFlag 8 /Decode [%s %s %s %s %s]",
		cs, num(c[1]-huge), num(c[1]+huge), num(c[2]-huge), num(c[2]+huge), decode), data.Bytes(), true)
}

// tiling returns what an image pattern is drawn with: a tiling pattern, or
// for a fill with a pattern that does not repeat the image itself.
func (ex *exporter) tiling(sty style, m matrix, stroke bool) painted {
	p := sty.paint
	if sty.obj == nil || int(p.Image) >= len(sty.obj.Images) {
		ex.warnf("bad image reference")
		return painted{}
	}
	pm := matrix{float64(p.Matrix[0]), float64(p.Matrix[1]), float64(p.Matrix[2]), float64(p.Matrix[3]), float64(p.Matrix[4]), float64(p.Matrix[5])}
	if !pm.finite() {
		return painted{}
	}
	pt := m.mul(pm).areaScale()
	img := ex.imageOf(sty.obj.Images[p.Image], func(w, h float64) float64 { return pt }, true)
	if img == nil || !(img.w > 0) || !(img.h > 0) {
		return painted{}
	}
	// the image fills the cell of the pattern, upright (the space of the
	// pattern has its y axis down, as that of the document)
	cell := matrix{img.w, 0, 0, -img.h, 0, img.h}
	if p.Repeat == bdf.NoRepeat && !stroke {
		return painted{ok: true, image: img, matrix: pm.mul(cell), alpha: 1}
	}
	key := patKey{p, m}
	pat, ok := ex.patterns[key]
	if !ok {
		xs, ys := img.w, img.h
		if p.Repeat == bdf.RepeatY || p.Repeat == bdf.NoRepeat {
			xs = huge
		}
		if p.Repeat == bdf.RepeatX || p.Repeat == bdf.NoRepeat {
			ys = huge
		}
		in := newStream(identity, identity)
		in.op("%s cm %s Do", cell, in.use("XObject", img.name, img.ref))
		pat = named{name: ex.newName("P"), ref: ex.pw.addStream(fmt.Sprintf(" /Type /Pattern /PatternType 1 /PaintType 1 /TilingType 1 /BBox [0 0 %s] /XStep %s /YStep %s /Matrix [%s] /Resources %s",
			nums(img.w, img.h), num(xs), num(ys), m.mul(pm), in.resources()), in.b.Bytes(), true)}
		ex.patterns[key] = pat
	}
	return painted{ok: true, name: pat.name, ref: pat.ref, alpha: 1}
}
