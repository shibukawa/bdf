package ipp

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/print/ipp"
)

// A CUPS server (sink) is an exit: a queue of a CUPS server, a printer or
// a class there, that the job is printed to over IPP as a computer prints
// to it. The CUPS server's filters and drivers make what its printer
// takes, of the PDF (or PWG Raster) the manager sends, as print.Print
// chooses by what the queue says it supports. Its URL is the queue's
// (ipp://host:631/printers/<queue>); its credential, the user and
// password the server asks for.

// checkCupsOut takes the queue's URL and the credential of its server
// (user:password; none: no authentication).
func checkCupsOut(ctx context.Context, m *docflow.Manager, in docflow.Node, n *docflow.Node) error {
	n.URL, n.Credential = strings.TrimSpace(in.URL), strings.TrimSpace(in.Credential)
	if _, err := queueOf(n.URL); err != nil {
		return fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	return docflow.CheckCredential(n.Credential)
}

// queueOf is the queue a URL names on a CUPS server: /printers/<queue> or
// /classes/<class>, over ipp or ipps.
func queueOf(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err == nil && (u.Scheme == "ipp" || u.Scheme == "ipps") && u.Host != "" {
		for _, prefix := range []string{"/printers/", "/classes/"} {
			if q, ok := strings.CutPrefix(u.Path, prefix); ok && q != "" && !strings.Contains(q, "/") {
				return q, nil
			}
		}
	}
	return "", fmt.Errorf("url %q: a queue of a CUPS server, ipp://host:631/printers/<queue> (ipps://, /classes/<class>)", raw)
}

// cupsClient is a client of a node's CUPS server, at url, with the user
// and password of its credential.
func cupsClient(m *docflow.Manager, rawURL, credential string) (*ipp.Client, error) {
	c := &ipp.Client{URL: rawURL, User: "bdf-docflow"}
	if credential == "" {
		return c, nil
	}
	v, err := m.Credential(credential)
	if err != nil {
		return nil, err
	}
	var ok bool
	if c.Username, c.Password, ok = strings.Cut(v, ":"); !ok {
		return nil, docflow.Final(fmt.Errorf("credential %s: user:password", credential))
	}
	return c, nil
}

// sendCupsOut prints a document to the node's queue, as send prints to a
// printer.
func sendCupsOut(ctx context.Context, o *docflow.Out) (string, error) {
	c, err := cupsClient(o.Manager(), o.Node.URL, o.Node.Credential)
	if err != nil {
		return "", err
	}
	return printTo(ctx, o, c)
}

// Queue is a queue of a CUPS server, as discoverQueues lists it.
type Queue struct {
	Name string `json:"name"`
	// Path is where the queue is on its server: /printers/<name>, or
	// /classes/<name> of a class.
	Path      string `json:"path"`
	Info      string `json:"info,omitempty"`
	Location  string `json:"location,omitempty"`
	MakeModel string `json:"make_model,omitempty"`
	// State is idle, processing or stopped.
	State     string `json:"state"`
	Accepting bool   `json:"accepting"`
	Class     bool   `json:"class,omitempty"`
}

// discoverQueues lists the queues of the CUPS server of a node's URL (its
// scheme, host and port: the path may be left out), asked with the
// node's credential.
func discoverQueues(ctx context.Context, m *docflow.Manager, in docflow.Node) (any, error) {
	u, err := url.Parse(strings.TrimSpace(in.URL))
	if err != nil || (u.Scheme != "ipp" && u.Scheme != "ipps") || u.Host == "" {
		return nil, fmt.Errorf("%w: url %q: a CUPS server, ipp://host:631 (or ipps://)", docflow.ErrInvalid, in.URL)
	}
	cred := strings.TrimSpace(in.Credential)
	if err := docflow.CheckCredential(cred); err != nil {
		return nil, err
	}
	c, err := cupsClient(m, u.Scheme+"://"+u.Host+"/", cred)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", docflow.ErrInvalid, err)
	}
	ps, err := c.CupsGetPrinters(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %v", docflow.ErrPrint, c.URL, err)
	}
	out := []Queue{}
	for _, p := range ps {
		q := Queue{Name: p.Name, Path: "/printers/" + p.Name, Info: p.Info, Location: p.Location, MakeModel: p.MakeModel, State: p.State.String(), Accepting: p.Accepting, Class: p.IsClass}
		if p.IsClass {
			q.Path = "/classes/" + p.Name
		}
		out = append(out, q)
	}
	return out, nil
}
