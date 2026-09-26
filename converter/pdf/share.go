package pdf

import (
	"crypto/sha256"
	"fmt"

	"github.com/shibukawa/bdf"
)

// Pages exported from one master (a slide layout, a letterhead) begin with
// the same instructions, but PDF has already expanded them into every page.
// Identical objects dedupe by hash on their own; this recovers the shared
// part of pages that differ later: the longest common instruction prefix
// of two or more page bodies becomes one shared object that each page
// USEs (bdf.SharePrefixes).

// sharePrefixMinBytes is the smallest op-stream prefix worth sharing: a USE
// plus the state re-emitted after it costs a few dozen bytes per page.
const sharePrefixMinBytes = 512

// sharePagePrefixes rewrites the page bodies that share a prefix. Fonts and
// child objects are still placeholders at this point (finalize fills them
// in), so they are first made distinct per resource for the content
// comparison, then mapped back to the pending entries of the new objects.
func (c *converter) sharePagePrefixes() {
	if len(c.pageBodies) < 2 {
		return
	}
	objs := make([]*bdf.Object, len(c.pageBodies))
	for i, pr := range c.pageBodies {
		p := pr.body
		for ref, f := range p.fonts {
			p.obj.UpdateFont(ref, bdf.SystemFont(fmt.Sprintf("\x00pdf2bdf:font:%p", f), 0, 0))
		}
		for ref, child := range p.children {
			p.obj.UpdateObject(ref, pointerHash(child))
		}
		objs[i] = p.obj
	}
	shares, err := bdf.SharePrefixes(objs, sharePrefixMinBytes)
	if err != nil {
		c.warnf("prefix sharing: %v", err)
		return
	}
	replaced := map[*pending]bool{}
	for _, sh := range shares {
		first := c.pageBodies[sh.Members[0].Index].body
		prefix := c.rebindPending(sh.Prefix, first, sh.PrefixFonts, sh.PrefixObjects, nil)
		c.pendings = append(c.pendings, prefix)
		for _, m := range sh.Members {
			old := c.pageBodies[m.Index].body
			rest := c.rebindPending(m.Rest, old, m.RestFonts, m.RestObjects, prefix)
			c.pageBodies[m.Index].body = rest
			c.pendings = append(c.pendings, rest)
			replaced[old] = true
			c.sharedBytes += m.Saved
		}
		c.sharedPrefixes++
	}
	if len(replaced) > 0 {
		kept := c.pendings[:0]
		for _, p := range c.pendings {
			if !replaced[p] {
				kept = append(kept, p)
			}
		}
		c.pendings = kept
	}
}

// rebindPending wraps a rewritten object in a pending whose font and child
// maps follow the new references (fonts[i] / objects[i] give the reference
// in old that the new reference i came from; -1 is the prefix placeholder).
func (c *converter) rebindPending(obj *bdf.Object, old *pending, fonts []bdf.FontRef, objects []int, prefix *pending) *pending {
	p := &pending{obj: obj, fonts: map[bdf.FontRef]*pdfFont{}, fontRefs: map[*pdfFont]bdf.FontRef{}, children: map[bdf.ObjRef]*pending{}, bbox: old.bbox}
	for i, oldRef := range fonts {
		f := old.fonts[oldRef]
		p.fonts[bdf.FontRef(i)] = f
		p.fontRefs[f] = bdf.FontRef(i)
	}
	for i, oldRef := range objects {
		if oldRef < 0 {
			p.children[bdf.ObjRef(i)] = prefix
		} else {
			p.children[bdf.ObjRef(i)] = old.children[bdf.ObjRef(oldRef)]
		}
	}
	return p
}

// pointerHash is a placeholder hash unique to a pending object.
func pointerHash(p *pending) bdf.Hash {
	sum := sha256.Sum256([]byte(fmt.Sprintf("pdf2bdf:pending:%p", p)))
	var h bdf.Hash
	copy(h[:], sum[:])
	return h
}
