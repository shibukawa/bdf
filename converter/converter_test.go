package converter

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
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

// fakeStream is a stream of two empty pages.
type fakeStream struct{}

func (s *fakeStream) Outline() *bdf.Document { return bdf.NewDocument() }
func (s *fakeStream) Pages() int             { return 2 }
func (s *fakeStream) Page(i int) (*bdf.Document, error) {
	return bdf.NewDocument(), nil
}
func (s *fakeStream) Finish() (*Result, error) {
	return &Result{Summary: "streamed", Warnings: []string{"from the stream"}}, nil
}

func TestOpenStream(t *testing.T) {
	whole := &Format{
		Name:   "test-whole",
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return bytes.HasPrefix(head, []byte("WHOLE")) },
		Convert: func(r io.ReaderAt, size int64, o *Options) (*Result, error) {
			return &Result{Doc: bdf.NewDocument(), Summary: "whole"}, nil
		},
	}
	paged := &Format{
		Name:   "test-paged",
		Detect: func(head []byte, r io.ReaderAt, size int64) bool { return bytes.HasPrefix(head, []byte("PAGED")) },
		Convert: func(r io.ReaderAt, size int64, o *Options) (*Result, error) {
			return nil, errors.New("converted whole")
		},
		Stream: func(r io.ReaderAt, size int64, o *Options) (Stream, error) { return &fakeStream{}, nil },
	}
	Register(whole)
	Register(paged)
	defer func() {
		mu.Lock()
		delete(formats, whole.Name)
		delete(formats, paged.Name)
		mu.Unlock()
	}()

	in := []byte("WHOLE document")
	s, err := OpenStream(bytes.NewReader(in), int64(len(in)), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := s.Finish()
	if err != nil || s.Pages() != 0 || res.Summary != "whole" || s.Outline() != res.Doc {
		t.Fatalf("a format without Stream: %d pages, %v, %v", s.Pages(), res, err)
	}
	if _, err := s.Page(0); err == nil {
		t.Error("Page of a document converted whole")
	}

	in = []byte("PAGED document")
	s, err = OpenStream(bytes.NewReader(in), int64(len(in)), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if s.Pages() != 2 {
		t.Fatalf("Pages() = %d", s.Pages())
	}
	if res, err = s.Finish(); err != nil || res.Summary != "streamed" || len(res.Warnings) != 1 {
		t.Fatalf("Finish: %v, %v", res, err)
	}

	junk := []byte("junk")
	if _, err := OpenStream(bytes.NewReader(junk), int64(len(junk)), "", nil); !errors.Is(err, ErrUnknownFormat) {
		t.Errorf("junk: %v", err)
	}
}
