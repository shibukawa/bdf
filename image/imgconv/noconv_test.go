//go:build bdf_noconv

package imgconv

import "testing"

func TestNotAvailable(t *testing.T) {
	if Available() {
		t.Fatal("codecs should be absent")
	}
	img := uiImage(64, 64)
	r, err := EncodeImage(img, true, Options{Mode: Convert})
	if err != ErrNotAvailable || r.Format != "png" || len(r.Data) == 0 {
		t.Fatalf("Convert without codecs: %v %+v", err, r.Format)
	}
	r, err = Optimize([]byte("\x89PNG\r\n\x1a\n"), Options{Mode: Convert})
	if err != ErrNotAvailable || r.Converted {
		t.Fatalf("Optimize without codecs: %v", err)
	}
}
