package imgconv

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
)

// Orientation returns the EXIF orientation (1–8) that browsers apply when
// they draw a JPEG or PNG image (from a JPEG APP1 segment or a PNG eXIf
// chunk), or 1 when there is none. Browsers draw WebP images as stored.
func Orientation(data []byte) int {
	var exif []byte
	switch Sniff(data) {
	case "jpeg":
		for pos := 2; pos+4 <= len(data) && data[pos] == 0xff; {
			marker := data[pos+1]
			if marker == 0xda || marker == 0xd9 { // start of scan, end of image
				break
			}
			n := int(binary.BigEndian.Uint16(data[pos+2:]))
			if n < 2 || pos+2+n > len(data) {
				break
			}
			seg := data[pos+4 : pos+2+n]
			if marker == 0xe1 && bytes.HasPrefix(seg, []byte("Exif\x00\x00")) {
				exif = seg[6:]
				break
			}
			pos += 2 + n
		}
	case "png":
		for pos := 8; pos+12 <= len(data); {
			n := int(binary.BigEndian.Uint32(data[pos:]))
			if n < 0 || pos+12+n > len(data) {
				break
			}
			if string(data[pos+4:pos+8]) == "eXIf" {
				exif = bytes.TrimPrefix(data[pos+8:pos+8+n], []byte("Exif\x00\x00"))
				break
			}
			pos += 12 + n
		}
	}
	if len(exif) < 8 {
		return 1
	}
	var order binary.ByteOrder
	switch string(exif[:4]) {
	case "II*\x00":
		order = binary.LittleEndian
	case "MM\x00*":
		order = binary.BigEndian
	default:
		return 1
	}
	ifd := int(order.Uint32(exif[4:]))
	if ifd < 0 || ifd+2 > len(exif) {
		return 1
	}
	count := int(order.Uint16(exif[ifd:]))
	for i := 0; i < count; i++ {
		e := ifd + 2 + 12*i
		if e+12 > len(exif) {
			break
		}
		if order.Uint16(exif[e:]) == 0x0112 && order.Uint16(exif[e+2:]) == 3 { // Orientation, SHORT
			if o := int(order.Uint16(exif[e+8:])); o >= 1 && o <= 8 {
				return o
			}
			break
		}
	}
	return 1
}

// Orient turns an image as its EXIF orientation says it is displayed.
func Orient(img image.Image, o int) image.Image {
	if o <= 1 || o > 8 {
		return img
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	src := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(src, src.Bounds(), img, b.Min, draw.Src)
	ow, oh := w, h
	if o >= 5 {
		ow, oh = h, w
	}
	dst := image.NewRGBA(image.Rect(0, 0, ow, oh))
	for y := 0; y < oh; y++ {
		for x := 0; x < ow; x++ {
			var sx, sy int
			switch o {
			case 2:
				sx, sy = w-1-x, y
			case 3:
				sx, sy = w-1-x, h-1-y
			case 4:
				sx, sy = x, h-1-y
			case 5:
				sx, sy = y, x
			case 6:
				sx, sy = y, h-1-x
			case 7:
				sx, sy = w-1-y, h-1-x
			case 8:
				sx, sy = w-1-y, x
			}
			copy(dst.Pix[y*dst.Stride+4*x:y*dst.Stride+4*x+4], src.Pix[sy*src.Stride+4*sx:])
		}
	}
	return dst
}
