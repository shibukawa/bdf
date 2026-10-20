package server

import (
	"embed"
	"net/http"
)

//go:embed web/index.html web/app.js
var web embed.FS

// page serves the page of the state and its script. The page only reads
// the API, with the token the viewer gives.
func (s *Server) page(w http.ResponseWriter, r *http.Request) {
	name, typ := "web/index.html", "text/html; charset=utf-8"
	if r.URL.Path == "/app.js" {
		name, typ = "web/app.js", "text/javascript; charset=utf-8"
	}
	b, err := web.ReadFile(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h := w.Header()
	h.Set("Content-Type", typ)
	h.Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'unsafe-inline'; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'; form-action 'none'")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Referrer-Policy", "no-referrer")
	w.Write(b)
}
