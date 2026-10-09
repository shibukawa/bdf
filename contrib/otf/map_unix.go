//go:build unix

package otf

import (
	"os"
	"syscall"
)

// mapFile maps size bytes of f read-only. The mapping outlives f.
func mapFile(f *os.File, size int64) ([]byte, error) {
	return syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
}

func unmap(b []byte) error { return syscall.Munmap(b) }
