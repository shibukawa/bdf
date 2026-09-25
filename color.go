package bdf

import "fmt"

// Color is a packed 0xRRGGBBAA value.
type Color uint32

// RGB builds an opaque color.
func RGB(r, g, b uint8) Color { return RGBA(r, g, b, 255) }

// RGBA builds a color with alpha.
func RGBA(r, g, b, a uint8) Color {
	return Color(uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a))
}

// CSS returns the color as a CSS color string.
func (c Color) CSS() string {
	r, g, b, a := uint8(c>>24), uint8(c>>16), uint8(c>>8), uint8(c)
	if a == 255 {
		return fmt.Sprintf("#%02x%02x%02x", r, g, b)
	}
	return fmt.Sprintf("rgba(%d,%d,%d,%.4g)", r, g, b, float64(a)/255)
}
