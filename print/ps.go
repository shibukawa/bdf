package print

import (
	"bytes"

	"github.com/shibukawa/bdf/print/raster"
	"github.com/shibukawa/bdf/vector/psbdf"
)

// ps writes the selected pages of the view as a PostScript file, with the
// copies the printer does not make: placed on the paper of the media as
// raster pages are, in gray when the job is not printed in colour.
func (j *job) ps(media string) ([]byte, error) {
	name := media
	if name == "" {
		name = j.p.MediaDefault
	}
	scaling := j.t.Scaling
	if media == "" && scaling == "" {
		scaling = raster.ScaleFit
	}
	dpi := j.t.Resolution
	if dpi < 300 {
		dpi = 300
	}
	opts := &psbdf.Options{Views: []string{j.view.ID}, SkipRoles: j.o.SkipRoles, Pages: j.order(),
		FontFS: j.o.FontFS, FontDirs: j.o.FontDirs, NoSystemFonts: j.o.NoSystemFonts,
		Place: PlacePS(name, scaling, j.t.Orientation), Sides: j.t.Sides, Gray: !colour(j.p, j.t), RasterDPI: dpi, Title: j.t.JobName}
	if j.raw {
		// what IPP asks in attributes, the file asks of the printer
		if n := j.t.copies(); n > 1 {
			opts.Copies, opts.Collate = n, j.t.Collate != "uncollated"
		}
		opts.MediaType = j.t.MediaType
	}
	var buf bytes.Buffer
	res, err := psbdf.Write(&buf, j.doc, opts)
	if err != nil {
		return nil, err
	}
	j.res.Sheets = res.Pages
	j.res.Warnings = append(j.res.Warnings, res.Warnings...)
	return buf.Bytes(), nil
}

// PlacePS returns how the pages of a PostScript file go on the paper of a
// media (a PWG media name), for psbdf.Options.Place: scaled, turned and
// centred as raster pages are (raster.Fit). It is nil for a media that is
// not known: each page is then printed on a sheet of its size.
func PlacePS(media string, scaling raster.Scaling, orientation raster.Orientation) func(w, h float64) (float64, float64, [6]float64) {
	m, err := ParseMedia(media)
	if err != nil {
		return nil
	}
	return func(w, h float64) (float64, float64, [6]float64) {
		// the placement in hundredths of points, which raster.Fit centres
		// on whole ones
		const k = 100
		pl, err := raster.Fit(w, h, raster.Sheet{Width: m.Width, Height: m.Height, DPI: 72 * k}, scaling, orientation)
		if err != nil {
			return 0, 0, [6]float64{}
		}
		s, x, y := pl.Scale/k, pl.X/k, pl.Y/k
		pw, ph := w*s, h*s
		// from the page (its y axis down) to the sheet with its y axis
		// down, as imagebdf places a page
		var t [6]float64
		switch pl.Turns & 3 {
		case 1:
			t = [6]float64{0, s, -s, 0, x + ph, y}
		case 2:
			t = [6]float64{-s, 0, 0, -s, x + pw, y + ph}
		case 3:
			t = [6]float64{0, -s, s, 0, x, y + pw}
		default:
			t = [6]float64{s, 0, 0, s, x, y}
		}
		// and to the sheet of PostScript, its y axis up
		return m.Width, m.Height, [6]float64{t[0], -t[1], t[2], -t[3], t[4], m.Height - t[5]}
	}
}
