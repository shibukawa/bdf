package music

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/converter/internal/canvas"
	"github.com/shibukawa/bdf/converter/internal/fontset"
	"github.com/shibukawa/bdf/internal/fontdb"
)

// Page geometry defaults.
const (
	a4W       = 595.2756 // pt
	a4H       = 841.8898
	mm        = 72 / 25.4
	margin    = 14 * mm
	defaultSp = 1.75 * mm // a 7 mm staff
)

// maxNotes bounds the notes engraved (hostile inputs).
const maxNotes = 2_000_000

// maxStaves bounds the staves of a system.
const maxStaves = 4096

// maxCues bounds the cues of the playing position: repeats play the
// measures, and so their cues, many times. It is a variable for the tests.
var maxCues = 4_000_000

func build(s *Score, o *Options) (*Result, error) {
	res := &Result{}
	warn := func(msg string) {
		if o.Warn != nil {
			o.Warn(msg)
		} else {
			res.Warnings = append(res.Warnings, msg)
		}
	}
	if s.err != nil {
		return nil, s.err
	}
	if len(s.Parts) == 0 || len(s.Measures) == 0 {
		return nil, errors.New("music: the score has no notes")
	}
	staves := 0
	for _, p := range s.Parts {
		staves += max(p.Staves, 1)
	}
	if staves > maxStaves {
		return nil, errStaves()
	}
	if staves*len(s.Measures) > maxStaffMeasures {
		return nil, errStaffMeasures(staves, "staves", len(s.Measures))
	}
	count := 0
	for _, p := range s.Parts {
		for _, m := range p.Measures {
			if m != nil {
				count += len(m.Events)
			}
		}
	}
	if count > maxNotes {
		return nil, fmt.Errorf("music: the score has more than %d notes", maxNotes)
	}
	for _, w := range s.warnings {
		warn(w)
	}
	pageW, pageH := o.PageWidth, o.PageHeight
	if pageW <= 0 || pageH <= 0 {
		pageW, pageH = a4W, a4H
	}
	sp := o.StaffSpace
	if sp <= 0 {
		sp = defaultSp
	}

	doc := bdf.NewDocument()
	doc.Meta.Source = o.Source
	title := s.Title
	if o.Title != "" {
		title = o.Title
	}
	if title != "" {
		doc.Meta.DC.Title = bdf.DCValues{title}
	}
	if s.Composer != "" {
		doc.Meta.DC.Creator = bdf.DCValues{s.Composer}
	}
	for _, c := range []string{s.Lyricist, s.Arranger} {
		if c != "" {
			doc.Meta.DC.Contributor = append(doc.Meta.DC.Contributor, c)
		}
	}
	if s.Rights != "" {
		doc.Meta.DC.Rights = bdf.DCValues{s.Rights}
	}
	db := fontdb.New(o.FontFS, o.FontDirs, !o.NoSystemFonts)
	set := fontset.New(db, warn)

	e := &engraver{s: s, o: o, text: &texter{set: set}, warn: warn}
	e.prepare()
	// a system must fit a page: smaller staves for many of them
	usableH := pageH - 2*margin
	if n := float64(len(e.staves)); (n*4+(n-1)*6.5+10)*sp > usableH {
		sp = max(usableH/(n*4+(n-1)*6.5+10), 2.0)
	}
	e.sp = sp
	e.def = fontDefaults()
	e.columns()

	// the left edge: part names on the first system, brackets and braces
	names := len(s.Parts) > 1
	nameW, abbrW := 0.0, 0.0
	if names {
		for _, p := range s.Parts {
			nameW = max(nameW, e.text.width(p.Name, 2.2*sp, false, false)/sp)
			abbrW = max(abbrW, e.text.width(p.Abbrev, 2.2*sp, false, false)/sp)
		}
	}
	bracket := e.bracketWidth()
	indent := bracket
	if nameW > 0 {
		indent += nameW + 1.2
	}
	later := bracket
	if abbrW > 0 {
		later += abbrW + 1.2
	}
	width := (pageW-2*margin)/sp - later
	systems := e.breakSystems(width, indent-later)
	for i, sys := range systems {
		sys.left = later
		if i == 0 {
			sys.left = indent
		}
	}

	// staves of each system and their distances
	for _, sys := range systems {
		sys.staves = make([]*staffLayout, len(e.staves))
		for si := range e.staves {
			sys.staves[si] = e.layoutStaff(sys, si)
		}
		sys.y = make([]float64, len(e.staves))
		for si := 1; si < len(e.staves); si++ {
			minGap := 6.5
			if e.staves[si].part == e.staves[si-1].part {
				minGap = 5.5
			}
			gap := max(minGap, sys.staves[si-1].below-4-sys.staves[si].above+1.0)
			sys.y[si] = sys.y[si-1] + 4 + gap
		}
		sys.above = sys.staves[0].above
		last := len(e.staves) - 1
		sys.below = sys.y[last] + sys.staves[last].below
	}

	// pages
	type pageInfo struct {
		systems []*system
		titleH  float64 // pt
	}
	var pages []*pageInfo
	titleH := e.titleHeight(title)
	cur := &pageInfo{titleH: titleH}
	pages = append(pages, cur)
	y := margin + titleH
	bottom := pageH - margin
	for i, sys := range systems {
		// systems 2.5 spaces apart, their staves at least 8 spaces apart
		top := y - sys.above*sp
		if n := len(cur.systems); n > 0 {
			prev := cur.systems[n-1]
			top = max(top+2.5*sp, prev.top+(prev.y[len(prev.y)-1]+4+8)*sp)
		}
		brk := i > 0 && s.Measures[sys.first].Break == BreakPage
		if len(cur.systems) > 0 && (top+sys.below*sp > bottom || brk) {
			cur = &pageInfo{}
			pages = append(pages, cur)
			top = margin - sys.above*sp
		}
		sys.top = top
		sys.page = len(pages) - 1
		y = sys.top + sys.below*sp
		cur.systems = append(cur.systems, sys)
	}
	// spread the systems of full pages over the page
	for pi, pg := range pages {
		if pi == len(pages)-1 || len(pg.systems) < 2 {
			continue
		}
		last := pg.systems[len(pg.systems)-1]
		free := bottom - (last.top + last.below*sp)
		start := margin + pg.titleH
		used := last.top + last.below*sp - start
		if used < 0.6*(bottom-start) {
			continue
		}
		step := min(free/float64(len(pg.systems)-1), 10*sp)
		for k, sys := range pg.systems {
			sys.top += step * float64(k)
		}
	}

	// the document
	view := doc.NewView("score", bdf.ViewFixed, title)
	cvs := canvas.NewBuilder(doc, set)
	glyphs := e.glyphCollection(systems)
	var ext bdf.Hash
	if len(glyphs.paths) > 0 {
		ext = doc.AddPaths(glyphs.paths)
	}
	var cues bdf.Cues
	points := make([][]cuePoint, len(s.Measures))
	type laid struct {
		pg *bdf.Page
		cv *canvas.Canvas
	}
	var out []laid
	for pi, pg := range pages {
		p := view.AddPage(float32(pageW), float32(pageH))
		cv := cvs.New()
		cv.Obj.SetBBox(0, 0, p.W, p.H)
		cv.Obj.FillColor(bdf.RGB(255, 255, 255))
		cv.Obj.FillRect(0, 0, p.W, p.H)
		cv.Obj.FillColor(bdf.RGB(0, 0, 0))
		cv.Obj.StrokeColor(bdf.RGB(0, 0, 0))
		cv.Drawn = true
		w := &writer{e: e, cv: cv, glyphs: glyphs, ext: ext, refs: map[int]bdf.PathRef{}}
		if pi == 0 {
			e.drawTitle(w, title, pageW)
		}
		for _, sys := range pg.systems {
			ox := margin + sys.left*sp
			oy := sys.top
			w.obj().Mark(bdf.MarkFigure, e.systemAlt(sys))
			e.drawSystem(w, sys, ox, oy)
			w.obj().Mark(bdf.MarkEnd, "")
			// cue points
			si := uint32(len(cues.Systems))
			last := len(e.staves) - 1
			cues.Systems = append(cues.Systems, bdf.CueSystem{Page: uint32(pi),
				X: float32(ox), Y: float32(oy - sp), W: float32(sys.width * sp),
				H: float32((sys.y[last] + 6) * sp)})
			for mi := sys.first; mi <= sys.last; mi++ {
				col := e.cols[mi]
				for _, sl := range col.slots {
					points[mi] = append(points[mi], cuePoint{sl.offset, si, ox + (col.x+sl.x+0.6)*sp})
				}
				points[mi] = append(points[mi], cuePoint{s.Measures[mi].Length, si, ox + (col.x+col.width)*sp})
			}
		}
		if pi > 0 {
			e.drawPageNumber(w, pi+1, pageW, pageH)
		}
		out = append(out, laid{p, cv})
	}
	if !o.SystemFonts {
		res.EmbeddedFonts = set.Embed(doc, fontset.EmbedOptions{NoSubset: o.NoSubset, NoWOFF2: o.NoWOFF2, IgnoreFSType: o.IgnoreFSType})
	}
	cvs.Encode()
	set.ReportMissing()
	for _, l := range out {
		l.pg.Layers = []bdf.Layer{{Role: bdf.RoleBody, Obj: l.cv.Hash()}}
	}
	keep := make([]bool, len(view.Pages))
	if o.Pages != nil {
		var sel []*bdf.Page
		for _, n := range o.Pages.Numbers(len(view.Pages)) {
			if n < 1 || n > len(view.Pages) {
				warn(fmt.Sprintf("page %d does not exist (the score has %d)", n, len(view.Pages)))
				continue
			}
			sel = append(sel, view.Pages[n-1])
			keep[n-1] = true
		}
		view.Pages = sel
	} else {
		for i := range keep {
			keep[i] = true
		}
	}

	if s.Play != nil && !o.NoPlay {
		e.play(doc, view, &cues, points, keep)
		res.Seconds = s.Play.seconds()
	}
	if !o.NoTextIndex {
		if _, err := doc.BuildTextIndex(view); err != nil {
			warn(fmt.Sprintf("text index: %v", err))
		}
	}
	res.Doc = doc
	res.Pages = len(view.Pages)
	res.Systems = len(systems)
	res.Measures = len(s.Measures)
	res.Parts = len(s.Parts)
	res.Staves = len(e.staves)
	res.Notes = e.notes
	return res, nil
}

// cuePoint is a place of a measure the playing position passes.
type cuePoint struct {
	offset int
	system uint32
	x      float64
}

// play stores the music played and the cues.
func (e *engraver) play(doc *bdf.Document, view *bdf.View, cues *bdf.Cues, points [][]cuePoint, keep []bool) {
	s := e.s
	smf := s.SMF
	if smf == nil {
		smf = s.Play.SMF()
	}
	h := doc.AddPart(bdf.PartSeq, smf)
	view.Play = &bdf.Play{Seq: h.String()}
	order := s.PlayOrder
	if order == nil {
		for mi := range s.Measures {
			order = append(order, PlayedMeasure{Tick: e.starts[mi], Measure: mi})
		}
	}
	div := s.Play.Division
	if div <= 0 {
		div = PPQ
	}
	// pages left out of a selection lose their systems
	pageIndex := make([]int, len(keep))
	n := 0
	for i, k := range keep {
		pageIndex[i] = -1
		if k {
			pageIndex[i] = n
			n++
		}
	}
	sysIndex := make([]int, len(cues.Systems))
	var systems []bdf.CueSystem
	for i, cs := range cues.Systems {
		sysIndex[i] = -1
		if p := int(cs.Page); p < len(pageIndex) && pageIndex[p] >= 0 {
			cs.Page = uint32(pageIndex[p])
			sysIndex[i] = len(systems)
			systems = append(systems, cs)
		}
	}
	var list []bdf.Cue
cues:
	for _, pm := range order {
		if pm.Measure < 0 || pm.Measure >= len(points) {
			continue
		}
		for _, p := range points[pm.Measure] {
			si := sysIndex[p.system]
			if si < 0 {
				continue
			}
			if len(list) >= maxCues {
				e.warn(fmt.Sprintf("the music has more than %d cues; the rest are left out", maxCues))
				break cues
			}
			t := pm.Tick + p.offset
			tick := uint32((int64(t)*int64(div) + PPQ/2) / PPQ)
			list = append(list, bdf.Cue{Tick: tick, System: uint32(si), X: float32(p.x)})
		}
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].Tick < list[j].Tick })
	if len(systems) > 0 {
		hc := doc.AddPart(bdf.PartIndex, bdf.EncodeCues(&bdf.Cues{Systems: systems, Cues: list}))
		view.Play.Cues = hc.String()
	}
}

// seconds is the length of the performance with its tempo changes.
func (p *Performance) seconds() float64 {
	end := p.end()
	tempos := append([]Tempo(nil), p.Tempo...)
	sort.SliceStable(tempos, func(i, j int) bool { return tempos[i].Tick < tempos[j].Tick })
	bpm, last, sec := 120.0, 0, 0.0
	for _, t := range tempos {
		if t.Tick >= end {
			break
		}
		if t.BPM <= 0 {
			continue
		}
		sec += float64(t.Tick-last) / PPQ * 60 / bpm
		bpm, last = t.BPM, t.Tick
	}
	sec += float64(end-last) / PPQ * 60 / bpm
	return math.Round(sec*10) / 10
}

// systemAlt is the alternative text of a system.
func (e *engraver) systemAlt(sys *system) string {
	a, b := e.measureNumber(sys.first), e.measureNumber(sys.last)
	if a == b {
		return "Measure " + a
	}
	return "Measures " + a + "–" + b
}

// bracketWidth is the room left of the staves for brackets and braces.
func (e *engraver) bracketWidth() float64 {
	w := 0.0
	if len(e.s.Parts) > 1 && !e.grand {
		w = 1.6
	}
	if e.grand {
		w = 1.2
	}
	return w
}

// titleHeight is the height of the title block on the first page.
func (e *engraver) titleHeight(title string) float64 {
	sp := e.sp
	h := 0.0
	if title != "" {
		h += 4.4*sp*1.2 + sp
	}
	if e.s.Subtitle != "" {
		h += 2.8 * sp * 1.3
	}
	if e.s.Composer != "" || e.s.Lyricist != "" {
		h += 2.2 * sp * 1.4
	}
	if e.s.Arranger != "" {
		h += 2.2 * sp * 1.3
	}
	if h > 0 {
		h += 3 * sp
	}
	return h
}

// summary describes a result for the command line.
func summary(r *Result) string {
	s := fmt.Sprintf("%d pages, %d systems, %d measures, %d parts on %d staves, %d notes", r.Pages, r.Systems,
		r.Measures, r.Parts, r.Staves, r.Notes)
	if r.Seconds > 0 {
		s += fmt.Sprintf(", %d:%02d played", int(r.Seconds)/60, int(r.Seconds)%60)
	}
	if r.EmbeddedFonts > 0 {
		s += fmt.Sprintf(", %d embedded font(s)", r.EmbeddedFonts)
	}
	return s
}
