package serve

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type watchTestMeta struct{}

func (*watchTestMeta) ID() string { return "watch_test_client" }

func watchTestTable() *schema.Table {
	table := &schema.Table{
		Name: "watch_test_table",
		Resolver: func(_ context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
			res <- map[string]any{"id": "full-sync", "name": "from full sync"}
			return nil
		},
		Columns: schema.ColumnList{
			{Name: "id", Type: arrow.BinaryTypes.String, PrimaryKey: true, Resolver: schema.PathResolver("id")},
			{Name: "name", Type: arrow.BinaryTypes.String, Resolver: schema.PathResolver("name")},
		},
	}
	schema.AddCqIDs(table)
	return table
}

// watchTestClient is a source client that implements both optional event-sync
// interfaces: it replays canned changes over Watch and answers Refresh from a
// canned map.
type watchTestClient struct {
	plugin.UnimplementedDestination

	tables schema.Tables
	sched  *scheduler.Scheduler

	// changes are emitted by Watch, in order.
	changes []plugin.Change
	// follow keeps Watch open until the context is done instead of returning
	// once the changes have been emitted.
	follow bool
	// refreshItems maps schema.KeysString(keys) to the items Refresh returns.
	// A key that is absent is reported as gone.
	refreshItems map[string][]any

	mu           sync.Mutex
	acked        []string
	watchStarted chan struct{}
	readyCalls   []plugin.WatchOptions
}

func newWatchTestClient() *watchTestClient {
	return &watchTestClient{
		tables:       schema.Tables{watchTestTable()},
		sched:        scheduler.NewScheduler(),
		refreshItems: map[string][]any{},
		watchStarted: make(chan struct{}),
	}
}

func (*watchTestClient) Close(context.Context) error { return nil }

func (c *watchTestClient) Tables(_ context.Context, options plugin.TableOptions) (schema.Tables, error) {
	return c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
}

func (c *watchTestClient) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	if len(options.Refresh) > 0 {
		return c.sched.Refresh(ctx, c.tables, c, options.Refresh, res)
	}
	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return err
	}
	return c.sched.Sync(ctx, &watchTestMeta{}, tt, res, scheduler.WithSyncDeterministicCQID(options.DeterministicCQID))
}

func (c *watchTestClient) Refresh(_ context.Context, _ *schema.Table, keys arrow.RecordBatch) (schema.ClientMeta, []any, error) {
	items, ok := c.refreshItems[schema.KeysString(keys)]
	if !ok {
		return nil, nil, scheduler.ErrNotFound
	}
	return &watchTestMeta{}, items, nil
}

func (c *watchTestClient) WatchReady(_ context.Context, opts plugin.WatchOptions) (plugin.Ready, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.readyCalls = append(c.readyCalls, opts)
	return plugin.Ready{HasCheckpoint: true, CheckpointTime: time.Unix(1700000000, 0).UTC()}, nil
}

func (c *watchTestClient) Watch(ctx context.Context, _ plugin.Position, changes chan<- plugin.Change) error {
	close(c.watchStarted)
	for _, change := range c.changes {
		select {
		case changes <- change:
		case <-ctx.Done():
			return nil
		}
	}
	if c.follow {
		<-ctx.Done()
	}
	return nil
}

func (c *watchTestClient) Ack(_ context.Context, eventID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.acked = append(c.acked, eventID)
	return nil
}

func (c *watchTestClient) ackedEvents() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]string(nil), c.acked...)
}

func (c *watchTestClient) watchOptions() []plugin.WatchOptions {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]plugin.WatchOptions(nil), c.readyCalls...)
}

// serveTestPlugin starts the plugin over a bufconn listener and returns a
// connected gRPC client. The server is stopped when the test ends.
func serveTestPlugin(t *testing.T, newClient plugin.NewClientFunc) pb.PluginClient {
	t.Helper()

	p := plugin.NewPlugin("test_watch_plugin", "v1.0.0", newClient)
	p.SetLogger(zerolog.Nop())
	srv := Plugin(p, WithArgs("serve"), WithTestListener())

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	var serveErr error
	go func() {
		defer wg.Done()
		serveErr = srv.Serve(ctx)
	}()
	t.Cleanup(func() {
		cancel()
		wg.Wait()
		require.NoError(t, serveErr)
	})

	// nolint:staticcheck // bufconn dialing still needs the deprecated DialContext
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(srv.bufPluginDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	c := pb.NewPluginClient(conn)
	_, err = c.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)
	return c
}

func serveWatchTestClient(t *testing.T, cl *watchTestClient) pb.PluginClient {
	t.Helper()
	return serveTestPlugin(t, func(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
		return cl, nil
	})
}

func watchTestKeys(t *testing.T, id string) arrow.RecordBatch {
	t.Helper()
	keys, err := schema.NewKeysRecord(watchTestTable(), map[string]any{"id": id})
	require.NoError(t, err)
	return keys
}

func TestWatch_HandshakeAndDrain(t *testing.T) {
	cl := newWatchTestClient()
	cl.changes = []plugin.Change{
		{Table: "watch_test_table", Keys: watchTestKeys(t, "a"), Op: plugin.OpUpsert, EventID: "e1", Time: time.Unix(1700000001, 0).UTC()},
		{Table: "watch_test_table", Keys: watchTestKeys(t, "b"), Op: plugin.OpDelete, EventID: "e2", Time: time.Unix(1700000002, 0).UTC()},
	}
	c := serveWatchTestClient(t, cl)
	ctx := context.Background()

	stream, err := c.Watch(ctx)
	require.NoError(t, err)

	require.NoError(t, stream.Send(&pb.Watch_Request{Message: &pb.Watch_Request_Start{Start: &pb.Watch_Start{
		Tables:  []string{"*"},
		Drain:   true,
		Backend: &pb.Sync_BackendOptions{TableName: "cq_state", Connection: "@@plugin.test.connection@@"},
	}}}))

	resp, err := stream.Recv()
	require.NoError(t, err)
	ready, ok := resp.Message.(*pb.Watch_Response_Ready)
	require.True(t, ok, "expected Ready first, got %T", resp.Message)
	require.True(t, ready.Ready.HasCheckpoint)
	require.Equal(t, time.Unix(1700000000, 0).UTC(), ready.Ready.CheckpointTime.AsTime())

	opts := cl.watchOptions()
	require.Len(t, opts, 1)
	require.Equal(t, []string{"*"}, opts[0].Tables)
	require.True(t, opts[0].Drain)
	require.NotNil(t, opts[0].BackendOptions)
	require.Equal(t, "cq_state", opts[0].BackendOptions.TableName)

	require.NoError(t, stream.Send(&pb.Watch_Request{Message: &pb.Watch_Request_Position{Position: &pb.Watch_Position{}}}))

	var got []*pb.Watch_Change
	drained := false
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		switch m := resp.Message.(type) {
		case *pb.Watch_Response_Change:
			got = append(got, m.Change)
		case *pb.Watch_Response_Drained:
			drained = true
		default:
			t.Fatalf("unexpected message %T", resp.Message)
		}
		if drained {
			break
		}
	}

	require.True(t, drained, "expected a Drained message")
	require.Len(t, got, 2)
	require.Equal(t, "e1", got[0].EventId)
	require.Equal(t, pb.Op_OP_UPSERT, got[0].Op)
	require.Equal(t, time.Unix(1700000001, 0).UTC(), got[0].Time.AsTime())
	require.Equal(t, "e2", got[1].EventId)
	require.Equal(t, pb.Op_OP_DELETE, got[1].Op)

	keys, err := pb.NewRecordFromBytes(got[0].Keys)
	require.NoError(t, err)
	require.Equal(t, "id=a", schema.KeysString(keys))
}

func TestWatch_AckIsForwarded(t *testing.T) {
	cl := newWatchTestClient()
	cl.follow = true
	cl.changes = []plugin.Change{
		{Table: "watch_test_table", Keys: watchTestKeys(t, "a"), Op: plugin.OpUpsert, EventID: "e1"},
	}
	c := serveWatchTestClient(t, cl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.Watch(ctx)
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.Watch_Request{Message: &pb.Watch_Request_Start{Start: &pb.Watch_Start{Tables: []string{"*"}}}}))
	_, err = stream.Recv()
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.Watch_Request{Message: &pb.Watch_Request_Position{Position: &pb.Watch_Position{}}}))

	resp, err := stream.Recv()
	require.NoError(t, err)
	change, ok := resp.Message.(*pb.Watch_Response_Change)
	require.True(t, ok)

	require.NoError(t, stream.Send(&pb.Watch_Request{Message: &pb.Watch_Request_Ack{Ack: &pb.Watch_Ack{EventId: change.Change.EventId}}}))

	require.Eventually(t, func() bool {
		return len(cl.ackedEvents()) == 1
	}, 5*time.Second, 10*time.Millisecond, "expected the ack to reach the client")
	require.Equal(t, []string{"e1"}, cl.ackedEvents())
}

func TestWatch_FirstMessageMustBeStart(t *testing.T) {
	c := serveWatchTestClient(t, newWatchTestClient())

	stream, err := c.Watch(context.Background())
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.Watch_Request{Message: &pb.Watch_Request_Position{Position: &pb.Watch_Position{}}}))

	_, err = stream.Recv()
	require.Error(t, err)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Contains(t, status.Convert(err).Message(), "must be Start")
}

func TestWatch_UnimplementedWithoutWatchClient(t *testing.T) {
	c := serveTestPlugin(t, memdb.NewMemDBClient)

	stream, err := c.Watch(context.Background())
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.Watch_Request{Message: &pb.Watch_Request_Start{Start: &pb.Watch_Start{Tables: []string{"*"}}}}))

	_, err = stream.Recv()
	require.Error(t, err)
	require.Equal(t, codes.Unimplemented, status.Code(err))
}

// syncRefresh runs a refresh sync and collects the messages it produces.
func syncRefresh(ctx context.Context, t *testing.T, c pb.PluginClient, resources []*pb.Sync_Resource) []*pb.Sync_Response {
	t.Helper()
	syncClient, err := c.Sync(ctx, &pb.Sync_Request{Refresh: resources})
	require.NoError(t, err)

	var msgs []*pb.Sync_Response
	for {
		resp, err := syncClient.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		msgs = append(msgs, resp)
	}
	return msgs
}

func watchTestResource(t *testing.T, id string, op pb.Op) *pb.Sync_Resource {
	t.Helper()
	keys, err := pb.RecordToBytes(watchTestKeys(t, id))
	require.NoError(t, err)
	return &pb.Sync_Resource{Table: "watch_test_table", Keys: keys, Op: op}
}

func TestSyncRefresh_InsertAndDelete(t *testing.T) {
	cl := newWatchTestClient()
	cl.refreshItems["id=here"] = []any{map[string]any{"id": "here", "name": "still around"}}
	c := serveWatchTestClient(t, cl)

	msgs := syncRefresh(context.Background(), t, c, []*pb.Sync_Resource{
		watchTestResource(t, "here", pb.Op_OP_UPSERT),
		watchTestResource(t, "gone", pb.Op_OP_UPSERT),
	})

	var inserts, deletes, migrates int
	for _, msg := range msgs {
		switch m := msg.Message.(type) {
		case *pb.Sync_Response_Insert:
			inserts++
			record, err := pb.NewRecordFromBytes(m.Insert.Record)
			require.NoError(t, err)
			require.EqualValues(t, 1, record.NumRows())
		case *pb.Sync_Response_DeleteRecord:
			deletes++
			require.Equal(t, "watch_test_table", m.DeleteRecord.TableName)
			require.Len(t, m.DeleteRecord.WhereClause, 1)
			require.Equal(t, pb.PredicatesGroup_AND, m.DeleteRecord.WhereClause[0].GroupingType)
			require.Len(t, m.DeleteRecord.WhereClause[0].Predicates, 1)
			require.Equal(t, pb.Predicate_EQ, m.DeleteRecord.WhereClause[0].Predicates[0].Operator)
			require.Equal(t, "id", m.DeleteRecord.WhereClause[0].Predicates[0].Column)
		case *pb.Sync_Response_MigrateTable:
			migrates++
		}
	}

	require.Equal(t, 1, inserts, "expected one insert for the resource that still exists")
	require.Equal(t, 1, deletes, "expected one delete for the resource that is gone")
	require.Zero(t, migrates, "a refresh sync must not emit migrate messages")
}

func TestSyncRefresh_WhileWatchStreamIsOpen(t *testing.T) {
	cl := newWatchTestClient()
	cl.follow = true
	cl.refreshItems["id=here"] = []any{map[string]any{"id": "here", "name": "still around"}}
	c := serveWatchTestClient(t, cl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.Watch(ctx)
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.Watch_Request{Message: &pb.Watch_Request_Start{Start: &pb.Watch_Start{Tables: []string{"*"}}}}))
	_, err = stream.Recv()
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.Watch_Request{Message: &pb.Watch_Request_Position{Position: &pb.Watch_Position{}}}))

	// The watch is running: a refresh sync on the same plugin must still go through.
	select {
	case <-cl.watchStarted:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for the watch to start")
	}

	msgs := syncRefresh(ctx, t, c, []*pb.Sync_Resource{watchTestResource(t, "here", pb.Op_OP_UPSERT)})

	var inserts int
	for _, msg := range msgs {
		if _, ok := msg.Message.(*pb.Sync_Response_Insert); ok {
			inserts++
		}
	}
	require.Equal(t, 1, inserts, "refresh sync must succeed while a watch stream is open")
}
