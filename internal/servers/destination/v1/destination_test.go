package destination

import (
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/destination/v1"
	pbSource "github.com/cloudquery/plugin-pb-go/pb/source/v2"
	"github.com/cloudquery/plugin-pb-go/specs"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
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
	destinationSpec := specs.Destination{
		Name: "test",
	}
	destinationSpecBytes, err := json.Marshal(destinationSpec)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Configure(ctx, &pb.Configure_Request{
		Config: destinationSpecBytes,
	})
	if err != nil {
		t.Fatal(err)
	}

	writeMockServer := &mockWriteServer{}
	if err := s.Write(writeMockServer); err != nil {
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
	schemas := schema.Tables{table}.ToArrowSchemas()
	schemaBytes, err := pbSource.SchemasToBytes(schemas)
	if err != nil {
		t.Fatal(err)
	}
	sc := table.ToArrowSchema()
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	bldr.Field(0).(*array.StringBuilder).Append("test")
	record := bldr.NewRecord()
	recordBytes, err := pbSource.RecordToBytes(record)
	if err != nil {
		t.Fatal(err)
	}

	sourceSpec := specs.Source{
		Name: "source_test",
	}
	sourceSpecBytes, err := json.Marshal(sourceSpec)
	if err != nil {
		t.Fatal(err)
	}

	writeMockServer.messages = []*pb.Write_Request{
		{
			Tables:     schemaBytes,
			Resource:   recordBytes,
			SourceSpec: sourceSpecBytes,
		},
	}
	if err := s.Write(writeMockServer); err != nil {
		t.Fatal(err)
	}

	if _, err := s.Close(ctx, &pb.Close_Request{}); err != nil {
		t.Fatal(err)
	}
}
