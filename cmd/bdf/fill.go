package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/flate"
	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/bdf/form"
	"github.com/shibukawa/bdf/form/tag"
	"github.com/shibukawa/bdf/image/imgconv"
	"github.com/shibukawa/bdf/raster/imagebdf"
)

// formFlags are the flags of the commands that fill templates.
type formFlags struct {
	fs            *flag.FlagSet
	cmd           string
	data          stringList
	files         *string
	strict        *bool
	pages         *string
	quiet         *bool
	images        *string
	quality       *int
	noSubset      *bool
	noWOFF2       *bool
	ignoreFSType  *bool
	fonts         *string
	fontDirs      stringList
	noSystemFonts *bool
	params        stringList
	compression   *string
	// pictures is the directory of the pictures that the data names, once
	// it is opened (see options)
	pictures fs.FS
}

func addFormFlags(fs *flag.FlagSet, cmd, dataFlag, dataUsage string) *formFlags {
	f := &formFlags{fs: fs, cmd: cmd}
	fs.Var(&f.data, dataFlag, dataUsage)
	f.files = fs.String("files", "", "directory of the pictures that the data names (default: that of the template)")
	f.strict = fs.Bool("strict", false, "fail when the data lacks a name of a tag (default: leave the tag empty, with a warning)")
	f.pages = fs.String("pages", "", "sheets or slides of the template to fill, e.g. 1,3; for a Word template the pages of the form to keep (default: all)")
	f.quiet = fs.Bool("q", false, "do not print warnings")
	f.images = fs.String("images", "convert", "pictures: keep (store as they are) or convert (try WebP, keep when smaller)")
	f.quality = fs.Int("quality", 80, "lossy WebP quality (1-100)")
	f.noSubset = fs.Bool("no-subset", false, "embed whole fonts instead of the glyphs in use")
	f.noWOFF2 = fs.Bool("no-woff2", false, "store embedded fonts as TrueType/OpenType instead of WOFF2")
	f.ignoreFSType = fs.Bool("ignore-fstype", false, "embed fonts whose OS/2 fsType forbids embedding or subsetting (only with the rights to do so)")
	f.fonts = fs.String("fonts", "", "embed (subset and embed the fonts used for layout; the default) or system (refer to fonts by name)")
	fs.Var(&f.fontDirs, "font-dir", "directory searched for fonts before the system ones (repeatable)")
	f.noSystemFonts = fs.Bool("no-system-fonts", false, "use only the fonts under -font-dir")
	fs.Var(&f.params, "param", "option of the format of the template as name=value, e.g. paper=a4 (repeatable; see bdf generate -h)")
	f.compression = fs.String("compression", "normal", "deflate of the parts: normal (level 6), best (level 9) or fast (level 1)")
	return f
}

func (f *formFlags) fail(msg string) {
	fmt.Fprintln(os.Stderr, "bdf "+f.cmd+":", msg)
	os.Exit(2)
}

// convert returns the options of the conversion that the flags give.
func (f *formFlags) convert() *converter.Options {
	img := imgconv.Options{Quality: *f.quality}
	switch *f.images {
	case "keep":
		img.Mode = imgconv.Keep
	case "convert":
		img.Mode = imgconv.Convert
	default:
		f.fail("-images must be keep or convert")
	}
	sel, err := converter.ParsePages(*f.pages)
	if err != nil {
		f.fail("-pages: " + err.Error())
	}
	opts := &converter.Options{Pages: sel, Images: img, FontDirs: f.fontDirs, NoSystemFonts: *f.noSystemFonts, NoSubset: *f.noSubset,
		NoWOFF2: *f.noWOFF2, IgnoreFSType: *f.ignoreFSType, Params: map[string]string{}}
	switch *f.fonts {
	case "", "embed":
	case "system":
		opts.SystemFonts = true
	default:
		f.fail("-fonts must be embed or system")
	}
	for _, kv := range f.params {
		name, value, ok := strings.Cut(kv, "=")
		if !ok || name == "" {
			f.fail("-param " + kv + ": want name=value")
		}
		opts.Params[name] = value
	}
	return opts
}

// sources reads the files of data that flags name: file, or name=file.
func (f *formFlags) sources(flags []string) []form.Source {
	var out []form.Source
	for _, v := range flags {
		name, file, ok := strings.Cut(v, "=")
		if !ok || strings.ContainsAny(name, `/\`) {
			name, file = "", v
		}
		b, err := os.ReadFile(file)
		check(err)
		out = append(out, form.Source{Name: name, FileName: file, Data: b})
	}
	return out
}

// options returns the options of filling the template at path with data.
// The pictures that the data names are read from the directory of -files
// (by default that of the template) and from it alone: a path that leaves
// it, by .. or by a link, is not read.
func (f *formFlags) options(path string, data map[string]any) *form.Options {
	if f.pictures == nil {
		dir := *f.files
		if dir == "" {
			dir = filepath.Dir(path)
		}
		root, err := os.OpenRoot(dir)
		if err != nil {
			f.fail("-files: " + err.Error())
		}
		f.pictures = root.FS()
	}
	opts := f.convert()
	opts.FileName = filepath.Base(path)
	return &form.Options{Data: data, Strict: *f.strict, Files: f.pictures, Convert: opts}
}

// finish gives a document the compression the flags ask for.
func (f *formFlags) finish(doc *bdf.Document) {
	switch *f.compression {
	case "normal":
	case "best":
		doc.CompressionLevel = flate.BestCompression
	case "fast":
		doc.CompressionLevel = flate.BestSpeed
	default:
		f.fail("-compression must be normal, best or fast")
	}
}

func (f *formFlags) warn(warnings []string) {
	if !*f.quiet {
		for _, w := range warnings {
			fmt.Fprintln(os.Stderr, "warning:", w)
		}
	}
}

// fill fills a template with data and writes the form.
func fill(args []string) {
	fs := flag.NewFlagSet("fill", flag.ExitOnError)
	f := addFormFlags(fs, "fill", "data", "data of the template: a JSON or YAML object as file, or a value under a name as name=file: a list of a CSV or TSV file or of the sheets of a workbook (.xlsx), a JSON or YAML value (repeatable; a later object replaces the values of the same names)")
	each := fs.String("each", "", "fill the template once for each item of a list, given as -data gives one (file, or name=file to put the item under a name instead of making its values those of the data); the output path names each form with tags: out/{{no}}.bdf")
	preview := addPreviewFlags(fs)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bdf fill [flags] <template.xlsx | .docx | .pptx | .pdf | .json | .bdft> <out.bdf | outdir/>\n  fills a template with data, and writes the form. A template is an Excel workbook, a Word document or a\n  PowerPoint presentation whose text holds tags such as {{name}} (the form is the pages the workbook prints,\n  the pages of the document or the slides), a PDF with form fields or a template of pdfme (.json), whose\n  fields are filled in their places, or a compiled template")
		fs.PrintDefaults()
	}
	fs.Parse(args)
	if fs.NArg() != 2 {
		fs.Usage()
		os.Exit(2)
	}
	in, out := fs.Arg(0), fs.Arg(1)
	tpl, err := os.ReadFile(in)
	check(err)
	data, err := form.LoadData(f.sources(f.data)...)
	check(err)
	one := func(data map[string]any, out string) {
		res, err := form.Fill(bytes.NewReader(tpl), int64(len(tpl)), f.options(in, data))
		check(err)
		f.finish(res.Doc)
		f.warn(res.Warnings)
		pages := 0
		for _, v := range res.Doc.Views {
			pages += len(v.Pages)
		}
		check(write(res.Doc, out))
		fmt.Fprintf(os.Stderr, "%s: %d page(s), %d warning(s)\n", out, pages, len(res.Warnings))
		preview.write(res.Doc, imagebdf.Options{FontDirs: f.fontDirs, NoSystemFonts: *f.noSystemFonts}, *f.quiet)
	}
	if *each == "" {
		one(data, out)
		return
	}
	// a form for each item
	name, items := eachItems(f, *each)
	toks, err := tag.Parse(out)
	if err != nil {
		f.fail("the output path: " + err.Error())
	}
	tagged := false
	for _, t := range toks {
		if tagged = tagged || t.Kind == tag.Value; t.Kind != tag.Text && t.Kind != tag.Value {
			f.fail("the output path takes tags of values only ({{name}})")
		}
	}
	if !tagged {
		f.fail("with -each, the output path names each form with a tag: " + strings.TrimSuffix(out, filepath.Ext(out)) + "-{{_index}}" + filepath.Ext(out))
	}
	seen := map[string]bool{}
	for i, it := range items {
		d := make(map[string]any, len(data)+4)
		for k, v := range data {
			d[k] = v
		}
		if obj, ok := it.(map[string]any); ok && name == "" {
			for k, v := range obj {
				d[k] = v
			}
		} else if name != "" {
			d[name] = it
		} else {
			f.fail(fmt.Sprintf("-each: item %d is not an object; give the items a name (-each item=%s)", i+1, *each))
		}
		d["_index"], d["_count"] = float64(i+1), float64(len(items))
		sc := tag.NewScope(d)
		var b strings.Builder
		for _, t := range toks {
			if t.Kind == tag.Text {
				b.WriteString(t.Text)
				continue
			}
			v, ok := sc.Lookup(t.Path)
			if !ok || v == nil {
				check(fmt.Errorf("-each: item %d has no %q to name its form with", i+1, t.Name))
			}
			b.WriteString(fileNamePart(v))
		}
		path := b.String()
		if seen[path] {
			check(fmt.Errorf("-each: items %d and another are both written to %s", i+1, path))
		}
		seen[path] = true
		check(os.MkdirAll(filepath.Dir(path), 0o755))
		one(d, path)
	}
}

// eachItems reads the list of -each: its name ("" to merge the items into
// the data) and its items.
func eachItems(f *formFlags, flag string) (string, []any) {
	name, file, ok := strings.Cut(flag, "=")
	if !ok || strings.ContainsAny(name, `/\`) {
		name, file = "", flag
	}
	b, err := os.ReadFile(file)
	check(err)
	d, err := form.LoadData(form.Source{Name: "each", FileName: file, Data: b})
	check(err)
	items, ok := d["each"].([]any)
	if !ok {
		f.fail("-each: " + file + " is not a list")
	}
	return name, items
}

// fileNamePart writes a value as a part of a file name: without what
// separates directories.
func fileNamePart(v any) string {
	s := fmt.Sprint(v)
	if n, ok := v.(float64); ok && n == float64(int64(n)) {
		s = fmt.Sprint(int64(n))
	}
	s = strings.Map(func(r rune) rune {
		switch {
		case r == '/' || r == '\\' || r == ':' || r < ' ':
			return '_'
		}
		return r
	}, strings.TrimSpace(s))
	if s == "" || s == "." || s == ".." {
		return "_"
	}
	return s
}

// templateCmd compiles templates and tells what they ask of data.
func templateCmd(args []string) {
	usage := func() {
		fmt.Fprintln(os.Stderr, `usage:
  bdf template compile [flags] <template.xlsx | .docx | .pptx | .pdf | .json> <out.bdft>
                                     check a template and write it as a compiled template: a BDF document
                                     that shows the template (filled with -sample data) and holds it
  bdf template schema [-json-schema] <template.xlsx | .docx | .pptx | .pdf | .json | .bdft> [out.json]
                                     write the names that a template asks of data (the names its tags look
                                     up, or its fields), as JSON (with -json-schema: a JSON Schema of the data)

  A template is an Office document with tags, a PDF with form fields, or a template of pdfme (.json).`)
		os.Exit(2)
	}
	if len(args) < 1 {
		usage()
	}
	switch args[0] {
	case "compile":
		fs := flag.NewFlagSet("template compile", flag.ExitOnError)
		f := addFormFlags(fs, "template compile", "sample", "sample data that the pages of the compiled template show the template filled with, given as bdf fill -data gives data (repeatable; default: the pages show the tags, and those of a template of pdfme what its designer left in its fields)")
		fs.Usage = func() {
			fmt.Fprintln(os.Stderr, "usage: bdf template compile [flags] <template.xlsx | .docx | .pptx | .pdf | .json> <out.bdft>")
			fs.PrintDefaults()
		}
		fs.Parse(args[1:])
		if fs.NArg() != 2 {
			fs.Usage()
			os.Exit(2)
		}
		in, out := fs.Arg(0), fs.Arg(1)
		tpl, err := os.ReadFile(in)
		check(err)
		var sample map[string]any
		if len(f.data) > 0 {
			sample, err = form.LoadData(f.sources(f.data)...)
			check(err)
		}
		c, err := form.Compile(bytes.NewReader(tpl), int64(len(tpl)), filepath.Base(in), f.options(in, sample))
		check(err)
		f.finish(c.Doc)
		f.warn(c.Warnings)
		if isDir(out) {
			f.fail("a compiled template is one file")
		}
		check(writeSingle(c.Doc, out))
		pages := 0
		for _, v := range c.Doc.Views {
			pages += len(v.Pages)
		}
		fmt.Fprintf(os.Stderr, "%s: %d field(s), %d page(s), %d warning(s)\n", out, len(c.Template.Fields), pages, len(c.Warnings))
	case "schema":
		fs := flag.NewFlagSet("template schema", flag.ExitOnError)
		jsonSchema := fs.Bool("json-schema", false, "write a JSON Schema of the data instead of the list of the names: what a program makes a form for the data of")
		var samples stringList
		fs.Var(&samples, "sample", "with -json-schema: sample data, given as bdf fill -data gives data (repeatable), which tells what the tags do not: where a name is, and the type of its value (default: the schema that a compiled template holds, that of its sample)")
		fs.Usage = func() {
			fmt.Fprintln(os.Stderr, "usage: bdf template schema [flags] <template.xlsx | .docx | .pptx | .pdf | .json | .bdft> [out.json]")
			fs.PrintDefaults()
		}
		fs.Parse(args[1:])
		if fs.NArg() < 1 || fs.NArg() > 2 {
			fs.Usage()
			os.Exit(2)
		}
		tpl, err := os.ReadFile(fs.Arg(0))
		check(err)
		fields, err := form.Fields(bytes.NewReader(tpl), int64(len(tpl)))
		check(err)
		if fields == nil {
			fields = []tag.Field{}
		}
		var b []byte
		if *jsonSchema {
			// (a compiled template has its schema; a sample makes another)
			var sample map[string]any
			if len(samples) > 0 {
				f := &formFlags{}
				sample, err = form.LoadData(f.sources(samples)...)
				check(err)
			}
			b, err = form.SchemaOf(bytes.NewReader(tpl), int64(len(tpl)), sample)
			check(err)
		} else {
			b, err = json.MarshalIndent(map[string]any{"fields": fields}, "", "  ")
			check(err)
			b = append(b, '\n')
		}
		if fs.NArg() == 2 && fs.Arg(1) != "-" {
			check(os.WriteFile(fs.Arg(1), b, 0o644))
			return
		}
		os.Stdout.Write(b)
	default:
		usage()
	}
}
