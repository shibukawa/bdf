package ipp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raw"
)

// Bonjour: the IPP printers DNS-SD finds on the local links of the
// process, AirPrint printers and the queues computers and CUPS servers
// share among them. The editor picks a printer's URL, or a CUPS server
// (sink)'s queue, from them (Kind.Discover), and a node picked so keeps
// the instance name (Node.Bonjour): when its address stops answering (a
// printer's address from DHCP changes), it is looked for again by the
// name, and the job printed where it is now. A serverless deployment has
// no local links: the address is given there.

// browseTimeout is how long Bonjour is asked.
const browseTimeout = 2 * time.Second

// browse asks Bonjour for the IPP printers (tests change it).
var browse = func(ctx context.Context) ([]ipp.Service, error) { return ipp.Browse(ctx, browseTimeout) }

// Found is a printer, or a queue of a CUPS server, that Bonjour found, as
// the editor picks from.
type Found struct {
	// Name is the instance name it is advertised by ("Canon MF Series"),
	// what a node picked from it keeps.
	Name string `json:"name"`
	// URL is where it is printed to: ipp or ipps, its address and port,
	// and its path (rp).
	URL string `json:"url"`
	// Host is its host name on the local links ("printer.local").
	Host string `json:"host"`
	// Model and Note are what it says of itself (the TXT keys ty and
	// note: the note is often where it is).
	Model string `json:"model,omitempty"`
	Note  string `json:"note,omitempty"`
	// Queue tells a queue of a CUPS server (/printers/<queue>,
	// /classes/<class>): a printer a computer or a CUPS server shares.
	Queue bool `json:"queue,omitempty"`
	// AirPrint tells a printer that takes Apple raster (URF), as AirPrint
	// printers do; Secure, one found as _ipps._tcp (TLS).
	AirPrint bool `json:"airprint,omitempty"`
	Secure   bool `json:"secure,omitempty"`
}

// discoverPrinters lists the IPP printers Bonjour finds, for a printer.
func discoverPrinters(ctx context.Context, m *docflow.Manager, in docflow.Node) (any, error) {
	return bonjour(ctx, m, false)
}

// bonjour lists what Bonjour finds on the local links, sorted by name: the
// queues of CUPS servers alone, or every IPP printer. The CUPS server
// nodes of this manager are left out: they are reached along edges.
func bonjour(ctx context.Context, m *docflow.Manager, queues bool) ([]Found, error) {
	if !m.Ports() {
		return nil, fmt.Errorf("%w: Bonjour finds printers on the local network, which a serverless deployment does not have: give the address", docflow.ErrUnsupported)
	}
	ss, err := browse(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: Bonjour: %v", docflow.ErrPrint, err)
	}
	own := map[string]bool{}
	if nodes, err := m.Nodes(ctx); err == nil {
		for _, n := range nodes {
			if n.Kind == KindCups && n.Port != 0 {
				own[ownKey(strings.TrimPrefix(printerUUID(n.Name), "urn:uuid:"), n.Port)] = true
			}
		}
	}
	out := []Found{}
	for _, s := range ss {
		if own[ownKey(s.TXT["UUID"], s.Port)] {
			continue
		}
		_, qerr := queueOf(s.URL)
		f := Found{Name: s.Name, URL: s.URL, Host: s.Host, Model: s.TXT["ty"], Note: s.TXT["note"], Queue: qerr == nil, AirPrint: s.TXT["URF"] != "", Secure: s.Secure}
		if queues && !f.Queue {
			continue
		}
		out = append(out, f)
	}
	return out, nil
}

func ownKey(uuid string, port int) string { return strings.ToLower(uuid) + "@" + strconv.Itoa(port) }

// checkBonjour takes the instance name a node's URL was picked by: an
// IPP printer's (Bonjour does not list 9100 and LPD), a label of DNS at
// most.
func checkBonjour(in docflow.Node, n *docflow.Node) error {
	n.Bonjour = strings.TrimSpace(in.Bonjour)
	switch {
	case n.Bonjour == "":
	case len(n.Bonjour) > 63:
		return fmt.Errorf("%w: bonjour %q: an instance name is 63 bytes at most", docflow.ErrInvalid, n.Bonjour)
	case raw.IsURL(n.URL):
		return fmt.Errorf("%w: bonjour: Bonjour finds IPP printers, not %s", docflow.ErrInvalid, n.URL)
	}
	return nil
}

// moved is where the printers picked from Bonjour were found again after
// their addresses stopped answering: the node's name and URL → the URL
// now.
var moved sync.Map

// printFound prints to the IPP printer of a node (a printer, a CUPS
// server's queue), the client of a URL made by client. One picked from
// Bonjour that takes no connection at its address is looked for by its
// name, and, found at another, printed to there: nothing reached the
// first address, so nothing is printed twice. The address found is
// kept, for the next jobs to go there first.
func printFound(ctx context.Context, o *docflow.Out, client func(url string) (*ipp.Client, error)) (string, error) {
	key := o.Node.Name + "\x00" + o.Node.URL
	at := o.Node.URL
	if v, ok := moved.Load(key); ok {
		at = v.(string)
	}
	c, err := client(at)
	if err != nil {
		return "", err
	}
	id, err := printTo(ctx, o, c)
	if err == nil || o.Node.Bonjour == "" || !o.Manager().Ports() || !unreachable(err) {
		return id, err
	}
	now, ok := findAgain(ctx, o.Node)
	if !ok || now == at {
		return id, err
	}
	if c, err = client(now); err != nil {
		return "", err
	}
	if id, err = printTo(ctx, o, c); err == nil || !unreachable(err) {
		moved.Store(key, now)
	}
	return id, err
}

// findAgain is the URL a node's printer picked from Bonjour is at now: the
// service of its instance name (a CUPS server's queue, for a CUPS server
// (sink)).
func findAgain(ctx context.Context, n *docflow.Node) (string, bool) {
	ss, err := browse(ctx)
	if err != nil {
		return "", false
	}
	for _, s := range ss {
		if !strings.EqualFold(s.Name, n.Bonjour) {
			continue
		}
		if _, err := queueOf(s.URL); n.Kind == KindCupsOut && err != nil {
			continue
		}
		return s.URL, true
	}
	return "", false
}

// unreachable tells the error of a printer that took no connection: the
// address answers no more.
func unreachable(err error) bool {
	var op *net.OpError
	return errors.As(err, &op) && op.Op == "dial"
}
