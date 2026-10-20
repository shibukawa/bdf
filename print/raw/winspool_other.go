//go:build !windows

package raw

import "fmt"

func checkSpool(name string) error {
	return fmt.Errorf("%w: a queue of Windows (winspool://) is printed to on Windows", ErrNotHere)
}

func openSpool(name, job string) (device, error) { return nil, checkSpool(name) }
