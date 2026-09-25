//go:build bdf_noconv

package woff2

// Available reports whether the Brotli encoder is compiled in.
func Available() bool { return false }

func compress([]byte) ([]byte, error) { return nil, ErrNotAvailable }
