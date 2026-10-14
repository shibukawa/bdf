// Package pdfme fills templates of pdfme (https://pdfme.com) with data: a
// template made with the designer of pdfme (a JSON file: a base PDF or a
// blank page, and the fields on its pages with their places) is filled as
// the generator of pdfme fills it, into a BDF document.
//
// The fields have places in millimetres and a type: text (with its font,
// size, alignment and a size that fits), text with variables, a choice, a
// date, an image, an SVG, a signature, QR codes and barcodes, lines,
// rectangles and ellipses, checkboxes and radio buttons, lists, and tables,
// whose rows are the data: a table grows with its rows and moves what is
// under it, onto further pages of a blank base. The text of a read-only
// field may hold expressions ("{total * 1.1}", "{currentPage}").
//
// The data is an object: the value of each field by its name. It is what
// pdfme calls an input, with the values as the values of JSON and not
// their texts: a table takes a list of rows (lists of cells, or objects
// by the names of the columns), a checkbox true or false, a text with
// variables an object of them.
//
// What is not drawn as pdfme draws it: the inline markdown of text,
// text that makes its field taller (overflow "expand"), the codes that
// package barcode does not make (japanpost, upce, gs1datamatrix, pdf417),
// and the digits under barcodes. Text is laid out with the fonts of the
// machine by the names the template gives (pdfme has the fonts given to
// it); the default font of pdfme is Roboto.
package pdfme

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math"
	"strings"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/formbox"
	"github.com/shibukawa/bdf/converter/pdf"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/image/imgconv"
)

// TemplateData is what a template is filled with.
type TemplateData struct {
	// Data is the data: an object with the value of each field by its
	// name.
	Data map[string]any
	// Strict makes a field that the data has no value for an error;
	// without it the field is left empty, with a warning.
	Strict bool
	// Files holds the pictures that the data names by their paths.
	Files fs.FS
	// Format writes a number that fills a text (nil: as the General
	// format of a spreadsheet writes it).
	Format tag.Formatter
	// Now is the date and the time that the expressions date and dateTime
	// are ("2026/10/14 09:30"), as pdfme has those of the clock. Without it
	// they are values of the data, like the other names that expressions
	// read: the same data then fills to the same document every time.
	Now string
}

// Options controls the fonts of the text and the conversion of the base
// PDF.
type Options struct {
	Pages         conv.Pages
	Title         string
	FontFS        fs.FS
	FontDirs      []string
	NoSystemFonts bool
	SystemFonts   bool
	NoSubset      bool
	NoWOFF2       bool
	IgnoreFSType  bool
	Images        imgconv.Options
	NoTextIndex   bool
}

// OptionsOf returns the options that those of a conversion give.
func OptionsOf(o *conv.Options) *Options {
	if o == nil {
		return &Options{}
	}
	return &Options{Pages: o.Pages, Title: o.Title, FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts, SystemFonts: o.SystemFonts,
		NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, Images: o.Images, NoTextIndex: o.NoTextIndex}
}

// Result is a filled template.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	Fields   []tag.Field
	Pages    int
}

// Limits of a template (it is not trusted).
const (
	maxTemplate = 96 << 20
	maxPages    = 2000
	maxSchemas  = 1 << 16
	maxRows     = 1 << 17
	maxCells    = 1 << 20
	maxText     = 1 << 20
)

// mm is the points of a millimetre, as pdfme has them.
const mm = 2.8346

// schema is a field of a template: its properties as the template has
// them.
type schema struct {
	m          map[string]any
	name, typ  string
	x, y, w, h float64 // mm
	readOnly   bool
	// of a table or a list that is laid out over pages: the rows that
	// this part of it has, and whether it is one that goes on
	from, to int
	part     bool
	split    bool
	// total is the number of the lines of a text that is laid out over
	// pages (1 for one that is not in parts)
	total int
	// origin is the field of the template that a part is of (nil for the
	// field itself)
	origin *schema
}

// key is the field of the template that a field is, or is a part of.
func (s *schema) key() *schema {
	if s.origin != nil {
		return s.origin
	}
	return s
}

func (s *schema) num(key string, def float64) float64 {
	if v, ok := s.m[key].(float64); ok && !math.IsNaN(v) && !math.IsInf(v, 0) {
		return v
	}
	return def
}

func (s *schema) str(key, def string) string {
	if v, ok := s.m[key].(string); ok {
		return v
	}
	return def
}

func (s *schema) flag(key string) bool {
	v, _ := s.m[key].(bool)
	return v
}

func (s *schema) obj(key string) map[string]any {
	v, _ := s.m[key].(map[string]any)
	return v
}

// blankPage is a base that is a blank page: its size and the room at its
// edges, in millimetres.
type blankPage struct {
	w, h float64
	pad  [4]float64 // top, right, bottom, left
}

type template struct {
	pages  [][]*schema
	static []*schema
	blank  *blankPage
	base   []byte // a PDF
}

// IsTemplate reports whether a file looks like a template of pdfme by its
// start: an object of JSON with schemas or a base PDF.
func IsTemplate(head []byte) bool {
	t := bytes.TrimLeft(bytes.TrimPrefix(head, []byte("\xef\xbb\xbf")), " \t\r\n")
	return len(t) > 0 && t[0] == '{' && (bytes.Contains(t, []byte(`"schemas"`)) || bytes.Contains(t, []byte(`"basePdf"`)))
}

func newSchema(m map[string]any, name string) *schema {
	s := &schema{m: m, name: name}
	if n, ok := m["name"].(string); ok && name == "" {
		s.name = n
	}
	s.typ, _ = m["type"].(string)
	if pos, ok := m["position"].(map[string]any); ok {
		s.x, _ = pos["x"].(float64)
		s.y, _ = pos["y"].(float64)
	}
	s.w, s.h = s.num("width", 0), s.num("height", 0)
	s.readOnly = s.flag("readOnly")
	for _, v := range []*float64{&s.x, &s.y, &s.w, &s.h} {
		if math.IsNaN(*v) || math.IsInf(*v, 0) || math.Abs(*v) > 1e5 {
			*v = 0
		}
	}
	return s
}

// orderedPage reads a page of the older form of templates: an object of
// the fields by their names, in the order of the file.
func orderedPage(raw json.RawMessage) ([]*schema, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	var out []*schema
	for dec.More() {
		key, err := dec.Token()
		if err != nil {
			return nil, err
		}
		var m map[string]any
		if err := dec.Decode(&m); err != nil {
			return nil, fmt.Errorf("the field %v: %w", key, err)
		}
		name, _ := key.(string)
		if m != nil {
			out = append(out, newSchema(m, name))
		}
	}
	return out, nil
}

func schemaList(raw json.RawMessage) ([]*schema, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) > 0 && raw[0] == '{' {
		return orderedPage(raw)
	}
	var list []map[string]any
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}
	out := make([]*schema, 0, len(list))
	for _, m := range list {
		if m != nil {
			out = append(out, newSchema(m, ""))
		}
	}
	return out, nil
}

// parse reads a template.
func parse(r io.ReaderAt, size int64) (*template, error) {
	if size > maxTemplate {
		return nil, fmt.Errorf("pdfme: the template is larger than %d MiB", maxTemplate>>20)
	}
	data := make([]byte, size)
	if _, err := r.ReadAt(data, 0); err != nil && err != io.EOF {
		return nil, fmt.Errorf("pdfme: %w", err)
	}
	var top struct {
		Schemas []json.RawMessage `json:"schemas"`
		BasePdf json.RawMessage   `json:"basePdf"`
	}
	if err := json.Unmarshal(bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")), &top); err != nil {
		return nil, fmt.Errorf("pdfme: the template is not JSON: %w", err)
	}
	if top.Schemas == nil {
		return nil, fmt.Errorf("pdfme: the file is not a template of pdfme (it has no schemas)")
	}
	if len(top.Schemas) > maxPages {
		return nil, fmt.Errorf("pdfme: the template has more than %d pages", maxPages)
	}
	t := &template{}
	n := 0
	for i, raw := range top.Schemas {
		page, err := schemaList(raw)
		if err != nil {
			return nil, fmt.Errorf("pdfme: page %d of the template: %w", i+1, err)
		}
		if n += len(page); n > maxSchemas {
			return nil, fmt.Errorf("pdfme: the template has more than %d fields", maxSchemas)
		}
		t.pages = append(t.pages, page)
	}
	base := bytes.TrimSpace(top.BasePdf)
	switch {
	case len(base) > 0 && base[0] == '"':
		var s string
		if err := json.Unmarshal(base, &s); err != nil {
			return nil, fmt.Errorf("pdfme: the base PDF: %w", err)
		}
		i := strings.Index(s, "base64,")
		if !strings.HasPrefix(s, "data:") || i < 0 {
			return nil, fmt.Errorf("pdfme: the base PDF is not in the template (a data: URL of it is); the address of a file is not read")
		}
		pdfBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(s[i+len("base64,"):]))
		if err != nil {
			return nil, fmt.Errorf("pdfme: the base PDF is not base64: %w", err)
		}
		t.base = pdfBytes
	case len(base) > 0 && base[0] == '{':
		var b struct {
			Width, Height float64
			Padding       []float64
			StaticSchema  json.RawMessage `json:"staticSchema"`
		}
		if err := json.Unmarshal(base, &b); err != nil {
			return nil, fmt.Errorf("pdfme: the base page: %w", err)
		}
		if !(b.Width > 1) || !(b.Height > 1) || b.Width > 5000 || b.Height > 5000 {
			return nil, fmt.Errorf("pdfme: the base page is %g × %g mm", b.Width, b.Height)
		}
		t.blank = &blankPage{w: b.Width, h: b.Height}
		for i := 0; i < 4 && i < len(b.Padding); i++ {
			if p := b.Padding[i]; p > 0 && p < b.Height {
				t.blank.pad[i] = p
			}
		}
		if len(bytes.TrimSpace(b.StaticSchema)) > 0 {
			st, err := schemaList(b.StaticSchema)
			if err != nil {
				return nil, fmt.Errorf("pdfme: the static fields of the base page: %w", err)
			}
			t.static = st
		}
	default:
		return nil, fmt.Errorf("pdfme: the template has no base PDF")
	}
	return t, nil
}

// Barcodes of pdfme, by the symbologies of package barcode ("" for one
// that is not made).
var codeKinds = map[string]string{
	"qrcode": "qr", "ean13": "ean13", "ean8": "ean8", "code39": "code39", "code128": "code128", "nw7": "nw7", "itf14": "itf14", "upca": "upca",
	"japanpost": "", "upce": "", "gs1datamatrix": "", "pdf417": "",
}

// Kinds of fields, by what they are filled with.
func isText(typ string) bool {
	switch typ {
	case "text", "select", "date", "time", "dateTime":
		return true
	}
	return false
}

func isMark(typ string) bool { return typ == "checkbox" || typ == "radioGroup" || typ == "circleMark" }

// heads are the names of the columns of a table.
func (s *schema) heads() []string {
	var out []string
	if l, ok := s.m["head"].([]any); ok {
		for _, h := range l {
			out = append(out, tag.ValueText(tag.Normalize(h), "", nil))
		}
	}
	return out
}

// variables are the names of the variables of a text with variables.
func (s *schema) variables() []string {
	var out []string
	if l, ok := s.m["variables"].([]any); ok {
		for _, v := range l {
			if name, ok := v.(string); ok && name != "" {
				out = append(out, name)
			}
		}
	}
	return out
}

// fields tells the fields of a template as the names that data has.
func (t *template) fields() []tag.Field {
	var out []tag.Field
	index := map[string]int{}
	add := func(f tag.Field, at string) {
		key := strings.Join(f.In, "\x00") + "\x00\x00" + f.Name
		if i, ok := index[key]; ok {
			if o := &out[i]; len(o.At) < 16 && o.At[len(o.At)-1] != at {
				o.At = append(o.At, at)
			}
			return
		}
		index[key] = len(out)
		f.At = []string{at}
		out = append(out, f)
	}
	for pi, page := range t.pages {
		at := fmt.Sprintf("page %d", pi+1)
		for _, s := range page {
			if s.readOnly || s.name == "" {
				continue
			}
			content := s.str("content", "")
			f := tag.Field{Name: s.name, Kind: tag.FieldValue}
			switch kind, isCode := codeKinds[s.typ]; {
			case isCode:
				f.Kind, f.Code = tag.FieldImage, kind
				if kind == "" {
					f.Kind, f.Type = tag.FieldValue, "string"
				}
				if content != "" {
					f.Default = content
				}
			case s.typ == "image" || s.typ == "signature" || s.typ == "svg":
				f.Kind = tag.FieldImage
			case isMark(s.typ):
				f.Type, f.Default = "boolean", content == "true"
			case s.typ == "table":
				sec := f
				sec.Kind, sec.Mode = tag.FieldSection, tag.ModeRows
				add(sec, at)
				for _, h := range s.heads() {
					add(tag.Field{Name: h, Kind: tag.FieldValue, In: []string{s.name}}, at) // (a text or a number)
				}
				continue
			case s.typ == "list":
				sec := f
				sec.Kind = tag.FieldSection
				add(sec, at)
				add(tag.Field{Name: ".", Kind: tag.FieldValue, In: []string{s.name}}, at)
				continue
			case s.typ == "multiVariableText":
				for _, v := range s.variables() {
					add(tag.Field{Name: s.name + "." + v, Kind: tag.FieldValue, Type: "string"}, at)
				}
				continue
			case isText(s.typ):
				f.Type = "string"
				switch s.typ {
				case "date":
					f.Type = "date"
				case "dateTime":
					f.Type = "date-time"
				case "select":
					if l, ok := s.m["options"].([]any); ok {
						for _, o := range l {
							f.Options = append(f.Options, tag.ValueText(tag.Normalize(o), "", nil))
						}
					}
				}
				f.Multiline = s.typ == "text" && (strings.Contains(content, "\n") || s.h > 2.2*s.num("fontSize", 13)*0.3528)
				if content != "" && s.typ != "date" && s.typ != "dateTime" && s.typ != "time" {
					f.Default = content
				}
			default:
				continue // (a shape, or a type that is not known: nothing to fill)
			}
			add(f, at)
		}
	}
	for _, e := range t.expressionNames() {
		add(tag.Field{Name: e.name, Kind: tag.FieldValue}, e.at)
	}
	return out
}

// exprName is a name that an expression of a text of the template reads
// of the data, and where.
type exprName struct{ name, at string }

// expressionNames returns the names that the expressions of the read-only
// texts read and no field of the template has: values of the data that the
// template shows in its texts ("Total {total}"). The date and the time of
// pdfme (date, dateTime) are among them: the data says them, so that the
// same data fills to the same form.
func (t *template) expressionNames() []exprName {
	known := map[string]bool{"currentPage": true, "totalPages": true}
	for _, list := range append(append([][]*schema{}, t.pages...), t.static) {
		for _, s := range list {
			known[s.name] = true
		}
	}
	var out []exprName
	scan := func(list []*schema, at string) {
		for _, s := range list {
			if !s.readOnly || s.typ == "table" || s.typ == "multiVariableText" {
				continue
			}
			for _, name := range exprNames(s.str("content", "")) {
				if !known[name] {
					known[name] = true
					out = append(out, exprName{name, at})
				}
			}
		}
	}
	for i, page := range t.pages {
		scan(page, fmt.Sprintf("page %d", i+1))
	}
	scan(t.static, "every page")
	return out
}

// Inspect reads the fields of a template: the names that data for it has.
func Inspect(r io.ReaderAt, size int64) ([]tag.Field, error) {
	t, err := parse(r, size)
	if err != nil {
		return nil, err
	}
	return t.fields(), nil
}

// Defaults returns the data that a template has for itself: the content
// that its designer left in each field, as the value of the field. It is a
// sample of the data of the template.
func Defaults(r io.ReaderAt, size int64) (map[string]any, error) {
	t, err := parse(r, size)
	if err != nil {
		return nil, err
	}
	out := map[string]any{}
	for _, page := range t.pages {
		for _, s := range page {
			if s.readOnly || s.name == "" {
				continue
			}
			if _, ok := out[s.name]; ok {
				continue
			}
			content := s.str("content", "")
			switch {
			case isMark(s.typ):
				out[s.name] = content == "true"
			case s.typ == "table":
				rows := []any{}
				heads := s.heads()
				for _, row := range tableRows(content, nil) {
					obj := map[string]any{}
					for i, h := range heads {
						if i < len(row) {
							obj[h] = row[i]
						} else {
							obj[h] = ""
						}
					}
					rows = append(rows, obj)
				}
				out[s.name] = rows
			case s.typ == "list":
				items := []any{}
				for _, it := range listItems(content) {
					items = append(items, it)
				}
				out[s.name] = items
			case s.typ == "multiVariableText":
				vars := map[string]any{}
				var m map[string]any
				json.Unmarshal([]byte(content), &m)
				for _, v := range s.variables() {
					vars[v] = tag.ValueText(tag.Normalize(m[v]), "", nil)
				}
				out[s.name] = vars
			case s.typ == "image" || s.typ == "signature" || s.typ == "svg":
				// (a picture of the template is in the template: it is no
				// value of data)
			case isText(s.typ) || codeKinds[s.typ] != "" || s.typ == "japanpost" || s.typ == "upce" || s.typ == "gs1datamatrix" || s.typ == "pdf417":
				out[s.name] = content
			}
		}
	}
	for _, e := range t.expressionNames() {
		if _, ok := out[e.name]; !ok {
			out[e.name] = ""
		}
	}
	return out, nil
}

// tableRows reads the rows of a table: of a value of data (a list of
// lists of cells, or of objects by the names of the columns), or of the
// text that pdfme has them as (the JSON of a list of lists).
func tableRows(v any, heads []string) [][]string {
	if s, ok := v.(string); ok {
		var parsed any
		if json.Unmarshal([]byte(s), &parsed) != nil {
			return nil
		}
		v = tag.Normalize(parsed)
	}
	list, _ := v.([]any)
	if len(list) > maxRows {
		list = list[:maxRows]
	}
	out := make([][]string, 0, len(list))
	cells := 0
	for _, r := range list {
		var row []string
		switch x := r.(type) {
		case []any:
			for _, c := range x {
				row = append(row, tag.ValueText(c, "", nil))
			}
		case map[string]any:
			for _, h := range heads {
				v, ok := x[h]
				if !ok && strings.Contains(h, ".") {
					// (a name with dots is an object in an object)
					var at any = x
					for _, part := range strings.Split(h, ".") {
						m, _ := at.(map[string]any)
						at = m[part]
					}
					v = at
				}
				row = append(row, tag.ValueText(v, "", nil))
			}
		default:
			continue
		}
		if cells += len(row); cells > maxCells {
			break
		}
		out = append(out, row)
	}
	return out
}

// listItems reads the items of a list: of a list of data, of the JSON of
// one, or of the lines of a text.
func listItems(v any) []string {
	if s, ok := v.(string); ok {
		var parsed []any
		if json.Unmarshal([]byte(s), &parsed) != nil {
			if strings.TrimSpace(s) == "" {
				return nil
			}
			return strings.Split(strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(s), "\n")
		}
		v = tag.Normalize(parsed)
	}
	list, _ := v.([]any)
	if len(list) > maxRows {
		list = list[:maxRows]
	}
	out := make([]string, 0, len(list))
	for _, it := range list {
		out = append(out, tag.ValueText(it, "", nil))
	}
	return out
}

// filler fills a template.
type filler struct {
	t        *template
	td       *TemplateData
	data     map[string]any
	r        *formbox.Renderer
	warnings []string
	warned   map[string]bool
	// vars are the values that expressions have, and pages the number of
	// the pages of the filled template
	vars  map[string]any
	pages int
	// measuring is set while the pages are laid out: the number of a page
	// is then no value of an expression yet
	measuring bool
	// rows and items are the rows of the tables and the items of the lists
	// that were read, by their fields
	rows  map[*schema][][]string
	items map[*schema][]string
	// pictures are the pictures that the data names, read once each: a
	// picture of the base page is on every page
	pictures map[string]loaded
}

// rowsOf reads the rows of a table once.
func (f *filler) rowsOf(s *schema, v any) [][]string {
	if rows, ok := f.rows[s.key()]; ok {
		return rows
	}
	rows := tableRows(v, s.heads())
	if f.rows == nil {
		f.rows = map[*schema][][]string{}
	}
	f.rows[s.key()] = rows
	return rows
}

// itemsOf reads the items of a list once.
func (f *filler) itemsOf(s *schema, v any) []string {
	if items, ok := f.items[s.key()]; ok {
		return items
	}
	items := listItems(v)
	if f.items == nil {
		f.items = map[*schema][]string{}
	}
	f.items[s.key()] = items
	return items
}

func (f *filler) warnf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	if !f.warned[s] {
		f.warned[s] = true
		f.warnings = append(f.warnings, s)
	}
}

// value finds the value of a field in the data.
func (f *filler) value(s *schema) (any, bool) {
	if s.typ == "multiVariableText" {
		if v, ok := f.data[s.name]; ok {
			return v, true
		}
		// (its variables may be values of the data by their own names)
		vars := map[string]any{}
		for _, name := range s.variables() {
			if v, ok := f.data[name]; ok {
				vars[name] = v
			}
		}
		return vars, len(vars) > 0
	}
	v, ok := f.data[s.name]
	return v, ok
}

// check tells the fields that the data has no value for.
func (f *filler) check() error {
	seen := map[string]bool{}
	for _, page := range f.t.pages {
		for _, s := range page {
			if s.readOnly || s.name == "" || seen[s.name] {
				continue
			}
			seen[s.name] = true
			switch {
			case isText(s.typ), isMark(s.typ), s.typ == "table", s.typ == "list", s.typ == "multiVariableText", s.typ == "image", s.typ == "signature", s.typ == "svg":
			default:
				if _, ok := codeKinds[s.typ]; !ok {
					continue
				}
			}
			if _, ok := f.value(s); ok {
				continue
			}
			if f.td.Strict {
				return fmt.Errorf("pdfme: the data has no %q", s.name)
			}
			f.warnf("the data has no %q; it is left empty", s.name)
		}
	}
	for _, e := range f.t.expressionNames() {
		if _, ok := f.data[e.name]; ok || f.td.Now != "" && (e.name == "date" || e.name == "dateTime") {
			continue
		}
		if f.td.Strict {
			return fmt.Errorf("pdfme: the data has no %q", e.name)
		}
		f.warnf("the data has no %q, which a text of the template shows", e.name)
	}
	return nil
}

// expressions sets the values that expressions have: the content of the
// read-only fields and the values of the data, as pdfme gives them (a text
// that is the JSON of a value is that value).
func (f *filler) expressions() {
	vars := map[string]any{}
	parsed := func(v any) any {
		if s, ok := v.(string); ok {
			var x any
			if json.Unmarshal([]byte(s), &x) == nil {
				return tag.Normalize(x)
			}
		}
		return v
	}
	for _, list := range append(append([][]*schema{}, f.t.pages...), f.t.static) {
		for _, s := range list {
			if s.name == "" {
				continue
			}
			if s.readOnly {
				vars[s.name] = parsed(s.str("content", ""))
			} else if _, ok := vars[s.name]; !ok {
				vars[s.name] = ""
			}
		}
	}
	for k, v := range f.data {
		vars[k] = parsed(v)
	}
	if now := f.td.Now; now != "" {
		vars["dateTime"] = now
		vars["date"] = strings.SplitN(now, " ", 2)[0]
	}
	f.vars = vars
}

// Fill fills a template with data.
func Fill(r io.ReaderAt, size int64, td *TemplateData, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if td == nil {
		td = &TemplateData{}
	}
	t, err := parse(r, size)
	if err != nil {
		return nil, err
	}
	f := &filler{t: t, td: td, warned: map[string]bool{}}
	f.data, _ = tag.Normalize(td.Data).(map[string]any)
	if f.data == nil {
		f.data = map[string]any{}
	}
	if err := f.check(); err != nil {
		return nil, err
	}
	f.expressions()

	var doc *bdf.Document
	var view *bdf.View
	fopts := formbox.Options{FontFS: opts.FontFS, FontDirs: opts.FontDirs, NoSystemFonts: opts.NoSystemFonts, SystemFonts: opts.SystemFonts,
		NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType, Images: opts.Images,
		Warn: func(msg string) { f.warnf("%s", msg) }}
	type placed struct {
		page    *bdf.Page
		schemas []*schema
	}
	var pages []placed
	if t.blank != nil {
		doc = bdf.NewDocument()
		doc.Meta.Source = "pdfme"
		view = doc.NewView("pages", bdf.ViewFixed, opts.Title)
		f.r = formbox.New(doc, fopts)
		f.measuring = true
		laid := f.layout()
		f.measuring = false
		if len(laid) > maxPages {
			f.warnf("the filled template has more than %d pages: the rest are left out", maxPages)
			laid = laid[:maxPages]
		}
		for _, list := range laid {
			pg := view.AddPage(float32(t.blank.w*mm), float32(t.blank.h*mm))
			pages = append(pages, placed{pg, list})
		}
		if sel := opts.Pages.Numbers(len(pages)); sel != nil {
			var kept []placed
			var keptPages []*bdf.Page
			for _, n := range sel {
				if n >= 1 && n <= len(pages) {
					kept, keptPages = append(kept, pages[n-1]), append(keptPages, pages[n-1].page)
				}
			}
			pages, view.Pages = kept, keptPages
		}
	} else {
		var pw []string
		res, err := pdf.Convert(bytes.NewReader(t.base), &pdf.Options{Pages: opts.Pages, Title: opts.Title, NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2,
			IgnoreFSType: opts.IgnoreFSType, Images: opts.Images, NoTextIndex: true, Warn: func(msg string) { pw = append(pw, msg) }})
		if err != nil {
			return nil, fmt.Errorf("pdfme: the base PDF: %w", err)
		}
		for _, w := range pw {
			f.warnf("the base PDF: %s", w)
		}
		doc = res.Doc
		doc.Meta.Source = "pdfme"
		if len(doc.Views) == 0 {
			return nil, fmt.Errorf("pdfme: the base PDF has no pages")
		}
		view = doc.Views[0]
		f.r = formbox.New(doc, fopts)
		// the fields of a page of the template are on the page of the base
		// of its number
		for i, pg := range view.Pages {
			var list []*schema
			if i < len(res.PageNumbers) {
				if n := res.PageNumbers[i] - 1; n >= 0 && n < len(t.pages) {
					list = t.pages[n]
				}
			}
			pages = append(pages, placed{pg, list})
		}
	}
	if opts.Title != "" {
		doc.Meta.DC.Title = bdf.DCValues{opts.Title}
	}
	f.pages = len(pages)
	var layers []*formbox.Layer
	for i, p := range pages {
		l := f.r.Layer(float64(p.page.W), float64(p.page.H))
		f.vars["currentPage"], f.vars["totalPages"] = float64(i+1), float64(len(pages))
		for _, s := range t.static {
			f.draw(l, s, true)
		}
		for _, s := range p.schemas {
			f.draw(l, s, false)
		}
		layers = append(layers, l)
	}
	f.r.Finish()
	for i, p := range pages {
		if layers[i].Drawn() || t.blank != nil {
			p.page.Layers = append(p.page.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: layers[i].Hash()})
		}
	}
	if !opts.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			f.warnf("text index: %v", err)
		}
	}
	return &Result{Doc: doc, Warnings: f.warnings, Fields: t.fields(), Pages: len(pages)}, nil
}
