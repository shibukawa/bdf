package ipp

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/docflow"
	"github.com/shibukawa/bdf/print"
	"github.com/shibukawa/bdf/print/ipp"
)

// A CUPS server node is a printer computers and phones add as AirPrint,
// IPP Everywhere and Mopria printers are added: its attributes tell what
// those clients ask before they print (the paper, the rasters, PDF), so
// that a Mac makes an AirPrint queue of it and sends PDF, not PostScript
// (which bdf does not read), Android's print services (Mopria) and
// Windows's IPP class driver take it as theirs, and CUPS makes an IPP
// Everywhere queue of it; it is advertised with DNS-SD on the local
// links unless the node says not to (Node.DNSSD false; cmd/docflow
// announces it). On its own port it is an IPP Everywhere printer, with
// Create-Job (created.go); not page-ranges, as clients pick the pages
// they send (a document goes on whole).
// What it receives is made .bdf and sent on, so what it can do is not a
// printer's: any of the paper below, in colour, with no margins; one
// side and one copy, as the document goes on as it is.

// The rasters clients draw for it: PWG raster's resolutions and types,
// Apple raster's (urf-supported: version 1.4, 8-bit grey and sRGB, the
// resolutions, one copy, one input slot, plain paper, normal quality),
// and PCLm's (its strips of 16 rows, flate preferred: the PDF converter
// reads any).
var (
	rasterResolutions = []int{150, 300, 600}
	rasterTypes       = []string{"black_1", "sgray_8", "srgb_8"}
	urfSupported      = []string{"V1.4", "W8", "SRGB24", "RS150-300-600", "CP1", "IS1", "MT1", "PQ4"}
)

// paper is the paper offered, A4 first; paperReady the paper a client starts
// with.
var (
	paper = []string{
		"iso_a4_210x297mm", "iso_a3_297x420mm", "iso_a5_148x210mm", "jis_b4_257x364mm", "jis_b5_182x257mm", "na_letter_8.5x11in",
		"na_legal_8.5x14in", "jpn_hagaki_100x148mm", "na_index-4x6_4x6in", "iso_dl_110x220mm", "jpn_chou3_120x235mm",
	}
	paperReady = []string{"iso_a4_210x297mm", "na_letter_8.5x11in"}
)

// paperMax is the class of the largest paper (printer-paper-max, the TXT
// key PaperMax): A3.
const paperMax = "tabloid-A3"

// txtPDL is the pdl key of the TXT record: PDF first, so that clients
// that take the first they can send vectors, then the rasters of
// AirPrint, IPP Everywhere and Mopria.
const txtPDL = "application/pdf,image/urf,image/pwg-raster,application/PCLm,image/jpeg,image/png"

// makeModel is the CUPS server nodes' printer-make-and-model: what a
// client names the driver it makes ("BDF Docflow-AirPrint" on a Mac).
const makeModel = "BDF Docflow"

// mopriaVersion is the version of Mopria the nodes answer as
// (mopria-certified, and the TXT key): Android's print services and
// Windows's IPP class driver take a printer that says it.
const mopriaVersion = "2.0"

// started is when the process started: the printers' configuration and
// state changed then, as far as a client is told.
var started = time.Now()

// airPrintAttrs are a CUPS server node's attributes as an AirPrint, IPP
// Everywhere and Mopria printer: q is the request (its destination, the
// host it came to), queued the jobs not ended, auth and security the
// keywords of how it is reached. On the manager's own server, which
// takes no Create-Job, it does not say IPP Everywhere
// (ipp-features-supported) and has no icons.
func (s *server) airPrintAttrs(q *ippRequest, queued int, auth, security string) []ipp.Attribute {
	ops := make([]ipp.Value, len(q.operations()))
	for i, op := range q.operations() {
		ops[i] = ipp.Enum(int(op))
	}
	var mediaDB, mediaReady, mediaSizes []ipp.Value
	for _, m := range paper {
		pm, _ := print.ParseMedia(m)
		w, h := int(math.Round(pm.Width/72*2540)), int(math.Round(pm.Height/72*2540))
		size := ipp.Collection(ipp.Attr("x-dimension", ipp.Integer(w)), ipp.Attr("y-dimension", ipp.Integer(h)))
		mediaSizes = append(mediaSizes, size)
		col := ipp.Collection(ipp.Attr("media-key", ipp.Keyword(m)), ipp.Attr("media-size", size), ipp.Attr("media-size-name", ipp.Keyword(m)),
			ipp.Attr("media-source", ipp.Keyword("auto")), ipp.Attr("media-type", ipp.Keyword("stationery")),
			ipp.Attr("media-top-margin", ipp.Integer(0)), ipp.Attr("media-bottom-margin", ipp.Integer(0)),
			ipp.Attr("media-left-margin", ipp.Integer(0)), ipp.Attr("media-right-margin", ipp.Integer(0)))
		mediaDB = append(mediaDB, col)
		for _, r := range paperReady {
			if r == m {
				mediaReady = append(mediaReady, col)
			}
		}
	}
	resolutions := make([]ipp.Value, len(rasterResolutions))
	for i, dpi := range rasterResolutions {
		resolutions[i] = ipp.NewResolution(dpi, dpi, ipp.UnitsDPI)
	}
	now := s.m.Now()
	features := []ipp.Value{ipp.Keyword("none")}
	var port []ipp.Attribute
	if q.server != "" {
		features = []ipp.Value{ipp.Keyword("ipp-everywhere")}
		var icons []ipp.Value
		for _, n := range ipp.IconSizes {
			icons = append(icons, ipp.URI(q.page()+"icons/"+strconv.Itoa(n)+".png"))
		}
		port = []ipp.Attribute{
			ipp.Attr("printer-icons", icons...),
			ipp.Attr("multiple-operation-time-out", ipp.Integer(int(createdTimeout/time.Second))),
			ipp.Attr("multiple-operation-time-out-action", ipp.Keyword("abort-job")),
		}
	}
	return append([]ipp.Attribute{
		ipp.Attr("printer-uri-supported", ipp.URI(q.printerURI())),
		ipp.Attr("uri-security-supported", ipp.Keyword(security)),
		ipp.Attr("uri-authentication-supported", ipp.Keyword(auth)),
		ipp.Attr("printer-name", ipp.Name(q.dest)),
		ipp.Attr("printer-info", ipp.Text(q.title)),
		ipp.Attr("printer-location", ipp.Text(q.title)),
		ipp.Attr("printer-make-and-model", ipp.Text(makeModel)),
		ipp.Attr("printer-more-info", ipp.URI(q.moreInfo())),
		ipp.Attr("printer-device-id", ipp.Text("MFG:BDF;MDL:Docflow;CMD:PDF,URF,PWG,PCLM,JPEG,PNG;")),
		ipp.Attr("printer-uuid", ipp.URI(printerUUID(q.dest))),
		ipp.Attr("printer-kind", ipp.Keyword("document")),
		ipp.Attr("printer-type", ipp.Integer(0)),
		ipp.Attr("printer-state", ipp.Enum(3)), // idle
		ipp.Attr("printer-state-reasons", ipp.Keyword("none")),
		ipp.Attr("printer-is-accepting-jobs", ipp.Boolean(true)),
		ipp.Attr("queued-job-count", ipp.Integer(queued)),
		ipp.Attr("printer-up-time", ipp.Integer(int(now.Unix()&0x7fffffff))),
		ipp.Attr("printer-current-time", ipp.DateTime(now)),
		ipp.Attr("printer-state-change-time", ipp.Integer(int(started.Unix()&0x7fffffff))),
		ipp.Attr("printer-state-change-date-time", ipp.DateTime(started)),
		ipp.Attr("printer-config-change-time", ipp.Integer(int(started.Unix()&0x7fffffff))),
		ipp.Attr("printer-config-change-date-time", ipp.DateTime(started)),
		ipp.Attr("printer-organization", ipp.Text("")),
		ipp.Attr("printer-organizational-unit", ipp.Text("")),
		ipp.Attr("printer-geo-location", ipp.OutOfBand(ipp.TagUnknown)),
		// no supplies: what is printed goes on as documents
		ipp.Attr("printer-supply", ipp.OctetString([]byte("index=1;class=unknown;type=other;unit=percent;maxcapacity=100;level=-2;colorantname=unknown;"))),
		ipp.Attr("printer-supply-description", ipp.Text("supplies of the printers the node sends to")),
		ipp.Attr("printer-supply-info-uri", ipp.URI(q.moreInfo())),
		ipp.Attr("mopria-certified", ipp.Text(mopriaVersion)),
		ipp.Attr("printer-firmware-name", ipp.Name("bdf")),
		ipp.Attr("printer-firmware-string-version", ipp.Text(makeModel)),
		ipp.Attr("printer-firmware-version", ipp.OctetString([]byte{1})),
		ipp.Attr("ipp-versions-supported", ipp.Keyword("1.1"), ipp.Keyword("2.0")),
		ipp.Attr("ipp-features-supported", features...),
		ipp.Attr("operations-supported", ops...),
		ipp.Attr("charset-configured", ipp.Charset("utf-8")),
		ipp.Attr("charset-supported", ipp.Charset("utf-8")),
		ipp.Attr("natural-language-configured", ipp.Language("en")),
		ipp.Attr("generated-natural-language-supported", ipp.Language("en")),
		ipp.Attr("document-format-default", ipp.MimeType("application/octet-stream")),
		ipp.Attr("document-format-supported", ipp.MimeTypes(ippFormats)...),
		ipp.Attr("document-format-preferred", ipp.MimeType("application/pdf")),
		ipp.Attr("pdf-versions-supported", ipp.Keywords([]string{"adobe-1.3", "adobe-1.4", "adobe-1.5", "adobe-1.6", "adobe-1.7", "iso-32000-1_2008", "iso-32000-2_2017"})...),
		ipp.Attr("compression-supported", ipp.Keyword("none")),
		ipp.Attr("pdl-override-supported", ipp.Keyword("not-attempted")),
		ipp.Attr("multiple-document-jobs-supported", ipp.Boolean(false)),
		ipp.Attr("job-ids-supported", ipp.Boolean(true)),
		ipp.Attr("preferred-attributes-supported", ipp.Boolean(false)),
		ipp.Attr("printer-get-attributes-supported", ipp.Keyword("document-format")),
		ipp.Attr("identify-actions-supported", ipp.Keyword("display")),
		ipp.Attr("identify-actions-default", ipp.Keyword("display")),
		ipp.Attr("which-jobs-supported", ipp.Keywords([]string{"completed", "not-completed", "all"})...),
		ipp.Attr("job-creation-attributes-supported", ipp.Keywords([]string{"copies", "job-account-id", "job-accounting-user-id", "media", "media-col",
			"orientation-requested", "print-color-mode", "print-quality", "printer-resolution", "sides"})...),
		ipp.Attr("job-k-octets-supported", ipp.NewRange(0, int(min(s.m.Options().MaxDocument/1024, math.MaxInt32)))),
		ipp.Attr("job-account-id-supported", ipp.Boolean(true)),
		ipp.Attr("job-account-id-default", ipp.OutOfBand(ipp.TagNoValue)),
		ipp.Attr("job-accounting-user-id-supported", ipp.Boolean(true)),
		ipp.Attr("job-accounting-user-id-default", ipp.OutOfBand(ipp.TagNoValue)),
		ipp.Attr("page-ranges-supported", ipp.Boolean(false)),
		ipp.Attr("copies-default", ipp.Integer(1)),
		ipp.Attr("copies-supported", ipp.NewRange(1, 1)),
		ipp.Attr("finishings-default", ipp.Enum(3)),
		ipp.Attr("finishings-supported", ipp.Enum(3)),
		ipp.Attr("output-bin-default", ipp.Keyword("face-up")),
		ipp.Attr("output-bin-supported", ipp.Keyword("face-up")),
		ipp.Attr("sides-default", ipp.Keyword("one-sided")),
		ipp.Attr("sides-supported", ipp.Keyword("one-sided")),
		ipp.Attr("orientation-requested-default", ipp.OutOfBand(ipp.TagNoValue)),
		ipp.Attr("orientation-requested-supported", ipp.Enums([]int{3, 4, 5, 6})...),
		ipp.Attr("landscape-orientation-requested-preferred", ipp.Enum(5)),
		ipp.Attr("print-quality-default", ipp.Enum(4)),
		ipp.Attr("print-quality-supported", ipp.Enum(4)),
		ipp.Attr("printer-resolution-default", ipp.NewResolution(300, 300, ipp.UnitsDPI)),
		ipp.Attr("printer-resolution-supported", resolutions...),
		ipp.Attr("print-color-mode-default", ipp.Keyword("auto")),
		ipp.Attr("print-color-mode-supported", ipp.Keyword("auto"), ipp.Keyword("color"), ipp.Keyword("monochrome")),
		ipp.Attr("print-content-optimize-default", ipp.Keyword("auto")),
		ipp.Attr("print-content-optimize-supported", ipp.Keyword("auto")),
		ipp.Attr("print-rendering-intent-default", ipp.Keyword("auto")),
		ipp.Attr("print-rendering-intent-supported", ipp.Keyword("auto")),
		ipp.Attr("color-supported", ipp.Boolean(true)),
		ipp.Attr("pages-per-minute", ipp.Integer(60)),
		ipp.Attr("pages-per-minute-color", ipp.Integer(60)),
		ipp.Attr("media-default", ipp.Keyword(paperReady[0])),
		ipp.Attr("media-supported", ipp.Keywords(paper)...),
		ipp.Attr("media-ready", ipp.Keywords(paperReady)...),
		ipp.Attr("media-col-supported", ipp.Keywords([]string{"media-bottom-margin", "media-left-margin", "media-right-margin", "media-size", "media-size-name",
			"media-source", "media-top-margin", "media-type"})...),
		ipp.Attr("media-col-database", mediaDB...),
		ipp.Attr("media-col-default", mediaReady[0]),
		ipp.Attr("media-col-ready", mediaReady...),
		ipp.Attr("media-size-supported", mediaSizes...),
		ipp.Attr("media-source-supported", ipp.Keyword("auto")),
		ipp.Attr("media-type-supported", ipp.Keyword("stationery")),
		ipp.Attr("media-top-margin-supported", ipp.Integer(0)),
		ipp.Attr("media-bottom-margin-supported", ipp.Integer(0)),
		ipp.Attr("media-left-margin-supported", ipp.Integer(0)),
		ipp.Attr("media-right-margin-supported", ipp.Integer(0)),
		ipp.Attr("printer-paper-max", ipp.Keyword(paperMax)),
		ipp.Attr("pwg-raster-document-resolution-supported", resolutions...),
		ipp.Attr("pwg-raster-document-type-supported", ipp.Keywords(rasterTypes)...),
		ipp.Attr("pwg-raster-document-sheet-back", ipp.Keyword("normal")),
		ipp.Attr("urf-supported", ipp.Keywords(urfSupported)...),
		ipp.Attr("pclm-source-resolution-supported", resolutions...),
		ipp.Attr("pclm-strip-height-preferred", ipp.Integer(16)),
		ipp.Attr("pclm-strip-height-supported", ipp.Integer(16)),
		ipp.Attr("pclm-compression-method-preferred", ipp.Keyword("flate")),
		ipp.Attr("pclm-raster-back-side", ipp.Keyword("normal")),
	}, port...)
}

// moreInfo is printer-more-info, the page a client opens for the
// printer: on a CUPS server node's port, what the port tells of the node
// (GET /); on the manager's own server, the topology.
func (q *ippRequest) moreInfo() string {
	if q.server != "" {
		return q.page()
	}
	return q.page() + "topology"
}

// page is the root of the server the request came to, over HTTP (https
// for ipps).
func (q *ippRequest) page() string {
	scheme := "http"
	if q.scheme() == "ipps" {
		scheme = "https"
	}
	return scheme + "://" + q.r.Host + "/"
}

// Advertisements are the DNS-SD services of the CUPS server nodes that
// are advertised (Node.Advertised) on a port of their own, of nodes (those
// whose port the process listens on): _ipps._tcp for one with a
// certificate, else _ipp._tcp, with the subtypes _print (IPP
// Everywhere) and _universal (AirPrint) and the TXT keys of Mopria, as
// "<title> @ <machine>" (machine is the computer's name, a label). The
// services are on the host the caller advertises them on
// (ipp.Advertise).
func Advertisements(m *docflow.Manager, nodes []docflow.Node, machine string) []ipp.Advertisement {
	air := "none"
	if len(m.Options().SubmitTokens) > 0 {
		air = "username,password"
	}
	var out []ipp.Advertisement
	for _, n := range nodes {
		if n.Kind != KindCups || !n.Advertised() {
			continue
		}
		title := or(n.Title, n.Name)
		txt := []string{
			"txtvers=1",
			"qtotal=1",
			"rp=printers/" + n.Name,
			"ty=" + makeModel,
			"product=(" + makeModel + ")",
			"usb_MFG=BDF",
			"usb_MDL=Docflow",
			"note=" + clipTXT(title),
			"priority=50",
			"pdl=" + txtPDL,
			"URF=" + strings.Join(urfSupported, ","),
			"UUID=" + strings.TrimPrefix(printerUUID(n.Name), "urn:uuid:"),
			"Color=T",
			"Duplex=F",
			"Copies=F",
			"Collate=F",
			"Scan=F",
			"Fax=F",
			"kind=document",
			"PaperMax=" + paperMax,
			"air=" + air,
			"mopria-certified=" + mopriaVersion,
			// Android's print service lists it (F opts out)
			"print_wfds=T",
		}
		typ := "_ipp._tcp"
		if n.TLSCert != "" {
			typ = "_ipps._tcp"
			txt = append(txt, "TLS=1.2")
		}
		out = append(out, ipp.Advertisement{Instance: instanceName(title, machine), Type: typ, Subtypes: []string{"_print", "_universal"}, Port: n.Port, TXT: txt})
	}
	return out
}

// instanceName is a node's DNS-SD instance name: its title and the
// computer's name, a label of 63 bytes at most.
func instanceName(title, machine string) string {
	name := strings.ReplaceAll(title+" @ "+machine, ".", " ")
	for len(name) > 63 {
		r := []rune(name)
		name = string(r[:len(r)-1])
	}
	return name
}

// clipTXT bounds a value of a TXT key: a string of the record has 255
// bytes at most.
func clipTXT(s string) string {
	if len(s) <= 200 {
		return s
	}
	r := []rune(s)
	for len(string(r)) > 200 {
		r = r[:len(r)-1]
	}
	return string(r)
}
