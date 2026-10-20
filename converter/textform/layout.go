package textform

import (
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/form/textpart"
	"github.com/shibukawa/bdf/internal/fontdb"
)

// MonoFamilies are the monospace fonts that the pages are drawn with when
// the template names none, in the order they are looked for: BIZ UDGothic
// (of Google Fonts, under the SIL Open Font License), IPAGothic, Noto Sans
// Mono CJK JP, and those of the systems. Without any of them the pages are
// drawn with what the fonts at hand have; the characters stay in their
// columns all the same.
var MonoFamilies = []string{"BIZ UDGothic", "BIZ UDゴシック", "IPAGothic", "IPAゴシック", "Noto Sans Mono CJK JP",
	"TakaoGothic", "VL Gothic", "MS Gothic", "Osaka-Mono", "Noto Sans Mono", "DejaVu Sans Mono", "Liberation Mono", "Courier New", "monospace"}

// baseline is where the base line of a line is, from its top, in its
// height: that of the em box of the ideographs of CJK fonts, whose box
// drawing characters then join from line to line.
const baseline = 0.88

// families returns the fonts that the pages are drawn with: the one the
// template names first, then the first of MonoFamilies that the fonts at
// hand have, then the others.
func families(db *fontdb.DB, named string) []string {
	var out []string
	if named != "" {
		out = append(out, named)
	}
	for _, f := range MonoFamilies {
		if len(db.Family(f)) > 0 {
			out = append(out, f)
			break
		}
	}
	return append(out, MonoFamilies...)
}

// run is characters drawn together: of one font, size and layer, and of
// the same width in the font for each of their columns, so that the
// advance of the run, which is that of its columns, puts each of them in
// its own.
type run struct {
	text        string
	col, cols   int
	paper, tall bool
	face        *fontset.Choice
	// perCol is the width of the characters in the font, in ems, for each
	// of their columns
	perCol float64
}

// draw lays the pages of the print image out on a fixed view: each
// character in its columns of the grid, the ranges of the paper on a layer
// of the role "background". tables gives the table of each line of each
// page (0 for none).
func draw(doc *bdf.Document, view *bdf.View, im *textpart.Image, tables [][]int, fams []string, set *fontset.Set, embed func()) {
	cw, rh := 72/im.CPI, 72/im.LPI
	x0, y0 := float64(im.Margin.Left)*cw, float64(im.Margin.Top)*rh
	w := float64(im.Columns) * cw
	cvs := canvas.NewBuilder(doc, set)
	type layers struct {
		page        *bdf.Page
		body, paper *canvas.Canvas
	}
	var all []layers
	for pi, p := range im.Pages {
		h := float64(im.Lines) * rh
		if im.Lines == 0 {
			h = float64(im.Margin.Top+max(1, len(p.Lines))) * rh
		}
		pg := view.AddPage(float32(w), float32(h))
		ls := layers{page: pg, body: cvs.New(), paper: cvs.New()}
		for _, cv := range []*canvas.Canvas{ls.body, ls.paper} {
			cv.Obj.SetBBox(0, 0, float32(w), float32(h))
			cv.Obj.FillColor(bdf.RGB(0, 0, 0))
			cv.Obj.TextStyle(bdf.AlignLeft, bdf.BaselineAlphabetic, bdf.DirLTR, 0)
		}
		table, cellRow := 0, 0
		for li, l := range p.Lines {
			runs := runsOf(im, l, fams, set)
			if len(runs) == 0 {
				continue
			}
			row := li + 1
			if l.Tall() {
				row++
			}
			y := y0 + float64(row)*rh - (1-baseline)*rh
			var tid int
			if pi < len(tables) && li < len(tables[pi]) {
				tid = tables[pi][li]
			}
			body, paper := false, false
			for _, r := range runs {
				body, paper = body || !r.paper, paper || r.paper
			}
			if body {
				o := ls.body.Obj
				if tid != table {
					if table != 0 {
						o.Mark(bdf.MarkEnd, "")
					}
					if table, cellRow = tid, 0; tid != 0 {
						o.Mark(bdf.MarkTable, "")
					}
				}
				if table != 0 {
					cellRow++
					o.Mark(bdf.MarkCell, fmt.Sprintf("A%d", cellRow))
				} else {
					o.Mark(bdf.MarkParagraph, "")
				}
			}
			if paper {
				ls.paper.Obj.Mark(bdf.MarkParagraph, "")
			}
			for _, r := range runs {
				cv := ls.body
				if r.paper {
					cv = ls.paper
				}
				size := rh
				if r.tall {
					size *= 2
				}
				cv.Obj.Font(cv.Font(r.face.Use), float32(size))
				cv.Obj.FillText(r.text, float32(x0+float64(r.col)*cw), float32(y), float32(float64(r.cols)*cw))
				cv.Drawn = true
			}
		}
		if table != 0 {
			ls.body.Obj.Mark(bdf.MarkEnd, "")
		}
		all = append(all, ls)
	}
	embed()
	cvs.Encode()
	for _, ls := range all {
		if ls.paper.Drawn {
			ls.page.Layers = append(ls.page.Layers, bdf.Layer{Role: bdf.RoleBackground, Obj: ls.paper.Hash()})
		}
		if ls.body.Drawn || !ls.paper.Drawn {
			ls.page.Layers = append(ls.page.Layers, bdf.Layer{Role: bdf.RoleBody, Obj: ls.body.Hash()})
		}
	}
}

// runsOf splits a line into the runs it is drawn in. A space is drawn in
// the run of the characters before it when it is as wide as they are for
// its columns; spaces that are not, and those at the start of the line,
// are not drawn (they are room between runs).
func runsOf(im *textpart.Image, l textpart.Line, fams []string, set *fontset.Set) []run {
	var out []run
	var cur *run
	for _, c := range im.Cells(l) {
		if c.Cols == 0 {
			if cur != nil {
				cur.text += c.Text
			}
			continue
		}
		r, _ := utf8.DecodeRuneInString(c.Text)
		paper, tall := c.Paper, c.Size == textpart.SizeTall || c.Size == textpart.SizeBig
		joins := func(face *fontset.Choice) bool {
			return cur != nil && cur.face == face && cur.paper == paper && cur.tall == tall && cur.col+cur.cols == c.Col &&
				math.Abs(set.Advance(face, r)/float64(c.Cols)-cur.perCol) < 1e-6
		}
		if r == ' ' {
			if cur != nil && joins(cur.face) {
				cur.text += c.Text
				cur.cols += c.Cols
				continue
			}
			cur = nil
			continue
		}
		face := set.FaceForFamilies(fams, false, false, r)
		if face.Loaded != nil && face.Loaded.Has(r) && joins(face) {
			cur.text += c.Text
			cur.cols += c.Cols
			continue
		}
		out = append(out, run{text: c.Text, col: c.Col, cols: c.Cols, paper: paper, tall: tall, face: face, perCol: set.Advance(face, r) / float64(c.Cols)})
		cur = &out[len(out)-1]
		// (a character that no font has is a run of its own: the advance
		// of a run would stretch the others over its room)
		if face.Loaded == nil || !face.Loaded.Has(r) {
			cur = nil
		}
	}
	return out
}
