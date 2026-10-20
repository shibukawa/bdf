//go:build !(darwin || linux)

package raw_test

import (
	"os"
	"testing"
)

func openPTY(t *testing.T) (*os.File, string) {
	t.Skip("pseudo-terminals are tested on macOS and Linux")
	return nil, ""
}
