package bdf

import (
	"bytes"
	"slices"
	"testing"
)

// TestByteOrderMark checks that a byte order mark at the start of a string
// of a part is no text: it is not written, and one that another writer
// wrote is dropped when the part is read, as the viewer's decoder drops
// it. A reader that kept it had another text than the viewer, a character
// longer.
func TestByteOrderMark(t *testing.T) {
	o := NewObject()
	o.FillText("\ufeffHello", 0, 10, 0)
	o.FillText("Hello", 0, 20, 0) // the same string
	o.FillText("\ufeff\ufeff", 0, 30, 0)
	o.FillText("a\ufeffb", 0, 40, 0) // inside a string it is a character
	o.Link(0, 0, 1, 1, "\ufeffhttps://example.com/")
	data := o.Encode()
	p, err := DecodeObject(data)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"Hello", "", "a\ufeffb", "https://example.com/"}; !slices.Equal(p.Strings, want) {
		t.Errorf("the strings of the object: %q, want %q", p.Strings, want)
	}
	if bytes.Contains(data, []byte("\xef\xbb\xbfHello")) || bytes.Contains(data, []byte("\xef\xbb\xbfhttps")) {
		t.Error("a mark is written at the start of a string")
	}

	// what another writer wrote: one mark is dropped, and only at the start
	foreign := bytes.Replace(data, []byte("\x05Hello"), []byte("\x08\xef\xbb\xbfHello"), 1)
	foreign = bytes.Replace(foreign, []byte("\x05a\xef\xbb\xbfb"), []byte("\x0b\xef\xbb\xbf\xef\xbb\xbfa\xef\xbb\xbfb"), 1)
	if len(foreign) != len(data)+3+6 {
		t.Fatalf("the object was not rewritten: %d bytes for %d", len(foreign), len(data))
	}
	p, err = DecodeObject(foreign)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"Hello", "", "\ufeffa\ufeffb", "https://example.com/"}; !slices.Equal(p.Strings, want) {
		t.Errorf("the strings of another writer's object: %q, want %q", p.Strings, want)
	}

	// the runs of a text index
	index := EncodeTextIndex([]IndexRun{{Text: "\ufeffabc"}, {Sep: 1, Text: "d\ufeff"}})
	runs, err := DecodeTextIndex(index)
	if err != nil || len(runs) != 2 || runs[0].Text != "abc" || runs[1].Text != "d\ufeff" {
		t.Errorf("the runs of a text index: %+v, %v", runs, err)
	}
	runs, err = DecodeTextIndex(bytes.Replace(index, []byte("\x03abc"), []byte("\x06\xef\xbb\xbfabc"), 1))
	if err != nil || len(runs) != 2 || runs[0].Text != "abc" {
		t.Errorf("the runs of another writer's text index: %+v, %v", runs, err)
	}
}
