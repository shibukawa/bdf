package docflow

import (
	"net/http"
)

// An entry of an adapter may have a port of its own, besides the
// manager's own server, for a protocol of its own (Kind.Serve: a CUPS
// server takes IPP there, docflow/io/ipp); the job follows the node's
// edges (To) from there. What applications call, the uploaders, the
// inboxes and the downloaders, is the API at /api/<kind>/<name> on the
// manager's own server alone. The process (cmd/docflow) listens on the
// ports of the nodes that have one and serves PortHandler there; a
// serverless deployment has no ports (Options.NoPorts refuses them).

// PortHandler serves a node on its own port: what it does is read from
// the node at each request, so that its settings changing need no
// restart; a node that is gone, or has no port any more, answers 404.
func PortHandler(m *Manager, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, err := m.Node(r.Context(), name)
		if err != nil || n.Port == 0 {
			http.Error(w, "no node listens here any more", http.StatusNotFound)
			return
		}
		if k := kindNamed(n.Kind); k != nil && k.Serve != nil {
			k.Serve(m, n).ServeHTTP(w, r)
			return
		}
		http.Error(w, "no node listens here any more", http.StatusNotFound)
	})
}
