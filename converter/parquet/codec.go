package parquet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/gzip"

	"github.com/shibukawa/tinygodriver/compress/zstd"
)

// maxPageSize bounds the size of a page, compressed or not.
const maxPageSize = 1 << 29

var errCorrupt = errors.New("corrupt compressed data")

var zstdReaders sync.Pool

func decompressZstd(src []byte, size int) ([]byte, error) {
	d, _ := zstdReaders.Get().(*zstd.Reader)
	if d == nil {
		var err error
		d, err = zstd.NewReader(bytes.NewReader(src), zstd.WithMaxWindow(maxPageSize), zstd.WithMaxOutput(maxPageSize))
		if err != nil {
			return nil, err
		}
	} else if err := d.Reset(bytes.NewReader(src)); err != nil {
		zstdReaders.Put(d)
		return nil, err
	}
	defer func() {
		// Do not keep the compressed page alive while the reader is pooled.
		_ = d.Reset(bytes.NewReader(nil))
		zstdReaders.Put(d)
	}()
	return readSized(d, size)
}

// decompress decompresses a page (or the values of a v2 data page) of a
// codec into its size. The size a page says it has is not taken on trust:
// the output grows as the data decompresses, or is bounded by how much
// the codec can expand its input.
func decompress(codec int32, src []byte, size int) ([]byte, error) {
	if size < 0 || size > maxPageSize {
		return nil, fmt.Errorf("a page of %d bytes", size)
	}
	var out []byte
	var err error
	switch codec {
	case codecNone:
		return src, nil
	case codecSnappy:
		// a copy of 64 bytes takes 3
		if n, k := binary.Uvarint(src); k <= 0 || n != uint64(size) || size > 22*len(src) {
			return nil, fmt.Errorf("Snappy: %w", errCorrupt)
		}
		out = make([]byte, size)
		err = snappyBlock(out, src)
	case codecGzip:
		var zr *gzip.Reader
		if zr, err = gzip.NewReader(bytes.NewReader(src)); err == nil {
			out, err = readSized(zr, size)
		}
	case codecBrotli:
		out, err = brotliDecompress(src, size)
	case codecZstd:
		out, err = decompressZstd(src, size)
	case codecLZ4Raw, codecLZ4:
		// a byte of 255 adds 255 to the length of a match
		if size > 256*len(src)+64 {
			return nil, fmt.Errorf("LZ4: %w", errCorrupt)
		}
		if codec == codecLZ4 {
			out, err = lz4Parquet(src, size)
			break
		}
		out = make([]byte, size)
		var n int
		if n, err = lz4Block(out, 0, src); err == nil && n != size {
			err = errCorrupt
		}
	case codecLZO:
		return nil, errors.New("LZO compression is not supported")
	default:
		return nil, fmt.Errorf("unknown compression codec %d", codec)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", codecNames[codec], err)
	}
	if len(out) != size {
		return nil, fmt.Errorf("%s: %w", codecNames[codec], errCorrupt)
	}
	return out, nil
}

// readSized reads a stream that decompresses into size bytes.
func readSized(r io.Reader, size int) ([]byte, error) {
	out, err := io.ReadAll(io.LimitReader(r, int64(size)+1))
	if err != nil {
		return nil, err
	}
	if len(out) != size {
		return nil, errCorrupt
	}
	return out, nil
}

// snappyBlock decompresses a Snappy block (its length, then literals and
// copies of the bytes before) into dst, which is as long as it says.
func snappyBlock(dst, src []byte) error {
	n, k := binary.Uvarint(src)
	if k <= 0 || n != uint64(len(dst)) {
		return errCorrupt
	}
	d, s := 0, k
	for s < len(src) {
		tag := src[s]
		var length, offset int
		switch tag & 3 {
		case 0: // a literal
			length = int(tag >> 2)
			s++
			if length >= 60 {
				nb := length - 59
				if len(src)-s < nb {
					return errCorrupt
				}
				length = int(leUint(src[s:], nb))
				s += nb
			}
			length++
			if length <= 0 || length > len(src)-s || length > len(dst)-d {
				return errCorrupt
			}
			d += copy(dst[d:], src[s:s+length])
			s += length
			continue
		case 1:
			if len(src)-s < 2 {
				return errCorrupt
			}
			length, offset = 4+int(tag>>2&7), int(tag>>5)<<8|int(src[s+1])
			s += 2
		case 2:
			if len(src)-s < 3 {
				return errCorrupt
			}
			length, offset = 1+int(tag>>2), int(binary.LittleEndian.Uint16(src[s+1:]))
			s += 3
		case 3:
			if len(src)-s < 5 {
				return errCorrupt
			}
			length, offset = 1+int(tag>>2), int(binary.LittleEndian.Uint32(src[s+1:]))
			s += 5
		}
		if offset <= 0 || offset > d || length > len(dst)-d {
			return errCorrupt
		}
		if offset >= length {
			d += copy(dst[d:d+length], dst[d-offset:])
		} else {
			for i := range length {
				dst[d+i] = dst[d-offset+i]
			}
			d += length
		}
	}
	if d != len(dst) {
		return errCorrupt
	}
	return nil
}

// lz4Parquet decompresses the deprecated LZ4 codec, which writers took for
// different framings: Hadoop's (blocks each after their decompressed and
// compressed sizes, big-endian), the LZ4 frame format, and a raw block. Like
// Arrow, it tries them in that order.
func lz4Parquet(src []byte, size int) ([]byte, error) {
	out := make([]byte, size)
	if lz4Hadoop(out, src) {
		return out, nil
	}
	if len(src) >= 4 && binary.LittleEndian.Uint32(src) == lz4FrameMagic {
		n, err := lz4Frame(out, src)
		if err != nil || n != size {
			return nil, errCorrupt
		}
		return out, nil
	}
	n, err := lz4Block(out, 0, src)
	if err != nil || n != size {
		return nil, errCorrupt
	}
	return out, nil
}

func lz4Hadoop(out, src []byte) bool {
	pos := 0
	for len(src) > 0 {
		if len(src) < 8 {
			return false
		}
		dsize, csize := binary.BigEndian.Uint32(src), binary.BigEndian.Uint32(src[4:])
		src = src[8:]
		if uint64(csize) > uint64(len(src)) || uint64(dsize) > uint64(len(out)-pos) {
			return false
		}
		n, err := lz4Block(out[:pos+int(dsize)], pos, src[:csize])
		if err != nil || n != int(dsize) {
			return false
		}
		pos += n
		src = src[csize:]
	}
	return pos == len(out)
}

const lz4FrameMagic = 0x184D2204

// lz4Frame decompresses an LZ4 frame into out and returns its size.
func lz4Frame(out, src []byte) (int, error) {
	if len(src) < 7 {
		return 0, errCorrupt
	}
	flg := src[4]
	if flg>>6 != 1 {
		return 0, errCorrupt // version
	}
	p := 7 // magic, FLG, BD, header checksum
	if flg&0x08 != 0 {
		p += 8 // content size
	}
	if flg&0x01 != 0 {
		p += 4 // dictionary id
	}
	pos := 0
	for {
		if p+4 > len(src) {
			return 0, errCorrupt
		}
		bsize := binary.LittleEndian.Uint32(src[p:])
		p += 4
		if bsize == 0 {
			return pos, nil // the end mark (a content checksum may follow)
		}
		raw := bsize&0x80000000 != 0
		bsize &= 0x7fffffff
		if uint64(bsize) > uint64(len(src)-p) {
			return 0, errCorrupt
		}
		block := src[p : p+int(bsize)]
		p += int(bsize)
		if flg&0x10 != 0 {
			p += 4 // block checksum
		}
		if raw {
			if len(block) > len(out)-pos {
				return 0, errCorrupt
			}
			pos += copy(out[pos:], block)
			continue
		}
		// blocks may refer to the data of the blocks before them
		n, err := lz4Block(out, pos, block)
		if err != nil {
			return 0, err
		}
		pos += n
	}
}

// lz4Block decompresses an LZ4 block into dst from dst[start:], where
// matches may reach back before start, and returns the size decompressed.
func lz4Block(dst []byte, start int, src []byte) (int, error) {
	d, s := start, 0
	for s < len(src) {
		token := src[s]
		s++
		lit := int(token >> 4)
		if lit == 15 {
			for {
				if s >= len(src) {
					return 0, errCorrupt
				}
				c := src[s]
				s++
				lit += int(c)
				if lit > len(src) {
					return 0, errCorrupt
				}
				if c != 255 {
					break
				}
			}
		}
		if lit > len(src)-s || lit > len(dst)-d {
			return 0, errCorrupt
		}
		d += copy(dst[d:], src[s:s+lit])
		s += lit
		if s == len(src) {
			break // the last sequence has only literals
		}
		if s+2 > len(src) {
			return 0, errCorrupt
		}
		off := int(src[s]) | int(src[s+1])<<8
		s += 2
		if off == 0 || off > d {
			return 0, errCorrupt
		}
		ml := int(token & 15)
		if ml == 15 {
			for {
				if s >= len(src) {
					return 0, errCorrupt
				}
				c := src[s]
				s++
				ml += int(c)
				if ml > len(dst) {
					return 0, errCorrupt
				}
				if c != 255 {
					break
				}
			}
		}
		ml += 4
		if ml > len(dst)-d {
			return 0, errCorrupt
		}
		if off >= ml {
			d += copy(dst[d:d+ml], dst[d-off:d-off+ml])
		} else {
			for i := range ml {
				dst[d+i] = dst[d-off+i]
			}
			d += ml
		}
	}
	return d - start, nil
}
