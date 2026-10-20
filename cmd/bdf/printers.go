package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// printerPasswordEnv names the environment variable that holds the
// password of a printer that asks for one.
const printerPasswordEnv = "BDF_PRINTER_PASSWORD"

// connFlags are the flags of talking to a printer, shared by bdf print
// and bdf printers.
type connFlags struct {
	user         *string
	username     *string
	passwordFile *string
	insecure     *bool
	trustStore   *string
	lang         *string
}

func addConnFlags(fs *flag.FlagSet) *connFlags {
	return &connFlags{
		user:         fs.String("user", "", "requesting-user-name of the jobs (default: the user of the process)"),
		username:     fs.String("username", "", "user name for a printer that asks for Basic authentication"),
		passwordFile: fs.String("password-file", "", "file holding the printer's password (- for the standard input; default: $"+printerPasswordEnv+")"),
		insecure:     fs.Bool("insecure", false, "take any certificate of an ipps printer (default: the system's roots, else the fingerprint recorded on first use)"),
		trustStore:   fs.String("trust-store", "", "file of the recorded fingerprints of printers' certificates (default: bdf/printers.json in the user's config directory)"),
		lang:         fs.String("lang", "en", "attributes-natural-language of the requests"),
	}
}

func (f *connFlags) client(url string) *ipp.Client {
	c := &ipp.Client{URL: url, User: *f.user, Username: *f.username, Language: *f.lang}
	if *f.username != "" {
		if *f.passwordFile != "" {
			pw, err := readPassword(*f.passwordFile)
			check(err)
			c.Password = pw
		} else {
			c.Password = os.Getenv(printerPasswordEnv)
		}
	}
	if *f.insecure {
		c.TLS = ipp.TLSPolicy{Insecure: true}
	} else {
		c.TLS = ipp.TLSPolicy{Store: &ipp.TrustStore{Path: *f.trustStore}}
	}
	return c
}

// printersCmd lists the printers on the network, the queues of a CUPS
// server, or what one printer supports.
func printersCmd(args []string) {
	fs := flag.NewFlagSet("printers", flag.ExitOnError)
	url := fs.String("url", "", "show what the printer at this URL supports (ipp://host/ipp/print, ipps://…)")
	cups := fs.String("cups", "", "list the queues of a CUPS server at this URL, e.g. ipp://localhost:631/ (\"local\" for that one)")
	timeout := fs.Duration("timeout", 2*time.Second, "how long to listen for printers on the network")
	asJSON := fs.Bool("json", false, "write JSON")
	sflags := addServerFlags(fs, false)
	conn := addConnFlags(fs)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: bdf printers [flags]
  lists the IPP printers that DNS-SD finds on the local network (mDNS),
  or with -url what one printer supports, or with -cups the queues of a CUPS server,
  or with -server the queues and printers of a print server (bdf printd)`)
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 0 || (*url != "" && *cups != "") || (*sflags.url != "" && (*url != "" || *cups != "")) {
		fs.Usage()
		os.Exit(2)
	}
	ctx := context.Background()
	switch {
	case *sflags.url != "":
		serverPrinters(ctx, sflags, *asJSON)
	case *url != "":
		p, err := conn.client(*url).GetPrinterAttributes(ctx)
		check(err)
		if *asJSON {
			writeJSON(os.Stdout, p)
			return
		}
		for _, l := range p.Summary() {
			fmt.Println(l)
		}
	case *cups != "":
		u := *cups
		if u == "local" {
			u = "ipp://localhost:631/"
		}
		qs, err := conn.client(u).CupsGetPrinters(ctx)
		check(err)
		if *asJSON {
			writeJSON(os.Stdout, qs)
			return
		}
		for _, q := range qs {
			state := q.State.String()
			if !q.Accepting {
				state += ", not accepting"
			}
			fmt.Printf("%s\t%s\t%s\t%s\n", q.Name, q.URI, state, firstNonEmpty(q.MakeModel, q.Info))
		}
	default:
		svcs, err := ipp.Browse(ctx, *timeout)
		check(err)
		if *asJSON {
			writeJSON(os.Stdout, svcs)
			return
		}
		if len(svcs) == 0 {
			fmt.Fprintln(os.Stderr, "no IPP printers answered on the local network; give a printer's URL with -url")
			return
		}
		for _, s := range svcs {
			fmt.Printf("%s\t%s\t%s\t%s\n", s.Name, s.URL, firstNonEmpty(s.TXT["ty"], s.TXT["product"]), strings.Join(s.Formats(), ","))
		}
	}
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func writeJSON(w io.Writer, v any) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	check(enc.Encode(v))
}

// serverPrinters lists the queues and printers of a print server.
func serverPrinters(ctx context.Context, sf *serverFlags, asJSON bool) {
	c := sf.client("")
	fail := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, "bdf:", err)
			os.Exit(serverExit(err))
		}
	}
	qs, err := c.Queues(ctx)
	fail(err)
	ps, err := c.Printers(ctx)
	fail(err)
	if asJSON {
		writeJSON(os.Stdout, map[string]any{"queues": qs, "printers": ps})
		return
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "QUEUE\tSTATE\tPOLICY\tPRINTERS\tPENDING\tHELD\tPROCESSING")
	for _, q := range qs {
		var state []string
		if q.Paused {
			state = append(state, "paused")
		}
		if q.Rejecting {
			state = append(state, "rejecting")
		}
		if state == nil {
			state = []string{"running"}
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\t%d\t%d\n", q.Name, strings.Join(state, ","), q.Policy, strings.Join(q.Printers, ","),
			q.Jobs["pending"], q.Jobs["pending-held"], q.Jobs["processing"]+q.Jobs["processing-stopped"])
	}
	fmt.Fprintln(tw)
	fmt.Fprintln(tw, "PRINTER\tSTATE\tJOBS\tMODEL\tURL")
	for _, p := range ps {
		state := p.State
		switch {
		case !p.Up:
			state = "down"
			if p.Error != "" {
				state += ": " + p.Error
			}
		case !p.Accepting:
			state += ", not accepting"
		}
		if r := strings.Join(p.Reasons, ","); r != "" && r != "none" {
			state += " (" + r + ")"
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\n", p.Name, state, p.Jobs, p.MakeModel, p.URL)
	}
	tw.Flush()
}
