package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Math in Markdown is written as GitHub writes it: $…$ and $`…`$ in a line,
// $$…$$ for a formula set on its own (in a paragraph, or as lines of their
// own), and fenced code blocks of the language math. The formulas are
// rendered as MathML math elements that hold the LaTeX as an annotation,
// which the HTML converter lays out.

var (
	kindMathInline = ast.NewNodeKind("MathInline")
	kindMathBlock  = ast.NewNodeKind("MathBlock")
)

type mathInline struct {
	ast.BaseInline
	src     []byte
	display bool
}

func (n *mathInline) Kind() ast.NodeKind { return kindMathInline }

func (n *mathInline) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Src": string(n.src)}, nil)
}

type mathBlock struct {
	ast.BaseBlock
	src []byte
}

func (n *mathBlock) Kind() ast.NodeKind { return kindMathBlock }

func (n *mathBlock) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Src": string(n.src)}, nil)
}

// mathExtension adds the math syntax to a goldmark.Markdown.
type mathExtension struct{}

func (mathExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(util.Prioritized(mathBlockParser{}, 690)),
		parser.WithInlineParsers(util.Prioritized(mathInlineParser{}, 90)),
		parser.WithASTTransformers(util.Prioritized(mathFences{}, 100)),
	)
	m.Renderer().AddOptions(renderer.WithNodeRenderers(util.Prioritized(mathRenderer{}, 100)))
}

type mathInlineParser struct{}

func (mathInlineParser) Trigger() []byte { return []byte{'$'} }

func (mathInlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
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
		// the first dollar sign that follows closes the formula; when it
		// cannot (a space before it, a digit after it: "$5 and $10"),
		// there is no formula
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

func isSpaceByte(b byte) bool { return b == ' ' || b == '\t' || b == '\n' || b == '\r' }

// mathBlockParser reads formulas between lines that start and end with $$.
type mathBlockParser struct{}

func (mathBlockParser) Trigger() []byte { return []byte{'$'} }

func (mathBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	pos := pc.BlockOffset()
	if pos < 0 || !bytes.HasPrefix(line[pos:], []byte("$$")) {
		return nil, parser.NoChildren
	}
	rest := line[pos+2:]
	if end := bytes.Index(rest, []byte("$$")); end >= 0 {
		// one line: a block only when nothing follows the formula
		if len(bytes.TrimSpace(rest[end+2:])) > 0 {
			return nil, parser.NoChildren
		}
		reader.AdvanceToEOL()
		return &mathBlock{src: append([]byte(nil), rest[:end]...)}, parser.Close
	}
	reader.AdvanceToEOL()
	return &mathBlock{src: append([]byte(nil), rest...)}, parser.NoChildren
}

func (mathBlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, _ := reader.PeekLine()
	if line == nil {
		return parser.Close
	}
	m := node.(*mathBlock)
	if end := bytes.Index(line, []byte("$$")); end >= 0 {
		m.src = append(m.src, line[:end]...)
		reader.AdvanceToEOL()
		return parser.Close
	}
	m.src = append(m.src, line...)
	reader.AdvanceToEOL()
	return parser.Continue | parser.NoChildren
}

func (mathBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {}
func (mathBlockParser) CanInterruptParagraph() bool                                { return true }
func (mathBlockParser) CanAcceptIndentedLine() bool                                { return false }

// mathFences turns fenced code blocks of the language math into formulas.
type mathFences struct{}

func (mathFences) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	src := reader.Source()
	var fences []*ast.FencedCodeBlock
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if f, ok := n.(*ast.FencedCodeBlock); ok && entering && string(f.Language(src)) == "math" {
			fences = append(fences, f)
		}
		return ast.WalkContinue, nil
	})
	for _, f := range fences {
		var b bytes.Buffer
		for i := 0; i < f.Lines().Len(); i++ {
			seg := f.Lines().At(i)
			b.Write(seg.Value(src))
		}
		f.Parent().ReplaceChild(f.Parent(), f, &mathBlock{src: b.Bytes()})
	}
}

type mathRenderer struct{}

func (mathRenderer) RegisterFuncs(r renderer.NodeRendererFuncRegisterer) {
	r.Register(kindMathInline, func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			m := n.(*mathInline)
			writeMath(w, m.src, m.display)
		}
		return ast.WalkSkipChildren, nil
	})
	r.Register(kindMathBlock, func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			writeMath(w, n.(*mathBlock).src, true)
			w.WriteString("\n")
		}
		return ast.WalkSkipChildren, nil
	})
}

// writeMath writes a formula as a math element holding its LaTeX.
func writeMath(w util.BufWriter, src []byte, display bool) {
	w.WriteString("<math")
	if display {
		w.WriteString(` display="block"`)
	}
	w.WriteString(`><semantics><annotation encoding="application/x-tex">`)
	w.Write(util.EscapeHTML(bytes.TrimSpace(src)))
	w.WriteString("</annotation></semantics></math>")
}
