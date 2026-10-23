package ipp_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/docflow/internal/flowtest"
	dipp "github.com/shibukawa/bdf/docflow/io/ipp"
	"github.com/shibukawa/bdf/print/ipp"
)

// A CUPS server node answers as an AirPrint printer does, so that a Mac
// makes an AirPrint queue of it (which sends PDF) and not a PostScript
// one (whose jobs bdf cannot read): Apple and PWG rasters, PDF
// preferred, the paper, a page about it; the media-col-database when
// asked by name. A document format it does not take is refused.
func TestAirPrintAttributes(t *testing.T) {
	e := flowtest.New(t, nil)
	ctx := context.Background()
	e.Put(docflow.Node{Name: "keep", Kind: docflow.KindSink},
		docflow.Node{Name: "office", Kind: dipp.KindCups, Title: "Office", To: []string{"keep"}})
	cups := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		docflow.PortHandler(e.M(), "office").ServeHTTP(w, r)
	}))
	defer cups.Close()
	host := strings.TrimPrefix(cups.URL, "http://")
	c := &ipp.Client{URL: "ipp://" + host + "/printers/office"}
	ask := func(requested ...string) *ipp.Group {
		t.Helper()
		m := c.NewRequest(ipp.OpGetPrinterAttributes)
		if len(requested) > 0 {
			m.Add(ipp.TagOperation, ipp.Attr("requested-attributes", ipp.Keywords(requested)...))
		}
		res, err := c.Do(ctx, m, nil)
		if err != nil {
			t.Fatal(err)
		}
		return res.Group(ipp.TagPrinter)
	}
	all := ask()
	text := func(name string) string {
		if a := all.Get(name); a != nil {
			return a.Text()
		}
		return ""
	}
	for name, want := range map[string]string{
		"printer-uri-supported":     "ipp://" + host + "/printers/office",
		"printer-more-info":         "http://" + host + "/",
		"printer-make-and-model":    "BDF Docflow",
		"printer-info":              "Office",
		"document-format-preferred": "application/pdf",
		"media-default":             "iso_a4_210x297mm",
		"printer-kind":              "document",
	} {
		if got := text(name); got != want {
			t.Errorf("%s = %q, want %q", name, got, want)
		}
	}
	if a := all.Get("urf-supported"); a == nil || !slices.Contains(a.Strings(), "SRGB24") || !slices.Contains(a.Strings(), "RS150-300-600") {
		t.Errorf("urf-supported %v", a)
	}
	if a := all.Get("document-format-supported"); a == nil || !slices.Contains(a.Strings(), "image/urf") || !slices.Contains(a.Strings(), "application/pdf") ||
		slices.Contains(a.Strings(), "application/postscript") {
		t.Errorf("document-format-supported %v", a)
	}
	if a := all.Get("media-supported"); a == nil || len(a.Values) != 11 {
		t.Errorf("media-supported %v", a)
	}
	if all.Get("media-col-database") != nil {
		t.Error("media-col-database without being asked")
	}
	// asked by name, it alone
	db := ask("media-col-database")
	if a := db.Get("media-col-database"); a == nil || len(a.Values) != 11 || len(db.Attrs) != 1 {
		t.Errorf("media-col-database alone: %v", db.Attrs)
	}
	// a group of them
	tmpl := ask("job-template")
	if tmpl.Get("media-supported") == nil || tmpl.Get("printer-name") != nil {
		t.Errorf("job-template: %v", tmpl.Attrs)
	}
	// PostScript, which a Mac's generic driver sends, is refused before
	// a job is
	m := c.NewRequest(ipp.OpGetPrinterAttributes)
	m.Add(ipp.TagOperation, ipp.Attr("document-format", ipp.MimeType("application/postscript")))
	var se *ipp.StatusError
	if _, err := c.Do(ctx, m, nil); !errors.As(err, &se) || se.Status != ipp.StatusClientErrorDocumentFormatNotSupported {
		t.Errorf("document-format application/postscript: %v", err)
	}
}

// The CUPS server nodes that say so are advertised with DNS-SD as
// AirPrint printers: _ipp._tcp, or _ipps._tcp with a certificate, with
// the subtypes of IPP Everywhere and AirPrint and the TXT keys their
// clients read; a node that does not say so is not.
func TestAdvertisements(t *testing.T) {
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "keep", Kind: docflow.KindSink},
		docflow.Node{Name: "office", Kind: dipp.KindCups, Title: "Office 2.F", To: []string{"keep"}, Port: 16310, DNSSD: true},
		docflow.Node{Name: "secure", Kind: dipp.KindCups, To: []string{"keep"}, Port: 16311, DNSSD: true, TLSCert: "c.crt", TLSKey: "c.key"},
		docflow.Node{Name: "quiet", Kind: dipp.KindCups, To: []string{"keep"}, Port: 16312})
	nodes, err := e.M().Nodes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	ads := dipp.Advertisements(e.M(), nodes, "mac")
	if len(ads) != 2 {
		t.Fatalf("advertisements %+v", ads)
	}
	office, secure := ads[0], ads[1]
	if office.Instance != "Office 2 F @ mac" || office.Type != "_ipp._tcp" || office.Port != 16310 || !slices.Equal(office.Subtypes, []string{"_print", "_universal"}) {
		t.Errorf("office %+v", office)
	}
	for _, kv := range []string{"rp=printers/office", "pdl=application/pdf,image/urf,image/pwg-raster,image/jpeg,image/png", "URF=V1.4,W8,SRGB24,RS150-300-600,CP1,IS1,MT1,PQ4",
		"ty=BDF Docflow", "note=Office 2.F", "kind=document", "air=none"} {
		if !slices.Contains(office.TXT, kv) {
			t.Errorf("office's TXT has no %s: %q", kv, office.TXT)
		}
	}
	if secure.Instance != "secure @ mac" || secure.Type != "_ipps._tcp" || !slices.Contains(secure.TXT, "TLS=1.2") {
		t.Errorf("secure %+v", secure)
	}
	// with submit tokens, the clients ask for a name and a password
	e.Opts.SubmitTokens = map[string]string{"office": "secret"}
	if ads := dipp.Advertisements(e.M(), nodes, "mac"); !slices.Contains(ads[0].TXT, "air=username,password") {
		t.Errorf("TXT with submit tokens %q", ads[0].TXT)
	}
}
