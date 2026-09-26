package drawingml

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// rgba is a color with components in 0..1 (sRGB).
type rgba struct{ R, G, B, A float64 }

var (
	black = rgba{0, 0, 0, 1}
	white = rgba{1, 1, 1, 1}
)

func (c rgba) bdf() bdf.Color {
	q := func(v float64) uint8 { return uint8(math.Round(clamp01(v) * 255)) }
	return bdf.RGBA(q(c.R), q(c.G), q(c.B), q(c.A))
}

func clamp01(v float64) float64 {
	if v < 0 || math.IsNaN(v) {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func hexColor(s string) (rgba, bool) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(s) != 6 {
		return rgba{}, false
	}
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return rgba{}, false
	}
	return rgba{float64(v>>16&0xff) / 255, float64(v>>8&0xff) / 255, float64(v&0xff) / 255, 1}, true
}

// colorCtx resolves DrawingML colors: the theme's scheme, the color map of
// the page (bg1 → lt1 …) and the placeholder color of style references.
type colorCtx struct {
	scheme map[string]rgba
	clrMap map[string]string
	phClr  *rgba
}

func (cc *colorCtx) withPh(c rgba) *colorCtx {
	n := *cc
	n.phClr = &c
	return &n
}

var defaultClrMap = map[string]string{
	"bg1": "lt1", "tx1": "dk1", "bg2": "lt2", "tx2": "dk2",
	"accent1": "accent1", "accent2": "accent2", "accent3": "accent3", "accent4": "accent4",
	"accent5": "accent5", "accent6": "accent6", "hlink": "hlink", "folHlink": "folHlink",
}

func (cc *colorCtx) scheme1(name string) (rgba, bool) {
	if name == "phClr" {
		if cc.phClr != nil {
			return *cc.phClr, true
		}
		return black, true
	}
	if m, ok := cc.clrMap[name]; ok {
		name = m
	} else if m, ok := defaultClrMap[name]; ok {
		name = m
	}
	c, ok := cc.scheme[name]
	return c, ok
}

// isColor reports whether an element is one of the color choices.
func isColor(n *ooxml.Node) bool {
	switch n.Name {
	case "srgbClr", "schemeClr", "sysClr", "prstClr", "scrgbClr", "hslClr":
		return true
	}
	return false
}

// color resolves the first color choice among n's children.
func (cc *colorCtx) color(n *ooxml.Node) (rgba, bool) {
	if n == nil {
		return rgba{}, false
	}
	for _, k := range n.Kids {
		if isColor(k) {
			return cc.resolve(k)
		}
	}
	return rgba{}, false
}

// resolve evaluates a color choice element and its modifiers.
func (cc *colorCtx) resolve(n *ooxml.Node) (rgba, bool) {
	var c rgba
	ok := true
	switch n.Name {
	case "srgbClr":
		c, ok = hexColor(n.AttrStr("val", ""))
	case "schemeClr":
		c, ok = cc.scheme1(n.AttrStr("val", ""))
	case "sysClr":
		if c, ok = hexColor(n.AttrStr("lastClr", "")); !ok {
			switch n.AttrStr("val", "") {
			case "window", "btnHighlight", "highlightText", "menu", "info":
				c, ok = white, true
			default:
				c, ok = black, true
			}
		}
	case "prstClr":
		c, ok = presetColor(n.AttrStr("val", ""))
	case "scrgbClr":
		c = rgba{fromLinear(n.AttrPct("r", 0)), fromLinear(n.AttrPct("g", 0)), fromLinear(n.AttrPct("b", 0)), 1}
	case "hslClr":
		c = hslToRGB(float64(n.AttrInt("hue", 0))/60000/360, n.AttrPct("sat", 0), n.AttrPct("lum", 0))
		c.A = 1
	default:
		return rgba{}, false
	}
	if !ok {
		return rgba{}, false
	}
	return applyMods(c, n), true
}

func applyMods(c rgba, n *ooxml.Node) rgba {
	for _, m := range n.Kids {
		v := m.AttrPct("val", 0)
		switch m.Name {
		case "alpha":
			c.A = v
		case "alphaMod":
			c.A *= v
		case "alphaOff":
			c.A += v
		case "lumMod", "lumOff", "satMod", "satOff", "hueMod", "hueOff", "hue", "sat", "lum":
			h, s, l := rgbToHSL(c)
			switch m.Name {
			case "lumMod":
				l *= v
			case "lumOff":
				l += v
			case "lum":
				l = v
			case "satMod":
				s *= v
			case "satOff":
				s += v
			case "sat":
				s = v
			case "hueMod":
				h *= v
			case "hueOff":
				h += float64(m.AttrInt("val", 0)) / 60000 / 360
			case "hue":
				h = float64(m.AttrInt("val", 0)) / 60000 / 360
			}
			h -= math.Floor(h)
			a := c.A
			c = hslToRGB(h, clamp01(s), clamp01(l))
			c.A = a
		case "tint":
			// Office mixes with white in linear RGB.
			c.R = fromLinear(1 - (1-toLinear(c.R))*v)
			c.G = fromLinear(1 - (1-toLinear(c.G))*v)
			c.B = fromLinear(1 - (1-toLinear(c.B))*v)
		case "shade":
			c.R = fromLinear(toLinear(c.R) * v)
			c.G = fromLinear(toLinear(c.G) * v)
			c.B = fromLinear(toLinear(c.B) * v)
		case "comp":
			h, s, l := rgbToHSL(c)
			a := c.A
			c = hslToRGB(math.Mod(h+0.5, 1), s, l)
			c.A = a
		case "inv":
			c.R, c.G, c.B = 1-c.R, 1-c.G, 1-c.B
		case "gray":
			y := 0.2126*c.R + 0.7152*c.G + 0.0722*c.B
			c.R, c.G, c.B = y, y, y
		case "gamma":
			c.R, c.G, c.B = fromLinear(c.R), fromLinear(c.G), fromLinear(c.B)
		case "invGamma":
			c.R, c.G, c.B = toLinear(c.R), toLinear(c.G), toLinear(c.B)
		case "red":
			c.R = v
		case "redMod":
			c.R *= v
		case "redOff":
			c.R += v
		case "green":
			c.G = v
		case "greenMod":
			c.G *= v
		case "greenOff":
			c.G += v
		case "blue":
			c.B = v
		case "blueMod":
			c.B *= v
		case "blueOff":
			c.B += v
		}
		c.R, c.G, c.B, c.A = clamp01(c.R), clamp01(c.G), clamp01(c.B), clamp01(c.A)
	}
	return c
}

func toLinear(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func fromLinear(v float64) float64 {
	v = clamp01(v)
	if v <= 0.0031308 {
		return v * 12.92
	}
	return 1.055*math.Pow(v, 1/2.4) - 0.055
}

func rgbToHSL(c rgba) (h, s, l float64) {
	mx := math.Max(c.R, math.Max(c.G, c.B))
	mn := math.Min(c.R, math.Min(c.G, c.B))
	l = (mx + mn) / 2
	if mx == mn {
		return 0, 0, l
	}
	d := mx - mn
	if l > 0.5 {
		s = d / (2 - mx - mn)
	} else {
		s = d / (mx + mn)
	}
	switch mx {
	case c.R:
		h = (c.G - c.B) / d
		if c.G < c.B {
			h += 6
		}
	case c.G:
		h = (c.B-c.R)/d + 2
	default:
		h = (c.R-c.G)/d + 4
	}
	return h / 6, s, l
}

func hslToRGB(h, s, l float64) rgba {
	if s == 0 {
		return rgba{l, l, l, 1}
	}
	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q
	hue := func(t float64) float64 {
		t -= math.Floor(t)
		switch {
		case t < 1.0/6:
			return p + (q-p)*6*t
		case t < 0.5:
			return q
		case t < 2.0/3:
			return p + (q-p)*(2.0/3-t)*6
		}
		return p
	}
	return rgba{hue(h + 1.0/3), hue(h), hue(h - 1.0/3), 1}
}

// presetColor maps ST_PresetColorVal names ("dkBlue", "ltGray" …) to the
// CSS named colors they stand for.
func presetColor(name string) (rgba, bool) {
	n := strings.ToLower(name)
	switch {
	case strings.HasPrefix(n, "dk"):
		n = "dark" + n[2:]
	case strings.HasPrefix(n, "lt"):
		n = "light" + n[2:]
	case strings.HasPrefix(n, "med"):
		n = "medium" + n[3:]
	}
	if v, ok := cssColors[n]; ok {
		return rgba{float64(v>>16&0xff) / 255, float64(v>>8&0xff) / 255, float64(v&0xff) / 255, 1}, true
	}
	return rgba{}, false
}

var cssColors = map[string]uint32{
	"aliceblue": 0xf0f8ff, "antiquewhite": 0xfaebd7, "aqua": 0x00ffff, "aquamarine": 0x7fffd4, "azure": 0xf0ffff,
	"beige": 0xf5f5dc, "bisque": 0xffe4c4, "black": 0x000000, "blanchedalmond": 0xffebcd, "blue": 0x0000ff,
	"blueviolet": 0x8a2be2, "brown": 0xa52a2a, "burlywood": 0xdeb887, "cadetblue": 0x5f9ea0, "chartreuse": 0x7fff00,
	"chocolate": 0xd2691e, "coral": 0xff7f50, "cornflowerblue": 0x6495ed, "cornsilk": 0xfff8dc, "crimson": 0xdc143c,
	"cyan": 0x00ffff, "darkblue": 0x00008b, "darkcyan": 0x008b8b, "darkgoldenrod": 0xb8860b, "darkgray": 0xa9a9a9,
	"darkgrey": 0xa9a9a9, "darkgreen": 0x006400, "darkkhaki": 0xbdb76b, "darkmagenta": 0x8b008b, "darkolivegreen": 0x556b2f,
	"darkorange": 0xff8c00, "darkorchid": 0x9932cc, "darkred": 0x8b0000, "darksalmon": 0xe9967a, "darkseagreen": 0x8fbc8f,
	"darkslateblue": 0x483d8b, "darkslategray": 0x2f4f4f, "darkslategrey": 0x2f4f4f, "darkturquoise": 0x00ced1,
	"darkviolet": 0x9400d3, "deeppink": 0xff1493, "deepskyblue": 0x00bfff, "dimgray": 0x696969, "dimgrey": 0x696969,
	"dodgerblue": 0x1e90ff, "firebrick": 0xb22222, "floralwhite": 0xfffaf0, "forestgreen": 0x228b22, "fuchsia": 0xff00ff,
	"gainsboro": 0xdcdcdc, "ghostwhite": 0xf8f8ff, "gold": 0xffd700, "goldenrod": 0xdaa520, "gray": 0x808080,
	"grey": 0x808080, "green": 0x008000, "greenyellow": 0xadff2f, "honeydew": 0xf0fff0, "hotpink": 0xff69b4,
	"indianred": 0xcd5c5c, "indigo": 0x4b0082, "ivory": 0xfffff0, "khaki": 0xf0e68c, "lavender": 0xe6e6fa,
	"lavenderblush": 0xfff0f5, "lawngreen": 0x7cfc00, "lemonchiffon": 0xfffacd, "lightblue": 0xadd8e6, "lightcoral": 0xf08080,
	"lightcyan": 0xe0ffff, "lightgoldenrodyellow": 0xfafad2, "lightgray": 0xd3d3d3, "lightgrey": 0xd3d3d3, "lightgreen": 0x90ee90,
	"lightpink": 0xffb6c1, "lightsalmon": 0xffa07a, "lightseagreen": 0x20b2aa, "lightskyblue": 0x87cefa, "lightslategray": 0x778899,
	"lightslategrey": 0x778899, "lightsteelblue": 0xb0c4de, "lightyellow": 0xffffe0, "lime": 0x00ff00, "limegreen": 0x32cd32,
	"linen": 0xfaf0e6, "magenta": 0xff00ff, "maroon": 0x800000, "mediumaquamarine": 0x66cdaa, "mediumblue": 0x0000cd,
	"mediumorchid": 0xba55d3, "mediumpurple": 0x9370db, "mediumseagreen": 0x3cb371, "mediumslateblue": 0x7b68ee,
	"mediumspringgreen": 0x00fa9a, "mediumturquoise": 0x48d1cc, "mediumvioletred": 0xc71585, "midnightblue": 0x191970,
	"mintcream": 0xf5fffa, "mistyrose": 0xffe4e1, "moccasin": 0xffe4b5, "navajowhite": 0xffdead, "navy": 0x000080,
	"oldlace": 0xfdf5e6, "olive": 0x808000, "olivedrab": 0x6b8e23, "orange": 0xffa500, "orangered": 0xff4500,
	"orchid": 0xda70d6, "palegoldenrod": 0xeee8aa, "palegreen": 0x98fb98, "paleturquoise": 0xafeeee, "palevioletred": 0xdb7093,
	"papayawhip": 0xffefd5, "peachpuff": 0xffdab9, "peru": 0xcd853f, "pink": 0xffc0cb, "plum": 0xdda0dd,
	"powderblue": 0xb0e0e6, "purple": 0x800080, "red": 0xff0000, "rosybrown": 0xbc8f8f, "royalblue": 0x4169e1,
	"saddlebrown": 0x8b4513, "salmon": 0xfa8072, "sandybrown": 0xf4a460, "seagreen": 0x2e8b57, "seashell": 0xfff5ee,
	"sienna": 0xa0522d, "silver": 0xc0c0c0, "skyblue": 0x87ceeb, "slateblue": 0x6a5acd, "slategray": 0x708090,
	"slategrey": 0x708090, "snow": 0xfffafa, "springgreen": 0x00ff7f, "steelblue": 0x4682b4, "tan": 0xd2b48c,
	"teal": 0x008080, "thistle": 0xd8bfd8, "tomato": 0xff6347, "turquoise": 0x40e0d0, "violet": 0xee82ee,
	"wheat": 0xf5deb3, "white": 0xffffff, "whitesmoke": 0xf5f5f5, "yellow": 0xffff00, "yellowgreen": 0x9acd32,
}
