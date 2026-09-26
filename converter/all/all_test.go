package all

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/shibukawa/bdf/converter"
)

func TestDetect(t *testing.T) {
	var pptx bytes.Buffer
	zw := zip.NewWriter(&pptx)
	w, _ := zw.Create("ppt/presentation.xml")
	w.Write([]byte("<p:presentation/>"))
	zw.Close()
	var xlsx bytes.Buffer
	zw = zip.NewWriter(&xlsx)
	w, _ = zw.Create("xl/workbook.xml")
	w.Write([]byte("<workbook/>"))
	zw.Close()
	var docx bytes.Buffer
	zw = zip.NewWriter(&docx)
	w, _ = zw.Create("word/document.xml")
	w.Write([]byte("<w:document/>"))
	zw.Close()
	var vsdx bytes.Buffer
	zw = zip.NewWriter(&vsdx)
	w, _ = zw.Create("visio/document.xml")
	w.Write([]byte("<VisioDocument/>"))
	zw.Close()
	vdx := []byte(`<?xml version='1.0' encoding='utf-8' ?><VisioDocument xmlns='http://schemas.microsoft.com/visio/2003/core'>`)
	emf := make([]byte, 88)
	binary.LittleEndian.PutUint32(emf, 1)
	copy(emf[40:], " EMF")
	wmf := []byte{1, 0, 9, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for _, c := range []struct {
		name string
		data []byte
		want string
	}{
		{"pdf", []byte("%PDF-1.7\n..."), "pdf"},
		{"pptx", pptx.Bytes(), "pptx"},
		{"xlsx", xlsx.Bytes(), "xlsx"},
		{"docx", docx.Bytes(), "docx"},
		{"vsdx", vsdx.Bytes(), "visio"},
		{"vdx", vdx, "visio"},
		{"emf", emf, "emf"},
		{"wmf", wmf, "emf"},
		{"csv", []byte("id,name\n1,Ann\n2,Bob\n"), "csv"},
		{"drawio", read(t, "multipage.drawio"), "drawio"},
		{"drawio svg", read(t, "embedded.drawio.svg"), "drawio"},
		{"drawio png", read(t, "embedded.drawio.png"), "drawio"},
		{"mxGraphModel", []byte("\ufeff<?xml version=\"1.0\"?>\n<mxGraphModel><root/></mxGraphModel>"), "drawio"},
		// images come after the formats that read some of them (draw.io's exports)
		{"plain svg", []byte(`<svg xmlns="http://www.w3.org/2000/svg"/>`), "image"},
		{"svg with commas", []byte("<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0,0,10,10\">\n<path d=\"M1,1 L2,2\"/>\n<path d=\"M3,3 L4,4\"/>\n</svg>\n"), "image"},
		{"png", readImage(t, "tags.png"), "image"},
		{"jpeg", readImage(t, "photo.jpg"), "image"},
		{"gif", readImage(t, "anim.gif"), "image"},
		{"webp", readImage(t, "scene.webp"), "image"},
		{"avif", readImage(t, "rotated.avif"), "image"},
		{"bmp", readImage(t, "flag.bmp"), "image"},
		{"ico", readImage(t, "icon.ico"), "image"},
		{"tsv", []byte("id\tname\n1\tAnn\n"), "csv"},
		{"dxf", []byte("  0\r\nSECTION\r\n  2\r\nHEADER\r\n"), "dxf"},
		{"binary dxf", []byte("AutoCAD Binary DXF\r\n\x1a\x00\x00\x00"), "dxf"},
		{"junk", []byte("hello"), ""},
	} {
		got := ""
		if f := converter.Detect(bytes.NewReader(c.data), int64(len(c.data))); f != nil {
			got = f.Name
		}
		if got != c.want {
			t.Errorf("%s detected as %q, want %q", c.name, got, c.want)
		}
	}
}

// read returns a file of the draw.io converter's test data.
func read(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "drawio", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// readImage returns a file of the image converter's test data.
func readImage(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "image", "testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}
