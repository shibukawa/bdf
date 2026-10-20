package spool

import (
	"fmt"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

// Compact rewrites a spool's file without the space that removed jobs and
// documents left: the file does not shrink by itself. The spool must be
// closed (a running server holds the file's lock, and Compact fails). It
// returns the sizes before and after.
//
// This is the one place the spool knows that yuruidb keeps it in bbolt.
func Compact(path string) (before, after int64, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, 0, err
	}
	tmp := path + ".compact"
	if err := compactTo(path, tmp, fi.Mode().Perm()); err != nil {
		os.Remove(tmp)
		return 0, 0, err
	}
	nfi, err := os.Stat(tmp)
	if err == nil {
		// both files are closed: Windows renames over a file no one has
		// open
		err = os.Rename(tmp, path)
	}
	if err != nil {
		os.Remove(tmp)
		return 0, 0, err
	}
	return fi.Size(), nfi.Size(), nil
}

func compactTo(path, tmp string, perm os.FileMode) error {
	src, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: time.Second, ReadOnly: true})
	if err != nil {
		return fmt.Errorf("spool: %s is in use or not a spool: %w", path, err)
	}
	defer src.Close()
	os.Remove(tmp)
	dst, err := bolt.Open(tmp, perm, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return err
	}
	if err := bolt.Compact(dst, src, 64<<20); err != nil {
		dst.Close()
		return err
	}
	if err := dst.Sync(); err != nil {
		dst.Close()
		return err
	}
	return dst.Close()
}
