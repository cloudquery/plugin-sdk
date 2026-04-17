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

	// Apply transforms IN PLACE on the actual table tree (not copies from
	// FlattenTables). Transforms populate columns; if we transform copies
	// the originals stay column-less and inserts produce zero rows.
	var applyTransforms func([]*schema.Table)
	applyTransforms = func(ts []*schema.Table) {
		for _, tbl := range ts {
			if tbl.Transform != nil {
				require.NoError(t, tbl.Transform(tbl))
			}
			applyTransforms(tbl.Relations)
		}
	}
	applyTransforms(tables)

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

	// Count ROWS across all SyncInsert messages, not just message count.
	// Batching differences between backends may produce different message
	// counts for the same logical data — row count is the true equivalence.
	countRows := func(ms []message.SyncMessage) int64 {
		var n int64
		for _, m := range ms {
			if im, ok := m.(*message.SyncInsert); ok {
				n += im.Record.NumRows()
			}
		}
		return n
	}
	inMemRows := countRows(inMemMsgs)
	badgerRows := countRows(badgerMsgs)

	// Expected: 2 roots + (2 roots * 2 children each) = 6 rows total.
	require.Equal(t, int64(6), inMemRows, "in-memory row count should match expected fixture output")
	require.Equal(t, inMemRows, badgerRows, "badger row count should match in-memory: in-memory=%d badger=%d", inMemRows, badgerRows)
}
