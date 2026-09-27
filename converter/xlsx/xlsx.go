// Package xlsx converts Excel workbooks (.xlsx) into BDF documents.
//
// Every worksheet becomes a sheet view: the cells are laid out here (BDF has
// no layout engine) on the sheet's grid of rows and columns and drawn into
// square tiles, with their fills, borders, number formats, alignment,
// merged cells, conditional formats and table styles; pictures, shapes and
// charts are drawn by the DrawingML renderer that the Office converters
// share (converter/internal/ooxml/drawingml), each once as an object that
// the tiles it covers use. Text is measured with the fonts that are then
// embedded as subsets. Chart sheets become one-page fixed views. See
// docs/design.md §3.6.
package xlsx

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontdb"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
	"github.com/shibukawa/bdf/imgconv"
)

// Options controls the conversion.
type Options struct {
	// Sheets selects 1-based sheets (in workbook order, chart sheets
	// included; see converter.Pages); nil converts every sheet (hidden ones
	// only with Hidden).
	Sheets conv.Pages
	// Hidden includes hidden sheets when Sheets is nil.
	Hidden bool
	// Title overrides the document title.
	Title string
	// Images controls whether raster images are re-encoded (see imgconv).
	// The zero value keeps images as they are.
	Images imgconv.Options
	// FontFS holds fonts that are not in the local file system; it is
	// searched before FontDirs (see converter.Options.FontFS).
	FontFS fs.FS
	// FontDirs are searched for fonts before the system font directories.
	FontDirs []string
	// NoSystemFonts restricts font lookup to FontFS and FontDirs.
	NoSystemFonts bool
	// SystemFonts refers to fonts by family name instead of embedding the
	// fonts used for layout. Viewers then substitute their own fonts; the
	// text advances keep the text at the widths computed here.
	SystemFonts bool
	// NoSubset embeds whole fonts instead of the glyphs in use.
	NoSubset bool
	// NoWOFF2 stores embedded fonts as TrueType/OpenType instead of WOFF2
	// (builds tagged bdf_noconv never produce WOFF2).
	NoWOFF2 bool
	// IgnoreFSType embeds fonts whose OS/2 fsType forbids embedding or
	// subsetting. Set it only when you hold the rights to embed the fonts.
	IgnoreFSType bool
	// NoTextIndex skips building the text index parts.
	NoTextIndex bool
	// Warn receives non-fatal problems; when nil they are collected in Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	Sheets        int
	EmbeddedFonts int
}

type converter struct {
	pkg      *ooxml.Package
	opts     *Options
	doc      *bdf.Document
	fonts    *fontset.Set
	cvs      *canvas.Builder
	r        *drawingml.Renderer
	d        *drawingml.Drawing // the workbook's theme and color map
	warnings []string
	warned   map[string]bool

	wbPart   string
	date1904 bool
	sst      []*richText
	st       *styles
	locale   string // "ja" when the workbook is Japanese (built-in formats differ)
	eaScript string // script of the East Asian theme fonts (Jpan, Hang …)
	eaLang   string // language of East Asian text without kana or hangul ("" unknown)
	mdw      float64
	patterns map[string]bdf.Hash

	sheets       []sheetRef
	parsed       map[string]*worksheet // worksheets by name, for chart references
	chartsFilled map[string]bool

	embeddedFonts int
}

// sheetRef is a sheet of the workbook.
type sheetRef struct {
	name   string
	part   string
	chart  bool
	hidden bool
}

// ConvertFile converts a .xlsx file.
func ConvertFile(path string, opts *Options) (*Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	return Convert(f, st.Size(), opts)
}

// Convert converts a workbook read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	p, err := ooxml.Open(r, size)
	if err != nil {
		return nil, fmt.Errorf("xlsx: %w", err)
	}
	c := newConverter(p, opts)
	c.r = drawingml.New(drawingml.Config{Package: p, Doc: c.doc, Fonts: c.fonts, Images: opts.Images, Warn: func(msg string) { c.warnf("%s", msg) }})

	c.wbPart = "xl/workbook.xml"
	if r, ok := p.RelOfType("", "/officeDocument"); ok {
		c.wbPart = r.Target
	}
	wb, err := p.XML(c.wbPart)
	if err != nil {
		return nil, fmt.Errorf("xlsx: %w", err)
	}
	c.date1904 = wb.Child("workbookPr").AttrBool("date1904", false)
	c.d = c.r.NewDrawing(c.wbPart, nil, host{})
	c.loadSharedStrings()
	sn, err := c.partOfType(c.wbPart, "/styles")
	var stylesNode *ooxml.Node
	if err == nil {
		stylesNode = sn
	}
	c.st = loadStyles(c, stylesNode, func(name string) (rgb, bool) {
		r, g, b, ok := c.d.ThemeColor(name)
		return rgb{r, g, b}, ok
	})
	c.doc.Meta.Source = "xlsx"
	c.doc.Meta.DC = p.CoreProperties()
	if opts.Title != "" {
		c.doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	c.setLocale()
	c.mdw = c.maxDigitWidth(c.st.font(0))

	var sheets []sheetRef
	for _, s := range wb.Path("sheets").Children("sheet") {
		r, ok := p.Target(c.wbPart, s.RelID("id"))
		if !ok {
			continue
		}
		state := s.AttrStr("state", "visible")
		ref := sheetRef{name: s.AttrStr("name", ""), part: r.Target, hidden: state != "visible"}
		switch {
		case strings.HasSuffix(r.Type, "/worksheet"):
		case strings.HasSuffix(r.Type, "/chartsheet"):
			ref.chart = true
		default:
			c.warnOnce("sheettype:"+r.Type, "sheet %q: %s sheets are not supported", ref.name, r.Type[strings.LastIndexByte(r.Type, '/')+1:])
			continue
		}
		sheets = append(sheets, ref)
	}
	c.sheets = sheets
	sel := opts.Sheets.Numbers(len(sheets))
	if sel == nil {
		for i, s := range sheets {
			if !s.hidden || opts.Hidden {
				sel = append(sel, i+1)
			}
		}
	}
	var views []pendingView
	for _, n := range sel {
		if n < 1 || n > len(sheets) {
			return nil, fmt.Errorf("xlsx: sheet %d out of range (1-%d)", n, len(sheets))
		}
		ref := sheets[n-1]
		id := "sheet" + strconv.Itoa(n)
		if ref.chart {
			v, layers := c.chartSheetSafe(ref, id)
			views = append(views, pendingView{view: v, layers: layers})
			continue
		}
		v, tiles, err := c.worksheetSafe(ref, id)
		if err != nil {
			return nil, fmt.Errorf("xlsx: sheet %q: %w", ref.name, err)
		}
		views = append(views, pendingView{view: v, tiles: tiles})
	}
	c.finish(views)
	return &Result{Doc: c.doc, Warnings: c.warnings, Sheets: len(sel), EmbeddedFonts: c.embeddedFonts}, nil
}

// newConverter sets up a conversion: the document, the fonts and the
// object builder. p is nil for grids, which come from no workbook.
func newConverter(p *ooxml.Package, opts *Options) *converter {
	c := &converter{pkg: p, opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{}, patterns: map[string]bdf.Hash{},
		parsed: map[string]*worksheet{}, chartsFilled: map[string]bool{}}
	db := fontdb.New(opts.FontFS, opts.FontDirs, !opts.NoSystemFonts)
	c.fonts = fontset.New(db, func(msg string) { c.warnf("%s", msg) })
	c.cvs = canvas.NewBuilder(c.doc, c.fonts)
	if len(db.Faces) == 0 && !opts.SystemFonts {
		c.warnf("no fonts found; text is laid out with estimated metrics and not embedded")
	}
	return c
}

// pendingView is a converted sheet whose objects get their hashes when the
// fonts are known.
type pendingView struct {
	view   *bdf.View
	tiles  map[string]*canvas.Canvas
	layers []*canvas.Canvas
}

// finish embeds the fonts, encodes the objects, points the views at them
// and indexes their text.
func (c *converter) finish(views []pendingView) {
	opts := c.opts
	if !opts.SystemFonts {
		c.embeddedFonts = c.fonts.Embed(c.doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType})
	}
	c.cvs.Encode()
	c.fonts.ReportMissing()
	for _, pv := range views {
		for key, cv := range pv.tiles {
			pv.view.Tiles[key] = cv.Hash().String()
		}
		if len(pv.layers) > 0 {
			for i, cv := range pv.layers {
				pv.view.Pages[0].Layers[i].Obj = cv.Hash()
			}
		}
		if !opts.NoTextIndex {
			if _, err := c.doc.BuildTextIndex(pv.view); err != nil {
				c.warnf("text index: %v", err)
			}
		}
	}
}

func (c *converter) warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if c.opts.Warn != nil {
		c.opts.Warn(msg)
		return
	}
	c.warnings = append(c.warnings, msg)
}

func (c *converter) warnOnce(key, format string, args ...any) {
	if c.warned[key] {
		return
	}
	c.warned[key] = true
	c.warnf(format, args...)
}

// partOfType returns the XML of the part that part relates to by a
// relationship type (suffix).
func (c *converter) partOfType(part, suffix string) (*ooxml.Node, error) {
	r, ok := c.pkg.RelOfType(part, suffix)
	if !ok {
		return nil, fmt.Errorf("no %s part", strings.TrimPrefix(suffix, "/"))
	}
	return c.pkg.XML(r.Target)
}

func (c *converter) loadSharedStrings() {
	n, err := c.partOfType(c.wbPart, "/sharedStrings")
	if err != nil {
		return
	}
	for _, si := range n.Children("si") {
		c.sst = append(c.sst, c.readRich(si))
	}
}

// East Asian fonts that tell the language of a workbook's text.
var eaFonts = []struct {
	names  []string
	script string
	lang   string
}{
	{[]string{"游ゴシック", "游明朝", "yu gothic", "yu mincho", "ｍｓ ｐゴシック", "ｍｓ ゴシック", "ms pgothic", "ms gothic", "ｍｓ 明朝", "ｍｓ ｐ明朝", "ms mincho", "メイリオ", "meiryo", "bizud", "biz ud"}, "Jpan", "ja"},
	{[]string{"맑은 고딕", "malgun gothic", "굴림", "gulim", "돋움", "dotum", "바탕", "batang"}, "Hang", "ko"},
	{[]string{"等线", "dengxian", "宋体", "simsun", "微软雅黑", "microsoft yahei"}, "Hans", "zh-CN"},
	{[]string{"新細明體", "pmingliu", "微軟正黑體", "microsoft jhenghei"}, "Hant", "zh-TW"},
}

// setLocale guesses the workbook's East Asian language from its default
// font (Excel writes no language): Japanese workbooks have their own
// built-in date and currency formats, and their Han-only text is Japanese.
func (c *converter) setLocale() {
	c.eaScript = "Jpan"
	name := strings.ToLower(c.st.font(0).name)
	if l := c.doc.Meta.DC.Language.First(); l != "" {
		name = ""
		if s := fontset.Script(l); s != "" {
			c.eaScript, c.eaLang = s, l
		}
	}
	for _, f := range eaFonts {
		for _, n := range f.names {
			if name != "" && strings.HasPrefix(name, n) {
				c.eaScript, c.eaLang = f.script, f.lang
			}
		}
	}
	if c.eaLang == "" && c.d.ThemeFont("+mn-ea", "") == "" {
		// a theme whose minor East Asian font is set only per script
		if f := c.d.ThemeFont("+mn-ea", "Jpan"); f != "" && strings.EqualFold(f, c.st.font(0).name) {
			c.eaLang = "ja"
		}
	}
	if strings.HasPrefix(c.eaLang, "ja") {
		c.locale = "ja"
	}
}

// fontName returns the family a font draws Latin text with: the theme's
// font for the major and minor scheme fonts.
func (c *converter) fontName(f *xfont) string {
	switch f.scheme {
	case "minor":
		if n := c.d.ThemeFont("+mn-lt", c.eaScript); n != "" {
			return n
		}
	case "major":
		if n := c.d.ThemeFont("+mj-lt", c.eaScript); n != "" {
			return n
		}
	}
	if f.name == "" {
		return "Calibri"
	}
	return f.name
}

// eaFontName returns the family for East Asian text: the theme's East
// Asian font for scheme fonts, else none (the Latin font and the fallbacks).
func (c *converter) eaFontName(f *xfont) string {
	switch f.scheme {
	case "minor":
		return c.d.ThemeFont("+mn-ea", c.eaScript)
	case "major":
		return c.d.ThemeFont("+mj-ea", c.eaScript)
	}
	return ""
}

// host implements drawingml.Host for spreadsheet drawings, which have no
// placeholders; shape text defaults to 11 points.
type host struct{}

var shapeTextStyle = &ooxml.Node{Name: "lstStyle", Kids: []*ooxml.Node{{Name: "defPPr", Kids: []*ooxml.Node{
	{Name: "defRPr", Attrs: attrs("sz", "1100")}}}}}

func (host) Placeholder(*ooxml.Node, string) ([]drawingml.Inherited, bool) { return nil, true }
func (host) TextStyle(string) *ooxml.Node                                  { return shapeTextStyle }
func (host) Field(string) (string, bool)                                   { return "", false }
func (host) Link(string, ooxml.Rel) string                                 { return "" }

func attrs(kv ...string) []xml.Attr {
	var out []xml.Attr
	for i := 0; i+1 < len(kv); i += 2 {
		out = append(out, xml.Attr{Name: xml.Name{Local: kv[i]}, Value: kv[i+1]})
	}
	return out
}
