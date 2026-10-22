package form

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/shibukawa/bdf"
	"github.com/shibukawa/bdf/form/tag"
)

// OpenAll in Options.Open leaves open every name that the data does not
// give.
const OpenAll = "*"

// A partly filled form (.bdff) is a compiled template that holds values of
// some of the names of its data (Template.Filled) and leaves the others
// for later (Template.Open): one person fills in what is theirs and hands
// it on, and another fills in the rest. Its pages show the template filled
// with the values it holds, with the names left open empty; its schema is
// that of the names left open. Filling it with the data of those names
// fills the template with both, once: the document is the one that the
// template filled with all of the data at once is.
//
// The names are those of the schema of the data, with dots between the
// names of objects ("party_b", "party_b.name"). A list is filled, or left
// open, as a whole.

// isObject reports whether a name is an object of names (not a list, nor a
// value).
func (n *schemaNode) isObject() bool { return !n.list() && len(n.kids) > 0 }

// sortedKids are the names of a node in the order of the template.
func (n *schemaNode) sortedKids() []*schemaNode {
	kids := append([]*schemaNode(nil), n.kids...)
	sort.SliceStable(kids, func(a, b int) bool { return kids[a].first < kids[b].first })
	return kids
}

// names lists the names of the data that are filled or left open: those
// of values, pictures, conditions and lists, in objects, with dots.
func (n *schemaNode) names(prefix string, out []string) []string {
	for _, k := range n.sortedKids() {
		p := join(prefix, k.name)
		if k.isObject() {
			out = k.names(p, out)
		} else {
			out = append(out, p)
		}
	}
	return out
}

// node returns the node of a name with dots outside lists.
func (n *schemaNode) node(name string) *schemaNode {
	at := n
	for _, p := range strings.Split(name, ".") {
		if !at.isObject() {
			return nil
		}
		at = at.byName[p]
		if at == nil {
			return nil
		}
	}
	return at
}

func join(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

// under reports whether a name is one of open or in an object of open.
func under(open []string, name string) bool {
	for _, o := range open {
		if name == o || strings.HasPrefix(name, o+".") {
			return true
		}
	}
	return false
}

// openObject writes the schema of the names of a node that are open.
func (n *schemaNode) openObject(prefix string, open []string) *ordered {
	o := &ordered{}
	o.set("type", "object")
	props := &ordered{}
	for _, k := range n.sortedKids() {
		p := join(prefix, k.name)
		switch {
		case under(open, p):
			props.set(k.name, k.schema())
		case k.isObject() && within(open, p):
			ko := k.openObject(p, open)
			if k.label != "" {
				ko.set("title", k.label)
			} else {
				ko.set("title", k.name)
			}
			props.set(k.name, ko)
		}
	}
	o.set("properties", props)
	return o
}

// within reports whether an open name is in an object.
func within(open []string, name string) bool {
	for _, o := range open {
		if strings.HasPrefix(o, name+".") {
			return true
		}
	}
	return false
}

// lookupPath finds the value of a name with dots in data.
func lookupPath(data map[string]any, name string) (any, bool) {
	var v any = data
	for _, p := range strings.Split(name, ".") {
		m, ok := v.(map[string]any)
		if !ok {
			return nil, false
		}
		if v, ok = m[p]; !ok {
			return nil, false
		}
	}
	return v, true
}

// checkGiven checks that the data gives only names left open, and none
// that are filled: prefix is where data is.
func checkGiven(data map[string]any, prefix string, open []string, filled map[string]any) error {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		p := join(prefix, k)
		if under(open, p) {
			continue
		}
		if within(open, p) {
			m, ok := data[k].(map[string]any)
			if !ok {
				return fmt.Errorf("form: %q of the data is not an object: the form leaves names in it open", p)
			}
			if err := checkGiven(m, p, open, filled); err != nil {
				return err
			}
			continue
		}
		if _, ok := lookupPath(filled, p); ok {
			return fmt.Errorf("form: %q is filled already; the data may give only the names that the form leaves open", p)
		}
		return fmt.Errorf("form: the form does not leave %q open; the data may give only the names that it leaves open", p)
	}
	return nil
}

// merge copies the values of src into dst, into the objects of the same
// names.
func merge(dst, src map[string]any) {
	for k, v := range src {
		if m, ok := v.(map[string]any); ok {
			if d, ok := dst[k].(map[string]any); ok {
				merge(d, m)
				continue
			}
		}
		dst[k] = v
	}
}

// jsonCopy is data as JSON keeps it: what a partly filled form holds.
func jsonCopy(data map[string]any) (map[string]any, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("form: the data is not data of JSON: %w", err)
	}
	out := map[string]any{}
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("form: the data is not data of JSON: %w", err)
	}
	tag.Normalize(out)
	return out, nil
}

// inlinePictures replaces the paths of pictures in the data with data:
// URLs, so that a partly filled form holds its pictures: the person that
// fills in the rest has not the files.
func inlinePictures(n *schemaNode, v any, opts *Options) (any, error) {
	switch x := v.(type) {
	case map[string]any:
		for k, e := range x {
			kid := n.byName[k]
			if kid == nil {
				continue
			}
			var err error
			if x[k], err = inlinePictures(kid, e, opts); err != nil {
				return nil, err
			}
		}
		return x, nil
	case []any:
		for i, e := range x {
			var err error
			if x[i], err = inlinePictures(n, e, opts); err != nil {
				return nil, err
			}
		}
		return x, nil
	case string:
		if !n.image || n.code != "" || x == "" || strings.HasPrefix(x, "data:") {
			return x, nil
		}
		b, err := tag.PictureData(opts.Files, x)
		if err != nil {
			return nil, fmt.Errorf("form: the picture of %q: %w", n.name, err)
		}
		return "data:" + pictureType(b) + ";base64," + base64.StdEncoding.EncodeToString(b), nil
	}
	return v, nil
}

// pictureType is the media type of a picture.
func pictureType(b []byte) string {
	head := bytes.TrimLeft(b[:min(len(b), 512)], "\ufeff \t\r\n")
	if bytes.HasPrefix(head, []byte("<")) {
		return "image/svg+xml"
	}
	return http.DetectContentType(b)
}

// partial is a fill of a template that is partly filled, or makes one.
type partial struct {
	// data fills the template: the values filled and those given.
	data map[string]any
	// open are the names left open after the fill (none: the form is
	// complete).
	open []string
	// filled are the values of the partly filled form that the fill
	// makes.
	filled map[string]any
}

// prepare works out a fill of a template whose description is def (nil for
// a template that is not compiled), with the fields of the template.
// It returns nil for a fill that has nothing to do with partly filled
// forms.
func prepare(def *Template, fields []tag.Field, opts *Options) (*partial, error) {
	var filled map[string]any
	var open []string
	if def != nil {
		filled, open = def.Filled, def.Open
	}
	if filled == nil && len(open) == 0 && len(opts.Open) == 0 {
		return nil, nil
	}
	if filled == nil {
		filled = map[string]any{}
	}
	given, err := jsonCopy(opts.Data)
	if err != nil {
		return nil, err
	}
	if def != nil && def.Filled != nil {
		if err := checkGiven(given, "", open, filled); err != nil {
			return nil, err
		}
	}
	data, err := jsonCopy(filled)
	if err != nil {
		return nil, err
	}
	merge(data, given)
	root, err := schemaTree(fields, data)
	if err != nil {
		return nil, err
	}
	names := root.names("", nil)
	// the names left open after this fill
	var left []string
	for _, o := range opts.Open {
		if o == OpenAll {
			for _, name := range names {
				if _, ok := lookupPath(data, name); !ok && (def == nil || def.Filled == nil || under(open, name)) {
					left = append(left, name)
				}
			}
			continue
		}
		o = strings.TrimSpace(o)
		if root.node(o) == nil {
			return nil, fmt.Errorf("form: %q is not a name of the data of the template that can be left open (a name in a list is left open with its list)", o)
		}
		if def != nil && def.Filled != nil && !under(open, o) {
			return nil, fmt.Errorf("form: %q is filled already; it cannot be left open again", o)
		}
		left = append(left, o)
	}
	left = uniq(left)
	for _, o := range left {
		if _, ok := lookupPath(given, o); ok {
			return nil, fmt.Errorf("form: the data gives %q, which it leaves open", o)
		}
	}
	for _, name := range names {
		if _, ok := lookupPath(data, name); ok || under(left, name) {
			continue
		}
		if def != nil && def.Filled != nil {
			return nil, fmt.Errorf("form: the data has no %q, which the form leaves open; give it (null for none), or leave it open again", name)
		}
		return nil, fmt.Errorf("form: the data has no %q; give it (null for none), or leave it open", name)
	}
	if len(left) > 0 {
		if _, err := inlinePictures(root, data, opts); err != nil {
			return nil, err
		}
	}
	return &partial{data: data, open: left, filled: data}, nil
}

func uniq(names []string) []string {
	seen := map[string]bool{}
	out := names[:0]
	for _, n := range names {
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	return out
}

// templateDigest is the SHA-256 of the description of a compiled template:
// what a partly filled form names the template it was made of by.
func templateDigest(b []byte) string {
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// partlyFilled makes the document of a partly filled form: doc is the
// template filled with the values of p.
func partlyFilled(doc *bdf.Document, src []byte, def *Template, format string, fields []tag.Field, p *partial) (*Template, error) {
	root, err := schemaTree(fields, p.filled)
	if err != nil {
		return nil, err
	}
	schema, err := writeSchema(root.openObject("", p.open))
	if err != nil {
		return nil, err
	}
	t := &Template{Version: TemplateVersion, Format: format, Fields: fields,
		Schema: schema, Filled: p.filled, Open: p.open}
	if t.Fields == nil {
		t.Fields = []tag.Field{}
	}
	if def != nil {
		t.Name, t.From = def.Name, def.From
		if t.From == "" {
			t.From = def.digest
		}
	}
	t.Source = doc.AddPart(bdf.PartSource, src).String()
	b, err := json.Marshal(t)
	if err != nil {
		return nil, fmt.Errorf("form: %w", err)
	}
	doc.AddPart(bdf.PartTemplate, b)
	return t, nil
}
