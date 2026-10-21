package manager

import (
	"context"
	"fmt"
	"slices"
	"strings"
)

// Kinds of a payload: what flows along an edge. Every node takes some of
// them at its input (takes) and sends some on at its output (emits). The
// graph is typed when a node is put (typed): the kinds that reach every
// node are inferred from the entries along the edges, and a kind reaching
// a node that does not take it (data a printer, say) is refused. A job's
// payload is checked again at every node it reaches (route), so that a
// graph changed while a job was on its way does not break it.
const (
	// PayloadBDF is a .bdf document: what printers, displays and the
	// pages take.
	PayloadBDF = "bdf"
	// PayloadData is data: JSON, YAML, CSV, an Excel workbook or a
	// Parquet file, which a template node fills a template with.
	PayloadData = "data"
	// PayloadAny is any other file as it was received (a PDF, an Office
	// document, an image …), which a convert node makes a .bdf of.
	PayloadAny = "any"
)

var payloadKinds = []string{PayloadBDF, PayloadData, PayloadAny}

// dataFormats are the formats of the kind data, as the entries name them.
var dataFormats = []string{"json", "yaml", "csv", "xlsx", "parquet"}

// kindOfFormat tells the kind of a payload of a format.
func kindOfFormat(format string) string {
	switch {
	case format == "bdf":
		return PayloadBDF
	case slices.Contains(dataFormats, format):
		return PayloadData
	}
	return PayloadAny
}

// kinds is a set of payload kinds.
type kinds uint8

const (
	kindBDF kinds = 1 << iota
	kindData
	kindAny
	allKinds = kindBDF | kindData | kindAny
)

func kindBit(kind string) kinds {
	switch kind {
	case PayloadBDF:
		return kindBDF
	case PayloadData:
		return kindData
	case PayloadAny:
		return kindAny
	}
	return 0
}

func (k kinds) has(kind string) bool { return k&kindBit(kind) != 0 }

// names are the kinds of the set, in order.
func (k kinds) names() []string {
	var out []string
	for _, kind := range payloadKinds {
		if k&kindBit(kind) != 0 {
			out = append(out, kind)
		}
	}
	return out
}

func (k kinds) String() string { return strings.Join(k.names(), ", ") }

// takes tells the kinds a node of a kind takes at its input: printers and
// displays .bdf documents, template nodes data, every other node
// anything.
func takes(kind string) kinds {
	switch kind {
	case KindPrinter, KindDisplay:
		return kindBDF
	case KindTemplate:
		return kindData
	}
	return allKinds
}

// takesName is the input of a kind of node as the API shows it: bdf,
// data, or any.
func takesName(kind string) string {
	switch takes(kind) {
	case kindBDF:
		return PayloadBDF
	case kindData:
		return PayloadData
	}
	return PayloadAny
}

// entryKind is the kind an entry sends on of what it receives from
// outside: data, the file as received (any), or the .bdf made of it.
func (n *Node) entryKind() kinds {
	switch {
	case n.Accept == AcceptData:
		return kindData
	case n.Output == OutputFile:
		return kindAny
	}
	return kindBDF
}

// emits tells the kinds a node sends on along its edges, given the kinds
// that reach it (in): an entry what it receives from outside besides
// what reaches it, a template or a convert node the .bdf it makes, a
// display the .bdf its fallback printer prints, every other node what
// reaches it.
func (n *Node) emits(in kinds) kinds {
	switch n.Kind {
	case KindEndpoint, KindCups:
		return in | n.entryKind()
	case KindTemplate, KindConvert, KindDisplay:
		return kindBDF
	}
	return in
}

// typing is the kinds that reach each node, and from where.
type typing struct {
	reach map[string]kinds
	// from names, for a node and a kind, the node the kind came from.
	from map[string]map[string]string
}

// inferKinds computes the kinds that reach every node from the entries
// along the edges, to a fixpoint (what a node emits grows with its
// input alone).
func inferKinds(nodes map[string]*Node) *typing {
	t := &typing{reach: map[string]kinds{}, from: map[string]map[string]string{}}
	for range len(nodes) + 2 {
		changed := false
		for _, n := range nodes {
			out := n.emits(t.reach[n.Name])
			for _, to := range n.edges() {
				if nodes[to] == nil {
					continue
				}
				add := out &^ t.reach[to]
				if add == 0 {
					continue
				}
				t.reach[to] |= add
				if t.from[to] == nil {
					t.from[to] = map[string]string{}
				}
				for _, kind := range add.names() {
					t.from[to][kind] = n.Name
				}
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return t
}

// path tells where a kind reaching a node came from: the nodes from the
// entry on, in order.
func (t *typing) path(node, kind string) []string {
	path := []string{node}
	seen := map[string]bool{node: true}
	for {
		prev := t.from[node][kind]
		if prev == "" || seen[prev] {
			return path
		}
		seen[prev] = true
		path = append([]string{prev}, path...)
		node = prev
	}
}

// typed refuses a graph in which a kind reaches a node that does not take
// it, with n put as it is: every node is checked, since a change of n
// (an entry taking data now, say) changes what reaches the nodes after it.
func (m *Manager) typed(ctx context.Context, n *Node) error {
	list, err := m.Nodes(ctx)
	if err != nil {
		return err
	}
	nodes := map[string]*Node{}
	for i := range list {
		nodes[list[i].Name] = &list[i]
	}
	nodes[n.Name] = n
	t := inferKinds(nodes)
	names := make([]string, 0, len(nodes))
	for name := range nodes {
		names = append(names, name)
	}
	slices.Sort(names)
	for _, name := range names {
		x := nodes[name]
		if bad := t.reach[name] &^ takes(x.Kind); bad != 0 {
			kind := bad.names()[0]
			return fmt.Errorf("%w: %s (%s) takes %s, but %s reaches it: %s", ErrInvalid, name, x.Kind, takes(x.Kind), kind, strings.Join(t.path(name, kind), " → "))
		}
	}
	return nil
}
