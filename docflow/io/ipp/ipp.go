// Package ipp is the IPP input and output of docflow: the printers a job
// is printed on (the kind printer, over IPP, port 9100 or LPD), the CUPS
// servers jobs are printed to (the kind cups, an entry: a source), and
// the queues of CUPS servers jobs are printed on (the kind cups-out, an
// exit: a sink), as docflow's kinds of node are registered. A program that imports it
// serves them.
//
// A CUPS server node is the one way IPP enters the topology (the other
// entries, an uploader or an S3 trigger, take files over HTTP, at the
// API). Its own
// port is a CUPS server with one printer, the node, at /printers/<node>;
// the node is also an IPP printer at /ipp/print/<name> of the manager's
// own server (an ingress: on a serverless deployment, with no ports, the
// only one). A sender (a WMS, a business system, an office PC) prints to
// it as to any printer, and the job follows the node's edges from there.
// The manager's own server takes the operations of one request alone
// (Print-Job, not Create-Job and Send-Document), so that a request is all
// there is to a job: no instance waits for a document. A node's own port,
// where one process listens, takes Create-Job and Send-Document too
// (created.go), as IPP Everywhere asks: there it is an AirPrint, IPP
// Everywhere and Mopria printer that Macs, iPhones, Android and Windows
// add without a driver. /ipp/print itself answers CUPS-Get-Printers with
// the CUPS server nodes, as a CUPS server lists its destinations.
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
	// KindCupsOut is an exit: a queue of a CUPS server the job is printed
	// to over IPP (cupsout.go).
	KindCupsOut = "cups-out"
)

func init() {
	docflow.RegisterKind(docflow.Kind{Name: KindCups, Title: "CUPS server (source)", Role: docflow.RoleEntry, PassesOn: true, Check: checkCups, Serve: serveCups})
	docflow.RegisterKind(docflow.Kind{Name: KindCupsOut, Title: "CUPS server (sink)", Role: docflow.RoleExit, Takes: []string{docflow.PayloadBDF}, Fixed: true,
		Check: checkCupsOut, Send: sendCupsOut, Discover: discoverQueues})
	docflow.RegisterKind(docflow.Kind{Name: KindPrinter, Role: docflow.RoleExit, Takes: []string{docflow.PayloadBDF}, Fixed: true,
		Check: checkPrinter, Ready: ready, Send: send, Discover: discoverPrinters})
	docflow.RegisterIngress(docflow.Ingress{Name: "ipp", Prefix: Path, Mount: mount, Address: func(n *docflow.Node) string {
		if n.Kind != KindCups {
			return "" // IPP enters at a CUPS server alone
		}
		return Path + "/" + n.Name
	}})
}

// mount adds the CUPS server nodes' IPP printers to the engine's handler.
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
