package raster

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"io"
	"strings"
)

// ReadPage is a page read back from PWG Raster: its header and its
// pixels.
type ReadPage struct {
	Page
	// Image is an *image.Gray for sgray_8, an *image.RGBA for srgb_8,
	// and an *image.Gray of 0 and 255 for the types of a bit a pixel
	// (black where the bit is ink).
	Image image.Image
	// Lines are the lines written, each with its repeat count: how
	// much the compression of lines gained.
	Lines int
}

// Read reads a PWG Raster stream: for tests and tools that look at what
// was sent. It refuses pages of more than maxPixels pixels (64 Mi when
// 0).
func Read(r io.Reader, maxPixels int) ([]*ReadPage, error) {
	if maxPixels <= 0 {
		maxPixels = 64 << 20
	}
	br := bufio.NewReader(r)
	sync := make([]byte, 4)
	if _, err := io.ReadFull(br, sync); err != nil || string(sync) != Sync {
		return nil, errors.New("raster: not a PWG Raster stream")
	}
	var pages []*ReadPage
	var h [HeaderSize]byte
	for {
		if _, err := io.ReadFull(br, h[:]); err != nil {
			if err == io.EOF {
				return pages, nil
			}
			return nil, fmt.Errorf("raster: page %d: %w", len(pages)+1, err)
		}
		p, err := parseHeader(h)
		if err != nil {
			return nil, fmt.Errorf("raster: page %d: %w", len(pages)+1, err)
		}
		if p.Width*p.Height > maxPixels {
			return nil, fmt.Errorf("raster: page %d is %d × %d pixels", len(pages)+1, p.Width, p.Height)
		}
		rp := &ReadPage{Page: *p}
		if err := rp.decode(br); err != nil {
			return nil, fmt.Errorf("raster: page %d: %w", len(pages)+1, err)
		}
		pages = append(pages, rp)
	}
}

func cstr(b []byte) string {
	if i := strings.IndexByte(string(b), 0); i >= 0 {
		return string(b[:i])
	}
	return string(b)
}

func parseHeader(h [HeaderSize]byte) (*Page, error) {
	u32 := func(off int) int { return int(int32(binary.BigEndian.Uint32(h[off:]))) }
	if cstr(h[offMediaClass:offMediaClass+64]) != "PwgRaster" {
		return nil, errors.New("not a PWG Raster header")
	}
	p := &Page{Width: u32(offWidth), Height: u32(offHeight), XDPI: u32(offHWResolution), YDPI: u32(offHWResolution + 4),
		PageWidth: u32(offPageSize), PageHeight: u32(offPageSize + 4), PageSizeName: cstr(h[offPageSizeName : offPageSizeName+64]),
		MediaType: cstr(h[offMediaType : offMediaType+64]), MediaColor: cstr(h[offMediaColor : offMediaColor+64]),
		Duplex: u32(offDuplex) != 0, Tumble: u32(offTumble) != 0, CrossFeed: u32(offCrossFeed), Feed: u32(offFeed),
		Copies: u32(offNumCopies), TotalPages: u32(offTotalPages), Quality: u32(offQuality), MediaPosition: u32(offMediaPos)}
	cs, bpc, bpp := u32(offColorSpace), u32(offBitsPerColor), u32(offBitsPerPixel)
	switch {
	case cs == csBlack && bpc == 1 && bpp == 1:
		p.Type = Black1
	case cs == csSGray && bpc == 1 && bpp == 1:
		p.Type = SGray1
	case cs == csSGray && bpc == 8 && bpp == 8:
		p.Type = SGray8
	case cs == csSRGB && bpc == 8 && bpp == 24:
		p.Type = SRGB8
	default:
		return nil, fmt.Errorf("colour space %d of %d bits a colour, %d a pixel", cs, bpc, bpp)
	}
	if p.Width <= 0 || p.Height <= 0 {
		return nil, fmt.Errorf("a page of %d × %d pixels", p.Width, p.Height)
	}
	if bpl := u32(offBytesPerLine); bpl != (p.Width*bpp+7)/8 {
		return nil, fmt.Errorf("%d bytes a line for %d pixels of %d bits", bpl, p.Width, bpp)
	}
	return p, nil
}

func (rp *ReadPage) decode(br *bufio.Reader) error {
	_, _, bits, _, _ := rp.Type.layout()
	unit := max(bits/8, 1)
	n := (rp.Width*bits + 7) / 8
	line := make([]byte, n)
	var gray *image.Gray
	var rgba *image.RGBA
	if rp.Type == SRGB8 {
		rgba = image.NewRGBA(image.Rect(0, 0, rp.Width, rp.Height))
		rp.Image = rgba
	} else {
		gray = image.NewGray(image.Rect(0, 0, rp.Width, rp.Height))
		rp.Image = gray
	}
	put := func(y int) {
		switch rp.Type {
		case SGray8:
			copy(gray.Pix[y*gray.Stride:], line)
		case SRGB8:
			row := rgba.Pix[y*rgba.Stride:]
			for x := 0; x < rp.Width; x++ {
				row[4*x], row[4*x+1], row[4*x+2], row[4*x+3] = line[3*x], line[3*x+1], line[3*x+2], 255
			}
		default:
			row := gray.Pix[y*gray.Stride:]
			for x := 0; x < rp.Width; x++ {
				bit := line[x>>3]&(0x80>>(x&7)) != 0
				ink := bit == (rp.Type == Black1)
				row[x] = 255
				if ink {
					row[x] = 0
				}
			}
		}
	}
	for y := 0; y < rp.Height; {
		rep, err := br.ReadByte()
		if err != nil {
			return err
		}
		rp.Lines++
		for i := 0; i < n; {
			c, err := br.ReadByte()
			if err != nil {
				return err
			}
			if c < 128 {
				count := (int(c) + 1) * unit
				if i+count > n {
					return errors.New("a run past the end of the line")
				}
				if _, err := io.ReadFull(br, line[i:i+unit]); err != nil {
					return err
				}
				for k := unit; k < count; k++ {
					line[i+k] = line[i+k-unit]
				}
				i += count
			} else {
				count := (257 - int(c)) * unit
				if i+count > n {
					return errors.New("pixels past the end of the line")
				}
				if _, err := io.ReadFull(br, line[i:i+count]); err != nil {
					return err
				}
				i += count
			}
		}
		for k := 0; k <= int(rep) && y < rp.Height; k++ {
			put(y)
			y++
		}
	}
	return nil
}
