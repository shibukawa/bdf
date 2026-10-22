package manager

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/shibukawa/bdf/converter"
	"github.com/shibukawa/yuruidb"
)

// The topology is a graph of nodes in one namespace: entries (endpoints,
// CUPS servers), routing nodes (queues, classes, filters, inboxes),
// transforming nodes (templates, converts) and exits (displays, printers,
// archives, sinks), with edges from queues, endpoints, CUPS servers,
// template and convert nodes and inboxes (To), filters (Rules,
// Otherwise), classes (Members) and displays (Fallback). A job sent to a
// destination follows the edges (route). The graph is checked when a
// node is put: every node it names exists, no path leads back to it, and
// the kinds of payload that reach every node are ones it takes (typed).

// nodeKinds are the kinds of a node, in the order of a job's way.
var nodeKinds = []string{KindEndpoint, KindCups, KindQueue, KindClass, KindFilter, KindTemplate, KindConvert, KindInbox, KindArchive, KindDisplay, KindPrinter, KindSink}

// routingKinds are the kinds a node changes between, keeping its name. A
// display, a printer or a sink keeps its kind: its slips, devices,
// deliveries or kept jobs name it.
var routingKinds = []string{KindEndpoint, KindCups, KindQueue, KindClass, KindFilter, KindTemplate, KindConvert, KindInbox, KindArchive}

// PutNode adds a node or replaces its settings; the settings of other
// kinds are dropped. An endpoint added gets a key, returned once (Key);
// PutNode keeps the key and the templates of an endpoint that is one
// already (RotateKey makes a new key).
func (m *Manager) PutNode(ctx context.Context, in Node) (key string, err error) {
	if !validName(in.Name) {
		return "", fmt.Errorf("%w: node name %q", ErrInvalid, in.Name)
	}
	if !slices.Contains(nodeKinds, in.Kind) {
		return "", fmt.Errorf("%w: kind %q: %s", ErrInvalid, in.Kind, strings.Join(nodeKinds, ", "))
	}
	if err := checkScope(in.Scope, in.Owner); err != nil {
		return "", err
	}
	n := Node{Name: in.Name, Kind: in.Kind, Title: clip(strings.TrimSpace(in.Title), 255), Description: clip(strings.TrimSpace(in.Description), 1000), Scope: in.Scope, Owner: in.Owner}
	switch n.Kind {
	case KindQueue:
		n.To = cleanNames(in.To)
		if len(n.To) == 0 {
			return "", fmt.Errorf("%w: a queue forwards to at least one node (to)", ErrInvalid)
		}
	case KindEndpoint, KindCups:
		n.To = cleanNames(in.To)
		if len(n.To) == 0 {
			return "", fmt.Errorf("%w: a %s forwards what it receives to at least one node (to)", ErrInvalid, kindName(n.Kind))
		}
		n.Accept = or(in.Accept, AcceptFile)
		if n.Accept != AcceptFile && n.Accept != AcceptData {
			return "", fmt.Errorf("%w: accept %q: file or data", ErrInvalid, in.Accept)
		}
		n.Output = or(in.Output, OutputBDF)
		if n.Output != OutputBDF && n.Output != OutputFile {
			return "", fmt.Errorf("%w: output %q: bdf (the document made of the file) or file (as received)", ErrInvalid, in.Output)
		}
		if n.Accept == AcceptData {
			n.Output = "" // data goes on as it is
		}
		n.Formats = cleanFormats(in.Formats)
		if err := checkFormats(n.Formats); err != nil {
			return "", err
		}
		if n.Meta, err = claimedMeta(in.Meta); err != nil {
			return "", err
		}
		n.Port = in.Port
		if n.Kind == KindEndpoint {
			n.Apps = cleanNames(in.Apps)
			for _, id := range n.Apps {
				if _, err := m.App(ctx, id); err != nil {
					return "", fmt.Errorf("%w: no app %q", ErrInvalid, id)
				}
			}
		} else {
			n.TLSCert, n.TLSKey = strings.TrimSpace(in.TLSCert), strings.TrimSpace(in.TLSKey)
			if n.Port == 0 {
				return "", fmt.Errorf("%w: a CUPS server listens on a port of its own (port)", ErrInvalid)
			}
			if (n.TLSCert == "") != (n.TLSKey == "") {
				return "", fmt.Errorf("%w: a CUPS server's certificate takes both tls_cert and tls_key (files)", ErrInvalid)
			}
		}
	case KindTemplate:
		n.To = cleanNames(in.To)
		if len(n.To) == 0 {
			return "", fmt.Errorf("%w: a template node sends the form it fills to at least one node (to)", ErrInvalid)
		}
		n.Template, n.TemplateFrom = strings.TrimSpace(in.Template), strings.TrimSpace(in.TemplateFrom)
		if (n.Template == "") == (n.TemplateFrom == "") {
			return "", fmt.Errorf("%w: a template node names the template it fills (template), or the metadata key that names one (template_from)", ErrInvalid)
		}
		if n.Template != "" && !templateIDRE.MatchString(n.Template) {
			return "", fmt.Errorf("%w: template id %q: letters, digits, _ and -", ErrInvalid, n.Template)
		}
		n.ListName, n.Strict = strings.TrimSpace(in.ListName), in.Strict
	case KindConvert:
		n.To = cleanNames(in.To)
		if len(n.To) == 0 {
			return "", fmt.Errorf("%w: a convert node sends the document it makes to at least one node (to)", ErrInvalid)
		}
		n.Formats = cleanFormats(in.Formats)
		if err := checkFormats(n.Formats); err != nil {
			return "", err
		}
	case KindSink:
		n.Latest = in.Latest
		if n.Latest < 0 || n.Latest > 1000 {
			return "", fmt.Errorf("%w: latest %d: a sink keeps 1 to 1000 documents", ErrInvalid, in.Latest)
		}
	case KindClass:
		n.Members = cleanNames(in.Members)
		if len(n.Members) == 0 {
			return "", fmt.Errorf("%w: a class has at least one member", ErrInvalid)
		}
		n.Policy = or(in.Policy, PolicyRoundRobin)
		if n.Policy != PolicyRoundRobin && n.Policy != PolicyFailover {
			return "", fmt.Errorf("%w: policy %q: round-robin or failover", ErrInvalid, in.Policy)
		}
	case KindFilter:
		rules, err := checkRules(in.Rules)
		if err != nil {
			return "", err
		}
		n.Rules, n.Otherwise = rules, cleanNames(in.Otherwise)
		if len(n.Rules) == 0 && len(n.Otherwise) == 0 {
			return "", fmt.Errorf("%w: a filter has rules, or nodes it sends to otherwise", ErrInvalid)
		}
	case KindInbox:
		n.To, n.Keep = cleanNames(in.To), in.Keep
	case KindArchive:
		n.Store, n.Collection = strings.TrimSpace(in.Store), strings.TrimSpace(in.Collection)
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
	case KindDisplay:
		n.Keep, n.OfflineAfter = in.Keep, in.OfflineAfter
		if in.Fallback != nil && strings.TrimSpace(in.Fallback.Printer) != "" {
			n.Fallback = &Fallback{Printer: strings.TrimSpace(in.Fallback.Printer), After: in.Fallback.After}
			if k, err := m.kindOf(ctx, n.Fallback.Printer); err != nil || k != KindPrinter {
				return "", fmt.Errorf("%w: fallback printer %q is not a printer", ErrInvalid, n.Fallback.Printer)
			}
		}
	case KindPrinter:
		n.URL, n.Route = strings.TrimSpace(in.URL), or(in.Route, RouteDirect)
		if n.Route != RouteDirect && n.Route != RouteBranch {
			return "", fmt.Errorf("%w: route %q: direct or branch", ErrInvalid, in.Route)
		}
		if err := checkPrinterURL(n.URL); err != nil {
			return "", err
		}
		n.Preview = strings.TrimSpace(in.Preview)
		if err := checkPreviewURL(n.Preview); err != nil {
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
	if err := m.typed(ctx, &n); err != nil {
		return "", err
	}
	cur, err := m.Node(ctx, n.Name)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return "", err
	}
	var hash string
	if n.Kind == KindEndpoint {
		// a new endpoint, or one that was another kind, gets a key; an
		// endpoint keeps its
		if cur != nil && cur.Kind == KindEndpoint && cur.KeyHash != "" {
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
	if have != kind && !(slices.Contains(routingKinds, kind) && slices.Contains(routingKinds, have)) {
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
	// IPP is the destination's path; In an endpoint's.
	IPP string `json:"ipp"`
	In  string `json:"in,omitempty"`
	// Takes is the kind of payload the node takes at its input (bdf,
	// data or any), and Emits the kinds it sends on along its edges, as
	// inferred from the entries (kinds.go). Accepts are the kinds its
	// input may take now: what it takes, within what every node after it
	// takes (a queue in front of a printer takes .bdf alone; none: nothing
	// fits). Broken are the edges along which the node sends a kind the
	// node after it does not take, by the node's name: left by a change
	// of settings, drawn broken by the editor, and a job gets a leg of
	// the error there.
	Takes   string              `json:"takes,omitempty"`
	Emits   []string            `json:"emits,omitempty"`
	Accepts []string            `json:"accepts"`
	Broken  map[string][]string `json:"broken,omitempty"`
	// Templates are the ids of a template node's templates.
	Templates []string `json:"templates,omitempty"`
	// Held counts the items an inbox holds; Open the slips open at a
	// display and Devices its devices; Pending the deliveries waiting at
	// a printer, an archive or a class; Kept the jobs a sink keeps, and
	// Latest the newest of them.
	Held    int    `json:"held,omitempty"`
	Open    int    `json:"open,omitempty"`
	Devices int    `json:"devices,omitempty"`
	Pending int    `json:"pending,omitempty"`
	Kept    int    `json:"kept,omitempty"`
	Latest  string `json:"latest,omitempty"`
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
	ty := inferKinds(byName)
	acc := accepts(byName)
	broken := ty.broken(byName)
	out := make([]TopoNode, 0, len(nodes))
	for _, n := range nodes {
		n.KeyHash = ""
		t := TopoNode{Node: n, IPP: IPPPath + "/" + n.Name, Takes: takesName(n.Kind), Accepts: acc[n.Name].names()}
		if t.Accepts == nil {
			t.Accepts = []string{}
		}
		if passesOn(n.Kind) || n.Kind == KindTemplate || n.Kind == KindConvert || n.Kind == KindDisplay {
			t.Emits = n.emits(ty.reach[n.Name]).names()
		}
		for to, bad := range broken[n.Name] {
			if t.Broken == nil {
				t.Broken = map[string][]string{}
			}
			t.Broken[to] = bad.names()
		}
		switch n.Kind {
		case KindArchive:
			t.Collection = n.collection()
			t.Pending = pending[n.Name]
		case KindEndpoint:
			t.In = InPath + "/" + n.Name
		case KindTemplate:
			t.Templates = byNode[n.Name]
			slices.Sort(t.Templates)
			t.Pending = pending[n.Name]
		case KindConvert:
			t.Pending = pending[n.Name]
		case KindInbox:
			if st, err := m.st.db.Get[inboxState](ctx, n.Name); err == nil {
				t.Held = len(st.Data.Held)
			}
		case KindDisplay:
			if st, err := m.st.db.Get[locationState](ctx, n.Name); err == nil {
				t.Open, t.Devices = len(st.Data.Open), len(st.Data.Devices)
			}
		case KindPrinter, KindClass:
			t.Pending = pending[n.Name]
		case KindSink:
			if ids := m.sinkJobs(ctx, n.Name); len(ids) > 0 {
				t.Kept, t.Latest = len(ids), ids[len(ids)-1]
			}
		}
		out = append(out, t)
	}
	slices.SortFunc(out, func(a, b TopoNode) int { return strings.Compare(a.Name, b.Name) })
	return out, nil
}

// kindName is a kind as it is said: "CUPS server" for cups.
func kindName(kind string) string {
	if kind == KindCups {
		return "CUPS server"
	}
	return kind
}
