// Package otf opens OpenType font files, TrueType-flavored (.ttf, .ttc) and
// CFF-flavored (.otf, .otc) alike, for reading with little memory.
//
// A font file is mostly glyph data. A Japanese font holds twenty thousand
// glyphs or more, of which a document uses a few hundred, and a collection
// holds several fonts of which one is wanted. Reading such a file into
// memory to lay a page out keeps all of it resident. Open maps the file
// into memory instead, where the platform allows (Unix and Windows): the
// bytes look like a slice of the whole file, but only the pages that are
// read take memory, the table directory, the small tables and the glyphs
// in use, and the kernel may drop them again under pressure since the file
// backs them. Elsewhere, and when a file cannot be mapped, the file is read
// whole.
//
// The bytes are read-only. A mapped file that another process truncates
// while it is mapped faults when the truncated pages are read. Fonts in
// system and font directories are not rewritten in place, which is the use
// this package is for.
package otf

import (
	"errors"
	"io"
	"io/fs"
	"os"
)

// File is an open font file.
type File struct {
	data   []byte
	mapped bool
}

// maxInt is the largest slice the platform addresses.
const maxInt = int64(^uint(0) >> 1)

var errTooLarge = errors.New("otf: the file is too large to address")

// Open opens the font file at path: mapped into memory on Unix and
// Windows, read whole elsewhere and when the file cannot be mapped.
func Open(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return fromOSFile(f)
}

// OpenFS is Open for the file name of fsys. The file is mapped when fsys
// opens it as an *os.File (os.DirFS does) and read whole otherwise (fonts
// embedded in the program, say).
func OpenFS(fsys fs.FS, name string) (*File, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if osf, ok := f.(*os.File); ok {
		return fromOSFile(osf)
	}
	data, err := readAll(f)
	if err != nil {
		return nil, err
	}
	return &File{data: data}, nil
}

// fromOSFile maps f when it is a regular file that is not empty, and
// reads it otherwise. The result outlives f, which the caller closes.
func fromOSFile(f *os.File) (*File, error) {
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if st.Mode().IsRegular() && st.Size() > 0 {
		if st.Size() > maxInt {
			return nil, errTooLarge
		}
		if data, err := mapFile(f, st.Size()); err == nil {
			return &File{data: data, mapped: true}, nil
		}
	}
	data, err := readAll(f)
	if err != nil {
		return nil, err
	}
	return &File{data: data}, nil
}

// readAll reads a file whole, in one allocation of the size Stat reports
// when it reports one.
func readAll(f fs.File) ([]byte, error) {
	size := int64(-1)
	if st, err := f.Stat(); err == nil && st.Mode().IsRegular() {
		size = st.Size()
	}
	if size < 0 || size > maxInt {
		return io.ReadAll(f)
	}
	buf := make([]byte, size)
	n, err := io.ReadFull(f, buf)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return buf[:n], nil // shorter than Stat said
	}
	if err != nil {
		return nil, err
	}
	rest, err := io.ReadAll(f) // longer than Stat said
	if err != nil {
		return nil, err
	}
	return append(buf, rest...), nil
}

// Bytes returns the contents of the file. Close invalidates them.
func (f *File) Bytes() []byte { return f.data }

// Size returns the size of the file in bytes.
func (f *File) Size() int { return len(f.data) }

// Mapped reports whether the file is mapped into memory rather than read.
func (f *File) Mapped() bool { return f.mapped }

// Close releases the mapping, or the bytes. The slice Bytes returned, and
// whatever points into it, must not be used afterwards: a mapped one
// faults. Closing a closed File does nothing.
func (f *File) Close() error {
	data, mapped := f.data, f.mapped
	f.data, f.mapped = nil, false
	if mapped && len(data) > 0 {
		return unmap(data)
	}
	return nil
}
