package imgconv

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// scansJPEG returns a grey progressive JPEG image of w × h pixels whose
// coefficients are all zero, sent in n scans, and where the marker of each
// scan is. A scan of it says "the blocks from here on have nothing more"
// as often as it takes to pass every block, in two bytes for every 16384
// of them: what the decoder does for it is to visit the blocks.
func scansJPEG(w, h, n int) (data []byte, scans []int) {
	var b bytes.Buffer
	b.Write([]byte{0xff, 0xd8})
	b.Write([]byte{0xff, 0xdb, 0, 67, 0}) // a quantization table of ones
	b.Write(bytes.Repeat([]byte{1}, 64))
	b.Write([]byte{0xff, 0xc2, 0, 11, 8, byte(h >> 8), byte(h), byte(w >> 8), byte(w), 1, 1, 0x11, 0})
	// a Huffman table for AC coefficients with one code, of one bit: the
	// end of the bands of 16384 blocks, and of as many more as the 14
	// bits after it say
	b.Write([]byte{0xff, 0xc4, 0, 20, 0x10, 1})
	b.Write(make([]byte, 15))
	b.WriteByte(0xe0)
	runs := (((w+7)/8)*((h+7)/8) + 16383) / 16384
	bits := make([]byte, (15*runs+7)/8)
	for i := 0; i < n; i++ {
		scans = append(scans, b.Len())
		b.Write([]byte{0xff, 0xda, 0, 8, 1, 1, 0, 1, 63, 0}) // the AC coefficients of the one component
		b.Write(bits)
	}
	b.Write([]byte{0xff, 0xd9})
	return b.Bytes(), scans
}

func TestJPEGScans(t *testing.T) {
	grey := func(data []byte) bool {
		img, err := Decode(data)
		if err != nil {
			t.Error(err)
			return false
		}
		g, ok := img.(*image.Gray)
		return ok && g.Bounds() == image.Rect(0, 0, 64, 48) && g.GrayAt(0, 0).Y == 128 && g.GrayAt(63, 47).Y == 128
	}
	// an image of no more scans than are read is read as it is
	for _, n := range []int{1, 4, MaxJPEGScans} {
		data, _ := scansJPEG(64, 48, n)
		if got, cut := LimitJPEGScans(data); cut || &got[0] != &data[0] || !grey(data) {
			t.Errorf("%d scans: cut %v, decoded %v", n, cut, grey(data))
		}
	}
	// one of more is the image an encoder had written with the first of them
	want, _ := scansJPEG(64, 48, MaxJPEGScans)
	for _, n := range []int{MaxJPEGScans + 1, 5000} {
		data, at := scansJPEG(64, 48, n)
		got, cut := LimitJPEGScans(data)
		if !cut || !bytes.Equal(got, want) || jpegScansEnd(data, MaxJPEGScans) != at[MaxJPEGScans] || !grey(data) {
			t.Errorf("%d scans: cut %v to %d bytes (want %d), the next scan at %d (want %d)", n, cut, len(got), len(want), jpegScansEnd(data, MaxJPEGScans), at[MaxJPEGScans])
		}
		for _, k := range []int{0, 1, 7, n - 1} {
			if end := jpegScansEnd(data, k); end != at[k] {
				t.Errorf("%d scans: the one after %d at %d, want %d", n, k, end, at[k])
			}
		}
		if end := jpegScansEnd(data, n); end != -1 {
			t.Errorf("%d scans: one after them at %d", n, end)
		}
	}
	// what is no JPEG image is not looked into
	other := []byte("\x89PNG\r\n\x1a\n\xff\xda\x00\x08")
	if got, cut := LimitJPEGScans(other); cut || &got[0] != &other[0] {
		t.Error("a PNG file is cut")
	}
}

// TestJPEGScansMarkers checks that scans are counted where the decoder
// finds them: not in a segment, and after the data of a scan with what
// such data has of 0xff bytes.
func TestJPEGScansMarkers(t *testing.T) {
	scan := []byte{0xff, 0xda, 0, 8, 1, 1, 0, 1, 63, 0}
	var b bytes.Buffer
	b.Write([]byte{0xff, 0xd8})
	// a segment with a scan marker among its bytes, as that of a thumbnail has
	b.Write([]byte{0xff, 0xe1, 0, 2 + 12})
	b.Write(append([]byte{'E', 'x'}, scan...))
	var at []int
	for i, entropy := range [][]byte{
		{0x12, 0xff, 0x00, 0x34, 0xff, 0xd0, 0x56}, // a 0xff of the data, and a restart marker
		{0xff, 0xff, 0xd7, 0x78, 0xff, 0x00},       // fill bytes before a restart marker
		{},                                         // no data at all
		{0x9a},
	} {
		if i == 2 {
			at = append(at, b.Len())
			b.Write([]byte{0xff, 0xff, 0xff}) // fill bytes before the marker of a scan: cut before them
		} else {
			at = append(at, b.Len())
		}
		b.Write(scan)
		b.Write(entropy)
	}
	end := b.Len()
	b.Write([]byte{0xff, 0xd9})
	b.Write(scan) // after the end of the image: not of it
	data := b.Bytes()
	for k, want := range at {
		if got := jpegScansEnd(data, k); got != want {
			t.Errorf("the scan after %d at %d, want %d", k, got, want)
		}
	}
	if got := jpegScansEnd(data, len(at)); got != -1 {
		t.Errorf("a scan after the end of the image at %d (the end is at %d)", got, end)
	}
	// cut short anywhere, there is a scan to stop at or there is none
	for n := 0; n <= len(data); n++ {
		for k := 0; k <= len(at); k++ {
			if got := jpegScansEnd(data[:n:n], k); got != -1 && (k >= len(at) || got != at[k] || got+4 > n) {
				t.Fatalf("%d bytes: the scan after %d at %d", n, k, got)
			}
		}
	}
	// a segment shorter than its own length: the decoder gives up there
	short := append([]byte{0xff, 0xd8, 0xff, 0xe0, 0, 1}, scan...)
	if got := jpegScansEnd(short, 0); got != -1 {
		t.Errorf("a scan after a segment of no length at %d", got)
	}
}

// TestJPEGScansTime decodes an image of a megapixel in 20,000 scans, which
// took seven seconds: now what a hundred scans take.
func TestJPEGScansTime(t *testing.T) {
	data, _ := scansJPEG(1024, 1024, 20000)
	start := time.Now()
	img, err := Decode(data)
	if took := time.Since(start); err != nil || img.Bounds().Dx() != 1024 || took > 3*time.Second {
		t.Errorf("%d bytes: %v after %v", len(data), err, took)
	}
}

// TestJPEGScansOfEncoders cuts progressive images as encoders write them
// (testdata: made-up pictures, one grey with restart markers, one in
// colour with its chroma at half the size) after each number of scans:
// every cut is an image, and the cut after the last is the whole.
func TestJPEGScansOfEncoders(t *testing.T) {
	for name, scans := range map[string]int{"progressive-grey-restarts.jpg": 6, "progressive-420.jpg": 10} {
		data, err := os.ReadFile(filepath.Join("testdata", name))
		if err != nil {
			t.Fatal(err)
		}
		whole, err := jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		if end := jpegScansEnd(data, scans); end != -1 {
			t.Errorf("%s: a scan after its %d at %d", name, scans, end)
		}
		for k := 1; k < scans; k++ {
			end := jpegScansEnd(data, k)
			if end < 0 {
				t.Fatalf("%s: no scan after the first %d", name, k)
			}
			cut, err := jpeg.Decode(bytes.NewReader(append(data[:end:end], 0xff, 0xd9)))
			if err != nil || cut.Bounds() != whole.Bounds() {
				t.Errorf("%s cut after %d scans: %v", name, k, err)
			}
		}
		// and Decode, which cuts nothing of it, reads what the decoder reads
		got, err := Decode(data)
		if err != nil || !reflect.DeepEqual(got, whole) {
			t.Errorf("%s: %v", name, err)
		}
	}
}
