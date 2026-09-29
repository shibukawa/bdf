package music

import (
	"fmt"
	"strconv"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
)

const maxTabParts = 4
const maxTabNotes = 20000

type tabPage struct {
	page *bdf.Page
	cv   *canvas.Canvas
}

type builtTab struct {
	view   *bdf.View
	pages  []tabPage
	cues   bdf.Cues
	points [][]cuePoint
	keep   []bool
}

// tabViews draws four fingering alternatives for each of the first four
// pitched parts. Keeping them as BDF views makes the existing score-view
// controls the UI for selecting a playing approach.
func (e *engraver) tabViews(doc *bdf.Document, cvs *canvas.Builder, systems []*system, pageW, pageH float64) []builtTab {
	var out []builtTab
	parts := 0
	for pi, part := range e.s.Parts {
		if part.Percussion {
			continue
		}
		onsets := tabOnsets(e.s, part)
		if len(onsets) == 0 {
			continue
		}
		count := 0
		for _, onset := range onsets {
			count += len(onset.notes)
		}
		if count > maxTabNotes {
			e.warn(fmt.Sprintf("guitar TAB for %s is omitted: more than %d notes", part.Name, maxTabNotes))
			continue
		}
		if parts >= maxTabParts {
			e.warn(fmt.Sprintf("guitar TAB shows the first %d pitched parts; later parts are omitted", maxTabParts))
			break
		}
		parts++
		name := part.Name
		if name == "" {
			name = fmt.Sprintf("Part %d", pi+1)
		}
		for _, preset := range tabPresets {
			positions, omitted := chooseTab(onsets, preset)
			if preset.id == "low" && omitted > 0 {
				e.warn(fmt.Sprintf("guitar TAB for %s leaves %d notes without a playable six-string position", name, omitted))
			}
			if len(positions) == 0 {
				continue
			}
			out = append(out, e.tabView(doc, cvs, systems, pageW, pageH, pi, name, preset, positions))
		}
	}
	return out
}

func (e *engraver) tabView(doc *bdf.Document, cvs *canvas.Builder, systems []*system, pageW, pageH float64,
	partIndex int, name string, preset tabPreset, positions map[*Note]TabPosition) builtTab {
	view := doc.NewView(fmt.Sprintf("guitar-tab-p%d-%s", partIndex+1, preset.id), bdf.ViewFixed,
		fmt.Sprintf("TAB %d: %s · %s (estimated)", partIndex+1, preset.name, name))
	built := builtTab{view: view, points: make([][]cuePoint, len(e.s.Measures))}
	sp := e.sp
	step := 1.35 * sp
	staffH := 5 * step
	blockH := staffH + 6.5*sp
	// A short TAB should read like a compact chart rather than a mostly
	// empty A4 sheet. Longer pieces retain ordinary page breaks.
	tabPageH := min(pageH, max(270, 2*margin+13*sp+float64(len(systems))*blockH))
	y := 0.0
	pageIndex := -1
	var w *writer
	newPage := func() {
		p := view.AddPage(float32(pageW), float32(tabPageH))
		cv := cvs.New()
		cv.Obj.SetBBox(0, 0, p.W, p.H)
		cv.Obj.FillColor(bdf.RGB(255, 255, 255))
		cv.Obj.FillRect(0, 0, p.W, p.H)
		cv.Obj.FillColor(bdf.RGB(0, 0, 0))
		cv.Obj.StrokeColor(bdf.RGB(0, 0, 0))
		cv.Drawn = true
		built.pages = append(built.pages, tabPage{p, cv})
		pageIndex++
		w = &writer{e: e, cv: cv}
		y = margin + 5*sp
		if pageIndex == 0 {
			title := e.text.line(fmt.Sprintf("Guitar TAB — %s — %s", name, preset.name), 3*sp, true, false)
			w.text(title, margin, margin+title.asc, 0)
			detail := e.text.line(preset.explanation+" · estimated fingering", 1.75*sp, false, false)
			w.text(detail, margin, margin+title.asc+3*sp, 0)
			w.text(e.text.line("Standard tuning E A D G B E · other positions may work", 1.75*sp, false, false),
				margin, margin+title.asc+5.2*sp, 0)
			y = margin + 13*sp
		}
	}
	newPage()
	for _, sys := range systems {
		if y+staffH+3*sp > tabPageH-margin {
			newPage()
		}
		x0 := margin + sys.left*sp
		x1 := x0 + sys.width*sp
		w.obj().Mark(bdf.MarkFigure, fmt.Sprintf("Guitar TAB, %s, %s", name, e.systemAlt(sys)))
		for stringNo, label := range []string{"e", "B", "G", "D", "A", "E"} {
			sy := y + float64(stringNo)*step
			w.text(e.text.line(label, 1.8*sp, false, false), x0-2.0*sp, sy+0.35*sp, 0)
			w.draw([]prim{{kind: primLine, x: x0 / sp, y: sy / sp, w: x1 / sp, h: sy / sp, width: 0.08}}, 0, 0)
		}
		w.text(e.text.line(e.measureNumber(sys.first), 1.6*sp, false, false), x0+0.4*sp, y-2*sp, 0)
		for mi := sys.first; mi <= sys.last; mi++ {
			col := e.cols[mi]
			var events []*Event
			if mi < len(e.s.Parts[partIndex].Measures) && e.s.Parts[partIndex].Measures[mi] != nil {
				events = e.s.Parts[partIndex].Measures[mi].Events
			}
			bx := x0 + (col.x+col.width)*sp
			w.draw([]prim{{kind: primLine, x: bx / sp, y: y / sp, w: bx / sp, h: (y + staffH) / sp, width: 0.12}}, 0, 0)
			if mi > sys.first && mi%4 == 0 {
				w.text(e.text.line(e.measureNumber(mi), 1.5*sp, false, false), x0+(col.x+0.3)*sp, y-2*sp, 0)
			}
			for _, sl := range col.slots {
				x := x0 + (col.x+sl.x+0.6)*sp
				for _, event := range events {
					if event == nil || event.Offset != sl.offset || event.Hidden || event.Rest {
						continue
					}
					for _, note := range event.Notes {
						pos, ok := positions[note]
						if !ok {
							continue
						}
						lineY := y + float64(pos.String-1)*step
						label := e.text.line(strconv.Itoa(pos.Fret), 1.85*sp, true, false)
						w.obj().FillColor(bdf.RGB(255, 255, 255))
						w.obj().FillRect(float32(x-label.width/2-0.2*sp), float32(lineY-0.8*sp),
							float32(label.width+0.4*sp), float32(1.2*sp))
						w.obj().FillColor(bdf.RGB(0, 0, 0))
						w.text(label, x, lineY+0.38*sp, 1)
					}
				}
				built.points[mi] = append(built.points[mi], cuePoint{sl.offset, uint32(len(built.cues.Systems)), x})
			}
			built.points[mi] = append(built.points[mi], cuePoint{e.s.Measures[mi].Length, uint32(len(built.cues.Systems)), bx})
		}
		w.obj().Mark(bdf.MarkEnd, "")
		built.cues.Systems = append(built.cues.Systems, bdf.CueSystem{Page: uint32(pageIndex),
			X: float32(x0), Y: float32(y - sp), W: float32(x1 - x0), H: float32(staffH + 2*sp)})
		y += blockH
	}
	// Page selection applies independently to this view's page count.
	built.keep = make([]bool, len(view.Pages))
	if e.o.Pages == nil {
		for i := range built.keep {
			built.keep[i] = true
		}
	} else {
		var selected []*bdf.Page
		for _, n := range e.o.Pages.Numbers(len(view.Pages)) {
			if n >= 1 && n <= len(view.Pages) {
				selected = append(selected, view.Pages[n-1])
				built.keep[n-1] = true
			}
		}
		view.Pages = selected
	}
	return built
}
