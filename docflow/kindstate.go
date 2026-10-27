package docflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/shibukawa/yuruidb"
)

// kindState is what a kind keeps of a node between ticks (Kind.Tick):
// a JSON value of the kind's own under the node's name (kind-states),
// which every instance sees, as a node's settings are.
type kindState struct {
	Kind string          `json:"kind"`
	Data json.RawMessage `json:"data"`
}

// KindState reads what a kind kept of a node (PutKindState) into v: a
// pointer to the kind's own value. A node with nothing kept leaves v as
// it is; what another kind kept under the name (the node was that kind
// before) is left too.
func (m *Manager) KindState(ctx context.Context, n *Node, v any) error {
	rec, err := m.st.db.Get[kindState](ctx, n.Name)
	if err != nil {
		if errors.Is(err, yuruidb.ErrNotFound) {
			return nil
		}
		return mapErr(err)
	}
	if rec.Data.Kind != n.Kind || len(rec.Data.Data) == 0 {
		return nil
	}
	if err := json.Unmarshal(rec.Data.Data, v); err != nil {
		return fmt.Errorf("docflow: the state of %s: %w", n.Name, err)
	}
	return nil
}

// PutKindState keeps v (a JSON value of the kind's own, of 64 KiB at
// most) of a node between ticks, over what was kept before.
func (m *Manager) PutKindState(ctx context.Context, n *Node, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("docflow: the state of %s: %w", n.Name, err)
	}
	if len(b) > maxKindState {
		return fmt.Errorf("docflow: the state of %s is %d bytes; at most %d", n.Name, len(b), maxKindState)
	}
	_, err = change(ctx, m.st.db, n.Name, true, func(s *kindState) error {
		s.Kind, s.Data = n.Kind, b
		return nil
	})
	return err
}

// maxKindState bounds what a kind keeps of a node.
const maxKindState = 64 << 10

// dropKindState removes what a kind kept of a node (DeleteNode).
func (m *Manager) dropKindState(ctx context.Context, name string) {
	rec, err := m.st.db.Get[kindState](ctx, name)
	if err != nil {
		return
	}
	m.st.db.Delete[kindState](ctx, name, rec.Revision())
}
