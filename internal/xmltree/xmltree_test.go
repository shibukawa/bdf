package xmltree

import (
	"encoding/xml"
	"errors"
	"strings"
	"testing"
)

func TestElementBudget(t *testing.T) {
	// a root with four children: five elements
	doc := `<r><a/><b><c/></b><d>text</d></r>`
	for _, budget := range []int{5, 100} {
		b := budget
		n, err := ParseCounting(strings.NewReader(doc), nil, &b)
		if err != nil || n == nil || len(n.Kids) != 3 {
			t.Fatalf("budget %d: %v, %+v", budget, err, n)
		}
		if used := budget - b; used != 5 {
			t.Errorf("budget %d: %d elements counted, want 5", budget, used)
		}
	}
	for _, budget := range []int{0, 4} {
		b := budget
		if _, err := ParseCounting(strings.NewReader(doc), nil, &b); !errors.Is(err, ErrTooManyElements) {
			t.Errorf("budget %d: %v, want ErrTooManyElements", budget, err)
		}
	}
	// a streamed element has a budget of its own
	d := xml.NewDecoder(strings.NewReader(doc))
	tok, _ := d.Token()
	if _, err := ReadElement(d, tok.(xml.StartElement)); err != nil {
		t.Fatal(err)
	}
}

// TestManyElements checks that a part of many small elements is refused,
// whatever its markup takes.
func TestManyElements(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("<r>")
	for range MaxElements {
		sb.WriteString("<a/>")
	}
	sb.WriteString("</r>")
	if _, err := Parse([]byte(sb.String())); !errors.Is(err, ErrTooManyElements) {
		t.Fatalf("%d elements: %v, want ErrTooManyElements", MaxElements+1, err)
	}
}
