package markdown

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yuin/goldmark/ast"
)

// Many repeated headings still receive unique IDs in order.
func TestRepeatedHeadings(t *testing.T) {
	const size = 12000
	g := &githubIDs{used: map[string]bool{}, next: map[string]int{}}
	var last []byte
	for range size {
		last = g.Generate([]byte("Same Heading"), ast.KindHeading)
	}
	if want := fmt.Sprintf("same-heading-%d", size-1); string(last) != want {
		t.Errorf("the last of %d headings has the id %s, want %s", size, last, want)
	}
}

// Unclosed formulas remain text even in a long document.
func TestUnclosedFormulas(t *testing.T) {
	const size = 80000
	out, err := ToHTML([]byte(strings.Repeat("$`x ", size) + "\n\n$`y`$ " + strings.Repeat("$`x ", size) + "\n"))
	if err != nil {
		t.Fatal(err)
	}
	if n := bytes.Count(out, []byte("$")); n != 2*size || bytes.Count(out, []byte("<math")) != 1 {
		t.Errorf("%d dollar signs of %d, %d formulas", n, 2*size, bytes.Count(out, []byte("<math")))
	}
}
