package pdf

import (
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/shibukawa/bdf"
)

// softMask is an ExtGState /SMask: a transparency group whose alpha or
// luminosity multiplies the opacity of what is painted while it is set.
type softMask struct {
	sd       *types.StreamDict // the group XObject /G
	key      string
	res      types.Dict // resources where it was set (for a group without its own)
	kind     byte       // bdf.MaskAlpha or bdf.MaskLuminosity
	backdrop bdf.Color  // /BC for a luminosity mask
	transfer []byte     // /TR sampled at 256 mask values; nil for the identity
	ctm      matrix     // the CTM when it was set: the mask's coordinate system
	bbox     rect       // the group's bounding box in object space
}

// maskGroup is the GROUP_BEGIN opened for what is painted under a soft
// mask. Painting operations are masked one by one in PDF; the group masks
// everything painted at its level while the mask stays set, which is the
// same for the usual q /GS gs … Do Q.
type maskGroup struct {
	sm    *softMask
	depth int     // len(in.stack) when it was opened
	ctm   matrix  // CTM at GROUP_BEGIN
	em    emitted // emitted state of the enclosing context
}

// loadSoftMask reads a soft mask dictionary; the mask takes the current CTM.
func (in *interp) loadSoftMask(v types.Object) *softMask {
	p := in.c.pdf
	d := p.dict(v)
	g := d["G"]
	sd := p.stream(g)
	if sd == nil {
		in.c.warnOnce("smask-g", "soft mask without a group XObject ignored")
		return nil
	}
	sm := &softMask{sd: sd, key: objKey(g), res: in.res, kind: bdf.MaskAlpha, ctm: in.gs.ctm}
	if p.name(d["S"]) == "Luminosity" {
		sm.kind = bdf.MaskLuminosity
		sm.backdrop = bdf.RGB(0, 0, 0)
		if bc := p.nums(d["BC"]); len(bc) > 0 {
			var cs *colorSpace
			if csObj := p.dict(sd.Dict["Group"])["CS"]; csObj != nil {
				cs = in.c.loadColorSpace(csObj, in.res)
			}
			if cs == nil || cs.n != len(bc) {
				switch len(bc) {
				case 1:
					cs = csGray
				case 4:
					cs = csCMYK
				default:
					cs = csRGB
				}
			}
			sm.backdrop = cs.rgb(bc)
		}
	}
	if tr := d["TR"]; tr != nil && p.name(tr) != "Identity" {
		// Chrome inverts masks with {1 exch sub}.
		if fn := p.loadFunction(tr); fn != nil {
			sm.transfer = make([]byte, 256)
			ident := true
			for i := range sm.transfer {
				v := 0.0
				if out := fn.eval(float64(i) / 255); len(out) > 0 {
					v = out[0]
				}
				sm.transfer[i] = to8(v)
				ident = ident && sm.transfer[i] == byte(i)
			}
			if ident {
				sm.transfer = nil
			}
		} else {
			in.c.warnOnce("smask-tr", "soft mask transfer function not readable; ignored")
		}
	}
	mat := p.matrixOr(sd.Dict["Matrix"], identity)
	sm.bbox = mat.mul(sm.ctm).transformRect(p.rectOr(sd.Dict["BBox"], rect{}))
	return sm
}

// setSoftMask changes the soft mask; a group open for the old one at this
// level ends here. (A change at a deeper level leaves the outer group open,
// so what follows is masked by both.)
func (in *interp) setSoftMask(sm *softMask) {
	if n := len(in.masks); n > 0 && in.masks[n-1].depth == len(in.stack) {
		in.closeMask(false)
	}
	in.gs.softMask = sm
}

// enterMask opens the group for the soft mask in effect before a painting
// operation, unless it is already open.
func (in *interp) enterMask() {
	sm := in.gs.softMask
	if sm == nil {
		return
	}
	if n := len(in.masks); n > 0 && in.masks[n-1].sm == sm {
		return
	}
	inv, ok := in.gs.ctm.inverse()
	if !ok {
		return
	}
	in.closeTextBlock()
	bb := in.p.bbox
	r := rect{float64(bb.X), float64(bb.Y), float64(bb.X + bb.W), float64(bb.Y + bb.H)}.intersect(in.gs.clip)
	if sm.outside() == 0 {
		// Outside the group's box the mask is 0.
		r = r.intersect(sm.bbox)
	}
	if r.empty() {
		r = rect{}
	}
	u := inv.transformRect(r)
	in.masks = append(in.masks, maskGroup{sm: sm, depth: len(in.stack), ctm: in.gs.ctm, em: in.em})
	in.obj.GroupBegin(1, in.gs.blend, float32(u.x0), float32(u.y0), float32(u.x1-u.x0), float32(u.y1-u.y0))
	// The group's canvas starts from the initial state; force every state out again.
	in.em = initialEmitted()
}

// closeMask draws the innermost group's mask and ends the group.
// restoring says that a RESTORE follows, which discards the transform.
func (in *interp) closeMask(restoring bool) {
	n := len(in.masks)
	if n == 0 {
		return
	}
	mg := in.masks[n-1]
	in.masks = in.masks[:n-1]
	in.closeTextBlock()
	sm := mg.sm
	p := in.c.pdf
	in.obj.MaskBegin(sm.kind, sm.backdrop, sm.transfer)
	if child := in.c.formObject(sm.key, sm.sd, sm.res, in.depth+1, nil); child != nil {
		if inv, ok := in.gs.ctm.inverse(); ok {
			if rel := sm.ctm.mul(inv); !rel.isIdentity() {
				in.obj.Transform(float32(rel[0]), float32(rel[1]), float32(rel[2]), float32(rel[3]), float32(rel[4]), float32(rel[5]))
			}
			if mat := p.matrixOr(sm.sd.Dict["Matrix"], identity); !mat.isIdentity() {
				in.obj.Transform(float32(mat[0]), float32(mat[1]), float32(mat[2]), float32(mat[3]), float32(mat[4]), float32(mat[5]))
			}
			if bbox := p.rectOr(sm.sd.Dict["BBox"], rect{}); !bbox.empty() {
				in.obj.ClipRect(float32(bbox.x0), float32(bbox.y0), float32(bbox.x1-bbox.x0), float32(bbox.y1-bbox.y0))
			}
			in.obj.Use(in.p.childRef(child))
		}
	}
	in.obj.MaskEnd()
	in.obj.GroupEnd()
	in.em = mg.em
	if !restoring && in.gs.ctm != mg.ctm {
		// A cm inside the group changed the transform at its level; the
		// enclosing context still has the one of GROUP_BEGIN. (A clip set
		// there is lost.)
		if inv, ok := mg.ctm.inverse(); ok {
			d := in.gs.ctm.mul(inv)
			in.obj.Transform(float32(d[0]), float32(d[1]), float32(d[2]), float32(d[3]), float32(d[4]), float32(d[5]))
		}
	}
}

// outside is the mask value outside the group's box: the transparent or
// backdrop value through the transfer function.
func (sm *softMask) outside() byte {
	v := byte(0)
	if sm.kind == bdf.MaskLuminosity {
		v = to8(luminosity(sm.backdrop))
	}
	if sm.transfer != nil {
		v = sm.transfer[v]
	}
	return v
}

// luminosity is the luminosity of a colour with the weights of the PDF
// non-separable blend modes, as the renderer computes it for a mask.
func luminosity(c bdf.Color) float64 {
	r, g, b := float64(uint8(c>>24)), float64(uint8(c>>16)), float64(uint8(c>>8))
	return (0.3*r + 0.59*g + 0.11*b) / 255
}
