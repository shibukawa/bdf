package rowenc

import (
	"encoding/binary"
	"hash/crc32"
	"image"
	"io"
	"math"

	"github.com/klauspost/compress/zlib"
)

// PNG writes a PNG file from the rows of an image (see Writer).
type PNG struct {
	w   io.Writer
	r   *rower
	z   *zlib.Writer
	out *idatWriter
	// prev is the bytes of the row before, for the filters, and cur
	// those of the row being written with its filter byte; filtered
	// holds each filter tried
	prev, cur []byte
	filtered  [5][]byte
	err       error
	closed    bool
}

// idatBytes is the compressed bytes of an IDAT chunk: the stream is cut
// into chunks of this many.
const idatBytes = 1 << 20

// NewPNG returns a writer of a PNG file of width × height pixels of a
// kind to w (see Options).
func NewPNG(w io.Writer, width, height int, kind Kind, opts *Options) (*PNG, error) {
	r, err := newRower(width, height, kind)
	if err != nil {
		return nil, err
	}
	var o Options
	if opts != nil {
		o = *opts
	}
	level := o.Level
	switch {
	case level == 0:
		level = zlib.DefaultCompression
	case level < zlib.BestSpeed || level > zlib.BestCompression:
		return nil, errBadLevel
	}
	p := &PNG{w: w, r: r}
	if _, err := w.Write([]byte("\x89PNG\r\n\x1a\n")); err != nil {
		return nil, err
	}
	var ihdr [13]byte
	binary.BigEndian.PutUint32(ihdr[0:4], uint32(width))
	binary.BigEndian.PutUint32(ihdr[4:8], uint32(height))
	switch kind {
	case RGBA:
		ihdr[8], ihdr[9] = 8, 6
	case Gray:
		ihdr[8], ihdr[9] = 8, 0
	case Bilevel:
		ihdr[8], ihdr[9] = 1, 3
	}
	if err := writeChunk(w, "IHDR", ihdr[:]); err != nil {
		return nil, err
	}
	if o.DPI > 0 {
		var phys [9]byte
		ppm := uint32(math.Round(o.DPI / 0.0254))
		binary.BigEndian.PutUint32(phys[0:4], ppm)
		binary.BigEndian.PutUint32(phys[4:8], ppm)
		phys[8] = 1
		if err := writeChunk(w, "pHYs", phys[:]); err != nil {
			return nil, err
		}
	}
	if kind == Bilevel {
		if err := writeChunk(w, "PLTE", []byte{0, 0, 0, 255, 255, 255}); err != nil {
			return nil, err
		}
	}
	p.out = &idatWriter{w: w}
	z, err := zlib.NewWriterLevel(p.out, level)
	if err != nil {
		return nil, err
	}
	p.z = z
	n := r.rowBytes()
	p.prev, p.cur = make([]byte, n), make([]byte, n+1)
	for i := range p.filtered {
		p.filtered[i] = make([]byte, n+1)
	}
	return p, nil
}

// WriteRows writes the next rows (see Writer).
func (p *PNG) WriteRows(rows image.Image) error {
	if p.err != nil {
		return p.err
	}
	if p.closed {
		return errClosed
	}
	bpp := 1
	if p.r.kind == RGBA {
		bpp = 4
	}
	err := p.r.each(rows, true, func(row []byte) error {
		f := p.filter(row, bpp)
		if _, err := p.z.Write(f); err != nil {
			return err
		}
		copy(p.prev, row)
		return nil
	})
	if err != nil {
		p.err = err
	}
	return err
}

// filter returns a row with the filter byte before it: the filter whose
// bytes are smallest, as image/png chooses, among None, Sub, Up, Average
// and Paeth. The row before is prev (zeros for the first).
func (p *PNG) filter(row []byte, bpp int) []byte {
	n := len(row)
	prev := p.prev[:n]
	// none
	none := p.filtered[0][:n+1]
	none[0] = 0
	copy(none[1:], row)
	if p.r.kind == Bilevel {
		// packed bits: Up pays off for rows alike, the rest rarely
		up := p.filtered[2][:n+1]
		up[0] = 2
		for i := 0; i < n; i++ {
			up[i+1] = row[i] - prev[i]
		}
		if sumAbs(up[1:]) < sumAbs(none[1:]) {
			return up
		}
		return none
	}
	best, bestSum := none, sumAbs(none[1:])
	try := func(k int, fn func(i int) byte) {
		f := p.filtered[k][:n+1]
		f[0] = byte(k)
		for i := 0; i < n; i++ {
			f[i+1] = fn(i)
		}
		if s := sumAbs(f[1:]); s < bestSum {
			best, bestSum = f, s
		}
	}
	try(1, func(i int) byte {
		if i < bpp {
			return row[i]
		}
		return row[i] - row[i-bpp]
	})
	try(2, func(i int) byte { return row[i] - prev[i] })
	try(3, func(i int) byte {
		var a byte
		if i >= bpp {
			a = row[i-bpp]
		}
		return row[i] - byte((int(a)+int(prev[i]))/2)
	})
	try(4, func(i int) byte {
		var a, c byte
		if i >= bpp {
			a, c = row[i-bpp], prev[i-bpp]
		}
		return row[i] - paeth(a, prev[i], c)
	})
	return best
}

// sumAbs sums the bytes as signed values' magnitudes, as image/png
// weighs a filter.
func sumAbs(b []byte) int {
	s := 0
	for _, v := range b {
		s += int(int8(v)) * sign(int8(v))
	}
	return s
}

func sign(v int8) int {
	if v < 0 {
		return -1
	}
	return 1
}

// paeth is the Paeth predictor of PNG.
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

// Close ends the image (see Writer).
func (p *PNG) Close() error {
	if p.closed {
		return errClosed
	}
	p.closed = true
	if p.err != nil {
		return p.err
	}
	if err := p.r.done(); err != nil {
		return err
	}
	if err := p.z.Close(); err != nil {
		return err
	}
	if err := p.out.flush(); err != nil {
		return err
	}
	return writeChunk(p.w, "IEND", nil)
}

// idatWriter cuts the compressed stream into IDAT chunks.
type idatWriter struct {
	w   io.Writer
	buf []byte
}

func (c *idatWriter) Write(b []byte) (int, error) {
	n := len(b)
	for len(b) > 0 {
		room := idatBytes - len(c.buf)
		take := min(room, len(b))
		c.buf = append(c.buf, b[:take]...)
		b = b[take:]
		if len(c.buf) >= idatBytes {
			if err := c.flush(); err != nil {
				return n - len(b), err
			}
		}
	}
	return n, nil
}

func (c *idatWriter) flush() error {
	if len(c.buf) == 0 {
		return nil
	}
	err := writeChunk(c.w, "IDAT", c.buf)
	c.buf = c.buf[:0]
	return err
}

// writeChunk writes a PNG chunk: length, type, data and CRC.
func writeChunk(w io.Writer, kind string, data []byte) error {
	var head [8]byte
	binary.BigEndian.PutUint32(head[0:4], uint32(len(data)))
	copy(head[4:8], kind)
	crc := crc32.NewIEEE()
	crc.Write(head[4:8])
	crc.Write(data)
	var tail [4]byte
	binary.BigEndian.PutUint32(tail[:], crc.Sum32())
	if _, err := w.Write(head[:]); err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	_, err := w.Write(tail[:])
	return err
}
