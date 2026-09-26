package pptx

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"io"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/imgconv"
)

// emfRecord appends one EMF record made of little-endian 32-bit fields
// followed by raw bytes.
func emfRecord(buf *bytes.Buffer, typ uint32, fields []int32, raw []byte) {
	size := 8 + 4*len(fields) + len(raw)
	pad := (4 - size%4) % 4
	binary.Write(buf, binary.LittleEndian, typ)
	binary.Write(buf, binary.LittleEndian, uint32(size+pad))
	binary.Write(buf, binary.LittleEndian, fields)
	buf.Write(raw)
	buf.Write(make([]byte, pad))
}

// testEMF draws a white background, a red square and a line of text into a
// 100 × 100 pixel frame.
func testEMF() []byte {
	var b bytes.Buffer
	// header: bounds, frame (0.01 mm), " EMF", version, bytes, records,
	// handles, description, palette, device and millimetre sizes
	emfRecord(&b, 1, []int32{0, 0, 99, 99, 0, 0, 2646, 2646, 0x464D4520, 0x10000, 0, 0, 3, 0, 0, 0, 3780, 3780, 1000, 1000}, nil)
	emfRecord(&b, 39, []int32{1, 0, 0xFFFFFF, 0}, nil) // CREATEBRUSHINDIRECT 1: white
	emfRecord(&b, 37, []int32{1}, nil)                 // SELECTOBJECT
	emfRecord(&b, 37, []int32{-0x7FFFFFF8}, nil)       // SELECTOBJECT NULL_PEN
	emfRecord(&b, 43, []int32{0, 0, 100, 100}, nil)    // RECTANGLE
	emfRecord(&b, 39, []int32{2, 0, 0x0000FF, 0}, nil) // CREATEBRUSHINDIRECT 2: red
	emfRecord(&b, 37, []int32{2}, nil)                 // SELECTOBJECT
	emfRecord(&b, 43, []int32{10, 10, 50, 50}, nil)    // RECTANGLE
	emfRecord(&b, 24, []int32{0x008000}, nil)          // SETTEXTCOLOR: green
	emfRecord(&b, 22, []int32{24}, nil)                // SETTEXTALIGN: baseline
	text := utf16.Encode([]rune("EMF text"))
	raw := make([]byte, len(text)*2)
	for i, u := range text {
		binary.LittleEndian.PutUint16(raw[i*2:], u)
	}
	// EXTTEXTOUTW: bounds, graphics mode, scales, reference, count,
	// string offset, options, rectangle, dx offset
	emfRecord(&b, 84, []int32{0, 0, 0, 0, 1, 0, 0, 10, 80, int32(len(text)), 76, 0, 0, 0, 0, 0, 0}, raw)
	emfRecord(&b, 14, []int32{0, 16, 20}, nil) // EOF
	return b.Bytes()
}

// testWMF draws a blue ellipse and a line of text in a 1000 × 1000 window.
func testWMF() []byte {
	var recs bytes.Buffer
	rec := func(fn uint16, args ...int16) {
		binary.Write(&recs, binary.LittleEndian, uint32(3+len(args)))
		binary.Write(&recs, binary.LittleEndian, fn)
		binary.Write(&recs, binary.LittleEndian, args)
	}
	rec(0x020B, 0, 0)                                        // SETWINDOWORG y x
	rec(0x020C, 1000, 1000)                                  // SETWINDOWEXT y x
	rec(0x02FC, 0, 0, 0xFF, 0)                               // CREATEBRUSHINDIRECT: solid blue (0x00FF0000)
	rec(0x012D, 0)                                           // SELECTOBJECT
	rec(0x0418, 900, 900, 100, 100)                          // ELLIPSE bottom right top left
	rec(0x0521, 8, 0x4D57, 0x2046, 0x6574, 0x7478, 500, 100) // TEXTOUT "WMF text" y x
	rec(0x0000)
	var b bytes.Buffer
	// standard header: type, header size, version, size, objects, max record, members
	binary.Write(&b, binary.LittleEndian, []uint16{1, 9, 0x300})
	binary.Write(&b, binary.LittleEndian, uint32(9+recs.Len()/2))
	binary.Write(&b, binary.LittleEndian, []uint16{1})
	binary.Write(&b, binary.LittleEndian, uint32(8))
	binary.Write(&b, binary.LittleEndian, []uint16{0})
	b.Write(recs.Bytes())
	return b.Bytes()
}

func testPNG() []byte {
	img := image.NewNRGBA(image.Rect(0, 0, 4, 4))
	for i := range img.Pix {
		img.Pix[i] = 255 // white
	}
	img.SetNRGBA(1, 1, color.NRGBA{0, 0, 0, 255})
	var b bytes.Buffer
	png.Encode(&b, img)
	return b.Bytes()
}

// pictureDeck replaces the first slide of basic.pptx with pictures.
func pictureDeck(t *testing.T, slide string, media map[string][]byte) []byte {
	t.Helper()
	zr, err := zip.OpenReader("testdata/basic.pptx")
	if err != nil {
		t.Fatal(err)
	}
	defer zr.Close()
	var out bytes.Buffer
	zw := zip.NewWriter(&out)
	put := func(name string, data []byte) {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write(data)
	}
	rels := `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` +
		`<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/slideLayout" Target="../slideLayouts/slideLayout7.xml"/>`
	for name := range media {
		id := strings.TrimSuffix(name, name[strings.LastIndex(name, "."):])
		rels += `<Relationship Id="` + id + `" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="../media/` + name + `"/>`
		put("ppt/media/"+name, media[name])
	}
	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(rc)
		rc.Close()
		switch f.Name {
		case "ppt/slides/slide1.xml":
			data = []byte(`<p:sld xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" xmlns:p="http://schemas.openxmlformats.org/presentationml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><p:cSld><p:spTree><p:nvGrpSpPr><p:cNvPr id="1" name=""/><p:cNvGrpSpPr/><p:nvPr/></p:nvGrpSpPr><p:grpSpPr/>` + slide + `</p:spTree></p:cSld></p:sld>`)
		case "ppt/slides/_rels/slide1.xml.rels":
			data = []byte(rels + `</Relationships>`)
		case "[Content_Types].xml":
			data = bytes.Replace(data, []byte(`<Default Extension="png"`), []byte(`<Default Extension="emf" ContentType="image/x-emf"/><Default Extension="wmf" ContentType="image/x-wmf"/><Default Extension="png"`), 1)
		}
		put(f.Name, data)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}

func pic(id, rid string, x int, effects string) string {
	return `<p:pic><p:nvPicPr><p:cNvPr id="` + id + `" name="Picture ` + id + `"/><p:cNvPicPr/><p:nvPr/></p:nvPicPr>` +
		`<p:blipFill><a:blip r:embed="` + rid + `">` + effects + `</a:blip><a:stretch><a:fillRect/></a:stretch></p:blipFill>` +
		`<p:spPr><a:xfrm><a:off x="` + itoa(x) + `" y="914400"/><a:ext cx="1905000" cy="1905000"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></p:spPr></p:pic>`
}

func TestMetafilesAndRecolor(t *testing.T) {
	transparentWhite := `<a:clrChange><a:clrFrom><a:srgbClr val="FFFFFF"/></a:clrFrom><a:clrTo><a:srgbClr val="FFFFFF"><a:alpha val="0"/></a:srgbClr></a:clrTo></a:clrChange>`
	deck := pictureDeck(t,
		pic("2", "emf1", 457200, transparentWhite)+
			pic("3", "wmf1", 2743200, `<a:duotone><a:srgbClr val="000000"/><a:srgbClr val="00FF00"/></a:duotone>`)+
			pic("4", "png1", 5029200, transparentWhite),
		map[string][]byte{"emf1.emf": testEMF(), "wmf1.wmf": testWMF(), "png1.png": testPNG()})
	res, err := Convert(bytes.NewReader(deck), int64(len(deck)), testOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("warnings: %v", res.Warnings)
	}
	var buf bytes.Buffer
	if err := res.Doc.WriteSingle(&buf); err != nil {
		t.Fatal(err)
	}
	r, err := bdf.OpenSingle(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if text := plainText(t, r); !strings.Contains(text, "EMF text") || !strings.Contains(text, "WMF text") {
		t.Errorf("metafile text missing: %q", text)
	}
	o, err := r.Object(body(r.Manifest.Views[0].Pages[0]))
	if err != nil {
		t.Fatal(err)
	}
	colors := map[uint64]bool{}
	images := 0
	o.Walk(func(in bdf.Instr) {
		switch in.Op {
		case bdf.OpFillColor:
			colors[in.Args[0].(uint64)] = true
		case bdf.OpImage, bdf.OpImageSub:
			images++
		}
	})
	for _, c := range []struct {
		rgba uint64
		what string
	}{
		{0xFFFFFF00, "the EMF background made transparent"},
		{0xFF0000FF, "the red EMF square"},
		{0x008000FF, "the green EMF text"},
		// blue has a luminance of 0.114: 11.4 % of the way to green
		{0x001D00FF, "the WMF ellipse through the duotone"},
	} {
		if !colors[c.rgba] {
			t.Errorf("no fill color %08X (%s); have %v", c.rgba, c.what, colors)
		}
	}
	if colors[0xFFFFFFFF] {
		t.Error("opaque white left in the recolored EMF")
	}
	if images != 1 || len(o.Images) != 1 {
		t.Fatalf("images drawn = %d, parts = %d", images, len(o.Images))
	}
	data, err := r.Part(o.Images[0])
	if err != nil {
		t.Fatal(err)
	}
	img, err := imgconv.Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, _, a := img.At(0, 0).RGBA(); a != 0 {
		t.Errorf("white pixel alpha = %d, want transparent", a)
	}
	if _, _, _, a := img.At(1, 1).RGBA(); a == 0 {
		t.Error("black pixel became transparent")
	}
}
