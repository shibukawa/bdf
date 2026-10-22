package pdfbdf

import (
	"fmt"
	"image"
	"image/draw"
	"math"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/imgconv"
	"github.com/shibukawa/bdf/raster/imagebdf"
	"github.com/shibukawa/bdf/vector/internal/export"
)

// pdfImage is an image of the document as an image XObject of the file.
type pdfImage struct {
	name string
	ref  ref
	// the size of the image in the units that IMAGE_SUB and patterns
	// count in: its pixels, or the natural size of an SVG image
	w, h float64
}

// svgDPI is the resolution SVG images are drawn at, for the size they are
// first shown at: a PDF holds pictures of them.
const svgDPI = 300

// imageOf returns the image XObject of an image part, or nil for one that
// cannot be written. scale gives how many points a unit of an image of a
// size is shown as, where it is first used (an SVG image is drawn for it).
func (ex *exporter) imageOf(h bdf.Hash, scale func(w, h float64) float64, smooth bool) *pdfImage {
	if img, ok := ex.images[h]; ok {
		return img
	}
	ex.images[h] = nil
	part := ex.Doc.Part(h)
	if part == nil {
		ex.Warnf("image %s is missing", h)
		return nil
	}
	data := part.Data
	out := &pdfImage{name: ex.newName("Im")}
	interpolate := ""
	if smooth {
		interpolate = ex.interpolate()
	}
	switch format := imgconv.Sniff(data); {
	case export.IsSVG(data):
		px := func(w, h float64) float64 {
			k := scale(w, h) * svgDPI / 72
			if !(k > 0) || math.IsInf(k, 0) {
				k = 1
			}
			return k
		}
		pix, w, hh, warnings, err := imagebdf.SVGImage(data, px, &imagebdf.Options{FontFS: ex.opts.FontFS, FontDirs: ex.opts.FontDirs, NoSystemFonts: ex.opts.NoSystemFonts})
		for _, w := range warnings {
			ex.Warnf("%s", w)
		}
		if err != nil {
			ex.Warnf("an SVG image is not drawn: %v", err)
			return nil
		}
		out.w, out.h = w, hh
		out.ref = ex.pixels(pix, interpolate)
	case format == "jpeg" && imgconv.Orientation(data) <= 1:
		w, hh, comps, ok := jpegInfo(data)
		if !ok {
			return ex.decoded(h, out, data, interpolate)
		}
		cs := "/DeviceRGB"
		if comps == 1 {
			cs = "/DeviceGray"
		}
		out.w, out.h = float64(w), float64(hh)
		// the file as it is: a PDF holds JPEG images
		out.ref = ex.pw.addStream(fmt.Sprintf(" /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace %s /BitsPerComponent 8 /Filter /DCTDecode%s", w, hh, cs, interpolate), data, false)
	default:
		return ex.decoded(h, out, data, interpolate)
	}
	ex.images[h] = out
	return out
}

// interpolate is the entry that asks readers to smooth an image that is
// shown larger than it is. A PDF/A file does not ask: what its readers
// show is to be the same.
func (ex *exporter) interpolate() string {
	if ex.pdfa != nil {
		return ""
	}
	return " /Interpolate true"
}

// decoded writes an image as its pixels.
func (ex *exporter) decoded(h bdf.Hash, out *pdfImage, data []byte, interpolate string) *pdfImage {
	img, err := imgconv.Decode(data)
	if err != nil {
		switch format := imgconv.Sniff(data); {
		case format == "avif":
			ex.Warnf("AVIF images are not written")
		case err == imgconv.ErrTooLarge:
			ex.Warnf("an image of more than %d pixels is not written", imgconv.MaxDecodePixels)
		default:
			ex.Warnf("an image (%s) cannot be decoded", format)
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
	out.w, out.h = float64(b.Dx()), float64(b.Dy())
	out.ref = ex.pixels(pix, interpolate)
	ex.images[h] = out
	return out
}

// pixels writes an image XObject of pixels: gray when all are, with a
// soft mask of their alphas when not all are opaque.
func (ex *exporter) pixels(pix *image.NRGBA, interpolate string) ref {
	w, h := pix.Rect.Dx(), pix.Rect.Dy()
	gray, opaque := true, true
	for y := 0; y < h && (gray || opaque); y++ {
		row := pix.Pix[y*pix.Stride : y*pix.Stride+4*w]
		for x := 0; x < 4*w; x += 4 {
			if row[x] != row[x+1] || row[x] != row[x+2] {
				gray = false
			}
			if row[x+3] != 0xff {
				opaque = false
			}
		}
	}
	// rows with the filter of PNG that takes the row above from each
	// (they deflate smaller)
	plane := func(comps int, at func(row []byte, x int, out []byte)) []byte {
		out := make([]byte, 0, h*(w*comps+1))
		prev := make([]byte, w*comps)
		cur := make([]byte, w*comps)
		for y := 0; y < h; y++ {
			row := pix.Pix[y*pix.Stride : y*pix.Stride+4*w]
			for x := 0; x < w; x++ {
				at(row, x, cur[x*comps:])
			}
			out = append(out, 2)
			for i, v := range cur {
				out = append(out, v-prev[i])
			}
			prev, cur = cur, prev
		}
		return out
	}
	parms := func(comps int) string {
		return fmt.Sprintf(" /DecodeParms << /Predictor 15 /Colors %d /BitsPerComponent 8 /Columns %d >>", comps, w)
	}
	smask := ""
	if !opaque {
		data := plane(1, func(row []byte, x int, out []byte) { out[0] = row[4*x+3] })
		if !ex.pw.compress {
			data = unfilter(data, w, 1)
		}
		m := ex.pw.addStream(fmt.Sprintf(" /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceGray /BitsPerComponent 8%s%s", w, h, ex.parms(parms(1)), interpolate), data, true)
		smask = " /SMask " + m.String()
	}
	cs, comps := "/DeviceRGB", 3
	var data []byte
	if gray {
		cs, comps = "/DeviceGray", 1
		data = plane(1, func(row []byte, x int, out []byte) { out[0] = row[4*x] })
	} else {
		data = plane(3, func(row []byte, x int, out []byte) { copy(out[:3], row[4*x:4*x+3]) })
	}
	if !ex.pw.compress {
		data = unfilter(data, w, comps)
	}
	return ex.pw.addStream(fmt.Sprintf(" /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace %s /BitsPerComponent 8%s%s%s", w, h, cs, ex.parms(parms(comps)), smask, interpolate), data, true)
}

// parms returns the decode parameters of a deflated image: none for a
// file that is not compressed, whose images are their pixels.
func (ex *exporter) parms(p string) string {
	if !ex.pw.compress {
		return ""
	}
	return p
}

// unfilter takes the filter of rows away again.
func unfilter(data []byte, w, comps int) []byte {
	n := w * comps
	out := make([]byte, 0, len(data))
	prev := make([]byte, n)
	for len(data) >= n+1 {
		row := data[1 : n+1]
		for i := range row {
			prev[i] += row[i]
		}
		out = append(out, prev...)
		data = data[n+1:]
	}
	return out
}

// jpegInfo reads the size and the number of components of a JPEG image
// that a PDF can hold as it is: of 8 bits, coded with Huffman tables,
// gray or of three components.
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
		case 0xc0, 0xc1, 0xc2: // baseline, extended and progressive
			if n < 8 || data[pos+4] != 8 {
				return 0, 0, 0, false
			}
			h = int(data[pos+5])<<8 | int(data[pos+6])
			w = int(data[pos+7])<<8 | int(data[pos+8])
			comps = int(data[pos+9])
			return w, h, comps, w > 0 && h > 0 && (comps == 1 || comps == 3)
		case 0xc3, 0xc5, 0xc6, 0xc7, 0xc9, 0xca, 0xcb, 0xcd, 0xce, 0xcf, 0xda, 0xd9:
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
		ex.Warnf("bad image reference")
		return
	}
	if w == 0 || h == 0 {
		return
	}
	if ex.sigNext {
		ex.sigNext = false
		ex.fieldRect(x, y, w, h)
	}
	s := ex.s
	st := &s.st
	// how large a unit of the image is on the page
	pt := s.toPage.Mul(st.ctm).AreaScale()
	img := ex.imageOf(ex.obj.Images[i], func(nw, nh float64) float64 {
		if src != nil {
			nw, nh = src[2], src[3]
		}
		return pt * math.Max(math.Abs(w/nw), math.Abs(h/nh))
	}, st.smooth)
	if img == nil || !(img.w > 0) || !(img.h > 0) {
		return
	}
	if st.shadowOn() {
		hash := ex.obj.Images[i]
		ex.castShadow([4]float64{math.Min(x, x+w), math.Min(y, y+h), math.Max(x, x+w), math.Max(y, y+h)}, 0, func(o *bdf.Object) {
			ex.shadows.part(ex, hash)
			if src == nil {
				o.Image(o.AddImage(hash), float32(x), float32(y), float32(w), float32(h))
			} else {
				o.ImageSub(o.AddImage(hash), float32(src[0]), float32(src[1]), float32(src[2]), float32(src[3]), float32(x), float32(y), float32(w), float32(h))
			}
		})
	}
	ex.blendMode()
	if st.alpha != st.pCa {
		s.op("%s gs", ex.extGState("/ca "+export.Fine(st.alpha)))
		st.pCa = st.alpha
	}
	name := s.use("XObject", img.name, img.ref)
	if src == nil {
		// the unit square of an image has its y axis up
		s.op("q %s cm %s Do Q", export.Mat(matrix{w, 0, 0, -h, x, y + h}), name)
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
	s.op("q %s re W n %s cm %s Do Q", export.Nums(x, y, w, h), export.Mat(matrix{fw, 0, 0, -fh, fx, fy + fh}), name)
}
