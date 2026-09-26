package drawio

import (
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// labelOf lays out the label of the cell with the given id in a one-page model.
func labelOf(t *testing.T, cells, id string) (*converter, *labelBox) {
	t.Helper()
	root, err := parseXML([]byte(`<mxGraphModel><root><mxCell id="0"/><mxCell id="1" parent="0"/>` + cells + `</root></mxGraphModel>`))
	if err != nil {
		t.Fatal(err)
	}
	c := newConverter(testOptions())
	c.m = parseModel(root)
	c.page = &pageInfo{name: "P", number: 1, count: 1}
	v := newView(c.m, nil, nil)
	st := v.state(c.m.cells[id])
	return c, c.layoutLabel(st, c.newShape(st))
}

func lineTexts(l *labelBox) []string {
	var out []string
	for _, ln := range l.lo.lines {
		var b strings.Builder
		for _, it := range ln.items {
			b.WriteRune(it.r)
		}
		out = append(out, b.String())
	}
	return out
}

func TestLabelWrap(t *testing.T) {
	_, l := labelOf(t, `<mxCell id="a" value="A long label that should wrap inside the box nicely" style="whiteSpace=wrap;html=1;" vertex="1" parent="1"><mxGeometry x="0" y="0" width="120" height="60" as="geometry"/></mxCell>`, "a")
	lines := lineTexts(l)
	if len(lines) < 3 {
		t.Fatalf("lines %q", lines)
	}
	// wrapping width: the box less the spacing, plus the foreignObject padding
	for i, ln := range l.lo.lines {
		if ln.width > 120-4+2 {
			t.Errorf("line %d %q is %.1f wide", i, lines[i], ln.width)
		}
		if strings.HasPrefix(lines[i], " ") || strings.HasSuffix(lines[i], " ") {
			t.Errorf("line %q keeps a collapsible space", lines[i])
		}
	}
	if strings.Join(lines, " ") != "A long label that should wrap inside the box nicely" {
		t.Errorf("lines %q", lines)
	}
	// centered in the box: the middle of the text is the middle of the cell
	mid := l.top + l.lo.height/2
	if mid < 28 || mid > 32 {
		t.Errorf("text middle at %.1f, want about 30", mid)
	}

	_, l = labelOf(t, `<mxCell id="a" value="A long label that does not wrap" style="html=1;" vertex="1" parent="1"><mxGeometry x="0" y="0" width="60" height="40" as="geometry"/></mxCell>`, "a")
	if n := len(l.lo.lines); n != 1 || l.lo.width < 100 {
		t.Errorf("unwrapped label: %d lines, %.1f wide", n, l.lo.width)
	}
}

func TestLabelKinsoku(t *testing.T) {
	_, l := labelOf(t, `<mxCell id="a" value="日本語のラベルは文字ごとに折り返されます。句読点は行頭に来ません。" style="whiteSpace=wrap;html=1;" vertex="1" parent="1"><mxGeometry x="0" y="0" width="100" height="80" as="geometry"/></mxCell>`, "a")
	lines := lineTexts(l)
	if len(lines) < 3 {
		t.Fatalf("lines %q", lines)
	}
	for i, s := range lines {
		if i > 0 && (strings.HasPrefix(s, "。") || strings.HasPrefix(s, "、")) {
			t.Errorf("line %d starts with closing punctuation: %q", i, lines)
		}
		if i > 0 && !l.lo.lines[i].cjkWrap {
			t.Errorf("line %d not marked as an East Asian wrap", i)
		}
	}
}

func TestPlainLabel(t *testing.T) {
	_, l := labelOf(t, `<mxCell id="a" value="Plain&#xa;two   lines" style="rounded=0;" vertex="1" parent="1"><mxGeometry x="0" y="0" width="100" height="60" as="geometry"/></mxCell>`, "a")
	if got := lineTexts(l); len(got) != 2 || got[1] != "two lines" || !l.plain {
		t.Fatalf("lines %q plain %v", got, l.plain)
	}
}

func TestHTMLStructure(t *testing.T) {
	c, l := labelOf(t, `<mxCell id="a" value="&lt;h2&gt;Title&lt;/h2&gt;one&lt;br&gt;two&lt;ul&gt;&lt;li&gt;x&lt;/li&gt;&lt;li&gt;y&lt;/li&gt;&lt;/ul&gt;" style="text;html=1;" vertex="1" parent="1"><mxGeometry x="0" y="0" width="200" height="100" as="geometry"/></mxCell>`, "a")
	cv := c.objs.New()
	c.drawLabel(newC2D(cv), l)
	c.finalize()
	o, err := bdf.DecodeObject(c.doc.Part(cv.Hash()).Data)
	if err != nil {
		t.Fatal(err)
	}
	var marks []string
	o.Walk(func(in bdf.Instr) {
		switch in.Op {
		case bdf.OpMark:
			marks = append(marks, map[uint64]string{0: "P", 1: "LINE", 3: "BOX", 5: "WRAP", 6: "H", 7: "LIST", 8: "LI", 11: "END"}[in.Args[0].(uint64)])
		case bdf.OpFillText:
			marks = append(marks, in.Args[0].(string))
		}
	})
	got := strings.Join(marks, " ")
	if want := "BOX H Title P one P two LIST LI x LI y END"; got != want {
		t.Errorf("marks\n got %s\nwant %s", got, want)
	}
}
