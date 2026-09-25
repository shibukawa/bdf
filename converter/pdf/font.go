package pdf

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/sfnt"
)

// fontProgram is an embedded font file.
type fontProgram struct {
	kind string // "ttf", "cff", "otf", "type1"
	data []byte
	sf   *sfnt.Font
	cff  *cffFont
}

// usedGlyph records how a character code is rendered and what text it carries.
type usedGlyph struct {
	gid     int
	unicode string // text for search/selection ("" if unknown)
	char    rune   // code point used in FILL_TEXT (maps to gid in the rebuilt cmap)
}

// pdfFont is a loaded PDF font.
type pdfFont struct {
	key       string
	dict      types.Dict
	subtype   string
	composite bool
	type3     bool
	baseFont  string

	enc   *cmap // composite fonts: code → CID
	toUni *cmap

	firstChar    int
	widths       []float64 // simple fonts, glyph space (/1000)
	missingWidth float64
	defaultWidth float64
	cidWidths    map[uint32]float64
	fontMatrix   matrix // Type3

	baseNames   [256]string
	diffNames   [256]string
	hasDiffs    bool
	hasBaseEnc  bool
	symbolic    bool
	flags       int
	prog        *fontProgram
	cid2gid     []byte // CIDToGIDMap stream (nil = identity)
	charProcs   types.Dict
	t3Resources types.Dict
	t3Glyphs    map[string]*pending

	family  string
	weight  uint16
	style   byte
	ascent  float64
	descent float64

	used      map[uint32]*usedGlyph
	uniOwner  map[rune]int // rune → gid that owns it in the rebuilt cmap
	puaNext   rune
	final     bdf.Font
	finalized bool
	warned    map[string]bool
}

// loadFont loads (and caches) a font resource.
func (c *converter) loadFont(ref types.Object, res types.Dict) *pdfFont {
	key := objKey(ref)
	if key != "" {
		if f, ok := c.fonts[key]; ok {
			return f
		}
	}
	p := c.pdf
	d := p.dict(ref)
	if d == nil {
		return nil
	}
	f := &pdfFont{key: key, dict: d, subtype: p.name(d["Subtype"]), baseFont: p.name(d["BaseFont"]),
		used: map[uint32]*usedGlyph{}, uniOwner: map[rune]int{}, puaNext: 0xE000, warned: map[string]bool{},
		fontMatrix: matrix{0.001, 0, 0, 0.001, 0, 0}, defaultWidth: 1000}
	if key != "" {
		c.fonts[key] = f
	}
	if tu := p.stream(d["ToUnicode"]); tu != nil {
		if data, _, err := p.decodeStream(tu); err == nil {
			f.toUni = parseCMap(data)
		}
	}
	desc := p.dict(d["FontDescriptor"])
	switch f.subtype {
	case "Type0":
		f.composite = true
		c.loadType0(f, d)
		if dfs := p.array(d["DescendantFonts"]); len(dfs) > 0 {
			if cid := p.dict(dfs[0]); cid != nil {
				desc = p.dict(cid["FontDescriptor"])
				c.loadCIDFont(f, cid)
			}
		}
	case "Type3":
		f.type3 = true
		f.fontMatrix = p.matrixOr(d["FontMatrix"], f.fontMatrix)
		f.charProcs = p.dict(d["CharProcs"])
		f.t3Resources = p.dict(d["Resources"])
		if f.t3Resources == nil {
			f.t3Resources = res
		}
		f.t3Glyphs = map[string]*pending{}
		c.loadSimpleWidths(f, d)
		c.loadSimpleEncoding(f, d, nil)
	default: // Type1, MMType1, TrueType
		c.loadSimpleWidths(f, d)
		c.loadDescriptor(f, desc)
		c.loadSimpleEncoding(f, d, desc)
	}
	if f.composite {
		c.loadDescriptor(f, desc)
	}
	c.styleFromName(f)
	return f
}

func (c *converter) loadType0(f *pdfFont, d types.Dict) {
	p := c.pdf
	switch e := p.deref(d["Encoding"]).(type) {
	case types.Name:
		switch e.Value() {
		case "Identity-H":
			f.enc = identityCMap()
		case "Identity-V":
			f.enc = identityCMap()
			f.enc.vertical = true
		default:
			c.warnf("predefined CMap %s is not available; treating as Identity", e.Value())
			f.enc = identityCMap()
		}
	default:
		if sd := p.stream(d["Encoding"]); sd != nil {
			if data, _, err := p.decodeStream(sd); err == nil {
				f.enc = parseCMap(data)
				if um := p.deref(sd.Dict["UseCMap"]); um != nil {
					if n, ok := um.(types.Name); ok && strings.HasPrefix(n.Value(), "Identity") {
						f.enc.identity = true
					}
				}
			}
		}
		if f.enc == nil {
			f.enc = identityCMap()
		}
	}
}

func (c *converter) loadCIDFont(f *pdfFont, cid types.Dict) {
	p := c.pdf
	f.defaultWidth = p.numOr(cid["DW"], 1000)
	f.cidWidths = map[uint32]float64{}
	w := p.array(cid["W"])
	for i := 0; i < len(w); {
		start, ok := p.num(w[i])
		if !ok {
			i++
			continue
		}
		if i+1 < len(w) {
			if arr := p.array(w[i+1]); arr != nil {
				for k, e := range arr {
					f.cidWidths[uint32(start)+uint32(k)] = p.numOr(e, 1000)
				}
				i += 2
				continue
			}
			if i+2 < len(w) {
				end := p.numOr(w[i+1], start)
				width := p.numOr(w[i+2], 1000)
				if end-start < 65536 {
					for cc := start; cc <= end; cc++ {
						f.cidWidths[uint32(cc)] = width
					}
				}
				i += 3
				continue
			}
		}
		i++
	}
	if sd := p.stream(cid["CIDToGIDMap"]); sd != nil {
		if data, _, err := p.decodeStream(sd); err == nil {
			f.cid2gid = data
		}
	}
}

func (c *converter) loadSimpleWidths(f *pdfFont, d types.Dict) {
	p := c.pdf
	f.firstChar = p.intOr(d["FirstChar"], 0)
	f.widths = p.nums(d["Widths"])
	if f.type3 {
		// Type3 widths are in glyph space; convert through FontMatrix to text space ×1000.
		for i, w := range f.widths {
			x, _ := f.fontMatrix.apply(w, 0)
			x0, _ := f.fontMatrix.apply(0, 0)
			f.widths[i] = (x - x0) * 1000
		}
	}
}

func (c *converter) loadDescriptor(f *pdfFont, desc types.Dict) {
	p := c.pdf
	if desc == nil {
		return
	}
	f.flags = p.intOr(desc["Flags"], 0)
	f.symbolic = f.flags&4 != 0 && f.flags&32 == 0
	f.missingWidth = p.numOr(desc["MissingWidth"], 0)
	f.ascent = p.numOr(desc["Ascent"], 0)
	f.descent = p.numOr(desc["Descent"], 0)
	if w := p.numOr(desc["FontWeight"], 0); w > 0 {
		f.weight = uint16(w)
	} else if p.numOr(desc["StemV"], 0) >= 120 {
		f.weight = 700
	}
	if f.flags&(1<<6) != 0 || p.numOr(desc["ItalicAngle"], 0) != 0 {
		f.style = bdf.StyleItalic
	}
	for _, key := range []string{"FontFile2", "FontFile3", "FontFile"} {
		sd := p.stream(desc[key])
		if sd == nil {
			continue
		}
		data, _, err := p.decodeStream(sd)
		if err != nil {
			c.warnf("font %s: %s: %v", f.baseFont, key, err)
			continue
		}
		prog := &fontProgram{data: data}
		switch key {
		case "FontFile2":
			prog.kind = "ttf"
		case "FontFile3":
			prog.kind = "cff"
			if st := p.name(sd.Dict["Subtype"]); st == "OpenType" {
				prog.kind = "otf"
			}
		case "FontFile":
			prog.kind = "type1"
		}
		// Sniff the actual format; PDF producers mislabel these regularly.
		if len(data) >= 4 {
			switch {
			case data[0] == 0 && data[1] == 1 && data[2] == 0 && data[3] == 0, string(data[:4]) == "true", string(data[:4]) == "ttcf":
				prog.kind = "ttf"
			case string(data[:4]) == "OTTO":
				prog.kind = "otf"
			case data[0] == 1 && data[1] == 0 && data[2] <= 4:
				prog.kind = "cff"
			case data[0] == '%' && data[1] == '!', data[0] == 0x80:
				prog.kind = "type1"
			}
		}
		switch prog.kind {
		case "ttf", "otf":
			sf, err := sfnt.Parse(data)
			if err != nil {
				c.warnf("font %s: %v", f.baseFont, err)
				continue
			}
			prog.sf = sf
			if sf.IsCFF {
				if cffData, ok := sf.Tables["CFF "]; ok {
					if cf, err := parseCFF(cffData); err == nil {
						prog.cff = cf
					}
				}
			}
		case "cff":
			cf, err := parseCFF(data)
			if err != nil {
				c.warnf("font %s: %v", f.baseFont, err)
				continue
			}
			prog.cff = cf
		case "type1":
			c.warnf("font %s: Type1 font programs are not embedded; using a system font", f.baseFont)
			continue
		}
		f.prog = prog
		break
	}
}

func (c *converter) loadSimpleEncoding(f *pdfFont, d types.Dict, desc types.Dict) {
	p := c.pdf
	// Base encoding: font's built-in for symbolic/embedded fonts, Standard otherwise.
	builtin := f.symbolic && f.prog != nil
	if !builtin && !f.type3 {
		f.baseNames = standardEncoding
		f.hasBaseEnc = true
	}
	apply := func(name string) {
		switch name {
		case "WinAnsiEncoding":
			f.baseNames = winAnsiEncoding
			f.hasBaseEnc = true
		case "MacRomanEncoding":
			f.baseNames = macRomanEncoding
			f.hasBaseEnc = true
		case "StandardEncoding", "MacExpertEncoding":
			f.baseNames = standardEncoding
			f.hasBaseEnc = true
		}
	}
	switch e := p.deref(d["Encoding"]).(type) {
	case types.Name:
		apply(e.Value())
	case types.Dict:
		if be := p.name(e["BaseEncoding"]); be != "" {
			apply(be)
		} else if !builtin && !f.type3 && !f.hasBaseEnc {
			f.baseNames = standardEncoding
			f.hasBaseEnc = true
		}
		code := 0
		for _, item := range p.array(e["Differences"]) {
			if n, ok := p.num(item); ok {
				code = int(n)
				continue
			}
			if name := p.name(item); name != "" && code >= 0 && code < 256 {
				f.diffNames[code] = name
				f.hasDiffs = true
				code++
			}
		}
	}
	// Non-symbolic fonts without any encoding still use Standard.
	if !f.hasBaseEnc && !builtin && !f.type3 {
		f.baseNames = standardEncoding
		f.hasBaseEnc = true
	}
}

// styleFromName derives family/weight/style for system fallback fonts.
func (c *converter) styleFromName(f *pdfFont) {
	name := f.baseFont
	if i := strings.Index(name, "+"); i == 6 {
		name = name[7:]
	}
	lower := strings.ToLower(name)
	if f.weight == 0 {
		switch {
		case strings.Contains(lower, "black"), strings.Contains(lower, "heavy"), strings.Contains(lower, "extrabold"):
			f.weight = 900
		case strings.Contains(lower, "bold"), strings.Contains(lower, "semibold"), strings.Contains(lower, "demibold"):
			f.weight = 700
		case strings.Contains(lower, "light"):
			f.weight = 300
		default:
			f.weight = 400
		}
	}
	if f.style == bdf.StyleNormal && (strings.Contains(lower, "italic") || strings.Contains(lower, "oblique")) {
		f.style = bdf.StyleItalic
	}
	switch {
	case strings.Contains(lower, "courier"), strings.Contains(lower, "mono"), f.flags&1 != 0:
		f.family = "monospace"
	case strings.Contains(lower, "times"), strings.Contains(lower, "serif") && !strings.Contains(lower, "sans"), strings.Contains(lower, "georgia"), strings.Contains(lower, "book"), strings.Contains(lower, "garamond"), strings.Contains(lower, "mincho"), strings.Contains(lower, "roman"):
		f.family = "serif"
	case f.flags&2 != 0 && !strings.Contains(lower, "arial") && !strings.Contains(lower, "helvetica") && !strings.Contains(lower, "gothic"):
		f.family = "serif"
	default:
		f.family = "sans-serif"
	}
	if strings.Contains(lower, "symbol") || strings.Contains(lower, "dingbat") {
		f.family = "serif"
	}
}

func (f *pdfFont) warnOnce(c *converter, key, format string, args ...any) {
	if f.warned[key] {
		return
	}
	f.warned[key] = true
	c.warnf(format, args...)
}

// decode splits string bytes into codes.
func (f *pdfFont) decode(s []byte) []glyphCode {
	if f.composite && f.enc != nil {
		return f.enc.decode(s)
	}
	out := make([]glyphCode, len(s))
	for i, b := range s {
		out[i] = glyphCode{code: uint32(b), nbytes: 1, cid: uint32(b)}
	}
	return out
}

// width returns the horizontal advance of a code in text space units (per unit font size).
func (f *pdfFont) width(g glyphCode) float64 {
	if f.composite {
		if w, ok := f.cidWidths[g.cid]; ok {
			return w / 1000
		}
		return f.defaultWidth / 1000
	}
	i := int(g.code) - f.firstChar
	if i >= 0 && i < len(f.widths) {
		return f.widths[i] / 1000
	}
	if f.missingWidth > 0 || len(f.widths) > 0 {
		return f.missingWidth / 1000
	}
	// No Widths at all (standard 14 fonts): use the embedded program's advances when present.
	if f.prog != nil && f.prog.sf != nil {
		if gid, ok := f.gid(g); ok && gid < len(f.prog.sf.Advances) {
			return float64(f.prog.sf.Advances[gid]) / float64(f.prog.sf.UnitsPerEm)
		}
	}
	return standardWidth(f, g.code)
}

// glyphName returns the encoding glyph name for a simple-font code.
func (f *pdfFont) glyphName(code uint32) string {
	if code >= 256 {
		return ""
	}
	if n := f.diffNames[code]; n != "" {
		return n
	}
	if f.hasBaseEnc {
		return f.baseNames[code]
	}
	return ""
}

// unicode returns the text carried by a code.
func (f *pdfFont) unicode(g glyphCode) string {
	if f.toUni != nil {
		if s := f.toUni.toUnicode(g.code); s != "" {
			return s
		}
	}
	if !f.composite {
		name := f.glyphName(g.code)
		if name == "" && f.symbolic {
			name = standardEncoding[g.code]
		}
		if r, ok := glyphNameToRune(name); ok {
			return string(r)
		}
		if name == "" && g.code >= 32 && g.code < 127 {
			return string(rune(g.code))
		}
	}
	return ""
}

// gid resolves the glyph index of a code within the embedded program.
func (f *pdfFont) gid(g glyphCode) (int, bool) {
	prog := f.prog
	if prog == nil {
		return 0, false
	}
	if f.composite {
		cid := int(g.cid)
		if f.cid2gid != nil {
			if 2*cid+1 < len(f.cid2gid) {
				return int(f.cid2gid[2*cid])<<8 | int(f.cid2gid[2*cid+1]), true
			}
			return 0, false
		}
		if prog.cff != nil && prog.cff.isCID {
			gid, ok := prog.cff.cidToGID[cid]
			return gid, ok
		}
		return cid, true
	}
	name := f.glyphName(g.code)
	code := int(g.code)
	if prog.cff != nil {
		cf := prog.cff
		if name != "" {
			if gid, ok := cf.nameToGID[name]; ok {
				return gid, true
			}
			if r, ok := glyphNameToRune(name); ok {
				// Try the standard name for the same code point.
				for n2, gid := range cf.nameToGID {
					if r2, ok := glyphNameToRune(n2); ok && r2 == r {
						return gid, true
					}
				}
			}
			if i, ok := glyphNameIndex(name); ok && i < cf.numGlyphs {
				return i, true
			}
		}
		if gid, ok := cf.encoding[code]; ok {
			return gid, true
		}
		if name == "" {
			if std := standardEncoding[code]; std != "" {
				if gid, ok := cf.nameToGID[std]; ok {
					return gid, true
				}
			}
		}
		return 0, false
	}
	sf := prog.sf
	if sf == nil {
		return 0, false
	}
	// TrueType simple font lookup (PDF 32000-1 §9.6.6.4), with pragmatic fallbacks.
	tryUnicode := func() (int, bool) {
		if name == "" || sf.Cmap == nil {
			return 0, false
		}
		if r, ok := glyphNameToRune(name); ok {
			if gid, ok := sf.Cmap[uint32(r)]; ok {
				return int(gid), true
			}
		}
		return 0, false
	}
	trySymbol := func() (int, bool) {
		if sf.CmapSymbol == nil {
			return 0, false
		}
		for _, cc := range []uint32{uint32(code), 0xF000 + uint32(code), 0xF100 + uint32(code), 0xF200 + uint32(code)} {
			if gid, ok := sf.CmapSymbol[cc]; ok {
				return int(gid), true
			}
		}
		return 0, false
	}
	tryMac := func() (int, bool) {
		if sf.CmapMac == nil {
			return 0, false
		}
		if gid, ok := sf.CmapMac[uint32(code)]; ok {
			return int(gid), true
		}
		return 0, false
	}
	tryPost := func() (int, bool) {
		if name == "" || sf.PostNames == nil {
			return 0, false
		}
		if gid, ok := sf.PostNames[name]; ok {
			return int(gid), true
		}
		return 0, false
	}
	tryIndex := func() (int, bool) {
		if i, ok := glyphNameIndex(name); ok && i < sf.NumGlyphs {
			return i, true
		}
		return 0, false
	}
	order := []func() (int, bool){trySymbol, tryMac, tryUnicode, tryPost, tryIndex}
	if !f.symbolic || f.hasDiffs {
		order = []func() (int, bool){tryUnicode, tryPost, trySymbol, tryMac, tryIndex}
	}
	for _, fn := range order {
		if gid, ok := fn(); ok {
			return gid, true
		}
	}
	if sf.Cmap == nil && sf.CmapSymbol == nil && sf.CmapMac == nil {
		// No cmap at all: codes index glyphs directly.
		if code < sf.NumGlyphs {
			return code, true
		}
	}
	if sf.Cmap != nil && name == "" && code >= 32 {
		if gid, ok := sf.Cmap[uint32(code)]; ok {
			return int(gid), true
		}
	}
	return 0, false
}

// use records a code as used and returns the text to draw and to index.
// draw is the string for FILL_TEXT; text is the Unicode for search (may differ).
func (f *pdfFont) use(g glyphCode) (draw string, text string) {
	if u, ok := f.used[g.code]; ok {
		if u.char != 0 {
			return string(u.char), u.unicode
		}
		return u.unicode, u.unicode
	}
	uni := f.unicode(g)
	u := &usedGlyph{unicode: uni}
	f.used[g.code] = u
	if f.prog == nil || f.type3 {
		// System font: draw the Unicode text itself.
		if uni == "" && !f.composite {
			uni = string(rune(g.code))
			u.unicode = ""
			u.char = rune(g.code)
			return uni, ""
		}
		return uni, uni
	}
	gid, ok := f.gid(g)
	if !ok {
		gid = 0
	}
	u.gid = gid
	// Choose the code point that will map to this glyph in the rebuilt cmap.
	runes := []rune(uni)
	if len(runes) == 1 && runes[0] >= 0x20 && !unicode.Is(unicode.Co, runes[0]) {
		r := runes[0]
		if owner, taken := f.uniOwner[r]; !taken || owner == gid {
			f.uniOwner[r] = gid
			u.char = r
			return string(r), uni
		}
	}
	// Conflicting or multi-character mapping: private use code point.
	for {
		r := f.puaNext
		f.puaNext++
		if _, taken := f.uniOwner[r]; !taken {
			f.uniOwner[r] = gid
			u.char = r
			return string(r), uni
		}
	}
}

// finalize builds the bdf.Font (and font part) after all pages were processed.
func (c *converter) finalizeFont(f *pdfFont) bdf.Font {
	if f.finalized {
		return f.final
	}
	f.finalized = true
	f.final = bdf.SystemFont(f.family, f.weight, f.style)
	if f.prog == nil || f.type3 {
		return f.final
	}
	cm := map[uint32]uint16{}
	for _, u := range f.used {
		if u.char != 0 && u.gid > 0 {
			cm[uint32(u.char)] = uint16(u.gid)
		}
	}
	cm[0] = 0
	if len(cm) == 1 {
		return f.final
	}
	family := f.baseFont
	if i := strings.Index(family, "+"); i == 6 {
		family = family[7:]
	}
	if family == "" {
		family = "PDFFont"
	}
	var data []byte
	prog := f.prog
	switch {
	case prog.sf != nil:
		if prog.sf.IsCFF && prog.sf.Tables["CFF "] == nil {
			return f.final
		}
		if !prog.sf.IsCFF && !c.opts.NoSubset {
			keep := make(map[uint16]bool, len(cm))
			for _, gid := range cm {
				keep[gid] = true
			}
			prog.sf.PruneGlyphs(keep)
		}
		data = prog.sf.Rebuild(cm, family, int(f.weight), f.style != bdf.StyleNormal)
	case prog.cff != nil:
		sf := &sfnt.Font{Tables: map[string][]byte{"CFF ": prog.data}, IsCFF: true, NumGlyphs: prog.cff.numGlyphs, UnitsPerEm: 1000}
		if m := prog.cff.fontMatrix; m[0] != 0 {
			sf.UnitsPerEm = int(1/m[0] + 0.5)
		}
		sf.Ascent, sf.Descent = int16(f.ascent), int16(f.descent)
		if sf.Ascent == 0 {
			sf.Ascent, sf.Descent = int16(sf.UnitsPerEm*8/10), int16(-sf.UnitsPerEm*2/10)
		}
		adv := make([]uint16, sf.NumGlyphs)
		for i := range adv {
			adv[i] = uint16(sf.UnitsPerEm / 2)
		}
		for code, u := range f.used {
			if u.gid < len(adv) {
				w := f.width(glyphCode{code: code, cid: f.cidOf(code)})
				adv[u.gid] = uint16(w * float64(sf.UnitsPerEm))
			}
		}
		sf.Advances = adv
		data = sf.Rebuild(cm, family, int(f.weight), f.style != bdf.StyleNormal)
	default:
		return f.final
	}
	h := c.doc.AddFont(data)
	f.final = bdf.EmbeddedFont(h, f.weight, f.style)
	f.final.Family = f.family
	return f.final
}

func (f *pdfFont) cidOf(code uint32) uint32 {
	if f.composite && f.enc != nil {
		return f.enc.cid(code)
	}
	return code
}

// standardWidth approximates advances for the standard 14 fonts without Widths.
func standardWidth(f *pdfFont, code uint32) float64 {
	if f.family == "monospace" {
		return 0.6
	}
	if code == 32 {
		return 0.278
	}
	if code >= 'A' && code <= 'Z' {
		return 0.667
	}
	if code >= '0' && code <= '9' {
		return 0.556
	}
	return 0.5
}

func (f *pdfFont) String() string {
	return fmt.Sprintf("%s(%s)", f.baseFont, f.subtype)
}
