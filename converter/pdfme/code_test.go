package pdfme

import (
	"reflect"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/formbox"
)

// pdfme owns the text below its bars, including whether it is drawn.
// The barcode package's default JAN digits must not add a second label.
func TestEANText(t *testing.T) {
	for _, include := range []bool{false, true} {
		t.Run(map[bool]string{false: "bars", true: "labelled"}[include], func(t *testing.T) {
			r := formbox.New(bdf.NewDocument(), formbox.Options{
				FontDirs: []string{"../docx/testdata/fonts"}, NoSystemFonts: true,
			})
			f := &filler{r: r}
			layer := r.Layer(200, 100)
			s := &schema{name: "jan", typ: "ean13", m: map[string]any{"includetext": include}}
			f.code(layer, s, "490123456789", box{10, 10, 113, 36})
			obj, err := bdf.DecodeObject(layer.Obj.Encode())
			if err != nil {
				t.Fatal(err)
			}
			if len(obj.Paths) != 1 || obj.Paths[0].Inline == nil {
				t.Fatalf("barcode paths: %v", obj.Paths)
			}
			for _, verb := range obj.Paths[0].Inline.Verbs {
				if verb != bdf.VerbRect {
					t.Fatal("the bar path includes a second label")
				}
			}
			var texts []string
			if err := obj.Walk(func(in bdf.Instr) {
				if in.Op == bdf.OpFillText {
					texts = append(texts, in.Args[0].(string))
				}
			}); err != nil {
				t.Fatal(err)
			}
			var want []string
			if include {
				want = []string{"4901234567894"}
			}
			if !reflect.DeepEqual(texts, want) {
				t.Fatalf("labels = %q, want %q", texts, want)
			}
		})
	}
}
