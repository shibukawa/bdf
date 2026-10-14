package tag

import (
	"bytes"
	"io"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

func TestPictureData(t *testing.T) {
	files := fstest.MapFS{"images/photo.png": {Data: []byte("picture")}}
	for _, name := range []string{"images/photo.png", `./images\photo.png`, "data:image/png;base64,cGljdHVyZQ=="} {
		got, err := PictureData(files, name)
		if err != nil || !bytes.Equal(got, []byte("picture")) {
			t.Errorf("PictureData(%q) = %q, %v", name, got, err)
		}
	}
	for _, name := range []string{"../secret", "/secret", "images/../secret", "images", "data:image/png,plain", "data:image/png;base64,!"} {
		if _, err := PictureData(files, name); err == nil {
			t.Errorf("PictureData(%q) succeeded", name)
		}
	}
}

// The size from Stat can be stale (or supplied by a filesystem that does
// not know it). Bound the read itself, including when ReadFile is offered.
func TestPictureDataReadLimit(t *testing.T) {
	f := &pictureFile{remaining: MaxPicture + 100}
	files := &pictureFS{file: f}
	data, err := PictureData(files, "photo.png")
	if err == nil || !strings.Contains(err.Error(), "larger") || data != nil {
		t.Fatalf("PictureData = %d bytes, %v; want an oversized-picture error", len(data), err)
	}
	if files.readFile || f.read > MaxPicture+1 || !f.closed {
		t.Errorf("ReadFile called: %v; bytes read: %d; closed: %v", files.readFile, f.read, f.closed)
	}
}

type pictureFS struct {
	file     *pictureFile
	readFile bool
}

func (f *pictureFS) Open(string) (fs.File, error) { return f.file, nil }
func (f *pictureFS) ReadFile(string) ([]byte, error) {
	f.readFile = true
	return io.ReadAll(f.file)
}

type pictureFile struct {
	remaining, read int
	closed          bool
}

func (f *pictureFile) Read(b []byte) (int, error) {
	if f.remaining == 0 {
		return 0, io.EOF
	}
	n := min(len(b), f.remaining)
	clear(b[:n])
	f.remaining -= n
	f.read += n
	return n, nil
}
func (f *pictureFile) Stat() (fs.FileInfo, error) { return pictureInfo{}, nil }
func (f *pictureFile) Close() error               { f.closed = true; return nil }

type pictureInfo struct{}

func (pictureInfo) Name() string       { return "photo.png" }
func (pictureInfo) Size() int64        { return 1 }
func (pictureInfo) Mode() fs.FileMode  { return 0 }
func (pictureInfo) ModTime() time.Time { return time.Time{} }
func (pictureInfo) IsDir() bool        { return false }
func (pictureInfo) Sys() any           { return nil }
