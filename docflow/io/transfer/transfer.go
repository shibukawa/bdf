// Package transfer is what the file-transfer kinds of docflow share
// (docflow/io/ftp, docflow/io/sftp): a far side's files behind one
// interface (FS), a node's far side as its URL says it (Site), the
// settings and the fields of an exit that uploads under a temporary
// name and renames into place (Send), and of an entry that polls a
// directory at every tick, claims a file by renaming it, takes it in as
// a job and deletes or moves it (Poll). A kind of one's own over another
// protocol (WebDAV) uses the same.
package transfer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/docflow"
)

// Entry is a file of a far side's directory as it is listed.
type Entry struct {
	Name    string
	Size    int64
	ModTime time.Time
	Dir     bool
}

// FS is a far side's files, as a protocol's client gives them. Paths are
// slash-separated, as the server has them (Site.Path makes them of a
// node's base and a name). A file not there is fs.ErrNotExist.
type FS interface {
	// List lists a directory (not . and ..).
	List(ctx context.Context, dir string) ([]Entry, error)
	// Stat tells a file.
	Stat(ctx context.Context, name string) (*Entry, error)
	// Open reads a file.
	Open(ctx context.Context, name string) (io.ReadCloser, error)
	// Create writes a file, made or truncated; Close ends the write and
	// tells whether it all went.
	Create(ctx context.Context, name string) (io.WriteCloser, error)
	// Rename renames a file; with replace, over a file at the new name
	// (without it, that is an error on most servers).
	Rename(ctx context.Context, from, to string, replace bool) error
	// Remove removes a file.
	Remove(ctx context.Context, name string) error
	// Mkdir makes a directory, whose parent is there; one there already
	// is an error.
	Mkdir(ctx context.Context, dir string) error
	Close() error
}

// Dialer connects to a node's far side: with its credential read, which
// an exit or an entry does at every send or poll (never at Check).
type Dialer func(ctx context.Context) (FS, error)

// Site is a node's far side as its URL says: <scheme>://user@host:port/dir.
// The directory is as the server has it (absolute); /~/dir is one under
// the login's own directory; none is the login's own.
type Site struct {
	Scheme string
	Host   string
	Port   int
	User   string
	// Base is the directory, with no trailing slash: absolute, relative
	// (after /~/), or "" (the login's own).
	Base string
	port int // the scheme's default, left out of URL()
}

// ParseURL reads a node's URL: one of the schemes given, with its
// default port; a password in the URL is refused (the credential is
// named apart).
func ParseURL(raw string, ports map[string]int) (*Site, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil, fmt.Errorf("url %q: %v", raw, err)
	}
	port, ok := ports[u.Scheme]
	if !ok || u.Hostname() == "" || u.RawQuery != "" || u.Fragment != "" {
		var schemes []string
		for s := range ports {
			schemes = append(schemes, s)
		}
		sort.Strings(schemes)
		return nil, fmt.Errorf("url %q: %s://user@host:port/dir", raw, strings.Join(schemes, "|"))
	}
	if _, has := u.User.Password(); has {
		return nil, fmt.Errorf("url %q has a password: name a credential instead", raw)
	}
	s := &Site{Scheme: u.Scheme, Host: u.Hostname(), Port: port, User: u.User.Username(), port: port}
	if p := u.Port(); p != "" {
		if s.Port, err = strconv.Atoi(p); err != nil || s.Port < 1 || s.Port > 65535 {
			return nil, fmt.Errorf("url %q: port %q", raw, p)
		}
	}
	dir := u.Path
	switch {
	case dir == "" || dir == "/~" || dir == "/~/":
		s.Base = ""
	case strings.HasPrefix(dir, "/~/"):
		s.Base = docflow.CleanPath(dir[3:])
	default:
		s.Base = "/" + docflow.CleanPath(dir)
	}
	return s, nil
}

// Address is host:port.
func (s *Site) Address() string { return s.Host + ":" + strconv.Itoa(s.Port) }

// Dir is the base directory as a client lists it: "." for the login's
// own.
func (s *Site) Dir() string {
	if s.Base == "" {
		return "."
	}
	return s.Base
}

// Path is a file's path on the server: rel (clean, relative) under the
// base.
func (s *Site) Path(rel string) string {
	switch {
	case rel == "":
		return s.Dir()
	case s.Base == "", s.Base == "/":
		return s.Base + rel
	}
	return s.Base + "/" + rel
}

// URL is the URL of a file on the server (no user), the way a delivery
// records where it put it.
func (s *Site) URL(p string) string {
	host := s.Host
	if s.Port != s.port {
		host = s.Address()
	}
	if !strings.HasPrefix(p, "/") {
		p = "/~/" + p
	}
	return s.Scheme + "://" + host + p
}

// What an exit's settings are (Node.Settings), besides a kind's own.
type ExitSettings struct {
	URL        string `json:"url"`
	Credential string `json:"credential"`
	// Path is the file's path under the base, expanding {{meta.<key>}}
	// (DefaultPath: the job's number, and the payload's name).
	Path string `json:"path"`
	// Mkdirs makes the directories of the path that are not there.
	Mkdirs bool `json:"mkdirs"`
	// Overwrite writes over a file that is there; otherwise one that is
	// there is left (a retry after a write finds it and is done).
	Overwrite bool `json:"overwrite"`
}

// DefaultPath is the path an exit writes by default.
const DefaultPath = "{{meta.$job}}/{{meta.$name}}"

// ExitFields are the fields of an exit's settings: the URL and the
// credential (a kind's own words), the kind's own fields, then the
// path, mkdirs and overwrite.
func ExitFields(url, credential docflow.Field, own ...docflow.Field) []docflow.Field {
	out := []docflow.Field{url, credential}
	out = append(out, own...)
	return append(out,
		docflow.Field{Key: "path", Type: docflow.FieldPath, Default: DefaultPath, Expands: true,
			Title: docflow.Text{"en": "Path", "ja": "パス"},
			Help:  docflow.Text{"en": "The file's path under the URL's directory: the job's number and the file's name by default.", "ja": "URL のディレクトリの下のファイルのパス。既定はジョブの番号とファイルの名前"}},
		docflow.Field{Key: "mkdirs", Type: docflow.FieldBool, Default: true,
			Title: docflow.Text{"en": "Make directories", "ja": "ディレクトリを作る"},
			Help:  docflow.Text{"en": "Make the directories of the path that are not there.", "ja": "パスのディレクトリが無ければ作る"}},
		docflow.Field{Key: "overwrite", Type: docflow.FieldBool,
			Title: docflow.Text{"en": "Overwrite", "ja": "上書き"},
			Help:  docflow.Text{"en": "Write over a file that is there; otherwise it is left as it is, as written by an earlier try.", "ja": "既にあるファイルに上書きする。しないなら、あるものを前の試みが書いたものとしてそのまま残す"}})
}

// URLField is the URL field of a kind, with its placeholder and help.
func URLField(placeholder string, help docflow.Text) docflow.Field {
	return docflow.Field{Key: "url", Type: docflow.FieldString, Required: true, Title: docflow.Text{"en": "URL", "ja": "URL"}, Help: help, Placeholder: docflow.Text{"en": placeholder}}
}

// CredentialField is the credential field of a kind, with its help.
func CredentialField(help docflow.Text, required bool) docflow.Field {
	return docflow.Field{Key: "credential", Type: docflow.FieldCredential, Required: required, Title: docflow.Text{"en": "Credential", "ja": "資格"}, Help: help}
}

// Send writes what reaches an exit node to its far side: under a
// temporary name (.<delivery>.part), renamed into place, so that a
// reader never sees half a file. A file at the final name already is
// taken as written by a try whose outcome was not known, unless the
// node overwrites. It returns the file's URL, which the delivery
// records.
func Send(ctx context.Context, o *docflow.Out, s ExitSettings, site *Site, dial Dialer) (string, error) {
	rel := o.ExpandPath(firstOr(s.Path, DefaultPath))
	if rel == "" {
		return "", docflow.Final(errors.New("the path is empty for this job"))
	}
	dest := site.Path(rel)
	f, err := dial(ctx)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if !s.Overwrite {
		if _, err := f.Stat(ctx, dest); err == nil {
			return site.URL(dest), nil // written before
		}
	}
	dir := path.Dir(dest)
	if s.Mkdirs {
		if err := MkdirAll(ctx, f, dir); err != nil {
			return "", err
		}
	}
	tmp := path.Join(dir, "."+o.ID+".part")
	if err := upload(ctx, o, f, tmp); err != nil {
		f.Remove(ctx, tmp)
		return "", err
	}
	if err := f.Rename(ctx, tmp, dest, s.Overwrite); err != nil {
		if !s.Overwrite {
			if _, serr := f.Stat(ctx, dest); serr == nil {
				f.Remove(ctx, tmp)
				return site.URL(dest), nil // written meanwhile by another try
			}
		}
		f.Remove(ctx, tmp)
		return "", err
	}
	return site.URL(dest), nil
}

// upload writes the payload to a file of the far side.
func upload(ctx context.Context, o *docflow.Out, f FS, name string) error {
	r, err := o.Open(ctx)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := f.Create(ctx, name)
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, r); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

// MkdirAll makes a directory and its parents that are not there: a
// Mkdir that fails is checked with Stat.
func MkdirAll(ctx context.Context, f FS, dir string) error {
	if dir == "" || dir == "." || dir == "/" {
		return nil
	}
	if e, err := f.Stat(ctx, dir); err == nil {
		if !e.Dir {
			return fmt.Errorf("%s is a file, not a directory", dir)
		}
		return nil
	}
	parent := path.Dir(dir)
	if parent != dir {
		if err := MkdirAll(ctx, f, parent); err != nil {
			return err
		}
	}
	if err := f.Mkdir(ctx, dir); err != nil {
		if e, serr := f.Stat(ctx, dir); serr == nil && e.Dir {
			return nil // made meanwhile
		}
		return err
	}
	return nil
}

// What an entry's settings are (Node.Settings), besides a kind's own.
type PollSettings struct {
	URL        string `json:"url"`
	Credential string `json:"credential"`
	// Pattern is the names taken (path.Match's form; * every file).
	Pattern string `json:"pattern"`
	// Every is how often the directory is looked at, however often Tick
	// is called.
	Every docflow.Duration `json:"every"`
	// After is what is done with a file taken: AfterDelete or AfterMove
	// (to DoneDir, under the base unless absolute).
	After   string `json:"after"`
	DoneDir string `json:"done_dir"`
	// Settle is how long a file is left after it was last changed, so
	// that one still being written is not taken; a file seen unchanged
	// at the last poll is taken whatever the time.
	Settle docflow.Duration `json:"settle"`
	// MaxPerTick bounds the files taken at one tick.
	MaxPerTick int `json:"max_per_tick"`
}

// What is done with a file taken (PollSettings.After).
const (
	AfterDelete = "delete"
	AfterMove   = "move"
)

// The defaults of an entry's settings.
const (
	DefaultEvery      = "1m"
	DefaultSettle     = "10s"
	DefaultMaxPerTick = 20
)

// PollFields are the fields of an entry's settings: the URL and the
// credential (a kind's own words), the kind's own fields, then the
// pattern, every, after, done_dir, settle and max_per_tick.
func PollFields(url, credential docflow.Field, own ...docflow.Field) []docflow.Field {
	out := []docflow.Field{url, credential}
	out = append(out, own...)
	return append(out,
		docflow.Field{Key: "pattern", Type: docflow.FieldString, Default: "*",
			Title: docflow.Text{"en": "Pattern", "ja": "パターン"},
			Help:  docflow.Text{"en": "The names taken: * every file, *.pdf, order-*.csv.", "ja": "取る名前。* はすべてのファイル、*.pdf、order-*.csv など"}},
		docflow.Field{Key: "every", Type: docflow.FieldDuration, Default: DefaultEvery,
			Title: docflow.Text{"en": "Every", "ja": "間隔"},
			Help:  docflow.Text{"en": "How often the directory is looked at (30s, 5m, 1h), however often the tick comes.", "ja": "ディレクトリを見回る間隔（30s、5m、1h）。tick がそれより密でも守る"}},
		docflow.Field{Key: "after", Type: docflow.FieldEnum, Default: AfterDelete,
			Choices: []docflow.Choice{{Value: AfterDelete, Label: docflow.Text{"en": "delete", "ja": "消す"}}, {Value: AfterMove, Label: docflow.Text{"en": "move to a directory", "ja": "ディレクトリへ移す"}}},
			Title:   docflow.Text{"en": "After taking", "ja": "取った後"},
			Help:    docflow.Text{"en": "What is done with a file taken in as a job: deleted, or moved to the done directory (what must be kept goes to an archive).", "ja": "ジョブにしたファイルをどうするか。消すか、済みのディレクトリへ移す（残す約束が要るものはアーカイブへ）"}},
		docflow.Field{Key: "done_dir", Type: docflow.FieldString, Required: true, If: map[string]string{"after": AfterMove},
			Title:       docflow.Text{"en": "Done directory", "ja": "済みのディレクトリ"},
			Help:        docflow.Text{"en": "Where a file taken is moved: under the URL's directory, or absolute. It is made if it is not there.", "ja": "取ったファイルを移す先。URL のディレクトリの下か絶対パス。無ければ作る"},
			Placeholder: docflow.Text{"en": "done"}},
		docflow.Field{Key: "settle", Type: docflow.FieldDuration, Default: DefaultSettle,
			Title: docflow.Text{"en": "Settle", "ja": "落ち着く時間"},
			Help:  docflow.Text{"en": "A file changed more recently than this is left for the next time, unless it was seen the same at the last look: one still being written is not taken.", "ja": "更新がこの時間より新しいファイルは次に回す（前回と同じまま見えたものは取る）。書いている途中のものを取らないため"}},
		docflow.Field{Key: "max_per_tick", Type: docflow.FieldInt, Default: DefaultMaxPerTick,
			Title: docflow.Text{"en": "Files per tick", "ja": "1 回に取る数"},
			Help:  docflow.Text{"en": "At most this many files are taken at one tick; the rest wait for the next.", "ja": "1 回の tick で取るファイルの数の上限。残りは次の tick"}})
}

// CheckPoll checks an entry's settings beyond their fields: the pattern,
// and what is done after.
func CheckPoll(s PollSettings) error {
	if _, err := path.Match(firstOr(s.Pattern, "*"), "x"); err != nil {
		return fmt.Errorf("pattern %q: %v", s.Pattern, err)
	}
	switch s.After {
	case "", AfterDelete:
	case AfterMove:
		if strings.TrimSpace(s.DoneDir) == "" {
			return errors.New("a done directory (done_dir) is where a file is moved to")
		}
	default:
		return fmt.Errorf("after %q: delete or move", s.After)
	}
	if s.MaxPerTick < 0 {
		return errors.New("files per tick (max_per_tick) is at least 1")
	}
	return nil
}

// pollState is what an entry keeps between ticks (docflow.KindState).
type pollState struct {
	// Last is when the directory was last looked at.
	Last time.Time `json:"last,omitzero"`
	// Seen are the files seen and left at the last look, by name: their
	// size and time, for a file unchanged since to be taken.
	Seen map[string]string `json:"seen,omitempty"`
}

// maxSeen bounds the files remembered between ticks.
const maxSeen = 2000

// The suffixes of the files an entry renames: one it is taking (claimed,
// not yet taken), and one the engine refused (a format it does not
// read), left for a person.
const (
	takingSuffix  = ".taking"
	refusedSuffix = ".refused"
)

// Poll looks at an entry node's directory, every so often, and takes
// the files there in as jobs: one is claimed by renaming it to
// <name>.docflow-<node>.taking (so that another instance polling at
// the same time skips it), read, sent to the node as a job (Submit,
// with the file's path, size and time as <kind>.path, .size and .time
// in the metadata, the login as the claimed user), then deleted or
// moved. A file the engine refuses for good (not a format it reads) is
// renamed .refused and left; one it cannot take now stays claimed and
// is tried again at a later tick. The files taken are counted in the
// report (Polled), what was left and why is noted. kind is the kind's
// name (the job's source), prefix its key in the metadata ("sftp").
func Poll(ctx context.Context, m *docflow.Manager, n *docflow.Node, rep *docflow.TickReport, s PollSettings, site *Site, kind, prefix string, dial Dialer) error {
	var st pollState
	if err := m.KindState(ctx, n, &st); err != nil {
		return err
	}
	now := m.Now()
	if every := time.Duration(s.Every); every > 0 && !st.Last.IsZero() && now.Sub(st.Last) < every {
		return nil
	}
	f, err := dial(ctx)
	if err != nil {
		rep.Note("%s: %v", n.Name, err)
		return err
	}
	defer f.Close()
	entries, err := f.List(ctx, site.Dir())
	if err != nil {
		rep.Note("%s: list %s: %v", n.Name, site.Dir(), err)
		return err
	}
	mine := ".docflow-" + n.Name
	pattern := firstOr(s.Pattern, "*")
	max := s.MaxPerTick
	if max <= 0 {
		max = DefaultMaxPerTick
	}
	settle := time.Duration(s.Settle)
	limit := m.Options().MaxDocument
	// what to take: the files claimed earlier and left (an earlier tick
	// fell over), then the new ones that fit, in order of name
	type take struct {
		name, claimed string
	}
	var claimed, fresh []take
	seen := map[string]string{}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name < entries[j].Name })
	for _, e := range entries {
		if e.Dir || e.Name == "" || strings.HasPrefix(e.Name, ".") {
			continue
		}
		if i := strings.LastIndex(e.Name, ".docflow-"); i >= 0 {
			if strings.HasSuffix(e.Name, mine+takingSuffix) {
				claimed = append(claimed, take{name: e.Name[:i], claimed: e.Name})
			}
			continue // another node's, or refused
		}
		if strings.HasSuffix(e.Name, ".part") || strings.HasSuffix(e.Name, ".tmp") || strings.HasSuffix(e.Name, ".filepart") {
			continue
		}
		if ok, _ := path.Match(pattern, e.Name); !ok {
			continue
		}
		if e.Size > limit {
			rep.Note("%s: %s is %d bytes, more than the %d taken: left", n.Name, e.Name, e.Size, limit)
			continue
		}
		key := strconv.FormatInt(e.Size, 10) + "/" + strconv.FormatInt(e.ModTime.Unix(), 10)
		if settle > 0 && (e.ModTime.IsZero() || now.Sub(e.ModTime) < settle) && st.Seen[e.Name] != key {
			if len(seen) < maxSeen {
				seen[e.Name] = key
			}
			continue
		}
		fresh = append(fresh, take{name: e.Name})
	}
	var errs []error
	taken := 0
	for _, t := range append(claimed, fresh...) {
		if taken >= max {
			break
		}
		if err := ctx.Err(); err != nil {
			errs = append(errs, err)
			break
		}
		claimed := t.claimed
		if claimed == "" {
			claimed = t.name + mine + takingSuffix
			if err := f.Rename(ctx, site.Path(t.name), site.Path(claimed), false); err != nil {
				continue // another instance took it, or it is gone
			}
		}
		taken++ // a try counts, whatever came of it: the bound is on the work
		took, err := takeIn(ctx, m, n, f, site, t.name, claimed, s, kind, prefix, limit)
		if err != nil {
			rep.Note("%s: %s: %v", n.Name, t.name, err)
			errs = append(errs, fmt.Errorf("%s: %w", t.name, err))
		}
		if took {
			rep.Polled++
		}
	}
	st.Last, st.Seen = now, seen
	if err := m.PutKindState(ctx, n, st); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// takeIn reads a claimed file, makes a job of it and does with the file
// what the settings say. took tells a job was made (whatever happened
// to the file after).
func takeIn(ctx context.Context, m *docflow.Manager, n *docflow.Node, f FS, site *Site, name, claimed string, s PollSettings, kind, prefix string, limit int64) (took bool, err error) {
	src := site.Path(claimed)
	e, err := f.Stat(ctx, src)
	if err != nil {
		return false, err
	}
	tmp, err := os.CreateTemp(m.Options().TempDir, "bdf-poll-*")
	if err != nil {
		return false, err
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()
	r, err := f.Open(ctx, src)
	if err != nil {
		return false, err
	}
	size, err := io.Copy(tmp, io.LimitReader(r, limit+1))
	r.Close()
	if err != nil {
		return false, err
	}
	if size > limit {
		// grown since it was listed: back to its name, for a later look
		f.Rename(ctx, src, site.Path(name), false)
		return false, fmt.Errorf("%d bytes, more than the %d taken: left", size, limit)
	}
	meta := map[string]any{prefix + ".path": site.Path(name), prefix + ".size": size}
	if !e.ModTime.IsZero() {
		meta[prefix+".time"] = e.ModTime.UTC().Format(time.RFC3339)
	}
	j, err := m.Submit(ctx, n.Name, docflow.Input{Name: name, ClaimedUser: site.User, Source: kind, Doc: tmp, Size: size, Meta: meta})
	if err != nil {
		var refusal *docflow.Refusal
		if errors.Is(err, docflow.ErrFormat) || errors.Is(err, docflow.ErrInvalid) || errors.As(err, &refusal) {
			// nothing a later try mends: left for a person, as .refused
			m.Options().Log.Warn("poll: refused", "node", n.Name, "file", site.Path(name), "err", err)
			if rerr := f.Rename(ctx, src, site.Path(name+".docflow-"+n.Name+refusedSuffix), true); rerr != nil {
				return false, errors.Join(err, rerr)
			}
			return false, fmt.Errorf("refused, left as %s: %w", refusedSuffix, err)
		}
		return false, err
	}
	m.Options().Log.Info("poll: took", "node", n.Name, "file", site.Path(name), "job", j.Number)
	switch s.After {
	case AfterMove:
		done := s.DoneDir
		if !strings.HasPrefix(done, "/") {
			done = site.Path(docflow.CleanPath(done))
		}
		if err := MkdirAll(ctx, f, done); err != nil {
			return true, fmt.Errorf("job %d made; the done directory: %w", j.Number, err)
		}
		if err := f.Rename(ctx, src, path.Join(done, name), true); err != nil {
			return true, fmt.Errorf("job %d made; moving it: %w", j.Number, err)
		}
	default:
		if err := f.Remove(ctx, src); err != nil {
			return true, fmt.Errorf("job %d made; deleting it: %w", j.Number, err)
		}
	}
	return true, nil
}

func firstOr(s, def string) string {
	if s != "" {
		return s
	}
	return def
}
