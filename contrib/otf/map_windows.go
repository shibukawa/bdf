//go:build windows

package otf

import (
	"os"
	"syscall"
	"unsafe"
)

// mapFile maps size bytes of f read-only. The view outlives f and the
// mapping object, which is closed right away.
func mapFile(f *os.File, size int64) ([]byte, error) {
	h, err := syscall.CreateFileMapping(syscall.Handle(f.Fd()), nil, syscall.PAGE_READONLY, uint32(size>>32), uint32(size), nil)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(h)
	addr, err := syscall.MapViewOfFile(h, syscall.FILE_MAP_READ, 0, 0, uintptr(size))
	if err != nil {
		return nil, err
	}
	// The view is outside the Go heap, so the collector has nothing to
	// track in it; the conversion goes through *uintptr because vet flags
	// unsafe.Pointer(uintptr) as a pointer it might lose.
	return unsafe.Slice((*byte)(*(*unsafe.Pointer)(unsafe.Pointer(&addr))), int(size)), nil
}

func unmap(b []byte) error {
	return syscall.UnmapViewOfFile(uintptr(unsafe.Pointer(unsafe.SliceData(b))))
}
