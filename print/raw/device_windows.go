package raw

import "fmt"

const nonblock = 0

// usbHere tells whether USB printers have device files here: Windows
// gives them a queue instead.
func usbHere() error {
	return fmt.Errorf("%w: Windows gives a USB printer a queue, not a device file; print to the queue of its driver (winspool://<name>)", ErrNotHere)
}
