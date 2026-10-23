package docflow

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/shibukawa/bdf/converter"
)

// The topology is a graph of nodes in one namespace: entries (endpoints,
// CUPS servers), routing nodes (queues, classes, filters, inboxes),
// transforming nodes (templates, converts) and exits (displays, printers,
// archives, sinks), with edges from queues, endpoints, CUPS servers,
// template and convert nodes and inboxes (To; a template or a convert
// node's Original, too), filters (Rules,
// Otherwise), classes (Members) and displays (Fallback). A job sent to a
// destination follows the edges (route). The graph is checked when a
// node is put: every node it names exists, no path leads back to it, and
// the kinds of payload that reach every node are ones it takes (typed).

// PutNode adds a node or replaces its settings; the settings of other
// kinds are dropped. An endpoint added gets a key, returned once (Key);
// PutNode keeps the key and the templates of an endpoint that is one
// already (RotateKey makes a new key).
func (m *Manager) PutNode(ctx context.Context, in Node) (key string, err error) {
	if !validName(in.Name) {
		return "", fmt.Errorf("%w: node name %q", ErrInvalid, in.Name)
	}
	k := kindNamed(in.Kind)
	if k == nil {
		return "", fmt.Errorf("%w: kind %q: %s", ErrInvalid, in.Kind, strings.Join(kindNames(), ", "))
	}
	if err := checkScope(in.Scope, in.Owner); err != nil {
		return "", err
	}
	n := Node{Name: in.Name, Kind: in.Kind, Title: clip(strings.TrimSpace(in.Title), 255), Description: clip(strings.TrimSpace(in.Description), 1000), Scope: in.Scope, Owner: in.Owner,
		NoIPP: in.NoIPP, OnError: cleanNames(in.OnError)}
	if k.Role == RoleEntry {
		if err := checkEntry(k, in, &n); err != nil {
			return "", err
		}
	}
	if k.Check != nil {
		if err := k.Check(ctx, m, in, &n); err != nil {
			return "", err
		}
	}
	if n.Port != 0 {
		if m.o.NoPorts {
			return "", fmt.Errorf("%w: ports of their own are not available here (serverless): applications call the API instead", ErrInvalid)
		}
		if n.Port < 1024 || n.Port > 65535 {
			return "", fmt.Errorf("%w: port %d: 1024 to 65535", ErrInvalid, n.Port)
		}
		others, err := m.Nodes(ctx)
		if err != nil {
			return "", err
		}
		for _, o := range others {
			if o.Port == n.Port && o.Name != n.Name {
				return "", fmt.Errorf("%w: port %d is %s's", ErrConflict, n.Port, o.Name)
			}
		}
	}
	for _, name := range n.edges() {
		if name == n.Name {
			return "", fmt.Errorf("%w: %s forwards to itself", ErrInvalid, n.Name)
		}
		if _, err := m.kindOf(ctx, name); err != nil {
			return "", fmt.Errorf("%w: %s: no node %q", ErrInvalid, n.Name, name)
		}
	}
	if err := m.nameFree(ctx, n.Name, n.Kind); err != nil {
		return "", err
	}
	if err := m.acyclic(ctx, &n); err != nil {
		return "", err
	}
	if err := m.fromStock(ctx, &n); err != nil {
		return "", err
	}
	if err := m.typed(ctx, &n); err != nil {
		return "", err
	}
	cur, err := m.Node(ctx, n.Name)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return "", err
	}
	var hash string
	if keyed(n.Kind) {
		// a new endpoint, or one that was another kind, gets a key; an
		// endpoint keeps its (and so does a keyed entry of an adapter)
		if cur != nil && cur.Kind == n.Kind && cur.KeyHash != "" {
			hash = cur.KeyHash
		} else {
			key = newKey()
			if hash, err = hashKey(key); err != nil {
				return "", err
			}
		}
	}
	_, err = change(ctx, m.st.db, n.Name, true, func(r *nodeRecord) error {
		if hash == "" && r.Kind == n.Kind && keyed(n.Kind) {
			hash = r.KeyHash // put twice at once
		}
		*r = nodeRecord(n)
		r.KeyHash = hash
		return nil
	})
	if err != nil {
		return "", err
	}
	if cur != nil && cur.Kind == KindTemplate && n.Kind != KindTemplate {
		// no longer a template node: its templates go
		if err := m.removeTemplatesOf(ctx, n.Name); err != nil {
			m.o.Log.Warn("remove the templates of a node", "node", n.Name, "err", err)
		}
	}
	m.changed()
	return key, nil
}

// RotateKey gives an endpoint a new key, returned once; requests with the
// old one are refused from now.
func (m *Manager) RotateKey(ctx context.Context, name string) (string, error) {
	key := newKey()
	hash, err := hashKey(key)
	if err != nil {
		return "", err
	}
	_, err = change(ctx, m.st.db, name, false, func(r *nodeRecord) error {
		if !keyed(r.Kind) {
			return fmt.Errorf("%w: %s is a %s, not an endpoint", ErrInvalid, name, r.Kind)
		}
		r.KeyHash = hash
		return nil
	})
	if err != nil {
		return "", err
	}
	m.o.Log.Info("endpoint key rotated", "endpoint", name)
	return key, nil
}

func cleanNames(names []string) []string {
	var out []string
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n != "" && !slices.Contains(out, n) {
			out = append(out, n)
		}
	}
	return out
}

// formatRE is the shape of a format's name.
var formatRE = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,31}$`)

// extraFormats are the formats the manager takes besides the converter's:
// its own documents, the rasters of driverless printing (which
// print/printconv reads), and the data it reads itself (JSON, YAML).
var extraFormats = []string{"bdf", "pwg-raster", "urf", "json", "yaml"}

// KnownFormats names the formats an entry may take: the converter's
// (those registered: converter/all has every one) and the manager's own.
func KnownFormats() []string {
	out := slices.Clone(extraFormats)
	for _, f := range converter.Formats() {
		out = append(out, f.Name)
	}
	slices.Sort(out)
	return out
}

func cleanFormats(fs []string) []string {
	var out []string
	for _, f := range fs {
		f = strings.ToLower(strings.TrimSpace(f))
		if f != "" && !slices.Contains(out, f) {
			out = append(out, f)
		}
	}
	return out
}

// checkFormats refuses a format that is not known. With no converter
// registered (a program without converter/all), any name of the shape
// goes.
func checkFormats(fs []string) error {
	known := len(converter.Formats()) > 0
	for _, f := range fs {
		if !formatRE.MatchString(f) {
			return fmt.Errorf("%w: format %q", ErrInvalid, f)
		}
		if known && !slices.Contains(extraFormats, f) && converter.Lookup(f) == nil {
			return fmt.Errorf("%w: format %q is not one the converter reads: %s", ErrInvalid, f, strings.Join(KnownFormats(), ", "))
		}
	}
	return nil
}

// Node returns a node.
func (m *Manager) Node(ctx context.Context, name string) (*Node, error) {
	if !validName(name) {
		return nil, fmt.Errorf("%w: node %q", ErrNotFound, name)
	}
	rec, err := m.st.db.Get[nodeRecord](ctx, name)
	if err != nil {
		return nil, mapErr(err)
	}
	n := Node(rec.Data)
	n.Name = name
	return &n, nil
}

// Nodes returns every node, by name.
func (m *Manager) Nodes(ctx context.Context) ([]Node, error) {
	recs, err := list[nodeRecord](ctx, m.st.db)
	if err != nil {
		return nil, err
	}
	out := make([]Node, len(recs))
	for i, r := range recs {
		out[i] = Node(r.Data)
		out[i].Name = r.ID()
	}
	return out, nil
}

// DeleteNode removes a node, with a template node's templates. One
// another node sends to stays; an inbox with items held stays. A
// display's slips and devices stay until it is added again.
func (m *Manager) DeleteNode(ctx context.Context, name string) error {
	rec, err := m.st.db.Get[nodeRecord](ctx, name)
	if err != nil {
		return mapErr(err)
	}
	if err := m.unreferenced(ctx, name); err != nil {
		return err
	}
	if rec.Data.Kind == KindInbox {
		if st, err := m.st.db.Get[inboxState](ctx, name); err == nil && len(st.Data.Held) > 0 {
			return fmt.Errorf("%w: inbox %s holds %d items", ErrConflict, name, len(st.Data.Held))
		}
	}
	if err := m.st.db.Delete[nodeRecord](ctx, name, rec.Revision()); err != nil {
		return mapErr(err)
	}
	m.changed()
	switch rec.Data.Kind {
	case KindTemplate:
		return m.removeTemplatesOf(ctx, name)
	case KindSink:
		return m.dropSink(ctx, name)
	case KindExport:
		return m.dropExport(ctx, name)
	}
	return nil
}

// kindOf tells the kind of a name.
func (m *Manager) kindOf(ctx context.Context, name string) (string, error) {
	n, err := m.Node(ctx, name)
	if err != nil {
		return "", err
	}
	return n.Kind, nil
}

// nameFree refuses a name another kind of node has: the namespace is one,
// so that a destination and an edge name a node without its kind. A
// routing node changes kind; a display or a printer does not.
func (m *Manager) nameFree(ctx context.Context, name, kind string) error {
	have, err := m.kindOf(ctx, name)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if have != kind && !(changes(kind) && changes(have)) {
		return fmt.Errorf("%w: %q is a %s already", ErrConflict, name, have)
	}
	return nil
}

// referrers are the nodes that send to a name: by an edge, or as a
// display's fallback printer.
func (m *Manager) referrers(ctx context.Context, name string) ([]string, error) {
	var out []string
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		if slices.Contains(n.edges(), name) {
			out = append(out, n.Name)
		}
	}
	return out, nil
}

func (m *Manager) unreferenced(ctx context.Context, name string) error {
	refs, err := m.referrers(ctx, name)
	if err != nil {
		return err
	}
	if len(refs) > 0 {
		return fmt.Errorf("%w: %s is sent to by %s", ErrConflict, name, strings.Join(refs, ", "))
	}
	return nil
}

// maxDepth bounds a route, against a cycle that slipped in.
const maxDepth = 16

// acyclic refuses a node from which a path leads back to it.
func (m *Manager) acyclic(ctx context.Context, n *Node) error {
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return err
	}
	edges := map[string][]string{}
	for _, x := range nodes {
		edges[x.Name] = x.edges()
	}
	edges[n.Name] = n.edges()
	seen := map[string]bool{}
	var walk func(from string, path []string) error
	walk = func(from string, path []string) error {
		for _, to := range edges[from] {
			if to == n.Name {
				return fmt.Errorf("%w: a cycle: %s", ErrInvalid, strings.Join(append(path, to), " → "))
			}
			if !seen[to] {
				seen[to] = true
				if err := walk(to, append(path, to)); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return walk(n.Name, []string{n.Name})
}

// TopoNode is a node of the topology as the API shows it: the node (its
// key's hash left out) and what is at it now.
type TopoNode struct {
	Node
	// Addresses are where the node is reached, by ingress (ipp: its IPP
	// printer's path); In is an endpoint's path.
	Addresses map[string]string `json:"addresses,omitempty"`
	In        string            `json:"in,omitempty"`
	// Takes is the kind of payload the node takes at its input (bdf,
	// data or any), and Emits the kinds it sends on along its edges, as
	// inferred from the entries (kinds.go). Accepts are the kinds its
	// input may take now: what it takes, within what every node after it
	// takes (a queue in front of a printer takes .bdf alone; none: nothing
	// fits). Broken are the edges along which the node sends a kind the
	// node after it does not take, by the node's name: left by a change
	// of settings, drawn broken by the editor, and a job gets a leg of
	// the error there.
	// Reaches, of a template or a convert node, is the kinds that reach
	// it: what it sends on along Original.
	Takes   string              `json:"takes,omitempty"`
	Emits   []string            `json:"emits,omitempty"`
	Reaches []string            `json:"reaches,omitempty"`
	Accepts []string            `json:"accepts"`
	Broken  map[string][]string `json:"broken,omitempty"`
	// Templates are the ids of a template node's templates.
	Templates []string `json:"templates,omitempty"`
	// Held counts the items an inbox holds; Open the slips open at a
	// display and Devices its devices; Pending the deliveries waiting at
	// a printer, an archive or a class; Kept the jobs a sink keeps, and
	// Latest the newest of them; Listed the entries an export lists (of
	// every owner).
	Held    int    `json:"held,omitempty"`
	Open    int    `json:"open,omitempty"`
	Devices int    `json:"devices,omitempty"`
	Pending int    `json:"pending,omitempty"`
	Kept    int    `json:"kept,omitempty"`
	Latest  string `json:"latest,omitempty"`
	Listed  int    `json:"listed,omitempty"`
}

// Topology returns every node, by name.
func (m *Manager) Topology(ctx context.Context) ([]TopoNode, error) {
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return nil, err
	}
	tpls, err := list[templateRecord](ctx, m.st.db)
	if err != nil {
		return nil, err
	}
	byNode := map[string][]string{}
	for _, t := range tpls {
		byNode[t.Data.Node] = append(byNode[t.Data.Node], t.Data.ID)
	}
	pending := map[string]int{}
	if ds, err := m.Pending(ctx); err == nil {
		for _, d := range ds {
			pending[d.Node]++
		}
	}
	byName := make(map[string]*Node, len(nodes))
	for i := range nodes {
		byName[nodes[i].Name] = &nodes[i]
	}
	ins := ingresses()
	ty := inferKinds(byName)
	acc := accepts(byName)
	broken := ty.broken(byName)
	out := make([]TopoNode, 0, len(nodes))
	for _, n := range nodes {
		n.KeyHash = ""
		t := TopoNode{Node: n, Takes: takesName(n.Kind), Accepts: acc[n.Name].names(), Pending: pending[n.Name]}
		for _, in := range ins {
			if a := in.Address(&n); a != "" {
				if t.Addresses == nil {
					t.Addresses = map[string]string{}
				}
				t.Addresses[in.Name] = a
			}
		}
		if t.Accepts == nil {
			t.Accepts = []string{}
		}
		if k := kindNamed(n.Kind); k != nil && (k.PassesOn || k.Makes != "") {
			t.Emits = n.emits(ty.reach[n.Name]).names()
		}
		if n.Kind == KindTemplate || n.Kind == KindConvert {
			t.Reaches = ty.reach[n.Name].names()
		}
		for to, bad := range broken[n.Name] {
			if t.Broken == nil {
				t.Broken = map[string][]string{}
			}
			t.Broken[to] = bad.names()
		}
		if keyed(n.Kind) {
			t.In = InPath + "/" + n.Name
		}
		switch n.Kind {
		case KindArchive:
			t.Collection = n.collection()
		case KindTemplate:
			t.Templates = byNode[n.Name]
			slices.Sort(t.Templates)
		case KindInbox:
			if st, err := m.st.db.Get[inboxState](ctx, n.Name); err == nil {
				t.Held = len(st.Data.Held)
			}
		case KindDisplay:
			if st, err := m.st.db.Get[locationState](ctx, n.Name); err == nil {
				t.Open, t.Devices = len(st.Data.Open), len(st.Data.Devices)
			}
		case KindSink:
			if ids := m.sinkJobs(ctx, n.Name); len(ids) > 0 {
				t.Kept, t.Latest = len(ids), ids[len(ids)-1]
			}
		case KindExport:
			t.Listed = len(m.exportJobs(ctx, exportListID(n.Name, "")))
		}
		out = append(out, t)
	}
	slices.SortFunc(out, func(a, b TopoNode) int { return strings.Compare(a.Name, b.Name) })
	return out, nil
}

// changes tells a kind a node may change from and to, keeping its name:
// one registered that is not Fixed.
func changes(kind string) bool {
	k := kindNamed(kind)
	return k != nil && !k.Fixed
}
