package manager

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Metadata goes with a job through the topology: what the sender gave
// (an API request's X-Meta-* headers and ?meta.<key>= parameters, an IPP
// request's job attributes as ipp.<name>), what the entry adds (Node.Meta)
// and what the manager verified. A sender's keys are claims: a filter
// routes by them and a template node picks its template by them, as the
// administrator set up; nothing is authorized by them. The keys the
// manager verifies start with $, which a sender's cannot: $job,
// $received, $entry, $source, $sender, $sender_kind and $file. At a node,
// rules and templates also see name and user (the job's), $kind,
// $format, $pages, $size, $sha256 and $payload of the payload there, and
// $trace, $trace_role, $trace_node and $trace_peer of the step it came
// out of or went into, and $fingerprints of its lineage (metaView,
// trace.go).

const (
	maxMetaKeys  = 64
	maxMetaBytes = 16 << 10
)

// metaKeyRE is the shape of a sender's or an entry's key: letters,
// digits, _ . - and /, not starting with $.
var metaKeyRE = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_./-]{0,63}$`)

// claimedMeta cleans the metadata a sender (or an entry's settings) gave:
// the keys of the shape, and not too much of it (ErrInvalid).
func claimedMeta(in map[string]any) (map[string]any, error) {
	if len(in) == 0 {
		return nil, nil
	}
	if len(in) > maxMetaKeys {
		return nil, fmt.Errorf("%w: %d metadata keys; at most %d", ErrInvalid, len(in), maxMetaKeys)
	}
	out := make(map[string]any, len(in))
	size := 0
	for k, v := range in {
		if !metaKeyRE.MatchString(k) {
			return nil, fmt.Errorf("%w: metadata key %q: letters, digits, _ . - and /, not starting with $", ErrInvalid, k)
		}
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("%w: metadata %s: %v", ErrInvalid, k, err)
		}
		if size += len(k) + len(b); size > maxMetaBytes {
			return nil, fmt.Errorf("%w: the metadata is larger than %d bytes", ErrInvalid, maxMetaBytes)
		}
		out[k] = v
	}
	return out, nil
}

// metaOfRequest is the metadata of an API request: ?meta.<key>=value and
// X-Meta-<Key>: value headers (the key in lower case).
func metaOfRequest(r *http.Request) map[string]any {
	out := map[string]any{}
	for k, vs := range r.URL.Query() {
		if key, ok := strings.CutPrefix(k, "meta."); ok && len(vs) > 0 {
			out[key] = vs[0]
		}
	}
	for k, vs := range r.Header {
		if key, ok := strings.CutPrefix(http.CanonicalHeaderKey(k), "X-Meta-"); ok && len(vs) > 0 {
			out[strings.ToLower(key)] = vs[0]
		}
	}
	return out
}

// jobMeta assembles a job's metadata: the sender's, the entry's over it,
// and what the manager verified over both.
func jobMeta(claimed map[string]any, entry *Node, j *Job) map[string]any {
	out := map[string]any{}
	maps.Copy(out, claimed)
	if entry != nil {
		maps.Copy(out, entry.Meta)
	}
	out["$job"] = j.Number
	out["$received"] = j.Received.UTC().Format(time.RFC3339Nano)
	out["$entry"] = j.Dest
	out["$source"] = j.Source
	if j.SubmittedBy != nil {
		out["$sender"] = j.SubmittedBy.Label()
		out["$sender_kind"] = j.SubmittedBy.Kind
	}
	if j.Name != "" {
		out["$file"] = j.Name
	}
	return out
}

// claimedOf is the part of a job's metadata its sender gave (or the
// entry): every key but the manager's.
func claimedOf(meta map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range meta {
		if !strings.HasPrefix(k, "$") {
			out[k] = v
		}
	}
	return out
}

// metaView is a job's metadata as rules and templates see it at a node
// with a payload: every value as text, with name and user (the job's),
// the payload's $kind, $format, $pages, $size, $sha256 and $payload, and
// of its trace (nil: none) the step's id ($trace), the side
// ($trace_role: source or output), the node ($trace_node) and the key
// the other side is archived by ($trace_peer), and the ids of the steps
// of its lineage, comma-separated ($fingerprints).
func metaView(j *Job, p *Payload, t *Trace) map[string]string {
	out := make(map[string]string, len(j.Meta)+12)
	for k, v := range j.Meta {
		out[k] = metaText(v)
	}
	out["name"], out["user"] = j.Name, j.ClaimedUser
	if p != nil {
		out["$kind"], out["$format"], out["$payload"], out["$sha256"] = p.Kind, p.Format, p.ID, p.SHA256
		out["$pages"], out["$size"] = strconv.Itoa(p.Pages), strconv.FormatInt(p.Size, 10)
	}
	if t != nil {
		out["$trace"], out["$trace_role"], out["$trace_node"], out["$trace_peer"] = t.ID, t.Role, t.Node, t.peer().Key
		out["$fingerprints"] = t.fingerprintIDs()
	}
	return out
}

// metaText is a metadata value as text: a string as it is, a number as it
// is written, anything else as JSON.
func metaText(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case json.Number:
		return x.String()
	case int:
		return strconv.Itoa(x)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case bool:
		return strconv.FormatBool(x)
	case nil:
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}
