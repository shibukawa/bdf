package drawio

import (
	"math"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// rgba is a color with 8-bit channels.
type rgba struct{ r, g, b, a uint8 }

func (c rgba) bdf() bdf.Color { return bdf.RGBA(c.r, c.g, c.b, c.a) }

// withAlpha multiplies the color's alpha by f.
func (c rgba) withAlpha(f float64) rgba {
	c.a = uint8(math.Round(float64(c.a) * clamp(f, 0, 1)))
	return c
}

// parseColor reads a draw.io color value: #rgb, #rrggbb, #rrggbbaa, rgb()
// and rgba(), CSS color names, and light-dark(light, dark), of which the
// light color is used (BDF pages are drawn on white). ok is false for
// "none", "transparent", empty and unknown values.
func parseColor(s string) (rgba, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return rgba{}, false
	}
	low := strings.ToLower(s)
	if strings.HasPrefix(low, "light-dark(") && strings.HasSuffix(low, ")") {
		args := splitArgs(s[len("light-dark(") : len(s)-1])
		if len(args) > 0 {
			return parseColor(args[0])
		}
		return rgba{}, false
	}
	switch low {
	case "none", "transparent", "inherit", "initial", "unset":
		return rgba{}, false
	}
	if low[0] == '#' {
		h := low[1:]
		switch len(h) {
		case 3, 4:
			var v [4]uint8
			v[3] = 255
			for i := range len(h) {
				n, ok := hexNibble(h[i])
				if !ok {
					return rgba{}, false
				}
				v[i] = n * 17
			}
			return rgba{v[0], v[1], v[2], v[3]}, true
		case 6, 8:
			n, err := strconv.ParseUint(h, 16, 32)
			if err != nil {
				return rgba{}, false
			}
			if len(h) == 6 {
				return rgba{uint8(n >> 16), uint8(n >> 8), uint8(n), 255}, true
			}
			return rgba{uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(n)}, true
		}
		return rgba{}, false
	}
	if strings.HasPrefix(low, "rgb") {
		open, close := strings.IndexByte(low, '('), strings.LastIndexByte(low, ')')
		if open < 0 || close < open {
			return rgba{}, false
		}
		args := splitArgs(strings.ReplaceAll(low[open+1:close], "/", ","))
		if len(args) < 3 {
			return rgba{}, false
		}
		var ch [3]uint8
		for i := range 3 {
			a := strings.TrimSpace(args[i])
			f, ok := parseFloat(strings.TrimSuffix(a, "%"))
			if !ok {
				return rgba{}, false
			}
			if strings.HasSuffix(a, "%") {
				f = f * 255 / 100
			}
			ch[i] = uint8(math.Round(clamp(f, 0, 255)))
		}
		alpha := uint8(255)
		if len(args) > 3 {
			a := strings.TrimSpace(args[3])
			f, ok := parseFloat(strings.TrimSuffix(a, "%"))
			if ok {
				if strings.HasSuffix(a, "%") {
					f /= 100
				}
				alpha = uint8(math.Round(clamp(f, 0, 1) * 255))
			}
		}
		return rgba{ch[0], ch[1], ch[2], alpha}, true
	}
	if v, ok := cssColors[low]; ok {
		return rgba{uint8(v >> 16), uint8(v >> 8), uint8(v), 255}, true
	}
	return rgba{}, false
}

// splitArgs splits a comma separated argument list at the top level of
// parentheses.
func splitArgs(s string) []string {
	var out []string
	depth, start := 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	return append(out, strings.TrimSpace(s[start:]))
}

func hexNibble(c byte) (uint8, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	}
	return 0, false
}

// cssColors are the CSS named colors.
var cssColors = map[string]uint32{
	"aliceblue": 0xf0f8ff, "antiquewhite": 0xfaebd7, "aqua": 0x00ffff, "aquamarine": 0x7fffd4, "azure": 0xf0ffff,
	"beige": 0xf5f5dc, "bisque": 0xffe4c4, "black": 0x000000, "blanchedalmond": 0xffebcd, "blue": 0x0000ff,
	"blueviolet": 0x8a2be2, "brown": 0xa52a2a, "burlywood": 0xdeb887, "cadetblue": 0x5f9ea0, "chartreuse": 0x7fff00,
	"chocolate": 0xd2691e, "coral": 0xff7f50, "cornflowerblue": 0x6495ed, "cornsilk": 0xfff8dc, "crimson": 0xdc143c,
	"cyan": 0x00ffff, "darkblue": 0x00008b, "darkcyan": 0x008b8b, "darkgoldenrod": 0xb8860b, "darkgray": 0xa9a9a9,
	"darkgreen": 0x006400, "darkgrey": 0xa9a9a9, "darkkhaki": 0xbdb76b, "darkmagenta": 0x8b008b, "darkolivegreen": 0x556b2f,
	"darkorange": 0xff8c00, "darkorchid": 0x9932cc, "darkred": 0x8b0000, "darksalmon": 0xe9967a, "darkseagreen": 0x8fbc8f,
	"darkslateblue": 0x483d8b, "darkslategray": 0x2f4f4f, "darkslategrey": 0x2f4f4f, "darkturquoise": 0x00ced1,
	"darkviolet": 0x9400d3, "deeppink": 0xff1493, "deepskyblue": 0x00bfff, "dimgray": 0x696969, "dimgrey": 0x696969,
	"dodgerblue": 0x1e90ff, "firebrick": 0xb22222, "floralwhite": 0xfffaf0, "forestgreen": 0x228b22, "fuchsia": 0xff00ff,
	"gainsboro": 0xdcdcdc, "ghostwhite": 0xf8f8ff, "gold": 0xffd700, "goldenrod": 0xdaa520, "gray": 0x808080,
	"green": 0x008000, "greenyellow": 0xadff2f, "grey": 0x808080, "honeydew": 0xf0fff0, "hotpink": 0xff69b4,
	"indianred": 0xcd5c5c, "indigo": 0x4b0082, "ivory": 0xfffff0, "khaki": 0xf0e68c, "lavender": 0xe6e6fa,
	"lavenderblush": 0xfff0f5, "lawngreen": 0x7cfc00, "lemonchiffon": 0xfffacd, "lightblue": 0xadd8e6, "lightcoral": 0xf08080,
	"lightcyan": 0xe0ffff, "lightgoldenrodyellow": 0xfafad2, "lightgray": 0xd3d3d3, "lightgreen": 0x90ee90, "lightgrey": 0xd3d3d3,
	"lightpink": 0xffb6c1, "lightsalmon": 0xffa07a, "lightseagreen": 0x20b2aa, "lightskyblue": 0x87cefa,
	"lightslategray": 0x778899, "lightslategrey": 0x778899, "lightsteelblue": 0xb0c4de, "lightyellow": 0xffffe0,
	"lime": 0x00ff00, "limegreen": 0x32cd32, "linen": 0xfaf0e6, "magenta": 0xff00ff, "maroon": 0x800000,
	"mediumaquamarine": 0x66cdaa, "mediumblue": 0x0000cd, "mediumorchid": 0xba55d3, "mediumpurple": 0x9370db,
	"mediumseagreen": 0x3cb371, "mediumslateblue": 0x7b68ee, "mediumspringgreen": 0x00fa9a, "mediumturquoise": 0x48d1cc,
	"mediumvioletred": 0xc71585, "midnightblue": 0x191970, "mintcream": 0xf5fffa, "mistyrose": 0xffe4e1, "moccasin": 0xffe4b5,
	"navajowhite": 0xffdead, "navy": 0x000080, "oldlace": 0xfdf5e6, "olive": 0x808000, "olivedrab": 0x6b8e23,
	"orange": 0xffa500, "orangered": 0xff4500, "orchid": 0xda70d6, "palegoldenrod": 0xeee8aa, "palegreen": 0x98fb98,
	"paleturquoise": 0xafeeee, "palevioletred": 0xdb7093, "papayawhip": 0xffefd5, "peachpuff": 0xffdab9, "peru": 0xcd853f,
	"pink": 0xffc0cb, "plum": 0xdda0dd, "powderblue": 0xb0e0e6, "purple": 0x800080, "rebeccapurple": 0x663399,
	"red": 0xff0000, "rosybrown": 0xbc8f8f, "royalblue": 0x4169e1, "saddlebrown": 0x8b4513, "salmon": 0xfa8072,
	"sandybrown": 0xf4a460, "seagreen": 0x2e8b57, "seashell": 0xfff5ee, "sienna": 0xa0522d, "silver": 0xc0c0c0,
	"skyblue": 0x87ceeb, "slateblue": 0x6a5acd, "slategray": 0x708090, "slategrey": 0x708090, "snow": 0xfffafa,
	"springgreen": 0x00ff7f, "steelblue": 0x4682b4, "tan": 0xd2b48c, "teal": 0x008080, "thistle": 0xd8bfd8,
	"tomato": 0xff6347, "turquoise": 0x40e0d0, "violet": 0xee82ee, "wheat": 0xf5deb3, "white": 0xffffff,
	"whitesmoke": 0xf5f5f5, "yellow": 0xffff00, "yellowgreen": 0x9acd32,
}
