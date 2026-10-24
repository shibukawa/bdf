package ipp_test

import (
	"bytes"
	"context"
	"errors"
	"io"
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
	// an IPP Everywhere and Mopria printer on its own port, with the
	// raster of Mopria and the operations of two requests
	if a := all.Get("document-format-supported"); a == nil || !slices.Contains(a.Strings(), "application/PCLm") || !slices.Contains(a.Strings(), "image/pwg-raster") {
		t.Errorf("document-format-supported %v", a)
	}
	if text("ipp-features-supported") != "ipp-everywhere" || text("mopria-certified") != "2.0" {
		t.Errorf("ipp-features-supported %q, mopria-certified %q", text("ipp-features-supported"), text("mopria-certified"))
	}
	if a := all.Get("operations-supported"); a == nil || !slices.Contains(a.Ints(), int32(ipp.OpCreateJob)) || !slices.Contains(a.Ints(), int32(ipp.OpSendDocument)) ||
		!slices.Contains(a.Ints(), int32(ipp.OpCancelMyJobs)) {
		t.Errorf("operations-supported %v", a)
	}
	if b, ok := all.Get("page-ranges-supported").Bool(); !ok || b {
		t.Errorf("page-ranges-supported %v", all.Get("page-ranges-supported"))
	}
	// its icons, on its port
	icons := all.Get("printer-icons").Strings()
	if len(icons) != 3 || icons[1] != "http://"+host+"/icons/128.png" {
		t.Fatalf("printer-icons %q", icons)
	}
	res, err := http.Get(icons[1])
	if err != nil {
		t.Fatal(err)
	}
	png, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusOK || res.Header.Get("Content-Type") != "image/png" || !bytes.HasPrefix(png, []byte("\x89PNG")) {
		t.Errorf("icon: %d %s", res.StatusCode, res.Header.Get("Content-Type"))
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

	// at the manager's own server, a request is all there is to a job:
	// no Create-Job, so no IPP Everywhere, and no icons
	own := &ipp.Client{URL: strings.Replace(e.Srv.URL, "http://", "ipp://", 1) + dipp.Path + "/office"}
	res2, err := own.Do(ctx, own.NewRequest(ipp.OpGetPrinterAttributes), nil)
	if err != nil {
		t.Fatal(err)
	}
	g := res2.Group(ipp.TagPrinter)
	if a := g.Get("operations-supported"); a == nil || slices.Contains(a.Ints(), int32(ipp.OpCreateJob)) || !slices.Contains(a.Ints(), int32(ipp.OpPrintJob)) {
		t.Errorf("operations-supported at the manager's server %v", a)
	}
	if g.Get("ipp-features-supported").Text() != "none" || g.Get("printer-icons") != nil || g.Get("mopria-certified").Text() != "2.0" {
		t.Errorf("at the manager's server: %v %v", g.Get("ipp-features-supported"), g.Get("printer-icons"))
	}
	if _, err := own.Do(ctx, own.NewRequest(ipp.OpCreateJob), nil); !errors.As(err, &se) || se.Status != ipp.StatusServerErrorOperationNotSupported {
		t.Errorf("Create-Job at the manager's server: %v", err)
	}
}

// The CUPS server nodes that say so are advertised with DNS-SD as
// AirPrint, IPP Everywhere and Mopria printers: _ipp._tcp, or _ipps._tcp
// with a certificate, with the subtypes of IPP Everywhere and AirPrint
// and the TXT keys their clients read (Mopria's among them); a node that
// says not to (dnssd false) is not.
func TestAdvertisements(t *testing.T) {
	e := flowtest.New(t, nil)
	e.Put(docflow.Node{Name: "keep", Kind: docflow.KindSink},
		docflow.Node{Name: "office", Kind: dipp.KindCups, Title: "Office 2.F", To: []string{"keep"}, Port: 16310},
		docflow.Node{Name: "secure", Kind: dipp.KindCups, To: []string{"keep"}, Port: 16311, TLSCert: "c.crt", TLSKey: "c.key"},
		docflow.Node{Name: "quiet", Kind: dipp.KindCups, To: []string{"keep"}, Port: 16312, DNSSD: new(false)})
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
	for _, kv := range []string{"rp=printers/office", "pdl=application/pdf,image/urf,image/pwg-raster,application/PCLm,image/jpeg,image/png",
		"URF=V1.4,W8,SRGB24,RS150-300-600,CP1,IS1,MT1,PQ4", "ty=BDF Docflow", "note=Office 2.F", "kind=document", "air=none", "mopria-certified=2.0", "print_wfds=T"} {
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
