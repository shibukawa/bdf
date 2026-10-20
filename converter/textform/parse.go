package textform

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/shibukawa/bdf/form/tag"
)

// Names of sections that are not looked up in the data: they say how the
// lines or characters between their tags are printed.
const (
	secWide  = "_wide"
	secTall  = "_tall"
	secBig   = "_big"
	secPaper = "preprinted"
)

// attribute flags of characters and lines
const (
	attrWide uint8 = 1 << iota
	attrTall
	attrPaper
)

// attrOf returns the flags of a section of attributes (0 for a section of
// the data).
func attrOf(name string) uint8 {
	switch name {
	case secWide:
		return attrWide
	case secTall:
		return attrTall
	case secBig:
		return attrWide | attrTall
	case secPaper:
		return attrPaper
	}
	return 0
}

// line is a line of the body of a template.
type line struct {
	// no is its number in the file (from 1), and src its text.
	no   int
	src  string
	toks []tag.Token
	// match is the token that closes each Open or Inverted of the same
	// line (-1 for one that closes on another line or with a mode: a
	// section of lines).
	match []int
	// lineLevel marks the tags of the sections of lines.
	lineLevel []bool
	// brk marks a page break (a form feed) instead of a line.
	brk bool
	// tagOnly marks a line of the tags of sections of lines and nothing
	// else: it is not printed.
	tagOnly bool
	// tagged marks a line with tags other than those of sections of lines.
	tagged bool
	// attrs are the flags of the sections of attributes that the line is
	// in.
	attrs uint8
}

func (l *line) at() string { return fmt.Sprintf("line %d", l.no) }

// section is a section of lines: the lines from that of its opening tag to
// that of its closing tag, repeated for the items of its value.
type section struct {
	name     string
	path     []string
	inverted bool
	mode     string
	// a and b are the lines (indexes of template.lines) of its tags, and
	// open the line and the token of its opening tag.
	a, b     int
	open     [2]int
	parent   *section
	children []*section
	at       string
}

// slots returns the lines of a section of slots: those that are printed.
func (s *section) slots(t *template) []int {
	var out []int
	for i := s.a; i <= s.b; i++ {
		if !t.lines[i].tagOnly && !t.lines[i].brk {
			out = append(out, i)
		}
	}
	return out
}

// template is a template of text read.
type template struct {
	grid  grid
	lines []*line
	root  *section
	// grows is set when a section adds lines: pages then break by
	// themselves when they are full.
	grows bool
	// sections maps the opening tags of sections of lines to them, by
	// line and token.
	sections map[[2]int]*section
	warnings []string
}

func (t *template) warnf(format string, args ...any) {
	t.warnings = append(t.warnings, fmt.Sprintf(format, args...))
}

// Limits of a template (it is not trusted).
const (
	maxTemplate = 16 << 20
	maxColumns  = 4096
	maxLines    = 4096
	// maxBodyLines bounds the lines of the body of a template.
	maxBodyLines = 1 << 18
)

// readAll reads a template, bounded.
func readAll(r io.ReaderAt, size int64) ([]byte, error) {
	if size < 0 {
		return nil, fmt.Errorf("text: negative template size")
	}
	if size > maxTemplate {
		return nil, fmt.Errorf("text: the template is larger than %d MiB", maxTemplate>>20)
	}
	data := make([]byte, size)
	if _, err := io.ReadFull(io.NewSectionReader(r, 0, size), data); err != nil {
		return nil, fmt.Errorf("text: %w", err)
	}
	return data, nil
}

// IsText reports whether the start of a file is text of UTF-8 (without
// NUL bytes; its last character may be cut by the end of the start).
func IsText(head []byte) bool {
	_, ok := text(head)
	return ok
}

func text(head []byte) ([]byte, bool) {
	t := bytes.TrimPrefix(head, []byte("\xef\xbb\xbf"))
	if len(t) == 0 || bytes.IndexByte(t, 0) >= 0 {
		return nil, false
	}
	for i := 0; i < 3 && len(t) > 0 && !utf8.Valid(t); i++ {
		t = t[:len(t)-1]
	}
	return t, utf8.Valid(t)
}

// IsTemplate reports whether a file looks like a template of text by its
// start: text of UTF-8 (IsText) that starts with front matter (a line
// "---") or has a tag.
func IsTemplate(head []byte) bool {
	t, ok := text(head)
	return ok && (bytes.HasPrefix(t, []byte("---\n")) || bytes.HasPrefix(t, []byte("---\r\n")) || tag.Has(string(t)))
}

// parse reads a template: its front matter, the lines of its body and
// their tags, and its sections of lines.
func parse(data []byte) (*template, error) {
	s := strings.TrimPrefix(string(data), "\ufeff")
	if !utf8.ValidString(s) {
		return nil, fmt.Errorf("text: the template is not text of UTF-8")
	}
	s = strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
	t := &template{sections: map[[2]int]*section{}}
	front, body, first, _ := splitFrontMatter(s)
	g, err := parseFrontMatter(front, t.warnf)
	if err != nil {
		return nil, err
	}
	t.grid = g
	body = strings.TrimSuffix(body, "\n")
	raw := strings.Split(body, "\n")
	if body == "" {
		raw = nil
	}
	if len(raw) > maxBodyLines {
		return nil, fmt.Errorf("text: the template has more than %d lines", maxBodyLines)
	}
	for i, r := range raw {
		no := first + i
		// a form feed breaks the page; the text around it is on the lines
		// before and after the break
		parts := strings.Split(r, "\f")
		for k, p := range parts {
			if k > 0 {
				t.lines = append(t.lines, &line{no: no, brk: true})
			}
			if p == "" && len(parts) > 1 {
				continue
			}
			toks, err := tag.ParseText(p)
			if err != nil {
				if e, ok := err.(*tag.Error); ok {
					return nil, fmt.Errorf("text: line %d, column %d: %s", no, columnOf(p, e.Pos), e.Msg)
				}
				return nil, fmt.Errorf("text: line %d: %w", no, err)
			}
			for _, tk := range toks {
				if tk.Kind == tag.Image {
					return nil, fmt.Errorf("text: line %d: {{%%%s}}: a template of text has no pictures nor codes (its lines are printed in the text mode of printers)", no, tk.Name)
				}
			}
			t.lines = append(t.lines, &line{no: no, src: p, toks: toks})
		}
	}
	if err := t.sectionsOf(); err != nil {
		return nil, err
	}
	return t, nil
}

// columnOf is the column (from 1) of a byte of a line, for a message.
func columnOf(s string, pos int) int {
	if pos > len(s) {
		pos = len(s)
	}
	return utf8.RuneCountInString(s[:pos]) + 1
}

// open is a tag that opens a section, waiting for its closing tag.
type open struct {
	li, ti int
	tok    tag.Token
}

// sectionsOf matches the tags of the sections, makes the tree of the
// sections of lines and gives the lines in sections of attributes their
// flags.
func (t *template) sectionsOf() error {
	t.root = &section{a: 0, b: len(t.lines) - 1}
	var stack []open
	for li, l := range t.lines {
		l.match = make([]int, len(l.toks))
		l.lineLevel = make([]bool, len(l.toks))
		for ti, tk := range l.toks {
			l.match[ti] = -1
			switch tk.Kind {
			case tag.Open, tag.Inverted:
				if a := attrOf(tk.Name); a != 0 && (tk.Kind == tag.Inverted || tk.Mode != "") {
					return fmt.Errorf("text: %s: {{#%s}} marks how characters are printed; it is not inverted and takes no %q", l.at(), tk.Name, tk.Mode)
				}
				stack = append(stack, open{li, ti, tk})
			case tag.Close:
				if len(stack) == 0 {
					return fmt.Errorf("text: %s: {{/%s}} closes no section", l.at(), tk.Name)
				}
				o := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if tk.Name != "" && tk.Name != o.tok.Name {
					return fmt.Errorf("text: %s: {{/%s}} closes {{#%s}} of %s", l.at(), tk.Name, o.tok.Name, t.lines[o.li].at())
				}
				if o.li == li && o.tok.Mode == "" {
					// a section of the characters of a line
					l.match[o.ti] = ti
					continue
				}
				l.lineLevel[ti], t.lines[o.li].lineLevel[o.ti] = true, true
				if a := attrOf(o.tok.Name); a != 0 {
					for k := o.li; k <= li; k++ {
						t.lines[k].attrs |= a
					}
					continue
				}
				sec := &section{name: o.tok.Name, path: o.tok.Path, inverted: o.tok.Kind == tag.Inverted, mode: o.tok.Mode, a: o.li, b: li, open: [2]int{o.li, o.ti}, at: t.lines[o.li].at()}
				t.sections[[2]int{o.li, o.ti}] = sec
				t.place(sec)
			}
		}
	}
	if len(stack) > 0 {
		o := stack[len(stack)-1]
		return fmt.Errorf("text: %s: the section {{#%s}} is not closed with {{/%s}}", t.lines[o.li].at(), o.tok.Name, o.tok.Name)
	}
	for _, l := range t.lines {
		if l.brk {
			continue
		}
		l.tagOnly = len(l.toks) > 0
		for ti, tk := range l.toks {
			switch {
			case l.lineLevel[ti]:
			case tk.Kind == tag.Text && strings.TrimSpace(tk.Text) == "":
			default:
				l.tagOnly = false
			}
			if !l.lineLevel[ti] && tk.Kind != tag.Text {
				l.tagged = true
			}
		}
	}
	return t.check(t.root)
}

// place puts a section of lines in the tree when its closing tag is read.
// The sections it holds closed before it, and are children of the root
// until then: those whose opening tags follow its own.
func (t *template) place(sec *section) {
	var kept []*section
	for _, ch := range t.root.children {
		if ch.open[0] > sec.open[0] || ch.open[0] == sec.open[0] && ch.open[1] > sec.open[1] {
			ch.parent = sec
			sec.children = append(sec.children, ch)
			continue
		}
		kept = append(kept, ch)
	}
	sec.parent = t.root
	t.root.children = append(kept, sec)
}

// check gives the sections of lines their modes and refuses those that
// cannot be filled.
func (t *template) check(sec *section) error {
	for i, ch := range sec.children {
		switch {
		case ch.inverted:
			if ch.mode == tag.ModeSlots || ch.mode == tag.ModePages {
				return fmt.Errorf("text: %s: the inverted section {{^%s}} does not repeat: it takes no %q", ch.at, ch.name, ch.mode)
			}
			ch.mode = tag.ModeRows
		case ch.mode == "" && ch.a == ch.b:
			ch.mode = tag.ModeRows
		case ch.mode == "":
			ch.mode = tag.ModeSlots
		}
		if ch.mode == tag.ModeRows && !ch.inverted {
			t.grows = true
		}
		if ch.mode == tag.ModeSlots && len(ch.children) > 0 {
			in := ch.children[0]
			return fmt.Errorf("text: %s: the section {{#%s}} is in the lines of {{#%s}}, which are slots (a line for an item); write {{#%s | rows}} to repeat the lines instead",
				in.at, in.name, ch.name, ch.name)
		}
		if ch.mode == tag.ModeSlots && len(ch.slots(t)) == 0 {
			return fmt.Errorf("text: %s: the section {{#%s}} has no line to fill", ch.at, ch.name)
		}
		if i > 0 && ch.a <= sec.children[i-1].b {
			prev := sec.children[i-1]
			return fmt.Errorf("text: %s: the sections {{#%s}} and {{#%s}} both repeat line %d; a line is repeated by one section, or by sections one inside the other",
				ch.at, prev.name, ch.name, t.lines[ch.a].no)
		}
		if err := t.check(ch); err != nil {
			return err
		}
	}
	return nil
}

// sectionAt returns the innermost section of lines that a line is in (the
// root for none).
func (t *template) sectionAt(li int) *section {
	sec := t.root
	for {
		var in *section
		for _, ch := range sec.children {
			if ch.a <= li && li <= ch.b {
				in = ch
				break
			}
		}
		if in == nil {
			return sec
		}
		sec = in
	}
}

// fields lists the names that the tags of the template look up, in the
// order of its lines.
func (t *template) fields() []tag.Field {
	var out []tag.Field
	index := map[string]int{}
	add := func(fd tag.Field, at string) {
		// (the variables of sections are not of the data; "." is the item
		// of the section around it, and tells that its items are values)
		if strings.HasPrefix(fd.Name, "_") || fd.Name == "." && len(fd.In) == 0 || fd.Kind == tag.FieldSection && attrOf(fd.Name) != 0 {
			return
		}
		key := fd.Kind + "\x00" + strings.Join(fd.In, "\x00") + "\x00\x00" + fd.Name
		if i, ok := index[key]; ok {
			old := &out[i]
			if old.Format == "" {
				old.Format = fd.Format
			}
			if old.Mode == "" {
				old.Mode, old.Slots = fd.Mode, fd.Slots
			}
			old.Inverted = old.Inverted && fd.Inverted
			if n := len(old.At); n < 16 && old.At[n-1] != at {
				old.At = append(old.At, at)
			}
			return
		}
		fd.In = append([]string(nil), fd.In...)
		fd.At = []string{at}
		index[key] = len(out)
		out = append(out, fd)
	}
	for li, l := range t.lines {
		if l.brk {
			continue
		}
		var in []string
		for sec := t.sectionAt(li); sec != t.root; sec = sec.parent {
			in = append([]string{sec.name}, in...)
		}
		for ti, tk := range l.toks {
			switch tk.Kind {
			case tag.Value:
				add(tag.Field{Name: tk.Name, Kind: tag.FieldValue, In: in, Format: tk.Format}, l.at())
			case tag.Open, tag.Inverted:
				fd := tag.Field{Name: tk.Name, Kind: tag.FieldSection, In: in, Inverted: tk.Kind == tag.Inverted}
				if sec := t.sections[[2]int{li, ti}]; sec != nil {
					// the section of lines this tag opens: it is looked up
					// in the sections around it
					fd.In = nil
					for p := sec.parent; p != t.root; p = p.parent {
						fd.In = append([]string{p.name}, fd.In...)
					}
					if fd.Mode = sec.mode; sec.mode == tag.ModeSlots {
						fd.Slots = len(sec.slots(t))
					}
					add(fd, l.at())
					continue
				}
				add(fd, l.at())
				if !l.lineLevel[ti] && attrOf(tk.Name) == 0 {
					in = append(append([]string(nil), in...), tk.Name)
				}
			case tag.Close:
				if !l.lineLevel[ti] && len(in) > 0 && l.match != nil {
					// (the close of a section of characters of the data)
					for k := ti - 1; k >= 0; k-- {
						if l.match[k] == ti {
							if attrOf(l.toks[k].Name) == 0 {
								in = in[:len(in)-1]
							}
							break
						}
					}
				}
			}
		}
	}
	return out
}
