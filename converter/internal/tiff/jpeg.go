package tiff

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"slices"
)

// A JPEG-compressed TIFF keeps each strip or tile as a JPEG stream of its
// own, usually without the quantization and Huffman tables, which are in
// the JPEGTables tag. A strip becomes a decodable JPEG by putting the
// tables back in. The strips of a page join into one stream without
// re-encoding when they share their tables and frame and each holds whole
// MCU rows: each strip's data then becomes a restart interval (a strip
// starts with fresh DC predictions and ends byte-aligned, as a restart
// interval does), and restart markers go between them. tiff2pdf does the
// same.

// jpegStream is a JPEG stream split at its scan.
type jpegStream struct {
	pre  [][]byte // the marker segments before the frame and scan, except SOI
	sof  []byte   // the frame header segment
	sos  []byte   // the scan header segment
	scan []byte   // the entropy-coded data, without EOI
	// single reports one scan followed by EOI (or the end of the data):
	// what joining needs.
	single bool
}

// parseJPEG splits a JPEG stream (or the tables-only stream of JPEGTables,
// which has no frame or scan).
func parseJPEG(b []byte) (*jpegStream, bool) {
	if len(b) < 2 || b[0] != 0xff || b[1] != 0xd8 {
		return nil, false
	}
	s := &jpegStream{}
	for i := 2; i+1 < len(b); {
		if b[i] != 0xff {
			return nil, false
		}
		m := b[i+1]
		switch {
		case m == 0xff: // fill byte
			i++
			continue
		case m == 0xd8 || m == 0x01 || (m >= 0xd0 && m <= 0xd7):
			i += 2
			continue
		case m == 0xd9:
			return s, true
		}
		if i+4 > len(b) {
			return nil, false
		}
		n := int(binary.BigEndian.Uint16(b[i+2:]))
		if n < 2 || i+2+n > len(b) {
			return nil, false
		}
		seg := b[i : i+2+n]
		i += 2 + n
		switch {
		case m == 0xda:
			s.sos = seg
			s.scan, s.single = scanData(b[i:])
			return s, true
		case m >= 0xc0 && m <= 0xcf && m != 0xc4 && m != 0xc8 && m != 0xcc:
			s.sof = seg
		default:
			s.pre = append(s.pre, seg)
		}
	}
	return s, true
}

// scanData returns the entropy-coded data at the start of b, and whether
// it ends with EOI or the end of b (and not with another marker segment).
func scanData(b []byte) ([]byte, bool) {
	for i := 0; i+1 < len(b); i++ {
		if b[i] != 0xff {
			continue
		}
		switch m := b[i+1]; {
		case m == 0x00 || m == 0xff || (m >= 0xd0 && m <= 0xd7):
			// stuffed byte, fill, or restart marker
		case m == 0xd9:
			return b[:i], true
		default:
			return b[:i], false
		}
	}
	return b, true
}

// JPEG returns the page as one baseline JPEG stream that browsers decode,
// made from its strips without re-encoding, when the page allows it: JPEG
// compression of grey, RGB or YCbCr pixels without extra samples, in
// strips (not tiles) that share their tables and frame, each a whole
// number of MCU rows high, with no restart intervals of their own.
func (d *IFD) JPEG() ([]byte, bool) {
	l, err := d.layout()
	if err != nil {
		return nil, false
	}
	return l.joinJPEG()
}

func (l *layout) joinJPEG() ([]byte, bool) {
	if l.comp != CompressionJPEG || l.tiled || l.planar != 1 {
		return nil, false
	}
	switch {
	case l.photo == PhotometricBlackIsZero && l.spp == 1:
	case (l.photo == PhotometricRGB || l.photo == PhotometricYCbCr) && l.spp == 3:
	default:
		return nil, false
	}
	tables, _ := l.tables()
	var first *jpegStream
	var scans [][]byte
	var mcuW, mcuH int
	for i := range l.down {
		b, ok := l.rawBlock(i)
		if !ok {
			return nil, false
		}
		s, ok := parseJPEG(b)
		if !ok || s.sof == nil || s.sos == nil || !s.single || len(s.sof) < 10 {
			return nil, false
		}
		if m := s.sof[1]; m != 0xc0 && m != 0xc1 {
			return nil, false // progressive, lossless, arithmetic
		}
		for _, p := range s.pre {
			if p[1] == 0xdd {
				return nil, false // restart intervals inside the strip
			}
		}
		height := int(binary.BigEndian.Uint16(s.sof[5:]))
		if int(binary.BigEndian.Uint16(s.sof[7:])) != l.w || height < l.rows(i) || (i < l.down-1 && height != l.bh) {
			return nil, false
		}
		if first == nil {
			first = s
			nf := int(s.sof[9])
			if nf != l.spp || len(s.sof) < 10+3*nf {
				return nil, false
			}
			mcuW, mcuH = 8, 8
			if nf > 1 {
				for c := range nf {
					hv := s.sof[11+3*c]
					mcuW, mcuH = max(mcuW, 8*int(hv>>4)), max(mcuH, 8*int(hv&15))
				}
			}
			if l.down > 1 && l.bh%mcuH != 0 {
				return nil, false
			}
		} else if !sameSegments(s.pre, first.pre) || !bytes.Equal(s.sos, first.sos) || !sameFrame(s.sof, first.sof) {
			return nil, false
		}
		scans = append(scans, s.scan)
	}
	interval := (l.w + mcuW - 1) / mcuW * (l.bh / mcuH)
	if l.down > 1 && interval > 0xffff || l.h > 0xffff {
		return nil, false
	}
	for _, p := range tables {
		if p[1] == 0xdd {
			return nil, false
		}
	}

	var out bytes.Buffer
	out.Write([]byte{0xff, 0xd8})
	pre := append(append([][]byte(nil), tables...), first.pre...)
	if l.photo == PhotometricRGB {
		pre = slices.DeleteFunc(pre, isJFIF) // JFIF means YCbCr
		if m := adobeMarker(pre, first.sof, l.photo); m != nil {
			pre = append(pre, m)
		}
	}
	for _, p := range pre {
		out.Write(p)
	}
	if l.down > 1 {
		out.Write([]byte{0xff, 0xdd, 0, 4, byte(interval >> 8), byte(interval)})
	}
	sof := append([]byte(nil), first.sof...)
	binary.BigEndian.PutUint16(sof[5:], uint16(l.h))
	out.Write(sof)
	out.Write(first.sos)
	for i, s := range scans {
		if i > 0 {
			out.Write([]byte{0xff, 0xd0 + byte((i-1)%8)})
		}
		out.Write(s)
	}
	out.Write([]byte{0xff, 0xd9})
	return out.Bytes(), true
}

// tables returns the marker segments of JPEGTables.
func (l *layout) tables() ([][]byte, bool) {
	b := l.d.Bytes(TagJPEGTables)
	if b == nil {
		return nil, true
	}
	s, ok := parseJPEG(b)
	if !ok {
		return nil, false
	}
	return s.pre, true
}

// rawBlock returns the stored bytes of strip or tile i.
func (l *layout) rawBlock(i int) ([]byte, bool) {
	f := l.d.f
	off, n := l.offsets[i], l.counts[i]
	if off >= uint64(f.size) || n == 0 {
		return nil, false
	}
	n = min(n, uint64(f.size)-off)
	b := make([]byte, n)
	if _, err := f.r.ReadAt(b, int64(off)); err != nil && err != io.EOF {
		return nil, false
	}
	return b, true
}

func sameSegments(a, b [][]byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !bytes.Equal(a[i], b[i]) {
			return false
		}
	}
	return true
}

// sameFrame compares frame headers but for their heights.
func sameFrame(a, b []byte) bool {
	return len(a) == len(b) && bytes.Equal(a[:5], b[:5]) && bytes.Equal(a[7:], b[7:])
}

func isJFIF(p []byte) bool { return p[1] == 0xe0 && len(p) >= 9 && string(p[4:9]) == "JFIF\x00" }

// adobeMarker returns the Adobe marker (APP14, transform 0) that makes a
// decoder take the components of an RGB or CMYK page as they are, not as
// YCbCr or YCCK, or nil when the stream has one already or names its
// components R, G and B. libtiff writes CMYK JPEG data without the marker.
func adobeMarker(pre [][]byte, sof []byte, photo int) []byte {
	if photo != PhotometricRGB && photo != PhotometricSeparated {
		return nil
	}
	for _, p := range pre {
		if p[1] == 0xee {
			return nil
		}
	}
	if photo == PhotometricRGB && len(sof) >= 17 && sof[10] == 'R' && sof[13] == 'G' && sof[16] == 'B' {
		return nil
	}
	return []byte{0xff, 0xee, 0, 14, 'A', 'd', 'o', 'b', 'e', 0, 100, 0, 0, 0, 0, 0}
}

// decodeJPEG decodes a JPEG page: the joined stream when there is one,
// otherwise strip by strip or tile by tile.
func (l *layout) decodeJPEG() (image.Image, bool, error) {
	if l.planar != 1 {
		return nil, false, UnsupportedError("JPEG compression of planar pixels")
	}
	if l.photo == PhotometricPalette {
		return nil, false, UnsupportedError("JPEG compression of palette pixels")
	}
	if int64(l.w)*int64(l.h)*4 > MaxDecodeBytes {
		return nil, false, fmt.Errorf("tiff: a %d×%d image is too large to decode", l.w, l.h)
	}
	if b, ok := l.joinJPEG(); ok {
		if img, err := jpeg.Decode(bytes.NewReader(b)); err == nil {
			return img, false, nil
		}
	}
	tables, _ := l.tables()
	var dst draw.Image
	gray := l.ncolor == 1 && l.alpha == 0
	if gray {
		dst = image.NewGray(image.Rect(0, 0, l.w, l.h))
	} else {
		m := image.NewRGBA(image.Rect(0, 0, l.w, l.h))
		fillAlpha(m.Pix)
		dst = m
	}
	damaged := false
	for by := range l.down {
		for bx := range l.across {
			b, ok := l.rawBlock(by*l.across + bx)
			var img image.Image
			var err error
			added := false
			if ok {
				var sb []byte
				sb, added = l.standalone(tables, b)
				// The frame header of a damaged strip may claim any size.
				if c, cerr := jpeg.DecodeConfig(bytes.NewReader(sb)); cerr != nil || c.Width > l.bw+16 || c.Height > l.bh+16 {
					ok = false
				} else {
					img, err = jpeg.Decode(bytes.NewReader(sb))
				}
			}
			if !ok || err != nil {
				damaged = true
				continue
			}
			r := image.Rect(bx*l.bw, by*l.bh, bx*l.bw+l.bw, by*l.bh+l.bh).Intersect(dst.Bounds())
			if c, ok := img.(*image.CMYK); ok && added {
				img = invertCMYK(c)
			}
			draw.Draw(dst, r, img, img.Bounds().Min, draw.Src)
		}
	}
	if gray && l.photo == PhotometricWhiteIsZero {
		for i, v := range dst.(*image.Gray).Pix {
			dst.(*image.Gray).Pix[i] = 255 - v
		}
	}
	return dst, damaged, nil
}

// standalone turns a strip or tile into a JPEG stream of its own: the
// tables go back in, and an Adobe marker when the page needs one (added
// reports that).
func (l *layout) standalone(tables [][]byte, b []byte) (out []byte, added bool) {
	s, ok := parseJPEG(b)
	if !ok || s.sof == nil {
		return b, false
	}
	var buf bytes.Buffer
	buf.Write([]byte{0xff, 0xd8})
	if m := adobeMarker(append(append([][]byte(nil), tables...), s.pre...), s.sof, l.photo); m != nil {
		buf.Write(m)
		added = true
	}
	for _, p := range tables {
		buf.Write(p)
	}
	buf.Write(b[2:])
	return buf.Bytes(), added
}

// invertCMYK undoes the inversion Go's decoder applies to CMYK JPEG data
// with an Adobe marker (Adobe applications store the inks inverted), for
// data that only has the marker because adobeMarker added it.
func invertCMYK(c *image.CMYK) *image.CMYK {
	out := image.NewCMYK(c.Rect)
	for i, v := range c.Pix {
		out.Pix[i] = 255 - v
	}
	return out
}
