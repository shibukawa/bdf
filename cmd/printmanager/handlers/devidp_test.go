package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/print/manager"
	"github.com/shibukawa/popcornweb/plugin/auth"
	"github.com/shibukawa/popcornweb/pw"
	"github.com/shibukawa/popcornweb/sessionconfig"
	"github.com/shibukawa/popcornweb/testutil"
	"github.com/shibukawa/tinybind-go/minitoml"
)

// devPolicy reads the [oidc] rules of config.dev.toml at the module's root,
// the ones pw dev runs with.
func devPolicy(t *testing.T) OIDCConfig {
	t.Helper()
	b, err := os.ReadFile("../../../config.dev.toml")
	if err != nil {
		t.Fatal(err)
	}
	doc, err := minitoml.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	var o OIDCConfig
	if v, ok := doc.Get("oidc.groups_claim"); ok {
		o.GroupsClaim, _ = v.AsString()
	}
	if v, ok := doc.Get("oidc.admin_groups"); ok {
		o.AdminGroups, _ = v.AsStringSlice()
	}
	if v, ok := doc.Get("oidc.admins"); ok {
		o.Admins, _ = v.AsStringSlice()
	}
	return o
}

// devServer serves the application as pw dev does, signing in at pw's
// development IdP with the roster, admin pre-selected (no browser picks).
func devServer(t *testing.T, policy OIDCConfig) (*testutil.Server, *App) {
	t.Helper()
	cfg := ManagerConfig{Records: "bbolt:" + filepath.Join(t.TempDir(), "pm.db"), MaxDocumentMB: 1}
	a, err := Open(context.Background(), cfg, policy)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { a.Close() })
	srv := testutil.TestRun(t, a.Handler(), func(c *testutil.Config) {
		c.Update(func(sv *pw.ServerConfig) { sv.Public.Enabled = false })
	}, testutil.WithIdentityProvider(
		testutil.WithIdPConfig("../devidp.toml"),
		testutil.WithLoginUser("admin"),
		testutil.WithIdPBinding(func(c *testutil.Config, idp testutil.IdPInfo) {
			ac := c.Get[auth.Config]()
			sc := c.Get[sessionconfig.SessionConfig]()
			ac.OIDC.Issuer, ac.OIDC.ClientID, ac.OIDC.ClientSecret = idp.Issuer, idp.ClientID, idp.ClientSecret
			if err := ConfigureAuth(&ac, &sc, "a secret of thirty-two bytes or more", false, policy); err != nil {
				t.Fatal(err)
			}
			c.Set(ac)
			c.Set(sc)
		}),
	))
	return srv, a
}

// pw dev signs people in at pw's development IdP with the roster
// cmd/printmanager/devidp.toml (popcornweb.toml [dev.idp]): its admin, a
// member of warehouse-a and an outsider hit the rules of config.dev.toml.
// The same provider runs here (testutil.WithIdentityProvider), so the
// roster and the rules cannot drift apart. print/fakeidp stays for what
// the roster cannot say (a nested groups claim, an unverified email, a
// token tampered with).
func TestDevRoster(t *testing.T) {
	policy := devPolicy(t)
	srv, a := devServer(t, policy)
	if _, err := a.m.PutNode(context.Background(), manager.Node{Name: "wh-a", Kind: manager.KindDisplay, Scope: manager.ScopeGroup, Owner: "warehouse-a"}); err != nil {
		t.Fatal(err)
	}
	if _, err := a.m.PutNode(context.Background(), manager.Node{Name: "office", Kind: manager.KindDisplay}); err != nil {
		t.Fatal(err)
	}
	for who, want := range map[string]struct {
		admin int
		dests string
	}{
		"admin":    {http.StatusOK, "office,wh-a"},
		"member":   {http.StatusForbidden, "office,wh-a"},
		"outsider": {http.StatusForbidden, "office"},
	} {
		srv.LoginAs(t, who)
		jar, _ := cookiejar.New(nil)
		c := &http.Client{Jar: jar}
		res, err := c.Get(srv.URL + "/auth/login?next=%2Fjobs")
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusOK || res.Request.URL.Path != "/jobs" {
			t.Fatalf("%s's sign-in: %d at %s", who, res.StatusCode, res.Request.URL)
		}
		res, err = c.Get(srv.URL + "/admin")
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != want.admin {
			t.Errorf("%s on /admin: %d; want %d", who, res.StatusCode, want.admin)
		}
		res, err = c.Get(srv.URL + "/api/destinations")
		if err != nil {
			t.Fatal(err)
		}
		var ds []manager.Destination
		json.NewDecoder(res.Body).Decode(&ds)
		res.Body.Close()
		var names []string
		for _, d := range ds {
			names = append(names, d.Name)
		}
		if got := strings.Join(names, ","); got != want.dests {
			t.Errorf("%s sends to %s; want %s", who, got, want.dests)
		}
	}
}

// A sign-in after a sign-out lets the person choose who they are, even
// though the IdP would sign the same person in again without a word
// (pw dev's provider keeps its own session; so does a real IdP's single
// sign-on): popcornweb marks the sign-in to carry prompt=select_account and
// login until one completes. The requests a browser makes on its own right
// after the sign-out must not get in the way: Safari asks for the touch
// icons and /favicon.ico (with the cookies), which the guard answers 401
// rather than sending to the sign-in, and a start of a sign-in nobody
// finishes (a second tab, a prefetch) does not use the mark up
// (popcornweb 0.5.14; before it, the first of these spent the mark and the
// person's own click was answered from the IdP's session).
func TestSignOutThenChoose(t *testing.T) {
	srv, _ := devServer(t, devPolicy(t))
	jar, _ := cookiejar.New(nil)
	c := &http.Client{Jar: jar}
	me := func() string {
		t.Helper()
		res, err := c.Get(srv.URL + "/api/me")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		var v struct{ Name string }
		json.NewDecoder(res.Body).Decode(&v)
		return v.Name
	}
	res, err := c.Get(srv.URL + "/auth/login?next=%2Fjobs")
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusOK || res.Request.URL.Path != "/jobs" || me() != "Ada Admin" {
		t.Fatalf("admin's sign-in: %d at %s as %q", res.StatusCode, res.Request.URL, me())
	}
	req, _ := http.NewRequest("POST", srv.URL+"/auth/logout", nil)
	req.Header.Set("Origin", srv.URL)
	if res, err = c.Do(req); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.Request.URL.Path != "/" {
		t.Fatalf("signing out landed on %s", res.Request.URL)
	}
	// the browser's own requests on the landing page, as Safari makes them
	// (no Accept of text/html): not pages, so not sent to the sign-in
	for _, p := range []string{"/apple-touch-icon-precomposed.png", "/apple-touch-icon.png", "/favicon.ico"} {
		res, err := c.Get(srv.URL + p)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusUnauthorized || res.Request.URL.Path != p {
			t.Errorf("%s: %d at %s; want 401 there", p, res.StatusCode, res.Request.URL)
		}
	}
	// a sign-in begun and left (a second tab, a prefetch) does not use the
	// mark up: the person's own click still asks the IdP to let them choose
	c.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	for i := range 2 {
		if res, err = c.Get(srv.URL + "/auth/login?next=%2Fjobs"); err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusSeeOther {
			t.Fatalf("sign-in %d after signing out: %d", i, res.StatusCode)
		}
	}
	authorize, err := url.Parse(res.Header.Get("Location"))
	if err != nil || res.StatusCode != http.StatusSeeOther || !strings.HasPrefix(authorize.String(), srv.IdPInfo().Issuer) {
		t.Fatalf("sign-in after signing out: %d -> %q", res.StatusCode, res.Header.Get("Location"))
	}
	if got := authorize.Query().Get("prompt"); !slices.Contains(strings.Fields(got), "select_account") {
		t.Fatalf("prompt=%q; want select_account", got)
	}
	// which the IdP does, instead of signing the admin in again
	if res, err = c.Get(authorize.String()); err != nil {
		t.Fatal(err)
	}
	page, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusOK || !strings.Contains(string(page), "Ada Admin") || !strings.Contains(string(page), "Mina Member") {
		t.Fatalf("the IdP after signing out: %d %.300s", res.StatusCode, page)
	}
	// and the person picks someone else
	hidden := func(name string) string {
		m := regexp.MustCompile(`name="` + name + `" value="([^"]+)"`).FindStringSubmatch(string(page))
		if m == nil {
			t.Fatalf("no %s in the IdP's page", name)
		}
		return m[1]
	}
	form := url.Values{"auth": {hidden("auth")}, "csrf": {hidden("csrf")}, "subject": {"member"}}
	c.CheckRedirect = nil
	if res, err = c.PostForm(srv.IdPInfo().Issuer+"/login", form); err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusOK || res.Request.URL.Path != "/jobs" || me() != "Mina Member (warehouse-a)" {
		t.Errorf("after choosing: %d at %s as %q", res.StatusCode, res.Request.URL, me())
	}
}
