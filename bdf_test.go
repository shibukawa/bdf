package bdf

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func sample() *Object {
	o := NewObject()
	p := o.AddPath((&Path{}).MoveTo(0, 0).LineTo(10, 0).LineTo(10, 10).Close())
	f := o.AddFont(SystemFont("sans-serif", 400, StyleNormal))
	o.FillColor(RGB(1, 2, 3)).FillPath(p, NonZero)
	o.Dash([]float32{1, 2}, 0.5)
	o.Font(f, 12).FillText("hi", 1, 2, 3)
	o.FillPathRun(EvenOdd, []Glyph{{p, 20, 20}, {p, 30, 20}})
	o.Ext([]byte{1, 2, 3})
	return o
}

func TestRoundTripObject(t *testing.T) {
	data := sample().Encode()
	if !bytes.HasPrefix(data, []byte("BOBJ")) {
		t.Fatal("bad magic")
	}
	obj, err := DecodeObject(data)
	if err != nil {
		t.Fatal(err)
	}
	if obj.BBox != (Rect{0, 0, 40, 30}) {
		t.Fatalf("bbox = %v", obj.BBox)
	}
	ins, err := obj.Instructions()
	if err != nil {
		t.Fatal(err)
	}
	if len(ins) != 7 {
		t.Fatalf("got %d instructions", len(ins))
	}
	s, err := Disassemble(data)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`str[0] = "hi"`, "FILL_COLOR 16909311", "DASH [1 2] 0.5", `FILL_TEXT "hi" 1 2 3`, "FILL_PATH_RUN 1 [{0 20 20} {0 30 20}]", "EXT 3 bytes"} {
		if !strings.Contains(s, want) {
			t.Errorf("disassembly missing %q:\n%s", want, s)
		}
	}
}

func TestDecodeErrors(t *testing.T) {
	data := sample().Encode()
	for n := 0; n < len(data); n++ {
		if _, err := DecodeObject(data[:n]); err == nil {
			t.Fatalf("truncated to %d bytes decoded without error", n)
		}
	}
	bad := append([]byte(nil), data...)
	bad[len(bad)-4] = 0x7E // corrupt EXT opcode into an unknown one
	if o, err := DecodeObject(bad); err == nil {
		if _, err := o.Instructions(); err == nil {
			t.Fatal("unknown opcode accepted")
		}
	}
}

func TestDocumentDedup(t *testing.T) {
	d := NewDocument()
	a, _ := d.AddObject(sample())
	b, _ := d.AddObject(sample())
	if a != b {
		t.Fatal("identical objects got different hashes")
	}
	if len(d.Parts()) != 1 {
		t.Fatalf("parts = %d", len(d.Parts()))
	}
}

func TestSingleAndSplit(t *testing.T) {
	d := NewDocument()
	d.Meta.DC.Title = DCValues{"t"}
	d.Meta.DC.Creator = DCValues{"a", "b"}
	big := NewObject()
	for i := 0; i < 200; i++ {
		big.FillRect(float32(i), 0, 1, 1)
	}
	h, _ := d.AddObject(big)
	v := d.NewView("v", ViewFixed, "")
	v.AddPage(100, 100, Layer{Role: RoleBody, Obj: h})

	var buf bytes.Buffer
	if err := d.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	e, _ := r.Entry(h)
	if e.Enc != EncDeflateRaw {
		t.Fatalf("big object not compressed: %+v", e)
	}
	got, err := r.Part(h)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, big.Encode()) {
		t.Fatal("part bytes differ after round trip")
	}
	if r.Manifest.Views[0].Pages[0].Layers[0].Obj != h {
		t.Fatal("manifest layer hash mismatch")
	}
	if dc := r.Manifest.Meta.DC; dc.Title.First() != "t" || len(dc.Creator) != 2 || dc.Creator[1] != "b" {
		t.Fatalf("dc = %+v", dc)
	}

	dir := t.TempDir()
	d2, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	if err := d2.WriteSplit(dir); err != nil {
		t.Fatal(err)
	}
	r2, err := OpenSplit(dir)
	if err != nil {
		t.Fatal(err)
	}
	d3, err := r2.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	var buf2 bytes.Buffer
	if err := d3.WriteSingle(&buf2); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), buf2.Bytes()) {
		t.Fatal("single -> split -> single is not byte-identical")
	}
}

func TestExtractTextAndIndex(t *testing.T) {
	d := NewDocument()
	child := NewObject()
	cf := child.AddFont(SystemFont("serif", 400, StyleNormal))
	child.Font(cf, 10).FillText("child", 0, 0, 20)
	childH, childBB := d.AddObject(child)

	o := NewObject()
	f := o.AddFont(SystemFont("sans-serif", 400, StyleNormal))
	o.Font(f, 12)
	o.Mark(MarkParagraph, "").FillText("Hello", 0, 0, 30).FillText("World", 35, 0, 30) // gap 5 > 0.2*size => space
	o.Mark(MarkLine, "").FillText("second", 0, 15, 40)
	o.FillText("-line", 40, 15, 20)                     // adjacent => none
	o.Mark(MarkWrap, "").FillText("wrapped", 0, 30, 40) // next line, joined without a separator
	sq := o.AddPath((&Path{}).Rect(0, 0, 1, 1))
	o.Mark(MarkAltText, "outlined").FillPathAt(sq, NonZero, 5, 45) // alt text describes the next drawing op
	o.Save().Translate(100, 100).Use(o.AddObject(childH, childBB)).Restore()
	o.Mark(MarkAltText, "drawn by a child").UseAt(o.AddObject(childH, childBB), 50, 150) // spans the child's bbox
	h, _ := d.AddObject(o)

	obj, _ := DecodeObject(d.Part(h).Data)
	runs, err := ExtractText(obj, func(hh Hash) *ObjectPart {
		p, _ := DecodeObject(d.Part(hh).Data)
		return p
	})
	if err != nil {
		t.Fatal(err)
	}
	var texts []string
	var seps []byte
	for _, r := range runs {
		texts = append(texts, r.Text)
		seps = append(seps, r.Sep)
	}
	wantTexts := []string{"Hello", "World", "second", "-line", "wrapped", "outlined", "child", "drawn by a child"}
	wantSeps := []byte{SepBreak, SepSpace, SepSpace, SepNone, SepNone, SepSpace, SepSpace, SepSpace}
	if strings.Join(texts, "|") != strings.Join(wantTexts, "|") {
		t.Fatalf("texts = %v", texts)
	}
	for i := range wantSeps {
		if seps[i] != wantSeps[i] {
			t.Errorf("run %d (%q) sep = %d, want %d", i, texts[i], seps[i], wantSeps[i])
		}
	}
	if runs[6].X != 100 || runs[6].Y != 100 || runs[6].Font.Family != "serif" {
		t.Fatalf("child run = %+v", runs[6])
	}
	if runs[5].X != 5 || runs[5].Y != 45 || !runs[5].AltText {
		t.Fatalf("alt run = %+v", runs[5])
	}
	if r := runs[7]; r.X != 50+childBB.X || r.Y != 150 || r.Advance != childBB.W || !r.AltText {
		t.Fatalf("alt run of a child = %+v, child bbox %+v", r, childBB)
	}

	v := d.NewView("v", ViewFixed, "")
	v.AddPage(200, 200, Layer{Role: RoleBody, Obj: h})
	ih, err := d.BuildTextIndex(v)
	if err != nil {
		t.Fatal(err)
	}
	if v.TextIndex != ih.String() {
		t.Fatal("view does not reference the index")
	}
	idx, err := DecodeTextIndex(d.Part(ih).Data)
	if err != nil {
		t.Fatal(err)
	}
	if got := PlainText(idx); got != "Hello World second-linewrapped outlined child drawn by a child" {
		t.Fatalf("plain text = %q", got)
	}
	if idx[5].Ordinal != 5 || idx[0].A != 0 || idx[0].B != 0 {
		t.Fatalf("index entry = %+v", idx[5])
	}
}

// Structure marks separate runs like PARAGRAPH; LANG leaves the separator to
// the position guess and keeps a pending ALT_TEXT for its drawing op.
func TestExtractTextStructureMarks(t *testing.T) {
	o := NewObject()
	f := o.AddFont(SystemFont("sans-serif", 400, StyleNormal))
	sq := o.AddPath((&Path{}).Rect(0, 0, 1, 1))
	o.Font(f, 10)
	o.Mark(MarkHeading, "1").FillText("Title", 0, 0, 25)
	o.Mark(MarkList, "").Mark(MarkListItem, "").FillText("one", 0, 15, 15)
	o.Mark(MarkLang, "ja").FillText("two", 17, 15, 15) // same line, gap 2 < 0.2*size => none
	o.Mark(MarkEnd, "").Mark(MarkTable, "").Mark(MarkCell, "A1 col").FillText("head", 0, 30, 20)
	o.Mark(MarkCell, "A2").FillText("cell", 0, 45, 20)
	o.Mark(MarkEnd, "").Mark(MarkAltText, "icon").Mark(MarkLang, "").FillPathAt(sq, NonZero, 0, 60)
	o.Mark(MarkFigure, "chart").FillText("label", 0, 75, 25).Mark(MarkEnd, "")
	obj, _ := DecodeObject(o.Encode())
	runs, err := ExtractText(obj, func(Hash) *ObjectPart { return nil })
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range runs {
		got = append(got, fmt.Sprintf("%s:%d:%v", r.Text, r.Sep, r.AltText))
	}
	want := "Title:2:false one:2:false two:0:false head:2:false cell:2:false icon:2:true label:2:false"
	if strings.Join(got, " ") != want {
		t.Fatalf("runs = %v\nwant  %s", got, want)
	}
	if runs[5].Y != 60 {
		t.Fatalf("ALT_TEXT run lost its drawing op across LANG: %+v", runs[5])
	}
}

func TestExtractContent(t *testing.T) {
	o := NewObject()
	f := o.AddFont(SystemFont("sans-serif", 400, StyleNormal))
	img := o.AddImage(HashOf([]byte("picture")))
	o.Font(f, 10)
	o.Mark(MarkHeading, "2").FillText("Title", 0, 0, 25)
	o.Mark(MarkList, "").Mark(MarkListItem, "").FillText("one", 0, 15, 15)
	o.Mark(MarkList, "").Mark(MarkListItem, "").FillText("nested", 10, 30, 30).Mark(MarkEnd, "")
	o.Mark(MarkListItem, "").FillText("two", 0, 45, 15).Link(0, 37, 15, 10, "https://example.com/").Mark(MarkEnd, "")
	o.Mark(MarkTable, "").FillText("caption", 0, 60, 35)
	o.Mark(MarkCell, "A1 col").FillText("head", 0, 75, 20).Mark(MarkCell, "B1:C1").FillText("wide", 30, 75, 20).Mark(MarkEnd, "")
	o.Mark(MarkParagraph, "").FillText("before", 0, 90, 30).Image(img, 0, 100, 50, 50)
	o.Mark(MarkFigure, "chart").Mark(MarkData, "chart-1").Image(img, 0, 160, 50, 50).FillText("label", 0, 220, 25).Mark(MarkEnd, "")
	o.Mark(MarkLang, "ja").Mark(MarkCell, "B12").FillText("セル", 0, 240, 20)
	o.Mark(MarkLang, "").Mark(MarkBox, "P quote").FillText("said", 0, 255, 20)
	o.Mark(MarkParagraph, "code quote quote").FillText("x := 1", 0, 270, 30)
	o.Mark(MarkList, "ordered quote").Mark(MarkListItem, "").FillText("step", 0, 285, 20).Mark(MarkEnd, "")
	obj, _ := DecodeObject(o.Encode())
	c, err := ExtractContent(obj, func(Hash) *ObjectPart { return nil })
	if err != nil {
		t.Fatal(err)
	}
	var nodes []string
	for _, n := range c.Nodes {
		s := fmt.Sprintf("%s:%d", n.Kind, n.Parent)
		if n.Level > 0 {
			s += fmt.Sprintf(":h%d", n.Level)
		}
		if n.Cell != nil {
			s += fmt.Sprintf(":%d,%d+%d,%d%s", n.Cell.Row, n.Cell.Col, n.Cell.Rows, n.Cell.Cols, n.Cell.Scope)
		}
		if n.Alt != "" {
			s += ":" + n.Alt
		}
		if n.Data != "" {
			s += "=" + n.Data
		}
		if n.Code {
			s += ":code"
		}
		if n.Quote > 0 {
			s += fmt.Sprintf(":quote%d", n.Quote)
		}
		if n.List != "" {
			s += ":" + n.List
		}
		nodes = append(nodes, s)
	}
	want := "heading:-1:h2 list:-1 item:1 paragraph:2 list:2 item:4 paragraph:5 item:1 paragraph:7 table:-1 caption:9 cell:9:0,0+1,1col paragraph:11 cell:9:0,1+1,2 paragraph:13 paragraph:-1 figure:-1:chart=chart-1 paragraph:16 cell:-1:11,1+1,1 paragraph:-1:quote1 paragraph:-1:code:quote2 list:-1:quote1:ordered item:21 paragraph:22"
	if got := strings.Join(nodes, " "); got != want {
		t.Fatalf("nodes = %s\nwant    %s", got, want)
	}
	var runs []string
	for _, r := range c.Runs {
		runs = append(runs, fmt.Sprintf("%s@%d%s", r.Text, r.Node, r.Lang))
	}
	if got, want := strings.Join(runs, " "), "Title@0 one@3 nested@6 two@8 caption@10 head@12 wide@14 before@15 label@17 セル@18ja said@19 x := 1@20 step@23"; got != want {
		t.Fatalf("runs = %s\nwant   %s", got, want)
	}
	if len(c.Links) != 1 || c.Links[0].URL != "https://example.com/" || c.Links[0].After != 3 || c.Links[0].Node != 8 {
		t.Fatalf("links = %+v", c.Links)
	}
	if len(c.Images) != 2 {
		t.Fatalf("images = %+v", c.Images)
	}
	// the first image interrupts a paragraph; the second is in the figure, before its text
	if im := c.Images[0]; im.Node != 15 || im.After != 7 || im.Before != 16 || im.X != 0 || im.Y != 100 || im.W != 50 {
		t.Errorf("image 1 = %+v", im)
	}
	if im := c.Images[1]; im.Node != 16 || im.After != 7 || im.Before != 17 {
		t.Errorf("image 2 = %+v", im)
	}
	// ExtractText gives the same runs
	runs2, err := ExtractText(obj, func(Hash) *ObjectPart { return nil })
	if err != nil || len(runs2) != len(c.Runs) {
		t.Fatalf("ExtractText: %d runs, %v", len(runs2), err)
	}
}

func TestParseCellRef(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want CellRef
		ok   bool
	}{
		{"A1", CellRef{0, 0, 1, 1, ""}, true},
		{"b12 row", CellRef{11, 1, 1, 1, "row"}, true},
		{"B2:C3 col", CellRef{1, 1, 2, 2, "col"}, true},
		{"AA1", CellRef{0, 26, 1, 1, ""}, true},
		{"", CellRef{}, false},
		{"12", CellRef{}, false},
		{"A0", CellRef{}, false},
		{"A1 top", CellRef{}, false},
	} {
		got, ok := ParseCellRef(tc.in)
		if ok != tc.ok || got != tc.want {
			t.Errorf("ParseCellRef(%q) = %+v, %v; want %+v, %v", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}
