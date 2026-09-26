// Command gen-drawio-awsmap writes the table the drawio converter draws the
// shapes of older AWS icon generations with (converter/drawio/
// aws_legacy_table.go): for every shape name of mxgraph.aws (stencils/aws),
// mxgraph.aws2 (stencils/aws2), mxgraph.aws3 (stencils/aws3.xml) and
// mxgraph.aws3d (stencils/aws3d.xml and mxAWS3D.js), the style that
// draw.io's current AWS palette (Sidebar-AWS4.js) draws the same service
// or resource with, or none.
//
// The names are matched with the entries of the current palette by the
// titles and tags of the old palettes (Sidebar-AWS.js, Sidebar-AWS3.js,
// Sidebar-AWS3D.js) and of the current one (see match.go); overrides.txt
// holds the curated mappings for what the matching gets wrong or misses,
// and the names left without a counterpart. Every name must be decided:
// the generator fails on names that are neither matched nor listed.
//
// Usage (from the repository root):
//
//	go run ./tools/gen-drawio-awsmap [-repo URL-or-dir] [-out file] [-v]
//
// -repo is the root of the draw.io repository: by default the release tag
// on GitHub, or a local checkout (or the files extracted from the desktop
// app's app.asar under drawio/src/main/webapp). -v prints the decision for
// every name; -sheets writes diagrams of the old icons next to their
// targets for looking at the mapping in draw.io, and -targets a diagram of
// palette entries for choosing overrides. -samples writes the converter's
// test diagrams of the old icons:
//
//	go run ./tools/gen-drawio-awsmap -samples converter/drawio/testdata/aws_legacy
package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// release is the draw.io version the palettes and stencils are read from.
const release = "v31.4.4"

func main() {
	repo := flag.String("repo", "https://raw.githubusercontent.com/jgraph/drawio/"+release, "draw.io repository root (URL or directory)")
	out := flag.String("out", "converter/drawio/aws_legacy_table.go", "the Go file to write")
	overridesPath := flag.String("overrides", "tools/gen-drawio-awsmap/overrides.txt", "the curated mappings")
	verbose := flag.Bool("v", false, "print the decision for every old name")
	sheets := flag.String("sheets", "", "write diagrams of the old icons and of the mapping into this directory")
	samples := flag.String("samples", "", "write the converter's sample diagrams of the old icons into this directory (converter/drawio/testdata/aws_legacy)")
	targetSheet := flag.String("targets", "", "write a diagram of the current palette's entries (those named by -only) to this file")
	only := flag.String("only", "", "comma-separated target labels (as in overrides.txt) for -targets")
	flag.Parse()

	src, err := load(*repo)
	if err != nil {
		fail(err)
	}
	ovText, err := os.ReadFile(*overridesPath)
	if err != nil {
		fail(err)
	}
	ovs, err := parseOverrides(string(ovText))
	if err != nil {
		fail(err)
	}
	targets := src.targets()
	if *targetSheet != "" {
		var labels []string
		if *only != "" {
			labels = strings.Split(*only, ",")
		}
		if err := writeTargetSheet(targets, labels, *targetSheet); err != nil {
			fail(err)
		}
	}
	ds, unclassified, err := classify(src, ovs, targets)
	if err != nil {
		fail(err)
	}
	report(ds, unclassified, *verbose)
	if *sheets != "" {
		if err := writeSheets(src, *sheets); err != nil {
			fail(err)
		}
		if err := writeCompareSheets(src, ds, *sheets); err != nil {
			fail(err)
		}
	}
	if len(unclassified) > 0 {
		fail(fmt.Errorf("%d names are neither matched nor in overrides.txt", len(unclassified)))
	}
	if *samples != "" {
		if err := writeSamples(src, ds, *samples); err != nil {
			fail(err)
		}
	}
	code, err := writeTable(ds)
	if err != nil {
		fail(err)
	}
	if err := os.WriteFile(*out, code, 0o644); err != nil {
		fail(err)
	}
}

// oldName is a shape name of an older AWS generation.
type oldName struct {
	gen     string // aws, aws2, aws3, aws3d
	name    string // the registry name, e.g. mxgraph.aws3.ec2
	display string // the name in the stencil file or JavaScript class
	w, h    float64
}

// sources are the inputs read from the draw.io repository.
type sources struct {
	aws4  []paletteEntry
	old   []paletteEntry
	names []oldName
	// aws4Stencils are the names of the aws4 stencils.
	aws4Stencils map[string]stencilInfo
}

type stencilInfo struct {
	w, h float64
}

// oldStencilFiles are the stencil files of the older generations,
// relative to src/main/webapp/stencils.
var oldStencilFiles = []string{
	"aws/compute.xml", "aws/content_delivery.xml", "aws/database.xml", "aws/deployment_management.xml",
	"aws/groups.xml", "aws/messaging.xml", "aws/misc.xml", "aws/networking.xml",
	"aws/non_service_specific.xml", "aws/on_demand_workforce.xml", "aws/storage.xml",
	"aws2/administration_and_security.xml", "aws2/analytics.xml", "aws2/app_services.xml",
	"aws2/compute_and_networking.xml", "aws2/database.xml", "aws2/deployment_and_management.xml",
	"aws2/developer_tools.xml", "aws2/enterprise_applications.xml", "aws2/game_development.xml",
	"aws2/internet_of_things.xml", "aws2/management_tools.xml", "aws2/mobile_services.xml",
	"aws2/networking.xml", "aws2/non-service_specific.xml", "aws2/on-demand_workforce.xml",
	"aws2/sdks.xml", "aws2/security_and_identity.xml", "aws2/storage_and_content_delivery.xml",
	"aws3.xml", "aws3d.xml",
}

// oldSidebars are the sidebar scripts of the older palettes.
var oldSidebars = []string{"Sidebar-AWS.js", "Sidebar-AWS3.js", "Sidebar-AWS3D.js"}

func load(repo string) (*sources, error) {
	s := &sources{aws4Stencils: map[string]stencilInfo{}}
	for _, f := range oldStencilFiles {
		data, err := read(repo, "src/main/webapp/stencils/"+f)
		if err != nil {
			return nil, err
		}
		gen := strings.TrimSuffix(strings.SplitN(f, "/", 2)[0], ".xml")
		for _, sh := range stencilShapes(string(data)) {
			s.names = append(s.names, oldName{gen: gen, name: sh.name, display: sh.display, w: sh.w, h: sh.h})
		}
	}
	// aws3d shapes drawn by JavaScript (mxAWS3D.js)
	js, err := read(repo, "src/main/webapp/shapes/mxAWS3D.js")
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, n := range s.names {
		seen[n.name] = true
	}
	for _, m := range regexp.MustCompile(`'(mxgraph\.aws3d\.[A-Za-z0-9_]+)'`).FindAllStringSubmatch(string(js), -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			s.names = append(s.names, oldName{gen: "aws3d", name: m[1], display: strings.TrimPrefix(m[1], "mxgraph.aws3d.")})
		}
	}
	sort.SliceStable(s.names, func(i, j int) bool {
		if s.names[i].gen != s.names[j].gen {
			return s.names[i].gen < s.names[j].gen
		}
		return s.names[i].name < s.names[j].name
	})

	for _, f := range oldSidebars {
		es, err := readSidebar(repo, f)
		if err != nil {
			return nil, err
		}
		s.old = append(s.old, es...)
	}
	s.aws4, err = readSidebar(repo, "Sidebar-AWS4.js")
	if err != nil {
		return nil, err
	}
	data, err := read(repo, "src/main/webapp/stencils/aws4.xml")
	if err != nil {
		return nil, err
	}
	for _, sh := range stencilShapes(string(data)) {
		s.aws4Stencils[sh.name] = stencilInfo{sh.w, sh.h}
	}
	return s, nil
}

func readSidebar(repo, name string) ([]paletteEntry, error) {
	data, err := read(repo, "src/main/webapp/js/diagramly/sidebar/"+name)
	if err != nil {
		return nil, err
	}
	es, err := readPalettes(string(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return es, nil
}

type stencilShape struct {
	name, display string
	w, h          float64
}

var (
	shapesRe = regexp.MustCompile(`<shapes\s+name="([^"]*)"`)
	shapeRe  = regexp.MustCompile(`<shape\s[^>]*>`)
	attrRe   = regexp.MustCompile(`(\w+)="([^"]*)"`)
)

// stencilShapes returns the shapes of a stencil file with their registry
// names (mxStencilRegistry.parseStencilSet: the lowercased package and
// shape names, spaces in shape names as "_").
func stencilShapes(xml string) []stencilShape {
	pkg := ""
	if m := shapesRe.FindStringSubmatch(xml); m != nil {
		pkg = strings.ToLower(m[1]) + "."
	}
	var out []stencilShape
	for _, tag := range shapeRe.FindAllString(xml, -1) {
		attrs := map[string]string{}
		for _, a := range attrRe.FindAllStringSubmatch(tag, -1) {
			attrs[a[1]] = a[2]
		}
		name, ok := attrs["name"]
		if !ok {
			continue
		}
		sh := stencilShape{name: pkg + strings.ToLower(strings.ReplaceAll(unescapeXML(name), " ", "_")), display: unescapeXML(name)}
		fmt.Sscan(attrs["w"], &sh.w)
		fmt.Sscan(attrs["h"], &sh.h)
		out = append(out, sh)
	}
	return out
}

func unescapeXML(s string) string {
	return strings.NewReplacer("&amp;", "&", "&lt;", "<", "&gt;", ">", "&quot;", `"`, "&apos;", "'").Replace(s)
}

// read returns a file of the repository.
func read(repo, name string) ([]byte, error) {
	if !strings.HasPrefix(repo, "https://") && !strings.HasPrefix(repo, "http://") {
		return os.ReadFile(filepath.Join(repo, filepath.FromSlash(name)))
	}
	resp, err := http.Get(strings.TrimSuffix(repo, "/") + "/" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", name, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "gen-drawio-awsmap:", err)
	os.Exit(1)
}
