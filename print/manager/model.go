package manager

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

// Location is a place where slips are shown instead of printed: a
// warehouse, and the tablets enrolled there. Its slips come through the
// queue of the same name (/ipp/print/<name>).
type Location struct {
	Name  string `json:"name"`
	Title string `json:"title,omitempty"`
	// Fallback prints the slips nobody opened in time, or all of them
	// while no device of the location is connected, so that shipping does
	// not stop.
	Fallback *Fallback `json:"fallback,omitempty"`
	// Keep is how long a slip that was done or canceled stays on the
	// devices (default 14 days); it is then removed.
	Keep Duration `json:"keep,omitempty"`
	// OfflineAfter is how long a device may be silent and still count as
	// connected (default 2 minutes).
	OfflineAfter Duration `json:"offline_after,omitempty"`
}

// Fallback is the printer that prints a location's slips on paper, and
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

func (l *Location) keep() time.Duration { return time.Duration(or(l.Keep, defaultKeep)) }
func (l *Location) offlineAfter() time.Duration {
	return time.Duration(or(l.OfflineAfter, defaultOfflineAfter))
}

func or[T comparable](v, def T) T {
	var zero T
	if v == zero {
		return def
	}
	return v
}

// Kinds of the nodes of a topology: what a job sent to a node becomes.
// Locations (display) and printers are nodes too; the others are routing
// nodes (Node), which an administrator composes.
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
	// KindDisplay shows the job on the tablets of a location.
	KindDisplay = "display"
	// KindPrinter prints the job on paper.
	KindPrinter = "printer"
)

// Policies of a class.
const (
	PolicyRoundRobin = "round-robin"
	PolicyFailover   = "failover"
)

// Node is a routing node of the topology: a queue, a class, an inbox or
// an archive. Its name is a destination (/ipp/print/<name>), in the one
// namespace of locations, printers and nodes.
type Node struct {
	Name  string `json:"name"`
	Kind  string `json:"kind"`
	Title string `json:"title,omitempty"`
	// To are the nodes a queue forwards to (all of them), or the nodes
	// an inbox's items may be released to.
	To []string `json:"to,omitempty"`
	// Members and Policy are a class's.
	Members []string `json:"members,omitempty"`
	Policy  string   `json:"policy,omitempty"`
	// Keep is how long an inbox shows its released and discarded items
	// (default 14 days).
	Keep Duration `json:"keep,omitempty"`
	// Store is the archive's store, a URL as OpenStore takes for documents
	// (bbolt, s3, r2, gcs); Collection the collection (the key prefix)
	// its objects are in, "archive-<name>" by default.
	Store      string `json:"store,omitempty"`
	Collection string `json:"collection,omitempty"`
	// Templates are the templates an endpoint fills (all when none);
	// KeyHash is the hash of its key (argon2id), never the key.
	Templates []string `json:"templates,omitempty"`
	KeyHash   string   `json:"key_hash,omitempty"`
}

func (n *Node) keep() time.Duration { return time.Duration(or(n.Keep, defaultKeep)) }

func (n *Node) collection() string { return or(n.Collection, "archive-"+n.Name) }

// Template is a template kept for the endpoints to fill: a compiled
// template (.bdft), or one of Excel, Word, PowerPoint, PDF forms, pdfme
// or text, as form.Fill takes them.
type Template struct {
	ID       string      `json:"id"`
	Name     string      `json:"name,omitempty"`
	Format   string      `json:"format"`
	Size     int64       `json:"size"`
	SHA256   string      `json:"sha256"`
	Uploaded time.Time   `json:"uploaded"`
	Fields   []tag.Field `json:"fields"`
	// Sample is data the template shows itself filled with, when it
	// has some: a start for the data an endpoint is sent.
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
	// JobAborted: a delivery failed for good.
	JobAborted JobState = "aborted"
)

func (s JobState) ended() bool {
	return s == JobCompleted || s == JobCanceled || s == JobAborted
}

// Job is a document a sender sent to a destination, and its way through
// the topology: a leg for each node it reached. The document is kept once
// (by the job's number) for every slip, item and delivery of the job.
type Job struct {
	ID     string `json:"id"`
	Number int32  `json:"number"`
	// Dest is the node the sender sent it to.
	Dest     string    `json:"dest"`
	Name     string    `json:"name,omitempty"`
	User     string    `json:"user,omitempty"`
	Source   string    `json:"source"`
	Format   string    `json:"format,omitempty"`
	Pages    int       `json:"pages"`
	Size     int64     `json:"size"`
	SHA256   string    `json:"sha256"`
	Received time.Time `json:"received"`
	State    JobState  `json:"state"`
	Ended    time.Time `json:"ended,omitzero"`
	Canceled *Mark     `json:"canceled,omitempty"`
	Legs     []Leg     `json:"legs"`
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
	// (printer, archive), or error when the node could not take it.
	State string `json:"state"`
	// Ref is the slip, the item or the delivery made here.
	Ref string `json:"ref,omitempty"`
	// Member is the member a class chose.
	Member string `json:"member,omitempty"`
	Error  string `json:"error,omitempty"`
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
			case LegFailed:
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
	ID       string    `json:"id"`
	Job      int32     `json:"job"`
	Inbox    string    `json:"inbox"`
	Name     string    `json:"name,omitempty"`
	User     string    `json:"user,omitempty"`
	Dest     string    `json:"dest"`
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

// Delivery is a job on its way to a printer or an archive: made when the
// job is routed, done by a tick (or right away, in the background of a
// process) under a lease, tried again after a failure.
type Delivery struct {
	ID   string `json:"id"`
	Job  int32  `json:"job"`
	Name string `json:"name,omitempty"`
	// Node is the printer, the archive or the class the job was routed
	// to; Printer the printer it goes to now (a member, of a class).
	Node    string `json:"node"`
	Kind    string `json:"kind"`
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

// Routes of a printer: how its jobs reach it.
const (
	// RouteDirect: the manager sends the job itself (the printer is on a
	// network the manager reaches: the same LAN, or a VPC with Direct
	// Connect or a VPN).
	RouteDirect = "direct"
	// RouteBranch: a bdf printd of the branch fetches the job (INFRA).
	// Not done yet.
	RouteBranch = "branch"
)

// Printer is a printer a manager prints on.
type Printer struct {
	Name string `json:"name"`
	// URL is ipp://, ipps://, socket:// or lpd:// (a printer the manager
	// reaches over the network).
	URL   string `json:"url"`
	Route string `json:"route,omitempty"`
	Info  string `json:"info,omitempty"`
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

// Slip is a job shown at a location: a document as .bdf, and what
// happened to it.
type Slip struct {
	ID       string `json:"id"`
	Job      int32  `json:"job"`
	Location string `json:"location"`
	Name     string `json:"name,omitempty"`
	// User is requesting-user-name as the sender told it: shown, never
	// trusted.
	User     string    `json:"user,omitempty"`
	Source   string    `json:"source"` // "ipp" or "api"
	Format   string    `json:"format,omitempty"`
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

// Device is a tablet enrolled at a location. It proves itself with a key
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
)

// names are those of locations and printers: they are queue names in
// URLs, and record ids.
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
