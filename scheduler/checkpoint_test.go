package scheduler

import (
	"context"
	"errors"
	"fmt"
	"net"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/checkpoint"
	stateV3 "github.com/cloudquery/plugin-sdk/v4/internal/clients/state/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	serversv3 "github.com/cloudquery/plugin-sdk/v4/internal/servers/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
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
	require.ErrorIs(t, scope.Checkpoint(context.Background(), &testMarker{name: "after"}), checkpoint.ErrSyncAborted)
	close(res)
	require.Empty(t, requireMarkersAfterRows(t, wait()))
}

// cancellingBatcher cancels the sync while it flushes.
type cancellingBatcher struct {
	cancel  context.CancelFunc
	flushed bool
}

func (*cancellingBatcher) process(*schema.Resource) {}
func (b *cancellingBatcher) flush() bool {
	b.cancel()
	return b.flushed
}
func (*cancellingBatcher) close() {}

func TestCheckpointNotSentWhenSyncIsCancelledDuringFlush(t *testing.T) {
	for _, flushed := range []bool{false, true} {
		t.Run(fmt.Sprintf("flushed=%t", flushed), func(t *testing.T) {
			// sending is a random choice between res and ctx.Done, so repeat to catch a checkpoint going out
			for range 50 {
				ctx, cancel := context.WithCancel(context.Background())
				res := make(chan message.SyncMessage, 1)
				s := &syncClient{checkpoints: newSyncCheckpoints(res), msgChan: res, logger: zerolog.Nop()}
				resources := make(chan *schema.Resource)
				errCh := make(chan error, 1)
				go func() {
					errCh <- s.processResources(ctx, &cancellingBatcher{cancel: cancel, flushed: flushed}, resources)
				}()
				require.NoError(t, s.checkpoints.send(context.Background(), []message.SyncMessage{&testMarker{name: "checkpoint"}}))
				close(resources)
				require.ErrorIs(t, <-errCh, context.Canceled)
				require.Empty(t, res)
			}
		})
	}
}

func cursorCheckpoint(value string, version uint64) *message.SyncCheckpoint {
	sc := (&schema.Table{Name: "cq_state", Columns: schema.ColumnList{
		{Name: stateKeyColumn, Type: arrow.BinaryTypes.String, PrimaryKey: true},
		{Name: stateValueColumn, Type: arrow.BinaryTypes.String},
		{Name: stateVersionColumn, Type: arrow.PrimitiveTypes.Uint64},
	}}).ToArrowSchema()
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	defer bldr.Release()
	bldr.Field(0).(*array.StringBuilder).Append("cursor")
	bldr.Field(1).(*array.StringBuilder).Append(value)
	bldr.Field(2).(*array.Uint64Builder).Append(version)
	return &message.SyncCheckpoint{Record: bldr.NewRecordBatch()}
}

func stateValues(t *testing.T, msgs message.SyncMessages) map[string]string {
	t.Helper()
	values := map[string]string{}
	for _, msg := range msgs {
		if m, ok := msg.(*message.SyncCheckpoint); ok {
			keys, _ := column[*array.String](m.Record, stateKeyColumn)
			vals, _ := column[*array.String](m.Record, stateValueColumn)
			for i := range keys.Len() {
				values[keys.Value(i)] = vals.Value(i)
			}
		}
	}
	return values
}

func TestHeldCheckpointsAreCoalesced(t *testing.T) {
	table := &schema.Table{
		Name:    checkpointParentTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			scope := checkpoint.FromContext(ctx)
			for page := range 5 {
				res <- []int64{int64(page)}
				if err := scope.Checkpoint(ctx, cursorCheckpoint(strconv.Itoa(page), uint64(page+1))); err != nil {
					return err
				}
			}
			// an older version does not replace the newer one
			return scope.Checkpoint(ctx, cursorCheckpoint("stale", 1))
		},
	}
	s := NewScheduler(WithStrategy(StrategyShuffleQueue))
	res := make(chan message.SyncMessage)
	wait := collectSyncMessages(res)
	require.NoError(t, s.Sync(context.Background(), &testExecutionClient{}, schema.Tables{table}, res))
	close(res)
	msgs := wait()

	var checkpoints int
	for _, msg := range msgs {
		if _, ok := msg.(*message.SyncCheckpoint); ok {
			checkpoints++
		}
	}
	require.Equal(t, 1, checkpoints)
	require.Equal(t, map[string]string{"cursor": "4"}, stateValues(t, msgs))
	require.IsType(t, &message.SyncCheckpoint{}, msgs[len(msgs)-1])
}

func TestPendingCheckpointsDroppedWhenSyncIsCancelled(t *testing.T) {
	for _, strategy := range []Strategy{StrategyDFS, StrategyShuffleQueue} {
		t.Run(strategy.String(), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			table := &schema.Table{
				Name:    checkpointParentTable,
				Columns: []schema.Column{idColumn()},
				Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
					res <- []int64{1, 2}
					err := checkpoint.FromContext(ctx).Checkpoint(ctx, &testMarker{name: "pending", parents: []int64{1, 2}})
					cancel()
					return err
				},
			}
			// the batch never fills, so the rows are dropped when the sync is cancelled
			s := NewScheduler(WithStrategy(strategy), WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour)))
			res := make(chan message.SyncMessage)
			wait := collectSyncMessages(res)
			_ = s.Sync(ctx, &testExecutionClient{}, schema.Tables{table}, res)
			close(res)
			require.Empty(t, requireMarkersAfterRows(t, wait()))
		})
	}
}

func TestJoinedCheckpointFollowsEveryScope(t *testing.T) {
	var (
		parentScope = make(chan checkpoint.Scope, 1)
		joinedSent  = make(chan struct{})
	)
	parent := &schema.Table{
		Name:    checkpointParentTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			res <- []int64{1, 2, 3}
			parentScope <- checkpoint.FromContext(ctx)
			<-joinedSent
			return nil
		},
	}
	other := &schema.Table{
		Name:    checkpointChildTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			res <- []int64{10, 11}
			scope := checkpoint.Join(<-parentScope, checkpoint.FromContext(ctx))
			defer close(joinedSent)
			return scope.Checkpoint(ctx, &testMarker{name: "joined", parents: []int64{1, 2, 3}, children: []int64{10, 11}})
		},
	}
	for range 20 {
		joinedSent = make(chan struct{})
		s := NewScheduler(WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour)))
		res := make(chan message.SyncMessage)
		wait := collectSyncMessages(res)
		require.NoError(t, s.Sync(context.Background(), &testExecutionClient{}, schema.Tables{parent, other}, res))
		close(res)
		require.Len(t, requireMarkersAfterRows(t, wait()), 1)
	}
}

func TestHeldCheckpointsAreMarkedSentWhenTheSyncCompletes(t *testing.T) {
	var sent []*checkpoint.Sent
	table := &schema.Table{
		Name:    checkpointParentTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			for page := range 3 {
				res <- []int64{int64(page)}
				s, err := checkpoint.Send(ctx, nil, checkpoint.FromContext(ctx), nil, cursorCheckpoint(strconv.Itoa(page), uint64(page+1)))
				if err != nil {
					return err
				}
				sent = append(sent, s)
			}
			return nil
		},
	}
	s := NewScheduler(WithStrategy(StrategyShuffleQueue))
	res := make(chan message.SyncMessage)
	wait := collectSyncMessages(res)
	require.NoError(t, s.Sync(context.Background(), &testExecutionClient{}, schema.Tables{table}, res))
	close(res)
	msgs := wait()

	require.Equal(t, map[string]string{"cursor": "2"}, stateValues(t, msgs))
	require.Len(t, sent, 3)
	for _, s := range sent {
		require.True(t, s.Done())
	}
}

func TestTableScopeCheckpointAfterCancelledSync(t *testing.T) {
	for range 20 {
		var scope checkpoint.Scope
		ctx, cancel := context.WithCancel(context.Background())
		table := &schema.Table{
			Name:    checkpointParentTable,
			Columns: []schema.Column{idColumn()},
			Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
				scope = checkpoint.FromContext(ctx)
				res <- []int64{1, 2}
				time.Sleep(20 * time.Millisecond)
				cancel()
				return nil
			},
		}
		s := NewScheduler(WithStrategy(StrategyDFS), WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour)))
		res := make(chan message.SyncMessage)
		wait := collectSyncMessages(res)
		_ = s.Sync(ctx, &testExecutionClient{}, schema.Tables{table}, res)
		// the resolver ended, so its scope sends through the sync, which was aborted
		err := scope.Checkpoint(context.Background(), &testMarker{name: "after", parents: []int64{1, 2}})
		close(res)
		require.ErrorIs(t, err, checkpoint.ErrSyncAborted)
		require.Empty(t, requireMarkersAfterRows(t, wait()))
	}
}

func TestShuffleQueueMergesCheckpointsOfKeysAlsoSetOutsideResolvers(t *testing.T) {
	var sent []*checkpoint.Sent
	table := &schema.Table{
		Name:    checkpointParentTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			scope := checkpoint.FromContext(ctx)
			var previous []*checkpoint.Sent
			for page := range 6 {
				res <- []int64{int64(page)}
				// every other checkpoint is sent as if the key was set outside a resolver
				var pageScope checkpoint.Scope
				if page%2 == 0 {
					pageScope = scope
				}
				s, err := checkpoint.Send(ctx, nil, pageScope, previous, cursorCheckpoint(strconv.Itoa(page), uint64(page+1)))
				if err != nil {
					return err
				}
				previous = []*checkpoint.Sent{s}
				sent = append(sent, s)
			}
			return nil
		},
	}
	s := NewScheduler(WithStrategy(StrategyShuffleQueue))
	res := make(chan message.SyncMessage)
	wait := collectSyncMessages(res)
	require.NoError(t, s.Sync(context.Background(), &testExecutionClient{}, schema.Tables{table}, res))
	close(res)
	msgs := wait()

	var checkpoints int
	for _, msg := range msgs {
		if _, ok := msg.(*message.SyncCheckpoint); ok {
			checkpoints++
		}
	}
	require.Equal(t, 1, checkpoints, "every checkpoint is merged")
	require.Equal(t, map[string]string{"cursor": "5"}, stateValues(t, msgs))
	for _, s := range sent {
		require.True(t, s.Done())
	}
}

func TestAbortDropsCheckpointsOfStuckResolvers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var (
		release = make(chan struct{})
		queued  = make(chan struct{})
		gate    = make(chan struct{})
		resume  = make(chan struct{})
		res     = make(chan message.SyncMessage)
	)
	defer close(release)
	stateClient := newCheckpointStateClient(t, res)
	stuck := &schema.Table{
		Name:    checkpointParentTable,
		Columns: []schema.Column{idColumn()},
		// the resolver goroutine never gets past this resource, so it never forwards the checkpoint after it
		PreResourceResolver: func(context.Context, schema.ClientMeta, *schema.Resource) error {
			<-release
			return errors.New("released")
		},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			res <- int64(1)
			defer close(queued)
			if err := stateClient.SetKey(ctx, "cursor", "1"); err != nil {
				return err
			}
			return stateClient.Flush(ctx)
		},
	}
	other := &schema.Table{
		Name:    checkpointChildTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			res <- int64(2)
			return checkpoint.FromContext(ctx).Checkpoint(ctx, &testMarker{name: "other"})
		},
	}

	var cursorSent atomic.Bool
	go func() {
		for msg := range res {
			if _, ok := msg.(*message.SyncCheckpoint); ok {
				cursorSent.Store(true)
			}
			// stop reading once the other table's row arrives, so the sync is cancelled while sending its checkpoint
			if insert, ok := msg.(*message.SyncInsert); ok && insert.GetTable().Name == checkpointChildTable {
				close(gate)
				<-resume
			}
		}
	}()
	defer close(res)
	defer close(resume)

	s := NewScheduler(WithStrategy(StrategyDFS), WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour)))
	syncErr := make(chan error, 1)
	go func() {
		syncErr <- s.Sync(ctx, &testExecutionClient{}, schema.Tables{stuck, other}, res)
	}()
	receiveWithin(t, queued)
	receiveWithin(t, gate)
	cancel()
	require.ErrorIs(t, receiveWithin(t, syncErr), context.Canceled)
	require.ErrorIs(t, stateClient.Flush(context.Background()), checkpoint.ErrSyncAborted, "the checkpoint behind the stuck resolver is aborted")
	require.False(t, cursorSent.Load())
}

func TestShuffleQueueAbortDropsCheckpointsOfStuckResolvers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var (
		release   = make(chan struct{})
		queued    = make(chan struct{})
		cancelled = make(chan struct{})
		done      = make(chan struct{})
		res       = make(chan message.SyncMessage)
	)
	defer close(release)
	defer close(done)
	stateClient := newCheckpointStateClient(t, res)
	stuck := &schema.Table{
		Name:    checkpointParentTable,
		Columns: []schema.Column{idColumn()},
		PreResourceResolver: func(context.Context, schema.ClientMeta, *schema.Resource) error {
			<-release
			return errors.New("released")
		},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			res <- int64(1)
			defer close(queued)
			if err := stateClient.SetKey(ctx, "cursor", "1"); err != nil {
				return err
			}
			return stateClient.Flush(ctx)
		},
	}
	// sends resources after the cancel, so that the sync notices it while the stuck resolver never returns
	waker := &schema.Table{
		Name:    checkpointChildTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			<-cancelled
			ids := make([]int64, 64)
			for i := range ids {
				ids[i] = int64(i)
			}
			res <- ids
			return nil
		},
	}

	var cursorSent atomic.Bool
	go func() {
		for {
			select {
			case msg := <-res:
				if _, ok := msg.(*message.SyncCheckpoint); ok {
					cursorSent.Store(true)
				}
			case <-done:
				return
			}
		}
	}()

	s := NewScheduler(WithStrategy(StrategyShuffleQueue), WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour)))
	syncErr := make(chan error, 1)
	go func() {
		syncErr <- s.Sync(ctx, &testExecutionClient{}, schema.Tables{stuck, waker}, res)
	}()
	receiveWithin(t, queued)
	cancel()
	close(cancelled)
	require.ErrorIs(t, receiveWithin(t, syncErr), context.Canceled)
	require.ErrorIs(t, stateClient.Flush(context.Background()), checkpoint.ErrSyncAborted, "the held checkpoint is aborted")
	require.False(t, cursorSent.Load())
}

func TestTableScopeStartedAfterAbortDropsCheckpoints(t *testing.T) {
	c := newSyncCheckpoints(make(chan message.SyncMessage))
	close(c.loopDone)
	c.abort()

	sent, err := checkpoint.Send(context.Background(), nil, newTableScope(c), nil, &testMarker{name: "late"})
	require.ErrorIs(t, err, checkpoint.ErrSyncAborted)
	require.True(t, sent.Aborted())
}

func TestEndedTableScopesCollapse(t *testing.T) {
	scopes := make(chan checkpoint.Scope, 2)
	collapsedWhileRunning := make(chan bool, 2)
	table := func(name string) *schema.Table {
		return &schema.Table{
			Name:    name,
			Columns: []schema.Column{idColumn()},
			Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
				scope := checkpoint.FromContext(ctx)
				collapsedWhileRunning <- checkpoint.Collapse(scope) != scope
				scopes <- scope
				res <- int64(1)
				return nil
			},
		}
	}
	s := NewScheduler(WithStrategy(StrategyDFS))
	res := make(chan message.SyncMessage)
	wait := collectSyncMessages(res)
	require.NoError(t, s.Sync(context.Background(), &testExecutionClient{}, schema.Tables{table(checkpointParentTable), table(checkpointChildTable)}, res))
	close(res)
	wait()

	require.False(t, <-collapsedWhileRunning)
	require.False(t, <-collapsedWhileRunning)
	a, b := <-scopes, <-scopes
	require.NotSame(t, a, b)
	require.Same(t, checkpoint.Collapse(a), checkpoint.Collapse(b))
	require.Same(t, checkpoint.Collapse(a), checkpoint.Join(a, b), "joining ended scopes waits on the sync once")
}

func TestSyncAllDropsCheckpointsAfterItReturns(t *testing.T) {
	var scope checkpoint.Scope
	table := &schema.Table{
		Name:    checkpointParentTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			scope = checkpoint.FromContext(ctx)
			res <- int64(1)
			return nil
		},
	}
	for _, strategy := range []Strategy{StrategyDFS, StrategyShuffleQueue} {
		t.Run(strategy.String(), func(t *testing.T) {
			s := NewScheduler(WithStrategy(strategy))
			_, err := s.SyncAll(context.Background(), &testExecutionClient{}, schema.Tables{table})
			require.NoError(t, err)
			// SyncAll closed its channel, so sending would panic
			sent, err := checkpoint.Send(context.Background(), nil, scope, nil, &testMarker{name: "late"})
			require.ErrorIs(t, err, errResClosed)
			require.True(t, sent.Aborted())
		})
	}
}

// receiveWithin fails the test instead of hanging when nothing arrives.
func receiveWithin[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for a value")
		var zero T
		return zero
	}
}

func TestCheckpointOfEndedResolverFollowsItsRows(t *testing.T) {
	res := make(chan message.SyncMessage)
	stateClient := newCheckpointStateClient(t, res)
	endedScope := make(chan checkpoint.Scope, 1)
	ended := &schema.Table{
		Name:    checkpointParentTable,
		Columns: []schema.Column{idColumn()},
		// ends with its row still in the batcher, and the cursor not yet flushed
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			endedScope <- checkpoint.FromContext(ctx)
			res <- int64(1)
			return stateClient.SetKey(ctx, "cursor", "1")
		},
	}
	flushing := &schema.Table{
		Name:    checkpointChildTable,
		Columns: []schema.Column{idColumn()},
		Resolver: func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, _ chan<- any) error {
			scope := <-endedScope
			deadline := time.Now().Add(10 * time.Second)
			for checkpoint.Collapse(scope) == scope {
				if time.Now().After(deadline) {
					return errors.New("the other resolver did not end")
				}
				time.Sleep(time.Millisecond)
			}
			return stateClient.Flush(ctx)
		},
	}

	s := NewScheduler(WithStrategy(StrategyDFS), WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour)))
	wait := collectSyncMessages(res)
	require.NoError(t, s.Sync(context.Background(), &testExecutionClient{}, schema.Tables{ended, flushing}, res))
	close(res)
	msgs := wait()

	rowAt, cursorAt := -1, -1
	for i, msg := range msgs {
		switch m := msg.(type) {
		case *message.SyncError:
			t.Fatal(m.Error)
		case *message.SyncInsert:
			if slices.Contains(recordIDs(m.Record), 1) {
				rowAt = i
			}
		case *message.SyncCheckpoint:
			if stateValues(t, message.SyncMessages{m})["cursor"] == "1" {
				cursorAt = i
			}
		}
	}
	require.NotEqual(t, -1, cursorAt, "the cursor is sent")
	require.Less(t, rowAt, cursorAt, "the cursor follows the row of the resolver that set it")
}

func TestSyncReportsCancellationWhileTheBatcherCloses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	table := func(name string) *schema.Table {
		return &schema.Table{
			Name:    name,
			Columns: []schema.Column{idColumn()},
			Resolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
				res <- []int64{1, 2}
				return nil
			},
		}
	}
	// the batches never fill, so their rows go out only as the batcher closes. The first one cancels the sync before
	// the other can be sent.
	s := NewScheduler(WithBatchOptions(WithBatchMaxRows(1000), WithBatchTimeout(time.Hour)))
	res := make(chan message.SyncMessage)
	go func() {
		for msg := range res {
			if _, ok := msg.(*message.SyncInsert); ok {
				cancel()
			}
		}
	}()
	err := s.Sync(ctx, &testExecutionClient{}, schema.Tables{table(checkpointParentTable), table(checkpointChildTable)}, res)
	close(res)
	require.ErrorIs(t, err, context.Canceled, "the sync may have dropped rows, so it must not report success")
}

func TestReturnedWaitsForCheckpointsBeingSent(t *testing.T) {
	res := make(chan message.SyncMessage)
	c := newSyncCheckpoints(res)
	c.closeRes = true
	close(c.loopDone)
	require.NoError(t, c.complete(context.Background()))

	sent := make(chan error, 1)
	go func() {
		sent <- c.send(context.Background(), []message.SyncMessage{&testMarker{name: "late"}})
	}()
	// wait until the send holds the read lock, blocked on res
	for c.resMu.TryLock() {
		c.resMu.Unlock()
		time.Sleep(time.Millisecond)
	}
	returned := make(chan struct{})
	go func() {
		c.returned()
		close(returned)
	}()
	select {
	case <-returned:
		t.Fatal("returned did not wait for the checkpoint being sent")
	case <-time.After(100 * time.Millisecond):
	}
	require.IsType(t, &testMarker{}, receiveWithin(t, res))
	require.NoError(t, receiveWithin(t, sent))
	receiveWithin(t, returned)
	require.ErrorIs(t, c.send(context.Background(), []message.SyncMessage{&testMarker{name: "after"}}), errResClosed)
}

// newCheckpointStateClient returns a state client that sends its checkpoints on res, backed by an in-memory plugin.
func newCheckpointStateClient(t *testing.T, res chan<- message.SyncMessage) *stateV3.Client {
	t.Helper()
	p := plugin.NewPlugin("memdb", "v1.0.0", memdb.NewMemDBClient)
	require.NoError(t, p.Init(context.Background(), nil, plugin.NewClientOptions{}))
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	pb.RegisterPluginServer(server, &serversv3.Server{Plugin: p, Logger: zerolog.Nop()})
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	client, err := stateV3.NewClient(context.Background(), conn, "cq_state", stateV3.WithCheckpoints(res))
	require.NoError(t, err)
	return client
}
