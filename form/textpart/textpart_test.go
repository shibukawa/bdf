package textpart

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

func TestWidth(t *testing.T) {
	for _, c := range []struct {
		r          rune
		narrow, ea int
	}{
		{'A', 1, 1}, {' ', 1, 1}, {'ｱ', 1, 1}, {'漢', 2, 2}, {'あ', 2, 2}, {'Ａ', 2, 2}, {'　', 2, 2},
		{'─', 1, 2}, {'│', 1, 2}, {'○', 1, 2}, {'※', 1, 2}, {'α', 1, 2},
		{'́', 0, 0}, {'​', 0, 0}, {'️', 0, 0}, {'\t', 0, 0}, {0x7f, 0, 0},
		{0x1f600, 2, 2}, {0x20b9f, 2, 2},
	} {
		if got := Width(c.r, false); got != c.narrow {
			t.Errorf("Width(%U, narrow) = %d, want %d", c.r, got, c.narrow)
		}
		if got := Width(c.r, true); got != c.ea {
			t.Errorf("Width(%U, wide) = %d, want %d", c.r, got, c.ea)
		}
	}
}

func sample() *Image {
	return &Image{Version: Version, CPI: 10, LPI: 6, Columns: 40, Lines: 12, Margin: Margin{Left: 2, Top: 1}, Skip: 1, Pages: []Page{
		{Lines: []Line{
			{Text: "納品書 No.12", Size: []Range{{Col: 0, Cols: 12, Size: SizeBig}}},
			{},
			{Text: "┌──┐", Paper: []Range{{Col: 0, Cols: 4}}},
			{Text: "│ペン│ 2", Paper: []Range{{Col: 0, Cols: 1}, {Col: 5, Cols: 1}}},
			{Text: "éx"},
		}},
		{Lines: []Line{{Text: "AB", Size: []Range{{Col: 0, Cols: 2, Size: SizeWide}}}}},
	}}
}

func TestCells(t *testing.T) {
	im := sample()
	got := im.Cells(im.Pages[0].Lines[0])
	want := []Cell{
		{Text: "納", Col: 0, Cols: 4, Size: SizeBig}, {Text: "品", Col: 4, Cols: 4, Size: SizeBig}, {Text: "書", Col: 8, Cols: 4, Size: SizeBig},
		{Text: " ", Col: 12, Cols: 1}, {Text: "N", Col: 13, Cols: 1}, {Text: "o", Col: 14, Cols: 1}, {Text: ".", Col: 15, Cols: 1},
		{Text: "1", Col: 16, Cols: 1}, {Text: "2", Col: 17, Cols: 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Cells = %+v\nwant %+v", got, want)
	}
	got = im.Cells(im.Pages[0].Lines[3])
	if !got[0].Paper || got[1].Paper || got[1].Col != 1 || got[3].Col != 5 || !got[3].Paper || got[4].Paper {
		t.Errorf("Cells of the paper: %+v", got)
	}
	got = im.Cells(im.Pages[0].Lines[4])
	if len(got) != 2 || got[0].Text != "é" || got[1].Col != 1 {
		t.Errorf("Cells of a combining mark: %+v", got)
	}
	if !im.Pages[0].Lines[0].Tall() || im.Pages[1].Lines[0].Tall() || im.Rows(im.Pages[0]) != 10 {
		t.Error("Tall or Rows")
	}
}

func TestText(t *testing.T) {
	im := sample()
	got := im.Text(TextOptions{})
	want := "納  品  書   No.12\n\n┌──┐\n│ペン│ 2\néx\n\fA B\n"
	if got != want {
		t.Errorf("Text = %q\nwant %q", got, want)
	}
	got = im.Text(TextOptions{Margins: true, NoPaper: true, CRLF: true})
	want = "\r\n  納  品  書   No.12\r\n\r\n\r\n   ペン  2\r\n  éx\r\n\f\r\n  A B\r\n"
	if got != want {
		t.Errorf("Text with margins, without the paper = %q\nwant %q", got, want)
	}
}

func TestEncodeDecode(t *testing.T) {
	im := sample()
	b, err := im.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(b, []byte(`"version":1,"cpi":10,"lpi":6,"columns":40,"lines":12,"margin":{"left":2,"top":1},"skip":1,"pages":[{"lines":[{"text":"納品書 No.12","size":[{"col":0,"cols":12,"size":"big"}]},{"text":""}`)) {
		t.Errorf("the JSON: %s", b)
	}
	back, err := Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(back, im) {
		t.Errorf("Decode(Encode) = %+v", back)
	}
	for name, s := range map[string]string{
		"version":       `{"version":2,"cpi":10,"lpi":6,"columns":80,"lines":66,"margin":{"left":0,"top":0},"pages":[]}`,
		"pitch":         `{"version":1,"cpi":0,"lpi":6,"columns":80,"lines":66,"margin":{"left":0,"top":0},"pages":[]}`,
		"columns":       `{"version":1,"cpi":10,"lpi":6,"columns":0,"lines":66,"margin":{"left":0,"top":0},"pages":[]}`,
		"margins":       `{"version":1,"cpi":10,"lpi":6,"columns":80,"lines":6,"margin":{"left":0,"top":4},"skip":2,"pages":[]}`,
		"long page":     `{"version":1,"cpi":10,"lpi":6,"columns":80,"lines":2,"margin":{"left":0,"top":0},"pages":[{"lines":[{"text":""},{"text":""},{"text":""}]}]}`,
		"overlap":       `{"version":1,"cpi":10,"lpi":6,"columns":80,"lines":0,"margin":{"left":0,"top":0},"pages":[{"lines":[{"text":"ab","paper":[{"col":0,"cols":2},{"col":1,"cols":1}]}]}]}`,
		"empty range":   `{"version":1,"cpi":10,"lpi":6,"columns":80,"lines":0,"margin":{"left":0,"top":0},"pages":[{"lines":[{"text":"ab","paper":[{"col":0,"cols":0}]}]}]}`,
		"unknown size":  `{"version":1,"cpi":10,"lpi":6,"columns":80,"lines":0,"margin":{"left":0,"top":0},"pages":[{"lines":[{"text":"ab","size":[{"col":0,"cols":2,"size":"huge"}]}]}]}`,
		"not an object": `[]`,
	} {
		if _, err := Decode([]byte(s)); err == nil {
			t.Errorf("Decode accepts %s", name)
		}
	}
}

func TestRead(t *testing.T) {
	doc := bdf.NewDocument()
	v := doc.NewView("pages", bdf.ViewFixed, "")
	v.AddPage(288, 144)
	if _, err := FromDocument(doc); err != ErrNone {
		t.Errorf("FromDocument without the part: %v", err)
	}
	im := sample()
	if _, err := im.AddTo(doc); err != nil {
		t.Fatal(err)
	}
	if got, err := FromDocument(doc); err != nil || got.Pages[1].Lines[0].Text != "AB" {
		t.Errorf("FromDocument: %v", err)
	}
	var buf bytes.Buffer
	if err := doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	// a reader of documents opens it, the part aside
	rd, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if len(rd.Manifest.Views) != 1 || len(rd.Manifest.Parts) != 1 || rd.Manifest.Parts[0].T != bdf.PartText {
		t.Errorf("the manifest: %+v", rd.Manifest.Parts)
	}
	got, err := Read(rd)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, im) {
		t.Errorf("Read = %+v", got)
	}
	if !strings.Contains(got.Text(TextOptions{}), "│ペン│ 2") {
		t.Error("the text of the part read")
	}
}
