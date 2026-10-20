// Package client talks to a print server (bdf printd) over its HTTP API:
// it submits documents, follows and changes jobs, and reads the queues,
// the printers and the accounting. It depends on nothing of bdf: its
// types have the shape of the API's JSON.
package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Client is a client of a print server.
type Client struct {
	// URL is the server's: http://127.0.0.1:8640 by default.
	URL string
	// Token is the bearer token; "" sends none (a server without tokens).
	Token string
	// User names the user of the jobs (X-BDF-User): a server takes it
	// from a token that may act for others, or when it has no tokens.
	User string
	// HTTP is the HTTP client; http.DefaultClient by default.
	HTTP *http.Client
}

// Error is an error the server answered.
type Error struct {
	Status  int    // the HTTP status
	Code    string // "not-found", "not-accepting", "full" …
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("server: %s (%d %s)", e.Message, e.Status, e.Code)
}

// IsCode tells an error the server answered with a code.
func IsCode(err error, code string) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}

func (c *Client) base() string {
	if c.URL == "" {
		return "http://127.0.0.1:8640"
	}
	return strings.TrimRight(c.URL, "/")
}

func (c *Client) http() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}

// request sends a request and returns the answer of a status below 300,
// or the server's error.
func (c *Client) request(ctx context.Context, method, path string, body io.Reader, typ string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.base()+path, body)
	if err != nil {
		return nil, err
	}
	if typ != "" {
		req.Header.Set("Content-Type", typ)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	if c.User != "" {
		req.Header.Set("X-BDF-User", c.User)
	}
	res, err := c.http().Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 300 {
		return res, nil
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
	var e struct {
		Error struct{ Code, Message string }
	}
	if json.Unmarshal(b, &e) != nil || e.Error.Code == "" {
		return nil, &Error{Status: res.StatusCode, Code: "http", Message: strings.TrimSpace(string(b))}
	}
	return nil, &Error{Status: res.StatusCode, Code: e.Error.Code, Message: e.Error.Message}
}

// call sends a request with a JSON body (nil: none) and decodes the JSON
// answer into out (nil: drops it).
func (c *Client) call(ctx context.Context, method, path string, in, out any) error {
	var body io.Reader
	typ := ""
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body, typ = bytes.NewReader(b), "application/json"
	}
	res, err := c.request(ctx, method, path, body, typ)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if out == nil {
		_, err = io.Copy(io.Discard, res.Body)
		return err
	}
	return json.NewDecoder(res.Body).Decode(out)
}

// Submission is a job to submit.
type Submission struct {
	Queue string
	// Name is the document's file name, which the server's converter
	// reads the format of some inputs from (.csv, .md).
	Name     string
	Document io.Reader
	// Ticket is what to print: a print.Ticket, or anything that is the
	// JSON of one ({"copies": 2, "sides": "two-sided-long-edge"}).
	Ticket   any
	Options  *JobOptions
	Password string
	// Priority is 1 to 100; 0 is the server's (50).
	Priority      int
	HoldUntil     string
	HoldUntilTime time.Time
}

// Submit submits a job. The document is streamed after the other fields.
func (c *Client) Submit(ctx context.Context, sub Submission) (*Job, error) {
	if sub.Queue == "" {
		return nil, errors.New("client: no queue")
	}
	if sub.Document == nil {
		return nil, errors.New("client: no document")
	}
	fields := [][2]string{}
	add := func(name string, v any) error {
		if v == nil {
			return nil
		}
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("client: %s: %w", name, err)
		}
		if string(b) != "{}" && string(b) != "null" {
			fields = append(fields, [2]string{name, string(b)})
		}
		return nil
	}
	if err := add("ticket", sub.Ticket); err != nil {
		return nil, err
	}
	if sub.Options != nil {
		if err := add("options", sub.Options); err != nil {
			return nil, err
		}
	}
	if sub.Name != "" {
		fields = append(fields, [2]string{"name", sub.Name})
	}
	if sub.Password != "" {
		fields = append(fields, [2]string{"password", sub.Password})
	}
	if sub.Priority != 0 {
		fields = append(fields, [2]string{"priority", strconv.Itoa(sub.Priority)})
	}
	if sub.HoldUntil != "" {
		fields = append(fields, [2]string{"hold-until", sub.HoldUntil})
	}
	if !sub.HoldUntilTime.IsZero() {
		fields = append(fields, [2]string{"hold-until-time", sub.HoldUntilTime.Format(time.RFC3339)})
	}
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	go func() {
		err := func() error {
			for _, f := range fields {
				if err := mw.WriteField(f[0], f[1]); err != nil {
					return err
				}
			}
			name := sub.Name
			if name == "" {
				name = "document"
			}
			w, err := mw.CreateFormFile("document", name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(w, sub.Document); err != nil {
				return err
			}
			return mw.Close()
		}()
		pw.CloseWithError(err)
	}()
	res, err := c.request(ctx, "POST", "/api/v1/queues/"+url.PathEscape(sub.Queue)+"/jobs", pr, mw.FormDataContentType())
	pr.CloseWithError(errors.New("client: the request ended"))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var j Job
	return &j, json.NewDecoder(res.Body).Decode(&j)
}

// Job returns a job.
func (c *Client) Job(ctx context.Context, id int32) (*Job, error) {
	var j Job
	return &j, c.call(ctx, "GET", fmt.Sprintf("/api/v1/jobs/%d", id), nil, &j)
}

// JobFilter selects jobs to list.
type JobFilter struct {
	Queue, User string
	// State is "" (the jobs not done), "all" (with those done within the
	// server's retention of documents) or a job-state.
	State string
	// After lists the jobs of greater ids; Limit how many (100 by
	// default, at most 1000).
	After int32
	Limit int
}

// JobList is a page of jobs; Next is the After of the next page, when
// there may be one.
type JobList struct {
	Jobs []Job `json:"jobs"`
	Next int32 `json:"next,omitempty"`
}

// Jobs lists a page of jobs.
func (c *Client) Jobs(ctx context.Context, f JobFilter) (*JobList, error) {
	q := url.Values{}
	for k, v := range map[string]string{"queue": f.Queue, "user": f.User, "state": f.State} {
		if v != "" {
			q.Set(k, v)
		}
	}
	if f.After != 0 {
		q.Set("after", strconv.Itoa(int(f.After)))
	}
	if f.Limit != 0 {
		q.Set("limit", strconv.Itoa(f.Limit))
	}
	var l JobList
	return &l, c.call(ctx, "GET", "/api/v1/jobs?"+q.Encode(), nil, &l)
}

func (c *Client) jobOp(ctx context.Context, id int32, op string, body any) (*Job, error) {
	var j Job
	return &j, c.call(ctx, "POST", fmt.Sprintf("/api/v1/jobs/%d/%s", id, op), body, &j)
}

// Cancel cancels a job, at its printer too.
func (c *Client) Cancel(ctx context.Context, id int32) (*Job, error) {
	return c.jobOp(ctx, id, "cancel", nil)
}

// Hold holds a job that is not sent yet.
func (c *Client) Hold(ctx context.Context, id int32) (*Job, error) {
	return c.jobOp(ctx, id, "hold", nil)
}

// Release releases a held job; a password opens an encrypted document
// that waits for one.
func (c *Client) Release(ctx context.Context, id int32, password string) (*Job, error) {
	return c.jobOp(ctx, id, "release", map[string]string{"password": password})
}

// SetPassword gives the password of an encrypted document.
func (c *Client) SetPassword(ctx context.Context, id int32, password string) (*Job, error) {
	return c.jobOp(ctx, id, "password", map[string]string{"password": password})
}

// Retry sends again a job that was aborted or canceled.
func (c *Client) Retry(ctx context.Context, id int32) (*Job, error) {
	return c.jobOp(ctx, id, "retry", nil)
}

// Move moves a job not sent yet to another queue.
func (c *Client) Move(ctx context.Context, id int32, queue string) (*Job, error) {
	return c.jobOp(ctx, id, "move", map[string]string{"queue": queue})
}

// SetPriority changes a job's priority (1 to 100).
func (c *Client) SetPriority(ctx context.Context, id int32, priority int) (*Job, error) {
	var j Job
	return &j, c.call(ctx, "PATCH", fmt.Sprintf("/api/v1/jobs/%d", id), map[string]int{"priority": priority}, &j)
}

// Queues returns the queues.
func (c *Client) Queues(ctx context.Context) ([]Queue, error) {
	var out struct{ Queues []Queue }
	return out.Queues, c.call(ctx, "GET", "/api/v1/queues", nil, &out)
}

// Queue returns a queue.
func (c *Client) Queue(ctx context.Context, name string) (*Queue, error) {
	var q Queue
	return &q, c.call(ctx, "GET", "/api/v1/queues/"+url.PathEscape(name), nil, &q)
}

// QueueOp pauses ("pause"), resumes ("resume"), rejects ("reject") or
// accepts ("accept") a queue.
func (c *Client) QueueOp(ctx context.Context, name, op string) (*Queue, error) {
	var q Queue
	return &q, c.call(ctx, "POST", "/api/v1/queues/"+url.PathEscape(name)+"/"+url.PathEscape(op), nil, &q)
}

// Printers returns the printers.
func (c *Client) Printers(ctx context.Context) ([]Printer, error) {
	var out struct{ Printers []Printer }
	return out.Printers, c.call(ctx, "GET", "/api/v1/printers", nil, &out)
}

// Printer returns a printer, with what it supports.
func (c *Client) Printer(ctx context.Context, name string) (*Printer, error) {
	var p Printer
	return &p, c.call(ctx, "GET", "/api/v1/printers/"+url.PathEscape(name), nil, &p)
}

// Accounting writes the jobs done from to to (zero times do not bound;
// to excludes) as "jsonl" or "csv".
func (c *Client) Accounting(ctx context.Context, from, to time.Time, format string, w io.Writer) error {
	q := url.Values{}
	if format != "" {
		q.Set("format", format)
	}
	if !from.IsZero() {
		q.Set("from", from.Format(time.RFC3339))
	}
	if !to.IsZero() {
		q.Set("to", to.Format(time.RFC3339))
	}
	return c.download(ctx, "/api/v1/accounting?"+q.Encode(), w)
}

// Backup writes a copy of the server's spool.db.
func (c *Client) Backup(ctx context.Context, w io.Writer) error {
	return c.download(ctx, "/api/v1/backup", w)
}

func (c *Client) download(ctx context.Context, path string, w io.Writer) error {
	res, err := c.request(ctx, "GET", path, nil, "")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(w, res.Body)
	return err
}

// Event is a change the server tells: of a job, a queue or a printer.
type Event struct {
	Job     *Job
	Queue   *Queue
	Printer *Printer
}

// Events calls f with the changes the server tells until ctx ends, f
// returns an error, or the stream ends (io.EOF); job (0: all) follows one
// job. Changes the server could not send in time are not told: a client
// that comes back reads the state again.
func (c *Client) Events(ctx context.Context, job int32, f func(Event) error) error {
	return c.events(ctx, job, nil, f)
}

// events is Events; opened is called once the server has the stream.
func (c *Client) events(ctx context.Context, job int32, opened func(), f func(Event) error) error {
	path := "/api/v1/events"
	if job != 0 {
		path += fmt.Sprintf("?job=%d", job)
	}
	res, err := c.request(ctx, "GET", path, nil, "")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if opened != nil {
		opened()
	}
	sc := bufio.NewScanner(res.Body)
	sc.Buffer(make([]byte, 64<<10), 4<<20)
	var name string
	var data []byte
	for sc.Scan() {
		line := sc.Text()
		switch {
		case line == "":
			if len(data) > 0 {
				var ev Event
				var err error
				switch name {
				case "job":
					ev.Job = new(Job)
					err = json.Unmarshal(data, ev.Job)
				case "queue":
					ev.Queue = new(Queue)
					err = json.Unmarshal(data, ev.Queue)
				case "printer":
					ev.Printer = new(Printer)
					err = json.Unmarshal(data, ev.Printer)
				}
				if err != nil {
					return err
				}
				if ev != (Event{}) {
					if err := f(ev); err != nil {
						return err
					}
				}
			}
			name, data = "", nil
		case strings.HasPrefix(line, "event:"):
			name = strings.TrimSpace(line[6:])
		case strings.HasPrefix(line, "data:"):
			data = append(data, strings.TrimPrefix(line[5:], " ")...)
		}
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err := sc.Err(); err != nil {
		return err
	}
	return io.EOF
}

// Wait follows a job until it is done and returns it; each change of its
// state, reasons or printer goes to changed (nil: none). It follows the
// server's events, and asks for the job when the stream opens, every
// while, and when the stream is lost, so that no change is missed.
func (c *Client) Wait(ctx context.Context, id int32, changed func(*Job)) (*Job, error) {
	var mu sync.Mutex
	var last *Job
	// tell takes a state of the job and tells whether it is done
	tell := func(n *Job) bool {
		mu.Lock()
		defer mu.Unlock()
		if last == nil || n.State != last.State || !slices.Equal(n.Reasons, last.Reasons) || printerName(n) != printerName(last) {
			if changed != nil {
				changed(n)
			}
		}
		last = n
		return n.Done()
	}
	get := func() (bool, error) {
		n, err := c.Job(ctx, id)
		if err != nil {
			return false, err
		}
		return tell(n), nil
	}
	if done, err := get(); err != nil || done {
		return last, err
	}
	errDone := errors.New("done")
	for {
		sctx, cancel := context.WithCancel(ctx)
		stop := make(chan struct{})
		var wg sync.WaitGroup
		err := c.events(sctx, id, func() {
			wg.Go(func() {
				t := time.NewTicker(waitCheck)
				defer t.Stop()
				for {
					if done, _ := get(); done {
						cancel()
						return
					}
					select {
					case <-stop:
						return
					case <-t.C:
					}
				}
			})
		}, func(ev Event) error {
			if ev.Job != nil && ev.Job.ID == id && tell(ev.Job) {
				return errDone
			}
			return nil
		})
		close(stop)
		cancel()
		wg.Wait()
		mu.Lock()
		l := last
		mu.Unlock()
		switch {
		case l.Done():
			return l, nil
		case ctx.Err() != nil:
			return l, ctx.Err()
		}
		var e *Error
		if errors.As(err, &e) && e.Status != http.StatusServiceUnavailable {
			return l, err
		}
		// the stream ended or the server keeps no more: ask after a
		// while, then open it again
		select {
		case <-ctx.Done():
			return l, ctx.Err()
		case <-time.After(waitPoll):
		}
		if done, err := get(); done {
			return last, nil
		} else if IsCode(err, "not-found") {
			return l, err
		}
	}
}

// waitPoll is how long Wait waits before it asks again when the events
// stream is lost; waitCheck how often it asks while the stream is open.
var (
	waitPoll  = 3 * time.Second
	waitCheck = 30 * time.Second
)

func printerName(j *Job) string {
	if j.Printer == nil {
		return ""
	}
	return j.Printer.Name
}
