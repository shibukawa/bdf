package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// heartbeat is how often a quiet event stream says it is alive, so that
// proxies keep it open.
var heartbeat = 25 * time.Second

// events streams the changes of jobs, queues and printers as Server-Sent
// Events: "job", "queue" and "printer" with the JSON of each. job=<id>
// follows one job. A caller sees the jobs it may see. Changes are dropped
// for a client that does not keep up: a client that comes back reads the
// state again.
func (s *Server) events(w http.ResponseWriter, r *http.Request, c *caller) {
	var only int32
	if v := r.URL.Query().Get("job"); v != "" {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil || n <= 0 {
			badRequest(w, "job: a job id")
			return
		}
		only = int32(n)
	}
	select {
	case s.streams <- struct{}{}:
		defer func() { <-s.streams }()
	default:
		w.Header().Set("Retry-After", "10")
		writeError(w, http.StatusServiceUnavailable, "busy", "the server has as many event streams as it keeps")
		return
	}
	ch, stop := s.sp.Subscribe(64)
	defer stop()
	rc := http.NewResponseController(w)
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "retry: 3000\n\n")
	if rc.Flush() != nil {
		return
	}
	send := func(event string, v any) error {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, b); err != nil {
			return err
		}
		return rc.Flush()
	}
	tick := time.NewTicker(heartbeat)
	defer tick.Stop()
	for {
		var err error
		select {
		case <-r.Context().Done():
			return
		case <-s.done:
			return
		case <-tick.C:
			if _, err = fmt.Fprint(w, ": alive\n\n"); err == nil {
				err = rc.Flush()
			}
		case ev, ok := <-ch:
			switch {
			case !ok:
				return // the spool closed
			case ev.Job != nil:
				if c.sees(ev.Job) && (only == 0 || ev.Job.ID == only) {
					err = send("job", ev.Job)
				}
			case only != 0:
			case ev.Queue != nil:
				err = send("queue", ev.Queue)
			case ev.Printer != nil:
				p := *ev.Printer
				p.URL = redact(p.URL)
				err = send("printer", p)
			}
		}
		if err != nil {
			return
		}
	}
}
