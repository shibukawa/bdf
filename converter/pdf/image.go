package pdf

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"

	"github.com/shibukawa/bdf/converter/internal/jbig2"
	"github.com/shibukawa/bdf/image/imgconv"

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

// maxImageSide bounds the width and the height of an image in pixels, and
// maxImageBytes the bytes of the samples that its dictionary calls for (its
// rows × the bytes of a row): the samples are read whole before they become
// pixels. An image of more pixels than imgconv.MaxDecodePixels, the limit
// of the images that imgconv decodes, is reduced while they do (reduced).
const (
	maxImageSide  = 65535
	maxImageBytes = 1 << 30
)

// maxImageComponents bounds the colour components of an image (DeviceN has
// at most 32 colorants).
const maxImageComponents = 32

// loadImage converts an image XObject (or inline image dict+data) into PNG/JPEG bytes.
// fill is the current fill colour, used for stencil masks.
func (c *converter) loadImage(d types.Dict, raw []byte, filters []types.PDFFilter, res types.Dict, fill bdf.Color) (*decodedImage, error) {
	px, err := c.loadPixels(d, raw, filters, res, fill)
	if err != nil {
		return nil, err
	}
	return c.store(px), nil
}

// imagePixels is a decoded image: its pixels, or the bytes to store for an
// image that is kept as it is encoded (a JPEG without masks).
type imagePixels struct {
	img      *image.NRGBA
	lossless bool // the pixels come from a lossless source
	isMask   bool
	reduced  bool // the pixels are averages of those of the image (reduced)
	encoded  *decodedImage
}

// store encodes the pixels of an image according to the image options.
func (c *converter) store(px *imagePixels) *decodedImage {
	c.stored++
	if px.encoded != nil {
		return px.encoded
	}
	di := c.storePixels(px.img, px.lossless)
	di.isMask = px.isMask
	return di
}

// imageCodec returns the image codec among the filters of an image ("" for
// an image of samples).
func imageCodec(filters []types.PDFFilter) string {
	for _, f := range filters {
		switch f.Name {
		case filter.DCT, filter.JPX, filter.JBIG2, filter.CCITTFax:
			return f.Name
		}
	}
	return ""
}

// sampleColorSpace returns the colour space of the samples of an image.
func (c *converter) sampleColorSpace(d types.Dict, res types.Dict) (*colorSpace, error) {
	csObj := d["ColorSpace"]
	if csObj == nil {
		csObj = d["CS"]
	}
	cs := c.loadColorSpace(csObj, res)
	if cs.family == "Pattern" || cs.n == 0 {
		return nil, errf("image with unsupported colour space")
	}
	if cs.n > maxImageComponents {
		return nil, fmt.Errorf("image with %d colour components", cs.n)
	}
	return cs, nil
}

// loadPixels decodes an image XObject (or inline image dict+data).
func (c *converter) loadPixels(d types.Dict, raw []byte, filters []types.PDFFilter, res types.Dict, fill bdf.Color) (*imagePixels, error) {
	p := c.pdf
	w := p.intOr(d["Width"], p.intOr(d["W"], 0))
	h := p.intOr(d["Height"], p.intOr(d["H"], 0))
	if w <= 0 || h <= 0 {
		return nil, errf("image without size")
	}
	if w > maxImageSide || h > maxImageSide {
		return nil, fmt.Errorf("image of %d × %d pixels is too large", w, h)
	}
	bpc := p.intOr(d["BitsPerComponent"], p.intOr(d["BPC"], 8))
	isMask := p.boolOr(d["ImageMask"], p.boolOr(d["IM"], false))
	decodeArr := p.nums(d["Decode"])
	if decodeArr == nil {
		decodeArr = p.nums(d["D"])
	}

	// Apply the non-image filters; find out whether an image codec remains.
	// Samples are decoded up to the rows of the image: the colour space
	// says how many bytes a row has.
	sd := &types.StreamDict{Raw: raw, FilterPipeline: filters}
	var cs *colorSpace
	var data []byte
	var err error
	codec := imageCodec(filters)
	// the bytes of the samples: of one bit a pixel for a mask and for what
	// the decoders of JBIG2 and CCITT give
	size := int64((w+7)/8) * int64(h)
	switch {
	case codec == filter.DCT, codec == filter.JPX:
		// the decoders of these have their own limits
		data, codec, err = p.decodeStream(sd)
	case size > p.sampleLimit():
		err = fmt.Errorf("image of %d × %d pixels is too large (%d bytes of samples)", w, h, size)
	case codec != "":
		data, codec, err = p.decodeStream(sd)
	case isMask:
		data, _, err = p.decodeStreamUpTo(sd, size, true)
	default:
		if cs, err = c.sampleColorSpace(d, res); err != nil {
			return nil, err
		}
		if bpc < 1 || bpc > 16 {
			return nil, fmt.Errorf("image with %d bits per component", bpc)
		}
		if size = (int64(w)*int64(cs.n)*int64(bpc) + 7) / 8 * int64(h); size > p.sampleLimit() {
			return nil, fmt.Errorf("image of %d × %d pixels is too large (%d bytes of samples)", w, h, size)
		}
		data, _, err = p.decodeStreamUpTo(sd, size, true)
	}
	if err != nil {
		return nil, err
	}
	switch codec {
	case filter.DCT:
		if !isMask && d["SMask"] == nil && d["Mask"] == nil && len(decodeArr) == 0 {
			// Pass JPEG bytes through untouched (or let imgconv try a smaller encoding).
			if cfg, err := jpeg.DecodeConfig(bytes.NewReader(data)); err == nil {
				return &imagePixels{encoded: c.storeEncoded(data, cfg.Width, cfg.Height)}, nil
			}
		}
		img, err := imgconv.Decode(data)
		if err != nil {
			return nil, err
		}
		nrgba := toNRGBA(img)
		c.applyMasks(nrgba, d, res)
		return &imagePixels{img: nrgba}, nil
	case filter.JPX:
		if isMask {
			return nil, errf("JPXDecode image masks are not allowed")
		}
		return c.loadJPX(d, data, res)
	case filter.JBIG2:
		var globals []byte
		for _, f := range filters {
			if f.Name == filter.JBIG2 {
				if gs := p.stream(f.DecodeParms["JBIG2Globals"]); gs != nil {
					if globals, _, err = p.decodeStream(gs); err != nil {
						return nil, err
					}
				}
			}
		}
		bm, err := jbig2.Decode(data, globals)
		if err != nil {
			return nil, err
		}
		// The filter's samples are 1 = white (JBIG2's 1 is black), one row
		// of the image dictionary's width after another.
		rowBytes := (w + 7) / 8
		data = make([]byte, rowBytes*h)
		for y := 0; y < h; y++ {
			row := data[y*rowBytes : (y+1)*rowBytes]
			if y < bm.Height {
				copy(row, bm.Data[y*bm.Stride:y*bm.Stride+min(bm.Stride, rowBytes)])
			}
			for i := range row {
				row[i] = ^row[i]
			}
		}
		bpc = 1
		if !isMask && d["ColorSpace"] == nil {
			d = d.Clone().(types.Dict)
			d["ColorSpace"] = types.Name("DeviceGray")
		}
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
		// The decoder makes every row it is told of, which one bit of data
		// may stand for: the image has no use for more rows than its own.
		if parms["Columns"] > maxImageSide {
			return nil, fmt.Errorf("CCITTFaxDecode with %d columns", parms["Columns"])
		}
		parms["Rows"] = min(parms["Rows"], h)
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

	var rows imageRows
	if isMask {
		rows = stencilRows(data, w, h, decodeArr, fill)
	} else {
		if cs == nil {
			if cs, err = c.sampleColorSpace(d, res); err != nil {
				return nil, err
			}
		}
		rows = sampleRows(data, w, h, bpc, cs, decodeArr)
	}
	k := reduction(w, h, p.pixelLimit())
	if k > 1 {
		out := reduced(rows, k)
		c.warnOnce("image-reduced", "an image of %d × %d pixels is stored at %d × %d", w, h, out.Bounds().Dx(), out.Bounds().Dy())
		if !isMask {
			c.applyMasks(out, d, res)
		}
		return &imagePixels{img: out, lossless: true, isMask: isMask, reduced: true}, nil
	}
	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h && rows.row(y, out.Pix[y*out.Stride:y*out.Stride+4*w]); y++ {
	}
	if isMask {
		return &imagePixels{img: out, lossless: true, isMask: true}, nil
	}
	c.applyMasks(out, d, res)
	return &imagePixels{img: out, lossless: true}, nil
}

// imageRows is an image of samples, or a stencil mask, read a row of
// pixels at a time.
type imageRows struct {
	w, h int
	// inData is the number of rows that the data has, whole or in part:
	// the rows after them are all the same.
	inData int
	// row writes the w pixels of row y to dst (4 bytes each, as image.NRGBA
	// has them). It reports whether the image goes on: after the last row
	// the data has, the rows are transparent.
	row func(y int, dst []uint8) bool
	// fast gives the pixels that row gives, from a table of the colours
	// of the sample values (nil without one).
	fast func(y int, dst []uint8) bool
}

// stencilRows returns the rows of a stencil mask: the pixels it paints have
// the fill colour, the others are transparent.
func stencilRows(data []byte, w, h int, decodeArr []float64, fill bdf.Color) imageRows {
	// Stencil mask: sample 0 paints (unless Decode [1 0]).
	paintOnZero := true
	if len(decodeArr) >= 1 && decodeArr[0] == 1 {
		paintOnZero = false
	}
	rowBytes := (w + 7) / 8
	r, g, b := uint8(fill>>24), uint8(fill>>16), uint8(fill>>8)
	return imageRows{w: w, h: h, inData: (len(data) + rowBytes - 1) / rowBytes, row: func(y int, dst []uint8) bool {
		for x := 0; x < w; x++ {
			i := y*rowBytes + x/8
			bit := uint8(0)
			if i < len(data) {
				bit = data[i] >> (7 - uint(x%8)) & 1
			}
			px := dst[4*x : 4*x+4 : 4*x+4]
			if (bit == 0) == paintOnZero {
				px[0], px[1], px[2], px[3] = r, g, b, 255
			} else {
				px[0], px[1], px[2], px[3] = 0, 0, 0, 0
			}
		}
		return true
	}}
}

// sampleRows returns the rows of an image of samples of bpc bits in the
// colour space cs.
func sampleRows(data []byte, w, h, bpc int, cs *colorSpace, decodeArr []float64) imageRows {
	n := cs.n
	maxv := float64(int(1)<<uint(bpc) - 1)
	rowBits := w * n * bpc
	rowBytes := max((rowBits+7)/8, 1)
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
	// colour returns the colour of a pixel from the values of its samples
	// (in comps).
	colour := func() (uint8, uint8, uint8) {
		for k, raw := range comps {
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
		return to8(r), to8(g), to8(b)
	}
	rowOf := func(y int, dst []uint8) []byte {
		start := y * rowBytes
		if start >= len(data) {
			clear(dst)
			return nil
		}
		return data[start:min(len(data), start+rowBytes)]
	}
	rows := imageRows{w: w, h: h, inData: (len(data) + rowBytes - 1) / rowBytes, row: func(y int, dst []uint8) bool {
		row := rowOf(y, dst)
		if row == nil {
			return false
		}
		for x := 0; x < w; x++ {
			for k := 0; k < n; k++ {
				comps[k] = sample(row, x*n+k)
			}
			px := dst[4*x : 4*x+4 : 4*x+4]
			px[0], px[1], px[2] = colour()
			px[3] = 255
		}
		return true
	}}
	if n == 1 && bpc >= 1 && bpc <= 8 {
		// One sample a pixel of 256 values at most: the colours of the
		// values are found once.
		table := make([][3]uint8, 1<<uint(bpc))
		for v := range table {
			comps[0] = float64(v)
			table[v][0], table[v][1], table[v][2] = colour()
		}
		mask, per := uint8(int(1)<<uint(bpc)-1), 8/bpc // per: the samples of a byte, when bpc divides 8
		rows.fast = func(y int, dst []uint8) bool {
			row := rowOf(y, dst)
			if row == nil {
				return false
			}
			for x := 0; x < w; x++ {
				var v uint8
				switch bpc {
				case 8:
					if x < len(row) {
						v = row[x]
					}
				case 1, 2, 4:
					if i := x / per; i < len(row) {
						v = row[i] >> uint(8-bpc*(x%per+1)) & mask
					}
				default:
					v = uint8(sample(row, x))
				}
				px := dst[4*x : 4*x+4 : 4*x+4]
				px[0], px[1], px[2], px[3] = table[v][0], table[v][1], table[v][2], 255
			}
			return true
		}
	}
	return rows
}

// reduction returns the factor by which an image of w × h pixels is reduced
// to have limit pixels at most: 1 for an image that has, else the smallest
// whole number k for which the image of ⌈w/k⌉ × ⌈h/k⌉ pixels has.
func reduction(w, h int, limit int64) int {
	k := 1
	for int64((w+k-1)/k)*int64((h+k-1)/k) > max(limit, 1) {
		k++
	}
	return k
}

// reduced makes the pixels of an image reduced by the factor k, from its
// rows, k of them at a time: the pixels of the image itself are never held.
// A pixel is the average of the block of k × k pixels it stands for (of
// fewer at the right and the bottom of an image whose size k does not
// divide), the colours weighted by the opacity: a stencil mask becomes the
// fill colour at the opacity of what the mask paints of the block. The
// image is drawn where it was, since images are drawn in the unit square
// whatever pixels they have.
func reduced(rows imageRows, k int) *image.NRGBA {
	w, h := rows.w, rows.h
	ow, oh := (w+k-1)/k, (h+k-1)/k
	out := image.NewNRGBA(image.Rect(0, 0, ow, oh))
	row := rows.row
	if rows.fast != nil {
		row = rows.fast
	}
	line := make([]uint8, 4*w)
	sums := make([]uint32, 4*ow) // of a block: red, green and blue × opacity, and the opacity
	after := -1                  // the first row of the result made of k rows after those of the data
	for oy := 0; oy < oh; oy++ {
		y0, y1 := oy*k, min((oy+1)*k, h)
		if after >= 0 && y1-y0 == k {
			// as that row: the size of an image does not say how much work it is
			copy(out.Pix[oy*out.Stride:oy*out.Stride+4*ow], out.Pix[after*out.Stride:])
			continue
		}
		if y0 >= rows.inData && y1-y0 == k {
			after = oy
		}
		clear(sums)
		for y := y0; y < y1; y++ {
			row(y, line) // a row the data does not have is transparent
			for ox, x := 0, 0; ox < ow; ox++ {
				s := sums[4*ox : 4*ox+4 : 4*ox+4]
				for x1 := min(x+k, w); x < x1; x++ {
					px := line[4*x : 4*x+4 : 4*x+4]
					a := uint32(px[3])
					s[0] += uint32(px[0]) * a
					s[1] += uint32(px[1]) * a
					s[2] += uint32(px[2]) * a
					s[3] += a
				}
			}
		}
		o := out.Pix[oy*out.Stride : oy*out.Stride+4*ow]
		for ox := 0; ox < ow; ox++ {
			s := sums[4*ox : 4*ox+4 : 4*ox+4]
			n := uint32((y1 - y0) * (min((ox+1)*k, w) - ox*k))
			a := (s[3] + n/2) / n
			if a == 0 {
				continue
			}
			px := o[4*ox : 4*ox+4 : 4*ox+4]
			px[0] = uint8((s[0] + s[3]/2) / s[3])
			px[1] = uint8((s[1] + s[3]/2) / s[3])
			px[2] = uint8((s[2] + s[3]/2) / s[3])
			px[3] = uint8(a)
		}
	}
	return out
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

// maskPixels decodes the image of a mask (nil when it cannot be read back).
// The pixels are taken as they are decoded when storing them would keep
// them as they are, which saves encoding them only to decode them again.
// Otherwise the image is stored and read back as before: the values of a
// JPEG, or of pixels that imgconv stores lossy, are those of the stored
// bytes. The masks that the image of a mask has itself are not read: it
// must not have any, and may name itself as its mask.
func (c *converter) maskPixels(sd *types.StreamDict, res types.Dict, fill bdf.Color) (*imageMask, error) {
	c.inMask++
	px, err := c.loadPixels(sd.Dict, sd.Raw, sd.FilterPipeline, res, fill)
	c.inMask--
	if err != nil {
		return nil, err
	}
	if px.encoded == nil && (px.reduced || px.lossless && c.storedExactly(px.img)) {
		b := px.img.Bounds()
		return &imageMask{img: px.img, w: b.Dx(), h: b.Dy(), reduced: px.reduced}, nil
	}
	di := c.store(px)
	img, err := decodeStored(di)
	if err != nil {
		return nil, nil
	}
	return &imageMask{img: img, w: di.w, h: di.h}, nil
}

// storedExactly reports whether storePixels keeps the pixels of a lossless
// source as they are: always but in Convert mode, where imgconv may choose
// a lossy encoding for an image of many colours.
func (c *converter) storedExactly(img *image.NRGBA) bool {
	o := c.opts.Images
	if o.Mode != imgconv.Convert || !imgconv.Available() || o.NoLossyForLossless {
		return true
	}
	// imgconv counts more colours than this as a photograph unless the
	// options say fewer; a mask has 256 greys at most.
	colors := 256
	if o.PhotoColors > 0 {
		colors = min(colors, o.PhotoColors)
	}
	return !imgconv.PhotoLike(img, colors)
}

// imageMask is the mask of an image: the image of its /SMask or /Mask.
type imageMask struct {
	img     image.Image
	w, h    int
	stencil bool // /Mask: what the mask paints is left out (else its values are the opacity)
	reduced bool // the pixels are averages: those of a stencil say how much of them it paints
}

// at returns the red and the opacity of a pixel of the mask, as
// image.Image.At(x, y).RGBA() does.
func (m *imageMask) at(x, y int) (r, a uint32) {
	if img, ok := m.img.(*image.NRGBA); ok && x >= 0 && y >= 0 && x < m.w && y < m.h {
		px := img.Pix[y*img.Stride+4*x:]
		r, _, _, a = color.NRGBA{R: px[0], G: px[1], B: px[2], A: px[3]}.RGBA()
		return r, a
	}
	r, _, _, a = m.img.At(x, y).RGBA()
	return r, a
}

// loadMask reads the mask of an image: its /SMask (soft mask) or /Mask
// (stencil; a colour key is not read). nil for an image without one.
func (c *converter) loadMask(d types.Dict, res types.Dict) *imageMask {
	p := c.pdf
	if c.inMask > 0 {
		return nil
	}
	if sm := p.stream(d["SMask"]); sm != nil {
		m, err := c.maskPixels(sm, res, 0)
		if err != nil {
			c.warnf("soft mask: %v", err)
		}
		return m
	}
	switch m := p.deref(d["Mask"]).(type) {
	case types.StreamDict, *types.StreamDict:
		sd := p.stream(d["Mask"])
		if m, err := c.maskPixels(sd, res, bdf.RGB(0, 0, 0)); err == nil && m != nil {
			m.stencil = true
			return m
		}
	case types.Array:
		// Colour key masking is applied on decoded samples; approximate by
		// comparing the final colour against the key converted through the space.
		_ = m
	}
	return nil
}

// applyMasks applies /SMask (soft mask) or /Mask (stencil or colour key) to img.
func (c *converter) applyMasks(img *image.NRGBA, d types.Dict, res types.Dict) {
	m := c.loadMask(d, res)
	if m == nil {
		return
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			mx, my := x*m.w/w, y*m.h/h
			r, a := m.at(mx, my)
			i := img.PixOffset(x, y)
			switch {
			case !m.stencil:
				img.Pix[i+3] = uint8(r >> 8)
			case m.reduced:
				// the part of the pixel that the stencil leaves
				img.Pix[i+3] = uint8(uint32(img.Pix[i+3]) * (0xffff - a) / 0xffff)
			case a > 0:
				// Stencil painted (a>0) means masked OUT.
				img.Pix[i+3] = 0
			}
		}
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
