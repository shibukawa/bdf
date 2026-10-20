package server

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"html"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// AirPrintProfile is a configuration profile (.mobileconfig) that adds the
// IPP printers to iPhones and iPads that do not see them with mDNS (on
// another subnet): an AirPrint payload of each queue (the address, its
// path, the port, TLS), and the server's certificate as a root to trust
// when the server signed it itself. An MDM sends it, or a person opens it
// in Safari (the page of the printers links it). host is the name or the
// address the devices reach the server at; port the IPP server's.
func AirPrintProfile(cfg *Config, host string, port int) ([]byte, error) {
	if cfg.IPP == nil {
		return nil, fmt.Errorf("the configuration has no ipp")
	}
	x := &ippServer{s: &Server{cfg: cfg, log: slog.New(slog.DiscardHandler)}, cfg: cfg.IPP, queues: map[string]*ippQueue{},
		secureOnly: cfg.IPP.Auth == authBasic || cfg.IPP.RequireTLS}
	if err := x.init(); err != nil {
		return nil, err
	}
	cert, err := x.certificate()
	if err != nil {
		return nil, err
	}
	return x.profile(host, port, cert), nil
}

// profile writes the configuration profile of the queues at host:port.
func (x *ippServer) profile(host string, port int, cert tls.Certificate) []byte {
	id := func(parts ...string) string {
		sum := sha256.Sum256([]byte("bdf-printd-profile\x00" + x.s.cfg.Spool + "\x00" + host + "\x00" + strings.Join(parts, "\x00")))
		b := sum[:16]
		b[6] = b[6]&0x0f | 0x50
		b[8] = b[8]&0x3f | 0x80
		h := strings.ToUpper(hex.EncodeToString(b))
		return h[:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
	}
	reverse := "local.bdf.printd." + strings.NewReplacer(".", "-", ":", "-").Replace(host)
	esc := html.EscapeString
	var b bytes.Buffer
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>PayloadContent</key>
	<array>
		<dict>
			<key>PayloadType</key>
			<string>com.apple.airprint</string>
			<key>PayloadVersion</key>
			<integer>1</integer>
`)
	fmt.Fprintf(&b, "\t\t\t<key>PayloadIdentifier</key>\n\t\t\t<string>%s.airprint</string>\n", esc(reverse))
	fmt.Fprintf(&b, "\t\t\t<key>PayloadUUID</key>\n\t\t\t<string>%s</string>\n", id("airprint"))
	fmt.Fprintf(&b, "\t\t\t<key>PayloadDisplayName</key>\n\t\t\t<string>AirPrint (%s)</string>\n", esc(host))
	b.WriteString("\t\t\t<key>AirPrint</key>\n\t\t\t<array>\n")
	for _, q := range x.order {
		fmt.Fprintf(&b, "\t\t\t\t<dict>\n\t\t\t\t\t<key>IPAddress</key>\n\t\t\t\t\t<string>%s</string>\n", esc(host))
		fmt.Fprintf(&b, "\t\t\t\t\t<key>ResourcePath</key>\n\t\t\t\t\t<string>%s</string>\n", esc(strings.TrimPrefix(q.path, "/")))
		fmt.Fprintf(&b, "\t\t\t\t\t<key>Port</key>\n\t\t\t\t\t<integer>%d</integer>\n", port)
		fmt.Fprintf(&b, "\t\t\t\t\t<key>ForceTLS</key>\n\t\t\t\t\t<%t/>\n\t\t\t\t</dict>\n", x.secureOnly || len(cert.Certificate) > 0)
	}
	b.WriteString("\t\t\t</array>\n\t\t</dict>\n")
	if x.cfg.TLS == nil && len(cert.Certificate) > 0 {
		// the certificate the server signed itself, to trust
		fmt.Fprintf(&b, `		<dict>
			<key>PayloadType</key>
			<string>com.apple.security.root</string>
			<key>PayloadVersion</key>
			<integer>1</integer>
			<key>PayloadIdentifier</key>
			<string>%s.certificate</string>
			<key>PayloadUUID</key>
			<string>%s</string>
			<key>PayloadDisplayName</key>
			<string>bdf printd (%s)</string>
			<key>PayloadCertificateFileName</key>
			<string>bdf-printd.cer</string>
			<key>PayloadContent</key>
			<data>%s</data>
		</dict>
`, esc(reverse), id("certificate", fingerprint(cert)), esc(host), base64.StdEncoding.EncodeToString(cert.Certificate[0]))
	}
	fmt.Fprintf(&b, `	</array>
	<key>PayloadType</key>
	<string>Configuration</string>
	<key>PayloadVersion</key>
	<integer>1</integer>
	<key>PayloadIdentifier</key>
	<string>%s</string>
	<key>PayloadUUID</key>
	<string>%s</string>
	<key>PayloadDisplayName</key>
	<string>Printers of %s</string>
	<key>PayloadDescription</key>
	<string>Adds the printers of bdf printd at %s to AirPrint.</string>
	<key>PayloadRemovalDisallowed</key>
	<false/>
</dict>
</plist>
`, esc(reverse), id("profile"), esc(host), esc(host))
	return b.Bytes()
}

// serveProfile serves the configuration profile of the printers at the
// host the request came to: a person opens it on an iPhone or iPad.
func (x *ippServer) serveProfile(w http.ResponseWriter, r *http.Request) {
	if !isSecure(r) {
		// over plain HTTP a machine of the network could put its own
		// certificate in it: the profile is served over https only
		http.Redirect(w, r, "https://"+r.Host+r.URL.Path, http.StatusFound)
		return
	}
	host, portStr, err := net.SplitHostPort(r.Host)
	if err != nil {
		host, portStr = r.Host, strconv.Itoa(x.port)
	}
	port, _ := strconv.Atoi(portStr)
	cert, err := x.certificate()
	if err != nil {
		http.Error(w, "no certificate", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-apple-aspen-config")
	w.Header().Set("Content-Disposition", `attachment; filename="bdf-printd.mobileconfig"`)
	w.Write(x.profile(strings.Trim(host, "[]"), port, cert))
}
