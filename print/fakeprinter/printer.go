// Package fakeprinter is printers that are not there, for tests and for
// trying bdf print and bdf printd without paper. A Printer takes jobs over
// IPP (it is an http.Handler that answers as an IPP Everywhere printer
// does), over the raw port 9100, over LPD and over a serial line, keeps
// what it is sent, and can be told what to support, what to refuse and
// how to fail: it stops, forgets its jobs, restarts its job ids, cuts a
// connection, answers slowly. A CUPS is a CUPS server of such printers,
// with classes. Presets are printers of the kinds bdf prints to; the
// package readback turns what they received back into pages.
//
// The command tools/fakeprinter runs a farm of them from a YAML file, with
// a page of what they received.
package fakeprinter

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// Printer is a printer that is not there. New gives one with the
// capabilities of a plain office printer; change the fields before it is
// served, and through the methods while it is.
type Printer struct {
	// Name is printer-name; Info, Location and MakeModel the printer's
	// description; UUID printer-uuid.
	Name, Info, Location, MakeModel string
	UUID                            string
	// Formats are document-format-supported.
	Formats []string
	// Reject are formats refused with client-error-document-format-
	// not-supported although listed, to test a client's fallback.
	Reject      []string
	Media       []string
	MediaReady  []string
	Sources     []string // media-source-supported
	Types       []string // media-type-supported
	Sides       []string
	ColorModes  []string
	Qualities   []int
	Resolutions []ipp.Resolution
	Copies      ipp.Range
	PageRanges  bool
	Finishings  []int
	Scaling     []string
	OutputBins  []string

	PWGRasterTypes       []string
	PWGRasterResolutions []ipp.Resolution
	SheetBack            string
	URF                  []string
	// PCLmStripHeight, when not 0, makes the printer say it takes PCLm,
	// with these strips, back side, resolutions and compressions.
	PCLmStripHeight int
	PCLmBackSide    string
	PCLmResolutions []ipp.Resolution
	PCLmCompression []string

	// Language is what the printer takes on its raw transports (9100,
	// LPD, serial): a language of print/pdl ("zpl", "escpos" …),
	// "postscript", "pdf", "pcl", or "" for anything. With Width and DPI
	// (its dots across and resolution) it tells readback how to read a
	// job, and it makes the printer answer the status queries of its
	// language (ZPL ~HS, ESC/POS DLE EOT).
	Language   string
	Width, DPI int
	// Profile is the print/pdl profile a printer of a language is of
	// ("zd421-203"), which a client names in its URL (?profile=).
	Profile string

	State     ipp.PrinterState
	Reasons   []string
	Accepting bool

	// Username and Password, when set, are required as Basic
	// authentication.
	Username, Password string
	// IgnoreUnsupported takes a job whose attributes the printer does not
	// support without ipp-attribute-fidelity, ignoring them
	// (successful-ok-ignored-or-substituted-attributes), as RFC 8011
	// says; by default the printer refuses it, as many printers do.
	IgnoreUnsupported bool
	// Only11 refuses IPP/2.x with server-error-version-not-supported.
	Only11 bool
	// NoChunked refuses requests without a length (chunked), with HTTP
	// 411, as old printers do.
	NoChunked bool
	// Processing is how long a job stays processing after its last
	// document before it is completed.
	Processing time.Duration
	// Busy refuses so many of the next Create-Job and Print-Job with
	// server-error-busy, as a printer that is printing for others does.
	Busy int
	// Latency delays each IPP answer and each raw connection's first
	// byte; Bandwidth (bytes a second, 0 for no bound) slows what the
	// printer reads, as a slow printer's buffer does.
	Latency   time.Duration
	Bandwidth int
	// NeverComplete keeps every job processing.
	NeverComplete bool
	// Dir saves each document to a file "<job>-<n>.<ext>" when set.
	Dir string
	// Limits bound the decoding of requests.
	Limits *ipp.Limits
	// MaxDocument bounds a document the printer keeps (256 MiB when 0).
	MaxDocument int64

	// Requests counts the operations received, by name.
	Requests map[string]int

	mu   sync.Mutex
	ids  *idSource
	jobs []*Job
	// cutAfter cuts the connection of the next document after so many
	// bytes (-1: no); dropCreate cuts it before the answer to the next
	// Create-Job, once the job is made; cutControl cuts an LPD job in its
	// control file; refuse makes LPD refuse its queue; back is what 9100
	// says back once it has a job.
	cutAfter   int64
	dropCreate bool
	cutControl bool
	refuse     bool
	back       string
	abortNext  bool
	faults     []*Fault
	certs      *certificate
	stateSince time.Time
	started    time.Time
	watchers   []chan<- Event
}

// idSource gives job ids: a CUPS shares one between its printers.
type idSource struct {
	mu   sync.Mutex
	next int32
}

func (s *idSource) take() int32 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	return s.next
}

func (s *idSource) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next = 0
}

// Job is a job a printer received, over any transport.
type Job struct {
	ID int32
	// Transport is "ipp", "socket", "lpd" or "serial".
	Transport string
	// Printer is the name of the printer that has it.
	Printer string
	Name    string
	User    string
	// Format is the document-format of the last document (IPP), or what
	// the data was told to be (raw: a language, "postscript", "pdf" …).
	Format string
	// Attrs are the job template attributes of an IPP job.
	Attrs     []ipp.Attribute
	Documents [][]byte
	// Queue and Control are the queue and the control file of an LPD job.
	Queue, Control string
	Created        time.Time
	UUID           string
	// doneAt is when the last document arrived; the state is derived
	// from it and Processing.
	doneAt   time.Time
	canceled bool
	aborted  bool
	held     bool
}

// Document returns the bytes of a document of a job, for a test to look
// at.
func (j *Job) Document(n int) []byte {
	if n < 0 || n >= len(j.Documents) {
		return nil
	}
	return bytes.Clone(j.Documents[n])
}

// Data returns the documents of a job, one after the other: what a raw
// printer was sent.
func (j *Job) Data() []byte { return bytes.Join(j.Documents, nil) }

// Attr returns a job template attribute of an IPP job as a string, or "".
func (j *Job) Attr(name string) string {
	if a := ipp.Get(j.Attrs, name); a != nil {
		return a.String()
	}
	return ""
}

// New returns a printer with the capabilities of a plain office printer:
// PDF, PWG Raster and JPEG, A4, Letter, B5 and A3, two-sided, colour, 300
// and 600 dpi, up to 99 copies.
func New() *Printer {
	return &Printer{
		Name: "fakeprinter", MakeModel: "BDF Fake Printer", UUID: "urn:uuid:4e7b5a0e-5c3e-4f1a-9f6b-000000000001",
		Formats:              []string{"application/pdf", "image/pwg-raster", "image/jpeg", "application/octet-stream"},
		Media:                []string{"iso_a4_210x297mm", "na_letter_8.5x11in", "jis_b5_182x257mm", "iso_a3_297x420mm"},
		MediaReady:           []string{"iso_a4_210x297mm"},
		Sources:              []string{"auto", "tray-1"},
		Types:                []string{"auto", "stationery"},
		Sides:                []string{"one-sided", "two-sided-long-edge", "two-sided-short-edge"},
		ColorModes:           []string{"auto", "color", "monochrome"},
		Qualities:            []int{3, 4, 5},
		Resolutions:          []ipp.Resolution{{X: 300, Y: 300, Units: ipp.UnitsDPI}, {X: 600, Y: 600, Units: ipp.UnitsDPI}},
		Copies:               ipp.Range{Lo: 1, Hi: 99},
		PageRanges:           true,
		Finishings:           []int{3},
		Scaling:              []string{"auto", "auto-fit", "fill", "fit", "none"},
		OutputBins:           []string{"face-down"},
		PWGRasterTypes:       []string{"black_1", "sgray_8", "srgb_8"},
		PWGRasterResolutions: []ipp.Resolution{{X: 300, Y: 300, Units: ipp.UnitsDPI}, {X: 600, Y: 600, Units: ipp.UnitsDPI}},
		SheetBack:            "normal",
		URF:                  []string{"V1.4", "W8", "SRGB24", "RS300-600", "DM1"},
		State:                ipp.PrinterIdle, Reasons: []string{"none"}, Accepting: true,
		Requests: map[string]int{}, cutAfter: -1,
	}
}

// NewUUID returns a fresh urn:uuid, for a printer or a job.
func NewUUID() string {
	var b [16]byte
	rand.Read(b[:])
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("urn:uuid:%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// init sets what New sets for a Printer made as a literal.
func (p *Printer) init() {
	if p.Requests == nil {
		p.Requests = map[string]int{}
	}
	if p.ids == nil {
		p.ids = &idSource{}
	}
	if p.started.IsZero() {
		p.started = time.Now()
		p.stateSince = p.started
		if p.cutAfter == 0 && !p.dropCreate {
			p.cutAfter = -1
		}
	}
	if p.State == 0 {
		p.State = ipp.PrinterIdle
	}
	if len(p.Reasons) == 0 {
		p.Reasons = []string{"none"}
	}
}

// lock locks the printer, making it ready first.
func (p *Printer) lock() {
	p.mu.Lock()
	p.init()
}

// Stop stops the printer for reasons ("media-empty-error" …): it answers
// as stopped and does not finish its jobs until Resume.
func (p *Printer) Stop(reasons ...string) {
	p.lock()
	defer p.mu.Unlock()
	if len(reasons) == 0 {
		reasons = []string{"paused"}
	}
	p.State, p.Reasons, p.stateSince = ipp.PrinterStopped, reasons, time.Now()
	p.notify(Event{Kind: "state"})
}

// Resume makes a stopped printer idle again.
func (p *Printer) Resume() {
	p.lock()
	defer p.mu.Unlock()
	p.State, p.Reasons, p.stateSince = ipp.PrinterIdle, []string{"none"}, time.Now()
	p.notify(Event{Kind: "state"})
}

// SetAccepting makes the printer take jobs or refuse them
// (server-error-not-accepting-jobs, and LPD's refusal).
func (p *Printer) SetAccepting(ok bool) {
	p.lock()
	defer p.mu.Unlock()
	p.Accepting = ok
	p.notify(Event{Kind: "state"})
}

// Forget drops the jobs the printer knows, as a printer that keeps no
// history does: their ids are not reused.
func (p *Printer) Forget() {
	p.lock()
	defer p.mu.Unlock()
	p.jobs = nil
	p.notify(Event{Kind: "jobs"})
}

// Restart forgets the jobs and counts the job ids from 1 again, as a
// printer that was switched off and on does: a client that kept a job id
// finds another job, or none, under it.
func (p *Printer) Restart() {
	p.lock()
	defer p.mu.Unlock()
	p.jobs = nil
	p.ids.reset()
	p.State, p.Reasons, p.stateSince, p.started = ipp.PrinterIdle, []string{"none"}, time.Now(), time.Now()
	p.notify(Event{Kind: "jobs"})
}

// CutAfter cuts the connection of the next document sent after n bytes
// of it (n >= 0), as a network that fails does; the job keeps no
// document. Over LPD the connection is cut once the data file is sent,
// before the control file, whatever n is.
func (p *Printer) CutAfter(n int64) {
	p.lock()
	defer p.mu.Unlock()
	p.cutAfter = n
}

// CutControl cuts the connection of the next LPD job in its control file.
func (p *Printer) CutControl() {
	p.lock()
	defer p.mu.Unlock()
	p.cutControl = true
}

// DropCreateResponse makes the next Create-Job make its job and then cut
// the connection before answering: the client does not learn the job.
func (p *Printer) DropCreateResponse() {
	p.lock()
	defer p.mu.Unlock()
	p.dropCreate = true
}

// Refuse makes the printer refuse jobs on LPD (its queue is disabled), or
// take them again.
func (p *Printer) Refuse(refuse bool) {
	p.lock()
	defer p.mu.Unlock()
	p.refuse = refuse
}

// Back sets what a printer of the raw port says back once it has a job.
func (p *Printer) Back(back string) {
	p.lock()
	defer p.mu.Unlock()
	p.back = back
}

// AbortNext makes the next job the printer receives aborted.
func (p *Printer) AbortNext() {
	p.lock()
	defer p.mu.Unlock()
	p.abortNext = true
}

// Abort makes a job aborted, to test a client's handling.
func (p *Printer) Abort(id int32) {
	p.lock()
	defer p.mu.Unlock()
	if j := p.job(id); j != nil {
		j.aborted = true
		p.notify(Event{Kind: "job", Job: id})
	}
}

// Status is a printer's state at a moment.
type Status struct {
	State     ipp.PrinterState
	Reasons   []string
	Accepting bool
}

// Snapshot returns the printer's state now.
func (p *Printer) Snapshot() Status {
	p.lock()
	defer p.mu.Unlock()
	return Status{State: p.State, Reasons: append([]string(nil), p.Reasons...), Accepting: p.Accepting}
}

// Jobs returns the jobs received, oldest first.
func (p *Printer) Jobs() []*Job {
	p.lock()
	defer p.mu.Unlock()
	return append([]*Job(nil), p.jobs...)
}

// Job returns a job by id, or nil.
func (p *Printer) Job(id int32) *Job {
	p.lock()
	defer p.mu.Unlock()
	return p.job(id)
}

func (p *Printer) job(id int32) *Job {
	for _, j := range p.jobs {
		if j.ID == id {
			return j
		}
	}
	return nil
}

// JobState returns the state of a job now.
func (p *Printer) JobState(j *Job) ipp.JobState {
	p.lock()
	defer p.mu.Unlock()
	return p.state(j)
}

func (p *Printer) state(j *Job) ipp.JobState {
	switch {
	case j.canceled:
		return ipp.JobCanceled
	case j.aborted:
		return ipp.JobAborted
	case j.held:
		return ipp.JobPendingHeld
	case j.doneAt.IsZero():
		return ipp.JobPending
	case p.NeverComplete, time.Since(j.doneAt) < p.Processing, p.State == ipp.PrinterStopped:
		return ipp.JobProcessing
	}
	return ipp.JobCompleted
}

// add records a new job: its id, and the faults that wait for a number of
// jobs. The printer is locked.
func (p *Printer) add(j *Job) {
	j.ID = p.ids.take()
	j.Printer = p.Name
	j.Created = time.Now()
	if j.UUID == "" {
		j.UUID = NewUUID()
	}
	if p.abortNext {
		p.abortNext = false
		j.aborted = true
	}
	p.jobs = append(p.jobs, j)
	p.notify(Event{Kind: "job", Job: j.ID})
}

// done marks the last document of a job arrived. The printer is locked.
func (p *Printer) done(j *Job) {
	j.doneAt = time.Now()
	p.notify(Event{Kind: "job", Job: j.ID})
	p.fire(true)
}

// Event tells a watcher of a printer that something changed.
type Event struct {
	// Kind is "job" (a job came or changed), "jobs" (they were dropped)
	// or "state".
	Kind    string
	Printer string
	Job     int32
}

// Watch sends the printer's events to c until stop is called; an event
// that c has no room for is dropped.
func (p *Printer) Watch(c chan<- Event) (stop func()) {
	p.lock()
	defer p.mu.Unlock()
	p.watchers = append(p.watchers, c)
	return func() {
		p.lock()
		defer p.mu.Unlock()
		for i, w := range p.watchers {
			if w == c {
				p.watchers = append(p.watchers[:i], p.watchers[i+1:]...)
				break
			}
		}
	}
}

// notify sends an event to the watchers. The printer is locked.
func (p *Printer) notify(e Event) {
	e.Printer = p.Name
	for _, w := range p.watchers {
		select {
		case w <- e:
		default:
		}
	}
}
