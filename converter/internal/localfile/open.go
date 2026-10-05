//go:build !tinygo

// Package localfile opens document resources within a directory.
package localfile

import (
	"io"
	"os"
)

// Open resolves each component within dir, including symlinks, without a
// check-then-open race. A reference that leaves dir is rejected.
func Open(dir, name string) (io.ReadCloser, error) {
	return os.OpenInRoot(dir, name)
}
