package ipp

import (
	"context"
	"testing"
	"time"
)

// TestBrowseRuns checks that Browse opens its sockets and returns in
// time on this host; it finds what the network has, usually nothing.
func TestBrowseRuns(t *testing.T) {
	start := time.Now()
	svcs, err := Browse(context.Background(), 300*time.Millisecond)
	if err != nil {
		t.Skipf("no multicast here: %v", err)
	}
	if d := time.Since(start); d > 2*time.Second {
		t.Fatalf("Browse took %v", d)
	}
	t.Logf("%d printers", len(svcs))
}
