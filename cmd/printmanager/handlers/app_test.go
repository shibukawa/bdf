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
	cfg := ManagerConfig{Records: "bbolt:" + filepath.Join(t.TempDir(), "pm.db"), AdminToken: "t", MaxDocumentMB: 1}

	serverless := cfg
	serverless.Serverless = true
	if a, err := Open(ctx, serverless); err == nil {
		a.Close()
		t.Fatal("serverless on a bbolt file")
	}
	bad := cfg
	bad.SubmitTokens = []string{"no-secret"}
	if a, err := Open(ctx, bad); err == nil {
		a.Close()
		t.Fatal("a submit token without its secret")
	}

	a, err := Open(ctx, cfg)
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
	for path, want := range map[string]string{"/": "/tablet", "/tablet": "/public/tablet.js", "/admin": "/public/admin.js"} {
		code, body := get(path, "")
		if code != http.StatusOK || !strings.Contains(body, want) {
			t.Errorf("%s: %d, no %s", path, code, want)
		}
	}
	// the manager's API is mounted
	if code, _ := get("/api/locations", ""); code != http.StatusUnauthorized {
		t.Errorf("api without a token: %d", code)
	}
	if code, body := get("/api/locations", "t"); code != http.StatusOK || strings.TrimSpace(body) != "[]" {
		t.Errorf("api: %d %s", code, body)
	}
}
