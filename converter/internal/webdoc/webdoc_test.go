package webdoc

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func find(n *html.Node, tag string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		return n
	}
	for k := n.FirstChild; k != nil; k = k.NextSibling {
		if f := find(k, tag); f != nil {
			return f
		}
	}
	return nil
}

func text(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			walk(k)
		}
	}
	walk(n)
	return b.String()
}

func attr(n *html.Node, ns, key string) (string, bool) {
	for _, a := range n.Attr {
		if a.Namespace == ns && a.Key == key {
			return a.Val, true
		}
	}
	return "", false
}

const xhtmlDoc = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" xml:lang="ja">
<head><title/><script src="x.js"/></head>
<body>
<p><a id="p5"/>本文&nbsp;です。</p>
<section epub:type="chapter"><p>次</p></section>
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 10 10">
<image width="10" height="10" xlink:href="cover.jpg"/><foreignObject/></svg>
</body>
</html>`

// Elements written empty stay empty, as XML says; the HTML parser would
// have put the rest of the document in the script.
func TestXHTML(t *testing.T) {
	doc, err := Parse([]byte(xhtmlDoc), "")
	if err != nil {
		t.Fatal(err)
	}
	body := find(doc, "body")
	if body == nil || body.DataAtom != atom.Body {
		t.Fatalf("no body")
	}
	p := find(body, "p")
	if got := text(p); got != "本文 です。" {
		t.Errorf("paragraph text %q", got)
	}
	if a := find(p, "a"); a == nil || a.FirstChild != nil {
		t.Errorf("empty a element: %+v", a)
	}
	root := find(doc, "html")
	if v, _ := attr(root, "", "lang"); v != "ja" {
		t.Errorf("lang %q", v)
	}
	if v, _ := attr(root, "", "xml:lang"); v != "ja" {
		t.Errorf("xml:lang %q", v)
	}
	if v, _ := attr(find(doc, "section"), "", "epub:type"); v != "chapter" {
		t.Errorf("epub:type %q", v)
	}
	svg := find(doc, "svg")
	if svg.Namespace != "svg" || svg.DataAtom != atom.Svg {
		t.Errorf("svg element %q %v", svg.Namespace, svg.DataAtom)
	}
	img := find(svg, "image")
	if v, _ := attr(img, "xlink", "href"); v != "cover.jpg" {
		t.Errorf("xlink:href %q (%+v)", v, img.Attr)
	}
	if fo := find(svg, "foreignObject"); fo == nil || fo.Namespace != "svg" {
		t.Errorf("foreignObject keeps its case: %+v", fo)
	}
}

// Documents that are not well-formed are read as HTML.
func TestFallback(t *testing.T) {
	doc, err := Parse([]byte(`<?xml version="1.0"?><html><body><p>a<br>b</p></body></html>`), "")
	if err != nil {
		t.Fatal(err)
	}
	if got := text(find(doc, "p")); got != "ab" {
		t.Errorf("text %q", got)
	}
	if _, err := ParseXHTML([]byte(`<p>a<br>b</p>`)); err == nil {
		t.Errorf("unclosed br parsed as XML")
	}
	if _, err := ParseXHTML([]byte(`<a/><b/>`)); err == nil {
		t.Errorf("two root elements parsed")
	}
	// what the HTML parser reads better: entities XHTML does not have, an
	// ampersand alone, attributes without quotes, bytes that are not text
	for _, body := range []string{
		`<p>&check;</p>`, `<p>AT&T</p>`, `<p>&nbsp</p>`, `<p a="&bogus;"/>`, `<p a="x&y"/>`, `<p a="<"/>`, `<p a=b/>`, `<p hidden/>`,
		"<p>\xff</p>", "<p a=\"\xc3(\"/>", "<p>\x01</p>", "<p><![CDATA[\x00]]></p>", `<p>]]></p>`, `<p><!-- a -- b --></p>`, `<p>`, `</p>`,
	} {
		if _, err := ParseXHTML([]byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body>` + body + `</body></html>`)); err == nil {
			t.Errorf("%q parsed as XML", body)
		}
	}
	// and what is XML
	for _, body := range []string{
		`<p a="&amp;&lt;&#65;&#x1F600;&#0000065;&nbsp;">&apos;&quot;&gt;&eacute;</p>`, `<p><![CDATA[AT&T <b> &bogus;]]></p>`, "<p>\t\r\n</p>", `<p><!-- AT&T - --></p>`,
	} {
		if _, err := ParseXHTML([]byte(`<html xmlns="http://www.w3.org/1999/xhtml"><body>` + body + `</body></html>`)); err != nil {
			t.Errorf("%q: %v", body, err)
		}
	}
	// the entities a DOCTYPE declares are XHTML's own
	doc, err = ParseXHTML([]byte(`<!DOCTYPE html [<!ENTITY me "A &amp; B">]><html xmlns="http://www.w3.org/1999/xhtml"><body><p title="&me;">&me;</p></body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	if p := find(doc, "p"); text(p) != "A & B" || p.Attr[0].Val != "A & B" {
		t.Errorf("declared entity: %q, %q", text(p), p.Attr[0].Val)
	}
}

func TestEncoding(t *testing.T) {
	// 日本 in Shift_JIS
	sjis := []byte("<?xml version=\"1.0\" encoding=\"Shift_JIS\"?><html><body><p>\x93\xfa\x96\x7b</p></body></html>")
	doc, err := ParseXHTML(sjis)
	if err != nil {
		t.Fatal(err)
	}
	if got := text(find(doc, "p")); got != "日本" {
		t.Errorf("text %q", got)
	}
}

func TestIsXML(t *testing.T) {
	for _, c := range []struct {
		data, ct string
		want     bool
	}{
		{"<?xml version='1.0'?><html/>", "", true},
		{"\ufeff\n <?xml version='1.0'?><html/>", "", true},
		{"<!DOCTYPE html><html>", "", false},
		{"<html>", "application/xhtml+xml; charset=utf-8", true},
		{"<html>", "text/html", false},
	} {
		if got := IsXML([]byte(c.data), c.ct); got != c.want {
			t.Errorf("IsXML(%q, %q) = %v", c.data, c.ct, got)
		}
	}
}

func TestDataURL(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"data:text/plain;base64,aGVsbG8=", "hello"},
		{"data:text/plain;base64,aGVs\nbG8", "hello"},
		{"data:,a%20b", "a b"},
	} {
		b, err := DataURL(c.in)
		if err != nil || string(b) != c.want {
			t.Errorf("DataURL(%q) = %q, %v", c.in, b, err)
		}
	}
	if _, err := DataURL("data:nocomma"); err == nil {
		t.Errorf("malformed URL decoded")
	}
	if !IsDataURL("DATA:,x") || IsDataURL("dat") {
		t.Errorf("IsDataURL")
	}
}
