package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/spool"
)

// upstream is a server this one fetches jobs from, as a proxy of IPP
// INFRA: it tells the upstream its printer (the queue the jobs go to),
// fetches the jobs offered, puts them in its spool, and reports their
// states. A job is taken once: the upstream binds a job to the proxy that
// acknowledges it, and the proxy finds a job it took (by the upstream
// job's uuid, kept as the job's Origin) before it fetches one.
type upstream struct {
	s      *Server
	name   string
	cfg    UpstreamConfig
	c      *ipp.Client
	device string
	poll   time.Duration
	// reported are the states last reported of the jobs, by local id.
	reported map[int32]string
	failing  bool
}

// upstreams makes the upstreams of the configuration.
func (s *Server) upstreams() ([]*upstream, error) {
	var out []*upstream
	for _, name := range sortedKeys(s.cfg.Upstreams) {
		u := s.cfg.Upstreams[name]
		c := &ipp.Client{URL: u.URL, User: s.cfg.User, Username: u.Username, Timeout: 2 * time.Minute}
		switch {
		case u.PasswordFile != "":
			b, err := os.ReadFile(u.PasswordFile)
			if err != nil {
				return nil, fmt.Errorf("upstream %s: %w", name, err)
			}
			c.Password = strings.TrimSpace(string(b))
		case u.PasswordEnv != "":
			c.Password = os.Getenv(u.PasswordEnv)
		}
		c.TLS.Insecure = u.Insecure
		if !u.Insecure {
			c.TLS.Store = s.sp.TrustStore()
		}
		device := u.UUID
		if device == "" {
			device = queueUUID(s.cfg.Spool, "upstream:"+name)
		}
		poll := time.Duration(u.Poll)
		if poll <= 0 {
			poll = 10 * time.Second
		}
		out = append(out, &upstream{s: s, name: name, cfg: u, c: c, device: device, poll: poll, reported: map[int32]string{}})
	}
	return out, nil
}

// run fetches and reports until ctx ends, then tells the upstream the
// proxy stopped.
func (u *upstream) run(ctx context.Context) {
	u.s.log.Info("upstream: fetching jobs", "upstream", u.name, "url", redact(u.cfg.URL), "queue", u.cfg.Queue, "device", u.device)
	timer := time.NewTimer(0)
	defer timer.Stop()
	events, stop := u.s.sp.Subscribe(64)
	defer stop()
	for {
		select {
		case <-ctx.Done():
			dctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			m := u.request(ipp.OpDeregisterOutputDevice)
			u.c.Do(dctx, m, nil)
			cancel()
			return
		case ev := <-events:
			// a job of the upstream changed here: reported now
			if ev.Job == nil || ev.Job.Origin == nil || ev.Job.Origin.Upstream != u.name {
				continue
			}
			u.report(ctx, *ev.Job)
			continue
		case <-timer.C:
		}
		err := u.cycle(ctx)
		switch {
		case err != nil && ctx.Err() == nil && !u.failing:
			u.failing = true
			u.s.log.Warn("upstream: not reached", "upstream", u.name, "err", err)
		case err == nil && u.failing:
			u.failing = false
			u.s.log.Info("upstream: reached again", "upstream", u.name)
		}
		timer.Reset(u.poll)
	}
}

// request makes a request of the proxy.
func (u *upstream) request(op ipp.Op) *ipp.Message {
	m := u.c.NewRequest(op)
	m.Add(ipp.TagOperation, ipp.Attr("output-device-uuid", ipp.URI(u.device)))
	return m
}

// cycle tells the upstream the printer, puts the jobs right, fetches the
// jobs offered and reports those it has.
func (u *upstream) cycle(ctx context.Context) error {
	if err := u.register(ctx); err != nil {
		return err
	}
	if err := u.activeJobs(ctx); err != nil {
		return err
	}
	if err := u.fetch(ctx); err != nil {
		return err
	}
	for _, j := range u.s.sp.Origins(u.name) {
		u.report(ctx, j)
	}
	return nil
}

// printerAttributesUp are the attributes of the printer that are not
// told to the upstream: those of where this server is reached.
var printerAttributesUp = []string{"printer-uri-supported", "uri-authentication-supported", "uri-security-supported", "printer-icons",
	"printer-more-info", "printer-supply-info-uri", "printer-uuid", "media-col-database"}

// register tells the upstream the state of the queue the jobs go to, and
// what its printer can do (Update-Output-Device-Attributes).
func (u *upstream) register(ctx context.Context) error {
	st, reasons, msg, accepting := u.s.queueState(u.cfg.Queue)
	m := u.request(ipp.OpUpdateOutputDeviceAttributes)
	attrs := []ipp.Attribute{
		ipp.Attr("printer-state", ipp.Enum(int(st))),
		ipp.Attr("printer-state-reasons", keywordValues(reasons)...),
		ipp.Attr("printer-state-message", ipp.Text(clip127(msg))),
		ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(accepting)),
	}
	if qc := u.s.cfg.Queues[u.cfg.Queue]; len(qc.Printers) > 0 {
		if caps, _ := u.s.sp.PrinterCapabilities(qc.Printers[0]); caps != nil {
			for _, a := range caps.Attrs {
				if !slices.Contains(printerAttributesUp, a.Name) && ipp.Get(attrs, a.Name) == nil {
					attrs = append(attrs, a)
				}
			}
		}
	}
	m.Add(ipp.TagPrinter, attrs...)
	_, err := u.c.Do(ctx, m, nil)
	return err
}

// activeJobs puts the jobs right with the upstream (Update-Active-Jobs):
// a job canceled there is canceled here.
func (u *upstream) activeJobs(ctx context.Context) error {
	jobs := u.s.sp.Origins(u.name)
	m := u.request(ipp.OpUpdateActiveJobs)
	var ids, states []ipp.Value
	for _, j := range jobs {
		ids = append(ids, ipp.Integer(int(j.Origin.JobID)))
		states = append(states, ipp.Enum(int(j.State)))
	}
	if len(ids) > 0 {
		m.Add(ipp.TagOperation, ipp.Attr("job-ids", ids...), ipp.Attr("output-device-job-states", states...))
	}
	res, err := u.c.Do(ctx, m, nil)
	if err != nil {
		return err
	}
	ended := res.Attr(ipp.TagOperation, "job-ids").Ints()
	for _, j := range jobs {
		if slices.Contains(ended, j.Origin.JobID) && !j.State.Done() {
			u.s.log.Info("upstream: the job ended upstream: canceled here", "upstream", u.name, "job", j.ID, "upstream-job", j.Origin.JobID)
			u.s.sp.Cancel(ctx, j.ID, "job-canceled-by-operator")
		}
	}
	return nil
}

// fetch fetches the jobs offered: a job taken before is acknowledged
// again (the proxy stopped before it said it had the document); another
// is fetched, put in the spool, and acknowledged.
func (u *upstream) fetch(ctx context.Context) error {
	m := u.request(ipp.OpGetJobs)
	m.Add(ipp.TagOperation, ipp.Attr("which-jobs", ipp.Keyword("fetchable")), ipp.Attr("limit", ipp.Integer(20)),
		ipp.Attr("requested-attributes", ipp.Keyword("job-id"), ipp.Keyword("job-uuid")))
	res, err := u.c.Do(ctx, m, nil)
	if err != nil {
		return err
	}
	for _, g := range res.Groups {
		if g.Tag != ipp.TagJob {
			continue
		}
		id, _ := g.Get("job-id").Int()
		uuid := g.Get("job-uuid").Text()
		if id <= 0 || uuid == "" {
			continue
		}
		if j, ok := u.s.sp.FindOrigin(ctx, u.name, uuid); ok {
			u.acknowledgeDocument(ctx, id)
			u.report(ctx, j)
			continue
		}
		if err := u.take(ctx, id, uuid); err != nil {
			if ctx.Err() != nil {
				return err
			}
			u.s.log.Warn("upstream: a job could not be taken", "upstream", u.name, "upstream-job", id, "err", err)
		}
	}
	return nil
}

// take fetches a job and its document and puts them in the spool.
func (u *upstream) take(ctx context.Context, id int32, uuid string) error {
	m := u.request(ipp.OpFetchJob)
	m.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(id))))
	res, err := u.c.Do(ctx, m, nil)
	if err != nil {
		return err
	}
	g := res.Group(ipp.TagJob)
	if g == nil {
		return errors.New("Fetch-Job told no job")
	}
	sub := spool.Submission{Queue: u.cfg.Queue, Name: g.Get("job-name").Text(), User: g.Get("job-originating-user-name").Text(),
		Origin: &spool.Origin{Upstream: u.name, JobID: id, UUID: uuid}}
	if p, ok := g.Get("job-priority").Int(); ok {
		sub.Priority = int(p)
	}
	if t := g.Get("bdf-ticket").Text(); t != "" {
		if err := json.Unmarshal([]byte(t), &sub.Ticket); err != nil {
			return fmt.Errorf("bdf-ticket: %w", err)
		}
	} else {
		sub.Ticket = ticketOf(g.Attrs)
	}
	if o := g.Get("bdf-job-options").Text(); o != "" {
		if err := json.Unmarshal([]byte(o), &sub.Options); err != nil {
			return fmt.Errorf("bdf-job-options: %w", err)
		}
	}
	// the hold was the upstream's: here it is printed
	sub.Ticket.HoldUntil = ""
	ack := u.request(ipp.OpAcknowledgeJob)
	ack.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(id))))
	if _, err := u.c.Do(ctx, ack, nil); err != nil {
		return err
	}
	fd := u.request(ipp.OpFetchDocument)
	fd.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(id))), ipp.Attr("document-number", ipp.Integer(1)),
		ipp.Attr("compression-accepted", ipp.Keyword("none")),
		ipp.Attr("document-format-accepted", ipp.MimeType(bdfFormat), ipp.MimeType("application/pdf")))
	res, body, err := u.c.DoStream(ctx, fd)
	if err != nil {
		return err
	}
	defer body.Close()
	sub.Type = res.Attr(ipp.TagOperation, "document-format").Text()
	if pw := res.Attr(ipp.TagOperation, "document-password"); pw != nil {
		sub.Password = pw.Text()
	}
	up, err := u.s.sp.Receive(ctx, body)
	if err != nil {
		return err
	}
	j, err := u.s.sp.Submit(ctx, spool.Submission{Queue: sub.Queue, Name: sub.Name, User: sub.User, Type: sub.Type, Upload: up,
		Ticket: sub.Ticket, Options: sub.Options, Priority: sub.Priority, Password: sub.Password, Origin: sub.Origin})
	if err != nil {
		return err
	}
	u.s.log.Info("upstream: job taken", "upstream", u.name, "upstream-job", id, "job", j.ID, "user", j.User, "name", j.Name)
	u.acknowledgeDocument(ctx, id)
	return nil
}

func (u *upstream) acknowledgeDocument(ctx context.Context, id int32) {
	m := u.request(ipp.OpAcknowledgeDocument)
	m.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(id))), ipp.Attr("document-number", ipp.Integer(1)))
	if _, err := u.c.Do(ctx, m, nil); err != nil {
		u.s.log.Debug("upstream: Acknowledge-Document", "upstream", u.name, "upstream-job", id, "err", err)
	}
}

// report tells the upstream the state of a job, when it changed
// (Update-Job-Status); a job that ended upstream is canceled here.
func (u *upstream) report(ctx context.Context, j spool.Job) {
	sig := fmt.Sprintf("%d %v %d %s", j.State, j.Reasons, j.Impressions, j.Message)
	if u.reported[j.ID] == sig {
		return
	}
	m := u.request(ipp.OpUpdateJobStatus)
	m.Add(ipp.TagOperation, ipp.Attr("job-id", ipp.Integer(int(j.Origin.JobID))))
	reasons := publicReasons(orNone(j.Reasons))
	attrs := []ipp.Attribute{
		ipp.Attr("output-device-job-state", ipp.Enum(int(j.State))),
		ipp.Attr("output-device-job-state-reasons", keywordValues(reasons)...),
		ipp.Attr("job-impressions-completed", ipp.Integer(j.Impressions)),
	}
	if j.Message != "" {
		attrs = append(attrs, ipp.Attr("output-device-job-state-message", ipp.Text(clip127(j.Message))))
	}
	m.Add(ipp.TagJob, attrs...)
	res, err := u.c.Do(ctx, m, nil)
	if err != nil {
		if ipp.IsStatus(err, ipp.StatusClientErrorNotFound) {
			u.reported[j.ID] = sig // the upstream forgot it
		}
		return
	}
	u.reported[j.ID] = sig
	if n, ok := res.Attr(ipp.TagJob, "job-state").Int(); ok && ipp.JobState(n).Done() && !j.State.Done() {
		u.s.log.Info("upstream: the job ended upstream: canceled here", "upstream", u.name, "job", j.ID, "upstream-job", j.Origin.JobID)
		u.s.sp.Cancel(ctx, j.ID, "job-canceled-by-operator")
	}
}

// ticketOf reads a ticket from the job template attributes of a job
// (Fetch-Job of an upstream that is not bdf printd).
func ticketOf(attrs []ipp.Attribute) print.Ticket {
	var t print.Ticket
	g := &ipp.Group{Attrs: attrs}
	t.JobName = g.Get("job-name").Text()
	if n, ok := g.Get("copies").Int(); ok && n > 0 && n < 10000 {
		t.Copies = int(n)
	}
	switch v := g.Get("sides").Text(); v {
	case "one-sided", "two-sided-long-edge", "two-sided-short-edge":
		t.Sides = v
	}
	if m := g.Get("media").Text(); m != "" {
		if _, err := print.ParseMedia(m); err == nil {
			t.Media = m
		}
	}
	switch v := g.Get("print-color-mode").Text(); v {
	case "color", "monochrome":
		t.ColorMode = v
	}
	if pages, ok := pageRanges(*orEmpty(g.Get("page-ranges"))); ok {
		t.Pages = pages
	}
	return t
}

func orEmpty(a *ipp.Attribute) *ipp.Attribute {
	if a == nil {
		return &ipp.Attribute{}
	}
	return a
}
