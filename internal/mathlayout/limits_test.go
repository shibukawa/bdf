package mathlayout

import (
	"strings"
	"testing"

	"github.com/shibukawa/bdf/internal/xmltree"
	"golang.org/x/net/html"
)

// depthOf returns how deeply the nodes of a formula nest.
func depthOf(n Node) int {
	deepest := 0
	kids := func(list ...Node) {
		for _, k := range list {
			if k != nil {
				deepest = max(deepest, depthOf(k))
			}
		}
	}
	switch n := n.(type) {
	case *Row:
		kids(n.Kids...)
	case *Frac:
		kids(n.Num, n.Den)
	case *Radical:
		kids(n.Base, n.Degree)
	case *Scripts:
		kids(n.Base, n.Sub, n.Sup, n.PreSub, n.PreSup)
	case *UnderOver:
		kids(n.Base, n.Under, n.Over)
	case *Table:
		for _, r := range n.Rows {
			kids(r...)
		}
	case *Styled:
		kids(n.Kid)
	case *Enclose:
		kids(n.Kid)
	case *Phantom:
		kids(n.Kid)
	case *Bar:
		kids(n.Kid)
	}
	return deepest + 1
}

// Groups and arguments nest 200 deep: what is nested deeper is read as a
// part of what it is in, so that reading and laying out a formula need
// neither a stack as deep as the formula is long, nor the time of its
// depth for each of its levels.
func TestTeXDepth(t *testing.T) {
	const levels = 3000
	for _, c := range []struct{ name, open, close string }{
		{"group", "{", "}"},
		{"frac", `\frac1{`, "}"},
		{"sqrt", `\sqrt `, ""},
		{"sqrt degree", `\sqrt[{`, "}]2"},
		{"over", `1\over `, ""},
		{"left", `\left(`, `\right)`},
		{"matrix", `\begin{pmatrix}1&`, `\end{pmatrix}`},
		{"root", `\root 3 \of `, ""},
		{"style", `{\displaystyle `, "}"},
		{"color", `\color{red} `, ""},
		{"font", `\bf `, ""},
		{"script", `x^{`, "}"},
		{"hat", `\hat `, ""},
		{"arrow", `\xrightarrow[{`, `}]{}`},
	} {
		// a level makes up to four nodes, one in the other (the arrow)
		n := ParseTeX(strings.Repeat(c.open, levels) + "y" + strings.Repeat(c.close, levels))
		if d := depthOf(n); d > 2*200+10 {
			t.Errorf("%s: %d levels make a formula %d deep", c.name, levels, d)
		}
		if s := Linear(n); !strings.Contains(s, "y") {
			t.Errorf("%s: the innermost is lost: %.80s", c.name, s)
		}
		// less deep than it may be, all is read as it is written
		n = ParseTeX(strings.Repeat(c.open, 40) + "y" + strings.Repeat(c.close, 40))
		if d := depthOf(n); d < 41 {
			t.Errorf("%s: 40 levels make a formula %d deep", c.name, d)
		}
	}
	// \not takes what follows it, which may be \not
	if s := Linear(ParseTeX(strings.Repeat(`\not `, levels) + "=")); s != "=" {
		t.Errorf("not: %q", s)
	}
	if s := Linear(ParseTeX(strings.Repeat(`\not `, 3) + "=")); s != "≠\u0338\u0338" {
		t.Errorf("not: %q", s)
	}
	// the braces of the groups read over close nothing else
	n := ParseTeX(strings.Repeat("{", 250) + "a" + strings.Repeat("}", 250) + "b}c")
	if s := Linear(n); s != "ab" {
		t.Errorf("groups 250 deep, then b, a brace and c: %q", s)
	}
}

// The elements of MathML nest 200 deep; of those nested deeper the tokens
// are read.
func TestMathMLDepth(t *testing.T) {
	const levels = 400
	src := `<math>` + strings.Repeat(`<msqrt><mn>1</mn>`, levels) + `<mfrac><mi>y</mi><mrow><mi>z</mi><annotation>no</annotation></mrow></mfrac>` +
		strings.Repeat(`</msqrt>`, levels) + `</math>`
	doc, err := html.Parse(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	n, _ := ParseMathML(findElement(doc, "math"))
	if d := depthOf(n); d > 2*200+4 || d < 200 {
		t.Errorf("%d levels make a formula %d deep", levels, d)
	}
	s := Linear(n)
	if !strings.HasSuffix(s, strings.Repeat("1", levels-199)+"yz"+strings.Repeat(")", 200)) || strings.Count(s, "1") != levels || strings.Contains(s, "no") {
		t.Errorf("formula %.60s … %s", s, s[max(len(s)-300, 0):])
	}
}

// The arguments of Office Math nest 200 deep; of those nested deeper the
// runs are read.
func TestOMMLDepth(t *testing.T) {
	const levels = 400
	src := `<m:oMath ` + mns + `>` + strings.Repeat(`<m:rad><m:radPr><m:degHide m:val="1"/></m:radPr><m:deg/><m:e><m:r><m:t>1</m:t></m:r>`, levels) +
		`<m:f><m:num><m:r><m:t>y</m:t></m:r></m:num><m:den><m:r><m:t>z</m:t></m:r></m:den></m:f>` + strings.Repeat(`</m:e></m:rad>`, levels) + `</m:oMath>`
	root, err := xmltree.Parse([]byte(src))
	if err != nil {
		t.Skipf("the document is not read: %v", err)
	}
	o := &OMML{}
	n := o.ParseOMML(root)
	if d := depthOf(n); d > 2*200+4 || d < 200 {
		t.Errorf("%d levels make a formula %d deep", levels, d)
	}
	if s := Linear(n); strings.Count(s, "1") != levels || !strings.Contains(s, "yz") {
		t.Errorf("formula %.60s … %s", s, s[max(len(s)-300, 0):])
	}
	if o.depth != 0 {
		t.Errorf("depth %d after the formula", o.depth)
	}
}
