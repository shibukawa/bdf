package rowenc

import (
	"encoding/binary"
	"errors"
	"image"
	"io"
	"math"
)

var (
	errClosed   = errors.New("rowenc: the writer is closed")
	errBadLevel = errors.New("rowenc: the compression level must be 1 to 9")
)

// TIFF writes a TIFF file from the rows of an image (see Writer): not
// compressed, in one strip, with the directory before the pixels, so
// that the file is written in order without seeking. The pixels are 8
// bits a sample (RGB with an associated alpha sample for RGBA, gray), or
// a bit a pixel with 0 black for Bilevel.
type TIFF struct {
	w      io.Writer
	r      *rower
	err    error
	closed bool
}

// NewTIFF returns a writer of a TIFF file of width × height pixels of a
// kind to w (see Options; the compression level does not apply).
func NewTIFF(w io.Writer, width, height int, kind Kind, opts *Options) (*TIFF, error) {
	r, err := newRower(width, height, kind)
	if err != nil {
		return nil, err
	}
	dpi := 72.0
	if opts != nil && opts.DPI > 0 {
		dpi = opts.DPI
	}
	if int64(r.rowBytes())*int64(height) > math.MaxUint32-1024 {
		return nil, errors.New("rowenc: the image is too large for a TIFF file")
	}
	// the directory: its entries in the order of their tags, with the
	// values that do not fit in 4 bytes after it
	type entry struct {
		tag, typ uint16
		count    uint32
		value    []byte // the value, or the bytes it points to
	}
	const (
		typShort    = 3
		typLong     = 4
		typRational = 5
	)
	short := func(v ...uint16) []byte {
		b := make([]byte, 2*len(v))
		for i, x := range v {
			binary.LittleEndian.PutUint16(b[2*i:], x)
		}
		return b
	}
	long := func(v uint32) []byte {
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, v)
		return b
	}
	rational := func(dpi float64) []byte {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint32(b[0:4], uint32(math.Round(dpi*1000)))
		binary.LittleEndian.PutUint32(b[4:8], 1000)
		return b
	}
	samples, bits, photometric := uint16(1), []uint16{8}, uint16(1) // min-is-black
	switch kind {
	case RGBA:
		samples, bits, photometric = 4, []uint16{8, 8, 8, 8}, 2 // RGB
	case Bilevel:
		bits = []uint16{1}
	}
	strip := uint32(r.rowBytes()) * uint32(height)
	entries := []entry{
		{256, typLong, 1, long(uint32(width))},
		{257, typLong, 1, long(uint32(height))},
		{258, typShort, uint32(len(bits)), short(bits...)},
		{259, typShort, 1, short(1)}, // no compression
		{262, typShort, 1, short(photometric)},
		{273, typLong, 1, nil}, // strip offset: after the directory
		{277, typShort, 1, short(samples)},
		{278, typLong, 1, long(uint32(height))},
		{279, typLong, 1, long(strip)},
		{282, typRational, 1, rational(dpi)},
		{283, typRational, 1, rational(dpi)},
		{284, typShort, 1, short(1)}, // chunky
		{296, typShort, 1, short(2)}, // inches
	}
	if kind == RGBA {
		entries = append(entries, entry{338, typShort, 1, short(1)}) // associated alpha
	}
	// the header, the directory, the values after it, then the pixels
	ifd := uint32(8)
	after := ifd + 2 + 12*uint32(len(entries)) + 4
	var tail []byte
	data := after
	for _, e := range entries {
		if len(e.value) > 4 {
			data += uint32(len(e.value))
		}
	}
	buf := make([]byte, 0, data)
	buf = append(buf, 'I', 'I', 42, 0)
	buf = binary.LittleEndian.AppendUint32(buf, ifd)
	buf = binary.LittleEndian.AppendUint16(buf, uint16(len(entries)))
	for _, e := range entries {
		buf = binary.LittleEndian.AppendUint16(buf, e.tag)
		buf = binary.LittleEndian.AppendUint16(buf, e.typ)
		buf = binary.LittleEndian.AppendUint32(buf, e.count)
		switch {
		case e.tag == 273:
			buf = binary.LittleEndian.AppendUint32(buf, data)
		case len(e.value) > 4:
			buf = binary.LittleEndian.AppendUint32(buf, after+uint32(len(tail)))
			tail = append(tail, e.value...)
		default:
			var v [4]byte
			copy(v[:], e.value)
			buf = append(buf, v[:]...)
		}
	}
	buf = binary.LittleEndian.AppendUint32(buf, 0) // no directory after
	buf = append(buf, tail...)
	if _, err := w.Write(buf); err != nil {
		return nil, err
	}
	return &TIFF{w: w, r: r}, nil
}

// WriteRows writes the next rows (see Writer).
func (t *TIFF) WriteRows(rows image.Image) error {
	if t.err != nil {
		return t.err
	}
	if t.closed {
		return errClosed
	}
	err := t.r.each(rows, false, func(row []byte) error {
		_, err := t.w.Write(row)
		return err
	})
	if err != nil {
		t.err = err
	}
	return err
}

// Close ends the image (see Writer): nothing follows the pixels.
func (t *TIFF) Close() error {
	if t.closed {
		return errClosed
	}
	t.closed = true
	if t.err != nil {
		return t.err
	}
	return t.r.done()
}
