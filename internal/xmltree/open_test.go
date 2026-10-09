package xmltree

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"
	"unicode/utf8"

	"github.com/shibukawa/tinygodriver/encoding/xmlro"
	"github.com/shibukawa/tinygodriver/encoding/xmlro/htmlentity"
)

// tokens writes the tokens of a document as the readers of the formats
// other than Office's take them: the names and attributes with their
// namespaces, and the character data. The error that ends the reading is
// its last line.
func tokens(r *xmlro.Reader) string {
	var b strings.Builder
	var text []byte
	flush := func() {
		if len(text) > 0 {
			fmt.Fprintf(&b, "text %q\n", text)
			text = text[:0]
		}
	}
	for {
		k, err := r.Next()
		if err != nil {
			flush()
			b.WriteString("error\n")
			return b.String()
		}
		switch k {
		case xmlro.EOF:
			flush()
			return b.String()
		case xmlro.StartElement:
			flush()
			n := ElementName(r)
			fmt.Fprintf(&b, "start %s|%s", n.Space, n.Local)
			for _, a := range Attrs(r) {
				fmt.Fprintf(&b, " %s|%s=%q", a.Name.Space, a.Name.Local, a.Value)
			}
			b.WriteString("\n")
		case xmlro.EndElement:
			flush()
			fmt.Fprintf(&b, "end %s\n", Local(r.Name()))
		case xmlro.Text, xmlro.CData:
			text = AppendText(text, r)
		}
	}
}

// stdTokens writes the same of an encoding/xml decoder.
func stdTokens(d *xml.Decoder) string {
	var b strings.Builder
	var text []byte
	flush := func() {
		if len(text) > 0 {
			fmt.Fprintf(&b, "text %q\n", text)
			text = text[:0]
		}
	}
	for {
		tok, err := d.Token()
		if err == io.EOF {
			flush()
			return b.String()
		}
		if err != nil {
			flush()
			b.WriteString("error\n")
			return b.String()
		}
		switch t := tok.(type) {
		case xml.StartElement:
			flush()
			fmt.Fprintf(&b, "start %s|%s", t.Name.Space, t.Name.Local)
			for _, a := range t.Attr {
				fmt.Fprintf(&b, " %s|%s=%q", a.Name.Space, a.Name.Local, a.Value)
			}
			b.WriteString("\n")
		case xml.EndElement:
			flush()
			fmt.Fprintf(&b, "end %s\n", t.Name.Local)
		case xml.CharData:
			text = append(text, t...)
		}
	}
}

// TestOpenMatchesEncodingXML reads the markup that the formats other than
// Office's are written with, with the readers they open and with the
// decoders of encoding/xml they were read with before.
func TestOpenMatchesEncodingXML(t *testing.T) {
	docs := map[string]string{
		"namespaces":  `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xml:lang="en"><a xlink:href="#b" href="#c"/><x:y xmlns:x="urn:x" x:z="1"/><un:bound un:a="1"/><g xmlns=""><h/></g></svg>`,
		"declaration": `<?xml version="1.0" encoding="UTF-8"?>` + "\n" + `<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">` + "\n<svg><!-- c --><?pi x?><g/></svg>\n",
		"text":        "<r>a&amp;b &lt;c&gt; &#65;&#x42; <![CDATA[<raw>&amp;\r\nline]]> line\r\nend\rtoo<k/>tail</r>",
		"html":        `<r>&nbsp;&eacute;&hellip;&unknown; AT&T &amp<br><img src="a.png"><hr/><BR>x</BR><p>one<p>two</r>`,
		"attributes":  `<r a='single' b="dou&quot;ble" c="&#10;nl" d="  spaced  " e="tab	tab" checked f=bare g="&eacute;"/>`,
		"mismatched":  `<a><b><c>text</b>more</a>`,
		"stray end":   `<a><b></c></b></a>`,
		"truncated":   `<a><b>text`,
		"unclosed":    `<a><b attr="1"`,
		"two roots":   `<a/>text<b/>`,
		"empty":       ``,
		"colons":      `<:a b:="1" :c="2"><d:/></:a>`,
		"prefixed":    `<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns:dc="http://purl.org/dc/elements/1.1/"><rdf:Description rdf:about="" dc:format="image/png"><dc:title><rdf:Alt><rdf:li xml:lang="x-default">T</rdf:li></rdf:Alt></dc:title></rdf:Description></rdf:RDF>`,
		"latin-1":     "<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><r a=\"caf\xe9\">na\xefve</r>",
		"bom":         "\xef\xbb\xbf<r>x</r>",
		"deep":        strings.Repeat("<a>", 200) + "x" + strings.Repeat("</a>", 200),
	}
	latin1 := func(label string, r io.Reader) (io.Reader, error) { return Latin1(label, r) }
	for name, doc := range docs {
		for _, html := range []bool{false, true} {
			d := xml.NewDecoder(bytes.NewReader(bytes.TrimPrefix([]byte(doc), []byte("\xef\xbb\xbf"))))
			d.Strict = false
			d.CharsetReader = latin1
			o := xmlro.Options{Lenient: true}
			if html {
				d.Entity, d.AutoClose = xml.HTMLEntity, xml.HTMLAutoClose
				o.Entities, o.AutoClose = htmlentity.Lookup, htmlentity.AutoClose
			}
			want, got := stdTokens(d), tokens(Open([]byte(doc), o))
			if got != want {
				t.Errorf("%s (html %v):\n%s--- encoding/xml:\n%s", name, html, got, want)
			}
		}
	}
}

// TestOpenEncodings: UTF-16 with a byte order mark is decoded whatever the
// declaration says, and what a lenient reader gives is valid UTF-8.
func TestOpenEncodings(t *testing.T) {
	utf16 := func(s string, bigEndian bool) []byte {
		out := []byte{0xff, 0xfe}
		if bigEndian {
			out = []byte{0xfe, 0xff}
		}
		for _, c := range s {
			if bigEndian {
				out = append(out, byte(c>>8), byte(c))
			} else {
				out = append(out, byte(c), byte(c>>8))
			}
		}
		return out
	}
	const doc = `<?xml version="1.0" encoding="UTF-16"?><r a="é">日本</r>`
	want := "start |r |a=\"é\"\ntext \"日本\"\nend r\n"
	for _, bigEndian := range []bool{false, true} {
		if got := tokens(Open(utf16(doc, bigEndian), xmlro.Options{Lenient: true})); got != want {
			t.Errorf("UTF-16 (big endian %v): %s", bigEndian, got)
		}
	}
	// 日本 in Shift_JIS, declared after white space
	sjis := func(label string, r io.Reader) (io.Reader, error) {
		b, err := io.ReadAll(r)
		return bytes.NewReader(bytes.ReplaceAll(b, []byte("\x93\xfa\x96\x7b"), []byte("日本"))), err
	}
	if got := tokens(Open([]byte("\n <?xml version=\"1.0\" encoding=\"Shift_JIS\"?><r>\x93\xfa\x96\x7b</r>"), xmlro.Options{Lenient: true, CharsetReader: sjis})); !strings.Contains(got, `text "日本"`) {
		t.Errorf("declared after white space: %s", got)
	}
	got := tokens(Open([]byte("<r a=\"\xff\">a\xc3(b</r>"), xmlro.Options{Lenient: true}))
	if !utf8.ValidString(got) || !strings.Contains(got, "a\ufffd(b") {
		t.Errorf("bytes that are not UTF-8: %q", got)
	}
	// a strict reader gives them as they are, for its caller to refuse
	r := Open([]byte("<r>a\xc3(b</r>"), xmlro.Options{})
	r.Next()
	r.Next()
	if utf8.Valid(r.Text()) {
		t.Errorf("a strict reader replaced bytes: %q", r.Text())
	}
}

// TestOpenDeclarations: a document is converted once, however many
// declarations of an encoding it has (each had the rest of the document
// converted again).
func TestOpenDeclarations(t *testing.T) {
	const n = 20000
	doc := []byte(`<?xml version="1.0" encoding="latin1"?><r a="caf` + "\xe9" + `">` + strings.Repeat(`<?xml version="1.0" encoding="latin1"?><a/>`, n) + `</r>`)
	calls := 0
	r := Open(doc, xmlro.Options{Lenient: true, CharsetReader: func(label string, src io.Reader) (io.Reader, error) {
		calls++
		return Latin1(label, src)
	}})
	elements, value := 0, ""
	for {
		k, err := r.Next()
		if err != nil {
			t.Fatal(err)
		}
		if k == xmlro.EOF {
			break
		}
		if k == xmlro.StartElement {
			if elements++; elements == 1 {
				v, _ := r.Attr("a")
				value = v.String()
			}
		}
	}
	if calls != 1 || elements != n+1 || value != "café" {
		t.Errorf("%d conversions, %d elements, attribute %q", calls, elements, value)
	}
	// the text of comments and CDATA sections stays as it is
	got := tokens(Open([]byte(`<?xml version="1.0" encoding="latin1"?><r><!-- <?xml version="1.0"?> --><![CDATA[<?xml version="1.0" encoding="x"?>]]><?xml version="1.0" encoding="latin1"?><a/></r>`), xmlro.Options{Lenient: true}))
	if want := "start |r\ntext \"<?xml version=\\\"1.0\\\" encoding=\\\"x\\\"?>\"\nstart |a\nend a\nend r\n"; got != want {
		t.Errorf("declarations in a comment and a CDATA section: %s", got)
	}
	// declarations hidden from that by markup that is not well formed end
	// the reading at the first of them
	calls = 0
	hidden := []byte(`<?xml version="1.0" encoding="latin1"?><r><a t="<!--"/>` + strings.Repeat(`<?xml version="1.0" encoding="latin1"?><b/>`, n) + `<!-- --></r>`)
	got = tokens(Open(hidden, xmlro.Options{Lenient: true, CharsetReader: func(label string, src io.Reader) (io.Reader, error) {
		calls++
		return Latin1(label, src)
	}}))
	if calls != 1 || !strings.HasSuffix(got, "end a\nerror\n") {
		t.Errorf("hidden declarations: %d conversions, %s", calls, got)
	}
}

// TestBoundEntities: the entities of a DOCTYPE are replaced while they
// stand for MaxEntityText at most, and stay references past it.
func TestBoundEntities(t *testing.T) {
	big := strings.Repeat("A", 50<<10)
	doc := func(refs int) []byte {
		return []byte(`<!DOCTYPE r [<!-- <!ENTITY c "--> <!ENTITY a "` + big + `"> <!ENTITY % p "x"><!ENTITY ns 'urn:x'>]>` +
			`<r xmlns="&ns;">` + strings.Repeat("&a;", refs) + `</r>`)
	}
	within, dropped := BoundEntities(doc(20))
	if dropped || !bytes.Equal(within, doc(20)) {
		t.Errorf("20 references are dropped")
	}
	if got := tokens(Open(doc(20), xmlro.Options{Lenient: true})); !strings.HasPrefix(got, "start urn:x|r") || len(got) < 20*len(big) {
		t.Errorf("20 references: %d bytes, %.40q", len(got), got)
	}
	past, dropped := BoundEntities(doc(21))
	if !dropped || bytes.Contains(past, entityDecl) || len(past) != len(doc(21)) {
		t.Errorf("21 references are kept")
	}
	if got := tokens(Open(doc(21), xmlro.Options{Lenient: true})); !strings.HasPrefix(got, `start &ns;|r`) || !strings.Contains(got, `text "&a;&a;`) || len(got) > 1000 {
		t.Errorf("21 references: %d bytes, %.60q", len(got), got)
	}
	// in UTF-16, and in an encoding that is converted
	u := []byte{0xff, 0xfe}
	for _, c := range doc(21) {
		u = append(u, c, 0)
	}
	if got := tokens(Open(u, xmlro.Options{Lenient: true})); len(got) > 1000 {
		t.Errorf("UTF-16: %d bytes", len(got))
	}
	l1 := append([]byte(`<?xml version="1.0" encoding="latin1"?>`), doc(21)...)
	if got := tokens(Open(l1, xmlro.Options{Lenient: true})); len(got) > 1000 {
		t.Errorf("Latin-1: %d bytes", len(got))
	}
	swapped := func(label string, r io.Reader) (io.Reader, error) {
		// an encoding whose bytes are not those of the declarations
		b, err := io.ReadAll(r)
		return bytes.NewReader(bytes.ReplaceAll(b, []byte("<!NOTATION"), []byte("<!ENTITY"))), err
	}
	hidden := append([]byte(`<?xml version="1.0" encoding="x-swapped"?>`), bytes.ReplaceAll(doc(21), []byte("<!ENTITY"), []byte("<!NOTATION"))...)
	if got := tokens(Open(hidden, xmlro.Options{Lenient: true, CharsetReader: swapped})); len(got) > 1000 {
		t.Errorf("converted: %d bytes", len(got))
	}
	for doc, want := range map[string]int{
		`<!ENTITY a "12345">`:    5,
		`<!ENTITY a'12345'>`:     5,
		`<!ENTITY  a.b:c-d ''>`:  0,
		`<!ENTITY % a "12345">`:  -1,
		`<!ENTITY a SYSTEM "x">`: -1,
		`<!ENTITY "12345">`:      -1,
		`<!ENTITYa "12345">`:     -1,
		`<!ENTITY a "12345`:      -1,
	} {
		_, size, ok := declaredEntity([]byte(doc)[len(entityDecl):])
		if !ok {
			size = -1
		}
		if size != want {
			t.Errorf("%s: %d, want %d", doc, size, want)
		}
	}
}

// TestEntitiesRefused: a part of an Office package that declares entities
// is refused, however it is read.
func TestEntitiesRefused(t *testing.T) {
	doc := []byte(`<?xml version="1.0"?><!DOCTYPE r [<!ENTITY a "` + strings.Repeat("A", 1000) + `">]><r>&a;&a;</r>`)
	if _, err := Parse(doc); !errors.Is(err, ErrEntities) {
		t.Errorf("Parse: %v", err)
	}
	for name, src := range map[string]io.Reader{
		"whole":     bytes.NewReader(doc),
		"one byte":  iotest.OneByteReader(bytes.NewReader(doc)),
		"half":      iotest.HalfReader(bytes.NewReader(doc)),
		"data, EOF": iotest.DataErrReader(bytes.NewReader(doc)),
	} {
		if _, err := ParsePickingReader(src, nil); !errors.Is(err, ErrEntities) {
			t.Errorf("%s: %v", name, err)
		}
	}
	// a reader that is used for one part after another
	r := NewReader(strings.NewReader(`<r><a/></r>`))
	budget := MaxElements
	if _, err := ParseFrom(r, nil, &budget); err != nil {
		t.Fatal(err)
	}
	Reset(r, iotest.OneByteReader(bytes.NewReader(doc)))
	if _, err := ParseFrom(r, nil, &budget); !errors.Is(err, ErrEntities) {
		t.Errorf("after Reset: %v", err)
	}
	// a DOCTYPE without them is read past, in pieces too
	plainDoc := `<!DOCTYPE r SYSTEM "r.dtd"><r>&lt;!ENTIT<a/>Y</r>`
	if n, err := ParsePickingReader(iotest.OneByteReader(strings.NewReader(plainDoc)), nil); err != nil || n.Text != "<!ENTITY" {
		t.Errorf("DOCTYPE without entities: %v, %+v", err, n)
	}
}
