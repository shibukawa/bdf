package spool

import (
	"fmt"
	"time"
)

// Holds of job-hold-until, as CUPS defines them. A job held for a time of
// day is released when the time comes, by the server's local time.
var holdWindows = map[string]struct{ from, to int }{ // hours, to past midnight when to <= from
	"day-time":     {6, 18},
	"evening":      {18, 6},
	"night":        {18, 6},
	"second-shift": {16, 24},
	"third-shift":  {0, 8},
}

// checkHold refuses a job-hold-until the spool does not know.
func checkHold(hold string) error {
	switch hold {
	case "", "no-hold", "indefinite", "weekend":
		return nil
	}
	if _, ok := holdWindows[hold]; ok {
		return nil
	}
	return fmt.Errorf("%w: job-hold-until %q", ErrInvalid, hold)
}

// holdEnd returns when a job held as hold says is released, from now: now
// itself when the hold is over or there is none, and the zero time for
// "indefinite" (released by hand). at is job-hold-until-time, which wins
// when set.
func holdEnd(hold string, at time.Time, now time.Time) time.Time {
	if !at.IsZero() {
		if at.After(now) {
			return at
		}
		return now
	}
	switch hold {
	case "", "no-hold":
		return now
	case "indefinite":
		return time.Time{}
	case "weekend":
		day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		for day.Weekday() != time.Saturday && day.Weekday() != time.Sunday {
			day = day.AddDate(0, 0, 1)
		}
		if day.After(now) {
			return day
		}
		return now
	}
	w, ok := holdWindows[hold]
	if !ok {
		return now
	}
	h := now.Hour()
	in := w.from <= h && h < w.to
	if w.to <= w.from {
		in = h >= w.from || h < w.to
	}
	if in {
		return now
	}
	start := time.Date(now.Year(), now.Month(), now.Day(), w.from, 0, 0, 0, now.Location())
	if !start.After(now) {
		start = start.AddDate(0, 0, 1)
	}
	return start
}
