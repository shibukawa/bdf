package raw

import (
	"fmt"

	"golang.org/x/sys/unix"
)

const (
	getTermios = unix.TCGETS
	setTermios = unix.TCSETS
)

var speeds = map[int]uint32{
	1200: unix.B1200, 2400: unix.B2400, 4800: unix.B4800, 9600: unix.B9600, 19200: unix.B19200,
	38400: unix.B38400, 57600: unix.B57600, 115200: unix.B115200, 230400: unix.B230400,
}

func setSpeed(t *unix.Termios, baud int) error {
	b, ok := speeds[baud]
	if !ok {
		return fmt.Errorf("%d bauds", baud)
	}
	t.Cflag &^= unix.CBAUD
	t.Cflag |= b
	t.Ispeed, t.Ospeed = b, b
	return nil
}

// drain waits until what was written to the port is sent (tcdrain).
func drain(fd int) error { return unix.IoctlSetInt(fd, unix.TCSBRK, 1) }

// discard drops what was written to the port and not sent (tcflush).
func discard(fd int) error { return unix.IoctlSetInt(fd, unix.TCFLSH, unix.TCOFLUSH) }
