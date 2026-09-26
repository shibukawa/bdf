package markdown

import (
	"bytes"
	stdhtml "html"
	"strings"
)

// frontMatter is the metadata block at the top of a document: YAML
// between --- lines, or TOML between +++ lines. Only what describes the
// document is read (flat keys, lists, and the name of an author), by the
// Dublin Core element it stands for.
type frontMatter struct {
	terms []string            // DC elements in the order they appear
	dc    map[string][]string // values by DC element
}

// terms maps front matter keys to Dublin Core elements.
var terms = map[string]string{
	"title": "title", "author": "creator", "authors": "creator", "creator": "creator",
	"contributor": "contributor", "contributors": "contributor",
	"date": "created", "published": "created", "created": "created", "publishdate": "created",
	"lastmod": "modified", "updated": "modified", "modified": "modified",
	"lang": "language", "language": "language", "locale": "language",
	"description": "description", "summary": "description", "abstract": "description", "excerpt": "description",
	"tags": "subject", "keywords": "subject", "categories": "subject", "subject": "subject",
	"publisher": "publisher", "license": "rights", "rights": "rights", "copyright": "rights",
	"id": "identifier", "url": "identifier", "identifier": "identifier", "permalink": "identifier",
}

// splitFrontMatter separates the front matter from the document.
func splitFrontMatter(src []byte) (*frontMatter, []byte) {
	fm := &frontMatter{dc: map[string][]string{}}
	var delim string
	switch {
	case bytes.HasPrefix(src, []byte("---\n")) || bytes.HasPrefix(src, []byte("---\r\n")):
		delim = "---"
	case bytes.HasPrefix(src, []byte("+++\n")) || bytes.HasPrefix(src, []byte("+++\r\n")):
		delim = "+++"
	default:
		return fm, src
	}
	rest := src[bytes.IndexByte(src, '\n')+1:]
	var lines []string
	for off := 0; off < len(rest); {
		end := bytes.IndexByte(rest[off:], '\n')
		next := len(rest)
		if end >= 0 {
			next = off + end + 1
		} else {
			end = len(rest) - off
		}
		line := strings.TrimRight(string(rest[off:off+end]), "\r")
		if line == delim || delim == "---" && line == "..." {
			if delim == "+++" {
				fm.toml(lines)
			} else {
				fm.yaml(lines)
			}
			return fm, rest[next:]
		}
		lines = append(lines, line)
		off = next
	}
	// no closing line: not front matter (a thematic break)
	return fm, src
}

func (fm *frontMatter) add(key, value string) {
	term, ok := terms[strings.ToLower(key)]
	value = strings.TrimSpace(value)
	if !ok || value == "" {
		return
	}
	if _, seen := fm.dc[term]; !seen {
		fm.terms = append(fm.terms, term)
	}
	fm.dc[term] = append(fm.dc[term], value)
}

// yaml reads the lines of YAML front matter.
func (fm *frontMatter) yaml(lines []string) {
	key := ""
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		indented := line[0] == ' ' || line[0] == '\t'
		if item, ok := strings.CutPrefix(t, "- "); ok && key != "" {
			// an item of the list under key (or the name of one)
			if k, v, ok := strings.Cut(item, ":"); ok && strings.TrimSpace(k) == "name" {
				item = v
			}
			fm.add(key, unquote(item))
			continue
		}
		k, v, ok := strings.Cut(t, ":")
		if !ok {
			continue
		}
		if indented {
			// a mapping under key: its name
			if key != "" && strings.TrimSpace(k) == "name" {
				fm.add(key, unquote(v))
			}
			continue
		}
		key, v = strings.TrimSpace(k), strings.TrimSpace(v)
		switch {
		case v == "":
		case v == "|" || v == ">" || v == "|-" || v == ">-":
			// a block scalar: the indented lines that follow
			var parts []string
			for i+1 < len(lines) && (strings.TrimSpace(lines[i+1]) == "" || lines[i+1][0] == ' ' || lines[i+1][0] == '\t') {
				i++
				if s := strings.TrimSpace(lines[i]); s != "" {
					parts = append(parts, s)
				}
			}
			fm.add(key, strings.Join(parts, " "))
		case strings.HasPrefix(v, "["):
			for _, s := range splitList(v) {
				fm.add(key, s)
			}
		default:
			fm.add(key, unquote(v))
		}
	}
}

// toml reads the lines of TOML front matter (keys before the first table).
func (fm *frontMatter) toml(lines []string) {
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "[") {
			return
		}
		k, v, ok := strings.Cut(t, "=")
		if !ok {
			continue
		}
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		if strings.HasPrefix(v, "[") {
			for _, s := range splitList(v) {
				fm.add(k, s)
			}
			continue
		}
		fm.add(k, unquote(v))
	}
}

// splitList splits an inline list ([a, "b", 'c']).
func splitList(v string) []string {
	v = strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(v), "["), "]")
	var out []string
	for _, s := range strings.Split(v, ",") {
		if s = unquote(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// unquote removes the quotes of a scalar, or a trailing comment of an
// unquoted one.
func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && (s[0] == '"' && s[len(s)-1] == '"' || s[0] == '\'' && s[len(s)-1] == '\'') {
		inner := s[1 : len(s)-1]
		if s[0] == '"' {
			inner = strings.NewReplacer(`\"`, `"`, `\\`, `\`, `\n`, " ", `\t`, " ").Replace(inner)
		}
		return inner
	}
	if i := strings.Index(s, " #"); i >= 0 {
		s = strings.TrimSpace(s[:i])
	}
	return s
}

// writeHead writes the html start tag and the head: the title, and the
// other elements as Dublin Core meta elements (converter/html reads them).
func (fm *frontMatter) writeHead(out *bytes.Buffer) {
	out.WriteString("<html")
	if l := fm.dc["language"]; len(l) > 0 {
		out.WriteString(` lang="` + stdhtml.EscapeString(l[0]) + `"`)
	}
	out.WriteString(">\n<head>\n<meta charset=\"utf-8\">\n")
	if t := fm.dc["title"]; len(t) > 0 {
		out.WriteString("<title>" + stdhtml.EscapeString(t[0]) + "</title>\n")
	}
	for _, term := range fm.terms {
		for _, v := range fm.dc[term] {
			out.WriteString(`<meta name="dc.` + term + `" content="` + stdhtml.EscapeString(v) + "\">\n")
		}
	}
	out.WriteString("</head>\n")
}
