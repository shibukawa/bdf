package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print/fakeidp"
	"github.com/shibukawa/bdf/print/manager"
	"github.com/shibukawa/popcornweb/plugin/auth"
	"github.com/shibukawa/popcornweb/pw"
	"github.com/shibukawa/popcornweb/pwconfig"
	"github.com/shibukawa/popcornweb/sessionconfig"
	"github.com/shibukawa/popcornweb/testutil"
)

// The people of the fake IdP: groups in Keycloak's realm_access.roles.
var people = []fakeidp.User{
	{Subject: "alice", Name: "Alice", Email: "alice@example.com", Groups: []string{"warehouse-a"}},
	{Subject: "bob", Name: "Bob", Email: "bob@example.com", Groups: []string{"warehouse-b"}},
	{Subject: "root", Name: "Root", Email: "root@example.com"},
	{Subject: "lead", Name: "Lead", Email: "lead@example.com", Groups: []string{"warehouse-a", "pm-admins"}},
	{Subject: "carol", Name: "Carol", Email: "carol@elsewhere.example"},
	{Subject: "mallory", Name: "Mallory", Email: "root@example.com", EmailUnverified: true},
}

var testPolicy = OIDCConfig{GroupsClaim: "realm_access.roles", Admins: []string{"root@example.com"}, AdminGroups: []string{"pm-admins"},
	Allowed: AllowedConfig{Domains: []string{"example.com"}}}

// site is the application served as pw serves it (the middleware, the
// sign-in), signing in at a fake IdP.
type site struct {
	t   *testing.T
	idp *fakeidp.IdP
	app *App
	url string
}

func newSite(t *testing.T) *site {
	t.Helper()
	idp := fakeidp.New(people...)
	idp.GroupsClaim = testPolicy.GroupsClaim
	issuer := fakeidp.Start(t, idp)
	cfg := ManagerConfig{Records: "bbolt:" + filepath.Join(t.TempDir(), "pm.db"), AdminToken: "t", MaxDocumentMB: 1,
		SubmitTokens: []string{"wms:s3cret"}}
	a, err := Open(context.Background(), cfg, testPolicy)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { a.Close() })
	srv := testutil.TestRun(t, a.Handler(), func(c *testutil.Config) {
		ac := c.Get[auth.Config]()
		sc := c.Get[sessionconfig.SessionConfig]()
		ac.OIDC.Issuer, ac.OIDC.ClientID, ac.OIDC.ClientSecret = issuer, fakeidp.ClientID, fakeidp.ClientSecret
		if err := ConfigureAuth(&ac, &sc, "a secret of thirty-two bytes or more", false, testPolicy); err != nil {
			t.Fatal(err)
		}
		c.Set(ac)
		c.Set(sc)
		// the public tree is package main's
		c.Update(func(sv *pw.ServerConfig) { sv.Public.Enabled = false })
	})
	return &site{t: t, idp: idp, app: a, url: srv.URL}
}

// browser is a person's browser: a cookie jar.
type browser struct {
	s *site
	c *http.Client
}

func (s *site) browser() *browser {
	jar, _ := cookiejar.New(nil)
	return &browser{s: s, c: &http.Client{Jar: jar}}
}

// login signs a person in as the top page's button does, through the
// IdP, and lands on /jobs.
func (s *site) login(who string) *browser {
	s.t.Helper()
	b := s.browser()
	s.idp.LoginAs(who)
	code, body := b.get("/auth/login?next=%2Fjobs")
	if code != http.StatusOK || !strings.Contains(body, `id="account-name"`) {
		s.t.Fatalf("%s's sign-in: %d %.400s", who, code, body)
	}
	return b
}

// send asks as a browser does for a page (Accept names text/html): the
// guard sends only such requests to the sign-in, and answers the rest 401.
func (b *browser) send(method, path string, body []byte, edit func(*http.Request)) (int, string) {
	b.s.t.Helper()
	req, _ := http.NewRequest(method, b.s.url+path, bytes.NewReader(body))
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	if edit != nil {
		edit(req)
	}
	res, err := b.c.Do(req)
	if err != nil {
		b.s.t.Fatal(err)
	}
	defer res.Body.Close()
	out, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(out)
}

func (b *browser) get(path string) (int, string) { return b.send("GET", path, nil, nil) }

// post sends as the pages do: from the site's origin.
func (b *browser) post(path string, body []byte) (int, string) {
	return b.send("POST", path, body, func(r *http.Request) { r.Header.Set("Origin", b.s.url) })
}

func doc(t *testing.T) []byte {
	d := bdf.NewDocument()
	v := d.NewView("pages", bdf.ViewFixed, "")
	h, _ := d.AddObject(bdf.NewObject())
	v.AddPage(595, 842, bdf.Layer{Role: bdf.RoleBody, Obj: h})
	var b bytes.Buffer
	if err := d.WriteSingle(&b); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestLogin(t *testing.T) {
	s := newSite(t)
	api := func(method, path string, body any, token string, status int) {
		t.Helper()
		var r io.Reader
		switch b := body.(type) {
		case []byte:
			r = bytes.NewReader(b)
		case nil:
		default:
			j, _ := json.Marshal(b)
			r = bytes.NewReader(j)
		}
		req, _ := http.NewRequest(method, s.url+path, r)
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		out, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode != status {
			t.Fatalf("%s %s: %d %s; want %d", method, path, res.StatusCode, out, status)
		}
	}
	// the settings, and jobs of the WMS (a submit token) to each warehouse
	api("PUT", "/api/locations/wh-a", manager.Location{Title: "Warehouse A", Scope: manager.ScopeGroup, Owner: "warehouse-a"}, "t", http.StatusOK)
	api("PUT", "/api/locations/wh-b", manager.Location{Title: "Warehouse B", Scope: manager.ScopeGroup, Owner: "warehouse-b"}, "t", http.StatusOK)
	api("PUT", "/api/locations/office", manager.Location{Title: "Office"}, "t", http.StatusOK)
	api("POST", "/api/destinations/wh-a/jobs?name=picking-a&user=sato", doc(t), "wms:s3cret", http.StatusCreated)
	api("POST", "/api/destinations/wh-b/jobs?name=picking-b", doc(t), "wms:s3cret", http.StatusCreated)

	// nobody signed in: the top page is a sign-in, every other page but
	// /tablet goes to the IdP when a browser asks for the page, and is 401
	// to anything else (an icon a browser fetches on its own, a script)
	if code, body := s.browser().get("/"); code != http.StatusOK || !strings.Contains(body, `href="/auth/login?next=%2Fjobs"`) {
		t.Errorf("/ signed out: %d %.300s", code, body)
	}
	noRedirect := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	for _, path := range []string{"/jobs", "/admin", "/audit"} {
		for _, accept := range []string{"text/html,*/*;q=0.8", "image/avif,image/webp,*/*"} {
			req, _ := http.NewRequest("GET", s.url+path, nil)
			req.Header.Set("Accept", accept)
			res, err := noRedirect.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			res.Body.Close()
			if strings.HasPrefix(accept, "text/html") {
				if res.StatusCode != http.StatusSeeOther || !strings.HasPrefix(res.Header.Get("Location"), "/auth/login") {
					t.Errorf("%s without a sign-in: %d %s", path, res.StatusCode, res.Header.Get("Location"))
				}
			} else if res.StatusCode != http.StatusUnauthorized {
				t.Errorf("%s asked for as an image without a sign-in: %d %s", path, res.StatusCode, res.Header.Get("Location"))
			}
		}
	}
	if code, body := s.browser().get("/tablet"); code != http.StatusOK || !strings.Contains(body, "/public/tablet.js") {
		t.Errorf("/tablet: %d", code)
	}

	// a member of warehouse-a sends a job from the page, as herself
	alice := s.login("alice")
	if code, body := alice.post("/api/destinations/office/jobs?name=alice-memo&user=someone", doc(t)); code != http.StatusCreated {
		t.Fatalf("alice's upload: %d %s", code, body)
	}
	if code, _ := alice.post("/api/destinations/wh-b/jobs?name=alice-wh-b", doc(t)); code != http.StatusForbidden {
		t.Errorf("alice prints to warehouse-b: %d", code)
	}
	bob := s.login("bob")
	if code, body := bob.post("/api/destinations/office/jobs?name=bob-memo", doc(t)); code != http.StatusCreated {
		t.Fatalf("bob's upload: %d %s", code, body)
	}

	jobs := func(b *browser) []manager.Job {
		t.Helper()
		code, body := b.get("/api/my/jobs")
		if code != http.StatusOK {
			t.Fatalf("/api/my/jobs: %d %s", code, body)
		}
		var js []manager.Job
		json.Unmarshal([]byte(body), &js)
		return js
	}
	for _, j := range jobs(alice) {
		if j.Name == "alice-memo" {
			if by := j.SubmittedBy; by == nil || by.Kind != manager.PrincipalUser || by.Via != manager.ViaOIDC || by.Name != "Alice" || j.ClaimedUser != "someone" {
				t.Errorf("alice's job by %+v, claimed %q", j.SubmittedBy, j.ClaimedUser)
			}
		}
	}

	// each sees their own and their groups' jobs, an admin all
	sees := func(b *browser, want []string) {
		t.Helper()
		var names []string
		for _, j := range jobs(b) {
			names = append(names, j.Name)
		}
		if strings.Join(names, ",") != strings.Join(want, ",") {
			t.Errorf("sees %q; want %q", names, want)
		}
	}
	sees(alice, []string{"alice-memo", "picking-a"})
	sees(bob, []string{"bob-memo", "picking-b"})
	if code, body := alice.get("/jobs"); code != http.StatusOK || strings.Contains(body, `href="/admin"`) || !strings.Contains(body, "/public/jobs.js") {
		t.Errorf("alice's /jobs: %d", code)
	}
	root := s.login("root")
	sees(root, []string{"bob-memo", "alice-memo", "picking-b", "picking-a"})
	if _, body := root.get("/jobs"); !strings.Contains(body, `href="/admin"`) {
		t.Error("root's /jobs has no Admin")
	}

	// /admin and /audit: the admins (by email, by a group of a nested
	// claim), 403 for others
	lead := s.login("lead")
	for who, b := range map[string]*browser{"alice": alice, "bob": bob, "root": root, "lead": lead} {
		want := http.StatusForbidden
		if who == "root" || who == "lead" {
			want = http.StatusOK
		}
		for _, path := range []string{"/admin", "/audit"} {
			code, body := b.get(path)
			if code != want {
				t.Errorf("%s on %s: %d; want %d", who, path, code, want)
			}
			if want == http.StatusOK && (strings.Contains(body, `name="token"`) || !strings.Contains(body, "/public"+path+".js")) {
				t.Errorf("%s's %s: %.300s", who, path, body)
			}
		}
	}
	// the API: an admin's session works as the admin token did
	if code, _ := alice.get("/api/topology"); code != http.StatusForbidden {
		t.Errorf("alice on the API: %d", code)
	}
	if code, body := root.get("/api/audit"); code != http.StatusOK || !strings.Contains(body, `"who":"Alice","action":"login"`) {
		t.Errorf("root on the API: %d %s", code, body)
	}
	if code, _ := root.send("DELETE", "/api/locations/office", nil, func(r *http.Request) { r.Header.Set("Origin", "https://evil.example") }); code != http.StatusForbidden {
		t.Errorf("an admin's session from another origin: %d", code)
	}

	// who the IdP signs in but the policy keeps out: another domain, or
	// an admin's email the IdP did not verify
	for _, who := range []string{"carol", "mallory"} {
		b := s.browser()
		s.idp.LoginAs(who)
		if code, _ := b.get("/jobs"); code != http.StatusForbidden {
			t.Errorf("%s got in: %d", who, code)
		}
	}

	// a session whose person's record is gone is still the person
	// (without groups): the IdP's sign-in would answer it as signed in
	if err := s.app.m.DeletePerson(context.Background(), personID(s.idp.Issuer, "root")); err != nil {
		t.Fatal(err)
	}
	if code, _ := root.get("/admin"); code != http.StatusOK {
		t.Errorf("root without a record: %d", code)
	}

	if n := s.idp.Logins(); n != 6 {
		t.Errorf("%d sign-ins at the IdP; want 6", n)
	}

	// signing out ends the session
	alice.c.CheckRedirect = noRedirect.CheckRedirect
	if code, _ := alice.post("/auth/logout", nil); code != http.StatusSeeOther {
		t.Errorf("logout: %d", code)
	}
	if code, _ := alice.get("/jobs"); code != http.StatusSeeOther {
		t.Errorf("/jobs after signing out: %d", code)
	}
}

func TestConfigureAuth(t *testing.T) {
	defer pwconfig.SwapEnv("production", true)()
	ok := auth.Config{OIDC: auth.OIDCConfig{Issuer: "https://idp.example", ClientID: "pm", ClientSecret: "x"}}
	for name, c := range map[string]struct {
		a auth.Config
		o OIDCConfig
	}{
		"no IdP":         {auth.Config{}, testPolicy},
		"no secret":      {auth.Config{OIDC: auth.OIDCConfig{Issuer: "https://idp.example", ClientID: "pm"}}, testPolicy},
		"no admin rule":  {ok, OIDCConfig{GroupsClaim: "groups"}},
		"an admin typed": {ok, OIDCConfig{Admins: []string{"root"}}},
	} {
		var s sessionconfig.SessionConfig
		if err := ConfigureAuth(&c.a, &s, "secret", true, c.o); err == nil {
			t.Errorf("%s: started", name)
		}
	}
	a, s := ok, sessionconfig.SessionConfig{Backend: sessionconfig.SessionBackendRDB, Cookie: sessionconfig.SessionCookieConfig{Secure: true}}
	if err := ConfigureAuth(&a, &s, "secret", true, OIDCConfig{Admins: []string{"https://idp.example#123"}}); err != nil {
		t.Fatal(err)
	}
	if !a.Enabled || a.Mode != auth.ModeOIDCOnly || a.Backend != authBackend || a.OIDC.AllowLoopbackHTTP ||
		!s.Enabled || s.Backend != sessionconfig.SessionBackendCookie || s.Keyring.Secret == "" || !s.Cookie.Secure {
		t.Errorf("auth %+v, session %+v", a, s)
	}
	b, s2 := ok, sessionconfig.SessionConfig{}
	ConfigureAuth(&b, &s2, "secret", true, OIDCConfig{AdminGroups: []string{"x"}})
	if s2.Keyring.Secret != s.Keyring.Secret {
		t.Error("the session's keys differ between instances of one secret")
	}
}

func TestGroups(t *testing.T) {
	for _, c := range []struct {
		claim, prefix string
		claims        map[string]any
		want          string
	}{
		{"groups", "", map[string]any{"groups": []any{"a", "b", 3, "a"}}, "a,b"},
		{"groups", "", map[string]any{"groups": "solo"}, "solo"},
		{"realm_access.roles", "", map[string]any{"realm_access": map[string]any{"roles": []any{"x", "pm-y"}}}, "x,pm-y"},
		{"realm_access.roles", "pm-", map[string]any{"realm_access": map[string]any{"roles": []any{"x", "pm-y"}}}, "pm-y"},
		{"realm_access.roles", "", map[string]any{"realm_access": "flat"}, ""},
		{"roles", "", map[string]any{}, ""},
	} {
		claims := func(name string) (json.RawMessage, bool) {
			v, ok := c.claims[name]
			if !ok {
				return nil, false
			}
			b, _ := json.Marshal(v)
			return b, true
		}
		p := policy{OIDCConfig{GroupsClaim: c.claim, GroupsPrefix: c.prefix}}
		if got := strings.Join(p.groups(claims), ","); got != c.want {
			t.Errorf("%s %v: %q; want %q", c.claim, c.claims, got, c.want)
		}
	}
}
