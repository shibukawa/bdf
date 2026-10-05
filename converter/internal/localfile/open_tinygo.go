//go:build tinygo

package localfile

import (
	"errors"
	"io"
)

// The browser converter has no local filesystem. TinyGo does not provide
// os.OpenInRoot; reject local resources instead of opening them unconfined.
func Open(dir, name string) (io.ReadCloser, error) {
	return nil, errors.New("local resource files are unavailable in TinyGo; use data URLs or archive resources")
}
