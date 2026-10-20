//go:build !(darwin || linux)

package fakeprinter

import (
	"errors"
	"os"
)

// OpenPTY opens a pseudo-terminal; there are none here.
func OpenPTY() (far *os.File, name string, err error) {
	return nil, "", errors.New("fakeprinter: pseudo-terminals are made on macOS and Linux")
}
