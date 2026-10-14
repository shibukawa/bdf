package xlsx

import (
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// The header and the footer of a printed sheet are strings of text and
// codes (ECMA-376 §18.3.1.39): &L, &C and &R start the left, the center and
// the right section, &P and &N are the page number and the number of pages,
// &D, &T, &F and &A the date, the time, the file and the sheet, and the
// others format the text after them (&"font,style", &size, &B, &I, &U, &E,
// &S, &X, &Y, &Kcolor).

// headerFooter is the headers and footers of a sheet.
type headerFooter struct {
	// odd, even and first pages: the header and the footer
	odd, even, first [2]string
	differentOddEven bool
	differentFirst   bool
	scaleWithDoc     bool
	alignWithMargins bool
}

func readHeaderFooter(n *ooxml.Node) *headerFooter {
	// line breaks are written as _x000a_
	text := func(name string) string { return decodeXString(n.Child(name).Content()) }
	return &headerFooter{
		odd:              [2]string{text("oddHeader"), text("oddFooter")},
		even:             [2]string{text("evenHeader"), text("evenFooter")},
		first:            [2]string{text("firstHeader"), text("firstFooter")},
		differentOddEven: n.AttrBool("differentOddEven", false),
		differentFirst:   n.AttrBool("differentFirst", false),
		scaleWithDoc:     n.AttrBool("scaleWithDoc", true),
		alignWithMargins: n.AttrBool("alignWithMargins", true),
	}
}

// codes returns the header and the footer of the page numbered n, when the
// first page has the number first.
func (hf *headerFooter) codes(n, first int) (header, footer string) {
	if hf == nil {
		return "", ""
	}
	switch {
	case hf.differentFirst && n == first:
		return hf.first[0], hf.first[1]
	case hf.differentOddEven && n%2 == 0:
		return hf.even[0], hf.even[1]
	}
	return hf.odd[0], hf.odd[1]
}

// hfRun is text of a header or a footer in one style.
type hfRun struct {
	text string
	st   *tstyle
}

// maxHFText bounds the characters of a header or a footer that are drawn.
const maxHFText = 4096

// hfSections reads the codes of a header or a footer into its left, center
// and right sections: the lines of each, as runs of text.
func (p *printPlan) hfSections(code string, page, pages int) [3][][]hfRun {
	s := p.s
	base := *s.textStyle(s.c.st.font(0), nil, nil)
	if p.ps.hf != nil && p.ps.hf.scaleWithDoc {
		base.size *= p.scale
	}
	var out [3][][]hfRun
	sec, st := 1, base
	var buf strings.Builder
	flush := func() {
		if buf.Len() == 0 {
			return
		}
		cp := st
		if len(out[sec]) == 0 {
			out[sec] = [][]hfRun{nil}
		}
		k := len(out[sec]) - 1
		out[sec][k] = append(out[sec][k], hfRun{buf.String(), &cp})
		buf.Reset()
	}
	rs := []rune(code)
	if len(rs) > maxHFText {
		rs = rs[:maxHFText]
	}
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if r == '\n' {
			flush()
			if len(out[sec]) == 0 {
				out[sec] = [][]hfRun{nil}
			}
			out[sec] = append(out[sec], nil)
			continue
		}
		if r == '\r' {
			continue
		}
		if r != '&' || i+1 >= len(rs) {
			buf.WriteRune(r)
			continue
		}
		i++
		switch c := rs[i]; c {
		case '&':
			buf.WriteRune('&')
		case 'L', 'C', 'R':
			flush()
			sec, st = strings.IndexRune("LCR", c), base
		case 'P':
			n := page
			// &P+n and &P-n add to the page number
			if i+2 < len(rs) && (rs[i+1] == '+' || rs[i+1] == '-') && rs[i+2] >= '0' && rs[i+2] <= '9' {
				j, d := i+2, 0
				for ; j < len(rs) && rs[j] >= '0' && rs[j] <= '9' && d < 1<<20; j++ {
					d = d*10 + int(rs[j]-'0')
				}
				if rs[i+1] == '-' {
					d = -d
				}
				n, i = n+d, j-1
			}
			buf.WriteString(strconv.Itoa(n))
		case 'N':
			buf.WriteString(strconv.Itoa(pages))
		case 'D':
			if t := s.c.opts.PrintTime; !t.IsZero() {
				layout := "1/2/2006"
				if s.c.locale == "ja" {
					layout = "2006/1/2"
				}
				buf.WriteString(t.Format(layout))
			}
		case 'T':
			if t := s.c.opts.PrintTime; !t.IsZero() {
				buf.WriteString(t.Format("15:04"))
			}
		case 'F':
			buf.WriteString(s.c.opts.FileName)
		case 'A':
			buf.WriteString(s.ws.name)
		case 'Z':
			// the path of the file: not the reader's to see
		case 'G':
			s.c.warnOnce("hfpicture", "sheet %q: the pictures of headers and footers are not drawn", s.ws.name)
		case 'B':
			flush()
			st.bold = !st.bold
		case 'I':
			flush()
			st.italic = !st.italic
		case 'U', 'E':
			flush()
			want := "single"
			if c == 'E' {
				want = "double"
			}
			if st.underline == want {
				want = ""
			}
			st.underline = want
		case 'S':
			flush()
			st.strike = !st.strike
		case 'X', 'Y':
			flush()
			want := "superscript"
			if c == 'Y' {
				want = "subscript"
			}
			if st.vert == want {
				want = ""
			}
			st.vert = want
		case 'O', 'H':
			// outline and shadow: not drawn
		case '"':
			j := i + 1
			for j < len(rs) && rs[j] != '"' {
				j++
			}
			if j >= len(rs) {
				break
			}
			flush()
			name, style, _ := strings.Cut(string(rs[i+1:j]), ",")
			if name = strings.TrimSpace(name); name != "" && name != "-" {
				st.latin, st.ea = name, name
			} else {
				st.latin, st.ea = base.latin, base.ea
			}
			style = strings.ToLower(style)
			st.bold = strings.Contains(style, "bold") || strings.Contains(style, "太字")
			st.italic = strings.Contains(style, "italic") || strings.Contains(style, "斜体")
			i = j
		case 'K':
			// a color: RRGGBB, or a theme color and its tint (kept as it is)
			if i+6 < len(rs) {
				if v, err := strconv.ParseUint(string(rs[i+1:i+7]), 16, 32); err == nil {
					flush()
					st.color = bdf.RGB(uint8(v>>16), uint8(v>>8), uint8(v))
					i += 6
					break
				}
			}
			if i+6 < len(rs) && (rs[i+3] == '+' || rs[i+3] == '-') {
				i += 6
			}
		default:
			if c >= '0' && c <= '9' {
				j, size := i, 0
				for ; j < len(rs) && rs[j] >= '0' && rs[j] <= '9' && size < 1000; j++ {
					size = size*10 + int(rs[j]-'0')
				}
				if size >= 1 && size <= 409 {
					flush()
					st.size = float64(size)
					if p.ps.hf != nil && p.ps.hf.scaleWithDoc {
						st.size *= p.scale
					}
				}
				i = j - 1
			}
			// an unknown code is passed over
		}
	}
	flush()
	return out
}

// hfLayer draws a header or a footer of the page numbered page of a sheet
// whose last page has the number pages; nil when there is none.
func (p *printPlan) hfLayer(code string, header bool, page, pages int) *canvas.Canvas {
	if code == "" {
		return nil
	}
	s := p.s
	secs := p.hfSections(code, page, pages)
	base := s.textStyle(s.c.st.font(0), nil, nil)
	cv := s.c.cvs.New()
	t := &tileCv{cv: cv}
	left, right := p.ml, p.pw-p.mr
	if p.ps.hf != nil && !p.ps.hf.alignWithMargins {
		left, right = 0.7*72, p.pw-0.7*72
	}
	for k, lines := range secs {
		var tls []*tline
		height, any := 0.0, false
		for _, runs := range lines {
			ln := &tline{}
			for _, run := range runs {
				for _, r := range run.text {
					ln.items = append(ln.items, s.item(run.st, r))
				}
			}
			any = any || len(ln.items) > 0
			lineMetrics(ln, base, s)
			height += ln.asc + ln.desc
			tls = append(tls, ln)
		}
		if !any {
			continue
		}
		y := p.mh
		if !header {
			y = p.ph - p.mf - height
		}
		cv.Obj.Mark(bdf.MarkParagraph, "")
		for i, ln := range tls {
			switch k {
			case 0:
				ln.x = left
			case 1:
				ln.x = (left+right)/2 - ln.width/2
			default:
				ln.x = right - ln.width
			}
			ln.base = y + ln.asc
			if i > 0 {
				cv.Obj.Mark(bdf.MarkLine, "")
			}
			s.emitLine(t, ln, 0, 0)
			y += ln.asc + ln.desc
		}
	}
	t.letterSpacing(0)
	return cv
}
