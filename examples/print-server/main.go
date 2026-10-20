// Command print-server is a sample (docs/examples/print-server): a
// business application that prints invoices through a print server it
// embeds (package print/server, the server of bdf printd). It fills an
// invoice template (an Excel workbook with tags, package form) with what
// is typed on its page, puts the .bdf in the spool's queue, and shows the
// jobs as the spool sends them and the printer prints them (the server's
// /api/v1/events).
//
// The spool keeps the jobs in one file, sends them to the printer, and
// retries a printer that is not reached; a job is printed once even when
// the program is stopped while it is being sent. Without a printer at
// hand, -fake-printer prints on bdf's fake IPP printer, in the same
// process.
//
//	cd examples/print-server && go run . -fake-printer
//	open http://127.0.0.1:8083/
//
//	go run . -printer ipps://printer.example.local/ipp/print
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/form"
	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/server"
	"github.com/shibukawa/bdf/print/spool"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:8083", "address to listen on")
	printer := flag.String("printer", "", "the printer's URL (ipp://…, ipps://…)")
	fake := flag.Bool("fake-printer", false, "print on bdf's fake IPP printer, in this process")
	dir := flag.String("dir", "./data", "directory of the spool's file")
	tmpl := flag.String("template", "../../converter/xlsx/testdata/form-rows.xlsx", "the invoice template")
	fonts := flag.String("font-dir", "../../converter/pptx/testdata/fonts", "directory of the template's fonts")
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	url := *printer
	switch {
	case *fake && url == "":
		u, err := fakePrinter()
		if err != nil {
			log.Fatal(err)
		}
		url = u
	case url == "":
		log.Fatal("no printer: -printer ipp://… or -fake-printer")
	}
	tpl, err := os.ReadFile(*tmpl)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(*dir, 0o755); err != nil {
		log.Fatal(err)
	}

	// the print server's configuration, as bdf printd reads it from YAML;
	// it serves on this program's address, so it listens nowhere itself
	cfg, err := server.ParseConfig([]byte(fmt.Sprintf(`
printers: { printer: { url: %q } }
queues: { invoices: { printers: [printer], defaults: { sides: two-sided-long-edge } } }
`, url)))
	if err != nil {
		log.Fatal(err)
	}
	cfg.Spool = filepath.Join(*dir, "spool.db")
	srv, err := server.New(ctx, cfg, server.Options{Log: logger})
	if err != nil {
		log.Fatal(err)
	}
	srv.Spool().Start()
	defer srv.Close()

	a := &app{spool: srv.Spool(), template: tpl, fonts: *fonts}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", a.page)
	mux.HandleFunc("POST /invoices", a.invoice)
	// the print server's API: the page lists the jobs and follows them
	mux.Handle("/api/", srv.Handler())

	hs := &http.Server{Addr: *listen, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	go func() {
		<-ctx.Done()
		sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		hs.Shutdown(sctx)
	}()
	logger.Info("serving", "url", "http://"+*listen+"/", "printer", url)
	if err := hs.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// fakePrinter serves bdf's fake IPP printer on a port of this machine and
// returns its URL. It takes a few seconds to print a job.
func fakePrinter() (string, error) {
	ps, _ := fakeprinter.LookupPreset("office-mfp")
	p := ps.New("office")
	p.Processing = 3 * time.Second
	s, err := p.ServeIPP("127.0.0.1:0")
	if err != nil {
		return "", err
	}
	return s.URL, nil
}

type app struct {
	spool    *spool.Spool
	template []byte
	fonts    string
	next     atomic.Int64 // the number of the last invoice
}

// invoice fills the template with the form's data and submits the
// invoice to the queue.
func (a *app) invoice(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	items, err := parseItems(r.FormValue("items"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	no := fmt.Sprintf("A-%d", a.next.Add(1))
	data := map[string]any{"no": no, "customer": r.FormValue("customer"), "note": r.FormValue("note"), "items": items}
	res, err := form.Fill(bytes.NewReader(a.template), int64(len(a.template)), &form.Options{Data: data,
		Convert: &converter.Options{FontDirs: []string{a.fonts}, NoSystemFonts: true}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var doc bytes.Buffer
	if err := res.Doc.WriteSingle(&doc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// the application submits for the user who asked: here, the one the
	// page names
	j, err := a.spool.Submit(r.Context(), spool.Submission{Queue: "invoices", Name: "invoice-" + no + ".bdf",
		User: r.FormValue("user"), Document: &doc})
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/?job=%d", j.ID), http.StatusSeeOther)
}

// parseItems reads lines of "name, quantity, price".
func parseItems(s string) ([]map[string]any, error) {
	var items []map[string]any
	for i, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		f := strings.Split(line, ",")
		if len(f) != 3 {
			return nil, fmt.Errorf("line %d: name, quantity, price", i+1)
		}
		qty, err1 := strconv.Atoi(strings.TrimSpace(f[1]))
		price, err2 := strconv.Atoi(strings.TrimSpace(f[2]))
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d: the quantity and the price are numbers", i+1)
		}
		items = append(items, map[string]any{"name": strings.TrimSpace(f[0]), "qty": qty, "price": price})
	}
	return items, nil
}

func (a *app) page(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pageTemplate.Execute(w, nil)
}

var pageTemplate = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Print server sample</title>
<style>
:root { color-scheme: light dark; --line: #8884; }
body { font: 14px/1.5 system-ui, sans-serif; margin: 0; padding: 16px; max-width: 960px; }
h1 { font-size: 18px; }
form { display: grid; gap: 8px; max-width: 480px; }
label { display: grid; gap: 2px; }
input, textarea, button { font: inherit; padding: 4px 8px; }
table { border-collapse: collapse; width: 100%; margin-top: 16px; }
th, td { text-align: left; padding: 4px 8px; border-bottom: 1px solid var(--line); }
.scroll { overflow-x: auto; }
</style>
</head>
<body>
<h1>Invoices</h1>
<form method="post" action="/invoices">
  <label>Customer <input name="customer" value="ACME" required></label>
  <label>Items (name, quantity, price; a line each)
    <textarea name="items" rows="5">Widget, 2, 100
Gadget, 3, 250</textarea></label>
  <label>Note <input name="note" value="Thank you."></label>
  <label>User <input name="user" value="sato"></label>
  <button>Print</button>
</form>
<div class="scroll"><table>
  <thead><tr><th>Job</th><th>Name</th><th>User</th><th>State</th><th>Printer</th></tr></thead>
  <tbody id="jobs"></tbody>
</table></div>
<script>
const jobs = new Map();
function render() {
  const rows = [...jobs.values()].sort((a, b) => b.id - a.id).map((j) => {
    const tr = document.createElement("tr");
    for (const v of [j.id, j.name, j.user, j.state + (j.reasons ? " (" + j.reasons.join(", ") + ")" : ""), j.printer ? j.printer.name : ""]) {
      const td = document.createElement("td");
      td.textContent = v;
      tr.append(td);
    }
    return tr;
  });
  document.getElementById("jobs").replaceChildren(...rows);
}
fetch("/api/v1/jobs?state=all").then((r) => r.json()).then((l) => {
  for (const j of l.jobs) jobs.set(j.id, j);
  render();
});
new EventSource("/api/v1/events").addEventListener("job", (e) => {
  const j = JSON.parse(e.data);
  jobs.set(j.id, j);
  render();
});
</script>
</body>
</html>
`))
