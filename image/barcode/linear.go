package barcode

import (
	"fmt"
	"strings"
)

// --- JAN / EAN / UPC ---

// The bars of the digits of EAN codes at the left of the symbol (odd
// parity), as widths of space, bar, space, bar. At the right the bars are
// those of the same widths with bars and spaces changed over, and the left
// digits of even parity have the widths in the other order.
var eanDigits = [10][4]int{
	{3, 2, 1, 1}, {2, 2, 2, 1}, {2, 1, 2, 2}, {1, 4, 1, 1}, {1, 1, 3, 2},
	{1, 2, 3, 1}, {1, 1, 1, 4}, {1, 3, 1, 2}, {1, 2, 1, 3}, {3, 1, 1, 2},
}

// eanParity says which of the six left digits of an EAN-13 code have even
// parity, by its first digit (which has no bars of its own).
var eanParity = [10]string{"LLLLLL", "LLGLGG", "LLGGLG", "LLGGGL", "LGLLGG", "LGGLLG", "LGGGLL", "LGLGLG", "LGLGGL", "LGGLGL"}

func encodeEAN(kind, text string, quiet int) (*Symbol, error) {
	n, left, right := 13, 11, 7
	switch kind {
	case EAN8:
		n, left = 8, 7
	case UPCA:
		n, left, right = 12, 9, 9
	}
	code, err := gtin(kind, text, n)
	if err != nil {
		return nil, err
	}
	// (a UPC-A code is the EAN-13 code with a 0 before it)
	digs, parity := code, "LLLL"
	switch kind {
	case UPCA:
		parity = eanParity[0]
	case EAN13:
		digs, parity = code[1:], eanParity[code[0]-'0']
	}
	half := len(digs) / 2
	widths := []int{1, 1, 1} // the guard: bar, space, bar
	for i := 0; i < half; i++ {
		d := eanDigits[digs[i]-'0']
		if parity[i] == 'G' {
			d = [4]int{d[3], d[2], d[1], d[0]}
		}
		widths = append(widths, d[:]...)
	}
	widths = append(widths, 1, 1, 1, 1, 1) // the centre guard: space, bar, space, bar, space
	for i := half; i < len(digs); i++ {
		d := eanDigits[digs[i]-'0']
		widths = append(widths, d[:]...) // bar, space, bar, space
	}
	widths = append(widths, 1, 1, 1)
	if quiet >= 0 {
		left, right = quiet, quiet
	}
	s := linear(widths, left, right)
	s.Text = code
	return s, nil
}

// --- Code 39 ---

// The elements of the characters of Code 39: five bars and the four spaces
// between them, a bar first; 1 is a wide element.
var code39Chars = map[byte]string{
	'0': "000110100", '1': "100100001", '2': "001100001", '3': "101100000", '4': "000110001",
	'5': "100110000", '6': "001110000", '7': "000100101", '8': "100100100", '9': "001100100",
	'A': "100001001", 'B': "001001001", 'C': "101001000", 'D': "000011001", 'E': "100011000",
	'F': "001011000", 'G': "000001101", 'H': "100001100", 'I': "001001100", 'J': "000011100",
	'K': "100000011", 'L': "001000011", 'M': "101000010", 'N': "000010011", 'O': "100010010",
	'P': "001010010", 'Q': "000000111", 'R': "100000110", 'S': "001000110", 'T': "000010110",
	'U': "110000001", 'V': "011000001", 'W': "111000000", 'X': "010010001", 'Y': "110010000",
	'Z': "011010000", '-': "010000101", '.': "110000100", ' ': "011000100", '*': "010010100",
	'$': "010101000", '/': "010100010", '+': "010001010", '%': "000101010",
}

func encodeCode39(text string, quiet int) (*Symbol, error) {
	// (the asterisks around a text are the start and stop characters, which the symbol has anyway)
	t := strings.TrimSuffix(strings.TrimPrefix(text, "*"), "*")
	if t == "" {
		return nil, fmt.Errorf("barcode: nothing to make a code39 symbol of")
	}
	var widths []int
	for i, c := range []byte("*" + t + "*") {
		p, ok := code39Chars[c]
		if !ok || c == '*' && i != 0 && i != len(t)+1 {
			return nil, fmt.Errorf("barcode: code39 has no %q (it has digits, capital letters, the space and - . $ / + %%)", string(rune(c)))
		}
		if i > 0 {
			widths = append(widths, 1) // the gap between characters
		}
		for _, e := range []byte(p) {
			widths = append(widths, 1+2*int(e-'0')) // a wide element is three narrow ones
		}
	}
	if quiet < 0 {
		quiet = 10
	}
	s := linear(widths, quiet, quiet)
	s.Text = t
	return s, nil
}

// --- Code 128 and GS1-128 ---

// The widths of the bars and spaces of the characters of Code 128, a bar
// first: the values 0 to 102, the start characters A, B and C, and the
// stop character.
var code128Chars = [107]string{
	"212222", "222122", "222221", "121223", "121322", "131222", "122213", "122312", "132212", "221213",
	"221312", "231212", "112232", "122132", "122231", "113222", "123122", "123221", "223211", "221132",
	"221231", "213212", "223112", "312131", "311222", "321122", "321221", "312212", "322112", "322211",
	"212123", "212321", "232121", "111323", "131123", "131321", "112313", "132113", "132311", "211313",
	"231113", "231311", "112133", "112331", "132131", "113123", "113321", "133121", "313121", "211331",
	"231131", "213113", "213311", "213131", "311123", "311321", "331121", "312113", "312311", "332111",
	"314111", "221411", "431111", "111224", "111422", "121124", "121421", "141122", "141221", "112214",
	"112412", "122114", "122411", "142112", "142211", "241211", "221114", "413111", "241112", "134111",
	"111242", "121142", "121241", "114212", "124112", "124211", "411212", "421112", "421211", "212141",
	"214121", "412121", "111143", "111341", "131141", "114113", "114311", "411113", "411311", "113141",
	"114131", "311141", "411131", "211412", "211214", "211232", "2331112",
}

// Values of Code 128 that are not characters of data.
const (
	c128FNC1   = 102
	c128ToC    = 99
	c128ToB    = 100
	c128ToA    = 101
	c128StartA = 103
	c128StartB = 104
	c128StartC = 105
	c128Stop   = 106
	// fnc1 stands for the function character 1 among the bytes of a text
	fnc1 = 0x80
)

// gs1Fixed are the application identifiers of GS1 whose fields have a
// fixed length, by their first two digits: no separator follows them.
var gs1Fixed = map[string]bool{"00": true, "01": true, "02": true, "03": true, "04": true, "11": true, "12": true, "13": true, "14": true,
	"15": true, "16": true, "17": true, "18": true, "19": true, "20": true, "31": true, "32": true, "33": true, "34": true, "35": true, "36": true, "41": true}

// gs1Elements reads the element strings of a GS1-128 symbol,
// "(01)04912345678904(10)ABC", as the bytes of the symbol: FNC1 first, and
// after a field of variable length that another follows.
func gs1Elements(text string) ([]byte, error) {
	out := []byte{fnc1}
	if !strings.HasPrefix(text, "(") {
		// as it is, with the group separator for FNC1
		for _, c := range []byte(text) {
			if c == 0x1d {
				c = fnc1
			}
			out = append(out, c)
		}
		return out, nil
	}
	for rest := text; rest != ""; {
		end := strings.IndexByte(rest, ')')
		if rest[0] != '(' || end < 3 || end > 5 || !digits(rest[1:end]) {
			return nil, fmt.Errorf("barcode: gs1-128 takes application identifiers in parentheses, as in (01)04912345678904(10)ABC, not %q", text)
		}
		ai := rest[1:end]
		rest = rest[end+1:]
		next := strings.IndexByte(rest, '(')
		if next < 0 {
			next = len(rest)
		}
		if next == 0 {
			return nil, fmt.Errorf("barcode: the application identifier (%s) of %q has no data", ai, text)
		}
		out = append(out, ai...)
		out = append(out, rest[:next]...)
		rest = rest[next:]
		if rest != "" && !gs1Fixed[ai[:2]] {
			out = append(out, fnc1)
		}
	}
	return out, nil
}

func encodeCode128(text string, gs1 bool, quiet int) (*Symbol, error) {
	for _, c := range []byte(text) {
		if c > 127 {
			return nil, fmt.Errorf("barcode: code128 has the characters of ASCII alone, and %q has others", text)
		}
	}
	data := []byte(text)
	if gs1 {
		var err error
		if data, err = gs1Elements(text); err != nil {
			return nil, err
		}
	}
	// the digits from a place on
	run := func(i int) int {
		n := 0
		for i+n < len(data) && data[i+n] >= '0' && data[i+n] <= '9' {
			n++
		}
		return n
	}
	// the code set for what is not pairs of digits: A for control
	// characters, B for small letters, the one whose character comes first
	ab := func(i int) int {
		for _, c := range data[i:] {
			if c < 32 {
				return c128ToA
			}
			if c >= 96 && c != fnc1 {
				return c128ToB
			}
		}
		return c128ToB
	}
	var vals []int
	set := 0
	to := func(s int) {
		if set == s {
			return
		}
		if set == 0 {
			vals = append(vals, map[int]int{c128ToA: c128StartA, c128ToB: c128StartB, c128ToC: c128StartC}[s])
		} else {
			vals = append(vals, s)
		}
		set = s
	}
	for i := 0; i < len(data); {
		c := data[i]
		n := run(i)
		// pairs of digits: four or more at the start or the end of the
		// data, six or more inside it, or all of it
		if set != c128ToC && (n >= 6 || n >= 4 && (i == 0 || i+n == len(data)) || n == len(data) && n%2 == 0) {
			if n%2 == 1 && i+n == len(data) && set != 0 {
				// an odd run at the end: its first digit in the set it is in
				vals = append(vals, int(c)-32)
				i++
			}
			to(c128ToC)
			continue
		}
		switch {
		case set == c128ToC && n >= 2:
			vals = append(vals, int(c-'0')*10+int(data[i+1]-'0'))
			i += 2
			continue
		case c == fnc1:
			if set == 0 {
				// (pairs of digits after it: the symbol starts with them)
				if k := run(i + 1); k >= 4 || k >= 2 && i+1+k == len(data) {
					to(c128ToC)
				} else {
					to(ab(i))
				}
			}
			vals = append(vals, c128FNC1)
		case c < 32:
			to(c128ToA)
			vals = append(vals, int(c)+64)
		case c >= 96:
			to(c128ToB)
			vals = append(vals, int(c)-32)
		default:
			if set == 0 || set == c128ToC {
				to(ab(i))
			}
			vals = append(vals, int(c)-32)
		}
		i++
	}
	sum := vals[0]
	for i, v := range vals[1:] {
		sum += (i + 1) * v
	}
	vals = append(vals, sum%103, c128Stop)
	var widths []int
	for _, v := range vals {
		for _, e := range []byte(code128Chars[v]) {
			widths = append(widths, int(e-'0'))
		}
	}
	if quiet < 0 {
		quiet = 10
	}
	s := linear(widths, quiet, quiet)
	s.Text = text
	return s, nil
}

// --- ITF (interleaved 2 of 5) ---

// The elements of the digits of ITF: five, two of them wide.
var itfDigits = [10]string{"00110", "10001", "01001", "11000", "00101", "10100", "01100", "00011", "10010", "01010"}

func encodeITF(text string, itf14 bool, quiet int) (*Symbol, error) {
	t := strings.TrimSpace(text)
	if itf14 {
		var err error
		if t, err = gtin(ITF14, t, 14); err != nil {
			return nil, err
		}
	} else if !digits(t) {
		return nil, fmt.Errorf("barcode: itf takes digits, not %q", text)
	}
	if len(t)%2 == 1 {
		t = "0" + t // the digits go in pairs
	}
	// a narrow element is two modules and a wide one five
	w := func(e byte) int { return 2 + 3*int(e-'0') }
	widths := []int{2, 2, 2, 2} // the start: narrow bars and spaces
	for i := 0; i < len(t); i += 2 {
		bars, spaces := itfDigits[t[i]-'0'], itfDigits[t[i+1]-'0']
		for k := 0; k < 5; k++ {
			widths = append(widths, w(bars[k]), w(spaces[k]))
		}
	}
	widths = append(widths, 5, 2, 2) // the stop: a wide bar, a space, a narrow bar
	if quiet < 0 {
		quiet = 10
	}
	s := linear(widths, 2*quiet, 2*quiet)
	s.Text = t
	return s, nil
}

// --- NW-7 (Codabar) ---

// The elements of the characters of NW-7: four bars and the three spaces
// between them, a bar first; 1 is a wide element.
var nw7Chars = map[byte]string{
	'0': "0000011", '1': "0000110", '2': "0001001", '3': "1100000", '4': "0010010",
	'5': "1000010", '6': "0100001", '7': "0100100", '8': "0110000", '9': "1001000",
	'-': "0001100", '$': "0011000", ':': "1000101", '/': "1010001", '.': "1010100", '+': "0010101",
	'A': "0011010", 'B': "0101001", 'C': "0001011", 'D': "0001110",
}

func encodeNW7(text string, quiet int) (*Symbol, error) {
	t := strings.ToUpper(strings.TrimSpace(text))
	start := func(c byte) bool { return c >= 'A' && c <= 'D' }
	// a text without start and stop characters is between A and A
	if len(t) < 2 || !start(t[0]) || !start(t[len(t)-1]) {
		t = "A" + t + "A"
	}
	var widths []int
	for i, c := range []byte(t) {
		p, ok := nw7Chars[c]
		if !ok || start(c) != (i == 0 || i == len(t)-1) {
			return nil, fmt.Errorf("barcode: nw7 takes digits and - $ : / . + between the start and stop characters A to D, not %q", text)
		}
		if i > 0 {
			widths = append(widths, 1)
		}
		for _, e := range []byte(p) {
			widths = append(widths, 1+2*int(e-'0'))
		}
	}
	if quiet < 0 {
		quiet = 10
	}
	s := linear(widths, quiet, quiet)
	s.Text = t
	return s, nil
}
