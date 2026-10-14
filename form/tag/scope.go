package tag

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// The values of the data of a template are those of JSON, and dates: nil,
// bool, float64, string, time.Time, []any and map[string]any (see
// Normalize).

// Scope is the data that the tags at a place of a template see: the item of
// the section they are in, inside the items of the sections around it, in
// the data of the template. A name is looked up from the innermost item
// outwards.
type Scope struct {
	parent *Scope
	value  any
	vars   map[string]any
	// blank is the scope of a place that no item fills (an empty slot):
	// the names of the data are not found from it, the variables are
	blank bool
}

// NewScope returns the scope of the data of a template.
func NewScope(root any) *Scope { return &Scope{value: root} }

// Item returns the scope of an item in s, with variables about its place
// (their names start with "_": "_index", "_first" …).
func (s *Scope) Item(value any, vars map[string]any) *Scope {
	return &Scope{parent: s, value: value, vars: vars}
}

// Blank returns the scope of a place in s that no item fills: only the
// variables, its own and those around it, are found from it.
func (s *Scope) Blank(vars map[string]any) *Scope {
	return &Scope{parent: s, vars: vars, blank: true}
}

// Value returns the item of the scope.
func (s *Scope) Value() any { return s.value }

// IsBlank reports whether the scope is of a place that no item fills, or
// is in one: names that are not found from it are not missing from the
// data.
func (s *Scope) IsBlank() bool {
	for ; s != nil; s = s.parent {
		if s.blank {
			return true
		}
	}
	return false
}

// Lookup finds the value of the name of a tag. The first part of the name
// is looked up in the item of the scope, then in those around it; the
// others in the value found.
func (s *Scope) Lookup(path []string) (any, bool) {
	if len(path) == 0 {
		return nil, false
	}
	if len(path) == 1 && path[0] == "." {
		return s.value, !s.blank
	}
	var v any
	found, varsOnly := false, false
	for sc := s; sc != nil && !found; sc = sc.parent {
		if x, ok := sc.vars[path[0]]; ok {
			v, found = x, true
			break
		}
		if m, ok := sc.value.(map[string]any); ok && !varsOnly {
			if x, ok := m[path[0]]; ok {
				v, found = x, true
				break
			}
		}
		varsOnly = varsOnly || sc.blank
	}
	if !found {
		return nil, false
	}
	for _, p := range path[1:] {
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

// Truthy reports whether a value shows a section: not nil, false, zero, an
// empty string or an empty list.
func Truthy(v any) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	case float64:
		return x != 0
	case string:
		return x != ""
	case []any:
		return len(x) > 0
	}
	return true
}

// Items returns what a section of a value is repeated for: the items of a
// list, the value itself when it is another true value, and nothing when
// it is false.
func Items(v any) []any {
	if l, ok := v.([]any); ok {
		return l
	}
	if Truthy(v) {
		return []any{v}
	}
	return nil
}

// Normalize returns data in the values that templates are filled with: the
// numbers as float64, the lists as []any, the objects as map[string]any
// with keys that are strings, whatever the decoder that read them made of
// them. Values of other kinds become their text.
func Normalize(v any) any {
	switch x := v.(type) {
	case nil, bool, float64, string, time.Time:
		return x
	case json.Number:
		if f, err := x.Float64(); err == nil {
			return f
		}
		return x.String()
	case map[string]any:
		for k, e := range x {
			x[k] = Normalize(e)
		}
		return x
	case []any:
		for i, e := range x {
			x[i] = Normalize(e)
		}
		return x
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return rv.Float()
	case reflect.Map:
		out := make(map[string]any, rv.Len())
		for it := rv.MapRange(); it.Next(); {
			out[fmt.Sprint(it.Key().Interface())] = Normalize(it.Value().Interface())
		}
		return out
	case reflect.Slice, reflect.Array:
		out := make([]any, rv.Len())
		for i := range out {
			out[i] = Normalize(rv.Index(i).Interface())
		}
		return out
	case reflect.Pointer, reflect.Interface:
		if rv.IsNil() {
			return nil
		}
		return Normalize(rv.Elem().Interface())
	}
	return fmt.Sprint(v)
}

// ParseDate reads the text of a date as data writes it: a date
// ("2026-10-13", "2026/10/13"), or a date and a time ("2026-10-13T09:30:00",
// with or without a zone, or with a space for the T). The time is that of
// the clock on the wall: the zone is dropped, as a spreadsheet has none.
func ParseDate(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if len(s) < 8 || len(s) > 40 || s[0] < '0' || s[0] > '9' {
		return time.Time{}, false
	}
	for _, layout := range []string{"2006-01-02", "2006/1/2", time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02T15:04",
		"2006-01-02 15:04:05", "2006-01-02 15:04", "2006/1/2 15:04:05", "2006/1/2 15:04"} {
		if t, err := time.Parse(layout, s); err == nil {
			return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC), true
		}
	}
	return time.Time{}, false
}

// Formatter writes a value in a number format: the format codes of Excel
// ("#,##0", "yyyy/mm/dd"). An empty format is the General format, as a cell
// without a format shows a number. (converter/xlsx has one: Formatter.)
type Formatter func(v any, format string) string

// ValueText writes a value as the text that replaces a tag in the text of
// a document: in the format of the tag, with the formatter; without a format,
// a number as the General format writes it, a date as yyyy-mm-dd (with
// hh:mm, and :ss, when it has a time), true and false as TRUE and FALSE,
// nothing for no value. Lists and objects have no text. Without a
// formatter, a format is not applied.
func ValueText(v any, format string, f Formatter) string {
	if format != "" && f != nil {
		switch v.(type) {
		case float64, string, time.Time, bool:
			return f(v, format)
		}
	}
	switch x := v.(type) {
	case string:
		return x
	case bool:
		if x {
			return "TRUE"
		}
		return "FALSE"
	case float64:
		if f != nil {
			return f(x, "")
		}
		return strconv.FormatFloat(x, 'g', 15, 64)
	case time.Time:
		switch {
		case x.Hour() == 0 && x.Minute() == 0 && x.Second() == 0:
			return x.Format("2006-01-02")
		case x.Second() == 0:
			return x.Format("2006-01-02 15:04")
		}
		return x.Format("2006-01-02 15:04:05")
	}
	return ""
}

// Vars returns the variables of item k of n (see Scope.Item): "_index" and
// "_count", "_first" and "_last", and, for an item that is a page, "_page"
// and "_pages".
func Vars(k, n int, page bool) map[string]any {
	v := map[string]any{"_index": float64(k + 1), "_count": float64(n), "_first": k == 0, "_last": k == n-1}
	if page {
		v["_page"], v["_pages"] = float64(k+1), float64(n)
	}
	return v
}
