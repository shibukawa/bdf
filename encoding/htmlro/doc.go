// Package htmlro is a read-only pull tokenizer for HTML and, on top of it,
// a parser that builds the tree golang.org/x/net/html builds, with a third
// of the allocations.
//
// The tokenizer is that of golang.org/x/net/html, carried over state for
// state so that the two tokenize alike, and what it returns is read in
// place: every name, attribute and text it hands out is a slice into its
// own buffer, valid until the next call to Next, and a caller copies only
// what it keeps. In steady state it allocates nothing.
//
// # Reading
//
// Next advances one token at a time. On a StartTag the caller reads the
// name and the attributes it wants, as the tag was indexed when it was
// scanned; on a Text, Comment or Doctype token the content, as a Value
// that decodes on demand:
//
//	r := htmlro.NewReader(src, htmlro.Options{})
//	for {
//		k, err := r.Next()
//		if err != nil {
//			return err
//		}
//		if k == htmlro.EOF {
//			break
//		}
//		if k == htmlro.StartTag && r.NameIs("a") {
//			if href, ok := r.Attr("href"); ok {
//				links = append(links, href.String())
//			}
//		}
//	}
//
// # Trees
//
// Parse, ParseBytes and ParseFragment run the tree construction stage of
// HTML5 over the tokens, as golang.org/x/net/html does, into that package's
// html.Node, so that a tree reads, renders and is walked with that package
// and the libraries built on it. The parser makes the strings a node keeps
// once each, interns the tag and attribute names that are not atoms, and
// gathers the text of a text node in one buffer before it becomes a
// string. The html5lib test suite passes as it does there, and on real
// pages the trees are the same, node for node.
//
// # Bounds
//
// A token may be as long as the input unless Options.MaxBufferBytes bounds
// it, and elements nest at most Options.MaxDepth deep, 512 by default. A
// document that needs more is refused, with ErrTooLarge or ErrTooDeep.
package htmlro
