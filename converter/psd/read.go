package psd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"unicode/utf16"
)

// The structure of a Photoshop file (Adobe Photoshop File Formats
// Specification): a header, colour mode data, image resources, layer and
// mask information, and the image data, the composite of the document.
// PSB, the large document format (version 2), widens some lengths to 64
// bits. Pixel data is read only when it is decoded.

// Colour modes.
const (
	modeBitmap       = 0
	modeGrayscale    = 1
	modeIndexed      = 2
	modeRGB          = 3
	modeCMYK         = 4
	modeMultichannel = 7
	modeDuotone      = 8
	modeLab          = 9
)

var modeNames = map[int]string{modeBitmap: "bitmap", modeGrayscale: "grayscale", modeIndexed: "indexed", modeRGB: "RGB",
	modeCMYK: "CMYK", modeMultichannel: "multichannel", modeDuotone: "duotone", modeLab: "Lab"}

// Image resource IDs.
const (
	resResolution   = 1005
	resTransparency = 1047 // the transparent colour of an indexed image
	resVersion      = 1057
	resXMP          = 1060
)

// Channel IDs of layer records.
const (
	chAlpha    = -1
	chMask     = -2 // the layer mask (or the rendered vector mask, with a real mask)
	chRealMask = -3 // the layer mask of a layer with a vector mask too
)

var errFormat = errors.New("not a Photoshop file")

type header struct {
	psb      bool // version 2: the large document format
	channels int
	w, h     int
	depth    int
	mode     int
}

// baseChannels is the number of colour channels of the mode.
func (h *header) baseChannels() int {
	switch h.mode {
	case modeRGB, modeLab:
		return 3
	case modeCMYK:
		return 4
	case modeMultichannel:
		return h.channels
	}
	return 1
}

type file struct {
	r   io.ReaderAt
	hdr header
	// palette is the colour table of an indexed image (256 reds, greens, blues).
	palette   []byte
	resources map[int][]byte
	layers    []*layer // bottom to top, as stored
	// globalAlpha: the first channel after the colour channels of the
	// composite is its transparency.
	globalAlpha bool
	imageData   int64 // offset of the image data section
	size        int64
}

type channel struct {
	id     int
	off    int64 // compression code and data
	length int64
}

type mask struct {
	top, left, bottom, right int
	defaultColor             uint8
	flags                    uint8
}

func (m *mask) disabled() bool { return m.flags&2 != 0 }
func (m *mask) empty() bool    { return m.bottom <= m.top || m.right <= m.left }

type layer struct {
	top, left, bottom, right int
	channels                 []channel
	blend                    string // blend mode key
	opacity                  uint8
	fill                     uint8 // fill opacity
	clipping                 bool  // clipped to the layer below
	hidden                   bool
	name                     string
	mask, realMask           *mask
	// section is the section divider type: 0 a layer, 1 or 2 the top of a
	// group (open or closed), 3 the bottom of a group.
	section      int
	sectionBlend string
	artboard     *artboard
	fillColor    *[3]float64 // a solid colour fill layer (0-255 RGB)
	effects      bool        // has layer effects turned on
	keys         []string    // the additional information keys
}

func (l *layer) w() int { return l.right - l.left }
func (l *layer) h() int { return l.bottom - l.top }

type artboard struct {
	top, left, bottom, right int
	// background: 1 white, 2 black, 3 transparent, 4 the colour.
	background int
	color      [3]float64
}

// cursor reads big-endian values from a section of the file. Errors stick:
// after one, reads return zeros and err keeps the first error.
type cursor struct {
	r        io.ReaderAt
	off, end int64
	err      error
	buf      [8]byte
}

func (c *cursor) read(n int) []byte {
	var b []byte
	if n <= len(c.buf) {
		b = c.buf[:n]
	} else {
		b = make([]byte, n)
	}
	if c.err != nil {
		clear(b)
		return b
	}
	if int64(n) > c.end-c.off {
		c.err = io.ErrUnexpectedEOF
		clear(b)
		return b
	}
	if _, err := c.r.ReadAt(b, c.off); err != nil && !(err == io.EOF && c.off+int64(n) <= c.end) {
		c.err = err
		clear(b)
		return b
	}
	c.off += int64(n)
	return b
}

func (c *cursor) bytes(n int64) []byte {
	if n < 0 || n > c.end-c.off || n > 1<<30 {
		if c.err == nil {
			c.err = io.ErrUnexpectedEOF
		}
		return nil
	}
	return append([]byte(nil), c.read(int(n))...)
}

func (c *cursor) u8() uint8 { return c.read(1)[0] }

// peek returns the next byte without reading it.
func (c *cursor) peek() uint8 {
	v := c.u8()
	if c.err == nil {
		c.off--
	}
	return v
}
func (c *cursor) u16() uint16 { return binary.BigEndian.Uint16(c.read(2)) }
func (c *cursor) u32() uint32 { return binary.BigEndian.Uint32(c.read(4)) }
func (c *cursor) u64() uint64 { return binary.BigEndian.Uint64(c.read(8)) }
func (c *cursor) i32() int    { return int(int32(c.u32())) }
func (c *cursor) f64() float64 {
	return math.Float64frombits(c.u64())
}

// length reads a 32-bit length, or a 64-bit one when wide.
func (c *cursor) length(wide bool) int64 {
	if wide {
		v := c.u64()
		if v > math.MaxInt64/2 {
			c.err = errFormat
			return 0
		}
		return int64(v)
	}
	return int64(c.u32())
}

func (c *cursor) skip(n int64) {
	if c.err != nil {
		return
	}
	if n < 0 || n > c.end-c.off {
		c.err = io.ErrUnexpectedEOF
		return
	}
	c.off += n
}

// sub returns a cursor over the next n bytes and moves past them.
func (c *cursor) sub(n int64) *cursor {
	s := &cursor{r: c.r, off: c.off, end: c.off + n, err: c.err}
	if n < 0 || n > c.end-c.off {
		s.end = c.off
		if c.err == nil {
			c.err = io.ErrUnexpectedEOF
		}
		s.err = c.err
	}
	c.skip(n)
	return s
}

func (c *cursor) left() int64 { return c.end - c.off }

// pascal reads a Pascal string padded to a multiple of pad bytes.
func (c *cursor) pascal(pad int) string {
	n := int(c.u8())
	s := string(c.bytes(int64(n)))
	total := 1 + n
	if r := total % pad; r != 0 {
		c.skip(int64(pad - r))
	}
	return s
}

// unicode reads a length-prefixed UTF-16 string.
func (c *cursor) unicode() string {
	n := int64(c.u32())
	if n > c.left()/2 {
		c.err = errFormat
		return ""
	}
	b := c.bytes(2 * n)
	u := make([]uint16, n)
	for i := range u {
		u[i] = binary.BigEndian.Uint16(b[2*i:])
	}
	for len(u) > 0 && u[len(u)-1] == 0 {
		u = u[:len(u)-1]
	}
	return string(utf16.Decode(u))
}

// readFile reads the structure of a Photoshop file.
func readFile(r io.ReaderAt, size int64) (*file, error) {
	f := &file{r: r, size: size, resources: map[int][]byte{}}
	c := &cursor{r: r, end: size}
	if string(c.read(4)) != "8BPS" {
		return nil, errFormat
	}
	switch c.u16() {
	case 1:
	case 2:
		f.hdr.psb = true
	default:
		return nil, fmt.Errorf("unknown version")
	}
	c.skip(6)
	f.hdr.channels = int(c.u16())
	f.hdr.h = int(c.u32())
	f.hdr.w = int(c.u32())
	f.hdr.depth = int(c.u16())
	f.hdr.mode = int(c.u16())
	if c.err != nil {
		return nil, c.err
	}
	limit := 30000
	if f.hdr.psb {
		limit = 300000
	}
	if f.hdr.w < 1 || f.hdr.h < 1 || f.hdr.w > limit || f.hdr.h > limit || f.hdr.channels < 1 || f.hdr.channels > 56 {
		return nil, fmt.Errorf("bad image size or channel count")
	}
	switch f.hdr.depth {
	case 1, 8, 16, 32:
	default:
		return nil, fmt.Errorf("unsupported bit depth %d", f.hdr.depth)
	}
	if _, ok := modeNames[f.hdr.mode]; !ok {
		return nil, fmt.Errorf("unsupported colour mode %d", f.hdr.mode)
	}
	if f.hdr.channels < f.hdr.baseChannels() {
		return nil, fmt.Errorf("%d channel(s) for %s", f.hdr.channels, modeNames[f.hdr.mode])
	}

	// Colour mode data: the colour table of an indexed image.
	cm := c.sub(int64(c.u32()))
	if f.hdr.mode == modeIndexed {
		f.palette = cm.bytes(768)
		if cm.err != nil {
			return nil, fmt.Errorf("colour table: %w", cm.err)
		}
	}

	// Image resources.
	rs := c.sub(int64(c.u32()))
	for rs.left() >= 12 && rs.err == nil {
		sig := string(rs.read(4))
		if sig != "8BIM" && sig != "MeSa" && sig != "AgHg" && sig != "PHUT" && sig != "DCSR" {
			break
		}
		id := int(rs.u16())
		rs.pascal(2)
		n := int64(rs.u32())
		data := rs.bytes(n)
		if n%2 != 0 {
			rs.skip(1)
		}
		if _, ok := f.resources[id]; !ok && sig == "8BIM" && rs.err == nil {
			f.resources[id] = data
		}
	}

	// Layer and mask information. A broken section loses the layers, not
	// the composite after it.
	lm := c.sub(c.length(f.hdr.psb))
	if c.err != nil {
		return nil, c.err
	}
	f.imageData = c.off
	if lm.left() > 0 {
		if err := f.readLayerSection(lm); err != nil {
			f.layers = nil
			return f, fmt.Errorf("layers: %w", err)
		}
	}
	return f, nil
}

func (f *file) readLayerSection(lm *cursor) error {
	li := lm.sub(lm.length(f.hdr.psb))
	if li.left() > 0 {
		if err := f.readLayerInfo(li); err != nil {
			return err
		}
	}
	lm.sub(int64(lm.u32())) // global layer mask information
	// Additional information of the document: 16- and 32-bit documents keep
	// their layers in Lr16 or Lr32, and say that the composite has
	// transparency with Mt16, Mt32 or Mtrn.
	f.taggedBlocks(lm, func(key string, c *cursor) {
		switch key {
		case "Lr16", "Lr32", "Layr":
			if len(f.layers) == 0 {
				if err := f.readLayerInfo(c); err != nil && lm.err == nil {
					lm.err = err
				}
			}
		case "Mtrn", "Mt16", "Mt32":
			f.globalAlpha = true
		}
	})
	if lm.err != nil && lm.err != io.ErrUnexpectedEOF {
		return lm.err
	}
	return nil
}

// wideKeys are the additional information keys whose length is 64 bits in
// a PSB file.
var wideKeys = map[string]bool{"LMsk": true, "Lr16": true, "Lr32": true, "Layr": true, "Mt16": true, "Mt32": true,
	"Mtrn": true, "Alph": true, "FMsk": true, "lnk2": true, "FEid": true, "FXid": true, "PxSD": true}

// taggedBlocks walks additional information blocks ("8BIM" or "8B64", a
// key, a length and the data, padded to an even length). Writers pad some
// blocks further, to a multiple of four: zero bytes before a block are
// skipped.
func (f *file) taggedBlocks(c *cursor, fn func(key string, c *cursor)) {
	for c.err == nil {
		for c.left() > 0 && c.peek() == 0 {
			c.skip(1)
		}
		if c.left() < 12 {
			return
		}
		sig := string(c.read(4))
		if sig != "8BIM" && sig != "8B64" {
			return
		}
		key := string(c.read(4))
		n := c.length(f.hdr.psb && (sig == "8B64" || wideKeys[key]))
		if c.err != nil {
			return
		}
		fn(key, c.sub(n))
		if n%2 != 0 && c.left() > 0 {
			c.skip(1)
		}
	}
}

func (f *file) readLayerInfo(c *cursor) error {
	count := int(int16(c.u16()))
	if count < 0 {
		count = -count
		f.globalAlpha = true
	}
	layers := make([]*layer, 0, min(count, 4096))
	for i := 0; i < count && c.err == nil; i++ {
		l := &layer{opacity: 255, fill: 255}
		l.top, l.left, l.bottom, l.right = c.i32(), c.i32(), c.i32(), c.i32()
		n := int(c.u16())
		if n > 56 {
			return fmt.Errorf("layer %d has %d channels", i+1, n)
		}
		for j := 0; j < n; j++ {
			id := int(int16(c.u16()))
			l.channels = append(l.channels, channel{id: id, length: c.length(f.hdr.psb)})
		}
		if string(c.read(4)) != "8BIM" {
			return fmt.Errorf("layer %d: bad blend mode signature", i+1)
		}
		l.blend = string(c.read(4))
		l.opacity = c.u8()
		l.clipping = c.u8() != 0
		flags := c.u8()
		l.hidden = flags&2 != 0
		c.skip(1)
		extra := c.sub(int64(c.u32()))
		md := extra.sub(int64(extra.u32()))
		if md.left() >= 18 {
			m := &mask{}
			m.top, m.left, m.bottom, m.right = md.i32(), md.i32(), md.i32(), md.i32()
			m.defaultColor, m.flags = md.u8(), md.u8()
			if m.flags&16 != 0 {
				params := md.u8()
				for _, bit := range []uint8{1, 2, 4, 8} {
					if params&bit != 0 {
						if bit == 1 || bit == 4 {
							md.skip(1)
						} else {
							md.skip(8)
						}
					}
				}
			}
			l.mask = m
			if md.left() >= 18 {
				rm := &mask{}
				rm.flags, rm.defaultColor = md.u8(), md.u8()
				rm.top, rm.left, rm.bottom, rm.right = md.i32(), md.i32(), md.i32(), md.i32()
				l.realMask = rm
			}
		}
		extra.sub(int64(extra.u32())) // blending ranges
		l.name = extra.pascal(4)
		f.taggedBlocks(extra, func(key string, b *cursor) {
			l.keys = append(l.keys, key)
			switch key {
			case "luni":
				if s := b.unicode(); b.err == nil {
					l.name = s
				}
			case "lsct", "lsdk":
				l.section = int(b.u32())
				if b.left() >= 8 && string(b.read(4)) == "8BIM" {
					l.sectionBlend = string(b.read(4))
				}
			case "iOpa":
				l.fill = b.u8()
			case "artb", "artd", "abdd":
				if l.artboard == nil {
					l.artboard = readArtboard(b)
				}
			case "lfx2", "lmfx":
				b.skip(8) // versions
				d := readDescriptor(b)
				on, ok := true, false
				if d != nil {
					on, ok = d.items["masterFXSwitch"].(bool)
				}
				l.effects = l.effects || d == nil || (on || !ok) && enabled(d.items)
			case "SoCo":
				b.skip(4) // version
				if d := readDescriptor(b); d != nil {
					if clr, ok := d.items["Clr "].(*descriptor); ok {
						if rgb, ok := clr.rgb(); ok {
							l.fillColor = &rgb
						}
					}
				}
			}
		})
		if c.err != nil {
			break
		}
		layers = append(layers, l)
	}
	if c.err != nil {
		return c.err
	}
	// Channel image data follows the records, layer by layer.
	for _, l := range layers {
		for i := range l.channels {
			ch := &l.channels[i]
			ch.off = c.off
			c.skip(ch.length)
		}
	}
	if c.err != nil {
		return fmt.Errorf("channel data: %w", c.err)
	}
	f.layers = layers
	return nil
}

// enabled reports whether a value of an effects descriptor holds an effect
// that is turned on.
func enabled(v any) bool {
	switch v := v.(type) {
	case map[string]any:
		for _, e := range v {
			if enabled(e) {
				return true
			}
		}
	case *descriptor:
		if on, ok := v.items["enab"].(bool); ok {
			return on
		}
		return enabled(v.items)
	case []any:
		for _, e := range v {
			if enabled(e) {
				return true
			}
		}
	}
	return false
}

func readArtboard(b *cursor) *artboard {
	b.skip(4) // version
	d := readDescriptor(b)
	if d == nil {
		return nil
	}
	rect, ok := d.items["artboardRect"].(*descriptor)
	if !ok {
		return nil
	}
	num := func(key string) int {
		v, _ := rect.items[key].(float64)
		return int(math.Round(v))
	}
	a := &artboard{top: num("Top "), left: num("Left"), bottom: num("Btom"), right: num("Rght"), background: 1}
	if v, ok := d.items["artboardBackgroundType"].(int64); ok {
		a.background = int(v)
	}
	if clr, ok := d.items["Clr "].(*descriptor); ok {
		a.color, _ = clr.rgb()
	}
	return a
}

// resolution returns the pixels per inch of the image (72 when the file does
// not say).
func (f *file) resolution() (x, y float64) {
	x, y = 72, 72
	if b := f.resources[resResolution]; len(b) >= 16 {
		if v := float64(binary.BigEndian.Uint32(b[0:])) / 65536; v >= 1 && v < 1e6 {
			x = v
		}
		if v := float64(binary.BigEndian.Uint32(b[8:])) / 65536; v >= 1 && v < 1e6 {
			y = v
		}
	}
	return x, y
}

// hasComposite reports whether the image data holds the composite of the
// layers. Photoshop leaves it out when "Maximize Compatibility" is off and
// says so in its version information.
func (f *file) hasComposite() bool {
	b := f.resources[resVersion]
	return len(b) < 5 || b[4] != 0
}
