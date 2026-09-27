package hpgl

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/cad"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/converter/internal/hpglsniff"
)

// A plot file is read as a stream in two contexts, as the devices read it:
// HP-GL/2 instructions, and escape sequences (PCL and HP RTL, with PJL
// lines after a universal exit). Escape sequences switch between them.

type mode int

const (
	modeGL  mode = iota // HP-GL/2 instructions
	modeESC             // PCL or HP RTL escape sequences (and PJL)
)

type converter struct {
	opts     *Options
	warnings []string
	warned   map[string]bool
	set      *fontset.Set
	fonts    *cad.Fonts
	doc      *bdf.Document
	title    string // the picture name of BP
	jobName  string // the job name of PJL
	hasText  bool

	pages []*page
	cur   *page
	gl    *gl
	rtl   *rtl
	pcl   *pcl // the page set-up of a PCL printer, when the job is one

	b    []byte
	i    int
	mode mode

	instructions int // HP-GL/2 instructions read
	rasters      int // raster rows read
	images       int // raster images stored
	items        int // items drawn
	stopped      bool
}

// page is a page being plotted or plotted.
type page struct {
	d cad.Drawing
	// hard is the plot size (the hard-clip limits) in the drawing, when the
	// plot set one.
	hard cad.Rect
	// pcl is the physical page of a PCL printer (the drawing is then in its
	// coordinates, plotter units with y down), nil for a plotter.
	pcl *pclPage
	// marked records HP-GL/2 marks. view turns the drawing to show it:
	// seen in the coordinate system that RO turned, in which the plot was
	// drawn upright (a driver turns HP-GL/2 to the paper that way), or for
	// a page of raster images only as their rows run.
	marked bool
	view   canvas.Matrix
}

func newConverter(opts *Options) *converter {
	c := &converter{opts: opts, warned: map[string]bool{}}
	c.gl = newGL(c)
	c.rtl = newRTL(c)
	c.cur = &page{}
	return c
}

// run reads a whole plot.
func (c *converter) run(b []byte) {
	c.b = b
	c.i = 0
	// a job starts with an escape sequence or PJL; a bare plot starts with
	// its instructions
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t' || b[i] == '\r' || b[i] == '\n' || b[i] == 0) {
		i++
	}
	if i < len(b) && (b[i] == 0x1b && !(i+1 < len(b) && b[i+1] == '.') || bytes.HasPrefix(b[i:], []byte("@PJL"))) {
		c.mode = modeESC
	}
	for c.i < len(b) && !c.stopped {
		switch c.mode {
		case modeGL:
			c.glLoop()
		default:
			c.escLoop()
		}
	}
}

// budget counts an item drawn and reports whether it may be drawn.
func (c *converter) budget() bool {
	if c.items++; c.items > maxItems {
		c.warnOnce("budget", "only the first %d items are drawn", maxItems)
		return false
	}
	return true
}

// endPage ends the page being plotted: a page with something on it is
// kept.
func (c *converter) endPage() {
	c.gl.flush()
	c.gl.closeWindow()
	p := c.cur
	if !p.d.Empty() {
		if len(c.pages) >= maxPages {
			c.warnOnce("pages", "only the first %d pages are converted", maxPages)
			c.stopped = true
		} else {
			if c.printer() {
				p.pcl = c.pcl.page()
			} else if c.gl.psSet {
				p.hard = c.gl.hardRect()
			}
			p.view = invert(c.gl.rotation())
			if !p.marked {
				if v, ok := c.rtl.view(); ok {
					p.view = v
				}
			}
			c.pages = append(c.pages, p)
		}
	}
	c.cur = &page{}
	c.gl.pageAdvanced()
	c.rtl.pageAdvanced()
	c.gl.openWindow()
}

// --- HP-GL/2 context ---

func isLetter(c byte) bool { return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' }

// glLoop reads HP-GL/2 instructions until an escape sequence leaves the
// context or the input ends.
func (c *converter) glLoop() {
	b := c.b
	for c.i < len(b) && !c.stopped {
		ch := b[c.i]
		switch {
		case ch == 0x1b:
			if c.glEscape() {
				return
			}
		case isLetter(ch):
			if c.i+1 >= len(b) {
				c.i = len(b)
				return
			}
			if !isLetter(b[c.i+1]) {
				c.i++
				continue
			}
			m := string([]byte{upper(ch), upper(b[c.i+1])})
			c.i += 2
			if hpglsniff.Known(m) {
				c.instructions++
			}
			c.gl.instruction(m)
		case ch == '!':
			// the device-specific instructions of cutting plotters
			c.i++
			for c.i < len(b) && isLetter(b[c.i]) {
				c.i++
			}
			c.params()
			c.warnOnce("!", "device-specific instructions (!…) are ignored")
		default:
			c.i++
		}
	}
}

func upper(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 32
	}
	return c
}

// glEscape handles an escape sequence met in the HP-GL/2 context and
// reports whether it left the context.
func (c *converter) glEscape() bool {
	b := c.b
	if c.i+1 >= len(b) {
		c.i = len(b)
		return false
	}
	switch b[c.i+1] {
	case '.':
		c.i = hpglsniff.SkipDeviceControl(b, c.i)
		return false
	case 'E':
		c.i += 2
		c.reset()
		c.mode = modeESC
		return true
	case '%':
		v, term, n, ok := escNumber(b[c.i+2:])
		if !ok {
			c.i += 2
			return false
		}
		c.i += 2 + n
		switch term {
		case 'A':
			// enter PCL or HP RTL; 1 moves the cursor to the pen
			c.gl.flush()
			c.rtl.enter(v == 1 || v == 3)
			c.mode = modeESC
			return true
		case 'X':
			c.gl.flush()
			c.mode = modeESC
			return true
		}
		return false
	}
	// another escape sequence: HP RTL or PCL without the switch
	c.gl.flush()
	c.mode = modeESC
	return true
}

// escNumber reads the signed number and terminator of ESC % # X.
func escNumber(b []byte) (v int, term byte, n int, ok bool) {
	i := 0
	neg := false
	if i < len(b) && (b[i] == '-' || b[i] == '+') {
		neg = b[i] == '-'
		i++
	}
	j := i
	for i < len(b) && b[i] >= '0' && b[i] <= '9' && i-j < 9 {
		v = v*10 + int(b[i]-'0')
		i++
	}
	if i >= len(b) {
		return 0, 0, i, false
	}
	if neg {
		v = -v
	}
	return v, b[i], i + 1, true
}

// params reads the numeric parameters of an instruction, up to its
// terminator or the next instruction.
func (c *converter) params() []float64 {
	b := c.b
	var vals []float64
	for c.i < len(b) {
		ch := b[c.i]
		switch {
		case ch == ' ' || ch == ',' || ch == '\t' || ch == '\r' || ch == '\n' || ch == 0:
			c.i++
		case ch >= '0' && ch <= '9' || ch == '.' || ch == '+' || ch == '-':
			v, n := number(b[c.i:])
			c.i += max(n, 1)
			if n > 0 {
				vals = append(vals, v)
			}
		case ch == ';':
			c.i++
			return vals
		default:
			return vals
		}
	}
	return vals
}

// number reads a number: an optional sign, digits and a decimal fraction.
// It returns 0 bytes read when there is no digit.
func number(b []byte) (float64, int) {
	i := 0
	if i < len(b) && (b[i] == '+' || b[i] == '-') {
		i++
	}
	digits := 0
	for i < len(b) && (b[i] >= '0' && b[i] <= '9') {
		i++
		digits++
	}
	if i < len(b) && b[i] == '.' {
		i++
		for i < len(b) && b[i] >= '0' && b[i] <= '9' {
			i++
			digits++
		}
	}
	if digits == 0 {
		return 0, i
	}
	v, err := strconv.ParseFloat(string(b[:i]), 64)
	if err != nil {
		return 0, i
	}
	return v, i
}

// quoted reads a quoted string parameter ("…", with "" for a quote), or
// the characters up to a semicolon when it is not quoted.
func (c *converter) quoted() string {
	b := c.b
	for c.i < len(b) && (b[c.i] == ' ' || b[c.i] == '\t' || b[c.i] == ',') {
		c.i++
	}
	if c.i >= len(b) || b[c.i] != '"' {
		k := bytes.IndexByte(b[c.i:], ';')
		if k < 0 {
			k = len(b) - c.i
		}
		s := string(b[c.i : c.i+k])
		c.i += k
		return s
	}
	c.i++
	var sb strings.Builder
	for c.i < len(b) {
		ch := b[c.i]
		c.i++
		if ch == '"' {
			if c.i < len(b) && b[c.i] == '"' {
				sb.WriteByte('"')
				c.i++
				continue
			}
			break
		}
		sb.WriteByte(ch)
	}
	return sb.String()
}

// --- escape sequence context (PCL, HP RTL, PJL) ---

// escLoop reads escape sequences and PJL until HP-GL/2 is entered or the
// input ends.
func (c *converter) escLoop() {
	b := c.b
	lineStart := true
	for c.i < len(b) && !c.stopped {
		ch := b[c.i]
		switch {
		case ch == 0x1b:
			if c.escape() {
				return
			}
			// PJL follows a universal exit
			lineStart = bytes.HasSuffix(b[:c.i], []byte("%-12345X"))
			continue
		case ch == '@' && lineStart && bytes.HasPrefix(b[c.i:], []byte("@PJL")):
			if c.pjl() {
				return
			}
			lineStart = true
			continue
		case ch == '\r' || ch == '\n':
			lineStart = true
		case ch == 0x0c:
			// a form feed ejects the page of a PCL printer
			if c.printer() {
				c.rtl.end()
				c.endPage()
			}
		case ch == ' ' || ch == '\t' || ch == 0:
		default:
			lineStart = false
			if ch >= 0x20 && ch < 0x7f {
				// text: a bare plot after the PJL of its job
				if isLetter(ch) && c.pcl == nil && hpglsniff.Instructions(b[c.i:min(len(b), c.i+512)]) >= 2 {
					c.mode = modeGL
					return
				}
				c.warnOnce("pcltext", "the text of the PCL job is not drawn")
			}
		}
		c.i++
	}
}

// pjl reads PJL lines and reports whether they entered HP-GL/2.
func (c *converter) pjl() bool {
	b := c.b
	for c.i < len(b) && bytes.HasPrefix(b[c.i:], []byte("@PJL")) {
		end := bytes.IndexAny(b[c.i:], "\r\n")
		if end < 0 {
			end = len(b) - c.i
		}
		line := string(b[c.i : c.i+end])
		c.i += end
		for c.i < len(b) && (b[c.i] == '\r' || b[c.i] == '\n') {
			c.i++
		}
		f := strings.Fields(strings.ReplaceAll(line[4:], "=", " = "))
		if len(f) == 0 {
			continue
		}
		upperAll := func(i int) string {
			if i < len(f) {
				return strings.ToUpper(f[i])
			}
			return ""
		}
		switch strings.ToUpper(f[0]) {
		case "ENTER":
			if upperAll(1) == "LANGUAGE" {
				lang := upperAll(3)
				if upperAll(2) != "=" {
					lang = upperAll(2)
				}
				switch {
				case strings.HasPrefix(lang, "HPGL"):
					c.gl.standalone()
					c.mode = modeGL
					return true
				case lang == "RTL":
					c.gl.standalone()
				case lang == "PCL":
					c.gl.enterPrinter()
				}
			}
		case "JOB":
			if m := jobName.FindStringSubmatch(line); m != nil {
				c.jobName = strings.ToValidUTF8(m[1]+m[2], "")
			}
		case "SET":
			key := upperAll(1)
			val := upperAll(3)
			if upperAll(2) != "=" {
				val = upperAll(2)
			}
			switch key {
			case "RESOLUTION":
				if v, err := strconv.Atoi(val); err == nil && v > 0 && v <= 4800 {
					c.rtl.native = float64(v)
				}
			case "RENDERMODE":
				c.gl.grayscale = val == "GRAYSCALE"
			}
		}
	}
	return false
}

// jobName matches the name of a PJL job, quoted or not.
var jobName = regexp.MustCompile(`(?i)\bNAME\s*=\s*(?:"([^"]*)"|(\S+))`)

// escape reads an escape sequence and reports whether it entered HP-GL/2.
func (c *converter) escape() bool {
	b := c.b
	if c.i+1 >= len(b) {
		c.i = len(b)
		return false
	}
	p := b[c.i+1]
	switch p {
	case 'E':
		c.i += 2
		c.reset()
		return false
	case '%':
		v, term, n, ok := escNumber(b[c.i+2:])
		c.i += 2 + n
		if !ok {
			return false
		}
		switch term {
		case 'B':
			c.rtl.end()
			c.rtl.leave()
			c.gl.enter(v)
			c.mode = modeGL
			return true
		case 'A':
			c.rtl.enter(v == 1 || v == 3)
		}
		return false
	case '.':
		c.i = hpglsniff.SkipDeviceControl(b, c.i)
		return false
	case '&', '*', '(', ')':
	default:
		// a two-character sequence (ESC 9, ESC =, ESC Y …)
		c.i += 2
		return false
	}
	c.i += 2
	if c.i >= len(b) {
		return false
	}
	g := b[c.i]
	if (p == '(' || p == ')') && (g >= '0' && g <= '9') {
		// a symbol set designation: ESC ( 8U
		c.escValues(p, 0)
		return false
	}
	c.i++
	c.escValues(p, g)
	return false
}

// escValues reads the values and terminators of a parameterized escape
// sequence (ESC p g #t#t…#T) and carries out each command.
func (c *converter) escValues(p, g byte) {
	b := c.b
	for c.i < len(b) {
		// value: sign, digits, fraction
		start := c.i
		rel := false
		if c.i < len(b) && (b[c.i] == '+' || b[c.i] == '-') {
			rel = true
			c.i++
		}
		for c.i < len(b) && (b[c.i] >= '0' && b[c.i] <= '9' || b[c.i] == '.') {
			c.i++
		}
		v, _ := strconv.ParseFloat(string(b[start:c.i]), 64)
		if c.i >= len(b) {
			return
		}
		t := b[c.i]
		c.i++
		last := t >= 'A' && t <= 'Z'
		cmd := upper(t)
		if !last && !(t >= 'a' && t <= 'z') {
			// not a terminator: the sequence was malformed
			c.i--
			return
		}
		var data []byte
		if cmd == 'W' || p == '*' && g == 'b' && cmd == 'V' || p == '&' && g == 'p' && cmd == 'X' {
			n := int(v)
			if n < 0 {
				n = 0
			}
			if n > len(b)-c.i {
				n = len(b) - c.i
			}
			data = b[c.i : c.i+n]
			c.i += n
		}
		c.command(p, g, cmd, v, rel, data)
		if last {
			return
		}
	}
}

// command carries out an escape sequence command of PCL or HP RTL.
func (c *converter) command(p, g, cmd byte, v float64, rel bool, data []byte) {
	switch p {
	case '*':
		switch g {
		case 'r', 'b', 't', 'v':
			c.rtl.command(g, cmd, v, rel, data)
		case 'p':
			c.rtl.cursor(g, cmd, v, rel)
		case 'c':
			// the picture frame (X, Y, K, L, T); patterns are not drawn
			if cmd == 'X' || cmd == 'Y' || cmd == 'K' || cmd == 'L' || cmd == 'T' {
				p := c.pclSetup()
				p.cursor = c.rtl.cap
				p.pictureFrame(cmd, v)
			}
		}
	case '&':
		switch g {
		case 'a':
			c.rtl.cursor(g, cmd, v, rel)
		case 'l', 'u':
			c.pageCommand(g, cmd, v)
		}
	}
	// the rest (fonts, text) is not drawn
}

// reset is ESC E: the page ends, and every setting returns to its
// default.
func (c *converter) reset() {
	c.rtl.end()
	c.endPage()
	c.gl.reset()
	c.rtl.reset()
	if c.pcl != nil {
		c.pcl.reset()
	}
}
