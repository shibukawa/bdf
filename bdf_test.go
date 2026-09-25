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
