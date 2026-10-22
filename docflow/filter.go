package docflow

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// A filter sends a job on by rules: the first rule whose match fits the
// job (its name, who sent it, its format, its pages, its metadata, the
// kind of its payload) forwards it to the rule's nodes; a job no rule fits
// goes to Otherwise, or ends there (dropped) when there is none. Routing
// by what a job is, where an inbox routes by hand. What it sends on is
// what reached it (the kind too).

// Rule is a rule of a filter: where a job that matches goes.
type Rule struct {
	Match Match    `json:"match"`
	To    []string `json:"to"`
}

// Match is what a rule looks at; every condition given must hold. Name,
// User, Format and the values of Meta are regular expressions (RE2), User
// against the name the sender gave and the verified sender's label, Meta
// against the job's metadata by key (metaView: the sender's keys, the
// entry's, the manager's $ keys and the payload's; a key the job lacks
// does not fit). Kind is the payload's (bdf, data or any).
type Match struct {
	Name     string            `json:"name,omitempty"`
	User     string            `json:"user,omitempty"`
	Format   string            `json:"format,omitempty"`
	MinPages int               `json:"min_pages,omitempty"`
	MaxPages int               `json:"max_pages,omitempty"`
	Meta     map[string]string `json:"meta,omitempty"`
	Kind     string            `json:"kind,omitempty"`
}

// empty tells a match with no condition.
func (mt Match) empty() bool {
	return mt.Name == "" && mt.User == "" && mt.Format == "" && mt.MinPages == 0 && mt.MaxPages == 0 && len(mt.Meta) == 0 && mt.Kind == ""
}

// compiled is a match's expressions, compiled.
type compiled struct {
	name, user, format *regexp.Regexp
	meta               map[string]*regexp.Regexp
	m                  Match
}

func (mt Match) compile() (*compiled, error) {
	c := &compiled{m: mt}
	var err error
	for _, f := range []struct {
		what, expr string
		re         **regexp.Regexp
	}{{"name", mt.Name, &c.name}, {"user", mt.User, &c.user}, {"format", mt.Format, &c.format}} {
		if f.expr == "" {
			continue
		}
		if *f.re, err = regexp.Compile(f.expr); err != nil {
			return nil, fmt.Errorf("%w: match %s %q: %v", ErrInvalid, f.what, f.expr, err)
		}
	}
	if mt.MinPages < 0 || mt.MaxPages < 0 || (mt.MaxPages > 0 && mt.MinPages > mt.MaxPages) {
		return nil, fmt.Errorf("%w: match pages %d..%d", ErrInvalid, mt.MinPages, mt.MaxPages)
	}
	if len(mt.Meta) > 0 {
		c.meta = map[string]*regexp.Regexp{}
		for k, expr := range mt.Meta {
			if strings.TrimSpace(k) == "" {
				return nil, fmt.Errorf("%w: match meta: a key is empty", ErrInvalid)
			}
			re, err := regexp.Compile(expr)
			if err != nil {
				return nil, fmt.Errorf("%w: match meta %s %q: %v", ErrInvalid, k, expr, err)
			}
			c.meta[k] = re
		}
	}
	if mt.Kind != "" && !slices.Contains(payloadKinds, mt.Kind) {
		return nil, fmt.Errorf("%w: match kind %q: bdf, data or any", ErrInvalid, mt.Kind)
	}
	if mt.empty() {
		return nil, fmt.Errorf("%w: a rule matches something (name, user, format, pages, meta or kind)", ErrInvalid)
	}
	return c, nil
}

// fits tells a job with a payload the match fits.
func (c *compiled) fits(j *Job, p *Payload, t *Trace) bool {
	if c.m.Kind != "" && (p == nil || p.Kind != c.m.Kind) {
		return false
	}
	if len(c.meta) > 0 {
		view := metaView(j, p, t)
		for k, re := range c.meta {
			v, ok := view[k]
			if !ok || !re.MatchString(v) {
				return false
			}
		}
	}
	if c.name != nil && !c.name.MatchString(j.Name) {
		return false
	}
	if c.user != nil && !c.user.MatchString(j.ClaimedUser) && !c.user.MatchString(j.SubmittedBy.Label()) {
		return false
	}
	if c.format != nil && !c.format.MatchString(j.Format) {
		return false
	}
	if c.m.MinPages > 0 && j.Pages < c.m.MinPages {
		return false
	}
	if c.m.MaxPages > 0 && j.Pages > c.m.MaxPages {
		return false
	}
	return true
}

// checkRules compiles a filter's rules and cleans their nodes.
func checkRules(rules []Rule) ([]Rule, error) {
	out := make([]Rule, 0, len(rules))
	for i, r := range rules {
		if _, err := r.Match.compile(); err != nil {
			return nil, fmt.Errorf("rule %d: %w", i+1, err)
		}
		r.To = cleanNames(r.To)
		if len(r.To) == 0 {
			return nil, fmt.Errorf("%w: rule %d forwards to at least one node (to)", ErrInvalid, i+1)
		}
		out = append(out, r)
	}
	return out, nil
}

// ruleTargets are the nodes a filter's rules and Otherwise name.
func (n *Node) ruleTargets() []string {
	var out []string
	for _, r := range n.Rules {
		for _, to := range r.To {
			if !slices.Contains(out, to) {
				out = append(out, to)
			}
		}
	}
	for _, to := range n.Otherwise {
		if !slices.Contains(out, to) {
			out = append(out, to)
		}
	}
	return out
}

// routeFilter sends the job on by the first rule that fits, or
// Otherwise; none ends it here.
func (m *Manager) routeFilter(ctx context.Context, r *routing, name, from string, depth int) {
	n, err := m.nodeOfKind(ctx, name, KindFilter)
	if err != nil {
		r.add(Leg{Node: name, Kind: KindFilter, From: from, State: LegError, Error: err.Error()})
		return
	}
	for i, rule := range n.Rules {
		c, err := rule.Match.compile()
		if err != nil {
			r.add(Leg{Node: name, Kind: KindFilter, From: from, State: LegError, Error: err.Error()})
			return
		}
		if c.fits(r.job, r.payload, r.trace) {
			r.add(Leg{Node: name, Kind: KindFilter, From: from, State: LegRouted, Member: fmt.Sprintf("rule %d", i+1)})
			for _, to := range rule.To {
				m.route(ctx, r, to, name, depth+1)
			}
			return
		}
	}
	if len(n.Otherwise) == 0 {
		r.add(Leg{Node: name, Kind: KindFilter, From: from, State: LegDropped, Error: "no rule fits: " + describe(r.job)})
		return
	}
	r.add(Leg{Node: name, Kind: KindFilter, From: from, State: LegRouted, Member: "otherwise"})
	for _, to := range n.Otherwise {
		m.route(ctx, r, to, name, depth+1)
	}
}

// describe tells what a job is, for a leg that no rule fit.
func describe(j *Job) string {
	parts := []string{fmt.Sprintf("%d page(s)", j.Pages)}
	if j.Kind != "" && j.Kind != PayloadBDF {
		parts = append(parts, j.Kind)
	}
	if j.Format != "" {
		parts = append(parts, j.Format)
	}
	if j.Name != "" {
		parts = append(parts, quoted(j.Name))
	}
	return strings.Join(parts, ", ")
}

func quoted(s string) string { return fmt.Sprintf("%q", clip(s, 60)) }
