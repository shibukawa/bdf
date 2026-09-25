package fixture

import (
	"bytes"
	"testing"

	"github.com/shibukawa/bdf"
)

func TestDeterministic(t *testing.T) {
	var a, b bytes.Buffer
	for _, out := range []*bytes.Buffer{&a, &b} {
		d, err := Demo()
		if err != nil {
			t.Fatal(err)
		}
		if err := d.WriteSingle(out); err != nil {
			t.Fatal(err)
		}
	}
	if !bytes.Equal(a.Bytes(), b.Bytes()) {
		t.Fatal("fixture output is not deterministic")
	}
}

func TestAllObjectsDecode(t *testing.T) {
	d, err := Demo()
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := d.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	objs := 0
	for _, e := range r.Manifest.Parts {
		if e.T != bdf.PartObject {
			continue
		}
		o, err := r.Object(e.H)
		if err != nil {
			t.Fatalf("%s: %v", e.H, err)
		}
		if _, err := o.Instructions(); err != nil {
			t.Fatalf("%s: %v", e.H, err)
		}
		for _, dep := range o.Objects {
			if _, ok := r.Entry(dep); !ok {
				t.Fatalf("%s references missing object %s", e.H, dep)
			}
		}
		objs++
	}
	if objs < 8 {
		t.Fatalf("only %d objects", objs)
	}
	// header and footer are shared across the two flow pages
	doc := r.Manifest.Views[1]
	if doc.Pages[0].Layers[0].Obj != doc.Pages[1].Layers[0].Obj {
		t.Fatal("header not shared")
	}
	if doc.Pages[0].Layers[1].Obj == doc.Pages[1].Layers[1].Obj {
		t.Fatal("bodies unexpectedly identical")
	}
}

func TestAdvance(t *testing.T) {
	d := bdf.NewDocument()
	f, err := LoadFonts(d)
	if err != nil {
		t.Fatal(err)
	}
	w := f.RegularM.Advance("Hello", 12)
	if w < 25 || w > 45 {
		t.Fatalf("advance of Hello at 12pt = %g", w)
	}
}
