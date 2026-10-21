package psbdf

import (
	"fmt"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/vector/internal/export"
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
	stops := export.StopsOf(p)
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
		same = same && s == stop{Off: s.Off, R: stops[0].R, G: stops[0].G, B: stops[0].B, A: stops[0].A}
		opaque = opaque && s.A == 1
	}
	if same {
		s := stops[0]
		sh.solid, sh.color = true, bdf.RGBA(uint8(math.Round(s.R*255)), uint8(math.Round(s.G*255)), uint8(math.Round(s.B*255)), uint8(math.Round(s.A*255)))
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
	sh.dict = fmt.Sprintf("<< /ShadingType %d /ColorSpace %s /Coords [%s] /Function %s /Extend [true true] >>", typ, cs, export.Nums(sh.coords...), ex.function(sh.stops))
	return sh.dict
}

// function writes the function of the colors of stops: from 0 to 1, with
// the first and the last color outside the stops.
// stop is a stop of a gradient (see export.Stop).
type stop = export.Stop

func (ex *exporter) function(stops []stop) string {
	val := func(s stop) string {
		if ex.opts.Gray {
			return export.Fine(luma(s.R, s.G, s.B))
		}
		return export.Fine(s.R) + " " + export.Fine(s.G) + " " + export.Fine(s.B)
	}
	seg := func(a, b stop) string {
		return fmt.Sprintf("<< /FunctionType 2 /Domain [0 1] /C0 [%s] /C1 [%s] /N 1 >>", val(a), val(b))
	}
	// the stops from 0 to 1
	if stops[0].Off > 0 {
		first := stops[0]
		first.Off = 0
		stops = append([]stop{first}, stops...)
	}
	if n := len(stops); stops[n-1].Off < 1 {
		last := stops[n-1]
		last.Off = 1
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
			bounds.WriteString(export.Fine(stops[i].Off) + " ")
		}
	}
	return fmt.Sprintf("<< /FunctionType 3 /Domain [0 1] /Functions [%s\n] /Bounds [%s] /Encode [%s] >>", fns.String(), bounds.String(), encode.String())
}
