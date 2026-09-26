package main

import (
	"fmt"
	"sort"
	"strings"
)

// decision is how an old name is drawn.
type decision struct {
	name oldName
	// target is nil for names left unmapped
	target *target
	// how: "auto" (matched by the key), "override", or "unmapped"
	how, key, comment string
}

// classify decides every old name: an override, else the automatic match.
// Names with neither are returned in unclassified.
func classify(s *sources, ovs []override, targets []*target) (ds []decision, unclassified []oldName, err error) {
	byName := map[string]override{}
	for _, o := range ovs {
		byName[o.name] = o
	}
	known := map[string]bool{}
	for _, n := range s.names {
		known[n.name] = true
		if o, ok := byName[n.name]; ok {
			d := decision{name: n, how: "override", comment: o.comment}
			if o.spec == "-" {
				d.how = "unmapped"
			} else if d.target, err = s.resolveSpec(targets, o.spec); err != nil {
				return nil, nil, fmt.Errorf("overrides.txt:%d: %w", o.line, err)
			}
			ds = append(ds, d)
			continue
		}
		if t, k := s.match(n, targets); t != nil {
			ds = append(ds, decision{name: n, target: t, how: "auto", key: k})
			continue
		}
		unclassified = append(unclassified, n)
	}
	for _, o := range ovs {
		if !known[o.name] {
			return nil, nil, fmt.Errorf("overrides.txt:%d: %s is not a shape of an older generation", o.line, o.name)
		}
	}
	return ds, unclassified, nil
}

// resolveSpec returns the target of an override: a palette entry (see
// findTarget), groupicon:<grIcon> for the icon of a group drawn alone
// (colored like the group's icon), or builtin:<style> for a draw.io
// shape that keeps the cell's colors.
func (s *sources) resolveSpec(targets []*target, spec string) (*target, error) {
	if st, ok := strings.CutPrefix(spec, "builtin:"); ok {
		return &target{kind: "builtin", id: st}, nil
	}
	if icon, ok := strings.CutPrefix(spec, "groupicon:"); ok {
		g, err := findTarget(targets, "group:"+icon)
		if err != nil {
			return nil, err
		}
		info, ok := s.aws4Stencils["mxgraph.aws4."+icon]
		if !ok {
			return nil, fmt.Errorf("no stencil mxgraph.aws4.%s", icon)
		}
		return &target{kind: "groupicon", id: icon, entry: g.entry, category: g.category,
			iconW: info.w, iconH: info.h}, nil
	}
	return findTarget(targets, spec)
}

// report prints the decisions and the coverage per generation.
func report(ds []decision, unclassified []oldName, verbose bool) {
	type counts struct{ total, auto, override, unmapped int }
	byGen := map[string]*counts{}
	var gens []string
	var unmapped []string
	for _, d := range ds {
		c := byGen[d.name.gen]
		if c == nil {
			c = &counts{}
			byGen[d.name.gen] = c
			gens = append(gens, d.name.gen)
		}
		c.total++
		switch {
		case d.how == "unmapped":
			c.unmapped++
			unmapped = append(unmapped, d.name.name+"  # "+d.comment)
		case d.how == "auto":
			c.auto++
		default:
			c.override++
		}
		if verbose {
			t := "-"
			if d.target != nil {
				t = d.target.label()
				if d.target.entry != nil && d.target.kind != "groupicon" {
					t += " (" + d.target.entry.title + ")"
				}
			}
			fmt.Printf("%-9s %-70s %s  %s%s\n", d.how, d.name.name, t, d.key, d.comment)
		}
	}
	sort.Strings(gens)
	fmt.Printf("%-6s %6s %6s %9s %9s\n", "gen", "names", "auto", "override", "unmapped")
	for _, g := range gens {
		c := byGen[g]
		fmt.Printf("%-6s %6d %6d %9d %9d\n", g, c.total, c.auto, c.override, c.unmapped)
	}
	if len(unmapped) > 0 {
		fmt.Println("unmapped:")
		for _, u := range unmapped {
			fmt.Println("  " + u)
		}
	}
	if len(unclassified) > 0 {
		fmt.Println("unclassified (add them to overrides.txt):")
		for _, n := range unclassified {
			fmt.Printf("  %s  (%s)\n", n.name, n.display)
		}
	}
}
