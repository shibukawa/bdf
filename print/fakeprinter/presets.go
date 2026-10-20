package fakeprinter

import (
	"fmt"
	"slices"
	"sort"

	"github.com/shibukawa/bdf/print/ipp"
	"github.com/shibukawa/bdf/print/pdl"
)

// Preset is a kind of printer bdf prints to.
type Preset struct {
	Name, Title string
	// Transport is how the kind is reached: "ipp", "socket", "lpd" or
	// "serial".
	Transport string
	make      func() *Printer
}

// New makes a printer of the preset, named name.
func (ps Preset) New(name string) *Printer {
	p := ps.make()
	p.Name = name
	p.UUID = NewUUID()
	return p
}

var dpi = func(ns ...int) []ipp.Resolution {
	var out []ipp.Resolution
	for _, n := range ns {
		out = append(out, ipp.Resolution{X: int32(n), Y: int32(n), Units: ipp.UnitsDPI})
	}
	return out
}

// officeMedia are the media of an office printer in Japan.
var officeMedia = []string{
	"iso_a4_210x297mm", "iso_a3_297x420mm", "iso_a5_148x210mm", "jis_b4_257x364mm", "jis_b5_182x257mm", "na_letter_8.5x11in",
	"na_legal_8.5x14in", "jpn_hagaki_100x148mm", "iso_dl_110x220mm", "jpn_chou3_120x235mm",
}

var presets = []Preset{
	{Name: "office-mfp", Title: "Office multifunction printer (IPP Everywhere: PDF, PWG Raster, Apple Raster, JPEG)", Transport: "ipp", make: func() *Printer {
		p := New()
		p.MakeModel = "BDF Fake Office MFP"
		p.Formats = []string{"application/pdf", "image/urf", "image/pwg-raster", "image/jpeg", "application/octet-stream"}
		p.Media = officeMedia
		p.Sources = []string{"auto", "main", "manual", "tray-1", "tray-2"}
		p.Types = []string{"auto", "stationery", "envelope", "labels", "photographic"}
		p.Finishings = []int{3, 4, 20}
		p.OutputBins = []string{"face-down", "top", "stacker-1"}
		p.URF = []string{"V1.4", "W8", "SRGB24", "RS300-600", "DM1", "CP99", "IS1-2-4", "MT1-2-3-7", "FN3-4-20"}
		return p
	}},
	{Name: "airprint-only", Title: "Printer of AirPrint only (Apple Raster and JPEG, no PDF)", Transport: "ipp", make: func() *Printer {
		p := New()
		p.MakeModel = "BDF Fake AirPrint Printer"
		p.Formats = []string{"image/urf", "image/jpeg", "application/octet-stream"}
		p.PWGRasterTypes, p.PWGRasterResolutions = nil, nil
		return p
	}},
	{Name: "pwg-only", Title: "Printer of PWG Raster only, which makes no copies", Transport: "ipp", make: func() *Printer {
		p := New()
		p.MakeModel = "BDF Fake Raster Printer"
		p.Formats = []string{"image/pwg-raster", "application/octet-stream"}
		p.Copies = ipp.Range{Lo: 1, Hi: 1}
		p.URF = nil
		p.ColorModes = []string{"auto", "monochrome"}
		p.PWGRasterTypes = []string{"black_1", "sgray_8"}
		p.SheetBack = "flipped"
		return p
	}},
	{Name: "ps-office", Title: "Office printer of PostScript and PDF", Transport: "ipp", make: func() *Printer {
		p := New()
		p.MakeModel = "BDF Fake PostScript Printer"
		p.Formats = []string{"application/postscript", "application/pdf", "application/octet-stream"}
		p.PWGRasterTypes, p.PWGRasterResolutions, p.URF = nil, nil, nil
		p.Media = officeMedia
		return p
	}},
	{Name: "pclm-inkjet", Title: "Inkjet printer of PCLm (Mopria) and JPEG", Transport: "ipp", make: func() *Printer {
		p := New()
		p.MakeModel = "BDF Fake Inkjet"
		p.Formats = []string{"application/PCLm", "image/jpeg", "application/octet-stream"}
		p.PWGRasterTypes, p.PWGRasterResolutions, p.URF = nil, nil, nil
		p.PCLmStripHeight, p.PCLmBackSide = 16, "rotated"
		p.PCLmResolutions = dpi(300, 600)
		p.PCLmCompression = []string{"flate", "jpeg"}
		p.Media = []string{"iso_a4_210x297mm", "jpn_hagaki_100x148mm", "na_index-4x6_4x6in", "iso_a5_148x210mm"}
		return p
	}},
	{Name: "legacy-ipp11", Title: "Old printer of IPP/1.1 only, without chunked requests", Transport: "ipp", make: func() *Printer {
		p := New()
		p.MakeModel = "BDF Fake Old Printer"
		p.Only11, p.NoChunked = true, true
		p.Formats = []string{"application/pdf", "application/postscript", "application/octet-stream"}
		p.PWGRasterTypes, p.PWGRasterResolutions, p.URF = nil, nil, nil
		return p
	}},
	{Name: "jetdirect-ps", Title: "Printer of PostScript on the raw port 9100 (PJL)", Transport: "socket", make: func() *Printer {
		p := New()
		p.MakeModel, p.Language = "BDF Fake JetDirect", "postscript"
		p.Formats = []string{"application/postscript"}
		return p
	}},
	{Name: "lpd-legacy", Title: "LPD queue that takes anything", Transport: "lpd", make: func() *Printer {
		p := New()
		p.MakeModel = "BDF Fake LPD Queue"
		return p
	}},
}

// profilePresets are the label, receipt and dot-matrix printers: a preset
// of each profile of print/pdl the farm offers, named as the maker's model.
var profilePresets = []struct{ name, profile, transport string }{
	{"zebra-zd421-203", "zd421-203", "socket"},
	{"zebra-zd421-300", "zd421-300", "socket"},
	{"sato-cl4nx", "cl4nx-203", "socket"},
	{"tec-bex4t1", "b-ex4t1-203", "socket"},
	{"epson-tm-t88", "tm-t88", "socket"},
	{"pos-58", "pos-58", "socket"},
	{"epson-vp-f2000", "vp-f2000", "serial"},
}

func init() {
	for _, pp := range profilePresets {
		prof, ok := pdl.LookupProfile(pp.profile)
		if !ok {
			panic("fakeprinter: no profile " + pp.profile)
		}
		presets = append(presets, Preset{Name: pp.name, Title: prof.Title + " (" + prof.Language + ")", Transport: pp.transport,
			make: func() *Printer { return FromProfile(prof) }})
	}
}

// FromProfile makes a printer of a profile of print/pdl: its language,
// dots and resolution, over the raw transports.
func FromProfile(prof pdl.Profile) *Printer {
	p := New()
	p.MakeModel = prof.Title
	p.Language, p.Width, p.DPI = prof.Language, prof.Width, prof.DPI
	p.Formats = []string{"application/octet-stream", "application/vnd.cups-raw"}
	p.PWGRasterTypes, p.PWGRasterResolutions, p.URF = nil, nil, nil
	p.Resolutions = dpi(prof.DPI)
	p.ColorModes = []string{"monochrome"}
	p.Sides = []string{"one-sided"}
	widthMM := float64(prof.Width) * 25.4 / float64(prof.DPI)
	switch prof.Language {
	case pdl.ESCPOS:
		p.Media = []string{fmt.Sprintf("roll_max_%.0fx3000mm", widthMM)}
	case pdl.ESCP:
		p.Media = []string{"na_fanfold-us_11x14.875in", "iso_a4_210x297mm"}
	default:
		p.Media = []string{"oe_4x6-label_4x6in", "om_100x150-label_100x150mm", "om_60x40-label_60x40mm"}
	}
	p.MediaReady = p.Media[:1]
	return p
}

// LookupPreset finds a preset by name.
func LookupPreset(name string) (Preset, bool) {
	i := slices.IndexFunc(presets, func(p Preset) bool { return p.Name == name })
	if i < 0 {
		return Preset{}, false
	}
	return presets[i], true
}

// Presets are the presets, by name.
func Presets() []Preset {
	out := slices.Clone(presets)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// CUPSPreset is a CUPS server of three printers and a class: "office"
// (office-mfp), "laser" (pwg-only) and "ps" (ps-office); the class
// "floor" of office and laser. office is the default.
func CUPSPreset() *CUPS {
	c := NewCUPS()
	for _, x := range [][2]string{{"office", "office-mfp"}, {"laser", "pwg-only"}, {"ps", "ps-office"}} {
		ps, _ := LookupPreset(x[1])
		c.Add(ps.New(x[0]))
	}
	c.AddClass("floor", "office", "laser")
	return c
}
