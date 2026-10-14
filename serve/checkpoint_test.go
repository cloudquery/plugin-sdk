package serve

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"slices"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	pbdiscoveryv1 "github.com/cloudquery/plugin-pb-go/pb/discovery/v1"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/checkpoint"
	stateV3 "github.com/cloudquery/plugin-sdk/v4/internal/clients/state/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	discoveryServerV1 "github.com/cloudquery/plugin-sdk/v4/internal/servers/discovery/v1"
	serversv3 "github.com/cloudquery/plugin-sdk/v4/internal/servers/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/state"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// servePlugin serves p over an in-memory listener and returns a client connection to it.
func servePlugin(t *testing.T, p *plugin.Plugin) *grpc.ClientConn {
	t.Helper()
	srv := Plugin(p, WithArgs("serve"), WithTestListener())
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.Serve(ctx); err != nil {
			t.Error(err)
		}
	}()
	t.Cleanup(func() {
		cancel()
		wg.Wait()
	})

	// nolint:staticcheck
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(srv.bufPluginDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

// serveBackend serves p on a TCP address, as state clients dial backends by address.
func serveBackend(t *testing.T, p *plugin.Plugin) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	grpcServer := grpc.NewServer()
	pb.RegisterPluginServer(grpcServer, &serversv3.Server{Plugin: p, Logger: zerolog.Nop()})
	pbdiscoveryv1.RegisterDiscoveryServer(grpcServer, &discoveryServerV1.Server{Versions: []int32{3}})
	go func() { _ = grpcServer.Serve(listener) }()
	t.Cleanup(grpcServer.Stop)
	return listener.Addr().String()
}

func insertRecord(t *testing.T, table *schema.Table, ids ...int64) []byte {
	t.Helper()
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, table.ToArrowSchema())
	defer bldr.Release()
	bldr.Field(0).(*array.Int64Builder).AppendValues(ids, nil)
	recordBytes, err := pb.RecordToBytes(bldr.NewRecordBatch())
	require.NoError(t, err)
	return recordBytes
}

func TestInitReportsFlushSupport(t *testing.T) {
	for _, supported := range []bool{false, true} {
		t.Run(fmt.Sprintf("supported=%t", supported), func(t *testing.T) {
			var opts []plugin.Option
			if supported {
				opts = append(opts, plugin.WithFlushSupport())
			}
			c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("memdb", "v1.0.0", memdb.NewMemDBClient, opts...)))
			res, err := c.Init(context.Background(), &pb.Init_Request{})
			require.NoError(t, err)
			require.Equal(t, supported, res.SupportsFlush)
		})
	}
}

func TestWriteStreamFlush(t *testing.T) {
	ctx := context.Background()
	c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("memdb", "v1.0.0", memdb.NewMemDBClient, plugin.WithFlushSupport())))
	_, err := c.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)

	table := &schema.Table{Name: "test_flush", Columns: []schema.Column{{Name: "id", Type: arrow.PrimitiveTypes.Int64}}}
	tableBytes, err := pb.SchemaToBytes(table.ToArrowSchema())
	require.NoError(t, err)

	stream, err := c.WriteStream(ctx)
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.Write_Request{Message: &pb.Write_Request_MigrateTable{MigrateTable: &pb.Write_MessageMigrateTable{Table: tableBytes}}}))
	require.NoError(t, stream.Send(&pb.Write_Request{Message: &pb.Write_Request_Insert{Insert: &pb.Write_MessageInsert{Record: insertRecord(t, table, 1, 2)}}}))
	require.NoError(t, stream.Send(&pb.Write_Request{Message: &pb.Write_Request_Flush{Flush: &pb.Write_MessageFlush{Id: 7}}}))
	res, err := stream.Recv()
	require.NoError(t, err)
	require.EqualValues(t, 7, res.FlushId)

	readClient, err := c.Read(ctx, &pb.Read_Request{Table: tableBytes})
	require.NoError(t, err)
	var rows int64
	for {
		res, err := readClient.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		record, err := pb.NewRecordFromBytes(res.Record)
		require.NoError(t, err)
		rows += record.NumRows()
	}
	require.EqualValues(t, 2, rows)

	require.NoError(t, stream.CloseSend())
	_, err = stream.Recv()
	require.ErrorIs(t, err, io.EOF)
}

func TestWriteStreamWithoutFlushSupport(t *testing.T) {
	ctx := context.Background()
	c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("memdb", "v1.0.0", memdb.NewMemDBClient)))
	_, err := c.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)

	stream, err := c.WriteStream(ctx)
	require.NoError(t, err)
	_, err = stream.Recv()
	require.Equal(t, codes.Unimplemented, status.Code(err))
}

func TestWriteStreamFlushFailure(t *testing.T) {
	ctx := context.Background()
	c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("memdb", "v1.0.0", memdb.GetNewClient(memdb.WithErrOnWrite()), plugin.WithFlushSupport())))
	_, err := c.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)

	stream, err := c.WriteStream(ctx)
	require.NoError(t, err)
	// the plugin may fail before or after receiving the flush
	_ = stream.Send(&pb.Write_Request{Message: &pb.Write_Request_Flush{Flush: &pb.Write_MessageFlush{Id: 1}}})
	_, err = stream.Recv()
	require.Error(t, err)
	require.NotErrorIs(t, err, io.EOF)
}

func TestWriteRejectsFlush(t *testing.T) {
	ctx := context.Background()
	c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("memdb", "v1.0.0", memdb.NewMemDBClient, plugin.WithFlushSupport())))
	_, err := c.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)

	stream, err := c.Write(ctx)
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.Write_Request{Message: &pb.Write_Request_Flush{Flush: &pb.Write_MessageFlush{Id: 1}}}))
	_, err = stream.CloseAndRecv()
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

const (
	checkpointPages    = 3
	checkpointPageSize = 4
)

var checkpointTable = &schema.Table{
	Name: "test_checkpoint_items",
	Columns: []schema.Column{{
		Name: "id",
		Type: arrow.PrimitiveTypes.Int64,
		Resolver: func(_ context.Context, _ schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
			return resource.Set(c.Name, resource.Item.(int64))
		},
	}},
}

type checkpointSourceClient struct {
	plugin.UnimplementedDestination
	strategy scheduler.Strategy
	pages    int
}

func (*checkpointSourceClient) ID() string { return "checkpoint_source" }

func (*checkpointSourceClient) Close(context.Context) error { return nil }

func (*checkpointSourceClient) Tables(context.Context, plugin.TableOptions) (schema.Tables, error) {
	return schema.Tables{checkpointTable}, nil
}

// Sync syncs pages of items, checkpointing the page number as the cursor after each page. It sets one more key in the
// resolver and flushes it after the scheduler returns.
func (c *checkpointSourceClient) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	stateClient, err := state.NewConnectedClientWithOptions(ctx, options.BackendOptions, state.ConnectionOptions{}, state.ClientOptions{Checkpoints: res})
	if err != nil {
		return err
	}
	defer stateClient.Close()

	table := *checkpointTable
	table.Resolver = func(ctx context.Context, _ schema.ClientMeta, _ *schema.Resource, res chan<- any) error {
		for page := range c.pages {
			ids := make([]int64, checkpointPageSize)
			for i := range ids {
				ids[i] = int64(page*checkpointPageSize + i)
			}
			res <- ids
			if err := stateClient.SetKey(ctx, "cursor", strconv.Itoa(page)); err != nil {
				return err
			}
			if err := stateClient.Flush(ctx); err != nil {
				return err
			}
		}
		return stateClient.SetKey(ctx, "finished", "true")
	}
	s := scheduler.NewScheduler(scheduler.WithStrategy(c.strategy), scheduler.WithBatchOptions(scheduler.WithBatchMaxRows(1000), scheduler.WithBatchTimeout(time.Hour)))
	if err := s.Sync(ctx, c, schema.Tables{&table}, res); err != nil {
		return err
	}
	return stateClient.Flush(ctx)
}

func newCheckpointSource(ctx context.Context, logger zerolog.Logger, spec []byte, opts plugin.NewClientOptions) (plugin.Client, error) {
	return newCheckpointSourceWith(scheduler.StrategyDFS, checkpointPages)(ctx, logger, spec, opts)
}

func newCheckpointSourceWith(strategy scheduler.Strategy, pages int) plugin.NewClientFunc {
	return func(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
		return &checkpointSourceClient{strategy: strategy, pages: pages}, nil
	}
}

func readState(t *testing.T, conn *grpc.ClientConn, tableName string) map[string]string {
	t.Helper()
	tableBytes, err := pb.SchemaToBytes(stateV3.Table(tableName).ToArrowSchema())
	require.NoError(t, err)
	readClient, err := pb.NewPluginClient(conn).Read(context.Background(), &pb.Read_Request{Table: tableBytes})
	require.NoError(t, err)
	values := map[string]string{}
	for {
		res, err := readClient.Recv()
		if errors.Is(err, io.EOF) {
			return values
		}
		require.NoError(t, err)
		record, err := pb.NewRecordFromBytes(res.Record)
		require.NoError(t, err)
		addStateRows(values, record)
	}
}

func addStateRows(values map[string]string, record arrow.RecordBatch) {
	keys, vals := record.Column(0).(*array.String), record.Column(1).(*array.String)
	for i := range keys.Len() {
		values[keys.Value(i)] = vals.Value(i)
	}
}

func TestSyncCheckpoints(t *testing.T) {
	for _, supportsCheckpoints := range []bool{true, false} {
		t.Run(fmt.Sprintf("supports_checkpoints=%t", supportsCheckpoints), func(t *testing.T) {
			ctx := context.Background()
			backendPlugin := plugin.NewPlugin("memdb", "v1.0.0", memdb.NewMemDBClient, plugin.WithFlushSupport())
			require.NoError(t, backendPlugin.Init(ctx, nil, plugin.NewClientOptions{}))
			backendAddr := serveBackend(t, backendPlugin)

			c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("checkpoint_source", "v1.0.0", newCheckpointSource)))
			_, err := c.Init(ctx, &pb.Init_Request{})
			require.NoError(t, err)

			syncClient, err := c.Sync(ctx, &pb.Sync_Request{
				Backend:             &pb.Sync_BackendOptions{TableName: "cq_state", Connection: backendAddr},
				SupportsCheckpoints: supportsCheckpoints,
			})
			require.NoError(t, err)

			sent := map[int64]bool{}
			checkpointed := map[string]string{}
			for {
				res, err := syncClient.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				require.NoError(t, err)
				switch m := res.Message.(type) {
				case *pb.Sync_Response_Insert:
					record, err := pb.NewRecordFromBytes(m.Insert.Record)
					require.NoError(t, err)
					for _, id := range record.Column(record.Schema().FieldIndices("id")[0]).(*array.Int64).Int64Values() {
						sent[id] = true
					}
				case *pb.Sync_Response_Checkpoint:
					require.True(t, supportsCheckpoints)
					record, err := pb.NewRecordFromBytes(m.Checkpoint.Record)
					require.NoError(t, err)
					values := map[string]string{}
					addStateRows(values, record)
					if cursor, ok := values["cursor"]; ok {
						page, err := strconv.Atoi(cursor)
						require.NoError(t, err)
						for id := range int64((page + 1) * checkpointPageSize) {
							require.Truef(t, sent[id], "checkpoint for page %d sent before row %d", page, id)
						}
					}
					if _, ok := values["finished"]; ok {
						require.Len(t, sent, checkpointPages*checkpointPageSize)
					}
					for k, v := range values {
						checkpointed[k] = v
					}
				}
			}
			require.Len(t, sent, checkpointPages*checkpointPageSize)

			backendState := readState(t, servePluginConn(t, backendAddr), "cq_state")
			want := map[string]string{"cursor": strconv.Itoa(checkpointPages - 1), "finished": "true"}
			if supportsCheckpoints {
				require.Equal(t, want, checkpointed)
				require.Empty(t, backendState, "the client persists checkpoints")
			} else {
				require.Empty(t, checkpointed)
				require.Equal(t, want, backendState)
			}
		})
	}
}

func servePluginConn(t *testing.T, addr string) *grpc.ClientConn {
	t.Helper()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })
	return conn
}

type recordingScope struct {
	msgs []message.SyncMessage
	err  error
}

func (s *recordingScope) Checkpoint(_ context.Context, msg message.SyncMessage) error {
	if s.err != nil {
		return s.err
	}
	s.msgs = append(s.msgs, msg)
	return nil
}

func checkpointValues(t *testing.T, msgs ...message.SyncMessage) map[string]string {
	t.Helper()
	values := map[string]string{}
	for _, msg := range msgs {
		if tracked, ok := msg.(*checkpoint.Tracked); ok {
			msg = tracked.Message
		}
		require.IsType(t, &message.SyncCheckpoint{}, msg)
		addStateRows(values, msg.(*message.SyncCheckpoint).Record)
	}
	return values
}

// deliver sends what a scope received as the scheduler would, once the scope's resources are sent.
func deliver(t *testing.T, res chan<- message.SyncMessage, msgs ...message.SyncMessage) {
	t.Helper()
	for _, msg := range msgs {
		require.NoError(t, checkpoint.Deliver(context.Background(), res, msg))
	}
}

func newCheckpointStateClient(t *testing.T, tableName string) (*stateV3.Client, chan message.SyncMessage, *grpc.ClientConn) {
	t.Helper()
	ctx := context.Background()
	conn := servePlugin(t, plugin.NewPlugin("memdb", "v1.0.0", memdb.NewMemDBClient))
	_, err := pb.NewPluginClient(conn).Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)
	res := make(chan message.SyncMessage, 10)
	stateClient, err := stateV3.NewClient(ctx, conn, tableName, stateV3.WithCheckpoints(res))
	require.NoError(t, err)
	return stateClient, res, conn
}

func TestStateClientCheckpointsPerScope(t *testing.T) {
	ctx := context.Background()
	stateClient, res, conn := newCheckpointStateClient(t, "test_scopes")

	scopeA, scopeB := &recordingScope{}, &recordingScope{}
	ctxA, ctxB := checkpoint.WithScope(ctx, scopeA), checkpoint.WithScope(ctx, scopeB)
	require.NoError(t, stateClient.SetKey(ctxA, "a", "1"))
	require.NoError(t, stateClient.SetKey(ctxB, "b", "1"))
	require.NoError(t, stateClient.SetKey(ctx, "c", "1"))
	require.NoError(t, stateClient.SetKey(ctxA, "shared", "1"))
	require.NoError(t, stateClient.SetKey(ctxB, "shared", "2"))

	require.NoError(t, stateClient.Flush(ctxA))
	require.Len(t, res, 1)
	require.Equal(t, map[string]string{"c": "1"}, checkpointValues(t, <-res))

	// the shared key waits for both scopes: each gets a barrier, and the last one to reach it sends the checkpoint
	require.Len(t, scopeA.msgs, 2)
	require.Len(t, scopeB.msgs, 2)
	var barriers []message.SyncMessage
	for _, scope := range []*recordingScope{scopeA, scopeB} {
		var tracked []message.SyncMessage
		for _, msg := range scope.msgs {
			if _, ok := msg.(*checkpoint.Barrier); ok {
				barriers = append(barriers, msg)
			} else {
				tracked = append(tracked, msg)
			}
		}
		require.Len(t, tracked, 1)
	}
	require.Equal(t, map[string]string{"a": "1"}, checkpointValues(t, scopeA.msgs[slices.IndexFunc(scopeA.msgs, isCheckpoint)]))
	require.Equal(t, map[string]string{"b": "1"}, checkpointValues(t, scopeB.msgs[slices.IndexFunc(scopeB.msgs, isCheckpoint)]))
	require.Len(t, barriers, 2)
	deliver(t, res, barriers[0])
	require.Empty(t, res)
	deliver(t, res, barriers[1])
	require.Equal(t, map[string]string{"shared": "2"}, checkpointValues(t, <-res))

	require.NoError(t, stateClient.Flush(ctx))
	require.Len(t, scopeA.msgs, 2)
	require.Len(t, scopeB.msgs, 2)
	require.Empty(t, res)

	require.Empty(t, readState(t, conn, "test_scopes"), "checkpoints are not written to the backend")
}

func isCheckpoint(msg message.SyncMessage) bool {
	_, ok := msg.(*checkpoint.Tracked)
	return ok
}

func TestStateClientCheckpointWaitsForPreviousCheckpointOfKey(t *testing.T) {
	ctx := context.Background()
	stateClient, res, _ := newCheckpointStateClient(t, "test_previous")

	scopeA, scopeB := &recordingScope{}, &recordingScope{}
	require.NoError(t, stateClient.SetKey(checkpoint.WithScope(ctx, scopeA), "cursor", "1"))
	require.NoError(t, stateClient.Flush(ctx))
	require.NoError(t, stateClient.SetKey(checkpoint.WithScope(ctx, scopeB), "cursor", "2"))
	require.NoError(t, stateClient.Flush(ctx))
	require.Len(t, scopeA.msgs, 1)
	require.Len(t, scopeB.msgs, 1)

	// scope B finishes first, but its checkpoint must not overtake scope A's resources
	deliver(t, res, scopeB.msgs...)
	require.Empty(t, res)
	deliver(t, res, scopeA.msgs...)
	require.Equal(t, map[string]string{"cursor": "1"}, checkpointValues(t, <-res))
	require.Equal(t, map[string]string{"cursor": "2"}, checkpointValues(t, <-res))

	// once sent, a checkpoint holds back nothing
	scopeC := &recordingScope{}
	require.NoError(t, stateClient.SetKey(checkpoint.WithScope(ctx, scopeC), "cursor", "3"))
	require.NoError(t, stateClient.Flush(ctx))
	require.Len(t, scopeC.msgs, 1)
	require.IsType(t, &checkpoint.Tracked{}, scopeC.msgs[0])
}

func TestStateClientKeepsKeysOfFailedCheckpoints(t *testing.T) {
	ctx := context.Background()
	stateClient, res, _ := newCheckpointStateClient(t, "test_failed_scopes")

	aborted := &recordingScope{err: checkpoint.ErrSyncAborted}
	require.NoError(t, stateClient.SetKey(checkpoint.WithScope(ctx, aborted), "a", "1"))
	require.NoError(t, stateClient.SetKey(ctx, "b", "1"))

	// the failed scope does not stop other keys from being sent
	require.ErrorIs(t, stateClient.Flush(ctx), state.ErrSyncAborted)
	require.Equal(t, map[string]string{"b": "1"}, checkpointValues(t, <-res))

	aborted.err = nil
	require.NoError(t, stateClient.Flush(ctx))
	require.Equal(t, map[string]string{"a": "1"}, checkpointValues(t, aborted.msgs...))
	require.Empty(t, res)
}

// heldFlushClient hands each flush to the test instead of completing it.
type heldFlushClient struct {
	plugin.UnimplementedDestination
	flushes chan *message.WriteFlush
}

func (*heldFlushClient) Close(context.Context) error { return nil }
func (*heldFlushClient) Tables(context.Context, plugin.TableOptions) (schema.Tables, error) {
	return nil, nil
}
func (*heldFlushClient) Sync(context.Context, plugin.SyncOptions, chan<- message.SyncMessage) error {
	return plugin.ErrNotImplemented
}

func (c *heldFlushClient) Write(_ context.Context, msgs <-chan message.WriteMessage) error {
	for msg := range msgs {
		if flush, ok := msg.(*message.WriteFlush); ok {
			c.flushes <- flush
		}
	}
	return nil
}

func TestWriteStreamAcksFlushOnlyAfterCommit(t *testing.T) {
	ctx := context.Background()
	client := &heldFlushClient{flushes: make(chan *message.WriteFlush)}
	c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("held", "v1.0.0", func(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
		return client, nil
	}, plugin.WithFlushSupport())))
	_, err := c.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)

	stream, err := c.WriteStream(ctx)
	require.NoError(t, err)
	responses := make(chan *pb.Write_Response)
	errs := make(chan error, 1)
	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				errs <- err
				return
			}
			responses <- res
		}
	}()

	require.NoError(t, stream.Send(&pb.Write_Request{Message: &pb.Write_Request_Flush{Flush: &pb.Write_MessageFlush{Id: 3}}}))
	flush := receive(t, client.flushes)
	select {
	case res := <-responses:
		t.Fatalf("flush %d acknowledged before it completed", res.FlushId)
	case err := <-errs:
		t.Fatal(err)
	case <-time.After(200 * time.Millisecond):
	}
	flush.Done(nil)
	require.EqualValues(t, 3, receive(t, responses).FlushId)

	require.NoError(t, stream.Send(&pb.Write_Request{Message: &pb.Write_Request_Flush{Flush: &pb.Write_MessageFlush{Id: 4}}}))
	receive(t, client.flushes).Done(errors.New("disk full"))
	select {
	case err := <-errs:
		require.Equal(t, codes.Internal, status.Code(err))
		require.ErrorContains(t, err, "flush failed: disk full")
	case res := <-responses:
		t.Fatalf("failed flush %d acknowledged", res.FlushId)
	case <-time.After(5 * time.Second):
		t.Fatal("stream did not fail after the flush failed")
	}
}

func TestWriteStreamRejectsEmptyMessage(t *testing.T) {
	ctx := context.Background()
	c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("memdb", "v1.0.0", memdb.NewMemDBClient, plugin.WithFlushSupport())))
	_, err := c.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)

	stream, err := c.WriteStream(ctx)
	require.NoError(t, err)
	require.NoError(t, stream.Send(&pb.Write_Request{}))
	_, err = stream.Recv()
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

type noFlushClient struct {
	plugin.Client
}

func (noFlushClient) SupportsFlush() bool { return false }

func TestClientCanTurnOffFlushSupport(t *testing.T) {
	ctx := context.Background()
	c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("memdb", "v1.0.0", func(ctx context.Context, logger zerolog.Logger, spec []byte, opts plugin.NewClientOptions) (plugin.Client, error) {
		client, err := memdb.NewMemDBClient(ctx, logger, spec, opts)
		return noFlushClient{Client: client}, err
	}, plugin.WithFlushSupport())))
	res, err := c.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)
	require.False(t, res.SupportsFlush)

	stream, err := c.WriteStream(ctx)
	require.NoError(t, err)
	_, err = stream.Recv()
	require.Equal(t, codes.Unimplemented, status.Code(err))
}

func TestShuffleQueueCoalescesFlushedCheckpoints(t *testing.T) {
	const pages = 100
	ctx := context.Background()
	backendPlugin := plugin.NewPlugin("memdb", "v1.0.0", memdb.NewMemDBClient, plugin.WithFlushSupport())
	require.NoError(t, backendPlugin.Init(ctx, nil, plugin.NewClientOptions{}))
	backendAddr := serveBackend(t, backendPlugin)

	c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("checkpoint_source", "v1.0.0", newCheckpointSourceWith(scheduler.StrategyShuffleQueue, pages))))
	_, err := c.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)
	syncClient, err := c.Sync(ctx, &pb.Sync_Request{
		Backend:             &pb.Sync_BackendOptions{TableName: "cq_state", Connection: backendAddr},
		SupportsCheckpoints: true,
	})
	require.NoError(t, err)

	var rows int64
	var cursors []string
	for {
		res, err := syncClient.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		switch m := res.Message.(type) {
		case *pb.Sync_Response_Insert:
			record, err := pb.NewRecordFromBytes(m.Insert.Record)
			require.NoError(t, err)
			rows += record.NumRows()
		case *pb.Sync_Response_Checkpoint:
			record, err := pb.NewRecordFromBytes(m.Checkpoint.Record)
			require.NoError(t, err)
			values := map[string]string{}
			addStateRows(values, record)
			if cursor, ok := values["cursor"]; ok {
				require.EqualValues(t, pages*checkpointPageSize, rows, "the cursor follows every row")
				cursors = append(cursors, cursor)
			}
		}
	}
	require.Equal(t, []string{strconv.Itoa(pages - 1)}, cursors, "flushes of the cursor are merged into one checkpoint")
}

func TestStateClientJoinsScopesOfSameValueWrites(t *testing.T) {
	ctx := context.Background()
	stateClient, res, _ := newCheckpointStateClient(t, "test_same_value")

	scopeA, scopeB := &recordingScope{}, &recordingScope{}
	require.NoError(t, stateClient.SetKey(checkpoint.WithScope(ctx, scopeA), "cursor", "1"))
	// the same value from another resolver summarizes its resources too
	require.NoError(t, stateClient.SetKey(checkpoint.WithScope(ctx, scopeB), "cursor", "1"))
	require.NoError(t, stateClient.Flush(ctx))

	require.Len(t, scopeA.msgs, 1)
	require.Len(t, scopeB.msgs, 1)
	deliver(t, res, scopeA.msgs...)
	require.Empty(t, res)
	deliver(t, res, scopeB.msgs...)
	require.Equal(t, map[string]string{"cursor": "1"}, checkpointValues(t, <-res))
}

func TestStateClientKeepsKeysWhoseCheckpointWasDropped(t *testing.T) {
	ctx := context.Background()
	stateClient, res, _ := newCheckpointStateClient(t, "test_dropped")

	scope := &recordingScope{}
	require.NoError(t, stateClient.SetKey(checkpoint.WithScope(ctx, scope), "cursor", "1"))
	require.NoError(t, stateClient.Flush(ctx))
	require.Len(t, scope.msgs, 1)

	// the sync is cancelled before the checkpoint is sent
	checkpoint.Drop(scope.msgs[0])
	require.ErrorIs(t, stateClient.Flush(ctx), state.ErrSyncAborted)
	require.NoError(t, stateClient.SetKey(ctx, "cursor", "2"))
	require.ErrorIs(t, stateClient.Flush(ctx), state.ErrSyncAborted)
	require.Empty(t, res)
}

// collapsingScope stands for a table resolver's scope, which collapses into the sync's once the resolver ends.
type collapsingScope struct {
	recordingScope
	into checkpoint.Scope
}

func (s *collapsingScope) Collapsed() checkpoint.Scope { return s.into }

func TestStateClientCollapsesEndedScopes(t *testing.T) {
	ctx := context.Background()
	stateClient, res, _ := newCheckpointStateClient(t, "test_collapsed")

	shared := &recordingScope{}
	for i := range 1000 {
		scope := &collapsingScope{into: shared}
		require.NoError(t, stateClient.SetKey(checkpoint.WithScope(ctx, scope), "cursor", strconv.Itoa(i)))
	}
	require.NoError(t, stateClient.Flush(ctx))
	require.Len(t, shared.msgs, 1, "one checkpoint through the shared scope")
	deliver(t, res, shared.msgs...)
	require.Equal(t, map[string]string{"cursor": "999"}, checkpointValues(t, <-res))
}

// receive fails the test instead of hanging when nothing arrives.
func receive[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for a value")
		var zero T
		return zero
	}
}

// stoppingClient's Write returns without waiting for its messages to end, after taking takeFirst messages.
type stoppingClient struct {
	heldFlushClient
	takeFirst int
}

func (c *stoppingClient) Write(_ context.Context, msgs <-chan message.WriteMessage) error {
	for range c.takeFirst {
		<-msgs
	}
	return nil
}

func TestWriteStreamWhenThePluginStopsWriting(t *testing.T) {
	for _, tc := range []struct {
		name      string
		takeFirst int
		err       string
	}{
		{name: "before taking the flush", takeFirst: 0, err: "plugin stopped writing"},
		{name: "before completing the flush", takeFirst: 1, err: "plugin stopped writing before the flush completed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			client := &stoppingClient{takeFirst: tc.takeFirst}
			c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("stopping", "v1.0.0", func(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
				return client, nil
			}, plugin.WithFlushSupport())))
			_, err := c.Init(ctx, &pb.Init_Request{})
			require.NoError(t, err)

			stream, err := c.WriteStream(ctx)
			require.NoError(t, err)
			require.NoError(t, stream.Send(&pb.Write_Request{Message: &pb.Write_Request_Flush{Flush: &pb.Write_MessageFlush{Id: 1}}}))
			errs := make(chan error, 1)
			go func() {
				_, err := stream.Recv()
				errs <- err
			}()
			err = receive(t, errs)
			require.Equal(t, codes.Internal, status.Code(err))
			require.ErrorContains(t, err, tc.err)
		})
	}
}

func TestInitDoesNotReportFlushSupportWithoutConnection(t *testing.T) {
	c := pb.NewPluginClient(servePlugin(t, plugin.NewPlugin("memdb", "v1.0.0", memdb.NewMemDBClient, plugin.WithFlushSupport())))
	res, err := c.Init(context.Background(), &pb.Init_Request{NoConnection: true})
	require.NoError(t, err)
	require.False(t, res.SupportsFlush)
}
