package bdf

// View kinds.
const (
	ViewFixed = "fixed"
	ViewFlow  = "flow"
	ViewSheet = "sheet"
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
)

// Encodings.
const (
	EncIdentity   = "identity"
	EncDeflateRaw = "deflate-raw"
)

// Manifest is the JSON document directory.
type Manifest struct {
	BDF   int         `json:"bdf"`
	Opset int         `json:"opset"`
	Unit  string      `json:"unit"`
	Meta  Meta        `json:"meta,omitempty"`
	Views []*View     `json:"views"`
	Parts []PartEntry `json:"parts"`
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

	// sheet
	Tile      float32           `json:"tile,omitempty"`
	Cols      []Run             `json:"cols,omitempty"`
	Rows      []Run             `json:"rows,omitempty"`
	Freeze    *Freeze           `json:"freeze,omitempty"`
	Gridlines bool              `json:"gridlines,omitempty"`
	Tiles     map[string]string `json:"tiles,omitempty"`
	TilesRef  string            `json:"tilesRef,omitempty"`
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
}

// AddPage appends a page to a fixed/flow view.
func (v *View) AddPage(w, h float32, layers ...Layer) *Page {
	p := &Page{W: w, H: h, Layers: layers}
	v.Pages = append(v.Pages, p)
	return p
}
