package visio

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"os"
	"regexp"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/tiff"
	"github.com/shibukawa/bdf/imgconv"
)

// TestTIFFForeignData replaces the PNG picture of flow.vdx with TIFF
// pictures, read by converter/internal/tiff: a BlackIsZero Group 4 page
// (which golang.org/x/image/tiff shows inverted) and YCbCr JPEG strips
// (which it does not read).
func TestTIFFForeignData(t *testing.T) {
	vdx, err := os.ReadFile("testdata/flow.vdx")
	if err != nil {
		t.Fatal(err)
	}
	read := func(name string) []byte {
		b, err := os.ReadFile("../internal/tiff/testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	foreign := regexp.MustCompile(`<ForeignData ForeignType='Bitmap' CompressionType='PNG'>[^<]*</ForeignData>`)
	if n := len(foreign.FindAll(vdx, -1)); n != 1 {
		t.Fatalf("flow.vdx has %d PNG pictures, want 1", n)
	}
	// images converts flow.vdx with the picture replaced and returns the
	// image parts stored.
	images := func(pic []byte) [][]byte {
		data := foreign.ReplaceAll(vdx, []byte(`<ForeignData ForeignType='Bitmap' CompressionType='TIFF'>`+base64.StdEncoding.EncodeToString(pic)+`</ForeignData>`))
		res, err := Convert(bytes.NewReader(data), int64(len(data)), testOptions())
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Warnings) != 0 {
			t.Errorf("warnings: %v", res.Warnings)
		}
		var out [][]byte
		for _, p := range res.Doc.Parts() {
			if p.Type == bdf.PartImage {
				out = append(out, p.Data)
			}
		}
		return out
	}

	ref, err := png.Decode(bytes.NewReader(read("ref-bw.png")))
	if err != nil {
		t.Fatal(err)
	}
	var g4 image.Image
	for _, b := range images(read("g4.tif")) {
		if img, err := imgconv.Decode(b); err == nil && img.Bounds().Size() == ref.Bounds().Size() {
			g4 = img
		}
	}
	if g4 == nil {
		t.Fatal("the Group 4 picture is not stored")
	}
	for y := range ref.Bounds().Dy() {
		for x := range ref.Bounds().Dx() {
			rr, _, _, _ := ref.At(x, y).RGBA()
			gr, _, _, _ := g4.At(x, y).RGBA()
			if rr>>8 != gr>>8 {
				t.Fatalf("Group 4 picture differs from libtiff's reading at %d,%d", x, y)
			}
		}
	}

	jpg := read("jpeg.tif")
	d, err := tiff.FirstPage(jpg)
	if err != nil {
		t.Fatal(err)
	}
	joined, ok := d.JPEG()
	if !ok {
		t.Fatal("jpeg.tif's first page does not join")
	}
	found := false
	for _, b := range images(jpg) {
		found = found || bytes.Equal(b, joined)
	}
	if !found {
		t.Error("the JPEG picture is not stored as the joined strips")
	}
}
