package scheduler

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/cloudquery/plugin-sdk/v4/internal/checkpoint"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

type testMarker struct {
	name string
	// rows that must be sent before the marker
	parents  []int64
	children []int64
}

func (*testMarker) GetTable() *schema.Table { return nil }
func (*testMarker) IsSyncMessage() bool     { return true }

const (
	checkpointParentTable = "checkpoint_parent"
	checkpointChildTable  = "checkpoint_child"
)

func idColumn() schema.Column {
	return schema.Column{
		Name: "id",
		Type: arrow.PrimitiveTypes.Int64,
		Resolver: func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
			return resource.Set(c.Name, resource.Item.(int64))
		},
	}
}

func childIDs(parent int64) []int64 {
	return []int64{parent * 100, parent*100 + 1}
}

func sendSlice(res chan<- any, ids []int64) {
	res <- ids
}

func sendEach(res chan<- any, ids []int64) {
	for _, id := range ids {
		res <- id
	}
}

// checkpointTables returns a parent table that checkpoints after each page, and a child table that checkpoints after
// its rows.
func checkpointTables(t *testing.T, pages, pageSize int, sendPage func(chan<- any, []int64)) schema.Tables {
	child := &schema.Table{
		Name:    checkpointChildTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, parent *schema.Resource, res chan<- any) error {
			ids := childIDs(parent.Item.(int64))
			res <- ids
			return checkpoint.FromContext(ctx).Checkpoint(ctx, &testMarker{
				name:     fmt.Sprintf("child-%d", parent.Item.(int64)),
				children: ids,
			})
		},
	}
	parent := &schema.Table{
		Name:      checkpointParentTable,
		Columns:   []schema.Column{idColumn()},
		Relations: schema.Tables{child},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			scope := checkpoint.FromContext(ctx)
			require.NotNil(t, scope)
			var parents, children []int64
			for page := range pages {
				ids := make([]int64, pageSize)
				for i := range ids {
					ids[i] = int64(page*pageSize + i + 1)
					parents = append(parents, ids[i])
					children = append(children, childIDs(ids[i])...)
				}
				sendPage(res, ids)
				err := scope.Checkpoint(ctx, &testMarker{
					name:     fmt.Sprintf("page-%d", page),
					parents:  append([]int64(nil), parents...),
					children: append([]int64(nil), children...),
				})
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	child.Parent = parent
	return schema.Tables{parent}
}

func collectSyncMessages(res <-chan message.SyncMessage) (wait func() message.SyncMessages) {
	var msgs message.SyncMessages
	done := make(chan struct{})
	go func() {
		defer close(done)
		for msg := range res {
			msgs = append(msgs, msg)
		}
	}()
	return func() message.SyncMessages {
		<-done
		return msgs
	}
}

func recordIDs(record arrow.RecordBatch) []int64 {
	col := record.Column(record.Schema().FieldIndices("id")[0]).(*array.Int64)
	return col.Int64Values()
}

// requireMarkersAfterRows checks every marker follows the rows it lists, and returns the markers in order.
func requireMarkersAfterRows(t *testing.T, msgs message.SyncMessages) []*testMarker {
	t.Helper()
	sent := map[string]map[int64]bool{
		checkpointParentTable: {},
		checkpointChildTable:  {},
	}
	var markers []*testMarker
	for _, msg := range msgs {
		switch m := msg.(type) {
		case *message.SyncInsert:
			table := m.GetTable().Name
			for _, id := range recordIDs(m.Record) {
				sent[table][id] = true
			}
		case *testMarker:
			for _, id := range m.parents {
				require.Truef(t, sent[checkpointParentTable][id], "%s sent before parent %d", m.name, id)
			}
			for _, id := range m.children {
				require.Truef(t, sent[checkpointChildTable][id], "%s sent before child %d", m.name, id)
			}
			markers = append(markers, m)
		}
	}
	return markers
}

func TestCheckpointsFollowTheirResources(t *testing.T) {
	const pages, pageSize = 4, 3
	for _, strategy := range []Strategy{StrategyDFS, StrategyRoundRobin, StrategyShuffle} {
		for _, batching := range []struct {
			name string
			opt  Option
		}{
			{name: "batching", opt: WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour))},
			{name: "no_batching", opt: WithoutBatching()},
		} {
			for _, send := range []struct {
				name string
				fn   func(chan<- any, []int64)
			}{
				{name: "slices", fn: sendSlice},
				{name: "single_items", fn: sendEach},
			} {
				t.Run(fmt.Sprintf("%s/%s/%s", strategy.String(), batching.name, send.name), func(t *testing.T) {
					for range 20 {
						s := NewScheduler(WithStrategy(strategy), batching.opt)
						res := make(chan message.SyncMessage)
						wait := collectSyncMessages(res)
						err := s.Sync(context.Background(), &testExecutionClient{}, checkpointTables(t, pages, pageSize, send.fn), res)
						close(res)
						msgs := wait()
						require.NoError(t, err)

						markers := requireMarkersAfterRows(t, msgs)
						require.Len(t, markers, pages+pages*pageSize)
						var pageMarkers []string
						for _, m := range markers {
							if m.parents != nil {
								pageMarkers = append(pageMarkers, m.name)
							}
						}
						require.Equal(t, []string{"page-0", "page-1", "page-2", "page-3"}, pageMarkers)
					}
				})
			}
		}
	}
}

func TestCheckpointsWaitForSyncWithShuffleQueue(t *testing.T) {
	const pages, pageSize = 3, 2
	s := NewScheduler(WithStrategy(StrategyShuffleQueue), WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour)))
	res := make(chan message.SyncMessage)
	wait := collectSyncMessages(res)
	err := s.Sync(context.Background(), &testExecutionClient{}, checkpointTables(t, pages, pageSize, sendSlice), res)
	close(res)
	msgs := wait()
	require.NoError(t, err)

	markers := requireMarkersAfterRows(t, msgs)
	require.Len(t, markers, pages+pages*pageSize)
	for _, msg := range msgs[len(msgs)-len(markers):] {
		require.IsType(t, &testMarker{}, msg, "checkpoints must follow every resource")
	}
}

func TestCheckpointAfterSync(t *testing.T) {
	var (
		mu     sync.Mutex
		scopes []checkpoint.Scope
	)
	table := &schema.Table{
		Name:    checkpointParentTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			mu.Lock()
			scopes = append(scopes, checkpoint.FromContext(ctx))
			mu.Unlock()
			res <- []int64{1, 2}
			return nil
		},
	}

	for _, strategy := range []Strategy{StrategyDFS, StrategyShuffleQueue} {
		t.Run(strategy.String(), func(t *testing.T) {
			scopes = nil
			s := NewScheduler(WithStrategy(strategy), WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour)))
			res := make(chan message.SyncMessage)
			wait := collectSyncMessages(res)
			require.NoError(t, s.Sync(context.Background(), &testExecutionClient{}, schema.Tables{table}, res))
			require.Len(t, scopes, 1)
			require.NoError(t, scopes[0].Checkpoint(context.Background(), &testMarker{name: "after", parents: []int64{1, 2}}))
			close(res)

			markers := requireMarkersAfterRows(t, wait())
			require.Len(t, markers, 1)
		})
	}
}

func TestCheckpointAfterCancelledSync(t *testing.T) {
	var scope checkpoint.Scope
	ctx, cancel := context.WithCancel(context.Background())
	table := &schema.Table{
		Name:    checkpointParentTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			scope = checkpoint.FromContext(ctx)
			res <- []int64{1, 2}
			cancel()
			return nil
		},
	}

	// with a table scope, Checkpoint could still be queued while the resolver goroutine winds down
	s := NewScheduler(WithStrategy(StrategyShuffleQueue), WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour)))
	res := make(chan message.SyncMessage)
	wait := collectSyncMessages(res)
	_ = s.Sync(ctx, &testExecutionClient{}, schema.Tables{table}, res)
	require.ErrorIs(t, scope.Checkpoint(context.Background(), &testMarker{name: "after"}), errSyncAborted)
	close(res)
	require.Empty(t, requireMarkersAfterRows(t, wait()))
}
