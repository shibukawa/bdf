package docx

import (
	"strconv"
	"strings"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// numLevel is one level of a list definition (w:lvl).
type numLevel struct {
	start   int
	format  string // numFmt: decimal, bullet, upperRoman …
	text    string // lvlText: "%1.", "•" …
	jc      string // left, center, right
	suffix  string // tab, space, nothing
	restart int    // lvlRestart: restart after this level (0-based), -1: after any higher level
	legal   bool   // isLgl: numbers of higher levels in decimal
	pPr     *ooxml.Node
	rPr     *ooxml.Node
}

type abstractNum struct {
	id     string
	levels [9]*numLevel
	link   string // numStyleLink: the numbering style whose list this is
}

type numInst struct {
	abs      *abstractNum
	levels   [9]*numLevel // lvlOverride/lvl
	starts   [9]int       // lvlOverride/startOverride (0: none)
	hasStart [9]bool
	started  bool
}

// numbering is the document's list definitions and the counters of the
// lists while the document is read in order.
type numbering struct {
	abstract map[string]*abstractNum
	nums     map[string]*numInst
	counters map[*abstractNum]*[9]int
	used     map[*abstractNum]*[9]bool
}

func (c *converter) loadNumbering() {
	c.num = &numbering{abstract: map[string]*abstractNum{}, nums: map[string]*numInst{},
		counters: map[*abstractNum]*[9]int{}, used: map[*abstractNum]*[9]bool{}}
	r, ok := c.pkg.RelOfType(c.main, "/numbering")
	if !ok {
		return
	}
	root, err := c.pkg.XML(r.Target)
	if err != nil {
		c.warnf("numbering: %v", err)
		return
	}
	parseLvl := func(l *ooxml.Node) *numLevel {
		lv := &numLevel{start: int(l.Child("start").AttrInt("val", 0)), format: "decimal", text: "", jc: "left", suffix: "tab",
			restart: -1, pPr: l.Child("pPr"), rPr: l.Child("rPr")}
		if f := l.Child("numFmt"); f != nil {
			lv.format = f.AttrStr("val", "decimal")
			if lv.format == "custom" {
				lv.format = "decimal"
			}
		}
		if t := l.Child("lvlText"); t != nil {
			lv.text = t.AttrStr("val", "")
		}
		if j := l.Child("lvlJc"); j != nil {
			lv.jc = jcValue(val(j))
		}
		if s := l.Child("suff"); s != nil {
			lv.suffix = s.AttrStr("val", "tab")
		}
		if r := l.Child("lvlRestart"); r != nil {
			lv.restart = int(r.AttrInt("val", 0)) - 1
		}
		lv.legal = onOff(l.Child("isLgl"))
		return lv
	}
	for _, a := range root.Children("abstractNum") {
		an := &abstractNum{id: a.AttrStr("abstractNumId", ""), link: val(a.Child("numStyleLink"))}
		for _, l := range a.Children("lvl") {
			if i := int(l.AttrInt("ilvl", 0)); i >= 0 && i < 9 {
				an.levels[i] = parseLvl(l)
			}
		}
		c.num.abstract[an.id] = an
	}
	for _, n := range root.Children("num") {
		ni := &numInst{abs: c.num.abstract[val(n.Child("abstractNumId"))]}
		for _, o := range n.Children("lvlOverride") {
			i := int(o.AttrInt("ilvl", 0))
			if i < 0 || i >= 9 {
				continue
			}
			if s := o.Child("startOverride"); s != nil {
				ni.starts[i], ni.hasStart[i] = int(s.AttrInt("val", 0)), true
			}
			if l := o.Child("lvl"); l != nil {
				ni.levels[i] = parseLvl(l)
			}
		}
		c.num.nums[n.AttrStr("numId", "")] = ni
	}
	// numStyleLink: the list is the one of the numbering style's numbering
	for _, an := range c.num.abstract {
		if an.link == "" || an.levels[0] != nil {
			continue
		}
		if st := c.st.byID[an.link]; st != nil {
			if id := val(st.pPr.Path("numPr", "numId")); id != "" {
				if ni := c.num.nums[id]; ni != nil && ni.abs != nil && ni.abs != an {
					an.levels = ni.abs.levels
				}
			}
		}
	}
}

// level returns the definition of a list level (nil for no list).
func (n *numbering) level(numID string, ilvl int) *numLevel {
	if n == nil || numID == "" || numID == "0" || ilvl < 0 || ilvl >= 9 {
		return nil
	}
	ni := n.nums[numID]
	if ni == nil || ni.abs == nil {
		return nil
	}
	if l := ni.levels[ilvl]; l != nil {
		return l
	}
	return ni.abs.levels[ilvl]
}

// next counts a paragraph of a list and returns the text of its number.
func (n *numbering) next(numID string, ilvl int) (string, *numLevel) {
	lvl := n.level(numID, ilvl)
	if lvl == nil {
		return "", nil
	}
	ni := n.nums[numID]
	cnt := n.counters[ni.abs]
	used := n.used[ni.abs]
	if cnt == nil {
		cnt, used = &[9]int{}, &[9]bool{}
		n.counters[ni.abs], n.used[ni.abs] = cnt, used
	}
	if !ni.started {
		// the overrides of a list instance restart its levels once
		ni.started = true
		for i := range 9 {
			if ni.hasStart[i] {
				used[i] = false
			}
		}
	}
	levelOf := func(i int) *numLevel {
		if l := n.level(numID, i); l != nil {
			return l
		}
		return &numLevel{start: 1, format: "decimal"}
	}
	if !used[ilvl] {
		start := lvl.start
		if ni.hasStart[ilvl] {
			start = ni.starts[ilvl]
		}
		cnt[ilvl] = start
		used[ilvl] = true
	} else {
		cnt[ilvl]++
	}
	// deeper levels restart
	for i := ilvl + 1; i < 9; i++ {
		if l := levelOf(i); l.restart == -1 || l.restart >= ilvl {
			used[i] = false
		}
	}
	text := lvl.text
	if lvl.format == "bullet" {
		return text, lvl
	}
	if lvl.format == "none" {
		// only the literal text around the numbers
		for i := range 9 {
			text = strings.ReplaceAll(text, "%"+strconv.Itoa(i+1), "")
		}
		return text, lvl
	}
	for i := 8; i >= 0; i-- {
		ph := "%" + strconv.Itoa(i+1)
		if !strings.Contains(text, ph) {
			continue
		}
		v := cnt[i]
		if !used[i] {
			v = levelOf(i).start
		}
		f := levelOf(i).format
		if lvl.legal && i != ilvl {
			f = "decimal"
		}
		text = strings.ReplaceAll(text, ph, formatNumber(v, f))
	}
	return text, lvl
}

// formatNumber formats a list or page number in an ST_NumberFormat.
func formatNumber(n int, format string) string {
	switch format {
	case "upperRoman":
		return roman(n)
	case "lowerRoman":
		return strings.ToLower(roman(n))
	case "upperLetter":
		return letters(n)
	case "lowerLetter":
		return strings.ToLower(letters(n))
	case "decimalZero":
		if n >= 0 && n < 10 {
			return "0" + strconv.Itoa(n)
		}
	case "decimalFullWidth", "decimalFullWidth2":
		return fullWidth(strconv.Itoa(n))
	case "decimalEnclosedCircle", "decimalEnclosedCircleChinese":
		if n >= 1 && n <= 20 {
			return string(rune(0x2460 + n - 1))
		}
		if n >= 21 && n <= 35 {
			return string(rune(0x3251 + n - 21))
		}
		if n >= 36 && n <= 50 {
			return string(rune(0x32b1 + n - 36))
		}
	case "decimalEnclosedParen":
		if n >= 1 && n <= 20 {
			return string(rune(0x2474 + n - 1))
		}
		return "(" + strconv.Itoa(n) + ")"
	case "decimalEnclosedFullstop":
		if n >= 1 && n <= 20 {
			return string(rune(0x2488 + n - 1))
		}
	case "ideographDigital", "japaneseDigitalTenThousand":
		return kanjiDigits(n)
	case "japaneseCounting", "chineseCounting", "chineseCountingThousand", "japaneseLegal", "taiwaneseCounting", "taiwaneseCountingThousand":
		return kanjiNumber(n)
	case "aiueo":
		return kanaSeq(n, aiueo, false)
	case "aiueoFullWidth":
		return kanaSeq(n, aiueo, true)
	case "iroha":
		return kanaSeq(n, iroha, false)
	case "irohaFullWidth":
		return kanaSeq(n, iroha, true)
	case "ideographTraditional":
		if n >= 1 && n <= 10 {
			return string([]rune("甲乙丙丁戊己庚辛壬癸")[n-1])
		}
	case "ideographZodiac":
		if n >= 1 && n <= 12 {
			return string([]rune("子丑寅卯辰巳午未申酉戌亥")[n-1])
		}
	case "ordinal":
		return strconv.Itoa(n) + ordinalSuffix(n)
	case "numberInDash":
		return "- " + strconv.Itoa(n) + " -"
	case "none":
		return ""
	}
	return strconv.Itoa(n)
}

func ordinalSuffix(n int) string {
	if n%100 >= 11 && n%100 <= 13 {
		return "th"
	}
	switch n % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	}
	return "th"
}

func roman(n int) string {
	if n < 1 || n > 3999 {
		return strconv.Itoa(n)
	}
	vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	syms := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	var b strings.Builder
	for i, v := range vals {
		for n >= v {
			b.WriteString(syms[i])
			n -= v
		}
	}
	return b.String()
}

// letters numbers A…Z, AA…ZZ as Word does.
func letters(n int) string {
	if n < 1 {
		return strconv.Itoa(n)
	}
	return strings.Repeat(string(rune('A'+(n-1)%26)), (n-1)/26+1)
}

func fullWidth(s string) string {
	var out []rune
	for _, r := range s {
		if r >= '0' && r <= '9' {
			r += 0xff10 - '0'
		}
		out = append(out, r)
	}
	return string(out)
}

func kanjiDigits(n int) string {
	d := []rune("〇一二三四五六七八九")
	var out []rune
	for _, r := range strconv.Itoa(n) {
		if r >= '0' && r <= '9' {
			out = append(out, d[r-'0'])
		}
	}
	return string(out)
}

// kanjiNumber writes n with 十, 百, 千 and 万.
func kanjiNumber(n int) string {
	if n <= 0 {
		return kanjiDigits(n)
	}
	d := []rune("〇一二三四五六七八九")
	group := func(n int) string {
		var b strings.Builder
		for _, u := range []struct {
			v int
			s string
		}{{1000, "千"}, {100, "百"}, {10, "十"}} {
			if q := n / u.v; q > 0 {
				if q > 1 {
					b.WriteRune(d[q])
				}
				b.WriteString(u.s)
				n %= u.v
			}
		}
		if n > 0 {
			b.WriteRune(d[n])
		}
		return b.String()
	}
	if n >= 10000 {
		return group(n/10000) + "万" + group(n%10000)
	}
	return group(n)
}

const (
	aiueo = "アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン"
	iroha = "イロハニホヘトチリヌルヲワカヨタレソツネナラムウヰノオクヤマケフコエテアサキユメミシヱヒモセス"
)

func kanaSeq(n int, seq string, full bool) string {
	rs := []rune(seq)
	if n < 1 {
		return strconv.Itoa(n)
	}
	r := rs[(n-1)%len(rs)]
	if !full {
		if h, ok := halfWidthKana[r]; ok {
			return h
		}
	}
	return string(r)
}

// halfWidthKana maps the katakana of the list sequences to their half-width
// forms (Word's aiueo and iroha formats).
var halfWidthKana = func() map[rune]string {
	full := []rune("アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲンヰヱ")
	half := []string{"ｱ", "ｲ", "ｳ", "ｴ", "ｵ", "ｶ", "ｷ", "ｸ", "ｹ", "ｺ", "ｻ", "ｼ", "ｽ", "ｾ", "ｿ", "ﾀ", "ﾁ", "ﾂ", "ﾃ", "ﾄ",
		"ﾅ", "ﾆ", "ﾇ", "ﾈ", "ﾉ", "ﾊ", "ﾋ", "ﾌ", "ﾍ", "ﾎ", "ﾏ", "ﾐ", "ﾑ", "ﾒ", "ﾓ", "ﾔ", "ﾕ", "ﾖ", "ﾗ", "ﾘ", "ﾙ", "ﾚ", "ﾛ", "ﾜ", "ｦ", "ﾝ", "ｲ", "ｴ"}
	m := map[rune]string{}
	for i, r := range full {
		m[r] = half[i]
	}
	return m
}()
