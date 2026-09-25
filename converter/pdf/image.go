package pdf

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"

	"github.com/shibukawa/bdf/imgconv"

	"github.com/pdfcpu/pdfcpu/pkg/filter"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
)

// decodedImage is an image ready to be stored as a part.
type decodedImage struct {
	data   []byte // encoded bytes (PNG/JPEG, or WebP/AVIF after conversion)
	format string
	w, h   int
	isMask bool
}

// loadImage converts an image XObject (or inline image dict+data) into PNG/JPEG bytes.
// fill is the current fill colour, used for stencil masks.
func (c *converter) loadImage(d types.Dict, raw []byte, filters []types.PDFFilter, res types.Dict, fill bdf.Color) (*decodedImage, error) {
	p := c.pdf
	w := p.intOr(d["Width"], p.intOr(d["W"], 0))
	h := p.intOr(d["Height"], p.intOr(d["H"], 0))
	if w <= 0 || h <= 0 {
		return nil, errf("image without size")
	}
	bpc := p.intOr(d["BitsPerComponent"], p.intOr(d["BPC"], 8))
	isMask := p.boolOr(d["ImageMask"], p.boolOr(d["IM"], false))
	decodeArr := p.nums(d["Decode"])
	if decodeArr == nil {
		decodeArr = p.nums(d["D"])
	}

	// Apply the non-image filters; find out whether an image codec remains.
	sd := &types.StreamDict{Raw: raw, FilterPipeline: filters}
	data, codec, err := p.decodeStream(sd)
	if err != nil {
		return nil, err
	}
	switch codec {
	case filter.DCT:
		if !isMask && d["SMask"] == nil && d["Mask"] == nil && len(decodeArr) == 0 {
			// Pass JPEG bytes through untouched (or let imgconv try a smaller encoding).
			if cfg, err := jpeg.DecodeConfig(bytes.NewReader(data)); err == nil {
				return c.storeEncoded(data, cfg.Width, cfg.Height), nil
			}
		}
		img, err := jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		nrgba := toNRGBA(img)
		c.applyMasks(nrgba, d, res)
		return c.storePixels(nrgba, false), nil
	case filter.JPX:
		return nil, errf("JPXDecode images are not supported")
	case filter.JBIG2:
		return nil, errf("JBIG2Decode images are not supported")
	case filter.CCITTFax:
		parms := map[string]int{"Columns": w, "Rows": h}
		var dp types.Dict
		for _, f := range filters {
			if f.Name == filter.CCITTFax {
				dp = f.DecodeParms
			}
		}
		for k, v := range dp {
			if n, ok := p.num(v); ok {
				parms[k] = int(n)
			}
			if b, ok := p.deref(v).(types.Boolean); ok && bool(b) {
				parms[k] = 1
			}
		}
		fi, err := filter.NewFilter(filter.CCITTFax, parms)
		if err != nil {
			return nil, err
		}
		r, err := fi.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		if data, err = readAll(r); err != nil {
			return nil, err
		}
		bpc = 1
		if !isMask && d["ColorSpace"] == nil {
			d = d.Clone().(types.Dict)
			d["ColorSpace"] = types.Name("DeviceGray")
		}
		// CCITT decoders produce 1 = black by default (BlackIs1 false means 0 bits are black in the encoded data,
		// but the decoded output follows the DeviceGray convention where 0 is black).
	}

	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	if isMask {
		// Stencil mask: sample 0 paints (unless Decode [1 0]).
		paintOnZero := true
		if len(decodeArr) >= 1 && decodeArr[0] == 1 {
			paintOnZero = false
		}
		rowBytes := (w + 7) / 8
		r, g, b := uint8(fill>>24), uint8(fill>>16), uint8(fill>>8)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				i := y*rowBytes + x/8
				bit := uint8(0)
				if i < len(data) {
					bit = data[i] >> (7 - uint(x%8)) & 1
				}
				paint := (bit == 0) == paintOnZero
				if paint {
					out.SetNRGBA(x, y, color.NRGBA{r, g, b, 255})
				}
			}
		}
		di := c.storePixels(out, true)
		di.isMask = true
		return di, nil
	}

	csObj := d["ColorSpace"]
	if csObj == nil {
		csObj = d["CS"]
	}
	cs := c.loadColorSpace(csObj, res)
	n := cs.n
	if cs.family == "Pattern" || n == 0 {
		return nil, errf("image with unsupported colour space")
	}
	maxv := float64(int(1)<<uint(bpc) - 1)
	rowBits := w * n * bpc
	rowBytes := (rowBits + 7) / 8
	comps := make([]float64, n)
	// Decode array handling for the common cases.
	dec := decodeArr
	if len(dec) != 2*n {
		dec = nil
	}
	sample := func(row []byte, i int) float64 {
		// i-th sample in the row
		switch bpc {
		case 8:
			if i < len(row) {
				return float64(row[i])
			}
		case 16:
			if 2*i+1 < len(row) {
				return float64(uint16(row[2*i])<<8 | uint16(row[2*i+1]))
			}
		default:
			bit := i * bpc
			var v uint32
			for k := 0; k < bpc; k++ {
				b := bit + k
				if b/8 >= len(row) {
					return 0
				}
				v = v<<1 | uint32(row[b/8]>>(7-uint(b%8))&1)
			}
			return float64(v)
		}
		return 0
	}
	for y := 0; y < h; y++ {
		start := y * rowBytes
		if start >= len(data) {
			break
		}
		row := data[start:min(len(data), start+rowBytes)]
		for x := 0; x < w; x++ {
			for k := 0; k < n; k++ {
				raw := sample(row, x*n+k)
				var v float64
				switch {
				case cs.family == "Indexed":
					v = raw
					if dec != nil {
						v = dec[0] + raw*(dec[1]-dec[0])/maxv
					}
				case cs.family == "Lab":
					v = raw / maxv
				default:
					v = raw / maxv
					if dec != nil {
						v = dec[2*k] + v*(dec[2*k+1]-dec[2*k])
					}
				}
				comps[k] = v
			}
			var r, g, b float64
			if cs.family == "Lab" {
				r, g, b = cs.rgbf(cs.decodeLabBytes(comps))
			} else {
				r, g, b = cs.rgbf(comps)
			}
			out.SetNRGBA(x, y, color.NRGBA{to8(r), to8(g), to8(b), 255})
		}
	}
	c.applyMasks(out, d, res)
	return c.storePixels(out, true), nil
}

// storePixels encodes decoded pixels according to the image options.
func (c *converter) storePixels(img *image.NRGBA, lossless bool) *decodedImage {
	r, err := imgconv.EncodeImage(img, lossless, c.opts.Images)
	if err != nil {
		c.warnOnce("imgconv", "image conversion: %v", err)
	}
	b := img.Bounds()
	return &decodedImage{data: r.Data, format: r.Format, w: b.Dx(), h: b.Dy()}
}

// storeEncoded keeps already encoded bytes, letting imgconv shrink them in Convert mode.
func (c *converter) storeEncoded(data []byte, w, h int) *decodedImage {
	r, err := imgconv.Optimize(data, c.opts.Images)
	if err != nil {
		c.warnOnce("imgconv", "image conversion: %v", err)
	}
	return &decodedImage{data: r.Data, format: r.Format, w: w, h: h}
}

// applyMasks applies /SMask (soft mask) or /Mask (stencil or colour key) to img.
func (c *converter) applyMasks(img *image.NRGBA, d types.Dict, res types.Dict) {
	p := c.pdf
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if sm := p.stream(d["SMask"]); sm != nil {
		mask, err := c.loadImage(sm.Dict, sm.Raw, sm.FilterPipeline, res, 0)
		if err == nil {
			if mimg, err := decodeStored(mask); err == nil {
				for y := 0; y < h; y++ {
					for x := 0; x < w; x++ {
						mx, my := x*mask.w/w, y*mask.h/h
						r, _, _, _ := mimg.At(mx, my).RGBA()
						i := img.PixOffset(x, y)
						img.Pix[i+3] = uint8(r >> 8)
					}
				}
			}
		} else {
			c.warnf("soft mask: %v", err)
		}
		return
	}
	switch m := p.deref(d["Mask"]).(type) {
	case types.StreamDict, *types.StreamDict:
		sd := p.stream(d["Mask"])
		mask, err := c.loadImage(sd.Dict, sd.Raw, sd.FilterPipeline, res, bdf.RGB(0, 0, 0))
		if err == nil {
			if mimg, err := decodeStored(mask); err == nil {
				for y := 0; y < h; y++ {
					for x := 0; x < w; x++ {
						mx, my := x*mask.w/w, y*mask.h/h
						_, _, _, a := mimg.At(mx, my).RGBA()
						// Stencil painted (a>0) means masked OUT.
						i := img.PixOffset(x, y)
						if a > 0 {
							img.Pix[i+3] = 0
						}
					}
				}
			}
		}
	case types.Array:
		// Colour key masking is applied on decoded samples; approximate by
		// comparing the final colour against the key converted through the space.
		_ = m
	}
}

func toNRGBA(img image.Image) *image.NRGBA {
	if n, ok := img.(*image.NRGBA); ok {
		return n
	}
	b := img.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			out.Set(x, y, img.At(b.Min.X+x, b.Min.Y+y))
		}
	}
	return out
}

// decodeStored decodes an image produced by storePixels (PNG or WebP).
func decodeStored(di *decodedImage) (image.Image, error) {
	switch di.format {
	case "png":
		return png.Decode(bytes.NewReader(di.data))
	case "jpeg":
		return jpeg.Decode(bytes.NewReader(di.data))
	}
	return imgconv.Decode(di.data)
}

type convError string

func (e convError) Error() string { return string(e) }

func errf(s string) error { return convError(s) }
