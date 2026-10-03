package mathlayout

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Linear writes a formula in a linear notation close to the one Office
// calls linear format (UnicodeMath): x^2, x_(i+1), (a+b)/(c+d), √(x+1),
// √(n&x), ∑_(i=1)^n x_i. It is the text that extraction, search and copy
// see for the formula.
func Linear(n Node) string {
	w := &linearWriter{}
	w.node(n)
	return strings.TrimSpace(w.b.String())
}

// linearWriter writes the linear notation of a formula.
type linearWriter struct {
	b       strings.Builder
	variant Variant // of the styles the node written is in
}

// sub writes n on its own and returns its text.
func (w *linearWriter) sub(n Node) string {
	s := &linearWriter{variant: w.variant}
	s.node(n)
	return s.b.String()
}

func (w *linearWriter) node(n Node) {
	b := &w.b
	switch n := n.(type) {
	case *Row:
		if n == nil {
			return
		}
		for i, k := range n.Kids {
			t := w.sub(k)
			if prev := n.Kids[max(i-1, 0)]; i > 0 && isFunctionNode(prev) && (startsAlnum(t) || hasScripts(prev)) {
				// sin θ, det A
				b.WriteString(" ")
			}
			b.WriteString(t)
		}
	case *Atom:
		if n == nil {
			return
		}
		b.WriteString(atomText(n, w.variant))
	case *Frac:
		if n.Kind == FracLinear || n.Kind == FracSkewed || !n.NoBar {
			w.group(n.Num)
			b.WriteString("/")
			w.group(n.Den)
			return
		}
		b.WriteString("(")
		w.node(n.Num)
		b.WriteString("¦")
		w.node(n.Den)
		b.WriteString(")")
	case *Radical:
		if n.Degree == nil {
			b.WriteString("√")
			w.group(n.Base)
			return
		}
		b.WriteString("√(")
		w.node(n.Degree)
		b.WriteString("&")
		w.node(n.Base)
		b.WriteString(")")
	case *Scripts:
		if n.PreSub != nil || n.PreSup != nil {
			b.WriteString(" ")
			w.scripts(n.PreSub, n.PreSup)
		}
		w.group(n.Base)
		w.scripts(n.Sub, n.Sup)
		if isLargeOpNode(n.Base) {
			b.WriteString(" ")
		}
	case *UnderOver:
		if a := coreAtom(n.Over); a != nil && a.Accent && n.Under == nil && isAccentChar(a.Text) {
			// an accent: the base and the combining form of the accent
			w.group(n.Base)
			b.WriteRune(combiningAccent([]rune(normalizeOp(a.Text))[0]))
			return
		}
		if a := coreAtom(n.Over); a != nil && a.Stretchy && n.Under == nil {
			b.WriteString(atomText(a, w.variant))
			w.group(n.Base)
			return
		}
		if a := coreAtom(n.Under); a != nil && a.Stretchy && n.Over == nil {
			b.WriteString(atomText(a, w.variant))
			w.group(n.Base)
			return
		}
		w.group(n.Base)
		w.scripts(n.Under, n.Over)
		if isLargeOpNode(n.Base) {
			b.WriteString(" ")
		}
	case *Table:
		for i, row := range n.Rows {
			if i > 0 {
				if n.Aligned || n.Display {
					b.WriteString("\n")
				} else {
					b.WriteString("; ")
				}
			}
			for j, cell := range row {
				if j > 0 && !n.Aligned {
					b.WriteString(", ")
				}
				w.node(cell)
			}
		}
	case *Space:
		if n.Width >= 0.15 {
			b.WriteString(" ")
		}
	case *Styled:
		if n.Variant != VarAuto {
			saved := w.variant
			w.variant = n.Variant
			defer func() { w.variant = saved }()
		}
		w.node(n.Kid)
	case *Enclose:
		w.node(n.Kid)
	case *Phantom:
		if n.Show {
			w.node(n.Kid)
		}
	case *Bar:
		if n.Under {
			b.WriteString("▁")
		} else {
			b.WriteString("¯")
		}
		w.group(n.Kid)
	}
}

// atomText is the text of a token: letters in their plain form unless their
// style changes what they mean (ℝ, 𝔤).
func atomText(a *Atom, inherited Variant) string {
	text := a.Text
	v := a.Variant
	if v == VarAuto {
		v = inherited
	}
	if a.Kind == Op {
		text = normalizeOp(text)
		if r, ok := singleRune(text); ok && r >= 0x2061 && r <= 0x2064 {
			return ""
		}
	}
	switch v {
	case VarDoubleStruck, VarScript, VarBoldScript, VarFraktur, VarBoldFraktur:
		var s strings.Builder
		for _, r := range text {
			s.WriteRune(styled(r, v))
		}
		return s.String()
	}
	if a.Kind == Text {
		return text
	}
	return text
}

func (w *linearWriter) scripts(sub, sup Node) {
	b := &w.b
	if sub != nil {
		b.WriteString("_")
		w.group(sub)
	}
	if sup != nil {
		b.WriteString("^")
		w.group(sup)
	}
}

// group writes a node, in parentheses when it is more than one token.
func (w *linearWriter) group(n Node) {
	b := &w.b
	t := w.sub(n)
	if simple(n) || utf8.RuneCountInString(t) <= 1 || isFenced(t) {
		b.WriteString(t)
		return
	}
	b.WriteString("(")
	b.WriteString(t)
	b.WriteString(")")
}

// simple reports whether a node is one token (or a token with scripts).
func simple(n Node) bool {
	switch n := n.(type) {
	case *Atom:
		return n != nil && n.Kind != Text && !strings.ContainsRune(n.Text, ' ')
	case *Styled:
		return simple(n.Kid)
	case *Row:
		return n != nil && len(n.Kids) == 1 && simple(n.Kids[0])
	case *Radical:
		return n.Degree == nil
	}
	return false
}

// isFenced reports whether s is wrapped in one pair of brackets.
func isFenced(s string) bool {
	if len(s) < 2 {
		return false
	}
	pairs := map[rune]rune{'(': ')', '[': ']', '{': '}', '⟨': '⟩', '|': '|', '‖': '‖', '⌊': '⌋', '⌈': '⌉'}
	rs := []rune(s)
	cl, ok := pairs[rs[0]]
	if !ok || rs[len(rs)-1] != cl {
		return false
	}
	depth := 0
	for i, r := range rs {
		switch {
		case r == rs[0] && (cl != r || i == 0):
			depth++
		case r == cl:
			depth--
			if depth == 0 && i != len(rs)-1 {
				return false
			}
		}
	}
	return true
}

func isLargeOpNode(n Node) bool {
	a := coreAtom(n)
	return a != nil && a.LargeOp
}

// hasScripts reports whether a node has limits or scripts (lim_(x→0)).
func hasScripts(n Node) bool {
	switch n := n.(type) {
	case *Scripts, *UnderOver:
		return true
	case *Row:
		return n != nil && len(n.Kids) == 1 && hasScripts(n.Kids[0])
	case *Styled:
		return hasScripts(n.Kid)
	}
	return false
}

func isAccentChar(s string) bool {
	r, ok := singleRune(normalizeOp(s))
	return ok && (isCombining(r) || combiningAccent(r) != r)
}

// isFunctionNode reports whether a node is a function name (sin, det), or a
// function name with scripts (cos², log₂).
func isFunctionNode(n Node) bool {
	switch n := n.(type) {
	case *Row:
		if n != nil && n.Class == LargeOp && len(n.Kids) == 1 {
			return isFunctionNode(n.Kids[0]) || coreAtom(n.Kids[0]) != nil
		}
		return n != nil && len(n.Kids) == 1 && isFunctionNode(n.Kids[0])
	case *Scripts:
		return isFunctionNode(n.Base)
	case *UnderOver:
		return isFunctionNode(n.Base)
	case *Styled:
		return isFunctionNode(n.Kid)
	case *Atom:
		return n != nil && n.Kind == Ident && (n.Class == LargeOp || IsFunctionName(n.Text))
	}
	return false
}

func startsAlnum(s string) bool {
	for _, r := range s {
		return unicode.IsLetter(r) || unicode.IsDigit(r)
	}
	return false
}
