package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raster"
	"github.com/shibukawa/bdf/print/server/client"
)

// Exit statuses of bdf print, besides 0 (printed) and 2 (usage).
const (
	exitFailed      = 1 // the job was aborted or canceled, or another error
	exitConnection  = 3 // the printer could not be reached, or its certificate or password was refused
	exitUnsupported = 4 // the printer refused the job or, with -strict, what the ticket asks
	exitStopped     = 5 // the printer is stopped or not accepting jobs
	exitTimeout     = 6 // the job was sent and is not done within -wait
)

const printUsage = `usage: bdf print [flags] -printer <url | name> <file.bdf | dir | input>
       bdf print [flags] -server <url> [-queue <queue>] <file.bdf | dir | input>
  prints a document on an IPP printer, or a queue of CUPS, without a driver:
  as a PDF file when the printer takes PDF (the document stays vectors), else
  as PWG Raster drawn here; an input in a format bdf generate converts is
  converted first. bdf printers lists the printers and what they support.

  With -server (or $BDF_SERVER, without -printer) the document goes as it is
  to a queue of a print server (bdf printd), which converts and prints it;
  bdf print writes the job's id and returns, or follows the job with -wait.
  bdf jobs lists and changes the jobs.

exit status: 0 printed (or queued), 1 aborted, canceled or failed, 2 usage,
3 the printer or server could not be reached (or refused the certificate,
password or token), 4 the printer or server refused the job (or, with
-strict, what it asks), 5 the printer is stopped (or the queue does not
take jobs), 6 the job is not done within -wait.`

// printCmd prints a document.
func printCmd(args []string) {
	fs := flag.NewFlagSet("print", flag.ExitOnError)
	printer := fs.String("printer", "", "the printer: its URL (ipp://host/ipp/print, ipps://…, ipp://localhost:631/printers/<queue>) or its name on the network (bdf printers)")
	copies := fs.Int("copies", 1, "copies")
	collate := fs.String("collate", "", "collated (1 2 3 1 2 3) or uncollated (1 1 2 2 3 3) copies")
	sides := fs.String("sides", "", "one-sided, two-sided-long-edge (long), two-sided-short-edge (short)")
	media := fs.String("media", "", "paper: a PWG name (iso_a4_210x297mm) or a4, a3, b4, b5 (JIS), letter, legal, hagaki, chou3, kaku2 …; default: the paper of the size of the first page")
	source := fs.String("source", "", "input tray (media-source): tray-1, manual …")
	mtype := fs.String("media-type", "", "media-type: stationery, labels, envelope …")
	color := fs.String("color", "", "print-color-mode: color, monochrome or auto")
	mono := fs.Bool("mono", false, "print in black (the same as -color monochrome)")
	bilevel := fs.Bool("bilevel", false, "with -mono, a raster page in black and white without gray (line drawings and text)")
	quality := fs.String("quality", "", "draft, normal or high")
	resolution := fs.Int("resolution", 0, "dots per inch (default: 600 for black, 300 for colour, for raster pages)")
	pages := fs.String("pages", "", "pages to print, e.g. 1-3,5 (default: all)")
	scaling := fs.String("scaling", "", "how a page is fitted to the paper: auto, auto-fit (shrink, never enlarge), fit, fill or none")
	orientation := fs.String("orientation", "", "portrait, landscape, reverse-landscape or reverse-portrait (default: landscape pages turn onto portrait paper)")
	finishings := fs.String("finishings", "", "finishings, comma-separated: staple, punch, fold, or their enum numbers")
	bin := fs.String("output-bin", "", "output-bin")
	hold := fs.String("hold", "", "job-hold-until: indefinite, night, weekend …")
	priority := fs.Int("priority", 0, "job-priority, 1 to 100")
	format := fs.String("format", "auto", "what to send: auto (PDF when the printer takes it, else PostScript, PWG Raster, Apple Raster or PCLm), pdf, ps, pwg-raster, urf or pclm")
	rasterIfLossy := fs.Bool("raster-if-lossy", false, "send PWG Raster when the PDF file would leave out what the page draws (FILTER, composite operations)")
	view := fs.String("view", "", "id of the view to print (default: the first view of pages)")
	var skip stringList
	fs.Var(&skip, "skip", "role of the layers to leave out, e.g. background or master, for paper that has them printed already (repeatable)")
	strict := fs.Bool("strict", false, "fail when the printer does not support what is asked (default: leave it out, with a warning)")
	wait := fs.Duration("wait", 10*time.Minute, "how long to wait for the job to be done (0: return once it is sent)")
	maxImpressions := fs.Int("max-impressions", 10000, "the most pages to send (pages times copies)")
	jobName := fs.String("job-name", "", "job-name (default: the file's name)")
	timeout := fs.Duration("timeout", 2*time.Second, "how long to look for a printer named on the network")
	allow := fs.Bool("allow-plaintext", false, "print an encrypted document or password-protected input (read with $"+passwordEnv+") over a connection without TLS")
	var params stringList
	fs.Var(&params, "param", "option of the format of an input to convert, as name=value (repeatable; see bdf generate -h)")
	quiet := fs.Bool("q", false, "do not print progress and warnings")
	sflags := addServerFlags(fs, true)
	queue := fs.String("queue", "", "the server's queue (default: its only queue)")
	conn := addConnFlags(fs)
	ropts := fontFlags(fs)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, printUsage)
		fs.PrintDefaults()
	}
	fs.Parse(args)
	set := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { set[f.Name] = true })
	if set["server"] && *printer != "" {
		badUsage("print", "-printer prints here and -server on a print server: one of them")
	}
	toServer := *printer == "" && *sflags.url != ""
	if fs.NArg() != 1 || *printer == "" && !toServer {
		fs.Usage()
		os.Exit(2)
	}
	if !toServer && set["queue"] {
		badUsage("print", "-queue is a queue of a print server (-server)")
	}
	in := fs.Arg(0)
	t := print.Ticket{Copies: *copies, Collate: *collate, MediaSource: *source, MediaType: *mtype, ColorMode: *color, Bilevel: *bilevel,
		Resolution: *resolution, Scaling: raster.Scaling(*scaling), OutputBin: *bin, HoldUntil: *hold, Priority: *priority,
		JobName: *jobName}
	switch *sides {
	case "", "one-sided", "two-sided-long-edge", "two-sided-short-edge":
		t.Sides = *sides
	case "one":
		t.Sides = "one-sided"
	case "long", "two-sided":
		t.Sides = "two-sided-long-edge"
	case "short":
		t.Sides = "two-sided-short-edge"
	default:
		badUsage("print", "-sides must be one-sided, two-sided-long-edge or two-sided-short-edge")
	}
	if *mono {
		if *color != "" && *color != "monochrome" {
			badUsage("print", "-mono and -color "+*color)
		}
		t.ColorMode = "monochrome"
	}
	if *bilevel && t.ColorMode != "monochrome" {
		badUsage("print", "-bilevel is for -mono")
	}
	switch *quality {
	case "":
	case "draft":
		t.Quality = 3
	case "normal":
		t.Quality = 4
	case "high":
		t.Quality = 5
	default:
		badUsage("print", "-quality must be draft, normal or high")
	}
	switch *orientation {
	case "":
	case "portrait":
		t.Orientation = raster.Portrait
	case "landscape":
		t.Orientation = raster.Landscape
	case "reverse-landscape":
		t.Orientation = raster.ReverseLandscape
	case "reverse-portrait":
		t.Orientation = raster.ReversePortrait
	default:
		badUsage("print", "-orientation must be portrait, landscape, reverse-landscape or reverse-portrait")
	}
	if *finishings != "" {
		for _, f := range strings.Split(*finishings, ",") {
			n, ok := finishingNames[strings.TrimSpace(f)]
			if !ok {
				var err error
				if n, err = strconv.Atoi(strings.TrimSpace(f)); err != nil {
					badUsage("print", "-finishings: "+f+" is not a finishing")
				}
			}
			t.Finishings = append(t.Finishings, n)
		}
	}
	if *media != "" {
		m, err := print.ResolveMedia(*media)
		if err != nil {
			badUsage("print", "-media: "+err.Error())
		}
		t.Media = m.Name
	}
	var selection converter.Pages
	if *pages != "" && *pages != "all" {
		var err error
		if selection, err = converter.ParsePages(*pages); err != nil {
			badUsage("print", err.Error())
		}
	}
	if t.JobName == "" {
		t.JobName = filepath.Base(strings.TrimRight(in, "/"+string(filepath.Separator)))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if toServer {
		for _, name := range []string{"param", "max-impressions", "font-dir", "no-system-fonts", "timeout", "username", "password-file", "insecure", "trust-store"} {
			if set[name] {
				badUsage("print", "-"+name+" is not for -server: the server converts and prints with its configuration")
			}
		}
		if !set["copies"] {
			t.Copies = 0 // the queue's
		}
		sp := serverPrint{queue: *queue, in: in, ticket: t, selection: selection, allow: *allow, quiet: *quiet,
			options: client.JobOptions{View: *view, SkipRoles: skip, Strict: *strict}}
		if *format != "auto" {
			sp.options.Format = *format
		}
		if set["raster-if-lossy"] {
			sp.options.RasterIfLossy = rasterIfLossy
		}
		// a job on a server is left to it unless asked
		if set["wait"] {
			sp.wait = *wait
		}
		if sp.queue == "" {
			sp.queue = onlyQueue(ctx, sflags)
		}
		printToServer(ctx, sflags, *conn.user, sp)
		return
	}
	target := *printer
	if !strings.Contains(target, "://") {
		target = findPrinter(ctx, target, *timeout)
	}
	c := conn.client(target)
	secure := strings.HasPrefix(strings.ToLower(target), "ipps:") || strings.HasPrefix(strings.ToLower(target), "https:")

	ro := ropts()
	var doc *bdf.Document
	if isStored(in) {
		r, err := open(in)
		check(err)
		if r.Locked() {
			check(fmt.Errorf("%s is encrypted: set $%s", in, passwordEnv))
		}
		if r.Encrypted() && !secure && !*allow {
			check(fmt.Errorf("%s is encrypted, and the printer's connection has no TLS: use ipps://, or -allow-plaintext", in))
		}
		doc, err = r.ToDocument()
		r.Close()
		check(err)
		if selection != nil {
			v, err := print.PrintView(doc, *view)
			check(err)
			t.Pages = selection.Numbers(len(v.Pages))
		}
	} else {
		opts := &converter.Options{FontDirs: ro.FontDirs, NoSystemFonts: ro.NoSystemFonts, Password: os.Getenv(passwordEnv), Pages: selection,
			Params: map[string]string{}}
		if f, err := converter.DetectFile(in); err == nil && printViews[f.Name] != "" {
			opts.Params["views"] = printViews[f.Name] // views of pages
		}
		for _, kv := range params {
			name, value, ok := strings.Cut(kv, "=")
			if !ok || name == "" {
				badUsage("print", "-param "+kv+": want name=value")
			}
			opts.Params[name] = value
		}
		res, err := converter.ConvertFile(in, "", opts)
		switch {
		case errors.Is(err, converter.ErrUnknownFormat):
			badUsage("print", in+": neither a bdf document nor an input format bdf generate converts")
		case errors.Is(err, converter.ErrPasswordRequired):
			check(fmt.Errorf("%s is encrypted: set $%s", in, passwordEnv))
		case errors.Is(err, converter.ErrWrongPassword):
			check(fmt.Errorf("the password does not open %s", in))
		}
		check(err)
		if res.Protected && !secure && !*allow {
			check(fmt.Errorf("%s is password-protected, and the printer's connection has no TLS: use ipps://, or -allow-plaintext", in))
		}
		if !*quiet {
			for _, w := range res.Warnings {
				fmt.Fprintln(os.Stderr, "warning:", w)
			}
		}
		doc = res.Doc
	}

	say := func(format string, args ...any) {
		if !*quiet {
			fmt.Fprintf(os.Stderr, format+"\n", args...)
		}
	}
	w := *wait
	if w < 0 {
		w = 0
	}
	res, err := print.Print(ctx, c, doc, &print.Options{Ticket: t, Format: *format, RasterIfLossy: *rasterIfLossy, View: *view, SkipRoles: skip,
		Strict: *strict, FontDirs: ro.FontDirs, NoSystemFonts: ro.NoSystemFonts, Wait: w, MaxImpressions: *maxImpressions,
		Progress: func(e print.Event) {
			switch e.Kind {
			case "created":
				say("job %d created", e.Job.ID)
			case "sent":
				say("sent %s", size(e.Bytes))
			case "state":
				say("%s", e.Job)
			}
		}})
	if res != nil && !*quiet {
		for _, u := range res.Unsupported {
			fmt.Fprintln(os.Stderr, "warning: the printer does not support", u+": left out")
		}
		for _, w := range res.Warnings {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "bdf:", err)
		os.Exit(exitCode(err))
	}
	what := print.FormatTitle(res.Format)
	state := "sent"
	if res.Job != nil && res.Job.State.Done() {
		state = res.Job.State.String()
	}
	id := int32(0)
	if res.Job != nil {
		id = res.Job.ID
	}
	say("%s: job %d %s, %d page(s) as %s on %s", c.PrinterURI(), id, state, res.Sheets, what, firstNonEmpty(res.Media, "the printer's default paper"))
}

// finishingNames are the enums of finishings by name.
var finishingNames = map[string]int{"none": 3, "staple": 4, "punch": 5, "cover": 6, "bind": 7, "saddle-stitch": 8, "edge-stitch": 9, "fold": 10,
	"trim": 11, "bale": 12, "booklet-maker": 13, "jog-offset": 14, "coat": 15, "laminate": 16,
	"staple-top-left": 20, "staple-bottom-left": 21, "staple-top-right": 22, "staple-bottom-right": 23,
	"staple-dual-left": 28, "staple-dual-top": 29, "punch-dual-left": 74, "punch-dual-top": 75}

// exitCode is the exit status of an error of a job.
func exitCode(err error) int {
	var se *ipp.StatusError
	var netErr *net.OpError
	var urlErr *url.Error
	switch {
	case errors.Is(err, print.ErrPrinterStopped), errors.As(err, &se) && (se.Status == ipp.StatusServerErrorNotAcceptingJobs || se.Status == ipp.StatusServerErrorPrinterIsDeactivated):
		return exitStopped
	case errors.Is(err, print.ErrUnsupported), errors.As(err, &se) && se.Status.ClientError():
		return exitUnsupported
	case errors.Is(err, ipp.ErrCertificateChanged), errors.Is(err, ipp.ErrAuth), errors.Is(err, ipp.ErrProtocol), errors.As(err, &netErr), errors.As(err, &urlErr):
		return exitConnection
	case errors.Is(err, print.ErrTimeout):
		return exitTimeout
	}
	return exitFailed
}

func size(n int64) string {
	switch {
	case n >= 10<<20:
		return strconv.FormatInt(n>>20, 10) + " MiB"
	case n >= 10<<10:
		return strconv.FormatInt(n>>10, 10) + " KiB"
	}
	return strconv.FormatInt(n, 10) + " bytes"
}

// findPrinter finds the URL of a printer by its name on the network: the
// one of that name, or the only one whose name holds it.
func findPrinter(ctx context.Context, name string, timeout time.Duration) string {
	svcs, err := ipp.Browse(ctx, timeout)
	check(err)
	var near []ipp.Service
	for _, s := range svcs {
		if strings.EqualFold(s.Name, name) {
			return s.URL
		}
		if strings.Contains(strings.ToLower(s.Name), strings.ToLower(name)) {
			near = append(near, s)
		}
	}
	switch len(near) {
	case 1:
		return near[0].URL
	case 0:
		fmt.Fprintf(os.Stderr, "bdf: no printer named %q on the network (bdf printers lists them; give a URL with -printer ipp://…)\n", name)
	default:
		fmt.Fprintf(os.Stderr, "bdf: %d printers are named like %q:\n", len(near), name)
		for _, s := range near {
			fmt.Fprintf(os.Stderr, "  %s\t%s\n", s.Name, s.URL)
		}
	}
	os.Exit(exitConnection)
	return ""
}
