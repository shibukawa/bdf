package jbig2

// Combination operators (7.4.1.5, 7.4.3.1.1).
const (
	opOr = iota
	opAnd
	opXor
	opXnor
	opReplace
)

// newBitmap allocates a white bitmap; callers check the size first.
func newBitmap(w, h int) *Bitmap {
	stride := (w + 7) >> 3
	return &Bitmap{Width: w, Height: h, Stride: stride, Data: make([]byte, stride*h)}
}

// pixel returns the pixel at (x, y), or 0 outside the bitmap.
func (b *Bitmap) pixel(x, y int) int {
	if x < 0 || y < 0 || x >= b.Width || y >= b.Height {
		return 0
	}
	return int(b.Data[y*b.Stride+x>>3]>>(7-uint(x&7))) & 1
}

// set sets the pixel at (x, y), which must lie in the bitmap, to 1.
func (b *Bitmap) set(x, y int) {
	b.Data[y*b.Stride+x>>3] |= 0x80 >> uint(x&7)
}

// fill sets every pixel to v.
func (b *Bitmap) fill(v int) {
	c := byte(0)
	if v != 0 {
		c = 0xff
	}
	for i := range b.Data {
		b.Data[i] = c
	}
	b.clearPadding()
}

// clearPadding clears the bits past the width at the end of each row.
func (b *Bitmap) clearPadding() {
	if b.Width&7 == 0 {
		return
	}
	mask := byte(0xff) << uint(8-b.Width&7)
	for i := b.Stride - 1; i < len(b.Data); i += b.Stride {
		b.Data[i] &= mask
	}
}

// row returns row y, or nil outside the bitmap.
func (b *Bitmap) row(y int) []byte {
	if y < 0 || y >= b.Height {
		return nil
	}
	return b.Data[y*b.Stride : (y+1)*b.Stride]
}

// bit returns pixel x of a row of width w, or 0 outside it or for a nil row.
func bit(row []byte, w, x int) int {
	if row == nil || x < 0 || x >= w {
		return 0
	}
	return int(row[x>>3]>>(7-uint(x&7))) & 1
}

// sub copies the w×h rectangle at (x, y), which must lie in the bitmap.
func (b *Bitmap) sub(x, y, w, h int) *Bitmap {
	s := newBitmap(w, h)
	for r := 0; r < h; r++ {
		composeRow(s.row(r), 0, b.row(y+r), x, w, opReplace)
	}
	return s
}

// compose combines src into b with its top left corner at (x, y), clipped to b.
func (b *Bitmap) compose(src *Bitmap, x, y, op int) {
	sx, sy := 0, 0
	if x < 0 {
		sx, x = -x, 0
	}
	if y < 0 {
		sy, y = -y, 0
	}
	w := min(src.Width-sx, b.Width-x)
	h := min(src.Height-sy, b.Height-y)
	if w <= 0 || h <= 0 {
		return
	}
	for r := 0; r < h; r++ {
		composeRow(b.row(y+r), x, src.row(sy+r), sx, w, op)
	}
}

// composeRow combines the w bits of s starting at bit sx into d starting at bit dx.
func composeRow(d []byte, dx int, s []byte, sx, w, op int) {
	for w > 0 {
		off := uint(dx & 7)
		n := 8 - int(off)
		if n > w {
			n = w
		}
		// The next 8 source bits, aligned to the destination bit position.
		i, sh := sx>>3, uint(sx&7)
		v := s[i] << sh
		if sh != 0 && i+1 < len(s) {
			v |= s[i+1] >> (8 - sh)
		}
		v >>= off
		mask := byte(0xff) >> off &^ (byte(0xff) >> (off + uint(n)))
		p := &d[dx>>3]
		var r byte
		switch op {
		case opOr:
			r = *p | v
		case opAnd:
			r = *p & v
		case opXor:
			r = *p ^ v
		case opXnor:
			r = ^(*p ^ v)
		default:
			r = v
		}
		*p = *p&^mask | r&mask
		dx += n
		sx += n
		w -= n
	}
}
