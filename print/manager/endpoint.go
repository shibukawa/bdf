package manager

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"

	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/form"
	"github.com/shibukawa/yuruidb"
)

// An endpoint is an entry of the network for applications, over REST:
// POST /in/<name> with the endpoint's key takes a document (PDF, Office,
// images, .bdf …), or the id of one of the endpoint's templates and the
// data to fill it with (JSON), and the job flows from the endpoint to its
// nodes. What it takes is its Accept and Formats (PutNode). Each endpoint
// has a key of its own, kept as a hash (argon2id); the key is shown once,
// when the endpoint is made or its key rotated. Its templates are its
// own: they go with it.

// InPath is where the endpoints are.
const InPath = "/in"

// OAuthTokenPath is where the apps take their tokens (apps.go).
const OAuthTokenPath = "/oauth/token"

// newKey makes an endpoint's key: 32 random bytes, in base64url.
func newKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "bdfk_" + base64.RawURLEncoding.EncodeToString(b)
}

// Parameters of argon2id (OWASP's first recommendation).
const (
	argonTime    = 2
	argonMemory  = 19 * 1024
	argonThreads = 1
	argonLen     = 32
)

// hashKey hashes a key as a PHC string of argon2id.
func hashKey(key string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	sum := argon2.IDKey([]byte(key), salt, argonTime, argonMemory, argonThreads, argonLen)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(sum)), nil
}

// checkKey tells whether a key is the one a PHC string hashes.
func checkKey(phc, key string) bool {
	parts := strings.Split(phc, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false
	}
	var memory, t uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &t, &threads); err != nil {
		return false
	}
	salt, err1 := base64.RawStdEncoding.DecodeString(parts[4])
	want, err2 := base64.RawStdEncoding.DecodeString(parts[5])
	if err1 != nil || err2 != nil || len(want) == 0 || memory == 0 || t == 0 || threads == 0 {
		return false
	}
	got := argon2.IDKey([]byte(key), salt, t, memory, threads, uint32(len(want)))
	return subtle.ConstantTimeCompare(got, want) == 1
}

// endpoint returns an endpoint node.
func (m *Manager) endpoint(ctx context.Context, name string) (*Node, error) {
	return m.nodeOfKind(ctx, name, KindEndpoint)
}

// Endpoint finds the endpoint a key opens.
func (m *Manager) Endpoint(ctx context.Context, name, key string) (*Node, error) {
	n, err := m.endpoint(ctx, name)
	if err != nil {
		return nil, err
	}
	if key == "" || n.KeyHash == "" || !checkKey(n.KeyHash, key) {
		return nil, fmt.Errorf("%w: the endpoint's key", ErrUnauthorized)
	}
	return n, nil
}

// MaxTemplate bounds a template's file.
const MaxTemplate = 32 << 20

// templateIDRE is the shape of a template's id: a name without a dot, so
// that the record's id (<endpoint>.<id>) splits at its last dot.
var templateIDRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)

// templateID is the record id (and the file's) of a template of an
// endpoint.
func templateID(endpoint, id string) string { return endpoint + "." + id }

// PutTemplate keeps a template of an endpoint under an id, replacing one
// there. name is its file name, which tells the format of a plain file
// (a compiled template tells its own).
func (m *Manager) PutTemplate(ctx context.Context, endpoint, id, name string, r io.ReaderAt, size int64) (*Template, error) {
	if _, err := m.endpoint(ctx, endpoint); err != nil {
		return nil, err
	}
	if !templateIDRE.MatchString(id) {
		return nil, fmt.Errorf("%w: template id %q: letters, digits, _ and -", ErrInvalid, id)
	}
	if size > MaxTemplate {
		return nil, fmt.Errorf("%w: %d bytes; a template is at most %d", ErrInvalid, size, MaxTemplate)
	}
	format, sample, err := form.Describe(r, size)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFormat, err)
	}
	fields, err := form.Fields(r, size)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFormat, err)
	}
	h := sha256.New()
	if _, err := io.Copy(h, io.NewSectionReader(r, 0, size)); err != nil {
		return nil, err
	}
	rid := templateID(endpoint, id)
	// the file is replaced: removed, then written
	if err := m.st.tpls.Remove(ctx, rid); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return nil, mapErr(err)
	}
	meta, err := m.st.tpls.Insert(ctx, rid, io.NewSectionReader(r, 0, size))
	if err != nil {
		return nil, mapErr(err)
	}
	t := Template{ID: id, Endpoint: endpoint, Name: clip(name, 255), Format: format, Size: meta.Size, SHA256: hex.EncodeToString(h.Sum(nil)), Uploaded: m.now(), Fields: fields, Sample: sample}
	if _, err := change(ctx, m.st.db, rid, true, func(rec *templateRecord) error {
		*rec = templateRecord(t)
		return nil
	}); err != nil {
		return nil, err
	}
	m.o.Log.Info("template", "endpoint", endpoint, "id", id, "format", format, "fields", len(fields))
	return &t, nil
}

// Template returns a template of an endpoint.
func (m *Manager) Template(ctx context.Context, endpoint, id string) (*Template, error) {
	if !validName(endpoint) || !templateIDRE.MatchString(id) {
		return nil, fmt.Errorf("%w: template %q of %q", ErrNotFound, id, endpoint)
	}
	rec, err := m.st.db.Get[templateRecord](ctx, templateID(endpoint, id))
	if err != nil {
		return nil, mapErr(err)
	}
	t := Template(rec.Data)
	t.ID, t.Endpoint = id, endpoint
	return &t, nil
}

// Templates returns the templates of an endpoint, by id.
func (m *Manager) Templates(ctx context.Context, endpoint string) ([]Template, error) {
	if _, err := m.endpoint(ctx, endpoint); err != nil {
		return nil, err
	}
	return m.templatesOf(ctx, endpoint)
}

func (m *Manager) templatesOf(ctx context.Context, endpoint string) ([]Template, error) {
	recs, err := list[templateRecord](ctx, m.st.db)
	if err != nil {
		return nil, err
	}
	out := []Template{}
	for _, r := range recs {
		if r.Data.Endpoint == endpoint {
			out = append(out, Template(r.Data))
		}
	}
	slices.SortFunc(out, func(a, b Template) int { return strings.Compare(a.ID, b.ID) })
	return out, nil
}

// DeleteTemplate removes a template of an endpoint.
func (m *Manager) DeleteTemplate(ctx context.Context, endpoint, id string) error {
	if !validName(endpoint) || !templateIDRE.MatchString(id) {
		return fmt.Errorf("%w: template %q of %q", ErrNotFound, id, endpoint)
	}
	rid := templateID(endpoint, id)
	rec, err := m.st.db.Get[templateRecord](ctx, rid)
	if err != nil {
		return mapErr(err)
	}
	if err := m.st.db.Delete[templateRecord](ctx, rid, rec.Revision()); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return mapErr(err)
	}
	if err := m.st.tpls.Remove(ctx, rid); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return mapErr(err)
	}
	return nil
}

// removeTemplatesOf removes every template of an endpoint that is gone.
func (m *Manager) removeTemplatesOf(ctx context.Context, endpoint string) error {
	ts, err := m.templatesOf(ctx, endpoint)
	if err != nil {
		return err
	}
	var errs []error
	for _, t := range ts {
		if err := m.DeleteTemplate(ctx, endpoint, t.ID); err != nil && !errors.Is(err, ErrNotFound) {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// WriteTemplate writes a template's file.
func (m *Manager) WriteTemplate(ctx context.Context, endpoint, id string, w io.Writer) error {
	r, err := m.st.tpls.Open(ctx, templateID(endpoint, id))
	if err != nil {
		return mapErr(err)
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// Fill is a request to fill a template: its id, the data, and the job's
// name and user.
type Fill struct {
	Template string         `json:"template"`
	Data     map[string]any `json:"data"`
	Name     string         `json:"name,omitempty"`
	User     string         `json:"user,omitempty"`
	// Strict refuses data that lacks a name the template looks up.
	Strict bool `json:"strict,omitempty"`
}

// SubmitFill fills a template of an endpoint and sends the form from the
// endpoint: its way of printing from data. An endpoint that takes
// documents only refuses (ErrFormat).
func (m *Manager) SubmitFill(ctx context.Context, ep *Node, by *Principal, f Fill) (*Job, error) {
	if ep.Accept == AcceptFile {
		return nil, fmt.Errorf("%w: %s takes documents, not templates to fill", ErrFormat, ep.Name)
	}
	t, err := m.Template(ctx, ep.Name, f.Template)
	if err != nil {
		return nil, err
	}
	r, err := m.st.tpls.Open(ctx, templateID(ep.Name, t.ID))
	if err != nil {
		return nil, mapErr(err)
	}
	b, err := io.ReadAll(io.LimitReader(r, MaxTemplate+1))
	r.Close()
	if err != nil {
		return nil, err
	}
	data := f.Data
	if data == nil {
		data = map[string]any{}
	}
	res, err := form.Fill(bytes.NewReader(b), int64(len(b)), &form.Options{Data: data, Strict: f.Strict,
		Convert: &converter.Options{FontDirs: m.o.FontDirs, NoSystemFonts: m.o.NoSystemFonts, FileName: t.Name, Params: map[string]string{"remote": "false"}}})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalid, err)
	}
	for _, w := range res.Warnings {
		m.o.Log.Warn("fill", "endpoint", ep.Name, "template", t.ID, "warning", w)
	}
	name := firstOr(f.Name, t.ID)
	return m.submitDoc(ctx, ep.Name, res.Doc, "form:"+t.Format, Input{Name: name, By: by, ClaimedUser: f.User, Source: "endpoint"})
}

// serveIn takes a request to an endpoint: a JSON request to fill a
// template, or a document. What the endpoint refuses (ErrFormat) is
// recorded in the audit log.
func (h *handler) serveIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.PathValue("name")
	var ep *Node
	var by *Principal
	var app *Viewer // an app sends with its token: the endpoint must list it
	if v, ok := h.viewerOf(r); ok && (v.Admin || v.Kind == PrincipalApp) && r.Header.Get("Authorization") != "" {
		by = v.Principal()
		n, err := h.m.endpoint(ctx, name)
		if err != nil {
			h.writeErr(w, ErrNotFound)
			return
		}
		ep = n
		if !v.Admin {
			app = v
		}
	} else {
		n, err := h.m.Endpoint(ctx, name, bearer(r))
		if err != nil {
			if errors.Is(err, ErrUnauthorized) {
				w.Header().Set("WWW-Authenticate", `Bearer realm="bdf print manager endpoint"`)
			}
			h.writeErr(w, err)
			return
		}
		ep = n
		by = &Principal{Kind: PrincipalEndpoint, ID: name, Via: ViaEndpointKey}
	}
	refused := func(err error) {
		if errors.Is(err, ErrFormat) {
			h.m.audit(ctx, by.Label(), "refuse-document", ep.Name, err.Error())
		}
		h.writeErr(w, err)
	}
	user := r.URL.Query().Get("user")
	f, err := os.CreateTemp(h.m.o.TempDir, "bdf-pm-in-*")
	if err != nil {
		h.writeErr(w, err)
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()
	n, err := io.Copy(f, io.LimitReader(r.Body, h.m.o.MaxDocument+1))
	if err != nil {
		h.writeErr(w, errors.Join(ErrInvalid, err))
		return
	}
	if n > h.m.o.MaxDocument {
		writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "the document is too large"})
		return
	}
	// a request to fill a template is JSON: said so, or an object
	ct := r.Header.Get("Content-Type")
	head := make([]byte, 1)
	f.ReadAt(head, 0)
	if strings.HasPrefix(ct, "application/json") || (ct == "" && head[0] == '{') {
		var fill Fill
		dec := json.NewDecoder(io.NewSectionReader(f, 0, n))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&fill); err != nil {
			h.writeErr(w, errors.Join(ErrInvalid, err))
			return
		}
		if fill.Template == "" {
			h.writeErr(w, fmt.Errorf("%w: a JSON request names a template", ErrInvalid))
			return
		}
		fill.User = firstOr(fill.User, user)
		if err := h.m.appMaySend(ctx, app, ep.Name, true); err != nil {
			h.writeErr(w, err)
			return
		}
		j, err := h.m.SubmitFill(ctx, ep, by, fill)
		if err != nil {
			refused(err)
			return
		}
		writeJSON(w, http.StatusCreated, j)
		return
	}
	if err := h.m.appMaySend(ctx, app, ep.Name, false); err != nil {
		h.writeErr(w, err)
		return
	}
	j, err := h.m.Submit(ctx, ep.Name, Input{Name: r.URL.Query().Get("name"), By: by, ClaimedUser: user, Source: "endpoint", Doc: f, Size: n})
	if err != nil {
		refused(err)
		return
	}
	writeJSON(w, http.StatusCreated, j)
}

// The templates' API: under their endpoint.

func (h *handler) listTemplates(w http.ResponseWriter, r *http.Request) {
	ts, err := h.m.Templates(r.Context(), r.PathValue("name"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ts)
}

func (h *handler) getTemplate(w http.ResponseWriter, r *http.Request) {
	t, err := h.m.Template(r.Context(), r.PathValue("name"), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *handler) templateFile(w http.ResponseWriter, r *http.Request) {
	t, err := h.m.Template(r.Context(), r.PathValue("name"), r.PathValue("id"))
	if err != nil {
		h.writeErr(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(t.Size, 10))
	w.Header().Set("Content-Disposition", `attachment; filename="`+firstOr(t.Name, t.ID)+`"`)
	if err := h.m.WriteTemplate(r.Context(), t.Endpoint, t.ID, w); err != nil {
		panic(http.ErrAbortHandler)
	}
}

// putTemplate takes a template's file as the body: ?name= is its file
// name.
func (h *handler) putTemplate(w http.ResponseWriter, r *http.Request) {
	f, err := os.CreateTemp(h.m.o.TempDir, "bdf-pm-tpl-*")
	if err != nil {
		h.writeErr(w, err)
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()
	n, err := io.Copy(f, io.LimitReader(r.Body, MaxTemplate+1))
	if err != nil {
		h.writeErr(w, errors.Join(ErrInvalid, err))
		return
	}
	t, err := h.m.PutTemplate(r.Context(), r.PathValue("name"), r.PathValue("id"), r.URL.Query().Get("name"), f, n)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), h.who(r), "put-template", t.Endpoint+"."+t.ID, t.Name)
	writeJSON(w, http.StatusOK, t)
}

func (h *handler) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	if err := h.m.DeleteTemplate(r.Context(), r.PathValue("name"), r.PathValue("id")); err != nil {
		h.writeErr(w, err)
		return
	}
	h.m.audit(r.Context(), h.who(r), "delete-template", r.PathValue("name")+"."+r.PathValue("id"), "")
	w.WriteHeader(http.StatusNoContent)
}

// formats answers the formats an endpoint may take.
func (h *handler) formats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, KnownFormats())
}

// The layout of the topology's canvas.

const layoutID = "admin"

// GetLayout returns where the nodes are drawn.
func (m *Manager) GetLayout(ctx context.Context) (*Layout, error) {
	rec, err := m.st.db.Get[layoutRecord](ctx, layoutID)
	if errors.Is(err, yuruidb.ErrNotFound) {
		return &Layout{Positions: map[string]Point{}}, nil
	}
	if err != nil {
		return nil, mapErr(err)
	}
	l := Layout(rec.Data)
	if l.Positions == nil {
		l.Positions = map[string]Point{}
	}
	return &l, nil
}

// PutLayout moves nodes: the positions given replace those of the same
// names; a position of a node that is gone is dropped.
func (m *Manager) PutLayout(ctx context.Context, l Layout) (*Layout, error) {
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return nil, err
	}
	names := map[string]bool{}
	for _, n := range nodes {
		names[n.Name] = true
	}
	rec, err := change(ctx, m.st.db, layoutID, true, func(r *layoutRecord) error {
		if r.Positions == nil {
			r.Positions = map[string]Point{}
		}
		for name, p := range l.Positions {
			r.Positions[name] = p
		}
		for name := range r.Positions {
			if !names[name] {
				delete(r.Positions, name)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	out := Layout(rec)
	return &out, nil
}

func (h *handler) getLayout(w http.ResponseWriter, r *http.Request) {
	l, err := h.m.GetLayout(r.Context())
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, l)
}

func (h *handler) putLayout(w http.ResponseWriter, r *http.Request) {
	var l Layout
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&l); err != nil {
		h.writeErr(w, errors.Join(ErrInvalid, err))
		return
	}
	got, err := h.m.PutLayout(r.Context(), l)
	if err != nil {
		h.writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, got)
}
