// Package markdownbdf writes a document as Markdown: the structure its MARKs
// carry (headings, paragraphs, lists, tables, figures with their
// alternative texts; docs/spec.md §7.8) as the blocks of GitHub Flavored
// Markdown, the fonts of its runs as emphasis and code, its LINKs as links
// and its images as files beside the text.
//
// The text is the document's in plain text: a document that was encrypted
// should not have it written where the password does not protect it.
package markdownbdf

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/shibukawa/bdf"
)

// Options say what to write and how.
type Options struct {
	// Views are the ids of the views to write. By default the scroll views
	// are written when the document has one (the same text laid out
	// without pages: Word, HTML, Markdown and EPUB documents), else every
	// view.
	Views []string
	// SkipRoles are the roles of the layers that are left out. Nil leaves
	// out DefaultSkipRoles: what is around the content on every page.
	SkipRoles []string
	// InternalLinks keeps the links into the document itself (#page=N),
	// which have no target in the Markdown. By default their text stands
	// alone.
	InternalLinks bool
	// KeepRunning keeps the running heads and feet of the pages of a
	// paged view: the short text at the same place in the top or bottom
	// margin of two pages or more (a journal's name, a paper's title, a
	// notice; the digits of "Page 3 of 9" aside), and page numbers alone
	// there. By default they are left out, as the header and footer
	// layers of a Word document are.
	KeepRunning bool
	// Image is called once for each image the Markdown refers to, with a
	// file name made of the hash of its part and its format (the same
	// image is one file), and the bytes of the part. It returns the path
	// the Markdown refers to the image by: the name, when the images are
	// written beside the text. Nil leaves the images out, and a figure
	// keeps its alternative text.
	Image func(name string, data []byte) (string, error)
	// Tables says how the tables the document holds beside its drawing
	// (docs/spec.md §4.6; the data of charts) are written: "inline" (the
	// default: a table of the figure's data after its text, up to
	// InlineRows rows, and the whole as a CSV file beside the text when
	// there are more and File is set), "files" (as CSV files only, linked
	// from the figure; the tables of sheets too) or "none".
	Tables string
	// InlineRows is the most rows of a table written inline (100 when 0).
	InlineRows int
	// File is called for each CSV file of a table, as Image is for images.
	// Nil writes no files: a table of more rows than InlineRows is cut.
	File func(name string, data []byte) (string, error)
}

// DefaultInlineRows is the rows of a table written inline when
// Options.InlineRows is 0.
const DefaultInlineRows = 100

// DefaultSkipRoles are the roles of the layers that are not the content of
// a page: what is under it and around it on every page.
var DefaultSkipRoles = []string{bdf.RoleHeader, bdf.RoleFooter, bdf.RoleMaster, bdf.RoleBackground}

// Result tells what was written.
type Result struct {
	Views    int
	Images   int
	Tables   int // the tables written, inline or as files
	Warnings []string
}

// Write writes the document as Markdown.
func Write(w io.Writer, doc *bdf.Document, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	views, err := selectViews(doc, opts.Views)
	if err != nil {
		return nil, err
	}
	switch opts.Tables {
	case "", "inline", "files", "none":
	default:
		return nil, fmt.Errorf("markdownbdf: Tables must be inline, files or none, not %q", opts.Tables)
	}
	wr := &writer{doc: doc, opts: opts, objs: map[bdf.Hash]*bdf.ObjectPart{}, faces: map[bdf.Hash]*face{},
		images: map[bdf.Hash]string{}, csvs: map[string]string{}, res: &Result{Views: len(views)}}
	skip := map[string]bool{}
	roles := opts.SkipRoles
	if roles == nil {
		roles = DefaultSkipRoles
	}
	for _, r := range roles {
		skip[r] = true
	}
	var out []string
	if fm := frontMatter(&doc.Meta.DC); fm != "" {
		out = append(out, fm)
	}
	for _, v := range views {
		if len(views) > 1 {
			title := v.Title
			if title == "" {
				title = v.ID
			}
			out = append(out, "# "+escape(title))
		}
		if v.Kind == bdf.ViewSheet {
			root := &node{kind: kindRoot, index: -1}
			tile := v.TileSize()
			keep := func(r *bdf.TextRun) bool { return r.X >= 0 && r.X < tile && r.Y >= 0 && r.Y < tile }
			for _, k := range tileKeys(v) {
				h, err := bdf.ParseHash(v.Tiles[k])
				if err != nil {
					return nil, fmt.Errorf("markdownbdf: view %s: tile %s: %w", v.ID, k, err)
				}
				if err := wr.extract(root, h, keep); err != nil {
					return nil, err
				}
			}
			wr.guess = guessCode([]*node{root})
			out = append(out, wr.sheet(root)...)
			if opts.Tables == "files" {
				// the tables of the sheet, as files
				for _, t := range doc.Tables {
					if t.View == v.ID && t.Page == 0 {
						if link := wr.tableFile(t); link != "" {
							out = append(out, link)
						}
					}
				}
			}
			continue
		}
		roots := make([]*node, len(v.Pages))
		for i, p := range v.Pages {
			roots[i] = &node{kind: kindRoot, index: -1}
			for _, l := range p.Layers {
				if skip[l.Role] {
					continue
				}
				if err := wr.extract(roots[i], l.Obj, nil); err != nil {
					return nil, err
				}
			}
		}
		if !opts.KeepRunning && v.Kind != bdf.ViewScroll {
			dropRunning(roots, v.Pages)
		}
		wr.guess = guessCode(roots)
		for _, root := range roots {
			out = append(out, wr.blocks(root.kids, false, 0)...)
		}
	}
	_, err = io.WriteString(w, strings.Join(out, "\n\n")+"\n")
	return wr.res, err
}

// WriteDir writes the document as dir/name.md, with its images beside it
// (the directory is made when it is missing). opts.Image is replaced.
func WriteDir(dir, name string, doc *bdf.Document, opts *Options) (*Result, error) {
	if opts == nil {
		opts = &Options{}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	o := *opts
	o.Image = func(file string, data []byte) (string, error) {
		return file, os.WriteFile(filepath.Join(dir, file), data, 0o644)
	}
	o.File = o.Image
	path := filepath.Join(dir, name+".md")
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	res, err := Write(f, doc, &o)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		os.Remove(path)
		return nil, err
	}
	return res, nil
}

// selectViews picks the views to write.
func selectViews(doc *bdf.Document, ids []string) ([]*bdf.View, error) {
	if len(ids) > 0 {
		var out []*bdf.View
		for _, id := range ids {
			found := false
			for _, v := range doc.Views {
				if v.ID == id {
					out, found = append(out, v), true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("markdownbdf: no view %q", id)
			}
		}
		return out, nil
	}
	var scroll, all []*bdf.View
	for _, v := range doc.Views {
		all = append(all, v)
		if v.Kind == bdf.ViewScroll {
			scroll = append(scroll, v)
		}
	}
	if len(scroll) > 0 {
		return scroll, nil
	}
	if len(all) == 0 {
		return nil, errors.New("markdownbdf: the document has no views")
	}
	return all, nil
}

// tileKeys returns the keys of the tiles of a sheet in reading order: by
// rows of tiles, left to right.
func tileKeys(v *bdf.View) []string {
	type key struct {
		x, y int
		s    string
	}
	keys := make([]key, 0, len(v.Tiles))
	for k := range v.Tiles {
		var kk key
		fmt.Sscanf(k, "%d,%d", &kk.x, &kk.y)
		kk.s = k
		keys = append(keys, kk)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].y != keys[j].y {
			return keys[i].y < keys[j].y
		}
		return keys[i].x < keys[j].x
	})
	out := make([]string, len(keys))
	for i, k := range keys {
		out[i] = k.s
	}
	return out
}

// writer holds what the writing of one document shares.
type writer struct {
	doc    *bdf.Document
	opts   *Options
	objs   map[bdf.Hash]*bdf.ObjectPart
	faces  map[bdf.Hash]*face
	images map[bdf.Hash]string // the path of each image written ("" when it could not be)
	csvs   map[string]string   // the path of the CSV file of each table written ("" when it could not be)
	res    *Result
	// guess says the monospaced paragraphs of the view are code (guessCode)
	guess bool
}

func (w *writer) warnf(format string, args ...any) {
	w.res.Warnings = append(w.res.Warnings, fmt.Sprintf(format, args...))
}

func (w *writer) resolve(h bdf.Hash) *bdf.ObjectPart {
	if o, ok := w.objs[h]; ok {
		return o
	}
	var o *bdf.ObjectPart
	if p := w.doc.Part(h); p != nil {
		if d, err := bdf.DecodeObject(p.Data); err == nil {
			o = d
		}
	}
	w.objs[h] = o
	return o
}

// The kinds of the nodes of the tree, besides those of the structure.
const (
	kindRoot  bdf.NodeKind = 254
	kindImage bdf.NodeKind = 255
)

// node is a node of the tree of a page: a node of the structure, an image,
// or the root.
type node struct {
	kind   bdf.NodeKind
	parent *node
	kids   []*node
	level  int
	cell   *bdf.CellRef
	alt    string
	data   string // the id of the table a figure is the data of
	code   bool   // a line of code (the payload's word)
	quote  int    // how deep in quotations (the payload's words)
	list   string // "bullet", "ordered" or "" (the payload's word)
	runs   []*run
	image  bdf.Hash
	// index is the index of the node in the content it came from, offset
	// by the nodes before it in the tree; -1 for the root and images.
	index int
}

func (n *node) add(k *node) {
	k.parent = n
	n.kids = append(n.kids, k)
}

func (n *node) insert(at int, k *node) {
	k.parent = n
	n.kids = append(n.kids, nil)
	copy(n.kids[at+1:], n.kids[at:])
	n.kids[at] = k
}

// leaf reports whether a node holds runs rather than nodes.
func (n *node) leaf() bool {
	switch n.kind {
	case bdf.NodeParagraph, bdf.NodeHeading, bdf.NodeCaption:
		return true
	case bdf.NodeCell:
		return n.parent == nil || n.parent.kind != bdf.NodeTable
	}
	return false
}

// run is a run of text with what the Markdown makes of its font and the
// link over it.
type run struct {
	bdf.TextRun
	bold, italic, mono bool
	link               string
}

// extract adds the content of a top-level object to a tree. keep, when
// set, tells the runs that belong to the object (a tile of a sheet keeps
// the runs whose anchor lies in it).
func (w *writer) extract(root *node, h bdf.Hash, keep func(*bdf.TextRun) bool) error {
	o := w.resolve(h)
	if o == nil {
		return fmt.Errorf("markdownbdf: object %s missing or invalid", h)
	}
	c, err := bdf.ExtractContent(o, w.resolve)
	if err != nil {
		return err
	}
	base := 0
	for _, k := range root.kids {
		if k.index >= base {
			base = k.index + 1
		}
	}
	// (the tree may hold nodes of tiles before this one)
	base = max(base, countNodes(root))
	nodes := make([]*node, len(c.Nodes))
	for i, tn := range c.Nodes {
		n := &node{kind: tn.Kind, level: tn.Level, cell: tn.Cell, alt: tn.Alt, data: tn.Data, code: tn.Code, quote: tn.Quote, list: tn.List, index: base + i}
		nodes[i] = n
		parent := root
		if tn.Parent >= 0 {
			parent = nodes[tn.Parent]
		}
		parent.add(n)
	}
	var all []*run
	for i := range c.Runs {
		r := &c.Runs[i]
		if keep != nil && !keep(r) {
			continue
		}
		rr := &run{TextRun: *r}
		rr.bold, rr.italic, rr.mono = w.style(r.Font)
		nodes[r.Node].runs = append(nodes[r.Node].runs, rr)
		all = append(all, rr)
	}
	for _, l := range c.Links {
		if strings.HasPrefix(l.URL, "#") && !w.opts.InternalLinks {
			continue
		}
		for _, r := range all {
			if r.link != "" {
				continue
			}
			// the runs a LINK covers are its text: those anchored in it
			px, py := r.X+r.Size*0.3, r.Y-r.Size*0.35
			if px >= l.X && px <= l.X+l.W && py >= l.Y && py <= l.Y+l.H {
				r.link = l.URL
			}
		}
	}
	for _, im := range c.Images {
		n := &node{kind: kindImage, image: im.Hash, index: -1}
		switch {
		case im.Node >= 0 && nodes[im.Node].leaf():
			// after the leaf it interrupts, and the images there already
			leaf := nodes[im.Node]
			p := leaf.parent
			at := len(p.kids)
			for i, k := range p.kids {
				if k == leaf {
					at = i + 1
					for at < len(p.kids) && p.kids[at].kind == kindImage {
						at++
					}
					break
				}
			}
			p.insert(at, n)
		default:
			p := root
			if im.Node >= 0 {
				p = nodes[im.Node]
			}
			at := len(p.kids)
			for i, k := range p.kids {
				if k.index >= base+im.Before {
					at = i
					break
				}
			}
			p.insert(at, n)
		}
	}
	if len(c.Nodes) == 1 && c.Nodes[0].Kind == bdf.NodeParagraph && len(nodes[0].runs) > 1 {
		// an object without MARKs (an untagged PDF): guess its structure
		guessStructure(root, nodes[0])
	}
	prune(root)
	return nil
}

func countNodes(n *node) int {
	c := 0
	for _, k := range n.kids {
		if k.index >= 0 {
			c++
		}
		c += countNodes(k)
	}
	return c
}

// prune drops the nodes that hold nothing: the cells of a tile that
// belong to another, empty paragraphs.
func prune(n *node) {
	kept := n.kids[:0]
	for _, k := range n.kids {
		prune(k)
		if len(k.kids) == 0 && len(k.runs) == 0 && k.kind != kindImage && !(k.kind == bdf.NodeFigure && k.alt != "") {
			continue
		}
		kept = append(kept, k)
	}
	n.kids = kept
}

// guessStructure splits the one paragraph of an object without MARKs by the
// positions of its runs, and takes the lines in a larger or bolder font
// for headings: what an untagged PDF leaves to guess.
func guessStructure(root, p *node) {
	runs := p.runs
	sizes := make([]float64, len(runs))
	for i, r := range runs {
		sizes[i] = float64(r.Size)
	}
	sort.Float64s(sizes)
	body := sizes[len(sizes)/2]
	var paras []*node
	var cur *node
	var ref *run // the last run on the baseline: a superscript or subscript is not a new line
	for _, r := range runs {
		split := cur == nil
		if ref == nil {
			ref = r
		} else {
			dy := float64(r.Y - ref.Y)
			small := r.Size < ref.Size*0.85 && abs(dy) < float64(ref.Size)
			if !small {
				if abs(dy) > float64(ref.Size)*1.6 || (dy > float64(ref.Size)*0.5 && (r.Size != ref.Size || r.bold != ref.bold)) {
					split = true
				}
				ref = r
			}
		}
		if split {
			cur = &node{kind: bdf.NodeParagraph, index: p.index}
			paras = append(paras, cur)
			r.Sep = bdf.SepBreak
		}
		cur.runs = append(cur.runs, r)
	}
	for _, q := range paras {
		allBold, size, n := true, float64(0), 0
		for _, r := range q.runs {
			allBold = allBold && r.bold
			size += float64(r.Size)
			n += len([]rune(r.Text))
		}
		size /= float64(len(q.runs))
		switch {
		case n > 120:
		case size >= body*1.5:
			q.kind, q.level = bdf.NodeHeading, 1
		case size >= body*1.2:
			q.kind, q.level = bdf.NodeHeading, 2
		case allBold && size > body*1.02 && len(q.runs) <= 3:
			q.kind, q.level = bdf.NodeHeading, 3
		}
	}
	// the paragraphs take the place of the one they were
	kids := make([]*node, 0, len(root.kids)+len(paras))
	for _, k := range root.kids {
		if k != p {
			kids = append(kids, k)
			continue
		}
		for _, q := range paras {
			q.parent = root
			kids = append(kids, q)
		}
	}
	root.kids = kids
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

// frontMatter writes the Dublin Core elements of the document as YAML.
func frontMatter(dc *bdf.DublinCore) string {
	var lines []string
	for _, name := range []string{"title", "creator", "contributor", "publisher", "subject", "description", "date", "created", "modified", "language", "identifier", "rights"} {
		vs := *dc.Field(name)
		switch len(vs) {
		case 0:
		case 1:
			lines = append(lines, name+": "+yamlScalar(vs[0]))
		default:
			items := make([]string, len(vs))
			for i, v := range vs {
				items[i] = yamlScalar(v)
			}
			lines = append(lines, name+": ["+strings.Join(items, ", ")+"]")
		}
	}
	if len(lines) == 0 {
		return ""
	}
	return "---\n" + strings.Join(lines, "\n") + "\n---"
}

// yamlScalar writes a string as a YAML scalar: as it is when nothing in it
// means something else to YAML, else quoted.
func yamlScalar(s string) string {
	plain := s != "" && !strings.ContainsAny(s, ":#\"'[]{},&*!|>%@`\\") && strings.TrimSpace(s) == s && s[0] != '-' && s[0] != '?'
	if plain {
		switch strings.ToLower(s) {
		case "true", "false", "yes", "no", "on", "off", "null", "~":
			plain = false
		}
		if _, err := strconv.ParseFloat(s, 64); err == nil {
			plain = false
		}
	}
	if plain {
		return s
	}
	return strconv.Quote(s)
}

// The margin of a page that running heads and feet are looked for in, as
// a fraction of its height, and the height of the places that count as
// the same.
const (
	marginFrac = 0.1
	placeFrac  = 0.02
)

var (
	digitsRe = regexp.MustCompile(`\d+`)
	// pageNumberRe matches a page number alone: "3", "- 3 -", "Page 3",
	// "3 / 12", "3 of 12", "p. 3", "3 ページ", "iv".
	pageNumberRe = regexp.MustCompile(`^(?i)(page|p\.|pg\.?|第)?\s*\d{1,4}\s*(ページ|頁)?(\s*(of|/|／)\s*\d{1,4}\s*(ページ|頁)?)?$|^(?i)[ivxlcdm]{1,7}$`)
)

// dropRunning drops the running heads and feet of the pages of a view: the
// paragraphs in the top or bottom margin of a page whose text, digits
// aside, is at the same place on another page too, the page numbers alone
// in a margin, and, on the pages after the first, the paragraphs at a
// place where those were found.
func dropRunning(roots []*node, pages []*bdf.Page) {
	type key struct {
		band, place int
		text        string
	}
	type cand struct {
		n      *node
		k      key
		number bool
	}
	counts := map[key]int{}
	cands := make([][]cand, len(roots))
	for i, root := range roots {
		h := pages[i].H
		if h <= 0 {
			continue
		}
		for _, k := range root.kids {
			if k.kind != bdf.NodeParagraph && k.kind != bdf.NodeHeading || len(k.runs) == 0 {
				continue
			}
			top, bottom := float32(math.Inf(1)), float32(math.Inf(-1))
			var sb strings.Builder
			for j, r := range k.runs {
				top, bottom = min(top, r.Y-r.Size*0.8), max(bottom, r.Y)
				if j > 0 && r.Sep != bdf.SepNone {
					sb.WriteByte(' ')
				}
				sb.WriteString(r.Text)
			}
			band := 0
			switch {
			case bottom <= h*marginFrac:
				band = -1
			case top >= h*(1-marginFrac):
				band = 1
			default:
				continue
			}
			text := strings.TrimSpace(sb.String())
			if text == "" || len([]rune(text)) > 120 {
				continue
			}
			c := cand{n: k, k: key{band, int(math.Round(float64(k.runs[0].Y / (h * placeFrac)))), digitsRe.ReplaceAllString(text, "#")}}
			c.number = pageNumberRe.MatchString(strings.Trim(text, " -–—·|[]()"))
			counts[c.k]++
			cands[i] = append(cands[i], c)
		}
	}
	// a place that running text was found at holds running text on the
	// other pages too (the heads of odd and even pages differ)
	type place struct{ band, place int }
	places := map[place]bool{}
	for _, cs := range cands {
		for _, c := range cs {
			if c.number || counts[c.k] >= 2 {
				places[place{c.k.band, c.k.place}] = true
			}
		}
	}
	for i, cs := range cands {
		drop := map[*node]bool{}
		for _, c := range cs {
			if c.number || counts[c.k] >= 2 || i > 0 && places[place{c.k.band, c.k.place}] {
				drop[c.n] = true
			}
		}
		if len(drop) == 0 {
			continue
		}
		kept := roots[i].kids[:0]
		for _, k := range roots[i].kids {
			if !drop[k] {
				kept = append(kept, k)
			}
		}
		roots[i].kids = kept
	}
}
