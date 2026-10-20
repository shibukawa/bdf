//go:build darwin || linux

package raw_test

import (
	"os"
	"testing"

	"golang.org/x/sys/unix"
)

// openPTY opens a pseudo-terminal: its far end, and the name of the
// terminal a serial printer would be.
func openPTY(t *testing.T) (*os.File, string) {
	t.Helper()
	far, err := os.OpenFile("/dev/ptmx", os.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		t.Skip("no pseudo-terminals:", err)
	}
	fd := int(far.Fd())
	name, err := ptyName(fd)
	if err != nil {
		t.Fatal(err)
	}
	// the far end raw too, so that it reads every byte as it was sent
	tio, err := unix.IoctlGetTermios(fd, getTermios)
	if err == nil {
		tio.Iflag, tio.Oflag, tio.Lflag = 0, 0, 0
		unix.IoctlSetTermios(fd, setTermios, tio)
	}
	return far, name
}
