package tiff

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"math/bits"

	"golang.org/x/image/tiff/lzw"
)

// MaxDecodeBytes bounds the memory of a decoded page (the pixels of the
// image Decode returns).
var MaxDecodeBytes int64 = 1 << 30

// UnsupportedError reports a valid TIFF feature this package does not read.
type UnsupportedError string

func (e UnsupportedError) Error() string { return "tiff: unsupported " + string(e) }

// layout is what decoding a page needs from its tags.
type layout struct {
	d                    *IFD
	w, h                 int
	spp, bps             int
	photo, comp          int
	planar               int // 1 chunky, 2 planar
	pred                 int
	fill                 int
	ncolor               int // colour samples: 1 grey or palette, 3 RGB, 4 CMYK
	alpha                int // 0 none, 1 associated (premultiplied), 2 unassociated
	tiled                bool
	bw, bh, across, down int
	offsets, counts      []uint64
}

func (d *IFD) layout() (*layout, error) {
	l := &layout{d: d, photo: d.Photometric(), comp: d.Compression()}
	l.w, l.h = d.Size()
	if l.w <= 0 || l.h <= 0 {
		return nil, errors.New("tiff: the image has no size")
	}
	l.spp = int(d.Uint(TagSamplesPerPixel, 1))
	if l.spp < 1 || l.spp > 16 {
		return nil, fmt.Errorf("tiff: %d samples per pixel", l.spp)
	}
	bps := d.Uints(TagBitsPerSample)
	if len(bps) == 0 {
		bps = []uint64{1}
	}
	for _, b := range bps[1:] {
		if b != bps[0] {
			return nil, UnsupportedError("samples of different sizes")
		}
	}
	l.bps = int(bps[0])
	switch l.bps {
	case 1, 2, 4, 8, 16:
	default:
		return nil, UnsupportedError(fmt.Sprintf("%d bits per sample", l.bps))
	}
	if f := d.Uint(TagSampleFormat, 1); f != 1 {
		return nil, UnsupportedError("sample format " + map[uint64]string{2: "signed integer", 3: "floating point"}[f])
	}
	l.planar = int(d.Uint(TagPlanarConfig, 1))
	if l.spp == 1 {
		l.planar = 1
	}
	if l.planar != 1 && l.planar != 2 {
		return nil, fmt.Errorf("tiff: planar configuration %d", l.planar)
	}
	l.pred = int(d.Uint(TagPredictor, 1))
	switch {
	case l.pred == 1:
	case l.pred == 2 && (l.bps == 8 || l.bps == 16):
	default:
		return nil, UnsupportedError(fmt.Sprintf("predictor %d with %d bits per sample", l.pred, l.bps))
	}
	l.fill = int(d.Uint(TagFillOrder, 1))

	switch l.photo {
	case PhotometricWhiteIsZero, PhotometricBlackIsZero:
		l.ncolor = 1
	case PhotometricRGB:
		l.ncolor = 3
	case PhotometricPalette:
		l.ncolor = 1
		if l.bps > 8 {
			return nil, UnsupportedError("16-bit palette")
		}
		if len(d.Uints(TagColorMap)) < 3<<l.bps {
			return nil, errors.New("tiff: palette image without a colour map")
		}
	case PhotometricSeparated:
		if d.Uint(TagInkSet, 1) != 1 {
			return nil, UnsupportedError("inks other than CMYK")
		}
		l.ncolor = 4
	case PhotometricYCbCr:
		if l.comp != CompressionJPEG {
			return nil, UnsupportedError("YCbCr pixels outside JPEG")
		}
		l.ncolor = 3
	default:
		return nil, UnsupportedError(fmt.Sprintf("photometric interpretation %d", l.photo))
	}
	if l.spp < l.ncolor {
		return nil, fmt.Errorf("tiff: %d samples per pixel for photometric interpretation %d", l.spp, l.photo)
	}
	if l.spp > l.ncolor && l.photo != PhotometricPalette && l.photo != PhotometricSeparated {
		if es := d.Uints(TagExtraSamples); len(es) > 0 && (es[0] == 1 || es[0] == 2) {
			l.alpha = int(es[0])
		}
	}

	if d.Has(TagTileWidth) {
		l.tiled = true
		l.bw, l.bh = int(min(d.Uint(TagTileWidth, 0), 1<<30)), int(min(d.Uint(TagTileLength, 0), 1<<30))
		l.offsets, l.counts = d.Uints(TagTileOffsets), d.Uints(TagTileByteCounts)
		if l.bw <= 0 || l.bh <= 0 {
			return nil, errors.New("tiff: bad tile size")
		}
	} else {
		l.bw, l.bh = l.w, int(min(d.Uint(TagRowsPerStrip, uint64(l.h)), uint64(l.h)))
		if l.bh <= 0 {
			l.bh = l.h
		}
		l.offsets, l.counts = d.Uints(TagStripOffsets), d.Uints(TagStripByteCounts)
	}
	l.across, l.down = (l.w+l.bw-1)/l.bw, (l.h+l.bh-1)/l.bh
	// A tile may be larger than the image; a damaged file must not make
	// decoding one allocate more than the page could take.
	if int64(l.rowBytes())*int64(l.bh) > MaxDecodeBytes {
		return nil, fmt.Errorf("tiff: %d×%d tiles are too large", l.bw, l.bh)
	}
	n := l.blocks()
	if len(l.offsets) < n {
		return nil, fmt.Errorf("tiff: %d strips or tiles for %d", len(l.offsets), n)
	}
	if len(l.counts) < n {
		if n != 1 || l.comp != CompressionNone {
			return nil, errors.New("tiff: missing strip or tile byte counts")
		}
		l.counts = []uint64{uint64(l.rowBytes()) * uint64(l.bh)}
	}
	return l, nil
}

// blocks returns the number of strips or tiles.
func (l *layout) blocks() int {
	n := l.across * l.down
	if l.planar == 2 {
		n *= l.spp
	}
	return n
}

// samples returns the samples per pixel in a block (1 when planar).
func (l *layout) samples() int {
	if l.planar == 2 {
		return 1
	}
	return l.spp
}

// rowBytes returns the bytes of a row of a block.
func (l *layout) rowBytes() int { return (l.bw*l.samples()*l.bps + 7) / 8 }

// rows returns the rows of the block at block row by: tiles are always
// whole, the last strip only has the rows left.
func (l *layout) rows(by int) int {
	if l.tiled {
		return l.bh
	}
	return min(l.bh, l.h-by*l.bh)
}

// Decode decodes the page's pixels into an *image.Paletted (bilevel and
// palette images), *image.Gray, *image.RGBA, *image.NRGBA (unassociated
// alpha), *image.CMYK, or *image.YCbCr (JPEG). damaged reports pixel data
// that ended early or did not decode: what is missing is left blank.
func (d *IFD) Decode() (img image.Image, damaged bool, err error) {
	l, err := d.layout()
	if err != nil {
		return nil, false, err
	}
	switch l.comp {
	case CompressionJPEG:
		return l.decodeJPEG()
	case CompressionOldJPEG:
		return nil, false, UnsupportedError("old-style JPEG compression")
	case CompressionCCITTRLE, CompressionG3, CompressionG4:
		if l.bps != 1 || l.spp != 1 {
			return nil, false, errors.New("tiff: CCITT compression of non-bilevel pixels")
		}
	case CompressionNone, CompressionPackBits, CompressionLZW, CompressionDeflate, CompressionDeflateOld:
	default:
		return nil, false, UnsupportedError(fmt.Sprintf("compression %d", l.comp))
	}
	dst, err := l.newImage()
	if err != nil {
		return nil, false, err
	}
	planes := 1
	if l.planar == 2 {
		planes = l.spp
	}
	vals := make([]uint8, l.bw*l.samples())
	for s := range planes {
		for by := range l.down {
			for bx := range l.across {
				i := (s*l.down+by)*l.across + bx
				rows := l.rows(by)
				raw, bad := l.readBlock(i, rows)
				damaged = damaged || bad
				l.unpredict(raw, rows)
				l.put(dst, raw, vals, s, bx*l.bw, by*l.bh, rows)
			}
		}
	}
	return dst, damaged, nil
}

// newImage allocates the image the page decodes into.
func (l *layout) newImage() (image.Image, error) {
	r := image.Rect(0, 0, l.w, l.h)
	px := int64(l.w) * int64(l.h)
	bpp := int64(4)
	if l.ncolor == 1 && l.alpha == 0 {
		bpp = 1
	}
	if px*bpp > MaxDecodeBytes {
		return nil, fmt.Errorf("tiff: a %d×%d image is too large to decode", l.w, l.h)
	}
	switch {
	case l.photo == PhotometricPalette:
		cm := l.d.Uints(TagColorMap)
		n := 1 << l.bps
		pal := make(color.Palette, n)
		for i := range n {
			pal[i] = color.RGBA{uint8(cm[i] >> 8), uint8(cm[n+i] >> 8), uint8(cm[2*n+i] >> 8), 255}
		}
		return image.NewPaletted(r, pal), nil
	case l.ncolor == 1 && l.alpha == 0 && l.bps == 1:
		return image.NewPaletted(r, color.Palette{color.Gray{0}, color.Gray{255}}), nil
	case l.ncolor == 1 && l.alpha == 0:
		return image.NewGray(r), nil
	case l.ncolor == 4:
		return image.NewCMYK(r), nil
	case l.alpha == 2:
		m := image.NewNRGBA(r)
		fillAlpha(m.Pix)
		return m, nil
	}
	m := image.NewRGBA(r)
	fillAlpha(m.Pix)
	return m, nil
}

func fillAlpha(pix []byte) {
	for i := 3; i < len(pix); i += 4 {
		pix[i] = 255
	}
}

// readBlock reads and decompresses strip or tile i, padded with zeros to
// the size it should have; short reports that it was not all there.
func (l *layout) readBlock(i, rows int) (raw []byte, short bool) {
	need := l.rowBytes() * rows
	f := l.d.f
	off, n := l.offsets[i], l.counts[i]
	if off >= uint64(f.size) {
		return make([]byte, need), true
	}
	n = min(n, uint64(f.size)-off)
	src := make([]byte, n)
	if _, err := f.r.ReadAt(src, int64(off)); err != nil && err != io.EOF {
		return make([]byte, need), true
	}
	switch l.comp {
	case CompressionCCITTRLE:
		return decodeCCITT(src, ccittRLE, false, l.fill == 2, l.bw, rows)
	case CompressionG3:
		return decodeCCITT(src, ccittT4, l.d.Uint(TagT4Options, 0)&1 != 0, l.fill == 2, l.bw, rows)
	case CompressionG4:
		return decodeCCITT(src, ccittT6, false, l.fill == 2, l.bw, rows)
	}
	if l.fill == 2 {
		for k, b := range src {
			src[k] = bits.Reverse8(b)
		}
	}
	var r io.Reader
	switch l.comp {
	case CompressionNone:
		raw = src
	case CompressionPackBits:
		raw = unpackBits(src, need)
	case CompressionLZW:
		r = lzw.NewReader(bytes.NewReader(src), lzw.MSB, 8)
	case CompressionDeflate, CompressionDeflateOld:
		zr, err := zlib.NewReader(bytes.NewReader(src))
		if err != nil {
			return make([]byte, need), true
		}
		r = zr
	}
	if r != nil {
		raw = make([]byte, need)
		k, _ := io.ReadFull(r, raw)
		return raw, k < need
	}
	if len(raw) < need {
		return append(raw, make([]byte, need-len(raw))...), true
	}
	return raw[:need], false
}

// unpackBits expands PackBits data up to need bytes.
func unpackBits(src []byte, need int) []byte {
	out := make([]byte, 0, need)
	for i := 0; i < len(src) && len(out) < need; {
		n := int(int8(src[i]))
		i++
		switch {
		case n >= 0:
			end := min(i+n+1, len(src))
			out = append(out, src[i:end]...)
			i = end
		case n != -128 && i < len(src):
			for range 1 - n {
				out = append(out, src[i])
			}
			i++
		}
	}
	return out[:min(len(out), need)]
}

// unpredict undoes horizontal differencing (predictor 2) in place.
func (l *layout) unpredict(raw []byte, rows int) {
	if l.pred != 2 {
		return
	}
	ns, rb := l.samples(), l.rowBytes()
	n := l.bw * ns
	for r := range rows {
		row := raw[r*rb : (r+1)*rb]
		if l.bps == 8 {
			for i := ns; i < n; i++ {
				row[i] += row[i-ns]
			}
			continue
		}
		bo := l.d.f.bo
		for i := ns; i < n; i++ {
			bo.PutUint16(row[2*i:], bo.Uint16(row[2*i:])+bo.Uint16(row[2*(i-ns):]))
		}
	}
}

// put stores a decoded block (of plane s) at x0, y0 of dst.
func (l *layout) put(dst image.Image, raw, vals []byte, s, x0, y0, rows int) {
	ns, rb := l.samples(), l.rowBytes()
	cols := min(l.bw, l.w-x0)
	if cols <= 0 {
		return
	}
	scale := grayScale(l.bps)
	inv := l.photo == PhotometricWhiteIsZero
	for r := range min(rows, l.h-y0) {
		y := y0 + r
		v := vals[:cols*ns]
		unpack(v, raw[r*rb:(r+1)*rb], l.bps, l.d.f.bo)
		switch m := dst.(type) {
		// One-channel images take the first sample of each pixel; extra
		// samples of no known kind (and their planes) are left out.
		case *image.Paletted:
			if s > 0 {
				return
			}
			out := m.Pix[y*m.Stride+x0:][:cols]
			for x := range out {
				b := v[x*ns]
				if inv && l.photo != PhotometricPalette {
					b ^= 1
				}
				out[x] = b
			}
		case *image.Gray:
			if s > 0 {
				return
			}
			out := m.Pix[y*m.Stride+x0:][:cols]
			for x := range out {
				g := scale[v[x*ns]]
				if inv {
					g = 255 - g
				}
				out[x] = g
			}
		case *image.CMYK:
			out := m.Pix[y*m.Stride+4*x0:][:4*cols]
			l.channels(out, v, s, ns, cols, scale, false)
		case *image.RGBA:
			l.channels(m.Pix[y*m.Stride+4*x0:][:4*cols], v, s, ns, cols, scale, inv)
		case *image.NRGBA:
			l.channels(m.Pix[y*m.Stride+4*x0:][:4*cols], v, s, ns, cols, scale, inv)
		}
	}
}

// channels stores the samples of a row (v, ns per pixel; plane s when
// planar) into four-channel pixels: grey or RGB, then alpha, or CMYK.
func (l *layout) channels(out, v []byte, s, ns, cols int, scale *[256]uint8, inv bool) {
	for x := range cols {
		p := out[4*x : 4*x+4]
		for k := range ns {
			sample := k
			if l.planar == 2 {
				sample = s
			}
			b := scale[v[x*ns+k]]
			switch {
			case l.ncolor == 4:
				if sample < 4 {
					p[sample] = b
				}
			case sample < l.ncolor:
				if inv {
					b = 255 - b
				}
				if l.ncolor == 1 {
					p[0], p[1], p[2] = b, b, b
				} else {
					p[sample] = b
				}
			case sample == l.ncolor && l.alpha != 0:
				p[3] = b
			}
		}
	}
}

// unpack reads n = len(dst) samples of bps bits from row: the values of
// samples up to 8 bits, the high byte of 16-bit ones.
func unpack(dst, row []byte, bps int, bo binary.ByteOrder) {
	switch bps {
	case 8:
		copy(dst, row)
	case 16:
		hi := 0
		if bo == binary.LittleEndian {
			hi = 1
		}
		for i := range dst {
			dst[i] = row[2*i+hi]
		}
	default:
		mask := byte(1<<bps - 1)
		for i := range dst {
			bit := i * bps
			dst[i] = row[bit>>3] >> (8 - bps - bit&7) & mask
		}
	}
}

var grayScales [17]*[256]uint8

func init() {
	for _, bps := range []int{1, 2, 4, 8, 16} {
		t := new([256]uint8)
		maxv := 1<<min(bps, 8) - 1
		for i := range 256 {
			t[i] = uint8(min(i, maxv) * 255 / maxv)
		}
		grayScales[bps] = t
	}
}

// grayScale maps the unpacked samples of bps bits to 0-255.
func grayScale(bps int) *[256]uint8 { return grayScales[bps] }
