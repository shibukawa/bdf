package ipp

import (
	"slices"
	"strings"
)

// JobTemplateAttrs are the job template attributes of a request that makes
// a job (RFC 8011 §5.2, with those of IPP Everywhere): what "job-template"
// stands for in requested-attributes, and what a printer lists in
// job-creation-attributes-supported.
var JobTemplateAttrs = []string{
	"copies", "finishings", "ipp-attribute-fidelity", "job-account-id", "job-accounting-user-id", "job-hold-until", "job-name",
	"job-priority", "media", "media-col",
	"multiple-document-handling", "number-up", "orientation-requested", "output-bin", "overrides", "page-ranges",
	"print-color-mode", "print-content-optimize", "print-quality", "print-rendering-intent", "print-scaling",
	"printer-resolution", "sides",
}

// Requested is the requested-attributes of the operation group of a
// request, or def without it.
func Requested(op *Group, def []string) []string {
	if a := op.Get("requested-attributes"); a != nil {
		return a.Strings()
	}
	return def
}

// Filter keeps the attributes requested: by name, or by the names of
// groups ("all", "job-template" …) that groups maps to the names in them.
func Filter(attrs []Attribute, req []string, groups map[string][]string) []Attribute {
	want := map[string]bool{}
	for _, r := range req {
		if r == "all" {
			return attrs
		}
		if names, ok := groups[r]; ok {
			for _, n := range names {
				want[n] = true
			}
			continue
		}
		want[r] = true
	}
	var out []Attribute
	for _, a := range attrs {
		if want[a.Name] {
			out = append(out, a)
		}
	}
	return out
}

// PrinterGroups are the groups of a printer's attributes that
// requested-attributes names: "job-template" the -default, -supported,
// -ready and -database of the job template attributes, and
// "printer-description" the others.
func PrinterGroups(attrs []Attribute) map[string][]string {
	var tmpl, desc []string
	for _, a := range attrs {
		base := a.Name
		for _, suf := range []string{"-default", "-supported", "-ready", "-database"} {
			base = strings.TrimSuffix(base, suf)
		}
		if base != a.Name && (slices.Contains(JobTemplateAttrs, base) || base == "media-col" || base == "media") {
			tmpl = append(tmpl, a.Name)
		} else {
			desc = append(desc, a.Name)
		}
	}
	return map[string][]string{"job-template": tmpl, "printer-description": desc}
}

// Keywords are the values of an attribute of keywords.
func Keywords(s []string) []Value {
	out := make([]Value, len(s))
	for i, v := range s {
		out[i] = Keyword(v)
	}
	return out
}

// MimeTypes are the values of an attribute of MIME types.
func MimeTypes(s []string) []Value {
	out := make([]Value, len(s))
	for i, v := range s {
		out[i] = MimeType(v)
	}
	return out
}

// Enums are the values of an attribute of enums.
func Enums(ns []int) []Value {
	out := make([]Value, len(ns))
	for i, n := range ns {
		out[i] = Enum(n)
	}
	return out
}
