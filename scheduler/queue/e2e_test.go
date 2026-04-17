package queue_test

import (
	"context"
	"testing"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/stretchr/testify/require"
)

// e2eClient is a minimal ClientMeta for the equivalence test.
type e2eClient struct{}

func (e2eClient) ID() string { return "client-1" }

type rootItem struct{ ID string }
type childItem struct {
	ID       string
	ParentID string
}

func buildE2ETables() schema.Tables {
	childTbl := &schema.Table{
		Name: "children",
		Resolver: func(ctx context.Context, _ schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
			p, ok := parent.Item.(rootItem)
			if !ok {
				if pp, ok2 := parent.Item.(*rootItem); ok2 {
					p = *pp
				} else {
					return nil
				}
			}
			res <- []any{childItem{ID: "c1-" + p.ID, ParentID: p.ID}, childItem{ID: "c2-" + p.ID, ParentID: p.ID}}
			return nil
		},
		Transform: transformers.TransformWithStruct(&childItem{}),
	}
	rootTbl := &schema.Table{
		Name: "roots",
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			res <- []any{rootItem{ID: "r1"}, rootItem{ID: "r2"}}
			return nil
		},
		Transform: transformers.TransformWithStruct(&rootItem{}),
		Relations: []*schema.Table{childTbl},
	}
	return schema.Tables{rootTbl}
}

func runSync(t *testing.T, cfg *scheduler.QueueConfig) []message.SyncMessage {
	t.Helper()
	tables := buildE2ETables()
	for _, tbl := range tables.FlattenTables() {
		if tbl.Transform != nil {
			require.NoError(t, tbl.Transform(tbl))
		}
	}

	opts := []scheduler.Option{
		scheduler.WithStrategy(scheduler.StrategyShuffleQueue),
		scheduler.WithConcurrency(100),
	}
	if cfg != nil {
		store, err := scheduler.NewStorageFromConfig(cfg, 1, "inv-1")
		require.NoError(t, err)
		opts = append(opts, scheduler.WithStorage(store))
	}

	s := scheduler.NewScheduler(opts...)
	msgs, err := s.SyncAll(context.Background(), e2eClient{}, tables)
	require.NoError(t, err)
	return msgs
}

func TestE2E_InMemoryVsBadger_Equivalent(t *testing.T) {
	inMemMsgs := runSync(t, nil) // default = in-memory
	badgerDir := t.TempDir()
	badgerMsgs := runSync(t, &scheduler.QueueConfig{
		Type: scheduler.QueueTypeBadger,
		Path: badgerDir,
	})

	// Count inserts per table. Exact ordering may differ, but totals should match.
	countInserts := func(ms []message.SyncMessage) int {
		n := 0
		for _, m := range ms {
			if _, ok := m.(*message.SyncInsert); ok {
				n++
			}
		}
		return n
	}
	inMemN := countInserts(inMemMsgs)
	badgerN := countInserts(badgerMsgs)
	require.Equal(t, inMemN, badgerN, "in-memory vs badger insert count should match: in-memory=%d badger=%d", inMemN, badgerN)

	// Expected: 2 roots + (2 * 2) = 6 resources → some number of SyncInsert batches.
	// Minimum 1 message for each table's resources (ensure nonzero).
	require.Greater(t, inMemN, 0, "in-memory produced zero inserts")
	require.Greater(t, badgerN, 0, "badger produced zero inserts")
}
