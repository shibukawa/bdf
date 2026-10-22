package bdf

import "fmt"

// View kinds.
const (
	ViewFixed  = "fixed"
	ViewFlow   = "flow"
	ViewSheet  = "sheet"
	ViewScroll = "scroll"
)

// Page progression directions of fixed and flow views.
const (
	DirectionLTR = "ltr"
	DirectionRTL = "rtl"
)

// Layer roles.
const (
	RoleBackground = "background"
	RoleMaster     = "master"
	RoleHeader     = "header"
	RoleFooter     = "footer"
	RoleBody       = "body"
	RoleNotes      = "notes"
	RoleAnnotation = "annotation"
)

// Part types.
const (
	PartObject = "obj"
	PartFont   = "font"
	PartImage  = "img"
	PartPath   = "path"
	PartIndex  = "idx"
	// PartTable is a table the document holds beside its drawing: an
	// Apache Parquet file (docs/spec.md §4.6).
	PartTable = "table"
	// PartSeq is a Standard MIDI File that a view plays (docs/spec.md §4.4).
	PartSeq = "seq"
	// PartAudio is an audio file that a view plays, as it is (docs/spec.md
	// §4.4).
	PartAudio = "audio"
	// PartTemplate describes a compiled template of forms as JSON, and
	// PartSource is the template it was compiled from (package form;
	// docs/spec.md §4.5). No view refers to them: a reader of documents
	// passes them over.
	PartTemplate = "tpl"
	PartSource   = "src"
	// PartFill is the record of a form completed from a partly filled
	// form whose values were signed: the values, and the signatures of
	// those who gave them, as JSON (package form; docs/spec.md §4.5). No
	// view refers to it.
	PartFill = "fill"
	// PartText is the print image of a form filled from a template of
	// plain text: its grid of columns and lines, and the characters of
	// each line, as JSON (package form/textpart; docs/spec.md §4). No view
	// refers to it: a reader of documents passes it over, and the writers
	// of the text modes of printers (ESC/P, ESC/POS) print from it.
	PartText = "text"
	// PartSealed is the type of the parts of an encrypted document's outer
	// manifest (docs/spec.md §3.5).
	PartSealed = "sealed"
)

// Encodings.
const (
	EncIdentity   = "identity"
	EncDeflateRaw = "deflate-raw"
)

// Manifest is the JSON document directory.
//
// An encrypted document stores an outer manifest that holds only BDF,
// Encryption and the sealed parts; the manifest it is read with is sealed
// (docs/spec.md §3.5).
type Manifest struct {
	BDF        int         `json:"bdf"`
	Encryption *Encryption `json:"encryption,omitempty"`
	Opset      int         `json:"opset,omitempty"`
	Unit       string      `json:"unit,omitempty"`
	Meta       Meta        `json:"meta,omitzero"`
	// Segment, in a segment document, names the pages it carries
	// (docs/spec.md §3.6).
	Segment *Segment `json:"segment,omitempty"`
	Views   []*View  `json:"views,omitempty"`
	// Tables are the tables the document holds beside its drawing
	// (docs/spec.md §4.6).
	Tables []*Table    `json:"tables,omitempty"`
	Parts  []PartEntry `json:"parts"`
}

// Table is a table a document holds beside its drawing (docs/spec.md
// §4.6): the data of a chart, the cells of a sheet, the rows of a CSV or
// Parquet file, as an Apache Parquet file in a PartTable part.
type Table struct {
	// ID names the table in the document ("chart-1", "sheet1", "data").
	ID    string `json:"id"`
	Title string `json:"title,omitempty"`
	// Part is the hash of the part that holds the Parquet file.
	Part string `json:"part"`
	// Rows is the count of rows, and Columns the columns, as the Parquet
	// file's schema has them: a listing without decoding the file.
	Rows    int64         `json:"rows"`
	Columns []TableColumn `json:"columns"`
	// View and Page tell where the figure the table is the data of is:
	// the view's id and the page in it, from 1 (0 for a view without
	// pages: the table of a sheet); "" when the table is not shown.
	View string `json:"view,omitempty"`
	Page int    `json:"page,omitempty"`
}

// TableColumn is a column of a table: its name, and its type ("string",
// "number", "integer", "boolean", "date", "datetime" or "other";
// encoding/parquet.Type).
type TableColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Table returns the table with an id, or nil.
func (m *Manifest) Table(id string) *Table {
	for _, t := range m.Tables {
		if t != nil && t.ID == id {
			return t
		}
	}
	return nil
}

// Meta holds document metadata.
type Meta struct {
	// DC describes the document itself (title, creator, dates, ...).
	DC DublinCore `json:"dc,omitzero"`
	// Source is the input format the document was converted from ("pdf", "pptx", ...).
	Source    string `json:"source,omitempty"`
	Generator string `json:"generator,omitempty"`
}

// View is a viewing unit: a slide deck, a document, or a sheet.
type View struct {
	ID    string `json:"id"`
	Kind  string `json:"kind"`
	Title string `json:"title,omitempty"`

	// optional text index part (docs/spec.md §7.9)
	TextIndex string `json:"textIndex,omitempty"`

	// fixed / flow
	Pages      []*Page     `json:"pages,omitempty"`
	Continuous *Continuous `json:"continuous,omitempty"`
	// Direction is the order pages are put side by side in (spreads, a
	// horizontal row): DirectionLTR (the default when "") or DirectionRTL,
	// for books bound on the right (docs/spec.md §4.1).
	Direction string `json:"direction,omitempty"`

	// sheet
	Tile      float32           `json:"tile,omitempty"`
	Cols      []Run             `json:"cols,omitempty"`
	Rows      []Run             `json:"rows,omitempty"`
	Freeze    *Freeze           `json:"freeze,omitempty"`
	Gridlines bool              `json:"gridlines,omitempty"`
	Tiles     map[string]string `json:"tiles,omitempty"`
	TilesRef  string            `json:"tilesRef,omitempty"`

	// Play is the music a view plays (docs/spec.md §4.4).
	Play *Play `json:"play,omitempty"`
}

// Play is what a view plays: the music of a score as a Standard MIDI File
// (Seq), or a recording as the audio file itself (Audio), and cues that tie
// its time to places on the pages (docs/spec.md §4.4). A view has one of
// Seq and Audio.
type Play struct {
	// Seq is the hash of the PartSeq part.
	Seq string `json:"seq,omitempty"`
	// Audio is the hash of the PartAudio part, and Type the media type of
	// the file it holds ("audio/mpeg").
	Audio string `json:"audio,omitempty"`
	Type  string `json:"type,omitempty"`
	// Cues is the hash of the cue index part (PartIndex), if any. Its
	// ticks are those of the Standard MIDI File, or milliseconds of the
	// audio.
	Cues string `json:"cues,omitempty"`
}

// Run is a run-length entry [count, size] for sheet rows/columns.
type Run [2]float32

// Freeze describes frozen panes.
type Freeze struct {
	Cols int `json:"cols,omitempty"`
	Rows int `json:"rows,omitempty"`
}

// Continuous configures continuous mode of a flow view.
type Continuous struct {
	Gap float32 `json:"gap"`
}

// Page is a fixed-size page.
type Page struct {
	W      float32  `json:"w"`
	H      float32  `json:"h"`
	Body   *RectDef `json:"body,omitempty"`
	Layers []Layer  `json:"layers"`
}

// RectDef is a rectangle in JSON.
type RectDef struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	W float32 `json:"w"`
	H float32 `json:"h"`
}

// Layer is one drawable layer of a page.
type Layer struct {
	Role string `json:"role"`
	Obj  Hash   `json:"obj"`
}

// PartEntry describes one part in the manifest.
type PartEntry struct {
	H    Hash   `json:"h"`
	T    string `json:"t"`
	Enc  string `json:"enc"`
	Len  int    `json:"len"`  // stored (possibly compressed) length
	Size int    `json:"size"` // decoded length
	Off  int64  `json:"off"`  // offset from the start of the parts region (single form)
	// Sealed names the part of the outer manifest that holds this part
	// sealed (encrypted documents only).
	Sealed Hash `json:"sealed,omitzero"`
}

// Limits of the numbers of a manifest. They decide how long the loops of a
// reader run and how much it allocates (the tiles of a region, the
// gridlines of a sheet, the bands of a continuous layout), so a reader
// checks them once, when it opens the document. What they refuse no writer
// has a use for.
const (
	// MaxPageSize bounds the sides of pages and of their bodies, in units:
	// the coordinates of objects are f32, which cannot tell whole units
	// apart beyond it.
	MaxPageSize = 1 << 24
	// MaxSheetEntries bounds the rows of a sheet, and its columns (Excel
	// has 1,048,576 rows).
	MaxSheetEntries = 1 << 24
	// MinEntrySize is the smallest size of a row or a column that has one
	// (0: hidden), in units.
	MinEntrySize = 1.0 / 64
	// DefaultTile is the tile size of a sheet view that states none.
	DefaultTile = 2048
	// MaxTables bounds the tables of a document, and MaxTableColumns the
	// columns of a table (docs/spec.md §4.6).
	MaxTables       = 1 << 16
	MaxTableColumns = 1 << 14
)

// TileSize returns the size of the tiles of a sheet view: Tile, or
// DefaultTile for a view that states none.
func (v *View) TileSize() float32 {
	if v.Tile > 0 {
		return v.Tile
	}
	return DefaultTile
}

// check reports what is wrong with a manifest that was read: unsupported
// versions, null views and pages, and numbers out of range.
func (m *Manifest) check() error {
	if m.BDF != FormatVersion {
		return &FormatError{Msg: fmt.Sprintf("unsupported manifest format version %d", m.BDF)}
	}
	if m.Opset < 0 || m.Opset > OpsetVersion {
		return &FormatError{Msg: fmt.Sprintf("unsupported manifest opset %d", m.Opset)}
	}
	size := func(v float32) bool { return v >= 0 && v <= MaxPageSize } // false for NaN
	for _, v := range m.Views {
		if v == nil {
			return &FormatError{Msg: "null view"}
		}
		for _, p := range v.Pages {
			switch {
			case p == nil:
				return &FormatError{Msg: fmt.Sprintf("view %s: null page", v.ID)}
			case !size(p.W) || !size(p.H):
				return &FormatError{Msg: fmt.Sprintf("view %s: page size %g × %g out of range", v.ID, p.W, p.H)}
			}
			if b := p.Body; b != nil && !(size(b.W) && size(b.H) && size(max(b.X, -b.X)) && size(max(b.Y, -b.Y))) {
				return &FormatError{Msg: fmt.Sprintf("view %s: page body out of range", v.ID)}
			}
		}
		if c := v.Continuous; c != nil && !size(max(c.Gap, -c.Gap)) {
			return &FormatError{Msg: fmt.Sprintf("view %s: gap %g out of range", v.ID, c.Gap)}
		}
		if !(v.Tile <= 0 || v.Tile >= 1 && v.Tile <= MaxPageSize) {
			return &FormatError{Msg: fmt.Sprintf("view %s: tile size %g out of range", v.ID, v.Tile)}
		}
		for _, runs := range [][]Run{v.Cols, v.Rows} {
			total := 0
			for _, r := range runs {
				n := int(r[0])
				if !(r[0] >= 0 && r[0] <= MaxSheetEntries) || float32(n) != r[0] {
					return &FormatError{Msg: fmt.Sprintf("view %s: %g rows or columns", v.ID, r[0])}
				}
				if total += n; total > MaxSheetEntries {
					return &FormatError{Msg: fmt.Sprintf("view %s: more than %d rows or columns", v.ID, MaxSheetEntries)}
				}
				if !size(r[1]) || r[1] > 0 && r[1] < MinEntrySize {
					return &FormatError{Msg: fmt.Sprintf("view %s: rows or columns of size %g", v.ID, r[1])}
				}
			}
		}
		if f := v.Freeze; f != nil && (f.Cols < 0 || f.Rows < 0) {
			return &FormatError{Msg: fmt.Sprintf("view %s: bad frozen panes", v.ID)}
		}
	}
	if len(m.Tables) > MaxTables {
		return &FormatError{Msg: fmt.Sprintf("%d tables, more than %d", len(m.Tables), MaxTables)}
	}
	ids := make(map[string]bool, len(m.Tables))
	for _, t := range m.Tables {
		switch {
		case t == nil:
			return &FormatError{Msg: "null table"}
		case t.ID == "" || ids[t.ID]:
			return &FormatError{Msg: fmt.Sprintf("table %q: no id, or one of another table", t.ID)}
		case t.Rows < 0 || len(t.Columns) > MaxTableColumns || t.Page < 0:
			return &FormatError{Msg: fmt.Sprintf("table %s: numbers out of range", t.ID)}
		}
		if _, err := ParseHash(t.Part); err != nil {
			return &FormatError{Msg: fmt.Sprintf("table %s: bad part %q", t.ID, t.Part)}
		}
		ids[t.ID] = true
	}
	return nil
}

// AddPage appends a page to a fixed/flow view.
func (v *View) AddPage(w, h float32, layers ...Layer) *Page {
	p := &Page{W: w, H: h, Layers: layers}
	v.Pages = append(v.Pages, p)
	return p
}
