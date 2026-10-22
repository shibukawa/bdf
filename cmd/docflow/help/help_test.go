package help

import (
	"io/fs"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/docflow"
)

// Every chapter is in every language, with a title, and nothing a
// language has the other lacks: the files, and the ids of the headings
// (so that a link to a section, and the pages' own links to the manual,
// work in any language).
func TestChapters(t *testing.T) {
	for _, lang := range Languages {
		names, err := fs.Glob(files, lang+"/*.md")
		if err != nil {
			t.Fatal(err)
		}
		var found []string
		for _, n := range names {
			found = append(found, strings.TrimSuffix(strings.TrimPrefix(n, lang+"/"), ".md"))
		}
		want := slices.Clone(Chapters)
		slices.Sort(want)
		slices.Sort(found)
		if !slices.Equal(found, want) {
			t.Errorf("%s: the files are %v; the chapters are %v", lang, found, want)
		}
		if got := len(Pages(lang)); got != len(Chapters) {
			t.Errorf("%s: %d pages of %d chapters", lang, got, len(Chapters))
		}
	}
	for _, name := range Chapters {
		en, ok := Get("en", name)
		if !ok {
			t.Fatalf("en/%s is missing", name)
		}
		if en.Title == "" || !strings.Contains(en.HTML, "<h1") {
			t.Errorf("en/%s: no title", name)
		}
		if strings.Contains(en.HTML, "raw HTML omitted") {
			t.Errorf("en/%s: raw HTML in the source (the renderer drops it)", name)
		}
		if len(slices.Compact(slices.Sorted(slices.Values(en.IDs)))) != len(en.IDs) {
			t.Errorf("en/%s: a heading id is used twice: %v", name, en.IDs)
		}
		for _, id := range en.IDs {
			if strings.HasPrefix(id, "heading") || strings.ContainsFunc(id, func(r rune) bool { return r > 0x7f }) {
				t.Errorf("en/%s: heading %q needs an explicit {#id}", name, id)
			}
		}
		for _, lang := range Languages[1:] {
			p, ok := Get(lang, name)
			if !ok {
				t.Fatalf("%s/%s is missing", lang, name)
			}
			if strings.Contains(p.HTML, "raw HTML omitted") {
				t.Errorf("%s/%s: raw HTML in the source", lang, name)
			}
			if !slices.Equal(p.IDs, en.IDs) {
				t.Errorf("%s/%s: the headings' ids are %v; en's are %v", lang, name, p.IDs, en.IDs)
			}
		}
	}
}

// Every link between chapters reaches a chapter and, with a fragment, a
// heading of it; a link to a section of the same page reaches one; and
// the links are rewritten to the pages' URLs when rendered.
func TestLinks(t *testing.T) {
	for _, lang := range Languages {
		for _, name := range Chapters {
			links, err := Links(lang, name)
			if err != nil {
				t.Fatal(err)
			}
			for _, l := range links {
				file, fragment, _ := strings.Cut(l, "#")
				target := name
				switch {
				case strings.Contains(file, "://"):
					continue
				case file == "":
				case strings.HasSuffix(file, ".md") && !strings.Contains(file, "/"):
					target = strings.TrimSuffix(file, ".md")
				default:
					t.Errorf("%s/%s: link %q: a chapter is linked as chapter.md, elsewhere with a scheme", lang, name, l)
					continue
				}
				p, ok := Get(lang, target)
				if !ok {
					t.Errorf("%s/%s: link %q: no chapter %s", lang, name, l, target)
					continue
				}
				if fragment != "" && !slices.Contains(p.IDs, fragment) {
					t.Errorf("%s/%s: link %q: %s has no heading #%s (its ids: %v)", lang, name, l, target, fragment, p.IDs)
				}
			}
			p, _ := Get(lang, name)
			if strings.Contains(p.HTML, `href="`+name) || strings.Contains(p.HTML, `.md"`) || strings.Contains(p.HTML, `.md#`) {
				t.Errorf("%s/%s: a link to a chapter was not rewritten", lang, name)
			}
		}
	}
	if got := resolve("ja", "nodes.md#printer"); got != "/help/ja/nodes#printer" {
		t.Errorf("resolve: %s", got)
	}
	if got := resolve("en", "index.md"); got != "/help/en" {
		t.Errorf("resolve index: %s", got)
	}
	for _, keep := range []string{"#common", "https://example.com/a.md", "../README.md"} {
		if got := resolve("en", keep); got != keep {
			t.Errorf("resolve(%q) = %q", keep, got)
		}
	}
}

// The chapter of the node kinds has a section for every kind of node the
// manager has, under the kind's name: the editor links to them.
func TestKinds(t *testing.T) {
	kinds := []string{docflow.KindEndpoint, docflow.KindCups, docflow.KindQueue, docflow.KindClass, docflow.KindFilter, docflow.KindTemplate,
		docflow.KindConvert, docflow.KindInbox, docflow.KindDisplay, docflow.KindPrinter, docflow.KindArchive, docflow.KindSink, docflow.KindExport}
	for _, lang := range Languages {
		p, ok := Get(lang, "nodes")
		if !ok {
			t.Fatal("no nodes chapter")
		}
		for _, k := range kinds {
			if !slices.Contains(p.IDs, k) {
				t.Errorf("%s/nodes: no section #%s", lang, k)
			}
		}
	}
}

func TestLang(t *testing.T) {
	for accept, want := range map[string]string{
		"":                              "en",
		"ja,en-US;q=0.9,en;q=0.8":       "ja",
		"ja-JP":                         "ja",
		"fr-CH, fr;q=0.9, en;q=0.8, ja": "en",
		"de":                            "en",
		"*":                             "en",
		"JA":                            "ja",
	} {
		if got := Lang(accept); got != want {
			t.Errorf("Lang(%q) = %s, want %s", accept, got, want)
		}
	}
	if URL("ja", "index") != "/help/ja" || URL("en", "jobs") != "/help/en/jobs" {
		t.Error("URL")
	}
}

// A page renders to HTML a browser shows: the headings carry their ids,
// the links between chapters are the pages' URLs, and the tables are
// tables.
func TestRender(t *testing.T) {
	p, err := render("en", "x", []byte("# A title\n\nSee [jobs](jobs.md#print) and [below](#two).\n\n## Two {#two}\n\n| a | b |\n|---|---|\n| 1 | 2 |\n\n### Three {#three}\n"))
	if err != nil {
		t.Fatal(err)
	}
	if p.Title != "A title" || len(p.Sections) != 1 || p.Sections[0] != (Section{ID: "two", Title: "Two"}) || !slices.Equal(p.IDs, []string{"a-title", "two", "three"}) {
		t.Errorf("%+v", p)
	}
	for _, want := range []string{`<h2 id="two">Two</h2>`, `href="/help/en/jobs#print"`, `href="#two"`, "<table>", `<h3 id="three">`} {
		if !strings.Contains(p.HTML, want) {
			t.Errorf("no %s in\n%s", want, p.HTML)
		}
	}
	if _, err := render("en", "x", []byte("no heading\n")); err == nil {
		t.Error("a chapter without a title")
	}
	if _, err := render("en", "x", []byte("# t\n\n![a](b.png)\n")); err == nil {
		t.Error("an image")
	}
}
