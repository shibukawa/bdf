package pptx

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/tiff"
	"github.com/shibukawa/bdf/imgconv"
)

// TestTIFFPictures draws TIFF pictures, read by converter/internal/tiff: a
// BlackIsZero Group 4 page (which golang.org/x/image/tiff shows inverted),
// YCbCr JPEG strips (which it does not read), and the Group 4 page again
// with its white made transparent.
func TestTIFFPictures(t *testing.T) {
	read := func(name string) []byte {
		b, err := os.ReadFile("../internal/tiff/testdata/" + name)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	g4, jpg := read("g4.tif"), read("jpeg.tif")
	ref, err := png.Decode(bytes.NewReader(read("ref-bw.png")))
	if err != nil {
		t.Fatal(err)
	}
	transparentWhite := `<a:clrChange><a:clrFrom><a:srgbClr val="FFFFFF"/></a:clrFrom><a:clrTo><a:srgbClr val="FFFFFF"><a:alpha val="0"/></a:srgbClr></a:clrTo></a:clrChange>`
	deck := pictureDeck(t,
		pic("2", "g4", 457200, "")+pic("3", "jpg", 2743200, "")+pic("4", "g4", 5029200, transparentWhite),
		map[string][]byte{"g4.tif": g4, "jpg.tif": jpg})
	res, err := Convert(bytes.NewReader(deck), int64(len(deck)), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	o, err := r.Object(body(r.Manifest.Views[0].Pages[0]))
	if err != nil {
		t.Fatal(err)
	}
	var drawn [][]byte
	o.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpImage || in.Op == bdf.OpImageSub {
			b, err := r.Part(o.Images[in.Args[0].(uint64)])
			if err != nil {
				t.Fatal(err)
			}
			drawn = append(drawn, b)
		}
	})
	if len(drawn) != 3 {
		t.Fatalf("%d pictures drawn, want 3", len(drawn))
	}
	decode := func(b []byte) image.Image {
		img, err := imgconv.Decode(b)
		if err != nil {
			t.Fatal(err)
		}
		return img
	}
	// The Group 4 page: a PNG with the pixels libtiff reads.
	if imgconv.Sniff(drawn[0]) != "png" {
		t.Errorf("Group 4 picture stored as %s", imgconv.Sniff(drawn[0]))
	}
	plain := decode(drawn[0])
	black := image.Point{-1, -1}
	for y := range ref.Bounds().Dy() {
		for x := range ref.Bounds().Dx() {
			rr, _, _, _ := ref.At(x, y).RGBA()
			pr, _, _, _ := plain.At(x, y).RGBA()
			if rr>>8 != pr>>8 {
				t.Fatalf("Group 4 picture differs from libtiff's reading at %d,%d", x, y)
			}
			if rr == 0 && black.X < 0 {
				black = image.Point{x, y}
			}
		}
	}
	// The JPEG strips: stored as the JPEG they join into.
	d, err := tiff.FirstPage(jpg)
	if err != nil {
		t.Fatal(err)
	}
	if joined, ok := d.JPEG(); !ok || !bytes.Equal(drawn[1], joined) {
		t.Errorf("JPEG picture is not the joined strips (%s, %d bytes)", imgconv.Sniff(drawn[1]), len(drawn[1]))
	}
	// Recoloured: white turned transparent, black kept.
	rc := decode(drawn[2])
	if _, _, _, a := rc.At(0, 0).RGBA(); a != 0 {
		t.Errorf("white corner has alpha %d after clrChange", a>>8)
	}
	if cr, _, _, a := rc.At(black.X, black.Y).RGBA(); a>>8 != 255 || cr != 0 {
		t.Errorf("black pixel %v is %d/%d after clrChange", black, cr>>8, a>>8)
	}
}
