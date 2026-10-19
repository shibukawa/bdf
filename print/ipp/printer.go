package ipp

import (
	"sort"
	"strconv"
	"strings"
)

// PrinterAttributes are the attributes GetPrinterAttributes asks for
// when given none: what Printer reads. media-col-database is not among
// them (it is large on some printers); ask for it by name.
var PrinterAttributes = []string{
	"printer-uri-supported", "uri-authentication-supported", "uri-security-supported",
	"printer-name", "printer-info", "printer-make-and-model", "printer-location", "printer-uuid", "printer-device-id",
	"printer-state", "printer-state-reasons", "printer-state-message", "printer-is-accepting-jobs",
	"ipp-versions-supported", "operations-supported", "printer-config-change-time",
	"document-format-supported", "document-format-default", "pdf-versions-supported",
	"media-supported", "media-ready", "media-default", "media-source-supported", "media-type-supported", "media-col-supported",
	"sides-supported", "sides-default", "print-color-mode-supported", "print-color-mode-default",
	"print-quality-supported", "print-quality-default", "printer-resolution-supported", "printer-resolution-default",
	"copies-supported", "copies-default", "page-ranges-supported", "print-scaling-supported", "print-scaling-default",
	"number-up-supported", "orientation-requested-supported", "multiple-document-handling-supported",
	"finishings-supported", "finishings-default", "output-bin-supported", "output-bin-default",
	"pwg-raster-document-resolution-supported", "pwg-raster-document-type-supported", "pwg-raster-document-sheet-back",
	"urf-supported", "pclm-strip-height-preferred", "pclm-raster-back-side", "pclm-source-resolution-supported",
	"job-password-supported", "job-password-encryption-supported", "job-account-id-supported", "job-accounting-user-id-supported",
	"identify-actions-supported", "job-creation-attributes-supported",
}

// Printer is what a printer says of itself in Get-Printer-Attributes:
// its names and state, and what it supports, read into fields for the
// attributes a job needs, with all of them in Attrs.
type Printer struct {
	URI      string   // the URL asked
	URIs     []string // printer-uri-supported
	Auth     []string // uri-authentication-supported, parallel to URIs
	Security []string // uri-security-supported, parallel to URIs

	Name      string
	Info      string
	MakeModel string
	Location  string
	UUID      string
	DeviceID  string

	State         PrinterState
	Reasons       []string
	StateMessage  string
	AcceptingJobs bool

	Versions         []string
	Operations       []Op
	ConfigChangeTime int32

	Formats       []string // document-format-supported
	DefaultFormat string
	PDFVersions   []string

	Media        []string // media-supported
	MediaReady   []string
	MediaDefault string
	MediaSources []string
	MediaTypes   []string
	MediaCol     bool    // media-col is accepted (media-col-supported names members)
	MediaColDB   []Media // media-col-database, when asked for

	Sides             []string
	SidesDefault      string
	ColorModes        []string
	ColorModeDefault  string
	Qualities         []int32 // print-quality-supported: 3 draft, 4 normal, 5 high
	QualityDefault    int32
	Resolutions       []Resolution
	ResolutionDefault Resolution
	Copies            Range // copies-supported; {1,1} when not said
	PageRanges        bool
	Scaling           []string
	ScalingDefault    string
	NumberUp          []int32
	Orientations      []int32
	DocumentHandling  []string
	Finishings        []int32
	FinishingsDefault []int32
	OutputBins        []string
	OutputBinDefault  string

	PWGRasterResolutions []Resolution
	PWGRasterTypes       []string
	PWGRasterSheetBack   string
	URF                  []string
	PCLm                 bool

	JobPassword           int32 // job-password-supported: the length, 0 when not supported
	JobPasswordEncryption []string
	JobAccountID          bool
	JobAccountingUserID   bool
	IdentifyActions       []string
	JobCreationAttrs      []string

	// Attrs are the printer attributes as received.
	Attrs []Attribute
}

// Media is a size of paper a printer has or takes (an entry of
// media-col-database or media-ready in collection form).
type Media struct {
	SizeName string // media-size-name, when given
	// X and Y are the width and height in hundredths of a millimetre.
	X, Y   int32
	Source string
	Type   string
	Color  string
	// Margins in hundredths of a millimetre: left, right, top, bottom.
	Left, Right, Top, Bottom int32
}

// ParsePrinter reads a Get-Printer-Attributes response.
func ParsePrinter(m *Message) *Printer {
	p := &Printer{Copies: Range{1, 1}}
	g := m.Group(TagPrinter)
	if g == nil {
		return p
	}
	p.Attrs = g.Attrs
	p.URIs = g.Get("printer-uri-supported").Strings()
	p.Auth = g.Get("uri-authentication-supported").Strings()
	p.Security = g.Get("uri-security-supported").Strings()
	p.Name = g.Get("printer-name").Text()
	p.Info = g.Get("printer-info").Text()
	p.MakeModel = g.Get("printer-make-and-model").Text()
	p.Location = g.Get("printer-location").Text()
	p.UUID = g.Get("printer-uuid").Text()
	p.DeviceID = g.Get("printer-device-id").Text()
	if n, ok := g.Get("printer-state").Int(); ok {
		p.State = PrinterState(n)
	}
	p.Reasons = g.Get("printer-state-reasons").Strings()
	p.StateMessage = g.Get("printer-state-message").Text()
	p.AcceptingJobs, _ = g.Get("printer-is-accepting-jobs").Bool()
	p.Versions = g.Get("ipp-versions-supported").Strings()
	for _, n := range g.Get("operations-supported").Ints() {
		p.Operations = append(p.Operations, Op(n))
	}
	p.ConfigChangeTime, _ = g.Get("printer-config-change-time").Int()
	p.Formats = g.Get("document-format-supported").Strings()
	p.DefaultFormat = g.Get("document-format-default").Text()
	p.PDFVersions = g.Get("pdf-versions-supported").Strings()
	p.Media = g.Get("media-supported").Strings()
	p.MediaReady = g.Get("media-ready").Strings()
	p.MediaDefault = g.Get("media-default").Text()
	p.MediaSources = g.Get("media-source-supported").Strings()
	p.MediaTypes = g.Get("media-type-supported").Strings()
	p.MediaCol = len(g.Get("media-col-supported").Strings()) > 0
	for _, c := range g.Get("media-col-database").Collections() {
		p.MediaColDB = append(p.MediaColDB, ParseMedia(c))
	}
	p.Sides = g.Get("sides-supported").Strings()
	p.SidesDefault = g.Get("sides-default").Text()
	p.ColorModes = g.Get("print-color-mode-supported").Strings()
	p.ColorModeDefault = g.Get("print-color-mode-default").Text()
	p.Qualities = g.Get("print-quality-supported").Ints()
	p.QualityDefault, _ = g.Get("print-quality-default").Int()
	p.Resolutions = g.Get("printer-resolution-supported").Resolutions()
	if r := g.Get("printer-resolution-default").Resolutions(); len(r) > 0 {
		p.ResolutionDefault = r[0]
	}
	if r := g.Get("copies-supported").Ranges(); len(r) > 0 {
		p.Copies = r[0]
	}
	p.PageRanges, _ = g.Get("page-ranges-supported").Bool()
	p.Scaling = g.Get("print-scaling-supported").Strings()
	p.ScalingDefault = g.Get("print-scaling-default").Text()
	p.NumberUp = g.Get("number-up-supported").Ints()
	for _, r := range g.Get("number-up-supported").Ranges() {
		for n := r.Lo; n <= r.Hi && n <= 16; n++ {
			p.NumberUp = append(p.NumberUp, n)
		}
	}
	p.Orientations = g.Get("orientation-requested-supported").Ints()
	p.DocumentHandling = g.Get("multiple-document-handling-supported").Strings()
	p.Finishings = g.Get("finishings-supported").Ints()
	p.FinishingsDefault = g.Get("finishings-default").Ints()
	p.OutputBins = g.Get("output-bin-supported").Strings()
	p.OutputBinDefault = g.Get("output-bin-default").Text()
	p.PWGRasterResolutions = g.Get("pwg-raster-document-resolution-supported").Resolutions()
	p.PWGRasterTypes = g.Get("pwg-raster-document-type-supported").Strings()
	p.PWGRasterSheetBack = g.Get("pwg-raster-document-sheet-back").Text()
	p.URF = g.Get("urf-supported").Strings()
	p.PCLm = g.Get("pclm-strip-height-preferred") != nil
	p.JobPassword, _ = g.Get("job-password-supported").Int()
	p.JobPasswordEncryption = g.Get("job-password-encryption-supported").Strings()
	p.JobAccountID = g.Get("job-account-id-supported") != nil
	p.JobAccountingUserID = g.Get("job-accounting-user-id-supported") != nil
	p.IdentifyActions = g.Get("identify-actions-supported").Strings()
	p.JobCreationAttrs = g.Get("job-creation-attributes-supported").Strings()
	return p
}

// ParseMedia reads a media-col collection.
func ParseMedia(c []Attribute) Media {
	m := Media{SizeName: Get(c, "media-size-name").Text(), Source: Get(c, "media-source").Text(), Type: Get(c, "media-type").Text(),
		Color: Get(c, "media-color").Text()}
	if sz := Get(c, "media-size").Collections(); len(sz) > 0 {
		m.X, _ = Get(sz[0], "x-dimension").Int()
		m.Y, _ = Get(sz[0], "y-dimension").Int()
	}
	m.Left, _ = Get(c, "media-left-margin").Int()
	m.Right, _ = Get(c, "media-right-margin").Int()
	m.Top, _ = Get(c, "media-top-margin").Int()
	m.Bottom, _ = Get(c, "media-bottom-margin").Int()
	return m
}

// Supports reports whether the printer takes documents of the format
// (a media type such as "application/pdf" or "image/pwg-raster").
func (p *Printer) Supports(format string) bool {
	for _, f := range p.Formats {
		if strings.EqualFold(f, format) {
			return true
		}
	}
	return false
}

// SupportsOp reports whether the printer names the operation among
// operations-supported.
func (p *Printer) SupportsOp(op Op) bool {
	for _, o := range p.Operations {
		if o == op {
			return true
		}
	}
	return false
}

// Has reports whether an attribute of the name is among the printer's,
// and returns it.
func (p *Printer) Has(name string) (*Attribute, bool) {
	for i := range p.Attrs {
		if p.Attrs[i].Name == name {
			return &p.Attrs[i], true
		}
	}
	return nil, false
}

// Summary lines the printer up for people: a name, the model, the
// state, the formats and the main capabilities, one a line.
func (p *Printer) Summary() []string {
	var lines []string
	add := func(k, v string) {
		if v != "" {
			lines = append(lines, k+": "+v)
		}
	}
	add("name", p.Name)
	add("model", p.MakeModel)
	add("info", p.Info)
	add("location", p.Location)
	add("uri", strings.Join(p.URIs, ", "))
	st := p.State.String()
	if len(p.Reasons) > 0 && !(len(p.Reasons) == 1 && p.Reasons[0] == "none") {
		st += " (" + strings.Join(p.Reasons, ", ") + ")"
	}
	if !p.AcceptingJobs {
		st += ", not accepting jobs"
	}
	add("state", st)
	add("ipp", strings.Join(p.Versions, ", "))
	add("auth", strings.Join(p.Auth, ", "))
	add("formats", strings.Join(p.Formats, ", "))
	add("media", strings.Join(p.Media, ", "))
	add("media-ready", strings.Join(p.MediaReady, ", "))
	add("media-default", p.MediaDefault)
	add("sources", strings.Join(p.MediaSources, ", "))
	add("types", strings.Join(p.MediaTypes, ", "))
	add("sides", strings.Join(p.Sides, ", "))
	add("color", strings.Join(p.ColorModes, ", "))
	add("quality", ints(p.Qualities))
	add("resolution", resolutions(p.Resolutions))
	if p.Copies != (Range{1, 1}) {
		add("copies", p.Copies.String())
	}
	if p.PageRanges {
		add("page-ranges", "yes")
	}
	add("scaling", strings.Join(p.Scaling, ", "))
	add("finishings", ints(p.Finishings))
	add("output-bins", strings.Join(p.OutputBins, ", "))
	add("pwg-raster", strings.TrimSpace(resolutions(p.PWGRasterResolutions)+" "+strings.Join(p.PWGRasterTypes, ", ")+" "+p.PWGRasterSheetBack))
	add("urf", strings.Join(p.URF, ","))
	if p.JobPassword > 0 {
		add("job-password", "up to "+itoa(int(p.JobPassword))+" characters")
	}
	if p.JobAccountID {
		add("job-account-id", "yes")
	}
	return lines
}

func ints(ns []int32) string {
	s := make([]string, len(ns))
	for i, n := range ns {
		s[i] = itoa(int(n))
	}
	return strings.Join(s, ", ")
}

func resolutions(rs []Resolution) string {
	s := make([]string, len(rs))
	for i, r := range rs {
		s[i] = r.String()
	}
	return strings.Join(s, ", ")
}

func itoa(n int) string { return strconv.Itoa(n) }

// SortedFormats returns the formats the printer supports, sorted.
func (p *Printer) SortedFormats() []string {
	out := append([]string(nil), p.Formats...)
	sort.Strings(out)
	return out
}
