// Command fakeprinter runs a farm of fake printers (print/fakeprinter)
// from a YAML file: printers of IPP, ipps, the raw port 9100, LPD and
// serial lines (pseudo-terminals), and CUPS servers with classes, each of
// a preset kind, failing as their faults say. A page shows what they
// received, drawn back (print/fakeprinter/readback), and a JSON API lets
// tests and scripts list the jobs and inject faults.
//
//	go run ./tools/fakeprinter -config test/print/farm.yaml
//	open http://127.0.0.1:8650/
//
// The printers listen on the loopback and are not advertised unless the
// file says so (dns-sd: true, with a host name of the farm's own).
package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image/png"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"go.yaml.in/yaml/v3"

	"github.com/shibukawa/bdf/print/fakeprinter"
	"github.com/shibukawa/bdf/print/fakeprinter/readback"
	"github.com/shibukawa/bdf/print/ipp"
)

//go:embed web
var web embed.FS

// config is the farm's YAML file.
type config struct {
	// Listen is the address the printers listen at: 127.0.0.1 by default.
	Listen string `yaml:"listen"`
	// Web is the address of the page and the API: 127.0.0.1:8650.
	Web string `yaml:"web"`
	// Store keeps the documents received in a directory of each printer;
	// without it they are kept in memory only.
	Store string `yaml:"store"`
	// DNSSD advertises the IPP printers on the local links, from Host
	// (<machine>-fake.local by default).
	DNSSD bool   `yaml:"dns-sd"`
	Host  string `yaml:"host"`
	// Fonts are the directories of the fonts of the PDF drawn back.
	Fonts    []string        `yaml:"fonts"`
	Printers []printerConfig `yaml:"printers"`
}

// printerConfig is a printer of the farm.
type printerConfig struct {
	Name   string `yaml:"name"`
	Preset string `yaml:"preset"`
	// Transport overrides the preset's: ipp, ipps, socket, lpd, serial,
	// or cups (a CUPS server of the preset cups-server).
	Transport string `yaml:"transport"`
	// Port is the port (0: any free one); Queue the LPD queue's name.
	Port       int           `yaml:"port"`
	Queue      string        `yaml:"queue"`
	Info       string        `yaml:"info"`
	Location   string        `yaml:"location"`
	Username   string        `yaml:"username"`
	Password   string        `yaml:"password"`
	Processing time.Duration `yaml:"processing"`
	// Formats replace the preset's document-format-supported.
	Formats []string            `yaml:"formats"`
	Faults  []fakeprinter.Fault `yaml:"faults"`
}

// farmPrinter is a printer running.
type farmPrinter struct {
	Name      string `json:"name"`
	Preset    string `json:"preset"`
	Transport string `json:"transport"`
	URL       string `json:"url"`
	p         *fakeprinter.Printer
	cups      *fakeprinter.CUPS
	close     func()
}

type farm struct {
	cfg      config
	printers []*farmPrinter
	mu       sync.Mutex
	cache    map[string]*readback.Result
}

func main() {
	path := flag.String("config", "test/print/farm.yaml", "the farm's YAML file")
	list := flag.Bool("presets", false, "list the presets and exit")
	flag.Parse()
	if *list {
		for _, p := range fakeprinter.Presets() {
			fmt.Printf("%s\t%s\t%s\n", p.Name, p.Transport, p.Title)
		}
		fmt.Printf("cups-server\tcups\tCUPS server: office (office-mfp), laser (pwg-only), ps (ps-office), class floor\n")
		return
	}
	b, err := os.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}
	var cfg config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatalf("%s: %v", *path, err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	f, err := start(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer f.close()
	for _, p := range f.printers {
		fmt.Printf("%-16s %-18s %-6s %s\n", p.Name, p.Preset, p.Transport, p.URL)
	}
	hs := &http.Server{Addr: cfg.Web, Handler: f.handler(), ReadHeaderTimeout: 10 * time.Second}
	go func() {
		<-ctx.Done()
		hs.Close()
	}()
	fmt.Printf("the farm's page: http://%s/\n", cfg.Web)
	if err := hs.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// start starts the printers of a configuration.
func start(ctx context.Context, cfg config) (*farm, error) {
	if cfg.Listen == "" {
		cfg.Listen = "127.0.0.1"
	}
	if cfg.Web == "" {
		cfg.Web = "127.0.0.1:8650"
	}
	f := &farm{cfg: cfg, cache: map[string]*readback.Result{}}
	seen := map[string]bool{}
	var ads []ipp.Advertisement
	for _, pc := range cfg.Printers {
		if pc.Name == "" || seen[pc.Name] {
			f.close()
			return nil, fmt.Errorf("a printer without a name, or %q twice", pc.Name)
		}
		seen[pc.Name] = true
		fp, err := f.startPrinter(pc)
		if err != nil {
			f.close()
			return nil, fmt.Errorf("%s: %w", pc.Name, err)
		}
		f.printers = append(f.printers, fp)
		if fp.Transport == "ipp" || fp.Transport == "ipps" {
			ads = append(ads, advertisement(fp))
		}
	}
	if cfg.DNSSD && len(ads) > 0 {
		host := cfg.Host
		if host == "" {
			name, _ := os.Hostname()
			host = strings.TrimSuffix(strings.Split(name, ".")[0], ".local") + "-fake.local"
		}
		r, err := ipp.Advertise(ctx, ads, ipp.ResponderOptions{Host: host})
		if err != nil {
			f.close()
			return nil, fmt.Errorf("dns-sd: %w", err)
		}
		f.printers = append(f.printers, &farmPrinter{Name: "", close: r.Close})
	}
	return f, nil
}

// advertisement is the DNS-SD service of an IPP printer of the farm.
func advertisement(fp *farmPrinter) ipp.Advertisement {
	typ := "_ipp._tcp"
	if fp.Transport == "ipps" {
		typ = "_ipps._tcp"
	}
	_, port, _ := net.SplitHostPort(strings.TrimPrefix(strings.TrimPrefix(strings.Split(fp.URL, "/ipp/")[0], "ipps://"), "ipp://"))
	n, _ := strconv.Atoi(port)
	return ipp.Advertisement{Instance: fp.Name + " (fake)", Type: typ, Subtypes: []string{"_print"}, Port: n,
		TXT: []string{"txtvers=1", "rp=ipp/print", "ty=" + fp.p.MakeModel, "pdl=" + strings.Join(fp.p.Formats, ","), "UUID=" + strings.TrimPrefix(fp.p.UUID, "urn:uuid:")}}
}

func (f *farm) startPrinter(pc printerConfig) (*farmPrinter, error) {
	addr := net.JoinHostPort(f.cfg.Listen, strconv.Itoa(pc.Port))
	if pc.Preset == "cups-server" || pc.Transport == "cups" {
		c := fakeprinter.CUPSPreset()
		for _, p := range c.Printers() {
			f.setup(p, pc, pc.Name+"."+p.Name)
		}
		s, err := c.Serve(addr)
		if err != nil {
			return nil, err
		}
		return &farmPrinter{Name: pc.Name, Preset: "cups-server", Transport: "cups", URL: s.URL + "/", cups: c, close: func() { s.Close() }}, nil
	}
	ps, ok := fakeprinter.LookupPreset(pc.Preset)
	if !ok {
		return nil, fmt.Errorf("no preset %q (fakeprinter -presets lists them)", pc.Preset)
	}
	p := ps.New(pc.Name)
	f.setup(p, pc, pc.Name)
	transport := pc.Transport
	if transport == "" {
		transport = ps.Transport
	}
	fp := &farmPrinter{Name: pc.Name, Preset: pc.Preset, Transport: transport, p: p}
	switch transport {
	case "ipp", "ipps":
		serve := p.ServeIPP
		if transport == "ipps" {
			serve = p.ServeIPPS
		}
		s, err := serve(addr)
		if err != nil {
			return nil, err
		}
		fp.URL, fp.close = s.URL, func() { s.Close() }
	case "socket", "lpd":
		var l *fakeprinter.Listener
		var err error
		if transport == "socket" {
			l, err = p.ListenSocket(addr)
		} else {
			q := pc.Queue
			if q == "" {
				q = pc.Name
			}
			l, err = p.ListenLPD(addr, q)
		}
		if err != nil {
			return nil, err
		}
		fp.URL, fp.close = l.URL+profileQuery(p), l.Close
	case "serial":
		far, name, err := fakeprinter.OpenPTY()
		if err != nil {
			return nil, err
		}
		go p.ServeSerial(far, 300*time.Millisecond)
		fp.URL, fp.close = "serial://"+name+profileQuery(p), func() { far.Close() }
	default:
		return nil, fmt.Errorf("no transport %q: ipp, ipps, socket, lpd, serial or cups", transport)
	}
	return fp, nil
}

// profileQuery is what a raw printer's URL says of its language: the
// profile of print/pdl it is of.
func profileQuery(p *fakeprinter.Printer) string {
	if p.Profile != "" {
		return "?profile=" + p.Profile
	}
	return ""
}

// setup sets what a printer's configuration says.
func (f *farm) setup(p *fakeprinter.Printer, pc printerConfig, dir string) {
	if pc.Info != "" {
		p.Info = pc.Info
	}
	if pc.Location != "" {
		p.Location = pc.Location
	}
	p.Username, p.Password = pc.Username, pc.Password
	if pc.Processing > 0 {
		p.Processing = pc.Processing
	}
	if len(pc.Formats) > 0 {
		p.Formats = pc.Formats
	}
	if f.cfg.Store != "" {
		p.Dir = filepath.Join(f.cfg.Store, filepath.FromSlash(dir))
	}
	for _, fault := range pc.Faults {
		p.Inject(fault)
	}
}

func (f *farm) close() {
	for _, p := range f.printers {
		if p.close != nil {
			p.close()
		}
	}
}

// all are the printers of the farm, those of its CUPS servers each as
// "<server>.<printer>".
func (f *farm) all() []*farmPrinter {
	var out []*farmPrinter
	for _, fp := range f.printers {
		switch {
		case fp.cups != nil:
			for _, p := range fp.cups.Printers() {
				out = append(out, &farmPrinter{Name: fp.Name + "." + p.Name, Preset: "cups-server", Transport: "cups", URL: fp.URL + "printers/" + p.Name, p: p})
			}
		case fp.p != nil:
			out = append(out, fp)
		}
	}
	return out
}

func (f *farm) printer(name string) *farmPrinter {
	for _, fp := range f.all() {
		if fp.Name == name {
			return fp
		}
	}
	return nil
}

type printerJSON struct {
	*farmPrinter
	MakeModel string   `json:"make-and-model"`
	State     string   `json:"state"`
	Reasons   []string `json:"reasons"`
	Accepting bool     `json:"accepting"`
	Jobs      int      `json:"jobs"`
}

type jobJSON struct {
	Printer   string    `json:"printer"`
	ID        int32     `json:"id"`
	Transport string    `json:"transport"`
	Name      string    `json:"name"`
	User      string    `json:"user"`
	Format    string    `json:"format"`
	Bytes     int       `json:"bytes"`
	Created   time.Time `json:"created"`
	State     string    `json:"state"`
	Attrs     []string  `json:"attributes,omitempty"`
	Queue     string    `json:"queue,omitempty"`
	Control   string    `json:"control,omitempty"`
}

func (f *farm) handler() http.Handler {
	mux := http.NewServeMux()
	sub, _ := fs.Sub(web, "web")
	mux.Handle("GET /", http.FileServerFS(sub))
	mux.HandleFunc("GET /api/printers", func(w http.ResponseWriter, r *http.Request) {
		var out []printerJSON
		for _, fp := range f.all() {
			st := fp.p.Snapshot()
			out = append(out, printerJSON{farmPrinter: fp, MakeModel: fp.p.MakeModel, State: st.State.String(), Reasons: st.Reasons, Accepting: st.Accepting, Jobs: len(fp.p.Jobs())})
		}
		writeJSON(w, out)
	})
	mux.HandleFunc("GET /api/jobs", func(w http.ResponseWriter, r *http.Request) {
		want := r.URL.Query().Get("printer")
		out := []jobJSON{}
		for _, fp := range f.all() {
			if want != "" && fp.Name != want {
				continue
			}
			for _, j := range fp.p.Jobs() {
				jj := jobJSON{Printer: fp.Name, ID: j.ID, Transport: j.Transport, Name: j.Name, User: j.User, Format: j.Format,
					Bytes: len(j.Data()), Created: j.Created, State: fp.p.JobState(j).String(), Queue: j.Queue, Control: j.Control}
				for _, a := range j.Attrs {
					jj.Attrs = append(jj.Attrs, a.Name+"="+a.String())
				}
				out = append(out, jj)
			}
		}
		slices.SortFunc(out, func(a, b jobJSON) int { return b.Created.Compare(a.Created) })
		writeJSON(w, out)
	})
	job := func(w http.ResponseWriter, r *http.Request) (*farmPrinter, *fakeprinter.Job) {
		fp := f.printer(r.PathValue("printer"))
		if fp == nil {
			http.Error(w, "no such printer", http.StatusNotFound)
			return nil, nil
		}
		id, _ := strconv.Atoi(r.PathValue("id"))
		j := fp.p.Job(int32(id))
		if j == nil {
			http.Error(w, "no such job", http.StatusNotFound)
			return nil, nil
		}
		return fp, j
	}
	mux.HandleFunc("GET /api/printers/{printer}/jobs/{id}/data", func(w http.ResponseWriter, r *http.Request) {
		if _, j := job(w, r); j != nil {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="job-%d.prn"`, j.ID))
			w.Write(j.Data())
		}
	})
	mux.HandleFunc("GET /api/printers/{printer}/jobs/{id}/pages", func(w http.ResponseWriter, r *http.Request) {
		fp, j := job(w, r)
		if j == nil {
			return
		}
		res, err := f.readBack(fp, j)
		if err != nil {
			writeJSON(w, map[string]any{"error": err.Error()})
			return
		}
		type page struct {
			Width, Height, DPI int
			Symbols            []string `json:",omitempty"`
			Cut                string   `json:",omitempty"`
		}
		var pages []page
		for _, pg := range res.Pages {
			b := pg.Image.Bounds()
			pages = append(pages, page{b.Dx(), b.Dy(), pg.DPI, pg.Symbols, pg.Cut})
		}
		writeJSON(w, map[string]any{"format": res.Format, "pages": pages, "commands": res.Commands, "drawer": res.Drawer})
	})
	mux.HandleFunc("GET /api/printers/{printer}/jobs/{id}/pages/{n}", func(w http.ResponseWriter, r *http.Request) {
		fp, j := job(w, r)
		if j == nil {
			return
		}
		res, err := f.readBack(fp, j)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		n, _ := strconv.Atoi(strings.TrimSuffix(r.PathValue("n"), ".png"))
		if n < 1 || n > len(res.Pages) {
			http.Error(w, "no such page", http.StatusNotFound)
			return
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, res.Pages[n-1].Image); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(buf.Bytes())
	})
	mux.HandleFunc("POST /api/printers/{printer}/faults", func(w http.ResponseWriter, r *http.Request) {
		fp := f.printer(r.PathValue("printer"))
		if fp == nil {
			http.Error(w, "no such printer", http.StatusNotFound)
			return
		}
		var fault fakeprinter.Fault
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&fault); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fp.p.Inject(fault)
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /api/reset", func(w http.ResponseWriter, r *http.Request) {
		for _, fp := range f.all() {
			fp.p.Inject(fakeprinter.Fault{Restart: true, Resume: true, Accept: true})
		}
		f.mu.Lock()
		f.cache = map[string]*readback.Result{}
		f.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}

// readBack reads a job back, once.
func (f *farm) readBack(fp *farmPrinter, j *fakeprinter.Job) (*readback.Result, error) {
	key := fmt.Sprintf("%s\x00%s\x00%d", fp.Name, j.UUID, j.ID)
	f.mu.Lock()
	res := f.cache[key]
	f.mu.Unlock()
	if res != nil {
		return res, nil
	}
	if fp.p.JobState(j) == ipp.JobPending {
		return nil, errors.New("the job has not all its documents")
	}
	res, err := readback.Job(fp.p, j, &readback.Options{FontDirs: f.cfg.Fonts})
	if err != nil {
		return nil, err
	}
	f.mu.Lock()
	f.cache[key] = res
	f.mu.Unlock()
	return res, nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
