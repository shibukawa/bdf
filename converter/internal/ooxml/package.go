// Package ooxml reads Office Open XML documents (.pptx, .docx, .xlsx):
// the parts of their Open Packaging Conventions container (ECMA-376 Part 2),
// the relationships between the parts, and part XML as generic element
// trees.
package ooxml

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/shibukawa/bdf/internal/xmltree"
)

// maxPartSize bounds the decompressed size of one package part.
const maxPartSize = 1 << 30

// Package is an Open Packaging Conventions container (the zip file of an
// Office document). Part names are matched case-insensitively. A nil
// *Package is an empty package, for markup that comes without one (such as
// the single XML file of a Visio 2003 drawing).
type Package struct {
	// Supported, when set, makes the parts' mc:AlternateContent resolve to
	// the choices whose required namespaces it accepts (see ParseChoosing).
	Supported func(prefix string) bool
	// Choose, when set, makes mc:AlternateContent resolve to the choices it
	// accepts as well (such as MathChoice).
	Choose func(choice *Node) bool

	files map[string]*zip.File // by lower-cased part name without leading slash
	xmls  map[string]*Node
	rels  map[string]map[string]Rel
}

// Rel is a relationship with its target resolved to a part name (or kept
// as is when external).
type Rel struct {
	ID, Type, Target string
	External         bool
}

// Open reads the zip directory of a package.
func Open(r io.ReaderAt, size int64) (*Package, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	p := &Package{files: map[string]*zip.File{}, xmls: map[string]*Node{}, rels: map[string]map[string]Rel{}}
	for _, f := range zr.File {
		// some producers write Windows separators in the zip (xl\workbook.xml)
		name := strings.ReplaceAll(f.Name, "\\", "/")
		p.files[strings.ToLower(strings.TrimPrefix(name, "/"))] = f
	}
	return p, nil
}

// Has reports whether the package has a part.
func (p *Package) Has(name string) bool {
	if p == nil {
		return false
	}
	_, ok := p.files[strings.ToLower(strings.TrimPrefix(name, "/"))]
	return ok
}

// PartSize reports the uncompressed size recorded in the ZIP directory.
func (p *Package) PartSize(name string) (uint64, bool) {
	if p == nil {
		return 0, false
	}
	f, ok := p.files[strings.ToLower(strings.TrimPrefix(name, "/"))]
	if !ok {
		return 0, false
	}
	return f.UncompressedSize64, true
}

// Read returns the bytes of a part.
func (p *Package) Read(name string) ([]byte, error) {
	rc, err := p.OpenPart(name)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return b, nil
}

// OpenPart opens a part without keeping its uncompressed bytes in memory.
// The reader enforces the same size limit as Read, including when the ZIP
// directory understates the uncompressed size.
func (p *Package) OpenPart(name string) (io.ReadCloser, error) {
	if p == nil {
		return nil, fmt.Errorf("missing part %s", name)
	}
	f, ok := p.files[strings.ToLower(strings.TrimPrefix(name, "/"))]
	if !ok {
		return nil, fmt.Errorf("missing part %s", name)
	}
	if f.UncompressedSize64 > maxPartSize {
		return nil, fmt.Errorf("%s: part too large", name)
	}
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	return &limitedPart{ReadCloser: rc, remaining: maxPartSize, name: name}, nil
}

type limitedPart struct {
	io.ReadCloser
	remaining int64
	name      string
}

func (r *limitedPart) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	if r.remaining == 0 {
		var one [1]byte
		n, err := r.ReadCloser.Read(one[:])
		if n > 0 {
			return 0, fmt.Errorf("%s: part too large", r.name)
		}
		return 0, err
	}
	if int64(len(b)) > r.remaining {
		b = b[:int(r.remaining)]
	}
	n, err := r.ReadCloser.Read(b)
	r.remaining -= int64(n)
	return n, err
}

// XML returns the parsed XML of a part (cached).
func (p *Package) XML(name string) (*Node, error) {
	if p == nil {
		return nil, fmt.Errorf("missing part %s", name)
	}
	key := strings.ToLower(strings.TrimPrefix(name, "/"))
	if n, ok := p.xmls[key]; ok {
		return n, nil
	}
	rc, err := p.OpenPart(name)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	n, err := ParsePickingReader(rc, p.choice())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	if _, err := io.Copy(io.Discard, rc); err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	p.xmls[key] = n
	return n, nil
}

// ReadElement reads one streamed element with this package's markup
// compatibility choices, as XML does for a whole part.
func (p *Package) ReadElement(d *xml.Decoder, start xml.StartElement) (*Node, error) {
	return xmltree.ReadElementPicking(d, start, p.choice())
}

func (p *Package) choice() func(*Node) bool {
	pick := xmltree.SupportedChoice(p.Supported)
	if p.Choose != nil {
		supported := pick
		pick = func(c *Node) bool { return supported != nil && supported(c) || p.Choose(c) }
	}
	return pick
}

// DiscardXML releases a cached tree after its caller has built a more compact
// representation. Nodes still referenced by that representation stay valid.
func (p *Package) DiscardXML(name string) {
	if p != nil {
		delete(p.xmls, strings.ToLower(strings.TrimPrefix(name, "/")))
	}
}

// Rels returns the relationships of a part by ID ("" for the package's
// own relationships).
func (p *Package) Rels(part string) map[string]Rel {
	if p == nil {
		return nil
	}
	part = strings.TrimPrefix(part, "/")
	if m, ok := p.rels[part]; ok {
		return m
	}
	dir, base := path.Split(part)
	m := map[string]Rel{}
	p.rels[part] = m
	root, err := p.XML(dir + "_rels/" + base + ".rels")
	if err != nil {
		return m
	}
	for _, r := range root.Children("Relationship") {
		rl := Rel{ID: r.AttrStr("Id", ""), Type: r.AttrStr("Type", ""), Target: r.AttrStr("Target", "")}
		if strings.EqualFold(r.AttrStr("TargetMode", ""), "External") {
			rl.External = true
		} else {
			rl.Target = ResolvePart(dir, rl.Target)
		}
		m[rl.ID] = rl
	}
	return m
}

// ResolvePart resolves a relationship target against the source part's directory.
func ResolvePart(dir, target string) string {
	if strings.HasPrefix(target, "/") {
		return path.Clean(target)[1:]
	}
	return strings.TrimPrefix(path.Clean("/"+dir+target), "/")
}

// Target returns the part a relationship ID points to.
func (p *Package) Target(part, id string) (Rel, bool) {
	r, ok := p.Rels(part)[id]
	return r, ok && r.Target != ""
}

// RelOfType returns the first relationship whose type ends with suffix
// (e.g. "/slideLayout"), in ID order for determinism.
func (p *Package) RelOfType(part, suffix string) (Rel, bool) {
	var best Rel
	found := false
	for _, r := range p.Rels(part) {
		if strings.HasSuffix(r.Type, suffix) && (!found || relIDLess(r.ID, best.ID)) {
			best, found = r, true
		}
	}
	return best, found
}

// relIDLess orders "rId2" before "rId10".
func relIDLess(a, b string) bool {
	if len(a) != len(b) {
		return len(a) < len(b)
	}
	return a < b
}
