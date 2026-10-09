package audio

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shibukawa/bdf"
)

// The card page: its width comes from the cover (a cover narrower than
// coverMin or wider than coverMax points is scaled), or is the side of the
// placeholder; the text keeps a margin.
const (
	coverMin        = 360
	coverMax        = 720
	placeholderSide = 480
	margin          = 32
	labelWidth      = 112
	maxLyricLines   = 2000
)

// Colors of the card.
var (
	colorText        = bdf.RGB(0x1f, 0x23, 0x28)
	colorMuted       = bdf.RGB(0x5c, 0x66, 0x70)
	colorLabel       = bdf.RGB(0x80, 0x89, 0x92)
	colorPlaceholder = bdf.RGB(0xe6, 0xea, 0xee)
	colorIcon        = bdf.RGB(0x9a, 0xa4, 0xad)
	colorRule        = bdf.RGB(0xe0, 0xe4, 0xe8)
)

// cover is the picture the card shows, as it is drawn.
type coverImage struct {
	data   []byte
	format string
	w, h   int // pixels
}

// card lays the page out: the cover (or a placeholder) across the top,
// then the title, the artist and the album, a table of the other fields,
// the lyrics and the chapters. The text is set in the viewer's sans-serif
// font, so lines are wrapped on estimated widths and left unmeasured.
type card struct {
	doc   *bdf.Document
	o     *bdf.Object
	w     float32 // the page width
	y     float32 // the baseline of the next line
	reg   bdf.FontRef
	bold  bdf.FontRef
	table int // rows of the details table
}

// buildCard draws the card of a track and returns the page size and the
// object.
func buildCard(doc *bdf.Document, t *track, cover *coverImage, title, alt string) (w, h float32, obj bdf.Hash) {
	c := &card{doc: doc, o: bdf.NewObject()}
	c.reg = c.o.AddFont(bdf.SystemFont("sans-serif", 400, bdf.StyleNormal))
	c.bold = c.o.AddFont(bdf.SystemFont("sans-serif", 700, bdf.StyleNormal))
	c.drawCover(cover, alt)
	c.y += 26
	tg := t.tags
	c.heading(title, 22, 28, 3)
	if s := strings.Join(tg[keyArtist], ", "); s != "" {
		c.paragraph(s, c.reg, 14, colorText, 20)
	}
	if s := tg.first(keyAlbum); s != "" {
		c.paragraph(s, c.reg, 13, colorMuted, 18)
	}
	c.y += 6
	c.details(t)
	if s := tg.first(keyLyrics); s != "" {
		c.lyrics(s)
	}
	if len(t.chapters) > 0 {
		c.chapters(t.chapters)
	}
	c.y += margin - 8
	h = c.y
	if h < c.w {
		h = c.w
	}
	obj, _ = doc.AddObject(c.o)
	return c.w, h, obj
}

// drawCover draws the cover as wide as the page, or a placeholder square
// with a note on it.
func (c *card) drawCover(cover *coverImage, alt string) {
	if cover == nil {
		c.w = placeholderSide
		c.o.Mark(bdf.MarkFigure, "")
		c.o.FillColor(colorPlaceholder).FillRect(0, 0, c.w, c.w)
		c.o.FillColor(colorIcon).FillPath(c.o.AddPath(notePath(c.w/2, c.w/2, c.w*0.22)), bdf.NonZero)
		c.o.Mark(bdf.MarkEnd, "")
		c.y = c.w
		return
	}
	w := float32(cover.w) * 0.75
	c.w = min(max(w, coverMin), coverMax)
	h := c.w * float32(cover.h) / float32(cover.w)
	img := c.o.AddImage(c.doc.AddImage(cover.data))
	c.o.Mark(bdf.MarkFigure, alt)
	c.o.Smoothing(true, 2)
	c.o.Image(img, 0, 0, c.w, h)
	c.o.Mark(bdf.MarkEnd, "")
	c.y = h
}

// notePath is an eighth note of about side r, centered at (x, y).
func notePath(x, y, r float32) *bdf.Path {
	p := &bdf.Path{}
	// the head, tilted
	p.Ellipse(x-r*0.35, y+r*0.55, r*0.42, r*0.3, -0.45, 0, 6.2832, false)
	p.Close()
	// the stem
	p.Rect(x-r*0.35+r*0.3, y-r*0.9, r*0.1, r*1.45)
	// the flag
	sx, sy := x-r*0.35+r*0.4, y-r*0.9
	p.MoveTo(sx, sy)
	p.CubicTo(sx+r*0.55, sy+r*0.25, sx+r*0.65, sy+r*0.75, sx+r*0.3, sy+r*1.15)
	p.CubicTo(sx+r*0.55, sy+r*0.7, sx+r*0.35, sy+r*0.4, sx, sy+r*0.3)
	p.Close()
	return p
}

// heading sets a heading of at most lines lines.
func (c *card) heading(s string, size, leading float32, lines int) {
	c.o.Mark(bdf.MarkHeading, "1")
	c.lines(wrap(s, size, c.w-2*margin, lines), c.bold, size, colorText, margin, leading)
}

// paragraph sets a paragraph of one field.
func (c *card) paragraph(s string, font bdf.FontRef, size float32, color bdf.Color, leading float32) {
	c.o.Mark(bdf.MarkParagraph, "")
	c.lines(wrap(s, size, c.w-2*margin, 8), font, size, color, margin, leading)
}

// lines draws lines one under the other. A line wrapped from the one
// before is marked WRAP when the break is between wide characters (no
// space was dropped there), else LINE (spec §7.8).
func (c *card) lines(lines []string, font bdf.FontRef, size float32, color bdf.Color, x, leading float32) {
	c.o.Font(font, size).FillColor(color)
	for i, l := range lines {
		if i > 0 {
			mark := bdf.MarkLine
			if wide(lastRune(lines[i-1])) && wide(firstRune(l)) {
				mark = bdf.MarkWrap
			}
			c.o.Mark(mark, "")
		}
		c.y += leading
		c.o.FillText(l, x, c.y-leading*0.28, 0)
	}
}

func firstRune(s string) rune {
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

func lastRune(s string) rune {
	r, _ := utf8.DecodeLastRuneInString(s)
	return r
}

// details sets the table of the other fields.
func (c *card) details(t *track) {
	tg := t.tags
	type row struct{ label, value string }
	var rows []row
	add := func(label string, values ...string) {
		if s := strings.Join(values, ", "); s != "" {
			rows = append(rows, row{label, s})
		}
	}
	if artists := tg[keyArtist]; len(tg[keyAlbumArtist]) > 0 && strings.Join(tg[keyAlbumArtist], ", ") != strings.Join(artists, ", ") {
		add("Album artist", tg[keyAlbumArtist]...)
	}
	add("Subtitle", tg[keySubtitle]...)
	add("Composer", tg[keyComposer]...)
	add("Lyricist", tg[keyLyricist]...)
	add("Conductor", tg[keyConductor]...)
	add("Arranger", tg[keyArranger]...)
	add("Remixer", tg[keyRemixer]...)
	add("Grouping", tg[keyGrouping]...)
	add("Genre", tg[keyGenre]...)
	if n := tg.first(keyTrack); n != "" {
		if total := tg.first(keyTrackTotal); total != "" {
			n += " of " + total
		}
		add("Track", n)
	}
	if n := tg.first(keyDisc); n != "" {
		if total := tg.first(keyDiscTotal); total != "" {
			n += " of " + total
		}
		add("Disc", n)
	}
	add("Date", tg.first(keyDate))
	add("Original date", tg.first(keyOriginalDate))
	add("Publisher", tg[keyPublisher]...)
	add("ISRC", tg[keyISRC]...)
	add("Copyright", tg[keyCopyright]...)
	if langs := tg[keyLanguage]; len(langs) > 0 {
		var names []string
		for _, l := range langs {
			if b := bcp47(l); b != "" {
				names = append(names, b)
			}
		}
		add("Language", names...)
	}
	add("BPM", tg.first(keyBPM))
	add("Comment", tg.first(keyComment))
	if t.duration > 0 {
		add("Duration", clock(t.duration))
	}
	add("Format", formatLine(t))
	if len(rows) == 0 {
		return
	}
	c.y += 6
	c.o.FillColor(colorRule).FillRect(margin, c.y, c.w-2*margin, 0.75)
	c.y += 8
	c.o.Mark(bdf.MarkTable, "")
	for i, r := range rows {
		c.o.Mark(bdf.MarkCell, fmt.Sprintf("A%d row", i+1))
		y := c.y
		c.lines([]string{r.label}, c.reg, 10, colorLabel, margin, 17)
		c.y = y
		c.o.Mark(bdf.MarkCell, fmt.Sprintf("B%d", i+1))
		c.lines(wrap(r.value, 11, c.w-2*margin-labelWidth, 12), c.reg, 11, colorText, margin+labelWidth, 17)
	}
	c.o.Mark(bdf.MarkEnd, "")
}

// lyrics sets the lyrics, stanza by stanza.
func (c *card) lyrics(s string) {
	c.section("Lyrics")
	lines := strings.Split(s, "\n")
	if len(lines) > maxLyricLines {
		lines = lines[:maxLyricLines]
	}
	newStanza := true
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			if !newStanza {
				c.y += 8
			}
			newStanza = true
			continue
		}
		if newStanza {
			c.o.Mark(bdf.MarkParagraph, "")
		} else {
			c.o.Mark(bdf.MarkLine, "")
		}
		newStanza = false
		c.lines(wrap(l, 11, c.w-2*margin, 8), c.reg, 11, colorText, margin, 16)
	}
}

// chapters sets the list of chapters.
func (c *card) chapters(chapters []chapter) {
	c.section("Chapters")
	c.o.Mark(bdf.MarkList, "")
	for _, ch := range chapters {
		c.o.Mark(bdf.MarkListItem, "")
		y := c.y
		c.lines([]string{clock(ch.start)}, c.reg, 11, colorLabel, margin, 17)
		c.y = y
		c.lines(wrap(ch.title, 11, c.w-2*margin-labelWidth, 4), c.reg, 11, colorText, margin+labelWidth, 17)
	}
	c.o.Mark(bdf.MarkEnd, "")
}

// section sets the heading of a section.
func (c *card) section(name string) {
	c.y += 14
	c.o.Mark(bdf.MarkHeading, "2")
	c.lines([]string{name}, c.bold, 12, colorMuted, margin, 18)
	c.y += 2
}

// clock writes seconds as m:ss or h:mm:ss.
func clock(seconds float64) string {
	s := int(seconds + 0.5)
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s/60%60, s%60)
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

// formatLine describes the coding ("MP3, 320 kbps, 44.1 kHz, stereo").
func formatLine(t *track) string {
	var parts []string
	codec := t.codec
	if codec == "" {
		codec = formatNames[t.format]
	}
	if t.protected {
		codec += " (protected)"
	}
	parts = append(parts, codec)
	if t.bits > 0 {
		parts = append(parts, fmt.Sprintf("%d-bit", t.bits))
	}
	if t.bitrate > 0 {
		s := fmt.Sprintf("%d kbps", (t.bitrate+500)/1000)
		if t.vbr {
			s += " (variable)"
		}
		parts = append(parts, s)
	}
	if t.sampleRate > 0 {
		parts = append(parts, kHz(t.sampleRate))
	}
	switch t.channels {
	case 0:
	case 1:
		parts = append(parts, "mono")
	case 2:
		parts = append(parts, "stereo")
	default:
		parts = append(parts, fmt.Sprintf("%d channels", t.channels))
	}
	return strings.Join(parts, ", ")
}

// kHz writes a sample rate ("44.1 kHz", "48 kHz", "22.05 kHz").
func kHz(rate int) string {
	s := fmt.Sprintf("%.3f", float64(rate)/1000)
	s = strings.TrimRight(strings.TrimRight(s, "0"), ".")
	return s + " kHz"
}

// estWidth estimates the width of a string in a sans-serif font of a size,
// without a font to measure with.
func estWidth(s string, size float32) float32 {
	var w float32
	for _, r := range s {
		w += runeWidth(r)
	}
	return w * size
}

// runeWidth is the width of a character in ems, roughly.
func runeWidth(r rune) float32 {
	switch {
	case r == ' ':
		return 0.28
	case r < 0x80 && unicode.IsDigit(r):
		return 0.56
	case r < 0x80 && unicode.IsUpper(r):
		return 0.68
	case r < 0x80 && unicode.IsLower(r):
		return 0.53
	case r < 0x80:
		return 0.34
	case wide(r):
		return 1
	}
	return 0.62
}

// wide reports a character set on a full-width cell: CJK and the like.
func wide(r rune) bool {
	return unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) ||
		unicode.Is(unicode.Hangul, r) || r >= 0x3000 && r <= 0x303F || r >= 0xFF00 && r <= 0xFF60 || r >= 0x1F300 && r <= 0x1FAFF
}

// wrap breaks a string into lines of at most width points, at spaces and
// between wide characters, and inside a word that is too long for a line;
// at most max lines, the last one ending with an ellipsis when the text
// goes on.
func wrap(s string, size, width float32, max int) []string {
	s = strings.TrimSpace(s)
	if s == "" || max <= 0 {
		return nil
	}
	var lines []string
	var line []rune
	var lineW float32
	// flush ends the line; false when no more lines may be added
	flush := func() bool {
		if t := strings.TrimSpace(string(line)); t != "" {
			lines = append(lines, t)
		}
		line, lineW = line[:0], 0
		return len(lines) < max
	}
	runes := []rune(s)
	var piece []rune
	cut := false
	i := 0
	for i < len(runes) && !cut {
		// the next unbreakable piece: a word, or one wide character
		j := i
		var w float32
		for j < len(runes) {
			r := runes[j]
			if j > i && (r == ' ' || wide(r) || wide(runes[j-1])) {
				break
			}
			w += runeWidth(r) * size
			j++
			if r == ' ' {
				break
			}
		}
		piece = runes[i:j]
		i = j
		if piece[0] == ' ' && len(line) == 0 {
			continue
		}
		if lineW+w > width && len(line) > 0 {
			if !flush() {
				cut = true
				break
			}
			if piece[0] == ' ' {
				continue
			}
		}
		if w > width { // a word longer than a line: cut it
			for k, r := range piece {
				rw := runeWidth(r) * size
				if lineW+rw > width && len(line) > 0 && !flush() {
					cut = true
					piece = piece[k:]
					break
				}
				line = append(line, r)
				lineW += rw
			}
			continue
		}
		line = append(line, piece...)
		lineW += w
	}
	if !cut {
		flush()
		return lines
	}
	if strings.TrimSpace(string(piece)+string(runes[i:])) == "" {
		return lines
	}
	last := []rune(lines[max-1])
	if len(last) > 1 {
		last = last[:len(last)-1]
	}
	lines[max-1] = strings.TrimSpace(string(last)) + "…"
	return lines
}
