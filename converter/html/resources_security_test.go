package html

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalImagesStayInDirectory(t *testing.T) {
	parent := t.TempDir()
	dir := filepath.Join(parent, "document")
	if err := os.MkdirAll(filepath.Join(dir, "images"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{filepath.Join(parent, "private.png"), filepath.Join(dir, "safe.png")} {
		if err := os.WriteFile(name, onePixel(), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Symlink(filepath.Join(parent, "private.png"), filepath.Join(dir, "link.png")); err != nil {
		t.Fatal(err)
	}
	r := newResources(&Options{Dir: dir, NoRemote: true})
	for _, name := range []string{"safe.png", "images/../safe.png"} {
		b, err := r.local(&url.URL{Path: name})
		if err != nil || !bytes.Equal(b, onePixel()) {
			t.Errorf("%q: bytes %d, error %v", name, len(b), err)
		}
	}
	for _, name := range []string{"../private.png", "images/../../private.png", "link.png"} {
		if b, err := r.local(&url.URL{Path: name}); err == nil || b != nil {
			t.Errorf("%q escaped the document directory: bytes %d, error %v", name, len(b), err)
		}
	}
}
