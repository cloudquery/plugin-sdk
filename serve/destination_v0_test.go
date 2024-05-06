package serve

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/apache/arrow/go/v16/arrow/array"
	"github.com/apache/arrow/go/v16/arrow/memory"
	pbBase "github.com/cloudquery/plugin-pb-go/pb/base/v0"
	pb "github.com/cloudquery/plugin-pb-go/pb/destination/v0"
	"github.com/cloudquery/plugin-pb-go/specs"
	schemav2 "github.com/cloudquery/plugin-sdk/v2/schema"
	"github.com/cloudquery/plugin-sdk/v2/testdata"
	"github.com/cloudquery/plugin-sdk/v4/internal/deprecated"
	"github.com/cloudquery/plugin-sdk/v4/internal/memdb"
	serversDestination "github.com/cloudquery/plugin-sdk/v4/internal/servers/destination/v0"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDestination(t *testing.T) {
	p := plugin.NewPlugin("testDestinationPlugin", "development", memdb.NewMemDBClient)
	srv := Plugin(p, WithArgs("serve"), WithDestinationV0V1Server(), WithTestListener())
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
	conn, err := grpc.DialContext(ctx, "bufnet1", grpc.WithContextDialer(srv.bufPluginDialer), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	c := pb.NewDestinationClient(conn)
	spec := specs.Destination{
		WriteMode: specs.WriteModeAppend,
	}
	specBytes, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := c.Configure(ctx, &pbBase.Configure_Request{Config: specBytes}); err != nil {
		t.Fatal(err)
	}
	getNameRes, err := c.GetName(ctx, &pbBase.GetName_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getNameRes.Name != "testDestinationPlugin" {
		t.Fatalf("expected name to be testDestinationPlugin but got %s", getNameRes.Name)
	}

	getVersionRes, err := c.GetVersion(ctx, &pbBase.GetVersion_Request{})
	if err != nil {
		t.Fatal(err)
	}
	if getVersionRes.Version != "development" {
		t.Fatalf("expected version to be development but got %s", getVersionRes.Version)
	}

	tableName := "test_destination_serve"
	sourceName := "test_destination_serve_source"
	syncTime := time.Now()
	tableV2 := testdata.TestTable(tableName)
	tablesV2 := schemav2.Tables{tableV2}
	sourceSpec := specs.Source{
		Name: sourceName,
	}
	tablesV2Bytes, err := json.Marshal(tablesV2)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.Migrate(ctx, &pb.Migrate_Request{
		Tables: tablesV2Bytes,
	}); err != nil {
		t.Fatal(err)
	}

	destResource := schemav2.DestinationResource{
		TableName: tableName,
		Data:      deprecated.GenTestData(tableV2),
	}
	_ = destResource.Data[0].Set(sourceName)
	_ = destResource.Data[1].Set(syncTime)
	destResourceBytes, err := json.Marshal(destResource)
	if err != nil {
		t.Fatal(err)
	}
	sourceSpecBytes, err := json.Marshal(sourceSpec)
	if err != nil {
		t.Fatal(err)
	}
	writeClient, err := c.Write2(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeClient.Send(&pb.Write2_Request{
		SourceSpec: sourceSpecBytes,
		Source:     sourceSpec.Name,
		Timestamp:  timestamppb.New(syncTime.Truncate(time.Microsecond)),
		Tables:     tablesV2Bytes,
	}); err != nil {
		t.Fatal(err)
	}
	if err := writeClient.Send(&pb.Write2_Request{
		Resource: destResourceBytes,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := writeClient.CloseAndRecv(); err != nil {
		t.Fatal(err)
	}

	// serversDestination
	table := serversDestination.TableV2ToV3(tableV2)
	msgs, err := p.SyncAll(ctx, plugin.SyncOptions{
		Tables: []string{tableName},
	})
	if err != nil {
		t.Fatal(err)
	}
	totalResources := 0
	destRecord := serversDestination.CQTypesOneToRecord(memory.DefaultAllocator, destResource.Data, table.ToArrowSchema())
	for _, msg := range msgs {
		totalResources++
		m := msg.(*message.SyncInsert)
		if !array.RecordEqual(destRecord, m.Record) {
			// diff := destination.RecordDiff(destRecord, resource)
			t.Fatalf("expected %v but got %v", destRecord, m.Record)
		}
	}
	if totalResources != 1 {
		t.Fatalf("expected 1 resource but got %d", totalResources)
	}

	if _, err := c.DeleteStale(ctx, &pb.DeleteStale_Request{
		Source:    "testSource",
		Timestamp: timestamppb.New(time.Now().Truncate(time.Microsecond)),
		Tables:    tablesV2Bytes,
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := c.Close(ctx, &pb.Close_Request{}); err != nil {
		t.Fatalf("failed to call Close: %v", err)
	}
	cancel()
	wg.Wait()
	if serverErr != nil {
		t.Fatal(serverErr)
	}
}
