package main

import "strings"

// styleKV is a key=value entry of a style string.
type styleKV struct{ key, value string }

// styleEntries returns the key=value entries of a style string in order,
// the last of repeated keys winning in place of the first (bare named
// styles are dropped).
func styleEntries(s string) []styleKV {
	var out []styleKV
	idx := map[string]int{}
	for _, tok := range strings.Split(s, ";") {
		k, v, ok := strings.Cut(tok, "=")
		if !ok || k == "" {
			continue
		}
		if i, dup := idx[k]; dup {
			out[i].value = v
			continue
		}
		idx[k] = len(out)
		out = append(out, styleKV{k, v})
	}
	return out
}

// parseStyle returns the key=value entries of a style string as a map.
func parseStyle(s string) map[string]string {
	m := map[string]string{}
	for _, kv := range styleEntries(s) {
		m[kv.key] = kv.value
	}
	return m
}
