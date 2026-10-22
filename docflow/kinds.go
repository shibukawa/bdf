package docflow

import (
	"context"
	"fmt"
	"slices"
	"strings"
)

// Kinds of a payload: what flows along an edge. Every node takes some of
// them at its input (takes) and sends some on at its output (emits). The
// graph is typed when a node is put (typed): the kinds that reach every
// node are inferred from the entries along the edges, and a new edge that
// brings a kind to a node that does not take it (data to a printer, say)
// is refused. A change of a node's settings may leave an edge that no
// longer carries what the node after it takes (an entry taking data now,
// in front of a printer): the edge stays, broken, for the editor to show
// and the administrator to mend. A job's payload is checked again at
// every node it reaches (route), so that a broken edge, or a graph
// changed while a job was on its way, gives the job a leg of the error
// there and does not break it.
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

// passesOn tells a kind of node that sends on what reaches it (and, an
// entry, what it receives): what it may take at its input is what every
// node after it takes, so that its first edge, on either side, types it
// as the wildcard pin of a visual editor is typed; so does an archive,
// which sends on what it wrote. A template or a convert node sends on the
// .bdf it makes whatever reaches it (and what reached it along Original,
// which accepts types), a display the .bdf it shows; the others send
// nothing on.
func passesOn(kind string) bool {
	switch kind {
	case KindQueue, KindClass, KindFilter, KindInbox, KindEndpoint, KindCups, KindArchive:
		return true
	}
	return false
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
// what reaches it, a template or a convert node the .bdf it makes (and
// what reaches it, when it sends that on along Original: emitsTo), a
// display the .bdf its fallback printer prints, every other node what
// reaches it.
func (n *Node) emits(in kinds) kinds {
	switch n.Kind {
	case KindEndpoint, KindCups:
		return in | n.entryKind()
	case KindTemplate, KindConvert:
		if len(n.Original) > 0 {
			return kindBDF | in
		}
		return kindBDF
	case KindDisplay:
		return kindBDF
	}
	return in
}

// emitsTo tells the kinds a node sends on along its edge to a node: of a
// template or a convert node, what reaches it along Original (as it took
// it) and the .bdf it makes along To; of the others, what it emits.
func (n *Node) emitsTo(in kinds, to string) kinds {
	if n.Kind != KindTemplate && n.Kind != KindConvert {
		return n.emits(in)
	}
	var out kinds
	if slices.Contains(n.To, to) {
		out |= kindBDF
	}
	if slices.Contains(n.Original, to) {
		out |= in
	}
	return out
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
			for _, to := range n.edges() {
				if nodes[to] == nil {
					continue
				}
				add := n.emitsTo(t.reach[n.Name], to) &^ t.reach[to]
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

// accepts tells, for every node, the kinds its input may take now: what
// the node takes, and, for a node that passes on what reaches it (or a
// template or a convert node, along Original), what every node it sends
// it to accepts in turn (a queue in front of a printer
// takes .bdf alone; one in front of a printer and a template node takes
// nothing). A graph of nodes that are not there yet is not refused here:
// an edge to a name that is no node is left out.
func accepts(nodes map[string]*Node) map[string]kinds {
	acc := make(map[string]kinds, len(nodes))
	for name, n := range nodes {
		acc[name] = takes(n.Kind)
	}
	for range len(nodes) + 2 {
		changed := false
		for name, n := range nodes {
			along := n.edges()
			if !passesOn(n.Kind) {
				if len(n.Original) == 0 {
					continue
				}
				along = n.Original
			}
			a := takes(n.Kind)
			for _, to := range along {
				if nodes[to] != nil {
					a &= acc[to]
				}
			}
			if a != acc[name] {
				acc[name] = a
				changed = true
			}
		}
		if !changed {
			break
		}
	}
	return acc
}

// broken tells the edges along which a node sends a kind the node after it
// does not take: by the node that sends, then the node sent to, the kinds.
func (t *typing) broken(nodes map[string]*Node) map[string]map[string]kinds {
	out := map[string]map[string]kinds{}
	for _, n := range nodes {
		for _, to := range n.edges() {
			x := nodes[to]
			if x == nil {
				continue
			}
			if bad := n.emitsTo(t.reach[n.Name], to) &^ takes(x.Kind); bad != 0 {
				if out[n.Name] == nil {
					out[n.Name] = map[string]kinds{}
				}
				out[n.Name][to] |= bad
			}
		}
	}
	return out
}

// withoutEdges is a copy of the node without its edges to the names
// given.
func (n *Node) withoutEdges(names map[string]bool) *Node {
	keep := func(list []string) []string {
		var out []string
		for _, x := range list {
			if !names[x] {
				out = append(out, x)
			}
		}
		return out
	}
	c := *n
	c.To, c.Original, c.Members, c.Otherwise = keep(n.To), keep(n.Original), keep(n.Members), keep(n.Otherwise)
	c.Rules = slices.Clone(n.Rules)
	for i := range c.Rules {
		c.Rules[i].To = keep(n.Rules[i].To)
	}
	if c.Fallback != nil && names[c.Fallback.Printer] {
		c.Fallback = nil
	}
	return &c
}

// typed refuses a new edge of n (one the node did not have) along which a
// kind reaches a node that does not take it, with n put as it is: every
// node is checked, since an edge of n changes what reaches the nodes
// after it. A change of the settings of n alone (an entry taking data
// now, say; a queue made a template node) is not refused for the edges
// it leaves broken: the topology shows them (TopoNode.Broken), and a job
// gets a leg of the error where it reaches a node that does not take it.
func (m *Manager) typed(ctx context.Context, n *Node) error {
	list, err := m.Nodes(ctx)
	if err != nil {
		return err
	}
	nodes := map[string]*Node{}
	for i := range list {
		nodes[list[i].Name] = &list[i]
	}
	fresh := map[string]bool{}
	for _, to := range n.edges() {
		fresh[to] = true
	}
	if cur := nodes[n.Name]; cur != nil {
		for _, to := range cur.edges() {
			delete(fresh, to)
		}
	}
	nodes[n.Name] = n.withoutEdges(fresh)
	before := inferKinds(nodes).broken(nodes)
	nodes[n.Name] = n
	t := inferKinds(nodes)
	after := t.broken(nodes)
	froms := make([]string, 0, len(after))
	for from := range after {
		froms = append(froms, from)
	}
	slices.Sort(froms)
	for _, from := range froms {
		tos := make([]string, 0, len(after[from]))
		for to := range after[from] {
			tos = append(tos, to)
		}
		slices.Sort(tos)
		for _, to := range tos {
			bad := after[from][to] &^ before[from][to]
			if bad == 0 {
				continue
			}
			kind := bad.names()[0]
			x := nodes[to]
			return fmt.Errorf("%w: %s (%s) takes %s, but %s reaches it: %s", ErrInvalid, to, x.Kind, takes(x.Kind), kind, strings.Join(append(t.path(from, kind), to), " → "))
		}
	}
	return nil
}
