// Package ipp reads and writes the messages of the Internet Printing
// Protocol (RFC 8010) and speaks it to printers (RFC 8011, PWG 5100.12:
// IPP/2.x). It is what bdf print and the print server talk to a printer,
// to a CUPS queue, and to each other with.
//
// A Message is a request or a response: a version, an operation or a
// status, a request id, and groups of attributes, each attribute a name
// with one or more values. Encode writes one and Decode reads one from
// the start of a body, leaving a document that follows it unread.
//
// A message from the network is as untrusted as any other input: Decode
// is bounded by Limits and refuses what passes them.
package ipp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"
)

// Version is the version of the protocol a message names.
type Version struct{ Major, Minor uint8 }

// Versions of the protocol.
var (
	Version11 = Version{1, 1}
	Version20 = Version{2, 0}
)

func (v Version) String() string {
	return strconv.Itoa(int(v.Major)) + "." + strconv.Itoa(int(v.Minor))
}

// Message is a request or a response.
type Message struct {
	Version Version
	// Code is the operation of a request (Op) or the status of a
	// response (Status).
	Code      uint16
	RequestID uint32
	Groups    []Group
}

// Op returns the code as the operation of a request.
func (m *Message) Op() Op { return Op(m.Code) }

// Status returns the code as the status of a response.
func (m *Message) Status() Status { return Status(m.Code) }

// Group is a group of attributes, begun by a group tag.
type Group struct {
	Tag   Tag
	Attrs []Attribute
}

// Attribute is a named value, or several (a 1setOf).
type Attribute struct {
	Name   string
	Values []Value
}

// Attr makes an attribute.
func Attr(name string, values ...Value) Attribute { return Attribute{Name: name, Values: values} }

// Resolution is a resolution value: dots per inch (UnitsDPI) or per
// centimetre (UnitsDPCM), across the feed and along it.
type Resolution struct {
	X, Y  int32
	Units uint8
}

// Units of a Resolution.
const (
	UnitsDPI  uint8 = 3
	UnitsDPCM uint8 = 4
)

func (r Resolution) String() string {
	u := "dpi"
	if r.Units == UnitsDPCM {
		u = "dpcm"
	} else if r.Units != UnitsDPI {
		u = "units-" + strconv.Itoa(int(r.Units))
	}
	return strconv.Itoa(int(r.X)) + "x" + strconv.Itoa(int(r.Y)) + u
}

// Range is a rangeOfInteger value.
type Range struct{ Lo, Hi int32 }

func (r Range) String() string { return strconv.Itoa(int(r.Lo)) + "-" + strconv.Itoa(int(r.Hi)) }

// Value is one value of an attribute. Tag tells its type and which of
// the fields holds it: Int for integer and enum (and boolean, as 0 or
// 1); Str for the character strings, octetString and the text of
// textWithLanguage and nameWithLanguage (whose language is Lang); Time
// for dateTime; Res for resolution; Range for rangeOfInteger; Coll for
// a collection; none for an out-of-band value. Raw holds the bytes of a
// value whose tag the package does not know, as they were.
type Value struct {
	Tag   Tag
	Int   int32
	Str   string
	Lang  string
	Time  time.Time
	Res   Resolution
	Range Range
	Coll  []Attribute
	Raw   []byte
}

// Integer makes an integer value.
func Integer(n int) Value { return Value{Tag: TagInteger, Int: int32(n)} }

// Boolean makes a boolean value.
func Boolean(b bool) Value {
	v := Value{Tag: TagBoolean}
	if b {
		v.Int = 1
	}
	return v
}

// Enum makes an enum value.
func Enum(n int) Value { return Value{Tag: TagEnum, Int: int32(n)} }

// Keyword makes a keyword value.
func Keyword(s string) Value { return Value{Tag: TagKeyword, Str: s} }

// Name makes a nameWithoutLanguage value.
func Name(s string) Value { return Value{Tag: TagName, Str: s} }

// Text makes a textWithoutLanguage value.
func Text(s string) Value { return Value{Tag: TagText, Str: s} }

// URI makes a uri value.
func URI(s string) Value { return Value{Tag: TagURI, Str: s} }

// URIScheme makes a uriScheme value.
func URIScheme(s string) Value { return Value{Tag: TagURIScheme, Str: s} }

// Charset makes a charset value.
func Charset(s string) Value { return Value{Tag: TagCharset, Str: s} }

// Language makes a naturalLanguage value.
func Language(s string) Value { return Value{Tag: TagLanguage, Str: s} }

// MimeType makes a mimeMediaType value.
func MimeType(s string) Value { return Value{Tag: TagMimeType, Str: s} }

// OctetString makes an octetString value.
func OctetString(b []byte) Value { return Value{Tag: TagOctetString, Str: string(b)} }

// TextLang makes a textWithLanguage value.
func TextLang(lang, s string) Value { return Value{Tag: TagTextLang, Lang: lang, Str: s} }

// NameLang makes a nameWithLanguage value.
func NameLang(lang, s string) Value { return Value{Tag: TagNameLang, Lang: lang, Str: s} }

// DateTime makes a dateTime value.
func DateTime(t time.Time) Value { return Value{Tag: TagDateTime, Time: t} }

// NewResolution makes a resolution value.
func NewResolution(x, y int, units uint8) Value {
	return Value{Tag: TagResolution, Res: Resolution{X: int32(x), Y: int32(y), Units: units}}
}

// NewRange makes a rangeOfInteger value.
func NewRange(lo, hi int) Value {
	return Value{Tag: TagRange, Range: Range{Lo: int32(lo), Hi: int32(hi)}}
}

// Collection makes a collection value of member attributes.
func Collection(members ...Attribute) Value { return Value{Tag: TagBeginCollection, Coll: members} }

// OutOfBand makes a value that is not there (TagUnsupported, TagUnknown,
// TagNoValue …).
func OutOfBand(tag Tag) Value { return Value{Tag: tag} }

// Bool returns a boolean value, or whether an integer is not zero.
func (v Value) Bool() bool { return v.Int != 0 }

// String renders the value for people: a number, a keyword, a time, a
// resolution, a range, a collection in braces, the tag of an out-of-band
// value, or the hex of unknown bytes.
func (v Value) String() string {
	switch {
	case v.Tag == TagBoolean:
		return strconv.FormatBool(v.Bool())
	case v.Tag == TagInteger || v.Tag == TagEnum:
		return strconv.Itoa(int(v.Int))
	case v.Tag == TagDateTime:
		return v.Time.Format(time.RFC3339)
	case v.Tag == TagResolution:
		return v.Res.String()
	case v.Tag == TagRange:
		return v.Range.String()
	case v.Tag == TagBeginCollection:
		var b strings.Builder
		b.WriteByte('{')
		for i, a := range v.Coll {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(a.Name)
			b.WriteByte('=')
			b.WriteString(a.String())
		}
		b.WriteByte('}')
		return b.String()
	case v.Tag.IsOutOfBand():
		return v.Tag.String()
	case v.Raw != nil:
		return v.Tag.String() + ":" + hex.EncodeToString(v.Raw)
	case v.Tag == TagTextLang || v.Tag == TagNameLang:
		return v.Str
	case v.Tag == TagOctetString:
		return hex.EncodeToString([]byte(v.Str))
	}
	return v.Str
}

// String renders the values of the attribute, separated by commas.
func (a Attribute) String() string {
	if len(a.Values) == 1 {
		return a.Values[0].String()
	}
	parts := make([]string, len(a.Values))
	for i, v := range a.Values {
		parts[i] = v.String()
	}
	return strings.Join(parts, ",")
}

// NewRequest makes a request of the operation, with the operation
// group begun by the charset and the language that RFC 8011 requires
// first ("utf-8" and "en" when "").
func NewRequest(op Op, id uint32, charset, lang string) *Message {
	if charset == "" {
		charset = "utf-8"
	}
	if lang == "" {
		lang = "en"
	}
	return &Message{Version: Version20, Code: uint16(op), RequestID: id, Groups: []Group{{Tag: TagOperation, Attrs: []Attribute{
		Attr("attributes-charset", Charset(charset)), Attr("attributes-natural-language", Language(lang))}}}}
}

// NewResponse makes a response of the status to a request, with the
// same first two attributes.
func NewResponse(st Status, id uint32, charset, lang string) *Message {
	m := NewRequest(Op(st), id, charset, lang)
	return m
}

// Add adds attributes to the first group of the tag, making the group
// when the message has none.
func (m *Message) Add(tag Tag, attrs ...Attribute) {
	for i := range m.Groups {
		if m.Groups[i].Tag == tag {
			m.Groups[i].Attrs = append(m.Groups[i].Attrs, attrs...)
			return
		}
	}
	m.Groups = append(m.Groups, Group{Tag: tag, Attrs: attrs})
}

// Encoding.

// maxLength is the largest name or value the encoding can carry (a
// 16-bit length).
const maxLength = math.MaxUint16

// ErrEncode is wrapped by the errors of encoding a message that the
// encoding cannot carry.
var ErrEncode = errors.New("ipp: the message cannot be encoded")

// Marshal returns the encoding of the message.
func (m *Message) Marshal() ([]byte, error) {
	return m.AppendTo(nil)
}

// Encode writes the encoding of the message.
func (m *Message) Encode(w io.Writer) error {
	b, err := m.Marshal()
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// AppendTo appends the encoding of the message to buf.
func (m *Message) AppendTo(buf []byte) ([]byte, error) {
	buf = append(buf, m.Version.Major, m.Version.Minor)
	buf = binary.BigEndian.AppendUint16(buf, m.Code)
	buf = binary.BigEndian.AppendUint32(buf, m.RequestID)
	var err error
	for gi := range m.Groups {
		g := &m.Groups[gi]
		if !g.Tag.IsGroup() || g.Tag == TagEnd {
			return nil, fmt.Errorf("%w: %s is not a group tag", ErrEncode, g.Tag)
		}
		buf = append(buf, byte(g.Tag))
		for ai := range g.Attrs {
			if buf, err = appendAttribute(buf, &g.Attrs[ai], 0); err != nil {
				return nil, err
			}
		}
	}
	return append(buf, byte(TagEnd)), nil
}

// appendAttribute appends an attribute: its first value with the name,
// the others with an empty name. A member of a collection (depth above
// 0) is written as memberAttrName and values without a name.
func appendAttribute(buf []byte, a *Attribute, depth int) ([]byte, error) {
	if len(a.Name) > maxLength || (depth == 0 && a.Name == "") {
		return nil, fmt.Errorf("%w: attribute name %q", ErrEncode, clip(a.Name))
	}
	if len(a.Values) == 0 {
		return nil, fmt.Errorf("%w: attribute %q has no value", ErrEncode, a.Name)
	}
	if depth > 0 {
		buf = append(buf, byte(TagMemberName), 0, 0)
		buf = appendLength(buf, len(a.Name))
		buf = append(buf, a.Name...)
	}
	var err error
	for i := range a.Values {
		v := &a.Values[i]
		if v.Tag.IsGroup() || v.Tag == TagMemberName || v.Tag == TagEndCollection {
			return nil, fmt.Errorf("%w: attribute %q has a value tagged %s", ErrEncode, a.Name, v.Tag)
		}
		buf = append(buf, byte(v.Tag))
		if i == 0 && depth == 0 {
			buf = appendLength(buf, len(a.Name))
			buf = append(buf, a.Name...)
		} else {
			buf = append(buf, 0, 0)
		}
		if buf, err = appendValue(buf, v, a.Name, depth); err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func appendLength(buf []byte, n int) []byte { return binary.BigEndian.AppendUint16(buf, uint16(n)) }

func appendValue(buf []byte, v *Value, name string, depth int) ([]byte, error) {
	if v.Raw != nil && v.Tag != TagBeginCollection {
		if len(v.Raw) > maxLength {
			return nil, fmt.Errorf("%w: a value of %q is %d bytes", ErrEncode, name, len(v.Raw))
		}
		buf = appendLength(buf, len(v.Raw))
		return append(buf, v.Raw...), nil
	}
	switch {
	case v.Tag.IsOutOfBand():
		return append(buf, 0, 0), nil
	case v.Tag == TagInteger || v.Tag == TagEnum:
		buf = append(buf, 0, 4)
		return binary.BigEndian.AppendUint32(buf, uint32(v.Int)), nil
	case v.Tag == TagBoolean:
		b := byte(0)
		if v.Int != 0 {
			b = 1
		}
		return append(buf, 0, 1, b), nil
	case v.Tag == TagDateTime:
		t := v.Time
		y := t.Year()
		if y < 0 || y > math.MaxUint16 {
			return nil, fmt.Errorf("%w: the year %d of %q", ErrEncode, y, name)
		}
		_, off := t.Zone()
		sign := byte('+')
		if off < 0 {
			sign, off = '-', -off
		}
		buf = append(buf, 0, 11)
		buf = binary.BigEndian.AppendUint16(buf, uint16(y))
		return append(buf, byte(t.Month()), byte(t.Day()), byte(t.Hour()), byte(t.Minute()), byte(t.Second()),
			byte(t.Nanosecond()/1e8), sign, byte(off/3600), byte(off%3600/60)), nil
	case v.Tag == TagResolution:
		buf = append(buf, 0, 9)
		buf = binary.BigEndian.AppendUint32(buf, uint32(v.Res.X))
		buf = binary.BigEndian.AppendUint32(buf, uint32(v.Res.Y))
		return append(buf, v.Res.Units), nil
	case v.Tag == TagRange:
		buf = append(buf, 0, 8)
		buf = binary.BigEndian.AppendUint32(buf, uint32(v.Range.Lo))
		return binary.BigEndian.AppendUint32(buf, uint32(v.Range.Hi)), nil
	case v.Tag == TagTextLang || v.Tag == TagNameLang:
		n := 2 + len(v.Lang) + 2 + len(v.Str)
		if len(v.Lang) > maxLength || len(v.Str) > maxLength || n > maxLength {
			return nil, fmt.Errorf("%w: a value of %q is %d bytes", ErrEncode, name, n)
		}
		buf = appendLength(buf, n)
		buf = appendLength(buf, len(v.Lang))
		buf = append(buf, v.Lang...)
		buf = appendLength(buf, len(v.Str))
		return append(buf, v.Str...), nil
	case v.Tag == TagBeginCollection:
		buf = append(buf, 0, 0)
		var err error
		for i := range v.Coll {
			if buf, err = appendAttribute(buf, &v.Coll[i], depth+1); err != nil {
				return nil, err
			}
		}
		return append(buf, byte(TagEndCollection), 0, 0, 0, 0), nil
	default:
		// The character strings, octetString, and unknown tags whose
		// value is in Str.
		if len(v.Str) > maxLength {
			return nil, fmt.Errorf("%w: a value of %q is %d bytes", ErrEncode, name, len(v.Str))
		}
		buf = appendLength(buf, len(v.Str))
		return append(buf, v.Str...), nil
	}
}

// Decoding.

// Limits bound what Decode accepts of a message: its bytes, its
// attributes (the members of collections included), the bytes of one
// value, the values of one attribute, and the nesting of collections.
// A zero field takes the default.
type Limits struct {
	MaxBytes      int
	MaxAttributes int
	MaxValueBytes int
	MaxValues     int
	MaxDepth      int
}

// DefaultLimits are the limits of Decode when none are given: 4 MiB,
// 10,000 attributes, 64 KiB a value, 1,000 values an attribute, and
// collections 8 deep.
var DefaultLimits = Limits{MaxBytes: 4 << 20, MaxAttributes: 10000, MaxValueBytes: 64 << 10, MaxValues: 1000, MaxDepth: 8}

func (l Limits) orDefault() Limits {
	d := DefaultLimits
	if l.MaxBytes > 0 {
		d.MaxBytes = l.MaxBytes
	}
	if l.MaxAttributes > 0 {
		d.MaxAttributes = l.MaxAttributes
	}
	if l.MaxValueBytes > 0 {
		d.MaxValueBytes = l.MaxValueBytes
	}
	if l.MaxValues > 0 {
		d.MaxValues = l.MaxValues
	}
	if l.MaxDepth > 0 {
		d.MaxDepth = l.MaxDepth
	}
	return d
}

// ErrLimit is wrapped by the error of a message that passes a limit.
var ErrLimit = errors.New("ipp: the message passes a limit")

// ErrMalformed is wrapped by the error of a message that is not an IPP
// message.
var ErrMalformed = errors.New("ipp: malformed message")

// Unmarshal decodes a message from bytes. Anything after the message is
// ignored.
func Unmarshal(b []byte, lim *Limits) (*Message, error) {
	return Decode(bytes.NewReader(b), lim)
}

// Decode reads a message from the start of r: its header, its groups,
// and the end-of-attributes tag, leaving what follows (the document of
// a Print-Job or Send-Document) unread. r is read with a bufio.Reader
// unless it is one, so the caller that reads on should read from the
// returned reader.
func Decode(r io.Reader, lim *Limits) (*Message, error) {
	m, _, err := DecodeFrom(r, lim)
	return m, err
}

// DecodeFrom is Decode that also returns the buffered reader positioned
// after the message, for the document that follows it.
func DecodeFrom(r io.Reader, lim *Limits) (*Message, *bufio.Reader, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}
	var l Limits
	if lim != nil {
		l = *lim
	}
	d := &decoder{r: br, lim: l.orDefault()}
	m, err := d.message()
	if err != nil {
		return nil, br, err
	}
	return m, br, nil
}

type decoder struct {
	r     *bufio.Reader
	lim   Limits
	n     int // bytes read
	attrs int // attributes read, members included
	buf   []byte
}

// token is one tag of the stream with what follows it.
type token struct {
	tag     Tag
	hasName bool
	name    string
	data    []byte // the value; valid until the next token
}

func (d *decoder) read(n int) ([]byte, error) {
	if d.n+n > d.lim.MaxBytes {
		return nil, fmt.Errorf("%w: more than %d bytes", ErrLimit, d.lim.MaxBytes)
	}
	if cap(d.buf) < n {
		d.buf = make([]byte, n, n+256)
	}
	d.buf = d.buf[:n]
	if _, err := io.ReadFull(d.r, d.buf); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, fmt.Errorf("%w: the message ends at byte %d", ErrMalformed, d.n)
		}
		return nil, err
	}
	d.n += n
	return d.buf, nil
}

func (d *decoder) uint16() (int, error) {
	b, err := d.read(2)
	if err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint16(b)), nil
}

// next reads a tag and, for a value tag, its name and value.
func (d *decoder) next() (token, error) {
	b, err := d.read(1)
	if err != nil {
		return token{}, err
	}
	t := token{tag: Tag(b[0])}
	if t.tag.IsGroup() {
		return t, nil
	}
	n, err := d.uint16()
	if err != nil {
		return token{}, err
	}
	if n > 0 {
		b, err := d.read(n)
		if err != nil {
			return token{}, err
		}
		t.hasName, t.name = true, string(b)
	}
	if n, err = d.uint16(); err != nil {
		return token{}, err
	}
	if n > d.lim.MaxValueBytes {
		return token{}, fmt.Errorf("%w: a value of %d bytes", ErrLimit, n)
	}
	if t.data, err = d.read(n); err != nil {
		return token{}, err
	}
	return t, nil
}

func (d *decoder) message() (*Message, error) {
	h, err := d.read(8)
	if err != nil {
		return nil, err
	}
	m := &Message{Version: Version{h[0], h[1]}, Code: binary.BigEndian.Uint16(h[2:]), RequestID: binary.BigEndian.Uint32(h[4:])}
	var g *Group
	var a *Attribute
	for {
		t, err := d.next()
		if err != nil {
			return nil, err
		}
		switch {
		case t.tag == TagEnd:
			return m, nil
		case t.tag.IsGroup():
			m.Groups = append(m.Groups, Group{Tag: t.tag})
			g, a = &m.Groups[len(m.Groups)-1], nil
			continue
		case g == nil:
			return nil, fmt.Errorf("%w: a value before any group", ErrMalformed)
		case t.tag == TagMemberName || t.tag == TagEndCollection:
			return nil, fmt.Errorf("%w: %s outside a collection", ErrMalformed, t.tag)
		case t.hasName:
			if err := d.countAttr(); err != nil {
				return nil, err
			}
			g.Attrs = append(g.Attrs, Attribute{Name: t.name})
			a = &g.Attrs[len(g.Attrs)-1]
		case a == nil:
			return nil, fmt.Errorf("%w: a value without a name", ErrMalformed)
		}
		if err := d.addValue(a, t, 0); err != nil {
			return nil, err
		}
	}
}

func (d *decoder) countAttr() error {
	d.attrs++
	if d.attrs > d.lim.MaxAttributes {
		return fmt.Errorf("%w: more than %d attributes", ErrLimit, d.lim.MaxAttributes)
	}
	return nil
}

// addValue parses the value of a token and adds it to the attribute.
func (d *decoder) addValue(a *Attribute, t token, depth int) error {
	if len(a.Values) >= d.lim.MaxValues {
		return fmt.Errorf("%w: attribute %q has more than %d values", ErrLimit, a.Name, d.lim.MaxValues)
	}
	v, err := d.value(t, a.Name, depth)
	if err != nil {
		return err
	}
	a.Values = append(a.Values, v)
	return nil
}

func (d *decoder) value(t token, name string, depth int) (Value, error) {
	v := Value{Tag: t.tag}
	b := t.data
	bad := func() (Value, error) {
		return Value{}, fmt.Errorf("%w: attribute %q: a %s value of %d bytes", ErrMalformed, name, t.tag, len(b))
	}
	switch {
	case t.tag.IsOutOfBand():
		if len(b) > 0 {
			v.Raw = bytes.Clone(b)
		}
	case t.tag == TagInteger || t.tag == TagEnum:
		if len(b) != 4 {
			return bad()
		}
		v.Int = int32(binary.BigEndian.Uint32(b))
	case t.tag == TagBoolean:
		if len(b) != 1 {
			return bad()
		}
		if b[0] != 0 {
			v.Int = 1
		}
	case t.tag == TagDateTime:
		if len(b) != 11 {
			return bad()
		}
		off := (int(b[9])*60 + int(b[10])) * 60
		if b[8] == '-' {
			off = -off
		}
		v.Time = time.Date(int(binary.BigEndian.Uint16(b)), time.Month(b[2]), int(b[3]), int(b[4]), int(b[5]), int(b[6]),
			int(b[7])*1e8, time.FixedZone("", off))
	case t.tag == TagResolution:
		if len(b) != 9 {
			return bad()
		}
		v.Res = Resolution{X: int32(binary.BigEndian.Uint32(b)), Y: int32(binary.BigEndian.Uint32(b[4:])), Units: b[8]}
	case t.tag == TagRange:
		if len(b) != 8 {
			return bad()
		}
		v.Range = Range{Lo: int32(binary.BigEndian.Uint32(b)), Hi: int32(binary.BigEndian.Uint32(b[4:]))}
	case t.tag == TagTextLang || t.tag == TagNameLang:
		if len(b) < 4 {
			return bad()
		}
		n := int(binary.BigEndian.Uint16(b))
		if 2+n+2 > len(b) {
			return bad()
		}
		v.Lang = string(b[2 : 2+n])
		rest := b[2+n:]
		k := int(binary.BigEndian.Uint16(rest))
		if 2+k != len(rest) {
			return bad()
		}
		v.Str = string(rest[2:])
	case t.tag == TagBeginCollection:
		if depth+1 > d.lim.MaxDepth {
			return Value{}, fmt.Errorf("%w: collections more than %d deep", ErrLimit, d.lim.MaxDepth)
		}
		coll, err := d.collection(name, depth+1)
		if err != nil {
			return Value{}, err
		}
		v.Coll = coll
	case t.tag.IsString() || t.tag == TagOctetString:
		v.Str = string(b)
	default:
		v.Raw = bytes.Clone(b)
		if v.Raw == nil {
			v.Raw = []byte{}
		}
	}
	return v, nil
}

// collection reads the members of a collection after its begCollection
// value, through its endCollection.
func (d *decoder) collection(name string, depth int) ([]Attribute, error) {
	var attrs []Attribute
	var a *Attribute
	for {
		t, err := d.next()
		if err != nil {
			return nil, err
		}
		switch {
		case t.tag == TagEndCollection:
			return attrs, nil
		case t.tag == TagEnd || t.tag.IsGroup():
			return nil, fmt.Errorf("%w: the collection of %q is not closed", ErrMalformed, name)
		case t.tag == TagMemberName:
			if err := d.countAttr(); err != nil {
				return nil, err
			}
			attrs = append(attrs, Attribute{Name: string(t.data)})
			a = &attrs[len(attrs)-1]
			continue
		case a == nil:
			return nil, fmt.Errorf("%w: a value before any member of the collection of %q", ErrMalformed, name)
		}
		if err := d.addValue(a, t, depth); err != nil {
			return nil, err
		}
	}
}

func clip(s string) string {
	if len(s) > 32 {
		return s[:32] + "…"
	}
	return s
}
