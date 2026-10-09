package ooxml

import (
	"archive/zip"
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/internal/xmltree"
)

// TestLiveElementBudget checks that the trees a package keeps may not
// have more than maxLiveElements elements together, and that discarding a
// tree gives its elements back.
func TestLiveElementBudget(t *testing.T) {
	defer func(n int) { maxLiveElements = n }(maxLiveElements)
	maxLiveElements = 10
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	part := func(name string, elements int) {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte("<r>" + strings.Repeat("<a/>", elements-1) + "</r>"))
	}
	part("six.xml", 6)
	part("five.xml", 5)
	part("four.xml", 4)
	part("eleven.xml", 11)
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	p, err := Open(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.XML("six.xml"); err != nil {
		t.Fatal(err)
	}
	if _, err := p.XML("five.xml"); !errors.Is(err, xmltree.ErrTooManyElements) {
		t.Fatalf("six and five elements: %v, want ErrTooManyElements", err)
	}
	if _, err := p.XML("four.xml"); err != nil {
		t.Fatalf("six and four elements: %v", err)
	}
	p.DiscardXML("six.xml")
	if _, err := p.XML("five.xml"); err != nil {
		t.Fatalf("four and five elements, six given back: %v", err)
	}
	// a part past the budget on its own, also past the budget of one tree
	p.DiscardXML("four.xml")
	p.DiscardXML("five.xml")
	if _, err := p.XML("eleven.xml"); !errors.Is(err, xmltree.ErrTooManyElements) {
		t.Fatalf("eleven elements: %v, want ErrTooManyElements", err)
	}
	// the tree of a part that was refused takes nothing
	if _, err := p.XML("six.xml"); err != nil {
		t.Fatalf("six elements after refusals: %v", err)
	}
}
