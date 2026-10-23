package ipp

import (
	"context"
	"testing"

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
