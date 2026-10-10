package tiff

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"runtime"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/image/imgconv"
)

// Files that state more than they hold, or have the same data read over
// and over: the reader checks what a page states against its data, reads
// no more of a tag than it needs, and stops when a file has had its share.

// made builds a classic little-endian TIFF: the header, data (at offset 8),
// then an IFD per page of entries {tag, type, count, value or offset}.
func made(data []byte, pages ...[][4]uint32) *File {
	var b bytes.Buffer
	b.WriteString("II*\x00")
	if len(data)%2 != 0 {
		data = append(data, 0)
	}
	binary.Write(&b, binary.LittleEndian, uint32(8+len(data)))
	b.Write(data)
	for i, entries := range pages {
		binary.Write(&b, binary.LittleEndian, uint16(len(entries)))
		for _, e := range entries {
			binary.Write(&b, binary.LittleEndian, [2]uint16{uint16(e[0]), uint16(e[1])})
			binary.Write(&b, binary.LittleEndian, [2]uint32{e[2], e[3]})
		}
		next := uint32(b.Len() + 4)
		if i == len(pages)-1 {
			next = 0
		}
		binary.Write(&b, binary.LittleEndian, next)
	}
	f, err := Open(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		panic(err)
	}
	return f
}

// grey returns the entries of a page of w × h 8-bit grey pixels in one
// strip of n bytes at offset 8.
func grey(w, h, compression, n int) [][4]uint32 {
	return [][4]uint32{{TagImageWidth, 4, 1, uint32(w)}, {TagImageLength, 4, 1, uint32(h)}, {TagBitsPerSample, 3, 1, 8},
		{TagCompression, 3, 1, uint32(compression)}, {TagPhotometric, 3, 1, PhotometricBlackIsZero},
		{TagStripOffsets, 4, 1, 8}, {TagRowsPerStrip, 4, 1, uint32(h)}, {TagStripByteCounts, 4, 1, uint32(n)}}
}

// allocated returns the bytes f allocates.
func allocated(f func()) uint64 {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	f()
	runtime.ReadMemStats(&after)
	return after.TotalAlloc - before.TotalAlloc
}

// read returns the bytes read from a file since it was opened.
func read(f *File) int64 { return readFactor*f.size + readSlack - f.toRead }

func TestNextEOL(t *testing.T) {
	for _, c := range []struct {
		data  []byte
		start int
		found bool
		end   int
	}{
		{[]byte{0x80, 0x08}, 0, true, 13},       // skip one bit, then 11 zeros and a 1
		{[]byte{0x80, 0x08}, 1, true, 13},       // start at the EOL
		{[]byte{0x80, 0x08}, 2, false, 16},      // ten zeros are too few
		{[]byte{0x40, 0x00, 0x10}, 0, true, 20}, // skip a short run first
		{[]byte{0, 0}, 0, false, 16},            // no terminating 1
		{[]byte{0xff}, 0, false, 8},
	} {
		r := &bitReader{b: c.data, pos: c.start}
		if got := r.nextEOL(); got != c.found || r.pos != c.end {
			t.Errorf("% x from bit %d: EOL %v, then bit %d; want %v, %d", c.data, c.start, got, r.pos, c.found, c.end)
		}
	}
	// A damaged strip of zeros cannot be decoded as a fax image.
	if _, damaged := DecodeFax(make([]byte, 4<<10), false, false, 8, 1); !damaged {
		t.Error("a strip of zeros decoded")
	}
}

func TestPlausible(t *testing.T) {
	// A page of 4096 × 4096 pixels with 16 bytes of them is not decoded.
	d := made(make([]byte, 16), grey(4096, 4096, CompressionNone, 16)).Pages()[0]
	var err error
	n := allocated(func() { _, _, err = d.Decode() })
	if err == nil || !strings.Contains(err.Error(), "16 bytes of pixel data for a 4096×4096 image") || n > 64<<10 {
		t.Errorf("error %v, %d bytes allocated", err, n)
	}
	if _, ok := d.JPEG(); ok {
		t.Error("a JPEG stream of it")
	}
	// A small page that lacks most of its pixels is: they are left blank.
	d = made(make([]byte, 16), grey(100, 100, CompressionNone, 16)).Pages()[0]
	if img, damaged, err := d.Decode(); err != nil || !damaged || img.Bounds().Dx() != 100 {
		t.Errorf("small page: damaged %v, %v", damaged, err)
	}
	// So is a page that compresses as far as deflate goes,
	var z bytes.Buffer
	zw, _ := zlib.NewWriterLevel(&z, zlib.BestCompression)
	zw.Write(make([]byte, 1024*1024))
	zw.Close()
	d = made(z.Bytes(), grey(1024, 1024, CompressionDeflate, z.Len())).Pages()[0]
	if _, damaged, err := d.Decode(); err != nil || damaged {
		t.Errorf("zeros deflated to %d bytes: damaged %v, %v", z.Len(), damaged, err)
	}
	// and a white fax page, whose rows take a bit each in Group 4.
	fax := grey(4000, 800, CompressionG4, 100)
	fax[2][3], fax[4][3] = 1, PhotometricWhiteIsZero
	d = made(bytes.Repeat([]byte{0xff}, 100), fax).Pages()[0]
	if img, damaged, err := d.Decode(); err != nil || damaged || nrgba(img).Pix[0] != 255 {
		t.Errorf("white fax page: damaged %v, %v", damaged, err)
	}
}

func TestTagValues(t *testing.T) {
	// Tags of one value that state as many as the file has bytes: one is
	// read. Here the width, the height and the orientation are 8, the byte
	// at offset 4 of the file (where the first IFD is).
	var pages [][][4]uint32
	const size = 8 + 50*(2+9*12+4)
	for range 50 {
		pages = append(pages, [][4]uint32{{TagNewSubfileType, 1, size - 5, 5}, {TagImageWidth, 1, size - 4, 4}, {TagImageLength, 1, size - 4, 4},
			{TagBitsPerSample, 3, 1, 8}, {TagPhotometric, 3, 1, PhotometricBlackIsZero}, {TagStripOffsets, 4, 1, 8},
			{TagOrientation, 1, size - 4, 4}, {TagRowsPerStrip, 4, 1, 8}, {TagStripByteCounts, 4, 1, 64}})
	}
	f := made(nil, pages...)
	if f.size != size {
		t.Fatalf("file of %d bytes, want %d", f.size, size)
	}
	n := allocated(func() {
		for _, d := range f.Pages() {
			img, damaged, err := d.Decode()
			if w, h := d.Size(); err != nil || damaged || w != 8 || h != 8 || d.Orientation() != 8 || img.Bounds().Dx() != 8 {
				t.Fatalf("page of %d × %d: damaged %v, %v", w, h, damaged, err)
			}
		}
	})
	if got := read(f); len(f.Pages()) != 50 || got > f.size || n > 1<<20 {
		t.Errorf("%d pages: %d bytes read of a file of %d, %d allocated", len(f.Pages()), got, f.size, n)
	}
	// The values of a tag that has them all read are still all there.
	if v := f.IFDs[0].Uints(TagImageWidth); len(v) != size-4 || v[0] != 8 {
		t.Errorf("%d values", len(v))
	}

	// JPEG tables are a few hundred bytes: what states more is not read.
	d := &IFD{f: f, entries: map[uint16]entry{TagJPEGTables: {typ: 7, count: maxTables + 1, off: 8}}}
	before := f.toRead
	if tables, ok := (&layout{d: d}).tables(); ok || tables != nil || f.toRead != before {
		t.Errorf("tables of %d bytes read", maxTables+1)
	}
}

func TestSharedData(t *testing.T) {
	// Strips without compression: what a strip states beyond its rows is
	// not read.
	const strips = 3000
	data := make([]byte, 2*4*strips)
	for i := range strips {
		binary.LittleEndian.PutUint32(data[4*i:], 8)              // offset
		binary.LittleEndian.PutUint32(data[4*(strips+i):], 1<<30) // byte count
	}
	page := grey(1, strips, CompressionNone, 0)
	page[5], page[6], page[7] = [4]uint32{TagStripOffsets, 4, strips, 8}, [4]uint32{TagRowsPerStrip, 4, 1, 1},
		[4]uint32{TagStripByteCounts, 4, strips, 8 + 4*strips}
	f := made(data, page)
	if _, damaged, err := f.Pages()[0].Decode(); err != nil || damaged || read(f) > 2*f.size {
		t.Errorf("uncompressed: damaged %v, %v, %d bytes read of a file of %d", damaged, err, read(f), f.size)
	}

	// Compressed strips that each are the whole file: the file has had its
	// share read before they all are, the rest of the page is blank and
	// the next page is not decoded.
	page[3][3] = CompressionDeflate
	f = made(data, page, grey(8, 8, CompressionNone, 64))
	var damaged bool
	var err error
	n := allocated(func() { _, damaged, err = f.Pages()[0].Decode() })
	if err != nil || !damaged || !f.Spent() || n > 4<<20 {
		t.Errorf("compressed: damaged %v, %v, spent %v, %d bytes allocated", damaged, err, f.Spent(), n)
	}
	if _, _, err := f.Pages()[1].Decode(); !errors.Is(err, ErrSpent) {
		t.Errorf("next page: %v", err)
	}
	if _, ok := f.Pages()[1].JPEG(); ok {
		t.Error("next page: a JPEG stream")
	}
}

func TestDecodeBudget(t *testing.T) {
	// Forty pages that share a megabyte of zeros deflated to a kilobyte:
	// the file has the pages decoded that its size allows.
	var z bytes.Buffer
	zw, _ := zlib.NewWriterLevel(&z, zlib.BestCompression)
	zw.Write(make([]byte, 1024*1024))
	zw.Close()
	var pages [][][4]uint32
	for range 40 {
		pages = append(pages, grey(1024, 1024, CompressionDeflate, z.Len()))
	}
	f := made(z.Bytes(), pages...)
	decoded := 0
	for _, d := range f.Pages() {
		if _, _, err := d.Decode(); err == nil {
			decoded++
		} else if !errors.Is(err, ErrSpent) {
			t.Fatal(err)
		}
	}
	if want := int((maxRatio*f.size + decodeSlack) >> 20); decoded != want || !f.Spent() {
		t.Errorf("%d pages of a file of %d bytes decoded, want %d", decoded, f.size, want)
	}
}

func TestPictureSize(t *testing.T) {
	old := MaxPicturePixels
	defer func() { MaxPicturePixels = old }()
	b := open(t, "rgb.tif")
	w, h := b.Pages()[0].Size()
	data := make([]byte, b.size)
	b.r.ReadAt(data, 0)
	MaxPicturePixels = int64(w*h) - 1
	if _, _, err := Picture(data, imgconv.Options{}); err == nil || !strings.Contains(err.Error(), "is too large") {
		t.Errorf("picture of %d×%d pixels: %v", w, h, err)
	}
	MaxPicturePixels++
	if _, _, err := Picture(data, imgconv.Options{}); err != nil {
		t.Error(err)
	}
}
