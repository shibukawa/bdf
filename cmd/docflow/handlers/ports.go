package handlers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/shibukawa/bdf/docflow"
)

// ports are the listeners of the entry nodes that have a port of their
// own (an endpoint's REST, a CUPS server's IPP): one per node, opened and
// closed as the nodes change (the manager's OnChange, and every tick in
// case a port could not be opened), and closed with the app. A serverless
// deployment has none (docflow.Options.NoPorts).
type ports struct {
	a       *App
	mu      sync.Mutex
	servers map[int]*portServer
	kicked  chan struct{}
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
}

func tlsKey(n docflow.Node) string { return n.TLSCert + "|" + n.TLSKey }

// closeAll closes every listener (the app closes).
func (p *ports) closeAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
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
