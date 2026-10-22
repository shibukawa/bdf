// Package tag reads the tags of form templates: the marks in the text of a
// Word, Excel or PowerPoint file that data replaces when the template is
// filled (package form).
//
// A tag is written between double braces, in the manner of Mustache:
//
//	{{name}}                 a value; {{a.b}} a value of a value, {{.}} the item itself
//	{{total | fmt:"#,##0"}}  a value in a number format (the format codes of Excel)
//	{{#items}} … {{/items}}  a section: once for each item of a list, or when the value is true
//	{{^items}} … {{/items}}  an inverted section: when the list is empty or the value is false
//	{{%photo}}               a picture ({{%photo | fit:cover}})
//	{{%sign | signature}}    a picture that a person draws: a signature ({{%seal | seal}} the picture of a seal)
//	{{%url | qr}}            a QR code of a value ({{%url | qr:H}}, with the level of error correction)
//	{{%jan | barcode:ean13}} a barcode of a value (the symbologies of image/barcode)
//
// In a template of plain text (ParseText) a value may take a number of
// columns: {{name | width:20}} (on the left of them), {{no | right:8}},
// {{title | center:30}}, with | cut or | wrap for a value that is wider.
//
// A section may say how it repeats what a template has between its tags
// ({{#items | pages}}; see Token.Mode). The braces, the marks after them
// and the quotes may be the full-width or typographic ones that input
// methods and word processors put in their place. \{{ is {{ as text.
//
// Tags never run code: there are no expressions, only names looked up in
// the data. A template is as untrusted as any other input.
package tag

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shibukawa/bdf/image/barcode"
)

// Kind is the kind of a token.
type Kind uint8

const (
	// Text is text outside tags.
	Text Kind = iota
	// Value is {{name}}.
	Value
	// Open is {{#name}}.
	Open
	// Inverted is {{^name}}.
	Inverted
	// Close is {{/name}}.
	Close
	// Image is {{%name}}.
	Image
)

// How a section repeats the rows (or other units) of the template between
// its tags (Token.Mode).
const (
	// ModeRows copies the unit once for each item, moving down what is
	// below. It is the default of a section whose tags are in the same
	// row. ("block" is another name for it.)
	ModeRows = "rows"
	// ModeSlots fills the rows that are there, an item each, and leaves the
	// others empty: nothing moves. It is the default of a section whose
	// tags are in different rows.
	ModeSlots = "slots"
	// ModePages copies the unit once for each item as ModeRows does, each
	// copy on a printed page of its own.
	ModePages = "pages"
)

// Alignments of values in their columns (Token.Align), in templates of
// text.
const (
	AlignLeft   = "left"
	AlignRight  = "right"
	AlignCenter = "center"
)

// What is done with a value wider than its columns (Token.Overflow), in
// templates of text.
const (
	// OverflowKeep writes it whole, with a warning: the line is longer.
	OverflowKeep = ""
	// OverflowCut writes the columns it has room for.
	OverflowCut = "cut"
	// OverflowWrap writes the rest of it on the lines that follow, in the
	// same columns.
	OverflowWrap = "wrap"
)

// MaxWidth bounds the columns of a value (Token.Width).
const MaxWidth = 4096

// Fits of pictures (Token.Fit).
const (
	FitContain = "contain"
	FitCover   = "cover"
	FitStretch = "stretch"
)

// Inputs of pictures (Token.Input, Field.Input): what a person gives for
// the picture, which a page that asks for the data asks with. Either is a
// picture as any other is.
const (
	// InputSignature is a signature that the person draws by hand.
	InputSignature = "signature"
	// InputSeal is the picture of a seal (a hanko): a picture that the
	// person picks, of which the white is taken out.
	InputSeal = "seal"
)

// Token is a piece of template text: text, or a tag.
type Token struct {
	Kind Kind
	// Text is the text of a Text token.
	Text string
	// Name is the name of a tag as written ("customer.name", "."), and
	// Path its parts. A Close may have none: it closes the open section.
	Name string
	Path []string
	// Format is the number format of a Value (fmt:"…"); HasFormat tells
	// an empty one from none.
	Format    string
	HasFormat bool
	// Width is the columns a Value takes in a template of text (width:n,
	// right:n, center:n; 0 for the columns of its value), Align how it is
	// placed in them (AlignLeft …), and Overflow what is done with a
	// value wider than them (OverflowKeep …). Only ParseText reads them.
	Width    int
	Align    string
	Overflow string
	// Mode is how an Open repeats ("" for the default of where it is).
	Mode string
	// Fit is how an Image fills its place ("" is FitContain).
	Fit string
	// Input is what a person gives for the picture of an Image
	// (InputSignature, InputSeal; "" for any picture).
	Input string
	// Code is the symbology of an Image that is a QR code or a barcode of
	// its value, not a picture that the value names (barcode.QR,
	// barcode.EAN13 …; "" for a picture). Level is the level of error
	// correction of a QR code ("" for the default), and Quiet the modules
	// of its quiet zone when HasQuiet is set (else that of the symbology).
	Code     string
	Level    string
	Quiet    int
	HasQuiet bool
	// Pos and End are the bytes of the token in the text.
	Pos, End int
}

// Error is a mistake in the tags of a text, at a byte of it.
type Error struct {
	Pos int
	Msg string
}

func (e *Error) Error() string { return e.Msg }

func errorf(pos int, format string, args ...any) *Error {
	return &Error{Pos: pos, Msg: fmt.Sprintf(format, args...)}
}

// The braces: ASCII and full-width.
var (
	openers = []string{"{{", "｛｛"}
	closers = []string{"}}", "｝｝"}
)

// MaxTagLen bounds the bytes of a tag between its braces.
const MaxTagLen = 1024

// Has reports whether text may hold a tag.
func Has(s string) bool {
	return strings.Contains(s, openers[0]) || strings.Contains(s, openers[1])
}

// Parse splits text into text and tags.
func Parse(s string) ([]Token, error) { return parse(s, false) }

// ParseText splits the text of a template of plain text into text and
// tags: those of Parse, whose values may also take a number of columns
// (width:n, right:n, center:n, with cut or wrap).
func ParseText(s string) ([]Token, error) { return parse(s, true) }

func parse(s string, text bool) ([]Token, error) {
	var out []Token
	var lit strings.Builder
	litPos := 0
	flush := func(end int) {
		if lit.Len() > 0 {
			out = append(out, Token{Kind: Text, Text: lit.String(), Pos: litPos, End: end})
			lit.Reset()
		}
	}
	for i := 0; i < len(s); {
		op := prefix(s[i:], openers)
		if op == "" {
			if s[i] == '\\' && prefix(s[i+1:], openers) != "" {
				// an escaped brace is text
				if lit.Len() == 0 {
					litPos = i
				}
				op = prefix(s[i+1:], openers)
				lit.WriteString(op)
				i += 1 + len(op)
				continue
			}
			if lit.Len() == 0 {
				litPos = i
			}
			_, n := utf8.DecodeRuneInString(s[i:])
			lit.WriteString(s[i : i+n])
			i += n
			continue
		}
		flush(i)
		start := i + len(op)
		end, next := closing(s, start)
		if end < 0 {
			return nil, errorf(i, "the tag %s is not closed with }}", clip(s[i:]))
		}
		if end-start > MaxTagLen {
			return nil, errorf(i, "the tag %s is longer than %d bytes", clip(s[i:]), MaxTagLen)
		}
		tok, err := parseTag(s[start:end], i, text)
		if err != nil {
			return nil, err
		}
		tok.Pos, tok.End = i, next
		out = append(out, tok)
		i = next
	}
	flush(len(s))
	return out, nil
}

// prefix returns the one of the strings that s starts with, or "".
func prefix(s string, of []string) string {
	for _, p := range of {
		if strings.HasPrefix(s, p) {
			return p
		}
	}
	return ""
}

// closing finds the braces that close a tag whose content starts at start:
// where the content ends and where the text goes on (-1 when none). Braces
// in a quoted argument are not the tag's.
func closing(s string, start int) (end, next int) {
	var quote rune
	for i := start; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			}
		case closeQuote(r) != 0:
			quote = closeQuote(r)
		default:
			if c := prefix(s[i:], closers); c != "" {
				return i, i + len(c)
			}
			if prefix(s[i:], openers) != "" {
				return -1, -1 // another tag starts: this one was not closed
			}
		}
		i += n
	}
	return -1, -1
}

// closeQuote returns the quote that closes a quoted argument that r opens,
// or 0 when r opens none: the ASCII quotes and the typographic ones that
// word processors put in their place.
func closeQuote(r rune) rune {
	switch r {
	case '"', '\'', '＂', '＇':
		return r
	case '“', '”', '„', '‟':
		return '”'
	case '‘', '’', '‚', '‛':
		return '’'
	case '「':
		return '」'
	}
	return 0
}

// clip shortens the text of a tag for a message.
func clip(s string) string {
	const n = 24
	if len(s) <= n {
		return s
	}
	i := n
	for i > 0 && !utf8.RuneStart(s[i]) {
		i--
	}
	return s[:i] + "…"
}

// space reports the white space around the parts of a tag.
func space(r rune) bool { return unicode.IsSpace(r) || r == '　' }

// parseTag reads what is between the braces of a tag at pos.
func parseTag(body string, pos int, text bool) (Token, error) {
	tok := Token{Kind: Value}
	body = strings.TrimFunc(body, space)
	if r, n := utf8.DecodeRuneInString(body); n > 0 {
		kind := Value
		switch r {
		case '#', '＃':
			kind = Open
		case '^', '＾':
			kind = Inverted
		case '/', '／':
			kind = Close
		case '%', '％':
			kind = Image
		}
		if kind != Value {
			tok.Kind, body = kind, strings.TrimFunc(body[n:], space)
		}
	}
	parts := splitFilters(body)
	tok.Name = strings.TrimFunc(parts[0], space)
	switch {
	case tok.Name == "" && tok.Kind == Close:
	case tok.Name == "":
		return tok, errorf(pos, "a tag without a name")
	case tok.Name == ".":
		tok.Path = []string{"."}
	default:
		for _, p := range strings.Split(tok.Name, ".") {
			p = strings.TrimFunc(p, space)
			if p == "" {
				return tok, errorf(pos, "the name %q of a tag has an empty part", tok.Name)
			}
			if strings.ContainsAny(p, "{}#^/%\"'") {
				return tok, errorf(pos, "the name %q of a tag has a character that names do not have", tok.Name)
			}
			tok.Path = append(tok.Path, p)
		}
		tok.Name = strings.Join(tok.Path, ".")
	}
	for _, f := range parts[1:] {
		name, arg, hasArg := cutFilter(f)
		switch {
		case name == "fmt" && tok.Kind == Value && hasArg:
			tok.Format, tok.HasFormat = arg, true
		case name == "fit" && tok.Kind == Image && (arg == FitContain || arg == FitCover || arg == FitStretch):
			tok.Fit = arg
		case (name == InputSignature || name == InputSeal) && tok.Kind == Image && !hasArg && tok.Input == "":
			tok.Input = name
		case name == "qr" && tok.Kind == Image && (!hasArg || len(arg) == 1 && strings.Contains("LMQH", strings.ToUpper(arg))):
			tok.Code, tok.Level = barcode.QR, strings.ToUpper(arg)
		case name == "barcode" && tok.Kind == Image && hasArg && barcode.Kind(arg) != "":
			tok.Code = barcode.Kind(arg)
		case name == "quiet" && tok.Kind == Image && hasArg && quietArg(arg) >= 0:
			tok.Quiet, tok.HasQuiet = quietArg(arg), true
		case text && tok.Kind == Value && hasArg && (name == "width" || name == AlignRight || name == AlignCenter) && tok.Width == 0:
			n, err := strconv.Atoi(arg)
			if err != nil || n < 1 || n > MaxWidth {
				return tok, errorf(pos, "the tag {{%s}} takes columns from 1 to %d: %s", tok.Name, MaxWidth, strings.TrimFunc(f, space))
			}
			tok.Width, tok.Align = n, AlignLeft
			if name != "width" {
				tok.Align = name
			}
		case text && tok.Kind == Value && !hasArg && (name == OverflowCut || name == OverflowWrap) && tok.Overflow == "":
			tok.Overflow = name
		case (tok.Kind == Open || tok.Kind == Inverted) && !hasArg && (name == ModeRows || name == ModeSlots || name == ModePages || name == "block"):
			if name == "block" {
				name = ModeRows
			}
			tok.Mode = name
		default:
			return tok, errorf(pos, "the tag {{%s}} has an option it does not take: %s", tok.Name, strings.TrimFunc(f, space))
		}
	}
	if tok.Input != "" && tok.Code != "" {
		return tok, errorf(pos, "the tag {{%%%s}} is a code; it does not take %s", tok.Name, tok.Input)
	}
	if tok.Overflow != "" && tok.Width == 0 {
		return tok, errorf(pos, "the tag {{%s}} takes %s with its columns only (width:n, right:n or center:n)", tok.Name, tok.Overflow)
	}
	return tok, nil
}

// quietArg reads the width of a quiet zone: 0 to 64 modules (-1 for what
// is none).
func quietArg(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 || n > 64 {
		return -1
	}
	return n
}

// Symbol makes the QR code or the barcode that an Image tag asks for
// (Code) of a value of the data: a text, or a number (written without an
// exponent; a code that starts with zeros is a text in the data).
func (t Token) Symbol(v any) (*barcode.Symbol, error) {
	var text string
	switch x := v.(type) {
	case string:
		text = x
	case float64:
		text = strconv.FormatFloat(x, 'f', -1, 64)
	default:
		return nil, fmt.Errorf("the value of %q is not a text nor a number", t.Name)
	}
	o := barcode.Options{Level: t.Level}
	if t.HasQuiet {
		if o.Quiet = t.Quiet; t.Quiet == 0 {
			o.Quiet = barcode.NoQuiet
		}
	}
	return barcode.Encode(t.Code, text, &o)
}

// splitFilters splits the content of a tag at the bars outside quotes.
func splitFilters(s string) []string {
	var out []string
	var quote rune
	start := 0
	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			}
		case closeQuote(r) != 0:
			quote = closeQuote(r)
		case r == '|' || r == '｜':
			out = append(out, s[start:i])
			start = i + n
		}
		i += n
	}
	return append(out, s[start:])
}

// cutFilter reads an option of a tag: its name and its argument, unquoted.
func cutFilter(s string) (name, arg string, hasArg bool) {
	s = strings.TrimFunc(s, space)
	i := strings.IndexAny(s, ":：")
	if i < 0 {
		return strings.ToLower(s), "", false
	}
	_, n := utf8.DecodeRuneInString(s[i:])
	name = strings.ToLower(strings.TrimFunc(s[:i], space))
	arg = strings.TrimFunc(s[i+n:], space)
	if r, n := utf8.DecodeRuneInString(arg); n > 0 && closeQuote(r) != 0 {
		end, m := utf8.DecodeLastRuneInString(arg)
		if len(arg) > n && end == closeQuote(r) {
			arg = arg[n : len(arg)-m]
		}
	}
	return name, arg, true
}

// Kinds of fields (Field.Kind).
const (
	FieldValue   = "value"
	FieldSection = "section"
	FieldImage   = "image"
)

// Field is a name that the tags of a template look up in the data: what a
// template asks of the data that fills it.
type Field struct {
	// Name is the name as the tags write it ("customer.name").
	Name string `json:"name"`
	// Kind is FieldValue, FieldSection or FieldImage.
	Kind string `json:"kind"`
	// In names the sections the tags are in, outermost first: the name is
	// looked up in their items first.
	In []string `json:"in,omitempty"`
	// Format is the number format of a value, when its tags give one.
	Format string `json:"format,omitempty"`
	// Mode is how a section repeats (ModeRows, ModeSlots or ModePages; ""
	// for a section inside a text), and Slots the items a section of
	// ModeSlots has room for.
	Mode  string `json:"mode,omitempty"`
	Slots int    `json:"slots,omitempty"`
	// Inverted is set for a name that only inverted sections look up.
	Inverted bool `json:"inverted,omitempty"`
	// Code is the symbology of an image that is a QR code or a barcode of
	// the value ("qr", "ean13" …): the value is then what the symbol
	// says, not the name of a picture.
	Code string `json:"code,omitempty"`
	// Input is what a person gives for a picture (InputSignature,
	// InputSeal; "" for any picture).
	Input string `json:"input,omitempty"`
	// At is where the tags are ("Sheet1!B5").
	At []string `json:"at"`

	// What a template whose fields are more than tags says of a value (a
	// field of the form of a PDF, a field of a template of pdfme):

	// Type is the type of the value: "string", "number", "boolean",
	// "date" or "date-time" ("": the tags do not tell).
	Type string `json:"type,omitempty"`
	// Options are the values that a value is one of, and Labels what a
	// person is shown for each (the same number of them, or none).
	Options []string `json:"options,omitempty"`
	Labels  []string `json:"labels,omitempty"`
	// Label names the field for a person, and Default is the value that
	// the template has for it.
	Label   string `json:"label,omitempty"`
	Default any    `json:"default,omitempty"`
	// MaxLen is the most characters that a text has room for (0: any
	// number), and Multiline says that it may be of several lines.
	MaxLen    int  `json:"maxLength,omitempty"`
	Multiline bool `json:"multiline,omitempty"`
}
