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
	"fmt"
	"io"
	"io/fs"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter" // the name converter is taken by the conversion state
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
	"github.com/shibukawa/bdf/converter/internal/ooxml/drawingml"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/image/imgconv"
	"github.com/shibukawa/bdf/internal/fontdb"
	"github.com/shibukawa/tinygodriver/encoding/xmlro"
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

	// Views selects the views of a worksheet: ViewsSheet (the default when
	// ""), the sheet as Excel shows it; ViewsPrint, the pages Excel prints
	// (the print areas cut by the page setup, with the print titles, the
	// header and the footer); or ViewsBoth.
	Views string
	// Paper names the paper of the printed pages of sheets whose page setup
	// names none (see PaperNames): "" is A4 for a Japanese workbook and
	// Letter for the others.
	Paper string
	// PrintTime is the date and the time that the headers and footers of
	// printed pages show for &D and &T; the zero time shows none, so that
	// the conversion of a file is the same whenever it is done.
	PrintTime time.Time
	// FileName is the name that headers and footers show for &F.
	FileName string

	// fill, when set, fills the workbook as a template (Fill)
	fill *filler
}

// Result is the outcome of a conversion.
type Result struct {
	Doc           *bdf.Document
	Warnings      []string
	Sheets        int
	EmbeddedFonts int
	// Fields are the names that the tags of a filled template look up
	// (Fill).
	Fields []tag.Field
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
	embeddedGen   uint64
	hasEmbedded   bool

	// printNames are the print areas and titles of the sheets, by the index
	// of the sheet in the workbook
	printNames map[int]*printNames
	// keepCells keeps the cells of a sheet after its text is drawn: its
	// other view is drawn from them
	keepCells   bool
	sheetWarned map[string]bool
}

// sheetRef is a sheet of the workbook.
type sheetRef struct {
	name   string
	part   string
	chart  bool
	hidden bool
	idx    int // its index among the sheets of the workbook (localSheetId)
}

// ConvertFile converts a .xlsx file.
func ConvertFile(path string, opts *Options) (*Result, error) {
	return conv.OpenFile(path, func(r io.ReaderAt, size int64) (*Result, error) { return Convert(r, size, opts) })
}

// openWorkbook opens a workbook for a conversion: its shared strings, its
// styles and its sheets. It returns the root of the workbook part.
func openWorkbook(r io.ReaderAt, size int64, opts *Options) (*converter, *ooxml.Node, error) {
	p, err := ooxml.Open(r, size)
	if err != nil {
		return nil, nil, fmt.Errorf("xlsx: %w", err)
	}
	// equations are laid out from their Office Math, not drawn from the
	// text or pictures of their fallbacks
	p.Choose = ooxml.MathChoice
	c := newConverter(p, opts)
	c.r = drawingml.New(drawingml.Config{Package: p, Doc: c.doc, Fonts: c.fonts, Images: opts.Images, Warn: func(msg string) { c.warnf("%s", msg) }})

	c.wbPart = "xl/workbook.xml"
	if r, ok := p.RelOfType("", "/officeDocument"); ok {
		c.wbPart = r.Target
	}
	wb, err := p.XML(c.wbPart)
	if err != nil {
		return nil, nil, fmt.Errorf("xlsx: %w", err)
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

	for i, s := range wb.Path("sheets").Children("sheet") {
		r, ok := p.Target(c.wbPart, s.RelID("id"))
		if !ok {
			continue
		}
		state := s.AttrStr("state", "visible")
		ref := sheetRef{name: s.AttrStr("name", ""), part: r.Target, hidden: state != "visible", idx: i}
		switch {
		case strings.HasSuffix(r.Type, "/worksheet"):
		case strings.HasSuffix(r.Type, "/chartsheet"):
			ref.chart = true
		default:
			c.warnOnce("sheettype:"+r.Type, "sheet %q: %s sheets are not supported", ref.name, r.Type[strings.LastIndexByte(r.Type, '/')+1:])
			continue
		}
		c.sheets = append(c.sheets, ref)
	}
	return c, wb, nil
}

// Convert converts a workbook read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	views := opts.Views
	switch views {
	case "":
		views = ViewsSheet
		if opts.fill != nil {
			views = ViewsPrint // a filled template is a printed form
		}
	case ViewsSheet, ViewsPrint, ViewsBoth:
	default:
		return nil, fmt.Errorf("xlsx: views %q: want %s, %s or %s", views, ViewsSheet, ViewsPrint, ViewsBoth)
	}
	if opts.Paper != "" {
		if _, ok := paperNames[strings.ToLower(opts.Paper)]; !ok {
			return nil, fmt.Errorf("xlsx: paper %q: want one of %s", opts.Paper, strings.Join(PaperNames(), ", "))
		}
	}
	c, wb, err := openWorkbook(r, size, opts)
	if err != nil {
		return nil, err
	}
	p := c.pkg
	c.keepCells = views == ViewsBoth || opts.fill != nil
	c.loadPrintNames(wb)
	sheets := c.sheets
	if f := opts.fill; f != nil {
		// the template is filled before any sheet is drawn: what a sheet
		// shows of another (a formula, a chart) is of the filled one
		if err := f.run(c); err != nil {
			return nil, err
		}
	}
	sel := opts.Sheets.Numbers(len(sheets))
	if sel == nil {
		for i, s := range sheets {
			if !s.hidden || opts.Hidden {
				sel = append(sel, i+1)
			}
		}
	}
	// On a large multi-sheet book, finalize canvases in batches instead of
	// retaining every sheet's unencoded objects until the last one is drawn.
	var sheetBytes uint64
	for _, n := range sel {
		if n >= 1 && n <= len(sheets) {
			if size, ok := p.PartSize(sheets[n-1].part); ok {
				sheetBytes += min(size, uint64(2<<20)-sheetBytes)
			}
		}
	}
	batched := len(sel) > 1 && sheetBytes >= 2<<20
	var pending []pendingView
	var batchBytes uint64
	flush := func() {
		if len(pending) == 0 {
			return
		}
		c.encodePending(pending)
		pending = nil
		batchBytes = 0
	}
	for i, n := range sel {
		if n < 1 || n > len(sheets) {
			return nil, fmt.Errorf("xlsx: sheet %d out of range (1-%d)", n, len(sheets))
		}
		ref := sheets[n-1]
		id := "sheet" + strconv.Itoa(n)
		if ref.chart {
			v, layers := c.chartSheetSafe(ref, id)
			pending = append(pending, pendingView{view: v, layers: layers})
		} else {
			if views != ViewsPrint {
				v, tiles, err := c.worksheetSafe(ref, id)
				if err != nil {
					return nil, fmt.Errorf("xlsx: sheet %q: %w", ref.name, err)
				}
				pending = append(pending, pendingView{view: v, tiles: tiles})
			}
			if views != ViewsSheet {
				title := ref.name
				if views == ViewsBoth {
					title += " (print)"
				}
				v, pages, err := c.printSheetSafe(ref, "print"+strconv.Itoa(n), title)
				if err != nil {
					return nil, fmt.Errorf("xlsx: sheet %q: %w", ref.name, err)
				}
				pending = append(pending, pendingView{view: v, pages: pages})
			}
		}
		// Charts can read other sheets while this one is rendered. Once the
		// sheet is drawn, only the canvases are needed; release its cells and
		// any chart source sheets before the next sheet is read. (The sheets
		// of a filled template are kept: they are not read again.)
		if opts.fill == nil {
			clear(c.parsed)
		}
		if size, ok := p.PartSize(ref.part); ok {
			batchBytes += min(size, uint64(1<<20)-batchBytes)
		}
		if batched && batchBytes >= 1<<20 {
			flush()
			if i+1 < len(sel) {
				// The previous sheet's cells and drawing instructions have
				// become unreachable. Collect them before parsing the next one.
				runtime.GC()
			}
		}
	}
	// Rendering is complete. Font embedding and canvas encoding need none of
	// the workbook's parsed values, styles or package parts.
	c.pkg, c.r, c.d = nil, nil, nil
	c.sst, c.st, c.sheets, c.parsed = nil, nil, nil, nil
	flush()
	if !c.hasEmbedded {
		c.encodePending(nil)
	}
	c.fonts.ReportMissing()
	c.indexViews()
	return &Result{Doc: c.doc, Warnings: c.warnings, Sheets: len(sel), EmbeddedFonts: c.embeddedFonts}, nil
}

// newConverter sets up a conversion: the document, the fonts and the
// object builder. p is nil for grids, which come from no workbook.
func newConverter(p *ooxml.Package, opts *Options) *converter {
	c := &converter{pkg: p, opts: opts, doc: bdf.NewDocument(), warned: map[string]bool{}, patterns: map[string]bdf.Hash{},
		parsed: map[string]*worksheet{}, chartsFilled: map[string]bool{}, sheetWarned: map[string]bool{}}
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
	// pages are the objects of the layers of each page of a print view
	pages [][]*canvas.Canvas
}

// encodePending embeds any newly measured glyphs, encodes the canvases and
// points their views at the resulting object hashes.
func (c *converter) encodePending(views []pendingView) {
	opts := c.opts
	if !c.hasEmbedded || c.fonts.Generation() != c.embeddedGen {
		// With system fonts, the formula font is embedded still: its layout
		// depends on its glyphs.
		c.embeddedFonts = c.fonts.Embed(c.doc, fontset.EmbedOptions{NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2,
			IgnoreFSType: opts.IgnoreFSType, PinnedOnly: opts.SystemFonts})
		c.embeddedGen = c.fonts.Generation()
		c.hasEmbedded = true
	}
	c.cvs.EncodeAndRelease()
	for _, pv := range views {
		for key, cv := range pv.tiles {
			pv.view.Tiles[key] = cv.Hash().String()
		}
		if len(pv.layers) > 0 {
			for i, cv := range pv.layers {
				pv.view.Pages[0].Layers[i].Obj = cv.Hash()
			}
		}
		for i, layers := range pv.pages {
			for k, cv := range layers {
				pv.view.Pages[i].Layers[k].Obj = cv.Hash()
			}
		}
	}
}

// finish handles a single grid, which has no workbook or sheet batches.
func (c *converter) finish(views []pendingView) {
	c.encodePending(views)
	c.fonts.ReportMissing()
	c.indexViews()
}

func (c *converter) indexViews() {
	if c.opts.NoTextIndex {
		return
	}
	for _, v := range c.doc.Views {
		if _, err := c.doc.BuildTextIndex(v); err != nil {
			c.warnf("text index: %v", err)
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
	r, ok := c.pkg.RelOfType(c.wbPart, "/sharedStrings")
	if !ok {
		return
	}
	rc, err := c.pkg.OpenPart(r.Target)
	if err != nil {
		return
	}
	defer rc.Close()
	// the strings are read a tree at a time; a table that cannot be read
	// whole is not used at all, so that no cell shows another's string
	xr := ooxml.NewReader(rc)
	root, err := rootElement(xr)
	if err != nil {
		return
	}
	for {
		ok, err := xr.NextChild(root)
		if err != nil {
			c.sst = nil
			return
		}
		if !ok {
			return
		}
		if !xmlro.Equal(xr.LocalName(), "si") {
			continue
		}
		n, err := c.pkg.ReadFrom(xr)
		if err != nil {
			c.sst = nil
			return
		}
		c.sst = append(c.sst, c.readRich(n))
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

func attrs(kv ...string) []ooxml.Attr {
	var out []ooxml.Attr
	for i := 0; i+1 < len(kv); i += 2 {
		out = append(out, ooxml.Attr{Name: ooxml.Name{Local: kv[i]}, Value: kv[i+1]})
	}
	return out
}
