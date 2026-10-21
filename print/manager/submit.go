package manager

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image/png"
	"io"
	"os"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// Input is a document sent to a destination.
type Input struct {
	// Name is the job's name (job-name, or the file's).
	Name string
	// By is who sent it, as verified; nil is an anonymous sender.
	By *Principal
	// ClaimedUser is who the sender says sent it (requesting-user-name).
	ClaimedUser string
	// Source is "ipp" or "api".
	Source string
	// Doc is the document, of Size bytes: .bdf, or what the converter
	// reads.
	Doc  io.ReaderAt
	Size int64
}

// ErrFormat is the error of a document neither .bdf nor of a format the
// converter reads.
var ErrFormat = errors.New("manager: the document's format is not known")

// Submit takes a document for a destination: it is made a .bdf document,
// kept once as the job's, and the job is routed from the destination
// through the topology (route). Deliveries to printers and archives are
// sent after (Tick, or the background).
func (m *Manager) Submit(ctx context.Context, dest string, in Input) (*Job, error) {
	if _, err := m.kindOf(ctx, dest); err != nil {
		return nil, err
	}
	doc, format, err := m.document(ctx, in)
	if err != nil {
		return nil, err
	}
	return m.submitDoc(ctx, dest, doc, format, in)
}

// submitDoc keeps a document as a job and routes it from a destination.
func (m *Manager) submitDoc(ctx context.Context, dest string, doc *bdf.Document, format string, in Input) (*Job, error) {
	if _, err := m.kindOf(ctx, dest); err != nil {
		return nil, err
	}
	v, err := print.PrintView(doc, "")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFormat, err)
	}
	f, err := os.CreateTemp(m.o.TempDir, "bdf-pm-*.bdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	h := sha256.New()
	if err := doc.WriteSingle(io.MultiWriter(f, h)); err != nil {
		return nil, err
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
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
	j := &Job{ID: jobID(number), Number: number, Dest: dest, Name: clip(in.Name, 255), SubmittedBy: by, ClaimedUser: clip(in.ClaimedUser, 255),
		Group: group, Source: in.Source, Format: format, Pages: len(v.Pages), SHA256: hex.EncodeToString(h.Sum(nil)), Received: now,
		State: JobProcessing, Legs: []Leg{}}
	meta, err := m.st.docs.Insert(ctx, j.ID, f)
	if err != nil {
		return nil, mapErr(err)
	}
	j.Size = meta.Size
	// the job is there before its slips, items and deliveries refer to it
	if _, err := m.st.db.Insert(ctx, j.ID, jobRecord(*j)); err != nil {
		m.st.docs.Remove(context.WithoutCancel(ctx), j.ID)
		return nil, mapErr(err)
	}
	r := &routing{job: j, now: now}
	m.route(ctx, r, dest, "", 0)
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
	m.o.Log.Info("job", "dest", dest, "job", number, "name", j.Name, "pages", j.Pages, "source", j.Source, "state", j.State, "legs", len(j.Legs),
		"by", by.Kind+":"+by.ID, "via", by.Via)
	m.kick(r.deliveries)
	return j, nil
}

// document reads or converts what a sender sent.
func (m *Manager) document(ctx context.Context, in Input) (*bdf.Document, string, error) {
	if in.Size > m.o.MaxDocument {
		return nil, "", fmt.Errorf("%w: %d bytes; at most %d", ErrInvalid, in.Size, m.o.MaxDocument)
	}
	head := make([]byte, len(bdf.Magic))
	if n, _ := in.Doc.ReadAt(head, 0); n == len(head) && bytes.Equal(head, bdf.Magic[:]) {
		r, err := bdf.OpenSingle(in.Doc, in.Size)
		if err != nil {
			return nil, "", fmt.Errorf("%w: %v", ErrFormat, err)
		}
		defer r.Close()
		if r.Locked() {
			return nil, "", fmt.Errorf("%w: the document is locked with a password", ErrFormat)
		}
		doc, err := r.ToDocument()
		if err != nil {
			return nil, "", fmt.Errorf("%w: %v", ErrFormat, err)
		}
		return doc, "bdf", nil
	}
	if m.o.Converter == nil {
		return nil, "", ErrFormat
	}
	format, ok := m.o.Converter.Detect(in.Doc, in.Size, in.Name)
	if !ok {
		return nil, "", ErrFormat
	}
	doc, _, err := m.o.Converter.Convert(ctx, in.Doc, in.Size, format, in.Name, "")
	if err != nil {
		return nil, "", fmt.Errorf("%w: %s: %v", ErrFormat, format, err)
	}
	if doc.Lock != nil {
		return nil, "", fmt.Errorf("%w: the document needs a password", ErrFormat)
	}
	return doc, format, nil
}

// Document returns the document of a job.
func (m *Manager) Document(ctx context.Context, job int32) (*bdf.Document, error) {
	r, err := m.st.docs.Open(ctx, jobID(job))
	if err != nil {
		return nil, mapErr(err)
	}
	defer r.Close()
	br, err := bdf.OpenSingle(r, r.Size)
	if err != nil {
		return nil, err
	}
	defer br.Close()
	return br.ToDocument()
}

// WriteDocument writes the .bdf file of a job.
func (m *Manager) WriteDocument(ctx context.Context, job int32, w io.Writer) error {
	r, err := m.st.docs.Open(ctx, jobID(job))
	if err != nil {
		return mapErr(err)
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// MaxPageWidth bounds the width of a page drawn for a device.
const MaxPageWidth = 2400

// WritePage draws a page (from 0) of a slip as PNG, width pixels wide.
func (m *Manager) WritePage(ctx context.Context, s *Slip, page, width int, w io.Writer) error {
	if page < 0 || page >= s.Pages {
		return fmt.Errorf("%w: page %d of %d", ErrNotFound, page+1, s.Pages)
	}
	width = min(max(width, 200), MaxPageWidth)
	doc, err := m.Document(ctx, s.Job)
	if err != nil {
		return err
	}
	v, err := print.PrintView(doc, "")
	if err != nil {
		return err
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
