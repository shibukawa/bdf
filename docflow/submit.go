package docflow

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"io"
	"mime"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// Input is a document, or data, sent to a destination.
type Input struct {
	// Name is the job's name (job-name, or the file's).
	Name string
	// By is who sent it, as verified; nil is an anonymous sender.
	By *Principal
	// ClaimedUser is who the sender says sent it (requesting-user-name).
	ClaimedUser string
	// Source is "ipp", "api", "endpoint" or "reprint".
	Source string
	// Doc is what was sent, of Size bytes: .bdf, data, or what the
	// converter reads.
	Doc  io.ReaderAt
	Size int64
	// ContentType is the request's, when the sender gave one: it tells
	// data apart (application/json, text/csv …) where the content and
	// the name do not.
	ContentType string
	// Meta is the metadata the sender gave (meta.go).
	Meta map[string]any
}

// ErrFormat is the error of a document neither .bdf nor of a format the
// converter reads, or one the destination does not take.
var ErrFormat = errors.New("manager: the document's format is not known")

// Submit takes a document or data for a destination. An entry's settings
// say what it takes (Accept, Formats) and what it sends on of a document
// (Output: the .bdf made of it, or the file as received); a template node
// takes data, a convert node a file as it is, every other node a document
// made .bdf. What was received is kept once as the job's first payload,
// and the job is routed from the destination through the topology
// (route). What the nodes send and make (deliveries, the steps of
// template and convert nodes) is done after (Tick, or the background).
// A job the asserts right after the destination refuse is kept, aborted,
// and the error is a *Refusal (gate).
func (m *Manager) Submit(ctx context.Context, dest string, in Input) (*Job, error) {
	n, err := m.Node(ctx, dest)
	if err != nil {
		return nil, err
	}
	accept, output, formats := AcceptFile, OutputBDF, []string(nil)
	var entry *Node
	switch {
	case isEntry(n.Kind):
		accept, output, formats, entry = or(n.Accept, AcceptFile), or(n.Output, OutputBDF), n.Formats, n
	case n.Kind == KindTemplate:
		accept = AcceptData
	case n.Kind == KindConvert:
		output = OutputFile
	case n.Kind == KindExport:
		// an export lists what an archive wrote: a job sent to it has none
		return nil, fmt.Errorf("%w: %s is an export, which lists what an archive before it wrote: send to an entry", ErrInvalid, dest)
	}
	claimed, err := claimedMeta(in.Meta)
	if err != nil {
		return nil, err
	}
	rc, err := m.receive(ctx, in, accept, output, formats)
	if err != nil {
		return nil, err
	}
	defer rc.close()
	return m.submitPayload(ctx, dest, entry, rc, in, claimed)
}

// received is what an entry made of what a sender sent: the bytes to
// keep as the job's first payload, and what they are.
type received struct {
	r    io.ReaderAt
	size int64
	p    Payload
	// format is the format received (Job.Format): the file's, before an
	// entry converted it.
	format string
	tmp    *os.File
}

func (rc *received) close() {
	if rc.tmp != nil {
		rc.tmp.Close()
		os.Remove(rc.tmp.Name())
	}
}

// receive reads what a sender sent: a .bdf document, data, or a document
// the converter reads, as the destination takes it (accept, formats),
// converted to .bdf when the destination sends that on (output).
func (m *Manager) receive(ctx context.Context, in Input, accept, output string, formats []string) (*received, error) {
	if in.Size > m.o.MaxDocument {
		return nil, fmt.Errorf("%w: %d bytes; at most %d", ErrInvalid, in.Size, m.o.MaxDocument)
	}
	format := "bdf"
	head := make([]byte, len(bdf.Magic))
	if n, _ := in.Doc.ReadAt(head, 0); n != len(head) || !bytes.Equal(head, bdf.Magic[:]) {
		if format = m.detect(in); format == "" {
			return nil, ErrFormat
		}
	}
	kind := kindOfFormat(format)
	if accept == AcceptData && kind != PayloadData {
		return nil, fmt.Errorf("%w: %s is not data (json, yaml, csv, xlsx or parquet)", ErrFormat, format)
	}
	if !accepted(formats, format) {
		return nil, fmt.Errorf("%w: %s is not taken here (%s)", ErrFormat, format, strings.Join(formats, ", "))
	}
	rc := &received{r: in.Doc, size: in.Size, format: format, p: Payload{Kind: kind, Format: format, Name: clip(in.Name, 255)}}
	switch {
	case format == "bdf":
		doc, err := readBDF(in.Doc, in.Size)
		if err != nil {
			return nil, err
		}
		if rc.p.Pages, err = pagesOf(doc); err != nil {
			return nil, err
		}
	case accept == AcceptData || output == OutputFile:
		// as received
	default:
		doc, err := m.convert(ctx, in.Doc, in.Size, format, in.Name)
		if err != nil {
			return nil, err
		}
		pages, err := pagesOf(doc)
		if err != nil {
			return nil, err
		}
		f, size, err := m.writeDocument(doc)
		if err != nil {
			return nil, err
		}
		rc.tmp, rc.r, rc.size = f, f, size
		rc.p.Kind, rc.p.Format, rc.p.Pages, rc.p.Name = PayloadBDF, "bdf", pages, bdfName(rc.p.Name)
	}
	return rc, nil
}

// submitPayload keeps what was received as a job and routes it from a
// destination.
func (m *Manager) submitPayload(ctx context.Context, dest string, entry *Node, rc *received, in Input, claimed map[string]any) (*Job, error) {
	number, err := m.nextJob(ctx)
	if err != nil {
		return nil, err
	}
	now := m.now()
	by := in.By
	if by == nil {
		by = &Principal{Kind: PrincipalAnonymous, Via: ViaNone}
	}
	_, _, group, err := m.destScope(ctx, dest)
	if err != nil {
		return nil, err
	}
	if group == "" && by.Kind == PrincipalApp {
		// an app's jobs are its group's, when the destination gives none
		if a, err := m.App(ctx, by.ID); err == nil {
			group = a.Group
		}
	}
	j := &Job{ID: jobID(number), Number: number, Dest: dest, Name: clip(in.Name, 255), SubmittedBy: by, ClaimedUser: clip(in.ClaimedUser, 255),
		Group: group, Source: in.Source, Kind: rc.p.Kind, Format: rc.format, Pages: rc.p.Pages, Received: now, State: JobProcessing, Legs: []Leg{}}
	h := sha256.New()
	meta, err := m.st.docs.Insert(ctx, j.ID, io.TeeReader(io.NewSectionReader(rc.r, 0, rc.size), h))
	if err != nil {
		return nil, mapErr(err)
	}
	p := rc.p
	p.ID, p.Size, p.SHA256, p.At = j.ID, meta.Size, hex.EncodeToString(h.Sum(nil)), now
	j.Payloads = []Payload{p}
	j.Size, j.SHA256 = p.Size, p.SHA256
	j.Meta = jobMeta(claimed, entry, j)
	// the job is there before its slips, items and deliveries refer to it
	if _, err := m.st.db.Insert(ctx, j.ID, jobRecord(*j)); err != nil {
		m.st.docs.Remove(context.WithoutCancel(ctx), j.ID)
		return nil, mapErr(err)
	}
	r := &routing{job: j, now: now, payload: &j.Payloads[0]}
	refused := m.gate(ctx, r, dest)
	if refused == nil {
		m.route(ctx, r, dest, "", 0)
	}
	j, err = m.changeJob(ctx, j.ID, func(j *Job) error {
		j.Legs = append(j.Legs, r.legs...)
		j.Live += r.live
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := m.touchDest(ctx, dest, j.ID); err != nil {
		return j, err
	}
	if err := m.touchDest(ctx, allJobs, j.ID); err != nil {
		return j, err
	}
	m.o.Log.Info("job", "dest", dest, "job", number, "name", j.Name, "kind", j.Kind, "format", j.Format, "pages", j.Pages, "source", j.Source, "state", j.State, "legs", len(j.Legs),
		"by", by.Kind+":"+by.ID, "via", by.Via)
	m.kick(r.deliveries)
	if refused != nil {
		refused.Job = j
		return j, refused
	}
	return j, nil
}

// accepted tells a format one of the formats accepted (none: every one).
func accepted(formats []string, format string) bool {
	return len(formats) == 0 || slices.Contains(formats, format)
}

// maxSniff bounds what is read to tell JSON from the rest.
const maxSniff = 8 << 20

// detect names the format of what a sender sent that is not .bdf: data by
// its name or the request's content type (JSON, YAML, CSV), JSON by its
// content, else what the converter reads; "" is unknown.
func (m *Manager) detect(in Input) string {
	ext := strings.ToLower(path.Ext(strings.ReplaceAll(in.Name, "\\", "/")))
	ct, _, _ := mime.ParseMediaType(in.ContentType)
	switch {
	case ext == ".json" || ct == "application/json":
		return "json"
	case ext == ".yaml" || ext == ".yml" || ct == "application/yaml" || ct == "application/x-yaml" || ct == "text/yaml":
		return "yaml"
	case ext == ".csv" || ext == ".tsv" || ct == "text/csv" || ct == "text/tab-separated-values":
		return "csv"
	}
	if looksJSON(in.Doc, in.Size) {
		return "json"
	}
	if m.o.Converter != nil {
		if format, ok := m.o.Converter.Detect(in.Doc, in.Size, in.Name); ok {
			return format
		}
	}
	return ""
}

// looksJSON tells a JSON object or array (a bare number or string is not
// taken as data): one of maxSniff bytes or less is checked whole.
func looksJSON(r io.ReaderAt, size int64) bool {
	head := make([]byte, 64)
	n, _ := r.ReadAt(head, 0)
	t := bytes.TrimLeft(head[:n], " \t\r\n")
	if len(t) == 0 || (t[0] != '{' && t[0] != '[') {
		return false
	}
	if size > maxSniff {
		return true
	}
	b := make([]byte, size)
	if n, err := r.ReadAt(b, 0); int64(n) != size && err != nil && err != io.EOF {
		return false
	}
	return json.Valid(b)
}

// readBDF opens a .bdf document that is not locked.
func readBDF(r io.ReaderAt, size int64) (*bdf.Document, error) {
	br, err := bdf.OpenSingle(r, size)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFormat, err)
	}
	defer br.Close()
	if br.Locked() {
		return nil, fmt.Errorf("%w: the document is locked with a password", ErrFormat)
	}
	doc, err := br.ToDocument()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFormat, err)
	}
	return doc, nil
}

// pagesView is the view of a document that is printed and shown: the
// first view of pages (fixed or flow), or the first scroll view, whose
// strips are pages too; as print.PrintView picks it.
func pagesView(doc *bdf.Document) (*bdf.View, error) {
	sheets := false
	for _, kinds := range [][]string{{bdf.ViewFixed, bdf.ViewFlow}, {bdf.ViewScroll}} {
		for _, v := range doc.Views {
			sheets = sheets || v.Kind == bdf.ViewSheet
			if slices.Contains(kinds, v.Kind) && len(v.Pages) > 0 {
				return v, nil
			}
		}
	}
	if sheets {
		return nil, errors.New("the document has no pages to print: its views are sheets, which have no pages (an Excel workbook has its print view with -param views=print)")
	}
	return nil, errors.New("the document has no pages to print")
}

// pagesOf counts the pages of a document as it is printed.
func pagesOf(doc *bdf.Document) (int, error) {
	v, err := pagesView(doc)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrFormat, err)
	}
	return len(v.Pages), nil
}

// convert makes a .bdf document of a file of a format the converter
// reads.
func (m *Manager) convert(ctx context.Context, r io.ReaderAt, size int64, format, name string) (*bdf.Document, error) {
	if m.o.Converter == nil {
		return nil, fmt.Errorf("%w: no converter", ErrFormat)
	}
	doc, _, err := m.o.Converter.Convert(ctx, r, size, format, name, "")
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrFormat, format, err)
	}
	if doc.Lock != nil {
		return nil, fmt.Errorf("%w: the document needs a password", ErrFormat)
	}
	return doc, nil
}

// writeDocument writes a document to a file of its own, rewound.
func (m *Manager) writeDocument(doc *bdf.Document) (*os.File, int64, error) {
	f, err := os.CreateTemp(m.o.TempDir, "bdf-pm-*.bdf")
	if err != nil {
		return nil, 0, err
	}
	var size int64
	err = doc.WriteSingle(f)
	if err == nil {
		size, err = f.Seek(0, io.SeekCurrent)
	}
	if err == nil {
		_, err = f.Seek(0, io.SeekStart)
	}
	if err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, 0, err
	}
	return f, size, nil
}

// bdfName names the .bdf document made of a file: the file's name with
// the extension changed.
func bdfName(name string) string {
	if name == "" {
		return ""
	}
	return strings.TrimSuffix(name, path.Ext(name)) + ".bdf"
}

// contentTypeOf is the media type of a payload, for a download: its file's
// extension's, else its format's.
func contentTypeOf(format, name string) string {
	if ext := path.Ext(name); ext != "" && format != "bdf" {
		if t := mime.TypeByExtension(ext); t != "" {
			return t
		}
	}
	switch format {
	case "bdf":
		return "application/x-bdf"
	case "pdf":
		return "application/pdf"
	case "json":
		return "application/json"
	case "yaml":
		return "application/yaml"
	case "csv":
		return "text/csv"
	case "xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "parquet":
		return "application/vnd.apache.parquet"
	case "png":
		return "image/png"
	}
	return "application/octet-stream"
}

// Reprint sends what a job received again, as a new job of a principal,
// to a destination: a reprint from the pages. The sender's metadata goes
// with it.
func (m *Manager) Reprint(ctx context.Context, j *Job, dest string, by *Principal) (*Job, error) {
	p := j.received()
	if p == nil {
		return nil, fmt.Errorf("%w: job %s has no payload", ErrNotFound, j.ID)
	}
	r, err := m.st.docs.Open(ctx, p.ID)
	if err != nil {
		return nil, mapErr(err)
	}
	defer r.Close()
	return m.Submit(ctx, dest, Input{Name: firstOr(j.Name, "job-"+j.ID), By: by, ClaimedUser: j.ClaimedUser, Source: "reprint", Doc: r, Size: r.Size,
		ContentType: contentTypeOf(p.Format, p.Name), Meta: claimedOf(j.Meta)})
}

// Document returns the .bdf document of a payload.
func (m *Manager) Document(ctx context.Context, payload string) (*bdf.Document, error) {
	r, err := m.st.docs.Open(ctx, payload)
	if err != nil {
		return nil, mapErr(err)
	}
	defer r.Close()
	return readBDF(r, r.Size)
}

// WritePayload writes the bytes of a payload.
func (m *Manager) WritePayload(ctx context.Context, payload string, w io.Writer) error {
	r, err := m.st.docs.Open(ctx, payload)
	if err != nil {
		return mapErr(err)
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// MaxPageWidth bounds the width of a page drawn for a device.
const MaxPageWidth = 2400

// WritePage draws a page (from 0, of pages) of a .bdf payload as PNG,
// width pixels wide.
func (m *Manager) WritePage(ctx context.Context, payload string, pages, page, width int, w io.Writer) error {
	if page < 0 || page >= pages {
		return fmt.Errorf("%w: page %d of %d", ErrNotFound, page+1, pages)
	}
	width = min(max(width, 64), MaxPageWidth)
	doc, err := m.Document(ctx, payload)
	if err != nil {
		return err
	}
	v, err := pagesView(doc)
	if err != nil {
		return err
	}
	if page >= len(v.Pages) {
		return fmt.Errorf("%w: page %d of %d", ErrNotFound, page+1, len(v.Pages))
	}
	p := v.Pages[page]
	if p == nil || p.W <= 0 {
		return fmt.Errorf("%w: page %d has no size", ErrInvalid, page+1)
	}
	r := imagebdf.New(doc, &imagebdf.Options{FontDirs: m.o.FontDirs, NoSystemFonts: m.o.NoSystemFonts})
	img, err := r.Page(v, page, float64(width)/float64(p.W))
	if err != nil {
		return err
	}
	return (&png.Encoder{CompressionLevel: png.BestSpeed}).Encode(w, img)
}

// clip cuts a string to n bytes, at a rune.
func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	for n > 0 && s[n]&0xc0 == 0x80 {
		n--
	}
	return s[:n]
}
