package barcode

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// TestDump writes symbols for a reader of barcodes to read (a check from
// outside, by hand: BARCODE_DUMP names the file).
func TestDump(t *testing.T) {
	out := os.Getenv("BARCODE_DUMP")
	if out == "" {
		t.Skip()
	}
	type sym struct {
		Kind, Text, Level, Want string
		W, H                    int
		Rows                    []string
	}
	var syms []sym
	add := func(kind, text, level string) {
		s, err := Encode(kind, text, &Options{Level: level})
		if err != nil {
			t.Errorf("%s %q: %v", kind, text, err)
			return
		}
		e := sym{Kind: s.Kind, Text: text, Level: level, Want: s.Text, W: s.W, H: s.H}
		for y := 0; y < s.H; y++ {
			var b strings.Builder
			for x := 0; x < s.W; x++ {
				if s.At(x, y) {
					b.WriteByte('1')
				} else {
					b.WriteByte('0')
				}
			}
			e.Rows = append(e.Rows, b.String())
		}
		syms = append(syms, e)
	}
	long := strings.Repeat("The quick brown fox jumps over the lazy dog. 0123456789 ", 40)
	for _, level := range []string{"L", "M", "Q", "H"} {
		for _, text := range []string{"0", "01234567", "HELLO WORLD", "https://example.com/orders/1024?x=1", "日本語のテキスト、帳票", strings.Repeat("7", 300), strings.Repeat("AB12 ", 120),
			long[:100], long[:400], long[:900], long[:1200]} {
			add(QR, text, level)
		}
	}
	// every version at every level: a text of bytes that fills it
	for lv, level := range []string{"L", "M", "Q", "H"} {
		for ver := 1; ver <= 40; ver++ {
			n := (qrDataCodewords(ver, lv)*8 - 4 - qrCountBits(qrByte, ver)) / 8
			text := strings.Repeat("abcdefghij klmnopqrstuvwxyz,", n/28+1)[:n]
			add(QR, text, level)
			if got := (syms[len(syms)-1].W - 8 - 17) / 4; got != ver {
				t.Errorf("level %s: %d bytes make version %d, want %d", level, n, got, ver)
			}
		}
	}
	// every mask, to compare with another writer of QR codes
	for mask := 0; mask < 8; mask++ {
		for i, text := range []string{"01234567890123456789", "HELLO WORLD 123 $%*+-./:", "https://example.com/orders/1024", "日本語のテキスト", long[:300], long[:1500], strings.Repeat("5", 2000)} {
			level := []string{"M", "Q", "H", "L", "Q", "M", "H"}[i]
			q, err := encodeQR(text, level, 4, mask)
			if err != nil {
				t.Fatal(err)
			}
			e := sym{Kind: "qr-mask", Text: text, Level: level, Want: string(rune('0' + mask)), W: q.W, H: q.H}
			for y := 0; y < q.H; y++ {
				var b strings.Builder
				for x := 0; x < q.W; x++ {
					if q.At(x, y) {
						b.WriteByte('1')
					} else {
						b.WriteByte('0')
					}
				}
				e.Rows = append(e.Rows, b.String())
			}
			syms = append(syms, e)
		}
	}
	add(QR, long[:2200], "L")
	add(QR, strings.Repeat("9", 7089), "L") // all a QR code holds
	add(EAN13, "4901234567894", "")
	add(EAN13, "490123456789", "")
	add(EAN13, "0012345678905", "")
	add(EAN13, "9784873119038", "")
	add(EAN8, "4901234", "")
	add(EAN8, "96385074", "")
	add(UPCA, "03600029145", "")
	add(Code39, "CODE-39 TEST.$/+%", "")
	add(Code39, "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	for _, text := range []string{"Code 128", "1234567890", "12345", "A1234B", "ab123456cd", "Hello, World! ~{|}", "\tTAB\x01ctl", "ABC\x1fabc", "AB12345", "123", "12", "9", "x1234", "1234x", "!\"#$%&'()*+,-./:;<=>?@[\\]^_`"} {
		add(Code128, text, "")
	}
	for v := 0; v < 100; v += 7 { // every pair of digits
		add(Code128, string([]byte{'0' + byte(v/10), '0' + byte(v%10), '0' + byte((v+1)%100/10), '0' + byte((v+1)%100%10), '0' + byte((v+2)%100/10), '0' + byte((v+2)%100%10),
			'0' + byte((v+3)%100/10), '0' + byte((v+3)%100%10), '0' + byte((v+4)%100/10), '0' + byte((v+4)%100%10), '0' + byte((v+5)%100/10), '0' + byte((v+5)%100%10), '0' + byte((v+6)%100/10), '0' + byte((v+6)%100%10)}), "")
	}
	add(GS1128, "(01)04912345678904(15)261231(10)LOT123", "")
	add(GS1128, "(10)AB12(21)12345678", "")
	add(GS1128, "(00)123456789012345675", "")
	add(ITF, "1234567890", "")
	add(ITF, "12345", "")
	add(ITF14, "1490123456789", "")
	add(ITF14, "04912345678904", "")
	add(NW7, "A1234567890A", "")
	add(NW7, "12-34$56:78/9.0+", "")
	add(NW7, "b123456d", "")
	b, _ := json.Marshal(syms)
	if err := os.WriteFile(out, b, 0o644); err != nil {
		t.Fatal(err)
	}
}
