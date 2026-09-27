package equation

import (
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// colorNames are the color names of LaTeX's xcolor and the CSS colors
// formulas use.
var colorNames = map[string]uint32{
	"black": 0x000000, "white": 0xffffff, "red": 0xff0000, "green": 0x00ff00, "blue": 0x0000ff,
	"cyan": 0x00ffff, "magenta": 0xff00ff, "yellow": 0xffff00, "gray": 0x808080, "grey": 0x808080,
	"darkgray": 0x404040, "lightgray": 0xbfbfbf, "brown": 0xbf8040, "lime": 0xbfff00, "olive": 0x808000,
	"orange": 0xff8000, "pink": 0xffbfbf, "purple": 0xbf0040, "teal": 0x008080, "violet": 0x800080,
	"maroon": 0x800000, "navy": 0x000080, "silver": 0xc0c0c0, "fuchsia": 0xff00ff, "aqua": 0x00ffff,
	"darkgreen": 0x006400, "darkblue": 0x00008b, "darkred": 0x8b0000, "gold": 0xffd700, "indigo": 0x4b0082,
}

// ParseColor reads a color: a name, #rgb, #rrggbb or rgb(r, g, b).
func ParseColor(s string) (bdf.Color, bool) {
	s = strings.ToLower(strings.TrimSpace(s))
	if v, ok := colorNames[s]; ok {
		return bdf.RGBA(uint8(v>>16), uint8(v>>8), uint8(v), 255), true
	}
	if strings.HasPrefix(s, "#") {
		h := s[1:]
		if len(h) == 3 {
			h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
		}
		if v, err := strconv.ParseUint(h, 16, 32); err == nil && len(h) == 6 {
			return bdf.RGBA(uint8(v>>16), uint8(v>>8), uint8(v), 255), true
		}
		return 0, false
	}
	if strings.HasPrefix(s, "rgb(") && strings.HasSuffix(s, ")") {
		parts := strings.Split(s[4:len(s)-1], ",")
		if len(parts) != 3 {
			return 0, false
		}
		var c [3]uint8
		for i, p := range parts {
			v, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil {
				return 0, false
			}
			c[i] = uint8(max(0, min(255, v)))
		}
		return bdf.RGBA(c[0], c[1], c[2], 255), true
	}
	return 0, false
}
