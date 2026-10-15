package imagebdf

import "testing"

import "github.com/shibukawa/bdf"

func TestAntialiasUnknownModeWarnsAndPreservesDefault(t *testing.T) {
	o := bdf.NewObject().FillColor(0xff0000ff).FillRect(.25, .25, 3.5, 3.5)
	doc, view := onePage(o, 4, 4)
	r := New(doc, &Options{Antialias: AntialiasMode(255)})
	img, err := r.Page(view, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	want, err := New(doc, nil).Page(view, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	assertImagesEqual(t, img, want)
	if !warned(r, "unknown antialias mode") {
		t.Fatalf("missing invalid mode warning: %v", r.Warnings())
	}
}
