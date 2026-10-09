//go:build tinygo

package html

import (
	"net/http"
	"time"
)

// newClient returns the client the images of documents are fetched with.
// In the browser, what a page may reach is for the browser to decide, and
// TinyGo's net has no dialer to check the addresses in.
func newClient(bool) *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}
