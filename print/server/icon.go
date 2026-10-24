package server

import (
	"html/template"
	"net/http"

	"github.com/shibukawa/bdf/print/ipp"
)

// icon serves an icon of printer-icons.
func (x *ippServer) icon(w http.ResponseWriter, r *http.Request) {
	ipp.ServeIcon(w, r, r.PathValue("file"))
}

var ippPage = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{.Title}}</title>
<style>
:root{color-scheme:light dark;--fg:#1d2733;--bg:#fff;--muted:#5b6876;--line:#d5dbe1}
@media (prefers-color-scheme:dark){:root{--fg:#e6ebf0;--bg:#14191f;--muted:#9aa7b4;--line:#2e3946}}
body{font:15px/1.5 system-ui,sans-serif;color:var(--fg);background:var(--bg);margin:0;padding:24px 16px;max-width:760px;margin:auto}
h1{font-size:20px;margin:0 0 16px;display:flex;gap:12px;align-items:center}
table{border-collapse:collapse;width:100%}td,th{text-align:left;padding:8px 6px;border-bottom:1px solid var(--line)}
th{font-weight:600;color:var(--muted);font-size:13px}code{font-size:13px;word-break:break-all}
</style></head><body>
<h1><img src="/icons/48.png" width="32" height="32" alt="">{{.Title}}</h1>
<table><tr><th>Printer</th><th>State</th><th>Jobs</th><th>Address</th></tr>
{{range .Queues}}<tr><td>{{.Info}}{{if .Location}}<br><small>{{.Location}}</small>{{end}}</td><td>{{.State}}{{if .Message}}<br><small>{{.Message}}</small>{{end}}</td><td>{{.Jobs}}</td><td><code>{{.URI}}</code></td></tr>
{{end}}</table>
{{if .Secure}}<p><a href="/airprint.mobileconfig">Add these printers to an iPhone or iPad</a> that does not list them (a configuration profile: open it in Safari, then install it in Settings).</p>{{else}}<p>To add these printers to an iPhone or iPad that does not list them, open this page over <a href="https://{{.Host}}/">https</a>.</p>{{end}}
</body></html>
`))

// page shows the printers and their states: printer-more-info.
func (x *ippServer) page(w http.ResponseWriter, r *http.Request) {
	type row struct {
		Info, Location, State, Message, URI string
		Jobs                                int
	}
	data := struct {
		Title, Host string
		Secure      bool
		Queues      []row
	}{Title: "bdf printd on " + x.host, Host: r.Host, Secure: isSecure(r)}
	scheme := "ipp"
	if isSecure(r) {
		scheme = "ipps"
	}
	for _, q := range x.order {
		if name := r.PathValue("queue"); name != "" && name != q.name {
			continue
		}
		st, _, msg, accepting := x.queueState(q.name)
		state := map[ipp.PrinterState]string{ipp.PrinterIdle: "idle", ipp.PrinterProcessing: "printing", ipp.PrinterStopped: "stopped"}[st]
		if !accepting {
			state += ", not accepting jobs"
		}
		n := 0
		if info, err := x.s.sp.Queue(q.name); err == nil {
			for _, c := range info.Jobs {
				n += c
			}
		}
		data.Queues = append(data.Queues, row{Info: q.info, Location: q.location, State: state, Message: msg, Jobs: n, URI: scheme + "://" + r.Host + q.path})
	}
	if len(data.Queues) == 0 {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; style-src 'unsafe-inline'")
	ippPage.Execute(w, data)
}
