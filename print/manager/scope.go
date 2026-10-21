package manager

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"
)

// Who sent a job and who may see it. A slip records the principal the
// manager verified (SubmittedBy) apart from the name the sender gave
// (ClaimedUser): ownership and scope come from the first alone. A person
// signed in to the pages is a Viewer: the web application signs them in
// (OIDC at an external IdP) and keeps the groups the IdP gave them; the
// manager decides here what they may see and print to.

// Kinds of principals.
const (
	// PrincipalUser is a person signed in through the IdP (OIDC).
	PrincipalUser = "user"
	// PrincipalToken is a sender with a token of the configuration: a
	// submit token's name, or "admin" for the admin token.
	PrincipalToken = "token"
	// PrincipalAnonymous is a sender nobody verified (no submit tokens
	// are configured).
	PrincipalAnonymous = "anonymous"
)

// How a principal was verified.
const (
	ViaOIDC        = "oidc"
	ViaAdminToken  = "admin-token"
	ViaSubmitToken = "submit-token"
	ViaNone        = "none"
)

// Principal is who did something, as the manager verified it.
type Principal struct {
	Kind string `json:"kind"`
	// ID is the person's (Viewer.ID) or the token's name; empty for an
	// anonymous sender.
	ID string `json:"id,omitempty"`
	// Name is shown: the person's name or email at the time.
	Name string `json:"name,omitempty"`
	Via  string `json:"via"`
}

// Label is how the principal is shown.
func (p *Principal) Label() string {
	if p == nil {
		return ""
	}
	return firstOr(p.Name, firstOr(p.ID, p.Kind))
}

// Kinds of locations and printers: who may print to them, and the scope
// their jobs have by default.
const (
	// KindShared: anyone signed in prints to it; a job is its submitter's
	// alone. The kind of one that names none.
	KindShared = "shared"
	// KindGroup: the members of the Owner group print to it, and see its
	// jobs.
	KindGroup = "group"
	// KindPersonal: the Owner person alone (their ID or email).
	KindPersonal = "personal"
)

func checkKind(kind, owner string) error {
	switch kind {
	case "", KindShared:
		if owner != "" && kind == KindShared {
			return fmt.Errorf("%w: a shared one has no owner", ErrInvalid)
		}
		return nil
	case KindGroup, KindPersonal:
		if strings.TrimSpace(owner) == "" {
			return fmt.Errorf("%w: a %s one needs its owner", ErrInvalid, kind)
		}
		return nil
	}
	return fmt.Errorf("%w: kind %q: personal, group or shared", ErrInvalid, kind)
}

// Viewer is a person signed in to the pages.
type Viewer struct {
	// ID is stable for the person (the web application derives it from
	// the IdP's issuer and subject).
	ID    string
	Name  string
	Email string
	// Groups are the values of the IdP's groups claim at the last login.
	Groups []string
	// Admin sees everything and changes the settings.
	Admin bool
}

// Principal is the viewer as the submitter of a job.
func (v *Viewer) Principal() *Principal {
	return &Principal{Kind: PrincipalUser, ID: v.ID, Name: firstOr(v.Name, v.Email), Via: ViaOIDC}
}

// Member tells a member of a group.
func (v *Viewer) Member(group string) bool {
	return group != "" && slices.Contains(v.Groups, group)
}

// is tells the person a personal one's owner names: their ID, or their
// email (as the IdP verified it).
func (v *Viewer) is(owner string) bool {
	return owner != "" && (owner == v.ID || v.Email != "" && strings.EqualFold(owner, v.Email))
}

// CanSee tells whether the viewer may see a slip: the submitter, the
// members of its group, and the admins.
func (v *Viewer) CanSee(s *Slip) bool {
	if v == nil {
		return false
	}
	if v.Admin || v.Member(s.Group) {
		return true
	}
	by := s.SubmittedBy
	return by != nil && by.Kind == PrincipalUser && by.ID == v.ID
}

// CanPrintTo tells whether the viewer may print to a location or printer
// of a kind and owner.
func (v *Viewer) CanPrintTo(kind, owner string) bool {
	if v == nil {
		return false
	}
	switch kind {
	case "", KindShared:
		return true
	case KindGroup:
		return v.Admin || v.Member(owner)
	case KindPersonal:
		return v.Admin || v.is(owner)
	}
	return false
}

// slipGroup is the group a slip of a location gets: the location's, or
// else its fallback printer's, when that is a group's.
func (m *Manager) slipGroup(ctx context.Context, l *Location) string {
	if l.Kind == KindGroup {
		return l.Owner
	}
	if l.Fallback != nil {
		if p, err := m.Printer(ctx, l.Fallback.Printer); err == nil && p.Kind == KindGroup {
			return p.Owner
		}
	}
	return ""
}

// Slips returns every slip, the latest first.
func (m *Manager) Slips(ctx context.Context) ([]Slip, error) {
	recs, err := list[slipRecord](ctx, m.st.db)
	if err != nil {
		return nil, err
	}
	out := make([]Slip, len(recs))
	for i, r := range recs {
		out[i] = Slip(r.Data)
	}
	slices.SortFunc(out, func(a, b Slip) int { return strings.Compare(b.ID, a.ID) })
	return out, nil
}

// SlipsOf returns the slips a viewer may see, the latest first.
func (m *Manager) SlipsOf(ctx context.Context, v *Viewer) ([]Slip, error) {
	all, err := m.Slips(ctx)
	if err != nil {
		return nil, err
	}
	return slices.DeleteFunc(all, func(s Slip) bool { return !v.CanSee(&s) }), nil
}

// Places are the locations and printers a viewer may print to.
type Places struct {
	Locations []Location `json:"locations"`
	Printers  []Printer  `json:"printers"`
}

// PlacesOf returns the locations and printers a viewer may print to.
func (m *Manager) PlacesOf(ctx context.Context, v *Viewer) (*Places, error) {
	ls, err := m.Locations(ctx)
	if err != nil {
		return nil, err
	}
	ps, err := m.Printers(ctx)
	if err != nil {
		return nil, err
	}
	return &Places{
		Locations: slices.DeleteFunc(ls, func(l Location) bool { return !v.CanPrintTo(l.Kind, l.Owner) }),
		Printers:  slices.DeleteFunc(ps, func(p Printer) bool { return !v.CanPrintTo(p.Kind, p.Owner) }),
	}, nil
}

// Person is what the web application keeps of a person who signed in:
// who the IdP said they are, and their groups at the last login.
type Person struct {
	Issuer  string   `json:"iss"`
	Subject string   `json:"sub"`
	Name    string   `json:"name,omitempty"`
	Email   string   `json:"email,omitempty"`
	Groups  []string `json:"groups,omitempty"`
	// Login is when they last signed in.
	Login time.Time `json:"login"`
}

// PutPerson keeps a person under their ID.
func (m *Manager) PutPerson(ctx context.Context, id string, p Person) error {
	_, err := change(ctx, m.st.db, id, true, func(r *personRecord) error {
		*r = personRecord(p)
		return nil
	})
	return err
}

// Person returns a person kept under an ID.
func (m *Manager) Person(ctx context.Context, id string) (*Person, error) {
	rec, err := m.st.db.Get[personRecord](ctx, id)
	if err != nil {
		return nil, mapErr(err)
	}
	p := Person(rec.Data)
	return &p, nil
}
