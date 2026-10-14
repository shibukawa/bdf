package barcode

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
)

// row writes a row of modules as 0 and 1.
func row(s *Symbol, y int) string {
	var b strings.Builder
	for x := 0; x < s.W; x++ {
		if s.At(x, y) {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
	return b.String()
}

func mustEncode(t testing.TB, kind, text string, opts *Options) *Symbol {
	t.Helper()
	s, err := Encode(kind, text, opts)
	if err != nil {
		t.Fatalf("%s %q: %v", kind, text, err)
	}
	return s
}

// readQR reads the text of a QR code back from its modules (it corrects no
// errors: the error correction codewords are those of the data, or it
// fails). It returns the version, the level and the mask with it.
func readQR(t testing.TB, s *Symbol, quiet int) (text string, ver, level, mask int) {
	t.Helper()
	n := s.W - 2*quiet
	ver = (n - 17) / 4
	if s.W != s.H || ver < 1 || ver > 40 || n != ver*4+17 {
		t.Fatalf("a QR code of %d × %d modules", s.W, s.H)
	}
	at := func(x, y int) bool { return s.At(x+quiet, y+quiet) }
	// the format information, at its two places
	read := func(places [][2]int) int {
		v := 0
		for i, p := range places {
			if at(p[0], p[1]) {
				v |= 1 << i
			}
		}
		return v ^ 0x5412
	}
	var first, second [][2]int
	for i := 0; i <= 5; i++ {
		first = append(first, [2]int{8, i})
	}
	first = append(first, [2]int{8, 7}, [2]int{8, 8}, [2]int{7, 8})
	for i := 9; i < 15; i++ {
		first = append(first, [2]int{14 - i, 8})
	}
	for i := 0; i < 8; i++ {
		second = append(second, [2]int{n - 1 - i, 8})
	}
	for i := 8; i < 15; i++ {
		second = append(second, [2]int{8, n - 15 + i})
	}
	f := read(first)
	if g := read(second); g != f {
		t.Fatalf("the format information is %015b and %015b", f, g)
	}
	rem := f >> 10
	for i := 0; i < 10; i++ {
		rem = rem<<1 ^ rem>>9*0x537
	}
	if rem&0x3ff != f&0x3ff {
		t.Fatalf("the format information %015b has another error correction", f)
	}
	mask = f >> 10 & 7
	for i, b := range qrLevelBits {
		if b == f>>13 {
			level = i
		}
	}
	if !at(8, n-8) {
		t.Fatal("the dark module is light")
	}
	// the codewords, from the modules that are no function patterns
	q := &qrMatrix{ver: ver, size: n, mods: make([]bool, n*n), fn: make([]bool, n*n)}
	q.function()
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			if q.fn[y*n+x] {
				if x != 8 && y != 8 && q.mods[y*n+x] != at(x, y) {
					t.Fatalf("the function pattern at %d, %d is wrong", x, y)
				}
				continue
			}
			q.mods[y*n+x] = at(x, y)
		}
	}
	q.mask(mask)
	raw := make([]byte, qrRawModules(ver)/8)
	i := 0
	for right := n - 1; right >= 1; right -= 2 {
		if right == 6 {
			right = 5
		}
		for v := 0; v < n; v++ {
			for j := 0; j < 2; j++ {
				x, y := right-j, v
				if (right+1)&2 == 0 {
					y = n - 1 - v
				}
				if !q.fn[y*n+x] && i < len(raw)*8 {
					if q.mods[y*n+x] {
						raw[i>>3] |= 0x80 >> (i & 7)
					}
					i++
				}
			}
		}
	}
	// the blocks, and their error correction
	blocks, ecc := int(qrBlocks[level][ver]), int(qrEccLen[level][ver])
	short, shortLen := blocks-len(raw)%blocks, len(raw)/blocks-ecc
	parts := make([][]byte, blocks)
	k := 0
	for c := 0; c <= shortLen; c++ {
		for b := 0; b < blocks; b++ {
			if c < shortLen || b >= short {
				parts[b] = append(parts[b], raw[k])
				k++
			}
		}
	}
	div := rsDivisor(ecc)
	var data []byte
	eccs := raw[k:]
	for b := 0; b < blocks; b++ {
		for c, e := range rsRemainder(parts[b], div) {
			if eccs[c*blocks+b] != e {
				t.Fatalf("version %d level %d: the error correction of block %d is not that of its data", ver, level, b)
			}
		}
		data = append(data, parts[b]...)
	}
	// the segment
	pos := 0
	get := func(w int) int {
		v := 0
		for ; w > 0; w, pos = w-1, pos+1 {
			v = v<<1 | int(data[pos>>3]>>(7-pos&7)&1)
		}
		return v
	}
	var out strings.Builder
	switch m := get(4); m {
	case 1:
		for c := get(qrCountBits(qrNumeric, ver)); c > 0; c -= 3 {
			w := min(c, 3)
			fmt.Fprintf(&out, "%0*d", w, get(3*w+1))
		}
	case 2:
		c := get(qrCountBits(qrAlnum, ver))
		for ; c >= 2; c -= 2 {
			v := get(11)
			out.WriteByte(qrAlnumChars[v/45])
			out.WriteByte(qrAlnumChars[v%45])
		}
		if c == 1 {
			out.WriteByte(qrAlnumChars[get(6)])
		}
	case 4:
		for c := get(qrCountBits(qrByte, ver)); c > 0; c-- {
			out.WriteByte(byte(get(8)))
		}
	default:
		t.Fatalf("mode %04b", m)
	}
	return out.String(), ver, level, mask
}

func TestQR(t *testing.T) {
	long := strings.Repeat("The quick brown fox jumps over the lazy dog. 0123456789 ", 60)
	for _, c := range []struct {
		text, level string
		ver         int
	}{
		{"0", "", 1}, {"01234567", "M", 1}, {"HELLO WORLD", "Q", 1}, {"HELLO WORLD", "H", 2},
		{"https://example.com/orders/1024", "M", 3}, {"日本語のテキスト、帳票", "L", 3},
		{strings.Repeat("7", 300), "H", 11}, {strings.Repeat("AB12 ", 120), "Q", 19}, {long[:1200], "M", 29}, {long[:2900], "L", 40},
		{strings.Repeat("9", 7089), "L", 40}, {strings.Repeat("A", 4296), "l", 40},
	} {
		s := mustEncode(t, "QR", c.text, &Options{Level: c.level})
		text, ver, level, _ := readQR(t, s, 4)
		want := strings.Index("LMQH", strings.ToUpper(c.level))
		if c.level == "" {
			want = 1
		}
		if text != c.text || ver != c.ver || level != want || s.Linear || s.Text != c.text || s.Kind != QR {
			t.Errorf("%.20q at %q: read %.20q, version %d (want %d), level %d (want %d)", c.text, c.level, text, ver, c.ver, level, want)
		}
	}
	// every version at every level, filled, and every mask
	for lv := 0; lv < 4; lv++ {
		for ver := 1; ver <= 40; ver++ {
			n := (qrDataCodewords(ver, lv)*8 - 4 - qrCountBits(qrByte, ver)) / 8
			text := strings.Repeat("abcdefghij klmnopqrstuvwxyz,", n/28+1)[:n]
			s, err := encodeQR(text, "LMQH"[lv:lv+1], 0, (ver+lv)%8)
			if err != nil {
				t.Fatal(err)
			}
			got, v, l, m := readQR(t, s, 0)
			if got != text || v != ver || l != lv || m != (ver+lv)%8 {
				t.Errorf("version %d level %d: read version %d, level %d, mask %d and %d of %d bytes", ver, lv, v, l, m, len(got), n)
			}
			// one byte more does not fit
			if s, _ := encodeQR(text+"x", "LMQH"[lv:lv+1], 0, 0); s != nil && (s.W-17)/4 <= ver {
				t.Errorf("version %d level %d holds more than %d bytes", ver, lv, n)
			}
		}
	}
	// the symbols are those of other writers of QR codes (test/barcode/read.py
	// compares them with python-qrcode module by module): what they are now
	for text, want := range map[string]string{
		"HELLO WORLD":                     "7076f3ab",
		"https://example.com/orders/1024": "9bba8b5a",
		"日本語のテキスト、帳票":                     "492ad646",
	} {
		s := mustEncode(t, QR, text, &Options{Level: "Q"})
		var rows strings.Builder
		for y := 0; y < s.H; y++ {
			rows.WriteString(row(s, y) + "\n")
		}
		if got := fmt.Sprintf("%x", sha256.Sum256([]byte(rows.String())))[:8]; got != want {
			t.Errorf("the QR code of %q is another than it was: %s, was %s", text, got, want)
		}
	}
	// the quiet zone
	a, b, c := mustEncode(t, QR, "1", nil), mustEncode(t, QR, "1", &Options{Quiet: NoQuiet}), mustEncode(t, QR, "1", &Options{Quiet: 2})
	if a.W != 29 || b.W != 21 || c.W != 25 || !b.At(0, 0) || a.At(3, 3) || !a.At(4, 4) {
		t.Errorf("sizes %d, %d and %d", a.W, b.W, c.W)
	}
	for _, bad := range []struct{ text, level, want string }{
		{strings.Repeat("9", 7090), "L", "more than a QR code of level L holds"}, {"x", "X", "not a level of error correction"}, {"", "", "nothing to make"},
		{strings.Repeat("x", MaxText+1), "", "more than a symbol holds"},
	} {
		if _, err := Encode(QR, bad.text, &Options{Level: bad.level}); err == nil || !strings.Contains(err.Error(), bad.want) {
			t.Errorf("%.10q at %q: %v, want %q", bad.text, bad.level, err, bad.want)
		}
	}
}

func TestLinear(t *testing.T) {
	for _, c := range []struct{ kind, text, said, want string }{
		// JAN: the check digit is added, and the bars are those of the standard
		{"jan", "490123456789", "4901234567894", "00000000000" + "101" + "0001011" + "0100111" + "0011001" + "0010011" + "0100001" + "0011101" + "01010" +
			"1001110" + "1010000" + "1000100" + "1001000" + "1110100" + "1011100" + "101" + "0000000"},
		{"EAN-8", "9638507", "96385074", "0000000" + "101" + "0001011" + "0101111" + "0111101" + "0110111" + "01010" + "1001110" + "1110010" + "1000100" + "1011100" + "101" + "0000000"},
		{"upc", "03600029145", "036000291452", ""},
		// (the start character, A and the stop character: wide elements of three modules)
		{"code39", "A", "A", "0000000000" + "1000101110111010" + "1110101000101110" + "100010111011101" + "0000000000"},
		{"itf", "123", "0123", ""},
		{"ITF-14", "1490123456789", "14901234567891", ""},
		{"nw7", "123", "A123A", ""},
		{"Codabar", "b12d", "B12D", ""},
		{"code128", "Code 128", "Code 128", ""},
	} {
		s := mustEncode(t, c.kind, c.text, nil)
		if s.Text != c.said || !s.Linear || s.H != 1 || s.Kind != Kind(c.kind) {
			t.Errorf("%s %q says %q (want %q), linear %v, %d rows", c.kind, c.text, s.Text, c.said, s.Linear, s.H)
		}
		if got := row(s, 0); c.want != "" && got != c.want {
			t.Errorf("%s %q:\n got %s\nwant %s", c.kind, c.text, got, c.want)
		}
	}
	// the characters of Code 128: of 11 modules (the stop character of 13), each another
	seen := map[string]bool{}
	for v, p := range code128Chars {
		sum := 0
		for _, e := range []byte(p) {
			sum += int(e - '0')
		}
		if want := map[bool]int{false: 11, true: 13}[v == c128Stop]; sum != want || seen[p] || len(p) != map[bool]int{false: 6, true: 7}[v == c128Stop] {
			t.Errorf("the character %d of Code 128, %s, has %d modules", v, p, sum)
		}
		seen[p] = true
	}
	// the code sets of Code 128: pairs of digits where there are four at an end or six inside
	vals := func(kind, text string) string {
		s := mustEncode(t, kind, text, &Options{Quiet: NoQuiet})
		var out []string
		r := row(s, 0)
		for x := 0; x+11 <= len(r); x += 11 {
			var w strings.Builder
			for i, run := x, 1; i < x+11; i, run = i+1, run+1 {
				if i+1 == x+11 || r[i+1] != r[i] {
					fmt.Fprint(&w, run)
					run = 0
				}
			}
			for v, p := range code128Chars {
				if p == w.String() {
					out = append(out, fmt.Sprint(v))
				}
			}
		}
		return strings.Join(out, " ")
	}
	for _, c := range []struct{ kind, text, want string }{
		{"code128", "1234567890", "105 12 34 56 78 90"},
		{"code128", "AB", "104 33 34"},
		{"code128", "AB12345", "104 33 34 17 99 23 45"},
		{"code128", "1234AB", "105 12 34 100 33 34"},
		{"code128", "A123456B", "104 33 99 12 34 56 100 34"},
		{"code128", "A1234B", "104 33 17 18 19 20 34"},
		{"code128", "a\tb", "104 65 101 73 100 66"},
		{"code128", "\tA", "103 73 33"},
		{"gs1-128", "(01)04912345678904(10)AB", "105 102 1 4 91 23 45 67 89 4 10 100 33 34"},
		{"gs1-128", "(10)AB(21)1234", "104 102 17 16 33 34 102 99 21 12 34"},
	} {
		got := vals(c.kind, c.text)
		// (without the check character; the stop character, of 13 modules, is not read)
		if f := strings.Fields(got); len(f) < 2 || strings.Join(f[:len(f)-1], " ") != c.want {
			t.Errorf("%s %q: %s, want %s and the check", c.kind, c.text, got, c.want)
		}
	}
	// the check character of Code 128: the sum of the values times their
	// places, with the start character, modulo 103
	if got := vals("code128", "PJJ123C"); got != "104 48 42 42 17 18 19 35 55" {
		t.Errorf("PJJ123C: %s", got)
	}
	for _, c := range []struct{ kind, text, want string }{
		{"jan", "4901234567890", "the check digit of the ean13 code 4901234567890 is 4, not 0"},
		{"jan", "49012345", "takes 13 digits"}, {"ean8", "abc", "takes 8 digits"}, {"upca", "1", "takes 12 digits"},
		{"code39", "abc", `code39 has no "a"`}, {"code39", "A*B", `code39 has no "*"`}, {"code39", "**", "nothing"},
		{"code128", "日本", "the characters of ASCII alone"}, {"gs1-128", "(01", "application identifiers in parentheses"}, {"gs1-128", "(01)(10)A", "has no data"},
		{"itf", "12a", "takes digits"}, {"itf14", "123", "takes 14 digits"}, {"nw7", "12A3", "between the start and stop characters"},
		{"nw7", "1x", "takes digits"}, {"datamatrix", "1", "is not a symbology"},
	} {
		if _, err := Encode(c.kind, c.text, nil); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s %q: %v, want %q", c.kind, c.text, err, c.want)
		}
	}
	if Kind("JAN-13") != EAN13 || Kind("GS1-128") != GS1128 || Kind("NW-7") != NW7 || Kind("QR Code") != QR || Kind("pdf417") != "" || len(Kinds()) != 10 {
		t.Error("the names of the symbologies")
	}
}

func TestPath(t *testing.T) {
	// a finder pattern is three rectangles of rows and columns, not 33 modules
	s := mustEncode(t, QR, "1", &Options{Quiet: NoQuiet})
	rects := s.Rects()
	dark := 0
	for y := 0; y < s.H; y++ {
		for x := 0; x < s.W; x++ {
			if s.At(x, y) {
				dark++
			}
		}
	}
	covered := 0
	for _, r := range rects {
		covered += r[2] * r[3]
		for y := r[1]; y < r[1]+r[3]; y++ {
			for x := r[0]; x < r[0]+r[2]; x++ {
				if !s.At(x, y) {
					t.Fatalf("the rectangle %v holds a light module", r)
				}
			}
		}
	}
	if covered != dark || len(rects) >= dark || rects[0] != [4]int{0, 0, 7, 1} {
		t.Errorf("%d rectangles of %d modules for %d dark modules; the first is %v", len(rects), covered, dark, rects[0])
	}
	// the bars of a barcode are one rectangle each, as tall as the box
	b := mustEncode(t, EAN13, "4901234567894", nil)
	p := b.Path(10, 20, 113, 30)
	if len(p.Verbs) != 30 || p.Verbs[0] != bdf.VerbRect || p.Args[0] != 10+11 || p.Args[1] != 20 || p.Args[2] != 1 || p.Args[3] != 30 {
		t.Errorf("%d bars; the first is %v", len(p.Verbs), p.Args[:4])
	}
	o := b.Object(113, 30)
	if o.BBox.W != 113 || o.BBox.H != 30 || len(o.Deps()) != 0 {
		t.Errorf("the object is %v", o.BBox)
	}
}

// FuzzEncode makes symbols of any text: a QR code says the text it was
// made of, and no text makes a symbology fail of itself.
func FuzzEncode(f *testing.F) {
	f.Add(uint8(0), "HELLO WORLD", uint8(1))
	f.Add(uint8(0), "日本語 text 0123", uint8(3))
	f.Add(uint8(1), "4901234567894", uint8(0))
	f.Add(uint8(5), "AB12345\tx", uint8(0))
	f.Add(uint8(6), "(01)04912345678904(10)A", uint8(0))
	f.Add(uint8(9), "A12-34A", uint8(0))
	f.Fuzz(func(t *testing.T, kind uint8, text string, level uint8) {
		kinds := Kinds()
		k := kinds[int(kind)%len(kinds)]
		s, err := Encode(k, text, &Options{Level: "LMQH"[level%4 : level%4+1], Quiet: int(level%8) - 1})
		if err != nil {
			return
		}
		if s.W <= 0 || s.H <= 0 || len(s.Rects()) == 0 {
			t.Fatalf("%s %q: a symbol of %d × %d modules", k, text, s.W, s.H)
		}
		if k == QR {
			quiet := int(level%8) - 1
			switch {
			case quiet < 0:
				quiet = 0
			case quiet == 0:
				quiet = 4
			}
			if got, _, lv, _ := readQR(t, s, quiet); got != text || lv != int(level%4) {
				t.Fatalf("the QR code of %q says %q at level %d", text, got, lv)
			}
		}
	})
}
