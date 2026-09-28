package drawio

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf"
)

// Tests that malformed or hostile input ends in an error or a warning
// instead of a hang, a stack overflow or unbounded memory. Each uses the
// test fonts so results do not depend on the installed fonts.

func limitsModel(cells string) []byte {
	return []byte(`<mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>` + cells + `</root></mxGraphModel>`)
}

func limitsOptions() *Options {
	return &Options{FontDirs: testFonts, NoSystemFonts: true, NoTextIndex: true}
}

// convertWithin converts data and fails if it does not finish in limit
// (a hang). It returns the result of the conversion.
func convertWithin(t *testing.T, data []byte, limit time.Duration) *Result {
	t.Helper()
	type outcome struct {
		res *Result
		err error
	}
	ch := make(chan outcome, 1)
	go func() {
		res, err := Convert(data, limitsOptions())
		ch <- outcome{res, err}
	}()
	select {
	case o := <-ch:
		if o.err != nil {
			t.Fatalf("Convert: %v", o.err)
		}
		return o.res
	case <-time.After(limit):
		t.Fatalf("Convert did not finish in %v (hang)", limit)
		return nil
	}
}

// deflate encodes s the way draw.io compresses a diagram or stencil.
func deflate(s string) string {
	var b bytes.Buffer
	w, _ := flate.NewWriter(&b, flate.BestCompression)
	w.Write([]byte(url.PathEscape(s)))
	w.Close()
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// deflateRaw encodes raw bytes with no URL encoding (for a bomb payload).
func deflateRaw(n int) string {
	var b bytes.Buffer
	w, _ := flate.NewWriter(&b, flate.BestCompression)
	chunk := bytes.Repeat([]byte{'A'}, 1<<16)
	for n > 0 {
		k := min(n, len(chunk))
		w.Write(chunk[:k])
		n -= k
	}
	w.Close()
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// A cell that is its own grandparent (a parent cycle) with an edge into it
// once hung visibleTerminal; the model now cuts the unreachable chain.
func TestParentCycleTerminatesEdge(t *testing.T) {
	data := limitsModel(
		`<mxCell id="a" parent="b" vertex="1"><mxGeometry width="10" height="10" as="geometry"/></mxCell>` +
			`<mxCell id="b" parent="a" vertex="1"><mxGeometry width="10" height="10" as="geometry"/></mxCell>` +
			`<mxCell id="e" parent="1" edge="1" source="a"><mxGeometry relative="1" as="geometry"><mxPoint x="50" y="50" as="targetPoint"/></mxGeometry></mxCell>`)
	convertWithin(t, data, 10*time.Second)
}

// A swimlane color that points into a parent cycle once hung swimlaneOf.
func TestParentCycleTerminatesSwimlane(t *testing.T) {
	data := limitsModel(
		`<mxCell id="a" parent="b" vertex="1"><mxGeometry width="10" height="10" as="geometry"/></mxCell>` +
			`<mxCell id="b" parent="a" vertex="1"><mxGeometry width="10" height="10" as="geometry"/></mxCell>` +
			`<mxCell id="v" parent="1" vertex="1" target="a" style="fillColor=swimlane"><mxGeometry width="10" height="10" as="geometry"/></mxCell>`)
	convertWithin(t, data, 10*time.Second)
}

// Two swimlanes whose colors refer to each other once recursed until the
// stack overflowed; resolveColor now stops at maxColorDepth.
func TestSwimlaneColorCycle(t *testing.T) {
	data := limitsModel(
		`<mxCell id="s1" parent="1" vertex="1" target="s2" style="swimlane;fillColor=swimlane"><mxGeometry width="100" height="100" as="geometry"/></mxCell>` +
			`<mxCell id="s2" parent="1" vertex="1" target="s1" style="swimlane;fillColor=swimlane"><mxGeometry x="200" width="100" height="100" as="geometry"/></mxCell>`)
	convertWithin(t, data, 10*time.Second)
}

// A compressed diagram body larger than the limit is rejected, not read
// into memory whole.
func TestDecompressionBombDiagram(t *testing.T) {
	old := maxDecompressed
	maxDecompressed = 1 << 20
	defer func() { maxDecompressed = old }()
	data := []byte(`<mxfile><diagram id="x" name="p">` + deflateRaw(4<<20) + `</diagram></mxfile>`)
	if _, err := Convert(data, limitsOptions()); err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("want a size-limit error, got %v", err)
	}
}

// A shape=stencil(...) style larger than the limit is rejected.
func TestDecompressionBombStencil(t *testing.T) {
	old := maxDecompressed
	maxDecompressed = 1 << 20
	defer func() { maxDecompressed = old }()
	data := limitsModel(`<mxCell id="v" parent="1" vertex="1" style="shape=stencil(` + deflateRaw(4<<20) + `)"><mxGeometry width="10" height="10" as="geometry"/></mxCell>`)
	res := convertWithin(t, data, 10*time.Second)
	// the stencil is dropped with a warning; the cell is still drawn
	if !hasWarning(res, "cannot be read") {
		t.Fatalf("want a warning that the stencil could not be read, got %v", res.Warnings)
	}
}

// A stencil whose elements nest past the limit is truncated with a warning
// instead of overflowing the stack.
func TestStencilNestingLimit(t *testing.T) {
	n := maxStencilNest + 50
	shape := `<shape name="deep" w="100" h="100"><foreground>` +
		strings.Repeat("<path>", n) + strings.Repeat("</path>", n) + `</foreground></shape>`
	data := limitsModel(`<mxCell id="v" parent="1" vertex="1" style="shape=stencil(` + deflate(shape) + `)"><mxGeometry width="100" height="100" as="geometry"/></mxCell>`)
	res := convertWithin(t, data, 10*time.Second)
	if !hasWarning(res, "deeply nested stencil") {
		t.Fatalf("want a deep-stencil warning, got %v", res.Warnings)
	}
}

// A label whose markup nests past the limit is truncated with a warning
// instead of overflowing the stack.
func TestLabelNestingLimit(t *testing.T) {
	n := maxLabelDepth + 50
	// draw.io stores the label markup XML-escaped in the value attribute
	label := strings.Repeat("&lt;b&gt;", n) + "x"
	data := limitsModel(`<mxCell id="v" parent="1" vertex="1" value="` + label + `" style="html=1"><mxGeometry width="100" height="100" as="geometry"/></mxCell>`)
	res := convertWithin(t, data, 10*time.Second)
	if !hasWarning(res, "deeply nested label") {
		t.Fatalf("want a deep-label warning, got %v", res.Warnings)
	}
}

// A formula longer than the limit is drawn as text with a warning instead
// of being handed to the recursive formula engine.
func TestFormulaLengthLimit(t *testing.T) {
	tex := strings.Repeat("{", maxFormulaLen) + "x" + strings.Repeat("}", maxFormulaLen)
	data := []byte(`<mxGraphModel math="1"><root><mxCell id="0"/><mxCell id="1" parent="0"/>` +
		`<mxCell id="v" parent="1" vertex="1" value="$$` + tex + `$$" style="html=1"><mxGeometry width="100" height="100" as="geometry"/></mxCell></root></mxGraphModel>`)
	res := convertWithin(t, data, 10*time.Second)
	if !hasWarning(res, "too long to typeset") {
		t.Fatalf("want a formula-length warning, got %v", res.Warnings)
	}
}

// indexCloseTag finds the case-insensitive end tag without lowercasing the
// whole input (the fix for the quadratic scan of script and style tags).
func TestIndexCloseTag(t *testing.T) {
	cases := []struct {
		s, tag string
		want   int
	}{
		{"abc</SCRIPT>", "</script", 3},
		{"a<b></Style>c", "</style", 4},
		{"no end here", "</script", -1},
		{"</script", "</script", 0},
	}
	for _, c := range cases {
		if got := indexCloseTag(c.s, c.tag); got != c.want {
			t.Errorf("indexCloseTag(%q, %q) = %d, want %d", c.s, c.tag, got, c.want)
		}
	}
	// many style elements parse in reasonable time (was quadratic before)
	label := strings.Repeat("<style>A</style>", 40000)
	done := make(chan struct{})
	go func() { parseHTML(label); close(done) }()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatalf("parseHTML of %d style elements did not finish in 10s", 40000)
	}
}

func hasWarning(res *Result, substr string) bool {
	for _, w := range res.Warnings {
		if strings.Contains(w, substr) {
			return true
		}
	}
	return false
}

// TestPlaceholderLimit converts a label made of placeholders of a long
// attribute: what they stand for is not written out beyond the limit.
func TestPlaceholderLimit(t *testing.T) {
	const n = 2000 // 2000 placeholders of 2000 bytes each: 4 MB of text
	cells := `<object id="2" placeholders="1" a="` + strings.Repeat("x", n) + `" label="` + strings.Repeat("%a%", n) + `">` +
		`<mxCell style="whiteSpace=wrap" vertex="1" parent="1"><mxGeometry width="120" height="60" as="geometry"/></mxCell></object>`
	m := parseModel(mustXML(t, limitsModel(cells)))
	got, cut := m.label(m.cells["2"], nil, nil)
	if !cut || len(got) > maxLabel+3*n {
		t.Errorf("a label of %d bytes, cut %v", len(got), cut)
	}
	if !strings.HasPrefix(got, "xxx") || !strings.HasSuffix(got, "%a%") {
		t.Errorf("the label starts %q and ends %q", got[:3], got[len(got)-3:])
	}
	res := convertWithin(t, limitsModel(cells), 20*time.Second)
	if !hasWarning(res, "placeholders") {
		t.Errorf("warnings %q", res.Warnings)
	}
	// a label within the limit is whole (%c% stands for nothing, %% for %)
	short := `<object id="2" placeholders="1" a="1" b="2" label="%a%+%b%=%c%%"><mxCell vertex="1" parent="1"/></object>`
	m = parseModel(mustXML(t, limitsModel(short)))
	if got, cut := m.label(m.cells["2"], nil, nil); got != "1+2=%c%" || cut {
		t.Errorf("label %q, cut %v", got, cut)
	}
}

func mustXML(t *testing.T, data []byte) *xmlNode {
	t.Helper()
	n, err := parseXML(data)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

// TestPatternBudget converts a page of pattern fills of many lines each:
// the lines of the page have a limit, not only those of every shape.
func TestPatternBudget(t *testing.T) {
	defer func(n int) { maxPageRepeats = n }(maxPageRepeats)
	maxPageRepeats = 10_000
	var b strings.Builder
	for i := range 40 {
		// 900 lines each
		b.WriteString(`<mxCell id="p` + strings.Repeat("i", i+1) + `" style="shape=mxgraph.basic.patternFillRect;fillStyle=hor;step=1;" vertex="1" parent="1">` +
			`<mxGeometry width="100" height="900" as="geometry"/></mxCell>`)
	}
	res := convertWithin(t, limitsModel(b.String()), 20*time.Second)
	if !hasWarning(res, "lines of patterns") {
		t.Errorf("warnings %q", res.Warnings)
	}
	lines := 0
	for _, p := range res.Doc.Parts() {
		if p.Type != bdf.PartObject {
			continue
		}
		o, err := bdf.DecodeObject(p.Data)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range o.Paths {
			if e.Inline == nil {
				continue
			}
			for _, v := range e.Inline.Verbs {
				if v == bdf.VerbMove {
					lines++
				}
			}
		}
	}
	// 40 shapes of 900 lines are 36,000 lines
	if lines < 9_000 || lines > 12_000 {
		t.Errorf("%d lines drawn", lines)
	}
}
