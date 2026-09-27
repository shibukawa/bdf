// Package gerber converts the fabrication data of printed circuit boards
// into BDF documents: Gerber files (RS-274X, as the Ucamco Gerber Layer
// Format Specification defines it, with the X2 attributes and the
// deprecated commands of older files), Excellon drill files, and a zip
// archive of the files of a board (with its Gerber job file, when it has
// one).
//
// A board becomes a view of its top and one of its bottom, drawn the way
// the board looks (substrate, copper, solder mask, the exposed copper's
// finish, silkscreen and holes), and a view for each file, drawn in the
// colors of PCB design tools on their dark background. A single file
// becomes a view of its layer. What each file is (top copper, bottom
// solder mask, ...) comes from its X2 file attributes, the job file, the
// naming conventions of KiCad, Altium (Protel), Eagle and others, or its
// content (drill files). See docs/design.md §3.19.
package gerber

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/shibukawa/bdf"
)

// Views selects the views of a board.
type Views int

const (
	// ViewsAll makes the top and bottom of the board and a view for each
	// file.
	ViewsAll Views = iota
	// ViewsBoard makes the top and bottom of the board only.
	ViewsBoard
	// ViewsLayers makes a view for each file only.
	ViewsLayers
)

// Options controls the conversion.
type Options struct {
	// Title overrides the document title.
	Title string
	// FileName is the name of a single file converted, which tells its
	// layer when its content does not.
	FileName string
	// Views selects the views of a board.
	Views Views
	// Mask, Silkscreen and Finish are the colors of the solder mask
	// (green, red, blue, black, white, yellow, purple or #rrggbb), the
	// silkscreen (white, black, yellow or #rrggbb) and the exposed copper
	// (gold, silver or copper); "" takes those of the job file, else
	// green, white and gold.
	Mask, Silkscreen, Finish string
	// Warn receives non-fatal problems; when nil they are collected in
	// Result.Warnings.
	Warn func(msg string)
}

// Result is the outcome of a conversion.
type Result struct {
	Doc      *bdf.Document
	Warnings []string
	// Layers is the number of files drawn.
	Layers int
	// Board reports that the top and bottom views were made.
	Board bool
	// W and H are the size of the board (or of the drawing) in mm.
	W, H float64
	// Scale is how many times larger than the board its pages are.
	Scale float64
}

const maxSize = 1 << 30

// ConvertFile converts a file: a Gerber or drill file, or a zip archive of
// the files of a board.
func ConvertFile(name string, opts *Options) (*Result, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.FileName == "" {
		o.FileName = path.Base(strings.ReplaceAll(name, "\\", "/"))
	}
	return Convert(bytes.NewReader(data), int64(len(data)), &o)
}

// input is a file to convert.
type input struct {
	name string
	data []byte
}

// Convert converts a Gerber or drill file, or a zip archive of the files
// of a board, read from r.
func Convert(r io.ReaderAt, size int64, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if size > maxSize {
		return nil, fmt.Errorf("gerber: the file is larger than %d bytes", maxSize)
	}
	data := make([]byte, size)
	if _, err := r.ReadAt(data, 0); err != nil && err != io.EOF {
		return nil, fmt.Errorf("gerber: %w", err)
	}
	c := &conv{opts: opts, warned: map[string]bool{}, count: new(int)}
	var inputs []input
	archive := bytes.HasPrefix(data, []byte("PK\x03\x04"))
	if archive {
		var err error
		if inputs, err = c.unzip(data); err != nil {
			return nil, fmt.Errorf("gerber: %w", err)
		}
	} else {
		inputs = []input{{name: opts.FileName, data: data}}
	}
	for _, in := range inputs {
		c.read(in, !archive)
	}
	if len(c.layers) == 0 {
		switch {
		case c.files > 0:
			return nil, fmt.Errorf("gerber: the files draw nothing")
		case archive:
			return nil, fmt.Errorf("gerber: the archive holds no Gerber or Excellon files")
		}
		return nil, fmt.Errorf("gerber: not a Gerber or Excellon file")
	}
	c.identify()
	res, err := c.build()
	if err != nil {
		return nil, err
	}
	res.Warnings = c.warnings
	return res, nil
}

// conv is the state of a conversion.
type conv struct {
	opts     *Options
	warnings []string
	warned   map[string]bool
	count    *int
	layers   []*layer
	// files counts the Gerber and drill files read, drawing or not
	files int
	job   *jobFile
	// file is the file warnings are about
	file string
}

func (c *conv) warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if c.file != "" {
		msg = c.file + ": " + msg
	}
	if c.warned[msg] {
		return
	}
	c.warned[msg] = true
	if c.opts.Warn != nil {
		c.opts.Warn(msg)
		return
	}
	c.warnings = append(c.warnings, msg)
}

// maxEntries and maxUnzipped bound what is read from an archive.
const (
	maxEntries  = 1000
	maxUnzipped = 1 << 30
)

// unzip returns the files of a zip archive (not its directories, nor the
// resource forks and hidden files of macOS).
func (c *conv) unzip(data []byte) ([]input, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	var out []input
	total := int64(0)
	for _, f := range zr.File {
		name := f.Name
		base := path.Base(name)
		if f.FileInfo().IsDir() || strings.HasPrefix(name, "__MACOSX/") || strings.HasPrefix(base, ".") {
			continue
		}
		if len(out) >= maxEntries {
			c.warnf("the archive holds more than %d files; the rest are left out", maxEntries)
			break
		}
		rc, err := f.Open()
		if err != nil {
			c.warnf("%s: %v", name, err)
			continue
		}
		b, err := io.ReadAll(io.LimitReader(rc, maxUnzipped-total+1))
		rc.Close()
		if err != nil {
			c.warnf("%s: %v", name, err)
			continue
		}
		total += int64(len(b))
		if total > maxUnzipped {
			return nil, fmt.Errorf("the archive holds more than %d bytes", maxUnzipped)
		}
		out = append(out, input{name: name, data: b})
	}
	// the files in the order of their names
	slices.SortStableFunc(out, func(a, b input) int { return strings.Compare(a.name, b.name) })
	return out, nil
}

// read reads a file into a layer; single is set for a file converted on
// its own, which is read as a Gerber file when its content tells nothing.
func (c *conv) read(in input, single bool) {
	c.file = ""
	if !single {
		c.file = in.name
	}
	defer func() { c.file = "" }()
	switch sniff(in.name, in.data) {
	case "job":
		var j jobFile
		if err := json.Unmarshal(in.data, &j); err != nil {
			c.warnf("the job file is not read: %v", err)
			return
		}
		c.job = &j
		return
	case "excellon":
		l := parseExcellon(in.data, c.count, c.warnf)
		l.name = in.name
		c.add(l)
		return
	case "gerber":
	default:
		// drill files without a header, and for a file on its own, whatever
		// Gerber commands it holds
		if looksLikeDrill(in.data) && (single || drillExtension(in.name)) {
			l := parseExcellon(in.data, c.count, c.warnf)
			l.name = in.name
			c.add(l)
			return
		}
		if !single {
			c.warnf("not a Gerber or Excellon file; left out")
			return
		}
	}
	l, err := parseGerber(in.data, c.count, c.warnf)
	if err != nil {
		c.warnf("%v", err)
		return
	}
	l.name = in.name
	c.add(l)
}

func (c *conv) add(l *layer) {
	c.files++
	if l.img.empty() {
		name := l.name
		if name == "" {
			name = "the file"
		}
		c.warnf("%s draws nothing", name)
		return
	}
	c.layers = append(c.layers, l)
}

// identify tells what each layer is and sorts them from the top of the
// board down.
func (c *conv) identify() {
	for _, l := range c.layers {
		l.fn = c.function(l)
	}
	slices.SortStableFunc(c.layers, func(a, b *layer) int {
		if d := a.fn.order() - b.fn.order(); d != 0 {
			return d
		}
		return strings.Compare(a.name, b.name)
	})
}

func (c *conv) function(l *layer) function {
	base := strings.ToLower(path.Base(strings.ReplaceAll(l.name, "\\", "/")))
	if c.job != nil {
		for _, fa := range c.job.FilesAttributes {
			if strings.ToLower(path.Base(strings.ReplaceAll(fa.Path, "\\", "/"))) == base {
				if f, ok := fileFunction(fa.FileFunction); ok {
					return c.plating(f, l)
				}
			}
		}
	}
	if f, ok := fileFunction(l.attrs["FileFunction"]); ok {
		return c.plating(f, l)
	}
	if l.drill {
		f, _ := byName(l.name, true)
		return c.plating(f, l)
	}
	if f, ok := byName(l.name, false); ok {
		return f
	}
	return function{}
}

// plating tells the plating of a drill layer that its function leaves
// open.
func (c *conv) plating(f function, l *layer) function {
	if f.kind == kDrill && f.plated == "" {
		f.plated = l.plated
		if f.plated == "" {
			g, _ := byName(l.name, true)
			f.plated = g.plated
		}
	}
	return f
}

// jobFile is what the converter reads of a Gerber job file (.gbrjob).
type jobFile struct {
	Header struct {
		CreationDate string
	}
	GeneralSpecs struct {
		ProjectId struct {
			Name string
		}
		Finish string
	}
	FilesAttributes []struct {
		Path         string
		FileFunction string
		FilePolarity string
	}
	MaterialStackup []struct {
		Type  string
		Color string
		Name  string
	}
}

// sniff tells what a file is from its content: "gerber", "excellon",
// "job" (a Gerber job file) or "".
func sniff(name string, data []byte) string {
	head := data[:min(len(data), 64<<10)]
	switch {
	case strings.EqualFold(path.Ext(name), ".gbrjob") || isJob(head):
		return "job"
	case IsGerber(head):
		return "gerber"
	case IsExcellon(head):
		return "excellon"
	}
	return ""
}

func isJob(head []byte) bool {
	t := bytes.TrimLeft(bytes.TrimPrefix(head, []byte("\xef\xbb\xbf")), " \t\r\n")
	return bytes.HasPrefix(t, []byte("{")) && bytes.Contains(head, []byte(`"GeneralSpecs"`))
}

// text reports whether b looks like text: no control characters but
// white space.
func text(b []byte) bool {
	for _, c := range b {
		if c < ' ' && c != '\n' && c != '\r' && c != '\t' || c == 0x7f {
			return false
		}
	}
	return true
}

// IsGerber reports whether the first bytes of a file (up to 64 KiB, so
// that long comment headers are passed) are those of a Gerber file in the
// extended format (RS-274X).
func IsGerber(head []byte) bool {
	t := bytes.TrimLeft(head, " \t\r\n")
	if len(t) == 0 || !text(head) {
		return false
	}
	// the first command: an extended command of the format, or a word
	// command (a comment, most often)
	switch t[0] {
	case '%':
		if len(t) < 3 || !gerberCommands[string(t[1:3])] {
			return false
		}
	case 'G', 'M', 'D', 'N':
		if len(t) < 2 || !isDigit(t[1]) {
			return false
		}
	default:
		return false
	}
	return bytes.Contains(head, []byte("%FS")) || bytes.Contains(head, []byte("%MO")) || bytes.Contains(head, []byte("%AD"))
}

// gerberCommands are the extended commands a Gerber file can start with.
var gerberCommands = map[string]bool{
	"FS": true, "MO": true, "TF": true, "TA": true, "TO": true, "TD": true, "AD": true, "AM": true, "IN": true, "IP": true,
	"LN": true, "LP": true, "OF": true, "SF": true, "AS": true, "MI": true, "IR": true, "IC": true, "IJ": true, "IO": true,
}

// IsExcellon reports whether the first bytes of a file are those of an
// Excellon drill file: a header (M48), after comments and a rewind stop
// (%).
func IsExcellon(head []byte) bool {
	if !text(head) {
		return false
	}
	for n, line := range bytes.Split(head, []byte("\n")) {
		s := bytes.TrimSpace(line)
		switch {
		case n > 200:
			return false
		case len(s) == 0 || s[0] == ';' || bytes.Equal(s, []byte("%")):
			continue
		case bytes.HasPrefix(s, []byte("M48")):
			return true
		}
		return false
	}
	return false
}

// drillExtension reports whether a file name has an extension of drill
// files.
func drillExtension(name string) bool {
	switch strings.ToLower(path.Ext(name)) {
	case ".drl", ".drd", ".xln", ".exc", ".txt", ".nc", ".tap", ".cnc", ".drill":
		return true
	}
	return false
}

// looksLikeDrill reports whether a file without a header looks like drill
// data: most of its lines choose tools (T01) or give coordinates (X1.0Y2.0),
// and it holds both.
func looksLikeDrill(data []byte) bool {
	head := data[:min(len(data), 64<<10)]
	if !text(head) || bytes.IndexByte(head, '*') >= 0 {
		return false
	}
	tools, coords, lines := 0, 0, 0
	for _, line := range bytes.Split(head, []byte("\n")) {
		s := bytes.TrimSpace(line)
		if len(s) == 0 || s[0] == ';' || s[0] == '%' {
			continue
		}
		lines++
		switch {
		case len(s) > 1 && s[0] == 'T' && isDigit(s[1]):
			tools++
		case len(s) > 1 && (s[0] == 'X' || s[0] == 'Y') && (isDigit(s[1]) || s[1] == '-' || s[1] == '+' || s[1] == '.'):
			coords++
		}
	}
	return tools > 0 && coords > 0 && 2*(tools+coords) >= lines
}
