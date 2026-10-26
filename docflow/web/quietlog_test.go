package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/shibukawa/popcornweb/middlewares"
	"github.com/shibukawa/popcornweb/pwruntime"
)

type records struct {
	mu   sync.Mutex
	recs []pwruntime.Record
}

func (s *records) Emit(_ context.Context, r pwruntime.Record) {
	s.mu.Lock()
	s.recs = append(s.recs, r)
	s.mu.Unlock()
}

// The topology page's polls answered 200 or 304 leave no access log;
// their other answers, the other requests and what the handler logs do.
func TestQuietPolls(t *testing.T) {
	sink := &records{}
	backend := pwruntime.NewLogBackend(pwruntime.LevelInfo, sink)
	status := http.StatusNotModified
	h := QuietPolls(middlewares.Track(middlewares.AccessLog()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status == http.StatusInternalServerError {
			pwruntime.ReadLogger(r.Context()).Error("live failed")
		}
		w.WriteHeader(status)
	}))))
	do := func(method, path string, st int) {
		status = st
		r := httptest.NewRequest(method, path, nil)
		h.ServeHTTP(httptest.NewRecorder(), r.WithContext(pwruntime.WithLogBackend(r.Context(), backend)))
	}
	do("GET", livePath, http.StatusNotModified)
	do("GET", livePath, http.StatusOK)
	do("GET", "/api/topology", http.StatusOK)
	do("GET", livePath, http.StatusUnauthorized)
	do("GET", livePath, http.StatusInternalServerError)
	var got []string
	for _, r := range sink.recs {
		line := r.Message
		for _, a := range r.Attributes {
			if a.Key == "path" {
				p, _ := a.Value.AsString()
				line += " " + p
			}
			if a.Key == "status" {
				st, _ := a.Value.AsInt64()
				line += " " + http.StatusText(int(st))
			}
		}
		got = append(got, line)
	}
	want := []string{
		"request completed /api/topology OK",
		"request completed " + livePath + " Unauthorized",
		"live failed",
		"request completed " + livePath + " Internal Server Error",
	}
	if len(got) != len(want) {
		t.Fatalf("records %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("record %d: %q, want %q", i, got[i], want[i])
		}
	}
}
