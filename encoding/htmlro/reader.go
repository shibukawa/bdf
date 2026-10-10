package htmlro

import (
	"errors"
	"io"
	"strconv"
)

// Kind is the kind of token the reader is positioned on.
type Kind uint8

const (
	// None is the kind before the first Next and after an error.
	None Kind = iota
	// StartTag is an opening tag, <a href="x">. Its attributes are read
	// with NextAttr or Attr.
	StartTag
	// EndTag is a closing tag, </a>.
	EndTag
	// SelfClosingTag is a tag that closes itself, <br/>. The tree
	// construction rules decide whether that means anything.
	SelfClosingTag
	// Text is character data, as written: references are not decoded and
	// whitespace is not trimmed. The text of a script or style element is
	// one Text token.
	Text
	// Comment is the body of a comment, <!--x-->, and also what the
	// tokenizer makes of bogus markup: <?x>, </>, <!x>.
	Comment
	// Doctype is a document type declaration, <!DOCTYPE x>, as the text
	// after the keyword.
	Doctype
	// EOF is reported once at the end of the input.
	EOF
)

var kindNames = [...]string{"None", "StartTag", "EndTag", "SelfClosingTag", "Text", "Comment", "Doctype", "EOF"}

func (k Kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "Kind(" + strconv.Itoa(int(k)) + ")"
}

// Options bounds a Reader and a parse. The zero value selects the defaults
// noted on each field.
type Options struct {
	// BufferSize is the initial buffer of a Reader over a stream, 64 KiB by
	// default. It grows only when one token does not fit, and only up to
	// MaxBufferBytes.
	BufferSize int
	// MaxBufferBytes bounds the buffer, and so the largest single token: a
	// tag with its attributes, a run of text, the whole text of a script or
	// style element. 0, the default, leaves it unbounded, as the HTML
	// parser of golang.org/x/net does; a reader of untrusted input sets it.
	MaxBufferBytes int
	// MaxDepth bounds the stack of open elements of a parse, 512 by
	// default, which is also the bound of golang.org/x/net/html.
	MaxDepth int
	// NoScripting parses as a browser with scripting disabled does: the
	// content of noscript elements is parsed as markup.
	NoScripting bool
	// Fragment is the tag of the element whose inner HTML a Reader reads,
	// "script" or "textarea" for instance, so that the first token is read
	// as that element's text. ParseFragment sets it from its context.
	Fragment string
}

const defaultBufferSize = 64 << 10

var (
	// ErrTooLarge is returned when a token does not fit in
	// Options.MaxBufferBytes.
	ErrTooLarge = errors.New("html: token exceeds MaxBufferBytes")
	// ErrTooDeep is returned by a parse whose open elements nest deeper
	// than Options.MaxDepth.
	ErrTooDeep = errors.New("html: elements nest deeper than MaxDepth")
)

// span is a range of bytes in a Reader's buffer. The start is inclusive,
// the end is exclusive.
type span struct {
	start, end int
}

// rawKind is the element whose text the reader is inside: the ones whose
// content is not markup, and the end tag that ends it.
type rawKind uint8

const (
	rawNone rawKind = iota
	rawIframe
	rawNoembed
	rawNoframes
	rawNoscript
	rawPlaintext
	rawScript
	rawStyle
	rawTextarea
	rawTitle
	rawXmp
)

var rawTagNames = [...]string{"", "iframe", "noembed", "noframes", "noscript", "plaintext", "script", "style", "textarea", "title", "xmp"}

// rawKindOf returns the rawKind of a lower-case tag name, rawNone when the
// element's content is markup.
func rawKindOf(name string) rawKind {
	for k, n := range rawTagNames {
		if k != 0 && n == name {
			return rawKind(k)
		}
	}
	return rawNone
}

// Reader reads HTML tokens from an io.Reader or from a byte slice.
//
// Every slice a Reader returns aliases its buffer and is valid until the
// next call to Next. Callers that keep a name or a value copy it.
//
// A Reader is not safe for concurrent use. Reuse one across documents with
// Reset rather than allocating one per document.
type Reader struct {
	src io.Reader
	// buf[raw.start:raw.end] holds the raw bytes of the current token;
	// buf[raw.end:] is buffered input that will yield future tokens.
	buf   []byte
	raw   span
	owned bool  // buf is the reader's own: it may be compacted and grown
	base  int64 // bytes discarded before buf[0]
	// readErr is the error the source returned, kept apart from err
	// because a Read may return bytes and an error together.
	readErr error
	// err is the first error met while tokenizing. Once set it stays set.
	err      error
	maxBuf   int
	bufSize  int
	fragment rawKind
	kind     Kind
	// buf[data.start:data.end] holds the current token's data: a text
	// token's text, a tag token's name, a comment's body.
	data span
	// pendingAttr is the attribute being scanned; when complete it is
	// pushed onto attr. nAttr counts the attributes NextAttr returned.
	pendingAttr [2]span
	attr        [][2]span
	nAttr       int
	// rawTag is the element that closes the next token, when the next
	// token is the text of a raw text or RCDATA element.
	rawTag rawKind
	// textIsRaw is whether the current text token's data is not escaped.
	textIsRaw bool
	// convertNUL is whether NUL bytes in the current token's data become
	// U+FFFD.
	convertNUL bool
	// allowCDATA is whether CDATA sections are recognized.
	allowCDATA bool
	// alt holds the names of the current token that had to be folded or
	// cleaned of NULs, which cannot be done in place.
	alt []byte
}

// NewReader returns a Reader over src.
func NewReader(src io.Reader, opts Options) *Reader {
	r := &Reader{}
	r.init(opts)
	r.Reset(src)
	return r
}

// NewBytesReader returns a Reader over data, which it reads in place and
// does not modify.
func NewBytesReader(data []byte, opts Options) *Reader {
	r := &Reader{}
	r.init(opts)
	r.ResetBytes(data)
	return r
}

func (r *Reader) init(opts Options) {
	r.bufSize = opts.BufferSize
	if r.bufSize <= 0 {
		r.bufSize = defaultBufferSize
	}
	r.maxBuf = opts.MaxBufferBytes
	r.fragment = fragmentRawKind(opts.Fragment)
}

// fragmentRawKind is the state the reader starts in for the inner HTML of
// the element named tag. As the "Parsing HTML Fragments" section of the
// specification allows an implementation that does not report errors, the
// raw text and script data states are the plaintext state here: there is
// no appropriate end tag in the fragment case, so they read the same.
func fragmentRawKind(tag string) rawKind {
	var tmp [16]byte
	if len(tag) > len(tmp) {
		return rawNone
	}
	switch k := rawKindOf(string(lower(append(tmp[:0], tag...)))); k {
	case rawTitle, rawTextarea:
		return k
	case rawNone:
		return rawNone
	default:
		return rawPlaintext
	}
}

// Reset makes the Reader read src from the start, keeping its buffer.
func (r *Reader) Reset(src io.Reader) {
	r.src = src
	if r.owned && r.buf != nil {
		r.buf = r.buf[:0]
	} else {
		r.buf = make([]byte, 0, r.bufSize)
		r.owned = true
	}
	r.readErr = nil
	r.reset()
}

// ResetBytes makes the Reader read data, in place, from the start.
func (r *Reader) ResetBytes(data []byte) {
	r.src = nil
	r.buf = data
	r.owned = false
	r.readErr = io.EOF
	r.reset()
}

func (r *Reader) reset() {
	r.raw = span{}
	r.base = 0
	r.err = nil
	r.kind = None
	r.data = span{}
	r.attr = r.attr[:0]
	r.nAttr = 0
	r.rawTag = r.fragment
	r.textIsRaw = false
	r.convertNUL = false
	r.allowCDATA = false
	r.alt = r.alt[:0]
}

// AllowCDATA sets whether the reader recognizes <![CDATA[x]]> as the text
// "x". Off, the default, reads it as a bogus comment. The HTML tokenizer
// recognizes CDATA only in foreign content, SVG and MathML, which the
// reader cannot know by itself; the parser sets this as it goes.
func (r *Reader) AllowCDATA(allow bool) {
	r.allowCDATA = allow
}

// NextIsNotRawText tells the reader that the next token is markup although
// the last start tag was one of an element whose content is text: a title
// inside svg, a textarea the parser does not open. The parser calls it as
// the tree construction rules require.
func (r *Reader) NextIsNotRawText() {
	r.rawTag = rawNone
}

// Kind returns the kind of the current token.
func (r *Reader) Kind() Kind { return r.kind }

// Offset returns the byte offset of the current token from the start of
// the input.
func (r *Reader) Offset() int64 { return r.base + int64(r.raw.start) }

// Next advances to the next token and returns its kind. At the end of the
// input it returns EOF with a nil error, once; a token the input ends
// inside is dropped, as the tokenizer specification has it. After an error
// it returns None and the same error.
func (r *Reader) Next() (Kind, error) {
	k := r.next()
	if k == None {
		if r.err == io.EOF {
			r.kind = EOF
			return EOF, nil
		}
		r.kind = None
		return None, r.err
	}
	r.kind = k
	return k, nil
}

// Name returns the lower-cased name of a tag token: the "img" of
// <IMG SRC="x">. A NUL in it is U+FFFD.
func (r *Reader) Name() []byte {
	switch r.kind {
	case StartTag, EndTag, SelfClosingTag:
		return r.foldName(r.buf[r.data.start:r.data.end])
	}
	return nil
}

// NameIs reports whether the current token is a tag named s, lower case.
func (r *Reader) NameIs(s string) bool {
	switch r.kind {
	case StartTag, EndTag, SelfClosingTag:
		return EqualFold(r.buf[r.data.start:r.data.end], s)
	}
	return false
}

// NextAttr returns the next attribute of a start tag not yet returned, its
// name lower-cased and its value as written, and whether there was one. A
// name that a tag repeats is returned once, with its first value.
func (r *Reader) NextAttr() (name []byte, value Value, ok bool) {
	if r.nAttr >= len(r.attr) {
		return nil, Value{}, false
	}
	switch r.kind {
	case StartTag, SelfClosingTag:
		x := r.attr[r.nAttr]
		r.nAttr++
		name = r.foldName(r.buf[x[0].start:x[0].end])
		value = Value{raw: r.buf[x[1].start:x[1].end], flags: valReader | valEntities | valAttr | valNUL}
		return name, value, true
	}
	return nil, Value{}, false
}

// Attr returns the value of the attribute named name, lower case, of a
// start tag, and whether the tag has it.
func (r *Reader) Attr(name string) (Value, bool) {
	switch r.kind {
	case StartTag, SelfClosingTag:
		for _, x := range r.attr {
			if EqualFold(r.buf[x[0].start:x[0].end], name) {
				return Value{raw: r.buf[x[1].start:x[1].end], flags: valReader | valEntities | valAttr | valNUL}, true
			}
		}
	}
	return Value{}, false
}

// Text returns the content of a Text, Comment or Doctype token, as
// written; its Value decodes as the tokenizer would.
func (r *Reader) Text() Value {
	switch r.kind {
	case Text, Comment, Doctype:
		flags := uint8(valReader)
		if r.convertNUL || r.kind == Comment {
			flags |= valNUL
		}
		if !r.textIsRaw {
			flags |= valEntities
		}
		return Value{raw: r.buf[r.data.start:r.data.end], flags: flags}
	}
	return Value{}
}

// attrCount returns how many attributes the current start tag has, repeats
// dropped.
func (r *Reader) attrCount() int {
	return len(r.attr)
}

// foldName returns b lower-cased with its NULs replaced by U+FFFD: b itself
// when that changes nothing, else a copy in alt, which lives until the next
// token.
func (r *Reader) foldName(b []byte) []byte {
	for i, c := range b {
		if 'A' <= c && c <= 'Z' || c == 0 {
			start := len(r.alt)
			r.alt = append(r.alt, b[:i]...)
			for _, c := range b[i:] {
				switch {
				case 'A' <= c && c <= 'Z':
					r.alt = append(r.alt, c+'a'-'A')
				case c == 0:
					r.alt = append(r.alt, "�"...)
				default:
					r.alt = append(r.alt, c)
				}
			}
			return r.alt[start:]
		}
	}
	return b
}

// lower lower-cases the ASCII letters of b in place.
func lower(b []byte) []byte {
	for i, c := range b {
		if 'A' <= c && c <= 'Z' {
			b[i] = c + 'a' - 'A'
		}
	}
	return b
}

// readByte returns the next byte of the input, reading from the source
// into the buffer when it is exhausted. buf[raw.start:raw.end] stays one
// contiguous slice of the bytes read so far for the current token. It sets
// err if the source returns an error, or if the token outgrows maxBuf.
// Precondition: err == nil.
func (r *Reader) readByte() byte {
	if r.raw.end >= len(r.buf) {
		// The buffer is exhausted: check whether the previous read ended
		// the input.
		if r.readErr != nil {
			r.err = r.readErr
			return 0
		}
		// Move the current token to the start of the buffer, into a buffer
		// twice as big when it takes more than half of this one, and read
		// after it.
		c := cap(r.buf)
		d := r.raw.end - r.raw.start
		var buf1 []byte
		if 2*d > c {
			buf1 = make([]byte, d, 2*c)
		} else {
			buf1 = r.buf[:d]
		}
		copy(buf1, r.buf[r.raw.start:r.raw.end])
		if x := r.raw.start; x != 0 {
			// The spans refer to the same bytes after the move.
			r.data.start -= x
			r.data.end -= x
			r.pendingAttr[0].start -= x
			r.pendingAttr[0].end -= x
			r.pendingAttr[1].start -= x
			r.pendingAttr[1].end -= x
			for i := range r.attr {
				r.attr[i][0].start -= x
				r.attr[i][0].end -= x
				r.attr[i][1].start -= x
				r.attr[i][1].end -= x
			}
			r.base += int64(x)
		}
		r.raw.start, r.raw.end, r.buf = 0, d, buf1[:d]
		var n int
		n, r.readErr = readAtLeastOneByte(r.src, buf1[d:cap(buf1)])
		if n == 0 {
			r.err = r.readErr
			return 0
		}
		r.buf = buf1[:d+n]
	}
	x := r.buf[r.raw.end]
	r.raw.end++
	if r.maxBuf > 0 && r.raw.end-r.raw.start >= r.maxBuf {
		r.err = ErrTooLarge
		return 0
	}
	return x
}

// readAtLeastOneByte reads from src so that a read cannot return (0, nil):
// it returns io.ErrNoProgress when src does so too many times in a row.
func readAtLeastOneByte(src io.Reader, b []byte) (int, error) {
	for i := 0; i < 100; i++ {
		if n, err := src.Read(b); n != 0 || err != nil {
			return n, err
		}
	}
	return 0, io.ErrNoProgress
}
