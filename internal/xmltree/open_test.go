package xmltree

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"
	"time"
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

func TestOpenTokens(t *testing.T) {
	for _, tc := range []struct {
		name, doc, want string
		html            bool
	}{
		{"namespaces", `<r xmlns="urn:r" xmlns:p="urn:p"><p:x p:y="1" plain="2"/><none xmlns=""/></r>`, "start urn:r|r |xmlns=\"urn:r\" xmlns|p=\"urn:p\"\nstart urn:p|x urn:p|y=\"1\" |plain=\"2\"\nend x\nstart |none |xmlns=\"\"\nend none\nend r\n", false},
		{"declaration", `<?xml version="1.0"?><!DOCTYPE svg><svg><!-- ignored --><?pi ignored?><g/></svg>`, "start |svg\nstart |g\nend g\nend svg\n", false},
		{"text", "<r>a&amp;b &lt;c&gt; &#65;&#x42; <![CDATA[<raw>&amp;\r\nline]]> line\r\nend\rtoo<k/>tail</r>", "start |r\ntext \"a&b <c> AB <raw>&amp;\\nline line\\nend\\ntoo\"\nstart |k\nend k\ntext \"tail\"\nend r\n", false},
		{"html entities and void elements", `<r>a&nbsp;&unknown;<br>tail</r>`, "start |r\ntext \"a\\u00a0&unknown;\"\nstart |br\nend br\ntext \"tail\"\nend r\n", true},
		{"mismatched", `<a><b><c>text</b>more</a>`, "start |a\nstart |b\nstart |c\ntext \"text\"\nend c\nend b\ntext \"more\"\nend a\n", false},
		{"truncated", `<a><b>text`, "start |a\nstart |b\ntext \"text\"\nerror\n", false},
		{"unclosed attribute", `<a><b attr="1"`, "start |a\nerror\n", false},
		{"two roots", `<a/>text<b/>`, "start |a\nend a\ntext \"text\"\nstart |b\nend b\n", false},
		{"empty", "", "", false},
		{"latin-1", "<?xml version=\"1.0\" encoding=\"ISO-8859-1\"?><r a=\"caf\xe9\">na\xefve</r>", "start |r |a=\"café\"\ntext \"naïve\"\nend r\n", false},
		{"bom", "\xef\xbb\xbf<r>x</r>", "start |r\ntext \"x\"\nend r\n", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := xmlro.Options{Lenient: true}
			if tc.html {
				opts.Entities, opts.AutoClose = htmlentity.Lookup, htmlentity.AutoClose
			}
			if got := tokens(Open([]byte(tc.doc), opts)); got != tc.want {
				t.Errorf("tokens:\n%s\nwant:\n%s", got, tc.want)
			}
		})
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

// TestOpenDeclarations: a document is converted once, for the declaration
// it starts with: one that follows the root element's start is a
// processing instruction and nothing more (each had the rest of the
// document converted again).
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
	// wherever they are, and whatever is around them
	calls = 0
	hidden := []byte(`<?xml version="1.0" encoding="latin1"?><r><a t="<!--"/>` + strings.Repeat(`<?xml version="1.0" encoding="latin1"?><b/>`, n) + `<!-- --></r>`)
	got = tokens(Open(hidden, xmlro.Options{Lenient: true, CharsetReader: func(label string, src io.Reader) (io.Reader, error) {
		calls++
		return Latin1(label, src)
	}}))
	if calls != 1 || strings.Count(got, "start |b\n") != n || !strings.HasSuffix(got, "end b\nend r\n") {
		t.Errorf("declarations after markup that is not well formed: %d conversions, %d elements, …%.40q", calls, strings.Count(got, "start |b\n"), got[max(0, len(got)-40):])
	}
	// a declaration after the root element has begun converts nothing,
	// neither the first of a document
	calls = 0
	late := []byte("<r a=\"caf\xc3\xa9\"><?xml version=\"1.0\" encoding=\"latin1\"?><b>\xc3\xa9</b></r>")
	got = tokens(Open(late, xmlro.Options{Lenient: true, CharsetReader: func(label string, src io.Reader) (io.Reader, error) {
		calls++
		return Latin1(label, src)
	}}))
	if want := "start |r |a=\"café\"\nstart |b\ntext \"é\"\nend b\nend r\n"; calls != 0 || got != want {
		t.Errorf("a declaration inside the root element: %d conversions, %s", calls, got)
	}
	// nor does the second of two before it
	got = tokens(Open([]byte("<?xml version=\"1.0\"?><?xml version=\"1.0\" encoding=\"latin1\"?><r>\xc3\xa9</r>"), xmlro.Options{Lenient: true}))
	if want := "start |r\ntext \"é\"\nend r\n"; got != want {
		t.Errorf("a second declaration: %s", got)
	}
}

// TestOpenDoctypes: a document has one DOCTYPE, before its root element:
// the reading ends at another, and at one inside the root element, whose
// entities a document could otherwise declare where they are used.
func TestOpenDoctypes(t *testing.T) {
	for doc, want := range map[string]string{
		`<!DOCTYPE r [<!ENTITY a "b">]><r>&a;</r>`:                                  "start |r\ntext \"b\"\nend r\n",
		`<!DOCTYPE r [<!ENTITY a "b">]><!DOCTYPE r [<!ENTITY c "d">]><r>&a;&c;</r>`: "error\n",
		`<r>&a;<!DOCTYPE r [<!ENTITY a "b">]>&a;</r>`:                               "start |r\ntext \"&a;\"\nerror\n",
	} {
		if got := tokens(Open([]byte(doc), xmlro.Options{Lenient: true})); got != want {
			t.Errorf("%s: %s", doc, got)
		}
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

// TestBoundEntitiesTime: the references of a document are counted in time
// in proportion to it. With a name declared that is as long as the
// document, the ";" of every "&" was looked for that far: the square of a
// document without one.
func TestBoundEntitiesTime(t *testing.T) {
	const n = 300000
	doc := []byte(`<!DOCTYPE r [<!ENTITY ` + strings.Repeat("n", n) + ` "v"><!ENTITY a "b">]><r>` + strings.Repeat("&a", n) + `</r>`)
	start := time.Now()
	if _, dropped := BoundEntities(doc); dropped || time.Since(start) > 2*time.Second {
		t.Errorf("%d references without an end: dropped %v, %v", n, dropped, time.Since(start))
	}
	// a name longer than the reader takes for one is no reference: what is
	// declared for it is not counted
	long := strings.Repeat("n", maxEntityName+1)
	doc = []byte(`<!DOCTYPE r [<!ENTITY ` + long + ` "` + strings.Repeat("A", MaxEntityText) + `"><!ENTITY a "b">]><r>&` + long + `;&` + long + `;&a;</r>`)
	if _, dropped := BoundEntities(doc); dropped {
		t.Errorf("references by a name of %d bytes are counted", len(long))
	}
	if got := tokens(Open(doc, xmlro.Options{Lenient: true})); got != "start |r\ntext \"&"+long+";&"+long+";b\"\nend r\n" {
		t.Errorf("references by a name of %d bytes: %.80q", len(long), got)
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
