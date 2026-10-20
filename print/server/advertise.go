package server

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/print/ipp"
)

// instanceName is the DNS-SD instance name of a queue: its info and the
// server's host, a label of 63 bytes at most.
func (x *ippServer) instanceName(q *ippQueue) string {
	short, _, _ := strings.Cut(x.host, ".")
	name := strings.ReplaceAll(q.info+" @ "+short, ".", " ")
	for len(name) > 63 {
		r := []rune(name)
		name = string(r[:len(r)-1])
	}
	return name
}

// mdnsHost is the host name the printers are advertised on with mDNS:
// one of the server's own ("<host>-bdf.local"), whose addresses it
// answers itself. The system's name is the system responder's: records
// of it from another responder are a conflict, which renames the
// machine.
func (x *ippServer) mdnsHost() string {
	label, _, _ := strings.Cut(x.host, ".")
	return label + "-bdf.local"
}

// pdl is the pdl key of the TXT record: PDF first, so that clients that
// take the first they can send vectors.
const pdl = "application/pdf,image/urf,image/pwg-raster,image/jpeg,image/png"

// serviceTXT is the TXT record of a queue as IPP Everywhere (PWG 5100.14
// §4.2) and AirPrint (whose URF key it must have) ask: its path,
// formats, model, uuid, what it can do and how it is reached.
func (x *ippServer) serviceTXT(q *ippQueue, caps *ipp.Printer, secure bool, port int) []string {
	media, _ := mediaOf(caps)
	color, duplex := true, true
	if caps != nil {
		if len(caps.ColorModes) > 0 && !slices.Contains(caps.ColorModes, "color") {
			color = false
		}
		if len(caps.Sides) > 0 && !slices.Contains(caps.Sides, "two-sided-long-edge") {
			duplex = false
		}
	}
	tf := func(b bool) string {
		if b {
			return "T"
		}
		return "F"
	}
	air := "none"
	if x.cfg.Auth == authBasic {
		air = "username,password"
	}
	scheme := "http"
	if secure {
		scheme = "https"
	}
	ty := "BDF Print Server"
	if caps != nil && caps.MakeModel != "" {
		ty += " (" + caps.MakeModel + ")"
	}
	txt := []string{
		"txtvers=1",
		"qtotal=1",
		"rp=" + strings.TrimPrefix(q.path, "/"),
		"ty=" + clipTXT(ty),
		"product=(BDF Print Server)",
		"usb_MFG=BDF",
		"usb_MDL=Print Server",
		"note=" + clipTXT(q.location),
		"adminurl=" + clipTXT(fmt.Sprintf("%s://%s:%d%s", scheme, x.adminHost(), port, q.path)),
		"priority=50",
		"pdl=" + pdl,
		"URF=" + strings.Join(urfSupported, ","),
		"UUID=" + strings.TrimPrefix(q.uuid, "urn:uuid:"),
		"TLS=1.2",
		"Color=" + tf(color),
		"Duplex=" + tf(duplex),
		"Copies=T",
		"Collate=T",
		"Scan=F",
		"Fax=F",
		"kind=document",
		"air=" + air,
		"mopria-certified=2.0",
		"print_wfds=T",
	}
	if pm := paperMax(media); pm != "" {
		txt = append(txt, "PaperMax="+pm)
	}
	return txt
}

// clipTXT bounds a value of a TXT key: a string of the record has 255
// bytes at most.
func clipTXT(s string) string {
	if len(s) <= 200 {
		return s
	}
	r := []rune(s)
	for len(string(r)) > 200 {
		r = r[:len(r)-1]
	}
	return string(r)
}

// advertisements are the services of the queues: _ipps._tcp for each,
// and _ipp._tcp unless ipp is refused; the subtypes _print (IPP
// Everywhere) and _universal (AirPrint).
func (x *ippServer) advertisements(caps func(queue string) *ipp.Printer, port int) []ipp.Advertisement {
	var out []ipp.Advertisement
	for _, q := range x.order {
		c := caps(q.name)
		name := x.instanceName(q)
		if !x.secureOnly {
			out = append(out, ipp.Advertisement{Instance: name, Type: "_ipp._tcp", Subtypes: []string{"_print", "_universal"}, Port: port,
				TXT: x.serviceTXT(q, c, false, port)})
		}
		out = append(out, ipp.Advertisement{Instance: name, Type: "_ipps._tcp", Subtypes: []string{"_print", "_universal"}, Port: port,
			TXT: x.serviceTXT(q, c, true, port)})
	}
	return out
}

// advertise advertises the queues with DNS-SD on the local links until
// ctx ends, unless dns-sd is false.
func (x *ippServer) advertise(ctx context.Context) error {
	if x.cfg.DNSSD != nil && !*x.cfg.DNSSD {
		return nil
	}
	ads := x.advertisements(x.queueCaps, x.port)
	// the instances' names are probed in the background: the server
	// serves meanwhile
	go func() {
		r, err := ipp.Advertise(ctx, ads, ipp.ResponderOptions{Host: x.mdnsHost(), Log: x.s.log})
		if err != nil {
			x.s.log.Warn("ipp: DNS-SD", "err", err)
			return
		}
		x.s.log.Info("ipp: advertised", "host", r.Host(), "services", strings.Join(r.Names(), ", "))
		x.mu.Lock()
		x.responder = r
		x.mu.Unlock()
		x.readvertise()
	}()
	return nil
}

// adminHost is the host of adminurl: the name the clients of the zone
// (ZoneRecords) or of mDNS reach the server by.
func (x *ippServer) adminHost() string {
	if x.zone {
		return x.host
	}
	return x.mdnsHost()
}

// readvertise announces the TXT records again, once the printers have
// told what they can do (colour, both sides, the largest paper).
func (x *ippServer) readvertise() {
	x.mu.Lock()
	r := x.responder
	x.mu.Unlock()
	if r == nil {
		return
	}
	var txt [][]string
	for _, a := range x.advertisements(x.queueCaps, x.port) {
		txt = append(txt, a.TXT)
	}
	r.SetTXT(txt)
}

// ZoneRecords are the records of the queues for a DNS zone (RFC 6763 §11,
// wide-area DNS-SD), in the zone file's format: clients of other subnets
// find the printers in the domain, as they find them with mDNS on theirs.
// The host is the IPP server's (ipp.host), a name of the zone; caps give
// the queues' printers' capabilities (CheckPrinters), nil when unknown.
func ZoneRecords(cfg *Config, domain string, port int, caps map[string]*ipp.Printer) (string, error) {
	if cfg.IPP == nil {
		return "", fmt.Errorf("the configuration has no ipp")
	}
	domain = strings.TrimSuffix(domain, ".") + "."
	x := &ippServer{s: &Server{cfg: cfg}, cfg: cfg.IPP, queues: map[string]*ippQueue{}, secureOnly: cfg.IPP.Auth == authBasic || cfg.IPP.RequireTLS, zone: true}
	if err := x.init(); err != nil {
		return "", err
	}
	if !strings.Contains(strings.TrimSuffix(x.host, "."), ".") || strings.HasSuffix(x.host, ".local") {
		return "", fmt.Errorf("ipp.host %q is not a name of %s: set it to the server's name in the domain", x.host, domain)
	}
	host := strings.TrimSuffix(x.host, ".") + "."
	queueCaps := func(q string) *ipp.Printer {
		qc := cfg.Queues[q]
		if len(qc.Printers) == 0 {
			return nil
		}
		return caps[qc.Printers[0]]
	}
	ads := x.advertisements(queueCaps, port)
	var b strings.Builder
	fmt.Fprintf(&b, "; DNS-SD records of the printers of bdf printd (RFC 6763 §11) for %s\n", domain)
	fmt.Fprintf(&b, "b._dns-sd._udp.%s\tIN PTR\t%s\n", domain, domain)
	fmt.Fprintf(&b, "lb._dns-sd._udp.%s\tIN PTR\t%s\n", domain, domain)
	types := map[string]bool{}
	for _, a := range ads {
		types[a.Type] = true
	}
	for _, t := range sortedKeys(types) {
		fmt.Fprintf(&b, "_services._dns-sd._udp.%s\tIN PTR\t%s.%s\n", domain, t, domain)
	}
	for _, a := range ads {
		inst := zoneEscape(a.Instance) + "." + a.Type + "." + domain
		fmt.Fprintf(&b, "\n%s.%s\tIN PTR\t%s\n", a.Type, domain, inst)
		for _, s := range a.Subtypes {
			fmt.Fprintf(&b, "%s._sub.%s.%s\tIN PTR\t%s\n", s, a.Type, domain, inst)
		}
		fmt.Fprintf(&b, "%s\tIN SRV\t0 0 %d %s\n", inst, a.Port, host)
		fmt.Fprintf(&b, "%s\tIN TXT\t", inst)
		for i, t := range a.TXT {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.Quote(t))
		}
		b.WriteByte('\n')
	}
	return b.String(), nil
}

// zoneEscape escapes a label for a zone file.
func zoneEscape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == ' ':
			b.WriteString(`\032`)
		case r == '.' || r == '\\' || r == '"' || r == ';' || r == '(' || r == ')' || r == '@' || r == '$':
			b.WriteByte('\\')
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

var _ = sort.Strings
