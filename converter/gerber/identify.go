package gerber

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"unicode"
)

// kind is what a layer of a board is.
type kind uint8

const (
	kOther kind = iota
	kCopper
	kMask
	kSilk
	kPaste
	kOutline
	kDrill
)

// side is the side of the board a layer is on.
type side uint8

const (
	sNone side = iota
	sTop
	sInner
	sBottom
)

// function is what a layer is, and where.
type function struct {
	kind kind
	side side
	// index numbers inner copper layers (1: the first under the top)
	index int
	// plated tells drill layers: "plated", "nonplated" or ""
	plated string
	// what the file attributes call a layer this converter does not draw
	// in the board views ("Drillmap", "AssemblyDrawing,Top")
	other string
}

// order is the place of a layer in the list of layers: top to bottom
// through the board, then the outline, the drills and the rest.
func (f function) order() int {
	switch f.kind {
	case kSilk, kPaste, kMask, kCopper:
		base := map[kind]int{kSilk: 0, kPaste: 1, kMask: 2, kCopper: 3}[f.kind]
		switch f.side {
		case sInner:
			return 1000 + f.index
		case sBottom:
			return 3000 - base
		}
		return base
	case kOutline:
		return 4000
	case kDrill:
		return 5000
	}
	return 6000
}

// title names the layer.
func (f function) title() string {
	sideName := map[side]string{sTop: "Top", sBottom: "Bottom", sNone: "Top"}[f.side]
	switch f.kind {
	case kCopper:
		if f.side == sInner {
			return fmt.Sprintf("Inner copper %d", f.index)
		}
		return sideName + " copper"
	case kMask:
		return sideName + " solder mask"
	case kSilk:
		return sideName + " silkscreen"
	case kPaste:
		return sideName + " paste"
	case kOutline:
		return "Outline"
	case kDrill:
		switch f.plated {
		case "plated":
			return "Drill (plated)"
		case "nonplated":
			return "Drill (non-plated)"
		}
		return "Drill"
	}
	return ""
}

// id is the view ID of the layer ("top-copper").
func (f function) id() string {
	if f.kind == kOther {
		return ""
	}
	return strings.NewReplacer(" ", "-", "(", "", ")", "").Replace(strings.ToLower(f.title()))
}

// fileFunction reads the value of a TF.FileFunction attribute
// ("Copper,L1,Top,Signal", "Soldermask,Bot", "Plated,1,2,PTH").
func fileFunction(v string) (function, bool) {
	fs := strings.Split(v, ",")
	for i := range fs {
		fs[i] = strings.TrimSpace(fs[i])
	}
	sideOf := func(i int) side {
		if i < len(fs) {
			switch strings.ToLower(fs[i]) {
			case "top":
				return sTop
			case "bot":
				return sBottom
			case "inr":
				return sInner
			}
		}
		return sNone
	}
	switch strings.ToLower(fs[0]) {
	case "copper":
		f := function{kind: kCopper, side: sideOf(2)}
		if len(fs) > 1 && strings.HasPrefix(strings.ToUpper(fs[1]), "L") {
			n, _ := strconv.Atoi(fs[1][1:])
			f.index = n - 1 // L2 is the first inner layer
			if f.side == sNone && n == 1 {
				f.side = sTop
			}
		}
		if f.side == sNone {
			f.side = sInner
		}
		return f, true
	case "soldermask":
		return function{kind: kMask, side: sideOf(1)}, true
	case "legend":
		return function{kind: kSilk, side: sideOf(1)}, true
	case "paste", "solderpaste":
		return function{kind: kPaste, side: sideOf(1)}, true
	case "profile":
		return function{kind: kOutline}, true
	case "plated":
		return function{kind: kDrill, plated: "plated"}, true
	case "nonplated":
		return function{kind: kDrill, plated: "nonplated"}, true
	case "mixedplating":
		return function{kind: kDrill}, true
	case "":
		return function{}, false
	}
	return function{kind: kOther, other: v}, true
}

// extensions are the file name extensions that tell a layer: those of
// Protel and Altium (which KiCad and EasyEDA can write too), of Eagle's
// older CAM jobs, and of OrCAD.
var extensions = map[string]function{
	"gtl": {kind: kCopper, side: sTop}, "gbl": {kind: kCopper, side: sBottom},
	"gto": {kind: kSilk, side: sTop}, "gbo": {kind: kSilk, side: sBottom},
	"gts": {kind: kMask, side: sTop}, "gbs": {kind: kMask, side: sBottom},
	"gtp": {kind: kPaste, side: sTop}, "gbp": {kind: kPaste, side: sBottom},
	"gko": {kind: kOutline}, "gml": {kind: kOutline}, "gm1": {kind: kOutline},
	"cmp": {kind: kCopper, side: sTop}, "sol": {kind: kCopper, side: sBottom},
	"plc": {kind: kSilk, side: sTop}, "pls": {kind: kSilk, side: sBottom},
	"stc": {kind: kMask, side: sTop}, "sts": {kind: kMask, side: sBottom},
	"crc": {kind: kPaste, side: sTop}, "crs": {kind: kPaste, side: sBottom},
	"dim": {kind: kOutline}, "mil": {kind: kOutline},
	"top": {kind: kCopper, side: sTop}, "bot": {kind: kCopper, side: sBottom},
	"smt": {kind: kMask, side: sTop}, "smb": {kind: kMask, side: sBottom},
	"sst": {kind: kSilk, side: sTop}, "ssb": {kind: kSilk, side: sBottom},
	"spt": {kind: kPaste, side: sTop}, "spb": {kind: kPaste, side: sBottom},
	"oln": {kind: kOutline}, "out": {kind: kOutline},
}

// words splits a file name into lower-case words: at punctuation, between
// letters and digits and where lower case turns upper ("Gerber_TopSilkLayer"
// is gerber, top, silk, layer).
func words(s string) []string {
	var out []string
	var cur []rune
	flush := func() {
		if len(cur) > 0 {
			out = append(out, strings.ToLower(string(cur)))
			cur = cur[:0]
		}
	}
	rs := []rune(s)
	for i, r := range rs {
		switch {
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			flush()
			continue
		case i > 0 && len(cur) > 0:
			prev := rs[i-1]
			if unicode.IsDigit(r) != unicode.IsDigit(prev) || unicode.IsUpper(r) && unicode.IsLower(prev) {
				flush()
			}
		}
		cur = append(cur, r)
	}
	flush()
	return out
}

// byName tells a layer by its file name.
func byName(name string, drill bool) (function, bool) {
	base := path.Base(strings.ReplaceAll(name, "\\", "/"))
	stem, ext := base, ""
	if i := strings.LastIndexByte(base, '.'); i > 0 {
		stem, ext = base[:i], strings.ToLower(base[i+1:])
	}
	if !drill {
		if f, ok := extensions[ext]; ok {
			return f, true
		}
		// Protel inner layers: .g1, .g2 (signal), .gp1 (plane); mechanical
		// layers .gm2 and up hold drawings as often as the outline
		for _, pre := range []string{"g", "gp", "gl", "ly", "l"} {
			if n, err := strconv.Atoi(strings.TrimPrefix(ext, pre)); err == nil && strings.HasPrefix(ext, pre) && n > 0 && n < 100 {
				if pre == "l" || pre == "ly" {
					n-- // Eagle numbers the top layer 1
				}
				return function{kind: kCopper, side: sInner, index: n}, true
			}
		}
	}
	if f, ok := kicadSuffix(stem); ok && !drill {
		return f, true
	}
	ws := words(stem)
	has := func(names ...string) bool {
		for _, w := range ws {
			for _, n := range names {
				if w == n {
					return true
				}
			}
		}
		return false
	}
	if drill {
		f := function{kind: kDrill}
		switch {
		case has("npth", "nonplated", "unplated") || strings.Contains(strings.ToLower(stem), "non_plated") || strings.Contains(strings.ToLower(stem), "non-plated"):
			f.plated = "nonplated"
		case has("pth", "plated"):
			f.plated = "plated"
		}
		return f, true
	}
	var f function
	switch {
	case has("map", "drl", "drill", "drills", "holes", "fab", "fabrication", "assembly", "crtyd", "courtyard", "adhes", "adhesive", "glue", "dwgs", "drawings", "cmts", "comments", "eco1", "eco2", "user", "margin", "keepout", "job", "notes", "readme"):
		return function{}, false
	case has("outline", "edge", "cuts", "profile", "boardoutline", "dimension", "mill", "milling", "contour", "route", "rout", "border"):
		return function{kind: kOutline}, true
	case has("paste", "cream", "stencil", "solderpaste", "pastemask"):
		f.kind = kPaste
	case has("mask", "soldermask", "stop", "stopmask", "solderstop", "resist", "solderresist"):
		f.kind = kMask
	case has("silk", "silks", "silkscreen", "legend", "overlay", "nomenclature", "marking", "ss"):
		f.kind = kSilk
	case has("cu", "copper", "signal", "layer", "conductor", "plane", "etch", "trace", "traces"):
		f.kind = kCopper
	default:
		return function{}, false
	}
	switch {
	case has("top", "front", "f", "cmp", "component", "upper"):
		f.side = sTop
	case has("bottom", "bot", "back", "b", "sold", "lower"):
		f.side = sBottom
	default:
		// inner layers: In1.Cu (KiCad), InnerLayer1 (EasyEDA), Mid1, L2 (Eagle)
		for i, w := range ws {
			n, err := strconv.Atoi(w)
			if err != nil || n <= 0 || n >= 100 || i == 0 || f.kind != kCopper {
				continue
			}
			prev := ws[i-1]
			if prev == "layer" && i >= 2 && (ws[i-2] == "inner" || ws[i-2] == "mid") {
				prev = ws[i-2]
			}
			switch prev {
			case "in", "inner", "mid", "inr", "gp", "innerlayer":
				return function{kind: kCopper, side: sInner, index: n}, true
			case "l":
				if n > 1 {
					// Eagle numbers the top layer 1
					return function{kind: kCopper, side: sInner, index: n - 1}, true
				}
			}
		}
		f.side = sTop
	}
	return f, true
}

// kicadLayers are the names KiCad gives its layers, which end the names of
// the files it plots.
var kicadLayers = map[string]function{
	"f_cu": {kind: kCopper, side: sTop}, "b_cu": {kind: kCopper, side: sBottom},
	"f_mask": {kind: kMask, side: sTop}, "b_mask": {kind: kMask, side: sBottom},
	"f_silks": {kind: kSilk, side: sTop}, "b_silks": {kind: kSilk, side: sBottom},
	"f_silkscreen": {kind: kSilk, side: sTop}, "b_silkscreen": {kind: kSilk, side: sBottom},
	"f_paste": {kind: kPaste, side: sTop}, "b_paste": {kind: kPaste, side: sBottom},
	"edge_cuts": {kind: kOutline},
}

// kicadSuffix tells a layer by the KiCad layer name that ends a file's
// name ("board-F_Cu", "board-In1_Cu", "board-Edge_Cuts").
func kicadSuffix(stem string) (function, bool) {
	s := strings.ToLower(strings.ReplaceAll(stem, ".", "_"))
	for name, f := range kicadLayers {
		if strings.HasSuffix(s, name) && (len(s) == len(name) || strings.ContainsRune("-_ ", rune(s[len(s)-len(name)-1]))) {
			return f, true
		}
	}
	// In1_Cu to In30_Cu
	if t, ok := strings.CutSuffix(s, "_cu"); ok {
		if i := strings.LastIndex(t, "in"); i >= 0 && (i == 0 || strings.ContainsRune("-_ ", rune(t[i-1]))) {
			if n, err := strconv.Atoi(t[i+2:]); err == nil && n > 0 && n < 100 {
				return function{kind: kCopper, side: sInner, index: n}, true
			}
		}
	}
	return function{}, false
}
