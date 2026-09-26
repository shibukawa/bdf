// Package ooxml reads Office Open XML documents (.pptx, .docx, .xlsx):
// the parts of their Open Packaging Conventions container (ECMA-376 Part 2),
// the relationships between the parts, and part XML as generic element
// trees.
package ooxml

import (
	"archive/zip"
	"fmt"
	"io"
	"path"
	"strings"
)

// maxPartSize bounds the decompressed size of one package part.
const maxPartSize = 1 << 30

// Package is an Open Packaging Conventions container (the zip file of an
// Office document). Part names are matched case-insensitively. A nil
// *Package is an empty package, for markup that comes without one (such as
// the single XML file of a Visio 2003 drawing).
type Package struct {
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
		p.files[strings.ToLower(strings.TrimPrefix(f.Name, "/"))] = f
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

// Read returns the bytes of a part.
func (p *Package) Read(name string) ([]byte, error) {
	if p == nil {
		return nil, fmt.Errorf("missing part %s", name)
	}
	f, ok := p.files[strings.ToLower(strings.TrimPrefix(name, "/"))]
	if !ok {
		return nil, fmt.Errorf("missing part %s", name)
	}
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	b, err := io.ReadAll(io.LimitReader(rc, maxPartSize+1))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	if len(b) > maxPartSize {
		return nil, fmt.Errorf("%s: part too large", name)
	}
	return b, nil
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
	b, err := p.Read(name)
	if err != nil {
		return nil, err
	}
	n, err := Parse(b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	p.xmls[key] = n
	return n, nil
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
