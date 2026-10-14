// Package pdfform fills the interactive form of a PDF (its AcroForm) with
// data, as a template of a form: the file is the template, its fields are
// what the data has values for, and the filled form is a BDF document of
// its pages with the values drawn in the places of the fields.
//
// The name of a field is the name the data has its value by: the full name
// of the field, whose parts are the names of an object in an object
// ("applicant.name"). A text field takes a text or a number, a checkbox
// true or false, a group of radio buttons the value of one of its buttons,
// and a choice field one of its choices. A field that the data has no
// value for is left as the file has it.
//
// The values are drawn here, not by the file: a text in the font, the size,
// the color and the alignment that the field says (in a font of the
// machine that has its characters, which the fonts of a form often lack),
// and the buttons with the appearances that the file has for their states.
// The filled form is flat: its fields are no fields any more.
package pdfform

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"sort"
	"strings"

	"github.com/shibukawa/bdf"
	conv "github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/internal/formbox"
	"github.com/shibukawa/bdf/converter/pdf"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/image/imgconv"
)

// TemplateData is what a form is filled with.
type TemplateData struct {
	// Data is the data: an object of the values of JSON.
	Data map[string]any
	// Strict makes a field that the data has no value for an error;
	// without it the field is left as the file has it, with a warning.
	Strict bool
	// Format writes a number that fills a text field (nil: as the General
	// format of a spreadsheet writes it).
	Format tag.Formatter
}

// Options controls the conversion of the file and the fonts of the values.
type Options struct {
	// Pages selects the pages of the file (nil: all), and Password opens
	// a file that is encrypted.
	Pages    conv.Pages
	Password string
	Title    string
	// The fonts that the values are drawn with (see converter.Options).
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
	return &Options{Pages: o.Pages, Password: o.Password, Title: o.Title, FontFS: o.FontFS, FontDirs: o.FontDirs, NoSystemFonts: o.NoSystemFonts,
		SystemFonts: o.SystemFonts, NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType, Images: o.Images, NoTextIndex: o.NoTextIndex}
}

// Result is a filled form.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	// Fields are the fields of the form: the names that the data has
	// values for.
	Fields []tag.Field
	Pages  int
	// Protected reports that the file opened only with the password.
	Protected bool
}

func (o *Options) pdf() *pdf.Options {
	return &pdf.Options{Pages: o.Pages, Password: o.Password, Title: o.Title, NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2,
		IgnoreFSType: o.IgnoreFSType, Images: o.Images, NoTextIndex: true}
}

// ErrNoForm is returned for a PDF without the fields of a form.
var ErrNoForm = fmt.Errorf("pdfform: the PDF has no form fields to fill")

// Inspect reads the fields of the form of a PDF: the names that data for
// it has.
func Inspect(rs io.ReadSeeker, opts *Options) ([]tag.Field, error) {
	if opts == nil {
		opts = &Options{}
	}
	fields, err := pdf.FormFields(rs, opts.pdf())
	if err != nil {
		return nil, err
	}
	out := fieldList(fields)
	if len(out) == 0 {
		return nil, ErrNoForm
	}
	return out, nil
}

// fieldList tells the fields of a form as the names of a template.
func fieldList(fields []pdf.FormField) []tag.Field {
	var out []tag.Field
	for i := range fields {
		fd := &fields[i]
		f := tag.Field{Name: fd.Name, Kind: tag.FieldValue, Label: fd.Label}
		switch fd.Kind {
		case pdf.FormText:
			f.Type, f.MaxLen, f.Multiline = "string", fd.MaxLen, fd.Multiline
			if fd.Value != "" {
				f.Default = fd.Value
			}
		case pdf.FormCheckbox:
			f.Type, f.Default = "boolean", fd.Value != "" && fd.Value != "Off"
		case pdf.FormRadio, pdf.FormChoice:
			f.Type = "string"
			labelled := false
			for _, o := range fd.Options {
				f.Options, f.Labels = append(f.Options, o.Value), append(f.Labels, o.Label)
				labelled = labelled || o.Label != o.Value
			}
			if !labelled {
				f.Labels = nil
			}
			if fd.Value != "" && fd.Value != "Off" {
				f.Default = fd.Value
			}
		default:
			continue // (a button to press, a signature: nothing that data fills)
		}
		seen := map[int]bool{}
		for _, w := range fd.Widgets {
			if !seen[w.Page] {
				seen[w.Page] = true
				f.At = append(f.At, fmt.Sprintf("page %d", w.Page+1))
			}
		}
		out = append(out, f)
	}
	return out
}

// lookup finds the value of a field in the data: by its full name, or as
// the names of objects in objects that the parts of its name are.
func lookup(data map[string]any, name string) (any, bool) {
	if v, ok := data[name]; ok {
		return v, true
	}
	var v any = data
	for _, part := range strings.Split(name, ".") {
		m, ok := v.(map[string]any)
		if !ok {
			return nil, false
		}
		if v, ok = m[part]; !ok {
			return nil, false
		}
	}
	return v, true
}

// Fill fills the form of a PDF with data.
func Fill(rs io.ReadSeeker, t *TemplateData, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if t == nil {
		t = &TemplateData{}
	}
	fields, err := pdf.FormFields(rs, opts.pdf())
	if err != nil {
		return nil, err
	}
	list := fieldList(fields)
	if len(list) == 0 {
		return nil, ErrNoForm
	}
	data, _ := tag.Normalize(t.Data).(map[string]any)
	var warnings []string
	warn := func(format string, args ...any) { warnings = append(warnings, fmt.Sprintf(format, args...)) }

	// what each field shows: the state of a button, the text of a field
	show := map[string]pdf.FormShow{}
	texts := map[string]string{}
	for i := range fields {
		fd := &fields[i]
		if fd.Kind == pdf.FormButton || fd.Kind == pdf.FormSignature {
			continue
		}
		v, ok := lookup(data, fd.Name)
		if !ok {
			if t.Strict {
				return nil, fmt.Errorf("pdfform: the data has no %q", fd.Name)
			}
			warn("the data has no %q; the field is left as the file has it", fd.Name)
			continue
		}
		switch fd.Kind {
		case pdf.FormCheckbox:
			state := "Off"
			if s, isText := v.(string); isText && s != "" {
				// (the state of the box by its name, as a group of boxes has them)
				state = s
				switch strings.ToLower(s) {
				case "off", "false", "no", "0":
					state = "Off"
				case "on", "true", "yes", "1":
					state = onState(fd)
				}
			} else if tag.Truthy(v) {
				state = onState(fd)
			}
			show[fd.Name] = pdf.FormShow{State: state}
		case pdf.FormRadio:
			state := "Off"
			if s := tag.ValueText(v, "", t.Format); s != "" {
				known := false
				for _, o := range fd.Options {
					known = known || o.Value == s
				}
				if known {
					state = s
				} else if t.Strict {
					return nil, fmt.Errorf("pdfform: %q is not a button of %q (%s)", s, fd.Name, optionList(fd))
				} else {
					warn("%q is not a button of %q (%s); none is on", s, fd.Name, optionList(fd))
				}
			}
			show[fd.Name] = pdf.FormShow{State: state}
		case pdf.FormChoice:
			s := tag.ValueText(v, "", t.Format)
			for _, o := range fd.Options {
				if o.Value == s {
					s = o.Label // what the field shows of the choice
					break
				}
			}
			show[fd.Name], texts[fd.Name] = pdf.FormShow{Blank: true}, s
		default:
			s := tag.ValueText(v, "", t.Format)
			if fd.Password {
				s = strings.Repeat("*", len([]rune(s)))
			}
			show[fd.Name], texts[fd.Name] = pdf.FormShow{Blank: true}, s
		}
	}

	po := opts.pdf()
	po.Form = show
	po.Warn = func(msg string) { warnings = append(warnings, msg) }
	res, err := pdf.Convert(rs, po)
	if err != nil {
		return nil, err
	}
	doc := res.Doc
	if len(doc.Views) == 0 {
		return nil, fmt.Errorf("pdfform: the PDF has no pages")
	}
	view := doc.Views[0]

	// the texts, over the pages
	r := formbox.New(doc, formbox.Options{FontFS: opts.FontFS, FontDirs: opts.FontDirs, NoSystemFonts: opts.NoSystemFonts, SystemFonts: opts.SystemFonts,
		NoSubset: opts.NoSubset, NoWOFF2: opts.NoWOFF2, IgnoreFSType: opts.IgnoreFSType, Images: opts.Images,
		Warn: func(msg string) { warnings = append(warnings, msg) }})
	layers := map[int]*formbox.Layer{}
	for i := range fields {
		fd := &fields[i]
		text, ok := texts[fd.Name]
		if !ok || text == "" {
			continue
		}
		for wi := range fd.Widgets {
			w := &fd.Widgets[wi]
			if w.Page < 0 || w.Page >= len(view.Pages) {
				continue
			}
			l := layers[w.Page]
			if l == nil {
				pg := view.Pages[w.Page]
				l = r.Layer(float64(pg.W), float64(pg.H))
				layers[w.Page] = l
			}
			drawText(r, l, fd, w, text)
		}
	}
	r.Finish()
	pages := make([]int, 0, len(layers))
	for p := range layers {
		pages = append(pages, p)
	}
	sort.Ints(pages)
	for _, p := range pages {
		if l := layers[p]; l.Drawn() {
			view.Pages[p].Layers = append(view.Pages[p].Layers, bdf.Layer{Role: bdf.RoleBody, Obj: l.Hash()})
		}
	}
	if !opts.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			warn("text index: %v", err)
		}
	}
	return &Result{Doc: doc, Warnings: warnings, Fields: list, Pages: len(view.Pages), Protected: res.Protected}, nil
}

// onState is the state that the boxes of a checkbox show when it is on.
func onState(fd *pdf.FormField) string {
	for _, w := range fd.Widgets {
		if w.State != "" {
			return w.State
		}
	}
	return "Yes"
}

func optionList(fd *pdf.FormField) string {
	var names []string
	for _, o := range fd.Options {
		names = append(names, o.Value)
	}
	return strings.Join(names, ", ")
}

// fontOf reads the name of a font of a form as a family and its style:
// "Helvetica-BoldOblique" is Helvetica, bold and slanted.
func fontOf(name string) (family string, bold, italic bool) {
	family = name
	if i := strings.IndexAny(name, "-,"); i > 0 {
		family = name[:i]
	}
	style := strings.ToLower(name[len(family):])
	bold = strings.Contains(style, "bold") || strings.Contains(style, "black") || strings.Contains(style, "heavy")
	italic = strings.Contains(style, "italic") || strings.Contains(style, "oblique")
	switch strings.ToLower(family) {
	case "helv", "helvetica", "arial", "arialmt":
		family = "Helvetica"
	case "tiro", "times", "timesnewroman", "timesnewromanpsmt":
		family = "Times New Roman"
	case "cour", "courier", "couriernew", "couriernewpsmt":
		family = "Courier New"
	case "":
		family = "Helvetica"
	}
	return family, bold, italic
}

// The room of the text of a field, as readers of PDF files leave it: at
// the sides of the text, and the height of a line of a size.
const (
	sidePad    = 2.0
	lineFactor = 1.35 // of a field of one line: the line fills the field
	leading    = 1.2  // of a field of several lines
	minSize    = 4.0
	autoMulti  = 12.0 // the size that a field of several lines starts from
)

// drawText draws the text of a field in a widget: by the rules that
// readers of PDF files make the appearance of a field by.
func drawText(r *formbox.Renderer, l *formbox.Layer, fd *pdf.FormField, w *pdf.FormWidget, text string) {
	st := &formbox.Style{Size: fd.FontSize, Color: fd.Color}
	st.Font, st.Bold, st.Italic = fontOf(fd.FontName)
	b := math.Max(w.Border, 0)
	pad := sidePad + b
	room := w.W - 2*pad
	if room <= 0 || w.H <= 2*b {
		return
	}
	o := l.Obj
	o.Save()
	l.Transform(w.Matrix)
	o.ClipRect(float32(b), float32(b), float32(w.W-2*b), float32(w.H-2*b))
	l.Mark(bdf.MarkParagraph, "")
	defer o.Restore()

	if fd.Multiline {
		// lines from the top; a size that fits is the largest that keeps
		// them in the field
		if st.Size <= 0 {
			for st.Size = autoMulti; st.Size > minSize; st.Size -= 0.5 {
				if float64(len(r.Wrap(text, st, room)))*leading*st.Size <= w.H-2*b-2*sidePad {
					break
				}
			}
		}
		for i, line := range r.Wrap(text, st, room) {
			if i > 0 {
				l.Mark(bdf.MarkLine, "")
			}
			y := b + sidePad + st.Size*(0.95+leading*float64(i))
			l.Text(line, st, pad+align(fd.Align, room, r.Width(line, st)), y)
		}
		return
	}
	text = strings.NewReplacer("\r\n", " ", "\n", " ", "\r", " ").Replace(text)
	if st.Size <= 0 {
		// as tall as the field has room for, and no wider than it
		st.Size = math.Max((w.H-2*b)/lineFactor, minSize)
		if tw := r.Width(text, st); tw > room && !fd.Comb {
			st.Size = math.Max(st.Size*room/tw, minSize)
		}
	}
	y := w.H/2 + 0.325*st.Size
	if fd.Comb && fd.MaxLen > 0 {
		// a cell for each character
		cell := w.W / float64(fd.MaxLen)
		for i, c := range []rune(text) {
			if i >= fd.MaxLen {
				break
			}
			s := string(c)
			l.Text(s, st, cell*float64(i)+(cell-r.Width(s, st))/2, y)
		}
		return
	}
	l.Text(text, st, pad+align(fd.Align, room, r.Width(text, st)), y)
}

// align is where a line starts in the room it has.
func align(q int, room, width float64) float64 {
	switch q {
	case 1:
		return (room - width) / 2
	case 2:
		return room - width
	}
	return 0
}
