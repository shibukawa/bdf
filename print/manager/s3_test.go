package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"
)

// TestS3 runs a location on S3-compatible storage, as a serverless
// deployment keeps it: the records and the documents in buckets, every
// request on a manager of its own, ticks at once. It runs when
// BDF_PM_S3_ENDPOINT names the server (AWS_ACCESS_KEY_ID and
// AWS_SECRET_ACCESS_KEY its keys), which has the buckets pm-records and
// pm-docs.
func TestS3(t *testing.T) {
	ep := os.Getenv("BDF_PM_S3_ENDPOINT")
	if ep == "" {
		t.Skip("BDF_PM_S3_ENDPOINT is not set")
	}
	q := "?region=us-east-1&endpoint=" + url.QueryEscape(ep)
	e := newEnvAt(t, "s3://pm-records"+q, "s3://pm-docs"+q, nil)
	if !e.st.Shared {
		t.Fatal("S3 is not shared")
	}
	// the buckets outlive a run: a location of its own
	loc := fmt.Sprintf("s3-%d", time.Now().UnixNano())
	p := e.printer("laser-" + loc)
	e.admin("PUT", "/api/locations/"+loc, Location{Fallback: &Fallback{Printer: "laser-" + loc, After: Duration(10 * time.Minute)}}, http.StatusOK)
	d := e.enroll(loc, "tablet")
	if b := d.board(); b == nil || len(b.Open) != 0 {
		t.Fatalf("board %+v", b)
	}
	for i := range 3 {
		e.printJob(loc, fmt.Sprint("slip ", i), 1)
	}
	b := d.board()
	if b == nil || len(b.Open) != 3 {
		t.Fatalf("board %+v", b)
	}
	d.action(b.Open[0].ID, "done", http.StatusOK)

	// nobody opens the other two: ticks at once print each once
	e.clock.add(11 * time.Minute)
	d.login()
	d.board()
	var wg sync.WaitGroup
	for range 6 {
		wg.Go(func() {
			if _, err := e.m().Tick(context.Background()); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if n := len(p.Jobs()); n != 2 {
		t.Fatalf("%d jobs printed; want 2", n)
	}
	var rep TickReport
	json.Unmarshal(e.admin("POST", "/api/tick", nil, http.StatusOK), &rep)
	if rep.Printed != 0 {
		t.Fatalf("printed again: %+v", rep)
	}
	e.admin("DELETE", "/api/locations/"+loc, nil, http.StatusNoContent)
	e.admin("DELETE", "/api/printers/laser-"+loc, nil, http.StatusNoContent)
}
