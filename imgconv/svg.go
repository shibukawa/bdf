package imgconv

import (
	"bytes"
	"encoding/xml"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// SVG images are stored as they are and drawn by the browser (docs/spec.md
// §6.2); what converters need is their size, which the root element's
// width, height and viewBox give.

// SVGSize is the size an SVG image says it has.
type SVGSize struct {
	// W and H are the natural width and height in CSS px (0: none). One that
	// is missing follows from the other and the view box's proportions.
	W, H float64
	// ViewW and ViewH are the width and height of the view box (0: none).
	ViewW, ViewH float64
}

// ParseSVGSize reads the width, height and viewBox attributes of an SVG
// root element. Relative lengths (percentages) are not a natural size.
func ParseSVGSize(width, height, viewBox string) SVGSize {
	var s SVGSize
	s.W, _ = CSSLength(width)
	s.H, _ = CSSLength(height)
	if f := strings.FieldsFunc(viewBox, func(r rune) bool { return r == ' ' || r == ',' || r == '\t' || r == '\n' || r == '\r' }); len(f) == 4 {
		vw, err1 := strconv.ParseFloat(f[2], 64)
		vh, err2 := strconv.ParseFloat(f[3], 64)
		if err1 == nil && err2 == nil && vw > 0 && vh > 0 {
			s.ViewW, s.ViewH = vw, vh
		}
	}
	if s.ViewW > 0 {
		switch {
		case s.W > 0 && s.H == 0:
			s.H = s.W * s.ViewH / s.ViewW
		case s.H > 0 && s.W == 0:
			s.W = s.H * s.ViewW / s.ViewH
		}
	}
	return s
}

// Ratio returns the natural aspect ratio (width / height): the view box's,
// else that of the natural size; 0 when there is none.
func (s SVGSize) Ratio() float64 {
	switch {
	case s.ViewW > 0:
		return s.ViewW / s.ViewH
	case s.W > 0 && s.H > 0:
		return s.W / s.H
	}
	return 0
}

// Pixels returns the size of the image's pixels in CSS px (docs/spec.md
// §6.2): the natural size, else the view box's size, else 300 × 150 (the
// default size of replaced elements) for what is missing.
func (s SVGSize) Pixels() (w, h float64) {
	switch {
	case s.W > 0 && s.H > 0:
		return s.W, s.H
	case s.W == 0 && s.H == 0 && s.ViewW > 0:
		return s.ViewW, s.ViewH
	}
	w, h = s.W, s.H
	if w == 0 {
		w = 300
	}
	if h == 0 {
		h = 150
	}
	return w, h
}

// SVGRoot returns the size attributes of an SVG document's root element;
// ok is false when the document's first element is not svg.
func SVGRoot(data []byte) (width, height, viewBox string, ok bool) {
	d := xml.NewDecoder(bytes.NewReader(bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))))
	d.Strict = false
	d.CharsetReader = func(_ string, r io.Reader) (io.Reader, error) { return r, nil } // the attributes read are ASCII
	for {
		tok, err := d.Token()
		if err != nil {
			return "", "", "", false
		}
		if s, isStart := tok.(xml.StartElement); isStart {
			if s.Name.Local != "svg" {
				return "", "", "", false
			}
			for _, a := range s.Attr {
				if a.Name.Space != "" {
					continue
				}
				switch a.Name.Local {
				case "width":
					width = a.Value
				case "height":
					height = a.Value
				case "viewBox":
					viewBox = a.Value
				}
			}
			return width, height, viewBox, true
		}
	}
}

// lengthRE matches a CSS length: a number and a unit.
var lengthRE = regexp.MustCompile(`^\s*([+-]?(?:\d+\.?\d*|\.\d+)(?:[eE][+-]?\d+)?)\s*([A-Za-z%]*)\s*$`)

// CSSLength reads an absolute CSS length in CSS px: a number alone is px,
// and em and ex are those of the default font (16 px). Percentages,
// lengths that are not positive and what it cannot read are not lengths.
func CSSLength(s string) (float64, bool) {
	m := lengthRE.FindStringSubmatch(s)
	if m == nil {
		return 0, false
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil || v <= 0 {
		return 0, false
	}
	switch strings.ToLower(m[2]) {
	case "", "px":
	case "pt":
		v *= 96.0 / 72
	case "pc":
		v *= 16
	case "in":
		v *= 96
	case "cm":
		v *= 96 / 2.54
	case "mm":
		v *= 96 / 25.4
	case "q":
		v *= 96 / 101.6
	case "em":
		v *= 16
	case "ex":
		v *= 8
	default:
		return 0, false
	}
	return v, true
}
