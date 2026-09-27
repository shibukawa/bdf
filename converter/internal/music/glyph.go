package music

import (
	"math"
	"strconv"
	"sync"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/music/smufl"
)

// glyph is a SMuFL glyph with its metrics in staff spaces (y down, as on
// the page).
type glyph struct {
	name    string
	advance float64
	// x0, y0, x1, y1 is the bounding box (y down: y0 is the top).
	x0, y0, x1, y1 float64
	anchors        map[string][2]float64 // y down
	outline        *outline
}

// outline is a glyph's path in staff spaces, y down, parsed once.
type outline struct {
	verbs []byte
	args  []float64
}

var (
	glyphMu    sync.Mutex
	glyphCache = map[string]*glyph{}
)

// sym returns the glyph named name; it panics on a name the table lacks,
// which is a programming error.
func sym(name string) *glyph {
	glyphMu.Lock()
	defer glyphMu.Unlock()
	if g, ok := glyphCache[name]; ok {
		return g
	}
	src := smufl.Lookup(name)
	if src == nil {
		panic("music: no glyph " + name)
	}
	const u = smufl.Space
	g := &glyph{name: name, advance: float64(src.Advance) / u,
		x0: src.BBox[0] / u, y0: -src.BBox[3] / u, x1: src.BBox[2] / u, y1: -src.BBox[1] / u}
	if len(src.Anchors) > 0 {
		g.anchors = map[string][2]float64{}
		for k, v := range src.Anchors {
			g.anchors[k] = [2]float64{v[0] / u, -v[1] / u}
		}
	}
	g.outline = parseOutline(src.Path, 1.0/u)
	glyphCache[name] = g
	return g
}

// width is the width of the glyph's bounding box.
func (g *glyph) width() float64 { return g.x1 - g.x0 }

// anchor returns a named anchor, or (0, 0).
func (g *glyph) anchor(name string) (x, y float64) {
	a := g.anchors[name]
	return a[0], a[1]
}

// parseOutline parses the SVG path data of the table (absolute M, L, Q,
// C, Z in font units, y up) into staff spaces, y down.
func parseOutline(d string, k float64) *outline {
	o := &outline{}
	i := 0
	num := func() float64 {
		j := i
		if j < len(d) && d[j] == '-' {
			j++ // a minus starts a number: "10-20" is two
		}
		for j < len(d) && (d[j] == '.' || d[j] >= '0' && d[j] <= '9') {
			j++
		}
		v, _ := strconv.ParseFloat(d[i:j], 64)
		i = j
		for i < len(d) && d[i] == ' ' {
			i++
		}
		return v
	}
	pt := func() {
		x := num()
		y := num()
		o.args = append(o.args, x*k, -y*k)
	}
	for i < len(d) {
		c := d[i]
		i++
		switch c {
		case 'M':
			o.verbs = append(o.verbs, bdf.VerbMove)
			pt()
		case 'L':
			o.verbs = append(o.verbs, bdf.VerbLine)
			pt()
		case 'Q':
			o.verbs = append(o.verbs, bdf.VerbQuad)
			pt()
			pt()
		case 'C':
			o.verbs = append(o.verbs, bdf.VerbCubic)
			pt()
			pt()
			pt()
		case 'Z':
			o.verbs = append(o.verbs, bdf.VerbClose)
		}
	}
	return o
}

// path returns the outline scaled by s (pt per staff space).
func (o *outline) path(s float64) *bdf.Path {
	p := &bdf.Path{Verbs: append([]byte(nil), o.verbs...), Args: make([]float32, len(o.args))}
	for i, v := range o.args {
		p.Args[i] = f32(v * s)
	}
	return p
}

func f32(v float64) float32 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return 0
	}
	return float32(v)
}
