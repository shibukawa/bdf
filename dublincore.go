package bdf

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DublinCore describes the document with the fifteen elements of the Dublin
// Core Metadata Element Set 1.1, plus the DCMI Terms created and modified
// (docs/spec.md §4.3). Every element is optional and repeatable.
type DublinCore struct {
	Title       DCValues `json:"title,omitempty"`
	Creator     DCValues `json:"creator,omitempty"`
	Subject     DCValues `json:"subject,omitempty"`
	Description DCValues `json:"description,omitempty"`
	Publisher   DCValues `json:"publisher,omitempty"`
	Contributor DCValues `json:"contributor,omitempty"`
	Date        DCValues `json:"date,omitempty"`
	Type        DCValues `json:"type,omitempty"`
	Format      DCValues `json:"format,omitempty"`
	Identifier  DCValues `json:"identifier,omitempty"`
	Source      DCValues `json:"source,omitempty"`
	Language    DCValues `json:"language,omitempty"`
	Relation    DCValues `json:"relation,omitempty"`
	Coverage    DCValues `json:"coverage,omitempty"`
	Rights      DCValues `json:"rights,omitempty"`
	Created     DCValues `json:"created,omitempty"`
	Modified    DCValues `json:"modified,omitempty"`
}

// DCTerms lists the element names DublinCore holds, in field order.
var DCTerms = []string{
	"title", "creator", "subject", "description", "publisher", "contributor", "date", "type",
	"format", "identifier", "source", "language", "relation", "coverage", "rights", "created", "modified",
}

// Field returns the values of the element with the given name (one of
// DCTerms), or nil for an unknown name.
func (dc *DublinCore) Field(name string) *DCValues {
	switch name {
	case "title":
		return &dc.Title
	case "creator":
		return &dc.Creator
	case "subject":
		return &dc.Subject
	case "description":
		return &dc.Description
	case "publisher":
		return &dc.Publisher
	case "contributor":
		return &dc.Contributor
	case "date":
		return &dc.Date
	case "type":
		return &dc.Type
	case "format":
		return &dc.Format
	case "identifier":
		return &dc.Identifier
	case "source":
		return &dc.Source
	case "language":
		return &dc.Language
	case "relation":
		return &dc.Relation
	case "coverage":
		return &dc.Coverage
	case "rights":
		return &dc.Rights
	case "created":
		return &dc.Created
	case "modified":
		return &dc.Modified
	}
	return nil
}

// IsZero reports whether no element has a value.
func (dc DublinCore) IsZero() bool {
	for _, name := range DCTerms {
		if len(*dc.Field(name)) > 0 {
			return false
		}
	}
	return true
}

// DCValues holds the values of one element. In JSON a single value is a
// string and several values are an array of strings.
type DCValues []string

// First returns the first value, or "" when there is none.
func (v DCValues) First() string {
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// SplitKeywords splits a keyword list as PDF and Office documents store it
// into subjects. Keywords are separated by commas or semicolons (ASCII or
// full-width) or by the ideographic comma.
func SplitKeywords(s string) DCValues {
	var out DCValues
	for _, k := range strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case ',', ';', '，', '；', '、':
			return true
		}
		return false
	}) {
		if k = strings.TrimSpace(k); k != "" {
			out = append(out, k)
		}
	}
	return out
}

// MarshalJSON writes one value as a string and several as an array.
func (v DCValues) MarshalJSON() ([]byte, error) {
	if len(v) == 1 {
		return json.Marshal(v[0])
	}
	return json.Marshal([]string(v))
}

// UnmarshalJSON accepts a string, an array of strings or null.
func (v *DCValues) UnmarshalJSON(b []byte) error {
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		*v = DCValues{s}
		return nil
	}
	var a []string
	if err := json.Unmarshal(b, &a); err != nil {
		return fmt.Errorf("bdf: a Dublin Core value must be a string or an array of strings: %s", b)
	}
	*v = a
	return nil
}
