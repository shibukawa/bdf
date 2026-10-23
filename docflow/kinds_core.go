package docflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/shibukawa/yuruidb"
)

// The core's kinds of node. Their settings are checked here; what a job
// that reaches one does is in route.go and the file of the kind.
func init() {
	RegisterKind(Kind{Name: KindEndpoint, Role: RoleEntry, PassesOn: true})
	RegisterKind(Kind{Name: KindQueue, Role: RoleRoute, PassesOn: true, Check: checkForward})
	RegisterKind(Kind{Name: KindClass, Role: RoleRoute, PassesOn: true, Check: checkClass, route: (*Manager).routeClass})
	RegisterKind(Kind{Name: KindFilter, Role: RoleRoute, PassesOn: true, Check: checkFilter, route: (*Manager).routeFilter})
	RegisterKind(Kind{Name: KindAssert, Role: RoleRoute, PassesOn: true, Check: checkAssert, route: (*Manager).routeAssert})
	RegisterKind(Kind{Name: KindBranch, Role: RoleRoute, PassesOn: true, Check: checkBranch, route: (*Manager).routeBranch})
	RegisterKind(Kind{Name: KindInbox, Role: RoleRoute, PassesOn: true, Stocks: true, Check: checkInbox,
		route: func(m *Manager, ctx context.Context, r *routing, name, from string, depth int) {
			m.routeInbox(ctx, r, name, from, depth)
		}})
	RegisterKind(Kind{Name: KindTemplate, Role: RoleStep, Takes: []string{PayloadData}, Makes: PayloadBDF, Check: checkTemplate, route: routeStep(KindTemplate)})
	RegisterKind(Kind{Name: KindConvert, Role: RoleStep, Makes: PayloadBDF, Check: checkConvert, route: routeStep(KindConvert)})
	RegisterKind(Kind{Name: KindSign, Role: RoleStep, Takes: []string{PayloadBDF}, Makes: PayloadBDF, Check: checkSign, route: routeStep(KindSign)})
	RegisterKind(Kind{Name: KindArchive, Role: RoleExit, PassesOn: true, Stocks: true, Check: checkArchive, route: routeStep(KindArchive)})
	RegisterKind(Kind{Name: KindDisplay, Role: RoleExit, Takes: []string{PayloadBDF}, Makes: PayloadBDF, Fixed: true, Check: checkDisplay,
		route: func(m *Manager, ctx context.Context, r *routing, name, from string, _ int) {
			m.routeDisplay(ctx, r, name, from)
		}})
	RegisterKind(Kind{Name: KindSink, Role: RoleExit, Fixed: true, Check: checkSink,
		route: func(m *Manager, ctx context.Context, r *routing, name, from string, _ int) {
			m.routeSink(ctx, r, name, from)
		}})
	RegisterKind(Kind{Name: KindExport, Title: "HTTP (pull)", Role: RoleExit, Takes: []string{PayloadStocked}, Fixed: true, Check: func(ctx context.Context, m *Manager, in Node, n *Node) error { return m.checkExport(ctx, in, n) },
		route: func(m *Manager, ctx context.Context, r *routing, name, from string, _ int) {
			m.routeExport(ctx, r, name, from)
		}})
}

// routeStep makes a job that reaches a node of the kind a delivery: the
// step of a template or a convert node, an archive's write.
func routeStep(kind string) func(m *Manager, ctx context.Context, r *routing, name, from string, _ int) {
	return func(m *Manager, ctx context.Context, r *routing, name, from string, _ int) {
		m.routeDelivery(ctx, r, Leg{Node: name, Kind: kind, From: from}, "", "", kind)
	}
}

// checkForward takes where a node forwards to. A node is put before the
// nodes after it, so none is fine: a job that reaches it meanwhile gets
// an error there (route).
func checkForward(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	return nil
}

// checkEntry checks the settings every entry has: where it forwards, what
// it accepts and sends on, the formats it takes, the metadata it adds,
// its own port, and the apps that may send to it with their own
// credentials (an endpoint's REST, a CUPS server's IPP Basic).
func checkEntry(ctx context.Context, m *Manager, k *Kind, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	n.Accept = or(in.Accept, AcceptFile)
	if n.Accept != AcceptFile && n.Accept != AcceptData {
		return fmt.Errorf("%w: accept %q: file or data", ErrInvalid, in.Accept)
	}
	n.Output = or(in.Output, OutputBDF)
	if n.Output != OutputBDF && n.Output != OutputFile {
		return fmt.Errorf("%w: output %q: bdf (the document made of the file) or file (as received)", ErrInvalid, in.Output)
	}
	if n.Accept == AcceptData {
		n.Output = "" // data goes on as it is
	}
	n.Formats = cleanFormats(in.Formats)
	if err := checkFormats(n.Formats); err != nil {
		return err
	}
	var err error
	if n.Meta, err = claimedMeta(in.Meta); err != nil {
		return err
	}
	n.Port = in.Port
	n.Apps = cleanNames(in.Apps)
	for _, id := range n.Apps {
		if _, err := m.App(ctx, id); err != nil {
			return fmt.Errorf("%w: no app %q", ErrInvalid, id)
		}
	}
	return nil
}

func checkTemplate(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	// its templates are uploaded once it is there: none named yet is
	// fine, and a job that reaches it meanwhile fails there (fill)
	n.Template, n.TemplateFrom = strings.TrimSpace(in.Template), strings.TrimSpace(in.TemplateFrom)
	if n.Template != "" && n.TemplateFrom != "" {
		return fmt.Errorf("%w: a template node names the template it fills (template), or the metadata key that names one (template_from), not both", ErrInvalid)
	}
	if n.Template != "" && !templateIDRE.MatchString(n.Template) {
		return fmt.Errorf("%w: template id %q: letters, digits, _ and -", ErrInvalid, n.Template)
	}
	n.ListName, n.Strict = strings.TrimSpace(in.ListName), in.Strict
	n.Original = cleanNames(in.Original)
	return nil
}

func checkConvert(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	n.Formats = cleanFormats(in.Formats)
	n.Original = cleanNames(in.Original)
	return checkFormats(n.Formats)
}

func checkSign(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To = cleanNames(in.To)
	n.SignRole = strings.TrimSpace(in.SignRole)
	if len(n.SignRole) > 64 {
		return fmt.Errorf("%w: sign_role is 64 bytes at most", ErrInvalid)
	}
	return nil
}

func checkSink(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.Latest = in.Latest
	if n.Latest < 0 || n.Latest > 1000 {
		return fmt.Errorf("%w: latest %d: a sink keeps 1 to 1000 documents", ErrInvalid, in.Latest)
	}
	return nil
}

func checkClass(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.Members = cleanNames(in.Members)
	n.Policy = or(in.Policy, PolicyRoundRobin)
	if n.Policy != PolicyRoundRobin && n.Policy != PolicyFailover {
		return fmt.Errorf("%w: policy %q: round-robin or failover", ErrInvalid, in.Policy)
	}
	return nil
}

func checkFilter(ctx context.Context, m *Manager, in Node, n *Node) error {
	rules, err := checkRules(in.Rules)
	if err != nil {
		return err
	}
	n.Rules, n.Otherwise = rules, cleanNames(in.Otherwise)
	if len(n.Rules) == 0 && len(n.Otherwise) == 0 {
		return fmt.Errorf("%w: a filter has rules, or nodes it sends to otherwise", ErrInvalid)
	}
	return nil
}

func checkInbox(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.To, n.Keep = cleanNames(in.To), in.Keep
	return nil
}

func checkArchive(ctx context.Context, m *Manager, in Node, n *Node) error {
	// what it wrote goes on to To, once written (deliver)
	n.To = cleanNames(in.To)
	n.Store, n.Collection = strings.TrimSpace(in.Store), strings.TrimSpace(in.Collection)
	if n.Store == "" {
		return fmt.Errorf("%w: an archive has a store (bbolt, s3, r2 or gcs)", ErrInvalid)
	}
	if n.Collection != "" {
		if err := yuruidb.ValidateCollection(n.Collection); err != nil {
			return fmt.Errorf("%w: collection %q: letters, digits, _ and -", ErrInvalid, n.Collection)
		}
	}
	// the store is opened now, so that a wrong URL is refused here
	if _, err := m.st.archive(ctx, n.Store, n.collection()); err != nil {
		return fmt.Errorf("%w: store: %v", ErrInvalid, err)
	}
	return nil
}

// checkDisplay takes a fallback that is an exit sending .bdf out (a
// printer).
func checkDisplay(ctx context.Context, m *Manager, in Node, n *Node) error {
	n.Keep, n.OfflineAfter = in.Keep, in.OfflineAfter
	if in.Fallback != nil && strings.TrimSpace(in.Fallback.Printer) != "" {
		n.Fallback = &Fallback{Printer: strings.TrimSpace(in.Fallback.Printer), After: in.Fallback.After}
		k, err := m.kindOf(ctx, n.Fallback.Printer)
		if x := kindNamed(k); err != nil || x == nil || x.Send == nil || !x.takes().has(PayloadBDF) {
			return fmt.Errorf("%w: fallback printer %q is not a printer", ErrInvalid, n.Fallback.Printer)
		}
	}
	return nil
}
