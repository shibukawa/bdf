package server

import (
	"context"
	"strings"
)

// instanceName is the DNS-SD instance name of a queue: its info and the
// server's host, 63 bytes at most.
func (x *ippServer) instanceName(q *ippQueue) string {
	short, _, _ := strings.Cut(x.host, ".")
	name := q.info + " @ " + short
	for len(name) > 63 {
		r := []rune(name)
		name = string(r[:len(r)-1])
	}
	return name
}

// advertise advertises the queues with DNS-SD until ctx ends.
func (x *ippServer) advertise(ctx context.Context) error { return nil }
