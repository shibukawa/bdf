package manager

import (
	"net/http"
	"strings"
)

// An entry node of a process may have a port of its own, besides the
// manager's own server: an endpoint takes its REST there (every POST is
// to it), a CUPS server takes IPP there as a CUPS server does (the nodes
// at /printers/<name>, the classes at /classes/<name>, the server itself
// at /), so that computers and business systems print to it as to a CUPS
// server, and applications send to an endpoint on a port of its own. The
// process (cmd/printmanager) listens on the ports of the nodes that have
// one and serves PortHandler there; a serverless deployment has no ports
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
				writeJSON(w, http.StatusOK, map[string]any{"node": name, "kind": n.Kind, "port": n.Port, "ipp": "POST application/ipp to /printers/<node>, or to / for the server"})
				return
			}
			dest, ok := cupsDest(r.URL.Path)
			if !ok {
				http.Error(w, "no such printer", http.StatusNotFound)
				return
			}
			h.serveIPPTo(w, r, dest)
		default:
			http.Error(w, "no node listens here any more", http.StatusNotFound)
		}
	})
}

// cupsDest is the node an IPP request on a CUPS server's port is to: the
// server itself at /, /ipp/print and /admin, a node at /printers/<name>,
// /classes/<name> or /ipp/print/<name>.
func cupsDest(path string) (dest string, ok bool) {
	path = strings.Trim(path, "/")
	switch path {
	case "", "ipp/print", "admin":
		return "", true
	}
	for _, prefix := range []string{"printers/", "classes/", "ipp/print/"} {
		if rest, found := strings.CutPrefix(path, prefix); found && validName(rest) {
			return rest, true
		}
	}
	return "", false
}
