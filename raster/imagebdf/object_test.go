package imagebdf

import (
	"errors"
	"image"
	"image/color"
	"testing"

	"github.com/shibukawa/bdf"
)

// TestObject draws an object of its own: the image is the bounding box in
// whole pixels, placed from the origin of the object, on no background.
func TestObject(t *testing.T) {
	o := bdf.NewObject()
	o.FillColor(bdf.RGBA(255, 0, 0, 255))
	o.FillRect(-2, -3, 2, 3) // up and left of the origin
	o.FillColor(bdf.RGBA(0, 0, 255, 255))
	o.FillRect(0.25, 0.25, 1.75, 0.75) // starts half a pixel in at scale 2
	img, err := Object(o, 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	if want := image.Rect(-4, -6, 4, 2); img.Bounds() != want {
		t.Fatalf("bounds %v, want %v", img.Bounds(), want)
	}
	for _, c := range []struct {
		x, y int
		want color.RGBA
	}{
		{-4, -6, color.RGBA{255, 0, 0, 255}},
		{-1, -1, color.RGBA{255, 0, 0, 255}},
		{0, -1, color.RGBA{}}, // nothing drawn: transparent
		{-1, 0, color.RGBA{}},
		{0, 0, color.RGBA{0, 0, 64, 64}}, // a quarter of the pixel
		{1, 0, color.RGBA{0, 0, 128, 128}},
		{1, 1, color.RGBA{0, 0, 255, 255}},
		{3, 1, color.RGBA{0, 0, 255, 255}},
	} {
		if got := img.RGBAAt(c.x, c.y); got != c.want {
			t.Errorf("pixel (%d, %d): %v, want %v", c.x, c.y, got, c.want)
		}
	}

	img, err = Object(o, 2, &Options{Background: color.White})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := img.RGBAAt(0, -1), (color.RGBA{255, 255, 255, 255}); got != want {
		t.Errorf("background: %v, want %v", got, want)
	}
}

func TestObjectErrors(t *testing.T) {
	if _, err := Object(nil, 1, nil); err == nil {
		t.Error("no object: no error")
	}
	if _, err := Object(bdf.NewObject(), 1, nil); err == nil {
		t.Error("empty object: no error")
	}
	o := bdf.NewObject()
	o.FillRect(0, 0, 10, 10)
	if _, err := Object(o, 0, nil); err == nil {
		t.Error("scale 0: no error")
	}
	if _, err := Object(o, 1e4, nil); !errors.Is(err, ErrTooLarge) {
		t.Errorf("10⁵ × 10⁵ pixels: %v, want ErrTooLarge", err)
	}
	dep := bdf.NewObject()
	dep.Image(dep.AddImage(bdf.Hash{1}), 0, 0, 10, 10)
	if _, err := Object(dep, 1, nil); !errors.Is(err, ErrNeedsDocument) {
		t.Errorf("an object with an image: %v, want ErrNeedsDocument", err)
	}
}
