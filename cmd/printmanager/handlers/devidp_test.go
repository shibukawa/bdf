package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
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

// pw dev signs people in at pw's development IdP with the roster
// cmd/printmanager/devidp.toml (popcornweb.toml [dev.idp]): its admin, a
// member of warehouse-a and an outsider hit the rules of config.dev.toml.
// The same provider runs here (testutil.WithIdentityProvider), so the
// roster and the rules cannot drift apart. print/fakeidp stays for what
// the roster cannot say (a nested groups claim, an unverified email, a
// token tampered with).
func TestDevRoster(t *testing.T) {
	policy := devPolicy(t)
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
	if err := a.m.PutLocation(context.Background(), manager.Location{Name: "wh-a", Scope: manager.ScopeGroup, Owner: "warehouse-a"}); err != nil {
		t.Fatal(err)
	}
	if err := a.m.PutLocation(context.Background(), manager.Location{Name: "office"}); err != nil {
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
