package psbdf

import (
	"bytes"
	"compress/zlib"
	"encoding/ascii85"
	"fmt"
	"image"
	"image/draw"
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/imgconv"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// psImage is an image of the document as PostScript draws it: its data,
// deflated (or a JPEG file as it is) and in ASCII85, which is written
// where it is drawn.
type psImage struct {
	// the size of the image in the units that IMAGE_SUB counts in: its
	// pixels, or the natural size of an SVG image
	w, h float64
	// opaque is set for an image that PostScript draws: one without alpha
	opaque bool
	// the pixels of the data, gray or RGB, and its filter
	pw, ph int
	gray   bool
	filter string
	data   []byte
}

// svgDPI is the resolution SVG images are drawn at, for the size they are
// first shown at.
const svgDPI = 300

// imageOf returns an image part as PostScript draws it, or nil for one that
// cannot be drawn. scale gives how many points a unit of an image of a
// size is shown as, where it is first used (an SVG image is drawn for it).
func (ex *exporter) imageOf(h bdf.Hash, scale func(w, h float64) float64) *psImage {
	if img, ok := ex.images[h]; ok {
		return img
	}
	ex.images[h] = nil
	part := ex.doc.Part(h)
	if part == nil {
		ex.warnf("image %s is missing", h)
		return nil
	}
	data := part.Data
	var out *psImage
	switch format := imgconv.Sniff(data); {
	case isSVG(data):
		px := func(w, h float64) float64 {
			k := scale(w, h) * svgDPI / 72
			if !(k > 0) || math.IsInf(k, 0) {
				k = 1
			}
			return k
		}
		pix, w, hh, warnings, err := imagebdf.SVGImage(data, px, &imagebdf.Options{FontFS: ex.opts.FontFS, FontDirs: ex.opts.FontDirs, NoSystemFonts: ex.opts.NoSystemFonts})
		for _, w := range warnings {
			ex.warnf("%s", w)
		}
		if err != nil {
			ex.warnf("an SVG image is not drawn: %v", err)
			return nil
		}
		out = ex.pixels(pix)
		out.w, out.h = w, hh
	case format == "jpeg" && imgconv.Orientation(data) <= 1:
		w, hh, comps, ok := jpegInfo(data)
		if !ok || comps == 3 && ex.opts.Gray {
			out = ex.decoded(data)
			break
		}
		// the file as it is: PostScript reads JPEG images
		var b bytes.Buffer
		writeA85(&b, data)
		out = &psImage{w: float64(w), h: float64(hh), opaque: true, pw: w, ph: hh, gray: comps == 1, filter: "/DCTDecode", data: b.Bytes()}
	default:
		out = ex.decoded(data)
	}
	ex.images[h] = out
	return out
}

// decoded returns an image of its pixels.
func (ex *exporter) decoded(data []byte) *psImage {
	img, err := imgconv.Decode(data)
	if err != nil {
		switch format := imgconv.Sniff(data); {
		case format == "avif":
			ex.warnf("AVIF images are not written")
		case err == imgconv.ErrTooLarge:
			ex.warnf("an image of more than %d pixels is not written", imgconv.MaxDecodePixels)
		default:
			ex.warnf("an image (%s) cannot be decoded", format)
		}
		return nil
	}
	// browsers turn JPEG and PNG images by their EXIF orientation
	img = imgconv.Orient(img, imgconv.Orientation(data))
	b := img.Bounds()
	if b.Empty() {
		return nil
	}
	pix, ok := img.(*image.NRGBA)
	if !ok || b.Min != (image.Point{}) {
		pix = image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(pix, pix.Bounds(), img, b.Min, draw.Src)
	}
	out := ex.pixels(pix)
	out.w, out.h = float64(b.Dx()), float64(b.Dy())
	return out
}

// pixels returns an image of pixels: its data when all are opaque (one
// with alpha is drawn as a picture of the page, and its data is not
// written), gray when all are gray.
func (ex *exporter) pixels(pix *image.NRGBA) *psImage {
	w, h := pix.Rect.Dx(), pix.Rect.Dy()
	out := &psImage{pw: w, ph: h, opaque: true, gray: true, filter: "/FlateDecode"}
	for y := 0; y < h && out.opaque; y++ {
		row := pix.Pix[y*pix.Stride : y*pix.Stride+4*w]
		for x := 0; x < 4*w; x += 4 {
			if row[x] != row[x+1] || row[x] != row[x+2] {
				out.gray = false
			}
			if row[x+3] != 0xff {
				out.opaque = false
				break
			}
		}
	}
	if !out.opaque {
		return out
	}
	out.gray = out.gray || ex.opts.Gray
	out.data = ex.encode(w, h, out.gray, func(y int, rgb []byte) {
		row := pix.Pix[y*pix.Stride:]
		for x := 0; x < w; x++ {
			copy(rgb[3*x:3*x+3], row[4*x:4*x+3])
		}
	})
	return out
}

// encode deflates the rows of an image, gray or RGB, as ASCII85: row gives
// the RGB pixels of a row.
func (ex *exporter) encode(w, h int, gray bool, row func(y int, rgb []byte)) []byte {
	var raw bytes.Buffer
	zw, _ := zlib.NewWriterLevel(&raw, zlib.BestSpeed)
	rgb := make([]byte, 3*w)
	g := make([]byte, w)
	for y := 0; y < h; y++ {
		row(y, rgb)
		if gray {
			for x := range g {
				r, gg, b := float64(rgb[3*x])/255, float64(rgb[3*x+1])/255, float64(rgb[3*x+2])/255
				g[x] = byte(math.Round(luma(r, gg, b) * 255))
			}
			zw.Write(g)
		} else {
			zw.Write(rgb)
		}
	}
	zw.Close()
	var b bytes.Buffer
	writeA85(&b, raw.Bytes())
	return b.Bytes()
}

// picture returns a picture of a part of a page, which the renderer drew
// on white.
func (ex *exporter) picture(img *image.RGBA) *psImage {
	w, h := img.Rect.Dx(), img.Rect.Dy()
	out := &psImage{w: float64(w), h: float64(h), pw: w, ph: h, opaque: true, gray: ex.opts.Gray, filter: "/FlateDecode"}
	out.data = ex.encode(w, h, out.gray, func(y int, rgb []byte) {
		row := img.Pix[y*img.Stride:]
		for x := 0; x < w; x++ {
			p := row[4*x : 4*x+4]
			// (premultiplied: over white)
			for i := 0; i < 3; i++ {
				rgb[3*x+i] = p[i] + (0xff - p[3])
			}
		}
	})
	return out
}

// writeImage writes an image drawn into the unit square under m, its first
// row at the top (y 0).
func writeImage(b *bytes.Buffer, img *psImage, m matrix, smooth bool) {
	cs, decode := "/DeviceRGB", "0 1 0 1 0 1"
	if img.gray {
		cs, decode = "/DeviceGray", "0 1"
	}
	interpolate := ""
	if smooth {
		interpolate = " /Interpolate true"
	}
	fmt.Fprintf(b, "q %s cm\n<< /ImageType 1 /Width %d /Height %d /BitsPerComponent 8 /Decode [%s] /ImageMatrix [%d 0 0 %d 0 0]%s >> %s %s bdfimage\n",
		m, img.pw, img.ph, decode, img.pw, img.ph, interpolate, cs, img.filter)
	b.Write(img.data)
	b.WriteString("Q\n")
}

// a85Line is the length of the lines of ASCII85 data.
const a85Line = 76

// writeA85 writes data in ASCII85, in lines, and its end (~>). A line does
// not start with %, which a reader of the comments of the file would take
// for one of them.
func writeA85(b *bytes.Buffer, data []byte) {
	enc := make([]byte, ascii85.MaxEncodedLen(len(data)))
	enc = enc[:ascii85.Encode(enc, data)]
	for len(enc) > 0 {
		n := min(a85Line, len(enc))
		if enc[0] == '%' {
			b.WriteByte(' ')
		}
		b.Write(enc[:n])
		b.WriteByte('\n')
		enc = enc[n:]
	}
	b.WriteString("~>\n")
}

// deflate compresses data for the FlateDecode filter.
func deflate(data []byte) []byte {
	var b bytes.Buffer
	zw, _ := zlib.NewWriterLevel(&b, zlib.BestCompression)
	zw.Write(data)
	zw.Close()
	return b.Bytes()
}

// isSVG reports whether image bytes are SVG markup (after a byte order
// mark and white space) rather than an encoded bitmap, as the viewer tells.
func isSVG(b []byte) bool {
	b = bytes.TrimPrefix(b, []byte("\xef\xbb\xbf"))
	b = bytes.TrimLeft(b, " \t\r\n")
	return len(b) > 0 && b[0] == '<'
}

// jpegInfo reads the size and the number of components of a JPEG image
// that the DCTDecode filter of every interpreter reads: baseline, of 8
// bits, gray or of three components.
func jpegInfo(data []byte) (w, h, comps int, ok bool) {
	if len(data) < 4 || data[0] != 0xff || data[1] != 0xd8 {
		return 0, 0, 0, false
	}
	for pos := 2; pos+4 <= len(data); {
		if data[pos] != 0xff {
			return 0, 0, 0, false
		}
		marker := data[pos+1]
		if marker == 0xff {
			pos++
			continue
		}
		if marker == 0xd8 || marker >= 0xd0 && marker <= 0xd7 || marker == 0x01 {
			pos += 2
			continue
		}
		n := int(data[pos+2])<<8 | int(data[pos+3])
		if n < 2 || pos+2+n > len(data) {
			return 0, 0, 0, false
		}
		switch marker {
		case 0xc0, 0xc1: // baseline and extended
			if n < 8 || data[pos+4] != 8 {
				return 0, 0, 0, false
			}
			h = int(data[pos+5])<<8 | int(data[pos+6])
			w = int(data[pos+7])<<8 | int(data[pos+8])
			comps = int(data[pos+9])
			return w, h, comps, w > 0 && h > 0 && (comps == 1 || comps == 3)
		case 0xc2, 0xc3, 0xc5, 0xc6, 0xc7, 0xc9, 0xca, 0xcb, 0xcd, 0xce, 0xcf, 0xda, 0xd9:
			return 0, 0, 0, false
		}
		pos += 2 + n
	}
	return 0, 0, 0, false
}

// image draws an image, or the part src (x, y, w, h in the units of the
// image) of it, into the rectangle (x, y, w, h).
func (ex *exporter) image(i int, src *[4]float64, x, y, w, h float64) {
	if i < 0 || i >= len(ex.obj.Images) {
		ex.warnf("bad image reference")
		return
	}
	if w == 0 || h == 0 {
		return
	}
	st := &ex.st
	visible, plain := ex.look(nil)
	if !visible {
		return
	}
	// how large a unit of the image is on the page
	pt := st.ctm.areaScale()
	img := ex.imageOf(ex.obj.Images[i], func(nw, nh float64) float64 {
		if src != nil {
			nw, nh = src[2], src[3]
		}
		return pt * math.Max(math.Abs(w/nw), math.Abs(h/nh))
	})
	if img == nil || !(img.w > 0) || !(img.h > 0) {
		return
	}
	if !ex.drawn([4]float64{math.Min(x, x+w), math.Min(y, y+h), math.Max(x, x+w), math.Max(y, y+h)}, 0, plain && img.opaque) {
		return
	}
	if src == nil {
		writeImage(ex.b, img, matrix{w, 0, 0, h, x, y}, st.smooth)
		return
	}
	sx, sy, sw, sh := src[0], src[1], src[2], src[3]
	if sw == 0 || sh == 0 {
		return
	}
	// the whole image where the part of it is in its rectangle, cut to
	// the rectangle
	fw, fh := img.w*w/sw, img.h*h/sh
	fx, fy := x-sx*w/sw, y-sy*h/sh
	ex.op("q %s re W n", nums(x, y, w, h))
	writeImage(ex.b, img, matrix{fw, 0, 0, fh, fx, fy}, st.smooth)
	ex.op("Q")
}
