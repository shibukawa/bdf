package pdf

import (
	"io"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"

	"github.com/shibukawa/bdf"
)

// The fields of the interactive form of a PDF (its AcroForm): the places
// that a reader types values into. A file with them is a form to fill:
// FormFields tells what the fields are and where their widgets are on the
// pages, and Options.Form says what the widgets show when the file is
// converted, so that a caller draws the values of data in their places
// (converter/pdfform does).

// Kinds of form fields.
const (
	FormText      = "text"
	FormCheckbox  = "checkbox"
	FormRadio     = "radio"
	FormChoice    = "choice"
	FormButton    = "button"
	FormSignature = "signature"
)

// FormField is a field of the form of a PDF.
type FormField struct {
	// Name is the full name of the field: the names of the fields it is
	// in and its own, joined by dots ("customer.name").
	Name string
	// Label is the name of the field for people (its tooltip), when it
	// has one.
	Label string
	Kind  string
	// Value and Default are the value the file has for the field and the
	// one it is reset to: the text of a text field, the state of a
	// checkbox or of a group of radio buttons ("Off" for none), the
	// choice of a choice field.
	Value, Default string
	// Of a text field: more than one line, a cell for each character of
	// MaxLen, and a password.
	Multiline, Comb, Password bool
	MaxLen                    int
	ReadOnly, Required        bool
	// Options are the choices of a choice field, and the states of the
	// buttons of a group of radio buttons.
	Options []FormOption
	// The text of the field, as its default appearance says: the font by
	// its name (""), the size in points (0: as large as the field has room
	// for), the color, and the alignment (0 left, 1 centered, 2 right).
	FontName string
	FontSize float64
	Color    bdf.Color
	Align    int
	Widgets  []FormWidget
}

// FormOption is a choice of a field: the value that the data has for it,
// and what a reader sees.
type FormOption struct {
	Value, Label string
}

// FormWidget is a place of a field on a page.
type FormWidget struct {
	// Page is the index of the page among those of the document that the
	// file converts to, from 0.
	Page int
	// W and H are the size of what the widget holds, upright, and Matrix
	// maps that space (its origin at the top left, y down) to the space
	// of the page: a widget, or its page, may be turned.
	W, H   float64
	Matrix [6]float64
	// State is the state that a checkbox or a radio button shows when it
	// is on ("Yes", or the value of the button).
	State string
	// Border is the width of the border that the widget draws.
	Border float64
}

// FormShow says what the widgets of a field show when the file is
// converted, instead of what the file has for them.
type FormShow struct {
	// State is the state of a checkbox or of a group of radio buttons:
	// a button shows its appearance of that state when it has one, and
	// that of "Off" when not.
	State string
	// Blank leaves the text of a text or choice field out, and keeps
	// its border and background: the caller draws the value.
	Blank bool
}

const (
	maxFieldDepth = 32
	maxFormFields = 1 << 16
)

// field is what the dictionaries of a widget and of the fields above it
// say of a field.
type fieldInfo struct {
	name string
	attr func(key string) types.Object
}

// fieldOf reads the field of a widget annotation: its full name, and its
// attributes, which the fields above it give it.
func (p *pdf) fieldOf(ad types.Dict) fieldInfo {
	chain := []types.Dict{ad}
	for d, i := p.dict(ad["Parent"]), 0; d != nil && i < maxFieldDepth; d, i = p.dict(d["Parent"]), i+1 {
		chain = append(chain, d)
	}
	var names []string
	for i := len(chain) - 1; i >= 0; i-- {
		if t, ok := chain[i]["T"]; ok {
			names = append(names, p.text(t))
		}
	}
	return fieldInfo{name: strings.Join(names, "."), attr: func(key string) types.Object {
		for _, d := range chain {
			if v, ok := d[key]; ok {
				return v
			}
		}
		return nil
	}}
}

// onState returns the state that the appearance of a button has beside
// "Off".
func (p *pdf) onState(ad types.Dict) string {
	nd := p.dict(p.dict(ad["AP"])["N"])
	state := ""
	for k := range nd {
		if k != "Off" && (state == "" || k < state) {
			state = k
		}
	}
	return state
}

// valueText reads the value of a field as text: a string, a name, or the
// first of an array of them.
func (p *pdf) valueText(o types.Object) string {
	switch v := p.deref(o).(type) {
	case types.Name:
		return v.Value()
	case types.Array:
		if len(v) > 0 {
			return p.valueText(v[0])
		}
		return ""
	case nil:
		return ""
	}
	return p.text(o)
}

// defaultAppearance reads a default appearance string ("/Helv 10 Tf 0 g"):
// the name of the font among the resources of the form, the size, and the
// color.
func defaultAppearance(da string) (font string, size float64, color bdf.Color) {
	color = bdf.RGB(0, 0, 0)
	toks := strings.Fields(da)
	num := func(i int) float64 {
		if i < 0 {
			return 0
		}
		v, _ := strconv.ParseFloat(toks[i], 64)
		return v
	}
	c8 := func(v float64) uint8 { return uint8(min(max(v, 0), 1)*255 + 0.5) }
	for i, t := range toks {
		switch t {
		case "Tf":
			if i >= 2 {
				font, size = strings.TrimPrefix(toks[i-2], "/"), num(i-1)
			}
		case "g":
			if i >= 1 {
				g := c8(num(i - 1))
				color = bdf.RGB(g, g, g)
			}
		case "rg":
			if i >= 3 {
				color = bdf.RGB(c8(num(i-3)), c8(num(i-2)), c8(num(i-1)))
			}
		case "k":
			if i >= 4 {
				c, m, y, k := num(i-4), num(i-3), num(i-2), num(i-1)
				color = bdf.RGB(c8((1-c)*(1-k)), c8((1-m)*(1-k)), c8((1-y)*(1-k)))
			}
		}
	}
	if size < 0 || size > 1000 {
		size = 0
	}
	return font, size, color
}

// FormFields reads the fields of the form of a PDF, in the order of their
// first widgets on the pages. The pages are those that the options select
// (Options.Pages): a field with no widget on them is left out. A file
// without a form has none.
func FormFields(rs io.ReadSeeker, opts *Options) (fields []FormField, err error) {
	defer recovered(&err)
	c, _, err := newConverter(rs, opts)
	if err != nil {
		return nil, err
	}
	return c.formFields(), nil
}

func (c *converter) formFields() []FormField {
	p := c.pdf
	var form types.Dict
	if cat, err := p.ctx.Catalog(); err == nil && cat != nil {
		form = p.dict(cat["AcroForm"])
	}
	fonts := p.dict(p.dict(form["DR"])["Font"])
	var out []FormField
	index := map[string]int{}
	for pi, pr := range c.pageBodies {
		if pr.err != nil {
			continue
		}
		for _, a := range p.array(pr.dict["Annots"]) {
			ad := p.dict(a)
			if ad == nil || p.name(ad["Subtype"]) != "Widget" {
				continue
			}
			r := p.rectOr(ad["Rect"], rect{})
			if r.empty() {
				continue
			}
			fi := p.fieldOf(ad)
			if fi.name == "" {
				continue
			}
			i, ok := index[fi.name]
			if !ok {
				if len(out) >= maxFormFields {
					continue
				}
				f, ok := c.formField(fi, ad, form, fonts)
				if !ok {
					continue
				}
				i = len(out)
				index[fi.name] = i
				out = append(out, f)
			}
			f := &out[i]
			w := FormWidget{Page: pi, Border: 1}
			// the space of what the widget holds, in that of the page of
			// the file (y up): turned by the quarter turns of the widget
			rot := ((p.intOr(p.dict(ad["MK"])["R"], 0) % 360) + 360) % 360 / 90 * 90
			user := matrix{1, 0, 0, -1, r.x0, r.y1}
			w.W, w.H = r.x1-r.x0, r.y1-r.y0
			switch rot {
			case 90:
				user = matrix{0, 1, 1, 0, r.x0, r.y0}
				w.W, w.H = w.H, w.W
			case 180:
				user = matrix{-1, 0, 0, 1, r.x1, r.y0}
			case 270:
				user = matrix{0, -1, -1, 0, r.x1, r.y1}
				w.W, w.H = w.H, w.W
			}
			w.Matrix = user.mul(pr.base)
			if bs := p.dict(ad["BS"]); bs != nil {
				w.Border = p.numOr(bs["W"], 1)
			} else if b := p.nums(ad["Border"]); len(b) >= 3 {
				w.Border = b[2]
			}
			if p.dict(ad["MK"])["BC"] == nil {
				w.Border = 0 // (no color: no border is drawn)
			}
			if f.Kind == FormCheckbox || f.Kind == FormRadio {
				w.State = p.onState(ad)
				if f.Kind == FormRadio && w.State != "" {
					known := false
					for _, o := range f.Options {
						known = known || o.Value == w.State
					}
					if !known {
						f.Options = append(f.Options, FormOption{Value: w.State, Label: w.State})
					}
				}
			}
			f.Widgets = append(f.Widgets, w)
		}
	}
	return out
}

// formField reads what a field is: ad is its first widget.
func (c *converter) formField(fi fieldInfo, ad, form, fonts types.Dict) (FormField, bool) {
	p := c.pdf
	f := FormField{Name: fi.name, Label: p.text(fi.attr("TU"))}
	flags := p.intOr(fi.attr("Ff"), 0)
	bit := func(n int) bool { return flags&(1<<(n-1)) != 0 }
	f.ReadOnly, f.Required = bit(1), bit(2)
	switch p.name(fi.attr("FT")) {
	case "Tx":
		f.Kind = FormText
		f.Multiline, f.Password, f.Comb = bit(13), bit(14), bit(25)
		f.MaxLen = max(0, p.intOr(fi.attr("MaxLen"), 0))
		if f.MaxLen == 0 {
			f.Comb = false
		}
	case "Btn":
		switch {
		case bit(17):
			f.Kind = FormButton
		case bit(16):
			f.Kind = FormRadio
		default:
			f.Kind = FormCheckbox
		}
	case "Ch":
		f.Kind = FormChoice
		for _, o := range p.array(fi.attr("Opt")) {
			if pair := p.array(o); len(pair) >= 2 {
				f.Options = append(f.Options, FormOption{Value: p.text(pair[0]), Label: p.text(pair[1])})
			} else {
				s := p.text(o)
				f.Options = append(f.Options, FormOption{Value: s, Label: s})
			}
			if len(f.Options) >= 4096 {
				break
			}
		}
	case "Sig":
		f.Kind = FormSignature
	default:
		return f, false
	}
	f.Value, f.Default = p.valueText(fi.attr("V")), p.valueText(fi.attr("DV"))
	da := fi.attr("DA")
	if da == nil {
		da = form["DA"]
	}
	var res string
	res, f.FontSize, f.Color = defaultAppearance(string(p.str(da)))
	fd := p.dict(fonts[res])
	if fd == nil {
		// (a font that the appearance of the widget has, and the form not)
		if sd := p.stream(p.dict(ad["AP"])["N"]); sd != nil {
			fd = p.dict(p.dict(p.dict(sd.Dict["Resources"])["Font"])[res])
		}
	}
	if fd != nil {
		// the font of the resources of the form, by the name it is known by
		f.FontName = p.name(fd["BaseFont"])
		if i := strings.IndexByte(f.FontName, '+'); i == 6 {
			f.FontName = f.FontName[7:]
		}
	}
	q := fi.attr("Q")
	if q == nil {
		q = form["Q"]
	}
	if f.Align = p.intOr(q, 0); f.Align < 0 || f.Align > 2 {
		f.Align = 0
	}
	return f, true
}

// widgetShow tells what a widget annotation shows, when the options have
// something for its field.
func (c *converter) widgetShow(ad types.Dict) (FormShow, bool) {
	if len(c.form) == 0 {
		return FormShow{}, false
	}
	show, ok := c.form[c.pdf.fieldOf(ad).name]
	return show, ok
}
