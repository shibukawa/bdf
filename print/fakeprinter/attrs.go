package fakeprinter

import (
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/print/ipp"
)

// printerGroups are the groups of printer attributes requested-attributes
// names: "job-template" the -default, -supported, -ready and -database of
// job template attributes, "printer-description" the others.
func printerGroups(attrs []ipp.Attribute) map[string][]string {
	var tmpl, desc []string
	for _, a := range attrs {
		base := a.Name
		for _, suf := range []string{"-default", "-supported", "-ready", "-database"} {
			base = strings.TrimSuffix(base, suf)
		}
		if base != a.Name && (slices.Contains(jobTemplate, base) || base == "media-col" || base == "media") {
			tmpl = append(tmpl, a.Name)
		} else {
			desc = append(desc, a.Name)
		}
	}
	return map[string][]string{"job-template": tmpl, "printer-description": desc}
}

func (p *Printer) getPrinterAttributes(q *request) {
	if f := q.op.Get("document-format"); f != nil && !p.supports(f.Text()) {
		q.fail(ipp.StatusClientErrorDocumentFormatNotSupported, "document-format %s", f.Text())
		q.res.Add(ipp.TagUnsupportedGroup, *f)
		return
	}
	attrs := p.attributes(q)
	req := requested(q.op, nil)
	// media-col-database is told when asked by name, as printers do
	if !slices.Contains(req, "media-col-database") {
		attrs = slices.DeleteFunc(attrs, func(a ipp.Attribute) bool { return a.Name == "media-col-database" })
	}
	q.res.Add(ipp.TagPrinter, filter(attrs, req, printerGroups(attrs))...)
}

// mediaSize is the size of a media keyword of PWG 5101.1 (its last part,
// "210x297mm" or "8.5x11in") in hundredths of a millimetre.
func mediaSize(name string) (w, h int, ok bool) {
	i := strings.LastIndexByte(name, '_')
	dims := name[i+1:]
	unit := 0.0
	switch {
	case strings.HasSuffix(dims, "mm"):
		unit, dims = 100, strings.TrimSuffix(dims, "mm")
	case strings.HasSuffix(dims, "in"):
		unit, dims = 2540, strings.TrimSuffix(dims, "in")
	default:
		return 0, 0, false
	}
	ws, hs, found := strings.Cut(dims, "x")
	if !found {
		return 0, 0, false
	}
	fw, err1 := strconv.ParseFloat(ws, 64)
	fh, err2 := strconv.ParseFloat(hs, 64)
	if err1 != nil || err2 != nil || fw <= 0 || fh <= 0 {
		return 0, 0, false
	}
	return int(math.Round(fw * unit)), int(math.Round(fh * unit)), true
}

// margins are the printer's margins in hundredths of a millimetre: none
// for a printer of labels or receipts, a sixth of an inch for others.
func (p *Printer) margins() int {
	if p.Language != "" && p.Language != "postscript" && p.Language != "pdf" && p.Language != "pcl" {
		return 0
	}
	return 423
}

// attributes are the printer's attributes, as an IPP Everywhere printer
// tells them. The printer is locked.
func (p *Printer) attributes(q *request) []ipp.Attribute {
	uri := q.printerURI()
	auth := "requesting-user-name"
	if p.Username != "" {
		auth = "basic"
	}
	sec := "none"
	if q.secure {
		sec = "tls"
	}
	m := p.margins()
	var mediaDB, mediaReady, mediaSizes []ipp.Value
	for _, name := range p.Media {
		w, h, ok := mediaSize(name)
		if !ok {
			continue
		}
		size := ipp.Collection(ipp.Attr("x-dimension", ipp.Integer(w)), ipp.Attr("y-dimension", ipp.Integer(h)))
		mediaSizes = append(mediaSizes, size)
		col := ipp.Collection(ipp.Attr("media-key", ipp.Keyword(name)), ipp.Attr("media-size", size), ipp.Attr("media-size-name", ipp.Keyword(name)),
			ipp.Attr("media-source", ipp.Keyword(first(p.Sources, "auto"))), ipp.Attr("media-type", ipp.Keyword(first(p.Types, "auto"))),
			ipp.Attr("media-top-margin", ipp.Integer(m)), ipp.Attr("media-bottom-margin", ipp.Integer(m)),
			ipp.Attr("media-left-margin", ipp.Integer(m)), ipp.Attr("media-right-margin", ipp.Integer(m)))
		mediaDB = append(mediaDB, col)
		if slices.Contains(p.MediaReady, name) {
			mediaReady = append(mediaReady, col)
		}
	}
	if len(mediaReady) == 0 && len(mediaDB) > 0 {
		mediaReady = mediaDB[:1]
	}
	ops := make([]ipp.Value, len(operations))
	for i, op := range operations {
		ops[i] = ipp.Enum(int(op))
	}
	color := slices.Contains(p.ColorModes, "color")
	ready := p.MediaReady
	if len(ready) == 0 && len(p.Media) > 0 {
		ready = p.Media[:1]
	}
	info := p.Info
	if info == "" {
		info = p.Name
	}
	cmd := []string{}
	for _, f := range p.Formats {
		switch f {
		case "application/pdf":
			cmd = append(cmd, "PDF")
		case "image/pwg-raster":
			cmd = append(cmd, "PWG")
		case "image/urf":
			cmd = append(cmd, "URF")
		case "image/jpeg":
			cmd = append(cmd, "JPEG")
		case "application/postscript":
			cmd = append(cmd, "POSTSCRIPT")
		case "application/PCLm":
			cmd = append(cmd, "PCLM")
		}
	}
	attrs := []ipp.Attribute{
		ipp.Attr("printer-uri-supported", ipp.URI(uri)),
		ipp.Attr("uri-authentication-supported", ipp.Keyword(auth)),
		ipp.Attr("uri-security-supported", ipp.Keyword(sec)),
		ipp.Attr("printer-name", ipp.Name(p.Name)),
		ipp.Attr("printer-info", ipp.Text(info)),
		ipp.Attr("printer-location", ipp.Text(p.Location)),
		ipp.Attr("printer-make-and-model", ipp.Text(p.MakeModel)),
		ipp.Attr("printer-more-info", ipp.URI(strings.Replace(uri, "ipp", "http", 1))),
		ipp.Attr("printer-device-id", ipp.Text("MFG:BDF;MDL:Fake Printer;CMD:"+strings.Join(cmd, ",")+";")),
		ipp.Attr("printer-uuid", ipp.URI(p.UUID)),
		ipp.Attr("printer-kind", ipp.Keyword("document")),
		ipp.Attr("printer-geo-location", ipp.OutOfBand(ipp.TagUnknown)),
		ipp.Attr("printer-organization", ipp.Text("")),
		ipp.Attr("printer-organizational-unit", ipp.Text("")),
		ipp.Attr("printer-icons", ipp.URI(strings.Replace(uri, "ipp", "http", 1)+"/icon.png")),
		ipp.Attr("printer-state", ipp.Enum(int(p.State))),
		ipp.Attr("printer-state-reasons", keywords(p.Reasons)...),
		ipp.Attr("printer-state-message", ipp.Text("")),
		ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(p.Accepting)),
		ipp.Attr("printer-up-time", ipp.Integer(p.upTime(time.Now()))),
		ipp.Attr("printer-current-time", ipp.DateTime(time.Now())),
		ipp.Attr("printer-state-change-time", ipp.Integer(p.upTime(p.stateSince))),
		ipp.Attr("printer-state-change-date-time", ipp.DateTime(p.stateSince)),
		ipp.Attr("printer-config-change-time", ipp.Integer(1)),
		ipp.Attr("printer-config-change-date-time", ipp.DateTime(p.started)),
		ipp.Attr("queued-job-count", ipp.Integer(p.queued())),
		ipp.Attr("ipp-versions-supported", p.versions()...),
		ipp.Attr("ipp-features-supported", ipp.Keyword("ipp-everywhere")),
		ipp.Attr("operations-supported", ops...),
		ipp.Attr("charset-configured", ipp.Charset("utf-8")),
		ipp.Attr("charset-supported", ipp.Charset("utf-8"), ipp.Charset("us-ascii")),
		ipp.Attr("natural-language-configured", ipp.Language("en")),
		ipp.Attr("generated-natural-language-supported", ipp.Language("en")),
		ipp.Attr("compression-supported", ipp.Keyword("none"), ipp.Keyword("deflate"), ipp.Keyword("gzip")),
		ipp.Attr("document-format-supported", mimes(p.Formats)...),
		ipp.Attr("document-format-default", ipp.MimeType(first(p.Formats, "application/octet-stream"))),
		ipp.Attr("multiple-document-jobs-supported", ipp.Boolean(false)),
		ipp.Attr("multiple-operation-time-out", ipp.Integer(300)),
		ipp.Attr("multiple-operation-time-out-action", ipp.Keyword("abort-job")),
		ipp.Attr("job-ids-supported", ipp.Boolean(true)),
		ipp.Attr("preferred-attributes-supported", ipp.Boolean(false)),
		ipp.Attr("which-jobs-supported", ipp.Keyword("completed"), ipp.Keyword("not-completed"), ipp.Keyword("all")),
		ipp.Attr("job-creation-attributes-supported", keywords(jobTemplate)...),
		ipp.Attr("printer-get-attributes-supported", ipp.Keyword("document-format")),
		ipp.Attr("overrides-supported", ipp.Keyword("document-number"), ipp.Keyword("pages")),
		ipp.Attr("identify-actions-supported", ipp.Keyword("display"), ipp.Keyword("sound")),
		ipp.Attr("identify-actions-default", ipp.Keyword("display")),
		ipp.Attr("copies-supported", ipp.NewRange(int(p.Copies.Lo), int(p.Copies.Hi))),
		ipp.Attr("copies-default", ipp.Integer(1)),
		ipp.Attr("page-ranges-supported", ipp.Boolean(p.PageRanges)),
		ipp.Attr("job-hold-until-supported", ipp.Keyword("no-hold"), ipp.Keyword("indefinite")),
		ipp.Attr("job-hold-until-default", ipp.Keyword("no-hold")),
		ipp.Attr("job-priority-supported", ipp.Integer(100)),
		ipp.Attr("job-priority-default", ipp.Integer(50)),
		ipp.Attr("job-sheets-supported", ipp.Keyword("none")),
		ipp.Attr("job-sheets-default", ipp.Keyword("none")),
		ipp.Attr("multiple-document-handling-supported", ipp.Keyword("separate-documents-collated-copies"), ipp.Keyword("separate-documents-uncollated-copies")),
		ipp.Attr("multiple-document-handling-default", ipp.Keyword("separate-documents-collated-copies")),
		ipp.Attr("number-up-supported", ipp.Integer(1)),
		ipp.Attr("number-up-default", ipp.Integer(1)),
		ipp.Attr("orientation-requested-supported", ipp.Enum(3), ipp.Enum(4), ipp.Enum(5), ipp.Enum(6)),
		ipp.Attr("orientation-requested-default", ipp.OutOfBand(ipp.TagNoValue)),
		ipp.Attr("print-quality-supported", enums(p.Qualities)...),
		ipp.Attr("print-quality-default", ipp.Enum(4)),
		ipp.Attr("printer-resolution-supported", resolutions(p.Resolutions)...),
		ipp.Attr("printer-resolution-default", resolutions(p.Resolutions[:1])...),
		ipp.Attr("print-scaling-supported", keywords(p.Scaling)...),
		ipp.Attr("print-scaling-default", ipp.Keyword(first(p.Scaling, "auto"))),
		ipp.Attr("print-color-mode-supported", keywords(p.ColorModes)...),
		ipp.Attr("print-color-mode-default", ipp.Keyword(first(p.ColorModes, "auto"))),
		ipp.Attr("color-supported", ipp.Boolean(color)),
		ipp.Attr("print-content-optimize-supported", ipp.Keyword("auto")),
		ipp.Attr("print-content-optimize-default", ipp.Keyword("auto")),
		ipp.Attr("print-rendering-intent-supported", ipp.Keyword("auto")),
		ipp.Attr("print-rendering-intent-default", ipp.Keyword("auto")),
		ipp.Attr("sides-supported", keywords(p.Sides)...),
		ipp.Attr("sides-default", ipp.Keyword(first(p.Sides, "one-sided"))),
		ipp.Attr("finishings-supported", enums(p.Finishings)...),
		ipp.Attr("finishings-default", ipp.Enum(3)),
		ipp.Attr("media-supported", keywords(p.Media)...),
		ipp.Attr("media-default", ipp.Keyword(first(ready, "iso_a4_210x297mm"))),
		ipp.Attr("media-ready", keywords(ready)...),
		ipp.Attr("media-col-supported", ipp.Keyword("media-bottom-margin"), ipp.Keyword("media-left-margin"), ipp.Keyword("media-right-margin"),
			ipp.Keyword("media-size"), ipp.Keyword("media-size-name"), ipp.Keyword("media-source"), ipp.Keyword("media-top-margin"), ipp.Keyword("media-type")),
		ipp.Attr("media-source-supported", keywords(p.Sources)...),
		ipp.Attr("media-type-supported", keywords(p.Types)...),
		ipp.Attr("media-top-margin-supported", ipp.Integer(m)),
		ipp.Attr("media-bottom-margin-supported", ipp.Integer(m)),
		ipp.Attr("media-left-margin-supported", ipp.Integer(m)),
		ipp.Attr("media-right-margin-supported", ipp.Integer(m)),
		ipp.Attr("output-bin-supported", keywords(p.OutputBins)...),
		ipp.Attr("output-bin-default", ipp.Keyword(first(p.OutputBins, "face-down"))),
		ipp.Attr("pages-per-minute", ipp.Integer(20)),
		ipp.Attr("printer-supply", ipp.OctetString([]byte("index=1;class=supplyThatIsConsumed;type=toner;unit=percent;maxcapacity=100;level=80;colorantname=black;"))),
		ipp.Attr("printer-supply-description", ipp.Text("Black Toner")),
		ipp.Attr("printer-supply-info-uri", ipp.URI(strings.Replace(uri, "ipp", "http", 1)+"/supplies")),
		ipp.Attr("printer-firmware-name", ipp.Name("fakeprinter")),
		ipp.Attr("printer-firmware-string-version", ipp.Text("bdf fakeprinter")),
		ipp.Attr("printer-firmware-version", ipp.OctetString([]byte{1})),
		ipp.Attr("job-k-octets-supported", ipp.NewRange(0, 262144)),
		ipp.Attr("pdl-override-supported", ipp.Keyword("attempted")),
	}
	if len(mediaDB) > 0 {
		attrs = append(attrs, ipp.Attr("media-col-database", mediaDB...), ipp.Attr("media-col-ready", mediaReady...),
			ipp.Attr("media-col-default", mediaReady[0]), ipp.Attr("media-size-supported", mediaSizes...))
	}
	if color {
		attrs = append(attrs, ipp.Attr("pages-per-minute-color", ipp.Integer(20)))
	}
	if slices.Contains(p.Formats, "application/pdf") {
		attrs = append(attrs, ipp.Attr("pdf-versions-supported", ipp.Keyword("adobe-1.3"), ipp.Keyword("adobe-1.4"), ipp.Keyword("adobe-1.5"),
			ipp.Keyword("adobe-1.6"), ipp.Keyword("adobe-1.7"), ipp.Keyword("iso-32000-1_2008"), ipp.Keyword("pwg-5102.3")),
			ipp.Attr("document-format-preferred", ipp.MimeType("application/pdf")))
	}
	if len(p.PWGRasterTypes) > 0 {
		attrs = append(attrs, ipp.Attr("pwg-raster-document-resolution-supported", resolutions(p.PWGRasterResolutions)...),
			ipp.Attr("pwg-raster-document-type-supported", keywords(p.PWGRasterTypes)...),
			ipp.Attr("pwg-raster-document-sheet-back", ipp.Keyword(first([]string{p.SheetBack}, "normal"))))
	}
	if len(p.URF) > 0 {
		attrs = append(attrs, ipp.Attr("urf-supported", keywords(p.URF)...))
	}
	if p.PCLmStripHeight > 0 {
		attrs = append(attrs, ipp.Attr("pclm-strip-height-preferred", ipp.Integer(p.PCLmStripHeight)),
			ipp.Attr("pclm-raster-back-side", ipp.Keyword(p.PCLmBackSide)),
			ipp.Attr("pclm-source-resolution-supported", resolutions(p.PCLmResolutions)...),
			ipp.Attr("pclm-compression-method-preferred", keywords(p.PCLmCompression)...))
	}
	return attrs
}

// cupsAttrs are the attributes CUPS-Get-Printers tells of a printer.
func (p *Printer) cupsAttrs(q *request, uri string) []ipp.Attribute {
	info := p.Info
	if info == "" {
		info = p.Name
	}
	return []ipp.Attribute{
		ipp.Attr("printer-name", ipp.Name(p.Name)),
		ipp.Attr("printer-uri-supported", ipp.URI(uri)),
		ipp.Attr("printer-info", ipp.Text(info)),
		ipp.Attr("printer-location", ipp.Text(p.Location)),
		ipp.Attr("printer-make-and-model", ipp.Text(p.MakeModel)),
		ipp.Attr("printer-state", ipp.Enum(int(p.State))),
		ipp.Attr("printer-state-reasons", keywords(p.Reasons)...),
		ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(p.Accepting)),
		ipp.Attr("printer-type", ipp.Integer(cupsPrinterType(p))),
		ipp.Attr("device-uri", ipp.URI(uri)),
		ipp.Attr("printer-uuid", ipp.URI(p.UUID)),
		ipp.Attr("queued-job-count", ipp.Integer(p.queued())),
	}
}

// cupsPrinterType is CUPS's printer-type of a printer: its abilities as
// bits (CUPS_PRINTER_BW, COLOR, DUPLEX, STAPLE, COPIES, COLLATE …).
func cupsPrinterType(p *Printer) int {
	t := 0x0004 // CUPS_PRINTER_BW
	if slices.Contains(p.ColorModes, "color") {
		t |= 0x0008
	}
	if len(p.Sides) > 1 {
		t |= 0x0010
	}
	if p.Copies.Hi > 1 {
		t |= 0x0100 | 0x0200 // COPIES, COLLATE
	}
	if slices.Contains(p.Finishings, 4) {
		t |= 0x0020
	}
	return t
}

// queued is the number of jobs not done. The printer is locked.
func (p *Printer) queued() int {
	n := 0
	for _, j := range p.jobs {
		if !p.state(j).Done() {
			n++
		}
	}
	return n
}

// versions are ipp-versions-supported.
func (p *Printer) versions() []ipp.Value {
	if p.Only11 {
		return []ipp.Value{ipp.Keyword("1.0"), ipp.Keyword("1.1")}
	}
	return []ipp.Value{ipp.Keyword("1.0"), ipp.Keyword("1.1"), ipp.Keyword("2.0")}
}

func first(s []string, def string) string {
	if len(s) == 0 || s[0] == "" {
		return def
	}
	return s[0]
}

func keywords(s []string) []ipp.Value {
	out := make([]ipp.Value, len(s))
	for i, v := range s {
		out[i] = ipp.Keyword(v)
	}
	return out
}

func mimes(s []string) []ipp.Value {
	out := make([]ipp.Value, len(s))
	for i, v := range s {
		out[i] = ipp.MimeType(v)
	}
	return out
}

func enums(ns []int) []ipp.Value {
	out := make([]ipp.Value, len(ns))
	for i, n := range ns {
		out[i] = ipp.Enum(n)
	}
	return out
}

func resolutions(rs []ipp.Resolution) []ipp.Value {
	out := make([]ipp.Value, len(rs))
	for i, r := range rs {
		out[i] = ipp.NewResolution(int(r.X), int(r.Y), r.Units)
	}
	return out
}
