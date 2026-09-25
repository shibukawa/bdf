//go:build !bdf_noconv

package imgconv

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"testing"
)

func TestSniff(t *testing.T) {
	cases := map[string][]byte{
		"png":  []byte("\x89PNG\r\n\x1a\n...."),
		"jpeg": []byte{0xff, 0xd8, 0xff, 0xe0, 0, 0},
		"webp": []byte("RIFF\x00\x00\x00\x00WEBPVP8L"),
		"avif": []byte("\x00\x00\x00\x1cftypavif...."),
		"gif":  []byte("GIF89a......"),
		"svg":  []byte("<?xml version=\"1.0\"?><svg xmlns=\"http://www.w3.org/2000/svg\"/>"),
		"":     []byte("hello"),
	}
	for want, b := range cases {
		if got := Sniff(b); got != want {
			t.Errorf("Sniff(%q) = %q, want %q", b[:min(8, len(b))], got, want)
		}
	}
}

func TestEncodeImageLosslessGraphics(t *testing.T) {
	img := uiImage(400, 300)
	keep, err := EncodeImage(img, true, Options{})
	if err != nil || keep.Format != "png" || keep.Converted {
		t.Fatalf("keep mode: %+v %v", keep.Format, err)
	}
	r, err := EncodeImage(img, true, Options{Mode: Convert})
	if err != nil {
		t.Fatal(err)
	}
	if r.Format != "webp" || !r.Lossless || !r.Converted {
		t.Fatalf("graphics should become lossless WebP: %+v", r)
	}
	if len(r.Data) >= len(keep.Data)/5 {
		t.Fatalf("webp %d bytes not much smaller than png %d", len(r.Data), len(keep.Data))
	}
	dec, err := Decode(r.Data)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range [][2]int{{0, 0}, {33, 5}, {150, 150}, {399, 299}} {
		want := img.NRGBAAt(p[0], p[1])
		cr, cg, cb, ca := dec.At(p[0], p[1]).RGBA()
		if uint8(cr>>8) != want.R || uint8(cg>>8) != want.G || uint8(cb>>8) != want.B || uint8(ca>>8) != want.A {
			t.Fatalf("pixel %v changed: %v %v %v %v want %v", p, cr>>8, cg>>8, cb>>8, ca>>8, want)
		}
	}
}

func TestEncodeImagePhoto(t *testing.T) {
	img := photoImage(400, 300)
	if !PhotoLike(img, 4096) || PhotoLike(uiImage(400, 300), 4096) {
		t.Fatal("PhotoLike heuristic is wrong")
	}
	r, err := EncodeImage(img, true, Options{Mode: Convert, Quality: 80, Method: 4})
	if err != nil {
		t.Fatal(err)
	}
	if r.Format != "webp" || r.Lossless {
		t.Fatalf("photo should become lossy WebP: format=%s lossless=%v", r.Format, r.Lossless)
	}
	r2, err := EncodeImage(img, true, Options{Mode: Convert, NoLossyForLossless: true, Method: 1})
	if err != nil {
		t.Fatal(err)
	}
	if !r2.Lossless {
		t.Fatal("NoLossyForLossless ignored")
	}
}

func TestOptimize(t *testing.T) {
	img := photoImage(300, 200)
	var jpg bytes.Buffer
	if err := jpeg.Encode(&jpg, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatal(err)
	}
	r, err := Optimize(jpg.Bytes(), Options{Mode: Convert, Quality: 80, Method: 4})
	if err != nil {
		t.Fatal(err)
	}
	if r.Format != "webp" || !r.Converted || len(r.Data) >= jpg.Len() {
		t.Fatalf("jpeg should shrink to webp: %+v", r.Format)
	}
	// Candidates that are not smaller than the original are rejected.
	if best, err := encodeBest(uiImage(64, 64), true, Options{Mode: Convert}, 10); err != nil || best != nil {
		t.Fatalf("expected no candidate below 10 bytes, got %+v %v", best, err)
	}
	var small bytes.Buffer
	if err := png.Encode(&small, uiImage(4, 4)); err != nil {
		t.Fatal(err)
	}
	r, err = Optimize(small.Bytes(), Options{Mode: Convert})
	if err != nil || r.Format != "webp" && r.Format != "png" {
		t.Fatalf("small png: %+v %v", r.Format, err)
	}
	// Formats browsers decode natively are never touched.
	for _, b := range [][]byte{[]byte("RIFF\x00\x00\x00\x00WEBPVP8L"), []byte("\x00\x00\x00\x1cftypavif"), []byte("<svg/>")} {
		r, err := Optimize(b, Options{Mode: Convert})
		if err != nil || r.Converted || !bytes.Equal(r.Data, b) {
			t.Fatalf("%q was changed: %+v %v", b, r, err)
		}
	}
	// Keep mode returns the input.
	r, _ = Optimize(jpg.Bytes(), Options{})
	if r.Converted || r.Format != "jpeg" {
		t.Fatal("keep mode changed the image")
	}
}
