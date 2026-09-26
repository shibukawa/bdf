package metafile

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"slices"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"golang.org/x/text/encoding/japanese"
)

// output is what replaying a metafile produced.
type output struct {
	doc      *bdf.Document
	drawn    bool
	warnings []string
	colors   []bdf.Color
	texts    []string
	langs    []string
	fonts    []string // families of the font records
	ops      map[byte]int
}

func replay(t *testing.T, data []byte, rc Recolor) *output {
	t.Helper()
	d := &output{doc: bdf.NewDocument(), ops: map[byte]int{}}
	fonts := fontset.New(fontdb.New(nil, nil, false), nil)
	b := canvas.NewBuilder(d.doc, fonts)
	cv := b.New()
	Draw(cv, data, 0, 0, 100, 100, &Options{Doc: d.doc, Fonts: fonts, Recolor: rc,
		Warn: func(msg string) { d.warnings = append(d.warnings, msg) }})
	d.drawn = cv.Drawn
	b.Encode()
	o, err := bdf.DecodeObject(d.doc.Part(cv.Hash()).Data)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range o.Fonts {
		d.fonts = append(d.fonts, f.Family)
	}
	err = o.Walk(func(in bdf.Instr) {
		d.ops[in.Op]++
		switch in.Op {
		case bdf.OpFillColor:
			d.colors = append(d.colors, bdf.Color(in.Args[0].(uint64)))
		case bdf.OpFillText:
			d.texts = append(d.texts, in.Args[0].(string))
		case bdf.OpMark:
			if byte(in.Args[0].(uint64)) == bdf.MarkLang {
				d.langs = append(d.langs, in.Args[1].(string))
			}
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if d.drawn != (d.ops[bdf.OpFillPath]+d.ops[bdf.OpStrokePath]+d.ops[bdf.OpFillText]+d.ops[bdf.OpImageSub] > 0) {
		t.Errorf("Drawn = %v", d.drawn)
	}
	return d
}

// emfRecord appends one EMF record made of little-endian 32-bit fields
// followed by raw bytes.
func emfRecord(buf *bytes.Buffer, typ uint32, fields []int32, raw []byte) {
	size := 8 + 4*len(fields) + len(raw)
	pad := (4 - size%4) % 4
	binary.Write(buf, binary.LittleEndian, typ)
	binary.Write(buf, binary.LittleEndian, uint32(size+pad))
	binary.Write(buf, binary.LittleEndian, fields)
	buf.Write(raw)
	buf.Write(make([]byte, pad))
}

// emfHeader starts an EMF whose frame is 100 × 100 pixels.
func emfHeader() *bytes.Buffer {
	var b bytes.Buffer
	// bounds, frame (0.01 mm), " EMF", version, bytes, records, handles,
	// description, palette, device and millimetre sizes
	emfRecord(&b, 1, []int32{0, 0, 99, 99, 0, 0, 2646, 2646, 0x464D4520, 0x10000, 0, 0, 3, 0, 0, 0, 3780, 3780, 1000, 1000}, nil)
	return &b
}

func utf16LE(s string) []byte {
	u := utf16.Encode([]rune(s))
	b := make([]byte, (len(u)*2+3)/4*4)
	for i, c := range u {
		binary.LittleEndian.PutUint16(b[i*2:], c)
	}
	return b
}

// extTextOutW is an EXTTEXTOUTW record drawing s at x,y.
func extTextOutW(b *bytes.Buffer, x, y int32, opts int32, s string) {
	// bounds, graphics mode, scales, reference, count, string offset,
	// options, rectangle, dx offset
	emfRecord(b, 84, []int32{0, 0, 0, 0, 1, 0, 0, x, y, int32(len([]rune(s))), 76, opts, 0, 0, 0, 0, 0}, utf16LE(s))
}

// dib24 is a packed 2 × 2 24-bit bottom-up DIB: red, green / blue, white.
func dib24() []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, []int32{40, 2, 2})
	binary.Write(&b, binary.LittleEndian, []uint16{1, 24})
	binary.Write(&b, binary.LittleEndian, []int32{0, 0, 0, 0, 0, 0})
	// rows bottom-up, BGR, padded to 4 bytes
	b.Write([]byte{255, 0, 0, 255, 255, 255, 0, 0})
	b.Write([]byte{0, 0, 255, 0, 255, 0, 0, 0})
	return b.Bytes()
}

func testEMF() []byte {
	b := emfHeader()
	emfRecord(b, 39, []int32{1, 0, 0xFFFFFF, 0}, nil) // CREATEBRUSHINDIRECT 1: white
	emfRecord(b, 37, []int32{1}, nil)                 // SELECTOBJECT
	emfRecord(b, 37, []int32{-0x7FFFFFF8}, nil)       // SELECTOBJECT NULL_PEN
	emfRecord(b, 43, []int32{0, 0, 100, 100}, nil)    // RECTANGLE
	emfRecord(b, 18, []int32{1}, nil)                 // SETBKMODE: transparent
	emfRecord(b, 39, []int32{2, 2, 0x0000FF, 4}, nil) // CREATEBRUSHINDIRECT 2: red cross hatch
	emfRecord(b, 37, []int32{2}, nil)
	emfRecord(b, 43, []int32{10, 10, 50, 50}, nil)
	face := make([]byte, 8+64)
	face[1] = 1 // underline
	copy(face[8:], utf16LE("Arial"))
	emfRecord(b, 82, []int32{3, -12, 0, 0, 0, 400}, face) // EXTCREATEFONTINDIRECTW
	emfRecord(b, 37, []int32{3}, nil)
	emfRecord(b, 24, []int32{0x800000}, nil) // SETTEXTCOLOR: navy
	emfRecord(b, 22, []int32{24}, nil)       // SETTEXTALIGN: baseline
	extTextOutW(b, 10, 80, 0, "EMF テキスト")
	// STRETCHDIBITS: bounds, destination, source, header and bits, usage,
	// raster operation, destination size
	emfRecord(b, 81, []int32{0, 0, 0, 0, 60, 60, 0, 0, 2, 2, 80, 40, 120, 16, 0, 0xCC0020, 20, 20}, dib24())
	emfRecord(b, 14, []int32{0, 16, 20}, nil) // EOF
	return b.Bytes()
}

func TestKind(t *testing.T) {
	placeable := make([]byte, 40)
	binary.LittleEndian.PutUint32(placeable, 0x9AC6CDD7)
	for _, c := range []struct {
		data []byte
		want string
	}{
		{testEMF(), "emf"},
		{testWMF(t), "wmf"},
		{placeable, "wmf"},
		{[]byte("\x89PNG\r\n\x1a\n and more bytes than a header"), ""},
	} {
		if got := Kind(c.data); got != c.want {
			t.Errorf("Kind = %q, want %q", got, c.want)
		}
	}
}

func TestDrawEMF(t *testing.T) {
	d := replay(t, testEMF(), nil)
	if len(d.warnings) != 0 {
		t.Errorf("warnings: %v", d.warnings)
	}
	for _, c := range []bdf.Color{bdf.RGBA(255, 255, 255, 255), bdf.RGBA(0, 0, 128, 255)} {
		if !slices.Contains(d.colors, c) {
			t.Errorf("no fill color %08X in %v", c, d.colors)
		}
	}
	if got := strings.Join(d.texts, ""); got != "EMF テキスト" {
		t.Errorf("text = %q (%q)", got, d.texts)
	}
	if !slices.Equal(d.langs, []string{"ja"}) {
		t.Errorf("languages = %q", d.langs)
	}
	if len(d.fonts) == 0 || d.fonts[0] != `"Arial", sans-serif` {
		t.Errorf("fonts = %q", d.fonts)
	}
	if d.ops[bdf.OpFillPaint] != 1 || d.ops[bdf.OpImageSub] != 1 || d.ops[bdf.OpFillRect] != 2 {
		t.Errorf("hatch fills %d, images %d, underlines %d", d.ops[bdf.OpFillPaint], d.ops[bdf.OpImageSub], d.ops[bdf.OpFillRect])
	}
	// the hatch tile: a red cross on a transparent background
	var imgs []image.Image
	for _, p := range d.doc.Parts() {
		img, err := png.Decode(bytes.NewReader(p.Data))
		if err == nil {
			imgs = append(imgs, img)
		}
	}
	if len(imgs) != 2 || imgs[0].Bounds().Dx() != 8 {
		t.Fatalf("images: %d", len(imgs))
	}
	if c := color.NRGBAModel.Convert(imgs[0].At(3, 0)).(color.NRGBA); c != (color.NRGBA{255, 0, 0, 255}) {
		t.Errorf("hatch line = %v", c)
	}
	if _, _, _, a := imgs[0].At(3, 3).RGBA(); a != 0 {
		t.Error("hatch background is not transparent")
	}
	// the bitmap, top row first
	if c := color.NRGBAModel.Convert(imgs[1].At(0, 0)).(color.NRGBA); c != (color.NRGBA{255, 0, 0, 255}) {
		t.Errorf("bitmap pixel 0,0 = %v", c)
	}
}

// swapRB is a Recolor that swaps red and blue.
type swapRB struct{ images int }

func (s *swapRB) Color(c color.NRGBA) color.NRGBA { return color.NRGBA{c.B, c.G, c.R, c.A} }
func (s *swapRB) Image(img *image.NRGBA)          { s.images++ }

func TestRecolor(t *testing.T) {
	rc := &swapRB{}
	d := replay(t, testEMF(), rc)
	if !slices.Contains(d.colors, bdf.RGBA(128, 0, 0, 255)) || slices.Contains(d.colors, bdf.RGBA(0, 0, 128, 255)) {
		t.Errorf("text color not recolored: %v", d.colors)
	}
	if rc.images != 1 {
		t.Errorf("bitmaps recolored: %d", rc.images)
	}
}

// testWMF writes "日本語" in a Shift_JIS font into a 1000 × 1000 window.
func testWMF(t *testing.T) []byte {
	var recs bytes.Buffer
	rec := func(fn uint16, args []byte) {
		if len(args)%2 == 1 {
			args = append(args, 0)
		}
		binary.Write(&recs, binary.LittleEndian, uint32(3+len(args)/2))
		binary.Write(&recs, binary.LittleEndian, fn)
		recs.Write(args)
	}
	words := func(v ...int16) []byte {
		var b bytes.Buffer
		binary.Write(&b, binary.LittleEndian, v)
		return b.Bytes()
	}
	sjis := func(s string) []byte {
		b, err := japanese.ShiftJIS.NewEncoder().Bytes([]byte(s))
		if err != nil {
			t.Fatal(err)
		}
		return b
	}
	rec(0x020B, words(0, 0))       // SETWINDOWORG y x
	rec(0x020C, words(1000, 1000)) // SETWINDOWEXT y x
	// CREATEFONTINDIRECT: height, width, escapement, orientation, weight,
	// italic, underline, strike out, charset (SHIFTJIS), …, face name
	lf := append(words(-100, 0, 0, 0, 400), 0, 0, 0, 128, 0, 0, 0, 0)
	rec(0x02FB, append(append(lf, sjis("ＭＳ ゴシック")...), 0))
	rec(0x012D, words(0)) // SELECTOBJECT
	text := sjis("日本語")
	rec(0x0521, append(append(words(int16(len(text))), text...), words(500, 100)...)) // TEXTOUT: count string y x
	rec(0x0000, nil)
	var b bytes.Buffer
	// standard header: type, header size, version, size, objects, max record, members
	binary.Write(&b, binary.LittleEndian, []uint16{1, 9, 0x300})
	binary.Write(&b, binary.LittleEndian, uint32(9+recs.Len()/2))
	binary.Write(&b, binary.LittleEndian, []uint16{1})
	binary.Write(&b, binary.LittleEndian, uint32(32))
	binary.Write(&b, binary.LittleEndian, []uint16{0})
	b.Write(recs.Bytes())
	return b.Bytes()
}

func TestDrawWMF(t *testing.T) {
	d := replay(t, testWMF(t), nil)
	if !slices.Equal(d.texts, []string{"日本語"}) {
		t.Errorf("text = %q", d.texts)
	}
	if len(d.fonts) != 1 || !strings.HasPrefix(d.fonts[0], `"ＭＳ ゴシック"`) {
		t.Errorf("fonts = %q", d.fonts)
	}
	// Han characters alone do not say which language they are
	if len(d.langs) != 0 {
		t.Errorf("languages = %q", d.langs)
	}
}

func TestWarnings(t *testing.T) {
	b := emfHeader()
	for range 2 {
		// GDICOMMENT: an EMF+ header record (type 0x4001) without the dual flag
		emfRecord(b, 70, []int32{12, 0x2B464D45, 0x00004001, 0}, nil)
		extTextOutW(b, 10, 10, 0x10, "\x01\x02") // ETO_GLYPH_INDEX
		// STRETCHDIBITS with a truncated bitmap
		emfRecord(b, 81, []int32{0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 80, 40, 120, 4, 0, 0xCC0020, 20, 20}, append(dib24()[:40], 1, 2, 3, 4))
	}
	emfRecord(b, 14, []int32{0, 16, 20}, nil)
	d := replay(t, b.Bytes(), nil)
	want := []string{
		"EMF+ only metafiles are not supported",
		"metafile text drawn by glyph index is not supported",
		"a metafile bitmap could not be decoded",
	}
	if !slices.Equal(d.warnings, want) {
		t.Errorf("warnings = %q, want %q", d.warnings, want)
	}
}

func TestDecodeDIB(t *testing.T) {
	// 1 bit per pixel, 2 × 2, top-down, palette black / yellow
	var bmi bytes.Buffer
	binary.Write(&bmi, binary.LittleEndian, []int32{40, 2, -2})
	binary.Write(&bmi, binary.LittleEndian, []uint16{1, 1})
	binary.Write(&bmi, binary.LittleEndian, []int32{0, 0, 0, 0, 2, 0})
	bmi.Write([]byte{0, 0, 0, 0, 0, 255, 255, 0})
	bits := []byte{0x40, 0, 0, 0, 0x80, 0, 0, 0}
	img, err := decodeDIB(bmi.Bytes(), bits)
	if err != nil {
		t.Fatal(err)
	}
	yellow, black := color.NRGBA{255, 255, 0, 255}, color.NRGBA{0, 0, 0, 255}
	for _, c := range []struct {
		x, y int
		want color.NRGBA
	}{{0, 0, black}, {1, 0, yellow}, {0, 1, yellow}, {1, 1, black}} {
		if got := img.NRGBAAt(c.x, c.y); got != c.want {
			t.Errorf("pixel %d,%d = %v, want %v", c.x, c.y, got, c.want)
		}
	}
	// a packed DIB: the bits follow the header and the palette
	packed := append(bmi.Bytes(), bits...)
	if got := dibBits(packed); !bytes.Equal(got, bits) {
		t.Errorf("dibBits = %v", got)
	}
	if _, err := decodeDIB(bmi.Bytes(), bits[:4]); err == nil {
		t.Error("truncated bits decoded")
	}
}
