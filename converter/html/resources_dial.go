//go:build !tinygo

package html

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"time"
)

// newClient returns the client the images of documents are fetched with.
// Unless private, it connects to public addresses only (publicOnly): a
// document that a server converts could otherwise reach, through the
// server, what is private to it, such as the metadata service of its cloud
// or the other servers of its network. The check is skipped when the
// environment names a proxy, which is where the connections go then.
func newClient(private bool) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	d := &net.Dialer{Timeout: 15 * time.Second}
	if !private && !proxied() {
		d.Control = publicOnly
	}
	t.DialContext = d.DialContext
	return &http.Client{Timeout: 30 * time.Second, Transport: t}
}

// proxied reports whether the environment names a proxy for HTTP requests.
func proxied() bool {
	for _, scheme := range []string{"http", "https"} {
		u, err := http.ProxyFromEnvironment(&http.Request{URL: &url.URL{Scheme: scheme, Host: "example.com"}})
		if err == nil && u != nil {
			return true
		}
	}
	return false
}

// publicOnly refuses connections to addresses that are not public on the
// Internet: those of the host itself, of its network and of its link
// (where the metadata services of clouds answer). It is asked about the
// address connected to, after the name was resolved, so a name that
// resolves to such an address is refused too, and so is a redirection to
// one.
func publicOnly(_, address string, _ syscall.RawConn) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return errors.New("not an IP address")
	}
	if !ip.IsGlobalUnicast() || ip.IsPrivate() || sharedAddress(ip) {
		return fmt.Errorf("%s is not a public address", host)
	}
	return nil
}

// sharedAddressSpace is the carrier-grade NAT range (RFC 6598), which
// IsPrivate does not count.
var sharedAddressSpace = net.IPNet{IP: net.IPv4(100, 64, 0, 0), Mask: net.CIDRMask(10, 32)}

func sharedAddress(ip net.IP) bool {
	return sharedAddressSpace.Contains(ip)
}
