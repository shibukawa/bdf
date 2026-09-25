package converter

import (
	"archive/zip"
	"bytes"
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
	junk := []byte("hello")
	if f := Detect(bytes.NewReader(junk), int64(len(junk))); f != Unknown {
		t.Fatalf("junk detected as %q", f)
	}
}
