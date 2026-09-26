package imgconv

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
)

// EncodePixels stores decoded pixels of any image type. lossless says
// whether they come from a lossless source.
//
// In Keep mode, pixels from a lossless source become PNG in the smallest
// colour type that holds them (1 bit per pixel for a two-colour image,
// greyscale, a palette of up to 256 colours), and pixels from a lossy
// source become JPEG (at Options.Quality) unless they are transparent. In
// Convert mode, WebP replaces that when it is smaller, as in EncodeImage.
func EncodePixels(img image.Image, lossless bool, opts Options) (Result, error) {
	keep := encodeKeep(img, lossless, opts)
	if opts.Mode != Convert {
		return keep, nil
	}
	if !Available() {
		return keep, ErrNotAvailable
	}
	best, err := encodeBest(compactNRGBA(img), lossless, opts, len(keep.Data))
	if err != nil {
		return keep, err
	}
	if best == nil {
		return keep, nil
	}
	best.Converted = true
	return *best, nil
}

func encodeKeep(img image.Image, lossless bool, opts Options) Result {
	if !lossless && opaque(img) {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: opts.quality()}); err == nil {
			return Result{Data: buf.Bytes(), Format: "jpeg"}
		}
	}
	var buf bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	_ = enc.Encode(&buf, smallest(img))
	return Result{Data: buf.Bytes(), Format: "png", Lossless: true}
}

// smallest returns img in the image type the PNG encoder writes with the
// fewest bits per pixel: a two-colour palette (1 bit), grey, or a palette.
func smallest(img image.Image) image.Image {
	switch m := img.(type) {
	case *image.Paletted:
		return m
	case *image.Gray:
		b := m.Rect
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for _, v := range m.Pix[m.PixOffset(b.Min.X, y):][:b.Dx()] {
				if v != 0 && v != 255 {
					return m
				}
			}
		}
		out := image.NewPaletted(image.Rect(0, 0, b.Dx(), b.Dy()), color.Palette{color.Gray{0}, color.Gray{255}})
		for y := b.Min.Y; y < b.Max.Y; y++ {
			row := out.Pix[(y-b.Min.Y)*out.Stride:]
			for x, v := range m.Pix[m.PixOffset(b.Min.X, y):][:b.Dx()] {
				row[x] = v >> 7
			}
		}
		return out
	}
	n := compactNRGBA(img)
	gray := true
	for i := 0; i < len(n.Pix); i += 4 {
		if n.Pix[i+3] != 255 || n.Pix[i] != n.Pix[i+1] || n.Pix[i+1] != n.Pix[i+2] {
			gray = false
			break
		}
	}
	if gray {
		g := image.NewGray(n.Rect)
		for i := range g.Pix {
			g.Pix[i] = n.Pix[4*i]
		}
		return smallest(g)
	}
	// Up to 256 colours: a palette, in the order the colours appear.
	pal := map[uint32]uint8{}
	idx := make([]uint8, len(n.Pix)/4)
	for i := range idx {
		j := 4 * i
		c := uint32(n.Pix[j])<<24 | uint32(n.Pix[j+1])<<16 | uint32(n.Pix[j+2])<<8 | uint32(n.Pix[j+3])
		k, ok := pal[c]
		if !ok {
			if len(pal) == 256 {
				return n
			}
			k = uint8(len(pal))
			pal[c] = k
		}
		idx[i] = k
	}
	p := &image.Paletted{Pix: idx, Stride: n.Rect.Dx(), Rect: n.Rect, Palette: make(color.Palette, len(pal))}
	for c, k := range pal {
		p.Palette[k] = color.NRGBA{uint8(c >> 24), uint8(c >> 16), uint8(c >> 8), uint8(c)}
	}
	return p
}

// compactNRGBA returns img as NRGBA pixels whose rows follow each other,
// with the origin at (0, 0).
func compactNRGBA(img image.Image) *image.NRGBA {
	if n, ok := img.(*image.NRGBA); ok && n.Rect.Min == (image.Point{}) && n.Stride == 4*n.Rect.Dx() {
		return n
	}
	b := img.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), img, b.Min, draw.Src)
	return out
}

// opaque reports whether every pixel of img is opaque.
func opaque(img image.Image) bool {
	if o, ok := img.(interface{ Opaque() bool }); ok {
		return o.Opaque()
	}
	return false
}
