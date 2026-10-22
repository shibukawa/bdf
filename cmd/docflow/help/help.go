// Package help is Docflow's manual: chapters written in
// Markdown beside this file, in English (en) and Japanese (ja), embedded
// in the binary and served by the web application at /help for everyone
// signed in. The same files read on GitHub: a link between chapters is
// written as jobs.md and rewritten to the page's URL when served.
//
// Every chapter has the same name in every language, and the same
// section ids (an explicit {#id} on each heading), so that a link to a
// section works in any language, and the pages' own links to the manual
// (/help/nodes#printer) need no language.
package help

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

//go:embed en ja
var files embed.FS

// Path is where the manual is served: /help/<lang>/<chapter>.
const Path = "/help"

// Languages the manual is written in, the first the default: the
// directories beside this file.
var Languages = []string{"en", "ja"}

// Chapters in the order of the manual: the files' names, without .md,
// the same in every language. The first is the overview, served at
// /help/<lang> as well.
var Chapters = []string{"index", "jobs", "topology", "nodes", "flow", "templates", "entries", "exports", "displays", "inboxes", "roles", "api", "operate"}

// Page is a chapter rendered: its title (the heading at the top), the
// body as HTML, and its sections (the second-level headings) for a table
// of contents.
type Page struct {
	Lang     string
	Name     string
	Title    string
	HTML     string
	Sections []Section
	// IDs are the ids of every heading of the page, for links to check
	// against.
	IDs []string
}

// Section is a heading of a page: its id (the fragment of a link to it)
// and its text.
type Section struct {
	ID    string
	Title string
}

// URL is the path of a chapter: /help/<lang>/<name>, the overview at
// /help/<lang>.
func URL(lang, name string) string {
	if name == Chapters[0] {
		return path.Join(Path, lang)
	}
	return path.Join(Path, lang, name)
}

// IsLang tells a language the manual is written in.
func IsLang(s string) bool { return slices.Contains(Languages, s) }

// Lang picks the language of a browser from its Accept-Language header:
// the first of the list the manual is written in (browsers list them in
// order of preference), else the default.
func Lang(acceptLanguage string) string {
	for _, part := range strings.Split(acceptLanguage, ",") {
		tag, _, _ := strings.Cut(strings.TrimSpace(part), ";")
		primary, _, _ := strings.Cut(strings.ToLower(tag), "-")
		if IsLang(primary) {
			return primary
		}
	}
	return Languages[0]
}

// Get returns a chapter in a language, rendered; false when there is no
// such chapter or language.
func Get(lang, name string) (*Page, bool) {
	p, ok := pages()[lang+"/"+name]
	return p, ok
}

// Pages returns the chapters of a language, in the manual's order.
func Pages(lang string) []*Page {
	var out []*Page
	for _, name := range Chapters {
		if p, ok := Get(lang, name); ok {
			out = append(out, p)
		}
	}
	return out
}

// pages renders every chapter of every language once. The sources are
// embedded and checked by the tests, so a failure is a bug of the build,
// not of a request.
var pages = sync.OnceValue(func() map[string]*Page {
	out := map[string]*Page{}
	for _, lang := range Languages {
		for _, name := range Chapters {
			src, err := fs.ReadFile(files, lang+"/"+name+".md")
			if err != nil {
				panic(fmt.Sprintf("help: %s/%s: %v", lang, name, err))
			}
			p, err := render(lang, name, src)
			if err != nil {
				panic(fmt.Sprintf("help: %s/%s: %v", lang, name, err))
			}
			out[lang+"/"+name] = p
		}
	}
	return out
})

// markdown is the chapters' dialect: GitHub's tables and the rest of
// GFM, an id on every heading (an explicit {#id}, else one made of the
// text), and no raw HTML (the renderer leaves a comment where it was
// written, which the tests refuse).
var markdown = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(parser.WithAutoHeadingID(), parser.WithAttribute()),
)

// render makes a page of a chapter's source: the Markdown parsed, the
// links between chapters rewritten to their URLs, the headings collected,
// and the tree rendered.
func render(lang, name string, src []byte) (*Page, error) {
	p := &Page{Lang: lang, Name: name}
	doc := markdown.Parser().Parse(text.NewReader(src))
	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch v := n.(type) {
		case *ast.Heading:
			id, _ := v.AttributeString("id")
			sid := string(id.([]byte))
			if v.Level == 1 && p.Title == "" {
				p.Title = textOf(v, src)
			}
			if v.Level == 2 {
				p.Sections = append(p.Sections, Section{ID: sid, Title: textOf(v, src)})
			}
			p.IDs = append(p.IDs, sid)
		case *ast.Link:
			v.Destination = []byte(resolve(lang, string(v.Destination)))
		case *ast.Image:
			return ast.WalkStop, fmt.Errorf("an image, %s: the manual is text", v.Destination)
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}
	if p.Title == "" {
		return nil, fmt.Errorf("no title (a heading at the top)")
	}
	var buf bytes.Buffer
	if err := markdown.Renderer().Render(&buf, src, doc); err != nil {
		return nil, err
	}
	p.HTML = buf.String()
	return p, nil
}

// resolve rewrites a link to another chapter (jobs.md, nodes.md#printer:
// what reads on GitHub) to the page's URL; a link to a section of the
// same page (#id) and a link elsewhere (https://…) stay.
func resolve(lang, dest string) string {
	file, fragment, _ := strings.Cut(dest, "#")
	name, ok := strings.CutSuffix(file, ".md")
	if !ok || strings.Contains(name, "/") || strings.Contains(name, ":") {
		return dest
	}
	u := URL(lang, name)
	if fragment != "" {
		u += "#" + fragment
	}
	return u
}

// textOf is the text of a node: what a heading says.
func textOf(n ast.Node, src []byte) string {
	var b strings.Builder
	ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch v := c.(type) {
			case *ast.Text:
				b.Write(v.Segment.Value(src))
			case *ast.String:
				b.Write(v.Value)
			}
		}
		return ast.WalkContinue, nil
	})
	return b.String()
}

// Links lists the destinations of the links of a chapter's source, for
// the tests to check: each as written.
func Links(lang, name string) ([]string, error) {
	src, err := fs.ReadFile(files, lang+"/"+name+".md")
	if err != nil {
		return nil, err
	}
	var out []string
	doc := markdown.Parser().Parse(text.NewReader(src))
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if l, ok := n.(*ast.Link); ok && entering {
			out = append(out, string(l.Destination))
		}
		return ast.WalkContinue, nil
	})
	return out, nil
}
