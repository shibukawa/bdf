// Command gen-drawio-stencils writes the draw.io stencil libraries that the
// drawio converter embeds (converter/drawio/stencils): it reads the stencil
// XML files of a draw.io release, drops what drawing does not use (the
// <connections> elements, comments and the whitespace between elements),
// rounds coordinates to a precision the stencil's size cannot show, writes
// the steps of each <path> as a compact d attribute ("M44 11L44 9C…", which
// the converter expands back into elements when it parses the file), and
// stores each file gzip-compressed, together with the license.
//
// Usage (from the repository root):
//
//	go run ./tools/gen-drawio-stencils [-repo URL-or-dir] [-out dir]
//
// -repo is the root of the draw.io repository: by default the release tag
// on GitHub the libraries were taken from, or a local checkout.
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// release is the draw.io version the embedded libraries come from.
const release = "v31.4.4"

// libraries are the stencil files embedded, relative to
// src/main/webapp/stencils. They are the general-purpose libraries
// (flowchart, basic shapes, arrows), the AWS icons (aws4: the current
// generation, which the converter also draws older AWS icons with), the
// network, rack, Cisco and legacy Azure icons of everyday IT diagrams,
// BPMN, integration patterns, lean mapping, floor plans, electrical
// circuits and P&ID. Other big vendor sets (aws3, cisco19, gcp2, office,
// ...) are left out to keep the package small.
var libraries = []string{
	"aws4.xml",
	"flowchart.xml",
	"basic.xml",
	"arrows.xml",
	"networks.xml",
	"bpmn.xml",
	"eip.xml",
	"lean_mapping.xml",
	"floorplan.xml",
	"azure.xml",
	"rack/general.xml",
	"cisco/buildings.xml",
	"cisco/computers_and_peripherals.xml",
	"cisco/controllers_and_modules.xml",
	"cisco/directors.xml",
	"cisco/hubs_and_gateways.xml",
	"cisco/misc.xml",
	"cisco/modems_and_phones.xml",
	"cisco/people.xml",
	"cisco/routers.xml",
	"cisco/security.xml",
	"cisco/servers.xml",
	"cisco/storage.xml",
	"cisco/switches.xml",
	"cisco/wireless.xml",
	"electrical/abstract.xml",
	"electrical/capacitors.xml",
	"electrical/diodes.xml",
	"electrical/electro-mechanical.xml",
	"electrical/iec417.xml",
	"electrical/iec_logic_gates.xml",
	"electrical/inductors.xml",
	"electrical/instruments.xml",
	"electrical/logic_gates.xml",
	"electrical/miscellaneous.xml",
	"electrical/mosfets1.xml",
	"electrical/mosfets2.xml",
	"electrical/op_amps.xml",
	"electrical/opto_electronics.xml",
	"electrical/plc_ladder.xml",
	"electrical/power_semiconductors.xml",
	"electrical/radio.xml",
	"electrical/resistors.xml",
	"electrical/rot_mech.xml",
	"electrical/signal_sources.xml",
	"electrical/thermionic_devices.xml",
	"electrical/transistors.xml",
	"electrical/transmission.xml",
	"electrical/waveforms.xml",
	"pid/agitators.xml",
	"pid/apparatus_elements.xml",
	"pid/centrifuges.xml",
	"pid/compressors.xml",
	"pid/compressors_iso.xml",
	"pid/crushers_grinding.xml",
	"pid/driers.xml",
	"pid/engines.xml",
	"pid/feeders.xml",
	"pid/filters.xml",
	"pid/fittings.xml",
	"pid/flow_sensors.xml",
	"pid/heat_exchangers.xml",
	"pid/instruments.xml",
	"pid/misc.xml",
	"pid/mixers.xml",
	"pid/piping.xml",
	"pid/pumps.xml",
	"pid/pumps_din.xml",
	"pid/pumps_iso.xml",
	"pid/separators.xml",
	"pid/shaping_machines.xml",
	"pid/valves.xml",
	"pid/vessels.xml",
}

func main() {
	repo := flag.String("repo", "https://raw.githubusercontent.com/jgraph/drawio/"+release, "draw.io repository root (URL or directory)")
	out := flag.String("out", "converter/drawio/stencils", "output directory")
	flag.Parse()

	total := 0
	for _, lib := range libraries {
		src, err := read(*repo, "src/main/webapp/stencils/"+lib)
		if err != nil {
			fail(err)
		}
		min, err := minify(src)
		if err != nil {
			fail(fmt.Errorf("%s: %w", lib, err))
		}
		var buf bytes.Buffer
		zw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
		zw.Write(min)
		zw.Close()
		dst := filepath.Join(*out, filepath.FromSlash(lib)+".gz")
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			fail(err)
		}
		if err := os.WriteFile(dst, buf.Bytes(), 0o644); err != nil {
			fail(err)
		}
		total += buf.Len()
		fmt.Printf("%8d %8d %s\n", len(src), buf.Len(), lib)
	}
	fmt.Printf("%8s %8d total\n", "", total)

	license, err := read(*repo, "LICENSE")
	if err != nil {
		fail(err)
	}
	notice, err := read(*repo, "src/main/webapp/stencils/LICENSE")
	if err != nil {
		fail(err)
	}
	if err := os.WriteFile(filepath.Join(*out, "LICENSE"), license, 0o644); err != nil {
		fail(err)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "The *.xml.gz files in this directory are draw.io's stencil libraries\n")
	fmt.Fprintf(&b, "(https://github.com/jgraph/drawio, %s, src/main/webapp/stencils),\n", release)
	fmt.Fprintf(&b, "Copyright (c) JGraph Ltd, licensed under the Apache License 2.0 (see\n")
	fmt.Fprintf(&b, "LICENSE). They were modified by tools/gen-drawio-stencils: the\n")
	fmt.Fprintf(&b, "<connections> elements, comments and the whitespace between elements\n")
	fmt.Fprintf(&b, "were removed, coordinates were rounded, the steps of paths were written\n")
	fmt.Fprintf(&b, "in a compact form, and each file was compressed with gzip.\n\n")
	fmt.Fprintf(&b, "draw.io's stencils/LICENSE reads:\n\n")
	b.Write(bytes.TrimSpace(notice))
	b.WriteString("\n")
	if err := os.WriteFile(filepath.Join(*out, "NOTICE"), []byte(b.String()), 0o644); err != nil {
		fail(err)
	}
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

// minify rewrites a stencil file without <connections>, comments,
// processing instructions and whitespace-only text; empty elements are
// written self-closing.
func minify(src []byte) ([]byte, error) {
	d := xml.NewDecoder(bytes.NewReader(src))
	var toks []xml.Token
	skip := 0
	prec := 2
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			if skip > 0 || t.Name.Local == "connections" {
				skip++
				continue
			}
			t = t.Copy()
			if t.Name.Local == "shape" {
				prec = precision(t)
			}
			roundAttrs(t, prec)
			toks = append(toks, t)
		case xml.EndElement:
			if skip > 0 {
				skip--
				continue
			}
			toks = append(toks, t)
		case xml.CharData:
			if skip == 0 && len(bytes.TrimSpace(t)) > 0 {
				toks = append(toks, t.Copy())
			}
		}
	}
	toks = compactPaths(toks)
	var b bytes.Buffer
	for i := 0; i < len(toks); i++ {
		switch t := toks[i].(type) {
		case xml.StartElement:
			b.WriteString("<" + t.Name.Local)
			for _, a := range t.Attr {
				b.WriteString(" " + a.Name.Local + `="`)
				xml.EscapeText(&b, []byte(a.Value))
				b.WriteString(`"`)
			}
			if i+1 < len(toks) {
				if _, ok := toks[i+1].(xml.EndElement); ok {
					b.WriteString("/>")
					i++
					continue
				}
			}
			b.WriteString(">")
		case xml.EndElement:
			b.WriteString("</" + t.Name.Local + ">")
		case xml.CharData:
			xml.EscapeText(&b, t)
		}
	}
	b.WriteString("\n")
	return b.Bytes(), nil
}

// precision returns the decimals coordinates of a <shape> are rounded to:
// two, or more for small stencils, so that the rounding stays below a
// two-thousandth of the stencil's size.
func precision(shape xml.StartElement) int {
	size := 0.0
	for _, a := range shape.Attr {
		if a.Name.Local == "w" || a.Name.Local == "h" {
			if v, err := strconv.ParseFloat(a.Value, 64); err == nil {
				size = math.Max(size, v)
			}
		}
	}
	if size <= 0 {
		return 4
	}
	return max(2, int(math.Ceil(math.Log10(2000/size))))
}

// geometryAttrs are the coordinate attributes of the drawing elements.
var geometryAttrs = map[string]map[string]bool{
	"move": {"x": true, "y": true}, "line": {"x": true, "y": true},
	"quad":  {"x1": true, "y1": true, "x2": true, "y2": true},
	"curve": {"x1": true, "y1": true, "x2": true, "y2": true, "x3": true, "y3": true},
	"arc":   {"rx": true, "ry": true, "x": true, "y": true},
	"rect":  {"x": true, "y": true, "w": true, "h": true}, "roundrect": {"x": true, "y": true, "w": true, "h": true},
	"ellipse": {"x": true, "y": true, "w": true, "h": true},
}

// roundAttrs rounds the coordinates of a drawing element.
func roundAttrs(t xml.StartElement, prec int) {
	keys := geometryAttrs[t.Name.Local]
	for i, a := range t.Attr {
		if keys[a.Name.Local] {
			if r, ok := roundNum(a.Value, prec); ok {
				t.Attr[i].Value = r
			}
		}
	}
}

// roundNum rounds a plain decimal number; ok is false for anything else,
// which is kept as it is.
func roundNum(s string, prec int) (string, bool) {
	s = strings.TrimSpace(s)
	if s == "" || strings.ContainsAny(s, "eExX") {
		return "", false
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return "", false
	}
	r := strconv.FormatFloat(v, 'f', prec, 64)
	if strings.Contains(r, ".") {
		r = strings.TrimRight(strings.TrimRight(r, "0"), ".")
	}
	if r == "-0" || r == "" {
		r = "0"
	}
	return r, true
}

// pathSteps are the steps a compact path holds, with their letters and
// the attributes of their numbers in order.
var pathSteps = map[string]struct {
	letter byte
	attrs  []string
}{
	"move":  {'M', []string{"x", "y"}},
	"line":  {'L', []string{"x", "y"}},
	"quad":  {'Q', []string{"x1", "y1", "x2", "y2"}},
	"curve": {'C', []string{"x1", "y1", "x2", "y2", "x3", "y3"}},
	"arc":   {'A', []string{"rx", "ry", "x-axis-rotation", "large-arc-flag", "sweep-flag", "x", "y"}},
	"close": {'Z', nil},
}

// compactPaths writes the steps of each <path> as its d attribute: a
// letter per step followed by its numbers, separated by spaces. Paths with
// other elements, attributes a step does not have, or values that are not
// plain numbers are kept as elements.
func compactPaths(toks []xml.Token) []xml.Token {
	var out []xml.Token
	for i := 0; i < len(toks); i++ {
		start, ok := toks[i].(xml.StartElement)
		if !ok || start.Name.Local != "path" {
			out = append(out, toks[i])
			continue
		}
		var d strings.Builder
		j := i + 1
		for ; j < len(toks); j++ {
			if _, end := toks[j].(xml.EndElement); end {
				break
			}
			step, ok := toks[j].(xml.StartElement)
			spec, known := pathSteps[step.Name.Local]
			if !ok || !known || j+1 >= len(toks) {
				d.Reset()
				break
			}
			if _, end := toks[j+1].(xml.EndElement); !end {
				d.Reset()
				break
			}
			vals := map[string]string{}
			for _, a := range step.Attr {
				vals[a.Name.Local] = a.Value
			}
			if len(vals) > len(spec.attrs) {
				d.Reset()
				break
			}
			d.WriteByte(spec.letter)
			valid := true
			for k, name := range spec.attrs {
				v, ok := vals[name]
				if !ok {
					v = "0" // Number("") is 0, as for a missing attribute
				} else if _, err := strconv.ParseFloat(v, 64); err != nil || strings.ContainsAny(v, "eExX ") {
					valid = false
					break
				}
				if k > 0 {
					d.WriteByte(' ')
				}
				d.WriteString(v)
			}
			if !valid {
				d.Reset()
				break
			}
			j++ // the step's end element
		}
		if d.Len() == 0 || j >= len(toks) {
			out = append(out, toks[i])
			continue
		}
		start.Attr = append(append([]xml.Attr(nil), start.Attr...), xml.Attr{Name: xml.Name{Local: "d"}, Value: d.String()})
		out = append(out, start, toks[j])
		i = j
	}
	return out
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "gen-drawio-stencils:", err)
	os.Exit(1)
}
