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
