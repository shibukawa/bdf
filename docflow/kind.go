package docflow

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/shibukawa/bdf"
)

// The kinds of node are registered, as the drivers of database/sql are:
// the core's (queues, classes, filters, inboxes, archives, template and
// convert nodes, displays, sinks, exports and endpoints) by this package,
// and the inputs and outputs by the packages of docflow/io, in their init
// (docflow/io/ipp: printers, CUPS servers, and every node's IPP printer at
// /ipp/print). A program imports the ones it serves. The engine knows a
// kind by its entry here alone: what it takes and sends on, how its
// settings are checked, and what a job that reaches it does. The settings
// of every kind are the fields of Node.

// Role is what a kind of node does with a job, as the editor colors it.
type Role string

const (
	// RoleEntry takes jobs from outside the topology and forwards them to
	// every node of To.
	RoleEntry Role = "entry"
	// RoleRoute sends on what reaches it.
	RoleRoute Role = "route"
	// RoleStep sends on a payload it makes of what reaches it.
	RoleStep Role = "step"
	// RoleExit is where a job ends, or leaves the topology.
	RoleExit Role = "exit"
)

// Kind is a kind of node.
type Kind struct {
	Name string
	// Title is the kind as it is said ("CUPS server"); Name by default.
	Title string
	Role  Role
	// Takes are the kinds of payload its input takes (PayloadBDF,
	// PayloadData, PayloadAny); none, every one.
	Takes []string
	// PassesOn tells a kind that sends on what reaches it (a router, an
	// entry, an archive): its input takes what the nodes after it take,
	// and it is typed by its first edge.
	PassesOn bool
	// Makes is the kind of payload it sends on whatever reaches it (a
	// template node's .bdf); none, what reaches it.
	Makes string
	// Fixed keeps a node of the kind from becoming another kind, and
	// another from becoming one of it: its records name it (a display's
	// slips, a printer's deliveries).
	Fixed bool
	// Check checks the settings of a node of the kind and copies them
	// from in to n, which PutNode made of the settings every kind has
	// (and, of an entry, those of every entry). An error of ErrInvalid is
	// the caller's.
	Check func(ctx context.Context, m *Manager, in Node, n *Node) error
	// Ready tells whether the engine can send to a node of an exit now
	// (a printer a branch fetches from: not yet). Its error stops a send
	// for good.
	Ready func(n *Node) error
	// Send sends a payload out of the topology (an exit): a job that
	// reaches the node becomes a delivery, done in the background at
	// least once and tried again after a failure (Final: never again).
	// It returns what to record of the send (the printer's job id).
	Send func(ctx context.Context, o *Out) (result string, err error)
	// Serve serves a node of an entry on its own port (Node.Port).
	Serve func(m *Manager, n *Node) http.Handler
	// Step is the step of a node of the kind (a render or a notify
	// node): a job that reaches the node becomes a delivery, done in the
	// background at least once as Send is, and the job goes on from the
	// node's To with what Step made, or, with nil, with what reached it.
	Step func(ctx context.Context, o *Out) (*Made, error)
	// Keyed gives a node of an entry a key, as an endpoint has: POST
	// /in/<name> with it sends to the node (an S3 Trigger).
	Keyed bool
	// Fetch takes what a request to a keyed entry carries, a reference
	// (an S3 event), and reads the document it names; none, the body is
	// the document.
	Fetch func(ctx context.Context, o *Fetch) (*Fetched, error)

	// route takes a job that reaches a node of a core kind; a kind with
	// none sends to an exit (Send), or forwards along To.
	route func(m *Manager, ctx context.Context, r *routing, name, from string, depth int)
}

func (k *Kind) title() string { return or(k.Title, k.Name) }

// takes is the set of the payload kinds its input takes.
func (k *Kind) takes() kinds {
	if len(k.Takes) == 0 {
		return allKinds
	}
	var out kinds
	for _, t := range k.Takes {
		out |= kindBit(t)
	}
	return out
}

var registry = struct {
	sync.RWMutex
	kinds     []*Kind
	ingresses []*Ingress
}{}

// RegisterKind makes a kind of node known. It is called once for a kind,
// in an init; a name registered twice panics.
func RegisterKind(k Kind) {
	if k.Name == "" || !validName(k.Name) {
		panic(fmt.Sprintf("docflow: kind name %q", k.Name))
	}
	if k.Makes != "" && kindBit(k.Makes) == 0 {
		panic(fmt.Sprintf("docflow: kind %s makes %q", k.Name, k.Makes))
	}
	registry.Lock()
	defer registry.Unlock()
	for _, x := range registry.kinds {
		if x.Name == k.Name {
			panic("docflow: kind " + k.Name + " registered twice")
		}
	}
	registry.kinds = append(registry.kinds, &k)
}

// Kinds returns the kinds registered, in the order of a job's way:
// entries, routers, steps, exits.
func Kinds() []Kind {
	registry.RLock()
	defer registry.RUnlock()
	out := make([]Kind, 0, len(registry.kinds))
	for _, role := range []Role{RoleEntry, RoleRoute, RoleStep, RoleExit} {
		for _, k := range registry.kinds {
			if k.Role == role {
				out = append(out, *k)
			}
		}
	}
	return out
}

// kindNamed returns a registered kind, or nil.
func kindNamed(name string) *Kind {
	registry.RLock()
	defer registry.RUnlock()
	for _, k := range registry.kinds {
		if k.Name == name {
			return k
		}
	}
	return nil
}

// kindNames are the names of the kinds registered, in order.
func kindNames() []string {
	var out []string
	for _, k := range Kinds() {
		out = append(out, k.Name)
	}
	return out
}

// isEntry tells a kind of entry (Kind.Role).
func isEntry(kind string) bool {
	k := kindNamed(kind)
	return k != nil && k.Role == RoleEntry
}

// ServesPort tells a kind whose node listens on a port of its own
// (Node.Port, PortHandler): an endpoint, and an entry of an adapter that
// serves one (Kind.Serve: a CUPS server).
func ServesPort(kind string) bool {
	if kind == KindEndpoint {
		return true
	}
	k := kindNamed(kind)
	return k != nil && k.Serve != nil
}

// keyed tells a kind of entry whose node has a key that sends to it at
// POST /in/<name>: an endpoint, and a keyed entry of an adapter.
func keyed(kind string) bool {
	if kind == KindEndpoint {
		return true
	}
	k := kindNamed(kind)
	return k != nil && k.Keyed && k.Role == RoleEntry
}

// steps tells a kind whose node does a step of an adapter's (Kind.Step).
func steps(kind string) bool {
	k := kindNamed(kind)
	return k != nil && k.Step != nil
}

// sendsOut tells a kind that sends what reaches it out of the topology:
// the member of a class, the fallback of a display (with .bdf).
func sendsOut(kind string) bool {
	k := kindNamed(kind)
	return k != nil && k.Send != nil
}

// Ingress is a way into the topology besides the entries: a protocol on
// which every node (or every one that does not turn it off) is a
// destination, as every node is an IPP printer at /ipp/print/<name>.
type Ingress struct {
	Name string
	// Prefix is the path its requests are under, which the web
	// application mounts Handler at.
	Prefix string
	// Mount adds its routes to the engine's handler.
	Mount func(m *Manager, mux *http.ServeMux)
	// Address is where a node is reached by it; "" none.
	Address func(n *Node) string
}

// RegisterIngress makes an ingress known, once, in an init.
func RegisterIngress(in Ingress) {
	registry.Lock()
	defer registry.Unlock()
	for _, x := range registry.ingresses {
		if x.Name == in.Name {
			panic("docflow: ingress " + in.Name + " registered twice")
		}
	}
	registry.ingresses = append(registry.ingresses, &in)
}

func ingresses() []*Ingress {
	registry.RLock()
	defer registry.RUnlock()
	return slices.Clone(registry.ingresses)
}

// Prefixes are the paths of the ingresses registered, for the web
// application to mount Handler at.
func Prefixes() []string {
	var out []string
	for _, in := range ingresses() {
		out = append(out, in.Prefix)
	}
	return out
}

// Out is a payload an exit sends out of the topology, or a step takes.
type Out struct {
	// ID is the delivery's: <job>.<leg>, the same on every try, the key
	// a far side tells a second send of the same by (a Message-ID, an
	// Idempotency-Key).
	ID      string
	Node    *Node
	Job     *Job
	Payload *Payload
	// Name is the job's name for the far side (the job's, or the slip's).
	Name string
	// Meta is the job's metadata, with $error and $error_node when the
	// job came along an on_error edge (Cause).
	Meta map[string]any
	// Cause is the failure that sent the job here along an on_error
	// edge, or nil.
	Cause *Failure
	m     *Manager
}

// Failure is a send or a step that failed for good, which the job left
// along an on_error edge.
type Failure struct {
	Node  string `json:"node"`
	Error string `json:"error"`
}

// Made is a payload a step made: Body is read once, into the job's
// documents.
type Made struct {
	Kind   string
	Format string
	Name   string
	Pages  int
	Body   io.Reader
}

// Fetch is a request to a keyed entry that carries a reference (Kind.Fetch).
type Fetch struct {
	Node *Node
	// Body is the request's body, of 1 MiB at most.
	Body []byte
	m    *Manager
}

// Manager is the engine the request came to.
func (f *Fetch) Manager() *Manager { return f.m }

// Fetched is the document a reference names: its name, the metadata it
// comes with (as a sender's), and its bytes, read once and closed.
type Fetched struct {
	Name string
	Meta map[string]any
	Body io.ReadCloser
}

// Document reads the payload as a .bdf document.
func (o *Out) Document(ctx context.Context) (*bdf.Document, error) {
	return o.m.Document(ctx, o.Payload.ID)
}

// Open reads the payload's bytes.
func (o *Out) Open(ctx context.Context) (io.ReadCloser, error) {
	r, err := o.m.st.docs.Open(ctx, o.Payload.ID)
	return r, mapErr(err)
}

// Manager is the engine the payload is sent by.
func (o *Out) Manager() *Manager { return o.m }

// view is the metadata as the templates of a node see it: metaView,
// with $node, $delivery and $name (the payload's name, else the job's),
// and $error and $error_node of a job that
// came along an on_error edge.
func (o *Out) view() map[string]string {
	j := o.Job
	if j == nil {
		j = &Job{Name: o.Name} // a slip whose job is gone
	}
	v := metaView(j, o.Payload)
	for k, x := range o.Meta {
		if _, ok := v[k]; !ok {
			v[k] = metaText(x)
		}
	}
	v["$node"], v["$delivery"] = o.Node.Name, o.ID
	if o.Payload != nil {
		v["$name"] = firstOr(o.Payload.Name, o.Name)
	}
	if o.Cause != nil {
		v["$error"], v["$error_node"] = o.Cause.Error, o.Cause.Node
	}
	return v
}

// Expand replaces every {{meta.<key>}} in s with the value of the key
// (Out.view): a sender's key, name and user, or one of the manager's ($job,
// $file, $node, $error …); a key not there is empty. There is no other
// expression.
func (o *Out) Expand(s string) string { return expand(s, o.view()) }

// ExpandFunc expands s as Expand does, with every value escaped by esc
// first (url.QueryEscape, for a URL).
func (o *Out) ExpandFunc(s string, esc func(string) string) string {
	v := o.view()
	for k, x := range v {
		v[k] = esc(x)
	}
	return expand(s, v)
}

// ExpandHeader expands s for a header of a mail or a request: without
// line breaks or other control characters, so that a value cannot add a
// header.
func (o *Out) ExpandHeader(s string) string {
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return ' '
		}
		return r
	}, o.Expand(s)))
}

// ExpandPath expands s for a path (an object's key, a file's name): every
// value is cleaned of control characters, backslashes and .. segments,
// and the path of a leading slash, so that a value cannot leave the
// prefix it is put under.
func (o *Out) ExpandPath(s string) string {
	v := o.view()
	for k, x := range v {
		v[k] = cleanPathPart(x)
	}
	return CleanPath(expand(s, v))
}

var expandRE = regexp.MustCompile(`\{\{\s*meta\.([^{}\s]+)\s*\}\}`)

func expand(s string, view map[string]string) string {
	return expandRE.ReplaceAllStringFunc(s, func(m string) string {
		return view[expandRE.FindStringSubmatch(m)[1]]
	})
}

// cleanPathPart cleans a value put into a path: no control characters,
// no backslashes, no .. or . segment.
func cleanPathPart(s string) string {
	s = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f || r == '\\' {
			return -1
		}
		return r
	}, s)
	parts := strings.Split(s, "/")
	parts = slices.DeleteFunc(parts, func(p string) bool { return p == ".." || p == "." })
	return strings.Join(parts, "/")
}

// CleanPath makes a path relative and clean: no leading slash, no empty,
// . or .. segment.
func CleanPath(s string) string {
	parts := slices.DeleteFunc(strings.Split(cleanPathPart(s), "/"), func(p string) bool { return p == "" })
	return strings.Join(parts, "/")
}

// ErrUnsupported is the error of what a node cannot do yet (a printer a
// branch fetches from).
var ErrUnsupported = errors.New("docflow: not supported")

// Final marks an error no retry mends: a send that fails with it fails
// for good.
func Final(err error) error {
	if err == nil {
		return nil
	}
	return finalErr{err}
}

// What the engine gives the adapters of docflow/io besides its API.

// Options are the options the manager runs with, their defaults filled.
func (m *Manager) Options() Options { return m.o }

// Now is the manager's clock.
func (m *Manager) Now() time.Time { return m.now() }

// Sender finds who sends a job with a request (an ingress's, an entry's
// port's): the admin token, a submit token, an app, or, with no
// SubmitTokens, anyone. ok is false for credentials not known, or none
// where some are required.
func (m *Manager) Sender(r *http.Request) (by *Principal, ok bool) {
	return (&handler{m: m}).submitter(r)
}

// MaySend refuses an app a destination its endpoints do not let it send
// to; anyone else may.
func (m *Manager) MaySend(ctx context.Context, r *http.Request, dest string) error {
	if v, ok := (&handler{m: m}).viewerOf(r); ok && v.Kind == PrincipalApp {
		return m.appMaySend(ctx, v, dest, false)
	}
	return nil
}

// CleanNames trims a list of node names, without empty ones and
// duplicates: the edges of a node's settings.
func CleanNames(names []string) []string { return cleanNames(names) }

// PagesView is the view of a document whose pages are printed, drawn and
// counted: its first fixed or flow view of pages, else a scroll view.
func PagesView(doc *bdf.Document) (*bdf.View, error) { return pagesView(doc) }

// Ended tells a state a job does not leave.
func (s JobState) Ended() bool { return s.ended() }

// MetaKeyOK tells the shape of a metadata key: letters, digits, _ . / -.
func MetaKeyOK(k string) bool { return metaKeyRE.MatchString(k) }

// MetaText is a metadata value as text: a string as it is, a number as it
// is written, anything else as JSON.
func MetaText(v any) string { return metaText(v) }

// CheckPreviewURL takes the URL of a printer on a fake printer farm's API
// (tools/fakeprinter), whose pages are shown beside a node.
func CheckPreviewURL(raw string) error { return checkPreviewURL(raw) }

// KindInfo is a kind as GET /api/kinds tells it.
type KindInfo struct {
	Name     string   `json:"name"`
	Title    string   `json:"title"`
	Role     Role     `json:"role"`
	Takes    []string `json:"takes"`
	PassesOn bool     `json:"passes_on,omitempty"`
	Makes    string   `json:"makes,omitempty"`
	Fixed    bool     `json:"fixed,omitempty"`
	// SendsOut tells an exit whose jobs are sent out of the topology
	// (Kind.Send: a printer); Port an entry that listens on a port of
	// its own (an endpoint, a CUPS server).
	SendsOut bool `json:"sends_out,omitempty"`
	Port     bool `json:"port,omitempty"`
	// Step tells a kind whose step is done in the background (render,
	// notify); Keyed an entry with a key of its own at /in/<name>.
	Step  bool `json:"step,omitempty"`
	Keyed bool `json:"keyed,omitempty"`
}

// KindInfos are the kinds registered, as the API tells them.
func KindInfos() []KindInfo {
	ks := Kinds()
	out := make([]KindInfo, 0, len(ks))
	for _, k := range ks {
		out = append(out, KindInfo{Name: k.Name, Title: k.title(), Role: k.Role, Takes: k.takes().names(), PassesOn: k.PassesOn, Makes: k.Makes, Fixed: k.Fixed,
			SendsOut: k.Send != nil, Port: ServesPort(k.Name), Step: k.Step != nil, Keyed: keyed(k.Name)})
	}
	return out
}
