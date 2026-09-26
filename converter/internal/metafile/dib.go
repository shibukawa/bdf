package metafile

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg" // BI_JPEG bitmaps
	_ "image/png"  // BI_PNG bitmaps
)

var errBadDIB = errors.New("bad DIB")

// decodeDIB decodes a device-independent bitmap (BITMAPINFO + bits).
func decodeDIB(bmi, bits []byte) (*image.NRGBA, error) {
	if len(bmi) < 12 {
		return nil, errBadDIB
	}
	hsize := int(le32(bmi, 0))
	var w, h, bpp, compression int
	if hsize == 12 { // BITMAPCOREHEADER
		w, h, bpp = int(le16(bmi, 4)), int(int16(le16(bmi, 6))), int(le16(bmi, 10))
	} else {
		if len(bmi) < 40 {
			return nil, errBadDIB
		}
		w, h = int(int32(le32(bmi, 4))), int(int32(le32(bmi, 8)))
		bpp, compression = int(le16(bmi, 14)), int(le32(bmi, 16))
	}
	if compression == 4 || compression == 5 { // BI_JPEG, BI_PNG
		img, _, err := image.Decode(bytes.NewReader(bits))
		if err != nil {
			return nil, err
		}
		return toNRGBA(img), nil
	}
	topDown := h < 0
	if h < 0 {
		h = -h
	}
	if w <= 0 || h <= 0 || w*h > 1<<26 {
		return nil, errBadDIB
	}
	// color table
	var pal []color.NRGBA
	if bpp <= 8 {
		n := 1 << bpp
		if hsize >= 40 {
			if used := int(le32(bmi, 32)); used > 0 && used < n {
				n = used
			}
		}
		entry := 4
		if hsize == 12 {
			entry = 3
		}
		off := hsize
		for i := 0; i < n && off+entry <= len(bmi); i++ {
			pal = append(pal, color.NRGBA{bmi[off+2], bmi[off+1], bmi[off], 255})
			off += entry
		}
	}
	masks := [3]uint32{0x7c00, 0x03e0, 0x001f}
	if bpp == 32 {
		masks = [3]uint32{0xff0000, 0xff00, 0xff}
	}
	if compression == 3 && len(bmi) >= hsize+12 {
		masks = [3]uint32{le32(bmi, hsize), le32(bmi, hsize+4), le32(bmi, hsize+8)}
		if hsize >= 52 {
			masks = [3]uint32{le32(bmi, 40), le32(bmi, 44), le32(bmi, 48)}
		}
	} else if compression != 0 {
		return nil, errBadDIB // RLE
	}
	stride := ((w*bpp + 31) / 32) * 4
	if len(bits) < stride*h {
		return nil, errBadDIB
	}
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	field := func(v, mask uint32) uint8 {
		if mask == 0 {
			return 0
		}
		shift := 0
		for mask&1 == 0 {
			mask >>= 1
			shift++
		}
		return uint8((v >> shift & mask) * 255 / mask)
	}
	for row := 0; row < h; row++ {
		src := bits[row*stride:]
		y := h - 1 - row
		if topDown {
			y = row
		}
		for x := 0; x < w; x++ {
			var c color.NRGBA
			switch bpp {
			case 1, 2, 4, 8:
				idx := int(src[x*bpp/8]>>(8-bpp-(x*bpp%8))) & (1<<bpp - 1)
				if idx < len(pal) {
					c = pal[idx]
				} else {
					c = color.NRGBA{0, 0, 0, 255}
				}
			case 16:
				v := uint32(le16(src, x*2))
				c = color.NRGBA{field(v, masks[0]), field(v, masks[1]), field(v, masks[2]), 255}
			case 24:
				c = color.NRGBA{src[x*3+2], src[x*3+1], src[x*3], 255}
			case 32:
				v := le32(src, x*4)
				c = color.NRGBA{field(v, masks[0]), field(v, masks[1]), field(v, masks[2]), 255}
			default:
				return nil, errBadDIB
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img, nil
}

// dibBits returns the pixel data following a packed DIB header.
func dibBits(bmi []byte) []byte {
	if len(bmi) < 12 {
		return nil
	}
	hsize := int(le32(bmi, 0))
	bpp, compression, colors := 0, uint32(0), 0
	if hsize == 12 {
		bpp = int(le16(bmi, 10))
		if bpp <= 8 {
			colors = (1 << bpp) * 3
		}
	} else {
		if len(bmi) < 40 {
			return nil
		}
		bpp, compression = int(le16(bmi, 14)), le32(bmi, 16)
		n := int(le32(bmi, 32))
		if n == 0 && bpp <= 8 {
			n = 1 << bpp
		}
		colors = n * 4
		if compression == 3 && hsize == 40 {
			colors += 12
		}
	}
	off := hsize + colors
	if off > len(bmi) {
		return nil
	}
	return bmi[off:]
}

func toNRGBA(img image.Image) *image.NRGBA {
	if n, ok := img.(*image.NRGBA); ok {
		return n
	}
	b := img.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), img, b.Min, draw.Src)
	return out
}
