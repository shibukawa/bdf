package drawio

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"strings"
	"sync"
)

// The stencil registry (mxStencilRegistry with draw.io's dynamic loading
// in Graph.js): a shape name such as mxgraph.flowchart.decision names a
// library file (stencils/flowchart.xml) from its middle parts, and the
// file registers its shapes under the lowercased name of its <shapes>
// element and the lowercased shape names with "_" for spaces.
//
// The libraries are embedded gzip-compressed (stencils/*.xml.gz, made by
// tools/gen-drawio-stencils from draw.io's stencil files; see
// stencils/NOTICE for their license). A file is decompressed and parsed
// the first time one of its names is looked up, and kept for the life of
// the process. Names of other libraries, and the many mxgraph.* shapes
// that draw.io implements in JavaScript, resolve to nil.

//go:generate go run ../../tools/gen-drawio-stencils -out stencils

//go:embed stencils/*.xml.gz stencils/*/*.xml.gz
var stencilData embed.FS

// stencilLibraries are the stencil files of the libraries draw.io loads
// by a name other than the file's (mxStencilRegistry.libraries in
// Editor.js, without the JavaScript files). A library listed here loads
// only these files.
var stencilLibraries = map[string][]string{
	"atlassian":                     {"atlassian.xml"},
	"bpmn":                          {"bpmn.xml"},
	"bpmn2":                         {"bpmn.xml"},
	"cisco19":                       {"cisco19.xml"},
	"cisco_safe":                    {"cisco_safe/architecture.xml", "cisco_safe/business_icons.xml", "cisco_safe/capability.xml", "cisco_safe/design.xml", "cisco_safe/iot_things_icons.xml", "cisco_safe/people_places_things_icons.xml", "cisco_safe/security_icons.xml", "cisco_safe/technology_icons.xml", "cisco_safe/threat.xml"},
	"kubernetes":                    {"kubernetes.xml", "kubernetes2.xml"},
	"flowchart":                     {"flowchart.xml"},
	"rackGeneral":                   {"rack/general.xml"},
	"rackF5":                        {"rack/f5.xml"},
	"lean_mapping":                  {"lean_mapping.xml"},
	"basic":                         {"basic.xml"},
	"ios7icons":                     {"ios7/icons.xml"},
	"ios7ui":                        {"ios7/misc.xml"},
	"android":                       {"android/android.xml"},
	"electrical/abstract":           {"electrical/abstract.xml"},
	"electrical/logic_gates":        {"electrical/logic_gates.xml"},
	"electrical/miscellaneous":      {"electrical/miscellaneous.xml"},
	"electrical/signal_sources":     {"electrical/signal_sources.xml"},
	"electrical/electro-mechanical": {"electrical/electro-mechanical.xml"},
	"electrical/transmission":       {"electrical/transmission.xml"},
	"mockup/graphics":               {"mockup/misc.xml"},
	"mockup/misc":                   {"mockup/misc.xml"},
	"mockup/navigation":             {"mockup/misc.xml"},
	"floorplan":                     {"floorplan.xml"},
	"bootstrap":                     {"bootstrap.xml"},
	"gmdl":                          {"gmdl.xml"},
	"gcp2":                          {"gcp2.xml"},
	"ibm":                           {"ibm.xml"},
	"ibmcloud":                      {"ibm_cloud.xml"},
	"cabinets":                      {"cabinets.xml"},
	"eip":                           {"eip.xml"},
	"networks":                      {"networks.xml"},
	"networks2":                     {"networks2.xml"},
	"atlassian2":                    {"atlassian2.xml"},
	"aws3d":                         {"aws3d.xml"},
	"aws4":                          {"aws4.xml"},
	"aws4b":                         {"aws4.xml"},
	"veeam":                         {"veeam/2d.xml", "veeam/3d.xml", "veeam/veeam.xml"},
	"veeam2":                        {"veeam/2d.xml", "veeam/3d.xml", "veeam/veeam2.xml"},
	"pid2misc":                      {"pid/misc.xml"},
	"pidFlowSensors":                {"pid/flow_sensors.xml"},
	"salesforce":                    {"salesforce.xml"},
	"sap":                           {"sap.xml"},
	// libraries of JavaScript shapes only
	"mockup": nil, "arrows2": nil, "c4": nil, "dfd": nil, "er": nil, "ios": nil,
	"infographic": nil, "mockup/buttons": nil, "mockup/containers": nil,
	"mockup/forms": nil, "mockup/markup": nil, "mockup/text": nil,
	"archimate": nil, "archimate3": nil, "archimate4": nil, "sysml": nil,
	"uml25": nil, "pid2inst": nil, "pid2valves": nil, "emoji": nil,
}

// stencilFile is an embedded library file, parsed on first use.
type stencilFile struct {
	once   sync.Once
	shapes map[string]*stencil
}

var (
	stencilFilesMu sync.Mutex
	stencilFiles   = map[string]*stencilFile{}
)

// lookupStencil returns the stencil of a library shape name, or nil
// (mxStencilRegistry.getStencil as draw.io overrides it).
func lookupStencil(name string) *stencil {
	base := stencilBasename(name)
	if base == "" {
		return nil
	}
	files, ok := stencilLibraries[base]
	if !ok {
		files = []string{strings.Replace(base, "_-_", "_", 1) + ".xml"}
	}
	for _, f := range files {
		if st := loadStencilFile(f)[name]; st != nil {
			return st
		}
	}
	return nil
}

// stencilBasename returns the library a stencil name belongs to: the
// parts between "mxgraph" and the shape name, joined with "/"
// (mxStencilRegistry.getBasenameForStencil); "" for other names.
func stencilBasename(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) < 2 || parts[0] != "mxgraph" {
		return ""
	}
	base := parts[1]
	for i := 2; i < len(parts)-1; i++ {
		base += "/" + parts[i]
	}
	return base
}

// loadStencilFile returns the shapes of an embedded library file by name
// (nil when the file is not embedded).
func loadStencilFile(file string) map[string]*stencil {
	path := "stencils/" + file + ".gz"
	stencilFilesMu.Lock()
	f := stencilFiles[file]
	if f == nil {
		// only embedded files get an entry (names come from the diagram)
		if _, err := fs.Stat(stencilData, path); err != nil {
			stencilFilesMu.Unlock()
			return nil
		}
		f = &stencilFile{}
		stencilFiles[file] = f
	}
	stencilFilesMu.Unlock()
	f.once.Do(func() {
		data, err := fs.ReadFile(stencilData, path)
		if err != nil {
			return
		}
		zr, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return
		}
		xmlData, err := io.ReadAll(zr)
		if err != nil {
			return
		}
		root, err := parseStencilXML(xmlData)
		if err != nil {
			return
		}
		f.shapes = map[string]*stencil{}
		parseStencilSet(root, f.shapes)
	})
	return f.shapes
}

// parseStencilSet adds the shapes of a <shapes> element (or of the
// <shapes> in a <stencils> element) to m under their registry names
// (mxStencilRegistry.parseStencilSet).
func parseStencilSet(root *stencilNode, m map[string]*stencil) {
	if root.name == "stencils" {
		for _, k := range root.kids {
			if k.name == "shapes" {
				parseStencilSet(k, m)
			}
		}
		return
	}
	pkg := ""
	if v, ok := root.attr("name"); ok {
		pkg = strings.ToLower(v + ".")
	}
	for _, k := range root.kids {
		name, ok := k.attr("name")
		if !ok {
			continue
		}
		m[pkg+strings.ToLower(strings.ReplaceAll(name, " ", "_"))] = newStencil(k)
	}
}

// stencil returns the stencil for a shape name, or nil: a library shape
// or a shape=stencil(...) style.
func (c *converter) stencil(name string) *stencil {
	if strings.HasPrefix(name, "stencil(") {
		return c.inlineStencil(name)
	}
	return lookupStencil(name)
}

// Stencils defined in the style (shape=stencil(...)) are cached by their
// definition; the cache is dropped when it grows beyond maxInlineStencils.
const maxInlineStencils = 256

type inlineResult struct {
	st  *stencil
	err error
}

var (
	inlineMu       sync.Mutex
	inlineStencils = map[string]inlineResult{}
)

// inlineStencil returns the stencil of a shape=stencil(...) style: the
// shape XML, URI-encoded, raw-deflated and base64-encoded (draw.io's
// mxCellRenderer.createShape override in Graph.js, with Graph.decompress).
func (c *converter) inlineStencil(name string) *stencil {
	inlineMu.Lock()
	r, ok := inlineStencils[name]
	inlineMu.Unlock()
	if !ok {
		st, err := decodeInlineStencil(name)
		r = inlineResult{st, err}
		inlineMu.Lock()
		if len(inlineStencils) >= maxInlineStencils {
			clear(inlineStencils)
		}
		inlineStencils[name] = r
		inlineMu.Unlock()
	}
	if r.err != nil {
		c.warnOnce("stencil("+truncate(name, 40), "shape=stencil(...) cannot be read: %v", r.err)
	}
	return r.st
}

// decodeInlineStencil parses the value of a shape=stencil(...) style.
func decodeInlineStencil(name string) (*stencil, error) {
	data := strings.TrimSuffix(strings.TrimPrefix(name, "stencil("), ")")
	xmlText, err := decompressGraph(data)
	if err != nil {
		return nil, err
	}
	root, err := parseStencilXML([]byte(xmlText))
	if err != nil {
		return nil, err
	}
	return newStencil(root), nil
}

// decompressGraph reverses Graph.compress: base64, raw deflate and URI
// encoding (Graph.decompress, with zapGremlins).
func decompressGraph(data string) (string, error) {
	// atob's forgiving base64: whitespace and missing padding are fine
	data = strings.Map(func(r rune) rune {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f' {
			return -1
		}
		return r
	}, data)
	raw, err := base64.RawStdEncoding.DecodeString(strings.TrimRight(data, "="))
	if err != nil {
		return "", fmt.Errorf("base64: %w", err)
	}
	inflated, err := io.ReadAll(flate.NewReader(bytes.NewReader(raw)))
	if err != nil {
		return "", fmt.Errorf("inflate: %w", err)
	}
	text, err := url.PathUnescape(string(inflated))
	if err != nil {
		return "", err
	}
	return zapGremlins(text), nil
}

// zapGremlins removes the characters XML does not allow (mxUtils.zapGremlins).
func zapGremlins(s string) string {
	return strings.Map(func(r rune) rune {
		if (r < 32 && r != '\t' && r != '\n' && r != '\r') || r == 0xFFFE || r == 0xFFFF {
			return -1
		}
		return r
	}, s)
}
