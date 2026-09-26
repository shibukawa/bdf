package canvas

import (
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
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
	fonts := fontset.New(fontdb.New(nil, false), nil)
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
