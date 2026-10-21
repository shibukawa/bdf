package manager

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/shibukawa/yuruidb"
)

// The topology is a graph: locations, printers and routing nodes in one
// namespace, with edges from queues and inboxes (To) and classes
// (Members). A job sent to a destination follows the edges (route). The
// graph is checked when a node is put: every node it names exists, and
// no path leads back to it.

// nodeKinds are the kinds of a Node.
var nodeKinds = []string{KindQueue, KindClass, KindInbox, KindArchive, KindEndpoint}

// PutNode adds a routing node or replaces its settings. An endpoint
// added gets a key, returned once (Key); PutNode keeps the key of an
// endpoint that has one (RotateKey makes a new one).
func (m *Manager) PutNode(ctx context.Context, n Node) (key string, err error) {
	if !validName(n.Name) {
		return "", fmt.Errorf("%w: node name %q", ErrInvalid, n.Name)
	}
	if !slices.Contains(nodeKinds, n.Kind) {
		return "", fmt.Errorf("%w: kind %q: %s", ErrInvalid, n.Kind, strings.Join(nodeKinds, ", "))
	}
	if err := checkScope(n.Scope, n.Owner); err != nil {
		return "", err
	}
	n.To = cleanNames(n.To)
	n.Members = cleanNames(n.Members)
	n.Templates = cleanNames(n.Templates)
	n.KeyHash = ""
	switch n.Kind {
	case KindQueue, KindEndpoint:
		if len(n.To) == 0 {
			return "", fmt.Errorf("%w: a %s forwards to at least one node (to)", ErrInvalid, n.Kind)
		}
		n.Members, n.Policy, n.Store, n.Collection, n.Keep = nil, "", "", "", 0
		if n.Kind == KindQueue {
			n.Templates = nil
		}
		for _, id := range n.Templates {
			if _, err := m.Template(ctx, id); err != nil {
				return "", fmt.Errorf("%w: no template %q", ErrInvalid, id)
			}
		}
	case KindClass:
		if len(n.Members) == 0 {
			return "", fmt.Errorf("%w: a class has at least one member", ErrInvalid)
		}
		if n.Policy == "" {
			n.Policy = PolicyRoundRobin
		}
		if n.Policy != PolicyRoundRobin && n.Policy != PolicyFailover {
			return "", fmt.Errorf("%w: policy %q: round-robin or failover", ErrInvalid, n.Policy)
		}
		n.To, n.Store, n.Collection, n.Keep, n.Templates = nil, "", "", 0, nil
	case KindInbox:
		n.Members, n.Policy, n.Store, n.Collection, n.Templates = nil, "", "", "", nil
	case KindArchive:
		if n.Store == "" {
			return "", fmt.Errorf("%w: an archive has a store (bbolt, s3, r2 or gcs)", ErrInvalid)
		}
		if n.Collection != "" {
			if err := yuruidb.ValidateCollection(n.Collection); err != nil {
				return "", fmt.Errorf("%w: collection %q: letters, digits, _ and -", ErrInvalid, n.Collection)
			}
		}
		// the store is opened now, so that a wrong URL is refused here
		if _, err := m.st.archive(ctx, n.Store, n.collection()); err != nil {
			return "", fmt.Errorf("%w: store: %v", ErrInvalid, err)
		}
		n.To, n.Members, n.Policy, n.Keep, n.Templates = nil, nil, "", 0, nil
	}
	for _, name := range append(slices.Clone(n.To), n.Members...) {
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
	var hash string
	if n.Kind == KindEndpoint {
		// a new endpoint, or one that was another kind, gets a key; an
		// endpoint keeps its
		if cur, err := m.Node(ctx, n.Name); err == nil && cur.Kind == KindEndpoint && cur.KeyHash != "" {
			hash = cur.KeyHash
		} else {
			key = newKey()
			if hash, err = hashKey(key); err != nil {
				return "", err
			}
		}
	}
	_, err = change(ctx, m.st.db, n.Name, true, func(r *nodeRecord) error {
		if hash == "" && r.Kind == KindEndpoint && n.Kind == KindEndpoint {
			hash = r.KeyHash // put twice at once
		}
		*r = nodeRecord(n)
		r.KeyHash = hash
		return nil
	})
	if err != nil {
		return "", err
	}
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
		if r.Kind != KindEndpoint {
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

// Node returns a routing node.
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

// Nodes returns the routing nodes, by name.
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

// DeleteNode removes a routing node. One another node forwards to stays;
// an inbox with items held stays.
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
	return mapErr(m.st.db.Delete[nodeRecord](ctx, name, rec.Revision()))
}

// kindOf tells the kind of a name: a location (display), a printer, or
// a routing node's kind.
func (m *Manager) kindOf(ctx context.Context, name string) (string, error) {
	if !validName(name) {
		return "", fmt.Errorf("%w: %q", ErrNotFound, name)
	}
	if _, err := m.st.db.Get[locationRecord](ctx, name); err == nil {
		return KindDisplay, nil
	} else if !errors.Is(err, yuruidb.ErrNotFound) {
		return "", mapErr(err)
	}
	if _, err := m.st.db.Get[printerRecord](ctx, name); err == nil {
		return KindPrinter, nil
	} else if !errors.Is(err, yuruidb.ErrNotFound) {
		return "", mapErr(err)
	}
	rec, err := m.st.db.Get[nodeRecord](ctx, name)
	if err != nil {
		return "", mapErr(err)
	}
	return rec.Data.Kind, nil
}

// nameFree refuses a name another kind of node has: the namespace is one,
// so that a destination and an edge name a node without its kind.
func (m *Manager) nameFree(ctx context.Context, name, kind string) error {
	have, err := m.kindOf(ctx, name)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	same := have == kind
	if slices.Contains(nodeKinds, kind) && slices.Contains(nodeKinds, have) {
		same = true // a node changes kind
	}
	if !same {
		return fmt.Errorf("%w: %q is a %s already", ErrConflict, name, have)
	}
	return nil
}

// referrers are the nodes and locations that send to a name.
func (m *Manager) referrers(ctx context.Context, name string) ([]string, error) {
	var out []string
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		if slices.Contains(n.To, name) || slices.Contains(n.Members, name) {
			out = append(out, n.Name)
		}
	}
	locs, err := m.Locations(ctx)
	if err != nil {
		return nil, err
	}
	for _, l := range locs {
		if l.Fallback != nil && l.Fallback.Printer == name {
			out = append(out, l.Name)
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
		edges[x.Name] = append(slices.Clone(x.To), x.Members...)
	}
	edges[n.Name] = append(slices.Clone(n.To), n.Members...)
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

// TopoNode is a node of the topology as the API shows it: a location, a
// printer or a routing node, with the edges from it.
type TopoNode struct {
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Title string `json:"title,omitempty"`
	// IPP is the destination's path.
	IPP string `json:"ipp"`
	// To and Members are the edges; Policy a class's.
	To      []string `json:"to,omitempty"`
	Members []string `json:"members,omitempty"`
	Policy  string   `json:"policy,omitempty"`
	// Fallback is a location's printer; URL and Route a printer's.
	Fallback *Fallback `json:"fallback,omitempty"`
	URL      string    `json:"url,omitempty"`
	Route    string    `json:"route,omitempty"`
	// Store and Collection are an archive's; Keep an inbox's.
	Store      string   `json:"store,omitempty"`
	Collection string   `json:"collection,omitempty"`
	Keep       Duration `json:"keep,omitempty"`
	// Templates are an endpoint's, and In its path.
	Templates []string `json:"templates,omitempty"`
	In        string   `json:"in,omitempty"`
	// Scope and Owner are who may send to the node (Node.Scope).
	Scope string `json:"scope,omitempty"`
	Owner string `json:"owner,omitempty"`
	// Held counts the items an inbox holds.
	Held int `json:"held,omitempty"`
}

// Topology returns every node, by name.
func (m *Manager) Topology(ctx context.Context) ([]TopoNode, error) {
	var out []TopoNode
	locs, err := m.Locations(ctx)
	if err != nil {
		return nil, err
	}
	for _, l := range locs {
		out = append(out, TopoNode{Name: l.Name, Kind: KindDisplay, Title: l.Title, IPP: IPPPath + "/" + l.Name, Fallback: l.Fallback, Keep: l.Keep,
			Scope: l.Scope, Owner: l.Owner})
	}
	printers, err := m.Printers(ctx)
	if err != nil {
		return nil, err
	}
	for _, p := range printers {
		out = append(out, TopoNode{Name: p.Name, Kind: KindPrinter, Title: p.Info, IPP: IPPPath + "/" + p.Name, URL: p.URL, Route: p.Route,
			Scope: p.Scope, Owner: p.Owner})
	}
	nodes, err := m.Nodes(ctx)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		t := TopoNode{Name: n.Name, Kind: n.Kind, Title: n.Title, IPP: IPPPath + "/" + n.Name, To: n.To, Members: n.Members, Policy: n.Policy,
			Store: n.Store, Keep: n.Keep, Templates: n.Templates, Scope: n.Scope, Owner: n.Owner}
		if n.Kind == KindArchive {
			t.Collection = n.collection()
		}
		if n.Kind == KindEndpoint {
			t.In = InPath + "/" + n.Name
		}
		if n.Kind == KindInbox {
			if st, err := m.st.db.Get[inboxState](ctx, n.Name); err == nil {
				t.Held = len(st.Data.Held)
			}
		}
		out = append(out, t)
	}
	slices.SortFunc(out, func(a, b TopoNode) int { return strings.Compare(a.Name, b.Name) })
	return out, nil
}
