//go:build !unix && !windows

package raw

import "fmt"

const nonblock = 0

func usbHere() error { return fmt.Errorf("%w: no device files of USB printers", ErrNotHere) }
