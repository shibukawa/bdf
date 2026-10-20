package server

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"slices"
	"strings"
	"unicode"

	"github.com/shibukawa/bdf/print/spool"
)

// Roles of tokens.
const (
	// roleSubmit submits jobs and acts on one's own.
	roleSubmit = "submit"
	// roleManage acts on others' jobs and reads the accounting.
	roleManage = "manage"
	// roleAdmin pauses and rejects queues, discovers printers and takes
	// backups.
	roleAdmin = "admin"
	// roleProxy fetches the jobs of the queues of printers of proxies
	// (IPP INFRA), and reports them: the server of a branch.
	roleProxy = "proxy"
)

// caller is who makes a request.
type caller struct {
	// token is the token's name; "" when the server has no tokens.
	token string
	// user is the user of the jobs it submits and sees as its own.
	user   string
	roles  []string
	queues []string // nil: all
}

func (c *caller) can(role string) bool { return slices.Contains(c.roles, role) }

// mayQueue tells a queue the caller may submit to.
func (c *caller) mayQueue(q string) bool { return c.queues == nil || slices.Contains(c.queues, q) }

// owns tells a job the caller submitted (as its user).
func (c *caller) owns(j *spool.Job) bool { return j.User == c.user }

// sees tells a job the caller may see: its own, or any with manage.
func (c *caller) sees(j *spool.Job) bool { return c.can(roleManage) || c.owns(j) }

// acts tells a job the caller may change: its own with submit, or any
// with manage.
func (c *caller) acts(j *spool.Job) bool {
	return c.can(roleManage) || c.owns(j) && c.can(roleSubmit)
}

// token is a token as the server compares it.
type token struct {
	TokenConfig
	sum [32]byte
}

// authenticator finds who makes a request.
type authenticator struct {
	tokens []token
}

func newAuthenticator(c *Config) *authenticator {
	a := &authenticator{}
	for _, t := range c.Tokens {
		a.tokens = append(a.tokens, token{TokenConfig: t, sum: sha256.Sum256([]byte(t.secret))})
	}
	return a
}

// userHeader names the user of a job a token submits for (act-as).
const userHeader = "X-BDF-User"

// caller returns who makes a request, or nil when it is not allowed. A
// server without tokens (on a loopback address) lets anyone of the
// machine do anything, under the user the request names.
func (a *authenticator) caller(r *http.Request) *caller {
	named := r.Header.Get(userHeader)
	if named != "" && !validUser(named) {
		return nil
	}
	if len(a.tokens) == 0 {
		if named == "" {
			named = "anonymous"
		}
		return &caller{user: named, roles: []string{roleSubmit, roleManage, roleAdmin}}
	}
	secret, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if secret = strings.TrimSpace(secret); !ok || secret == "" {
		return nil
	}
	// the sums have one length: the comparison tells nothing of the
	// secrets' lengths, and every token is compared
	sum := sha256.Sum256([]byte(secret))
	var found *token
	for i := range a.tokens {
		if subtle.ConstantTimeCompare(sum[:], a.tokens[i].sum[:]) == 1 {
			found = &a.tokens[i]
		}
	}
	if found == nil {
		return nil
	}
	c := &caller{token: found.Name, user: found.Name, roles: found.Roles, queues: found.Queues}
	if len(c.roles) == 0 {
		c.roles = []string{roleSubmit}
	}
	if named != "" && found.ActAs {
		c.user = named
	}
	return c
}

// validUser tells a user name a request may give: printable, not too long.
func validUser(s string) bool {
	if len(s) > 255 {
		return false
	}
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
