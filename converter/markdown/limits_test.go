package markdown

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// countingIDs are the ids of headings as they were made when the numbers
// of a repeated heading were tried from 1 for each of them.
type countingIDs struct {
	used map[string]bool
}

func (g *countingIDs) Generate(value []byte, kind ast.NodeKind) []byte {
	var b strings.Builder
	for _, r := range strings.ToLower(string(value)) {
		switch {
		case r == ' ':
			b.WriteByte('-')
		case r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.Is(unicode.M, r):
			b.WriteRune(r)
		}
	}
	slug := b.String()
	if slug == "" {
		slug = "heading"
	}
	id := slug
	for n := 1; g.used[id]; n++ {
		id = slug + "-" + strconv.Itoa(n)
	}
	g.used[id] = true
	return []byte(id)
}

func (g *countingIDs) Put(value []byte) { g.used[string(value)] = true }

// rescanningMath reads formulas in a line as they were read when the rest
// of the line was looked through for each dollar sign.
type rescanningMath struct{}

func (rescanningMath) Trigger() []byte { return []byte{'$'} }

func (rescanningMath) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 2 {
		return nil
	}
	var src []byte
	n := 0
	display := false
	switch {
	case line[1] == '$':
		end := bytes.Index(line[2:], []byte("$$"))
		if end <= 0 {
			return nil
		}
		src, n, display = line[2:2+end], end+4, true
	case line[1] == '`':
		end := bytes.Index(line[2:], []byte("`$"))
		if end < 0 {
			return nil
		}
		src, n = line[2:2+end], end+4
	default:
		if isSpaceByte(line[1]) {
			return nil
		}
		end := -1
	scan:
		for i := 1; i < len(line); i++ {
			switch line[i] {
			case '\\':
				i++
			case '$':
				if isSpaceByte(line[i-1]) || i+1 < len(line) && line[i+1] >= '0' && line[i+1] <= '9' {
					return nil
				}
				end = i
				break scan
			}
		}
		if end < 0 {
			return nil
		}
		src, n = line[1:end], end+1
	}
	block.Advance(n)
	return &mathInline{src: append([]byte(nil), src...), display: display}
}

type rescanningExtension struct{}

func (rescanningExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(util.Prioritized(mathBlockParser{}, 690)),
		parser.WithInlineParsers(util.Prioritized(rescanningMath{}, 90)),
		parser.WithASTTransformers(util.Prioritized(mathFences{}, 100)),
	)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(util.Prioritized(mathRenderer{}, 100)))
}

// toHTMLBefore renders a document with the ids and the formulas read as
// they were.
func toHTMLBefore(src []byte) ([]byte, error) {
	return toHTML(src, rescanningExtension{}, &countingIDs{used: map[string]bool{}})
}

// Documents render to the HTML they rendered to when the numbers of
// repeated headings were counted from 1 and the lines looked through for
// each dollar sign.
func TestHTMLAsBefore(t *testing.T) {
	pieces := []string{"# a\n", "# a\n", "## a-1\n", "# a 1\n", "# A\n", "# b\n", "# \n", "# !\n", "# heading\n", "# a-2\n", "# a-1-1\n", "#  日本語\n", "# x {#a}\n", "\n",
		"$", "$$", "$`", "`$", " ", "x", "1", "\\", "\\$", "$x$", "$$x$$", "$`x`$", "a $5 and $10", "`", "*", "[", "](#a)", "|", "| a | $`b | c`$ |\n", "|-|-|-|\n", "> ", "- ", "<b>", "\n\n"}
	rnd := rand.New(rand.NewSource(3))
	for round := range 3000 {
		var b strings.Builder
		for range 1 + rnd.Intn(30) {
			b.WriteString(pieces[rnd.Intn(len(pieces))])
		}
		got, err := ToHTML([]byte(b.String()))
		if err != nil {
			t.Fatal(err)
		}
		want, err := toHTMLBefore([]byte(b.String()))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("round %d: %q renders to\n%s\nwant\n%s", round, b.String(), got, want)
		}
	}
}

// fastest returns the shortest of three runs of fn.
func fastest(fn func()) time.Duration {
	best := time.Duration(1 << 62)
	for range 3 {
		start := time.Now()
		fn()
		best = min(best, time.Since(start))
	}
	return best
}

// growsWithSize fails the test when fn takes much more than four times as
// long for a size four times as large: the time grows with the square of
// the size (sixteen times), not with the size.
func growsWithSize(t *testing.T, size int, fn func(size int)) {
	t.Helper()
	small := fastest(func() { fn(size) })
	large := fastest(func() { fn(4 * size) })
	if large > 8*small+20*time.Millisecond {
		t.Errorf("size %d took %v, size %d took %v", size, small, 4*size, large)
	}
}

// Headings that repeat get their numbers in a time that grows with their
// number: the numbers given are not tried again for each of them.
func TestRepeatedHeadings(t *testing.T) {
	growsWithSize(t, 3000, func(size int) {
		g := &githubIDs{used: map[string]bool{}, next: map[string]int{}}
		var last []byte
		for range size {
			last = g.Generate([]byte("Same Heading"), ast.KindHeading)
		}
		if want := fmt.Sprintf("same-heading-%d", size-1); string(last) != want {
			t.Errorf("the last of %d headings has the id %s, want %s", size, last, want)
		}
	})
}

// A line in which many formulas in backticks start and none ends is read
// in a time that grows with its length: the end is looked for once.
func TestUnclosedFormulas(t *testing.T) {
	growsWithSize(t, 20000, func(size int) {
		out, err := ToHTML([]byte(strings.Repeat("$`x ", size) + "\n\n$`y`$ " + strings.Repeat("$`x ", size) + "\n"))
		if err != nil {
			t.Fatal(err)
		}
		if n := bytes.Count(out, []byte("$")); n != 2*size || bytes.Count(out, []byte("<math")) != 1 {
			t.Errorf("%d dollar signs of %d, %d formulas", n, 2*size, bytes.Count(out, []byte("<math")))
		}
	})
}
