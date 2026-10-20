package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/shibukawa/bdf/print/server/client"
)

const jobsUsage = `usage: bdf jobs [flags]                       the jobs not done (-all: with those done)
       bdf jobs [flags] show <id>             a job
       bdf jobs [flags] cancel|hold|release|retry <id>
       bdf jobs [flags] move <id> <queue>
       bdf jobs [flags] priority <id> <1-100>
       bdf jobs [flags] password <id>         the password of an encrypted document ($BDF_PASSWORD,
                                              or -password-file)
       bdf jobs [flags] accounting [-from date] [-to date] [-csv] [-o file]
  lists and changes the jobs of a print server (bdf printd). Without the
  manage role a token sees and changes its own jobs.`

// jobsCmd lists and changes the jobs of a print server.
func jobsCmd(args []string) {
	fs := flag.NewFlagSet("jobs", flag.ExitOnError)
	sflags := addServerFlags(fs, true)
	all := fs.Bool("all", false, "list the jobs done too (those the server keeps the documents of)")
	queue := fs.String("queue", "", "list the jobs of a queue")
	userName := fs.String("user", "", "list the jobs of a user, and act as the user with a token that may (act-as) (default: the user of the process)")
	state := fs.String("state", "", "list the jobs in a state: pending, pending-held, processing, completed …")
	limit := fs.Int("limit", 100, "list at most so many jobs")
	asJSON := fs.Bool("json", false, "write JSON")
	passwordFile := fs.String("password-file", "", "file holding the password of release and password (- for the standard input; default: $"+passwordEnv+")")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, jobsUsage)
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if *sflags.url == "" {
		badUsage("jobs", "no server: -server or $"+serverEnv)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	c := sflags.client(*userName)
	fail := func(err error) {
		if err != nil {
			fmt.Fprintln(os.Stderr, "bdf:", err)
			os.Exit(serverExit(err))
		}
	}
	rest := fs.Args()
	op := "list"
	if len(rest) > 0 {
		op, rest = rest[0], rest[1:]
	}
	id := func() int32 {
		if len(rest) == 0 {
			badUsage("jobs", op+": no job id")
		}
		n, err := strconv.ParseInt(rest[0], 10, 32)
		if err != nil || n <= 0 {
			badUsage("jobs", op+": "+rest[0]+" is not a job id")
		}
		return int32(n)
	}
	args1 := func(n int) {
		if len(rest) != n {
			fs.Usage()
			os.Exit(2)
		}
	}
	password := func() string {
		pw, err := readPassword(*passwordFile)
		check(err)
		if pw != "" && plaintext(c.URL) {
			check(fmt.Errorf("the password would cross the network without TLS: use an https:// server"))
		}
		return pw
	}
	var j *client.Job
	var err error
	switch op {
	case "list":
		args1(0)
		f := client.JobFilter{Queue: *queue, User: *userName, State: *state, Limit: min(*limit, 1000)}
		if *all && f.State == "" {
			f.State = "all"
		}
		jobs := []client.Job{}
		for len(jobs) < *limit {
			page, err := c.Jobs(ctx, f)
			fail(err)
			jobs = append(jobs, page.Jobs...)
			if page.Next == 0 {
				break
			}
			f.After = page.Next
		}
		if len(jobs) > *limit {
			jobs = jobs[:*limit]
		}
		if *asJSON {
			writeJSON(os.Stdout, jobs)
			return
		}
		tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(tw, "ID\tSTATE\tREASONS\tQUEUE\tNAME\tUSER\tPRINTER\tSUBMITTED")
		for _, j := range jobs {
			p := ""
			if j.Printer != nil {
				p = j.Printer.Name
			}
			fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", j.ID, j.State, strings.Join(j.Reasons, ","), j.Queue, j.Name, j.User, p,
				j.Submitted.Local().Format("2006-01-02 15:04"))
		}
		tw.Flush()
		return
	case "show":
		args1(1)
		j, err = c.Job(ctx, id())
	case "cancel":
		args1(1)
		j, err = c.Cancel(ctx, id())
	case "hold":
		args1(1)
		j, err = c.Hold(ctx, id())
	case "release":
		args1(1)
		j, err = c.Release(ctx, id(), password())
	case "retry":
		args1(1)
		j, err = c.Retry(ctx, id())
	case "password":
		args1(1)
		n := id()
		pw := password()
		if pw == "" {
			badUsage("jobs", "password: no password ($"+passwordEnv+" or -password-file)")
		}
		j, err = c.SetPassword(ctx, n, pw)
	case "move":
		args1(2)
		j, err = c.Move(ctx, id(), rest[1])
	case "priority":
		args1(2)
		n := id()
		p, perr := strconv.Atoi(rest[1])
		if perr != nil {
			badUsage("jobs", "priority: "+rest[1]+" is not a number")
		}
		j, err = c.SetPriority(ctx, n, p)
	case "accounting":
		accounting(ctx, c, rest)
		return
	default:
		fs.Usage()
		os.Exit(2)
	}
	fail(err)
	if *asJSON {
		writeJSON(os.Stdout, j)
		return
	}
	if op != "show" {
		fmt.Printf("job %d %s\n", j.ID, describe(j))
		return
	}
	showJob(os.Stdout, j)
}

// showJob writes a job for people.
func showJob(w io.Writer, j *client.Job) {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	line := func(k string, v any) {
		if s := fmt.Sprint(v); s != "" && s != "0" {
			fmt.Fprintf(tw, "%s:\t%s\n", k, s)
		}
	}
	when := func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Local().Format("2006-01-02 15:04:05")
	}
	line("job", j.ID)
	line("uuid", j.UUID)
	line("name", j.Name)
	line("queue", j.Queue)
	line("user", j.User)
	line("state", describe(j))
	line("priority", j.Priority)
	line("hold-until", j.HoldUntil)
	doc := fmt.Sprintf("%s, %s, %d bytes", j.Document.Name, j.Document.Format, j.Document.Size)
	if j.Document.Encrypted {
		doc += ", encrypted"
	}
	if j.DocumentsRemoved {
		doc += " (removed)"
	}
	line("document", doc)
	if len(j.Ticket) > 0 && string(j.Ticket) != "{}" {
		line("ticket", string(j.Ticket))
	}
	if j.Printer != nil {
		p := j.Printer.Name
		if j.Printer.JobID != 0 {
			p += fmt.Sprintf(" (job %d)", j.Printer.JobID)
		}
		line("printer", p)
	}
	line("attempts", j.Attempts)
	if !j.Due.IsZero() && !j.Done() {
		line("next", when(j.Due))
	}
	line("submitted", when(j.Submitted))
	line("finished", when(j.Finished))
	line("sheets", j.Sheets)
	line("format", j.Format)
	line("media", j.Media)
	for _, u := range j.Unsupported {
		line("unsupported", u)
	}
	for _, w := range j.Warnings {
		line("warning", w)
	}
	tw.Flush()
}

// accounting writes the accounting of the jobs done.
func accounting(ctx context.Context, c *client.Client, args []string) {
	fs := flag.NewFlagSet("jobs accounting", flag.ExitOnError)
	from := fs.String("from", "", "the first day (2026-10-01) or time (RFC 3339) of the jobs done")
	to := fs.String("to", "", "the last day, included, or the time before which the jobs were done")
	csv := fs.Bool("csv", false, "write CSV (default: JSON Lines)")
	out := fs.String("o", "", "the file to write (default: the standard output)")
	fs.Parse(args)
	if fs.NArg() != 0 {
		fs.Usage()
		os.Exit(2)
	}
	parse := func(name, v string, end bool) time.Time {
		if v == "" {
			return time.Time{}
		}
		if t, err := time.ParseInLocation(time.DateOnly, v, time.Local); err == nil {
			if end {
				t = t.AddDate(0, 0, 1)
			}
			return t
		}
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			badUsage("jobs", "accounting: -"+name+" "+v+": a day (2026-10-01) or an RFC 3339 time")
		}
		return t
	}
	f, t := parse("from", *from, false), parse("to", *to, true)
	format := "jsonl"
	if *csv {
		format = "csv"
	}
	w := io.Writer(os.Stdout)
	if *out != "" {
		file, err := os.Create(*out)
		check(err)
		defer file.Close()
		w = file
	}
	if err := c.Accounting(ctx, f, t, format, w); err != nil {
		fmt.Fprintln(os.Stderr, "bdf:", err)
		os.Exit(serverExit(err))
	}
}
