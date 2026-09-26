package ai

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
)

func read(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// colors returns the fill and stroke colours an object and its children set.
func colors(t *testing.T, r *bdf.Reader, h bdf.Hash, out map[bdf.Color]bool) {
	t.Helper()
	o, err := r.Object(h)
	if err != nil {
		t.Fatal(err)
	}
	err = o.Walk(func(in bdf.Instr) {
		switch in.Op {
		case bdf.OpFillColor, bdf.OpStrokeColor:
			out[bdf.Color(in.Args[0].(uint64))] = true
		case bdf.OpUse, bdf.OpUseAt:
			colors(t, r, o.Objects[int(in.Args[0].(uint64))], out)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
}

func open(t *testing.T, doc *bdf.Document) *bdf.Reader {
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

func TestConvertArtboards(t *testing.T) {
	data := read(t, "artboards.ai")
	if f := converter.Detect(bytes.NewReader(data), int64(len(data))); f == nil || f.Name != "ai" {
		t.Fatalf("detected as %v", f)
	}
	res, err := Convert(bytes.NewReader(data), int64(len(data)), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	r := open(t, res.Doc)
	m := r.Manifest
	if m.Meta.Source != "ai" || m.Meta.DC.Title.First() != "artboards" {
		t.Errorf("meta = %+v", m.Meta)
	}
	// One page per artboard, the size of the artboard (the trim box), not
	// of the media box with the bleed.
	pages := m.Views[0].Pages
	want := [][2]float32{{200, 200}, {300, 150}, {120, 240}}
	if res.Artboards != 3 || len(pages) != 3 {
		t.Fatalf("%d artboards, %d pages", res.Artboards, len(pages))
	}
	for i, p := range pages {
		if p.W != want[i][0] || p.H != want[i][1] {
			t.Errorf("artboard %d: %v × %v, want %v", i+1, p.W, p.H, want[i])
		}
	}
	// The notes layer is hidden: its red is drawn nowhere.
	for i, p := range pages {
		seen := map[bdf.Color]bool{}
		colors(t, r, p.Layers[0].Obj, seen)
		if seen[bdf.RGB(255, 0, 0)] {
			t.Errorf("artboard %d draws the hidden layer", i+1)
		}
		if len(seen) < 2 {
			t.Errorf("artboard %d draws %d colours", i+1, len(seen))
		}
	}

	// The bleed box includes the bleed; pages select artboards.
	res, err = Convert(bytes.NewReader(data), int64(len(data)), &Options{Box: "bleed", Pages: []int{2}})
	if err != nil {
		t.Fatal(err)
	}
	if p := res.Doc.Views[0].Pages; len(p) != 1 || p[0].W != 318 || p[0].H != 168 {
		t.Errorf("artboard 2 with bleed: %+v", p)
	}
}

func TestUnsupported(t *testing.T) {
	for _, c := range []struct {
		file string
		kind string
		err  error
	}{
		{"nopdf.ai", "pdf", ErrNoPDFContent},
		{"legacy.ai", "ps", ErrPostScript},
	} {
		data := read(t, c.file)
		if k := Kind(bytes.NewReader(data), int64(len(data))); k != c.kind {
			t.Errorf("%s: kind %q, want %q", c.file, k, c.kind)
		}
		// Both are Illustrator files, which the registry hands to this
		// converter rather than reporting an unknown format.
		_, err := converter.Convert(bytes.NewReader(data), int64(len(data)), "", nil)
		if !errors.Is(err, c.err) {
			t.Errorf("%s: %v, want %v", c.file, err, c.err)
		}
	}
}
