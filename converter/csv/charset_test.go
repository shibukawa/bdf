package csv

import (
	"strings"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	xunicode "golang.org/x/text/encoding/unicode"
)

func encode(t *testing.T, enc encoding.Encoding, s string) []byte {
	t.Helper()
	b, err := enc.NewEncoder().Bytes([]byte(s))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestCharset(t *testing.T) {
	ja := "名前,説明\n山田,ひらがなとカタカナの文章です。\n佐藤,漢字だけ\n"
	names := "氏名,住所\n山田太郎,東京都千代田区\n佐藤花子,大阪府大阪市\n" // no kana
	latin := "Name;City;Note\nMüller;Köln;Straße\nRenée;Besançon;“naïve” café – ok\n"
	for _, c := range []struct {
		name string
		data []byte
		want string
		bom  bool
		text string
	}{
		{"ascii", []byte("a,b\n1,2\n"), "UTF-8", false, "a,b\n1,2\n"},
		{"utf-8", []byte(ja), "UTF-8", false, ja},
		{"utf-8 bom", append([]byte{0xEF, 0xBB, 0xBF}, ja...), "UTF-8", true, ja},
		{"utf-16le bom", append([]byte{0xFF, 0xFE}, encode(t, xunicode.UTF16(xunicode.LittleEndian, xunicode.IgnoreBOM), ja)...), "UTF-16LE", true, ja},
		{"utf-16be bom", append([]byte{0xFE, 0xFF}, encode(t, xunicode.UTF16(xunicode.BigEndian, xunicode.IgnoreBOM), ja)...), "UTF-16BE", true, ja},
		{"utf-32le bom", []byte{0xFF, 0xFE, 0, 0, 'a', 0, 0, 0, ',', 0, 0, 0, 'b', 0, 0, 0}, "UTF-32LE", true, "a,b"},
		{"utf-16le", encode(t, xunicode.UTF16(xunicode.LittleEndian, xunicode.IgnoreBOM), ja), "UTF-16LE", false, ja},
		{"utf-16be", encode(t, xunicode.UTF16(xunicode.BigEndian, xunicode.IgnoreBOM), latin), "UTF-16BE", false, latin},
		{"shift_jis", encode(t, japanese.ShiftJIS, ja), "Shift_JIS", false, ja},
		{"shift_jis kanji", encode(t, japanese.ShiftJIS, names), "Shift_JIS", false, names},
		{"shift_jis half-width", encode(t, japanese.ShiftJIS, "ｶﾅ,ﾒｲ\nﾔﾏﾀﾞ,ﾀﾛｳ\n"), "Shift_JIS", false, "ｶﾅ,ﾒｲ\nﾔﾏﾀﾞ,ﾀﾛｳ\n"},
		{"euc-jp", encode(t, japanese.EUCJP, ja), "EUC-JP", false, ja},
		{"euc-jp kanji", encode(t, japanese.EUCJP, names), "EUC-JP", false, names},
		{"iso-2022-jp", encode(t, japanese.ISO2022JP, ja), "ISO-2022-JP", false, ja},
		{"windows-1252", encode(t, charmap.Windows1252, latin), "WINDOWS-1252", false, latin},
		// a UTF-8 file with a stray Latin-1 byte stays UTF-8
		{"utf-8 stray", []byte(strings.Repeat("名前,説明\n", 5) + "caf\xe9\n"), "UTF-8", false, strings.Repeat("名前,説明\n", 5) + "caf�\n"},
	} {
		cs, bom := bomCharset(c.data)
		if !bom {
			cs = detectCharset(c.data)
		}
		if cs.name != c.want || bom != c.bom {
			t.Errorf("%s: detected %s (bom %v), want %s (bom %v)", c.name, cs.name, bom, c.want, c.bom)
			continue
		}
		if got := decode(c.data, cs); got != c.text {
			t.Errorf("%s: decoded %q, want %q", c.name, got, c.text)
		}
	}
}

func TestNamedCharset(t *testing.T) {
	for label, want := range map[string]string{"utf-8": "UTF-8", "UTF8": "UTF-8", "sjis": "Shift_JIS", "cp932": "Shift_JIS",
		"euc-kr": "EUC-KR", "gb2312": "GBK", "big5": "BIG5", "latin1": "WINDOWS-1252"} {
		cs, err := namedCharset(label)
		if err != nil || cs.name != want {
			t.Errorf("%s: %q %v, want %q", label, cs.name, err, want)
		}
	}
	if _, err := namedCharset("klingon"); err == nil {
		t.Error("an unknown encoding was accepted")
	}
}
