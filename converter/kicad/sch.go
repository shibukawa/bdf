package kicad

import (
	"fmt"
	"path"
	"slices"
	"strconv"
	"strings"
)

// schFile is a schematic file: the root sheet or the file of hierarchical
// sheets.
type schFile struct {
	path string // relative to the input's directory (slash-separated)
	n    *node
	libs map[string]*node // lib_symbols by name
	uuid string
	size int // bytes
}

// sheetInst is a sheet of the hierarchy as it appears in the schematic: a
// file used twice is two instances, each a page.
type sheetInst struct {
	file *schFile
	// path is the instance path of KiCad 7 and later: the UUIDs of the
	// root and of the sheets down to this one ("/root/sheet/..."); pathV6
	// is that of KiCad 6, without the root ("" for the root sheet)
	path, pathV6 string
	page         string // page number as shown
	name         string // sheet name ("" for the root)
	namePath     string // "/", "/name/", "/name/sub/"
	parent       *sheetInst
	node         *node // the (sheet) in the parent
	index        int   // 1-based page of the view
}

// maxSheetDepth bounds the nesting of sheets (a sheet cannot contain
// itself, but a malformed file could).
const maxSheetDepth = 32

// maxSheets bounds the sheet instances of a schematic.
const maxSheets = 2000

// loadSchematic reads the hierarchy of sheets under the root schematic.
func (c *conv) loadSchematic(rootPath string, data []byte) ([]*sheetInst, error) {
	files := map[string]*schFile{}
	open := func(p string, data []byte) (*schFile, error) {
		if f, ok := files[p]; ok {
			return f, nil
		}
		n, err := parse(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		if n.name != "kicad_sch" {
			return nil, fmt.Errorf("%s: not a KiCad schematic", p)
		}
		if v := n.child("version").int(0); v > 0 && v < 20211014 {
			return nil, fmt.Errorf("%s: a schematic of KiCad 5 or earlier (version %d); open it in KiCad 6 or later and save it", p, v)
		}
		f := &schFile{path: p, n: n, libs: map[string]*node{}, uuid: n.str("uuid"), size: len(data)}
		for _, s := range n.child("lib_symbols").children("symbol") {
			f.libs[s.arg(0)] = s
		}
		files[p] = f
		return f, nil
	}
	rootFile, err := open(rootPath, data)
	if err != nil {
		return nil, err
	}
	root := &sheetInst{file: rootFile, path: "/" + rootFile.uuid, namePath: "/", page: "1"}
	if rootFile.uuid == "" {
		root.path = "/"
	}
	// KiCad 6 keeps the page numbers of every sheet in the root
	v6pages := map[string]string{}
	for _, p := range rootFile.n.child("sheet_instances").children("path") {
		v6pages[p.arg(0)] = p.str("page")
	}
	if pg := v6pages["/"]; pg != "" {
		root.page = pg
	}
	out := []*sheetInst{root}
	var walk func(s *sheetInst, depth int)
	walk = func(s *sheetInst, depth int) {
		if depth >= maxSheetDepth {
			c.warnOnce("sheet-depth", "sheets nested deeper than %d levels are left out", maxSheetDepth)
			return
		}
		for _, sh := range s.file.n.children("sheet") {
			if len(out) >= maxSheets {
				c.warnOnce("sheet-count", "more than %d sheets: the rest are left out", maxSheets)
				return
			}
			rel := sheetFile(sh)
			uuid := sh.str("uuid")
			child := &sheetInst{
				path:   s.path + "/" + uuid,
				pathV6: s.pathV6 + "/" + uuid,
				name:   sheetName(sh),
				parent: s,
				node:   sh,
			}
			if s.path == "/" {
				child.path = "/" + uuid
			}
			child.namePath = s.namePath + child.name + "/"
			child.page = sheetPage(sh, s.path, c.project)
			if child.page == "" {
				child.page = v6pages[child.pathV6]
			}
			if rel == "" {
				c.warnf("sheet %q names no file", child.name)
				continue
			}
			if c.refs == nil {
				c.warnOnce("sheet-files", "the sheets' files were not given with the schematic (%s): their pages are left out", rel)
				continue
			}
			b, p, err := c.readRef(s.file.path, rel)
			if err != nil {
				c.warnf("sheet %q: %s is not among the files given: its page is left out", child.name, rel)
				continue
			}
			if s.uses(p) {
				// KiCad refuses such a hierarchy
				c.warnf("sheet %q: %s contains itself: its page is left out", child.name, rel)
				continue
			}
			f, err := open(p, b)
			if err != nil {
				c.warnf("sheet %q: %v", child.name, err)
				continue
			}
			child.file = f
			out = append(out, child)
			walk(child, depth+1)
		}
	}
	walk(root, 0)
	repairPages(out)
	// pages in the order of their numbers, as KiCad lists and plots them
	order := make([]*sheetInst, len(out))
	copy(order, out)
	slices.SortStableFunc(order, func(a, b *sheetInst) int { return comparePages(a.page, b.page) })
	for i, s := range order {
		s.index = i + 1
	}
	return order, nil
}

// repairPages numbers the sheets without a page number, and those whose
// number an earlier sheet (in hierarchy order) has, with the lowest numbers
// no sheet has (SCH_SHEET_LIST::RepairPageNumbers).
func repairPages(sheets []*sheetInst) {
	reserved := map[string]bool{}
	for _, s := range sheets {
		if s.page != "" {
			reserved[s.page] = true
		}
	}
	assigned := map[string]bool{}
	next := 1
	for _, s := range sheets {
		if s.page != "" && !assigned[s.page] {
			assigned[s.page] = true
			continue
		}
		for reserved[strconv.Itoa(next)] || assigned[strconv.Itoa(next)] {
			next++
		}
		s.page = strconv.Itoa(next)
		assigned[s.page] = true
		next++
	}
}

// comparePages orders page numbers: numerically when both are numbers,
// else as strings; pages without a number last.
func comparePages(a, b string) int {
	switch {
	case a == "" && b == "":
		return 0
	case a == "":
		return 1
	case b == "":
		return -1
	}
	na, ea := strconv.Atoi(a)
	nb, eb := strconv.Atoi(b)
	if ea == nil && eb == nil {
		return na - nb
	}
	return strings.Compare(a, b)
}

// sheetProperty returns the value of a property of a sheet by its name in
// KiCad 7 and later, or in KiCad 6.
func sheetProperty(sh *node, name, v6 string) *node {
	for _, p := range sh.children("property") {
		switch p.arg(0) {
		case name, v6:
			return p
		}
	}
	return nil
}

func sheetFile(sh *node) string {
	return strings.ReplaceAll(sheetProperty(sh, "Sheetfile", "Sheet file").arg(1), `\`, "/")
}

func sheetName(sh *node) string { return sheetProperty(sh, "Sheetname", "Sheet name").arg(1) }

// sheetPage returns a sheet's page number from its instances (KiCad 7 and
// later): the one of the project, else the first.
func sheetPage(sh *node, parentPath, project string) string {
	inst := sh.child("instances")
	var first string
	for _, pr := range inst.children("project") {
		for _, p := range pr.children("path") {
			if p.arg(0) != parentPath {
				continue
			}
			pg := p.str("page")
			if pr.arg(0) == project {
				return pg
			}
			if first == "" {
				first = pg
			}
		}
	}
	return first
}

// symbolInstance returns the reference and unit of a symbol on a sheet:
// from the symbol's instances (KiCad 7 and later), from the root's list
// (KiCad 6), or from the symbol itself.
func (c *conv) symbolInstance(sym *node, s *sheetInst, root *schFile, v6refs map[string]*node) (ref string, unit int) {
	unit = sym.child("unit").int(0)
	ref = symbolProperty(sym, "Reference").arg(1)
	var first *node
	for _, pr := range sym.child("instances").children("project") {
		for _, p := range pr.children("path") {
			if p.arg(0) != s.path {
				continue
			}
			if pr.arg(0) == c.project {
				first = p
				break
			}
			if first == nil {
				first = p
			}
		}
	}
	if first == nil {
		if p := v6refs[s.pathV6+"/"+sym.str("uuid")]; p != nil {
			first = p
		}
	}
	if first != nil {
		if r := first.str("reference"); r != "" {
			ref = r
		}
		if u := first.child("unit"); u != nil {
			unit = u.int(0)
		}
	}
	return ref, unit
}

// symbolProperty returns a property of a symbol by name.
func symbolProperty(sym *node, name string) *node {
	for _, p := range sym.children("property") {
		if p.arg(0) == name {
			return p
		}
	}
	return nil
}

// displayPath is a sheet's path as KiCad shows it in the title block.
// uses reports whether the sheet or one above it is drawn from the file.
func (s *sheetInst) uses(file string) bool {
	for ; s != nil; s = s.parent {
		if s.file.path == file {
			return true
		}
	}
	return false
}

func (s *sheetInst) displayPath() string { return s.namePath }

// fileBase is the file name of a sheet's file.
func (s *sheetInst) fileBase() string { return path.Base(s.file.path) }
