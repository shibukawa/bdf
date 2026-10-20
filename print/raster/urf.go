package raster

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// URFSync is the first bytes of an Apple Raster stream; the number of its
// pages follows, in four bytes.
const URFSync = "UNIRAST\x00"

// URFHeaderSize is the bytes of the header of a page of Apple Raster.
const URFHeaderSize = 32

// NewURFWriter returns a writer of an Apple Raster stream (image/urf) of
// pages pages to w. Its pages are of SGray8 (W8) or SRGB8 (SRGB24); their
// rows are compressed as those of PWG Raster. The header of a page holds
// what CUPS writes in one (raster-stream.c): the bits of a pixel, the
// colour space, the sides, the quality, the input tray, the size in pixels
// and the resolution.
func NewURFWriter(w io.Writer, pages int) *Writer {
	return &Writer{w: bufio.NewWriterSize(w, 64<<10), urf: true, urfPages: pages}
}

// colour spaces of the header of a page of Apple Raster
const (
	urfSGray = 0
	urfSRGB  = 1
)

// urf sides
const (
	urfOneSided = 1
	urfShort    = 2
	urfLong     = 3
)

func (p *Page) urfHeader() ([URFHeaderSize]byte, error) {
	var h [URFHeaderSize]byte
	if p.Width <= 0 || p.Height <= 0 || p.XDPI <= 0 {
		return h, errors.New("raster: an empty page")
	}
	switch p.Type {
	case SGray8:
		h[0], h[1] = 8, urfSGray
	case SRGB8:
		h[0], h[1] = 24, urfSRGB
	default:
		return h, fmt.Errorf("raster: Apple Raster has no type %q (sgray_8 or srgb_8)", p.Type)
	}
	switch {
	case !p.Duplex:
		h[2] = urfOneSided
	case p.Tumble:
		h[2] = urfShort
	default:
		h[2] = urfLong
	}
	h[3] = byte(min(max(p.Quality, 0), 255))
	h[5] = byte(min(max(p.MediaPosition, 0), 255))
	binary.BigEndian.PutUint32(h[12:], uint32(p.Width))
	binary.BigEndian.PutUint32(h[16:], uint32(p.Height))
	binary.BigEndian.PutUint32(h[20:], uint32(p.XDPI))
	return h, nil
}

// ReadURF reads an Apple Raster stream: for tests and tools that look at
// what was sent. It refuses pages of more than maxPixels pixels (64 Mi
// when 0), and more pages than the stream says it has.
func ReadURF(r io.Reader, maxPixels int) ([]*ReadPage, error) {
	if maxPixels <= 0 {
		maxPixels = 64 << 20
	}
	br := bufio.NewReader(r)
	var head [12]byte
	if _, err := io.ReadFull(br, head[:]); err != nil || string(head[:8]) != URFSync {
		return nil, errors.New("raster: not an Apple Raster stream")
	}
	count := binary.BigEndian.Uint32(head[8:])
	var pages []*ReadPage
	var h [URFHeaderSize]byte
	for {
		if _, err := io.ReadFull(br, h[:]); err != nil {
			if err == io.EOF {
				if uint32(len(pages)) != count {
					return nil, fmt.Errorf("raster: %d pages of %d", len(pages), count)
				}
				return pages, nil
			}
			return nil, fmt.Errorf("raster: page %d: %w", len(pages)+1, err)
		}
		if uint32(len(pages)) >= count {
			return nil, fmt.Errorf("raster: more than the %d pages the stream has", count)
		}
		p := Page{Width: int(binary.BigEndian.Uint32(h[12:])), Height: int(binary.BigEndian.Uint32(h[16:])),
			Quality: int(h[3]), MediaPosition: int(h[5]), CrossFeed: 1, Feed: 1, Copies: 1, TotalPages: int(count)}
		p.XDPI = int(binary.BigEndian.Uint32(h[20:]))
		p.YDPI = p.XDPI
		switch {
		case h[0] == 8 && h[1] == urfSGray:
			p.Type = SGray8
		case h[0] == 24 && h[1] == urfSRGB:
			p.Type = SRGB8
		default:
			return nil, fmt.Errorf("raster: page %d: %d bits a pixel in colour space %d", len(pages)+1, h[0], h[1])
		}
		p.Duplex, p.Tumble = h[2] == urfShort || h[2] == urfLong, h[2] == urfShort
		if p.Width <= 0 || p.Height <= 0 || p.Width > 1<<20 || p.Height > 1<<20 || p.Width*p.Height > maxPixels {
			return nil, fmt.Errorf("raster: page %d is %d × %d pixels", len(pages)+1, p.Width, p.Height)
		}
		rp := &ReadPage{Page: p}
		if err := rp.decode(br); err != nil {
			return nil, fmt.Errorf("raster: page %d: %w", len(pages)+1, err)
		}
		pages = append(pages, rp)
	}
}
