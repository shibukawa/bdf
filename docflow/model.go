package docflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/shibukawa/bdf/form/tag"
)

// Duration is a time.Duration written as text: "90s", "10m", "2h", and
// "14d" for days.
type Duration time.Duration

// ParseDuration reads a Duration: time.ParseDuration's form, or a whole
// number of days ("14d").
func ParseDuration(s string) (Duration, error) {
	if days, ok := strings.CutSuffix(s, "d"); ok {
		n, err := strconv.Atoi(days)
		if err != nil || n < 0 {
			return 0, fmt.Errorf("manager: duration %q", s)
		}
		return Duration(time.Duration(n) * 24 * time.Hour), nil
	}
	d, err := time.ParseDuration(s)
	if err != nil || d < 0 {
		return 0, fmt.Errorf("manager: duration %q", s)
	}
	return Duration(d), nil
}

func (d Duration) String() string {
	t := time.Duration(d)
	if t > 0 && t%(24*time.Hour) == 0 {
		return strconv.Itoa(int(t/(24*time.Hour))) + "d"
	}
	s := t.String()
	if t%time.Minute == 0 {
		s = strings.TrimSuffix(s, "0s")
		if t%time.Hour == 0 {
			s = strings.TrimSuffix(s, "0m")
		}
	}
	return s
}

func (d Duration) d() time.Duration { return time.Duration(d) }

func (d Duration) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	v, err := ParseDuration(s)
	*d = v
	return err
}

// Fallback is the printer that prints a display's slips on paper, and
// after how long an unopened slip is printed (0: only while no device is
// connected).
type Fallback struct {
	Printer string   `json:"printer"`
	After   Duration `json:"after,omitempty"`
}

const (
	defaultKeep         = Duration(14 * 24 * time.Hour)
	defaultOfflineAfter = Duration(2 * time.Minute)
)

func or[T comparable](v, def T) T {
	var zero T
	if v == zero {
		return def
	}
	return v
}

// The core's kinds of node (kinds_core.go): what a job sent to a node
// becomes. Endpoints are entries; queues, classes, filters and inboxes
// route; template and convert nodes make a payload; displays, archives,
// sinks and exports are where a job ends. The kinds of the adapters of
// docflow/io (printers, CUPS servers) are theirs.
const (
	// KindQueue forwards a job to every node of To: a virtual printer
	// in front of others, and a fan-out when To has several.
	KindQueue = "queue"
	// KindClass sends a job to one of its Members, as a class of CUPS
	// does: in turn (round-robin), or the first one (failover). A member
	// that does not take it is skipped for the next.
	KindClass = "class"
	// KindInbox keeps a job until a person releases it to one of the
	// nodes of To (or discards it): a filter worked by hand.
	KindInbox = "inbox"
	// KindArchive writes the job's .bdf and a record of it to a store
	// that is written once and never changed by the manager: an audit
	// trail (the vault of pm-2 adds the retention of S3 Object Lock).
	KindArchive = "archive"
	// KindEndpoint takes documents over HTTP (POST /in/<name>, with the
	// endpoint's own key): a document of any format the converter reads,
	// or a template's id and the data to fill it with. It forwards the
	// job to every node of To, as a queue does.
	KindEndpoint = "endpoint"
	// KindDisplay shows the job as a slip on the devices enrolled at the
	// display (the tablets of a warehouse), and prints on paper what
	// nobody opened in time (Fallback).
	KindDisplay = "display"
	// KindSink keeps the latest documents that reach it, to look at on
	// the topology (sink.go).
	KindSink = "sink"
	// KindFilter sends a job on by rules on what it is (filter.go).
	KindFilter = "filter"
	// KindAssert checks the job's metadata and sends on the job that
	// passes; one right after the node a job is sent to refuses it to the
	// sender (assert.go).
	KindAssert = "assert"
	// KindBranch sends the job along the first of its cases whose checks
	// it passes, else along its default, as a switch does (branch.go).
	KindBranch = "branch"
	// KindTemplate fills one of its templates with the data that reaches
	// it (kind data) and sends the form on as .bdf (transform.go).
	KindTemplate = "template"
	// KindConvert makes a .bdf document of whatever reaches it (a PDF, an
	// Office document, a CSV file …), as an entry does when it converts,
	// and sends it on (transform.go).
	KindConvert = "convert"
	// KindExport is an exit applications and people fetch documents from:
	// it lists what an archive before it wrote, each entry under the
	// owner its metadata names, and serves them at GET /out/<name>
	// (export.go).
	KindExport = "export"
	// KindSign signs the .bdf document that reaches it with the key of
	// the manager, as a witness of who sent it, and sends it on
	// (signnode.go).
	KindSign = "sign"
)

// What an entry (an endpoint, a CUPS server) accepts: documents (PDF,
// Office, images, .bdf …; the default), or data (JSON, YAML, CSV, Excel,
// Parquet: what a template node fills).
const (
	AcceptFile = "file"
	AcceptData = "data"
)

// What an entry sends on of a document: the .bdf document made of it
// (the default), or the file as it was received (kind any).
const (
	OutputBDF  = "bdf"
	OutputFile = "file"
)

// Policies of a class.
const (
	PolicyRoundRobin = "round-robin"
	PolicyFailover   = "failover"
)

// Node is a node of the topology, of any kind. Its name is a destination
// in the one namespace (/in/<name>, a CUPS server's /printers/<name>);
// the settings are the kind's (PutNode drops the others).
type Node struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	// Title is shown instead of the name; Description is a note of what
	// the node is for, shown in the node as the description of a C4
	// element is (its first lines on the canvas, all of it beside).
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	// To are the nodes a queue, an endpoint or a CUPS server forwards to
	// (all of them), the nodes an archive sends on to once it wrote, or
	// the nodes an inbox's items may be released to; of a template or a
	// convert node, the nodes the document it makes goes to, and Original
	// those the payload it took goes to as it was (trace.go).
	To       []string `json:"to,omitempty"`
	Original []string `json:"original,omitempty"`
	// Members and Policy are a class's.
	Members []string `json:"members,omitempty"`
	Policy  string   `json:"policy,omitempty"`
	// Rules and Otherwise are a filter's: the first rule that fits a job
	// forwards it; none, Otherwise (empty: the job ends there).
	Rules     []Rule   `json:"rules,omitempty"`
	Otherwise []string `json:"otherwise,omitempty"`
	// Checks are an assert's: a job must pass every one to go on to To;
	// one that does not goes along OnError (assert.go).
	Checks []Check `json:"checks,omitempty"`
	// Cases and Default are a branch's: the first case whose checks the
	// job passes sends it to the case's nodes; none, Default (empty: the
	// job ends there).
	Cases   []Case   `json:"cases,omitempty"`
	Default []string `json:"default,omitempty"`
	// Keep is how long an inbox shows its released and discarded items,
	// and a display its slips done or canceled (default 14 days).
	Keep Duration `json:"keep,omitempty"`
	// Store is the archive's store, a URL as OpenStore takes for documents
	// (bbolt, s3, r2, gcs); Collection the collection (the key prefix)
	// its objects are in, "archive-<name>" by default.
	Store      string `json:"store,omitempty"`
	Collection string `json:"collection,omitempty"`
	// Accept is what an entry takes: AcceptFile (documents; the default)
	// or AcceptData (data for a template node). Output is what it sends
	// on of a document: OutputBDF (the .bdf made of it; the default) or
	// OutputFile (as received). Formats are the formats it takes, as the
	// converter names them (pdf, docx, image, csv … and bdf), of a
	// convert node those it converts; none takes any. A document of
	// another format is refused (ErrFormat).
	Accept  string   `json:"accept,omitempty"`
	Output  string   `json:"output,omitempty"`
	Formats []string `json:"formats,omitempty"`
	// Meta is metadata an entry adds to every job it takes in: its keys
	// win over the sender's (meta.go).
	Meta map[string]any `json:"meta,omitempty"`
	// Template is the id of the template a template node fills, or
	// TemplateFrom the metadata key whose value names it (a sender's
	// "doc_type", say). ListName is the name the data gets when it is a
	// list of records (a CSV file, a sheet, a JSON array): the file's name
	// by default. Strict refuses data that lacks a name the template
	// looks up.
	Template     string `json:"template,omitempty"`
	TemplateFrom string `json:"template_from,omitempty"`
	ListName     string `json:"list_name,omitempty"`
	// SignRole is the role a sign node signs in (RoleWitness when "").
	SignRole string `json:"sign_role,omitempty"`
	Strict   bool   `json:"strict,omitempty"`
	// KeyHash is the hash of an endpoint's key (argon2id), never the key.
	// Apps are the apps that may send to the endpoint with their own
	// credentials (apps.go), besides the key.
	KeyHash string   `json:"key_hash,omitempty"`
	Apps    []string `json:"apps,omitempty"`
	// Fallback prints a display's slips nobody opened in time, or all of
	// them while no device of the display is connected, so that shipping
	// does not stop. OfflineAfter is how long a device may be silent and
	// still count as connected (default 2 minutes).
	Fallback     *Fallback `json:"fallback,omitempty"`
	OfflineAfter Duration  `json:"offline_after,omitempty"`
	// URL is a printer's: ipp://, ipps://, socket:// or lpd:// (a printer
	// the manager reaches over the network); Route how its jobs reach it
	// (direct, branch: docflow/io/ipp). Preview is a printer's fake printer on a
	// farm (tools/fakeprinter): the URL of the printer on the farm's API,
	// whose pages are shown beside the node (preview.go).
	URL     string `json:"url,omitempty"`
	Route   string `json:"route,omitempty"`
	Preview string `json:"preview,omitempty"`
	// Port is a port of the process's own that an endpoint (REST) or a
	// CUPS server (IPP) listens on (ports.go); TLSCert and TLSKey are the
	// files of a CUPS server's certificate (none: plain ipp).
	Port    int    `json:"port,omitempty"`
	TLSCert string `json:"tls_cert,omitempty"`
	TLSKey  string `json:"tls_key,omitempty"`
	// Latest is how many of the latest documents a sink keeps (default
	// 10).
	Latest int `json:"latest,omitempty"`
	// By, Subject, Fields and Entries are an export's: By is the metadata
	// key whose value says whose a document is (a customer's id the
	// sender gave, say), Subject what of a person is compared with it
	// (SubjectEmail, the default, SubjectSub, SubjectID or SubjectGroup),
	// Fields the metadata keys copied into the entries listed, and
	// Entries how many of the latest entries are listed for one owner,
	// and for everyone (default 1000). Apps are the apps that may read
	// it, as an endpoint's are those that may send.
	By      string   `json:"by,omitempty"`
	Subject string   `json:"subject,omitempty"`
	Fields  []string `json:"fields,omitempty"`
	Entries int      `json:"entries,omitempty"`
	// Scope is who may send to the node from the pages and whose its
	// jobs are (ScopeShared, ScopeGroup, ScopePersonal; empty is
	// shared), and Owner the group (a value of the IdP's groups claim)
	// or the person (ID or email) it belongs to. The jobs sent to a
	// group's node are the group's: its members see them.
	Scope string `json:"scope,omitempty"`
	Owner string `json:"owner,omitempty"`
	// OnError are the nodes a job goes to, with the payload that reached
	// the node, when the node's send or step fails for good (a notify
	// node, say); none, the manager's (Options.OnError), or the job stops
	// there, aborted. Every node has them.
	OnError []string `json:"on_error,omitempty"`
	// Credential names the secret of an exit or a notify node: the
	// environment's DOCFLOW_CRED_<name>, or a file (credential.go). It
	// is never the secret.
	Credential string `json:"credential,omitempty"`
	// The settings of the exits and steps of docflow/io and
	// docflow/render, which expand {{meta.<key>}} (Out.Expand): Format and
	// Page are what a render node makes (pdf, png or svg; the page of a
	// png or an svg); From, Recipients and MailSubject a mail's (smtp, and a
	// notify node's by mail); Text a mail's text or a notification's;
	// Path the key an s3 exit writes; Method and Headers an http exit's
	// request; Via is how a notify node tells (slack, mail or http).
	Format      string            `json:"format,omitempty"`
	Page        int               `json:"page,omitempty"`
	From        string            `json:"from,omitempty"`
	Recipients  []string          `json:"recipients,omitempty"`
	MailSubject string            `json:"mail_subject,omitempty"`
	Text        string            `json:"text,omitempty"`
	Path        string            `json:"path,omitempty"`
	Method      string            `json:"method,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Via         string            `json:"via,omitempty"`
}

func (n *Node) keep() time.Duration { return time.Duration(or(n.Keep, defaultKeep)) }
func (n *Node) offlineAfter() time.Duration {
	return time.Duration(or(n.OfflineAfter, defaultOfflineAfter))
}

func (n *Node) collection() string { return or(n.Collection, "archive-"+n.Name) }

// edges are the nodes a node sends to: To, Original, Members, a filter's
// rules and Otherwise, a branch's cases and Default, a display's fallback
// printer, and OnError.
func (n *Node) edges() []string {
	return append(n.flowEdges(), n.OnError...)
}

// flowEdges are the edges along which a node sends what it emits (or, along
// Original, what it took): all but OnError, which carry what reached it.
func (n *Node) flowEdges() []string {
	out := append(append(append([]string{}, n.To...), n.Original...), n.Members...)
	out = append(out, n.ruleTargets()...)
	out = append(out, n.caseTargets()...)
	if n.Fallback != nil && n.Fallback.Printer != "" {
		out = append(out, n.Fallback.Printer)
	}
	return out
}

// sends tells the kinds a node sends along each of its edges, given the
// kinds that reach it: what it emits along the edges of its flow, what
// reached it along OnError.
func (n *Node) sends(in kinds) map[string]kinds {
	out := map[string]kinds{}
	for _, to := range n.flowEdges() {
		out[to] |= n.emitsTo(in, to)
	}
	for _, to := range n.OnError {
		out[to] |= in
	}
	return out
}

// Template is a template of a template node, which fills it with the data
// that reaches it: a compiled template (.bdft), or one of Excel, Word,
// PowerPoint, PDF forms, pdfme or text, as form.Fill takes them. Its id
// is unique within its node.
type Template struct {
	ID       string      `json:"id"`
	Node     string      `json:"node"`
	Name     string      `json:"name,omitempty"`
	Format   string      `json:"format"`
	Size     int64       `json:"size"`
	SHA256   string      `json:"sha256"`
	Uploaded time.Time   `json:"uploaded"`
	Fields   []tag.Field `json:"fields"`
	// Sample is data the template shows itself filled with, when it
	// has some: a start for the data a sender sends.
	Sample map[string]any `json:"sample,omitempty"`
}

// Point is where a node is drawn on the topology's canvas.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Layout is where the nodes are drawn, by name.
type Layout struct {
	Positions map[string]Point `json:"positions"`
}

// JobState is where a job is, over all the nodes it reached.
type JobState string

const (
	// JobProcessing: a delivery to a printer or an archive is pending.
	JobProcessing JobState = "processing"
	// JobHeld: an inbox holds it, and nothing else is pending.
	JobHeld JobState = "held"
	// JobCompleted: every node took it.
	JobCompleted JobState = "completed"
	// JobCanceled: the sender or an administrator canceled it.
	JobCanceled JobState = "canceled"
	// JobAborted: a delivery failed for good, or an assert refused it.
	JobAborted JobState = "aborted"
)

func (s JobState) ended() bool {
	return s == JobCompleted || s == JobCanceled || s == JobAborted
}

// Payload is one of the things a job carries: what was received (a .bdf
// document, data, or any file), and what the nodes made of it (the .bdf
// a template or a convert node made). Its bytes are kept under its ID
// (the job's for the one received, <job>.<leg> for one a node made).
type Payload struct {
	ID string `json:"id"`
	// Kind is PayloadBDF, PayloadData or PayloadAny; Format the format of
	// the bytes (bdf, pdf, csv, json …); Name the file's name, when the
	// sender gave one.
	Kind   string `json:"kind"`
	Format string `json:"format,omitempty"`
	Name   string `json:"name,omitempty"`
	// Pages counts the pages of a .bdf document.
	Pages  int    `json:"pages,omitempty"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
	// Made is the node that made it ("" for the one received), and At
	// when.
	Made string    `json:"made,omitempty"`
	At   time.Time `json:"at"`
}

// Job is what a sender sent to a destination (a document, or data), and
// its way through the topology: a leg for each node it reached, and the
// payloads the nodes made on the way. Kind, Format, Pages, Size and
// SHA256 are the received payload's (Format the format received, before
// an entry converted it). The payloads are kept for every slip, item and
// delivery of the job.
type Job struct {
	ID     string `json:"id"`
	Number int32  `json:"number"`
	// Dest is the node the sender sent it to.
	Dest string `json:"dest"`
	Name string `json:"name,omitempty"`
	// SubmittedBy is who sent it, as the manager verified it; ownership
	// and scope come from it alone.
	SubmittedBy *Principal `json:"submitted_by,omitempty"`
	// ClaimedUser is who the sender said sent it (requesting-user-name,
	// ?user= of the API, the user of an endpoint's data): shown, never
	// trusted.
	ClaimedUser string `json:"claimed_user,omitempty"`
	// Group is the group whose members see the job besides its
	// submitter: the destination's (or, of a location, its fallback
	// printer's) when it was received. Empty: the submitter's alone.
	Group    string    `json:"group,omitempty"`
	Source   string    `json:"source"`
	Kind     string    `json:"kind,omitempty"`
	Format   string    `json:"format,omitempty"`
	Pages    int       `json:"pages"`
	Size     int64     `json:"size"`
	SHA256   string    `json:"sha256"`
	Received time.Time `json:"received"`
	// Meta is the job's metadata: what the sender gave (a claim, shown
	// and routed by, never trusted), what the entry added, and what the
	// manager verified (the keys that start with $). See meta.go.
	Meta     map[string]any `json:"meta,omitempty"`
	Payloads []Payload      `json:"payloads"`
	State    JobState       `json:"state"`
	Ended    time.Time      `json:"ended,omitzero"`
	Canceled *Mark          `json:"canceled,omitempty"`
	Legs     []Leg          `json:"legs"`
	// Live counts the records that still need the document (slips,
	// items, pending deliveries): the job and the document are removed
	// when it ended and nothing needs it.
	Live int `json:"live"`
}

// Leg is a node a job reached, and what became of it there.
type Leg struct {
	Node string `json:"node"`
	Kind string `json:"kind"`
	// From is the node that forwarded it ("" at the destination).
	From string    `json:"from,omitempty"`
	At   time.Time `json:"at"`
	// State is routed (queue, class), shown or canceled (display), held,
	// released or discarded (inbox), pending, done, failed or canceled
	// (printer, archive), kept (sink), listed (export), refused (assert),
	// or error when the node could not take it.
	State string `json:"state"`
	// Ref is the slip, the item or the delivery made here.
	Ref string `json:"ref,omitempty"`
	// Member is the member a class chose, the rule of a filter, the case
	// of a branch, or the owner an export listed the job for.
	Member string `json:"member,omitempty"`
	// Payload is the payload the node got, and Made the one it made of it
	// (a template or a convert node).
	Payload string `json:"payload,omitempty"`
	Made    string `json:"made,omitempty"`
	Error   string `json:"error,omitempty"`
}

// payload finds a payload of the job by its id; "" is the one received.
func (j *Job) payload(id string) *Payload {
	for i := range j.Payloads {
		if j.Payloads[i].ID == id || (id == "" && i == 0) {
			return &j.Payloads[i]
		}
	}
	return nil
}

// received is the payload received, or nil.
func (j *Job) received() *Payload {
	if len(j.Payloads) == 0 {
		return nil
	}
	return &j.Payloads[0]
}

// document is the .bdf payload to show for the job: the last one made,
// or the one received when it is a .bdf document; nil when the job has
// none.
func (j *Job) document() *Payload {
	for i := len(j.Payloads) - 1; i >= 0; i-- {
		if j.Payloads[i].Kind == PayloadBDF {
			return &j.Payloads[i]
		}
	}
	return nil
}

// States of a leg.
const (
	LegRouted    = "routed"
	LegShown     = "shown"
	LegHeld      = "held"
	LegReleased  = "released"
	LegDiscarded = "discarded"
	LegPending   = "pending"
	LegDone      = "done"
	LegFailed    = "failed"
	LegCanceled  = "canceled"
	LegKept      = "kept"
	LegListed    = "listed"
	LegDropped   = "dropped"
	LegRefused   = "refused"
	LegError     = "error"
)

// derive sets the job's state from its legs, and when it ended.
func (j *Job) derive(now time.Time) {
	if j.Canceled != nil {
		j.State = JobCanceled
	} else {
		pending, held, failed := false, false, false
		for _, l := range j.Legs {
			switch l.State {
			case LegPending:
				pending = true
			case LegHeld:
				held = true
			case LegFailed, LegRefused:
				failed = true
			}
		}
		switch {
		case pending:
			j.State = JobProcessing
		case held:
			j.State = JobHeld
		case failed:
			j.State = JobAborted
		default:
			j.State = JobCompleted
		}
	}
	if j.State.ended() {
		if j.Ended.IsZero() {
			j.Ended = now
		}
	} else {
		j.Ended = time.Time{}
	}
}

// leg finds the leg of a ref.
func (j *Job) leg(ref string) *Leg {
	for i := range j.Legs {
		if j.Legs[i].Ref == ref {
			return &j.Legs[i]
		}
	}
	return nil
}

// ItemState is where an item of an inbox is.
type ItemState string

const (
	ItemHeld      ItemState = "held"
	ItemReleased  ItemState = "released"
	ItemDiscarded ItemState = "discarded"
)

// Item is a job held in an inbox until a person releases it.
type Item struct {
	ID          string     `json:"id"`
	Job         int32      `json:"job"`
	Inbox       string     `json:"inbox"`
	Name        string     `json:"name,omitempty"`
	SubmittedBy *Principal `json:"submitted_by,omitempty"`
	ClaimedUser string     `json:"claimed_user,omitempty"`
	Dest        string     `json:"dest"`
	// Payload is the job's payload held here, and Trace the step it came
	// out of (trace.go), which goes on with it when it is released.
	Payload  string    `json:"payload,omitempty"`
	Trace    *Trace    `json:"trace,omitempty"`
	Pages    int       `json:"pages"`
	Received time.Time `json:"received"`
	State    ItemState `json:"state"`
	// Released tells where it went and who sent it; Discarded who
	// discarded it (or canceled the job).
	Released  *Release `json:"released,omitempty"`
	Discarded *Mark    `json:"discarded,omitempty"`
}

// Release is an item sent on from an inbox.
type Release struct {
	To     []string  `json:"to"`
	Worker string    `json:"worker,omitempty"`
	At     time.Time `json:"at"`
}

func (it *Item) endedAt() time.Time {
	switch {
	case it.Released != nil:
		return it.Released.At
	case it.Discarded != nil:
		return it.Discarded.At
	}
	return time.Time{}
}

// DeliveryState is where a delivery is.
type DeliveryState string

const (
	DeliveryPending  DeliveryState = "pending"
	DeliveryDone     DeliveryState = "done"
	DeliveryFailed   DeliveryState = "failed"
	DeliveryCanceled DeliveryState = "canceled"
)

// Delivery is a job on its way to a printer or an archive, or a step of a
// template or a convert node: made when the job is routed, done by a
// tick (or right away, in the background of a process) under a lease,
// tried again after a failure.
type Delivery struct {
	ID   string `json:"id"`
	Job  int32  `json:"job"`
	Name string `json:"name,omitempty"`
	// Node is the printer, the archive, the class or the transforming
	// node the job was routed to; Printer the printer it goes to now (a
	// member, of a class). Payload is the job's payload it takes.
	Node    string `json:"node"`
	Kind    string `json:"kind"`
	Payload string `json:"payload,omitempty"`
	Printer string `json:"printer,omitempty"`
	Class   string `json:"class,omitempty"`
	// Tried are the members of the class that did not take it, so that
	// the next try goes to another.
	Tried     []string      `json:"tried,omitempty"`
	State     DeliveryState `json:"state"`
	Created   time.Time     `json:"created"`
	Attempts  int           `json:"attempts,omitempty"`
	LastError string        `json:"last_error,omitempty"`
	RetryAt   time.Time     `json:"retry_at,omitzero"`
	Lease     *Lease        `json:"lease,omitempty"`
	// Result is the printer's job id, or the archive's key.
	Result string    `json:"result,omitempty"`
	Done   time.Time `json:"done,omitzero"`
	// Trace is the step the payload came out of, or went into, on its
	// way here (trace.go): an archive writes it into its record.
	Trace *Trace `json:"trace,omitempty"`
	// Cause is the failure that sent the job here along an on_error edge.
	Cause *Failure `json:"cause,omitempty"`
}

// inboxState is what changes in an inbox: its held items and the ended
// ones still shown.
type inboxState struct {
	Seq   int64    `json:"seq"`
	Held  []string `json:"held,omitempty"`
	Ended []string `json:"ended,omitempty"`
}

// destState is the jobs sent to a destination, latest last, the last
// recentJobs of them.
type destState struct {
	Jobs []string `json:"jobs,omitempty"`
}

const recentJobs = 500

// dayRecord is the jobs that ended on a day, for their removal.
type dayRecord struct {
	Jobs []string `json:"jobs,omitempty"`
}

// pendingRecord marks a delivery as pending: the ticks list them.
type pendingRecord struct {
	Job int32 `json:"job"`
}

// SlipState is where a slip is.
type SlipState string

const (
	SlipNew      SlipState = "new"
	SlipOpened   SlipState = "opened"
	SlipDone     SlipState = "done"
	SlipCanceled SlipState = "canceled"
)

// ended tells a slip that is no longer to be worked on.
func (s SlipState) ended() bool { return s == SlipDone || s == SlipCanceled }

// Slip is a job shown at a display (Location is its name): a document as
// .bdf, and what happened to it.
type Slip struct {
	ID       string `json:"id"`
	Job      int32  `json:"job"`
	Location string `json:"location"`
	Name     string `json:"name,omitempty"`
	// SubmittedBy and ClaimedUser are the job's.
	SubmittedBy *Principal `json:"submitted_by,omitempty"`
	ClaimedUser string     `json:"claimed_user,omitempty"`
	Source      string     `json:"source"` // "ipp" or "api"
	Format      string     `json:"format,omitempty"`
	// Payload is the job's .bdf payload shown here.
	Payload  string    `json:"payload,omitempty"`
	Pages    int       `json:"pages"`
	Size     int64     `json:"size"`
	SHA256   string    `json:"sha256"`
	Received time.Time `json:"received"`
	State    SlipState `json:"state"`
	Opened   *Mark     `json:"opened,omitempty"`
	Done     *Mark     `json:"done,omitempty"`
	Canceled *Mark     `json:"canceled,omitempty"`
	// FallbackAt is when an unopened slip is printed on paper; zero
	// when the location has no fallback.
	FallbackAt time.Time `json:"fallback_at,omitzero"`
	Printed    []Printed `json:"printed,omitempty"`
	// Lease is held while a manager sends the slip to a printer, so that
	// two do not send it at once.
	Lease *Lease `json:"lease,omitempty"`
	// Attempts, LastError and RetryAt are those of the prints that
	// failed.
	Attempts  int       `json:"attempts,omitempty"`
	LastError string    `json:"last_error,omitempty"`
	RetryAt   time.Time `json:"retry_at,omitzero"`
}

// Mark is who did something to a slip, and when: a device and the worker
// who said their name there (a record, not an authentication), or the
// sender ("" device).
type Mark struct {
	Device string    `json:"device,omitempty"`
	Worker string    `json:"worker,omitempty"`
	At     time.Time `json:"at"`
}

// Printed is a print of a slip on paper.
type Printed struct {
	Printer string    `json:"printer"`
	Why     string    `json:"why"` // "fallback" or "device"
	Device  string    `json:"device,omitempty"`
	Job     int32     `json:"job,omitempty"` // the printer's job id
	At      time.Time `json:"at"`
}

// Lease is a manager sending a slip.
type Lease struct {
	Owner   string    `json:"owner"`
	Until   time.Time `json:"until"`
	Printer string    `json:"printer"`
	Why     string    `json:"why"`
	Device  string    `json:"device,omitempty"`
}

// fallbackPrinted tells a slip the fallback printed already.
func (s *Slip) fallbackPrinted() bool {
	for _, p := range s.Printed {
		if p.Why == whyFallback {
			return true
		}
	}
	return false
}

const (
	whyFallback = "fallback"
	whyDevice   = "device"
)

// Device is a tablet enrolled at a display (Location is its name). It proves itself with a key
// it made and cannot export (WebCrypto, non-extractable): the manager
// keeps the public key.
type Device struct {
	ID       string    `json:"id"`
	Location string    `json:"location"`
	Name     string    `json:"name"`
	Key      string    `json:"key"` // ECDSA P-256 public key, SPKI DER in base64
	Enrolled time.Time `json:"enrolled"`
	LastSeen time.Time `json:"last_seen,omitzero"`
}

// enrollment is a one-time code that adds a device to a location.
type enrollment struct {
	Location string    `json:"location"`
	Expires  time.Time `json:"expires"`
}

// locationState is what changes at a location: its open slips, the ended
// ones still shown, and its devices. Every change of a slip of the
// location changes it too (Seq), so that its revision tells a device
// polling that something changed.
type locationState struct {
	Seq     int64    `json:"seq"`
	Open    []string `json:"open,omitempty"`
	Ended   []string `json:"ended,omitempty"`
	Devices []string `json:"devices,omitempty"`
}

var (
	// ErrNotFound is the error of what is not there.
	ErrNotFound = errors.New("manager: not found")
	// ErrInvalid is the error of a request that cannot be.
	ErrInvalid = errors.New("manager: invalid")
	// ErrConflict is the error of a change that another one beat.
	ErrConflict = errors.New("manager: changed meanwhile")
	// ErrUnauthorized is the error of a device or token not known.
	ErrUnauthorized = errors.New("manager: not authorized")
	// ErrForbidden is the error of what a known principal may not do.
	ErrForbidden = errors.New("manager: not allowed")
)

// names are those of nodes: they are queue names in URLs, and record
// ids.
var nameRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,63}$`)

func validName(s string) bool { return nameRE.MatchString(s) }

// jobID is the record id of a job: the number with zeros before it, so
// that the ids sort as the numbers do.
func jobID(job int32) string { return fmt.Sprintf("%010d", job) }

// slipID is the record id of the slip of a job at a location: a job
// routed to two locations has a slip at each.
func slipID(job int32, loc string) string { return jobID(job) + "." + loc }

// itemID is the record id of the item of a job in an inbox.
func itemID(job int32, inbox string) string { return jobID(job) + "." + inbox }
