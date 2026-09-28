package epub

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"
)

// limits sets what the tests make small, and puts it back after the test.
func limits(t *testing.T, inflated, text int64, rules int) {
	t.Helper()
	i, x, r, g := maxInflated, maxText, maxRules, maxTags
	t.Cleanup(func() { maxInflated, maxText, maxRules, maxTags = i, x, r, g })
	maxInflated, maxText, maxRules = inflated, text, rules
}

// The files read from a publication are, decompressed, at most maxInflated
// larger than the publication: a ZIP file holds any number of files that
// compress to nearly nothing. What is read after that is left out, with
// one warning.
func TestInflated(t *testing.T) {
	limits(t, 100<<10, maxText, maxRules)
	words := strings.Repeat("word ", 8000) // 40 KB, some 200 bytes in the publication
	var names []string
	var files []file
	for i := range 6 {
		name := fmt.Sprintf("c%d.xhtml", i)
		names = append(names, name)
		files = append(files, file{"OEBPS/" + name, chapterDoc("en", "", `<link rel="stylesheet" href="s.css"/>`,
			fmt.Sprintf(`<p>chapter%d <img src="p%d.png" alt="picture%d"/> %s</p>`, i, i, i, words))})
		files = append(files, file{fmt.Sprintf("OEBPS/p%d.png", i), string(tinyPNG)})
	}
	files = append(files, file{"OEBPS/s.css", "p { text-align: center }"})
	data := makeEPUB(t, opf3("", "", names...), files...)
	if len(data) > 8<<10 {
		t.Fatalf("a publication of %d bytes", len(data))
	}
	res, r := convertBytes(t, data, nil)
	// two chapters of 40 KB are read, and the pictures of the first (the
	// second is read after the chapters)
	text := ""
	for _, p := range viewPages(t, r, "pages") {
		text += p.text + "\n"
	}
	for i := range 6 {
		if got, want := strings.Contains(text, fmt.Sprintf("chapter%d", i)), i < 2; got != want {
			t.Errorf("chapter %d read: %v", i, got)
		}
	}
	if res.Chapters != 2 || res.Images != 2 || strings.Contains(text, "picture") {
		t.Errorf("%d chapters, %d images", res.Chapters, res.Images)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "more than 0 MiB larger than it, decompressed") {
		t.Errorf("warnings %q", res.Warnings)
	}

	// no chapter can be read: an error
	limits(t, 1<<10, maxText, maxRules)
	if _, err := Convert(bytes.NewReader(data), int64(len(data)), testOptions()); err == nil {
		t.Errorf("converted with nothing read")
	}
}

// Content documents and style sheets of more than maxText are left out:
// they take many times their size when they are read.
func TestTextSize(t *testing.T) {
	limits(t, maxInflated, 32<<10, maxRules)
	long := strings.Repeat("word ", 8000)
	data := makeEPUB(t, opf3("", "", "a.xhtml", "b.xhtml"),
		file{"OEBPS/a.xhtml", chapterDoc("en", "", `<link rel="stylesheet" href="s.css"/>`, `<p class="c">short <img src="p.png" alt=""/></p>`)},
		file{"OEBPS/b.xhtml", chapterDoc("en", "", "", "<p>long "+long+"</p>")},
		file{"OEBPS/s.css", ".c { display: none } /* " + long + " */"},
		file{"OEBPS/p.png", string(tinyPNG) + long}) // pictures may be larger
	res, r := convertBytes(t, data, nil)
	pages := viewPages(t, r, "pages")
	if res.Chapters != 1 || res.Images != 1 || len(pages) != 1 || !strings.Contains(pages[0].text, "short") {
		t.Errorf("%d chapters, %d images, %d pages", res.Chapters, res.Images, len(pages))
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "OEBPS/b.xhtml: larger than 0 MiB; left out") {
		t.Errorf("warnings %q", res.Warnings)
	}
}

// The content documents of a publication have at most maxTags tags
// together: their elements take memory, however little text they hold, and
// compress to nearly nothing. The content documents after that are left
// out, with one warning.
func TestTags(t *testing.T) {
	limits(t, maxInflated, maxText, maxRules)
	maxTags = 1000
	var names []string
	var files []file
	for i := range 6 {
		name := fmt.Sprintf("c%d.xhtml", i)
		names = append(names, name)
		// 300 tags, and those of the document around them
		files = append(files, file{"OEBPS/" + name, chapterDoc("en", "", "", fmt.Sprintf("<p>chapter%d%s</p>", i, strings.Repeat("<b/>", 300)))})
	}
	res, r := convertBytes(t, makeEPUB(t, opf3("", "", names...), files...), nil)
	pages := viewPages(t, r, "pages")
	if res.Chapters != 3 || len(pages) != 3 || !strings.Contains(pages[2].text, "chapter2") {
		t.Errorf("%d chapters, %d pages", res.Chapters, len(pages))
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "more than 1000 tags") {
		t.Errorf("warnings %q", res.Warnings)
	}
}

// A chapter whose elements nest deeper than a parser lets them is left
// out: the readers of the tree would need a stack as deep.
func TestDeepChapter(t *testing.T) {
	deep := strings.Repeat("<b>", 520) + "deep" + strings.Repeat("</b>", 520)
	data := makeEPUB(t, opf3("", "", "a.xhtml", "b.xhtml"),
		file{"OEBPS/a.xhtml", chapterDoc("en", "", "", "<p>"+deep+"</p>")},
		file{"OEBPS/b.xhtml", chapterDoc("en", "", "", "<p>"+strings.Repeat("<b>", 400)+"flat"+strings.Repeat("</b>", 400)+"</p>")})
	res, r := convertBytes(t, data, nil)
	pages := viewPages(t, r, "pages")
	if res.Chapters != 1 || len(pages) != 1 || !strings.Contains(pages[0].text, "flat") {
		t.Errorf("%d chapters, %d pages", res.Chapters, len(pages))
	}
	if len(res.Warnings) != 2 || !strings.Contains(res.Warnings[0], "nested deeper than 512") || !strings.Contains(res.Warnings[1], "OEBPS/a.xhtml") ||
		!strings.Contains(res.Warnings[1], "left out") {
		t.Errorf("warnings %q", res.Warnings)
	}
}

// The SVG documents made of the pages of a fixed-layout book are together
// at most twice what may be read from the publication; the pages after
// that are left blank, with one warning. A page can make a document many
// times its size: of groups that are in each other and that it uses, each
// is copied with what it holds.
func TestFixedSVGSize(t *testing.T) {
	limits(t, 16<<10, maxText, maxRules)
	const groups = 20
	var page strings.Builder
	page.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="0" height="0">`)
	for i := range groups {
		fmt.Fprintf(&page, `<g id="g%d">`, i)
	}
	page.WriteString(`<path d="` + strings.Repeat("M0 0L9 9", 120) + `"/>` + strings.Repeat("</g>", groups) + `</svg>` +
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 600 800">`)
	for i := range groups {
		fmt.Fprintf(&page, `<use href="#g%d"/>`, i)
	}
	page.WriteString(`</svg>`)
	var names []string
	var files []file
	for i := range 3 {
		name := fmt.Sprintf("p%d.xhtml", i)
		names = append(names, name)
		files = append(files, file{"OEBPS/" + name, chapterDoc("en", "", `<meta name="viewport" content="width=600, height=800"/>`, page.String())})
	}
	data := makeEPUB(t, opf3(`<meta property="rendition:layout">pre-paginated</meta>`, "", names...), files...)
	// the pages are 2 KB each, their SVG documents 23 KB: the first of them
	// is less than twice the 17 KB that may be read, two of them are more
	res, r := convertBytes(t, data, nil)
	var drawn []int
	for _, p := range viewPages(t, r, "pages") {
		drawn = append(drawn, len(p.images))
	}
	if !res.FixedLayout || res.Pages != 3 || fmt.Sprint(drawn) != "[1 0 0]" {
		t.Errorf("fixed layout %v, %d pages, drawn %v", res.FixedLayout, res.Pages, drawn)
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "decompressed") {
		t.Errorf("warnings %q", res.Warnings)
	}
	if parts := svgParts(res.Doc); len(parts) != 1 || len(parts[0]) < 20<<10 || len(parts[0]) > 30<<10 {
		t.Errorf("%d SVG documents", len(parts))
	}
}

// styleBefore is what the rules gave an element when every rule of the
// style sheets was kept as a rule of its own.
func styleBefore(rules []cssRule, n *html.Node) string {
	classes := strings.Fields(attrVal(n, "class"))
	type winner struct {
		v           string
		spec, order int
	}
	won := map[string]winner{}
	var keys []string
	for i := range rules {
		r := &rules[i]
		if !r.sel.matches(n, classes) {
			continue
		}
		for _, d := range r.decls {
			if pictureOnly[d[0]] && n.Data != "img" {
				continue
			}
			if d[0] == "display" && d[1] != "none" {
				continue
			}
			w, ok := won[d[0]]
			if ok && (w.spec > r.spec || w.spec == r.spec && w.order > r.order) {
				continue
			}
			if !ok {
				keys = append(keys, d[0])
			}
			won[d[0]] = winner{d[1], r.spec, r.order}
		}
	}
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k + ": " + won[k].v + "; ")
	}
	return b.String()
}

// randomSheet writes a style sheet of rules whose selectors and
// declarations repeat.
func randomSheet(rnd *rand.Rand, rules int) string {
	pick := func(list ...string) string { return list[rnd.Intn(len(list))] }
	var b strings.Builder
	for range rules {
		for i := range 1 + rnd.Intn(2) {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(pick("", "", "p", "img", "span", "html", "*") + pick("", "", "#i", "#j"))
			for range rnd.Intn(3) {
				b.WriteString("." + pick("a", "b", "c", "d"))
			}
			b.WriteString(pick("", "", "", ":root", " p", ":hover"))
		}
		b.WriteString(" { ")
		for range 1 + rnd.Intn(3) {
			b.WriteString(pick("text-align", "display", "width", "height", "writing-mode", "text-combine-upright", "color", "-epub-text-emphasis-style") + ": " +
				pick("none", "block", "center", "1em", "2em", "vertical-rl", "all", "left !important", "") + "; ")
		}
		b.WriteString("}\n")
	}
	return b.String()
}

// The rules of a selector, made one rule, give the elements the styles
// the rules gave them one by one.
func TestStylesAsBefore(t *testing.T) {
	rnd := rand.New(rand.NewSource(5))
	styles := 0
	for round := range 1500 {
		var rules []cssRule
		parseCSS(randomSheet(rnd, 1+rnd.Intn(40)), nil, 0, &rules)
		st := newStyler(rules)
		for range 40 {
			n := newNode([]string{"p", "img", "span", "html", "div"}[rnd.Intn(5)], []string{"", "a", "b", "a b", "d c  a", "b b", "e"}[rnd.Intn(7)],
				[]string{"", "", "color: red"}[rnd.Intn(3)])
			if id := []string{"", "", "i", "j", "k"}[rnd.Intn(5)]; id != "" {
				n.Attr = append(n.Attr, html.Attribute{Key: "id", Val: id})
			}
			if rnd.Intn(8) == 0 {
				n.Namespace = "svg"
			}
			own := attrVal(n, "style")
			want := styleBefore(rules, n) + own
			if own == "" {
				want = strings.TrimSuffix(want, " ")
			}
			if n.Namespace != "" {
				want = own
			}
			st.apply(n)
			if got := attrVal(n, "style"); got != want {
				t.Fatalf("round %d: <%s id=%q class=%q>: style %q, want %q", round, n.Data, attrVal(n, "id"), attrVal(n, "class"), got, want)
			}
			if want != own {
				styles++
			}
		}
	}
	if styles < 15000 {
		t.Errorf("%d elements with a style", styles)
	}
}

// A style sheet has at most maxRules rules of selectors that differ; the
// rules after them are left out with a warning.
func TestRulesLimit(t *testing.T) {
	limits(t, maxInflated, maxText, 50)
	var css strings.Builder
	for i := range 80 {
		fmt.Fprintf(&css, ".c%d { text-align: center } .c%d { display: none } .c0 { text-align: right }\n", i, i)
	}
	body := `<p class="c0">zero</p><p class="c49">last</p><p class="c50">after</p><p class="c79">end</p>`
	data := makeEPUB(t, opf3("", "", "a.xhtml", "b.xhtml"),
		file{"OEBPS/a.xhtml", chapterDoc("en", "", `<link rel="stylesheet" href="s.css"/>`, body)},
		file{"OEBPS/b.xhtml", chapterDoc("en", "", `<link rel="stylesheet" href="s.css"/>`, body)},
		file{"OEBPS/s.css", css.String()})
	res, r := convertBytes(t, data, nil)
	for i, p := range viewPages(t, r, "pages") {
		// the page number, and the paragraphs that are not hidden
		if want := fmt.Sprintf("%d\nafter\nend", i+1); p.text != want {
			t.Errorf("page %d: %q", i+1, p.text)
		}
	}
	if len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "more than 50 rules") {
		t.Errorf("warnings %q", res.Warnings)
	}
}

// fastest returns the shortest of three runs of fn.
func fastest(fn func()) time.Duration {
	best := time.Duration(1 << 62)
	for range 3 {
		start := time.Now()
		fn()
		best = min(best, time.Since(start))
	}
	return best
}

// The elements of a chapter take their styles in a time that does not
// grow with the number of times the style sheet repeats a rule, and the
// chapters of a book read the style sheet they share once.
func TestRepeatedRules(t *testing.T) {
	book := func(rules int) []byte {
		chapter := chapterDoc("en", "", `<link rel="stylesheet" href="s.css"/>`, strings.Repeat(`<p class="c">x</p>`, 300))
		return makeEPUB(t, opf3("", "", "a.xhtml", "b.xhtml", "c.xhtml", "d.xhtml"), file{"OEBPS/a.xhtml", chapter}, file{"OEBPS/b.xhtml", chapter},
			file{"OEBPS/c.xhtml", chapter}, file{"OEBPS/d.xhtml", chapter}, file{"OEBPS/s.css", strings.Repeat(".c { text-align: center }\n", rules)})
	}
	load := func(data []byte) func() {
		return func() {
			pub, err := open(bytes.NewReader(data), int64(len(data)))
			if err != nil {
				t.Fatal(err)
			}
			c := &converter{pub: pub, opts: &Options{}, byPath: map[string]*chapter{}, sheets: map[string]string{}, warned: map[string]bool{},
				parsed: map[string][]cssRule{}, stylers: map[string]*styler{}}
			c.dc, c.props = pub.metadata()
			if err := c.spine(); err != nil {
				t.Fatal(err)
			}
			for _, ch := range c.chapters {
				if err := c.load(ch); err != nil {
					t.Fatal(err)
				}
				if p := findElement(ch.body, "p"); attrVal(p, "style") != "text-align: center;" {
					t.Fatalf("style %q", attrVal(p, "style"))
				}
			}
		}
	}
	few, many := fastest(load(book(10))), fastest(load(book(3000)))
	// 3000 rules were tried for each of the 1200 elements, and read for
	// each of the 4 chapters
	if many > 4*few+20*time.Millisecond {
		t.Errorf("10 rules took %v, 3000 rules %v", few, many)
	}
}
