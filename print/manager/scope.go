package manager

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"
)

// Who sent a job and who may see it. A job records the principal the
// manager verified (SubmittedBy) apart from the name the sender gave
// (ClaimedUser): ownership and scope come from the first alone. Whoever a
// request is from is a Viewer: a person signed in to the pages (the web
// application signs them in with OIDC at an external IdP and keeps the
// groups the IdP gave them), or a token of the configuration (the admin
// token, a submit token). The manager decides here what they may see and
// send to.

// Kinds of principals.
const (
	// PrincipalUser is a person signed in through the IdP (OIDC).
	PrincipalUser = "user"
	// PrincipalToken is a sender with a token of the configuration: a
	// submit token's name, or "admin" for the admin token.
	PrincipalToken = "token"
	// PrincipalEndpoint is an endpoint's key (POST /in/<name>).
	PrincipalEndpoint = "endpoint"
	// PrincipalAnonymous is a sender nobody verified (no submit tokens
	// are configured).
	PrincipalAnonymous = "anonymous"
)

// How a principal was verified.
const (
	ViaOIDC        = "oidc"
	ViaAdminToken  = "admin-token"
	ViaSubmitToken = "submit-token"
	ViaEndpointKey = "endpoint-key"
	ViaNone        = "none"
)

// Principal is who sent a job, as the manager verified it.
type Principal struct {
	Kind string `json:"kind"`
	// ID is the person's (Viewer.ID), the token's or the endpoint's
	// name; empty for an anonymous sender.
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

// Scopes of nodes: who may send to them from the pages, and whose their
// jobs are.
const (
	// ScopeShared: anyone signed in sends to it; a job is its
	// submitter's alone. The scope of one that names none.
	ScopeShared = "shared"
	// ScopeGroup: the members of the Owner group send to it, and see its
	// jobs.
	ScopeGroup = "group"
	// ScopePersonal: the Owner person alone (their ID or email).
	ScopePersonal = "personal"
)

func checkScope(scope, owner string) error {
	switch scope {
	case "", ScopeShared:
		if owner != "" && scope == ScopeShared {
			return fmt.Errorf("%w: a shared node has no owner", ErrInvalid)
		}
		return nil
	case ScopeGroup, ScopePersonal:
		if strings.TrimSpace(owner) == "" {
			return fmt.Errorf("%w: a %s node needs its owner", ErrInvalid, scope)
		}
		return nil
	}
	return fmt.Errorf("%w: scope %q: personal, group or shared", ErrInvalid, scope)
}

// Viewer is whoever a request is from: a person signed in to the pages,
// or a token.
type Viewer struct {
	// Kind and Via are as a Principal's (PrincipalUser and ViaOIDC for a
	// person).
	Kind string
	Via  string
	// ID is stable for the person (the web application derives it from
	// the IdP's issuer and subject), or the token's name.
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
	return &Principal{Kind: v.Kind, ID: v.ID, Name: firstOr(v.Name, v.Email), Via: v.Via}
}

// Label is how the viewer is shown, and named in the audit log.
func (v *Viewer) Label() string { return v.Principal().Label() }

// Member tells a member of a group.
func (v *Viewer) Member(group string) bool {
	return group != "" && slices.Contains(v.Groups, group)
}

// is tells the person a personal node's owner names: their ID, or their
// email (as the IdP verified it).
func (v *Viewer) is(owner string) bool {
	return owner != "" && v.Kind == PrincipalUser && (owner == v.ID || v.Email != "" && strings.EqualFold(owner, v.Email))
}

// Submitted tells the job's submitter.
func (v *Viewer) Submitted(j *Job) bool {
	by := j.SubmittedBy
	return v != nil && by != nil && by.Kind == v.Kind && by.ID != "" && by.ID == v.ID
}

// CanSee tells whether the viewer may see a job: its submitter, the
// members of its group, and the admins.
func (v *Viewer) CanSee(j *Job) bool {
	return v != nil && (v.Admin || v.Member(j.Group) || v.Submitted(j))
}

// CanCancel tells whether the viewer may cancel a job: its submitter and
// the admins (the members of its group only see it).
func (v *Viewer) CanCancel(j *Job) bool {
	return v != nil && (v.Admin || v.Submitted(j))
}

// CanSendTo tells whether the viewer may send to a node of a scope and
// owner from the pages.
func (v *Viewer) CanSendTo(scope, owner string) bool {
	if v == nil {
		return false
	}
	switch scope {
	case "", ScopeShared:
		return true
	case ScopeGroup:
		return v.Admin || v.Member(owner)
	case ScopePersonal:
		return v.Admin || v.is(owner)
	}
	return false
}

// destScope tells a destination's scope and owner, and the group its
// jobs get: the destination's, or else (a location's) its fallback
// printer's, when that is a group's.
func (m *Manager) destScope(ctx context.Context, name string) (scope, owner, group string, err error) {
	kind, err := m.kindOf(ctx, name)
	if err != nil {
		return "", "", "", err
	}
	switch kind {
	case KindDisplay:
		l, err := m.Location(ctx, name)
		if err != nil {
			return "", "", "", err
		}
		scope, owner = l.Scope, l.Owner
		if scope != ScopeGroup && l.Fallback != nil {
			if p, err := m.Printer(ctx, l.Fallback.Printer); err == nil && p.Scope == ScopeGroup {
				group = p.Owner
			}
		}
	case KindPrinter:
		p, err := m.Printer(ctx, name)
		if err != nil {
			return "", "", "", err
		}
		scope, owner = p.Scope, p.Owner
	default:
		n, err := m.Node(ctx, name)
		if err != nil {
			return "", "", "", err
		}
		scope, owner = n.Scope, n.Owner
	}
	if scope == ScopeGroup {
		group = owner
	}
	return scope, owner, group, nil
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
