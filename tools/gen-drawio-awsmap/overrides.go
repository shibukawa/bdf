package main

import (
	"fmt"
	"strings"
)

// override is a curated line of overrides.txt: an old name and the
// target it maps to, or "-" for a name deliberately left unmapped.
type override struct {
	name, spec, comment string
	line                int
}

// parseOverrides reads overrides.txt: one old name and a target per
// line, "#" starting a comment.
func parseOverrides(text string) ([]override, error) {
	var out []override
	seen := map[string]bool{}
	for i, line := range strings.Split(text, "\n") {
		line, comment, _ := strings.Cut(line, "#")
		f := strings.Fields(line)
		if len(f) == 0 {
			continue
		}
		if len(f) != 2 {
			return nil, fmt.Errorf("overrides.txt:%d: want <old name> <target>", i+1)
		}
		if seen[f[0]] {
			return nil, fmt.Errorf("overrides.txt:%d: %s is listed twice", i+1, f[0])
		}
		seen[f[0]] = true
		out = append(out, override{name: f[0], spec: f[1], comment: strings.TrimSpace(comment), line: i + 1})
	}
	return out, nil
}
