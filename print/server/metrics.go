package server

import (
	"bufio"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/spool"
)

// liveStates are the states of jobs not done, which bdf_print_jobs tells
// for each queue (0 too).
var liveStates = []ipp.JobState{ipp.JobPending, ipp.JobPendingHeld, ipp.JobProcessing, ipp.JobProcessingStopped}

// metrics writes the server's metrics in Prometheus' text format. The
// counters count since the server started.
func (s *Server) metrics(w http.ResponseWriter, r *http.Request, c *caller) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	b := bufio.NewWriter(w)
	defer b.Flush()
	family := func(name, typ, help string) {
		fmt.Fprintf(b, "# HELP %s %s\n# TYPE %s %s\n", name, help, name, typ)
	}
	sample := func(name string, v any, labels ...string) {
		b.WriteString(name)
		if len(labels) > 0 {
			b.WriteByte('{')
			for i := 0; i < len(labels); i += 2 {
				if i > 0 {
					b.WriteByte(',')
				}
				fmt.Fprintf(b, "%s=\"%s\"", labels[i], escapeLabel(labels[i+1]))
			}
			b.WriteByte('}')
		}
		fmt.Fprintf(b, " %v\n", v)
	}
	bit := func(ok bool) int {
		if ok {
			return 1
		}
		return 0
	}

	queues := s.sp.Queues()
	stats := s.sp.Stats()
	family("bdf_print_jobs", "gauge", "Jobs not done, by queue and state.")
	for _, q := range queues {
		for _, st := range liveStates {
			sample("bdf_print_jobs", q.Jobs[st.String()], "queue", q.Name, "state", st.String())
		}
	}
	family("bdf_print_jobs_total", "counter", "Jobs done, by queue and state.")
	done := make([]spool.QueueState, 0, len(stats.Done))
	for k := range stats.Done {
		done = append(done, k)
	}
	sort.Slice(done, func(a, b int) bool {
		if done[a].Queue != done[b].Queue {
			return done[a].Queue < done[b].Queue
		}
		return done[a].State < done[b].State
	})
	for _, k := range done {
		sample("bdf_print_jobs_total", stats.Done[k], "queue", k.Queue, "state", k.State)
	}
	family("bdf_print_sheets_total", "counter", "Sheets of the jobs done, by queue.")
	for _, q := range queues {
		sample("bdf_print_sheets_total", stats.Sheets[q.Name], "queue", q.Name)
	}
	family("bdf_print_queue_paused", "gauge", "1 when the queue does not send its jobs.")
	for _, q := range queues {
		sample("bdf_print_queue_paused", bit(q.Paused), "queue", q.Name)
	}
	family("bdf_print_queue_accepting", "gauge", "1 when the queue takes jobs.")
	for _, q := range queues {
		sample("bdf_print_queue_accepting", bit(!q.Rejecting), "queue", q.Name)
	}

	printers := s.sp.Printers()
	family("bdf_print_printer_up", "gauge", "1 when the printer answered last.")
	for _, p := range printers {
		sample("bdf_print_printer_up", bit(p.Up), "printer", p.Name)
	}
	family("bdf_print_printer_state", "gauge", "The printer's printer-state: 3 idle, 4 processing, 5 stopped; 0 before it answered.")
	for _, p := range printers {
		st := 0
		if caps, _ := s.sp.PrinterCapabilities(p.Name); caps != nil {
			st = int(caps.State)
		}
		sample("bdf_print_printer_state", st, "printer", p.Name)
	}
	family("bdf_print_printer_jobs", "gauge", "Jobs sent to the printer and not done.")
	for _, p := range printers {
		sample("bdf_print_printer_jobs", p.Jobs, "printer", p.Name)
	}
	family("bdf_print_sent_bytes_total", "counter", "Bytes the printer took.")
	for _, p := range printers {
		sample("bdf_print_sent_bytes_total", stats.SentBytes[p.Name], "printer", p.Name)
	}
	family("bdf_print_retries_total", "counter", "Jobs put back because the printer did not take them.")
	for _, p := range printers {
		sample("bdf_print_retries_total", stats.Retries[p.Name], "printer", p.Name)
	}

	family("bdf_print_spool_bytes", "gauge", "Bytes of the documents of the jobs not done.")
	sample("bdf_print_spool_bytes", stats.Bytes)
	family("bdf_print_spool_file_bytes", "gauge", "Size of the spool's file.")
	sample("bdf_print_spool_file_bytes", stats.FileSize)
}

func escapeLabel(s string) string {
	return strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`).Replace(s)
}
