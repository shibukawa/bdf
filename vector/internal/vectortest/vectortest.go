// Package vectortest holds what the tests of the vector exporters share:
// the fonts and documents they write, and how their pages are compared.
package vectortest

import (
	"bytes"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// Fonts are the fonts of the tests: text is laid out with them alone, so
// that the files do not depend on the machine.
var Fonts = []string{"../../converter/pptx/testdata/fonts"}

// Load reads a document of the testdata directory at the root.
func Load(t testing.TB, name string) *bdf.Document {
	t.Helper()
	r, err := bdf.OpenSingleFile(filepath.Join("../../testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	doc, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	return doc
}

// One makes a document of one page of w × h, of what draw draws.
func One(t testing.TB, w, h float32, draw func(doc *bdf.Document, o *bdf.Object)) *bdf.Document {
	t.Helper()
	doc := bdf.NewDocument()
	o := bdf.NewObject()
	draw(doc, o)
	o.SetBBox(0, 0, w, h)
	hash, _ := doc.AddObject(o)
	v := doc.NewView("pages", bdf.ViewFixed, "")
	p := v.AddPage(w, h)
	p.Layers = append(p.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: hash})
	return doc
}

// DeepUses makes a document of objects that draw each other twice, forty
// deep: more drawings than a writer writes.
func DeepUses() *bdf.Document {
	doc := bdf.NewDocument()
	leaf := bdf.NewObject()
	leaf.FillRect(0, 0, 1, 1).SetBBox(0, 0, 1, 1)
	h, box := doc.AddObject(leaf)
	for i := 0; i < 40; i++ {
		o := bdf.NewObject()
		r := o.AddObject(h, box)
		o.Use(r).UseAt(r, 1, 1).SetBBox(0, 0, 2, 2)
		h, box = doc.AddObject(o)
	}
	v := doc.NewView("pages", bdf.ViewFixed, "")
	p := v.AddPage(10, 10)
	p.Layers = append(p.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	return doc
}

// Render draws a page of a document with the image renderer, with the
// fonts of the tests.
func Render(t testing.TB, doc *bdf.Document, v *bdf.View, page int) *image.RGBA {
	t.Helper()
	img, err := imagebdf.New(doc, &imagebdf.Options{FontDirs: Fonts, NoSystemFonts: true}).Page(v, page, 1)
	if err != nil {
		t.Fatal(err)
	}
	return img
}

// Distance is the mean difference of the channels of two images, of 255,
// where both are; images whose sizes differ by more than a pixel (a PDF
// reader or Ghostscript rounds the size of a page down, the image renderer
// up) are as far apart as can be.
func Distance(a, b *image.RGBA) float64 {
	w, h := min(a.Rect.Dx(), b.Rect.Dx()), min(a.Rect.Dy(), b.Rect.Dy())
	if w == 0 || h == 0 || max(a.Rect.Dx(), b.Rect.Dx())-w > 1 || max(a.Rect.Dy(), b.Rect.Dy())-h > 1 {
		return 255
	}
	sum := 0.0
	for y := 0; y < h; y++ {
		for x := 0; x < 4*w; x++ {
			sum += math.Abs(float64(a.Pix[y*a.Stride+x]) - float64(b.Pix[y*b.Stride+x]))
		}
	}
	return sum / float64(4*w*h)
}

// Save writes images of a page that differ as PNG files, in the directory
// the environment variable env names (nowhere without it).
func Save(env, name string, imgs ...*image.RGBA) {
	dir := os.Getenv(env)
	if dir == "" {
		return
	}
	for i, img := range imgs {
		var b bytes.Buffer
		png.Encode(&b, img)
		os.WriteFile(filepath.Join(dir, name+"-"+string(rune('a'+i))+".png"), b.Bytes(), 0o644)
	}
}
