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
	// PayloadStocked is a file of any format that an inbox or an archive
	// keeps, sent along its edge to a node that reads it there: HTTP
	// (pull), where applications fetch it. Where it is kept is where it
	// comes from and how long it is there: an inbox's while it holds the
	// job, an archive's for good. Only those send it (Kind.Stocks), and
	// only along their To, to such a node.
	PayloadStocked = "stocked"
)

var payloadKinds = []string{PayloadBDF, PayloadData, PayloadAny, PayloadStocked}

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
	kindStocked
	// allKinds is what a node that takes anything takes: every kind but
	// a stocked file, which only a node that reads one takes.
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
	case PayloadStocked:
		return kindStocked
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

// takes tells the kinds a node of a kind takes at its input (Kind.Takes:
// printers and displays .bdf documents, template nodes data, every other
// node anything). A kind not registered takes anything, so that a node of
// it is not shown broken; a job that reaches it gets a leg of the error.
func takes(kind string) kinds {
	if k := kindNamed(kind); k != nil {
		return k.takes()
	}
	return allKinds
}

// passesOn tells a kind of node that sends on what reaches it (and, an
// entry, what it receives): what it may take at its input is what every
// node after it takes, so that its first edge, on either side, types it
// as the wildcard pin of a visual editor is typed; so does an archive,
// which sends on what it wrote (Kind.PassesOn). A template or a convert
// node sends on the .bdf it makes whatever reaches it (Kind.Makes), and
// what reached it along Original, which accepts types; a display the
// .bdf it shows; the others send nothing on.
func passesOn(kind string) bool {
	k := kindNamed(kind)
	return k != nil && k.PassesOn
}

// stocks tells a kind of node that keeps what reaches it (an inbox, an
// archive): along its To, to a node that takes a stocked file, it sends
// that (Kind.Stocks).
func stocks(kind string) bool {
	k := kindNamed(kind)
	return k != nil && k.Stocks
}

// stockedTo tells whether a node sends a stocked file along its edge to
// x: a node that stocks, along its To, to one that takes it.
func (n *Node) stockedTo(x *Node) bool {
	return x != nil && stocks(n.Kind) && takes(x.Kind).has(PayloadStocked) && slices.Contains(n.To, x.Name)
}

// takesName is the input of a kind of node as the API shows it: bdf,
// data, stocked, or any.
func takesName(kind string) string {
	switch takes(kind) {
	case kindBDF:
		return PayloadBDF
	case kindData:
		return PayloadData
	case kindStocked:
		return PayloadStocked
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
	k := kindNamed(n.Kind)
	switch {
	case k == nil:
		return in
	case k.Role == RoleEntry:
		return in | n.entryKind()
	case k.Makes != "":
		if len(n.Original) > 0 {
			return kindBit(k.Makes) | in
		}
		return kindBit(k.Makes)
	}
	return in
}

// emitsTo tells the kinds a node sends on along its edge to a node (x,
// nil when it is not there): of a template or a convert node, what
// reaches it along Original (as it took it) and the .bdf it makes along
// To; of an inbox or an archive, a stocked file to a node that takes one
// (stockedTo); of the others, what it emits.
func (n *Node) emitsTo(in kinds, to string, x *Node) kinds {
	if n.stockedTo(x) {
		return kindStocked
	}
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
			for to, out := range n.sends(t.reach[n.Name], nodes) {
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
			// what reaches a node goes on along OnError as it is, along
			// Original of a template or a convert node, and along every
			// edge of one that passes it on, but those along which an
			// inbox or an archive sends a stocked file instead
			along := append(slices.Clone(n.Original), n.OnError...)
			if passesOn(n.Kind) {
				along = n.edges()
			}
			if len(along) == 0 {
				continue
			}
			a := takes(n.Kind)
			for _, to := range along {
				if x := nodes[to]; x != nil && !n.stockedTo(x) {
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
		for to, sent := range n.sends(t.reach[n.Name], nodes) {
			x := nodes[to]
			if x == nil {
				continue
			}
			if bad := sent &^ takes(x.Kind); bad != 0 {
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
	c.To, c.Original, c.Members, c.Otherwise, c.OnError = keep(n.To), keep(n.Original), keep(n.Members), keep(n.Otherwise), keep(n.OnError)
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
