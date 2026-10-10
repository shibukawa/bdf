package imagebdf

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"math"

	"github.com/shibukawa/bdf/image/imgconv"
)

// picture is a decoded image with its mipmaps: level k is the image scaled
// down 2^k times, premultiplied.
type picture struct {
	w, h   int
	levels []*image.RGBA

	// used is the image drawn that used the picture last, and drop takes
	// it from where a renderer keeps it (see Renderer.keep)
	used int
	drop func()
}

// decodePicture decodes a bitmap image part. AVIF images, and images of
// more than imgconv.MaxDecodePixels pixels (4 bytes a pixel: more memory
// than a preview is worth), are not decoded (ok is false).
func decodePicture(data []byte) (*picture, string, bool) {
	img, err := imgconv.Decode(data)
	if errors.Is(err, imgconv.ErrTooLarge) {
		return nil, "too large", false
	}
	if err != nil {
		return nil, imgconv.Sniff(data), false
	}
	// browsers turn JPEG and PNG images by their EXIF orientation
	img = imgconv.Orient(img, imgconv.Orientation(data))
	b := img.Bounds()
	rgba, ok := img.(*image.RGBA)
	if !ok || b.Min != (image.Point{}) {
		rgba = image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(rgba, rgba.Bounds(), img, b.Min, draw.Src)
	}
	if rgba.Bounds().Empty() {
		return nil, imgconv.Sniff(data), false
	}
	return &picture{w: b.Dx(), h: b.Dy(), levels: []*image.RGBA{rgba}}, "", true
}

// isSVG reports whether image bytes are SVG markup (after a byte order
// mark and white space) rather than an encoded bitmap, as the viewer tells.
func isSVG(b []byte) bool {
	if len(b) >= 2 && (b[0] == 0xff && b[1] == 0xfe || b[0] == 0xfe && b[1] == 0xff) {
		return true // UTF-16
	}
	b = bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	b = bytes.TrimLeft(b, " \t\r\n")
	return len(b) > 0 && b[0] == '<'
}

// levelFor picks the mip level for sampling through inv (device → image
// pixels): the one where a device pixel spans at most about two image
// pixels.
func (p *picture) levelFor(inv matrix) int {
	f := inv.maxScale()
	if !(f > 2) {
		return 0
	}
	// bounded as a float64: a scale no int holds is that of an image
	// drawn in no place at all
	l := int(math.Min(math.Floor(math.Log2(f)), 31))
	for len(p.levels) <= l {
		prev := p.levels[len(p.levels)-1]
		if prev.Rect.Dx() <= 1 && prev.Rect.Dy() <= 1 {
			return len(p.levels) - 1
		}
		p.levels = append(p.levels, halve(prev))
	}
	return l
}

// halve scales an image down by two with a box filter.
func halve(src *image.RGBA) *image.RGBA {
	sw, sh := src.Rect.Dx(), src.Rect.Dy()
	w, h := max(1, (sw+1)/2), max(1, (sh+1)/2)
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var sum [4]int
			n := 0
			for dy := 0; dy < 2; dy++ {
				for dx := 0; dx < 2; dx++ {
					sx, sy := 2*x+dx, 2*y+dy
					if sx >= sw || sy >= sh {
						continue
					}
					o := sy*src.Stride + 4*sx
					for c := 0; c < 4; c++ {
						sum[c] += int(src.Pix[o+c])
					}
					n++
				}
			}
			o := y*dst.Stride + 4*x
			for c := 0; c < 4; c++ {
				dst.Pix[o+c] = uint8((sum[c] + n/2) / n)
			}
		}
	}
	return dst
}

// sample reads the colour at (x, y) in full-size image pixels from a mip
// level, interpolating bilinearly; wrapX and wrapY repeat the image instead
// of clamping at its edges.
func (p *picture) sample(level int, x, y float64, wrapX, wrapY bool) [4]float32 {
	img := p.levels[min(level, len(p.levels)-1)]
	w, h := img.Rect.Dx(), img.Rect.Dy()
	x = x*float64(w)/float64(p.w) - 0.5
	y = y*float64(h)/float64(p.h) - 0.5
	x0, y0 := math.Floor(x), math.Floor(y)
	fx, fy := float32(x-x0), float32(y-y0)
	ix, iy := int(x0), int(y0)
	at := func(px, py int) [4]float32 {
		if wrapX {
			px = ((px % w) + w) % w
		} else {
			px = min(max(px, 0), w-1)
		}
		if wrapY {
			py = ((py % h) + h) % h
		} else {
			py = min(max(py, 0), h-1)
		}
		o := py*img.Stride + 4*px
		return [4]float32{float32(img.Pix[o]) / 255, float32(img.Pix[o+1]) / 255, float32(img.Pix[o+2]) / 255, float32(img.Pix[o+3]) / 255}
	}
	a, b, c, d := at(ix, iy), at(ix+1, iy), at(ix, iy+1), at(ix+1, iy+1)
	var out [4]float32
	for k := 0; k < 4; k++ {
		top := a[k] + (b[k]-a[k])*fx
		bot := c[k] + (d[k]-c[k])*fx
		out[k] = top + (bot-top)*fy
	}
	return out
}

// nearest reads the pixel under (x, y) of the full-size image.
func (p *picture) nearest(x, y float64) [4]float32 {
	img := p.levels[0]
	px := min(max(int(math.Floor(x)), 0), p.w-1)
	py := min(max(int(math.Floor(y)), 0), p.h-1)
	o := py*img.Stride + 4*px
	return [4]float32{float32(img.Pix[o]) / 255, float32(img.Pix[o+1]) / 255, float32(img.Pix[o+2]) / 255, float32(img.Pix[o+3]) / 255}
}
