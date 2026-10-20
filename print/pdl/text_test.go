package pdl

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/form/textpart"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
)

// formDoc opens a form filled from a template of text (testdata/form).
func formDoc(t *testing.T, name string) (*bdf.Document, *textpart.Image) {
	t.Helper()
	f, err := os.Open("../../testdata/form/" + name)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	fi, _ := f.Stat()
	r, err := bdf.OpenSingle(f, fi.Size())
	if err != nil {
		t.Fatal(err)
	}
	d, err := r.ToDocument()
	if err != nil {
		t.Fatal(err)
	}
	im, err := textpart.FromDocument(d)
	if err != nil {
		t.Fatal(err)
	}
	return d, im
}

// imageDoc is a document of a print image, its pages empty.
func imageDoc(t *testing.T, im *textpart.Image) *bdf.Document {
	t.Helper()
	d := bdf.NewDocument()
	v := d.NewView("v", bdf.ViewFixed, "")
	for range im.Pages {
		v.AddPage(72*8, 72*11)
	}
	if _, err := im.AddTo(d); err != nil {
		t.Fatal(err)
	}
	return d
}

// textGrid is the pages of characters a reader of the text mode prints:
// for each line, the characters at their columns.
type textGrid struct {
	pages [][]map[int]string
}

func (g *textGrid) put(page, line, col int, s string) {
	for len(g.pages) <= page {
		g.pages = append(g.pages, nil)
	}
	for len(g.pages[page]) <= line {
		g.pages[page] = append(g.pages[page], map[int]string{})
	}
	g.pages[page][line][col] = s
}

// text writes the pages as textpart.Image.Text writes them, its margins
// included.
func (g *textGrid) text(ambiguousWide bool) string {
	var b strings.Builder
	for pi, p := range g.pages {
		if pi > 0 {
			b.WriteByte('\f')
		}
		for _, l := range p {
			var line strings.Builder
			end := 0
			for c := range l {
				end = max(end, c+1)
			}
			for c := 0; c < end; {
				s, ok := l[c]
				if !ok {
					line.WriteByte(' ')
					c++
					continue
				}
				line.WriteString(s)
				r := []rune(s)[0]
				c += max(textpart.Width(r, ambiguousWide), 1)
			}
			b.WriteString(strings.TrimRight(line.String(), " "))
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// trimPages drops the empty lines at the end of each page of a text.
func trimPages(s string) string {
	pages := strings.Split(s, "\f")
	for i, p := range pages {
		pages[i] = strings.TrimRight(p, "\n")
	}
	return strings.Join(pages, "\f")
}

func jisRune(j1, j2 byte) string {
	s1 := (int(j1)+1)/2 + 0x70
	if j1 > 0x5e {
		s1 += 0x40
	}
	s2 := int(j2) + 0x7e
	if j1%2 == 1 {
		s2 = int(j2) + 0x1f
		if j2 >= 0x60 {
			s2++
		}
	}
	b, err := japanese.ShiftJIS.NewDecoder().Bytes([]byte{byte(s1), byte(s2)})
	if err != nil {
		return "�"
	}
	return string(b)
}

// readESCPText prints what writeESCPText writes: the characters at their
// columns (the left margin included) and lines, as ESC/P-J84 prints them
// (Canon's and Brother's emulations): ESC W is FS ! bit 2, ESC w enlarges
// upward from the base line, FS W (FS ! bit 3) prints the upper halves of
// Kanji, feeds 24/180 inch and prints the lower halves and the rest.
func readESCPText(data []byte, ambiguousWide bool, lpi float64) (*textGrid, error) {
	g := &textGrid{}
	lineDots := escpDots / lpi
	page, col, margin, y := 0, 0, 0, 0.0
	kanji, half, wide, tall, big, japan := false, false, false, false, false, false
	type placed struct {
		col       int
		s         string
		tall, big bool
	}
	var buf []placed
	lf := func() error {
		anyTall, anyBig := false, false
		for _, p := range buf {
			anyTall = anyTall || p.tall
			anyBig = anyBig || p.big
		}
		at := y
		switch {
		case anyBig:
			at = y + escpKanjiDots - lineDots // the lower halves are on the next line
		case anyTall:
			at = y - lineDots // enlarged upward into the line above
		}
		line := int(math.Round(at / lineDots))
		if len(buf) > 0 && math.Abs(float64(line)*lineDots-at) > 0.5 {
			return fmt.Errorf("a line at %g/180 inch, off the grid", at)
		}
		for _, p := range buf {
			g.put(page, line, p.col, p.s)
		}
		buf = nil
		y += lineDots
		if anyBig {
			y += escpKanjiDots
		}
		return nil
	}
	need := func(i, n int) error {
		if i+n > len(data) {
			return fmt.Errorf("a command cut at %d", i)
		}
		return nil
	}
	for i := 0; i < len(data); {
		c := data[i]
		switch {
		case c == 0x1b:
			if err := need(i, 2); err != nil {
				return nil, err
			}
			switch data[i+1] {
			case '@', 'P', 'M', 'g', '2', '0', 'O':
				i += 2
			case 'R':
				japan = data[i+2] == 8
				i += 3
			case '3', 'N':
				i += 3
			case 'J':
				if len(buf) > 0 {
					return nil, fmt.Errorf("ESC J after characters at %d", i)
				}
				y += float64(data[i+2])
				i += 3
			case 'l':
				margin = int(data[i+2])
				col = margin
				i += 3
			case 'W':
				wide = data[i+2] == 1
				i += 3
			case 'w':
				tall = data[i+2] == 1
				i += 3
			case 'C':
				if data[i+2] == 0 {
					i += 4
				} else {
					i += 3
				}
			default:
				return nil, fmt.Errorf("ESC %q at %d", data[i+1], i)
			}
		case c == 0x1c:
			if err := need(i, 2); err != nil {
				return nil, err
			}
			switch data[i+1] {
			case '&':
				kanji = true
				i += 2
			case '.':
				kanji = false
				i += 2
			case '!':
				n := data[i+2]
				half, wide, big = n&2 != 0, n&4 != 0, n&8 != 0
				i += 3
			case 'S', 'T':
				i += 4
			default:
				return nil, fmt.Errorf("FS %q at %d", data[i+1], i)
			}
		case c == 0x12 || c == 0x0f:
			i++
		case c == '\r':
			col = margin
			i++
		case c == '\n':
			if err := lf(); err != nil {
				return nil, err
			}
			i++
		case c == 0x0c:
			if len(buf) > 0 {
				return nil, fmt.Errorf("FF after characters at %d", i)
			}
			page, col, y = page+1, margin, 0
			i++
		case kanji:
			if err := need(i, 2); err != nil {
				return nil, err
			}
			w := 2
			if half {
				w = 1
			}
			buf = append(buf, placed{col, jisRune(data[i], data[i+1]), false, big})
			if wide || big {
				w *= 2
			}
			col += w
			i += 2
		case c >= 0x20:
			s := string(rune(c))
			switch {
			case c >= 0xa1 && c <= 0xdf:
				s = string(rune(0xff61 + int(c) - 0xa1))
			case c == '\\' && japan:
				s = "¥"
			case c > 0x7e:
				return nil, fmt.Errorf("a byte 0x%02x at %d", c, i)
			}
			if c != ' ' {
				buf = append(buf, placed{col, s, tall, false})
			}
			col++
			if wide || big {
				col++ // (FS W widens the characters of a byte)
			}
			i++
		default:
			return nil, fmt.Errorf("a byte 0x%02x at %d", c, i)
		}
	}
	return g, nil
}

// readESCPOSText prints what writeESCPOSText writes, at dots of a pitch
// and a line of the grid: a receipt ends at each cut.
func readESCPOSText(data []byte, dpi int, cpi, lpi float64) (*textGrid, int, error) {
	g := &textGrid{}
	pitch, lineDots := float64(dpi)/cpi, float64(dpi)/lpi
	page, x, y, adv, kspace, size, cuts := 0, 0, 0, 12, 0, byte(0), 0
	kanji, japan, table := false, false, 1
	type placed struct {
		x int
		s string
	}
	var buf []placed
	flush := func() error {
		line := int(math.Round(float64(y) / lineDots))
		if math.Abs(float64(line)*lineDots-float64(y)) >= 1 {
			return fmt.Errorf("a line at %d dots, off the grid", y)
		}
		for _, p := range buf {
			col := int(math.Round(float64(p.x) / pitch))
			if int(math.Round(float64(col)*pitch)) != p.x {
				return fmt.Errorf("%q at %d dots, off the grid", p.s, p.x)
			}
			g.put(page, line, col, p.s)
		}
		buf = nil
		x = 0
		return nil
	}
	for i := 0; i < len(data); {
		c := data[i]
		n := func(k int) int {
			if i+k >= len(data) {
				return -1
			}
			return int(data[i+k])
		}
		mul := 1
		if size&0x10 != 0 {
			mul = 2
		}
		switch {
		case c == 0x1b && n(1) == '@':
			i += 2
		case c == 0x1b && n(1) == 'R':
			japan = n(2) == 8
			i += 3
		case c == 0x1b && n(1) == 'M':
			i += 3
		case c == 0x1b && n(1) == ' ':
			adv = 12 + n(2)
			i += 3
		case c == 0x1b && n(1) == 't':
			table = n(2)
			i += 3
		case c == 0x1b && n(1) == '$':
			x = n(2) | n(3)<<8
			i += 4
		case c == 0x1b && n(1) == 'J':
			if len(buf) > 0 {
				if err := flush(); err != nil {
					return nil, 0, err
				}
			}
			y += n(2)
			i += 3
		case c == 0x1b && n(1) == 'p':
			i += 5
		case c == 0x1d && n(1) == 'P':
			if n(2) != dpi {
				return nil, 0, fmt.Errorf("GS P %d", n(2))
			}
			i += 4
		case c == 0x1d && n(1) == '!':
			size = byte(n(2))
			i += 3
		case c == 0x1d && n(1) == 'V':
			page, y, cuts = page+1, 0, cuts+1
			i += 4
		case c == 0x1c && n(1) == 'C':
			i += 3
		case c == 0x1c && n(1) == '!':
			if want := byte(size>>2&4 | size<<3&8); byte(n(2)) != want {
				return nil, 0, fmt.Errorf("FS ! %d after GS ! %d", n(2), size)
			}
			i += 3
		case c == 0x1c && n(1) == 'S':
			kspace = n(2) + n(3)
			i += 4
		case c == 0x1c && n(1) == '&':
			kanji = true
			i += 2
		case c == 0x1c && n(1) == '.':
			kanji = false
			i += 2
		case kanji && c >= 0x21 && n(1) >= 0x21:
			buf = append(buf, placed{x, jisRune(c, byte(n(1)))})
			x += (24 + kspace) * mul
			i += 2
		case c >= 0x20 && !kanji:
			s := string(rune(c))
			switch {
			case c == '\\' && japan:
				s = "¥"
			case c > 0x7e && table == 1:
				s = string(rune(0xff61 + int(c) - 0xa1))
			case c > 0x7e && table == 2:
				s = string(charmap.CodePage850.DecodeByte(c))
			case c > 0x7e:
				s = string(charmap.CodePage437.DecodeByte(c))
			}
			if c != ' ' {
				buf = append(buf, placed{x, s})
			}
			x += adv * mul
			i++
		default:
			return nil, 0, fmt.Errorf("a byte 0x%02x at %d", c, i)
		}
	}
	if len(buf) > 0 {
		return nil, 0, fmt.Errorf("characters not printed")
	}
	return g, cuts, nil
}

func TestTextESCP(t *testing.T) {
	prof, _ := LookupProfile("vp-f2000")
	for _, name := range []string{"text-slip.bdf", "text-rows.bdf", "text-pages.bdf"} {
		for _, noPaper := range []bool{false, true} {
			d, im := formDoc(t, name)
			o := &Options{Profile: prof}
			if noPaper {
				o.SkipRoles = []string{bdf.RoleBackground}
			}
			var buf bytes.Buffer
			res, err := Write(&buf, d, d.Views[0], o)
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			if res.Pages != len(im.Pages) || len(res.Unsupported) > 0 {
				t.Errorf("%s: %+v", name, res)
			}
			data := buf.Bytes()
			switch name {
			case "text-slip.bdf":
				// 納品書 of double size: lowered 6/180 inch, FS ! 8
				if !bytes.Contains(data, []byte("\x1bJ\x06")) || !bytes.Contains(data, []byte("\x1c!\x08\x1c&")) {
					t.Errorf("%s: no Kanji of double size", name)
				}
			case "text-pages.bdf":
				// 前頁より繰越 of double height: printed at its size
				if !strings.Contains(strings.Join(res.Warnings, "\n"), "double height alone") || !bytes.Contains(data, []byte("\x1bw\x01")) {
					t.Errorf("%s: %q", name, res.Warnings)
				}
			}
			if bytes.Contains(data, []byte("\x1b*")) {
				t.Errorf("%s: dots in the text mode", name)
			}
			length := []byte{0x1b, 'C', byte(im.Lines)}
			if !bytes.Contains(data, length) || im.Skip > 0 && !bytes.Contains(data, []byte{0x1b, 'N', byte(im.Skip)}) {
				t.Errorf("%s: no page length or skip", name)
			}
			if got := bytes.Count(data, []byte("\r\n\x0c")); got != len(im.Pages) {
				t.Errorf("%s: %d form feeds", name, got)
			}
			g, err := readESCPText(data, im.AmbiguousWide, im.LPI)
			if err != nil {
				t.Fatalf("%s: %v", name, err)
			}
			want := trimPages(im.Text(textpart.TextOptions{Margins: true, NoPaper: noPaper}))
			if got := trimPages(g.text(im.AmbiguousWide)); got != want {
				t.Errorf("%s (no paper %v): printed\n%s\nwant\n%s", name, noPaper, got, want)
			}
		}
	}
}

func TestTextESCPHalfWidth(t *testing.T) {
	// the box drawing of a template of "ambiguous: narrow" is Kanji at
	// half width, of "ambiguous: wide" full width
	_, im := formDoc(t, "text-slip.bdf")
	if im.AmbiguousWide {
		t.Fatal("the slip is of ambiguous: narrow")
	}
	d := imageDoc(t, im)
	prof, _ := LookupProfile("vp-f2000")
	var buf bytes.Buffer
	if _, err := Write(&buf, d, d.Views[0], &Options{Profile: prof, Settings: Settings{Mode: ModeText}}); err != nil {
		t.Fatal(err)
	}
	// FS ! 2 (half width), FS &, ┌ (JIS 0x2823) and ─ (0x2821); FS T 3 3
	// at 10 CPI, FS S 6 6
	for _, want := range []string{"\x1c!\x02\x1c&\x28\x23\x28\x21", "\x1cT\x03\x03", "\x1cS\x06\x06"} {
		if !bytes.Contains(buf.Bytes(), []byte(want)) {
			t.Errorf("no %q", want)
		}
	}
}

func TestTextESCPOS(t *testing.T) {
	for _, name := range []string{"tm-t88", "pos-80"} {
		d, im := formDoc(t, "text-receipt.bdf")
		prof, _ := LookupProfile(name)
		var buf bytes.Buffer
		o := &Options{Profile: prof, Settings: Settings{Mode: ModeText, Drawer: 1}, Copies: 2}
		res, err := Write(&buf, d, d.Views[0], o)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if res.Pages != 2 {
			t.Errorf("%s: %+v", name, res)
		}
		data := buf.Bytes()
		if bytes.Contains(data, []byte("\x1dv0")) {
			t.Errorf("%s: dots in the text mode", name)
		}
		if !bytes.HasSuffix(data, []byte("\x1bp\x00\x19\xfa")) {
			t.Errorf("%s: no cash drawer at the end", name)
		}
		g, cuts, err := readESCPOSText(data, prof.DPI, im.CPI, im.LPI)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if cuts != 2 {
			t.Errorf("%s: %d cuts", name, cuts)
		}
		text := trimPages(im.Text(textpart.TextOptions{Margins: true}))
		want := text + "\f" + text
		if got := trimPages(g.text(im.AmbiguousWide)); got != want {
			t.Errorf("%s: printed\n%q\nwant\n%q", name, got, want)
		}
	}
}

func TestTextModes(t *testing.T) {
	d, im := formDoc(t, "text-receipt.bdf")
	pos, _ := LookupProfile("pos-80")
	escp, _ := LookupProfile("vp-f2000")
	zd, _ := LookupProfile("zd421-203")
	write := func(p Profile, s Settings, doc *bdf.Document) ([]byte, *Result, error) {
		var buf bytes.Buffer
		res, err := Write(&buf, doc, doc.Views[0], &Options{Profile: p, Settings: s, NoSystemFonts: true})
		return buf.Bytes(), res, err
	}
	// receipts are dots unless the text mode is asked
	if data, _, err := write(pos, Settings{}, d); err != nil || !bytes.Contains(data, []byte("\x1dv0")) {
		t.Errorf("ESC/POS by default: %v", err)
	}
	if data, _, err := write(escp, Settings{Mode: ModeRaster}, d); err != nil || !bytes.Contains(data, []byte("\x1b*")) {
		t.Errorf("ESC/P in the raster mode: %v", err)
	}
	if _, _, err := write(zd, Settings{Mode: ModeText}, d); err == nil {
		t.Error("ZPL in the text mode")
	}
	plain, _ := testDoc(t, 200, 200, 1)
	if _, _, err := write(escp, Settings{Mode: ModeText}, plain); err == nil {
		t.Error("the text mode of a document without a print image")
	}
	// a pitch ESC/P has not: dots by default, an error in the text mode
	odd := *im
	odd.CPI = 13
	od := imageDoc(t, &odd)
	if data, res, err := write(escp, Settings{}, od); err != nil || bytes.Contains(data, []byte("\x1c&")) || !strings.Contains(strings.Join(res.Warnings, "\n"), "printed as dots") {
		t.Errorf("13 CPI by default: %v %+v", err, res)
	}
	if _, _, err := write(escp, Settings{Mode: ModeText}, od); err == nil {
		t.Error("13 CPI in the text mode")
	}
	// wider than the roll
	p58, _ := LookupProfile("pos-58")
	if _, _, err := write(p58, Settings{Mode: ModeText}, d); err == nil || !strings.Contains(err.Error(), "wide") {
		t.Errorf("42 columns of 15 CPI on 58 mm: %v", err)
	}
}

func TestTextEncode(t *testing.T) {
	im := &textpart.Image{Version: 1, CPI: 10, LPI: 6, Columns: 40, Pages: []textpart.Page{{Lines: []textpart.Line{
		{Text: "┏━┓ ○ ¥ ｱ 亜 😀"},
	}}}}
	for _, c := range []struct {
		lang  string
		kinds []int
		warns int
	}{
		{ESCP, []int{gHalfKanji, gHalfKanji, gHalfKanji, gASCII, gHalfKanji, gASCII, gYen, gASCII, gKana, gASCII, gKanji, gASCII, gKanji}, 1},
		{ESCPOS, []int{gPC437, gPC437, gPC437, gASCII, gASCII, gASCII, gYen, gASCII, gKana, gASCII, gKanji, gASCII, gKanji}, 3},
	} {
		res := &Result{}
		e := newTextEncoder(im, &Options{Profile: Profile{Language: c.lang}}, res)
		gs := e.line(im.Pages[0].Lines[0])
		var kinds []int
		for _, g := range gs {
			kinds = append(kinds, g.kind)
		}
		if fmt.Sprint(kinds) != fmt.Sprint(c.kinds) || len(res.Warnings) != c.warns {
			t.Errorf("%s: %v %q", c.lang, kinds, res.Warnings)
		}
		if c.lang == ESCPOS && (gs[0].b[0] != 0xda || gs[1].b[0] != 0xc4) {
			t.Errorf("┌─ of code page 437: % x % x", gs[0].b[0], gs[1].b[0])
		}
		if last := gs[len(gs)-1]; last.b != [2]byte{0x22, 0x2e} {
			t.Errorf("%s: 😀 as % x", c.lang, last.b)
		}
	}
	enc := japanese.ShiftJIS.NewEncoder()
	for r, want := range map[rune][2]byte{'亜': {0x30, 0x21}, '─': {0x28, 0x21}, '①': {0x2d, 0x21}, '腕': {0x4f, 0x53}, '弌': {0x50, 0x21}, '熙': {0x74, 0x26}} {
		if got, ok := jisX0208(enc, r); !ok || got != want {
			t.Errorf("%c: % x %v", r, got, ok)
		}
		if s := jisRune(want[0], want[1]); s != string(r) {
			t.Errorf("%x %x: %q", want[0], want[1], s)
		}
	}
	if _, ok := jisX0208(enc, '髙'); ok {
		t.Error("髙 is of IBM's extensions")
	}
}
