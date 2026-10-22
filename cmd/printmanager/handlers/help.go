package handlers

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/shibukawa/bdf/cmd/printmanager/help"
	"github.com/shibukawa/popcornweb/pw"
)

// The manual (package help) is a page of everyone signed in. /help goes
// to the overview in the browser's language (Accept-Language), and
// /help/<chapter> to that chapter, keeping the fragment (the browser
// does): the pages link to the manual without a language,
// /help/nodes#printer. /help/<lang>/<chapter> is a chapter; /help/<lang>
// the overview.
func (a *App) help(w http.ResponseWriter, r *http.Request, ac Account) {
	first, second := r.PathValue("first"), r.PathValue("second")
	if !help.IsLang(first) {
		if second != "" || first != "" && !slices.Contains(help.Chapters, first) {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, help.URL(help.Lang(r.Header.Get("Accept-Language")), firstOr(first, help.Chapters[0])), http.StatusSeeOther)
		return
	}
	p, ok := help.Get(first, firstOr(second, help.Chapters[0]))
	if !ok {
		http.NotFound(w, r)
		return
	}
	pw.WriteHTML(w, r, Help(helpParams(ac, p)))
}

// helpLangs names the manual's languages, each in itself.
var helpLangs = map[string]string{"en": "English", "ja": "日本語"}

// helpParams is a chapter as the page shows it: the chapters of its
// language, with the sections of the one shown, and the other languages
// at the same chapter.
func helpParams(ac Account, p *help.Page) HelpParams {
	var chapters []HelpChapter
	for _, c := range help.Pages(p.Lang) {
		ch := HelpChapter{Title: c.Title, Href: url.URL{Path: help.URL(p.Lang, c.Name)}, Current: c.Name == p.Name}
		if ch.Current {
			for _, s := range c.Sections {
				ch.Sections = append(ch.Sections, HelpSection{Title: s.Title, Href: url.URL{Path: help.URL(p.Lang, c.Name), Fragment: s.ID}})
			}
		}
		chapters = append(chapters, ch)
	}
	var langs []HelpLang
	for _, l := range help.Languages {
		if l != p.Lang {
			langs = append(langs, HelpLang{Label: firstOr(helpLangs[l], l), Href: url.URL{Path: help.URL(l, p.Name)}})
		}
	}
	return HelpParams{Account: ac, Title: p.Title + " · Help · BDF Print Manager", Language: p.Lang, Chapters: chapters, Langs: langs, Body: p.HTML}
}
