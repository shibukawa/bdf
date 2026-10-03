package mathlayout

// The control sequences of LaTeX's math mode (with amsmath and amssymb)
// that stand for one character.

// texOrd are symbols that are ordinary atoms (letters and signs).
var texOrd = map[string]string{
	"alpha": "α", "beta": "β", "gamma": "γ", "delta": "δ", "epsilon": "ϵ", "varepsilon": "ε", "zeta": "ζ",
	"eta": "η", "theta": "θ", "vartheta": "ϑ", "iota": "ι", "kappa": "κ", "varkappa": "ϰ", "lambda": "λ", "mu": "μ",
	"nu": "ν", "xi": "ξ", "omicron": "ο", "pi": "π", "varpi": "ϖ", "rho": "ρ", "varrho": "ϱ", "sigma": "σ",
	"varsigma": "ς", "tau": "τ", "upsilon": "υ", "phi": "ϕ", "varphi": "φ", "chi": "χ", "psi": "ψ", "omega": "ω",
	"digamma": "ϝ",
	"Gamma":   "Γ", "Delta": "Δ", "Theta": "Θ", "Lambda": "Λ", "Xi": "Ξ", "Pi": "Π", "Sigma": "Σ", "Upsilon": "Υ",
	"Phi": "Φ", "Psi": "Ψ", "Omega": "Ω",
	"infty": "∞", "partial": "∂", "nabla": "∇", "emptyset": "∅", "varnothing": "∅", "hbar": "ℏ", "hslash": "ℏ",
	"ell": "ℓ", "wp": "℘", "Re": "ℜ", "Im": "ℑ", "aleph": "ℵ", "beth": "ℶ", "gimel": "ℷ", "daleth": "ℸ",
	"forall": "∀", "exists": "∃", "nexists": "∄", "neg": "¬", "lnot": "¬", "prime": "′", "backprime": "‵",
	"angle": "∠", "measuredangle": "∡", "sphericalangle": "∢", "triangle": "△", "bot": "⊥", "top": "⊤",
	"flat": "♭", "natural": "♮", "sharp": "♯", "clubsuit": "♣", "diamondsuit": "♢", "heartsuit": "♡",
	"spadesuit": "♠", "imath": "ı", "jmath": "ȷ", "surd": "√", "complement": "∁", "backslash": "\\",
	"degree": "°", "checkmark": "✓", "square": "□", "blacksquare": "■", "Box": "□", "Diamond": "◇",
	"mho": "℧", "eth": "ð", "Finv": "Ⅎ", "Game": "⅁", "Bbbk": "𝕜", "circledS": "Ⓢ", "diagup": "╱",
	"diagdown": "╲", "lozenge": "◊", "blacklozenge": "⧫", "bigstar": "★", "blacktriangle": "▲",
	"blacktriangledown": "▼", "triangledown": "▽", "vartriangle": "△", "S": "§", "P": "¶", "dag": "†",
	"ddag": "‡", "pounds": "£", "yen": "¥", "copyright": "©", "circledR": "®",
	"ldots": "…", "dots": "…", "dotsc": "…", "dotso": "…", "cdots": "⋯", "dotsb": "⋯", "dotsm": "⋯",
	"dotsi": "⋯", "vdots": "⋮", "ddots": "⋱", "iddots": "⋰", "cdotp": "·", "ldotp": ".",
	"|": "‖", "Vert": "‖", "vert": "|", "lvert": "|", "rvert": "|", "lVert": "‖", "rVert": "‖",
	"#": "#", "$": "$", "%": "%", "&": "&", "_": "_",
}

// texVarGreek are the italic capital Greek letters (\varGamma …).
var texVarGreek = map[string]string{
	"varGamma": "Γ", "varDelta": "Δ", "varTheta": "Θ", "varLambda": "Λ", "varXi": "Ξ", "varPi": "Π",
	"varSigma": "Σ", "varUpsilon": "Υ", "varPhi": "Φ", "varPsi": "Ψ", "varOmega": "Ω",
}

// texOps are operators, relations and punctuation.
var texOps = map[string]string{
	"pm": "±", "mp": "∓", "times": "×", "div": "÷", "cdot": "⋅", "ast": "∗", "star": "⋆", "circ": "∘",
	"bullet": "∙", "oplus": "⊕", "ominus": "⊖", "otimes": "⊗", "oslash": "⊘", "odot": "⊙", "circledcirc": "⊚",
	"circledast": "⊛", "circleddash": "⊝", "cup": "∪", "cap": "∩", "sqcup": "⊔", "sqcap": "⊓", "uplus": "⊎",
	"wedge": "∧", "land": "∧", "vee": "∨", "lor": "∨", "setminus": "∖", "smallsetminus": "∖", "wr": "≀",
	"diamond": "⋄", "bigtriangleup": "△", "bigtriangledown": "▽", "triangleleft": "◁", "triangleright": "▷",
	"lhd": "⊲", "rhd": "⊳", "unlhd": "⊴", "unrhd": "⊵", "dagger": "†", "ddagger": "‡", "amalg": "⨿",
	"ltimes": "⋉", "rtimes": "⋊", "boxplus": "⊞", "boxminus": "⊟", "boxtimes": "⊠", "boxdot": "⊡",
	"Cup": "⋓", "Cap": "⋒", "curlyvee": "⋎", "curlywedge": "⋏", "intercal": "⊺", "divideontimes": "⋇",
	"dotplus": "∔", "barwedge": "⊼", "veebar": "⊻", "centerdot": "⋅",
	"leq": "≤", "le": "≤", "geq": "≥", "ge": "≥", "neq": "≠", "ne": "≠", "equiv": "≡", "approx": "≈",
	"sim": "∼", "simeq": "≃", "cong": "≅", "propto": "∝", "varpropto": "∝", "in": "∈", "notin": "∉",
	"ni": "∋", "owns": "∋", "subset": "⊂", "supset": "⊃", "subseteq": "⊆", "supseteq": "⊇", "subsetneq": "⊊",
	"supsetneq": "⊋", "nsubseteq": "⊈", "nsupseteq": "⊉", "sqsubset": "⊏", "sqsupset": "⊐", "sqsubseteq": "⊑",
	"sqsupseteq": "⊒", "ll": "≪", "gg": "≫", "lll": "⋘", "ggg": "⋙", "prec": "≺", "succ": "≻",
	"preceq": "⪯", "succeq": "⪰", "parallel": "∥", "nparallel": "∦", "perp": "⟂", "mid": "∣", "nmid": "∤",
	"vdash": "⊢", "dashv": "⊣", "models": "⊨", "vDash": "⊨", "Vdash": "⊩", "doteq": "≐", "asymp": "≍",
	"bowtie": "⋈", "smile": "⌣", "frown": "⌢", "lesssim": "≲", "gtrsim": "≳", "leqslant": "⩽", "geqslant": "⩾",
	"leqq": "≦", "geqq": "≧", "lessgtr": "≶", "gtrless": "≷", "approxeq": "≊", "triangleq": "≜",
	"coloneqq": "≔", "eqqcolon": "≕", "coloneq": "≔", "therefore": "∴", "because": "∵", "nleq": "≰",
	"ngeq": "≱", "nless": "≮", "ngtr": "≯", "ncong": "≇", "nsim": "≁", "lessdot": "⋖",
	"gtrdot": "⋗", "backsim": "∽", "eqsim": "≂", "Doteq": "≑", "fallingdotseq": "≒", "risingdotseq": "≓",
	"circeq": "≗", "bumpeq": "≏", "Bumpeq": "≎", "between": "≬", "pitchfork": "⋔", "trianglelefteq": "⊴",
	"trianglerighteq": "⊵", "vartriangleleft": "⊲", "vartriangleright": "⊳", "Subset": "⋐", "Supset": "⋑",
	"to": "→", "rightarrow": "→", "leftarrow": "←", "gets": "←", "leftrightarrow": "↔", "Rightarrow": "⇒",
	"Leftarrow": "⇐", "Leftrightarrow": "⇔", "implies": "⟹", "impliedby": "⟸", "iff": "⟺", "mapsto": "↦",
	"longmapsto": "⟼", "longrightarrow": "⟶", "longleftarrow": "⟵", "longleftrightarrow": "⟷",
	"Longrightarrow": "⟹", "Longleftarrow": "⟸", "Longleftrightarrow": "⟺", "uparrow": "↑", "downarrow": "↓",
	"updownarrow": "↕", "Uparrow": "⇑", "Downarrow": "⇓", "Updownarrow": "⇕", "nearrow": "↗", "searrow": "↘",
	"swarrow": "↙", "nwarrow": "↖", "hookrightarrow": "↪", "hookleftarrow": "↩", "rightharpoonup": "⇀",
	"rightharpoondown": "⇁", "leftharpoonup": "↼", "leftharpoondown": "↽", "rightleftharpoons": "⇌",
	"leftrightharpoons": "⇋", "leftrightarrows": "⇆", "rightleftarrows": "⇄", "leadsto": "⇝",
	"rightsquigarrow": "⇝", "twoheadrightarrow": "↠", "twoheadleftarrow": "↞", "rightarrowtail": "↣",
	"leftarrowtail": "↢", "curvearrowright": "↷", "curvearrowleft": "↶", "circlearrowright": "↻",
	"circlearrowleft": "↺", "Lsh": "↰", "Rsh": "↱", "upharpoonright": "↾", "downharpoonright": "⇂",
	"nrightarrow": "↛", "nleftarrow": "↚", "nRightarrow": "⇏", "nLeftarrow": "⇍", "nleftrightarrow": "↮",
	"nLeftrightarrow": "⇎",
	"colon":           ":", "vcentcolon": ":", "ratio": "∶",
}

// texLargeOps are the large operators.
var texLargeOps = map[string]string{
	"sum": "∑", "prod": "∏", "coprod": "∐", "int": "∫", "iint": "∬", "iiint": "∭", "iiiint": "⨌", "oint": "∮",
	"oiint": "∯", "oiiint": "∰", "intop": "∫", "smallint": "∫", "bigcup": "⋃", "bigcap": "⋂", "bigvee": "⋁",
	"bigwedge": "⋀", "bigoplus": "⨁", "bigotimes": "⨂", "bigodot": "⨀", "biguplus": "⨄", "bigsqcup": "⨆",
}

// texDelims are the delimiters \left, \right and \big take.
var texDelims = map[string]string{
	"langle": "⟨", "rangle": "⟩", "lbrace": "{", "rbrace": "}", "{": "{", "}": "}", "lbrack": "[",
	"rbrack": "]", "lfloor": "⌊", "rfloor": "⌋", "lceil": "⌈", "rceil": "⌉", "vert": "|", "Vert": "‖",
	"lvert": "|", "rvert": "|", "lVert": "‖", "rVert": "‖", "|": "‖", "backslash": "\\", "lgroup": "⟮",
	"rgroup": "⟯", "llbracket": "⟦", "rrbracket": "⟧", "ulcorner": "⌜", "urcorner": "⌝", "llcorner": "⌞",
	"lrcorner": "⌟", "uparrow": "↑", "downarrow": "↓", "updownarrow": "↕", "Uparrow": "⇑",
	"Downarrow": "⇓", "Updownarrow": "⇕", "lmoustache": "⎰", "rmoustache": "⎱",
}

// texAccents are the accents over a base; the wide ones grow to it.
var texAccents = map[string]struct {
	char    string
	stretch bool
}{
	"hat": {"^", false}, "check": {"ˇ", false}, "tilde": {"~", false}, "acute": {"´", false},
	"grave": {"`", false}, "dot": {"˙", false}, "ddot": {"¨", false}, "dddot": {"⃛", false},
	"breve": {"˘", false}, "bar": {"¯", false}, "vec": {"→", false}, "mathring": {"˚", false},
	"widehat": {"^", true}, "widetilde": {"~", true}, "widecheck": {"ˇ", true},
	"overrightarrow": {"→", true}, "overleftarrow": {"←", true}, "overleftrightarrow": {"↔", true},
	"Overrightarrow": {"⇒", true}, "overgroup": {"⏠", true}, "overarc": {"⏜", true},
}

// texUnderAccents go under their base.
var texUnderAccents = map[string]string{
	"underrightarrow": "→", "underleftarrow": "←", "underleftrightarrow": "↔", "undergroup": "⏡",
	"utilde": "~",
}

// texNegated are the negations of relations that \not makes.
var texNegated = map[string]string{
	"=": "≠", "<": "≮", ">": "≯", "≤": "≰", "≥": "≱", "∈": "∉", "∋": "∌", "⊂": "⊄", "⊃": "⊅",
	"⊆": "⊈", "⊇": "⊉", "≡": "≢", "∼": "≁", "≃": "≄", "≅": "≇", "≈": "≉", "∣": "∤", "∥": "∦",
	"≺": "⊀", "≻": "⊁", "⊢": "⊬", "⊨": "⊭", "→": "↛", "←": "↚", "↔": "↮", "⇒": "⇏", "⇐": "⇍", "⇔": "⇎",
	"∃": "∄",
}

// texSpaces are the spacing commands, in em.
var texSpaces = map[string]float64{
	",": 3.0 / 18, ":": 4.0 / 18, ">": 4.0 / 18, ";": 5.0 / 18, "!": -3.0 / 18, " ": 0.25, "quad": 1,
	"qquad": 2, "enspace": 0.5, "thinspace": 3.0 / 18, "medspace": 4.0 / 18, "thickspace": 5.0 / 18,
	"negthinspace": -3.0 / 18, "negmedspace": -4.0 / 18, "negthickspace": -5.0 / 18, "enskip": 0.5,
	"nobreakspace": 0.25, "space": 0.25,
}

// texVariants are the font commands and the variants they set.
var texVariants = map[string]Variant{
	"mathrm": VarNormal, "mathup": VarNormal, "mathit": VarItalic, "mathbf": VarBold, "mathbfit": VarBoldItalic,
	"boldsymbol": VarBoldItalic, "bm": VarBoldItalic, "mathbb": VarDoubleStruck, "Bbb": VarDoubleStruck,
	"mathcal": VarScript, "mathscr": VarScript, "mathfrak": VarFraktur, "frak": VarFraktur,
	"mathsf": VarSansSerif, "mathsfit": VarSansSerifItalic, "mathbfsf": VarSansSerifBold, "mathtt": VarMonospace,
	"mathnormal": VarAuto,
}

// texOldFonts are the font declarations of plain TeX.
var texOldFonts = map[string]Variant{
	"rm": VarNormal, "it": VarItalic, "bf": VarBold, "cal": VarScript, "sf": VarSansSerif, "tt": VarMonospace,
	"mit": VarAuto,
}

// texBig are the sizes of \big and its larger forms, in em.
var texBig = map[string]float64{"big": 1.2, "Big": 1.8, "bigg": 2.4, "Bigg": 3.0}
