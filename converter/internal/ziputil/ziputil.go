// Package ziputil opens zip archives (Office packages, EPUB publications,
// archives of KiCad, Gerber and SXF files) with the inflate of
// klauspost/compress, which reads a third faster than the standard
// library's and gives the same bytes. The decompressors are pooled, as
// archive/zip pools its own: one holds a 32 KiB window.
package ziputil

import (
	"archive/zip"
	"io"
	"sync"

	"github.com/klauspost/compress/flate"
)

// NewReader is zip.NewReader with the faster inflate registered.
func NewReader(r io.ReaderAt, size int64) (*zip.Reader, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	zr.RegisterDecompressor(zip.Deflate, newInflater)
	return zr, nil
}

var inflaters sync.Pool

// inflater is a pooled flate reader; Close returns it to the pool.
type inflater struct {
	fr io.ReadCloser
}

func newInflater(r io.Reader) io.ReadCloser {
	if reused := inflaters.Get(); reused != nil {
		fr := reused.(io.ReadCloser)
		if err := fr.(flate.Resetter).Reset(r, nil); err == nil {
			return &inflater{fr: fr}
		}
	}
	return &inflater{fr: flate.NewReader(r)}
}

func (i *inflater) Read(p []byte) (int, error) {
	if i.fr == nil {
		return 0, io.ErrClosedPipe
	}
	return i.fr.Read(p)
}

func (i *inflater) Close() error {
	if i.fr == nil {
		return nil
	}
	err := i.fr.Close()
	i.fr.(flate.Resetter).Reset(nil, nil) // not the archive's reader, until the next
	inflaters.Put(i.fr)
	i.fr = nil
	return err
}
