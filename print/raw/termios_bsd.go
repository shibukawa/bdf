//go:build darwin || freebsd || netbsd || openbsd || dragonfly

package raw

import "golang.org/x/sys/unix"

const (
	getTermios = unix.TIOCGETA
	setTermios = unix.TIOCSETA
)

// setSpeed sets the speed: the BSDs keep it in bauds.
func setSpeed(t *unix.Termios, baud int) error {
	set(&t.Ispeed, baud)
	set(&t.Ospeed, baud)
	return nil
}

// set sets a field of the type it has on this system.
func set[T ~int32 | ~uint32 | ~uint64](p *T, v int) { *p = T(v) }

// drain waits until what was written to the port is sent (tcdrain).
func drain(fd int) error {
	_, _, e := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.TIOCDRAIN), 0)
	if e != 0 {
		return e
	}
	return nil
}

// discard drops what was written to the port and not sent (tcflush).
func discard(fd int) error { return unix.IoctlSetPointerInt(fd, unix.TIOCFLUSH, 2) } // FWRITE
