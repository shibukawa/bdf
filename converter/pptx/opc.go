package pptx

import (
	"archive/zip"
	"fmt"
	"io"
	"path"
	"strings"
)

// maxPartSize bounds the decompressed size of one package part.
const maxPartSize = 1 << 30

// pkg is an Open Packaging Conventions container (the zip of a .pptx).
type pkg struct {
	files map[string]*zip.File // by lower-cased part name without leading slash
	xmls  map[string]*node
	rels  map[string]map[string]rel
}

// rel is a relationship with its target resolved to a part name (or kept
// as is when external).
type rel struct {
	ID, Type, Target string
	External         bool
}

func openPkg(r io.ReaderAt, size int64) (*pkg, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	p := &pkg{files: map[string]*zip.File{}, xmls: map[string]*node{}, rels: map[string]map[string]rel{}}
	for _, f := range zr.File {
		p.files[strings.ToLower(strings.TrimPrefix(f.Name, "/"))] = f
	}
	return p, nil
}

func (p *pkg) has(name string) bool {
	_, ok := p.files[strings.ToLower(strings.TrimPrefix(name, "/"))]
	return ok
}

// read returns the bytes of a part.
func (p *pkg) read(name string) ([]byte, error) {
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

// xml returns the parsed XML of a part (cached).
func (p *pkg) xml(name string) (*node, error) {
	key := strings.ToLower(strings.TrimPrefix(name, "/"))
	if n, ok := p.xmls[key]; ok {
		return n, nil
	}
	b, err := p.read(name)
	if err != nil {
		return nil, err
	}
	n, err := parseXML(b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	p.xmls[key] = n
	return n, nil
}

// relsOf returns the relationships of a part by ID.
func (p *pkg) relsOf(part string) map[string]rel {
	part = strings.TrimPrefix(part, "/")
	if m, ok := p.rels[part]; ok {
		return m
	}
	dir, base := path.Split(part)
	m := map[string]rel{}
	p.rels[part] = m
	root, err := p.xml(dir + "_rels/" + base + ".rels")
	if err != nil {
		return m
	}
	for _, r := range root.children("Relationship") {
		rl := rel{ID: r.attrStr("Id", ""), Type: r.attrStr("Type", ""), Target: r.attrStr("Target", "")}
		if strings.EqualFold(r.attrStr("TargetMode", ""), "External") {
			rl.External = true
		} else {
			rl.Target = resolvePart(dir, rl.Target)
		}
		m[rl.ID] = rl
	}
	return m
}

// resolvePart resolves a relationship target against the source part's directory.
func resolvePart(dir, target string) string {
	if strings.HasPrefix(target, "/") {
		return path.Clean(target)[1:]
	}
	return strings.TrimPrefix(path.Clean("/"+dir+target), "/")
}

// target returns the part a relationship ID points to.
func (p *pkg) target(part, id string) (rel, bool) {
	r, ok := p.relsOf(part)[id]
	return r, ok && r.Target != ""
}

// relOfType returns the first relationship whose type ends with suffix
// (e.g. "/slideLayout"), in ID order for determinism.
func (p *pkg) relOfType(part, suffix string) (rel, bool) {
	var best rel
	found := false
	for _, r := range p.relsOf(part) {
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
