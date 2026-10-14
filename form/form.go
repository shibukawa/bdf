// Package form fills templates with data and makes the forms that
// business systems print (invoices, slips, certificates) as BDF documents.
//
// A template is an Excel workbook, a Word document or a PowerPoint
// presentation whose text holds tags (package form/tag): {{name}} for a
// value, {{#items}} … {{/items}} for what a list repeats. Filling it lays
// the file out with the tags replaced by the data: a workbook as the pages
// Excel prints, a document as its pages, a presentation as its slides. See
// Fill of converter/xlsx, converter/docx and converter/pptx for what the
// tags do in each.
//
// A template is also a file whose fields have places: a PDF with the
// fields of a form (its AcroForm; converter/pdfform), and a template of
// pdfme (the JSON file that its designer saves; converter/pdfme). Their
// fields are the names that the data has values for.
//
//	data, err := form.LoadData(form.Source{FileName: "order.json", Data: order},
//		form.Source{Name: "items", FileName: "lines.csv", Data: lines})
//	res, err := form.Fill(template, size, &form.Options{Data: data})
//	err = res.Doc.WriteSingle(out)
//
// The data is an object of the values of JSON. Sums, taxes, what is
// carried over to the next page and the rounding of all of them are the
// data's: a template looks values up and formats them, and does not
// calculate with them (the formulas of a workbook aside).
//
// Compile makes a compiled template (.bdft) of a template: a BDF document
// that holds the template with what it asks of data, and shows it filled
// with sample data. It is checked once, and is filled like the template.
package form

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/converter/docx"
	"github.com/shibukawa/bdf/converter/pdfform"
	"github.com/shibukawa/bdf/converter/pdfme"
	"github.com/shibukawa/bdf/converter/pptx"
	"github.com/shibukawa/bdf/converter/xlsx"
	"github.com/shibukawa/bdf/form/tag"
)

// Options controls filling a template.
type Options struct {
	// Data is the data: an object of the values of JSON (see LoadData).
	Data map[string]any
	// Strict makes a tag whose name the data does not have an error;
	// without it the tag is replaced by nothing, with a warning.
	Strict bool
	// Files holds the pictures that the data names by their paths.
	Files fs.FS
	// Convert are the options of the conversion: the fonts, the images,
	// the sheets or slides to fill (Pages; for a Word document the pages
	// to keep) and the parameters of the format ("views", "paper" …). nil
	// converts with the defaults.
	Convert *converter.Options
}

// Result is a filled template.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	// Fields are the names that the tags of the template look up.
	Fields []tag.Field
	// Format is the format of the template (FormatXLSX …).
	Format string
}

// Formats of templates: the documents of Office with tags in their text,
// and the files whose fields are what is filled: a PDF with the fields of
// a form, and a template of pdfme.
const (
	FormatXLSX  = "xlsx"
	FormatDOCX  = "docx"
	FormatPPTX  = "pptx"
	FormatPDF   = "pdf"
	FormatPDFME = "pdfme"
)

const notTemplate = "form: the template is not an Excel workbook (.xlsx), a Word document (.docx), a PowerPoint presentation (.pptx), a PDF with form fields (.pdf), a template of pdfme (.json) nor a compiled template (.bdft)"

// source finds what a template is: a compiled template gives the template
// it holds.
func source(r io.ReaderAt, size int64) (format string, src io.ReaderAt, n int64, def *Template, err error) {
	if size < 0 {
		return "", nil, 0, nil, fmt.Errorf("form: negative template size")
	}
	head := make([]byte, min(size, 1024))
	got, err := r.ReadAt(head, 0)
	if err != nil && err != io.EOF {
		return "", nil, 0, nil, fmt.Errorf("form: %w", err)
	}
	head = head[:got]
	if bytes.HasPrefix(head, bdf.Magic[:]) {
		def, data, err := ReadTemplate(r, size)
		if err != nil {
			return "", nil, 0, nil, err
		}
		return def.Format, bytes.NewReader(data), int64(len(data)), def, nil
	}
	if bytes.Contains(head, []byte("%PDF-")) {
		return FormatPDF, r, size, nil, nil
	}
	if pdfme.IsTemplate(head) {
		return FormatPDFME, r, size, nil, nil
	}
	if !bytes.HasPrefix(head, []byte("PK\x03\x04")) {
		return "", nil, 0, nil, fmt.Errorf(notTemplate)
	}
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return "", nil, 0, nil, fmt.Errorf("form: %w", err)
	}
	for _, f := range zr.File {
		switch strings.ToLower(strings.ReplaceAll(f.Name, "\\", "/")) {
		case "xl/workbook.xml":
			return FormatXLSX, r, size, nil, nil
		case "word/document.xml":
			return FormatDOCX, r, size, nil, nil
		case "ppt/presentation.xml":
			return FormatPPTX, r, size, nil, nil
		}
	}
	return "", nil, 0, nil, fmt.Errorf(notTemplate)
}

// Fill fills a template, or a compiled template, with data.
func Fill(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	format, src, n, _, err := source(r, size)
	if err != nil {
		return nil, err
	}
	return fill(format, src, n, opts)
}

func fill(format string, r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	switch format {
	case FormatXLSX:
		xo, err := xlsx.OptionsOf(opts.Convert)
		if err != nil {
			return nil, err
		}
		res, err := xlsx.Fill(r, size, &xlsx.TemplateData{Data: opts.Data, Strict: opts.Strict, Files: opts.Files}, xo)
		if err != nil {
			return nil, err
		}
		return &Result{Doc: res.Doc, Warnings: res.Warnings, Fields: res.Fields, Format: format}, nil
	case FormatDOCX:
		// (the number formats of tags are those of Excel: its formatter)
		res, err := docx.Fill(r, size, &docx.TemplateData{Data: opts.Data, Strict: opts.Strict, Files: opts.Files, Format: xlsx.Formatter()}, docx.OptionsOf(opts.Convert))
		if err != nil {
			return nil, err
		}
		return &Result{Doc: res.Doc, Warnings: res.Warnings, Fields: res.Fields, Format: format}, nil
	case FormatPPTX:
		po, err := pptx.OptionsOf(opts.Convert)
		if err != nil {
			return nil, err
		}
		res, err := pptx.Fill(r, size, &pptx.TemplateData{Data: opts.Data, Strict: opts.Strict, Files: opts.Files, Format: xlsx.Formatter()}, po)
		if err != nil {
			return nil, err
		}
		return &Result{Doc: res.Doc, Warnings: res.Warnings, Fields: res.Fields, Format: format}, nil
	case FormatPDF:
		res, err := pdfform.Fill(io.NewSectionReader(r, 0, size), &pdfform.TemplateData{Data: opts.Data, Strict: opts.Strict, Format: xlsx.Formatter()}, pdfform.OptionsOf(opts.Convert))
		if err != nil {
			return nil, err
		}
		return &Result{Doc: res.Doc, Warnings: res.Warnings, Fields: res.Fields, Format: format}, nil
	case FormatPDFME:
		res, err := pdfme.Fill(r, size, &pdfme.TemplateData{Data: opts.Data, Strict: opts.Strict, Files: opts.Files, Format: xlsx.Formatter()}, pdfme.OptionsOf(opts.Convert))
		if err != nil {
			return nil, err
		}
		return &Result{Doc: res.Doc, Warnings: res.Warnings, Fields: res.Fields, Format: format}, nil
	}
	return nil, fmt.Errorf("form: templates of the format %q are not filled", format)
}

// Fields reads the tags of a template without filling it: the names that
// data for it has, or the first mistake in its tags.
func Fields(r io.ReaderAt, size int64) ([]tag.Field, error) {
	format, src, n, def, err := source(r, size)
	if err != nil {
		return nil, err
	}
	if def != nil {
		return def.Fields, nil
	}
	switch format {
	case FormatXLSX:
		return xlsx.Inspect(src, n, &xlsx.Options{NoSystemFonts: true, SystemFonts: true})
	case FormatDOCX:
		return docx.Inspect(src, n)
	case FormatPPTX:
		return pptx.Inspect(src, n, nil)
	case FormatPDF:
		return pdfform.Inspect(io.NewSectionReader(src, 0, n), nil)
	case FormatPDFME:
		return pdfme.Inspect(src, n)
	}
	return nil, fmt.Errorf("form: templates of the format %q are not read", format)
}

// Describe tells what a template is: its format (FormatXLSX …; that of the
// template a compiled template holds), and the sample of its data that it
// has for itself: of a compiled template the sample it was compiled with,
// of a template of pdfme what its designer left in its fields. sample is
// nil for a template without one.
func Describe(r io.ReaderAt, size int64) (format string, sample map[string]any, err error) {
	format, src, n, def, err := source(r, size)
	if err != nil {
		return "", nil, err
	}
	switch {
	case def != nil:
		return format, def.Sample, nil
	case format == FormatPDFME:
		sample, err = pdfme.Defaults(src, n)
		return format, sample, err
	}
	return format, nil, nil
}

// TemplateVersion is the version of the description of compiled templates.
const TemplateVersion = 1

// MaxTemplateSize bounds the bytes of the source stored in a compiled
// template. Both compiling and reading a compiled template enforce it.
const MaxTemplateSize = 256 << 20

// Template describes a compiled template: the part of the type
// bdf.PartTemplate of a .bdft, as JSON.
type Template struct {
	// Version is TemplateVersion.
	Version int `json:"version"`
	// Format is the format of the template the document holds (FormatXLSX
	// …), Name its file name when it had one, and Source the hash of the part
	// (of the type bdf.PartSource) that is the template.
	Format string `json:"format"`
	Name   string `json:"name,omitempty"`
	Source string `json:"source"`
	// Fields are the names that the tags of the template look up.
	Fields []tag.Field `json:"fields"`
	// Sample is the data the pages of the document show the template
	// filled with (none: they show the tags).
	Sample map[string]any `json:"sample,omitempty"`
	// Schema is the JSON Schema of the data (Schema, made of the fields
	// and the sample): what a program makes a form for the data of,
	// without reading the template.
	Schema json.RawMessage `json:"schema,omitempty"`
}

// Compiled is a compiled template.
type Compiled struct {
	// Doc is the document to write as the .bdft: its pages are the
	// template filled with the sample data.
	Doc      *bdf.Document
	Template *Template
	Warnings []string
}

// Compile checks a template and makes its compiled form. The data of the
// options is the sample that the pages of the compiled template show the
// template filled with; without data they show the template with its tags
// (a template of pdfme with what its designer left in its fields, which is
// then the sample). name is the file name of the template, when it has one.
// The source may be at most MaxTemplateSize bytes.
func Compile(r io.ReaderAt, size int64, name string, opts *Options) (*Compiled, error) {
	if size > MaxTemplateSize {
		return nil, fmt.Errorf("form: the template is larger than %d MiB", MaxTemplateSize>>20)
	}
	if opts == nil {
		opts = &Options{}
	}
	format, src, n, def, err := source(r, size)
	if err != nil {
		return nil, err
	}
	if def != nil {
		return nil, fmt.Errorf("form: the template is a compiled template already")
	}
	data := make([]byte, n)
	if _, err := io.ReadFull(io.NewSectionReader(src, 0, n), data); err != nil {
		return nil, fmt.Errorf("form: %w", err)
	}
	out := &Compiled{Template: &Template{Version: TemplateVersion, Format: format, Name: name}}
	if opts.Data == nil && format == FormatPDFME {
		sample, err := pdfme.Defaults(bytes.NewReader(data), n)
		if err != nil {
			return nil, err
		}
		o := *opts
		o.Data, o.Strict = sample, false
		opts = &o
	}
	if opts.Data != nil {
		res, err := fill(format, bytes.NewReader(data), n, opts)
		if err != nil {
			return nil, err
		}
		out.Doc, out.Warnings, out.Template.Fields, out.Template.Sample = res.Doc, res.Warnings, res.Fields, opts.Data
	} else {
		// the template as it is, with its tags, after they are checked
		fields, err := Fields(bytes.NewReader(data), n)
		if err != nil {
			return nil, err
		}
		o := converter.Options{}
		if opts.Convert != nil {
			o = *opts.Convert
		}
		o.Params = map[string]string{}
		if opts.Convert != nil {
			for k, v := range opts.Convert.Params {
				o.Params[k] = v
			}
		}
		// (as a filled template has them: the pages)
		if o.Params["views"] == "" {
			switch format {
			case FormatXLSX:
				o.Params["views"] = xlsx.ViewsPrint
			case FormatDOCX:
				o.Params["views"] = docx.ViewsPages
			}
		}
		res, err := converter.Convert(bytes.NewReader(data), n, format, &o)
		if err != nil {
			return nil, err
		}
		out.Doc, out.Warnings, out.Template.Fields = res.Doc, res.Warnings, fields
	}
	if out.Template.Fields == nil {
		out.Template.Fields = []tag.Field{}
	}
	if out.Template.Schema, err = Schema(out.Template.Fields, out.Template.Sample); err != nil {
		return nil, fmt.Errorf("form: the sample data is not data of JSON: %w", err)
	}
	out.Template.Source = out.Doc.AddPart(bdf.PartSource, data).String()
	def2, err := json.Marshal(out.Template)
	if err != nil {
		return nil, fmt.Errorf("form: the sample data is not data of JSON: %w", err)
	}
	out.Doc.AddPart(bdf.PartTemplate, def2)
	return out, nil
}

// maxTemplatePart bounds the description of a compiled template.
const maxTemplatePart = 64 << 20

// ReadTemplate reads a compiled template: its description, and the
// template it holds.
func ReadTemplate(r io.ReaderAt, size int64) (*Template, []byte, error) {
	rd, err := bdf.OpenSingle(r, size)
	if err != nil {
		return nil, nil, fmt.Errorf("form: %w", err)
	}
	if rd.Locked() {
		return nil, nil, fmt.Errorf("form: the compiled template is encrypted")
	}
	for _, e := range rd.Manifest.Parts {
		if e.T != bdf.PartTemplate {
			continue
		}
		b, err := templatePart(rd, e.H, bdf.PartTemplate, "description of the compiled template", maxTemplatePart)
		if err != nil {
			return nil, nil, err
		}
		def := &Template{}
		if err := json.Unmarshal(b, def); err != nil {
			return nil, nil, fmt.Errorf("form: the description of the compiled template: %w", err)
		}
		if def.Version != TemplateVersion {
			return nil, nil, fmt.Errorf("form: the compiled template is of version %d; this reads version %d (compile the template again)", def.Version, TemplateVersion)
		}
		h, err := bdf.ParseHash(def.Source)
		if err != nil {
			return nil, nil, fmt.Errorf("form: the compiled template does not name its template: %w", err)
		}
		src, err := templatePart(rd, h, bdf.PartSource, "template of the compiled template", MaxTemplateSize)
		if err != nil {
			return nil, nil, err
		}
		for i := range def.Fields {
			if def.Fields[i].At == nil {
				def.Fields[i].At = []string{}
			}
		}
		if def.Sample != nil {
			tag.Normalize(def.Sample)
		}
		return def, src, nil
	}
	return nil, nil, fmt.Errorf("form: the document is not a compiled template (it has no %q part)", bdf.PartTemplate)
}

// templatePart checks the entry Reader.Part will actually load. Duplicate
// hashes in a manifest may replace the entry encountered in the loop.
// Bound both stored and decoded lengths before either is allocated.
func templatePart(rd *bdf.Reader, h bdf.Hash, typ, name string, limit int) ([]byte, error) {
	e, ok := rd.Entry(h)
	if !ok || e.T != typ {
		return nil, fmt.Errorf("form: the compiled template does not hold its %s", name)
	}
	if e.Size < 0 || e.Len < 0 {
		return nil, fmt.Errorf("form: the %s has a negative size", name)
	}
	if e.Size > limit || e.Len > limit {
		return nil, fmt.Errorf("form: the %s is larger than %d MiB", name, limit>>20)
	}
	if (e.Enc == bdf.EncIdentity || e.Enc == "") && e.Len != e.Size {
		return nil, fmt.Errorf("form: the %s has inconsistent stored and decoded sizes", name)
	}
	data, err := rd.Part(h)
	if err != nil {
		return nil, fmt.Errorf("form: %w", err)
	}
	if len(data) != e.Size || bdf.HashOf(data) != h {
		return nil, fmt.Errorf("form: the %s is damaged", name)
	}
	return data, nil
}
