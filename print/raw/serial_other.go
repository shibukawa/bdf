//go:build !(linux || darwin || freebsd || netbsd || openbsd || dragonfly || windows)

package raw

import "fmt"

func openSerial(path string, s Serial) (device, error) {
	return nil, fmt.Errorf("%w: no serial ports", ErrNotHere)
}
