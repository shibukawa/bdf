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

// login opens / as a person and follows the sign-in through the IdP.
func (s *site) login(who string) *browser {
	s.t.Helper()
	b := s.browser()
	s.idp.LoginAs(who)
	code, body := b.get("/")
	if code != http.StatusOK || !strings.Contains(body, `id="account-name"`) {
		s.t.Fatalf("%s's sign-in: %d %.400s", who, code, body)
	}
	return b
}

func (b *browser) send(method, path string, body []byte, edit func(*http.Request)) (int, string) {
	b.s.t.Helper()
	req, _ := http.NewRequest(method, b.s.url+path, bytes.NewReader(body))
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
	api("PUT", "/api/locations/wh-a", manager.Location{Title: "第 1 倉庫", Kind: manager.KindGroup, Owner: "warehouse-a"}, "t", http.StatusOK)
	api("PUT", "/api/locations/wh-b", manager.Location{Title: "第 2 倉庫", Kind: manager.KindGroup, Owner: "warehouse-b"}, "t", http.StatusOK)
	api("PUT", "/api/locations/office", manager.Location{Title: "事務所"}, "t", http.StatusOK)
	api("POST", "/api/locations/wh-a/slips?name=picking-a&user=sato", doc(t), "wms:s3cret", http.StatusCreated)
	api("POST", "/api/locations/wh-b/slips?name=picking-b", doc(t), "wms:s3cret", http.StatusCreated)

	// nobody signed in: every page but /tablet goes to the IdP
	noRedirect := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	for _, path := range []string{"/", "/admin"} {
		res, err := noRedirect.Get(s.url + path)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusSeeOther || !strings.HasPrefix(res.Header.Get("Location"), "/auth/login") {
			t.Errorf("%s without a sign-in: %d %s", path, res.StatusCode, res.Header.Get("Location"))
		}
	}
	if code, body := s.browser().get("/tablet"); code != http.StatusOK || !strings.Contains(body, "/public/tablet.js") {
		t.Errorf("/tablet: %d", code)
	}

	// a member of warehouse-a sends a job from the page, as herself
	alice := s.login("alice")
	if code, body := alice.post("/api/locations/office/slips?name=alice-memo&user=someone", doc(t)); code != http.StatusCreated {
		t.Fatalf("alice's upload: %d %s", code, body)
	}
	if code, _ := alice.post("/api/locations/wh-b/slips?name=alice-wh-b", doc(t)); code != http.StatusForbidden {
		t.Errorf("alice prints to warehouse-b: %d", code)
	}
	bob := s.login("bob")
	if code, body := bob.post("/api/locations/office/slips?name=bob-memo", doc(t)); code != http.StatusCreated {
		t.Fatalf("bob's upload: %d %s", code, body)
	}

	slips, err := s.app.m.Slips(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for _, sl := range slips {
		if sl.Name == "alice-memo" {
			if by := sl.SubmittedBy; by == nil || by.Kind != manager.PrincipalUser || by.Via != manager.ViaOIDC || by.Name != "Alice" || sl.ClaimedUser != "someone" {
				t.Errorf("alice's slip by %+v, claimed %q", sl.SubmittedBy, sl.ClaimedUser)
			}
		}
	}

	// each sees their own and their groups' jobs, an admin all
	sees := func(b *browser, want, not []string) {
		t.Helper()
		code, body := b.get("/")
		if code != http.StatusOK {
			t.Fatalf("/: %d", code)
		}
		for _, n := range want {
			if !strings.Contains(body, n) {
				t.Errorf("does not see %s", n)
			}
		}
		for _, n := range not {
			if strings.Contains(body, n) {
				t.Errorf("sees %s", n)
			}
		}
	}
	sees(alice, []string{"alice-memo", "picking-a", "warehouse-a", "第 1 倉庫"}, []string{"bob-memo", "picking-b", "第 2 倉庫", `href="/admin"`})
	sees(bob, []string{"bob-memo", "picking-b"}, []string{"alice-memo", "picking-a"})
	root := s.login("root")
	sees(root, []string{"alice-memo", "bob-memo", "picking-a", "picking-b", `href="/admin"`, "管理者"}, nil)

	// /admin: the admins (by email, by a group of a nested claim), 403 for others
	for who, b := range map[string]*browser{"alice": alice, "bob": bob, "root": root, "lead": s.login("lead")} {
		want := http.StatusForbidden
		if who == "root" || who == "lead" {
			want = http.StatusOK
		}
		code, body := b.get("/admin")
		if code != want {
			t.Errorf("%s on /admin: %d; want %d", who, code, want)
		}
		if want == http.StatusOK && (strings.Contains(body, `id="token"`) || !strings.Contains(body, "/public/admin.js")) {
			t.Errorf("%s's /admin: %.300s", who, body)
		}
	}
	// the API: an admin's session works as the admin token did
	if code, _ := alice.get("/api/locations"); code != http.StatusForbidden {
		t.Errorf("alice on the API: %d", code)
	}
	if code, body := root.get("/api/slips"); code != http.StatusOK || !strings.Contains(body, "alice-memo") {
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
		if code, _ := b.get("/"); code != http.StatusForbidden {
			t.Errorf("%s got in: %d", who, code)
		}
	}

	if n := s.idp.Logins(); n != 6 {
		t.Errorf("%d sign-ins at the IdP; want 6", n)
	}

	// signing out ends the session
	alice.c.CheckRedirect = noRedirect.CheckRedirect
	if code, _ := alice.post("/auth/logout", nil); code != http.StatusSeeOther {
		t.Errorf("logout: %d", code)
	}
	if code, _ := alice.get("/"); code != http.StatusSeeOther {
		t.Errorf("/ after signing out: %d", code)
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
