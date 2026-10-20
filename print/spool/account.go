package spool

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Accounting writes the jobs done between from and to (by when they were
// done; a zero time does not bound) as JSON Lines ("jsonl") or CSV
// ("csv"), oldest first: who printed what, where, and how much.
func (s *Spool) Accounting(ctx context.Context, from, to time.Time, format string, w io.Writer) error {
	var write func(Job) error
	var flush func() error
	switch format {
	case "", "jsonl":
		enc := json.NewEncoder(w)
		write = func(j Job) error { return enc.Encode(accountOf(j)) }
		flush = func() error { return nil }
	case "csv":
		cw := csv.NewWriter(w)
		if err := cw.Write(accountColumns); err != nil {
			return err
		}
		write = func(j Job) error { return cw.Write(accountOf(j).row()) }
		flush = func() error { cw.Flush(); return cw.Error() }
	default:
		return fmt.Errorf("%w: accounting format %q: jsonl or csv", ErrInvalid, format)
	}
	var werr error
	err := s.st.eachDone(ctx, func(j Job, _ string) bool {
		if (!from.IsZero() && j.Finished.Before(from)) || (!to.IsZero() && !j.Finished.Before(to)) {
			return true
		}
		if werr = write(j); werr != nil {
			return false
		}
		return ctx.Err() == nil
	})
	if err != nil {
		return err
	}
	if werr != nil {
		return werr
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return flush()
}

// account is one line of the accounting.
type account struct {
	ID          int32     `json:"id"`
	UUID        string    `json:"uuid"`
	Queue       string    `json:"queue"`
	Printer     string    `json:"printer,omitempty"`
	User        string    `json:"user"`
	Name        string    `json:"name"`
	State       string    `json:"state"`
	Reasons     string    `json:"reasons,omitempty"`
	Submitted   time.Time `json:"submitted"`
	Finished    time.Time `json:"finished"`
	Sheets      int       `json:"sheets"`
	Impressions int       `json:"impressions"`
	Media       string    `json:"media,omitempty"`
	Sides       string    `json:"sides,omitempty"`
	ColorMode   string    `json:"print-color-mode,omitempty"`
	Copies      int       `json:"copies"`
	Format      string    `json:"format,omitempty"`
	Bytes       int64     `json:"bytes"`
	Attempts    int       `json:"attempts"`
}

var accountColumns = []string{"id", "uuid", "queue", "printer", "user", "name", "state", "reasons", "submitted", "finished",
	"sheets", "impressions", "media", "sides", "print-color-mode", "copies", "format", "bytes", "attempts"}

func accountOf(j Job) account {
	a := account{ID: j.ID, UUID: j.UUID, Queue: j.Queue, User: j.User, Name: j.Name, State: j.State.String(), Reasons: strings.Join(j.Reasons, " "),
		Submitted: j.Submitted, Finished: j.Finished, Sheets: j.Sheets, Impressions: j.Impressions, Media: j.Media, Sides: j.Ticket.Sides,
		ColorMode: j.Ticket.ColorMode, Copies: max(j.Ticket.Copies, 1), Format: j.Format, Bytes: j.Bytes, Attempts: j.Attempts}
	if j.Printer != nil {
		a.Printer = j.Printer.Name
	}
	return a
}

func (a account) row() []string {
	return []string{strconv.Itoa(int(a.ID)), a.UUID, a.Queue, a.Printer, a.User, a.Name, a.State, a.Reasons,
		a.Submitted.Format(time.RFC3339), a.Finished.Format(time.RFC3339), strconv.Itoa(a.Sheets), strconv.Itoa(a.Impressions),
		a.Media, a.Sides, a.ColorMode, strconv.Itoa(a.Copies), a.Format, strconv.FormatInt(a.Bytes, 10), strconv.Itoa(a.Attempts)}
}
