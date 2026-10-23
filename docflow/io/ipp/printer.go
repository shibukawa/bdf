package ipp

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raw"
)

// Routes of a printer: how its jobs reach it.
const (
	// RouteDirect: the manager sends the job itself (the printer is on a
	// network the manager reaches: the same LAN, or a VPC with Direct
	// Connect or a VPN).
	RouteDirect = "direct"
	// RouteBranch: a bdf printd of the branch fetches the job (INFRA).
	// Not done yet.
	RouteBranch = "branch"
)

// ErrBranch is the error of a printer that a branch's bdf printd is to
// fetch the jobs of: not done yet.
var ErrBranch = fmt.Errorf("%w: printers of a branch (route branch) are not supported yet", docflow.ErrUnsupported)

func checkPrinter(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	n.URL, n.Route = strings.TrimSpace(in.URL), or(in.Route, RouteDirect)
	if n.Route != RouteDirect && n.Route != RouteBranch {
		return fmt.Errorf("%w: route %q: direct or branch", docflow.ErrInvalid, in.Route)
	}
	if err := checkPrinterURL(n.URL); err != nil {
		return err
	}
	n.Preview = strings.TrimSpace(in.Preview)
	return docflow.CheckPreviewURL(n.Preview)
}

// checkPrinterURL takes the printers a manager reaches over the network:
// IPP, and 9100 and LPD. Serial, USB and Windows queues are a branch
// printd's.
func checkPrinterURL(u string) error {
	if raw.IsURL(u) {
		p, err := raw.Parse(u)
		if err != nil {
			return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
		}
		if p.Scheme != "socket" && p.Scheme != "lpd" {
			return fmt.Errorf("%w: %s: a printer of a manager is reached over the network (ipp, ipps, socket, lpd)", docflow.ErrInvalid, u)
		}
		return nil
	}
	if _, _, err := ipp.ParseURL(u); err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	return nil
}

// ready refuses a printer a branch fetches from.
func ready(n *docflow.Node) error {
	if n.Route == RouteBranch {
		return ErrBranch
	}
	return nil
}

// send prints a document on a printer, with what the job's metadata says
// for the printer (printMeta), and returns the printer's job id (none for
// a raw printer).
func send(ctx context.Context, o *docflow.Out) (string, error) {
	if raw.IsURL(o.Node.URL) {
		rp, err := raw.Parse(o.Node.URL)
		if err != nil {
			return "", err
		}
		doc, opts, err := printJob(ctx, o)
		if err != nil {
			return "", err
		}
		_, err = print.PrintRaw(ctx, rp, doc, opts)
		return "", err
	}
	return printTo(ctx, o, &ipp.Client{URL: o.Node.URL, User: "bdf-docflow"})
}

// printJob is the document of a send and how it is printed.
func printJob(ctx context.Context, o *docflow.Out) (*bdf.Document, *print.Options, error) {
	doc, err := o.Document(ctx)
	if err != nil {
		return nil, nil, err
	}
	mo := o.Manager().Options()
	return doc, &print.Options{Ticket: print.Ticket{JobName: o.Name, Meta: printMeta(o.Meta)}, FontDirs: mo.FontDirs, NoSystemFonts: mo.NoSystemFonts}, nil
}

// printTo prints a document on an IPP printer (a printer, a CUPS server's
// queue), and returns the printer's job id.
func printTo(ctx context.Context, o *docflow.Out, c *ipp.Client) (string, error) {
	doc, opts, err := printJob(ctx, o)
	if err != nil {
		return "", err
	}
	res, err := print.Print(ctx, c, doc, opts)
	if err != nil {
		return "", err
	}
	if res.Job != nil && res.Job.ID != 0 {
		return strconv.Itoa(int(res.Job.ID)), nil
	}
	return "", nil
}

// printMeta is what a job's metadata says for the printer, as print.Meta
// carries it: the keys ipp.job-account-id, ipp.job-accounting-user-id,
// ipp.document-metadata and ipp.x-… (what came in as those IPP
// attributes, or was given as metadata under these keys) become the
// attributes again. A value that is not one the attribute takes is left
// out, as print.Meta.Take leaves it.
func printMeta(meta map[string]any) print.Meta {
	var m print.Meta
	for k, v := range meta {
		name, ok := strings.CutPrefix(k, "ipp.")
		if !ok {
			continue
		}
		var vals []ipp.Value
		switch x := v.(type) {
		case []string:
			for _, s := range x {
				vals = append(vals, ipp.Name(s))
			}
		case []any:
			for _, e := range x {
				vals = append(vals, ipp.Name(docflow.MetaText(e)))
			}
		default:
			vals = []ipp.Value{ipp.Name(docflow.MetaText(v))}
		}
		m.Take(ipp.Attr(name, vals...))
	}
	return m
}
