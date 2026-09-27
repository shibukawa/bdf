package psd

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// Channel data is decoded into planes of 8-bit samples, whatever the depth
// of the file: 16-bit samples are scaled, 32-bit ones (linear light) are
// encoded with the sRGB curve, 1-bit ones become black (1) and white (0).

// Compression codes.
const (
	compRaw     = 0
	compRLE     = 1 // PackBits rows
	compZip     = 2
	compZipPred = 3 // ZIP with delta prediction
)

// maxPixels bounds the images decoded (a 30000 × 30000 PSD is 900 M).
const maxPixels = 1 << 30

func rowBytes(w, depth int) int { return (w*depth + 7) / 8 }

// plane decodes one channel of w × h samples. linear marks the colour
// channels of a 32-bit image, which are linear light.
type planeSpec struct {
	w, h   int
	depth  int
	psb    bool
	linear bool
}

// decodeRows turns decompressed rows into 8-bit samples.
func (s planeSpec) decodeRows(raw []byte) []uint8 {
	out := make([]uint8, s.w*s.h)
	rb := rowBytes(s.w, s.depth)
	for y := 0; y < s.h; y++ {
		row := raw[y*rb : (y+1)*rb]
		dst := out[y*s.w : (y+1)*s.w]
		switch s.depth {
		case 1:
			for x := range dst {
				if row[x>>3]&(0x80>>(x&7)) == 0 {
					dst[x] = 255
				}
			}
		case 8:
			copy(dst, row)
		case 16:
			for x := range dst {
				dst[x] = uint8((uint32(binary.BigEndian.Uint16(row[2*x:]))*255 + 32767) / 65535)
			}
		case 32:
			for x := range dst {
				v := float64(math.Float32frombits(binary.BigEndian.Uint32(row[4*x:])))
				dst[x] = floatSample(v, s.linear)
			}
		}
	}
	return out
}

func floatSample(v float64, linear bool) uint8 {
	if !(v > 0) {
		return 0
	}
	if v >= 1 {
		return 255
	}
	if linear {
		if v <= 0.0031308 {
			v *= 12.92
		} else {
			v = 1.055*math.Pow(v, 1/2.4) - 0.055
		}
	}
	return uint8(math.Round(v * 255))
}

// packBits decodes one PackBits row into dst, which it fills exactly.
func packBits(src, dst []byte) error {
	i, o := 0, 0
	for o < len(dst) {
		if i >= len(src) {
			return fmt.Errorf("RLE row too short")
		}
		n := int(int8(src[i]))
		i++
		switch {
		case n >= 0:
			if i+n+1 > len(src) || o+n+1 > len(dst) {
				return fmt.Errorf("RLE run overflows its row")
			}
			copy(dst[o:], src[i:i+n+1])
			i += n + 1
			o += n + 1
		case n > -128:
			if i >= len(src) || o+1-n > len(dst) {
				return fmt.Errorf("RLE run overflows its row")
			}
			for k := 0; k < 1-n; k++ {
				dst[o+k] = src[i]
			}
			i++
			o += 1 - n
		}
	}
	return nil
}

// unpredict undoes the delta prediction of ZIP-with-prediction rows.
func (s planeSpec) unpredict(raw []byte) {
	rb := rowBytes(s.w, s.depth)
	for y := 0; y < s.h; y++ {
		row := raw[y*rb : (y+1)*rb]
		switch s.depth {
		case 8:
			for x := 1; x < len(row); x++ {
				row[x] += row[x-1]
			}
		case 16:
			for x := 2; x+1 < len(row); x += 2 {
				v := binary.BigEndian.Uint16(row[x:]) + binary.BigEndian.Uint16(row[x-2:])
				binary.BigEndian.PutUint16(row[x:], v)
			}
		case 32:
			// The bytes of a row are stored plane by plane (all first bytes,
			// then all second bytes...) and delta coded as one run.
			for x := 1; x < len(row); x++ {
				row[x] += row[x-1]
			}
			tmp := make([]byte, len(row))
			for x := 0; x < s.w; x++ {
				for b := 0; b < 4; b++ {
					tmp[4*x+b] = row[b*s.w+x]
				}
			}
			copy(row, tmp)
		}
	}
}

// inflate decompresses ZIP channel data into want bytes.
func inflate(data []byte, want int) ([]byte, error) {
	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	out := make([]byte, want)
	if _, err := io.ReadFull(zr, out); err != nil {
		return nil, err
	}
	return out, nil
}

// decodeChannel decodes the data of one layer channel: a compression code
// and the compressed samples.
func (s planeSpec) decodeChannel(data []byte) ([]uint8, error) {
	if s.w <= 0 || s.h <= 0 {
		return nil, nil
	}
	if int64(s.w)*int64(s.h) > maxPixels {
		return nil, fmt.Errorf("channel of %d × %d pixels", s.w, s.h)
	}
	if len(data) < 2 {
		return nil, fmt.Errorf("channel data missing")
	}
	comp := int(binary.BigEndian.Uint16(data))
	data = data[2:]
	rb := rowBytes(s.w, s.depth)
	switch comp {
	case compRaw:
		if len(data) < rb*s.h {
			return nil, fmt.Errorf("raw channel data too short")
		}
		return s.decodeRows(data), nil
	case compRLE:
		cw := 2
		if s.psb {
			cw = 4
		}
		if len(data) < cw*s.h {
			return nil, fmt.Errorf("RLE row counts missing")
		}
		counts := data[:cw*s.h]
		data = data[cw*s.h:]
		raw := make([]byte, rb*s.h)
		for y := 0; y < s.h; y++ {
			n := int(binary.BigEndian.Uint16(counts[cw*y+cw-2:]))
			if cw == 4 {
				n = int(binary.BigEndian.Uint32(counts[4*y:]))
			}
			if n > len(data) {
				return nil, fmt.Errorf("RLE data too short")
			}
			if err := packBits(data[:n], raw[y*rb:(y+1)*rb]); err != nil {
				return nil, err
			}
			data = data[n:]
		}
		return s.decodeRows(raw), nil
	case compZip, compZipPred:
		raw, err := inflate(data, rb*s.h)
		if err != nil {
			return nil, fmt.Errorf("ZIP channel data: %w", err)
		}
		if comp == compZipPred {
			s.unpredict(raw)
		}
		return s.decodeRows(raw), nil
	}
	return nil, fmt.Errorf("unknown compression %d", comp)
}

// readPlanes decodes the first n channels of the image data section: the
// composite. The section holds one compression code for all channels; RLE
// row counts for every channel come before the rows.
func (f *file) readPlanes(n int) ([][]uint8, error) {
	h := &f.hdr
	if int64(h.w)*int64(h.h) > maxPixels {
		return nil, fmt.Errorf("image of %d × %d pixels", h.w, h.h)
	}
	c := &cursor{r: f.r, off: f.imageData, end: f.size}
	comp := int(c.u16())
	if c.err != nil {
		return nil, fmt.Errorf("image data missing")
	}
	rb := rowBytes(h.w, h.depth)
	spec := func(i int) planeSpec {
		return planeSpec{w: h.w, h: h.h, depth: h.depth, psb: h.psb, linear: i < h.baseChannels() && h.mode != modeIndexed}
	}
	planes := make([][]uint8, n)
	switch comp {
	case compRaw:
		for i := range planes {
			raw := c.bytes(int64(rb) * int64(h.h))
			if c.err != nil {
				return nil, fmt.Errorf("image data too short")
			}
			planes[i] = spec(i).decodeRows(raw)
		}
	case compRLE:
		cw := 2
		if h.psb {
			cw = 4
		}
		counts := c.bytes(int64(cw) * int64(h.h) * int64(h.channels))
		if c.err != nil {
			return nil, fmt.Errorf("RLE row counts missing")
		}
		count := func(i int) int64 {
			if cw == 4 {
				return int64(binary.BigEndian.Uint32(counts[4*i:]))
			}
			return int64(binary.BigEndian.Uint16(counts[2*i:]))
		}
		for i := range planes {
			var total int64
			for y := 0; y < h.h; y++ {
				total += count(i*h.h + y)
			}
			data := c.bytes(total)
			if c.err != nil {
				return nil, fmt.Errorf("RLE data too short")
			}
			raw := make([]byte, rb*h.h)
			for y := 0; y < h.h; y++ {
				k := count(i*h.h + y)
				if err := packBits(data[:k], raw[y*rb:(y+1)*rb]); err != nil {
					return nil, err
				}
				data = data[k:]
			}
			planes[i] = spec(i).decodeRows(raw)
		}
	case compZip, compZipPred:
		raw, err := inflate(c.bytes(c.left()), rb*h.h*n)
		if err != nil {
			return nil, fmt.Errorf("ZIP image data: %w", err)
		}
		for i := range planes {
			s := spec(i)
			p := raw[i*rb*h.h : (i+1)*rb*h.h]
			if comp == compZipPred {
				s.unpredict(p)
			}
			planes[i] = s.decodeRows(p)
		}
	default:
		return nil, fmt.Errorf("unknown compression %d", comp)
	}
	return planes, nil
}

// layerPlane decodes a channel of a layer, w × h samples.
func (f *file) layerPlane(ch channel, w, h int, linear bool) ([]uint8, error) {
	if w <= 0 || h <= 0 {
		return nil, nil
	}
	c := &cursor{r: f.r, off: ch.off, end: f.size}
	data := c.bytes(ch.length)
	if c.err != nil {
		return nil, fmt.Errorf("channel data too short")
	}
	return planeSpec{w: w, h: h, depth: f.hdr.depth, psb: f.hdr.psb, linear: linear}.decodeChannel(data)
}
