package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
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
	conn := addConnFlags(fs)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: bdf printers [flags]
  lists the IPP printers that DNS-SD finds on the local network (mDNS),
  or with -url what one printer supports, or with -cups the queues of a CUPS server`)
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 0 || (*url != "" && *cups != "") {
		fs.Usage()
		os.Exit(2)
	}
	ctx := context.Background()
	switch {
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
