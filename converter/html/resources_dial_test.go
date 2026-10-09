//go:build !tinygo

package html

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// The images of a document are fetched from public addresses only, unless
// the document may reach private ones: the test server listens on the
// loopback address, which the host's own services do.
func TestRemoteImagesFromPublicAddressesOnly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write(onePixel()) }))
	defer srv.Close()
	deadline := time.Now().Add(time.Minute)
	if b, err := httpGet(client, srv.URL+"/a.png", deadline); err == nil || !strings.Contains(err.Error(), "not a public address") {
		t.Errorf("fetched %d bytes from the loopback address: %v", len(b), err)
	}
	if b, err := httpGet(privateClient, srv.URL+"/a.png", deadline); err != nil || !bytes.Equal(b, onePixel()) {
		t.Errorf("allowed private addresses: bytes %d, %v", len(b), err)
	}
	for _, o := range []Options{{}, {AllowPrivate: true}} {
		r := newResources(&o)
		_, err := r.get(srv.URL + "/b.png")
		if refused := err != nil && strings.Contains(err.Error(), "not a public address"); refused == o.AllowPrivate {
			t.Errorf("AllowPrivate %v: %v", o.AllowPrivate, err)
		}
	}
	for _, addr := range []string{"127.0.0.1:80", "[::1]:443", "10.1.2.3:80", "172.16.0.1:80", "192.168.1.1:80", "169.254.169.254:80", "[fe80::1]:80", "[fd00::1]:80", "100.64.0.1:80", "0.0.0.0:80", "[::]:80"} {
		if err := publicOnly("tcp", addr, nil); err == nil {
			t.Errorf("%s: connected", addr)
		}
	}
	for _, addr := range []string{"93.184.216.34:443", "[2606:2800:220:1:248:1893:25c8:1946]:443", "8.8.8.8:53"} {
		if err := publicOnly("tcp", addr, nil); err != nil {
			t.Errorf("%s: %v", addr, err)
		}
	}
}
