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
		{"tsv", []byte("id\tname\n1\tAnn\n"), "csv"},
		{"dxf", []byte("  0\r\nSECTION\r\n  2\r\nHEADER\r\n"), "dxf"},
		{"binary dxf", []byte("AutoCAD Binary DXF\r\n\x1a\x00\x00\x00"), "dxf"},
		{"html", []byte("\n<!DOCTYPE html>\n<html><body>x</body></html>"), "html"},
		{"xhtml", []byte(`<?xml version="1.0" encoding="utf-8"?>` + "\n<!-- c -->\n" + `<html xmlns="http://www.w3.org/1999/xhtml">`), "html"},
		{"mhtml", []byte("From: <Saved by Blink>\r\nMIME-Version: 1.0\r\nContent-Type: multipart/related;\r\n\ttype=\"text/html\";\r\n\tboundary=\"b\"\r\n\r\n--b\r\n"), "html"},
		// any other text is Markdown, also an HTML fragment (a README's raw
		// HTML), but not an XML document
		{"markdown", []byte("# Title\n\ntext"), "markdown"},
		{"fragment", []byte(`<p align="center"><img src="logo.png"></p>`), "markdown"},
		{"text", []byte("hello"), "markdown"},
		{"plain svg", []byte(`<svg xmlns="http://www.w3.org/2000/svg"/>`), ""},
		{"xml", []byte("<?xml version=\"1.0\"?>\n<data/>"), ""},
		{"junk", []byte("\x00\x01\x02\x03"), ""},
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

// TestDetectFile checks that a file's extension decides between the
// fallback format (Markdown takes any text) and the format that lists it.
func TestDetectFile(t *testing.T) {
	dir := t.TempDir()
	for _, c := range []struct{ name, data, want string }{
		{"fragment.html", "<div>\n\n    <p>indented</p>\n</div>", "html"},
		{"notes.md", "<div>\n\n    <p>indented</p>\n</div>", "markdown"},
		{"notes.txt", "plain text", "markdown"},
		// a CSV file of one line is text that only its extension tells
		{"one.csv", "a,b,c", "csv"},
		{"page.htm", "<!DOCTYPE html><p>x", "html"},
	} {
		path := filepath.Join(dir, c.name)
		if err := os.WriteFile(path, []byte(c.data), 0o644); err != nil {
			t.Fatal(err)
		}
		f, err := converter.DetectFile(path)
		if err != nil || f == nil || f.Name != c.want {
			t.Errorf("%s detected as %v (%v), want %s", c.name, f, err, c.want)
		}
	}
}
