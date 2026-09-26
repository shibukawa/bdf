package drawingml

import (
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/metafile"
)

// Windows metafiles (EMF and WMF) are replayed into the object instead of
// being stored as images, which browsers cannot decode (see package
// metafile).

// drawMetafile replays a metafile into the rectangle x,y,w,h with the
// picture effects rc (may be nil).
func (s *Drawing) drawMetafile(cv *canvas.Canvas, data []byte, rc *recolor, x, y, w, h float64) {
	c := s.c
	opts := &metafile.Options{Doc: c.doc, Fonts: c.fonts, Images: c.imgOpts,
		Warn: func(msg string) { c.warnOnce("mf:"+msg, "%s", msg) }}
	if rc != nil {
		opts.Recolor = rc
	}
	metafile.Draw(cv, data, x, y, w, h, opts)
}
