package rowenc

import (
	"bytes"
	"encoding/binary"
	"hash/adler32"
	"hash/crc32"
	"image"
	"io"
	"math"
	"runtime"
	"sync"

	"github.com/klauspost/compress/flate"
)

// PNG writes a PNG file from the rows of an image (see Writer).
//
// The rows are filtered and compressed in chunks of a fixed number of
// rows, each by a goroutine (see Options.Threads): the deflate stream of
// the image is the streams of the chunks one after the other, each
// ended with a sync flush, so the file is the same however the rows come
// and however many goroutines write it.
type PNG struct {
	w   io.Writer
	r   *rower
	out *idatWriter
	err error

	level, threads int
	bpp            int // bytes a pixel for the filters, 1 at least
	chunkRows      int // the rows of a chunk
	chunks         int // the chunks begun

	// cur is the raw rows of the chunk being gathered, after the raw row
	// before it (for the filters), and free the buffers of the chunks
	// written, for the next ones
	cur  *chunk
	free chan *chunk
	// jobs takes the chunks to compress and slots the compressed chunk
	// of each, in order; sink writes them to the file in order and
	// reports to done
	jobs  chan *chunk
	slots []chan *chunk
	wg    sync.WaitGroup
	done  chan error
	quit  chan struct{}
	// adler is the Adler-32 of the filtered rows written so far
	adler  uint32
	closed bool
}

// chunk is the rows of a chunk: raw, then filtered and compressed.
type chunk struct {
	index int
	// prev is the raw row before the chunk (zeros before the first), raw
	// the raw rows of the chunk, and rows how many there are
	prev, raw []byte
	rows      int
	// filtered holds the rows with their filter bytes, and each filter
	// tried; packed is the compressed stream of the chunk, ending with a
	// sync flush, and adler the Adler-32 of filtered
	filtered []byte
	try      [5][]byte
	packed   bytes.Buffer
	adler    uint32
	err      error
}

// idatBytes is the compressed bytes of an IDAT chunk: the stream is cut
// into chunks of this many.
const idatBytes = 1 << 20

// chunkBytes is the raw bytes of the rows of a chunk compressed as one,
// about: a chunk is as many whole rows as fit, one at least.
const chunkBytes = 256 << 10

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
		level = flate.DefaultCompression
	case level < flate.BestSpeed || level > flate.BestCompression:
		return nil, errBadLevel
	}
	threads := o.Threads
	if threads <= 0 {
		threads = min(maxThreads, runtime.GOMAXPROCS(0))
	}
	threads = max(1, min(threads, maxThreads))
	p := &PNG{w: w, r: r, level: level, threads: threads, bpp: 1}
	if kind == RGBA {
		p.bpp = 4
	}
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
	// the zlib header: deflate, a 32 KiB window, no preset dictionary
	if _, err := p.out.Write([]byte{0x78, 0x9c}); err != nil {
		return nil, err
	}
	p.adler = 1
	n := r.rowBytes()
	p.chunkRows = max(1, chunkBytes/(n+1))
	p.chunks = 0
	// as many chunks in flight as there are goroutines, and one gathered
	p.free = make(chan *chunk, threads+1)
	for i := 0; i < threads+1; i++ {
		p.free <- &chunk{prev: make([]byte, n), raw: make([]byte, 0, p.chunkRows*n)}
	}
	p.jobs = make(chan *chunk, threads)
	p.slots = make([]chan *chunk, 0, (height+p.chunkRows-1)/p.chunkRows)
	for i := 0; i < cap(p.slots); i++ {
		p.slots = append(p.slots, make(chan *chunk, 1))
	}
	p.done = make(chan error, 1)
	p.quit = make(chan struct{})
	for i := 0; i < threads; i++ {
		p.wg.Add(1)
		go p.compress()
	}
	go p.sink()
	return p, nil
}

// maxThreads bounds the goroutines of a PNG writer.
const maxThreads = 8

// WriteRows writes the next rows (see Writer).
func (p *PNG) WriteRows(rows image.Image) error {
	if p.err != nil {
		return p.err
	}
	if p.closed {
		return errClosed
	}
	err := p.r.each(rows, true, func(row []byte) error {
		if p.cur == nil {
			p.cur = <-p.free
			p.cur.index = p.chunks
			p.cur.raw, p.cur.rows = p.cur.raw[:0], 0
			p.chunks++
		}
		c := p.cur
		c.raw = append(c.raw, row...)
		c.rows++
		if c.rows == p.chunkRows {
			p.dispatch(c)
		}
		return nil
	})
	if err != nil {
		p.err = err
	}
	return err
}

// dispatch hands the chunk gathered to the goroutines, with the raw row
// before the next one.
func (p *PNG) dispatch(c *chunk) {
	n := p.r.rowBytes()
	next := <-p.free
	copy(next.prev, c.raw[len(c.raw)-n:])
	next.index = p.chunks
	next.raw, next.rows = next.raw[:0], 0
	p.chunks++
	p.jobs <- c
	p.cur = next
}

// compress filters and compresses the chunks it takes, and puts each in
// its slot (see PNG).
func (p *PNG) compress() {
	defer p.wg.Done()
	fw, _ := flate.NewWriter(nil, p.level)
	for c := range p.jobs {
		p.filterChunk(c)
		c.packed.Reset()
		fw.Reset(&c.packed)
		if _, err := fw.Write(c.filtered); err != nil {
			c.err = err
		} else if err := fw.Flush(); err != nil {
			c.err = err
		}
		c.adler = adler32.Checksum(c.filtered)
		p.slots[c.index] <- c
	}
}

// filterChunk filters the raw rows of a chunk: each with the filter whose
// bytes are smallest, as image/png chooses, among None, Sub, Up, Average
// and Paeth (None and Up for packed bits, which the rest rarely help).
func (p *PNG) filterChunk(c *chunk) {
	n := p.r.rowBytes()
	for i := range c.try {
		if cap(c.try[i]) < n+1 {
			c.try[i] = make([]byte, n+1)
		}
	}
	c.filtered = c.filtered[:0]
	prev := c.prev
	for i := 0; i < c.rows; i++ {
		row := c.raw[i*n : (i+1)*n]
		c.filtered = append(c.filtered, p.filter(c, row, prev)...)
		prev = row
	}
}

// filter returns a row with the filter byte before it, in the buffers of
// the chunk; prev is the row before (zeros for the first).
func (p *PNG) filter(c *chunk, row, prev []byte) []byte {
	n, bpp := len(row), p.bpp
	prev = prev[:n]
	// none
	none := c.try[0][:n+1]
	none[0] = 0
	copy(none[1:], row)
	if p.r.kind == Bilevel {
		// packed bits: Up pays off for rows alike, the rest rarely
		up := c.try[2][:n+1]
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
	// sub
	f := c.try[1][:n+1]
	f[0] = 1
	copy(f[1:1+bpp], row[:bpp])
	for i := bpp; i < n; i++ {
		f[i+1] = row[i] - row[i-bpp]
	}
	if s := sumAbs(f[1:]); s < bestSum {
		best, bestSum = f, s
	}
	// up
	f = c.try[2][:n+1]
	f[0] = 2
	for i := 0; i < n; i++ {
		f[i+1] = row[i] - prev[i]
	}
	if s := sumAbs(f[1:]); s < bestSum {
		best, bestSum = f, s
	}
	// average
	f = c.try[3][:n+1]
	f[0] = 3
	for i := 0; i < bpp; i++ {
		f[i+1] = row[i] - prev[i]/2
	}
	for i := bpp; i < n; i++ {
		f[i+1] = row[i] - byte((int(row[i-bpp])+int(prev[i]))/2)
	}
	if s := sumAbs(f[1:]); s < bestSum {
		best, bestSum = f, s
	}
	// paeth
	f = c.try[4][:n+1]
	f[0] = 4
	for i := 0; i < bpp; i++ {
		f[i+1] = row[i] - paeth(0, prev[i], 0)
	}
	for i := bpp; i < n; i++ {
		f[i+1] = row[i] - paeth(row[i-bpp], prev[i], prev[i-bpp])
	}
	if s := sumAbs(f[1:]); s < bestSum {
		best, bestSum = f, s
	}
	return best
}

// sink writes the compressed chunks to the file in order, and sums their
// Adler-32 into the image's (see adlerJoin); it reports when the slots
// are all taken, or at the first error.
func (p *PNG) sink() {
	for i := range p.slots {
		var c *chunk
		select {
		case c = <-p.slots[i]:
		case <-p.quit:
			return
		}
		if c.err != nil {
			p.done <- c.err
			return
		}
		if _, err := p.out.Write(c.packed.Bytes()); err != nil {
			p.done <- err
			return
		}
		p.adler = adlerJoin(p.adler, c.adler, len(c.filtered))
		p.free <- c
	}
	p.done <- nil
}

// adlerJoin returns the Adler-32 of the bytes of a and then those of b,
// from the Adler-32 of each and the length of b (as zlib's
// adler32_combine does).
func adlerJoin(a, b uint32, lenB int) uint32 {
	const mod = 65521
	a1, b1 := uint64(a&0xffff), uint64(a>>16)
	a2, b2 := uint64(b&0xffff), uint64(b>>16)
	rem := uint64(lenB % mod)
	sa := (a1 + a2 + mod - 1) % mod
	sb := (b1 + b2 + rem*(a1+mod-1)) % mod
	return uint32(sb<<16 | sa)
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
	if p.err == nil {
		p.err = p.r.done()
	}
	if p.err != nil {
		// the goroutines end with their chunks, and the sink with them
		close(p.jobs)
		p.wg.Wait()
		close(p.quit)
		return p.err
	}
	if p.cur != nil && p.cur.rows > 0 {
		p.jobs <- p.cur
		p.cur = nil
	}
	close(p.jobs)
	p.wg.Wait()
	if err := <-p.done; err != nil {
		close(p.quit)
		return err
	}
	// the final block of the stream, empty, and the Adler-32 of the rows
	var tail bytes.Buffer
	fw, _ := flate.NewWriter(&tail, p.level)
	if err := fw.Close(); err != nil {
		return err
	}
	var sum [4]byte
	binary.BigEndian.PutUint32(sum[:], p.adler)
	tail.Write(sum[:])
	if _, err := p.out.Write(tail.Bytes()); err != nil {
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
