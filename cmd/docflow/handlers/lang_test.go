package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// The pages are in the language of their reader: the one they chose with
// the switch (POST /lang, pw's cookie), else their browser's
// (Accept-Language), else English. The API is not translated, and does
// not vary on either.
func TestLanguage(t *testing.T) {
	s := newSite(t)
	noRedirect := func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	ask := func(c *http.Client, method, path, accept string, form url.Values) (*http.Response, string) {
		t.Helper()
		var body io.Reader
		if form != nil {
			body = strings.NewReader(form.Encode())
		}
		req, _ := http.NewRequest(method, s.url+path, body)
		req.Header.Set("Accept", "text/html,*/*;q=0.8")
		req.Header.Set("Origin", s.url)
		if form != nil {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		if accept != "" {
			req.Header.Set("Accept-Language", accept)
		}
		res, err := c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		out, _ := io.ReadAll(res.Body)
		return res, string(out)
	}
	has := func(what, body string, wants ...string) {
		t.Helper()
		for _, w := range wants {
			if !strings.Contains(body, w) {
				t.Errorf("%s: no %s", what, w)
			}
		}
	}

	// signed out: the top page and the device's page in the browser's
	// language, varying on what they negotiate from
	anon := s.browser().c
	res, body := ask(anon, "GET", "/", "ja,en;q=0.8", nil)
	has("/ (ja)", body, `<html lang="ja"`, `>サインイン<`, `name="lang" value="en" lang="en"`, `>English<`, `name="next" value="/"`)
	if res.Header.Get("Content-Language") != "ja" || !strings.Contains(strings.Join(res.Header.Values("Vary"), ","), "Cookie") {
		t.Errorf("/ (ja): Content-Language %q, Vary %q", res.Header.Get("Content-Language"), res.Header.Values("Vary"))
	}
	_, body = ask(anon, "GET", "/", "fr", nil)
	has("/ (fr)", body, `<html lang="en"`, `>Sign in<`, `>日本語<`)
	_, body = ask(anon, "GET", "/tablet?code=ABCDE12345", "ja-JP", nil)
	has("/tablet (ja)", body, `<html lang="ja"`, `ディスプレイボード`, `name="next" value="/tablet?code=ABCDE12345"`)
	// the switch is theirs too: no sign-in on the way
	device := &http.Client{Jar: anon.Jar, CheckRedirect: noRedirect}
	if res, _ := ask(device, "POST", "/lang", "", url.Values{"lang": {"en"}, "next": {"/tablet?code=ABCDE12345"}}); res.StatusCode != http.StatusSeeOther || res.Header.Get("Location") != "/tablet?code=ABCDE12345" {
		t.Errorf("POST /lang signed out: %d %s", res.StatusCode, res.Header.Get("Location"))
	}
	_, body = ask(anon, "GET", "/tablet", "ja-JP", nil)
	has("/tablet after the switch", body, `<html lang="en"`, `>Display board<`)

	// signed in: the browser's language until the person chooses
	alice := s.login("alice")
	mine := &http.Client{Jar: alice.c.Jar, CheckRedirect: noRedirect}
	_, body = ask(mine, "GET", "/jobs", "", nil)
	has("/jobs", body, `<html lang="en"`, `>Print a document<`, `name="next" value="/jobs"`, `>日本語<`)
	res, _ = ask(mine, "POST", "/lang", "", url.Values{"lang": {"ja"}, "next": {"/admin/jobs?dest=x"}})
	if res.StatusCode != http.StatusSeeOther || res.Header.Get("Location") != "/admin/jobs?dest=x" {
		t.Errorf("POST /lang: %d %s", res.StatusCode, res.Header.Get("Location"))
	}
	// the choice outranks the browser's
	_, body = ask(mine, "GET", "/jobs", "en-US,en;q=0.9", nil)
	has("/jobs after the switch", body, `<html lang="ja"`, `>文書を印刷する<`, `>English<`, `/public/jobs.js`)
	if res, _ := ask(mine, "GET", "/help/nodes", "en", nil); res.Header.Get("Location") != "/help/ja/nodes" {
		t.Errorf("/help/nodes after the switch: %s", res.Header.Get("Location"))
	}
	_, body = ask(mine, "GET", "/audit", "", nil)
	has("/audit (forbidden, ja)", body, `このページは開けません`, `Alice には役割 audit がありません`)

	// a language the pages are not in changes nothing; next stays on
	// this site
	ask(mine, "POST", "/lang", "", url.Values{"lang": {"fr"}, "next": {"/jobs"}})
	if _, body = ask(mine, "GET", "/jobs", "", nil); !strings.Contains(body, `<html lang="ja"`) {
		t.Error("POST /lang fr changed the language")
	}
	for _, next := range []string{"//evil.example/", "https://evil.example/", `/\evil.example`, "jobs", ""} {
		if res, _ := ask(mine, "POST", "/lang", "", url.Values{"lang": {"en"}, "next": {next}}); res.Header.Get("Location") != "/" {
			t.Errorf("POST /lang next=%q: to %s", next, res.Header.Get("Location"))
		}
	}

	// the API is the same in every language
	res, _ = ask(mine, "GET", "/api/me", "ja", nil)
	if res.StatusCode != http.StatusOK || res.Header.Get("Content-Language") != "" || strings.Contains(strings.Join(res.Header.Values("Vary"), ","), "Cookie") {
		t.Errorf("/api/me: %d, Content-Language %q, Vary %q", res.StatusCode, res.Header.Get("Content-Language"), res.Header.Values("Vary"))
	}
}
