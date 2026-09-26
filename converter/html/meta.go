package html

import (
	"strings"

	"codeberg.org/readeck/go-readability/v2"
	"github.com/shibukawa/bdf"
	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// metadata is what a document says about itself in its head.
type metadata struct {
	dc          bdf.DublinCore // Dublin Core meta elements (DC.title, dcterms.created …)
	title       string
	creator     []string
	description string
	subject     []string
	publisher   string
	lang        string
	created     string
	modified    string
	identifier  string
}

// readMeta reads the language of the root element, the title, the meta
// elements (names, Open Graph properties and Dublin Core) and the
// canonical link of a document.
func readMeta(doc *xhtml.Node) *metadata {
	m := &metadata{}
	if h := findElement(doc, atom.Html); h != nil {
		m.lang = attr(h, "lang")
		if m.lang == "" {
			m.lang = attr(h, "xml:lang")
		}
	}
	head := findElement(doc, atom.Head)
	if head == nil {
		return m
	}
	var walk func(n *xhtml.Node)
	walk = func(n *xhtml.Node) {
		if n.Type == xhtml.ElementNode {
			switch n.DataAtom {
			case atom.Title:
				if m.title == "" {
					m.title = collapse(textOf(n))
				}
			case atom.Meta:
				m.meta(n)
			case atom.Link:
				if strings.EqualFold(attr(n, "rel"), "canonical") && m.identifier == "" {
					m.identifier = attr(n, "href")
				}
			}
		}
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			walk(k)
		}
	}
	walk(head)
	return m
}

// meta reads a meta element.
func (m *metadata) meta(n *xhtml.Node) {
	content := collapse(attr(n, "content"))
	if content == "" {
		return
	}
	name := strings.ToLower(attr(n, "name"))
	if name == "" {
		name = strings.ToLower(attr(n, "property"))
	}
	if eq := strings.ToLower(attr(n, "http-equiv")); eq == "content-language" && m.lang == "" {
		m.lang = strings.TrimSpace(strings.Split(content, ",")[0])
		return
	}
	for _, prefix := range []string{"dc.", "dcterms.", "dc:", "dcterms:"} {
		if term, ok := strings.CutPrefix(name, prefix); ok {
			if f := m.dc.Field(term); f != nil {
				*f = append(*f, content)
			}
			return
		}
	}
	switch name {
	case "author", "article:author":
		if !strings.HasPrefix(content, "http") {
			m.creator = append(m.creator, content)
		}
	case "description", "og:description":
		if m.description == "" {
			m.description = content
		}
	case "keywords":
		m.subject = append(m.subject, bdf.SplitKeywords(content)...)
	case "og:site_name":
		m.publisher = content
	case "og:title":
		if m.title == "" {
			m.title = content
		}
	case "article:published_time", "date", "citation_publication_date":
		if m.created == "" {
			m.created = content
		}
	case "article:modified_time", "last-modified":
		if m.modified == "" {
			m.modified = content
		}
	case "og:url":
		if m.identifier == "" {
			m.identifier = content
		}
	case "language", "og:locale":
		if m.lang == "" {
			m.lang = strings.ReplaceAll(content, "_", "-")
		}
	}
}

// fromArticle takes what Readability found about the article: its title
// (without the site's name), byline, site, language, excerpt and dates.
func (m *metadata) fromArticle(a readability.Article) {
	if t := collapse(a.Title()); t != "" {
		m.title = t
	}
	if b := collapse(a.Byline()); b != "" && len(m.creator) == 0 {
		m.creator = []string{b}
	}
	if s := collapse(a.SiteName()); s != "" && m.publisher == "" {
		m.publisher = s
	}
	if l := collapse(a.Language()); l != "" && m.lang == "" {
		m.lang = l
	}
	if e := collapse(a.Excerpt()); e != "" && m.description == "" {
		m.description = e
	}
	if t, err := a.PublishedTime(); err == nil && !t.IsZero() && m.created == "" {
		m.created = t.Format("2006-01-02T15:04:05Z07:00")
	}
	if t, err := a.ModifiedTime(); err == nil && !t.IsZero() && m.modified == "" {
		m.modified = t.Format("2006-01-02T15:04:05Z07:00")
	}
}

// dublinCore returns the document's metadata: the Dublin Core meta
// elements, completed by the rest (the title falls back to the first
// heading of the body).
func (m *metadata) dublinCore(body *xhtml.Node) bdf.DublinCore {
	dc := m.dc
	set := func(f *bdf.DCValues, v ...string) {
		if len(*f) > 0 {
			return
		}
		for _, s := range v {
			if s != "" {
				*f = append(*f, s)
			}
		}
	}
	title := m.title
	if title == "" {
		if h := findElement(body, atom.H1); h != nil {
			title = collapse(textOf(h))
		}
	}
	set(&dc.Title, title)
	set(&dc.Creator, m.creator...)
	set(&dc.Description, m.description)
	set(&dc.Subject, m.subject...)
	set(&dc.Publisher, m.publisher)
	set(&dc.Language, m.lang)
	set(&dc.Created, m.created)
	set(&dc.Modified, m.modified)
	set(&dc.Identifier, m.identifier)
	return dc
}

func textOf(n *xhtml.Node) string {
	var b strings.Builder
	var walk func(n *xhtml.Node)
	walk = func(n *xhtml.Node) {
		if n.Type == xhtml.TextNode {
			b.WriteString(n.Data)
		}
		for k := n.FirstChild; k != nil; k = k.NextSibling {
			walk(k)
		}
	}
	walk(n)
	return b.String()
}

// collapse collapses the white space of a string.
func collapse(s string) string { return strings.Join(strings.Fields(s), " ") }
