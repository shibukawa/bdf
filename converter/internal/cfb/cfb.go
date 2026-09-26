// Package cfb reads Compound File Binary files ([MS-CFB]): the container of
// legacy Office documents, and of encrypted Office Open XML documents, which
// keep their encrypted package in a stream of one.
//
// It reads the streams directly under the root storage, which is all the
// encrypted documents need.
package cfb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf16"
)

// Signature starts every compound file.
var Signature = []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}

// ErrNotExist is returned for a stream the file does not have.
var ErrNotExist = errors.New("cfb: no such stream")

const (
	maxRegSect = 0xfffffffa
	endOfChain = 0xfffffffe
	noStream   = 0xffffffff

	typeStorage = 1
	typeStream  = 2
	typeRoot    = 5
)

// File is an open compound file.
type File struct {
	r          io.ReaderAt
	size       int64
	sectorSize int64
	cutoff     uint64
	fat        []uint32
	miniFAT    []uint32
	dir        []entry
	mini       []byte // the mini stream, read on first use
}

type entry struct {
	name               string
	typ                byte
	left, right, child uint32
	start              uint32
	size               uint64
}

func errorf(format string, args ...any) error { return fmt.Errorf("cfb: "+format, args...) }

// Open reads the header, the allocation tables and the directory.
func Open(r io.ReaderAt, size int64) (*File, error) {
	h := make([]byte, 512)
	if _, err := r.ReadAt(h, 0); err != nil {
		return nil, errorf("reading header: %v", err)
	}
	if !bytes.Equal(h[:8], Signature) {
		return nil, errorf("not a compound file")
	}
	le := binary.LittleEndian
	major := le.Uint16(h[26:])
	shift := le.Uint16(h[30:])
	switch {
	case major == 3 && shift == 9, major == 4 && shift == 12:
	default:
		return nil, errorf("unsupported version %d with sector shift %d", major, shift)
	}
	if le.Uint16(h[28:]) != 0xfffe || le.Uint16(h[32:]) != 6 {
		return nil, errorf("bad header")
	}
	f := &File{r: r, size: size, sectorSize: 1 << shift, cutoff: uint64(le.Uint32(h[56:]))}
	sectors := (size + f.sectorSize - 1) / f.sectorSize // bounds every table and chain

	// The FAT sectors: 109 in the header, the rest in the DIFAT chain.
	nFAT := int64(le.Uint32(h[44:]))
	if nFAT > sectors {
		return nil, errorf("bad FAT sector count")
	}
	var fatSects []uint32
	for i := range 109 {
		if int64(len(fatSects)) == nFAT {
			break
		}
		fatSects = append(fatSects, le.Uint32(h[76+4*i:]))
	}
	next := le.Uint32(h[68:])
	per := int(f.sectorSize/4) - 1
	for n := int64(0); int64(len(fatSects)) < nFAT; n++ {
		if next > maxRegSect || n > sectors {
			return nil, errorf("bad DIFAT chain")
		}
		b, err := f.sector(next)
		if err != nil {
			return nil, err
		}
		for i := 0; i < per && int64(len(fatSects)) < nFAT; i++ {
			fatSects = append(fatSects, le.Uint32(b[4*i:]))
		}
		next = le.Uint32(b[4*per:])
	}
	for _, s := range fatSects {
		b, err := f.sector(s)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(b); i += 4 {
			f.fat = append(f.fat, le.Uint32(b[i:]))
		}
	}

	d, err := f.chain(f.fat, le.Uint32(h[48:]), -1, f.sectorSize, f.sector)
	if err != nil {
		return nil, errorf("directory: %v", err)
	}
	for i := 0; i+128 <= len(d); i += 128 {
		e := d[i : i+128]
		n := min(int(le.Uint16(e[64:])), 64)
		u := make([]uint16, 0, 32)
		for j := 0; j+1 < n; j += 2 {
			u = append(u, le.Uint16(e[j:]))
		}
		for len(u) > 0 && u[len(u)-1] == 0 {
			u = u[:len(u)-1]
		}
		sz := le.Uint64(e[120:])
		if major == 3 {
			sz &= 0xffffffff // the high half is undefined in version 3
		}
		f.dir = append(f.dir, entry{name: string(utf16.Decode(u)), typ: e[66],
			left: le.Uint32(e[68:]), right: le.Uint32(e[72:]), child: le.Uint32(e[76:]),
			start: le.Uint32(e[116:]), size: sz})
	}
	if len(f.dir) == 0 || f.dir[0].typ != typeRoot {
		return nil, errorf("no root storage")
	}

	if n := le.Uint32(h[64:]); n > 0 {
		b, err := f.chain(f.fat, le.Uint32(h[60:]), -1, f.sectorSize, f.sector)
		if err != nil {
			return nil, errorf("mini FAT: %v", err)
		}
		for i := 0; i+4 <= len(b); i += 4 {
			f.miniFAT = append(f.miniFAT, le.Uint32(b[i:]))
		}
	}
	return f, nil
}

// sector returns the bytes of a regular sector.
func (f *File) sector(n uint32) ([]byte, error) {
	off := (int64(n) + 1) * f.sectorSize
	if n > maxRegSect || off+f.sectorSize > f.size {
		return nil, errorf("sector %d out of range", n)
	}
	b := make([]byte, f.sectorSize)
	if _, err := f.r.ReadAt(b, off); err != nil {
		return nil, err
	}
	return b, nil
}

// miniSector returns the bytes of a sector of the mini stream.
func (f *File) miniSector(n uint32) ([]byte, error) {
	if f.mini == nil {
		root := f.dir[0]
		if root.size > uint64(f.size) {
			return nil, errorf("mini stream too large")
		}
		b, err := f.chain(f.fat, root.start, int64(root.size), f.sectorSize, f.sector)
		if err != nil {
			return nil, errorf("mini stream: %v", err)
		}
		f.mini = b
	}
	off := int64(n) * 64
	if off+64 > int64(len(f.mini)) {
		return nil, errorf("mini sector %d out of range", n)
	}
	return f.mini[off : off+64], nil
}

// chain reads the sectors of a chain, up to size bytes (-1: the whole chain).
func (f *File) chain(table []uint32, start uint32, size int64, sectorSize int64, read func(uint32) ([]byte, error)) ([]byte, error) {
	var out []byte
	if size > 0 {
		out = make([]byte, 0, size)
	}
	for s, n := start, 0; s != endOfChain; n++ {
		if size >= 0 && int64(len(out)) >= size {
			break
		}
		if int(s) >= len(table) || n > len(table) {
			return nil, errorf("bad sector chain")
		}
		b, err := read(s)
		if err != nil {
			return nil, err
		}
		out = append(out, b...)
		s = table[s]
	}
	if size >= 0 {
		if int64(len(out)) < size {
			return nil, errorf("stream shorter than its size")
		}
		out = out[:size]
	}
	return out, nil
}

// Stream returns a stream directly under the root storage. Names compare
// case-insensitively, as in the format.
func (f *File) Stream(name string) ([]byte, error) {
	e, ok := f.find(name)
	if !ok || e.typ != typeStream {
		return nil, ErrNotExist
	}
	if e.size < f.cutoff {
		return f.chain(f.miniFAT, e.start, int64(e.size), 64, f.miniSector)
	}
	if e.size > uint64(f.size) {
		return nil, errorf("stream %q larger than the file", name)
	}
	return f.chain(f.fat, e.start, int64(e.size), f.sectorSize, f.sector)
}

// Has reports whether the root storage holds a stream or storage by name.
func (f *File) Has(name string) bool {
	_, ok := f.find(name)
	return ok
}

// find walks the root storage's tree of children.
func (f *File) find(name string) (entry, bool) {
	seen := map[uint32]bool{}
	stack := []uint32{f.dir[0].child}
	for len(stack) > 0 {
		i := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if i == noStream || int(i) >= len(f.dir) || seen[i] {
			continue
		}
		seen[i] = true
		e := f.dir[i]
		if (e.typ == typeStream || e.typ == typeStorage) && strings.EqualFold(e.name, name) {
			return e, true
		}
		stack = append(stack, e.left, e.right)
	}
	return entry{}, false
}
