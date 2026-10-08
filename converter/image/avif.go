package image

import (
	"bytes"
	"errors"

	"github.com/shibukawa/bdf/converter/internal/isobmff"
)

// reader reads big-endian fields, remembering that it ran out of data.
type reader struct {
	b   []byte
	bad bool
}

func (r *reader) take(n int) []byte {
	if r.bad || n > len(r.b) {
		r.bad = true
		return make([]byte, n)
	}
	v := r.b[:n]
	r.b = r.b[n:]
	return v
}

func (r *reader) u8() uint64  { return uint64(r.take(1)[0]) }
func (r *reader) u16() uint64 { return uint64(be.Uint16(r.take(2))) }
func (r *reader) u32() uint64 { return uint64(be.Uint32(r.take(4))) }

// uint reads an unsigned integer of 0, 4 or 8 bytes (the sizes of iloc).
func (r *reader) uint(size uint64) uint64 {
	switch size {
	case 0:
		return 0
	case 4:
		return r.u32()
	case 8:
		return be.Uint64(r.take(8))
	}
	r.bad = true
	return 0
}

// cstring reads a NUL-terminated string.
func (r *reader) cstring() string {
	i := bytes.IndexByte(r.b, 0)
	if r.bad || i < 0 {
		r.bad = true
		return ""
	}
	s := string(r.b[:i])
	r.b = r.b[i+1:]
	return s
}

// readAVIF reads an AVIF image (HEIF, ISO/IEC 23008-12): the size of the
// primary item as its rotation (irot) leaves it, and its Exif and XMP
// items. An image sequence without a primary item takes the size of its
// track.
func readAVIF(b []byte) (*picture, error) {
	p := &picture{format: "avif"}
	top := isobmff.Boxes(b)
	// an image sequence: "avis" as the major or a compatible brand, and a movie
	if ftyp := isobmff.Find(top, "ftyp"); isobmff.Find(top, "moov") != nil {
		for i := 0; i+4 <= len(ftyp); i += 4 {
			p.animated = p.animated || i != 4 && string(ftyp[i:i+4]) == "avis"
		}
	}
	if meta := isobmff.Find(top, "meta"); len(meta) >= 4 {
		readHEIFMeta(p, b, isobmff.Boxes(meta[4:]))
	}
	if p.w == 0 {
		// a sequence: the track header's width and height (16.16)
		if trak := isobmff.Find(isobmff.Boxes(isobmff.Find(top, "moov")), "trak"); trak != nil {
			if tkhd := isobmff.Find(isobmff.Boxes(trak), "tkhd"); len(tkhd) >= 84 {
				at := 76 // version 0
				if tkhd[0] == 1 {
					at = 88
				}
				if len(tkhd) >= at+8 {
					p.w, p.h = float64(be.Uint32(tkhd[at:])>>16), float64(be.Uint32(tkhd[at+4:])>>16)
				}
			}
		}
	}
	if p.w == 0 || p.h == 0 {
		return nil, errors.New("the AVIF image has no size")
	}
	return p, nil
}

// readHEIFMeta reads the items of a meta box: the primary item's size and
// rotation, and the Exif and XMP items.
func readHEIFMeta(p *picture, file []byte, meta []isobmff.Box) {
	primary := uint64(0)
	if pitm := isobmff.Find(meta, "pitm"); len(pitm) >= 4 {
		r := &reader{b: pitm[4:]}
		if pitm[0] == 0 {
			primary = r.u16()
		} else {
			primary = r.u32()
		}
	}
	// item types
	type item struct{ typ, contentType string }
	items := map[uint64]item{}
	if iinf := isobmff.Find(meta, "iinf"); len(iinf) >= 4 {
		r := &reader{b: iinf[4:]}
		if iinf[0] == 0 {
			r.u16()
		} else {
			r.u32()
		}
		for _, infe := range isobmff.Boxes(r.b) {
			if infe.Type != "infe" || len(infe.Data) < 4 || infe.Data[0] < 2 {
				continue
			}
			e := &reader{b: infe.Data[4:]}
			var id uint64
			if infe.Data[0] == 2 {
				id = e.u16()
			} else {
				id = e.u32()
			}
			e.u16() // protection index
			it := item{typ: string(e.take(4))}
			e.cstring() // name
			if it.typ == "mime" {
				it.contentType = e.cstring()
			}
			if !e.bad || it.typ != "mime" {
				items[id] = it
			}
		}
	}
	// the primary item's properties
	rotated := false
	if iprp := isobmff.Boxes(isobmff.Find(meta, "iprp")); iprp != nil {
		props := isobmff.Boxes(isobmff.Find(iprp, "ipco"))
		for _, ipma := range iprp {
			if ipma.Type != "ipma" || len(ipma.Data) < 4 {
				continue
			}
			version, flags := ipma.Data[0], ipma.Data[3]
			r := &reader{b: ipma.Data[4:]}
			for n := r.u32(); n > 0 && !r.bad; n-- {
				var id uint64
				if version < 1 {
					id = r.u16()
				} else {
					id = r.u32()
				}
				for k := r.u8(); k > 0 && !r.bad; k-- {
					var index uint64
					if flags&1 != 0 {
						index = r.u16() & 0x7fff
					} else {
						index = r.u8() & 0x7f
					}
					if id != primary || index == 0 || index > uint64(len(props)) {
						continue
					}
					switch prop := props[index-1]; prop.Type {
					case "ispe":
						if len(prop.Data) >= 12 {
							p.w, p.h = float64(be.Uint32(prop.Data[4:])), float64(be.Uint32(prop.Data[8:]))
						}
					case "irot":
						rotated = len(prop.Data) >= 1 && prop.Data[0]&1 != 0 // 90° or 270°
					}
				}
			}
		}
	}
	if rotated {
		p.w, p.h = p.h, p.w
	}
	// Exif and XMP items
	iloc := isobmff.Find(meta, "iloc")
	if len(iloc) < 6 {
		return
	}
	idat := isobmff.Find(meta, "idat")
	version := iloc[0]
	r := &reader{b: iloc[4:]}
	sizes := r.u8()
	offsetSize, lengthSize := sizes>>4, sizes&15
	sizes = r.u8()
	baseSize, indexSize := sizes>>4, uint64(0)
	if version == 1 || version == 2 {
		indexSize = sizes & 15
	}
	var count uint64
	if version < 2 {
		count = r.u16()
	} else {
		count = r.u32()
	}
	// The items read are parts of the file: together they are no larger
	// than it is, however often the box lists them.
	room := uint64(len(file))
	for ; count > 0 && !r.bad; count-- {
		var id, method uint64
		if version < 2 {
			id = r.u16()
		} else {
			id = r.u32()
		}
		if version == 1 || version == 2 {
			method = r.u16() & 15
		}
		r.u16() // data reference index
		base := r.uint(baseSize)
		it := items[id]
		exif, xmp := it.typ == "Exif", it.typ == "mime" && it.contentType == "application/rdf+xml"
		var data []byte
		src := file
		if method == 1 {
			src = idat
		}
		for n := r.u16(); n > 0 && !r.bad; n-- {
			if indexSize > 0 {
				r.uint(indexSize)
			}
			off, length := base+r.uint(offsetSize), r.uint(lengthSize)
			if !exif && !xmp {
				continue // the extents of other items are only read past
			}
			if length == 0 && off <= uint64(len(src)) {
				length = uint64(len(src)) - off
			}
			if method > 1 || off > uint64(len(src)) || length > uint64(len(src))-off || length > room {
				exif, xmp = false, false
				continue
			}
			room -= length
			data = append(data, src[off:off+length]...)
		}
		switch {
		case exif && len(data) >= 4:
			// the offset of the TIFF header, then the Exif data
			if at := uint64(be.Uint32(data)) + 4; at < uint64(len(data)) {
				p.exif = append(p.exif, data[at:])
			}
		case xmp && data != nil:
			p.xmp = append(p.xmp, data)
		}
	}
}
