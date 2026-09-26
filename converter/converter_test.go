package converter

import (
	"bytes"
	"io"
	"strings"
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

func TestRegistry(t *testing.T) {
	var got *Options
	f := &Format{
		Name:   "test-magic",
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return bytes.HasPrefix(head, []byte("MAGIC")) },
		Convert: func(r io.ReaderAt, size int64, o *Options) (*Result, error) {
			got = o
			return &Result{Summary: "ok"}, nil
		},
	}
	Register(f)
	defer func() {
		mu.Lock()
		delete(formats, f.Name)
		mu.Unlock()
	}()
	if Lookup("test-magic") != f {
		t.Error("Lookup")
	}
	found := false
	for _, g := range Formats() {
		found = found || g == f
	}
	if !found {
		t.Error("Formats does not list the format")
	}
	in := []byte("MAGIC and more")
	if Detect(bytes.NewReader(in), int64(len(in))) != f {
		t.Error("not detected")
	}
	junk := []byte("hello")
	if g := Detect(bytes.NewReader(junk), int64(len(junk))); g != nil {
		t.Errorf("junk detected as %s", g.Name)
	}
	res, err := f.Convert(bytes.NewReader(in), int64(len(in)), &Options{Params: map[string]string{"flag": "true", "bad": "maybe"}})
	if err != nil || res.Summary != "ok" {
		t.Fatal(res, err)
	}
	if v, err := got.BoolParam("flag"); !v || err != nil {
		t.Errorf("BoolParam(flag) = %v, %v", v, err)
	}
	if v, err := got.BoolParam("unset"); v || err != nil {
		t.Errorf("BoolParam(unset) = %v, %v", v, err)
	}
	if _, err := got.BoolParam("bad"); err == nil {
		t.Error("BoolParam accepted maybe")
	}
	func() {
		defer func() {
			if r := recover(); r == nil || !strings.Contains(r.(string), "twice") {
				t.Errorf("registering twice: %v", r)
			}
		}()
		Register(f)
	}()
}
