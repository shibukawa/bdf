// Package convtest holds what the tests of the converters share: reading
// the documents they make back, their text, the files of testdata, and
// the pictures and symbols the tests of forms look for.
package convtest

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/shibukawa/bdf"
)

// Reopen writes a document as a single file and opens it with a reader,
// as a reader of the file would.
func Reopen(t testing.TB, doc *bdf.Document) *bdf.Reader {
	t.Helper()
	var buf bytes.Buffer
	if err := doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// PlainText is the text of the text index of a view, as a search sees it.
func PlainText(t testing.TB, r *bdf.Reader, view int) string {
	t.Helper()
	h, err := bdf.ParseHash(r.Manifest.Views[view].TextIndex)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Part(h)
	if err != nil {
		t.Fatal(err)
	}
	idx, err := bdf.DecodeTextIndex(b)
	if err != nil {
		t.Fatal(err)
	}
	return bdf.PlainText(idx)
}

// Read reads a file of the testdata directory of the package tested.
func Read(t testing.TB, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Allocated returns the bytes allocated on the heap while f runs.
func Allocated(f func()) uint64 {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	f()
	runtime.ReadMemStats(&after)
	return after.TotalAlloc - before.TotalAlloc
}

// Within runs f and fails the test when it takes longer than d.
func Within(t testing.TB, d time.Duration, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() { f(); close(done) }()
	select {
	case <-done:
	case <-time.After(d):
		t.Fatalf("%s did not finish in %v", what, d)
	}
}

// HasWarning tells whether one of the warnings has part in it.
func HasWarning(warnings []string, part string) bool {
	return slices.ContainsFunc(warnings, func(w string) bool { return strings.Contains(w, part) })
}

// PNG makes a picture of a size, of a gradient.
func PNG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w*h; i++ {
		img.Set(i%w, i/w, color.RGBA{200, uint8(i), 30, 255})
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

// CountingFS counts the files opened of it (tag.PictureData opens a file
// to read it).
type CountingFS struct {
	fstest.MapFS
	Reads int
}

func (c *CountingFS) Open(name string) (fs.File, error) {
	c.Reads++
	return c.MapFS.Open(name)
}

// CodeBounds recognizes the symbol of a code in a path: bars or modules
// (rectangles, 8 at least) followed optionally by closed digit outlines,
// which must use curves and stay below the bars. The box holds both, so
// the frame of the symbol is checked by the caller's expected symbol.
func CodeBounds(p *bdf.Path) (rects, contours int, box [4]float32) {
	if p == nil || len(p.Verbs) < 8 || p.Verbs[0] != bdf.VerbRect {
		return
	}
	box = [4]float32{p.Args[0], p.Args[1], p.Args[0], p.Args[1]}
	grow := func(x, y float32) {
		box[0], box[1] = min(box[0], x), min(box[1], y)
		box[2], box[3] = max(box[2], x), max(box[3], y)
	}
	at := 0
	for rects < len(p.Verbs) && p.Verbs[rects] == bdf.VerbRect {
		grow(p.Args[at], p.Args[at+1])
		grow(p.Args[at]+p.Args[at+2], p.Args[at+1]+p.Args[at+3])
		at += 4
		rects++
	}
	if rects < 8 {
		return 0, 0, [4]float32{}
	}
	bottom, curves, closed := box[3], 0, 0
	for _, verb := range p.Verbs[rects:] {
		points := 0
		switch verb {
		case bdf.VerbMove:
			contours++
			points = 1
		case bdf.VerbLine:
			points = 1
		case bdf.VerbQuad:
			curves++
			points = 2
		case bdf.VerbClose:
			closed++
		default:
			return 0, 0, [4]float32{}
		}
		for ; points > 0; points-- {
			x, y := p.Args[at], p.Args[at+1]
			if y <= bottom {
				return 0, 0, [4]float32{}
			}
			grow(x, y)
			at += 2
		}
	}
	if contours > 0 && (curves == 0 || closed != contours) {
		return 0, 0, [4]float32{}
	}
	return
}

// RectBox summarizes the symbol of a code in a path: its bars or modules,
// its digit contours and its size.
func RectBox(p *bdf.Path) string {
	rects, contours, box := CodeBounds(p)
	if rects == 0 {
		return ""
	}
	label := ""
	if contours > 0 {
		label = fmt.Sprintf(" + %d contours below bars", contours)
	}
	return fmt.Sprintf("%d rects of %.2f × %.2f%s", rects, box[2]-box[0], box[3]-box[1], label)
}
