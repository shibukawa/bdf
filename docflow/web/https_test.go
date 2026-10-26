package web

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/shibukawa/bdf/docflow"
)

// An enrollment's URL is at an origin the devices on the network reach:
// the setting's, else the request's, at this computer's address for
// localhost and on the HTTPS port when there is one.
func TestDeviceOrigin(t *testing.T) {
	for _, tc := range []struct {
		deviceURL, lan string
		httpsPort      int
		host, proto    string
		tls            bool
		want           string
	}{
		{host: "localhost:8080", lan: "192.168.1.7", want: "http://192.168.1.7:8080"},
		{host: "127.0.0.1:8080", lan: "192.168.1.7", httpsPort: 8443, want: "https://192.168.1.7:8443"},
		{host: "[::1]:8080", lan: "192.168.1.7", want: "http://192.168.1.7:8080"},
		{host: "localhost:8080", want: "http://localhost:8080"},
		{host: "[::1]", want: "http://[::1]"},
		// a name the admin reached it at is the devices' too
		{host: "docflow.example.com", lan: "192.168.1.7", proto: "https", want: "https://docflow.example.com"},
		{host: "printhost:8080", lan: "192.168.1.7", httpsPort: 8443, want: "https://printhost:8443"},
		{host: "192.168.1.7:8443", lan: "192.168.1.7", httpsPort: 8443, tls: true, want: "https://192.168.1.7:8443"},
		{deviceURL: "https://docflow.example.com/", host: "localhost:8080", lan: "192.168.1.7", httpsPort: 8443, want: "https://docflow.example.com"},
	} {
		a := &App{cfg: Config{DeviceURL: tc.deviceURL}, httpsPort: tc.httpsPort, lanAddr: func() string { return tc.lan }}
		r := httptest.NewRequest("POST", "/api/nodes/wh-a/enrollments", nil)
		r.Host = tc.host
		if tc.proto != "" {
			r.Header.Set("X-Forwarded-Proto", tc.proto)
		}
		if tc.tls {
			r.TLS = &tls.ConnectionState{}
		}
		if got := a.deviceOrigin(r); got != tc.want {
			t.Errorf("%+v: %s, want %s", tc, got, tc.want)
		}
	}
}

func TestCheckHTTPS(t *testing.T) {
	for _, cfg := range []Config{
		{DeviceURL: "docflow.example.com"},
		{DeviceURL: "https://docflow.example.com/tablet"},
		{DeviceURL: "ftp://docflow.example.com"},
		{HTTPSPort: 70000},
		{HTTPSPort: 8443, Serverless: true},
		{HTTPSPort: 8443, TLSCertFile: "cert.pem"},
	} {
		if err := checkHTTPS(cfg); err == nil {
			t.Errorf("%+v: no error", cfg)
		}
	}
	for _, cfg := range []Config{
		{},
		{DeviceURL: "https://docflow.example.com/"},
		{DeviceURL: "http://192.168.1.7:8080"},
		// development: a certificate of its own
		{HTTPSPort: 8443},
		{HTTPSPort: 8443, TLSCertFile: "cert.pem", TLSKeyFile: "key.pem"},
	} {
		if err := checkHTTPS(cfg); err != nil {
			t.Errorf("%+v: %v", cfg, err)
		}
	}
}

// The development certificate names this computer, is kept, and is made
// again when it would not do.
func TestDevCertificate(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 10, 23, 9, 0, 0, 0, time.UTC)
	c, err := devCertificate(dir, "", now)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"localhost", "127.0.0.1", "::1"} {
		if err := c.Leaf.VerifyHostname(name); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
	serial := func(c tls.Certificate) string { return c.Leaf.SerialNumber.String() }
	again, err := devCertificate(dir, "", now.Add(24*time.Hour))
	if err != nil || serial(again) != serial(c) {
		t.Fatalf("not kept: %v", err)
	}
	// an address it does not name, or an end within 30 days
	if moved, _ := devCertificate(dir, "192.0.2.10", now); serial(moved) == serial(c) {
		t.Error("kept for an address it does not name")
	}
	kept, _ := devCertificate(dir, "", now)
	if late, _ := devCertificate(dir, "", kept.Leaf.NotAfter.Add(-29*24*time.Hour)); serial(late) == serial(kept) {
		t.Error("kept near its end")
	}
	// with nowhere to keep it, it is made for the run
	if c, err := devCertificate("", "", now); err != nil || c.Leaf == nil {
		t.Fatal(err)
	}
}

// The HTTPS port serves the framework's chain, and enrollments made over
// http send the devices there.
func TestStartHTTPS(t *testing.T) {
	cache := t.TempDir()
	t.Setenv("HOME", cache)
	t.Setenv("XDG_CACHE_HOME", cache)
	t.Setenv("LocalAppData", cache)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	ctx := context.Background()
	a, err := Open(ctx, Config{Records: "bbolt:" + filepath.Join(t.TempDir(), "pm.db"), AdminToken: "t", MaxDocumentMB: 1, HTTPSPort: port}, OIDCConfig{}, Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	a.lanAddr = func() string { return "192.0.2.10" }
	stop, err := a.StartHTTPS()
	if err != nil {
		t.Fatal(err)
	}
	defer stop()
	c := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	get := func() (int, string) {
		res, err := c.Get("https://127.0.0.1:" + strconv.Itoa(a.httpsPort) + "/tablet")
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		return res.StatusCode, string(b)
	}
	if code, _ := get(); code != http.StatusServiceUnavailable {
		t.Errorf("before the chain: %d", code)
	}
	a.Outermost(a.Handler())
	if code, body := get(); code != http.StatusOK || !strings.Contains(body, "/_docflow/tablet.js") {
		t.Errorf("/tablet over TLS: %d", code)
	}

	srv := httptest.NewServer(a.Handler())
	defer srv.Close()
	if _, err := a.m.PutNode(ctx, docflow.Node{Name: "wh-a", Kind: docflow.KindDisplay}); err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest("POST", strings.Replace(srv.URL, "127.0.0.1", "localhost", 1)+"/api/nodes/wh-a/enrollments", nil)
	req.Header.Set("Authorization", "Bearer t")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	var en struct{ URL string }
	json.NewDecoder(res.Body).Decode(&en)
	if want := "https://192.0.2.10:" + strconv.Itoa(a.httpsPort) + "/tablet?code="; !strings.HasPrefix(en.URL, want) {
		t.Errorf("enrollment %d %q, want %s…", res.StatusCode, en.URL, want)
	}
}
