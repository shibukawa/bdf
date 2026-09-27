// Package xmp reads the Dublin Core of an XMP metadata packet (ISO 16684-1),
// as image files embed it.
package xmp

import (
	"bytes"
	"encoding/xml"
	"slices"
	"strings"

	"github.com/shibukawa/bdf"
)

const (
	rdfNS = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	dcNS  = "http://purl.org/dc/elements/1.1/"
	xmpNS = "http://ns.adobe.com/xap/1.0/"
)

// DublinCore returns the dc: properties of a packet, with xmp:CreateDate and
// xmp:ModifyDate as created and modified. A language alternative gives its
// x-default value, or its first. A malformed packet gives what was read
// before the error.
func DublinCore(packet []byte) bdf.DublinCore {
	var dc bdf.DublinCore
	field := func(name xml.Name) *bdf.DCValues {
		switch {
		case name.Space == dcNS && name.Local != "created" && name.Local != "modified":
			return dc.Field(name.Local)
		case name.Space == xmpNS && name.Local == "CreateDate":
			return &dc.Created
		case name.Space == xmpNS && name.Local == "ModifyDate":
			return &dc.Modified
		}
		return nil
	}
	d := xml.NewDecoder(bytes.NewReader(packet))
	d.Strict = false
	for {
		tok, err := d.Token()
		if err != nil {
			return dc
		}
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		if start.Name.Space == rdfNS && start.Name.Local == "Description" {
			// Simple properties may be attributes of the description.
			for _, a := range start.Attr {
				if f := field(a.Name); f != nil && len(*f) == 0 {
					if v := strings.TrimSpace(a.Value); v != "" {
						*f = bdf.DCValues{v}
					}
				}
			}
			continue
		}
		if f := field(start.Name); f != nil {
			if vals := property(d); len(*f) == 0 && len(vals) > 0 {
				*f = vals
			}
		}
	}
}

// property reads the value of the property element just started: its text,
// or the items of the rdf:Alt, rdf:Seq or rdf:Bag in it.
func property(d *xml.Decoder) bdf.DCValues {
	var vals bdf.DCValues
	var text, li strings.Builder
	depth, inLi, alt, def := 0, false, false, -1
	for {
		tok, err := d.Token()
		if err != nil {
			return vals
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			if t.Name.Space != rdfNS {
				continue
			}
			switch t.Name.Local {
			case "Alt":
				alt = true
			case "li":
				inLi = true
				li.Reset()
				if slices.ContainsFunc(t.Attr, func(a xml.Attr) bool { return a.Name.Local == "lang" && a.Value == "x-default" }) {
					def = len(vals)
				}
			}
		case xml.EndElement:
			if depth == 0 {
				if len(vals) == 0 {
					if s := strings.TrimSpace(text.String()); s != "" {
						vals = bdf.DCValues{s}
					}
				}
				if alt && len(vals) > 1 {
					if def < 0 {
						def = 0
					}
					vals = vals[def : def+1]
				}
				return vals
			}
			if inLi && t.Name.Space == rdfNS && t.Name.Local == "li" {
				if s := strings.TrimSpace(li.String()); s != "" {
					vals = append(vals, s)
				} else if def == len(vals) {
					def = -1
				}
				inLi = false
			}
			depth--
		case xml.CharData:
			if inLi {
				li.Write(t)
			} else if depth == 0 {
				text.Write(t)
			}
		}
	}
}
