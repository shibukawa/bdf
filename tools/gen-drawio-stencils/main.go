// Command gen-drawio-stencils writes the draw.io stencil libraries that the
// drawio converter embeds (converter/drawio/stencils): it reads the stencil
// XML files of a draw.io release, drops what drawing does not use (the
// <connections> elements, comments and the whitespace between elements),
// and stores each file gzip-compressed, together with the license.
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
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// release is the draw.io version the embedded libraries come from.
const release = "v31.4.4"

// libraries are the stencil files embedded, relative to
// src/main/webapp/stencils. They are the general-purpose libraries
// (flowchart, basic shapes, arrows), the network, rack, Cisco and legacy
// Azure icons of everyday IT diagrams, BPMN, integration patterns, lean
// mapping, floor plans, electrical circuits and P&ID. Big vendor sets
// (aws3, aws4, cisco19, gcp2, office, ...) are left out to keep the
// package small; most of the newer ones draw their icons through
// JavaScript shapes (such as mxgraph.aws4.resourceIcon) anyway.
var libraries = []string{
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
	fmt.Fprintf(&b, "were removed, and each file was compressed with gzip.\n\n")
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
			toks = append(toks, t.Copy())
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

func fail(err error) {
	fmt.Fprintln(os.Stderr, "gen-drawio-stencils:", err)
	os.Exit(1)
}
