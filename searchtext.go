package bdf

import (
	"sort"
	"strings"
)

// SearchText is the text of a document for a search index outside it: its
// metadata and the text of each view, page by page.
type SearchText struct {
	Meta  Meta       `json:"meta"`
	Views []ViewText `json:"views"`
}

// ViewText is the text of a view.
type ViewText struct {
	ID    string     `json:"id"`
	Kind  string     `json:"kind"`
	Title string     `json:"title,omitempty"`
	Pages []PageText `json:"pages"`
}

// PageText is the text of a page (numbered from 1, for links like
// #page=N), or of a whole sheet (Page 0).
type PageText struct {
	Page int    `json:"page,omitempty"`
	Text string `json:"text"`
}

// SearchText returns the text of the document for a search index: that of
// the views' text index parts (docs/spec.md §7.9), or extracted from the
// objects of a view without one. Pages without text are left out, and so
// is a scroll view when the document also has a flow view (Word, Markdown
// and EPUB documents may lay the same text out both ways).
//
// The text is the document's in plain text; a document that was encrypted
// should not have it stored where the password does not protect it.
func (d *Document) SearchText() (*SearchText, error) {
	out := &SearchText{Meta: d.Meta, Views: []ViewText{}}
	hasFlow := false
	for _, v := range d.Views {
		hasFlow = hasFlow || v.Kind == ViewFlow
	}
	for _, v := range d.Views {
		if v.Kind == ViewScroll && hasFlow {
			continue
		}
		runs, err := d.viewRuns(v)
		if err != nil {
			return nil, err
		}
		vt := ViewText{ID: v.ID, Kind: v.Kind, Title: v.Title, Pages: []PageText{}}
		if v.Kind == ViewSheet {
			if t := strings.TrimSpace(PlainText(runs)); t != "" {
				vt.Pages = append(vt.Pages, PageText{Text: t})
			}
		} else {
			// the runs of a page are in order, those of its layers one after another
			byPage := map[uint32][]IndexRun{}
			var pages []uint32
			for _, r := range runs {
				if _, ok := byPage[r.A]; !ok {
					pages = append(pages, r.A)
				}
				byPage[r.A] = append(byPage[r.A], r)
			}
			sort.Slice(pages, func(i, j int) bool { return pages[i] < pages[j] })
			for _, p := range pages {
				if t := strings.TrimSpace(PlainText(byPage[p])); t != "" {
					vt.Pages = append(vt.Pages, PageText{Page: int(p) + 1, Text: t})
				}
			}
		}
		out.Views = append(out.Views, vt)
	}
	return out, nil
}

// viewRuns returns the runs of a view's text index part, or extracts them.
func (d *Document) viewRuns(v *View) ([]IndexRun, error) {
	if v.TextIndex != "" {
		if h, err := ParseHash(v.TextIndex); err == nil {
			if p := d.parts[h]; p != nil {
				return DecodeTextIndex(p.Data)
			}
		}
	}
	return d.indexRuns(v)
}
