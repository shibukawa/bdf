package all

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
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
	var p2z bytes.Buffer
	zw = zip.NewWriter(&p2z)
	w, _ = zw.Create("D0PL001Z.P21")
	w.Write([]byte("ISO-10303-21;\nHEADER;\n"))
	zw.Close()
	var gerberZip bytes.Buffer
	zw = zip.NewWriter(&gerberZip)
	w, _ = zw.Create("board/board-F_Cu.gbr")
	w.Write([]byte("%FSLAX46Y46*%\n%MOMM*%\n"))
	zw.Close()
	var drillZip bytes.Buffer
	zw = zip.NewWriter(&drillZip)
	w, _ = zw.Create("NCDRILL.TXT")
	w.Write([]byte("M48\nMETRIC\nT1C0.8\n%\n"))
	zw.Close()
	var epub bytes.Buffer
	zw = zip.NewWriter(&epub)
	w, _ = zw.CreateHeader(&zip.FileHeader{Name: "mimetype", Method: zip.Store})
	w.Write([]byte("application/epub+zip"))
	w, _ = zw.Create("META-INF/container.xml")
	w.Write([]byte("<container/>"))
	zw.Close()
	step := "ISO-10303-21;\nHEADER;\nFILE_DESCRIPTION((''),'2;1');\nFILE_NAME('a','',(''),(''),'','','');\n"
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
		{"csv", []byte("id,name\n1,Ann\n2,Bob\n"), "csv"},
		{"drawio", read(t, "multipage.drawio"), "drawio"},
		{"drawio svg", read(t, "embedded.drawio.svg"), "drawio"},
		{"drawio png", read(t, "embedded.drawio.png"), "drawio"},
		{"mxGraphModel", []byte("\ufeff<?xml version=\"1.0\"?>\n<mxGraphModel><root/></mxGraphModel>"), "drawio"},
		{"cgm", []byte{0x00, 0x23, 0x02, 'm', 'f', 0x00, 0x10, 0x22, 0x00, 0x04, 0x00, 0x40}, "cgm"},
		{"cgm clear text", []byte("BEGMF 'drawing';\nMFVERSION 1;\n"), "cgm"},
		// a comma on every line, as CSV has
		{"cgm clear text of points", []byte("BEGMF 'a,b';\nVDCEXT (0,0) (100,100);\nLINE (0,0) (10,10);\nLINE (5,0) (5,10);\n"), "cgm"},
		// draw.io's PNG and SVG exports refine images: they are asked first
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
		{"jww", []byte("JwwData.\xbc\x02\x00\x00"), "jww"},
		{"sxf p21", []byte(step + "FILE_SCHEMA(('ASSOCIATIVE_DRAUGHTING'));\nENDSEC;\n"), "sxf"},
		{"sxf p2z", p2z.Bytes(), "sxf"},
		{"sxf p21 of regular lines", []byte(step + "FILE_SCHEMA(('ASSOCIATIVE_DRAUGHTING'));\nENDSEC;\nDATA;\n" +
			strings.Repeat("#10=CARTESIAN_POINT('',(1.,2.));\n", 40)), "sxf"},
		{"step ap214", []byte(step + "FILE_SCHEMA(('AUTOMOTIVE_DESIGN'));\nENDSEC;\n"), ""},
		{"gerber", []byte("%TF.GenerationSoftware,KiCad,Pcbnew,7.0.10*%\n%FSLAX46Y46*%\n%MOMM*%\n"), "gerber"},
		{"gerber of Eagle", []byte("G75*\nG70*\n%OFA0B0*%\n%FSLAX24Y24*%\n%IPPOS*%\n"), "gerber"},
		// a file of aperture definitions, one comma to a line, is not CSV
		{"gerber of apertures", []byte(strings.Repeat("%ADD10R,0.0500X0.0550*%\n", 60) + "%FSLAX24Y24*%\n"), "gerber"},
		{"excellon with commas", []byte("M48\nINCH,LZ,00.0000\n" + strings.Repeat("T1C0.0100\n", 30)), "gerber"},
		{"gerber after a long header", []byte(strings.Repeat("G04 Altium header comment*\n", 60) + "%FSLAX25Y25*%\n%MOIN*%\n"), "gerber"},
		{"excellon", []byte("M48\n; DRILL file {KiCad 7.0.10}\nFMAT,2\nMETRIC\nT1C0.300\n%\n"), "gerber"},
		{"excellon of Eagle", []byte("%\nM48\nM72\nT01C0.0236\n%\n"), "gerber"},
		{"gerber zip", gerberZip.Bytes(), "gerber"},
		{"drill zip", drillZip.Bytes(), "gerber"},
		{"hpgl", []byte("IN;SP1;PA0,0;PD1000,0,1000,1000;PU;"), "hpgl"},
		// semicolons end the instructions, but a plot is not CSV
		{"hpgl in lines", []byte("IN;\nSP1;\nPU0,0;\nPD1000,0;\nPD1000,1000;\nPU;\n"), "hpgl"},
		{"hpgl job", []byte("\x1b%-12345X@PJL ENTER LANGUAGE=HPGL2\r\n\x1bE\x1b%-1BBP;IN;PS1000,1000;"), "hpgl"},
		{"csv with semicolons", []byte("id;name;price\n1;apple;120\n2;pear;90\n"), "csv"},
		{"tiff", []byte("II*\x00\x08\x00\x00\x00"), "tiff"},
		{"big-endian tiff", []byte("MM\x00*\x00\x00\x00\x08"), "tiff"},
		{"bigtiff", []byte("II+\x00\x08\x00\x00\x00\x10\x00\x00\x00\x00\x00\x00\x00"), "tiff"},
		{"epub", epub.Bytes(), "epub"},
		{"parquet", readParquet(t, "basic.parquet"), "parquet"},
		// a file whose footer is encrypted, turned down by the converter
		{"encrypted parquet", []byte("PARE\x00\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00PARE"), "parquet"},
		// the magic number at one end only
		{"truncated parquet", readParquet(t, "basic.parquet")[:1000], ""},
		{"html", []byte("\n<!DOCTYPE html>\n<html><body>x</body></html>"), "html"},
		{"xhtml", []byte(`<?xml version="1.0" encoding="utf-8"?>` + "\n<!-- c -->\n" + `<html xmlns="http://www.w3.org/1999/xhtml">`), "html"},
		{"mhtml", []byte("From: <Saved by Blink>\r\nMIME-Version: 1.0\r\nContent-Type: multipart/related;\r\n\ttype=\"text/html\";\r\n\tboundary=\"b\"\r\n\r\n--b\r\n"), "html"},
		// Markdown and HTML fragments (a README's raw HTML) are text that
		// only a file's extension tells (see TestDetectFile)
		{"markdown", []byte("# Title\n\ntext"), ""},
		{"fragment", []byte(`<p align="center"><img src="logo.png"></p>`), ""},
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

// readParquet returns a file of the Parquet converter's test data.
func readParquet(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "parquet", "testdata", name))
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

// TestDetectFile checks that a file's extension tells the format of text
// whose content no format recognizes.
func TestDetectFile(t *testing.T) {
	dir := t.TempDir()
	for _, c := range []struct{ name, data, want string }{
		{"fragment.html", "<div>\n\n    <p>indented</p>\n</div>", "html"},
		{"notes.md", "<div>\n\n    <p>indented</p>\n</div>", "markdown"},
		{"README.md", "# Title\n\ntext", "markdown"},
		// a Markdown file that starts with an SVG picture is not an SVG file
		{"logo.md", "<svg viewBox=\"0 0 1 1\"><rect/></svg>\n\n# Title\n\ntext", "markdown"},
		{"notes.txt", "plain text", ""},
		// a CSV file of one line is text that only its extension tells
		{"one.csv", "a,b,c", "csv"},
		{"page.htm", "<!DOCTYPE html><p>x", "html"},
		// a drill file without a header
		{"NCDRILL.drl", "T1C0.8\nX1.0Y1.0\nM30\n", "gerber"},
		// by its extension (the HP-GL/2 converter then turns down a gnuplot
		// script)
		{"graph.plt", "set terminal png\nplot sin(x)\n", "hpgl"},
	} {
		path := filepath.Join(dir, c.name)
		if err := os.WriteFile(path, []byte(c.data), 0o644); err != nil {
			t.Fatal(err)
		}
		f, err := converter.DetectFile(path)
		got := ""
		if f != nil {
			got = f.Name
		}
		if err != nil || got != c.want {
			t.Errorf("%s detected as %q (%v), want %q", c.name, got, err, c.want)
		}
	}
}
