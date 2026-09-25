package bdf

import (
	"bytes"
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
	d.Meta.Title = "t"
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
	o.FillText("-line", 40, 15, 20) // adjacent => none
	sq := o.AddPath((&Path{}).Rect(0, 0, 1, 1))
	o.Mark(MarkAltText, "outlined").FillPathAt(sq, NonZero, 5, 30) // alt text describes the next drawing op
	o.Save().Translate(100, 100).Use(o.AddObject(childH, childBB)).Restore()
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
	wantTexts := []string{"Hello", "World", "second", "-line", "outlined", "child"}
	wantSeps := []byte{SepBreak, SepSpace, SepSpace, SepNone, SepSpace, SepSpace}
	if strings.Join(texts, "|") != strings.Join(wantTexts, "|") {
		t.Fatalf("texts = %v", texts)
	}
	for i := range wantSeps {
		if seps[i] != wantSeps[i] {
			t.Errorf("run %d (%q) sep = %d, want %d", i, texts[i], seps[i], wantSeps[i])
		}
	}
	if runs[5].X != 100 || runs[5].Y != 100 || runs[5].Font.Family != "serif" {
		t.Fatalf("child run = %+v", runs[5])
	}
	if runs[4].X != 5 || runs[4].Y != 30 || !runs[4].AltText {
		t.Fatalf("alt run = %+v", runs[4])
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
	if got := PlainText(idx); got != "Hello World second-line outlined child" {
		t.Fatalf("plain text = %q", got)
	}
	if idx[5].Ordinal != 5 || idx[0].A != 0 || idx[0].B != 0 {
		t.Fatalf("index entry = %+v", idx[5])
	}
}
