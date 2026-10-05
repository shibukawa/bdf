package imagebdf

import (
	"errors"
	"image"
	"image/color"
	"math"
	"testing"

	"github.com/shibukawa/bdf"
)

func TestPageDPIAndSize(t *testing.T) {
	d := bdf.NewDocument()
	o := bdf.NewObject()
	o.FillColor(0xff0000ff).FillRect(0, 0, 36, 36)
	hash, _ := d.AddObject(o)
	v := d.NewView("page", bdf.ViewFixed, "")
	v.AddPage(72, 36, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	r := New(d, nil)

	for _, tc := range []struct {
		name            string
		draw            func() (*image.RGBA, error)
		bounds          image.Rectangle
		redX, whiteX, y int
	}{
		{"dpi", func() (*image.RGBA, error) { return r.PageDPI(v, 0, 144) }, image.Rect(0, 0, 144, 72), 30, 120, 20},
		{"size", func() (*image.RGBA, error) { return r.PageSize(v, 0, 30, 45) }, image.Rect(0, 0, 30, 45), 5, 25, 20},
	} {
		t.Run(tc.name, func(t *testing.T) {
			img, err := tc.draw()
			if err != nil {
				t.Fatal(err)
			}
			var _ image.Image = img
			if img.Bounds() != tc.bounds {
				t.Fatalf("bounds = %v, want %v", img.Bounds(), tc.bounds)
			}
			if got := img.RGBAAt(tc.redX, tc.y); got != (color.RGBA{R: 255, A: 255}) {
				t.Errorf("red pixel = %v", got)
			}
			if got := img.RGBAAt(tc.whiteX, tc.y); got != (color.RGBA{R: 255, G: 255, B: 255, A: 255}) {
				t.Errorf("white pixel = %v", got)
			}
		})
	}
}

func TestPageSizeScrollStrip(t *testing.T) {
	d := bdf.NewDocument()
	o := bdf.NewObject()
	o.FillColor(0x0000ffff).FillRect(0, 0, 10, 10)
	hash, _ := d.AddObject(o)
	v := d.NewView("scroll", bdf.ViewScroll, "")
	v.AddPage(10, 10)
	v.AddPage(10, 10, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	img, err := New(d, nil).PageSize(v, 1, 20, 20)
	if err != nil {
		t.Fatal(err)
	}
	if got := img.RGBAAt(10, 10); got != (color.RGBA{B: 255, A: 255}) {
		t.Errorf("second strip pixel = %v", got)
	}
}

func TestPageDPIAndSizeErrors(t *testing.T) {
	d := bdf.NewDocument()
	v := d.NewView("page", bdf.ViewFixed, "")
	v.AddPage(72, 36)
	r := New(d, nil)
	for _, dpi := range []float64{0, -1, math.NaN(), math.Inf(1)} {
		if img, err := r.PageDPI(v, 0, dpi); err == nil || img != nil {
			t.Errorf("dpi %v: image %v, error %v", dpi, img, err)
		}
	}
	for _, size := range [][2]int{{0, 10}, {10, -1}, {MaxPixels + 1, 1}, {int(^uint(0) >> 1), 3}} {
		if img, err := r.PageSize(v, 0, size[0], size[1]); err == nil || img != nil {
			t.Errorf("size %v: image %v, error %v", size, img, err)
		}
	}
	if _, err := r.PageSize(v, 1, 10, 10); !errors.Is(err, ErrNoPage) {
		t.Errorf("missing page: %v", err)
	}
	if _, err := r.PageSize(nil, 0, 10, 10); err == nil {
		t.Error("nil view: expected error")
	}
}
