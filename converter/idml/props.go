package idml

import (
	"strings"

	"github.com/shibukawa/bdf/converter/internal/ooxml"
)

// IDML writes an object's properties two ways: simple values as attributes
// (PointSize="12"), and typed values as children of a Properties element
// (<AppliedFont type="string">Minion Pro</AppliedFont>, lists as ListItem
// children, <Leading type="enumeration">Auto</Leading>). A style's
// BasedOn is one of the latter.

// prop returns a property of an element, from its attributes or its
// Properties.
func prop(n *ooxml.Node, name string) (string, bool) {
	if n == nil {
		return "", false
	}
	if v, ok := n.Attr(name); ok {
		return v, true
	}
	if k := n.Child("Properties").Child(name); k != nil {
		if k.AttrStr("type", "") == "list" {
			var parts []string
			for _, it := range k.Children("ListItem") {
				parts = append(parts, strings.TrimSpace(it.Content()))
			}
			return strings.Join(parts, " "), true
		}
		return strings.TrimSpace(k.Content()), true
	}
	return "", false
}

// listValues returns the items of a list property (its ListItem children,
// or the fields of its attribute).
func listValues(n *ooxml.Node, name string) []string {
	if k := n.Child("Properties").Child(name); k != nil {
		var out []string
		for _, it := range k.Children("ListItem") {
			out = append(out, strings.TrimSpace(it.Content()))
		}
		return out
	}
	return strings.Fields(n.AttrStr(name, ""))
}

// styleChain lists an element's applied style and the styles it is based
// on, nearest first. attr names the attribute holding the reference
// (AppliedParagraphStyle, AppliedCharacterStyle); styles is the table of
// the styles by Self.
func styleChain(n *ooxml.Node, attr string, styles map[string]*ooxml.Node) []*ooxml.Node {
	var out []*ooxml.Node
	ref, _ := prop(n, attr)
	// a BasedOn names the style without its kind ("$ID/[No paragraph
	// style]") as often as with it ("ParagraphStyle/$ID/NormalParagraphStyle")
	prefix := strings.TrimPrefix(attr, "Applied") + "/"
	seen := map[string]bool{}
	for ref != "" && !seen[ref] && len(out) < 32 {
		seen[ref] = true
		s := styles[ref]
		if s == nil {
			s = styles[prefix+ref]
		}
		if s == nil {
			break
		}
		out = append(out, s)
		ref, _ = prop(s, "BasedOn")
	}
	return out
}

// lookup returns the first value of a property along a chain of elements.
func lookup(chain []*ooxml.Node, name string) (string, bool) {
	for _, n := range chain {
		if v, ok := prop(n, name); ok {
			return v, true
		}
	}
	return "", false
}
