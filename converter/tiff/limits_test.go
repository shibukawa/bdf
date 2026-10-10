package tiff

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"slices"
	"testing"

	"github.com/shibukawa/bdf/converter/internal/tiff"
	"github.com/shibukawa/bdf/image/imgconv"
)

// made builds a classic little-endian TIFF: the header, data (at offset 8),
// then an IFD per page of entries {tag, type, count, value or offset}.
func made(data []byte, pages ...[][4]uint32) []byte {
	var b bytes.Buffer
	b.WriteString("II*\x00")
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
	return b.Bytes()
}

// grey returns the entries of a page of w × h 8-bit grey pixels in strips
// of a row each: count of them, their offsets at offset 8 and their byte
// counts after those.
func grey(w, h, compression, count int) [][4]uint32 {
	return [][4]uint32{{tiff.TagImageWidth, 4, 1, uint32(w)}, {tiff.TagImageLength, 4, 1, uint32(h)}, {tiff.TagBitsPerSample, 3, 1, 8},
		{tiff.TagCompression, 3, 1, uint32(compression)}, {tiff.TagPhotometric, 3, 1, tiff.PhotometricBlackIsZero},
		{tiff.TagStripOffsets, 4, uint32(count), 8}, {tiff.TagRowsPerStrip, 4, 1, 1}, {tiff.TagStripByteCounts, 4, uint32(count), uint32(8 + 4*count)}}
}

// TestLimits converts files that state more than they hold: the pages that
// do are left blank or left out, with a warning each.
func TestLimits(t *testing.T) {
	// 16 strips of a byte each, for a page that states 4096 rows of 4096
	// pixels.
	data := make([]byte, 2*4*4096)
	for i := range 4096 {
		binary.LittleEndian.PutUint32(data[4*i:], 8)
		binary.LittleEndian.PutUint32(data[4*(4096+i):], uint32(max(0, 1-i/16)))
	}
	b := made(data, grey(4096, 4096, tiff.CompressionNone, 4096), grey(4, 4, tiff.CompressionNone, 4))
	res, err := Convert(bytes.NewReader(b), int64(len(b)), nil)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"page 1: tiff: 16 bytes of pixel data for a 4096×4096 image; the page is left blank"}
	if p := res.Doc.Views[0].Pages; len(p) != 2 || len(p[0].Layers) != 0 || len(p[1].Layers) != 1 || !slices.Equal(res.Warnings, want) {
		t.Errorf("%d pages, warnings %q", len(p), res.Warnings)
	}

	// 3000 compressed strips that each are the whole file: the file has had
	// its share read before they all are, and the page after is left out.
	data = make([]byte, 2*4*3000)
	for i := range 3000 {
		binary.LittleEndian.PutUint32(data[4*i:], 8)
		binary.LittleEndian.PutUint32(data[4*(3000+i):], 1<<30)
	}
	b = made(data, grey(1, 3000, tiff.CompressionDeflate, 3000), grey(4, 4, tiff.CompressionNone, 4))
	res, err = Convert(bytes.NewReader(b), int64(len(b)), nil)
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"page 1: the pixel data is damaged; what is missing is left blank",
		"tiff: the file has its data read or decoded over and over: page 2 and those after it are left out"}
	if p := res.Doc.Views[0].Pages; len(p) != 1 || res.Pages != 1 || !slices.Equal(res.Warnings, want) {
		t.Errorf("%d pages, warnings %q", len(p), res.Warnings)
	}

	// Forty pages that share a megabyte of zeros deflated to a kilobyte:
	// the file has the pages decoded that its size allows.
	var z bytes.Buffer
	zw, _ := zlib.NewWriterLevel(&z, zlib.BestCompression)
	zw.Write(make([]byte, 1024*1024))
	zw.Close()
	page := grey(1024, 1024, tiff.CompressionDeflate, 1)
	page[5][3], page[6][3], page[7][3] = 8, 1024, uint32(z.Len()) // one strip
	var pages [][][4]uint32
	for range 40 {
		pages = append(pages, page)
	}
	b = made(z.Bytes(), pages...)
	res, err = Convert(bytes.NewReader(b), int64(len(b)), &Options{Images: imgconv.Options{MaxPixels: 64 * 64}})
	if err != nil {
		t.Fatal(err)
	}
	want = []string{"tiff: the file has its data read or decoded over and over: page 11 and those after it are left out"}
	if res.Pages != 10 || !slices.Equal(res.Warnings, want) {
		t.Errorf("%d pages, warnings %q", res.Pages, res.Warnings)
	}
}
