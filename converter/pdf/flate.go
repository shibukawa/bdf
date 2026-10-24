package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/flate"
)

// FlateDecode is read here, not by pdfcpu's filter. That filter makes a zlib
// reader for every stream and copies what it reads through a buffer that
// grows as it goes, which was four tenths of the bytes a conversion
// allocated. Here the inflaters (klauspost's, as everywhere else) are
// pooled, one holding its 32 KiB window and tables, and the bytes are read
// straight into the slice returned.
//
// Streams are read as pdfcpu read them, with three differences that only a
// damaged file can show: the Adler-32 checksum of the zlib wrapper is not
// checked (the bytes were read; a wrong checksum, which some writers leave,
// is no reason to drop them), data without the zlib header is read as bare
// deflate, and a predictor's last row cut short is dropped instead of
// failing the stream. As before, a stream cut short of its final block
// gives the bytes before the cut, and corrupt data fails the stream.

var inflaters sync.Pool

// inflater is a pooled flate reader with the reader it reads from, so that
// a Reset allocates nothing.
type inflater struct {
	in bytes.Reader
	r  io.ReadCloser
}

// inflateUpTo decodes the data of a FlateDecode stream to at most max+1
// bytes, applying the predictor of parms (checked by checkPredictor): more
// than max bytes tells the caller the stream goes on. With sized set, max
// is what the stream is expected to hold (the samples of an image), and
// the output is allocated for it at once.
func inflateUpTo(data []byte, parms map[string]int, max int64, sized bool) ([]byte, error) {
	if max < 0 {
		return nil, nil
	}
	limit := max + 1
	p, err := newPredictor(parms)
	if err != nil {
		return nil, err
	}
	if p.kind == predictNone {
		return inflate(data, limit, sized)
	}
	// Rows are read whole, and the output must reach limit bytes when the
	// stream has them: the raw data of as many rows as that takes.
	rows := (limit + int64(p.rowSize) - 1) / int64(p.rowSize)
	raw, err := inflate(data, rows*int64(p.rowLen), sized)
	if err != nil {
		return nil, err
	}
	out, err := p.undo(raw)
	if err != nil {
		return nil, err
	}
	if int64(len(out)) > limit {
		out = out[:limit]
	}
	return out, nil
}

// inflate decodes zlib or bare deflate data to at most limit bytes, which
// is the size expected when sized is set.
func inflate(data []byte, limit int64, sized bool) ([]byte, error) {
	if len(data) == 0 {
		return nil, io.ErrUnexpectedEOF
	}
	src := data
	// A zlib header: deflate with a window of at most 32 KiB, and the pair
	// a multiple of 31. Anything else is taken for bare deflate.
	if len(src) >= 2 && src[0]&0x0f == 8 && src[0]>>4 <= 7 && (uint16(src[0])<<8|uint16(src[1]))%31 == 0 {
		if src[1]&0x20 != 0 {
			return nil, errors.New("zlib data with a preset dictionary")
		}
		src = src[2:]
	}
	inf, _ := inflaters.Get().(*inflater)
	if inf == nil {
		inf = &inflater{}
		inf.in.Reset(src)
		inf.r = flate.NewReader(&inf.in)
	} else {
		inf.in.Reset(src)
		if err := inf.r.(flate.Resetter).Reset(&inf.in, nil); err != nil {
			return nil, err
		}
	}
	defer func() {
		inf.in.Reset(nil) // not the stream's bytes, while pooled
		inf.r.(flate.Resetter).Reset(&inf.in, nil)
		inflaters.Put(inf)
	}()
	// The output starts at the size of the input, which data that does
	// not compress has. When that fills, the ratio so far says how much the
	// stream holds, so that the second slice is usually the last: doubling
	// the slice each time it filled was most of the bytes a conversion
	// allocated, and left the last one up to half empty.
	first := int64(len(src)) + 1024
	if sized {
		first = limit
	}
	out := make([]byte, 0, min(limit, first))
	for int64(len(out)) < limit {
		if len(out) == cap(out) {
			consumed := max(len(src)-inf.in.Len(), 1)
			projected := int64(len(out)) * int64(len(src)) / int64(consumed)
			want := max(projected+projected/8+1024, int64(cap(out))+int64(cap(out))/2)
			grown := make([]byte, len(out), min(limit, want))
			copy(grown, out)
			out = grown
		}
		space := out[len(out):cap(out)]
		if rest := limit - int64(len(out)); int64(len(space)) > rest {
			space = space[:rest]
		}
		n, err := inf.r.Read(space)
		out = out[:len(out)+n]
		if err == io.EOF || errors.Is(err, io.ErrUnexpectedEOF) {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// Predictors of FlateDecode and LZWDecode (7.4.4.4): TIFF's horizontal
// differencing, and PNG's filters, chosen by the byte that starts each row.
const (
	predictNone = iota
	predictTIFF
	predictPNG
)

type predictor struct {
	kind                 int
	colors, bpc, columns int
	bpp                  int // bytes per pixel, at least one, for the PNG filters
	rowSize              int // bytes of a row of samples
	rowLen               int // bytes of a row as stored: a filter byte more for PNG
}

// newPredictor reads the predictor parameters, which checkPredictor has
// bounded. Predictor values other than 1, 2 and 10 to 15 are an error, as
// they were for pdfcpu.
func newPredictor(parms map[string]int) (*predictor, error) {
	get := func(key string, def int) int {
		if v, ok := parms[key]; ok {
			return v
		}
		return def
	}
	p := &predictor{colors: get("Colors", 1), bpc: get("BitsPerComponent", 8), columns: get("Columns", 1)}
	switch pr := get("Predictor", 1); {
	case pr <= 1:
		return p, nil
	case pr == 2:
		p.kind = predictTIFF
	case pr >= 10 && pr <= 15:
		p.kind = predictPNG
	default:
		return nil, fmt.Errorf("predictor %d", pr)
	}
	p.bpp = (p.bpc*p.colors + 7) / 8
	p.rowSize = (p.bpc*p.colors*p.columns + 7) / 8
	p.rowLen = p.rowSize
	if p.kind == predictPNG {
		p.rowLen++
	}
	return p, nil
}

// undo reverses the prediction of the whole rows of raw, in place, and
// returns them. A last row cut short is left out.
func (p *predictor) undo(raw []byte) ([]byte, error) {
	rows := len(raw) / p.rowLen
	if p.kind == predictTIFF {
		for r := 0; r < rows; r++ {
			p.undoTIFF(raw[r*p.rowSize : (r+1)*p.rowSize])
		}
		return raw[:rows*p.rowSize], nil
	}
	var prev []byte // the row above, reconstructed; nil above the first
	for r := 0; r < rows; r++ {
		row := raw[r*p.rowLen : (r+1)*p.rowLen]
		cur := row[1:]
		if err := p.undoPNG(row[0], cur, prev); err != nil {
			return nil, err
		}
		// The rows close up as the filter bytes go: the row's place is at
		// or before its samples, and the row above is before it.
		prev = raw[r*p.rowSize : (r+1)*p.rowSize]
		copy(prev, cur)
	}
	return raw[:rows*p.rowSize], nil
}

// undoTIFF adds each component to the same component of the sample before
// it, modulo 2^bpc: 16-bit components are big-endian, and narrower ones are
// packed from the most significant bit, their padding bits left alone.
func (p *predictor) undoTIFF(row []byte) {
	switch p.bpc {
	case 8:
		for i := p.colors; i < len(row); i++ {
			row[i] += row[i-p.colors]
		}
	case 16:
		for i := 2 * p.colors; i+1 < len(row); i += 2 {
			v := uint16(row[i])<<8 | uint16(row[i+1])
			v += uint16(row[i-2*p.colors])<<8 | uint16(row[i-2*p.colors+1])
			row[i], row[i+1] = byte(v>>8), byte(v)
		}
	default:
		mask := byte(1<<p.bpc - 1)
		get := func(k int) byte {
			bit := k * p.bpc
			return row[bit/8] >> (8 - p.bpc - bit%8) & mask
		}
		for k, n := p.colors, p.colors*p.columns; k < n; k++ {
			bit := k * p.bpc
			shift := 8 - p.bpc - bit%8
			v := (get(k) + get(k-p.colors)) & mask
			row[bit/8] = row[bit/8]&^(mask<<shift) | v<<shift
		}
	}
}

// undoPNG reconstructs a row of samples from its filter type and the row
// above (nil for a row of zeros).
func (p *predictor) undoPNG(filter byte, cur, prev []byte) error {
	bpp := p.bpp
	switch filter {
	case 0: // None
	case 1: // Sub
		for i := bpp; i < len(cur); i++ {
			cur[i] += cur[i-bpp]
		}
	case 2: // Up
		for i := range prev {
			cur[i] += prev[i]
		}
	case 3: // Average
		if prev == nil {
			for i := bpp; i < len(cur); i++ {
				cur[i] += cur[i-bpp] / 2
			}
			break
		}
		for i := 0; i < bpp && i < len(cur); i++ {
			cur[i] += prev[i] / 2
		}
		for i := bpp; i < len(cur); i++ {
			cur[i] += byte((int(cur[i-bpp]) + int(prev[i])) / 2)
		}
	case 4: // Paeth
		if prev == nil {
			for i := bpp; i < len(cur); i++ {
				cur[i] += cur[i-bpp] // paeth(a, 0, 0) is a
			}
			break
		}
		for i := 0; i < bpp && i < len(cur); i++ {
			cur[i] += prev[i] // paeth(0, b, 0) is b
		}
		for i := bpp; i < len(cur); i++ {
			cur[i] += paeth(cur[i-bpp], prev[i], prev[i-bpp])
		}
	default:
		return fmt.Errorf("PNG filter type %d", filter)
	}
	return nil
}

// paeth is the Paeth predictor of PNG: whichever of the left, above and
// upper-left bytes is nearest to left+above-upperLeft.
func paeth(a, b, c byte) byte {
	pa, pb, pc := int(b)-int(c), int(a)-int(c), 0
	pc = pa + pb
	if pa < 0 {
		pa = -pa
	}
	if pb < 0 {
		pb = -pb
	}
	if pc < 0 {
		pc = -pc
	}
	if pa <= pb && pa <= pc {
		return a
	}
	if pb <= pc {
		return b
	}
	return c
}
