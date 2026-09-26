package psd

import (
	"encoding/binary"
	"image"
	"math"
)

// toNRGBA converts the colour planes of a mode (and an alpha plane, or nil)
// into RGB. Colour profiles are not applied: RGB and grayscale samples are
// taken as sRGB, CMYK is converted without ink profiles, and duotone and
// multichannel images show their first channel as gray.
func (f *file) toNRGBA(planes [][]uint8, alpha []uint8, w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	n := w * h
	var trans = -1
	if b := f.resources[resTransparency]; len(b) >= 2 && f.hdr.mode == modeIndexed {
		trans = int(binary.BigEndian.Uint16(b))
	}
	for i := 0; i < n; i++ {
		var r, g, b uint8
		a := uint8(255)
		switch f.hdr.mode {
		case modeRGB:
			r, g, b = planes[0][i], planes[1][i], planes[2][i]
		case modeCMYK:
			// Samples are stored inverted: 255 is no ink.
			k := uint32(planes[3][i])
			r = uint8((uint32(planes[0][i])*k + 127) / 255)
			g = uint8((uint32(planes[1][i])*k + 127) / 255)
			b = uint8((uint32(planes[2][i])*k + 127) / 255)
		case modeLab:
			r, g, b = labToRGB(planes[0][i], planes[1][i], planes[2][i])
		case modeIndexed:
			c := int(planes[0][i])
			if len(f.palette) >= 768 {
				r, g, b = f.palette[c], f.palette[256+c], f.palette[512+c]
			}
			if c == trans {
				a = 0
			}
		default:
			r = planes[0][i]
			g, b = r, r
		}
		if alpha != nil {
			a = alpha[i]
		}
		img.Pix[4*i], img.Pix[4*i+1], img.Pix[4*i+2], img.Pix[4*i+3] = r, g, b, a
	}
	return img
}

// unmatte undoes the white matte of a composite with transparency:
// Photoshop stores its colours blended over white, the colour of no ink or
// full lightness (and neutral a and b in Lab).
func (f *file) unmatte(planes [][]uint8, alpha []uint8) {
	for c, p := range planes {
		m := 255.0
		if f.hdr.mode == modeLab && c > 0 {
			m = 128
		}
		for i, v := range p {
			a := float64(alpha[i]) / 255
			switch {
			case a == 0:
				p[i] = uint8(m)
			case a < 1:
				p[i] = clamp8((float64(v) - m*(1-a)) / a)
			}
		}
	}
}

func clamp8(v float64) uint8 {
	switch {
	case v <= 0:
		return 0
	case v >= 255:
		return 255
	}
	return uint8(math.Round(v))
}

// labToRGB converts CIELAB (D50, as Photoshop stores it: L scaled to 0-255,
// a and b offset by 128) to sRGB.
func labToRGB(l8, a8, b8 uint8) (uint8, uint8, uint8) {
	l := float64(l8) * 100 / 255
	a, b := float64(a8)-128, float64(b8)-128
	fy := (l + 16) / 116
	fx := fy + a/500
	fz := fy - b/200
	inv := func(t float64) float64 {
		if t > 6.0/29 {
			return t * t * t
		}
		return 3 * (6.0 / 29) * (6.0 / 29) * (t - 4.0/29)
	}
	// D50 white
	x, y, z := 0.9642*inv(fx), inv(fy), 0.8249*inv(fz)
	// XYZ (D50) to linear sRGB, with Bradford adaptation to D65
	rl := 3.1338561*x - 1.6168667*y - 0.4906146*z
	gl := -0.9787684*x + 1.9161415*y + 0.0334540*z
	bl := 0.0719453*x - 0.2289914*y + 1.4052427*z
	return floatSample(rl, true), floatSample(gl, true), floatSample(bl, true)
}
