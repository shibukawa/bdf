// Package server is the print server: a spool (print/spool) of the
// queues and printers of a configuration file, its HTTP API (/api/v1),
// the changes as Server-Sent Events, /metrics, and a page of its state.
// bdf printd runs it.
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"

	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/pdl"
	"github.com/shibukawa/bdf/print/raw"
	"github.com/shibukawa/bdf/print/spool"
)

// Config is the configuration of a server, as its YAML file has it.
type Config struct {
	// Listen is the address of the HTTP API: 127.0.0.1:8640 by default.
	Listen string `yaml:"listen"`
	// Spool is the spool's file: spool.db next to the configuration
	// file by default.
	Spool string `yaml:"spool"`
	// User is the requesting-user-name sent to the printers: "bdf" by
	// default.
	User string `yaml:"user"`
	// TLS serves HTTPS; it is required when Listen is not a loopback
	// address, unless Insecure.
	TLS      *TLSConfig `yaml:"tls"`
	Insecure bool       `yaml:"insecure"`
	// Tokens are who may use the API. None on a loopback address lets
	// anyone of the machine use it.
	Tokens    []TokenConfig            `yaml:"tokens"`
	Printers  map[string]PrinterConfig `yaml:"printers"`
	Queues    map[string]QueueConfig   `yaml:"queues"`
	Retention RetentionConfig          `yaml:"retention"`
	Limits    LimitsConfig             `yaml:"limits"`
	Fonts     FontsConfig              `yaml:"fonts"`
	// TempDir holds documents while they arrive.
	TempDir string `yaml:"temp-dir"`
	// IPP shows the queues as IPP printers, which the systems add without
	// a driver (IPP Everywhere, AirPrint, Mopria).
	IPP *IPPConfig `yaml:"ipp"`
	// Upstreams are the servers this one fetches jobs from, as a proxy
	// of a branch (IPP INFRA): each into a queue of this server.
	Upstreams map[string]UpstreamConfig `yaml:"upstreams"`

	dir string // of the file, for relative paths
}

// UpstreamConfig is a server this one fetches jobs from: the URL of its
// queue (a queue of a printer of proxies there), the queue here the jobs
// go to, and the token's name and secret there (Basic).
type UpstreamConfig struct {
	URL          string `yaml:"url"`
	Queue        string `yaml:"queue"`
	Username     string `yaml:"username"`
	PasswordEnv  string `yaml:"password-env"`
	PasswordFile string `yaml:"password-file"`
	Fingerprint  string `yaml:"fingerprint"`
	Insecure     bool   `yaml:"insecure"`
	// Poll is how often the jobs are asked for: 10 seconds by default.
	Poll Duration `yaml:"poll"`
	// UUID is the output-device-uuid of this proxy; one made of the
	// spool and the upstream's name by default.
	UUID string `yaml:"uuid"`
}

// IPPConfig is the IPP server of the queues: a printer of each at
// /ipp/print/<queue>, ipp and ipps on one port, advertised with DNS-SD.
type IPPConfig struct {
	// Listen is the address: ":8631" by default.
	Listen string `yaml:"listen"`
	// Host is the name of the server in the printers' URIs, the
	// certificate and DNS-SD: the machine's name, in .local, by default.
	Host string `yaml:"host"`
	// Auth is "none" (anyone of the network prints, under the
	// requesting-user-name it gives) or "basic" (a token's name and
	// secret): basic when there are tokens.
	Auth string `yaml:"auth"`
	// RequireTLS refuses ipp (HTTP 426): only ipps. It is so with basic.
	RequireTLS bool `yaml:"require-tls"`
	// TLS is the certificate of ipps; without it the server makes one,
	// signed by itself, kept next to the spool (ipp-tls.pem).
	TLS *TLSConfig `yaml:"tls"`
	// Queues are the queues shown; all by default. Default is also shown
	// at /ipp/print.
	Queues  []string `yaml:"queues"`
	Default string   `yaml:"default"`
	// Location, Organization and OrganizationalUnit describe the
	// printers (printer-location …), a queue's location over Location.
	Location           string `yaml:"location"`
	Organization       string `yaml:"organization"`
	OrganizationalUnit string `yaml:"organizational-unit"`
	// DNSSD advertises the printers on the local links (mDNS): true by
	// default.
	DNSSD *bool `yaml:"dns-sd"`
}

// TLSConfig is the certificate and key of HTTPS.
type TLSConfig struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

// TokenConfig is a token: its name (the user of its jobs), its secret
// (from a file, the environment, or the value itself), the queues it may
// submit to (all when none), and its roles: submit, manage (others' jobs,
// accounting), admin (queues, discovery, backup).
type TokenConfig struct {
	Name   string   `yaml:"name"`
	File   string   `yaml:"file"`
	Env    string   `yaml:"env"`
	Value  string   `yaml:"value"`
	Queues []string `yaml:"queues"`
	Roles  []string `yaml:"roles"`
	// ActAs lets the token name the user of a job (X-BDF-User): an
	// application that submits for its users.
	ActAs bool `yaml:"act-as"`

	secret string
}

// PrinterConfig is a printer: its URL (ipp://, ipps://, http://,
// https://), its user and password when it asks for them, and the
// fingerprint of its certificate when it is pinned; or a printer of 9100
// or LPD (socket://host:9100, lpd://host/queue), the format it takes (ps,
// the default, or pdf) and whether a PJL header goes before the document
// (by default, for PDF only).
type PrinterConfig struct {
	// Proxy makes a printer that proxies of branches fetch the jobs of
	// (IPP INFRA), with no URL: the bdf printd of a branch, behind a NAT,
	// whose upstreams name the queue.
	Proxy bool   `yaml:"proxy"`
	URL   string `yaml:"url"`
	// Format is ps or pdf for a printer of 9100 or LPD, or the language
	// of a printer of labels, receipts or forms (zpl, sbpl, tpcl, escpos,
	// escp) on any of them: a raw queue of CUPS too. Profile names its
	// model (bdf print -profile list); the URL's lang, profile, dpi and
	// width say the same of a raw printer.
	Format       string `yaml:"format"`
	Profile      string `yaml:"profile"`
	PJL          *bool  `yaml:"pjl"`
	Username     string `yaml:"username"`
	PasswordEnv  string `yaml:"password-env"`
	PasswordFile string `yaml:"password-file"`
	// Fingerprint ("sha256:…") pins the certificate of an ipps printer;
	// without it, the first certificate seen is recorded and kept.
	Fingerprint string `yaml:"fingerprint"`
	// Insecure takes any certificate.
	Insecure bool `yaml:"insecure"`
}

// QueueConfig is a queue of printers.
type QueueConfig struct {
	Printers []string `yaml:"printers"`
	Policy   string   `yaml:"policy"`
	// Defaults are the ticket's defaults, with IPP's names (sides,
	// print-color-mode, media …).
	Defaults      map[string]any `yaml:"defaults"`
	Format        string         `yaml:"format"`
	RasterIfLossy bool           `yaml:"raster-if-lossy"`
	MaxInflight   int            `yaml:"max-inflight"`
	// Info and Location describe the queue as an IPP printer
	// (printer-info, the name it is advertised with, and
	// printer-location).
	Info     string `yaml:"info"`
	Location string `yaml:"location"`
}

// RetentionConfig is how long what is done is kept ("168h", "7d").
type RetentionConfig struct {
	Documents Duration `yaml:"documents"`
	Records   Duration `yaml:"records"`
}

// LimitsConfig bounds what the server takes.
type LimitsConfig struct {
	MaxDocument    Size     `yaml:"max-document"`
	MaxJobs        int      `yaml:"max-jobs"`
	MaxBytes       Size     `yaml:"max-bytes"`
	MaxImpressions int      `yaml:"max-impressions"`
	Renders        int      `yaml:"renders"`
	RetryFor       Duration `yaml:"retry-for"`
	// Uploads and Events bound the submissions and the event streams at
	// once: 8 and 64 by default.
	Uploads int `yaml:"uploads"`
	Events  int `yaml:"events"`
}

// FontsConfig finds the fonts of documents that name them.
type FontsConfig struct {
	Dirs     []string `yaml:"dirs"`
	NoSystem bool     `yaml:"no-system"`
}

// Duration reads "90s", "24h" and "7d".
type Duration time.Duration

func (d *Duration) UnmarshalYAML(n *yaml.Node) error {
	v, err := parseDuration(n.Value)
	if err != nil {
		return err
	}
	*d = Duration(v)
	return nil
}

func parseDuration(s string) (time.Duration, error) {
	if days, ok := strings.CutSuffix(s, "d"); ok {
		n, err := strconv.ParseFloat(days, 64)
		if err != nil {
			return 0, fmt.Errorf("duration %q", s)
		}
		return time.Duration(n * float64(24*time.Hour)), nil
	}
	v, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("duration %q", s)
	}
	return v, nil
}

// Size reads bytes, "512KiB", "256MiB", "10GiB" (or KB, MB, GB).
type Size int64

func (s *Size) UnmarshalYAML(n *yaml.Node) error {
	v, err := parseSize(n.Value)
	if err != nil {
		return err
	}
	*s = Size(v)
	return nil
}

func parseSize(s string) (int64, error) {
	units := []struct {
		suffix string
		n      int64
	}{{"KiB", 1 << 10}, {"MiB", 1 << 20}, {"GiB", 1 << 30}, {"TiB", 1 << 40}, {"KB", 1e3}, {"MB", 1e6}, {"GB", 1e9}, {"TB", 1e12}, {"B", 1}}
	t := strings.TrimSpace(s)
	for _, u := range units {
		if num, ok := strings.CutSuffix(t, u.suffix); ok {
			f, err := strconv.ParseFloat(strings.TrimSpace(num), 64)
			if err != nil || f < 0 {
				return 0, fmt.Errorf("size %q", s)
			}
			return int64(f * float64(u.n)), nil
		}
	}
	n, err := strconv.ParseInt(t, 10, 64)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("size %q", s)
	}
	return n, nil
}

// LoadConfig reads a configuration file, with its secrets, and checks it.
func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c, err := ParseConfig(b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	c.dir = filepath.Dir(abs)
	if err := c.resolve(); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return c, nil
}

// ParseConfig reads a configuration from YAML and checks it; paths are
// relative to the working directory.
func ParseConfig(b []byte) (*Config, error) {
	var c Config
	dec := yaml.NewDecoder(strings.NewReader(string(b)))
	dec.KnownFields(true)
	// an empty file is a configuration of defaults
	if err := dec.Decode(&c); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if c.Listen == "" {
		c.Listen = "127.0.0.1:8640"
	}
	if c.User == "" {
		c.User = "bdf"
	}
	return &c, c.check()
}

// check checks what does not need the files and the environment.
func (c *Config) check() error {
	if len(c.Printers) == 0 {
		return errors.New("no printers")
	}
	if len(c.Queues) == 0 {
		return errors.New("no queues")
	}
	for name, p := range c.Printers {
		if p.Proxy {
			if p.URL != "" || p.Format != "" || p.Profile != "" || p.PJL != nil || p.Username != "" || p.PasswordEnv != "" || p.PasswordFile != "" || p.Fingerprint != "" || p.Insecure {
				return fmt.Errorf("printer %s: a printer of proxies has no url, format, pjl, user, password or certificate", name)
			}
			continue
		}
		if raw.IsURL(p.URL) {
			if _, err := p.raw(c.User); err != nil {
				return fmt.Errorf("printer %s: %w", name, err)
			}
			if p.Username != "" || p.PasswordEnv != "" || p.PasswordFile != "" || p.Fingerprint != "" || p.Insecure {
				return fmt.Errorf("printer %s: a printer of 9100, LPD or a device has no user, password or certificate", name)
			}
			continue
		}
		if p.PJL != nil {
			return fmt.Errorf("printer %s: pjl is for printers of 9100 or LPD (socket://, lpd://)", name)
		}
		if p.Format != "" || p.Profile != "" {
			if _, err := p.profile(); err != nil {
				return fmt.Errorf("printer %s: %w", name, err)
			}
		}
		if _, _, err := ipp.ParseURL(p.URL); err != nil {
			return fmt.Errorf("printer %s: %w", name, err)
		}
		if p.Fingerprint != "" && !strings.HasPrefix(p.Fingerprint, "sha256:") {
			return fmt.Errorf("printer %s: fingerprint %q: sha256:…", name, p.Fingerprint)
		}
	}
	for name, q := range c.Queues {
		if len(q.Printers) == 0 {
			return fmt.Errorf("queue %s: no printers", name)
		}
		for _, p := range q.Printers {
			if _, ok := c.Printers[p]; !ok {
				return fmt.Errorf("queue %s: no printer %s", name, p)
			}
			if c.Printers[p].Proxy && len(q.Printers) > 1 {
				return fmt.Errorf("queue %s: printer %s of proxies is the only printer of its queue", name, p)
			}
		}
		if _, err := q.ticket(); err != nil {
			return fmt.Errorf("queue %s: defaults: %w", name, err)
		}
		switch q.Policy {
		case "", "any", "failover":
		default:
			return fmt.Errorf("queue %s: policy %q: any or failover", name, q.Policy)
		}
		if _, ok := print.FormatName(q.Format); !ok {
			return fmt.Errorf("queue %s: format %q: auto, pdf, ps, pwg-raster, urf or pclm", name, q.Format)
		}
		if q.MaxInflight < 0 {
			return fmt.Errorf("queue %s: max-inflight %d", name, q.MaxInflight)
		}
	}
	seen := map[string]bool{}
	for _, t := range c.Tokens {
		if t.Name == "" {
			return errors.New("a token without name")
		}
		if seen[t.Name] {
			return fmt.Errorf("token %s twice", t.Name)
		}
		seen[t.Name] = true
		for _, r := range t.Roles {
			if r != roleSubmit && r != roleManage && r != roleAdmin && r != roleProxy {
				return fmt.Errorf("token %s: role %q: submit, manage, admin or proxy", t.Name, r)
			}
		}
		for _, q := range t.Queues {
			if _, ok := c.Queues[q]; !ok {
				return fmt.Errorf("token %s: no queue %s", t.Name, q)
			}
		}
	}
	if c.TLS != nil && (c.TLS.Cert == "" || c.TLS.Key == "") {
		return errors.New("tls: cert and key")
	}
	if err := c.checkIPP(); err != nil {
		return err
	}
	for name, u := range c.Upstreams {
		if _, _, err := ipp.ParseURL(u.URL); err != nil {
			return fmt.Errorf("upstream %s: %w", name, err)
		}
		if _, ok := c.Queues[u.Queue]; !ok {
			return fmt.Errorf("upstream %s: no queue %q", name, u.Queue)
		}
		if u.Fingerprint != "" && !strings.HasPrefix(u.Fingerprint, "sha256:") {
			return fmt.Errorf("upstream %s: fingerprint %q: sha256:…", name, u.Fingerprint)
		}
		if u.UUID != "" && !strings.HasPrefix(u.UUID, "urn:uuid:") {
			return fmt.Errorf("upstream %s: uuid %q: urn:uuid:…", name, u.UUID)
		}
	}
	proxyQueue := map[string]string{}
	for _, qn := range sortedKeys(c.Queues) {
		for _, p := range c.Queues[qn].Printers {
			if !c.Printers[p].Proxy {
				continue
			}
			if c.IPP == nil {
				return fmt.Errorf("printer %s of proxies needs ipp: the proxies fetch its jobs over IPP", p)
			}
			// (a token of a queue fetches the jobs of its printer: of that
			// queue only)
			if other, ok := proxyQueue[p]; ok {
				return fmt.Errorf("printer %s of proxies is the printer of queues %s and %s: of one queue only", p, other, qn)
			}
			proxyQueue[p] = qn
		}
	}
	if !c.Insecure && !loopback(c.Listen) && (c.TLS == nil || len(c.Tokens) == 0) {
		return fmt.Errorf("listen %s is not a loopback address: it needs tls and tokens (or insecure: true)", c.Listen)
	}
	return nil
}

// resolve reads the secrets and makes the paths absolute.
func (c *Config) resolve() error {
	abs := func(p string) string {
		if p == "" || filepath.IsAbs(p) {
			return p
		}
		return filepath.Join(c.dir, p)
	}
	if c.Spool == "" {
		c.Spool = "spool.db"
	}
	c.Spool = abs(c.Spool)
	c.TempDir = abs(c.TempDir)
	if c.TLS != nil {
		c.TLS.Cert, c.TLS.Key = abs(c.TLS.Cert), abs(c.TLS.Key)
	}
	if c.IPP != nil && c.IPP.TLS != nil {
		c.IPP.TLS.Cert, c.IPP.TLS.Key = abs(c.IPP.TLS.Cert), abs(c.IPP.TLS.Key)
	}
	for i := range c.Fonts.Dirs {
		c.Fonts.Dirs[i] = abs(c.Fonts.Dirs[i])
	}
	for i := range c.Tokens {
		t := &c.Tokens[i]
		switch {
		case t.File != "":
			b, err := os.ReadFile(abs(t.File))
			if err != nil {
				return fmt.Errorf("token %s: %w", t.Name, err)
			}
			t.secret = strings.TrimSpace(string(b))
		case t.Env != "":
			t.secret = os.Getenv(t.Env)
		default:
			t.secret = t.Value
		}
		if len(t.secret) < 16 {
			return fmt.Errorf("token %s: the secret is shorter than 16 characters", t.Name)
		}
	}
	for name, p := range c.Printers {
		if p.PasswordFile != "" {
			p.PasswordFile = abs(p.PasswordFile)
			c.Printers[name] = p
		}
	}
	for name, u := range c.Upstreams {
		if u.PasswordFile != "" {
			u.PasswordFile = abs(u.PasswordFile)
			c.Upstreams[name] = u
		}
	}
	return nil
}

// checkIPP checks the IPP server's configuration and fills its defaults.
func (c *Config) checkIPP() error {
	p := c.IPP
	if p == nil {
		return nil
	}
	if p.Listen == "" {
		p.Listen = ":8631"
	}
	if _, _, err := net.SplitHostPort(p.Listen); err != nil {
		return fmt.Errorf("ipp: listen %q: %w", p.Listen, err)
	}
	switch p.Auth {
	case "":
		p.Auth = authNone
		if len(c.Tokens) > 0 {
			p.Auth = authBasic
		}
	case authNone:
	case authBasic:
		if len(c.Tokens) == 0 {
			return errors.New("ipp: auth basic needs tokens")
		}
	default:
		return fmt.Errorf("ipp: auth %q: none or basic", p.Auth)
	}
	if p.TLS != nil && (p.TLS.Cert == "" || p.TLS.Key == "") {
		return errors.New("ipp: tls: cert and key")
	}
	for _, q := range p.Queues {
		if _, ok := c.Queues[q]; !ok {
			return fmt.Errorf("ipp: no queue %s", q)
		}
	}
	if p.Default != "" {
		if _, ok := c.Queues[p.Default]; !ok {
			return fmt.Errorf("ipp: default: no queue %s", p.Default)
		}
		if len(p.Queues) > 0 && !slices.Contains(p.Queues, p.Default) {
			return fmt.Errorf("ipp: default %s is not among the queues shown", p.Default)
		}
	}
	for name, q := range c.Queues {
		if !validQueueName(name) {
			return fmt.Errorf("queue %s: a queue shown by IPP is named with letters, digits, - and _", name)
		}
		if len(q.Info) > 127 || len(q.Location) > 127 {
			return fmt.Errorf("queue %s: info and location have 127 bytes at most", name)
		}
	}
	return nil
}

// validQueueName tells a name that is a path segment as it is.
func validQueueName(s string) bool {
	if s == "" || len(s) > 127 {
		return false
	}
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.') {
			return false
		}
	}
	return true
}

// ticket reads a queue's defaults as a ticket, with the JSON names.
func (q QueueConfig) ticket() (print.Ticket, error) {
	var t print.Ticket
	if len(q.Defaults) == 0 {
		return t, nil
	}
	b, err := json.Marshal(q.Defaults)
	if err != nil {
		return t, err
	}
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&t); err != nil {
		return t, err
	}
	return t, t.Check()
}

// client makes the IPP client of a printer.
// raw returns the printer of 9100 or LPD of the configuration.
func (p PrinterConfig) raw(user string) (*raw.Printer, error) {
	r, err := raw.Parse(p.URL)
	if err != nil {
		return nil, err
	}
	switch p.Format {
	case "":
	case "ps", "postscript":
		r.Format = raw.FormatPostScript
	case "pdf":
		r.Format = raw.FormatPDF
	default:
		if pdl.Language(p.Format) == "" {
			return nil, fmt.Errorf("format %q: ps, pdf, zpl, sbpl, tpcl, escpos or escp", p.Format)
		}
		r.Format = pdl.Language(p.Format)
	}
	if p.Profile != "" {
		if !raw.IsLanguage(r.Format) {
			prof, ok := pdl.LookupProfile(p.Profile)
			if !ok {
				return nil, fmt.Errorf("profile %q: not a model (bdf print -profile list)", p.Profile)
			}
			if p.Format != "" {
				return nil, fmt.Errorf("profile %s is of %s, not %s", p.Profile, pdl.Title(prof.Language), p.Format)
			}
			r.Format = prof.Language
		}
		r.Profile = p.Profile
	}
	if raw.IsLanguage(r.Format) {
		// the URL's profile, dpi and width, checked now
		if _, err := pdl.ParseProfile(r.Profile, r.Format, r.DPI, r.Width); err != nil {
			return nil, err
		}
	}
	if p.PJL != nil {
		r.PJL = fmt.Sprint(*p.PJL)
	}
	if r.User == "" {
		r.User = user
	}
	return r, nil
}

// profile returns the profile of an IPP printer of a language (a raw
// queue of CUPS): its profile, or the default one of its format.
func (p PrinterConfig) profile() (pdl.Profile, error) {
	if p.Format != "" && pdl.Language(p.Format) == "" {
		return pdl.Profile{}, fmt.Errorf("format %q: an IPP printer is given the language of a printer of labels, receipts or forms (zpl, sbpl, tpcl, escpos, escp), or none", p.Format)
	}
	return pdl.ParseProfile(p.Profile, p.Format, "", "")
}

func (p PrinterConfig) client(user string) (*ipp.Client, error) {
	c := &ipp.Client{URL: p.URL, User: user, Username: p.Username, Timeout: 5 * time.Minute}
	switch {
	case p.PasswordFile != "":
		b, err := os.ReadFile(p.PasswordFile)
		if err != nil {
			return nil, err
		}
		c.Password = strings.TrimSpace(string(b))
	case p.PasswordEnv != "":
		c.Password = os.Getenv(p.PasswordEnv)
	}
	c.TLS.Insecure = p.Insecure
	return c, nil
}

// spoolOptions turns the configuration into the spool's.
func (c *Config) spoolOptions() (spool.Options, error) {
	o := spool.Options{Path: c.Spool, TempDir: c.TempDir, FontDirs: c.Fonts.Dirs, NoSystemFonts: c.Fonts.NoSystem,
		Retention: spool.Retention{Documents: time.Duration(c.Retention.Documents), Records: time.Duration(c.Retention.Records)},
		Limits: spool.Limits{MaxDocument: int64(c.Limits.MaxDocument), MaxJobs: c.Limits.MaxJobs, MaxBytes: int64(c.Limits.MaxBytes),
			MaxImpressions: c.Limits.MaxImpressions, Renders: c.Limits.Renders, RetryFor: time.Duration(c.Limits.RetryFor)}}
	for _, name := range sortedKeys(c.Printers) {
		pc := c.Printers[name]
		if pc.Proxy {
			o.Printers = append(o.Printers, spool.Printer{Name: name, Proxy: true})
			continue
		}
		if raw.IsURL(pc.URL) {
			r, err := pc.raw(c.User)
			if err != nil {
				return o, fmt.Errorf("printer %s: %w", name, err)
			}
			o.Printers = append(o.Printers, spool.Printer{Name: name, Raw: r})
			continue
		}
		cl, err := pc.client(c.User)
		if err != nil {
			return o, fmt.Errorf("printer %s: %w", name, err)
		}
		sp := spool.Printer{Name: name, Client: cl}
		if pc.Format != "" || pc.Profile != "" {
			prof, err := pc.profile()
			if err != nil {
				return o, fmt.Errorf("printer %s: %w", name, err)
			}
			sp.Format, sp.Profile = prof.Language, &prof
		}
		o.Printers = append(o.Printers, sp)
	}
	for _, name := range sortedKeys(c.Queues) {
		q := c.Queues[name]
		t, _ := q.ticket()
		o.Queues = append(o.Queues, spool.Queue{Name: name, Printers: q.Printers, Policy: q.Policy, Defaults: t, Format: q.Format,
			RasterIfLossy: q.RasterIfLossy, MaxInflight: q.MaxInflight})
	}
	return o, nil
}

// pinned are the fingerprints set in the configuration, by host.
func (c *Config) pinned() map[string]string {
	m := map[string]string{}
	for _, p := range c.Printers {
		if p.Fingerprint == "" {
			continue
		}
		if u, err := url.Parse(p.URL); err == nil {
			m[u.Hostname()] = p.Fingerprint
		}
	}
	for _, p := range c.Upstreams {
		if p.Fingerprint == "" {
			continue
		}
		if u, err := url.Parse(p.URL); err == nil {
			m[u.Hostname()] = p.Fingerprint
		}
	}
	return m
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// loopback tells an address that only the machine reaches.
func loopback(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
