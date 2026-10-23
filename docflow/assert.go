package docflow

import (
	"context"
	"errors"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// An assert checks a job's metadata: every check must hold, or the job
// goes no further. One right after the node a job is sent to (an entry,
// a queue), or that node itself, is the sender's gate: the job is
// checked before it goes anywhere, and the sender's request fails with
// what failed (Refusal: 422 of REST, client-error-bad-request of IPP);
// the job is kept, aborted, and goes to the assert's on_error nodes (a
// notify node). An assert further on checks the job as it reaches it,
// once the sender's request is answered: a job it refuses is aborted
// there, and goes to its on_error nodes. A job that passes goes on to
// To, as it reached the assert. A branch's cases are checks too
// (branch.go).

// Check is a condition on a key of a job's metadata, as the nodes see
// it (metaView: the sender's keys, the entry's, the manager's $ keys,
// name and user, and the payload's $kind, $format, $pages …). The key
// must have a value that is not empty, unless Optional, which passes a
// job without it; a value must then match Pattern (RE2, anywhere in it:
// ^…$ for the whole), be a number from Min to Max, and be one of Enum,
// those given. Message is what a failure says instead of the check's own
// words.
type Check struct {
	Key      string   `json:"key"`
	Optional bool     `json:"optional,omitempty"`
	Pattern  string   `json:"pattern,omitempty"`
	Min      *float64 `json:"min,omitempty"`
	Max      *float64 `json:"max,omitempty"`
	Enum     []string `json:"enum,omitempty"`
	Message  string   `json:"message,omitempty"`
}

// CheckFailure is a check a job failed at a node.
type CheckFailure struct {
	Node  string `json:"node"`
	Key   string `json:"key"`
	Error string `json:"error"`
}

// ErrRefused is the error of a job an assert refused to its sender
// (Refusal).
var ErrRefused = errors.New("manager: refused")

// Refusal is a job the asserts right after the node it was sent to
// refused, and the checks it failed there. The job is kept, aborted.
type Refusal struct {
	Job      *Job
	Failures []CheckFailure
}

func (e *Refusal) Error() string {
	var parts []string
	for _, f := range e.Failures {
		parts = append(parts, f.Node+": "+f.Error)
	}
	return "refused: " + strings.Join(parts, "; ")
}

func (e *Refusal) Unwrap() error { return ErrRefused }

// apiRefusal is the answer of REST to a job refused.
type apiRefusal struct {
	Error string `json:"error"`
	// Job is the number of the job kept, aborted: /api/jobs/{id}.
	Job      int32          `json:"job"`
	Failures []CheckFailure `json:"failures"`
}

func (e *Refusal) body() apiRefusal {
	b := apiRefusal{Error: e.Error(), Failures: e.Failures}
	if e.Job != nil {
		b.Job = e.Job.Number
	}
	return b
}

const (
	maxChecks = 64
	maxCases  = 64
)

// checkKeyRE is the shape of the key a check looks at: a sender's or an
// entry's key, or one of the manager's ($job, $pages …).
var checkKeyRE = regexp.MustCompile(`^\$?[A-Za-z0-9_][A-Za-z0-9_./-]{0,63}$`)

// checkChecks cleans checks and refuses those that cannot be: a key not
// of the shape, an expression that does not compile, bounds the wrong
// way round.
func checkChecks(in []Check) ([]Check, error) {
	if len(in) > maxChecks {
		return nil, fmt.Errorf("%w: %d checks; at most %d", ErrInvalid, len(in), maxChecks)
	}
	out := make([]Check, 0, len(in))
	for i, c := range in {
		c.Key = strings.TrimSpace(c.Key)
		if !checkKeyRE.MatchString(c.Key) {
			return nil, fmt.Errorf("%w: check %d: key %q: letters, digits, _ . - and /, after a $ for the manager's", ErrInvalid, i+1, c.Key)
		}
		if c.Pattern != "" {
			if _, err := regexp.Compile(c.Pattern); err != nil {
				return nil, fmt.Errorf("%w: check %d (%s): pattern %q: %v", ErrInvalid, i+1, c.Key, c.Pattern, err)
			}
		}
		for _, b := range []*float64{c.Min, c.Max} {
			if b != nil && (math.IsNaN(*b) || math.IsInf(*b, 0)) {
				return nil, fmt.Errorf("%w: check %d (%s): a bound is a finite number", ErrInvalid, i+1, c.Key)
			}
		}
		if c.Min != nil && c.Max != nil && *c.Min > *c.Max {
			return nil, fmt.Errorf("%w: check %d (%s): min %v is more than max %v", ErrInvalid, i+1, c.Key, *c.Min, *c.Max)
		}
		var enum []string
		for _, v := range c.Enum {
			if v = strings.TrimSpace(v); v != "" && !slices.Contains(enum, v) {
				enum = append(enum, v)
			}
		}
		c.Enum = enum
		c.Message = clip(strings.TrimSpace(c.Message), 255)
		out = append(out, c)
	}
	return out, nil
}

// failure tells why a value fails the check ("" when it holds): of a key
// the job has (ok), its value as the nodes see it.
func (c *Check) failure(v string, ok bool) string {
	why := ""
	switch {
	case !ok || v == "":
		if !c.Optional {
			why = "required"
		}
	case c.Pattern != "" && !patternOf(c.Pattern).MatchString(v):
		why = fmt.Sprintf("%s does not match %s", quoted(v), c.Pattern)
	case c.Min != nil || c.Max != nil:
		x, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		switch {
		case err != nil || math.IsNaN(x) || math.IsInf(x, 0):
			why = quoted(v) + " is not a number"
		case c.Min != nil && x < *c.Min:
			why = fmt.Sprintf("%s is less than %s", v, num(*c.Min))
		case c.Max != nil && x > *c.Max:
			why = fmt.Sprintf("%s is more than %s", v, num(*c.Max))
		}
	}
	if why == "" && ok && v != "" && len(c.Enum) > 0 && !slices.Contains(c.Enum, v) {
		why = fmt.Sprintf("%s is not one of %s", quoted(v), strings.Join(c.Enum, ", "))
	}
	switch {
	case why == "":
		return ""
	case c.Message != "":
		return c.Message
	}
	return c.Key + ": " + why
}

// patternOf compiles the expression of a check, which compiled when the
// node was put (checkChecks).
func patternOf(expr string) *regexp.Regexp {
	re, err := regexp.Compile(expr)
	if err != nil {
		return regexp.MustCompile(`$.^`) // matches nothing
	}
	return re
}

func num(x float64) string { return strconv.FormatFloat(x, 'f', -1, 64) }

// test checks a job with a payload against checks: the failures, at a
// node; none passes.
func test(checks []Check, node string, view map[string]string) []CheckFailure {
	var out []CheckFailure
	for i := range checks {
		c := &checks[i]
		v, ok := view[c.Key]
		if why := c.failure(v, ok); why != "" {
			out = append(out, CheckFailure{Node: node, Key: c.Key, Error: why})
		}
	}
	return out
}

func failureText(fails []CheckFailure) string {
	parts := make([]string, len(fails))
	for i, f := range fails {
		parts[i] = f.Error
	}
	return strings.Join(parts, "; ")
}

func checkAssert(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	var err error
	n.Checks, err = checkChecks(in.Checks)
	return err
}

// routeAssert sends on a job that passes the assert's checks, as it
// reached it; one that does not is refused there.
func (m *Manager) routeAssert(ctx context.Context, r *routing, name, from string, depth int) {
	n, err := m.nodeOfKind(ctx, name, KindAssert)
	if err != nil {
		r.add(Leg{Node: name, Kind: KindAssert, From: from, State: LegError, Error: err.Error()})
		return
	}
	if fails := test(n.Checks, name, metaView(r.job, r.payload, r.trace)); len(fails) > 0 {
		m.refuse(ctx, r, n, from, fails, depth)
		return
	}
	if len(n.To) == 0 {
		r.add(Leg{Node: name, Kind: KindAssert, From: from, State: LegError, Error: "sends to no node yet: connect it to the nodes after it"})
		return
	}
	r.add(Leg{Node: name, Kind: KindAssert, From: from, State: LegRouted})
	for _, to := range n.To {
		m.route(ctx, r, to, name, depth+1)
	}
}

// refuse gives a job the leg of what an assert refused, and sends it,
// as it reached the assert, to the assert's on_error nodes (none: the
// manager's, unless the job is on an error's way already) with the
// failure, as a send that failed for good is (routeError).
func (m *Manager) refuse(ctx context.Context, r *routing, n *Node, from string, fails []CheckFailure, depth int) {
	why := failureText(fails)
	r.add(Leg{Node: n.Name, Kind: KindAssert, From: from, State: LegRefused, Error: why})
	to := n.OnError
	if len(to) == 0 && r.cause == nil {
		to = m.o.OnError
	}
	to = slices.DeleteFunc(slices.Clone(to), func(name string) bool { return name == n.Name })
	if len(to) == 0 {
		return
	}
	was := r.cause
	r.cause = &Failure{Node: n.Name, Error: why}
	for _, x := range to {
		m.route(ctx, r, x, n.Name, depth+1)
	}
	r.cause = was
}

// gate checks a job against the asserts right after the node it was
// sent to: that node itself when it is one, else those an entry or a
// queue forwards to. A job one of them refuses goes no further than the
// asserts that refused it (and their on_error nodes): the Refusal is the
// sender's answer. nil, the job goes on as it is routed.
func (m *Manager) gate(ctx context.Context, r *routing, dest string) *Refusal {
	n, err := m.Node(ctx, dest)
	if err != nil {
		return nil // route tells
	}
	from, heads := dest, n.To
	switch {
	case n.Kind == KindAssert:
		from, heads = "", []string{dest}
	case !isEntry(n.Kind) && n.Kind != KindQueue:
		return nil
	}
	type refused struct {
		n     *Node
		fails []CheckFailure
	}
	var out []refused
	view := metaView(r.job, r.payload, r.trace)
	for _, name := range heads {
		a, err := m.nodeOfKind(ctx, name, KindAssert)
		if err != nil {
			continue
		}
		if fails := test(a.Checks, name, view); len(fails) > 0 {
			out = append(out, refused{a, fails})
		}
	}
	if len(out) == 0 {
		return nil
	}
	if from != "" {
		r.add(Leg{Node: dest, Kind: n.Kind, State: LegRouted})
	}
	ref := &Refusal{}
	for _, x := range out {
		m.refuse(ctx, r, x.n, from, x.fails, 1)
		ref.Failures = append(ref.Failures, x.fails...)
	}
	return ref
}
