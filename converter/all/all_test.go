package all

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
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
		{"ai", []byte("%PDF-1.6\n1 0 obj\n<x:xmpmeta><illustrator:Type>Document</illustrator:Type>"), "ai"},
		// A PDF saved by Illustrator is a PDF, even with Illustrator's data in it.
		{"pdf from Illustrator", []byte("%PDF-1.6\n<xmp:CreatorTool>Adobe Illustrator 30.6</xmp:CreatorTool><illustrator:CreatorSubTool>"), "pdf"},
		{"ai 8", []byte("%!PS-Adobe-3.0 \r%%Creator: Adobe Illustrator(R) 8.0\r"), "ai"},
		{"psd", []byte("8BPS\x00\x01\x00\x00\x00\x00\x00\x00"), "psd"},
		{"psb", []byte("8BPS\x00\x02\x00\x00\x00\x00\x00\x00"), "psd"},
		{"pptx", pptx.Bytes(), "pptx"},
		{"xlsx", xlsx.Bytes(), "xlsx"},
		{"docx", docx.Bytes(), "docx"},
		{"vsdx", vsdx.Bytes(), "visio"},
		{"vdx", vdx, "visio"},
		{"emf", emf, "emf"},
		{"wmf", wmf, "emf"},
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
