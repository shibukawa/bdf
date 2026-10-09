package idml

import (
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// Swatches (Resources/Graphic.xml) are referred to by their Self: a Color
// (CMYK, RGB or Lab values), a Tint of a color, a Gradient of stops, a
// MixedInk, or the swatch None. Colors are converted to sRGB the simple
// way, without profiles: CMYK by multiplying out the inks, Lab through
// XYZ with the D50 white point.

// rgba is a color with alpha in 0..1.
type rgba struct{ r, g, b, a float64 }

func (c rgba) bdf() bdf.Color {
	u := func(v float64) uint8 { return uint8(math.Round(math.Max(0, math.Min(1, v)) * 255)) }
	return bdf.RGBA(u(c.r), u(c.g), u(c.b), u(c.a))
}

// tint mixes the color with white: 100 is the color, 0 white.
func (c rgba) tint(pct float64) rgba {
	t := math.Max(0, math.Min(100, pct)) / 100
	return rgba{c.r*t + (1 - t), c.g*t + (1 - t), c.b*t + (1 - t), c.a}
}

func (c rgba) alpha(a float64) rgba { c.a *= a; return c }

var black = rgba{0, 0, 0, 1}

// colorValue converts a Color element's values to sRGB.
func colorValue(n *ooxml.Node) (rgba, bool) {
	v := nums(n.AttrStr("ColorValue", ""))
	switch strings.ToUpper(n.AttrStr("Space", "")) {
	case "CMYK":
		if len(v) != 4 {
			return rgba{}, false
		}
		k := 1 - v[3]/100
		return rgba{(1 - v[0]/100) * k, (1 - v[1]/100) * k, (1 - v[2]/100) * k, 1}, true
	case "RGB":
		if len(v) != 3 {
			return rgba{}, false
		}
		return rgba{v[0] / 255, v[1] / 255, v[2] / 255, 1}, true
	case "LAB":
		if len(v) != 3 {
			return rgba{}, false
		}
		return labToRGB(v[0], v[1], v[2]), true
	}
	return rgba{}, false
}

// labToRGB converts CIE L*a*b* (D50) to sRGB.
func labToRGB(l, a, b float64) rgba {
	fy := (l + 16) / 116
	fx := fy + a/500
	fz := fy - b/200
	finv := func(t float64) float64 {
		if t > 6.0/29 {
			return t * t * t
		}
		return 3 * (6.0 / 29) * (6.0 / 29) * (t - 4.0/29)
	}
	// D50 white, then Bradford-adapted to D65 for sRGB
	x, y, z := 0.9642*finv(fx), 1.0*finv(fy), 0.8251*finv(fz)
	r := 3.1339*x - 1.6170*y - 0.4906*z
	g := -0.9785*x + 1.9160*y + 0.0333*z
	bb := 0.0720*x - 0.2290*y + 1.4057*z
	gamma := func(c float64) float64 {
		c = math.Max(0, math.Min(1, c))
		if c <= 0.0031308 {
			return 12.92 * c
		}
		return 1.055*math.Pow(c, 1/2.4) - 0.055
	}
	return rgba{gamma(r), gamma(g), gamma(bb), 1}
}

// swatch resolves a swatch reference to a solid color: false for None, a
// gradient, or an unknown swatch (which is black, with a warning).
func (c *converter) swatch(ref string) (rgba, bool) {
	return c.swatchDepth(ref, 0)
}

func (c *converter) swatchDepth(ref string, depth int) (rgba, bool) {
	if ref == "" || ref == "Swatch/None" || ref == "n" || depth > 8 {
		return rgba{}, false
	}
	n := c.d.swatches[ref]
	switch swatchKind(n) {
	case "Color":
		if col, ok := colorValue(n); ok {
			return col, true
		}
	case "Tint":
		if base, ok := c.swatchDepth(n.AttrStr("BaseColor", ""), depth+1); ok {
			return base.tint(n.AttrFloat("TintValue", 100)), true
		}
	case "Gradient":
		return rgba{}, false
	case "MixedInk":
		// the inks' mix is not computed: its first ink stands for it
		c.warnOnce("mixedink", "mixed ink swatches are drawn as black")
		return black, true
	}
	switch {
	case strings.HasPrefix(ref, "Color/Black") || ref == "Color/Registration":
		return black, true
	case strings.HasPrefix(ref, "Color/Paper"):
		return rgba{1, 1, 1, 1}, true
	}
	c.warnOnce("swatch:"+ref, "unknown swatch %s: drawn as black", ref)
	return black, true
}

// swatchKind returns the element name of a swatch ("" for none).
func swatchKind(n *ooxml.Node) string {
	if n == nil {
		return ""
	}
	return n.Name
}

// gradient returns the Gradient element a reference names, or nil.
func (c *converter) gradient(ref string) *ooxml.Node {
	if n := c.d.swatches[ref]; n != nil && n.Name == "Gradient" {
		return n
	}
	return nil
}

// stops reads the stops of a gradient, tinted and with alpha.
func (c *converter) stops(g *ooxml.Node, tint, alpha float64) []bdf.Stop {
	var out []bdf.Stop
	for _, s := range g.Children("GradientStop") {
		col, ok := c.swatch(s.AttrStr("StopColor", ""))
		if !ok {
			col = rgba{1, 1, 1, 0}
		}
		if tint >= 0 {
			col = col.tint(tint)
		}
		out = append(out, bdf.Stop{Offset: f32(math.Max(0, math.Min(1, s.AttrFloat("Location", 0)/100))), Color: col.alpha(alpha).bdf()})
	}
	if len(out) == 1 {
		out = append(out, out[0])
		out[1].Offset = 1
	}
	return out
}

// paintOf is a fill or a stroke color: solid, or a gradient.
type paintOf struct {
	solid rgba
	grad  *ooxml.Node
	tint  float64 // -1 for none
	alpha float64
}

func (p paintOf) ok() bool { return p.grad != nil || p.solid.a > 0 }

// average returns a color standing for the paint (a gradient's first stop).
func (c *converter) average(p paintOf) rgba {
	if p.grad == nil {
		return p.solid
	}
	st := c.stops(p.grad, p.tint, p.alpha)
	if len(st) == 0 {
		return black.alpha(p.alpha)
	}
	v := st[0].Color
	return rgba{float64(uint8(v>>24)) / 255, float64(uint8(v>>16)) / 255, float64(uint8(v>>8)) / 255, float64(uint8(v)) / 255}
}

// paint resolves a swatch reference with its tint and opacity.
func (c *converter) paint(ref string, tint, opacity float64) paintOf {
	p := paintOf{tint: tint, alpha: opacity}
	if g := c.gradient(ref); g != nil {
		p.grad = g
		return p
	}
	if col, ok := c.swatch(ref); ok {
		if tint >= 0 {
			col = col.tint(tint)
		}
		p.solid = col.alpha(opacity)
	}
	return p
}

// opacity reads the opacity (0..1) of a transparency setting element
// (TransparencySetting, FillTransparencySetting, StrokeTransparencySetting).
func opacity(n *ooxml.Node) float64 {
	if n == nil {
		return 1
	}
	return math.Max(0, math.Min(100, n.Child("BlendingSetting").AttrFloat("Opacity", 100))) / 100
}

// gradientPaint builds the paint of a gradient fill over a box in item
// coordinates: a linear gradient along GradientFillAngle (counterclockwise
// from the x axis, as InDesign measures it) starting at GradientFillStart
// for GradientFillLength, or across the box when the item gives neither;
// a radial one from the start point with the length as its radius.
func (c *converter) gradientPaint(p paintOf, n *ooxml.Node, box rect, prefix string) bdf.Paint {
	stops := c.stops(p.grad, p.tint, p.alpha)
	angle := n.AttrFloat(prefix+"Angle", 0) * math.Pi / 180
	dx, dy := math.Cos(angle), -math.Sin(angle)
	cx, cy := (box.x0+box.x1)/2, (box.y0+box.y1)/2
	// the extent of the box along the direction
	ext := math.Abs(box.w()*dx) + math.Abs(box.h()*dy)
	sx, sy := cx-dx*ext/2, cy-dy*ext/2
	length := ext
	if s := nums(n.AttrStr(prefix+"Start", "")); len(s) == 2 {
		sx, sy = s[0], s[1]
	}
	if l, ok := n.Attr(prefix + "Length"); ok {
		if v := nums(l); len(v) == 1 && v[0] > 0 {
			length = v[0]
		}
	}
	if p.grad.AttrStr("Type", "Linear") == "Radial" {
		if s := nums(n.AttrStr(prefix+"Start", "")); len(s) != 2 {
			sx, sy = cx, cy
			length = math.Hypot(box.w(), box.h()) / 2
		}
		return bdf.RadialGradient(f32(sx), f32(sy), 0, f32(sx), f32(sy), f32(length), stops...)
	}
	return bdf.LinearGradient(f32(sx), f32(sy), f32(sx+dx*length), f32(sy+dy*length), stops...)
}

// strokeStyle is how a line is drawn.
type strokeStyle struct {
	width float64
	paint paintOf
	dash  []float32
	cap   byte
	join  byte
	miter float64
	align string // CenterAlignment, InsideAlignment or OutsideAlignment
}

// stroke reads the stroke of an item; false when it has none.
func (c *converter) stroke(n *ooxml.Node) (strokeStyle, bool) {
	s := strokeStyle{width: n.AttrFloat("StrokeWeight", 0), miter: n.AttrFloat("MiterLimit", 4), join: bdf.JoinMiter}
	if s.width <= 0 {
		return s, false
	}
	alpha := opacity(n.Child("TransparencySetting")) * opacity(n.Child("StrokeTransparencySetting"))
	s.paint = c.paint(n.AttrStr("StrokeColor", ""), n.AttrFloat("StrokeTint", -1), alpha)
	if !s.paint.ok() {
		return s, false
	}
	switch n.AttrStr("EndCap", "ButtEndCap") {
	case "RoundEndCap":
		s.cap = bdf.CapRound
	case "ProjectingEndCap":
		s.cap = bdf.CapSquare
	}
	switch n.AttrStr("EndJoin", "MiterEndJoin") {
	case "RoundEndJoin":
		s.join = bdf.JoinRound
	case "BevelEndJoin":
		s.join = bdf.JoinBevel
	}
	s.align = n.AttrStr("StrokeAlignment", "CenterAlignment")
	s.dash = c.dashes(n.AttrStr("StrokeType", "StrokeStyle/$ID/Solid"), s.width, &s.cap)
	return s, true
}

// dashes returns the dash pattern of a stroke style for a line of width w
// (nil for a solid line). The canned dashed and dotted styles are in
// multiples of the width; custom dashed styles give their dashes in points.
func (c *converter) dashes(ref string, w float64, cap *byte) []float32 {
	name := strings.TrimPrefix(ref, "StrokeStyle/")
	if n := c.d.strokeStyles[ref]; n != nil {
		switch n.Name {
		case "DashedStrokeStyle":
			var out []float32
			for _, v := range listValues(n, "DashArray") {
				if f := nums(v); len(f) == 1 {
					out = append(out, f32(f[0]))
				}
			}
			if len(out) == 0 {
				out = []float32{12, 4}
			}
			if n.AttrStr("LineCap", "ButtEndCap") == "RoundEndCap" {
				*cap = bdf.CapRound
			}
			return out
		case "DottedStrokeStyle":
			*cap = bdf.CapRound
			return []float32{0, f32(math.Max(n.AttrFloat("DotGap", 3*w), w))}
		case "StripedStrokeStyle":
			c.warnOnce("striped", "striped strokes are drawn as solid lines")
			return nil
		case "StrokeStyle":
			name = strings.TrimPrefix(n.AttrStr("Name", name), "$ID/")
		}
	}
	name = strings.TrimPrefix(name, "$ID/")
	switch {
	case name == "Solid":
		return nil
	case name == "Dashed":
		return []float32{12, 4}
	case name == "Dotted" || strings.EqualFold(name, "Canned dotted") || name == "Japanese Dots":
		*cap = bdf.CapRound
		return []float32{0, f32(2 * w)}
	case strings.HasPrefix(name, "Canned Dashed "):
		if f := nums(strings.ReplaceAll(strings.TrimPrefix(name, "Canned Dashed "), "x", " ")); len(f) == 2 {
			return []float32{f32(f[0] * w), f32(f[1] * w)}
		}
	case name == "ThinThin" || name == "ThickThin" || name == "ThinThick" || name == "ThickThick" || strings.Contains(name, "Thin") || strings.Contains(name, "Thick") || name == "Triple":
		c.warnOnce("striped", "striped strokes are drawn as solid lines")
		return nil
	case name == "White Diamond" || name == "Left Slant Hash" || name == "Right Slant Hash" || name == "Straight Hash" || name == "Wavy":
		c.warnOnce("stroke:"+name, "%s strokes are drawn as solid lines", name)
		return nil
	}
	c.warnOnce("stroke:"+ref, "unknown stroke style %s: drawn as a solid line", ref)
	return nil
}
