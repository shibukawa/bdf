// Package tiff reads TIFF files: the image file directories (IFDs) of a
// multi-page file, their tags and their pixels (docs/design.md §3.10).
//
// It covers what scanners, fax software and image editors write: classic
// TIFF and BigTIFF; strips and tiles, chunky or planar; no compression,
// PackBits, LZW, Deflate, JPEG, and CCITT fax coding (modified Huffman,
// Group 3 one- and two-dimensional, with or without fill bits, and Group
// 4; see ccitt.go); bilevel, grey, palette, RGB and CMYK pixels of 1 to 16
// bits, with or without alpha, and YCbCr in JPEG. A JPEG page whose strips
// allow it comes out as one baseline JPEG stream, without decoding
// (IFD.JPEG).
//
// CCITT data is read as libtiff reads it: the codes of black runs give 1
// bits, which the photometric interpretation then maps (WhiteIsZero, the
// usual value, makes them black). golang.org/x/image/tiff maps them the
// other way round and shows a BlackIsZero fax inverted.
package tiff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"unicode/utf8"
)

// Tags.
const (
	TagNewSubfileType   = 254
	TagImageWidth       = 256
	TagImageLength      = 257
	TagBitsPerSample    = 258
	TagCompression      = 259
	TagPhotometric      = 262
	TagFillOrder        = 266
	TagDocumentName     = 269
	TagImageDescription = 270
	TagStripOffsets     = 273
	TagOrientation      = 274
	TagSamplesPerPixel  = 277
	TagRowsPerStrip     = 278
	TagStripByteCounts  = 279
	TagXResolution      = 282
	TagYResolution      = 283
	TagPlanarConfig     = 284
	TagPageName         = 285
	TagT4Options        = 292
	TagT6Options        = 293
	TagResolutionUnit   = 296
	TagPageNumber       = 297
	TagSoftware         = 305
	TagDateTime         = 306
	TagArtist           = 315
	TagPredictor        = 317
	TagColorMap         = 320
	TagTileWidth        = 322
	TagTileLength       = 323
	TagTileOffsets      = 324
	TagTileByteCounts   = 325
	TagInkSet           = 332
	TagExtraSamples     = 338
	TagSampleFormat     = 339
	TagJPEGTables       = 347
	TagCopyright        = 33432
	TagDNGVersion       = 50706
)

// Compression schemes.
const (
	CompressionNone       = 1
	CompressionCCITTRLE   = 2
	CompressionG3         = 3
	CompressionG4         = 4
	CompressionLZW        = 5
	CompressionOldJPEG    = 6
	CompressionJPEG       = 7
	CompressionDeflate    = 8
	CompressionPackBits   = 32773
	CompressionDeflateOld = 32946
)

// Photometric interpretations.
const (
	PhotometricWhiteIsZero = 0
	PhotometricBlackIsZero = 1
	PhotometricRGB         = 2
	PhotometricPalette     = 3
	PhotometricMask        = 4
	PhotometricSeparated   = 5
	PhotometricYCbCr       = 6
	PhotometricCIELab      = 8
)

// Resolution units.
const (
	ResUnitNone       = 1
	ResUnitInch       = 2
	ResUnitCentimeter = 3
)

// Limits on what a file may make the reader do.
const (
	maxIFDs    = 1 << 16
	maxEntries = 1 << 16
)

// File is an open TIFF file.
type File struct {
	r    io.ReaderAt
	size int64
	bo   binary.ByteOrder
	big  bool
	// IFDs are the directories of the main chain, in file order.
	IFDs []*IFD
	// ChainErr, when not nil, is why the chain ends early: IFDs holds the
	// directories read before the damaged one.
	ChainErr error
}

// Sniff reports whether head starts a TIFF file, classic or BigTIFF.
func Sniff(head []byte) bool {
	if len(head) < 4 {
		return false
	}
	switch string(head[:4]) {
	case "II*\x00", "MM\x00*", "II+\x00", "MM\x00+":
		return true
	}
	return false
}

// Open reads the header and the chain of IFDs. It fails only when the
// first IFD cannot be read.
func Open(r io.ReaderAt, size int64) (*File, error) {
	var h [16]byte
	n, _ := r.ReadAt(h[:], 0)
	if n < 8 || !Sniff(h[:n]) {
		return nil, errors.New("tiff: not a TIFF file")
	}
	f := &File{r: r, size: size, bo: binary.LittleEndian}
	if h[0] == 'M' {
		f.bo = binary.BigEndian
	}
	var next int64
	if f.bo.Uint16(h[2:]) == 43 {
		if n < 16 || f.bo.Uint16(h[4:]) != 8 {
			return nil, errors.New("tiff: bad BigTIFF header")
		}
		f.big = true
		next = int64(f.bo.Uint64(h[8:]))
	} else {
		next = int64(f.bo.Uint32(h[4:]))
	}
	seen := map[int64]bool{}
	for next != 0 {
		if len(f.IFDs) == maxIFDs {
			f.ChainErr = fmt.Errorf("tiff: more than %d IFDs", maxIFDs)
			break
		}
		if seen[next] {
			f.ChainErr = errors.New("tiff: the IFD chain loops")
			break
		}
		seen[next] = true
		d, n, err := f.readIFD(next)
		if err != nil {
			if len(f.IFDs) == 0 {
				return nil, err
			}
			f.ChainErr = err
			break
		}
		f.IFDs = append(f.IFDs, d)
		next = n
	}
	if len(f.IFDs) == 0 {
		return nil, errors.New("tiff: no IFD")
	}
	return f, nil
}

// Pages returns the IFDs that are pages: those not marked (NewSubfileType)
// as a reduced-resolution version or a transparency mask of another image.
func (f *File) Pages() []*IFD {
	var out []*IFD
	for _, d := range f.IFDs {
		if d.Uint(TagNewSubfileType, 0)&5 == 0 {
			out = append(out, d)
		}
	}
	return out
}

// entry is a tag of an IFD: its type, value count and where the value is.
type entry struct {
	typ    uint16
	count  uint64
	off    int64  // file offset of the value, when it is not inline
	inline []byte // the value, when it fits in the entry
}

// IFD is an image file directory.
type IFD struct {
	f *File
	// Offset is where the IFD starts in the file.
	Offset  int64
	entries map[uint16]entry
}

var typeSize = map[uint16]uint64{1: 1, 2: 1, 3: 2, 4: 4, 5: 8, 6: 1, 7: 1, 8: 2, 9: 4, 10: 8, 11: 4, 12: 8, 13: 4, 16: 8, 17: 8, 18: 8}

func (f *File) readIFD(off int64) (*IFD, int64, error) {
	cw, ew, vw := int64(2), int64(12), 4 // count width, entry width, inline value width
	if f.big {
		cw, ew, vw = 8, 20, 8
	}
	if off < 8 || off+cw > f.size {
		return nil, 0, fmt.Errorf("tiff: IFD offset %d outside the file", off)
	}
	b := make([]byte, cw)
	if _, err := f.r.ReadAt(b, off); err != nil {
		return nil, 0, fmt.Errorf("tiff: IFD at %d: %w", off, err)
	}
	var n int64
	if f.big {
		n = int64(f.bo.Uint64(b))
	} else {
		n = int64(f.bo.Uint16(b))
	}
	if n <= 0 || n > maxEntries || off+cw+n*ew+int64(vw) > f.size {
		return nil, 0, fmt.Errorf("tiff: IFD at %d: bad entry count %d", off, n)
	}
	b = make([]byte, n*ew+int64(vw))
	if _, err := f.r.ReadAt(b, off+cw); err != nil {
		return nil, 0, fmt.Errorf("tiff: IFD at %d: %w", off, err)
	}
	d := &IFD{f: f, Offset: off, entries: make(map[uint16]entry, n)}
	for i := range n {
		e := b[i*ew:]
		tag, typ := f.bo.Uint16(e), f.bo.Uint16(e[2:])
		var count uint64
		var val []byte
		if f.big {
			count, val = f.bo.Uint64(e[4:]), e[12:20]
		} else {
			count, val = uint64(f.bo.Uint32(e[4:])), e[8:12]
		}
		size, ok := typeSize[typ]
		if !ok || count == 0 || count > uint64(f.size) {
			continue // unknown types are skipped, as the specification says
		}
		en := entry{typ: typ, count: count}
		if total := size * count; total <= uint64(vw) {
			en.inline = append([]byte(nil), val[:total]...)
		} else {
			var o uint64
			if f.big {
				o = f.bo.Uint64(val)
			} else {
				o = uint64(f.bo.Uint32(val))
			}
			if o > uint64(f.size) || total > uint64(f.size)-o {
				continue // the value lies outside the file
			}
			en.off = int64(o)
		}
		d.entries[tag] = en
	}
	var next int64
	if f.big {
		next = int64(f.bo.Uint64(b[n*ew:]))
	} else {
		next = int64(f.bo.Uint32(b[n*ew:]))
	}
	return d, next, nil
}

// Has reports whether the IFD has the tag.
func (d *IFD) Has(tag uint16) bool {
	_, ok := d.entries[tag]
	return ok
}

// raw returns the bytes of a tag's value.
func (d *IFD) raw(tag uint16) (entry, []byte, bool) {
	e, ok := d.entries[tag]
	if !ok {
		return e, nil, false
	}
	if e.inline != nil {
		return e, e.inline, true
	}
	b := make([]byte, typeSize[e.typ]*e.count)
	if _, err := d.f.r.ReadAt(b, e.off); err != nil && err != io.EOF {
		return e, nil, false
	}
	return e, b, true
}

// Uints returns the values of an integer tag (nil when it is missing or
// not an integer).
func (d *IFD) Uints(tag uint16) []uint64 {
	e, b, ok := d.raw(tag)
	if !ok {
		return nil
	}
	bo := d.f.bo
	out := make([]uint64, e.count)
	for i := range out {
		switch e.typ {
		case 1, 7:
			out[i] = uint64(b[i])
		case 6:
			out[i] = uint64(int64(int8(b[i])))
		case 3:
			out[i] = uint64(bo.Uint16(b[2*i:]))
		case 8:
			out[i] = uint64(int64(int16(bo.Uint16(b[2*i:]))))
		case 4, 13:
			out[i] = uint64(bo.Uint32(b[4*i:]))
		case 9:
			out[i] = uint64(int64(int32(bo.Uint32(b[4*i:]))))
		case 16, 17, 18:
			out[i] = bo.Uint64(b[8*i:])
		default:
			return nil
		}
	}
	return out
}

// Uint returns the first value of an integer tag, or def.
func (d *IFD) Uint(tag uint16, def uint64) uint64 {
	if v := d.Uints(tag); len(v) > 0 {
		return v[0]
	}
	return def
}

// Float returns the first value of a numeric tag (a rational for most).
func (d *IFD) Float(tag uint16) (float64, bool) {
	e, b, ok := d.raw(tag)
	if !ok {
		return 0, false
	}
	bo := d.f.bo
	switch e.typ {
	case 5:
		num, den := bo.Uint32(b), bo.Uint32(b[4:])
		if den == 0 {
			return 0, false
		}
		return float64(num) / float64(den), true
	case 10:
		num, den := int32(bo.Uint32(b)), int32(bo.Uint32(b[4:]))
		if den == 0 {
			return 0, false
		}
		return float64(num) / float64(den), true
	case 11:
		return float64(math.Float32frombits(bo.Uint32(b))), true
	case 12:
		return math.Float64frombits(bo.Uint64(b)), true
	}
	if v := d.Uints(tag); len(v) > 0 {
		return float64(v[0]), true
	}
	return 0, false
}

// Bytes returns the value of a BYTE or UNDEFINED tag.
func (d *IFD) Bytes(tag uint16) []byte {
	e, b, ok := d.raw(tag)
	if !ok || (e.typ != 1 && e.typ != 7) {
		return nil
	}
	return b
}

// Strings returns the strings of an ASCII tag: the NUL-separated values,
// without empty ones and surrounding space. Text that is not UTF-8 is read
// as Latin-1.
func (d *IFD) Strings(tag uint16) []string {
	e, b, ok := d.raw(tag)
	if !ok || (e.typ != 2 && e.typ != 7 && e.typ != 1) {
		return nil
	}
	var out []string
	for _, s := range strings.Split(string(b), "\x00") {
		if !utf8.ValidString(s) {
			r := make([]rune, len(s))
			for i := range len(s) {
				r[i] = rune(s[i])
			}
			s = string(r)
		}
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// Size returns the width and height of the image in pixels.
func (d *IFD) Size() (w, h int) {
	return int(min(d.Uint(TagImageWidth, 0), math.MaxInt32)), int(min(d.Uint(TagImageLength, 0), math.MaxInt32))
}

// Resolution returns the pixels per unit along x and y (0 when a tag is
// missing or not positive) and the unit (ResUnitInch by default).
func (d *IFD) Resolution() (x, y float64, unit int) {
	x, _ = d.Float(TagXResolution)
	y, _ = d.Float(TagYResolution)
	if !(x > 0) || math.IsInf(x, 0) {
		x = 0
	}
	if !(y > 0) || math.IsInf(y, 0) {
		y = 0
	}
	return x, y, int(d.Uint(TagResolutionUnit, ResUnitInch))
}

// Orientation returns the Orientation tag (1 to 8; 1, rows top to bottom
// and columns left to right, when it is missing or out of range).
func (d *IFD) Orientation() int {
	if o := d.Uint(TagOrientation, 1); o >= 1 && o <= 8 {
		return int(o)
	}
	return 1
}

// Compression returns the compression scheme (CompressionNone by default).
func (d *IFD) Compression() int { return int(d.Uint(TagCompression, CompressionNone)) }

// Photometric returns the photometric interpretation. When the tag is
// missing, as it sometimes is in fax files, it is guessed: WhiteIsZero
// for CCITT data, RGB for three samples, BlackIsZero otherwise.
func (d *IFD) Photometric() int {
	if v := d.Uints(TagPhotometric); len(v) > 0 {
		return int(v[0])
	}
	switch c := d.Compression(); {
	case c == CompressionCCITTRLE || c == CompressionG3 || c == CompressionG4:
		return PhotometricWhiteIsZero
	case d.Uint(TagSamplesPerPixel, 1) >= 3:
		return PhotometricRGB
	}
	return PhotometricBlackIsZero
}
