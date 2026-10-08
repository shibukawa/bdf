// Package isobmff reads the boxes of the ISO base media file format (ISO/IEC
// 14496-12): the container of AVIF and HEIF images, of MP4 and M4A audio,
// and of MP4 and QuickTime video. It only splits a byte string into boxes
// and walks the top-level boxes of a file; what each box means is the
// business of the reader of the format.
package isobmff

import (
	"encoding/binary"
	"io"
)

// Box is a box: its type and contents (without the header).
type Box struct {
	Type string
	Data []byte
}

// Boxes splits b into boxes. A box that claims more than b holds ends the
// list.
func Boxes(b []byte) []Box {
	var out []Box
	for len(b) >= 8 {
		size, hdr := uint64(binary.BigEndian.Uint32(b)), uint64(8)
		typ := string(b[4:8])
		switch size {
		case 0: // to the end
			size = uint64(len(b))
		case 1:
			if len(b) < 16 {
				return out
			}
			size, hdr = binary.BigEndian.Uint64(b[8:]), 16
		}
		if size < hdr || size > uint64(len(b)) {
			return out
		}
		out = append(out, Box{typ, b[hdr:size]})
		b = b[size:]
	}
	return out
}

// Find returns the contents of the first box of a type, or nil.
func Find(bs []Box, typ string) []byte {
	for _, b := range bs {
		if b.Type == typ {
			return b.Data
		}
	}
	return nil
}

// FindAll returns the contents of every box of a type.
func FindAll(bs []Box, typ string) [][]byte {
	var out [][]byte
	for _, b := range bs {
		if b.Type == typ {
			out = append(out, b.Data)
		}
	}
	return out
}

// Top is a top-level box of a file: where it is, not its contents, so that
// a file is walked without reading its media data.
type Top struct {
	Type string
	// Off is the offset of the box in the file, Header the length of its
	// header and Size its whole length, header included.
	Off, Header, Size int64
}

// Contents reads the contents of the box from r.
func (t Top) Contents(r io.ReaderAt) ([]byte, error) {
	b := make([]byte, t.Size-t.Header)
	if _, err := r.ReadAt(b, t.Off+t.Header); err != nil && err != io.EOF {
		return nil, err
	}
	return b, nil
}

// Walk lists the top-level boxes of a file of size bytes, reading only
// their headers. A box that claims more than the file holds, or that is
// smaller than its header, ends the list.
func Walk(r io.ReaderAt, size int64) []Top {
	var out []Top
	var hdr [16]byte
	for off := int64(0); off+8 <= size; {
		n, err := r.ReadAt(hdr[:], off)
		if n < 8 && err != nil {
			return out
		}
		box := Top{Type: string(hdr[4:8]), Off: off, Header: 8, Size: int64(binary.BigEndian.Uint32(hdr[:]))}
		switch box.Size {
		case 0:
			box.Size = size - off
		case 1:
			if n < 16 {
				return out
			}
			large := binary.BigEndian.Uint64(hdr[8:])
			if large > uint64(size-off) {
				return out
			}
			box.Size, box.Header = int64(large), 16
		}
		if box.Size < box.Header || box.Size > size-off {
			return out
		}
		out = append(out, box)
		off += box.Size
	}
	return out
}
