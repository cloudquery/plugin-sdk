package plugin

import (
	"context"
	"io"
	"testing"

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
	record := bldr.NewRecord()
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

func (*mockSourceColumnAdderPluginClient) Transform(context.Context, <-chan arrow.Record, chan<- arrow.Record) error {
	return nil
}
func (*mockSourceColumnAdderPluginClient) TransformSchema(_ context.Context, old *arrow.Schema) (*arrow.Schema, error) {
	return old.AddField(1, arrow.Field{Name: "source", Type: arrow.BinaryTypes.String})
}
func (*mockSourceColumnAdderPluginClient) Close(context.Context) error { return nil }
