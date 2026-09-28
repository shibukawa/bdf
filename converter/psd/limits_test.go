package psd

import (
	"bytes"
	"cmp"
	"compress/zlib"
	"encoding/binary"
	"math"
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/converter"
)

// Files that state sizes their data does not bear out, and documents whose
// layers would take more memory and work than a picture is worth: the
// converter checks the data before it makes room, and leaves out with a
// warning what goes beyond its limits.

func be16(v int) []byte { return binary.BigEndian.AppendUint16(nil, uint16(v)) }
func be32(v int) []byte { return binary.BigEndian.AppendUint32(nil, uint32(v)) }

// testLayer is a layer of a document made by layered: a rectangle of one
// colour, the bounding divider of a group (section 3) or a group (section
// 1), which may be clipped to the layer below or an artboard (top, left,
// bottom, right). shade, when set, gives the samples of the colour
// channels instead of rgb, and alpha those of a transparency channel, for
// the pixel at x, y of the canvas. An opacity of 0 is 255, a blend mode
// of "" is normal.
type testLayer struct {
	top, left, bottom, right int
	rgb                      [3]uint8
	section                  int
	clip                     bool
	artboard                 *[4]int
	shade                    func(c, x, y int) uint8
	alpha                    func(x, y int) uint8
	opacity                  uint8
	blend                    string
}

// id is a key or class ID of a descriptor.
func id(s string) []byte {
	if len(s) == 4 {
		return append(be32(0), s...)
	}
	return append(be32(len(s)), s...)
}

// artboardInfo is the data of an artb block: a version and a descriptor
// whose artboardRect holds the sides.
func artboardInfo(r [4]int) []byte {
	b := slices.Concat(be32(16), be32(0), id("artboard"), be32(1), // version; name, class and items of the descriptor
		id("artboardRect"), []byte("Objc"), be32(0), id("classFloatRect"), be32(4))
	for i, key := range []string{"Top ", "Left", "Btom", "Rght"} {
		b = slices.Concat(b, id(key), []byte("doub"), binary.BigEndian.AppendUint64(nil, math.Float64bits(float64(r[i]))))
	}
	return b
}

// layered makes an 8-bit RGB document of w × h pixels without a composite.
func layered(w, h int, layers []testLayer) []byte {
	var b, rec, data bytes.Buffer
	b.WriteString("8BPS\x00\x01\x00\x00\x00\x00\x00\x00")
	b.Write(be16(3))
	b.Write(be32(h))
	b.Write(be32(w))
	b.Write(be16(8))
	b.Write(be16(modeRGB))
	b.Write(be32(0)) // colour mode data
	// image resources: the version information, without real merged data
	version := []byte{0, 0, 0, 1, 0, 0}
	b.Write(be32(12 + len(version)))
	b.WriteString("8BIM\x04\x21\x00\x00")
	b.Write(be32(len(version)))
	b.Write(version)
	for _, l := range layers {
		for _, v := range []int{l.top, l.left, l.bottom, l.right} {
			rec.Write(be32(v))
		}
		n := (l.bottom - l.top) * (l.right - l.left)
		if n == 0 {
			rec.Write(be16(0))
		} else {
			ids := []int{0, 1, 2}
			if l.alpha != nil {
				ids = append(ids, chAlpha)
			}
			rec.Write(be16(len(ids)))
			for _, c := range ids {
				rec.Write(be16(c))
				rec.Write(be32(2 + n))
				data.Write(be16(compRaw))
				for y := l.top; y < l.bottom; y++ {
					for x := l.left; x < l.right; x++ {
						switch {
						case c == chAlpha:
							data.WriteByte(l.alpha(x, y))
						case l.shade != nil:
							data.WriteByte(l.shade(c, x, y))
						default:
							data.WriteByte(l.rgb[c])
						}
					}
				}
			}
		}
		rec.WriteString("8BIM")
		rec.WriteString(cmp.Or(l.blend, "norm"))
		rec.WriteByte(cmp.Or(l.opacity, 255))
		if l.clip {
			rec.WriteString("\x01\x00\x00")
		} else {
			rec.WriteString("\x00\x00\x00")
		}
		extra := slices.Concat(be32(0), be32(0), []byte{0, 0, 0, 0}) // mask, blending ranges, name
		if l.section != 0 {
			extra = slices.Concat(extra, []byte("8BIMlsct"), be32(4), be32(l.section))
		}
		if l.artboard != nil {
			info := artboardInfo(*l.artboard)
			extra = slices.Concat(extra, []byte("8BIMartb"), be32(len(info)), info)
		}
		rec.Write(be32(len(extra)))
		rec.Write(extra)
	}
	info := slices.Concat(be16(len(layers)), rec.Bytes(), data.Bytes())
	if len(info)%2 != 0 {
		info = append(info, 0)
	}
	b.Write(be32(4 + len(info) + 4))
	b.Write(be32(len(info)))
	b.Write(info)
	b.Write(be32(0)) // global layer mask information
	return b.Bytes()
}

// limit sets a limit of the converter for a test.
func limit[T any](t *testing.T, v *T, n T) {
	t.Helper()
	old := *v
	*v = n
	t.Cleanup(func() { *v = old })
}

func hasWarning(warnings []string, part string) bool {
	return slices.ContainsFunc(warnings, func(w string) bool { return strings.Contains(w, part) })
}

// red is a red square in the middle of a canvas of 8 × 8 pixels.
var red = testLayer{top: 2, left: 2, bottom: 6, right: 6, rgb: [3]uint8{255, 0, 0}}

// nested makes a document of 8 × 8 pixels with the red square inside groups,
// one in the other.
func nested(depth int) []byte {
	var ls []testLayer
	for range depth {
		ls = append(ls, testLayer{section: 3})
	}
	ls = append(ls, red)
	for range depth {
		ls = append(ls, testLayer{section: 1})
	}
	return layered(8, 8, ls)
}

func TestCompositeLimits(t *testing.T) {
	c := convert(t, nested(maxGroupNest), nil)
	if px := rgbaAt(c.pages[0], 3, 3); px != [4]uint8{255, 0, 0, 255} || len(c.res.Warnings) != 1 {
		t.Errorf("%d groups: %v at 3,3, warnings %v", maxGroupNest, px, c.res.Warnings)
	}
	// One group more: it is left out, and no buffer is made for it.
	c = convert(t, nested(maxGroupNest+1), nil)
	if px := rgbaAt(c.pages[0], 3, 3); px[3] != 0 || !hasWarning(c.res.Warnings, "groups more than 16 deep are not drawn") {
		t.Errorf("%d groups: %v at 3,3, warnings %v", maxGroupNest+1, px, c.res.Warnings)
	}

	// Groups side by side: each takes a buffer the size of the canvas, and a
	// layer one of its own size, until the document has had its share.
	var ls []testLayer
	for i := range 8 {
		ls = append(ls, testLayer{section: 3}, testLayer{top: i, bottom: i + 1, right: 8, rgb: [3]uint8{0, 0, 255}}, testLayer{section: 1})
	}
	limit(t, &maxCompositeWork, 3*(64+8))
	c = convert(t, layered(8, 8, ls), nil)
	if !hasWarning(c.res.Warnings, "the layers are too many or too large to composite them all; the rest are not drawn") {
		t.Errorf("warnings %v", c.res.Warnings)
	}
	for y, want := range []uint8{255, 255, 255, 0, 0, 0, 0, 0} {
		if px := rgbaAt(c.pages[0], 4, y); px[3] != want {
			t.Errorf("row %d: %v, want alpha %d", y, px, want)
		}
	}

	// A layer larger than the largest layer is left out,
	limit(t, &maxCompositeWork, 1<<32)
	limit(t, &maxLayerPixels, 8*8)
	c = convert(t, layered(8, 8, []testLayer{{left: -1, bottom: 8, right: 8, rgb: [3]uint8{0, 255, 0}}, red}), nil)
	if px := rgbaAt(c.pages[0], 3, 3); px != [4]uint8{255, 0, 0, 255} || rgbaAt(c.pages[0], 0, 0)[3] != 0 || !hasWarning(c.res.Warnings, "too large to draw") {
		t.Errorf("large layer: %v at 3,3, warnings %v", px, c.res.Warnings)
	}
	// and a canvas too large to composite is an error, before any buffer
	// is made.
	limit(t, &maxCompositePixels, 8*8-1)
	data := layered(8, 8, []testLayer{red})
	if _, err := Convert(bytes.NewReader(data), int64(len(data)), nil); err == nil || !strings.Contains(err.Error(), "8 × 8 pixels are too large to composite") {
		t.Errorf("error = %v", err)
	}
}

// TestLiveBuffers: the canvas and the groups that are drawn into buffers of
// their own at the same time hold maxLivePixels pixels, which bounds how
// deep the groups of a document go.
func TestLiveBuffers(t *testing.T) {
	limit(t, &maxLivePixels, 3*8*8) // the canvas and two groups
	c := convert(t, nested(2), nil)
	if px := rgbaAt(c.pages[0], 3, 3); px != [4]uint8{255, 0, 0, 255} || len(c.res.Warnings) != 1 {
		t.Errorf("2 groups: %v at 3,3, warnings %v", px, c.res.Warnings)
	}
	c = convert(t, nested(3), nil)
	if px := rgbaAt(c.pages[0], 3, 3); px[3] != 0 || !hasWarning(c.res.Warnings, "groups more than 2 deep are not drawn") {
		t.Errorf("3 groups: %v at 3,3, warnings %v", px, c.res.Warnings)
	}

	// A group is alive while the group clipped to it is drawn into it:
	// blue over the canvas, inside the red square.
	data := layered(8, 8, []testLayer{{section: 3}, red, {section: 1},
		{section: 3}, {bottom: 8, right: 8, rgb: [3]uint8{0, 0, 255}}, {section: 1, clip: true}})
	c = convert(t, data, nil)
	if px := rgbaAt(c.pages[0], 3, 3); px != [4]uint8{0, 0, 255, 255} || rgbaAt(c.pages[0], 0, 0)[3] != 0 || len(c.res.Warnings) != 1 {
		t.Errorf("clipped group: %v at 3,3, warnings %v", px, c.res.Warnings)
	}
	limit(t, &maxLivePixels, 2*8*8)
	c = convert(t, data, nil)
	if px := rgbaAt(c.pages[0], 3, 3); px != [4]uint8{255, 0, 0, 255} || !hasWarning(c.res.Warnings, "groups more than 1 deep are not drawn") {
		t.Errorf("clipped group without room: %v at 3,3, warnings %v", px, c.res.Warnings)
	}

	// What the limit leaves the canvases of documents.
	limit(t, &maxLivePixels, 1<<28)
	for _, c := range []struct{ w, h, want int }{{4000, 3000, 16}, {8192, 8192, 3}, {8192, 4096, 7}, {16, 16, 16}} {
		if got := groupDepth(int64(c.w) * int64(c.h)); got != c.want {
			t.Errorf("canvas of %d × %d: groups %d deep, want %d", c.w, c.h, got, c.want)
		}
	}
}

// TestArtboardPages: the pages of a document together have the pixels of
// four canvases, which artboards that lie on each other go beyond.
func TestArtboardPages(t *testing.T) {
	var ls []testLayer
	for range 6 {
		ls = append(ls, testLayer{section: 3}, red, testLayer{section: 1, artboard: &[4]int{0, 0, 8, 8}})
	}
	data := layered(8, 8, ls)
	c := convert(t, data, nil)
	want := "the artboards cover more than 4 times the canvas: artboard 5 and those after it are left out"
	if c.res.Pages != 4 || c.res.Artboards != 4 || len(c.pages) != 4 || !slices.Contains(c.res.Warnings, want) {
		t.Errorf("%d pages, %d artboards, warnings %q", c.res.Pages, c.res.Artboards, c.res.Warnings)
	}
	if px := rgbaAt(c.pages[3], 3, 3); px != [4]uint8{255, 0, 0, 255} {
		t.Errorf("page 4 at 3,3: %v", px)
	}
	// A page that is asked for again is the same artboard.
	c = convert(t, data, &Options{Pages: converter.PageList(6, 6, 6, 1, 1, 2)})
	if c.res.Pages != 6 || len(c.res.Warnings) != 1 {
		t.Errorf("pages asked for again: %d pages, warnings %q", c.res.Pages, c.res.Warnings)
	}

	// Artboards beside each other are all pages.
	ls = nil
	for i := range 8 {
		ls = append(ls, testLayer{section: 3}, red, testLayer{section: 1, artboard: &[4]int{0, i, 8, i + 1}})
	}
	c = convert(t, layered(8, 8, ls), nil)
	if c.res.Pages != 8 || c.res.Artboards != 8 || len(c.res.Warnings) != 1 {
		t.Errorf("artboards beside each other: %d pages, warnings %q", c.res.Pages, c.res.Warnings)
	}
}

// allocated returns the bytes f allocates.
func allocated(f func()) uint64 {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	f()
	runtime.ReadMemStats(&after)
	return after.TotalAlloc - before.TotalAlloc
}

func TestCompressedSize(t *testing.T) {
	// A composite of 1024 × 1024 pixels with next to no data: the file is
	// refused for that, before room is made for its samples.
	head := slices.Concat([]byte("8BPS\x00\x01\x00\x00\x00\x00\x00\x00"), be16(3), be32(1024), be32(1024), be16(8), be16(modeRGB),
		be32(0), be32(0), be32(0)) // no colour mode data, resources or layers
	for _, c := range []struct {
		name string
		data []byte
	}{
		{"ZIP", slices.Concat(head, be16(compZip), []byte{0x78, 0x9c, 3, 0})},
		{"RLE", slices.Concat(head, be16(compRLE), make([]byte, 2*1024*3))},
	} {
		var err error
		n := allocated(func() { _, err = Convert(bytes.NewReader(c.data), int64(len(c.data)), nil) })
		if err == nil || !strings.Contains(err.Error(), "too little data for the size of the image") {
			t.Errorf("%s composite: error %v", c.name, err)
		}
		if n > 1<<20 {
			t.Errorf("%s composite: %d bytes allocated for a file of %d", c.name, n, len(c.data))
		}
	}
	// The channels of layers.
	s := planeSpec{w: 1024, h: 1024, depth: 8}
	for _, data := range [][]byte{
		slices.Concat(be16(compZip), []byte{0x78, 0x9c, 3, 0}),
		slices.Concat(be16(compZipPred), []byte{0x78, 0x9c, 3, 0}),
		slices.Concat(be16(compRLE), make([]byte, 2*1024)),
	} {
		var err error
		n := allocated(func() { _, err = s.decodeChannel(data) })
		if err == nil || !strings.Contains(err.Error(), "too little data for the size of the image") || n > 1<<16 {
			t.Errorf("channel with compression %d: error %v, %d bytes allocated", data[1], err, n)
		}
	}

	// What compresses best still decodes: rows of one run in PackBits, and
	// zeros in ZIP.
	s = planeSpec{w: 128, h: 2, depth: 8}
	rows, err := s.decodeChannel(slices.Concat(be16(compRLE), be16(2), be16(2), []byte{0x81, 7, 0x81, 9}))
	if err != nil || len(rows) != 256 || rows[0] != 7 || rows[255] != 9 {
		t.Errorf("PackBits runs: %v, %d samples", err, len(rows))
	}
	var z bytes.Buffer
	zw, _ := zlib.NewWriterLevel(&z, zlib.BestCompression)
	zw.Write(make([]byte, 1024*1024))
	zw.Close()
	s = planeSpec{w: 1024, h: 1024, depth: 8}
	if rows, err = s.decodeChannel(slices.Concat(be16(compZip), z.Bytes())); err != nil || len(rows) != 1024*1024 {
		t.Errorf("ZIP of zeros (%d bytes): %v, %d samples", z.Len(), err, len(rows))
	}
}

// TestDescriptorList reads a list that states a million items and holds
// none that can be read.
func TestDescriptorList(t *testing.T) {
	data := slices.Concat([]byte("VlLs"), be32(1<<20), make([]byte, 1<<20))
	var v any
	var ok bool
	n := allocated(func() {
		v, ok = readValue(&cursor{r: bytes.NewReader(data), end: int64(len(data))}, 0)
	})
	if ok || v != nil || n > 1<<16 {
		t.Errorf("list read: %v, %v, %d bytes allocated", v, ok, n)
	}
	// A list longer than the room made for it at first.
	data = slices.Concat([]byte("VlLs"), be32(3000))
	for i := range 3000 {
		data = append(append(data, "long"...), be32(i)...)
	}
	v, ok = readValue(&cursor{r: bytes.NewReader(data), end: int64(len(data))}, 0)
	if list, _ := v.([]any); !ok || len(list) != 3000 || list[2999] != int64(2999) {
		t.Errorf("list of 3000: %v items, %v", len(list), ok)
	}
}
