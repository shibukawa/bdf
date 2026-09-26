package emf

import (
	"bytes"
	"encoding/binary"
	"math"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
)

// testFonts are the M PLUS 1p subsets of the PowerPoint converter's tests.
const testFonts = "../pptx/testdata/fonts"

// record appends one EMF record of little-endian 32-bit fields and raw bytes.
func record(buf *bytes.Buffer, typ uint32, fields []int32, raw []byte) {
	size := 8 + 4*len(fields) + len(raw)
	pad := (4 - size%4) % 4
	binary.Write(buf, binary.LittleEndian, typ)
	binary.Write(buf, binary.LittleEndian, uint32(size+pad))
	binary.Write(buf, binary.LittleEndian, fields)
	buf.Write(raw)
	buf.Write(make([]byte, pad))
}

func utf16LE(s string) []byte {
	u := utf16.Encode([]rune(s))
	b := make([]byte, (len(u)*2+3)/4*4)
	for i, c := range u {
		binary.LittleEndian.PutUint16(b[i*2:], c)
	}
	return b
}

// testEMF is a picture of 200 × 100 pixels at 96 dpi (150 × 75 pt) named
// "Sample" with a red square and a line of text.
func testEMF() []byte {
	var b bytes.Buffer
	desc := utf16LE("Test\x00Sample\x00\x00")
	// bounds, frame (0.01 mm), " EMF", version, bytes, records, handles,
	// description length and offset, palette, device and millimetre sizes
	record(&b, 1, []int32{0, 0, 199, 99, 0, 0, 5292, 2646, 0x464D4520, 0x10000, 0, 0, 3,
		int32(len(desc) / 2), 88, 0, 3780, 3780, 1000, 1000}, desc)
	record(&b, 39, []int32{1, 0, 0x0000FF, 0}, nil) // CREATEBRUSHINDIRECT: red
	record(&b, 37, []int32{1}, nil)                 // SELECTOBJECT
	record(&b, 43, []int32{10, 10, 60, 60}, nil)    // RECTANGLE
	face := make([]byte, 8+64)
	copy(face[8:], utf16LE("M PLUS 1p"))
	record(&b, 82, []int32{2, -20, 0, 0, 0, 400}, face) // EXTCREATEFONTINDIRECTW
	record(&b, 37, []int32{2}, nil)
	record(&b, 22, []int32{24}, nil) // SETTEXTALIGN: baseline
	s := "EMF テキスト"
	// EXTTEXTOUTW: bounds, graphics mode, scales, reference, count, string
	// offset, options, rectangle, dx offset
	record(&b, 84, []int32{0, 0, 0, 0, 1, 0, 0, 80, 50, int32(len([]rune(s))), 76, 0, 0, 0, 0, 0, 0}, utf16LE(s))
	record(&b, 14, []int32{0, 16, 20}, nil) // EOF
	return b.Bytes()
}

// testWMF is a placeable WMF of 1440 × 720 twips (72 × 36 pt).
func testWMF() []byte {
	var recs bytes.Buffer
	rec := func(fn uint16, args ...int16) {
		binary.Write(&recs, binary.LittleEndian, uint32(3+len(args)))
		binary.Write(&recs, binary.LittleEndian, fn)
		binary.Write(&recs, binary.LittleEndian, args)
	}
	rec(0x020B, 0, 0)                // SETWINDOWORG y x
	rec(0x020C, 720, 1440)           // SETWINDOWEXT y x
	rec(0x041B, 700, 1400, 100, 100) // RECTANGLE bottom right top left
	rec(0x0000)
	var b bytes.Buffer
	// placeable header: key, handle, bounding box, units per inch, reserved, checksum
	binary.Write(&b, binary.LittleEndian, uint32(0x9AC6CDD7))
	binary.Write(&b, binary.LittleEndian, []int16{0, 0, 0, 1440, 720, 1440})
	binary.Write(&b, binary.LittleEndian, []uint16{0, 0, 0})
	// standard header: type, header size, version, size, objects, max record, members
	binary.Write(&b, binary.LittleEndian, []uint16{1, 9, 0x300})
	binary.Write(&b, binary.LittleEndian, uint32(9+recs.Len()/2))
	binary.Write(&b, binary.LittleEndian, []uint16{0})
	binary.Write(&b, binary.LittleEndian, uint32(7))
	binary.Write(&b, binary.LittleEndian, []uint16{0})
	b.Write(recs.Bytes())
	return b.Bytes()
}

func reopen(t *testing.T, doc *bdf.Document) *bdf.Reader {
	t.Helper()
	var buf bytes.Buffer
	if err := doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestConvertEMF(t *testing.T) {
	data := testEMF()
	res, err := Convert(bytes.NewReader(data), int64(len(data)), &Options{FontDirs: []string{testFonts}, NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	if math.Abs(res.W-150) > 0.1 || math.Abs(res.H-75) > 0.1 {
		t.Errorf("page = %g × %g pt, want 150 × 75", res.W, res.H)
	}
	if res.EmbeddedFonts != 1 {
		t.Errorf("embedded fonts = %d", res.EmbeddedFonts)
	}
	r := reopen(t, res.Doc)
	m := r.Manifest
	if m.Meta.Source != "emf" || m.Meta.DC.Title.First() != "Sample" {
		t.Errorf("source %q, title %q", m.Meta.Source, m.Meta.DC.Title.First())
	}
	if len(m.Views) != 1 || len(m.Views[0].Pages) != 1 {
		t.Fatalf("views = %+v", m.Views)
	}
	h, err := bdf.ParseHash(m.Views[0].TextIndex)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Part(h)
	if err != nil {
		t.Fatal(err)
	}
	idx, err := bdf.DecodeTextIndex(b)
	if err != nil {
		t.Fatal(err)
	}
	if text := bdf.PlainText(idx); !strings.Contains(text, "EMF テキスト") {
		t.Errorf("text = %q", text)
	}
}

func TestConvertWMF(t *testing.T) {
	data := testWMF()
	res, err := Convert(bytes.NewReader(data), int64(len(data)), &Options{Title: "Box", NoSystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.W != 72 || res.H != 36 {
		t.Errorf("page = %g × %g pt, want 72 × 36", res.W, res.H)
	}
	if got := res.Doc.Meta.DC.Title.First(); got != "Box" || res.Doc.Meta.Source != "wmf" {
		t.Errorf("title %q, source %q", got, res.Doc.Meta.Source)
	}
}

func TestNotAMetafile(t *testing.T) {
	data := []byte("%PDF-1.7 and more bytes than a metafile header needs....")
	if _, err := Convert(bytes.NewReader(data), int64(len(data)), nil); err == nil {
		t.Error("a PDF converted")
	}
}

func TestRegistered(t *testing.T) {
	data := testEMF()
	f := converter.Detect(bytes.NewReader(data), int64(len(data)))
	if f == nil || f.Name != "emf" {
		t.Fatalf("detected as %v", f)
	}
	res, err := f.Convert(bytes.NewReader(data), int64(len(data)), &converter.Options{NoSystemFonts: true, SystemFonts: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(res.Summary, "1 page") {
		t.Errorf("summary = %q", res.Summary)
	}
}
