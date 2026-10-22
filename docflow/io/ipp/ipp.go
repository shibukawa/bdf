// Package ipp is the IPP input and output of docflow: the printers a job
// is printed on (the kind printer, over IPP, port 9100 or LPD), CUPS
// servers on ports of their own (the kind cups, an entry), and every
// node's IPP printer at /ipp/print/<name> (an ingress), as docflow's
// kinds of node are registered. A program that imports it serves them.
//
// Every node is an IPP printer at /ipp/print/<name>, unless it turned IPP
// off (Node.NoIPP): a sender (a WMS, a business system, an office PC)
// prints to it as to any printer, over HTTP(S) on the manager's own port,
// and the job follows the topology from there. Only the operations of one
// request are taken (Print-Job, not Create-Job and Send-Document), so
// that a request is all there is to a job: no instance waits for a
// document. /ipp/print itself answers CUPS-Get-Printers and
// CUPS-Get-Classes with every node, as a CUPS server lists its
// destinations. A CUPS server node's own port is a CUPS server with one
// printer, the node, at /printers/<node>: what is printed there follows
// the node's edges.
package ipp

import (
	"io"
	"net/http"
	"strings"

	"github.com/shibukawa/bdf/docflow"
)

// The kinds of node of this package.
const (
	// KindPrinter prints the job on paper: a printer the manager reaches
	// over the network (Node.URL), the way its jobs reach it (Node.Route).
	KindPrinter = "printer"
	// KindCups is an entry for computers and business systems: a CUPS
	// server on a port of its own, with one printer there, the node
	// itself. It forwards what it receives to every node of To.
	KindCups = "cups"
)

func init() {
	docflow.RegisterKind(docflow.Kind{Name: KindCups, Title: "CUPS server", Role: docflow.RoleEntry, PassesOn: true, Check: checkCups, Serve: serveCups})
	docflow.RegisterKind(docflow.Kind{Name: KindPrinter, Role: docflow.RoleExit, Takes: []string{docflow.PayloadBDF}, Fixed: true,
		Check: checkPrinter, Ready: ready, Send: send})
	docflow.RegisterIngress(docflow.Ingress{Name: "ipp", Prefix: Path, Mount: mount, Address: func(n *docflow.Node) string {
		if n.NoIPP || n.Kind == docflow.KindExport {
			return "" // an export is read from, never sent to
		}
		return Path + "/" + n.Name
	}})
}

// mount adds the nodes' IPP printers to the engine's handler.
func mount(m *docflow.Manager, mux *http.ServeMux) {
	s := &server{m: m}
	mux.HandleFunc("POST "+Path+"/{dest}", s.serveIPP)
	mux.HandleFunc("POST "+Path, s.serveIPP)
	mux.HandleFunc("POST "+Path+"/{$}", s.serveIPP)
}

type server struct {
	m *docflow.Manager
}

func or[T comparable](v, def T) T {
	var zero T
	if v == zero {
		return def
	}
	return v
}

func firstOr(s, or string) string {
	if strings.TrimSpace(s) == "" {
		return or
	}
	return s
}

func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for n > 0 && s[n]&0xc0 == 0x80 {
		n--
	}
	return s[:n]
}

func discard(r io.Reader) { io.Copy(io.Discard, io.LimitReader(r, 1<<20)) }
