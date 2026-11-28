package client

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
)

// Change is one canned event the example replays over the Watch stream.
type Change struct {
	// Table is the name of the table the changed resource belongs to.
	Table string `json:"table"`
	// Keys maps every primary key column of the table to its value.
	Keys map[string]any `json:"keys"`
	// Op is "upsert" (the default) or "delete".
	Op string `json:"op"`
	// EventID identifies the change. Acking it checkpoints everything up to it.
	EventID string `json:"event_id"`
}

// PluginOp returns the SDK operation the change describes.
func (c Change) PluginOp() plugin.Op {
	if c.Op == OpDelete {
		return plugin.OpDelete
	}
	return plugin.OpUpsert
}

const (
	OpUpsert = "upsert"
	OpDelete = "delete"
)

type Spec struct {
	NumClients int `json:"num_clients"`
	// Changes are replayed, in order, over the Watch stream.
	Changes []Change `json:"changes,omitempty"`
	// RefreshData holds the rows a refresh sync can find, per table. A row must
	// carry every primary key column of its table.
	RefreshData map[string][]map[string]any `json:"refresh_data,omitempty"`
}

func (s *Spec) SetDefaults() {
	if s.NumClients <= 0 {
		s.NumClients = 1
	}
	for i := range s.Changes {
		if s.Changes[i].Op == "" {
			s.Changes[i].Op = OpUpsert
		}
	}
}

func (s *Spec) Validate() error {
	seen := make(map[string]struct{}, len(s.Changes))
	for i, change := range s.Changes {
		switch {
		case change.Table == "":
			return fmt.Errorf("changes[%d]: table is required", i)
		case change.EventID == "":
			return fmt.Errorf("changes[%d]: event_id is required", i)
		case len(change.Keys) == 0:
			return fmt.Errorf("changes[%d]: keys is required", i)
		case change.Op != OpUpsert && change.Op != OpDelete:
			return fmt.Errorf("changes[%d]: op must be %q or %q, got %q", i, OpUpsert, OpDelete, change.Op)
		}
		if _, ok := seen[change.EventID]; ok {
			return fmt.Errorf("changes[%d]: duplicate event_id %q", i, change.EventID)
		}
		seen[change.EventID] = struct{}{}
	}
	return nil
}
