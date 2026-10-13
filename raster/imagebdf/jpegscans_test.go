package imagebdf

import (
	"bytes"
	"image/color"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/image/imgconv"
)

// TestJPEGOfManyScans draws a JPEG image that is nothing but scans, each a
// few bytes that send the decoder through every block of the image again:
// it is read as far as imgconv.MaxJPEGScans of them, with a warning.
func TestJPEGOfManyScans(t *testing.T) {
	jpegOf := func(scans int) []byte {
		// grey, 1024 × 1024, progressive, every coefficient zero: see
		// scansJPEG of the tests of imgconv
		var b bytes.Buffer
		b.Write([]byte{0xff, 0xd8, 0xff, 0xdb, 0, 67, 0})
		b.Write(bytes.Repeat([]byte{1}, 64))
		b.Write([]byte{0xff, 0xc2, 0, 11, 8, 4, 0, 4, 0, 1, 1, 0x11, 0, 0xff, 0xc4, 0, 20, 0x10, 1})
		b.Write(make([]byte, 15))
		b.WriteByte(0xe0)
		for i := 0; i < scans; i++ {
			b.Write([]byte{0xff, 0xda, 0, 8, 1, 1, 0, 1, 63, 0, 0, 0})
		}
		b.Write([]byte{0xff, 0xd9})
		return b.Bytes()
	}
	for _, scans := range []int{imgconv.MaxJPEGScans, 20000} {
		d := bdf.NewDocument()
		o := bdf.NewObject()
		o.Image(o.AddImage(d.AddImage(jpegOf(scans))), 0, 0, 8, 8)
		h, _ := d.AddObject(o)
		v := d.NewView("v", bdf.ViewFixed, "")
		v.AddPage(8, 8, bdf.Layer{Role: bdf.RoleBody, Obj: h})
		r := New(d, &Options{NoSystemFonts: true})
		start := time.Now()
		img, err := r.Page(v, 0, 1)
		if err != nil {
			t.Fatal(err)
		}
		if got := img.RGBAAt(4, 4); got != (color.RGBA{128, 128, 128, 255}) || warned(r, "scans") != (scans > imgconv.MaxJPEGScans) || len(r.Warnings()) > 1 {
			t.Errorf("%d scans: drawn %v, warnings %v", scans, got, r.Warnings())
		}
		if took := time.Since(start); took > 3*time.Second {
			t.Errorf("%d scans took %v", scans, took)
		}
	}
}
