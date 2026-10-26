package web

import (
	"context"
	"net/http"

	"github.com/shibukawa/popcornweb/pwruntime"
)

// livePath is what the topology page polls, every few seconds while it is
// shown (POLL_MS of topology.js).
const livePath = "/api/topology/live"

// QuietPolls keeps the polls of the topology page out of the access log:
// an answer that tells the page what it asked (200, or 304: nothing new)
// would be a line every few seconds for every page open, and the log
// would show nothing else. Another answer (a 401, a 5xx) is logged as
// every request is. It runs before the access log
// (pw.SlotAccessLog - 5), which writes through the backend it puts in
// the request's context.
func QuietPolls(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == livePath {
			to := pwruntime.DeriveResources(r.Context()).Log
			r = r.WithContext(pwruntime.WithLogBackend(r.Context(), pwruntime.NewLogBackend(to.Minimum(), quietPoll{to})))
		}
		next.ServeHTTP(w, r)
	})
}

// quietPoll passes on the records of a poll but its access log's.
type quietPoll struct{ to *pwruntime.LogBackend }

func (s quietPoll) Emit(ctx context.Context, rec pwruntime.Record) {
	if rec.Message == "request completed" && routineStatus(rec.Attributes) {
		return
	}
	pwruntime.NewLogger(ctx, s.to).Log(ctx, rec.Level, rec.Message, rec.Attributes...)
}

func routineStatus(attrs []pwruntime.Attribute) bool {
	for _, a := range attrs {
		if a.Key == "status" {
			st, _ := a.Value.AsInt64()
			return st == http.StatusOK || st == http.StatusNotModified
		}
	}
	return false
}
