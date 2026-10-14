package formbox

import (
	"bytes"
	"testing"

	"github.com/shibukawa/bdf"
)

func TestPictureCache(t *testing.T) {
	doc := bdf.NewDocument()
	r := New(doc, Options{NoSystemFonts: true})
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg"/>`)
	pic, ok := r.Picture(svg)
	if !ok {
		t.Fatal("SVG was not added")
	}
	copyOfSVG := bytes.Clone(svg)
	again, ok := r.Picture(copyOfSVG)
	if !ok || pic != again || len(r.images) != 1 || len(doc.Parts()) != 1 {
		t.Fatalf("same picture was not reused: %v, %v, %d cache entries, %d parts", again, ok, len(r.images), len(doc.Parts()))
	}
	for range 2 {
		if _, ok := r.Picture([]byte("invalid")); ok {
			t.Fatal("invalid picture was added")
		}
	}
	if len(r.images) != 2 || len(doc.Parts()) != 1 {
		t.Errorf("invalid-picture cache: %d entries, %d parts", len(r.images), len(doc.Parts()))
	}
}

func BenchmarkPictureCache(b *testing.B) {
	data := append([]byte("<svg>"), bytes.Repeat([]byte(" "), 1<<20)...)
	data = append(data, []byte("</svg>")...)
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for b.Loop() {
		r := New(bdf.NewDocument(), Options{NoSystemFonts: true})
		r.Picture(data)
	}
}
