//go:build !unix && !windows

package otf

import (
	"errors"
	"os"
)

var errNoMap = errors.New("otf: memory mapping is not available on this platform")

func mapFile(f *os.File, size int64) ([]byte, error) { return nil, errNoMap }

func unmap(b []byte) error { return nil }
