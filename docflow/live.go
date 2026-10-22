package docflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"time"
)

// What moves in the topology, for the pages that draw it: how many items,
// slips and deliveries are at each node, how many jobs reached it lately,
// and the legs of the last minutes, which a canvas runs along the edges.
// A viewer gets the activity of the jobs they may see (CanSee): a person
// their own and their groups', an administrator every one's.

// Live is the topology's activity as a viewer sees it.
type Live struct {
	// Rev changes when anything below does: a client polls with it
	// (If-None-Match), and gets 304 while it is the same.
	Rev   string              `json:"rev"`
	At    time.Time           `json:"at"`
	Nodes map[string]LiveNode `json:"nodes"`
	// Flows are the legs of the last flowWindow, oldest first.
	Flows []Flow `json:"flows"`
}

// LiveNode is what is at a node now, of the jobs the viewer sees.
type LiveNode struct {
	// Held counts the items an inbox holds, Open the slips open on a
	// display, Pending the deliveries waiting at a printer, an archive or
	// a class, and Failed those given up.
	Held    int `json:"held,omitempty"`
	Open    int `json:"open,omitempty"`
	Pending int `json:"pending,omitempty"`
	Failed  int `json:"failed,omitempty"`
	// Recent counts the jobs that reached the node in the last
	// recentWindow.
	Recent int `json:"recent,omitempty"`
	// Kept counts the jobs a sink keeps, and Latest is the newest of
	// them (its id); Listed counts the jobs an export lists.
	Kept   int    `json:"kept,omitempty"`
	Latest string `json:"latest,omitempty"`
	Listed int    `json:"listed,omitempty"`
}

// Flow is a job that went from one node to the next.
type Flow struct {
	Job   int32     `json:"job"`
	Name  string    `json:"name,omitempty"`
	From  string    `json:"from"`
	To    string    `json:"to"`
	State string    `json:"state"`
	At    time.Time `json:"at"`
}

const (
	recentWindow = time.Hour
	flowWindow   = 10 * time.Minute
	maxFlows     = 200
)

// Live computes the activity a viewer may see, from the recent jobs.
func (m *Manager) Live(ctx context.Context, v *Viewer) (*Live, error) {
	jobs, err := m.jobsOfState(ctx, allJobs, 0)
	if err != nil {
		return nil, err
	}
	now := m.now()
	live := &Live{At: now, Nodes: map[string]LiveNode{}, Flows: []Flow{}}
	h := sha256.New()
	// the recent counts age: the revision changes every minute
	fmt.Fprintf(h, "%d\n", now.Unix()/60)
	kept := map[string]map[string]bool{}   // by sink, the jobs it keeps
	listed := map[string]map[string]bool{} // by export, the jobs it lists now
	for i := range jobs {
		j := &jobs[i]
		if !v.CanSee(j) {
			continue
		}
		fmt.Fprintf(h, "%s|%s|%d|%d\n", j.ID, j.State, len(j.Legs), j.Live)
		for _, l := range j.Legs {
			n := live.Nodes[l.Node]
			if now.Sub(l.At) < recentWindow {
				n.Recent++
			}
			if l.From != "" && now.Sub(l.At) < flowWindow {
				live.Flows = append(live.Flows, Flow{Job: j.Number, Name: j.Name, From: l.From, To: l.Node, State: l.State, At: l.At})
			}
			switch l.State {
			case LegHeld:
				n.Held++
			case LegPending:
				n.Pending++
			case LegFailed:
				n.Failed++
			case LegListed:
				// the leg stays listed: the export's list tells whether it
				// lists the job still (an export deleted and put again does not)
				if listed[l.Node] == nil {
					listed[l.Node] = map[string]bool{}
					for _, id := range m.exportJobs(ctx, exportListID(l.Node, "")) {
						listed[l.Node][id] = true
					}
				}
				if listed[l.Node][j.ID] {
					n.Listed++
				}
			case LegKept:
				if kept[l.Node] == nil {
					kept[l.Node] = map[string]bool{}
					for _, id := range m.sinkJobs(ctx, l.Node) {
						kept[l.Node][id] = true
					}
				}
				if kept[l.Node][j.ID] {
					n.Kept++
					if n.Latest == "" {
						n.Latest = j.ID // the jobs come newest first
					}
				}
			case LegShown:
				// the leg stays shown: the slip tells whether it is still
				// open (a job that needs no record any more has none)
				if j.Live > 0 && l.Ref != "" {
					if s, err := m.Slip(ctx, l.Ref); err == nil {
						fmt.Fprintf(h, "%s|%s\n", s.ID, s.State)
						if !s.State.ended() {
							n.Open++
						}
					}
				}
			}
			live.Nodes[l.Node] = n
		}
	}
	for name, n := range live.Nodes {
		if n == (LiveNode{}) {
			delete(live.Nodes, name)
		}
	}
	slices.SortFunc(live.Flows, func(a, b Flow) int { return a.At.Compare(b.At) })
	if len(live.Flows) > maxFlows {
		live.Flows = live.Flows[len(live.Flows)-maxFlows:]
	}
	live.Rev = hex.EncodeToString(h.Sum(nil))[:16]
	return live, nil
}

// publicTopology is the topology as everyone signed in sees it: the
// nodes, their kinds and edges, without the settings (who may send,
// stores, printers' addresses, formats) and without the counts, which
// Live gives within the viewer's scope.
func publicTopology(ns []TopoNode) []TopoNode {
	for i := range ns {
		n := &ns[i]
		n.Node = Node{Name: n.Name, Kind: n.Kind, Title: n.Title, Description: n.Description, To: n.To, Original: n.Original, Members: n.Members, Policy: n.Policy, Fallback: n.Fallback}
		n.In, n.Templates, n.Held, n.Open, n.Devices, n.Pending, n.Kept, n.Latest, n.Listed = "", nil, 0, 0, 0, 0, 0, "", 0
	}
	return ns
}
