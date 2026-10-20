//go:build darwin || linux

package fakeprinter

import (
	"os"

	"golang.org/x/sys/unix"
)

// OpenPTY opens a pseudo-terminal: its far end, which a serial printer
// reads (ServeSerial), and the name of the terminal a client opens as the
// printer's serial port. Both ends are raw: every byte arrives as it was
// sent.
func OpenPTY() (far *os.File, name string, err error) {
	far, err = os.OpenFile("/dev/ptmx", os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		return nil, "", err
	}
	fd := int(far.Fd())
	name, err = ptyName(fd)
	if err != nil {
		far.Close()
		return nil, "", err
	}
	if tio, err := unix.IoctlGetTermios(fd, getTermios); err == nil {
		tio.Iflag, tio.Oflag, tio.Lflag = 0, 0, 0
		unix.IoctlSetTermios(fd, setTermios, tio)
	}
	return far, name, nil
}
