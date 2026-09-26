// Package jbig2 decodes JBIG2 bi-level images (ITU-T T.88) in the form PDF
// embeds them with the JBIG2Decode filter: the sequential organisation without
// the file header, optionally preceded by a stream of global segments.
//
// It decodes generic regions (arithmetic templates 0–3 with typical prediction,
// and MMR), generic refinement regions, text regions with their symbol
// dictionaries (arithmetic or Huffman coded, with refinement and aggregation),
// halftone regions with their pattern dictionaries, custom Huffman tables and
// striped pages of initially unknown height. Extended templates and colour
// extensions (T.88 Amendments 2 and 3) are not supported.
//
// The input is untrusted: sizes are checked against fixed limits, and malformed
// data yields an error rather than a panic.
package jbig2

import (
	"errors"
	"fmt"
)

// Bitmap is a decoded page. Pixels are packed 8 to a byte, most significant bit
// first, and every row starts on a byte boundary; 1 is black and 0 is white,
// the JBIG2 convention. Padding bits at the end of a row are 0.
//
// PDF's JBIG2Decode filter delivers the inverse: in the decoded sample stream a
// 1 bit is white and a 0 bit is black, as in a 1-bit DeviceGray image with the
// default Decode array. A PDF reader therefore inverts Data.
type Bitmap struct {
	Width, Height int
	Stride        int    // bytes per row
	Data          []byte // 1 bit per pixel, MSB first, 1 = black (the JBIG2 convention)
}

// Limits on sizes taken from the data.
const (
	maxDimension = 1 << 24 // width or height of a bitmap
	maxPixels    = 1 << 28 // pixels of a bitmap (32 MiB)
	maxWork      = 1 << 30 // units of work (about pixels allocated, decoded or drawn) in one Decode
	maxSymbols   = 1 << 20 // symbols in a dictionary or available to a text region
)

var errTruncated = errors.New("jbig2: unexpected end of data")

func errorf(format string, args ...any) error { return fmt.Errorf("jbig2: "+format, args...) }

// Decode decodes a JBIG2 stream as embedded in PDF (JBIG2Decode): the sequential
// organisation without the file header, page 1, with the optional JBIG2Globals
// stream whose segments are processed first. The result uses the JBIG2
// convention, 1 = black; see Bitmap.
//
// A stream cut short still decodes: a region whose arithmetic coded data ends
// early keeps what was decoded before, and a segment cut off is decoded from
// the data there is.
func Decode(data, globals []byte) (bm *Bitmap, err error) {
	defer func() {
		if r := recover(); r != nil {
			bm, err = nil, errorf("malformed data (%v)", r)
		}
	}()
	d := &decoder{segments: map[uint32]*segment{}}
	if err := d.run(globals, true); err != nil {
		return nil, err
	}
	d.ended = false
	if err := d.run(data, false); err != nil {
		return nil, err
	}
	if d.page == nil {
		return nil, errorf("no page information segment")
	}
	bm = d.page.bm
	if bm.Width == 0 || bm.Height == 0 {
		return nil, errorf("empty page")
	}
	return bm, nil
}

// decoder holds the state of one Decode call.
type decoder struct {
	segments map[uint32]*segment // processed segments, for referral
	page     *page
	pageNum  uint32 // page association of the page being decoded
	ended    bool   // end of page or end of file seen
	work     int64
}

// segment is a processed segment and its result: a *symbolDict,
// *patternDict, *huffTable or, for an intermediate region, a *Bitmap.
type segment struct {
	segmentHeader
	result any
}

// page is the page bitmap under construction.
type page struct {
	bm       *Bitmap
	defPixel int
	grow     bool // height unknown: grown by stripes
}

// spend accounts for n units of work (roughly pixels) against maxWork.
func (d *decoder) spend(n int64) error {
	d.work += n
	if d.work > maxWork {
		return errorf("image too large or too complex")
	}
	return nil
}

// newBitmap allocates a white bitmap after checking its size against the limits.
func (d *decoder) newBitmap(w, h int) (*Bitmap, error) {
	if w < 0 || h < 0 || int64(w) > maxDimension || int64(h) > maxDimension || int64(w)*int64(h) > maxPixels {
		return nil, errorf("bitmap size %dx%d out of range", w, h)
	}
	if err := d.spend(int64(w)*int64(h) + 64); err != nil {
		return nil, err
	}
	return newBitmap(w, h), nil
}

// dim converts a width or height read from the data.
func dim(v uint32) (int, error) {
	if int64(v) > maxDimension {
		return 0, errorf("size %d out of range", v)
	}
	return int(v), nil
}

// run processes the segments of a stream.
func (d *decoder) run(data []byte, global bool) error {
	for len(data) > 0 && !d.ended {
		h, n, err := parseSegmentHeader(data)
		if err == errTruncated {
			return nil // a cut-off header or trailing padding
		}
		if err != nil {
			return err
		}
		data = data[n:]
		length := int64(h.length)
		if h.length == unknownLength {
			// Only immediate generic regions may have an unknown length. Some
			// encoders mark the last segment this way, so take the rest.
			length = int64(len(data))
			found := false
			if isGeneric(h.typ) {
				var l int
				if l, found = genericRegionLength(data); found {
					length = int64(l)
				}
			}
			if !found {
				h.length = uint32(length)
			}
		}
		if length > int64(len(data)) {
			length = int64(len(data)) // truncated stream: decode what there is
		}
		body := data[:length]
		data = data[length:]
		if err := d.segment(&segment{segmentHeader: h}, body, global); err != nil {
			return err
		}
	}
	return nil
}

// segment processes one segment.
func (d *decoder) segment(s *segment, data []byte, global bool) error {
	if !global && s.page != 0 && d.pageNum != 0 && s.page != d.pageNum {
		if s.typ == typePageInfo {
			d.ended = true // the next page begins
		}
		return nil
	}
	var err error
	switch s.typ {
	case typeSymbolDict:
		err = d.symbolDictSegment(s, data)
	case typeIntermediateText, typeImmediateText, typeImmediateLosslessText:
		err = d.textRegionSegment(s, data)
	case typePatternDict:
		err = d.patternDictSegment(s, data)
	case typeIntermediateHalftone, typeImmediateHalftone, typeImmediateLosslessHalftone:
		err = d.halftoneRegionSegment(s, data)
	case typeIntermediateGeneric, typeImmediateGeneric, typeImmediateLosslessGeneric:
		err = d.genericRegionSegment(s, data)
	case typeIntermediateRefinement, typeImmediateRefinement, typeImmediateLosslessRefinement:
		err = d.refinementRegionSegment(s, data)
	case typePageInfo:
		err = d.pageInfoSegment(s, data)
	case typeEndOfPage:
		d.ended = true
	case typeEndOfStripe:
		err = d.endOfStripeSegment(data)
	case typeEndOfFile:
		d.ended = true
	case typeTables:
		s.result, err = parseTableSegment(data)
	}
	if err != nil {
		return err
	}
	d.segments[s.number] = s
	return nil
}

// referred returns the results of the segments s refers to.
func (d *decoder) referred(s *segment) []any {
	var r []any
	for _, n := range s.refs {
		if rs := d.segments[n]; rs != nil && rs.result != nil {
			r = append(r, rs.result)
		}
	}
	return r
}

// pageInfoSegment starts the page (7.4.8).
func (d *decoder) pageInfoSegment(s *segment, data []byte) error {
	if d.page != nil {
		d.ended = true
		return nil
	}
	r := reader{b: data}
	w, h := r.u32(), r.u32()
	r.skip(8) // resolution
	flags := r.u8()
	r.u16() // striping
	if r.err != nil {
		return r.err
	}
	width, err := dim(w)
	if err != nil {
		return err
	}
	p := &page{defPixel: int(flags>>2) & 1}
	height := 0
	if h == 0xffffffff {
		p.grow = true
	} else if height, err = dim(h); err != nil {
		return err
	}
	if p.bm, err = d.newBitmap(width, height); err != nil {
		return err
	}
	p.bm.fill(p.defPixel)
	d.page = p
	d.pageNum = s.page
	return nil
}

// endOfStripeSegment extends a page of unknown height to the stripe's end row (7.4.10).
func (d *decoder) endOfStripeSegment(data []byte) error {
	r := reader{b: data}
	end := r.u32()
	if r.err != nil {
		return r.err
	}
	if d.page == nil || !d.page.grow {
		return nil
	}
	return d.growPage(int64(end) + 1)
}

// growPage extends a page of unknown height to h rows.
func (d *decoder) growPage(h int64) error {
	bm := d.page.bm
	if h <= int64(bm.Height) {
		return nil
	}
	if h > maxDimension || int64(bm.Width)*h > maxPixels {
		return errorf("page height %d out of range", h)
	}
	if err := d.spend((h - int64(bm.Height)) * int64(bm.Width)); err != nil {
		return err
	}
	old := len(bm.Data)
	bm.Data = append(bm.Data, make([]byte, int(h)*bm.Stride-old)...)
	if d.page.defPixel != 0 {
		for i := old; i < len(bm.Data); i++ {
			bm.Data[i] = 0xff
		}
	}
	bm.Height = int(h)
	bm.clearPadding()
	return nil
}

// placeRegion keeps the bitmap of an intermediate region segment for later
// refinement, or draws the bitmap of an immediate one onto the page.
func (d *decoder) placeRegion(s *segment, bm *Bitmap, ri regionInfo) error {
	if isIntermediate(s.typ) {
		s.result = bm
		return nil
	}
	if d.page == nil {
		return errorf("region segment %d before the page information", s.number)
	}
	if d.page.grow {
		if err := d.growPage(int64(ri.y) + int64(bm.Height)); err != nil {
			return err
		}
	}
	return d.draw(d.page.bm, bm, int64(ri.x), int64(ri.y), ri.op)
}

// draw combines src into dst at (x, y), accounting for the work.
func (d *decoder) draw(dst, src *Bitmap, x, y int64, op int) error {
	if err := d.spend(int64(src.Width)*int64(src.Height)/4 + 64); err != nil {
		return err
	}
	dst.compose(src, clampCoord(x), clampCoord(y), op)
	return nil
}

// clampCoord limits a coordinate to a range in which any bitmap placed there
// either lies partly on the page or lies clear off it, without overflow.
func clampCoord(v int64) int {
	const lim = 1 << 30
	if v < -lim {
		return -lim
	}
	if v > lim {
		return lim
	}
	return int(v)
}

// ceilLog2 returns the number of bits needed to number n values.
func ceilLog2(n int) int {
	b := 0
	for b < 31 && 1<<b < n {
		b++
	}
	return b
}
