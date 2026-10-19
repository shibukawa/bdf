package ipp

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

// sample is a message with a value of every type: a 1setOf, nested
// collections, an out-of-band value and a value of a tag the package
// does not know.
func sample() *Message {
	m := NewRequest(OpGetPrinterAttributes, 42, "", "ja")
	m.Add(TagOperation,
		Attr("printer-uri", URI("ipp://printer.local:631/ipp/print")),
		Attr("requesting-user-name", Name("yoshiki")),
		Attr("requested-attributes", Keyword("media-supported"), Keyword("sides-supported"), Keyword("printer-state")),
	)
	m.Add(TagJob,
		Attr("copies", Integer(3)),
		Attr("sides", Keyword("two-sided-long-edge")),
		Attr("job-name", NameLang("ja", "請求書")),
		Attr("job-note", TextLang("en", "note")),
		Attr("multiple-document-handling", Keyword("separate-documents-collated-copies")),
		Attr("print-quality", Enum(5)),
		Attr("ipp-attribute-fidelity", Boolean(true)),
		Attr("job-hold-until", OutOfBand(TagNoValue)),
		Attr("printer-resolution", NewResolution(600, 600, UnitsDPI)),
		Attr("page-ranges", NewRange(1, 3), NewRange(5, 5)),
		Attr("date-time-at-creation", DateTime(time.Date(2026, 10, 19, 12, 34, 56, 700000000, time.FixedZone("", 9*3600)))),
		Attr("octets", OctetString([]byte{0, 1, 2, 255})),
		Attr("media-col", Collection(
			Attr("media-size", Collection(Attr("x-dimension", Integer(21000)), Attr("y-dimension", Integer(29700)))),
			Attr("media-source", Keyword("tray-1")),
			Attr("media-type", Keyword("stationery"), Keyword("plain")),
		)),
		Attr("vendor", Value{Tag: 0x60, Raw: []byte{1, 2, 3}}),
		Attr("vendor-string", Value{Tag: 0x4b, Str: "unknown string"}),
		Attr("scheme", URIScheme("ipps")),
		Attr("mime", MimeType("application/pdf")),
		Attr("plain", Text("plain text")),
	)
	m.Add(TagUnsupportedGroup, Attr("finishings", OutOfBand(TagUnsupported)))
	return m
}

func TestRoundTrip(t *testing.T) {
	m := sample()
	b, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := Unmarshal(b, nil)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := got.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, b2) {
		t.Fatalf("the message does not survive a round trip:\n%x\n%x", b, b2)
	}
	if got.Version != Version20 || got.Op() != OpGetPrinterAttributes || got.RequestID != 42 {
		t.Fatalf("header: %+v", got)
	}
	if len(got.Groups) != 3 || got.Groups[0].Tag != TagOperation || got.Groups[1].Tag != TagJob || got.Groups[2].Tag != TagUnsupportedGroup {
		t.Fatalf("groups: %+v", got.Groups)
	}
	op := got.Group(TagOperation)
	if op.Get("attributes-charset").Text() != "utf-8" || op.Get("attributes-natural-language").Text() != "ja" {
		t.Fatalf("the first attributes: %+v", op.Attrs[:2])
	}
	if s := op.Get("requested-attributes").Strings(); len(s) != 3 || s[2] != "printer-state" {
		t.Fatalf("1setOf: %v", s)
	}
	job := got.Group(TagJob)
	if n, ok := job.Get("copies").Int(); !ok || n != 3 {
		t.Fatalf("copies: %d %v", n, ok)
	}
	if v, _ := job.Get("job-name").First(); v.Tag != TagNameLang || v.Lang != "ja" || v.Str != "請求書" {
		t.Fatalf("nameWithLanguage: %+v", v)
	}
	if b, ok := job.Get("ipp-attribute-fidelity").Bool(); !ok || !b {
		t.Fatalf("boolean: %v %v", b, ok)
	}
	if job.Get("job-hold-until").OutOfBand() != TagNoValue {
		t.Fatalf("out of band: %+v", job.Get("job-hold-until"))
	}
	if r := job.Get("printer-resolution").Resolutions(); len(r) != 1 || r[0] != (Resolution{600, 600, UnitsDPI}) {
		t.Fatalf("resolution: %v", r)
	}
	if r := job.Get("page-ranges").Ranges(); len(r) != 2 || r[0] != (Range{1, 3}) || r[1] != (Range{5, 5}) {
		t.Fatalf("ranges: %v", r)
	}
	if v, _ := job.Get("date-time-at-creation").First(); !v.Time.Equal(time.Date(2026, 10, 19, 3, 34, 56, 700000000, time.UTC)) {
		t.Fatalf("dateTime: %v", v.Time)
	} else if _, off := v.Time.Zone(); off != 9*3600 {
		t.Fatalf("dateTime zone: %d", off)
	}
	if v, _ := job.Get("octets").First(); v.Str != "\x00\x01\x02\xff" {
		t.Fatalf("octetString: %q", v.Str)
	}
	cols := job.Get("media-col").Collections()
	if len(cols) != 1 {
		t.Fatalf("collections: %d", len(cols))
	}
	size := Get(cols[0], "media-size").Collections()
	if len(size) != 1 {
		t.Fatalf("nested collection: %+v", cols[0])
	}
	if x, _ := Get(size[0], "x-dimension").Int(); x != 21000 {
		t.Fatalf("x-dimension: %d", x)
	}
	if s := Get(cols[0], "media-type").Strings(); len(s) != 2 || s[1] != "plain" {
		t.Fatalf("1setOf in a collection: %v", s)
	}
	if v, _ := job.Get("vendor").First(); v.Tag != 0x60 || !bytes.Equal(v.Raw, []byte{1, 2, 3}) {
		t.Fatalf("unknown tag: %+v", v)
	}
	if v, _ := job.Get("vendor-string").First(); v.Tag != 0x4b || v.Str != "unknown string" {
		t.Fatalf("unknown string tag: %+v", v)
	}
	if got.Group(TagUnsupportedGroup).Get("finishings").OutOfBand() != TagUnsupported {
		t.Fatalf("unsupported group: %+v", got.Groups[2])
	}
	if s := job.Get("media-col").String(); s != "{media-size={x-dimension=21000, y-dimension=29700}, media-source=tray-1, media-type=stationery,plain}" {
		t.Fatalf("String of a collection: %s", s)
	}
}

// TestBytes checks the encoding of a request against bytes assembled by
// hand from RFC 8010 §3.1.
func TestBytes(t *testing.T) {
	m := &Message{Version: Version20, Code: uint16(OpGetPrinterAttributes), RequestID: 1, Groups: []Group{{Tag: TagOperation, Attrs: []Attribute{
		Attr("attributes-charset", Charset("utf-8")),
		Attr("attributes-natural-language", Language("en")),
		Attr("printer-uri", URI("ipp://p/ipp/print")),
		Attr("requested-attributes", Keyword("a"), Keyword("bc")),
		Attr("limit", Integer(1)),
		Attr("my-jobs", Boolean(true)),
	}}}}
	want := []byte{
		2, 0, 0x00, 0x0b, 0, 0, 0, 1,
		0x01,
		0x47, 0, 18} // charset, name-length
	want = append(want, "attributes-charset"...)
	want = append(want, 0, 5)
	want = append(want, "utf-8"...)
	want = append(want, 0x48, 0, 27)
	want = append(want, "attributes-natural-language"...)
	want = append(want, 0, 2, 'e', 'n')
	want = append(want, 0x45, 0, 11)
	want = append(want, "printer-uri"...)
	want = append(want, 0, 17)
	want = append(want, "ipp://p/ipp/print"...)
	want = append(want, 0x44, 0, 20)
	want = append(want, "requested-attributes"...)
	want = append(want, 0, 1, 'a', 0x44, 0, 0, 0, 2, 'b', 'c')
	want = append(want, 0x21, 0, 5)
	want = append(want, "limit"...)
	want = append(want, 0, 4, 0, 0, 0, 1)
	want = append(want, 0x22, 0, 7)
	want = append(want, "my-jobs"...)
	want = append(want, 0, 1, 1, 0x03)
	got, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("encoding:\n got %x\nwant %x", got, want)
	}
	back, err := Unmarshal(want, nil)
	if err != nil {
		t.Fatal(err)
	}
	if back.Op() != OpGetPrinterAttributes || len(back.Groups) != 1 || len(back.Groups[0].Attrs) != 6 {
		t.Fatalf("decoding: %+v", back)
	}
}

func TestCollectionBytes(t *testing.T) {
	// RFC 8010 §3.1.6: a collection is begCollection, then memberAttrName
	// and values with empty names, then endCollection.
	m := &Message{Version: Version11, Code: 0, RequestID: 7, Groups: []Group{{Tag: TagJob, Attrs: []Attribute{
		Attr("media-col", Collection(Attr("media-type", Keyword("plain")))),
	}}}}
	want := []byte{1, 1, 0, 0, 0, 0, 0, 7, 0x02,
		0x34, 0, 9}
	want = append(want, "media-col"...)
	want = append(want, 0, 0,
		0x4a, 0, 0, 0, 10)
	want = append(want, "media-type"...)
	want = append(want, 0x44, 0, 0, 0, 5)
	want = append(want, "plain"...)
	want = append(want, 0x37, 0, 0, 0, 0, 0x03)
	got, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("encoding:\n got %x\nwant %x", got, want)
	}
}

func TestDecodeLeavesDocument(t *testing.T) {
	m := NewRequest(OpPrintJob, 1, "", "")
	b, _ := m.Marshal()
	doc := "%PDF-1.7 the document"
	r := bytes.NewReader(append(b, doc...))
	got, br, err := DecodeFrom(r, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Op() != OpPrintJob {
		t.Fatalf("op: %v", got.Op())
	}
	rest, err := io.ReadAll(br)
	if err != nil {
		t.Fatal(err)
	}
	if string(rest) != doc {
		t.Fatalf("the document after the message: %q", rest)
	}
}

func TestLimits(t *testing.T) {
	many := NewRequest(OpGetJobs, 1, "", "")
	for i := 0; i < 50; i++ {
		many.Add(TagJob, Attr("a", Integer(i)))
	}
	b, _ := many.Marshal()
	if _, err := Unmarshal(b, &Limits{MaxAttributes: 20}); !errors.Is(err, ErrLimit) {
		t.Fatalf("attributes: %v", err)
	}
	if _, err := Unmarshal(b, &Limits{MaxBytes: 100}); !errors.Is(err, ErrLimit) {
		t.Fatalf("bytes: %v", err)
	}
	long := NewRequest(OpGetJobs, 1, "", "")
	long.Add(TagJob, Attr("s", Text(strings.Repeat("x", 300))))
	b, _ = long.Marshal()
	if _, err := Unmarshal(b, &Limits{MaxValueBytes: 200}); !errors.Is(err, ErrLimit) {
		t.Fatalf("value bytes: %v", err)
	}
	vals := NewRequest(OpGetJobs, 1, "", "")
	a := Attr("k")
	for i := 0; i < 30; i++ {
		a.Values = append(a.Values, Keyword("v"))
	}
	vals.Add(TagJob, a)
	b, _ = vals.Marshal()
	if _, err := Unmarshal(b, &Limits{MaxValues: 10}); !errors.Is(err, ErrLimit) {
		t.Fatalf("values: %v", err)
	}
	deep := Collection(Attr("x", Integer(1)))
	for i := 0; i < 10; i++ {
		deep = Collection(Attr("c", deep))
	}
	nest := NewRequest(OpGetJobs, 1, "", "")
	nest.Add(TagJob, Attr("deep", deep))
	b, _ = nest.Marshal()
	if _, err := Unmarshal(b, &Limits{MaxDepth: 4}); !errors.Is(err, ErrLimit) {
		t.Fatalf("depth: %v", err)
	}
	if _, err := Unmarshal(b, &Limits{MaxDepth: 12}); err != nil {
		t.Fatalf("depth within the limit: %v", err)
	}
}

func TestMalformed(t *testing.T) {
	ok, _ := sample().Marshal()
	for name, b := range map[string][]byte{
		"empty":                     nil,
		"header only":               ok[:8],
		"truncated":                 ok[:len(ok)/2],
		"no end":                    ok[:len(ok)-1],
		"value before a group":      {2, 0, 0, 1, 0, 0, 0, 1, 0x21, 0, 1, 'a', 0, 4, 0, 0, 0, 1, 3},
		"value without a name":      {2, 0, 0, 1, 0, 0, 0, 1, 0x01, 0x21, 0, 0, 0, 4, 0, 0, 0, 1, 3},
		"integer of 2 bytes":        {2, 0, 0, 1, 0, 0, 0, 1, 0x01, 0x21, 0, 1, 'a', 0, 2, 0, 1, 3},
		"member outside collection": {2, 0, 0, 1, 0, 0, 0, 1, 0x01, 0x4a, 0, 0, 0, 1, 'a', 3},
		"collection not closed":     {2, 0, 0, 1, 0, 0, 0, 1, 0x01, 0x34, 0, 1, 'c', 0, 0, 0x4a, 0, 0, 0, 1, 'm', 0x21, 0, 0, 0, 4, 0, 0, 0, 1, 3},
		"text with language short":  {2, 0, 0, 1, 0, 0, 0, 1, 0x01, 0x35, 0, 1, 'a', 0, 3, 0, 2, 'e', 3},
	} {
		_, err := Unmarshal(b, nil)
		if !errors.Is(err, ErrMalformed) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

func TestEncodeErrors(t *testing.T) {
	for name, m := range map[string]*Message{
		"no name":      {Groups: []Group{{Tag: TagJob, Attrs: []Attribute{Attr("", Integer(1))}}}},
		"no value":     {Groups: []Group{{Tag: TagJob, Attrs: []Attribute{Attr("a")}}}},
		"group as tag": {Groups: []Group{{Tag: TagInteger}}},
		"member tag":   {Groups: []Group{{Tag: TagJob, Attrs: []Attribute{Attr("a", Value{Tag: TagMemberName})}}}},
		"long value":   {Groups: []Group{{Tag: TagJob, Attrs: []Attribute{Attr("a", Text(strings.Repeat("x", 70000)))}}}},
	} {
		if _, err := m.Marshal(); !errors.Is(err, ErrEncode) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

func TestNames(t *testing.T) {
	if OpPrintJob.String() != "Print-Job" || Op(0x1234).String() != "operation-0x1234" {
		t.Fatal(OpPrintJob, Op(0x1234))
	}
	if StatusClientErrorDocumentFormatNotSupported.String() != "client-error-document-format-not-supported" {
		t.Fatal(StatusClientErrorDocumentFormatNotSupported)
	}
	if !StatusOKIgnoredOrSubstituted.OK() || !StatusClientErrorNotFound.ClientError() || !StatusServerErrorBusy.ServerError() {
		t.Fatal("status classes")
	}
	if JobCompleted.String() != "completed" || !JobAborted.Done() || JobProcessing.Done() {
		t.Fatal("job states")
	}
	if PrinterStopped.String() != "stopped" || TagKeyword.String() != "keyword" || Tag(0x60).String() != "tag-0x60" {
		t.Fatal("names")
	}
	if NewResolution(300, 600, UnitsDPI).String() != "300x600dpi" || NewRange(1, 5).String() != "1-5" || OutOfBand(TagUnknown).String() != "unknown" {
		t.Fatal("value strings")
	}
}
