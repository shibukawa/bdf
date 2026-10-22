package markdownbdf

import (
	"bytes"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/encoding/parquet"
)

var (
	// bulletRe matches the marker of an item of a bulleted list, as the
	// converters draw them.
	bulletRe = regexp.MustCompile(`^[•◦▪▫●○■□◆◇‣⁃·\-–—*►➢✓✔☐☑]$`)
	// orderedRe matches the marker of an item of a numbered list.
	orderedRe = regexp.MustCompile(`^(\(?\d{1,3}[.)]|\(\d{1,3}\)|[a-zA-Z][.)]|[ivxlcIVXLC]{1,5}[.)]|[①-⑳])$`)
	numRe     = regexp.MustCompile(`\d+`)
	listRe    = regexp.MustCompile(`^\s*(-|\d+\.) `)
	startRe   = regexp.MustCompile(`^\d{1,9}[.)] `)
	// orderedMarkerRe matches the marker of an item of a list that says it
	// is numbered: a number or letters with a period or a parenthesis, a
	// number alone (also full-width), a circled number.
	orderedMarkerRe = regexp.MustCompile(`^(\(?[\p{N}\p{L}]{1,4}[.)．）、]|\p{N}{1,4}|[①-⑳㉑-㉟])$`)
)

// bulletMarker reports whether a run is the marker of an item of a list
// that says it is bulleted: a few characters, none a letter or a digit.
func bulletMarker(s string) bool {
	rs := []rune(s)
	if len(rs) == 0 || len(rs) > 3 {
		return false
	}
	for _, r := range rs {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// escape keeps the characters of text from being read as Markdown.
func escape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\', '*', '_', '`', '[', ']', '<', '>', '|', '~':
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// escapeStart keeps the first characters of a block from starting another
// kind of block.
func escapeStart(s string) string {
	t := strings.TrimLeft(s, " ")
	if t == "" {
		return s
	}
	switch {
	case strings.HasPrefix(t, "#"), strings.HasPrefix(t, ">"), strings.HasPrefix(t, "+ "), strings.HasPrefix(t, "- "), strings.HasPrefix(t, "---"):
		return "\\" + s
	}
	if m := startRe.FindString(t); m != "" {
		return s[:len(m)-2] + "\\" + s[len(m)-2:]
	}
	return s
}

// style is what the Markdown wraps a stretch of runs in.
type style struct {
	bold, italic, mono bool
	shift              int // 1 a superscript, -1 a subscript
	link               string
}

// shifts tells the runs that are superscripts (1) or subscripts (-1): the
// runs in a smaller font than the run before or after them, raised or
// lowered from its baseline.
func shifts(runs []*run) []int {
	out := make([]int, len(runs))
	for i, r := range runs {
		if r.AltText || r.Size <= 0 {
			continue
		}
		var ref *run
		switch {
		case i > 0 && runs[i-1].Size > r.Size*1.15 && !runs[i-1].AltText:
			ref = runs[i-1]
		case i+1 < len(runs) && runs[i+1].Size > r.Size*1.15 && !runs[i+1].AltText:
			ref = runs[i+1]
		default:
			continue
		}
		switch dy := ref.Y - r.Y; {
		case dy > ref.Size*0.15 && dy < ref.Size*1.2:
			out[i] = 1
		case dy < -ref.Size*0.05 && dy > -ref.Size*0.8:
			out[i] = -1
		}
	}
	return out
}

// adjacent reports whether a run starts right where the one before it ends.
func adjacent(prev, cur *run) bool {
	return prev.Advance > 0 && cur.X-(prev.X+prev.Advance) < prev.Size*0.3 && abs(float64(cur.X-(prev.X+prev.Advance))) < float64(prev.Size)
}

// inline writes the runs of a leaf. In a heading, or a header cell, the
// weight of the whole is not emphasis.
func inline(runs []*run, heading bool) string {
	var b strings.Builder
	var cur style
	var seg strings.Builder
	shift := shifts(runs)
	flush := func() {
		t := seg.String()
		seg.Reset()
		if t == "" {
			return
		}
		// the spaces around a stretch stay outside its markers
		lead := t[:len(t)-len(strings.TrimLeft(t, " "))]
		trail := t[len(strings.TrimRight(t, " ")):]
		core := strings.TrimSpace(t)
		b.WriteString(lead)
		if core != "" {
			open, close := "", ""
			if cur.mono {
				open, close = "`", "`"
				core = strings.ReplaceAll(core, "`", "``")
			} else {
				core = escape(core)
				if cur.bold && !heading {
					open, close = "**"+open, close+"**"
				}
				if cur.italic {
					open, close = "*"+open, close+"*"
				}
			}
			switch cur.shift {
			case 1:
				open, close = "<sup>"+open, close+"</sup>"
			case -1:
				open, close = "<sub>"+open, close+"</sub>"
			}
			if cur.link != "" {
				open, close = "["+open, close+"]("+cur.link+")"
			}
			b.WriteString(open + core + close)
		}
		b.WriteString(trail)
	}
	for i, r := range runs {
		s := style{bold: r.bold, italic: r.italic, mono: r.mono, shift: shift[i], link: r.link}
		sep := ""
		if i > 0 && r.Sep != bdf.SepNone {
			sep = " "
			// a superscript or subscript joins the word it is on (its
			// baseline differs, which looks like a line break)
			if (shift[i] != 0 || shift[i-1] != 0) && adjacent(runs[i-1], r) {
				sep = ""
			}
		}
		if s != cur {
			flush()
			b.WriteString(sep)
			cur = s
		} else {
			seg.WriteString(sep)
		}
		seg.WriteString(r.Text)
	}
	flush()
	return b.String()
}

// codeLine writes the runs of a paragraph of code as they are: a line of
// a code block, whose wrapped lines (LINE marks) join with a space.
func codeLine(runs []*run) string {
	var b strings.Builder
	for i, r := range runs {
		if i > 0 && r.Sep != bdf.SepNone {
			b.WriteByte(' ')
		}
		b.WriteString(r.Text)
	}
	return strings.TrimRight(b.String(), " ")
}

// isCode reports whether a paragraph is a line of code: one its mark says
// is, or, when guess is set, one whose text is all in a monospaced font.
func isCode(n *node, guess bool) bool {
	if n.kind != bdf.NodeParagraph || len(n.runs) == 0 {
		return false
	}
	if n.code {
		return true
	}
	if !guess {
		return false
	}
	for _, r := range n.runs {
		if !r.mono && strings.TrimSpace(r.Text) != "" {
			return false
		}
	}
	return true
}

// guessCode reports whether the monospaced paragraphs of a view, whose
// marks do not say they are code, are code: when most of its paragraphs
// are not monospaced (a document set in a monospaced font is not a
// listing).
func guessCode(roots []*node) bool {
	mono, all := 0, 0
	var count func(n *node)
	count = func(n *node) {
		for _, k := range n.kids {
			if k.kind == bdf.NodeParagraph && len(k.runs) > 0 {
				all++
				if isCode(k, true) {
					mono++
				}
			}
			count(k)
		}
	}
	for _, r := range roots {
		count(r)
	}
	return mono > 0 && (all < 4 || mono*2 < all)
}

// imagePath writes an image once and returns the path the Markdown refers
// to it by ("" when it is not written).
func (w *writer) imagePath(h bdf.Hash) string {
	if p, ok := w.images[h]; ok {
		return p
	}
	path := ""
	if w.opts.Image != nil {
		if part := w.doc.Part(h); part == nil {
			w.warnf("image %s is missing", h)
		} else {
			name := h.String()[:16] + imageExt(part.Data)
			p, err := w.opts.Image(name, part.Data)
			if err != nil {
				w.warnf("image %s: %v", name, err)
			} else {
				path = p
				w.res.Images++
			}
		}
	}
	w.images[h] = path
	return path
}

// imageExt tells the extension of an image file by its first bytes.
func imageExt(d []byte) string {
	switch {
	case bytes.HasPrefix(d, []byte("\x89PNG")):
		return ".png"
	case bytes.HasPrefix(d, []byte("\xff\xd8")):
		return ".jpg"
	case bytes.HasPrefix(d, []byte("GIF8")):
		return ".gif"
	case bytes.HasPrefix(d, []byte("RIFF")) && len(d) > 12 && string(d[8:12]) == "WEBP":
		return ".webp"
	case len(d) > 12 && string(d[4:8]) == "ftyp":
		return ".avif"
	case bytes.HasPrefix(d, []byte("BM")):
		return ".bmp"
	case bytes.Contains(d[:min(len(d), 1024)], []byte("<svg")):
		return ".svg"
	}
	return ".bin"
}

// blocks writes a sequence of nodes as blocks: paragraphs, headings, lists,
// tables, figures, fences of the lines of code among them, and quotations
// of the nodes deeper in quotations than depth, the depth of the
// quotation the blocks are in. inCell says the blocks are in a cell or an
// item, where a heading is a bold line.
func (w *writer) blocks(kids []*node, inCell bool, depth int) []string {
	var out []string
	guess := w.guess
	for i := 0; i < len(kids); i++ {
		k := kids[i]
		if k.quote > depth {
			// a quotation: the nodes in it, written one level down
			j := i
			for ; j < len(kids) && kids[j].quote > depth; j++ {
			}
			inner := strings.Join(w.blocks(kids[i:j], inCell, depth+1), "\n\n")
			i = j - 1
			if inner == "" {
				continue
			}
			out = append(out, "> "+strings.ReplaceAll(strings.ReplaceAll(inner, "\n", "\n> "), "> \n", ">\n"))
			continue
		}
		if isCode(k, guess) {
			var lines []string
			j := i
			for ; j < len(kids) && kids[j].quote <= depth && isCode(kids[j], guess); j++ {
				lines = append(lines, codeLine(kids[j].runs))
			}
			i = j - 1
			fence := "```"
			for _, l := range lines {
				for strings.Contains(l, fence) {
					fence += "`"
				}
			}
			out = append(out, fence+"\n"+strings.Join(lines, "\n")+"\n"+fence)
			continue
		}
		if s := w.block(k, inCell, depth); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func (w *writer) block(n *node, inCell bool, depth int) string {
	switch n.kind {
	case bdf.NodeParagraph, bdf.NodeCaption, bdf.NodeCell:
		t := strings.TrimSpace(inline(n.runs, false))
		if t == "" {
			return ""
		}
		return escapeStart(t)
	case bdf.NodeHeading:
		t := strings.TrimSpace(inline(n.runs, true))
		if t == "" {
			return ""
		}
		if inCell {
			return "**" + t + "**"
		}
		level := n.level
		if level == 0 {
			level = 1
		}
		return strings.Repeat("#", level) + " " + t
	case kindImage:
		if p := w.imagePath(n.image); p != "" {
			return "![](" + p + ")"
		}
		return ""
	case bdf.NodeFigure:
		return w.figure(n, inCell, depth)
	case bdf.NodeList:
		return w.list(n, depth)
	case bdf.NodeTable:
		return w.table(n)
	case bdf.NodeItem:
		return strings.Join(w.blocks(n.kids, true, depth), "\n\n")
	}
	return ""
}

// figure writes a figure: its images with its alternative text, or that
// text alone. The text drawn in a figure with alternative text (a chart's
// labels) is left to it, as a reader of the document aloud leaves it
// (docs/spec.md §7.8); the text of a figure without stands.
func (w *writer) figure(n *node, inCell bool, depth int) string {
	var images, rest []*node
	for _, k := range n.kids {
		if k.kind == kindImage {
			images = append(images, k)
		} else {
			rest = append(rest, k)
		}
	}
	alt := strings.TrimSpace(n.alt)
	var parts []string
	for _, im := range images {
		p := w.imagePath(im.image)
		if p == "" {
			continue
		}
		a := ""
		if len(parts) == 0 {
			a = escape(alt)
		}
		parts = append(parts, "!["+a+"]("+p+")")
	}
	if len(parts) == 0 && alt != "" {
		parts = append(parts, "*"+escape(alt)+"*")
	}
	if alt == "" {
		parts = append(parts, w.blocks(rest, inCell, depth)...)
	}
	if n.data != "" {
		// the table the figure is the data of
		if t := w.doc.Table(n.data); t != nil {
			parts = append(parts, w.tableBlocks(t, inCell)...)
		}
	}
	return strings.Join(parts, "\n\n")
}

// tableBlocks writes a table the document holds (docs/spec.md §4.6) as
// Options.Tables says: inline, up to InlineRows rows and the whole as a
// CSV file when there are more; or as a CSV file alone.
func (w *writer) tableBlocks(t *bdf.Table, inCell bool) []string {
	mode := w.opts.Tables
	if mode == "" {
		mode = "inline"
	}
	if mode == "none" {
		return nil
	}
	data := w.tableData(t)
	if data == nil {
		return nil
	}
	var out []string
	if mode == "inline" {
		limit := w.opts.InlineRows
		if limit <= 0 {
			limit = DefaultInlineRows
		}
		tb, err := parquet.Read(data, parquet.Limits{MaxRows: limit + 1})
		if err != nil {
			w.warnf("table %s: %v", t.ID, err)
			return nil
		}
		more := len(tb.Rows) > limit
		if more {
			tb.Rows = tb.Rows[:limit]
		}
		if g := gfmTable(tb, inCell); g != "" {
			out = append(out, g)
			w.res.Tables++
		}
		if !more {
			return out
		}
		if link := w.tableFile(t); link != "" {
			return append(out, link)
		}
		return append(out, fmt.Sprintf("(%d of %d rows)", limit, t.Rows))
	}
	if link := w.tableFile(t); link != "" {
		out = append(out, link)
	}
	return out
}

// tableData is the Parquet file of a table, nil when it is missing.
func (w *writer) tableData(t *bdf.Table) []byte {
	h, err := bdf.ParseHash(t.Part)
	if err != nil {
		w.warnf("table %s: %v", t.ID, err)
		return nil
	}
	part := w.doc.Part(h)
	if part == nil {
		w.warnf("table %s: its part is missing", t.ID)
		return nil
	}
	return part.Data
}

// tableFile writes a table as a CSV file (once) and returns the line that
// links to it, "" when files are not written.
func (w *writer) tableFile(t *bdf.Table) string {
	if w.opts.File == nil {
		return ""
	}
	path, done := w.csvs[t.ID]
	if !done {
		path = ""
		if data := w.tableData(t); data != nil {
			tb, err := parquet.Read(data, parquet.Limits{})
			if err != nil {
				w.warnf("table %s: %v", t.ID, err)
			} else {
				var b bytes.Buffer
				if err := parquet.WriteCSV(&b, tb); err == nil {
					if p, err := w.opts.File(t.ID+".csv", b.Bytes()); err != nil {
						w.warnf("table %s: %v", t.ID, err)
					} else {
						path = p
						w.res.Tables++
					}
				}
			}
		}
		w.csvs[t.ID] = path
	}
	if path == "" {
		return ""
	}
	title := t.Title
	if title == "" {
		title = t.ID
	}
	return fmt.Sprintf("[%s: %d rows × %d columns](%s)", escape(title), t.Rows, len(t.Columns), path)
}

// gfmTable writes a table in memory as a GitHub Flavored Markdown table.
func gfmTable(t *parquet.Table, inCell bool) string {
	if len(t.Columns) == 0 {
		return ""
	}
	cell := func(s string) string {
		s = escape(s)
		s = strings.ReplaceAll(s, "\r\n", "<br>")
		s = strings.ReplaceAll(s, "\n", "<br>")
		return s
	}
	var lines []string
	names := make([]string, len(t.Columns))
	for i, c := range t.Columns {
		names[i] = cell(c.Name)
	}
	lines = append(lines, "| "+strings.Join(names, " | ")+" |", "|"+strings.Repeat("---|", len(t.Columns)))
	for _, r := range t.Rows {
		cols := make([]string, len(t.Columns))
		for j, c := range t.Columns {
			var v any
			if j < len(r) {
				v = r[j]
			}
			cols[j] = cell(parquet.Text(c, v))
		}
		lines = append(lines, "| "+strings.Join(cols, " | ")+" |")
	}
	s := strings.Join(lines, "\n")
	if inCell {
		s = strings.ReplaceAll(s, "\n", "<br>")
	}
	return s
}

// marker tells whether the first run of an item is the marker the
// converter drew for it, and whether it numbers the item. kind is what
// the list's mark says its markers are ("bullet", "ordered" or "").
func marker(it *node, kind string) (first *node, ordered bool, text string) {
	if len(it.kids) == 0 {
		return nil, false, ""
	}
	first = it.kids[0]
	if !first.leaf() || len(first.runs) == 0 {
		return first, false, ""
	}
	t := strings.TrimSpace(first.runs[0].Text)
	switch kind {
	case "bullet":
		if bulletMarker(t) || bulletRe.MatchString(t) {
			return first, false, t
		}
		return first, false, ""
	case "ordered":
		if orderedMarkerRe.MatchString(t) || orderedRe.MatchString(t) {
			return first, true, t
		}
		return first, true, ""
	}
	switch {
	case bulletRe.MatchString(t):
		return first, false, t
	case orderedRe.MatchString(t):
		return first, true, t
	}
	return first, false, ""
}

// list writes a list: an item a line, what follows its first paragraph
// indented under it, a nested list right after the item it is in.
func (w *writer) list(n *node, depth int) string {
	var lines []string
	num := 0
	ordered := n.list == "ordered"
	for i, it := range n.kids {
		if it.kind != bdf.NodeItem {
			if s := w.block(it, true, depth); s != "" {
				lines = append(lines, s)
			}
			continue
		}
		first, ord, mk := marker(it, n.list)
		if i == 0 && n.list == "" {
			ordered = ord
		}
		kids := it.kids
		head := ""
		if first != nil && first.leaf() {
			runs := first.runs
			if mk != "" {
				if ordered && i == 0 {
					if m := numRe.FindString(mk); m != "" {
						num, _ = strconv.Atoi(m)
						num--
					}
				}
				runs = runs[1:]
			}
			head = strings.TrimSpace(inline(runs, first.kind == bdf.NodeHeading))
			kids = kids[1:]
		}
		bullet := "- "
		if ordered {
			num++
			bullet = strconv.Itoa(num) + ". "
		}
		pad := strings.Repeat(" ", len(bullet))
		item := bullet + head
		for _, b := range w.blocks(kids, true, depth) {
			b = pad + strings.ReplaceAll(b, "\n", "\n"+pad)
			b = strings.ReplaceAll(b, "\n"+pad+"\n", "\n\n")
			if listRe.MatchString(b) {
				item += "\n" + b
			} else {
				item += "\n\n" + b
			}
		}
		lines = append(lines, item)
	}
	return strings.Join(lines, "\n")
}

// cellText writes what is in a cell on one line of a table.
func (w *writer) cellText(n *node) string {
	header := n.cell != nil && n.cell.Scope != ""
	var parts []string
	if n.leaf() {
		if t := strings.TrimSpace(inline(n.runs, header)); t != "" {
			parts = append(parts, t)
		}
	} else if header {
		for _, k := range n.kids {
			if t := strings.TrimSpace(inline(k.runs, true)); t != "" {
				parts = append(parts, t)
			}
		}
	} else {
		for _, b := range w.blocks(n.kids, true, 0) {
			parts = append(parts, strings.ReplaceAll(b, "\n", "<br>"))
		}
	}
	return strings.Join(parts, "<br>")
}

// grid writes cells as a table: the cells with references in their
// places, a cell that spans columns or rows in its top-left one. A table
// without a header row gets an empty one, as the format wants one.
func (w *writer) grid(cells []*node) string {
	type slot struct {
		n *node
		c bdf.CellRef
	}
	var slots []slot
	rowSet := map[int]bool{}
	minC, maxC := math.MaxInt, -1
	for _, c := range cells {
		if c.cell == nil {
			continue
		}
		slots = append(slots, slot{c, *c.cell})
		rowSet[c.cell.Row] = true
		minC, maxC = min(minC, c.cell.Col), max(maxC, c.cell.Col+c.cell.Cols-1)
	}
	if len(slots) == 0 {
		// cells without places: one a row
		var lines []string
		for _, c := range cells {
			lines = append(lines, "| "+w.cellText(c)+" |")
		}
		if len(lines) == 0 {
			return ""
		}
		return strings.Join(append([]string{"| |", "|---|"}, lines...), "\n")
	}
	rows := make([]int, 0, len(rowSet))
	for r := range rowSet {
		rows = append(rows, r)
	}
	sort.Ints(rows)
	at := map[[2]int]*slot{}
	for i := range slots {
		s := &slots[i]
		if _, ok := at[[2]int{s.c.Row, s.c.Col}]; !ok {
			at[[2]int{s.c.Row, s.c.Col}] = s
		}
	}
	header := false
	for _, s := range slots {
		if s.c.Row == rows[0] && s.c.Scope == "col" {
			header = true
		}
	}
	width := maxC - minC + 1
	var lines []string
	if !header {
		lines = append(lines, "|"+strings.Repeat(" |", width), "|"+strings.Repeat("---|", width))
	}
	for i, r := range rows {
		cols := make([]string, width)
		for c := minC; c <= maxC; c++ {
			if s := at[[2]int{r, c}]; s != nil {
				cols[c-minC] = w.cellText(s.n)
			}
		}
		lines = append(lines, "| "+strings.Join(cols, " | ")+" |")
		if i == 0 && header {
			lines = append(lines, "|"+strings.Repeat("---|", width))
		}
	}
	return strings.Join(lines, "\n")
}

// table writes a table: its caption, then its cells.
func (w *writer) table(n *node) string {
	var caption, cells []*node
	for _, k := range n.kids {
		switch k.kind {
		case bdf.NodeCell:
			cells = append(cells, k)
		case bdf.NodeCaption:
			caption = append(caption, k)
		}
	}
	out := w.grid(cells)
	if len(caption) > 0 {
		var ts []string
		for _, c := range caption {
			if t := strings.TrimSpace(inline(c.runs, false)); t != "" {
				ts = append(ts, t)
			}
		}
		if len(ts) > 0 {
			out = escapeStart(strings.Join(ts, " ")) + "\n\n" + out
		}
	}
	return out
}

// sheet writes the cells of a sheet (CELLs outside a table) as one table,
// and what else is drawn on it as blocks.
func (w *writer) sheet(root *node) []string {
	var cells, rest []*node
	for _, k := range root.kids {
		if k.kind == bdf.NodeCell && k.cell != nil {
			cells = append(cells, k)
		} else {
			rest = append(rest, k)
		}
	}
	var out []string
	if len(cells) > 0 {
		out = append(out, w.grid(cells))
	}
	return append(out, w.blocks(rest, false, 0)...)
}
