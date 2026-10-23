package docflow

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/form"
	"github.com/shibukawa/yuruidb"
	"github.com/shibukawa/yuruidb/blob"
)

// A template node and a convert node make a .bdf document of the payload
// that reaches them: a template node fills one of its templates with
// data (JSON, YAML, CSV, an Excel workbook), a convert node converts a
// file (a PDF, an Office document, a CSV file …) as an entry does. The
// step is a delivery (deliver.go): done under a lease by the background
// of a process or the next tick, tried again after a failure a retry may
// mend, never after one it cannot (finalErr). The document made is a new
// payload of the job, kept under the delivery's id, and the job goes on
// from the node's To with it (continueWith). A template node's templates
// are its own: they go with it.

// finalErr is a step's failure that no retry mends: a template missing,
// data that is not an object, a format with no converter.
type finalErr struct{ err error }

func (e finalErr) Error() string { return e.err.Error() }
func (e finalErr) Unwrap() error { return e.err }

func final(format string, args ...any) error { return finalErr{fmt.Errorf(format, args...)} }

func isFinal(err error) bool {
	var f finalErr
	return errors.As(err, &f)
}

// MaxTemplate bounds a template's file.
const MaxTemplate = 32 << 20

// templateIDRE is the shape of a template's id: a name without a dot, so
// that the record's id (<node>.<id>) splits at its last dot.
var templateIDRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)

// templateID is the record id (and the file's) of a template of a node.
func templateID(node, id string) string { return node + "." + id }

// templateNode returns a template node.
func (m *Manager) templateNode(ctx context.Context, name string) (*Node, error) {
	return m.nodeOfKind(ctx, name, KindTemplate)
}

// PutTemplate keeps a template of a template node under an id, replacing
// one there. name is its file name, which tells the format of a plain
// file (a compiled template tells its own).
func (m *Manager) PutTemplate(ctx context.Context, node, id, name string, r io.ReaderAt, size int64) (*Template, error) {
	if _, err := m.templateNode(ctx, node); err != nil {
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
	rid := templateID(node, id)
	// the file is replaced: removed, then written
	if err := m.st.tpls.Remove(ctx, rid); err != nil && !errors.Is(err, yuruidb.ErrNotFound) {
		return nil, mapErr(err)
	}
	meta, err := m.st.tpls.Insert(ctx, rid, io.NewSectionReader(r, 0, size))
	if err != nil {
		return nil, mapErr(err)
	}
	t := Template{ID: id, Node: node, Name: clip(name, 255), Format: format, Size: meta.Size, SHA256: hex.EncodeToString(h.Sum(nil)), Uploaded: m.now(), Fields: fields, Sample: sample}
	if _, err := change(ctx, m.st.db, rid, true, func(rec *templateRecord) error {
		*rec = templateRecord(t)
		return nil
	}); err != nil {
		return nil, err
	}
	m.o.Log.Info("template", "node", node, "id", id, "format", format, "fields", len(fields))
	return &t, nil
}

// Template returns a template of a node.
func (m *Manager) Template(ctx context.Context, node, id string) (*Template, error) {
	if !validName(node) || !templateIDRE.MatchString(id) {
		return nil, fmt.Errorf("%w: template %q of %q", ErrNotFound, id, node)
	}
	rec, err := m.st.db.Get[templateRecord](ctx, templateID(node, id))
	if err != nil {
		return nil, mapErr(err)
	}
	t := Template(rec.Data)
	t.ID, t.Node = id, node
	return &t, nil
}

// Templates returns the templates of a template node, by id.
func (m *Manager) Templates(ctx context.Context, node string) ([]Template, error) {
	if _, err := m.templateNode(ctx, node); err != nil {
		return nil, err
	}
	return m.templatesOf(ctx, node)
}

func (m *Manager) templatesOf(ctx context.Context, node string) ([]Template, error) {
	recs, err := list[templateRecord](ctx, m.st.db)
	if err != nil {
		return nil, err
	}
	out := []Template{}
	for _, r := range recs {
		if r.Data.Node == node {
			out = append(out, Template(r.Data))
		}
	}
	slices.SortFunc(out, func(a, b Template) int { return strings.Compare(a.ID, b.ID) })
	return out, nil
}

// DeleteTemplate removes a template of a node.
func (m *Manager) DeleteTemplate(ctx context.Context, node, id string) error {
	if !validName(node) || !templateIDRE.MatchString(id) {
		return fmt.Errorf("%w: template %q of %q", ErrNotFound, id, node)
	}
	rid := templateID(node, id)
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

// removeTemplatesOf removes every template of a node that is gone, or is
// no template node any more.
func (m *Manager) removeTemplatesOf(ctx context.Context, node string) error {
	ts, err := m.templatesOf(ctx, node)
	if err != nil {
		return err
	}
	var errs []error
	for _, t := range ts {
		if err := m.DeleteTemplate(ctx, node, t.ID); err != nil && !errors.Is(err, ErrNotFound) {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// WriteTemplate writes a template's file.
func (m *Manager) WriteTemplate(ctx context.Context, node, id string, w io.Writer) error {
	r, err := m.st.tpls.Open(ctx, templateID(node, id))
	if err != nil {
		return mapErr(err)
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// readBlob reads a blob whole, of limit bytes at most.
func readBlob(ctx context.Context, st *blob.Store, id string, limit int64) ([]byte, error) {
	r, err := st.Open(ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	defer r.Close()
	if r.Size > limit {
		return nil, fmt.Errorf("%w: %s is %d bytes; at most %d", ErrInvalid, id, r.Size, limit)
	}
	return io.ReadAll(io.LimitReader(r, limit))
}

// transform does the step of a template or a convert node for a
// delivery: the payload it makes, kept under the delivery's id. A convert
// node sends a .bdf document on as it is.
func (m *Manager) transform(ctx context.Context, d *Delivery, n *Node, j *Job) (*Payload, error) {
	in := j.payload(d.Payload)
	if in == nil {
		return nil, final("the job has no payload %q", d.Payload)
	}
	var doc *bdf.Document
	var err error
	switch n.Kind {
	case KindTemplate:
		if in.Kind != PayloadData {
			return nil, final("%s takes data; this is %s", n.Name, in.Kind)
		}
		doc, err = m.fill(ctx, n, j, in)
	case KindConvert:
		if in.Kind == PayloadBDF {
			return in, nil
		}
		doc, err = m.convertPayload(ctx, n, in)
	case KindSign:
		return m.signPayload(ctx, d, n, j, in)
	default:
		return nil, final("%s is a %s: it transforms nothing", n.Name, n.Kind)
	}
	if err != nil {
		return nil, err
	}
	return m.keepDocument(ctx, d.ID, n.Name, bdfName(in.Name), doc)
}

// fill fills the template a template node names (its Template, or the
// one the metadata names at TemplateFrom) with the job's data.
func (m *Manager) fill(ctx context.Context, n *Node, j *Job, in *Payload) (*bdf.Document, error) {
	id := n.Template
	if id == "" && n.TemplateFrom == "" {
		return nil, final("%s names no template yet: choose one of its templates, or the metadata key that names one", n.Name)
	}
	if n.TemplateFrom != "" {
		if id = metaView(j, in, nil)[n.TemplateFrom]; id == "" {
			return nil, final("the metadata has no %s to name the template", n.TemplateFrom)
		}
	}
	t, err := m.Template(ctx, n.Name, id)
	if errors.Is(err, ErrNotFound) {
		return nil, final("%s has no template %q", n.Name, id)
	}
	if err != nil {
		return nil, err
	}
	tpl, err := readBlob(ctx, m.st.tpls, templateID(n.Name, t.ID), MaxTemplate)
	if err != nil {
		return nil, err
	}
	data, err := readBlob(ctx, m.st.docs, in.ID, m.o.MaxDocument)
	if err != nil {
		return nil, err
	}
	values, err := readData(in.Format, in.Name, data, n.ListName)
	if err != nil {
		return nil, final("the data: %v", err)
	}
	res, err := form.Fill(bytes.NewReader(tpl), int64(len(tpl)), &form.Options{Data: values, Strict: n.Strict,
		Convert: &converter.Options{FontDirs: m.o.FontDirs, NoSystemFonts: m.o.NoSystemFonts, FileName: t.Name, Params: map[string]string{"remote": "false"}}})
	if err != nil {
		return nil, final("template %s: %v", t.ID, err)
	}
	for _, w := range res.Warnings {
		m.o.Log.Warn("fill", "node", n.Name, "template", t.ID, "job", j.Number, "warning", w)
	}
	return res.Doc, nil
}

// readData reads data for a template: a JSON or YAML object is the data
// itself; a list (a CSV file, the sheets of a workbook, a JSON array) is
// the value of listName, or of the file's name without its extension.
func readData(format, name string, b []byte, listName string) (map[string]any, error) {
	ext, ok := map[string]string{"json": ".json", "yaml": ".yaml", "csv": ".csv", "xlsx": ".xlsx"}[format]
	if !ok {
		return nil, fmt.Errorf("%s data is not filled into a template: convert it to a document instead", format)
	}
	base := path.Base(strings.ReplaceAll(name, "\\", "/"))
	stem := strings.TrimSuffix(base, path.Ext(base))
	if stem == "" || stem == "." || stem == "/" {
		stem = "data"
	}
	file := stem + ext
	switch format {
	case "json", "yaml":
		data, err := form.LoadData(form.Source{FileName: file, Data: b})
		if err == nil {
			return data, nil
		}
		if named, err2 := form.LoadData(form.Source{Name: firstOr(listName, stem), FileName: file, Data: b}); err2 == nil {
			return named, nil
		}
		return nil, err
	}
	return form.LoadData(form.Source{Name: listName, FileName: file, Data: b})
}

// convertPayload makes a .bdf document of a payload at a convert node.
func (m *Manager) convertPayload(ctx context.Context, n *Node, in *Payload) (*bdf.Document, error) {
	if !accepted(n.Formats, in.Format) {
		return nil, final("%s is not converted at %s (%s)", in.Format, n.Name, strings.Join(n.Formats, ", "))
	}
	if in.Format == "json" || in.Format == "yaml" {
		return nil, final("there is no document of %s data: fill a template with it", in.Format)
	}
	r, err := m.st.docs.Open(ctx, in.ID)
	if err != nil {
		return nil, mapErr(err)
	}
	defer r.Close()
	doc, err := m.convert(ctx, r, r.Size, in.Format, in.Name)
	if err != nil {
		return nil, final("%v", err)
	}
	return doc, nil
}

// keepDocument keeps a document a node made as a payload under an id: one
// there already (a retry after the outcome of the write was not known) is
// taken as kept.
func (m *Manager) keepDocument(ctx context.Context, id, node, name string, doc *bdf.Document) (*Payload, error) {
	pages, err := pagesOf(doc)
	if err != nil {
		return nil, final("%v", err)
	}
	f, size, err := m.writeDocument(doc)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	h := sha256.New()
	meta, err := m.st.docs.Insert(ctx, id, io.TeeReader(f, h))
	switch {
	case errors.Is(err, yuruidb.ErrAlreadyExists):
		r, err := m.st.docs.Open(ctx, id)
		if err != nil {
			return nil, mapErr(err)
		}
		defer r.Close()
		h = sha256.New()
		if _, err := io.Copy(h, r); err != nil {
			return nil, err
		}
		size = r.Size
	case err != nil:
		return nil, mapErr(err)
	default:
		size = meta.Size
	}
	return &Payload{ID: id, Kind: PayloadBDF, Format: "bdf", Name: name, Pages: pages, Size: size, SHA256: hex.EncodeToString(h.Sum(nil)), Made: node, At: m.now()}, nil
}

// keepMade keeps what a step of an adapter's node made as a payload
// under an id, as keepDocument does.
func (m *Manager) keepMade(ctx context.Context, id, node string, made *Made) (*Payload, error) {
	if kindBit(made.Kind) == 0 {
		return nil, final("%s made a payload of no kind (%q)", node, made.Kind)
	}
	h := sha256.New()
	var size int64
	meta, err := m.st.docs.Insert(ctx, id, io.TeeReader(io.LimitReader(made.Body, m.o.MaxDocument+1), h))
	switch {
	case errors.Is(err, yuruidb.ErrAlreadyExists):
		r, err := m.st.docs.Open(ctx, id)
		if err != nil {
			return nil, mapErr(err)
		}
		defer r.Close()
		h = sha256.New()
		if _, err := io.Copy(h, r); err != nil {
			return nil, err
		}
		size = r.Size
	case err != nil:
		return nil, mapErr(err)
	default:
		size = meta.Size
	}
	if size > m.o.MaxDocument {
		m.st.docs.Remove(ctx, id)
		return nil, final("%s made %d bytes; a document is at most %d", node, size, m.o.MaxDocument)
	}
	return &Payload{ID: id, Kind: made.Kind, Format: made.Format, Name: made.Name, Pages: made.Pages, Size: size, SHA256: hex.EncodeToString(h.Sum(nil)), Made: node, At: m.now()}, nil
}

// continueWith sends a job on from a node with the payload the node
// made: the payload is added to the job, the node's leg says what it
// made, and the job follows the node's To with it and its Original with
// the payload the node took, each with the trace of the step.
func (m *Manager) continueWith(ctx context.Context, d *Delivery, n *Node, p *Payload) error {
	return m.continueFrom(ctx, d, n, p, true, nil)
}

// continueFrom routes the job on from a node whose step is done, with a
// payload: one the node made (made: the leg's Made, and the payload the
// node took goes along Original), or the one it took (an archive, which
// sends on what it wrote, archived, with the trace it came with).
func (m *Manager) continueFrom(ctx context.Context, d *Delivery, n *Node, p *Payload, made bool, archived *archiveRef) error {
	j, err := m.JobOfNumber(ctx, d.Job)
	if err != nil {
		return err
	}
	if j.Canceled != nil {
		return nil
	}
	if j.payload(p.ID) == nil {
		j.Payloads = append(j.Payloads, *p)
	}
	r := &routing{job: j, now: m.now(), payload: j.payload(p.ID), archived: archived, trace: d.Trace, cause: d.Cause}
	if made {
		in := j.payload(d.Payload)
		if in == nil {
			in = r.payload
		}
		r.trace = newTrace(d, n, j, in, r.payload, TraceOutput, d.Trace)
		for _, to := range n.To {
			m.route(ctx, r, to, n.Name, 0)
		}
		r.payload, r.trace = in, r.trace.as(TraceSource)
		for _, to := range n.Original {
			m.route(ctx, r, to, n.Name, 0)
		}
	} else {
		for _, to := range n.To {
			m.route(ctx, r, to, n.Name, 0)
		}
	}
	if _, err := m.changeJob(ctx, j.ID, func(j *Job) error {
		if j.payload(p.ID) == nil {
			j.Payloads = append(j.Payloads, *p)
		}
		if l := j.leg(d.ID); l != nil && made {
			l.Made = p.ID
		}
		j.Legs = append(j.Legs, r.legs...)
		j.Live += r.live
		return nil
	}); err != nil {
		return err
	}
	m.kick(r.deliveries)
	return nil
}

// The templates' API: under their template node.

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
	if err := h.m.WriteTemplate(r.Context(), t.Node, t.ID, w); err != nil {
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
	h.m.audit(r.Context(), h.who(r), "put-template", t.Node+"."+t.ID, t.Name)
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
