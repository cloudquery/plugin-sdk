package serve

import (
	"bytes"
	"context"
	"io"
	"sync"
	"testing"

	"github.com/apache/arrow/go/v16/arrow"
	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/ipc"
	"github.com/apache/arrow/go/v16/arrow/memory"
	pb "github.com/cloudquery/plugin-pb-go/pb/plugin/v3"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestPluginServe(t *testing.T) {
	p := plugin.NewPlugin(
		"testPluginV3",
		"v1.0.0",
		memdb.NewMemDBClient)
	srv := Plugin(p, WithArgs("serve"), WithTestListener())
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	var serverErr error
	go func() {
		defer wg.Done()
		serverErr = srv.Serve(ctx)
	}()
	defer func() {
		cancel()
		wg.Wait()
	}()

	// https://stackoverflow.com/questions/42102496/testing-a-grpc-service
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(srv.bufPluginDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	c := pb.NewPluginClient(conn)

	getNameRes, err := c.GetName(ctx, &pb.GetName_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getNameRes.Name != "testPluginV3" {
		t.Fatalf("expected name to be testPluginV3 but got %s", getNameRes.Name)
	}

	getVersionResponse, err := c.GetVersion(ctx, &pb.GetVersion_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getVersionResponse.Version != "v1.0.0" {
		t.Fatalf("Expected version to be v1.0.0 but got %s", getVersionResponse.Version)
	}

	if _, err := c.Init(ctx, &pb.Init_Request{}); err != nil {
		t.Fatal(err)
	}

	getTablesRes, err := c.GetTables(ctx, &pb.GetTables_Request{Tables: []string{"*"}})
	if err != nil {
		t.Fatal(err)
	}
	schemas, err := pb.NewSchemasFromBytes(getTablesRes.Tables)
	if err != nil {
		t.Fatal(err)
	}
	tables, err := schema.NewTablesFromArrowSchemas(schemas)
	if err != nil {
		t.Fatal(err)
	}

	if len(tables) != 3 {
		t.Fatalf("Expected 2 tables but got %d", len(tables))
	}
	testTable := schema.Table{
		Name: "test_table",
		Columns: []schema.Column{
			{
				Name: "col1",
				Type: arrow.BinaryTypes.String,
			},
		},
	}
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, testTable.ToArrowSchema())
	bldr.Field(0).(*array.StringBuilder).Append("test")
	record := bldr.NewRecord()

	recordBytes, err := pb.RecordToBytes(record)
	if err != nil {
		t.Fatal(err)
	}
	sc := testTable.ToArrowSchema()
	tableBytes, err := pb.SchemaToBytes(sc)
	if err != nil {
		t.Fatal(err)
	}
	writeClient, err := c.Write(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err := writeClient.Send(&pb.Write_Request{
		Message: &pb.Write_Request_MigrateTable{
			MigrateTable: &pb.Write_MessageMigrateTable{
				Table: tableBytes,
			},
		},
	}); err != nil {
		t.Fatal(err)
	}
	if err := writeClient.Send(&pb.Write_Request{
		Message: &pb.Write_Request_Insert{
			Insert: &pb.Write_MessageInsert{
				Record: recordBytes,
			},
		},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := writeClient.CloseAndRecv(); err != nil {
		t.Fatal(err)
	}

	syncClient, err := c.Sync(ctx, &pb.Sync_Request{
		Tables: []string{"test_table"},
	})
	if err != nil {
		t.Fatal(err)
	}
	var resources []arrow.Record
	for {
		r, err := syncClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		m := r.Message.(*pb.Sync_Response_Insert)
		rdr, err := ipc.NewReader(bytes.NewReader(m.Insert.Record))
		if err != nil {
			t.Fatal(err)
		}
		for rdr.Next() {
			rec := rdr.Record()
			rec.Retain()
			resources = append(resources, rec)
		}
	}

	totalResources := 0
	for _, resource := range resources {
		sc := resource.Schema()
		tableName, ok := sc.Metadata().GetValue(schema.MetadataTableName)
		if !ok {
			t.Fatal("Expected table name metadata to be set")
		}
		if tableName != "test_table" {
			t.Fatalf("Expected resource with table name test_table. got: %s", tableName)
		}
		if len(resource.Columns()) != 1 {
			t.Fatalf("Expected resource with data length 1 but got %d", len(resource.Columns()))
		}
		totalResources++
	}
	if totalResources != 1 {
		t.Fatalf("Expected 1 resource on channel but got %d", totalResources)
	}

	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}
