package psbdf

import (
	"fmt"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
)

// shading is a gradient or a pattern of a document, as PostScript draws
// it: linear and radial gradients of opaque stops are smooth shadings, the
// rest is drawn as pictures.
type shading struct {
	ok bool
	// solid is set for a gradient of one color
	solid bool
	color bdf.Color
	// opaque is set for a gradient that PostScript draws: linear or
	// radial, with stops that are all opaque
	opaque bool
	kind   byte
	coords []float64
	stops  []stop
	// dict is the shading dictionary, once written
	dict string
}

// stop is a color stop of a gradient at an offset from 0 to 1.
type stop struct {
	off     float64
	r, g, b float64
	a       float64
}

// stopsOf returns the stops of a gradient as Canvas uses them: in the
// order of their offsets, from 0 to 1.
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

// shading returns how a gradient or a pattern is drawn.
func (ex *exporter) shading(p *bdf.Paint) *shading {
	if ex.shadings == nil {
		ex.shadings = map[*bdf.Paint]*shading{}
	}
	if sh, ok := ex.shadings[p]; ok {
		return sh
	}
	sh := &shading{kind: p.Kind}
	ex.shadings[p] = sh
	if p.Kind == bdf.PaintPattern {
		// an image pattern is a picture
		sh.ok = true
		return sh
	}
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
	degenerate := p.Kind == bdf.PaintLinear && c[0] == c[2] && c[1] == c[3] ||
		p.Kind == bdf.PaintRadial && c[0] == c[3] && c[1] == c[4] && c[2] == c[5]
	if degenerate {
		// Canvas draws nothing with a gradient between a point and itself
		return sh
	}
	sh.ok = true
	same, opaque := true, true
	for _, s := range stops {
		same = same && s == stop{s.off, stops[0].r, stops[0].g, stops[0].b, stops[0].a}
		opaque = opaque && s.a == 1
	}
	if same {
		s := stops[0]
		sh.solid, sh.color = true, bdf.RGBA(uint8(math.Round(s.r*255)), uint8(math.Round(s.g*255)), uint8(math.Round(s.b*255)), uint8(math.Round(s.a*255)))
		return sh
	}
	if p.Kind == bdf.PaintRadial {
		c[2], c[5] = math.Max(c[2], 0), math.Max(c[5], 0)
	}
	// a conic gradient has no shading of its own: it is a picture
	sh.opaque = opaque && p.Kind != bdf.PaintConic
	sh.coords, sh.stops = c, stops
	return sh
}

// shadingDict writes the shading dictionary of a gradient that is drawn
// as one.
func (ex *exporter) shadingDict(sh *shading) string {
	if sh.dict != "" {
		return sh.dict
	}
	cs := "/DeviceRGB"
	if ex.opts.Gray {
		cs = "/DeviceGray"
	}
	typ := 2
	if sh.kind == bdf.PaintRadial {
		typ = 3
	}
	sh.dict = fmt.Sprintf("<< /ShadingType %d /ColorSpace %s /Coords [%s] /Function %s /Extend [true true] >>", typ, cs, nums(sh.coords...), ex.function(sh.stops))
	return sh.dict
}

// function writes the function of the colors of stops: from 0 to 1, with
// the first and the last color outside the stops.
func (ex *exporter) function(stops []stop) string {
	val := func(s stop) string {
		if ex.opts.Gray {
			return fine(luma(s.r, s.g, s.b))
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
		return seg(stops[0], stops[1])
	}
	var fns, bounds, encode strings.Builder
	for i := 0; i+1 < len(stops); i++ {
		fns.WriteString("\n" + seg(stops[i], stops[i+1]))
		encode.WriteString("0 1 ")
		if i > 0 {
			bounds.WriteString(fine(stops[i].off) + " ")
		}
	}
	return fmt.Sprintf("<< /FunctionType 3 /Domain [0 1] /Functions [%s\n] /Bounds [%s] /Encode [%s] >>", fns.String(), bounds.String(), encode.String())
}
