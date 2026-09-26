package drawingml

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Preset and custom geometry (DrawingML shape definitions): guides are
// evaluated with the shape size and adjust values, then the path list is
// turned into BDF paths in the shape's local coordinates (points).

//go:embed presets.xml.gz
var presetsGz []byte

var (
	presetsOnce sync.Once
	presets     map[string]*shapeDef
)

type guide struct{ name, fmla string }

type shapeDef struct {
	av    []guide
	gd    []guide
	rect  [4]string // l t r b guide names; empty = whole box
	paths []*ooxml.Node
}

func loadPresets() map[string]*shapeDef {
	presetsOnce.Do(func() {
		presets = map[string]*shapeDef{}
		zr, err := gzip.NewReader(bytes.NewReader(presetsGz))
		if err != nil {
			return
		}
		data, err := io.ReadAll(zr)
		if err != nil {
			return
		}
		root, err := ooxml.Parse(data)
		if err != nil {
			return
		}
		for _, s := range root.Kids {
			presets[s.Name] = parseShapeDef(s)
		}
	})
	return presets
}

// parseShapeDef reads avLst/gdLst/rect/pathLst of a preset or a:custGeom.
func parseShapeDef(n *ooxml.Node) *shapeDef {
	d := &shapeDef{}
	for _, g := range n.Child("avLst").Children("gd") {
		d.av = append(d.av, guide{g.AttrStr("name", ""), g.AttrStr("fmla", "")})
	}
	for _, g := range n.Child("gdLst").Children("gd") {
		d.gd = append(d.gd, guide{g.AttrStr("name", ""), g.AttrStr("fmla", "")})
	}
	if r := n.Child("rect"); r != nil {
		d.rect = [4]string{r.AttrStr("l", "l"), r.AttrStr("t", "t"), r.AttrStr("r", "r"), r.AttrStr("b", "b")}
	}
	d.paths = n.Child("pathLst").Children("path")
	return d
}

// guideEnv evaluates guide formulas. Angles are in 60000ths of a degree.
type guideEnv struct {
	vals map[string]float64
}

func newGuideEnv(w, h float64) *guideEnv {
	ss, ls := math.Min(w, h), math.Max(w, h)
	e := &guideEnv{vals: map[string]float64{
		"w": w, "h": h, "l": 0, "t": 0, "r": w, "b": h, "hc": w / 2, "vc": h / 2,
		"ls": ls, "ss": ss,
		"cd2": 10800000, "cd4": 5400000, "cd8": 2700000, "3cd4": 16200000, "3cd8": 8100000, "5cd8": 13500000, "7cd8": 18900000,
	}}
	for _, d := range []int{2, 3, 4, 5, 6, 8, 10, 12, 16, 32} {
		s := strconv.Itoa(d)
		e.vals["wd"+s] = w / float64(d)
		e.vals["hd"+s] = h / float64(d)
		e.vals["ssd"+s] = ss / float64(d)
	}
	return e
}

func (e *guideEnv) get(tok string) float64 {
	if v, ok := e.vals[tok]; ok {
		return v
	}
	if f, err := strconv.ParseFloat(tok, 64); err == nil {
		return f
	}
	return 0
}

const angUnit = math.Pi / 180 / 60000 // radians per ST_Angle unit

func (e *guideEnv) eval(fmla string) float64 {
	f := strings.Fields(fmla)
	if len(f) == 0 {
		return 0
	}
	arg := func(i int) float64 {
		if i < len(f) {
			return e.get(f[i])
		}
		return 0
	}
	x, y, z := arg(1), arg(2), arg(3)
	switch f[0] {
	case "val":
		return x
	case "*/":
		if z == 0 {
			return 0
		}
		return x * y / z
	case "+-":
		return x + y - z
	case "+/":
		if z == 0 {
			return 0
		}
		return (x + y) / z
	case "?:":
		if x > 0 {
			return y
		}
		return z
	case "abs":
		return math.Abs(x)
	case "at2":
		return math.Atan2(y, x) / angUnit
	case "cat2":
		return x * math.Cos(math.Atan2(z, y))
	case "sat2":
		return x * math.Sin(math.Atan2(z, y))
	case "cos":
		return x * math.Cos(y*angUnit)
	case "sin":
		return x * math.Sin(y*angUnit)
	case "tan":
		return x * math.Tan(y*angUnit)
	case "max":
		return math.Max(x, y)
	case "min":
		return math.Min(x, y)
	case "mod":
		return math.Sqrt(x*x + y*y + z*z)
	case "pin":
		if y < x {
			return x
		}
		if y > z {
			return z
		}
		return y
	case "sqrt":
		if x < 0 {
			return 0
		}
		return math.Sqrt(x)
	}
	return 0
}

// geomPath is one evaluated sub-path of a shape.
type geomPath struct {
	path   *bdf.Path
	fill   string // "none", "norm", "lighten", "lightenLess", "darken", "darkenLess"
	stroke bool
	closed bool
	// end tangents for arrowheads: start point and the direction into the
	// path, end point and the direction out of it
	start, startDir, end, endDir [2]float64
	empty                        bool
}

// geometry is an evaluated shape outline.
type geometry struct {
	paths    []geomPath
	textRect [4]float64 // l t r b
}

// shapeGeometry evaluates spPr's preset or custom geometry at size w×h.
// A shape with neither is a rectangle.
func shapeGeometry(spPr *ooxml.Node, w, h float64) (*geometry, string) {
	if pg := spPr.Child("prstGeom"); pg != nil {
		name := pg.AttrStr("prst", "rect")
		def := loadPresets()[name]
		if def == nil {
			return rectGeometry(w, h), name
		}
		return evalGeometry(def, pg.Child("avLst"), w, h), ""
	}
	if cg := spPr.Child("custGeom"); cg != nil {
		return evalGeometry(parseShapeDef(cg), nil, w, h), ""
	}
	return rectGeometry(w, h), ""
}

func rectGeometry(w, h float64) *geometry {
	p := &bdf.Path{}
	p.Rect(0, 0, float32(w), float32(h))
	return &geometry{
		paths:    []geomPath{{path: p, fill: "norm", stroke: true, closed: true}},
		textRect: [4]float64{0, 0, w, h},
	}
}

func evalGeometry(def *shapeDef, override *ooxml.Node, w, h float64) *geometry {
	e := newGuideEnv(w, h)
	for _, g := range def.av {
		e.vals[g.name] = e.eval(g.fmla)
	}
	for _, g := range override.Children("gd") {
		e.vals[g.AttrStr("name", "")] = e.eval(g.AttrStr("fmla", ""))
	}
	for _, g := range def.gd {
		e.vals[g.name] = e.eval(g.fmla)
	}
	geo := &geometry{textRect: [4]float64{0, 0, w, h}}
	if def.rect[0] != "" {
		geo.textRect = [4]float64{e.get(def.rect[0]), e.get(def.rect[1]), e.get(def.rect[2]), e.get(def.rect[3])}
	}
	for _, pn := range def.paths {
		geo.paths = append(geo.paths, buildPath(e, pn, w, h))
	}
	return geo
}

// buildPath converts one a:path. Its coordinates are in its own w×h space
// when those attributes are set, scaled to the shape size.
func buildPath(e *guideEnv, pn *ooxml.Node, w, h float64) geomPath {
	sx, sy := 1.0, 1.0
	if pw := pn.AttrFloat("w", 0); pw > 0 {
		sx = w / pw
	}
	if ph := pn.AttrFloat("h", 0); ph > 0 {
		sy = h / ph
	}
	gp := geomPath{path: &bdf.Path{}, fill: pn.AttrStr("fill", "norm"), stroke: pn.AttrBool("stroke", true), empty: true}
	var cx, cy float64       // current point (shape space)
	var sx0, sy0 float64     // sub-path start
	var prevX, prevY float64 // previous point, for the end tangent
	started := false
	pt := func(n *ooxml.Node) (float64, float64) {
		return e.get(n.AttrStr("x", "0")) * sx, e.get(n.AttrStr("y", "0")) * sy
	}
	setStart := func(x, y, dx, dy float64) {
		if !started {
			gp.start = [2]float64{x, y}
			gp.startDir = [2]float64{dx, dy}
			started = true
		}
	}
	setEnd := func(x, y, dx, dy float64) {
		gp.end = [2]float64{x, y}
		gp.endDir = [2]float64{dx, dy}
	}
	for _, c := range pn.Kids {
		switch c.Name {
		case "moveTo":
			pts := c.Children("pt")
			if len(pts) == 0 {
				continue
			}
			cx, cy = pt(pts[0])
			sx0, sy0 = cx, cy
			gp.path.MoveTo(float32(cx), float32(cy))
		case "lnTo":
			pts := c.Children("pt")
			if len(pts) == 0 {
				continue
			}
			x, y := pt(pts[0])
			setStart(cx, cy, x-cx, y-cy)
			gp.path.LineTo(float32(x), float32(y))
			prevX, prevY = cx, cy
			cx, cy = x, y
			setEnd(x, y, x-prevX, y-prevY)
			gp.empty = false
		case "quadBezTo":
			pts := c.Children("pt")
			if len(pts) < 2 {
				continue
			}
			x1, y1 := pt(pts[0])
			x, y := pt(pts[1])
			setStart(cx, cy, x1-cx, y1-cy)
			gp.path.QuadTo(float32(x1), float32(y1), float32(x), float32(y))
			cx, cy = x, y
			setEnd(x, y, x-x1, y-y1)
			gp.empty = false
		case "cubicBezTo":
			pts := c.Children("pt")
			if len(pts) < 3 {
				continue
			}
			x1, y1 := pt(pts[0])
			x2, y2 := pt(pts[1])
			x, y := pt(pts[2])
			setStart(cx, cy, x1-cx, y1-cy)
			gp.path.CubicTo(float32(x1), float32(y1), float32(x2), float32(y2), float32(x), float32(y))
			cx, cy = x, y
			setEnd(x, y, x-x2, y-y2)
			gp.empty = false
		case "arcTo":
			wR := e.get(c.AttrStr("wR", "0")) * sx
			hR := e.get(c.AttrStr("hR", "0")) * sy
			st := e.get(c.AttrStr("stAng", "0")) * angUnit
			sw := e.get(c.AttrStr("swAng", "0")) * angUnit
			x, y, a0, a1, ok := arcEnd(cx, cy, wR, hR, st, sw)
			if !ok {
				continue
			}
			ex, ey := cx-wR*math.Cos(a0), cy-hR*math.Sin(a0) // center
			dir := 1.0
			if sw < 0 {
				dir = -1
			}
			setStart(cx, cy, dir*-wR*math.Sin(a0), dir*hR*math.Cos(a0))
			gp.path.Ellipse(float32(ex), float32(ey), float32(wR), float32(hR), 0, float32(a0), float32(a1), sw < 0)
			cx, cy = x, y
			setEnd(x, y, dir*-wR*math.Sin(a1), dir*hR*math.Cos(a1))
			gp.empty = false
		case "close":
			gp.path.Close()
			gp.closed = true
			cx, cy = sx0, sy0
		}
	}
	return gp
}

// arcEnd converts an arcTo from the current point. stAng and swAng are
// visual angles (the direction from the center of the ellipse), which the
// preset formulas use; Canvas ellipses take parametric angles, so both ends
// are converted and the sweep keeps its winding.
func arcEnd(cx, cy, wR, hR, st, sw float64) (x, y, a0, a1 float64, ok bool) {
	if wR <= 0 || hR <= 0 {
		return 0, 0, 0, 0, false
	}
	param := func(v float64) float64 { return math.Atan2(wR*math.Sin(v), hR*math.Cos(v)) }
	a0 = param(st)
	// param(v) - v stays within ±π/2, so the parametric sweep is sw plus the
	// wrapped difference of the end angles.
	d := param(st+sw) - (a0 + sw)
	d -= 2 * math.Pi * math.Round(d/(2*math.Pi))
	a1 = a0 + sw + d
	ecx, ecy := cx-wR*math.Cos(a0), cy-hR*math.Sin(a0)
	return ecx + wR*math.Cos(a1), ecy + hR*math.Sin(a1), a0, a1, true
}
