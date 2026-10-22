package manager

// A template or a convert node sends on the document it makes along its
// To, and, along its Original, the payload it took as it was (the file
// received, the data): the two may go to archives of their own. Both go
// with a trace of the step: its id (the step's delivery id, the same on
// both sides), the node, which side this is, and both payloads with the
// keys an archive names them by (archiveKey). An archive writes the trace
// into its record (ArchiveRecord.Trace), so that from either object the
// other is found by its name in the archive it went to; the manager does
// not look for it, nor know whether it was written. Rules and exports see
// the trace as $trace, $trace_role, $trace_node and $trace_peer (metaView).
//
// A trace goes with the payload until the next template or convert node
// makes one of its own, through queues, filters, classes, inboxes and an
// archive (on to an export after it).

// Sides of a trace.
const (
	// TraceSource is the payload a step took, sent on along Original.
	TraceSource = "source"
	// TraceOutput is the payload it made, sent on along To.
	TraceOutput = "output"
)

// Trace ties the payloads a step took and made.
type Trace struct {
	// ID is the step's: the id of its delivery, unique to it.
	ID string `json:"id"`
	// Node is the template or convert node that did the step, Kind its
	// kind.
	Node string `json:"node"`
	Kind string `json:"kind"`
	// Role is the side of the payload this trace goes with: TraceSource
	// or TraceOutput.
	Role   string    `json:"role"`
	Source TraceSide `json:"source"`
	Output TraceSide `json:"output"`
}

// TraceSide is a payload of a step, and how an archive names it: Key is
// the key of its objects (the record is Key + ".json"), Object its file.
type TraceSide struct {
	Payload string `json:"payload"`
	Kind    string `json:"kind"`
	Format  string `json:"format,omitempty"`
	Name    string `json:"name,omitempty"`
	SHA256  string `json:"sha256"`
	Key     string `json:"key"`
	Object  string `json:"object"`
}

// newTrace is the trace of a step done by delivery d at node n, which
// took in and made out, as the side role sees it.
func newTrace(d *Delivery, n *Node, j *Job, in, out *Payload, role string) *Trace {
	side := func(p *Payload) TraceSide {
		key := archiveKey(j, p)
		return TraceSide{Payload: p.ID, Kind: p.Kind, Format: p.Format, Name: p.Name, SHA256: p.SHA256, Key: key, Object: key + extensionOf(p)}
	}
	return &Trace{ID: d.ID, Node: n.Name, Kind: n.Kind, Role: role, Source: side(in), Output: side(out)}
}

// as is the trace as the other side sees it.
func (t *Trace) as(role string) *Trace {
	c := *t
	c.Role = role
	return &c
}

// peer is the other side's payload.
func (t *Trace) peer() TraceSide {
	if t.Role == TraceSource {
		return t.Output
	}
	return t.Source
}
