package otlayout

import (
	"strconv"
	"strings"
)

// FeatureName returns the name of a registered feature tag ("Standard
// Ligatures" for "liga"), or "" for tags the OpenType feature registry
// does not list.
func FeatureName(tag string) string {
	if n, ok := featureNames[tag]; ok {
		return n
	}
	if len(tag) == 4 && digits(tag[2:]) {
		n, _ := strconv.Atoi(tag[2:])
		switch {
		case tag[:2] == "ss" && n >= 1 && n <= 20:
			return "Stylistic Set " + strconv.Itoa(n)
		case tag[:2] == "cv" && n >= 1 && n <= 99:
			return "Character Variant " + strconv.Itoa(n)
		}
	}
	return ""
}

// Feature defaults: when browsers apply a feature to text without being
// asked to (with HarfBuzz, as Chrome and Firefox do).
const (
	// Optional features apply only when the text asks for them
	// (font-feature-settings, font-variant).
	Optional = iota
	// Default features apply to all horizontal text.
	Default
	// ScriptDefault features apply where the shaper of a script calls for
	// them (the joining forms of Arabic, the conjuncts of Indic scripts).
	ScriptDefault
	// Vertical features apply to vertical text.
	Vertical
)

// FeatureDefault tells when browsers apply a feature.
func FeatureDefault(tag string) int {
	switch tag {
	case "abvm", "blwm", "calt", "ccmp", "clig", "curs", "dist", "kern", "liga", "locl", "mark", "mkmk", "rclt", "rlig", "rvrn":
		return Default
	case "isol", "init", "medi", "med2", "fina", "fin2", "fin3", "mset", "stch",
		"nukt", "akhn", "rphf", "rkrf", "pref", "blwf", "abvf", "half", "pstf", "vatu", "cjct", "cfar",
		"pres", "abvs", "blws", "psts", "haln", "ljmo", "vjmo", "tjmo", "ltra", "ltrm", "rtla", "rtlm":
		return ScriptDefault
	case "vert", "vrt2", "vrtr", "vkrn":
		return Vertical
	}
	return Optional
}

var featureNames = map[string]string{
	"aalt": "Access All Alternates",
	"abvf": "Above-base Forms",
	"abvm": "Above-base Mark Positioning",
	"abvs": "Above-base Substitutions",
	"afrc": "Alternative Fractions",
	"akhn": "Akhand",
	"blwf": "Below-base Forms",
	"blwm": "Below-base Mark Positioning",
	"blws": "Below-base Substitutions",
	"c2pc": "Petite Capitals From Capitals",
	"c2sc": "Small Capitals From Capitals",
	"calt": "Contextual Alternates",
	"case": "Case-Sensitive Forms",
	"ccmp": "Glyph Composition / Decomposition",
	"cfar": "Conjunct Form After Ro",
	"chws": "Contextual Half-width Spacing",
	"cjct": "Conjunct Forms",
	"clig": "Contextual Ligatures",
	"cpct": "Centered CJK Punctuation",
	"cpsp": "Capital Spacing",
	"cswh": "Contextual Swash",
	"curs": "Cursive Positioning",
	"dist": "Distances",
	"dlig": "Discretionary Ligatures",
	"dnom": "Denominators",
	"dtls": "Dotless Forms",
	"expt": "Expert Forms",
	"falt": "Final Glyph on Line Alternates",
	"fin2": "Terminal Forms #2",
	"fin3": "Terminal Forms #3",
	"fina": "Terminal Forms",
	"flac": "Flattened Accent Forms",
	"frac": "Fractions",
	"fwid": "Full Widths",
	"half": "Half Forms",
	"haln": "Halant Forms",
	"halt": "Alternate Half Widths",
	"hist": "Historical Forms",
	"hkna": "Horizontal Kana Alternates",
	"hlig": "Historical Ligatures",
	"hngl": "Hangul",
	"hojo": "Hojo Kanji Forms",
	"hwid": "Half Widths",
	"init": "Initial Forms",
	"isol": "Isolated Forms",
	"ital": "Italics",
	"jalt": "Justification Alternates",
	"jp04": "JIS2004 Forms",
	"jp78": "JIS78 Forms",
	"jp83": "JIS83 Forms",
	"jp90": "JIS90 Forms",
	"kern": "Kerning",
	"lfbd": "Left Bounds",
	"liga": "Standard Ligatures",
	"ljmo": "Leading Jamo Forms",
	"lnum": "Lining Figures",
	"locl": "Localized Forms",
	"ltra": "Left-to-right Alternates",
	"ltrm": "Left-to-right Mirrored Forms",
	"mark": "Mark Positioning",
	"med2": "Medial Forms #2",
	"medi": "Medial Forms",
	"mgrk": "Mathematical Greek",
	"mkmk": "Mark to Mark Positioning",
	"mset": "Mark Positioning via Substitution",
	"nalt": "Alternate Annotation Forms",
	"nlck": "NLC Kanji Forms",
	"nukt": "Nukta Forms",
	"numr": "Numerators",
	"onum": "Oldstyle Figures",
	"opbd": "Optical Bounds",
	"ordn": "Ordinals",
	"ornm": "Ornaments",
	"palt": "Proportional Alternate Widths",
	"pcap": "Petite Capitals",
	"pkna": "Proportional Kana",
	"pnum": "Proportional Figures",
	"pref": "Pre-base Forms",
	"pres": "Pre-base Substitutions",
	"pstf": "Post-base Forms",
	"psts": "Post-base Substitutions",
	"pwid": "Proportional Widths",
	"qwid": "Quarter Widths",
	"rand": "Randomize",
	"rclt": "Required Contextual Alternates",
	"rkrf": "Rakar Forms",
	"rlig": "Required Ligatures",
	"rphf": "Reph Form",
	"rtbd": "Right Bounds",
	"rtla": "Right-to-left Alternates",
	"rtlm": "Right-to-left Mirrored Forms",
	"ruby": "Ruby Notation Forms",
	"rvrn": "Required Variation Alternates",
	"salt": "Stylistic Alternates",
	"sinf": "Scientific Inferiors",
	"size": "Optical Size",
	"smcp": "Small Capitals",
	"smpl": "Simplified Forms",
	"ssty": "Math Script-style Alternates",
	"stch": "Stretching Glyph Decomposition",
	"subs": "Subscript",
	"sups": "Superscript",
	"swsh": "Swash",
	"titl": "Titling",
	"tjmo": "Trailing Jamo Forms",
	"tnam": "Traditional Name Forms",
	"tnum": "Tabular Figures",
	"trad": "Traditional Forms",
	"twid": "Third Widths",
	"unic": "Unicase",
	"valt": "Alternate Vertical Metrics",
	"vatu": "Vattu Variants",
	"vchw": "Vertical Contextual Half-width Spacing",
	"vert": "Vertical Alternates",
	"vhal": "Alternate Vertical Half Metrics",
	"vjmo": "Vowel Jamo Forms",
	"vkna": "Vertical Kana Alternates",
	"vkrn": "Vertical Kerning",
	"vpal": "Proportional Alternate Vertical Metrics",
	"vrt2": "Vertical Alternates and Rotation",
	"vrtr": "Vertical Alternates for Rotation",
	"zero": "Slashed Zero",
}

// ScriptName returns the name of a script tag ("Latin" for "latn"), or ""
// for tags it does not know.
func ScriptName(tag string) string { return scriptNames[tag] }

var scriptNames = map[string]string{
	"DFLT": "Default",
	"adlm": "Adlam",
	"arab": "Arabic",
	"armn": "Armenian",
	"bali": "Balinese",
	"beng": "Bengali",
	"bng2": "Bengali v2",
	"bopo": "Bopomofo",
	"brai": "Braille",
	"byzm": "Byzantine Music",
	"cans": "Canadian Syllabics",
	"cher": "Cherokee",
	"copt": "Coptic",
	"cyrl": "Cyrillic",
	"deva": "Devanagari",
	"dev2": "Devanagari v2",
	"ethi": "Ethiopic",
	"geor": "Georgian",
	"gjr2": "Gujarati v2",
	"glag": "Glagolitic",
	"goth": "Gothic",
	"grek": "Greek",
	"gujr": "Gujarati",
	"gur2": "Gurmukhi v2",
	"guru": "Gurmukhi",
	"hang": "Hangul",
	"hani": "CJK Ideographic",
	"hebr": "Hebrew",
	"java": "Javanese",
	"jamo": "Hangul Jamo",
	"kana": "Hiragana and Katakana",
	"khmr": "Khmer",
	"knd2": "Kannada v2",
	"knda": "Kannada",
	"lao ": "Lao",
	"latn": "Latin",
	"math": "Mathematical Alphanumeric Symbols",
	"mlm2": "Malayalam v2",
	"mlym": "Malayalam",
	"mong": "Mongolian",
	"musc": "Musical Symbols",
	"mym2": "Myanmar v2",
	"mymr": "Myanmar",
	"nko ": "N'Ko",
	"ogam": "Ogham",
	"ory2": "Odia v2",
	"orya": "Odia",
	"runr": "Runic",
	"sinh": "Sinhala",
	"sund": "Sundanese",
	"syrc": "Syriac",
	"taml": "Tamil",
	"tel2": "Telugu v2",
	"telu": "Telugu",
	"tfng": "Tifinagh",
	"tglg": "Tagalog",
	"thaa": "Thaana",
	"thai": "Thai",
	"tibt": "Tibetan",
	"tml2": "Tamil v2",
	"vai ": "Vai",
	"yi  ": "Yi",
}

// LangName returns the name of a language system tag ("Japanese" for
// "JAN"), or "" for tags it does not know.
func LangName(tag string) string { return langNames[strings.TrimRight(tag, " ")] }

var langNames = map[string]string{
	"AFK":  "Afrikaans",
	"ARA":  "Arabic",
	"AZE":  "Azerbaijani",
	"BEL":  "Belarusian",
	"BGR":  "Bulgarian",
	"BOS":  "Bosnian",
	"BSH":  "Bashkir",
	"CAT":  "Catalan",
	"CHU":  "Chuvash",
	"CRT":  "Crimean Tatar",
	"CSY":  "Czech",
	"DAN":  "Danish",
	"DEU":  "German",
	"ELL":  "Greek",
	"ENG":  "English",
	"ESP":  "Spanish",
	"ETI":  "Estonian",
	"FAR":  "Persian",
	"FIN":  "Finnish",
	"FRA":  "French",
	"GAG":  "Gagauz",
	"GUA":  "Guarani",
	"HIN":  "Hindi",
	"HRV":  "Croatian",
	"HUN":  "Hungarian",
	"HYE":  "Armenian",
	"IPPH": "Phonetic transcription (IPA)",
	"IRI":  "Irish",
	"IRT":  "Irish Traditional",
	"ISL":  "Icelandic",
	"ITA":  "Italian",
	"IWR":  "Hebrew",
	"JAN":  "Japanese",
	"KAZ":  "Kazakh",
	"KOR":  "Korean",
	"KRK":  "Karakalpak",
	"KSH":  "Kashmiri",
	"KUR":  "Kurdish",
	"LSM":  "Lule Sami",
	"LTH":  "Lithuanian",
	"LVI":  "Latvian",
	"MAR":  "Marathi",
	"MKD":  "Macedonian",
	"MNG":  "Mongolian",
	"MOL":  "Moldavian",
	"MTS":  "Maltese",
	"NAV":  "Navajo",
	"NEP":  "Nepali",
	"NLD":  "Dutch",
	"NOR":  "Norwegian",
	"NSM":  "Northern Sami",
	"PAS":  "Pashto",
	"PLK":  "Polish",
	"PTG":  "Portuguese",
	"ROM":  "Romanian",
	"RUS":  "Russian",
	"SAN":  "Sanskrit",
	"SKS":  "Skolt Sami",
	"SKY":  "Slovak",
	"SLV":  "Slovenian",
	"SND":  "Sindhi",
	"SQI":  "Albanian",
	"SRB":  "Serbian",
	"SVE":  "Swedish",
	"TAT":  "Tatar",
	"TRK":  "Turkish",
	"UKR":  "Ukrainian",
	"URD":  "Urdu",
	"UYG":  "Uyghur",
	"VIT":  "Vietnamese",
	"YID":  "Yiddish",
	"ZHH":  "Chinese, Hong Kong",
	"ZHS":  "Chinese, Simplified",
	"ZHT":  "Chinese, Traditional",
	"ZHTM": "Chinese, Traditional, Macao",
}
