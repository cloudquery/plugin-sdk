package plugin

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestGetName(t *testing.T) {
	ctx := context.Background()
	s := Server{
		Plugin: plugin.NewPlugin("test", "development", memdb.NewMemDBClient),
	}
	res, err := s.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Name != "test" {
		t.Fatalf("expected test, got %s", res.GetName())
	}
}

func TestGetVersion(t *testing.T) {
	ctx := context.Background()
	s := Server{
		Plugin: plugin.NewPlugin("test", "development", memdb.NewMemDBClient),
	}
	resVersion, err := s.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if resVersion.Version != "development" {
		t.Fatalf("expected development, got %s", resVersion.GetVersion())
	}
}

func TestGetTables(t *testing.T) {
	ctx := context.Background()
	pluginVersion := "v1.2.3"
	s := Server{
		Plugin: plugin.NewPlugin("test", pluginVersion, memdb.NewMemDBClient),
	}

	_, err := s.Init(ctx, &pb.Init_Request{})
	require.NoError(t, err)

	res, err := s.GetTables(ctx, &pb.GetTables_Request{Tables: []string{"*"}})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Greater(t, len(res.Tables), 0, "expected at least one table")

	// Verify that the plugin version is included in the schema metadata
	for _, tableBytes := range res.Tables {
		sc, err := pb.NewSchemaFromBytes(tableBytes)
		require.NoError(t, err)

		version, found := sc.Metadata().GetValue(schema.MetadataTablePluginVersion)
		require.True(t, found, "expected plugin version to be in schema metadata")
		require.Equal(t, pluginVersion, version, "expected plugin version to match")
	}
}

type mockSyncServer struct {
	grpc.ServerStream
	messages []*pb.Sync_Response
}

func (s *mockSyncServer) Send(*pb.Sync_Response) error {
	s.messages = append(s.messages, &pb.Sync_Response{})
	return nil
}

func (*mockSyncServer) SetHeader(metadata.MD) error {
	return nil
}
func (*mockSyncServer) SendHeader(metadata.MD) error {
	return nil
}
func (*mockSyncServer) SetTrailer(metadata.MD) {
}
func (*mockSyncServer) Context() context.Context {
	return context.Background()
}
func (*mockSyncServer) SendMsg(any) error {
	return nil
}
func (*mockSyncServer) RecvMsg(any) error {
	return nil
}

type mockWriteServer struct {
	grpc.ServerStream
	messages []*pb.Write_Request
}

func (*mockWriteServer) SendAndClose(*pb.Write_Response) error {
	return nil
}
func (s *mockWriteServer) Recv() (*pb.Write_Request, error) {
	if len(s.messages) > 0 {
		msg := s.messages[0]
		s.messages = s.messages[1:]
		return msg, nil
	}
	return nil, io.EOF
}
func (*mockWriteServer) SetHeader(metadata.MD) error {
	return nil
}
func (*mockWriteServer) SendHeader(metadata.MD) error {
	return nil
}
func (*mockWriteServer) SetTrailer(metadata.MD) {
}
func (*mockWriteServer) Context() context.Context {
	return context.Background()
}
func (*mockWriteServer) SendMsg(any) error {
	return nil
}
func (*mockWriteServer) RecvMsg(any) error {
	return nil
}

func TestPluginSync(t *testing.T) {
	ctx := context.Background()
	s := Server{
		Plugin: plugin.NewPlugin("test", "development", memdb.NewMemDBClient),
	}

	_, err := s.Init(ctx, &pb.Init_Request{})
	if err != nil {
		t.Fatal(err)
	}

	streamSyncServer := &mockSyncServer{}
	if err := s.Sync(&pb.Sync_Request{}, streamSyncServer); err != nil {
		t.Fatal(err)
	}
	if len(streamSyncServer.messages) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(streamSyncServer.messages))
	}
	writeMockServer := &mockWriteServer{}

	table := &schema.Table{
		Name: "test",
		Columns: []schema.Column{
			{
				Name: "test",
				Type: arrow.BinaryTypes.String,
			},
		},
	}
	sc := table.ToArrowSchema()
	b, err := pb.SchemaToBytes(sc)
	if err != nil {
		t.Fatal(err)
	}
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	bldr.Field(0).(*array.StringBuilder).Append("test")
	record := bldr.NewRecordBatch()
	recordBytes, err := pb.RecordToBytes(record)
	if err != nil {
		t.Fatal(err)
	}

	writeMockServer.messages = []*pb.Write_Request{
		{
			Message: &pb.Write_Request_MigrateTable{
				MigrateTable: &pb.Write_MessageMigrateTable{
					Table: b,
				},
			},
		},
		{
			Message: &pb.Write_Request_Insert{
				Insert: &pb.Write_MessageInsert{
					Record: recordBytes,
				},
			},
		},
	}

	if err := s.Write(writeMockServer); err != nil {
		t.Fatal(err)
	}

	streamSyncServer = &mockSyncServer{}
	if err := s.Sync(&pb.Sync_Request{
		Tables: []string{"*"},
	}, streamSyncServer); err != nil {
		t.Fatal(err)
	}
	if len(streamSyncServer.messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(streamSyncServer.messages))
	}

	if _, err := s.Close(ctx, &pb.Close_Request{}); err != nil {
		t.Fatal(err)
	}
}

func TestTransformSchema(t *testing.T) {
	ctx := context.Background()
	s := Server{
		Plugin: plugin.NewPlugin("test", "development", getColumnAdderPlugin()),
	}

	_, err := s.Init(ctx, &pb.Init_Request{})
	if err != nil {
		t.Fatal(err)
	}

	table := &schema.Table{
		Name: "test",
		Columns: []schema.Column{
			{
				Name: "test",
				Type: arrow.BinaryTypes.String,
			},
		},
	}
	sc := table.ToArrowSchema()

	schemaBytes, err := pb.SchemaToBytes(sc)
	require.NoError(t, err)

	resp, err := s.TransformSchema(ctx, &pb.TransformSchema_Request{Schema: schemaBytes})
	if err != nil {
		t.Fatal(err)
	}

	newSchema, err := pb.NewSchemaFromBytes(resp.Schema)
	require.NoError(t, err)

	require.Len(t, newSchema.Fields(), 2)
	require.Equal(t, "test", newSchema.Fields()[0].Name)
	require.Equal(t, "source", newSchema.Fields()[1].Name)
	require.Equal(t, "utf8", newSchema.Fields()[1].Type.(*arrow.StringType).Name())

	if _, err := s.Close(ctx, &pb.Close_Request{}); err != nil {
		t.Fatal(err)
	}
}

type mockSourceColumnAdderPluginClient struct {
	plugin.UnimplementedDestination
	plugin.UnimplementedSource
}

func getColumnAdderPlugin(...plugin.Option) plugin.NewClientFunc {
	c := &mockSourceColumnAdderPluginClient{}
	return func(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
		return c, nil
	}
}

func (*mockSourceColumnAdderPluginClient) Transform(context.Context, <-chan arrow.RecordBatch, chan<- arrow.RecordBatch) error {
	return nil
}
func (*mockSourceColumnAdderPluginClient) TransformSchema(_ context.Context, old *arrow.Schema) (*arrow.Schema, error) {
	return old.AddField(1, arrow.Field{Name: "source", Type: arrow.BinaryTypes.String})
}
func (*mockSourceColumnAdderPluginClient) Close(context.Context) error { return nil }

type testTransformPluginClient struct {
	plugin.UnimplementedDestination
	plugin.UnimplementedSource
	recordsSent int32
}

func (c *testTransformPluginClient) Transform(ctx context.Context, recvRecords <-chan arrow.RecordBatch, sendRecords chan<- arrow.RecordBatch) error {
	for record := range recvRecords {
		select {
		default:
			time.Sleep(1 * time.Second)
			sendRecords <- record
			atomic.AddInt32(&c.recordsSent, 1)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (*testTransformPluginClient) TransformSchema(_ context.Context, old *arrow.Schema) (*arrow.Schema, error) {
	return old, nil
}

func (*testTransformPluginClient) Close(context.Context) error {
	return nil
}

func TestTransformNoDeadlockOnSendError(t *testing.T) {
	client := &testTransformPluginClient{}
	p := plugin.NewPlugin("test", "development", func(context.Context, zerolog.Logger, []byte, plugin.NewClientOptions) (plugin.Client, error) {
		return client, nil
	})
	s := Server{
		Plugin: p,
	}
	_, err := s.Init(context.Background(), &pb.Init_Request{})
	require.NoError(t, err)

	// Create a channel to signal when Send was called
	sendCalled := make(chan struct{})
	// Create a channel to signal when we should return from the test
	done := make(chan struct{})
	defer close(done)

	stream := &mockTransformServerWithBlockingSend{
		incomingMessages: makeRequests(3), // Multiple messages to ensure Transform tries to keep sending
		sendCalled:       sendCalled,
		done:             done,
	}

	// Run Transform in a goroutine with a timeout
	errCh := make(chan error)
	go func() {
		errCh <- s.Transform(stream)
	}()

	// Wait for the first Send to be called
	select {
	case <-sendCalled:
		// Send was called, good
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for Send to be called")
	}

	// Now wait for Transform to complete or timeout
	select {
	case err := <-errCh:
		require.Error(t, err)
		// Check for either the simulated error or context cancellation
		if !strings.Contains(err.Error(), "simulated stream send error") &&
			!strings.Contains(err.Error(), "context canceled") {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Transform got deadlocked")
	}
}

type mockTransformServerWithBlockingSend struct {
	grpc.ServerStream
	incomingMessages []*pb.Transform_Request
	sendCalled       chan struct{}
	done             chan struct{}
	sendCount        int32
}

func (s *mockTransformServerWithBlockingSend) Recv() (*pb.Transform_Request, error) {
	if len(s.incomingMessages) > 0 {
		msg := s.incomingMessages[0]
		s.incomingMessages = s.incomingMessages[1:]
		return msg, nil
	}
	return nil, io.EOF
}

func (s *mockTransformServerWithBlockingSend) Send(*pb.Transform_Response) error {
	// Signal that Send was called
	select {
	case s.sendCalled <- struct{}{}:
	default:
	}

	// Return error on first send
	if atomic.AddInt32(&s.sendCount, 1) == 1 {
		return errors.New("simulated stream send error")
	}

	// Block until test is done
	<-s.done
	return nil
}

func (*mockTransformServerWithBlockingSend) Context() context.Context {
	return context.Background()
}

func makeRequests(i int) []*pb.Transform_Request {
	requests := make([]*pb.Transform_Request, i)
	for i := range i {
		requests[i] = makeRequestFromString("test")
	}
	return requests
}

func makeRequestFromString(s string) *pb.Transform_Request {
	record := makeRecordFromString(s)
	bs, _ := pb.RecordToBytes(record)
	return &pb.Transform_Request{Record: bs}
}

func makeRecordFromString(s string) arrow.RecordBatch {
	str := array.NewStringBuilder(memory.DefaultAllocator)
	str.AppendString(s)
	arr := str.NewStringArray()
	sch := arrow.NewSchema([]arrow.Field{{Name: "col1", Type: arrow.BinaryTypes.String}}, nil)

	return array.NewRecordBatch(sch, []arrow.Array{arr}, 1)
}
