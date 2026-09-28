package pdf

import (
	"errors"
	"fmt"
	"math"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// maxPageTreeDepth bounds the page tree nodes above a page, and
// maxRepeatedPages the pages of a document that are pages listed before: a
// page is converted for each place it has in the tree, and a reference to
// it takes six bytes.
const (
	maxPageTreeDepth = 64
	maxRepeatedPages = 16
)

// pageLeaf is a page of the page tree, with the attributes it inherits.
type pageLeaf struct {
	dict  types.Dict
	ref   *types.IndirectRef // nil for a page written in its parent's /Kids
	attrs *model.InheritedPageAttrs
	err   error // an attribute of the page or of a node above it is damaged
}

// pageLeaves returns the pages of the document in order. The tree is walked
// once, with the attributes pdfcpu's PageDict gives a page (the media and
// crop box, the rotation and the resources of the page or of the nearest
// node above it that has them). pdfcpu walks it from the root for every
// page, believes the /Count of the root, and does not end on a node that is
// among its own kids. A node met again is left out, deeper nodes than
// maxPageTreeDepth too, and so are the pages listed again once
// maxRepeatedPages of them are pages: left is how many.
func (p *pdf) pageLeaves() (leaves []pageLeaf, left int, err error) {
	cat, err := p.ctx.Catalog()
	if err != nil {
		return nil, 0, err
	}
	if cat == nil || p.dict(cat["Pages"]) == nil {
		return nil, 0, errors.New("no page tree")
	}
	seen := map[int]bool{}
	listed := map[int]bool{}
	repeated := 0
	var walk func(o types.Object, attrs model.InheritedPageAttrs, err error, depth int)
	walk = func(o types.Object, attrs model.InheritedPageAttrs, err error, depth int) {
		d := p.dict(o)
		if d == nil {
			return
		}
		var ref *types.IndirectRef
		if r, ok := o.(types.IndirectRef); ok {
			ref = &r
		}
		_, hasKids := d["Kids"]
		switch typ := p.name(d["Type"]); {
		case typ == "Pages", typ == "" && hasKids:
			if depth > maxPageTreeDepth {
				return
			}
			if ref != nil {
				if seen[ref.ObjectNumber.Value()] {
					return
				}
				seen[ref.ObjectNumber.Value()] = true
			}
			if e := p.inheritPageAttrs(d, &attrs); e != nil && err == nil {
				err = e
			}
			for _, kid := range p.array(d["Kids"]) {
				walk(kid, attrs, err, depth+1)
			}
		case typ == "Page", typ == "":
			if ref != nil {
				if n := ref.ObjectNumber.Value(); !listed[n] {
					listed[n] = true
				} else if repeated++; repeated > maxRepeatedPages {
					left++
					return
				}
			}
			if e := p.inheritPageAttrs(d, &attrs); e != nil && err == nil {
				err = e
			}
			leaves = append(leaves, pageLeaf{dict: d, ref: ref, attrs: &attrs, err: err})
		}
	}
	walk(cat["Pages"], model.InheritedPageAttrs{}, nil, 0)
	return leaves, left, nil
}

// inheritPageAttrs replaces the attributes that a page tree node or a page
// defines.
func (p *pdf) inheritPageAttrs(d types.Dict, attrs *model.InheritedPageAttrs) error {
	for _, box := range []struct {
		key string
		r   **types.Rectangle
	}{{"MediaBox", &attrs.MediaBox}, {"CropBox", &attrs.CropBox}} {
		o, ok := d[box.key]
		if !ok {
			continue
		}
		a := p.array(o)
		if len(a) < 4 {
			return fmt.Errorf("/%s is not a rectangle", box.key)
		}
		var v [4]float64
		for i := range v {
			if v[i], ok = p.num(a[i]); !ok {
				return fmt.Errorf("/%s is not a rectangle", box.key)
			}
		}
		*box.r = types.NewRectangle(v[0], v[1], v[2], v[3])
	}
	if o, ok := d["Rotate"]; ok {
		v, ok := p.num(o)
		if !ok {
			return errors.New("/Rotate is not a number")
		}
		attrs.Rotate = int(math.Round(v))
	}
	if o, ok := d["Resources"]; ok {
		switch v := p.deref(o).(type) {
		case nil:
			attrs.Resources = nil
		case types.Dict:
			attrs.Resources = v
		default:
			return errors.New("/Resources is not a dictionary")
		}
	}
	return nil
}
