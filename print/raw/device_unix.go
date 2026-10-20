//go:build unix

package raw

import (
	"fmt"
	"runtime"
	"syscall"
)

// nonblock opens a device file so that its writes can be bounded (the
// runtime's poller takes it).
const nonblock = syscall.O_NONBLOCK

// usbHere tells whether USB printers have device files here: macOS keeps
// them to itself.
func usbHere() error {
	if runtime.GOOS == "darwin" || runtime.GOOS == "ios" {
		return fmt.Errorf("%w: macOS gives a USB printer no device file; print to its CUPS queue (ipp://localhost:631/printers/<name>, which takes the printer's language as it is)", ErrNotHere)
	}
	return nil
}
