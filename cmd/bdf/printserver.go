package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/server/client"
)

// serverPrint is what bdf print -server does with its flags.
type serverPrint struct {
	queue     string
	in        string
	ticket    print.Ticket
	options   client.JobOptions
	selection converter.Pages
	wait      time.Duration
	allow     bool
	quiet     bool
}

// printToServer submits a document to a print server's queue, as it is:
// the server converts what is not a .bdf. With a wait it follows the job
// until it is done.
func printToServer(ctx context.Context, sf *serverFlags, userName string, sp serverPrint) {
	say := func(format string, args ...any) {
		if !sp.quiet {
			fmt.Fprintf(os.Stderr, format+"\n", args...)
		}
	}
	c := sf.client(userName)
	password := os.Getenv(passwordEnv)
	t := sp.ticket
	// the spool holds and orders the job; the printer is not asked to
	sub := client.Submission{Queue: sp.queue, Name: t.JobName, Priority: t.Priority, HoldUntil: t.HoldUntil, Options: &sp.options}
	t.Priority, t.HoldUntil = 0, ""

	var doc io.Reader
	if isStored(sp.in) {
		r, err := openStored(sp.in)
		check(err)
		defer r.Close()
		if r.Encrypted() {
			if password == "" {
				say("%s is encrypted: the job waits for its password (bdf jobs password, or set $%s)", sp.in, passwordEnv)
			} else {
				check(unlock(r, password, sp.in))
				sub.Password = password
			}
		}
		if sp.selection != nil {
			if r.Locked() {
				check(fmt.Errorf("-pages of an encrypted document needs $%s", passwordEnv))
			}
			d, err := r.ToDocument()
			check(err)
			v, err := print.PrintView(d, sp.options.View)
			check(err)
			t.Pages = sp.selection.Numbers(len(v.Pages))
		}
		if st, err := os.Stat(sp.in); err == nil && st.IsDir() {
			// the split form goes as one file
			pr, pw := io.Pipe()
			go func() { pw.CloseWithError(r.WriteSingle(pw)) }()
			doc = pr
		} else {
			f, err := os.Open(sp.in)
			check(err)
			defer f.Close()
			doc = f
		}
	} else {
		if sp.selection != nil {
			last := 0
			for _, s := range sp.selection {
				if s.To == 0 {
					badUsage("print", "-pages with an open end (5-) is for .bdf documents with -server: the server converts the input")
				}
				last = max(last, s.To)
			}
			t.Pages = sp.selection.Numbers(last)
		}
		f, err := os.Open(sp.in)
		check(err)
		defer f.Close()
		st, err := f.Stat()
		check(err)
		// an input that needs a password takes the one set, which is
		// checked here
		switch protected, err := converter.CheckPassword(f, st.Size(), password); {
		case errors.Is(err, converter.ErrPasswordRequired):
			say("%s is encrypted: the job waits for its password (bdf jobs password, or set $%s)", sp.in, passwordEnv)
		case errors.Is(err, converter.ErrWrongPassword):
			check(fmt.Errorf("the password does not open %s", sp.in))
		case err != nil:
			check(err)
		case protected:
			sub.Password = password
		}
		doc = f
	}
	if sub.Password != "" && plaintext(c.URL) && !sp.allow {
		check(fmt.Errorf("the password of %s would cross the network without TLS: use an https:// server, or -allow-plaintext", sp.in))
	}
	if sub.Name == "" {
		sub.Name = filepath.Base(strings.TrimRight(sp.in, "/"+string(filepath.Separator)))
	}
	sub.Ticket, sub.Document = t, doc

	j, err := c.Submit(ctx, sub)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bdf:", err)
		os.Exit(serverExit(err))
	}
	fmt.Println(j.ID)
	say("job %d %s in queue %s", j.ID, describe(j), j.Queue)
	if sp.wait <= 0 {
		return
	}
	wctx, cancel := context.WithTimeout(ctx, sp.wait)
	defer cancel()
	first := true
	j, err = c.Wait(wctx, j.ID, func(j *client.Job) {
		if !first {
			say("job %d %s", j.ID, describe(j))
		}
		first = false
	})
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		fmt.Fprintf(os.Stderr, "bdf: job %d is not done within %v; it stays queued (bdf jobs show %d)\n", j.ID, sp.wait, j.ID)
		os.Exit(exitTimeout)
	case ctx.Err() != nil:
		fmt.Fprintf(os.Stderr, "bdf: stopped waiting; job %d stays queued (bdf jobs cancel %d)\n", j.ID, j.ID)
		os.Exit(exitFailed)
	case err != nil:
		fmt.Fprintln(os.Stderr, "bdf:", err)
		os.Exit(serverExit(err))
	}
	if !sp.quiet {
		for _, u := range j.Unsupported {
			fmt.Fprintln(os.Stderr, "warning: the printer does not support", u+": left out")
		}
		for _, w := range j.Warnings {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
	if j.State != "completed" {
		fmt.Fprintf(os.Stderr, "bdf: job %d %s\n", j.ID, describe(j))
		os.Exit(exitFailed)
	}
	what := print.FormatTitle(j.Format)
	say("job %d completed on %s: %d page(s) as %s on %s", j.ID, printerName(j), j.Sheets, what, firstNonEmpty(j.Media, "the printer's default paper"))
}

// describe tells a job's state, why, and where.
func describe(j *client.Job) string {
	s := j.State
	if r := strings.Join(j.Reasons, ", "); r != "" && r != "none" {
		s += " (" + r + ")"
	}
	if j.Printer != nil && !j.Done() {
		s += " at " + j.Printer.Name
	}
	if j.Message != "" {
		s += ": " + j.Message
	}
	return s
}

func printerName(j *client.Job) string {
	if j.Printer == nil {
		return "no printer"
	}
	return j.Printer.Name
}

// onlyQueue is the queue of a server that has one.
func onlyQueue(ctx context.Context, sf *serverFlags) string {
	qs, err := sf.client("").Queues(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bdf:", err)
		os.Exit(serverExit(err))
	}
	if len(qs) == 1 {
		return qs[0].Name
	}
	names := make([]string, len(qs))
	for i, q := range qs {
		names[i] = q.Name
	}
	badUsage("print", "the server has the queues "+strings.Join(names, ", ")+": name one with -queue")
	return ""
}
