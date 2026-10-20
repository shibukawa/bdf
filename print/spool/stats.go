package spool

import (
	"maps"
	"os"
)

// Stats count what the spool did since it opened, and what it holds: a
// server's metrics.
type Stats struct {
	// Done counts the jobs done, by queue and state.
	Done map[QueueState]int64
	// Sheets counts the sheets of the jobs done, by queue.
	Sheets map[string]int64
	// SentBytes counts the bytes printers took, and Retries the jobs put
	// back because a printer did not take them, by printer.
	SentBytes map[string]int64
	Retries   map[string]int64
	// Bytes is the size of the documents of the jobs not done; FileSize
	// the size of the spool's file.
	Bytes    int64
	FileSize int64
}

// QueueState is a queue and a job state.
type QueueState struct {
	Queue, State string
}

func newStats() Stats {
	return Stats{Done: map[QueueState]int64{}, Sheets: map[string]int64{}, SentBytes: map[string]int64{}, Retries: map[string]int64{}}
}

// Stats returns the spool's counts.
func (s *Spool) Stats() Stats {
	s.mu.Lock()
	st := Stats{Done: maps.Clone(s.stats.Done), Sheets: maps.Clone(s.stats.Sheets), SentBytes: maps.Clone(s.stats.SentBytes),
		Retries: maps.Clone(s.stats.Retries)}
	for _, e := range s.jobs {
		st.Bytes += e.job.Document.Size
		if e.job.Document.BDF != nil {
			st.Bytes += e.job.Document.BDF.Size
		}
	}
	s.mu.Unlock()
	if fi, err := os.Stat(s.path); err == nil {
		st.FileSize = fi.Size()
	}
	return st
}

// count changes the counts.
func (s *Spool) count(f func(st *Stats)) {
	s.mu.Lock()
	f(&s.stats)
	s.mu.Unlock()
}
