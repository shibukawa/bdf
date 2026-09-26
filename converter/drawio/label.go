package drawio

// labelBox is a laid-out label.
type labelBox struct {
	box rect
}

// layoutLabel lays out the label of a cell; nil when it has none.
func (c *converter) layoutLabel(st *cellState, s *shape) *labelBox { return nil }

// extent is the rectangle the label paints over.
func (l *labelBox) extent() rect { return l.box }

// drawLabel draws a laid-out label.
func (c *converter) drawLabel(c2 *c2d, l *labelBox) {}
