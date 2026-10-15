# encoding/htmlro

A read-only pull tokenizer for HTML, and on top of it a parser that builds
the tree `golang.org/x/net/html` builds.
The tokenizer reads names, attributes and text in place and allocates
nothing in steady state; the parser makes the strings a node keeps, once
each, and little else.

The name says what it is: HTML, read only. It never writes; `html.Render` of
`golang.org/x/net/html` writes the trees it builds.

## Reading

```go
r := htmlro.NewReader(src, htmlro.Options{})   // or htmlro.NewBytesReader(data, ...)
for {
	k, err := r.Next()
	if err != nil {
		return err
	}
	if k == htmlro.EOF {
		break
	}
	if k == htmlro.StartTag && r.NameIs("a") {
		if href, ok := r.Attr("href"); ok {
			links = append(links, href.String())
		}
	}
}
```

On a `StartTag` the caller reads `Name`, and the attributes with `Attr` by
name or all of them in order with `NextAttr`; they were indexed as the tag
was scanned, so neither copies anything. A name a tag repeats is reported
once, with its first value. On a `Text`, `Comment` or `Doctype` token,
`Text` returns the content.

Every slice the reader returns aliases its buffer and is valid until the
next `Next`. A caller that keeps a name or a value copies it.

## Values

`Attr`, `NextAttr` and `Text` return a `Value`: the bytes as written, with a
note of what decoding them means where they were read. `AppendTo`, `String`
and `Equal` decode on demand: "\r\n" and "\r" become "\n"; a NUL becomes
U+FFFD in an attribute value, a comment and the text of a script or style
element (the parser deals with the NULs of other text, as the specification
has it); and in an attribute value or text that is not raw, a character
reference becomes what it names, numeric or one of the 2,231 named ones,
the legacy references that need no semicolon included, with the rule that
keeps `&copy=1` as written in an attribute value. Content with none of
those, which is nearly all of it, is never copied: `NeedsDecoding` says.

`UnescapeString` is the same decoding for a string outside the reader,
what `golang.org/x/net/html.UnescapeString` does, so a program that has
this package needs neither that nor the standard library's `html` for it.

## Trees

```go
doc, err := htmlro.ParseBytes(data, htmlro.Options{MaxBufferBytes: len(data)})
```

`Parse`, `ParseBytes` and `ParseFragment` are `html.Parse` and
`html.ParseFragment` of `golang.org/x/net/html`, the same tree construction
stage over this tokenizer, and they build that package's `html.Node`, so a
tree reads, renders and is walked with that package and with the libraries
built on it (go-readability, cascadia, goquery). `Options.NoScripting` is
its `ParseOptionEnableScripting(false)`.

The two give the same tree. The html5lib tree construction suite passes
here as it does there (the same six cases skipped), the html5lib tokenizer
inputs and that package's own tokenizer cases tokenize the same through
both. `TestParseMatchesXNet` checks the resulting trees and rendered
markup, while `FuzzReaderMatchesTokenizer` and `FuzzParseMatchesXNet`
exercise additional inputs.

## Bounds

A token may be as long as the input unless `Options.MaxBufferBytes` bounds
it, since the text of a script element often is; a reader of untrusted
input sets it, and gets `ErrTooLarge` past it. Elements nest at most
`Options.MaxDepth` deep, 512 by default as in `golang.org/x/net/html`, and
a document that needs more is refused with `ErrTooDeep` where that package
panics.

## What it is made of

The tokenizer states (`scan.go`), the tree construction stage (`parse.go`,
`foreign.go`, `doctype.go`, `const.go`) and the test suites are those of
`golang.org/x/net/html`, carried over function for function so that the
two read alike, under its license ([LICENSE](./LICENSE)). The entity table
is generated from that package's by `gen.go`. What is new is the reader
around the states, the `Value`, and the way the parser takes its tokens.

Run `go test ./encoding/htmlro` to check tokenization, tree construction,
value decoding, reset behavior and input limits.
