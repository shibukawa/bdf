// Package barcode makes QR codes and barcodes as drawings: the symbols that
// forms carry for the machines that read them (slips, labels, invoices).
//
// Encode makes a symbol of a text: a QR code (model 2, versions 1 to 40,
// the four levels of error correction, in the numeric, alphanumeric or
// byte mode that the text needs), or a linear barcode: JAN/EAN-13, EAN-8,
// UPC-A, Code 39, Code 128, GS1-128, ITF (interleaved 2 of 5, ITF-14) and
// NW-7 (Codabar). A symbol is a grid of modules with its quiet zone, and
// Path draws its dark modules in a box, with human-readable digits for
// JAN/EAN and UPC-A: as vectors, which stay sharp at any size and in print.
//
// The package writes symbols; it does not read them. It does not check
// what a text means (that the digits of a JAN code are those of a product),
// only that the symbology can hold it, and the check digits that a
// symbology has.
package barcode

import (
	"fmt"
	"strings"

	"github.com/shibukawa/bdf"
)

// The symbologies.
const (
	QR      = "qr"
	EAN13   = "ean13" // JAN (13 digits)
	EAN8    = "ean8"  // JAN (8 digits)
	UPCA    = "upca"
	Code39  = "code39"
	Code128 = "code128"
	GS1128  = "gs1-128"
	ITF     = "itf"
	ITF14   = "itf14"
	NW7     = "nw7" // Codabar
)

// names are the names that symbologies go by.
var names = map[string]string{
	"qr": QR, "qrcode": QR,
	"ean13": EAN13, "ean": EAN13, "jan": EAN13, "jan13": EAN13,
	"ean8": EAN8, "jan8": EAN8,
	"upca": UPCA, "upc": UPCA,
	"code39":  Code39,
	"code128": Code128,
	"gs1128":  GS1128, "ean128": GS1128,
	"itf": ITF, "itf14": ITF14,
	"nw7": NW7, "codabar": NW7,
}

// Kinds lists the symbologies by the names Kind takes for them.
func Kinds() []string {
	return []string{QR, EAN13, EAN8, UPCA, Code39, Code128, GS1128, ITF, ITF14, NW7}
}

// Kind returns the symbology of a name ("jan", "JAN-13", "Codabar" …), or
// "" for a name that is none.
func Kind(name string) string {
	return names[strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToLower(name))]
}

// Options are the options of a symbol.
type Options struct {
	// Level is the level of error correction of a QR code: "L" (7% of
	// the symbol may be lost), "M" (15%, the default), "Q" (25%) or "H"
	// (30%).
	Level string
	// Quiet is the width of the quiet zone around the symbol in modules:
	// the white margin readers need. 0 is that of the symbology (4
	// modules around a QR code, 10 at the sides of most barcodes), and
	// NoQuiet none: the margin is then for what is around the symbol to
	// leave.
	Quiet int
	// NoText leaves out the human-readable digits of JAN/EAN and UPC-A.
	// The digits are drawn by Path and Object by default, including a
	// check digit added by Encode. The modules are unchanged either way.
	NoText bool
}

// NoQuiet is the Options.Quiet of a symbol without a quiet zone.
const NoQuiet = -1

// Symbol is a symbol: a grid of modules, dark or light.
type Symbol struct {
	// Kind is the symbology.
	Kind string
	// Text is what the symbol says: the text, with the check digit of
	// the symbology when the text came without it.
	Text string
	// W and H are the size of the grid with its quiet zone, in modules.
	// A linear barcode is a module high; Path leaves room for the digits
	// of JAN/EAN and UPC-A inside the box where it draws the bars.
	W, H int
	// Linear is set for a barcode of bars.
	Linear bool
	mods   []bool
	noText bool
	opts   Options // what the symbol was encoded with (see Mark)
}

// At reports whether the module at column x and row y is dark.
func (s *Symbol) At(x, y int) bool {
	if x < 0 || y < 0 || x >= s.W || y >= s.H {
		return false
	}
	return s.mods[y*s.W+x]
}

// Rects returns the dark modules as rectangles of modules (x, y, width,
// height): the runs of dark modules of the rows, with those that are the
// same in rows that follow each other as one.
func (s *Symbol) Rects() [][4]int {
	var out [][4]int
	open := map[[2]int]int{} // the runs of the row before: their rectangles
	for y := 0; y < s.H; y++ {
		next := map[[2]int]int{}
		for x := 0; x < s.W; {
			if !s.mods[y*s.W+x] {
				x++
				continue
			}
			x0 := x
			for x < s.W && s.mods[y*s.W+x] {
				x++
			}
			key := [2]int{x0, x}
			if i, ok := open[key]; ok {
				out[i][3]++
				next[key] = i
				continue
			}
			next[key] = len(out)
			out = append(out, [4]int{x0, y, x - x0, 1})
		}
		open = next
	}
	return out
}

// Path returns the symbol (and its quiet zone) in the box at (x, y) of
// the size w × h. Its dark modules are rectangles. JAN/EAN and UPC-A
// also have their digits below the bars unless Options.NoText was set;
// the bars are shortened to leave room inside the box. The digits are
// outlines, so no font is needed. Fill the path with the nonzero rule.
func (s *Symbol) Path(x, y, w, h float32) *bdf.Path {
	p := &bdf.Path{}
	if s.W == 0 || s.H == 0 {
		return p
	}
	barHeight, digitScale := h, float32(0)
	if !s.noText && (s.Kind == EAN13 || s.Kind == EAN8 || s.Kind == UPCA) && w > 0 && h > 0 && len(s.Text) > 0 {
		// The ink of the digits is 750 units high, their advance 620.
		// Keep their aspect ratio even in a tall, narrow box. At most
		// 22% of its height is ink, with a quarter of that as a gap.
		digitScale = min(h*0.22/750, w/float32(len(s.Text)*620))
		barHeight -= float32(float32(digitScale*750) * 1.25)
	}
	mx, my := w/float32(s.W), barHeight/float32(s.H)
	for _, r := range s.Rects() {
		// (the products are rounded before they are added: the same on every
		// machine, where one that multiplies and adds at once would differ)
		p.Rect(x+float32(float32(r[0])*mx), y+float32(float32(r[1])*my), float32(r[2])*mx, float32(r[3])*my)
	}
	if digitScale > 0 {
		x0 := x + (w-float32(float32(len(s.Text)*620)*digitScale))/2
		// M PLUS's lowest digit outline is 10 units below its baseline.
		baseline := y + h - float32(10*digitScale)
		for _, c := range []byte(s.Text) {
			appendDigit(p, c-'0', x0, baseline, digitScale)
			x0 += float32(620 * digitScale)
		}
	}
	return p
}

// Object returns the symbol as an object of the size w × h: its dark
// modules and, for JAN/EAN and UPC-A, its digits, in black on nothing.
// Options.NoText leaves out the digits. The object needs no font.
func (s *Symbol) Object(w, h float32) *bdf.Object {
	o := bdf.NewObject()
	o.FillColor(bdf.RGB(0, 0, 0)).FillPath(o.AddPath(s.Path(0, 0, w, h)), bdf.NonZero)
	o.SetBBox(0, 0, w, h)
	return o
}

// MaxText bounds the text of a symbol, in bytes: more than a QR code
// holds.
const MaxText = 8192

// Encode makes the symbol of a text in a symbology (a name Kind knows).
func Encode(kind, text string, opts *Options) (*Symbol, error) {
	k := Kind(kind)
	if k == "" {
		return nil, fmt.Errorf("barcode: %q is not a symbology (%s)", kind, strings.Join(Kinds(), ", "))
	}
	if text == "" {
		return nil, fmt.Errorf("barcode: nothing to make a %s symbol of", k)
	}
	if len(text) > MaxText {
		return nil, fmt.Errorf("barcode: a text of %d bytes is more than a symbol holds", len(text))
	}
	o := Options{}
	if opts != nil {
		o = *opts
	}
	quiet := -1 // that of the symbology
	switch {
	case o.Quiet < 0:
		quiet = 0
	case o.Quiet > 0:
		quiet = min(o.Quiet, 64)
	}
	var s *Symbol
	var err error
	switch k {
	case QR:
		s, err = encodeQR(text, o.Level, quiet, -1)
	case EAN13, EAN8, UPCA:
		s, err = encodeEAN(k, text, quiet)
	case Code39:
		s, err = encodeCode39(text, quiet)
	case Code128, GS1128:
		s, err = encodeCode128(text, k == GS1128, quiet)
	case ITF, ITF14:
		s, err = encodeITF(text, k == ITF14, quiet)
	case NW7:
		s, err = encodeNW7(text, quiet)
	}
	if err != nil {
		return nil, err
	}
	s.Kind = k
	s.noText = o.NoText
	s.opts = o
	return s, nil
}

// linear makes the symbol of bars: the widths of the bars and the spaces
// between them in modules, a bar first, with a quiet zone at the sides.
func linear(widths []int, left, right int) *Symbol {
	n := left + right
	for _, w := range widths {
		n += w
	}
	s := &Symbol{W: n, H: 1, Linear: true, mods: make([]bool, n)}
	x := left
	for i, w := range widths {
		if i%2 == 0 {
			for k := 0; k < w; k++ {
				s.mods[x+k] = true
			}
		}
		x += w
	}
	return s
}

// digits reports whether a text is digits alone.
func digits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return s != ""
}

// mod10 returns the check digit of the digits of a GTIN (JAN, EAN, UPC,
// ITF-14): with the weights 3 and 1 from the right.
func mod10(s string) byte {
	sum := 0
	for i := 0; i < len(s); i++ {
		d := int(s[len(s)-1-i] - '0')
		if i%2 == 0 {
			d *= 3
		}
		sum += d
	}
	return byte('0' + (10-sum%10)%10)
}

// gtin returns the digits of a code of n digits with its check digit: the
// text may come without it (n-1 digits), and is refused with another.
func gtin(kind, text string, n int) (string, error) {
	t := strings.TrimSpace(text)
	if !digits(t) || len(t) != n && len(t) != n-1 {
		return "", fmt.Errorf("barcode: %s takes %d digits (or %d, without the check digit), not %q", kind, n, n-1, text)
	}
	if len(t) == n-1 {
		return t + string(mod10(t)), nil
	}
	if c := mod10(t[:n-1]); c != t[n-1] {
		return "", fmt.Errorf("barcode: the check digit of the %s code %s is %c, not %c", kind, t, c, t[n-1])
	}
	return t, nil
}
