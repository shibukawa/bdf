package handlers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/bdf/docflow"
	docflowipp "github.com/shibukawa/bdf/docflow/io/ipp"
	"github.com/shibukawa/bdf/print/ipp"
)

// ports are the listeners of the entry nodes that have a port of their
// own (a CUPS server's IPP): one per node, opened and closed as the nodes
// change (the manager's OnChange, and every tick in case a port could not
// be opened), and closed with the app. A serverless
// deployment has none (docflow.Options.NoPorts). The CUPS servers that
// say so (dnssd) are advertised with DNS-SD on the local links while
// their ports are open.
type ports struct {
	a       *App
	mu      sync.Mutex
	servers map[int]*portServer
	kicked  chan struct{}
	// ads are the services advertised now, as a key, and adv their
	// responder's; mdns configures the responders (tests: a port of
	// their own on the loopback).
	ads  string
	adv  *advertising
	mdns ipp.ResponderOptions
}

// advertising is a responder of mDNS announcing services: stop says
// goodbye, done is closed once it did.
type advertising struct {
	stop context.CancelFunc
	done chan struct{}
}

type portServer struct {
	name, tls string
	srv       *http.Server
}

func newPorts(a *App) *ports {
	return &ports{a: a, servers: map[int]*portServer{}, kicked: make(chan struct{}, 1)}
}

// kick asks for a reconcile soon (the nodes changed).
func (p *ports) kick() {
	select {
	case p.kicked <- struct{}{}:
	default:
	}
}

// run reconciles now, when kicked and every interval, until ctx ends.
func (p *ports) run(ctx context.Context, every time.Duration) {
	if every <= 0 {
		every = 30 * time.Second
	}
	t := time.NewTicker(every)
	defer t.Stop()
	p.reconcile(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.kicked:
			p.reconcile(ctx)
		case <-t.C:
			p.reconcile(ctx)
		}
	}
}

// reconcile opens the listeners of the nodes with a port and closes
// those of the nodes that lost theirs (or changed their certificate).
func (p *ports) reconcile(ctx context.Context) {
	nodes, err := p.a.m.Nodes(ctx)
	if err != nil {
		p.a.log.Warn("ports: nodes", "err", err)
		return
	}
	want := map[int]docflow.Node{}
	for _, n := range nodes {
		if n.Port > 0 && docflow.ServesPort(n.Kind) {
			want[n.Port] = n
		}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for port, s := range p.servers {
		if n, ok := want[port]; !ok || n.Name != s.name || tlsKey(n) != s.tls {
			s.srv.Close()
			delete(p.servers, port)
			p.a.log.Info("port closed", "node", s.name, "port", port)
		}
	}
	for port, n := range want {
		if _, ok := p.servers[port]; ok {
			continue
		}
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			p.a.log.Warn("port", "node", n.Name, "port", port, "err", err)
			continue
		}
		srv := &http.Server{Handler: docflow.PortHandler(p.a.m, n.Name), ReadHeaderTimeout: 10 * time.Second}
		p.servers[port] = &portServer{name: n.Name, tls: tlsKey(n), srv: srv}
		go func() {
			var err error
			if n.TLSCert != "" {
				err = srv.ServeTLS(ln, n.TLSCert, n.TLSKey)
			} else {
				err = srv.Serve(ln)
			}
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				p.a.log.Warn("port", "node", n.Name, "port", port, "err", err)
			}
		}()
		p.a.log.Info("port open", "node", n.Name, "kind", n.Kind, "port", port, "tls", n.TLSCert != "")
	}
	var open []docflow.Node
	for port, n := range want {
		if _, ok := p.servers[port]; ok {
			open = append(open, n)
		}
	}
	slices.SortFunc(open, func(a, b docflow.Node) int { return strings.Compare(a.Name, b.Name) })
	p.advertise(ctx, open)
}

// advertise announces the CUPS servers of nodes (those listening) that
// say so, with DNS-SD, on a host name of the process's own
// ("<computer>-docflow.local"): the computer's own is its system
// responder's, which takes another's records of it as a conflict and
// renames the computer. When they change, the old services say goodbye
// before the new are announced (a goodbye after would take the new
// away). A responder that cannot start is tried again at the next
// reconcile.
func (p *ports) advertise(ctx context.Context, nodes []docflow.Node) {
	machine := machineName()
	ads := docflowipp.Advertisements(p.a.m, nodes, machine)
	key := fmt.Sprint(ads)
	if key == p.ads {
		return
	}
	p.stopAdvertising()
	p.ads = key
	if len(ads) == 0 {
		return
	}
	actx, stop := context.WithCancel(ctx)
	adv := &advertising{stop: stop, done: make(chan struct{})}
	p.adv = adv
	o := p.mdns
	o.Host, o.Log = machine+"-docflow.local", p.a.log
	go func() {
		defer close(adv.done)
		r, err := ipp.Advertise(actx, ads, o)
		if err != nil {
			p.a.log.Warn("DNS-SD", "err", err)
			p.mu.Lock()
			if p.adv == adv {
				p.ads = ""
			}
			p.mu.Unlock()
			return
		}
		p.a.log.Info("DNS-SD: advertised", "host", r.Host(), "services", strings.Join(r.Names(), ", "))
		<-actx.Done()
		r.Close()
	}()
}

// stopAdvertising says goodbye to the services advertised, and waits
// until it is said.
func (p *ports) stopAdvertising() {
	if p.adv == nil {
		return
	}
	p.adv.stop()
	<-p.adv.done
	p.adv, p.ads = nil, ""
}

// machineName is the computer's name, a label ("office-mac"): what the
// services are named after and their host name is made of.
func machineName() string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "docflow"
	}
	h, _, _ = strings.Cut(h, ".")
	// the host name's label ("<name>-docflow") is 63 bytes at most
	for len(h) > 55 {
		r := []rune(h)
		h = string(r[:len(r)-1])
	}
	return h
}

func tlsKey(n docflow.Node) string { return n.TLSCert + "|" + n.TLSKey }

// closeAll closes every listener (the app closes), and says goodbye to
// the services advertised.
func (p *ports) closeAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stopAdvertising()
	for port, s := range p.servers {
		s.srv.Close()
		delete(p.servers, port)
	}
}

// open tells the ports open now, in order.
func (p *ports) open() []int {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]int, 0, len(p.servers))
	for port := range p.servers {
		out = append(out, port)
	}
	slices.Sort(out)
	return out
}
