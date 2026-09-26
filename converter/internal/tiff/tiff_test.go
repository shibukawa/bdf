package tiff

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf/imgconv"
)

// The files in testdata come from test/tiff/gen.sh: each holds one picture
// stored in several ways, one per page, and ref-*.png is that picture as
// ImageMagick decodes it through libtiff.

func open(t *testing.T, name string) *File {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	f, err := Open(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return f
}

func reference(t *testing.T, name string) *image.NRGBA {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "ref-"+name+".png"))
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return nrgba(img)
}

func nrgba(img image.Image) *image.NRGBA {
	b := img.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), img, b.Min, draw.Src)
	return out
}

// diff returns the largest and the mean difference of the channels.
func diff(a, b *image.NRGBA) (maxd int, mean float64) {
	if a.Rect != b.Rect {
		return 256, 256
	}
	sum := 0
	for i := range a.Pix {
		d := int(a.Pix[i]) - int(b.Pix[i])
		if d < 0 {
			d = -d
		}
		sum += d
		maxd = max(maxd, d)
	}
	return maxd, float64(sum) / float64(len(a.Pix))
}

func TestDecode(t *testing.T) {
	for _, c := range []struct {
		file, ref string
		pages     int
		// the largest channel difference allowed: ImageMagick's CMYK
		// conversion and palette rounding differ from Go's by one.
		tolerance int
	}{
		{"rgb.tif", "rgb", 10, 0},        // none, LZW, LZW+predictor, Deflate, Deflate+predictor, PackBits, planar, tiles, 7-row strips, 16-bit
		{"rgb-be.tif", "rgb", 2, 0},      // big-endian: 16-bit with predictor, 8-bit
		{"rgb-bigtiff.tif", "rgb", 2, 0}, // BigTIFF: LZW, planar Deflate
		{"gray.tif", "gray", 5, 0},       // none, LZW+predictor, Deflate, PackBits, 16-bit Deflate+predictor
		{"bw.tif", "bw", 15, 0},          // BlackIsZero: none, PackBits, LZW, Deflate, G3, G3 fill, G3 2-D, G3 2-D fill, G4, G4 LSB first, G4 16-row strips, G4 tiles; WhiteIsZero: none, G3 2-D fill, G4
		{"pal.tif", "pal", 2, 1},
		{"rgba.tif", "rgba", 1, 0},
		{"cmyk.tif", "cmyk", 3, 1}, // none, Deflate (and JPEG, which TestJPEG checks)
	} {
		f := open(t, c.file)
		want := reference(t, c.ref)
		pages := f.Pages()
		if len(pages) != c.pages {
			t.Errorf("%s: %d pages, want %d", c.file, len(pages), c.pages)
		}
		for i, d := range pages {
			if d.Compression() == CompressionJPEG {
				continue
			}
			img, damaged, err := d.Decode()
			if err != nil || damaged {
				t.Errorf("%s page %d: %v (damaged %v)", c.file, i+1, err, damaged)
				continue
			}
			if maxd, _ := diff(nrgba(img), want); maxd > c.tolerance {
				t.Errorf("%s page %d (%T): differs from the reference by up to %d", c.file, i+1, img, maxd)
			}
		}
	}
}

// TestDecodeTypes checks the image types pages decode into, which decide
// how they are stored.
func TestDecodeTypes(t *testing.T) {
	for _, c := range []struct {
		file string
		want image.Image
	}{
		{"bw.tif", &image.Paletted{}},
		{"gray.tif", &image.Gray{}},
		{"pal.tif", &image.Paletted{}},
		{"rgb.tif", &image.RGBA{}},
		{"rgba.tif", &image.NRGBA{}},
		{"cmyk.tif", &image.CMYK{}},
		{"jpeg.tif", &image.YCbCr{}},
	} {
		img, _, err := open(t, c.file).Pages()[0].Decode()
		if err != nil {
			t.Fatal(err)
		}
		if got, want := typeName(img), typeName(c.want); got != want {
			t.Errorf("%s decodes into %s, want %s", c.file, got, want)
		}
		if p, ok := img.(*image.Paletted); ok && c.file == "bw.tif" && len(p.Palette) != 2 {
			t.Errorf("bilevel page with %d colours", len(p.Palette))
		}
	}
}

func typeName(img image.Image) string {
	switch img.(type) {
	case *image.Paletted:
		return "Paletted"
	case *image.Gray:
		return "Gray"
	case *image.RGBA:
		return "RGBA"
	case *image.NRGBA:
		return "NRGBA"
	case *image.CMYK:
		return "CMYK"
	case *image.YCbCr:
		return "YCbCr"
	}
	return "?"
}

func TestJPEG(t *testing.T) {
	rgb, gray := reference(t, "rgb"), reference(t, "gray")
	for _, c := range []struct {
		file   string
		page   int
		joined bool
		ref    *image.NRGBA
	}{
		{"jpeg.tif", 0, true, rgb},        // YCbCr, 16-row strips
		{"jpeg.tif", 1, true, rgb},        // YCbCr, 32-row strips
		{"jpeg.tif", 2, true, rgb},        // RGB
		{"jpeg.tif", 3, true, gray},       // grey
		{"jpeg.tif", 4, false, rgb},       // tiles
		{"jpeg-magick.tif", 0, true, rgb}, // ImageMagick's RGB strips
		{"cmyk.tif", 2, false, reference(t, "cmyk")},
	} {
		d := open(t, c.file).Pages()[c.page]
		img, damaged, err := d.Decode()
		if err != nil || damaged {
			t.Fatalf("%s page %d: %v (damaged %v)", c.file, c.page+1, err, damaged)
		}
		// JPEG is lossy (and the chroma of these small pictures is
		// subsampled), so this only catches gross errors such as reading
		// RGB as YCbCr.
		if _, mean := diff(nrgba(img), c.ref); mean > 8 {
			t.Errorf("%s page %d: differs from the reference by %.2f on average", c.file, c.page+1, mean)
		}
		b, ok := d.JPEG()
		if ok != c.joined {
			t.Errorf("%s page %d: joined %v, want %v", c.file, c.page+1, ok, c.joined)
		}
		if !ok {
			continue
		}
		// The joined stream holds the strips' data as they are: it decodes
		// to what decoding the strips one by one gives.
		joined, err := jpeg.Decode(bytes.NewReader(b))
		if err != nil {
			t.Fatalf("%s page %d: joined stream: %v", c.file, c.page+1, err)
		}
		l, _ := d.layout()
		tables, _ := l.tables()
		for i := range l.down {
			raw, _ := l.rawBlock(i)
			sb, _ := l.standalone(tables, raw)
			strip, err := jpeg.Decode(bytes.NewReader(sb))
			if err != nil {
				t.Fatal(err)
			}
			want := nrgba(strip)
			got := nrgba(joined).SubImage(image.Rect(0, i*l.bh, l.w, min(l.h, (i+1)*l.bh))).(*image.NRGBA)
			for y := range want.Rect.Dy() - max(0, (i+1)*l.bh-l.h) {
				for x := range l.w {
					if want.NRGBAAt(x, y) != got.NRGBAAt(x, i*l.bh+y) {
						t.Fatalf("%s page %d: joined stream differs from strip %d at %d,%d", c.file, c.page+1, i, x, y)
					}
				}
			}
		}
	}
}

func TestPages(t *testing.T) {
	f := open(t, "thumb.tif")
	if len(f.IFDs) != 3 || len(f.Pages()) != 2 {
		t.Fatalf("%d IFDs, %d pages; want 3 IFDs and 2 pages (the reduced-resolution copy left out)", len(f.IFDs), len(f.Pages()))
	}
	if w, _ := f.Pages()[1].Size(); w != 67 {
		t.Errorf("second page is %d pixels wide, want the grey page's 67", w)
	}
}

// tiny builds a classic little-endian TIFF of 1×1 grey pixels with the
// given IFD offsets chained in order (next[i] is IFD i's next offset).
func tiny(next func(i int, offs []uint32) uint32, n int) []byte {
	var b bytes.Buffer
	b.WriteString("II*\x00")
	binary.Write(&b, binary.LittleEndian, uint32(8))
	offs := make([]uint32, n)
	for i := range n {
		offs[i] = uint32(8 + i*(2+5*12+4))
	}
	for i := range n {
		binary.Write(&b, binary.LittleEndian, uint16(5))
		for _, e := range [][3]uint32{{256, 1, 1}, {257, 1, 1}, {258, 1, 8}, {262, 1, 1}, {273, 1, 0}} {
			binary.Write(&b, binary.LittleEndian, uint16(e[0]))
			binary.Write(&b, binary.LittleEndian, uint16(3))
			binary.Write(&b, binary.LittleEndian, e[1])
			binary.Write(&b, binary.LittleEndian, e[2])
		}
		binary.Write(&b, binary.LittleEndian, next(i, offs))
	}
	return b.Bytes()
}

func TestChain(t *testing.T) {
	loop := tiny(func(i int, offs []uint32) uint32 { return offs[(i+1)%len(offs)] }, 3)
	f, err := Open(bytes.NewReader(loop), int64(len(loop)))
	if err != nil {
		t.Fatal(err)
	}
	if len(f.IFDs) != 3 || f.ChainErr == nil {
		t.Errorf("looping chain: %d IFDs, error %v; want 3 and an error", len(f.IFDs), f.ChainErr)
	}
	bad := tiny(func(i int, offs []uint32) uint32 {
		if i == 1 {
			return 1 << 30
		}
		return offs[i+1]
	}, 2)
	f, err = Open(bytes.NewReader(bad), int64(len(bad)))
	if err != nil || len(f.IFDs) != 2 || f.ChainErr == nil {
		t.Errorf("chain leaving the file: %v, %d IFDs, %v", err, len(f.IFDs), f.ChainErr)
	}
	if _, err := Open(bytes.NewReader([]byte("II*\x00\xff\xff\x00\x00")), 8); err == nil {
		t.Error("no error for a first IFD outside the file")
	}
}

func TestSniff(t *testing.T) {
	for head, want := range map[string]bool{
		"II*\x00\x08\x00\x00\x00": true, "MM\x00*\x00\x00\x00\x08": true, "II+\x00\x08\x00\x00\x00": true,
		"MM\x00+": true, "II*": false, "%PDF-1.7": false, "\x89PNG": false,
	} {
		if got := Sniff([]byte(head)); got != want {
			t.Errorf("Sniff(%q) = %v", head, got)
		}
	}
}

func TestTooLarge(t *testing.T) {
	old := MaxDecodeBytes
	defer func() { MaxDecodeBytes = old }()
	MaxDecodeBytes = 1000
	if _, _, err := open(t, "rgb.tif").Pages()[0].Decode(); err == nil {
		t.Error("no error for a page over MaxDecodeBytes")
	}
}

func TestCCITT(t *testing.T) {
	// Modified Huffman (compression 2): each row is byte-aligned. A row of
	// 3 white and 5 black pixels is 1000 (white 3) 0011 (black 5).
	out, damaged := decodeCCITT([]byte{0x83, 0x83}, ccittRLE, false, false, 8, 2)
	if damaged || !bytes.Equal(out, []byte{0x1f, 0x1f}) {
		t.Errorf("RLE: % x (damaged %v), want 1f 1f", out, damaged)
	}
	// A damaged row of T.4 data is left white, and decoding carries on at
	// the next EOL.
	f := open(t, "bw.tif")
	d := f.Pages()[5] // G3 with fill bits
	l, _ := d.layout()
	src, _ := l.rawBlock(0)
	want, _ := decodeCCITT(src, ccittT4, false, false, l.bw, l.h)
	src = append([]byte(nil), src...)
	mid := len(src) / 2
	for i := mid; i < mid+4; i++ {
		src[i] = 0xff
	}
	got, damaged := decodeCCITT(src, ccittT4, false, false, l.bw, l.h)
	if !damaged {
		t.Error("damaged data not reported")
	}
	rb := (l.bw + 7) / 8
	same := 0
	for y := range l.h {
		if bytes.Equal(got[y*rb:(y+1)*rb], want[y*rb:(y+1)*rb]) {
			same++
		}
	}
	if same < l.h-8 {
		t.Errorf("only %d of %d rows decoded after damage", same, l.h)
	}
}

func TestMetadata(t *testing.T) {
	b, err := os.ReadFile("../../tiff/testdata/scan.tif")
	if err != nil {
		t.Fatal(err)
	}
	f, err := Open(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		t.Fatal(err)
	}
	d := f.Pages()[0]
	if s := d.Strings(TagArtist); len(s) != 1 || s[0] != "Alice; Bob" {
		t.Errorf("Artist %q", s)
	}
	if x, y, unit := d.Resolution(); x != 300 || y != 300 || unit != ResUnitInch {
		t.Errorf("resolution %v %v %v", x, y, unit)
	}
	if x, y, _ := f.Pages()[2].Resolution(); x != 0 || y != 0 {
		t.Errorf("page 3 has resolution %v %v, want none", x, y)
	}
}

func TestPicture(t *testing.T) {
	read := func(name string) []byte {
		b, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	// Lossless pages become PNG in any mode (the caller converts), with the
	// pixels libtiff reads: the BlackIsZero Group 4 page is not inverted.
	for _, c := range []struct{ file, ref string }{{"g4.tif", "bw"}, {"rgb.tif", "rgb"}, {"pal.tif", "pal"}} {
		for _, mode := range []imgconv.Mode{imgconv.Keep, imgconv.Convert} {
			pic, damaged, err := Picture(read(c.file), imgconv.Options{Mode: mode})
			if err != nil || damaged || imgconv.Sniff(pic) != "png" {
				t.Fatalf("%s: %s, damaged %v, %v", c.file, imgconv.Sniff(pic), damaged, err)
			}
			img, err := png.Decode(bytes.NewReader(pic))
			if err != nil {
				t.Fatal(err)
			}
			if maxd, _ := diff(nrgba(img), reference(t, c.ref)); maxd > 1 {
				t.Errorf("%s: differs from the reference by up to %d", c.file, maxd)
			}
		}
	}
	// JPEG strips that join are the joined stream.
	pic, _, err := Picture(read("jpeg.tif"), imgconv.Options{})
	d, _ := FirstPage(read("jpeg.tif"))
	joined, _ := d.JPEG()
	if err != nil || !bytes.Equal(pic, joined) {
		t.Errorf("JPEG page: %s, %v", imgconv.Sniff(pic), err)
	}
	if _, _, err := Picture([]byte("II*\x00\x00\x00\x00\x00"), imgconv.Options{}); err == nil {
		t.Error("no error for a TIFF file without IFDs")
	}
}
