package drawio

// Stencils: shapes defined in draw.io's XML stencil libraries
// (mxStencil). Not yet supported: every name resolves to nil.

// stencil is a parsed stencil shape.
type stencil struct {
	// w0, h0 are the stencil's own size (the w and h attributes).
	w0, h0 float64
	// aspect is "fixed" when the stencil keeps its proportions
	// ("variable" otherwise).
	aspect string
}

// stencil returns the stencil for a shape name, or nil.
func (c *converter) stencil(name string) *stencil { return nil }

// labelMargins returns the insets of the label rectangle a stencil's
// labelBounds define (draw.io's mxShape.getLabelMargins override), or nil.
func (st *stencil) labelMargins(s *shape, r rect) *rect { return nil }

// drawShape paints the stencil in the box x, y, w, h (mxStencil.drawShape).
func (st *stencil) drawShape(c *c2d, s *shape, x, y, w, h float64) {}
