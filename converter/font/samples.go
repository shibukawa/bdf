package font

import (
	"slices"
	"strings"
	"unicode"
)

// script is a writing system the overview shows sample text of, when the
// font has every character of the text.
type script struct {
	name     string
	code     string   // ISO 15924, as the meta table names design languages
	lang     string   // BCP 47, for the LANG marks
	text     string   // a pangram or a well-known line
	specimen string   // the few characters shown large
	charset  []string // lines of characters shown beside them
	rtl      bool
}

var scripts = []script{
	{name: "Latin", code: "Latn", lang: "en", text: "The quick brown fox jumps over the lazy dog.", specimen: "Aa",
		charset: []string{"ABCDEFGHIJKLM", "NOPQRSTUVWXYZ", "abcdefghijklm", "nopqrstuvwxyz", "0123456789"}},
	{name: "Japanese", code: "Jpan", lang: "ja", text: "色は匂へど散りぬるを 我が世誰ぞ常ならむ いろはにほへと", specimen: "あ永",
		charset: []string{"あいうえおかきくけこ", "アイウエオカキクケコ", "永東京春夏秋冬日本語"}},
	{name: "Chinese (Simplified)", code: "Hans", lang: "zh-Hans", text: "天地玄黄，宇宙洪荒。日月盈昃，辰宿列张。", specimen: "永",
		charset: []string{"天地玄黄宇宙洪荒", "日月盈昃辰宿列张"}},
	{name: "Chinese (Traditional)", code: "Hant", lang: "zh-Hant", text: "天地玄黃，宇宙洪荒。日月盈昃，辰宿列張。", specimen: "永",
		charset: []string{"天地玄黃宇宙洪荒", "日月盈昃辰宿列張"}},
	{name: "Korean", code: "Kore", lang: "ko", text: "다람쥐 헌 쳇바퀴에 타고파", specimen: "가",
		charset: []string{"가나다라마바사", "아자차카타파하"}},
	{name: "Cyrillic", code: "Cyrl", lang: "ru", text: "Съешь же ещё этих мягких французских булок, да выпей чаю.", specimen: "Жж",
		charset: []string{"АБВГДЕЖЗИЙКЛМНОП", "абвгдежзийклмноп"}},
	{name: "Greek", code: "Grek", lang: "el", text: "Τάχιστη αλώπηξ βαφής ψημένη γη, δρασκελίζει υπέρ νωθρού κυνός.", specimen: "Ωω",
		charset: []string{"ΑΒΓΔΕΖΗΘΙΚΛΜ", "αβγδεζηθικλμ"}},
	{name: "Arabic", code: "Arab", lang: "ar", text: "نص حكيم له سر قاطع وذو شأن عظيم مكتوب على ثوب أخضر ومغلف بجلد أزرق", specimen: "ع", rtl: true,
		charset: []string{"ا ب ت ث ج ح خ د ذ ر ز", "س ش ص ض ط ظ ع غ ف ق"}},
	{name: "Hebrew", code: "Hebr", lang: "he", text: "דג סקרן שט בים מאוכזב ולפתע מצא חברה", specimen: "א", rtl: true,
		charset: []string{"אבגדהוזחטיכלמנסעפצקרשת"}},
	{name: "Thai", code: "Thai", lang: "th", text: "เป็นมนุษย์สุดประเสริฐเลิศคุณค่า กว่าบรรดาฝูงสัตว์เดรัจฉาน", specimen: "ก",
		charset: []string{"กขคงจฉชซญดตถทนบปผพฟมยรลวศสหอฮ"}},
	{name: "Devanagari", code: "Deva", lang: "hi", text: "ऋषियों को सताने वाले दुष्ट राक्षसों के राजा रावण का सर्वनाश करने वाले विष्णुवतार भगवान श्रीराम", specimen: "क",
		charset: []string{"अआइईउऊएऐओऔ", "कखगघङचछजझञ"}},
	{name: "Bengali", code: "Beng", lang: "bn", text: "আমার সোনার বাংলা, আমি তোমায় ভালোবাসি", specimen: "ক",
		charset: []string{"অআইঈউঊএঐওঔ", "কখগঘঙচছজঝঞ"}},
	{name: "Tamil", code: "Taml", lang: "ta", text: "யாமறிந்த மொழிகளிலே தமிழ்மொழி போல் இனிதாவது எங்கும் காணோம்", specimen: "க",
		charset: []string{"அஆஇஈஉஊஎஏஐஒஓஔ", "கஙசஞடணதநபமயரலவ"}},
	{name: "Armenian", code: "Armn", lang: "hy", text: "Բել դղյակի ձախ ժամն օֆ ազգությանը ցպահանջ չճշտած վնաս էր եւ փառք։", specimen: "Աա",
		charset: []string{"ԱԲԳԴԵԶԷԸԹԺԻԼԽԾԿ", "աբգդեզէըթժիլխծկ"}},
	{name: "Georgian", code: "Geor", lang: "ka", text: "აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ", specimen: "ა",
		charset: []string{"აბგდევზთიკლმნოპჟ", "რსტუფქღყშჩცძწჭხჯჰ"}},
}

// samplePlan is the text the overview draws in the font.
type samplePlan struct {
	primary   *script  // the first script the font covers (nil: none)
	scripts   []script // the scripts the font covers
	specimen  string
	charset   []string
	waterfall string
	rtl       bool
	lang      string
}

// plan picks the sample text of a face; text, when not "", replaces the
// pangram of the waterfall.
func (fc *face) plan(text string) samplePlan {
	var p samplePlan
	han := false // one sample of Chinese characters: Japanese, or one form of Chinese
	for _, s := range scripts {
		cjk := s.lang == "ja" || s.lang[:2] == "zh"
		if cjk && han || !fc.has(s.text) || !fc.has(s.specimen) {
			continue
		}
		han = han || cjk
		p.scripts = append(p.scripts, s)
	}
	if len(p.scripts) > 0 {
		p.primary = fc.primary(p.scripts)
	}
	if p.primary != nil {
		pr := p.primary
		p.specimen, p.waterfall, p.rtl, p.lang = pr.specimen, pr.text, pr.rtl, pr.lang
		for _, l := range pr.charset {
			if fc.has(l) {
				p.charset = append(p.charset, l)
			}
		}
	} else {
		// a font of symbols, or of a script without sample text: its first
		// characters
		var rs []rune
		for _, r := range fc.runes {
			if unicode.IsGraphic(r) && !unicode.IsSpace(r) && !unicode.Is(unicode.Mn, r) {
				rs = append(rs, r)
			}
			if len(rs) == 48 {
				break
			}
		}
		if len(rs) > 0 {
			p.specimen = string(rs[:min(2, len(rs))])
			p.waterfall = string(rs[:min(24, len(rs))])
			for i := 0; i < len(rs) && len(p.charset) < 4; i += 12 {
				p.charset = append(p.charset, string(rs[i:min(i+12, len(rs))]))
			}
		}
	}
	if text != "" {
		p.waterfall = text
	}
	return p
}

// primary picks the script a font shows large and at every size: the
// first design language of its meta table that it covers, else the first
// script it covers beyond the European alphabets (a font for Devanagari
// has Latin letters too), else Latin.
func (fc *face) primary(covered []script) *script {
	design, _ := metaLangs(fc.f.Tables["meta"])
	for _, tag := range strings.Split(design, ", ") {
		// a ScriptLangTag: a script, or a language with its script
		code := tag[strings.LastIndex(tag, "-")+1:]
		for i := range covered {
			if strings.EqualFold(covered[i].code, code) {
				return &covered[i]
			}
		}
	}
	for i := range covered {
		switch covered[i].code {
		case "Latn", "Grek", "Cyrl", "Armn", "Geor":
		default:
			return &covered[i]
		}
	}
	return &covered[0]
}

// all is every piece of text the plan draws, for the sample font.
func (p *samplePlan) all() string {
	s := p.specimen + p.waterfall
	for _, l := range p.charset {
		s += l
	}
	for _, sc := range p.scripts {
		s += sc.text
	}
	return s
}

// scriptCounts counts the characters of the font by Unicode script, most
// first.
func (fc *face) scriptCounts() []scriptCount {
	counts := map[string]int{}
	last := ""
	for _, r := range fc.runes {
		name := "Unknown"
		// the characters are in order, and runs of them in one script
		if t := unicode.Scripts[last]; t != nil && unicode.Is(t, r) {
			name = last
		} else {
			for n, t := range unicode.Scripts {
				if unicode.Is(t, r) {
					name = n
					break
				}
			}
		}
		last = name
		counts[name]++
	}
	var out []scriptCount
	for n, c := range counts {
		out = append(out, scriptCount{n, c})
	}
	slices.SortFunc(out, func(a, b scriptCount) int {
		if a.n != b.n {
			return b.n - a.n
		}
		if a.name < b.name {
			return -1
		}
		return 1
	})
	return out
}

type scriptCount struct {
	name string
	n    int
}
