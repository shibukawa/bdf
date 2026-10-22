package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestApp(t *testing.T) {
	ctx := context.Background()
	cfg := DocflowConfig{Records: "bbolt:" + filepath.Join(t.TempDir(), "pm.db"), AdminToken: "t", MaxDocumentMB: 1}

	serverless := cfg
	serverless.Serverless = true
	if a, err := Open(ctx, serverless, OIDCConfig{}); err == nil {
		a.Close()
		t.Fatal("serverless on a bbolt file")
	}
	bad := cfg
	bad.SubmitTokens = []string{"no-secret"}
	if a, err := Open(ctx, bad, OIDCConfig{}); err == nil {
		a.Close()
		t.Fatal("a submit token without its secret")
	}

	a, err := Open(ctx, cfg, OIDCConfig{})
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()
	get := func(path, token string) (int, string) {
		req, _ := http.NewRequest("GET", srv.URL+path, nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		return res.StatusCode, string(b)
	}
	// without an IdP (dev), nobody signs in: the top page says so, and
	// the pages of people send to it (with one, the top page signs in at
	// the IdP and the others go there: login_test)
	for path, want := range map[string]string{"/": "AUTH_OIDC_ISSUER", "/tablet": "/public/tablet.js"} {
		code, body := get(path, "")
		if code != http.StatusOK || !strings.Contains(body, want) {
			t.Errorf("%s: %d, no %s", path, code, want)
		}
	}
	noRedirect := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	for _, path := range []string{"/jobs", "/topology", "/admin", "/audit"} {
		res, err := noRedirect.Get(srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusSeeOther || res.Header.Get("Location") != "/" {
			t.Errorf("%s without a session: %d %s", path, res.StatusCode, res.Header.Get("Location"))
		}
	}
	// the manager's API is mounted
	if code, _ := get("/api/nodes", ""); code != http.StatusUnauthorized {
		t.Errorf("api without a token: %d", code)
	}
	if code, body := get("/api/nodes", "t"); code != http.StatusOK || strings.TrimSpace(body) != "[]" {
		t.Errorf("api: %d %s", code, body)
	}
	// the endpoints and the IPP server are mounted too: the manager
	// answers (JSON), not the framework's page
	for _, path := range []string{"/in/nowhere", "/ipp/print"} {
		req, _ := http.NewRequest("POST", srv.URL+path, strings.NewReader("{}"))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode == http.StatusOK || !strings.Contains(string(b), "error") && !strings.Contains(string(b), "IPP") {
			t.Errorf("%s: %d %s", path, res.StatusCode, b)
		}
	}
}
