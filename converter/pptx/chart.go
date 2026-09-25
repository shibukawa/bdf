package pptx

// drawChart draws a chart part.
func (s *slideCtx) drawChart(cv *canvas, sh *shape, xf xform, rid string) {
	s.c.warnOnce("chart", "charts are not supported yet")
}
