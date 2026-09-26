package converter

import (
	"bytes"
	"io"
	"slices"
	"strings"
	"testing"
)

func TestParsePages(t *testing.T) {
	for spec, want := range map[string]Pages{
		"1-2,5,7-":   {{1, 2}, {5, 5}, {7, 0}},
		" 3 , 1-1 ,": {{3, 3}, {1, 1}},
		"-3":         {{1, 3}},
		"-":          {{1, 0}},
		"":           nil,
		",":          nil,
	} {
		got, err := ParsePages(spec)
		if err != nil || !slices.Equal(got, want) {
			t.Errorf("ParsePages(%q) = %v, %v; want %v", spec, got, err, want)
		}
	}
	for _, spec := range []string{"x", "0", "-0", "5-3", "1-x", "2--3", "1.5", "3x"} {
		if got, err := ParsePages(spec); err == nil {
			t.Errorf("ParsePages(%q) = %v, want an error", spec, got)
		}
	}
	if s := (Pages{{1, 2}, {5, 5}, {7, 0}}).String(); s != "1-2,5,7-" {
		t.Errorf("String() = %q", s)
	}
}

func TestPagesNumbers(t *testing.T) {
	for _, c := range []struct {
		pages Pages
		count int
		want  []int
	}{
		{nil, 5, nil},
		{Pages{{1, 2}, {5, 5}, {7, 0}}, 8, []int{1, 2, 5, 7, 8}},
		{Pages{{2, 0}}, 3, []int{2, 3}},
		// Order and repeats are the converter's business.
		{PageList(3, 1, 3), 5, []int{3, 1, 3}},
		// Past the last page: what the converter reports (or skips), and
		// no more: the first page after the end, whatever the range says.
		{Pages{{9, 9}}, 5, []int{9}},
		{Pages{{3, 7}}, 5, []int{3, 4, 5, 6}},
		{Pages{{8, 9}}, 5, []int{8}},
		{Pages{{1, 1 << 30}}, 3, []int{1, 2, 3, 4}},
		{Pages{{7, 0}}, 5, []int{7}},
	} {
		got := c.pages.Numbers(c.count)
		if !slices.Equal(got, c.want) || (got == nil) != (c.want == nil) {
			t.Errorf("%v.Numbers(%d) = %v, want %v", c.pages, c.count, got, c.want)
		}
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
