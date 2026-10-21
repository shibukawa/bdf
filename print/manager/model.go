package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	// Kind is who may print to it (KindShared, KindGroup, KindPersonal;
	// empty is shared), and Owner the group (a value of the IdP's groups
	// claim) or the person (ID or email) it belongs to. The slips of a
	// group's location are the group's: its members see them.
	Kind  string `json:"kind,omitempty"`
	Owner string `json:"owner,omitempty"`
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
	// Kind and Owner are as a Location's. A location's slips are the
	// group's of its fallback printer when the location names none.
	Kind  string `json:"kind,omitempty"`
	Owner string `json:"owner,omitempty"`
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
	// SubmittedBy is who sent it, as the manager verified it; ownership
	// and scope come from it alone.
	SubmittedBy *Principal `json:"submitted_by,omitempty"`
	// ClaimedUser is who the sender said sent it (requesting-user-name,
	// or ?user= of the API): shown, never trusted.
	ClaimedUser string `json:"claimed_user,omitempty"`
	// Group is the group whose members see the slip besides its
	// submitter: the location's, or its fallback printer's, when it
	// was received. Empty: the submitter's alone.
	Group    string    `json:"group,omitempty"`
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

// slipID is the record id of the slip of a job: the number with zeros
// before it, so that the ids sort as the numbers do.
func slipID(job int32) string { return fmt.Sprintf("%010d", job) }
