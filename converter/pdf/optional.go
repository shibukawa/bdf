package pdf

import (
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// Optional content (PDF 32000-1 §8.11): the layers of a PDF, which viewers
// let people show and hide. A BDF document has no such switch, so the
// converter draws what a viewer shows when it opens the document: the
// groups the default configuration (/OCProperties /D) turns on. Content in
// a hidden group, and XObjects and annotations that belong to one, are left
// out. Illustrator writes its layers this way, hidden ones included.
type optionalContent struct {
	p      *pdf
	baseOn bool
	state  map[int]bool // object number of a group → on, for the groups /ON and /OFF list
}

// newOptionalContent reads the default configuration of the catalog's
// /OCProperties; nil when the document has none (everything is shown).
func newOptionalContent(p *pdf, cat types.Dict) *optionalContent {
	props := p.dict(cat["OCProperties"])
	if props == nil {
		return nil
	}
	d := p.dict(props["D"])
	oc := &optionalContent{p: p, baseOn: p.name(d["BaseState"]) != "OFF", state: map[int]bool{}}
	for _, list := range []struct {
		key string
		on  bool
	}{{"ON", true}, {"OFF", false}} {
		for _, g := range p.array(d[list.key]) {
			if r, ok := g.(types.IndirectRef); ok {
				oc.state[r.ObjectNumber.Value()] = list.on
			}
		}
	}
	return oc
}

// visible reports whether content that belongs to o (an optional content
// group or membership dictionary, as a /OC entry or a marked-content
// property list gives it) is shown.
func (oc *optionalContent) visible(o types.Object) bool {
	if oc == nil || o == nil {
		return true
	}
	d := oc.p.dict(o)
	if d == nil {
		return true
	}
	switch oc.p.name(d["Type"]) {
	case "OCG":
		return oc.groupOn(o)
	case "OCMD":
		if v, ok := oc.expression(d["VE"], 0); ok {
			return v
		}
		var groups []types.Object
		if a := oc.p.array(d["OCGs"]); a != nil {
			groups = a
		} else if d["OCGs"] != nil {
			groups = []types.Object{d["OCGs"]}
		}
		var on, off int
		for _, g := range groups {
			if oc.p.dict(g) == nil {
				continue // null or deleted groups are ignored
			}
			if oc.groupOn(g) {
				on++
			} else {
				off++
			}
		}
		if on+off == 0 {
			return true
		}
		switch oc.p.name(d["P"]) {
		case "AllOn":
			return off == 0
		case "AnyOff":
			return off > 0
		case "AllOff":
			return on == 0
		}
		return on > 0 // AnyOn
	}
	return true
}

// groupOn returns the state of a group in the default configuration.
func (oc *optionalContent) groupOn(o types.Object) bool {
	if r, ok := o.(types.IndirectRef); ok {
		if on, listed := oc.state[r.ObjectNumber.Value()]; listed {
			return on
		}
	}
	return oc.baseOn
}

// expression evaluates a visibility expression (/VE, PDF 1.6): an array of
// /And, /Or or /Not and operands that are groups or expressions. ok is
// false when there is no expression or it is malformed.
func (oc *optionalContent) expression(o types.Object, depth int) (v, ok bool) {
	a := oc.p.array(o)
	if len(a) < 2 || depth > 16 {
		return false, false
	}
	var vals []bool
	for _, e := range a[1:] {
		if sub := oc.p.array(e); sub != nil {
			if v, ok := oc.expression(sub, depth+1); ok {
				vals = append(vals, v)
			}
		} else if oc.p.dict(e) != nil {
			vals = append(vals, oc.groupOn(e))
		}
	}
	if len(vals) == 0 {
		return false, false
	}
	switch oc.p.name(a[0]) {
	case "Not":
		return !vals[0], true
	case "And":
		for _, v := range vals {
			if !v {
				return false, true
			}
		}
		return true, true
	case "Or":
		for _, v := range vals {
			if v {
				return true, true
			}
		}
		return false, true
	}
	return false, false
}
