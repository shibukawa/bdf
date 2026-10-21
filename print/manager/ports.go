package manager

import (
	"net/http"
	"strings"
)

// An entry node of a process may have a port of its own, besides the
// manager's own server: an endpoint takes its REST there (every POST is
// to it), a CUPS server takes IPP there as a CUPS server with one printer
// does (the node itself at /printers/<name> and at /ipp/print, the server
// itself at /), so that computers and business systems print to it as to
// a CUPS server, and the job follows the node's edges (To) from there;
// applications send to an endpoint on a port of its own. The process
// (cmd/printmanager) listens on the ports of the nodes that have one and
// serves PortHandler there; a serverless deployment has no ports
// (Options.NoPorts refuses them).

// PortHandler serves a node on its own port: what it does is read from
// the node at each request, so that its settings changing need no
// restart; a node that is gone, or has no port any more, answers 404.
func PortHandler(m *Manager, name string) http.Handler {
	h := &handler{m: m}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, err := m.Node(r.Context(), name)
		if err != nil || n.Port == 0 {
			http.Error(w, "no node listens here any more", http.StatusNotFound)
			return
		}
		switch n.Kind {
		case KindEndpoint:
			if r.Method != http.MethodPost {
				writeJSON(w, http.StatusOK, map[string]any{"node": name, "kind": n.Kind, "port": n.Port, "post": "a document, or JSON naming a template and its data, with the endpoint's key"})
				return
			}
			r.SetPathValue("name", name)
			h.serveIn(w, r)
		case KindCups:
			if r.Method != http.MethodPost {
				writeJSON(w, http.StatusOK, map[string]any{"node": name, "kind": n.Kind, "port": n.Port, "to": n.To, "ipp": "POST application/ipp to /printers/" + name + " (the one printer here, which forwards to the nodes of to), or to / for the server"})
				return
			}
			dest, ok := cupsDest(r.URL.Path, name)
			if !ok {
				http.Error(w, "no such printer: this server's one printer is /printers/"+name, http.StatusNotFound)
				return
			}
			h.serveIPPTo(w, r, dest, name)
		default:
			http.Error(w, "no node listens here any more", http.StatusNotFound)
		}
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
