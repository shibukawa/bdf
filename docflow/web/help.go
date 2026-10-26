package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strings"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/web/help"
	"github.com/shibukawa/bdf/docflow/webmessages"
	"github.com/shibukawa/popcornweb/pw"
)

// The manual (package help) is a page of everyone signed in. /help goes
// to the overview in the language of the pages (the person's choice, else
// the browser's), and /help/<chapter> to that chapter, keeping the
// fragment (the browser does): the pages link to the manual without a
// language, /help/nodes#printer. /help/<lang>/<chapter> is a chapter;
// /help/<lang> the overview.
func (a *App) help(w http.ResponseWriter, r *http.Request, ac Account) {
	first, second := r.PathValue("first"), r.PathValue("second")
	if !help.IsLang(first) {
		if second != "" || first != "" && !slices.Contains(help.Chapters, first) {
			http.NotFound(w, r)
			return
		}
		lang := pw.RequestLocale(r).Tag()
		if !help.IsLang(lang) {
			lang = help.Languages[0]
		}
		http.Redirect(w, r, help.URL(lang, firstOr(first, help.Chapters[0])), http.StatusSeeOther)
		return
	}
	p, ok := help.Get(first, firstOr(second, help.Chapters[0]))
	if !ok {
		http.NotFound(w, r)
		return
	}
	if p.Name == "nodes" {
		p = withKindDocs(p)
	}
	pw.WriteHTML(w, r, Help(helpParams(r, ac, p)))
}

// withKindDocs adds to the chapter of the nodes the chapters of the kinds
// a program registered with a Doc (docflow.Kind.Doc), one section each,
// with the kind's name as its id (the pages link to /help/nodes#<kind>),
// in the language of the page.
func withKindDocs(p *help.Page) *help.Page {
	var md strings.Builder
	for _, k := range docflow.Kinds() {
		if k.Doc == nil {
			continue
		}
		fmt.Fprintf(&md, "## %s {#%s}\n\n%s\n\n", k.Titles.In(p.Lang), k.Name, strings.TrimSpace(k.Doc.In(p.Lang)))
	}
	if md.Len() == 0 {
		return p
	}
	q, err := help.Appended(p, []byte(md.String()))
	if err != nil {
		slog.Warn("help: the kinds' chapters", "err", err)
		return p
	}
	return q
}

// helpParams is a chapter as the page shows it: the chapters of its
// language, with the sections of the one shown. The switch of the
// language comes back to the chapter without a language, which is the
// chapter in the language chosen.
func helpParams(r *http.Request, ac Account, p *help.Page) HelpParams {
	var chapters []HelpChapter
	for _, c := range help.Pages(p.Lang) {
		ch := HelpChapter{Title: c.Title, Href: url.URL{Path: help.URL(p.Lang, c.Name)}, Current: c.Name == p.Name}
		if ch.Current {
			// the page's own sections: the chapter's, and those added
			// to it (withKindDocs)
			for _, s := range p.Sections {
				ch.Sections = append(ch.Sections, HelpSection{Title: s.Title, Href: url.URL{Path: help.URL(p.Lang, c.Name), Fragment: s.ID}})
			}
		}
		chapters = append(chapters, ch)
	}
	ac.Here = path.Join(help.Path, p.Name)
	title := webmessages.PagesHelpTitle(pw.RequestLocale(r), p.Title)
	return HelpParams{Account: ac, Title: title, Language: p.Lang, Chapters: chapters, Body: p.HTML}
}
