package bdf_test

import (
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

func searchText(t *testing.T, path string) *bdf.SearchText {
	t.Helper()
	r, err := bdf.OpenSingleFile(path)
	if err != nil {
		t.Fatal(err)
	}
	d, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	st, err := d.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestSearchText(t *testing.T) {
	// slides: one entry per page, numbered from 1
	st := searchText(t, "testdata/pptx/basic.bdf")
	if st.Meta.Source != "pptx" || st.Meta.DC.Title.First() != "PowerPoint test deck" {
		t.Errorf("meta %+v", st.Meta)
	}
	if len(st.Views) != 1 || len(st.Views[0].Pages) == 0 {
		t.Fatalf("views %+v", st.Views)
	}
	if p := st.Views[0].Pages[1]; p.Page != 2 || !strings.HasPrefix(p.Text, "Bullets and levels\n") {
		t.Errorf("page 2: %+v", p)
	}
	// Word: the pages, without the scroll view that repeats them
	st = searchText(t, "testdata/docx/basic.bdf")
	if len(st.Views) != 1 || st.Views[0].Kind != bdf.ViewFlow {
		t.Errorf("Word views %+v", st.Views)
	}
	if !strings.Contains(st.Views[0].Pages[0].Text, "文書の行は、変換器が折り返します。") {
		t.Errorf("Word page 1: %q", st.Views[0].Pages[0].Text)
	}
	// a workbook: each sheet whole
	st = searchText(t, "testdata/xlsx/basic.bdf")
	if len(st.Views) != 2 || len(st.Views[0].Pages) != 1 || st.Views[0].Pages[0].Page != 0 || !strings.Contains(st.Views[0].Pages[0].Text, "North ノース") {
		t.Errorf("sheets %+v", st.Views)
	}
	// images: metadata, and no pages of text
	st = searchText(t, "testdata/image/photo.bdf")
	if st.Meta.DC.Description.First() != "A red sun over a green field" || len(st.Views[0].Pages) != 0 {
		t.Errorf("image %+v", st)
	}
}

// TestSearchTextWithoutIndex extracts the text of a view that has no text
// index part from its objects.
func TestSearchTextWithoutIndex(t *testing.T) {
	d := bdf.NewDocument()
	o := bdf.NewObject()
	o.Font(o.AddFont(bdf.SystemFont("serif", 400, bdf.StyleNormal)), 12).FillText("Hello", 10, 20, 0).FillText("world", 50, 20, 0)
	h, _ := d.AddObject(o)
	v := d.NewView("v", bdf.ViewFixed, "")
	v.AddPage(100, 100)
	v.AddPage(100, 100, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	st, err := d.SearchText()
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Views) != 1 || len(st.Views[0].Pages) != 1 || st.Views[0].Pages[0] != (bdf.PageText{Page: 2, Text: "Hello world"}) {
		t.Errorf("%+v", st.Views)
	}
}
