package print

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shibukawa/bdf/print/ipp"
)

// Meta is what a job says about itself besides how it is printed: the
// account it is charged to, what describes its document, and attributes
// of one's own. A ticket carries it under IPP's names; the client sends
// it to the printer with the job, and the IPP server keeps it with a job
// and tells it back (Get-Job-Attributes), as a printer does.
type Meta struct {
	// AccountID is job-account-id and AccountingUserID
	// job-accounting-user-id (PWG 5100.7): the account a job is charged
	// to (a department, a project, a cost centre) and the user of the
	// account. A printer that says it takes them
	// (job-account-id-supported) is asked them.
	AccountID        string `json:"job-account-id,omitempty"`
	AccountingUserID string `json:"job-accounting-user-id,omitempty"`
	// Metadata is document-metadata (PWG 5100.13): "key=value" strings
	// that describe the document, with the keys of Dublin Core ("title",
	// "creator", "identifier", "subject", "description", "date" …) or
	// keys of one's own that begin with "x-". It goes with the document
	// (Print-Job, Send-Document); a printer keeps it with the job.
	Metadata []string `json:"document-metadata,omitempty"`
	// Custom are job attributes of one's own, sent as they are: their
	// names begin with "x-" (RFC 8011 §7.2 leaves vendors their own
	// prefixes), their values are strings (name on the wire). A printer
	// that does not know one says so, and prints the job without it.
	Custom map[string]Strings `json:"bdf-attributes,omitempty"`
}

// Strings are the values of an attribute: in JSON a string, or a list
// of strings.
type Strings []string

// UnmarshalJSON reads a string or a list of strings.
func (s *Strings) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '"' {
		var v string
		if err := json.Unmarshal(b, &v); err != nil {
			return err
		}
		*s = Strings{v}
		return nil
	}
	var v []string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*s = v
	return nil
}

// MarshalJSON writes one value as a string, more as a list.
func (s Strings) MarshalJSON() ([]byte, error) {
	if len(s) == 1 {
		return json.Marshal(s[0])
	}
	return json.Marshal([]string(s))
}

// The bounds of the meta.
const (
	// maxMetaValues bounds the values of document-metadata and the
	// attributes of one's own.
	maxMetaValues = 64
	// maxMetaValue is the bytes of a value (octetString(MAX), text(MAX)).
	maxMetaValue = 1023
	// maxMetaName is the bytes of an account (name(MAX)) and of the name
	// of an attribute (a keyword).
	maxMetaName = 255
	// maxCustomValues bounds the values of an attribute of one's own.
	maxCustomValues = 16
)

// dublinCore are the keys of document-metadata that PWG 5100.13 names:
// the elements of the Dublin Core Metadata Element Set (RFC 5013) and
// the DCMI Metadata Terms.
var dublinCore = []string{
	"contributor", "coverage", "creator", "date", "description", "format", "identifier", "language", "publisher",
	"relation", "rights", "source", "subject", "title", "type",
	"abstract", "accessRights", "accrualMethod", "accrualPeriodicity", "accrualPolicy", "alternative", "audience",
	"available", "bibliographicCitation", "conformsTo", "created", "dateAccepted", "dateCopyrighted", "dateSubmitted",
	"educationLevel", "extent", "hasFormat", "hasPart", "hasVersion", "instructionalMethod", "isFormatOf", "isPartOf",
	"isReferencedBy", "isReplacedBy", "isRequiredBy", "issued", "isVersionOf", "license", "mediator", "medium",
	"modified", "provenance", "references", "replaces", "requires", "rightsHolder", "spatial", "tableOfContents",
	"temporal", "valid",
}

// Check validates the meta: the lengths, the keys of document-metadata
// and the names of the attributes of one's own.
func (m *Meta) Check() error { return m.check() }

// Empty tells a meta that says nothing.
func (m *Meta) Empty() bool {
	return m.AccountID == "" && m.AccountingUserID == "" && len(m.Metadata) == 0 && len(m.Custom) == 0
}

func (m *Meta) check() error {
	if err := checkAccount("job-account-id", m.AccountID); err != nil {
		return err
	}
	if err := checkAccount("job-accounting-user-id", m.AccountingUserID); err != nil {
		return err
	}
	if len(m.Metadata) > maxMetaValues {
		return fmt.Errorf("print: document-metadata: %d values, %d at most", len(m.Metadata), maxMetaValues)
	}
	for _, s := range m.Metadata {
		if err := checkMetadata(s); err != nil {
			return err
		}
	}
	if len(m.Custom) > maxMetaValues {
		return fmt.Errorf("print: %d attributes of one's own, %d at most", len(m.Custom), maxMetaValues)
	}
	for name, vs := range m.Custom {
		if err := checkCustom(name, vs); err != nil {
			return err
		}
	}
	return nil
}

func checkAccount(attr, s string) error {
	if s != "" && (len(s) > maxMetaName || !printable(s)) {
		return fmt.Errorf("print: %s %.60q: up to %d bytes of text", attr, s, maxMetaName)
	}
	return nil
}

// printable tells a string of UTF-8 without control characters.
func printable(s string) bool {
	if !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

// checkMetadata validates a value of document-metadata: a key of Dublin
// Core or one of one's own that begins with x-, "=", and a value.
func checkMetadata(s string) error {
	key, _, ok := strings.Cut(s, "=")
	if !ok || len(s) > maxMetaValue || !printable(s) || !slices.Contains(dublinCore, key) && !customName(key, true) {
		return fmt.Errorf("print: document-metadata %.60q: key=value, with a key of Dublin Core (title, creator, identifier …) or one that begins with x-", s)
	}
	return nil
}

// customName tells a name of one's own: "x-" and letters, digits, ".",
// "-" and "_"; lower-case letters for the name of an attribute, as
// keywords are, upper too in a key of document-metadata.
func customName(s string, upper bool) bool {
	if !strings.HasPrefix(s, "x-") || len(s) == 2 || len(s) > maxMetaName {
		return false
	}
	for _, c := range s[2:] {
		switch {
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9', c == '.', c == '-', c == '_':
		case upper && c >= 'A' && c <= 'Z':
		default:
			return false
		}
	}
	return true
}

// checkCustom validates an attribute of one's own.
func checkCustom(name string, vs Strings) error {
	if !customName(name, false) {
		return fmt.Errorf("print: attribute %.60q: a name of one's own begins with x-, in lower-case letters, digits, '.', '-' and '_'", name)
	}
	if len(vs) == 0 || len(vs) > maxCustomValues {
		return fmt.Errorf("print: attribute %s: 1 to %d values", name, maxCustomValues)
	}
	for _, v := range vs {
		if len(v) > maxMetaValue || !printable(v) {
			return fmt.Errorf("print: attribute %s: a value is up to %d bytes of text", name, maxMetaValue)
		}
	}
	return nil
}

// merge fills what the meta lacks from defaults.
func (m Meta) merge(d Meta) Meta {
	if m.AccountID == "" {
		m.AccountID = d.AccountID
	}
	if m.AccountingUserID == "" {
		m.AccountingUserID = d.AccountingUserID
	}
	if m.Metadata == nil {
		m.Metadata = slices.Clone(d.Metadata)
	}
	if len(d.Custom) > 0 {
		c := maps.Clone(m.Custom)
		if c == nil {
			c = map[string]Strings{}
		}
		for k, v := range d.Custom {
			if _, ok := c[k]; !ok {
				c[k] = slices.Clone(v)
			}
		}
		m.Custom = c
	}
	return m
}

// Attrs are the attributes a job carries of the meta: job-account-id,
// job-accounting-user-id and the attributes of one's own (job template
// attributes), and document-metadata (a job status attribute, as a
// printer keeps it). An empty meta has none.
func (m *Meta) Attrs() []ipp.Attribute {
	var out []ipp.Attribute
	if m.AccountID != "" {
		out = append(out, ipp.Attr("job-account-id", ipp.Name(m.AccountID)))
	}
	if m.AccountingUserID != "" {
		out = append(out, ipp.Attr("job-accounting-user-id", ipp.Name(m.AccountingUserID)))
	}
	for _, name := range slices.Sorted(maps.Keys(m.Custom)) {
		out = append(out, ipp.Attr(name, names(m.Custom[name])...))
	}
	if len(m.Metadata) > 0 {
		out = append(out, ipp.Attr("document-metadata", octets(m.Metadata)...))
	}
	return out
}

// opAttrs are the operation attributes of the meta: document-metadata,
// which goes with the document (Print-Job, Send-Document).
func (m *Meta) opAttrs() []ipp.Attribute {
	if len(m.Metadata) == 0 {
		return nil
	}
	return []ipp.Attribute{ipp.Attr("document-metadata", octets(m.Metadata)...)}
}

// jobAttrs adds the job template attributes of the meta for a printer:
// the account when the printer says it takes one, the attributes of
// one's own always (Validate-Job tells those it does not know).
func (m *Meta) jobAttrs(p *ipp.Printer, a *attrs) {
	if m.AccountID != "" {
		a.add("job-account-id", boolSupported(p, "job-account-id-supported"), ipp.Name(m.AccountID))
	}
	if m.AccountingUserID != "" {
		a.add("job-accounting-user-id", boolSupported(p, "job-accounting-user-id-supported"), ipp.Name(m.AccountingUserID))
	}
	for _, name := range slices.Sorted(maps.Keys(m.Custom)) {
		a.add(name, true, names(m.Custom[name])...)
	}
}

// boolSupported tells a printer attribute of a boolean that is true.
func boolSupported(p *ipp.Printer, name string) bool {
	a, ok := p.Has(name)
	if !ok {
		return false
	}
	b, _ := a.Bool()
	return b
}

// Take reads an attribute of a request that makes a job into the meta
// when it is one of the meta's: job-account-id, job-accounting-user-id,
// one named x-… (job template attributes) or document-metadata (an
// operation attribute). known tells it is one of them; ok that its value
// was taken, else the attribute is unsupported as it came.
func (m *Meta) Take(a ipp.Attribute) (known, ok bool) {
	switch {
	case a.Name == "job-account-id", a.Name == "job-accounting-user-id":
		v, _ := a.First()
		if len(a.Values) != 1 || !stringValue(v) || v.Str == "" || checkAccount(a.Name, v.Str) != nil {
			return true, false
		}
		if a.Name == "job-account-id" {
			m.AccountID = v.Str
		} else {
			m.AccountingUserID = v.Str
		}
		return true, true
	case a.Name == "document-metadata":
		if len(a.Values) == 0 || len(a.Values) > maxMetaValues {
			return true, false
		}
		vs := make([]string, len(a.Values))
		for i, v := range a.Values {
			if !stringValue(v) || checkMetadata(v.Str) != nil {
				return true, false
			}
			vs[i] = v.Str
		}
		m.Metadata = vs
		return true, true
	case strings.HasPrefix(a.Name, "x-"):
		if _, had := m.Custom[a.Name]; !had && len(m.Custom) >= maxMetaValues {
			return true, false
		}
		vs := make(Strings, len(a.Values))
		for i, v := range a.Values {
			if v.Tag == ipp.TagBeginCollection || v.Tag.IsOutOfBand() {
				return true, false
			}
			vs[i] = v.String()
		}
		if checkCustom(a.Name, vs) != nil {
			return true, false
		}
		if m.Custom == nil {
			m.Custom = map[string]Strings{}
		}
		m.Custom[a.Name] = vs
		return true, true
	}
	return false, false
}

// stringValue tells a value that is a string: a name, text, keyword,
// URI, octetString …
func stringValue(v ipp.Value) bool {
	return v.Tag.IsString() || v.Tag == ipp.TagTextLang || v.Tag == ipp.TagNameLang || v.Tag == ipp.TagOctetString
}

func names(vs Strings) []ipp.Value {
	out := make([]ipp.Value, len(vs))
	for i, v := range vs {
		out[i] = ipp.Name(v)
	}
	return out
}

func octets(vs []string) []ipp.Value {
	out := make([]ipp.Value, len(vs))
	for i, v := range vs {
		out[i] = ipp.OctetString([]byte(v))
	}
	return out
}
