package font

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/internal/otlayout"
)

// featureGroup is a feature tag of GSUB or GPOS: every feature record of
// the tag, and what they apply.
type featureGroup struct {
	table   string
	t       *otlayout.Table
	tag     string
	lookups []int
	scripts []string // "script" or "script/LANG"
	uiName  string
	chars   []rune
}

// featureGroups returns the features of the font by table and tag.
func (fc *face) featureGroups() []*featureGroup {
	var out []*featureGroup
	for _, tb := range []struct {
		name string
		t    *otlayout.Table
	}{{"GSUB", fc.gsub}, {"GPOS", fc.gpos}} {
		if tb.t == nil {
			continue
		}
		users := tb.t.ScriptFeatures()
		byTag := map[string]*featureGroup{}
		type seen struct {
			lookups map[int]bool
			scripts map[string]bool
			chars   map[rune]bool
		}
		seenOf := map[*featureGroup]*seen{}
		for i, f := range tb.t.Features {
			g := byTag[f.Tag]
			if g == nil {
				g = &featureGroup{table: tb.name, t: tb.t, tag: f.Tag}
				byTag[f.Tag] = g
				seenOf[g] = &seen{map[int]bool{}, map[string]bool{}, map[rune]bool{}}
				out = append(out, g)
			}
			sn := seenOf[g]
			for _, l := range f.Lookups {
				if !sn.lookups[l] {
					sn.lookups[l] = true
					g.lookups = append(g.lookups, l)
				}
			}
			for _, u := range users[i] {
				if !sn.scripts[u] {
					sn.scripts[u] = true
					g.scripts = append(g.scripts, u)
				}
			}
			if f.UIName != 0 && g.uiName == "" {
				g.uiName = fc.name(f.UIName)
			}
			for _, r := range f.Chars {
				if !sn.chars[r] {
					sn.chars[r] = true
					g.chars = append(g.chars, r)
				}
			}
		}
	}
	return out
}

// scriptList describes the scripts and languages of a feature compactly:
// "latn (default, TRK, ROM), cyrl".
func scriptList(users []string) string {
	var order []string
	langs := map[string][]string{}
	for _, u := range users {
		sc, lang, _ := strings.Cut(u, "/")
		sc = strings.TrimSpace(sc)
		if _, ok := langs[sc]; !ok {
			order = append(order, sc)
			langs[sc] = nil
		}
		if lang == "" {
			lang = "default"
		}
		langs[sc] = append(langs[sc], strings.TrimSpace(lang))
	}
	var out []string
	for _, sc := range order {
		ls := langs[sc]
		if len(ls) == 1 && ls[0] == "default" {
			out = append(out, sc)
			continue
		}
		out = append(out, sc+" ("+strings.Join(ls, ", ")+")")
	}
	return strings.Join(out, ", ")
}

// scriptTags lists the scripts of a feature without their languages.
func scriptTags(users []string) string {
	var out []string
	seen := map[string]bool{}
	for _, u := range users {
		sc, _, _ := strings.Cut(u, "/")
		if sc = strings.TrimSpace(sc); !seen[sc] {
			seen[sc] = true
			out = append(out, sc)
		}
	}
	return strings.Join(out, " ")
}

// defaultNote tells when browsers apply a feature.
func defaultNote(tag string) string {
	switch otlayout.FeatureDefault(tag) {
	case otlayout.Default:
		return "on by default"
	case otlayout.ScriptDefault:
		return "applied by the script's shaper"
	case otlayout.Vertical:
		return "in vertical text"
	}
	return ""
}

// features adds the rows of the features view.
func (c *converter) features(s *scroll) {
	fc := c.fc
	groups := fc.featureGroups()
	s.add(30, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkHeading, "1")
		p.label("OpenType features", margin, y+22, 18, true, colText, 0)
	})
	var parts []string
	for _, tb := range []struct {
		name string
		t    *otlayout.Table
	}{{"GSUB", fc.gsub}, {"GPOS", fc.gpos}} {
		if tb.t != nil {
			n := 0
			for _, g := range groups {
				if g.table == tb.name {
					n++
				}
			}
			p := fmt.Sprintf("%s: %d features, %d lookups", tb.name, n, len(tb.t.Lookups))
			if tb.t.FeatureVariations > 0 {
				p += fmt.Sprintf(", %d feature variations", tb.t.FeatureVariations)
			}
			parts = append(parts, p)
		}
	}
	info := strings.Join(parts, " · ")
	s.add(18, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkParagraph, "")
		p.label(info, margin, y+12, bodySize, false, colGray, 0)
	})
	// the summary, with links to the details
	c.heading(s, "Summary")
	for i, g := range groups {
		s.add(15, func(p *pen, y float64) {
			p.obj().Mark(bdf.MarkParagraph, "")
			w := p.label(g.tag, margin, y+11, bodySize, true, colAccent, 0)
			p.linkTo(margin, y, w, 15, "feature"+strconv.Itoa(i))
			name := otlayout.FeatureName(g.tag)
			if g.uiName != "" {
				name = g.uiName
			}
			p.label(c.ui.fit(name, bodySize, false, 190), margin+42, y+11, bodySize, false, colText, 0)
			p.label(g.table, margin+240, y+11, smallSize, false, colGray, 0)
			p.label(defaultNote(g.tag), margin+275, y+11, smallSize, false, colAccent, 0)
			p.label(c.ui.fit(scriptTags(g.scripts), smallSize, false, contentW-420), margin+420, y+11, smallSize, false, colGray, 0)
		})
	}
	for i, g := range groups {
		c.featureDetail(s, g, "feature"+strconv.Itoa(i))
	}
}

// featureDetail adds a feature's heading, what it applies, and examples.
func (c *converter) featureDetail(s *scroll, g *featureGroup, anchor string) {
	s.space(20)
	s.anchor(anchor)
	name := otlayout.FeatureName(g.tag)
	h := s.add(24, func(p *pen, y float64) {
		p.obj().Mark(bdf.MarkHeading, "2")
		w := p.label(g.tag, margin, y+17, 14, true, colText, 0)
		if name != "" {
			w += p.label(name, margin+w+10, y+17, 12, false, colText, 0) + 10
		}
		if d := defaultNote(g.tag); d != "" {
			t := c.ui.line(d, smallSize, false)
			bx := margin + w + 12
			p.rect(bx, y+7, t.width+10, 13, colBadge)
			p.text(t, bx+5, y+16.5, colAccent)
		}
		p.rect(margin, y+22, contentW, 0.5, colRule)
	})
	h.keep = true
	var info []string
	if g.uiName != "" {
		info = append(info, "Name: "+g.uiName)
	}
	if len(g.chars) > 0 {
		var cs []string
		for _, r := range g.chars {
			cs = append(cs, string(r))
		}
		info = append(info, "Characters: "+strings.Join(cs, " "))
	}
	info = append(info, g.table+" scripts: "+scriptList(g.scripts))
	var lks []string
	for _, l := range g.lookups {
		lks = append(lks, c.lookupNote(g.t, l))
	}
	if len(lks) == 0 {
		lks = []string{"none"}
	}
	info = append(info, "Lookups: "+strings.Join(lks, "; "))
	for _, line := range info {
		for _, l := range c.ui.wrap(line, smallSize+0.5, false, contentW) {
			b := s.add(12, func(p *pen, y float64) {
				p.obj().Mark(bdf.MarkParagraph, "")
				p.label(l, margin, y+9, smallSize+0.5, false, colGray, 0)
			})
			b.keep = true
		}
	}
	s.space(6)
	items, total := c.featureExamples(g)
	c.itemRows(s, items)
	if total > len(items) {
		more := fmt.Sprintf("… and %s more", num(total-len(items)))
		s.add(14, func(p *pen, y float64) {
			p.obj().Mark(bdf.MarkParagraph, "")
			p.label(more, margin, y+10, smallSize+0.5, false, colGray, 0)
		})
	}
}

// lookupNote describes a lookup: its index, type and size.
func (c *converter) lookupNote(t *otlayout.Table, i int) string {
	l := t.Lookups[i]
	s := "#" + strconv.Itoa(i) + " " + t.TypeName(l.Type)
	if ctx, ok := t.Context(i); ok && len(ctx.Nested) > 0 {
		var ns []string
		for _, n := range ctx.Nested {
			ns = append(ns, "#"+strconv.Itoa(n))
		}
		rules := "rules"
		if ctx.Rules == 1 {
			rules = "rule"
		}
		s += fmt.Sprintf(" (%s %s, applying %s)", num(ctx.Rules), rules, strings.Join(ns, " "))
		return s
	}
	n := 0
	switch {
	case t == c.fc.gsub:
		t.Substs(i, func(otlayout.Subst) bool { n++; return true })
	case l.Type == otlayout.PosPair:
		n = t.PairCount(i)
	case l.Type == otlayout.PosSingle:
		t.Singles(i, func(uint16, otlayout.Value) bool { n++; return true })
	default:
		m, b := t.MarkCounts(i)
		if m > 0 {
			return s + fmt.Sprintf(" (%s marks, %s bases)", num(m), num(b))
		}
		n = b
	}
	if n > 0 {
		s += " (" + num(n) + ")"
	}
	return s
}

// placed is a glyph at a position, in em from the start of its group (y
// up).
type placed struct {
	g    uint16
	x, y float64
}

// item is an example of a feature: glyphs before and after it applies.
type item struct {
	before, after []placed
	bw, aw        float64 // widths in em
	// alts separates the glyphs after as choices.
	alts bool
	// boxes draws the advance of the single glyph before and after it is
	// adjusted.
	boxes  bool
	advB   float64
	advA   float64
	offset float64
	// note is a value shown under the glyphs after (a kerning value).
	note string
}

// extent returns the horizontal extent of placed glyphs in em: from the
// origin to the advance w, widened to their ink.
func (c *converter) extent(ps []placed, w float64) (x0, x1 float64) {
	x1 = w
	for _, pl := range ps {
		if b := c.fc.bounds[pl.g]; b.ok {
			x0 = min(x0, pl.x+b.x0)
			x1 = max(x1, pl.x+b.x1)
		}
	}
	return x0, x1
}

// seq places glyphs at their advances.
func (c *converter) seq(gs []uint16) ([]placed, float64) {
	out := make([]placed, len(gs))
	x := 0.0
	for i, g := range gs {
		out[i] = placed{g: g, x: x}
		x += c.fc.advances[g]
	}
	return out, x
}

// featureExamples collects the examples of a feature, at most c.examples of them,
// and how many there are.
func (c *converter) featureExamples(g *featureGroup) (items []item, total int) {
	limit := c.examples
	add := func(it item) {
		total++
		if limit == 0 || len(items) < limit {
			items = append(items, it)
		}
	}
	seen := map[int]bool{}
	var visit func(l int, depth int)
	visit = func(l int, depth int) {
		if seen[l] || depth > 8 || l < 0 || l >= len(g.t.Lookups) {
			return
		}
		seen[l] = true
		t := g.t
		if ctx, ok := t.Context(l); ok && len(ctx.Nested) > 0 {
			for _, n := range ctx.Nested {
				visit(n, depth+1)
			}
			return
		}
		if g.table == "GSUB" {
			t.Substs(l, func(s otlayout.Subst) bool {
				if !c.fc.valid(s.In...) || !c.fc.valid(s.Out...) {
					return true // a damaged table
				}
				it := item{}
				it.before, it.bw = c.seq(s.In)
				if t.Lookups[l].Type == otlayout.SubstAlternate {
					it.alts = true
					x := 0.0
					for _, o := range s.Out {
						it.after = append(it.after, placed{g: o, x: x})
						x += c.fc.advances[o] + 0.25
					}
					it.aw = max(x-0.25, 0)
				} else {
					it.after, it.aw = c.seq(s.Out)
				}
				add(it)
				return true
			})
			return
		}
		switch t.Lookups[l].Type {
		case otlayout.PosPair:
			c.pairExamples(t, l, add, &total)
		case otlayout.PosSingle:
			c.singleExamples(t, l, add)
		case otlayout.PosMarkToBase, otlayout.PosMarkToMark:
			c.markExamples(t, l, add)
		}
	}
	for _, l := range g.lookups {
		visit(l, 0)
	}
	return items, total
}

// pairExamples shows the pairs of a kerning lookup found among pairs of
// the font's sample characters, largest adjustments first; the total is
// the lookup's pair count.
func (c *converter) pairExamples(t *otlayout.Table, l int, add func(item), total *int) {
	fc := c.fc
	text := c.plan.all() + "AVAWAYATAvAwAyLTLVLWLYPAPaTaTeToTrTuTyVaVeVoWaWeWoYaYeYoFAFaKoRoLyffrvyvwkv.y.,T.V.Y.F.P.r,y,"
	var chars []rune
	for _, r := range text {
		if unicode.IsGraphic(r) && !unicode.IsSpace(r) && fc.glyph(r) != 0 {
			chars = append(chars, r)
		}
	}
	type cand struct {
		a, b   uint16
		v1, v2 otlayout.Value
	}
	var found []cand
	seen := map[[2]uint16]bool{}
	try := func(a, b uint16) {
		if seen[[2]uint16{a, b}] || !fc.valid(a, b) {
			return
		}
		seen[[2]uint16{a, b}] = true
		if v1, v2, ok := t.PairValue(l, a, b); ok && (!v1.Zero() || !v2.Zero()) {
			found = append(found, cand{a, b, v1, v2})
		}
	}
	for i := 0; i+1 < len(chars); i++ {
		try(fc.glyph(chars[i]), fc.glyph(chars[i+1]))
	}
	if len(found) < 24 {
		// pairs from the lookup itself, of glyphs mapped from characters
		n := 0
		t.Pairs(l, func(p otlayout.Pair) bool {
			if fc.valid(p.First, p.Second) && len(fc.chars[p.First]) > 0 && len(fc.chars[p.Second]) > 0 {
				try(p.First, p.Second)
			}
			n++
			return len(found) < 48 && n < 200000
		})
	}
	slices.SortStableFunc(found, func(x, y cand) int {
		return abs(y.v1.XAdvance+y.v2.XPlacement) - abs(x.v1.XAdvance+x.v2.XPlacement)
	})
	upem := fc.upem
	count := t.PairCount(l)
	shown := min(len(found), 48)
	for _, f := range found[:shown] {
		it := item{}
		it.before, it.bw = c.seq([]uint16{f.a, f.b})
		a0 := fc.advances[f.a]
		it.after = []placed{
			{g: f.a, x: float64(f.v1.XPlacement) / upem, y: float64(f.v1.YPlacement) / upem},
			{g: f.b, x: a0 + float64(f.v1.XAdvance+f.v2.XPlacement)/upem, y: float64(f.v2.YPlacement) / upem},
		}
		it.aw = a0 + float64(f.v1.XAdvance+f.v2.XAdvance)/upem + fc.advances[f.b]
		if v := f.v1.XAdvance + f.v2.XPlacement; v != 0 {
			it.note = strconv.Itoa(v)
		}
		add(it)
	}
	// the lookup has more pairs than those shown
	*total += max(count-shown, 0)
}

func abs(v int) int { return max(v, -v) }

// singleExamples shows how a single adjustment lookup moves glyphs, glyphs
// mapped from characters first.
func (c *converter) singleExamples(t *otlayout.Table, l int, add func(item)) {
	fc := c.fc
	var mapped, other []item
	t.Singles(l, func(g uint16, v otlayout.Value) bool {
		if !fc.valid(g) {
			return true
		}
		a := fc.advances[g]
		it := item{boxes: true, advB: a, advA: a + float64(v.XAdvance)/fc.upem, offset: float64(v.XPlacement) / fc.upem}
		it.before = []placed{{g: g}}
		it.after = []placed{{g: g, x: it.offset, y: float64(v.YPlacement) / fc.upem}}
		it.bw, it.aw = a, max(it.advA, a)
		if len(fc.chars[g]) > 0 {
			mapped = append(mapped, it)
		} else {
			other = append(other, it)
		}
		return true
	})
	for _, it := range append(mapped, other...) {
		add(it)
	}
}

// markExamples shows marks attached to bases: for each mark class, a mark
// on a few bases, those mapped from characters first.
func (c *converter) markExamples(t *otlayout.Table, l int, add func(item)) {
	fc := c.fc
	upem := fc.upem
	for _, at := range t.Attachments(l) {
		// a mark of each class, mapped ones first
		byClass := map[int]uint16{}
		var classes []int
		marks := make([]uint16, 0, len(at.Marks))
		for m := range at.Marks {
			if fc.valid(m) {
				marks = append(marks, m)
			}
		}
		slices.SortFunc(marks, func(a, b uint16) int {
			ma, mb := len(fc.chars[a]) > 0, len(fc.chars[b]) > 0
			if ma != mb {
				if ma {
					return -1
				}
				return 1
			}
			return int(a) - int(b)
		})
		for _, m := range marks {
			cl := at.Marks[m].Class
			if _, ok := byClass[cl]; !ok {
				byClass[cl] = m
				classes = append(classes, cl)
			}
		}
		slices.Sort(classes)
		bases := make([]uint16, 0, len(at.Bases))
		for b := range at.Bases {
			if fc.valid(b) {
				bases = append(bases, b)
			}
		}
		// letters of the sample text first, then other letters, then the
		// other glyphs mapped from characters
		sample := c.plan.all()
		rank := func(g uint16) int {
			r := firstChar(fc.chars[g])
			switch {
			case r == 0:
				return 3
			case unicode.IsLetter(r) && strings.ContainsRune(sample, r):
				return 0
			case unicode.IsLetter(r):
				return 1
			}
			return 2
		}
		slices.SortFunc(bases, func(a, b uint16) int {
			if ka, kb := rank(a), rank(b); ka != kb {
				return ka - kb
			}
			if ra, rb := firstChar(fc.chars[a]), firstChar(fc.chars[b]); ra != rb {
				return int(ra) - int(rb)
			}
			return int(a) - int(b)
		})
		for _, cl := range classes {
			m := byClass[cl]
			ma := at.Marks[m].Anchor
			shown := 0
			for _, b := range bases {
				an := at.Bases[b]
				if cl >= len(an) || an[cl] == nil {
					continue
				}
				it := item{}
				it.before, it.bw = c.seq([]uint16{b, m})
				ab := fc.advances[b]
				it.after = []placed{{g: b}, {g: m, x: float64(an[cl].X-ma.X) / upem, y: float64(an[cl].Y-ma.Y) / upem}}
				it.aw = ab
				if bx := fc.bounds[m]; bx.ok {
					it.aw = max(ab, it.after[1].x+bx.x1)
				}
				add(it)
				if shown++; shown == 4 {
					break
				}
			}
		}
	}
}

func firstChar(rs []rune) rune {
	if len(rs) == 0 {
		return 0
	}
	return rs[0]
}

// The examples.
const (
	itemSize = 22.0
	itemH    = 42.0
	arrowW   = 20.0
	itemPad  = 8.0
)

// itemRows lays the examples out in rows.
func (c *converter) itemRows(s *scroll, items []item) {
	type laid struct {
		it item
		x  float64
		w  float64
	}
	var row []laid
	x := 0.0
	flush := func() {
		if len(row) == 0 {
			return
		}
		r := row
		s.add(itemH+4, func(p *pen, y float64) {
			for _, l := range r {
				c.drawItem(p, l.it, margin+l.x, y, l.w)
			}
		})
		row, x = nil, 0
	}
	for _, it := range items {
		b0, b1 := c.extent(it.before, it.bw)
		a0, a1 := c.extent(it.after, it.aw)
		w := (b1-b0+a1-a0)*itemSize + arrowW + 2*itemPad
		w = max(w, 56)
		if x > 0 && x+w > contentW {
			flush()
		}
		if w > contentW {
			w = contentW
		}
		row = append(row, laid{it, x, w})
		x += w + 4
	}
	flush()
}

// drawItem draws an example in a box of width w: the glyphs before, an
// arrow, the glyphs after.
func (c *converter) drawItem(p *pen, it item, x, y, w float64) {
	fc := c.fc
	p.obj().Mark(bdf.MarkParagraph, "")
	p.rect(x, y, w, itemH, colEmpty)
	asc, desc := fc.vertical()
	size := itemSize
	base := y + (itemH-(asc+desc)*size)/2 + asc*size
	draw := func(ps []placed, ox float64, box bool, adv float64) {
		if box {
			p.rect(ox, base-asc*size, adv*size, (asc+desc)*size, colBadge)
		}
		for _, pl := range ps {
			p.glyph(pl.g, ox+pl.x*size, base-pl.y*size, size, altOf(fc, pl.g))
		}
	}
	b0, b1 := c.extent(it.before, it.bw)
	a0, _ := c.extent(it.after, it.aw)
	bx := x + itemPad - b0*size
	draw(it.before, bx, it.boxes, it.advB)
	ax := x + itemPad + (b1-b0)*size
	p.label("→", ax+arrowW/2, base-size*0.25, 10, false, colGray, 1)
	ox := ax + arrowW - a0*size
	draw(it.after, ox, it.boxes, it.advA)
	if it.alts && len(it.after) > 1 {
		// dots between the choices
		for i := 1; i < len(it.after); i++ {
			dx := ox + (it.after[i].x-0.125)*size
			p.rect(dx-0.75, base+desc*size*0.5, 1.5, 1.5, colGray)
		}
	}
	if it.note != "" {
		p.label(it.note, x+w-3, y+itemH-3, tinySize, false, colGray, 2)
	}
}

// altOf is the text a glyph stands for in an example: its character when
// it has one.
func altOf(fc *face, g uint16) string {
	if cs := fc.chars[g]; len(cs) > 0 && !unicode.IsControl(cs[0]) {
		return string(cs[0])
	}
	return ""
}
