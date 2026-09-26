package ooxml

import (
	"strings"

	"github.com/shibukawa/bdf"
)

// CoreProperties reads the core properties of the package, which are mostly
// Dublin Core already (ECMA-376 Part 2 §11); the keywords become subjects.
func (p *Package) CoreProperties() bdf.DublinCore {
	var dc bdf.DublinCore
	r, ok := p.RelOfType("", "/core-properties")
	if !ok {
		return dc
	}
	n, err := p.XML(r.Target)
	if err != nil {
		return dc
	}
	for _, e := range []struct {
		f    *bdf.DCValues
		name string
	}{
		{&dc.Title, "title"}, {&dc.Creator, "creator"}, {&dc.Subject, "subject"}, {&dc.Description, "description"},
		{&dc.Identifier, "identifier"}, {&dc.Language, "language"}, {&dc.Created, "created"}, {&dc.Modified, "modified"},
	} {
		for _, k := range n.Children(e.name) {
			if s := strings.TrimSpace(k.Content()); s != "" {
				*e.f = append(*e.f, s)
			}
		}
	}
	dc.Subject = append(dc.Subject, bdf.SplitKeywords(n.Child("keywords").Content())...)
	return dc
}
