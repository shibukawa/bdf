package font

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/internal/fontdb"
	"github.com/shibukawa/bdf/internal/otlayout"
	"github.com/shibukawa/bdf/woff2"
)

// converter makes the document of a font file.
type converter struct {
	o    *conv.Options
	warn func(string)
	doc  *bdf.Document
	set  *fontset.Set
	ui   *texter
	cvs  *canvas.Builder
	fl   *file

	// the font being drawn
	fc    *face
	paint *painter
	plan  samplePlan

	text     string // sample text of the waterfall ("" for the pangram)
	examples int    // substitutions shown per feature (0: all)

	views  []*bdf.View
	strips []*strip
	// glyphFonts are the glyph fonts embedded, by the glyph data they draw
	// (see glyphKey): the fonts of a collection share them.
	glyphFonts map[string]bdf.Font
	// fontParts counts the fonts embedded to draw the fonts shown.
	fontParts int
	// outlined counts the fonts drawn as outlines.
	outlined int
	// cut are the layout tables a warning was given for (see face).
	cut map[*otlayout.Table]bool
}

// defaultExamples is how many substitutions a feature shows by default.
const defaultExamples = 200

func convert(data []byte, o *conv.Options, warn func(string)) (*conv.Result, error) {
	fl, err := load(data)
	if err != nil {
		return nil, err
	}
	c := &converter{o: o, warn: warn, fl: fl, doc: bdf.NewDocument(), text: o.Param("text"), examples: defaultExamples}
	if v := o.Param("examples"); v != "" {
		if v == "all" {
			c.examples = 0
		} else if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			c.examples = n
		} else {
			return nil, fmt.Errorf("font: parameter examples: %q is not a number or all", v)
		}
	}
	faces, err := c.selectFaces(o.Param("font"))
	if err != nil {
		return nil, err
	}
	// the fonts of a collection that are not shown are not read further
	for _, fc := range faces {
		fc.prepare()
	}
	c.doc.Meta.Source = "font"
	c.set = fontset.New(fontdb.New(o.FontFS, o.FontDirs, !o.NoSystemFonts), warn)
	c.ui = &texter{set: c.set}
	c.cvs = canvas.NewBuilder(c.doc, c.set)

	first := faces[0]
	title := first.fullName()
	if len(faces) > 1 {
		title = first.family()
	}
	if o.Title != "" {
		title = o.Title
	}
	if title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{title}
	}
	for _, id := range []uint16{9, 8} {
		if s := first.name(id); s != "" {
			c.doc.Meta.DC.Creator = bdf.DCValues{s}
			break
		}
	}
	if s := first.name(0); s != "" {
		c.doc.Meta.DC.Rights = bdf.DCValues{s}
	}
	if s := first.name(10); s != "" {
		c.doc.Meta.DC.Description = bdf.DCValues{s}
	}
	c.doc.Meta.DC.Type = bdf.DCValues{"Font"}
	if head := first.f.Tables["head"]; len(head) >= 54 {
		if t, ok := headTime(head, 20); ok {
			c.doc.Meta.DC.Created = bdf.DCValues{t.Format(time.RFC3339)}
		}
		if t, ok := headTime(head, 28); ok {
			c.doc.Meta.DC.Modified = bdf.DCValues{t.Format(time.RFC3339)}
		}
	}

	for _, fc := range faces {
		c.face(fc, len(faces) > 1)
	}
	embedded := c.fontParts
	if !o.SystemFonts {
		embedded += c.set.Embed(c.doc, fontset.EmbedOptions{NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType})
	}
	c.cvs.Encode()
	c.set.ReportMissing()
	for _, s := range c.strips {
		s.page.Layers = []bdf.Layer{{Role: bdf.RoleBody, Obj: s.cv.Hash()}}
	}
	if !o.NoTextIndex {
		for _, v := range c.views {
			if _, err := c.doc.BuildTextIndex(v); err != nil {
				warn(fmt.Sprintf("text index: %v", err))
			}
		}
	}
	var sum []string
	if len(faces) > 1 {
		sum = append(sum, fmt.Sprintf("%d fonts", len(faces)))
	}
	sum = append(sum, fmt.Sprintf("%s glyphs, %s characters", num(first.f.NumGlyphs), num(len(first.runes))))
	if c.outlined > 0 {
		sum = append(sum, "drawn as outlines")
	}
	sum = append(sum, fmt.Sprintf("%d view(s), %d embedded font(s)", len(c.views), embedded))
	return &conv.Result{Doc: c.doc, Summary: strings.Join(sum, ", ")}, nil
}

// maxGlyphs bounds the glyphs of the fonts of a collection shown when the
// font parameter does not name them: those of four fonts of 65,536 glyphs.
const maxGlyphs = 1 << 18

// selectFaces picks the fonts of a collection the font parameter names
// (1-based numbers separated by commas, or all). By default the fonts are
// shown from the first while their glyphs number at most maxGlyphs.
func (c *converter) selectFaces(v string) ([]*face, error) {
	if v == "all" {
		return c.fl.faces, nil
	}
	if v == "" {
		n, out := 0, []*face(nil)
		for _, fc := range c.fl.faces {
			if n += fc.f.NumGlyphs; n > maxGlyphs && len(out) > 0 {
				c.warn(fmt.Sprintf("font: the collection has %d fonts; the first %d are shown (-param font=N or font=all shows others)", len(c.fl.faces), len(out)))
				break
			}
			out = append(out, fc)
		}
		return out, nil
	}
	var out []*face
	for _, f := range strings.Split(v, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(f))
		if err != nil {
			return nil, fmt.Errorf("font: parameter font: %q is not a number", f)
		}
		found := false
		for _, fc := range c.fl.faces {
			if fc.index == n-1 {
				out = append(out, fc)
				found = true
			}
		}
		if !found {
			return nil, fmt.Errorf("font: parameter font: the file has no font %d", n)
		}
	}
	return out, nil
}

// face makes the views of one font.
func (c *converter) face(fc *face, several bool) {
	c.fc = fc
	c.plan = fc.plan(c.text)
	c.paint = c.painterFor(fc)
	id, prefix := "", ""
	if several {
		id = "-" + strconv.Itoa(fc.index+1)
		prefix = fc.fullName() + ": "
		if prefix == ": " {
			prefix = "Font " + strconv.Itoa(fc.index+1) + ": "
		}
	}
	c.addView("overview"+id, prefix+"Overview", c.overview)
	if len(fc.runes) > 0 {
		c.addView("characters"+id, prefix+"Characters", c.characters)
	}
	c.addView("glyphs"+id, prefix+"Glyphs", c.glyphView)
	if fc.hasFeatures() {
		c.addView("features"+id, prefix+"Features", c.features)
	}
	if fc.namesCut {
		c.warn(fmt.Sprintf("%s: the records of the name table are longer than %d MB together; those past it are not shown", fc.fullName(), maxNames>>20))
	}
	cut := false
	for _, t := range []*otlayout.Table{fc.gsub, fc.gpos} {
		if t != nil && t.Truncated() && !c.cut[t] {
			if c.cut == nil {
				c.cut = map[*otlayout.Table]bool{}
			}
			c.cut[t], cut = true, true
		}
	}
	if cut {
		c.warn(fmt.Sprintf("%s: the lookups of the layout tables list more than any font does; the features are not shown in full", fc.fullName()))
	}
}

func (c *converter) addView(id, title string, build func(s *scroll)) {
	v := c.doc.NewView(id, bdf.ViewScroll, title)
	s := newScroll(c)
	build(s)
	c.strips = append(c.strips, s.finish(v)...)
	c.views = append(c.views, v)
}

// painterFor embeds what draws the glyphs of fc: the glyph font and the
// sample font, unless the font's license forbids embedding it, when the
// glyphs are drawn as outlines.
func (c *converter) painterFor(fc *face) *painter {
	p := &painter{fc: fc}
	embed, subset := true, true
	if os2, ok := fc.f.OS2(); ok {
		embed, subset = os2.Embeddable()
	}
	if !embed && !c.o.IgnoreFSType {
		c.outlined++
		c.warn(fmt.Sprintf("%s: the font's license (OS/2 fsType) does not allow embedding it; its glyphs are drawn as outlines", fc.fullName()))
		return p
	}
	key := glyphKey(fc)
	g, ok := c.glyphFonts[key]
	if !ok {
		g = bdf.EmbeddedFont(c.addFont(glyphProgram(fc)), 400, bdf.StyleNormal)
		if c.glyphFonts == nil {
			c.glyphFonts = map[string]bdf.Font{}
		}
		c.glyphFonts[key] = g
	}
	p.glyphFont = &g
	if data, ok := sampleProgram(fc, c.plan.all(), subset || c.o.IgnoreFSType); ok {
		s := bdf.EmbeddedFont(c.addFont(data), 400, bdf.StyleNormal)
		p.sampleFont = &s
	}
	return p
}

// addFont adds a font part, as WOFF2 unless told otherwise.
func (c *converter) addFont(data []byte) bdf.Hash {
	if !c.o.NoWOFF2 {
		if w, err := woff2.Encode(data); err == nil {
			data = w
		} else if err != woff2.ErrNotAvailable && err != woff2.ErrImplausible {
			c.warn(fmt.Sprintf("font: not stored as WOFF2: %v", err))
		}
	}
	c.fontParts++
	return c.doc.AddFont(data)
}

// hasFeatures reports whether the font has OpenType layout features.
func (fc *face) hasFeatures() bool {
	return fc.gsub != nil && len(fc.gsub.Features) > 0 || fc.gpos != nil && len(fc.gpos.Features) > 0
}

// glyphTables are the tables that decide how the glyph font draws glyphs.
var glyphTables = []string{"glyf", "loca", "CFF ", "CFF2", "hmtx", "cvt ", "fpgm", "prep", "gasp",
	"COLR", "CPAL", "sbix", "CBDT", "CBLC", "SVG ", "EBDT", "EBLC", "fvar", "avar"}

// glyphKey identifies the glyph data of a font: the fonts of a collection
// that share those tables (the same bytes of the file) and the numbers
// that read them draw with one glyph font.
func glyphKey(fc *face) string {
	var b strings.Builder
	for _, tag := range glyphTables {
		if t := fc.f.Tables[tag]; len(t) > 0 {
			fmt.Fprintf(&b, "%s:%p:%d;", tag, &t[0], len(t))
		}
	}
	head, hhea := fc.f.Tables["head"], fc.f.Tables["hhea"]
	fmt.Fprintf(&b, "%d:%d:%d:%d", fc.f.NumGlyphs, be16(head, 18), be16(head, 50), be16(hhea, 34))
	return b.String()
}
