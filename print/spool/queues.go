package spool

import (
	"context"
	"fmt"
	"runtime"
	"slices"
	"sort"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/raw"
)

func numCPU() int { return runtime.NumCPU() }

// QueueInfo is a queue and how many of its jobs are in each state.
type QueueInfo struct {
	Name     string   `json:"name"`
	Printers []string `json:"printers"`
	Policy   string   `json:"policy"`
	// Paused queues take jobs and do not send them; Rejecting ones do
	// not take jobs.
	Paused    bool `json:"paused"`
	Rejecting bool `json:"rejecting"`
	// Jobs counts the jobs not done, by state.
	Jobs map[string]int `json:"jobs"`
}

// Queues returns the queues, by name.
func (s *Spool) Queues() []QueueInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]QueueInfo, 0, len(s.queues))
	for _, q := range s.queues {
		out = append(out, s.queueInfo(q))
	}
	sort.Slice(out, func(a, b int) bool { return out[a].Name < out[b].Name })
	return out
}

// Queue returns a queue.
func (s *Spool) Queue(name string) (QueueInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	q := s.queues[name]
	if q == nil {
		return QueueInfo{}, fmt.Errorf("%w: %q", ErrNoQueue, name)
	}
	return s.queueInfo(q), nil
}

func (s *Spool) queueInfo(q *queueState) QueueInfo {
	info := QueueInfo{Name: q.Name, Printers: slices.Clone(q.Printers), Policy: q.Policy, Paused: q.rec.Paused, Rejecting: q.rec.Rejecting, Jobs: map[string]int{}}
	for _, e := range s.jobs {
		if e.job.Queue == q.Name {
			info.Jobs[e.job.State.String()]++
		}
	}
	return info
}

// PauseQueue stops sending the jobs of a queue (it still takes them);
// ResumeQueue sends them again.
func (s *Spool) PauseQueue(ctx context.Context, name string) (QueueInfo, error) {
	return s.setQueue(ctx, name, func(r *queueRecord) { r.Paused = true })
}

func (s *Spool) ResumeQueue(ctx context.Context, name string) (QueueInfo, error) {
	return s.setQueue(ctx, name, func(r *queueRecord) { r.Paused = false })
}

// RejectQueue stops taking jobs in a queue; AcceptQueue takes them again.
func (s *Spool) RejectQueue(ctx context.Context, name string) (QueueInfo, error) {
	return s.setQueue(ctx, name, func(r *queueRecord) { r.Rejecting = true })
}

func (s *Spool) AcceptQueue(ctx context.Context, name string) (QueueInfo, error) {
	return s.setQueue(ctx, name, func(r *queueRecord) { r.Rejecting = false })
}

func (s *Spool) setQueue(ctx context.Context, name string, f func(*queueRecord)) (QueueInfo, error) {
	s.mu.Lock()
	q := s.queues[name]
	if q == nil {
		s.mu.Unlock()
		return QueueInfo{}, fmt.Errorf("%w: %q", ErrNoQueue, name)
	}
	rec, rev := q.rec, q.rev
	s.mu.Unlock()
	f(&rec)
	ctx = context.WithoutCancel(ctx)
	if err := s.st.setQueue(ctx, name, rec, rev); err != nil {
		return QueueInfo{}, err
	}
	_, nrev, err := s.st.queue(ctx, name)
	if err != nil {
		return QueueInfo{}, err
	}
	s.mu.Lock()
	q.rec, q.rev = rec, nrev
	info := s.queueInfo(q)
	s.mu.Unlock()
	s.publish(Event{Queue: &info})
	s.wake()
	return info, nil
}

// printerState is a printer as the spool sees it.
type printerState struct {
	Printer
	client *ipp.Client
	raw    *raw.Printer
	// caps is what the printer said last (Get-Printer-Attributes); up
	// tells it answered, checked when, and err why not.
	caps    *ipp.Printer
	up      bool
	checked time.Time
	err     string
	// failing counts the failures in a row, for the delay of the next
	// check.
	failing int
	// sending is set while a job is sent to it: one at a time.
	sending bool
	pokes   chan struct{}
	// seenBy is the proxy that reported last, of a printer of proxies.
	seenBy string
}

func newPrinterState(p Printer) *printerState {
	return &printerState{Printer: p, client: p.Client, raw: p.Raw, pokes: make(chan struct{}, 1)}
}

// url is the printer's URL.
func (p *printerState) url() string {
	if p.Proxy {
		return "proxy:" + p.Name
	}
	if p.raw != nil {
		return p.raw.URL
	}
	return p.client.URL
}

func (s *Spool) printer(name string) *printerState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.printers[name]
}

// PrinterInfo is a printer and what it said last.
type PrinterInfo struct {
	Name      string   `json:"name"`
	URL       string   `json:"url"`
	Up        bool     `json:"up"`
	State     string   `json:"state,omitempty"`
	Reasons   []string `json:"reasons,omitempty"`
	Message   string   `json:"message,omitempty"`
	Accepting bool     `json:"accepting"`
	MakeModel string   `json:"make-and-model,omitempty"`
	// Checked is when the printer was last asked, and Error why it did
	// not answer.
	Checked time.Time `json:"checked,omitzero"`
	Error   string    `json:"error,omitempty"`
	// Jobs are the jobs sent to it and not done.
	Jobs int `json:"jobs"`
}

// Printers returns the printers, by name.
func (s *Spool) Printers() []PrinterInfo {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]PrinterInfo, 0, len(s.printers))
	for _, p := range s.printers {
		out = append(out, s.printerInfo(p))
	}
	sort.Slice(out, func(a, b int) bool { return out[a].Name < out[b].Name })
	return out
}

// PrinterCapabilities returns what a printer said last it supports, or
// nil before it answered.
func (s *Spool) PrinterCapabilities(name string) (*ipp.Printer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p := s.printers[name]
	if p == nil {
		return nil, fmt.Errorf("%w: printer %q", ErrInvalid, name)
	}
	return p.caps, nil
}

func (s *Spool) printerInfo(p *printerState) PrinterInfo {
	info := PrinterInfo{Name: p.Name, URL: p.url(), Up: p.up, Checked: p.checked, Error: p.err}
	if c := p.caps; c != nil {
		info.State, info.Reasons, info.Message, info.Accepting, info.MakeModel = c.State.String(), slices.Clone(c.Reasons), c.StateMessage, c.AcceptingJobs, c.MakeModel
	}
	for _, e := range s.jobs {
		if e.job.Printer != nil && e.job.Printer.Name == p.Name && (e.job.Step == StepWatch || e.job.Step == StepRemote) {
			info.Jobs++
		}
	}
	return info
}
