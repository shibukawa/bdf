package ipp

import (
	"bufio"
	"context"
	"maps"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/print/ipp"
)

// Create-Job, Send-Document and Close-Job: on a CUPS server node's own
// port, where one process listens, a job may come in two requests, as
// IPP Everywhere asks of a printer and as the clients of Mopria, Windows
// and CUPS send when a printer takes them. Create-Job takes the job's
// number (Manager.ReserveJobNumber) and keeps what it said (the job's
// name, its sender, its metadata) in the process; Send-Document brings
// the document, and the job is submitted with that number, as Print-Job
// does. A job whose document does not come within createdTimeout is
// aborted, its number skipped. A job has one document
// (multiple-document-jobs-supported false). The manager's own server
// does not take them: there, as on a serverless deployment, a request is
// all there is to a job.

// createdTimeout is how long a job Create-Job made waits for its
// document (multiple-operation-time-out); createdKept how long one ended
// without a document is still told of (tests change them).
var (
	createdTimeout = 5 * time.Minute
	createdKept    = 10 * time.Minute
)

// maxCreated bounds the jobs that wait for their documents.
const maxCreated = 1000

// created is a job Create-Job made: it waits for its document, or it
// ended without one (canceled, or aborted).
type created struct {
	number int32
	dest   string
	// in is the job to submit: its name, its sender and its metadata.
	in docflow.Input
	// user is who made it, as ippRequest.user tells: who sends its
	// document, closes it and cancels it as theirs.
	user     string
	at       time.Time
	deadline time.Time
	state    ipp.JobState
	reason   string
	ended    time.Time
	// sending tells a document arriving.
	sending bool
}

// createdJobs are the jobs Create-Job made of a store.
type createdJobs struct {
	mu   sync.Mutex
	jobs map[int32]*created
}

// createdBy holds the createdJobs of each store (Manager.StoreKey: a
// manager is made for a request, and a node's port served by a handler
// made for each).
var createdBy sync.Map // store → *createdJobs

func createdOf(m *docflow.Manager) *createdJobs {
	v, _ := createdBy.LoadOrStore(m.StoreKey(), &createdJobs{jobs: map[int32]*created{}})
	return v.(*createdJobs)
}

// expire aborts the jobs whose documents did not come in time, and
// forgets those ended long ago; with cs.mu held.
func (cs *createdJobs) expire(now time.Time) {
	for id, c := range cs.jobs {
		switch {
		case !c.state.Done() && !c.sending && now.After(c.deadline):
			c.state, c.reason, c.ended = ipp.JobAborted, "aborted-by-system", now
		case c.state.Done() && now.Sub(c.ended) > createdKept:
			delete(cs.jobs, id)
		}
	}
}

// get returns the job Create-Job made of a number, or nil.
func (cs *createdJobs) get(id int32) *created {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.expire(time.Now())
	return cs.jobs[id]
}

// of returns the jobs Create-Job made of a destination, newest first.
func (cs *createdJobs) of(dest string) []*created {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.expire(time.Now())
	var out []*created
	for _, c := range cs.jobs {
		if c.dest == dest {
			out = append(out, c)
		}
	}
	slices.SortFunc(out, func(a, b *created) int { return int(b.number - a.number) })
	return out
}

// end ends a job that waits for its document (canceled, aborted); false
// when it has ended already.
func (cs *createdJobs) end(c *created, st ipp.JobState, reason string) bool {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if c.state.Done() {
		return false
	}
	c.state, c.reason, c.ended = st, reason, time.Now()
	return true
}

// mine tells a job Create-Job made by the request's sender.
func (q *ippRequest) mine(c *created) bool {
	by := c.in.By
	return by.Kind == q.by.Kind && by.ID == q.by.ID && c.user == q.user()
}

// createJob takes a Create-Job: the job's number now, its document with
// Send-Document.
func (s *server) createJob(ctx context.Context, q *ippRequest) {
	if q.by.Kind == docflow.PrincipalApp {
		if err := s.m.MaySend(ctx, q.r, q.dest); err != nil {
			q.fail(ipp.StatusClientErrorNotAuthorized, "%v", err)
			return
		}
	}
	cs := createdOf(s.m)
	cs.mu.Lock()
	now := time.Now()
	cs.expire(now)
	waiting := 0
	for _, c := range cs.jobs {
		if !c.state.Done() {
			waiting++
		}
	}
	cs.mu.Unlock()
	if waiting >= maxCreated {
		q.fail(ipp.StatusServerErrorTooManyJobs, "too many jobs wait for their documents")
		return
	}
	number, err := s.m.ReserveJobNumber(ctx)
	if err != nil {
		s.m.Options().Log.Error("ipp: create-job", "dest", q.dest, "err", err)
		q.fail(ipp.StatusServerErrorInternal, "the job could not be made")
		return
	}
	c := &created{number: number, dest: q.dest, user: q.user(), at: now, deadline: now.Add(createdTimeout), state: ipp.JobPending, reason: "job-incoming",
		in: docflow.Input{Name: q.text("job-name"), By: q.by, ClaimedUser: q.text("requesting-user-name"), Meta: s.meta(q), Number: number}}
	cs.mu.Lock()
	cs.jobs[number] = c
	cs.mu.Unlock()
	q.res.Add(ipp.TagJob, s.createdAttributes(q, c)...)
}

// sendDocument takes the document of a job Create-Job made, and submits
// the job: with the first document, whether or not it is the last (a job
// has one). One without a document that is the last ends the job,
// aborted.
func (s *server) sendDocument(ctx context.Context, q *ippRequest, body *bufio.Reader) {
	c, j := s.jobOrCreated(ctx, q)
	if c == nil && j == nil {
		return
	}
	last, ok := q.attr("last-document").Bool()
	if q.attr("last-document") == nil || !ok {
		q.fail(ipp.StatusClientErrorBadRequest, "last-document is required")
		return
	}
	if j != nil {
		q.fail(ipp.StatusClientErrorNotPossible, "the job has its document: a job has one")
		return
	}
	if !q.mine(c) {
		q.fail(ipp.StatusClientErrorNotAuthorized, "the job is another's")
		return
	}
	cs := createdOf(s.m)
	cs.mu.Lock()
	switch {
	case c.state.Done():
		cs.mu.Unlock()
		q.fail(ipp.StatusClientErrorNotPossible, "the job is %s", c.state)
		return
	case c.sending:
		cs.mu.Unlock()
		q.fail(ipp.StatusServerErrorBusy, "a document of the job is arriving")
		return
	}
	c.sending = true
	cs.mu.Unlock()
	done := func() {
		cs.mu.Lock()
		c.sending = false
		cs.mu.Unlock()
	}
	if peek, _ := body.Peek(1); len(peek) == 0 {
		done()
		if last {
			cs.end(c, ipp.JobAborted, "job-data-insufficient")
		}
		q.res.Add(ipp.TagJob, s.createdAttributes(q, c)...)
		return
	}
	in := c.in
	in.Name = firstOr(in.Name, q.text("document-name"))
	// what Send-Document says of the document (document-format,
	// document-metadata) is kept as Print-Job's is (meta), within the
	// bound of the attributes kept
	in.Meta = maps.Clone(in.Meta)
	if in.Meta == nil {
		in.Meta = map[string]any{}
	}
	for k, v := range s.meta(q) {
		if _, had := in.Meta[k]; had || len(in.Meta) < maxIPPMeta {
			in.Meta[k] = v
		}
	}
	err := s.take(ctx, q, body, in)
	cs.mu.Lock()
	c.sending = false
	canceled := c.state.Done() // by Cancel-Job while the document came
	switch {
	case err == nil:
		// the job is a job now
		delete(cs.jobs, c.number)
	case !canceled:
		c.state, c.reason, c.ended = ipp.JobAborted, "aborted-by-system", time.Now()
	}
	cs.mu.Unlock()
	if err == nil && canceled {
		if j, err := s.m.JobOfNumber(ctx, c.number); err == nil {
			s.m.Cancel(ctx, j.ID, docflow.Mark{Worker: clip(c.user, 80)})
		}
	}
}

// closeJob ends the documents of a job: one Create-Job made that has no
// document is aborted; a job is told.
func (s *server) closeJob(ctx context.Context, q *ippRequest) {
	c, j := s.jobOrCreated(ctx, q)
	switch {
	case j != nil:
		q.res.Add(ipp.TagJob, s.jobStatus(q, j)...)
		return
	case c == nil:
		return
	case !q.mine(c):
		q.fail(ipp.StatusClientErrorNotAuthorized, "the job is another's")
		return
	}
	cs := createdOf(s.m)
	cs.mu.Lock()
	busy := c.sending
	cs.mu.Unlock()
	if busy {
		q.fail(ipp.StatusServerErrorBusy, "a document of the job is arriving")
		return
	}
	cs.end(c, ipp.JobAborted, "job-data-insufficient")
	q.res.Add(ipp.TagJob, s.createdAttributes(q, c)...)
}

// cancelMyJobs takes a Cancel-My-Jobs: the sender's jobs of the
// destination that have not ended (those of job-ids, when it says), as
// the sender is known: the token or the app it gives, else its
// requesting-user-name.
func (s *server) cancelMyJobs(ctx context.Context, q *ippRequest) {
	var ids []int32
	if a := q.attr("job-ids"); a != nil {
		ids = a.Ints()
	}
	pick := func(id int32) bool { return ids == nil || slices.Contains(ids, id) }
	cs := createdOf(s.m)
	for _, c := range cs.of(q.dest) {
		if q.mine(c) && pick(c.number) {
			cs.end(c, ipp.JobCanceled, "job-canceled-by-user")
		}
	}
	jobs, err := s.m.JobsOf(ctx, q.dest, 0)
	if err != nil {
		q.fail(ipp.StatusServerErrorInternal, "%v", err)
		return
	}
	user := q.text("requesting-user-name")
	for _, j := range jobs {
		by := j.SubmittedBy
		mine := by != nil && by.Kind == q.by.Kind && (q.by.ID != "" && by.ID == q.by.ID || q.by.ID == "" && user != "" && j.ClaimedUser == user)
		if !j.State.Ended() && mine && pick(j.Number) {
			s.m.Cancel(ctx, j.ID, docflow.Mark{Worker: clip(q.user(), 80)})
		}
	}
}

// createdAttributes are the attributes of a job Create-Job made, as
// Get-Job-Attributes tells a job (jobAttributes).
func (s *server) createdAttributes(q *ippRequest, c *created) []ipp.Attribute {
	cs := createdOf(s.m)
	cs.mu.Lock()
	st, reason, ended := c.state, c.reason, c.ended
	cs.mu.Unlock()
	uri := q.printerURI() + "/" + strconv.Itoa(int(c.number))
	attrs := []ipp.Attribute{
		ipp.Attr("job-id", ipp.Integer(int(c.number))),
		ipp.Attr("job-uri", ipp.URI(uri)),
		ipp.Attr("job-printer-uri", ipp.URI(q.printerURI())),
		ipp.Attr("job-name", ipp.Name(firstOr(c.in.Name, "job-"+strconv.Itoa(int(c.number))))),
		ipp.Attr("job-originating-user-name", ipp.Name(firstOr(c.in.ClaimedUser, firstOr(c.in.By.Label(), "anonymous")))),
		ipp.Attr("job-state", ipp.Enum(int(st))),
		ipp.Attr("job-state-reasons", ipp.Keyword(reason)),
		ipp.Attr("job-impressions-completed", ipp.Integer(0)),
		ipp.Attr("time-at-creation", ipp.Integer(int(c.at.Unix()&0x7fffffff))),
		ipp.Attr("time-at-processing", ipp.OutOfBand(ipp.TagNoValue)),
		ipp.Attr("job-printer-up-time", ipp.Integer(int(s.m.Now().Unix()&0x7fffffff))),
	}
	if st.Done() {
		attrs = append(attrs, ipp.Attr("time-at-completed", ipp.Integer(int(ended.Unix()&0x7fffffff))))
	} else {
		attrs = append(attrs, ipp.Attr("time-at-completed", ipp.OutOfBand(ipp.TagNoValue)))
	}
	return attrs
}
