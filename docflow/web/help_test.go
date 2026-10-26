package web

import (
	"net/http"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/docflow/web/help"
)

// The manual is a page of everyone signed in: /help goes to the overview
// in the browser's language, /help/<chapter> to the chapter, and
// /help/<lang>/<chapter> is a chapter with the sidebar; the bar of every
// page links it.
func TestHelp(t *testing.T) {
	s := newSite(t)
	noRedirect := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	head := func(c *http.Client, path, acceptLanguage string) (int, string) {
		t.Helper()
		req, _ := http.NewRequest("GET", s.url+path, nil)
		req.Header.Set("Accept", "text/html,*/*;q=0.8")
		if acceptLanguage != "" {
			req.Header.Set("Accept-Language", acceptLanguage)
		}
		res, err := c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		return res.StatusCode, res.Header.Get("Location")
	}
	// signed out: to the sign-in, as any page
	if code, loc := head(noRedirect, "/help", ""); code != http.StatusSeeOther || !strings.HasPrefix(loc, "/auth/login") {
		t.Errorf("/help signed out: %d %s", code, loc)
	}

	alice := s.login("alice")
	mine := &http.Client{Jar: alice.c.Jar, CheckRedirect: noRedirect.CheckRedirect}
	// the language is the browser's, and the chapter is kept
	for _, c := range []struct{ path, accept, want string }{
		{"/help", "", "/help/en"},
		{"/help", "ja,en;q=0.8", "/help/ja"},
		{"/help", "en-US,en;q=0.9,ja;q=0.8", "/help/en"},
		{"/help/nodes", "ja-JP", "/help/ja/nodes"},
		{"/help/index", "fr", "/help/en"},
		{"/help/en/", "", "/help/en"},
	} {
		if code, loc := head(mine, c.path, c.accept); code != http.StatusSeeOther || loc != c.want {
			t.Errorf("%s (%q): %d %s; want %s", c.path, c.accept, code, loc, c.want)
		}
	}
	for _, path := range []string{"/help/xx", "/help/en/nope", "/help/nope", "/help/nodes/printer"} {
		if code, _ := head(mine, path, ""); code != http.StatusNotFound {
			t.Errorf("%s: %d", path, code)
		}
	}

	// a chapter: the manual's HTML, the sidebar with the chapter marked
	// and its sections under it, and the bar's switch of the language,
	// which comes back to the chapter in the language chosen
	code, body := alice.get("/help/en/nodes")
	if code != http.StatusOK {
		t.Fatalf("/help/en/nodes: %d", code)
	}
	for _, want := range []string{`<article id="help" lang="en"`, `<h1 id="node-kinds">`, `<h2 id="printer">`, `class="menu-active"`, `href="/help/en/nodes#printer"`,
		`name="next" value="/help/nodes"`, `>日本語<`, `href="/help/en"`, `href="/help/en/jobs"`, `<title>Node kinds · Help · BDF Docflow</title>`, `id="account-name"`, "<table>"} {
		if !strings.Contains(body, want) {
			t.Errorf("/help/en/nodes: no %s", want)
		}
	}
	if strings.Contains(body, ".md") || strings.Contains(body, "raw HTML omitted") {
		t.Error("/help/en/nodes: a link to a file, or HTML dropped")
	}
	// the chapter's language is its path's; the bar's is the pages'
	code, body = alice.get("/help/ja")
	if code != http.StatusOK || !strings.Contains(body, `<article id="help" lang="ja"`) || !strings.Contains(body, `<html lang="en"`) || !strings.Contains(body, `href="/help/ja/jobs"`) {
		t.Errorf("/help/ja: %d %.300s", code, body)
	}
	// every chapter of every language serves
	for _, lang := range help.Languages {
		for _, name := range help.Chapters {
			if code, body := alice.get(help.URL(lang, name)); code != http.StatusOK || !strings.Contains(body, "<h1 id=") {
				t.Errorf("%s: %d", help.URL(lang, name), code)
			}
		}
	}
	// the bar of every page links the manual
	if code, body := alice.get("/jobs"); code != http.StatusOK || !strings.Contains(body, `href="/help"`) {
		t.Errorf("/jobs: %d, no link to the manual", code)
	}
}
