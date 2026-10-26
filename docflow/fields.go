package docflow

import (
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// A kind of node registered by a program of its own (a module that
// imports docflow, besides the core and its adapters) has settings of
// its own: Node.Settings, an object of the fields the kind declares
// (Kind.Fields). The editor makes a form of the fields, with what they
// are called and what they take in the language of the page; PutNode
// checks a node's settings against them (checkSettings) before the
// kind's Check sees them, so that a kind reads them with SettingsOf,
// filled in with their defaults, and trusts their shape.

// Text is a text in the languages the pages are written in, by language
// tag (BCP 47: "en", "ja"): English at least, which is read for a
// language that is not there.
type Text map[string]string

// In is the text in a language: the language's, else its primary
// language's ("ja" for "ja-JP"), else English, else the first there is.
func (t Text) In(lang string) string {
	if s := t[lang]; s != "" {
		return s
	}
	if base, _, ok := strings.Cut(lang, "-"); ok {
		if s := t[base]; s != "" {
			return s
		}
	}
	if s := t["en"]; s != "" {
		return s
	}
	for _, tag := range slices.Sorted(maps.Keys(t)) {
		if t[tag] != "" {
			return t[tag]
		}
	}
	return ""
}

// FieldType is what a field of a kind's own settings takes.
type FieldType string

const (
	// FieldString is a line of text.
	FieldString FieldType = "string"
	// FieldText is text of several lines.
	FieldText FieldType = "text"
	// FieldInt is a whole number.
	FieldInt FieldType = "int"
	// FieldBool is yes or no.
	FieldBool FieldType = "bool"
	// FieldDuration is a Duration, written as "30s", "10m", "2h" or "14d".
	FieldDuration FieldType = "duration"
	// FieldEnum is one of the field's Choices.
	FieldEnum FieldType = "enum"
	// FieldCredential is the name of a credential (credential.go), never
	// a secret.
	FieldCredential FieldType = "credential"
	// FieldPath is a path or a key with {{meta.<key>}} in it, which the
	// kind expands (Out.ExpandPath).
	FieldPath FieldType = "path"
	// FieldJSON is any JSON value, as it is.
	FieldJSON FieldType = "json"
)

// Choice is one of the values a field of FieldEnum takes.
type Choice struct {
	Value string `json:"value"`
	// Label is how it is said; none, the value.
	Label Text `json:"label,omitempty"`
}

// Field is a setting of a kind of node of its own: its key in
// Node.Settings, what it takes, and how the editor shows it.
type Field struct {
	// Key is the setting's key: letters, digits and _, starting with a
	// letter.
	Key  string    `json:"key"`
	Type FieldType `json:"type"`
	// Title is the field as it is said; Help what it is for, under the
	// field; Placeholder what an empty field shows.
	Title       Text `json:"title,omitempty"`
	Help        Text `json:"help,omitempty"`
	Placeholder Text `json:"placeholder,omitempty"`
	// Required refuses a node without the setting (empty is without),
	// when the settings of If have the values given (as Required.If);
	// always, with no If.
	Required bool              `json:"required,omitempty"`
	If       map[string]string `json:"if,omitempty"`
	// Choices are the values a FieldEnum takes.
	Choices []Choice `json:"choices,omitempty"`
	// Default is the value SettingsOf reads for a setting that is not
	// there (a JSON value of the field's type); the editor shows it.
	Default any `json:"default,omitempty"`
	// Expands tells a field whose value has {{meta.<key>}} replaced by
	// the job's metadata (the editor says so).
	Expands bool `json:"expands,omitempty"`
}

// fieldKeyRE is the shape of a field's key.
var fieldKeyRE = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)

// checkFields checks the fields a kind registers (RegisterKind).
func checkFields(fields []Field) error {
	seen := map[string]bool{}
	for _, f := range fields {
		if !fieldKeyRE.MatchString(f.Key) {
			return fmt.Errorf("field key %q: lowercase letters, digits and _, starting with a letter", f.Key)
		}
		if seen[f.Key] {
			return fmt.Errorf("field %s twice", f.Key)
		}
		seen[f.Key] = true
		switch f.Type {
		case FieldString, FieldText, FieldInt, FieldBool, FieldDuration, FieldCredential, FieldPath, FieldJSON:
		case FieldEnum:
			if len(f.Choices) == 0 {
				return fmt.Errorf("field %s: an enum has choices", f.Key)
			}
		default:
			return fmt.Errorf("field %s: type %q", f.Key, f.Type)
		}
		if f.Default != nil {
			if _, err := f.clean(f.Default); err != nil {
				return fmt.Errorf("field %s: default: %v", f.Key, err)
			}
		}
	}
	return nil
}

// fits tells whether the settings have the values a field's If names
// (none: always).
func (f *Field) fits(settings map[string]any) bool {
	for k, want := range f.If {
		if !strings.EqualFold(strings.TrimSpace(metaText(settings[k])), want) {
			return false
		}
	}
	return true
}

// clean checks a value against the field and returns it as it is kept:
// nil for one that is empty (not set).
func (f *Field) clean(v any) (any, error) {
	switch f.Type {
	case FieldString, FieldText, FieldPath, FieldCredential, FieldDuration, FieldEnum:
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%v is not text", v)
		}
		if f.Type != FieldText {
			s = strings.TrimSpace(s)
		}
		if s == "" {
			return nil, nil
		}
		switch f.Type {
		case FieldCredential:
			if !credentialRE.MatchString(s) {
				return nil, fmt.Errorf("credential %q: letters, digits, _ and -", s)
			}
		case FieldDuration:
			d, err := ParseDuration(s)
			if err != nil {
				return nil, fmt.Errorf("%q is not a duration (30s, 10m, 2h, 14d)", s)
			}
			s = d.String()
		case FieldEnum:
			if !slices.ContainsFunc(f.Choices, func(c Choice) bool { return c.Value == s }) {
				var values []string
				for _, c := range f.Choices {
					values = append(values, c.Value)
				}
				return nil, fmt.Errorf("%q is not one of %s", s, strings.Join(values, ", "))
			}
		}
		return s, nil
	case FieldInt:
		switch x := v.(type) {
		case float64:
			if x != math.Trunc(x) || math.Abs(x) > 1<<53 {
				return nil, fmt.Errorf("%v is not a whole number", v)
			}
			return int64(x), nil
		case int64:
			return x, nil
		case int:
			return int64(x), nil
		case json.Number:
			n, err := x.Int64()
			if err != nil {
				return nil, fmt.Errorf("%v is not a whole number", v)
			}
			return n, nil
		case string:
			x = strings.TrimSpace(x)
			if x == "" {
				return nil, nil
			}
			n, err := strconv.ParseInt(x, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%q is not a whole number", x)
			}
			return n, nil
		}
		return nil, fmt.Errorf("%v is not a whole number", v)
	case FieldBool:
		switch x := v.(type) {
		case bool:
			return x, nil
		case string:
			switch strings.ToLower(strings.TrimSpace(x)) {
			case "":
				return nil, nil
			case "true", "yes", "on", "1":
				return true, nil
			case "false", "no", "off", "0":
				return false, nil
			}
		}
		return nil, fmt.Errorf("%v is not true or false", v)
	case FieldJSON:
		return v, nil
	}
	return nil, fmt.Errorf("type %q", f.Type)
}

// checkSettings checks the settings of a node of a kind against the
// kind's fields and keeps them in n, canonical (the keys in order,
// without the empty ones); a kind without fields keeps none. It is the
// check before the kind's own (PutNode).
func checkSettings(k *Kind, in Node, n *Node) error {
	n.Settings = nil
	if len(k.Fields) == 0 {
		return nil
	}
	raw := map[string]any{}
	if len(in.Settings) > 0 && string(in.Settings) != "null" {
		dec := json.NewDecoder(strings.NewReader(string(in.Settings)))
		if err := dec.Decode(&raw); err != nil {
			return fmt.Errorf("%w: settings: %v", ErrInvalid, err)
		}
	}
	for key := range raw {
		if !slices.ContainsFunc(k.Fields, func(f Field) bool { return f.Key == key }) {
			return fmt.Errorf("%w: settings: a %s has no setting %q", ErrInvalid, k.Name, key)
		}
	}
	out := map[string]any{}
	for i := range k.Fields {
		f := &k.Fields[i]
		v, ok := raw[f.Key]
		if !ok || v == nil {
			continue
		}
		x, err := f.clean(v)
		if err != nil {
			return fmt.Errorf("%w: settings: %s: %v", ErrInvalid, f.Key, err)
		}
		if x != nil {
			out[f.Key] = x
		}
	}
	for i := range k.Fields {
		f := &k.Fields[i]
		if f.Required && out[f.Key] == nil && f.fits(out) {
			return fmt.Errorf("%w: settings: %s (%s) is required", ErrInvalid, f.Key, f.Title.In("en"))
		}
	}
	if len(out) == 0 {
		return nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("%w: settings: %v", ErrInvalid, err)
	}
	n.Settings = b
	return nil
}

// SettingsOf reads a node's settings (Node.Settings) into a value of the
// kind's own type, its fields' defaults first (a setting not there is its
// default), as encoding/json reads them: the fields' keys are the JSON
// keys. A node of a kind that is not registered reads its settings
// alone.
func SettingsOf[T any](n *Node) (T, error) {
	var out T
	if k := kindNamed(n.Kind); k != nil {
		defaults := map[string]any{}
		for _, f := range k.Fields {
			if f.Default != nil {
				defaults[f.Key] = f.Default
			}
		}
		if len(defaults) > 0 {
			b, err := json.Marshal(defaults)
			if err != nil {
				return out, fmt.Errorf("docflow: the defaults of %s: %w", n.Kind, err)
			}
			if err := json.Unmarshal(b, &out); err != nil {
				return out, fmt.Errorf("docflow: the defaults of %s: %w", n.Kind, err)
			}
		}
	}
	if len(n.Settings) > 0 {
		if err := json.Unmarshal(n.Settings, &out); err != nil {
			return out, fmt.Errorf("docflow: the settings of %s: %w", n.Name, err)
		}
	}
	return out, nil
}

// SetSettings writes a value of the kind's own type as a node's settings
// (a Check that cleans them, say): what encoding/json makes of it.
func SetSettings(n *Node, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	n.Settings = b
	return nil
}
