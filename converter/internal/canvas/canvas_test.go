package canvas

import (
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/internal/fontdb"
)

func TestIsDocLang(t *testing.T) {
	for _, c := range []struct {
		lang, doc string
		want      bool
	}{
		{"en-US", "en-US", true}, {"EN-us", "en-US", true}, {"ja", "ja-JP", true}, {"ja-JP", "ja", false},
		{"en-GB", "en-US", false}, {"j", "ja-JP", false}, {"", "en-US", false}, {"", "", true},
	} {
		if got := isDocLang(c.lang, c.doc); got != c.want {
			t.Errorf("isDocLang(%q, %q) = %v", c.lang, c.doc, got)
		}
	}
}

func TestBuilder(t *testing.T) {
	doc := bdf.NewDocument()
	doc.Meta.DC.Language = bdf.DCValues{"ja-JP"}
	fonts := fontset.New(fontdb.New(nil, nil, false), nil)
	b := NewBuilder(doc, fonts)
	cv := b.New()
	fc := fonts.Choose("Arial", false, false, false)
	ref := cv.Font(fc.Use)
	if cv.Font(fc.Use) != ref {
		t.Error("fonts are not shared")
	}
	cv.SetLang("ja") // the document's
	cv.SetLang("x-none")
	cv.SetLang("en-US")
	if cv.Lang() != "en-US" {
		t.Errorf("Lang = %q", cv.Lang())
	}
	ch, _ := cv.Child(bdf.Rect{W: 10, H: 10})
	ch.Obj.Font(ch.Font(fc.Use), 12)
	ch.Obj.FillText("x", 0, 0, 5)
	cv.Obj.Font(ref, 12)
	cv.Transform(Translate(1, 2))
	b.Encode()
	if cv.Hash() == (bdf.Hash{}) || ch.Hash() == (bdf.Hash{}) {
		t.Fatal("objects not encoded")
	}
	o, err := bdf.DecodeObject(doc.Part(cv.Hash()).Data)
	if err != nil {
		t.Fatal(err)
	}
	if len(o.Fonts) != 1 || o.Fonts[0].Family != `"Arial", sans-serif` {
		t.Errorf("fonts = %+v", o.Fonts)
	}
	if len(o.Objects) != 1 || o.Objects[0] != ch.Hash() {
		t.Errorf("child objects = %v", o.Objects)
	}
	var marks []string
	o.Walk(func(in bdf.Instr) {
		if in.Op == bdf.OpMark {
			marks = append(marks, in.Args[1].(string))
		}
	})
	if len(marks) != 1 || marks[0] != "en-US" {
		t.Errorf("LANG marks = %q", marks)
	}
}

// TestFixedFont mixes a font record of the converter's own with the fonts
// of the font set: each keeps its reference.
func TestFixedFont(t *testing.T) {
	doc := bdf.NewDocument()
	fonts := fontset.New(fontdb.New(nil, nil, false), nil)
	b := NewBuilder(doc, fonts)
	cv := b.New()
	own := bdf.EmbeddedFont(doc.AddFont([]byte("font")), 400, bdf.StyleNormal)
	fc := fonts.Choose("Arial", false, false, false)
	r1 := cv.Font(fc.Use)
	r2 := cv.FixedFont(own)
	if cv.FixedFont(own) != r2 || cv.Font(fc.Use) != r1 || r1 == r2 {
		t.Fatalf("references %d %d", r1, r2)
	}
	cv.Obj.Font(r1, 10)
	cv.Obj.FillText("a", 0, 0, 5)
	cv.Obj.Font(r2, 10)
	cv.Obj.FillText("b", 0, 0, 5)
	b.Encode()
	o, err := bdf.DecodeObject(doc.Part(cv.Hash()).Data)
	if err != nil {
		t.Fatal(err)
	}
	if o.Fonts[r2] != own || o.Fonts[r1] == own || o.Fonts[r1].Family == "\x00pending" {
		t.Errorf("fonts %+v", o.Fonts)
	}
}

func TestEncodeAndRelease(t *testing.T) {
	doc := bdf.NewDocument()
	b := NewBuilder(doc, fontset.New(fontdb.New(nil, nil, false), nil))
	parent := b.New()
	child, _ := parent.Child(bdf.Rect{W: 10, H: 10})
	child.Obj.FillRect(0, 0, 10, 10)
	b.EncodeAndRelease()
	for _, cv := range []*Canvas{parent, child} {
		if cv.Obj != nil || cv.children != nil || cv.Hash() == (bdf.Hash{}) || doc.Part(cv.Hash()) == nil {
			t.Fatalf("canvas was not encoded and released: %+v", cv)
		}
	}
	o, err := bdf.DecodeObject(doc.Part(parent.Hash()).Data)
	if err != nil || len(o.Objects) != 1 || o.Objects[0] != child.Hash() {
		t.Fatalf("child reference after release: %v, %v", o, err)
	}
	more := b.New()
	more.Obj.FillRect(0, 0, 2, 2)
	b.EncodeAndRelease()
	if more.Hash() == (bdf.Hash{}) || more.Obj != nil || len(b.canvases) != 0 {
		t.Fatal("builder cannot release a second batch")
	}
}
