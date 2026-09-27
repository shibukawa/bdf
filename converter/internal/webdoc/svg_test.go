package webdoc

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func wellFormed(s string) error {
	d := xml.NewDecoder(strings.NewReader(s))
	for {
		_, err := d.Token()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// An svg element becomes a well-formed SVG document, whether the page was
// read as XML or as HTML: the attributes of prefixes it does not declare
// (epub:type) and of namespaces it cannot name, and the elements of other
// vocabularies (an editor's metadata), are left out.
func TestSVGDocument(t *testing.T) {
	const svg = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
 xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape" xmlns:sodipodi="http://sodipodi.sourceforge.net/DTD/sodipodi-0.dtd"
 viewBox="0 0 10 10" inkscape:label="Layer" epub:type="figure" onclick="x()">
<sodipodi:namedview inkscape:zoom="1"/><rect width="10" height="10" fill="url(#g)" xml:space="preserve"/>
<use xlink:href="#s"/><image xlink:href="p.png"/><script>alert(1)</script>
<foreignObject><p xmlns="http://www.w3.org/1999/xhtml">html</p></foreignObject></svg>`
	xhtmlDoc := `<?xml version="1.0"?><html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops"><body>` + svg +
		`<svg xmlns="http://www.w3.org/2000/svg"><defs><linearGradient id="g"/><symbol id="s"/></defs></svg></body></html>`
	htmlDoc := `<!DOCTYPE html><body>` + svg + `<svg><defs><linearGradient id="g"/><symbol id="s"/></defs></svg>`
	for name, doc := range map[string]string{"xhtml": xhtmlDoc, "html": htmlDoc} {
		root, err := Parse([]byte(doc), "")
		if err != nil {
			t.Fatal(err)
		}
		ids := map[string]*html.Node{}
		WalkElements(root, func(e *html.Node) {
			if id := attrStr(e, "id"); id != "" {
				ids[id] = e
			}
		})
		out := string(SVGDocument(find(root, "svg"), "#123456", func(id string) *html.Node { return ids[id] },
			func(href string) string { return "data:image/png;base64,AA==" }))
		if err := wellFormed(out); err != nil {
			t.Errorf("%s: %v:\n%s", name, err, out)
		}
		for _, want := range []string{`xlink:href="data:image/png;base64,AA=="`, `xml:space="preserve"`, `<use xlink:href="#s">`,
			`<defs><linearGradient id="g"></linearGradient><symbol id="s"></symbol></defs>`, `color="#123456"`,
			`<p xmlns="http://www.w3.org/1999/xhtml">html</p>`} {
			if !strings.Contains(out, want) {
				t.Errorf("%s: lacks %s:\n%s", name, want, out)
			}
		}
		for _, not := range []string{"epub:type", "namedview", "script", "onclick"} {
			if strings.Contains(out, not) {
				t.Errorf("%s: keeps %s:\n%s", name, not, out)
			}
		}
		// a prefix the document declares is kept (the HTML parser keeps the declaration)
		if name == "html" && !strings.Contains(out, `inkscape:label="Layer"`) {
			t.Errorf("html: lacks the declared inkscape:label:\n%s", out)
		}
	}
}
