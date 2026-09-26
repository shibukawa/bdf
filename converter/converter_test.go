package converter

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestPageRange(t *testing.T) {
	if got, err := PageRange("1-2,5,7-", 8); err != nil || len(got) != 5 || got[4] != 8 {
		t.Fatalf("PageRange = %v %v", got, err)
	}
	if _, err := PageRange("x", 8); err == nil {
		t.Fatal("bad range accepted")
	}
}

func TestDetect(t *testing.T) {
	pdf := []byte("%PDF-1.7\n...")
	if f := Detect(bytes.NewReader(pdf), int64(len(pdf))); f != PDF {
		t.Fatalf("pdf detected as %q", f)
	}
	var b bytes.Buffer
	zw := zip.NewWriter(&b)
	w, _ := zw.Create("ppt/presentation.xml")
	w.Write([]byte("<p:presentation/>"))
	zw.Close()
	if f := Detect(bytes.NewReader(b.Bytes()), int64(b.Len())); f != PPTX {
		t.Fatalf("pptx detected as %q", f)
	}
	for _, name := range []string{"multipage.drawio", "embedded.drawio.svg", "embedded.drawio.png"} {
		data, err := os.ReadFile(filepath.Join("drawio", "testdata", name))
		if err != nil {
			t.Fatal(err)
		}
		if f := Detect(bytes.NewReader(data), int64(len(data))); f != Drawio {
			t.Fatalf("%s detected as %q", name, f)
		}
	}
	for s, want := range map[string]Format{
		"\ufeff<?xml version=\"1.0\"?>\n<mxGraphModel><root/></mxGraphModel>": Drawio,
		"<svg xmlns=\"http://www.w3.org/2000/svg\"/>":                         Unknown,
	} {
		if f := Detect(bytes.NewReader([]byte(s)), int64(len(s))); f != want {
			t.Fatalf("%q detected as %q", s, f)
		}
	}
	junk := []byte("hello")
	if f := Detect(bytes.NewReader(junk), int64(len(junk))); f != Unknown {
		t.Fatalf("junk detected as %q", f)
	}
}
