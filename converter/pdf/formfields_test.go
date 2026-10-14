package pdf

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// TestFormFields reads the fields of the form of a PDF: their kinds, what
// the file says of their text, and where their widgets are on the pages.
func TestFormFields(t *testing.T) {
	f, err := os.Open("../pdfform/testdata/application.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	fields, err := FormFields(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, fd := range fields {
		s := fmt.Sprintf("%s %s %q", fd.Name, fd.Kind, fd.Value)
		if fd.Kind == FormText || fd.Kind == FormChoice {
			s += fmt.Sprintf(" %s %g", fd.FontName, fd.FontSize)
		}
		for _, flag := range []struct {
			on   bool
			name string
		}{{fd.Multiline, "multiline"}, {fd.Comb, "comb"}, {fd.Required, "required"}} {
			if flag.on {
				s += " " + flag.name
			}
		}
		if fd.MaxLen > 0 && fd.Comb {
			s += fmt.Sprintf(" max %d", fd.MaxLen)
		}
		if len(fd.Options) > 0 {
			var o []string
			for _, op := range fd.Options {
				o = append(o, op.Value+"="+op.Label)
			}
			s += " [" + strings.Join(o, " ") + "]"
		}
		for _, w := range fd.Widgets {
			x, y := w.Matrix[4], w.Matrix[5]
			s += fmt.Sprintf(" @%d:%.0f,%.0f %.0fx%.0f", w.Page+1, x, y, w.W, w.H)
			if w.State != "" {
				s += "/" + w.State
			}
		}
		got = append(got, s)
	}
	want := []string{
		`applicant.name text "" Helvetica 11 @1:110,70 220x20 @2:110,70 220x20`,
		`applicant.kana text "" Helvetica 0 @1:110,100 220x20`,
		`zip text "" Courier 12 comb max 7 @1:110,130 126x20`,
		`address text "" Helvetica 10 multiline @1:110,160 260x52`,
		`amount text "0" Helvetica 11 @1:110,222 100x20`,
		`plan radio "basic" required [basic=basic pro=pro team=team] @1:110,262 14x14/basic @1:180,262 14x14/pro @1:250,262 14x14/team`,
		`country choice "jp" Helvetica 10 [jp=Japan us=United States de=Germany] @1:110,286 140x20`,
		`agree checkbox "Off" required @1:110,326 14x14/Yes`,
		`newsletter checkbox "Yes" required @1:110,350 14x14/Yes`,
		`notes text "" Helvetica 0 multiline @2:40,130 340x130`,
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("fields:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if len(fields) > 0 && (fields[0].Label != "Name of the applicant" || fields[4].Color != bdf.RGB(0, 0x33, 0x99)) {
		t.Errorf("the label %q and the color %08x", fields[0].Label, uint32(fields[4].Color))
	}
	// a file without a form has no fields
	g, err := os.Open("testdata/reportlab-master.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	if fields, err := FormFields(g, nil); err != nil || len(fields) != 0 {
		t.Errorf("a file without a form: %d fields, %v", len(fields), err)
	}
}

// TestOnForm converts a form with what its widgets show told once the
// fields are read: the file is read once for a filled form. An error from
// the caller ends the conversion.
func TestOnForm(t *testing.T) {
	f, err := os.Open("../pdfform/testdata/application.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var names []string
	res, err := Convert(f, &Options{NoTextIndex: true, OnForm: func(fields []FormField) (map[string]FormShow, error) {
		for _, fd := range fields {
			names = append(names, fd.Name)
		}
		return map[string]FormShow{"country": {Blank: true}}, nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if len(names) < 5 || names[0] != "applicant.name" {
		t.Errorf("the fields: %q", names)
	}
	// the choice of the file is left out of the field whose text is blank
	st, err := res.Doc.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	if text := st.Views[0].Pages[0].Text; strings.Contains(text, "Japan") || !strings.Contains(text, "Application") {
		t.Errorf("the page:\n%s", text)
	}
	if _, err := Convert(f, &Options{OnForm: func([]FormField) (map[string]FormShow, error) { return nil, fmt.Errorf("no form for me") }}); err == nil || !strings.Contains(err.Error(), "no form for me") {
		t.Errorf("the error of the caller: %v", err)
	}
}
