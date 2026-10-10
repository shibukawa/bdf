// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htmlro

// The tree construction stage, after the parser of golang.org/x/net/html
// (parse.go), which this file follows closely: the same insertion modes,
// the same functions, so that the two build the same tree. What differs is
// what reaches it: tokens read in place by a Reader, with a tag's name and
// attributes turned into strings once, into the node that keeps them, and
// text gathered into one buffer per text node.

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"golang.org/x/net/html"
	a "golang.org/x/net/html/atom"
)

// token is the token the parser is working on. Tag names, comments and
// doctypes are strings; the text of a Text token is Text, bytes that are
// valid until the next token.
type token struct {
	Type     Kind
	DataAtom a.Atom
	Data     string
	Text     []byte
	Attr     []html.Attribute
}

// A parser implements the HTML5 parsing algorithm:
// https://html.spec.whatwg.org/multipage/syntax.html#tree-construction
type parser struct {
	// r provides the tokens for the parser.
	r *Reader
	// tok is the most recently read token.
	tok token
	// maxDepth bounds the stack of open elements; err is set when a token
	// takes the stack past it, and ends the parse.
	maxDepth int
	err      error
	// textNode is the text node being written: its bytes gather in textBuf
	// and become its Data when another node is written or the parse ends.
	textNode *html.Node
	textBuf  []byte
	// scratch holds the decoded text of the current Text token, and attr
	// values while they are decoded; scratch2 the text of foreign content
	// with its NULs replaced.
	scratch, scratch2 []byte
	// strs holds the strings made so far for the tag and attribute names
	// that are not atoms, so that each is made once.
	strs map[string]string
	// Self-closing tags like <hr/> are treated as start tags, except that
	// hasSelfClosingToken is set while they are being processed.
	hasSelfClosingToken bool
	// doc is the document root element.
	doc *html.Node
	// The stack of open elements (section 12.2.4.2) and active formatting
	// elements (section 12.2.4.3).
	oe, afe nodeStack
	// Element pointers (section 12.2.4.4).
	head, form *html.Node
	// Other parsing state flags (section 12.2.4.5).
	scripting, framesetOK bool
	// The stack of template insertion modes
	templateStack insertionModeStack
	// im is the current insertion mode.
	im insertionMode
	// originalIM is the insertion mode to go back to after completing a text
	// or inTableText insertion mode.
	originalIM insertionMode
	// fosterParenting is whether new elements should be inserted according to
	// the foster parenting rules (section 12.2.6.1).
	fosterParenting bool
	// quirks is whether the parser is operating in "quirks mode."
	quirks bool
	// fragment is whether the parser is parsing an HTML fragment.
	fragment bool
	// context is the context element when parsing an HTML fragment
	// (section 12.4).
	context *html.Node
}

func (p *parser) top() *html.Node {
	if n := p.oe.top(); n != nil {
		return n
	}
	return p.doc
}

// Stop tags for use in popUntil. These come from section 12.2.4.2.
var (
	defaultScopeStopTags = map[string][]a.Atom{
		"":     {a.Applet, a.Caption, a.Html, a.Table, a.Td, a.Th, a.Marquee, a.Object, a.Template, a.Select},
		"math": {a.AnnotationXml, a.Mi, a.Mn, a.Mo, a.Ms, a.Mtext},
		"svg":  {a.Desc, a.ForeignObject, a.Title},
	}
)

type scope int

const (
	defaultScope scope = iota
	listItemScope
	buttonScope
	tableScope
	tableRowScope
	tableBodyScope
)

// popUntil pops the stack of open elements at the highest element whose tag
// is in matchTags, provided there is no higher element in the scope's stop
// tags (as defined in section 12.2.4.2). It returns whether or not there was
// such an element. If there was not, popUntil leaves the stack unchanged.
//
// For example, the set of stop tags for table scope is: "html", "table". If
// the stack was:
// ["html", "body", "font", "table", "b", "i", "u"]
// then popUntil(tableScope, "font") would return false, but
// popUntil(tableScope, "i") would return true and the stack would become:
// ["html", "body", "font", "table", "b"]
//
// If an element's tag is in both the stop tags and matchTags, then the stack
// will be popped and the function returns true (provided, of course, there was
// no higher element in the stack that was also in the stop tags). For example,
// popUntil(tableScope, "table") returns true and leaves:
// ["html", "body", "font"]
func (p *parser) popUntil(s scope, matchTags ...a.Atom) bool {
	if i := p.indexOfElementInScope(s, matchTags...); i != -1 {
		p.oe = p.oe[:i]
		return true
	}
	return false
}

// indexOfElementInScope returns the index in p.oe of the highest element whose
// tag is in matchTags that is in scope. If no matching element is in scope, it
// returns -1.
func (p *parser) indexOfElementInScope(s scope, matchTags ...a.Atom) int {
	for i := len(p.oe) - 1; i >= 0; i-- {
		tagAtom := p.oe[i].DataAtom
		if p.oe[i].Namespace == "" {
			for _, t := range matchTags {
				if t == tagAtom {
					return i
				}
			}
			switch s {
			case defaultScope:
				// No-op.
			case listItemScope:
				if tagAtom == a.Ol || tagAtom == a.Ul {
					return -1
				}
			case buttonScope:
				if tagAtom == a.Button {
					return -1
				}
			case tableScope:
				if tagAtom == a.Html || tagAtom == a.Table || tagAtom == a.Template {
					return -1
				}
			default:
				panic(fmt.Sprintf("html: internal error: indexOfElementInScope unknown scope: %d", s))
			}
		}
		switch s {
		case defaultScope, listItemScope, buttonScope:
			for _, t := range defaultScopeStopTags[p.oe[i].Namespace] {
				if t == tagAtom {
					return -1
				}
			}
		}
	}
	return -1
}

// elementInScope is like popUntil, except that it doesn't modify the stack of
// open elements.
func (p *parser) elementInScope(s scope, matchTags ...a.Atom) bool {
	return p.indexOfElementInScope(s, matchTags...) != -1
}

// clearStackToContext pops elements off the stack of open elements until a
// scope-defined element is found.
func (p *parser) clearStackToContext(s scope) {
	for i := len(p.oe) - 1; i >= 0; i-- {
		tagAtom := p.oe[i].DataAtom
		switch s {
		case tableScope:
			if tagAtom == a.Html || tagAtom == a.Table || tagAtom == a.Template {
				p.oe = p.oe[:i+1]
				return
			}
		case tableRowScope:
			if tagAtom == a.Html || tagAtom == a.Tr || tagAtom == a.Template {
				p.oe = p.oe[:i+1]
				return
			}
		case tableBodyScope:
			if tagAtom == a.Html || tagAtom == a.Tbody || tagAtom == a.Tfoot || tagAtom == a.Thead || tagAtom == a.Template {
				p.oe = p.oe[:i+1]
				return
			}
		default:
			panic(fmt.Sprintf("html: internal error: clearStackToContext unknown scope: %d", s))
		}
	}
}

// parseGenericRawTextElement implements the generic raw text element parsing
// algorithm defined in 12.2.6.2.
// https://html.spec.whatwg.org/multipage/parsing.html#parsing-elements-that-contain-only-text
// TODO: Since both RAWTEXT and RCDATA states are treated as tokenizer's part
// officially, need to make tokenizer consider both states.
func (p *parser) parseGenericRawTextElement() {
	p.addElement()
	p.originalIM = p.im
	p.im = textIM
}

// generateImpliedEndTags pops nodes off the stack of open elements as long as
// the top node has a tag name of dd, dt, li, optgroup, option, p, rb, rp, rt or rtc.
// If exceptions are specified, nodes with that name will not be popped off.
func (p *parser) generateImpliedEndTags(exceptions ...string) {
	var i int
loop:
	for i = len(p.oe) - 1; i >= 0; i-- {
		n := p.oe[i]
		if n.Type != html.ElementNode {
			break
		}
		switch n.DataAtom {
		case a.Dd, a.Dt, a.Li, a.Optgroup, a.Option, a.P, a.Rb, a.Rp, a.Rt, a.Rtc:
			for _, except := range exceptions {
				if n.Data == except {
					break loop
				}
			}
			continue
		}
		break
	}

	p.oe = p.oe[:i+1]
}

// addChild adds a child node n to the top element, and pushes n onto the stack
// of open elements if it is an element node.
func (p *parser) addChild(n *html.Node) {
	if p.shouldFosterParent() {
		p.fosterParent(n, nil)
	} else {
		p.top().AppendChild(n)
	}

	if n.Type == html.ElementNode {
		p.insertOpenElement(n)
	}
}

func (p *parser) insertOpenElement(n *html.Node) {
	p.oe = append(p.oe, n)
	if len(p.oe) > p.maxDepth && p.err == nil {
		// Noted, not raised: the parse stops after this token, since a
		// panic is not recoverable on every compiler.
		p.err = ErrTooDeep
	}
}

// shouldFosterParent returns whether the next node to be added should be
// foster parented.
func (p *parser) shouldFosterParent() bool {
	if p.fosterParenting {
		switch p.top().DataAtom {
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			return true
		}
	}
	return false
}

// fosterParent adds a child node according to the foster parenting rules.
// Section 12.2.6.1, "foster parenting".
//
// When n is a text node, text is what it holds, and joins the text node
// before the table when there is one.
func (p *parser) fosterParent(n *html.Node, text []byte) {
	var table, parent, prev, template *html.Node
	var i int
	for i = len(p.oe) - 1; i >= 0; i-- {
		if p.oe[i].DataAtom == a.Table {
			table = p.oe[i]
			break
		}
	}

	var j int
	for j = len(p.oe) - 1; j >= 0; j-- {
		if p.oe[j].DataAtom == a.Template {
			template = p.oe[j]
			break
		}
	}

	if template != nil && (table == nil || j > i) {
		template.AppendChild(n)
		return
	}

	if table == nil {
		// The foster parent is the html element.
		parent = p.oe[0]
	} else {
		parent = table.Parent
	}
	if parent == nil {
		parent = p.oe[i-1]
	}

	if table != nil {
		prev = table.PrevSibling
	} else {
		prev = parent.LastChild
	}
	if prev != nil && prev.Type == html.TextNode && n.Type == html.TextNode {
		if text == nil {
			// n holds its text already: it is moved by the adoption agency
			// algorithm, and may be the node being written.
			if p.textNode == n {
				p.flushText()
			}
			p.appendText(prev, nil)
			p.textBuf = append(p.textBuf, n.Data...)
			return
		}
		p.appendText(prev, text)
		return
	}

	parent.InsertBefore(n, table)
	if n.Type == html.TextNode && text != nil {
		p.appendText(n, text)
	}
}

// addText adds text to the preceding node if it is a text node, or else it
// calls addChild with a new text node.
func (p *parser) addText(text []byte) {
	if len(text) == 0 {
		return
	}

	if p.shouldFosterParent() {
		p.fosterParent(&html.Node{Type: html.TextNode}, text)
		return
	}

	t := p.top()
	if n := t.LastChild; n != nil && n.Type == html.TextNode {
		p.appendText(n, text)
		return
	}
	n := &html.Node{Type: html.TextNode}
	p.addChild(n)
	p.appendText(n, text)
}

// appendText adds text to the text node n, which becomes the node being
// written; the node written before it gets its Data.
func (p *parser) appendText(n *html.Node, text []byte) {
	if p.textNode != n {
		p.flushText()
		p.textNode = n
		p.textBuf = append(p.textBuf[:0], n.Data...)
	}
	p.textBuf = append(p.textBuf, text...)
}

// flushText gives the text node being written its Data.
func (p *parser) flushText() {
	if p.textNode != nil {
		p.textNode.Data = string(p.textBuf)
		p.textNode = nil
		p.textBuf = p.textBuf[:0]
	}
}

func attrCompare(a, b html.Attribute) int {
	return cmp.Or(
		cmp.Compare(a.Namespace, b.Namespace),
		cmp.Compare(a.Key, b.Key),
		cmp.Compare(a.Val, b.Val),
	)
}

// addElement adds a child element based on the current token.
func (p *parser) addElement() {
	p.addChild(&html.Node{
		Type:     html.ElementNode,
		DataAtom: p.tok.DataAtom,
		Data:     p.tok.Data,
		Attr:     p.tok.Attr,
	})
}

// Section 12.2.4.3.
func (p *parser) addFormattingElement() {
	tagAtom, attr := p.tok.DataAtom, p.tok.Attr
	p.addElement()

	// In order to optimize the search, we need the attributes to be sorted, so we
	// can just use slices.Equal.
	slices.SortFunc(attr, attrCompare)

	// Implement the Noah's Ark clause, but with three per family instead of two.
	identicalElements := 0
findIdenticalElements:
	for i := len(p.afe) - 1; i >= 0; i-- {
		n := p.afe[i]
		if n.Type == scopeMarkerNode {
			break
		}
		if n.Type != html.ElementNode {
			continue
		}
		if n.Namespace != "" {
			continue
		}
		if n.DataAtom != tagAtom {
			continue
		}
		if !slices.Equal(n.Attr, attr) {
			continue findIdenticalElements
		}

		identicalElements++
		if identicalElements >= 3 {
			p.afe.remove(n)
		}
	}

	// Sort the attributes to optimize future identical-element searches.
	top := p.top()
	slices.SortFunc(top.Attr, attrCompare)

	p.afe = append(p.afe, top)
}

// Section 12.2.4.3.
func (p *parser) clearActiveFormattingElements() {
	for {
		if n := p.afe.pop(); len(p.afe) == 0 || n.Type == scopeMarkerNode {
			return
		}
	}
}

// Section 12.2.4.3.
func (p *parser) reconstructActiveFormattingElements() {
	n := p.afe.top()
	if n == nil {
		return
	}
	if n.Type == scopeMarkerNode || p.oe.index(n) != -1 {
		return
	}
	i := len(p.afe) - 1
	for n.Type != scopeMarkerNode && p.oe.index(n) == -1 {
		if i == 0 {
			i = -1
			break
		}
		i--
		n = p.afe[i]
	}
	for {
		i++
		clone := cloneNode(p.afe[i])
		p.addChild(clone)
		p.afe[i] = clone
		if i == len(p.afe)-1 {
			break
		}
	}
}

// Section 12.2.5.
func (p *parser) acknowledgeSelfClosingTag() {
	p.hasSelfClosingToken = false
}

// An insertion mode (section 12.2.4.1) is the state transition function from
// a particular state in the HTML5 parser's state machine. It updates the
// parser's fields depending on parser.tok (where EOF means EOF).
// It returns whether the token was consumed.
type insertionMode func(*parser) bool

// setOriginalIM sets the insertion mode to return to after completing a text or
// inTableText insertion mode.
// Section 12.2.4.1, "using the rules for".
func (p *parser) setOriginalIM() {
	if p.originalIM != nil {
		panic("html: bad parser state: originalIM was set twice")
	}
	p.originalIM = p.im
}

// Section 12.2.4.1, "reset the insertion mode".
func (p *parser) resetInsertionMode() {
	for i := len(p.oe) - 1; i >= 0; i-- {
		n := p.oe[i]
		last := i == 0
		if last && p.context != nil {
			n = p.context
		}

		switch n.DataAtom {
		case a.Td, a.Th:
			// TODO: remove this divergence from the HTML5 spec.
			//
			// See https://bugs.chromium.org/p/chromium/issues/detail?id=829668
			p.im = inCellIM
		case a.Tr:
			p.im = inRowIM
		case a.Tbody, a.Thead, a.Tfoot:
			p.im = inTableBodyIM
		case a.Caption:
			p.im = inCaptionIM
		case a.Colgroup:
			p.im = inColumnGroupIM
		case a.Table:
			p.im = inTableIM
		case a.Template:
			// TODO: remove this divergence from the HTML5 spec.
			if n.Namespace != "" {
				continue
			}
			p.im = p.templateStack.top()
		case a.Head:
			// TODO: remove this divergence from the HTML5 spec.
			//
			// See https://bugs.chromium.org/p/chromium/issues/detail?id=829668
			p.im = inHeadIM
		case a.Body:
			p.im = inBodyIM
		case a.Frameset:
			p.im = inFramesetIM
		case a.Html:
			if p.head == nil {
				p.im = beforeHeadIM
			} else {
				p.im = afterHeadIM
			}
		default:
			if last {
				p.im = inBodyIM
				return
			}
			continue
		}
		return
	}
}

const whitespace = " \t\r\n\f"

// Section 12.2.6.4.1.
func initialIM(p *parser) bool {
	switch p.tok.Type {
	case Text:
		p.tok.Text = bytes.TrimLeft(p.tok.Text, whitespace)
		if len(p.tok.Text) == 0 {
			// It was all whitespace, so ignore it.
			return true
		}
	case Comment:
		p.doc.AppendChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
		return true
	case Doctype:
		n, quirks := parseDoctype(p.tok.Data)
		p.doc.AppendChild(n)
		p.quirks = quirks
		p.im = beforeHTMLIM
		return true
	}
	p.quirks = true
	p.im = beforeHTMLIM
	return false
}

// Section 12.2.6.4.2.
func beforeHTMLIM(p *parser) bool {
	switch p.tok.Type {
	case Doctype:
		// Ignore the token.
		return true
	case Text:
		p.tok.Text = bytes.TrimLeft(p.tok.Text, whitespace)
		if len(p.tok.Text) == 0 {
			// It was all whitespace, so ignore it.
			return true
		}
	case StartTag:
		if p.tok.DataAtom == a.Html {
			p.addElement()
			p.im = beforeHeadIM
			return true
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Head, a.Body, a.Html, a.Br:
			p.parseImpliedToken(StartTag, a.Html, a.Html.String())
			return false
		default:
			// Ignore the token.
			return true
		}
	case Comment:
		p.doc.AppendChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
		return true
	}
	p.parseImpliedToken(StartTag, a.Html, a.Html.String())
	return false
}

// Section 12.2.6.4.3.
func beforeHeadIM(p *parser) bool {
	switch p.tok.Type {
	case Text:
		p.tok.Text = bytes.TrimLeft(p.tok.Text, whitespace)
		if len(p.tok.Text) == 0 {
			// It was all whitespace, so ignore it.
			return true
		}
	case StartTag:
		switch p.tok.DataAtom {
		case a.Head:
			p.addElement()
			p.head = p.top()
			p.im = inHeadIM
			return true
		case a.Html:
			return inBodyIM(p)
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Head, a.Body, a.Html, a.Br:
			p.parseImpliedToken(StartTag, a.Head, a.Head.String())
			return false
		default:
			// Ignore the token.
			return true
		}
	case Comment:
		p.addChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
		return true
	case Doctype:
		// Ignore the token.
		return true
	}

	p.parseImpliedToken(StartTag, a.Head, a.Head.String())
	return false
}

// Section 12.2.6.4.4.
func inHeadIM(p *parser) bool {
	switch p.tok.Type {
	case Text:
		s := bytes.TrimLeft(p.tok.Text, whitespace)
		if len(s) < len(p.tok.Text) {
			// Add the initial whitespace to the current node.
			p.addText(p.tok.Text[:len(p.tok.Text)-len(s)])
			if len(s) == 0 {
				return true
			}
			p.tok.Text = s
		}
	case StartTag:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Base, a.Basefont, a.Bgsound, a.Link, a.Meta:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			return true
		case a.Noscript:
			if p.scripting {
				p.parseGenericRawTextElement()
				return true
			}
			p.addElement()
			p.im = inHeadNoscriptIM
			// Don't let the tokenizer go into raw text mode when scripting is disabled.
			p.r.NextIsNotRawText()
			return true
		case a.Script, a.Title:
			p.addElement()
			p.setOriginalIM()
			p.im = textIM
			return true
		case a.Noframes, a.Style:
			p.parseGenericRawTextElement()
			return true
		case a.Head:
			// Ignore the token.
			return true
		case a.Template:
			// TODO: remove this divergence from the HTML5 spec.
			//
			// We don't handle all of the corner cases when mixing foreign
			// content (i.e. <math> or <svg>) with <template>. Without this
			// early return, we can get into an infinite loop, possibly because
			// of the "TODO... further divergence" a little below.
			//
			// As a workaround, if we are mixing foreign content and templates,
			// just ignore the rest of the HTML. Foreign content is rare and a
			// relatively old HTML feature. Templates are also rare and a
			// relatively new HTML feature. Their combination is very rare.
			for _, e := range p.oe {
				if e.Namespace != "" {
					p.im = ignoreTheRemainingTokens
					return true
				}
			}

			p.addElement()
			p.afe = append(p.afe, &scopeMarker)
			p.framesetOK = false
			p.im = inTemplateIM
			p.templateStack = append(p.templateStack, inTemplateIM)
			return true
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Head:
			p.oe.pop()
			p.im = afterHeadIM
			return true
		case a.Body, a.Html, a.Br:
			p.parseImpliedToken(EndTag, a.Head, a.Head.String())
			return false
		case a.Template:
			if !p.oe.contains(a.Template) {
				return true
			}
			// TODO: remove this further divergence from the HTML5 spec.
			//
			// See https://bugs.chromium.org/p/chromium/issues/detail?id=829668
			p.generateImpliedEndTags()
			for i := len(p.oe) - 1; i >= 0; i-- {
				if n := p.oe[i]; n.Namespace == "" && n.DataAtom == a.Template {
					p.oe = p.oe[:i]
					break
				}
			}
			p.clearActiveFormattingElements()
			p.templateStack.pop()
			p.resetInsertionMode()
			return true
		default:
			// Ignore the token.
			return true
		}
	case Comment:
		p.addChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
		return true
	case Doctype:
		// Ignore the token.
		return true
	}

	p.parseImpliedToken(EndTag, a.Head, a.Head.String())
	return false
}

// Section 12.2.6.4.5.
func inHeadNoscriptIM(p *parser) bool {
	switch p.tok.Type {
	case Doctype:
		// Ignore the token.
		return true
	case StartTag:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Basefont, a.Bgsound, a.Link, a.Meta, a.Noframes, a.Style:
			return inHeadIM(p)
		case a.Head:
			// Ignore the token.
			return true
		case a.Noscript:
			// Don't let the tokenizer go into raw text mode even when a <noscript>
			// tag is in "in head noscript" insertion mode.
			p.r.NextIsNotRawText()
			// Ignore the token.
			return true
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Noscript, a.Br:
		default:
			// Ignore the token.
			return true
		}
	case Text:
		s := bytes.TrimLeft(p.tok.Text, whitespace)
		if len(s) == 0 {
			// It was all whitespace.
			return inHeadIM(p)
		}
	case Comment:
		return inHeadIM(p)
	}
	p.oe.pop()
	if p.top().DataAtom != a.Head {
		panic("html: the new current node will be a head element.")
	}
	p.im = inHeadIM
	if p.tok.DataAtom == a.Noscript {
		return true
	}
	return false
}

// Section 12.2.6.4.6.
func afterHeadIM(p *parser) bool {
	switch p.tok.Type {
	case Text:
		s := bytes.TrimLeft(p.tok.Text, whitespace)
		if len(s) < len(p.tok.Text) {
			// Add the initial whitespace to the current node.
			p.addText(p.tok.Text[:len(p.tok.Text)-len(s)])
			if len(s) == 0 {
				return true
			}
			p.tok.Text = s
		}
	case StartTag:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Body:
			p.addElement()
			p.framesetOK = false
			p.im = inBodyIM
			return true
		case a.Frameset:
			p.addElement()
			p.im = inFramesetIM
			return true
		case a.Base, a.Basefont, a.Bgsound, a.Link, a.Meta, a.Noframes, a.Script, a.Style, a.Template, a.Title:
			p.insertOpenElement(p.head)
			defer p.oe.remove(p.head)
			return inHeadIM(p)
		case a.Head:
			// Ignore the token.
			return true
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Body, a.Html, a.Br:
			// Drop down to creating an implied <body> tag.
		case a.Template:
			return inHeadIM(p)
		default:
			// Ignore the token.
			return true
		}
	case Comment:
		p.addChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
		return true
	case Doctype:
		// Ignore the token.
		return true
	}

	p.parseImpliedToken(StartTag, a.Body, a.Body.String())
	p.framesetOK = true
	if p.tok.Type == EOF {
		// Stop parsing.
		return true
	}
	return false
}

// copyAttributes copies attributes of src not found on dst to dst.
func copyAttributes(dst *html.Node, src token) {
	if len(src.Attr) == 0 {
		return
	}
	attr := map[string]string{}
	for _, t := range dst.Attr {
		attr[t.Key] = t.Val
	}
	for _, t := range src.Attr {
		if _, ok := attr[t.Key]; !ok {
			dst.Attr = append(dst.Attr, t)
			attr[t.Key] = t.Val
		}
	}
}

// Section 12.2.6.4.7.
func inBodyIM(p *parser) bool {
	switch p.tok.Type {
	case Text:
		d := p.tok.Text
		switch n := p.oe.top(); n.DataAtom {
		case a.Pre, a.Listing:
			if n.FirstChild == nil {
				// Ignore a newline at the start of a <pre> block.
				if len(d) != 0 && d[0] == '\r' {
					d = d[1:]
				}
				if len(d) != 0 && d[0] == '\n' {
					d = d[1:]
				}
			}
		}
		d = dropNUL(d)
		if len(d) == 0 {
			return true
		}
		p.reconstructActiveFormattingElements()
		p.addText(d)
		if p.framesetOK && len(bytes.TrimLeft(d, whitespace)) != 0 {
			// There were non-whitespace characters inserted.
			p.framesetOK = false
		}
	case StartTag:
		switch p.tok.DataAtom {
		case a.Html:
			if p.oe.contains(a.Template) {
				return true
			}
			copyAttributes(p.oe[0], p.tok)
		case a.Base, a.Basefont, a.Bgsound, a.Link, a.Meta, a.Noframes, a.Script, a.Style, a.Template, a.Title:
			return inHeadIM(p)
		case a.Body:
			if p.oe.contains(a.Template) {
				return true
			}
			if len(p.oe) >= 2 {
				body := p.oe[1]
				if body.Type == html.ElementNode && body.DataAtom == a.Body {
					p.framesetOK = false
					copyAttributes(body, p.tok)
				}
			}
		case a.Frameset:
			if !p.framesetOK || len(p.oe) < 2 || p.oe[1].DataAtom != a.Body {
				// Ignore the token.
				return true
			}
			body := p.oe[1]
			if body.Parent != nil {
				body.Parent.RemoveChild(body)
			}
			p.oe = p.oe[:1]
			p.addElement()
			p.im = inFramesetIM
			return true
		case a.Address, a.Article, a.Aside, a.Blockquote, a.Center, a.Details, a.Dialog, a.Dir, a.Div, a.Dl, a.Fieldset, a.Figcaption, a.Figure, a.Footer, a.Header, a.Hgroup, a.Main, a.Menu, a.Nav, a.Ol, a.P, a.Search, a.Section, a.Summary, a.Ul:
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.H1, a.H2, a.H3, a.H4, a.H5, a.H6:
			p.popUntil(buttonScope, a.P)
			switch n := p.top(); n.DataAtom {
			case a.H1, a.H2, a.H3, a.H4, a.H5, a.H6:
				p.oe.pop()
			}
			p.addElement()
		case a.Pre, a.Listing:
			p.popUntil(buttonScope, a.P)
			p.addElement()
			// The newline, if any, will be dealt with by the Text case.
			p.framesetOK = false
		case a.Form:
			if p.form != nil && !p.oe.contains(a.Template) {
				// Ignore the token
				return true
			}
			p.popUntil(buttonScope, a.P)
			p.addElement()
			if !p.oe.contains(a.Template) {
				p.form = p.top()
			}
		case a.Li:
			p.framesetOK = false
			for i := len(p.oe) - 1; i >= 0; i-- {
				node := p.oe[i]
				switch node.DataAtom {
				case a.Li:
					p.oe = p.oe[:i]
				case a.Address, a.Div, a.P:
					continue
				default:
					if !isSpecialElement(node) {
						continue
					}
				}
				break
			}
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.Dd, a.Dt:
			p.framesetOK = false
			for i := len(p.oe) - 1; i >= 0; i-- {
				node := p.oe[i]
				switch node.DataAtom {
				case a.Dd, a.Dt:
					p.oe = p.oe[:i]
				case a.Address, a.Div, a.P:
					continue
				default:
					if !isSpecialElement(node) {
						continue
					}
				}
				break
			}
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.Plaintext:
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.Button:
			if p.elementInScope(defaultScope, a.Button) {
				p.generateImpliedEndTags()
				p.popUntil(defaultScope, a.Button)
			}
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.framesetOK = false
		case a.A:
			for i := len(p.afe) - 1; i >= 0 && p.afe[i].Type != scopeMarkerNode; i-- {
				if n := p.afe[i]; n.Type == html.ElementNode && n.DataAtom == a.A {
					p.inBodyEndTagFormatting(a.A, "a")
					p.oe.remove(n)
					p.afe.remove(n)
					break
				}
			}
			p.reconstructActiveFormattingElements()
			p.addFormattingElement()
		case a.B, a.Big, a.Code, a.Em, a.Font, a.I, a.S, a.Small, a.Strike, a.Strong, a.Tt, a.U:
			p.reconstructActiveFormattingElements()
			p.addFormattingElement()
		case a.Nobr:
			p.reconstructActiveFormattingElements()
			if p.elementInScope(defaultScope, a.Nobr) {
				p.inBodyEndTagFormatting(a.Nobr, "nobr")
				p.reconstructActiveFormattingElements()
			}
			p.addFormattingElement()
		case a.Applet, a.Marquee, a.Object:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.afe = append(p.afe, &scopeMarker)
			p.framesetOK = false
		case a.Table:
			if !p.quirks {
				p.popUntil(buttonScope, a.P)
			}
			p.addElement()
			p.framesetOK = false
			p.im = inTableIM
			return true
		case a.Area, a.Br, a.Embed, a.Img, a.Keygen, a.Wbr:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			p.framesetOK = false
		case a.Input:
			if p.fragment && p.context.DataAtom == a.Select {
				// Ignore the token.
				return true
			}
			p.popUntil(defaultScope, a.Select)
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			if p.tok.DataAtom == a.Input {
				for _, t := range p.tok.Attr {
					if t.Key == "type" {
						if strings.EqualFold(t.Val, "hidden") {
							// Skip setting framesetOK = false
							return true
						}
					}
				}
			}
			p.framesetOK = false
		case a.Param, a.Source, a.Track:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
		case a.Hr:
			if p.elementInScope(buttonScope, a.P) {
				p.generateImpliedEndTags("p")
				p.popUntil(defaultScope, a.P)
			}
			if p.elementInScope(defaultScope, a.Select) {
				p.generateImpliedEndTags()
			}
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			p.framesetOK = false
		case a.Image:
			p.tok.DataAtom = a.Img
			p.tok.Data = a.Img.String()
			return false
		case a.Textarea:
			p.addElement()
			p.setOriginalIM()
			p.framesetOK = false
			p.im = textIM
		case a.Xmp:
			p.popUntil(buttonScope, a.P)
			p.reconstructActiveFormattingElements()
			p.framesetOK = false
			p.parseGenericRawTextElement()
		case a.Iframe:
			p.framesetOK = false
			p.parseGenericRawTextElement()
		case a.Noembed:
			p.parseGenericRawTextElement()
		case a.Noscript:
			if p.scripting {
				p.parseGenericRawTextElement()
				return true
			}
			p.reconstructActiveFormattingElements()
			p.addElement()
			// Don't let the tokenizer go into raw text mode when scripting is disabled.
			p.r.NextIsNotRawText()
		case a.Select:
			if p.fragment && p.context.DataAtom == a.Select {
				// Ignore the token.
				return true
			} else if p.popUntil(defaultScope, a.Select) {
				return true
			}
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.framesetOK = false
			return true
		case a.Option:
			if p.elementInScope(defaultScope, a.Select) {
				p.generateImpliedEndTags("optgroup")
				// If oe has option element in scope, parse error?
			} else if p.top().DataAtom == a.Option {
				p.oe.pop()
			}
			p.reconstructActiveFormattingElements()
			p.addElement()
		case a.Optgroup:
			if p.elementInScope(defaultScope, a.Select) {
				p.generateImpliedEndTags()
				// If oe has option or optgroup element in scope, parse error?
			} else if p.top().DataAtom == a.Option {
				p.oe.pop()
			}
			p.reconstructActiveFormattingElements()
			p.addElement()
		case a.Rb, a.Rtc:
			if p.elementInScope(defaultScope, a.Ruby) {
				p.generateImpliedEndTags()
			}
			p.addElement()
		case a.Rp, a.Rt:
			if p.elementInScope(defaultScope, a.Ruby) {
				p.generateImpliedEndTags("rtc")
			}
			p.addElement()
		case a.Math, a.Svg:
			p.reconstructActiveFormattingElements()
			if p.tok.DataAtom == a.Math {
				adjustAttributeNames(p.tok.Attr, mathMLAttributeAdjustments)
			} else {
				adjustAttributeNames(p.tok.Attr, svgAttributeAdjustments)
			}
			adjustForeignAttributes(p.tok.Attr)
			p.addElement()
			p.top().Namespace = p.tok.Data
			if p.hasSelfClosingToken {
				p.oe.pop()
				p.acknowledgeSelfClosingTag()
			}
			return true
		case a.Caption, a.Col, a.Colgroup, a.Frame, a.Head, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			// Ignore the token.
		default:
			p.reconstructActiveFormattingElements()
			p.addElement()
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Body:
			if p.elementInScope(defaultScope, a.Body) {
				p.im = afterBodyIM
			}
		case a.Html:
			if p.elementInScope(defaultScope, a.Body) {
				p.parseImpliedToken(EndTag, a.Body, a.Body.String())
				return false
			}
			return true
		case a.Address, a.Article, a.Aside, a.Blockquote, a.Button, a.Center, a.Details, a.Dialog, a.Dir, a.Div, a.Dl, a.Fieldset, a.Figcaption, a.Figure, a.Footer, a.Header, a.Hgroup, a.Listing, a.Main, a.Menu, a.Nav, a.Ol, a.Pre, a.Search, a.Section, a.Select, a.Summary, a.Ul:
			if !p.elementInScope(defaultScope, p.tok.DataAtom) {
				// Ignore the token.
				return true
			}
			p.generateImpliedEndTags()
			p.popUntil(defaultScope, p.tok.DataAtom)
		case a.Form:
			if p.oe.contains(a.Template) {
				i := p.indexOfElementInScope(defaultScope, a.Form)
				if i == -1 {
					// Ignore the token.
					return true
				}
				p.generateImpliedEndTags()
				if p.oe[i].DataAtom != a.Form {
					// Ignore the token.
					return true
				}
				p.popUntil(defaultScope, a.Form)
			} else {
				node := p.form
				p.form = nil
				i := p.indexOfElementInScope(defaultScope, a.Form)
				if node == nil || i == -1 || p.oe[i] != node {
					// Ignore the token.
					return true
				}
				p.generateImpliedEndTags()
				p.oe.remove(node)
			}
		case a.P:
			if !p.elementInScope(buttonScope, a.P) {
				p.parseImpliedToken(StartTag, a.P, a.P.String())
			}
			p.popUntil(buttonScope, a.P)
		case a.Li:
			p.popUntil(listItemScope, a.Li)
		case a.Dd, a.Dt:
			p.popUntil(defaultScope, p.tok.DataAtom)
		case a.H1, a.H2, a.H3, a.H4, a.H5, a.H6:
			p.popUntil(defaultScope, a.H1, a.H2, a.H3, a.H4, a.H5, a.H6)
		case a.A, a.B, a.Big, a.Code, a.Em, a.Font, a.I, a.Nobr, a.S, a.Small, a.Strike, a.Strong, a.Tt, a.U:
			p.inBodyEndTagFormatting(p.tok.DataAtom, p.tok.Data)
		case a.Applet, a.Marquee, a.Object:
			if p.popUntil(defaultScope, p.tok.DataAtom) {
				p.clearActiveFormattingElements()
			}
		case a.Br:
			p.tok.Type = StartTag
			return false
		case a.Template:
			return inHeadIM(p)
		default:
			p.inBodyEndTagOther(p.tok.DataAtom, p.tok.Data)
		}
	case Comment:
		p.addChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
	case EOF:
		// TODO: remove this divergence from the HTML5 spec.
		if len(p.templateStack) > 0 {
			p.im = inTemplateIM
			return false
		}
		for _, e := range p.oe {
			switch e.DataAtom {
			case a.Dd, a.Dt, a.Li, a.Optgroup, a.Option, a.P, a.Rb, a.Rp, a.Rt, a.Rtc, a.Tbody, a.Td, a.Tfoot, a.Th,
				a.Thead, a.Tr, a.Body, a.Html:
			default:
				return true
			}
		}
	}

	return true
}

func (p *parser) inBodyEndTagFormatting(tagAtom a.Atom, tagName string) {
	// This is the "adoption agency" algorithm, described at
	// https://html.spec.whatwg.org/multipage/syntax.html#adoptionAgency

	// TODO: this is a fairly literal line-by-line translation of that algorithm.
	// Once the code successfully parses the comprehensive test suite, we should
	// refactor this code to be more idiomatic.

	// Steps 1-2
	if current := p.oe.top(); current.Data == tagName && p.afe.index(current) == -1 {
		p.oe.pop()
		return
	}

	// Steps 3-5. The outer loop.
	for i := 0; i < 8; i++ {
		// Step 6. Find the formatting element.
		var formattingElement *html.Node
		for j := len(p.afe) - 1; j >= 0; j-- {
			if p.afe[j].Type == scopeMarkerNode {
				break
			}
			if p.afe[j].DataAtom == tagAtom {
				formattingElement = p.afe[j]
				break
			}
		}
		if formattingElement == nil {
			p.inBodyEndTagOther(tagAtom, tagName)
			return
		}

		// Step 7. Ignore the tag if formatting element is not in the stack of open elements.
		feIndex := p.oe.index(formattingElement)
		if feIndex == -1 {
			p.afe.remove(formattingElement)
			return
		}
		// Step 8. Ignore the tag if formatting element is not in the scope.
		if !p.elementInScope(defaultScope, tagAtom) {
			// Ignore the tag.
			return
		}

		// Step 9. This step is omitted because it's just a parse error but no need to return.

		// Steps 10-11. Find the furthest block.
		var furthestBlock *html.Node
		for _, e := range p.oe[feIndex:] {
			if isSpecialElement(e) {
				furthestBlock = e
				break
			}
		}
		if furthestBlock == nil {
			e := p.oe.pop()
			for e != formattingElement {
				e = p.oe.pop()
			}
			p.afe.remove(e)
			return
		}

		// Steps 12-13. Find the common ancestor and bookmark node.
		commonAncestor := p.oe[feIndex-1]
		bookmark := p.afe.index(formattingElement)

		// Step 14. The inner loop. Find the lastNode to reparent.
		lastNode := furthestBlock
		node := furthestBlock
		x := p.oe.index(node)
		// Step 14.1.
		j := 0
		for {
			// Step 14.2.
			j++
			// Step. 14.3.
			x--
			node = p.oe[x]
			// Step 14.4. Go to the next step if node is formatting element.
			if node == formattingElement {
				break
			}
			// Step 14.5. Remove node from the list of active formatting elements if
			// inner loop counter is greater than three and node is in the list of
			// active formatting elements.
			if ni := p.afe.index(node); j > 3 && ni > -1 {
				p.afe.remove(node)
				// If any element of the list of active formatting elements is removed,
				// we need to take care whether bookmark should be decremented or not.
				// This is because the value of bookmark may exceed the size of the
				// list by removing elements from the list.
				if ni <= bookmark {
					bookmark--
				}
				continue
			}
			// Step 14.6. Continue the next inner loop if node is not in the list of
			// active formatting elements.
			if p.afe.index(node) == -1 {
				p.oe.remove(node)
				continue
			}
			// Step 14.7.
			clone := cloneNode(node)
			p.afe[p.afe.index(node)] = clone
			p.oe[p.oe.index(node)] = clone
			node = clone
			// Step 14.8.
			if lastNode == furthestBlock {
				bookmark = p.afe.index(node) + 1
			}
			// Step 14.9.
			if lastNode.Parent != nil {
				lastNode.Parent.RemoveChild(lastNode)
			}
			node.AppendChild(lastNode)
			// Step 14.10.
			lastNode = node
		}

		// Step 15. Reparent lastNode to the common ancestor,
		// or for misnested table nodes, to the foster parent.
		if lastNode.Parent != nil {
			lastNode.Parent.RemoveChild(lastNode)
		}
		switch commonAncestor.DataAtom {
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			p.fosterParent(lastNode, nil)
		default:
			commonAncestor.AppendChild(lastNode)
		}

		// Steps 16-18. Reparent nodes from the furthest block's children
		// to a clone of the formatting element.
		clone := cloneNode(formattingElement)
		reparentChildren(clone, furthestBlock)
		furthestBlock.AppendChild(clone)

		// Step 19. Fix up the list of active formatting elements.
		if oldLoc := p.afe.index(formattingElement); oldLoc != -1 && oldLoc < bookmark {
			// Move the bookmark with the rest of the list.
			bookmark--
		}
		p.afe.remove(formattingElement)
		p.afe.insert(bookmark, clone)

		// Step 20. Fix up the stack of open elements.
		p.oe.remove(formattingElement)
		p.oe.insert(p.oe.index(furthestBlock)+1, clone)
	}
}

// inBodyEndTagOther performs the "any other end tag" algorithm for inBodyIM.
func (p *parser) inBodyEndTagOther(tagAtom a.Atom, tagName string) {
	for i := len(p.oe) - 1; i >= 0; i-- {
		// Two element nodes have the same tag if they have the same Data (a
		// string-typed field). As an optimization, for common HTML tags, each
		// Data string is assigned a unique, non-zero DataAtom (a uint32-typed
		// field), since integer comparison is faster than string comparison.
		// Uncommon (custom) tags get a zero DataAtom.
		//
		// The if condition here is equivalent to (p.oe[i].Data == tagName).
		if p.oe[i].Namespace == "" && (p.oe[i].DataAtom == tagAtom) &&
			((tagAtom != 0) || (p.oe[i].Data == tagName)) {
			p.oe = p.oe[:i]
			break
		}
		if isSpecialElement(p.oe[i]) {
			break
		}
	}
}

// Section 12.2.6.4.8.
func textIM(p *parser) bool {
	switch p.tok.Type {
	case EOF:
		p.oe.pop()
	case Text:
		d := p.tok.Text
		if n := p.oe.top(); n.DataAtom == a.Textarea && n.FirstChild == nil {
			// Ignore a newline at the start of a <textarea> block.
			if len(d) != 0 && d[0] == '\r' {
				d = d[1:]
			}
			if len(d) != 0 && d[0] == '\n' {
				d = d[1:]
			}
		}
		if len(d) == 0 {
			return true
		}
		p.addText(d)
		return true
	case EndTag:
		p.oe.pop()
	}
	p.im = p.originalIM
	p.originalIM = nil
	return p.tok.Type == EndTag
}

// Section 12.2.6.4.9.
func inTableIM(p *parser) bool {
	switch p.tok.Type {
	case Text:
		p.tok.Text = dropNUL(p.tok.Text)
		switch p.oe.top().DataAtom {
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			if len(bytes.Trim(p.tok.Text, whitespace)) == 0 {
				p.addText(p.tok.Text)
				return true
			}
		}
	case StartTag:
		switch p.tok.DataAtom {
		case a.Caption:
			p.clearStackToContext(tableScope)
			p.afe = append(p.afe, &scopeMarker)
			p.addElement()
			p.im = inCaptionIM
			return true
		case a.Colgroup:
			p.clearStackToContext(tableScope)
			p.addElement()
			p.im = inColumnGroupIM
			return true
		case a.Col:
			p.parseImpliedToken(StartTag, a.Colgroup, a.Colgroup.String())
			return false
		case a.Tbody, a.Tfoot, a.Thead:
			p.clearStackToContext(tableScope)
			p.addElement()
			p.im = inTableBodyIM
			return true
		case a.Td, a.Th, a.Tr:
			p.parseImpliedToken(StartTag, a.Tbody, a.Tbody.String())
			return false
		case a.Table:
			if p.popUntil(tableScope, a.Table) {
				p.resetInsertionMode()
				return false
			}
			// Ignore the token.
			return true
		case a.Style, a.Script, a.Template:
			return inHeadIM(p)
		case a.Input:
			for _, t := range p.tok.Attr {
				if t.Key == "type" && strings.EqualFold(t.Val, "hidden") {
					p.addElement()
					p.oe.pop()
					return true
				}
			}
			// Otherwise drop down to the default action.
		case a.Form:
			if p.oe.contains(a.Template) || p.form != nil {
				// Ignore the token.
				return true
			}
			p.addElement()
			p.form = p.oe.pop()
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Table:
			if p.popUntil(tableScope, a.Table) {
				p.resetInsertionMode()
				return true
			}
			// Ignore the token.
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			// Ignore the token.
			return true
		case a.Template:
			return inHeadIM(p)
		}
	case Comment:
		p.addChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
		return true
	case Doctype:
		// Ignore the token.
		return true
	case EOF:
		return inBodyIM(p)
	}

	p.fosterParenting = true
	defer func() { p.fosterParenting = false }()

	return inBodyIM(p)
}

// Section 12.2.6.4.11.
func inCaptionIM(p *parser) bool {
	switch p.tok.Type {
	case StartTag:
		switch p.tok.DataAtom {
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Td, a.Tfoot, a.Thead, a.Tr:
			if !p.popUntil(tableScope, a.Caption) {
				// Ignore the token.
				return true
			}
			p.clearActiveFormattingElements()
			p.im = inTableIM
			return false
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Caption:
			if p.popUntil(tableScope, a.Caption) {
				p.clearActiveFormattingElements()
				p.im = inTableIM
			}
			return true
		case a.Table:
			if !p.popUntil(tableScope, a.Caption) {
				// Ignore the token.
				return true
			}
			p.clearActiveFormattingElements()
			p.im = inTableIM
			return false
		case a.Body, a.Col, a.Colgroup, a.Html, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			// Ignore the token.
			return true
		}
	}
	return inBodyIM(p)
}

// Section 12.2.6.4.12.
func inColumnGroupIM(p *parser) bool {
	switch p.tok.Type {
	case Text:
		s := bytes.TrimLeft(p.tok.Text, whitespace)
		if len(s) < len(p.tok.Text) {
			// Add the initial whitespace to the current node.
			p.addText(p.tok.Text[:len(p.tok.Text)-len(s)])
			if len(s) == 0 {
				return true
			}
			p.tok.Text = s
		}
	case Comment:
		p.addChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
		return true
	case Doctype:
		// Ignore the token.
		return true
	case StartTag:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Col:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			return true
		case a.Template:
			return inHeadIM(p)
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Colgroup:
			if p.oe.top().DataAtom == a.Colgroup {
				p.oe.pop()
				p.im = inTableIM
			}
			return true
		case a.Col:
			// Ignore the token.
			return true
		case a.Template:
			return inHeadIM(p)
		}
	case EOF:
		return inBodyIM(p)
	}
	if p.oe.top().DataAtom != a.Colgroup {
		return true
	}
	p.oe.pop()
	p.im = inTableIM
	return false
}

// Section 12.2.6.4.13.
func inTableBodyIM(p *parser) bool {
	switch p.tok.Type {
	case StartTag:
		switch p.tok.DataAtom {
		case a.Tr:
			p.clearStackToContext(tableBodyScope)
			p.addElement()
			p.im = inRowIM
			return true
		case a.Td, a.Th:
			p.parseImpliedToken(StartTag, a.Tr, a.Tr.String())
			return false
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Tfoot, a.Thead:
			if p.popUntil(tableScope, a.Tbody, a.Thead, a.Tfoot) {
				p.im = inTableIM
				return false
			}
			// Ignore the token.
			return true
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Tbody, a.Tfoot, a.Thead:
			if p.elementInScope(tableScope, p.tok.DataAtom) {
				p.clearStackToContext(tableBodyScope)
				p.oe.pop()
				p.im = inTableIM
			}
			return true
		case a.Table:
			if p.popUntil(tableScope, a.Tbody, a.Thead, a.Tfoot) {
				p.im = inTableIM
				return false
			}
			// Ignore the token.
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html, a.Td, a.Th, a.Tr:
			// Ignore the token.
			return true
		}
	case Comment:
		p.addChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
		return true
	}

	return inTableIM(p)
}

// Section 13.2.6.4.14.
func inRowIM(p *parser) bool {
	switch p.tok.Type {
	case StartTag:
		switch p.tok.DataAtom {
		case a.Td, a.Th:
			p.clearStackToContext(tableRowScope)
			p.addElement()
			p.afe = append(p.afe, &scopeMarker)
			p.im = inCellIM
			return true
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			if p.elementInScope(tableScope, a.Tr) {
				p.clearStackToContext(tableRowScope)
				p.oe.pop()
				p.im = inTableBodyIM
				return false
			}
			// Ignore the token.
			return true
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Tr:
			if p.elementInScope(tableScope, a.Tr) {
				p.clearStackToContext(tableRowScope)
				p.oe.pop()
				p.im = inTableBodyIM
				return true
			}
			// Ignore the token.
			return true
		case a.Table:
			if p.elementInScope(tableScope, a.Tr) {
				p.clearStackToContext(tableRowScope)
				p.oe.pop()
				p.im = inTableBodyIM
				return false
			}
			// Ignore the token.
			return true
		case a.Tbody, a.Tfoot, a.Thead:
			if p.elementInScope(tableScope, p.tok.DataAtom) && p.elementInScope(tableScope, a.Tr) {
				p.clearStackToContext(tableRowScope)
				p.oe.pop()
				p.im = inTableBodyIM
				return false
			}
			// Ignore the token.
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html, a.Td, a.Th:
			// Ignore the token.
			return true
		}
	}

	return inTableIM(p)
}

// Section 12.2.6.4.15.
func inCellIM(p *parser) bool {
	switch p.tok.Type {
	case StartTag:
		switch p.tok.DataAtom {
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			if p.popUntil(tableScope, a.Td, a.Th) {
				// Close the cell and reprocess.
				p.clearActiveFormattingElements()
				p.im = inRowIM
				return false
			}
			// Ignore the token.
			return true
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Td, a.Th:
			if !p.popUntil(tableScope, p.tok.DataAtom) {
				// Ignore the token.
				return true
			}
			p.clearActiveFormattingElements()
			p.im = inRowIM
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html:
			// Ignore the token.
			return true
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			if !p.elementInScope(tableScope, p.tok.DataAtom) {
				// Ignore the token.
				return true
			}
			// Close the cell and reprocess.
			if p.popUntil(tableScope, a.Td, a.Th) {
				p.clearActiveFormattingElements()
			}
			p.im = inRowIM
			return false
		}
	}
	return inBodyIM(p)
}

// Section 12.2.6.4.18.
func inTemplateIM(p *parser) bool {
	switch p.tok.Type {
	case Text, Comment, Doctype:
		return inBodyIM(p)
	case StartTag:
		switch p.tok.DataAtom {
		case a.Base, a.Basefont, a.Bgsound, a.Link, a.Meta, a.Noframes, a.Script, a.Style, a.Template, a.Title:
			return inHeadIM(p)
		case a.Caption, a.Colgroup, a.Tbody, a.Tfoot, a.Thead:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inTableIM)
			p.im = inTableIM
			return false
		case a.Col:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inColumnGroupIM)
			p.im = inColumnGroupIM
			return false
		case a.Tr:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inTableBodyIM)
			p.im = inTableBodyIM
			return false
		case a.Td, a.Th:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inRowIM)
			p.im = inRowIM
			return false
		default:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inBodyIM)
			p.im = inBodyIM
			return false
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Template:
			return inHeadIM(p)
		default:
			// Ignore the token.
			return true
		}
	case EOF:
		if !p.oe.contains(a.Template) {
			// Ignore the token.
			return true
		}
		// TODO: remove this divergence from the HTML5 spec.
		//
		// See https://bugs.chromium.org/p/chromium/issues/detail?id=829668
		p.generateImpliedEndTags()
		for i := len(p.oe) - 1; i >= 0; i-- {
			if n := p.oe[i]; n.Namespace == "" && n.DataAtom == a.Template {
				p.oe = p.oe[:i]
				break
			}
		}
		p.clearActiveFormattingElements()
		p.templateStack.pop()
		p.resetInsertionMode()
		return false
	}
	return false
}

// Section 12.2.6.4.19.
func afterBodyIM(p *parser) bool {
	switch p.tok.Type {
	case EOF:
		// Stop parsing.
		return true
	case Text:
		s := bytes.TrimLeft(p.tok.Text, whitespace)
		if len(s) == 0 {
			// It was all whitespace.
			return inBodyIM(p)
		}
	case StartTag:
		if p.tok.DataAtom == a.Html {
			return inBodyIM(p)
		}
	case EndTag:
		if p.tok.DataAtom == a.Html {
			if !p.fragment {
				p.im = afterAfterBodyIM
			}
			return true
		}
	case Comment:
		// The comment is attached to the <html> element.
		if len(p.oe) < 1 || p.oe[0].DataAtom != a.Html {
			panic("html: bad parser state: <html> element not found, in the after-body insertion mode")
		}
		p.oe[0].AppendChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
		return true
	}
	p.im = inBodyIM
	return false
}

// Section 12.2.6.4.20.
func inFramesetIM(p *parser) bool {
	switch p.tok.Type {
	case Comment:
		p.addChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
	case Text:
		// Ignore all text but whitespace.
		s := keepWhitespace(p.tok.Text)
		if len(s) != 0 {
			p.addText(s)
		}
	case StartTag:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Frameset:
			p.addElement()
		case a.Frame:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
		case a.Noframes:
			return inHeadIM(p)
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Frameset:
			if p.oe.top().DataAtom != a.Html {
				p.oe.pop()
				if p.oe.top().DataAtom != a.Frameset {
					p.im = afterFramesetIM
					return true
				}
			}
		}
	default:
		// Ignore the token.
	}
	return true
}

// Section 12.2.6.4.21.
func afterFramesetIM(p *parser) bool {
	switch p.tok.Type {
	case Comment:
		p.addChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
	case Text:
		// Ignore all text but whitespace.
		s := keepWhitespace(p.tok.Text)
		if len(s) != 0 {
			p.addText(s)
		}
	case StartTag:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Noframes:
			return inHeadIM(p)
		}
	case EndTag:
		switch p.tok.DataAtom {
		case a.Html:
			p.im = afterAfterFramesetIM
			return true
		}
	default:
		// Ignore the token.
	}
	return true
}

// Section 12.2.6.4.22.
func afterAfterBodyIM(p *parser) bool {
	switch p.tok.Type {
	case EOF:
		// Stop parsing.
		return true
	case Text:
		s := bytes.TrimLeft(p.tok.Text, whitespace)
		if len(s) == 0 {
			// It was all whitespace.
			return inBodyIM(p)
		}
	case StartTag:
		if p.tok.DataAtom == a.Html {
			return inBodyIM(p)
		}
	case Comment:
		p.doc.AppendChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
		return true
	case Doctype:
		return inBodyIM(p)
	}
	p.im = inBodyIM
	return false
}

// Section 12.2.6.4.23.
func afterAfterFramesetIM(p *parser) bool {
	switch p.tok.Type {
	case Comment:
		p.doc.AppendChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
	case Text:
		// Ignore all text but whitespace.
		s := keepWhitespace(p.tok.Text)
		if len(s) != 0 {
			p.tok.Text = s
			return inBodyIM(p)
		}
	case StartTag:
		switch p.tok.DataAtom {
		case a.Html:
			return inBodyIM(p)
		case a.Noframes:
			return inHeadIM(p)
		}
	case Doctype:
		return inBodyIM(p)
	default:
		// Ignore the token.
	}
	return true
}

func ignoreTheRemainingTokens(p *parser) bool {
	return true
}

const whitespaceOrNUL = whitespace + "\x00"

// Section 13.2.6.5
func parseForeignContent(p *parser) bool {
	switch p.tok.Type {
	case Text:
		if p.framesetOK {
			p.framesetOK = len(bytes.TrimLeft(p.tok.Text, whitespaceOrNUL)) == 0
		}
		// tok.Text aliases scratch, so the replacement goes to the other
		// buffer.
		p.scratch2 = replaceNUL(p.scratch2[:0], p.tok.Text)
		p.tok.Text = p.scratch2
		p.addText(p.tok.Text)
	case Comment:
		p.addChild(&html.Node{
			Type: html.CommentNode,
			Data: p.tok.Data,
		})
	case StartTag:
		b := breakout[p.tok.Data]
		if p.tok.DataAtom == a.Font {
		loop:
			for _, attr := range p.tok.Attr {
				switch attr.Key {
				case "color", "face", "size":
					b = true
					break loop
				}
			}
		}
		if b {
			for i := len(p.oe) - 1; i >= 0; i-- {
				n := p.oe[i]
				if n.Namespace == "" || htmlIntegrationPoint(n) || mathMLTextIntegrationPoint(n) {
					p.oe = p.oe[:i+1]
					break
				}
			}
			return p.im(p)
		}
		current := p.adjustedCurrentNode()
		switch current.Namespace {
		case "math":
			adjustAttributeNames(p.tok.Attr, mathMLAttributeAdjustments)
		case "svg":
			// Adjust SVG tag names. The tokenizer lower-cases tag names, but
			// SVG wants e.g. "foreignObject" with a capital second "O".
			if x := svgTagNameAdjustments[p.tok.Data]; x != "" {
				p.tok.DataAtom = a.Lookup([]byte(x))
				p.tok.Data = x
			}
			adjustAttributeNames(p.tok.Attr, svgAttributeAdjustments)
		default:
			panic("html: bad parser state: unexpected namespace")
		}
		adjustForeignAttributes(p.tok.Attr)
		namespace := current.Namespace
		p.addElement()
		p.top().Namespace = namespace
		if namespace != "" {
			// Don't let the tokenizer go into raw text mode in foreign content
			// (e.g. in an SVG <title> tag).
			p.r.NextIsNotRawText()
		}
		if p.hasSelfClosingToken {
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
		}
	case EndTag:
		if strings.EqualFold(p.oe[len(p.oe)-1].Data, p.tok.Data) {
			p.oe = p.oe[:len(p.oe)-1]
			return true
		}
		for i := len(p.oe) - 1; i >= 0; i-- {
			if strings.EqualFold(p.oe[i].Data, p.tok.Data) {
				p.oe = p.oe[:i]
				return true
			}
			if i > 0 && p.oe[i-1].Namespace == "" {
				break
			}
		}
		return p.im(p)
	default:
		// Ignore the token.
	}
	return true
}

// Section 12.2.4.2.
func (p *parser) adjustedCurrentNode() *html.Node {
	if len(p.oe) == 1 && p.fragment && p.context != nil {
		return p.context
	}
	return p.oe.top()
}

// Section 12.2.6.
func (p *parser) inForeignContent() bool {
	if len(p.oe) == 0 {
		return false
	}
	n := p.adjustedCurrentNode()
	if n.Namespace == "" {
		return false
	}
	if mathMLTextIntegrationPoint(n) {
		if p.tok.Type == StartTag && p.tok.DataAtom != a.Mglyph && p.tok.DataAtom != a.Malignmark {
			return false
		}
		if p.tok.Type == Text {
			return false
		}
	}
	if n.Namespace == "math" && n.DataAtom == a.AnnotationXml && p.tok.Type == StartTag && p.tok.DataAtom == a.Svg {
		return false
	}
	if htmlIntegrationPoint(n) && (p.tok.Type == StartTag || p.tok.Type == Text) {
		return false
	}
	if p.tok.Type == EOF {
		return false
	}
	return true
}

// parseImpliedToken parses a token as though it had appeared in the parser's
// input.
func (p *parser) parseImpliedToken(t Kind, dataAtom a.Atom, data string) {
	realToken, selfClosing := p.tok, p.hasSelfClosingToken
	p.tok = token{
		Type:     t,
		DataAtom: dataAtom,
		Data:     data,
	}
	p.hasSelfClosingToken = false
	p.parseCurrentToken()
	p.tok, p.hasSelfClosingToken = realToken, selfClosing
}

// parseCurrentToken runs the current token through the parsing routines
// until it is consumed.
func (p *parser) parseCurrentToken() {
	if p.tok.Type == SelfClosingTag {
		p.hasSelfClosingToken = true
		p.tok.Type = StartTag
	}

	consumed := false
	for !consumed {
		if p.inForeignContent() {
			consumed = parseForeignContent(p)
		} else {
			consumed = p.im(p)
		}
	}

	if p.hasSelfClosingToken {
		// This is a parse error, but ignore it.
		p.hasSelfClosingToken = false
	}
}

func (p *parser) parse() (err error) {
	defer func() {
		if x := recover(); x != nil {
			if e, ok := x.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", x)
			}
		}
	}()
	// Iterate until EOF. Any other error will cause an early return.
	for {
		// CDATA sections are allowed only in foreign content.
		n := p.oe.top()
		p.r.AllowCDATA(n != nil && n.Namespace != "")
		// Read and parse the next token.
		k, err := p.r.Next()
		if err != nil {
			return err
		}
		p.readToken(k)
		p.parseCurrentToken()
		if p.err != nil {
			return p.err
		}
		if k == EOF {
			break
		}
	}
	p.flushText()
	return nil
}

// readToken makes tok from the token the reader is on: the name of a tag as
// an atom's string or one made once, its attributes as a slice of the size
// they need with their names and values as strings, the text of a Text token
// decoded into scratch.
func (p *parser) readToken(k Kind) {
	p.tok = token{Type: k}
	switch k {
	case StartTag, SelfClosingTag, EndTag:
		name := p.r.Name()
		if atom := a.Lookup(name); atom != 0 {
			p.tok.DataAtom, p.tok.Data = atom, atom.String()
		} else {
			p.tok.Data = p.intern(name)
		}
		if k != EndTag {
			if n := p.r.attrCount(); n > 0 {
				attr := make([]html.Attribute, 0, n)
				for {
					key, val, ok := p.r.NextAttr()
					if !ok {
						break
					}
					attr = append(attr, html.Attribute{Key: p.attrKey(key), Val: p.attrValue(val)})
				}
				p.tok.Attr = attr
			}
		}
	case Text:
		p.scratch = p.r.Text().AppendTo(p.scratch[:0])
		p.tok.Text = p.scratch
	case Comment, Doctype:
		p.tok.Data = p.r.Text().String()
	}
}

// intern returns b as a string, the one already made for these bytes when
// there is one. Tag and attribute names that are not atoms repeat; values
// are not interned, since on real pages the map costs more than the
// strings it saves.
func (p *parser) intern(b []byte) string {
	if s, ok := p.strs[string(b)]; ok {
		return s
	}
	if p.strs == nil {
		p.strs = make(map[string]string)
	}
	s := string(b)
	p.strs[s] = s
	return s
}

func (p *parser) attrKey(key []byte) string {
	if atom := a.Lookup(key); atom != 0 {
		return atom.String()
	}
	return p.intern(key)
}

func (p *parser) attrValue(v Value) string {
	return v.String()
}

// dropNUL returns b without its NUL bytes, in place.
func dropNUL(b []byte) []byte {
	i := bytes.IndexByte(b, 0)
	if i < 0 {
		return b
	}
	out := b[:i]
	for _, c := range b[i:] {
		if c != 0 {
			out = append(out, c)
		}
	}
	return out
}

// replaceNUL appends b to dst with each NUL byte replaced by U+FFFD.
func replaceNUL(dst, b []byte) []byte {
	for _, c := range b {
		if c == 0 {
			dst = append(dst, "\uFFFD"...)
		} else {
			dst = append(dst, c)
		}
	}
	return dst
}

// keepWhitespace returns the whitespace bytes of b, in place.
func keepWhitespace(b []byte) []byte {
	out := b[:0]
	for _, c := range b {
		switch c {
		case ' ', '\t', '\n', '\f', '\r':
			out = append(out, c)
		}
	}
	return out
}

// Parse returns the parse tree for the HTML read from src, as
// golang.org/x/net/html.Parse builds it: a tree of html.Node that the
// tree construction rules of HTML5 make, with the elements the rules
// imply and without the tags they drop. The input is assumed to be UTF-8.
//
// Parse refuses a document whose elements nest deeper than
// Options.MaxDepth with ErrTooDeep, and one with a token longer than
// Options.MaxBufferBytes with ErrTooLarge.
func Parse(src io.Reader, opts Options) (*html.Node, error) {
	return parseWith(NewReader(src, opts), opts)
}

// ParseBytes is Parse over data, read in place and not modified.
func ParseBytes(data []byte, opts Options) (*html.Node, error) {
	return parseWith(NewBytesReader(data, opts), opts)
}

func newParser(r *Reader, opts Options) *parser {
	p := &parser{
		r:          r,
		doc:        &html.Node{Type: html.DocumentNode},
		scripting:  !opts.NoScripting,
		framesetOK: true,
		im:         initialIM,
		maxDepth:   opts.MaxDepth,
	}
	if p.maxDepth <= 0 {
		p.maxDepth = defaultMaxDepth
	}
	return p
}

// defaultMaxDepth is the bound on open elements when Options.MaxDepth is
// zero, the bound golang.org/x/net/html has.
const defaultMaxDepth = 512

func parseWith(r *Reader, opts Options) (*html.Node, error) {
	p := newParser(r, opts)
	if err := p.parse(); err != nil {
		return nil, err
	}
	return p.doc, nil
}

// ParseFragment parses a fragment of HTML and returns the nodes that were
// found. If the fragment is the InnerHTML for an existing element, pass
// that element in context. It has the same intricacies as Parse.
func ParseFragment(src io.Reader, context *html.Node, opts Options) ([]*html.Node, error) {
	opts.Fragment = ""
	if context != nil {
		if context.Type != html.ElementNode {
			return nil, errors.New("html: ParseFragment of non-element Node")
		}
		// The next check isn't just context.DataAtom.String() == context.Data because
		// it is valid to pass an element whose tag isn't a known atom. For example,
		// DataAtom == 0 and Data = "tagfromthefuture" is perfectly consistent.
		if context.DataAtom != a.Lookup([]byte(context.Data)) {
			return nil, fmt.Errorf("html: inconsistent Node: DataAtom=%q, Data=%q", context.DataAtom, context.Data)
		}
		if context.Namespace == "" {
			opts.Fragment = context.DataAtom.String()
		}
	}
	p := newParser(NewReader(src, opts), opts)
	p.fragment = true
	p.context = context

	root := &html.Node{
		Type:     html.ElementNode,
		DataAtom: a.Html,
		Data:     a.Html.String(),
	}
	p.doc.AppendChild(root)
	p.oe = nodeStack{root}
	if context != nil && context.DataAtom == a.Template {
		p.templateStack = append(p.templateStack, inTemplateIM)
	}
	p.resetInsertionMode()

	for n := context; n != nil; n = n.Parent {
		if n.Type == html.ElementNode && n.DataAtom == a.Form {
			p.form = n
			break
		}
	}

	if err := p.parse(); err != nil {
		return nil, err
	}

	parent := p.doc
	if context != nil {
		parent = root
	}

	var result []*html.Node
	for c := parent.FirstChild; c != nil; {
		next := c.NextSibling
		parent.RemoveChild(c)
		result = append(result, c)
		c = next
	}
	return result, nil
}
