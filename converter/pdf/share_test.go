package pdf

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf"
)

// reportlab-master.pdf draws the same header, logo strip and footer on every
// page before the page's own text (see test/pdf/gen_master.py).
func TestSharePrefixAcrossPages(t *testing.T) {
	shared, rs := convert(t, "reportlab-master.pdf")
	if shared.SharedPrefixes != 1 {
		t.Fatalf("SharedPrefixes = %d, want 1 (saved %d bytes)", shared.SharedPrefixes, shared.SharedBytes)
	}
	plain, rp := convertWith(t, "reportlab-master.pdf", &Options{NoSharePrefix: true})
	if plain.SharedPrefixes != 0 {
		t.Fatalf("NoSharePrefix still shared")
	}
	view := shared.Doc.Views[0]
	var prefixes int
	for i, page := range view.Pages {
		obj := decodeObj(t, rs, page.Layers[0].Obj)
		ins, err := obj.Instructions()
		if err != nil {
			t.Fatal(err)
		}
		if ins[0].Op != bdf.OpUse {
			t.Errorf("page %d does not start with USE: %s", i, disasm(t, rs, page.Layers[0].Obj))
			continue
		}
		prefixes++
		// The same runs, in the same order and at the same places, as the
		// unshared conversion.
		want, err := bdf.ExtractText(decodeObj(t, rp, plain.Doc.Views[0].Pages[i].Layers[0].Obj), func(h bdf.Hash) *bdf.ObjectPart { return decodeObj(t, rp, h) })
		if err != nil {
			t.Fatal(err)
		}
		got, err := bdf.ExtractText(obj, func(h bdf.Hash) *bdf.ObjectPart { return decodeObj(t, rs, h) })
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != len(want) {
			t.Fatalf("page %d: %d runs, want %d", i, len(got), len(want))
		}
		for j := range want {
			if got[j].Text != want[j].Text || got[j].X != want[j].X || got[j].Y != want[j].Y || got[j].Sep != want[j].Sep || got[j].Ordinal != want[j].Ordinal {
				t.Fatalf("page %d run %d: %+v, want %+v", i, j, got[j], want[j])
			}
		}
	}
	if prefixes != len(view.Pages) {
		t.Fatalf("%d of %d pages share the prefix", prefixes, len(view.Pages))
	}
	if shared.SharedBytes < 3*512 {
		t.Fatalf("SharedBytes = %d", shared.SharedBytes)
	}
	t.Logf("shared %d bytes; sizes: shared %d, plain %d", shared.SharedBytes, len(encodeDoc(t, shared.Doc)), len(encodeDoc(t, plain.Doc)))
}

func convertWith(t *testing.T, name string, opts *Options) (*Result, *bdf.Reader) {
	t.Helper()
	res, err := ConvertFile(filepath.Join("testdata", name), opts)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return res, r
}

func decodeObj(t *testing.T, r *bdf.Reader, h bdf.Hash) *bdf.ObjectPart {
	t.Helper()
	b, err := r.Part(h)
	if err != nil {
		t.Fatal(err)
	}
	o, err := bdf.DecodeObject(b)
	if err != nil {
		t.Fatal(err)
	}
	return o
}

func disasm(t *testing.T, r *bdf.Reader, h bdf.Hash) string {
	t.Helper()
	b, err := r.Part(h)
	if err != nil {
		t.Fatal(err)
	}
	s, err := bdf.Disassemble(b)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func encodeDoc(t *testing.T, d *bdf.Document) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := d.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
