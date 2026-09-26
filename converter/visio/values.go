package visio

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
)

// Cell values are read through the page being drawn: cells set by the
// dynamic theme ("Themed") resolve against the theme with the QuickStyle
// cells of the shape and the variant the page selects.

// val returns a top-level cell value of a shape.
func (p *pageCtx) val(s *shape, name string) (string, bool) {
	c, ok := s.get(name)
	if !ok {
		return "", false
	}
	if c.v == "Themed" {
		return p.themed(s, "", "", name)
	}
	return c.v, true
}

func (p *pageCtx) num(s *shape, name string, def float64) float64 {
	v, ok := p.val(s, name)
	if !ok || strings.TrimSpace(v) == "" {
		return def
	}
	return num(v)
}

func (p *pageCtx) flag(s *shape, name string) bool { return p.num(s, name, 0) != 0 }

// dim reads a cell that follows the shape's size (a pin, the text block,
// the picture's rectangle): an inherited value whose formula only scales
// Width and Height is evaluated with the shape's own size (see formula.go).
func (p *pageCtx) dim(s *shape, name string, def float64) float64 {
	if _, own := s.cells[name]; !own {
		if v, ok := evalSize(s.formula(name), p.num(s, "Width", 0), p.num(s, "Height", 0)); ok {
			return v
		}
	}
	return p.num(s, name, def)
}

// color reads a color cell with its transparency cell (0-1) applied.
func (p *pageCtx) color(s *shape, name, trans string) (bdf.Color, bool) {
	v, ok := p.val(s, name)
	if !ok {
		return 0, false
	}
	c, ok := p.c.d.color(v)
	if !ok {
		return 0, false
	}
	return withTrans(c, p.num(s, trans, 0)), true
}

// rowVal returns a cell of a row of a single section (Character …).
func (p *pageCtx) rowVal(s *shape, sec, key, name string) (string, bool) {
	c, ok := s.rowCell(sec, key, name)
	if !ok {
		return "", false
	}
	if c.v == "Themed" {
		return p.themed(s, sec, key, name)
	}
	return c.v, true
}

func (p *pageCtx) rowNum(s *shape, sec, key, name string, def float64) float64 {
	v, ok := p.rowVal(s, sec, key, name)
	if !ok || strings.TrimSpace(v) == "" {
		return def
	}
	return num(v)
}

// withTrans applies a transparency (0 opaque … 1 transparent).
func withTrans(c bdf.Color, t float64) bdf.Color {
	if t <= 0 {
		return c
	}
	a := float64(c&0xff) * (1 - math.Min(t, 1))
	return c&^0xff | bdf.Color(math.Round(a))
}

func hexOf(c bdf.Color) string {
	const digits = "0123456789ABCDEF"
	b := []byte{'#', 0, 0, 0, 0, 0, 0}
	for i := 0; i < 6; i++ {
		b[1+i] = digits[c>>(28-4*uint(i))&0xf]
	}
	return string(b)
}

func alphaOf(c bdf.Color) float64 { return float64(c&0xff) / 255 }

func ftoa(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }

// quickStyle is the theme slice a shape selects.
type quickStyle struct {
	lineColor, fillColor, shadowColor, fontColor      int
	lineMatrix, fillMatrix, effectsMatrix, fontMatrix int
	variation                                         int
	varColor, varStyle                                int
	connector                                         bool
}

func (p *pageCtx) quickStyle(s *shape) *quickStyle {
	if q, ok := p.qs[s]; ok {
		return q
	}
	get := func(name string, def int) int {
		c, ok := s.get(name)
		if !ok || c.v == "Themed" {
			return def
		}
		return atoi(c.v, def)
	}
	q := &quickStyle{
		lineColor: get("QuickStyleLineColor", 100), fillColor: get("QuickStyleFillColor", 100),
		shadowColor: get("QuickStyleShadowColor", 100), fontColor: get("QuickStyleFontColor", 100),
		lineMatrix: get("QuickStyleLineMatrix", 100), fillMatrix: get("QuickStyleFillMatrix", 100),
		effectsMatrix: get("QuickStyleEffectsMatrix", 100), fontMatrix: get("QuickStyleFontMatrix", 100),
		variation: get("QuickStyleVariation", 0),
		varColor:  get("VariationColorIndex", 65534), varStyle: get("VariationStyleIndex", 65534),
	}
	// 65534: the page's variant
	if q.varColor == 65534 {
		q.varColor = atoi(p.pg.sheet.cells["VariationColorIndex"].v, 0)
	}
	if q.varStyle == 65534 {
		q.varStyle = atoi(p.pg.sheet.cells["VariationStyleIndex"].v, 0)
	}
	switch get("QuickStyleType", 0) {
	case 1, 2:
	case 3:
		q.connector = true
	default:
		for _, st := range s.styleChain(catOther) {
			if strings.EqualFold(st.name, "Connector") {
				q.connector = true
				break
			}
		}
	}
	p.qs[s] = q
	return q
}

// themeIndex returns the theme a shape uses: 0 for none, -1 when neither
// the shape nor its page says (the document's theme then applies).
func (p *pageCtx) themeIndex(s *shape) int {
	v := 65534
	if c, ok := s.get("ThemeIndex"); ok && c.v != "Themed" {
		v = atoi(c.v, 65534)
	}
	if v == 65534 {
		c, ok := p.pg.sheet.cells["ThemeIndex"]
		if !ok {
			return -1
		}
		v = atoi(c.v, -1)
	}
	return v
}

// themeColor returns the color a QuickStyle color index selects.
func (p *pageCtx) themeColor(q *quickStyle, idx int) bdf.Color {
	th := p.th
	names := []string{"dk1", "lt1", "accent1", "accent2", "accent3", "accent4", "accent5", "accent6"}
	switch {
	case idx >= 0 && idx < len(names):
		return th.colors[names[idx]]
	case idx == 8:
		return p.background()
	case idx >= 100 && idx <= 106, idx >= 200 && idx <= 206:
		// the colors of the page's variant (200-206 as 100-106)
		if len(th.varColors) == 0 {
			return th.colors["accent1"]
		}
		v := q.varColor
		if v < 0 || v >= len(th.varColors) {
			v = 0
		}
		if i := idx % 100; i < len(th.varColors[v]) {
			return th.varColors[v][i]
		}
		return th.colors["accent1"]
	}
	return th.colors["dk1"]
}

func (p *pageCtx) background() bdf.Color {
	if c, ok := p.th.colors["bkgnd"]; ok {
		return c
	}
	if c, ok := p.th.colors["lt1"]; ok {
		return c
	}
	return 0xffffffff
}

// matrix returns the style list index a QuickStyle matrix cell selects
// (0: the root style sheet's values).
func (p *pageCtx) matrix(q *quickStyle, idx int, pick func(varStyle) int) int {
	if idx >= 100 && idx <= 103 {
		if q.varStyle >= 0 && q.varStyle < len(p.th.varStyles) && idx-100 < len(p.th.varStyles[q.varStyle]) {
			return pick(p.th.varStyles[q.varStyle][idx-100])
		}
		return 1
	}
	return idx
}

// themed resolves a cell a dynamic theme sets. sec and key name the row of
// cells in sections (Character, FillGradient …).
func (p *pageCtx) themed(s *shape, sec, key, name string) (string, bool) {
	if p.th == nil || p.themeIndex(s) == 0 {
		// no theme applies: the root style sheet's values
		return p.rootValue(sec, key, name)
	}
	q := p.quickStyle(s)
	fs := &p.th.shapes
	if q.connector {
		fs = &p.th.conns
	}
	cc := p.th.colors
	switch {
	case sec == "Character":
		fm := p.matrix(q, q.fontMatrix, func(v varStyle) int { return v.font })
		if fm == 0 {
			return p.rootValue(sec, key, name)
		}
		fp := entry(fs.fontProps, fm)
		switch name {
		case "Color":
			c := p.themeColor(q, q.fontColor)
			if fp != nil {
				if r, ok := drawingml.ResolveColor(fp.Child("color"), cc, c); ok {
					c = r
				}
			}
			c = p.variationColor(q, 0x2, c)
			return hexOf(c), true
		case "ColorTrans":
			return "0", true
		case "Style":
			if fp != nil {
				return fp.AttrStr("style", "0"), true
			}
			return "0", true
		case "Font":
			return p.th.minor.latin, true
		case "AsianFont":
			if p.th.minor.ea != "" {
				return p.th.minor.ea, true
			}
			return p.th.minor.script("Jpan"), true
		case "ComplexScriptFont":
			return "", true
		}
	case isLineCell(name):
		lm := p.matrix(q, q.lineMatrix, func(v varStyle) int { return v.line })
		if lm == 0 {
			return p.rootValue(sec, key, name)
		}
		ln := entry(fs.lines, lm)
		lx := fs.lineExFor(lm)
		switch name {
		case "LineColor", "LineColorTrans":
			c := p.themeColor(q, q.lineColor)
			if f := fillOf(ln); f != nil && f.Name == "solidFill" {
				if r, ok := drawingml.ResolveColor(f, cc, c); ok {
					c = r
				}
			}
			c = p.variationColor(q, 0x4, c)
			if name == "LineColorTrans" {
				return ftoa(1 - alphaOf(c)), true
			}
			return hexOf(c), true
		case "LineWeight":
			return ftoa(ln.AttrEMU("w", 9525.0/ooxml.EMUPerPoint) / 72), true
		case "LinePattern":
			if f := fillOf(ln); f != nil && f.Name == "noFill" {
				return "0", true
			}
			if v, ok := lx.Attr("pattern"); ok {
				return v, true
			}
			return strconv.Itoa(dashPattern(ln.Child("prstDash").AttrStr("val", "solid"))), true
		case "LineCap":
			switch ln.AttrStr("cap", "sq") {
			case "rnd":
				return "0", true
			case "flat":
				return "1", true
			}
			return "2", true
		case "Rounding":
			return ftoa(lx.AttrEMU("rndg", 0) / 72), true
		case "BeginArrow":
			return lx.AttrStr("start", "0"), true
		case "EndArrow":
			return lx.AttrStr("end", "0"), true
		case "BeginArrowSize":
			return lx.AttrStr("startSize", "2"), true
		case "EndArrowSize":
			return lx.AttrStr("endSize", "2"), true
		case "CompoundType", "LineGradientEnabled", "LineGradientDir", "LineGradientAngle":
			return "0", true
		}
	case sec == "FillGradient" || isFillCell(name):
		fm := p.matrix(q, q.fillMatrix, func(v varStyle) int { return v.fill })
		if fm == 0 {
			return p.rootValue(sec, key, name)
		}
		f := entry(fs.fills, fm)
		ph := p.variationColor(q, 0x8, p.themeColor(q, q.fillColor))
		stops := gradientStops(f, cc, ph)
		switch name {
		case "FillForegnd", "FillForegndTrans", "FillBkgnd", "FillBkgndTrans":
			c := ph
			switch {
			case f != nil && f.Name == "solidFill":
				if r, ok := drawingml.ResolveColor(f, cc, ph); ok {
					c = r
				}
			case len(stops) > 0 && (name == "FillBkgnd" || name == "FillBkgndTrans"):
				c = stops[len(stops)-1].c
			case len(stops) > 0:
				c = stops[0].c
			}
			if strings.HasSuffix(name, "Trans") {
				return ftoa(1 - alphaOf(c)), true
			}
			return hexOf(c), true
		case "FillPattern":
			if f == nil || f.Name == "noFill" {
				return "0", true
			}
			if fp := entry(fs.fillProps, fm); fp != nil {
				if v, ok := fp.Attr("pattern"); ok {
					return v, true
				}
			}
			return "1", true
		case "FillGradientEnabled":
			if len(stops) > 1 {
				return "1", true
			}
			return "0", true
		case "FillGradientDir":
			if path := f.Child("path"); path != nil {
				if path.AttrStr("path", "circle") == "circle" {
					return "3", true
				}
				return "10", true
			}
			return "0", true
		case "FillGradientAngle":
			// DrawingML turns clockwise in a y-down space; Visio
			// measures counterclockwise from the x axis
			return ftoa(-float64(f.Child("lin").AttrInt("ang", 0)) / 60000 * math.Pi / 180), true
		case "RotateGradientWithShape":
			return strconv.Itoa(boolInt(f.AttrBool("rotWithShape", true))), true
		case "UseGroupGradient":
			return "0", true
		case "GradientStopColor", "GradientStopColorTrans", "GradientStopPosition":
			i := atoi(key, 0)
			if i < 0 || i >= len(stops) {
				return "", false
			}
			switch name {
			case "GradientStopColor":
				return hexOf(stops[i].c), true
			case "GradientStopColorTrans":
				return ftoa(1 - alphaOf(stops[i].c)), true
			}
			return ftoa(stops[i].pos), true
		}
	case isShadowCell(name):
		em := p.matrix(q, q.effectsMatrix, func(v varStyle) int { return v.effect })
		if em == 0 {
			return p.rootValue(sec, key, name)
		}
		sh := entry(fs.effects, em).Path("effectLst", "outerShdw")
		switch name {
		case "ShdwPattern":
			return strconv.Itoa(boolInt(sh != nil)), true
		case "ShdwForegnd", "ShdwForegndTrans":
			c := p.themeColor(q, q.shadowColor)
			if r, ok := drawingml.ResolveColor(sh, cc, c); ok {
				c = r
			}
			if name == "ShdwForegndTrans" {
				return ftoa(1 - alphaOf(c)), true
			}
			return hexOf(c), true
		case "ShapeShdwType":
			return "1", true
		case "ShapeShdwOffsetX", "ShapeShdwOffsetY":
			dist := sh.AttrEMU("dist", 0) / 72
			dir := float64(sh.AttrInt("dir", 0)) / 60000 * math.Pi / 180
			if name == "ShapeShdwOffsetX" {
				return ftoa(dist * math.Cos(dir)), true
			}
			return ftoa(-dist * math.Sin(dir)), true
		case "ShapeShdwBlur":
			return ftoa(sh.AttrEMU("blurRad", 0) / 72), true
		case "ShapeShdwScaleFactor":
			return "1", true
		case "ShapeShdwObliqueAngle":
			return "0", true
		}
	}
	return p.rootValue(sec, key, name)
}

// rootValue returns the value of the root style sheet ("No Style"), which
// cells the theme does not set fall back to.
func (p *pageCtx) rootValue(sec, key, name string) (string, bool) {
	root := p.c.d.styles[0]
	if root == nil {
		return "", false
	}
	var c cell
	var ok bool
	if sec == "" {
		c, ok = root.cells[name]
	} else if x := root.sections[secKey{sec, -1}]; x != nil {
		if r := x.rows[key]; r != nil {
			c, ok = r.cells[name]
		} else if r := x.rows["0"]; r != nil {
			c, ok = r.cells[name]
		}
	}
	if !ok || c.v == "Themed" {
		return "", false
	}
	return c.v, true
}

// variationColor applies the contrast rule of QuickStyleVariation: a color
// too close to the page background in luminance is replaced.
func (p *pageCtx) variationColor(q *quickStyle, bit int, c bdf.Color) bdf.Color {
	if q.variation&bit == 0 {
		return c
	}
	bg := luminance(p.background())
	if math.Abs(bg-luminance(c)) >= 0.1666 {
		return c
	}
	if bg <= 0.7292 {
		return 0xffffffff
	}
	// the most contrasting of the shape's text, fill and line colors
	best, bestD := c, 0.0
	for _, i := range []int{q.fontColor, q.fillColor, q.lineColor} {
		k := p.themeColor(q, i)
		if d := math.Abs(bg - luminance(k)); d > bestD {
			best, bestD = k, d
		}
	}
	return best
}

// luminance is the HSL luminance of a color.
func luminance(c bdf.Color) float64 {
	r, g, b := float64(c>>24)/255, float64(c>>16&0xff)/255, float64(c>>8&0xff)/255
	return (math.Max(r, math.Max(g, b)) + math.Min(r, math.Min(g, b))) / 2
}

func isLineCell(name string) bool {
	switch name {
	case "LineColor", "LineColorTrans", "LineWeight", "LinePattern", "LineCap", "Rounding", "BeginArrow", "EndArrow",
		"BeginArrowSize", "EndArrowSize", "CompoundType", "LineGradientEnabled", "LineGradientDir", "LineGradientAngle":
		return true
	}
	return false
}

func isFillCell(name string) bool {
	switch name {
	case "FillForegnd", "FillForegndTrans", "FillBkgnd", "FillBkgndTrans", "FillPattern", "FillGradientEnabled",
		"FillGradientDir", "FillGradientAngle", "RotateGradientWithShape", "UseGroupGradient":
		return true
	}
	return false
}

func isShadowCell(name string) bool {
	switch name {
	case "ShdwPattern", "ShdwForegnd", "ShdwForegndTrans", "ShapeShdwType", "ShapeShdwOffsetX", "ShapeShdwOffsetY",
		"ShapeShdwBlur", "ShapeShdwScaleFactor", "ShapeShdwObliqueAngle":
		return true
	}
	return false
}

// fillOf returns the fill choice of a DrawingML line or fill style.
func fillOf(n *ooxml.Node) *ooxml.Node {
	if n == nil {
		return nil
	}
	switch n.Name {
	case "noFill", "solidFill", "gradFill", "pattFill", "blipFill":
		return n
	}
	for _, k := range n.Kids {
		switch k.Name {
		case "noFill", "solidFill", "gradFill", "pattFill", "blipFill":
			return k
		}
	}
	return nil
}

type gradStop struct {
	pos float64
	c   bdf.Color
}

// gradientStops resolves the stops of a DrawingML gradient fill.
func gradientStops(f *ooxml.Node, scheme map[string]bdf.Color, ph bdf.Color) []gradStop {
	if f == nil || f.Name != "gradFill" {
		return nil
	}
	var out []gradStop
	for _, gs := range f.Child("gsLst").Children("gs") {
		if c, ok := drawingml.ResolveColor(gs, scheme, ph); ok {
			out = append(out, gradStop{gs.AttrPct("pos", 0), c})
		}
	}
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].pos < out[j-1].pos; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// dashPattern maps a DrawingML preset dash to a Visio line pattern.
func dashPattern(v string) int {
	switch v {
	case "dash", "sysDash":
		return 2
	case "dot", "sysDot":
		return 3
	case "dashDot", "sysDashDot":
		return 4
	case "lgDashDotDot", "sysDashDotDot":
		return 5
	case "lgDash":
		return 16
	case "lgDashDot":
		return 18
	}
	return 1
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
