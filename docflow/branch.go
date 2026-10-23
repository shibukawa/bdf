package docflow

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// A branch is a switch: the first of its cases whose checks the job
// passes (assert.go: Check) sends it to the case's nodes; a job no case
// fits goes to Default, or ends there (dropped) when there is none. Each
// case is an arm of its own on the canvas, as Default is. What it sends
// on is what reached it. A filter's rules (filter.go) look at the job by
// regular expressions; a branch's cases by the checks an assert makes.

// Case is an arm of a branch: where a job that passes every check of
// When goes. Name labels the arm ("case <n>" when it has none).
type Case struct {
	Name string   `json:"name,omitempty"`
	When []Check  `json:"when"`
	To   []string `json:"to,omitempty"`
}

// label is how a leg and the canvas name the case, the i-th of its
// branch.
func (c *Case) label(i int) string { return or(c.Name, "case "+strconv.Itoa(i+1)) }

func checkBranch(ctx context.Context, m *Manager, in Node, n *Node) error {
	if len(in.Cases) > maxCases {
		return fmt.Errorf("%w: %d cases; at most %d", ErrInvalid, len(in.Cases), maxCases)
	}
	n.Cases = make([]Case, 0, len(in.Cases))
	var names []string
	for i, c := range in.Cases {
		c.Name = clip(strings.TrimSpace(c.Name), 64)
		if c.Name == "default" || c.Name != "" && slices.Contains(names, c.Name) {
			return fmt.Errorf("%w: case %d: the name %q is another arm's", ErrInvalid, i+1, c.Name)
		}
		names = append(names, c.Name)
		var err error
		if c.When, err = checkChecks(c.When); err != nil {
			return fmt.Errorf("case %s: %w", c.label(i), err)
		}
		if len(c.When) == 0 {
			// a job that fits whatever reaches the branch goes to default
			return fmt.Errorf("%w: case %s checks nothing: give it a check, or send what no case fits to default", ErrInvalid, c.label(i))
		}
		c.To = cleanNames(c.To)
		n.Cases = append(n.Cases, c)
	}
	n.Default = cleanNames(in.Default)
	return nil
}

// caseTargets are the nodes a branch's cases and Default name.
func (n *Node) caseTargets() []string {
	var out []string
	for _, c := range n.Cases {
		for _, to := range c.To {
			if !slices.Contains(out, to) {
				out = append(out, to)
			}
		}
	}
	for _, to := range n.Default {
		if !slices.Contains(out, to) {
			out = append(out, to)
		}
	}
	return out
}

// routeBranch sends the job on along the first case it passes, or
// Default; none ends it here.
func (m *Manager) routeBranch(ctx context.Context, r *routing, name, from string, depth int) {
	n, err := m.nodeOfKind(ctx, name, KindBranch)
	if err != nil {
		r.add(Leg{Node: name, Kind: KindBranch, From: from, State: LegError, Error: err.Error()})
		return
	}
	view := metaView(r.job, r.payload, r.trace)
	arm, to := "default", n.Default
	fits := false
	for i := range n.Cases {
		c := &n.Cases[i]
		if len(test(c.When, name, view)) == 0 {
			arm, to, fits = c.label(i), c.To, true
			break
		}
	}
	switch {
	case len(to) > 0:
		r.add(Leg{Node: name, Kind: KindBranch, From: from, State: LegRouted, Member: arm})
		for _, x := range to {
			m.route(ctx, r, x, name, depth+1)
		}
	case fits:
		r.add(Leg{Node: name, Kind: KindBranch, From: from, State: LegError, Member: arm, Error: arm + " sends to no node yet: connect it to the nodes after it"})
	default:
		r.add(Leg{Node: name, Kind: KindBranch, From: from, State: LegDropped, Error: "no case fits: " + describe(r.job)})
	}
}
