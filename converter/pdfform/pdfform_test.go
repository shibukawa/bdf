package pdfform

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

func testOptions() *Options {
	return &Options{FontDirs: []string{"../pptx/testdata/fonts"}, NoSystemFonts: true}
}

func template(t testing.TB) *bytes.Reader {
	t.Helper()
	b, err := os.ReadFile("testdata/application.pdf")
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(b)
}

var applicant = map[string]any{
	"applicant":  map[string]any{"name": "Taro Yamada", "kana": "YAMADA TARO"},
	"zip":        "1000001",
	"address":    "1-1 Chiyoda, Chiyoda-ku, Tokyo\nSecond line of the address",
	"amount":     12800,
	"plan":       "pro",
	"country":    "us",
	"agree":      true,
	"newsletter": false,
	"notes":      "Notes of the applicant, long enough to be broken into lines by the width of the field that they are in.",
}

// pageText returns the text of the pages of a document.
func pageText(t *testing.T, doc *bdf.Document) []string {
	t.Helper()
	st, err := doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, pg := range st.Views[0].Pages {
		out = append(out, pg.Text)
	}
	return out
}

// TestInspect reads the fields of a form as the names that data has.
func TestInspect(t *testing.T) {
	fields, err := Inspect(template(t), nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, f := range fields {
		s := fmt.Sprintf("%s %s %v", f.Name, f.Type, f.At)
		if len(f.Options) > 0 {
			s += fmt.Sprintf(" %v %v", f.Options, f.Labels)
		}
		if f.Default != nil {
			s += fmt.Sprintf(" default=%v", f.Default)
		}
		if f.MaxLen > 0 && f.MaxLen < 100 {
			s += fmt.Sprintf(" max=%d", f.MaxLen)
		}
		if f.Multiline {
			s += " multiline"
		}
		if f.Label != "" {
			s += " (" + f.Label + ")"
		}
		got = append(got, s)
	}
	want := []string{
		"applicant.name string [page 1 page 2] (Name of the applicant)",
		"applicant.kana string [page 1]",
		"zip string [page 1] max=7 (Seven digits)",
		"address string [page 1] multiline",
		"amount string [page 1] default=0",
		"plan string [page 1] [basic pro team] [] default=basic",
		"country string [page 1] [jp us de] [Japan United States Germany] default=jp",
		"agree boolean [page 1] default=false (Agrees to the terms)",
		"newsletter boolean [page 1] default=true",
		"notes string [page 2] multiline",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("fields:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	for _, f := range fields {
		if f.Kind != tag.FieldValue {
			t.Errorf("%s is a %s", f.Name, f.Kind)
		}
	}
}

// TestFill fills the form: the values are the text of its pages, in the
// places of the fields.
func TestFill(t *testing.T) {
	res, err := Fill(template(t), &TemplateData{Data: applicant}, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %q", res.Warnings)
	}
	if out := os.Getenv("PDFFORM_OUT"); out != "" {
		var buf bytes.Buffer
		res.Doc.WriteSingle(&buf)
		os.WriteFile(out, buf.Bytes(), 0o644)
	}
	text := pageText(t, res.Doc)
	if len(text) != 2 {
		t.Fatalf("%d pages", len(text))
	}
	for _, want := range []string{"Application", "Taro Yamada", "YAMADA TARO", "1 0 0 0 0 0 1", "1-1 Chiyoda, Chiyoda-ku, Tokyo Second line of the address", "12800", "United States", "I agree to the terms"} {
		if !strings.Contains(text[0], want) {
			t.Errorf("the first page lacks %q:\n%s", want, text[0])
		}
	}
	// the field of two places has its value in both; a long text is broken into lines
	if !strings.Contains(text[1], "Taro Yamada") || !strings.Contains(text[1], "Notes of the applicant, long enough to be broken into lines") {
		t.Errorf("the second page:\n%s", text[1])
	}
	// the sizes: those of the fields, and for a field without one, a size
	// that fills a field of one line and that keeps the lines of the others in
	first, second := drawn(t, res.Doc, 0), drawn(t, res.Doc, 1)
	if first["Taro Yamada"] != 11 || first["12800"] != 11 || math.Abs(first["YAMADA TARO"]-19/1.35) > 0.01 || first["1"] != 12 {
		t.Errorf("the sizes of the first page: %v", first)
	}
	if second["Notes of the applicant, long enough to be broken into"] != 12 || second["lines by the width of the field that they are in."] != 12 {
		t.Errorf("the lines of the notes: %v", second)
	}
	// the buttons show the states of the data: the dot of the second
	// radio button, the mark of the first checkbox, and no other
	img, err := imagebdf.New(res.Doc, &imagebdf.Options{FontDirs: testOptions().FontDirs, NoSystemFonts: true}).Page(res.Doc.Views[0], 0, 2)
	if err != nil {
		t.Fatal(err)
	}
	ink := func(x, y, w, h float64) int {
		n := 0
		for py := int(y * 2); py < int((y+h)*2); py++ {
			for px := int(x * 2); px < int((x+w)*2); px++ {
				if c := img.RGBAAt(px, py); int(c.R)+int(c.G)+int(c.B) < 200 {
					n++
				}
			}
		}
		return n
	}
	// (the middles of the buttons: where a dot or a mark is)
	if a, b, c := ink(114, 266, 6, 6), ink(184, 266, 6, 6), ink(254, 266, 6, 6); a != 0 || b == 0 || c != 0 {
		t.Errorf("the dots of the radio buttons: %d %d %d", a, b, c)
	}
	if a, b := ink(113, 329, 8, 8), ink(113, 353, 8, 8); a == 0 || b != 0 {
		t.Errorf("the marks of the checkboxes: %d %d", a, b)
	}
	// (the value the file had for a field is not under the new one)
	if strings.Contains(text[0], "Japan") {
		t.Errorf("the choice of the file is left on the page:\n%s", text[0])
	}
}

// drawn returns the runs of text that the last layer of a page draws (the
// values), with the sizes of their fonts.
func drawn(t *testing.T, doc *bdf.Document, page int) map[string]float64 {
	t.Helper()
	layers := doc.Views[0].Pages[page].Layers
	o, err := bdf.DecodeObject(doc.Part(layers[len(layers)-1].Obj).Data)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]float64{}
	size := 0.0
	o.Walk(func(in bdf.Instr) {
		switch in.Op {
		case bdf.OpFont:
			size = float64(in.Args[1].(float32))
		case bdf.OpFillText:
			out[in.Args[0].(string)] = size
		}
	})
	return out
}

// TestFillValues fills the form with values of other kinds, and without
// some: what a field takes, and what is said of what it does not.
func TestFillValues(t *testing.T) {
	// the full names of the fields as the keys of the data; a box by the name of its state; a number in a text
	res, err := Fill(template(t), &TemplateData{Data: map[string]any{
		"applicant.name": "Flat", "applicant": map[string]any{"name": "Nested", "kana": 42}, "agree": "Yes", "newsletter": "off", "plan": "gold", "zip": 1234567890,
	}}, testOptions())
	if err != nil {
		t.Fatal(err)
	}
	text := strings.Join(pageText(t, res.Doc), "\n")
	if !strings.Contains(text, "Flat") || strings.Contains(text, "Nested") || !strings.Contains(text, "42") || !strings.Contains(text, "1 2 3 4 5 6 7") || strings.Contains(text, "8") {
		t.Errorf("the pages:\n%s", text)
	}
	// the fields without a value are as the file has them, and are told; so is a value that is no button
	var missing []string
	other := ""
	for _, w := range res.Warnings {
		if strings.Contains(w, "the data has no") {
			missing = append(missing, w[strings.Index(w, `"`)+1:strings.LastIndex(w[:strings.Index(w, ";")], `"`)])
		} else {
			other += w + "\n"
		}
	}
	if strings.Join(missing, " ") != "address amount country notes" || !strings.Contains(other, `"gold" is not a button of "plan" (basic, pro, team)`) {
		t.Errorf("warnings: %q", res.Warnings)
	}
	if !strings.Contains(text, "Japan") {
		t.Errorf("a field without a value does not show what the file has:\n%s", text)
	}
	// strictly, a field without a value is an error
	if _, err := Fill(template(t), &TemplateData{Data: map[string]any{}, Strict: true}, testOptions()); err == nil || !strings.Contains(err.Error(), `no "applicant.name"`) {
		t.Errorf("strict: %v", err)
	}
	// a PDF without a form is no template
	plain, err := os.ReadFile("../pdf/testdata/reportlab-master.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Fill(bytes.NewReader(plain), nil, testOptions()); err != ErrNoForm {
		t.Errorf("a file without a form: %v", err)
	}
	if _, err := Inspect(bytes.NewReader(plain), nil); err != ErrNoForm {
		t.Errorf("a file without a form: %v", err)
	}
}
