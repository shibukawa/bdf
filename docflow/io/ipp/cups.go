package ipp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/shibukawa/bdf/docflow"
)

// checkCups takes a CUPS server's port (the first free from
// docflow.FirstPort when none is given; none on a serverless deployment,
// where it is reached at /ipp/print/<name> alone) and its certificate's
// files (none: plain ipp); docflow checked what every entry has. What
// computers print is a document, made .bdf: a CUPS server has no
// accept, output or formats of its own.
func checkCups(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	n.Accept, n.Output, n.Formats = docflow.AcceptFile, docflow.OutputBDF, nil
	n.TLSCert, n.TLSKey = strings.TrimSpace(in.TLSCert), strings.TrimSpace(in.TLSKey)
	if n.Port == 0 && m.Ports() {
		p, err := m.FreePort(ctx, n.Name)
		if err != nil {
			return err
		}
		n.Port = p
	}
	if (n.TLSCert == "") != (n.TLSKey == "") {
		return fmt.Errorf("%w: a CUPS server's certificate takes both tls_cert and tls_key (files)", docflow.ErrInvalid)
	}
	return nil
}

// serveCups serves a CUPS server node on its own port: IPP there as a
// CUPS server with one printer does (the node itself at /printers/<name>
// and at /ipp/print, the server itself at /), so that computers and
// business systems print to it as to a CUPS server, and the job follows
// the node's edges (To) from there.
func serveCups(m *docflow.Manager, n *docflow.Node) http.Handler {
	s := &server{m: m}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"node": n.Name, "kind": n.Kind, "port": n.Port, "to": n.To,
				"ipp": "POST application/ipp to /printers/" + n.Name + " (the one printer here, which forwards to the nodes of to), or to / for the server"})
			return
		}
		dest, ok := cupsDest(r.URL.Path, n.Name)
		if !ok {
			http.Error(w, "no such printer: this server's one printer is /printers/"+n.Name, http.StatusNotFound)
			return
		}
		s.serveIPPTo(w, r, dest, n.Name)
	})
}

// cupsDest is the node an IPP request on a CUPS server's port is to: the
// server itself at / and /admin (CUPS-Get-Printers), and its one printer,
// the node, at /printers/<node>, /classes/<node>, /ipp/print/<node> and
// /ipp/print (a plain IPP printer's path, for clients that are no CUPS).
// Any other node's name is no printer here: what reaches the node goes
// on along its edges.
func cupsDest(path, node string) (dest string, ok bool) {
	path = strings.Trim(path, "/")
	switch path {
	case "", "admin":
		return "", true
	case "ipp/print":
		return node, true
	}
	for _, prefix := range []string{"printers/", "classes/", "ipp/print/"} {
		if rest, found := strings.CutPrefix(path, prefix); found && rest == node {
			return node, true
		}
	}
	return "", false
}
