//go:build bdf_noconv

package imgconv

import "image"

// Available reports whether the WebP codec is compiled in.
func Available() bool { return false }

func encodeWebP(*image.NRGBA, int, int, bool) ([]byte, error) { return nil, ErrNotAvailable }
