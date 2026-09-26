package pdf

import (
	"image"
	"image/color"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf/converter/internal/jpx"
)

// loadJPX decodes a JPXDecode image. The image dictionary's colour space
// wins over the JP2 header's (PDF 32000-1 §7.4.9), BitsPerComponent and
// Decode do not apply, and with /SMaskInData the JPEG 2000 opacity channel
// is the soft mask.
func (c *converter) loadJPX(d types.Dict, data []byte, res types.Dict) (*decodedImage, error) {
	p := c.pdf
	img, err := jpx.Decode(data)
	if err != nil {
		return nil, err
	}
	w, h := img.Width, img.Height
	var colour []jpx.Component
	for i, comp := range img.Components {
		if i != img.Alpha {
			colour = append(colour, comp)
		}
	}
	var cs *colorSpace
	ycc := false
	if csObj := d["ColorSpace"]; csObj != nil {
		cs = c.loadColorSpace(csObj, res)
		if cs.family == "Pattern" || cs.n == 0 || cs.n > len(colour) {
			cs = nil // a palette the decoder applied, or a mismatch: go by the data
		}
	}
	if cs == nil {
		switch {
		case img.ColorSpace == 17 || len(colour) < 3:
			cs = csGray
		case img.ColorSpace == 12 || len(colour) == 4 && img.ColorSpace != 16 && img.ColorSpace != 18:
			cs = csCMYK
		default:
			cs = csRGB
			ycc = img.ColorSpace == 18
		}
	}
	if len(colour) < cs.n {
		return nil, errf("JPX image with too few components")
	}
	smaskInData := p.intOr(d["SMaskInData"], 0)
	alpha := img.Alpha >= 0 && smaskInData != 0 && d["SMask"] == nil
	premul := alpha && (img.Premultiplied || smaskInData == 2)
	norm := func(comp *jpx.Component, i int) float64 {
		v := float64(comp.Data[i])
		if comp.Signed {
			v += float64(int64(1) << (comp.Precision - 1))
		}
		return v / float64(int64(1)<<comp.Precision-1)
	}
	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	comps := make([]float64, cs.n)
	for i := 0; i < w*h; i++ {
		a := 1.0
		if alpha {
			a = norm(&img.Components[img.Alpha], i)
		}
		for k := range comps {
			if cs.family == "Indexed" {
				comps[k] = float64(colour[k].Data[i]) // palette indices, not intensities
				continue
			}
			comps[k] = norm(&colour[k], i)
			if premul && a > 0 {
				comps[k] /= a
			}
		}
		var r, g, b float64
		if ycc {
			y, cb, cr := comps[0], comps[1]-0.5, comps[2]-0.5
			r, g, b = y+1.402*cr, y-0.344136*cb-0.714136*cr, y+1.772*cb
		} else {
			r, g, b = cs.rgbf(comps)
		}
		out.SetNRGBA(i%w, i/w, color.NRGBA{to8(r), to8(g), to8(b), to8(a)})
	}
	c.applyMasks(out, d, res)
	return c.storePixels(out, true), nil
}
