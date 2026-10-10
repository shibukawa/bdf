package htmlro

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// tok is one token as the reader and the tokenizer of golang.org/x/net/html
// both report it, for comparison: the kind, the decoded data, the
// attributes with their decoded values.
type tok struct {
	Kind string
	Data string
	Attr []html.Attribute
}

// xnetTokens tokenizes input with golang.org/x/net/html.
func xnetTokens(input string, cdata bool, fragment string) ([]tok, error) {
	var z *html.Tokenizer
	if fragment != "" {
		z = html.NewTokenizerFragment(strings.NewReader(input), fragment)
	} else {
		z = html.NewTokenizer(strings.NewReader(input))
	}
	z.AllowCDATA(cdata)
	var out []tok
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if err := z.Err(); err != io.EOF {
				return out, err
			}
			return out, nil
		}
		t := z.Token()
		out = append(out, tok{Kind: tt.String(), Data: t.Data, Attr: t.Attr})
	}
}

// readerTokens reads every token of r.
func readerTokens(r *Reader) ([]tok, error) {
	var out []tok
	for {
		k, err := r.Next()
		if err != nil {
			return out, err
		}
		if k == EOF {
			return out, nil
		}
		t := tok{Kind: k.String()}
		switch k {
		case StartTag, EndTag, SelfClosingTag:
			t.Data = string(r.Name())
			for {
				name, v, ok := r.NextAttr()
				if !ok {
					break
				}
				t.Attr = append(t.Attr, html.Attribute{Key: string(name), Val: v.String()})
			}
		default:
			t.Data = r.Text().String()
		}
		out = append(out, t)
	}
}

// oneByteReader hands out one byte per Read, so that every token crosses a
// buffer boundary.
type oneByteReader struct {
	s string
	i int
}

func (r *oneByteReader) Read(p []byte) (int, error) {
	if r.i >= len(r.s) {
		return 0, io.EOF
	}
	p[0] = r.s[r.i]
	r.i++
	return 1, nil
}

// feeds are the ways an input reaches a reader: in memory, from a stream,
// and from a stream one byte at a time into a buffer that starts at one
// byte, which exercises every move and growth of the buffer.
var feeds = []struct {
	name string
	open func(input string, cdata bool, fragment string) *Reader
}{
	{"bytes", func(input string, cdata bool, fragment string) *Reader {
		r := NewBytesReader([]byte(input), Options{Fragment: fragment})
		r.AllowCDATA(cdata)
		return r
	}},
	{"stream", func(input string, cdata bool, fragment string) *Reader {
		r := NewReader(strings.NewReader(input), Options{Fragment: fragment})
		r.AllowCDATA(cdata)
		return r
	}},
	{"byte-by-byte", func(input string, cdata bool, fragment string) *Reader {
		r := NewReader(&oneByteReader{s: input}, Options{BufferSize: 1, Fragment: fragment})
		r.AllowCDATA(cdata)
		return r
	}},
}

// compareTokens checks that the reader gives the tokens the tokenizer of
// golang.org/x/net/html gives for input, through every feed.
func compareTokens(t *testing.T, input string, cdata bool, fragment string) {
	t.Helper()
	want, wantErr := xnetTokens(input, cdata, fragment)
	for _, f := range feeds {
		got, err := readerTokens(f.open(input, cdata, fragment))
		if (err != nil) != (wantErr != nil) {
			t.Errorf("%s: %q: error %v, x/net %v", f.name, input, err, wantErr)
			continue
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %q:\ngot  %#v\nwant %#v", f.name, input, got, want)
		}
	}
}

func TestReaderMatchesTokenizer(t *testing.T) {
	for _, tt := range tokenTests {
		t.Run(tt.desc, func(t *testing.T) {
			compareTokens(t, tt.html, true, "")
			compareTokens(t, tt.html, false, "")
		})
	}
}

func TestReaderMatchesTokenizerOnFragments(t *testing.T) {
	inputs := []string{
		"<!-- inside </script> -->",
		"a<b>c</textarea>d</TEXTAREA>e",
		"&amp;<i>x</title>&lt;",
		"<p>x</style>y",
		"plain </plaintext> text",
	}
	for _, in := range inputs {
		for _, ctx := range []string{"script", "style", "title", "textarea", "plaintext", "xmp", "iframe", "noembed", "noframes", "noscript", "div", "SCRIPT"} {
			compareTokens(t, in, false, ctx)
		}
	}
}

// html5libTokenizerInputs returns the inputs of the html5lib tokenizer
// tests, their escapes undone.
func html5libTokenizerInputs(t *testing.T) []string {
	var inputs []string
	files, err := filepath.Glob("testdata/html5lib-tests/tokenizer/*.test")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		var tests struct {
			Tests []struct {
				Input         string
				DoubleEscaped bool
			}
		}
		if err := json.Unmarshal(data, &tests); err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		for _, tc := range tests.Tests {
			in := tc.Input
			if tc.DoubleEscaped {
				in = unescapeUnicode(in)
			}
			inputs = append(inputs, in)
		}
	}
	return inputs
}

func TestReaderMatchesTokenizerOnHTML5Lib(t *testing.T) {
	for _, in := range html5libTokenizerInputs(t) {
		compareTokens(t, in, false, "")
	}
}

func TestReaderMatchesTokenizerOnTreeTests(t *testing.T) {
	parseTestCases(t, func(file string, i int, ta *testAttrs) {
		compareTokens(t, ta.text, false, "")
		compareTokens(t, ta.text, true, "")
	})
}

var unicodeRegexp = regexp.MustCompile(`\\u[0-9a-fA-F]{4}`)

// unescapeUnicode undoes the \uXXXX escapes of a double-escaped html5lib
// test input, encoding each code unit as UTF-8 would, lone surrogates
// included, which is what the suite means by them.
func unescapeUnicode(s string) string {
	const (
		tx    = 0x80
		t2    = 0xC0
		t3    = 0xE0
		t4    = 0xF0
		maskx = 0x3F
	)
	return unicodeRegexp.ReplaceAllStringFunc(s, func(match string) string {
		n, err := strconv.ParseInt(match[2:], 16, 32)
		if err != nil {
			panic(err)
		}
		switch i := uint32(n); {
		case i <= 0x7F:
			return string(byte(n))
		case i <= 0x7FF:
			return string([]byte{t2 | byte(n>>6), tx | byte(n)&maskx})
		case i <= 0xFFFF:
			return string([]byte{t3 | byte(n>>12), tx | byte(n>>6)&maskx, tx | byte(n)&maskx})
		default:
			return string([]byte{t4 | byte(n>>18), tx | byte(n>>12)&maskx, tx | byte(n>>6)&maskx, tx | byte(n)&maskx})
		}
	})
}

func TestValueDecoding(t *testing.T) {
	cases := []struct {
		in, attr, text string
	}{
		{"plain", "plain", "plain"},
		{"a&amp;b", "a&b", "a&b"},
		{"a&AMP;b", "a&b", "a&b"},
		{"&lt;&gt;&quot;&apos;", `<>"'`, `<>"'`},
		{"&copy;", "©", "©"},
		{"&copy", "©", "©"},                 // a legacy reference needs no semicolon
		{"&copy=1", "&copy=1", "©=1"},       // but not before '=' in an attribute value
		{"&copyx", "&copyx", "©x"},          // nor before more name characters there
		{"&notin;", "∉", "∉"},               // the longest name wins
		{"&notin", "&notin", "¬in"},         // without the semicolon, the longest legacy prefix, in text
		{"&NotEqualTilde;", "≂̸", "≂̸"},     // two code points
		{"&#65;&#x42;&#X43;", "ABC", "ABC"}, // numeric, both bases, either case
		{"&#x80;&#153;", "€™", "€™"},        // Windows-1252 remapping
		{"&#0;&#xD800;&#x110000;", "\uFFFD\uFFFD\uFFFD", "\uFFFD\uFFFD\uFFFD"},
		{"&#;&#x;&#xg;&", "&#;&#x;&#xg;&", "&#;&#x;&#xg;&"},
		{"&&amp", "&&", "&&"},
		{"a\r\nb\rc\n", "a\nb\nc\n", "a\nb\nc\n"},
		{"&#13;", "\r", "\r"}, // a reference to CR is not a line end
		{"a\x00b", "a\uFFFDb", "a\x00b"},
		{"&am\x00p;", "&am\uFFFDp;", "&am\x00p;"},
	}
	for _, c := range cases {
		attr := Value{raw: []byte(c.in), flags: valReader | valEntities | valAttr | valNUL}
		text := Value{raw: []byte(c.in), flags: valReader | valEntities}
		if got := attr.String(); got != c.attr {
			t.Errorf("attribute %q: got %q, want %q", c.in, got, c.attr)
		}
		if got := text.String(); got != c.text {
			t.Errorf("text %q: got %q, want %q", c.in, got, c.text)
		}
		if got := string(attr.AppendTo([]byte("x"))); got != "x"+c.attr {
			t.Errorf("attribute %q AppendTo: got %q", c.in, got)
		}
		if !attr.Equal(c.attr) || attr.Equal(c.attr+"y") {
			t.Errorf("attribute %q Equal is wrong", c.in)
		}
		// Through the reader, as an attribute value and as text.
		r := NewBytesReader([]byte(`<p a="`+c.in+`">`+c.in), Options{})
		if k, _ := r.Next(); k != StartTag {
			t.Fatalf("%q: %v", c.in, k)
		}
		if v, ok := r.Attr("a"); !ok || v.String() != c.attr {
			t.Errorf("reader attribute %q: got %q, %v", c.in, v.String(), ok)
		}
		if k, _ := r.Next(); k != Text {
			t.Fatalf("%q: %v", c.in, k)
		}
		if got := r.Text().String(); got != c.text {
			t.Errorf("reader text %q: got %q, want %q", c.in, got, c.text)
		}
	}
	raw := Value{raw: []byte("a&amp;b"), flags: 0}
	if raw.NeedsDecoding() || raw.String() != "a&amp;b" {
		t.Errorf("raw text decodes: %q", raw.String())
	}
}

func TestNames(t *testing.T) {
	r := NewBytesReader([]byte("<DIV Class=\"a\" ID=b class=c dAtA-x=\"1\" ><P\x00Q a\x00B=1></Div>"), Options{})
	if k, _ := r.Next(); k != StartTag || !r.NameIs("div") || string(r.Name()) != "div" {
		t.Fatalf("got %v %q", k, r.Name())
	}
	if n := r.attrCount(); n != 3 {
		t.Errorf("attrCount: %d, want 3 (class repeated)", n)
	}
	var names []string
	for {
		name, v, ok := r.NextAttr()
		if !ok {
			break
		}
		names = append(names, string(name)+"="+v.String())
	}
	if want := []string{"class=a", "id=b", "data-x=1"}; !reflect.DeepEqual(names, want) {
		t.Errorf("attributes: %q, want %q", names, want)
	}
	if v, ok := r.Attr("class"); !ok || !v.Equal("a") {
		t.Errorf("Attr(class): %q %v", v.Raw(), ok)
	}
	if _, ok := r.Attr("missing"); ok {
		t.Error("Attr(missing) found")
	}
	if k, _ := r.Next(); k != StartTag || string(r.Name()) != "p\uFFFDq" {
		t.Errorf("NUL in a tag name: %v %q", k, r.Name())
	}
	if name, v, ok := r.NextAttr(); !ok || string(name) != "a\uFFFDb" || !v.Equal("1") {
		t.Errorf("NUL in an attribute name: %q %q %v", name, v.Raw(), ok)
	}
	if k, _ := r.Next(); k != EndTag || !r.NameIs("div") {
		t.Errorf("end tag: %v %q", k, r.Name())
	}
	if k, err := r.Next(); k != EOF || err != nil {
		t.Errorf("after the last token: %v %v", k, err)
	}
	if k, err := r.Next(); k != EOF || err != nil {
		t.Errorf("EOF again: %v %v", k, err)
	}
}

func TestOffset(t *testing.T) {
	in := "ab<p>c<!--d--></p>"
	for _, f := range feeds {
		r := f.open(in, false, "")
		var offsets []int64
		for {
			k, err := r.Next()
			if err != nil || k == EOF {
				break
			}
			offsets = append(offsets, r.Offset())
		}
		if want := []int64{0, 2, 5, 6, 14}; !reflect.DeepEqual(offsets, want) {
			t.Errorf("%s: offsets %v, want %v", f.name, offsets, want)
		}
	}
}

func TestReaderErrors(t *testing.T) {
	r := NewReader(strings.NewReader("<p>"+strings.Repeat("x", 100)), Options{MaxBufferBytes: 50})
	var last error
	for {
		k, err := r.Next()
		if err != nil {
			last = err
			break
		}
		if k == EOF {
			break
		}
	}
	if !errors.Is(last, ErrTooLarge) {
		t.Errorf("long text with MaxBufferBytes 50: %v", last)
	}
	if k, err := r.Next(); k != None || !errors.Is(err, ErrTooLarge) {
		t.Errorf("after the error: %v %v", k, err)
	}

	boom := errors.New("boom")
	r = NewReader(io.MultiReader(strings.NewReader("<p>x"), &errReader{boom}), Options{})
	for {
		k, err := r.Next()
		if err != nil {
			last = err
			break
		}
		if k == EOF {
			t.Fatal("EOF before the source's error")
		}
	}
	if last != boom {
		t.Errorf("source error: %v", last)
	}

	r = NewReader(&stuckReader{}, Options{})
	if _, err := r.Next(); err != io.ErrNoProgress {
		t.Errorf("stuck source: %v", err)
	}
}

type errReader struct{ err error }

func (r *errReader) Read([]byte) (int, error) { return 0, r.err }

type stuckReader struct{}

func (*stuckReader) Read([]byte) (int, error) { return 0, nil }

func TestResetKeepsBuffer(t *testing.T) {
	r := NewReader(strings.NewReader("<p>one</p>"), Options{BufferSize: 16})
	if _, err := readerTokens(r); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		src := strings.NewReader(fmt.Sprintf("<div id=%d>%s</div>", i, strings.Repeat("y", 40)))
		r.Reset(src)
		toks, err := readerTokens(r)
		if err != nil || len(toks) != 3 || toks[0].Attr[0].Val != strconv.Itoa(i) {
			t.Errorf("reset %d: %v %+v", i, err, toks)
		}
		r.ResetBytes([]byte("<b>z</b>"))
		toks, err = readerTokens(r)
		if err != nil || len(toks) != 3 || toks[1].Data != "z" {
			t.Errorf("reset bytes %d: %v %+v", i, err, toks)
		}
	}
}

func TestNextIsNotRawText(t *testing.T) {
	r := NewBytesReader([]byte("<title><b>x</b></title>"), Options{})
	r.Next()
	r.NextIsNotRawText()
	toks, err := readerTokens(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(toks) != 4 || toks[0].Kind != "StartTag" || toks[0].Data != "b" {
		t.Errorf("after NextIsNotRawText: %+v", toks)
	}
}

func TestKindString(t *testing.T) {
	if s := StartTag.String(); s != "StartTag" {
		t.Error(s)
	}
	if s := Kind(42).String(); s != "Kind(42)" {
		t.Error(s)
	}
	var b bytes.Buffer
	fmt.Fprint(&b, EOF)
	if b.String() != "EOF" {
		t.Error(b.String())
	}
}

// The standard library's html.UnescapeString is not the oracle: it lacks
// the two-code-point references (&nGt;), leaves "&#0" as written and wraps
// a numeric reference that overflows, where the tokenizer gives U+FFFD.
func TestUnescapeString(t *testing.T) {
	inputs := []string{"", "plain", "a&amp;b", "&copy", "&copy=1", "&notin", "&notin;", "&#x80;&#0;&#xD800;", "a\r\nb\x00&", "&&amp;&amp", "&NotEqualTilde;x"}
	for _, tt := range tokenTests {
		inputs = append(inputs, tt.html)
	}
	inputs = append(inputs, html5libTokenizerInputs(t)...)
	for _, in := range inputs {
		if got, want := UnescapeString(in), html.UnescapeString(in); got != want {
			t.Errorf("%q: got %q, want %q", in, got, want)
		}
	}
	if s := "no ampersand"; UnescapeString(s) != s {
		t.Error("changed")
	}
}
