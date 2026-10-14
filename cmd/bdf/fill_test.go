package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFillPicturesInDir reads the pictures that the data names from the
// directory of the template (or of -files) and from it alone: a link that
// leads out of it is not followed.
func TestFillPicturesInDir(t *testing.T) {
	dir := t.TempDir()
	tpl, err := os.ReadFile(filepath.Join("..", "..", "converter", "xlsx", "testdata", "form-fixed.xlsx"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "form.xlsx"), tpl, 0o644); err != nil {
		t.Fatal(err)
	}
	// the picture is outside the directory of the template, and a link in it leads there
	outside := t.TempDir()
	stamp := filepath.Join(outside, "stamp.png")
	if err := os.WriteFile(stamp, testPNG(), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(stamp, filepath.Join(dir, "link.png")); err != nil {
		t.Skip("no symbolic links here:", err)
	}
	data := filepath.Join(dir, "data.json")
	if err := os.WriteFile(data, []byte(`{"no":1,"customer":{"name":"ACME"},"photo":"link.png","paid":true,"tags":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	fonts := filepath.Join("..", "..", "converter", "pptx", "testdata", "fonts")
	status, stderr := run(t, "fill", "-font-dir", fonts, "-no-system-fonts", "-data", data, filepath.Join(dir, "form.xlsx"), filepath.Join(dir, "out.bdf"))
	if status != 0 || !strings.Contains(stderr, `the picture of "photo" is left out`) {
		t.Errorf("status %d:\n%s", status, stderr)
	}
	// a path out of the directory is not read either
	if err := os.WriteFile(data, []byte(`{"no":1,"customer":{"name":"ACME"},"photo":"../`+filepath.Base(outside)+`/stamp.png","paid":true,"tags":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	status, stderr = run(t, "fill", "-font-dir", fonts, "-no-system-fonts", "-data", data, filepath.Join(dir, "form.xlsx"), filepath.Join(dir, "out.bdf"))
	if status != 0 || !strings.Contains(stderr, `the picture of "photo" is left out`) {
		t.Errorf("status %d:\n%s", status, stderr)
	}
	// from the directory of -files it is read
	if err := os.WriteFile(data, []byte(`{"no":1,"customer":{"name":"ACME"},"photo":"stamp.png","paid":true,"tags":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	status, stderr = run(t, "fill", "-font-dir", fonts, "-no-system-fonts", "-files", outside, "-data", data, filepath.Join(dir, "form.xlsx"), filepath.Join(dir, "out.bdf"))
	if status != 0 || strings.Contains(stderr, `"photo"`) {
		t.Errorf("status %d:\n%s", status, stderr)
	}
	// a directory of -files that is none is told
	if status, stderr = run(t, "fill", "-files", filepath.Join(dir, "none"), "-data", data, filepath.Join(dir, "form.xlsx"), filepath.Join(dir, "out.bdf")); status != 2 || !strings.Contains(stderr, "-files:") {
		t.Errorf("status %d:\n%s", status, stderr)
	}
}

// testPNG makes a small picture.
func testPNG() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 8, 4))
	for i := 0; i < 32; i++ {
		img.Set(i%8, i/8, color.RGBA{200, uint8(i * 8), 30, 255})
	}
	var b bytes.Buffer
	png.Encode(&b, img)
	return b.Bytes()
}
