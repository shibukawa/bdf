package svg

import (
	"math"
	"strconv"
	"strings"
)

// Color is a colour: red, green, blue and alpha from 0 to 1, not
// premultiplied.
type Color struct {
	R, G, B, A float64
}

// ParseColor reads a CSS colour: a name, #rgb, #rgba, #rrggbb, #rrggbbaa,
// rgb(), rgba(), hsl(), hsla(), transparent, or currentColor, which is
// current. It reports false for what is none of them (a paint server such
// as url(#id), or none).
func ParseColor(s string, current Color) (Color, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	switch {
	case s == "currentcolor":
		return current, true
	case s == "transparent":
		return Color{}, true
	case strings.HasPrefix(s, "#"):
		h := s[1:]
		nib := func(i int) float64 {
			v, _ := strconv.ParseUint(h[i:i+1], 16, 8)
			return float64(v*17) / 255
		}
		byt := func(i int) float64 {
			v, _ := strconv.ParseUint(h[i:i+2], 16, 8)
			return float64(v) / 255
		}
		if _, err := strconv.ParseUint(h, 16, 64); err != nil {
			return Color{}, false
		}
		switch len(h) {
		case 3:
			return Color{R: nib(0), G: nib(1), B: nib(2), A: 1}, true
		case 4:
			return Color{R: nib(0), G: nib(1), B: nib(2), A: nib(3)}, true
		case 6:
			return Color{R: byt(0), G: byt(2), B: byt(4), A: 1}, true
		case 8:
			return Color{R: byt(0), G: byt(2), B: byt(4), A: byt(6)}, true
		}
		return Color{}, false
	case strings.HasPrefix(s, "rgb"), strings.HasPrefix(s, "hsl"):
		open, close := strings.IndexByte(s, '('), strings.LastIndexByte(s, ')')
		if open < 0 || close < open {
			return Color{}, false
		}
		parts := strings.FieldsFunc(s[open+1:close], func(r rune) bool { return r == ',' || r == ' ' || r == '/' })
		if len(parts) < 3 {
			return Color{}, false
		}
		num := func(p string, scale float64) float64 {
			if strings.HasSuffix(p, "%") {
				v, _ := strconv.ParseFloat(strings.TrimSuffix(p, "%"), 64)
				return v / 100
			}
			v, _ := strconv.ParseFloat(strings.TrimSuffix(p, "deg"), 64)
			return v / scale
		}
		a := 1.0
		if len(parts) >= 4 {
			a = num(parts[3], 1)
		}
		clamp := func(v float64) float64 { return math.Min(1, math.Max(0, v)) }
		if strings.HasPrefix(s, "rgb") {
			return Color{R: clamp(num(parts[0], 255)), G: clamp(num(parts[1], 255)), B: clamp(num(parts[2], 255)), A: clamp(a)}, true
		}
		hh := math.Mod(num(parts[0], 1), 360) / 360
		if hh < 0 {
			hh++
		}
		sat, l := clamp(num(parts[1], 100)), clamp(num(parts[2], 100))
		r, g, b := hslToRGB(hh, sat, l)
		return Color{R: r, G: g, B: b, A: clamp(a)}, true
	}
	if c, ok := namedColors[s]; ok {
		return Color{R: float64(c>>16) / 255, G: float64(c>>8&0xff) / 255, B: float64(c&0xff) / 255, A: 1}, true
	}
	return Color{}, false
}

func hslToRGB(h, s, l float64) (float64, float64, float64) {
	if s == 0 {
		return l, l, l
	}
	q := l * (1 + s)
	if l >= 0.5 {
		q = l + s - l*s
	}
	p := 2*l - q
	hue := func(t float64) float64 {
		if t < 0 {
			t++
		}
		if t > 1 {
			t--
		}
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
	return hue(h + 1.0/3), hue(h), hue(h - 1.0/3)
}
