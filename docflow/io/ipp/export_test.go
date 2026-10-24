package ipp

import (
	"context"
	"testing"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// SetBrowse makes f what Bonjour finds, until the test ends (and forgets
// the printers found again).
func SetBrowse(t testing.TB, f func(ctx context.Context) ([]ipp.Service, error)) {
	old := browse
	browse = f
	t.Cleanup(func() {
		browse = old
		moved.Clear()
	})
}

// PrinterUUID is the uuid a CUPS server node is advertised with.
var PrinterUUID = printerUUID

// SetCreatedTimeout makes d how long a job Create-Job made waits for its
// document, until the test ends.
func SetCreatedTimeout(t testing.TB, d time.Duration) {
	old := createdTimeout
	createdTimeout = d
	t.Cleanup(func() { createdTimeout = old })
}
