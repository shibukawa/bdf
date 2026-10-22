package form

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/shibukawa/bdf/form/tag"
)

// SchemaDialect is the dialect of JSON Schema that Schema writes.
const SchemaDialect = "https://json-schema.org/draft/2020-12/schema"

// SchemaKeyword is the keyword of the schemas of Schema that holds what a
// template says of a name: the kind of its tags ("value", "section" or
// "image"), their number format, how a section repeats ("mode", and
// "slots"), the symbology of a code ("code"), what a person gives for a
// picture ("input": "signature" or "seal"), where the tags are ("at"),
// what a person is shown for the values of a choice ("labels"), and that a
// text may have several lines ("multiline").
const SchemaKeyword = "x-bdf"

// Bound schema construction independently of how fields were obtained:
// compiled templates can provide fields without passing through tag.Parse.
const (
	maxSchemaDepth = 256
	maxSchemaNodes = 1 << 16
)

// Schema describes the data that a template asks for as a JSON Schema: an
// object with a property for each name that the tags of the template look
// up, in the order of the template. A program that makes a form for the
// data of a template (an editor, a page that asks a person for the values)
// makes it of the schema.
//
// A template does not say all of it, and the schema is what its tags
// suggest:
//
//   - A section with names in it, with {{.}}, or with a way to repeat (rows,
//     slots, pages) is a list: of objects with those names, or of values. A
//     section of slots holds as many items as it has slots, at most.
//   - Another section is a condition: true or false.
//   - A value with a number format is a number, with a date format a date
//     (a text as tag.ParseDate reads it), and otherwise a text or a number.
//   - A picture is a text: the path of a file, or a data: URL. A QR code or a
//     barcode is the text that the symbol says.
//   - A name with dots (customer.name) is a property of an object.
//
// A name in a section is looked up in the item of the section first and in
// what is around it next, so that the tags alone do not tell which of the
// two has it: it is a property of the item, unless the template looks the
// same name up around the section too. A sample of the data (that of a
// compiled template, or nil) settles it: the names are where the sample has
// them, a section is a list, an object or a condition as its value in the
// sample is, and a value has the type of its value there.
//
// Nothing is required: a tag whose name the data lacks is left empty.
// Construction is limited to 65,536 fields and nodes, nested at most 256 levels.
func Schema(fields []tag.Field, sample map[string]any) ([]byte, error) {
	root, err := schemaTree(fields, sample)
	if err != nil {
		return nil, err
	}
	return writeSchema(root.object())
}

// writeSchema writes the schema of the object of the data.
func writeSchema(out *ordered) ([]byte, error) {
	doc := &ordered{}
	doc.set("$schema", SchemaDialect)
	for i, k := range out.keys {
		doc.set(k, out.vals[i])
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// schemaTree places the fields of a template in the tree of the names of
// its data.
func schemaTree(fields []tag.Field, sample map[string]any) (*schemaNode, error) {
	if len(fields) > maxSchemaNodes {
		return nil, fmt.Errorf("form: the schema has more than %d fields", maxSchemaNodes)
	}
	root := &schemaNode{byName: map[string]*schemaNode{}, budget: &schemaBudget{nodes: 1}}
	if sample != nil {
		root.sample, root.sampled = tag.Normalize(sample), true
	}
	// the names around a section before those in it: the latter are looked
	// up in the former
	order := make([]int, len(fields))
	depth := make([][]string, len(fields))
	for i, fd := range fields {
		if len(fd.In) > maxSchemaDepth {
			return nil, fmt.Errorf("form: the schema is nested more than %d levels", maxSchemaDepth)
		}
		order[i] = i
		for _, s := range fd.In {
			if !strings.HasPrefix(s, "_") { // (a variable of a section: a condition, not a scope of data)
				depth[i] = append(depth[i], s)
			}
		}
	}
	sort.SliceStable(order, func(a, b int) bool { return len(depth[order[a]]) < len(depth[order[b]]) })
	for _, i := range order {
		if err := root.place(fields[i], depth[i], i); err != nil {
			return nil, err
		}
	}
	return root, nil
}

// SchemaOf returns the JSON Schema of the data of a template. A compiled
// template gives the schema it holds (that of its sample); with a sample,
// or for a template that is not compiled, the schema is made of the tags
// (Schema).
func SchemaOf(r io.ReaderAt, size int64, sample map[string]any) ([]byte, error) {
	_, _, _, def, err := source(r, size)
	if err != nil {
		return nil, err
	}
	if def != nil {
		if sample == nil && len(def.Schema) > 0 {
			var buf bytes.Buffer
			if err := json.Indent(&buf, def.Schema, "", "  "); err != nil {
				return nil, fmt.Errorf("form: the schema of the compiled template: %w", err)
			}
			buf.WriteByte('\n')
			return buf.Bytes(), nil
		}
		if sample == nil {
			sample = def.Sample
		}
		return Schema(def.Fields, sample)
	}
	fields, err := Fields(r, size)
	if err != nil {
		return nil, err
	}
	return Schema(fields, sample)
}

// schemaNode is a name of the data: a property of the object, or of the
// items of the list, that is around it.
type schemaNode struct {
	name   string
	kids   []*schemaNode
	byName map[string]*schemaNode
	depth  int
	budget *schemaBudget
	// first is the index of the first field of the name: the properties
	// are in the order of the template.
	first int

	value, image, section bool
	// sections counts the sections of the name that the template has.
	sections int
	// self is set for a section with {{.}} in it: its items are values.
	self     bool
	inverted bool
	code     string
	input    string
	format   string
	mode     string
	slots    int
	at       []string

	// what a template whose fields are more than tags says of the value
	typ       string
	options   []string
	labels    []string
	label     string
	def       any
	maxLen    int
	multiline bool

	// sample is the value of the name in the sample, when it has one.
	sample  any
	sampled bool
}

type schemaBudget struct{ nodes int }

func (n *schemaNode) kid(name string, field int) (*schemaNode, error) {
	if k := n.byName[name]; k != nil {
		return k, nil
	}
	if n.depth >= maxSchemaDepth {
		return nil, fmt.Errorf("form: the schema is nested more than %d levels", maxSchemaDepth)
	}
	if n.budget.nodes >= maxSchemaNodes {
		return nil, fmt.Errorf("form: the schema has more than %d nodes", maxSchemaNodes)
	}
	n.budget.nodes++
	k := &schemaNode{name: name, byName: map[string]*schemaNode{}, first: field, depth: n.depth + 1, budget: n.budget}
	if m, ok := n.item().(map[string]any); ok {
		if v, ok := m[name]; ok {
			k.sample, k.sampled = v, true
		}
	}
	n.byName[name] = k
	n.kids = append(n.kids, k)
	return k, nil
}

// item is what the names in a node are looked up in, in the sample: the
// value of an object, the first item of a list.
func (n *schemaNode) item() any {
	if l, ok := n.sample.([]any); ok {
		if len(l) == 0 {
			return nil
		}
		return l[0]
	}
	return n.sample
}

// has reports whether a name is one of a node: in the sample when the
// node has a value there, and otherwise when the template gave it one.
func (n *schemaNode) has(name string) bool {
	if m, ok := n.item().(map[string]any); ok {
		if _, ok := m[name]; ok {
			return true
		}
	}
	return n.byName[name] != nil
}

// find returns the node of a name with dots as it is looked up in a chain
// of scopes, the innermost last: in the innermost scope that has its first
// part, or in the innermost one when none has it. field is the index of
// the field that looks it up.
func find(chain []*schemaNode, name string, field int) (*schemaNode, error) {
	// Check before Split allocates an entry for every part of an untrusted name.
	if strings.Count(name, ".") >= maxSchemaDepth {
		return nil, fmt.Errorf("form: the schema is nested more than %d levels", maxSchemaDepth)
	}
	parts := strings.Split(name, ".")
	at := chain[len(chain)-1]
	for i := len(chain) - 1; i >= 0; i-- {
		if chain[i].has(parts[0]) {
			at = chain[i]
			break
		}
	}
	for _, p := range parts {
		if p == "" {
			continue
		}
		var err error
		at, err = at.kid(p, field)
		if err != nil {
			return nil, err
		}
	}
	return at, nil
}

// place gives the schema a field, the one at an index of the fields: in
// are the sections around its tags.
func (n *schemaNode) place(fd tag.Field, in []string, field int) error {
	chain := []*schemaNode{n}
	for _, s := range in {
		sec, err := find(chain, s, field)
		if err != nil {
			return err
		}
		sec.section = true
		chain = append(chain, sec)
	}
	if fd.Name == "." {
		if len(chain) > 1 {
			chain[len(chain)-1].self = true
		}
		return nil
	}
	if fd.Name == "" || strings.HasPrefix(fd.Name, "_") {
		return nil
	}
	leaf, err := find(chain, fd.Name, field)
	if err != nil {
		return err
	}
	if leaf == n {
		return nil
	}
	switch fd.Kind {
	case tag.FieldSection:
		// (inverted: no section of the name is another)
		leaf.inverted = fd.Inverted && (leaf.inverted || leaf.sections == 0)
		leaf.sections++
		leaf.section = true
		if leaf.mode == "" {
			leaf.mode, leaf.slots = fd.Mode, fd.Slots
		}
	case tag.FieldImage:
		leaf.image = true
		if leaf.code == "" {
			leaf.code = fd.Code
		}
		if leaf.input == "" {
			leaf.input = fd.Input
		}
	default:
		leaf.value = true
		if leaf.format == "" {
			leaf.format = fd.Format
		}
	}
	if leaf.typ == "" {
		leaf.typ = fd.Type
	}
	if leaf.options == nil {
		leaf.options, leaf.labels = fd.Options, fd.Labels
	}
	if leaf.label == "" {
		leaf.label = fd.Label
	}
	if leaf.def == nil {
		leaf.def = fd.Default
	}
	leaf.maxLen, leaf.multiline = max(leaf.maxLen, fd.MaxLen), leaf.multiline || fd.Multiline
	for _, a := range fd.At {
		if len(leaf.at) < 16 {
			leaf.at = append(leaf.at, a)
		}
	}
	return nil
}

// object writes the properties of a node.
func (n *schemaNode) object() *ordered {
	o := &ordered{}
	o.set("type", "object")
	props := &ordered{}
	kids := append([]*schemaNode(nil), n.kids...)
	sort.SliceStable(kids, func(a, b int) bool { return kids[a].first < kids[b].first })
	for _, k := range kids {
		props.set(k.name, k.schema())
	}
	o.set("properties", props)
	return o
}

// Types of the values of JSON.
var scalarTypes = []string{"string", "number"}

// schema writes the schema of a name.
func (n *schemaNode) schema() *ordered {
	info := &ordered{}
	kind := tag.FieldValue
	switch {
	case n.image:
		kind = tag.FieldImage
	case n.section:
		kind = tag.FieldSection
	}
	info.set("kind", kind)
	var o *ordered
	switch {
	case n.list():
		o = &ordered{}
		o.set("type", "array")
		switch {
		case len(n.kids) > 0:
			o.set("items", n.object())
		case n.self:
			item := &ordered{}
			item.set("type", n.scalar(n.item()))
			o.set("items", item)
		default:
			o.set("items", &ordered{})
		}
		if n.mode == tag.ModeSlots && n.slots > 0 {
			o.set("maxItems", n.slots)
		}
		if n.mode != "" {
			info.set("mode", n.mode)
		}
		if n.slots > 0 {
			info.set("slots", n.slots)
		}
	case len(n.kids) > 0:
		o = n.object()
	case n.section && !n.value && !n.image:
		// a condition: true or false, or the value that the sample has
		// for it (a text that is shown when there is one)
		o = &ordered{}
		if t := n.scalar(n.sample); t == "string" || t == "number" {
			o.set("type", t)
		} else {
			o.set("type", "boolean")
		}
		if n.inverted {
			o.set("description", "A part of the form is shown when this is false.")
			info.set("inverted", true)
		} else {
			o.set("description", "A part of the form is shown when this is true.")
		}
	case n.image && n.code == "":
		o = &ordered{}
		o.set("type", "string")
		o.set("format", "uri-reference")
		switch n.input {
		case tag.InputSignature:
			o.set("description", "A signature that a person draws: a picture, the path of a file or a data: URL.")
		case tag.InputSeal:
			o.set("description", "The picture of a seal: the path of a file, or a data: URL.")
		default:
			o.set("description", "A picture: the path of a file, or a data: URL.")
		}
		if n.input != "" {
			info.set("input", n.input)
		}
	case n.code != "":
		o = &ordered{}
		o.set("type", "string")
		o.set("description", "The text of the "+codeName(n.code)+" drawn in the form.")
		info.set("code", n.code)
	default:
		o = &ordered{}
		kind := n.typ
		if kind == "" {
			kind = formatKind(n.format)
		}
		switch kind {
		case "number", "boolean", "string":
			o.set("type", kind)
		case "date", "date-time":
			o.set("type", "string")
			o.set("format", kind)
		default:
			o.set("type", n.scalar(n.sample))
		}
		if len(n.options) > 0 {
			enum := make([]any, len(n.options))
			for i, v := range n.options {
				enum[i] = v
			}
			o.set("enum", enum)
			if len(n.labels) == len(n.options) {
				labels := make([]any, len(n.labels))
				for i, v := range n.labels {
					labels[i] = v
				}
				info.set("labels", labels)
			}
		}
		if n.maxLen > 0 {
			o.set("maxLength", n.maxLen)
		}
		if n.def != nil {
			o.set("default", n.def)
		}
		if n.multiline {
			info.set("multiline", true)
		}
		if n.format != "" {
			info.set("format", n.format)
		}
	}
	if n.label != "" {
		o.set("title", n.label)
	} else {
		o.set("title", n.name)
	}
	if len(n.at) > 0 {
		at := make([]any, len(n.at))
		for i, a := range n.at {
			at[i] = a
		}
		info.set("at", at)
	}
	o.set(SchemaKeyword, info)
	return o
}

// list reports whether a section is a list: by its value in the sample,
// and otherwise by what the template has in it.
func (n *schemaNode) list() bool {
	if !n.section {
		return false
	}
	if n.sampled && n.sample != nil {
		_, ok := n.sample.([]any)
		return ok
	}
	// (a section that the template also takes the value of is that value,
	// shown when there is one: {{#memo}}Memo: {{memo}}{{/memo}})
	return !n.value && !n.image && (len(n.kids) > 0 || n.self || n.mode != "")
}

// scalar is the type of a value: that of its value in the sample, or a
// text or a number.
func (n *schemaNode) scalar(sample any) any {
	switch sample.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	}
	out := make([]any, len(scalarTypes))
	for i, t := range scalarTypes {
		out[i] = t
	}
	return out
}

func codeName(code string) string {
	if code == "qr" {
		return "QR code"
	}
	return "barcode (" + code + ")"
}

// formatKind tells what a number format of Excel shows: "number", "date" or
// "date-time"; "" for the General format, a format of text or of a time of
// day, and what is none.
func formatKind(code string) string {
	if code == "" {
		return ""
	}
	// the first section, without what it shows as it is: quoted text,
	// escaped characters, colors and conditions in brackets, and the
	// characters that pad
	var b strings.Builder
	rs := []rune(code)
	for i := 0; i < len(rs); i++ {
		switch r := rs[i]; r {
		case '"':
			for i++; i < len(rs) && rs[i] != '"'; i++ {
			}
		case '\\', '_', '*':
			i++
		case '[':
			// ([h] and [mm] count hours and minutes; the others are colors,
			// conditions and locales)
			j := i + 1
			for j < len(rs) && rs[j] != ']' {
				j++
			}
			if in := strings.ToLower(string(rs[i+1 : min(j, len(rs))])); in != "" && strings.Trim(in, "hms") == "" {
				b.WriteString(in)
			}
			i = j
		case ';':
			i = len(rs)
		default:
			b.WriteRune(r)
		}
	}
	s := strings.ToLower(b.String())
	if strings.Contains(s, "general") || strings.Contains(s, "@") {
		return ""
	}
	digits := strings.ContainsAny(s, "0#?")
	if digits && (strings.Contains(s, "e+") || strings.Contains(s, "e-")) {
		return "number"
	}
	s = strings.NewReplacer("am/pm", "", "a/p", "").Replace(s)
	day := strings.ContainsAny(s, "ydge") || strings.Contains(s, "aaa") // (a weekday of Japanese formats)
	clock := strings.ContainsAny(s, "hs")
	if strings.Contains(s, "m") && !day && !clock {
		day = true // months alone
	}
	switch {
	case day && clock:
		return "date-time"
	case day:
		return "date"
	case clock:
		return "" // a time of day by itself: a number of a spreadsheet, a part of a day
	case digits:
		return "number"
	}
	return ""
}

// ordered is an object of JSON whose members keep their order.
type ordered struct {
	keys []string
	vals []any
}

func (o *ordered) set(k string, v any) {
	o.keys = append(o.keys, k)
	o.vals = append(o.vals, v)
}

func (o *ordered) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, k := range o.keys {
		if i > 0 {
			b.WriteByte(',')
		}
		if err := writeJSON(&b, k); err != nil {
			return nil, err
		}
		b.WriteByte(':')
		if err := writeJSON(&b, o.vals[i]); err != nil {
			return nil, err
		}
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func writeJSON(b *bytes.Buffer, v any) error {
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return err
	}
	b.Truncate(b.Len() - 1) // the line feed of Encode
	return nil
}
